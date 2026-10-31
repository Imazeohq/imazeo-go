package imazeo

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// TransformsResource creates and manages image transforms.
type TransformsResource struct {
	http      *httpTransport
	projectID func(context.Context) (string, error)
}

// Create runs a transform and returns binary output.
func (r *TransformsResource) Create(ctx context.Context, req TransformRequest) (*BinaryResult, error) {
	if req.ImageID == "" && req.URL == "" {
		return nil, fmt.Errorf("exactly one of imageId or url is required")
	}
	if req.ImageID != "" && req.URL != "" {
		return nil, fmt.Errorf("provide only one of imageId or url")
	}
	projectID, err := r.projectID(ctx)
	if err != nil {
		return nil, err
	}
	rawURL := r.http.path("projects", projectID, "transforms")
	return r.http.requestBinary(ctx, "POST", rawURL, req, map[string]string{
		"Content-Type": "application/json",
		"Accept":       "*/*",
	}, r.http.config.timeout)
}

// List returns a page of transformations.
func (r *TransformsResource) List(ctx context.Context, params PageParams) (*Page[TransformationRecord], error) {
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
	rawURL := r.http.path("projects", projectID, "transforms")
	if qs := q.Encode(); qs != "" {
		rawURL += "?" + qs
	}
	res, err := doJSON[[]TransformationRecord](r.http, ctx, "GET", rawURL, nil, r.http.config.timeout)
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
	return &Page[TransformationRecord]{
		Items:     res.Data,
		Page:      intPtrOr(res.Meta.Page, page),
		Limit:     intPtrOr(res.Meta.Limit, limit),
		Total:     intPtrOr(res.Meta.Total, len(res.Data)),
		Pages:     intPtrOr(res.Meta.Pages, 1),
		RequestID: res.Meta.RequestID,
	}, nil
}

// Get returns transform metadata by id.
func (r *TransformsResource) Get(ctx context.Context, id string) (*TransformationRecord, error) {
	projectID, err := r.projectID(ctx)
	if err != nil {
		return nil, err
	}
	rawURL := r.http.path("projects", projectID, "transforms", id)
	res, err := doJSON[TransformationRecord](r.http, ctx, "GET", rawURL, nil, r.http.config.timeout)
	if err != nil {
		return nil, err
	}
	return &res.Data, nil
}

// DownloadOptions configure transform content download.
type DownloadOptions struct {
	IfNoneMatch string
}

// Download returns transform output bytes. HTTP 304 yields NotModified=true.
func (r *TransformsResource) Download(ctx context.Context, id string, opts *DownloadOptions) (*BinaryResult, error) {
	projectID, err := r.projectID(ctx)
	if err != nil {
		return nil, err
	}
	rawURL := r.http.path("projects", projectID, "transforms", id, "content")
	headers := map[string]string{}
	if opts != nil && opts.IfNoneMatch != "" {
		headers["If-None-Match"] = opts.IfNoneMatch
	}
	return r.http.requestBinary(ctx, "GET", rawURL, nil, headers, r.http.config.timeout)
}

// Delete removes a stored transformation.
func (r *TransformsResource) Delete(ctx context.Context, id string) (*MessageResult, error) {
	projectID, err := r.projectID(ctx)
	if err != nil {
		return nil, err
	}
	rawURL := r.http.path("projects", projectID, "transforms", id)
	res, err := doJSON[MessageResult](r.http, ctx, "DELETE", rawURL, nil, r.http.config.timeout)
	if err != nil {
		return nil, err
	}
	return &res.Data, nil
}
