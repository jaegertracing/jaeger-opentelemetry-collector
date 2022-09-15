package metrics

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/jaegertracing/jaeger/pkg/metrics"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

func New(namespace string) metrics.Factory {
	factory := newFactory(&Factory{
		normalizer: strings.NewReplacer(".", "_", "-", "_"),
	},
		"",  // scope
		nil, // tags
	)
	return factory.Namespace(metrics.NSOptions{
		Name: namespace,
	})
}

func newFactory(parent *Factory, scope string, tags map[string]string) *Factory {
	return &Factory{
		normalizer: parent.normalizer,
		scope:      scope,
		tags:       tags,
	}
}

type Factory struct {
	scope      string
	tags       map[string]string
	normalizer *strings.Replacer
}

func (f *Factory) Counter(options metrics.Options) metrics.Counter {
	help := strings.TrimSpace(options.Help)
	if len(help) == 0 {
		help = options.Name
	}
	name := counterNamingConvention(f.subScope(options.Name))
	tags := f.mergeTags(options.Tags)
	tagKeys := tagKeys(tags)
	measure := stats.Int64(
		name,
		help,
		stats.UnitDimensionless,
	)
	v := &view.View{
		Name:        name,
		Measure:     measure,
		Description: help,
		Aggregation: view.Count(),
		TagKeys:     tagKeys,
	}
	err := view.Register(v)
	if err != nil {
		panic(err)
	}

	tagMutators := []tag.Mutator{}
	for k, v := range tags {
		key := tag.MustNewKey(k)
		tagMutators = append(tagMutators, tag.Insert(key, v))
	}

	ctx, err := tag.New(context.Background(), tagMutators...)
	if err != nil {
		panic(err)
	}
	return &counter{
		ctx,
		measure,
	}
}

type counter struct {
	ctx     context.Context
	counter *stats.Int64Measure
}

func (c *counter) Inc(v int64) {
	stats.Record(
		c.ctx,
		c.counter.M(v),
	)
}

func (f *Factory) Timer(options metrics.TimerOptions) metrics.Timer {
	help := strings.TrimSpace(options.Help)
	if len(help) == 0 {
		help = options.Name
	}
	name := f.subScope(options.Name)
	buckets := view.Distribution(0, 5, 10, 20, 50, 100, 200, 500, 1000, 2000, 5000) // Hardcoded :(
	tags := f.mergeTags(options.Tags)
	tagKeys := tagKeys(tags)
	measure := stats.Int64(
		name,
		help,
		stats.UnitMilliseconds,
	)
	v := &view.View{
		Name:        name,
		Measure:     measure,
		Description: help,
		Aggregation: buckets, // bad-name for opencensus
		TagKeys:     tagKeys,
	}
	// Error is currently ignored. Multiple registrations are attempted.
	// The counters are only registered once.
	// It does not seem to affect the end metrics when just ignored.
	// Maintainers might have an opinion if this is acceptable or not.
	_ = view.Register(v)

	tagMutators := []tag.Mutator{}
	for k, v := range tags {
		key := tag.MustNewKey(k)
		tagMutators = append(tagMutators, tag.Insert(key, v))
	}

	ctx, err := tag.New(context.Background(), tagMutators...)
	if err != nil {
		panic(err)
	}
	return &timer{
		ctx,
		measure,
	}
}

type timer struct {
	ctx          context.Context
	distribution *stats.Int64Measure
}

func (t *timer) Record(timeSpent time.Duration) {
	stats.Record(
		t.ctx,
		t.distribution.M(timeSpent.Milliseconds()),
	)
}

func (f *Factory) Gauge(options metrics.Options) metrics.Gauge {
	return metrics.NullGauge
}

func (f *Factory) Histogram(options metrics.HistogramOptions) metrics.Histogram {
	return metrics.NullHistogram
}

func (f *Factory) Namespace(scope metrics.NSOptions) metrics.Factory {
	return newFactory(f, f.subScope(scope.Name), f.mergeTags(scope.Tags))
}

func (f *Factory) subScope(name string) string {
	if f.scope == "" {
		return f.normalize(name)
	}
	if name == "" {
		return f.normalize(f.scope)
	}
	return f.normalize(f.scope + "_" + name) // Hard-coded
}

func (f *Factory) normalize(v string) string {
	return f.normalizer.Replace(v)
}

func (f *Factory) mergeTags(tags map[string]string) map[string]string {
	ret := make(map[string]string, len(f.tags)+len(tags))
	for k, v := range f.tags {
		ret[k] = v
	}
	for k, v := range tags {
		ret[k] = v
	}
	return ret
}

func counterNamingConvention(name string) string {
	if !strings.HasSuffix(name, "_total") {
		name += "_total"
	}
	return name
}

func tagNames(tags map[string]string) []string {
	ret := make([]string, 0, len(tags))
	for k := range tags {
		ret = append(ret, k)
	}
	sort.Strings(ret)
	return ret
}

func tagKeys(tags map[string]string) []tag.Key {
	tagNames := tagNames(tags)
	keys := make([]tag.Key, 0, len(tagNames))
	for _, tagName := range tagNames {
		keys = append(keys, tag.MustNewKey(tagName))
	}
	return keys
}
