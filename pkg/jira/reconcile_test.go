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
	extensions "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

// FailedConfigMapSDK is used to mock behavior for the ConfigMap resource.
type FailedConfigMapSDK struct {
	MockSDK
}

// Get will return an error only when given a ConfigMap resource.
func (m *FailedConfigMapSDK) Get(o runtime.Object) error {
	switch obj := o.(type) {
	case *v1.ConfigMap:
		_ = m.Called(obj)
		return errors.New("failed get configmap")
	}
	return nil
}

// FailedIngressSDK is used to mock behavior for the Ingress resource.
type FailedIngressSDK struct {
	MockSDK
}

// Get will return an error only when given an Ingress resource.
func (m *FailedIngressSDK) Get(o runtime.Object) error {
	switch obj := o.(type) {
	case *extensions.Ingress:
		_ = m.Called(obj)
		return errors.New("failed get ingress")
	}
	return nil
}

// FailedPodSDK is used to mock behavior for the Pod resources.
type FailedPodSDK struct {
	MockSDK
}

// Get will return an error only when given a Pod resource.
func (m *FailedPodSDK) Get(o runtime.Object) error {
	switch obj := o.(type) {
	case *v1.Pod:
		_ = m.Called(obj)
		return errors.New("failed get pod")
	}
	return nil
}

// FailedPVCSDK is used to mock behavior for the PersistentVolumeClaim resource.
type FailedPVCSDK struct {
	MockSDK
}

// Get will return an error only when given a PersistentVolumeClaim resource.
func (m *FailedPVCSDK) Get(o runtime.Object) error {
	switch obj := o.(type) {
	case *v1.PersistentVolumeClaim:
		_ = m.Called(obj)
		return errors.New("failed get pvc")
	}
	return nil
}

// FailedSecretSDK is used to mock behavior for the Secret resource.
type FailedSecretSDK struct {
	MockSDK
}

// Get will return an error only when given a Secret resource.
func (m *FailedSecretSDK) Get(o runtime.Object) error {
	switch obj := o.(type) {
	case *v1.Secret:
		_ = m.Called(obj)
		return errors.New("failed get secret")
	}
	return nil
}

// FailedServiceSDK is used to mock behavior for the Service resource.
type FailedServiceSDK struct {
	MockSDK
}

// Get will return an error only when given a Service resource.
func (m *FailedServiceSDK) Get(o runtime.Object) error {
	switch obj := o.(type) {
	case *v1.Service:
		_ = m.Called(obj)
		return errors.New("failed get service")
	}
	return nil
}

// TestNewReconcilerHandlesNilObject verifies that attempting to create a new
// Reconciler with a nil SDK value produces a nil result.
func TestNewReconcilerHandlesNilObject(t *testing.T) {
	r := NewReconciler(nil)
	assert.Nil(t, r)
}

// TestReconcileSetDefaultsChanged verifies that an update occurs for the Jira CRD when the defaults are initially set.
func TestReconcileSetDefaultsNotChanged(t *testing.T) {
	scName := "test-storage-class-name"
	jira := new(v1alpha1.Jira)
	jira.Spec.BaseImage = "test-base-image"
	jira.Spec.BaseImageVersion = "test-base-image-version"
	jira.Spec.ConfigMapName = "test-configmap-name"
	jira.Spec.DataMountPath = "test-data-mount-path"
	jira.Spec.Ingress = &v1alpha1.JiraIngressPolicy{
		Host:       "test-ingress-host",
		Path:       "/",
		TLS:        true,
		SecretName: "test-ingress-secret-name",
	}
	jira.Spec.Pod = &v1alpha1.JiraPodPolicy{
		PersistentVolumeClaimSpec: &v1.PersistentVolumeClaimSpec{
			StorageClassName: &scName,
		},
	}
	jira.Spec.SecretName = "test-secret-name"

	sdk := new(MockSDK)
	sdk.On("Get", mock.Anything).Return(nil)

	r := NewReconciler(sdk)
	err := r.Reconcile(jira)

	assert.Nil(t, err)
	sdk.AssertExpectations(t)
	sdk.AssertNumberOfCalls(t, "Get", 6)
}

