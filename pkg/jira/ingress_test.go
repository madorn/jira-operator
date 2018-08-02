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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestNewIngressMetadata verifies that an Ingress gets created with correct metadata.
func TestNewIngressMetadata(t *testing.T) {
	j := &v1alpha1.Jira{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test-namespace",
		},
		Spec: v1alpha1.JiraSpec{
			Ingress: &v1alpha1.JiraIngressPolicy{
				Host: "test-host",
			},
		},
	}

	ing := newIngress(j)

	assert.NotNil(t, ing)
	assert.Equal(t, "test", ing.ObjectMeta.Name)
	assert.Equal(t, "test-namespace", ing.ObjectMeta.Namespace)
}

// TestNewIngressSecretMetadata verifies that an Ingress Secret gets created with correct metadata.
func TestNewIngressSecretMetadata(t *testing.T) {
	j := &v1alpha1.Jira{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-jira",
			Namespace: "test-namespace",
		},
		Spec: v1alpha1.JiraSpec{
			Ingress: &v1alpha1.JiraIngressPolicy{
				SecretName: "test-secret-name",
			},
		},
	}

	sec := newIngressSecret(j)

	assert.NotNil(t, sec)
	assert.Equal(t, "test-secret-name", sec.ObjectMeta.Name)
	assert.Equal(t, "test-namespace", sec.ObjectMeta.Namespace)
}

// TestProcessIngressError verifies an unexpected error is returned when encountered.
func TestProcessIngressError(t *testing.T) {
	j := &v1alpha1.Jira{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-jira",
			Namespace: "test-namespace",
		},
		Spec: v1alpha1.JiraSpec{
			Ingress: &v1alpha1.JiraIngressPolicy{
				Host: "test-host",
			},
		},
	}

	s := new(MockSDK)
	s.On("Get", mock.Anything).Return(errors.New("test-error"))

	err := processIngress(j, s)

	if assert.Error(t, err) {
		assert.Equal(t, errors.New("test-error"), err)
	}
	s.AssertExpectations(t)
}
