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
	"fmt"

	"github.com/cloudflare/cfssl/log"
	"github.com/jmckind/jira-operator/pkg/apis/jira/v1alpha1"
	"github.com/jmckind/jira-operator/pkg/tls"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// ingressRules returns the rules for the Ingress resource.
func ingressRules(j *v1alpha1.Jira) []extensions.IngressRule {
	return []extensions.IngressRule{{
		Host: j.Spec.Ingress.Host,
		IngressRuleValue: extensions.IngressRuleValue{
			HTTP: &extensions.HTTPIngressRuleValue{
				Paths: []extensions.HTTPIngressPath{{
					Path: j.Spec.Ingress.Path,
					Backend: extensions.IngressBackend{
						ServiceName: j.Name,
						ServicePort: intstr.FromInt(8080),
					},
				},
				},
			},
		},
	},
	}
}

// ingressSecretCerts returns TLS certificates for the Ingress resource.
func ingressSecretCerts(j *v1alpha1.Jira) (caCert *x509.Certificate, key *rsa.PrivateKey, crt *x509.Certificate, err error) {
	caKey, caCrt, err := newCACertificate()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create ca certificate: %v", err)
	}

	config := tls.CertConfig{
		CommonName:   j.Spec.Ingress.Host,
		Organization: orgForTLSCert,
		AltNames:     tls.NewAltNames([]string{j.Spec.Ingress.Host}),
	}
	key, crt, err = newTLSCertificate(caCrt, caKey, config)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create tls certificate: %v", err)
	}
	return
}

// ingressTLS returns the TLS policy for the Ingress resource.
func ingressTLS(j *v1alpha1.Jira) []extensions.IngressTLS {
	return []extensions.IngressTLS{{
		Hosts:      []string{j.Spec.Ingress.Host},
		SecretName: j.Spec.Ingress.SecretName,
	},
	}
}

// newIngress returns a new Ingress resource.
func newIngress(j *v1alpha1.Jira) *extensions.Ingress {
	return &extensions.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      j.Name,
			Namespace: j.Namespace,
		},
	}
}

// newIngressSecret returns a new Secret resource used for Ingress.
func newIngressSecret(j *v1alpha1.Jira) *v1.Secret {
	return &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      j.Spec.Ingress.SecretName,
			Namespace: j.Namespace,
		},
	}
}

// processIngress manages the state of the Jira Ingress resource.
func processIngress(j *v1alpha1.Jira, s OperatorSDK) error {
	if !j.IsIngressEnabled() {
		logrus.Debug("ingress disabled for resource: %s", j.ObjectMeta.Name)
		return nil
	}

	ing := newIngress(j)
	err := s.Get(ing)
	if apierrors.IsNotFound(err) {
		log.Debugf("creating new ingress: %v", ing.ObjectMeta.Name)
		ing.ObjectMeta.OwnerReferences = ownerRef(j)
		ing.ObjectMeta.Labels = resourceLabels(j)
		ing.Spec.Rules = ingressRules(j)

		if j.IsIngressTLSEnabled() {
			ing.Spec.TLS = ingressTLS(j)
		}
		return s.Create(ing)
	}
	return err
}

// processIngressSecret manages the state of the Jira Ingress Secret resource.
func processIngressSecret(j *v1alpha1.Jira, s OperatorSDK) error {
	if !j.IsIngressTLSEnabled() {
		logrus.Debug("ingress tls disabled for resource: %s", j.ObjectMeta.Name)
		return nil
	}

	sec := newIngressSecret(j)
	err := s.Get(sec)
	if apierrors.IsNotFound(err) {
		log.Debugf("creating new ingress secret: %v", sec.ObjectMeta.Name)

		caCrt, key, crt, errx := ingressSecretCerts(j)
		if errx != nil {
			return fmt.Errorf("failed to create tls certificate: %v", err)
		}

		sec.ObjectMeta.OwnerReferences = ownerRef(j)
		sec.ObjectMeta.Labels = resourceLabels(j)
		sec.Data = map[string][]byte{
			"tls.key": tls.EncodePrivateKeyPEM(key),
			"tls.crt": tls.EncodeCertificatePEM(crt),
			"ca.crt":  tls.EncodeCertificatePEM(caCrt),
		}
		return s.Create(sec)
	}
	return err
}
