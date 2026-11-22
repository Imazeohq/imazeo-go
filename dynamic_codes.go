package imazeo

import "context"

type DynamicCodesResource struct {
	projectResource
}

func (r *DynamicCodesResource) Capabilities(ctx context.Context) (*DynamicCapabilities, error) {
	return projectJSON[DynamicCapabilities](&r.projectResource, ctx, "GET", []string{"dynamic-codes", "capabilities"}, nil, nil, nil)
}

func (r *DynamicCodesResource) List(ctx context.Context, query map[string]any) (*Page[DynamicCode], error) {
	return projectPage[DynamicCode](&r.projectResource, ctx, []string{"dynamic-codes"}, query)
}

func (r *DynamicCodesResource) Get(ctx context.Context, id string) (*DynamicCode, error) {
	return projectJSON[DynamicCode](&r.projectResource, ctx, "GET", []string{"dynamic-codes", id}, nil, nil, nil)
}

func (r *DynamicCodesResource) Create(ctx context.Context, input DynamicCreateInput, idempotencyKey string) (*DynamicCode, error) {
	return projectJSON[DynamicCode](&r.projectResource, ctx, "POST", []string{"dynamic-codes"}, input, nil, (ProcessingOptions{IdempotencyKey: idempotencyKey}).headers())
}

func (r *DynamicCodesResource) Update(ctx context.Context, id string, input DynamicEditInput, idempotencyKey string) (*DynamicCode, error) {
	return projectJSON[DynamicCode](&r.projectResource, ctx, "PATCH", []string{"dynamic-codes", id}, input, nil, (ProcessingOptions{IdempotencyKey: idempotencyKey}).headers())
}

func (r *DynamicCodesResource) Analytics(ctx context.Context, id string, query map[string]any) (*ScanAnalytics, error) {
	return projectJSON[ScanAnalytics](&r.projectResource, ctx, "GET", []string{"dynamic-codes", id, "analytics"}, nil, query, nil)
}

func (r *DynamicCodesResource) ExportAnalytics(ctx context.Context, id string, query map[string]any) (*BinaryResult, error) {
	return r.binary(ctx, []string{"dynamic-codes", id, "analytics", "export"}, query)
}
