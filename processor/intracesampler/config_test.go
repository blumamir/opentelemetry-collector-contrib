package intracesampler

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/confmap/confmaptest"
	"go.opentelemetry.io/collector/otelcol/otelcoltest"
)

func TestLoadConfig(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)
	processors, err := cm.Sub("processors")
	require.NoError(t, err)

	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()

	componentId := component.NewIDWithName(typeStr, "")
	sub, err := processors.Sub(componentId.String())
	require.NoError(t, err)
	require.NoError(t, component.UnmarshalConfig(sub, cfg))

	assert.NoError(t, component.ValidateConfig(cfg))
	expectedConfig := &Config{
		ProcessorSettings:  config.NewProcessorSettings(component.NewID(typeStr)),
		SamplingPercentage: 15.3,
		HashSeed:           22,
		ScopeLeaves:        []string{"foo", "bar"},
	}
	assert.Equal(t, expectedConfig, cfg)
}

func TestLoadInvalidConfig(t *testing.T) {
	factories, err := componenttest.NopFactories()
	require.NoError(t, err)

	factory := NewFactory()
	factories.Processors[typeStr] = factory

	_, err = otelcoltest.LoadConfigAndValidate(filepath.Join("testdata", "invalid.yaml"), factories)
	require.ErrorContains(t, err, "negative sampling rate: -15.30")
}
