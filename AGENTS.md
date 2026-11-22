# imazeo-go — agent notes

Independent public Go module (`github.com/Imazeohq/imazeo-go`).

- Contract: published API `v1` (see [GUIDE.md](GUIDE.md))
- Scope: project API-key context, media, credit quotes, operations, conversions, dynamic QR, artifacts and outputs
- Keep Go idioms (`context.Context`, functional options, `[]byte` uploads)
- Do not add admin or customer JWT portal APIs
- Do not publish without maintainer approval — see [RELEASE.md](RELEASE.md)
- Module path must stay `github.com/Imazeohq/imazeo-go`
- Release: `$env:RELEASE_MODE="true"; npm run release` (syncs `version.go`); Actions indexes the module
