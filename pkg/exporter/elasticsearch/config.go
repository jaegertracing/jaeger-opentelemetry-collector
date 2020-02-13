package elasticsearch

import (
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
}
