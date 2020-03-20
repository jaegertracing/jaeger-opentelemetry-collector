package cassandra

import (
	"time"

	"github.com/open-telemetry/opentelemetry-collector/config/configmodels"
)

// Config hold configuration of Jaeger Elasticserch exporter/storage.
// TODO consider using cassandra.Options directly from Jaeger
type Config struct {
	configmodels.ExporterSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct.
	Servers                       []string                 `mapstructure:"servers"`
	Port                          int                      `mapstructure:"port""`
	Keyspace                      string                   `mapstructure:"keyspace"`
	LocalDC                       string                   `mapstructure:"local_dc"`
	Consistency                   string                   `mapstructure:"consistency"`
	ProtocolVersion               int                      `mapstructure:"proto_version"`
	ConnectionsPerHost            int                      `mapstructure:"connections_per_host"`
	ConnectTimeout                time.Duration            `mapstructure:"connection_timeout"`
	ReconnectInterval             time.Duration            `mapstructure:"reconnect_interval"`
	SocketKeepAlive               time.Duration            `mapstructure:"socket_keep_alive"`
	SpanStoreWriteCacheTTL        time.Duration            `mapstructure:"span_store_write_cache_ttl"`
	MaxRetryAttempts              int                      `mapstructure:"max_retry_attempts"`
	DisableCompression            bool                     `mapstructure:"disable_compression"`
	Username                      string                   `mapstructure:"username"`
	Password                      string                   `mapstructure:"password"`
	Index                         IndexConfig              `mapstructure:"index"`
}

// IndexConfig configures indexing.
// By default all indexing is enabled.
type IndexConfig struct {
	IndexLogs        bool     `mapstructure:"logs"`
	IndexTags        bool     `mapstructure:"tags"`
	IndexProcessTags bool     `mapstructure:"process_tags"`
	TagBlackList     []string `mapstructure:"tag_back_list"`
	TagWhiteList     []string `mapstructure:"tag_white_list"`
}
