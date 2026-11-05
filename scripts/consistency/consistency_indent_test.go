//nolint:all
package main

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// indentFixtureDocs writes the given markdown body to a temporary
// website/docs-shaped file and parses it through parseResourceDocs, mirroring
// how the consistency gate reads a real resource doc. parseResourceDocs opens
// the file from docsPath+stem+".html.markdown" and recomputes the resource
// name from the file basename, so the stem must not contain a dot.
func indentFixtureDocs(t *testing.T, stem, body string) (map[string]ResourceAttribute, error) {
	t.Helper()
	dir := t.TempDir()
	filePath := filepath.Join(dir, stem+".html.markdown")
	if err := os.WriteFile(filePath, []byte(body), 0o644); err != nil {
		t.Fatalf("write fixture %s: %s", filePath, err)
	}
	attrs := make(map[string]ResourceAttribute)
	// parseResourceDocs joins docsPath+stem+".html.markdown"; the trailing
	// separator is required because it does not use filepath.Join.
	err := parseResourceDocs("alicloud_"+stem, dir+"/", true, attrs)
	return attrs, err
}

func indentAttrNames(attrs map[string]ResourceAttribute) []string {
	out := make([]string, 0, len(attrs))
	for k := range attrs {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// TestUnitIndentParseDocsTwoSpace is the core regression for the reported bug:
// a two-space nested Attributes Reference list. Before the fix, ascending from
// a grandchild (indent 4) back to a sibling (indent 2) left backIndex=2, which
// 2%4 != 0 rejected as "have not been formatted". The fix detects the
// two-space step and pops exactly one level, so every nested path is correct.
func TestUnitIndentParseDocsTwoSpace(t *testing.T) {
	body := "## Attributes Reference\n\n" +
		"The following attributes are exported:\n" +
		"* `id` - The ID.\n" +
		"* `endpoints` - List of endpoint details.\n" +
		"  * `computing_group_id` - The computing group ID.\n" +
		"  * `ports` - A list of port details.\n" +
		"    * `port` - The access port.\n" +
		"    * `protocol` - The protocol type.\n" +
		"  * `status` - Status.\n" +
		"  * `vpc_id` - The VPC ID.\n" +
		"* `engine_version` - The engine version.\n" +
		"## Timeouts\n"
	attrs, err := indentFixtureDocs(t, "fixture_two_space", body)
	if err != nil {
		t.Fatalf("parseResourceDocs failed for two-space fixture: %s", err)
	}
	want := []string{
		"endpoints",
		"endpoints.computing_group_id",
		"endpoints.ports",
		"endpoints.ports.port",
		"endpoints.ports.protocol",
		"endpoints.status",
		"endpoints.vpc_id",
		"engine_version",
		"id",
	}
	if got := indentAttrNames(attrs); !equalStrings(got, want) {
		t.Errorf("two-space paths =\n  %v\nwant\n  %v", got, want)
	}
}

// TestUnitIndentParseDocsFourSpace verifies legacy four-space nested lists are
// still parsed correctly after switching from a hard-coded /4 to a detected
// step. The same back-up arithmetic must work for both widths.
func TestUnitIndentParseDocsFourSpace(t *testing.T) {
	body := "## Attributes Reference\n\n" +
		"* `id` - The ID.\n" +
		"* `endpoints` - List of endpoint details.\n" +
		"    * `computing_group_id` - The computing group ID.\n" +
		"    * `ports` - A list of port details.\n" +
		"        * `port` - The access port.\n" +
		"        * `protocol` - The protocol type.\n" +
		"    * `status` - Status.\n" +
		"* `engine_version` - The engine version.\n" +
		"## Timeouts\n"
	attrs, err := indentFixtureDocs(t, "fixture_four_space", body)
	if err != nil {
		t.Fatalf("parseResourceDocs failed for four-space fixture: %s", err)
	}
	want := []string{
		"endpoints",
		"endpoints.computing_group_id",
		"endpoints.ports",
		"endpoints.ports.port",
		"endpoints.ports.protocol",
		"endpoints.status",
		"engine_version",
		"id",
	}
	if got := indentAttrNames(attrs); !equalStrings(got, want) {
		t.Errorf("four-space paths =\n  %v\nwant\n  %v", got, want)
	}
}

func TestUnitIndentParseDocsNonNestedInitialIndentFails(t *testing.T) {
	for _, tc := range []struct {
		name   string
		indent string
	}{
		{"one_space", " "},
		{"six_spaces", "      "},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// One space renders as a root sibling; six spaces renders
			// inside the parent paragraph (or as an indented code block).
			// Neither may satisfy a nested schema path root.child.
			body := "## Attributes Reference\n\n* `root` - desc.\n" +
				tc.indent + "* `child` - desc.\n* `sibling` - desc.\n## Timeouts\n"
			if _, err := indentFixtureDocs(t, "fixture_"+tc.name, body); err == nil {
				t.Fatal("non-nested Markdown must not be interpreted as a nested schema attribute")
			}
		})
	}
}

