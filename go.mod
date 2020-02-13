module github.com/jaegertracing/jaeger-opentelemetry-collector

go 1.13

require (
	github.com/jaegertracing/jaeger v1.14.0
	github.com/olivere/elastic v6.2.27+incompatible
	github.com/open-telemetry/opentelemetry-collector v0.2.6
	github.com/securego/gosec v0.0.0-20200203094520-d13bb6d2420c // indirect
	github.com/stretchr/testify v1.4.0
	github.com/uber/jaeger-lib v2.2.0+incompatible
	go.uber.org/zap v1.10.0
)

replace k8s.io/client-go => k8s.io/client-go v0.0.0-20190620085101-78d2af792bab
