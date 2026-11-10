# imazeo-go

Official **server-side** Go SDK for Imazeo.

|                |                                          |
| -------------- | ---------------------------------------- |
| Module         | `github.com/Imazeohq/imazeo-go`          |
| SDK version    | `1.0.0`                                  |
| Default API    | `v1`                                     |
| Supported APIs | `v1`                                     |
| Go             | **1.22+**                                |
| Registry       | [pkg.go.dev](https://pkg.go.dev/github.com/Imazeohq/imazeo-go) |
| Source         | [GitHub](https://github.com/Imazeohq/imazeo-go) |

> **Security:** Secret API keys (`X-API-Key`) are for trusted servers only.
> Do **not** ship this client with a secret key in browsers or end-user apps.

Integration contract and learning path: [GUIDE.md](GUIDE.md).

## Install

```bash
go get github.com/Imazeohq/imazeo-go@v1.0.0
```

Requires Go 1.22+. Create a project API key in the
[Imazeo Customer portal](https://app.imazeo.com).
Default API origin is `https://api.imazeo.com` (override with `WithBaseURL` when needed).

## Learning path

Read the self-contained [integration guide](GUIDE.md) for project/key setup,
language examples, PDF response types, pagination, output retention, troubleshooting
and production practices.

The default transport refuses redirects (`CheckRedirect` → `http.ErrUseLastResponse`)
so credentials stay on the configured API origin.

## Quick start

```go
package main

import (
	"context"
	"os"

	"github.com/Imazeohq/imazeo-go"
)

func main() {
	client, err := imazeo.New(
		imazeo.WithAPIKey(os.Getenv("IMAZEO_API_KEY")),
		imazeo.WithBaseURL(os.Getenv("IMAZEO_API_ORIGIN")),
		// imazeo.WithProjectID("…"),
	)
	if err != nil {
		panic(err)
	}

	qr, err := client.QR.Generate(context.Background(), imazeo.QRGenerateInput{
		Content: "https://imazeo.com",
		Format:  "png",
		Size:    256,
	})
	if err != nil {
		panic(err)
	}
	_ = os.WriteFile("qr.png", qr.Data, 0o644)
}
```

## Authentication

```go
imazeo.New(imazeo.WithAPIKey(os.Getenv("IMAZEO_API_KEY")))
```

Sends `X-API-Key` only. Never send `Authorization` together with an API key
(backend rejects that combination).

## Resources

| Resource            | Methods                                                                       |
| ------------------- | ----------------------------------------------------------------------------- |
| `client.Context`    | `Get` → project/key IDs                                                       |
| `client.Images`     | `Upload`, `List` (`Page`/`Limit`/`Q`), `Get`, `Delete`, `Download`            |
| `client.Transforms` | `Create`, `List`, `Get`, `Download` (`IfNoneMatch` → `NotModified`), `Delete` |
| `client.QR`         | `Generate`                                                                    |
| `client.Barcodes`   | `Generate`                                                                    |
| `client.PDF`        | `Info`, `Render`, `Merge`, `Extract`                                          |

### Upload + transform

```go
image, err := client.Images.Upload(ctx, bytes, &imazeo.UploadOptions{
	Filename:    "photo.jpg",
	ContentType: "image/jpeg",
})
width, quality := 1200, 85
out, err := client.Transforms.Create(ctx, imazeo.TransformRequest{
	ImageID: image.ID,
	Mode:    imazeo.TransformModeEphemeral,
	Transform: &imazeo.TransformSpec{
		Width:   &width,
		Format:  "webp",
		Quality: &quality,
	},
})
// out.Data → []byte
```

### PDF

```go
info, err := client.PDF.Info(ctx, pdfBytes, nil)
page, err := client.PDF.Render(ctx, pdfBytes, &imazeo.PDFRenderOptions{Pages: "1", Format: "png"})
merged, err := client.PDF.Merge(ctx, [][]byte{pdfA, pdfB})
extracted, err := client.PDF.Extract(ctx, pdfBytes, imazeo.PDFExtractOptions{Pages: "1,2"})
```

## Errors

```go
import (
	"errors"
	"github.com/Imazeohq/imazeo-go"
)

_, err := client.Images.Get(ctx, id)
var rate *imazeo.RateLimitError
if errors.As(err, &rate) {
	fmt.Println(rate.Code, rate.RequestID)
} else {
	var ie *imazeo.Error
	if errors.As(err, &ie) {
		fmt.Println(ie.Status, ie.Code, ie.Message)
	}
}
```

Typed errors: `AuthenticationError`, `AuthorizationError`, `ValidationError`,
`NotFoundError`, `ConflictError`, `RateLimitError`, `APIError`, `NetworkError`.

## Configuration

| Option             | Default                  | Notes                                  |
| ------------------ | ------------------------ | -------------------------------------- |
| `WithAPIKey`       | required                 | Secret key                             |
| `WithProjectID`    | from context             | Path project                           |
| `WithBaseURL`      | `https://api.imazeo.com` | Origin only; `/v1` stripped if present |
| `WithAPIVersion`   | `v1`                     | Allowlisted                            |
| `WithTimeout`      | `60s`                    | Image upload / transform / download    |
| `WithLightTimeout` | `30s`                    | Context, QR, barcode                   |
| `WithPDFTimeout`   | `120s`                   | PDF tools                              |
| `WithAppName`      | —                        | Appended to User-Agent                 |
| `WithHTTPClient`   | default client           | Tests / custom agents; redirects still refused |

User-Agent: `Imazeo-Go/1.0.0`.

## API version vs package version

- **API** path version: `v1` (via `WithAPIVersion`)
- **SDK** SemVer: `1.0.0` (`imazeo.SDKVersion`)

`imazeo.APICompatibility()` reports supported API versions.

## Not included

Portal JWT APIs (projects CRUD, API-key management, usage, plans, tickets, auth),
admin APIs, tool artifacts, project outputs, and browser secret-key usage.

## Development

```bash
go test ./...
go vet ./...
```

## License

MIT

## Barcode and PDF notes

Barcode generation requires `type` and string `content`. Optional `format` is
`png` or `svg` (default `png`); `includeText` defaults to `true`.
PDF merge uploads files in order as repeated multipart `file` parts.
Omit optional fields to use API defaults. Integer fields such as barcode scale
and image hue must be integers.

```go
showText := false
result, err := client.Barcodes.Generate(ctx, imazeo.BarcodeGenerateInput{Type: "code128", Content: "IMAZEO-2026", Format: "svg", IncludeText: &showText})
```