func TestUnitIndentParseDocsOtherValidIndentSteps(t *testing.T) {
	for _, tc := range []struct {
		name   string
		indent string
	}{
		{"three_spaces", "   "},
		{"five_spaces", "     "},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// These widths form nested lists in Markdown. Leave style
			// enforcement to MD007 while preserving their schema paths.
			body := "## Attributes Reference\n\n* `root` - desc.\n" +
				tc.indent + "* `child` - desc.\n" +
				tc.indent + tc.indent + "* `grandchild` - desc.\n" +
				tc.indent + "* `peer` - desc.\n* `sibling` - desc.\n## Timeouts\n"
			attrs, err := indentFixtureDocs(t, "fixture_"+tc.name, body)
			if err != nil {
				t.Fatalf("valid nested list was rejected: %s", err)
			}
			want := []string{"root", "root.child", "root.child.grandchild", "root.peer", "sibling"}
			if got := indentAttrNames(attrs); !equalStrings(got, want) {
				t.Fatalf("nested paths = %v; want %v", got, want)
			}
		})
	}
}

// TestUnitIndentParseDocsMultipleRootsAndAncestorFallback exercises several
// root fields, a grandchild that falls back past its parent to a sibling, and
// a deep-to-root ancestor fallback. This is the "multiple root fields" and
// "sibling/ancestor fallback" coverage required by the task.
func TestUnitIndentParseDocsMultipleRootsAndAncestorFallback(t *testing.T) {
	body := "## Attributes Reference\n\n" +
		"* `root1` - desc.\n" +
		"  * `child1` - desc.\n" +
		"    * `grandchild1` - desc.\n" +
		"  * `child2` - desc.\n" +
		"* `root2` - desc.\n" +
		"  * `child3` - desc.\n" +
		"* `root3` - desc.\n" +
		"## Timeouts\n"
	attrs, err := indentFixtureDocs(t, "fixture_roots", body)
	if err != nil {
		t.Fatalf("parseResourceDocs failed: %s", err)
	}
	want := []string{
		"root1",
		"root1.child1",
		"root1.child1.grandchild1",
		"root1.child2",
		"root2",
		"root2.child3",
		"root3",
	}
	if got := indentAttrNames(attrs); !equalStrings(got, want) {
		t.Errorf("roots paths =\n  %v\nwant\n  %v", got, want)
	}
}

// TestUnitIndentParseDocsBlockScopeResets verifies that a `### `block“ heading
// opens a fresh scope: the indentation context of the preceding nested list
// (depth, current field, detected step) does not bleed into the block's own
// list. Without the scope reset, the block's first root-level field would be
// treated as a back-up from the previous nested indent and mis-parsed.
func TestUnitIndentParseDocsBlockScopeResets(t *testing.T) {
	body := "## Attributes Reference\n\n" +
		"* `top` - desc.\n" +
		"  * `child` - desc.\n" +
		"### `block`\n\n" +
		"The block supports the following:\n" +
		"* `field1` - desc.\n" +
		"  * `field2` - desc.\n" +
		"## Timeouts\n"
	attrs, err := indentFixtureDocs(t, "fixture_block", body)
	if err != nil {
		t.Fatalf("parseResourceDocs failed: %s", err)
	}
	want := []string{
		"block.field1",
		"block.field1.field2",
		"top",
		"top.child",
	}
	if got := indentAttrNames(attrs); !equalStrings(got, want) {
		t.Errorf("block-scope paths =\n  %v\nwant\n  %v", got, want)
	}
}

