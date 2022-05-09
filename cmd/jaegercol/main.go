package main

import (
	"fmt"
	"log"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/service"

	"github.com/jaegertracing/jaeger-opentelemetry-collector/internal/components"
)

func main() {
	factories, err := components.Components()
	if err != nil {
		log.Fatalf("failed to build components: %v", err)
	}

	info := component.BuildInfo{
		Command:     "otelcontribcol",
		Description: "Jaeger OpenTelemetry Collector",
		Version:     "0.0.0",
	}

	if err = runInteractive(service.CollectorSettings{BuildInfo: info, Factories: factories}); err != nil {
		log.Fatal(err)
	}
}

func runInteractive(params service.CollectorSettings) error {
	cmd := service.NewCommand(params)
	if err := cmd.Execute(); err != nil {
		return fmt.Errorf("collector server run finished with error: %w", err)
	}

	return nil
}
