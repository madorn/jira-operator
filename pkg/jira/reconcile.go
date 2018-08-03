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

	"github.com/jmckind/jira-operator/pkg/apis/jira/v1alpha1"
	log "github.com/sirupsen/logrus"
)

// Reconciler ensures that the actual state of the Jira resource matches the desired state.
type Reconciler struct {
	resource *v1alpha1.Jira
	sdk      OperatorSDK
}

// NewReconciler will create a new Reconciler using the given sdk.
// If a nil sdk value is provided, the function will retun nil as well.
func NewReconciler(s OperatorSDK) *Reconciler {
	if s == nil {
		return nil
	}
	return &Reconciler{sdk: s}
}

// Reconcile the given Jira resource.
func (r *Reconciler) Reconcile(jr *v1alpha1.Jira) (err error) {
	if jr == nil {
		return errors.New("jira reference cannot be nil")
	}

	log.Debugf("reconciling resource: %s", jr.ObjectMeta.Name)
	r.resource = jr.DeepCopy()

	changed := r.resource.SetDefaults()
	if changed {
		log.Debug("simulating initializer")
		return r.sdk.Update(r.resource)
	}

	if err = processConfigMap(r.resource, r.sdk); err != nil {
		return
	}

	if err = processPVC(r.resource, r.sdk); err != nil {
		return
	}

	if err = processPods(r.resource, r.sdk); err != nil {
		return
	}

	if err = processService(r.resource, r.sdk); err != nil {
		return
	}

	if err = processIngressSecret(r.resource, r.sdk); err != nil {
		return
	}

	if err = processIngress(r.resource, r.sdk); err != nil {
		return
	}

	log.Debugf("finished reconciling resource: %s", jr.ObjectMeta.Name)
	return nil
}
