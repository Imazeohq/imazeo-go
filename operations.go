package imazeo

import (
	"encoding/json"
	"strings"
	"context"
)
type OperationsResource struct {
	projectResource
}

func (r *OperationsResource) Get(ctx context.Context, id string) (*ProcessingStatus, error) {
	return projectJSON[ProcessingStatus](&r.projectResource, ctx, "GET", []string{"operations", id}, nil, nil, nil)
}

func (r *OperationsResource) Lookup(ctx context.Context, input ProcessingLookupInput) (*ProcessingStatus, error) {
	return projectJSON[ProcessingStatus](&r.projectResource, ctx, "POST", []string{"operations", "lookup"}, input, nil, nil)
}
// Result returns either binary content or unwrapped structured JSON.
func (r *OperationsResource) Result(ctx context.Context, id string) (*BinaryResult, json.RawMessage, error) {
	result, err := r.binary(ctx, []string{"operations", id, "result"}, nil)
	if err != nil {
		return nil, nil, err
	}
	if strings.Contains(result.ContentType, "application/json") {
		data, _, err := parseSuccessEnvelope(result.Data, result.RequestID)
		return nil, data, err
	}
	return result, nil, nil
}
