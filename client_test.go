package imazeo_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	imazeo "github.com/Imazeohq/imazeo-go"
)

func jsonOK(t *testing.T, w http.ResponseWriter, status int, data any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("x-request-id", "req_test")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    data,
		"meta":    map[string]any{"requestId": "req_test"},
	})
}

func jsonErr(t *testing.T, w http.ResponseWriter, status int, code, message string) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("x-request-id", "req_test")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": false,
		"error": map[string]any{
			"code":    code,
			"message": message,
		},
		"meta": map[string]any{"requestId": "req_test"},
	})
}

func TestNormalizeBaseURL(t *testing.T) {
	got, err := imazeo.NormalizeBaseURL("https://api.example.com/")
	if err != nil {
		t.Fatal(err)
	}
	if got != "https://api.example.com" {
		t.Fatalf("got %q", got)
	}
	got, err = imazeo.NormalizeBaseURL("https://api.example.com/v1")
	if err != nil {
		t.Fatal(err)
	}
	if got != "https://api.example.com" {
		t.Fatalf("got %q", got)
	}
	if _, err := imazeo.NormalizeBaseURL("ftp://x"); err == nil {
		t.Fatal("expected error for ftp")
	}
}

func TestRejectUnsupportedAPIVersion(t *testing.T) {
	_, err := imazeo.New(imazeo.WithAPIKey("key"), imazeo.WithAPIVersion("v9"))
	if err == nil || !strings.Contains(err.Error(), "unsupported apiVersion") {
		t.Fatalf("expected unsupported apiVersion error, got %v", err)
	}
}

func TestSendsAPIKeyAndUserAgent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/api-key/context" {
			t.Fatalf("path %s", r.URL.Path)
		}
		if r.Header.Get("X-API-Key") != "secret-key" {
			t.Fatalf("X-API-Key %q", r.Header.Get("X-API-Key"))
		}
		wantUA := "Imazeo-Go/" + imazeo.SDKVersion
		if r.Header.Get("User-Agent") != wantUA {
			t.Fatalf("User-Agent %q want %q", r.Header.Get("User-Agent"), wantUA)
		}
		if r.Header.Get("Authorization") != "" {
			t.Fatal("Authorization must not be set")
		}
		jsonOK(t, w, 200, map[string]any{
			"projectId": "proj_1",
			"keyId":     "key_1",
		})
	}))
	defer srv.Close()

	client, err := imazeo.New(
		imazeo.WithAPIKey("secret-key"),
		imazeo.WithBaseURL(srv.URL),
	)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := client.Context.Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if ctx.ProjectID != "proj_1" {
		t.Fatalf("projectId %q", ctx.ProjectID)
	}
}

func TestResolvesProjectIDFromContext(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/api-key/context"):
			jsonOK(t, w, 200, map[string]any{
				"projectId": "proj_auto",
				"keyId":     "k",
			})
		case strings.Contains(r.URL.Path, "/projects/proj_auto/images"):
			jsonOK(t, w, 200, []any{})
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer srv.Close()

	client, err := imazeo.New(imazeo.WithAPIKey("k"), imazeo.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	page, err := client.Images.List(context.Background(), imazeo.PageParams{Page: 1, Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if page.RequestID != "req_test" {
		t.Fatalf("requestId %q", page.RequestID)
	}
	if len(page.Items) != 0 {
		t.Fatalf("items %#v", page.Items)
	}
}

func TestUploadXUploadLength(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Upload-Length") != "4" {
			t.Fatalf("X-Upload-Length %q", r.Header.Get("X-Upload-Length"))
		}
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "multipart/form-data") {
			t.Fatalf("Content-Type %q", ct)
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatal(err)
		}
		f, _, err := r.FormFile("file")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		b, _ := io.ReadAll(f)
		if string(b) != "\x01\x02\x03\x04" {
			t.Fatalf("file bytes %v", b)
		}
		jsonOK(t, w, 201, map[string]any{
			"_id":          "img_1",
			"projectId":    "p1",
			"sourceType":   "upload",
			"originalName": "a.png",
			"bytes":        4,
			"sha256":       "abc",
			"metadata": map[string]any{
				"format": "png", "mimeType": "image/png",
				"width": 1, "height": 1, "frames": 1, "pixelCount": 1,
				"hasAlpha": false, "orientation": nil,
			},
			"status":    "ready",
			"createdAt": "2026-01-01T00:00:00.000Z",
			"updatedAt": "2026-01-01T00:00:00.000Z",
		})
	}))
	defer srv.Close()

	client, err := imazeo.New(
		imazeo.WithAPIKey("k"),
		imazeo.WithProjectID("p1"),
		imazeo.WithBaseURL(srv.URL),
	)
	if err != nil {
		t.Fatal(err)
	}
	image, err := client.Images.Upload(context.Background(), []byte{1, 2, 3, 4}, &imazeo.UploadOptions{
		Filename:    "a.png",
		ContentType: "image/png",
	})
	if err != nil {
		t.Fatal(err)
	}
	if image.ID != "img_1" {
		t.Fatalf("id %q", image.ID)
	}
}

