// Minimal server-side QR example.
// Requires: IMAZEO_API_KEY; optional IMAZEO_BASE_URL, IMAZEO_PROJECT_ID
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Imazeohq/imazeo-go"
)

func main() {
	apiKey := os.Getenv("IMAZEO_API_KEY")
	if apiKey == "" {
		fatal("Set IMAZEO_API_KEY")
	}

	opts := []imazeo.Option{
		imazeo.WithAPIKey(apiKey),
		imazeo.WithBaseURL(envOr("IMAZEO_BASE_URL", "http://localhost:3001")),
	}
	if pid := os.Getenv("IMAZEO_PROJECT_ID"); pid != "" {
		opts = append(opts, imazeo.WithProjectID(pid))
	}

	client, err := imazeo.New(opts...)
	if err != nil {
		fatal(err)
	}

	ctxInfo, err := client.Context.Get(context.Background())
	if err != nil {
		fatal(err)
	}
	fmt.Println("project", ctxInfo.ProjectID, "key", ctxInfo.KeyID)

	qr, err := client.QR.Generate(context.Background(), imazeo.QRGenerateInput{
		Content: "https://imazeo.com",
		Size:    256,
		Format:  "png",
	})
	if err != nil {
		fatal(err)
	}
	if err := os.WriteFile("qr.png", qr.Data, 0o644); err != nil {
		fatal(err)
	}
	fmt.Println("wrote qr.png", qr.ContentType)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func fatal(v any) {
	fmt.Fprintln(os.Stderr, v)
	os.Exit(1)
}
