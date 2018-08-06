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

package jira

import (
	"github.com/coreos/jira-operator/pkg/apis/jira/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// defaultLabels returns the default set of labels based on the given name.
func defaultLabels(j *v1alpha1.Jira) map[string]string {
	return map[string]string{
		"app":     "jira",
		"cluster": j.Name,
	}
}

// ownerRef returns an owner reference for the Jira resource.
func ownerRef(j *v1alpha1.Jira) []metav1.OwnerReference {
	return []metav1.OwnerReference{
		*metav1.NewControllerRef(j, schema.GroupVersionKind{
			Group:   v1alpha1.SchemeGroupVersion.Group,
			Version: v1alpha1.SchemeGroupVersion.Version,
			Kind:    "Jira",
		}),
	}
}

// resourceLabels returns the labels for the resource.
func resourceLabels(j *v1alpha1.Jira) map[string]string {
	labels := defaultLabels(j)
	for key, val := range j.ObjectMeta.Labels {
		labels[key] = val
	}
	return labels
}
