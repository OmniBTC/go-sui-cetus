package types

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
