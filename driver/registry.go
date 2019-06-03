package driver

import (
	"github.com/InVisionApp/go-health"
	"github.com/open-identity/utils/tracing"
	"github.com/ory/herodot"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type RegistryLogger interface {
	Logger() logrus.FieldLogger
}

type RegistryWriter interface {
	Writer() herodot.Writer
}

type RegistryMetrics interface {
	Metrics() *prometheus.Registry
}

type RegistryHealth interface {
	Health() health.IHealth
}

type RegistryTracer interface {
	Tracer() *tracing.Tracer
}
