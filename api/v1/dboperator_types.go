/*
Copyright 2026.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DboperatorSpec defines the desired state of Dboperator
type DboperatorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Credentials of Dboperator. Edit dboperator_types.go to remove/update

	DbType   string `json:"dbtype"`
	Username string `json:"username"`
	Password string `json:"password"`
	DbName   string `json:"dbname,omitempty"`

	// +kubebuilder:validation:Enum=small;medium;big
	// +kubebuilder:default:=small
	Size    string `json:"size,omitempty"`
	Version string `json:"version,omitempty"`
}

// DboperatorStatus defines the observed state of Dboperator
type DboperatorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Phase      string             `json:"phase,omitempty"`
	Hostname   string             `json:"hostname,omitempty"`
	Port       string             `json:"port,omitempty"`
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Dboperator is the Schema for the dboperators API
type Dboperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DboperatorSpec   `json:"spec,omitempty"`
	Status DboperatorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DboperatorList contains a list of Dboperator
type DboperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Dboperator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Dboperator{}, &DboperatorList{})
}
