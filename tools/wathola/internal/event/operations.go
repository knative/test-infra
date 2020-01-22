package event

// StepsStore contains methods that register step event type
type StepsStore interface {
	RegisterStep(step *Step)
	Count() int
}

// FinishedStore registers a finished event type
type FinishedStore interface {
	RegisterFinished(finished *Finished)
	State() State
	Thrown() []string
}

// Typed says a type of an event
type Typed interface {
	Type() string
}
