package fluentdforward

import (
	"strings"
	"text/template"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
	. "github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/elements"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/helpers"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/security"
	"github.com/vimalk78/collector-conf-gen/internal/generator/url"
	urlhelper "github.com/vimalk78/collector-conf-gen/internal/generator/url"
	corev1 "k8s.io/api/core/v1"
)

const (
	defaultFluentdForwardPort = "24224"
)

type FluentdForward struct {
	Desc           string
	StoreID        string
	Host           string
	Port           string
	BufferConfig   []Element
	SecurityConfig []Element
}

func (ff FluentdForward) Name() string {
	return "fluentdForwardTemplate"
}

func (ff FluentdForward) Template() string {
	return `{{define "` + ff.Name() + `" -}}
# {{.Desc}}
@type forward
@id {{.StoreID}}
<server>
  host {{.Host}}
  port {{.Port}}
</server>
heartbeat_type none
keepalive true
{{compose .SecurityConfig}}
{{compose .BufferConfig}}
{{- end}}
`
}

func (ff FluentdForward) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(ff.Template()))
}

func (ff FluentdForward) Data() interface{} {
	return ff
}

func Conf(bufspec *logging.FluentdBufferSpec, secret *corev1.Secret, o logging.OutputSpec, op *Options) []Element {
	// URL is parasable, checked at input sanitization
	u, _ := urlhelper.Parse(o.URL)
	port := u.Port()
	if port == "" {
		port = defaultFluentdForwardPort
	}
	storeID := strings.ToLower(helpers.Replacer.Replace(o.Name))
	return []Element{
		FromLabel{
			Desc:    "Output to fluentdforward",
			InLabel: helpers.LabelName(o.Name),
			SubElements: []Element{
				Match{
					MatchTags: "**",
					MatchElement: FluentdForward{
						Desc:           "FluentdForward output",
						StoreID:        storeID,
						Host:           u.Hostname(),
						Port:           port,
						SecurityConfig: SecurityConfig(o, secret),
						BufferConfig:   output.Buffer(output.NOKEYS, bufspec, storeID, &o),
					},
				},
			},
		},
	}
}

func SecurityConfig(o logging.OutputSpec, secret *corev1.Secret) []Element {
	// URL is parasable, checked at input sanitization
	u, _ := urlhelper.Parse(o.URL)
	tls := TLS(url.IsTLSScheme(u.Scheme) || secret != nil)
	conf := []Element{
		tls,
	}
	if o.Secret != nil {
		if security.HasSharedKey(secret) {
			sk := SharedKey{
				KeyPath: security.SecretPath(o.Secret.Name, "shared_key"),
			}
			conf = append(conf, sk)
		}
		if security.HasTLSKeyAndCrt(secret) {
			kc := TLSKeyCert{
				// TODO: use constants.ClientCertKey
				KeyPath:  security.SecretPath(o.Secret.Name, "tls.key"),
				CertPath: security.SecretPath(o.Secret.Name, "tls.crt"),
			}
			conf = append(conf, kc)
		}
		if security.HasCABundle(secret) {
			ca := CAFile{
				// TODO: use constants.TrustedCABundleKey
				CAFilePath: security.SecretPath(o.Secret.Name, "ca-bundle.crt"),
			}
			conf = append(conf, ca)
		}
	}
	return conf
}
