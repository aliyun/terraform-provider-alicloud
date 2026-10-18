package dataworks

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const qualityCheckResultsTypeName = "alicloud_data_works_quality_check_results"

var (
	_ datasource.DataSource              = &qualityCheckResultsDataSource{}
	_ datasource.DataSourceWithConfigure = &qualityCheckResultsDataSource{}
)

// NewQualityCheckResultsDataSource returns the Framework-native data source that
// lists DataWorks quality check results.
func NewQualityCheckResultsDataSource() datasource.DataSource {
	return &qualityCheckResultsDataSource{}
}

type qualityCheckResultsDataSource struct {
	fwadapt.DataSourceBase
}

// rpcPostCaller is the minimal surface of *connectivity.AliyunClient this data
// source needs. It is an interface so the Read path can be exercised without a
// cloud account by injecting a fake caller in unit tests.
type rpcPostCaller interface {
	RpcPost(apiProductCode string, apiVersion string, apiName string, query map[string]interface{}, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error)
}

type qualityCheckParams struct {
	ProjectName string
	StartDate   string
	EndDate     string
	EntityId    string
	RuleId      string
	PageSize    int
	PageNumber  int
	IDs         []string
}

type qualityCheckResultModel struct {
	Id                   types.String `tfsdk:"id"`
	QualityCheckResultId types.String `tfsdk:"quality_check_result_id"`
	EntityId             types.String `tfsdk:"entity_id"`
	RuleId               types.String `tfsdk:"rule_id"`
	ProjectName          types.String `tfsdk:"project_name"`
	TableName            types.String `tfsdk:"table_name"`
	Op                   types.String `tfsdk:"op"`
	TaskId               types.String `tfsdk:"task_id"`
	CheckerName          types.String `tfsdk:"checker_name"`
	CheckerId            types.String `tfsdk:"checker_id"`
	CheckerType          types.String `tfsdk:"checker_type"`
	BeginTime            types.String `tfsdk:"begin_time"`
	EndTime              types.String `tfsdk:"end_time"`
	ExpectValue          types.String `tfsdk:"expect_value"`
	UpperValue           types.String `tfsdk:"upper_value"`
	LowerValue           types.String `tfsdk:"lower_value"`
	WarningThreshold     types.String `tfsdk:"warning_threshold"`
	CriticalThreshold    types.String `tfsdk:"critical_threshold"`
	Trend                types.String `tfsdk:"trend"`
	CheckResult          types.String `tfsdk:"check_result"`
	FixedCheck           types.String `tfsdk:"fixed_check"`
	IsPrediction         types.String `tfsdk:"is_prediction"`
	BlockType            types.String `tfsdk:"block_type"`
	MatchExpression      types.String `tfsdk:"match_expression"`
	ActualExpression     types.String `tfsdk:"actual_expression"`
	WhereCondition       types.String `tfsdk:"where_condition"`
	Property             types.String `tfsdk:"property"`
	MethodName           types.String `tfsdk:"method_name"`
	DateType             types.String `tfsdk:"date_type"`
	RuleName             types.String `tfsdk:"rule_name"`
	TemplateName         types.String `tfsdk:"template_name"`
	TemplateId           types.String `tfsdk:"template_id"`
	Comment              types.String `tfsdk:"comment"`
	ExternalId           types.String `tfsdk:"external_id"`
	ExternalType         types.String `tfsdk:"external_type"`
}

type qualityCheckResultsModel struct {
	Id          types.String              `tfsdk:"id"`
	IDs         []types.String            `tfsdk:"ids"`
	ProjectName types.String              `tfsdk:"project_name"`
	StartDate   types.String              `tfsdk:"start_date"`
	EndDate     types.String              `tfsdk:"end_date"`
	EntityId    types.String              `tfsdk:"entity_id"`
	RuleId      types.String              `tfsdk:"rule_id"`
	PageSize    types.Int64               `tfsdk:"page_size"`
	PageNumber  types.Int64               `tfsdk:"page_number"`
	Results     []qualityCheckResultModel `tfsdk:"results"`
}

