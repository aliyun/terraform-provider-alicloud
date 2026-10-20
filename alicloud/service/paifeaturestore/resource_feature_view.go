// Package paifeaturestore holds the Terraform Plugin Framework resources for
// the Alibaba Cloud PaiFeatureStore (PAI Feature Store) product. The product
// exposes a RESTful/ROA API on version 2023-06-21; there is no Go SDK, so the
// resources call the generic AliyunClient.Roa* helpers directly.
package paifeaturestore

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	featureViewTypeName = "alicloud_paifeaturestore_feature_view"
	paifeatureStoreAPI  = "PaiFeatureStore"
	paifeatureStoreVer  = "2023-06-21"
	idSeparator         = ":"
)

var (
	_ resource.Resource                = &featureViewResource{}
	_ resource.ResourceWithImportState = &featureViewResource{}
	_ resource.ResourceWithModifyPlan  = &featureViewResource{}
)

// NewFeatureViewResource returns the alicloud_paifeaturestore_feature_view resource.
func NewFeatureViewResource() resource.Resource {
	return &featureViewResource{}
}

type featureViewResource struct {
	fwadapt.ResourceBase
}

// featureViewResourceModel is the Terraform state shape. The composite identity
// is "InstanceId:FeatureViewId": InstanceId is the parent FeatureStore instance,
// FeatureViewId the primary key the API assigns at create.
type featureViewResourceModel struct {
	Id                   types.String            `tfsdk:"id"`
	InstanceId           types.String            `tfsdk:"instance_id"`
	FeatureViewId        types.String            `tfsdk:"feature_view_id"`
	FeatureViewName      types.String            `tfsdk:"feature_view_name"`
	ProjectId            types.String            `tfsdk:"project_id"`
	FeatureEntityId      types.String            `tfsdk:"feature_entity_id"`
	Type                 types.String            `tfsdk:"type"`
	WriteMethod          types.String            `tfsdk:"write_method"`
	SyncOnlineTable      types.Bool              `tfsdk:"sync_online_table"`
	Ttl                  types.Int64             `tfsdk:"ttl"`
	WriteToFeatureDb     types.Bool              `tfsdk:"write_to_feature_db"`
	RegisterTable        types.String            `tfsdk:"register_table"`
	RegisterDatasourceId types.String            `tfsdk:"register_datasource_id"`
	Config               types.String            `tfsdk:"config"`
	Fields               []featureViewFieldModel `tfsdk:"fields"`
	Tags                 types.Set               `tfsdk:"tags"`
	RegionId             types.String            `tfsdk:"region_id"`
	CreateTime           types.String            `tfsdk:"create_time"`
}

type featureViewFieldModel struct {
	Name       types.String                `tfsdk:"name"`
	Type       types.String                `tfsdk:"type"`
	Attributes types.Set                   `tfsdk:"attributes"`
	Transform  []featureViewTransformModel `tfsdk:"transform"`
	Dimension  types.Int64                 `tfsdk:"dimension"`
}

type featureViewTransformModel struct {
	Type        types.String                     `tfsdk:"type"`
	LLMConfigId types.Int64                      `tfsdk:"llm_config_id"`
	Input       []featureViewTransformInputModel `tfsdk:"input"`
}

type featureViewTransformInputModel struct {
	Name     types.String `tfsdk:"name"`
	Type     types.String `tfsdk:"type"`
	Modality types.String `tfsdk:"modality"`
}

func (r *featureViewResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = featureViewTypeName
}

