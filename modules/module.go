package modules

import (
	"context"

	"github.com/coming-chat/go-sui/client"
	suitypes "github.com/coming-chat/go-sui/types"
)

type baseModule struct {
	c      *client.Client
	config sdkParsedOptions
}

func (m *baseModule) dryRun(ctx context.Context,
	packageId suitypes.ObjectId,
	module, function string,
	typeArgs []string,
	arguments []any) (*suitypes.DryRunTransactionBlockResponse, error) {
	gasObj, err := getGasObject(m.c, m.config.simulationSigner, m.config.simulationGas)
	if err != nil {
		return nil, err
	}

	tx, err := m.c.MoveCall(ctx, *m.config.simulationSigner, packageId, module, function, []string{}, arguments, gasObj, suitypes.NewSafeSuiBigInt(m.config.simulationGas))
	if err != nil {
		return nil, err
	}

	return m.c.DryRunTransaction(ctx, tx)
}
