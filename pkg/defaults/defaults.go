package defaults

import (
	cass "github.com/jaegertracing/jaeger/plugin/storage/cassandra"
	"github.com/jaegertracing/jaeger/plugin/storage/es"
	"github.com/open-telemetry/opentelemetry-collector/config"
	"github.com/open-telemetry/opentelemetry-collector/defaults"
	"github.com/spf13/viper"

	"github.com/jaegertracing/jaeger-opentelemetry-collector/pkg/exporter/cassandra"

	"github.com/jaegertracing/jaeger-opentelemetry-collector/pkg/exporter/elasticsearch"
)

// Components creates default and Jaeger factories
func Components(v *viper.Viper) (config.Factories, error) {
	esExp := elasticsearch.Factory{Options: func() *es.Options {
		opts := elasticsearch.CreateOptions()
		opts.InitFromViper(v)
		return opts
	}}

	cassExp := cassandra.Factory{Options: func() *cass.Options {
		opts := cassandra.CreateOptions()
		opts.InitFromViper(v)
		return opts
	}}

	factories, err := defaults.Components()
	if err != nil {
		return factories, err
	}
	factories.Exporters[esExp.Type()] = esExp
	factories.Exporters[cassExp.Type()] = cassExp
	return factories, nil
}
