// Framework wire and state adapters shared by the framework-native resources. The helpers are
// pure functions over request maps, response maps and framework attribute values; keeping one
// copy keeps behaviour identical across product packages instead of relying on per-package
// generated duplicates.
package fwadapt

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WireField describes one request or response field mapping entry.
type WireField struct {
	Name   string
	Single bool
	Read   bool
	Input  bool
	Fields map[string]WireField
}

// Mapping binds a framework attribute path to a request field.
type Mapping struct {
	Source, Target, Encoding, Constant string
	Group, GroupKey, GroupValue        string
	Repeat                             bool
	IndexSources                       []int
	Values                             map[string]string
}

// MappedValue is one value found while walking a mapped path.
type MappedValue struct {
	Value   interface{}
	Indexes []int
}

// ReadIdentity identifies the response row a read must select.
type ReadIdentity struct {
	Path, Value  string
	AllowMissing bool
}

// SemanticValue lets a custom attribute value decide semantic equality. The methods are
// exported so a type in another package can satisfy the contract; unexported method names
// carry package identity in Go and cannot cross that boundary.
type SemanticValue interface {
	SemanticEqual(attr.Value) bool
}

// NullValue lets a custom attribute value return its own typed null placeholder.
type NullValue interface {
	NullValue() attr.Value
}

// GetIdentityParts splits a composite id and validates the required components.
func GetIdentityParts(id, separator string, count int, defaults map[int]string, essential []int) ([]string, error) {
	if separator == "" || count < 1 {
		return nil, fmt.Errorf("invalid composite identity definition")
	}
	parts := strings.Split(id, separator)
	if len(parts) > count {
		return nil, fmt.Errorf("composite identity has too many components; expected %d", count)
	}
	for len(parts) < count {
		parts = append(parts, "")
	}
	for index, part := range parts {
		if strings.TrimSpace(part) != "" {
			continue
		}
		required := false
		for _, identityIndex := range essential {
			if index == identityIndex {
				required = true
				break
			}
		}
		if value, declared := defaults[index]; !required && declared && strings.TrimSpace(value) != "" {
			parts[index] = value
			continue
		}
		return nil, fmt.Errorf("composite identity component %d is missing", index+1)
	}
	return parts, nil
}

// SelectReadObject picks the response object whose identity matches the request.
func SelectReadObject(raw interface{}, identities []ReadIdentity) (map[string]interface{}, error) {
	if raw == nil {
		return nil, nil
	}
	if len(identities) == 0 {
		return nil, fmt.Errorf("no readable resource identity was declared")
	}
	rows, isList := raw.([]interface{})
	if !isList {
		if object, ok := raw.(map[string]interface{}); ok {
			rows = []interface{}{object}
		} else {
			return nil, fmt.Errorf("unexpected resource response shape")
		}
	}
	var selected map[string]interface{}
	for _, row := range rows {
		object, ok := row.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("resource response item is not an object")
		}
		matches, readable := true, false
		for _, identity := range identities {
			values := WalkMappedPath(object, MappingPath(identity.Path), nil)
			if len(values) > 1 {
				return nil, fmt.Errorf("resource identity %s has multiple values", identity.Path)
			}
			if len(values) == 0 || values[0].Value == nil || fmt.Sprint(values[0].Value) == "" {
				if !identity.AllowMissing {
					matches = false
				}
				continue
			}
			readable = true
			if identity.Value == "" || fmt.Sprint(values[0].Value) != identity.Value {
				matches = false
			}
		}
		if !readable {
			if isList {
				return nil, fmt.Errorf("resource response item has no readable identity")
			}
			continue
		}
		if matches {
			if selected != nil {
				return nil, fmt.Errorf("resource response contains multiple matching identities")
			}
			selected = object
		}
	}
	return selected, nil
}

// ChangedUpdateMappings keeps only the update bindings whose source value changed.
func ChangedUpdateMappings(bindings []Mapping, desired, prior map[string]interface{}, requiredSources []string) []Mapping {
	selected := make([]Mapping, 0, len(bindings))
	for _, binding := range bindings {
		required := false
		for _, source := range requiredSources {
			if binding.Source == source || strings.HasPrefix(binding.Source, source+".") || strings.HasPrefix(binding.Source, source+"[*]") {
				required = true
				break
			}
		}
		if required || !ValuesEqual(WalkMappedPath(desired, MappingPath(binding.Source), nil), WalkMappedPath(prior, MappingPath(binding.Source), nil)) {
			selected = append(selected, binding)
		}
	}
	return selected
}

