package collector

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"

	sdtdclient "github.com/thelande/sdtd_client/pkg/sdtd_client"
)

const namespace = "sdtd"

func GetMetricDesc(name, unit, desc string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(namespace, name, unit),
		desc, nil, nil,
	)
}

var (
	upDesc       = GetMetricDesc("", "up", "Is the server reachable?")
	playersDesc  = GetMetricDesc("player", "count", "The number of online players.")
	zombiesDesc  = GetMetricDesc("zombie", "count", "The number active zombies.")
	animalsDesc  = GetMetricDesc("animal", "count", "The number active animals.")
	gametimeDesc = GetMetricDesc("gametime", "seconds", "The current time in the game.")
)

type Collector struct {
	Client *sdtdclient.SDTDClient
	Logger *log.Logger
}

func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range []*prometheus.Desc{
		upDesc,
		playersDesc,
		zombiesDesc,
		animalsDesc,
		gametimeDesc,
	} {
		ch <- desc
	}
}

func (c Collector) SetUp(ch chan<- prometheus.Metric, up bool) {
	val := 0.0
	if up {
		val = 1
	}
	ch <- prometheus.MustNewConstMetric(upDesc, prometheus.GaugeValue, val)
}

func (c Collector) Collect(ch chan<- prometheus.Metric) {
	if err := c.Client.Connect(); err != nil {
		level.Warn(*c.Logger).Log("msg", "Failed to connect to API server", "err", err)
		c.SetUp(ch, false)
		return
	}

	c.SetUp(ch, true)

	stats, err := c.Client.GetServerStats()
	if err != nil {
		level.Warn(*c.Logger).Log("msg", "Failed to get server stats", "err", err)
	} else {
		ch <- prometheus.MustNewConstMetric(playersDesc, prometheus.GaugeValue, float64(stats.Data.Players))
		ch <- prometheus.MustNewConstMetric(zombiesDesc, prometheus.GaugeValue, float64(stats.Data.Hostiles))
		ch <- prometheus.MustNewConstMetric(animalsDesc, prometheus.GaugeValue, float64(stats.Data.Animals))

		var gametime float64
		gametime += float64(stats.Data.GameTime.Days * 86400)
		gametime += float64(stats.Data.GameTime.Hours * 3600)
		gametime += float64(stats.Data.GameTime.Minutes * 60)
		ch <- prometheus.MustNewConstMetric(gametimeDesc, prometheus.CounterValue, gametime)
	}
}
