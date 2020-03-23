package cassandra

import (
	cassandraConfig "github.com/jaegertracing/jaeger/pkg/cassandra/config"
	"github.com/jaegertracing/jaeger/plugin/storage/cassandra/spanstore"
	"github.com/open-telemetry/opentelemetry-collector/exporter"
	"github.com/open-telemetry/opentelemetry-collector/exporter/exporterhelper"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"

	jexporter "github.com/jaegertracing/jaeger-opentelemetry-collector/pkg/exporter"
)

// New creates Cassandra exporter/storage.
func New(config *Config, log *zap.Logger) (exporter.TraceExporter, error) {
	cfg := cassandraConfig.Configuration{
		Servers:            config.Servers,
		Port:               config.Port,
		Keyspace:           config.Keyspace,
		LocalDC:            config.LocalDC,
		Consistency:        config.Consistency,
		ProtoVersion:       config.ProtocolVersion,
		MaxRetryAttempts:   config.MaxRetryAttempts,
		DisableCompression: config.DisableCompression,
		SocketKeepAlive:    config.SocketKeepAlive,
		ConnectTimeout:     config.ConnectTimeout,
		ReconnectInterval:  config.ReconnectInterval,
		ConnectionsPerHost: config.ConnectionsPerHost,
		Authenticator: cassandraConfig.Authenticator{Basic: cassandraConfig.BasicAuthenticator{
			Username: config.Username,
			Password: config.Password,
		}},
	}

	session, err := cfg.NewSession()
	if err != nil {
		return nil, err
	}

	// TODO configure indexing
	// make cassandra.Factory writeOptions public
	spanWriter := spanstore.NewSpanWriter(session, config.SpanStoreWriteCacheTTL, metrics.NullFactory, zap.NewNop())
	storage := jexporter.Storage{Writer: spanWriter}
	return exporterhelper.NewTraceExporter(
		config,
		storage.Store,
		exporterhelper.WithShutdown(func() error {
			session.Close()
			return nil
		}))
}
