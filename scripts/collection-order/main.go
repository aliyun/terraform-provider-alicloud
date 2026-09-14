// collection-order checks collection ordering through the provider's SDK diff.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

type testCase struct {
	Resource      string                 `json:"resource"`
	Path          string                 `json:"path"`
	Order         string                 `json:"order"`
	Reason        string                 `json:"reason"`
	Config        map[string]interface{} `json:"config"`
	Members       []interface{}          `json:"members"`
	NilMetaReason string                 `json:"nil_meta_reason,omitempty"`
}

func (c testCase) key() string { return c.Resource + "." + schemaPath(c.Path) }

// These exact legacy defects were reproduced with SDK Diff. They are reported,
// never silently skipped, and cannot excuse a change to an affected resource.
var knownDrift = map[string]bool{
	"alicloud_ecd_desktop.end_user_ids":       true,
	"alicloud_ecd_desktop_group.end_user_ids": true,
}

// SDK v1 recalculates these unset outputs when the known desktop ordering
// defect forces replacement. No other sibling changes belong to that baseline.
var knownReplacementOutputs = map[string]bool{
	"desktop_type": true, "payment_type": true, "status": true, "user_assign_mode": true,
}

func main() {
	base := flag.String("base", "HEAD", "base commit for affected-resource coverage")
	head := flag.String("head", "", "head commit; empty includes working-tree changes")
	manifest := flag.String("cases", "scripts/collection-order/cases.json", "collection fixtures")
	flag.Parse()
	log.SetOutput(io.Discard) // SDK logs are noisy; failures print the actual attribute diff.
	if err := run(*base, *head, *manifest); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(base, head, manifest string) error {
	cases, err := loadCases(manifest)
	if err != nil {
		return err
	}
	args := []string{"diff", "--no-renames", "--name-only", "--diff-filter=ACDMRT", base}
	if head != "" {
		args = append(args, head)
	}
	args = append(args, "--", "alicloud")
	// No shell expansion, and revision arguments may not be git options.
	if base == "" || strings.HasPrefix(base, "-") || strings.HasPrefix(head, "-") {
		return fmt.Errorf("invalid change range")
	}
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return fmt.Errorf("resolve changed files: %w", err)
	}
	changed := map[string]bool{}
	for _, name := range strings.Fields(string(out)) {
		changed[name] = true
	}
	current, err := os.ReadFile("alicloud/provider.go")
	if err != nil {
		return err
	}
	registry, err := providerRegistry(current)
	if err != nil {
		return err
	}
	previousRegistry := registry
	if changed["alicloud/provider.go"] {
		previous, err := exec.Command("git", "show", base+":alicloud/provider.go").Output()
		if err != nil {
			return fmt.Errorf("read base provider registry: %w", err)
		}
		previousRegistry, err = providerRegistry(previous)
		if err != nil {
			return err
		}
	}
	coverage, err := resourcesInFiles(".", changed, registry, previousRegistry)
	if err != nil {
		return err
	}
	for name, entry := range registry {
		if previousRegistry[name] != entry {
			coverage[name] = true
		}
	}
	affected := affectedResources(changed, coverage)
	p := alicloud.Provider().(*schema.Provider)
	resources := map[string]*schema.Resource{}
	for name, r := range p.ResourcesMap {
		resources[name] = r
	}
	for name, r := range p.DataSourcesMap {
		resources["data."+name] = r
	}
	errors := coverageErrors(resources, coverage, cases)
	passed, legacy := 0, 0
	for _, c := range cases {
		r := resources[c.Resource]
		if r == nil {
			errors = append(errors, c.key()+": resource no longer exists; remove the stale case")
			continue
		}
		drift, err := checkCase(r, c)
		if err != nil {
			errors = append(errors, c.key()+": "+err.Error())
			continue
		}
		if knownDrift[c.key()] {
			if drift == "" {
				errors = append(errors, c.key()+": legacy drift is fixed; remove its knownDrift entry")
			} else if affected[c.Resource] {
				errors = append(errors, c.key()+": affected legacy defect must be fixed, not baselined: "+drift)
			} else {
				legacy++
				fmt.Printf("::warning::KNOWN DRIFT %s: %s\n", c.key(), drift)
			}
		} else if drift != "" {
			errors = append(errors, c.key()+": "+drift)
		} else {
			passed++
			fmt.Printf("PASS %s\n", c.key())
		}
	}
	fmt.Printf("Collection order: %d passed, %d known drifts, %d failures; %d affected resources\n", passed, legacy, len(errors), len(affected))
	if len(errors) != 0 {
		return fmt.Errorf("%s", strings.Join(errors, "\n"))
	}
	if passed == 0 {
		return fmt.Errorf("no passing collection checks; refusing an empty or baseline-only success")
	}
	return nil
}

