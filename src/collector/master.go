package collector

import (
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/mtulio/azion-exporter/src/azion"
	"github.com/prometheus/client_golang/prometheus"
)

// CollectorMaster implements the prometheus.Collector interface.
type CollectorMaster struct {
	Collectors  map[string]Collector
	AzionClient *azion.Client
}

// Collector is the interface a collector has to implement.
type Collector interface {
	// Get new metrics and expose them via prometheus registry.
	Update(ch chan<- prometheus.Metric) error
}

const (
	// Namespace defines the common namespace to be used by all metrics.
	namespace       = "azion"
	defaultEnabled  = true
	defaultDisabled = false
)

var (
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "collector_duration_seconds"),
		"azion_exporter: Duration of a collector scrape.",
		[]string{"collector"},
		nil,
	)
	scrapeSuccessDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "collector_success"),
		"azion_exporter: Whether a collector succeeded.",
		[]string{"collector"},
		nil,
	)
)

// NewCollectorMaster creates a new NodeCollector.
func NewCollectorMaster(azionCli *azion.Client, metrics ...string) (*CollectorMaster, error) {
	var err error
	err = nil
	collectors := make(map[string]Collector)
	collectors["analytics"], err = NewCollectorAnalytics(azionCli, metrics...)
	if err != nil {
		panic(err)
	}

	return &CollectorMaster{
		Collectors:  collectors,
		AzionClient: azionCli,
	}, nil
}

//Describe is a Prometheus implementation to be called by collector.
//It essentially writes all descriptors to the prometheus desc channel.
func (cm *CollectorMaster) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- scrapeDurationDesc
	ch <- scrapeSuccessDesc
}

//Collect implements required collect function for all promehteus collectors
func (cm *CollectorMaster) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(cm.Collectors))
	for name, c := range cm.Collectors {
		go func(name string, c Collector) {
			execute(name, c, ch)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
}

// execute calls Update() function on subsystem to gather metrics
func execute(name string, c Collector, ch chan<- prometheus.Metric) {
	begin := time.Now()
	err := c.Update(ch)
	duration := time.Since(begin)
	var success float64

	if err != nil {
		log.Errorf("ERROR: %s collector failed after %fs: %s", name, duration.Seconds(), err)
		success = 0
	} else {
		log.Debugf("OK: %s collector succeeded after %fs.", name, duration.Seconds())
		success = 1
	}
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(), name)
	ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, success, name)
}
