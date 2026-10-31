package imazeo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// PDFResource exposes PDF info, render, merge, and extract tools.
type PDFResource struct {
	http      *httpTransport
	projectID func(context.Context) (string, error)
}

// Info returns PDF metadata for an uploaded document.
func (r *PDFResource) Info(ctx context.Context, data []byte, opts *UploadOptions) (*PDFInfo, error) {
	projectID, err := r.projectID(ctx)
	if err != nil {
		return nil, err
	}
	file := prepareBytes(data, opts, "document.pdf", "application/pdf")
	rawURL := r.http.path("projects", projectID, "tools", "pdf", "info")
	raw, _, err := r.http.requestMultipartJSON(ctx, "POST", rawURL, []multipartFile{file}, nil, nil, r.http.config.pdfTO)
	if err != nil {
		return nil, err
	}
	var info PDFInfo
	if err := json.Unmarshal(raw, &info); err != nil {
		return nil, newNetworkError("response was not valid JSON", err, "")
	}
	return &info, nil
}

// Render renders PDF pages. Single-page responses are binary; multi-page may be JSON.
func (r *PDFResource) Render(ctx context.Context, data []byte, opts *PDFRenderOptions) (*PDFRenderResult, error) {
	projectID, err := r.projectID(ctx)
	if err != nil {
		return nil, err
	}
	uploadOpts := &UploadOptions{Filename: "document.pdf", ContentType: "application/pdf"}
	if opts != nil && opts.Filename != "" {
		uploadOpts.Filename = opts.Filename
	}
	file := prepareBytes(data, uploadOpts, "document.pdf", "application/pdf")
	q := url.Values{}
	if opts != nil {
		if opts.Pages != "" {
			q.Set("pages", opts.Pages)
		}
		if opts.RangeStart != nil {
			q.Set("rangeStart", strconv.Itoa(*opts.RangeStart))
		}
		if opts.RangeEnd != nil {
			q.Set("rangeEnd", strconv.Itoa(*opts.RangeEnd))
		}
		if opts.All != nil {
			q.Set("all", strconv.FormatBool(*opts.All))
		}
		if opts.Format != "" {
			q.Set("format", opts.Format)
		}
		if opts.Scale != nil {
			q.Set("scale", strconv.FormatFloat(*opts.Scale, 'f', -1, 64))
		}
		if opts.Quality != nil {
			q.Set("quality", strconv.Itoa(*opts.Quality))
		}
	}
	rawURL := r.http.path("projects", projectID, "tools", "pdf", "render")
	if qs := q.Encode(); qs != "" {
		rawURL += "?" + qs
	}
	bin, jsonData, meta, err := r.http.requestMultipartBinary(ctx, "POST", rawURL, []multipartFile{file}, nil, r.http.config.pdfTO)
	if err != nil {
		return nil, err
	}
	if bin != nil {
		return &PDFRenderResult{Binary: bin, RequestID: bin.RequestID}, nil
	}
	var payload struct {
		Pages []PDFRenderPage `json:"pages"`
	}
	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return nil, newNetworkError("response was not valid JSON", err, meta.RequestID)
	}
	return &PDFRenderResult{Pages: payload.Pages, RequestID: meta.RequestID}, nil
}

// Merge concatenates one or more PDFs into a single PDF binary.
func (r *PDFResource) Merge(ctx context.Context, files [][]byte) (*BinaryResult, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("at least one PDF is required")
	}
	projectID, err := r.projectID(ctx)
	if err != nil {
		return nil, err
	}
	parts := make([]multipartFile, 0, len(files))
	for i, data := range files {
		parts = append(parts, prepareBytes(data, &UploadOptions{
			Filename:    fmt.Sprintf("file-%d.pdf", i+1),
			ContentType: "application/pdf",
		}, fmt.Sprintf("file-%d.pdf", i+1), "application/pdf"))
	}
	rawURL := r.http.path("projects", projectID, "tools", "pdf", "merge")
	bin, _, _, err := r.http.requestMultipartBinary(ctx, "POST", rawURL, parts, nil, r.http.config.pdfTO)
	if err != nil {
		return nil, err
	}
	if bin == nil {
		return nil, fmt.Errorf("expected PDF binary response from merge")
	}
	return bin, nil
}

// Extract returns a PDF containing selected pages.
func (r *PDFResource) Extract(ctx context.Context, data []byte, opts PDFExtractOptions) (*BinaryResult, error) {
	projectID, err := r.projectID(ctx)
	if err != nil {
		return nil, err
	}
	uploadOpts := &UploadOptions{Filename: "document.pdf", ContentType: "application/pdf"}
	if opts.Filename != "" {
		uploadOpts.Filename = opts.Filename
	}
	file := prepareBytes(data, uploadOpts, "document.pdf", "application/pdf")
	q := url.Values{}
	if opts.Pages != "" {
		q.Set("pages", opts.Pages)
	}
	if opts.RangeStart != nil {
		q.Set("rangeStart", strconv.Itoa(*opts.RangeStart))
	}
	if opts.RangeEnd != nil {
		q.Set("rangeEnd", strconv.Itoa(*opts.RangeEnd))
	}
	rawURL := r.http.path("projects", projectID, "tools", "pdf", "extract")
	if qs := q.Encode(); qs != "" {
		rawURL += "?" + qs
	}
	bin, _, _, err := r.http.requestMultipartBinary(ctx, "POST", rawURL, []multipartFile{file}, nil, r.http.config.pdfTO)
	if err != nil {
		return nil, err
	}
	if bin == nil {
		return nil, fmt.Errorf("expected PDF binary response from extract")
	}
	return bin, nil
}
