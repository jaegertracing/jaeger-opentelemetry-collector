package cassandra

import (
	"path"
	"testing"
	"time"

	"github.com/jaegertracing/jaeger/cmd/flags"
	jConfig "github.com/jaegertracing/jaeger/pkg/config"
	"github.com/jaegertracing/jaeger/plugin/storage/cassandra"
	"github.com/open-telemetry/opentelemetry-collector/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	factory := &Factory{Options: func() *cassandra.Options {
		return CreateOptions()
	}}
	defaultCfg := factory.CreateDefaultConfig().(*Config)
	assert.Equal(t, []string{"127.0.0.1"}, defaultCfg.Servers)
	assert.Equal(t, 2, defaultCfg.ConnectionsPerHost)
	assert.Equal(t, "jaeger_v1_test", defaultCfg.Keyspace)
	assert.Equal(t, 3, defaultCfg.MaxRetryAttempts)
	assert.Equal(t, 4, defaultCfg.ProtocolVersion)
	assert.Equal(t, time.Minute, defaultCfg.ReconnectInterval)
	assert.Equal(t, defaultWriteCacheTTL, defaultCfg.SpanStoreWriteCacheTTL)
	assert.Equal(t, true, defaultCfg.Index.IndexTags)
	assert.Equal(t, true, defaultCfg.Index.IndexLogs)
	assert.Equal(t, true, defaultCfg.Index.IndexProcessTags)
}

func TestLoadConfigAndFlags(t *testing.T) {
	factories, err := config.ExampleComponents()
	require.NoError(t, err)

	v, c := jConfig.Viperize(CreateOptions().AddFlags, flags.AddConfigFileFlag)
	err = c.ParseFlags([]string{"--cassandra.servers=bar", "--config-file=./testdata/jaeger-config.yaml"})
	require.NoError(t, err)

	err = flags.TryLoadConfigFile(v)
	require.NoError(t, err)

	factory := &Factory{Options: func() *cassandra.Options {
		opts := CreateOptions()
		opts.InitFromViper(v)
		require.Equal(t, []string{"bar"}, opts.GetPrimary().Servers)
		return opts
	}}

	factories.Exporters[TypeStr] = factory
	colConfig, err := config.LoadConfigFile(t, path.Join(".", "testdata", "config.yaml"), factories)
	require.NoError(t, err)
	require.NotNil(t, colConfig)

	cfg := colConfig.Exporters[TypeStr].(*Config)
	assert.Equal(t, TypeStr, cfg.Name())
	assert.Equal(t, []string{"first", "second"}, cfg.Servers)
	assert.Equal(t, false, cfg.Index.IndexTags)
}
