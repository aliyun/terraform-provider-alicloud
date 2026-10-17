package pai

import (
	"context"
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const defaultWorkspaceTypeName = "alicloud_pai_workspace_default_workspace"

// AIWorkSpace ROA API coordinates shared with the SDKv2 pai_workspace_workspace
// resource: the same product code, version and REST base. The default workspace
// is the singleton variant living at /api/v1/defaultWorkspaces.
//
// Note: this package deliberately avoids importing
// terraform-plugin-sdk/v2/helper/retry (used by the SDKv2 service helpers),
// because that package transitively registers the "sweep" flag, which collides
// with the terraform-plugin-testing "sweep" flag in a Framework test binary.
// Polling is done with plain time loops and alicloud/helper instead.
const (
	defaultWorkspaceProduct       = "AIWorkSpace"
	defaultWorkspaceVersion       = "2021-02-04"
	defaultWorkspacePath          = "/api/v1/defaultWorkspaces"
	defaultWorkspaceCreateTimeout = 5 * time.Minute
)

var (
	_ resource.Resource                = &defaultWorkspaceResource{}
	_ resource.ResourceWithImportState = &defaultWorkspaceResource{}
)

func NewDefaultWorkspaceResource() resource.Resource {
	return &defaultWorkspaceResource{}
}

// defaultWorkspaceResource is a Framework-native, create-only singleton: a PAI
// default workspace is one per region per account, has no public delete
// operation, and exposes no updatable business fields (UpdateDefaultWorkspace
// only accepts WorkspaceId and changes nothing). Description and EnvTypes are
// therefore RequiresReplace, and the Update/Delete handlers are no-ops by
// design. The routePlan freezes this createOnlyDesign.
type defaultWorkspaceResource struct {
	fwadapt.ResourceBase
}

type defaultWorkspaceResourceModel struct {
	Id              types.String `tfsdk:"id"`
	WorkspaceId     types.String `tfsdk:"workspace_id"`
	Description     types.String `tfsdk:"description"`
	EnvTypes        types.List   `tfsdk:"env_types"`
	WorkspaceName   types.String `tfsdk:"workspace_name"`
	DisplayName     types.String `tfsdk:"display_name"`
	Creator         types.String `tfsdk:"creator"`
	Status          types.String `tfsdk:"status"`
	RegionId        types.String `tfsdk:"region_id"`
	GmtCreateTime   types.String `tfsdk:"gmt_create_time"`
	GmtModifiedTime types.String `tfsdk:"gmt_modified_time"`
	OwnerUserId     types.String `tfsdk:"owner_user_id"`
	OwnerUserName   types.String `tfsdk:"owner_user_name"`
	OwnerUserKp     types.String `tfsdk:"owner_user_kp"`
}

func (r *defaultWorkspaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = defaultWorkspaceTypeName
}

func (r *defaultWorkspaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PAI Workspace default workspace resource. " +
			"The default workspace is a singleton per region per account: it is created once, " +
			"cannot be deleted, and its Description/EnvTypes are immutable after creation (RequiresReplace).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. Same as `workspace_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workspace_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the default workspace.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Required:    true,
				Description: "The description of the default workspace, up to 80 characters. Immutable after creation.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 80),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"env_types": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "The environments contained in the workspace. " +
					"Simple mode only contains `prod`; standard mode contains `dev` and `prod`. Immutable after creation.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"workspace_name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the workspace, unique within the region.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				Computed:    true,
				Description: "The display name of the workspace.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"creator": schema.StringAttribute{
				Computed:    true,
				Description: "The creator of the workspace.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the workspace. Valid values: `INITIALIZING`, `ENABLED`, `FAILURE`, `DISABLE`, `FROZEN`, `UPDATING`.",
				// No UseStateForUnknown: status converges from INITIALIZING to
				// ENABLED after creation, so the planned (unknown) value must
				// not be pinned to state.
			},
			"region_id": schema.StringAttribute{
				Computed:    true,
				Description: "The region ID of the workspace.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gmt_create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The UTC creation time of the workspace, in ISO 8601 format.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gmt_modified_time": schema.StringAttribute{
				Computed:    true,
				Description: "The UTC modification time of the workspace, in ISO 8601 format.",
			},
			"owner_user_id": schema.StringAttribute{
				Computed:    true,
				Description: "The user ID of the workspace owner.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"owner_user_name": schema.StringAttribute{
				Computed:    true,
				Description: "The user name of the workspace owner.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"owner_user_kp": schema.StringAttribute{
				Computed:    true,
				Description: "The user KP of the workspace owner.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *defaultWorkspaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan defaultWorkspaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var envTypes []string
	resp.Diagnostics.Append(plan.EnvTypes.ElementsAs(ctx, &envTypes, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{
		"Description": plan.Description.ValueString(),
		"EnvTypes":    envTypes,
	}

	client := r.Client()
	action := defaultWorkspacePath
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	deadline := time.Now().Add(defaultWorkspaceCreateTimeout)
	var response map[string]interface{}
	var err error
	for {
		response, err = client.RoaPost(defaultWorkspaceProduct, defaultWorkspaceVersion, action, nil, nil, body, true)
		if err == nil {
			break
		}
		if !helper.NeedRetry(err) {
			break
		}
		if time.Now().After(deadline) {
			break
		}
		wait()
	}
	logs.AddDebug(action, response, body)
	if err != nil {
		// Singleton idempotency: the default workspace is one per region per
		// account and cannot be deleted, so a prior run (or out-of-band
		// creation) may already have created it. CreateDefaultWorkspace may
		// surface that as an error; when it does, adopt the existing workspace
		// rather than failing, then wait for ENABLED.
		if existing, getErr := getDefaultWorkspace(client); getErr == nil && existing != nil {
			logs.AddDebug("CreateDefaultWorkspace: adopting an already-existing default workspace", existing, nil)
			object, waitErr := waitDefaultWorkspaceState(client, "ENABLED")
			if waitErr != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Creating %s: waiting for ENABLED", defaultWorkspaceTypeName),
					errs.WrapErrorf(waitErr, "Creating %s: waiting for ENABLED", defaultWorkspaceTypeName).Error(),
				)
				return
			}
			if object == nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Creating %s: waiting for ENABLED", defaultWorkspaceTypeName),
					"the default workspace was not visible via GetDefaultWorkspace after adoption",
				)
				return
			}
			resp.Diagnostics.Append(resp.State.Set(ctx, defaultWorkspaceState(ctx, object, client))...)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateDefaultWorkspace", defaultWorkspaceTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, "default-workspace", action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	// CreateDefaultWorkspace is asynchronous: on success it returns an
	// InstanceJobId (not a WorkspaceId), so the workspace is not yet visible.
	// Block on GetDefaultWorkspace until the default workspace converges to a
	// terminal state. This also adopts a pre-existing singleton: when the
	// default workspace already exists, CreateDefaultWorkspace is an idempotent
	// no-op and GetDefaultWorkspace returns the existing workspace.
	object, err := waitDefaultWorkspaceState(client, "ENABLED")
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for ENABLED", defaultWorkspaceTypeName),
			errs.WrapErrorf(err, "Creating %s: waiting for ENABLED", defaultWorkspaceTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for ENABLED", defaultWorkspaceTypeName),
			"the default workspace was not visible via GetDefaultWorkspace after creation",
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, defaultWorkspaceState(ctx, object, client))...)
}

func (r *defaultWorkspaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state defaultWorkspaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := getDefaultWorkspace(r.Client())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetDefaultWorkspace", defaultWorkspaceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetDefaultWorkspace", defaultWorkspaceTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, defaultWorkspaceState(ctx, object, r.Client()))...)
}

