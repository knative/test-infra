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

//test_util contains types defined and used by types and their corresponding verification methods.

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	basicTypeName = "BasicType"
	ptrTypeName = "PtrType"
	arrayTypeName = "ArrayType"
	otherTypeName = "OtherType"
	combinedTypeName = "CombinedType"
)

type baseType struct {
	field1 string
	field2 int16
}

type ptrType struct {
	structPtr *baseType
	basePtr *float32
}

type arrayType struct {
	structArr []baseType
	baseArr []bool
}

type otherType struct {
	structMap map[string]baseType
	baseMap map[string]string
}

type combinedNodeType struct {
	b baseType
	a arrayType
	p ptrType
}

func verifyBaseTypeNode(logPrefix string, data nodeData) error {
	if len(data.children) != 2 {
		return fmt.Errorf("%s Expected 2 children got only : %d", logPrefix, len(data.children))
	}

	if value, ok := data.children["field1"]; ok {
		n := value.getData()
		if !n.leafNode || n.fieldType.Kind() != reflect.String || n.fieldType.PkgPath() != "" || len(n.children) != 0 {
			return fmt.Errorf("%s Unexpected field: field1. Expected leafNode:true, Kind: %s, pkgPath: '' Children: 0 Found leafNode: %t Kind: %s pkgPath: %s Children:%d",
				logPrefix, reflect.String, n.leafNode, n.fieldType.Kind(), n.fieldType.PkgPath(), len(n.children))
		}
	} else {
		return fmt.Errorf("%s field1 child Not found", logPrefix)
	}

	return nil
}

func verifyPtrNode(data nodeData) error {
	if len(data.children) != 2 {
		return fmt.Errorf("Expected 2 children got: %d", len(data.children))
	}

	child := data.children["structPtr"]
	if len(child.getData().children) != 1 {
		return fmt.Errorf("Unexpected size for field:structPtr. Expected : 1, Found : %d", len(child.getData().children))
	}

	child = child.getData().children["structPtr-ptr"]
	if err := verifyBaseTypeNode("child structPtr-ptr: ", child.getData()); err != nil {
		return err
	}

	child = data.children["basePtr"]
	if len(child.getData().children) != 1 {
		return fmt.Errorf("Unexpected size for field:basePtr. Expected : 1 Found : %d", len(child.getData().children))
	}

	child = child.getData().children["basePtr-ptr"]
	d := child.getData()
	if d.fieldType.Kind() != reflect.Float32 || !d.leafNode || d.fieldType.PkgPath() != "" || len(d.children) != 0 {
		return fmt.Errorf("Unexpected field:basePtr-ptr: Expected: Kind: %s, leafNode: true, pkgPath: '' Children: 0 Found Kind: %s, leafNode: %t, pkgPath: %s children:%d",
			reflect.Float32, d.fieldType.Kind(), d.leafNode, d.fieldType.PkgPath(), len(d.children))
	}

	return nil
}

func verifyArrayNode(data nodeData) error {
	if len(data.children) != 2 {
		return fmt.Errorf("Expected 2 children got: %d", len(data.children))
	}

	child := data.children["structArr"]
	d := child.getData()
	if d.fieldType.Kind() != reflect.Slice {
		return fmt.Errorf("Unexpected kind for field:structArr: Expected : %s Found: %s", reflect.Slice, d.fieldType.Kind())
	} else if len(d.children) != 1 {
		return fmt.Errorf("Unexpected number of children for field:structArr: Expected : 1 Found : %d", len(d.children))
	}

	child = child.getData().children["structArr-arr"]
	if err := verifyBaseTypeNode("child structArr-arr:", child.getData()); err != nil {
		return err
	}

	child = data.children["baseArr"]
	d = child.getData()
	if d.fieldType.Kind() != reflect.Slice {
		return fmt.Errorf("Unexpected kind for field:baseArr: Expected : %s Found : %s", reflect.Slice, d.fieldType.Kind())
	} else if len(d.children) != 1 {
		return fmt.Errorf("Unexpected number of children for field:baseArr: Expected : 1 Found : %d", len(d.children))
	}

	child = child.getData().children["baseArr-arr"]
	d = child.getData()
	if d.fieldType.Kind() != reflect.Bool || !d.leafNode || d.fieldType.PkgPath() != "" || len(d.children) != 0 {
		return fmt.Errorf("Unexpected field:baseArr-arr Expected kind: %s, leafNode: true, pkgPath: '', children:0 Found: kind: %s, leafNode: %t, pkgPath: %s, children:%d",
			reflect.Bool, d.fieldType.Kind(), d.leafNode, d.fieldType.PkgPath(), len(d.children))
	}

	return nil
}

func verifyOtherTypeNode(data nodeData) error {
	if len(data.children) != 2 {
		return fmt.Errorf("OtherTypeVerification: Expected 2 children got: %d", len(data.children))
	}

	child := data.children["structMap"]
	d := child.getData()
	if d.fieldType.Kind() != reflect.Map || !d.leafNode || len(d.children) != 0 {
		return fmt.Errorf("Unexpected field:structMap - Expected Kind: %s, leafNode: true, children:0 Found Kind: %s, leafNode: %t, children: %d",
			reflect.Map, d.fieldType.Kind(), d.leafNode, len(d.children))
	}

	child = data.children["baseMap"]
	d = child.getData()
	if d.fieldType.Kind() != reflect.Map || !d.leafNode || len(d.children) != 0 {
		return fmt.Errorf("Unexpected field:structMap - Expected Kind: %s, leafNode: true, children: 0 Found kind: %s, leafNode: %t, children: %d",
			reflect.Map, d.fieldType.Kind(), d.leafNode, len(d.children))
	}

	return nil
}

func verifyResourceForest(forest *ResourceForest) error {
	if len(forest.ConnectedNodes) != 4 {
		return fmt.Errorf("Invalid number of connected nodes found. Expected : 4, Found : %d", len(forest.ConnectedNodes))
	}

	baseType := reflect.TypeOf(baseType{})
	if value, found := forest.ConnectedNodes[baseType.PkgPath() + "." + baseType.Name()]; !found {
		return errors.New("Cannot find baseType{} connectedNode")
	} else if value.Len() != 3 {
			return fmt.Errorf("Invalid length of baseType{} Node. Expected : 3 Found : %d", value.Len())
	}

	arrayType := reflect.TypeOf(arrayType{})
	if value, found := forest.ConnectedNodes[arrayType.PkgPath() + "." + arrayType.Name()]; !found {
		return errors.New("Cannot find arrayType{} connectedNode")
	} else if value.Len() != 1 {
		return fmt.Errorf("Invalid length of arrayType{} Node. Expected : 1 Found : %d", value.Len())
	}

	return nil
}