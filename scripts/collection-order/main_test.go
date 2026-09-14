package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func TestCollectionSemantics(t *testing.T) {
	for _, tc := range []struct {
		name, order         string
		typeOf              schema.ValueType
		forceNew, wantDrift bool
	}{
		{"set", "", schema.TypeSet, false, false},
		{"ordered list", "ordered", schema.TypeList, false, false},
		{"unordered list mutation", "unordered", schema.TypeList, false, true},
		{"replacement drift", "unordered", schema.TypeList, true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := &schema.Resource{Schema: map[string]*schema.Schema{"values": {
				Type: tc.typeOf, Optional: true, ForceNew: tc.forceNew, Elem: &schema.Schema{Type: schema.TypeString},
			}}}
			c := testCase{Path: "values", Order: tc.order, Members: []interface{}{"alice", "bob", "carol"}}
			drift, err := checkCase(r, c)
			if err != nil || (drift != "") != tc.wantDrift {
				t.Fatalf("drift=%q err=%v", drift, err)
			}
		})
	}
	t.Run("list with membership suppression", func(t *testing.T) {
		r := &schema.Resource{Schema: map[string]*schema.Schema{"values": {
			Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString},
			DiffSuppressFunc: func(_ string, _, _ string, d *schema.ResourceData) bool {
				old, new := d.GetChange("values")
				left, right := append([]interface{}{}, old.([]interface{})...), append([]interface{}{}, new.([]interface{})...)
				sort.Slice(left, func(i, j int) bool { return left[i].(string) < left[j].(string) })
				sort.Slice(right, func(i, j int) bool { return right[i].(string) < right[j].(string) })
				return reflect.DeepEqual(left, right)
			},
		}}}
		c := testCase{Path: "values", Order: "unordered", Members: []interface{}{"alice", "bob", "carol"}}
		if drift, err := checkCase(r, c); drift != "" || err != nil {
			t.Fatalf("drift=%q err=%v", drift, err)
		}
		// Mutation: suppressing everything must fail the membership controls.
		r.Schema["values"].DiffSuppressFunc = func(string, string, string, *schema.ResourceData) bool { return true }
		if _, err := checkCase(r, c); err == nil || !strings.Contains(err.Error(), "add: want diff=true") {
			t.Fatalf("overbroad suppressor escaped: %v", err)
		}
		// Resource-level CustomizeDiff must not be omitted from the check.
		r.Schema["values"].DiffSuppressFunc = nil
		r.Schema["values"].Computed = true
		r.CustomizeDiff = func(d *schema.ResourceDiff, _ interface{}) error { return d.Clear("values") }
		c.NilMetaReason = "The synthetic callback only clears the values field and never reads meta."
		if _, err := checkCase(r, c); err == nil {
			t.Fatal("CustomizeDiff swallowing membership changes escaped")
		}
	})
	t.Run("nested set and mandatory declaration", func(t *testing.T) {
		r := &schema.Resource{Schema: map[string]*schema.Schema{"block": {
			Type: schema.TypeList, Optional: true, MaxItems: 1,
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{"values": {
				Type: schema.TypeSet, Optional: true, Elem: &schema.Schema{Type: schema.TypeString},
			}}},
		}}}
		c := testCase{Resource: "example", Path: "block.0.values", Members: []interface{}{"alice", "bob", "carol"}, Config: map[string]interface{}{"block": []interface{}{map[string]interface{}{}}}}
		if drift, err := checkCase(r, c); drift != "" || err != nil {
			t.Fatalf("drift=%q err=%v", drift, err)
		}
		resources := map[string]*schema.Resource{"example": r}
		affected := map[string]bool{"example": true}
		if missing := coverageErrors(resources, affected, nil); len(missing) != 1 || !strings.Contains(missing[0], "example.block.values") {
			t.Fatalf("new nested collection escaped registry: %v", missing)
		}
		if missing := coverageErrors(resources, affected, []testCase{c}); len(missing) != 0 {
			t.Fatal(missing)
		}
	})
	t.Run("outer set preserves inner list ordering", func(t *testing.T) {
		r := &schema.Resource{Schema: map[string]*schema.Schema{"blocks": {
			Type: schema.TypeSet, Optional: true,
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{"values": {
				Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString},
			}}},
		}}}
		c := testCase{Path: "blocks.0.values", Order: "unordered", Members: []interface{}{"alice", "bob", "carol"}, Config: map[string]interface{}{"blocks": []interface{}{map[string]interface{}{}}}}
		if drift, err := checkCase(r, c); drift == "" || err != nil {
			t.Fatalf("outer set hash change escaped: drift=%q err=%v", drift, err)
		}
	})
	t.Run("fixed size and invalid fixture", func(t *testing.T) {
		r := &schema.Resource{Schema: map[string]*schema.Schema{"values": {
			Type: schema.TypeSet, Optional: true, MinItems: 2, MaxItems: 2, Elem: &schema.Schema{Type: schema.TypeString},
		}}}
		c := testCase{Path: "values", Members: []interface{}{"alice", "bob", "carol", "dave"}}
		if drift, err := checkCase(r, c); drift != "" || err != nil {
			t.Fatalf("valid fixed-size replacement rejected: drift=%q err=%v", drift, err)
		}
		r.Schema["values"].MinItems = 3
		if _, err := checkCase(r, c); err == nil || !strings.Contains(err.Error(), "invalid SDK fixture") {
			t.Fatalf("invalid fixture escaped: %v", err)
		}
	})
	t.Run("two legal enum values", func(t *testing.T) {
		r := &schema.Resource{Schema: map[string]*schema.Schema{"values": {
			Type: schema.TypeList, Optional: true, MaxItems: 2,
			Elem: &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.StringInSlice([]string{"A", "B"}, false)},
			DiffSuppressFunc: func(_ string, _, _ string, d *schema.ResourceData) bool {
				old, new := d.GetChange("values")
				left, right := old.([]interface{}), new.([]interface{})
				return len(left) == 2 && len(right) == 2 && left[0] == right[1] && left[1] == right[0]
			},
		}}}
		c := testCase{Path: "values", Order: "unordered", Members: []interface{}{"A", "B"}}
		if drift, err := checkCase(r, c); drift != "" || err != nil {
			t.Fatalf("two enum values rejected: drift=%q err=%v", drift, err)
		}
		// A suppressor hiding singleton replacement must still be caught.
		suppress := r.Schema["values"].DiffSuppressFunc
		r.Schema["values"].DiffSuppressFunc = func(k, old, new string, d *schema.ResourceData) bool {
			left, right := d.GetChange("values")
			return suppress(k, old, new, d) || (len(left.([]interface{})) == 1 && len(right.([]interface{})) == 1)
		}
		if _, err := checkCase(r, c); err == nil || !strings.Contains(err.Error(), "replace: want diff=true") {
			t.Fatalf("singleton replacement escaped: %v", err)
		}
		r.Schema["values"].DiffSuppressFunc = func(string, string, string, *schema.ResourceData) bool { return true }
		if _, err := checkCase(r, c); err == nil || !strings.Contains(err.Error(), "add: want diff=true") {
			t.Fatalf("two-value membership addition escaped: %v", err)
		}
		r.Schema["values"].DiffSuppressFunc = suppress
		r.Schema["values"].MinItems = 2
		if _, err := checkCase(r, c); err == nil || !strings.Contains(err.Error(), "replacement control") {
			t.Fatalf("fixed pair without replacement sample escaped: %v", err)
		}
	})
	t.Run("fixed pair needs a replacement sample", func(t *testing.T) {
		r := &schema.Resource{Schema: map[string]*schema.Schema{"values": {
			Type: schema.TypeList, Optional: true, MinItems: 2, MaxItems: 2,
			Elem:             &schema.Schema{Type: schema.TypeString},
			DiffSuppressFunc: func(string, string, string, *schema.ResourceData) bool { return true },
		}}}
		c := testCase{Path: "values", Order: "unordered", Members: []interface{}{"A", "B"}}
		if _, err := checkCase(r, c); err == nil || !strings.Contains(err.Error(), "replacement control") {
			t.Fatalf("all membership controls were skipped: %v", err)
		}
		c.Members = append(c.Members, "C")
		if _, err := checkCase(r, c); err == nil || !strings.Contains(err.Error(), "replace: want diff=true") {
			t.Fatalf("fixed-size membership suppression escaped: %v", err)
		}
	})

	t.Run("rotation catches reversal-only suppression", func(t *testing.T) {
		r := &schema.Resource{Schema: map[string]*schema.Schema{"values": {
			Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString},
			DiffSuppressFunc: func(_ string, _, _ string, d *schema.ResourceData) bool {
				old, new := d.GetChange("values")
				left, right := old.([]interface{}), new.([]interface{})
				if len(left) != len(right) {
					return false
				}
				for i := range left {
					if left[i] != right[len(right)-1-i] {
						return false
					}
				}
				return true
			},
		}}}
		c := testCase{Path: "values", Order: "unordered", Members: []interface{}{"alice", "bob", "carol", "dave"}}
		if drift, err := checkCase(r, c); err != nil || !strings.Contains(drift, "rotated") {
			t.Fatalf("rotation escaped: drift=%q err=%v", drift, err)
		}
	})
	t.Run("reorder changes a sibling through CustomizeDiff", func(t *testing.T) {
		r := &schema.Resource{Schema: map[string]*schema.Schema{
			"values":  {Type: schema.TypeList, Optional: true, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"sibling": {Type: schema.TypeString, Computed: true},
		}, CustomizeDiff: func(d *schema.ResourceDiff, _ interface{}) error {
			old, new := d.GetChange("values")
			if reflect.DeepEqual(old, new) {
				return nil
			}
			left, right := old.([]interface{}), new.([]interface{})
			if len(left) == 2 && len(right) == 2 && left[0] == right[1] && left[1] == right[0] {
				if err := d.Clear("values"); err != nil {
					return err
				}
				return d.SetNew("sibling", "changed-by-order")
			}
			return nil
		}}
		c := testCase{Path: "values", Order: "unordered", Members: []interface{}{"alice", "bob", "carol"}, NilMetaReason: "Synthetic callback only reads configuration."}
		if drift, err := checkCase(r, c); err != nil || !strings.Contains(drift, "sibling") {
			t.Fatalf("sibling diff escaped: drift=%q err=%v", drift, err)
		}
	})
}

