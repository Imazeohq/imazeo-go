package imazeo

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
)

// ImagesResource manages project images.
type ImagesResource struct {
	http      *httpTransport
	projectID func(context.Context) (string, error)
}

// Upload uploads image bytes as multipart field "file" with X-Upload-Length.
func (r *ImagesResource) Upload(ctx context.Context, data []byte, opts *UploadOptions) (*ImageRecord, error) {
	projectID, err := r.projectID(ctx)
	if err != nil {
		return nil, err
	}
	file := prepareBytes(data, opts, "file.bin", "application/octet-stream")
	uploadLen := int64(len(data))
	rawURL := r.http.path("projects", projectID, "images")
	raw, _, err := r.http.requestMultipartJSON(ctx, "POST", rawURL, []multipartFile{file}, nil, &uploadLen, r.http.config.timeout)
	if err != nil {
		return nil, err
	}
	var image ImageRecord
	if err := json.Unmarshal(raw, &image); err != nil {
		return nil, newNetworkError("response was not valid JSON", err, "")
	}
	return &image, nil
}

// List returns a page of images. Query param q filters by filename prefix.
func (r *ImagesResource) List(ctx context.Context, params PageParams) (*Page[ImageRecord], error) {
	projectID, err := r.projectID(ctx)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	if params.Page > 0 {
		q.Set("page", strconv.Itoa(params.Page))
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Q != "" {
		q.Set("q", params.Q)
	}
	rawURL := r.http.path("projects", projectID, "images")
	if qs := q.Encode(); qs != "" {
		rawURL += "?" + qs
	}
	res, err := doJSON[[]ImageRecord](r.http, ctx, "GET", rawURL, nil, r.http.config.timeout)
	if err != nil {
		return nil, err
	}
	page := params.Page
	if page <= 0 {
		page = 1
	}
	limit := params.Limit
	if limit <= 0 {
		limit = 20
	}
	return &Page[ImageRecord]{
		Items:     res.Data,
		Page:      intPtrOr(res.Meta.Page, page),
		Limit:     intPtrOr(res.Meta.Limit, limit),
		Total:     intPtrOr(res.Meta.Total, len(res.Data)),
		Pages:     intPtrOr(res.Meta.Pages, 1),
		RequestID: res.Meta.RequestID,
	}, nil
}

// Get returns image metadata by id.
func (r *ImagesResource) Get(ctx context.Context, id string) (*ImageRecord, error) {
	projectID, err := r.projectID(ctx)
	if err != nil {
		return nil, err
	}
	rawURL := r.http.path("projects", projectID, "images", id)
	res, err := doJSON[ImageRecord](r.http, ctx, "GET", rawURL, nil, r.http.config.timeout)
	if err != nil {
		return nil, err
	}
	return &res.Data, nil
}

// Delete removes an image by id.
func (r *ImagesResource) Delete(ctx context.Context, id string) (*MessageResult, error) {
	projectID, err := r.projectID(ctx)
	if err != nil {
		return nil, err
	}
	rawURL := r.http.path("projects", projectID, "images", id)
	res, err := doJSON[MessageResult](r.http, ctx, "DELETE", rawURL, nil, r.http.config.timeout)
	if err != nil {
		return nil, err
	}
	return &res.Data, nil
}

// Download returns original image bytes.
func (r *ImagesResource) Download(ctx context.Context, id string) (*BinaryResult, error) {
	projectID, err := r.projectID(ctx)
	if err != nil {
		return nil, err
	}
	rawURL := r.http.path("projects", projectID, "images", id, "content")
	return r.http.requestBinary(ctx, "GET", rawURL, nil, nil, r.http.config.timeout)
}
