// Package cr. This file is generated automatically. Please do not modify it manually, thank you!
// Service helpers for the framework-native alicloud_cr_repository resource.
// The resource embeds fwadapt.ResourceBase and reads the provider client through Meta().
// Cross-cutting helpers are taken from the shared leaf packages (alicloud/helper for client
// token and retry, alicloud/errs for error classification); this file only keeps the
// CR-specific response checks, identity mapping and state-wait logic.
package cr

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type fwTimeoutDurationValidator struct{}

type fwWireField struct {
	Name   string
	Single bool
	Read   bool
	Input  bool
	Fields map[string]fwWireField
}

type fwMapping struct {
	Source, Target, Encoding, Constant string
	Group, GroupKey, GroupValue        string
	Repeat                             bool
	IndexSources                       []int
	Values                             map[string]string
}

type fwMappedValue struct {
	Value   interface{}
	Indexes []int
}
type fwDeleteCondition struct {
	Property, Action string
	Values           []string
}
type fwReadIdentity struct {
	Path, Value  string
	AllowMissing bool
}

func fwTimeoutTypes() map[string]attr.Type {
	return map[string]attr.Type{"create": types.StringType, "update": types.StringType, "delete": types.StringType}
}

func fwTimeoutBlock() schema.SingleNestedBlock {
	attributes := map[string]schema.Attribute{}
	for _, operation := range []string{"create", "update", "delete"} {
		attributes[operation] = schema.StringAttribute{
			Optional:    true,
			Description: "Maximum retry and wait duration for this operation, for example 10m.",
			Validators:  []validator.String{fwTimeoutDurationValidator{}},
		}
	}
	return schema.SingleNestedBlock{Attributes: attributes}
}

func (fwTimeoutDurationValidator) Description(context.Context) string {
	return "A positive Go duration, for example 10m or 1h."
}
func (v fwTimeoutDurationValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}
func (fwTimeoutDurationValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	duration, err := time.ParseDuration(req.ConfigValue.ValueString())
	if err != nil || duration <= 0 {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid operation timeout", "Use a positive Go duration, for example 10m or 1h.")
	}
}

func fwTimeoutDuration(value types.Object, operation string, fallback time.Duration) (time.Duration, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	if value.IsNull() {
		return fallback, diagnostics
	}
	if value.IsUnknown() {
		diagnostics.AddError("Unknown operation timeout", "Resolve the timeouts block before applying.")
		return 0, diagnostics
	}
	raw, ok := value.Attributes()[operation]
	if !ok || raw.IsNull() {
		return fallback, diagnostics
	}
	configured, ok := raw.(types.String)
	if !ok || configured.IsUnknown() {
		diagnostics.AddError("Invalid operation timeout", "Resolve the operation timeout to a duration string before applying.")
		return 0, diagnostics
	}
	duration, err := time.ParseDuration(configured.ValueString())
	if err != nil || duration <= 0 {
		diagnostics.AddError("Invalid operation timeout", "Use a positive Go duration, for example 10m or 1h.")
		return 0, diagnostics
	}
	return duration, diagnostics
}

// retryRpcPost runs an RPC write through the shared retry wrapper, adding the CR-specific
// transient codes on top of helper.NeedRetry. ctx is the lifecycle context.
func retryRpcPost(ctx context.Context, client *connectivity.AliyunClient, product, version, action string, request map[string]interface{}, timeout time.Duration) (map[string]interface{}, error) {
	return helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
		return client.RpcPost(product, version, action, nil, request, true)
	}, timeout, repositoryRetryable)
}

// repositoryRetryable is the narrow CR wrapper over helper.NeedRetry: v1 retried the
// conflict/status/busy codes and DependencyViolation per call, which NeedRetry does not cover.
func repositoryRetryable(err error) bool {
	if helper.NeedRetry(err) {
		return true
	}
	sdkError, ok := err.(*tea.SDKError)
	if !ok || sdkError.Code == nil {
		return false
	}
	code := *sdkError.Code
	return strings.HasPrefix(code, "DependencyViolation") || code == "OperationConflict" ||
		code == "LastTokenProcessing" || code == "IncorrectStatus" || code == "SystemBusy"
}

