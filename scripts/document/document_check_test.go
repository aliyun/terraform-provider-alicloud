package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func validResourceDocument() []string {
	return []string{
		"---",
		`layout: "alicloud"`,
		`page_title: "Alicloud: alicloud_demo"`,
		"description: |-",
		"  Provides a Alicloud Demo resource.",
		"---",
		"",
		"# alicloud_demo",
		"",
		"-> **NOTE:** Available since v1.0.0.",
		"",
		"## Example Usage",
		"",
		"```terraform",
		`resource "alicloud_demo" "example" {`,
		`  name = "example"`,
		"}",
		"```",
		"",
		"## Argument Reference",
		"",
		"The following arguments are supported:",
		"",
		"* `name` - (Required) The name.",
		"",
		"## Attributes Reference",
		"",
		"* `id` - The ID of the resource.",
		"",
		"## Timeouts",
		"",
		"## Import",
		"",
		"```shell",
		"terraform import alicloud_demo.example <id>",
		"```",
	}
}

func writeResourceDocument(t *testing.T, lines []string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "demo.html.markdown")
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
		t.Fatalf("write document: %v", err)
	}
	return path
}

func TestUnitDocumentCheckAcceptsValidResourcePage(t *testing.T) {
	if code := checkDocument(writeResourceDocument(t, validResourceDocument())); code != 0 {
		t.Fatalf("valid resource page rejected with exit code %d", code)
	}
}

func TestUnitDocumentCheckRejectsMissingVersionNote(t *testing.T) {
	var document []string
	for _, line := range validResourceDocument() {
		if strings.HasPrefix(line, "-> **NOTE:** Available since v") {
			continue
		}
		document = append(document, line)
	}
	if code := checkDocument(writeResourceDocument(t, document)); code != 1 {
		t.Fatalf("page without a version note returned exit code %d, want 1", code)
	}
}

func TestUnitDocumentCheckRejectsNonTerraformExampleFence(t *testing.T) {
	document := validResourceDocument()
	for i, line := range document {
		if line == "```terraform" {
			document[i] = "```hcl"
			break
		}
	}
	if code := checkDocument(writeResourceDocument(t, document)); code != 1 {
		t.Fatalf("page with a non-terraform example fence returned exit code %d, want 1", code)
	}
}
