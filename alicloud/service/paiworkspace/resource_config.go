package paiworkspace

import (
	"context"
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const configTypeName = "alicloud_pai_workspace_config"

// AIWorkSpace product/version used by every Config operation.
const (
	configProduct = "AIWorkSpace"
	configVersion = "2021-02-04"
	configPathSep = ":"
)

var (
	_ resource.Resource                = (*configResource)(nil)
	_ resource.ResourceWithImportState = (*configResource)(nil)
)

func NewConfigResource() resource.Resource {
	return &configResource{}
}

type configResource struct {
	fwadapt.ResourceBase
}

// configLabelModel mirrors one element of the API Labels array (Key/Value).
type configLabelModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// configResourceModel is the Terraform state shape. The composite id is
// "<workspace_id>:<category_name>:<config_key>", carrying the three identity
// fields GetConfig (path + query) and DeleteConfig (path + query) need.
type configResourceModel struct {
	Id              types.String       `tfsdk:"id"`
	WorkspaceId     types.String       `tfsdk:"workspace_id"`
	CategoryName    types.String       `tfsdk:"category_name"`
	ConfigKey       types.String       `tfsdk:"config_key"`
	ConfigValue     types.String       `tfsdk:"config_value"`
	Labels          []configLabelModel `tfsdk:"labels"`
	GmtCreateTime   types.String       `tfsdk:"gmt_create_time"`
	GmtModifiedTime types.String       `tfsdk:"gmt_modified_time"`
}

func (r *configResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = configTypeName
}

func (r *configResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PAI Workspace Config resource.\n\n" +
			"Workspace configurations are keyed by workspace id, category name and config key. " +
			"There is no independent create API: UpdateConfig upserts the configuration, so this " +
			"resource creates and updates through the same `UpdateConfig` operation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource, in the format `<workspace_id>:<category_name>:<config_key>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workspace_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the PAI Workspace that owns the configuration.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"category_name": schema.StringAttribute{
				Required:    true,
				Description: "The category name of the configuration. Changing this creates a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"config_key": schema.StringAttribute{
				Required:    true,
				Description: "The key of the workspace configuration. Changing this creates a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"config_value": schema.StringAttribute{
				Required:    true,
				Description: "The value of the workspace configuration.",
			},
			"gmt_create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the configuration.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gmt_modified_time": schema.StringAttribute{
				Computed:    true,
				Description: "The last modification time of the configuration. It is bumped on every UpdateConfig, so the planned value never survives an apply.",
				// No UseStateForUnknown: the API refreshes the timestamp on update.
			},
		},
		Blocks: map[string]schema.Block{
			"labels": schema.SetNestedBlock{
				Description: "The labels of the workspace configuration. PAI does not honor a client-requested label order: it returns labels in a server-determined order, so the resource models labels as an unordered set. Adding, removing or changing a label value is detected as a change; reordering alone is not.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Optional:    true,
							Description: "The key of the label.",
						},
						"value": schema.StringAttribute{
							Optional:    true,
							Description: "The value of the label.",
						},
					},
				},
			},
		},
	}
}

func (r *configResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan configResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := configPut(r.Client(), &plan); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling UpdateConfig", configTypeName),
			errs.WrapErrorf(err, "Creating %s: calling UpdateConfig", configTypeName).Error(),
		)
		return
	}

	object, err := getConfig(r.Client(), plan.WorkspaceId.ValueString(), plan.CategoryName.ValueString(), plan.ConfigKey.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading back via GetConfig", configTypeName),
			errs.WrapErrorf(err, "Creating %s: reading back via GetConfig", configTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading back via GetConfig", configTypeName),
			fmt.Sprintf("The config %s is not visible via GetConfig right after creation", configID(&plan)),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, configState(object))...)
}

func (r *configResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state configResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceId, categoryName, configKey, ok := parseConfigID(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s", configTypeName),
			fmt.Sprintf("invalid resource id %q: expected %q", state.Id.ValueString(), "<workspace_id>:<category_name>:<config_key>"),
		)
		return
	}

	object, err := getConfig(r.Client(), workspaceId, categoryName, configKey)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetConfig", configTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetConfig", configTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, configState(object))...)
}

