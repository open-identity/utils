package metricsx

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/xhandler"
	"github.com/urfave/negroni"
)

const (
	LabelCode   = "code"
	LabelMethod = "method"
	LabelUrl    = "url"
)

func InstrumentHandler(r *prometheus.Registry) func(next xhandler.HandlerC) xhandler.HandlerC {
	return InstrumentHandlerWithIncludePath(r, "/")
}

func InstrumentHandlerWithIncludePath(r *prometheus.Registry, includePathPrefix string) func(next xhandler.HandlerC) xhandler.HandlerC {

	cnt := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of requests by HTTP status code.",
		},
		[]string{LabelCode, LabelMethod, LabelUrl},
	)

	gge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "http_requests_in_flight",
		Help: "Current number of requests being served.",
	}, []string{LabelMethod, LabelUrl})

	hist := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_requests_time",
		Help:    "Current number of scrapes being served.",
		Buckets: prometheus.ExponentialBuckets(1, 2, 5),
	}, []string{LabelCode, LabelMethod, LabelUrl})
	r.MustRegister(cnt, gge, hist)

	return func(next xhandler.HandlerC) xhandler.HandlerC {
		return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			if includePathPrefix == "" || !strings.HasPrefix(r.URL.Path, includePathPrefix) {
				next.ServeHTTPC(ctx, w, r)
				return
			}

			method := sanitizeMethod(r.Method)

			labels := prometheus.Labels{}
			labels[LabelMethod] = method
			labels[LabelUrl] = r.URL.Path
			gauge := gge.With(labels)
			gauge.Inc()
			defer gauge.Dec()

			start := time.Now()
			next.ServeHTTPC(ctx, w, r)
			timeElapsed := time.Since(start).Seconds()

			nw := w.(negroni.ResponseWriter)

			labels[LabelCode] = sanitizeCode(nw.Status())

			cnt.With(labels).Inc()
			hist.With(labels).Observe(timeElapsed)
		})
	}
}

func InstrumentHandlerCounter(counter *prometheus.CounterVec) func(next xhandler.HandlerC) xhandler.HandlerC {

	return func(next xhandler.HandlerC) xhandler.HandlerC {
		return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {

			next.ServeHTTPC(ctx, w, r)

			nw := w.(negroni.ResponseWriter)

			labels := prometheus.Labels{}
			labels[LabelCode] = sanitizeCode(nw.Status())
			labels[LabelMethod] = sanitizeMethod(r.Method)
			labels[LabelUrl] = r.URL.Path

			counter.With(labels).Inc()
		})
	}
}

func InstrumentHandlerInFlight(gauge *prometheus.GaugeVec) func(next xhandler.HandlerC) xhandler.HandlerC {

	return func(next xhandler.HandlerC) xhandler.HandlerC {
		return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			labels := prometheus.Labels{}
			labels[LabelMethod] = sanitizeMethod(r.Method)
			labels[LabelUrl] = r.URL.Path
			gauge := gauge.With(labels)
			gauge.Inc()
			defer gauge.Dec()
			next.ServeHTTPC(ctx, w, r)
		})
	}
}

func InstrumentTimeConsumed(hist *prometheus.HistogramVec) func(next xhandler.HandlerC) xhandler.HandlerC {

	return func(next xhandler.HandlerC) xhandler.HandlerC {
		return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {

			start := time.Now()
			next.ServeHTTPC(ctx, w, r)
			timeElapsed := time.Since(start).Seconds()

			nw := w.(negroni.ResponseWriter)

			labels := prometheus.Labels{}
			labels[LabelMethod] = sanitizeMethod(r.Method)
			labels[LabelUrl] = r.URL.Path
			labels[LabelCode] = sanitizeCode(nw.Status())

			hist.With(labels).Observe(timeElapsed)
		})
	}
}

func sanitizeMethod(m string) string {
	switch m {
	case "GET", "get":
		return "get"
	case "PUT", "put":
		return "put"
	case "HEAD", "head":
		return "head"
	case "POST", "post":
		return "post"
	case "DELETE", "delete":
		return "delete"
	case "CONNECT", "connect":
		return "connect"
	case "OPTIONS", "options":
		return "options"
	case "NOTIFY", "notify":
		return "notify"
	default:
		return strings.ToLower(m)
	}
}

func sanitizeCode(s int) string {
	switch s {
	case 100:
		return "100"
	case 101:
		return "101"

	case 200, 0:
		return "200"
	case 201:
		return "201"
	case 202:
		return "202"
	case 203:
		return "203"
	case 204:
		return "204"
	case 205:
		return "205"
	case 206:
		return "206"

	case 300:
		return "300"
	case 301:
		return "301"
	case 302:
		return "302"
	case 304:
		return "304"
	case 305:
		return "305"
	case 307:
		return "307"

	case 400:
		return "400"
	case 401:
		return "401"
	case 402:
		return "402"
	case 403:
		return "403"
	case 404:
		return "404"
	case 405:
		return "405"
	case 406:
		return "406"
	case 407:
		return "407"
	case 408:
		return "408"
	case 409:
		return "409"
	case 410:
		return "410"
	case 411:
		return "411"
	case 412:
		return "412"
	case 413:
		return "413"
	case 414:
		return "414"
	case 415:
		return "415"
	case 416:
		return "416"
	case 417:
		return "417"
	case 418:
		return "418"

	case 500:
		return "500"
	case 501:
		return "501"
	case 502:
		return "502"
	case 503:
		return "503"
	case 504:
		return "504"
	case 505:
		return "505"

	case 428:
		return "428"
	case 429:
		return "429"
	case 431:
		return "431"
	case 511:
		return "511"

	default:
		return strconv.Itoa(s)
	}
}
