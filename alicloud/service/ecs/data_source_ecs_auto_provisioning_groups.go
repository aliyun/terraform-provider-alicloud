package ecs

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	autoProvisioningGroupsTypeName   = "alicloud_ecs_auto_provisioning_groups"
	autoProvisioningGroupsAction     = "DescribeAutoProvisioningGroups"
	autoProvisioningGroupsPageSize   = 100
	autoProvisioningGroupsMaxPages   = 1000
	autoProvisioningGroupsMaxIds     = 20
	autoProvisioningGroupsReadBudget = 30 * time.Minute
)

var autoProvisioningGroupTypes = []string{
	"request",
	"maintain",
	"candidate",
}

var (
	_ datasource.DataSource              = (*autoProvisioningGroupsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*autoProvisioningGroupsDataSource)(nil)
)

func NewAutoProvisioningGroupsDataSource() datasource.DataSource {
	return &autoProvisioningGroupsDataSource{}
}

type autoProvisioningGroupsDataSource struct {
	fwadapt.DataSourceBase
}

type autoProvisioningGroupsModel struct {
	Id                        types.String `tfsdk:"id"`
	Ids                       types.List   `tfsdk:"ids"`
	AutoProvisioningGroupName types.String `tfsdk:"auto_provisioning_group_name"`
	ResourceGroupId           types.String `tfsdk:"resource_group_id"`
	Types                     types.List   `tfsdk:"types"`
	Tags                      types.Map    `tfsdk:"tags"`
	PageNumber                types.Int64  `tfsdk:"page_number"`
	PageSize                  types.Int64  `tfsdk:"page_size"`
	TotalCount                types.Int64  `tfsdk:"total_count"`
	Names                     types.List   `tfsdk:"names"`
	Groups                    types.List   `tfsdk:"groups"`
}

var autoProvisioningGroupLaunchTemplateConfigAttrTypes = map[string]attr.Type{
	"instance_type":     types.StringType,
	"max_price":         types.StringType,
	"vswitch_id":        types.StringType,
	"weighted_capacity": types.StringType,
	"priority":          types.StringType,
}

var autoProvisioningGroupTargetCapacitySpecificationAttrTypes = map[string]attr.Type{
	"default_target_capacity_type":  types.StringType,
	"total_target_capacity":         types.StringType,
	"pay_as_you_go_target_capacity": types.StringType,
	"spot_target_capacity":          types.StringType,
}

var autoProvisioningGroupSpotOptionsAttrTypes = map[string]attr.Type{
	"allocation_strategy":            types.StringType,
	"instance_interruption_behavior": types.StringType,
	"instance_pools_to_use_count":    types.Int64Type,
}

var autoProvisioningGroupPayAsYouGoOptionsAttrTypes = map[string]attr.Type{
	"allocation_strategy": types.StringType,
}

var autoProvisioningGroupCandidateOptionsAttrTypes = map[string]attr.Type{
	"timeout_minutes": types.Int64Type,
}

var autoProvisioningGroupCapacitySpecificationAttrTypes = map[string]attr.Type{
	"pay_as_you_go_capacity": types.StringType,
	"pre_paid_capacity":      types.StringType,
	"spot_capacity":          types.StringType,
	"total_capacity":         types.StringType,
}

var autoProvisioningGroupLaunchTemplateConfigObjectType = types.ObjectType{AttrTypes: autoProvisioningGroupLaunchTemplateConfigAttrTypes}
var autoProvisioningGroupTargetCapacitySpecificationObjectType = types.ObjectType{AttrTypes: autoProvisioningGroupTargetCapacitySpecificationAttrTypes}
var autoProvisioningGroupSpotOptionsObjectType = types.ObjectType{AttrTypes: autoProvisioningGroupSpotOptionsAttrTypes}
var autoProvisioningGroupPayAsYouGoOptionsObjectType = types.ObjectType{AttrTypes: autoProvisioningGroupPayAsYouGoOptionsAttrTypes}
var autoProvisioningGroupCandidateOptionsObjectType = types.ObjectType{AttrTypes: autoProvisioningGroupCandidateOptionsAttrTypes}
var autoProvisioningGroupCapacitySpecificationObjectType = types.ObjectType{AttrTypes: autoProvisioningGroupCapacitySpecificationAttrTypes}

