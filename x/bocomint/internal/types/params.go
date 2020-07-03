package types

import (
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	boco "github.com/Bococoin/core/types"
	"github.com/Bococoin/core/x/params"
)

// Parameter store keys
var (
	KeyMintDenom    = []byte("MintDenom")
	KeyRewardTable  = []byte("RewardTable")
	KeyMintInterval = []byte("MintInterval")
)

// mint parameters
type Params struct {
	MintDenom    string      `json:"mint_denom" yaml:"mint_denom"`       // type of coin to mint
	RewardTable  PeriodTable `json:"reward_table" yaml:"reward_table"`   // reward table by block height
	MintInterval int64       `json:"mint_interval" yaml:"mint_interval"` // interval in block when we mint coins
}

// ParamTable for minting module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(
	mintDenom string, rewardTable PeriodTable, mintInterval int64,
) Params {

	return Params{
		MintDenom:    mintDenom,
		RewardTable:  rewardTable,
		MintInterval: mintInterval,
	}
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		MintDenom:    boco.DefaultDenom,
		RewardTable:  DefaultRewardTable(),
		MintInterval: boco.DefaultMintingInterval,
	}
}

// validate params
func (p Params) Validate() error {
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validatePeriodTabel(p.RewardTable); err != nil {
		return err
	}
	if err := validateMintInterval(p.MintInterval); err != nil {
		return err
	}

	return nil
}

func (p Params) String() string {
	return fmt.Sprintf(`Minting Params:
  Mint Denom:             		%s
  Mint Interval:				%d
  Reward Table:					%s
`,
		p.MintDenom, p.MintInterval, p.RewardTable,
	)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyMintDenom, &p.MintDenom, validateMintDenom),
		params.NewParamSetPair(KeyRewardTable, &p.RewardTable, validatePeriodTabel),
		params.NewParamSetPair(KeyMintInterval, &p.MintInterval, validateMintInterval),
	}
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("mint denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}
func validatePeriodTabel(i interface{}) error {
	v, ok := i.(PeriodTable)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if len(v) == 0 {
		return errors.New("table can not be empty")
	}
	return nil
}
func validateMintInterval(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v <= 0 {
		return errors.New("minting interval must be positive")
	}
	return nil
}
