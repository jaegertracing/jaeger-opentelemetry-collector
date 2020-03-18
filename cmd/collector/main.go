package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jaegertracing/jaeger-opentelemetry-collector/pkg/defaults"

	"github.com/open-telemetry/opentelemetry-collector/config"
	"github.com/open-telemetry/opentelemetry-collector/config/configmodels"
	"github.com/open-telemetry/opentelemetry-collector/service/builder"

	jflags "github.com/jaegertracing/jaeger/cmd/flags"
	jconfig "github.com/jaegertracing/jaeger/pkg/config"
	"github.com/open-telemetry/opentelemetry-collector/service"
	"github.com/spf13/viper"

	"github.com/jaegertracing/jaeger-opentelemetry-collector/cmd/collector/app"
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

	factories, err := defaults.Components(v)
	handleErr(err)

	var cfgFactory service.ConfigFactory
	if getConfigFile() == "" {
		log.Println("Config file not provided, installing default Jaeger components")
		cfgFactory = func(*viper.Viper, config.Factories) (*configmodels.Config, error) {
			return app.DefaultConfig(factories), nil
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
	// parse flags to get the file
	f.Parse(os.Args)
	return builder.GetConfigFile()
}
