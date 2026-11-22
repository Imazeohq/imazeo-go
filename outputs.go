package imazeo

import "context"

type OutputsResource struct {
	projectResource
}

func (r *OutputsResource) Summary(ctx context.Context) (*OutputsSummary, error) {
	return projectJSON[OutputsSummary](&r.projectResource, ctx, "GET", []string{"outputs", "summary"}, nil, nil, nil)
}

func (r *OutputsResource) List(ctx context.Context, query map[string]any) (*Page[OutputRow], error) {
	return projectPage[OutputRow](&r.projectResource, ctx, []string{"outputs"}, query)
}
