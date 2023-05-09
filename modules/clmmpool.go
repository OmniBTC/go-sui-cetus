package modules

import (
	"math/big"
	"strconv"
)

type TickData struct {
	ObjectID               string
	Index                  int
	SqrtPrice              *big.Int
	LiquidityNet           *big.Int
	LiquidityGross         *big.Int
	FeeGrowthOutsideA      *big.Int
	FeeGrowthOutsideB      *big.Int
	RewardersGrowthOutside []*big.Int
}

type Tick struct {
	Index                  Bits
	SqrtPrice              string
	LiquidityNet           Bits
	LiquidityGross         string
	FeeGrowthOutsideA      string
	FeeGrowthOutsideB      string
	RewardersGrowthOutside [3]string
}

type Bits struct {
	Bits string
}

type ClmmpoolData struct {
	CoinA            string
	CoinB            string
	CurrentSqrtPrice *big.Int
	CurrentTickIndex int
	FeeGrowthGlobalA *big.Int
	FeeGrowthGlobalB *big.Int
	FeeProtocolCoinA *big.Int
	FeeProtocolCoinB *big.Int
	FeeRate          *big.Int
	Liquidity        *big.Int
	TickIndexes      []int
	TickSpacing      int
	Ticks            []TickData
	CollectionName   string
}

func TransClmmpoolDataWithoutTicks(pool Pool) (ClmmpoolData, error) {
	tickSpacing, err := strconv.Atoi(pool.TickSpacing)
	if err != nil {
		return ClmmpoolData{}, err
	}
	poolData := ClmmpoolData{
		CoinA:            string(pool.PoolImmutables.CoinTypeA),
		CoinB:            string(pool.PoolImmutables.CoinTypeB),
		CurrentSqrtPrice: new(big.Int).SetInt64(pool.CurrentSqrtPrice),
		CurrentTickIndex: pool.CurrentTickIndex,
		FeeGrowthGlobalA: new(big.Int).SetInt64(pool.FeeGrowthGlobalA),
		FeeGrowthGlobalB: new(big.Int).SetInt64(pool.FeeGrowthGlobalB),
		FeeProtocolCoinA: new(big.Int).SetInt64(pool.FeeProtocolCoinA),
		FeeProtocolCoinB: new(big.Int).SetInt64(pool.FeeProtocolCoinB),
		FeeRate:          new(big.Int).SetInt64(pool.FeeRate),
		Liquidity:        new(big.Int).SetInt64(pool.Liquidity),
		TickIndexes:      []int{},
		TickSpacing:      tickSpacing,
		Ticks:            []TickData{},
		CollectionName:   "",
	}
	return poolData, nil
}
