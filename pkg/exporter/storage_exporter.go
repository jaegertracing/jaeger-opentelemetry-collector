package exporter

import (
	"context"

	"github.com/jaegertracing/jaeger/storage/spanstore"
	"github.com/open-telemetry/opentelemetry-collector/consumer/consumerdata"
	"github.com/open-telemetry/opentelemetry-collector/consumer/consumererror"
	jaegertranslator "github.com/open-telemetry/opentelemetry-collector/translator/trace/jaeger"
)

// Storage wraps Jaeger span writer and implements OTEL exporter helper API.
type Storage struct {
	Writer spanstore.Writer
}

// Store stores data into storage
func (s *Storage) Store(ctx context.Context, td consumerdata.TraceData) (droppedSpans int, err error) {
	protoBatch, err := jaegertranslator.OCProtoToJaegerProto(td)
	if err != nil {
		return len(td.Spans), consumererror.Permanent(err)
	}
	dropped := 0
	for _, span := range protoBatch.Spans {
		span.Process = protoBatch.Process
		err := s.Writer.WriteSpan(span)
		// TODO should we wrap errors as we go and return?
		if err != nil {
			dropped++
		}
	}
	return dropped, nil
}
