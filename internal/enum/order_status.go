package enum

type OrderStatus int

const (
	// inbound
	Created OrderStatus = iota
	PaymentPending
	PaymentFailed
	Processing
	// outbound
	ShipmentPending
	EnRoute
	// abort
	Refunded
	Returned
	ReturnPending
	Canceled
	// post-delivery
	Delivered
	Complete
)