// TestUnitIndentParseDocsMalformedDescentFails: in a two-space list, a field
// indented by a non-step delta (3 spaces, not a multiple of 2) is malformed
// and must still be rejected — the fix does not relax real formatting errors.
func TestUnitIndentParseDocsMalformedDescentFails(t *testing.T) {
	body := "## Attributes Reference\n\n" +
		"* `root1` - desc.\n" +
		"  * `child1` - desc.\n" +
		"   * `malformed` - desc.\n" +
		"## Timeouts\n"
	_, err := indentFixtureDocs(t, "fixture_malformed_descent", body)
	if err == nil {
		t.Fatal("parseResourceDocs should fail for a non-step indentation delta")
	}
}

// TestUnitIndentParseDocsInconsistentStepFails: a two-space list that suddenly
// widens by four (mixed indentation) is inconsistent and must fail.
func TestUnitIndentParseDocsInconsistentStepFails(t *testing.T) {
	body := "## Attributes Reference\n\n" +
		"* `root1` - desc.\n" +
		"  * `child1` - desc.\n" +
		"      * `mixed` - desc.\n" +
		"## Timeouts\n"
	_, err := indentFixtureDocs(t, "fixture_inconsistent_step", body)
	if err == nil {
		t.Fatal("parseResourceDocs should fail for an inconsistent indentation step")
	}
}

// TestUnitIndentParseDocsNonMultipleBackDeltaFails: in a four-space list,
// ascending by a delta that is not a whole multiple of the step (e.g. back by
// 6 from indent 8 to indent 2) is malformed and must fail.
func TestUnitIndentParseDocsNonMultipleBackDeltaFails(t *testing.T) {
	body := "## Attributes Reference\n\n" +
		"* `root1` - desc.\n" +
		"    * `child1` - desc.\n" +
		"        * `grandchild1` - desc.\n" +
		"  * `malformed` - desc.\n" +
		"## Timeouts\n"
	_, err := indentFixtureDocs(t, "fixture_nonmultiple_back", body)
	if err == nil {
		t.Fatal("parseResourceDocs should fail for a non-multiple back-up delta")
	}
}

// schemaForTwoSpaceFixture mirrors the two-space fixture's parsed paths as a
// nested SDKv2 schema, used for end-to-end consistencyCheck assertions.
func schemaForTwoSpaceFixture() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id":             {Type: schema.TypeString},
		"engine_version": {Type: schema.TypeString},
		"endpoints": {
			Type: schema.TypeList,
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				"computing_group_id": {Type: schema.TypeString},
				"ports": {
					Type: schema.TypeList,
					Elem: &schema.Resource{Schema: map[string]*schema.Schema{
						"port":     {Type: schema.TypeString},
						"protocol": {Type: schema.TypeString},
					}},
				},
				"status": {Type: schema.TypeString},
				"vpc_id": {Type: schema.TypeString},
			}},
		},
	}
}

// TestUnitIndentConsistencyCheckMatchingPasses is the end-to-end proof that a
// two-space nested doc — the shape that triggered the original false failure —
// now passes consistencyCheck against a matching schema.
func TestUnitIndentConsistencyCheckMatchingPasses(t *testing.T) {
	body := "## Attributes Reference\n\n" +
		"* `id` - The ID.\n" +
		"* `endpoints` - List of endpoint details.\n" +
		"  * `computing_group_id` - The computing group ID.\n" +
		"  * `ports` - A list of port details.\n" +
		"    * `port` - The access port.\n" +
		"    * `protocol` - The protocol type.\n" +
		"  * `status` - Status.\n" +
		"  * `vpc_id` - The VPC ID.\n" +
		"* `engine_version` - The engine version.\n" +
		"## Timeouts\n"
	attrs, err := indentFixtureDocs(t, "fixture_consistency_pass", body)
	if err != nil {
		t.Fatalf("parseResourceDocs failed: %s", err)
	}
	if !consistencyCheck("alicloud_fixture_consistency_pass", attrs, schemaForTwoSpaceFixture()) {
		t.Error("consistencyCheck should PASS for a two-space nested doc with a matching schema")
	}
}

