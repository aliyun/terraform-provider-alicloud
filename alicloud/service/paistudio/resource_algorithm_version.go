package paistudio

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	algorithmVersionTypeName = "alicloud_pai_algorithm_version"

	// PaiStudio is an independent POP, distinct from AIWorkSpace (2021-02-04)
	// and PAIFlow. Its RESTful/ROA API version is 2022-01-12. Exported so the
	// external _test package can build a standalone client for fixtures.
	PaiStudioProduct = "PaiStudio"
	PaiStudioVersion = "2022-01-12"

	// An AlgorithmVersion is addressed by the algorithm it belongs to and the
	// version name; both form the path. The composite Terraform ID keeps the
	// pair together for import.
	idSeparator = ":"
)

var (
	_ resource.Resource                = &algorithmVersionResource{}
	_ resource.ResourceWithImportState = &algorithmVersionResource{}
)

func NewAlgorithmVersionResource() resource.Resource {
	return &algorithmVersionResource{}
}

type algorithmVersionResource struct {
	fwadapt.ResourceBase
}

type algorithmVersionResourceModel struct {
	Id                types.String `tfsdk:"id"`
	AlgorithmId       types.String `tfsdk:"algorithm_id"`
	AlgorithmVersion  types.String `tfsdk:"algorithm_version"`
	AlgorithmSpec     types.String `tfsdk:"algorithm_spec"`
	TenantId          types.String `tfsdk:"tenant_id"`
	UserId            types.String `tfsdk:"user_id"`
	AlgorithmProvider types.String `tfsdk:"algorithm_provider"`
	AlgorithmName     types.String `tfsdk:"algorithm_name"`
	GmtModifiedTime   types.String `tfsdk:"gmt_modified_time"`
	RegionId          types.String `tfsdk:"region_id"`
	CreateTime        types.String `tfsdk:"create_time"`
}

func (r *algorithmVersionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = algorithmVersionTypeName
}