func TestProviderCollections(t *testing.T) {
	log.SetOutput(io.Discard)
	cases, err := loadCases("cases.json")
	if err != nil {
		t.Fatal(err)
	}
	p := alicloud.Provider().(*schema.Provider)
	for _, c := range cases {
		t.Run(c.key(), func(t *testing.T) {
			r := p.ResourcesMap[c.Resource]
			if strings.HasPrefix(c.Resource, "data.") {
				r = p.DataSourcesMap[strings.TrimPrefix(c.Resource, "data.")]
			}
			if r == nil {
				t.Fatal("resource no longer exists")
			}
			drift, err := checkCase(r, c)
			if err != nil || (drift != "") != knownDrift[c.key()] {
				t.Fatalf("drift=%q err=%v (remove a fixed legacy baseline)", drift, err)
			}
			if drift != "" {
				t.Logf("KNOWN DRIFT: %s", drift)
			}
		})
	}
}

func TestLoadTwoMembers(t *testing.T) {
	cases, err := loadCases("cases.json")
	if err != nil {
		t.Fatal(err)
	}
	cases[0].Members = []interface{}{"A", "B"}
	data, err := json.Marshal(cases)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "cases.json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := loadCases(path); err != nil {
		t.Fatal(err)
	}
}

func TestAffectedResources(t *testing.T) {
	source, err := os.ReadFile("../../alicloud/provider.go")
	if err != nil {
		t.Fatal(err)
	}
	registry, err := providerRegistry(source)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		file string
		want map[string]bool
	}{
		{"service_alicloud_vpc.go", map[string]bool{}},
		{"service_alicloud_ecs.go", map[string]bool{}},
		{"common.go", map[string]bool{}},
		{"service_alicloud_ecd.go", map[string]bool{"alicloud_ecd_desktop": true}},
		{"service_alicloud_ecd_v2.go", map[string]bool{"alicloud_ecd_desktop_group": true}},
		{"resource_alicloud_ecd_desktop_group.go", map[string]bool{"alicloud_ecd_desktop_group": true}},
		{"resource_alicloud_ram_group_membership.go", map[string]bool{"alicloud_ram_group_membership": true}},
	} {
		t.Run(tc.file, func(t *testing.T) {
			changed := map[string]bool{"alicloud/" + tc.file: true}
			direct, err := resourcesInFiles("../..", changed, registry, registry)
			if err != nil {
				t.Fatal(err)
			}
			got := affectedResources(changed, direct)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %d resources %v, want %v; err=%v", len(got), got, tc.want, err)
			}
		})
	}
}

