package modules

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/coming-chat/go-sui/client"
	suitypes "github.com/coming-chat/go-sui/types"
	"github.com/omnibtc/go-sui-cetus/types"
)

type TokenFullList struct {
	FullList struct {
		ValueList []TokenInfo `json:"value_list"`
	} `json:"full_list"`
}

type PoolFullList struct {
	FullList struct {
		ValueList []PoolInfo `json:"value_list"`
	} `json:"full_list"`
}

type TokenInfo struct {
	Name           string `json:"name"`
	Symbol         string `json:"symbol"`
	OfficialSymbol string `json:"official_symbol"`
	CoingeckoID    string `json:"coingecko_id"`
	Decimals       int    `json:"decimals"`
	ProjectURL     string `json:"project_url"`
	LogoURL        string `json:"logo_url"`
	Address        string `json:"address"`
}

type PoolInfo struct {
	Symbol            string `json:"symbol"`
	Name              string `json:"name"`
	Decimals          int    `json:"decimals"`
	Fee               string `json:"fee"`
	TickSpacing       string `json:"tick_spacing"`
	Type              string `json:"type"`
	Address           string `json:"address"`
	CoinAAddress      string `json:"coin_a_address"`
	CoinBAddress      string `json:"coin_b_address"`
	ProjectURL        string `json:"project_url"`
	Sort              string `json:"sort"`
	IsDisplayRewarder bool   `json:"is_display_rewarder"`
	RewarderDisplay1  bool   `json:"rewarder_display_1"`
	RewarderDisplay2  bool   `json:"rewarder_display_2"`
	RewarderDisplay3  bool   `json:"rewarder_display_3"`
	IsStable          bool   `json:"is_stable"`

	// use FetchWarpPoolList get pool with TokenA & TokenB info
	TokenA *TokenInfo
	TokenB *TokenInfo
}

type TokenConfigEvent struct {
	CoinRegistryID types.SuiObjectIdType
	CoinListOwner  types.SuiObjectIdType
	PoolRegistryID types.SuiObjectIdType
	PoolListOwner  types.SuiObjectIdType
}

type TokenModule struct {
	baseModule
}

func NewTokenModule(c *client.Client, config sdkParsedOptions) *TokenModule {
	return &TokenModule{
		baseModule: baseModule{
			c:      c,
			config: config,
		},
	}
}

func (m *TokenModule) FetchTokenList(ctx context.Context, listOwnerAddr string, forceRefresh bool) ([]TokenInfo, error) {
	if !forceRefresh {
		tokens := getTokensCache()
		if tokens != nil {
			return tokens, nil
		}
	}

	var (
		effects   *suitypes.DryRunTransactionBlockResponse
		err       error
		ownerAddr *suitypes.HexData
	)
	if listOwnerAddr == "" {
		effects, err = m.dryRun(ctx,
			*m.config.tokenDisplay,
			"coin_list", "fetch_all_registered_coin_info",
			[]string{}, []any{*m.config.coinRegistryID})
	} else {
		ownerAddr, err = suitypes.NewHexData(listOwnerAddr)
		if err != nil {
			return nil, err
		}
		effects, err = m.dryRun(ctx,
			*m.config.tokenDisplay,
			"coin_list", "fetch_full_list",
			[]string{}, []any{*m.config.coinRegistryID, ownerAddr})
	}
	if err != nil {
		return nil, err
	}

	// parse event
	var tokens []TokenInfo
	err = parseEventWithContent(effects, "::coin_list::FetchCoinListEvent", func(event suitypes.SuiEvent) error {
		var fullList TokenFullList
		data, err := json.Marshal(event.ParsedJson)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &fullList)
		if err != nil {
			return err
		}
		tokens = fullList.FullList.ValueList
		return nil
	})

	if err != nil {
		setTokensCache(tokens)
	}
	return tokens, err
}

