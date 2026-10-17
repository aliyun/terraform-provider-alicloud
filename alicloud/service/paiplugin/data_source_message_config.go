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

const messageConfigDataSourceTypeName = "alicloud_pai_plugin_message_config"

var _ datasource.DataSource = &messageConfigDataSource{}

func NewMessageConfigDataSource() datasource.DataSource {
	return &messageConfigDataSource{}
}

type messageConfigDataSource struct {
	fwadapt.DataSourceBase
}

type messageConfigModel struct {
	Id           types.String `tfsdk:"id"`
	SmsReportUrl types.String `tfsdk:"sms_report_url"`
	SmsUpUrl     types.String `tfsdk:"sms_up_url"`
}

func (d *messageConfigDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = messageConfigDataSourceTypeName
}

func (d *messageConfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides the SMS message config of the current account (singleton).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source.",
			},
			"sms_report_url": schema.StringAttribute{
				Computed:    true,
				Description: "The URL that receives the SMS delivery report callback.",
			},
			"sms_up_url": schema.StringAttribute{
				Computed:    true,
				Description: "The URL that receives the uplink SMS message callback.",
			},
		},
	}
}

func (d *messageConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	response, err := roaGet(d.Client(), pathMsgConfig, nil)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Reading %s: calling GetMessageConfig", messageConfigDataSourceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetMessageConfig", messageConfigDataSourceTypeName).Error())
		return
	}

	data := dataObject(response)
	if data == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetMessageConfig", messageConfigDataSourceTypeName),
			fmt.Sprintf("The response carried no Data: %v", response),
		)
		return
	}

	state := messageConfigModel{
		Id:           types.StringValue(messageConfigDataSourceTypeName),
		SmsReportUrl: strValue(data, "SmsReportUrl"),
		SmsUpUrl:     strValue(data, "SmsUpUrl"),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
