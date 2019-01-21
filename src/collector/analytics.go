package collector

import (
	"fmt"

	"github.com/mtulio/azion-exporter/src/azion"
	"github.com/prometheus/client_golang/prometheus"
)

// Analytics keeps the collector info
type Analytics struct {
	AzionClient *azion.Client
	Metrics     []Metric
	// Metrics     map[string]*prometheus.Desc
}

type Metric struct {
	Prom *prometheus.Desc
	Name string
}

// NewCollectorAnalytics return the CollectorAnalytics object
func NewCollectorAnalytics(aCli *azion.Client) (*Analytics, error) {

	// metrics := make(map[string]*prometheus.Desc)
	var metrics []Metric

	m := Metric{
		Name: prometheus.BuildFQName(namespace, "requests", "total"),
	}
	m.Prom = prometheus.NewDesc(
		m.Name,
		"Azion Analytics Requests Total",
		nil, nil,
	)
	metrics = append(metrics, m)

	return &Analytics{
		AzionClient: aCli,
		Metrics:     metrics,
	}, nil
}

// Update implements Collector and exposes related metrics
func (ca *Analytics) Update(ch chan<- prometheus.Metric) error {
	if err := ca.updateMetrics(ch); err != nil {
		return err
	}
	return nil
}

func (ca *Analytics) updateMetrics(ch chan<- prometheus.Metric) error {
	for _, m := range ca.Metrics {
		fmt.Println(m.Name)
		ch <- prometheus.MustNewConstMetric(
			m.Prom,
			prometheus.GaugeValue,
			1,
		)
	}
	return nil
}
