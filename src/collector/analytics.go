package collector

import (
	"encoding/json"
	"fmt"
	"regexp"
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
		reStatusCode := regexp.MustCompile(`^cd_status_code_*`)
		reReq := regexp.MustCompile(`^cd_requests_*`)
		reBW := regexp.MustCompile(`^cd_bandwidth_*`)
		reDT := regexp.MustCompile(`^cd_data_transferred_*`)

		switch {
		case reReq.MatchString(mName):
			{
				m.Name = prometheus.BuildFQName(namespace, "cd", "requests_count")
				m.Description = "Azion Analytics Content Delivery Requests Count"
				m.Labels = []string{"type"}
				switch mName {
				case "cd_requests_total":
					{
						m.fCollector = ca.collectorWrapper("requests", "total")
						m.LabelsValue = []string{"total"}
					}
				case "cd_requests_saved":
					{
						m.fCollector = ca.collectorWrapper("requests", "saved")
						m.LabelsValue = []string{"saved"}
					}
				case "cd_requests_missed":
					{
						m.fCollector = ca.collectorWrapper("requests", "missed")
						m.LabelsValue = []string{"missed"}
					}
				}
			}
		case reBW.MatchString(mName):
			{
				m.Name = prometheus.BuildFQName(namespace, "cd", "bandwidth_gb")
				m.Description = "Azion Analytics Content Delivery Bandwidth Count"
				m.Labels = []string{"type"}
				switch mName {
				case "cd_bandwidth_total":
					{
						m.fCollector = ca.collectorWrapper("bandwidth", "total")
						m.LabelsValue = []string{"total"}
					}
				case "cd_bandwidth_saved":
					{
						m.fCollector = ca.collectorWrapper("bandwidth", "saved")
						m.LabelsValue = []string{"saved"}
					}
				case "cd_bandwidth_missed":
					{
						m.fCollector = ca.collectorWrapper("bandwidth", "missed")
						m.LabelsValue = []string{"missed"}
					}
				}
			}
		case reDT.MatchString(mName):
			{
				m.Name = prometheus.BuildFQName(namespace, "cd", "data_transferred_mb")
				m.Description = "Azion Analytics Content Delivery Data Transferred in MB"
				m.Labels = []string{"type"}
				switch mName {
				case "cd_data_transferred_total":
					{
						m.fCollector = ca.collectorWrapper("data_transferred", "total")
						m.LabelsValue = []string{"total"}
					}
				case "cd_data_transferred_saved":
					{
						m.fCollector = ca.collectorWrapper("data_transferred", "saved")
						m.LabelsValue = []string{"saved"}
					}
				case "cd_data_transferred_missed":
					{
						m.fCollector = ca.collectorWrapper("data_transferred", "missed")
						m.LabelsValue = []string{"missed"}
					}
				}
			}
		case reStatusCode.MatchString(mName):
			{
				m.Name = prometheus.BuildFQName(namespace, "cd", "status_code_total")
				m.Description = "Azion Analytics Content Delivery Status Code 5xx Total"
				m.Labels = []string{"code"}
				switch mName {
				case "cd_status_code_2xx":
					{
						m.fCollector = ca.collectorWrapper("status_code", "2xx")
						m.LabelsValue = []string{"2xx"}
					}
				case "cd_status_code_200":
					{
						m.fCollector = ca.collectorWrapper("status_code", "200")
						m.LabelsValue = []string{"200"}
					}
				case "cd_status_code_204":
					{
						m.fCollector = ca.collectorWrapper("status_code", "204")
						m.LabelsValue = []string{"204"}
					}
				case "cd_status_code_206":
					{
						m.fCollector = ca.collectorWrapper("status_code", "206")
						m.LabelsValue = []string{"206"}
					}
				case "cd_status_code_3xx":
					{
						m.fCollector = ca.collectorWrapper("status_code", "3xx")
						m.LabelsValue = []string{"3xx"}
					}
				case "cd_status_code_301":
					{
						m.fCollector = ca.collectorWrapper("status_code", "301")
						m.LabelsValue = []string{"301"}
					}
				case "cd_status_code_302":
					{
						m.fCollector = ca.collectorWrapper("status_code", "302")
						m.LabelsValue = []string{"302"}
					}
				case "cd_status_code_304":
					{
						m.fCollector = ca.collectorWrapper("status_code", "304")
						m.LabelsValue = []string{"304"}
					}
				case "cd_status_code_4xx":
					{
						m.fCollector = ca.collectorWrapper("status_code", "4xx")
						m.LabelsValue = []string{"4xx"}
					}
				case "cd_status_code_400":
					{
						m.fCollector = ca.collectorWrapper("status_code", "400")
						m.LabelsValue = []string{"400"}
					}
				case "cd_status_code_403":
					{
						m.fCollector = ca.collectorWrapper("status_code", "403")
						m.LabelsValue = []string{"403"}
					}
				case "cd_status_code_404":
					{
						m.fCollector = ca.collectorWrapper("status_code", "404")
						m.LabelsValue = []string{"404"}
					}
				case "cd_status_code_5xx":
					{
						m.fCollector = ca.collectorWrapper("status_code", "5xx")
						m.LabelsValue = []string{"5xx"}
					}
				case "cd_status_code_500":
					{
						m.fCollector = ca.collectorWrapper("status_code", "500")
						m.LabelsValue = []string{"500"}
					}
				case "cd_status_code_502":
					{
						m.fCollector = ca.collectorWrapper("status_code", "502")
						m.LabelsValue = []string{"502"}
					}
				case "cd_status_code_503":
					{
						m.fCollector = ca.collectorWrapper("status_code", "503")
						m.LabelsValue = []string{"503"}
					}

				}
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

func (ca *Analytics) collectorMetric(n, d string, args ...string) ([]byte, error) {

	mData, err := ca.AzionClient.Analytics.GetMetricDimension(n, d, args...)
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

func (ca *Analytics) collectorWrapper(metric, dimension string) func(m *Metric) error {
	return func(m *Metric) error {

		b, err := ca.collectorMetric(metric, dimension, "date_from=last-hour")
		if err != nil {
			return err
		}

		md := metric + "_" + dimension
		switch md {
		case "data_transferred_total":
			{
				var mI azion.MetricRespDTTotal
				// Casting metric payload
				json.Unmarshal(b, &mI)
				v, err := ca.metricAssertion(mI.Products.Prod1441740010.Metric.Dimension)
				if err != nil {
					return nil
				}
				m.Value = v
			}
		case "data_transferred_saved":
			{
				var mI azion.MetricRespDTSaved
				// Casting metric payload
				json.Unmarshal(b, &mI)
				v, err := ca.metricAssertion(mI.Products.Prod1441740010.Metric.Dimension)
				if err != nil {
					return nil
				}
				m.Value = v
			}
		case "data_transferred_missed":
			{
				var mI azion.MetricRespDTMissed
				// Casting metric payload
				json.Unmarshal(b, &mI)
				v, err := ca.metricAssertion(mI.Products.Prod1441740010.Metric.Dimension)
				if err != nil {
					return nil
				}
				m.Value = v
			}
		case "bandwidth_total":
			{
				var mI azion.MetricRespBWTotal
				// Casting metric payload
				json.Unmarshal(b, &mI)
				v, err := ca.metricAssertion(mI.Products.Prod1441740010.Metric.Dimension)
				if err != nil {
					return nil
				}
				m.Value = v
			}
		case "bandwidth_saved":
			{
				var mI azion.MetricRespBWSaved
				// Casting metric payload
				json.Unmarshal(b, &mI)
				v, err := ca.metricAssertion(mI.Products.Prod1441740010.Metric.Dimension)
				if err != nil {
					return nil
				}
				m.Value = v
			}
		case "bandwidth_missed":
			{
				var mI azion.MetricRespBWMissed
				// Casting metric payload
				json.Unmarshal(b, &mI)
				v, err := ca.metricAssertion(mI.Products.Prod1441740010.Metric.Dimension)
				if err != nil {
					return nil
				}
				m.Value = v
			}
		case "requests_total":
			{
				var mI azion.MetricRespReqTotal
				// Casting metric payload
				json.Unmarshal(b, &mI)
				v, err := ca.metricAssertion(mI.Products.Prod1441740010.Metric.Dimension)
				if err != nil {
					return nil
				}
				m.Value = v
			}
		case "requests_missed":
			{
				var mI azion.MetricRespReqMissed
				// Casting metric payload
				json.Unmarshal(b, &mI)
				v, err := ca.metricAssertion(mI.Products.Prod1441740010.Metric.Dimension)
				if err != nil {
					return nil
				}
				m.Value = v
			}
		case "requests_saved":
			{
				var mI azion.MetricRespReqSaved
				// Casting metric payload
				json.Unmarshal(b, &mI)
				v, err := ca.metricAssertion(mI.Products.Prod1441740010.Metric.Dimension)
				if err != nil {
					return nil
				}
				m.Value = v
			}
		case "status_code_2xx":
			{
				var mI azion.MetricRespSCode2xx
				// Casting metric payload
				json.Unmarshal(b, &mI)
				v, err := ca.metricAssertion(mI.Products.Prod1441740010.Metric.Dimension)
				if err != nil {
					return nil
				}
				m.Value = v
			}
		case "status_code_200":
			{
				var mI azion.MetricRespSCode200
				// Casting metric payload
				json.Unmarshal(b, &mI)
				v, err := ca.metricAssertion(mI.Products.Prod1441740010.Metric.Dimension)
				if err != nil {
					return nil
				}
				m.Value = v
			}
		case "status_code_204":
			{
				var mI azion.MetricRespSCode204
				// Casting metric payload
				json.Unmarshal(b, &mI)
				v, err := ca.metricAssertion(mI.Products.Prod1441740010.Metric.Dimension)
				if err != nil {
					return nil
				}
				m.Value = v
			}
		case "status_code_206":
			{
				var mI azion.MetricRespSCode206
				// Casting metric payload
				json.Unmarshal(b, &mI)
				v, err := ca.metricAssertion(mI.Products.Prod1441740010.Metric.Dimension)
				if err != nil {
					return nil
				}
				m.Value = v
			}
		case "status_code_3xx":
			{
				var mI azion.MetricRespSCode3xx
				// Casting metric payload
				json.Unmarshal(b, &mI)
				v, err := ca.metricAssertion(mI.Products.Prod1441740010.Metric.Dimension)
				if err != nil {
					return nil
				}
				m.Value = v
			}
		case "status_code_301":
			{
				var mI azion.MetricRespSCode301
				// Casting metric payload
				json.Unmarshal(b, &mI)
				v, err := ca.metricAssertion(mI.Products.Prod1441740010.Metric.Dimension)
				if err != nil {
					return nil
				}
				m.Value = v
			}
		case "status_code_302":
			{
				var mI azion.MetricRespSCode302
				// Casting metric payload
				json.Unmarshal(b, &mI)
				v, err := ca.metricAssertion(mI.Products.Prod1441740010.Metric.Dimension)
				if err != nil {
					return nil
				}
				m.Value = v
			}
		case "status_code_304":
			{
				var mI azion.MetricRespSCode304
				// Casting metric payload
				json.Unmarshal(b, &mI)
				v, err := ca.metricAssertion(mI.Products.Prod1441740010.Metric.Dimension)
				if err != nil {
					return nil
				}
				m.Value = v
			}
		case "status_code_4xx":
			{
				var mI azion.MetricRespSCode4xx
				// Casting metric payload
				json.Unmarshal(b, &mI)
				v, err := ca.metricAssertion(mI.Products.Prod1441740010.Metric.Dimension)
				if err != nil {
					return nil
				}
				m.Value = v
			}
		case "status_code_400":
			{
				var mI azion.MetricRespSCode400
				// Casting metric payload
				json.Unmarshal(b, &mI)
				v, err := ca.metricAssertion(mI.Products.Prod1441740010.Metric.Dimension)
				if err != nil {
					return nil
				}
				m.Value = v
			}
		case "status_code_403":
			{
				var mI azion.MetricRespSCode403
				// Casting metric payload
				json.Unmarshal(b, &mI)
				v, err := ca.metricAssertion(mI.Products.Prod1441740010.Metric.Dimension)
				if err != nil {
					return nil
				}
				m.Value = v
			}
		case "status_code_404":
			{
				var mI azion.MetricRespSCode404
				// Casting metric payload
				json.Unmarshal(b, &mI)
				v, err := ca.metricAssertion(mI.Products.Prod1441740010.Metric.Dimension)
				if err != nil {
					return nil
				}
				m.Value = v
			}
		case "status_code_5xx":
			{
				var mI azion.MetricRespSCode5xx
				// Casting metric payload
				json.Unmarshal(b, &mI)
				v, err := ca.metricAssertion(mI.Products.Prod1441740010.Metric.Dimension)
				if err != nil {
					return nil
				}
				m.Value = v
			}
		case "status_code_502":
			{
				var mI azion.MetricRespSCode502
				// Casting metric payload
				json.Unmarshal(b, &mI)
				v, err := ca.metricAssertion(mI.Products.Prod1441740010.Metric.Dimension)
				if err != nil {
					return nil
				}
				m.Value = v
			}
		case "status_code_503":
			{
				var mI azion.MetricRespSCode503
				// Casting metric payload
				json.Unmarshal(b, &mI)
				v, err := ca.metricAssertion(mI.Products.Prod1441740010.Metric.Dimension)
				if err != nil {
					return nil
				}
				m.Value = v
			}
		}
		return nil
	}
}
