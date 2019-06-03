package tracing

import (
	"io"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	jeagerConf "github.com/uber/jaeger-client-go/config"
)

// Tracer encapsulates tracing abilities.
type Tracer struct {
	ServiceName  string
	Provider     string
	Logger       logrus.FieldLogger
	JaegerConfig *JaegerConfig

	tracer opentracing.Tracer
	closer io.Closer
}

// JaegerConfig encapsulates jaeger's configuration.
type JaegerConfig struct {
	LocalAgentHostPort string
	SamplerType        string
	SamplerValue       float64
	SamplerServerURL   string
}

// Setup sets up the tracer. Currently supports jaeger.
func (t *Tracer) Setup() error {
	switch strings.ToLower(t.Provider) {
	case "jaeger":
		jc := jeagerConf.Configuration{
			Sampler: &jeagerConf.SamplerConfig{
				SamplingServerURL: t.JaegerConfig.SamplerServerURL,
				Type:              t.JaegerConfig.SamplerType,
				Param:             t.JaegerConfig.SamplerValue,
			},
			Reporter: &jeagerConf.ReporterConfig{
				LocalAgentHostPort: t.JaegerConfig.LocalAgentHostPort,
			},
		}

		closer, err := jc.InitGlobalTracer(
			t.ServiceName,
		)

		if err != nil {
			return err
		}

		t.closer = closer
		t.tracer = opentracing.GlobalTracer()
		t.Logger.Infof("Jaeger tracer configured!")
	case "":
		t.Logger.Infof("No tracer configured - skipping tracing setup")
	default:
		return errors.Errorf("unknown tracer: %s", t.Provider)
	}
	return nil
}

// IsLoaded returns true if the tracer has been loaded.
func (t *Tracer) IsLoaded() bool {
	if t == nil || t.tracer == nil {
		return false
	}
	return true
}

// Close closes the tracer.
func (t *Tracer) Close() {
	if t.closer != nil {
		err := t.closer.Close()
		if err != nil {
			t.Logger.Warn(err)
		}
	}
}
