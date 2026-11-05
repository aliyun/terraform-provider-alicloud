package oss

import (
	"context"
	"fmt"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	ossBucketLifecycleTypeName = "alicloud_oss_bucket_lifecycle"
)

var (
	_ resource.Resource                = &ossBucketLifecycleResource{}
	_ resource.ResourceWithConfigure   = &ossBucketLifecycleResource{}
	_ resource.ResourceWithImportState = &ossBucketLifecycleResource{}
)

// NewOssBucketLifecycleResource is the framework factory registered in the OSS
// ServicePackage.Resources list. It manages the lifecycle configuration of an
// existing OSS bucket, decoupling lifecycle_rule from alicloud_oss_bucket the
// same way cors/website/logging/versioning were already split out.
func NewOssBucketLifecycleResource() resource.Resource {
	return &ossBucketLifecycleResource{}
}

type ossBucketLifecycleResource struct {
	fwadapt.ResourceBase
}

// ossBucketLifecycleResourceModel is the top-level state/plan model. The id is
// the bucket name: a bucket has at most one lifecycle configuration.
type ossBucketLifecycleResourceModel struct {
	Id                                  types.String `tfsdk:"id"`
	Bucket                              types.String `tfsdk:"bucket"`
	LifecycleRuleAllowSameActionOverlap types.Bool   `tfsdk:"lifecycle_rule_allow_same_action_overlap"`
	LifecycleRule                       types.List   `tfsdk:"lifecycle_rule"`
}

type lifecycleRuleModel struct {
	Id                          types.String `tfsdk:"id"`
	Prefix                      types.String `tfsdk:"prefix"`
	Enabled                     types.Bool   `tfsdk:"enabled"`
	Expiration                  types.Set    `tfsdk:"expiration"`
	Transitions                 types.Set    `tfsdk:"transitions"`
	AbortMultipartUpload        types.Set    `tfsdk:"abort_multipart_upload"`
	NoncurrentVersionExpiration types.Set    `tfsdk:"noncurrent_version_expiration"`
	NoncurrentVersionTransition types.Set    `tfsdk:"noncurrent_version_transition"`
	Tags                        types.Map    `tfsdk:"tags"`
	Filter                      types.List   `tfsdk:"filter"`
}

type expirationModel struct {
	Date                      types.String `tfsdk:"date"`
	CreatedBeforeDate         types.String `tfsdk:"created_before_date"`
	Days                      types.Int64  `tfsdk:"days"`
	ExpiredObjectDeleteMarker types.Bool   `tfsdk:"expired_object_delete_marker"`
}

type transitionsModel struct {
	CreatedBeforeDate    types.String `tfsdk:"created_before_date"`
	Days                 types.Int64  `tfsdk:"days"`
	StorageClass         types.String `tfsdk:"storage_class"`
	IsAccessTime         types.Bool   `tfsdk:"is_access_time"`
	ReturnToStdWhenVisit types.Bool   `tfsdk:"return_to_std_when_visit"`
}

type abortMultipartUploadModel struct {
	CreatedBeforeDate types.String `tfsdk:"created_before_date"`
	Days              types.Int64  `tfsdk:"days"`
}

type noncurrentVersionExpirationModel struct {
	Days types.Int64 `tfsdk:"days"`
}

type noncurrentVersionTransitionModel struct {
	Days                 types.Int64  `tfsdk:"days"`
	StorageClass         types.String `tfsdk:"storage_class"`
	IsAccessTime         types.Bool   `tfsdk:"is_access_time"`
	ReturnToStdWhenVisit types.Bool   `tfsdk:"return_to_std_when_visit"`
}

type filterModel struct {
	Not                   types.List  `tfsdk:"not"`
	ObjectSizeGreaterThan types.Int64 `tfsdk:"object_size_greater_than"`
	ObjectSizeLessThan    types.Int64 `tfsdk:"object_size_less_than"`
}

type filterNotModel struct {
	Prefix types.String `tfsdk:"prefix"`
	Tag    types.List   `tfsdk:"tag"`
}

type filterNotTagModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func (r *ossBucketLifecycleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ossBucketLifecycleTypeName
}

