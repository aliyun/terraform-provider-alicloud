package fc

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const pipelineTemplateTypeName = "alicloud_fc_pipeline_template"

// pipelineTemplateTimeout replaces the SDKv2 resource-level default timeouts
// (Create/Update/Delete each 5 minutes); terraform-plugin-framework does not
// support user-configurable resource timeouts yet.
const pipelineTemplateTimeout = 5 * time.Minute

// FC serverless deployment API (serverless::2021-09-24, ServiceCode=FC). The
// five PipelineTemplate REST APIs live on the FC gateway at
// /apis/serverlessdeployment/v1/pipelinetemplates, distinct from the FC
// 2023-03-30 surface used by the alicloud_fcv3_* resources.
const (
	pipelineTemplateAPIVersion = "2021-09-24"
	pipelineTemplateCollection = "/apis/serverlessdeployment/v1/pipelinetemplates"
)

var (
	_ resource.Resource                = &pipelineTemplateResource{}
	_ resource.ResourceWithImportState = &pipelineTemplateResource{}
)

func NewPipelineTemplateResource() resource.Resource {
	return &pipelineTemplateResource{}
}

type pipelineTemplateResource struct {
	fwadapt.ResourceBase
}

type pipelineTemplateResourceModel struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	Spec            types.String `tfsdk:"spec"`
	Labels          types.Map    `tfsdk:"labels"`
	Uid             types.String `tfsdk:"uid"`
	Kind            types.String `tfsdk:"kind"`
	CreatedTime     types.String `tfsdk:"created_time"`
	DeletionTime    types.String `tfsdk:"deletion_time"`
	Generation      types.Int64  `tfsdk:"generation"`
	ResourceVersion types.Int64  `tfsdk:"resource_version"`
}

func (r *pipelineTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = pipelineTemplateTypeName
}

