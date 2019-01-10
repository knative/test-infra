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
	"container/list"
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

func getBaseTypeValue() baseType {
	return baseType{
		field1: "test",
	}
}

func getPtrTypeValueAllCovered() ptrType {
	f := new(float32)
	*f = 3.142
	b := getBaseTypeValue()
	return ptrType {
		basePtr: f,
		structPtr: &b,
	}
}

func getPtrTypeValueSomeCovered() ptrType {
	b := getBaseTypeValue()
	return ptrType {
		structPtr: &b,
	}
}

func getArrValueAllCovered() arrayType {
	b1 := getBaseTypeValue()
	b2 := baseType{
		field2: 32,
	}

	return arrayType{
		structArr: []baseType{b1, b2},
		baseArr: []bool{true, false},
	}
}

func getArrValueSomeCovered() arrayType {
	return arrayType{
		structArr: []baseType{getBaseTypeValue()},
	}
}

func getOtherTypeValue() otherType {
	m := make(map[string]baseType)
	m["test"] = getBaseTypeValue()
	return otherType {
		structMap: m,
	}
}

func getTestTree(treeName string, t reflect.Type) *ResourceTree {
	forest := ResourceForest{
		Version: "TestVersion",
		ConnectedNodes: make(map[string]*list.List),
		TopLevelTrees: make(map[string]ResourceTree),
	}

	tree := ResourceTree{
		ResourceName: treeName,
		Forest: &forest,
	}

	tree.BuildResourceTree(t)
	forest.TopLevelTrees[treeName] = tree
	return &tree
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

func verifyBaseTypeValue(logPrefix string, node NodeInterface) error {
	if !node.getData().covered {
		return errors.New(logPrefix + " Node marked as not-covered. Expected to be covered")
	}

	if !node.getData().children["field1"].getData().covered {
		return errors.New(logPrefix + " field1 marked as not-covered. Expected to be covered")
	}

	if node.getData().children["field2"].getData().covered {
		return errors.New(logPrefix + " field2 marked as covered. Expected to be not-covered")
	}

	return nil
}

func verifyPtrValueAllCovered(node NodeInterface) error {
	if !node.getData().covered {
		return errors.New("Node marked as not-covered. Expected to be covered")
	}

	child := node.getData().children["basePtr"]
	if !child.getData().covered {
		return errors.New("field:base_ptr marked as not-covered. Expected to be covered")
	}

	if !child.getData().children["basePtr" + ptrNodeNameSuffix].getData().covered {
		return errors.New("field:basePtr" + ptrNodeNameSuffix + "marked as not-covered. Expected to be covered" )
	}

	child = node.getData().children["structPtr"]
	if !child.getData().covered {
		return errors.New("field:structPtr marked as not-covered. Expected to be covered")
	}

	if err := verifyBaseTypeValue("field:structPtr" + ptrNodeNameSuffix, child.getData().children["structPtr" + ptrNodeNameSuffix]); err != nil {
		return err
	}

	return nil
}

func verifyPtrValueSomeCovered(node NodeInterface) error {
	if !node.getData().covered {
		return errors.New("Node marked as not-covered. Expected to be covered")
	}

	child := node.getData().children["basePtr"]
	if child.getData().covered {
		return errors.New("field:basePtr marked as covered. Expected to be not-covered")
	}

	if child.getData().children["basePtr" + ptrNodeNameSuffix].getData().covered {
		return errors.New("field:basePtr" + ptrNodeNameSuffix + "marked as covered. Expected to be not-covered" )
	}

	child = node.getData().children["structPtr"]
	if !child.getData().covered {
		return errors.New("field:structPtr marked as not-covered. Expected to be covered")
	}

	if err := verifyBaseTypeValue("field:structPtr" + ptrNodeNameSuffix, child.getData().children["structPtr" + ptrNodeNameSuffix]); err != nil {
		return err
	}

	return nil
}

func verifyArryValueAllCovered(node NodeInterface) error {
	if !node.getData().covered {
		return errors.New("Node marked as not-covered. Expected to be covered")
	}

	child := node.getData().children["baseArr"]
	if !child.getData().covered {
		return errors.New("field:baseArr marked as not-covered. Expected to be covered")
	}

	if !child.getData().children["baseArr" + arrayNodeNameSuffix].getData().covered {
		return errors.New("field:baseArr" + arrayNodeNameSuffix + " marked as not-covered. Expected to be covered" )
	}

	child = node.getData().children["structArr"]
	if !child.getData().covered {
		return errors.New("field:structArr marked as not-covered. Expected to be covered")
	}

	child = child.getData().children["structArr" + arrayNodeNameSuffix]
	if !child.getData().covered {
		return errors.New("structArr" + arrayNodeNameSuffix + " marked as not-covered. Expected to be covered")
	}

	if !child.getData().children["field1"].getData().covered {
		return errors.New("structArr" + arrayNodeNameSuffix + ".field1 marked as not-covered. Expected to be covered")
	}

	if !child.getData().children["field2"].getData().covered {
		return errors.New("structArr" + arrayNodeNameSuffix +".field1 marked as not-covered. Expected to be covered")
	}

	return nil
}

func verifyArrValueSomeCovered(node NodeInterface) error {
	if !node.getData().covered {
		return errors.New("Node marked as not-covered. Expected to be covered")
	}

	child := node.getData().children["baseArr"]
	if child.getData().covered {
		return errors.New("field:baseArr marked as covered. Expected to be not-covered")
	}

	if child.getData().children["baseArr" + arrayNodeNameSuffix].getData().covered {
		return errors.New("field:baseArr" + arrayNodeNameSuffix + " marked as covered. Expected to be not-covered" )
	}

	child = node.getData().children["structArr"]
	if !child.getData().covered {
		return errors.New("field:structArr marked as not-covered. Expected to be covered")
	}

	if err := verifyBaseTypeValue("field:structArr" + arrayNodeNameSuffix, child.getData().children["structArr" + arrayNodeNameSuffix]); err != nil {
		return err
	}

	return nil
}

func verifyOtherTypeValue(node NodeInterface) error {
	if !node.getData().covered {
		return errors.New("Node marked as not-covered. Expected to be covered")
	}

	if !node.getData().children["structMap"].getData().covered {
		return errors.New("field:structMap marked as not-covered. Expected to be covered")
	}

	if node.getData().children["baseMap"].getData().covered {
		return errors.New("field:baseMap marked as covered. Expected to be not-covered")
	}

	return nil
}