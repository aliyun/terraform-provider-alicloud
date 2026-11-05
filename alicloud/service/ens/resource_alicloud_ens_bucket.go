package ens

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const ensBucketTypeName = "alicloud_ens_bucket"

// Bounded retry durations for the ENS bucket lifecycle. The ENS control plane
// is eventually consistent: right after PutBucket the service may report
// NoSuchBucket on GetBucketInfo, and the APIs throttle under load, so each
// write/read is wrapped in an incremental retry. These match the ENS package's
// existing test setup helper durations.
const (
	ensBucketPutTimeout    = 2 * time.Minute
	ensBucketReadyTimeout  = 3 * time.Minute
	ensBucketAclTimeout    = 2 * time.Minute
	ensBucketDeleteTimeout = 1 * time.Minute
)

var (
	_ resource.Resource                = &ensBucketResource{}
	_ resource.ResourceWithImportState = &ensBucketResource{}
)

// NewEnsBucketResource is the native Framework factory registered through
// ServicePackage.Resources.
func NewEnsBucketResource() resource.Resource {
	return &ensBucketResource{}
}

// ensBucketResource is the native Framework implementation of the ENS bucket.
type ensBucketResource struct {
	fwadapt.ResourceBase
}

type ensBucketResourceModel struct {
	Id                types.String `tfsdk:"id"`
	BucketName        types.String `tfsdk:"bucket_name"`
	LogicalBucketType types.String `tfsdk:"logical_bucket_type"`
	Comment           types.String `tfsdk:"comment"`
	BucketAcl         types.String `tfsdk:"bucket_acl"`
	EnsRegionId       types.String `tfsdk:"ens_region_id"`
	DispatchScope     types.String `tfsdk:"dispatch_scope"`
	CreateTime        types.String `tfsdk:"create_time"`
	ModifyTime        types.String `tfsdk:"modify_time"`
}

func (r *ensBucketResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ensBucketTypeName
}

func (r *ensBucketResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an ENS Bucket resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource, same as `bucket_name`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bucket_name": schema.StringAttribute{
				Required: true,
				Description: "The name of the bucket. 3 to 50 characters in length, and can only contain lowercase letters, " +
					"digits and hyphens (-). It cannot start or end with a hyphen.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 50),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`),
						"must only contain lowercase letters, digits and hyphens (-), and must not start or end with a hyphen",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"logical_bucket_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "The logical bucket type. Valid values: `sink` (single-node, requires `ens_region_id`) and " +
					"`standard` (global, requires `dispatch_scope`). Defaults to `standard`.",
				Validators: []validator.String{
					stringvalidator.OneOf("sink", "standard"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"comment": schema.StringAttribute{
				Optional: true,
				Description: "The remark of the bucket. 0 to 128 characters in length. " +
					"Changing this field recreates the resource; the ENS API exposes no in-place update for the bucket comment.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"bucket_acl": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "The access control list of the bucket. Valid values: `private` (default), `public-read` and `public-read-write`. " +
					"Updatable in place via PutBucketAcl.",
				Validators: []validator.String{
					stringvalidator.OneOf("private", "public-read", "public-read-write"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ens_region_id": schema.StringAttribute{
				Optional: true,
				Description: "The node region id. Required when `logical_bucket_type` is `sink` (single-node bucket). " +
					"Write-only: not returned by GetBucketInfo, so the configured value is preserved in state.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dispatch_scope": schema.StringAttribute{
				Optional: true,
				Description: "The dispatch scope. Valid values: `domestic` and `oversea`. Required for a global (`standard`) bucket. " +
					"Write-only: not returned by GetBucketInfo, so the configured value is preserved in state.",
				Validators: []validator.String{
					stringvalidator.OneOf("domestic", "oversea"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the bucket was created, in RFC 3339 format (UTC).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modify_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the bucket was last modified, in RFC 3339 format (UTC).",
				// No UseStateForUnknown: PutBucketAcl bumps the timestamp, so the
				// planned state value can never survive an apply.
			},
		},
	}
}

func (r *ensBucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ensBucketResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"BucketName": plan.BucketName.ValueString(),
	}
	// logical_bucket_type is Optional+Computed: only forward it when the user set it.
	if !plan.LogicalBucketType.IsNull() && !plan.LogicalBucketType.IsUnknown() {
		request["LogicalBucketType"] = plan.LogicalBucketType.ValueString()
	}
	if !plan.Comment.IsNull() {
		request["Comment"] = plan.Comment.ValueString()
	}
	if !plan.BucketAcl.IsNull() && !plan.BucketAcl.IsUnknown() {
		request["BucketAcl"] = plan.BucketAcl.ValueString()
	}
	if !plan.EnsRegionId.IsNull() {
		request["EnsRegionId"] = plan.EnsRegionId.ValueString()
	}
	if !plan.DispatchScope.IsNull() {
		request["DispatchScope"] = plan.DispatchScope.ValueString()
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	if err := retry.Retry(ensBucketPutTimeout, func() *retry.RetryError {
		_, e := r.Meta().RpcPost("Ens", "2017-11-10", "PutBucket", nil, request, true)
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
			fmt.Sprintf("Creating %s: calling PutBucket", ensBucketTypeName),
			errs.WrapErrorf(err, "Creating %s: calling PutBucket", ensBucketTypeName).Error(),
		)
		return
	}

	// Persist the resource identity before the read-back: the bucket exists
	// remotely from here on, so a failing GetBucketInfo must not leave Terraform
	// without a resource to taint and destroy on the next apply.
	resp.Diagnostics.Append(resp.State.Set(ctx, ensBucketResourceModel{
		Id:         plan.BucketName,
		BucketName: plan.BucketName,
	})...)
	if resp.Diagnostics.HasError() {
		return
	}

	info, err := waitEnsBucketReady(r.Meta(), plan.BucketName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for GetBucketInfo", ensBucketTypeName),
			errs.WrapErrorf(err, "Creating %s: waiting for GetBucketInfo", ensBucketTypeName).Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, ensBucketState(info, plan))...)
}

func (r *ensBucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ensBucketResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// After `terraform import`, ImportStatePassthroughID only populates `id`,
	// so `bucket_name` is null. `id` always equals `bucket_name`, so fall back
	// to it as the GetBucketInfo lookup key.
	bucketName := state.BucketName.ValueString()
	if bucketName == "" {
		bucketName = state.Id.ValueString()
	}
	info, err := getEnsBucket(r.Meta(), bucketName)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetBucketInfo", ensBucketTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetBucketInfo", ensBucketTypeName).Error(),
		)
		return
	}
	if info == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, ensBucketState(info, state))...)
}