func (r *ossBucketLifecycleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	storageClassValidator := stringvalidator.OneOf(
		string(oss.StorageStandard),
		string(oss.StorageIA),
		string(oss.StorageArchive),
		string(oss.StorageColdArchive),
		string(oss.StorageDeepColdArchive),
	)

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
			"lifecycle_rule_allow_same_action_overlap": schema.BoolAttribute{
				Optional: true,
			},
		},
		Blocks: map[string]schema.Block{
			"lifecycle_rule": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
							Validators: []validator.String{
								stringvalidator.LengthBetween(0, 255),
							},
						},
						"prefix": schema.StringAttribute{
							Optional: true,
						},
						"enabled": schema.BoolAttribute{
							Required: true,
						},
						"tags": schema.MapAttribute{
							Optional:    true,
							ElementType: types.StringType,
						},
					},
					Blocks: map[string]schema.Block{
						"expiration": schema.SetNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"date": schema.StringAttribute{
										Optional: true,
										Validators: []validator.String{
											ossLifecycleDateValidator{},
										},
									},
									"created_before_date": schema.StringAttribute{
										Optional: true,
										Validators: []validator.String{
											ossLifecycleDateValidator{},
										},
									},
									"days": schema.Int64Attribute{
										Optional: true,
									},
									"expired_object_delete_marker": schema.BoolAttribute{
										Optional: true,
									},
								},
							},
						},
						"transitions": schema.SetNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"created_before_date": schema.StringAttribute{
										Optional: true,
										Validators: []validator.String{
											ossLifecycleDateValidator{},
										},
									},
									"days": schema.Int64Attribute{
										Optional: true,
									},
									"storage_class": schema.StringAttribute{
										Required:   true,
										Validators: []validator.String{storageClassValidator},
									},
									"is_access_time": schema.BoolAttribute{
										Optional: true,
									},
									"return_to_std_when_visit": schema.BoolAttribute{
										Optional: true,
									},
								},
							},
						},
						"abort_multipart_upload": schema.SetNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"created_before_date": schema.StringAttribute{
										Optional: true,
										Validators: []validator.String{
											ossLifecycleDateValidator{},
										},
									},
									"days": schema.Int64Attribute{
										Optional: true,
									},
								},
							},
						},
						"noncurrent_version_expiration": schema.SetNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"days": schema.Int64Attribute{
										Required: true,
									},
								},
							},
						},
						"noncurrent_version_transition": schema.SetNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"days": schema.Int64Attribute{
										Required: true,
									},
									"storage_class": schema.StringAttribute{
										Required:   true,
										Validators: []validator.String{storageClassValidator},
									},
									"is_access_time": schema.BoolAttribute{
										Optional: true,
									},
									"return_to_std_when_visit": schema.BoolAttribute{
										Optional: true,
									},
								},
							},
						},
						"filter": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"object_size_greater_than": schema.Int64Attribute{
										Optional: true,
										Validators: []validator.Int64{
											int64validator.AtLeast(1),
										},
									},
									"object_size_less_than": schema.Int64Attribute{
										Optional: true,
										Validators: []validator.Int64{
											int64validator.AtLeast(1),
										},
									},
								},
								Blocks: map[string]schema.Block{
									"not": schema.ListNestedBlock{
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"prefix": schema.StringAttribute{
													Optional: true,
												},
											},
											Blocks: map[string]schema.Block{
												"tag": schema.ListNestedBlock{
													NestedObject: schema.NestedBlockObject{
														Attributes: map[string]schema.Attribute{
															"key": schema.StringAttribute{
																Required: true,
															},
															"value": schema.StringAttribute{
																Required: true,
															},
														},
													},
												},
											},
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

func (r *ossBucketLifecycleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ossBucketLifecycleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucket := plan.Bucket.ValueString()
	if err := r.putLifecycle(ctx, plan); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling SetBucketLifecycle", ossBucketLifecycleTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, bucket, "SetBucketLifecycle", errs.AliyunOssGoSdk).Error(),
		)
		return
	}

	// Build state from the planned rules via a round-trip through the OSS SDK
	// types. Two alternatives were tried and rejected:
	//  1. resp.State.Set(ctx, &plan) — triggers a Value Conversion Error in
	//     terraform-plugin-framework v1.19.0 because zero-value attr.Value
	//     fields (types.Set{}, types.List{}, …) lose their element type
	//     during framework reflection.
	//  2. Reading back from GetBucketLifecycle — avoids the reflection path
	//     but introduces an eventual-consistency race: OSS may return
	//     NoSuchLifecycle or the previous configuration immediately after
	//     SetBucketLifecycle/DeleteBucketLifecycle, producing an
	//     "inconsistent result after apply" error.
	//
	// Converting the plan to oss.LifecycleRule and back through
	// ossRulesToModel — which constructs every nested object and collection
	// with explicit ObjectValue/SetValue/ListValue helpers — sidesteps both
	// problems: the state is correctly typed and always matches the plan.
	// The actual API state is reconciled during the next Read.
	rules, err := modelToOssRules(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: building state from plan", ossBucketLifecycleTypeName),
			err.Error(),
		)
		return
	}
	readBack, err := ossRulesToModel(ctx, bucket, rules)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: building state from plan", ossBucketLifecycleTypeName),
			err.Error(),
		)
		return
	}

	readBack.Bucket = plan.Bucket
	readBack.Id = types.StringValue(bucket)
	readBack.LifecycleRuleAllowSameActionOverlap = plan.LifecycleRuleAllowSameActionOverlap
	resp.Diagnostics.Append(resp.State.Set(ctx, &readBack)...)
}

