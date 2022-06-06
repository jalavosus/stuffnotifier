package flightawarepoller

type notificationType uint8

const (
	NoNotification notificationType = iota
	GateDepartureNotification
	TakeoffNotification
	LandingNotification
	GateArrivalNotification
	PreDepartureNotification
	PreArrivalNotification
	BaggageClaimNotification
)
