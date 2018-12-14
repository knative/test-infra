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
	"fmt"
	"reflect"
	"strconv"
)

//BasicTypeKindNode represents resource tree node of basic types like int, float, etc.
type BasicTypeKindNode struct {
	nodeData
	values map[string]bool // Values seen for this node. Useful for enum types.
	possibleEnum bool // Flag to indicate if this is a possible enum.
}

func (b *BasicTypeKindNode) getData() nodeData {
	return b.nodeData
}

func (b *BasicTypeKindNode) initialize(field string, parent INode, t reflect.Type, rt *ResourceTree) {
	b.nodeData.initialize(field, parent, t, rt)
	b.values = make(map[string]bool)
	b.nodeData.leafNode = true
}

func (b *BasicTypeKindNode) isNil(v reflect.Value) (bool, string) {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() != 0 {
			return false, strconv.Itoa(int(v.Int()))
		}
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.Uint() != 0 {
			return false, strconv.FormatUint(v.Uint(), 10)
		}
	case reflect.Float32, reflect.Float64:
		if v.Float() != 0 {
			return false, fmt.Sprintf("%f", v.Float())
		}
	case reflect.String:
		if v.Len() != 0 {
			return false, v.String()
		}
	case reflect.Bool:
		return false, strconv.FormatBool(v.Bool())
	}

	return true, ""
}