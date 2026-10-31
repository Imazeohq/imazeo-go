# SDK-4 plan — Go

**Status: implemented (1.0.0)**

1. Module path: `github.com/Imazeohq/imazeo-go` (`go 1.22`) — done
2. Functional options (`WithAPIKey`, `WithBaseURL`, …) — done
3. Typed structs; `[]byte` uploads with `X-Upload-Length` — done
4. Table-driven / httptest tests in `*_test.go` — done
5. Resources: Context, Images, Transforms, QR, Barcodes, PDF — done
6. Image list `q`; transform download 304 → `NotModified` — done
7. Default transport refuses redirects — done

Scope: API-key resources only. No admin/JWT/tool-artifacts/outputs.

Reference: `../sdk-typescript/` and `../docs/SDK_BEHAVIOR.md`.
Publish by tagging `v1.0.0` on `Imazeohq/imazeo-go` (Go proxy + pkg.go.dev).
