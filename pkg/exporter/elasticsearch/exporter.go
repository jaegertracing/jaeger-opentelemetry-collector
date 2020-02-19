package elasticsearch

import (
	"context"
	"strconv"
	"strings"

	eswrapper "github.com/jaegertracing/jaeger/pkg/es/wrapper"
	"github.com/jaegertracing/jaeger/plugin/storage/es"
	esSpanStore "github.com/jaegertracing/jaeger/plugin/storage/es/spanstore"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	"github.com/olivere/elastic"
	"github.com/open-telemetry/opentelemetry-collector/consumer/consumerdata"
	"github.com/open-telemetry/opentelemetry-collector/consumer/consumererror"
	"github.com/open-telemetry/opentelemetry-collector/exporter"
	"github.com/open-telemetry/opentelemetry-collector/exporter/exporterhelper"
	jaegertranslator "github.com/open-telemetry/opentelemetry-collector/translator/trace/jaeger"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

// New creates new Elasticsearch exporter/storage.
func New(config *Config, log *zap.Logger) (exporter.TraceExporter, error) {
	servers := strings.Split(config.Servers, ",")
	esRawClient, err := elastic.NewClient(
		elastic.SetURL(servers...),
		elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}
	bulk, err := esRawClient.BulkProcessor().
		BulkActions(config.bulkActions).
		BulkSize(config.bulkSize).
		Workers(config.bulkWorkers).
		FlushInterval(config.bulkFlushInterval).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	version := config.Version
	if version == 0 {
		version, err = getVersion(esRawClient, servers[0])
	}

	w := esSpanStore.NewSpanWriter(esSpanStore.SpanWriterParams{
		Logger:         log,
		MetricsFactory: metrics.NullFactory,
		Client:         eswrapper.WrapESClient(esRawClient, bulk, version),
		IndexPrefix:    config.IndexPrefix,
	})

	if config.CreateTemplates {
		spanMapping, serviceMapping := es.GetMappings(int64(config.Shards), int64(config.Shards), version)
		err := w.CreateTemplates(spanMapping, serviceMapping)
		if err != nil {
			return nil, err
		}
	}

	esStorage := &storage{
		w: w,
	}
	return exporterhelper.NewTraceExporter(
		config,
		esStorage.store,
		exporterhelper.WithTracing(true),
		exporterhelper.WithMetrics(true),
		exporterhelper.WithShutdown(func() error {
			return w.Close()
		}))
}

func getVersion(client *elastic.Client, server string) (uint, error) {
	pingResult, _, err := client.Ping(server).Do(context.Background())
	if err != nil {
		return 0, err
	}
	esVersion, err := strconv.Atoi(string(pingResult.Version.Number[0]))
	if err != nil {
		return 0, err
	}
	return uint(esVersion), nil
}

type storage struct {
	w spanstore.Writer
}

func (s *storage) store(ctx context.Context, td consumerdata.TraceData) (droppedSpans int, err error) {
	protoBatch, err := jaegertranslator.OCProtoToJaegerProto(td)
	if err != nil {
		return len(td.Spans), consumererror.Permanent(err)
	}
	dropped := 0
	for _, span := range protoBatch.Spans {
		span.Process = protoBatch.Process
		err := s.w.WriteSpan(span)
		// TODO should we wrap errors as we go and return?
		if err != nil {
			dropped++
		}
	}
	return 0, nil
}
