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
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Size",xDescriptors="urn:alm:descriptor:io.kubernetes:size"
	Size int32 `json:"size,omitempty"`
}

// ArcadeStatus defines the observed state of Arcade
// +k8s:openapi-gen=true
type ArcadeStatus struct {
	// Indicates the status of the Arcade instance; set to "OK" when Arcade instance is up
	// +operator-sdk:csv:customresourcedefinitions:type=status,displayName="ArcadeStatus",xDescriptors="urn:alm:descriptor:com.tectonic.ui:arcadeStatus"
	ArcadeStatus string `json:"arcadeStatus,omitempty"`
	// Provides additional information about a failure status
	// +operator-sdk:csv:customresourcedefinitions:type=status,displayName="Reason",xDescriptors="urn:alm:descriptor:io.kubernetes.phase:reason"
	Reason string `json:"reason,omitempty"`
}

// Different values for ArcadeStatus
const (
	ArcadeStatusOK      string = "OK"
	ArcadeStatusFailure string = "Failure"
	ArcadeStatusPending string = "Pending"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Arcade is the Schema for the arcades API
// +operator-sdk:csv:customresourcedefinitions:displayName="Arcade Instance",resources={{Pod,v1,arcade-sample},{Deployment,v1,arcade-sample},{Service,v1,arcade-sample}}
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
