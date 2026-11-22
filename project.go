package imazeo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type projectResource struct {
	http      *httpTransport
	projectID func(context.Context) (string, error)
}

func (r *projectResource) path(ctx context.Context, parts []string, query any) (string, error) {
	id, err := r.projectID(ctx)
	if err != nil {
		return "", err
	}
	raw := r.http.path(append([]string{"projects", id}, parts...)...)
	if query == nil {
		return raw, nil
	}
	encoded, err := json.Marshal(query)
	if err != nil {
		return "", err
	}
	var params map[string]any
	if err := json.Unmarshal(encoded, &params); err != nil {
		return "", err
	}
	values := url.Values{}
	for key, value := range params {
		if value != nil {
			values.Set(key, fmt.Sprint(value))
		}
	}
	if len(values) > 0 {
		raw += "?" + values.Encode()
	}
	return raw, nil
}

func projectJSON[T any](r *projectResource, ctx context.Context, method string, parts []string, body any, query any, headers map[string]string) (*T, error) {
	raw, err := r.path(ctx, parts, query)
	if err != nil {
		return nil, err
	}
	result, err := doJSON[T](r.http, ctx, method, raw, body, r.http.config.lightTO, headers)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func projectPage[T any](r *projectResource, ctx context.Context, parts []string, query map[string]any) (*Page[T], error) {
	raw, err := r.path(ctx, parts, query)
	if err != nil {
		return nil, err
	}
	result, err := doJSON[[]T](r.http, ctx, "GET", raw, nil, r.http.config.lightTO)
	if err != nil {
		return nil, err
	}
	return &Page[T]{
		Items:     result.Data,
		Page:      intPtrOr(result.Meta.Page, 1),
		Limit:     intPtrOr(result.Meta.Limit, 20),
		Total:     intPtrOr(result.Meta.Total, len(result.Data)),
		Pages:     intPtrOr(result.Meta.Pages, 1),
		RequestID: result.Meta.RequestID,
	}, nil
}

func (r *projectResource) binary(ctx context.Context, parts []string, query map[string]any) (*BinaryResult, error) {
	raw, err := r.path(ctx, parts, query)
	if err != nil {
		return nil, err
	}
	return r.http.requestBinary(ctx, "GET", raw, nil, nil, r.http.config.timeout)
}
