package types

import (
	"fmt"
	boco "github.com/Bococoin/core/types"
	"math"
)

type PeriodEntry struct {
	S int64 `json:"start" yaml:"start"`
	E int64 `json:"end" yaml:"end"`
	M int64 `json:"val" yaml:"val"`
}

type PeriodTable []PeriodEntry

func NewPeriodEntry(start int64, end int64, mult int64) *PeriodEntry {
	return &PeriodEntry{S: start, E: end, M: mult}
}

func (table PeriodTable) GetEntry(place int64) (*PeriodEntry, error) {
	if place < 0 {
		return nil, fmt.Errorf("place must be positive, got %d", place)
	}
	for _, entry := range table {
		if entry.CompareVal(place) == 0 {
			return &entry, nil
		}
	}
	return nil, fmt.Errorf("place %d is not in table", place)
}
func (table PeriodTable) GetValue(place int64) (int64, error) {
	if place < 0 {
		return 0, fmt.Errorf("place must be positive, got %d", place)
	}
	for _, entry := range table {
		if entry.CompareVal(place) == 0 {
			return entry.M, nil
		}
	}
	return 0, fmt.Errorf("place %d is not in table", place)
}

func (table PeriodTable) String() string {
	res := "\n"
	for _, entry := range table {
		res += fmt.Sprintln("    " + entry.String())
	}
	return res
}

func (mult PeriodEntry) String() string {
	return fmt.Sprintf("%10d%s - %10d%s: %d", mult.S/boco.OneCoin, boco.Coin, mult.E/boco.OneCoin, boco.Coin, mult.M)
}
func (mult PeriodEntry) CompareVal(val int64) int8 {
	if val < mult.S {
		return -1
	} else if val >= mult.S && val <= mult.E {
		return 0
	} else {
		return 1
	}
}

func DefaultRewardTable() PeriodTable {
	return PeriodTable{
		{
			S: 0,
			E: boco.BlocksPerQuarter*2 - 1,
			M: boco.DefaultMintStartValue,
		},
		{
			S: boco.BlocksPerQuarter * 2,
			E: boco.BlocksPerQuarter*3 - 1,
			M: boco.OneCoin * 5.2,
		},
		{
			S: boco.BlocksPerQuarter * 3,
			E: boco.BlocksPerQuarter*4 - 1,
			M: boco.OneCoin * 6.76,
		},
		{
			S: boco.BlocksPerQuarter * 4,
			E: boco.BlocksPerQuarter*5 - 1,
			M: boco.OneCoin * 8.788,
		},
		{
			S: boco.BlocksPerQuarter * 5,
			E: boco.BlocksPerQuarter*6 - 1,
			M: boco.OneCoin * 10.5456,
		},
		{
			S: boco.BlocksPerQuarter * 6,
			E: boco.BlocksPerQuarter*7 - 1,
			M: boco.OneCoin * 12.65472,
		},
		{
			S: boco.BlocksPerQuarter * 7,
			E: boco.BlocksPerQuarter*8 - 1,
			M: boco.OneCoin * 10.123776,
		},
		{
			S: boco.BlocksPerQuarter * 8,
			E: boco.BlocksPerQuarter*9 - 1,
			M: boco.OneCoin * 8.099,
		},
		{
			S: boco.BlocksPerQuarter * 9,
			E: boco.BlocksPerQuarter*10 - 1,
			M: boco.OneCoin * 5.669314,
		},
		{
			S: boco.BlocksPerQuarter * 10,
			E: boco.BlocksPerQuarter*11 - 1,
			M: boco.OneCoin * 3.96852,
		},
		{
			S: boco.BlocksPerQuarter * 11,
			E: boco.BlocksPerQuarter*12 - 1,
			M: boco.OneCoin * 2.77964,
		},
		{
			S: boco.BlocksPerQuarter * 12,
			E: boco.BlocksPerQuarter*13 - 1,
			M: boco.OneCoin * 1.944574,
		},
		{
			S: boco.BlocksPerQuarter * 13,
			E: math.MaxInt64,
			M: boco.OneCoin,
		},
	}
}
