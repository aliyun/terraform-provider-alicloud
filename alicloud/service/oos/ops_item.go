package oos

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
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

const opsItemTypeName = "alicloud_oos_ops_item"

// opsItemSolutionObjectAttrTypes is the Framework object type of a single
// solution. ListNestedBlock infers the same shape from opsItemSolutionModel's
// tfsdk tags; this explicit type is needed to build the list from API responses.
var opsItemSolutionObjectAttrTypes = map[string]attr.Type{
	"priority":    types.Int64Type,
	"type":        types.StringType,
	"url":         types.StringType,
	"description": types.StringType,
}

var (
	_ resource.Resource                = &opsItemResource{}
	_ resource.ResourceWithImportState = &opsItemResource{}
)

func NewOpsItemResource() resource.Resource {
	return &opsItemResource{}
}

type opsItemResource struct {
	fwadapt.ResourceBase
}

type opsItemSolutionModel struct {
	Priority    types.Int64  `tfsdk:"priority"`
	Type        types.String `tfsdk:"type"`
	Url         types.String `tfsdk:"url"`
	Description types.String `tfsdk:"description"`
}

type opsItemResourceModel struct {
	Id              types.String `tfsdk:"id"`
	OpsItemId       types.String `tfsdk:"ops_item_id"`
	Title           types.String `tfsdk:"title"`
	Severity        types.String `tfsdk:"severity"`
	Category        types.String `tfsdk:"category"`
	Source          types.String `tfsdk:"source"`
	Priority        types.Int64  `tfsdk:"priority"`
	Description     types.String `tfsdk:"description"`
	Status          types.String `tfsdk:"status"`
	DedupString     types.String `tfsdk:"dedup_string"`
	ResourceGroupId types.String `tfsdk:"resource_group_id"`
	Tags            types.Map    `tfsdk:"tags"`
	Solutions       types.List   `tfsdk:"solutions"`
	Resources       types.List   `tfsdk:"resources"`
	CreateDate      types.String `tfsdk:"create_date"`
	UpdateDate      types.String `tfsdk:"update_date"`
	CreateBy        types.String `tfsdk:"create_by"`
	LastModifiedBy  types.String `tfsdk:"last_modified_by"`
	Attributes      types.Map    `tfsdk:"attributes"`
}

func (r *opsItemResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = opsItemTypeName
}

func (r *opsItemResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an OOS Ops Item resource. An Ops Item tracks an operational event across resources. " +
			"OOS exposes no hard-delete API for Ops Items, so `terraform destroy` transitions the Ops Item to the " +
			"`Resolved` status rather than removing it from OOS.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. Same as `ops_item_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ops_item_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the Ops Item.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Required:    true,
				Description: "The title of the Ops Item.",
			},
			"severity": schema.StringAttribute{
				Required:    true,
				Description: "The severity of the Ops Item. Valid values: `Critical`, `High`, `Medium`, `Low`.",
				Validators: []validator.String{
					stringvalidator.OneOf("Critical", "High", "Medium", "Low"),
				},
			},
			"category": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The category of the Ops Item. Valid values: `Availability`, `Cost`, `Performance`, `Recovery`, `Security`.",
				Validators: []validator.String{
					stringvalidator.OneOf("Availability", "Cost", "Performance", "Recovery", "Security"),
				},
			},
			"source": schema.StringAttribute{
				Required:    true,
				Description: "The source of the Ops Item. Required by the CreateOpsItem API, which returns 400 MissingSource when it is omitted.",
			},
			"priority": schema.Int64Attribute{
				Required:    true,
				Description: "The priority of the Ops Item. Required by the CreateOpsItem API, which returns 400 InvalidParameter when it is omitted.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the Ops Item.",
			},
			"status": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "The status of the Ops Item. OOS creates Ops Items as `Open`; the status is applied on update, not at creation. " +
					"Valid values: `Open`, `InProgress`, `Resolved`. Destroying the resource sets `Resolved`.",
				Validators: []validator.String{
					stringvalidator.OneOf("Open", "InProgress", "Resolved"),
				},
			},
			"dedup_string": schema.StringAttribute{
				Optional: true,
				Description: "The deduplication string of the Ops Item. Create-only: changing it forces replacement. " +
					"OOS does not return this field on read, so the configured value is preserved in state.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"resource_group_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The ID of the resource group.",
			},
			"tags": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "The tags of the Ops Item.",
			},
			"resources": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "The list of associated resource ARNs.",
			},
			"create_date": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the Ops Item.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"update_date": schema.StringAttribute{
				Computed:    true,
				Description: "The last update time of the Ops Item.",
				// No UseStateForUnknown: UpdateOpsItem bumps the timestamp.
			},
			"create_by": schema.StringAttribute{
				Computed:    true,
				Description: "The creator of the Ops Item.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_modified_by": schema.StringAttribute{
				Computed:    true,
				Description: "The last modifier of the Ops Item.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"attributes": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "The custom attributes of the Ops Item.",
			},
		},
		Blocks: map[string]schema.Block{
			"solutions": schema.ListNestedBlock{
				Description: "The solutions of the Ops Item.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"priority": schema.Int64Attribute{
							Optional:    true,
							Computed:    true,
							Description: "The priority of the solution.",
						},
						"type": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The type of the solution.",
						},
						"url": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The URL of the solution.",
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The description of the solution.",
						},
					},
				},
			},
		},
	}
}

