package eas

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const paiEasResourceTypeName = "alicloud_pai_eas_resource"

var (
	_ resource.Resource                = &paiEasResource{}
	_ resource.ResourceWithImportState = &paiEasResource{}
)

func NewPaiEasResourceResource() resource.Resource {
	return &paiEasResource{}
}

// paiEasResource implements the alicloud_pai_eas_resource resource, an EAS
// (PAI) resource group backed by the eas/2021-07-01 Roa API.
type paiEasResource struct {
	fwadapt.ResourceBase
}

// paiEasResourceModel mirrors the alicloud_pai_eas_resource schema. The
// self-managed resource options (external_cluster_id, role_name,
// node_match_labels) are flattened to top-level attributes because the provider
// is served through the protocol v5 mux, which cannot represent nested
// attributes; the Create request still sends them under the nested
// SelfManagedResourceOptions object, and they are create-only (ForceNew) since
// DescribeResource does not return them.
type paiEasResourceModel struct {
	Id                types.String `tfsdk:"id"`
	ResourceName      types.String `tfsdk:"resource_name"`
	RegionId          types.String `tfsdk:"region_id"`
	Status            types.String `tfsdk:"status"`
	CreateTime        types.String `tfsdk:"create_time"`
	EcsInstanceCount  types.Int64  `tfsdk:"ecs_instance_count"`
	EcsInstanceType   types.String `tfsdk:"ecs_instance_type"`
	SystemDiskSize    types.Int64  `tfsdk:"system_disk_size"`
	ChargeType        types.String `tfsdk:"charge_type"`
	AutoRenewal       types.Bool   `tfsdk:"auto_renewal"`
	Zone              types.String `tfsdk:"zone"`
	ResourceType      types.String `tfsdk:"resource_type"`
	UsageMode         types.String `tfsdk:"usage_mode"`
	Labels            types.Map    `tfsdk:"labels"`
	ExternalClusterId types.String `tfsdk:"external_cluster_id"`
	RoleName          types.String `tfsdk:"role_name"`
	NodeMatchLabels   types.Map    `tfsdk:"node_match_labels"`
}

func (r *paiEasResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = paiEasResourceTypeName
}