func (r *defaultWorkspaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state defaultWorkspaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// UpdateDefaultWorkspace only accepts WorkspaceId and changes no business
	// field: Description and EnvTypes are RequiresReplace, so this handler is
	// never reached through a configuration change. It is kept to honour the
	// operation mapping and to no-op re-sync the state from the cloud.
	body := map[string]interface{}{
		"WorkspaceId": state.WorkspaceId.ValueString(),
	}
	if _, err := r.Client().RoaPut(defaultWorkspaceProduct, defaultWorkspaceVersion, defaultWorkspacePath, nil, nil, body, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateDefaultWorkspace", defaultWorkspaceTypeName),
			errs.WrapErrorf(err, "Updating %s: calling UpdateDefaultWorkspace", defaultWorkspaceTypeName).Error(),
		)
		return
	}

	object, err := getDefaultWorkspace(r.Client())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetDefaultWorkspace", defaultWorkspaceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetDefaultWorkspace", defaultWorkspaceTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, defaultWorkspaceState(ctx, object, r.Client()))...)
}

func (r *defaultWorkspaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No-op: the default workspace is a singleton with no public delete
	// operation. Removing it from Terraform state leaves the cloud resource
	// intact; it is not re-created on a later apply because Create is
	// idempotent against an already-enabled default workspace.
}

