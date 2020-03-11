package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/open-telemetry/opentelemetry-collector/config"
	"github.com/open-telemetry/opentelemetry-collector/config/configmodels"
	"github.com/open-telemetry/opentelemetry-collector/receiver"
	"github.com/open-telemetry/opentelemetry-collector/receiver/jaegerreceiver"
	"github.com/open-telemetry/opentelemetry-collector/service/builder"

	jflags "github.com/jaegertracing/jaeger/cmd/flags"
	jconfig "github.com/jaegertracing/jaeger/pkg/config"
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
	factories, err := defaults.Components()
	handleErr(err)
	factories.Exporters[esExp.Type()] = esExp

	var cfgFactory service.ConfigFactory
	if getConfigFile() == "" {
		log.Println("Config file not provided, installing default Jaeger components")
		cfgFactory = func(*viper.Viper, config.Factories) (*configmodels.Config, error) {
			return createConfig(factories), nil
		}
	}

	svc, err := service.New(service.Parameters{
		Factories:            factories,
		ApplicationStartInfo: info,
		ConfigFactory:        cfgFactory,
	})
	handleErr(err)

	cmd := svc.Command()
	opts := elasticsearch.CreateOptions()
	jconfig.AddFlags(v, cmd, opts.AddFlags, jflags.AddConfigFileFlag)

	// parse flags to propagate config file flag value to viper before service start
	cmd.ParseFlags(os.Args)
	err = jflags.TryLoadConfigFile(v)
	if err != nil {
		handleErr(fmt.Errorf("could not load Jaeger configuration file %w", err))
	}

	err = svc.Start()
	handleErr(err)
}

func getConfigFile() string {
	f := &flag.FlagSet{}
	builder.Flags(f)
	// parse flags to get file
	f.Parse(os.Args)
	return builder.GetConfigFile()
}

func createConfig(factories config.Factories) *configmodels.Config {
	cfg := &configmodels.Config{}
	jRec := factories.Receivers["jaeger"].CreateDefaultConfig().(*jaegerreceiver.Config)
	// TODO enable other protocols
	jRec.Protocols["grpc"] = &receiver.SecureReceiverSettings{
		ReceiverSettings: configmodels.ReceiverSettings{
			Endpoint: "localhost:14250",
		},
	}
	cfg.Receivers = map[string]configmodels.Receiver{
		"jaeger": jRec,
	}

	esCfg := factories.Exporters["jaeger_elasticsearch"].CreateDefaultConfig()
	cfg.Exporters = map[string]configmodels.Exporter{
		"jaeger_elasticsearch": esCfg,
	}

	cfg.Processors = map[string]configmodels.Processor{
		"batch": factories.Processors["batch"].CreateDefaultConfig(),
	}

	cfg.Service = configmodels.Service{
		Pipelines: map[string]*configmodels.Pipeline{
			"traces": {
				InputType:  configmodels.TracesDataType,
				Receivers:  []string{"jaeger"},
				Exporters:  []string{"jaeger_elasticsearch"},
				Processors: []string{"batch"},
			},
		},
	}
	return cfg
}
