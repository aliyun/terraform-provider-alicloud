//nolint:all
package main

import (
	"os"
	"path/filepath"
	"testing"
)

// structureFixtureDocs parses a docs body through parseResourceDocs with the
// given flattened schema keys, mirroring how the gate reads a real page.
func structureFixtureDocs(t *testing.T, stem, body string, schemaKeys map[string]struct{}) (map[string]ResourceAttribute, error) {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, stem+".html.markdown"), []byte(body), 0o644); err != nil {
		t.Fatalf("write fixture: %s", err)
	}
	attrs := make(map[string]ResourceAttribute)
	_, err := parseResourceDocs("alicloud_"+stem, dir+"/", true, attrs, schemaKeys)
	return attrs, err
}

func keySet(keys ...string) map[string]struct{} {
	out := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		out[key] = struct{}{}
	}
	return out
}

// TestUnitHeadingOwnsScope covers the ownership rule that separates a real
// nested block heading from a documentation-only grouping heading.
func TestUnitHeadingOwnsScope(t *testing.T) {
	keys := keySet("load_balancer_id", "x_forwarded_for", "x_forwarded_for.retrive_slb_ip")
	cases := []struct {
		name string
		keys map[string]struct{}
		want bool
	}{
		{"http", keys, false},                          // grouping heading on a flat schema
		{"x_forwarded_for", keys, true},                // a real attribute
		{"x_forwarded_for.retrive_slb_ip", keys, true}, // a real nested path
		{"", keys, false},                              // empty heading name
		{"http", nil, true},                            // unknown schema keeps the legacy assumption
	}
	for _, c := range cases {
		if got := headingOwnsScope(c.name, c.keys); got != c.want {
			t.Errorf("headingOwnsScope(%q) = %v, want %v", c.name, got, c.want)
		}
	}
}

// TestUnitGroupingHeadingDoesNotPrefix is the regression for the slb_listener
// false positives: `### `http“ groups fields of a flat schema, so its bullets
// must stay top-level instead of turning into `http.frontend_port`.
func TestUnitGroupingHeadingDoesNotPrefix(t *testing.T) {
	body := "## Argument Reference\n\n" +
		"* `HTTP Listener` - See [`http`](#http) below.\n\n" +
		"---\n\n" +
		"### `http`\n\n" +
		"The HTTP Listener supports the following:\n\n" +
		"* `frontend_port` - (Required, ForceNew) The frontend port.\n" +
		"* `protocol` - (Required, ForceNew) The protocol.\n" +
		"## Attributes Reference\n"
	attrs, err := structureFixtureDocs(t, "fixture_grouping", body, keySet("frontend_port", "protocol"))
	if err != nil {
		t.Fatalf("parseResourceDocs failed: %s", err)
	}
	if got := indentAttrNames(attrs); !equalStrings(got, []string{"frontend_port", "protocol"}) {
		t.Errorf("grouping heading paths = %v, want [frontend_port protocol]", got)
	}
}

// TestUnitAttributeHeadingStillPrefixes guards the other half of the ownership
// rule: a heading that does name a schema attribute keeps nesting its bullets.
func TestUnitAttributeHeadingStillPrefixes(t *testing.T) {
	body := "## Argument Reference\n\n" +
		"### `x_forwarded_for`\n\n" +
		"The mapping supports the following:\n\n" +
		"* `retrive_slb_ip` - (Optional) Whether the SLB-IP header is used.\n" +
		"## Attributes Reference\n"
	keys := keySet("x_forwarded_for", "x_forwarded_for.retrive_slb_ip")
	attrs, err := structureFixtureDocs(t, "fixture_attribute_heading", body, keys)
	if err != nil {
		t.Fatalf("parseResourceDocs failed: %s", err)
	}
	if got := indentAttrNames(attrs); !equalStrings(got, []string{"x_forwarded_for.retrive_slb_ip"}) {
		t.Errorf("attribute heading paths = %v, want [x_forwarded_for.retrive_slb_ip]", got)
	}
}

// TestUnitFrameworkNestedAttributeHeadingStillPrefixes reproduces the Framework
// path losing nested scopes when schema keys are truncated to their root names.
func TestUnitFrameworkNestedAttributeHeadingStillPrefixes(t *testing.T) {
	schemaAttributes := map[string]ResourceAttribute{
		"instances":          {Name: "instances"},
		"instances.tags":     {Name: "instances.tags"},
		"instances.tags.key": {Name: "instances.tags.key"},
	}
	body := "## Attributes Reference\n\n" +
		"### `instances.tags`\n\n" +
		"The tags support the following:\n\n" +
		"* `key` - The key of the tag.\n" +
		"## Timeouts\n"
	attrs, err := structureFixtureDocs(t, "fixture_framework_nested_heading", body, frameworkSchemaKeySet(schemaAttributes))
	if err != nil {
		t.Fatalf("parseResourceDocs failed: %s", err)
	}
	if got := indentAttrNames(attrs); !equalStrings(got, []string{"instances.tags.key"}) {
		t.Errorf("Framework nested heading paths = %v, want [instances.tags.key]", got)
	}
}

