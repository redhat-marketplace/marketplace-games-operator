/*


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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ArcadeSpec defines the desired state of Arcade, through defined fields
type ArcadeSpec struct {
	// Size field used to determine total number of Arcade deployments. This field is optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Size"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:io.kubernetes:size"
	// +optional
	Size int32 `json:"size,omitempty"`
}

// ArcadeStatus defines the observed state of Arcade
// +k8s:openapi-gen=true
type ArcadeStatus struct {
	// Indicates the status of the Arcade instance; set "OK" when Arcade instance is up
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	ArcadeStatus string `json:"arcadeStatus,omitempty"`
	// Provides additional information about a failure status
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	Reason string `json:"reason,omitempty"`
}

// Different values for ArcadeStatus
const (
	ArcadeStatusOK      string = "OK"
	ArcadeStatusFailure string = "Failure"
	ArcadeStatusPending string = "Pending"
)

// +operator-sdk:gen-csv:customresourcedefinitions.displayName="Arcade"
// +operator-sdk:gen-csv:customresourcedefinitions.resources="Deployment,v1,\"A Kubernetes Deployment\""
// +operator-sdk:gen-csv:customresourcedefinitions.resources="Service,v1,\"A Kubernetes Service\""

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// Arcade is the Schema for the arcades API
type Arcade struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ArcadeSpec   `json:"spec,omitempty"`
	Status ArcadeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ArcadeList contains a list of Arcade
type ArcadeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Arcade `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Arcade{}, &ArcadeList{})
}
