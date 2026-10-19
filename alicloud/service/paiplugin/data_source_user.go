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

const userDataSourceTypeName = "alicloud_pai_plugin_user"

var _ datasource.DataSource = &userDataSource{}

func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

type userDataSource struct {
	fwadapt.DataSourceBase
}

type userModel struct {
	Id            types.String `tfsdk:"id"`
	AccountStatus types.Int64  `tfsdk:"account_status"`
}

func (d *userDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = userDataSourceTypeName
}

func (d *userDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides the PAI Plugin user of the current account (singleton).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source.",
			},
			"account_status": schema.Int64Attribute{
				Computed:    true,
				Description: "The account status.",
			},
		},
	}
}

func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	response, err := roaGet(d.Client(), pathUsers, nil)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Reading %s: calling GetUser", userDataSourceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetUser", userDataSourceTypeName).Error())
		return
	}

	data := dataObject(response)
	if data == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetUser", userDataSourceTypeName),
			fmt.Sprintf("The response carried no Data: %v", response),
		)
		return
	}

	state := userModel{
		Id:            types.StringValue(userDataSourceTypeName),
		AccountStatus: types.Int64Value(intValue(data, "AccountStatus")),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
