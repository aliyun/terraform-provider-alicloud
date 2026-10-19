package cms

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const alertRobotsDataSourceTypeName = "alicloud_cms_alert_robots"

var (
	_ datasource.DataSource              = &alertRobotsDataSource{}
	_ datasource.DataSourceWithConfigure = &alertRobotsDataSource{}
)

func NewAlertRobotsDataSource() datasource.DataSource {
	return &alertRobotsDataSource{}
}

type alertRobotsDataSource struct {
	fwadapt.DataSourceBase
}

type alertRobotsDataSourceModel struct {
	Id         types.String `tfsdk:"id"`
	Ids        types.List   `tfsdk:"ids"`
	NameRegex  types.String `tfsdk:"name_regex"`
	Type       types.String `tfsdk:"type"`
	Workspace  types.String `tfsdk:"workspace"`
	OutputFile types.String `tfsdk:"output_file"`
	Robots     types.List   `tfsdk:"robots"`
}

type alertRobotDataSourceRobotModel struct {
	AlertRobotId        types.String `tfsdk:"alert_robot_id"`
	AlertRobotName      types.String `tfsdk:"alert_robot_name"`
	Type                types.String `tfsdk:"type"`
	Lang                types.String `tfsdk:"lang"`
	Url                 types.String `tfsdk:"url"`
	Workspace           types.String `tfsdk:"workspace"`
	DigitalEmployeeName types.String `tfsdk:"digital_employee_name"`
	RobotSignKey        types.String `tfsdk:"robot_sign_key"`
}

func (d *alertRobotsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = alertRobotsDataSourceTypeName
}

func (d *alertRobotsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a list of Cloud Monitor (CMS) Alert Robots in an Alibaba Cloud account according to the specified filters.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, a stable hash of the matched robot IDs.",
			},
			"ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "A list of alert robot IDs to filter. Also exported as the IDs of the matched robots.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string to filter alert robots by name.",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "The type of the webhook to filter. Valid values: `DING`, `WEIXIN`, `FEISHU`, `SLACK`, `TEAMS`, `DING_COOL_APP`.",
				Validators: []validator.String{
					stringvalidator.OneOf("DING", "WEIXIN", "FEISHU", "SLACK", "TEAMS", "DING_COOL_APP"),
				},
			},
			"workspace": schema.StringAttribute{
				Optional:    true,
				Description: "The workspace ID to filter.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "File path where to save the results.",
			},
			// robots is a list of objects modeled with an ObjectType rather than a
			// ListNestedAttribute: the latter carries nested Attributes, which the
			// Protocol v5 muxed server cannot serialize ("protocol version 5 cannot
			// have Attributes set"). An object-typed list is the v5-compatible form
			// and matches the state built by types.ListValueFrom.
			"robots": schema.ListAttribute{
				ElementType: robotObjectElementType(),
				Computed:    true,
				Description: "A list of alert robots. Each element contains: alert_robot_id, alert_robot_name, type, lang, url, workspace, digital_employee_name, robot_sign_key.",
			},
		},
	}
}

func (d *alertRobotsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config alertRobotsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := d.read(ctx, d.Client(), config)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Reading %s", alertRobotsDataSourceTypeName), err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (d *alertRobotsDataSource) read(ctx context.Context, client *connectivity.AliyunClient, config alertRobotsDataSourceModel) (alertRobotsDataSourceModel, error) {
	filterType := config.Type.ValueString()
	filterWorkspace := config.Workspace.ValueString()
	var filterIds []string
	if !config.Ids.IsNull() {
		_ = config.Ids.ElementsAs(ctx, &filterIds, false)
	}

	var nameRegex *regexp.Regexp
	if !config.NameRegex.IsNull() {
		pattern := config.NameRegex.ValueString()
		r, err := regexp.Compile(pattern)
		if err != nil {
			return alertRobotsDataSourceModel{}, fmt.Errorf("%s: name_regex format error: %s", alertRobotsDataSourceTypeName, err)
		}
		nameRegex = r
	}

	all, err := listAlertRobots(ctx, client, filterType, filterWorkspace, filterIds)
	if err != nil {
		return alertRobotsDataSourceModel{}, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, alertRobotsDataSourceTypeName, alertRobotsCollectionPath, errs.AlibabaCloudSdkGoERROR)
	}

	matched := make([]alertRobotDataSourceRobotModel, 0, len(all))
	matchedIds := make([]string, 0, len(all))
	for _, robot := range all {
		name, _ := robot["name"].(string)
		robotType, _ := robot["type"].(string)
		robotId, _ := robot["robotId"].(string)

		if nameRegex != nil && !nameRegex.MatchString(name) {
			continue
		}
		if filterType != "" && robotType != filterType {
			continue
		}
		if len(filterIds) > 0 && !contains(filterIds, robotId) {
			continue
		}

		matched = append(matched, alertRobotStateForDataSource(robot))
		matchedIds = append(matchedIds, robotId)
	}

	robotsList, diags := types.ListValueFrom(ctx, robotObjectElementType(), matched)
	if diags.HasError() {
		return alertRobotsDataSourceModel{}, fmt.Errorf("%s: building robots list: %v", alertRobotsDataSourceTypeName, diags)
	}
	idsList, diags := types.ListValueFrom(ctx, types.StringType, matchedIds)
	if diags.HasError() {
		return alertRobotsDataSourceModel{}, fmt.Errorf("%s: building ids list: %v", alertRobotsDataSourceTypeName, diags)
	}

	state := alertRobotsDataSourceModel{
		Id:         types.StringValue(dataSourceIdHash(matchedIds)),
		Ids:        idsList,
		NameRegex:  config.NameRegex,
		Type:       config.Type,
		Workspace:  config.Workspace,
		OutputFile: config.OutputFile,
		Robots:     robotsList,
	}

	if !config.OutputFile.IsNull() {
		if err := writeRobotsToFile(config.OutputFile.ValueString(), matched); err != nil {
			return alertRobotsDataSourceModel{}, fmt.Errorf("%s: writing output file: %s", alertRobotsDataSourceTypeName, err)
		}
	}
	return state, nil
}

