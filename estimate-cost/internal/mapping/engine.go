package mapping

import (
	"fmt"
	"os"
	"strings"
)

// MissingField means a required field could not be resolved from the plan,
// has no FallbackEnv, and no Default. Callers should treat the resource as
// unknown/incomplete and skip pricing for it.
type MissingField struct {
	Param string
	Path  string
}

func (e MissingField) Error() string {
	return fmt.Sprintf("required param %q (from %s) could not be resolved", e.Param, e.Path)
}

// Build applies a PricingTarget's rules to the plan's change.after / change.before
// maps and returns the corresponding CC API parameters.
//
// A rule using "$.x" reads from after; "$before.x" reads from before (which
// holds the refreshed prior values on an update plan). before may be nil
// (e.g. on a create plan).
func Build(t *PricingTarget, after, before map[string]interface{}) (map[string]interface{}, error) {
	out := map[string]interface{}{}
	for name, def := range t.Params {
		// Evaluate the per-param When predicate: every {path: value} pair must hold.
		if !matchWhen(def.When, after, before) {
			continue
		}

		switch {
		case def.Expand == "indexed":
			if err := expandIndexed(out, name, def, after, before); err != nil {
				return nil, err
			}
		default:
			val, ok := resolveScalar(def, after, before)
			if !ok {
				if def.Required {
					return nil, MissingField{Param: name, Path: def.From}
				}
				continue
			}
			if def.WrapArray {
				// Some pricing inputs are list-typed even for a single value
				// (e.g. ModifyInstanceChargeType's InstanceIds).
				val = []interface{}{val}
			}
			out[name] = val
		}
	}
	nestPricingContext(out)
	return out, nil
}

// nestPricingContext converts flat "PricingContext.a.b" keys into a nested
// object under the "PricingContext" key. GetApiPrice's PricingContext is a
// structured JSON object (unlike POP parameters such as "SystemDisk.Size",
// which stay flat) — CC API rejects the dotted-flat form for it. Mapping
// files write dotted keys for readability; this rewrites them on the way out.
func nestPricingContext(params map[string]interface{}) {
	const prefix = "PricingContext."
	var nested map[string]interface{}
	for k, v := range params {
		if !strings.HasPrefix(k, prefix) {
			continue
		}
		if nested == nil {
			nested = map[string]interface{}{}
		}
		cur := nested
		segs := strings.Split(k[len(prefix):], ".")
		for i, seg := range segs {
			if i == len(segs)-1 {
				cur[seg] = v
				break
			}
			next, ok := cur[seg].(map[string]interface{})
			if !ok {
				next = map[string]interface{}{}
				cur[seg] = next
			}
			cur = next
		}
		delete(params, k)
	}
	if nested != nil {
		params["PricingContext"] = nested
	}
}

// resolveScalar resolves a scalar field, in priority order:
//
//	const → from → fallbackEnv → default
//
// "$.x" reads from after, "$before.x" reads from before (before may be nil).
// Returns (value, true) on success, (nil, false) when nothing matched.
func resolveScalar(def ParamDef, after, before map[string]interface{}) (interface{}, bool) {
	if def.Const != nil {
		return def.Const, true
	}
	if def.From != "" {
		source, path := after, def.From
		if strings.HasPrefix(path, "$before.") {
			source = before
			path = "$." + path[len("$before."):]
		}
		if source != nil {
			if v, ok := getByPath(source, path); ok && !isEmpty(v) {
				return v, true
			}
		}
	}
	if def.FallbackEnv != "" {
		if v := os.Getenv(def.FallbackEnv); v != "" {
			return v, true
		}
	}
	if def.Default != nil {
		return def.Default, true
	}
	return nil, false
}

// expandIndexed flattens a list field into indexed parameters X.1.A, X.1.B,
// X.2.A, ... def.From must point to a list field in after; def.Fields gives
// the per-element sub-field mapping.
func expandIndexed(out map[string]interface{}, prefix string, def ParamDef, after, before map[string]interface{}) error {
	if def.From == "" || def.Fields == nil {
		return nil
	}
	raw, ok := getByPath(after, def.From)
	if !ok {
		return nil
	}
	list, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	for i, item := range list {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		// Build a synthetic "after" for sub-fields: "$item.foo" resolves to itemMap[foo].
		ctx := map[string]interface{}{"item": itemMap}
		for subName, subDef := range def.Fields {
			// Rewrite "$item.x" → "$.item.x" so the same getByPath works.
			rewritten := subDef
			rewritten.From = strings.Replace(subDef.From, "$item.", "$.item.", 1)
			val, ok := resolveScalar(rewritten, ctx, before)
			if !ok {
				if subDef.Required {
					return MissingField{Param: fmt.Sprintf("%s.%d.%s", prefix, i+1, subName), Path: subDef.From}
				}
				continue
			}
			out[fmt.Sprintf("%s.%d.%s", prefix, i+1, subName)] = val
		}
	}
	return nil
}

// matchWhen evaluates the When predicate: every {path: expected} pair must
// hold. An empty/nil When map is treated as true.
//
// Paths starting with "$." read from after; paths starting with "$before."
// read from before — useful for "field X changed from A to B" branching
// (e.g. detecting a PostPaid → PrePaid charge-type switch). before may be
// nil on create plans, in which case $before.* paths fail to resolve and
// the predicate returns false.
func matchWhen(when map[string]string, after, before map[string]interface{}) bool {
	if len(when) == 0 {
		return true
	}
	for path, expected := range when {
		source := after
		lookupPath := path
		if strings.HasPrefix(path, "$before.") {
			source = before
			lookupPath = "$." + path[len("$before."):]
		}
		if source == nil {
			return false
		}
		v, ok := getByPath(source, lookupPath)
		if !ok {
			// The When path is not present in the plan: we treat that as a
			// non-match (the simpler, more conservative behaviour, e.g. for
			// instance_charge_type-style discriminators).
			return false
		}
		if fmt.Sprint(v) != expected {
			return false
		}
	}
	return true
}

// getByPath is an intentionally minimal JSONPath implementation: only
// "$.a.b.c" style pure-dot paths are supported.
func getByPath(root map[string]interface{}, path string) (interface{}, bool) {
	if !strings.HasPrefix(path, "$.") {
		return nil, false
	}
	cur := interface{}(root)
	for _, seg := range strings.Split(path[2:], ".") {
		m, ok := cur.(map[string]interface{})
		if !ok {
			return nil, false
		}
		cur, ok = m[seg]
		if !ok {
			return nil, false
		}
	}
	return cur, true
}

// isEmpty classifies "unset in plan" values. In Terraform plan JSON, numbers
// default to 0 and strings default to "". We only filter obviously empty
// strings and nil; zero-valued numbers are left for per-param Default to
// handle (avoiding false negatives on legitimate zeros).
func isEmpty(v interface{}) bool {
	if v == nil {
		return true
	}
	if s, ok := v.(string); ok && s == "" {
		return true
	}
	return false
}
