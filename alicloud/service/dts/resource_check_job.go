package dts

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	checkJobTypeName = "alicloud_dts_check_job"

	checkJobProduct    = "Dts"
	checkJobAPIVersion = "2020-01-01"

	// terraform-plugin-framework has no user-configurable resource timeouts, so
	// the envelopes the SDKv2 implementation derived from d.Timeout are fixed
	// here: 20 minutes for create (purchase, configure and the first status
	// settle), 10 minutes for update and delete as before.
	checkJobCreateTimeout = 20 * time.Minute
	checkJobUpdateTimeout = 10 * time.Minute
	checkJobDeleteTimeout = 10 * time.Minute
)

var (
	_ resource.Resource                = &checkJobResource{}
	_ resource.ResourceWithImportState = &checkJobResource{}
)

// NewCheckJobResource is the native Framework factory registered through
// ServicePackage.Resources.
func NewCheckJobResource() resource.Resource {
	return &checkJobResource{}
}

// checkJobResource is the native Framework implementation of the DTS check job.
type checkJobResource struct {
	fwadapt.ResourceBase
}

func (r *checkJobResource) describe(id string) (map[string]interface{}, error) {
	return NewService(r.Client()).DescribeCheckJob(id)
}

type checkJobResourceModel struct {
	Id                              types.String `tfsdk:"id"`
	DtsInstanceId                   types.String `tfsdk:"dts_instance_id"`
	DtsJobName                      types.String `tfsdk:"dts_job_name"`
	DedicatedClusterId              types.String `tfsdk:"dedicated_cluster_id"`
	DataCheckConfigure              types.String `tfsdk:"data_check_configure"`
	DataInitialization              types.Bool   `tfsdk:"data_initialization"`
	DataSynchronization             types.Bool   `tfsdk:"data_synchronization"`
	StructureInitialization         types.Bool   `tfsdk:"structure_initialization"`
	Checkpoint                      types.String `tfsdk:"checkpoint"`
	InstanceClass                   types.String `tfsdk:"instance_class"`
	DbList                          types.String `tfsdk:"db_list"`
	Reserve                         types.String `tfsdk:"reserve"`
	SourceEndpointInstanceType      types.String `tfsdk:"source_endpoint_instance_type"`
	SourceEndpointEngineName        types.String `tfsdk:"source_endpoint_engine_name"`
	SourceEndpointInstanceId        types.String `tfsdk:"source_endpoint_instance_id"`
	SourceEndpointRegion            types.String `tfsdk:"source_endpoint_region"`
	SourceEndpointIp                types.String `tfsdk:"source_endpoint_ip"`
	SourceEndpointPort              types.String `tfsdk:"source_endpoint_port"`
	SourceEndpointOracleSid         types.String `tfsdk:"source_endpoint_oracle_sid"`
	SourceEndpointDatabaseName      types.String `tfsdk:"source_endpoint_database_name"`
	SourceEndpointUserName          types.String `tfsdk:"source_endpoint_user_name"`
	SourceEndpointPassword          types.String `tfsdk:"source_endpoint_password"`
	SourceEndpointOwnerId           types.String `tfsdk:"source_endpoint_owner_id"`
	SourceEndpointRole              types.String `tfsdk:"source_endpoint_role"`
	SourceEndpointVswitchId         types.String `tfsdk:"source_endpoint_vswitch_id"`
	SrcPrimaryVswitchId             types.String `tfsdk:"src_primary_vswitch_id"`
	SrcSecondaryVswitchId           types.String `tfsdk:"src_secondary_vswitch_id"`
	SourceEndpointSsl               types.String `tfsdk:"source_endpoint_ssl"`
	DestinationEndpointInstanceType types.String `tfsdk:"destination_endpoint_instance_type"`
	DestinationEndpointEngineName   types.String `tfsdk:"destination_endpoint_engine_name"`
	DestinationEndpointInstanceId   types.String `tfsdk:"destination_endpoint_instance_id"`
	DestinationEndpointRegion       types.String `tfsdk:"destination_endpoint_region"`
	DestinationEndpointIp           types.String `tfsdk:"destination_endpoint_ip"`
	DestinationEndpointPort         types.String `tfsdk:"destination_endpoint_port"`
	DestinationEndpointDatabaseName types.String `tfsdk:"destination_endpoint_database_name"`
	DestinationEndpointUserName     types.String `tfsdk:"destination_endpoint_user_name"`
	DestinationEndpointPassword     types.String `tfsdk:"destination_endpoint_password"`
	DestinationEndpointOracleSid    types.String `tfsdk:"destination_endpoint_oracle_sid"`
	DestinationEndpointOwnerId      types.String `tfsdk:"destination_endpoint_owner_id"`
	DestinationEndpointRole         types.String `tfsdk:"destination_endpoint_role"`
	DestinationEndpointSsl          types.String `tfsdk:"destination_endpoint_ssl"`
	DestPrimaryVswitchId            types.String `tfsdk:"dest_primary_vswitch_id"`
	DestSecondaryVswitchId          types.String `tfsdk:"dest_secondary_vswitch_id"`
	DtsBisLabel                     types.String `tfsdk:"dts_bis_label"`
	DelayNotice                     types.Bool   `tfsdk:"delay_notice"`
	DelayPhone                      types.String `tfsdk:"delay_phone"`
	DelayRuleTime                   types.String `tfsdk:"delay_rule_time"`
	ErrorNotice                     types.Bool   `tfsdk:"error_notice"`
	ErrorPhone                      types.String `tfsdk:"error_phone"`
	Status                          types.String `tfsdk:"status"`
}

func (r *checkJobResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = checkJobTypeName
}

