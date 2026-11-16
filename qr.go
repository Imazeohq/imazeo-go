package imazeo

import "context"

// QRResource generates QR codes.
type QRResource struct {
	http      *httpTransport
	projectID func(context.Context) (string, error)
}

// Generate creates a QR code image (binary response).
func (r *QRResource) Generate(ctx context.Context, input QRGenerateInput, processing ...ProcessingOptions) (*BinaryResult, error) {
	projectID, err := r.projectID(ctx)
	if err != nil {
		return nil, err
	}
	rawURL := r.http.path("projects", projectID, "tools", "qr")
	return r.http.requestBinary(ctx, "POST", rawURL, input, processingHeaders(processing), r.http.config.lightTO)
}