// TestUnitSpacedBulletIsNotAttribute pins the name contract: a backticked
// grouping or see-also bullet is neither a documented attribute nor a schema
// attribute, so it must not enter the docs map in either direction. Before the
// wider capture the line simply fell through both regexes; the explicit
// rejection keeps that outcome and makes it testable.
func TestUnitSpacedBulletIsNotAttribute(t *testing.T) {
	// The production attributeFieldRegex captures "Node Pool" as the name; the
	// contract check is what rejects it.
	words := []string{"  * `Node Pool` - See [`node_pool`](#node_pool) below.", "Node Pool", " See [`node_pool`](#node_pool) below."}
	if got := parseMatchLine(words, "Attribute", ""); got != nil {
		t.Errorf("spaced name must not parse as an attribute, got %+v", got)
	}
	if got := parseMatchLine(nil, "Argument", ""); got != nil {
		t.Errorf("nil match must not parse as an attribute, got %+v", got)
	}
	body := "## Attributes Reference\n\n" +
		"* `id` - The ID.\n" +
		"  * `Node Pool` - See [`node_pool`](#node_pool) below.\n" +
		"## Timeouts\n"
	attrs, err := structureFixtureDocs(t, "fixture_spaced_bullet", body, keySet("id", "node_pool"))
	if err != nil {
		t.Fatalf("parseResourceDocs failed: %s", err)
	}
	for name := range attrs {
		if name != "id" {
			t.Errorf("spaced bullet leaked into the docs map as %q", name)
		}
	}
}

// TestUnitSeeAlsoBulletStaysOutOfDocsMap checks the same contract end to end:
// the see-also lines that introduce a per-protocol section contribute no
// attribute, while the section body still contributes its fields.
func TestUnitSeeAlsoBulletStaysOutOfDocsMap(t *testing.T) {
	body := "## Argument Reference\n\n" +
		"The following arguments are supported:\n\n" +
		"* `HTTP Listener` - See [`http`](#http) below.\n\n" +
		"### `http`\n\n" +
		"* `frontend_port` - (Required) The frontend port.\n" +
		"## Attributes Reference\n"
	attrs, err := structureFixtureDocs(t, "fixture_see_also", body, keySet("frontend_port"))
	if err != nil {
		t.Fatalf("parseResourceDocs failed: %s", err)
	}
	if _, ok := attrs["HTTP Listener"]; ok {
		t.Error("see-also bullet leaked into the docs map")
	}
	if _, ok := attrs["frontend_port"]; !ok {
		t.Errorf("section body lost, docs map = %v", indentAttrNames(attrs))
	}
}

// TestUnitMissingDocsPageReturnsError is the unit-level proof for the silent
// skip this gate used to allow: an unreadable page must surface as an error so
// checkName can fail the run instead of passing it.
func TestUnitMissingDocsPageReturnsError(t *testing.T) {
	dir := t.TempDir()
	attrs := make(map[string]ResourceAttribute)
	_, err := parseResourceDocs("alicloud_absent_page", dir+"/", true, attrs, keySet("id"))
	if err == nil {
		t.Fatal("parseResourceDocs must report a missing docs page")
	}
}

// structureFixtureDuplicates parses a docs body and returns the attributes the
// page documents in both reference sections.
func structureFixtureDuplicates(t *testing.T, stem, body string, schemaKeys map[string]struct{}) []string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, stem+".html.markdown"), []byte(body), 0o644); err != nil {
		t.Fatalf("write fixture: %s", err)
	}
	attrs := make(map[string]ResourceAttribute)
	duplicates, err := parseResourceDocs("alicloud_"+stem, dir+"/", true, attrs, schemaKeys)
	if err != nil {
		t.Fatalf("parseResourceDocs failed: %s", err)
	}
	return duplicates
}

// TestUnitArgumentSectionRepeatsAreNotDuplicates is the regression for the
// slb_listener false positives: the flat fields of a per-protocol grouping are
// listed once per protocol inside the Argument Reference, which is not an input
// being re-documented as an output.
func TestUnitArgumentSectionRepeatsAreNotDuplicates(t *testing.T) {
	body := "## Argument Reference\n\n" +
		"### `http`\n\n" +
		"* `frontend_port` - (Required, ForceNew) The frontend port.\n" +
		"### `https`\n\n" +
		"* `frontend_port` - (Required, ForceNew) The frontend port.\n" +
		"## Attributes Reference\n"
	if got := structureFixtureDuplicates(t, "fixture_repeat_argument", body, keySet("frontend_port")); len(got) != 0 {
		t.Errorf("repeats inside the Argument Reference = %v, want none", got)
	}
}

// TestUnitAttributesRepeatIsDuplicate keeps the real check: a field already
// documented as an argument and listed again as an output is reported.
func TestUnitAttributesRepeatIsDuplicate(t *testing.T) {
	body := "## Argument Reference\n\n" +
		"* `instance_name` - (Optional) The name.\n" +
		"## Attributes Reference\n\n" +
		"* `id` - The ID.\n" +
		"* `instance_name` - The name.\n" +
		"## Timeouts\n"
	got := structureFixtureDuplicates(t, "fixture_repeat_attributes", body, keySet("instance_name"))
	if len(got) != 1 || got[0] != "instance_name" {
		t.Errorf("duplicates = %v, want [instance_name]", got)
	}
}

// TestUnitDeletedFileEntities covers the deletion derivation: a removed page
// names the entity it documented.
func TestUnitDeletedFileEntities(t *testing.T) {
	if got := entityFromDocsPage("website/docs/r/slb_listener.html.markdown"); got != "alicloud_slb_listener" {
		t.Errorf("entityFromDocsPage = %q, want alicloud_slb_listener", got)
	}
	if got := entityFromDocsPage("website/docs/d/elasticsearch_instances.html.markdown"); got != "alicloud_elasticsearch_instances" {
		t.Errorf("entityFromDocsPage = %q, want alicloud_elasticsearch_instances", got)
	}
}
