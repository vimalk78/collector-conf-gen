package fluentd

import (
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/vimalk78/collector-conf-gen/internal/generator"
)

func Concat(spec *logging.ClusterLogForwarderSpec, o *Options) []Element {
	return []Element{
		Pipeline{
			InLabel: labelName("CONCAT"),
			Desc:    "Concat log lines of container logs",
			SubElements: []Element{
				ConfLiteral{
					Desc:         "Concat container lines",
					TemplateName: "concatLines",
					TemplateStr:  ConcatLines,
					OutLabel:     labelName("INGRESS"),
				},
				Relabel{
					Desc:      "Kubernetes Logs go to INGRESS pipeline",
					MatchTags: "kubernetes.**",
					OutLabel:  labelName("INGRESS"),
				},
			},
		},
	}
}

func Ingress(spec *logging.ClusterLogForwarderSpec, o *Options) []Element {
	return []Element{
		Pipeline{
			InLabel: labelName("INGRESS"),
			Desc:    "Ingress pipeline",
			SubElements: MergeElements([]Element{
				ConfLiteral{
					Desc:         "Set Encodeing",
					TemplateName: "setEncoding",
					TemplateStr:  SetEncoding,
				},
				ConfLiteral{
					Desc:         "Filter out PRIORITY from journal logs",
					TemplateName: "filterJournalPRIORITY",
					TemplateStr:  FilterJournalPRIORITY,
				},
				ConfLiteral{
					Desc:         "Retag Journal logs to specific tags",
					OutLabel:     "INGRESS",
					TemplateName: "retagJournal",
					TemplateStr:  RetagJournalLogs,
				},
				ConfLiteral{
					Desc:         "Invoke kubernetes apiserver to get kunbernetes metadata",
					TemplateName: "kubernetesMetadata",
					TemplateStr:  KubernetesMetadataPlugin,
				},
				ConfLiteral{
					Desc:         "Parse Json fields for container, journal and eventrouter logs",
					TemplateName: "parseJsonFields",
					TemplateStr:  ParseJsonFields,
				},
				ConfLiteral{
					Desc:         "Clean kibana log fields",
					TemplateName: "cleanKibanaLogs",
					TemplateStr:  CleanKibanaLogs,
				},
				ConfLiteral{
					Desc:         "Fix level field in audit logs",
					TemplateName: "fixAuditLevel",
					TemplateStr:  FixAuditLevel,
				},
				ConfLiteral{
					Desc:         "Viaq Data Model: The big bad viaq model.",
					TemplateName: "viaqDataModel",
					TemplateStr:  ViaQDataModel,
				},
				ConfLiteral{
					Desc:         "Generate elasticsearch id",
					TemplateName: "genElasticsearchID",
					TemplateStr:  GenElasticsearchID,
				},
			},
				SourcesToInputs(spec, o)),
		},
	}
}

var SetEncoding string = `
{{define "setEncoding"}}
# {{.Desc}}
<filter **>
  @type record_modifier
  char_encoding utf-8
</filter>
{{- end}}
`

var FilterJournalPRIORITY string = `
{{define "filterJournalPRIORITY"}}
# {{.Desc}}
<filter journal>
  @type grep
  <exclude>
    key PRIORITY
    pattern ^7$
  </exclude>
</filter>
{{- end}}
`

