package cassandra

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/jaeger"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/jaegertracing/jaeger-opentelemetry-collector/internal/metrics"
	"github.com/jaegertracing/jaeger/pkg/cassandra"
	casConfig "github.com/jaegertracing/jaeger/pkg/cassandra/config"
	plugin "github.com/jaegertracing/jaeger/plugin/storage/cassandra"
	cSpanStore "github.com/jaegertracing/jaeger/plugin/storage/cassandra/spanstore"
	"github.com/jaegertracing/jaeger/plugin/storage/cassandra/spanstore/dbmodel"
)

func newTracesExporter(cfg *Config, set component.ExporterCreateSettings) (component.TracesExporter, error) {
	s := newSender(cfg, set.TelemetrySettings)
	return newTracesExporterWithSender(cfg, set, s)
}

func newTracesExporterWithSender(cfg *Config, set component.ExporterCreateSettings, sender *sender) (component.TracesExporter, error) {
	return exporterhelper.NewTracesExporter(
		context.TODO(), set, cfg, sender.pushTraces,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		exporterhelper.WithStart(sender.start),
		exporterhelper.WithShutdown(sender.shutdown),
		exporterhelper.WithTimeout(cfg.TimeoutSettings),
		exporterhelper.WithRetry(cfg.RetrySettings),
		exporterhelper.WithQueue(cfg.QueueSettings),
	)
}

type sender struct {
	name           string
	settings       component.TelemetrySettings
	cfg            *Config
	primaryConfig  casConfig.SessionBuilder
	primarySession cassandra.Session
	// factory  *plugin.Factory // Cannot use unless some functions or fields are exported
	writer *cSpanStore.SpanWriter
}

func newSender(cfg *Config, settings component.TelemetrySettings) *sender {
	return &sender{
		name:          cfg.ID().String(),
		settings:      settings,
		cfg:           cfg,
		primaryConfig: cfg.Options.GetPrimary(),
	}
}

func (s *sender) pushTraces(ctx context.Context, td ptrace.Traces) error {
	batches, err := jaeger.ProtoFromTraces(td)
	if err != nil {
		return consumererror.NewPermanent(fmt.Errorf("failed to push trace data via cassandra exporter: %w", err))
	}

	for _, batch := range batches {
		process := batch.GetProcess()
		if process == nil {
			s.settings.Logger.Debug("NO PROCESS")
		}
		for _, span := range batch.GetSpans() {
			span.Process = process
			err = s.writer.WriteSpan(ctx, span)
			if err != nil {
				s.settings.Logger.Debug("failed to push trace data to cassandra", zap.Error(err))
				return fmt.Errorf("failed to push trace data via cassandra exporter: %w", err)
			}
		}
	}

	return nil
}

func (s *sender) start(_ context.Context, host component.Host) error {
	primarySession, err := s.primaryConfig.NewSession(s.settings.Logger)
	if err != nil {
		return err
	}
	s.primarySession = primarySession

	options, err := writerOptions(s.cfg.Options)
	if err != nil {
		return err
	}

	s.writer = cSpanStore.NewSpanWriter(
		primarySession,
		s.cfg.Options.SpanStoreWriteCacheTTL,
		metrics.New("exporter_cassandra"),
		s.settings.Logger,
		options...,
	)

	/* We can't use the factory due to the way servers is populated in GetPrimary
	// todo(colsen) remove after feedback what to do. The factory would be preferable.
	s.settings.Logger.Info("configuration", zap.Any("options", s.cfg.Options))
	factory := plugin.NewFactory()
	factory.InitFromOptions(s.cfg.Options) // Don't do this, it'll override servers
	//factory.Options = s.cfg.Options        // Do this instead
	if err := factory.Initialize(metrics.NullFactory, s.settings.Logger); err != nil {
		return err
	}
	writer, err := factory.CreateSpanWriter()
	if err != nil {
		return err
	}
	s.factory = factory
	s.writer = writer
	*/

	return nil
}

func (s *sender) shutdown(_ context.Context) error {
	if s.writer == nil {
		return nil
	}
	if err := s.writer.Close(); err != nil {
		return err
	}
	return s.cfg.Options.Primary.TLS.Close()
}

// Copied from plugin/storage/cassandra/factory.go
// todo(colsen) can this be exported upstream?
func writerOptions(opts *plugin.Options) ([]cSpanStore.Option, error) {
	var tagFilters []dbmodel.TagFilter

	// drop all tag filters
	if !opts.Index.Tags || !opts.Index.ProcessTags || !opts.Index.Logs {
		tagFilters = append(tagFilters, dbmodel.NewTagFilterDropAll(!opts.Index.Tags, !opts.Index.ProcessTags, !opts.Index.Logs))
	}

	// black/white list tag filters
	tagIndexBlacklist := opts.TagIndexBlacklist()
	tagIndexWhitelist := opts.TagIndexWhitelist()
	if len(tagIndexBlacklist) > 0 && len(tagIndexWhitelist) > 0 {
		return nil, errors.New("only one of TagIndexBlacklist and TagIndexWhitelist can be specified")
	}
	if len(tagIndexBlacklist) > 0 {
		tagFilters = append(tagFilters, dbmodel.NewBlacklistFilter(tagIndexBlacklist))
	} else if len(tagIndexWhitelist) > 0 {
		tagFilters = append(tagFilters, dbmodel.NewWhitelistFilter(tagIndexWhitelist))
	}

	if len(tagFilters) == 0 {
		return nil, nil
	} else if len(tagFilters) == 1 {
		return []cSpanStore.Option{cSpanStore.TagFilter(tagFilters[0])}, nil
	}

	return []cSpanStore.Option{cSpanStore.TagFilter(dbmodel.NewChainedTagFilter(tagFilters...))}, nil
}
