package sender

// Sender will send messages continuously until process receives a SIGINT
type Sender interface {
	SendContinually()
}
