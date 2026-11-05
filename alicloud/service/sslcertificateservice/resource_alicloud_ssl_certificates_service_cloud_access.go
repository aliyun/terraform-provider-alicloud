package sslcertificateservice

import (
	"context"
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const cloudAccessTypeName = "alicloud_ssl_certificates_service_cloud_access"

// cloudAccessTimeout replaces the SDKv2 resource-level default timeouts
// (Create/Read/Delete each 5 minutes); terraform-plugin-framework does not
// support user-configurable resource timeouts yet.
const cloudAccessTimeout = 5 * time.Minute

var (
	_ resource.Resource                = &cloudAccessResource{}
	_ resource.ResourceWithImportState = &cloudAccessResource{}
)

func NewCloudAccessResource() resource.Resource {
	return &cloudAccessResource{}
}

type cloudAccessResource struct {
	fwadapt.ResourceBase
}

type cloudAccessResourceModel struct {
	Id            types.String `tfsdk:"id"`
	AccessId      types.String `tfsdk:"access_id"`
	CloudName     types.String `tfsdk:"cloud_name"`
	SecretId      types.String `tfsdk:"secret_id"`
	SecretKey     types.String `tfsdk:"secret_key"`
	ServiceStatus types.String `tfsdk:"service_status"`
}

func (r *cloudAccessResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = cloudAccessTypeName
}

func (r *cloudAccessResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"access_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cloud_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"secret_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"secret_key": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// service_status transitions server-side, so the planned state
			// value could never survive an apply.
			"service_status": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *cloudAccessResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan cloudAccessResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := "AddCloudAccess"
	request := map[string]interface{}{
		"RegionId":  r.Meta().RegionId,
		"CloudName": plan.CloudName.ValueString(),
		"SecretId":  plan.SecretId.ValueString(),
		"SecretKey": plan.SecretKey.ValueString(),
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(cloudAccessTimeout, func() *retry.RetryError {
		response, err = r.Meta().RpcPost("cas", "2020-04-07", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling AddCloudAccess", cloudAccessTypeName),
			errs.WrapErrorf(err, "Creating %s: calling AddCloudAccess", cloudAccessTypeName).Error(),
		)
		return
	}

	// AddCloudAccess does not return AccessId; read it back via ListCloudAccess.
	// The newly created entry may not be immediately visible, so retry.
	service := NewCasService(r.Meta())
	var entry map[string]interface{}
	readWait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(cloudAccessTimeout, func() *retry.RetryError {
		entries, listErr := service.ListCloudAccess(plan.CloudName.ValueString(), plan.SecretId.ValueString())
		if listErr != nil {
			return retry.NonRetryableError(listErr)
		}
		if len(entries) == 0 {
			readWait()
			return retry.RetryableError(fmt.Errorf("CloudAccess entry not yet visible after AddCloudAccess"))
		}
		entry = entries[0]
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading back AccessId via ListCloudAccess", cloudAccessTypeName),
			errs.WrapErrorf(err, "Creating %s: reading back AccessId via ListCloudAccess", cloudAccessTypeName).Error(),
		)
		return
	}

	state := cloudAccessState(entry)
	state.SecretKey = plan.SecretKey
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *cloudAccessResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state cloudAccessResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewCasService(r.Meta())
	entry, err := service.DescribeCloudAccess(state.Id.ValueString(), state.CloudName.ValueString(), state.SecretId.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListCloudAccess", cloudAccessTypeName),
			errs.WrapErrorf(err, "Reading %s: calling ListCloudAccess", cloudAccessTypeName).Error(),
		)
		return
	}

	model := cloudAccessState(entry)
	// SecretKey is write-only; ListCloudAccess never returns it, so refresh
	// must keep the previously configured value.
	model.SecretKey = state.SecretKey
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *cloudAccessResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state cloudAccessResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// All configurable attributes carry RequiresReplace, so any change triggers
	// destroy+create. Update is only reached for computed-only refresh; read
	// the current state and overlay the write-only SecretKey.
	service := NewCasService(r.Meta())
	entry, err := service.DescribeCloudAccess(state.Id.ValueString(), state.CloudName.ValueString(), state.SecretId.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling ListCloudAccess", cloudAccessTypeName),
			errs.WrapErrorf(err, "Updating %s: calling ListCloudAccess", cloudAccessTypeName).Error(),
		)
		return
	}

	model := cloudAccessState(entry)
	model.SecretKey = plan.SecretKey
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *cloudAccessResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state cloudAccessResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := "DeleteCloudAccess"
	request := map[string]interface{}{
		"RegionId": r.Meta().RegionId,
		"AccessId": state.Id.ValueString(),
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(cloudAccessTimeout, func() *retry.RetryError {
		response, err = r.Meta().RpcPost("cas", "2020-04-07", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteCloudAccess", cloudAccessTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteCloudAccess", cloudAccessTypeName).Error(),
		)
		return
	}
}

func (r *cloudAccessResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// cloudAccessState maps a ListCloudAccess entry onto the resource model.
// SecretKey is deliberately left null: it is write-only and every caller
// overlays the plan/state value before writing state.
func cloudAccessState(object map[string]interface{}) cloudAccessResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	accessId := formatAccessId(object["AccessId"])
	var idField types.String
	if accessId != "" {
		idField = types.StringValue(accessId)
	} else {
		idField = types.StringNull()
	}
	return cloudAccessResourceModel{
		Id:            idField,
		AccessId:      idField,
		CloudName:     str("CloudName"),
		SecretId:      str("SecretId"),
		ServiceStatus: str("ServiceStatus"),
	}
}
