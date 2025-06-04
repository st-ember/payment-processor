package stripe_adapter

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
)

var stripeSecret string

type CheckoutSessionUtil struct{}

func NewCheckoutSessionUtil() *CheckoutSessionUtil {
	return &CheckoutSessionUtil{}
}

func (s *CheckoutSessionUtil) StartSession(
	stripeParams []*stripe.CheckoutSessionLineItemParams) (string, string, error) {
	stripe.Key = stripeSecret
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
