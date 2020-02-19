package elasticsearch

import (
	"time"

	"github.com/open-telemetry/opentelemetry-collector/config/configmodels"
)

// Config hold configuration of Jaeger Elasticserch exporter/storage.
type Config struct {
	configmodels.ExporterSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct.
	// Servers define Elasticsearch server URLs e.g. http://localhost:9200
	Servers string `mapstructure:"server_urls"`
	// Shards define number of primary shards
	Shards uint `mapstructure:"shards"`
	// Replicas define number of replica shards
	Replicas uint `mapstructure:"replicas"`
	// Version defines Elasticsearch version. 0 means that ES version will be obtained from ping endpoint at startup.
	Version uint `mapstructure:"version"`
	// CreateTemplates defines whether Jaeger templates will be installed to Elasticsearch at startup.
	CreateTemplates bool `mapstructure:"create_mappings"`
	// IndexPrefix defines options prefix of Jaeger indices. For example "production" creates "production-jaeger-*"
	IndexPrefix string `mapstructure:"index_prefix"`

	// BulkActions defines the number of requests that can be enqueued before the bulk processor decides to commit.
	bulkActions int
	// BulkFlushInterval defines duration after which bulk requests are committed, regardless of other thresholds.
	bulkFlushInterval time.Duration
	// BulkSize defines the number of bytes that the bulk requests can take up before the bulk processor decides to commit.
	bulkSize int
	// BulkWorkers define the number of workers that are able to receive bulk requests and eventually commit them to Elasticsearch.
	bulkWorkers int
}