func (r *ossBucketLifecycleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ossBucketLifecycleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucket := state.Id.ValueString()
	model, err := r.getLifecycle(ctx, bucket)
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetBucketLifecycle", ossBucketLifecycleTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, bucket, "GetBucketLifecycle", errs.AliyunOssGoSdk).Error(),
		)
		return
	}

	model.Id = state.Id
	model.LifecycleRuleAllowSameActionOverlap = state.LifecycleRuleAllowSameActionOverlap
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *ossBucketLifecycleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ossBucketLifecycleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state ossBucketLifecycleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucket := state.Id.ValueString()
	if err := r.putLifecycle(ctx, plan); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling SetBucketLifecycle", ossBucketLifecycleTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, bucket, "SetBucketLifecycle", errs.AliyunOssGoSdk).Error(),
		)
		return
	}

	// Build state from the planned rules via a round-trip through the OSS SDK
	// types (see Create for the rationale).
	rules, err := modelToOssRules(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: building state from plan", ossBucketLifecycleTypeName),
			err.Error(),
		)
		return
	}
	readBack, err := ossRulesToModel(ctx, bucket, rules)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: building state from plan", ossBucketLifecycleTypeName),
			err.Error(),
		)
		return
	}

	readBack.Bucket = state.Bucket
	readBack.Id = state.Id
	readBack.LifecycleRuleAllowSameActionOverlap = plan.LifecycleRuleAllowSameActionOverlap
	resp.Diagnostics.Append(resp.State.Set(ctx, &readBack)...)
}

func (r *ossBucketLifecycleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ossBucketLifecycleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucket := state.Id.ValueString()
	if err := r.deleteLifecycle(ctx, bucket); err != nil {
		if errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteBucketLifecycle", ossBucketLifecycleTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, bucket, "DeleteBucketLifecycle", errs.AliyunOssGoSdk).Error(),
		)
		return
	}
}

func (r *ossBucketLifecycleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// putLifecycle sets the bucket lifecycle configuration to the planned rules, or
// deletes it when no rules are configured. This mirrors the existing
// resourceAlicloudOssBucketLifecycleRuleUpdate semantics embedded in
// alicloud_oss_bucket.
func (r *ossBucketLifecycleResource) putLifecycle(ctx context.Context, plan ossBucketLifecycleResourceModel) error {
	bucket := plan.Bucket.ValueString()
	rules, err := modelToOssRules(ctx, plan)
	if err != nil {
		return err
	}
	if len(rules) == 0 {
		return r.deleteLifecycle(ctx, bucket)
	}

	options := []oss.Option{}
	if plan.LifecycleRuleAllowSameActionOverlap.ValueBool() {
		options = append(options, oss.AllowSameActionOverLap(true))
	}

	_, err = r.Meta().WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return nil, ossClient.SetBucketLifecycle(bucket, rules, options...)
	})
	logs.AddDebug("SetBucketLifecycle", nil, bucket, map[string]interface{}{
		"bucketName": bucket,
		"rules":      rules,
	})
	return err
}

func (r *ossBucketLifecycleResource) getLifecycle(ctx context.Context, bucket string) (ossBucketLifecycleResourceModel, error) {
	raw, err := r.Meta().WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return ossClient.GetBucketLifecycle(bucket)
	})
	if err != nil {
		// NoSuchBucket means the bucket itself is gone: surface a NotFound error
		// so Read removes the resource. NoSuchLifecycle (bucket exists but has no
		// lifecycle configuration) is treated as an empty rule set below.
		if e, ok := err.(oss.ServiceError); ok {
			if e.Code == "NoSuchBucket" {
				return ossBucketLifecycleResourceModel{}, errs.WrapErrorf(errs.NotFoundErr("OssBucket", bucket), errs.NotFoundMsg, bucket)
			}
			if e.Code == "NoSuchLifecycle" {
				logs.AddDebug("GetBucketLifecycle", "no lifecycle configuration on bucket", bucket)
				// LifecycleRule must be a typed null List, not a zero-value
				// types.List{} (MISSING type), to avoid a Value Conversion
				// Error when the framework serializes state.
				return ossBucketLifecycleResourceModel{Id: types.StringValue(bucket), Bucket: types.StringValue(bucket), LifecycleRule: types.ListNull(lifecycleRuleObjectElemType())}, nil
			}
		}
		return ossBucketLifecycleResourceModel{}, errs.WrapErrorf(err, errs.DefaultErrorMsg, bucket, "GetBucketLifecycle", errs.AliyunOssGoSdk)
	}
	logs.AddDebug("GetBucketLifecycle", raw, bucket)
	lifecycle, _ := raw.(oss.GetBucketLifecycleResult)
	return ossRulesToModel(ctx, bucket, lifecycle.Rules)
}

func (r *ossBucketLifecycleResource) deleteLifecycle(ctx context.Context, bucket string) error {
	_, err := r.Meta().WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return nil, ossClient.DeleteBucketLifecycle(bucket)
	})
	logs.AddDebug("DeleteBucketLifecycle", nil, bucket, map[string]interface{}{
		"bucketName": bucket,
	})
	return err
}

