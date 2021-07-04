package fluentd

var ToStdOut = `
{{define "toStdout"}}
# {{.Desc}}
<match {{.Pattern}}>
 @type stdout
</match>
{{end -}}`