func (r *featureViewResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PaiFeatureStore Feature View resource. " +
			"A feature view belongs to a FeatureStore instance and project and maps an offline feature table to an online one.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource, in the format `<InstanceId>:<FeatureViewId>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the FeatureStore instance. It is the parent of the feature view.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"feature_view_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the feature view, assigned by the API at create time.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"feature_view_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the feature view. It must start with a letter or a number and can contain only letters, numbers, and underscores, up to 128 characters.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the FeatureStore project that the feature view belongs to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"feature_entity_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The ID of the feature entity.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the feature view. Valid values: `Batch`, `Stream`, `Sequence`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"write_method": schema.StringAttribute{
				Required:    true,
				Description: "The method of writing data to the online table. Valid values: `ByReadyMadeTable`, `Custom`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"sync_online_table": schema.BoolAttribute{
				Required:    true,
				Description: "Whether to synchronize the online feature table.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"ttl": schema.Int64Attribute{
				Required:    true,
				Description: "The data lifecycle, in seconds.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"write_to_feature_db": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether data is written to FeatureDB. Defaults to `false`.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolplanmodifier.RequiresReplace(),
				},
			},
			"register_table": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The offline registry table name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"register_datasource_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The ID of the data source where the registry table resides.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"config": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Advanced configuration, a JSON-format string.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			// Tags are unordered in the API model, so they are a Set: a List would
			// make element order provider-visible, and a server-side reorder would
			// then plan a spurious replacement.
			"tags": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "The tag list.",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
					setplanmodifier.RequiresReplace(),
				},
			},
			"region_id": schema.StringAttribute{
				Computed:    true,
				Description: "The region ID of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"fields": schema.ListNestedBlock{
				Description: "The list of feature view fields.",
				// The UpdateFeatureView API accepts only a subset of the field
				// sub-properties; rather than ship a silently lossy partial update,
				// the identifying field sub-attributes carry RequiresReplace so any
				// change to them forces replacement until a later revision implements
				// the mapped update subset precisely. Protocol v5 (tf5muxserver)
				// cannot express nested attributes, so fields/transform/input are
				// nested blocks (NestedBlockObject), which convert cleanly to v5.
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "The name of the feature field.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"type": schema.StringAttribute{
							Required:    true,
							Description: "The type of the feature field. Valid values: `INT32`, `INT64`, `FLOAT`, `DOUBLE`, `STRING`, `BOOLEAN`, `TIMESTAMP`, and the complex types `ARRAY` and `MAP`.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						// Field attributes are unordered in the API model, so they use a
						// Set for the same reason as tags.
						"attributes": schema.SetAttribute{
							Optional:    true,
							Computed:    true,
							ElementType: types.StringType,
							Description: "The attributes of the field. Valid values: `Partition`, `PrimaryKey`, `EventTime`.",
							PlanModifiers: []planmodifier.Set{
								setplanmodifier.UseStateForUnknown(),
								setplanmodifier.RequiresReplace(),
							},
						},
						"dimension": schema.Int64Attribute{
							Optional:    true,
							Computed:    true,
							Description: "The dimension.",
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
								int64planmodifier.RequiresReplace(),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"transform": schema.ListNestedBlock{
							Description: "The feature production configuration.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: "The feature production type. Valid value: `LLMEmbedding`.",
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
											stringplanmodifier.RequiresReplace(),
										},
									},
									"llm_config_id": schema.Int64Attribute{
										Optional:    true,
										Computed:    true,
										Description: "The LLM configuration ID.",
										PlanModifiers: []planmodifier.Int64{
											int64planmodifier.UseStateForUnknown(),
											int64planmodifier.RequiresReplace(),
										},
									},
								},
								Blocks: map[string]schema.Block{
									"input": schema.ListNestedBlock{
										Description: "The feature production inputs.",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"name": schema.StringAttribute{
													Optional:    true,
													Computed:    true,
													Description: "The input field name.",
													PlanModifiers: []planmodifier.String{
														stringplanmodifier.UseStateForUnknown(),
														stringplanmodifier.RequiresReplace(),
													},
												},
												"type": schema.StringAttribute{
													Optional:    true,
													Computed:    true,
													Description: "The input field type.",
													PlanModifiers: []planmodifier.String{
														stringplanmodifier.UseStateForUnknown(),
														stringplanmodifier.RequiresReplace(),
													},
												},
												"modality": schema.StringAttribute{
													Optional:    true,
													Computed:    true,
													Description: "The input field modality. Valid values: `NONE`, `TEXT`, `IMAGE`, `VIDEO`, `MULTI_IMAGES`.",
													PlanModifiers: []planmodifier.String{
														stringplanmodifier.UseStateForUnknown(),
														stringplanmodifier.RequiresReplace(),
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

func (r *featureViewResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan featureViewResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{
		"Name":            plan.FeatureViewName.ValueString(),
		"ProjectId":       plan.ProjectId.ValueString(),
		"Type":            plan.Type.ValueString(),
		"WriteMethod":     plan.WriteMethod.ValueString(),
		"SyncOnlineTable": plan.SyncOnlineTable.ValueBool(),
		"TTL":             plan.Ttl.ValueInt64(),
		"Fields":          fieldsToBody(plan.Fields),
	}
	if !plan.FeatureEntityId.IsNull() && !plan.FeatureEntityId.IsUnknown() {
		body["FeatureEntityId"] = plan.FeatureEntityId.ValueString()
	}
	if !plan.WriteToFeatureDb.IsNull() && !plan.WriteToFeatureDb.IsUnknown() {
		body["WriteToFeatureDB"] = plan.WriteToFeatureDb.ValueBool()
	}
	if !plan.RegisterTable.IsNull() && !plan.RegisterTable.IsUnknown() {
		body["RegisterTable"] = plan.RegisterTable.ValueString()
	}
	if !plan.RegisterDatasourceId.IsNull() && !plan.RegisterDatasourceId.IsUnknown() {
		body["RegisterDatasourceId"] = plan.RegisterDatasourceId.ValueString()
	}
	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		body["Config"] = plan.Config.ValueString()
	}
	if tags := stringValues(plan.Tags); len(tags) > 0 {
		body["Tags"] = tags
	}

	pathName := fmt.Sprintf("/api/v1/instances/%s/featureviews", plan.InstanceId.ValueString())
	response, err := r.Client().RoaPost(paifeatureStoreAPI, paifeatureStoreVer, pathName, nil, nil, body, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateFeatureView", featureViewTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateFeatureView", featureViewTypeName).Error(),
		)
		return
	}

	featureViewId := stringFromResponse(response, "FeatureViewId")
	if featureViewId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateFeatureView", featureViewTypeName),
			fmt.Sprintf("The CreateFeatureView response carried no FeatureViewId: %v", response),
		)
		return
	}

	// Mirror the ims Create flow: after the API returns the primary key, re-read
	// the resource and set state from the server response. Setting the plan
	// directly would leave Optional+Computed attributes null and trigger an
	// "inconsistent result after apply"; a fresh GetFeatureView populates every
	// attribute (echoed inputs, computed timestamps, server defaults) uniformly.
	object, err := getFeatureView(r.Client(), plan.InstanceId.ValueString(), featureViewId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling GetFeatureView", featureViewTypeName),
			errs.WrapErrorf(err, "Creating %s: calling GetFeatureView", featureViewTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateFeatureView", featureViewTypeName),
			fmt.Sprintf("The feature view %s is not visible via GetFeatureView right after creation", featureViewId),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, featureViewState(object, r.Client().RegionId, plan.InstanceId.ValueString(), featureViewId))...)
}