// TestReconcileHandlesNilObject verifies that a nil Jira value produces an error.
func TestReconcileHandlesNilObject(t *testing.T) {
	r := NewReconciler(new(MockSDK))
	err := r.Reconcile(nil)
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("jira reference cannot be nil"), err)
	}
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

func TestReconcileHandleIngressError(t *testing.T) {
	jira := new(v1alpha1.Jira)
	jira.Spec.BaseImage = "test"
	jira.Spec.BaseImageVersion = "test"
	jira.Spec.ConfigMapName = "test"
	jira.Spec.DataMountPath = "test"
	jira.Spec.SecretName = "test"
	jira.Spec.Ingress = &v1alpha1.JiraIngressPolicy{
		Host:       "test-host",
		Path:       "/",
		TLS:        true,
		SecretName: "test-secret-name",
	}

	sdk := new(FailedIngressSDK)
	sdk.On("Get", mock.Anything).Return(nil)

	r := NewReconciler(sdk)
	err := r.Reconcile(jira)

	if assert.Error(t, err) {
		assert.Equal(t, errors.New("failed get ingress"), err)
	}
	sdk.AssertExpectations(t)
}

func TestReconcileHandleIngressSecretError(t *testing.T) {
	jira := new(v1alpha1.Jira)
	jira.Spec.BaseImage = "test"
	jira.Spec.BaseImageVersion = "test"
	jira.Spec.ConfigMapName = "test"
	jira.Spec.DataMountPath = "test"
	jira.Spec.SecretName = "test"
	jira.Spec.Ingress = &v1alpha1.JiraIngressPolicy{
		Host:       "test-host",
		Path:       "/",
		TLS:        true,
		SecretName: "test-secret-name",
	}

	sdk := new(FailedSecretSDK)
	sdk.On("Get", mock.Anything).Return(nil)

	r := NewReconciler(sdk)
	err := r.Reconcile(jira)

	if assert.Error(t, err) {
		assert.Equal(t, errors.New("failed get secret"), err)
	}
	sdk.AssertExpectations(t)
}

func TestReconcileHandlePodError(t *testing.T) {
	jira := new(v1alpha1.Jira)
	jira.Spec.BaseImage = "test"
	jira.Spec.BaseImageVersion = "test"
	jira.Spec.ConfigMapName = "test"
	jira.Spec.DataMountPath = "test"
	jira.Spec.SecretName = "test"

	sdk := new(FailedPodSDK)
	sdk.On("Get", mock.Anything).Return(nil)

	r := NewReconciler(sdk)
	err := r.Reconcile(jira)

	if assert.Error(t, err) {
		assert.Equal(t, errors.New("failed get pod"), err)
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

func TestReconcileHandleServiceError(t *testing.T) {
	jira := new(v1alpha1.Jira)
	jira.Spec.BaseImage = "test"
	jira.Spec.BaseImageVersion = "test"
	jira.Spec.ConfigMapName = "test"
	jira.Spec.DataMountPath = "test"
	jira.Spec.SecretName = "test"

	sdk := new(FailedServiceSDK)
	sdk.On("Get", mock.Anything).Return(nil)

	r := NewReconciler(sdk)
	err := r.Reconcile(jira)

	if assert.Error(t, err) {
		assert.Equal(t, errors.New("failed get service"), err)
	}
	sdk.AssertExpectations(t)
}

// TestReconcileSetDefaultsChanged verifies that an update occurs for the Jira CRD when the defaults are initially set.
func TestReconcileSetDefaultsChanged(t *testing.T) {
	sdk := new(MockSDK)
	sdk.On("Update", mock.Anything).Return(nil)

	r := NewReconciler(sdk)
	err := r.Reconcile(new(v1alpha1.Jira))

	assert.Nil(t, err)
	sdk.AssertExpectations(t)
}
