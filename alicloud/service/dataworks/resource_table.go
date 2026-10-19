package dataworks

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const tableTypeName = "alicloud_dataworks_table"

// dataworks-public is the POP product for DataWorks; its RPC actions are called through the
// shared AliyunClient. The service package cannot import the root alicloud package, so retry,
// jsonpath and tea are used directly instead of incrementalWait/NeedRetry/WrapError.
const (
	dwProduct = "dataworks-public"
	dwVersion = "2020-05-18"
)

var (
	_ resource.Resource                = &tableResource{}
	_ resource.ResourceWithImportState = &tableResource{}
)

func NewTableResource() resource.Resource {
	return &tableResource{}
}

type tableResource struct {
	fwadapt.ResourceBase
}

// tableResourceModel mirrors the ACS::DataWorks::Table (dataworks-public 2020-05-18) pre
// resource model. The identity is TableGuid (required, composite form
// "odps.<engine>.<table>"); the plain table name is parsed from it for the TableName
// parameter of CreateTable/UpdateTable/DeleteTable, while the full TableGuid is used
// for GetMetaTableFullInfo. IsView/CategoryId/PhysicsLevelId are write-only (the read API
// does not return them), so they are preserved from state on Read to avoid drift.
type tableResourceModel struct {
	// Identity / required inputs.
	TableGuid types.String `tfsdk:"table_guid"`
	ProjectId types.Int64  `tfsdk:"project_id"`

	// Writable attributes returned by GetMetaTableFullInfo.
	Comment   types.String `tfsdk:"comment"`
	LifeCycle types.Int64  `tfsdk:"life_cycle"`
	Location  types.String `tfsdk:"location"`
	OwnerId   types.String `tfsdk:"owner_id"`

	// Write-only attributes (AC2): sent on create/update, never returned by the read API.
	IsView         types.Int64 `tfsdk:"is_view"`
	CategoryId     types.Int64 `tfsdk:"category_id"`
	PhysicsLevelId types.Int64 `tfsdk:"physics_level_id"`

	// AC3 resource-operation driving attributes (also write-only): drive
	// UpdateMetaTableIntroWiki and UpdateTableModelInfo, which are not part of the model
	// schema but required callable by the route plan.
	IntroWiki          types.String `tfsdk:"intro_wiki"`
	FirstLevelThemeId  types.Int64  `tfsdk:"first_level_theme_id"`
	SecondLevelThemeId types.Int64  `tfsdk:"second_level_theme_id"`
	LevelId            types.Int64  `tfsdk:"level_id"`
	LevelType          types.String `tfsdk:"level_type"`

	// Computed / read-only.
	Id             types.String `tfsdk:"id"`
	TableName      types.String `tfsdk:"table_name"`
	CreateTime     types.String `tfsdk:"create_time"`
	LastModifyTime types.String `tfsdk:"last_modify_time"`
	TenantId       types.Int64  `tfsdk:"tenant_id"`
	DataSize       types.Int64  `tfsdk:"data_size"`
	IsVisible      types.Int64  `tfsdk:"is_visible"`
	RegionId       types.String `tfsdk:"region_id"`
}

func (r *tableResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = tableTypeName
}