// Int converts a mapped response value to a types.Int64.
func Int(m map[string]interface{}, key string) types.Int64 {
	cur := interface{}(m)
	for _, seg := range strings.Split(key, ".") {
		mm, ok := cur.(map[string]interface{})
		if !ok {
			return types.Int64Null()
		}
		cur, ok = mm[seg]
		if !ok {
			return types.Int64Null()
		}
	}
	switch v := cur.(type) {
	case float64:
		return types.Int64Value(int64(v))
	case int:
		return types.Int64Value(int64(v))
	case int64:
		return types.Int64Value(v)
	case json.Number:
		if n, err := v.Int64(); err == nil {
			return types.Int64Value(n)
		}
	case string:
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return types.Int64Value(n)
		}
	}
	return types.Int64Null()
}

// PlannedChanged reports whether a planned value differs from the prior state.
func PlannedChanged(plan, state attr.Value) bool {
	if plan == nil || plan.IsUnknown() {
		return false
	}
	if state == nil || state.IsNull() || state.IsUnknown() || plan.IsNull() {
		return state == nil || !plan.Equal(state)
	}
	switch desired := plan.(type) {
	case types.Object:
		previous, ok := state.(types.Object)
		if !ok {
			return true
		}
		oldItems := previous.Attributes()
		for key, value := range desired.Attributes() {
			if PlannedChanged(value, oldItems[key]) {
				return true
			}
		}
		return false
	case types.List:
		previous, ok := state.(types.List)
		if !ok {
			return true
		}
		items, oldItems := desired.Elements(), previous.Elements()
		if len(items) != len(oldItems) {
			return true
		}
		for index, value := range items {
			if PlannedChanged(value, oldItems[index]) {
				return true
			}
		}
		return false
	case types.Set:
		previous, ok := state.(types.Set)
		if !ok {
			return true
		}
		return !SetMatches(desired.Elements(), previous.Elements(), func(a, b attr.Value) bool { return !PlannedChanged(a, b) })
	case types.Map:
		previous, ok := state.(types.Map)
		if !ok {
			return true
		}
		items, oldItems := desired.Elements(), previous.Elements()
		if len(items) != len(oldItems) {
			return true
		}
		for key, value := range items {
			if PlannedChanged(value, oldItems[key]) {
				return true
			}
		}
		return false
	}
	if semantic, ok := plan.(SemanticValue); ok {
		return !semantic.SemanticEqual(state)
	}
	return !plan.Equal(state)
}

// SetMatches compares two sets element by element with the supplied compatibility rule.
func SetMatches(left, right []attr.Value, compatible func(attr.Value, attr.Value) bool) bool {
	if len(left) != len(right) {
		return false
	}
	owners := make([]int, len(right))
	for index := range owners {
		owners[index] = -1
	}
	var assign func(int, []bool) bool
	assign = func(index int, seen []bool) bool {
		for candidate := range right {
			if seen[candidate] || !compatible(left[index], right[candidate]) {
				continue
			}
			seen[candidate] = true
			if owners[candidate] == -1 || assign(owners[candidate], seen) {
				owners[candidate] = index
				return true
			}
		}
		return false
	}
	for index := range left {
		if !assign(index, make([]bool, len(right))) {
			return false
		}
	}
	return true
}