func (r *opsItemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan opsItemResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"Title":    plan.Title.ValueString(),
		"Severity": plan.Severity.ValueString(),
		"RegionId": r.Client().RegionId,
	}
	if v, ok := stringVal(plan.Category); ok {
		request["Category"] = v
	}
	if v, ok := stringVal(plan.Source); ok {
		request["Source"] = v
	}
	if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
		request["Priority"] = plan.Priority.ValueInt64()
	}
	if v, ok := stringVal(plan.Description); ok {
		request["Description"] = v
	}
	if v, ok := stringVal(plan.DedupString); ok {
		request["DedupString"] = v
	}
	if v, ok := stringVal(plan.ResourceGroupId); ok {
		request["ResourceGroupId"] = v
	}
	if s, err := tagsToJSONString(plan.Tags); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Creating %s", opsItemTypeName), err.Error())
		return
	} else if s != "" {
		request["Tags"] = s
	}
	if s, err := stringListToJSONString(plan.Resources); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Creating %s", opsItemTypeName), err.Error())
		return
	} else if s != "" {
		request["Resources"] = s
	}
	if s, err := solutionsToJSONString(ctx, plan.Solutions); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Creating %s", opsItemTypeName), err.Error())
		return
	} else if s != "" {
		request["Solutions"] = s
	}

	response, err := r.Client().RpcPost("oos", "2019-06-01", "CreateOpsItem", nil, request, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateOpsItem", opsItemTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateOpsItem", opsItemTypeName).Error(),
		)
		return
	}

	opsItem, _ := response["OpsItem"].(map[string]interface{})
	opsItemId, _ := opsItem["OpsItemId"].(string)
	if opsItemId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateOpsItem", opsItemTypeName),
			fmt.Sprintf("The response carried no OpsItem.OpsItemId: %v", response),
		)
		return
	}

	// CreateOpsItem has no Status parameter, so OOS creates the Ops Item as
	// "Open". If the plan sets a non-default status, apply it via an immediate
	// UpdateOpsItem so the post-create state matches the plan and the framework
	// does not flag an inconsistent result after apply.
	if desired, ok := stringVal(plan.Status); ok && desired != "" && desired != "Open" {
		if _, err := r.Client().RpcPost("oos", "2019-06-01", "UpdateOpsItem", nil, map[string]interface{}{
			"OpsItemId": opsItemId,
			"Status":    desired,
			"RegionId":  r.Client().RegionId,
		}, true); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: applying initial status", opsItemTypeName),
				errs.WrapErrorf(err, "Creating %s: applying initial status", opsItemTypeName).Error(),
			)
			return
		}
	}

	object, err := getOpsItem(r.Client(), opsItemId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetOpsItem", opsItemTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetOpsItem", opsItemTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling GetOpsItem", opsItemTypeName),
			fmt.Sprintf("The Ops Item %s is not visible via GetOpsItem right after creation", opsItemId),
		)
		return
	}
	state, diags := opsItemState(ctx, object, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *opsItemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state opsItemResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Import populates only `id` (ImportStatePassthroughID); reconcile the
	// natural key so the first GetOpsItem after import is not called with "".
	if state.OpsItemId.IsNull() && !state.Id.IsNull() {
		state.OpsItemId = state.Id
	}

	object, err := getOpsItem(r.Client(), state.OpsItemId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetOpsItem", opsItemTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetOpsItem", opsItemTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	newState, diags := opsItemState(ctx, object, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *opsItemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state opsItemResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"OpsItemId": state.OpsItemId.ValueString(),
		"RegionId":  r.Client().RegionId,
	}
	// DedupString is create-only (forces replacement), so it is never sent here.
	if v, ok := stringVal(plan.Title); ok {
		request["Title"] = v
	}
	if v, ok := stringVal(plan.Severity); ok {
		request["Severity"] = v
	}
	if v, ok := stringVal(plan.Category); ok {
		request["Category"] = v
	}
	if v, ok := stringVal(plan.Source); ok {
		request["Source"] = v
	}
	if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
		request["Priority"] = plan.Priority.ValueInt64()
	}
	if v, ok := stringVal(plan.Description); ok {
		request["Description"] = v
	}
	if v, ok := stringVal(plan.Status); ok {
		request["Status"] = v
	}
	if v, ok := stringVal(plan.ResourceGroupId); ok {
		request["ResourceGroupId"] = v
	}
	// Tags/Resources/Solutions are Optional+Computed. A null plan value means
	// "leave as-is" (Computed carries the state value), so it is omitted; an
	// explicitly-empty collection serializes to {} or [] so the API clears it.
	if s, err := tagsToJSONString(plan.Tags); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Updating %s", opsItemTypeName), err.Error())
		return
	} else if s != "" {
		request["Tags"] = s
	}
	if s, err := stringListToJSONString(plan.Resources); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Updating %s", opsItemTypeName), err.Error())
		return
	} else if s != "" {
		request["Resources"] = s
	}
	if s, err := solutionsToJSONString(ctx, plan.Solutions); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Updating %s", opsItemTypeName), err.Error())
		return
	} else if s != "" {
		request["Solutions"] = s
	}

	if _, err := r.Client().RpcPost("oos", "2019-06-01", "UpdateOpsItem", nil, request, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateOpsItem", opsItemTypeName),
			errs.WrapErrorf(err, "Updating %s: calling UpdateOpsItem", opsItemTypeName).Error(),
		)
		return
	}

	object, err := getOpsItem(r.Client(), state.OpsItemId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetOpsItem", opsItemTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetOpsItem", opsItemTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	newState, diags := opsItemState(ctx, object, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *opsItemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state opsItemResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// OOS exposes no DeleteOpsItem API. Destroy transitions the Ops Item to the
	// terminal "Resolved" status so the operational record is closed; the Ops
	// Item remains in OOS and Terraform drops it from state.
	_, err := r.Client().RpcPost("oos", "2019-06-01", "UpdateOpsItem", nil, map[string]interface{}{
		"OpsItemId": state.OpsItemId.ValueString(),
		"Status":    "Resolved",
		"RegionId":  r.Client().RegionId,
	}, true)
	if err != nil && !isOpsItemNotFound(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling UpdateOpsItem", opsItemTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling UpdateOpsItem", opsItemTypeName).Error(),
		)
		return
	}
}