func (d *qualityCheckResultsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = qualityCheckResultsTypeName
}

func (d *qualityCheckResultsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides the Data Works Quality Check Results of the current Alibaba Cloud user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source. It is a hash of the quality check result IDs.",
			},
			"ids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of Quality Check Result IDs. When set, only results whose IDs are in this list are returned.",
			},
			"project_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the engine or data source.",
			},
			"start_date": schema.StringAttribute{
				Required:    true,
				Description: "The start of the business date range. The value must be in the `yyyy-MM-dd HH:mm:ss` format.",
			},
			"end_date": schema.StringAttribute{
				Required:    true,
				Description: "The end of the business date range. The value must be in the `yyyy-MM-dd HH:mm:ss` format.",
			},
			"entity_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the partition expression. Required when querying results by entity. Conflict with `rule_id`.",
			},
			"rule_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the monitoring rule. Required when querying results by rule. Conflict with `entity_id`.",
			},
			"page_size": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The number of entries per page. Default to 10. Maximum value is 20.",
			},
			"page_number": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The page number. Default to 1.",
			},
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: resultsNestedAttributes(),
				},
				Description: "A list of Data Works Quality Check Results.",
			},
		},
	}
}

func resultsNestedAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":                      schema.StringAttribute{Computed: true, Description: "The ID of the Quality Check Result."},
		"quality_check_result_id": schema.StringAttribute{Computed: true, Description: "The primary key ID of the Quality Check Result."},
		"entity_id":               schema.StringAttribute{Computed: true, Description: "The ID of the partition expression."},
		"rule_id":                 schema.StringAttribute{Computed: true, Description: "The ID of the monitoring rule."},
		"project_name":            schema.StringAttribute{Computed: true, Description: "The name of the project."},
		"table_name":              schema.StringAttribute{Computed: true, Description: "The name of the table."},
		"op":                      schema.StringAttribute{Computed: true, Description: "The operator."},
		"task_id":                 schema.StringAttribute{Computed: true, Description: "The ID of the task."},
		"checker_name":            schema.StringAttribute{Computed: true, Description: "The name of the checker."},
		"checker_id":              schema.StringAttribute{Computed: true, Description: "The ID of the checker."},
		"checker_type":            schema.StringAttribute{Computed: true, Description: "The type of the checker."},
		"begin_time":              schema.StringAttribute{Computed: true, Description: "The begin time."},
		"end_time":                schema.StringAttribute{Computed: true, Description: "The end time."},
		"expect_value":            schema.StringAttribute{Computed: true, Description: "The comparison value of the fixed value."},
		"upper_value":             schema.StringAttribute{Computed: true, Description: "The historical maximum value."},
		"lower_value":             schema.StringAttribute{Computed: true, Description: "The historical minimum value."},
		"warning_threshold":       schema.StringAttribute{Computed: true, Description: "The orange warning threshold."},
		"critical_threshold":      schema.StringAttribute{Computed: true, Description: "The red warning threshold."},
		"trend":                   schema.StringAttribute{Computed: true, Description: "The trend value."},
		"check_result":            schema.StringAttribute{Computed: true, Description: "The quality check result value."},
		"fixed_check":             schema.StringAttribute{Computed: true, Description: "Whether it is a fixed value check."},
		"is_prediction":           schema.StringAttribute{Computed: true, Description: "Whether it is a prediction."},
		"block_type":              schema.StringAttribute{Computed: true, Description: "Whether it is a strong rule."},
		"match_expression":        schema.StringAttribute{Computed: true, Description: "The partition expression."},
		"actual_expression":       schema.StringAttribute{Computed: true, Description: "The actual partition."},
		"where_condition":         schema.StringAttribute{Computed: true, Description: "The where condition."},
		"property":                schema.StringAttribute{Computed: true, Description: "The rule field."},
		"method_name":             schema.StringAttribute{Computed: true, Description: "The collection method."},
		"date_type":               schema.StringAttribute{Computed: true, Description: "The date type."},
		"rule_name":               schema.StringAttribute{Computed: true, Description: "The name of the rule."},
		"template_name":           schema.StringAttribute{Computed: true, Description: "The name of the template."},
		"template_id":             schema.StringAttribute{Computed: true, Description: "The ID of the template."},
		"comment":                 schema.StringAttribute{Computed: true, Description: "The comment of the rule."},
		"external_id":             schema.StringAttribute{Computed: true, Description: "The external node ID."},
		"external_type":           schema.StringAttribute{Computed: true, Description: "The external trigger type."},
	}
}