func TestMapsAuthenticationAndRateLimitAndValidation(t *testing.T) {
	authSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonErr(t, w, 401, "API_KEY_INVALID", "bad key")
	}))
	defer authSrv.Close()
	c1, _ := imazeo.New(imazeo.WithAPIKey("k"), imazeo.WithProjectID("p"), imazeo.WithBaseURL(authSrv.URL))
	_, err := c1.Images.Get(context.Background(), "x")
	var auth *imazeo.AuthenticationError
	if !errors.As(err, &auth) {
		t.Fatalf("expected AuthenticationError, got %T", err)
	}

	rateSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonErr(t, w, 429, "RATE_LIMIT_EXCEEDED", "slow down")
	}))
	defer rateSrv.Close()
	c2, _ := imazeo.New(imazeo.WithAPIKey("k"), imazeo.WithProjectID("p"), imazeo.WithBaseURL(rateSrv.URL))
	_, err = c2.Images.Get(context.Background(), "x")
	var rate *imazeo.RateLimitError
	if !errors.As(err, &rate) {
		t.Fatalf("expected RateLimitError, got %T", err)
	}

	valSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonErr(t, w, 400, "TOOL_INVALID", "bad")
	}))
	defer valSrv.Close()
	c3, _ := imazeo.New(imazeo.WithAPIKey("k"), imazeo.WithProjectID("p"), imazeo.WithBaseURL(valSrv.URL))
	_, err = c3.QR.Generate(context.Background(), imazeo.QRGenerateInput{Content: "hi"})
	var val *imazeo.ValidationError
	if !errors.As(err, &val) {
		t.Fatalf("expected ValidationError, got %T", err)
	}
}

func TestQRBinaryResult(t *testing.T) {
	png := []byte{137, 80, 78, 71}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("x-request-id", "req_bin")
		w.WriteHeader(200)
		_, _ = w.Write(png)
	}))
	defer srv.Close()

	client, err := imazeo.New(
		imazeo.WithAPIKey("k"),
		imazeo.WithProjectID("p"),
		imazeo.WithBaseURL(srv.URL),
	)
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.QR.Generate(context.Background(), imazeo.QRGenerateInput{Content: "https://imazeo.com"})
	if err != nil {
		t.Fatal(err)
	}
	if result.ContentType != "image/png" {
		t.Fatalf("contentType %q", result.ContentType)
	}
	if string(result.Data) != string(png) {
		t.Fatalf("data mismatch")
	}
	if result.RequestID != "req_bin" {
		t.Fatalf("requestId %q", result.RequestID)
	}
}

func TestImageListQFilter(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("q") != "photo" {
			t.Fatalf("q=%q", r.URL.Query().Get("q"))
		}
		jsonOK(t, w, 200, []any{})
	}))
	defer srv.Close()

	client, err := imazeo.New(
		imazeo.WithAPIKey("k"),
		imazeo.WithProjectID("p"),
		imazeo.WithBaseURL(srv.URL),
	)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Images.List(context.Background(), imazeo.PageParams{Q: "photo", Page: 1}); err != nil {
		t.Fatal(err)
	}
}

func TestTransformDownload304(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("If-None-Match") != `"abc"` {
			t.Fatalf("If-None-Match %q", r.Header.Get("If-None-Match"))
		}
		w.Header().Set("ETag", `"abc"`)
		w.Header().Set("x-request-id", "req_304")
		w.WriteHeader(http.StatusNotModified)
	}))
	defer srv.Close()

	client, err := imazeo.New(
		imazeo.WithAPIKey("k"),
		imazeo.WithProjectID("p"),
		imazeo.WithBaseURL(srv.URL),
	)
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.Transforms.Download(context.Background(), "tr_1", &imazeo.DownloadOptions{
		IfNoneMatch: `"abc"`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.NotModified {
		t.Fatal("expected NotModified")
	}
	if len(result.Data) != 0 {
		t.Fatalf("data len %d", len(result.Data))
	}
	if result.RequestID != "req_304" {
		t.Fatalf("requestId %q", result.RequestID)
	}
	if result.ETag != `"abc"` {
		t.Fatalf("etag %q", result.ETag)
	}
}

func TestTransformCreateBinary(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method %s", r.Method)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["imageId"] != "img" {
			t.Fatalf("body %#v", body)
		}
		w.Header().Set("Content-Type", "image/webp")
		w.Header().Set("x-transformation-id", "tr_1")
		w.Header().Set("x-cache", "MISS")
		w.Header().Set("x-request-id", "req_t")
		w.WriteHeader(200)
		_, _ = w.Write([]byte{1, 2})
	}))
	defer srv.Close()

	client, err := imazeo.New(
		imazeo.WithAPIKey("k"),
		imazeo.WithProjectID("p"),
		imazeo.WithBaseURL(srv.URL),
	)
	if err != nil {
		t.Fatal(err)
	}
	width := 100
	out, err := client.Transforms.Create(context.Background(), imazeo.TransformRequest{
		ImageID: "img",
		Transform: &imazeo.TransformSpec{
			Width:  &width,
			Format: "webp",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.TransformationID != "tr_1" {
		t.Fatalf("transformationId %q", out.TransformationID)
	}
	if out.Cache != "MISS" {
		t.Fatalf("cache %q", out.Cache)
	}
}

func TestVersionAndCompatibility(t *testing.T) {
	if imazeo.Version() != imazeo.SDKVersion {
		t.Fatalf("version %s", imazeo.Version())
	}
	def, supported := imazeo.APICompatibility()
	if def != "v1" || len(supported) != 1 || supported[0] != "v1" {
		t.Fatalf("%s %#v", def, supported)
	}
}

func TestRequiresAPIKey(t *testing.T) {
	_, err := imazeo.New()
	if err == nil || !strings.Contains(err.Error(), "apiKey is required") {
		t.Fatalf("got %v", err)
	}
}
