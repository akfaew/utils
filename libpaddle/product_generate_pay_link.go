package libpaddle

import (
	"context"
)

type ProductPayLink struct {
	URL string `json:"url"`
}

type ProductPayLinkResponse struct {
	Success  bool           `json:"success"`
	Response ProductPayLink `json:"response"`
}

// https://developer.paddle.com/api-reference/product-api/pay-links/createpaylink
type ProductGeneratePayLinkOptions struct {
	VendorID       int    `url:"vendor_id"`
	VendorAuthCode string `url:"vendor_auth_code"`

	ProductID               int      `url:"product_id,omitempty"`
	Title                   string   `url:"title,omitempty"`
	WebhookURL              string   `url:"webhook_url,omitempty"`
	Prices                  []string `url:"prices,brackets,omitempty"`
	RecurringPrices         []string `url:"recurring_prices,brackets,omitempty"`
	TrialDays               int      `url:"trial_days,omitempty"`
	CustomMessage           string   `url:"custom_message,omitempty"`
	CouponCode              string   `url:"coupon_code,omitempty"`
	Discountable            int      `url:"discountable,omitempty"`
	ImageURL                string   `url:"image_url,omitempty"`
	ReturnURL               string   `url:"return_url,omitempty"`
	QuantityVariable        int      `url:"quantity_variable,omitempty"`
	Quantity                int      `url:"quantity,omitempty"`
	Expires                 string   `url:"expires,omitempty"`
	Affiliates              []string `url:"affiliates,brackets,omitempty"`
	RecurringAffiliateLimit int      `url:"recurring_affiliate_limit,omitempty"`
	MarketingConsent        int      `url:"marketing_consent,omitempty"`
	CustomerEmail           string   `url:"customer_email,omitempty"`
	CustomerCountry         string   `url:"customer_country,omitempty"`
	CustomerPostcode        string   `url:"customer_postcode,omitempty"`
	Passthrough             string   `url:"passthrough,omitempty"`
	VatNumber               string   `url:"vat_number,omitempty"`
	VatCompanyName          string   `url:"vat_company_name,omitempty"`
	VatStreet               string   `url:"vat_street,omitempty"`
	VatCity                 string   `url:"vat_city,omitempty"`
	VatState                string   `url:"vat_state,omitempty"`
	VatCountry              string   `url:"vat_country,omitempty"`
	VatPostcode             string   `url:"vat_postcode,omitempty"`
}

func (s *ProductService) GeneratePayLink(ctx context.Context, options *ProductGeneratePayLinkOptions) (*ProductPayLinkResponse, error) {
	options.VendorID = s.client.conf.VendorID
	options.VendorAuthCode = s.client.conf.APIKey
	options.ProductID = s.client.conf.ProductID
	u, err := addOptions("product/generate_pay_link", options)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	paylink := new(ProductPayLinkResponse)
	_, err = s.client.Do(ctx, req, paylink)

	return paylink, err
}

func (s *ProductService) GeneratePayLinkCustom(ctx context.Context, options *ProductGeneratePayLinkOptions) (*ProductPayLinkResponse, error) {
	options.VendorID = s.client.conf.VendorID
	options.VendorAuthCode = s.client.conf.APIKey
	u, err := addOptions("product/generate_pay_link", options)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	paylink := new(ProductPayLinkResponse)
	_, err = s.client.Do(ctx, req, paylink)

	return paylink, err
}