// modelToOssRules converts the planned lifecycle_rule blocks into the
// aliyun-oss-go-sdk LifecycleRule slice, enforcing the same per-element
// mutual-exclusion rules that the embedded bucket implementation enforced in
// resourceAlicloudOssBucketLifecycleRuleUpdate.
func modelToOssRules(ctx context.Context, plan ossBucketLifecycleResourceModel) ([]oss.LifecycleRule, error) {
	if plan.LifecycleRule.IsNull() || plan.LifecycleRule.IsUnknown() {
		return nil, nil
	}
	var ruleModels []lifecycleRuleModel
	respDiags := plan.LifecycleRule.ElementsAs(ctx, &ruleModels, false)
	if respDiags.HasError() {
		return nil, fmt.Errorf("parsing lifecycle_rule: %v", respDiags)
	}

	rules := make([]oss.LifecycleRule, 0, len(ruleModels))
	for _, rm := range ruleModels {
		rule := oss.LifecycleRule{
			Prefix: rm.Prefix.ValueString(),
		}
		if id := rm.Id.ValueString(); id != "" {
			rule.ID = id
		}
		if rm.Enabled.ValueBool() {
			rule.Status = string(ExpirationStatusEnabled)
		} else {
			rule.Status = string(ExpirationStatusDisabled)
		}

		if !rm.Expiration.IsNull() && !rm.Expiration.IsUnknown() {
			var exps []expirationModel
			d := rm.Expiration.ElementsAs(ctx, &exps, false)
			if d.HasError() {
				return nil, fmt.Errorf("parsing expiration: %v", d)
			}
			if len(exps) > 0 {
				e := exps[0]
				i := oss.LifecycleExpiration{}
				valDate := e.Date.ValueString()
				valCreatedBeforeDate := e.CreatedBeforeDate.ValueString()
				valDays := int(e.Days.ValueInt64())
				if e.ExpiredObjectDeleteMarker.ValueBool() {
					if valDays > 0 || valDate != "" || valCreatedBeforeDate != "" {
						return nil, errs.WrapError(errs.Error("'date/created_before_date/days' conflicts with 'expired_object_delete_marker'. One and only one of them can be specified in one expiration configuration."))
					}
					val := e.ExpiredObjectDeleteMarker.ValueBool()
					i.ExpiredObjectDeleteMarker = &val
				} else {
					cnt := 0
					if valDate != "" {
						i.Date = fmt.Sprintf("%sT00:00:00.000Z", valDate)
						cnt++
					}
					if valCreatedBeforeDate != "" {
						i.CreatedBeforeDate = fmt.Sprintf("%sT00:00:00.000Z", valCreatedBeforeDate)
						cnt++
					}
					if valDays > 0 {
						i.Days = valDays
						cnt++
					}
					if cnt != 1 {
						return nil, errs.WrapError(errs.Error("One and only one of 'date', 'created_before_date' and 'days' can be specified in one expiration configuration."))
					}
				}
				rule.Expiration = &i
			}
		}

		if !rm.Transitions.IsNull() && !rm.Transitions.IsUnknown() {
			var ts []transitionsModel
			d := rm.Transitions.ElementsAs(ctx, &ts, false)
			if d.HasError() {
				return nil, fmt.Errorf("parsing transitions: %v", d)
			}
			for _, t := range ts {
				i := oss.LifecycleTransition{}
				valCreatedBeforeDate := t.CreatedBeforeDate.ValueString()
				valDays := int(t.Days.ValueInt64())
				valStorageClass := t.StorageClass.ValueString()
				if (valCreatedBeforeDate != "" && valDays > 0) || (valCreatedBeforeDate == "" && valDays <= 0) || valStorageClass == "" {
					return nil, errs.WrapError(errs.Error("'CreatedBeforeDate' conflicts with 'Days'. One and only one of them can be specified in one transition configuration. 'storage_class' must be set."))
				}
				if valCreatedBeforeDate != "" {
					i.CreatedBeforeDate = fmt.Sprintf("%sT00:00:00.000Z", valCreatedBeforeDate)
				}
				if valDays > 0 {
					i.Days = valDays
				}
				if valStorageClass != "" {
					i.StorageClass = oss.StorageClassType(valStorageClass)
				}
				if t.IsAccessTime.ValueBool() {
					val := t.IsAccessTime.ValueBool()
					i.IsAccessTime = &val
					val1 := t.ReturnToStdWhenVisit.ValueBool()
					i.ReturnToStdWhenVisit = &val1
				}
				rule.Transitions = append(rule.Transitions, i)
			}
		}

		if !rm.AbortMultipartUpload.IsNull() && !rm.AbortMultipartUpload.IsUnknown() {
			var ams []abortMultipartUploadModel
			d := rm.AbortMultipartUpload.ElementsAs(ctx, &ams, false)
			if d.HasError() {
				return nil, fmt.Errorf("parsing abort_multipart_upload: %v", d)
			}
			if len(ams) > 0 {
				am := ams[0]
				i := oss.LifecycleAbortMultipartUpload{}
				valCreatedBeforeDate := am.CreatedBeforeDate.ValueString()
				valDays := int(am.Days.ValueInt64())
				if (valCreatedBeforeDate != "" && valDays > 0) || (valCreatedBeforeDate == "" && valDays <= 0) {
					return nil, errs.WrapError(errs.Error("'CreatedBeforeDate' conflicts with 'days'. One and only one of them can be specified in one abort_multipart_upload configuration."))
				}
				if valCreatedBeforeDate != "" {
					i.CreatedBeforeDate = fmt.Sprintf("%sT00:00:00.000Z", valCreatedBeforeDate)
				}
				if valDays > 0 {
					i.Days = valDays
				}
				rule.AbortMultipartUpload = &i
			}
		}

		if !rm.NoncurrentVersionExpiration.IsNull() && !rm.NoncurrentVersionExpiration.IsUnknown() {
			var nves []noncurrentVersionExpirationModel
			d := rm.NoncurrentVersionExpiration.ElementsAs(ctx, &nves, false)
			if d.HasError() {
				return nil, fmt.Errorf("parsing noncurrent_version_expiration: %v", d)
			}
			if len(nves) > 0 {
				i := oss.LifecycleVersionExpiration{}
				i.NoncurrentDays = int(nves[0].Days.ValueInt64())
				rule.NonVersionExpiration = &i
			}
		}

		if !rm.NoncurrentVersionTransition.IsNull() && !rm.NoncurrentVersionTransition.IsUnknown() {
			var nvts []noncurrentVersionTransitionModel
			d := rm.NoncurrentVersionTransition.ElementsAs(ctx, &nvts, false)
			if d.HasError() {
				return nil, fmt.Errorf("parsing noncurrent_version_transition: %v", d)
			}
			for _, t := range nvts {
				i := oss.LifecycleVersionTransition{}
				i.NoncurrentDays = int(t.Days.ValueInt64())
				i.StorageClass = oss.StorageClassType(t.StorageClass.ValueString())
				if t.IsAccessTime.ValueBool() {
					val := t.IsAccessTime.ValueBool()
					i.IsAccessTime = &val
					val1 := t.ReturnToStdWhenVisit.ValueBool()
					i.ReturnToStdWhenVisit = &val1
				}
				rule.NonVersionTransitions = append(rule.NonVersionTransitions, i)
			}
		}

		if !rm.Tags.IsNull() && !rm.Tags.IsUnknown() {
			tagsMap := make(map[string]string)
			d := rm.Tags.ElementsAs(ctx, &tagsMap, false)
			if d.HasError() {
				return nil, fmt.Errorf("parsing tags: %v", d)
			}
			for k, v := range tagsMap {
				rule.Tags = append(rule.Tags, oss.Tag{
					Key:   k,
					Value: v,
				})
			}
		}

		if !rm.Filter.IsNull() && !rm.Filter.IsUnknown() {
			var filters []filterModel
			d := rm.Filter.ElementsAs(ctx, &filters, false)
			if d.HasError() {
				return nil, fmt.Errorf("parsing filter: %v", d)
			}
			if len(filters) > 0 {
				f := filters[0]
				filter := oss.LifecycleFilter{}
				if !f.Not.IsNull() && !f.Not.IsUnknown() {
					var nots []filterNotModel
					dn := f.Not.ElementsAs(ctx, &nots, false)
					if dn.HasError() {
						return nil, fmt.Errorf("parsing filter.not: %v", dn)
					}
					for _, n := range nots {
						filterNot := oss.LifecycleFilterNot{}
						filterNot.Prefix = n.Prefix.ValueString()
						if !n.Tag.IsNull() && !n.Tag.IsUnknown() {
							var tags []filterNotTagModel
							dt := n.Tag.ElementsAs(ctx, &tags, false)
							if dt.HasError() {
								return nil, fmt.Errorf("parsing filter.not.tag: %v", dt)
							}
							if len(tags) > 0 {
								tag := oss.Tag{}
								tag.Key = tags[0].Key.ValueString()
								tag.Value = tags[0].Value.ValueString()
								filterNot.Tag = &tag
							}
						}
						filter.Not = append(filter.Not, filterNot)
					}
				}
				if v := f.ObjectSizeGreaterThan.ValueInt64(); v > 0 {
					greater := int64(v)
					filter.ObjectSizeGreaterThan = &greater
				}
				if v := f.ObjectSizeLessThan.ValueInt64(); v > 0 {
					less := int64(v)
					filter.ObjectSizeLessThan = &less
				}
				rule.Filter = &filter
			}
		}

		rules = append(rules, rule)
	}
	return rules, nil
}