func TestCoverageFileNames(t *testing.T) {
	source, err := os.ReadFile("../../alicloud/provider.go")
	if err != nil {
		t.Fatal(err)
	}
	registry, err := providerRegistry(source)
	if err != nil {
		t.Fatal(err)
	}
	files, err := os.ReadDir("../../alicloud")
	if err != nil {
		t.Fatal(err)
	}
	scanned := 0
	for _, file := range files {
		name := file.Name()
		if file.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") ||
			(!strings.HasPrefix(name, "resource_alicloud_") && !strings.HasPrefix(name, "data_source_alicloud_")) {
			continue
		}
		scanned++
		direct, err := resourcesInFiles("../..", map[string]bool{"alicloud/" + name: true}, registry, registry)
		want := 1
		if name == "data_source_alicloud_common.go" {
			want = 0
		}
		if err != nil || len(direct) != want {
			t.Errorf("%s: got %v, want %d resources; err=%v", name, direct, want, err)
		}
		for resource := range direct {
			if registry[resource] == "" {
				t.Errorf("%s: selected unregistered resource %s", name, resource)
			}
		}
	}
	if scanned == 0 {
		t.Fatal("no implementation files scanned")
	}
	t.Logf("Classified %d non-test resource/data-source files", scanned)
	p := alicloud.Provider().(*schema.Provider)
	resources := map[string]*schema.Resource{}
	for name, r := range p.DataSourcesMap {
		resources["data."+name] = r
	}
	for file, want := range map[string]string{
		"data_source_alicloud_ess_lifecyclehooks.go":               "data.alicloud_ess_lifecycle_hooks",
		"data_source_alicloud_yundun_dbaudit_instances.go":         "data.alicloud_yundun_dbaudit_instance",
		"resource_alicloud_ess_lifecyclehook.go":                   "alicloud_ess_lifecycle_hook",
		"resource_alicloud_edas_application_package_attachment.go": "alicloud_edas_application_deployment",
	} {
		direct, err := resourcesInFiles("../..", map[string]bool{"alicloud/" + file: true}, registry, registry)
		if err != nil {
			t.Fatal(err)
		}
		if len(direct) != 1 || !direct[want] {
			t.Errorf("%s resolved to %v, want %s", file, direct, want)
		}
		if strings.HasPrefix(want, "data.") {
			missing := coverageErrors(resources, direct, nil)
			if len(missing) == 0 || !strings.Contains(strings.Join(missing, "\n"), want+".ids") {
				t.Errorf("configurable ids silently missed: %v", missing)
			}
		}
	}
	if missing := coverageErrors(resources, map[string]bool{"data.alicloud_unresolved": true}, nil); len(missing) == 0 {
		t.Fatal("unresolved resource was silently skipped")
	}
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "alicloud"), 0755); err != nil {
		t.Fatal(err)
	}
	unknown := "alicloud/resource_alicloud_unresolved.go"
	if err := os.WriteFile(filepath.Join(root, unknown), []byte("package alicloud"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := resourcesInFiles(root, map[string]bool{unknown: true}, registry, registry); err == nil {
		t.Fatal("existing unregistered implementation was silently skipped")
	}
	if err := os.Remove(filepath.Join(root, unknown)); err != nil {
		t.Fatal(err)
	}
	if _, err := resourcesInFiles(root, map[string]bool{unknown: true}, registry, registry); err == nil {
		t.Fatal("unknown missing filename was mistaken for a removed registration")
	}
	previous := map[string]string{"alicloud_unresolved": "removedResource()"}
	if direct, err := resourcesInFiles(root, map[string]bool{unknown: true}, registry, previous); err != nil || len(direct) != 0 {
		t.Fatalf("removed file and registration retained: %v %v", direct, err)
	}
	moved := "alicloud/data_source_alicloud_ess_lifecyclehooks.go"
	if direct, err := resourcesInFiles(root, map[string]bool{moved: true}, registry, registry); err != nil || !direct["data.alicloud_ess_lifecycle_hooks"] {
		t.Fatalf("moved implementation with live registration was skipped: %v %v", direct, err)
	}
	if direct, err := resourcesInFiles(root, map[string]bool{"alicloud/resource_alicloud_unresolved_test.go": true}, registry, registry); err != nil || len(direct) != 0 {
		t.Fatalf("test file selected as implementation: %v %v", direct, err)
	}

}

func TestProviderRegistryChange(t *testing.T) {
	source, err := os.ReadFile("../../alicloud/provider.go")
	if err != nil {
		t.Fatal(err)
	}
	baseEntries, err := providerRegistry(source)
	if err != nil {
		t.Fatal(err)
	}
	modified := strings.Replace(string(source), "ResourcesMap: map[string]*schema.Resource{", "ResourcesMap: map[string]*schema.Resource{\n\"alicloud_collection_example\": resourceCollectionExample(),", 1)
	modified = strings.Replace(modified, baseEntries["alicloud_ram_group_membership"], "changedGroupMembership()", 1)
	headEntries, err := providerRegistry([]byte(modified))
	if err != nil {
		t.Fatal(err)
	}
	selected := map[string]bool{}
	for name, entry := range headEntries {
		if baseEntries[name] != entry {
			selected[name] = true
		}
	}
	want := map[string]bool{"alicloud_collection_example": true, "alicloud_ram_group_membership": true}
	if !reflect.DeepEqual(selected, want) {
		t.Fatalf("changed registrations: got %v, want %v", selected, want)
	}
	unrelated, err := providerRegistry([]byte(string(source) + "\nfunc newProviderBehavior() {}\n"))
	if err != nil || !reflect.DeepEqual(unrelated, baseEntries) {
		t.Fatalf("unrelated provider behavior changed registrations: %v", err)
	}
}

func TestRunRenamedImplementation(t *testing.T) {
	manifest, err := filepath.Abs("cases.json")
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "alicloud"), 0755); err != nil {
		t.Fatal(err)
	}
	for _, file := range []string{"provider.go", "data_source_alicloud_ess_lifecyclehooks.go"} {
		source, err := os.ReadFile(filepath.Join("../../alicloud", file))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "alicloud", file), source, 0644); err != nil {
			t.Fatal(err)
		}
	}
	t.Chdir(root)
	git := func(args ...string) string {
		t.Helper()
		args = append([]string{"-c", "commit.gpgsign=false", "-c", "user.name=Collection Order Test", "-c", "user.email=collection-order@example.com"}, args...)
		out, err := exec.Command("git", args...).CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
		return string(out)
	}
	git("init", "-q")
	git("config", "diff.renames", "true")
	git("add", ".")
	git("commit", "-qm", "fixture before move")
	git("mv", "alicloud/data_source_alicloud_ess_lifecyclehooks.go", "alicloud/moved_lifecyclehooks.go")
	git("commit", "-qm", "move implementation")
	if diff := git("diff", "--name-status", "HEAD^", "HEAD"); !strings.HasPrefix(diff, "R100") {
		t.Fatalf("fixture must exercise real rename detection: %s", diff)
	}
	if err := run("HEAD^", "HEAD", manifest); err == nil || !strings.Contains(err.Error(), "data.alicloud_ess_lifecycle_hooks.ids: missing collection-order fixture") {
		t.Fatalf("renamed live data source escaped fixture coverage: %v", err)
	}
}
