// Package embedmappings ships the mapping JSON files inside the binary so
// users do not need to download or configure anything separately. The
// mapping files in data/ are the single source of truth for the project.
package embedmappings

import "embed"

//go:embed all:data
var FS embed.FS

// Dir is the directory name inside the embedded FS where mapping JSON files live.
const Dir = "data"
