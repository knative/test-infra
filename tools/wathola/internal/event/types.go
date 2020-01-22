package event

const (
	// StepType is a string type representation of step event
	StepType     = "dev.knative.wathola.step"
	// FinishedType os a string type representation of finished event
	FinishedType = "dev.knative.wathola.finished"
)

// Step is a event call at each step of verification
type Step struct {
	Number int
}

// Finished is step call after verification finishes
type Finished struct {
	Count int
}

// Type returns a type of a event
func (s Step) Type() string {
	return StepType
}

// Type returns a type of a event
func (f Finished) Type() string {
	return FinishedType
}

// State defines a state of event store
type State int

const (
	// Active == 1 (iota has been reset)
	Active  State = 1 << iota
	// Success == 2
	Success State = 1 << iota
	// Failed == 4
	Failed  State = 1 << iota
)
