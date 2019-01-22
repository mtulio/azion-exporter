package collector

import (
	"encoding/json"
	"fmt"
	"sync"
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
	LabelsValue []string
	LabelsConst prometheus.Labels
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
	// done := make(chan bool)
	wg := sync.WaitGroup{}
	wg.Add(len(ca.Metrics))

	for mID := range ca.Metrics {
		go func(m *Metric, ch chan<- prometheus.Metric) {
			if m.LabelsValue != nil {
				ch <- prometheus.MustNewConstMetric(
					m.Prom,
					prometheus.GaugeValue,
					m.Value,
					m.LabelsValue...,
				)
			} else {
				ch <- prometheus.MustNewConstMetric(
					m.Prom,
					prometheus.GaugeValue,
					m.Value,
					m.LabelsValue...,
				)
			}
			// done <- true
			wg.Done()
		}(&ca.Metrics[mID], ch)
	}

	// wait to finish all go routines
	wg.Wait()
	// <-done
	return nil
}

// InitMetrics initialize a list of metrics names and return error if fails.
func (ca *Analytics) InitMetrics(msEnabled ...string) error {

	for _, mName := range msEnabled {
		m := Metric{}
		switch mName {
		case "cd_requests_total":
			{
				m.Name = prometheus.BuildFQName(namespace, "cd", "requests_total")
				m.Description = "Azion Analytics Content Delivery Requests Total"
				m.fCollector = ca.collectorRequestsTotal()
			}
		case "cd_status_code_5xx":
			{
				m.Name = prometheus.BuildFQName(namespace, "cd", "status_code_total_5xx")
				m.Description = "Azion Analytics Content Delivery Status Code 5xx Total"
				m.fCollector = ca.collectorStatusCode5xx()
				// m.Labels = []string{"code"}
				// m.LabelsValue = []string{"5xx"}
				// m.LabelsConst = prometheus.Labels{"code": "5xx"}
			}
		case "cd_status_code_500":
			{
				m.Name = prometheus.BuildFQName(namespace, "cd", "status_code_total_500")
				m.Description = "Azion Analytics Content Delivery Status Code 500 Total"
				m.fCollector = ca.collectorStatusCode500()
				// m.Labels = []string{"code"}
				// m.LabelsValue = []string{"500"}
				// m.LabelsConst = prometheus.Labels{"code": "500"}
			}
		case "cd_status_code_502":
			{
				m.Name = prometheus.BuildFQName(namespace, "cd", "status_code_total_502")
				m.Description = "Azion Analytics Content Delivery Status Code 502 Total"
				m.fCollector = ca.collectorStatusCode502()
				// m.Labels = []string{"code"}
				// m.LabelsValue = []string{"502"}
				// m.LabelsConst = prometheus.Labels{"code": "502"}
			}
		case "cd_status_code_503":
			{
				m.Name = prometheus.BuildFQName(namespace, "cd", "status_code_total_503")
				m.Description = "Azion Analytics Content Delivery Status Code 503 Total"
				m.fCollector = ca.collectorStatusCode503()
				// m.Labels = []string{"code"}
				// m.LabelsValue = []string{"503"}
				// m.LabelsConst = prometheus.Labels{"code": "500"}
			}
		default:
			fmt.Println("Metric init Error, metric definition found: ", mName)
			continue
		}
		m.Prom = prometheus.NewDesc(
			m.Name,
			m.Description,
			m.Labels, m.LabelsConst,
		)
		ca.Metrics = append(ca.Metrics, m)
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
// Metrics mapping / parser / cast
//

// metricAssertion asserts the datapoints to retrieve last valid value.
// BUG Report: Azion Analytics API has delays to proccess latest datapoints,
// the last one is always lower, sometimes more than it,
// to prevent empty metrics, we will follow the strategy:
// - we consider >=2min datapoint an 'safe value'; if it's <=0, then
// - get the latest (>=2min) data point greater than 0;
// The value will be: >= 2 min ago && > 0.
func (ca *Analytics) metricAssertion(datapoints [][]interface{}) (float64, error) {

	value := 0.0
	posLatestDP := len(datapoints) - 2
	for i := posLatestDP; i >= 0; i-- {
		value = datapoints[i][1].(float64)
		if value > 0 {
			break
		}
	}
	return value, nil

}

// collectorRequestsTotal gather metrics from:
// - Product: Contend Delivery
// - Metric: requests
// - Dimension: total
// - Time: last-hour
func (ca *Analytics) collectorRequestsTotal() func(m *Metric) error {

	return func(m *Metric) error {

		var mI mIndexing
		mData, err := ca.AzionClient.Analytics.GetMetricDimensionProdCD("requests", "total", "date_from=last-hour")
		if err != nil {
			log.Info("Error getting metrics. Ignoring")
			return err
		}

		b, err := json.Marshal(mData)
		if err != nil {
			return err
		}
		// Asserting to ignore last datapoint that has "uncomplete" data.
		// Gathering '2 min ago' datapoint.
		// Casting metric payload {"products":{"1441740010":{"requests":{"total":[[T,V]]}}}}
		json.Unmarshal(b, &mI)

		// BUG Report: Azion Analytics API has delays to proccess latest datapoints,
		// the last one is always lower, sometimes more than it,
		// to prevent empty metrics, we will follow the strategy:
		// - we consider >=2min datapoint an 'safe value'; if it's <=0, then
		// - get the latest (>=2min) data point greater than 0;
		// The value will be: >= 2 min ago && > 0.
		posLatestDP := len(mI.Products.Num1441740010.Requests.Total) - 2
		for i := posLatestDP; i >= 0; i-- {
			m.Value = (mI.Products.Num1441740010.Requests.Total[i][1]).(float64)
			if m.Value > 0 {
				break
			}
		}
		return nil
	}
}

func (ca *Analytics) collectorStatusCode(n, d string, args ...string) ([]byte, error) {

	mData, err := ca.AzionClient.Analytics.GetMetricDimensionProdCD(n, d, args...)
	if err != nil {
		log.Info("Error getting metrics from API. Name ")
		return nil, err
	}
	b, err := json.Marshal(mData)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (ca *Analytics) collectorStatusCode5xx() func(m *Metric) error {
	return func(m *Metric) error {
		var mI mIndexing5xx
		b, err := ca.collectorStatusCode("status_code", "5xx", "date_from=last-hour")
		if err != nil {
			return err
		}

		// Casting metric payload
		json.Unmarshal(b, &mI)

		v, err := ca.metricAssertion(mI.Products.Num1441740010.StatusCode.Code)
		if err != nil {
			return nil
		}
		m.Value = v
		return nil
	}
}

func (ca *Analytics) collectorStatusCode500() func(m *Metric) error {
	return func(m *Metric) error {
		var mI mIndexing500
		b, err := ca.collectorStatusCode("status_code", "500", "date_from=last-hour")
		if err != nil {
			return err
		}

		// Casting metric payload
		json.Unmarshal(b, &mI)

		v, err := ca.metricAssertion(mI.Products.Num1441740010.StatusCode.Code)
		if err != nil {
			return nil
		}
		m.Value = v
		return nil
	}
}

func (ca *Analytics) collectorStatusCode502() func(m *Metric) error {
	return func(m *Metric) error {
		var mI mIndexing502
		b, err := ca.collectorStatusCode("status_code", "502", "date_from=last-hour")
		if err != nil {
			return err
		}

		// Casting metric payload
		json.Unmarshal(b, &mI)

		v, err := ca.metricAssertion(mI.Products.Num1441740010.StatusCode.Code)
		if err != nil {
			return nil
		}
		m.Value = v
		return nil
	}
}

func (ca *Analytics) collectorStatusCode503() func(m *Metric) error {
	return func(m *Metric) error {
		var mI mIndexing503
		b, err := ca.collectorStatusCode("status_code", "503", "date_from=last-hour")
		if err != nil {
			return err
		}

		// Casting metric payload
		json.Unmarshal(b, &mI)

		v, err := ca.metricAssertion(mI.Products.Num1441740010.StatusCode.Code)
		if err != nil {
			return nil
		}
		m.Value = v
		return nil
	}
}
