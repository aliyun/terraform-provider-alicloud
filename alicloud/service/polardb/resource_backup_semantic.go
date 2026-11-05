// Package polardb. Custom framework type for preferred_backup_period weekday-set
// semantic equality (F7). Not generated; hand-written to fix PR review F7.
package polardb

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// preferredBackupPeriodStringType is a custom string type for preferred_backup_period.
// PolarDB normalizes the configured backup period to a canonical weekday order
// (e.g. "Wednesday,Monday" -> "Monday,Wednesday"). A plain types.String compares the
// exact string, so the reordered value causes a persistent plan diff or a
// "Provider produced inconsistent result after apply" error.
//
// F7: with semantic equality, any weekday ordering is treated as equal so apply/refresh
// converge, while a genuine add/remove of weekdays still triggers an update.
type preferredBackupPeriodStringType struct {
	basetypes.StringType
}

var _ basetypes.StringTypable = preferredBackupPeriodStringType{}

// Equal reports type equivalence. The embedded StringType.Equal only matches
// basetypes.StringType, so it must be overridden to recognize this custom type itself.
func (t preferredBackupPeriodStringType) Equal(o attr.Type) bool {
	_, ok := o.(preferredBackupPeriodStringType)
	return ok
}

func (t preferredBackupPeriodStringType) String() string {
	return "preferredBackupPeriodString"
}

// ValueFromTerraform returns the custom value so plan/state values carry semantic equality.
func (t preferredBackupPeriodStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	v, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}
	sv, ok := v.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type %T for preferredBackupPeriodStringType", v)
	}
	return preferredBackupPeriodStringValue{sv}, nil
}

// ValueFromString returns the custom value from a StringValue.
func (t preferredBackupPeriodStringType) ValueFromString(_ context.Context, v basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return preferredBackupPeriodStringValue{v}, nil
}

// ValueType returns the custom value type for diagnostics/type detection.
func (t preferredBackupPeriodStringType) ValueType(_ context.Context) attr.Value {
	return preferredBackupPeriodStringValue{}
}

// preferredBackupPeriodStringValue is the custom value type. It embeds basetypes.StringValue
// (inheriting all StringValuable behavior) and adds weekday-set semantic equality.
type preferredBackupPeriodStringValue struct {
	basetypes.StringValue
}

var _ basetypes.StringValuableWithSemanticEquals = preferredBackupPeriodStringValue{}

// Type returns the custom type, overriding the embedded StringValue.Type which would return
// basetypes.StringType. State/Plan object marshaling (NewObjectValueMust / ObjectType.ValueFromTerraform)
// compares each attribute value's Type(ctx) against the schema attribute type; without this override
// the value reports basetypes.StringType while the schema expects preferredBackupPeriodStringType,
// causing "Invalid Object Attribute Type" panics during State.Set.
func (v preferredBackupPeriodStringValue) Type(_ context.Context) attr.Type {
	return preferredBackupPeriodStringType{}
}

// StringSemanticEquals returns true when the two weekday lists are the same set, regardless of
// order. Only known values reach this path (null/unknown are handled by the framework).
func (v preferredBackupPeriodStringValue) StringSemanticEquals(ctx context.Context, other basetypes.StringValuable) (bool, diag.Diagnostics) {
	if other == nil {
		return false, nil
	}
	ov, diags := other.ToStringValue(ctx)
	if diags.HasError() {
		return false, diags
	}
	return backupWeekdaySetEqual(v.ValueString(), ov.ValueString()), nil
}

// backupWeekdaySet splits a comma-separated weekday list into a lower-cased set, trimming
// whitespace and dropping empty entries.
func backupWeekdaySet(s string) map[string]struct{} {
	set := map[string]struct{}{}
	for _, part := range strings.Split(s, ",") {
		p := strings.TrimSpace(part)
		if p == "" {
			continue
		}
		set[strings.ToLower(p)] = struct{}{}
	}
	return set
}