func loadCases(path string) ([]testCase, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cases []testCase
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cases); err != nil {
		return nil, err
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("collection fixture registry is empty")
	}
	seen := map[string]bool{}
	for _, c := range cases {
		if seen[c.key()] || c.Resource == "" || c.Path == "" || strings.TrimSpace(c.Reason) == "" {
			return nil, fmt.Errorf("%s: duplicate case or missing resource/path/reason", c.key())
		}
		seen[c.key()] = true
		if len(c.Members) < 2 {
			return nil, fmt.Errorf("%s: provide at least two distinct members for reorder and membership controls", c.key())
		}
		for i, member := range c.Members {
			for _, previous := range c.Members[:i] {
				if reflect.DeepEqual(member, previous) {
					return nil, fmt.Errorf("%s: members must be distinct", c.key())
				}
			}
		}
	}
	for key := range knownDrift {
		if !seen[key] {
			return nil, fmt.Errorf("%s: known drift must retain an executable case", key)
		}
	}
	return cases, nil
}

// Directly edited resources and added/changed registrations require fixtures.
// Audited legacy dependencies only affect known-drift enforcement, not coverage.
func coverageErrors(resources map[string]*schema.Resource, affected map[string]bool, cases []testCase) []string {
	registered := map[string]bool{}
	for _, c := range cases {
		registered[c.key()] = true
	}
	var errors []string
	for name := range affected {
		r := resources[name]
		if r == nil {
			errors = append(errors, name+": selected resource is not registered")
			continue
		}
		for path := range collections(r.Schema, "") {
			if !registered[name+"."+path] {
				errors = append(errors, name+"."+path+": missing collection-order fixture; declare list ordering and provide at least two distinct members")
			}
		}
	}
	sort.Strings(errors)
	return errors
}

func collections(fields map[string]*schema.Schema, prefix string) map[string]*schema.Schema {
	found := map[string]*schema.Schema{}
	for name, s := range fields {
		if (!s.Optional && !s.Required) || s.Removed != "" {
			continue
		}
		path := prefix + name
		if s.Type == schema.TypeList || s.Type == schema.TypeSet {
			if s.MaxItems != 1 {
				found[path] = s
			}
			if nested, ok := s.Elem.(*schema.Resource); ok {
				for key, value := range collections(nested.Schema, path+".") {
					found[key] = value
				}
			}
		}
	}
	return found
}

func schemaPath(path string) string {
	var parts []string
	for _, part := range strings.Split(path, ".") {
		if _, err := strconv.Atoi(part); err != nil {
			parts = append(parts, part)
		}
	}
	return strings.Join(parts, ".")
}

