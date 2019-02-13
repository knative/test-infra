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

package webhook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"

	"github.com/knative/pkg/webhook"
	"github.com/knative/test-infra/tools/webhook-apicoverage/coveragecalculator"
	"github.com/knative/test-infra/tools/webhook-apicoverage/resourcetree"
	"github.com/knative/test-infra/tools/webhook-apicoverage/view"
	"go.uber.org/zap"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	decoder = serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer()
)

const (
	// ResourceQueryParam query param name to provide the resource.
	ResourceQueryParam = "resource"
)

// APICoverageRecorder type contains resource tree to record API coverage for resources.
type APICoverageRecorder struct {
	Logger *zap.SugaredLogger
	ResourceForest resourcetree.ResourceForest
	ResourceMap map[schema.GroupVersionKind]webhook.GenericCRD
	NodeRules resourcetree.NodeRules
	FieldRules resourcetree.FieldRules
	DisplayRules view.DisplayRules
}

// Init initializes the resources trees for set resources.
func (a *APICoverageRecorder) Init() {
	for resourceKind, resourceObj := range a.ResourceMap {
		a.ResourceForest.AddResourceTree(resourceKind.Kind, reflect.ValueOf(resourceObj).Elem().Type())
	}
}

// RecordResourceCoverage updates the resource tree with the request.
func (a *APICoverageRecorder) RecordResourceCoverage(w http.ResponseWriter, r *http.Request) {
	var (
		body []byte
		err error
	)

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		a.Logger.Errorf("Failed reading request body: %v", err)
		a.writeAdmissionResponse(a.getAdmissionResponse(false, "Admission Denied"), w)
		return
	}

	review := &v1beta1.AdmissionReview{}
	if _, _, err := decoder.Decode(body, nil, review); err != nil {
		a.Logger.Errorf("Unable to decode request: %v", err)
		a.writeAdmissionResponse(a.getAdmissionResponse(false, "Admission Denied"), w)
		return
	}

	gvk := schema.GroupVersionKind {
		Group:   review.Request.Kind.Group,
		Version: review.Request.Kind.Version,
		Kind:    review.Request.Kind.Kind,
	}
	if err := json.Unmarshal(review.Request.Object.Raw, a.ResourceMap[gvk]); err != nil {
		a.Logger.Errorf("Failed unmarshalling review.Request.Object.Raw for type: %s Error: %v", a.ResourceMap[gvk], err)
		a.writeAdmissionResponse(a.getAdmissionResponse(false, "Admission Denied"), w)
		return
	}
	resourceTree := a.ResourceForest.TopLevelTrees[gvk.Kind]
	resourceTree.UpdateCoverage(reflect.ValueOf(a.ResourceMap[gvk]).Elem())
	a.writeAdmissionResponse(a.getAdmissionResponse(true, "Welcome Aboard"), w)
}

func (a *APICoverageRecorder) getAdmissionResponse(allowed bool, message string) (*v1beta1.AdmissionResponse) {
	return &v1beta1.AdmissionResponse{
		Allowed: allowed,
		Result: &v1.Status{
			Message: message,
		},
	}
}

func (a *APICoverageRecorder) writeAdmissionResponse(admissionResp *v1beta1.AdmissionResponse, w http.ResponseWriter) {
	responseInBytes, err := json.Marshal(admissionResp)
	if err != nil {
		a.Logger.Errorf("Failing mashalling review response: %v", err)
	}

	if _, err := w.Write(responseInBytes); err != nil {
		a.Logger.Errorf("%v", err)
	}
}

// GetResourceCoverage retrieves resource coverage data for the passed in resource via query param.
func (a *APICoverageRecorder) GetResourceCoverage(w http.ResponseWriter, r *http.Request) {
	resource := r.URL.Query().Get(ResourceQueryParam)
	if _, ok := a.ResourceForest.TopLevelTrees[resource]; !ok {
		fmt.Fprintf(w, "Resource information not found for resource: %s", resource)
		return
	}

	var ignoredFields coveragecalculator.IgnoredFields
	ignoredFieldsFilePath := os.Getenv("KO_DATA_PATH") + "/ignoredfields.yaml"
	if err := ignoredFields.ReadFromFile(ignoredFieldsFilePath); err != nil {
		fmt.Fprintf(w, "Error reading file: %s", ignoredFieldsFilePath)
	}

	tree := a.ResourceForest.TopLevelTrees[resource]
	typeCoverage := tree.BuildCoverageData(a.NodeRules, a.FieldRules, ignoredFields)

	jsonLikeDisplay := view.GetJSONTypeDisplay(typeCoverage, a.DisplayRules)
	fmt.Fprint(w, jsonLikeDisplay)
}