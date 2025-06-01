package enum

import "fmt"

type StripeStatus int

const (
	Open StripeStatus = iota
	Complete
	Expired
)

func (s StripeStatus) String() string {
	switch s {
	case Open:
		return "open"
	case Complete:
		return "complete"
	case Expired:
		return "expired"
	default:
		return "unknown"
	}
}

func ParseStripeStatus(s string) (StripeStatus, error) {
	switch s {
	case "open":
		return Open, nil
	case "complete":
		return Complete, nil
	case "expired":
		return Expired, nil
	default:
		return -1, fmt.Errorf("invalid StripeStatus: %s", s)
	}
}