func checkCase(r *schema.Resource, c testCase) (string, error) {
	if r.CustomizeDiff != nil && strings.TrimSpace(c.NilMetaReason) == "" {
		return "", fmt.Errorf("CustomizeDiff requires an audited nil_meta_reason explaining why this fixture needs no configured client")
	}
	s := collections(r.Schema, "")[schemaPath(c.Path)]
	if s == nil {
		return "", fmt.Errorf("not a configurable multi-member collection; update the stale fixture")
	}
	ordered := c.Order == "ordered"
	if s.Type == schema.TypeSet {
		if c.Order != "" && c.Order != "unordered" {
			return "", fmt.Errorf("TypeSet cannot be declared ordered")
		}
	} else if c.Order != "unordered" && !ordered {
		return "", fmt.Errorf("TypeList must explicitly declare ordered or unordered semantics")
	}
	if len(c.Members) < 2 {
		return "", fmt.Errorf("at least two members required")
	}
	size := max(2, s.MinItems, len(c.Members)-1)
	if s.MaxItems > 0 {
		size = min(size, s.MaxItems)
	}
	if size > len(c.Members) {
		return "", fmt.Errorf("provide enough distinct members to satisfy MinItems=%d", s.MinItems)
	}
	canonical := c.Members[:size]
	background, err := collectionDiff(r, c, canonical, canonical)
	if err != nil {
		return "", err
	}
	for key, value := range background {
		if !value.NewComputed || inCollection(key, c.Path) {
			return "", fmt.Errorf("unchanged fixture already has a diff: %s", formatDiff(background))
		}
	}
	reverse := append([]interface{}{}, canonical...)
	for i, j := 0, len(reverse)-1; i < j; i, j = i+1, j-1 {
		reverse[i], reverse[j] = reverse[j], reverse[i]
	}
	additionState, added := canonical[:size-1], canonical
	if size < len(c.Members) && (s.MaxItems == 0 || size < s.MaxItems) {
		additionState = canonical
		added = c.Members[:size+1]
	}
	var replaced []interface{}
	replacementState := canonical
	if size == len(c.Members) && size > s.MinItems {
		replacementState = canonical[:size-1]
	}
	if len(replacementState) < len(c.Members) {
		replaced = append([]interface{}{}, replacementState...)
		replaced[len(replaced)-1] = c.Members[len(replacementState)]
	}
	// Canonical members model a normalized API result. This does not execute
	// Read: it tests SDK planning against that state without sorting config.
	type scenario struct {
		name          string
		state, config []interface{}
		wantDiff      bool
		reorder       bool
	}
	scenarios := []scenario{
		{"unchanged", canonical, canonical, false, false},
		{"reordered configuration", canonical, reverse, ordered, true},
		{"reordered state", reverse, canonical, ordered, true},
		{"add", additionState, added, true, false},
		{"remove", canonical, canonical[:len(canonical)-1], true, false},
		{"replace", replacementState, replaced, true, false},
		{"add in state", added, additionState, true, false},
		{"remove in state", canonical[:len(canonical)-1], canonical, true, false},
		{"replace in state", replaced, replacementState, true, false},
	}
	if len(canonical) >= 3 {
		rotated := append(append([]interface{}{}, canonical[1:]...), canonical[0])
		scenarios = append(scenarios, scenario{"rotated configuration", canonical, rotated, ordered, true})
	}
	var drift []string
	for _, scenario := range scenarios {
		if scenario.state == nil || scenario.config == nil {
			return "", fmt.Errorf("%s: missing executable replacement control; provide an additional distinct legal member", scenario.name)
		}
		if min(len(scenario.state), len(scenario.config)) < s.MinItems ||
			(s.MaxItems > 0 && max(len(scenario.state), len(scenario.config)) > s.MaxItems) {
			fmt.Printf("N/A %s %s: size change forbidden by MinItems=%d MaxItems=%d\n", c.key(), scenario.name, s.MinItems, s.MaxItems)
			continue
		}
		diff, err := collectionDiff(r, c, scenario.state, scenario.config)
		if err != nil {
			return "", fmt.Errorf("%s: %w", scenario.name, err)
		}
		// Only identical computed unknowns from the unchanged control are noise.
		// Retain every new parent/sibling diff caused by reordering the collection.
		for key, value := range background {
			if reflect.DeepEqual(diff[key], value) {
				delete(diff, key)
			}
		}
		for key, value := range diff {
			if value.Old == value.New && !value.NewComputed && !value.NewRemoved && !value.RequiresNew {
				delete(diff, key) // SDK v1 repeats unchanged inputs on a replacement.
			}
		}
		targetChanged := false
		for key, value := range diff {
			if inCollection(key, c.Path) {
				targetChanged = true
			} else if knownDrift[c.key()] {
				if c.key() == "alicloud_ecd_desktop.end_user_ids" && knownReplacementOutputs[key] && value.NewComputed && value.Old == "" && value.New == "" && !value.NewRemoved && !value.RequiresNew {
					continue
				}
				return "", fmt.Errorf("%s: unexpected diff outside the legacy attribute: %s", scenario.name, formatDiff(diff))
			}
		}
		if scenario.wantDiff && len(diff) != 0 && !targetChanged {
			return "", fmt.Errorf("%s: membership/order diff is missing from the collection itself: %s", scenario.name, formatDiff(diff))
		}
		if (len(diff) != 0) != scenario.wantDiff {
			message := fmt.Sprintf("%s: want diff=%t, got %s", scenario.name, scenario.wantDiff, formatDiff(diff))
			if scenario.reorder && !ordered && len(diff) != 0 {
				drift = append(drift, message)
			} else {
				return "", fmt.Errorf("%s", message)
			}
		}
	}
	return strings.Join(drift, "; "), nil
}

