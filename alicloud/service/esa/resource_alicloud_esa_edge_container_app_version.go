package esa

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	edgeContainerAppVersionTypeName = "alicloud_esa_edge_container_app_version"
	// edgeContainerAppVersionTimeout replaces the SDKv2 resource-level default
	// timeouts (5 minutes each); terraform-plugin-framework does not support
	// user-configurable resource timeouts yet.
	edgeContainerAppVersionTimeout = 5 * time.Minute
)

var (
	_ resource.Resource                = &edgeContainerAppVersionResource{}
	_ resource.ResourceWithImportState = &edgeContainerAppVersionResource{}
)

func NewEdgeContainerAppVersionResource() resource.Resource {
	return &edgeContainerAppVersionResource{}
}

type edgeContainerAppVersionResource struct {
	fwadapt.ResourceBase
}

// edgeContainerAppVersionResourceModel is the full state model. The create
// fields (app_id, name, remarks, containers) are replacement-only because the
// version is immutable once created; the publish fields drive the
// PublishEdgeContainerAppVersion call on update.
type edgeContainerAppVersionResourceModel struct {
	Id              types.String                            `tfsdk:"id"`
	AppId           types.String                            `tfsdk:"app_id"`
	VersionId       types.String                            `tfsdk:"version_id"`
	Name            types.String                            `tfsdk:"name"`
	Remarks         types.String                            `tfsdk:"remarks"`
	Containers      []edgeContainerAppVersionContainerModel `tfsdk:"containers"`
	PublishEnv      types.String                            `tfsdk:"publish_env"`
	PublishType     types.String                            `tfsdk:"publish_type"`
	Percentage      types.Int64                             `tfsdk:"percentage"`
	Regions         types.List                              `tfsdk:"regions"`
	FullRelease     types.Bool                              `tfsdk:"full_release"`
	StartTime       types.String                            `tfsdk:"start_time"`
	Status          types.String                            `tfsdk:"status"`
	CreateTime      types.String                            `tfsdk:"create_time"`
	UpdateTime      types.String                            `tfsdk:"update_time"`
	PublishTime     types.String                            `tfsdk:"publish_time"`
	LastPublishTime types.String                            `tfsdk:"last_publish_time"`
}

type edgeContainerAppVersionContainerModel struct {
	Name         types.String                              `tfsdk:"name"`
	Image        types.String                              `tfsdk:"image"`
	IsACRImage   types.Bool                                `tfsdk:"is_acr_image"`
	ACRImageInfo *edgeContainerAppVersionACRImageInfoModel `tfsdk:"acr_image_info"`
	Spec         types.String                              `tfsdk:"spec"`
	Command      types.String                              `tfsdk:"command"`
	Args         types.String                              `tfsdk:"args"`
	Storage      types.String                              `tfsdk:"storage"`
	EnvVariables types.String                              `tfsdk:"env_variables"`
	PostStart    types.String                              `tfsdk:"post_start"`
	PreStop      types.String                              `tfsdk:"pre_stop"`
	ProbeType    types.String                              `tfsdk:"probe_type"`
	ProbeContent *edgeContainerAppVersionProbeContentModel `tfsdk:"probe_content"`
}

type edgeContainerAppVersionACRImageInfoModel struct {
	IsEnterpriseRegistry types.Bool   `tfsdk:"is_enterprise_registry"`
	Domain               types.String `tfsdk:"domain"`
	InstanceId           types.String `tfsdk:"instance_id"`
	RegionId             types.String `tfsdk:"region_id"`
	RepoId               types.String `tfsdk:"repo_id"`
	RepoName             types.String `tfsdk:"repo_name"`
	RepoNamespace        types.String `tfsdk:"repo_namespace"`
	Tag                  types.String `tfsdk:"tag"`
	TagUrl               types.String `tfsdk:"tag_url"`
}

