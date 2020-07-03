package upgrade

// nolint

import (
	"github.com/Bococoin/core/x/upgrade/internal/keeper"
	"github.com/Bococoin/core/x/upgrade/internal/types"
	"github.com/Bococoin/core/x/upgrade/internal/upgrader"
)

const (
	ModuleName                        = types.ModuleName
	RouterKey                         = types.RouterKey
	StoreKey                          = types.StoreKey
	QuerierKey                        = types.QuerierKey
	PlanByte                          = types.PlanByte
	DoneByte                          = types.DoneByte
	ProposalTypeSoftwareUpgrade       = types.ProposalTypeSoftwareUpgrade
	ProposalTypeCancelSoftwareUpgrade = types.ProposalTypeCancelSoftwareUpgrade
	QueryCurrent                      = types.QueryCurrent
	QueryApplied                      = types.QueryApplied
)

var (
	// functions aliases
	RegisterCodec                    = types.RegisterCodec
	PlanKey                          = types.PlanKey
	NewSoftwareUpgradeProposal       = types.NewSoftwareUpgradeProposal
	NewCancelSoftwareUpgradeProposal = types.NewCancelSoftwareUpgradeProposal
	NewQueryAppliedParams            = types.NewQueryAppliedParams
	NewKeeper                        = keeper.NewKeeper
	NewQuerier                       = keeper.NewQuerier
	UpgradeRestartCmd                = upgrader.UpgradeRestartCmd
)

type (
	UpgradeHandler                = types.UpgradeHandler
	Plan                          = types.Plan
	SoftwareUpgradeProposal       = types.SoftwareUpgradeProposal
	CancelSoftwareUpgradeProposal = types.CancelSoftwareUpgradeProposal
	QueryAppliedParams            = types.QueryAppliedParams
	Keeper                        = keeper.Keeper
)
