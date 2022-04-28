module github.com/jaegertracing/jaeger-opentelemetry-collector

go 1.17

require github.com/jaegertracing/jaeger-opentelemetry-collector/exporter/elasticsearch v0.49.0

require (
	github.com/hashicorp/go-hclog v0.12.1 // indirect
	github.com/hashicorp/go-plugin v1.2.0 // indirect
	github.com/jaegertracing/jaeger v1.17.1-0.20200319151430-7304d868c02d
	github.com/magiconair/properties v1.8.1
	github.com/olivere/elastic v6.2.27+incompatible
	github.com/open-telemetry/opentelemetry-collector v0.2.8-0.20200318042533-55be0ec9ddc8
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.7.1
	github.com/uber/jaeger-lib v2.2.0+incompatible
	go.uber.org/zap v1.21.0
)

replace github.com/jaegertracing/jaeger-opentelemetry-collector/exporter/elasticsearch => ./exporter/elasticsearch
