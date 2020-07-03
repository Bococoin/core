package types // noalias

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/Bococoin/core/x/auth/exported"
	"github.com/Bococoin/core/x/supply/exported"
)

// StakingKeeper defines the expected staking keeper
type StakingKeeper interface {
	StakingTokenSupply(ctx sdk.Context) sdk.Int
	BondedRatio(ctx sdk.Context) sdk.Dec
}

// SupplyKeeper defines the expected supply keeper
type SupplyKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	SetModuleAccount(sdk.Context, exported.ModuleAccountI)

	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	GetSupply(ctx sdk.Context) (supply exported.SupplyI)
}

type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) acc.Account
	SetAccount(ctx sdk.Context, acc acc.Account)
	GetAllAccounts(ctx sdk.Context) (accounts []acc.Account)
	IterateAccounts(ctx sdk.Context, cb func(account acc.Account) (stop bool))
}
