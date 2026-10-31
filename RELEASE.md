# SDK release procedure (`github.com/Imazeohq/imazeo-go`)

Independent repository: [Imazeohq/imazeo-go](https://github.com/Imazeohq/imazeo-go).  
Module / docs: [pkg.go.dev/github.com/Imazeohq/imazeo-go](https://pkg.go.dev/github.com/Imazeohq/imazeo-go).  
Supported API prefix: `v1`. SemVer lives in this repo (`version.go` → `SDKVersion`)
and is published from git tag `v1.0.0`.

Go modules are published by tagging a public GitHub repo — there is no separate
registry upload (unlike npm/PyPI). After the tag is public, pkg.go.dev indexes it
(often within minutes; you can force a fetch — see below).

Do not publish automatically from CI until a protected release workflow is added.
Never overwrite a published version; ship a patch instead.

## Pre-publish checks

```powershell
cd D:\projects\imazeo\sdk-go
python scripts\check-release.py
gofmt -l .
go vet ./...
go test -race ./...
go build ./...
```

## Manual publish (first release `1.0.0`)

1. Create empty public GitHub repo `Imazeohq/imazeo-go` (no README/license from the UI).
2. Push this tree to `main`, then tag and push:
   ```powershell
   git tag -a v1.0.0 -m "imazeo-go 1.0.0"
   git push origin main
   git push origin v1.0.0
   ```
3. Confirm the module path in `go.mod` is exactly `github.com/Imazeohq/imazeo-go`
   (must match the GitHub owner/repo for `go get` to resolve).
4. Trigger pkg.go.dev indexing (optional but recommended):
   ```text
   https://pkg.go.dev/github.com/Imazeohq/imazeo-go@v1.0.0
   ```
   Or: `GOPROXY=https://proxy.golang.org go list -m github.com/Imazeohq/imazeo-go@v1.0.0`
5. Verify: `go get github.com/Imazeohq/imazeo-go@v1.0.0` in a temp module.

## After publish

Smoke-install:

```powershell
mkdir $env:TEMP\imazeo-go-smoke
cd $env:TEMP\imazeo-go-smoke
go mod init smoke
go get github.com/Imazeohq/imazeo-go@v1.0.0
```

If the release is bad, abandon the version and ship a reviewed patch tag — do not
reuse `v1.0.0`.