// backupWeekdaySetEqual reports whether two comma-separated weekday lists denote the same set.
func backupWeekdaySetEqual(a, b string) bool {
	as, bs := backupWeekdaySet(a), backupWeekdaySet(b)
	if len(as) != len(bs) {
		return false
	}
	for k := range as {
		if _, ok := bs[k]; !ok {
			return false
		}
	}
	return true
}

// backupImportedAdvancedPolicies builds advanced_data_policies state for an imported backup,
// which has no configuration to align with. The API reports the configured policies next to the
// ones it auto-created when advanced backup was enabled; only the configured ones belong in the
// imported state. It also keeps one entry per structural chain: after an in-place policy change
// (ModifyBackupPolicy with action_type CREATE) the service reports the earlier entry for the same
// chain as well, while the configuration block describes that chain once. The first entry the
// service returns for a chain is the current one. A list without any configured policy is
// reported as null, matching a resource that was applied without the block.
func backupImportedAdvancedPolicies(observed types.List) types.List {
	elementType := observed.ElementType(context.Background())
	if observed.IsNull() || observed.IsUnknown() {
		return types.ListNull(elementType)
	}
	kept := make([]attr.Value, 0, len(observed.Elements()))
	seenChains := map[string]bool{}
	for _, entry := range observed.Elements() {
		object, ok := entry.(types.Object)
		if !ok || object.IsNull() || object.IsUnknown() {
			continue
		}
		if autoCreated, exists := object.Attributes()["auto_created"]; exists && !autoCreated.IsNull() && !autoCreated.IsUnknown() {
			if flag, isBool := autoCreated.(types.Bool); isBool && flag.ValueBool() {
				continue
			}
		}
		if chain := backupAdvancedPolicyChainKey(object); chain != "" {
			if seenChains[chain] {
				continue
			}
			seenChains[chain] = true
		}
		kept = append(kept, object)
	}
	if len(kept) == 0 {
		return types.ListNull(elementType)
	}
	return types.ListValueMust(elementType, kept)
}

// backupAdvancedPolicyChainKey renders the structural chain of a policy (src/dest type and
// region) as a deduplication key. An entry without the full chain returns "" and is never
// deduplicated.
func backupAdvancedPolicyChainKey(policy types.Object) string {
	attributes := policy.Attributes()
	values := make([]string, 0, len(backupAdvancedPolicyChainKeys))
	for _, key := range backupAdvancedPolicyChainKeys {
		value, exists := attributes[key]
		if !exists || value.IsNull() || value.IsUnknown() {
			return ""
		}
		text, isString := value.(types.String)
		if !isString || text.ValueString() == "" {
			return ""
		}
		values = append(values, text.ValueString())
	}
	return strings.Join(values, "|")
}

// backupAlignAdvancedPolicies rebuilds advanced_data_policies state from the configured
// entries, filling each entry's computed attributes (policy_id, auto_created, ...) from the
// matching DescribeBackupPolicy entry.
//
// The nested block is configuration-owned and cannot be Computed, so the read result cannot
// be used as-is: the API also returns the policy it auto-creates when advanced backup is
// enabled, and using that list would change the block count between plan and apply. Using the
// plan values alone is not enough either, because omitted Optional+Computed sub-attributes are
// unknown at apply time and an unknown value must not be written back to state. Matching on the
// configured (read-back) fields keeps the count and every value known.
func backupAlignAdvancedPolicies(configured, observed types.List) types.List {
	if configured.IsNull() || configured.IsUnknown() {
		elementType := configured.ElementType(context.Background())
		if elementType == nil {
			return observed
		}
		return types.ListNull(elementType)
	}
	var candidates []attr.Value
	if !observed.IsNull() && !observed.IsUnknown() {
		candidates = observed.Elements()
	}
	entries := configured.Elements()
	aligned := make([]attr.Value, 0, len(entries))
	for _, entry := range entries {
		configuredObject, ok := entry.(types.Object)
		if !ok || configuredObject.IsNull() || configuredObject.IsUnknown() {
			aligned = append(aligned, fwKnownValue(entry))
			continue
		}
		match, found := backupMatchAdvancedPolicy(configuredObject, candidates)
		if !found {
			aligned = append(aligned, fwKnownValue(entry))
			continue
		}
		observedAttributes := match.Attributes()
		attributes := configuredObject.Attributes()
		for key, value := range attributes {
			if !value.IsUnknown() {
				continue
			}
			if observedValue, exists := observedAttributes[key]; exists && !observedValue.IsNull() && !observedValue.IsUnknown() {
				attributes[key] = observedValue
			} else {
				attributes[key] = fwKnownValue(value)
			}
		}
		aligned = append(aligned, types.ObjectValueMust(configuredObject.AttributeTypes(context.Background()), attributes))
	}
	return types.ListValueMust(configured.ElementType(context.Background()), aligned)
}

