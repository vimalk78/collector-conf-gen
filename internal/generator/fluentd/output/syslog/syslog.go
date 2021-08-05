package syslog

import (
	"fmt"
	"regexp"
	"strings"
	"text/template"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	corev1 "k8s.io/api/core/v1"

	. "github.com/vimalk78/collector-conf-gen/internal/generator"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd"
	. "github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/elements"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/helpers"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output"
	"github.com/vimalk78/collector-conf-gen/internal/generator/fluentd/output/security"
	"github.com/vimalk78/collector-conf-gen/internal/generator/url"
	urlhelper "github.com/vimalk78/collector-conf-gen/internal/generator/url"
)

type Syslog struct {
	Desc           string
	StoreID        string
	Host           string
	Port           string
	Facility       string
	Severity       string
	AppName        Element
	MsgID          Element
	ProcID         Element
	Tag            Element
	Protocol       string
	PayloadKey     string
	SecurityConfig []Element
	BufferConfig   []Element
}

func (s Syslog) Name() string {
	return "syslogTemplate"
}

func (s Syslog) Template() string {
	return `{{define "` + s.Name() + `" -}}
@type remote_syslog
@id {{.StoreID}}
host {{.Host}}
port {{.Port}}
rfc {{.Rfc}}
{{kv .Facility -}}
{{kv .Severity -}}
{{kv .AppName -}}
{{kv .MsgID -}}
{{kv .ProcID -}}
{{kv .Tag -}}
protocol {{.Protocol}}
packet_size 4096
hostname "#{ENV['NODE_NAME']}"
{{compose .SecurityConfig}}
{{if (eq .Protocol "tcp") -}}
timeout 60
timeout_exception true
keep_alive true
keep_alive_idle 75
keep_alive_cnt 9
keep_alive_intvl 7200
{{end -}}
{{if .PayloadKey -}}
<format>
  @type single_value
  message_key {{.PayloadKey}}
</format>
{{end -}}
{{compose .BufferConfig}}
</store>
{{end}}
`
}

func (s Syslog) Create(t *template.Template) *template.Template {
	return template.Must(t.Parse(s.Template()))
}

func (s Syslog) Data() interface{} {
	return s
}

func Conf(bufspec *logging.FluentdBufferSpec, secret *corev1.Secret, o logging.OutputSpec, op *Options) []Element {
	var addLogSource Element = Nil
	if o.Syslog != nil && o.Syslog.AddLogSource {
		addLogSource = AddLogSource(o, op)
	}
	return []Element{
		FromLabel{
			InLabel: helpers.LabelName(o.Name),
			SubElements: []Element{
				ParseJson(o, op),
				addLogSource,
				OutputConf(bufspec, secret, o, fluentd.JournalTags, op),
				OutputConf(bufspec, secret, o, "**", op),
			},
		},
	}
}

func OutputConf(bufspec *logging.FluentdBufferSpec, secret *corev1.Secret, o logging.OutputSpec, tags string, op *Options) Element {
	// URL is parasable, checked at input sanitization
	u, _ := urlhelper.Parse(o.URL)
	port := u.Port()
	if port == "" {
		port = ""
	}
	storeID := helpers.StoreID(o.Name, false)
	return Match{
		MatchTags: tags,
		MatchElement: Syslog{
			StoreID:        storeID,
			Host:           u.Hostname(),
			Port:           port,
			Facility:       Facility(o.Syslog),
			Severity:       Severity(o.Syslog),
			AppName:        AppName(o.Syslog, tags),
			MsgID:          MsgID(o.Syslog),
			ProcID:         ProcID(o.Syslog),
			Tag:            Tag(o.Syslog, tags),
			Protocol:       Protocol(o),
			SecurityConfig: SecurityConfig(o, secret),
			BufferConfig:   output.Buffer(output.NOKEYS, bufspec, storeID, &o),
		},
	}
}

func Facility(s *logging.Syslog) string {
	if s == nil || s.Facility == "" {
		return "user"
	}
	if IsKeyExpr(s.Facility) {
		return fmt.Sprintf("${%s}", s.Facility)
	}
	return s.Facility
}

func Severity(s *logging.Syslog) string {
	if s == nil || s.Severity == "" {
		return "debug"
	}
	if IsKeyExpr(s.Severity) {
		return fmt.Sprintf("${%s}", s.Severity)
	}
	return s.Severity
}

