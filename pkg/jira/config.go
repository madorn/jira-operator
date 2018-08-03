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

// DefaultDatabaseConfig is the default configuration for the JIRA database.
const DefaultDatabaseConfig = `<?xml version="1.0" encoding="UTF-8"?>
<jira-database-config>
	<name>defaultDS</name>
	<delegator-name>default</delegator-name>
	<database-type>h2</database-type>
	<schema-name>PUBLIC</schema-name>
	<jdbc-datasource>
		<url>jdbc:h2:file:/var/atlassian/jira/database/h2db</url>
		<driver-class>org.h2.Driver</driver-class>
		<username>sa</username>
		<password></password>
		<pool-min-size>20</pool-min-size>
		<pool-max-size>20</pool-max-size>
		<pool-max-wait>30000</pool-max-wait>
		<min-evictable-idle-time-millis>4000</min-evictable-idle-time-millis>
		<time-between-eviction-runs-millis>5000</time-between-eviction-runs-millis>
		<pool-max-idle>20</pool-max-idle>
		<pool-remove-abandoned>true</pool-remove-abandoned>
		<pool-remove-abandoned-timeout>300</pool-remove-abandoned-timeout>
	</jdbc-datasource>
</jira-database-config>
`

// newConfigMap returns a new ConfigMap resource.
func newConfigMap(j *v1alpha1.Jira) *v1.ConfigMap {
	return &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      j.Spec.ConfigMapName,
			Namespace: j.Namespace,
		},
	}
}

// ProcessConfigMap manages the state of the Jira ConfigMap resource.
func processConfigMap(j *v1alpha1.Jira, s OperatorSDK) error {
	cm := newConfigMap(j)
	log.Debugf("process configmap: %s", cm.ObjectMeta.Name)

	err := s.Get(cm)
	if apierrors.IsNotFound(err) {
		log.Debugf("creating new configmap: %s", cm.ObjectMeta.Name)
		cm.ObjectMeta.OwnerReferences = ownerRef(j)
		cm.ObjectMeta.Labels = resourceLabels(j)
		cm.Data = map[string]string{
			"dbconfig.xml": DefaultDatabaseConfig,
		}
		return s.Create(cm)
	}
	return err
}
