//nolint:all
package main

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	log "github.com/sirupsen/logrus"
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

// testArgsFieldRegex mirrors the package-internal argumentsFieldRegex used by
// parseResourceDocs for the Argument Reference phase. The original
// ssl_vpn_client_cert counterexample documented its Computed-only fields
// there, so the reverse-consistency tests must exercise that doc shape.
var testArgsFieldRegex = regexp.MustCompile("^\\* `([a-zA-Z_0-9]*)`[ ]*-? ?(\\(.*\\)) ?(.*)")

func argsWords(line string) []string {
	m := testArgsFieldRegex.FindStringSubmatch(line)
	if m == nil {
		return nil
	}
	return m
}

// captureConsistencyLogs redirects the package logger to a buffer while fn
// runs, so the tests can assert that consistencyCheck names the offending
// attribute in its error output.
func captureConsistencyLogs(t *testing.T, fn func()) string {
	t.Helper()
	var buf bytes.Buffer
	old := log.StandardLogger().Out
	log.SetOutput(&buf)
	defer log.SetOutput(old)
	fn()
	return buf.String()
}

// TestUnitConsistencyCheckComputedOnlyDocsOptionalRejected is the reverse
// negative for Optional: a Computed-only schema attribute documented with
// "(Optional)" must be rejected and the failure must name the attribute.
// Before the fix only the forward branch existed (schema Optional=true ⇒
// docs must say Optional), so this direction slipped through.
func TestUnitConsistencyCheckComputedOnlyDocsOptionalRejected(t *testing.T) {
	schemaDefined := map[string]*schema.Schema{
		"status": {Type: schema.TypeString, Computed: true},
	}
	docs := map[string]ResourceAttribute{
		"status": *parseMatchLine(
			argsWords("* `status` - (Optional) The status of the client certificate."), "Argument", ""),
	}
	if docs["status"].Optional != "true" {
		t.Fatalf("counterexample docs must parse Optional=\"true\", got %q", docs["status"].Optional)
	}
	out := captureConsistencyLogs(t, func() {
		if consistencyCheck("alicloud_ssl_vpn_client_cert", docs, schemaDefined) {
			t.Error("consistencyCheck should FAIL when a Computed-only schema attribute is documented as (Optional)")
		}
	})
	if !strings.Contains(out, "'status'") {
		t.Errorf("failure output must name the offending attribute 'status', got: %s", out)
	}
}

// TestUnitConsistencyCheckComputedOnlyDocsRequiredRejected is the reverse
// negative for Required: a Computed-only schema attribute documented with
// "(Required)" must be rejected and the failure must name the attribute.
func TestUnitConsistencyCheckComputedOnlyDocsRequiredRejected(t *testing.T) {
	schemaDefined := map[string]*schema.Schema{
		"ca_cert": {Type: schema.TypeString, Computed: true},
	}
	docs := map[string]ResourceAttribute{
		"ca_cert": *parseMatchLine(
			argsWords("* `ca_cert` - (Required) The client ca cert."), "Argument", ""),
	}
	if docs["ca_cert"].Required != "true" {
		t.Fatalf("counterexample docs must parse Required=\"true\", got %q", docs["ca_cert"].Required)
	}
	out := captureConsistencyLogs(t, func() {
		if consistencyCheck("alicloud_ssl_vpn_client_cert", docs, schemaDefined) {
			t.Error("consistencyCheck should FAIL when a Computed-only schema attribute is documented as (Required)")
		}
	})
	if !strings.Contains(out, "'ca_cert'") {
		t.Errorf("failure output must name the offending attribute 'ca_cert', got: %s", out)
	}
}