func Rfc(s *logging.Syslog) string {
	if s == nil || s.RFC == "" {
		return "rfc5424"
	}
	switch strings.ToLower(s.RFC) {
	case "rfc3164":
		return "rfc3164"
	case "rfc5424":
		return "rfc5424"
	}
	return "Unknown Rfc"
}

func AppName(s *logging.Syslog, matchtags string) Element {
	if s == nil {
		return Nil
	}
	appname := "appname"
	if s.AppName == "" {
		if matchtags == fluentd.JournalTags && Rfc(s) == "rfc5424" {
			return KV(appname, "${$.systemd.u.SYSLOG_IDENTIFIER}")
		} else {
			return Nil
		}
	}
	if IsKeyExpr(s.AppName) {
		return KV(appname, fmt.Sprintf("${%s}", s.AppName))
	}
	if IsTagExpr(s.AppName) {
		return KV(appname, s.AppName)
	}
	if s.AppName == "tag" {
		return KV(appname, "${tag}")
	}
	return KV(appname, s.AppName)
}

func Tag(s *logging.Syslog, matchtags string) Element {
	if s == nil {
		return Nil
	}
	program := "program"
	if s.Tag == "" {
		if matchtags == fluentd.JournalTags && Rfc(s) == "rfc3164" {
			return KV(program, "${$.systemd.u.SYSLOG_IDENTIFIER}")
		} else {
			return Nil
		}
	}
	if IsKeyExpr(s.Tag) {
		return KV(program, fmt.Sprintf("${%s}", s.Tag))
	}
	if IsTagExpr(s.Tag) {
		return KV(program, s.Tag)
	}
	if s.Tag == "tag" {
		return KV(program, "${tag}")
	}
	return KV(program, s.Tag)
}

func MsgID(s *logging.Syslog) Element {
	if s == nil || s.MsgID == "" {
		return Nil
	}
	msgid := "msgid"
	if IsKeyExpr(s.MsgID) {
		return KV(msgid, fmt.Sprintf("${%s}", s.MsgID))
	}
	return KV(msgid, s.MsgID)
}

func ProcID(s *logging.Syslog) Element {
	if s == nil || s.ProcID == "" {
		return Nil
	}
	procid := "procid"
	if IsKeyExpr(s.ProcID) {
		return KV(procid, fmt.Sprintf("${%s}", s.ProcID))
	}
	return KV(procid, s.ProcID)
}

func Protocol(o logging.OutputSpec) string {
	u, _ := urlhelper.Parse(o.URL)
	return urlhelper.PlainScheme(u.Scheme)
}

func BufferKeys(s *logging.Syslog, matchtags string) []string {
	if s == nil {
		return output.NOKEYS
	}
	keys := []string{}
	tagAdded := false
	if matchtags == fluentd.JournalTags {
		keys = append(keys, "$.systemd.u.SYSLOG_IDENTIFIER")
	}
	if IsKeyExpr(s.Tag) {
		keys = append(keys, s.Tag)
	}
	if IsTagExpr(s.Tag) && !tagAdded {
		keys = append(keys, "tag")
		tagAdded = true
	}
	if s.Tag == "tag" && !tagAdded {
		keys = append(keys, "tag")
		tagAdded = true
	}
	if IsKeyExpr(s.AppName) {
		keys = append(keys, s.AppName)
	}
	if IsTagExpr(s.AppName) && !tagAdded {
		keys = append(keys, "tag")
		tagAdded = true
	}
	if s.AppName == "tag" && !tagAdded {
		keys = append(keys, "tag")
	}
	if IsKeyExpr(s.MsgID) {
		keys = append(keys, s.MsgID)
	}
	if IsKeyExpr(s.PayloadKey) {
		keys = append(keys, s.ProcID)
	}
	if IsKeyExpr(s.Facility) {
		keys = append(keys, s.Facility)
	}
	if IsKeyExpr(s.Severity) {
		keys = append(keys, s.Severity)
	}
	return keys
}

func ParseJson(o logging.OutputSpec, op *Options) Element {
	return Filter{
		MatchTags: "**",
		Element: ConfLiteral{
			TemplateName: "syslogParseJson",
			TemplateStr: `{{define "syslogParseJson" -}}
@type parse_json_field
json_fields  message
merge_json_log false
replace_json_log true
{{end}}`,
		},
	}
}

