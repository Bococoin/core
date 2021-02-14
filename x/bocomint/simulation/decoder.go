package simulation

import (
	"bytes"
	"fmt"

	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/Bococoin/core/x/bocomint/internal/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding mint type
func DecodeStore(cdc *codec.Codec, kvA, kvB tmkv.Pair) string {
	switch {
	case bytes.Equal(kvA.Key, types.MinterKey):
		return fmt.Sprintf("minter")
	default:
		panic(fmt.Sprintf("invalid mint key %X", kvA.Key))
	}
}