var RetagJournalLogs string = `
# {{.Desc}}
{{define "retagJournal"}}
<match journal>
  @type rewrite_tag_filter
  # skip to @INGRESS label section
  @label @{{.OutLabel}}

  # see if this is a kibana container for special log handling
  # looks like this:
  # k8s_kibana.a67f366_logging-kibana-1-d90e3_logging_26c51a61-2835-11e6-ad29-fa163e4944d5_f0db49a2
  # we filter these logs through the kibana_transform.conf filter
  <rule>
    key CONTAINER_NAME
    pattern ^k8s_kibana\.
    tag kubernetes.journal.container.kibana
  </rule>

  <rule>
    key CONTAINER_NAME
    pattern ^k8s_[^_]+_logging-eventrouter-[^_]+_
    tag kubernetes.journal.container._default_.kubernetes-event
  </rule>

  # mark logs from default namespace for processing as k8s logs but stored as system logs
  <rule>
    key CONTAINER_NAME
    pattern ^k8s_[^_]+_[^_]+_default_
    tag kubernetes.journal.container._default_
  </rule>

  # mark logs from kube-* namespaces for processing as k8s logs but stored as system logs
  <rule>
    key CONTAINER_NAME
    pattern ^k8s_[^_]+_[^_]+_kube-(.+)_
    tag kubernetes.journal.container._kube-$1_
  </rule>

  # mark logs from openshift-* namespaces for processing as k8s logs but stored as system logs
  <rule>
    key CONTAINER_NAME
    pattern ^k8s_[^_]+_[^_]+_openshift-(.+)_
    tag kubernetes.journal.container._openshift-$1_
  </rule>

  # mark logs from openshift namespace for processing as k8s logs but stored as system logs
  <rule>
    key CONTAINER_NAME
    pattern ^k8s_[^_]+_[^_]+_openshift_
    tag kubernetes.journal.container._openshift_
  </rule>

  # mark fluentd container logs
  <rule>
    key CONTAINER_NAME
    pattern ^k8s_.*fluentd
    tag kubernetes.journal.container.fluentd
  </rule>

  # this is a kubernetes container
  <rule>
    key CONTAINER_NAME
    pattern ^k8s_
    tag kubernetes.journal.container
  </rule>

  # not kubernetes - assume a system log or system container log
  <rule>
    key _TRANSPORT
    pattern .+
    tag journal.system
  </rule>
</match>
{{- end}}
`

var KubernetesMetadataPlugin string = `
{{define "kubernetesMetadata"}}
# {{.Desc}}
<filter kubernetes.**>
  @type kubernetes_metadata
  kubernetes_url 'https://kubernetes.default.svc'
  cache_size '1000'
  watch 'false'
  use_journal 'nil'
  ssl_partial_chain 'true'
</filter>
{{- end}}
`

var ParseJsonFields string = `
{{define "parseJsonFields"}}
# {{.Desc}}
<filter kubernetes.journal.**>
  @type parse_json_field
  merge_json_log 'false'
  preserve_json_log 'true'
  json_fields 'log,MESSAGE'
</filter>

<filter kubernetes.var.log.containers.**>
  @type parse_json_field
  merge_json_log 'false'
  preserve_json_log 'true'
  json_fields 'log,MESSAGE'
</filter>

<filter kubernetes.var.log.containers.eventrouter-** kubernetes.var.log.containers.cluster-logging-eventrouter-**>
  @type parse_json_field
  merge_json_log true
  preserve_json_log true
  json_fields 'log,MESSAGE'
</filter>
{{- end}}
`

var CleanKibanaLogs string = `
{{define "cleanKibanaLogs"}}
# {{.Desc}}
<filter **kibana**>
  @type record_transformer
  enable_ruby
  <record>
    log ${record['err'] || record['msg'] || record['MESSAGE'] || record['log']}
  </record>
  remove_keys req,res,msg,name,level,v,pid,err
</filter>
{{- end}}
`

var FixAuditLevel string = `
{{define "fixAuditLevel"}}
# {{.Desc}}
<filter k8s-audit.log**>
  @type record_modifier
  <record>
    k8s_audit_level ${record['level']}
    level info
  </record>
</filter>
<filter openshift-audit.log**>
  @type record_modifier
  <record>
    openshift_audit_level ${record['level']}
    level info
  </record>
</filter>
{{end}}
`

