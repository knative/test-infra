# Resource Tree

resourcetree package contains types and interfaces that define a resource
tree(n-ary tree based representation of an API resource). Each resource
tree is composed of nodes which have data encapsulated inside [nodeData](node.go)
and operations that can be performed on each node in the interface [NodeInterface](node.go).
Type of a node is logically defined by reflect.Kind. Each node type is expected
to satisfy the [NodeInterface](node.go) interface.

## Resource Forest

[ResourceForest](resourceforest.go) groups all resource trees that are part
of an API version into a single construct and defines operations that span
across them. As an example for Knative Serving, we will have individual
resource trees for Configuration, Revision, Route and Service, and they
are encapsulated inside a resource forest under version v1alpha1. Example
of an operation that spans resource trees would be to get coverage details
for outlined types connected using ConnectedNodes.

ConnectedNodes represent connections between nodes that are of same
type(reflect.Type) and belong to same package but span across different
trees or branches of same tree. An example of ConnectedNodes would be
v1alpha1.Route.Spec.Traffic and v1alpha1.Route.Status.Traffic, both these
Traffic fields are of type v1alpha1.TrafficTarget, but are present in different
paths inside the resource tree. ConnectedNodes connects these two nodes, and
an outlining of this type would present the coverage across the two branches
and gives a unified view of what fields are covered. ConnectedNodes represent
connections between nodes that are of same type(reflect.Type) and belong to
same package but span across different trees or branches of same tree. An
example of ConnectedNodes would be v1alpha1.Route.Spec.Traffic and
v1alpha1.Route.Status.Traffic, both these Traffic fields are of type
v1alpha1.TrafficTarget, but are present in different paths inside the resource
tree. ConnectedNodes connects these two nodes, and an outlining of this type
would present the coverage across the two branches and gives a unified view of
what fields are covered.

## Type Analysis

A Resource tree is built using reflect.Type Each node type is expected to
implement NodeInterface method *buildChildNodes(t reflect.Type)*. Inside
this method each node creates child nodes based on its type, for e.g.
StructKindNode creates one child for each field defined in the struct.
Type analysis are defined inside [typeanalyzer_tests](buildChildNodes_test.go)

## Value Analysis

A Resource tree is updated using reflect.Value Each node type is expected
to implement NodeInterface method *updateCoverage(v reflect.Value)*.
Inisde this method each node updates its nodeData.covered field based on
whether the reflect.Value parameter being passed is set or not.