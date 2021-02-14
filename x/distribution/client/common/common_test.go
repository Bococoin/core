package common

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Bococoin/core/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
)

func TestQueryDelegationRewardsAddrValidation(t *testing.T) {
	cdc := codec.New()
	ctx := context.NewCLIContext().WithCodec(cdc)
	type args struct {
		delAddr string
		valAddr string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"invalid delegator address", args{"invalid", ""}, nil, true},
		{"empty delegator address", args{"", ""}, nil, true},
		{"invalid validator address", args{"boco1z98klpu4r50nm4pm8hewcryq55288r922rg5v9", "invalid"}, nil, true},
		{"empty validator address", args{"boco1z98klpu4r50nm4pm8hewcryq55288r922rg5v9", ""}, nil, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := QueryDelegationRewards(ctx, "", tt.args.delAddr, tt.args.valAddr)
			require.True(t, err != nil, tt.wantErr)
		})
	}
}
