package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

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
	azionEmail  *string
	azionPass   *string
	apiListen   string
	metricsPath string
	prom        *globalProm
	azionClient *azion.Client
	metricsName []string
}

const (
	exporterName        = "azion_exporter"
	exporterDescription = "Azion Exporter"
)

var (
	cfg = configParams{
		apiListen:   ":9801",
		metricsPath: "/metrics",
	}
)

// usage returns the command line usage sample.
func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [options]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	cfg.azionEmail = flag.String("azion.email", "", "API email address to get Authorization token")
	cfg.azionPass = flag.String("azion.password", "", "API password to get Authorization token")
	flag.Usage = usage
	flag.Parse()

	if *cfg.azionEmail == "" {
		*cfg.azionEmail = os.Getenv("AZION_EMAIL")
	}

	if *cfg.azionPass == "" {
		*cfg.azionPass = os.Getenv("AZION_PASSWORD")
	}

	// List of metrics to retrieve
	cfg.metricsName = append(cfg.metricsName, "cd_requests_total")

	cfg.azionClient = azion.NewClient(*cfg.azionEmail, *cfg.azionPass)

	initPromCollector()
	// Samples
	// sampleGetMetadata(c)
	// sampleGetMetricProdCDDimension(c)
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
	http.HandleFunc(cfg.metricsPath, handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>"` + exporterDescription + `"</title></head>
			<body>
			<h1>"` + exporterDescription + `"</h1>
			<p><a href="` + cfg.metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Info("Beginning to serve on port " + cfg.apiListen)
	log.Fatal(http.ListenAndServe(cfg.apiListen, nil))

}
