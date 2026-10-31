package imazeo

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Client is the official Imazeo server-side client.
//
// Uses secret X-API-Key authentication. Do not embed API keys in browsers.
type Client struct {
	Context    *ContextResource
	Images     *ImagesResource
	Transforms *TransformsResource
	QR         *QRResource
	Barcodes   *BarcodesResource
	PDF        *PDFResource

	http               *httpTransport
	configuredProjectID string
	resolvedProjectID   string
	projectMu           sync.Mutex
}

type clientOptions struct {
	apiKey     string
	projectID  string
	baseURL    string
	apiVersion string
	timeout    time.Duration
	lightTO    time.Duration
	pdfTO      time.Duration
	appName    string
	httpClient *http.Client
}

// Option configures a Client.
type Option func(*clientOptions) error

// WithAPIKey sets the secret API key (required). Sent as X-API-Key only.
func WithAPIKey(apiKey string) Option {
	return func(o *clientOptions) error {
		o.apiKey = strings.TrimSpace(apiKey)
		return nil
	}
}

// WithProjectID sets the path project id. If omitted, resolved via GET /v1/api-key/context.
func WithProjectID(projectID string) Option {
	return func(o *clientOptions) error {
		o.projectID = strings.TrimSpace(projectID)
		return nil
	}
}

// WithBaseURL sets the API origin (no /v1 suffix). Default: https://api.imazeo.com
func WithBaseURL(baseURL string) Option {
	return func(o *clientOptions) error {
		o.baseURL = strings.TrimSpace(baseURL)
		return nil
	}
}

// WithAPIVersion sets the allowlisted API version. Default: v1
func WithAPIVersion(apiVersion string) Option {
	return func(o *clientOptions) error {
		o.apiVersion = strings.TrimSpace(apiVersion)
		return nil
	}
}

// WithTimeout sets the default request timeout (image upload/transform/download). Default: 60s
func WithTimeout(d time.Duration) Option {
	return func(o *clientOptions) error {
		o.timeout = d
		return nil
	}
}

// WithLightTimeout sets the timeout for JSON / QR / barcode / context. Default: 30s
func WithLightTimeout(d time.Duration) Option {
	return func(o *clientOptions) error {
		o.lightTO = d
		return nil
	}
}

// WithPDFTimeout sets the timeout for PDF operations. Default: 120s
func WithPDFTimeout(d time.Duration) Option {
	return func(o *clientOptions) error {
		o.pdfTO = d
		return nil
	}
}

// WithAppName appends an application name to the User-Agent.
func WithAppName(appName string) Option {
	return func(o *clientOptions) error {
		o.appName = strings.TrimSpace(appName)
		return nil
	}
}

// WithHTTPClient overrides the HTTP client (tests / custom transports).
func WithHTTPClient(c *http.Client) Option {
	return func(o *clientOptions) error {
		o.httpClient = c
		return nil
	}
}

// New constructs an Imazeo client. WithAPIKey is required.
func New(opts ...Option) (*Client, error) {
	o := &clientOptions{
		baseURL:    DefaultBaseURL,
		apiVersion: DefaultAPIVersion,
		timeout:    DefaultTimeout,
		lightTO:    LightTimeout,
		pdfTO:      PDFTimeout,
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(o); err != nil {
			return nil, err
		}
	}
	apiKey := strings.TrimSpace(o.apiKey)
	if apiKey == "" {
		return nil, fmt.Errorf("apiKey is required")
	}
	if len(apiKey) > 128 {
		return nil, fmt.Errorf("apiKey is too long")
	}
	apiVersion, err := assertAPIVersion(o.apiVersion)
	if err != nil {
		return nil, err
	}
	baseURL := o.baseURL
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	baseURL, err = NormalizeBaseURL(baseURL)
	if err != nil {
		return nil, err
	}
	if o.timeout < time.Millisecond {
		return nil, fmt.Errorf("timeout must be a positive duration")
	}
	if o.lightTO < time.Millisecond {
		return nil, fmt.Errorf("light timeout must be a positive duration")
	}
	if o.pdfTO < time.Millisecond {
		return nil, fmt.Errorf("pdf timeout must be a positive duration")
	}

	ua := fmt.Sprintf("Imazeo-Go/%s", SDKVersion)
	if o.appName != "" {
		ua = ua + " " + o.appName
	}
	if len(ua) > 256 {
		ua = ua[:256]
	}

	tr := &httpTransport{config: transportConfig{
		apiKey:     apiKey,
		baseURL:    baseURL,
		apiVersion: apiVersion,
		timeout:    o.timeout,
		lightTO:    o.lightTO,
		pdfTO:      o.pdfTO,
		userAgent:  ua,
		httpClient: o.httpClient,
	}}

	c := &Client{
		http:                tr,
		configuredProjectID: o.projectID,
	}
	c.Context = &ContextResource{http: tr}
	c.Images = &ImagesResource{http: tr, projectID: c.resolveProjectID}
	c.Transforms = &TransformsResource{http: tr, projectID: c.resolveProjectID}
	c.QR = &QRResource{http: tr, projectID: c.resolveProjectID}
	c.Barcodes = &BarcodesResource{http: tr, projectID: c.resolveProjectID}
	c.PDF = &PDFResource{http: tr, projectID: c.resolveProjectID}
	return c, nil
}

// Version returns the SDK package SemVer.
func Version() string { return SDKVersion }

// APICompatibility reports the default and supported API versions.
func APICompatibility() (defaultAPI string, supported []string) {
	out := make([]string, len(SupportedAPIVersions))
	copy(out, SupportedAPIVersions)
	return DefaultAPIVersion, out
}

func (c *Client) resolveProjectID(ctx context.Context) (string, error) {
	if c.configuredProjectID != "" {
		return c.configuredProjectID, nil
	}
	c.projectMu.Lock()
	defer c.projectMu.Unlock()
	if c.resolvedProjectID != "" {
		return c.resolvedProjectID, nil
	}
	ctxInfo, err := c.Context.Get(ctx)
	if err != nil {
		return "", err
	}
	if ctxInfo.ProjectID == "" {
		return "", fmt.Errorf("api-key context did not include projectId")
	}
	c.resolvedProjectID = ctxInfo.ProjectID
	return c.resolvedProjectID, nil
}