// ossRulesToModel maps the aliyun-oss-go-sdk LifecycleRule slice back onto the
// resource model, mirroring the read path embedded in
// alicloud_oss_bucket (resource_alicloud_oss_bucket.go read lifecycle).
func ossRulesToModel(ctx context.Context, bucket string, rules []oss.LifecycleRule) (ossBucketLifecycleResourceModel, error) {
	model := ossBucketLifecycleResourceModel{
		Id:     types.StringValue(bucket),
		Bucket: types.StringValue(bucket),
	}
	if len(rules) == 0 {
		// An empty rule set still needs a properly typed null List; a
		// zero-value types.List{} carries a MISSING (dynamic) element
		// type that triggers a Value Conversion Error during State.Set.
		model.LifecycleRule = types.ListNull(lifecycleRuleObjectElemType())
		return model, nil
	}

	// Build each lifecycle_rule element as an explicit ObjectValue rather than
	// reflecting a Go struct slice via ListValueFrom/SetValueFrom. In
	// terraform-plugin-framework v1.19.0 the non-generic *ValueFrom helpers use
	// reflection that can produce Set/List values with a missing (dynamic)
	// element type when the source struct has zero-value attr.Value fields
	// (types.String{}, types.Int64{}, …). The resulting types.List passes
	// HasError checks but later triggers a "Value Conversion Error" at
	// resp.State.Set. Constructing every nested object with ObjectValue and
	// every collection with SetValue/ListValue (taking []attr.Value) sidesteps
	// the reflection path entirely.
	ruleElems := make([]attr.Value, 0, len(rules))
	for _, lr := range rules {
		expirationSet := types.SetNull(expirationObjectElemType())
		if lr.Expiration != nil {
			expAttrs := map[string]attr.Value{
				"date":                         types.StringNull(),
				"created_before_date":          types.StringNull(),
				"days":                         types.Int64Null(),
				"expired_object_delete_marker": types.BoolNull(),
			}
			if lr.Expiration.Date != "" {
				t, err := time.Parse("2006-01-02T15:04:05.000Z", lr.Expiration.Date)
				if err != nil {
					return model, errs.WrapError(err)
				}
				expAttrs["date"] = types.StringValue(t.Format("2006-01-02"))
			}
			if lr.Expiration.CreatedBeforeDate != "" {
				t, err := time.Parse("2006-01-02T15:04:05.000Z", lr.Expiration.CreatedBeforeDate)
				if err != nil {
					return model, errs.WrapError(err)
				}
				expAttrs["created_before_date"] = types.StringValue(t.Format("2006-01-02"))
			}
			if lr.Expiration.ExpiredObjectDeleteMarker != nil {
				expAttrs["expired_object_delete_marker"] = types.BoolValue(*lr.Expiration.ExpiredObjectDeleteMarker)
			}
			if lr.Expiration.Days > 0 {
				expAttrs["days"] = types.Int64Value(int64(lr.Expiration.Days))
			}
			obj, d := types.ObjectValue(expirationObjectElemType().AttrTypes, expAttrs)
			if d.HasError() {
				return model, fmt.Errorf("building expiration object: %v", d)
			}
			s, d := types.SetValue(expirationObjectElemType(), []attr.Value{obj})
			if d.HasError() {
				return model, fmt.Errorf("building expiration set: %v", d)
			}
			expirationSet = s
		}

		// transitions
		transitionsSet := types.SetNull(transitionsObjectElemType())
		if len(lr.Transitions) != 0 {
			tElems := make([]attr.Value, 0, len(lr.Transitions))
			for _, t := range lr.Transitions {
				tAttrs := map[string]attr.Value{
					"created_before_date":      types.StringNull(),
					"days":                     types.Int64Null(),
					"storage_class":            types.StringNull(),
					"is_access_time":           types.BoolNull(),
					"return_to_std_when_visit": types.BoolNull(),
				}
				if t.CreatedBeforeDate != "" {
					dt, err := time.Parse("2006-01-02T15:04:05.000Z", t.CreatedBeforeDate)
					if err != nil {
						return model, errs.WrapError(err)
					}
					tAttrs["created_before_date"] = types.StringValue(dt.Format("2006-01-02"))
				}
				if t.Days > 0 {
					tAttrs["days"] = types.Int64Value(int64(t.Days))
				}
				tAttrs["storage_class"] = types.StringValue(string(t.StorageClass))
				if t.IsAccessTime != nil {
					tAttrs["is_access_time"] = types.BoolValue(*t.IsAccessTime)
				}
				if t.ReturnToStdWhenVisit != nil {
					tAttrs["return_to_std_when_visit"] = types.BoolValue(*t.ReturnToStdWhenVisit)
				}
				obj, d := types.ObjectValue(transitionsObjectElemType().AttrTypes, tAttrs)
				if d.HasError() {
					return model, fmt.Errorf("building transitions object: %v", d)
				}
				tElems = append(tElems, obj)
			}
			s, d := types.SetValue(transitionsObjectElemType(), tElems)
			if d.HasError() {
				return model, fmt.Errorf("building transitions set: %v", d)
			}
			transitionsSet = s
		}

		// abort_multipart_upload
		amuSet := types.SetNull(abortMultipartUploadObjectElemType())
		if lr.AbortMultipartUpload != nil {
			amAttrs := map[string]attr.Value{
				"created_before_date": types.StringNull(),
				"days":                types.Int64Null(),
			}
			if lr.AbortMultipartUpload.CreatedBeforeDate != "" {
				t, err := time.Parse("2006-01-02T15:04:05.000Z", lr.AbortMultipartUpload.CreatedBeforeDate)
				if err != nil {
					return model, errs.WrapError(err)
				}
				amAttrs["created_before_date"] = types.StringValue(t.Format("2006-01-02"))
			}
			if lr.AbortMultipartUpload.Days > 0 {
				amAttrs["days"] = types.Int64Value(int64(lr.AbortMultipartUpload.Days))
			}
			obj, d := types.ObjectValue(abortMultipartUploadObjectElemType().AttrTypes, amAttrs)
			if d.HasError() {
				return model, fmt.Errorf("building abort_multipart_upload object: %v", d)
			}
			s, d := types.SetValue(abortMultipartUploadObjectElemType(), []attr.Value{obj})
			if d.HasError() {
				return model, fmt.Errorf("building abort_multipart_upload set: %v", d)
			}
			amuSet = s
		}

		// noncurrent_version_expiration
		nveSet := types.SetNull(noncurrentVersionExpirationObjectElemType())
		if lr.NonVersionExpiration != nil {
			nveAttrs := map[string]attr.Value{
				"days": types.Int64Null(),
			}
			nveAttrs["days"] = types.Int64Value(int64(lr.NonVersionExpiration.NoncurrentDays))
			obj, d := types.ObjectValue(noncurrentVersionExpirationObjectElemType().AttrTypes, nveAttrs)
			if d.HasError() {
				return model, fmt.Errorf("building noncurrent_version_expiration object: %v", d)
			}
			s, d := types.SetValue(noncurrentVersionExpirationObjectElemType(), []attr.Value{obj})
			if d.HasError() {
				return model, fmt.Errorf("building noncurrent_version_expiration set: %v", d)
			}
			nveSet = s
		}

		// noncurrent_version_transition
		nvtSet := types.SetNull(noncurrentVersionTransitionObjectElemType())
		if len(lr.NonVersionTransitions) != 0 {
			nvtElems := make([]attr.Value, 0, len(lr.NonVersionTransitions))
			for _, t := range lr.NonVersionTransitions {
				nvtAttrs := map[string]attr.Value{
					"days":                     types.Int64Null(),
					"storage_class":            types.StringNull(),
					"is_access_time":           types.BoolNull(),
					"return_to_std_when_visit": types.BoolNull(),
				}
				nvtAttrs["days"] = types.Int64Value(int64(t.NoncurrentDays))
				nvtAttrs["storage_class"] = types.StringValue(string(t.StorageClass))
				if t.IsAccessTime != nil {
					nvtAttrs["is_access_time"] = types.BoolValue(*t.IsAccessTime)
				}
				if t.ReturnToStdWhenVisit != nil {
					nvtAttrs["return_to_std_when_visit"] = types.BoolValue(*t.ReturnToStdWhenVisit)
				}
				obj, d := types.ObjectValue(noncurrentVersionTransitionObjectElemType().AttrTypes, nvtAttrs)
				if d.HasError() {
					return model, fmt.Errorf("building noncurrent_version_transition object: %v", d)
				}
				nvtElems = append(nvtElems, obj)
			}
			s, d := types.SetValue(noncurrentVersionTransitionObjectElemType(), nvtElems)
			if d.HasError() {
				return model, fmt.Errorf("building noncurrent_version_transition set: %v", d)
			}
			nvtSet = s
		}

		// tags
		tagsVal := types.MapNull(types.StringType)
		if len(lr.Tags) != 0 {
			tagElems := make(map[string]attr.Value, len(lr.Tags))
			for _, t := range lr.Tags {
				tagElems[t.Key] = types.StringValue(t.Value)
			}
			m, d := types.MapValue(types.StringType, tagElems)
			if d.HasError() {
				return model, fmt.Errorf("building tags map: %v", d)
			}
			tagsVal = m
		}

		// filter
		filterVal := types.ListNull(filterObjectElemType())
		if lr.Filter != nil {
			fAttrs := map[string]attr.Value{
				"not":                      types.ListNull(filterNotObjectElemType()),
				"object_size_greater_than": types.Int64Null(),
				"object_size_less_than":    types.Int64Null(),
			}
			if len(lr.Filter.Not) != 0 {
				notElems := make([]attr.Value, 0, len(lr.Filter.Not))
				for _, n := range lr.Filter.Not {
					nAttrs := map[string]attr.Value{
						"prefix": types.StringNull(),
						"tag":    types.ListNull(filterNotTagObjectElemType()),
					}
					if n.Prefix != "" {
						nAttrs["prefix"] = types.StringValue(n.Prefix)
					}
					if n.Tag != nil {
						tagObj, d := types.ObjectValue(filterNotTagObjectElemType().AttrTypes, map[string]attr.Value{
							"key":   types.StringValue(n.Tag.Key),
							"value": types.StringValue(n.Tag.Value),
						})
						if d.HasError() {
							return model, fmt.Errorf("building filter.not.tag object: %v", d)
						}
						tagList, d := types.ListValue(filterNotTagObjectElemType(), []attr.Value{tagObj})
						if d.HasError() {
							return model, fmt.Errorf("building filter.not.tag list: %v", d)
						}
						nAttrs["tag"] = tagList
					}
					notObj, d := types.ObjectValue(filterNotObjectElemType().AttrTypes, nAttrs)
					if d.HasError() {
						return model, fmt.Errorf("building filter.not object: %v", d)
					}
					notElems = append(notElems, notObj)
				}
				notList, d := types.ListValue(filterNotObjectElemType(), notElems)
				if d.HasError() {
					return model, fmt.Errorf("building filter.not list: %v", d)
				}
				fAttrs["not"] = notList
			}
			if lr.Filter.ObjectSizeGreaterThan != nil {
				fAttrs["object_size_greater_than"] = types.Int64Value(*lr.Filter.ObjectSizeGreaterThan)
			}
			if lr.Filter.ObjectSizeLessThan != nil {
				fAttrs["object_size_less_than"] = types.Int64Value(*lr.Filter.ObjectSizeLessThan)
			}
			filterObj, d := types.ObjectValue(filterObjectElemType().AttrTypes, fAttrs)
			if d.HasError() {
				return model, fmt.Errorf("building filter object: %v", d)
			}
			filterVal, d = types.ListValue(filterObjectElemType(), []attr.Value{filterObj})
			if d.HasError() {
				return model, fmt.Errorf("building filter list: %v", d)
			}
		}

		// Assemble the rule object from explicitly typed attr values.
		idVal := types.StringNull()
		if lr.ID != "" {
			idVal = types.StringValue(lr.ID)
		}
		prefixVal := types.StringNull()
		if lr.Prefix != "" {
			prefixVal = types.StringValue(lr.Prefix)
		}
		ruleAttrs := map[string]attr.Value{
			"id":                            idVal,
			"prefix":                        prefixVal,
			"enabled":                       types.BoolValue(LifecycleRuleStatus(lr.Status) == ExpirationStatusEnabled),
			"expiration":                    expirationSet,
			"transitions":                   transitionsSet,
			"abort_multipart_upload":        amuSet,
			"noncurrent_version_expiration": nveSet,
			"noncurrent_version_transition": nvtSet,
			"tags":                          tagsVal,
			"filter":                        filterVal,
		}
		ruleObj, d := types.ObjectValue(lifecycleRuleObjectElemType().AttrTypes, ruleAttrs)
		if d.HasError() {
			return model, fmt.Errorf("building lifecycle_rule object: %v", d)
		}
		ruleElems = append(ruleElems, ruleObj)
	}

	l, d := types.ListValue(lifecycleRuleObjectElemType(), ruleElems)
	if d.HasError() {
		return model, fmt.Errorf("building lifecycle_rule list: %v", d)
	}
	model.LifecycleRule = l
	return model, nil
}