func (r *paiEasResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PAI EAS Resource Group resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the EAS resource group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the resource group. The only attribute that can be updated after creation.",
			},
			"region_id": schema.StringAttribute{
				Computed:    true,
				Description: "The region (cluster) ID of the resource group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the resource group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the resource group was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ecs_instance_count": schema.Int64Attribute{
				Optional:    true,
				Description: "The number of instances. Required when `resource_type` is `Dedicated`.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"ecs_instance_type": schema.StringAttribute{
				Optional:    true,
				Description: "The ECS instance type.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"system_disk_size": schema.Int64Attribute{
				Optional:    true,
				Description: "The system disk size in GB. Valid values: 200 to 2000. Defaults to 200.",
				Validators: []validator.Int64{
					int64validator.Between(200, 2000),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"charge_type": schema.StringAttribute{
				Optional:    true,
				Description: "The charge type. Valid values: `PrePaid`, `PostPaid`.",
				Validators: []validator.String{
					stringvalidator.OneOf("PrePaid", "PostPaid"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"auto_renewal": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to auto renew. Valid when `charge_type` is `PrePaid`.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"zone": schema.StringAttribute{
				Optional:    true,
				Description: "The zone ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"resource_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The resource type. Valid values: `Dedicated`, `SelfManaged`.",
				Validators: []validator.String{
					stringvalidator.OneOf("Dedicated", "SelfManaged"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"usage_mode": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The usage mode.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"labels": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "The labels of the resource group.",
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			// self_managed_resource_options is flattened into the three
			// top-level attributes below because the provider is served through
			// the protocol v5 mux, which cannot represent nested attributes.
			// The Create request still sends them under the nested
			// SelfManagedResourceOptions object.
			"external_cluster_id": schema.StringAttribute{
				Optional:    true,
				Description: "SelfManagedResourceOptions: the external cluster ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role_name": schema.StringAttribute{
				Optional:    true,
				Description: "SelfManagedResourceOptions: the role name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"node_match_labels": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "SelfManagedResourceOptions: the node match labels.",
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *paiEasResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan paiEasResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// EAS-09: when ResourceType is Dedicated, the API requires EcsInstanceCount,
	// EcsInstanceType and ChargeType to provision the backing machines. Reject
	// the plan early with actionable attribute-level errors instead of letting
	// CreateResource fail with a less actionable server-side message.
	if plan.ResourceType.ValueString() == "Dedicated" {
		if plan.EcsInstanceCount.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ecs_instance_count"),
				"Missing required attribute when resource_type is Dedicated",
				"ecs_instance_count is required when resource_type is Dedicated.",
			)
		}
		if plan.EcsInstanceType.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ecs_instance_type"),
				"Missing required attribute when resource_type is Dedicated",
				"ecs_instance_type is required when resource_type is Dedicated.",
			)
		}
		if plan.ChargeType.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("charge_type"),
				"Missing required attribute when resource_type is Dedicated",
				"charge_type is required when resource_type is Dedicated.",
			)
		}
		if resp.Diagnostics.HasError() {
			return
		}
	}

	client := r.Client()
	svc := NewEasService(client)

	body := map[string]interface{}{
		"ResourceName": plan.ResourceName.ValueString(),
	}
	resp.Diagnostics.Append(addCreateOnlyFields(ctx, &plan, body)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := "/api/v2/resources"
	response, err := roaPost(client, action, body, plan.ResourceName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateResource", paiEasResourceTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, plan.ResourceName.ValueString(), action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	resourceId, _ := response["ResourceId"].(string)
	if resourceId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateResource", paiEasResourceTypeName),
			fmt.Sprintf("the CreateResource response carried no ResourceId: %v", response),
		)
		return
	}

	// EAS-01: persist the created resource identity before any wait or read so
	// a transient failure during provisioning does not lose Terraform's tracking
	// of the resource; the user can then refresh or destroy instead of an
	// orphan the provider no longer manages.
	plan.Id = types.StringValue(resourceId)

	// Wait until the resource group reaches the ready state. The EAS API may
	// provision the backing instances asynchronously.
	if err := svc.WaitForResource(resourceId, []string{"ResourceReady"}, []string{"ResourceCreateError"}, 20*time.Minute); err != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for ResourceReady", paiEasResourceTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, resourceId, "WaitForResource", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	object, err := svc.DescribeResource(resourceId)
	if err != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeResource", paiEasResourceTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, resourceId, "DescribeResource", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	setComputedFromResponse(&plan, object)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *paiEasResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state paiEasResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc := NewEasService(r.Client())
	object, err := svc.DescribeResource(state.Id.ValueString())
	if err != nil {
		if isResourceNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeResource", paiEasResourceTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), "DescribeResource", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	setComputedFromResponse(&state, object)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *paiEasResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state paiEasResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()
	// Only resource_name is updatable; every other configurable attribute is
	// ForceNew and would have triggered a replace instead of reaching Update.
	body := map[string]interface{}{
		"ResourceName": plan.ResourceName.ValueString(),
	}
	action := fmt.Sprintf("/api/v2/resources/%s/%s", client.RegionId, state.Id.ValueString())
	// EAS-06: a retryable error from UpdateResource must be retried within a
	// bounded window rather than dropped as a warning; surfacing it as a warning
	// and returning leaves the update un-applied while telling Terraform success,
	// so the plan drift is never reconciled. Retry then fail with an error.
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	if err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, e := client.RoaPut(easProductCode, easAPIVersion, action, nil, nil, body, true)
		if e != nil {
			if helper.NeedRetry(e) {
				wait()
				return retry.RetryableError(e)
			}
			return retry.NonRetryableError(e)
		}
		return nil
	}); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateResource", paiEasResourceTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	svc := NewEasService(client)
	object, err := svc.DescribeResource(state.Id.ValueString())
	if err != nil {
		if isResourceNotFound(err) {
			// EAS-07: the resource disappeared after UpdateResource. Do not drop
			// state here: removing the resource mid-Update would hide the
			// disappearance and leave the prior plan's changes unapplied. Surface
			// the error so the next Read reconciles state against the server.
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: resource disappeared after UpdateResource", paiEasResourceTypeName),
				fmt.Sprintf("the EAS resource group %s was not found after update; the next apply will refresh state", state.Id.ValueString()),
			)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeResource", paiEasResourceTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), "DescribeResource", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	setComputedFromResponse(&plan, object)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *paiEasResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state paiEasResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()
	svc := NewEasService(client)
	resourceId := state.Id.ValueString()

	// EAS-02: a Dedicated resource group owns the machine instances EAS
	// provisioned at create time; DeleteResource refuses to delete a group that
	// still has instances, so destroy would fail and the instances would keep
	// billing. Release them first. SelfManaged groups hold machines that belong
	// to an external cluster and are not released here — deleting them would
	// cross the resource's ownership boundary.
	if state.ResourceType.ValueString() == "Dedicated" {
		instances, err := svc.ListResourceInstances(resourceId)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Deleting %s: listing instances to release", paiEasResourceTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, resourceId, "ListResourceInstances", errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
		var postPaid, prePaid []string
		for _, inst := range instances {
			instanceId, _ := inst["InstanceId"].(string)
			if instanceId == "" {
				continue
			}
			chargeType, _ := inst["ChargeType"].(string)
			if chargeType == "PrePaid" {
				prePaid = append(prePaid, instanceId)
			} else {
				postPaid = append(postPaid, instanceId)
			}
		}
		if len(postPaid) > 0 {
			if err := svc.DeleteResourceInstances(resourceId, strings.Join(postPaid, ",")); err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Deleting %s: releasing post-paid instances", paiEasResourceTypeName),
					errs.WrapErrorf(err, errs.DefaultErrorMsg, resourceId, "DeleteResourceInstances", errs.AlibabaCloudSdkGoERROR).Error(),
				)
				return
			}
			if err := svc.WaitForInstancesGone(resourceId, 20*time.Minute); err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Deleting %s: waiting for instances to release", paiEasResourceTypeName),
					errs.WrapErrorf(err, errs.DefaultErrorMsg, resourceId, "WaitForInstancesGone", errs.AlibabaCloudSdkGoERROR).Error(),
				)
				return
			}
		}
		// The EAS DeleteResourceInstances API only allows ordinary users to
		// delete post-paid instances; prepaid instances must be unsubscribed
		// through the prepaid退订流程, which is outside this resource's
		// authorized scope. Surface the gap so the user unsubscribes them before
		// destroying the resource group, instead of silently leaving them to keep
		// billing.
		if len(prePaid) > 0 {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Deleting %s: prepaid instances cannot be released here", paiEasResourceTypeName),
				fmt.Sprintf("Resource group %s still has %d prepaid instance(s): %s. The EAS DeleteResourceInstances API does not allow ordinary users to delete prepaid instances; unsubscribe them through the prepaid退订流程 before destroying the resource group.", resourceId, len(prePaid), strings.Join(prePaid, ",")),
			)
			return
		}
	}

	action := fmt.Sprintf("/api/v2/resources/%s/%s", client.RegionId, resourceId)
	if _, err := client.RoaDelete(easProductCode, easAPIVersion, action, nil, nil, nil, true); err != nil {
		if isResourceNotFound(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteResource", paiEasResourceTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, resourceId, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	// EAS-03: wait for the resource group to disappear and surface a failure
	// (not a warning) when it does not, so the resource stays under Terraform
	// management instead of being silently orphaned with a dangling state.
	if err := svc.WaitForResourceGone(resourceId, 20*time.Minute); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: waiting for deletion", paiEasResourceTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, resourceId, "WaitForResourceGone", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
}

func (r *paiEasResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// setComputedFromResponse fills the server-side (computed and read-back)
// attributes from a DescribeResource response. Create-only attributes that the
// API does not return (EcsInstance*, SystemDiskSize, ChargeType, AutoRenewal,
// Zone, Labels and the flattened self-managed resource options
// external_cluster_id/role_name/node_match_labels) are left untouched so the
// planned or prior-state value is preserved.
func setComputedFromResponse(m *paiEasResourceModel, obj map[string]interface{}) {
	str := func(key string) types.String {
		if v, ok := obj[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	m.ResourceName = str("ResourceName")
	m.RegionId = str("ClusterId")
	m.Status = str("Status")
	m.CreateTime = str("CreateTime")
	m.ResourceType = str("ResourceType")
	m.UsageMode = str("UsageMode")
}

// addCreateOnlyFields fills the create request body with every configured
// create-only attribute. Each field is omitted when null so the server
// default applies, matching the API contract.
func addCreateOnlyFields(ctx context.Context, plan *paiEasResourceModel, body map[string]interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	if !plan.EcsInstanceCount.IsNull() {
		body["EcsInstanceCount"] = plan.EcsInstanceCount.ValueInt64()
	}
	if !plan.EcsInstanceType.IsNull() {
		body["EcsInstanceType"] = plan.EcsInstanceType.ValueString()
	}
	if !plan.SystemDiskSize.IsNull() {
		body["SystemDiskSize"] = plan.SystemDiskSize.ValueInt64()
	}
	if !plan.ChargeType.IsNull() {
		body["ChargeType"] = plan.ChargeType.ValueString()
	}
	if !plan.AutoRenewal.IsNull() {
		body["AutoRenewal"] = plan.AutoRenewal.ValueBool()
	}
	if !plan.Zone.IsNull() {
		body["Zone"] = plan.Zone.ValueString()
	}
	if !plan.ResourceType.IsNull() {
		body["ResourceType"] = plan.ResourceType.ValueString()
	}
	if !plan.UsageMode.IsNull() {
		body["UsageMode"] = plan.UsageMode.ValueString()
	}
	if !plan.Labels.IsNull() {
		labels := map[string]string{}
		if labelDiags := plan.Labels.ElementsAs(ctx, &labels, false); labelDiags.HasError() {
			return labelDiags
		}
		// EAS-05: the Labels API contract is type=object with
		// additionalProperties of type string; send the map directly rather than
		// a JSON-encoded string, which the server would reject or mis-store.
		body["Labels"] = labels
	}
	// self_managed_resource_options is flattened to the three top-level
	// attributes external_cluster_id, role_name and node_match_labels. They are
	// sent under the nested SelfManagedResourceOptions object, and only included
	// when at least one sub-key is configured, so the server default applies
	// otherwise.
	if !plan.ExternalClusterId.IsNull() || !plan.RoleName.IsNull() || !plan.NodeMatchLabels.IsNull() {
		opt := map[string]interface{}{}
		if !plan.ExternalClusterId.IsNull() {
			opt["ExternalClusterId"] = plan.ExternalClusterId.ValueString()
		}
		if !plan.RoleName.IsNull() {
			opt["RoleName"] = plan.RoleName.ValueString()
		}
		if !plan.NodeMatchLabels.IsNull() {
			nodeMatchLabels := map[string]string{}
			if nmlDiags := plan.NodeMatchLabels.ElementsAs(ctx, &nodeMatchLabels, false); nmlDiags.HasError() {
				return nmlDiags
			}
			opt["NodeMatchLabels"] = nodeMatchLabels
		}
		body["SelfManagedResourceOptions"] = opt
	}
	return diags
}

// roaPost wraps the create RoaPost call with the throttling retry loop used by
// the SDKv2 eas resources.
func roaPost(client *connectivity.AliyunClient, action string, body map[string]interface{}, name string) (map[string]interface{}, error) {
	// EAS-14: reuse helper.IncrementalWait (capped at 20s with jitter) instead
	// of a local copy whose doubling never actually stops at a ceiling.
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	var response map[string]interface{}
	var err error
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaPost(easProductCode, easAPIVersion, action, nil, nil, body, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, name, action, errs.AlibabaCloudSdkGoERROR)
	}
	return response, nil
}
