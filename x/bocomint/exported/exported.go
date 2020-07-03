package exported

import sdk "github.com/cosmos/cosmos-sdk/types"

type AccountRewardsInfo interface {
	GetAddress() sdk.AccAddress
	GetSelfPercent() sdk.Dec
}