func collectionDiff(r *schema.Resource, c testCase, state, config []interface{}) (map[string]*terraform.ResourceAttrDiff, error) {
	stateConfig, err := fixtureConfig(c, state)
	if err != nil {
		return nil, err
	}
	configMap, err := fixtureConfig(c, config)
	if err != nil {
		return nil, err
	}
	for _, config := range []map[string]interface{}{stateConfig, configMap} {
		if _, errors := r.Validate(terraform.NewResourceConfigRaw(config)); len(errors) != 0 {
			return nil, fmt.Errorf("invalid SDK fixture: %v", errors)
		}
	}
	d, err := schema.InternalMap(r.Schema).Data(nil, nil)
	if err != nil {
		return nil, err
	}
	d.SetId("collection-order-check")
	for key, value := range stateConfig {
		if err := d.Set(key, value); err != nil {
			return nil, fmt.Errorf("state %s: %w", key, err)
		}
	}
	// Exercise the real resource schema, field suppressors and CustomizeDiff.
	// Configure/Create/Read/Update/Delete are deliberately never called.
	diff, err := r.Diff(d.State(), terraform.NewResourceConfigRaw(configMap), nil)
	if err != nil || diff == nil {
		return nil, err
	}
	return diff.Attributes, nil
}

func inCollection(key, path string) bool {
	key, path = schemaPath(key), schemaPath(path)
	return key == path || strings.HasPrefix(key, path+".")
}

func formatDiff(diff map[string]*terraform.ResourceAttrDiff) string {
	var changes []string
	for key, value := range diff {
		changes = append(changes, fmt.Sprintf("%s: %q -> %q (requires_new=%t, computed=%t)", key, value.Old, value.New, value.RequiresNew, value.NewComputed))
	}
	sort.Strings(changes)
	return strings.Join(changes, ", ")
}

// Compare registration values so unrelated provider edits require no new fixtures.
func providerRegistry(source []byte) (map[string]string, error) {
	file, err := parser.ParseFile(token.NewFileSet(), "provider.go", source, 0)
	if err != nil {
		return nil, err
	}
	entries := map[string]string{}
	ast.Inspect(file, func(node ast.Node) bool {
		field, ok := node.(*ast.KeyValueExpr)
		if !ok {
			return true
		}
		key, ok := field.Key.(*ast.Ident)
		if !ok || (key.Name != "ResourcesMap" && key.Name != "DataSourcesMap") {
			return true
		}
		registry, ok := field.Value.(*ast.CompositeLit)
		if !ok {
			return true
		}
		for _, element := range registry.Elts {
			entry, ok := element.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			literal, ok := entry.Key.(*ast.BasicLit)
			if !ok {
				continue
			}
			name, err := strconv.Unquote(literal.Value)
			if err != nil {
				continue
			}
			if key.Name == "DataSourcesMap" {
				name = "data." + name
			}
			var value bytes.Buffer
			if err := format.Node(&value, token.NewFileSet(), entry.Value); err != nil {
				continue
			}
			entries[name] = value.String()
		}
		return false
	})
	if len(entries) == 0 {
		return nil, fmt.Errorf("provider registry could not be resolved")
	}
	return entries, nil
}

