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
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type FailedConfigMapSDK struct {
	MockSDK
}

func (m *FailedConfigMapSDK) Get(o runtime.Object) error {
	switch obj := o.(type) {
	case *v1.ConfigMap:
		_ = m.Called(obj)
		return errors.New("failed get configmap")
	}
	return nil
}

type FailedPVCSDK struct {
	MockSDK
}

func (m *FailedPVCSDK) Get(o runtime.Object) error {
	switch obj := o.(type) {
	case *v1.PersistentVolumeClaim:
		_ = m.Called(obj)
		return errors.New("failed get pvc")
	}
	return nil
}

func TestNewReconcilerHandlesNilObject(t *testing.T) {
	r := NewReconciler(nil)
	assert.Nil(t, r)
}

func TestReconcileHandlesNilObject(t *testing.T) {
	r := NewReconciler(new(MockSDK))
	err := r.Reconcile(nil)
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("jira reference cannot be nil"), err)
	}
}

func TestReconcileSetDefaultsChanged(t *testing.T) {
	sdk := new(MockSDK)
	sdk.On("Update", mock.Anything).Return(nil)

	r := NewReconciler(sdk)
	err := r.Reconcile(new(v1alpha1.Jira))

	assert.Nil(t, err)
	sdk.AssertExpectations(t)
}

func TestReconcileSetDefaultsNotChanged(t *testing.T) {
	jira := new(v1alpha1.Jira)
	jira.Spec.BaseImage = "test"
	jira.Spec.BaseImageVersion = "test"
	jira.Spec.ConfigMapName = "test"
	jira.Spec.DataMountPath = "test"
	jira.Spec.SecretName = "test"

	sdk := new(MockSDK)
	sdk.On("Get", mock.Anything).Return(nil)

	r := NewReconciler(sdk)
	err := r.Reconcile(jira)

	assert.Nil(t, err)
	sdk.AssertExpectations(t)
}

func TestReconcileHandleConfigMapError(t *testing.T) {
	jira := new(v1alpha1.Jira)
	jira.Spec.BaseImage = "test"
	jira.Spec.BaseImageVersion = "test"
	jira.Spec.ConfigMapName = "test"
	jira.Spec.DataMountPath = "test"
	jira.Spec.SecretName = "test"

	sdk := new(FailedConfigMapSDK)
	sdk.On("Get", mock.Anything)

	r := NewReconciler(sdk)
	err := r.Reconcile(jira)

	if assert.Error(t, err) {
		assert.Equal(t, errors.New("failed get configmap"), err)
	}
	sdk.AssertExpectations(t)
}

func TestReconcileHandlePVCError(t *testing.T) {
	scName := "test"
	jira := new(v1alpha1.Jira)
	jira.Spec.BaseImage = "test"
	jira.Spec.BaseImageVersion = "test"
	jira.Spec.ConfigMapName = "test"
	jira.Spec.DataMountPath = "test"
	jira.Spec.SecretName = "test"
	jira.Spec.Pod = &v1alpha1.JiraPodPolicy{
		PersistentVolumeClaimSpec: &v1.PersistentVolumeClaimSpec{
			StorageClassName: &scName,
		},
	}

	sdk := new(FailedPVCSDK)
	sdk.On("Get", mock.Anything).Return(nil)

	r := NewReconciler(sdk)
	err := r.Reconcile(jira)

	if assert.Error(t, err) {
		assert.Equal(t, errors.New("failed get pvc"), err)
	}
	sdk.AssertExpectations(t)
}
