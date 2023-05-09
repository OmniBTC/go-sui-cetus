package modules

import (
	suitypes "github.com/coming-chat/go-sui/types"
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
	TokenDisplay string
	Config       TokenConfig
}

type TokenConfig struct {
	CoinRegistryID string
	CoinListOwner  string
	PoolRegistryID string
	PoolListOwner  string
}

type LaunchpadOptions struct {
	IDODisplay  string
	IDORouter   string
	LockDisplay string
	LockRouter  string
	Config      LaunchpadConfig
}

type LaunchpadConfig struct {
	PoolsID       string
	AdminCapID    string
	LockManagerID string
	ConfigCapID   string
}

type XWhaleOptions struct {
	XWhaleDisplay    string
	XWhaleRouter     string
	DividendsDisplay string
	DividendsRouter  string
	BoosterDisplay   string
	BoosterRouter    string
	WhaleFaucet      string
	Config           XWhaleConfig
}

type XWhaleConfig struct {
	XWhaleManagerID   string
	LockManagerID     string
	DividendManagerID string
}

type CLMMOptions struct {
	CLMMDisplay string
	CLMMRouter  string
	Config      CLMMConfig
}

type CLMMConfig struct {
	GlobalConfigID string
	GlobalVaultID  string
	PoolsID        string
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
