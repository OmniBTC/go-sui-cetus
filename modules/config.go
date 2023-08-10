package modules

import (
	"github.com/coming-chat/go-sui/v2/sui_types"
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
	TokenDisplay           string
	CreatePoolEventPackage string
	Config                 TokenConfig
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
	simulationSigner       *sui_types.SuiAddress
	simulationGas          uint64
	coinRegistryID         *sui_types.ObjectID
	poolRegistryID         *sui_types.ObjectID
	coinListOwner          *sui_types.SuiAddress
	poolListOwner          *sui_types.SuiAddress
	tokenDisplay           *sui_types.ObjectID
	clmmRouter             *sui_types.ObjectID
	createPoolEventPackage string
}

func (s *SdkOptions) Parse() (options sdkParsedOptions, err error) {
	options.simulationGas = s.SimulationAccount.Gas
	if options.simulationSigner, err = sui_types.NewAddressFromHex(s.SimulationAccount.Address); err != nil {
		return options, err
	}
	if options.coinRegistryID, err = sui_types.NewObjectIdFromHex(string(s.Token.Config.CoinRegistryID)); err != nil {
		return options, err
	}
	if options.poolRegistryID, err = sui_types.NewObjectIdFromHex(string(s.Token.Config.PoolRegistryID)); err != nil {
		return options, err
	}
	if options.tokenDisplay, err = sui_types.NewObjectIdFromHex(string(s.Token.TokenDisplay)); err != nil {
		return options, err
	}
	if options.coinListOwner, err = sui_types.NewAddressFromHex(string(s.Token.Config.CoinListOwner)); err != nil {
		return options, err
	}
	if options.poolListOwner, err = sui_types.NewAddressFromHex(string(s.Token.Config.PoolListOwner)); err != nil {
		return options, err
	}
	if options.clmmRouter, err = sui_types.NewObjectIdFromHex(string(s.CLMM.CLMMRouter)); err != nil {
		return options, err
	}

	options.createPoolEventPackage = s.Token.CreatePoolEventPackage

	return
}
