package forwarder

// Forwarder perform waiting and receiving of events and forwarding them to other place
type Forwarder interface {
	Forward()
}
