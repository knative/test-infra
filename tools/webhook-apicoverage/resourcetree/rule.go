package resourcetree

// rule.go contains different rules that can be defined to control resource tree traversal.

// NodeRules encapsulates all the node level rules defined by a repo.
type NodeRules struct {
	rules []func(nodeInterface NodeInterface) bool
}

// Apply runs all the rules defined by a repo against a node.
func (n *NodeRules) Apply(node NodeInterface) bool {
	for _, rule := range n.rules {
		if !rule(node) {
			return false
		}
	}
	return true
}

// FieldRules encapsulates all the field level rules defined by a repo.
type FieldRules struct {
	rules []func(fieldName string) bool
}

// Apply runs all the rules defined by a repo against a field.
func (f *FieldRules) Apply(fieldName string) bool {
	for _, rule := range f.rules {
		if !rule(fieldName) {
			return false
		}
	}
	return true
}