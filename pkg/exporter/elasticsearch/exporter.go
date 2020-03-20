package elasticsearch

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
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
	httpClient := &http.Client{}
	options := []elastic.ClientOptionFunc{
		elastic.SetURL(config.Servers...),
		elastic.SetBasicAuth(config.Username, config.Password),
		elastic.SetSniff(config.Sniffer),
		elastic.SetHttpClient(httpClient),
	}
	if config.TokenFile != "" {
		token, err := loadToken(config.TokenFile)
		if err != nil {
			return nil, err
		}
		httpClient.Transport = &tokenAuthTransport{
			token:   token,
			wrapped: &http.Transport{},
		}
	}

	esRawClient, err := elastic.NewClient(options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client for %s, %v", config.Servers, err)
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
		version, err = getVersion(esRawClient, config.Servers[0])
	}
	var tags []string
	if config.TagsAsFields.AllAsFields && config.TagsAsFields.File != "" {
		tags, err = loadTagsFromFile(config.TagsAsFields.File)
		if err != nil {
			return nil, fmt.Errorf("failed to load tags file: %v", err)
		}
	}

	w := esSpanStore.NewSpanWriter(esSpanStore.SpanWriterParams{
		Logger:              log,
		MetricsFactory:      metrics.NullFactory,
		Client:              eswrapper.WrapESClient(esRawClient, bulk, version),
		IndexPrefix:         config.IndexPrefix,
		UseReadWriteAliases: config.UseWriteAlias,
		AllTagsAsFields:     config.TagsAsFields.AllAsFields,
		TagKeysAsFields:     tags,
		TagDotReplacement:   config.TagsAsFields.DotReplacement,
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
	return dropped, nil
}

func loadTagsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var tags []string
	for scanner.Scan() {
		line := scanner.Text()
		if tag := strings.TrimSpace(line); tag != "" {
			tags = append(tags, tag)
		}
	}
	return tags, nil
}

func loadToken(path string) (string, error) {
	b, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(b), "\r\n"), nil
}

// TokenAuthTransport
type tokenAuthTransport struct {
	token   string
	wrapped *http.Transport
}

func (tr *tokenAuthTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", "Bearer "+tr.token)
	return tr.wrapped.RoundTrip(r)
}
