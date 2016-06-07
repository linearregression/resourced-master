package application

import (
	"github.com/rcrowley/go-metrics"
	"time"
)

func (app *Application) NewMetricsRegistry() metrics.Registry {
	r := metrics.NewRegistry()
	metrics.RegisterDebugGCStats(r)
	metrics.RegisterRuntimeMemStats(r)

	go metrics.CaptureDebugGCStats(r, time.Second*60)
	go metrics.CaptureRuntimeMemStats(r, time.Second*60)

	return r
}