// ConfiguredMatch reports whether the configured value matches the observed state.
func ConfiguredMatch(config, state attr.Value) bool {
	if config == nil || config.IsNull() {
		return true
	}
	if config.IsUnknown() || state == nil || state.IsNull() || state.IsUnknown() {
		return false
	}
	switch value := config.(type) {
	case types.Object:
		previous, ok := state.(types.Object)
		if !ok {
			return false
		}
		fields := previous.Attributes()
		for key, child := range value.Attributes() {
			if !ConfiguredMatch(child, fields[key]) {
				return false
			}
		}
		return true
	case types.List:
		previous, ok := state.(types.List)
		if !ok || len(value.Elements()) != len(previous.Elements()) {
			return false
		}
		for index, child := range value.Elements() {
			if !ConfiguredMatch(child, previous.Elements()[index]) {
				return false
			}
		}
		return true
	case types.Set:
		previous, ok := state.(types.Set)
		return ok && SetMatches(value.Elements(), previous.Elements(), ConfiguredMatch)
	case types.Map:
		previous, ok := state.(types.Map)
		if !ok || len(value.Elements()) != len(previous.Elements()) {
			return false
		}
		for key, child := range value.Elements() {
			if !ConfiguredMatch(child, previous.Elements()[key]) {
				return false
			}
		}
		return true
	}
	if semantic, ok := config.(SemanticValue); ok {
		return semantic.SemanticEqual(state)
	}
	return config.Equal(state)
}

// SetMemberMatches reports whether one observed set member matches its prior value.
func SetMemberMatches(observed, prior attr.Value, fields map[string]WireField) bool {
	if prior == nil || prior.IsNull() || prior.IsUnknown() {
		return true
	}
	if observed == nil || observed.IsNull() || observed.IsUnknown() {
		return false
	}
	switch current := observed.(type) {
	case types.Object:
		previous, ok := prior.(types.Object)
		if !ok {
			return false
		}
		now, old := current.Attributes(), previous.Attributes()
		for key, field := range fields {
			if !field.Read || !field.Input {
				continue
			}
			if len(field.Fields) > 0 {
				if !SetMemberMatches(now[key], old[key], field.Fields) {
					return false
				}
			} else if !ConfiguredMatch(old[key], now[key]) {
				return false
			}
		}
		return true
	case types.List:
		previous, ok := prior.(types.List)
		if !ok || len(current.Elements()) != len(previous.Elements()) {
			return false
		}
		for index, value := range current.Elements() {
			if !SetMemberMatches(value, previous.Elements()[index], fields) {
				return false
			}
		}
		return true
	case types.Set:
		previous, ok := prior.(types.Set)
		return ok && SetMatches(current.Elements(), previous.Elements(), func(a, b attr.Value) bool { return SetMemberMatches(a, b, fields) })
	}
	return ConfiguredMatch(prior, observed)
}

// PreserveRequestValues restores request-scoped values the API does not echo back.
func PreserveRequestValues(observed, prior attr.Value, fields map[string]WireField) attr.Value {
	if observed == nil || prior == nil || prior.IsNull() || prior.IsUnknown() {
		return observed
	}
	switch current := observed.(type) {
	case types.List:
		previous, ok := prior.(types.List)
		if !ok || current.IsNull() || current.IsUnknown() {
			return observed
		}
		items, oldItems := current.Elements(), previous.Elements()
		for index := range items {
			if index < len(oldItems) {
				items[index] = PreserveRequestValues(items[index], oldItems[index], fields)
			}
		}
		return types.ListValueMust(current.ElementType(context.Background()), items)
	case types.Set:
		previous, ok := prior.(types.Set)
		if !ok || current.IsNull() || current.IsUnknown() {
			return observed
		}
		items, oldItems := current.Elements(), previous.Elements()
		for index, item := range items {
			var retained attr.Value
			ambiguous := false
			for _, old := range oldItems {
				if !SetMemberMatches(item, old, fields) {
					continue
				}
				candidate := PreserveRequestValues(item, old, fields)
				if retained != nil && !retained.Equal(candidate) {
					ambiguous = true
					break
				}
				retained = candidate
			}
			// Never attach a different member's request-only configuration on ambiguity.
			if retained != nil && !ambiguous {
				items[index] = retained
			}
		}
		return types.SetValueMust(current.ElementType(context.Background()), items)
	case types.Object:
		previous, ok := prior.(types.Object)
		if !ok || current.IsNull() || current.IsUnknown() {
			return observed
		}
		items, oldItems := current.Attributes(), previous.Attributes()
		for key, field := range fields {
			old, exists := oldItems[key]
			if !exists {
				continue
			}
			if !field.Read {
				items[key] = KnownValue(old)
			} else if len(field.Fields) > 0 {
				items[key] = PreserveRequestValues(items[key], old, field.Fields)
			}
		}
		return types.ObjectValueMust(current.AttributeTypes(context.Background()), items)
	}
	return observed
}

