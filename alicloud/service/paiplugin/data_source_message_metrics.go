package paiplugin

import (
	"context"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const messageMetricsDataSourceTypeName = "alicloud_pai_plugin_message_metrics"

var _ datasource.DataSource = &messageMetricsDataSource{}

func NewMessageMetricsDataSource() datasource.DataSource {
	return &messageMetricsDataSource{}
}

type messageMetricsDataSource struct {
	fwadapt.DataSourceBase
}

// messageMetricsModel flattens the per-date metrics collection into parallel
// arrays (dates, fails, pendings, rates, successes, totals). The array index
// aligns across all six lists, preserving the order returned by
// ListMessageMetrics; this avoids nested-object lists whose key order is a
// compile-time hazard in this no-local-build setting.
type messageMetricsModel struct {
	Id           types.String `tfsdk:"id"`
	StartDate    types.String `tfsdk:"start_date"`
	EndDate      types.String `tfsdk:"end_date"`
	GroupId      types.String `tfsdk:"group_id"`
	ScheduleId   types.String `tfsdk:"schedule_id"`
	Signature    types.String `tfsdk:"signature"`
	SignatureId  types.String `tfsdk:"signature_id"`
	TemplateCode types.String `tfsdk:"template_code"`
	TemplateId   types.String `tfsdk:"template_id"`
	TemplateType types.Int64  `tfsdk:"template_type"`
	Dates        types.List   `tfsdk:"dates"`
	Fails        types.List   `tfsdk:"fails"`
	Pendings     types.List   `tfsdk:"pendings"`
	Rates        types.List   `tfsdk:"rates"`
	Successes    types.List   `tfsdk:"successes"`
	Totals       types.List   `tfsdk:"totals"`
	TotalCount   types.Int64  `tfsdk:"total_count"`
}

func (d *messageMetricsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = messageMetricsDataSourceTypeName
}

func (d *messageMetricsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides PAI Plugin message delivery metrics between the given dates.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source.",
			},
			"start_date": schema.StringAttribute{
				Required:    true,
				Description: "The start date (yyyymmdd), inclusive.",
			},
			"end_date": schema.StringAttribute{
				Required:    true,
				Description: "The end date (yyyymmdd), inclusive.",
			},
			"group_id": schema.StringAttribute{
				Optional:    true,
				Description: "The group id filter.",
			},
			"schedule_id": schema.StringAttribute{
				Optional:    true,
				Description: "The schedule id filter.",
			},
			"signature": schema.StringAttribute{
				Optional:    true,
				Description: "The signature name filter.",
			},
			"signature_id": schema.StringAttribute{
				Optional:    true,
				Description: "The signature id filter.",
			},
			"template_code": schema.StringAttribute{
				Optional:    true,
				Description: "The template code filter.",
			},
			"template_id": schema.StringAttribute{
				Optional:    true,
				Description: "The template id filter.",
			},
			"template_type": schema.Int64Attribute{
				Optional:    true,
				Description: "The template type filter. 0: verification code, 1: notification, 2: marketing.",
			},
			"dates": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "The dates of the metrics, in the order returned by ListMessageMetrics.",
			},
			"fails": schema.ListAttribute{
				Computed:    true,
				ElementType: types.Int64Type,
				Description: "The failed counts, aligned by index with dates.",
			},
			"pendings": schema.ListAttribute{
				Computed:    true,
				ElementType: types.Int64Type,
				Description: "The pending counts, aligned by index with dates.",
			},
			"rates": schema.ListAttribute{
				Computed:    true,
				ElementType: types.Float64Type,
				Description: "The success rates, aligned by index with dates.",
			},
			"successes": schema.ListAttribute{
				Computed:    true,
				ElementType: types.Int64Type,
				Description: "The success counts, aligned by index with dates.",
			},
			"totals": schema.ListAttribute{
				Computed:    true,
				ElementType: types.Int64Type,
				Description: "The total counts, aligned by index with dates.",
			},
			"total_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The total count of metrics records returned.",
			},
		},
	}
}

func (d *messageMetricsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config messageMetricsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// template_type is optional and may be null; build the query only when set.
	var templateType string
	if !config.TemplateType.IsNull() && !config.TemplateType.IsUnknown() {
		templateType = fmt.Sprintf("%d", config.TemplateType.ValueInt64())
	}

	query := queryFrom(
		"StartDate", config.StartDate.ValueString(),
		"EndDate", config.EndDate.ValueString(),
		"GroupId", config.GroupId.ValueString(),
		"ScheduleId", config.ScheduleId.ValueString(),
		"Signature", config.Signature.ValueString(),
		"SignatureId", config.SignatureId.ValueString(),
		"TemplateCode", config.TemplateCode.ValueString(),
		"TemplateId", config.TemplateId.ValueString(),
		"TemplateType", templateType,
	)

	response, err := roaGet(d.Client(), pathMsgMetrics, query)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListMessageMetrics", messageMetricsDataSourceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling ListMessageMetrics", messageMetricsDataSourceTypeName).Error(),
		)
		return
	}

	data := dataObject(response)
	metrics := objectSlice(data, "Metrics")

	dates := make([]string, 0, len(metrics))
	fails := make([]int64, 0, len(metrics))
	pendings := make([]int64, 0, len(metrics))
	rates := make([]float64, 0, len(metrics))
	successes := make([]int64, 0, len(metrics))
	totals := make([]int64, 0, len(metrics))
	for _, m := range metrics {
		dates = append(dates, str(m, "Date"))
		fails = append(fails, intValue(m, "Fail"))
		pendings = append(pendings, intValue(m, "Pending"))
		rates = append(rates, floatValue(m, "Rate"))
		successes = append(successes, intValue(m, "Success"))
		totals = append(totals, intValue(m, "Total"))
	}

	state := messageMetricsModel{
		Id:           types.StringValue(messageMetricsDataSourceTypeName),
		StartDate:    config.StartDate,
		EndDate:      config.EndDate,
		GroupId:      config.GroupId,
		ScheduleId:   config.ScheduleId,
		Signature:    config.Signature,
		SignatureId:  config.SignatureId,
		TemplateCode: config.TemplateCode,
		TemplateId:   config.TemplateId,
		TemplateType: config.TemplateType,
		Dates:        stringList(dates),
		Fails:        intSliceValue(fails),
		Pendings:     intSliceValue(pendings),
		Rates:        floatSliceValue(rates),
		Successes:    intSliceValue(successes),
		Totals:       intSliceValue(totals),
		TotalCount:   types.Int64Value(intValue(data, "TotalCount")),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
