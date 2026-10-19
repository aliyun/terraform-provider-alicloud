package cms

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const entityStoreTypeName = "alicloud_cms_entity_store"

// entityStoreTimeout mirrors the SDKv2 line's default Create/Delete timeout.
// Framework resources in this repo do not expose a user-configurable timeouts
// block, so the value is fixed here to preserve the original behavior.
const entityStoreTimeout = 5 * time.Minute

var (
	_ resource.Resource                = &entityStoreResource{}
	_ resource.ResourceWithImportState = &entityStoreResource{}
)

func NewEntityStoreResource() resource.Resource {
	return &entityStoreResource{}
}

type entityStoreResource struct {
	fwadapt.ResourceBase
}

type entityStoreResourceModel struct {
	Id            types.String `tfsdk:"id"`
	WorkspaceName types.String `tfsdk:"workspace_name"`
	RegionId      types.String `tfsdk:"region_id"`
}

func (r *entityStoreResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = entityStoreTypeName
}

func (r *entityStoreResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a Cms EntityStore resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. It is the name of the workspace that the EntityStore belongs to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workspace_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the workspace that the EntityStore belongs to. EntityStore is a singleton resource identified by the workspace name; any change forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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

func (r *entityStoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan entityStoreResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceName := plan.WorkspaceName.ValueString()
	action := fmt.Sprintf("/workspace/%s/entitystore", workspaceName)

	// CreateEntityStore carries no body parameter; the empty map matches the
	// ROA wire shape used by the gateway. The SDKv2 line retried on
	// throttling with NeedRetry + incrementalWait; helper/retry provides the
	// same Retry semantics without registering the sweep flag, so the test
	// binary stays compatible with terraform-plugin-testing.
	var response map[string]interface{}
	err := retry.Retry(entityStoreTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = r.Client().RoaPost("Cms", "2024-03-30", action, nil, nil, map[string]interface{}{}, true)
		if callErr != nil {
			if helper.NeedRetry(callErr) {
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateEntityStore", entityStoreTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateEntityStore", entityStoreTypeName).Error(),
		)
		return
	}

	// The identity is the workspace name, echoed back on the response.
	if wn, _ := response["workspaceName"].(string); wn != "" {
		plan.Id = types.StringValue(wn)
	} else {
		plan.Id = types.StringValue(workspaceName)
	}

	object, err := getEntityStore(r.Client(), plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetEntityStore", entityStoreTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetEntityStore", entityStoreTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateEntityStore", entityStoreTypeName),
			fmt.Sprintf("The EntityStore for workspace %s is not visible via GetEntityStore right after creation", workspaceName),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, entityStoreState(plan.Id.ValueString(), object))...)
}

func (r *entityStoreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state entityStoreResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := getEntityStore(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetEntityStore", entityStoreTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetEntityStore", entityStoreTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, entityStoreState(state.Id.ValueString(), object))...)
}

func (r *entityStoreResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// workspace_name is ForceNew and region_id is Computed, so Terraform never
	// reaches Update with an in-place change. The method still refreshes state
	// from the live object to satisfy the resource.Resource interface and to
	// keep drift correction consistent with Read.
	var state entityStoreResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := getEntityStore(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetEntityStore", entityStoreTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetEntityStore", entityStoreTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, entityStoreState(state.Id.ValueString(), object))...)
}

func (r *entityStoreResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state entityStoreResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := fmt.Sprintf("/workspace/%s/entitystore", state.Id.ValueString())
	err := retry.Retry(entityStoreTimeout, func() *retry.RetryError {
		_, callErr := r.Client().RoaDelete("Cms", "2024-03-30", action, nil, nil, nil, true)
		if callErr != nil {
			// The EntityStore or its workspace is already gone; treat as
			// success so destroy stays idempotent, matching the SDKv2 line.
			if isEntityStoreNotExist(callErr) {
				return nil
			}
			if helper.NeedRetry(callErr) {
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteEntityStore", entityStoreTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteEntityStore", entityStoreTypeName).Error(),
		)
		return
	}
}

func (r *entityStoreResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getEntityStore queries GetEntityStore by workspace name and returns the
// response object, or nil when the EntityStore (or its workspace) no longer
// exists.
func getEntityStore(client *connectivity.AliyunClient, workspaceName string) (map[string]interface{}, error) {
	action := fmt.Sprintf("/workspace/%s/entitystore", workspaceName)
	response, err := client.RoaGet("Cms", "2024-03-30", action, nil, nil, nil)
	if err != nil {
		if isEntityStoreNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return response, nil
}

// entityStoreState maps a GetEntityStore response onto the resource model.
// Empty strings become null so absent API values read back as absent, not "".
func entityStoreState(id string, object map[string]interface{}) entityStoreResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	m := entityStoreResourceModel{
		Id:            types.StringValue(id),
		WorkspaceName: str("workspaceName"),
		RegionId:      str("regionId"),
	}
	// Prefer the echoed identity when the API carries it.
	if wn := str("workspaceName"); !wn.IsNull() {
		m.Id = wn
	}
	return m
}

// isEntityStoreNotExist reports whether err is the EntityStore / workspace
// not-found error. The code arrives on a *tea.SDKError from the ROA channel;
// the substring fallback covers wrapped forms, matching the SDKv2 line's
// IsExpectedErrors behavior on the same codes.
func isEntityStoreNotExist(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil {
		switch tea.StringValue(e.Code) {
		case "EntityStoreNotExist", "WorkspaceNotAvailable":
			return true
		}
	}
	msg := err.Error()
	return strings.Contains(msg, "EntityStoreNotExist") || strings.Contains(msg, "WorkspaceNotAvailable")
}
