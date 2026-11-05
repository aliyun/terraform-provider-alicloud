package oss

import (
	"context"
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	bucketEventNotificationTypeName    = "alicloud_oss_bucket_event_notification"
	bucketEventNotificationAction      = "/?eventNotification"
	bucketEventNotificationProduct     = "Oss"
	bucketEventNotificationAPIVersion  = "2019-05-17"
	bucketEventNotificationTimeout     = 5 * time.Minute
	bucketEventNotificationReadTimeout = 1 * time.Minute
)

var (
	_ resource.Resource                = &bucketEventNotificationResource{}
	_ resource.ResourceWithConfigure   = &bucketEventNotificationResource{}
	_ resource.ResourceWithImportState = &bucketEventNotificationResource{}
)

func NewBucketEventNotificationResource() resource.Resource {
	return &bucketEventNotificationResource{}
}

type bucketEventNotificationResource struct {
	fwadapt.ResourceBase
}

type bucketEventNotificationResourceModel struct {
	Id                           types.String                        `tfsdk:"id"`
	Bucket                       types.String                        `tfsdk:"bucket"`
	FunctionComputeConfiguration []functionComputeConfigurationModel `tfsdk:"function_compute_configuration"`
}

type functionComputeConfigurationModel struct {
	Id          types.String `tfsdk:"id"`
	Events      types.List   `tfsdk:"events"`
	Filter      *filterModel `tfsdk:"filter"`
	FunctionArn types.String `tfsdk:"function_arn"`
	AssumeRole  types.String `tfsdk:"assume_role"`
}

type filterModel struct {
	Key *filterKeyModel `tfsdk:"key"`
}

type filterKeyModel struct {
	Prefix types.String `tfsdk:"prefix"`
	Suffix types.String `tfsdk:"suffix"`
}

func (r *bucketEventNotificationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = bucketEventNotificationTypeName
}

func (r *bucketEventNotificationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bucket": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"function_compute_configuration": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"events": schema.ListAttribute{
							Required:    true,
							ElementType: types.StringType,
						},
						"function_arn": schema.StringAttribute{
							Required: true,
						},
						"assume_role": schema.StringAttribute{
							Optional: true,
						},
						"filter": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"key": schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{
										"prefix": schema.StringAttribute{
											Optional: true,
										},
										"suffix": schema.StringAttribute{
											Optional: true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *bucketEventNotificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan bucketEventNotificationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucket := plan.Bucket.ValueString()
	if err := r.putBucketEventNotification(bucket, plan.FunctionComputeConfiguration); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling PutBucketEventNotification", bucketEventNotificationTypeName),
			errs.WrapErrorf(err, "Creating %s: calling PutBucketEventNotification", bucketEventNotificationTypeName).Error(),
		)
		return
	}

	// The resource identity is the bucket name: PutBucketEventNotification
	// replaces the whole NotificationConfiguration, so one resource owns the
	// full configuration of one bucket, mirroring alicloud_oss_bucket_logging.
	id := bucket
	state := bucketEventNotificationResourceModel{
		Id:     types.StringValue(id),
		Bucket: plan.Bucket,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := r.describeBucketEventNotification(id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetBucketEventNotification", bucketEventNotificationTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetBucketEventNotification", bucketEventNotificationTypeName).Error(),
		)
		return
	}

	state = bucketEventNotificationState(object, id)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *bucketEventNotificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state bucketEventNotificationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := r.describeBucketEventNotification(state.Id.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetBucketEventNotification", bucketEventNotificationTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetBucketEventNotification", bucketEventNotificationTypeName).Error(),
		)
		return
	}

	model := bucketEventNotificationState(object, state.Id.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *bucketEventNotificationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state bucketEventNotificationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// PutBucketEventNotification replaces the whole NotificationConfiguration,
	// so any change to the list is applied by re-PUTting the full set.
	if err := r.putBucketEventNotification(state.Bucket.ValueString(), plan.FunctionComputeConfiguration); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling PutBucketEventNotification", bucketEventNotificationTypeName),
			errs.WrapErrorf(err, "Updating %s: calling PutBucketEventNotification", bucketEventNotificationTypeName).Error(),
		)
		return
	}

	object, err := r.describeBucketEventNotification(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetBucketEventNotification", bucketEventNotificationTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetBucketEventNotification", bucketEventNotificationTypeName).Error(),
		)
		return
	}

	model := bucketEventNotificationState(object, state.Id.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *bucketEventNotificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state bucketEventNotificationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.deleteBucketEventNotification(state.Bucket.ValueString(), state.Id.ValueString())
	if err != nil {
		// A missing configuration or a missing bucket both mean the resource is
		// already gone, mirroring logging/cors.
		if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"NoSuchBucket", "NoSuchNotificationConfiguration"}) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteBucketEventNotification", bucketEventNotificationTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteBucketEventNotification", bucketEventNotificationTypeName).Error(),
		)
		return
	}
}

