package main

import (
	"log"

	"github.com/jaegertracing/jaeger/pkg/config"
	"github.com/jaegertracing/jaeger/plugin/storage/es"
	"github.com/spf13/viper"

	"github.com/jaegertracing/jaeger-opentelemetry-collector/pkg/exporter/elasticsearch"

	"github.com/open-telemetry/opentelemetry-collector/defaults"
	"github.com/open-telemetry/opentelemetry-collector/service"
)

func main() {
	handleErr := func(err error) {
		if err != nil {
			log.Fatalf("Failed to run the service: %v", err)
		}
	}

	info := service.ApplicationStartInfo{
		ExeName:  "jaeger-opentelemetry-collector",
		LongName: "Jaeger OpenTelemetry Collector",
		// TODO
		//Version:  version.Version,
		//GitHash:  version.GitHash,
	}

	// TODO should we support jaeger conf file?
	v := viper.New()

	esExp := elasticsearch.Factory{Options: func() *es.Options {
		opts := elasticsearch.CreateOptions()
		opts.InitFromViper(v)
		return opts
	}}
	cmpts, err := defaults.Components()
	handleErr(err)
	cmpts.Exporters[esExp.Type()] = esExp

	svc, err := service.New(cmpts, info)
	handleErr(err)

	opts := elasticsearch.CreateOptions()
	config.AddFlags(v, svc.Command(), opts.AddFlags)

	err = svc.Start()
	handleErr(err)
}