var autoProvisioningGroupAttrTypes = map[string]attr.Type{
	"auto_provisioning_group_id":          types.StringType,
	"auto_provisioning_group_name":        types.StringType,
	"auto_provisioning_group_type":        types.StringType,
	"status":                              types.StringType,
	"state":                               types.StringType,
	"region_id":                           types.StringType,
	"resource_group_id":                   types.StringType,
	"launch_template_id":                  types.StringType,
	"launch_template_version":             types.StringType,
	"valid_from":                          types.StringType,
	"valid_until":                         types.StringType,
	"terminate_instances":                 types.BoolType,
	"terminate_instances_with_expiration": types.BoolType,
	"excess_capacity_termination_policy":  types.StringType,
	"max_spot_price":                      types.StringType,
	"creation_time":                       types.StringType,
	"tags":                                types.MapType{ElemType: types.StringType},
	"launch_template_configs":             types.ListType{ElemType: autoProvisioningGroupLaunchTemplateConfigObjectType},
	"target_capacity_specification":       types.ListType{ElemType: autoProvisioningGroupTargetCapacitySpecificationObjectType},
	"spot_options":                        types.ListType{ElemType: autoProvisioningGroupSpotOptionsObjectType},
	"pay_as_you_go_options":               types.ListType{ElemType: autoProvisioningGroupPayAsYouGoOptionsObjectType},
	"candidate_options":                   types.ListType{ElemType: autoProvisioningGroupCandidateOptionsObjectType},
	"capacity_specification":              types.ListType{ElemType: autoProvisioningGroupCapacitySpecificationObjectType},
	"suspended_processes":                 types.ListType{ElemType: types.StringType},
}

var autoProvisioningGroupObjectType = types.ObjectType{AttrTypes: autoProvisioningGroupAttrTypes}

func (d *autoProvisioningGroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = autoProvisioningGroupsTypeName
}

func (d *autoProvisioningGroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a list of auto provisioning groups in a region.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, derived from the matched auto provisioning group IDs.",
			},
			"ids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of auto provisioning group IDs to filter by. A maximum of 20 IDs may be specified: the DescribeAutoProvisioningGroups operation description limits the ID list to 20 entries even though the published schema metadata still declares `maxItems` 100, and the provider follows the documented operation limit. Blank IDs are dropped and repeated IDs are sent once, so the result does not depend on the server's duplicate handling. The attribute also reports the IDs that matched.",
				Validators: []validator.List{
					listvalidator.SizeAtMost(autoProvisioningGroupsMaxIds),
				},
			},
			"auto_provisioning_group_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the auto provisioning group. The filter is sent to the server.",
			},
			"resource_group_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the resource group to which the auto provisioning group belongs.",
			},
			"types": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "A list of auto provisioning group types to filter by. Valid values: `request`, `maintain` and `candidate`.",
				Validators: []validator.List{
					listvalidator.SizeAtMost(10),
					listvalidator.ValueStringsAre(stringvalidator.OneOf(autoProvisioningGroupTypes...)),
				},
			},
			"tags": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "A mapping of tags to filter by. Each entry is sent as a `Tag` filter. A maximum of 20 entries is allowed. A key must be 1 to 128 characters, must not start with `aliyun` or `acs:`, and must not contain a URL. A value may be empty, is limited to 128 characters and must not contain a URL.",
				Validators: []validator.Map{
					mapvalidator.SizeAtMost(20),
					apgTagMapValidator{},
				},
			},
			"page_number": schema.Int64Attribute{
				Optional:    true,
				Description: "The number of the first page to read. Valid values: 1 to 2147483647 (`PageNumber` is an int32 on the wire). Defaults to 1. The rows of the pages before `page_number` are not read and are not part of the result. When the API reports `TotalCount`, `total_count` retains the API's full matching total and includes those skipped rows; when the total is unknown, it falls back only to the rows actually read. The data source continues with the following pages until the rows that are readable from `page_number` have been read, or the read budget is exhausted.",
				Validators: []validator.Int64{
					int64validator.Between(1, math.MaxInt32),
				},
			},
			"page_size": schema.Int64Attribute{
				Optional:    true,
				Description: "The number of entries per page. Valid values: 1 to 100. Defaults to 100. A page shorter than this value ends the read only when the rows read account for every row that is readable from `page_number`; otherwise the read is reported as incomplete.",
				Validators: []validator.Int64{
					int64validator.Between(1, 100),
				},
			},
			"total_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of auto provisioning groups that match all configured filters, as reported by the `TotalCount` of the API response. When the API reports `TotalCount`, it counts every matching group in the region, including rows before `page_number`. A total that changes between pages, or that disagrees with the rows read, is reported as an error instead of being written. Only when a valid `AutoProvisioningGroups` container is present but `TotalCount` is omitted or null does the value fall back to the number of rows actually read from `page_number`; a response with neither the container nor a usable `TotalCount` is an error.",
			},
			"names": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of auto provisioning group names that match the configured filters.",
			},
			"groups": schema.ListAttribute{
				Computed:    true,
				ElementType: autoProvisioningGroupObjectType,
				Description: "A list of auto provisioning groups that match the configured filters.",
			},
		},
	}
}

