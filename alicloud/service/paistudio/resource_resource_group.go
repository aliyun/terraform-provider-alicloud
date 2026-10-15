package paistudio

import (
	"context"
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const resourceGroupTypeName = "alicloud_pai_resource_group"

// PaiStudio version/path constants for the ROA-style ResourceGroup APIs.
// CreateResourceGroup is POST /api/v1/resources; Get/Update/Delete are
// /api/v1/resources/{ResourceGroupID} with GET/PUT/DELETE respectively. The
// product/host is resolved by roaRequest (see client.go), not a product code.
const (
	paiStudioVersion  = "2022-01-12"
	resourceGroupPath = "/api/v1/resources"
)

var (
	_ resource.Resource                = &resourceGroupResource{}
	_ resource.ResourceWithImportState = &resourceGroupResource{}
)

func NewResourceGroupResource() resource.Resource {
	return &resourceGroupResource{}
}

type resourceGroupResource struct {
	fwadapt.ResourceBase
}

// resourceGroupModel is the tfsdk shape of alicloud_pai_resource_group.
//
// Several attributes are Optional+Computed because the API reads them back
// (Create-only fields like resource_type are also ForceNew). For such
// attributes a config-omitted value arrives as an unknown at apply time; the
// Create/Update methods therefore guard with IsUnknown and skip sending, so
// the API default applies and the computed value is refreshed on Read.
//
// The write-only `unbind` attribute is never returned by GetResourceGroup, so
// it is preserved from the prior state on Read and from the plan on
// Create/Update rather than refreshed from the API.
type resourceGroupModel struct {
	Id                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	Description               types.String `tfsdk:"description"`
	ResourceType              types.String `tfsdk:"resource_type"`
	ComputingResourceProvider types.String `tfsdk:"computing_resource_provider"`
	Version                   types.String `tfsdk:"version"`
	Unbind                    types.Bool   `tfsdk:"unbind"`
	UserVpc                   types.Object `tfsdk:"user_vpc"`
	Tags                      types.Map    `tfsdk:"tags"`
	Status                    types.String `tfsdk:"status"`
	ClusterId                 types.String `tfsdk:"cluster_id"`
	CreatorId                 types.String `tfsdk:"creator_id"`
	WorkspaceId               types.String `tfsdk:"workspace_id"`
	SupportRdma               types.Bool   `tfsdk:"support_rdma"`
	GmtCreatedTime            types.String `tfsdk:"gmt_created_time"`
	GmtModifiedTime           types.String `tfsdk:"gmt_modified_time"`
}

func (r *resourceGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = resourceGroupTypeName
}

func (r *resourceGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PAI Studio Resource Group resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the Resource Group (ResourceGroupID).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the Resource Group.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the Resource Group.",
			},
			"resource_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The compute resource type of the Resource Group. Valid values: `Ecs` (general compute) and `Lingjun` (Lingjun intelligent compute).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"computing_resource_provider": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The computing resource provider of the Resource Group. Valid values: `Ecs` and `Lingjun`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"version": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The version of the Resource Group. It takes effect only when `resource_type` is `Ecs`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"unbind": schema.BoolAttribute{
				Optional: true,
				Description: "Whether to unbind the Resource Group from the currently connected VPC. " +
					"Write-only: the API does not return this attribute on read, so the value is preserved from configuration.",
			},
			"tags": schema.MapAttribute{
				ElementType:   types.StringType,
				Optional:      true,
				Computed:      true,
				Description:   "The tags of the Resource Group. At most 20 entries.",
				PlanModifiers: []planmodifier.Map{mapRequiresReplace()},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the Resource Group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cluster_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the cluster to which the Resource Group belongs.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"creator_id": schema.StringAttribute{
				Computed:    true,
				Description: "The user ID of the Resource Group creator.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workspace_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the workspace.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"support_rdma": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the Resource Group supports RDMA.",
			},
			"gmt_created_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the Resource Group, in ISO 8601 format (UTC).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gmt_modified_time": schema.StringAttribute{
				Computed:    true,
				Description: "The modification time of the Resource Group, in ISO 8601 format (UTC).",
				// No UseStateForUnknown: UpdateResourceGroup bumps the
				// timestamp, so a planned state value can never survive apply.
			},
		},
		// user_vpc and default_forward_info are modeled as nested BLOCKS rather
		// than nested attributes: the provider is muxed behind tf5muxserver
		// (protocol v5), and protocol v5 cannot convert a nested attribute
		// (SingleNestedAttribute/Attributes). Nested blocks are proto5-native, so
		// the schema loads cleanly. The model field is still types.Object and the
		// reflection helpers in helpers.go are unchanged.
		Blocks: map[string]schema.Block{
			"user_vpc": schema.SingleNestedBlock{
				Description: "The VPC connection information of the Resource Group.",
				Attributes: map[string]schema.Attribute{
					"vpc_id": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The ID of the VPC.",
					},
					"switch_id": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The ID of the vSwitch.",
					},
					"security_group_id": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The ID of the security group.",
					},
					"role_arn": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The ARN of the role used for OSS access authorization.",
					},
					"default_route": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The default route. Valid values: `eth0` (public gateway) and `eth1` (user ENI).",
					},
					"extended_cidrs": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Computed:    true,
						Description: "The extended CIDR blocks.",
					},
				},
				Blocks: map[string]schema.Block{
					"default_forward_info": schema.SingleNestedBlock{
						Description: "The default public network gateway information.",
						Attributes: map[string]schema.Attribute{
							"eip_allocation_id": schema.StringAttribute{
								Optional:    true,
								Computed:    true,
								Description: "The EIP allocation ID.",
							},
							"nat_gateway_id": schema.StringAttribute{
								Optional:    true,
								Computed:    true,
								Description: "The NAT gateway ID.",
							},
						},
					},
				},
			},
		},
	}
}