var ViaQDataModel string = `
{{define "viaqDataModel" -}}
# {{.Desc}}
<filter **>
  @type viaq_data_model
  elasticsearch_index_prefix_field 'viaq_index_name'
  default_keep_fields CEE,time,@timestamp,aushape,ci_job,collectd,docker,fedora-ci,file,foreman,geoip,hostname,ipaddr4,ipaddr6,kubernetes,level,message,namespace_name,namespace_uuid,offset,openstack,ovirt,pid,pipeline_metadata,rsyslog,service,systemd,tags,testcase,tlog,viaq_msg_id
  extra_keep_fields ''
  keep_empty_fields 'message'
  use_undefined false
  undefined_name 'undefined'
  rename_time true
  rename_time_if_missing false
  src_time_name 'time'
  dest_time_name '@timestamp'
  pipeline_type 'collector'
  undefined_to_string 'false'
  undefined_dot_replace_char 'UNUSED'
  undefined_max_num_fields '-1'
  process_kubernetes_events 'false'
  <formatter>
    tag "system.var.log**"
    type sys_var_log
    remove_keys host,pid,ident
  </formatter>
  <formatter>
    tag "journal.system**"
    type sys_journal
    remove_keys log,stream,MESSAGE,_SOURCE_REALTIME_TIMESTAMP,__REALTIME_TIMESTAMP,CONTAINER_ID,CONTAINER_ID_FULL,CONTAINER_NAME,PRIORITY,_BOOT_ID,_CAP_EFFECTIVE,_CMDLINE,_COMM,_EXE,_GID,_HOSTNAME,_MACHINE_ID,_PID,_SELINUX_CONTEXT,_SYSTEMD_CGROUP,_SYSTEMD_SLICE,_SYSTEMD_UNIT,_TRANSPORT,_UID,_AUDIT_LOGINUID,_AUDIT_SESSION,_SYSTEMD_OWNER_UID,_SYSTEMD_SESSION,_SYSTEMD_USER_UNIT,CODE_FILE,CODE_FUNCTION,CODE_LINE,ERRNO,MESSAGE_ID,RESULT,UNIT,_KERNEL_DEVICE,_KERNEL_SUBSYSTEM,_UDEV_SYSNAME,_UDEV_DEVNODE,_UDEV_DEVLINK,SYSLOG_FACILITY,SYSLOG_IDENTIFIER,SYSLOG_PID
  </formatter>
  <formatter>
    tag "kubernetes.journal.container**"
    type k8s_journal
    remove_keys 'log,stream,MESSAGE,_SOURCE_REALTIME_TIMESTAMP,__REALTIME_TIMESTAMP,CONTAINER_ID,CONTAINER_ID_FULL,CONTAINER_NAME,PRIORITY,_BOOT_ID,_CAP_EFFECTIVE,_CMDLINE,_COMM,_EXE,_GID,_HOSTNAME,_MACHINE_ID,_PID,_SELINUX_CONTEXT,_SYSTEMD_CGROUP,_SYSTEMD_SLICE,_SYSTEMD_UNIT,_TRANSPORT,_UID,_AUDIT_LOGINUID,_AUDIT_SESSION,_SYSTEMD_OWNER_UID,_SYSTEMD_SESSION,_SYSTEMD_USER_UNIT,CODE_FILE,CODE_FUNCTION,CODE_LINE,ERRNO,MESSAGE_ID,RESULT,UNIT,_KERNEL_DEVICE,_KERNEL_SUBSYSTEM,_UDEV_SYSNAME,_UDEV_DEVNODE,_UDEV_DEVLINK,SYSLOG_FACILITY,SYSLOG_IDENTIFIER,SYSLOG_PID'
  </formatter>
  <formatter>
    tag "kubernetes.var.log.containers.eventrouter-** kubernetes.var.log.containers.cluster-logging-eventrouter-** k8s-audit.log** openshift-audit.log**"
    type k8s_json_file
    remove_keys log,stream,CONTAINER_ID_FULL,CONTAINER_NAME
    process_kubernetes_events 'true'
  </formatter>
  <formatter>
    tag "kubernetes.var.log.containers**"
    type k8s_json_file
    remove_keys log,stream,CONTAINER_ID_FULL,CONTAINER_NAME
  </formatter>
  <elasticsearch_index_name>
    enabled 'true'
    tag "journal.system** system.var.log** **_default_** **_kube-*_** **_openshift-*_** **_openshift_**"
    name_type static
    static_index_name infra-write
  </elasticsearch_index_name>
  <elasticsearch_index_name>
    enabled 'true'
    tag "linux-audit.log** k8s-audit.log** openshift-audit.log**"
    name_type static
    static_index_name audit-write
  </elasticsearch_index_name>
  <elasticsearch_index_name>
    enabled 'true'
    tag "**"
    name_type static
    static_index_name app-write
  </elasticsearch_index_name>
</filter>
{{end}}
`

var GenElasticsearchID string = `
{{define "genElasticsearchID" -}}
# {{.Desc}}
<filter **>
  @type elasticsearch_genid_ext
  hash_id_key viaq_msg_id
  alt_key kubernetes.event.metadata.uid
  alt_tags 'kubernetes.var.log.containers.logging-eventrouter-*.** kubernetes.var.log.containers.eventrouter-*.** kubernetes.var.log.containers.cluster-logging-eventrouter-*.** kubernetes.journal.container._default_.kubernetes-event'
</filter>
{{- end}}
`

var ConcatLines string = `
{{define "concatLines"}}
# {{.Desc}}
<filter kubernetes.**>
  @type concat
  key log
  partial_key logtag
  partial_value P
  separator ''
</filter>
{{- end}}
`
