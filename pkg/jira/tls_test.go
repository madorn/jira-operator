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

	"github.com/stretchr/testify/assert"
)

// TestNewCACertificate verifies that a new certificate and private key are generated.
func TestNewCACertificate(t *testing.T) {

	key, cert, err := newCACertificate()

	assert.NotNil(t, key)
	assert.NotNil(t, cert)
	assert.Nil(t, err)
}