func (r *resourceGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{}
	setString := func(key string, v types.String) {
		// Skip null and unknown: Optional+Computed attributes arrive as unknown
		// when the user omits them, and the API default should apply.
		if v.IsNull() || v.IsUnknown() {
			return
		}
		request[key] = v.ValueString()
	}
	setString("Name", plan.Name)
	setString("Description", plan.Description)
	setString("ResourceType", plan.ResourceType)
	setString("ComputingResourceProvider", plan.ComputingResourceProvider)
	setString("Version", plan.Version)

	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		request["Tag"] = tagsMapToTagArray(plan.Tags)
	}
	if !plan.UserVpc.IsNull() && !plan.UserVpc.IsUnknown() {
		if uv, ok := userVpcToObject(ctx, plan.UserVpc, &resp.Diagnostics); ok {
			request["UserVpc"] = uv
		}
		if resp.Diagnostics.HasError() {
			return
		}
	}

	response, err := roaRequest(r.Client(), "POST", resourceGroupPath, request, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateResourceGroup", resourceGroupTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateResourceGroup", resourceGroupTypeName).Error(),
		)
		return
	}

	id := resourceGroupIDFromResponse(response)
	if id == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateResourceGroup", resourceGroupTypeName),
			fmt.Sprintf("The response carried no ResourceGroupID: %v", response),
		)
		return
	}

	object, err := getResourceGroup(r.Client(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetResourceGroup", resourceGroupTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetResourceGroup", resourceGroupTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateResourceGroup", resourceGroupTypeName),
			fmt.Sprintf("The resource group %s is not visible via GetResourceGroup right after creation", id),
		)
		return
	}

	state := resourceGroupState(object, id)
	// Optional+Computed request attributes the API may not echo on Get (e.g.
	// ComputingResourceProvider is empty for Ecs resource groups): keep the
	// configured value so the post-apply state matches the plan.
	state.Name = echoString(state.Name, plan.Name)
	state.Description = echoString(state.Description, plan.Description)
	state.ResourceType = echoString(state.ResourceType, plan.ResourceType)
	state.ComputingResourceProvider = echoString(state.ComputingResourceProvider, plan.ComputingResourceProvider)
	state.Version = echoString(state.Version, plan.Version)
	// unbind is write-only: preserve the configured value; it is not part of
	// the Get response.
	state.Unbind = plan.Unbind
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *resourceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := getResourceGroup(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetResourceGroup", resourceGroupTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetResourceGroup", resourceGroupTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	next := resourceGroupState(object, state.Id.ValueString())
	// Re-apply the echo fallback against the prior state so a refresh does not
	// wipe values the API omits (e.g. ComputingResourceProvider).
	next.Name = echoString(next.Name, state.Name)
	next.Description = echoString(next.Description, state.Description)
	next.ResourceType = echoString(next.ResourceType, state.ResourceType)
	next.ComputingResourceProvider = echoString(next.ComputingResourceProvider, state.ComputingResourceProvider)
	next.Version = echoString(next.Version, state.Version)
	// unbind is write-only: carry the prior state value forward, never refresh.
	next.Unbind = state.Unbind
	resp.Diagnostics.Append(resp.State.Set(ctx, next)...)
}

