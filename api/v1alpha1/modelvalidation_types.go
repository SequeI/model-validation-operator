/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Model defines the details of the model to validate, including its path and
// the path to its corresponding signature file.
type Model struct {
	// Path to the model artifact. This could be a file path on a shared volume
	// or a URI to an object store.
	// +kubebuilder:validation:Pattern=`^(/|s3://|gs://|https?://)`
	Path string `json:"path"`
	// SignaturePath is the path to the cryptographic signature bundle for the model.
	// This is used by the various validation methods to verify the model's integrity.
	// +kubebuilder:validation:Pattern=`^(/|s3://|gs://|https?://)`
	SignaturePath string `json:"signaturePath"`
}

// SigstoreConfig defines the configuration for validating model signatures using
// Sigstore's certificate-based method, which requires a specific certificate
// identity and OIDC issuer.
type SigstoreConfig struct {
	// CertificateIdentity is the expected identity of the signing certificate.
	// For example, "email:jane.doe@example.com".
	// +kubebuilder:validation:Required
	CertificateIdentity string `json:"certificateIdentity,omitempty"`
	// CertificateOidcIssuer is the URL of the OIDC issuer that issued the signing certificate.
	// For example, "https://accounts.google.com".
	// +kubebuilder:validation:Required
	CertificateOidcIssuer string `json:"certificateOidcIssuer,omitempty"`
}

// PkiConfig defines the configuration for PKI-based verification.
// This method validates the signature using a trusted certificate authority (CA).
type PkiConfig struct {
	// CertificateAuthorityPath is the path to the trusted certificate authority (CA) file.
	// The signature's chain of trust will be verified against this CA.
	// +kubebuilder:validation:Required
	CertificateAuthority string `json:"certificateAuthority,omitempty"`
}

// PublicKeyConfig defines the configuration for public key-based verification.
// This method validates the signature directly against a public key.
type PublicKeyConfig struct {
	// KeyPath is the file path to the public key used for signature verification.
	// This should be the public key corresponding to the private key used for signing.
	// +kubebuilder:validation:Required
	KeyPath string `json:"keyPath,omitempty"`
}

// ValidationConfig defines the various methods available for validating model signatures.
// Only one validation method should be specified. The controller will use the first
// non-nil method it finds.
// +kubebuilder:validation:XValidation:rule="[has(self.sigstoreConfig), has(self.pkiConfig), has(self.publicKeyConfig)].filter(x, x).size() == 1", message="exactly one validation method must be specified"
type ValidationConfig struct {
	// SigstoreConfig is the configuration for Sigstore-based signature verification.
	// +kubebuilder:validation:Optional
	SigstoreConfig *SigstoreConfig `json:"sigstoreConfig,omitempty"`
	// +kubebuilder:validation:Optional
	// PkiConfig is the configuration for traditional PKI-based signature verification.
	PkiConfig *PkiConfig `json:"pkiConfig,omitempty"`
	// +kubebuilder:validation:Optional
	// PublicKeyConfig is the configuration for public key-based signature verification.
	PublicKeyConfig *PublicKeyConfig `json:"publicKeyConfig,omitempty"`
}

// ModelValidationSpec defines the desired state of ModelValidation
type ModelValidationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Model details.
	Model Model `json:"model"`
	// Configuration for validation methods.
	// Exactly one validation method must be specified.
	Config ValidationConfig `json:"config"`
}

// ModelValidationStatus defines the observed state of ModelValidation
type ModelValidationStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// ModelValidation is the Schema for the modelvalidations API.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=mv
// +kubebuilder:printcolumn:name="Model Path",type="string",JSONPath=".spec.model.path"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type ModelValidation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ModelValidationSpec   `json:"spec,omitempty"`
	Status ModelValidationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ModelValidationList contains a list of ModelValidation
type ModelValidationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ModelValidation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ModelValidation{}, &ModelValidationList{})
}
