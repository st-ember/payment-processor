package stripeadapter

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/checkout/session"
)

var stripeSecret string

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	stripeSecret = os.Getenv("STRIPE_SECRET")
}

type CheckoutSessionService struct{}

func NewCheckoutSessionService() *CheckoutSessionService {
	return &CheckoutSessionService{}
}

func (s *CheckoutSessionService) StartSession(stripeParams []*stripe.CheckoutSessionLineItemParams) (string, error) {
	stripe.Key = stripeSecret
	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String("htt"),
		LineItems:  stripeParams,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
	}
	result, err := session.New(params)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

type SessionStarter interface {
	StartSession(stripeParams []*stripe.CheckoutSessionLineItemParams) (string, error)
}
