package azion

// AnalyticsService handles communication with the Azion API methods related to
// Analytics.
type AnalyticsService struct {
	client *Client
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
type AnalyticsMetricDim map[string][]map[string]string

// GetMatadata returns the metadata path values.
//
// Azion API docs: https://www.azion.com.br/developers/api-v2/analytics/
func (a *AnalyticsService) GetMatadata() (*AnalyticsMetricDim, error) {
	req, err := a.client.NewRequest("GET", "/analytics/metadata", nil)
	if err != nil {
		return nil, err
	}

	var dimensionsResponse struct {
		Dimensions AnalyticsMetricDim
	}

	_, err = a.client.Do(req, &dimensionsResponse)
	if err != nil {
		return nil, err
	}

	return &dimensionsResponse.Dimensions, nil
}
