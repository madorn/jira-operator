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

	"github.com/jmckind/jira-operator/pkg/apis/jira/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// TestNewPodMetadata verifies that a Pod gets created with correct metadata.
func TestNewPodMetadata(t *testing.T) {
	j := &v1alpha1.Jira{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test-namespace",
		},
	}

	pod := newPod(j)

	assert.NotNil(t, pod)
	assert.Equal(t, "test", pod.ObjectMeta.Name)
	assert.Equal(t, "test-namespace", pod.ObjectMeta.Namespace)
}

// TestProcessPodsError verifies an unexpected error is returned when encountered.
func TestProcessPodsError(t *testing.T) {
	s := new(MockSDK)
	s.On("Get", mock.Anything).Return(errors.New("test-error"))

	err := processPods(new(v1alpha1.Jira), s)

	if assert.Error(t, err) {
		assert.Equal(t, errors.New("test-error"), err)
	}
	s.AssertExpectations(t)
}

// TestProcessPodsExists verifies a new Pod resource is not created when it already exists.
func TestProcessPodsExists(t *testing.T) {
	s := new(MockSDK)
	s.On("Get", mock.Anything).Return(nil)

	err := processPods(new(v1alpha1.Jira), s)

	assert.Nil(t, err)
	s.AssertExpectations(t)
}

// TestProcessPodsNew verifies a new Pod resource is created when it does not exist.
func TestProcessPodsNew(t *testing.T) {
	s := new(MockSDK)
	s.On("Get", mock.Anything).Return(apierrors.NewNotFound(schema.GroupResource{}, "test"))
	s.On("Create", mock.Anything).Return(nil)

	err := processPods(new(v1alpha1.Jira), s)

	assert.Nil(t, err)
	s.AssertExpectations(t)
}