func (r *algorithmVersionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PaiStudio Algorithm Version resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. It is formatted as `<algorithm_id>:<algorithm_version>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"algorithm_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the algorithm that the version belongs to. It is the parent algorithm registered in PaiStudio.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"algorithm_version": schema.StringAttribute{
				Required:    true,
				Description: "The version name of the algorithm version, e.g. `v1.0.0`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"algorithm_spec": schema.StringAttribute{
				Required: true,
				Description: "The algorithm spec, a JSON-encoded object matching the PaiStudio `AlgorithmSpec` schema " +
					("(JobType, Image, Command, InputChannels, OutputChannels, HyperParameters, etc.). " +
						"Use `jsonencode(...)` to build it from a Terraform object."),
			},
			"tenant_id": schema.StringAttribute{
				Computed:    true,
				Description: "The tenant ID of the algorithm.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.StringAttribute{
				Computed:    true,
				Description: "The user ID of the algorithm.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"algorithm_provider": schema.StringAttribute{
				Computed:    true,
				Description: "The algorithm provider.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"algorithm_name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the algorithm.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gmt_modified_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the algorithm version was last modified.",
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
				Description: "The time when the algorithm version was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *algorithmVersionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan algorithmVersionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spec, err := parseAlgorithmSpec(plan.AlgorithmSpec.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: parsing algorithm_spec", algorithmVersionTypeName),
			err.Error(),
		)
		return
	}

	path := algorithmVersionPath(plan.AlgorithmId.ValueString(), plan.AlgorithmVersion.ValueString())
	body := map[string]interface{}{"AlgorithmSpec": spec}
	if _, err := r.Client().RoaPost(PaiStudioProduct, PaiStudioVersion, path, nil, nil, body, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateAlgorithmVersion", algorithmVersionTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateAlgorithmVersion", algorithmVersionTypeName).Error(),
		)
		return
	}

	// Persist the created resource BEFORE reading. The composite id is
	// knowable from the plan identifiers (algorithm_id:algorithm_version) with
	// no Get response, so a later read failure or transient invisibility
	// cannot drop an already-created resource: the state keeps the id and the
	// original plan. The first Set is the recovery anchor; the final Set
	// refreshes computed fields from the read response.
	plan.Id = types.StringValue(algorithmVersionCompositeID(plan.AlgorithmId.ValueString(), plan.AlgorithmVersion.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := getAlgorithmVersion(r.Client(), plan.AlgorithmId.ValueString(), plan.AlgorithmVersion.ValueString())
	if err != nil {
		// Read failed after a successful create: keep the persisted plan state
		// (id + original config) and surface the error so the resource stays
		// tracked for the next refresh instead of being lost.
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetAlgorithmVersion", algorithmVersionTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetAlgorithmVersion", algorithmVersionTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling GetAlgorithmVersion", algorithmVersionTypeName),
			fmt.Sprintf("The algorithm version %s/%s is not visible via GetAlgorithmVersion right after creation; the resource is persisted in state for the next refresh",
				plan.AlgorithmId.ValueString(), plan.AlgorithmVersion.ValueString()),
		)
		return
	}
	if err := validateAlgorithmVersionResponse(object, plan.AlgorithmId.ValueString(), plan.AlgorithmVersion.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: validating GetAlgorithmVersion response", algorithmVersionTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, algorithmVersionState(object))...)
}

func (r *algorithmVersionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state algorithmVersionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := getAlgorithmVersion(r.Client(), state.AlgorithmId.ValueString(), state.AlgorithmVersion.ValueString())
	if err != nil {
		// A transient read error must not drop the existing state; return the
		// error without Set/RemoveResource so the prior state is preserved.
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetAlgorithmVersion", algorithmVersionTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetAlgorithmVersion", algorithmVersionTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	// Validate the response identity/type before flattening onto state. A
	// missing or wrong-typed AlgorithmId/AlgorithmVersion would otherwise
	// compose a ":" id, and a mismatched response could overwrite the
	// resource's identity. On such a malformed response, keep the prior state
	// and surface the error rather than writing a corrupt object.
	if err := validateAlgorithmVersionResponse(object, state.AlgorithmId.ValueString(), state.AlgorithmVersion.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: validating GetAlgorithmVersion response", algorithmVersionTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, algorithmVersionState(object))...)
}

func (r *algorithmVersionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state algorithmVersionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spec, err := parseAlgorithmSpec(plan.AlgorithmSpec.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: parsing algorithm_spec", algorithmVersionTypeName),
			err.Error(),
		)
		return
	}

	path := algorithmVersionPath(state.AlgorithmId.ValueString(), state.AlgorithmVersion.ValueString())
	body := map[string]interface{}{"AlgorithmSpec": spec}
	if _, err := r.Client().RoaPut(PaiStudioProduct, PaiStudioVersion, path, nil, nil, body, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateAlgorithmVersion", algorithmVersionTypeName),
			errs.WrapErrorf(err, "Updating %s: calling UpdateAlgorithmVersion", algorithmVersionTypeName).Error(),
		)
		return
	}

	object, err := getAlgorithmVersion(r.Client(), state.AlgorithmId.ValueString(), state.AlgorithmVersion.ValueString())
	if err != nil {
		// Keep the prior state on a read failure after update.
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetAlgorithmVersion", algorithmVersionTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetAlgorithmVersion", algorithmVersionTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	// Validate the response identity/type before flattening; a malformed
	// response keeps the prior state and surfaces an error.
	if err := validateAlgorithmVersionResponse(object, state.AlgorithmId.ValueString(), state.AlgorithmVersion.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: validating GetAlgorithmVersion response", algorithmVersionTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, algorithmVersionState(object))...)
}

func (r *algorithmVersionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state algorithmVersionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	path := algorithmVersionPath(state.AlgorithmId.ValueString(), state.AlgorithmVersion.ValueString())
	_, err := r.Client().RoaDelete(PaiStudioProduct, PaiStudioVersion, path, nil, nil, nil, true)
	if err != nil && !isAlgorithmVersionNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteAlgorithmVersion", algorithmVersionTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteAlgorithmVersion", algorithmVersionTypeName).Error(),
		)
		return
	}
}

func (r *algorithmVersionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, idSeparator, 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Importing %s", algorithmVersionTypeName),
			fmt.Sprintf("The import ID must be formatted as `<algorithm_id>%s<algorithm_version>`; got %q", idSeparator, req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, algorithmVersionResourceModel{
		Id:               types.StringValue(req.ID),
		AlgorithmId:      types.StringValue(parts[0]),
		AlgorithmVersion: types.StringValue(parts[1]),
	})...)
}

// algorithmVersionPath builds the ROA path that addresses a single algorithm
// version. Both identifiers live in the path, not the query, so they are
// substituted here rather than passed in the query map.
func algorithmVersionPath(algorithmId, algorithmVersion string) string {
	return fmt.Sprintf("/api/v1/algorithms/%s/versions/%s", algorithmId, algorithmVersion)
}

// getAlgorithmVersion queries GetAlgorithmVersion and returns the version
// object, or nil when it no longer exists (NoSuchObject).
func getAlgorithmVersion(client *connectivity.AliyunClient, algorithmId, algorithmVersion string) (map[string]interface{}, error) {
	response, err := client.RoaGet(PaiStudioProduct, PaiStudioVersion,
		algorithmVersionPath(algorithmId, algorithmVersion), nil, nil, nil)
	if err != nil {
		if isAlgorithmVersionNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return response, nil
}

// algorithmVersionState maps a GetAlgorithmVersion response onto the resource
// model. Empty strings become null: dates and optional fields the API omits
// read back as absent rather than "". The AlgorithmSpec object is stored as a
// canonical JSON string so the plan stays stable across refreshes.
func algorithmVersionState(object map[string]interface{}) algorithmVersionResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	spec := types.StringNull()
	if raw, ok := object["AlgorithmSpec"]; ok && raw != nil {
		if b, err := json.Marshal(raw); err == nil {
			spec = types.StringValue(string(b))
		}
	}
	algorithmId := rawString(object, "AlgorithmId")
	algorithmVersion := rawString(object, "AlgorithmVersion")
	return algorithmVersionResourceModel{
		Id:                types.StringValue(algorithmVersionCompositeID(algorithmId, algorithmVersion)),
		AlgorithmId:       str("AlgorithmId"),
		AlgorithmVersion:  str("AlgorithmVersion"),
		AlgorithmSpec:     spec,
		TenantId:          str("TenantId"),
		UserId:            str("UserId"),
		AlgorithmProvider: str("AlgorithmProvider"),
		AlgorithmName:     str("AlgorithmName"),
		GmtModifiedTime:   str("GmtModifiedTime"),
		RegionId:          str("RegionId"),
		CreateTime:        str("CreateTime"),
	}
}

// rawString returns object[key] as a string, or "" when absent.
func rawString(object map[string]interface{}, key string) string {
	if v, ok := object[key].(string); ok {
		return v
	}
	return ""
}

// parseAlgorithmSpec parses and canonicalizes the algorithm_spec JSON string.
// Canonicalizing (re-marshaling) keeps the stored value independent of the
// user's key order and whitespace, so the plan does not churn on refresh.
func parseAlgorithmSpec(raw string) (interface{}, error) {
	var spec interface{}
	if err := json.Unmarshal([]byte(raw), &spec); err != nil {
		return nil, fmt.Errorf("algorithm_spec is not valid JSON: %w", err)
	}
	return spec, nil
}

// isAlgorithmVersionNotExist reports whether err is the PaiStudio not-found
// error. The ROA call surfaces the product error contract on a *tea.SDKError
// whose Code is NoSuchObject; that is the only signal used here. A loose
// substring match is intentionally avoided so a non-NotFound message that
// merely contains the word "NoSuchObject" cannot be misread as not-found and
// silently drop the resource from state.
func isAlgorithmVersionNotExist(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil && *e.Code == "NoSuchObject" {
		return true
	}
	return false
}

// algorithmVersionCompositeID builds the Terraform composite id that keeps the
// algorithm/version pair together for import. It is used both when persisting
// state immediately after a successful create (before any Get response) and
// when flattening a read response, so the two cannot diverge.
func algorithmVersionCompositeID(algorithmId, algorithmVersion string) string {
	return fmt.Sprintf("%s%s%s", algorithmId, idSeparator, algorithmVersion)
}

// validateAlgorithmVersionResponse guards the GetAlgorithmVersion response
// before it reaches state. A missing or wrong-typed AlgorithmId/AlgorithmVersion
// would otherwise compose a ":" id (rawString falls back to ""), a mismatched
// response could overwrite the resource's identity, and the caller would have
// no signal that the object is not the one it asked for. On any such failure
// the caller must keep the prior state and surface this error rather than
// writing a corrupt object.
func validateAlgorithmVersionResponse(object map[string]interface{}, expectedAlgorithmId, expectedAlgorithmVersion string) error {
	if object == nil {
		return fmt.Errorf("GetAlgorithmVersion returned an empty response for %s%s%s", expectedAlgorithmId, idSeparator, expectedAlgorithmVersion)
	}
	algorithmId, ok := object["AlgorithmId"].(string)
	if !ok || algorithmId == "" {
		return fmt.Errorf("GetAlgorithmVersion response is missing the string AlgorithmId for %s%s%s: %#v", expectedAlgorithmId, idSeparator, expectedAlgorithmVersion, object["AlgorithmId"])
	}
	algorithmVersion, ok := object["AlgorithmVersion"].(string)
	if !ok || algorithmVersion == "" {
		return fmt.Errorf("GetAlgorithmVersion response is missing the string AlgorithmVersion for %s%s%s: %#v", expectedAlgorithmId, idSeparator, expectedAlgorithmVersion, object["AlgorithmVersion"])
	}
	if algorithmId != expectedAlgorithmId || algorithmVersion != expectedAlgorithmVersion {
		return fmt.Errorf("GetAlgorithmVersion response identity %s%s%s does not match the requested %s%s%s", algorithmId, idSeparator, algorithmVersion, expectedAlgorithmId, idSeparator, expectedAlgorithmVersion)
	}
	return nil
}