// KnownValue returns a known value, replacing an unknown or null placeholder.
func KnownValue(value attr.Value) attr.Value {
	if value != nil && value.IsUnknown() {
		if custom, ok := value.(NullValue); ok {
			return custom.NullValue()
		}
	}
	switch v := value.(type) {
	case types.String:
		if v.IsUnknown() {
			return types.StringNull()
		}
	case types.Bool:
		if v.IsUnknown() {
			return types.BoolNull()
		}
	case types.Int64:
		if v.IsUnknown() {
			return types.Int64Null()
		}
	case types.Float64:
		if v.IsUnknown() {
			return types.Float64Null()
		}
	case types.List:
		if v.IsUnknown() {
			return types.ListNull(v.ElementType(context.Background()))
		}
		if !v.IsNull() {
			items := []attr.Value{}
			for _, item := range v.Elements() {
				items = append(items, KnownValue(item))
			}
			return types.ListValueMust(v.ElementType(context.Background()), items)
		}
	case types.Set:
		if v.IsUnknown() {
			return types.SetNull(v.ElementType(context.Background()))
		}
		if !v.IsNull() {
			items := []attr.Value{}
			for _, item := range v.Elements() {
				items = append(items, KnownValue(item))
			}
			return types.SetValueMust(v.ElementType(context.Background()), items)
		}
	case types.Map:
		if v.IsUnknown() {
			return types.MapNull(v.ElementType(context.Background()))
		}
		if !v.IsNull() {
			items := map[string]attr.Value{}
			for key, item := range v.Elements() {
				items[key] = KnownValue(item)
			}
			return types.MapValueMust(v.ElementType(context.Background()), items)
		}
	case types.Object:
		if v.IsUnknown() {
			return types.ObjectNull(v.AttributeTypes(context.Background()))
		}
		if !v.IsNull() {
			items := map[string]attr.Value{}
			for key, item := range v.Attributes() {
				items[key] = KnownValue(item)
			}
			return types.ObjectValueMust(v.AttributeTypes(context.Background()), items)
		}
	}
	return value
}

// ValueToInterface converts a framework attribute value to a request payload value.
func ValueToInterface(value attr.Value, fields map[string]WireField, single bool) interface{} {
	if value == nil || value.IsNull() || value.IsUnknown() {
		return nil
	}
	switch v := value.(type) {
	case interface{ ValueString() string }:
		return v.ValueString()
	case types.Bool:
		return v.ValueBool()
	case types.Int64:
		return v.ValueInt64()
	case types.Float64:
		return v.ValueFloat64()
	case types.Object:
		out := map[string]interface{}{}
		for key, child := range v.Attributes() {
			field, mapped := fields[key]
			if !mapped {
				field.Name = key
			}
			if raw := ValueToInterface(child, field.Fields, field.Single); raw != nil {
				out[field.Name] = raw
			}
		}
		return out
	case types.Map:
		out := map[string]interface{}{}
		for key, child := range v.Elements() {
			if raw := ValueToInterface(child, nil, false); raw != nil {
				out[key] = raw
			}
		}
		return out
	case types.List:
		out := make([]interface{}, 0, len(v.Elements()))
		for _, child := range v.Elements() {
			if raw := ValueToInterface(child, fields, false); raw != nil {
				out = append(out, raw)
			}
		}
		if single {
			if len(out) == 0 {
				return nil
			}
			return out[0]
		}
		return out
	case types.Set:
		out := make([]interface{}, 0, len(v.Elements()))
		for _, child := range v.Elements() {
			if raw := ValueToInterface(child, fields, false); raw != nil {
				out = append(out, raw)
			}
		}
		return out
	}
	return nil
}

// MappingPath splits a JSONPath-like mapping path into segments.
func MappingPath(path string) []string {
	path = strings.TrimPrefix(path, "$.")
	path = strings.TrimPrefix(path, "$")
	path = strings.ReplaceAll(path, "[*]", ".*")
	if path == "" {
		return nil
	}
	return strings.Split(strings.TrimPrefix(path, "."), ".")
}

