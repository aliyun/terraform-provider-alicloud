//nolint:all
package main

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// testAttrFieldRegex mirrors the package-internal attributeFieldRegex used by
// readDoc for the Attribute phase. It is duplicated here (rather than reaching
// into readDoc's local declarations) so the test exercises the same doc-line
// shape that consistencyCheck sees in production.
var testAttrFieldRegex = regexp.MustCompile("^\\s*\\* `([a-zA-Z_0-9]*)`[ ]*-?(.*)")

func attrWords(line string) []string {
	m := testAttrFieldRegex.FindStringSubmatch(line)
	if m == nil {
		return nil
	}
	return m
}

// TestUnitParseMatchLineAttributeOptional reproduces the false positive observed
// on alicloud_nlb_load_balancers.balancers.tags: the schema marks the nested
// attribute Optional, and the doc correctly writes "(Optional)". Before the
// fix, the Attribute phase only parsed Deprecated, so attributeDocsValue.Optional
// was always "" and consistencyCheck reported "should be marked as Optional in
// the document" — a false positive, since the doc DOES mark it.
func TestUnitParseMatchLineAttributeOptional(t *testing.T) {
	words := attrWords("    * `tags` - (Optional) The tag of the resource.")
	if words == nil {
		t.Fatal("attribute line did not match")
	}
	got := parseMatchLine(words, "Attribute", "balancers")
	if got == nil {
		t.Fatal("parseMatchLine returned nil")
	}
	if got.Name != "balancers.tags" {
		t.Errorf("Name = %q, want %q", got.Name, "balancers.tags")
	}
	if got.Optional != "true" {
		t.Errorf("Optional = %q, want %q (doc marks (Optional))", got.Optional, "true")
	}
	if got.Required != "" || got.ForceNew || got.Deprecated != "" {
		t.Errorf("only Optional expected, got Required=%q ForceNew=%v Deprecated=%q", got.Required, got.ForceNew, got.Deprecated)
	}
}

// TestUnitParseMatchLineAttributeRequiredForceNew verifies Required and ForceNew
// are now recognized in the Attribute phase, mirroring the Argument phase.
func TestUnitParseMatchLineAttributeRequiredForceNew(t *testing.T) {
	words := attrWords("    * `id` - (Required, ForceNew) The identifier.")
	got := parseMatchLine(words, "Attribute", "")
	if got == nil {
		t.Fatal("parseMatchLine returned nil")
	}
	if got.Required != "true" {
		t.Errorf("Required = %q, want %q", got.Required, "true")
	}
	if !got.ForceNew {
		t.Errorf("ForceNew = false, want true")
	}
	if got.Optional != "" {
		t.Errorf("Optional = %q, want empty", got.Optional)
	}
}

// TestUnitParseMatchLineAttributeDeprecatedPreserved is a regression guard: the
// fix adds Optional/Required/ForceNew parsing but must keep Deprecated working.
func TestUnitParseMatchLineAttributeDeprecatedPreserved(t *testing.T) {
	words := attrWords("    * `old_field` - (Deprecated) use `new_field` instead.")
	got := parseMatchLine(words, "Attribute", "")
	if got == nil {
		t.Fatal("parseMatchLine returned nil")
	}
	if got.Deprecated != "Deprecated since" {
		t.Errorf("Deprecated = %q, want %q", got.Deprecated, "Deprecated since")
	}
	if got.Optional != "" || got.Required != "" || got.ForceNew {
		t.Errorf("Deprecated-only line should not set Optional/Required/ForceNew, got %+v", got)
	}
}

// TestUnitParseMatchLineAttributeNoModifier verifies a plain Computed-only output
// attribute (no parenthesized modifier) sets none of the markers.
func TestUnitParseMatchLineAttributeNoModifier(t *testing.T) {
	words := attrWords("    * `status` - The status of the NLB instance.")
	got := parseMatchLine(words, "Attribute", "")
	if got == nil {
		t.Fatal("parseMatchLine returned nil")
	}
	if got.Optional != "" || got.Required != "" || got.ForceNew || got.Deprecated != "" {
		t.Errorf("plain output attribute should have no modifiers, got %+v", got)
	}
}

// TestUnitParseMatchLineAttributeRootName covers nested attribute naming when a
// rootName prefix is supplied (the balancers.tags case with root active).
func TestUnitParseMatchLineAttributeRootName(t *testing.T) {
	words := attrWords("    * `tags` - (Optional) The tag of the resource.")
	got := parseMatchLine(words, "Attribute", "balancers")
	if got == nil || got.Name != "balancers.tags" {
		t.Fatalf("Name = %v, want %q", got, "balancers.tags")
	}
}

