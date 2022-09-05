package test

type Env map[string]string

func (t Env) Get(name string) string {
	return t[name]
}
