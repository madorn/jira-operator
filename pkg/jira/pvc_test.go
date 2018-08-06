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
	"errors"
	"testing"

	"github.com/coreos/jira-operator/pkg/apis/jira/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// TestNewPVCMetadata verifies that a PVC gets created with correct metadata.
func TestNewPVCMetadata(t *testing.T) {
	j := &v1alpha1.Jira{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test-namespace",
		},
		Spec: v1alpha1.JiraSpec{
			Pod: &v1alpha1.JiraPodPolicy{
				PersistentVolumeClaimSpec: &v1.PersistentVolumeClaimSpec{},
			},
		},
	}

	pvc := newPVC(j)

	assert.NotNil(t, pvc)
	assert.Equal(t, "test", pvc.ObjectMeta.Name)
	assert.Equal(t, "test-namespace", pvc.ObjectMeta.Namespace)
}

// TestProcessPVCError verifies an unexpected error is returned when encountered.
func TestProcessPVCError(t *testing.T) {
	j := new(v1alpha1.Jira)
	j.Spec.Pod = &v1alpha1.JiraPodPolicy{
		PersistentVolumeClaimSpec: &v1.PersistentVolumeClaimSpec{},
	}

	s := new(MockSDK)
	s.On("Get", mock.Anything).Return(errors.New("test-error"))

	err := processPVC(j, s)

	if assert.Error(t, err) {
		assert.Equal(t, errors.New("test-error"), err)
	}
	s.AssertExpectations(t)
}

// TestProcessPVCExists verifies a new PVC resource is not created when it already exists.
func TestProcessPVCExists(t *testing.T) {
	j := new(v1alpha1.Jira)
	j.Spec.Pod = &v1alpha1.JiraPodPolicy{
		PersistentVolumeClaimSpec: &v1.PersistentVolumeClaimSpec{},
	}

	s := new(MockSDK)
	s.On("Get", mock.Anything).Return(nil)

	err := processPVC(j, s)

	assert.Nil(t, err)
	s.AssertExpectations(t)
}

// TestProcessPVCNew verifies a new PVC resource is created when it does not exist.
func TestProcessPVCNew(t *testing.T) {
	j := new(v1alpha1.Jira)
	j.Spec.Pod = &v1alpha1.JiraPodPolicy{
		PersistentVolumeClaimSpec: &v1.PersistentVolumeClaimSpec{},
	}

	s := new(MockSDK)
	s.On("Get", mock.Anything).Return(apierrors.NewNotFound(schema.GroupResource{}, "test"))
	s.On("Create", mock.Anything).Return(nil)

	err := processPVC(j, s)

	assert.Nil(t, err)
	s.AssertExpectations(t)
}
