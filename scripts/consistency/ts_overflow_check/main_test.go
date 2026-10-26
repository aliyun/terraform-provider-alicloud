//nolint:all
package main

import (
	"bytes"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// parseInline parses src (declared as package alicloud so resource/data_source prefixes are
// irrelevant) and runs analyzeAST, returning the fileResult. file is the synthetic path used
// for diagnostics and docFor.
func parseInline(t *testing.T, file, src string) fileResult {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, src, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse %s: %v", file, err)
	}
	return analyzeAST(file, fset, f, []byte(src))
}

// r1PreFixSrc is the exact pre-fix base pattern from resource_alicloud_ots_tunnel.go: a
// TypeInt create_time fed by d.Set(..., t.UnixNano()). This MUST trigger R1.
const r1PreFixSrc = `package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOtsTunnel() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"create_time": {Type: schema.TypeInt, Computed: true},
			"tunnel_rpo":  {Type: schema.TypeInt, Computed: true},
			"label":       {Type: schema.TypeString, Optional: true},
		},
	}
}

func readTunnel(d *schema.ResourceData, t time.Time) {
	d.Set("create_time", t.UnixNano())
	d.Set("tunnel_rpo", t.UnixNano())
	d.Set("label", fmt.Sprint("v"))
}
`

// r1PostFixSrc is the post-fix pattern: TypeString + fmt.Sprint(UnixNano()). This MUST NOT
// trigger R1 (field is TypeString; even though UnixNano() appears, it is not fed to a
// TypeInt field).
const r1PostFixSrc = `package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOtsTunnel() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"create_time": {Type: schema.TypeString, Computed: true},
			"tunnel_rpo":  {Type: schema.TypeString, Computed: true},
			"label":       {Type: schema.TypeString, Optional: true},
		},
	}
}

func readTunnel(d *schema.ResourceData, t time.Time) {
	d.Set("create_time", fmt.Sprint(t.UnixNano()))
	d.Set("tunnel_rpo", fmt.Sprint(t.UnixNano()))
	d.Set("label", fmt.Sprint("v"))
}
`

func TestUnitR1DetectsPreFixPanicPattern(t *testing.T) {
	// Positive case (修复前 base 正例): TypeInt create_time fed from UnixNano() must hit.
	fr := parseInline(t, "alicloud/resource_alicloud_ots_tunnel.go", r1PreFixSrc)
	if len(fr.r1Hits) != 2 {
		t.Fatalf("expected 2 R1 hits (create_time + tunnel_rpo), got %d: %+v", len(fr.r1Hits), fr.r1Hits)
	}
	for _, h := range fr.r1Hits {
		if h.field != "create_time" && h.field != "tunnel_rpo" {
			t.Errorf("unexpected R1 field: %q", h.field)
		}
		if !strings.Contains(h.expr, "UnixNano()") {
			t.Errorf("R1 expr for %q should contain the raw UnixNano() call, got %q", h.field, h.expr)
		}
	}
}

func TestUnitR1NoHitOnPostFixTypeString(t *testing.T) {
	// Negative case (修复后): TypeString + fmt.Sprint(UnixNano()) must not hit.
	fr := parseInline(t, "alicloud/resource_alicloud_ots_tunnel.go", r1PostFixSrc)
	if len(fr.r1Hits) != 0 {
		t.Fatalf("expected 0 R1 hits on post-fix TypeString, got %d: %+v", len(fr.r1Hits), fr.r1Hits)
	}
}

func TestUnitR1NoHitWhenUnixTimeOnTypeStringField(t *testing.T) {
	// Negative case: UnixNano() fed to a TypeString field must not hit (field not in
	// typeIntFields). This guards against an over-broad R1 that keys only on the call.
	src := `package alicloud
import ("time"; "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema")
func r() *schema.Resource { return &schema.Resource{Schema: map[string]*schema.Schema{"create_time": {Type: schema.TypeString, Computed: true}}}}
func rd(d *schema.ResourceData, t time.Time) { d.Set("create_time", t.UnixNano()) }
`
	fr := parseInline(t, "alicloud/resource_alicloud_x.go", src)
	if len(fr.r1Hits) != 0 {
		t.Fatalf("expected 0 R1 hits when field is TypeString, got %d", len(fr.r1Hits))
	}
}

