package imazeo_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	imazeo "github.com/Imazeohq/imazeo-go"
)

func TestConversionQuoteAndUpload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/projects/project/conversions/xlsx" || r.URL.Query().Get("showGridLines") != "false" {
			t.Errorf("unexpected URL: %s", r.URL)
		}
		if r.Header.Get("Idempotency-Key") != "operation-key" || r.Header.Get("X-Quote-ID") != "quote-id" || r.Header.Get("X-Upload-Length") != "3" {
			t.Errorf("missing processing headers")
		}
		if err := r.ParseMultipartForm(1024); err != nil {
			t.Error(err)
		}
		if len(r.MultipartForm.File["file"]) != 1 {
			t.Error("expected one file")
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "data": map[string]any{"id": "job"}, "meta": map[string]any{"requestId": "request"}})
	}))
	defer server.Close()
	client, err := imazeo.New(imazeo.WithAPIKey("test"), imazeo.WithProjectID("project"), imazeo.WithBaseURL(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	grid := false
	job, err := client.Conversions.XLSX(context.Background(), []byte("abc"), imazeo.ProcessingOptions{IdempotencyKey: "operation-key", QuoteID: "quote-id"}, &imazeo.XlsxOptions{ShowGridLines: &grid})
	if err != nil {
		t.Fatal(err)
	}
	if job.ID != "job" {
		t.Fatalf("unexpected job: %#v", job)
	}
}

func TestDynamicRevisionHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		if body["expectedRevision"] != float64(0) || body["analyticsEnabled"] != false || r.Header.Get("Idempotency-Key") != "edit-key" {
			t.Errorf("unexpected edit: %#v", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"revision":1},"meta":{"requestId":"request"}}`))
	}))
	defer server.Close()
	client, err := imazeo.New(imazeo.WithAPIKey("test"), imazeo.WithProjectID("project"), imazeo.WithBaseURL(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	analytics := false
	_, err = client.DynamicCodes.Update(context.Background(), "code", imazeo.DynamicEditInput{ExpectedRevision: 0, AnalyticsEnabled: &analytics}, "edit-key")
	if err != nil {
		t.Fatal(err)
	}
}