// ossLifecycleDateValidator mirrors validateOssBucketDateTimestamp for the
// plugin-framework, rejecting values that are not YYYY-MM-DD dates.
type ossLifecycleDateValidator struct{}

func (ossLifecycleDateValidator) Description(_ context.Context) string {
	return "date must be in YYYY-MM-DD format"
}

func (ossLifecycleDateValidator) MarkdownDescription(_ context.Context) string {
	return "date must be in YYYY-MM-DD format"
}

func (ossLifecycleDateValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}
	if _, err := time.Parse("2006-01-02", req.ConfigValue.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid OSS Lifecycle Date",
			fmt.Sprintf("%q cannot be parsed as date YYYY-MM-DD Format", req.ConfigValue.ValueString()),
		)
	}
}

// objectElemType helpers return the framework ObjectType for each nested block
// model. They are used with the non-generic types.SetValueFrom / ListValueFrom
// helpers so the read path can rebuild nested blocks from the oss-go-sdk
// structs; terraform-plugin-framework v1.19.0 does not provide the generic
// SetValueFrom[T] / ListValueFrom[T] helpers.
func expirationObjectElemType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"date":                         types.StringType,
		"created_before_date":          types.StringType,
		"days":                         types.Int64Type,
		"expired_object_delete_marker": types.BoolType,
	}}
}

