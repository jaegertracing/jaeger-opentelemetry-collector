# EXPERIMENTAL - DO NOT USE

## Jaeger OpenTelemetry collector

This repository hosts Jaeger specific components for OpenTelemetry collector. The repository is inspired by [opentelemetry-collector-contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib) in a way that every component is a separate Golang module. The final distribution is assembled in [jaeger-opentelemetry-releases](https://github.com/jaegertracing/jaeger-opentelemetry-releases).

## Release

Jaeger OpenTelemetry collector versions are aligned with OpenTelemetry collector upstream versions.

To release `0.49.0` push the following tag:

```bash
git tag v0.49.0 && git push origin v.0.49.0
```
