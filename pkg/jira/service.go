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
	"github.com/jmckind/jira-operator/pkg/apis/jira/v1alpha1"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// DefaultServicePort is the default Jira port.
	DefaultServicePort = 8080

	// DefaultServiceName is the default standard port name/scheme.
	DefaultServiceName = "http"
)

// newService returns a new Service resource.
func newService(j *v1alpha1.Jira) *v1.Service {
	return &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      j.Name,
			Namespace: j.Namespace,
		},
	}
}

// servicePorts returns a new list of ServicePort resources.
func servicePorts(j *v1alpha1.Jira) []v1.ServicePort {
	return []v1.ServicePort{{
		Port: DefaultServicePort,
		Name: DefaultServiceName,
	}}
}

// processService manages the state of the Jira Service resource.
func processService(j *v1alpha1.Jira, s OperatorSDK) error {
	svc := newService(j)
	err := s.Get(svc)
	if apierrors.IsNotFound(err) {
		log.Debugf("creating new service: %v", svc.ObjectMeta.Name)
		svc.ObjectMeta.OwnerReferences = ownerRef(j)
		svc.ObjectMeta.Labels = resourceLabels(j)
		svc.Spec = v1.ServiceSpec{
			Selector:        resourceLabels(j),
			SessionAffinity: "ClientIP",
			Type:            "NodePort",
			Ports:           servicePorts(j),
		}
		return s.Create(svc)
	}
	return err
}
