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

	"github.com/cloudflare/cfssl/log"
	"github.com/jmckind/jira-operator/pkg/apis/jira/v1alpha1"
	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// containerLivenessProbe returns a new liveness Probe resource.
func containerLivenessProbe(j *v1alpha1.Jira) *v1.Probe {
	return &v1.Probe{
		Handler: v1.Handler{
			Exec: &v1.ExecAction{
				Command: []string{
					"curl",
					"--connect-timeout", "5",
					"--max-time", "10",
					"-k", "-s",
					fmt.Sprintf("http://localhost:%d/", 8080),
				},
			},
		},
		InitialDelaySeconds: 120,
		TimeoutSeconds:      10,
		PeriodSeconds:       120,
		FailureThreshold:    3,
	}
}

// containerReadinessProbe returns a new readiness Probe resource.
func containerReadinessProbe(j *v1alpha1.Jira) *v1.Probe {
	return &v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				Path:   "/",
				Port:   intstr.FromInt(8080),
				Scheme: v1.URISchemeHTTP,
			},
		},
		InitialDelaySeconds: 60,
		TimeoutSeconds:      10,
		PeriodSeconds:       30,
		FailureThreshold:    5,
	}
}

// containerResources returns a new ResourceRequirements resource.
func containerResources(j *v1alpha1.Jira) v1.ResourceRequirements {
	resources := v1.ResourceRequirements{}
	if j.Spec.Pod != nil {
		resources = j.Spec.Pod.Resources
	}
	return resources
}

// containerVolumeMounts returns a new list of VolumeMount resources.
func containerVolumeMounts(j *v1alpha1.Jira) (mounts []v1.VolumeMount) {
	mounts = make([]v1.VolumeMount, 0)
	if j.IsPVEnabled() {
		mounts = append(mounts, v1.VolumeMount{
			Name:      "jira-data",
			MountPath: j.Spec.DataMountPath,
		})
	}
	return
}

// initContainerVolumeMounts returns a new list of VolumeMount resources.
func initContainerVolumeMounts(j *v1alpha1.Jira) []v1.VolumeMount {
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

// newPod returns a new Pod resource.
func newPod(j *v1alpha1.Jira) *v1.Pod {
	return &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      j.Name,
			Namespace: j.Namespace,
		},
	}
}

// newPodSpec returns a new PodSpec resource.
func newPodSpec(j *v1alpha1.Jira) v1.PodSpec {
	return v1.PodSpec{
		InitContainers: podInitContainers(j),
		Containers:     podContainers(j),
		Volumes:        podVolumes(j),
	}
}

// podContainers returns a new list of Container resources.
func podContainers(j *v1alpha1.Jira) []v1.Container {
	return []v1.Container{{
		Name:  "jira",
		Image: fmt.Sprintf("%s:%s", j.Spec.BaseImage, j.Spec.BaseImageVersion),
		Ports: []v1.ContainerPort{{
			ContainerPort: 8080,
			Name:          "http",
		}},
		LivenessProbe:  containerLivenessProbe(j),
		ReadinessProbe: containerReadinessProbe(j),
		Resources:      containerResources(j),
		Stdin:          true,
		TTY:            true,
		VolumeMounts:   containerVolumeMounts(j),
	}}
}

// podContainers returns a new list of Init Container resources.
func podInitContainers(j *v1alpha1.Jira) []v1.Container {
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
		VolumeMounts: initContainerVolumeMounts(j),
	}
	result = append(result, ic)
	return result
}

// podContainers returns a new list of Volume resources.
func podVolumes(j *v1alpha1.Jira) []v1.Volume {
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

// processPods manages the state of the Jira Pod resources.
func processPods(j *v1alpha1.Jira, s OperatorSDK) error {
	pod := newPod(j)
	err := s.Get(pod)
	if apierrors.IsNotFound(err) {
		log.Debugf("creating new pod: %v", pod.ObjectMeta.Name)
		pod.ObjectMeta.OwnerReferences = ownerRef(j)
		pod.ObjectMeta.Labels = resourceLabels(j)
		pod.Spec = newPodSpec(j)
		return s.Create(pod)
	}
	return err
}