func (r *tableResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a DataWorks Table resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource, same as `table_guid`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"table_guid": schema.StringAttribute{
				Required: true,
				Description: "The composite guid of the table, in the form `odps.<engine>.<table_name>`. " +
					"The plain `<table_name>` segment is used as the TableName of CreateTable/UpdateTable/DeleteTable; " +
					"the full guid is used to read the table via GetMetaTableFullInfo.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.Int64Attribute{
				Required:    true,
				Description: "The ID of the DataWorks workspace that owns the table.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Description: "The comment of the table.",
			},
			"life_cycle": schema.Int64Attribute{
				Optional:    true,
				Description: "The lifecycle of the table, in days.",
			},
			"location": schema.StringAttribute{
				Optional:    true,
				Description: "The storage location of the external table.",
			},
			"owner_id": schema.StringAttribute{
				Optional:    true,
				Description: "The owner ID of the table.",
			},
			"is_view": schema.Int64Attribute{
				Optional:    true,
				Description: "Whether the table is a view. 0: table, 1: view. Write-only: not returned by the read API.",
			},
			"category_id": schema.Int64Attribute{
				Optional:    true,
				Description: "The category ID of the table. Write-only: not returned by the read API.",
			},
			"physics_level_id": schema.Int64Attribute{
				Optional:    true,
				Description: "The physics level ID of the table. Write-only: not returned by the read API.",
			},
			"intro_wiki": schema.StringAttribute{
				Optional:    true,
				Description: "The intro wiki content of the table. Drives UpdateMetaTableIntroWiki; write-only.",
			},
			"first_level_theme_id": schema.Int64Attribute{
				Optional:    true,
				Description: "The first-level theme ID. Drives UpdateTableModelInfo; write-only.",
			},
			"second_level_theme_id": schema.Int64Attribute{
				Optional:    true,
				Description: "The second-level theme ID. Drives UpdateTableModelInfo; write-only.",
			},
			"level_id": schema.Int64Attribute{
				Optional:    true,
				Description: "The model level ID. Drives UpdateTableModelInfo; write-only.",
			},
			"level_type": schema.StringAttribute{
				Optional:    true,
				Description: "The model level type. Drives UpdateTableModelInfo; write-only.",
			},
			"table_name": schema.StringAttribute{
				Computed:    true,
				Description: "The plain name of the table.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the table was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_modify_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the table was last modified.",
				// No UseStateForUnknown: UpdateTable bumps the timestamp, so the planned
				// value can never survive an apply.
			},
			"tenant_id": schema.Int64Attribute{
				Computed:    true,
				Description: "The tenant ID.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"data_size": schema.Int64Attribute{
				Computed:    true,
				Description: "The data size of the table.",
				// No UseStateForUnknown: the value moves outside provider control as data loads.
			},
			"is_visible": schema.Int64Attribute{
				Computed:    true,
				Description: "Whether the table is visible.",
			},
			"region_id": schema.StringAttribute{
				Computed:    true,
				Description: "The region ID of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *tableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tableGuid := plan.TableGuid.ValueString()
	tableName, err := parseTableName(tableGuid)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: parsing table_guid", tableTypeName),
			err.Error(),
		)
		return
	}

	request := map[string]interface{}{
		"TableName": tableName,
		"ProjectId": plan.ProjectId.ValueInt64(),
	}
	setCreateWriteFields(request, &plan)

	if _, err := r.Client().RpcPost(dwProduct, dwVersion, "CreateTable", nil, request, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateTable", tableTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateTable", tableTypeName).Error(),
		)
		return
	}

	// CreateTable is asynchronous (task based): poll GetMetaTableFullInfo until the table
	// is visible, then the composite TableGuid is the read identity.
	if _, err := waitTableReady(r.Client(), tableGuid); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for CreateTable", tableTypeName),
			errs.WrapErrorf(err, "Creating %s: table %s did not become visible", tableTypeName, tableGuid).Error(),
		)
		return
	}

	// AC3: resource-operation APIs are callable when their driving attributes are set.
	if err := applyTableOps(r.Client(), &plan, tableGuid, "create"); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling resource operation API", tableTypeName),
			err.Error(),
		)
		return
	}

	object, err := getMetaTableFullInfo(r.Client(), tableGuid)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetMetaTableFullInfo", tableTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetMetaTableFullInfo", tableTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateTable", tableTypeName),
			fmt.Sprintf("The table %s is not visible via GetMetaTableFullInfo right after creation", tableGuid),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, tableState(object, &plan, r.Client().RegionId))...)
}

