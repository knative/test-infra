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
	"container/list"
)

//ResourceForest represents the top-level forest that contains individual resource trees for top-level resource types and all connected nodes across resource trees.
type ResourceForest struct {
	Version string
	TopLevelTrees map[string]INode // Key is ResourceTree.ResourceName
	ConnectedNodes map[string]list.List // Head of the linked list keyed by nodeData.fieldType.pkg + nodeData.fieldType.Name()
}