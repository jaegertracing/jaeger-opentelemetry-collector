package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jaegertracing/jaeger/cmd/flags"
	"github.com/jaegertracing/jaeger/pkg/config"
	"github.com/jaegertracing/jaeger/plugin/storage/es"
	"github.com/open-telemetry/opentelemetry-collector/defaults"
	"github.com/open-telemetry/opentelemetry-collector/service"
	"github.com/spf13/viper"

	"github.com/jaegertracing/jaeger-opentelemetry-collector/pkg/exporter/elasticsearch"
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

	cmd := svc.Command()
	opts := elasticsearch.CreateOptions()
	config.AddFlags(v, cmd, opts.AddFlags, flags.AddConfigFileFlag)

	// parse flags to propagate config file flag value to viper before service start
	cmd.ParseFlags(os.Args)
	err = flags.TryLoadConfigFile(v)
	if err != nil {
		handleErr(fmt.Errorf("could not load Jaeger configuration file %w", err))
	}

	err = svc.Start()
	handleErr(err)
}