func (d *autoProvisioningGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// A configuration that still carries unknown values cannot be turned into a
	// complete request: every filter builder treats an unknown value as absent,
	// so the read would answer a wider question than the one asked and could
	// write that answer to state. Defer until the values are known when the
	// client supports it; otherwise stop before any request is sent.
	if !req.Config.Raw.IsFullyKnown() {
		if req.ClientCapabilities.DeferralAllowed {
			resp.Deferred = &datasource.Deferred{Reason: datasource.DeferredReasonDataSourceConfigUnknown}
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: the configuration is not known yet", autoProvisioningGroupsTypeName),
			"The configuration still contains unknown values, so the filters they define cannot be sent. Apply the configuration that defines them and read the data source again.",
		)
		return
	}

	var model autoProvisioningGroupsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// One deadline covers every page, every retry and the final state write.
	ctx, cancel := context.WithTimeout(ctx, autoProvisioningGroupsReadBudget)
	defer cancel()

	model, err := d.load(ctx, d.Meta(), d.Meta().RegionId, model)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", autoProvisioningGroupsTypeName, autoProvisioningGroupsAction),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// load performs the whole read: it builds the request, follows every page,
// decodes the records and returns the state model. Read only adds the deadline,
// diagnostics and the state write.
func (d *autoProvisioningGroupsDataSource) load(ctx context.Context, client *connectivity.AliyunClient, regionId string, model autoProvisioningGroupsModel) (autoProvisioningGroupsModel, error) {
	pageNumber := int64(1)
	if !model.PageNumber.IsNull() && !model.PageNumber.IsUnknown() {
		pageNumber = model.PageNumber.ValueInt64()
	}
	pageSize := int64(autoProvisioningGroupsPageSize)
	if !model.PageSize.IsNull() && !model.PageSize.IsUnknown() {
		pageSize = model.PageSize.ValueInt64()
	}

	request := d.autoProvisioningGroupsRequest(model, regionId, pageNumber, pageSize)

	objects, totalCount, err := d.autoProvisioningGroupsDescribeAll(ctx, client, request, pageNumber, pageSize)
	if err != nil {
		return model, err
	}

	type matchedGroup struct {
		id     string
		name   string
		object map[string]interface{}
	}
	sorted := make([]matchedGroup, 0, len(objects))
	for _, object := range objects {
		id, err := ecsRecordIdentity(object, autoProvisioningGroupsAction, "AutoProvisioningGroupId")
		if err != nil {
			return model, err
		}
		name, err := ecsDecodeStringMember(object, "AutoProvisioningGroupName")
		if err != nil {
			return model, fmt.Errorf("%s: %w", id, err)
		}
		sorted = append(sorted, matchedGroup{id: id, name: name.ValueString(), object: object})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].id < sorted[j].id
	})

	ids := make([]string, 0, len(sorted))
	names := make([]string, 0, len(sorted))
	elements := make([]attr.Value, 0, len(sorted))
	for _, group := range sorted {
		ids = append(ids, group.id)
		names = append(names, group.name)
		element, err := autoProvisioningGroupObject(group.object)
		if err != nil {
			return model, err
		}
		elements = append(elements, element)
	}

	// A cancelled or expired read must not return a partial result to state.
	if ctxErr := ctx.Err(); ctxErr != nil {
		return model, fmt.Errorf("the read was cancelled or ran out of time before the result was written: %v", ctxErr)
	}

	model.Id = types.StringValue(helper.DataResourceIdHash(ids))
	model.Ids = ecsStringListValue(ids)
	model.Names = ecsStringListValue(names)
	model.TotalCount = types.Int64Value(totalCount)
	model.Groups = types.ListValueMust(autoProvisioningGroupObjectType, elements)
	return model, nil
}