// repositoryCheckResponse rejects HTTP-200 failure payloads. The CR APIs report success in
// the body (IsSuccess/Code), so an unchecked failed call could look like an empty object
// on read or a silent no-op on update/delete. A response without IsSuccess is left to the
// call site so older shapes are not misclassified.
func repositoryCheckResponse(action string, response map[string]interface{}) error {
	if response == nil {
		return fmt.Errorf("%s returned an empty response", action)
	}
	success, ok := response["IsSuccess"].(bool)
	if !ok {
		return fmt.Errorf("%s response has no usable IsSuccess flag", action)
	}
	if success {
		return nil
	}
	code, _ := response["Code"].(string)
	if code == "" {
		code = "unknown"
	}
	return fmt.Errorf("%s failed: IsSuccess=false, Code=%s", action, code)
}

// repositoryResponseNotFound reports whether an HTTP-200 Get payload explicitly says the
// repository is gone. Only an explicit not-found code may be treated as absence; a released
// parent instance means the child repository cannot exist either.
func repositoryResponseNotFound(response map[string]interface{}) bool {
	if response == nil {
		return false
	}
	if success, ok := response["IsSuccess"].(bool); ok && success {
		return false
	}
	code, _ := response["Code"].(string)
	return code == "REPO_NOT_EXIST" || code == "INSTANCE_NOT_EXIST" ||
		code == "ResourceNotFound.Repository" || code == "ResourceNotFound"
}

// repositoryNotExist reports the exact NotFound codes for a repository. CR surfaces
// REPO_NOT_EXIST when Get/Delete find nothing and INSTANCE_NOT_EXIST once the parent
// registry instance is released, matching the SDKv2 CrService behaviour.
func repositoryNotExist(err error) bool {
	return errs.IsExpectedErrors(err, []string{"ResourceNotFound.Repository", "ResourceNotFound", "REPO_NOT_EXIST", "INSTANCE_NOT_EXIST"})
}

// repositorySelectReadResponse selects the repository object from a successful Get response.
// IsSuccess already said the call succeeded, so a payload without the repository identity is
// a contract violation, not evidence that the resource is gone.
func repositorySelectReadResponse(response map[string]interface{}, identities []fwReadIdentity, id string) (map[string]interface{}, error) {
	selected, err := fwSelectReadObject(response, identities)
	if err != nil {
		return nil, fmt.Errorf("selecting GetRepository resource: %w", err)
	}
	if selected == nil {
		return nil, fmt.Errorf("GetRepository succeeded but returned no repository identity for %q", id)
	}
	return selected, nil
}