func (r *checkJobResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dts_instance_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dts_job_name": schema.StringAttribute{Required: true},
			// A check job is purchased with a fixed class and cannot be resized in
			// place: TransferInstanceClass answers Success but the reported class
			// never changes, so the attribute is replacement-only.
			"instance_class": checkJobComputedReplaceableString(
				stringvalidator.OneOf("4xlarge", "2xlarge", "xlarge", "large", "medium", "small"),
			),
			"db_list":              checkJobRequiredReplaceableString(),
			"data_check_configure": checkJobReplaceableString(),
			"dedicated_cluster_id": checkJobReplaceableString(),
			"checkpoint":           checkJobComputedReplaceableString(),
			"reserve":              checkJobComputedReplaceableString(),
			"source_endpoint_instance_type": checkJobRequiredCaseInsensitiveString(
				stringvalidator.OneOfCaseInsensitive("CEN", "DG", "DISTRIBUTED_DMSLOGICDB", "ECS", "EXPRESS", "MONGODB", "OTHER", "PolarDB", "POLARDBX20", "RDS", "REDIS", "DISTRIBUTED_POLARDBX10"),
			),
			"source_endpoint_engine_name": checkJobRequiredCaseInsensitiveString(
				stringvalidator.OneOfCaseInsensitive("AS400", "DB2", "DMSPOLARDB", "HBASE", "MONGODB", "MSSQL", "MySQL", "ORACLE", "PolarDB", "POLARDBX20", "POLARDB_O", "POSTGRESQL", "TERADATA", "POLARDB_PG", "MARIADB", "POLARDBX10", "TiDB", "REDIS"),
			),
			"source_endpoint_instance_id":   checkJobComputedReplaceableString(),
			"source_endpoint_region":        checkJobComputedReplaceableString(),
			"source_endpoint_ip":            checkJobComputedReplaceableString(),
			"source_endpoint_port":          checkJobComputedReplaceableString(),
			"source_endpoint_oracle_sid":    checkJobComputedReplaceableString(),
			"source_endpoint_database_name": checkJobComputedReplaceableString(),
			"source_endpoint_user_name":     checkJobComputedReplaceableString(),
			"source_endpoint_password":      checkJobSensitiveString(),
			"source_endpoint_owner_id":      checkJobComputedReplaceableString(),
			"source_endpoint_role":          checkJobComputedReplaceableString(),
			"destination_endpoint_instance_type": checkJobRequiredCaseInsensitiveString(
				stringvalidator.OneOfCaseInsensitive("ads", "ADS", "CEN", "DATAHUB", "DG", "ECS", "EXPRESS", "GREENPLUM", "MONGODB", "OTHER", "PolarDB", "POLARDBX20", "RDS", "REDIS", "ELK", "Tablestore", "ODPS"),
			),
			"destination_endpoint_engine_name": checkJobRequiredCaseInsensitiveString(
				stringvalidator.OneOfCaseInsensitive("ADB20", "ADS", "ADB30", "AS400", "DATAHUB", "DB2", "GREENPLUM", "KAFKA", "MONGODB", "MSSQL", "MySQL", "ORACLE", "PolarDB", "POLARDBX20", "POLARDB_O", "PostgreSQL", "POLARDB_PG", "MARIADB", "POLARDBX10", "ODPS", "Tablestore", "ELK", "REDIS"),
			),
			"destination_endpoint_instance_id":   checkJobComputedReplaceableString(),
			"destination_endpoint_region":        checkJobComputedReplaceableString(),
			"destination_endpoint_ip":            checkJobComputedReplaceableString(),
			"destination_endpoint_port":          checkJobComputedReplaceableString(),
			"destination_endpoint_database_name": checkJobComputedReplaceableString(),
			"destination_endpoint_user_name":     checkJobComputedReplaceableString(),
			"destination_endpoint_password":      checkJobSensitiveString(),
			"destination_endpoint_oracle_sid":    checkJobComputedReplaceableString(),
			"destination_endpoint_owner_id":      checkJobComputedReplaceableString(),
			"destination_endpoint_role":          checkJobComputedReplaceableString(),
			"source_endpoint_vswitch_id":         checkJobReplaceableString(),
			"src_primary_vswitch_id":             checkJobReplaceableString(),
			"src_secondary_vswitch_id":           checkJobReplaceableString(),
			"dest_primary_vswitch_id":            checkJobReplaceableString(),
			"dest_secondary_vswitch_id":          checkJobReplaceableString(),
			"dts_bis_label":                      checkJobReplaceableString(),
			"delay_phone":                        checkJobReplaceableString(),
			"delay_rule_time":                    checkJobReplaceableString(),
			"error_phone":                        checkJobReplaceableString(),
			"delay_notice":                       checkJobReplaceableBool(),
			"error_notice":                       checkJobReplaceableBool(),
			// SSL travels inside Reserve, which ConfigureDtsJob accepts once:
			// ModifyDtsJob with ModifyTypeEnum=UPDATE_RESERVED answers
			// Forbidden.InstanceNotFound for CHECK instances, so a changed value
			// can only take effect on a new job.
			"source_endpoint_ssl": checkJobComputedReplaceableString(
				stringvalidator.OneOf("0", "1"),
			),
			"destination_endpoint_ssl": checkJobComputedReplaceableString(
				stringvalidator.OneOf("0", "1", "3"),
			),
			"data_initialization":      schema.BoolAttribute{Computed: true},
			"data_synchronization":     schema.BoolAttribute{Computed: true},
			"structure_initialization": schema.BoolAttribute{Computed: true},
			"status": schema.StringAttribute{
				Optional:   true,
				Computed:   true,
				Validators: []validator.String{stringvalidator.OneOf("Migrating", "Suspending")},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					checkJobFinishedStatusPlanModifier{},
				},
			},
		},
	}
}

func checkJobRequiredReplaceableString() schema.StringAttribute {
	return schema.StringAttribute{
		Required:      true,
		PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
	}
}

func checkJobRequiredCaseInsensitiveString(validators ...validator.String) schema.StringAttribute {
	return schema.StringAttribute{
		Required:   true,
		Validators: validators,
		PlanModifiers: []planmodifier.String{
			checkJobCaseInsensitivePlanModifier{},
			stringplanmodifier.RequiresReplace(),
		},
	}
}

func checkJobReplaceableString() schema.StringAttribute {
	return schema.StringAttribute{
		Optional:      true,
		PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
	}
}

func checkJobComputedReplaceableString(validators ...validator.String) schema.StringAttribute {
	return schema.StringAttribute{
		Optional:   true,
		Computed:   true,
		Validators: validators,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
			stringplanmodifier.RequiresReplace(),
		},
	}
}

