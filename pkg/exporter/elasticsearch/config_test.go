package elasticsearch

import (
	"path"
	"testing"

	"github.com/jaegertracing/jaeger/cmd/flags"

	jConfig "github.com/jaegertracing/jaeger/pkg/config"
	"github.com/jaegertracing/jaeger/plugin/storage/es"

	"github.com/open-telemetry/opentelemetry-collector/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	factory := &Factory{Options: func() *es.Options {
		return CreateOptions()
	}}
	defaultCfg := factory.CreateDefaultConfig().(*Config)
	assert.Equal(t, []string{"http://127.0.0.1:9200"}, defaultCfg.Servers)
	assert.Equal(t, true, defaultCfg.CreateTemplates)
}

func TestLoadConfigAndFlags(t *testing.T) {
	factories, err := config.ExampleComponents()
	require.NoError(t, err)

	v, c := jConfig.Viperize(CreateOptions().AddFlags, flags.AddConfigFileFlag)
	err = c.ParseFlags([]string{"--es.server-urls=bar", "--es.index-prefix=staging", "--config-file=./testdata/jaeger-config.yaml"})
	require.NoError(t, err)

	err = flags.TryLoadConfigFile(v)
	require.NoError(t, err)

	factory := &Factory{Options: func() *es.Options {
		opts := CreateOptions()
		opts.InitFromViper(v)
		require.Equal(t, []string{"bar"}, opts.GetPrimary().Servers)
		require.Equal(t, "staging", opts.GetPrimary().GetIndexPrefix())
		assert.Equal(t, int64(100), opts.GetPrimary().NumShards)
		return opts
	}}

	factories.Exporters[typeStr] = factory
	cfg, err := config.LoadConfigFile(t, path.Join(".", "testdata", "config.yaml"), factories)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	e1 := cfg.Exporters[typeStr]
	esCfg := e1.(*Config)
	assert.Equal(t, typeStr, esCfg.Name())
	assert.Equal(t, []string{"someUrl"}, esCfg.Servers)
	assert.Equal(t, true, esCfg.CreateTemplates)
	assert.Equal(t, "staging", esCfg.IndexPrefix)
	assert.Equal(t, uint(100), esCfg.Shards)
}
