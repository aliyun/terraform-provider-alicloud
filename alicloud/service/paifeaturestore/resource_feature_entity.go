package paifeaturestore

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const featureEntityTypeName = "alicloud_pai_feature_store_feature_entity"

var (
	_ resource.Resource                = &featureEntityResource{}
	_ resource.ResourceWithImportState = &featureEntityResource{}
)

func NewFeatureEntityResource() resource.Resource {
	return &featureEntityResource{}
}

type featureEntityResource struct {
	fwadapt.ResourceBase
}

type featureEntityResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	InstanceId            types.String `tfsdk:"instance_id"`
	ProjectId             types.String `tfsdk:"project_id"`
	FeatureEntityName     types.String `tfsdk:"feature_entity_name"`
	JoinId                types.String `tfsdk:"join_id"`
	ParentFeatureEntityId types.String `tfsdk:"parent_feature_entity_id"`
	FeatureEntityId       types.String `tfsdk:"feature_entity_id"`
	CreateTime            types.String `tfsdk:"create_time"`
	RegionId              types.String `tfsdk:"region_id"`
}

func (r *featureEntityResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = featureEntityTypeName
}

func (r *featureEntityResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a Pai Feature Store Feature Entity resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. Format: `<instance_id>:<feature_entity_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the FeatureStore instance, which can be obtained through the ListInstances interface.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the FeatureStore project, which can be obtained through the ListProjects interface.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"feature_entity_name": schema.StringAttribute{
				Required: true,
				Description: "The feature entity name. Start with a letter or number and contain only letters, numbers, " +
					"underscores, and decimal points. The maximum length cannot exceed 63.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.]*$`),
						"must start with a letter or number and contain only letters, numbers, underscores and decimal points",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"join_id": schema.StringAttribute{
				Required: true,
				Description: "The ID that associates multiple sets of feature sets. Generally, enter the primary key field " +
					"in the user/item table, such as user_id or item_id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"parent_feature_entity_id": schema.StringAttribute{
				Optional:    true,
				Description: "The parent feature entity ID. Empty or 0 indicates a root feature entity.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"feature_entity_id": schema.StringAttribute{
				Computed:    true,
				Description: "The first ID of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region_id": schema.StringAttribute{
				Computed:    true,
				Description: "The region ID of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *featureEntityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan featureEntityResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := map[string]interface{}{
		"InstanceId": plan.InstanceId.ValueString(),
		"RegionId":   r.Client().RegionId,
	}
	body := map[string]interface{}{
		"Name":      plan.FeatureEntityName.ValueString(),
		"ProjectId": plan.ProjectId.ValueString(),
		"JoinId":    plan.JoinId.ValueString(),
	}
	if !plan.ParentFeatureEntityId.IsNull() {
		body["ParentFeatureEntityId"] = plan.ParentFeatureEntityId.ValueString()
	}

	response, err := r.Client().RpcPost("PaiFeatureStore", "2023-06-21", "CreateFeatureEntity", query, body, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateFeatureEntity", featureEntityTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateFeatureEntity", featureEntityTypeName).Error(),
		)
		return
	}

	featureEntityId, _ := response["FeatureEntityId"].(string)
	if featureEntityId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateFeatureEntity", featureEntityTypeName),
			fmt.Sprintf("The response carried no FeatureEntityId: %v", response),
		)
		return
	}

	object, err := getFeatureEntity(r.Client(), plan.InstanceId.ValueString(), featureEntityId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetFeatureEntity", featureEntityTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetFeatureEntity", featureEntityTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateFeatureEntity", featureEntityTypeName),
			fmt.Sprintf("The feature entity %s is not visible via GetFeatureEntity right after creation", featureEntityId),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, featureEntityState(object, plan.InstanceId.ValueString(), featureEntityId, r.Client().RegionId))...)
}

