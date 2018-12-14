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

//PtrKindNode represents nodes in the resource tree of type reflect.Kind.Ptr, reflect.Kind.UnsafePointer, etc.
type PtrKindNode struct {
	nodeData
	objKind reflect.Kind // Type of the object being pointed to. Eg: *int will store reflect.Kind.Int. This is required for type-expansion and value-evaluation decisions.
}

func (p *PtrKindNode) getData() nodeData {
	return p.nodeData
}

func (p *PtrKindNode) initialize(field string, parent INode, t reflect.Type, rt *ResourceTree) {
	p.nodeData.initialize(field, parent, t, rt)
	p.objKind = t.Elem().Kind()
}