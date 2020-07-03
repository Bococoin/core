package upgrader

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/Bococoin/core/x/upgrade/internal/types"
	"github.com/hashicorp/go-getter"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func ProcessUpgrade(ctx sdk.Context, plan types.Plan) error {
	logger := ctx.Logger()
	logger.Info(fmt.Sprintf("\n\n--------------------------------------------------\n"+
		"Starting 'bocod' binary upgrade.\n"+
		"    Upgrade name: %s\n"+
		"    Upgrade height: %d\n"+
		"    Upgrade info: %s\n", plan.Name, plan.Height, plan.Info),
	)
	if ctx.BlockHeight() != plan.Height {
		return errors.New("unexpected BocoUpgradeHandler call, current block height is not equals to planned")

	}

	config, err := GetDownloadConfig(&plan)
	if err != nil {
		return errors.Wrap(err, "error while getting download configuration from plan")

	}

	newPath, err := DownloadNewBinary(config)
	if err != nil {
		return errors.Wrap(err, "error while update bocod binary")

	}

	logger.Info("New binary downloaded successfully. Starting....\n\n\n")
	err = StartProcess(newPath, "upg_restart")
	if err != nil {
		return errors.Wrap(err, "error while restarting")
	}
	return nil
}

// UpdateBinary will grab the binary and replace self with it
func DownloadNewBinary(config *UpgradeConfig) (string, error) {
	url, err := GetDownloadURL(config)
	if err != nil {
		return "", err
	}

	ver, err := GetNewVersion(config)
	if err != nil {
		return "", err
	}

	binPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	updateDir := filepath.Dir(binPath)
	oldname := filepath.Base(binPath)
	filename := ""
	if runtime.GOOS == "windows" {
		filename = fmt.Sprintf("%s_%s.exe", strings.ReplaceAll(oldname, ".exe", ""), ver)
	} else {
		filename = fmt.Sprintf("%s_%s", oldname, ver)
	}

	// Copy the contents of newbinary to a new executable file
	newBinPath := filepath.Join(updateDir, filename)

	err = getter.GetFile(newBinPath, url)
	if err != nil {
		return "", err
	}

	err = EnsureBinary(newBinPath)
	if err != nil {
		return "", err
	}

	// this is where we'll move the executable to so that we can swap in the updated replacement
	oldPath := filepath.Join(updateDir, fmt.Sprintf("%s.%s.old", oldname, version.Version))

	// move the existing executable to a new file in the same directory
	err = os.Rename(binPath, oldPath)
	if err != nil {
		return "", err
	}

	// if it is successful, let's ensure the binary is executable
	MarkExecutable(newBinPath)
	if err != nil {
		return "", err
	}
	return newBinPath, nil
}

// MarkExecutable will try to set the executable bits if not already set
// Fails if file doesn't exist or we cannot set those bits
func MarkExecutable(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return errors.Wrap(err, "stating binary")
	}
	// end early if world exec already set
	if info.Mode()&0001 == 1 {
		return nil
	}
	// now try to set all exec bits
	newMode := info.Mode().Perm() | 0111
	return os.Chmod(path, newMode)
}

// UpgradeConfig is expected format for the info field to allow auto-download
type UpgradeConfig struct {
	Version  string            `json:"version"`
	Binaries map[string]string `json:"binaries"`
}

func GetDownloadConfig(info *types.Plan) (*UpgradeConfig, error) {
	doc := strings.TrimSpace(info.Info)

	// check if it is the update config
	var config UpgradeConfig
	err := json.Unmarshal([]byte(doc), &config)
	if err == nil {
		return &config, nil
	} else {
		return nil, err
	}
}

// GetDownloadURL will check if there is an arch-dependent binary specified in Info
func GetDownloadURL(config *UpgradeConfig) (daemon string, err error) {

	// check if it is the update config
	daemon, ok := config.Binaries[osArch()]
	if !ok {
		return "", errors.Errorf("cannot find daemon binary for os/arch: %s", osArch())
	}

	return
}

func GetNewVersion(config *UpgradeConfig) (version string, err error) {
	if len(config.Version) == 0 {
		return "", errors.New("version is not defined")
	}
	return config.Version, nil
}

func osArch() string {
	return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
}

// EnsureBinary ensures the file exists and is executable, or returns an error
func EnsureBinary(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return errors.Wrap(err, "cannot stat home dir")
	}
	if !info.Mode().IsRegular() {
		return errors.Errorf("%s is not a regular file", info.Name())
	}
	if runtime.GOOS != "windows" {
		// this checks if the world-executable bit is set (we cannot check owner easily)
		exec := info.Mode().Perm() & 0001
		if exec == 0 {
			return errors.Errorf("%s is not world executable", info.Name())
		}
	}
	return nil
}
