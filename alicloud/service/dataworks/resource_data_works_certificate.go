// Package dataworks. This file is generated automatically. Please do not modify it manually, thank you!
// Nested attributes require Terraform Plugin Protocol 6.
package dataworks

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const certificateTypeName = "alicloud_data_works_certificate"

var (
	_ resource.Resource                = &certificateResource{}
	_ resource.ResourceWithImportState = &certificateResource{}
)

func NewCertificateResource() resource.Resource {
	return &certificateResource{}
}

type certificateResource struct {
	fwadapt.ResourceBase
}

type certificateResourceModel struct {
	Id              types.String `tfsdk:"id"`
	FwTimeouts      types.Object `tfsdk:"timeouts"`
	CertificateFile types.String `tfsdk:"certificate_file"`
	CreateTime      types.String `tfsdk:"create_time"`
	CreateUser      types.String `tfsdk:"create_user"`
	Description     types.String `tfsdk:"description"`
	FileSizeInBytes types.Int64  `tfsdk:"file_size_in_bytes"`
	Name            types.String `tfsdk:"name"`
	ProjectId       types.Int64  `tfsdk:"project_id"`
}

func (r *certificateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = certificateTypeName
}

func (r *certificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{"timeouts": fwTimeoutBlock()},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"certificate_file": schema.StringAttribute{
				Sensitive: true,
				Required:  true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The URL of the certificate file. The object must be readable by DataWorks. Changing this value replaces the certificate.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time of the certificate file, in RFC3339 format.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_user": schema.StringAttribute{
				Computed:    true,
				Description: "The creator ID of the certificate file.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The description of the certificate file.",
			},
			"file_size_in_bytes": schema.Int64Attribute{
				Computed:    true,
				Description: "File size in bytes.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The name of the certificate file, unique within the workspace.",
			},
			"project_id": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Description: "The ID of the DataWorks workspace to which the certificate file belongs.",
			},
		},
	}
}

func (r *certificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan certificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operationTimeout, timeoutDiagnostics := fwTimeoutDuration(plan.FwTimeouts, "create", 5*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"ClientToken": buildClientToken("ImportCertificate"),
	}
	// gap#9: RPC APIs require RegionId as a request param (separate from endpoint resolution);
	// the SDK RpcPost resolves the endpoint by region but does NOT inject RegionId into query/body.
	// v1 sets request["RegionId"] = client.RegionId authoritatively (provider region wins).
	request["RegionId"] = r.Client().RegionId
	createValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !plan.CertificateFile.IsNull() && !plan.CertificateFile.IsUnknown() {
		createValues["CertificateFile"] = fwValueToInterface(plan.CertificateFile, map[string]fwWireField{}, false)
	}
	if !plan.CreateTime.IsNull() && !plan.CreateTime.IsUnknown() {
		createValues["CreateTime"] = fwValueToInterface(plan.CreateTime, map[string]fwWireField{}, false)
	}
	if !plan.CreateUser.IsNull() && !plan.CreateUser.IsUnknown() {
		createValues["CreateUser"] = fwValueToInterface(plan.CreateUser, map[string]fwWireField{}, false)
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		createValues["Description"] = fwValueToInterface(plan.Description, map[string]fwWireField{}, false)
	}
	if !plan.FileSizeInBytes.IsNull() && !plan.FileSizeInBytes.IsUnknown() {
		createValues["FileSizeInBytes"] = fwValueToInterface(plan.FileSizeInBytes, map[string]fwWireField{}, false)
	}
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		createValues["Name"] = fwValueToInterface(plan.Name, map[string]fwWireField{}, false)
	}
	if !plan.ProjectId.IsNull() && !plan.ProjectId.IsUnknown() {
		createValues["ProjectId"] = fwValueToInterface(plan.ProjectId, map[string]fwWireField{}, false)
	}
	if err := fwApplyMappings(request, createValues, []fwMapping{
		{Source: "$.CertificateFile", Target: "CertificateFile", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Description", Target: "Description", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Name", Target: "Name", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ProjectId", Target: "ProjectId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	}, false); err != nil {
		resp.Diagnostics.AddError("Building create request", err.Error())
		return
	}

	response, err := retryRpcPost(r.Client(), "dataworks-public", "2024-05-18", "ImportCertificate", request, operationTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling ImportCertificate", certificateTypeName),
			errs.WrapErrorf(err, "Creating %s: calling ImportCertificate", certificateTypeName).Error(),
		)
		return
	}

	identityParts := []string{}
	identityValues0 := fwWalkMappedPath(request, fwMappingPath("ProjectId"), nil)
	if len(identityValues0) == 0 || identityValues0[0].Value == nil || fmt.Sprint(identityValues0[0].Value) == "" {
		resp.Diagnostics.AddError("Creating alicloud_data_works_certificate: missing identity", "The create operation did not return or accept every identity component")
		return
	}
	identityParts = append(identityParts, fmt.Sprint(identityValues0[0].Value))
	identityValues1 := fwWalkMappedPath(response, fwMappingPath("$.Id"), nil)
	if len(identityValues1) == 0 || identityValues1[0].Value == nil || fmt.Sprint(identityValues1[0].Value) == "" {
		resp.Diagnostics.AddError("Creating alicloud_data_works_certificate: missing identity", "The create operation did not return or accept every identity component")
		return
	}
	identityParts = append(identityParts, fmt.Sprint(identityValues1[0].Value))
	id := strings.Join(identityParts, ":")
	if id == "" || id == "<nil>" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling ImportCertificate", certificateTypeName),
			fmt.Sprintf("The response carried no ProjectId: %v", response),
		)
		return
	}

	// Persist the allocated identity before any asynchronous or configuration stage.
	trackedPlan := plan
	trackedPlan.Id = types.StringValue(id)
	trackedPlan.CertificateFile = fwKnownValue(trackedPlan.CertificateFile).(types.String)
	trackedPlan.CreateTime = fwKnownValue(trackedPlan.CreateTime).(types.String)
	trackedPlan.CreateUser = fwKnownValue(trackedPlan.CreateUser).(types.String)
	trackedPlan.Description = fwKnownValue(trackedPlan.Description).(types.String)
	trackedPlan.FileSizeInBytes = fwKnownValue(trackedPlan.FileSizeInBytes).(types.Int64)
	trackedPlan.Name = fwKnownValue(trackedPlan.Name).(types.String)
	trackedPlan.ProjectId = fwKnownValue(trackedPlan.ProjectId).(types.Int64)
	resp.Diagnostics.Append(resp.State.Set(ctx, trackedPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := certificateGet(r.Client(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading after create", certificateTypeName),
			errs.WrapErrorf(err, "Creating %s: reading after create (%s)", certificateTypeName, id).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: after create", certificateTypeName),
			fmt.Sprintf("The resource %s is not visible via GetCertificate right after creation", id),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, certificateState(object, &plan))...)
}

func (r *certificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state certificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := certificateGet(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetCertificate", certificateTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetCertificate", certificateTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, certificateState(object, &state))...)
}