func (r *configResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state configResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The identity fields are all ForceNew; Update only fires when config_value
	// or labels change. UpdateConfig is an upsert, so the full body is sent.
	if !plan.ConfigValue.Equal(state.ConfigValue) || !labelsEqual(plan.Labels, state.Labels) {
		if err := configPut(r.Client(), &plan); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling UpdateConfig", configTypeName),
				errs.WrapErrorf(err, "Updating %s: calling UpdateConfig", configTypeName).Error(),
			)
			return
		}
	}

	object, err := getConfig(r.Client(), plan.WorkspaceId.ValueString(), plan.CategoryName.ValueString(), plan.ConfigKey.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: reading back via GetConfig", configTypeName),
			errs.WrapErrorf(err, "Updating %s: reading back via GetConfig", configTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, configState(object))...)
}

func (r *configResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state configResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceId, categoryName, configKey, ok := parseConfigID(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s", configTypeName),
			fmt.Sprintf("invalid resource id %q: expected %q", state.Id.ValueString(), "<workspace_id>:<category_name>:<config_key>"),
		)
		return
	}

	// DeleteConfig path carries WorkspaceId and ConfigKey; CategoryName is a query param.
	path := fmt.Sprintf("/api/v1/workspaces/%s/configs/%s", workspaceId, configKey)
	query := map[string]*string{
		"CategoryName": ptrString(categoryName),
	}
	if _, err := r.Client().RoaDelete(configProduct, configVersion, path, query, nil, nil, true); err != nil {
		// A config that no longer upstream is a successful delete.
		if !isConfigNotExist(err) {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Deleting %s: calling DeleteConfig", configTypeName),
				errs.WrapErrorf(err, "Deleting %s: calling DeleteConfig", configTypeName).Error(),
			)
			return
		}
	}
}