func (r *opsItemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getOpsItem queries GetOpsItem by OpsItemId and returns the OpsItem object from
// the response, or nil when the Ops Item no longer exists.
func getOpsItem(client *connectivity.AliyunClient, opsItemId string) (map[string]interface{}, error) {
	response, err := client.RpcPost("oos", "2019-06-01", "GetOpsItem", nil, map[string]interface{}{
		"OpsItemId": opsItemId,
		"RegionId":  client.RegionId,
	}, true)
	if err != nil {
		if isOpsItemNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	opsItem, _ := response["OpsItem"].(map[string]interface{})
	if opsItem == nil {
		return nil, fmt.Errorf("the GetOpsItem response carried no OpsItem object: %v", response)
	}
	return opsItem, nil
}

// opsItemState maps the OpsItem object of a GetOpsItem response onto the
// resource model. DedupString is create-only and not returned by GetOpsItem,
// so the value carried by prior (the plan on create, the state on read/update)
// is preserved rather than being lost across reads.
func opsItemState(ctx context.Context, object map[string]interface{}, prior *opsItemResourceModel) (opsItemResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	m := opsItemResourceModel{
		Id:              str("OpsItemId"),
		OpsItemId:       str("OpsItemId"),
		Title:           str("Title"),
		Severity:        str("Severity"),
		Category:        str("Category"),
		Source:          str("Source"),
		Description:     str("Description"),
		Status:          str("Status"),
		ResourceGroupId: str("ResourceGroupId"),
		CreateDate:      str("CreateDate"),
		UpdateDate:      str("UpdateDate"),
		CreateBy:        str("CreateBy"),
		LastModifiedBy:  str("LastModifiedBy"),
	}
	if v, ok := toInt64(object["Priority"]); ok {
		m.Priority = types.Int64Value(v)
	} else {
		m.Priority = types.Int64Null()
	}
	// DedupString is not part of the GetOpsItem response; keep the configured
	// value from the plan/state so it survives refresh.
	m.DedupString = str("DedupString")
	if m.DedupString.IsNull() && prior != nil {
		m.DedupString = prior.DedupString
	}
	if v, ok := object["Tags"].(map[string]interface{}); ok && len(v) > 0 {
		tags := make(map[string]string, len(v))
		for k, val := range v {
			if s, ok := val.(string); ok {
				tags[k] = s
			}
		}
		if len(tags) > 0 {
			tagsMap, tdiags := types.MapValueFrom(ctx, types.StringType, tags)
			diags = append(diags, tdiags...)
			m.Tags = tagsMap
		} else {
			m.Tags = types.MapNull(types.StringType)
		}
	} else {
		m.Tags = types.MapNull(types.StringType)
	}
	if arr, ok := object["Resources"].([]interface{}); ok && len(arr) > 0 {
		res := make([]string, 0, len(arr))
		for _, raw := range arr {
			if s, ok := raw.(string); ok && s != "" {
				res = append(res, s)
			}
		}
		if len(res) > 0 {
			list, ldiags := types.ListValueFrom(ctx, types.StringType, res)
			diags = append(diags, ldiags...)
			m.Resources = list
		} else {
			m.Resources = types.ListNull(types.StringType)
		}
	} else {
		m.Resources = types.ListNull(types.StringType)
	}
	if arr, ok := object["Solutions"].([]interface{}); ok && len(arr) > 0 {
		sols := make([]opsItemSolutionModel, 0, len(arr))
		for _, raw := range arr {
			s, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}
			sol := opsItemSolutionModel{
				Type:        nullableString(s["Type"]),
				Url:         nullableString(s["Url"]),
				Description: nullableString(s["Description"]),
			}
			if v, ok := toInt64(s["Priority"]); ok {
				sol.Priority = types.Int64Value(v)
			} else {
				sol.Priority = types.Int64Null()
			}
			sols = append(sols, sol)
		}
		if len(sols) > 0 {
			list, ldiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: opsItemSolutionObjectAttrTypes}, sols)
			diags = append(diags, ldiags...)
			m.Solutions = list
		} else {
			m.Solutions = types.ListNull(types.ObjectType{AttrTypes: opsItemSolutionObjectAttrTypes})
		}
	} else {
		m.Solutions = types.ListNull(types.ObjectType{AttrTypes: opsItemSolutionObjectAttrTypes})
	}
	if v, ok := object["Attributes"].(map[string]interface{}); ok && len(v) > 0 {
		attrs := make(map[string]string, len(v))
		for k, val := range v {
			if s, ok := val.(string); ok {
				attrs[k] = s
			}
		}
		if len(attrs) > 0 {
			am, adiags := types.MapValueFrom(ctx, types.StringType, attrs)
			diags = append(diags, adiags...)
			m.Attributes = am
		} else {
			m.Attributes = types.MapNull(types.StringType)
		}
	} else {
		m.Attributes = types.MapNull(types.StringType)
	}
	return m, diags
}

