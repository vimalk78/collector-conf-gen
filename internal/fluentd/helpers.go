package fluentd

import (
	"fmt"
	"sort"
	"strings"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func labelNames(names []string) []string {
	asLabels := make([]string, len(names))
	for i, n := range names {
		asLabels[i] = labelName(n)
	}
	return asLabels
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
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	kv := make([]string, len(m))
	for i, k := range keys {
		kv[i] = fmt.Sprintf("%s:%s", k, m[k])
	}
	return kv
}