func (r *defaultWorkspaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getDefaultWorkspace calls GetDefaultWorkspace and returns the default
// workspace object at the response root, or nil when no default workspace
// exists (read as "gone" for Terraform state purposes).
func getDefaultWorkspace(client *connectivity.AliyunClient) (map[string]interface{}, error) {
	action := defaultWorkspacePath
	response, err := client.RoaGet(defaultWorkspaceProduct, defaultWorkspaceVersion, action, nil, nil, nil)
	logs.AddDebug(action, response, nil)
	if err != nil {
		if errs.NotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	if _, ok := response["WorkspaceId"].(string); !ok || response["WorkspaceId"] == "" {
		return nil, nil
	}
	return response, nil
}

// waitDefaultWorkspaceState polls GetDefaultWorkspace until Status reaches the
// target (e.g. ENABLED) or a terminal failure state.
func waitDefaultWorkspaceState(client *connectivity.AliyunClient, target string) (map[string]interface{}, error) {
	deadline := time.Now().Add(defaultWorkspaceCreateTimeout)
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	var object map[string]interface{}
	var lastErr error
	for {
		obj, e := getDefaultWorkspace(client)
		if e != nil {
			lastErr = e
			if !helper.NeedRetry(e) {
				return nil, e
			}
		} else if obj != nil {
			object = obj
			lastErr = nil
			status, _ := obj["Status"].(string)
			if status == target {
				return object, nil
			}
			if status == "FAILURE" || status == "FROZEN" || status == "DISABLE" {
				return object, fmt.Errorf("default workspace reached terminal state %q", status)
			}
		}
		if time.Now().After(deadline) {
			if lastErr != nil {
				return nil, lastErr
			}
			return object, fmt.Errorf("timeout waiting for default workspace status %q", target)
		}
		wait()
	}
}

// defaultWorkspaceState maps a GetDefaultWorkspace response onto the resource
// model. Empty strings become null, matching the IMS convention so absent
// computed fields read back as absent rather than "".
func defaultWorkspaceState(ctx context.Context, object map[string]interface{}, client *connectivity.AliyunClient) defaultWorkspaceResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	owner, _ := object["Owner"].(map[string]interface{})
	ownerStr := func(key string) types.String {
		if owner != nil {
			if v, ok := owner[key].(string); ok && v != "" {
				return types.StringValue(v)
			}
		}
		return types.StringNull()
	}

	workspaceId := str("WorkspaceId")

	var envList types.List
	if raw, ok := object["EnvTypes"].([]interface{}); ok {
		envs := make([]string, 0, len(raw))
		for _, v := range raw {
			if s, ok := v.(string); ok {
				envs = append(envs, s)
			}
		}
		lv, _ := types.ListValueFrom(ctx, types.StringType, envs)
		envList = lv
	} else {
		envList = types.ListNull(types.StringType)
	}

	regionId := types.StringNull()
	if client != nil && client.RegionId != "" {
		regionId = types.StringValue(client.RegionId)
	}

	return defaultWorkspaceResourceModel{
		Id:              workspaceId,
		WorkspaceId:     workspaceId,
		Description:     str("Description"),
		EnvTypes:        envList,
		WorkspaceName:   str("WorkspaceName"),
		DisplayName:     str("DisplayName"),
		Creator:         str("Creator"),
		Status:          str("Status"),
		RegionId:        regionId,
		GmtCreateTime:   str("GmtCreateTime"),
		GmtModifiedTime: str("GmtModifiedTime"),
		OwnerUserId:     ownerStr("UserId"),
		OwnerUserName:   ownerStr("UserName"),
		OwnerUserKp:     ownerStr("UserKp"),
	}
}
