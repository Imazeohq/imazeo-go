package imazeo

import "time"

// SDKVersion is the package SemVer (independent of the API path version).
const SDKVersion = "1.0.0"

// DefaultAPIVersion is the default path prefix under the API origin.
const DefaultAPIVersion = "v1"

// SupportedAPIVersions lists allowlisted apiVersion option values.
var SupportedAPIVersions = []string{"v1"}

// DefaultBaseURL is the production API origin (no /v1 suffix).
const DefaultBaseURL = "https://api.imazeo.com"

// Default request timeouts (JSON/light 30s, general/transform 60s, PDF 120s).
const (
	DefaultTimeout = 60 * time.Second
	LightTimeout   = 30 * time.Second
	PDFTimeout     = 120 * time.Second
)