// autoProvisioningGroupsRequest builds the full DescribeAutoProvisioningGroups
// request from the configuration, separate from the read loop that sends it.
func (d *autoProvisioningGroupsDataSource) autoProvisioningGroupsRequest(model autoProvisioningGroupsModel, regionId string, pageNumber, pageSize int64) map[string]interface{} {
	request := map[string]interface{}{
		"RegionId":   regionId,
		"PageNumber": int(pageNumber),
		"PageSize":   int(pageSize),
	}
	ecsSetRequestString(request, "AutoProvisioningGroupName", model.AutoProvisioningGroupName)
	ecsSetRequestString(request, "ResourceGroupId", model.ResourceGroupId)
	ecsApplyRepeatList(request, "AutoProvisioningGroupId", ecsElementsAsStrings(model.Ids))
	// The status filter is not exposed in this release. The group status is
	// still reported in `groups.*.status`.
	ecsApplyRepeatList(request, "AutoProvisioningGroupTypes", ecsElementsAsStrings(model.Types))
	ecsApplyTags(request, model.Tags)
	return request
}

func (d *autoProvisioningGroupsDataSource) autoProvisioningGroupsDescribeAll(ctx context.Context, client *connectivity.AliyunClient, request map[string]interface{}, startPage, pageSize int64) ([]map[string]interface{}, int64, error) {
	return ecsCollectNumberedPages(ctx, autoProvisioningGroupsAction, "AutoProvisioningGroupId", autoProvisioningGroupsMaxPages, startPage, pageSize, func(pageNumber int64) (ecsNumberedPage, error) {
		request["PageNumber"] = int(pageNumber)
		request["PageSize"] = int(pageSize)
		response, err := apgRpcPost(ctx, client, autoProvisioningGroupsTypeName, autoProvisioningGroupsAction, request)
		if err != nil {
			return ecsNumberedPage{}, err
		}
		if err := ecsValidateQueryEnvelope(response, autoProvisioningGroupsAction, "AutoProvisioningGroups", "TotalCount", pageNumber == 1); err != nil {
			return ecsNumberedPage{}, err
		}
		items, err := ecsDecodeRecordList(response, autoProvisioningGroupsAction, "AutoProvisioningGroups", "AutoProvisioningGroup")
		if err != nil {
			return ecsNumberedPage{}, err
		}
		totalCount, totalKnown, err := ecsDecodeCount(response, autoProvisioningGroupsAction, "TotalCount")
		if err != nil {
			return ecsNumberedPage{}, err
		}
		return ecsNumberedPage{items: items, total: totalCount, totalKnown: totalKnown}, nil
	})
}

func autoProvisioningGroupObject(object map[string]interface{}) (types.Object, error) {
	identity, err := ecsRecordIdentity(object, autoProvisioningGroupsAction, "AutoProvisioningGroupId")
	if err != nil {
		return types.ObjectNull(autoProvisioningGroupAttrTypes), err
	}
	attributes, err := autoProvisioningGroupAttributes(object)
	if err != nil {
		return types.ObjectNull(autoProvisioningGroupAttrTypes), fmt.Errorf("%s: %w", identity, err)
	}
	return types.ObjectValueMust(autoProvisioningGroupAttrTypes, attributes), nil
}

