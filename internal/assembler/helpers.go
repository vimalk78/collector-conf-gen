package assembler

import (
	"fmt"
	"sort"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var replacer = strings.NewReplacer(" ", "_", "-", "_", ".", "_")

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
