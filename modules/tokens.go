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

type PoolCreateEvent struct {
	CoinTypeA   string `json:"coin_type_a"`
	CoinTypeB   string `json:"coin_type_b"`
	PoolId      string `json:"pool_id"`
	TickSpacing int    `json:"tick_spacing"`
}

type PoolInfo struct {
	// Symbol            string `json:"symbol"`
	// Name              string `json:"name"`
	// Decimals          int    `json:"decimals"`
	// Fee               string `json:"fee"`
	// TickSpacing       string `json:"tick_spacing"`
	Type         string `json:"type"`
	Address      string `json:"address"`
	CoinAAddress string `json:"coin_a_address"`
	CoinBAddress string `json:"coin_b_address"`
	// ProjectURL        string `json:"project_url"`
	// Sort              string `json:"sort"`
	// IsDisplayRewarder bool   `json:"is_display_rewarder"`
	// RewarderDisplay1  bool   `json:"rewarder_display_1"`
	// RewarderDisplay2  bool   `json:"rewarder_display_2"`
	// RewarderDisplay3  bool   `json:"rewarder_display_3"`
	// IsStable          bool   `json:"is_stable"`

	// use FetchWarpPoolList get pool with TokenA & TokenB info
	TokenA *TokenInfo
	TokenB *TokenInfo
}

type TokenConfigEvent struct {
	CoinRegistryID string
	CoinListOwner  string
	PoolRegistryID string
	PoolListOwner  string
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

func (m *TokenModule) fetchPoolByDryRun(ctx context.Context, listOwnerAddr string, forceRefresh bool) ([]PoolInfo, error) {
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
	var tmpPools []PoolInfo
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
		tmpPools = fullList.FullList.ValueList
		return nil
	})
	if err != nil {
		return nil, err
	}

	objectIds := []suitypes.ObjectId{}
	for _, pool := range tmpPools {
		objId, err := suitypes.NewHexData(pool.Address)
		if err != nil {
			return nil, err
		}
		objectIds = append(objectIds, *objId)
	}
	objectInfos, err := m.c.MultiGetObjects(context.Background(), objectIds, &suitypes.SuiObjectDataOptions{
		ShowType:    true,
		ShowContent: true,
		ShowOwner:   true,
		ShowDisplay: true,
	})
	if err != nil {
		return nil, err
	}
	poolDetails := make([]PoolInfo, 0)
	for _, poolObject := range objectInfos {
		structTag, err := types.ParseMoveStructTag(*poolObject.Data.Type)
		if err != nil {
			continue
		}
		if len(structTag.TypeParams) != 2 {
			continue
		}

		poolDetails = append(poolDetails, PoolInfo{
			Address:      poolObject.Data.ObjectId.ShortString(),
			Type:         *poolObject.Data.Type,
			CoinAAddress: types.GetTypeTagFullName(structTag.TypeParams[0]),
			CoinBAddress: types.GetTypeTagFullName(structTag.TypeParams[1]),
		})
	}
	return poolDetails, nil
}

func (m *TokenModule) fetchPoolByEvents(ctx context.Context) ([]PoolInfo, error) {
	pageSize := uint(2000)
	hasMore := true
	events := make([]suitypes.SuiEvent, 0)
	var cursor *suitypes.EventId
	moveEventType := m.config.createPoolEventPackage + "::factory::CreatePoolEvent"
	for hasMore {
		data, err := m.c.QueryEvents(ctx, suitypes.EventFilter{
			MoveEventType: &moveEventType,
		}, cursor, &pageSize, false)
		if err != nil {
			return nil, err
		}
		cursor = data.NextCursor
		hasMore = data.HasNextPage
		events = append(events, data.Data...)
	}

	poolDetails := make([]PoolInfo, 0)
	for _, event := range events {
		var poolCreateEvent PoolCreateEvent
		data, err := json.Marshal(event.ParsedJson)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, &poolCreateEvent)
		if err != nil {
			return nil, err
		}
		poolDetails = append(poolDetails, PoolInfo{
			Address:      shortCoinTypeWithPrefix(poolCreateEvent.PoolId),
			CoinAAddress: shortCoinTypeWithPrefix(poolCreateEvent.CoinTypeA),
			CoinBAddress: shortCoinTypeWithPrefix(poolCreateEvent.CoinTypeB),
		})
	}
	return poolDetails, nil
}

func (m *TokenModule) FetchPoolList(ctx context.Context, listOwnerAddr string, forceRefresh bool) ([]PoolInfo, error) {
	if !forceRefresh {
		pools := getPoolsCache()
		if pools != nil {
			return pools, nil
		}
	}

	poolDetails, err := m.fetchPoolByEvents(ctx)
	if err != nil {
		setPoolsCache(poolDetails)
	}

	return poolDetails, err
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

	wrapPoolList := make([]PoolInfo, 0, len(poolList))
	for i := range poolList {
		poolItem := poolList[i]
		for j := range tokenList {
			if equalSuiCoinAddress(poolItem.CoinAAddress, tokenList[j].Address) {
				poolItem.TokenA = &tokenList[j]
			}
			if equalSuiCoinAddress(poolItem.CoinBAddress, tokenList[j].Address) {
				poolItem.TokenB = &tokenList[j]
			}
		}
		if poolItem.TokenA == nil || poolItem.TokenB == nil {
			continue
		}
		wrapPoolList = append(wrapPoolList, poolItem)
	}
	return wrapPoolList, nil
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

func equalSuiCoinAddress(x, y string) bool {
	if x == y {
		return true
	}
	x = strings.TrimLeft(x, "x0")
	y = strings.TrimLeft(y, "x0")
	return x == y
}

func shortCoinTypeWithPrefix(address string) string {
	return "0x" + strings.TrimLeft(address, "x0")
}
