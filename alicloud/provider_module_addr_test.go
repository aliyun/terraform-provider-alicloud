package alicloud

import (
	"os"
	"path/filepath"
	"testing"
)

// TestGetModuleAddrFromFile verifies that getModuleAddrFromFile handles every
// structural anomaly in .terraform/modules/modules.json gracefully (no panic,
// empty result) while still extracting registry module addresses for the
// well-formed case. This guards provider configure against a panic that
// previously arose from unsafe type assertions on the Modules field.
func TestGetModuleAddrFromFile(t *testing.T) {
	dir := t.TempDir()
	modPath := filepath.Join(dir, ".terraform", "modules", "modules.json")

	write := func(t *testing.T, content string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(modPath), 0o755); err != nil {
			t.Fatalf("mkdir modules dir: %v", err)
		}
		if err := os.WriteFile(modPath, []byte(content), 0o644); err != nil {
			t.Fatalf("write modules.json: %v", err)
		}
	}
	// removeModules makes the file absent so the "file missing" path is exercised.
	removeModules := func() { _ = os.Remove(modPath) }

	cases := []struct {
		name    string
		content string
		absent  bool
		want    string
	}{
		{
			name:    "well-formed registry modules",
			content: `{"Modules":[{"Source":"registry.terraform.io/aliyun/alicloud/1.209.0","Version":"1.209.0"},{"Source":"registry.terraform.io/aliyun/vpc/2.0.0","Version":"2.0.0"}]}`,
			// parts of first = ["registry.terraform.io","aliyun","alicloud","1.209.0"]
			//   -> "terraform-{parts[3]}-{parts[2]}/{version}" = "terraform-1.209.0-alicloud/1.209.0"
			// second -> "terraform-2.0.0-vpc/2.0.0", each prefixed with a space.
			want: " terraform-1.209.0-alicloud/1.209.0 terraform-2.0.0-vpc/2.0.0",
		},
		{
			// "Modules" is an object instead of an array: previously panicked on
			// the .([]interface{}) assertion; now a typed decode error, no panic.
			name:    "modules object instead of array",
			content: `{"Modules":{}}`,
			want:    "",
		},
		{
			// null array element: previously panicked on m.(map[string]interface{});
			// now decodes to a zero-value moduleRecord whose Source is "" -> skipped.
			name:    "null module element",
			content: `{"Modules":[null]}`,
			want:    "",
		},
		{
			// number where an object is expected: typed decode error, no panic.
			name:    "number module element",
			content: `{"Modules":[123]}`,
			want:    "",
		},
		{
			// wrong field type (Source as number): typed decode error, no panic.
			name:    "source wrong type number",
			content: `{"Modules":[{"Source":123}]}`,
			want:    "",
		},
		{
			name:    "empty object",
			content: `{}`,
			want:    "",
		},
		{
			name:    "modules key null",
			content: `{"Modules":null}`,
			want:    "",
		},
		{
			name:    "truncated json",
			content: `{"Modules":[{"Source":"reg`,
			want:    "",
		},
		{
			name:    "empty file",
			content: ``,
			want:    "",
		},
		{
			name:    "non-registry source ignored",
			content: `{"Modules":[{"Source":"github.com/foo/bar","Version":"1.0.0"}]}`,
			want:    "",
		},
		{
			// registry path with != 4 segments is ignored (preserves prior behavior).
			name:    "registry path with too many segments",
			content: `{"Modules":[{"Source":"registry.terraform.io/aliyun/alicloud/sub/x/1.0.0","Version":"1.0.0"}]}`,
			want:    "",
		},
		{
			name:   "file absent",
			absent: true,
			want:   "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.absent {
				removeModules()
			} else {
				write(t, tc.content)
			}
			got := getModuleAddrFromFile(modPath)
			if got != tc.want {
				t.Fatalf("getModuleAddrFromFile = %q, want %q", got, tc.want)
			}
		})
	}
}
