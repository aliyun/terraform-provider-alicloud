package dataworks

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const parameterTypeName = "alicloud_data_works_parameter"

var (
	_ resource.Resource                = &parameterResource{}
	_ resource.ResourceWithImportState = &parameterResource{}
)

// parameterPropertyAttrType is the element attr.Type of the properties set. It mirrors the
// default object type derived from the SetNestedAttribute NestedObject so that
// types.SetValueFrom can build state values without re-deriving the schema.
var parameterPropertyAttrType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"env_type": types.StringType,
		"value":    types.StringType,
	},
}

func NewParameterResource() resource.Resource {
	return &parameterResource{}
}

type parameterResource struct {
	fwadapt.ResourceBase
}

type parameterPropertyModel struct {
	EnvType types.String `tfsdk:"env_type"`
	Value   types.String `tfsdk:"value"`
}

type parameterResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	ProjectId   types.String `tfsdk:"project_id"`
	Scope       types.String `tfsdk:"scope"`
	Owner       types.String `tfsdk:"owner"`
	Description types.String `tfsdk:"description"`
	Properties  types.Set    `tfsdk:"properties"`
	CreateTime  types.String `tfsdk:"create_time"`
	CreateUser  types.String `tfsdk:"create_user"`
	ModifyTime  types.String `tfsdk:"modify_time"`
	ModifyUser  types.String `tfsdk:"modify_user"`
	Version     types.String `tfsdk:"version"`
}

func (r *parameterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = parameterTypeName
}

