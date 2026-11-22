package imazeo

// CodeDesign is an API v1 model.
type CodeDesign struct {
	ModuleShape		*string	`json:"moduleShape,omitempty"`
	FinderOuter		*string	`json:"finderOuter,omitempty"`
	FinderInner		*string	`json:"finderInner,omitempty"`
	Fill			*string	`json:"fill,omitempty"`
	TextAlign		*string	`json:"textAlign,omitempty"`
	FontFamily		*string	`json:"fontFamily,omitempty"`
	FontWeight		*string	`json:"fontWeight,omitempty"`
	BarTextAlign		*string	`json:"barTextAlign,omitempty"`
	BarTextPosition		*string	`json:"barTextPosition,omitempty"`
	ModuleScale		*int	`json:"moduleScale,omitempty"`
	GradientAngle		*int	`json:"gradientAngle,omitempty"`
	LogoScale		*int	`json:"logoScale,omitempty"`
	LogoPadding		*int	`json:"logoPadding,omitempty"`
	CaptionSize		*int	`json:"captionSize,omitempty"`
	QuietZone		*int	`json:"quietZone,omitempty"`
	BarTextSize		*int	`json:"barTextSize,omitempty"`
	FinderColor		*string	`json:"finderColor,omitempty"`
	FinderInnerColor	*string	`json:"finderInnerColor,omitempty"`
	GradientColor		*string	`json:"gradientColor,omitempty"`
	FrameColor		*string	`json:"frameColor,omitempty"`
	CaptionColor		*string	`json:"captionColor,omitempty"`
	Headline		*string	`json:"headline,omitempty"`
	Subtitle		*string	`json:"subtitle,omitempty"`
	SchemaVersion		int	`json:"schemaVersion"`
	Template		string	`json:"template"`
	TemplateVersion		int	`json:"templateVersion"`
	Frame			*string	`json:"frame,omitempty"`
	Caption			*string	`json:"caption,omitempty"`
	LogoImageID		*string	`json:"logoImageId,omitempty"`
}

// DocxOptions is an API v1 model.
type DocxOptions struct {
	Page				*int		`json:"page,omitempty"`
	Scale				*float64	`json:"scale,omitempty"`
	Format				*string		`json:"format,omitempty"`
	Quality				*int		`json:"quality,omitempty"`
	FitSinglePage			*bool		`json:"fitSinglePage,omitempty"`
	HonorWordPageBreaks		*bool		`json:"honorWordPageBreaks,omitempty"`
	RenderHeaders			*bool		`json:"renderHeaders,omitempty"`
	RenderFooters			*bool		`json:"renderFooters,omitempty"`
	RenderFootnotes			*bool		`json:"renderFootnotes,omitempty"`
	RenderEndnotes			*bool		`json:"renderEndnotes,omitempty"`
	RepairBrokenRelationships	*bool		`json:"repairBrokenRelationships,omitempty"`
	UnwrapSdt			*bool		`json:"unwrapSdt,omitempty"`
	CollapseMcAlternate		*bool		`json:"collapseMcAlternate,omitempty"`
	IgnoreImageWrap			*bool		`json:"ignoreImageWrap,omitempty"`
	IgnoreTableWrap			*bool		`json:"ignoreTableWrap,omitempty"`
	ApplyWordChrome			*bool		`json:"applyWordChrome,omitempty"`
	NormalizeTheme			*bool		`json:"normalizeTheme,omitempty"`
	RenderChanges			*bool		`json:"renderChanges,omitempty"`
}

// XlsxOptions is an API v1 model.
type XlsxOptions struct {
	Page		*int		`json:"page,omitempty"`
	Scale		*float64	`json:"scale,omitempty"`
	Format		*string		`json:"format,omitempty"`
	Quality		*int		`json:"quality,omitempty"`
	Engine		*string		`json:"engine,omitempty"`
	Sheet		*int		`json:"sheet,omitempty"`
	RenderMode	*string		`json:"renderMode,omitempty"`
	Range		*string		`json:"range,omitempty"`
	ShowGridLines	*bool		`json:"showGridLines,omitempty"`
	PaperSize	*string		`json:"paperSize,omitempty"`
}

