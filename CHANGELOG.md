# Changelog

## 1.0.0

### Added

- Initial official Go SDK (`github.com/Imazeohq/imazeo-go`)
- Resources: `Context`, `Images`, `Transforms`, `QR`, `Barcodes`, `PDF`
- Typed errors mapped from Imazeo API envelopes
- API version allowlist (`v1`) separate from SDK SemVer
- Multipart uploads with required `X-Upload-Length`
- Binary response helpers for images/PDF/QR/barcode
- Image list `q` filename-prefix filter
- Transform download `If-None-Match` → HTTP 304 as `BinaryResult.NotModified`
- Default transport refuses redirects so `X-API-Key` stays on the configured origin
- Reviewed API-key-only OpenAPI snapshot bundled as `openapi/backend.json`
- Native CI verification and independent SemVer/tag checks

### Media contract

- Barcode: required `type` and `content`; optional `format` (`png`/`svg`) and `includeText`
- PDF merge uses repeated multipart `file` parts in order
- Integer model types for barcode scale and image hue where applicable

### Out of scope for this SDK

- Tool artifacts CRUD/download
- Project outputs list

### API compatibility

- Default API: `v1`
- Supported API versions: `v1`
- Default base URL: `https://api.imazeo.com`
