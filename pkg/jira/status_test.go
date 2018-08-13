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
	"testing"

	"github.com/coreos/jira-operator/pkg/apis/jira/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestFormatEndpointWithoutIngress(t *testing.T) {
	j := &v1alpha1.Jira{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-jira",
			Namespace: "test-jira-namespace",
		},
	}

	e := formatEndpoint(j)

	assert.NotNil(t, e)
	assert.Equal(t, "http://test-jira:8080/", e)
}

func TestFormatEndpointIngressWithoutTLS(t *testing.T) {
	j := &v1alpha1.Jira{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-jira",
			Namespace: "test-jira-namespace",
		},
		Spec: v1alpha1.JiraSpec{
			Ingress: &v1alpha1.JiraIngressPolicy{
				Host: "test-ingress-host",
			},
		},
	}

	e := formatEndpoint(j)

	assert.NotNil(t, e)
	assert.Equal(t, "http://test-ingress-host", e)
}

func TestFormatEndpointIngressTLS(t *testing.T) {
	j := &v1alpha1.Jira{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-jira",
			Namespace: "test-jira-namespace",
		},
		Spec: v1alpha1.JiraSpec{
			Ingress: &v1alpha1.JiraIngressPolicy{
				Host: "test-ingress-host",
				TLS:  true,
			},
		},
	}

	e := formatEndpoint(j)

	assert.NotNil(t, e)
	assert.Equal(t, "https://test-ingress-host", e)
}

func TestFormatEndpointIngressPath(t *testing.T) {
	j := &v1alpha1.Jira{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-jira",
			Namespace: "test-jira-namespace",
		},
		Spec: v1alpha1.JiraSpec{
			Ingress: &v1alpha1.JiraIngressPolicy{
				Host: "test-ingress-host",
				Path: "/test-path",
			},
		},
	}

	e := formatEndpoint(j)

	assert.NotNil(t, e)
	assert.Equal(t, "http://test-ingress-host/test-path", e)
}

func TestFormatEndpointIngressWithDefaults(t *testing.T) {
	j := &v1alpha1.Jira{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-jira",
			Namespace: "test-jira-namespace",
		},
		Spec: v1alpha1.JiraSpec{
			Ingress: &v1alpha1.JiraIngressPolicy{},
		},
	}

	j.SetDefaults()
	e := formatEndpoint(j)

	assert.NotNil(t, e)
	assert.Equal(t, "http://test-jira/", e)
}