// WalkMappedPath returns every value reachable at the mapped path.
func WalkMappedPath(value interface{}, path []string, indexes []int) []MappedValue {
	if len(path) == 0 {
		return []MappedValue{{Value: value, Indexes: indexes}}
	}
	if path[0] == "*" {
		var out []MappedValue
		if items, ok := value.([]interface{}); ok {
			for i, item := range items {
				next := append(append([]int{}, indexes...), i)
				out = append(out, WalkMappedPath(item, path[1:], next)...)
			}
		}
		return out
	}
	if m, ok := value.(map[string]interface{}); ok {
		if child, exists := m[path[0]]; exists {
			return WalkMappedPath(child, path[1:], indexes)
		}
	}
	return nil
}

// SetMappedPath writes a raw value back at the mapped path.
func SetMappedPath(value interface{}, path []string, indexes []int, raw interface{}) interface{} {
	if len(path) == 0 {
		return raw
	}
	if path[0] == "*" {
		if len(indexes) == 0 {
			return value
		}
		items, _ := value.([]interface{})
		for len(items) <= indexes[0] {
			items = append(items, nil)
		}
		items[indexes[0]] = SetMappedPath(items[indexes[0]], path[1:], indexes[1:], raw)
		return items
	}
	m, _ := value.(map[string]interface{})
	if m == nil {
		m = map[string]interface{}{}
	}
	m[path[0]] = SetMappedPath(m[path[0]], path[1:], indexes, raw)
	return m
}