func (r *bucketEventNotificationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// describeBucketEventNotification fetches the event-notification configuration
// of a bucket. The id is the bucket name. A bucket with no configuration reads
// back as an empty list, which is mapped to NotFound so a refresh drops the
// resource as it does for logging/cors.
func (r *bucketEventNotificationResource) describeBucketEventNotification(id string) ([]interface{}, error) {
	hostMap := map[string]*string{"bucket": helper.StringPointer(id)}
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(bucketEventNotificationReadTimeout, func() *retry.RetryError {
		var requestErr error
		response, requestErr = r.Meta().Do(bucketEventNotificationProduct, helper.XmlParam("GET", bucketEventNotificationAPIVersion, "GetBucketEventNotification", bucketEventNotificationAction), nil, nil, nil, hostMap, true)
		if requestErr != nil {
			if helper.NeedRetry(requestErr) {
				wait()
				return retry.RetryableError(requestErr)
			}
			return retry.NonRetryableError(requestErr)
		}
		return nil
	})
	logs.AddDebug(bucketEventNotificationAction, response, id)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchBucket", "NoSuchNotificationConfiguration"}) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("BucketEventNotification", id), errs.NotFoundMsg, response)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, bucketEventNotificationAction, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return nil, errs.WrapErrorf(errs.NotFoundErr("BucketEventNotification", id), errs.NotFoundMsg, response)
	}

	v, jsonErr := jsonpath.Get("$.NotificationConfiguration.FunctionComputeConfiguration", response)
	if jsonErr != nil || v == nil {
		// No configuration set: treat as NotFound so the resource is removed on
		// refresh, the same semantics as logging's empty TargetBucket.
		return nil, errs.WrapErrorf(errs.NotFoundErr("BucketEventNotification", id), errs.NotFoundMsg, response)
	}
	configs, ok := v.([]interface{})
	if !ok || len(configs) == 0 {
		return nil, errs.WrapErrorf(errs.NotFoundErr("BucketEventNotification", id), errs.NotFoundMsg, response)
	}
	return configs, nil
}

func (r *bucketEventNotificationResource) putBucketEventNotification(bucket string, configs []functionComputeConfigurationModel) error {
	hostMap := map[string]*string{"bucket": helper.StringPointer(bucket)}
	items := make([]interface{}, 0, len(configs))
	for _, cfg := range configs {
		item := map[string]interface{}{}
		events := make([]interface{}, 0)
		if !cfg.Events.IsNull() && !cfg.Events.IsUnknown() {
			var eventStrings []string
			if diags := cfg.Events.ElementsAs(context.Background(), &eventStrings, false); !diags.HasError() {
				for _, e := range eventStrings {
					events = append(events, e)
				}
			}
		}
		item["Event"] = events

		if cfg.Filter != nil && cfg.Filter.Key != nil {
			keyMap := map[string]interface{}{}
			if !cfg.Filter.Key.Prefix.IsNull() {
				keyMap["Prefix"] = cfg.Filter.Key.Prefix.ValueString()
			}
			if !cfg.Filter.Key.Suffix.IsNull() {
				keyMap["Suffix"] = cfg.Filter.Key.Suffix.ValueString()
			}
			item["Filter"] = map[string]interface{}{"Key": keyMap}
		}

		funcMap := map[string]interface{}{}
		funcMap["Arn"] = cfg.FunctionArn.ValueString()
		if !cfg.AssumeRole.IsNull() {
			funcMap["AssumeRole"] = cfg.AssumeRole.ValueString()
		}
		item["Function"] = funcMap

		items = append(items, item)
	}

	body := map[string]interface{}{
		"NotificationConfiguration": map[string]interface{}{
			"FunctionComputeConfiguration": items,
		},
	}

	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(bucketEventNotificationTimeout, func() *retry.RetryError {
		var requestErr error
		response, requestErr = r.Meta().Do(bucketEventNotificationProduct, helper.XmlParam("PUT", bucketEventNotificationAPIVersion, "PutBucketEventNotification", bucketEventNotificationAction), nil, body, nil, hostMap, false)
		if requestErr != nil {
			if helper.NeedRetry(requestErr) {
				wait()
				return retry.RetryableError(requestErr)
			}
			return retry.NonRetryableError(requestErr)
		}
		return nil
	})
	logs.AddDebug(bucketEventNotificationAction, response, body)
	return err
}

