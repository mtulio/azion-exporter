package azion

// AnalyticsSvc handles communication with the Azion API methods related to
// Analytics.
type AnalyticsSvc struct {
	client  *Client
	BaseURI string
}

// Analytics represents a Azion Analytics.
type Analytics struct {
	Name *string `json:"name"`
	ID   *uint   `json:"id,omitempty"`
}

// AnalyticsMetric represents a Azion Analytics Metric.
type AnalyticsMetric struct {
	Name       *string `json:"name"`
	Dimensions []AnalyticsMetricDim
}

// AnalyticsMetricDim represents a Azion Analytics Metric dimensions.
type AnalyticsMetricDim map[string][]string

// MetricResp is a metric response payload returned by Azion API.
//
// Azion API docs: https://www.azion.com.br/developers/api-v2/analytics/
type MetricResp map[string]map[string]map[string]map[string][][]interface{}

// GetMatadata returns the metadata path values.
//
// Azion API docs: https://www.azion.com.br/developers/api-v2/analytics/
func (a *AnalyticsSvc) GetMatadata() (*AnalyticsMetricDim, error) {
	req, err := a.client.NewRequest("GET", "/analytics/metadata", nil)
	if err != nil {
		return nil, err
	}

	dimensionsResponse := new(AnalyticsMetricDim)

	_, err = a.client.Do(req, &dimensionsResponse)
	if err != nil {
		return nil, err
	}

	return dimensionsResponse, nil
}

// GetMetricDimension return the metric with dimensions
func (a *AnalyticsSvc) GetMetricDimension(pid, mc, dim string, qArgs ...string) (*MetricResp, error) {
	url := a.BaseURI + "/products/" + pid + "/aggregate/metrics/" + mc + "/dimensions/" + dim + "?"
	argCnt := 0
	for _, value := range qArgs {
		if argCnt == 0 {
			url += value
			argCnt++
		} else {
			url += "&" + value
		}
	}
	return a.getMetric(url)
}

// getMetric return the metric requested by URL
func (a *AnalyticsSvc) getMetric(url string) (*MetricResp, error) {

	req, err := a.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	metrics := new(MetricResp)

	_, err = a.client.Do(req, &metrics)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

// getProductID return the ProductID when an Alias is specified
func (a *AnalyticsSvc) getProductID(product string) string {
	switch product {
	case "ContentDelivery":
		return "1441740010"
	case "CloudStorage":
		return "1441740013"
	case "ImageOptimization":
		return "1441110021"
	case "LiveIngest":
		return "1467740028"
	case "MediaPackager":
		return "1441740014"
	default:
		return product
	}
}

//
// CD: Content Delivery (1441740010)
//

// GetMetricProdCDDimension return the metric with dimensions for product Content Delivery
func (a *AnalyticsSvc) GetMetricProdCDDimension(metric, dimension string, qArgs ...string) (*MetricResp, error) {
	return a.GetMetricDimension(a.getProductID("ContentDelivery"), metric, dimension, qArgs...)
}