// ApplyMappings writes the mapped values into the request map.
func ApplyMappings(request, values map[string]interface{}, bindings []Mapping, roa bool) error {
	mapped := map[string]interface{}{}
	repeated := map[string]bool{}
	for _, binding := range bindings {
		source, target := binding.Source, binding.Target
		if binding.Encoding == "comma" || binding.Encoding == "first" || binding.Encoding == "wrap-json" || binding.Encoding == "wrap-array" {
			source = strings.TrimSuffix(source, "[*]")
			target = strings.TrimSuffix(target, "[*]")
		}
		entries := WalkMappedPath(values, MappingPath(source), nil)
		if binding.Encoding == "scoped-const" {
			parent := source[:strings.LastIndex(source, ".")]
			entries = WalkMappedPath(values, MappingPath(parent), nil)
		}
		if strings.HasPrefix(binding.Encoding, "concat|") {
			parts := strings.Split(binding.Encoding, "|")
			fields := []string{}
			complete := true
			for _, propertyPath := range parts[2:] {
				found := WalkMappedPath(values, MappingPath(propertyPath), nil)
				if len(found) == 0 || found[0].Value == nil {
					complete = false
					break
				}
				fields = append(fields, fmt.Sprint(found[0].Value))
			}
			entries = nil
			if complete {
				entries = []MappedValue{{Value: strings.Join(fields, parts[1])}}
			}
		}
		if strings.Contains(source, "[*]") && !strings.Contains(target, "[*]") && binding.Encoding == "" {
			items := []interface{}{}
			for _, entry := range entries {
				if entry.Value != nil {
					items = append(items, entry.Value)
				}
			}
			if len(entries) > 0 {
				entries = []MappedValue{{Value: items}}
			}
		}
		concreteEntries := entries[:0]
		for _, entry := range entries {
			if entry.Value != nil {
				concreteEntries = append(concreteEntries, entry)
			}
		}
		entries = concreteEntries
		if binding.Constant != "" && !(binding.Encoding == "default" && len(entries) > 0) {
			if (binding.Encoding == "key" || binding.Encoding == "scoped-const") && len(entries) == 0 {
				continue
			}
			var raw interface{}
			if err := json.Unmarshal([]byte(binding.Constant), &raw); err != nil {
				return err
			}
			if len(entries) == 0 {
				entries = []MappedValue{{Value: raw}}
			} else {
				for index := range entries {
					entries[index].Value = raw
				}
			}
		}
		if strings.HasPrefix(binding.Encoding, "uniform-scalar|") && len(entries) > 0 {
			scalarType := strings.TrimPrefix(binding.Encoding, "uniform-scalar|")
			for _, entry := range entries {
				valid := false
				switch scalarType {
				case "boolean":
					_, valid = entry.Value.(bool)
				case "integer":
					switch value := entry.Value.(type) {
					case int, int8, int16, int32, int64, uint, uint8, uint16, uint32:
						valid = true
					case float64:
						valid = value == float64(int64(value))
					case json.Number:
						_, err := value.Int64()
						valid = err == nil
					}
				case "number":
					switch entry.Value.(type) {
					case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, float32, float64, json.Number:
						valid = true
					}
				}
				if !valid {
					return fmt.Errorf("configured value does not match scalar API parameter %s type %s", target, scalarType)
				}
			}
			first, err := json.Marshal(entries[0].Value)
			if err != nil {
				return err
			}
			for _, entry := range entries[1:] {
				next, err := json.Marshal(entry.Value)
				if err != nil {
					return err
				}
				if string(next) != string(first) {
					return fmt.Errorf("conflicting configured values for scalar API parameter %s", target)
				}
			}
			entries = entries[:1]
		}
		for _, entry := range entries {
			raw := entry.Value
			if raw == nil {
				continue
			}
			if mappedValue, ok := binding.Values[fmt.Sprint(raw)]; ok {
				raw = mappedValue
			}
			if strings.HasPrefix(binding.Encoding, "suffix|") {
				raw = fmt.Sprint(raw) + strings.TrimPrefix(binding.Encoding, "suffix|")
			}
			switch binding.Encoding {
			case "first":
				if list, ok := raw.([]interface{}); ok {
					if len(list) == 0 {
						continue
					}
					raw = list[0]
				}
			case "comma":
				if list, ok := raw.([]interface{}); ok {
					parts := make([]string, len(list))
					for i, item := range list {
						parts[i] = fmt.Sprint(item)
					}
					raw = strings.Join(parts, ",")
				}
			case "wrap-array":
				raw = []interface{}{raw}
			case "wrap-json":
				encoded, err := json.Marshal([]interface{}{raw})
				if err != nil {
					return err
				}
				raw = string(encoded)
			}
			indexes := append([]int(nil), entry.Indexes...)
			if len(binding.IndexSources) > 0 {
				if len(binding.IndexSources) != strings.Count(target, "[*]") {
					return fmt.Errorf("request wildcard alignment has an unexpected target dimension count")
				}
				indexes = make([]int, len(binding.IndexSources))
				for position, sourceIndex := range binding.IndexSources {
					if sourceIndex == -1 {
						continue
					}
					if sourceIndex < 0 || sourceIndex >= len(entry.Indexes) {
						return fmt.Errorf("request wildcard alignment is missing a source dimension")
					}
					indexes[position] = entry.Indexes[sourceIndex]
				}
			} else {
				for len(indexes) < strings.Count(target, "[*]") {
					indexes = append(indexes, 0)
				}
			}
			mapped = SetMappedPath(mapped, MappingPath(target), indexes, raw).(map[string]interface{})
		}
		if binding.Repeat {
			if parts := MappingPath(target); len(parts) > 0 {
				repeated[parts[0]] = true
			}
		}
	}
	if err := CompactKvPairs(mapped, bindings); err != nil {
		return err
	}
	for key, raw := range mapped {
		if roa {
			request[key] = raw
			continue
		}
		switch raw.(type) {
		case map[string]interface{}, []interface{}:
			if repeated[key] {
				FlattenRequest(request, key, raw)
			} else {
				encoded, err := json.Marshal(raw)
				if err != nil {
					return err
				}
				request[key] = string(encoded)
			}
		default:
			request[key] = raw
		}
	}
	return nil
}