func (r *featureEntityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state featureEntityResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId, featureEntityId, ok := parseFeatureEntityId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: parsing resource id", featureEntityTypeName),
			fmt.Sprintf("The resource id %q is not in the expected <instance_id>:<feature_entity_id> format", state.Id.ValueString()),
		)
		return
	}

	object, err := getFeatureEntity(r.Client(), instanceId, featureEntityId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetFeatureEntity", featureEntityTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetFeatureEntity", featureEntityTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, featureEntityState(object, instanceId, featureEntityId, r.Client().RegionId))...)
}

func (r *featureEntityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// FeatureEntity has no Update API: every input attribute is RequiresReplace,
	// so a planned change replaces the resource rather than reaching Update. Keep
	// the state in lockstep with the cloud by refreshing, mirroring Read.
	var state featureEntityResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId, featureEntityId, ok := parseFeatureEntityId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
		return
	}

	object, err := getFeatureEntity(r.Client(), instanceId, featureEntityId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetFeatureEntity", featureEntityTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetFeatureEntity", featureEntityTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, featureEntityState(object, instanceId, featureEntityId, r.Client().RegionId))...)
}

func (r *featureEntityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state featureEntityResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId, featureEntityId, ok := parseFeatureEntityId(state.Id.ValueString())
	if !ok {
		// Already gone or unparsable: treat as deleted.
		return
	}

	_, err := r.Client().RpcPost("PaiFeatureStore", "2023-06-21", "DeleteFeatureEntity", map[string]interface{}{
		"InstanceId":      instanceId,
		"FeatureEntityId": featureEntityId,
		"RegionId":        r.Client().RegionId,
	}, nil, true)
	if err != nil && !isFeatureEntityNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteFeatureEntity", featureEntityTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteFeatureEntity", featureEntityTypeName).Error(),
		)
		return
	}
}

func (r *featureEntityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getFeatureEntity queries GetFeatureEntity and returns the entity object from
// the response, or nil when the feature entity no longer exists.
func getFeatureEntity(client *connectivity.AliyunClient, instanceId, featureEntityId string) (map[string]interface{}, error) {
	response, err := client.RpcPost("PaiFeatureStore", "2023-06-21", "GetFeatureEntity", map[string]interface{}{
		"InstanceId":      instanceId,
		"FeatureEntityId": featureEntityId,
		"RegionId":        client.RegionId,
	}, nil, true)
	if err != nil {
		if isFeatureEntityNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	if response == nil {
		return nil, fmt.Errorf("the GetFeatureEntity response was empty")
	}
	return response, nil
}

// featureEntityState maps a GetFeatureEntity response onto the resource model.
// InstanceId, FeatureEntityId and RegionId are not present in the response and
// are carried in from the request; empty strings become null.
func featureEntityState(object map[string]interface{}, instanceId, featureEntityId, regionId string) featureEntityResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	return featureEntityResourceModel{
		Id:                    types.StringValue(fmt.Sprintf("%s:%s", instanceId, featureEntityId)),
		InstanceId:            types.StringValue(instanceId),
		FeatureEntityId:       types.StringValue(featureEntityId),
		ProjectId:             str("ProjectId"),
		FeatureEntityName:     str("Name"),
		JoinId:                str("JoinId"),
		ParentFeatureEntityId: str("ParentFeatureEntityId"),
		CreateTime:            str("GmtCreateTime"),
		RegionId:              types.StringValue(regionId),
	}
}

// parseFeatureEntityId splits the composite id "<instance_id>:<feature_entity_id>".
func parseFeatureEntityId(id string) (string, string, bool) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// isFeatureEntityNotExist reports whether err is the feature-entity-not-found
// error. DoConfigCenterRequestErrorProblem is the not-exist error code declared
// by the PaiFeatureStore/FeatureEntity CloudSpec model; it arrives on the
// *tea.SDKError returned by RpcPost. Match only the typed code so unrelated
// config-center errors are not mistaken for not-found.
func isFeatureEntityNotExist(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil && *e.Code == "DoConfigCenterRequestErrorProblem" {
		return true
	}
	return false
}