func (r *bucketEventNotificationResource) deleteBucketEventNotification(bucket, id string) error {
	hostMap := map[string]*string{"bucket": helper.StringPointer(bucket)}
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(bucketEventNotificationTimeout, func() *retry.RetryError {
		var requestErr error
		response, requestErr = r.Meta().Do(bucketEventNotificationProduct, helper.XmlParam("DELETE", bucketEventNotificationAPIVersion, "DeleteBucketEventNotification", bucketEventNotificationAction), nil, nil, nil, hostMap, false)
		if requestErr != nil {
			if helper.NeedRetry(requestErr) {
				wait()
				return retry.RetryableError(requestErr)
			}
			return retry.NonRetryableError(requestErr)
		}
		return nil
	})
	logs.AddDebug(bucketEventNotificationAction, response, id)
	return err
}

// bucketEventNotificationState maps a GetBucketEventNotification response onto
// the resource model. The configuration id of each item is server-assigned.
func bucketEventNotificationState(configs []interface{}, id string) bucketEventNotificationResourceModel {
	model := bucketEventNotificationResourceModel{
		Id:     types.StringValue(id),
		Bucket: types.StringValue(id),
	}
	items := make([]functionComputeConfigurationModel, 0, len(configs))
	for _, raw := range configs {
		cfg, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		item := functionComputeConfigurationModel{}

		if v, ok := cfg["ID"].(string); ok && v != "" {
			item.Id = types.StringValue(v)
		} else {
			item.Id = types.StringNull()
		}

		eventValues := make([]attr.Value, 0)
		if events, ok := cfg["Event"].([]interface{}); ok {
			for _, e := range events {
				if s, ok := e.(string); ok {
					eventValues = append(eventValues, types.StringValue(s))
				}
			}
		}
		if len(eventValues) > 0 {
			item.Events = types.ListValueMust(types.StringType, eventValues)
		} else {
			item.Events = types.ListNull(types.StringType)
		}

		item.FunctionArn = types.StringNull()
		item.AssumeRole = types.StringNull()
		if fn, ok := cfg["Function"].(map[string]interface{}); ok {
			if arn, ok := fn["Arn"].(string); ok {
				item.FunctionArn = types.StringValue(arn)
			}
			if role, ok := fn["AssumeRole"].(string); ok && role != "" {
				item.AssumeRole = types.StringValue(role)
			}
		}

		if filter, ok := cfg["Filter"].(map[string]interface{}); ok {
			if key, ok := filter["Key"].(map[string]interface{}); ok {
				keyModel := &filterKeyModel{}
				if p, ok := key["Prefix"].(string); ok && p != "" {
					keyModel.Prefix = types.StringValue(p)
				} else {
					keyModel.Prefix = types.StringNull()
				}
				if s, ok := key["Suffix"].(string); ok && s != "" {
					keyModel.Suffix = types.StringValue(s)
				} else {
					keyModel.Suffix = types.StringNull()
				}
				item.Filter = &filterModel{Key: keyModel}
			}
		}

		items = append(items, item)
	}
	model.FunctionComputeConfiguration = items
	return model
}
