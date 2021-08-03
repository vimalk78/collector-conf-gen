package cloudwatch

import (
	"fmt"
	"text/template"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd"
	. "github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/elements"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/helpers"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/security"
	corev1 "k8s.io/api/core/v1"
)

type CloudWatch struct {
	Region         string
	SecurityConfig []Element
}

func (cw CloudWatch) Name() string {
	return "cloudwatchTemplate"
}

func (cw CloudWatch) Template() string {
	return `{{define "` + cw.Name() + `" -}}
@type cloudwatch_logs
auto_create_stream true
region {{.Region }}
log_group_name_key cw_group_name
log_stream_name_key cw_stream_name
remove_log_stream_name_key true
remove_log_group_name_key true
auto_create_stream true
concurrency 2
{{compose_one .SecurityConfig}}
include_time_key true
log_rejected_request true
{{end}}`
}

func (cw CloudWatch) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(cw.Template()))
}

func (cw CloudWatch) Data() interface{} {
	return cw
}

func Conf(bufspec *logging.FluentdBufferSpec, secret *corev1.Secret, o logging.OutputSpec, op *Options) []Element {
	logGroupPrefix := ""
	logGroupName := ""
	return []Element{
		FromLabel{
			InLabel: helpers.LabelName(o.Name),
			SubElements: []Element{
				GroupNameStreamName(fmt.Sprintf("%s%s", logGroupPrefix, logGroupName),
					"${tag}",
					fluentd.ApplicationTags),
				GroupNameStreamName(fmt.Sprintf("%sinfrastructure", logGroupPrefix),
					"${record['hostname']}.${tag}",
					fluentd.InfraTags),
				GroupNameStreamName(fmt.Sprintf("%saudit", logGroupPrefix),
					"${record['hostname']}.${tag}",
					fluentd.AuditTags),
				OutputConf(bufspec, secret, o, op),
			},
		},
	}
}

func OutputConf(bufspec *logging.FluentdBufferSpec, secret *corev1.Secret, o logging.OutputSpec, op *Options) Element {
	return Match{
		MatchTags: "**",
		MatchElement: CloudWatch{
			Region:         o.Cloudwatch.Region,
			SecurityConfig: SecurityConfig(o, secret),
		},
	}
}

func SecurityConfig(o logging.OutputSpec, secret *corev1.Secret) []Element {
	return []Element{
		AWSKey{
			KeyIDPath: security.SecretPath(o.Secret.Name, "aws_access_key_id"),
			KeyPath:   security.SecretPath(o.Secret.Name, "aws_secret_access_key"),
		},
	}
}

func GroupNameStreamName(groupName, streamName, tag string) Element {
	return Filter{
		MatchTags: tag,
		Element: RecordModifier{
			Record: map[RecordKey]RubyExpression{
				"cw_group_name":  RubyExpression(groupName),
				"cw_stream_name": RubyExpression(streamName),
			},
		},
	}
}
