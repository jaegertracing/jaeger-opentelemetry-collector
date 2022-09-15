package cassandra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"

	casPkg "github.com/jaegertracing/jaeger/pkg/cassandra"
	casConfig "github.com/jaegertracing/jaeger/pkg/cassandra/config"
	"github.com/jaegertracing/jaeger/pkg/cassandra/mocks"
	"github.com/jaegertracing/jaeger/plugin/storage/cassandra"
)

func newOptions(configuration casConfig.Configuration) *cassandra.Options {
	options := cassandra.NewOptions("cassandra")
	options.Primary.Configuration = configuration
	return options
}

// https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/aa298997cb19718ee267fdc7f969573a594698f2/internal/coreinternal/testdata/resource.go#L19-L21
func initResource1(r pcommon.Resource) {
	r.Attributes().UpsertString("resource-attr", "resource-attr-val-1")
}

// https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/af4e1ef5c3109d4c9dc4700a0b2c2bcb17264e5e/internal/coreinternal/testdata/trace.go#L35-L46
func GenerateTracesOneEmptyResourceSpans() ptrace.Traces {
	td := ptrace.NewTraces()
	td.ResourceSpans().AppendEmpty()
	return td
}

func GenerateTracesNoLibraries() ptrace.Traces {
	td := GenerateTracesOneEmptyResourceSpans()
	rs0 := td.ResourceSpans().At(0)
	initResource1(rs0.Resource())
	return td
}

type testSessionBuilder struct{}

func (t *testSessionBuilder) NewSession(logger *zap.Logger) (casPkg.Session, error) {
	session := &mocks.Session{}
	session.On("Close").Return(nil)
	var result []interface{}
	query := &mocks.Query{}
	query.On("Exec").Return(nil)
	session.On("Query", "SELECT * from operation_names_v2 limit 1", result).Return(query)
	return session, nil
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "createExporter",
			config: Config{
				ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
				Options: newOptions(casConfig.Configuration{
					Servers:  []string{"foo.bar"},
					Keyspace: "test",
				}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sender := newSender(&tt.config, componenttest.NewNopExporterCreateSettings().TelemetrySettings)
			got, err := newTracesExporterWithSender(&tt.config, componenttest.NewNopExporterCreateSettings(), sender)
			assert.NoError(t, err)
			assert.NotNil(t, got)
			t.Cleanup(func() {
				require.NoError(t, got.Shutdown(context.Background()))
			})
			sender.primaryConfig = &testSessionBuilder{}

			err = got.Start(context.Background(), componenttest.NewNopHost())
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			err = got.ConsumeTraces(context.Background(), GenerateTracesNoLibraries())
			assert.NoError(t, err)
		})
	}
}
