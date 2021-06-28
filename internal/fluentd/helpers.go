package fluentd

import (
	"fmt"
	"sort"
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
	//fmt.Printf("map: %#v\n", m)
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	kv := make([]string, len(m))
	for i, k := range keys {
		//kv = append(kv, fmt.Sprintf("%s:%s", k, m[k]))
		kv[i] = fmt.Sprintf("%s:%s", k, m[k])
	}
	return kv
}