func (m *TokenModule) FetchPoolList(ctx context.Context, listOwnerAddr string, forceRefresh bool) ([]PoolInfo, error) {
	if !forceRefresh {
		pools := getPoolsCache()
		if pools != nil {
			return pools, nil
		}
	}

	var (
		effects   *suitypes.DryRunTransactionBlockResponse
		err       error
		ownerAddr *suitypes.HexData
	)
	if listOwnerAddr == "" {
		effects, err = m.dryRun(ctx,
			*m.config.tokenDisplay,
			"lp_list", "fetch_all_registered_coin_info",
			[]string{}, []any{*m.config.poolRegistryID})
	} else {
		ownerAddr, err = suitypes.NewHexData(listOwnerAddr)
		if err != nil {
			return nil, err
		}
		effects, err = m.dryRun(ctx,
			*m.config.tokenDisplay,
			"lp_list", "fetch_full_list",
			[]string{}, []any{*m.config.poolRegistryID, ownerAddr})
	}
	if err != nil {
		return nil, err
	}

	// parse event
	var pools []PoolInfo
	err = parseEventWithContent(effects, "::lp_list::FetchPoolListEvent", func(event suitypes.SuiEvent) error {
		var fullList PoolFullList
		data, err := json.Marshal(event.ParsedJson)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &fullList)
		if err != nil {
			return err
		}
		pools = fullList.FullList.ValueList
		return nil
	})

	if err != nil {
		setPoolsCache(pools)
	}

	return pools, err
}

// FetchWarpPoolList get pool & tokens, wrap token info into pool.TokenA & pool.TokenB
// when tokenA or tokenB not found, a error will return
func (m *TokenModule) FetchWarpPoolList(ctx context.Context, poolOwnerAddr, tokenOwnerAddr string, forceRefresh bool) ([]PoolInfo, error) {
	poolList, err := m.FetchPoolList(ctx, poolOwnerAddr, forceRefresh)
	if err != nil || len(poolList) == 0 {
		return poolList, err
	}

	tokenList, err := m.FetchTokenList(ctx, tokenOwnerAddr, forceRefresh)
	if err != nil {
		return poolList, err
	}

	for i := range poolList {
		for j := range tokenList {
			if poolList[i].CoinAAddress == tokenList[j].Address {
				poolList[i].TokenA = &tokenList[j]
			}
			if poolList[i].CoinBAddress == tokenList[j].Address {
				poolList[i].TokenB = &tokenList[j]
			}
		}
		if poolList[i].TokenA == nil {
			return poolList, errors.New("token a not found")
		}
		if poolList[i].TokenB == nil {
			return poolList, errors.New("token b not found")
		}
	}
	return poolList, nil
}

// getGasObject get a sui object for simulation gas
func getGasObject(c *client.Client, address *suitypes.Address, gas uint64) (*suitypes.HexData, error) {
	coins, err := c.GetSuiCoinsOwnedByAddress(context.Background(), *address)
	if err != nil {
		return nil, err
	}
	coin, err := coins.PickCoinNoLess(gas)
	if err != nil {
		return nil, err
	}
	return &coin.CoinObjectId, nil
}

// parseEventWithContent use function f to parse event in response, witch type has substr
func parseEventWithContent(dryRunResponse *suitypes.DryRunTransactionBlockResponse, substr string, f func(event suitypes.SuiEvent) error) (err error) {
	if !dryRunResponse.Effects.Data.IsSuccess() {
		if nil == dryRunResponse.Effects.Data.V1 {
			return errors.New("parse event failed, no effects")
		}
		return errors.New(dryRunResponse.Effects.Data.V1.Status.Error)
	}

	if len(dryRunResponse.Events) == 0 {
		return errors.New("invalid events")
	}

	defer func() {
		if merr := recover(); merr != nil {
			err = errors.New("event parse failed")
		}
	}()

	for _, event := range dryRunResponse.Events {
		if strings.Contains(event.Type, substr) {
			err = f(event)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
