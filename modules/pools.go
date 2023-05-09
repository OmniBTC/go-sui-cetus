package modules

import (
	"context"
	"strconv"

	"github.com/coming-chat/go-sui/client"
	suitypes "github.com/coming-chat/go-sui/types"
)

const (
	ClmmIntegrateModule = "pool_script"
	ClmmFetcherModule   = "fetcher_script"
)

type CreatePoolParams struct {
	TickSpacing         int
	InitializeSqrtPrice string
	URI                 string
	CoinPairType
}

type CreatePoolAddLiquidityParams struct {
	AmountA    int
	AmountB    int
	FixAmountA bool
	TickLower  int
	TickUpper  int
	CreatePoolParams
}

type FetchTickParams struct {
	PoolID string
	CoinPairType
}

type GetTickParams struct {
	Start []int
	Limit int
	FetchTickParams
}

type PoolModule struct {
	baseModule
}

func NewPoolModule(c *client.Client, config sdkParsedOptions) *PoolModule {
	return &PoolModule{
		baseModule: baseModule{
			c:      c,
			config: config,
		},
	}
}

func (m *PoolModule) FetchTicks(ctx context.Context, params FetchTickParams) ([]TickData, error) {
	ticks := []TickData{}
	start := []int{}
	limit := 512

	for {
		data, err := m.getTicks(ctx, GetTickParams{
			FetchTickParams: FetchTickParams{
				PoolID: params.PoolID,
				CoinPairType: CoinPairType{
					CoinTypeA: params.CoinTypeA,
					CoinTypeB: params.CoinTypeB,
				},
			},
			Start: start,
			Limit: limit,
		})
		if err != nil {
			return nil, err
		}

		ticks = append(ticks, data...)
		if len(data) < limit {
			break
		}
		start = []int{data[len(data)-1].Index}
	}

	return ticks, nil
}

func (m *PoolModule) getTicks(ctx context.Context, params GetTickParams) (ticks []TickData, err error) {
	typeArgs := []string{
		string(params.CoinTypeA), string(params.CoinTypeB),
	}
	args := []any{
		params.PoolID,
		params.Start,
		strconv.Itoa(params.Limit),
	}

	response, err := m.dryRun(ctx, *m.config.clmmRouter, ClmmFetcherModule, "fetch_ticks", typeArgs, args)
	if err != nil {
		return
	}

	err = parseEventWithContent(response, "FetchTicksResultEvent", func(event suitypes.SuiEvent) error {
		// todo
		return nil
	})
	return
}
