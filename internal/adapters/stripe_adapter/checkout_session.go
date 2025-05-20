package service

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
)

var stripeSecret string

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	stripeSecret = os.Getenv("STRIPE_SECRET")
}

type CheckoutSessionUtil struct{}

func NewCheckoutSessionUtil() *CheckoutSessionUtil {
	return &CheckoutSessionUtil{}
}

func (s *CheckoutSessionUtil) StartSession(
	stripeParams []*stripe.CheckoutSessionLineItemParams) (string, string, error) {
	stripe.Key = stripeSecret
	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String("htt"),
		LineItems:  stripeParams,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
	}
	result, err := session.New(params)
	if err != nil {
		return "", "", err
	}
	return result.ID, result.URL, nil
}

type SessionStarter interface {
	StartSession(stripeParams []*stripe.CheckoutSessionLineItemParams) (string, string, error)
}
