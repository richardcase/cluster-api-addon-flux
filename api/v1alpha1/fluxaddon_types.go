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
	FluxAddonFinalizer = "flux.addons.cluster.x-k8s.io"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FluxAddonSpec defines the desired state of FluxAddon
type FluxAddonSpec struct {
	// ClusterSelector selects Clusters in the same namespace with a label that matches the specified label selector. Flux
	// will be installed on all selected Clusters. If a Cluster is no longer selected, then Flux will be uninstalled.
	ClusterSelector metav1.LabelSelector `json:"clusterSelector"`

	RepositoryName string `json:"repositoryName,omitempty"`
}

// FluxAddonStatus defines the observed state of FluxAddon
type FluxAddonStatus struct {
	// Conditions defines current state of the HelmChartProxy.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`

	// MatchingClusters is the list of references to Clusters selected by the ClusterSelector.
	// +optional
	MatchingClusters []corev1.ObjectReference `json:"matchingClusters"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// FluxAddon is the Schema for the fluxaddons API
type FluxAddon struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FluxAddonSpec   `json:"spec,omitempty"`
	Status FluxAddonStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FluxAddonList contains a list of FluxAddon
type FluxAddonList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FluxAddon `json:"items"`
}

// GetConditions returns the list of conditions for an FluxAddon API object.
func (c *FluxAddon) GetConditions() clusterv1.Conditions {
	return c.Status.Conditions
}

// SetConditions will set the given conditions on an FluxAddon object.
func (c *FluxAddon) SetConditions(conditions clusterv1.Conditions) {
	c.Status.Conditions = conditions
}

func (c *FluxAddon) SetMatchingClusters(clusterList []clusterv1.Cluster) {
	matchingClusters := make([]corev1.ObjectReference, 0, len(clusterList))
	for _, cluster := range clusterList {
		matchingClusters = append(matchingClusters, corev1.ObjectReference{
			Kind:       cluster.Kind,
			APIVersion: cluster.APIVersion,
			Name:       cluster.Name,
			Namespace:  cluster.Namespace,
		})
	}
}

func init() {
	SchemeBuilder.Register(&FluxAddon{}, &FluxAddonList{})
}
