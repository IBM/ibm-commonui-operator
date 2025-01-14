//
// Copyright 2021 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package resources

import (
	"context"
	"time"

	operatorsv1alpha1 "github.com/IBM/ibm-commonui-operator/api/v1alpha1"
	// certmgr "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
	certmgr "github.com/ibm/ibm-cert-manager-operator/apis/cert-manager/v1"
	cmmeta "github.com/ibm/ibm-cert-manager-operator/apis/meta.cert-manager/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type CertificateData struct {
	Name      string
	Secret    string
	Common    string
	App       string
	Component string
}

// nolint
func getDesiredCertificate(ctx context.Context, client client.Client, instance *operatorsv1alpha1.CommonWebUI, data CertificateData) (*certmgr.Certificate, error) {
	reqLogger := log.WithValues("func", "getDesiredCertificate", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)

	metaLabels := map[string]string{
		"app":                          data.App,
		"component":                    data.Common,
		"release":                      ReleaseName,
		"app.kubernetes.io/instance":   "ibm-commonui-operator",
		"app.kubernetes.io/managed-by": "ibm-commonui-operator",
		"app.kubernetes.io/name":       UICertName,
		"manage-cert-rotation":         "yes",
	}

	certificate := &certmgr.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Labels:    metaLabels,
			Namespace: instance.Namespace,
		},
		Spec: certmgr.CertificateSpec{
			CommonName: data.Common,
			SecretName: data.Secret,
			IsCA:       false,
			DNSNames: []string{
				data.Common,
				data.Common + "." + instance.Namespace,
				data.Common + "." + instance.Namespace + ".svc.cluster.local",
			},
			// Organization: []string{"IBM"},
			IssuerRef: cmmeta.ObjectReference{
				Name: DefaultClusterIssuer,
				Kind: certmgr.IssuerKind,
			},
			Duration: &metav1.Duration{
				Duration: 9552 * time.Hour, /* 398 days */
			},
			RenewBefore: &metav1.Duration{
				Duration: 2880 * time.Hour, /* 120 days (3 months) */
			},
		},
	}

	err := controllerutil.SetControllerReference(instance, certificate, client.Scheme())
	if err != nil {
		reqLogger.Error(err, "Failed to set owner for certificate")
		return nil, err
	}

	return certificate, nil
}

func ReconcileCertificates(ctx context.Context, client client.Client, instance *operatorsv1alpha1.CommonWebUI, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileCertificates", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Reconciling certificates")

	certs := []CertificateData{
		UICertificateData,
	}

	for _, certData := range certs {
		reqLogger.Info("Checking certificate", "Certificate.Name", certData.Name)

		certificate := &certmgr.Certificate{}

		desiredCertificate, desiredErr := getDesiredCertificate(ctx, client, instance, certData)
		if desiredErr != nil {
			return desiredErr
		}

		err := client.Get(ctx, types.NamespacedName{Name: certData.Name, Namespace: instance.Namespace}, certificate)

		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating a new certificate", "Certificate.Namespace", desiredCertificate.Namespace, "Certificate.Name", desiredCertificate.Name)

			err = client.Create(ctx, desiredCertificate)
			if err != nil {
				if errors.IsAlreadyExists(err) {
					// Certificate already exists from a previous reconcile
					reqLogger.Info("Certificate already exists")
					*needToRequeue = true
				} else {
					// Failed to create a new certificate
					reqLogger.Info("Failed to create a new certificate", "Certificate.Namespace", desiredCertificate.Namespace, "Certificate.Name", desiredCertificate.Name)
					return err
				}
			} else {
				// Requeue after creating new certificate
				*needToRequeue = true
			}
		} else if err != nil && !errors.IsNotFound(err) {
			reqLogger.Error(err, "Failed to get certificate", "Certificate.Namespace", instance.Namespace, "Certificate.Name", certData.Name)
			return err
		} else {
			// Determine if current certificate has changed
			reqLogger.Info("Comparing current and desired certificates")

			if !IsCertificateEqual(certificate, desiredCertificate) {
				reqLogger.Info("Updating certificate", "Certificate.Namespace", certificate.Namespace, "Certificate.Name", certificate.Name)

				certificate.ObjectMeta.Name = desiredCertificate.ObjectMeta.Name
				certificate.ObjectMeta.Labels = desiredCertificate.ObjectMeta.Labels
				certificate.Spec = desiredCertificate.Spec

				err = client.Update(ctx, certificate)
				if err != nil {
					reqLogger.Error(err, "Failed to update certificate", "Certificate.Namespace", certificate.Namespace, "Certificate.Name", certificate.Name)
					return err
				}
			}
		}
	}

	return nil
}
