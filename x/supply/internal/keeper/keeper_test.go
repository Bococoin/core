package keeper_test

import (
	boco "github.com/Bococoin/core/types"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/Bococoin/core/x/supply/internal/types"
)

func TestSupply(t *testing.T) {
	initialPower := int64(100)
	initTokens := sdk.TokensFromConsensusPower(initialPower)

	app, ctx := createTestApp(false)
	totalSupply := sdk.NewCoins(sdk.NewCoin(boco.DefaultDenom, initTokens))
	app.SupplyKeeper.SetSupply(ctx, types.NewSupply(totalSupply))

	total := app.SupplyKeeper.GetSupply(ctx).GetTotal()

	require.Equal(t, totalSupply, total)
}

func TestValidatePermissions(t *testing.T) {
	app, _ := createTestApp(false)

	err := app.SupplyKeeper.ValidatePermissions(multiPermAcc)
	require.NoError(t, err)

	err = app.SupplyKeeper.ValidatePermissions(randomPermAcc)
	require.NoError(t, err)

	// unregistered permissions
	otherAcc := types.NewEmptyModuleAccount("other", "other")
	err = app.SupplyKeeper.ValidatePermissions(otherAcc)
	require.Error(t, err)
}
