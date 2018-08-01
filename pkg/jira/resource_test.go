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
	"fmt"
	"testing"

	"github.com/jmckind/jira-operator/pkg/apis/jira/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestDefaultLabels verifies that the default set of labels are created.
func TestDefaultLabels(t *testing.T) {
	j := &v1alpha1.Jira{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}
	labels := defaultLabels(j)

	assert.NotNil(t, labels)
	assert.Equal(t, "jira", labels["app"])
	assert.Equal(t, "test", labels["cluster"])
}

// TestOwnerRef verifies that the proper owner reference is returned.
func TestOwnerRef(t *testing.T) {
	expectedAPIVersion := fmt.Sprintf(
		"%s/%s",
		v1alpha1.SchemeGroupVersion.Group,
		v1alpha1.SchemeGroupVersion.Version,
	)

	o := ownerRef(&v1alpha1.Jira{})

	assert.NotNil(t, o)
	assert.Len(t, o, 1)
	assert.Equal(t, expectedAPIVersion, o[0].APIVersion)
	assert.Equal(t, "Jira", o[0].Kind)
}

// TestResourceLabels verifies that the correct set of labels are applied to the Jira resource.
func TestResourceLabels(t *testing.T) {
	j := &v1alpha1.Jira{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-namespace",
			Labels: map[string]string{
				"test1": "value1",
				"test2": "value2",
			},
		},
	}

	labels := resourceLabels(j)

	assert.NotNil(t, labels)
	assert.Equal(t, "value1", labels["test1"])
	assert.Equal(t, "value2", labels["test2"])
}
