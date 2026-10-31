package imazeo

// PageParams are common list query parameters.
type PageParams struct {
	Page  int
	Limit int
	// Q is a filename / display-name prefix filter where the resource supports it.
	Q string
}

// Page is a paginated list result.
type Page[T any] struct {
	Items     []T
	Page      int
	Limit     int
	Total     int
	Pages     int
	RequestID string
}

// BinaryResult is a non-enveloped binary API response.
type BinaryResult struct {
	Data             []byte
	ContentType      string
	RequestID        string
	Headers          map[string]string
	ETag             string
	Width            *int
	Height           *int
	Format           string
	TransformationID string
	Cache            string
	PDFPage          *int
	// NotModified is true when the server returned HTTP 304 for If-None-Match.
	NotModified bool
}

// ImageMetadata describes decoded image properties.
type ImageMetadata struct {
	Format      string `json:"format"`
	MimeType    string `json:"mimeType"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Frames      int    `json:"frames"`
	PixelCount  int    `json:"pixelCount"`
	HasAlpha    bool   `json:"hasAlpha"`
	Orientation *int   `json:"orientation"`
}

// ImageRecord is a stored project image.
type ImageRecord struct {
	ID           string        `json:"_id"`
	ProjectID    string        `json:"projectId"`
	SourceType   string        `json:"sourceType"`
	OriginalName string        `json:"originalName"`
	Bytes        int64         `json:"bytes"`
	SHA256       string        `json:"sha256"`
	Metadata     ImageMetadata `json:"metadata"`
	Status       string        `json:"status"`
	CreatedAt    string        `json:"createdAt"`
	UpdatedAt    string        `json:"updatedAt"`
}

// TransformSpec describes image transformation options.
type TransformSpec struct {
	Width         *int            `json:"width,omitempty"`
	Height        *int            `json:"height,omitempty"`
	Fit           string          `json:"fit,omitempty"`
	Crop          *CropRect       `json:"crop,omitempty"`
	Format        string          `json:"format,omitempty"`
	Quality       *int            `json:"quality,omitempty"`
	Rotate        *int            `json:"rotate,omitempty"`
	Flip          *bool           `json:"flip,omitempty"`
	Flop          *bool           `json:"flop,omitempty"`
	Optimize      *bool           `json:"optimize,omitempty"`
	StripMetadata *bool           `json:"stripMetadata,omitempty"`
	Grayscale     *bool           `json:"grayscale,omitempty"`
	Negate        *bool           `json:"negate,omitempty"`
	Blur          *float64        `json:"blur,omitempty"`
	Sharpen       *float64        `json:"sharpen,omitempty"`
	Brightness    *float64        `json:"brightness,omitempty"`
	Saturation    *float64        `json:"saturation,omitempty"`
	Hue           *int        `json:"hue,omitempty"`
	Trim          any             `json:"trim,omitempty"`
	Extend        *ExtendInsets   `json:"extend,omitempty"`
	Background    string          `json:"background,omitempty"`
	Position      string          `json:"position,omitempty"`
}

// CropRect is a pixel crop region.
type CropRect struct {
	Left   int `json:"left"`
	Top    int `json:"top"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// ExtendInsets pads an image on each side.
type ExtendInsets struct {
	Top    int `json:"top"`
	Right  int `json:"right"`
	Bottom int `json:"bottom"`
	Left   int `json:"left"`
}

// TransformMode controls persistence of transform outputs.
type TransformMode string

const (
	TransformModePersist   TransformMode = "PERSIST"
	TransformModeCache     TransformMode = "CACHE"
	TransformModeEphemeral TransformMode = "EPHEMERAL"
)

// TransformRequest creates a transform from an uploaded image or remote URL.
type TransformRequest struct {
	ImageID        string         `json:"imageId,omitempty"`
	URL            string         `json:"url,omitempty"`
	Mode           TransformMode  `json:"mode,omitempty"`
	MaxOutputBytes *int64         `json:"maxOutputBytes,omitempty"`
	Transform      *TransformSpec `json:"transform,omitempty"`
}

// TransformationRecord is stored transform metadata.
type TransformationRecord struct {
	ID            string         `json:"_id"`
	ProjectID     string         `json:"projectId"`
	ImageID       *string        `json:"imageId"`
	SourceHash    string         `json:"sourceHash"`
	TransformHash string         `json:"transformHash"`
	Spec          TransformSpec  `json:"spec"`
	Mode          TransformMode  `json:"mode"`
	Bytes         int64          `json:"bytes"`
	ProcessingMs  int            `json:"processingMs"`
	Status        string         `json:"status"`
	SHA256        string         `json:"sha256"`
	Metadata      ImageMetadata  `json:"metadata"`
	ExpiresAt     *string        `json:"expiresAt"`
	CreatedAt     string         `json:"createdAt"`
}

// APIKeyContext is returned by GET /v1/api-key/context.
type APIKeyContext struct {
	ProjectID   string `json:"projectId"`
	KeyID       string `json:"keyId"`
}

// QRGenerateInput generates a QR code.
type QRGenerateInput struct {
	Content         string `json:"content"`
	Format          string `json:"format,omitempty"`
	Size            int    `json:"size,omitempty"`
	Margin          *int   `json:"margin,omitempty"`
	ErrorCorrection string `json:"errorCorrection,omitempty"`
	Foreground      string `json:"foreground,omitempty"`
	Background      string `json:"background,omitempty"`
}

// BarcodeGenerateInput generates a barcode.
type BarcodeGenerateInput struct {
	Type        string `json:"type"`
	Content     string `json:"content"`
	Format      string `json:"format,omitempty"`
	Scale       *int   `json:"scale,omitempty"`
	Height      *int   `json:"height,omitempty"`
	IncludeText *bool  `json:"includeText,omitempty"`
	Foreground  string `json:"foreground,omitempty"`
	Background  string `json:"background,omitempty"`
}

// PDFInfo is returned by PDF info.
type PDFInfo struct {
	PageCount int            `json:"pageCount"`
	Encrypted bool           `json:"encrypted"`
	Pages     []PDFPageInfo  `json:"pages"`
}

// PDFPageInfo describes a single PDF page.
type PDFPageInfo struct {
	PageNumber int     `json:"pageNumber"`
	Width      float64 `json:"width"`
	Height     float64 `json:"height"`
}

// PDFRenderOptions control PDF page rendering.
type PDFRenderOptions struct {
	Pages      string
	RangeStart *int
	RangeEnd   *int
	All        *bool
	Format     string
	Scale      *float64
	Quality    *int
	Filename   string
}

// PDFRenderPage is one page in a multi-page JSON render response.
type PDFRenderPage struct {
	PageNumber  int    `json:"pageNumber"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	ContentType string `json:"contentType"`
	Data        string `json:"data"`
}

// PDFRenderResult is either a single-page binary or multi-page JSON payload.
type PDFRenderResult struct {
	Binary    *BinaryResult
	Pages     []PDFRenderPage
	RequestID string
}

// PDFExtractOptions control page extraction.
type PDFExtractOptions struct {
	Pages      string
	RangeStart *int
	RangeEnd   *int
	Filename   string
}

// UploadOptions describe a local file for multipart upload (never a filesystem path string).
type UploadOptions struct {
	Filename    string
	ContentType string
}

// MessageResult is a simple `{message}` JSON data payload.
type MessageResult struct {
	Message string `json:"message"`
}
