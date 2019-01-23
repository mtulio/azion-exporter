package azion

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	libraryVersion   = "0.1.0"
	apiVersion       = "2"
	defaultBaseURL   = "https://api.azionapi.net/"
	userAgent        = "azion-go-sdk/" + libraryVersion
	defaultMediaType = "application/json; version=" + apiVersion
)

// A Client manages communication with the API.
type Client struct {
	// HTTP client used to communicate with the API
	client *http.Client

	// Headers to attach to every request made with the client. Headers will be
	// used to provide API authentication details and other necessary
	// headers.
	Headers map[string]string

	// Email and Password contains the authentication details needed to authenticate
	// against the API.
	Email, Password string

	// Token contains manage the Token lifecycle to authenticate against the API.
	Token *clientToken

	// Base URL for API requests. Defaults to the public API, but can be
	// set to an alternate endpoint if necessary. BaseURL should always be
	// terminated by a slash.
	BaseURL *url.URL

	// User agent used when communicating with the API.
	UserAgent string

	// Services used to manipulate API entities.
	Analytics *AnalyticsSvc
	// CloudSecurity   *CloudSecurity
	// ContentDelivery *ContentDelivery
	// RealTimePurge   *RealTimePurge
}

// clientToken manages authorization token in the API
type clientToken struct {
	Token          string `json:"token"`
	CreatedAt      string `json:"created_at"`
	ExpiresAt      string `json:"expires_at"`
	ExpirationDate time.Time
}

// NewClient returns a new Azion API client bound to the public Azion API.
func NewClient(email, password string) *Client {
	bu, err := url.Parse(defaultBaseURL)
	if err != nil {
		panic("Default API base URL couldn't be parsed")
	}

	return NewClientWithBaseURL(bu, email, password)
}

// NewClientWithBaseURL returned a new Azion API client with a custom base URL.
func NewClientWithBaseURL(baseURL *url.URL, email, password string) *Client {
	headers := map[string]string{
		"Content-Type": defaultMediaType,
		"Accept":       defaultMediaType,
	}

	c := &Client{
		client:    http.DefaultClient,
		Headers:   headers,
		Email:     email,
		Password:  password,
		BaseURL:   baseURL,
		UserAgent: userAgent,
	}

	c.Analytics = &AnalyticsSvc{
		client:  c,
		BaseURI: "/analytics",
	}
	// c.CloudSecurity = &CloudSecurityService{client: c}
	// c.ContentDelivery = &ContentDeliveryService{client: c}
	// c.RealTimePurge = &RealTimePurgeService{client: c}

	return c
}

// getTokenBase64 return Base64 string from string arguments.
func getBase64(e, p string) string {
	data := []byte(e + ":" + p)
	return base64.StdEncoding.EncodeToString(data)
}

// timeParser return time struct to unstandarized Azion token's response.
func timeParser(strDate string) (time.Time, error) {

	dateTime := strings.Split(strDate, " ")
	ymd := strings.Split(dateTime[0], "-")
	hmsn := strings.Split(dateTime[1], ":")

	year, _ := strconv.Atoi(ymd[0])
	month, _ := strconv.Atoi(ymd[1])
	day, _ := strconv.Atoi(ymd[2])
	hour, _ := strconv.Atoi(hmsn[0])
	min, _ := strconv.Atoi(hmsn[1])

	t := time.Date(year, time.Month(month), day, hour, min, 0, 0, time.UTC)

	return t, nil
}

// tokenRequest make one http request to renew an Token.
//
// API doc: https://www.azion.com.br/developers/api-v2/authentication/
func (c *Client) tokenRequest(v interface{}) error {
	req, err := c.NewRequest("POST", "/tokens", nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.Email, c.Password)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		return err
	}

	if v != nil && resp.ContentLength != 0 {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return nil
}

// tokenRenew renew an Token and return error if it fails.
//
// API doc: https://www.azion.com.br/developers/api-v2/authentication/
func (c *Client) tokenRenew() error {

	type reqToken struct {
		Token     string `json:"token"`
		CreatedAt string `json:"created_at"`
		ExpiresAt string `json:"expires_at"`
	}
	tokenResp := new(reqToken)

	err := c.tokenRequest(tokenResp)
	if err != nil {
		return err
	}

	t, err := timeParser(tokenResp.ExpiresAt)
	if err != nil {
		return err
	}

	c.Token = new(clientToken)
	if c.Token == nil {
		panic("Unable to allocate Token object")
	}
	c.Token.ExpirationDate = t
	c.Token.Token = tokenResp.Token

	fmt.Println("Success getting an token: ", c.Token.Token)
	return nil
}

// tokenValidation return check if token is expired and renew it.
//
// API doc: https://www.azion.com.br/developers/api-v2/authentication/
func (c *Client) tokenValidation() error {
	if c.Token == nil {
		return c.tokenRenew()
	}
	// check if token is not valid
	tExpired := time.Now().After(c.Token.ExpirationDate)
	if tExpired {
		return c.tokenRenew()
	}

	return nil
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body. If specified, the map provided by headers will be used to
// update request headers.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		encodeErr := json.NewEncoder(buf).Encode(body)
		if encodeErr != nil {
			return nil, encodeErr
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

// Do sends an API request and returns the API response.  The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.  If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {

	// req.SetBasicAuth(c.Email, c.Password)
	errT := c.tokenValidation()
	if errT != nil {
		return nil, errT
	}
	req.Header.Set("Authorization", "Token "+c.Token.Token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil && resp.ContentLength != 0 {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return resp, err
}

// ErrorResponse reports an error caused by an API request.
// ErrorResponse implements the Error interface.
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response

	// Error messages produces by Azion API.
	Errors ErrorResponseMessages `json:"errors"`
}

// ErrorResponseMessages contains error messages returned from the Azion API.
type ErrorResponseMessages struct {
	Params  map[string]interface{} `json:"params,omitempty"`
	Request []string               `json:"request,omitempty"`
	System  []string               `json:"system,omitempty"`
}

// CheckResponse checks the API response for errors; and returns them if
// present. A Response is considered an error if it has a status code outside
// the 2XX range.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}

	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}

	return nil
}
