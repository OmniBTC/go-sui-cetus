package modules

import (
	"math/big"
	"time"

	"github.com/omnibtc/go-sui-cetus/types"
)

const (
	cacheTime10min = 10 * 60
)

var (
	tokensCache    []TokenInfo
	lastTokenCache uint64

	poolsCache    []PoolInfo
	lastPoolCache uint64
)

type PositionRewarder struct {
	GrowthInside string
	AmountOwed   string
}

type Position struct {
	PosObjectID         types.SuiObjectIdType
	Pool                types.SuiObjectIdType
	Type                types.SuiAddressType
	CoinTypeA           types.SuiAddressType
	CoinTypeB           types.SuiAddressType
	Index               int
	Liquidity           string
	TickLowerIndex      int
	TickUpperIndex      int
	FeeGrowthInsideA    string
	FeeOwedA            string
	FeeGrowthInsideB    string
	FeeOwedB            string
	RewardAmountOwed0   string
	RewardAmountOwed1   string
	RewardAmountOwed2   string
	RewardGrowthInside0 string
	RewardGrowthInside1 string
	RewardGrowthInside2 string
}

type CoinPairType struct {
	CoinTypeA types.SuiAddressType
	CoinTypeB types.SuiAddressType
}

type PoolImmutables struct {
	PoolAddress string
	TickSpacing string
	CoinPairType
}

type Pool struct {
	PoolType                string
	CoinAmountA             int
	CoinAmountB             int
	CurrentSqrtPrice        int64
	CurrentTickIndex        int
	FeeGrowthGlobalB        int64
	FeeGrowthGlobalA        int64
	FeeProtocolCoinA        int64
	FeeProtocolCoinB        int64
	FeeRate                 int64
	IsPause                 bool
	Liquidity               int64
	Index                   int
	PositionsHandle         string
	RewarderInfos           []Rewarder
	RewarderLastUpdatedTime string
	TicksHandle             string
	URI                     string
	Name                    string
	PoolImmutables
}

type Rewarder struct {
	CoinAddress        string
	EmissionsPerSecond int
	GrowthGlobal       int
	EmissionsEveryDay  int
}

type InitEvent struct {
	PoolsID        types.SuiObjectIdType
	GlobalConfigID types.SuiObjectIdType
	GlobalVaultID  types.SuiObjectIdType
}

type CreatePartnerEvent struct {
	Name         string
	Recipient    types.SuiAddressType
	PartnerID    types.SuiObjectIdType
	PartnerCapID types.SuiObjectIdType
	FeeRate      string
	StartEpoch   string
	EndEpoch     string
}

type FaucetEvent struct {
	ID   string
	Time int
}

type CoinAsset struct {
	CoinAddress  types.SuiAddressType
	CoinObjectId types.SuiObjectIdType
	Balance      *big.Int
}

type WarpSuiObject struct {
	CoinAddress types.SuiAddressType
	Balance     int
}

type FaucetCoin struct {
	TransactionModule string
	SuplyID           types.SuiObjectIdType
	Decimals          int
}

func getTokensCache() []TokenInfo {
	if lastTokenCache+uint64(cacheTime10min) < uint64(time.Now().Unix()) {
		return nil
	}
	return tokensCache
}

func setTokensCache(tokens []TokenInfo) {
	tokensCache = tokens
	lastTokenCache = uint64(time.Now().Unix())
}

func getPoolsCache() []PoolInfo {
	if lastPoolCache+uint64(cacheTime10min) < uint64(time.Now().Unix()) {
		return nil
	}
	return poolsCache
}

func setPoolsCache(pools []PoolInfo) {
	poolsCache = pools
	lastPoolCache = uint64(time.Now().Unix())
}
