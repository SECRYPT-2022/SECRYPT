package server

import (
	"github.com/SECRYPT-2022/SECRYPT/consensus"
	"github.com/SECRYPT-2022/SECRYPT/network"
	"github.com/SECRYPT-2022/SECRYPT/txpool"
)

// serverMetrics holds the metric instances of all sub systems
type serverMetrics struct {
	consensus *consensus.Metrics
	network   *network.Metrics
	txpool    *txpool.Metrics
}

// metricProvider serverMetric instance for the given ChainID and nameSpace
func metricProvider(nameSpace string, chainID string, metricsRequired bool) *serverMetrics {
	if metricsRequired {
		return &serverMetrics{
			consensus: consensus.GetPrometheusMetrics(nameSpace, "chain_id", chainID),
			network:   network.GetPrometheusMetrics(nameSpace, "chain_id", chainID),
			txpool:    txpool.GetPrometheusMetrics(nameSpace, "chain_id", chainID),
		}
	}

	return &serverMetrics{
		consensus: consensus.NilMetrics(),
		network:   network.NilMetrics(),
		txpool:    txpool.NilMetrics(),
	}
}