func (r *resourceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state resourceGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{}
	// For Optional+Computed updatable attributes, an unknown plan value means
	// the user removed the attribute from config; skip so the API keeps its value.
	if !plan.Name.IsUnknown() && !plan.Name.Equal(state.Name) {
		request["Name"] = plan.Name.ValueString()
	}
	if !plan.Description.IsUnknown() && !plan.Description.Equal(state.Description) {
		request["Description"] = plan.Description.ValueString()
	}
	// Unbind is a write-only action parameter present only on Update. Send it
	// whenever the configured value changes.
	if !plan.Unbind.IsUnknown() && !plan.Unbind.Equal(state.Unbind) {
		request["Unbind"] = plan.Unbind.ValueBool()
	}
	// UserVpc is updatable; send the whole object when it changes.
	if !plan.UserVpc.IsUnknown() && !plan.UserVpc.Equal(state.UserVpc) {
		if plan.UserVpc.IsNull() {
			request["UserVpc"] = map[string]interface{}{}
		} else if uv, ok := userVpcToObject(ctx, plan.UserVpc, &resp.Diagnostics); ok {
			request["UserVpc"] = uv
		}
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if len(request) > 0 {
		pathName := fmt.Sprintf("/api/v1/resources/%s", state.Id.ValueString())
		if _, err := roaRequest(r.Client(), "PUT", pathName, request, true); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling UpdateResourceGroup", resourceGroupTypeName),
				errs.WrapErrorf(err, "Updating %s: calling UpdateResourceGroup", resourceGroupTypeName).Error(),
			)
			return
		}
	}

	object, err := getResourceGroup(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetResourceGroup", resourceGroupTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetResourceGroup", resourceGroupTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	next := resourceGroupState(object, state.Id.ValueString())
	// Re-apply the echo fallback against the plan so values the API omits
	// survive the update.
	next.Name = echoString(next.Name, plan.Name)
	next.Description = echoString(next.Description, plan.Description)
	next.ResourceType = echoString(next.ResourceType, plan.ResourceType)
	next.ComputingResourceProvider = echoString(next.ComputingResourceProvider, plan.ComputingResourceProvider)
	next.Version = echoString(next.Version, plan.Version)
	// unbind is write-only: preserve the configured plan value after the update.
	next.Unbind = plan.Unbind
	resp.Diagnostics.Append(resp.State.Set(ctx, next)...)
}

func (r *resourceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pathName := fmt.Sprintf("/api/v1/resources/%s", state.Id.ValueString())
	_, err := roaRequest(r.Client(), "DELETE", pathName, nil, true)
	if err != nil && !isResourceGroupNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteResourceGroup", resourceGroupTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteResourceGroup", resourceGroupTypeName).Error(),
		)
		return
	}
}

