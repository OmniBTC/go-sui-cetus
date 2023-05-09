package modules

import (
	suitypes "github.com/coming-chat/go-sui/types"
	"github.com/omnibtc/go-sui-cetus/types"
)

type SdkOptions struct {
	SimulationAccount SimulationAccount
	Token             TokenOptions
	Launchpad         LaunchpadOptions
	XWhale            XWhaleOptions
	CLMM              CLMMOptions
}

type SimulationAccount struct {
	Address string
	Gas     uint64
}

type TokenOptions struct {
	TokenDisplay types.SuiObjectIdType
	Config       TokenConfig
}

type TokenConfig struct {
	CoinRegistryID types.SuiObjectIdType
	CoinListOwner  types.SuiObjectIdType
	PoolRegistryID types.SuiObjectIdType
	PoolListOwner  types.SuiObjectIdType
}

type LaunchpadOptions struct {
	IDODisplay  types.SuiObjectIdType
	IDORouter   types.SuiObjectIdType
	LockDisplay types.SuiObjectIdType
	LockRouter  types.SuiObjectIdType
	Config      LaunchpadConfig
}

type LaunchpadConfig struct {
	PoolsID       types.SuiObjectIdType
	AdminCapID    types.SuiObjectIdType
	LockManagerID types.SuiObjectIdType
	ConfigCapID   types.SuiObjectIdType
}

type XWhaleOptions struct {
	XWhaleDisplay    types.SuiObjectIdType
	XWhaleRouter     types.SuiObjectIdType
	DividendsDisplay types.SuiObjectIdType
	DividendsRouter  types.SuiObjectIdType
	BoosterDisplay   types.SuiObjectIdType
	BoosterRouter    types.SuiObjectIdType
	WhaleFaucet      types.SuiObjectIdType
	Config           XWhaleConfig
}

type XWhaleConfig struct {
	XWhaleManagerID   types.SuiObjectIdType
	LockManagerID     types.SuiObjectIdType
	DividendManagerID types.SuiObjectIdType
}

type CLMMOptions struct {
	CLMMDisplay types.SuiObjectIdType
	CLMMRouter  types.SuiObjectIdType
	Config      CLMMConfig
}

type CLMMConfig struct {
	GlobalConfigID types.SuiObjectIdType
	GlobalVaultID  types.SuiObjectIdType
	PoolsID        types.SuiObjectIdType
}

type sdkParsedOptions struct {
	simulationSigner *suitypes.Address
	simulationGas    uint64
	coinRegistryID   *suitypes.HexData
	poolRegistryID   *suitypes.HexData
	tokenDisplay     *suitypes.HexData
	clmmRouter       *suitypes.HexData
}

func (s *SdkOptions) Parse() (options sdkParsedOptions, err error) {
	options.simulationGas = s.SimulationAccount.Gas
	if options.simulationSigner, err = suitypes.NewAddressFromHex(s.SimulationAccount.Address); err != nil {
		return options, err
	}
	if options.coinRegistryID, err = suitypes.NewHexData(string(s.Token.Config.CoinRegistryID)); err != nil {
		return options, err
	}
	if options.poolRegistryID, err = suitypes.NewHexData(string(s.Token.Config.PoolRegistryID)); err != nil {
		return options, err
	}
	if options.tokenDisplay, err = suitypes.NewHexData(string(s.Token.TokenDisplay)); err != nil {
		return options, err
	}
	if options.clmmRouter, err = suitypes.NewHexData(string(s.CLMM.CLMMRouter)); err != nil {
		return options, err
	}

	return
}