func (r *tableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// ImportStatePassthroughID populates `id` only, so on a cold import `table_guid` is
	// empty; id always equals table_guid, so fall back to it.
	guid := state.TableGuid.ValueString()
	if guid == "" {
		guid = state.Id.ValueString()
	}
	object, err := getMetaTableFullInfo(r.Client(), guid)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetMetaTableFullInfo", tableTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetMetaTableFullInfo", tableTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, tableState(object, &state, r.Client().RegionId))...)
}

func (r *tableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state tableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tableGuid := plan.TableGuid.ValueString()
	tableName, err := parseTableName(tableGuid)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: parsing table_guid", tableTypeName),
			err.Error(),
		)
		return
	}

	request := map[string]interface{}{
		"TableName": tableName,
		"ProjectId": plan.ProjectId.ValueInt64(),
	}
	// Only the changed writable attributes are sent. Optional string attributes send an empty
	// string to clear, matching the legacy update behavior; integer attributes are sent only
	// when set, since a 0 is ambiguous for minimum-1 identifiers such as CategoryId.
	if !plan.Comment.Equal(state.Comment) {
		request["Comment"] = plan.Comment.ValueString()
	}
	if !plan.Location.Equal(state.Location) {
		request["Location"] = plan.Location.ValueString()
	}
	if !plan.OwnerId.Equal(state.OwnerId) {
		request["OwnerId"] = plan.OwnerId.ValueString()
	}
	if !plan.LifeCycle.Equal(state.LifeCycle) && !plan.LifeCycle.IsNull() {
		request["LifeCycle"] = plan.LifeCycle.ValueInt64()
	}
	if !plan.IsView.Equal(state.IsView) && !plan.IsView.IsNull() {
		request["IsView"] = plan.IsView.ValueInt64()
	}
	if !plan.CategoryId.Equal(state.CategoryId) && !plan.CategoryId.IsNull() {
		request["CategoryId"] = plan.CategoryId.ValueInt64()
	}
	if !plan.PhysicsLevelId.Equal(state.PhysicsLevelId) && !plan.PhysicsLevelId.IsNull() {
		request["PhysicsLevelId"] = plan.PhysicsLevelId.ValueInt64()
	}

	if len(request) > 2 {
		if _, err := r.Client().RpcPost(dwProduct, dwVersion, "UpdateTable", nil, request, true); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling UpdateTable", tableTypeName),
				errs.WrapErrorf(err, "Updating %s: calling UpdateTable", tableTypeName).Error(),
			)
			return
		}
	}

	// AC3: resource-operation APIs are callable when their driving attributes change.
	if err := applyTableOps(r.Client(), &plan, tableGuid, "update"); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling resource operation API", tableTypeName),
			err.Error(),
		)
		return
	}

	object, err := getMetaTableFullInfo(r.Client(), tableGuid)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetMetaTableFullInfo", tableTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetMetaTableFullInfo", tableTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, tableState(object, &plan, r.Client().RegionId))...)
}

func (r *tableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// id always equals table_guid; fall back to it for robustness.
	guid := state.TableGuid.ValueString()
	if guid == "" {
		guid = state.Id.ValueString()
	}
	tableName, err := parseTableName(guid)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: parsing table_guid", tableTypeName),
			err.Error(),
		)
		return
	}
	_, err = r.Client().RpcPost(dwProduct, dwVersion, "DeleteTable", nil, map[string]interface{}{
		"TableName": tableName,
		"ProjectId": state.ProjectId.ValueInt64(),
	}, true)
	if err != nil && !isTableNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteTable", tableTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteTable", tableTypeName).Error(),
		)
		return
	}
}

