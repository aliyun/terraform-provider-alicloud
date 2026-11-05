package oss

import (
	"context"
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	ossBucketTagsTypeName = "alicloud_oss_bucket_tags"
)

var (
	_ resource.Resource                = &ossBucketTagsResource{}
	_ resource.ResourceWithConfigure   = &ossBucketTagsResource{}
	_ resource.ResourceWithImportState = &ossBucketTagsResource{}
)

// NewOssBucketTagsResource is the Framework factory for the
// alicloud_oss_bucket_tags resource, which manages the full tag set of an OSS
// bucket through the PutBucketTags / GetBucketTags / DeleteBucketTags API
// family (Oss 2019-05-17).
func NewOssBucketTagsResource() resource.Resource {
	return &ossBucketTagsResource{}
}

type ossBucketTagsResource struct {
	fwadapt.ResourceBase
}

type ossBucketTagsResourceModel struct {
	Id     types.String `tfsdk:"id"`
	Bucket types.String `tfsdk:"bucket"`
	Tags   types.Map    `tfsdk:"tags"`
}

func (r *ossBucketTagsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ossBucketTagsTypeName
}

func (r *ossBucketTagsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			// tags holds the complete set of tags (key/value) applied to the
			// bucket. PutBucketTags is an idempotent overwrite, so the resource
			// always reconciles the bucket's tag set to the configured map; an
			// empty map clears the tags via DeleteBucketTags.
			"tags": schema.MapAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
		},
	}
}

func (r *ossBucketTagsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ossBucketTagsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucket := plan.Bucket.ValueString()
	if err := r.applyBucketTags(bucket, tagsFromModel(ctx, plan.Tags)); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling PutBucketTags", ossBucketTagsTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, ossBucketTagsTypeName, "PutBucketTags", errs.AliyunOssGoSdk).Error(),
		)
		return
	}

	// Persist the resource identity before the read-back: the tag set exists
	// remotely from here on, so a failing GET must not leave Terraform without
	// a resource to taint and destroy.
	resp.Diagnostics.Append(resp.State.Set(ctx, ossBucketTagsResourceModel{
		Id:     types.StringValue(bucket),
		Bucket: plan.Bucket,
		Tags:   plan.Tags,
	})...)
	if resp.Diagnostics.HasError() {
		return
	}

	model, err := r.readBucketTagsModel(ctx, bucket)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetBucketTags", ossBucketTagsTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, ossBucketTagsTypeName, "GetBucketTags", errs.AliyunOssGoSdk).Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *ossBucketTagsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ossBucketTagsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// bucket is the resource identity; fall back to id for an imported
	// resource whose bucket attribute has not been populated yet.
	bucket := state.Bucket.ValueString()
	if bucket == "" {
		bucket = state.Id.ValueString()
	}

	model, err := r.readBucketTagsModel(ctx, bucket)
	if err != nil {
		// A missing bucket means the underlying object is gone; a bucket that
		// simply has no tag set is reported as an empty map inside the helper
		// and keeps the resource alive.
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetBucketTags", ossBucketTagsTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, ossBucketTagsTypeName, "GetBucketTags", errs.AliyunOssGoSdk).Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *ossBucketTagsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ossBucketTagsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state ossBucketTagsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucket := state.Bucket.ValueString()
	if err := r.applyBucketTags(bucket, tagsFromModel(ctx, plan.Tags)); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling PutBucketTags", ossBucketTagsTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, ossBucketTagsTypeName, "PutBucketTags", errs.AliyunOssGoSdk).Error(),
		)
		return
	}

	model, err := r.readBucketTagsModel(ctx, bucket)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetBucketTags", ossBucketTagsTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, ossBucketTagsTypeName, "GetBucketTags", errs.AliyunOssGoSdk).Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *ossBucketTagsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ossBucketTagsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.deleteBucketTags(state.Bucket.ValueString()); err != nil {
		if errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteBucketTags", ossBucketTagsTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, ossBucketTagsTypeName, "DeleteBucketTags", errs.AliyunOssGoSdk).Error(),
		)
		return
	}
}

func (r *ossBucketTagsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// applyBucketTags reconciles a bucket's tag set to tags. An empty map clears
// the tags via DeleteBucketTags; a non-empty map overwrites via SetBucketTagging
// (PutBucketTags is idempotent).
func (r *ossBucketTagsResource) applyBucketTags(bucket string, tags map[string]string) error {
	if len(tags) == 0 {
		return r.deleteBucketTags(bucket)
	}
	var tagging oss.Tagging
	for k, v := range tags {
		tagging.Tags = append(tagging.Tags, oss.Tag{Key: k, Value: v})
	}
	_, err := r.Meta().WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return nil, ossClient.SetBucketTagging(bucket, tagging)
	})
	logs.AddDebug("PutBucketTags", tagging, bucket)
	return err
}

// deleteBucketTags removes the whole tag set from a bucket via DeleteBucketTags.
func (r *ossBucketTagsResource) deleteBucketTags(bucket string) error {
	_, err := r.Meta().WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return nil, ossClient.DeleteBucketTagging(bucket)
	})
	logs.AddDebug("DeleteBucketTags", nil, bucket)
	return err
}

// describeBucketTags fetches the bucket's current tag set. A bucket that has no
// tag set returns NoSuchTagSet and is reported as an empty map (the bucket still
// exists, so the resource stays alive); other errors propagate to the caller.
func (r *ossBucketTagsResource) describeBucketTags(bucket string) (map[string]string, error) {
	var tagsMap map[string]string
	_, err := r.Meta().WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		result, e := ossClient.GetBucketTagging(bucket)
		if e != nil {
			if errs.IsExpectedErrors(e, []string{"NoSuchTagSet"}) {
				tagsMap = map[string]string{}
				return nil, nil
			}
			return nil, e
		}
		tagsMap = make(map[string]string, len(result.Tags))
		for _, t := range result.Tags {
			tagsMap[t.Key] = t.Value
		}
		return result, nil
	})
	if err == nil && tagsMap == nil {
		tagsMap = map[string]string{}
	}
	return tagsMap, err
}

// readBucketTagsModel builds the resource model from the bucket's current tag
// set; used by Create, Read and Update to populate state from a real read-back.
func (r *ossBucketTagsResource) readBucketTagsModel(ctx context.Context, bucket string) (ossBucketTagsResourceModel, error) {
	tags, err := r.describeBucketTags(bucket)
	if err != nil {
		return ossBucketTagsResourceModel{}, err
	}
	stateTags, diags := types.MapValueFrom(ctx, types.StringType, tags)
	if diags.HasError() {
		return ossBucketTagsResourceModel{}, fmt.Errorf("building tags state: %v", diags.Errors())
	}
	return ossBucketTagsResourceModel{
		Id:     types.StringValue(bucket),
		Bucket: types.StringValue(bucket),
		Tags:   stateTags,
	}, nil
}

// tagsFromModel converts a planned/configured tags map into a plain
// map[string]string for the OSS SDK; null or unknown maps yield an empty set
// (which clears the bucket's tags on apply).
func tagsFromModel(ctx context.Context, m types.Map) map[string]string {
	tags := make(map[string]string)
	if m.IsNull() || m.IsUnknown() {
		return tags
	}
	_ = m.ElementsAs(ctx, &tags, false)
	return tags
}