// repositoryGet describes the repository by its "<instance_id>:<repository_id>" id.
// It returns (nil, nil) only when CR explicitly reports the repository or its parent
// instance as gone; every other failure is returned to the caller.
func repositoryGet(ctx context.Context, client *connectivity.AliyunClient, id string, includeSecondary ...bool) (map[string]interface{}, error) {
	completeRead := len(includeSecondary) == 0 || includeSecondary[0]
	getValues := map[string]interface{}{"RegionId": client.RegionId, "InstanceId": id}
	getValues["Status"] = "ALL"
	parts, identityPartsErr := fwGetIdentityParts(id, ":", 2, map[int]string{}, []int{1})
	if identityPartsErr != nil {
		return nil, fmt.Errorf("reading alicloud_cr_repository: %w", identityPartsErr)
	}
	if len(parts) > 0 {
		getValues["InstanceId"] = parts[0]
	}
	if len(parts) > 1 {
		getValues["RepositoryId"] = parts[1]
	}
	request := map[string]interface{}{}
	request["RegionId"] = client.RegionId
	if err := fwApplyMappings(request, getValues, []fwMapping{
		{Source: "$.InstanceId", Target: "InstanceId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RepoNamespaceName", Target: "RepoNamespaceName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RepositoryId", Target: "RepoId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RepositoryName", Target: "RepoName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	}, false); err != nil {
		return nil, fmt.Errorf("building GetRepository request: %w", err)
	}
	identities := []fwReadIdentity{}
	if len(parts) > 0 {
		identities = append(identities, fwReadIdentity{Path: "$.InstanceId", Value: parts[0], AllowMissing: true})
	}
	if len(parts) > 1 {
		identities = append(identities, fwReadIdentity{Path: "$.RepoId", Value: parts[1], AllowMissing: false})
	}
	response, err := retryRpcPost(ctx, client, "cr", "2018-12-01", "GetRepository", request, 1*time.Minute)
	if err != nil {
		if repositoryNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	// Only an explicit not-found payload counts as absence; any other 200 failure must
	// surface as an error instead of decaying into RemoveResource via a missing RepoId.
	if repositoryResponseNotFound(response) {
		return nil, nil
	}
	if err := repositoryCheckResponse("GetRepository", response); err != nil {
		return nil, err
	}
	object, err := repositorySelectReadResponse(response, identities, id)
	if err != nil {
		return nil, err
	}

	return repositoryCompleteRead(client, object, id, completeRead)
}

func repositoryCompleteRead(client *connectivity.AliyunClient, object map[string]interface{}, id string, includeSecondary bool) (map[string]interface{}, error) {
	object = fwNormalizeResponse(object, []fwMapping{
		// gap: GetRepository returns CreateTime/ModifiedTime as int64 UNIX-ms timestamps (json.Number),
		// but the schema declares create_time/modified_time as StringAttribute. The generated str()
		// reader only type-asserts string, so an integer value decodes to Null (failing NotNull checks).
		// Use the generator's scalar-string encoding to coerce numerics to their decimal string form.
		{Source: "$.CreateTime", Target: "$.CreateTime", Encoding: "scalar-string", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Detail", Target: "$.Detail", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.InstanceId", Target: "$.InstanceId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ModifiedTime", Target: "$.ModifiedTime", Encoding: "scalar-string", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RepoBuildType", Target: "$.RepoBuildType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RepoNamespaceName", Target: "$.RepoNamespaceName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RepoType", Target: "$.RepoType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RepoId", Target: "$.RepositoryId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RepoName", Target: "$.RepositoryName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ResourceGroupId", Target: "$.ResourceGroupId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RepoStatus", Target: "$.Status", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Summary", Target: "$.Summary", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.TagImmutability", Target: "$.TagImmutability", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.InstanceId", Target: "$.InstanceId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RepositoryId", Target: "$.RepositoryId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	})
	identityParts := strings.Split(id, ":")
	if len(identityParts) > 0 {
		object["InstanceId"] = identityParts[0]
	}
	if len(identityParts) > 1 {
		object["RepositoryId"] = identityParts[1]
	}
	if !includeSecondary {
		return object, nil
	}
	return object, nil
}

// repositoryStateRefreshFunc polls the repository status field for the gone-wait.
// A missing repository reports state "" so the empty-target wait converges.
func repositoryStateRefreshFunc(ctx context.Context, client *connectivity.AliyunClient, id, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := repositoryGet(ctx, client, id, false)
		if err != nil {
			return nil, "", err
		}
		if object == nil {
			return nil, "", nil
		}
		current := ""
		if values := fwWalkMappedPath(object, fwMappingPath(field), nil); len(values) > 0 {
			current = fmt.Sprint(values[0].Value)
		}
		for _, fs := range failStates {
			if current == fs {
				return object, current, fmt.Errorf("status reached fail state %q", current)
			}
		}
		return object, current, nil
	}
}

// repositoryWaitForState waits for the status field to reach target, mirroring v1 BuildStateConf.
// Returns the last polled object (nil once the resource disappears).
func repositoryWaitForState(ctx context.Context, client *connectivity.AliyunClient, id, field string, target, failStates []string, timeout time.Duration) (map[string]interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{},
		Target:     target,
		Refresh:    repositoryStateRefreshFunc(ctx, client, id, field, failStates),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	raw, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}
	if m, ok := raw.(map[string]interface{}); ok {
		return m, nil
	}
	return nil, nil
}

// Missing request-scope components may use only explicitly declared API defaults.
func fwGetIdentityParts(id, separator string, count int, defaults map[int]string, essential []int) ([]string, error) {
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

// Select by real response identity before CompleteRead can restore request-scoped fields.
func fwSelectReadObject(raw interface{}, identities []fwReadIdentity) (map[string]interface{}, error) {
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
			values := fwWalkMappedPath(object, fwMappingPath(identity.Path), nil)
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

// Optional update parameters only participate when their mapped value changed.
func fwChangedUpdateMappings(bindings []fwMapping, desired, prior map[string]interface{}, requiredSources []string) []fwMapping {
	selected := make([]fwMapping, 0, len(bindings))
	for _, binding := range bindings {
		required := false
		for _, source := range requiredSources {
			if binding.Source == source || strings.HasPrefix(binding.Source, source+".") || strings.HasPrefix(binding.Source, source+"[*]") {
				required = true
				break
			}
		}
		if required || !fwValuesEqual(fwWalkMappedPath(desired, fwMappingPath(binding.Source), nil), fwWalkMappedPath(prior, fwMappingPath(binding.Source), nil)) {
			selected = append(selected, binding)
		}
	}
	return selected
}

// Ignore cloud-managed tags in old state without sending them to UntagResources.
// fwPlannedChanged ignores computed unknown leaves while comparing configured nested values.
func fwPlannedChanged(plan, state attr.Value) bool {
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
			if fwPlannedChanged(value, oldItems[key]) {
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
			if fwPlannedChanged(value, oldItems[index]) {
				return true
			}
		}
		return false
	case types.Set:
		previous, ok := state.(types.Set)
		if !ok {
			return true
		}
		return !fwSetMatches(desired.Elements(), previous.Elements(), func(a, b attr.Value) bool { return !fwPlannedChanged(a, b) })
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
			if fwPlannedChanged(value, oldItems[key]) {
				return true
			}
		}
		return false
	}
	if semantic, ok := plan.(interface{ fwSemanticEqual(attr.Value) bool }); ok {
		return !semantic.fwSemanticEqual(state)
	}
	return !plan.Equal(state)
}

// fwSetMatches pairs unordered members without depending on response or plan order.
func fwSetMatches(left, right []attr.Value, compatible func(attr.Value, attr.Value) bool) bool {
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

func fwDeleteEligible(values map[string]interface{}, conditions []fwDeleteCondition) (bool, error) {
	return fwOperationEligible(values, conditions)
}

func fwOperationEligible(values map[string]interface{}, conditions []fwDeleteCondition) (bool, error) {
	eligible := true
	for _, condition := range conditions {
		found := fwWalkMappedPath(values, fwMappingPath(condition.Property), nil)
		if len(found) == 0 || found[0].Value == nil {
			return false, fmt.Errorf("cannot select an operation: %s is unknown", condition.Property)
		}
		matches := false
		for _, expected := range condition.Values {
			if fmt.Sprint(found[0].Value) == expected {
				matches = true
				break
			}
		}
		switch condition.Action {
		case "valid", "required":
			eligible = eligible && matches
		case "invalid":
			eligible = eligible && !matches
		default:
			return false, fmt.Errorf("unsupported operation dependency action: %s", condition.Action)
		}
	}
	return eligible, nil
}

func fwKnownValue(value attr.Value) attr.Value {
	if value != nil && value.IsUnknown() {
		if custom, ok := value.(interface{ fwNullValue() attr.Value }); ok {
			return custom.fwNullValue()
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
				items = append(items, fwKnownValue(item))
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
				items = append(items, fwKnownValue(item))
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
				items[key] = fwKnownValue(item)
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
				items[key] = fwKnownValue(item)
			}
			return types.ObjectValueMust(v.AttributeTypes(context.Background()), items)
		}
	}
	return value
}

func fwValueToInterface(value attr.Value, fields map[string]fwWireField, single bool) interface{} {
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
			if raw := fwValueToInterface(child, field.Fields, field.Single); raw != nil {
				out[field.Name] = raw
			}
		}
		return out
	case types.Map:
		out := map[string]interface{}{}
		for key, child := range v.Elements() {
			if raw := fwValueToInterface(child, nil, false); raw != nil {
				out[key] = raw
			}
		}
		return out
	case types.List:
		out := make([]interface{}, 0, len(v.Elements()))
		for _, child := range v.Elements() {
			if raw := fwValueToInterface(child, fields, false); raw != nil {
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
			if raw := fwValueToInterface(child, fields, false); raw != nil {
				out = append(out, raw)
			}
		}
		return out
	}
	return nil
}

func fwMappingPath(path string) []string {
	path = strings.TrimPrefix(path, "$.")
	path = strings.TrimPrefix(path, "$")
	path = strings.ReplaceAll(path, "[*]", ".*")
	if path == "" {
		return nil
	}
	return strings.Split(strings.TrimPrefix(path, "."), ".")
}

func fwWalkMappedPath(value interface{}, path []string, indexes []int) []fwMappedValue {
	if len(path) == 0 {
		return []fwMappedValue{{Value: value, Indexes: indexes}}
	}
	if path[0] == "*" {
		var out []fwMappedValue
		if items, ok := value.([]interface{}); ok {
			for i, item := range items {
				next := append(append([]int{}, indexes...), i)
				out = append(out, fwWalkMappedPath(item, path[1:], next)...)
			}
		}
		return out
	}
	if m, ok := value.(map[string]interface{}); ok {
		if child, exists := m[path[0]]; exists {
			return fwWalkMappedPath(child, path[1:], indexes)
		}
	}
	return nil
}

func fwSetMappedPath(value interface{}, path []string, indexes []int, raw interface{}) interface{} {
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
		items[indexes[0]] = fwSetMappedPath(items[indexes[0]], path[1:], indexes[1:], raw)
		return items
	}
	m, _ := value.(map[string]interface{})
	if m == nil {
		m = map[string]interface{}{}
	}
	m[path[0]] = fwSetMappedPath(m[path[0]], path[1:], indexes, raw)
	return m
}

func fwApplyMappings(request, values map[string]interface{}, bindings []fwMapping, roa bool) error {
	mapped := map[string]interface{}{}
	repeated := map[string]bool{}
	for _, binding := range bindings {
		source, target := binding.Source, binding.Target
		if binding.Encoding == "comma" || binding.Encoding == "first" || binding.Encoding == "wrap-json" || binding.Encoding == "wrap-array" {
			source = strings.TrimSuffix(source, "[*]")
			target = strings.TrimSuffix(target, "[*]")
		}
		entries := fwWalkMappedPath(values, fwMappingPath(source), nil)
		if binding.Encoding == "scoped-const" {
			parent := source[:strings.LastIndex(source, ".")]
			entries = fwWalkMappedPath(values, fwMappingPath(parent), nil)
		}
		if strings.HasPrefix(binding.Encoding, "concat|") {
			parts := strings.Split(binding.Encoding, "|")
			fields := []string{}
			complete := true
			for _, propertyPath := range parts[2:] {
				found := fwWalkMappedPath(values, fwMappingPath(propertyPath), nil)
				if len(found) == 0 || found[0].Value == nil {
					complete = false
					break
				}
				fields = append(fields, fmt.Sprint(found[0].Value))
			}
			entries = nil
			if complete {
				entries = []fwMappedValue{{Value: strings.Join(fields, parts[1])}}
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
				entries = []fwMappedValue{{Value: items}}
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
				entries = []fwMappedValue{{Value: raw}}
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
			mapped = fwSetMappedPath(mapped, fwMappingPath(target), indexes, raw).(map[string]interface{})
		}
		if binding.Repeat {
			if parts := fwMappingPath(target); len(parts) > 0 {
				repeated[parts[0]] = true
			}
		}
	}
	if err := fwCompactKvPairs(mapped, bindings); err != nil {
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
				fwFlattenRequest(request, key, raw)
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

// Only explicitly declared kvPairs groups are compacted; ordinary indexed lists retain shape.
func fwCompactKvPairs(mapped map[string]interface{}, bindings []fwMapping) error {
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

func fwFlattenRequest(request map[string]interface{}, key string, raw interface{}) {
	switch v := raw.(type) {
	case []interface{}:
		for i, item := range v {
			fwFlattenRequest(request, fmt.Sprintf("%s.%d", key, i+1), item)
		}
	case map[string]interface{}:
		for name, item := range v {
			fwFlattenRequest(request, key+"."+name, item)
		}
	default:
		if v != nil {
			request[key] = v
		}
	}
}

func fwNormalizeResponse(object map[string]interface{}, bindings []fwMapping) map[string]interface{} {
	// Whole-object mappings must precede their leaf aliases, regardless of AM order.
	// Stable insertion order preserves same-depth precedence and source array indexes.
	bindings = append([]fwMapping(nil), bindings...)
	for index := 1; index < len(bindings); index++ {
		for current := index; current > 0 && len(fwMappingPath(bindings[current].Target)) < len(fwMappingPath(bindings[current-1].Target)); current-- {
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
			for _, container := range fwWalkMappedPath(object, fwMappingPath(source), nil) {
				if values, ok := container.Value.([]interface{}); ok && len(values) == 0 {
					out = fwSetMappedPath(out, fwMappingPath(target), container.Indexes, []interface{}{}).(map[string]interface{})
				}
			}
		}
		for _, entry := range fwWalkMappedPath(object, fwMappingPath(binding.Source), nil) {
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
			out = fwSetMappedPath(out, fwMappingPath(target), entry.Indexes, raw).(map[string]interface{})
		}
	}
	return out
}

func fwValuesEqual(a, b interface{}) bool {
	aa, errA := json.Marshal(a)
	bb, errB := json.Marshal(b)
	return errA == nil && errB == nil && string(aa) == string(bb)
}
