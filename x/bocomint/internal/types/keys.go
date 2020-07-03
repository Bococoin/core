package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// the one key to use for the keeper store
var (
	MinterKey             = []byte{0x00}
	AccountRewardsInfoKey = []byte{0x01} // key for account rewards struct

)

// nolint
const (
	// module name
	ModuleName = "bocomint"

	// default paramspace for params keeper
	DefaultParamspace = ModuleName

	// StoreKey is the default store key for mint
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the minting store.
	QuerierRoute = StoreKey

	// Query endpoints supported by the minting querier
	QueryParameters  = "parameters"
	QueryBlockReward = "block_reward"
)

// AddressStoreKey turn an address to key used to get it from the account store
func GetAccountRewardsInfoStoreKey(addr sdk.AccAddress) []byte {
	return append(AccountRewardsInfoKey, addr.Bytes()...)
}

func GetAccountRewardsInfoAddress(key []byte) (valAddr sdk.AccAddress) {
	addr := key[1:]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	return sdk.AccAddress(addr)
}
