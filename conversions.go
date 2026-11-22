package imazeo

import (
	"encoding/json"
	"context"
)
type ConversionsResource struct {
	projectResource
}

func (r *ConversionsResource) Capabilities(ctx context.Context) (*ConversionCapabilities, error) {
	return projectJSON[ConversionCapabilities](&r.projectResource, ctx, "GET", []string{"conversions", "capabilities"}, nil, nil, nil)
}

func (r *ConversionsResource) Get(ctx context.Context, id string) (*ConversionJob, error) {
	return projectJSON[ConversionJob](&r.projectResource, ctx, "GET", []string{"conversions", id}, nil, nil, nil)
}

func (r *ConversionsResource) Cancel(ctx context.Context, id string) (*ConversionJob, error) {
	return projectJSON[ConversionJob](&r.projectResource, ctx, "POST", []string{"conversions", id, "cancel"}, nil, nil, nil)
}

func (r *ConversionsResource) Save(ctx context.Context, id string) (*ToolArtifact, error) {
	return projectJSON[ToolArtifact](&r.projectResource, ctx, "POST", []string{"conversions", id, "save"}, nil, nil, nil)
}

func (r *ConversionsResource) Download(ctx context.Context, id string) (*BinaryResult, error) {
	return r.binary(ctx, []string{"conversions", id, "download"}, nil)
}

func (r *ConversionsResource) List(ctx context.Context, engine string) ([]ConversionJob, error) {
	result, err := projectJSON[[]ConversionJob](&r.projectResource, ctx, "GET", []string{"conversions"}, nil, map[string]any{"engine": engine}, nil)
	if err != nil {
		return nil, err
	}
	return *result, nil
}

func (r *ConversionsResource) Lookup(ctx context.Context, key string) (*ConversionJob, error) {
	return projectJSON[ConversionJob](&r.projectResource, ctx, "GET", []string{"conversions", "lookup"}, nil, map[string]any{"key": key}, nil)
}

func (r *ConversionsResource) Website(ctx context.Context, sourceURL string, request ProcessingOptions, options *WebsiteOptions) (*ConversionJob, error) {
	return projectJSON[ConversionJob](&r.projectResource, ctx, "POST", []string{"conversions", "website"}, map[string]string{"url": sourceURL}, options, request.headers())
}

func (r *ConversionsResource) upload(ctx context.Context, engine string, data []byte, request ProcessingOptions, options any) (*ConversionJob, error) {
	parts := []string{"conversions"}
	if engine != "docx" {
		parts = append(parts, engine)
	}
	rawURL, err := r.path(ctx, parts, options)
	if err != nil {
		return nil, err
	}
	file := prepareBytes(data, nil, "document." + engine, "")
	length := int64(len(data))
	raw, _, err := r.http.requestMultipartJSON(ctx, "POST", rawURL, []multipartFile{file}, nil, &length, r.http.config.timeout, request.headers())
	if err != nil {
		return nil, err
	}
	var result ConversionJob
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *ConversionsResource) DOCX(ctx context.Context, data []byte, request ProcessingOptions, options *DocxOptions) (*ConversionJob, error) {
	return r.upload(ctx, "docx", data, request, options)
}

func (r *ConversionsResource) XLSX(ctx context.Context, data []byte, request ProcessingOptions, options *XlsxOptions) (*ConversionJob, error) {
	return r.upload(ctx, "xlsx", data, request, options)
}

func (r *ConversionsResource) PPTX(ctx context.Context, data []byte, request ProcessingOptions, options *PptxOptions) (*ConversionJob, error) {
	return r.upload(ctx, "pptx", data, request, options)
}