// isOpsItemNotFound reports whether err is the OOS Ops-Item-not-found error.
// The code arrives on a *tea.SDKError from RpcPost; the substring fallback
// covers wrapped forms, matching IsExpectedErrors' behavior on the SDK line.
func isOpsItemNotFound(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil && (*e.Code == "OpsItem.NotFound" || *e.Code == "InvalidOpsItem.NotFound") {
		return true
	}
	msg := err.Error()
	return strings.Contains(msg, "OpsItem.NotFound") || strings.Contains(msg, "InvalidOpsItem.NotFound")
}

// stringVal returns the underlying string for a non-null, non-unknown
// types.String and reports whether a value was present.
func stringVal(v types.String) (string, bool) {
	if v.IsNull() || v.IsUnknown() {
		return "", false
	}
	return v.ValueString(), true
}

// nullableString builds a types.String from a response value, mapping missing
// or non-string values to null.
func nullableString(v interface{}) types.String {
	if s, ok := v.(string); ok && s != "" {
		return types.StringValue(s)
	}
	return types.StringNull()
}

// toInt64 coerces numeric response values (json.Number or float64) to int64.
func toInt64(v interface{}) (int64, bool) {
	switch n := v.(type) {
	case json.Number:
		i, err := n.Int64()
		return i, err == nil
	case float64:
		return int64(n), true
	case int:
		return int64(n), true
	case int64:
		return n, true
	}
	return 0, false
}