// TestUnitParseMatchLineAttributeProseNotModifier guards the regression flagged in
// review: in the Attribute phase words[2] is the whole description tail, so a
// plain strings.Contains would mistake free-form prose for a modifier. The
// fix matches Optional/Required/ForceNew only inside a parenthesized token.
// Each line below has no parentheses and must not set any modifier.
func TestUnitParseMatchLineAttributeProseNotModifier(t *testing.T) {
	cases := []string{
		"    * `status` - Optional metadata returned by the service.",
		"    * `id` - Required for subsequent API calls.",
		"    * `field` - ForceNew is not supported for this field.",
		"    * `field` - Use Optional when configuring the block.",
		"    * `field` - This is Required for the resource to work.",
		"    * `field` - The ForceNew behavior is documented elsewhere.",
	}
	for _, line := range cases {
		words := attrWords(line)
		if words == nil {
			t.Errorf("line did not match regex: %q", line)
			continue
		}
		got := parseMatchLine(words, "Attribute", "")
		if got == nil {
			t.Errorf("parseMatchLine returned nil for %q", line)
			continue
		}
		if got.Optional != "" || got.Required != "" || got.ForceNew || got.Deprecated != "" {
			t.Errorf("prose %q must not set modifiers, got Optional=%q Required=%q ForceNew=%v Deprecated=%q",
				line, got.Optional, got.Required, got.ForceNew, got.Deprecated)
		}
	}
}

// TestUnitParseMatchLineAttributeSimilarWordsNotModifier checks near-miss tokens
// do not satisfy the \b<modifier>\b boundary inside a parenthesized group, so
// "(Optionals)", lowercase "(optional)" and split "(force new)" stay inert.
func TestUnitParseMatchLineAttributeSimilarWordsNotModifier(t *testing.T) {
	cases := []string{
		"    * `field` - (Optionals) not a valid modifier.",
		"    * `field` - (Requireds) not a valid modifier.",
		"    * `field` - (optional) lowercase is not the modifier token.",
		"    * `field` - (required) lowercase is not the modifier token.",
		"    * `field` - (force new) split words are not the modifier token.",
	}
	for _, line := range cases {
		words := attrWords(line)
		if words == nil {
			t.Errorf("line did not match regex: %q", line)
			continue
		}
		got := parseMatchLine(words, "Attribute", "")
		if got == nil {
			t.Errorf("parseMatchLine returned nil for %q", line)
			continue
		}
		if got.Optional != "" || got.Required != "" || got.ForceNew {
			t.Errorf("near-miss %q must not set modifiers, got Optional=%q Required=%q ForceNew=%v",
				line, got.Optional, got.Required, got.ForceNew)
		}
	}
}

// TestUnitParseMatchLineAttributeParenModifierWithProse verifies a real modifier in
// parentheses is still recognized when the following prose reuses the same
// modifier word — no double-counting, no false negative from the stricter
// paren-scoped match.
func TestUnitParseMatchLineAttributeParenModifierWithProse(t *testing.T) {
	cases := []struct {
		line     string
		opt, req string
		force    bool
	}{
		{"    * `field` - (Optional) This field is Optional in the config.", "true", "", false},
		{"    * `field` - (Required) The field is Required for the call.", "", "true", false},
		{"    * `field` - (Required, ForceNew) ForceNew applies on change.", "", "true", true},
		{"    * `field` - (ForceNew, Optional) Optional and ForceNew together.", "true", "", true},
	}
	for _, c := range cases {
		words := attrWords(c.line)
		if words == nil {
			t.Errorf("line did not match regex: %q", c.line)
			continue
		}
		got := parseMatchLine(words, "Attribute", "")
		if got == nil {
			t.Errorf("parseMatchLine returned nil for %q", c.line)
			continue
		}
		if got.Optional != c.opt || got.Required != c.req || got.ForceNew != c.force {
			t.Errorf("paren+prose %q: got Optional=%q Required=%q ForceNew=%v, want %q %q %v",
				c.line, got.Optional, got.Required, got.ForceNew, c.opt, c.req, c.force)
		}
	}
}

// TestUnitConsistencyCheckUnmarkedOptionalStillFails is the minimal end-to-end
// reproduction requested in review: when the schema marks a nested attribute
// Optional but the doc line has NO parenthesized modifier, parseMatchLine
// returns Optional="" and consistencyCheck reports "should be marked as
// Optional". The genuine missing-modifier failure must still fire after the
// fix — the stricter match does not relax real gaps.
func TestUnitConsistencyCheckUnmarkedOptionalStillFails(t *testing.T) {
	schemaDefined := map[string]*schema.Schema{
		"balancers.tags": {Optional: true},
	}
	unmarked := *parseMatchLine(
		attrWords("    * `tags` - The tag of the resource."), "Attribute", "balancers")
	if unmarked.Optional != "" {
		t.Fatalf("unmarked docs must parse Optional=\"\", got Optional=%q", unmarked.Optional)
	}
	docs := map[string]ResourceAttribute{"balancers.tags": unmarked}
	if consistencyCheck("alicloud_nlb_load_balancers", docs, schemaDefined) {
		t.Error("consistencyCheck should FAIL when schema Optional=true but docs has no (Optional) modifier")
	}
}

