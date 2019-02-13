package main

const templateSample1 = `## Projects
{{range .}}
### {{.Name}}
` + "`" + "`" + "`" + `
{{.Body}}
` + "`" + "`" + "`" + `

{{end}}`
