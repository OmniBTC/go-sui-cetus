package types

import (
	"errors"
	"math/big"
)

type BigNumber interface{}
type SuiResource interface{}
type SuiTxArg interface{}

const (
	CLOCK_ADDRESS         = "0x0000000000000000000000000000000000000000000000000000000000000006"
	ClmmIntegrateModule   = "pool_script"
	ClmmFetcherModule     = "fetcher_script"
	CoinInfoAddress       = "0x1::coin::CoinInfo"
	CoinStoreAddress      = "0x1::coin::CoinStore"
	PoolLiquidityCoinType = "PoolLiquidityCoin"
)

type NFT struct {
	Creator     string
	Description string
	ImageURL    string
	Link        string
	Name        string
	ProjectURL  string
}

type SuiStructTag struct {
	FullAddress   string
	SourceAddress string
	Address       string
	Module        string
	Name          string
	TypeArguments []string
}

type SuiBasicTypes string

const (
	Address SuiBasicTypes = "address"
	Bool    SuiBasicTypes = "bool"
	U8      SuiBasicTypes = "u8"
	U16     SuiBasicTypes = "u16"
	U32     SuiBasicTypes = "u32"
	U64     SuiBasicTypes = "u64"
	U128    SuiBasicTypes = "u128"
	U256    SuiBasicTypes = "u256"
)

const (
	Object SuiBasicTypes = "object"
)

func getDefaultSuiInputType(value interface{}) (SuiBasicTypes, error) {
	switch v := value.(type) {
	case string:
		if len(v) >= 2 && v[:2] == "0x" {
			return Object, nil
		}
	case int, int32, int64, *big.Int:
		return U64, nil
	case bool:
		return Bool, nil
	default:
		return "", errors.New("Unknown type for value")
	}
	return "", errors.New("Unknown type for value")
}