// backupMatchAdvancedPolicy finds the observed policy describing the same configured policy.
// action_type is an operation (CREATE/UPDATE/DELETE) and is not read back, so it is skipped;
// every other configured, known attribute must match exactly. A policy whose mutable leaves
// (retention, preserve flags) were just changed no longer matches exactly, so the structural
// chain is used as a fallback.
func backupMatchAdvancedPolicy(configured types.Object, candidates []attr.Value) (types.Object, bool) {
	if match, ok := backupMatchAdvancedPolicyExact(configured, candidates); ok {
		return match, true
	}
	return backupMatchAdvancedPolicyByChain(configured, candidates)
}

func backupMatchAdvancedPolicyExact(configured types.Object, candidates []attr.Value) (types.Object, bool) {
	attributes := configured.Attributes()
	for _, candidate := range candidates {
		observedObject, ok := candidate.(types.Object)
		if !ok || observedObject.IsNull() || observedObject.IsUnknown() {
			continue
		}
		observedAttributes := observedObject.Attributes()
		matched := true
		for key, value := range attributes {
			if key == "action_type" || value.IsNull() || value.IsUnknown() {
				continue
			}
			observedValue, exists := observedAttributes[key]
			if !exists || observedValue.IsUnknown() || observedValue.IsNull() || !observedValue.Equal(value) {
				matched = false
				break
			}
		}
		if matched {
			return observedObject, true
		}
	}
	return types.Object{}, false
}

// backupAdvancedPolicyChainKeys are the structural attributes describing what is backed up and
// where. Retention, preserve flags and the schedule value are mutable knobs: DBS resolves an
// ActionType=UPDATE entry by PolicyId, so an update that changes only those leaves still has to
// carry the id of the policy it describes.
var backupAdvancedPolicyChainKeys = []string{
	"src_type", "src_region", "dest_type", "dest_region", "bak_type", "filter_type", "filter_key",
}

// backupMatchAdvancedPolicyByChain matches on the structural chain only, and only when exactly
// one observed policy fits: attaching an ambiguous or unevidenced id would update a policy the
// configuration does not describe.
func backupMatchAdvancedPolicyByChain(configured types.Object, candidates []attr.Value) (types.Object, bool) {
	attributes := configured.Attributes()
	compared := 0
	for _, key := range backupAdvancedPolicyChainKeys {
		if value, exists := attributes[key]; exists && !value.IsNull() && !value.IsUnknown() {
			compared++
		}
	}
	if compared < 2 {
		return types.Object{}, false
	}
	var matches []types.Object
	for _, candidate := range candidates {
		observedObject, ok := candidate.(types.Object)
		if !ok || observedObject.IsNull() || observedObject.IsUnknown() {
			continue
		}
		observedAttributes := observedObject.Attributes()
		matched := true
		for _, key := range backupAdvancedPolicyChainKeys {
			value, exists := attributes[key]
			if !exists || value.IsNull() || value.IsUnknown() {
				continue
			}
			observedValue, exists := observedAttributes[key]
			// A field the read-back does not report carries no counter-evidence; only a
			// reported, different value rules the candidate out.
			if !exists || observedValue.IsUnknown() || observedValue.IsNull() {
				continue
			}
			if !observedValue.Equal(value) {
				matched = false
				break
			}
		}
		if matched {
			matches = append(matches, observedObject)
		}
	}
	if len(matches) != 1 {
		return types.Object{}, false
	}
	return matches[0], true
}
