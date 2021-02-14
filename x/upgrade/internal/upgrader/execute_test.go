package upgrader

import (
	"fmt"
	"github.com/Bococoin/core/x/upgrade/internal/types"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const (
	BocodUpgradeName = "bocod_upgrade"
)

func TestGetDownloadURL(t *testing.T) {
	plan := GetTestUpgradePlan()

	config, err := GetDownloadConfig(&plan)
	require.Nil(t, err)

	url, err := GetDownloadURL(config)

	require.Nil(t, err)
	fmt.Print(url)
}

func TestDoUpdateBinary(t *testing.T) {
	t.Skip("Skipping UpdateBinary testing")
	plan := GetTestUpgradePlan()

	config, err := GetDownloadConfig(&plan)
	require.Nil(t, err)

	newPath, err := DownloadNewBinary(config)
	require.Nil(t, err)
	require.NotNil(t, newPath)

	fmt.Print(newPath)
}

func TestEnsureBinary(t *testing.T) {
	binPath, err := os.Executable()
	require.Nil(t, err)

	err = EnsureBinary(binPath)
	require.Nil(t, err)
}

//TODO: bring it back to update_test
func GetTestUpgradePlan() types.Plan {
	return types.Plan{
		Name:   BocodUpgradeName,
		Height: 50,
		Info: "{" +
			"\"version\":\"1.0.0\"," +
			"\"binaries\":{" +
			"\"windows/amd64\":\"https://github.com/Bococoin/core/releases/download/1.0/bocod_wnd.zip?checksum=sha256:2cf51db906046e3e6743f4cd53b3029fd76ea4a8a5019a7106670fc82b9cc9a1\"," +
			"\"linux/amd64\":\"https://github.com/Bococoin/core/releases/download/1.0/bocod.zip?checksum=sha256:f53f49b1f0a2f53da19620ac2555077ac195b8c899d6c3a53644190b38bead11\"" +
			"}}",
	}
}