func (r *featureViewResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state featureViewResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId, featureViewId, ok := parseFeatureViewID(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: parsing resource id", featureViewTypeName),
			fmt.Sprintf("The resource id %q is not in the form `<InstanceId>:<FeatureViewId>`.", state.Id.ValueString()),
		)
		return
	}

	object, err := getFeatureView(r.Client(), instanceId, featureViewId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetFeatureView", featureViewTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetFeatureView", featureViewTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, featureViewState(object, r.Client().RegionId, instanceId, featureViewId))...)
}

// Update is never expected to be reached: every input attribute carries a
// RequiresReplace plan modifier and ModifyPlan forces replacement when a nested
// block element is removed, so any change plans a replacement. If the framework
// calls Update anyway, failing loudly is safer than writing the plan to state:
// the change was never sent to the API, so state would silently drift from the
// server on the next Read.
func (r *featureViewResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		fmt.Sprintf("Updating %s", featureViewTypeName),
		"The resource does not support in-place updates; every input change must replace the resource. "+
			"Reaching Update means the plan did not request replacement for the pending diff.",
	)
}

// ModifyPlan closes the gap left by nested RequiresReplace modifiers: the
// framework applies plan modifiers per planned list element, so removing a
// fields, transform or input element never visits the removed element's plan
// modifiers and Terraform would otherwise plan an in-place update. Comparing
// element counts at every nested level keeps the documented "any input change
// recreates the resource" contract.
func (r *featureViewResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	var plan, state featureViewResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if fieldsDroppedElements(plan.Fields, state.Fields) {
		resp.RequiresReplace = append(resp.RequiresReplace, path.Root("fields"))
	}
}