// listAlertRobots pages through GET /robots applying the server-side filters the
// API supports (types, workspace, robotIds), returning every page until exhausted.
func listAlertRobots(ctx context.Context, client *connectivity.AliyunClient, filterType, filterWorkspace string, filterIds []string) ([]map[string]interface{}, error) {
	const pageSize = 100
	pageNumber := 1
	var all []map[string]interface{}

	for {
		query := map[string]*string{
			"pageNumber": stringPointer(fmt.Sprintf("%d", pageNumber)),
			"pageSize":   stringPointer(fmt.Sprintf("%d", pageSize)),
		}
		if filterType != "" {
			typesJSON, _ := json.Marshal([]string{filterType})
			query["types"] = stringPointer(string(typesJSON))
		}
		if filterWorkspace != "" {
			query["workspace"] = stringPointer(filterWorkspace)
		}
		if len(filterIds) > 0 {
			robotIdsJSON, _ := json.Marshal(filterIds)
			query["robotIds"] = stringPointer(string(robotIdsJSON))
		}

		response, err := client.RoaGet(cmsProduct, cmsVersion, alertRobotsCollectionPath, query, nil, nil)
		if err != nil {
			return nil, err
		}

		robots, ok := response["robots"].([]interface{})
		if !ok {
			return all, nil
		}
		for _, raw := range robots {
			if robot, ok := raw.(map[string]interface{}); ok {
				all = append(all, robot)
			}
		}

		total := 0
		if v, ok := response["total"].(float64); ok {
			total = int(v)
		}
		if len(all) >= total || len(robots) == 0 {
			return all, nil
		}
		pageNumber++
	}
}

func alertRobotStateForDataSource(robot map[string]interface{}) alertRobotDataSourceRobotModel {
	str := func(key string) types.String {
		if v, ok := robot[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	return alertRobotDataSourceRobotModel{
		AlertRobotId:        str("robotId"),
		AlertRobotName:      str("name"),
		Type:                str("type"),
		Lang:                str("lang"),
		Url:                 str("webhook"),
		Workspace:           str("workspace"),
		DigitalEmployeeName: str("digitalEmployeeName"),
		RobotSignKey:        str("robotSignKey"),
	}
}

func robotObjectElementType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"alert_robot_id":        types.StringType,
			"alert_robot_name":      types.StringType,
			"type":                  types.StringType,
			"lang":                  types.StringType,
			"url":                   types.StringType,
			"workspace":             types.StringType,
			"digital_employee_name": types.StringType,
			"robot_sign_key":        types.StringType,
		},
	}
}

func dataSourceIdHash(ids []string) string {
	h := sha256.Sum256([]byte(strings.Join(ids, ",")))
	return hex.EncodeToString(h[:])
}

func writeRobotsToFile(path string, robots []alertRobotDataSourceRobotModel) error {
	out := make([]map[string]string, 0, len(robots))
	for _, r := range robots {
		out = append(out, map[string]string{
			"alert_robot_id":        r.AlertRobotId.ValueString(),
			"alert_robot_name":      r.AlertRobotName.ValueString(),
			"type":                  r.Type.ValueString(),
			"lang":                  r.Lang.ValueString(),
			"url":                   r.Url.ValueString(),
			"workspace":             r.Workspace.ValueString(),
			"digital_employee_name": r.DigitalEmployeeName.ValueString(),
			"robot_sign_key":        r.RobotSignKey.ValueString(),
		})
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}