func (r *certificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state certificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// DataWorks has no update API: every configurable attribute is ForceNew, so Update is only
	// reached for the timeouts block. Re-read instead of copying the plan, whose unconfigured
	// computed attributes are unknown at apply time.
	object, err := certificateGet(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling GetCertificate", certificateTypeName),
			errs.WrapErrorf(err, "Updating %s: calling GetCertificate", certificateTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, certificateState(object, &state))...)
}

func (r *certificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state certificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operationTimeout, timeoutDiagnostics := fwTimeoutDuration(state.FwTimeouts, "delete", 5*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !state.CertificateFile.IsNull() && !state.CertificateFile.IsUnknown() {
		deleteValues["CertificateFile"] = fwValueToInterface(state.CertificateFile, map[string]fwWireField{}, false)
	}
	if !state.CreateTime.IsNull() && !state.CreateTime.IsUnknown() {
		deleteValues["CreateTime"] = fwValueToInterface(state.CreateTime, map[string]fwWireField{}, false)
	}
	if !state.CreateUser.IsNull() && !state.CreateUser.IsUnknown() {
		deleteValues["CreateUser"] = fwValueToInterface(state.CreateUser, map[string]fwWireField{}, false)
	}
	if !state.Description.IsNull() && !state.Description.IsUnknown() {
		deleteValues["Description"] = fwValueToInterface(state.Description, map[string]fwWireField{}, false)
	}
	if !state.FileSizeInBytes.IsNull() && !state.FileSizeInBytes.IsUnknown() {
		deleteValues["FileSizeInBytes"] = fwValueToInterface(state.FileSizeInBytes, map[string]fwWireField{}, false)
	}
	if !state.Name.IsNull() && !state.Name.IsUnknown() {
		deleteValues["Name"] = fwValueToInterface(state.Name, map[string]fwWireField{}, false)
	}
	if !state.ProjectId.IsNull() && !state.ProjectId.IsUnknown() {
		deleteValues["ProjectId"] = fwValueToInterface(state.ProjectId, map[string]fwWireField{}, false)
	}
	identityParts := strings.Split(state.Id.ValueString(), ":")
	if len(identityParts) > 0 {
		deleteValues["ProjectId"] = identityParts[0]
	}
	if len(identityParts) > 1 {
		deleteValues["Id"] = identityParts[1]
	}
	selected := false
	deleteEligible0, deleteConditionErr0 := fwDeleteEligible(deleteValues, []fwDeleteCondition{})
	if deleteConditionErr0 != nil {
		resp.Diagnostics.AddError("Selecting deletion operation", deleteConditionErr0.Error())
		return
	}
	selected = selected || deleteEligible0
	if !selected {
		resp.Diagnostics.AddError("Selecting deletion operation", "No deletion operation matches the resource state; the resource remains tracked")
		return
	}
	if deleteEligible0 {
		request0 := map[string]interface{}{"ClientToken": buildClientToken("DeleteCertificate")}
		request0["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request0, deleteValues, []fwMapping{
			{Source: "$.Id", Target: "Id", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.ProjectId", Target: "ProjectId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, false); err != nil {
			resp.Diagnostics.AddError("Building DeleteCertificate request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "dataworks-public", "2024-05-18", "DeleteCertificate", request0, operationTimeout); err != nil && !isExpectedErrors(err, []string{"ResourceNotFound.Certificate", "ResourceNotFound"}) {
			resp.Diagnostics.AddError("Deleting alicloud_data_works_certificate: DeleteCertificate", err.Error())
			return
		}
	}
}

func (r *certificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// certificateState maps the Get response object onto the resource model. Fields present in the
// response are read via their responseKey (e.g. status -> ListenerStatus); request-only fields
// (responseKey empty, e.g. dry_run) are preserved from the prior plan/state instead.
func certificateState(object map[string]interface{}, prior *certificateResourceModel) certificateResourceModel {
	// Get has already merged and normalized every response mapping.
	str := func(key string) types.String {
		cur := interface{}(object)
		for _, seg := range strings.Split(key, ".") {
			m, ok := cur.(map[string]interface{})
			if !ok {
				return types.StringNull()
			}
			cur, ok = m[seg]
			if !ok {
				return types.StringNull()
			}
		}
		// gap#7c: present-but-empty strings must return ValueString(""), not Null.
		// v1 d.Set(field, "") treats empty as set (not absence); TestCheckResourceAttrSet
		// expects a non-null value for API-returned "" (e.g. DhcpOptionsSetId when VPC has
		// no DHCP set). The earlier && v != "" collapsed present-empty to Null, failing Set
		// checks. Absent fields still return Null via the !ok branches above.
		if v, ok := cur.(string); ok {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	intVal := func(key string) types.Int64 {
		cur := interface{}(object)
		for _, seg := range strings.Split(key, ".") {
			m, ok := cur.(map[string]interface{})
			if !ok {
				return types.Int64Null()
			}
			cur, ok = m[seg]
			if !ok {
				return types.Int64Null()
			}
		}
		switch v := cur.(type) {
		case float64:
			return types.Int64Value(int64(v))
		case int:
			return types.Int64Value(int64(v))
		case int64:
			return types.Int64Value(v)
		case json.Number:
			if n, err := v.Int64(); err == nil {
				return types.Int64Value(n)
			}
		case string:
			if n, err := strconv.ParseInt(v, 10, 64); err == nil {
				return types.Int64Value(n)
			}
		}
		return types.Int64Null()
	}
	// gap#11: composite-ID resource — build the id from all composite fields in the response,
	// joined by ":" (e.g. AccessGroupName:FileSystemType). The Get response carries every
	// part, so the state id round-trips through Read/Delete correctly (otherwise Read would
	// lose FileSystemType and DescribeAccessGroups would query the wrong namespace → 404).
	_cidParts := []string{
		str("ProjectId").ValueString(),
		str("Id").ValueString(),
	}
	m := certificateResourceModel{
		FwTimeouts:      types.ObjectNull(fwTimeoutTypes()),
		Id:              types.StringValue(strings.Join(_cidParts, ":")),
		CreateTime:      certificateCreateTime(object["CreateTime"]),
		CreateUser:      str("CreateUser"),
		Description:     str("Description"),
		FileSizeInBytes: intVal("FileSizeInBytes"),
		Name:            str("Name"),
		ProjectId:       intVal("ProjectId"),
	}
	if prior != nil {
		m.FwTimeouts = prior.FwTimeouts
		m.CertificateFile = prior.CertificateFile
	}
	return m
}

// certificateCreateTime renders the millisecond timestamp from GetCertificate as the RFC3339
// string the CloudSpec property declares. Hand-fix: the generated string mapper drops the raw
// int64 value, which would leave create_time unset and drift on every plan.
func certificateCreateTime(raw interface{}) types.String {
	switch v := raw.(type) {
	case nil:
		return types.StringNull()
	case float64:
		return types.StringValue(time.UnixMilli(int64(v)).UTC().Format(time.RFC3339))
	case int64:
		return types.StringValue(time.UnixMilli(v).UTC().Format(time.RFC3339))
	case int:
		return types.StringValue(time.UnixMilli(int64(v)).UTC().Format(time.RFC3339))
	case json.Number:
		if ms, err := v.Int64(); err == nil {
			return types.StringValue(time.UnixMilli(ms).UTC().Format(time.RFC3339))
		}
		return types.StringNull()
	case string:
		if v == "" {
			return types.StringValue("")
		}
		if ms, err := strconv.ParseInt(v, 10, 64); err == nil {
			return types.StringValue(time.UnixMilli(ms).UTC().Format(time.RFC3339))
		}
		return types.StringValue(v)
	default:
		return types.StringNull()
	}
}
