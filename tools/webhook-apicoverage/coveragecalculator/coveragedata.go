package coveragecalculator

// FieldCoverage represents coverage data for a field.
type FieldCoverage struct {
	Field string `json:"Field"`
	Values []string `json:"Values"`
	Coverage bool `json:"Covered"`
}

// Merge operation merges the field coverage data when multiple nodes represent the same type. (e.g. ConnectedNodes traversal)
func (f *FieldCoverage) Merge(coverage bool, values []string) {
	if coverage {
		f.Coverage = coverage
		f.Values = append(f.Values, values...)
	}
}

// TypeCoverage encapsulates type information and field coverage.
type TypeCoverage struct {
	Package string `json:"Package"`
	Type string `json:"Type"`
	Fields map[string]*FieldCoverage `json:"Fields"`
}