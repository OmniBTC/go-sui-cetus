package modules

import (
	"context"
	"testing"

	"github.com/coming-chat/go-sui/client"
)

const mainnet = "https://mainnet.sui.io"
const testnet = "https://fullnode.testnet.sui.io"
const useMainnet = false

func getMainnetConfig() (sdkParsedOptions, error) {
	options := SdkOptions{
		SimulationAccount: SimulationAccount{
			Address: "0x79e54dcebd85b45b6f447358d529a6c08687e3f98c6e9cd790238299fdedeabc",
			Gas:     100_000_000, // 0.1sui
		},
		Token: TokenOptions{
			TokenDisplay: "",
			Config: TokenConfig{
				CoinRegistryID: "0xe0b8cb7e56d465965cac5c5fe26cba558de35d88b9ec712c40f131f72c600151",
				CoinListOwner:  "0x1f6510ee7d8e2b39261bad012f0be0adbecfd75199450b7cbf28efab42dad083",
				PoolRegistryID: "0xab40481f926e686455edf819b4c6485fbbf147a42cf3b95f72ed88c94577e67a",
				PoolListOwner:  "0x6de133b609ea815e1f6a4d50785b798b134f567ec1f4ee113ae73f6900b4012d",
			},
		},
		CLMM: CLMMOptions{
			// CLMMDisplay: "",
			Config: CLMMConfig{
				PoolsID:        "0xf699e7f2276f5c9a75944b37a0c5b5d9ddfd2471bf6242483b03ab2887d198d0",
				GlobalConfigID: "0xdaa46292632c3c4d8f31f23ea0f9b36a28ff3677e9684980e4438403a67a3d8f",
				GlobalVaultID:  "0xce7bceef26d3ad1f6d9b6f13a953f053e6ed3ca77907516481ce99ae8e588f2b",
			},
		},
	}
	return options.Parse()
}

func getTestnetConfig() (sdkParsedOptions, error) {
	options := SdkOptions{
		SimulationAccount: SimulationAccount{
			Address: "0x79e54dcebd85b45b6f447358d529a6c08687e3f98c6e9cd790238299fdedeabc",
			Gas:     100_000_000, // 0.1sui
		},
		Token: TokenOptions{
			TokenDisplay: "0x9dac946be53cf3dd137fa9289a23ecf146f8b87f3eb1a91cc323c93cdd26f8b3",
			Config: TokenConfig{
				CoinRegistryID: "0x2cc70e515a70f3f953363febbab5e01027a25fc01b729d9de29df3db326cc302",
				CoinListOwner:  "0x44682678136f8e8b5a4ffbf9fe06a940b0ccbc9d46d9ae1a68aef2332c2e9cf1",
				PoolRegistryID: "0xfd15ad9a6493fc3ff6b2fc922daeda50bba3df760fda32e53382b7f8dbcbc133",
				PoolListOwner:  "0x494262448a7b8d07e6f00980b67f07a18432a5587d898c27651d18daa4c4c33f",
			},
		},
		CLMM: CLMMOptions{
			CLMMDisplay: "0xb7e7513751376aed2e21b267ef6edebe806a27c979870d3575658dd443ac4248",
			CLMMRouter:  "0x5b77ec28a4077acb46e27e2421aa36b6bbdbe14b4165bc8a7024f10f0fde6112",
			Config: CLMMConfig{
				PoolsID:        "0x67679ae85ea0f39f5c211330bb830f68aeccfbd0085f47f80fc27bef981cc678",
				GlobalConfigID: "0x28565a057d74e4c20d502002bdef5ecf8ff99e1bd9fc4dd11fe549c858ee99d7",
				GlobalVaultID:  "0x6d582d2fa147214d50a0e537084859403e96a3b4554e962d5e993ad5761251f4",
			},
		},
	}
	return options.Parse()
}

func getTestBaseModule() baseModule {
	if useMainnet {
		c, err := client.Dial(mainnet)
		assertNil(err)
		options, err := getMainnetConfig()
		assertNil(err)
		return baseModule{
			c:      c,
			config: options,
		}
	} else {
		c, err := client.Dial(testnet)
		assertNil(err)
		options, err := getTestnetConfig()
		assertNil(err)
		return baseModule{
			c:      c,
			config: options,
		}
	}
}

func TestTokenModule_FetchTokenList(t *testing.T) {
	type fields struct {
		baseModule baseModule
	}
	type args struct {
		ctx           context.Context
		listOwnerAddr string
		forceRefresh  bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []TokenInfo
		wantErr bool
	}{
		{
			name: "fetch all token",
			fields: fields{
				baseModule: getTestBaseModule(),
			},
			args:    args{ctx: context.Background(), listOwnerAddr: "", forceRefresh: false},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := &TokenModule{
				baseModule: tt.fields.baseModule,
			}
			_, err := tm.FetchTokenList(tt.args.ctx, tt.args.listOwnerAddr, tt.args.forceRefresh)
			if (err != nil) != tt.wantErr {
				t.Errorf("TokenModule.FetchTokenList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func assertNil(err error) {
	if err != nil {
		panic(err)
	}
}

func TestTokenModule_FetchWarpPoolList(t *testing.T) {
	type fields struct {
		baseModule baseModule
	}
	type args struct {
		ctx            context.Context
		poolOwnerAddr  string
		tokenOwnerAddr string
		forceRefresh   bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "fetch all token",
			fields: fields{
				baseModule: getTestBaseModule(),
			},
			args:    args{ctx: context.Background(), poolOwnerAddr: "", tokenOwnerAddr: "", forceRefresh: false},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TokenModule{
				baseModule: tt.fields.baseModule,
			}
			_, err := m.FetchWarpPoolList(tt.args.ctx, tt.args.poolOwnerAddr, tt.args.tokenOwnerAddr, tt.args.forceRefresh)
			if (err != nil) != tt.wantErr {
				t.Errorf("TokenModule.FetchWarpPoolList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
