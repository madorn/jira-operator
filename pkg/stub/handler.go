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

	"github.com/coreos/jira-operator/pkg/apis/jira/v1alpha1"
	"github.com/coreos/jira-operator/pkg/jira"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
)

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
	switch o := event.Object.(type) {
	case *v1alpha1.Jira:
		r := jira.NewReconciler(new(jira.SDKWrapper))
		return r.Reconcile(o)
	}
	return nil
}