func (d *qualityCheckResultsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config qualityCheckResultsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := qualityCheckParams{
		ProjectName: config.ProjectName.ValueString(),
		StartDate:   config.StartDate.ValueString(),
		EndDate:     config.EndDate.ValueString(),
		EntityId:    config.EntityId.ValueString(),
		RuleId:      config.RuleId.ValueString(),
	}
	params.PageSize = resolvePageSize(config.PageSize)
	params.PageNumber = resolvePageNumber(config.PageNumber)
	for _, v := range config.IDs {
		if v.IsNull() {
			continue
		}
		params.IDs = append(params.IDs, v.ValueString())
	}

	items, err := d.listQualityCheckResults(ctx, d.Client(), params)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Reading %s", qualityCheckResultsTypeName), err.Error())
		return
	}

	state := qualityCheckResultsModel{
		ProjectName: types.StringValue(params.ProjectName),
		StartDate:   types.StringValue(params.StartDate),
		EndDate:     types.StringValue(params.EndDate),
		EntityId:    config.EntityId,
		RuleId:      config.RuleId,
		PageSize:    types.Int64Value(int64(params.PageSize)),
		PageNumber:  types.Int64Value(int64(params.PageNumber)),
		Results:     make([]qualityCheckResultModel, 0, len(items)),
		IDs:         make([]types.String, 0, len(items)),
	}
	for _, item := range items {
		state.Results = append(state.Results, mapQualityCheckResult(item))
		id := fmt.Sprint(item["Id"])
		state.IDs = append(state.IDs, types.StringValue(id))
	}
	state.Id = types.StringValue(qualityCheckResultsID(params, state.IDs))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// listQualityCheckResults calls ListQualityResultsByEntity or ListQualityResultsByRule
// with pagination until a partial page is returned, applies the optional ids
// filter, and returns the raw API objects.
func (d *qualityCheckResultsDataSource) listQualityCheckResults(ctx context.Context, client rpcPostCaller, params qualityCheckParams) ([]map[string]interface{}, error) {
	const (
		product        = "dataworks-public"
		version        = "2020-05-18"
		actionByEntity = "ListQualityResultsByEntity"
		actionByRule   = "ListQualityResultsByRule"
	)

	action := actionByEntity
	if params.RuleId != "" {
		action = actionByRule
	}

	var all []map[string]interface{}
	pageNumber := params.PageNumber
	for {
		request := map[string]interface{}{
			"ProjectName": params.ProjectName,
			"StartDate":   params.StartDate,
			"EndDate":     params.EndDate,
			"PageSize":    params.PageSize,
			"PageNumber":  pageNumber,
		}
		if action == actionByEntity && params.EntityId != "" {
			request["EntityId"] = params.EntityId
		}
		if action == actionByRule {
			request["RuleId"] = params.RuleId
		}

		response, err := client.RpcPost(product, version, action, nil, request, true)
		if err != nil {
			return nil, fmt.Errorf("calling %s: %w", action, err)
		}

		items, err := extractRuleChecks(response)
		if err != nil {
			return nil, fmt.Errorf("extracting quality check results from %s response: %w", action, err)
		}

		for _, v := range items {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			if len(params.IDs) > 0 && !containsString(params.IDs, fmt.Sprint(item["Id"])) {
				continue
			}
			all = append(all, item)
		}

		if len(items) < params.PageSize {
			break
		}
		pageNumber++
	}
	return all, nil
}