func fixtureConfig(c testCase, members []interface{}) (map[string]interface{}, error) {
	// Clone through JSON so neither ResourceData nor SDK set hashing can mutate
	// another scenario. Fixtures use JSON-compatible values only.
	raw, err := json.Marshal(c.Config)
	if err != nil {
		return nil, err
	}
	config := map[string]interface{}{}
	if err := json.Unmarshal(raw, &config); err != nil {
		return nil, err
	}
	if config == nil {
		config = map[string]interface{}{}
	}
	var current interface{} = config
	parts := strings.Split(c.Path, ".")
	for _, part := range parts[:len(parts)-1] {
		switch value := current.(type) {
		case map[string]interface{}:
			current = value[part]
		case []interface{}:
			i, err := strconv.Atoi(part)
			if err != nil || i < 0 || i >= len(value) {
				return nil, fmt.Errorf("invalid fixture index in %s", c.Path)
			}
			current = value[i]
		default:
			return nil, fmt.Errorf("config must contain parent blocks for %s", c.Path)
		}
	}
	parent, ok := current.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("config must contain parent blocks for %s", c.Path)
	}
	memberJSON, err := json.Marshal(members)
	if err != nil {
		return nil, err
	}
	var clonedMembers []interface{}
	if err := json.Unmarshal(memberJSON, &clonedMembers); err != nil {
		return nil, err
	}
	parent[parts[len(parts)-1]] = clonedMembers
	return config, nil
}

// All fixtures run regardless of these selections. Only audited dependencies
// may promote a known drift to a failure for an indirectly changed resource.
func affectedResources(changed, direct map[string]bool) map[string]bool {
	affected := map[string]bool{}
	for name := range direct {
		affected[name] = true
	}
	// ponytail: use product-file boundaries; update these if the Read helpers move.
	for file, resource := range map[string]string{
		// DescribeEcdDesktop returns EndUserIds consumed by desktop Read.
		"alicloud/service_alicloud_ecd.go": "alicloud_ecd_desktop",
		// DescribeDesktopGroupDescribeUsersInGroup sorts EndUserIds for group Read.
		"alicloud/service_alicloud_ecd_v2.go": "alicloud_ecd_desktop_group",
	} {
		if changed[file] {
			affected[resource] = true
		}
	}
	return affected
}

func resourcesInFiles(root string, files map[string]bool, registry, previous map[string]string) (map[string]bool, error) {
	affected := map[string]bool{}
	// ponytail: explicit filename exceptions; unresolved existing files fail closed.
	aliases := map[string]string{
		"data.alicloud_ess_lifecyclehooks":             "data.alicloud_ess_lifecycle_hooks",
		"data.alicloud_yundun_dbaudit_instances":       "data.alicloud_yundun_dbaudit_instance",
		"alicloud_ess_lifecyclehook":                   "alicloud_ess_lifecycle_hook",
		"alicloud_edas_application_package_attachment": "alicloud_edas_application_deployment",
	}
	for file := range files {
		if filepath.Dir(file) != "alicloud" || !strings.HasSuffix(file, ".go") || strings.HasSuffix(file, "_test.go") {
			continue
		}
		if file == "alicloud/data_source_alicloud_common.go" {
			continue // Shared hash/schema helpers, not a data-source implementation.
		}
		name := strings.TrimSuffix(filepath.Base(file), ".go")
		if strings.HasPrefix(name, "resource_alicloud_") {
			name = strings.TrimPrefix(name, "resource_")
		} else if strings.HasPrefix(name, "data_source_alicloud_") {
			name = "data." + strings.TrimPrefix(name, "data_source_")
		} else {
			continue
		}
		if alias := aliases[name]; alias != "" {
			name = alias
		}
		if registry[name] == "" {
			_, err := os.Stat(filepath.Join(root, file))
			if os.IsNotExist(err) && previous[name] != "" {
				continue // Confirmed removal of both the file and its former registration.
			}
			if err != nil && !os.IsNotExist(err) {
				return nil, err
			}
			return nil, fmt.Errorf("%s: cannot resolve registered resource %s; add a filename mapping", file, name)
		}
		affected[name] = true // Keep coverage when the implementation moved but is still registered.
	}
	return affected, nil
}
