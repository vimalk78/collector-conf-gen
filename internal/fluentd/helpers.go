package fluentd

import (
	"fmt"
	"strings"
	"text/template"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	helperRegistry = &template.FuncMap{
		"applicationTag":      applicationTag,
		"labelName":           labelName,
		"sourceTypelabelName": sourceTypeLabelName,
		"routeMapValues":      routeMapValues,
	}
)
var replacer = strings.NewReplacer(" ", "_", "-", "_", ".", "_")

func applicationTag(namespace string) string {
	if "" == namespace {
		return "**"
	}
	return strings.ToLower(fmt.Sprintf("kubernetes.**_%s_**", namespace))
}

func labelName(name string) string {
	return strings.ToUpper(fmt.Sprintf("@%s", replacer.Replace(name)))
}

func sourceTypeLabelName(name string) string {
	return strings.ToUpper(fmt.Sprintf("@_%s", replacer.Replace(name)))
}

func routeMapValues(routeMap logging.RouteMap, key string) []string {
	if values, found := routeMap[key]; found {
		return values.List()
	}
	return []string{}
}

func comma_separated(arr []string) string {
	return strings.Join(arr, ",")
}

func LabelsKV(ls *metav1.LabelSelector) []string {
	m, _ := metav1.LabelSelectorAsMap(ls)
	kv := []string{}
	for k, v := range m {
		kv = append(kv, fmt.Sprintf("%s:%s", k, v))
	}
	return kv
}
