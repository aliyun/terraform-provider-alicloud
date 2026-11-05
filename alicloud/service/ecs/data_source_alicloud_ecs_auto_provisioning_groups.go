package ecs

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
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

var autoProvisioningGroupStatuses = []string{
	"submitted",
	"active",
	"deleted",
	"delete-running",
	"modifying",
}

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
	Status                    types.List   `tfsdk:"status"`
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
				Description: "A list of auto provisioning group IDs to filter by. A maximum of 20 IDs may be specified: the DescribeAutoProvisioningGroups operation description limits the ID list to 20 entries even though the published schema metadata still declares `maxItems` 100, and the provider follows the documented operation limit. Blank and repeated IDs are sent once, so the result does not depend on the server's duplicate handling. The attribute also reports the IDs that matched.",
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
			"status": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "A list of statuses to filter by. Valid values: `submitted`, `active`, `deleted`, `delete-running` and `modifying`.",
				Validators: []validator.List{
					listvalidator.SizeAtMost(10),
					listvalidator.ValueStringsAre(stringvalidator.OneOf(autoProvisioningGroupStatuses...)),
				},
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
					ecsFwTagMapValidator{allowEmptyValue: true},
				},
			},
			"page_number": schema.Int64Attribute{
				Optional:    true,
				Description: "The number of the first page to read. Valid values: 1 or greater. Defaults to 1. The data source continues with the following pages until the server total is reached or the read budget is exhausted.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"page_size": schema.Int64Attribute{
				Optional:    true,
				Description: "The number of entries per page. Valid values: 1 to 100. Defaults to 100. A page shorter than this value ends the read only when it accounts for every row the server total promised; otherwise the read is reported as incomplete.",
				Validators: []validator.Int64{
					int64validator.Between(1, 100),
				},
			},
			"total_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of auto provisioning groups that match all configured filters, as reported by the `TotalCount` of the API response. It counts every matching group in the region, not only the rows returned from `page_number`. A total that changes between pages, or that disagrees with the rows read, is reported as an error instead of being written.",
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
// diagnostics and the state write, so the same request-to-state path can be
// exercised directly in unit tests.
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

	sorted := append([]map[string]interface{}(nil), objects...)
	sort.Slice(sorted, func(i, j int) bool {
		return ecsFwString(sorted[i]["AutoProvisioningGroupId"]).ValueString() < ecsFwString(sorted[j]["AutoProvisioningGroupId"]).ValueString()
	})

	ids := make([]string, 0, len(sorted))
	names := make([]string, 0, len(sorted))
	elements := make([]attr.Value, 0, len(sorted))
	for _, object := range sorted {
		ids = append(ids, ecsFwString(object["AutoProvisioningGroupId"]).ValueString())
		names = append(names, ecsFwString(object["AutoProvisioningGroupName"]).ValueString())
		element, err := autoProvisioningGroupObject(object)
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
	model.Ids = ecsFwStringListValue(ids)
	model.Names = ecsFwStringListValue(names)
	model.TotalCount = types.Int64Value(totalCount)
	model.Groups = types.ListValueMust(autoProvisioningGroupObjectType, elements)
	return model, nil
}

// autoProvisioningGroupsRequest builds the full DescribeAutoProvisioningGroups
// request from the configuration. It is kept separate from Read so the
// request-to-state mapping can be asserted directly.
func (d *autoProvisioningGroupsDataSource) autoProvisioningGroupsRequest(model autoProvisioningGroupsModel, regionId string, pageNumber, pageSize int64) map[string]interface{} {
	request := map[string]interface{}{
		"RegionId":   regionId,
		"PageNumber": int(pageNumber),
		"PageSize":   int(pageSize),
	}
	ecsFwSetRequestString(request, "AutoProvisioningGroupName", model.AutoProvisioningGroupName)
	ecsFwSetRequestString(request, "ResourceGroupId", model.ResourceGroupId)
	ecsFwApplyRepeatList(request, "AutoProvisioningGroupId", ecsFwElementsAsStrings(model.Ids))
	ecsFwApplyRepeatList(request, "AutoProvisioningGroupStatus", ecsFwElementsAsStrings(model.Status))
	ecsFwApplyRepeatList(request, "AutoProvisioningGroupTypes", ecsFwElementsAsStrings(model.Types))
	ecsFwApplyTags(request, model.Tags)
	return request
}