func checkJobSensitiveString() schema.StringAttribute {
	return schema.StringAttribute{
		Optional:  true,
		Sensitive: true,
	}
}

func checkJobReplaceableBool() schema.BoolAttribute {
	return schema.BoolAttribute{
		Optional:      true,
		PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
	}
}

// checkJobCaseInsensitivePlanModifier keeps the configured spelling when it
// differs from state only by case, matching the DTS endpoint type and engine
// name values that the API upper-cases on read.
type checkJobCaseInsensitivePlanModifier struct{}

func (m checkJobCaseInsensitivePlanModifier) Description(ctx context.Context) string {
	return "Keeps the prior value when the only difference is letter case."
}

func (m checkJobCaseInsensitivePlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m checkJobCaseInsensitivePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() || req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}
	if strings.EqualFold(req.StateValue.ValueString(), req.PlanValue.ValueString()) {
		resp.PlanValue = req.StateValue
	}
}

// checkJobFinishedStatusPlanModifier treats a completed one-shot check as
// already satisfying a running or suspended request: the job cannot be started
// or suspended again after it finishes.
type checkJobFinishedStatusPlanModifier struct{}

func (m checkJobFinishedStatusPlanModifier) Description(ctx context.Context) string {
	return "Keeps a finished check job unchanged when a running or suspended status is requested."
}

func (m checkJobFinishedStatusPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m checkJobFinishedStatusPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() || req.StateValue.ValueString() != "Finished" {
		return
	}
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}
	switch req.PlanValue.ValueString() {
	case "Migrating", "Suspending":
		resp.PlanValue = req.StateValue
	}
}

func (r *checkJobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan checkJobResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The ConfigureDtsJob request is assembled before the instance is
	// purchased: assembling it parses the user supplied reserve JSON, and a
	// local failure there must abort the create before CreateDtsInstance buys
	// an instance that Terraform would have no way to read or delete.
	configureRequest, err := buildCheckJobConfigureRequest(plan, r.Client().RegionId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: building ConfigureDtsJob request", checkJobTypeName),
			errs.WrapError(err).Error(),
		)
		return
	}

	instanceRequest := map[string]interface{}{
		"Type":              "CHECK",
		"PayType":           "PostPaid",
		"RegionId":          r.Client().RegionId,
		"SourceRegion":      r.Client().RegionId,
		"DestinationRegion": r.Client().RegionId,
		"InstanceClass":     "small",
	}
	if v := plan.SourceEndpointRegion.ValueString(); !plan.SourceEndpointRegion.IsNull() && !plan.SourceEndpointRegion.IsUnknown() && v != "" {
		instanceRequest["SourceRegion"] = v
	}
	if v := plan.DestinationEndpointRegion.ValueString(); !plan.DestinationEndpointRegion.IsNull() && !plan.DestinationEndpointRegion.IsUnknown() && v != "" {
		instanceRequest["DestinationRegion"] = v
	}
	if v := plan.SourceEndpointEngineName.ValueString(); v != "" {
		instanceRequest["SourceEndpointEngineName"] = v
	}
	if v := plan.DestinationEndpointEngineName.ValueString(); v != "" {
		instanceRequest["DestinationEndpointEngineName"] = v
	}
	if v := plan.InstanceClass.ValueString(); !plan.InstanceClass.IsNull() && !plan.InstanceClass.IsUnknown() && v != "" {
		instanceRequest["InstanceClass"] = v
	}

	response, err := r.post("CreateDtsInstance", instanceRequest, checkJobCreateTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateDtsInstance", checkJobTypeName),
			err.Error(),
		)
		return
	}
	if fmt.Sprint(response["Success"]) == "false" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateDtsInstance", checkJobTypeName),
			fmt.Errorf("CreateDtsInstance failed, response: %v", response).Error(),
		)
		return
	}
	dtsInstanceId, _ := response["InstanceId"].(string)
	if dtsInstanceId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateDtsInstance", checkJobTypeName),
			errs.WrapError(fmt.Errorf("CreateDtsInstance did not return an instance id, response: %v", response)).Error(),
		)
		return
	}
	// The purchased instance is addressable from this point on. ConfigureDtsJob
	// is the next call and only it can return the job id, so a failure in
	// between would otherwise leave a running instance with no id in state.
	// Keeping the instance id as the resource id in this window lets Read and
	// Delete find it; it is replaced by the job id on success.
	plan.Id = types.StringValue(dtsInstanceId)
	plan.DtsInstanceId = types.StringValue(dtsInstanceId)
	resp.Diagnostics.Append(resp.State.Set(ctx, checkJobResolveUnknown(plan))...)
	if resp.Diagnostics.HasError() {
		return
	}

	configureRequest["DtsInstanceId"] = dtsInstanceId
	response, err = r.post("ConfigureDtsJob", configureRequest, checkJobCreateTimeout, "SQLExecuteError", "Forbidden.InstanceNotFound", "InvalidInstanceStatus")
	if err != nil {
		r.releasePurchasedInstance(ctx, resp, dtsInstanceId, err)
		return
	}
	if fmt.Sprint(response["Success"]) == "false" {
		r.releasePurchasedInstance(ctx, resp, dtsInstanceId, fmt.Errorf("ConfigureDtsJob failed, response: %v", response))
		return
	}
	dtsJobId, _ := response["DtsJobId"].(string)
	if dtsJobId == "" {
		r.releasePurchasedInstance(ctx, resp, dtsInstanceId, fmt.Errorf("ConfigureDtsJob did not return a job id, response: %v", response))
		return
	}
	plan.Id = types.StringValue(dtsJobId)
	resp.Diagnostics.Append(resp.State.Set(ctx, checkJobResolveUnknown(plan))...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateConf := helper.BuildStateConf([]string{}, []string{"NotStarted", "NotConfigured", "Prechecking", "PreCheckPass", "Migrating", "Finished"}, checkJobCreateTimeout, 5*time.Second, NewService(r.Client()).CheckJobStateRefreshFunc(dtsJobId, []string{"PrecheckFailed", "InitializeFailed", "Failed", "MigrationFailed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for the job to be configured", checkJobTypeName),
			errs.WrapErrorf(err, errs.IdMsg, dtsJobId).Error(),
		)
		return
	}

	// A configured status is an explicit intent for the new job; apply it here
	// instead of handing the freshly created resource to the update path.
	if target, ok := checkJobStringValue(plan.Status); ok {
		if err := r.statusFlow(dtsJobId, target, checkJobUpdateTimeout); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: applying status %q", checkJobTypeName, target),
				errs.WrapError(err).Error(),
			)
			return
		}
	}

	object, err := r.describe(dtsJobId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading the job after configuration", checkJobTypeName),
			errs.WrapError(err).Error(),
		)
		return
	}
	model := checkJobModelFromAPI(object, plan)
	model.Id = types.StringValue(dtsJobId)
	// A configured status is kept as planned even though the one-shot check can
	// settle on Finished before the read; Terraform rejects a known planned
	// value that changes during apply, and the next refresh reports the
	// server-side value.
	if !plan.Status.IsNull() && !plan.Status.IsUnknown() {
		model.Status = plan.Status
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, checkJobResolveUnknown(model))...)
}

