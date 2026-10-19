package cms

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const alertRobotTypeName = "alicloud_cms_alert_robot"

const (
	// /robot is the RESTful (ROA) resource root for a single AlertRobot.
	alertRobotResourcePath = "/robot"
	// /robots is the collection path used by List (Read) and Delete (bulk by ids).
	alertRobotsCollectionPath = "/robots"

	cmsProduct  = "Cms"
	cmsVersion  = "2024-03-30"
	langDefault = "zh_CN"
)

var (
	_ resource.Resource                = &alertRobotResource{}
	_ resource.ResourceWithImportState = &alertRobotResource{}
)

func NewAlertRobotResource() resource.Resource {
	return &alertRobotResource{}
}

type alertRobotResource struct {
	fwadapt.ResourceBase
}

type alertRobotResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	AlertRobotId        types.String `tfsdk:"alert_robot_id"`
	AlertRobotName      types.String `tfsdk:"alert_robot_name"`
	Type                types.String `tfsdk:"type"`
	Lang                types.String `tfsdk:"lang"`
	Url                 types.String `tfsdk:"url"`
	Workspace           types.String `tfsdk:"workspace"`
	DigitalEmployeeName types.String `tfsdk:"digital_employee_name"`
	RobotSignKey        types.String `tfsdk:"robot_sign_key"`
}

func (r *alertRobotResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = alertRobotTypeName
}