func (d *autoProvisioningGroupsDataSource) autoProvisioningGroupsDescribeAll(ctx context.Context, client *connectivity.AliyunClient, request map[string]interface{}, startPage, pageSize int64) ([]map[string]interface{}, int64, error) {
	return ecsFwCollectNumberedPages(ctx, autoProvisioningGroupsMaxPages, startPage, pageSize, autoProvisioningGroupsAction, "AutoProvisioningGroupId", func(pageNumber int64) ([]map[string]interface{}, int64, error) {
		request["PageNumber"] = int(pageNumber)
		request["PageSize"] = int(pageSize)
		response, err := ecsFwRpcPost(ctx, client, autoProvisioningGroupsTypeName, autoProvisioningGroupsAction, request)
		if err != nil {
			return nil, 0, err
		}
		items, err := ecsFwDecodeRecordList(response, "AutoProvisioningGroups", "AutoProvisioningGroup", autoProvisioningGroupsAction)
		if err != nil {
			return nil, 0, err
		}
		totalCount, _, err := ecsFwDecodeCount(response, "TotalCount", autoProvisioningGroupsAction)
		if err != nil {
			return nil, 0, err
		}
		return items, totalCount, nil
	})
}

func autoProvisioningGroupObject(object map[string]interface{}) (types.Object, error) {
	identity, err := ecsFwRecordIdentity(object, "AutoProvisioningGroupId", autoProvisioningGroupsAction)
	if err != nil {
		return types.ObjectNull(autoProvisioningGroupAttrTypes), err
	}
	tags, err := ecsFwDecodeTags(object["Tags"])
	if err != nil {
		return types.ObjectNull(autoProvisioningGroupAttrTypes), fmt.Errorf("%s: %w", identity, err)
	}
	launchTemplateConfigs, err := ecsFwOptionalObjectListMember("LaunchTemplateConfigs", "LaunchTemplateConfig", object["LaunchTemplateConfigs"], autoProvisioningGroupLaunchTemplateConfigObjectType, func(item map[string]interface{}) (attr.Value, error) {
		return types.ObjectValueMust(autoProvisioningGroupLaunchTemplateConfigAttrTypes, map[string]attr.Value{
			"instance_type":     ecsFwString(item["InstanceType"]),
			"max_price":         ecsFwString(item["MaxPrice"]),
			"vswitch_id":        ecsFwString(item["VSwitchId"]),
			"weighted_capacity": ecsFwString(item["WeightedCapacity"]),
			"priority":          ecsFwString(item["Priority"]),
		}), nil
	})
	if err != nil {
		return types.ObjectNull(autoProvisioningGroupAttrTypes), fmt.Errorf("%s: %w", identity, err)
	}
	targetCapacitySpecification, err := ecsFwOptionalObjectList("TargetCapacitySpecification", object["TargetCapacitySpecification"], autoProvisioningGroupTargetCapacitySpecificationObjectType, func(specification map[string]interface{}) map[string]attr.Value {
		return map[string]attr.Value{
			"default_target_capacity_type":  ecsFwString(specification["DefaultTargetCapacityType"]),
			"total_target_capacity":         ecsFwString(specification["TotalTargetCapacity"]),
			"pay_as_you_go_target_capacity": ecsFwString(specification["PayAsYouGoTargetCapacity"]),
			"spot_target_capacity":          ecsFwString(specification["SpotTargetCapacity"]),
		}
	})
	if err != nil {
		return types.ObjectNull(autoProvisioningGroupAttrTypes), fmt.Errorf("%s: %w", identity, err)
	}
	spotOptions, err := ecsFwOptionalObjectList("SpotOptions", object["SpotOptions"], autoProvisioningGroupSpotOptionsObjectType, func(options map[string]interface{}) map[string]attr.Value {
		return map[string]attr.Value{
			"allocation_strategy":            ecsFwString(options["AllocationStrategy"]),
			"instance_interruption_behavior": ecsFwString(options["InstanceInterruptionBehavior"]),
			"instance_pools_to_use_count":    ecsFwInt64(options["InstancePoolsToUseCount"]),
		}
	})
	if err != nil {
		return types.ObjectNull(autoProvisioningGroupAttrTypes), fmt.Errorf("%s: %w", identity, err)
	}
	payAsYouGoOptions, err := ecsFwOptionalObjectList("PayAsYouGoOptions", object["PayAsYouGoOptions"], autoProvisioningGroupPayAsYouGoOptionsObjectType, func(options map[string]interface{}) map[string]attr.Value {
		return map[string]attr.Value{
			"allocation_strategy": ecsFwString(options["AllocationStrategy"]),
		}
	})
	if err != nil {
		return types.ObjectNull(autoProvisioningGroupAttrTypes), fmt.Errorf("%s: %w", identity, err)
	}
	candidateOptions, err := ecsFwOptionalObjectList("CandidateOptions", object["CandidateOptions"], autoProvisioningGroupCandidateOptionsObjectType, func(options map[string]interface{}) map[string]attr.Value {
		return map[string]attr.Value{
			"timeout_minutes": ecsFwInt64(options["TimeoutMinutes"]),
		}
	})
	if err != nil {
		return types.ObjectNull(autoProvisioningGroupAttrTypes), fmt.Errorf("%s: %w", identity, err)
	}
	capacitySpecification, err := ecsFwOptionalObjectList("CapacitySpecification", object["CapacitySpecification"], autoProvisioningGroupCapacitySpecificationObjectType, func(specification map[string]interface{}) map[string]attr.Value {
		return map[string]attr.Value{
			"pay_as_you_go_capacity": ecsFwString(specification["PayAsYouGoCapacity"]),
			"pre_paid_capacity":      ecsFwString(specification["PrePaidCapacity"]),
			"spot_capacity":          ecsFwString(specification["SpotCapacity"]),
			"total_capacity":         ecsFwString(specification["TotalCapacity"]),
		}
	})
	if err != nil {
		return types.ObjectNull(autoProvisioningGroupAttrTypes), fmt.Errorf("%s: %w", identity, err)
	}
	suspendedProcesses, err := ecsFwOptionalStringList("SuspendedProcesses", "SuspendedProcess", object["SuspendedProcesses"])
	if err != nil {
		return types.ObjectNull(autoProvisioningGroupAttrTypes), fmt.Errorf("%s: %w", identity, err)
	}
	return types.ObjectValueMust(autoProvisioningGroupAttrTypes, map[string]attr.Value{
		"auto_provisioning_group_id":          ecsFwString(object["AutoProvisioningGroupId"]),
		"auto_provisioning_group_name":        ecsFwString(object["AutoProvisioningGroupName"]),
		"auto_provisioning_group_type":        ecsFwString(object["AutoProvisioningGroupType"]),
		"status":                              ecsFwString(object["Status"]),
		"state":                               ecsFwString(object["State"]),
		"region_id":                           ecsFwString(object["RegionId"]),
		"resource_group_id":                   ecsFwString(object["ResourceGroupId"]),
		"launch_template_id":                  ecsFwString(object["LaunchTemplateId"]),
		"launch_template_version":             ecsFwString(object["LaunchTemplateVersion"]),
		"valid_from":                          ecsFwString(object["ValidFrom"]),
		"valid_until":                         ecsFwString(object["ValidUntil"]),
		"terminate_instances":                 ecsFwBool(object["TerminateInstances"]),
		"terminate_instances_with_expiration": ecsFwBool(object["TerminateInstancesWithExpiration"]),
		"excess_capacity_termination_policy":  ecsFwString(object["ExcessCapacityTerminationPolicy"]),
		"max_spot_price":                      ecsFwString(object["MaxSpotPrice"]),
		"creation_time":                       ecsFwString(object["CreationTime"]),
		"tags":                                tags,
		"launch_template_configs":             launchTemplateConfigs,
		"target_capacity_specification":       targetCapacitySpecification,
		"spot_options":                        spotOptions,
		"pay_as_you_go_options":               payAsYouGoOptions,
		"candidate_options":                   candidateOptions,
		"capacity_specification":              capacitySpecification,
		"suspended_processes":                 suspendedProcesses,
	}), nil
}
