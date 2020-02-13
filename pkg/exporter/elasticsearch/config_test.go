package elasticsearch

import (
	"path"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	factory := &Factory{}
	defaultCfg := factory.CreateDefaultConfig().(*Config)
	assert.Equal(t, defaultServers, defaultCfg.Servers)
	assert.Equal(t, defaultCreateTemplate, defaultCfg.CreateTemplates)
}

func TestLoadConfig(t *testing.T) {
	factories, err := config.ExampleComponents()
	require.NoError(t, err)

	factory := &Factory{}
	factories.Exporters[typeStr] = factory
	cfg, err := config.LoadConfigFile(t, path.Join(".", "testdata", "config.yaml"), factories)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	e1 := cfg.Exporters[typeStr]
	esCfg := e1.(*Config)
	assert.Equal(t, typeStr, esCfg.Name())
	assert.Equal(t, "someUrl", esCfg.Servers)
	assert.Equal(t, true, esCfg.CreateTemplates)
}