func (r *parameterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a DataWorks Parameter resource.",
		Attributes: map[string]schema.Attribute{
			"properties": schema.SetNestedAttribute{
				Required:            true,
				MarkdownDescription: "The value configuration of the parameter. The production environment (Prod) is required. If duplicate environment configurations appear in the array, the later configuration does not take effect.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"env_type": schema.StringAttribute{
							Required: true,
							Description: "The project environment. Valid values: `Prod` (production) and `Dev` (development). " +
								"A Prod configuration is required for `type = SecretConstant`.",
							Validators: []validator.String{
								stringvalidator.OneOf("Prod", "Dev"),
							},
						},
						"value": schema.StringAttribute{
							Required: true,
							Description: "The parameter value. Only Chinese characters, letters, digits and some special characters " +
								"(/, :, ., [, ], ,, \\, ', \", _, =, ?, space, carriage return, line feed, +, -, *, %, &, @, !, $, #, {, }) are allowed.",
						},
					},
				},
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the DataWorks Parameter.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The parameter name, unique within the workspace. It must start with `workspace.`. The remaining part must start with a letter and can only contain letters, underscores and digits. The total length cannot exceed 255 characters.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the parameter. Valid values: `PlainConstant` (constant, plaintext), `SecretConstant` (constant, ciphertext) and `Variable` (variable).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("PlainConstant", "SecretConstant", "Variable"),
				},
			},
			"project_id": schema.StringAttribute{
				Required:    true,
				Description: "The workspace ID. Required when `scope = Project`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"scope": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The scope of the parameter. Valid value: `Project`. Defaults to `Project`. Currently no other type can be selected.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("Project"),
				},
			},
			"owner": schema.StringAttribute{
				Required:    true,
				Description: "The account ID of the responsible person.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the parameter.",
				PlanModifiers: []planmodifier.String{
					parameterEmptyStringToNullPlanModifier{},
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time. A 13-digit millisecond timestamp, for example `1640000000000`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_user": schema.StringAttribute{
				Computed:    true,
				Description: "The creator account ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modify_time": schema.StringAttribute{
				Computed:    true,
				Description: "The last modification time. A 13-digit millisecond timestamp, for example `1640000000000`.",
			},
			"modify_user": schema.StringAttribute{
				Computed:    true,
				Description: "The last modifier account ID.",
			},
			"version": schema.StringAttribute{
				Computed:    true,
				Description: "The version number.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *parameterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan parameterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	properties, err := parameterPropertiesFromSet(ctx, plan.Properties)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading properties", parameterTypeName),
			err.Error(),
		)
		return
	}

	propertiesJson, err := parameterPropertiesRequest(properties)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: marshalling properties", parameterTypeName),
			err.Error(),
		)
		return
	}

	request := map[string]interface{}{
		"RegionId":   r.Meta().RegionId,
		"Name":       plan.Name.ValueString(),
		"Type":       plan.Type.ValueString(),
		"ProjectId":  plan.ProjectId.ValueString(),
		"Owner":      plan.Owner.ValueString(),
		"Properties": propertiesJson,
	}
	if !plan.Scope.IsNull() && !plan.Scope.IsUnknown() {
		request["Scope"] = plan.Scope.ValueString()
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		request["Description"] = plan.Description.ValueString()
	}

	response, err := r.Meta().RpcPost("dataworks-public", "2024-05-18", "CreateParameter", nil, request, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateParameter", parameterTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateParameter", parameterTypeName).Error(),
		)
		return
	}

	id := fmt.Sprint(response["Id"])
	if id == "" || id == "<nil>" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateParameter", parameterTypeName),
			fmt.Sprintf("The response carried no Id: %v", response),
		)
		return
	}

	object, err := getParameter(r.Meta(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetParameter", parameterTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetParameter", parameterTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateParameter", parameterTypeName),
			fmt.Sprintf("The parameter %s is not visible via GetParameter right after creation", id),
		)
		return
	}
	state, stateDiags := parameterState(ctx, object)
	resp.Diagnostics.Append(stateDiags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *parameterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state parameterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := getParameter(r.Meta(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetParameter", parameterTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetParameter", parameterTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	newState, stateDiags := parameterState(ctx, object)
	resp.Diagnostics.Append(stateDiags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *parameterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state parameterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	properties, err := parameterPropertiesFromSet(ctx, plan.Properties)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: reading properties", parameterTypeName),
			err.Error(),
		)
		return
	}

	propertiesJson, err := parameterPropertiesRequest(properties)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: marshalling properties", parameterTypeName),
			err.Error(),
		)
		return
	}

	// UpdateParameter carries the full desired state of the updatable fields: Properties
	// (required) and Owner (required) are always sent, and Description is sent as an empty
	// string when the plan clears it so the server-side value is reset rather than left stale.
	request := map[string]interface{}{
		"RegionId":   r.Meta().RegionId,
		"Id":         state.Id.ValueString(),
		"Properties": propertiesJson,
		"Owner":      plan.Owner.ValueString(),
	}
	if plan.Description.IsNull() || plan.Description.IsUnknown() {
		request["Description"] = ""
	} else {
		request["Description"] = plan.Description.ValueString()
	}

	if _, err := r.Meta().RpcPost("dataworks-public", "2024-05-18", "UpdateParameter", nil, request, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateParameter", parameterTypeName),
			errs.WrapErrorf(err, "Updating %s: calling UpdateParameter", parameterTypeName).Error(),
		)
		return
	}

	object, err := getParameter(r.Meta(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetParameter", parameterTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetParameter", parameterTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	newState, stateDiags := parameterState(ctx, object)
	resp.Diagnostics.Append(stateDiags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *parameterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state parameterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.Meta().RpcPost("dataworks-public", "2024-05-18", "DeleteParameter", nil, map[string]interface{}{
		"RegionId": r.Meta().RegionId,
		"Id":       state.Id.ValueString(),
	}, true)
	if err != nil && !isParameterNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteParameter", parameterTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteParameter", parameterTypeName).Error(),
		)
		return
	}
}

