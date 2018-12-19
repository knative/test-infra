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
)

const (
	v1TimeType = "v1.Time"
	volatileTimeType = "apis.VolatileTime"
)

// StructKindNode represents nodes in the resource tree of type reflect.Kind.Struct
type StructKindNode struct {
	nodeData
}

func (s *StructKindNode) getData() nodeData {
	return s.nodeData
}

func (s *StructKindNode) initialize(field string, parent NodeInterface, t reflect.Type, rt *ResourceTree) {
	s.nodeData.initialize(field, parent, t, rt)
}

func (s *StructKindNode) buildChildNodes(t reflect.Type) {
	// For types that are part of the standard package, we treat them as leaf nodes and don't expand further.
	// https://golang.org/pkg/reflect/#StructField.
	if len(s.fieldType.PkgPath()) == 0 {
		s.leafNode = true
		return
	}

	for i := 0; i < t.NumField(); i++ {
		var childNode NodeInterface
		if s.isTimeNode(t.Field(i).Type) {
			childNode = new(TimeTypeNode)
			childNode.initialize(t.Field(i).Name, s, t.Field(i).Type, s.tree)
		} else {
			childNode = s.tree.createNode(t.Field(i).Name, s, t.Field(i).Type)
		}
		s.children[t.Field(i).Name] = childNode
		childNode.buildChildNodes(t.Field(i).Type)
	}
}

func (s *StructKindNode) isTimeNode(t reflect.Type) bool {
	if t.Kind() == reflect.Struct {
		return t.String() == v1TimeType || t.String() == volatileTimeType
	} else if t.Kind() == reflect.Ptr {
		return t.Elem().String() == v1TimeType || t.String() == volatileTimeType
	} else {
		return false
	}
}

func (s *StructKindNode) updateCoverage(v reflect.Value) {
	if v.IsValid() {
		s.covered = true
		if !s.leafNode {
			for i := 0; i < v.NumField(); i++ {
				s.children[v.Type().Field(i).Name].updateCoverage(v.Field(i))
			}
		}
	}
}