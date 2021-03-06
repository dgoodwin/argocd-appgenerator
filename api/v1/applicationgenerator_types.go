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

package v1

import (
	argoapi "github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ApplicationGeneratorSpec defines the desired state of ApplicationGenerator
type ApplicationGeneratorSpec struct {
	// ClusterSelector is meant to match ArgoCD clusters dynamically. Note that presently ArgoCD represents clusters
	// as a Secret with the label "argocd.argoproj.io/secret-type=cluster". Labels cannot be specified on this secret
	// with the argocd CLI, so this must be done out of band. You can however specify the ArgoCD secret-type label
	// itself as your selector to match all ArgoCD Clusters.
	ClusterSelector metav1.LabelSelector `json:"clusterSelector,omitempty"`

	// ApplicationSpec is an embedded ArgoCD ApplicationSpec definition. Any Clusters which match this AppGenerator will
	// have an ArgoCD Application resource created for them.
	ApplicationSpec argoapi.ApplicationSpec `json:"applicationSpec"`
}

// ApplicationGeneratorStatus defines the observed state of ApplicationGenerator
type ApplicationGeneratorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// ApplicationGenerator is the Schema for the applicationgenerators API
type ApplicationGenerator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationGeneratorSpec   `json:"spec,omitempty"`
	Status ApplicationGeneratorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ApplicationGeneratorList contains a list of ApplicationGenerator
type ApplicationGeneratorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApplicationGenerator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApplicationGenerator{}, &ApplicationGeneratorList{})
}
