package realtimecompute

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
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

const realtimeComputePackageTypeName = "alicloud_realtime_compute_package"

// realtimeComputePackageTimeout replaces the SDKv2 resource-level default
// timeouts (Create/Update/Delete each 5 minutes); terraform-plugin-framework
// does not support user-configurable resource timeouts yet.
const realtimeComputePackageTimeout = 5 * time.Minute

var (
	_ resource.Resource                = &realtimeComputePackageResource{}
	_ resource.ResourceWithImportState = &realtimeComputePackageResource{}
)

func NewRealtimeComputePackageResource() resource.Resource {
	return &realtimeComputePackageResource{}
}

type realtimeComputePackageResource struct {
	fwadapt.ResourceBase
}

// realtimeComputePackageResourceModel mirrors the schema. The id is the
// composite key "<projectName>:<packageName>".
type realtimeComputePackageResourceModel struct {
	Id          types.String `tfsdk:"id"`
	ProjectName types.String `tfsdk:"project_name"`
	PackageName types.String `tfsdk:"package_name"`
	OssBucket   types.String `tfsdk:"oss_bucket"`
	OssEndpoint types.String `tfsdk:"oss_endpoint"`
	OssOwner    types.String `tfsdk:"oss_owner"`
	OssPath     types.String `tfsdk:"oss_path"`
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`
	OriginName  types.String `tfsdk:"origin_name"`
	Md5         types.String `tfsdk:"md5"`
	Tag         types.String `tfsdk:"tag"`
	CreateTime  types.Int64  `tfsdk:"create_time"`
	Creator     types.String `tfsdk:"creator"`
	ModifyTime  types.Int64  `tfsdk:"modify_time"`
	Modifier    types.String `tfsdk:"modifier"`
}

func (r *realtimeComputePackageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = realtimeComputePackageTypeName
}

func (r *realtimeComputePackageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The name of the Realtime Compute (FOAS) project that owns the package. Part of the resource identity; changing it recreates the resource.",
			},
			"package_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The name of the package, unique within the project. Part of the resource identity; changing it recreates the resource.",
			},
			"oss_bucket": schema.StringAttribute{
				Required:    true,
				Description: "The OSS bucket that stores the package artifact.",
			},
			"oss_endpoint": schema.StringAttribute{
				Required:    true,
				Description: "The OSS endpoint used to access the bucket that stores the package artifact.",
			},
			"oss_owner": schema.StringAttribute{
				Required:    true,
				Description: "The owner of the OSS bucket that stores the package artifact.",
			},
			"oss_path": schema.StringAttribute{
				Required:    true,
				Description: "The OSS object path of the package artifact.",
			},
			"type": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("JAR", "DICTIONARY", "SCRIPT", "PYTHON"),
				},
				Description: "The package type. Valid values: `JAR`, `DICTIONARY`, `SCRIPT`, `PYTHON`. Changing it recreates the resource.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the package.",
			},
			"origin_name": schema.StringAttribute{
				Optional:    true,
				Description: "The alias (origin name) of the package artifact.",
			},
			"md5": schema.StringAttribute{
				Optional:    true,
				Description: "The MD5 checksum of the package artifact, used to validate the uploaded content.",
			},
			"tag": schema.StringAttribute{
				Optional:    true,
				Description: "The tag of the package.",
			},
			// create_time/creator/modify_time/modifier are API-managed computed
			// fields. As in the snapshot_lock Framework resource, they carry no
			// UseStateForUnknown planmodifier: modify_time/modifier follow
			// UpdatePackage server-side and the planned value could not survive an
			// apply, while create_time/creator are repopulated from GetPackage on
			// every refresh, so an unknown-in-plan never produces a real diff.
			"create_time": schema.Int64Attribute{
				Computed:    true,
				Description: "The creation time of the package.",
			},
			"creator": schema.StringAttribute{
				Computed:    true,
				Description: "The creator of the package.",
			},
			"modify_time": schema.Int64Attribute{
				Computed:    true,
				Description: "The latest modification time of the package.",
			},
			"modifier": schema.StringAttribute{
				Computed:    true,
				Description: "The latest modifier of the package.",
			},
		},
	}
}

func (r *realtimeComputePackageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan realtimeComputePackageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectName := plan.ProjectName.ValueString()
	action := fmt.Sprintf("/api/v2/projects/%s/packages", projectName)
	request := map[string]interface{}{
		"packageName": plan.PackageName.ValueString(),
		"ossBucket":   plan.OssBucket.ValueString(),
		"ossEndpoint": plan.OssEndpoint.ValueString(),
		"ossOwner":    plan.OssOwner.ValueString(),
		"ossPath":     plan.OssPath.ValueString(),
		"type":        plan.Type.ValueString(),
	}
	if !plan.Description.IsNull() {
		request["description"] = plan.Description.ValueString()
	}
	if !plan.OriginName.IsNull() {
		request["originName"] = plan.OriginName.ValueString()
	}
	if !plan.Md5.IsNull() {
		request["md5"] = plan.Md5.ValueString()
	}
	if !plan.Tag.IsNull() {
		request["tag"] = plan.Tag.ValueString()
	}
	query := map[string]*string{}
	header := map[string]*string{
		"RegionId": helper.StringPointer(r.Meta().RegionId),
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(realtimeComputePackageTimeout, func() *retry.RetryError {
		response, err = r.Meta().RoaPost("foas", foasPackageAPIVersion, action, query, header, request, true)
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
			fmt.Sprintf("Creating %s: calling CreatePackage", realtimeComputePackageTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, realtimeComputePackageTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	plan.Id = types.StringValue(fmt.Sprintf("%s:%s", projectName, plan.PackageName.ValueString()))

	object, err := NewService(r.Meta()).DescribeRealtimeComputePackage(plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetPackage", realtimeComputePackageTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, plan.Id.ValueString(), "GetPackage", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, realtimeComputePackageState(object, projectName, plan.PackageName.ValueString()))...)
}

func (r *realtimeComputePackageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state realtimeComputePackageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// After import only the id is set, so parse the composite identity keys
	// from it rather than relying on state fields that may still be unknown.
	projectName, packageName := parseRealtimeComputePackageID(state.Id.ValueString())

	object, err := NewService(r.Meta()).DescribeRealtimeComputePackage(state.Id.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetPackage", realtimeComputePackageTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), "GetPackage", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, realtimeComputePackageState(object, projectName, packageName))...)
}

func (r *realtimeComputePackageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state realtimeComputePackageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectName := state.ProjectName.ValueString()
	packageName := state.PackageName.ValueString()
	action := fmt.Sprintf("/api/v2/projects/%s/packages/%s", projectName, packageName)
	// The ACube operation mapping marks every updatable field of UpdatePackage
	// (Description/OssBucket/OssEndpoint/Md5/OriginName/OssOwner/OssPath/Tag) as
	// requestRequired:false, i.e. a partial update. Send all writable fields on
	// each update so partial and (hypothetical) full-replace semantics both leave
	// server state correct; the request is only issued when at least one writable
	// attribute actually changed.
	update := !plan.OssBucket.Equal(state.OssBucket) ||
		!plan.OssEndpoint.Equal(state.OssEndpoint) ||
		!plan.OssOwner.Equal(state.OssOwner) ||
		!plan.OssPath.Equal(state.OssPath) ||
		!plan.Description.Equal(state.Description) ||
		!plan.OriginName.Equal(state.OriginName) ||
		!plan.Md5.Equal(state.Md5) ||
		!plan.Tag.Equal(state.Tag)
	request := map[string]interface{}{}
	if update {
		request["ossBucket"] = plan.OssBucket.ValueString()
		request["ossEndpoint"] = plan.OssEndpoint.ValueString()
		request["ossOwner"] = plan.OssOwner.ValueString()
		request["ossPath"] = plan.OssPath.ValueString()
		request["description"] = plan.Description.ValueString()
		request["originName"] = plan.OriginName.ValueString()
		request["md5"] = plan.Md5.ValueString()
		request["tag"] = plan.Tag.ValueString()
	}

	if update {
		// Per the ACube operation mapping, RegionId is passed in the request
		// header (requestIn: header) on every package lifecycle API, including
		// UpdatePackage; keep it consistent with Create/Read/Delete rather than
		// placing it in the query string.
		query := map[string]*string{}
		header := map[string]*string{
			"RegionId": helper.StringPointer(r.Meta().RegionId),
		}
		var response map[string]interface{}
		var err error
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(realtimeComputePackageTimeout, func() *retry.RetryError {
			response, err = r.Meta().RoaPut("foas", foasPackageAPIVersion, action, query, header, request, true)
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
				fmt.Sprintf("Updating %s: calling UpdatePackage", realtimeComputePackageTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), action, errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
	}

	object, err := NewService(r.Meta()).DescribeRealtimeComputePackage(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetPackage", realtimeComputePackageTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), "GetPackage", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, realtimeComputePackageState(object, projectName, packageName))...)
}

func (r *realtimeComputePackageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state realtimeComputePackageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := fmt.Sprintf("/api/v2/projects/%s/packages/%s", state.ProjectName.ValueString(), state.PackageName.ValueString())
	query := map[string]*string{}
	header := map[string]*string{
		"RegionId": helper.StringPointer(r.Meta().RegionId),
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(realtimeComputePackageTimeout, func() *retry.RetryError {
		response, err = r.Meta().RoaDelete("foas", foasPackageAPIVersion, action, query, header, nil, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, nil)
	if err != nil {
		// DeletePackage is idempotent: a missing package is already gone.
		if errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeletePackage", realtimeComputePackageTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
}

func (r *realtimeComputePackageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// parseRealtimeComputePackageID splits the composite id "<projectName>:<packageName>"
// into its parts. A malformed id yields empty strings; the subsequent GetPackage call
// then surfaces the same invalid-id error DescribeRealtimeComputePackage would.
func parseRealtimeComputePackageID(id string) (string, string) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

// realtimeComputePackageState maps a GetPackage $.Package object onto the
// resource model. projectName and packageName are the resource identity keys;
// they are taken from the request state rather than the response because
// GetPackage may not echo ProjectName (it is only a path parameter), and a
// null identity field would collapse the composite id to ":<packageName>".
// The remaining fields use the FOAS 2018-11-11 GetPackage PascalCase property
// names defined by the ACube operation mapping.
func realtimeComputePackageState(object map[string]interface{}, projectName, packageName string) realtimeComputePackageResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	// The tea roa client may decode CreateTime/ModifyTime as json.Number, int
	// or float64; every numeric encoding the transport produces must map back
	// to Int64 or planned values would collapse to null.
	int64Val := func(key string) types.Int64 {
		switch v := object[key].(type) {
		case float64:
			return types.Int64Value(int64(v))
		case json.Number:
			if n, err := v.Int64(); err == nil {
				return types.Int64Value(n)
			}
		case int:
			return types.Int64Value(int64(v))
		case int64:
			return types.Int64Value(v)
		case string:
			if n, err := strconv.ParseInt(v, 10, 64); err == nil {
				return types.Int64Value(n)
			}
		}
		return types.Int64Null()
	}
	return realtimeComputePackageResourceModel{
		Id:          types.StringValue(fmt.Sprintf("%s:%s", projectName, packageName)),
		ProjectName: types.StringValue(projectName),
		PackageName: types.StringValue(packageName),
		OssBucket:   str("OssBucket"),
		OssEndpoint: str("OssEndpoint"),
		OssOwner:    str("OssOwner"),
		OssPath:     str("OssPath"),
		Type:        str("Type"),
		Description: str("Description"),
		OriginName:  str("OriginName"),
		Md5:         str("Md5"),
		Tag:         str("Tag"),
		CreateTime:  int64Val("CreateTime"),
		Creator:     str("Creator"),
		ModifyTime:  int64Val("ModifyTime"),
		Modifier:    str("Modifier"),
	}
}
