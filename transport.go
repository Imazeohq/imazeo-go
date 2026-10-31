package imazeo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type transportConfig struct {
	apiKey     string
	baseURL    string
	apiVersion string
	timeout    time.Duration
	lightTO    time.Duration
	pdfTO      time.Duration
	userAgent  string
	httpClient *http.Client
}

type jsonSuccess[T any] struct {
	Data T
	Meta responseMeta
}

type responseMeta struct {
	RequestID string
	Page      *int
	Limit     *int
	Total     *int
	Pages     *int
}

type envelope[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
	Meta    *struct {
		RequestID string `json:"requestId"`
		Page      *int   `json:"page"`
		Limit     *int   `json:"limit"`
		Total     *int   `json:"total"`
		Pages     *int   `json:"pages"`
	} `json:"meta"`
	Error *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Details any    `json:"details"`
	} `json:"error"`
}

func assertAPIVersion(value string) (string, error) {
	for _, v := range SupportedAPIVersions {
		if v == value {
			return value, nil
		}
	}
	return "", fmt.Errorf("unsupported apiVersion %q; supported: %s", value, strings.Join(SupportedAPIVersions, ", "))
}

// NormalizeBaseURL validates and normalizes an API origin (strips trailing slash and accidental /v1).
func NormalizeBaseURL(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("invalid baseUrl: %s", raw)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", fmt.Errorf("baseUrl must use http or https")
	}
	if u.User != nil {
		return "", fmt.Errorf("baseUrl must not include credentials")
	}
	path := strings.TrimRight(u.Path, "/")
	if path == "/v1" {
		path = ""
	}
	u.Path = path
	u.RawQuery = ""
	u.Fragment = ""
	if path == "" || path == "/" {
		return strings.TrimRight(u.Scheme+"://"+u.Host, "/"), nil
	}
	return strings.TrimRight(u.Scheme+"://"+u.Host+path, "/"), nil
}

type httpTransport struct {
	config transportConfig
}

func (t *httpTransport) path(segments ...string) string {
	parts := make([]string, 0, len(segments)+2)
	parts = append(parts, t.config.baseURL, t.config.apiVersion)
	for _, s := range segments {
		parts = append(parts, url.PathEscape(s))
	}
	// PathEscape encodes spaces etc.; join without extra encoding of slashes in segments.
	return strings.Join(parts, "/")
}

func (t *httpTransport) withTimeout(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	if d <= 0 {
		d = t.config.timeout
	}
	return context.WithTimeout(ctx, d)
}

func (t *httpTransport) requestJSON(ctx context.Context, method, rawURL string, body any, timeout time.Duration) (*http.Response, []byte, error) {
	var reader io.Reader
	headers := map[string]string{
		"Accept":     "application/json",
		"X-API-Key":  t.config.apiKey,
		"User-Agent": t.config.userAgent,
	}
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, nil, err
		}
		reader = bytes.NewReader(b)
		headers["Content-Type"] = "application/json"
	}
	return t.send(ctx, method, rawURL, reader, headers, timeout)
}

func (t *httpTransport) requestBinary(ctx context.Context, method, rawURL string, body any, extraHeaders map[string]string, timeout time.Duration) (*BinaryResult, error) {
	var reader io.Reader
	headers := map[string]string{
		"Accept":     "*/*",
		"X-API-Key":  t.config.apiKey,
		"User-Agent": t.config.userAgent,
	}
	for k, v := range extraHeaders {
		headers[k] = v
	}
	if body != nil {
		switch v := body.(type) {
		case []byte:
			reader = bytes.NewReader(v)
		case io.Reader:
			reader = v
		default:
			b, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			reader = bytes.NewReader(b)
			if _, ok := headers["Content-Type"]; !ok {
				headers["Content-Type"] = "application/json"
			}
		}
	}
	resp, raw, err := t.send(ctx, method, rawURL, reader, headers, timeout)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified {
		return binaryFromResponse(resp, nil, true, headers["If-None-Match"]), nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, t.errorFromBody(resp.StatusCode, raw, resp.Header.Get("x-request-id"))
	}
	return binaryFromResponse(resp, raw, false, ""), nil
}

type multipartFile struct {
	FieldName   string
	Filename    string
	ContentType string
	Data        []byte
}