// fieldsDroppedElements reports whether the planned fields list removes an
// element at any nesting level relative to state. The lists are ordered, so
// index-wise comparison is well defined; a reorder surfaces through the nested
// attribute modifiers because the elements at the shifted indexes differ.
func fieldsDroppedElements(plan, state []featureViewFieldModel) bool {
	if len(plan) != len(state) {
		return true
	}
	for i := range plan {
		if len(plan[i].Transform) != len(state[i].Transform) {
			return true
		}
		for j := range plan[i].Transform {
			if len(plan[i].Transform[j].Input) != len(state[i].Transform[j].Input) {
				return true
			}
		}
	}
	return false
}

func (r *featureViewResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state featureViewResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId, featureViewId, ok := parseFeatureViewID(state.Id.ValueString())
	if !ok {
		// An unparsable id cannot be deleted server-side; clearing local state is
		// the only sane action, matching Read's treatment of a malformed id.
		return
	}

	pathName := fmt.Sprintf("/api/v1/instances/%s/featureviews/%s", instanceId, featureViewId)
	_, err := r.Client().RoaDelete(paifeatureStoreAPI, paifeatureStoreVer, pathName, nil, nil, nil, true)
	if err != nil && !isFeatureViewNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteFeatureView", featureViewTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteFeatureView", featureViewTypeName).Error(),
		)
		return
	}
}