func transitionsObjectElemType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"created_before_date":      types.StringType,
		"days":                     types.Int64Type,
		"storage_class":            types.StringType,
		"is_access_time":           types.BoolType,
		"return_to_std_when_visit": types.BoolType,
	}}
}

func abortMultipartUploadObjectElemType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"created_before_date": types.StringType,
		"days":                types.Int64Type,
	}}
}

func noncurrentVersionExpirationObjectElemType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"days": types.Int64Type,
	}}
}

func noncurrentVersionTransitionObjectElemType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"days":                     types.Int64Type,
		"storage_class":            types.StringType,
		"is_access_time":           types.BoolType,
		"return_to_std_when_visit": types.BoolType,
	}}
}

func filterNotTagObjectElemType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	}}
}

func filterNotObjectElemType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"prefix": types.StringType,
		"tag":    types.ListType{ElemType: filterNotTagObjectElemType()},
	}}
}

func filterObjectElemType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"not":                      types.ListType{ElemType: filterNotObjectElemType()},
		"object_size_greater_than": types.Int64Type,
		"object_size_less_than":    types.Int64Type,
	}}
}

func lifecycleRuleObjectElemType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":                            types.StringType,
		"prefix":                        types.StringType,
		"enabled":                       types.BoolType,
		"expiration":                    types.SetType{ElemType: expirationObjectElemType()},
		"transitions":                   types.SetType{ElemType: transitionsObjectElemType()},
		"abort_multipart_upload":        types.SetType{ElemType: abortMultipartUploadObjectElemType()},
		"noncurrent_version_expiration": types.SetType{ElemType: noncurrentVersionExpirationObjectElemType()},
		"noncurrent_version_transition": types.SetType{ElemType: noncurrentVersionTransitionObjectElemType()},
		"tags":                          types.MapType{ElemType: types.StringType},
		"filter":                        types.ListType{ElemType: filterObjectElemType()},
	}}
}
