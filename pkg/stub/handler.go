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
	"fmt"

	"github.com/jmckind/jira-operator/pkg/apis/jira/v1alpha1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
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

// NewJiraHandler constructs JiraHandler objects
func NewJiraHandler() sdk.Handler {
	return &JiraHandler{}
}

// JiraHandler handles requests for Jira!
type JiraHandler struct {
	// Fill me
}

// Handle is the starting point for processing new events.
func (h *JiraHandler) Handle(ctx context.Context, event sdk.Event) error {
	log.Debug("handle event")
	switch o := event.Object.(type) {
	case *v1alpha1.Jira:
		err := handleJira(o)
		if err != nil {
			log.Errorf("Failed to handle jira: %v", err)
			return err
		}
	}
	return nil
}

// handleJira will create the resources for the JIRA deployment
func handleJira(j *v1alpha1.Jira) (err error) {
	log.Debug("handle jira")
	j.SetDefaults()

	if err = newJiraConfigMap(j); err != nil {
		return
	}
	if err = newJiraPVC(j); err != nil {
		return
	}
	if err = newJiraPod(j); err != nil {
		return
	}
	if err = newJiraService(j); err != nil {
		return
	}
	return nil
}

// newJiraConfigMap will create a JIRA ConfigMap
func newJiraConfigMap(j *v1alpha1.Jira) error {
	cm := &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            j.Spec.ConfigMapName,
			Namespace:       j.Namespace,
			OwnerReferences: ownerRef(j),
			Labels:          jiraLabels(j),
		},
		Data: map[string]string{
			"dbconfig.xml": DefaultDatabaseConfig,
		},
	}
	return createResource(j, cm)
}

// newJiraPod will create a JIRA Pod
func newJiraPod(j *v1alpha1.Jira) error {
	pod := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            j.Name,
			Namespace:       j.Namespace,
			OwnerReferences: ownerRef(j),
			Labels:          jiraLabels(j),
		},
		Spec: jiraPodSpec(j),
	}
	return createResource(j, pod)
}

// newJiraPV will create a JIRA PersistentVolume. There is no owner assigned to
// prevent loss of data. The user must manually claen up the PVC.
func newJiraPVC(j *v1alpha1.Jira) error {
	if !j.IsPVEnabled() {
		return nil
	}
	pvc := &v1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      j.Name,
			Namespace: j.Namespace,
			Labels:    jiraLabels(j),
		},
		Spec: *j.Spec.Pod.PersistentVolumeClaimSpec,
	}
	return createResource(j, pvc)
}

// newJiraService will create a JIRA Service
func newJiraService(j *v1alpha1.Jira) error {
	svc := &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            j.Name,
			Namespace:       j.Namespace,
			OwnerReferences: ownerRef(j),
			Labels:          jiraLabels(j),
		},
		Spec: v1.ServiceSpec{
			Selector:        jiraLabels(j),
			SessionAffinity: "ClientIP",
			Type:            "NodePort",
			Ports:           servicePorts(j),
		},
	}
	return createResource(j, svc)
}

// defaultLabels returns the default labels.
func defaultLabels(j *v1alpha1.Jira) map[string]string {
	return map[string]string{
		"app":     "jira",
		"cluster": j.Name,
	}
}

// jiraLabels returns the labels for the resource.
func jiraLabels(j *v1alpha1.Jira) map[string]string {
	labels := defaultLabels(j)
	for key, val := range j.ObjectMeta.Labels {
		labels[key] = val
	}
	return labels
}

// ownerRef returns an owner reference for the operator
func ownerRef(j *v1alpha1.Jira) []metav1.OwnerReference {
	return []metav1.OwnerReference{
		*metav1.NewControllerRef(j, schema.GroupVersionKind{
			Group:   v1alpha1.SchemeGroupVersion.Group,
			Version: v1alpha1.SchemeGroupVersion.Version,
			Kind:    "Jira",
		}),
	}
}

func initContainers(j *v1alpha1.Jira) []v1.Container {
	result := make([]v1.Container, 0)
	if !j.IsPVEnabled() {
		return result
	}

	mp := j.Spec.DataMountPath
	ic := v1.Container{
		Name:  "init",
		Image: "busybox",
		Command: []string{
			"/bin/sh",
			"-c",
			fmt.Sprintf("cp /etc/jira/dbconfig.xml %s/; chown -R 2:2 %s", mp, mp),
		},
		VolumeMounts: initVolumeMounts(j),
	}
	result = append(result, ic)
	return result
}

func jiraContainers(j *v1alpha1.Jira) []v1.Container {
	return []v1.Container{{
		Name:  "jira",
		Image: fmt.Sprintf("%s:%s", j.Spec.BaseImage, j.Spec.BaseImageVersion),
		Ports: []v1.ContainerPort{{
			ContainerPort: 8080,
			Name:          "http",
		}},
		Resources:    containerResources(j),
		Stdin:        true,
		TTY:          true,
		VolumeMounts: jiraVolumeMounts(j),
	}}
}

// jiraPodSpec returns a PodSpec for a JIRA container.
func jiraPodSpec(j *v1alpha1.Jira) v1.PodSpec {
	return v1.PodSpec{
		InitContainers: initContainers(j),
		Containers:     jiraContainers(j),
		Volumes:        jiraVolumes(j),
	}
}

// servicePorts returns the ports for the JIRA service.
func servicePorts(j *v1alpha1.Jira) []v1.ServicePort {
	return []v1.ServicePort{{
		Port: 8080,
		Name: "http",
	}}
}

// containerResources returns the resources requestd for the application.
func containerResources(j *v1alpha1.Jira) v1.ResourceRequirements {
	resources := v1.ResourceRequirements{}
	if j.Spec.Pod != nil {
		resources = j.Spec.Pod.Resources
	}
	return resources
}

func jiraVolumeMounts(j *v1alpha1.Jira) (mounts []v1.VolumeMount) {
	mounts = make([]v1.VolumeMount, 0)
	if j.IsPVEnabled() {
		mounts = append(mounts, v1.VolumeMount{
			Name:      "jira-data",
			MountPath: j.Spec.DataMountPath,
		})
	}
	return
}

func initVolumeMounts(j *v1alpha1.Jira) []v1.VolumeMount {
	return []v1.VolumeMount{
		{
			Name:      "jira-data",
			MountPath: j.Spec.DataMountPath,
		},
		{
			Name:      "jira-config",
			MountPath: "/etc/jira",
		},
	}
}

func jiraVolumes(j *v1alpha1.Jira) []v1.Volume {
	volumes := make([]v1.Volume, 0)
	cmv := v1.Volume{
		Name: "jira-config",
		VolumeSource: v1.VolumeSource{
			ConfigMap: &v1.ConfigMapVolumeSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name: j.Spec.ConfigMapName,
				},
				Items: []v1.KeyToPath{
					{Key: "dbconfig.xml", Path: "dbconfig.xml"},
				},
			},
		},
	}
	volumes = append(volumes, cmv)

	if j.IsPVEnabled() {
		pv := v1.Volume{
			Name: "jira-data",
			VolumeSource: v1.VolumeSource{
				PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
					ClaimName: j.Name},
			},
		}
		volumes = append(volumes, pv)
	}
	return volumes
}

func createResource(j *v1alpha1.Jira, o sdk.Object) error {
	err := sdk.Create(o)
	if errors.IsAlreadyExists(err) {
		log.Debug("resource already exists")
	} else if err != nil {
		log.Errorf("Failed to create resource: %v", err)
		return err
	}
	return nil
}
