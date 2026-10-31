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
	if len(os.Args) < 2 {
		fatal("Usage: go run ./examples/upload_transform <image>")
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

	bytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		fatal(err)
	}

	image, err := client.Images.Upload(context.Background(), bytes, &imazeo.UploadOptions{
		Filename:    "input.bin",
		ContentType: "application/octet-stream",
	})
	if err != nil {
		fatal(err)
	}
	fmt.Println("uploaded", image.ID, image.Metadata)

	width, quality := 800, 80
	transformed, err := client.Transforms.Create(context.Background(), imazeo.TransformRequest{
		ImageID: image.ID,
		Mode:    imazeo.TransformModeEphemeral,
		Transform: &imazeo.TransformSpec{
			Width:   &width,
			Format:  "webp",
			Quality: &quality,
		},
	})
	if err != nil {
		fatal(err)
	}
	if err := os.WriteFile("out.webp", transformed.Data, 0o644); err != nil {
		fatal(err)
	}
	fmt.Println("wrote out.webp", transformed.TransformationID)
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