// autoProvisioningGroupAttributes decodes every member of one group. The
// numeric members follow the API contract (MaxSpotPrice, the target and current
// capacities and the launch template prices); every other string member is
// decoded strictly by name, so a null or a container fails the record with its
// path instead of being written to state.
func autoProvisioningGroupAttributes(object map[string]interface{}) (map[string]attr.Value, error) {
	attributes := make(map[string]attr.Value, 25)
	if err := ecsDecodeStringFields(object, []ecsStringField{
		{"auto_provisioning_group_id", "AutoProvisioningGroupId", false},
		{"auto_provisioning_group_name", "AutoProvisioningGroupName", false},
		{"auto_provisioning_group_type", "AutoProvisioningGroupType", false},
		{"status", "Status", false},
		{"state", "State", false},
		{"region_id", "RegionId", false},
		{"resource_group_id", "ResourceGroupId", false},
		{"launch_template_id", "LaunchTemplateId", false},
		{"launch_template_version", "LaunchTemplateVersion", false},
		{"valid_from", "ValidFrom", false},
		{"valid_until", "ValidUntil", false},
		{"excess_capacity_termination_policy", "ExcessCapacityTerminationPolicy", false},
		{"max_spot_price", "MaxSpotPrice", true},
		{"creation_time", "CreationTime", false},
	}, attributes); err != nil {
		return nil, err
	}

	tags, err := ecsDecodeTags(object, ecsTagDecodeOptions{reservedKey: ecsProviderReservedTagKey, alternateMembers: true})
	if err != nil {
		return nil, err
	}
	attributes["tags"] = tags

	terminateInstances, err := ecsDecodeBoolMember(object, "TerminateInstances")
	if err != nil {
		return nil, err
	}
	attributes["terminate_instances"] = terminateInstances

	terminateInstancesWithExpiration, err := ecsDecodeBoolMember(object, "TerminateInstancesWithExpiration")
	if err != nil {
		return nil, err
	}
	attributes["terminate_instances_with_expiration"] = terminateInstancesWithExpiration

	launchTemplateConfigs, err := ecsOptionalObjectListMember(object, "LaunchTemplateConfigs", "LaunchTemplateConfig", autoProvisioningGroupLaunchTemplateConfigObjectType, func(item map[string]interface{}) (attr.Value, error) {
		attributes := make(map[string]attr.Value, 5)
		if err := ecsDecodeStringFields(item, []ecsStringField{
			{"instance_type", "InstanceType", false},
			{"max_price", "MaxPrice", true},
			{"vswitch_id", "VSwitchId", false},
			{"weighted_capacity", "WeightedCapacity", true},
			{"priority", "Priority", true},
		}, attributes); err != nil {
			return nil, err
		}
		return types.ObjectValueMust(autoProvisioningGroupLaunchTemplateConfigAttrTypes, attributes), nil
	})
	if err != nil {
		return nil, err
	}
	attributes["launch_template_configs"] = launchTemplateConfigs

	targetCapacitySpecification, err := ecsOptionalObjectList(object, "TargetCapacitySpecification", autoProvisioningGroupTargetCapacitySpecificationObjectType, func(specification map[string]interface{}) (map[string]attr.Value, error) {
		attributes := make(map[string]attr.Value, 4)
		if err := ecsDecodeStringFields(specification, []ecsStringField{
			{"default_target_capacity_type", "DefaultTargetCapacityType", false},
			{"total_target_capacity", "TotalTargetCapacity", true},
			{"pay_as_you_go_target_capacity", "PayAsYouGoTargetCapacity", true},
			{"spot_target_capacity", "SpotTargetCapacity", true},
		}, attributes); err != nil {
			return nil, err
		}
		return attributes, nil
	})
	if err != nil {
		return nil, err
	}
	attributes["target_capacity_specification"] = targetCapacitySpecification

	spotOptions, err := ecsOptionalObjectList(object, "SpotOptions", autoProvisioningGroupSpotOptionsObjectType, func(options map[string]interface{}) (map[string]attr.Value, error) {
		instancePoolsToUseCount, err := ecsDecodeInt64Member(options, "InstancePoolsToUseCount", "SpotOptions.InstancePoolsToUseCount")
		if err != nil {
			return nil, err
		}
		attributes := make(map[string]attr.Value, 3)
		if err := ecsDecodeStringFields(options, []ecsStringField{
			{"allocation_strategy", "AllocationStrategy", false},
			{"instance_interruption_behavior", "InstanceInterruptionBehavior", false},
		}, attributes); err != nil {
			return nil, err
		}
		attributes["instance_pools_to_use_count"] = instancePoolsToUseCount
		return attributes, nil
	})
	if err != nil {
		return nil, err
	}
	attributes["spot_options"] = spotOptions

	payAsYouGoOptions, err := ecsOptionalObjectList(object, "PayAsYouGoOptions", autoProvisioningGroupPayAsYouGoOptionsObjectType, func(options map[string]interface{}) (map[string]attr.Value, error) {
		attributes := make(map[string]attr.Value, 1)
		if err := ecsDecodeStringFields(options, []ecsStringField{
			{"allocation_strategy", "AllocationStrategy", false},
		}, attributes); err != nil {
			return nil, err
		}
		return attributes, nil
	})
	if err != nil {
		return nil, err
	}
	attributes["pay_as_you_go_options"] = payAsYouGoOptions

	candidateOptions, err := ecsOptionalObjectList(object, "CandidateOptions", autoProvisioningGroupCandidateOptionsObjectType, func(options map[string]interface{}) (map[string]attr.Value, error) {
		timeoutMinutes, err := ecsDecodeInt64Member(options, "TimeoutMinutes", "CandidateOptions.TimeoutMinutes")
		if err != nil {
			return nil, err
		}
		return map[string]attr.Value{"timeout_minutes": timeoutMinutes}, nil
	})
	if err != nil {
		return nil, err
	}
	attributes["candidate_options"] = candidateOptions

	capacitySpecification, err := ecsOptionalObjectList(object, "CapacitySpecification", autoProvisioningGroupCapacitySpecificationObjectType, func(specification map[string]interface{}) (map[string]attr.Value, error) {
		attributes := make(map[string]attr.Value, 4)
		if err := ecsDecodeStringFields(specification, []ecsStringField{
			{"pay_as_you_go_capacity", "PayAsYouGoCapacity", true},
			{"pre_paid_capacity", "PrePaidCapacity", true},
			{"spot_capacity", "SpotCapacity", true},
			{"total_capacity", "TotalCapacity", true},
		}, attributes); err != nil {
			return nil, err
		}
		return attributes, nil
	})
	if err != nil {
		return nil, err
	}
	attributes["capacity_specification"] = capacitySpecification

	suspendedProcesses, err := ecsOptionalStringList(object, "SuspendedProcesses", "SuspendedProcess")
	if err != nil {
		return nil, err
	}
	attributes["suspended_processes"] = suspendedProcesses

	return attributes, nil
}

