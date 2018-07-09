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

package stub

import (
	"context"

	"github.com/jmckind/jira-operator/pkg/apis/jira/v1alpha1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewJiraHandler() sdk.Handler {
	return &JiraHandler{}
}

type JiraHandler struct {
	// Fill me
}

func (h *JiraHandler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.Jira:
		err := handleJira(o)
		if err != nil {
			logrus.Errorf("Failed to handle jira: %v", err)
			return err
		}
	}
	return nil
}

// handleJira will create the resources for the cluster
func handleJira(cr *v1alpha1.Jira) error {
	err := newJiraConfigMap(cr)
	if err != nil && !errors.IsAlreadyExists(err) {
		logrus.Errorf("Failed to create jira configmap: %v", err)
		return err
	}

	err = newJiraPod(cr)
	if err != nil && !errors.IsAlreadyExists(err) {
		logrus.Errorf("Failed to create jira pod: %v", err)
		return err
	}

	err = newJiraService(cr)
	if err != nil && !errors.IsAlreadyExists(err) {
		logrus.Errorf("Failed to create jira service: %v", err)
		return err
	}

	return nil
}

// newJiraConfigMap will create a jira configmap
func newJiraConfigMap(cr *v1alpha1.Jira) error {
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:            cr.Name,
			Namespace:       cr.Namespace,
			OwnerReferences: ownerForCluster(cr),
			Labels:          labelsForCluster(cr),
		},
		Data: map[string]string{
			"dbconfig.xml":           "",
			"jira-config.properties": "",
		},
	}

	return sdk.Create(configMap)
}

// newJiraPod will create a jira pod
func newJiraPod(cr *v1alpha1.Jira) error {
	pod := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            cr.Name,
			Namespace:       cr.Namespace,
			OwnerReferences: ownerForCluster(cr),
			Labels:          labelsForCluster(cr),
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{{
				Name:  "jira",
				Image: "cptactionhank/atlassian-jira:7.10.0",
				Ports: []v1.ContainerPort{{
					ContainerPort: 8080,
					Name:          "http",
				},
				},
			},
			},
		},
	}

	return sdk.Create(pod)
}

// newJiraService will create a jira service
func newJiraService(cr *v1alpha1.Jira) error {
	svc := &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            cr.Name,
			Namespace:       cr.Namespace,
			OwnerReferences: ownerForCluster(cr),
			Labels:          labelsForCluster(cr),
		},
		Spec: v1.ServiceSpec{
			Selector:        labelsForCluster(cr),
			SessionAffinity: "ClientIP",
			Type:            "NodePort",
			Ports:           portsForService(cr),
		},
	}

	return sdk.Create(svc)
}

// labelsForCluster will create the labels for the cluster
func labelsForCluster(cr *v1alpha1.Jira) map[string]string {
	return map[string]string{
		"app":     "jira",
		"cluster": cr.Name,
	}
}

// ownerForCluster will create the owner references for the cluster
func ownerForCluster(cr *v1alpha1.Jira) []metav1.OwnerReference {
	return []metav1.OwnerReference{
		*metav1.NewControllerRef(cr, schema.GroupVersionKind{
			Group:   v1alpha1.SchemeGroupVersion.Group,
			Version: v1alpha1.SchemeGroupVersion.Version,
			Kind:    "Jira",
		}),
	}
}

// portsForService will create the ports for the jira service
func portsForService(cr *v1alpha1.Jira) []v1.ServicePort {
	return []v1.ServicePort{{
		Port: 8080,
		Name: "http",
	}}
}
