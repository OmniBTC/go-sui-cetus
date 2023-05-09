package suicetus

import (
	"github.com/coming-chat/go-sui/client"
	"github.com/omnibtc/go-sui-cetus/modules"
)

type SDK struct {
	TokenModule *modules.TokenModule
}

func NewCetusSDK(c *client.Client, config modules.SdkOptions) (*SDK, error) {
	options, err := config.Parse()
	if err != nil {
		return nil, err
	}
	return &SDK{
		TokenModule: modules.NewTokenModule(c, options),
	}, nil
}
