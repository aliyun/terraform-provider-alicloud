// Package ecs contains the Framework-native resources of the ECS product.
package ecs

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdkretry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const storageSetTypeName = "alicloud_ecs_storage_set"

var (
	_ resource.Resource                = &storageSetResource{}
	_ resource.ResourceWithImportState = &storageSetResource{}
)

func NewStorageSetResource() resource.Resource {
	return &storageSetResource{}
}

type storageSetResource struct {
	fwadapt.ResourceBase
}

type storageSetResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	Description        types.String `tfsdk:"description"`
	MaxPartitionNumber types.Int64  `tfsdk:"max_partition_number"`
	RegionId           types.String `tfsdk:"region_id"`
	ResourceGroupId    types.String `tfsdk:"resource_group_id"`
	StorageSetId       types.String `tfsdk:"storage_set_id"`
	StorageSetName     types.String `tfsdk:"storage_set_name"`
	Tags               types.Map    `tfsdk:"tags"`
	ZoneId             types.String `tfsdk:"zone_id"`
}

func (r *storageSetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = storageSetTypeName
}

func (r *storageSetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a ECS Storage Set resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource, same as `storage_set_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the Storage Set.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"max_partition_number": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The maximum number of partitions in the Storage Set. Changing this will create a new resource.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"region_id": schema.StringAttribute{
				Computed:    true,
				Description: "The region ID of the Storage Set.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_group_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The ID of the resource group to which the Storage Set belongs. ECS Storage Set is not a supported ResourceType of JoinResourceGroup (it returns InvalidResourceType.NotFound for storageset), so the resource group can only be set at creation via CreateStorageSet; changing it will create a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"storage_set_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the Storage Set.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"storage_set_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the Storage Set.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tags": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "The tags of the Storage Set. ECS Storage Set is not a supported ResourceType of TagResources/UntagResources (they return InvalidResourceType.NotFound for storageset), so tags can only be set at creation via CreateStorageSet; changing them will create a new resource.",
				// UseStateForUnknown keeps a computed tags value when the user
				// omits tags from config: the plan carries unknown, which is
				// replaced by the prior state so an omitted tags block is a
				// no-op (the cloud tags are preserved). RequiresReplace then
				// compares planned == prior and fires no replace. It does NOT
				// swallow an explicit known value: a non-empty map, or tags = {}
				// (a known empty map → recreate to clear), or tags = null (a
				// known null → recreate to clear), each differ from prior and
				// trigger RequiresReplace. UseStateForUnknown only applies to
				// unknown values, so the omitted-tags no-op is preserved while
				// every explicit change still recreates.
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
					mapplanmodifier.RequiresReplace(),
				},
			},
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "The zone ID of the Storage Set. Changing this will create a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *storageSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan storageSetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()
	action := "CreateStorageSet"
	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"ClientToken": helper.BuildClientToken(action),
		"ZoneId":      plan.ZoneId.ValueString(),
	}
	// Optional+Computed attributes arrive as unknown in the create plan when
	// the user omits them (UseStateForUnknown has no prior state to draw on
	// during create). Send only known, non-null values: ValueString/ValueInt64
	// on an unknown value yield "" / 0, which would otherwise send empty or
	// zero parameters the user never configured.
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		request["Description"] = plan.Description.ValueString()
	}
	if !plan.MaxPartitionNumber.IsNull() && !plan.MaxPartitionNumber.IsUnknown() {
		request["MaxPartitionNumber"] = plan.MaxPartitionNumber.ValueInt64()
	}
	if !plan.ResourceGroupId.IsNull() && !plan.ResourceGroupId.IsUnknown() {
		request["ResourceGroupId"] = plan.ResourceGroupId.ValueString()
	}
	if !plan.StorageSetName.IsNull() && !plan.StorageSetName.IsUnknown() {
		request["StorageSetName"] = plan.StorageSetName.ValueString()
	}
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		tagCount := 1
		for key, value := range plan.Tags.Elements() {
			strVal, ok := value.(types.String)
			if !ok {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Creating %s: invalid tag value", storageSetTypeName),
					fmt.Sprintf("the tag value for key %q has an unexpected type %T, expected a string", key, value),
				)
				return
			}
			if strVal.IsUnknown() {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Creating %s: unknown tag value", storageSetTypeName),
					fmt.Sprintf("the tag value for key %q is unknown and cannot be sent to CreateStorageSet", key),
				)
				return
			}
			// A null element value (e.g. a *string nil mapped via
			// MapValueFrom) is invalid and must not be sent: ValueString on a
			// null String yields "" and would silently create an empty tag the
			// user never supplied. An explicit empty string "" is a valid tag
			// value and is NOT null, so it is sent below.
			if strVal.IsNull() {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Creating %s: null tag value", storageSetTypeName),
					fmt.Sprintf("the tag value for key %q is null; use an empty string for an empty tag value", key),
				)
				return
			}
			request[fmt.Sprintf("Tag.%d.Key", tagCount)] = key
			request[fmt.Sprintf("Tag.%d.Value", tagCount)] = strVal.ValueString()
			tagCount++
		}
	}

	response, err := rpcPostWithRetry(client, action, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling %s", storageSetTypeName, action),
			errs.WrapErrorf(err, "Creating %s: calling %s", storageSetTypeName, action).Error(),
		)
		return
	}

	storageSetId, _ := response["StorageSetId"].(string)
	if storageSetId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling %s", storageSetTypeName, action),
			fmt.Sprintf("The response carried no StorageSetId: %v", response),
		)
		return
	}

	// Persist the known identity and planned attributes before refreshing, so a
	// transient DescribeStorageSets failure does not orphan the just-created
	// resource: Terraform keeps the state below and a later refresh or import
	// reconciles it. The full state is overwritten once Describe succeeds.
	resp.Diagnostics.Append(resp.State.Set(ctx, storageSetCreateState(plan, client, storageSetId))...)

	object, err := describeStorageSet(client, storageSetId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling DescribeStorageSets", storageSetTypeName),
			errs.WrapErrorf(err, "Creating %s: calling DescribeStorageSets", storageSetTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling %s", storageSetTypeName, action),
			fmt.Sprintf("The storage set %s is not visible via DescribeStorageSets right after creation", storageSetId),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, storageSetState(ctx, object))...)
}

func (r *storageSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state storageSetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := describeStorageSet(r.Client(), state.Id.ValueString())
	if err != nil {
		// Persist the existing state before returning the error so a transient
		// DescribeStorageSets failure does not orphan the resource: leaving
		// resp.State unset would make Terraform treat the resource as removed.
		// A malformed response (non-slice, non-object, missing or mismatched id)
		// reaches here as an error and must preserve the current state rather
		// than silently dropping it.
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeStorageSets", storageSetTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeStorageSets", storageSetTypeName).Error(),
		)
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, storageSetState(ctx, object))...)
}

func (r *storageSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state storageSetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()
	id := state.Id.ValueString()

	// ModifyStorageSetAttribute for description or storage_set_name changes.
	// description and storage_set_name are non-clearable (API min length 2),
	// so only send a field when it is known, non-null and actually different
	// from state. An unchanged Optional+Computed attribute may be unknown in
	// the update plan; sending ValueString on an unknown value would emit an
	// empty string and clear/validate-fail the field.
	descChanged := !plan.Description.IsUnknown() && !plan.Description.Equal(state.Description)
	nameChanged := !plan.StorageSetName.IsUnknown() && !plan.StorageSetName.Equal(state.StorageSetName)
	if descChanged || nameChanged {
		action := "ModifyStorageSetAttribute"
		request := map[string]interface{}{
			"StorageSetId": id,
			"RegionId":     client.RegionId,
			"ClientToken":  helper.BuildClientToken(action),
		}
		if descChanged && !plan.Description.IsNull() {
			request["Description"] = plan.Description.ValueString()
		}
		if nameChanged && !plan.StorageSetName.IsNull() {
			request["StorageSetName"] = plan.StorageSetName.ValueString()
		}
		if _, err := rpcPostWithRetry(client, action, request); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling %s", storageSetTypeName, action),
				errs.WrapErrorf(err, "Updating %s: calling %s", storageSetTypeName, action).Error(),
			)
			return
		}
	}

	// resource_group_id and tags are ForceNew (RequiresReplace): ECS Storage Set
	// is not a supported ResourceType of JoinResourceGroup, TagResources or
	// UntagResources — those APIs return InvalidResourceType.NotFound for
	// ResourceType=storageset (confirmed by ACC Step3 404, RequestId
	// 01A1437F-0FD2-51A8-9E3D-539E1176BD75, and the authoritative
	// JoinResourceGroup/TagResources ResourceType enum). Neither attribute can
	// be updated in place: CreateStorageSet accepts ResourceGroupId and Tag at
	// creation, and ModifyStorageSetAttribute only accepts StorageSetName and
	// Description. Terraform therefore recreates the resource on any change to
	// resource_group_id or tags instead of routing an unsupported Update here,
	// so the Update path issues only ModifyStorageSetAttribute for name/desc
	// and then reads back. Omitting tags (unknown plan) is a no-op that
	// preserves the cloud tags; an explicit known change (including tags = {}
	// to clear) differs from prior and triggers RequiresReplace.

	object, err := describeStorageSet(client, id)
	if err != nil {
		// Preserve the current state on a transient DescribeStorageSets error
		// after the update mutations; otherwise the resource would be dropped.
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeStorageSets", storageSetTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeStorageSets", storageSetTypeName).Error(),
		)
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, storageSetState(ctx, object))...)
}

func (r *storageSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state storageSetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()
	action := "DeleteStorageSet"
	request := map[string]interface{}{
		"StorageSetId": state.Id.ValueString(),
		"RegionId":     client.RegionId,
		"ClientToken":  helper.BuildClientToken(action),
	}

	_, err := rpcPostWithRetry(client, action, request)
	if err != nil && !isStorageSetNotFound(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling %s", storageSetTypeName, action),
			errs.WrapErrorf(err, "Deleting %s: calling %s", storageSetTypeName, action).Error(),
		)
		return
	}
}

func (r *storageSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// storageSetCreateState builds a partial resource model carrying the known
// identity and planned attributes right after CreateStorageSet returns. It is
// persisted before the post-create Describe so a transient read failure does
// not lose the freshly created StorageSetId from state.
func storageSetCreateState(plan storageSetResourceModel, client *connectivity.AliyunClient, storageSetId string) storageSetResourceModel {
	plan.Id = types.StringValue(storageSetId)
	plan.StorageSetId = types.StringValue(storageSetId)
	plan.RegionId = types.StringValue(client.RegionId)
	return plan
}

// describeStorageSet queries DescribeStorageSets by storage set ID and returns
// the matching StorageSet object, or (nil, nil) only when the resource no
// longer exists. Not-found is established strictly against the DescribeStorageSets
// contract (Ecs 2014-05-26): the authoritative signal is the 403
// InvalidStorageSetId.NotFound error code. A 200 response is treated as
// not-found only when the match count confirms zero — TotalCount ("存储集总数",
// always present in the responseDemo) is 0 and the top-level StorageSets key or
// the inner StorageSet key is absent, or the StorageSet array is explicitly
// empty. An explicit empty array is a legitimate not-found regardless of
// TotalCount: the API returned no matching entries for the requested ID.
//
// Every other 200 shape is inconsistent and returns an error so the caller
// preserves state instead of dropping the resource: a positive TotalCount
// without a StorageSets/StorageSet container (the API claims matches but
// returned none), a StorageSets value that is present but not an object, a
// StorageSet value that is present but not a list, a non-object entry, an entry
// missing StorageSetId, or an entry whose StorageSetId does not match the
// requested ID (the API filters by StorageSetIds, so a mismatched ID signals
// corruption, not absence). When TotalCount is omitted entirely the response is
// not a documented shape, so absence cannot be established from the count and
// the response is treated as inconsistent — an error preserves state rather than
// risking a silent drop. NotFound string matching is kept narrow
// (InvalidStorageSetId.NotFound) and never widened.
func describeStorageSet(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	action := "DescribeStorageSets"
	request := map[string]interface{}{
		"StorageSetIds": fmt.Sprintf("[\"%s\"]", id),
		"RegionId":      client.RegionId,
		"ClientToken":   helper.BuildClientToken(action),
	}

	response, err := rpcPostWithRetry(client, action, request)
	if err != nil {
		if isStorageSetNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	totalCountRaw, hasTotalCount := response["TotalCount"]
	totalCount, totalCountOk := parseTotalCount(totalCountRaw)

	// A count of zero corroborates an empty result, so an absent StorageSets or
	// StorageSet key is a legitimate not-found. A positive count without the
	// container, or a response that omits/malforms TotalCount so the count
	// cannot confirm absence, is inconsistent: return an error so the caller
	// preserves state/ID instead of dropping the resource.
	storageSetsRaw, hasKey := response["StorageSets"]
	if !hasKey {
		if totalCountOk && totalCount == 0 {
			return nil, nil
		}
		return nil, inconsistentStorageSetCountError(id, "StorageSets key", "absent", totalCountRaw, totalCount, totalCountOk, hasTotalCount)
	}
	storageSetsObj, ok := storageSetsRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("the DescribeStorageSets response StorageSets is not an object: %T", storageSetsRaw)
	}
	setsRaw, hasSets := storageSetsObj["StorageSet"]
	if !hasSets {
		// StorageSets present but the StorageSet key is absent: legitimate
		// not-found only when the count confirms zero; otherwise inconsistent.
		if totalCountOk && totalCount == 0 {
			return nil, nil
		}
		return nil, inconsistentStorageSetCountError(id, "StorageSet key", "absent", totalCountRaw, totalCount, totalCountOk, hasTotalCount)
	}
	sets, ok := setsRaw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("the DescribeStorageSets response StorageSet is not a list: %T", setsRaw)
	}
	if len(sets) == 0 {
		// An explicit empty array is the legitimate not-found result: the API
		// returned no matching entries for the requested ID.
		return nil, nil
	}
	mismatched := false
	for _, raw := range sets {
		object, ok := raw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("the DescribeStorageSets response contains a non-object StorageSet entry: %T", raw)
		}
		sid, hasID := object["StorageSetId"]
		if !hasID {
			return nil, fmt.Errorf("the DescribeStorageSets response contains a StorageSet entry without StorageSetId: %v", object)
		}
		gotID, _ := sid.(string)
		if gotID == id {
			return object, nil
		}
		// The response carried an id we did not request: do not treat that as
		// not-found; record it so the caller errors and preserves state.
		mismatched = true
	}
	if mismatched {
		return nil, fmt.Errorf("the DescribeStorageSets response returned no StorageSet with the requested id %q", id)
	}
	// Unreachable given len(sets) > 0 and every entry is validated above, but
	// keep a defensive error rather than a silent not-found.
	return nil, fmt.Errorf("the DescribeStorageSets response returned no StorageSet with the requested id %q", id)
}

// parseTotalCount extracts the TotalCount ("存储集总数") field from a
// DescribeStorageSets response. It returns the count and ok=true only when
// TotalCount is present AND a valid integer. The SDK RpcPost decodes responses
// with json.Decoder.UseNumber(), so TotalCount arrives as json.Number; float64,
// int, int64 and string encodings are also handled for robustness. A missing or
// non-integer TotalCount yields ok=false so describeStorageSet cannot use the
// count to establish absence and handles it conservatively.
func parseTotalCount(raw interface{}) (count int64, ok bool) {
	switch n := raw.(type) {
	case json.Number:
		if s := n.String(); s != "" {
			if c, err := n.Int64(); err == nil {
				return c, true
			}
		}
	case float64:
		// Reject non-integers: a fractional TotalCount is malformed and must
		// not be truncated to a misleading count.
		if n == float64(int64(n)) {
			return int64(n), true
		}
	case int:
		return int64(n), true
	case int64:
		return n, true
	case string:
		if n != "" {
			if c, err := strconv.ParseInt(n, 10, 64); err == nil {
				return c, true
			}
		}
	}
	return 0, false
}

// inconsistentStorageSetCountError reports a 200 DescribeStorageSets response
// whose match count does not corroborate the data container, so the resource
// must not be dropped. It distinguishes a positive count without a container
// (the API claims matches but returned none), a present-but-malformed TotalCount,
// and an omitted TotalCount (not a documented DescribeStorageSets shape, so
// absence cannot be established from the count).
func inconsistentStorageSetCountError(id, container, state string, totalCountRaw interface{}, totalCount int64, totalCountOk, hasTotalCount bool) error {
	switch {
	case totalCountOk && totalCount > 0:
		return fmt.Errorf("DescribeStorageSets returned TotalCount=%d for storage set %q but the %s is %s: a positive-count response must not be treated as not-found, the resource state is preserved", totalCount, id, container, state)
	case !totalCountOk && hasTotalCount:
		return fmt.Errorf("DescribeStorageSets returned a non-integer TotalCount (%T %v) for storage set %q while the %s is %s: the resource state is preserved", totalCountRaw, totalCountRaw, id, container, state)
	default:
		return fmt.Errorf("DescribeStorageSets omitted TotalCount for storage set %q while the %s is %s: absence cannot be established from the count, the resource state is preserved (DescribeStorageSets Ecs 2014-05-26 responseDemo always carries TotalCount)", id, container, state)
	}
}

// StorageSetExists reports whether a storage set with the given id is still
// visible via DescribeStorageSets. It is exported so acceptance-test destroy
// checks reuse the same strict not-found logic as Read instead of re-deriving
// it: not-found (the InvalidStorageSetId.NotFound error, TotalCount=0 with an
// absent StorageSets/StorageSet key, or an explicit empty array) yields
// (false, nil), while a positive-count response without a container, an omitted
// or malformed TotalCount, or a mismatched entry surfaces as (false, error) so
// state is never dropped on inconsistent data.
func StorageSetExists(client *connectivity.AliyunClient, id string) (bool, error) {
	object, err := describeStorageSet(client, id)
	if err != nil {
		return false, err
	}
	return object != nil, nil
}

// storageSetState maps a DescribeStorageSets StorageSet object onto the
// resource model. Empty strings become null for optional and computed fields.
func storageSetState(ctx context.Context, object map[string]interface{}) storageSetResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}

	tags := tagsFromObject(ctx, object)

	var maxPartition types.Int64
	if v, ok := object["StorageSetPartitionNumber"]; ok {
		switch n := v.(type) {
		case json.Number:
			// SDK RpcPost decodes responses with json.Decoder.UseNumber(), so
			// numeric fields arrive as json.Number (a string-backed type),
			// not float64. Parse as int64 without truncating non-integers or
			// coercing invalid values to 0: a non-integer or out-of-range
			// number stays null (field absent) rather than a wrong int.
			if s := n.String(); s != "" {
				if parsed, err := n.Int64(); err == nil {
					maxPartition = types.Int64Value(parsed)
				}
			}
		case float64:
			maxPartition = types.Int64Value(int64(n))
		case int:
			maxPartition = types.Int64Value(int64(n))
		case int64:
			maxPartition = types.Int64Value(n)
		case string:
			if n != "" {
				if parsed, err := strconv.ParseInt(n, 10, 64); err == nil {
					maxPartition = types.Int64Value(parsed)
				}
			}
		}
	}
	if maxPartition.IsNull() {
		maxPartition = types.Int64Null()
	}

	id := str("StorageSetId")
	return storageSetResourceModel{
		Id:                 id,
		Description:        str("Description"),
		MaxPartitionNumber: maxPartition,
		RegionId:           str("RegionId"),
		ResourceGroupId:    str("ResourceGroupId"),
		StorageSetId:       id,
		StorageSetName:     str("StorageSetName"),
		Tags:               tags,
		ZoneId:             str("ZoneId"),
	}
}

// tagsFromObject extracts tags from a DescribeStorageSets StorageSet object
// and returns a types.Map. System tags (aliyun, acs:, etc.) are filtered out.
func tagsFromObject(ctx context.Context, object map[string]interface{}) types.Map {
	tagsRaw, _ := object["Tags"].(map[string]interface{})
	tagList, _ := tagsRaw["Tag"].([]interface{})
	if len(tagList) == 0 {
		// No Tag entries: the cloud holds no tags. Return an empty KNOWN map,
		// not null — a null state value differs from a config that sets
		// tags = {} (an empty known map), forcing a perpetual diff and a
		// non-no-op plan. An empty known map equals an explicit {} config, so
		// the plan is stable; omitting tags still arrives as unknown in Update
		// and is preserved separately.
		empty, _ := types.MapValueFrom(ctx, types.StringType, map[string]string{})
		return empty
	}
	tagMap := make(map[string]string, len(tagList))
	for _, t := range tagList {
		tag, _ := t.(map[string]interface{})
		if tag == nil {
			continue
		}
		key, _ := tag["Key"].(string)
		value, _ := tag["Value"].(string)
		if key == "" || isIgnoredTagKey(key) {
			continue
		}
		tagMap[key] = value
	}
	if len(tagMap) == 0 {
		// All entries were system/empty keys: treat as no user tags — empty
		// known map, not null (see the len(tagList)==0 branch).
		empty, _ := types.MapValueFrom(ctx, types.StringType, map[string]string{})
		return empty
	}
	m, diags := types.MapValueFrom(ctx, types.StringType, tagMap)
	_ = diags // MapValueFrom does not error for valid string maps
	return m
}

// rpcPostWithRetry calls client.RpcPost with retry on transient errors,
// matching the legacy SDKv2 behavior.
func rpcPostWithRetry(client *connectivity.AliyunClient, action string, request map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = sdkretry.Retry(5*time.Minute, func() *sdkretry.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return sdkretry.RetryableError(err)
			}
			return sdkretry.NonRetryableError(err)
		}
		return nil
	})
	return response, err
}

// isStorageSetNotFound reports whether err is the storage-set-not-found error.
func isStorageSetNotFound(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "InvalidStorageSetId.NotFound")
}

// isIgnoredTagKey reports whether a tag key is an Alibaba Cloud system tag
// that should not be tracked in user state.
func isIgnoredTagKey(key string) bool {
	filters := []string{"^aliyun", "^acs:", "^http://", "^https://", "^sae.do.not.delete"}
	for _, f := range filters {
		if ok, _ := regexp.MatchString(f, key); ok {
			return true
		}
	}
	return false
}
