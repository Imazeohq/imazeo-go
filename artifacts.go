package imazeo

import (
	"context"
	"encoding/json"
)

type ArtifactsResource struct {
	projectResource
}

func (r *ArtifactsResource) List(ctx context.Context, query map[string]any) (*Page[ToolArtifact], error) {
	return projectPage[ToolArtifact](&r.projectResource, ctx, []string{"tools", "artifacts"}, query)
}

func (r *ArtifactsResource) Get(ctx context.Context, id string) (*ToolArtifact, error) {
	return projectJSON[ToolArtifact](&r.projectResource, ctx, "GET", []string{"tools", "artifacts", id}, nil, nil, nil)
}

func (r *ArtifactsResource) Delete(ctx context.Context, id string) (*MessageResult, error) {
	return projectJSON[MessageResult](&r.projectResource, ctx, "DELETE", []string{"tools", "artifacts", id}, nil, nil, nil)
}

func (r *ArtifactsResource) Download(ctx context.Context, id string) (*BinaryResult, error) {
	return r.binary(ctx, []string{"tools", "artifacts", id, "content"}, nil)
}

func (r *ArtifactsResource) Save(ctx context.Context, data []byte, input ArtifactSaveInput) (*ToolArtifact, error) {
	rawURL, err := r.path(ctx, []string{"tools", "artifacts"}, nil)
	if err != nil {
		return nil, err
	}
	fields := map[string]string{
		"kind":        input.Kind,
		"displayName": input.DisplayName,
		"contentType": input.ContentType,
	}
	if input.RequestSnapshot != nil {
		encoded, err := json.Marshal(input.RequestSnapshot)
		if err != nil {
			return nil, err
		}
		fields["requestSnapshot"] = string(encoded)
	}
	if input.SourceImageID != "" {
		fields["sourceImageId"] = input.SourceImageID
	}
	file := prepareBytes(data, nil, "output.bin", input.ContentType)
	length := int64(len(data))
	raw, _, err := r.http.requestMultipartJSON(ctx, "POST", rawURL, []multipartFile{file}, fields, &length, r.http.config.timeout)
	if err != nil {
		return nil, err
	}
	var result ToolArtifact
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
