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
	"reflect"

	"github.com/coreos/jira-operator/pkg/apis/jira/v1alpha1"
	log "github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

const (
	// StateAvailable inicates the Jira application is available for clients.
	StateAvailable = "Available"

	// StateInitializing inicates the Jira application is initializing.
	StateInitializing = "Initializing"

	// StateUnavailable inicates the Jira application is unavailable for clients.
	StateUnavailable = "Unavailable"
)

// formatEndpoint will return the URI for accessing the application.
func formatEndpoint(j *v1alpha1.Jira) string {
	scheme := "http"
	host := fmt.Sprintf("%s:%d", j.Name, DefaultServicePort)
	path := v1alpha1.DefaultIngressPath

	if j.IsIngressEnabled() {
		host = j.Spec.Ingress.Host
		path = j.Spec.Ingress.Path
	}
	if j.IsIngressTLSEnabled() {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s%s", scheme, host, path)
}

// jiraState will return the current state for the Jira application.
func jiraState(j *v1alpha1.Jira, s OperatorSDK) string {
	log.Debugf("begin determine state")
	pod := newPod(j)
	err := s.Get(pod)

	if apierrors.IsNotFound(err) {
		log.Debugf("pod not found: initializing")
		return StateInitializing
	} else if err != nil {
		log.Errorf("unable to get pod: %v", err)
		return StateUnavailable
	}

	if len(pod.Status.ContainerStatuses) > 0 {
		if pod.Status.ContainerStatuses[0].Ready {
			log.Debugf("pod ready: available")
			return StateAvailable
		} else {
			log.Debugf("pod not ready: initializing")
			return StateInitializing
		}
	}

	log.Debugf("unable to determine state, inidacting unavailable...")
	return StateUnavailable
}

// processStatus will process the status properties for the Jira resource.
func processStatus(j *v1alpha1.Jira, s OperatorSDK) error {
	status := v1alpha1.JiraStatus{
		Endpoint:    formatEndpoint(j),
		ServiceName: j.ObjectMeta.Name,
		State:       jiraState(j, s),
	}

	// don't update the status if there aren't any changes.
	if reflect.DeepEqual(j.Status, status) {
		return nil
	}

	log.Debugf("updating status for resource: %s", j.Name)
	j.Status = status
	return s.Update(j)
}
