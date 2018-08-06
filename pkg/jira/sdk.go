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
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"k8s.io/apimachinery/pkg/runtime"
)

// OperatorSDK interface represents the interaction with the Operator SDK.
type OperatorSDK interface {
	Create(runtime.Object) error
	Delete(runtime.Object) error
	Get(runtime.Object) error
	Update(runtime.Object) error
}

// SDKWrapper implements the OperatorSDK interface for the Jira Operator.
type SDKWrapper struct {
}

// Create will use the Operator SDK to create a new resource.
func (s *SDKWrapper) Create(o runtime.Object) (err error) {
	return sdk.Create(o)
}

// Delete will use the Operator SDK to delete an existing resource.
func (s *SDKWrapper) Delete(o runtime.Object) (err error) {
	return sdk.Delete(o)
}

// Get will use the Operator SDK to fetch an existing resource.
func (s *SDKWrapper) Get(o runtime.Object) (err error) {
	return sdk.Get(o)
}

// Update will use the Operator SDK to update an existing resource.
func (s *SDKWrapper) Update(o runtime.Object) (err error) {
	return sdk.Update(o)
}