func TestUnitR1NoHitWhenTypeIntFedNonUnixExpr(t *testing.T) {
	// Negative case: TypeInt field fed a non-Unix expression must not hit.
	src := `package alicloud
import ("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema")
func r() *schema.Resource { return &schema.Resource{Schema: map[string]*schema.Schema{"create_time": {Type: schema.TypeInt, Computed: true}}}}
func rd(d *schema.ResourceData) { d.Set("create_time", 42) }
`
	fr := parseInline(t, "alicloud/resource_alicloud_x.go", src)
	if len(fr.r1Hits) != 0 {
		t.Fatalf("expected 0 R1 hits for non-Unix expr, got %d", len(fr.r1Hits))
	}
}

func TestUnitUnixTimeCallInDetectsAllSelectors(t *testing.T) {
	cases := map[string]string{
		"milli": `package alicloud
import ("time"; "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema")
func r() *schema.Resource { return &schema.Resource{Schema: map[string]*schema.Schema{"create_time": {Type: schema.TypeInt, Computed: true}}}}
func rd(d *schema.ResourceData, t time.Time) { d.Set("create_time", t.UnixMilli()) }`,
		"micro": `package alicloud
import ("time"; "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema")
func r() *schema.Resource { return &schema.Resource{Schema: map[string]*schema.Schema{"create_time": {Type: schema.TypeInt, Computed: true}}}}
func rd(d *schema.ResourceData, t time.Time) { d.Set("create_time", t.UnixMicro()) }`,
		"chain": `package alicloud
import ("time"; "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema")
func r() *schema.Resource { return &schema.Resource{Schema: map[string]*schema.Schema{"create_time": {Type: schema.TypeInt, Computed: true}}}}
func rd(d *schema.ResourceData) { d.Set("create_time", time.Now().UnixNano()) }`,
	}
	for name, src := range cases {
		t.Run(name, func(t *testing.T) {
			fr := parseInline(t, "alicloud/resource_alicloud_x.go", src)
			if len(fr.r1Hits) != 1 {
				t.Fatalf("%s: expected 1 R1 hit, got %d", name, len(fr.r1Hits))
			}
		})
	}
}

// wafv3NestedSrc mirrors resource_alicloud_wafv3_defense_rule.go: time_periods.start/end
// are TypeString (the fix) and day_periods.start/end are TypeInt (the bounded exemption).
const wafv3NestedSrc = `package alicloud
import ("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema")
func r() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"config": {Type: schema.TypeList, Optional: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{
			"time_config": {Type: schema.TypeList, Optional: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				"time_periods": {Type: schema.TypeList, Optional: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"start": {Type: schema.TypeString, Optional: true},
					"end":   {Type: schema.TypeString, Optional: true},
				}}},
				"week_time_periods": {Type: schema.TypeList, Optional: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"day_periods": {Type: schema.TypeList, Optional: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{
						"start": {Type: schema.TypeInt, Optional: true},
						"end":   {Type: schema.TypeInt, Optional: true},
					}}},
				}}},
			}}},
		}}},
	}}
}
`

func TestUnitR2ExemptionKeyUsesNestedPath(t *testing.T) {
	// day_periods.start is TypeInt and its doc declares a sub-second timestamp; it must be
	// exempted by the "alicloud_wafv3_defense_rule::day_periods.start" key, NOT reported.
	fr := parseInline(t, "alicloud/resource_alicloud_wafv3_defense_rule.go", wafv3NestedSrc)
	doc := "### `config-time_config-week_time_periods-day_periods`\n" +
		"The day_periods supports:\n" +
		"* `end` - (Optional, Int) millisecond-level timestamp relative to 00:00. [0-86400000)\n" +
		"* `start` - (Optional, Int) millisecond-level timestamp relative to 00:00. [0-86400000)\n"
	got := checkR2Doc(fr, "alicloud_wafv3_defense_rule", doc, false)
	if len(got) != 0 {
		t.Fatalf("expected day_periods.start/end to be exempted, got %d violations: %v", len(got), got)
	}
	// Sanity: the path-aware collection produced the dotted nested path.
	var hasDayStart, hasTimeStart bool
	for _, fi := range fr.pathFields {
		if fi.path == "day_periods.start" && fi.typ == "Int" {
			hasDayStart = true
		}
		if fi.path == "time_periods.start" && fi.typ == "String" {
			hasTimeStart = true
		}
	}
	if !hasDayStart || !hasTimeStart {
		t.Fatalf("path-aware collection missing nested paths: hasDayStart=%v hasTimeStart=%v", hasDayStart, hasTimeStart)
	}
}