func (r *pipelineTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a FC Pipeline Template resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource, same as `name`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the pipeline template. It must be unique within the account and cannot be changed. 1-40 characters, must start with a lowercase letter, contain only lowercase letters, digits and hyphens, and not end with a hyphen.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 40),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-z]([a-z0-9-]*[a-z0-9])?$`),
						"must start with a lowercase letter, contain only lowercase letters, digits and hyphens, and not end with a hyphen"),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the pipeline template.",
			},
			// spec is the PipelineTemplateSpec object (context, tasks,
			// contextSchema). Its context.data and per-task context.data are
			// free-form JSON, so the whole spec is carried as a canonical
			// JSON string: it is parsed before send and re-marshalled on
			// read so key ordering never produces a perpetual diff.
			"spec": schema.StringAttribute{
				Optional:    true,
				Description: "The pipeline template spec, encoded as a JSON object with `context`, `tasks` and optional `contextSchema`. Use jsonencode() in HCL.",
			},
			"labels": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "The labels of the pipeline template.",
			},
			"uid": schema.StringAttribute{
				Computed:    true,
				Description: "The UID of the account that owns the template, read-only.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"kind": schema.StringAttribute{
				Computed:    true,
				Description: "The kind of the resource, read-only.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the template, read-only.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"deletion_time": schema.StringAttribute{
				Computed:    true,
				Description: "The deletion time of the template, read-only.",
				// No UseStateForUnknown: set only on deletion, the value
				// moves outside apply control.
			},
			"generation": schema.Int64Attribute{
				Computed:    true,
				Description: "The generation of the spec, incremented on spec update, read-only.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"resource_version": schema.Int64Attribute{
				Computed:    true,
				Description: "The resource version, incremented on update, read-only.",
				// No UseStateForUnknown: bumps on every update, the planned
				// value can never survive an apply.
			},
		},
	}
}

func (r *pipelineTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan pipelineTemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, ok := buildPipelineTemplateBody(ctx, &plan, &resp.Diagnostics)
	if !ok {
		return
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(pipelineTemplateTimeout, func() *retry.RetryError {
		response, err = r.Meta().RoaPost("FC", pipelineTemplateAPIVersion, pipelineTemplateCollection, nil, nil, body, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug("CreatePipelineTemplate", response, body)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreatePipelineTemplate", pipelineTemplateTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreatePipelineTemplate", pipelineTemplateTypeName).Error(),
		)
		return
	}

	object, err := DescribePipelineTemplate(r.Meta(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetPipelineTemplate", pipelineTemplateTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetPipelineTemplate", pipelineTemplateTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreatePipelineTemplate", pipelineTemplateTypeName),
			fmt.Sprintf("The pipeline template %s is not visible via GetPipelineTemplate right after creation", plan.Name.ValueString()),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, pipelineTemplateState(ctx, object))...)
}

func (r *pipelineTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state pipelineTemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := DescribePipelineTemplate(r.Meta(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetPipelineTemplate", pipelineTemplateTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetPipelineTemplate", pipelineTemplateTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, pipelineTemplateState(ctx, object))...)
}

func (r *pipelineTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state pipelineTemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, ok := buildPipelineTemplateBody(ctx, &plan, &resp.Diagnostics)
	if !ok {
		return
	}
	action := fmt.Sprintf("%s/%s", pipelineTemplateCollection, state.Name.ValueString())

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(pipelineTemplateTimeout, func() *retry.RetryError {
		response, err = r.Meta().RoaPut("FC", pipelineTemplateAPIVersion, action, nil, nil, body, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug("PutPipelineTemplate", response, body)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling PutPipelineTemplate", pipelineTemplateTypeName),
			errs.WrapErrorf(err, "Updating %s: calling PutPipelineTemplate", pipelineTemplateTypeName).Error(),
		)
		return
	}

	object, err := DescribePipelineTemplate(r.Meta(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetPipelineTemplate", pipelineTemplateTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetPipelineTemplate", pipelineTemplateTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, pipelineTemplateState(ctx, object))...)
}

func (r *pipelineTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state pipelineTemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := fmt.Sprintf("%s/%s", pipelineTemplateCollection, state.Id.ValueString())
	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(pipelineTemplateTimeout, func() *retry.RetryError {
		response, err = r.Meta().RoaDelete("FC", pipelineTemplateAPIVersion, action, nil, nil, nil, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug("DeletePipelineTemplate", response, nil)
	if err != nil {
		// The FC gateway surfaces a missing template as HTTP 404; tolerate it
		// so a destroy of an already-gone resource is idempotent.
		if errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeletePipelineTemplate", pipelineTemplateTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeletePipelineTemplate", pipelineTemplateTypeName).Error(),
		)
		return
	}
}

func (r *pipelineTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// DescribePipelineTemplate queries GetPipelineTemplate by name and returns the
// PipelineTemplate object, or nil when it no longer exists. The FC gateway
// surfaces a missing template as HTTP 404, which errs.NotFoundError maps.
func DescribePipelineTemplate(client *connectivity.AliyunClient, name string) (map[string]interface{}, error) {
	response, err := client.RoaGet("FC", pipelineTemplateAPIVersion, fmt.Sprintf("%s/%s", pipelineTemplateCollection, name), nil, nil, nil)
	if err != nil {
		if errs.NotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return response, nil
}

// buildPipelineTemplateBody constructs the PipelineTemplate request body from
// the planned model. spec is parsed from its JSON string so the API receives a
// real object; labels is expanded into a plain map. Returns ok=false when the
// plan carried malformed JSON or labels the framework could not decode.
func buildPipelineTemplateBody(ctx context.Context, plan *pipelineTemplateResourceModel, diags *diag.Diagnostics) (map[string]interface{}, bool) {
	body := map[string]interface{}{
		"name": plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() && plan.Description.ValueString() != "" {
		body["description"] = plan.Description.ValueString()
	}
	if !plan.Spec.IsNull() && plan.Spec.ValueString() != "" {
		var specObj interface{}
		if err := json.Unmarshal([]byte(plan.Spec.ValueString()), &specObj); err != nil {
			diags.AddError(
				fmt.Sprintf("Building %s request body: parsing spec", pipelineTemplateTypeName),
				fmt.Sprintf("spec is not valid JSON: %v", err),
			)
			return nil, false
		}
		body["spec"] = specObj
	}
	if !plan.Labels.IsNull() && len(plan.Labels.Elements()) > 0 {
		labelsMap := make(map[string]string, len(plan.Labels.Elements()))
		diags.Append(plan.Labels.ElementsAs(ctx, &labelsMap, false)...)
		if diags.HasError() {
			return nil, false
		}
		lmap := make(map[string]interface{}, len(labelsMap))
		for k, v := range labelsMap {
			lmap[k] = v
		}
		body["labels"] = lmap
	}
	return body, true
}

// pipelineTemplateState maps the GetPipelineTemplate response object onto the
// resource model. spec is re-marshalled to canonical (sorted-key) JSON so
// that key ordering in the API response never produces a perpetual diff.
func pipelineTemplateState(ctx context.Context, object map[string]interface{}) pipelineTemplateResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
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
	m := pipelineTemplateResourceModel{
		Id:              str("name"),
		Name:            str("name"),
		Description:     str("description"),
		Uid:             str("uid"),
		Kind:            str("kind"),
		CreatedTime:     str("createdTime"),
		DeletionTime:    str("deletionTime"),
		Generation:      int64Val("generation"),
		ResourceVersion: int64Val("resourceVersion"),
	}
	if spec, ok := object["spec"].(map[string]interface{}); ok && len(spec) > 0 {
		if b, err := json.Marshal(spec); err == nil {
			m.Spec = types.StringValue(string(b))
		}
	}
	if labels, ok := object["labels"].(map[string]interface{}); ok && len(labels) > 0 {
		lmap := make(map[string]string, len(labels))
		for k, v := range labels {
			if s, ok := v.(string); ok {
				lmap[k] = s
			} else {
				lmap[k] = fmt.Sprintf("%v", v)
			}
		}
		if lv, dgs := types.MapValueFrom(ctx, types.StringType, lmap); !dgs.HasError() {
			m.Labels = lv
		}
	}
	return m
}
