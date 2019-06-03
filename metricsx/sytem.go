package metricsx

import (
	"github.com/prometheus/client_golang/prometheus"
)

func RegisterSystemMetrics(r prometheus.Registerer) {
	r.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	r.MustRegister(prometheus.NewGoCollector())
}
