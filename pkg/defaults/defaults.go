package defaults

import (
	"github.com/jaegertracing/jaeger/plugin/storage/es"
	"github.com/open-telemetry/opentelemetry-collector/config"
	"github.com/open-telemetry/opentelemetry-collector/defaults"
	"github.com/spf13/viper"

	"github.com/jaegertracing/jaeger-opentelemetry-collector/pkg/exporter/elasticsearch"
)

func Components(v *viper.Viper) (config.Factories, error) {
	esExp := elasticsearch.Factory{Options: func() *es.Options {
		opts := elasticsearch.CreateOptions()
		opts.InitFromViper(v)
		return opts
	}}

	factories, err := defaults.Components()
	if err != nil {
		return factories, err
	}
	factories.Exporters[esExp.Type()] = esExp
	return factories, nil
}
