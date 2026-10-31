# imazeo-go — agent notes

Independent Go module. Backend OpenAPI/`/v1` is the contract source of truth.
Follow workspace [docs/SDK_BEHAVIOR.md](../docs/SDK_BEHAVIOR.md).
Mirror behavioral surface of `../sdk-typescript/`; keep Go idioms (functional
options, `context.Context`, `[]byte` uploads, `httptest` tests).
Do not add admin or JWT portal APIs to this client.
Publish only after the owner creates `Imazeohq/imazeo-go`, then follow
[RELEASE.md](RELEASE.md). Module path must stay `github.com/Imazeohq/imazeo-go`.