func (t *httpTransport) requestMultipartJSON(ctx context.Context, method, rawURL string, files []multipartFile, fields map[string]string, uploadLength *int64, timeout time.Duration) ([]byte, responseMeta, error) {
	body, contentType, totalLen, err := buildMultipart(files, fields)
	if err != nil {
		return nil, responseMeta{}, err
	}
	headers := map[string]string{
		"Accept":       "application/json",
		"X-API-Key":    t.config.apiKey,
		"User-Agent":   t.config.userAgent,
		"Content-Type": contentType,
	}
	if uploadLength != nil {
		headers["X-Upload-Length"] = strconv.FormatInt(*uploadLength, 10)
	} else if totalLen >= 0 {
		// For single-file uploads callers should pass exact file byte length.
	}
	resp, raw, err := t.send(ctx, method, rawURL, body, headers, timeout)
	if err != nil {
		return nil, responseMeta{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, responseMeta{}, t.errorFromBody(resp.StatusCode, raw, resp.Header.Get("x-request-id"))
	}
	data, meta, err := parseSuccessEnvelope(raw, resp.Header.Get("x-request-id"))
	if err != nil {
		return nil, responseMeta{}, err
	}
	return data, meta, nil
}

func (t *httpTransport) requestMultipartBinary(ctx context.Context, method, rawURL string, files []multipartFile, uploadLength *int64, timeout time.Duration) (*BinaryResult, json.RawMessage, responseMeta, error) {
	body, contentType, _, err := buildMultipart(files, nil)
	if err != nil {
		return nil, nil, responseMeta{}, err
	}
	headers := map[string]string{
		"Accept":       "*/*",
		"X-API-Key":    t.config.apiKey,
		"User-Agent":   t.config.userAgent,
		"Content-Type": contentType,
	}
	if uploadLength != nil {
		headers["X-Upload-Length"] = strconv.FormatInt(*uploadLength, 10)
	}
	resp, raw, err := t.send(ctx, method, rawURL, body, headers, timeout)
	if err != nil {
		return nil, nil, responseMeta{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, nil, responseMeta{}, t.errorFromBody(resp.StatusCode, raw, resp.Header.Get("x-request-id"))
	}
	ct := resp.Header.Get("Content-Type")
	if strings.Contains(ct, "application/json") {
		data, meta, err := parseSuccessEnvelope(raw, resp.Header.Get("x-request-id"))
		if err != nil {
			return nil, nil, responseMeta{}, err
		}
		return nil, data, meta, nil
	}
	return binaryFromResponse(resp, raw, false, ""), nil, responseMeta{}, nil
}

func doJSON[T any](t *httpTransport, ctx context.Context, method, rawURL string, body any, timeout time.Duration) (jsonSuccess[T], error) {
	var zero jsonSuccess[T]
	resp, raw, err := t.requestJSON(ctx, method, rawURL, body, timeout)
	if err != nil {
		return zero, err
	}
	defer resp.Body.Close()
	requestID := resp.Header.Get("x-request-id")
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return zero, t.errorFromBody(resp.StatusCode, raw, requestID)
	}
	var env envelope[T]
	if err := json.Unmarshal(raw, &env); err != nil {
		return zero, newNetworkError("response was not valid JSON", err, requestID)
	}
	if !env.Success {
		return zero, mapHTTPError(resp.StatusCode, "", "unexpected success envelope", requestID, nil)
	}
	meta := responseMeta{RequestID: requestID}
	if env.Meta != nil {
		if env.Meta.RequestID != "" {
			meta.RequestID = env.Meta.RequestID
		}
		meta.Page = env.Meta.Page
		meta.Limit = env.Meta.Limit
		meta.Total = env.Meta.Total
		meta.Pages = env.Meta.Pages
	}
	return jsonSuccess[T]{Data: env.Data, Meta: meta}, nil
}

func (t *httpTransport) send(ctx context.Context, method, rawURL string, body io.Reader, headers map[string]string, timeout time.Duration) (*http.Response, []byte, error) {
	ctx, cancel := t.withTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		return nil, nil, newNetworkError(err.Error(), err, "")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := t.config.httpClient
	if client == nil {
		client = http.DefaultClient
	}
	// Copy the client's settings without mutating a caller-owned client. API keys
	// must never follow a redirect, even to another endpoint on the same host.
	safeClient := *client
	safeClient.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := safeClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, nil, newNetworkError(fmt.Sprintf("request timed out after %s", timeout), err, "")
		}
		return nil, nil, newNetworkError(err.Error(), err, "")
	}
	raw, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, nil, newNetworkError("failed to read response body", err, resp.Header.Get("x-request-id"))
	}
	// Re-open body for callers that still want resp.Body (already closed/consumed).
	resp.Body = io.NopCloser(bytes.NewReader(raw))
	return resp, raw, nil
}