func (r *featureViewResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getFeatureView calls GetFeatureView and returns the feature view object, or
// nil when it no longer exists. The API reports a missing resource with the
// DoConfigCenterRequestErrorProblem code (per the CloudSpec resourceNotExistCondition).
func getFeatureView(client *connectivity.AliyunClient, instanceId, featureViewId string) (map[string]interface{}, error) {
	pathName := fmt.Sprintf("/api/v1/instances/%s/featureviews/%s", instanceId, featureViewId)
	response, err := client.RoaGet(paifeatureStoreAPI, paifeatureStoreVer, pathName, nil, nil, nil)
	if err != nil {
		if isFeatureViewNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return response, nil
}

// featureViewState maps a GetFeatureView response onto the resource model. Empty
// strings become null, matching the CloudSpec $exists defaults (e.g.
// WriteToFeatureDB defaults to false when absent).
func featureViewState(object map[string]interface{}, regionId, instanceId, featureViewId string) *featureViewResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	boolVal := func(key string, fallback bool) types.Bool {
		if v, ok := object[key].(bool); ok {
			return types.BoolValue(v)
		}
		return types.BoolValue(fallback)
	}
	return &featureViewResourceModel{
		Id:                   types.StringValue(buildFeatureViewID(instanceId, featureViewId)),
		InstanceId:           types.StringValue(instanceId),
		FeatureViewId:        types.StringValue(featureViewId),
		FeatureViewName:      str("Name"),
		ProjectId:            str("ProjectId"),
		FeatureEntityId:      str("FeatureEntityId"),
		Type:                 str("Type"),
		WriteMethod:          str("WriteMethod"),
		SyncOnlineTable:      boolVal("SyncOnlineTable", false),
		Ttl:                  int64FromResponse(object, "TTL"),
		WriteToFeatureDb:     boolVal("WriteToFeatureDB", false),
		RegisterTable:        str("RegisterTable"),
		RegisterDatasourceId: str("RegisterDatasourceId"),
		Config:               str("Config"),
		Fields:               fieldsFromResponse(object["Fields"]),
		Tags:                 stringSetFromResponse(object["Tags"]),
		RegionId:             types.StringValue(regionId),
		CreateTime:           str("GmtCreateTime"),
	}
}

// fieldsToBody converts the planned fields onto the CreateFeatureView request body.
func fieldsToBody(fields []featureViewFieldModel) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(fields))
	for _, f := range fields {
		item := map[string]interface{}{
			"Name": f.Name.ValueString(),
			"Type": f.Type.ValueString(),
		}
		if attrs := stringValues(f.Attributes); len(attrs) > 0 {
			item["Attributes"] = attrs
		}
		if len(f.Transform) > 0 {
			item["Transform"] = transformToBody(f.Transform)
		}
		if !f.Dimension.IsNull() && !f.Dimension.IsUnknown() {
			item["Dimension"] = f.Dimension.ValueInt64()
		}
		out = append(out, item)
	}
	return out
}

func transformToBody(transforms []featureViewTransformModel) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(transforms))
	for _, t := range transforms {
		item := map[string]interface{}{}
		if !t.Type.IsNull() && !t.Type.IsUnknown() {
			item["Type"] = t.Type.ValueString()
		}
		if !t.LLMConfigId.IsNull() && !t.LLMConfigId.IsUnknown() {
			item["LLMConfigId"] = t.LLMConfigId.ValueInt64()
		}
		if len(t.Input) > 0 {
			inputs := make([]map[string]interface{}, 0, len(t.Input))
			for _, in := range t.Input {
				entry := map[string]interface{}{}
				if !in.Name.IsNull() && !in.Name.IsUnknown() {
					entry["Name"] = in.Name.ValueString()
				}
				if !in.Type.IsNull() && !in.Type.IsUnknown() {
					entry["Type"] = in.Type.ValueString()
				}
				if !in.Modality.IsNull() && !in.Modality.IsUnknown() {
					entry["Modality"] = in.Modality.ValueString()
				}
				inputs = append(inputs, entry)
			}
			item["Input"] = inputs
		}
		out = append(out, item)
	}
	return out
}

// fieldsFromResponse rebuilds the fields model from the GetFeatureView response.
// The Fields list is ordered, so the slice preserves the API's element order in
// state (the unordered distinction in the model applies to Attributes only).
func fieldsFromResponse(raw interface{}) []featureViewFieldModel {
	arr, ok := raw.([]interface{})
	if !ok || len(arr) == 0 {
		return nil
	}
	out := make([]featureViewFieldModel, 0, len(arr))
	for _, el := range arr {
		field, ok := el.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, featureViewFieldModel{
			Name:       stringAttr(field, "Name"),
			Type:       stringAttr(field, "Type"),
			Attributes: stringSetFromResponse(field["Attributes"]),
			Transform:  transformFromResponse(field["Transform"]),
			Dimension:  int64Attr(field, "Dimension"),
		})
	}
	return out
}

func transformFromResponse(raw interface{}) []featureViewTransformModel {
	arr, ok := raw.([]interface{})
	if !ok || len(arr) == 0 {
		return nil
	}
	out := make([]featureViewTransformModel, 0, len(arr))
	for _, el := range arr {
		t, ok := el.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, featureViewTransformModel{
			Type:        stringAttr(t, "Type"),
			LLMConfigId: int64Attr(t, "LLMConfigId"),
			Input:       inputFromResponse(t["Input"]),
		})
	}
	return out
}