// TestUnitIndentConsistencyCheckMissingFieldFails: a schema field absent from
// the docs must still fail — the fix does not mask real gaps.
func TestUnitIndentConsistencyCheckMissingFieldFails(t *testing.T) {
	body := "## Attributes Reference\n\n" +
		"* `id` - The ID.\n" +
		"* `endpoints` - List of endpoint details.\n" +
		"  * `computing_group_id` - The computing group ID.\n" +
		"  * `status` - Status.\n" +
		"* `engine_version` - The engine version.\n" +
		"## Timeouts\n"
	attrs, err := indentFixtureDocs(t, "fixture_missing_field", body)
	if err != nil {
		t.Fatalf("parseResourceDocs failed: %s", err)
	}
	// Schema declares endpoints.vpc_id which the docs omit.
	sch := schemaForTwoSpaceFixture()
	if consistencyCheck("alicloud_fixture_missing_field", attrs, sch) {
		t.Error("consistencyCheck should FAIL when a schema field is missing from the docs")
	}
}

// TestUnitIndentConsistencyCheckExtraDocFieldFails: a field described in the
// docs but absent from the schema must still fail.
func TestUnitIndentConsistencyCheckExtraDocFieldFails(t *testing.T) {
	body := "## Attributes Reference\n\n" +
		"* `id` - The ID.\n" +
		"* `endpoints` - List of endpoint details.\n" +
		"  * `computing_group_id` - The computing group ID.\n" +
		"  * `ports` - A list of port details.\n" +
		"    * `port` - The access port.\n" +
		"    * `protocol` - The protocol type.\n" +
		"  * `status` - Status.\n" +
		"  * `vpc_id` - The VPC ID.\n" +
		"  * `ghost_field` - Not in schema.\n" +
		"* `engine_version` - The engine version.\n" +
		"## Timeouts\n"
	attrs, err := indentFixtureDocs(t, "fixture_extra_field", body)
	if err != nil {
		t.Fatalf("parseResourceDocs failed: %s", err)
	}
	if consistencyCheck("alicloud_fixture_extra_field", attrs, schemaForTwoSpaceFixture()) {
		t.Error("consistencyCheck should FAIL when the docs describe a field absent from the schema")
	}
}

// TestUnitIndentConsistencyCheckOptionalMismatchFails: a schema Optional
// attribute that the doc does not mark (Optional) must still fail — the
// modifier semantics are unchanged by the indentation fix.
func TestUnitIndentConsistencyCheckOptionalMismatchFails(t *testing.T) {
	body := "## Attributes Reference\n\n" +
		"* `id` - The ID.\n" +
		"* `endpoints` - List of endpoint details.\n" +
		"  * `computing_group_id` - The computing group ID.\n" +
		"  * `ports` - A list of port details.\n" +
		"    * `port` - The access port.\n" +
		"    * `protocol` - The protocol type.\n" +
		"  * `status` - Status.\n" +
		"  * `vpc_id` - The VPC ID.\n" +
		"* `engine_version` - The engine version.\n" +
		"## Timeouts\n"
	attrs, err := indentFixtureDocs(t, "fixture_optional_mismatch", body)
	if err != nil {
		t.Fatalf("parseResourceDocs failed: %s", err)
	}
	sch := schemaForTwoSpaceFixture()
	// Schema marks the nested port Optional, but the doc line carries no
	// (Optional) modifier, so consistencyCheck must report the mismatch.
	sch["endpoints"].Elem.(*schema.Resource).Schema["ports"].Elem.(*schema.Resource).Schema["port"].Optional = true
	if consistencyCheck("alicloud_fixture_optional_mismatch", attrs, sch) {
		t.Error("consistencyCheck should FAIL when schema Optional is unmarked in the docs")
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
