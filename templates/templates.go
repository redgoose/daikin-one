package templates

import "embed"

//go:embed tmpl/*.tmpl
var TemplatesFS embed.FS
