package bocomint

// nolint

import (
	"github.com/Bococoin/core/x/bocomint/internal/keeper"
	"github.com/Bococoin/core/x/bocomint/internal/types"
)

const (
	ModuleName        = types.ModuleName
	DefaultParamspace = types.DefaultParamspace
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
	QueryParameters   = types.QueryParameters
	QueryBlockReward  = types.QueryBlockReward
)

var (
	// functions aliases
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	ParamKeyTable       = types.ParamKeyTable
	NewParams           = types.NewParams
	DefaultParams       = types.DefaultParams

	// variable aliases
	ModuleCdc    = types.ModuleCdc
	MinterKey    = types.MinterKey
	KeyMintDenom = types.KeyMintDenom
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params
)