func (r *alertRobotResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a Cloud Monitor (CMS) Alert Robot resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the alert robot. Same as `alert_robot_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"alert_robot_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the alert robot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"alert_robot_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the alert robot.",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the webhook. Valid values: `DING`, `WEIXIN`, `FEISHU`, `SLACK`, `TEAMS`, `DING_COOL_APP`.",
				Validators: []validator.String{
					stringvalidator.OneOf("DING", "WEIXIN", "FEISHU", "SLACK", "TEAMS", "DING_COOL_APP"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"lang": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The language of the alert robot. Valid values: `zh_CN`, `en_US`. Defaults to `zh_CN` when omitted.",
				Validators: []validator.String{
					stringvalidator.OneOf("zh_CN", "en_US"),
				},
			},
			"url": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The webhook URL of the alert robot. Write-only at the API level: the List API never returns it, so the configured value is preserved in state and ignored on import verification.",
			},
			"workspace": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The workspace ID of the alert robot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"digital_employee_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the associated digital employee.",
			},
			"robot_sign_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The sign key of the alert robot. Sensitive: the API returns it when set and omits it when absent.",
			},
		},
	}
}

func (r *alertRobotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan alertRobotResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.create(ctx, r.Client(), plan)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Creating %s", alertRobotTypeName), err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *alertRobotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state alertRobotResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updated, gone, err := r.read(ctx, r.Client(), state)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Reading %s", alertRobotTypeName), err.Error())
		return
	}
	if gone {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, updated)...)
}

func (r *alertRobotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state alertRobotResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updated, err := r.update(ctx, r.Client(), plan, state)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Updating %s", alertRobotTypeName), err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, updated)...)
}

func (r *alertRobotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state alertRobotResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.delete(ctx, r.Client(), state); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Deleting %s", alertRobotTypeName), err.Error())
		return
	}
}

func (r *alertRobotResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// create builds the Create request body, calls POST /robot, then re-reads the
// robot to populate the computed fields. The webhook URL is write-only at the
// API level, so it is carried from the plan into the resulting state.
func (r *alertRobotResource) create(ctx context.Context, client *connectivity.AliyunClient, plan alertRobotResourceModel) (alertRobotResourceModel, error) {
	body := map[string]interface{}{
		"type":    plan.Type.ValueString(),
		"name":    plan.AlertRobotName.ValueString(),
		"webhook": plan.Url.ValueString(),
		"lang":    resolveLang(plan.Lang),
	}
	if !plan.Workspace.IsNull() {
		body["workspace"] = plan.Workspace.ValueString()
	}
	if !plan.DigitalEmployeeName.IsNull() {
		body["digitalEmployeeName"] = plan.DigitalEmployeeName.ValueString()
	}
	if !plan.RobotSignKey.IsNull() {
		body["robotSignKey"] = plan.RobotSignKey.ValueString()
	}

	response, err := client.RoaPost(cmsProduct, cmsVersion, alertRobotResourcePath, map[string]*string{}, nil, body, true)
	if err != nil {
		return alertRobotResourceModel{}, errs.WrapErrorf(err, errs.DefaultErrorMsg, alertRobotTypeName, alertRobotResourcePath, errs.AlibabaCloudSdkGoERROR)
	}

	id, _ := response["alertRobotId"].(string)
	if id == "" {
		return alertRobotResourceModel{}, fmt.Errorf("creating %s: the response carried no alertRobotId: %v", alertRobotTypeName, response)
	}

	robot, err := describeAlertRobot(client, id)
	if err != nil {
		return alertRobotResourceModel{}, errs.WrapErrorf(err, errs.DefaultErrorMsg, alertRobotTypeName, alertRobotResourcePath, errs.AlibabaCloudSdkGoERROR)
	}
	if robot == nil {
		return alertRobotResourceModel{}, fmt.Errorf("creating %s: the robot %s is not visible via %s right after creation", alertRobotTypeName, id, alertRobotsCollectionPath)
	}
	return mergePlan(alertRobotState(robot), plan), nil
}

func (r *alertRobotResource) read(ctx context.Context, client *connectivity.AliyunClient, state alertRobotResourceModel) (alertRobotResourceModel, bool, error) {
	robot, err := describeAlertRobot(client, state.Id.ValueString())
	if err != nil {
		return alertRobotResourceModel{}, false, errs.WrapErrorf(err, errs.DefaultErrorMsg, alertRobotTypeName, alertRobotsCollectionPath, errs.AlibabaCloudSdkGoERROR)
	}
	if robot == nil {
		return alertRobotResourceModel{}, true, nil
	}
	return mergePlan(alertRobotState(robot), state), false, nil
}

func (r *alertRobotResource) update(ctx context.Context, client *connectivity.AliyunClient, plan, state alertRobotResourceModel) (alertRobotResourceModel, error) {
	action := fmt.Sprintf("/robot/%s", state.Id.ValueString())
	body := map[string]interface{}{}

	if !plan.AlertRobotName.Equal(state.AlertRobotName) {
		body["name"] = plan.AlertRobotName.ValueString()
	}
	if !plan.Lang.Equal(state.Lang) {
		body["lang"] = resolveLang(plan.Lang)
	}
	if !plan.Url.Equal(state.Url) {
		body["webhook"] = plan.Url.ValueString()
	}
	if !plan.Workspace.Equal(state.Workspace) {
		body["workspace"] = plan.Workspace.ValueString()
	}
	if !plan.DigitalEmployeeName.Equal(state.DigitalEmployeeName) {
		body["digitalEmployeeName"] = plan.DigitalEmployeeName.ValueString()
	}
	if !plan.RobotSignKey.Equal(state.RobotSignKey) {
		body["robotSignKey"] = plan.RobotSignKey.ValueString()
	}

	if len(body) > 0 {
		if _, err := client.RoaPatch(cmsProduct, cmsVersion, action, map[string]*string{}, nil, body, true); err != nil {
			return alertRobotResourceModel{}, errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	robot, err := describeAlertRobot(client, state.Id.ValueString())
	if err != nil {
		return alertRobotResourceModel{}, errs.WrapErrorf(err, errs.DefaultErrorMsg, alertRobotTypeName, alertRobotsCollectionPath, errs.AlibabaCloudSdkGoERROR)
	}
	if robot == nil {
		// Still exists for the API's purposes of an in-place update, but the
		// robot vanished mid-update; surface the inconsistency rather than
		// silently dropping state.
		return alertRobotResourceModel{}, fmt.Errorf("updating %s: the robot %s is not visible via %s after the patch", alertRobotTypeName, state.Id.ValueString(), alertRobotsCollectionPath)
	}
	return mergePlan(alertRobotState(robot), plan), nil
}

func (r *alertRobotResource) delete(ctx context.Context, client *connectivity.AliyunClient, state alertRobotResourceModel) error {
	query := map[string]*string{}
	robotIds, _ := json.Marshal([]string{state.Id.ValueString()})
	query["robotIds"] = stringPointer(string(robotIds))
	if !state.Type.IsNull() {
		query["type"] = stringPointer(state.Type.ValueString())
	}

	_, err := client.RoaDelete(cmsProduct, cmsVersion, alertRobotsCollectionPath, query, nil, nil, true)
	if err != nil {
		// Idempotent delete: a missing robot (404 / ResourceNotFound / InvalidParameter
		// for an unknown id) is success, not a failure mid-destroy.
		if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"ResourceNotFound", "InvalidParameter"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), alertRobotsCollectionPath, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

// describeAlertRobot lists robots filtered by a single id and returns the matching
// object, or nil when the robot no longer exists. The List API is the only Read
// surface (there is no single-resource GET).
func describeAlertRobot(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	query := map[string]*string{}
	robotIds, _ := json.Marshal([]string{id})
	query["robotIds"] = stringPointer(string(robotIds))
	query["pageNumber"] = stringPointer("1")
	query["pageSize"] = stringPointer("100")

	response, err := client.RoaGet(cmsProduct, cmsVersion, alertRobotsCollectionPath, query, nil, nil)
	if err != nil {
		return nil, err
	}

	robots, ok := response["robots"].([]interface{})
	if !ok {
		return nil, nil
	}
	for _, raw := range robots {
		robot, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if robotId, _ := robot["robotId"].(string); robotId == id {
			return robot, nil
		}
	}
	return nil, nil
}

// alertRobotState maps a List API robot object onto the resource model. Empty
// strings become null so omitted optional fields read back as absent.
func alertRobotState(robot map[string]interface{}) alertRobotResourceModel {
	str := func(key string) types.String {
		if v, ok := robot[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	return alertRobotResourceModel{
		Id:                  str("robotId"),
		AlertRobotId:        str("robotId"),
		AlertRobotName:      str("name"),
		Type:                str("type"),
		Lang:                str("lang"),
		Workspace:           str("workspace"),
		DigitalEmployeeName: str("digitalEmployeeName"),
		RobotSignKey:        str("robotSignKey"),
		// url (webhook) is write-only at the API level: List never returns it,
		// so it is filled from plan/state by mergePlan.
		Url: types.StringNull(),
	}
}

// mergePlan restores the write-only and sensitive-as-planned fields the List API
// does not echo back, so the persisted state matches the configured plan.
func mergePlan(state, plan alertRobotResourceModel) alertRobotResourceModel {
	if state.Url.IsNull() && !plan.Url.IsNull() {
		state.Url = plan.Url
	}
	if state.RobotSignKey.IsNull() && !plan.RobotSignKey.IsNull() {
		state.RobotSignKey = plan.RobotSignKey
	}
	return state
}

// resolveLang returns the effective lang value, defaulting to zh_CN when the plan
// omits it (lang is Optional+Computed).
func resolveLang(v types.String) string {
	if v.IsNull() || v.ValueString() == "" {
		return langDefault
	}
	return v.ValueString()
}
