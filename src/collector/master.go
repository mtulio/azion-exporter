package collector

import (
	"sync"

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
func NewCollectorMaster(azionCli *azion.Client) (*CollectorMaster, error) {
	var err error
	err = nil
	collectors := make(map[string]Collector)
	collectors["analytics"], err = NewCollectorAnalytics(azionCli)
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
			// execute(name, c, ch)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
}