// TestUnitConsistencyCheckComputedOnlyCompliantDocsPass keeps the compliant
// shapes passing after the reverse branch is added: a Computed-only attribute
// documented as a read-only attribute (Attribute phase, no modifier), a real
// Optional, an Optional+Computed and a Required attribute keep passing. This
// guards against the reverse check over-reaching into forward-consistent
// documentation.
func TestUnitConsistencyCheckComputedOnlyCompliantDocsPass(t *testing.T) {
	cases := []struct {
		name        string
		schemaField string
		schemaDef   *schema.Schema
		docLine     string
		phase       string
	}{
		{
			name:        "computed-only documented read-only",
			schemaField: "status",
			schemaDef:   &schema.Schema{Type: schema.TypeString, Computed: true},
			docLine:     "  * `status` - The status of the client certificate.",
			phase:       "Attribute",
		},
		{
			name:        "real optional",
			schemaField: "name",
			schemaDef:   &schema.Schema{Type: schema.TypeString, Optional: true},
			docLine:     "* `name` - (Optional) The name of the client certificate.",
			phase:       "Argument",
		},
		{
			name:        "optional plus computed",
			schemaField: "name",
			schemaDef:   &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true},
			docLine:     "* `name` - (Optional) The name of the client certificate.",
			phase:       "Argument",
		},
		{
			name:        "required",
			schemaField: "ssl_vpn_server_id",
			schemaDef:   &schema.Schema{Type: schema.TypeString, Required: true},
			docLine:     "* `ssl_vpn_server_id` - (Required) The ID of the SSL-VPN server.",
			phase:       "Argument",
		},
	}
	for _, c := range cases {
		var words []string
		if c.phase == "Argument" {
			words = argsWords(c.docLine)
		} else {
			words = attrWords(c.docLine)
		}
		if words == nil {
			t.Errorf("%s: doc line did not match the %s-phase regex: %q", c.name, c.phase, c.docLine)
			continue
		}
		attribute := parseMatchLine(words, c.phase, "")
		if attribute == nil {
			t.Errorf("%s: parseMatchLine returned nil for %q", c.name, c.docLine)
			continue
		}
		schemaDefined := map[string]*schema.Schema{c.schemaField: c.schemaDef}
		docs := map[string]ResourceAttribute{c.schemaField: *attribute}
		out := captureConsistencyLogs(t, func() {
			if !consistencyCheck("alicloud_ssl_vpn_client_cert", docs, schemaDefined) {
				t.Error("compliant docs should PASS")
			}
		})
		if strings.Contains(out, "should not be marked") {
			t.Errorf("%s: reverse check must stay silent for compliant docs, got: %s", c.name, out)
		}
	}
}

// TestUnitConsistencyCheckComputedOnlyProseNotReverseFlagged guards the
// read-only documentation shape against prose false positives: an Attribute
// phase line whose description tail merely mentions "Optional" or "Required"
// in free-form prose sets no modifier (the closed modifier-group regex), so
// the reverse branch must not fire.
func TestUnitConsistencyCheckComputedOnlyProseNotReverseFlagged(t *testing.T) {
	schemaDefined := map[string]*schema.Schema{
		"status": {Type: schema.TypeString, Computed: true},
	}
	cases := []string{
		"  * `status` - Optional metadata returned by the service.",
		"  * `status` - Required for subsequent API calls.",
		"  * `status` - (Optional metadata returned by the service).",
	}
	for _, line := range cases {
		words := attrWords(line)
		if words == nil {
			t.Errorf("line did not match regex: %q", line)
			continue
		}
		attribute := parseMatchLine(words, "Attribute", "")
		if attribute == nil {
			t.Errorf("parseMatchLine returned nil for %q", line)
			continue
		}
		if attribute.Optional != "" || attribute.Required != "" {
			t.Errorf("prose %q must not parse modifiers, got Optional=%q Required=%q", line, attribute.Optional, attribute.Required)
			continue
		}
		docs := map[string]ResourceAttribute{"status": *attribute}
		if !consistencyCheck("alicloud_ssl_vpn_client_cert", docs, schemaDefined) {
			t.Errorf("prose %q must not trigger the reverse check", line)
		}
	}
}

