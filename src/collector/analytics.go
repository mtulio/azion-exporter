package collector

import (
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/mtulio/azion-exporter/src/azion"
	"github.com/prometheus/client_golang/prometheus"
)

// Analytics keeps the collector info
type Analytics struct {
	AzionClient *azion.Client
	Metrics     []Metric
}

// Metric describe the metric attributes
type Metric struct {
	Prom        *prometheus.Desc
	Name        string
	Description string
	fCollector  func(m *Metric) error
	Value       float64
	Labels      []string
}

// NewCollectorAnalytics return the CollectorAnalytics object
func NewCollectorAnalytics(aCli *azion.Client, msEnabled ...string) (*Analytics, error) {

	ca := &Analytics{
		AzionClient: aCli,
	}
	err := ca.InitMetrics(msEnabled...)
	if err != nil {
		log.Info("collector.Analytics: error initializing metrics")
	}
	go ca.InitCollectorsUpdater()
	return ca, nil
}

// Update implements Collector and exposes related metrics
func (ca *Analytics) Update(ch chan<- prometheus.Metric) error {
	done := make(chan bool)
	for mID := range ca.Metrics {
		go func(m *Metric, ch chan<- prometheus.Metric) {
			fmt.Println(m.Name)
			ch <- prometheus.MustNewConstMetric(
				m.Prom,
				prometheus.GaugeValue,
				m.Value,
				m.Labels...,
			)
			done <- true
		}(&ca.Metrics[mID], ch)
	}

	// wait to finish all go routines
	<-done
	return nil
}

// InitMetrics initialize a list of metrics names and return error if fails.
func (ca *Analytics) InitMetrics(msEnabled ...string) error {

	for _, mName := range msEnabled {
		switch mName {
		case "cd_requests_total":
			{
				m := Metric{
					Name:        prometheus.BuildFQName(namespace, "cd", "requests_total"),
					Description: "Azion Analytics Content Delivery Requests Total",
				}
				m.Prom = prometheus.NewDesc(
					m.Name,
					m.Description,
					m.Labels, nil,
				)
				m.fCollector = ca.collectorRequestsTotal()
				ca.Metrics = append(ca.Metrics, m)
			}
		}
	}
	return nil
}

// InitCollectorsUpdater start the paralel auto update for each collector
func (ca *Analytics) InitCollectorsUpdater() {
	for {
		for mID := range ca.Metrics {
			go func(m *Metric) {
				m.fCollector(m)
			}(&ca.Metrics[mID])
		}
		time.Sleep(time.Second * time.Duration(60))
	}
}

//
// Metrics mapping
//

// collectorRequestsTotal gather metrics from:
// - Product: Contend Delivery
// - Metric: requests
// - Dimension: total
// - Time: last-hour
func (ca *Analytics) collectorRequestsTotal() func(m *Metric) error {

	return func(m *Metric) error {
		mData, err := ca.AzionClient.Analytics.GetMetricDimensionProdCD("requests", "total", "date_from=last-hour")
		if err != nil {
			log.Info("Error getting metrics. Ignoring")
			return err
		}
		fmt.Println(mData)
		return nil
	}
}
