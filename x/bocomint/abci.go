package bocomint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	boco "github.com/Bococoin/core/types"
	"github.com/Bococoin/core/x/bocomint/internal/types"
)

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, k Keeper) {
	// fetch stored minter & params

	//Stop minting coins when Max Emission reached
	if k.GetTotalSupply(ctx) >= boco.DefaultMaxEmission {
		return
	}

	logger := ctx.Logger()
	blockHeight := ctx.BlockHeight()
	params := k.GetParams(ctx)

	amt, err := params.RewardTable.GetValue(blockHeight)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	//mint coins for params.MintInterval blocks
	mintedCoin := sdk.NewCoin(params.MintDenom, sdk.NewInt(amt))
	mintedCoins := sdk.NewCoins(mintedCoin)

	// mint coins, update supply
	err = k.MintCoins(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	// send the minted coins to the fee collector account
	err = k.AddCollectedFees(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyBlockReward, string(amt)),
		),
	)
}