// PptxOptions is an API v1 model.
type PptxOptions struct {
	Page	*int		`json:"page,omitempty"`
	Scale	*float64	`json:"scale,omitempty"`
	Format	*string		`json:"format,omitempty"`
	Quality	*int		`json:"quality,omitempty"`
	Engine	*string		`json:"engine,omitempty"`
}

// WebsiteOptions is an API v1 model.
type WebsiteOptions struct {
	Page			*int		`json:"page,omitempty"`
	Scale			*float64	`json:"scale,omitempty"`
	Format			*string		`json:"format,omitempty"`
	Quality			*int		`json:"quality,omitempty"`
	Engine			*string		`json:"engine,omitempty"`
	Mode			*string		`json:"mode,omitempty"`
	Width			*int		`json:"width,omitempty"`
	Height			*int		`json:"height,omitempty"`
	Selector		*string		`json:"selector,omitempty"`
	ClipX			*int		`json:"clipX,omitempty"`
	ClipY			*int		`json:"clipY,omitempty"`
	ClipWidth		*int		`json:"clipWidth,omitempty"`
	ClipHeight		*int		`json:"clipHeight,omitempty"`
	JavascriptEnabled	*bool		`json:"javascriptEnabled,omitempty"`
	BlockTrackers		*bool		`json:"blockTrackers,omitempty"`
	AutoScroll		*bool		`json:"autoScroll,omitempty"`
	FreezeAnimations	*bool		`json:"freezeAnimations,omitempty"`
	OmitBackground		*bool		`json:"omitBackground,omitempty"`
	ColorScheme		*string		`json:"colorScheme,omitempty"`
	ReducedMotion		*bool		`json:"reducedMotion,omitempty"`
	WaitUntil		*string		`json:"waitUntil,omitempty"`
	WaitMs			*int		`json:"waitMs,omitempty"`
}

// ConversionCapabilities is an API v1 model.
type ConversionCapabilities struct {
	Enabled			bool			`json:"enabled"`
	MaxUploadBytes		int64			`json:"maxUploadBytes"`
	MaxPixels		int64			`json:"maxPixels"`
	MaxOutputBytes		int64			`json:"maxOutputBytes"`
	Formats			[]string		`json:"formats"`
	MaxPage			int			`json:"maxPage"`
	TemporaryMinutes	int			`json:"temporaryMinutes"`
	Engines			[]ConversionEngineInfo	`json:"engines"`
}

// ConversionEngineInfo is an API v1 model.
type ConversionEngineInfo struct {
	ID		string	`json:"id"`
	Enabled		bool	`json:"enabled"`
	Label		string	`json:"label"`
	Unit		string	`json:"unit"`
	MaxUnits	int	`json:"maxUnits"`
}

// ConversionJob is an API v1 model.
type ConversionJob struct {
	ID		string			`json:"id"`
	ProjectID	string			`json:"projectId"`
	Status		string			`json:"status"`
	Engine		string			`json:"engine"`
	Options		map[string]any		`json:"options"`
	Credits		int			`json:"credits"`
	CreatedAt	string			`json:"createdAt"`
	ExpiresAt	*string			`json:"expiresAt"`
	ArtifactID	*string			`json:"artifactId"`
	ErrorCode	*string			`json:"errorCode"`
	Result		*ConversionResult	`json:"result"`
}

// ConversionResult is an API v1 model.
type ConversionResult struct {
	Bytes		int64		`json:"bytes"`
	Width		int		`json:"width"`
	Height		int		`json:"height"`
	PageCount	int		`json:"pageCount"`
	ContentType	string		`json:"contentType"`
	Warnings	[]string	`json:"warnings"`
}

