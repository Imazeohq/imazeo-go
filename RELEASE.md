# SDK release procedure (`github.com/Imazeohq/imazeo-go`)

Independent repository: [Imazeohq/imazeo-go](https://github.com/Imazeohq/imazeo-go).  
Module / docs: [pkg.go.dev/github.com/Imazeohq/imazeo-go](https://pkg.go.dev/github.com/Imazeohq/imazeo-go).  
Supported API prefix: `v1`. SemVer lives in `version.go` (`SDKVersion`) and `package.json`
(kept in sync by the release script). Published from git tag `v*`.

Never reuse a published tag; ship a patch instead.

## Release (PowerShell)

Commit your work on `main` first (clean working tree required). One-time: `npm install`.

```powershell
$env:RELEASE_MODE = "true"
npm run release
```

You will be prompted for version bump and GitHub Release.  
This syncs versions, commits, tags `vX.Y.Z`, and pushes.  
**Publish Go module** then verifies the tag and indexes `proxy.golang.org`.

## Release from GitHub only

If the version bump is already on `main`:

1. **Releases** → Draft a new release → tag `vX.Y.Z` → Publish  
   or
2. **Actions** → **Publish Go module** → **Run workflow**

## After publish

Verify: https://pkg.go.dev/github.com/Imazeohq/imazeo-go  
`go get github.com/Imazeohq/imazeo-go@vX.Y.Z`