func (r *ensBucketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ensBucketResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Only bucket_acl is updatable in place; every other mutable attribute is
	// ForceNew, so reaching Update means the ACL changed.
	if !plan.BucketAcl.Equal(state.BucketAcl) {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		if err := retry.Retry(ensBucketAclTimeout, func() *retry.RetryError {
			_, e := r.Meta().RpcPost("Ens", "2017-11-10", "PutBucketAcl", nil, map[string]interface{}{
				"BucketName": state.BucketName.ValueString(),
				"BucketAcl":  plan.BucketAcl.ValueString(),
			}, true)
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
				fmt.Sprintf("Updating %s: calling PutBucketAcl", ensBucketTypeName),
				errs.WrapErrorf(err, "Updating %s: calling PutBucketAcl", ensBucketTypeName).Error(),
			)
			return
		}
	}

	info, err := getEnsBucket(r.Meta(), state.BucketName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetBucketInfo", ensBucketTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetBucketInfo", ensBucketTypeName).Error(),
		)
		return
	}
	if info == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, ensBucketState(info, plan))...)
}

func (r *ensBucketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ensBucketResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := retry.Retry(ensBucketDeleteTimeout, func() *retry.RetryError {
		_, e := r.Meta().RpcPost("Ens", "2017-11-10", "DeleteBucket", nil, map[string]interface{}{
			"BucketName": state.BucketName.ValueString(),
		}, true)
		if e != nil {
			if isEnsBucketNotExist(e) {
				return nil
			}
			if helper.NeedRetry(e) {
				wait()
				return retry.RetryableError(e)
			}
			return retry.NonRetryableError(e)
		}
		return nil
	})
	if err != nil && !isEnsBucketNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteBucket", ensBucketTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteBucket", ensBucketTypeName).Error(),
		)
		return
	}

	// ENS deletes the bucket asynchronously after DeleteBucket returns, so a
	// refresh right after destroy can still observe it. Poll GetBucketInfo
	// until it reports NoSuchBucket, mirroring the create-side ready wait.
	if waitErr := waitEnsBucketDeleted(r.Meta(), state.BucketName.ValueString()); waitErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: waiting for deletion", ensBucketTypeName),
			errs.WrapErrorf(waitErr, "Deleting %s: waiting for deletion", ensBucketTypeName).Error(),
		)
		return
	}
}