// DynamicCapabilities is an API v1 model.
type DynamicCapabilities struct {
	Enabled		bool	`json:"enabled"`
	CanEdit		bool	`json:"canEdit"`
	Analytics	bool	`json:"analytics"`
	ActiveCount	int	`json:"activeCount"`
	MaxActiveCodes	int	`json:"maxActiveCodes"`
	RetentionDays	int	`json:"retentionDays"`
}

// DynamicCode is an API v1 model.
type DynamicCode struct {
	AnalyticsEnabled	bool			`json:"analyticsEnabled"`
	QrOptions		*QRGenerateInput	`json:"qrOptions"`
	ID			string			`json:"id"`
	ProjectID		string			`json:"projectId"`
	Name			string			`json:"name"`
	Slug			string			`json:"slug"`
	RedirectURL		string			`json:"redirectUrl"`
	Destination		*string			`json:"destination"`
	Status			string			`json:"status"`
	Suspended		bool			`json:"suspended"`
	Revision		int			`json:"revision"`
	CreatedAt		string			`json:"createdAt"`
	UpdatedAt		string			`json:"updatedAt"`
	RecordedScans		*int			`json:"recordedScans,omitempty"`
}

// DynamicCreateInput is an API v1 model.
type DynamicCreateInput struct {
	AnalyticsEnabled	*bool	`json:"analyticsEnabled,omitempty"`
	Name			string	`json:"name"`
	Destination		string	`json:"destination"`
}

// DynamicEditInput is an API v1 model.
type DynamicEditInput struct {
	AnalyticsEnabled	*bool			`json:"analyticsEnabled,omitempty"`
	QrOptions		*QRGenerateInput	`json:"qrOptions,omitempty"`
	ExpectedRevision	int			`json:"expectedRevision"`
	Name			*string			`json:"name,omitempty"`
	Destination		*string			`json:"destination,omitempty"`
	Status			*string			`json:"status,omitempty"`
}

// ScanAnalytics is an API v1 model.
type ScanAnalytics struct {
	Days		int		`json:"days"`
	Total		int		`json:"total"`
	RetentionDays	int		`json:"retentionDays"`
	BestEffort	bool		`json:"bestEffort"`
	Series		[]ScanPoint	`json:"series"`
	Devices		[]ScanBreakdown	`json:"devices"`
	Browsers	[]ScanBreakdown	`json:"browsers"`
	Countries	[]ScanBreakdown	`json:"countries"`
	Recent		[]ScanRecent	`json:"recent"`
}

// ScanPoint is an API v1 model.
type ScanPoint struct {
	Date	string	`json:"date"`
	Count	int	`json:"count"`
}

// ScanBreakdown is an API v1 model.
type ScanBreakdown struct {
	Key	string	`json:"key"`
	Count	int	`json:"count"`
}

// ScanRecent is an API v1 model.
type ScanRecent struct {
	ID		string	`json:"id"`
	CreatedAt	string	`json:"createdAt"`
	Device		string	`json:"device"`
	Browser		string	`json:"browser"`
	Bot		bool	`json:"bot"`
	Country		string	`json:"country"`
	Ip		*string	`json:"ip"`
	ForwardedFor	*string	`json:"forwardedFor"`
	UserAgent	*string	`json:"userAgent"`
	AcceptLanguage	*string	`json:"acceptLanguage"`
	Language	*string	`json:"language"`
	Referer		*string	`json:"referer"`
	Platform	string	`json:"platform"`
	ClientHints	*string	`json:"clientHints"`
	MobileHint	*bool	`json:"mobileHint"`
	Protocol	string	`json:"protocol"`
	Accept		*string	`json:"accept"`
}

