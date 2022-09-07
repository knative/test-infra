package test

// Env is a testable implementation of modules.Environment.
type Env map[string]string

// Get implements modules.Environment.
func (t Env) Get(name string) string {
	return t[name]
}