func (r *ensBucketResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getEnsBucket queries GetBucketInfo by bucket name and returns the BucketInfo
// object from the response, or nil when the bucket no longer exists. Used by
// Read/Update where NoSuchBucket means the resource is gone.
func getEnsBucket(client *connectivity.AliyunClient, bucketName string) (map[string]interface{}, error) {
	response, err := client.RpcPost("Ens", "2017-11-10", "GetBucketInfo", nil, map[string]interface{}{
		"BucketName": bucketName,
	}, true)
	if err != nil {
		if isEnsBucketNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	info, _ := response["BucketInfo"].(map[string]interface{})
	if info == nil {
		return nil, fmt.Errorf("the GetBucketInfo response carried no BucketInfo object: %v", response)
	}
	return info, nil
}

// waitEnsBucketReady polls GetBucketInfo until the bucket is queryable. Right
// after PutBucket the service may report NoSuchBucket during the creation
// window; that is retryable, while any other error aborts. Mirrors the ENS
// package's test setup helper.
func waitEnsBucketReady(client *connectivity.AliyunClient, bucketName string) (map[string]interface{}, error) {
	var info map[string]interface{}
	wait := helper.IncrementalWait(5*time.Second, 5*time.Second)
	err := retry.Retry(ensBucketReadyTimeout, func() *retry.RetryError {
		got, e := getEnsBucket(client, bucketName)
		if e != nil {
			// A transient failure (throttling, ServiceUnavailable, 5xx) during
			// the creation window must keep polling; only non-retryable errors
			// abort. Matches the ENS package's DescribeEnsInstance convention.
			if helper.NeedRetry(e) {
				wait()
				return retry.RetryableError(e)
			}
			return retry.NonRetryableError(e)
		}
		if got == nil {
			wait()
			return retry.RetryableError(fmt.Errorf("bucket %s is not visible via GetBucketInfo yet", bucketName))
		}
		info = got
		return nil
	})
	if err != nil {
		return nil, err
	}
	return info, nil
}

// waitEnsBucketDeleted polls GetBucketInfo until the bucket reports NoSuchBucket.
// ENS deletes buckets asynchronously after DeleteBucket returns, so without
// this wait a refresh immediately after destroy would still observe the bucket.
func waitEnsBucketDeleted(client *connectivity.AliyunClient, bucketName string) error {
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	return retry.Retry(ensBucketReadyTimeout, func() *retry.RetryError {
		got, e := getEnsBucket(client, bucketName)
		if e != nil {
			// Transient failures (throttling, ServiceUnavailable) keep polling;
			// any other error aborts, matching the create-side ready wait.
			if helper.NeedRetry(e) {
				wait()
				return retry.RetryableError(e)
			}
			return retry.NonRetryableError(e)
		}
		if got == nil {
			return nil
		}
		wait()
		return retry.RetryableError(fmt.Errorf("bucket %s is still deleting", bucketName))
	})
}

// ensBucketState maps the BucketInfo object of a GetBucketInfo response onto the
// resource model. Empty strings become null; the API-defaulted fields
// (LogicalBucketType, BucketAcl) fall back to their defaults so a refresh never
// desyncs from a config that relied on the API default. Write-only fields
// (EnsRegionId, DispatchScope) are not returned by GetBucketInfo, so they are
// preserved from the prior model passed by the caller.
func ensBucketState(info map[string]interface{}, prior ensBucketResourceModel) ensBucketResourceModel {
	str := func(key string) types.String {
		if v, ok := info[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	strDefault := func(key, def string) types.String {
		if v, ok := info[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringValue(def)
	}
	bucketName := str("BucketName")
	return ensBucketResourceModel{
		Id:                bucketName,
		BucketName:        bucketName,
		LogicalBucketType: strDefault("LogicalBucketType", "standard"),
		Comment:           str("Comment"),
		BucketAcl:         strDefault("BucketAcl", "private"),
		EnsRegionId:       prior.EnsRegionId,
		DispatchScope:     prior.DispatchScope,
		CreateTime:        str("CreateTime"),
		ModifyTime:        str("ModifyTime"),
	}
}

// isEnsBucketNotExist reports whether err is the ENS bucket-not-found error.
// GetBucketInfo and DeleteBucket surface NoSuchBucket; errs.NotFoundError
// covers wrapped forms, matching the ENS package convention.
func isEnsBucketNotExist(err error) bool {
	if err == nil {
		return false
	}
	return errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"NoSuchBucket"})
}
