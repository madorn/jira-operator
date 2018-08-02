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
	"crypto/rsa"
	"crypto/x509"

	"github.com/jmckind/jira-operator/pkg/tls"
)

var (
	caCommonName         = "Jira Operator CA"
	defaultClusterDomain = "cluster.local"
	orgForTLSCert        = []string{"redhat.com"}
)

// newCACertificate will return a new CA certificate and private key.
func newCACertificate() (*rsa.PrivateKey, *x509.Certificate, error) {
	key, err := tls.NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}

	config := tls.CertConfig{
		CommonName:   caCommonName,
		Organization: orgForTLSCert,
	}

	cert, err := tls.NewSelfSignedCACertificate(config, key)
	if err != nil {
		return nil, nil, err
	}

	return key, cert, err
}

func newTLSCertificate(caCert *x509.Certificate, caPrivKey *rsa.PrivateKey, config tls.CertConfig) (*rsa.PrivateKey, *x509.Certificate, error) {
	key, err := tls.NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}
	cert, err := tls.NewSignedCertificate(config, key, caCert, caPrivKey)
	if err != nil {
		return nil, nil, err
	}
	return key, cert, nil
}