// TestUnitConsistencyCheckSslVpnClientCertOriginalCounterexample reproduces
// the original counterexample from ticket 40430 / CR 31118684 at its
// pre-fix shape: five Computed-only schema attributes whose Argument
// Reference lines carried "(Optional)" (website/docs/r/ssl_vpn_client_cert.
// html.markdown on master). The reverse branch must reject the combination
// and name each of the five fields.
func TestUnitConsistencyCheckSslVpnClientCertOriginalCounterexample(t *testing.T) {
	schemaDefined := map[string]*schema.Schema{
		"ssl_vpn_server_id": {Type: schema.TypeString, Required: true, ForceNew: true},
		"name":              {Type: schema.TypeString, Optional: true},
		"status":            {Type: schema.TypeString, Computed: true},
		"ca_cert":           {Type: schema.TypeString, Computed: true},
		"client_cert":       {Type: schema.TypeString, Computed: true},
		"client_key":        {Type: schema.TypeString, Computed: true},
		"client_config":     {Type: schema.TypeString, Computed: true},
	}
	docLines := []string{
		"* `name` - (Optional) The name of the client certificate.",
		"* `ssl_vpn_server_id` - (Required, ForceNew) The ID of the SSL-VPN server.",
		"* `status` - (Optional) The status of the client certificate.",
		"* `ca_cert` - (Optional) The client ca cert.",
		"* `client_cert` - (Optional) The client cert.",
		"* `client_key` - (Optional) The client key.",
		"* `client_config` - (Optional) The vpn client config.",
	}
	docs := make(map[string]ResourceAttribute)
	for _, line := range docLines {
		words := argsWords(line)
		if words == nil {
			t.Fatalf("counterexample doc line did not match the Argument-phase regex: %q", line)
		}
		attribute := parseMatchLine(words, "Argument", "")
		if attribute == nil {
			t.Fatalf("parseMatchLine returned nil for %q", line)
		}
		docs[attribute.Name] = *attribute
	}
	var inconsistent bool
	out := captureConsistencyLogs(t, func() {
		inconsistent = !consistencyCheck("alicloud_ssl_vpn_client_cert", docs, schemaDefined)
	})
	if !inconsistent {
		t.Fatal("consistencyCheck should FAIL on the original ssl_vpn_client_cert counterexample")
	}
	for _, field := range []string{"status", "ca_cert", "client_cert", "client_key", "client_config"} {
		if !strings.Contains(out, "'"+field+"'") {
			t.Errorf("failure output must name the misdocumented attribute %q, got: %s", field, out)
		}
	}
}

// TestUnitConsistencyCheckSslVpnClientCertCorrectedDocsPass is the corrected
// counterpart of the counterexample: the five read-only attributes move to
// the Attributes Reference with no modifier while the configurable arguments
// keep theirs, and consistencyCheck passes — the shape the original CR
// 31118684 corrected the docs to.
func TestUnitConsistencyCheckSslVpnClientCertCorrectedDocsPass(t *testing.T) {
	schemaDefined := map[string]*schema.Schema{
		"ssl_vpn_server_id": {Type: schema.TypeString, Required: true, ForceNew: true},
		"name":              {Type: schema.TypeString, Optional: true},
		"status":            {Type: schema.TypeString, Computed: true},
		"ca_cert":           {Type: schema.TypeString, Computed: true},
		"client_cert":       {Type: schema.TypeString, Computed: true},
		"client_key":        {Type: schema.TypeString, Computed: true},
		"client_config":     {Type: schema.TypeString, Computed: true},
	}
	argumentLines := []string{
		"* `name` - (Optional) The name of the client certificate.",
		"* `ssl_vpn_server_id` - (Required, ForceNew) The ID of the SSL-VPN server.",
	}
	attributeLines := []string{
		"  * `status` - The status of the client certificate.",
		"  * `ca_cert` - The client ca cert.",
		"  * `client_cert` - The client cert.",
		"  * `client_key` - The client key.",
		"  * `client_config` - The vpn client config.",
	}
	docs := make(map[string]ResourceAttribute)
	for _, line := range argumentLines {
		attribute := parseMatchLine(argsWords(line), "Argument", "")
		if attribute == nil {
			t.Fatalf("parseMatchLine returned nil for %q", line)
		}
		docs[attribute.Name] = *attribute
	}
	for _, line := range attributeLines {
		attribute := parseMatchLine(attrWords(line), "Attribute", "")
		if attribute == nil {
			t.Fatalf("parseMatchLine returned nil for %q", line)
		}
		docs[attribute.Name] = *attribute
	}
	if !consistencyCheck("alicloud_ssl_vpn_client_cert", docs, schemaDefined) {
		t.Error("corrected ssl_vpn_client_cert docs should PASS")
	}
}
