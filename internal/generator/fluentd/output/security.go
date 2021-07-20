package output

import (
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"

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

type SecurityConf int

func (sc SecurityConf) Assemble(o *logging.OutputSpec, secret *corev1.Secret, op *Options) []Element {
	if o.Type == logging.OutputTypeElasticsearch {
		return []Element{}
	}
	return []Element{}
}
