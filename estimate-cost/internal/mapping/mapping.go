// Package mapping translates a resource's `change.after` block in a Terraform
// plan JSON into the (popCode, popVersion, apiName, parameters) tuple required
// by CC API GetApiPrice.
//
// Mapping rules come from mappings/*.json — one file per resource type. The
// supported fields in a rule are:
//
//	from / fallbackEnv / required / default / const / when / expand+fields
//
// See the estimate-cost README for the full field semantics.
package mapping

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// File is the decoded representation of a mapping JSON file.
type File struct {
	Resource       string           `json:"resource"`
	Billable       bool             `json:"billable"`
	PricingTargets []*PricingTarget `json:"pricingTargets"`
}

// PricingTarget describes "when these plan actions are matched, call this
// OpenAPI with these parameter mappings". When is evaluated at the target
// level to support conditional dispatch — e.g. for the same `update` action,
// route to ModifyInstanceSpec for PostPaid and ModifyPrepayInstanceSpec for
// PrePaid based on `instance_charge_type`.
//
// WhenChanged narrows the match further by requiring at least one of the
// listed schema fields to differ between change.before and change.after.
// This prevents false-positive triggers like "user updated only the disk
// size, but the spec-change target also fired because instance_charge_type
// matched". The strings are bare field names (no "$." prefix). When the
// list is empty/nil the constraint is treated as satisfied.
type PricingTarget struct {
	Name        string              `json:"name"`
	Actions     []string            `json:"actions"`
	When        map[string]string   `json:"when,omitempty"`
	WhenChanged []string            `json:"whenChanged,omitempty"`
	OpenAPI     OpenAPIRef          `json:"openapi"`
	Params      map[string]ParamDef `json:"params"`
}

// OpenAPIRef is the triple required to address a single OpenAPI operation
// inside CC API GetApiPrice.
type OpenAPIRef struct {
	PopCode    string `json:"popCode"`
	PopVersion string `json:"popVersion"`
	APIName    string `json:"apiName"`
}

// ParamDef describes how to resolve a single OpenAPI parameter. All JSON
// fields are omitempty; a single definition typically uses only a few of them.
type ParamDef struct {
	From        string              `json:"from,omitempty"`        // JSONPath like "$.field"
	FallbackEnv string              `json:"fallbackEnv,omitempty"` // env var to fall back to when From is unresolvable
	Default     interface{}         `json:"default,omitempty"`     // value to use when both From and FallbackEnv fail
	Const       interface{}         `json:"const,omitempty"`       // hard-coded value (overrides From)
	Required    bool                `json:"required,omitempty"`    // fail if unresolved
	When        map[string]string   `json:"when,omitempty"`        // include only when {"$.x": "V"} (i.e. x==V) holds
	Expand      string              `json:"expand,omitempty"`      // "indexed" expands a list to X.1.A, X.1.B …
	Fields      map[string]ParamDef `json:"fields,omitempty"`      // per-element field map used together with Expand=indexed
	WrapArray   bool                `json:"wrapArray,omitempty"`   // wrap the resolved scalar in a single-element list
}

// Load reads every *.json mapping file in dir and returns a resource-type
// index. It is a thin wrapper over LoadFS — production builds use LoadFS
// directly against the embedded FS.
func Load(dir string) (map[string]*File, error) {
	return LoadFS(os.DirFS(dir), ".")
}

// LoadFS reads mapping files from an arbitrary fs.FS (e.g. an embed.FS for the
// production build, or os.DirFS for development). Production and development
// share the same parsing logic this way.
func LoadFS(fsys fs.FS, dir string) (map[string]*File, error) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, fmt.Errorf("read mapping dir %s: %w", dir, err)
	}
	out := map[string]*File{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		raw, err := fs.ReadFile(fsys, path)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", path, err)
		}
		var f File
		if err := json.Unmarshal(raw, &f); err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}
		if f.Resource == "" {
			return nil, fmt.Errorf("%s: missing 'resource' field", path)
		}
		out[f.Resource] = &f
	}
	return out, nil
}

// PickTargetByAction is a convenience wrapper over PickTarget. It is mainly
// used by the fallback path where an update is not matched by any modify
// target and must fall back to the create target for a double-call diff.
func (f *File) PickTargetByAction(action string, after, before map[string]interface{}) *PricingTarget {
	return f.PickTarget([]string{action}, after, before)
}

// PickTarget returns the first PricingTarget whose actions intersect with the
// requested actions and whose target-level When predicate holds. Used for
// "one action → one OpenAPI call" cases such as create.
//
// before may be nil on create plans; targets whose When references $before.*
// will not match in that case.
func (f *File) PickTarget(actions []string, after, before map[string]interface{}) *PricingTarget {
	for _, t := range f.PricingTargets {
		if !actionsMatch(t.Actions, actions) {
			continue
		}
		if !matchWhen(t.When, after, before) {
			continue
		}
		if !matchWhenChanged(t.WhenChanged, after, before) {
			continue
		}
		return t
	}
	return nil
}

// PickTargets returns every matching target. Used for "one action may trigger
// multiple OpenAPI calls" cases — e.g. a PrePaid update that changes both
// instance_type and system_disk_size will call ModifyPrepayInstanceSpec and
// ResizeDisk separately, producing two separate billing orders.
func (f *File) PickTargets(actions []string, after, before map[string]interface{}) []*PricingTarget {
	var out []*PricingTarget
	for _, t := range f.PricingTargets {
		if !actionsMatch(t.Actions, actions) {
			continue
		}
		if !matchWhen(t.When, after, before) {
			continue
		}
		if !matchWhenChanged(t.WhenChanged, after, before) {
			continue
		}
		out = append(out, t)
	}
	return out
}

// matchWhenChanged returns true when at least one of the listed schema field
// names has a different value in after vs before. Used to gate "update_*"
// targets so that, e.g., a disk-only resize doesn't also fire the spec-
// change target. Empty/nil list passes through (no constraint).
func matchWhenChanged(fields []string, after, before map[string]interface{}) bool {
	if len(fields) == 0 {
		return true
	}
	if before == nil {
		return false
	}
	for _, f := range fields {
		if fmt.Sprint(after[f]) != fmt.Sprint(before[f]) {
			return true
		}
	}
	return false
}

func actionsMatch(target, actual []string) bool {
	for _, a := range actual {
		for _, ta := range target {
			if a == ta {
				return true
			}
		}
	}
	return false
}
