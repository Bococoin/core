package upgrader

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// UpgradeRestartCmd returns a command that
func UpgradeRestartCmd(ctx *server.Context) *cobra.Command { // nolint: golint
	cmd := &cobra.Command{
		Use:    "upg_restart",
		Short:  "Restart new binary version after automatic upgrade download.",
		Long:   `Restart new binary version after automatic upgrade download.`,
		Hidden: true,
		RunE: func(_ *cobra.Command, args []string) error {
			ctx.Logger.Info(fmt.Sprintf("New version %s.%s has been started for self rename", version.ServerName, version.Version))

			binPath, err := os.Executable()
			if err != nil {
				return err
			}

			updateDir := filepath.Dir(binPath)
			binname := filepath.Base(binPath)
			if strings.Count(binname, ".") <= 1 {
				return errors.Errorf("upg_restart parameter used only for restarting during automatic upgrade process")
			}

			newPath := ""
			if runtime.GOOS == "windows" {
				newPath = filepath.Join(updateDir, fmt.Sprintf("%s.exe", version.ServerName))
			} else {
				newPath = filepath.Join(updateDir, version.ServerName)
			}

			// move the existing executable to a new file in the same directory
			err = os.Rename(binPath, newPath)
			if err != nil {
				return err
			}

			ctx.Logger.Info("Rename completed. Starting for work")
			err = StartProcess(newPath, "start")
			return err
		},
	}

	return cmd
}
