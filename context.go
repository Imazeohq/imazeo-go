package imazeo

import "context"

// ContextResource exposes GET /v1/api-key/context.
type ContextResource struct {
	http *httpTransport
}

// Get returns the API key project and key IDs.
func (r *ContextResource) Get(ctx context.Context) (*APIKeyContext, error) {
	url := r.http.path("api-key", "context")
	res, err := doJSON[APIKeyContext](r.http, ctx, "GET", url, nil, r.http.config.lightTO)
	if err != nil {
		return nil, err
	}
	return &res.Data, nil
}