// releasePurchasedInstance releases the CHECK instance purchased by
// CreateDtsInstance when the following configuration failed. DeleteDtsJob
// accepts DtsInstanceId alone, which is also how the standalone dts instance
// resource releases an unconfigured instance.
//
// When the release itself fails the instance id stays as the resource id and
// both failures are reported: the state keeps a handle that Read can keep alive
// and Delete can act on. Clearing the id in that case would leave a billable
// instance that Terraform no longer knows about.
func (r *checkJobResource) releasePurchasedInstance(ctx context.Context, resp *resource.CreateResponse, dtsInstanceId string, cause error) {
	releaseErr := r.deleteByInstanceId(dtsInstanceId, checkJobCreateTimeout)
	if releaseErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: releasing the purchased instance failed", checkJobTypeName),
			errs.WrapError(fmt.Errorf("%v; releasing the purchased DTS instance %s failed: %v", cause, dtsInstanceId, releaseErr)).Error(),
		)
		return
	}
	resp.State.RemoveResource(ctx)
	resp.Diagnostics.AddError(
		fmt.Sprintf("Creating %s: configuring the purchased instance failed", checkJobTypeName),
		errs.WrapError(cause).Error(),
	)
}

func (r *checkJobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state checkJobResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	object, err := r.describe(id)
	if err != nil {
		if errs.NotFoundError(err) {
			// A create that failed after CreateDtsInstance could not return a
			// job id and kept the purchased instance id as the resource id.
			// Such an instance has no job to describe, so a missing job here
			// does not mean the instance is gone; dropping the id would lose
			// the only handle Delete has on it.
			if instanceId := state.DtsInstanceId.ValueString(); id != "" && instanceId != "" && id == instanceId {
				return
			}
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeDtsJobDetail", checkJobTypeName),
			errs.WrapError(err).Error(),
		)
		return
	}

	model := checkJobModelFromAPI(object, state)
	resp.Diagnostics.Append(resp.State.Set(ctx, checkJobResolveUnknown(model))...)
}

func (r *checkJobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state checkJobResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	jobId := state.Id.ValueString()
	if !plan.DtsJobName.Equal(state.DtsJobName) {
		request := map[string]interface{}{
			"DtsJobId":   jobId,
			"DtsJobName": plan.DtsJobName.ValueString(),
			"RegionId":   r.Client().RegionId,
		}
		response, err := r.post("ModifyDtsJobName", request, checkJobUpdateTimeout)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling ModifyDtsJobName", checkJobTypeName),
				err.Error(),
			)
			return
		}
		if fmt.Sprint(response["Success"]) == "false" {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling ModifyDtsJobName", checkJobTypeName),
				fmt.Errorf("ModifyDtsJobName failed, response: %v", response).Error(),
			)
			return
		}
	}

	// A one-shot check starts by itself right after ConfigureDtsJob and
	// ModifyDtsJobPassword answers OperationDenied.JobStatus while it runs, so
	// the call is retried until the check settles.
	if !plan.SourceEndpointPassword.Equal(state.SourceEndpointPassword) {
		if err := r.updatePassword(jobId, "src", plan); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling ModifyDtsJobPassword for the source endpoint", checkJobTypeName),
				err.Error(),
			)
			return
		}
		if target, ok := checkJobStringValue(plan.Status); ok {
			if err := r.statusFlow(jobId, target, checkJobUpdateTimeout); err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Updating %s: waiting for status %q", checkJobTypeName, target),
					errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, target)).Error(),
				)
				return
			}
		}
	}
	if !plan.DestinationEndpointPassword.Equal(state.DestinationEndpointPassword) {
		if err := r.updatePassword(jobId, "dest", plan); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling ModifyDtsJobPassword for the destination endpoint", checkJobTypeName),
				err.Error(),
			)
			return
		}
		if target, ok := checkJobStringValue(plan.Status); ok {
			if err := r.statusFlow(jobId, target, checkJobUpdateTimeout); err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Updating %s: waiting for status %q", checkJobTypeName, target),
					errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, target)).Error(),
				)
				return
			}
		}
	}

	if !plan.Status.Equal(state.Status) {
		if target, ok := checkJobStringValue(plan.Status); ok {
			if err := r.statusFlow(jobId, target, checkJobUpdateTimeout); err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Updating %s: applying status %q", checkJobTypeName, target),
					errs.WrapError(err).Error(),
				)
				return
			}
		}
	}

	object, err := r.describe(jobId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after update", checkJobTypeName),
			errs.WrapError(err).Error(),
		)
		return
	}
	model := checkJobModelFromAPI(object, plan)
	model.Id = types.StringValue(jobId)
	// A check job is a one-shot task: rotating an endpoint password retries
	// ModifyDtsJobPassword until the running check settles, so the job can
	// reach Finished between plan and apply. Terraform requires a known
	// planned value to be returned unchanged, so the planned status stays in
	// state and the next refresh reports the server-side value.
	if !plan.Status.IsNull() && !plan.Status.IsUnknown() {
		model.Status = plan.Status
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, checkJobResolveUnknown(model))...)
}

