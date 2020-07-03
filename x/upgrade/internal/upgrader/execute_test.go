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
			"\"windows/amd64\":\"https://github.com/Bococoin/core/releases/download/1.0.0/bocod.exe?checksum=sha256:d47042ebad23df6c170bf8aaef823b4da71e2859dda9d7bb63b63f23c115a786\"," +
			"\"linux/amd64\":\"https://github.com/Bococoin/core/releases/download/1.0.0/bocod?checksum=sha256:75a9725e73ee25d236f8ab8efc5f17c1085d24bd720e03f4bbf478ae0317d11f\"" +
			"}}",
	}
}
