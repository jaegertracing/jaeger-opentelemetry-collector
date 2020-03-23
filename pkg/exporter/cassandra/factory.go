package cassandra

import (
	"fmt"
	"time"

	"github.com/jaegertracing/jaeger/plugin/storage/cassandra"
	"github.com/open-telemetry/opentelemetry-collector/config/configerror"
	"github.com/open-telemetry/opentelemetry-collector/config/configmodels"
	"github.com/open-telemetry/opentelemetry-collector/exporter"
	"go.uber.org/zap"
)

const (
	// TypeStr defines type of the Cassandra exporter.
	TypeStr = "jaeger_cassandra"
)

var (
	defaultWriteCacheTTL = time.Hour * 12
)

// Options returns initialized cassandra.Options structure.
type Options func() *cassandra.Options

// CreateOptions creates Cassandra options supported by this exporter.
func CreateOptions() *cassandra.Options {
	return cassandra.NewOptions("cassandra")
}

// Factory is the factory for Jaeger Cassandra exporter.
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
	cfg := opts.GetPrimary()
	return &Config{
		Servers:                cfg.Servers,
		Port:                   cfg.Port,
		Keyspace:               cfg.Keyspace,
		LocalDC:                cfg.LocalDC,
		Consistency:            cfg.Consistency,
		ProtocolVersion:        cfg.ProtoVersion,
		ConnectionsPerHost:     cfg.ConnectionsPerHost,
		SocketKeepAlive:        cfg.SocketKeepAlive,
		ConnectTimeout:         cfg.ConnectTimeout,
		ReconnectInterval:      cfg.ReconnectInterval,
		MaxRetryAttempts:       cfg.MaxRetryAttempts,
		DisableCompression:     cfg.DisableCompression,
		SpanStoreWriteCacheTTL: defaultWriteCacheTTL,
		Password:               cfg.Authenticator.Basic.Password,
		Username:               cfg.Authenticator.Basic.Username,
		Index: IndexConfig{
			Logs:         true,
			Tags:         true,
			ProcessTags:  true,
			TagBlackList: nil,
			TagWhiteList: nil,
		},
		ExporterSettings: configmodels.ExporterSettings{
			TypeVal: TypeStr,
			NameVal: TypeStr,
		},
	}
}

// CreateTraceExporter creates Jaeger Cassandra trace exporter.
func (Factory) CreateTraceExporter(log *zap.Logger, cfg configmodels.Exporter) (exporter.TraceExporter, error) {
	config, ok := cfg.(*Config)
	if !ok {
		return nil, fmt.Errorf("could not cast configuration to %s", TypeStr)
	}
	return New(config, log)
}

// CreateMetricsExporter is not implemented.
func (Factory) CreateMetricsExporter(*zap.Logger, configmodels.Exporter) (exporter.MetricsExporter, error) {
	return nil, configerror.ErrDataTypeIsNotSupported
}