func (r *resourceGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// echoString returns the value read back from the API when the API actually
// echoes it, and falls back to the configured value when the API omits it
// (returns an empty string). GetResourceGroup returns ComputingResourceProvider
// as an empty string for Ecs resource groups even though Create accepts and
// stores the value; without this fallback the post-apply state diverges from
// the plan ("Provider produced inconsistent result after apply"). It is a
// no-op whenever the API echoes a non-empty value, so attributes that are
// reliably read back (Name, Description, ResourceType, Version) keep their
// real API value.
func echoString(apiVal, configuredVal types.String) types.String {
	if !apiVal.IsNull() && !apiVal.IsUnknown() && apiVal.ValueString() != "" {
		return apiVal
	}
	return configuredVal
}

// getResourceGroup queries GetResourceGroup by ResourceGroupID and returns the
// response object, or nil when the resource group no longer exists.
func getResourceGroup(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	pathName := fmt.Sprintf("/api/v1/resources/%s", id)
	response, err := roaRequest(client, "GET", pathName, nil, true)
	if err != nil {
		if isResourceGroupNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	if response == nil {
		return nil, fmt.Errorf("the GetResourceGroup response was empty for id %s", id)
	}
	return response, nil
}

// resourceGroupState maps a GetResourceGroup response onto the resource model.
// The id is the path ResourceGroupID (the Get response carries no ID field),
// so it is threaded in from state. Empty strings become null.
func resourceGroupState(obj map[string]interface{}, id string) resourceGroupModel {
	str := func(key string) types.String {
		if v, ok := obj[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	boolVal := func(key string) types.Bool {
		if v, ok := obj[key].(bool); ok {
			return types.BoolValue(v)
		}
		if v, ok := obj[key].(string); ok {
			return types.BoolValue(strings.EqualFold(v, "true"))
		}
		return types.BoolNull()
	}

	state := resourceGroupModel{
		Id:                        types.StringValue(id),
		Name:                      str("Name"),
		Description:               str("Description"),
		ResourceType:              str("ResourceType"),
		ComputingResourceProvider: str("ComputingResourceProvider"),
		Version:                   str("Version"),
		Status:                    str("Status"),
		ClusterId:                 str("ClusterID"),
		CreatorId:                 str("CreatorID"),
		WorkspaceId:               str("WorkspaceID"),
		SupportRdma:               boolVal("SupportRDMA"),
		GmtCreatedTime:            str("GmtCreatedTime"),
		GmtModifiedTime:           str("GmtModifiedTime"),
		// Default absent nested/complex attributes to explicit nulls rather
		// than zero values so state serialization is well-formed.
		UserVpc: types.ObjectNull(userVpcAttrTypes),
		Tags:    types.MapNull(types.StringType),
	}

	if uv, ok := obj["UserVpc"].(map[string]interface{}); ok && len(uv) > 0 {
		state.UserVpc = userVpcFromObject(uv)
	}
	if tags, ok := obj["Tags"].([]interface{}); ok && len(tags) > 0 {
		state.Tags = tagsArrayToMap(tags)
	}
	return state
}

// resourceGroupIDFromResponse extracts the ResourceGroupID from a
// CreateResourceGroup response, tolerating casing variants.
func resourceGroupIDFromResponse(response map[string]interface{}) string {
	for _, key := range []string{"ResourceGroupID", "ResourceGroupId", "Id", "id"} {
		if v, ok := response[key].(string); ok && v != "" {
			return v
		}
	}
	return ""
}

// isResourceGroupNotExist reports whether err is a not-found error from
// GetResourceGroup/DeleteResourceGroup. errs.NotFoundError covers the ROA 404
// and the NotFound message pattern; a substring fallback guards wrapped forms.
func isResourceGroupNotExist(err error) bool {
	if err == nil {
		return false
	}
	if errs.NotFoundError(err) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "notfound") || strings.Contains(msg, "not exist") || strings.Contains(msg, "notexist")
}