func (r *checkJobResource) updatePassword(jobId, endpoint string, plan checkJobResourceModel) error {
	request := map[string]interface{}{
		"DtsJobId": jobId,
		"Endpoint": endpoint,
		"RegionId": r.Client().RegionId,
	}
	if endpoint == "src" {
		if v, ok := checkJobStringValue(plan.SourceEndpointPassword); ok {
			request["Password"] = v
		}
		if v, ok := checkJobStringValue(plan.SourceEndpointUserName); ok {
			request["UserName"] = v
		}
	} else {
		if v, ok := checkJobStringValue(plan.DestinationEndpointPassword); ok {
			request["Password"] = v
		}
		if v, ok := checkJobStringValue(plan.DestinationEndpointUserName); ok {
			request["UserName"] = v
		}
	}
	response, err := r.post("ModifyDtsJobPassword", request, checkJobUpdateTimeout, "OperationDenied.JobStatus")
	if err != nil {
		return err
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return fmt.Errorf("ModifyDtsJobPassword failed, response: %v", response)
	}
	return nil
}

func (r *checkJobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state checkJobResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"RegionId": r.Client().RegionId,
	}
	if v := state.DtsInstanceId.ValueString(); v != "" {
		request["DtsInstanceId"] = v
	}
	// A create that failed before ConfigureDtsJob returned a job id keeps the
	// purchased instance id as the resource id; such a resource has no job and
	// is released by instance id alone.
	if id := state.Id.ValueString(); id != "" && id != state.DtsInstanceId.ValueString() {
		request["DtsJobId"] = id
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(checkJobDeleteTimeout, func() *retry.RetryError {
		response, err = r.Client().RpcPost(checkJobProduct, checkJobAPIVersion, "DeleteDtsJob", nil, request, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Forbidden.InstanceNotFound", "InvalidJobId"}) {
				return nil
			}
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug("DeleteDtsJob", response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteDtsJob", checkJobTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), "DeleteDtsJob", errs.AlibabaCloudSdkGoERROR).Error(),
		)
	}
}

func (r *checkJobResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// deleteByInstanceId releases an unconfigured CHECK instance. The documented
// not-found codes mean the instance is already gone.
func (r *checkJobResource) deleteByInstanceId(dtsInstanceId string, timeout time.Duration) error {
	request := map[string]interface{}{
		"DtsInstanceId": dtsInstanceId,
		"RegionId":      r.Client().RegionId,
	}
	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(timeout, func() *retry.RetryError {
		response, err = r.Client().RpcPost(checkJobProduct, checkJobAPIVersion, "DeleteDtsJob", nil, request, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Forbidden.InstanceNotFound", "InvalidJobId"}) {
				return nil
			}
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug("DeleteDtsJob", response, request)
	if err != nil {
		return err
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return fmt.Errorf("DeleteDtsJob failed, response: %v", response)
	}
	return nil
}

// statusFlow brings the job to the requested running state. A check job runs
// once and stays Finished afterwards, so a finished job already satisfies a
// requested running state.
func (r *checkJobResource) statusFlow(jobId, target string, timeout time.Duration) error {
	if target == "" {
		return nil
	}
	object, err := r.describe(jobId)
	if err != nil {
		return errs.WrapError(err)
	}
	current := fmt.Sprint(object["Status"])
	if current == target {
		return nil
	}
	if current == "Finished" && (target == "Migrating" || target == "Suspending") {
		return nil
	}
	switch target {
	case "Migrating":
		request := map[string]interface{}{
			"DtsJobId": jobId,
			"RegionId": r.Client().RegionId,
		}
		response, err := r.post("StartDtsJob", request, timeout)
		if err != nil {
			return err
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return fmt.Errorf("StartDtsJob failed, response: %v", response)
		}
		// Wait for a settled running state rather than the first observable
		// one: returning while the job is still Prechecking/PreCheckPass
		// leaves a server-driven status in the refreshed state that the
		// post-apply plan reports as drift.
		stateConf := helper.BuildStateConf([]string{}, []string{"Migrating", "Finished"}, timeout, 60*time.Second, NewService(r.Client()).CheckJobStateRefreshFunc(jobId, []string{"PrecheckFailed", "InitializeFailed", "Failed", "MigrationFailed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, jobId)
		}
	case "Suspending":
		request := map[string]interface{}{
			"DtsJobId": jobId,
			"RegionId": r.Client().RegionId,
		}
		response, err := r.post("SuspendDtsJob", request, timeout)
		if err != nil {
			return err
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return fmt.Errorf("SuspendDtsJob failed, response: %v", response)
		}
		stateConf := helper.BuildStateConf([]string{}, []string{"Suspending"}, timeout, 5*time.Second, NewService(r.Client()).CheckJobStateRefreshFunc(jobId, []string{"PrecheckFailed", "InitializeFailed", "Failed", "MigrationFailed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, jobId)
		}
	}
	return nil
}

// post issues a DTS call with the standard retry envelope. The optional
// expectedErrors are transient for the call and are retried until timeout.
func (r *checkJobResource) post(action string, request map[string]interface{}, timeout time.Duration, expectedErrors ...string) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(timeout, func() *retry.RetryError {
		response, err = r.Client().RpcPost(checkJobProduct, checkJobAPIVersion, action, nil, request, false)
		if err != nil {
			if (len(expectedErrors) > 0 && errs.IsExpectedErrors(err, expectedErrors)) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, checkJobTypeName, action, errs.AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

// buildCheckJobConfigureRequest assembles the ConfigureDtsJob request without
// the instance id, which is only known after CreateDtsInstance.
func buildCheckJobConfigureRequest(plan checkJobResourceModel, regionId string) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"DtsJobName": plan.DtsJobName.ValueString(),
		"DbList":     plan.DbList.ValueString(),
		"JobType":    "CHECK",
		"RegionId":   regionId,
		// JobType CHECK performs only data verification, the three
		// initialization flags must be false as required by ConfigureDtsJob.
		"DataInitialization":      false,
		"DataSynchronization":     false,
		"StructureInitialization": false,
	}
	request["DestinationEndpointInstanceType"] = plan.DestinationEndpointInstanceType.ValueString()
	request["SourceEndpointInstanceType"] = plan.SourceEndpointInstanceType.ValueString()

	setCheckJobString(request, "Checkpoint", plan.Checkpoint)
	setCheckJobString(request, "DelayPhone", plan.DelayPhone)
	setCheckJobString(request, "DelayRuleTime", plan.DelayRuleTime)
	setCheckJobBool(request, "DelayNotice", plan.DelayNotice)
	setCheckJobBool(request, "ErrorNotice", plan.ErrorNotice)
	setCheckJobString(request, "ErrorPhone", plan.ErrorPhone)
	setCheckJobString(request, "DestinationEndpointDataBaseName", plan.DestinationEndpointDatabaseName)
	setCheckJobString(request, "DestinationEndpointEngineName", plan.DestinationEndpointEngineName)
	setCheckJobString(request, "DestinationEndpointIP", plan.DestinationEndpointIp)
	setCheckJobString(request, "DestinationEndpointInstanceID", plan.DestinationEndpointInstanceId)
	setCheckJobString(request, "DestinationEndpointOracleSID", plan.DestinationEndpointOracleSid)
	setCheckJobString(request, "DestinationEndpointPassword", plan.DestinationEndpointPassword)
	setCheckJobString(request, "DestinationEndpointPort", plan.DestinationEndpointPort)
	setCheckJobString(request, "DestinationEndpointRegion", plan.DestinationEndpointRegion)
	setCheckJobString(request, "DestinationEndpointUserName", plan.DestinationEndpointUserName)
	setCheckJobString(request, "DestinationEndpointOwnerID", plan.DestinationEndpointOwnerId)
	setCheckJobString(request, "DestinationEndpointRole", plan.DestinationEndpointRole)
	setCheckJobString(request, "Reserve", plan.Reserve)
	if err := setCheckJobEndpointSSL(request, plan); err != nil {
		return nil, err
	}
	setCheckJobString(request, "SourceEndpointDatabaseName", plan.SourceEndpointDatabaseName)
	setCheckJobString(request, "SourceEndpointEngineName", plan.SourceEndpointEngineName)
	setCheckJobString(request, "SourceEndpointIP", plan.SourceEndpointIp)
	setCheckJobString(request, "SourceEndpointInstanceID", plan.SourceEndpointInstanceId)
	setCheckJobString(request, "SourceEndpointOracleSID", plan.SourceEndpointOracleSid)
	setCheckJobString(request, "SourceEndpointOwnerID", plan.SourceEndpointOwnerId)
	setCheckJobString(request, "SourceEndpointPassword", plan.SourceEndpointPassword)
	setCheckJobString(request, "SourceEndpointPort", plan.SourceEndpointPort)
	setCheckJobString(request, "SourceEndpointRegion", plan.SourceEndpointRegion)
	setCheckJobString(request, "SourceEndpointRole", plan.SourceEndpointRole)
	setCheckJobString(request, "SourceEndpointUserName", plan.SourceEndpointUserName)
	setCheckJobString(request, "DedicatedClusterId", plan.DedicatedClusterId)
	setCheckJobString(request, "DataCheckConfigure", plan.DataCheckConfigure)
	setCheckJobString(request, "SourceEndpointVSwitchID", plan.SourceEndpointVswitchId)
	setCheckJobString(request, "SrcPrimaryVswId", plan.SrcPrimaryVswitchId)
	setCheckJobString(request, "SrcSecondaryVswId", plan.SrcSecondaryVswitchId)
	setCheckJobString(request, "DestPrimaryVswId", plan.DestPrimaryVswitchId)
	setCheckJobString(request, "DestSecondaryVswId", plan.DestSecondaryVswitchId)
	setCheckJobString(request, "DtsBisLabel", plan.DtsBisLabel)
	return request, nil
}

// setCheckJobEndpointSSL folds source_endpoint_ssl and destination_endpoint_ssl
// into the DTS Reserve parameter as its srcSSL and destSSL keys.
// ConfigureDtsJob has no dedicated SSL parameter, so Reserve is the only place
// the connection mode can be sent. Values already present under those keys in
// a user-supplied reserve are overwritten by the dedicated attributes.
func setCheckJobEndpointSSL(request map[string]interface{}, plan checkJobResourceModel) error {
	srcSSL, srcOk := checkJobStringValue(plan.SourceEndpointSsl)
	destSSL, destOk := checkJobStringValue(plan.DestinationEndpointSsl)
	if !srcOk && !destOk {
		return nil
	}
	reserve := make(map[string]interface{})
	if v, ok := request["Reserve"]; ok {
		var decoded interface{}
		if err := json.Unmarshal([]byte(fmt.Sprint(v)), &decoded); err != nil {
			return errs.WrapError(err)
		}
		// json.Unmarshal accepts null and non-object documents without an
		// error, but only a JSON object can carry srcSSL/destSSL, and a nil map
		// would panic on the writes below.
		object, ok := decoded.(map[string]interface{})
		if !ok {
			return errs.WrapError(fmt.Errorf("reserve must be a JSON object, got %T", decoded))
		}
		reserve = object
	}
	if srcOk {
		reserve["srcSSL"] = srcSSL
	}
	if destOk {
		reserve["destSSL"] = destSSL
	}
	reserveJson, err := json.Marshal(reserve)
	if err != nil {
		return errs.WrapError(err)
	}
	request["Reserve"] = string(reserveJson)
	return nil
}

func setCheckJobString(request map[string]interface{}, key string, value types.String) {
	if v, ok := checkJobStringValue(value); ok {
		request[key] = v
	}
}

func setCheckJobBool(request map[string]interface{}, key string, value types.Bool) {
	if value.IsNull() || value.IsUnknown() {
		return
	}
	request[key] = value.ValueBool()
}

func checkJobStringValue(value types.String) (string, bool) {
	if value.IsNull() || value.IsUnknown() || value.ValueString() == "" {
		return "", false
	}
	return value.ValueString(), true
}

// checkJobModelFromAPI overlays the DescribeDtsJobDetail response onto base,
// which is the plan on create/update and the prior state on read. Attributes
// the API does not return (write-only settings, passwords, import-only
// switches) keep the base value, and endpoints are only overwritten when the
// response actually carries a value.
func checkJobModelFromAPI(object map[string]interface{}, base checkJobResourceModel) checkJobResourceModel {
	model := base
	if v, ok := object["DtsInstanceID"].(string); ok && v != "" {
		model.DtsInstanceId = types.StringValue(v)
	}
	if v, ok := object["DtsJobName"].(string); ok && v != "" {
		model.DtsJobName = types.StringValue(v)
	}
	if v, err := convert.FormatInt(object["Checkpoint"]); err == nil {
		model.Checkpoint = types.StringValue(fmt.Sprint(v))
	}
	if v, ok := object["DbObject"].(string); ok && v != "" {
		model.DbList = types.StringValue(v)
	}
	// DtsJobClass is the link specification reported back by
	// DescribeDtsJobDetail and the only API field carrying instance_class;
	// CHECK jobs use the same class names as other DTS jobs.
	if v, ok := object["DtsJobClass"].(string); ok && v != "" {
		model.InstanceClass = types.StringValue(v)
	}
	if v, ok := object["Status"]; ok && v != nil && fmt.Sprint(v) != "" {
		model.Status = types.StringValue(fmt.Sprint(convertCheckJobStatusResponse(v)))
	}

	source := checkJobObject(object["SourceEndpoint"])
	destination := checkJobObject(object["DestinationEndpoint"])
	model.SourceEndpointDatabaseName = checkJobOverlayString(source["DatabaseName"], model.SourceEndpointDatabaseName)
	model.SourceEndpointEngineName = checkJobOverlayCaseInsensitiveString(convertSourceEndpointEngineNameUppercaseResponse(source["EngineName"]), model.SourceEndpointEngineName)
	model.SourceEndpointIp = checkJobOverlayString(source["Ip"], model.SourceEndpointIp)
	model.SourceEndpointInstanceId = checkJobOverlayString(source["InstanceID"], model.SourceEndpointInstanceId)
	model.SourceEndpointInstanceType = checkJobOverlayCaseInsensitiveString(source["InstanceType"], model.SourceEndpointInstanceType)
	model.SourceEndpointOracleSid = checkJobOverlayString(source["OracleSID"], model.SourceEndpointOracleSid)
	model.SourceEndpointOwnerId = checkJobOverlayString(source["AliyunUid"], model.SourceEndpointOwnerId)
	model.SourceEndpointPort = checkJobOverlayString(source["Port"], model.SourceEndpointPort)
	model.SourceEndpointRegion = checkJobOverlayString(source["Region"], model.SourceEndpointRegion)
	model.SourceEndpointRole = checkJobOverlayString(source["RoleName"], model.SourceEndpointRole)
	model.SourceEndpointUserName = checkJobOverlayString(source["UserName"], model.SourceEndpointUserName)
	if ssl := convertEndpointSslResponse(source["SslSolutionEnum"]); ssl != nil {
		model.SourceEndpointSsl = types.StringValue(fmt.Sprint(ssl))
	}

	model.DestinationEndpointDatabaseName = checkJobOverlayString(destination["DatabaseName"], model.DestinationEndpointDatabaseName)
	model.DestinationEndpointEngineName = checkJobOverlayCaseInsensitiveString(destination["EngineName"], model.DestinationEndpointEngineName)
	model.DestinationEndpointIp = checkJobOverlayString(destination["Ip"], model.DestinationEndpointIp)
	model.DestinationEndpointInstanceId = checkJobOverlayString(destination["InstanceID"], model.DestinationEndpointInstanceId)
	model.DestinationEndpointInstanceType = checkJobOverlayCaseInsensitiveString(destination["InstanceType"], model.DestinationEndpointInstanceType)
	model.DestinationEndpointOracleSid = checkJobOverlayString(destination["OracleSID"], model.DestinationEndpointOracleSid)
	// The destination account and role identify the cross-account endpoint.
	// DescribeDtsJobDetail reports them as DestinationEndpoint.AliyunUid and
	// RoleName; without these the ForceNew values are lost on import and
	// restoring the configuration plans an unnecessary replacement.
	model.DestinationEndpointOwnerId = checkJobOverlayString(destination["AliyunUid"], model.DestinationEndpointOwnerId)
	model.DestinationEndpointRole = checkJobOverlayString(destination["RoleName"], model.DestinationEndpointRole)
	model.DestinationEndpointPort = checkJobOverlayString(destination["Port"], model.DestinationEndpointPort)
	model.DestinationEndpointRegion = checkJobOverlayString(destination["Region"], model.DestinationEndpointRegion)
	model.DestinationEndpointUserName = checkJobOverlayString(destination["UserName"], model.DestinationEndpointUserName)
	if ssl := convertEndpointSslResponse(destination["SslSolutionEnum"]); ssl != nil {
		model.DestinationEndpointSsl = types.StringValue(fmt.Sprint(ssl))
	}

	migrationMode := checkJobObject(object["MigrationMode"])
	// DataInitialization/DataSynchronization/StructureInitialization are fixed
	// to false for CHECK jobs; keep them in sync with the API response so the
	// state stays consistent.
	model.DataInitialization = checkJobComputedBool(migrationMode["DataInitialization"], model.DataInitialization)
	model.DataSynchronization = checkJobComputedBool(migrationMode["DataSynchronization"], model.DataSynchronization)
	model.StructureInitialization = checkJobComputedBool(migrationMode["StructureInitialization"], model.StructureInitialization)

	return model
}

func checkJobObject(value interface{}) map[string]interface{} {
	if object, ok := value.(map[string]interface{}); ok {
		return object
	}
	return map[string]interface{}{}
}

func checkJobOverlayString(value interface{}, base types.String) types.String {
	if v := fmt.Sprint(value); value != nil && v != "" && v != "<nil>" {
		return types.StringValue(v)
	}
	return base
}

// checkJobOverlayCaseInsensitiveString keeps the base spelling when the API
// value differs only by case, so a configured lower-case value stays stable.
func checkJobOverlayCaseInsensitiveString(value interface{}, base types.String) types.String {
	candidate := checkJobOverlayString(value, base)
	if base.IsNull() || base.IsUnknown() || !strings.EqualFold(candidate.ValueString(), base.ValueString()) {
		return candidate
	}
	return base
}

// checkJobComputedBool keeps a known bool and falls back to false, the only value the API accepts
// for the three initialization flags of a CHECK job, so a computed attribute is never left unknown.
func checkJobComputedBool(value interface{}, base types.Bool) types.Bool {
	if overlay := checkJobOverlayBool(value, base); !overlay.IsNull() && !overlay.IsUnknown() {
		return overlay
	}
	return types.BoolValue(false)
}

func checkJobOverlayBool(value interface{}, base types.Bool) types.Bool {
	switch v := value.(type) {
	case bool:
		return types.BoolValue(v)
	case string:
		if parsed, err := strconv.ParseBool(v); err == nil {
			return types.BoolValue(parsed)
		}
	}
	return base
}

// checkJobResolveUnknown rewrites unknown plan values to null before they are
// written to state: Terraform cannot store an unknown value, and every
// attribute that is still unknown at this point has no server counterpart.
func checkJobResolveUnknown(model checkJobResourceModel) checkJobResourceModel {
	model.Id = checkJobKnownString(model.Id)
	model.DtsInstanceId = checkJobKnownString(model.DtsInstanceId)
	model.DtsJobName = checkJobKnownString(model.DtsJobName)
	model.DedicatedClusterId = checkJobKnownString(model.DedicatedClusterId)
	model.DataCheckConfigure = checkJobKnownString(model.DataCheckConfigure)
	model.Checkpoint = checkJobKnownString(model.Checkpoint)
	model.InstanceClass = checkJobKnownString(model.InstanceClass)
	model.DbList = checkJobKnownString(model.DbList)
	model.Reserve = checkJobKnownString(model.Reserve)
	model.SourceEndpointInstanceType = checkJobKnownString(model.SourceEndpointInstanceType)
	model.SourceEndpointEngineName = checkJobKnownString(model.SourceEndpointEngineName)
	model.SourceEndpointInstanceId = checkJobKnownString(model.SourceEndpointInstanceId)
	model.SourceEndpointRegion = checkJobKnownString(model.SourceEndpointRegion)
	model.SourceEndpointIp = checkJobKnownString(model.SourceEndpointIp)
	model.SourceEndpointPort = checkJobKnownString(model.SourceEndpointPort)
	model.SourceEndpointOracleSid = checkJobKnownString(model.SourceEndpointOracleSid)
	model.SourceEndpointDatabaseName = checkJobKnownString(model.SourceEndpointDatabaseName)
	model.SourceEndpointUserName = checkJobKnownString(model.SourceEndpointUserName)
	model.SourceEndpointPassword = checkJobKnownString(model.SourceEndpointPassword)
	model.SourceEndpointOwnerId = checkJobKnownString(model.SourceEndpointOwnerId)
	model.SourceEndpointRole = checkJobKnownString(model.SourceEndpointRole)
	model.SourceEndpointVswitchId = checkJobKnownString(model.SourceEndpointVswitchId)
	model.SrcPrimaryVswitchId = checkJobKnownString(model.SrcPrimaryVswitchId)
	model.SrcSecondaryVswitchId = checkJobKnownString(model.SrcSecondaryVswitchId)
	model.SourceEndpointSsl = checkJobKnownString(model.SourceEndpointSsl)
	model.DestinationEndpointInstanceType = checkJobKnownString(model.DestinationEndpointInstanceType)
	model.DestinationEndpointEngineName = checkJobKnownString(model.DestinationEndpointEngineName)
	model.DestinationEndpointInstanceId = checkJobKnownString(model.DestinationEndpointInstanceId)
	model.DestinationEndpointRegion = checkJobKnownString(model.DestinationEndpointRegion)
	model.DestinationEndpointIp = checkJobKnownString(model.DestinationEndpointIp)
	model.DestinationEndpointPort = checkJobKnownString(model.DestinationEndpointPort)
	model.DestinationEndpointDatabaseName = checkJobKnownString(model.DestinationEndpointDatabaseName)
	model.DestinationEndpointUserName = checkJobKnownString(model.DestinationEndpointUserName)
	model.DestinationEndpointPassword = checkJobKnownString(model.DestinationEndpointPassword)
	model.DestinationEndpointOracleSid = checkJobKnownString(model.DestinationEndpointOracleSid)
	model.DestinationEndpointOwnerId = checkJobKnownString(model.DestinationEndpointOwnerId)
	model.DestinationEndpointRole = checkJobKnownString(model.DestinationEndpointRole)
	model.DestinationEndpointSsl = checkJobKnownString(model.DestinationEndpointSsl)
	model.DestPrimaryVswitchId = checkJobKnownString(model.DestPrimaryVswitchId)
	model.DestSecondaryVswitchId = checkJobKnownString(model.DestSecondaryVswitchId)
	model.DtsBisLabel = checkJobKnownString(model.DtsBisLabel)
	model.DelayNotice = checkJobKnownBool(model.DelayNotice)
	model.DelayPhone = checkJobKnownString(model.DelayPhone)
	model.DelayRuleTime = checkJobKnownString(model.DelayRuleTime)
	model.ErrorNotice = checkJobKnownBool(model.ErrorNotice)
	model.ErrorPhone = checkJobKnownString(model.ErrorPhone)
	model.Status = checkJobKnownString(model.Status)
	model.DataInitialization = checkJobKnownBool(model.DataInitialization)
	model.DataSynchronization = checkJobKnownBool(model.DataSynchronization)
	model.StructureInitialization = checkJobKnownBool(model.StructureInitialization)
	return model
}

func checkJobKnownString(value types.String) types.String {
	if value.IsUnknown() {
		return types.StringNull()
	}
	return value
}

func checkJobKnownBool(value types.Bool) types.Bool {
	if value.IsUnknown() {
		return types.BoolNull()
	}
	return value
}