func (r *tableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// parseTableName extracts the plain table name from a composite TableGuid of the form
// "odps.<engine>.<table_name>" (optionally "odps.<engine>.<schema>.<table_name>" for the
// three-layer model). CreateTable/UpdateTable/DeleteTable expect the plain name; the full
// guid is reserved for GetMetaTableFullInfo.
func parseTableName(tableGuid string) (string, error) {
	if tableGuid == "" {
		return "", fmt.Errorf("table_guid must be in the form odps.<engine>.<table_name>")
	}
	parts := strings.Split(tableGuid, ".")
	if len(parts) < 3 || parts[0] != "odps" {
		return "", fmt.Errorf("table_guid %q must be in the form odps.<engine>.<table_name>", tableGuid)
	}
	name := parts[len(parts)-1]
	if name == "" {
		return "", fmt.Errorf("table_guid %q has an empty table name segment", tableGuid)
	}
	return name, nil
}

// setCreateWriteFields populates the optional writable attributes on a CreateTable request.
func setCreateWriteFields(request map[string]interface{}, plan *tableResourceModel) {
	if !plan.Comment.IsNull() {
		request["Comment"] = plan.Comment.ValueString()
	}
	if !plan.LifeCycle.IsNull() {
		request["LifeCycle"] = plan.LifeCycle.ValueInt64()
	}
	if !plan.Location.IsNull() {
		request["Location"] = plan.Location.ValueString()
	}
	if !plan.OwnerId.IsNull() {
		request["OwnerId"] = plan.OwnerId.ValueString()
	}
	if !plan.IsView.IsNull() {
		request["IsView"] = plan.IsView.ValueInt64()
	}
	if !plan.CategoryId.IsNull() {
		request["CategoryId"] = plan.CategoryId.ValueInt64()
	}
	if !plan.PhysicsLevelId.IsNull() {
		request["PhysicsLevelId"] = plan.PhysicsLevelId.ValueInt64()
	}
}

// applyTableOps calls the AC3 resource-operation APIs (UpdateMetaTableIntroWiki and
// UpdateTableModelInfo) when their driving attributes are set. These attributes are not
// returned by GetMetaTableFullInfo, so they are write-only.
func applyTableOps(client *connectivity.AliyunClient, plan *tableResourceModel, tableGuid, op string) error {
	if !plan.IntroWiki.IsNull() {
		req := map[string]interface{}{
			"TableGuid": tableGuid,
			"Content":   plan.IntroWiki.ValueString(),
		}
		if _, err := client.RpcPost(dwProduct, dwVersion, "UpdateMetaTableIntroWiki", nil, req, true); err != nil {
			return errs.WrapErrorf(err, "%s %s: calling UpdateMetaTableIntroWiki", op, tableTypeName)
		}
	}
	if !plan.LevelId.IsNull() || !plan.LevelType.IsNull() ||
		!plan.FirstLevelThemeId.IsNull() || !plan.SecondLevelThemeId.IsNull() {
		req := map[string]interface{}{"TableGuid": tableGuid}
		if !plan.FirstLevelThemeId.IsNull() {
			req["FirstLevelThemeId"] = plan.FirstLevelThemeId.ValueInt64()
		}
		if !plan.SecondLevelThemeId.IsNull() {
			req["SecondLevelThemeId"] = plan.SecondLevelThemeId.ValueInt64()
		}
		if !plan.LevelId.IsNull() {
			req["LevelId"] = plan.LevelId.ValueInt64()
		}
		if !plan.LevelType.IsNull() {
			req["LevelType"] = plan.LevelType.ValueString()
		}
		if _, err := client.RpcPost(dwProduct, dwVersion, "UpdateTableModelInfo", nil, req, true); err != nil {
			return errs.WrapErrorf(err, "%s %s: calling UpdateTableModelInfo", op, tableTypeName)
		}
	}
	return nil
}

// waitTableReady polls GetMetaTableFullInfo until the table is visible or the timeout elapses.
func waitTableReady(client *connectivity.AliyunClient, tableGuid string) (map[string]interface{}, error) {
	var object map[string]interface{}
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		obj, err := getMetaTableFullInfo(client, tableGuid)
		if err != nil {
			if isRetryableError(err) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if obj == nil {
			return retry.RetryableError(fmt.Errorf("table %s is not yet visible", tableGuid))
		}
		object = obj
		return nil
	})
	if err != nil {
		return nil, err
	}
	return object, nil
}