func (t *httpTransport) errorFromBody(status int, raw []byte, requestID string) error {
	var env envelope[json.RawMessage]
	if err := json.Unmarshal(raw, &env); err == nil && env.Error != nil {
		rid := requestID
		if env.Meta != nil && env.Meta.RequestID != "" {
			rid = env.Meta.RequestID
		}
		return mapHTTPError(status, env.Error.Code, env.Error.Message, rid, env.Error.Details)
	}
	msg := http.StatusText(status)
	if msg == "" {
		msg = "request failed"
	}
	return mapHTTPError(status, "", msg, requestID, nil)
}

func parseSuccessEnvelope(raw []byte, requestID string) (json.RawMessage, responseMeta, error) {
	var env envelope[json.RawMessage]
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, responseMeta{}, newNetworkError("response was not valid JSON", err, requestID)
	}
	if !env.Success {
		return nil, responseMeta{}, mapHTTPError(http.StatusOK, "", "unexpected success envelope", requestID, nil)
	}
	meta := responseMeta{RequestID: requestID}
	if env.Meta != nil {
		if env.Meta.RequestID != "" {
			meta.RequestID = env.Meta.RequestID
		}
		meta.Page = env.Meta.Page
		meta.Limit = env.Meta.Limit
		meta.Total = env.Meta.Total
		meta.Pages = env.Meta.Pages
	}
	return env.Data, meta, nil
}

func binaryFromResponse(resp *http.Response, data []byte, notModified bool, fallbackETag string) *BinaryResult {
	headers := map[string]string{}
	for k, vals := range resp.Header {
		if len(vals) > 0 {
			headers[strings.ToLower(k)] = vals[0]
		}
	}
	ct := resp.Header.Get("Content-Type")
	if ct == "" {
		ct = "application/octet-stream"
	}
	etag := resp.Header.Get("ETag")
	if etag == "" {
		etag = fallbackETag
	}
	return &BinaryResult{
		Data:             data,
		ContentType:      ct,
		RequestID:        resp.Header.Get("x-request-id"),
		Headers:          headers,
		ETag:             etag,
		Width:            numHeader(resp.Header, "x-image-width"),
		Height:           numHeader(resp.Header, "x-image-height"),
		Format:           resp.Header.Get("x-image-format"),
		TransformationID: resp.Header.Get("x-transformation-id"),
		Cache:            resp.Header.Get("x-cache"),
		PDFPage:          numHeader(resp.Header, "x-pdf-page"),
		NotModified:      notModified,
	}
}

func numHeader(h http.Header, name string) *int {
	raw := h.Get(name)
	if raw == "" {
		return nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return nil
	}
	return &n
}

func buildMultipart(files []multipartFile, fields map[string]string) (*bytes.Buffer, string, int64, error) {
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	var fileBytes int64
	for _, f := range files {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes(f.FieldName), escapeQuotes(f.Filename)))
		if f.ContentType != "" {
			h.Set("Content-Type", f.ContentType)
		}
		part, err := w.CreatePart(h)
		if err != nil {
			return nil, "", 0, err
		}
		if _, err := part.Write(f.Data); err != nil {
			return nil, "", 0, err
		}
		fileBytes += int64(len(f.Data))
	}
	for k, v := range fields {
		if err := w.WriteField(k, v); err != nil {
			return nil, "", 0, err
		}
	}
	if err := w.Close(); err != nil {
		return nil, "", 0, err
	}
	return buf, w.FormDataContentType(), fileBytes, nil
}

func escapeQuotes(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}

func prepareBytes(data []byte, opts *UploadOptions, defaultName, defaultType string) multipartFile {
	filename := defaultName
	contentType := defaultType
	if opts != nil {
		if opts.Filename != "" {
			filename = opts.Filename
		}
		if opts.ContentType != "" {
			contentType = opts.ContentType
		}
	}
	if filename == "" {
		filename = "file.bin"
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return multipartFile{
		FieldName:   "file",
		Filename:    filename,
		ContentType: contentType,
		Data:        data,
	}
}

func intPtrOr(p *int, fallback int) int {
	if p != nil {
		return *p
	}
	return fallback
}