func TestUnitR2ReportsUnexemptedTimestampField(t *testing.T) {
	// A TypeInt field whose doc declares a sub-second timestamp and has NO exemption must
	// be reported, with the dotted path in the exemption key.
	src := `package alicloud
import ("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema")
func r() *schema.Resource { return &schema.Resource{Schema: map[string]*schema.Schema{
	"block": {Type: schema.TypeList, Optional: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{
		"create_time": {Type: schema.TypeInt, Computed: true},
	}}}}}}
`
	fr := parseInline(t, "alicloud/resource_alicloud_x.go", src)
	doc := "* `create_time` - create time, a UNIX timestamp in milliseconds.\n"
	got := checkR2Doc(fr, "alicloud_x", doc, false)
	if len(got) != 1 {
		t.Fatalf("expected 1 R2 violation for block.create_time, got %d: %v", len(got), got)
	}
	if !strings.Contains(got[0], "block.create_time") {
		t.Errorf("R2 violation should name the dotted path block.create_time, got %q", got[0])
	}
}

func TestUnitR2NoFalsePositiveFromNeighborAttribute(t *testing.T) {
	// image_size is TypeInt; its description says "in bytes" (no timestamp). The next
	// attribute image_create says "unix time in nanoseconds". The fixed-width window bug
	// flagged image_size; the list-item-line approach must NOT.
	src := `package alicloud
import ("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema")
func r() *schema.Resource { return &schema.Resource{Schema: map[string]*schema.Schema{
	"image_size":   {Type: schema.TypeInt, Computed: true},
	"image_create": {Type: schema.TypeString, Computed: true},
}}}
`
	fr := parseInline(t, "alicloud/data_source_alicloud_cr_repos.go", src)
	doc := "* `image_size` - Status of this image, in bytes.\n" +
		"* `image_create` - Create time of this image, unix time in nanoseconds.\n"
	got := checkR2Doc(fr, "alicloud_cr_repos", doc, false)
	if len(got) != 0 {
		t.Fatalf("image_size must not be flagged by its neighbour's description, got %d: %v", len(got), got)
	}
}

func TestUnitDocForPathIsHtmlMarkdown(t *testing.T) {
	cases := []struct {
		file, wantDir, wantRes string
	}{
		{"alicloud/resource_alicloud_ots_tunnel.go", "r", "alicloud_ots_tunnel"},
		{"alicloud/data_source_alicloud_cr_repos.go", "d", "alicloud_cr_repos"},
	}
	for _, c := range cases {
		path, res := docFor(c.file)
		if res != c.wantRes {
			t.Errorf("docFor(%q) resource = %q, want %q", c.file, res, c.wantRes)
		}
		if !strings.HasSuffix(path, ".html.markdown") {
			t.Errorf("docFor(%q) path = %q, must end with .html.markdown", c.file, path)
		}
		if !strings.Contains(path, filepath.Join("website", "docs", c.wantDir)) {
			t.Errorf("docFor(%q) path = %q, must be under website/docs/%s", c.file, path, c.wantDir)
		}
	}
}

func TestUnitR3FlagsNewTimestampNamedTypeInt(t *testing.T) {
	// A TypeInt create_time not in the baseline is a NEW R3 hit; a TypeString create_time is not.
	src := `package alicloud
import ("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema")
func r() *schema.Resource { return &schema.Resource{Schema: map[string]*schema.Schema{
	"create_time": {Type: schema.TypeInt, Computed: true},
	"timeout":     {Type: schema.TypeInt, Optional: true},
}}}
`
	fr := parseInline(t, "alicloud/resource_alicloud_x.go", src)
	if len(fr.r3Hits) != 1 || fr.r3Hits[0] != "create_time" {
		t.Fatalf("expected R3 to flag create_time only (timeout is blacklisted), got %v", fr.r3Hits)
	}
}

// withTempAlicloud creates a temp dir containing an alicloud/ subdir with the given .go
// files, chdirs into it, and restores the original working directory on cleanup.
func withTempAlicloud(t *testing.T, files map[string]string) string {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "alicloud"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, "alicloud", name), []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })
	return dir
}

