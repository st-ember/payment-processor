package stripe_adapter

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
)

type CheckoutSessionUtil struct {
	secret     string
	successUrl string
}

func NewCheckoutSessionUtil(secret, successUrl string) *CheckoutSessionUtil {
	return &CheckoutSessionUtil{
		secret:     secret,
		successUrl: successUrl,
	}
}

func (s *CheckoutSessionUtil) StartSession(
	stripeParams []*stripe.CheckoutSessionLineItemParams) (string, string, error) {
	stripe.Key = s.secret
	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String("http"),
		LineItems:  stripeParams,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
	}
	result, err := session.New(params)
	if err != nil {
		return "", "", err
	}
	return result.ID, result.URL, nil
}
