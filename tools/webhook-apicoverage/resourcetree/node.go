/*
Copyright 2018 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package resourcetree

import (
	"reflect"

	"github.com/knative/test-infra/tools/webhook-apicoverage/coveragecalculator"
)

// node.go contains types and interfaces pertaining to nodes inside resource tree.

// NodeInterface defines methods that can be performed on each node in the resource tree.
type NodeInterface interface {
	getData() nodeData
	initialize(field string, parent NodeInterface, t reflect.Type, rt *ResourceTree)
	buildChildNodes(t reflect.Type)
	updateCoverage(v reflect.Value)
	buildCoverageData(typeCoverage []coveragecalculator.TypeCoverage, nodeRules NodeRules, fieldRules FieldRules)
	getValues() ([]string)
}

//nodeData is the data stored in each node of the resource tree.
type nodeData struct {
	field string // Represents the Name of the field e.g. field name inside the struct.
	tree *ResourceTree // Reference back to the resource tree. Required for cross-tree traversal(connected nodes traversal)
	fieldType reflect.Type // Required as type information is not available during tree traversal.
	nodePath string // Path in the resource tree reaching this node.
	parent NodeInterface // Link back to parent.
	children map[string]NodeInterface // Child nodes are keyed using field names(nodeData.field).
	leafNode bool // Storing this as an additional field because type-analysis determines the value, which gets used later in value-evaluation
	covered bool
}

func (nd *nodeData) initialize(field string, parent NodeInterface, t reflect.Type, rt *ResourceTree) {
	nd.field = field
	nd.tree = rt
	nd.parent = parent
	nd.fieldType = t
	nd.children = make(map[string]NodeInterface)

	if parent != nil {
		nd.nodePath = parent.getData().nodePath + "." + field
	} else {
		nd.nodePath = field
	}
}