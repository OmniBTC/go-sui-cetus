package modules

import (
	"context"

	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/sui_types"
	suitypes "github.com/coming-chat/go-sui/v2/types"
)

type baseModule struct {
	c      *client.Client
	config sdkParsedOptions
}

func (m *baseModule) dryRun(ctx context.Context,
	packageId sui_types.ObjectID,
	module, function string,
	typeArgs []string,
	arguments []any) (*suitypes.DryRunTransactionBlockResponse, error) {
	gasObj, err := getGasObject(m.c, m.config.simulationSigner, m.config.simulationGas)
	if err != nil {
		return nil, err
	}

	tx, err := m.c.MoveCall(ctx, *m.config.simulationSigner, packageId, module, function, typeArgs, arguments, gasObj, suitypes.NewSafeSuiBigInt(m.config.simulationGas))
	if err != nil {
		return nil, err
	}

	return m.c.DryRunTransaction(ctx, tx.TxBytes)
}
