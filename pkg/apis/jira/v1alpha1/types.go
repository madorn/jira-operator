// Copyright 2018 Jira Operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// DefaultBaseImage is the default docker image to use for JIRA Pods.
	DefaultBaseImage = "cptactionhank/atlassian-jira"
	// DefaultBaseImageVersion is the default version to use for JIRA Pods.
	DefaultBaseImageVersion = "7.10.2"
	// DefaultDataMountPath is the default filesystem path for JIRA Home.
	DefaultDataMountPath = "/var/atlassian/jira"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// JiraList resource
type JiraList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Jira `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Jira resource
type Jira struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              JiraSpec   `json:"spec"`
	Status            JiraStatus `json:"status,omitempty"`
}

// JiraPodPolicy defines the policy for pods owned by rethinkdb operator.
type JiraPodPolicy struct {
	// Resources is the resource requirements for the jira container.
	Resources v1.ResourceRequirements `json:"resources,omitempty"`

	// PersistentVolumeClaimSpec is the spec to describe PVC for the jira container
	// This field is optional. If no PVC spec, jira container will use emptyDir as volume
	PersistentVolumeClaimSpec *v1.PersistentVolumeClaimSpec `json:"persistentVolumeClaimSpec,omitempty"`
}

// JiraSpec resource
type JiraSpec struct {
	// BaseImage image to use for a RethinkDB deployment.
	BaseImage string `json:"base_image"`

	// BaseImageVersion is the version of base image to use.
	BaseImageVersion string `json:"base_image_version"`

	// DataMountPath path for JIRA Home.
	DataMountPath string `json:"data_mount_path"`

	// ConfigMapName is the name of ConfigMap to use or create.
	ConfigMapName string `json:"configMapName"`

	// SecretName is the name of Secret to use or create.
	SecretName string `json:"secretName"`

	// Pod defines the policy for pods owned by rethinkdb operator.
	// This field cannot be updated once the CR is created.
	Pod *JiraPodPolicy `json:"pod,omitempty"`
}

// SetDefaults sets the default vaules for the cuberite spec and returns true if the spec was changed
func (j *Jira) SetDefaults() bool {
	changed := false
	if len(j.Spec.BaseImage) == 0 {
		j.Spec.BaseImage = DefaultBaseImage
		changed = true
	}
	if len(j.Spec.BaseImageVersion) == 0 {
		j.Spec.BaseImageVersion = DefaultBaseImageVersion
		changed = true
	}
	if len(j.Spec.ConfigMapName) == 0 {
		j.Spec.ConfigMapName = j.Name
		changed = true
	}
	if len(j.Spec.DataMountPath) == 0 {
		j.Spec.DataMountPath = DefaultDataMountPath
		changed = true
	}
	if len(j.Spec.SecretName) == 0 {
		j.Spec.SecretName = j.Name
		changed = true
	}
	return changed
}

// IsPVEnabled shortcut fucntion to determine PV status.
func (j *Jira) IsPVEnabled() bool {
	if podPolicy := j.Spec.Pod; podPolicy != nil {
		return podPolicy.PersistentVolumeClaimSpec != nil
	}
	return false
}

// JiraStatus resource
type JiraStatus struct {
	// Fill me
}
