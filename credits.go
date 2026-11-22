package imazeo

import "context"

type CreditsResource struct {
	projectResource
}

func (r *CreditsResource) Quote(ctx context.Context, input CreditQuoteInput) (*CreditQuote, error) {
	return projectJSON[CreditQuote](&r.projectResource, ctx, "POST", []string{"credits", "quote"}, input, nil, nil)
}
