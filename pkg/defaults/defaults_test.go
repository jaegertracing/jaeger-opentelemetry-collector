package defaults

import (
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestComponents(t *testing.T) {
	factories, err := Components(viper.New())
	require.NoError(t, err)
	assert.Equal(t, "jaeger_elasticsearch", factories.Exporters["jaeger_elasticsearch"].Type())
}
