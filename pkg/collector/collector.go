package collector

import "github.com/prometheus/client_golang/prometheus"

const namespace = "tmpl"

var (
	upDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the exporter able to contact the thing?",
		nil,
		nil,
	)
)

type Collector struct{}

func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- upDesc
}

func (c Collector) Collect(ch chan<- prometheus.Metric) {
	var success bool
	var up float64

	if success {
		up = 1
	}

	ch <- prometheus.MustNewConstMetric(upDesc, prometheus.GaugeValue, up)
}
