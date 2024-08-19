/*
Copyright 2024 Thomas Helander

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"net/http"
	"os"

	kingpin "github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	sdtdclient "github.com/thelande/sdtd_client/pkg/sdtd_client"
	"github.com/thelande/sdtd_exporter/pkg/collector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

const (
	exporterName  = "sdtd_exporter"
	exporterTitle = "Go Exporter Template"
)

var (
	apiUrl = kingpin.Flag(
		"server.url",
		"The base URL of the 7 Days to Die dedicated server web API (e.g., http://127.0.0.1:8080).",
	).Default("http://127.0.0.1:8080").Envar("SDTD_API_URL").String()

	tokenName = kingpin.Flag(
		"server.token-name",
		"The name of the API token to use to authenticate with the web server.",
	).Envar("SDTD_TOKEN_NAME").String()

	tokenSecret = kingpin.Flag(
		"server.token-secret",
		"The secret of the API token to use to authenticate with the web server.",
	).Envar("SDTD_TOKEN_SECRET").String()

	metricsPath = kingpin.Flag(
		"web.telemetry-path",
		"Path under which to expose metrics.",
	).Default("/metrics").String()
	webConfig = webflag.AddFlags(kingpin.CommandLine, ":9816")
	logger    log.Logger
)

func main() {
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Version(version.Print(exporterName))
	kingpin.Parse()

	logger = promlog.New(promlogConfig)
	level.Info(logger).Log("msg", fmt.Sprintf("Starting %s", exporterName), "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "build_context", version.BuildContext())

	client, err := sdtdclient.NewSDTDClient(
		*apiUrl,
		&sdtdclient.SDTDAuth{
			TokenName:   *tokenName,
			TokenSecret: *tokenSecret,
		},
		true,
		&logger,
	)
	if err != nil {
		panic(err)
	}
	collector := collector.Collector{Client: client, Logger: &logger}

	// Uncomment the following two lines and comment out prometheus.MustRegister(collector)
	// to exclude the go metrics. Make sure to swap line 88 and 89 as well.
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)
	// prometheus.MustRegister(collector)

	landingConfig := web.LandingConfig{
		Name:        exporterTitle,
		Description: "Prometheus go-based Exporter",
		Version:     version.Info(),
		Links: []web.LandingLinks{
			{
				Address: *metricsPath,
				Text:    "Metrics",
			},
		},
	}
	landingPage, err := web.NewLandingPage(landingConfig)
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	http.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	// http.Handle(*metricsPath, promhttp.Handler())
	http.Handle("/", landingPage)

	srv := &http.Server{}
	if err := web.ListenAndServe(srv, webConfig, logger); err != nil {
		level.Error(logger).Log("msg", "HTTP listener stopped", "error", err)
		os.Exit(1)
	}
}