// extractRuleChecks navigates the $.Data.RuleChecks path of a DataWorks
// quality results response without pulling in a jsonpath dependency.
func extractRuleChecks(response map[string]interface{}) ([]interface{}, error) {
	if response == nil {
		return nil, fmt.Errorf("empty response")
	}
	data, ok := response["Data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("response carried no Data object")
	}
	raw, ok := data["RuleChecks"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("response carried no Data.RuleChecks list")
	}
	return raw, nil
}

func mapQualityCheckResult(item map[string]interface{}) qualityCheckResultModel {
	return qualityCheckResultModel{
		Id:                   stringValue(item["Id"]),
		QualityCheckResultId: stringValue(item["Id"]),
		EntityId:             stringValue(item["EntityId"]),
		RuleId:               stringValue(item["RuleId"]),
		ProjectName:          stringValue(item["ProjectName"]),
		TableName:            stringValue(item["TableName"]),
		Op:                   stringValue(item["Op"]),
		TaskId:               stringValue(item["TaskId"]),
		CheckerName:          stringValue(item["CheckerName"]),
		CheckerId:            stringValue(item["CheckerId"]),
		CheckerType:          stringValue(item["CheckerType"]),
		BeginTime:            stringValue(item["BeginTime"]),
		EndTime:              stringValue(item["EndTime"]),
		ExpectValue:          stringValue(item["ExpectValue"]),
		UpperValue:           stringValue(item["UpperValue"]),
		LowerValue:           stringValue(item["LowerValue"]),
		WarningThreshold:     stringValue(item["WarningThreshold"]),
		CriticalThreshold:    stringValue(item["CriticalThreshold"]),
		Trend:                stringValue(item["Trend"]),
		CheckResult:          stringValue(item["CheckResult"]),
		FixedCheck:           stringValue(item["FixedCheck"]),
		IsPrediction:         stringValue(item["IsPrediction"]),
		BlockType:            stringValue(item["BlockType"]),
		MatchExpression:      stringValue(item["MatchExpression"]),
		ActualExpression:     stringValue(item["ActualExpression"]),
		WhereCondition:       stringValue(item["WhereCondition"]),
		Property:             stringValue(item["Property"]),
		MethodName:           stringValue(item["MethodName"]),
		DateType:             stringValue(item["DateType"]),
		RuleName:             stringValue(item["RuleName"]),
		TemplateName:         stringValue(item["TemplateName"]),
		TemplateId:           stringValue(item["TemplateId"]),
		Comment:              stringValue(item["Comment"]),
		ExternalId:           stringValue(item["ExternalId"]),
		ExternalType:         stringValue(item["ExternalType"]),
	}
}

func stringValue(v interface{}) types.String {
	return types.StringValue(fmt.Sprint(v))
}

func resolvePageSize(v types.Int64) int {
	if v.IsNull() || v.ValueInt64() <= 0 {
		return 10
	}
	size := int(v.ValueInt64())
	if size > 20 {
		return 10
	}
	return size
}

func resolvePageNumber(v types.Int64) int {
	if v.IsNull() || v.ValueInt64() <= 0 {
		return 1
	}
	return int(v.ValueInt64())
}

func containsString(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func qualityCheckResultsID(params qualityCheckParams, ids []types.String) string {
	if len(ids) == 0 {
		return fmt.Sprintf("%s|%s|%s", params.ProjectName, params.StartDate, params.EndDate)
	}
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, id.ValueString())
	}
	sum := sha1.Sum([]byte(strings.Join(parts, ",")))
	return hex.EncodeToString(sum[:])
}