func (r *parameterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getParameter queries GetParameter by parameter ID and returns the Parameter object from the
// response, or nil when the parameter no longer exists.
func getParameter(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	response, err := client.RpcGet("dataworks-public", "2024-05-18", "GetParameter", map[string]interface{}{
		"RegionId": client.RegionId,
		"Id":       id,
	}, nil)
	if err != nil {
		if isParameterNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	parameter, _ := response["Parameter"].(map[string]interface{})
	if parameter == nil {
		return nil, fmt.Errorf("the GetParameter response carried no Parameter object: %v", response)
	}
	return parameter, nil
}

// parameterPropertiesFromSet extracts the configured properties set into a slice for marshalling.
func parameterPropertiesFromSet(ctx context.Context, set types.Set) ([]parameterPropertyModel, error) {
	var properties []parameterPropertyModel
	if diags := set.ElementsAs(ctx, &properties, false); diags.HasError() {
		return nil, fmt.Errorf("reading properties set: %v", diags.Errors())
	}
	return properties, nil
}

// parameterPropertiesRequest marshals the configured properties block list into the JSON string the
// CreateParameter/UpdateParameter Properties parameter expects (style=json).
func parameterPropertiesRequest(properties []parameterPropertyModel) (string, error) {
	items := make([]map[string]interface{}, 0, len(properties))
	for _, p := range properties {
		items = append(items, map[string]interface{}{
			"EnvType": p.EnvType.ValueString(),
			"Value":   p.Value.ValueString(),
		})
	}
	raw, err := json.Marshal(items)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

// parameterState maps the Parameter object of a GetParameter response onto the resource model.
// int64 timestamps and IDs are stringified; absent or empty strings become null. Properties are
// rebuilt as a set so plan convergence does not depend on the server's array ordering.
func parameterState(ctx context.Context, parameter map[string]interface{}) (parameterResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	str := func(key string) types.String {
		if v, ok := parameter[key]; ok {
			switch s := v.(type) {
			case string:
				if s != "" {
					return types.StringValue(s)
				}
			case float64:
				if s != 0 {
					return types.StringValue(fmt.Sprintf("%.0f", s))
				}
			case json.Number:
				if s.String() != "" && s.String() != "0" {
					return types.StringValue(s.String())
				}
			case int:
				if s != 0 {
					return types.StringValue(fmt.Sprint(s))
				}
			case int64:
				if s != 0 {
					return types.StringValue(fmt.Sprint(s))
				}
			}
		}
		return types.StringNull()
	}

	properties := make([]parameterPropertyModel, 0)
	if raw, ok := parameter["Properties"].([]interface{}); ok {
		for _, item := range raw {
			m, _ := item.(map[string]interface{})
			properties = append(properties, parameterPropertyModel{
				EnvType: parameterStringField(m, "EnvType"),
				Value:   parameterStringField(m, "Value"),
			})
		}
	}
	propsSet, setDiags := types.SetValueFrom(ctx, parameterPropertyAttrType, properties)
	diags.Append(setDiags...)

	return parameterResourceModel{
		Id:          str("Id"),
		Name:        str("Name"),
		Type:        str("Type"),
		ProjectId:   str("ProjectId"),
		Scope:       str("Scope"),
		Owner:       str("Owner"),
		Description: str("Description"),
		Properties:  propsSet,
		CreateTime:  str("CreateTime"),
		CreateUser:  str("CreateUser"),
		ModifyTime:  str("ModifyTime"),
		ModifyUser:  str("ModifyUser"),
		Version:     str("Version"),
	}, diags
}

// parameterStringField reads a string field from a nested response object, turning an empty or
// absent string into null.
func parameterStringField(m map[string]interface{}, key string) types.String {
	if m == nil {
		return types.StringNull()
	}
	if v, ok := m[key].(string); ok && v != "" {
		return types.StringValue(v)
	}
	return types.StringNull()
}

// isParameterNotExist reports whether err is the DataWorks parameter-not-found error. errs.NotFoundError
// is the shared classifier used across the dataworks package (HTTP 404 and the canonical NotFound
// message). The exact DataWorks Parameter not-found code is not exposed by the OpenAPI; the
// anticipated Parameter not-exist code forms on the *tea.SDKError are matched as a fallback so
// Read and Delete stay idempotent, mirroring the SDK resources' reliance on the shared classifier.
func isParameterNotExist(err error) bool {
	if err == nil {
		return false
	}
	if errs.NotFoundError(err) {
		return true
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil {
		c := *e.Code
		if strings.Contains(c, "Parameter.NotExist") || strings.Contains(c, "NotExist.Parameter") ||
			strings.Contains(c, "ParameterNotExist") {
			return true
		}
	}
	return strings.Contains(err.Error(), "Parameter.NotExist") ||
		strings.Contains(err.Error(), "NotExist.Parameter") ||
		strings.Contains(err.Error(), "ParameterNotExist")
}

// parameterEmptyStringToNullPlanModifier normalizes an explicitly-empty configuration value to
// null so that it converges with the null representation parameterState stores when the API
// omits or clears the Description. DataWorks represents an absent Description as an empty string;
// parameterState turns that empty string into types.StringNull(). Without this modifier a config
// value of description = "" would plan as StringValue("") while the state holds StringNull(), and
// the two are not equivalent in the Framework, producing a perpetual diff and a redundant
// UpdateParameter on every apply. It only acts on the configured value; values sourced from state
// during Read are unaffected, and it never changes the null semantics of Required fields.
type parameterEmptyStringToNullPlanModifier struct{}

func (m parameterEmptyStringToNullPlanModifier) Description(ctx context.Context) string {
	return "Treats an empty string configuration value as null to converge with the API's empty representation."
}

func (m parameterEmptyStringToNullPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m parameterEmptyStringToNullPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if req.ConfigValue.ValueString() == "" {
		resp.PlanValue = types.StringNull()
	}
}
