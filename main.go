package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/mtulio/azion-exporter/src/azion"
	"github.com/mtulio/azion-exporter/src/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

type globalProm struct {
	Collector *collector.CollectorMaster
	Registry  *prometheus.Registry
	Gatherers *prometheus.Gatherers
}

type configParams struct {
	azionEmail     *string
	azionPass      *string
	apiListenAddr  *string
	apiMetricsPath *string
	prom           *globalProm
	azionClient    *azion.Client
	metricsName    []string
	metricInterval *int
}

const (
	exporterName        = "azion_exporter"
	exporterDescription = "Azion Exporter"
)

var (
	cfg               = configParams{}
	defAPIListenAddr  = ":9801"
	defAPIMetricsPath = "/metrics"
	defMetricInterval = 60
)

// usage returns the command line usage sample.
func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [options]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {

	cfg.apiListenAddr = flag.String("web.listen-address", defAPIListenAddr, "Address on which to expose metrics and web interface.")
	cfg.apiMetricsPath = flag.String("web.telemetry-path", defAPIMetricsPath, "Path under which to expose metrics.")

	cfg.azionEmail = flag.String("azion.email", "", "API email address to get Authorization token")
	cfg.azionPass = flag.String("azion.password", "", "API password to get Authorization token")

	fMetricsFilter := flag.String("metrics.filter", "", "List of metrics sepparated by comma")
	cfg.metricInterval = flag.Int("metrics.interval", defMetricInterval, "Interval in seconds to retrieve metrics from API")

	flag.Usage = usage
	flag.Parse()

	if *cfg.azionEmail == "" {
		*cfg.azionEmail = os.Getenv("AZION_EMAIL")
	}

	if *cfg.azionPass == "" {
		*cfg.azionPass = os.Getenv("AZION_PASSWORD")
	}

	if len(*fMetricsFilter) > 0 {
		for _, m := range strings.Split(*fMetricsFilter, ",") {
			cfg.metricsName = append(cfg.metricsName, m)
		}
	}

	cfg.azionClient = azion.NewClient(*cfg.azionEmail, *cfg.azionPass)

	initPromCollector()
}

// Main Prometheus handler
func handler(w http.ResponseWriter, r *http.Request) {

	// Delegate http serving to Prometheus client library, which will call collector.Collect.
	h := promhttp.InstrumentMetricHandler(
		cfg.prom.Registry,
		promhttp.HandlerFor(cfg.prom.Gatherers,
			promhttp.HandlerOpts{
				// ErrorLog:      log.NewErrorLogger(),
				ErrorHandling: promhttp.ContinueOnError,
			}),
	)
	h.ServeHTTP(w, r)
}

func main() {
	log.Infoln("Starting exporter ")

	// This section will start the HTTP server and expose
	// any metrics on the /metrics endpoint.
	http.HandleFunc(*cfg.apiMetricsPath, handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>"` + exporterDescription + `"</title></head>
			<body>
			<h1>` + exporterDescription + `</h1>
			<p><br> The metrics is available on the path:
			<a href="` + *cfg.apiMetricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Info("Beginning to serve on port " + *cfg.apiListenAddr)
	log.Fatal(http.ListenAndServe(*cfg.apiListenAddr, nil))

}