func AddLogSource(o logging.OutputSpec, op *Options) Element {
	return Filter{
		MatchTags: "**",
		Element: RecordModifier{
			Record: map[RecordKey]RubyExpression{
				"kubernetes_info": RubyExpression(`${if record.has_key?('kubernetes'); record['kubernetes']; else {}; end}`),
				"namespace_info":  RubyExpression(`${if record['kubernetes_info'] != nil && record['kubernetes_info'] != {}; "namespace_name=" + record['kubernetes_info']['namespace_name']; else nil; end}`),
				"pod_info":        RubyExpression(`${if record['kubernetes_info'] != nil && record['kubernetes_info'] != {}; "pod_name=" + record['kubernetes_info']['pod_name']; else nil; end}`),
				"container_info":  RubyExpression(`${if record['kubernetes_info'] != nil && record['kubernetes_info'] != {}; "container_name=" + record['kubernetes_info']['container_name']; else nil; end}`),
				"msg_key":         RubyExpression(`${if record.has_key?('message') && record['message'] != nil; record['message']; else nil; end}`),
				"msg_info":        RubyExpression(`${if record['msg_key'] != nil && record['msg_key'].is_a?(Hash); require 'json'; "message="+record['message'].to_json; elsif record['msg_key'] != nil; "message="+record['message']; else nil; end}`),
				"message":         RubyExpression(`${if record['msg_key'] != nil && record['kubernetes_info'] != nil && record['kubernetes_info'] != {}; record['namespace_info'] + ", " + record['container_info'] + ", " + record['pod_info'] + ", " + record['msg_info']; else record['message']; end}`),
				"systemd_info":    RubyExpression(`${if record.has_key?('systemd') && record['systemd']['t'].has_key?('PID'); record['systemd']['u']['SYSLOG_IDENTIFIER'] += "[" + record['systemd']['t']['PID'] + "]"; else {}; end}`),
			},
			RemoveKeys: []string{
				"kubernetes_info",
				"namespace_info",
				"pod_info",
				"container_info",
				"msg_key",
				"msg_info",
				"systemd_info",
			},
		},
	}
}

func SecurityConfig(o logging.OutputSpec, secret *corev1.Secret) []Element {
	u, _ := urlhelper.Parse(o.URL)
	tls := TLS(url.IsTLSScheme(u.Scheme) || secret != nil)
	conf := []Element{
		tls,
	}
	if security.HasCABundle(secret) {
		ca := CAFile{
			// TODO: use constants.TrustedCABundleKey
			CAFilePath: security.SecretPath(o.Secret.Name, "ca-bundle.crt"),
		}
		conf = append(conf, ca)
	}
	return conf
}

// The Syslog output fields can be set to an expression of the form $.abc.xyz
// If an expression is used, its value will be taken from corresponding key in the record
var keyre = regexp.MustCompile(`^\$(\.[[:word:]]*)+$`)

var tagre = regexp.MustCompile(`\${tag\[-??\d+\]}`)

func IsKeyExpr(str string) bool {
	return keyre.MatchString(str)
}

func IsTagExpr(str string) bool {
	return tagre.MatchString(str)
}

/*
//---
func isFacilityKeyExpr(facility string) bool {
	return conf.IsKeyExpr(facility)
}

func isSeverityKeyExpr(severity string) bool {
	return conf.IsKeyExpr(severity)
}

func isTagKeyExpr(tag string) bool {
	return conf.IsKeyExpr(tag)
}

func isTagTagExpr(tag string) bool {
	return conf.IsTagExpr(tag)
}

func IsPayloadKeyExpr(payload string) bool {
	return conf.IsKeyExpr(payload)
}

func isAppNameKeyExpr() bool {
	return conf.IsKeyExpr(conf.Target.Syslog.AppName)
}

func isAppNameTagExpr() bool {
	return conf.IsTagExpr(conf.Target.Syslog.AppName)
}

func isMsgIDKeyExpr() bool {
	return conf.IsKeyExpr(conf.Target.Syslog.MsgID)
}

func isProcIDKeyExpr() bool {
	return conf.IsKeyExpr(conf.Target.Syslog.ProcID)
}
*/
