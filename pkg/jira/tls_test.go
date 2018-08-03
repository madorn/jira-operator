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

	"github.com/jmckind/jira-operator/pkg/tls"
	"github.com/stretchr/testify/assert"
)

// TestNewCACertificate verifies that a new CA certificate and private key are generated.
func TestNewCACertificate(t *testing.T) {

	key, crt, err := newCACertificate()

	assert.Nil(t, err)
	assert.NotEmpty(t, key)
	assert.NotEmpty(t, crt)
}

// TestNewTLSCertificate verifies that a new certificate and private key are generated.
func TestNewTLSCertificate(t *testing.T) {
	caKey, caCrt, err := newCACertificate()
	config := tls.CertConfig{
		CommonName:   "test-ingress-host",
		Organization: orgForTLSCert,
		AltNames:     tls.NewAltNames([]string{"test-alt-name"}),
	}

	key, crt, err := newTLSCertificate(caCrt, caKey, config)

	assert.Nil(t, err)
	assert.NotEmpty(t, caCrt)
	assert.NotEmpty(t, key)
	assert.NotEmpty(t, crt)
}
