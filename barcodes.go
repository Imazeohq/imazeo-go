package imazeo

import "context"

// BarcodesResource generates barcodes.
type BarcodesResource struct {
	http      *httpTransport
	projectID func(context.Context) (string, error)
}

// Generate creates a barcode image (binary response).
func (r *BarcodesResource) Generate(ctx context.Context, input BarcodeGenerateInput, processing ...ProcessingOptions) (*BinaryResult, error) {
	projectID, err := r.projectID(ctx)
	if err != nil {
		return nil, err
	}
	rawURL := r.http.path("projects", projectID, "tools", "barcode")
	return r.http.requestBinary(ctx, "POST", rawURL, input, processingHeaders(processing), r.http.config.lightTO)
}
