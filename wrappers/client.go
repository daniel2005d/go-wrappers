package wrappers

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

// HTTPClient wraps http.Client with a timeout and redirect policy.
// Redirects are not followed so callers can inspect 3xx responses and Location.
type HTTPClient struct {
	client         *http.Client
	defaultHeaders http.Header
}

func NewHTTPClient(timeout time.Duration) *HTTPClient {
	return newHTTPClient(timeout, true)
}

// NewHTTPClientInsecureTLS creates an HTTP client that skips TLS certificate
// and hostname verification. Use only when the scan must inspect a site with
// an invalid or mismatched certificate.
func NewHTTPClientInsecureTLS(timeout time.Duration) *HTTPClient {
	return newHTTPClient(timeout, false)
}

func newHTTPClient(timeout time.Duration, verifyTLS bool) *HTTPClient {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	defaultHeaders := make(http.Header)
	defaultHeaders.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:157.0)")
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if !verifyTLS {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // #nosec G402: explicitly opt-in for technology fingerprinting
	}

	return &HTTPClient{
		defaultHeaders: defaultHeaders,
		client: &http.Client{
			Timeout:   timeout,
			Transport: transport,
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

// SetHeader configures a header to add to requests made by this wrapper.
// Configure headers before sending requests.
func (c *HTTPClient) SetHeader(name, value string) {
	c.defaultHeaders.Set(name, value)
}

// Do sends a prepared request and returns the response for any HTTP status.
// The caller is responsible for closing response.Body.
func (c *HTTPClient) Do(request *http.Request) (*http.Response, error) {
	request = request.Clone(request.Context())
	if request.Header == nil {
		request.Header = make(http.Header)
	}
	for name, values := range c.defaultHeaders {
		if request.Header.Values(name) == nil {
			for _, value := range values {
				request.Header.Add(name, value)
			}
		}
	}
	return c.client.Do(request)
}

// Request creates and sends an HTTP request with an optional byte body.
func (c *HTTPClient) Request(method, url string, data ...[]byte) (*http.Response, error) {
	var body io.Reader
	if len(data) > 0 {
		body = bytes.NewReader(data[0])
	}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	return c.Do(request)
}

// Post sends a POST request with the supplied content type and body.
// The caller is responsible for closing response.Body.
func (c *HTTPClient) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
	return c.Do(request)
}

func (c *HTTPClient) Get(url string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	return c.Do(request)
}

// Request sends a request using the default HTTP client.
func Request(method, url string, data ...[]byte) (*http.Response, error) {
	var body io.Reader
	if len(data) > 0 {
		body = bytes.NewReader(data[0])
	}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(request)
}

// Get fetches a URL using the default HTTP client and returns its response body.
