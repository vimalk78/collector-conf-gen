package logging

var ToStdOut = `
{{define "toStdout"}}
# {{.Desc}}
<match {{.Pattern}}>
 @type stdout
</match>
{{end -}}`