// CreditQuote is an API v1 model.
type CreditQuote struct {
	ReplayEnabled	bool	`json:"replayEnabled"`
	ID		string	`json:"id"`
	Kind		string	`json:"kind"`
	Units		int	`json:"units"`
	Credits		int	`json:"credits"`
	UnitCost	int	`json:"unitCost"`
	RemoteUnitCost	int	`json:"remoteUnitCost"`
	ScheduleID	string	`json:"scheduleId"`
	ExpiresAt	string	`json:"expiresAt"`
}

// CreditQuoteInput is an API v1 model.
type CreditQuoteInput struct {
	Kind	string	`json:"kind"`
	Units	int	`json:"units"`
	Remote	*bool	`json:"remote,omitempty"`
}

// ProcessingStatus is an API v1 model.
type ProcessingStatus struct {
	ID		string	`json:"id"`
	Kind		string	`json:"kind"`
	State		string	`json:"state"`
	ResultExpiresAt	string	`json:"resultExpiresAt"`
	HasResult	bool	`json:"hasResult"`
	FailureCode	*string	`json:"failureCode"`
}

// ProcessingLookupInput is an API v1 model.
type ProcessingLookupInput struct {
	Kind	string	`json:"kind"`
	Key	string	`json:"key"`
}

// ToolArtifact is an API v1 model.
type ToolArtifact struct {
	ID		string		`json:"_id"`
	ProjectID	string		`json:"projectId"`
	Kind		string		`json:"kind"`
	DisplayName	string		`json:"displayName"`
	Status		string		`json:"status"`
	Bytes		int64		`json:"bytes"`
	SHA256		string		`json:"sha256"`
	ContentType	string		`json:"contentType"`
	OutputType	string		`json:"outputType"`
	RequestSnapshot	map[string]any	`json:"requestSnapshot"`
	SourceImageID	*string		`json:"sourceImageId"`
	CreatedAt	string		`json:"createdAt"`
	UpdatedAt	string		`json:"updatedAt"`
}

// OutputsSummary is an API v1 model.
type OutputsSummary struct {
	Total		int	`json:"total"`
	Originals	int	`json:"originals"`
	Transforms	int	`json:"transforms"`
	Images		int	`json:"images"`
	Qr		int	`json:"qr"`
	Barcode		int	`json:"barcode"`
	Pdf		int	`json:"pdf"`
	Docx		*int	`json:"docx,omitempty"`
	Xlsx		*int	`json:"xlsx,omitempty"`
	Pptx		*int	`json:"pptx,omitempty"`
	Website		*int	`json:"website,omitempty"`
}

// OutputRow is an API v1 model.
type OutputRow struct {
	ID		string		`json:"id"`
	Kind		string		`json:"kind"`
	Name		string		`json:"name"`
	Type		string		`json:"type"`
	Status		string		`json:"status"`
	CreatedAt	string		`json:"createdAt"`
	Bytes		int64		`json:"bytes"`
	ThumbAvailable	bool		`json:"thumbAvailable"`
	SourceImageID	*string		`json:"sourceImageId,omitempty"`
	ContentType	*string		`json:"contentType,omitempty"`
	Meta		map[string]any	`json:"meta,omitempty"`
}

// ProcessingOptions binds an operation to its quote and idempotency key.
type ProcessingOptions struct {
	IdempotencyKey	string
	QuoteID		string
}

func (o ProcessingOptions) headers() map[string]string {
	headers := map[string]string{}
	if o.IdempotencyKey != "" {
		headers["Idempotency-Key"] = o.IdempotencyKey
	}
	if o.QuoteID != "" {
		headers["X-Quote-ID"] = o.QuoteID
	}
	return headers
}

func processingHeaders(options []ProcessingOptions) map[string]string {
	if len(options) > 0 {
		return options[0].headers()
	}
	return nil
}

type ArtifactSaveInput struct {
	Kind		string
	DisplayName	string
	ContentType	string
	RequestSnapshot	map[string]any
	SourceImageID	string
}