func TestUnitResolveFilesAcceptsNameOnlyAndFullDiff(t *testing.T) {
	// F3: the CI orchestrator feeds the FULL `git diff` as -fileNames (parity with
	// consistency_check.go), while a local run may pass `git diff --name-only`. Both must
	// resolve to the same alicloud/resource_* and alicloud/data_source_* file set; a full
	// diff must never be mistaken for an empty path list (silent ship = masked regression).
	nameOnly := "alicloud/resource_alicloud_ots_tunnel.go\nalicloud/data_source_alicloud_cr_repos.go\nalicloud/service/ecs/resource_skip.go\nscripts/consistency/ts_overflow_check/main.go\n"
	fullDiff := "diff --git a/alicloud/resource_alicloud_ots_tunnel.go b/alicloud/resource_alicloud_ots_tunnel.go\n" +
		"index aaa..bbb 100644\n" +
		"--- a/alicloud/resource_alicloud_ots_tunnel.go\n" +
		"+++ b/alicloud/resource_alicloud_ots_tunnel.go\n" +
		"@@ -1 +1 @@\n" +
		"-old\n" +
		"+new\n" +
		"diff --git a/alicloud/data_source_alicloud_cr_repos.go b/alicloud/data_source_alicloud_cr_repos.go\n" +
		"--- a/alicloud/data_source_alicloud_cr_repos.go\n" +
		"+++ b/alicloud/data_source_alicloud_cr_repos.go\n"
	want := []string{"alicloud/resource_alicloud_ots_tunnel.go", "alicloud/data_source_alicloud_cr_repos.go"}
	for name, body := range map[string]string{"name-only": nameOnly, "full-diff": fullDiff} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			f := filepath.Join(dir, "in.txt")
			if err := os.WriteFile(f, []byte(body), 0o644); err != nil {
				t.Fatalf("write: %v", err)
			}
			got := resolveFiles(config{fileNames: f})
			if len(got) != len(want) {
				t.Fatalf("%s: expected %d files, got %d (%v)", name, len(want), len(got), got)
			}
			for i, p := range want {
				if got[i] != p {
					t.Errorf("%s: file[%d] = %q, want %q", name, i, got[i], p)
				}
			}
		})
	}
}

func TestUnitR6AllDoesNotWriteBaseline(t *testing.T) {
	// F6: a plain -all must NOT rewrite the baseline file, so it can never silently mask
	// new hits.
	dir := withTempAlicloud(t, map[string]string{
		"resource_alicloud_x.go": r1PreFixSrc,
	})
	blPath := filepath.Join(dir, "baseline.txt")
	var out bytes.Buffer
	code := check(config{all: true, baseline: blPath}, &out, &out)
	if code == 0 {
		t.Fatalf("expected -all to fail (new TypeInt timestamp field not in empty baseline), got PASS:\n%s", out.String())
	}
	if _, err := os.Stat(blPath); !os.IsNotExist(err) {
		t.Fatalf("-all must not create the baseline file; stat err = %v", err)
	}
	if !strings.Contains(out.String(), "R3 NEW") {
		t.Errorf("expected -all to report new R3 hits against the baseline, got:\n%s", out.String())
	}
}

// r3OnlySrc has a TypeInt create_time (an R3 hit) but NO d.Set feeding it from a Unix*
// call and no doc — so R1=0 and R2=0, only the R3 baseline entry is recorded. This is the
// clean tree -update-baseline is meant for.
const r3OnlySrc = `package alicloud
import ("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema")
func r() *schema.Resource { return &schema.Resource{Schema: map[string]*schema.Schema{
	"create_time": {Type: schema.TypeInt, Computed: true},
	"label":       {Type: schema.TypeString, Optional: true},
}}}
`

func TestUnitR6UpdateBaselineWritesOnlyOnFlag(t *testing.T) {
	// F6: -update-baseline writes the baseline with the TypeInt timestamp field; and only
	// when the flag is set.
	dir := withTempAlicloud(t, map[string]string{
		"resource_alicloud_x.go": r3OnlySrc,
	})
	blPath := filepath.Join(dir, "baseline.txt")
	var out bytes.Buffer
	code := check(config{updateBL: true, baseline: blPath}, &out, &out)
	if code != 0 {
		t.Fatalf("expected -update-baseline on a clean tree to PASS, got %d:\n%s", code, out.String())
	}
	info, err := os.Stat(blPath)
	if err != nil {
		t.Fatalf("-update-baseline must write the baseline file: %v", err)
	}
	if info.Size() == 0 {
		t.Fatalf("baseline file must be non-empty")
	}
	b, _ := os.ReadFile(blPath)
	if !strings.Contains(string(b), "resource_alicloud_x.go::create_time") {
		t.Errorf("baseline should record the TypeInt timestamp field, got:\n%s", b)
	}
}
