package elasticsearch

import (
	"fmt"

	"github.com/jaegertracing/jaeger/plugin/storage/es"
	"github.com/open-telemetry/opentelemetry-collector/config/configerror"
	"github.com/open-telemetry/opentelemetry-collector/config/configmodels"
	"github.com/open-telemetry/opentelemetry-collector/exporter"
	"go.uber.org/zap"
)

const (
	// TypeStr defines type of the Elasticsearch exporter.
	TypeStr = "jaeger_elasticsearch"
)

// Options returns initialized es.Options structure.
type Options func() *es.Options

// CreateOptions creates Elasticsearch options supported by this exporter.
func CreateOptions() *es.Options {
	return es.NewOptions("es")
}

// Factory is the factory for Jaeger Elasticsearch exporter.
type Factory struct {
	Options Options
}

// Type gets the type of exporter.
func (Factory) Type() string {
	return TypeStr
}

// CreateDefaultConfig returns default configuration of Factory.
func (f Factory) CreateDefaultConfig() configmodels.Exporter {
	opts := f.Options()
	return &Config{
		Servers:           opts.GetPrimary().Servers,
		Shards:            uint(opts.GetPrimary().GetNumShards()),
		Replicas:          uint(opts.GetPrimary().GetNumReplicas()),
		IndexPrefix:       opts.GetPrimary().GetIndexPrefix(),
		CreateTemplates:   opts.GetPrimary().IsCreateIndexTemplates(),
		bulkActions:       opts.GetPrimary().BulkActions,
		bulkSize:          opts.GetPrimary().BulkSize,
		bulkWorkers:       opts.GetPrimary().BulkWorkers,
		bulkFlushInterval: opts.GetPrimary().BulkFlushInterval,
		Version:           opts.GetPrimary().Version,
		Sniffer:           opts.GetPrimary().Sniffer,
		UseWriteAlias:     opts.GetPrimary().UseReadWriteAliases,
		Username:          opts.GetPrimary().Username,
		Password:          opts.GetPrimary().Password,
		TokenFile:         opts.GetPrimary().TokenFilePath,
		TagsAsFields: TagsAsFields{
			AllAsFields:    opts.GetPrimary().AllTagsAsFields,
			DotReplacement: opts.GetPrimary().TagDotReplacement,
			File:           opts.GetPrimary().TagsFilePath,
		},

		ExporterSettings: configmodels.ExporterSettings{
			TypeVal: TypeStr,
			NameVal: TypeStr,
		},
	}
}

// CreateTraceExporter creates Jaeger Elasticsearch trace exporter.
func (Factory) CreateTraceExporter(log *zap.Logger, cfg configmodels.Exporter) (exporter.TraceExporter, error) {
	esCfg, ok := cfg.(*Config)
	if !ok {
		return nil, fmt.Errorf("could not cast configuration to %s", TypeStr)
	}
	return New(esCfg, log)
}

// CreateMetricsExporter is not implemented.
func (Factory) CreateMetricsExporter(*zap.Logger, configmodels.Exporter) (exporter.MetricsExporter, error) {
	return nil, configerror.ErrDataTypeIsNotSupported
}
