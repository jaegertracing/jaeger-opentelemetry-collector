package elasticsearch

import (
	"go.opentelemetry.io/collector/component"
)

type Config struct {
	// component.ExporterConfigSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct
}

var _ component.Config = (*Config)(nil)

// Validate checks if the exporter configuration is valid
func (cfg *Config) Validate() error {
	return nil
}