// CompactKvPairs compacts repeated key/value pairs into the wire shape.
func CompactKvPairs(mapped map[string]interface{}, bindings []Mapping) error {
	groups := map[string][2]string{}
	for _, binding := range bindings {
		if binding.Group == "" {
			continue
		}
		fields := [2]string{binding.GroupKey, binding.GroupValue}
		if fields[0] == "" || fields[1] == "" || fields[0] == fields[1] {
			return fmt.Errorf("Invalid kvPairs member declaration")
		}
		if previous, exists := groups[binding.Group]; exists && previous != fields {
			return fmt.Errorf("Conflicting kvPairs member declaration")
		}
		groups[binding.Group] = fields
	}
	for group, fields := range groups {
		raw, exists := mapped[group]
		if !exists {
			continue
		}
		members, ok := raw.(map[string]interface{})
		if !ok {
			return fmt.Errorf("kvPairs group is not an indexed object")
		}
		indexes := []int{}
		for key := range members {
			index, err := strconv.Atoi(key)
			if err != nil || index < 1 || strconv.Itoa(index) != key {
				return fmt.Errorf("Invalid kvPairs member index")
			}
			indexes = append(indexes, index)
		}
		sort.Ints(indexes)
		compact := map[string]interface{}{}
		codes := map[string]bool{}
		for position, index := range indexes {
			pair, valid := members[strconv.Itoa(index)].(map[string]interface{})
			if !valid || len(pair) != 2 || pair[fields[0]] == nil || pair[fields[1]] == nil {
				return fmt.Errorf("Incomplete kvPairs member")
			}
			code, valid := pair[fields[0]].(string)
			if !valid || code == "" || codes[code] {
				return fmt.Errorf("Invalid or duplicate kvPairs code")
			}
			codes[code] = true
			compact[strconv.Itoa(position+1)] = pair
		}
		mapped[group] = compact
	}
	return nil
}

// FlattenRequest expands one nested request value through its mapping.
func FlattenRequest(request map[string]interface{}, key string, raw interface{}) {
	switch v := raw.(type) {
	case []interface{}:
		for i, item := range v {
			FlattenRequest(request, fmt.Sprintf("%s.%d", key, i+1), item)
		}
	case map[string]interface{}:
		for name, item := range v {
			FlattenRequest(request, key+"."+name, item)
		}
	default:
		if v != nil {
			request[key] = v
		}
	}
}

// NormalizeResponse normalizes the response object for state conversion.
func NormalizeResponse(object map[string]interface{}, bindings []Mapping) map[string]interface{} {
	// Whole-object mappings must precede their leaf aliases, regardless of AM order.
	// Stable insertion order preserves same-depth precedence and source array indexes.
	bindings = append([]Mapping(nil), bindings...)
	for index := 1; index < len(bindings); index++ {
		for current := index; current > 0 && len(MappingPath(bindings[current].Target)) < len(MappingPath(bindings[current-1].Target)); current-- {
			bindings[current], bindings[current-1] = bindings[current-1], bindings[current]
		}
	}
	out := map[string]interface{}{}
	for key, value := range object {
		out[key] = value
	}
	for _, binding := range bindings {
		if strings.HasSuffix(binding.Source, "[*]") && strings.HasSuffix(binding.Target, "[*]") {
			source := strings.TrimSuffix(binding.Source, "[*]")
			target := strings.TrimSuffix(binding.Target, "[*]")
			for _, container := range WalkMappedPath(object, MappingPath(source), nil) {
				if values, ok := container.Value.([]interface{}); ok && len(values) == 0 {
					out = SetMappedPath(out, MappingPath(target), container.Indexes, []interface{}{}).(map[string]interface{})
				}
			}
		}
		for _, entry := range WalkMappedPath(object, MappingPath(binding.Source), nil) {
			raw := entry.Value
			target := binding.Target
			if binding.Encoding == "split-comma" {
				target = strings.TrimSuffix(target, "[*]")
				items := []interface{}{}
				if text, ok := raw.(string); ok && text != "" {
					for _, part := range strings.Split(text, ",") {
						items = append(items, strings.TrimSpace(part))
					}
				}
				raw = items
			}
			if binding.Encoding == "scalar-string" && raw != nil {
				switch raw.(type) {
				case string:
				case int, int32, int64, float32, float64, json.Number:
					raw = fmt.Sprint(raw)
				}
			}
			if mapped, ok := binding.Values[fmt.Sprint(raw)]; ok {
				raw = mapped
			}
			out = SetMappedPath(out, MappingPath(target), entry.Indexes, raw).(map[string]interface{})
		}
	}
	return out
}

// ValuesEqual compares two mapped values with numeric tolerance.
func ValuesEqual(a, b interface{}) bool {
	aa, errA := json.Marshal(a)
	bb, errB := json.Marshal(b)
	return errA == nil && errB == nil && string(aa) == string(bb)
}
