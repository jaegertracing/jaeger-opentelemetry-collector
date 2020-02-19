package elasticsearch

import (
	"fmt"
	"time"

	"github.com/open-telemetry/opentelemetry-collector/config/configerror"
	"github.com/open-telemetry/opentelemetry-collector/config/configmodels"
	"github.com/open-telemetry/opentelemetry-collector/exporter"
	"go.uber.org/zap"
)

const (
	typeStr = "jaeger_elasticsearch"

	defaultCreateTemplate    = true
	defaultServers           = "http://localhost:9200"
	defaultShards            = 5
	defaultReplicas          = 1
	defaultBulkActions       = 1000
	defaultBulkSizeBytes     = 5000000
	defaultBulkWorkers       = 1
	defaultBulkFlushInterval = 200 * time.Millisecond
)

// Factory is the factory for Jaeger Elasticsearch exporter.
type Factory struct {
}

// Type gets the type of exporter.
func (Factory) Type() string {
	return typeStr
}

// CreateDefaultConfig returns default configuration of Factory.
func (Factory) CreateDefaultConfig() configmodels.Exporter {
	return &Config{
		Servers:           defaultServers,
		Replicas:          defaultReplicas,
		Shards:            defaultShards,
		CreateTemplates:   defaultCreateTemplate,
		bulkActions:       defaultBulkActions,
		bulkSize:          defaultBulkSizeBytes,
		bulkWorkers:       defaultBulkWorkers,
		bulkFlushInterval: defaultBulkFlushInterval,
	}
}

// CreateTraceExporter creates Jaeger Elasticsearch trace exporter.
func (Factory) CreateTraceExporter(log *zap.Logger, cfg configmodels.Exporter) (exporter.TraceExporter, error) {
	esCfg, ok := cfg.(*Config)
	if !ok {
		return nil, fmt.Errorf("could not cast configuration to %s", typeStr)
	}
	return New(esCfg, log)
}

// CreateMetricsExporter is not implemented.
func (Factory) CreateMetricsExporter(*zap.Logger, configmodels.Exporter) (exporter.MetricsExporter, error) {
	return nil, configerror.ErrDataTypeIsNotSupported
}