// getMetaTableFullInfo reads a table by its full TableGuid. It returns (nil, nil) when the
// table does not exist so the caller can decide between RemoveResource and a retry.
func getMetaTableFullInfo(client *connectivity.AliyunClient, tableGuid string) (map[string]interface{}, error) {
	response, err := client.RpcPost(dwProduct, dwVersion, "GetMetaTableFullInfo", nil, map[string]interface{}{
		"TableGuid": tableGuid,
	}, true)
	if err != nil {
		if isTableNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return nil, fmt.Errorf("GetMetaTableFullInfo response carried no $.Data: %v", response)
	}
	object, ok := v.(map[string]interface{})
	if !ok || len(object) < 1 {
		return nil, nil
	}
	return object, nil
}

// tableState maps a GetMetaTableFullInfo $.Data object onto the resource model. Write-only
// attributes (IsView, CategoryId, PhysicsLevelId and the AC3 operation attributes) are not
// returned by the read API, so they are carried over from the prior state/plan to keep the
// state drift-free (AC2).
func tableState(object map[string]interface{}, prior *tableResourceModel, regionId string) tableResourceModel {
	str := func(key string, fallback types.String) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return fallback
	}
	i64 := func(key string, fallback types.Int64) types.Int64 {
		switch v := object[key].(type) {
		case float64:
			return types.Int64Value(int64(v))
		case int64:
			return types.Int64Value(v)
		case int:
			return types.Int64Value(int64(v))
		case json.Number:
			if n, e := v.Int64(); e == nil {
				return types.Int64Value(n)
			}
		}
		return fallback
	}

	// Identity is stable; prefer the response value and fall back to the prior guid.
	guid := str("TableGuid", prior.TableGuid)
	if guid.IsNull() {
		guid = prior.TableGuid
	}

	return tableResourceModel{
		Id:        guid,
		TableGuid: guid,

		ProjectId: i64("ProjectId", prior.ProjectId),
		Comment:   str("Comment", prior.Comment),
		LifeCycle: i64("LifeCycle", prior.LifeCycle),
		Location:  str("Location", prior.Location),
		OwnerId:   str("OwnerId", prior.OwnerId),

		// Write-only: preserved from prior state/plan.
		IsView:             prior.IsView,
		CategoryId:         prior.CategoryId,
		PhysicsLevelId:     prior.PhysicsLevelId,
		IntroWiki:          prior.IntroWiki,
		FirstLevelThemeId:  prior.FirstLevelThemeId,
		SecondLevelThemeId: prior.SecondLevelThemeId,
		LevelId:            prior.LevelId,
		LevelType:          prior.LevelType,

		// Computed from the read API.
		TableName:      str("TableName", prior.TableName),
		CreateTime:     str("CreateTime", prior.CreateTime),
		LastModifyTime: str("LastModifyTime", types.StringNull()),
		TenantId:       i64("TenantId", prior.TenantId),
		DataSize:       i64("DataSize", types.Int64Null()),
		IsVisible:      i64("IsVisible", types.Int64Null()),
		RegionId:       types.StringValue(regionId),
	}
}

// isTableNotExist reports whether err is the DataWorks table-not-found error. The code arrives
// on a *tea.SDKError from RpcPost; the substring fallback covers wrapped forms.
func isTableNotExist(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil && *e.Code == "InvalidParameter.Meta.NotExist" {
		return true
	}
	return strings.Contains(err.Error(), "InvalidParameter.Meta.NotExist")
}

// isRetryableError reports whether err is a throttling or internal-system error worth retrying.
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil {
		code := *e.Code
		if strings.HasPrefix(code, "Throttling") || strings.HasPrefix(code, "InternalError") {
			return true
		}
	}
	msg := err.Error()
	return strings.Contains(msg, "Throttling") || strings.Contains(msg, "InternalError") ||
		strings.Contains(msg, "try again")
}