func inputFromResponse(raw interface{}) []featureViewTransformInputModel {
	arr, ok := raw.([]interface{})
	if !ok || len(arr) == 0 {
		return nil
	}
	out := make([]featureViewTransformInputModel, 0, len(arr))
	for _, el := range arr {
		in, ok := el.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, featureViewTransformInputModel{
			Name:     stringAttr(in, "Name"),
			Type:     stringAttr(in, "Type"),
			Modality: stringAttr(in, "Modality"),
		})
	}
	return out
}

// buildFeatureViewID composes the composite identity "<InstanceId>:<FeatureViewId>".
func buildFeatureViewID(instanceId, featureViewId string) string {
	return instanceId + idSeparator + featureViewId
}

// parseFeatureViewID splits the composite identity. It reports ok=false when the
// id does not contain exactly one separator.
func parseFeatureViewID(id string) (instanceId, featureViewId string, ok bool) {
	parts := strings.SplitN(id, idSeparator, 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// isFeatureViewNotExist reports whether err is the PaiFeatureStore resource
// not-found condition. The CloudSpec model declares exactly
// resourceNotExistErrorCodes: ["DoConfigCenterRequestErrorProblem"], so only
// that typed SDK error code clears the resource from state; any other error
// (permission, throttling, internal) must surface instead of silently removing
// a resource that still exists.
func isFeatureViewNotExist(err error) bool {
	e, ok := err.(*tea.SDKError)
	return ok && e.Code != nil && *e.Code == "DoConfigCenterRequestErrorProblem"
}

// stringValues extracts the []string payload from a known (non-unknown, non-null)
// set of strings for the request body. Computed sets that are still unknown
// during planning are tolerated by returning nil.
func stringValues(in types.Set) []string {
	if in.IsNull() || in.IsUnknown() {
		return nil
	}
	elements := in.Elements()
	out := make([]string, 0, len(elements))
	for _, v := range elements {
		if s, ok := v.(types.String); ok && !s.IsNull() {
			out = append(out, s.ValueString())
		}
	}
	return out
}

func stringFromResponse(object map[string]interface{}, key string) string {
	if v, ok := object[key].(string); ok {
		return v
	}
	return ""
}

func stringAttr(object map[string]interface{}, key string) types.String {
	if v, ok := object[key].(string); ok && v != "" {
		return types.StringValue(v)
	}
	return types.StringNull()
}

// int64FromResponse reads a numeric response field. The tea ROA client decodes
// JSON numbers into json.Number; int64Attr therefore handles json.Number as
// well as the plain float64/int/int64 Go types.
func int64FromResponse(object map[string]interface{}, key string) types.Int64 {
	return int64Attr(object, key)
}

func int64Attr(object map[string]interface{}, key string) types.Int64 {
	switch v := object[key].(type) {
	case json.Number:
		if n, err := v.Int64(); err == nil {
			return types.Int64Value(n)
		}
	case float64:
		return types.Int64Value(int64(v))
	case int:
		return types.Int64Value(int64(v))
	case int64:
		return types.Int64Value(v)
	}
	return types.Int64Null()
}

func stringSetFromResponse(raw interface{}) types.Set {
	arr, ok := raw.([]interface{})
	if !ok || len(arr) == 0 {
		return types.SetNull(types.StringType)
	}
	elements := make([]attr.Value, 0, len(arr))
	for _, el := range arr {
		if s, ok := el.(string); ok && s != "" {
			elements = append(elements, types.StringValue(s))
		}
	}
	if len(elements) == 0 {
		return types.SetNull(types.StringType)
	}
	s, diags := types.SetValue(types.StringType, elements)
	if diags.HasError() {
		return types.SetNull(types.StringType)
	}
	return s
}