func (r *configResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// configPut upserts a configuration via UpdateConfig. The path is workspace
// scoped; ConfigKey, ConfigValue, CategoryName and Labels travel in the body.
func configPut(client *connectivity.AliyunClient, plan *configResourceModel) error {
	path := fmt.Sprintf("/api/v1/workspaces/%s/config", plan.WorkspaceId.ValueString())
	body := map[string]interface{}{
		"ConfigKey":    plan.ConfigKey.ValueString(),
		"ConfigValue":  plan.ConfigValue.ValueString(),
		"CategoryName": plan.CategoryName.ValueString(),
		"Labels":       configLabelsToAPI(plan.Labels),
	}
	_, err := client.RoaPut(configProduct, configVersion, path, nil, nil, body, true)
	return err
}

// getConfig reads a configuration via GetConfig. It returns the response object
// directly (the root mapping is the whole object), or nil when the configuration
// no longer exists. NotFound covers both a 404 and the SDK NotFound code set.
func getConfig(client *connectivity.AliyunClient, workspaceId, categoryName, configKey string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/workspaces/%s/config", workspaceId)
	query := map[string]*string{
		"CategoryName": ptrString(categoryName),
		"ConfigKey":    ptrString(configKey),
		"Verbose":      ptrString("true"),
	}
	response, err := client.RoaGet(configProduct, configVersion, path, query, nil, nil)
	if err != nil {
		if isConfigNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	if response == nil {
		return nil, nil
	}
	// A configuration that exists always carries its identity key. GetConfig
	// returns the object directly rather than a list, so an absent config that
	// the gateway still answers 200 to surfaces as an empty ConfigKey.
	if v, _ := response["ConfigKey"].(string); v == "" {
		return nil, nil
	}
	return response, nil
}

// configState maps a GetConfig response onto the resource model.
func configState(object map[string]interface{}) configResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	workspaceId := str("WorkspaceId")
	categoryName := str("CategoryName")
	configKey := str("ConfigKey")
	return configResourceModel{
		Id:              types.StringValue(configIDFromParts(workspaceId.ValueString(), categoryName.ValueString(), configKey.ValueString())),
		WorkspaceId:     workspaceId,
		CategoryName:    categoryName,
		ConfigKey:       configKey,
		ConfigValue:     str("ConfigValue"),
		Labels:          configLabelsFromAPI(object["Labels"]),
		GmtCreateTime:   str("GmtCreateTime"),
		GmtModifiedTime: str("GmtModifiedTime"),
	}
}

// configLabelsToAPI converts the Terraform labels list to the API body shape.
func configLabelsToAPI(labels []configLabelModel) []interface{} {
	out := make([]interface{}, 0, len(labels))
	for _, l := range labels {
		out = append(out, map[string]interface{}{
			"Key":   l.Key.ValueString(),
			"Value": l.Value.ValueString(),
		})
	}
	return out
}

// configSystemLabelPrefix marks labels PAI injects into every GetConfig
// response (e.g. "system.categoryName", whose value mirrors the config's
// CategoryName). They are not user-managed, so they are filtered out of state —
// otherwise Create/Update/Read would always surface one more label than the
// caller configured, breaking plan/apply consistency and ImportStateVerify.
const configSystemLabelPrefix = "system."

// configLabelsFromAPI converts the API Labels array to the Terraform list
// shape, dropping PAI-injected system labels (Key prefixed "system.").
func configLabelsFromAPI(v interface{}) []configLabelModel {
	if v == nil {
		return nil
	}
	arr, ok := v.([]interface{})
	if !ok || len(arr) == 0 {
		return nil
	}
	out := make([]configLabelModel, 0, len(arr))
	for _, item := range arr {
		m, _ := item.(map[string]interface{})
		key := stringFromAPI(m["Key"])
		if strings.HasPrefix(key.ValueString(), configSystemLabelPrefix) {
			continue
		}
		out = append(out, configLabelModel{
			Key:   key,
			Value: stringFromAPI(m["Value"]),
		})
	}
	return out
}

// stringFromAPI turns a JSON-decoded value into a String, null when empty.
func stringFromAPI(v interface{}) types.String {
	if s, ok := v.(string); ok && s != "" {
		return types.StringValue(s)
	}
	return types.StringNull()
}

// labelsEqual reports whether two label lists hold the same members, ignoring
// order. Labels are an unordered set (PAI returns them in a server-determined
// order), so a position-sensitive compare would flag a harmless reorder as a
// change and trigger a needless UpdateConfig.
func labelsEqual(a, b []configLabelModel) bool {
	if len(a) != len(b) {
		return false
	}
	for _, la := range a {
		found := false
		for _, lb := range b {
			if la.Key.Equal(lb.Key) && la.Value.Equal(lb.Value) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// configID builds the composite id from a plan's identity fields.
func configID(plan *configResourceModel) string {
	return configIDFromParts(plan.WorkspaceId.ValueString(), plan.CategoryName.ValueString(), plan.ConfigKey.ValueString())
}

// configIDFromParts joins the three identity fields with the separator.
func configIDFromParts(workspaceId, categoryName, configKey string) string {
	return strings.Join([]string{workspaceId, categoryName, configKey}, configPathSep)
}

// parseConfigID splits the composite id back into its three parts.
func parseConfigID(id string) (workspaceId, categoryName, configKey string, ok bool) {
	parts := strings.Split(id, configPathSep)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}

// ptrString returns a pointer to s for the ROA query map.
func ptrString(s string) *string {
	return &s
}

// configNotExistCodes are AIWorkSpace error codes that signal the workspace —
// and therefore its config — no longer exists. They are taken from the existing
// DescribePaiWorkspaceWorkspace not-found set; a config-specific not-found code
// has not yet been confirmed from runtime evidence, so the substring fallback
// in isConfigNotExist and the empty-ConfigKey defense in getConfig carry that
// case until an ACC run surfaces the precise code.
var configNotExistCodes = map[string]bool{
	"100400008": true,
	"100700008": true,
	"100400027": true,
}

// isConfigNotExist reports whether err means the config is gone. It combines
// errs.NotFoundError (HTTP 404 + "NotFound" in message), the AIWorkSpace
// workspace-not-found SDK codes, and a "notfound"/"notexist" substring fallback
// that mirrors the IMS isUserNotExist pattern for wrapped error forms.
func isConfigNotExist(err error) bool {
	if err == nil {
		return false
	}
	if errs.NotFoundError(err) {
		return true
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil && configNotExistCodes[*e.Code] {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "notfound") || strings.Contains(msg, "notexist")
}
