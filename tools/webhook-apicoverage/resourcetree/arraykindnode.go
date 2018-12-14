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

//ArrayKindNode represents resource tree node of types reflect.Kind.Array and reflect.Kind.Slice
type ArrayKindNode struct {
	nodeData
	arrKind reflect.Kind // Array type e.g. []int will store reflect.Kind.Int. This is required for type-expansion and value-evaluation decisions.
}

func (a *ArrayKindNode ) getData() nodeData {
	return a.nodeData
}

func (a *ArrayKindNode ) initialize(field string, parent INode, t reflect.Type, rt *ResourceTree) {
	a.nodeData.initialize(field, parent, t, rt)
	a.arrKind = t.Elem().Kind()
}