// TestUnitConsistencyCheckMarkedOptionalPasses verifies the fix end-to-end: with
// the doc line carrying "(Optional)", the same schema no longer false-flags
// the nested balancers.tags attribute (the original D1 false positive).
func TestUnitConsistencyCheckMarkedOptionalPasses(t *testing.T) {
	schemaDefined := map[string]*schema.Schema{
		"balancers.tags": {Optional: true},
	}
	marked := *parseMatchLine(
		attrWords("    * `tags` - (Optional) The tag of the resource."), "Attribute", "balancers")
	if marked.Optional != "true" {
		t.Fatalf("marked docs must parse Optional=\"true\", got Optional=%q", marked.Optional)
	}
	docs := map[string]ResourceAttribute{"balancers.tags": marked}
	if !consistencyCheck("alicloud_nlb_load_balancers", docs, schemaDefined) {
		t.Error("consistencyCheck should PASS when schema Optional=true and docs marks (Optional)")
	}
}

// TestUnitParseMatchLineAttributeParenProseAndUnclosedNotModifier guards the
// boundary flagged in the second review: the prior \([^)]*\b<modifier>\b would
// match an unclosed "(Optional" or a prose parenthesis like "(Required for
// subsequent API calls)" or "(Optional metadata returned by the service)".
// The closed-modifier-group regex rejects these: an unclosed group has no ")"
// and a prose parenthesis has non-modifier words after the token, so neither
// sets a modifier. None of these lines carry a real modifier group.
func TestUnitParseMatchLineAttributeParenProseAndUnclosedNotModifier(t *testing.T) {
	cases := []string{
		"    * `status` - (Optional metadata returned by the service).",
		"    * `id` - (Required for subsequent API calls).",
		"    * `field` - (ForceNew is not supported for this field).",
		"    * `field` - (Optional",
		"    * `field` - (Required",
		"    * `field` - note (Optional unclosed paren mid-sentence.",
	}
	for _, line := range cases {
		words := attrWords(line)
		if words == nil {
			t.Errorf("line did not match regex: %q", line)
			continue
		}
		got := parseMatchLine(words, "Attribute", "")
		if got == nil {
			t.Errorf("parseMatchLine returned nil for %q", line)
			continue
		}
		if got.Optional != "" || got.Required != "" || got.ForceNew || got.Deprecated != "" {
			t.Errorf("unclosed/prose paren %q must not set modifiers, got Optional=%q Required=%q ForceNew=%v Deprecated=%q",
				line, got.Optional, got.Required, got.ForceNew, got.Deprecated)
		}
	}
}

// TestUnitParseMatchLineAttributeClosedGroupStillRecognized confirms the
// stricter closed-group regex still recognizes every real modifier group
// shape: single tokens and comma-separated multi-token groups, in any order,
// with or without trailing prose, including two groups on one line. This
// guards against over-tightening that would reject the genuine "(Optional)"
// / "(Required, ForceNew)" contracts.
func TestUnitParseMatchLineAttributeClosedGroupStillRecognized(t *testing.T) {
	cases := []struct {
		line     string
		opt, req string
		force    bool
	}{
		{"    * `field` - (Optional) The tag.", "true", "", false},
		{"    * `field` - (Required) The id.", "", "true", false},
		{"    * `field` - (ForceNew) Recreates on change.", "", "", true},
		{"    * `field` - (Required, ForceNew) The id.", "", "true", true},
		{"    * `field` - (Optional, ForceNew) The tag.", "true", "", true},
		{"    * `field` - (ForceNew, Optional) The tag.", "true", "", true},
		{"    * `field` - (Optional, Required, ForceNew) compact.", "true", "true", true},
		{"    * `field` - (Optional) prose (Required) two groups.", "true", "true", false},
	}
	for _, c := range cases {
		words := attrWords(c.line)
		if words == nil {
			t.Errorf("line did not match regex: %q", c.line)
			continue
		}
		got := parseMatchLine(words, "Attribute", "")
		if got == nil {
			t.Errorf("parseMatchLine returned nil for %q", c.line)
			continue
		}
		if got.Optional != c.opt || got.Required != c.req || got.ForceNew != c.force {
			t.Errorf("closed group %q: got Optional=%q Required=%q ForceNew=%v, want %q %q %v",
				c.line, got.Optional, got.Required, got.ForceNew, c.opt, c.req, c.force)
		}
	}
}
