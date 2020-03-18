package app

import (
	"github.com/open-telemetry/opentelemetry-collector/config"
	"github.com/open-telemetry/opentelemetry-collector/config/configmodels"
	"github.com/open-telemetry/opentelemetry-collector/processor/batchprocessor"
	"github.com/open-telemetry/opentelemetry-collector/receiver"
	"github.com/open-telemetry/opentelemetry-collector/receiver/jaegerreceiver"
)

// DefaultConfig creates default configuration.
// It enabled default Jaeger receivers, processors and exporters.
func DefaultConfig(factories config.Factories) *configmodels.Config {
	return &configmodels.Config{
		Receivers:  createReceivers(factories),
		Exporters:  createExporters(factories),
		Processors: createProcessors(factories),
		Service: configmodels.Service{
			Pipelines: map[string]*configmodels.Pipeline{
				"traces": {
					InputType:  configmodels.TracesDataType,
					Receivers:  []string{"jaeger"},
					Exporters:  []string{"jaeger_elasticsearch"},
					Processors: []string{"batch"},
				},
			},
		},
	}
}

func createReceivers(factories config.Factories) configmodels.Receivers {
	rec := factories.Receivers["jaeger"].CreateDefaultConfig().(*jaegerreceiver.Config)
	// TODO load and serve sampling strategies
	rec.Protocols = map[string]*receiver.SecureReceiverSettings{
		"grpc": {
			ReceiverSettings: configmodels.ReceiverSettings{
				Endpoint: "localhost:14250",
			},
		},
		"thrift_http": {
			ReceiverSettings: configmodels.ReceiverSettings{
				Endpoint: "localhost:14268",
			},
		},
		"thrift_compact": {
			ReceiverSettings: configmodels.ReceiverSettings{
				Endpoint: "localhost:6831",
			},
		},
		"thrift_binary": {
			ReceiverSettings: configmodels.ReceiverSettings{
				Endpoint: "localhost:6832",
			},
		},
	}
	return map[string]configmodels.Receiver{
		"jaeger": rec,
	}
}

func createExporters(factories config.Factories) configmodels.Exporters {
	es := factories.Exporters["jaeger_elasticsearch"].CreateDefaultConfig()
	return map[string]configmodels.Exporter{
		"jaeger_elasticsearch": es,
	}
}

func createProcessors(factories config.Factories) configmodels.Processors {
	batch := factories.Processors["batch"].CreateDefaultConfig().(*batchprocessor.Config)
	return map[string]configmodels.Processor{
		"batch": batch,
	}
}
