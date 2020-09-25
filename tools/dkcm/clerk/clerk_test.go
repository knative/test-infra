/*
Copyright 2020 The Knative Authors

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
package clerk

import (
	"errors"
	"reflect"
	"testing"
)

var (
	fakeClusterParams  = NewClusterParams(AddZone("us-central1"), AddNodes(4), AddNodeType("e2-standard-4"))
	fakeCluster        = NewCluster(AddProjectID("knative-boskos-03"), AddStatus("WIP"))
	fakeClusterParams2 = NewClusterParams(AddNodes(4), AddNodeType("e2-standard-4"))
	fakeCluster2       = NewCluster(AddProjectID("knative-boskos-03"))
	fakeClusterParams3 = NewClusterParams(AddNodes(4))
)

// simulate a row of ClusterDB, some fields can be null
type fakeClusterScanner struct {
	ClusterID int64
	ProjectID string
	Status    string
	Zone      string
	Nodes     int64
	NodeType  string
}

func (fc fakeClusterScanner) Scan(dest ...interface{}) error {
	fcFields := []interface{}{fc.ClusterID, fc.ProjectID, fc.Status, fc.Zone, fc.Nodes, fc.NodeType}
	for idx, val := range dest {
		switch d := val.(type) {
		case *string:
			switch res := fcFields[idx].(type) {
			case string:
				*d = res
			}
		case *int64:
			switch res := fcFields[idx].(type) {
			case int64:
				*d = res
			}
		}
	}
	return nil

}

func fakeClusterSetUp() {
	fakeCluster.ClusterParams = fakeClusterParams
	fakeCluster2.ClusterParams = fakeClusterParams2
	fakeCluster.ID = 1
}

func TestPopulateCluster(t *testing.T) {
	fakeClusterSetUp()
	fc := &fakeClusterScanner{ClusterID: 1, ProjectID: "knative-boskos-03", Nodes: 4, NodeType: "e2-standard-4", Status: "WIP", Zone: "us-central1"} // full cluster
	fc2 := &fakeClusterScanner{ProjectID: "knative-boskos-03", Nodes: 4, NodeType: "e2-standard-4"}                                                  // partial cluster
	cases := []struct {
		row        *fakeClusterScanner
		wantResult *Cluster
		wantErr    error
	}{
		{fc, fakeCluster, nil},
		{fc2, fakeCluster2, nil},
	}
	for _, test := range cases {
		actualCluster, actualErr := populateCluster(test.row)
		if !reflect.DeepEqual(actualCluster, test.wantResult) {
			t.Errorf("get Cluster: got cluster '%v', want cluster '%v'", actualCluster, test.wantResult)
		}
		if !errors.Is(actualErr, test.wantErr) {
			t.Errorf("get Cluster: got err '%v', want err '%v'", actualErr, test.wantErr)
		}
	}
}

func TestGenerateAnd(t *testing.T) {
	cases := []struct {
		cp              *ClusterParams
		wantResult      string
		fieldStatements []string
	}{
		{fakeClusterParams, "Zone = 'us-central1' AND Nodes = 4 AND NodeType = 'e2-standard-4'", fakeClusterParams.generateParamsConditions(QueryZone(), QueryNodes(), QueryNodeType())},
		{fakeClusterParams2, "Zone = '' AND Nodes = 4 AND NodeType = 'e2-standard-4'", fakeClusterParams2.generateParamsConditions(QueryZone(), QueryNodes(), QueryNodeType())},
		{fakeClusterParams3, "Nodes = 4", fakeClusterParams3.generateParamsConditions(QueryNodes())},
	}
	for _, test := range cases {
		conditions := generateAND(test.fieldStatements)
		if !reflect.DeepEqual(conditions, test.wantResult) {
			t.Errorf("get condition: got condition '%s', want condition '%s'", conditions, test.wantResult)
		}
	}

}

func TestupdateQueryString(t *testing.T) {
	cases := []struct {
		dbName     string
		id         int64
		wantResult string
		update     []UpdateOption
	}{
		{"Clusters", 1, "UPDATE Clusters SET Zone = 'us-central1'WHERE ID = 1", []UpdateOption{UpdateStringField("Zone", "us-central1")}},
		{"Clusters", 1, "UPDATE Clusters SET Zone = 'us-central1',Nodes=6 WHERE ID = 1", []UpdateOption{UpdateStringField("Zone", "us-central1"), UpdateNumField("Nodes", 6)}},
	}
	for _, test := range cases {
		generatedString := updateQueryString(test.dbName, test.id, test.update...)
		if !reflect.DeepEqual(generatedString, test.wantResult) {
			t.Errorf("get string: got string '%s', want string '%s'", generatedString, test.wantResult)
		}
	}
}
