/*
Copyright 2022.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	FluxAddonLabelName = "tochangethis"

	FluxAddonInstanceFinalizer = "fluxinstance.addons.cluster.x-k8s.io"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FluxAddonInstanceSpec defines the desired state of FluxAddonInstance
type FluxAddonInstanceSpec struct {
	ClusterRef corev1.ObjectReference `json:"clusterRef"`
	RepoName   string                 `json:"repoName"`
}

// FluxAddonInstanceStatus defines the observed state of FluxAddonInstance
type FluxAddonInstanceStatus struct {
	// Conditions defines current state of the HelmReleaseProxy.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`

	// Status is the current status of the Helm release.
	// +optional
	Status string `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// FluxAddonInstance is the Schema for the fluxaddoninstances API
type FluxAddonInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FluxAddonInstanceSpec   `json:"spec,omitempty"`
	Status FluxAddonInstanceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FluxAddonInstanceList contains a list of FluxAddonInstance
type FluxAddonInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FluxAddonInstance `json:"items"`
}

// GetConditions returns the list of conditions for an FluxAddonInstance API object.
func (r *FluxAddonInstance) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions will set the given conditions on an FluxAddonInstance object.
func (r *FluxAddonInstance) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&FluxAddonInstance{}, &FluxAddonInstanceList{})
}
