package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
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
		if drift, err := checkCase(r, c); drift != "" || err != nil {
			t.Fatalf("fixed pair rejected: drift=%q err=%v", drift, err)
		}
		r.Schema["values"].DiffSuppressFunc = nil
		if drift, err := checkCase(r, c); drift == "" || err != nil {
			t.Fatalf("fixed pair reorder was not exercised: drift=%q err=%v", drift, err)
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
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "alicloud"), 0755); err != nil {
		t.Fatal(err)
	}
	for name, source := range map[string]string{
		"service_alicloud_example.go":       "package alicloud; func normalize() {}",
		"flatten.go":                        "package alicloud; func flatten() { normalize() }",
		"resource_alicloud_example.go":      "package alicloud; func readExample() { flatten() }",
		"data_source_alicloud_example.go":   "package alicloud; func readData() { normalize() }",
		"resource_alicloud_other.go":        "package alicloud; func readOther() {}",
		"resource_alicloud_example_test.go": "package alicloud; func testExample() { readOther() }",
	} {
		if err := os.WriteFile(filepath.Join(root, "alicloud", name), []byte(source), 0644); err != nil {
			t.Fatal(err)
		}
	}
	changed := map[string]bool{"alicloud/service_alicloud_example.go": true}
	got, err := affectedResources(root, changed)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]bool{"alicloud_example": true, "data.alicloud_example": true}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("transitive normalization change: got %v, want %v", got, want)
	}
	r := &schema.Resource{Schema: map[string]*schema.Schema{"values": {
		Type: schema.TypeSet, Optional: true, Elem: &schema.Schema{Type: schema.TypeString},
	}}}
	resources := map[string]*schema.Resource{"alicloud_example": r, "data.alicloud_example": r}
	if missing := coverageErrors(resources, resourcesInFiles(changed), nil); len(missing) != 0 {
		t.Fatalf("shared helper demanded new fixtures: %v", missing)
	}
	changed["alicloud/resource_alicloud_example.go"] = true
	changed["alicloud/data_source_alicloud_example.go"] = true
	changed["alicloud/resource_alicloud_other_test.go"] = true
	if direct := resourcesInFiles(changed); !reflect.DeepEqual(direct, want) || len(coverageErrors(resources, direct, nil)) != 2 {
		t.Fatalf("direct resource/data-source changes escaped coverage: %v", direct)
	}
	actual, err := affectedResources("../..", map[string]bool{"alicloud/service_alicloud_ecd_v2.go": true})
	if err != nil {
		t.Fatal(err)
	}
	if !actual["alicloud_ecd_desktop_group"] || actual["alicloud_ram_group_membership"] {
		t.Fatalf("ECD normalization consumer resolution is wrong: %d consumers, ECD selected=%t, unrelated RAM selected=%t", len(actual), actual["alicloud_ecd_desktop_group"], actual["alicloud_ram_group_membership"])
	}
	t.Logf("ECD shared normalization change selects %d resource/data-source consumers", len(actual))
	direct, err := affectedResources("../..", map[string]bool{"alicloud/resource_alicloud_ecd_desktop_group.go": true})
	if err != nil {
		t.Fatal(err)
	}
	if len(direct) != 1 || !direct["alicloud_ecd_desktop_group"] {
		t.Fatalf("direct resource change selected %d consumers", len(direct))
	}
	t.Logf("Direct ECD resource change selects %d consumer", len(direct))
}

func TestProviderRegistryChange(t *testing.T) {
	source, err := os.ReadFile("../../alicloud/provider.go")
	if err != nil {
		t.Fatal(err)
	}
	baseBody, baseEntries, err := providerRegistry(source)
	if err != nil {
		t.Fatal(err)
	}
	modified := strings.Replace(string(source), "ResourcesMap: map[string]*schema.Resource{", "ResourcesMap: map[string]*schema.Resource{\n\"alicloud_collection_example\": resourceCollectionExample(),", 1)
	headBody, headEntries, err := providerRegistry([]byte(modified))
	if err != nil {
		t.Fatal(err)
	}
	if baseBody != headBody {
		t.Fatal("plain resource registration treated as a global provider change")
	}
	changed := map[string]bool{"alicloud/resource_alicloud_collection_example.go": true}
	selected := resourcesInFiles(changed)
	for name, entry := range headEntries {
		if baseEntries[name] != entry {
			selected[name] = true
		}
	}
	if len(selected) != 1 || !selected["alicloud_collection_example"] {
		t.Fatalf("new resource registration selected %d resources", len(selected))
	}
	t.Logf("New resource plus real provider registration selects %d resource", len(selected))
	globalBody, _, err := providerRegistry([]byte(modified + "\nfunc newProviderBehavior() {}\n"))
	if err != nil {
		t.Fatal(err)
	}
	if globalBody == baseBody {
		t.Fatal("provider behavior change was mistaken for registration-only")
	}
}
