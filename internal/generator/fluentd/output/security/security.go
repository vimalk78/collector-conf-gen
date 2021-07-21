package security

import (
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
)

type TLS bool

type TLSKeyCert struct {
	KeyPath  string
	CertPath string
}

type UserNamePass struct {
	UsernamePath string
	PasswordPath string
}

type SharedKey struct {
	KeyPath string
}

type CAFile struct {
	CAFilePath string
}

var NoSecrets = map[string]*corev1.Secret{}

func HasUsernamePassword(secret *corev1.Secret) bool {
	if secret == nil {
		return false
	}

	// TODO: use constants.ClientUsername
	if _, ok := secret.Data["username"]; !ok {
		return false
	}
	if _, ok := secret.Data["password"]; !ok {
		return false
	}
	return true
}

func HasTLSKeyAndCrt(secret *corev1.Secret) bool {
	if secret == nil {
		return false
	}

	// TODO: use constants.ClientCertKey
	if _, ok := secret.Data["tls.crt"]; !ok {
		return false
	}
	if _, ok := secret.Data["tls.key"]; !ok {
		return false
	}
	return true
}

func HasCABundle(secret *corev1.Secret) bool {
	if secret == nil {
		return false
	}

	// TODO: use constants.TrustedCABundleKey
	if _, ok := secret.Data["ca-bundle.crt"]; !ok {
		return false
	}
	return true
}

func SecretPath(secret *corev1.Secret, file string) string {
	if secret != nil {
		// TODO use constants.CollectorSecretsDir
		return filepath.Join("/var/run/ocp-collector/secrets", secret.Name, file)
	}
	return ""
}