// tagsToJSONString marshals a tags map to the JSON object string the OOS API
// expects (CreateOpsItem/UpdateOpsItem Tags is type object, style json). A
// null/unknown map yields "" so the caller can omit the parameter; an
// explicitly-empty map yields "{}" so the API clears existing tags.
func tagsToJSONString(m types.Map) (string, error) {
	if m.IsNull() || m.IsUnknown() {
		return "", nil
	}
	tags := make(map[string]string, len(m.Elements()))
	for k, v := range m.Elements() {
		if s, ok := v.(types.String); ok && !s.IsNull() {
			tags[k] = s.ValueString()
		}
	}
	b, err := json.Marshal(tags)
	if err != nil {
		return "", fmt.Errorf("marshaling tags: %w", err)
	}
	return string(b), nil
}

// stringListToJSONString marshals a list of strings to the JSON array string
// the OOS API expects (Resources is type string, a JSON-encoded array). A
// null/unknown list yields "" (omit); an explicitly-empty list yields "[]"
// (clear).
func stringListToJSONString(l types.List) (string, error) {
	if l.IsNull() || l.IsUnknown() {
		return "", nil
	}
	res := make([]string, 0, len(l.Elements()))
	for _, v := range l.Elements() {
		if s, ok := v.(types.String); ok && !s.IsNull() {
			res = append(res, s.ValueString())
		}
	}
	b, err := json.Marshal(res)
	if err != nil {
		return "", fmt.Errorf("marshaling resources: %w", err)
	}
	return string(b), nil
}

// solutionsToJSONString marshals a list of solution objects to the JSON array
// string the OOS API expects (Solutions is type string, a JSON-encoded array).
// A null/unknown list yields "" (omit); an explicitly-empty list yields "[]"
// (clear).
func solutionsToJSONString(ctx context.Context, l types.List) (string, error) {
	if l.IsNull() || l.IsUnknown() {
		return "", nil
	}
	var sols []opsItemSolutionModel
	if diags := l.ElementsAs(ctx, &sols, false); diags.HasError() {
		return "", fmt.Errorf("reading solutions from plan: %v", diags)
	}
	out := make([]map[string]interface{}, 0, len(sols))
	for _, s := range sols {
		item := map[string]interface{}{}
		if !s.Priority.IsNull() && !s.Priority.IsUnknown() {
			item["priority"] = s.Priority.ValueInt64()
		}
		if v, ok := stringVal(s.Type); ok {
			item["type"] = v
		}
		if v, ok := stringVal(s.Url); ok {
			item["url"] = v
		}
		if v, ok := stringVal(s.Description); ok {
			item["description"] = v
		}
		out = append(out, item)
	}
	b, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("marshaling solutions: %w", err)
	}
	return string(b), nil
}
