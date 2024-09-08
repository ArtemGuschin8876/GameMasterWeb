package static

import "embed"

//go:embed ui/html/*.html
var FS embed.FS