type edgeContainerAppVersionProbeContentModel struct {
	InitialDelaySeconds types.Int64  `tfsdk:"initial_delay_seconds"`
	TimeoutSeconds      types.Int64  `tfsdk:"timeout_seconds"`
	PeriodSeconds       types.Int64  `tfsdk:"period_seconds"`
	SuccessThreshold    types.Int64  `tfsdk:"success_threshold"`
	FailureThreshold    types.Int64  `tfsdk:"failure_threshold"`
	Command             types.String `tfsdk:"command"`
	Port                types.Int64  `tfsdk:"port"`
	Host                types.String `tfsdk:"host"`
	Path                types.String `tfsdk:"path"`
	Scheme              types.String `tfsdk:"scheme"`
	HttpHeaders         types.String `tfsdk:"http_headers"`
}

func (r *edgeContainerAppVersionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = edgeContainerAppVersionTypeName
}

func (r *edgeContainerAppVersionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"app_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"version_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"remarks": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"containers": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
						"image": schema.StringAttribute{
							Required: true,
						},
						"is_acr_image": schema.BoolAttribute{
							Required: true,
						},
						"acr_image_info": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"is_enterprise_registry": schema.BoolAttribute{
									Optional: true,
								},
								"domain": schema.StringAttribute{
									Optional: true,
								},
								"instance_id": schema.StringAttribute{
									Optional: true,
								},
								"region_id": schema.StringAttribute{
									Optional: true,
								},
								"repo_id": schema.StringAttribute{
									Optional: true,
								},
								"repo_name": schema.StringAttribute{
									Optional: true,
								},
								"repo_namespace": schema.StringAttribute{
									Optional: true,
								},
								"tag": schema.StringAttribute{
									Optional: true,
								},
								"tag_url": schema.StringAttribute{
									Optional: true,
								},
							},
						},
						"spec": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("1C2G", "2C4G", "2C8G", "4C8G", "4C16G", "8C16G", "8C32G"),
							},
						},
						"command": schema.StringAttribute{
							Optional: true,
						},
						"args": schema.StringAttribute{
							Optional: true,
						},
						"storage": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("0.5G", "10G", "20G", "30G"),
							},
						},
						"env_variables": schema.StringAttribute{
							Optional: true,
						},
						"post_start": schema.StringAttribute{
							Optional: true,
						},
						"pre_stop": schema.StringAttribute{
							Optional: true,
						},
						"probe_type": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("exec", "tcpSocket", "httpGet"),
							},
						},
						"probe_content": schema.SingleNestedAttribute{
							Required: true,
							Attributes: map[string]schema.Attribute{
								"initial_delay_seconds": schema.Int64Attribute{
									Optional: true,
								},
								"timeout_seconds": schema.Int64Attribute{
									Optional: true,
								},
								"period_seconds": schema.Int64Attribute{
									Optional: true,
								},
								"success_threshold": schema.Int64Attribute{
									Optional: true,
								},
								"failure_threshold": schema.Int64Attribute{
									Optional: true,
								},
								"command": schema.StringAttribute{
									Optional: true,
								},
								"port": schema.Int64Attribute{
									Optional: true,
								},
								"host": schema.StringAttribute{
									Optional: true,
								},
								"path": schema.StringAttribute{
									Optional: true,
								},
								"scheme": schema.StringAttribute{
									Optional: true,
								},
								"http_headers": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"publish_env": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("prod", "staging"),
				},
			},
			"publish_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("percentage", "region"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"percentage": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.Between(1, 100),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"regions": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"full_release": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"start_time": schema.StringAttribute{
				Optional: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"update_time": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"publish_time": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_publish_time": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *edgeContainerAppVersionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan edgeContainerAppVersionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	containersJson, err := edgeContainerAppVersionContainersToJSON(plan.Containers)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: serializing containers", edgeContainerAppVersionTypeName),
			err.Error(),
		)
		return
	}

	body := map[string]interface{}{
		"AppId":      plan.AppId.ValueString(),
		"Name":       plan.Name.ValueString(),
		"Containers": containersJson,
		"RegionId":   r.Meta().RegionId,
	}
	if !plan.Remarks.IsNull() && plan.Remarks.ValueString() != "" {
		body["Remarks"] = plan.Remarks.ValueString()
	}

	action := "CreateEdgeContainerAppVersion"
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, edgeContainerAppVersionTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = r.Meta().RpcPost("ESA", "2024-09-10", action, nil, body, true)
		if callErr != nil {
			if helper.NeedRetry(callErr) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	logs.AddDebug(action, response, body)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateEdgeContainerAppVersion", edgeContainerAppVersionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, edgeContainerAppVersionTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	versionIdRaw, ok := response["VersionId"]
	if !ok || versionIdRaw == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading VersionId", edgeContainerAppVersionTypeName),
			errs.WrapErrorf(fmt.Errorf("VersionId is missing in the CreateEdgeContainerAppVersion response"), errs.FailedGetAttributeMsg, edgeContainerAppVersionTypeName, "$.VersionId", response).Error(),
		)
		return
	}
	versionId := fmt.Sprint(versionIdRaw)
	id := edgeContainerAppVersionJoinId(plan.AppId.ValueString(), versionId)

	// Wait for the version to become "created" before publishing or reading
	// back, so a "creating" or "failed" version is not exposed as ready.
	stateConf := helper.BuildStateConf(
		[]string{"creating"},
		[]string{"created"},
		edgeContainerAppVersionTimeout,
		5*time.Second,
		NewService(r.Meta()).EsaEdgeContainerAppVersionStateRefreshFunc(id, "$.Status", []string{"failed"}),
	)
	if _, waitErr := stateConf.WaitForState(); waitErr != nil {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for status to become created", edgeContainerAppVersionTypeName),
			errs.WrapErrorf(waitErr, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	// Publish the version when publish_env is configured. The publish API is
	// part of the version's lifecycle (CloudSpec update operate), so the first
	// publish happens here instead of requiring a separate apply.
	if publishEnv := plan.PublishEnv.ValueString(); publishEnv != "" && !plan.PublishEnv.IsNull() {
		if pubErr := NewService(r.Meta()).PublishEsaEdgeContainerAppVersion(id, edgeContainerAppVersionPublishRequest(&plan)); pubErr != nil {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: calling PublishEdgeContainerAppVersion", edgeContainerAppVersionTypeName),
				errs.WrapErrorf(pubErr, errs.DefaultErrorMsg, id, "PublishEdgeContainerAppVersion", errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
	}

	object, readErr := NewService(r.Meta()).DescribeEsaEdgeContainerAppVersionContext(ctx, id)
	if readErr != nil {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after create: calling GetEdgeContainerAppVersion", edgeContainerAppVersionTypeName),
			errs.WrapErrorf(readErr, errs.DefaultErrorMsg, id, "GetEdgeContainerAppVersion", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	model, mapErr := edgeContainerAppVersionState(object, &plan)
	if mapErr != nil {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after create: mapping GetEdgeContainerAppVersion response", edgeContainerAppVersionTypeName),
			mapErr.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *edgeContainerAppVersionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state edgeContainerAppVersionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	if id == "" {
		return
	}

	object, err := NewService(r.Meta()).DescribeEsaEdgeContainerAppVersionContext(ctx, id)
	if err != nil {
		if errs.IsProviderNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetEdgeContainerAppVersion", edgeContainerAppVersionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetEdgeContainerAppVersion", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	model, mapErr := edgeContainerAppVersionState(object, &state)
	if mapErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: mapping GetEdgeContainerAppVersion response", edgeContainerAppVersionTypeName),
			mapErr.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

// Update handles changes to the publish parameters. Every create field is
// replacement-only, so reaching this method means only publish_env,
// publish_type, percentage, regions, full_release or start_time changed; the
// new values are published via PublishEdgeContainerAppVersion.
func (r *edgeContainerAppVersionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan edgeContainerAppVersionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state edgeContainerAppVersionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()

	if publishEnv := plan.PublishEnv.ValueString(); publishEnv != "" && !plan.PublishEnv.IsNull() {
		if err := NewService(r.Meta()).PublishEsaEdgeContainerAppVersion(id, edgeContainerAppVersionPublishRequest(&plan)); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling PublishEdgeContainerAppVersion", edgeContainerAppVersionTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "PublishEdgeContainerAppVersion", errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
	}

	object, err := NewService(r.Meta()).DescribeEsaEdgeContainerAppVersionContext(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after update: calling GetEdgeContainerAppVersion", edgeContainerAppVersionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetEdgeContainerAppVersion", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	model, mapErr := edgeContainerAppVersionState(object, &plan)
	if mapErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after update: mapping GetEdgeContainerAppVersion response", edgeContainerAppVersionTypeName),
			mapErr.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *edgeContainerAppVersionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state edgeContainerAppVersionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	if id == "" {
		return
	}
	if err := NewService(r.Meta()).DeleteEsaEdgeContainerAppVersionContext(ctx, id); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteEdgeContainerAppVersion", edgeContainerAppVersionTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DeleteEdgeContainerAppVersion", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
}

func (r *edgeContainerAppVersionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// edgeContainerAppVersionContainersToJSON serializes the configured containers
// block into the JSON array string expected by the CreateEdgeContainerAppVersion
// Containers formData parameter.
func edgeContainerAppVersionContainersToJSON(containers []edgeContainerAppVersionContainerModel) (string, error) {
	items := make([]map[string]interface{}, 0, len(containers))
	for _, c := range containers {
		item := map[string]interface{}{
			"Name":       c.Name.ValueString(),
			"Image":      c.Image.ValueString(),
			"IsACRImage": c.IsACRImage.ValueBool(),
			"Spec":       c.Spec.ValueString(),
			"Storage":    c.Storage.ValueString(),
			"ProbeType":  c.ProbeType.ValueString(),
		}
		if !c.Command.IsNull() && c.Command.ValueString() != "" {
			item["Command"] = c.Command.ValueString()
		}
		if !c.Args.IsNull() && c.Args.ValueString() != "" {
			item["Args"] = c.Args.ValueString()
		}
		if !c.EnvVariables.IsNull() && c.EnvVariables.ValueString() != "" {
			item["EnvVariables"] = c.EnvVariables.ValueString()
		}
		if !c.PostStart.IsNull() && c.PostStart.ValueString() != "" {
			item["PostStart"] = c.PostStart.ValueString()
		}
		if !c.PreStop.IsNull() && c.PreStop.ValueString() != "" {
			item["PreStop"] = c.PreStop.ValueString()
		}
		if c.ACRImageInfo != nil {
			item["ACRImageInfo"] = edgeContainerAppVersionACRImageInfoToMap(c.ACRImageInfo)
		}
		if c.ProbeContent != nil {
			item["ProbeContent"] = edgeContainerAppVersionProbeContentToMap(c.ProbeContent)
		}
		items = append(items, item)
	}
	b, err := json.Marshal(items)
	if err != nil {
		return "", fmt.Errorf("marshaling containers: %w", err)
	}
	return string(b), nil
}

func edgeContainerAppVersionACRImageInfoToMap(m *edgeContainerAppVersionACRImageInfoModel) map[string]interface{} {
	out := map[string]interface{}{}
	if !m.IsEnterpriseRegistry.IsNull() {
		out["IsEnterpriseRegistry"] = m.IsEnterpriseRegistry.ValueBool()
	}
	if !m.Domain.IsNull() && m.Domain.ValueString() != "" {
		out["Domain"] = m.Domain.ValueString()
	}
	if !m.InstanceId.IsNull() && m.InstanceId.ValueString() != "" {
		out["InstanceId"] = m.InstanceId.ValueString()
	}
	if !m.RegionId.IsNull() && m.RegionId.ValueString() != "" {
		out["RegionId"] = m.RegionId.ValueString()
	}
	if !m.RepoId.IsNull() && m.RepoId.ValueString() != "" {
		out["RepoId"] = m.RepoId.ValueString()
	}
	if !m.RepoName.IsNull() && m.RepoName.ValueString() != "" {
		out["RepoName"] = m.RepoName.ValueString()
	}
	if !m.RepoNamespace.IsNull() && m.RepoNamespace.ValueString() != "" {
		out["RepoNamespace"] = m.RepoNamespace.ValueString()
	}
	if !m.Tag.IsNull() && m.Tag.ValueString() != "" {
		out["Tag"] = m.Tag.ValueString()
	}
	if !m.TagUrl.IsNull() && m.TagUrl.ValueString() != "" {
		out["TagUrl"] = m.TagUrl.ValueString()
	}
	return out
}

func edgeContainerAppVersionProbeContentToMap(m *edgeContainerAppVersionProbeContentModel) map[string]interface{} {
	out := map[string]interface{}{}
	if !m.InitialDelaySeconds.IsNull() {
		out["InitialDelaySeconds"] = m.InitialDelaySeconds.ValueInt64()
	}
	if !m.TimeoutSeconds.IsNull() {
		out["TimeoutSeconds"] = m.TimeoutSeconds.ValueInt64()
	}
	if !m.PeriodSeconds.IsNull() {
		out["PeriodSeconds"] = m.PeriodSeconds.ValueInt64()
	}
	if !m.SuccessThreshold.IsNull() {
		out["SuccessThreshold"] = m.SuccessThreshold.ValueInt64()
	}
	if !m.FailureThreshold.IsNull() {
		out["FailureThreshold"] = m.FailureThreshold.ValueInt64()
	}
	if !m.Command.IsNull() && m.Command.ValueString() != "" {
		out["Command"] = m.Command.ValueString()
	}
	if !m.Port.IsNull() {
		out["Port"] = m.Port.ValueInt64()
	}
	if !m.Host.IsNull() && m.Host.ValueString() != "" {
		out["Host"] = m.Host.ValueString()
	}
	if !m.Path.IsNull() && m.Path.ValueString() != "" {
		out["Path"] = m.Path.ValueString()
	}
	if !m.Scheme.IsNull() && m.Scheme.ValueString() != "" {
		out["Scheme"] = m.Scheme.ValueString()
	}
	if !m.HttpHeaders.IsNull() && m.HttpHeaders.ValueString() != "" {
		out["HttpHeaders"] = m.HttpHeaders.ValueString()
	}
	return out
}

// edgeContainerAppVersionPublishRequest builds the publish parameter map from
// the planned publish fields. Only non-empty/configured values are forwarded.
func edgeContainerAppVersionPublishRequest(plan *edgeContainerAppVersionResourceModel) map[string]interface{} {
	req := map[string]interface{}{
		"PublishEnv": plan.PublishEnv.ValueString(),
	}
	if !plan.Percentage.IsNull() && !plan.Percentage.IsUnknown() {
		req["Percentage"] = plan.Percentage.ValueInt64()
	}
	if !plan.PublishType.IsNull() && !plan.PublishType.IsUnknown() && plan.PublishType.ValueString() != "" {
		req["PublishType"] = plan.PublishType.ValueString()
	}
	if !plan.StartTime.IsNull() && plan.StartTime.ValueString() != "" {
		req["StartTime"] = plan.StartTime.ValueString()
	}
	if !plan.FullRelease.IsNull() && !plan.FullRelease.IsUnknown() {
		req["FullRelease"] = plan.FullRelease.ValueBool()
	}
	if !plan.Regions.IsNull() && len(plan.Regions.Elements()) > 0 {
		regions := make([]string, 0, len(plan.Regions.Elements()))
		for _, el := range plan.Regions.Elements() {
			if s, ok := el.(types.String); ok && !s.IsNull() {
				regions = append(regions, s.ValueString())
			}
		}
		if len(regions) > 0 {
			// ESA RPC encodes array parameters as JSON strings, matching the
			// Containers formData convention used by sibling resources.
			b, err := json.Marshal(regions)
			if err != nil {
				// fall back to the raw slice; the publish API will surface a
				// clear error if the encoding is unacceptable.
				req["Regions"] = regions
			} else {
				req["Regions"] = string(b)
			}
		}
	}
	return req
}

// edgeContainerAppVersionState maps the GetEdgeContainerAppVersion response
// (the Version object) into the state model, preserving the configured publish
// parameters from the existing state/plan since the API does not return them.
func edgeContainerAppVersionState(object map[string]interface{}, source *edgeContainerAppVersionResourceModel) (edgeContainerAppVersionResourceModel, error) {
	appId, versionId, err := edgeContainerAppVersionAppIdVersion(object, source)
	if err != nil {
		return edgeContainerAppVersionResourceModel{}, err
	}
	containers, err := edgeContainerAppVersionContainersFromResponse(object["Containers"])
	if err != nil {
		return edgeContainerAppVersionResourceModel{}, err
	}
	model := edgeContainerAppVersionResourceModel{
		Id:              types.StringValue(edgeContainerAppVersionJoinId(appId, versionId)),
		AppId:           types.StringValue(appId),
		VersionId:       types.StringValue(versionId),
		Name:            edgeContainerAppVersionString(object["Name"]),
		Remarks:         edgeContainerAppVersionString(object["Remarks"]),
		Containers:      containers,
		Status:          edgeContainerAppVersionString(object["Status"]),
		CreateTime:      edgeContainerAppVersionString(object["CreateTime"]),
		UpdateTime:      edgeContainerAppVersionString(object["UpdateTime"]),
		PublishTime:     edgeContainerAppVersionString(object["PublishTime"]),
		LastPublishTime: edgeContainerAppVersionString(object["LastPublishTime"]),
		// Publish parameters are configured, not server-returned; preserve them.
		PublishEnv:  source.PublishEnv,
		PublishType: source.PublishType,
		Percentage:  source.Percentage,
		Regions:     source.Regions,
		FullRelease: source.FullRelease,
		StartTime:   source.StartTime,
	}
	return model, nil
}

// edgeContainerAppVersionAppIdVersion resolves the app and version identity
// from the response, falling back to the configured values when the API omits
// them, so the identity stays stable across refreshes.
func edgeContainerAppVersionAppIdVersion(object map[string]interface{}, source *edgeContainerAppVersionResourceModel) (string, string, error) {
	appId := edgeContainerAppVersionString(object["AppId"]).ValueString()
	if appId == "" {
		appId = source.AppId.ValueString()
	}
	versionId := edgeContainerAppVersionString(object["VersionId"]).ValueString()
	if versionId == "" {
		versionId = source.VersionId.ValueString()
	}
	if appId == "" || versionId == "" {
		return "", "", fmt.Errorf("AppId or VersionId is missing in the GetEdgeContainerAppVersion response")
	}
	return appId, versionId, nil
}

func edgeContainerAppVersionContainersFromResponse(raw interface{}) ([]edgeContainerAppVersionContainerModel, error) {
	if raw == nil {
		return []edgeContainerAppVersionContainerModel{}, nil
	}
	list, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected Containers response type %T", raw)
	}
	out := make([]edgeContainerAppVersionContainerModel, 0, len(list))
	for i, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected Containers[%d] element type %T", i, item)
		}
		out = append(out, edgeContainerAppVersionContainerModel{
			Name:         edgeContainerAppVersionString(m["Name"]),
			Image:        edgeContainerAppVersionString(m["Image"]),
			IsACRImage:   edgeContainerAppVersionBool(m["IsACRImage"]),
			ACRImageInfo: edgeContainerAppVersionACRImageInfoFromResponse(m["ACRImageInfo"]),
			Spec:         edgeContainerAppVersionString(m["Spec"]),
			Command:      edgeContainerAppVersionString(m["Command"]),
			Args:         edgeContainerAppVersionString(m["Args"]),
			Storage:      edgeContainerAppVersionString(m["Storage"]),
			EnvVariables: edgeContainerAppVersionString(m["EnvVariables"]),
			PostStart:    edgeContainerAppVersionString(m["PostStart"]),
			PreStop:      edgeContainerAppVersionString(m["PreStop"]),
			ProbeType:    edgeContainerAppVersionString(m["ProbeType"]),
			ProbeContent: edgeContainerAppVersionProbeContentFromResponse(m["ProbeContent"]),
		})
	}
	return out, nil
}

func edgeContainerAppVersionACRImageInfoFromResponse(raw interface{}) *edgeContainerAppVersionACRImageInfoModel {
	if raw == nil {
		return nil
	}
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil
	}
	return &edgeContainerAppVersionACRImageInfoModel{
		IsEnterpriseRegistry: edgeContainerAppVersionBool(m["IsEnterpriseRegistry"]),
		Domain:               edgeContainerAppVersionString(m["Domain"]),
		InstanceId:           edgeContainerAppVersionString(m["InstanceId"]),
		RegionId:             edgeContainerAppVersionString(m["RegionId"]),
		RepoId:               edgeContainerAppVersionString(m["RepoId"]),
		RepoName:             edgeContainerAppVersionString(m["RepoName"]),
		RepoNamespace:        edgeContainerAppVersionString(m["RepoNamespace"]),
		Tag:                  edgeContainerAppVersionString(m["Tag"]),
		TagUrl:               edgeContainerAppVersionString(m["TagUrl"]),
	}
}

func edgeContainerAppVersionProbeContentFromResponse(raw interface{}) *edgeContainerAppVersionProbeContentModel {
	if raw == nil {
		return nil
	}
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil
	}
	return &edgeContainerAppVersionProbeContentModel{
		InitialDelaySeconds: edgeContainerAppVersionInt64(m["InitialDelaySeconds"]),
		TimeoutSeconds:      edgeContainerAppVersionInt64(m["TimeoutSeconds"]),
		PeriodSeconds:       edgeContainerAppVersionInt64(m["PeriodSeconds"]),
		SuccessThreshold:    edgeContainerAppVersionInt64(m["SuccessThreshold"]),
		FailureThreshold:    edgeContainerAppVersionInt64(m["FailureThreshold"]),
		Command:             edgeContainerAppVersionString(m["Command"]),
		Port:                edgeContainerAppVersionInt64(m["Port"]),
		Host:                edgeContainerAppVersionString(m["Host"]),
		Path:                edgeContainerAppVersionString(m["Path"]),
		Scheme:              edgeContainerAppVersionString(m["Scheme"]),
		HttpHeaders:         edgeContainerAppVersionString(m["HttpHeaders"]),
	}
}

func edgeContainerAppVersionString(raw interface{}) types.String {
	if v, ok := raw.(string); ok && v != "" {
		return types.StringValue(v)
	}
	return types.StringNull()
}

func edgeContainerAppVersionBool(raw interface{}) types.Bool {
	if v, ok := raw.(bool); ok {
		return types.BoolValue(v)
	}
	return types.BoolNull()
}

func edgeContainerAppVersionInt64(raw interface{}) types.Int64 {
	switch v := raw.(type) {
	case float64:
		return types.Int64Value(int64(v))
	case int:
		return types.Int64Value(int64(v))
	case int64:
		return types.Int64Value(v)
	case json.Number:
		if i, err := v.Int64(); err == nil {
			return types.Int64Value(i)
		}
	}
	return types.Int64Null()
}
