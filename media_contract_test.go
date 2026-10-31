package imazeo_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	imazeo "github.com/Imazeohq/imazeo-go"
)

func TestBarcodeContract(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
			return
		}
		if len(body) != 4 || body["type"] != "ean8" || body["content"] != "0123456" || body["format"] != "svg" || body["includeText"] != false {
			t.Errorf("body: %#v", body)
		}
		w.Header().Set("Content-Type", "image/svg+xml")
		_, _ = w.Write([]byte("<svg/>"))
	}))
	defer server.Close()
	client, err := imazeo.New(imazeo.WithAPIKey("test"), imazeo.WithProjectID("project"), imazeo.WithBaseURL(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	showText := false
	_, err = client.Barcodes.Generate(context.Background(), imazeo.BarcodeGenerateInput{
		Type: "ean8", Content: "0123456", Format: "svg", IncludeText: &showText,
	})
	if err != nil {
		t.Fatal(err)
	}
}