const apgRetryTimeout = 5 * time.Minute

// apgTagMapValidator enforces the key and value limits that the ECS tag
// filters document: the key must be present and must not use the reserved
// aliyun/acs: prefixes, the value length is capped, and neither side may carry
// a URL. DescribeAutoProvisioningGroups accepts an empty tag value, so an empty
// value only has to satisfy the same length and URL rules as any other value.
type apgTagMapValidator struct{}

func (v apgTagMapValidator) Description(context.Context) string {
	return "Each tag key and value must satisfy the tag limits documented for the API."
}

func (v apgTagMapValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateMap checks the key of every entry even when its value is null or
// unknown. A null value is rejected rather than serialized as an empty string:
// the two are different requests, and only an explicit empty string asks the
// API for an empty value.
func (v apgTagMapValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	for key, element := range req.ConfigValue.Elements() {
		if message := ecsValidateTagKey(key); message != "" {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid Tag Key", message)
			return
		}
		value, ok := element.(types.String)
		if !ok {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid Tag Value", fmt.Sprintf("The value of tag key %q must be a string.", key))
			return
		}
		if value.IsUnknown() {
			// Not known yet: it cannot be validated here, and the request
			// builder leaves it out instead of sending an empty value.
			continue
		}
		if value.IsNull() {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid Tag Value", fmt.Sprintf("The value of tag key %q must be set; omit the entry to leave the tag out of the request, or use an empty string to send an empty value.", key))
			return
		}
		if message := validateApgTagValue(value.ValueString()); message != "" {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid Tag Value", message)
			return
		}
	}
}

func validateApgTagValue(value string) string {
	if utf8.RuneCountInString(value) > ecsTagMaxLength {
		return fmt.Sprintf("A tag value must be at most %d characters.", ecsTagMaxLength)
	}
	if strings.Contains(value, "http://") || strings.Contains(value, "https://") {
		return "A tag value must not contain a URL."
	}
	return ""
}

// apgRpcPost sends one read-only ECS query through the provider-wide retry
// policy: transient failures are retried inside the bounded window, everything
// else fails on the first response. The context carries the read budget: it
// gates every attempt and the backoff, and the request itself is built from it,
// so an in-flight call is aborted instead of running on after the read returns.
func apgRpcPost(ctx context.Context, client *connectivity.AliyunClient, typeName, action string, request map[string]interface{}) (map[string]interface{}, error) {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return nil, errs.WrapErrorf(fmt.Errorf("%s was not sent: %v", action, ctxErr), errs.DataDefaultErrorMsg, typeName, action, errs.AlibabaCloudSdkGoERROR)
	}
	response, err := helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
		return client.RpcPostWithContext(ctx, "Ecs", "2014-05-26", action, nil, request)
	}, apgRetryTimeout, ecsNeedRetry)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, typeName, action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)
	return response, nil
}
