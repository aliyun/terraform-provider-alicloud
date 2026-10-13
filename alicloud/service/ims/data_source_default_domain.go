package ims

import (
	"context"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const defaultDomainTypeName = "alicloud_ims_default_domain"

var (
	_ datasource.DataSource              = &defaultDomainDataSource{}
	_ datasource.DataSourceWithConfigure = &defaultDomainDataSource{}
)

func NewDefaultDomainDataSource() datasource.DataSource {
	return &defaultDomainDataSource{}
}

type defaultDomainDataSource struct {
	fwadapt.DataSourceBase
}

type defaultDomainModel struct {
	Id            types.String `tfsdk:"id"`
	DefaultDomain types.String `tfsdk:"default_domain"`
}

func (d *defaultDomainDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = defaultDomainTypeName
}

func (d *defaultDomainDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides the default domain of the Alibaba Cloud account, which is the suffix of every RAM user's logon name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The default domain, same as `default_domain`.",
			},
			"default_domain": schema.StringAttribute{
				Computed:    true,
				Description: "The default domain of the Alibaba Cloud account.",
			},
		},
	}
}

func (d *defaultDomainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	const action = "GetDefaultDomain"

	response, err := d.Client().RpcPost("Ims", "2019-08-15", action, nil, map[string]interface{}{}, true)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Reading %s: calling %s", defaultDomainTypeName, action), err.Error())
		return
	}

	domain, _ := response["DefaultDomainName"].(string)
	if domain == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", defaultDomainTypeName, action),
			fmt.Sprintf("The response carried no DefaultDomainName: %v", response),
		)
		return
	}

	state := defaultDomainModel{
		Id:            types.StringValue(domain),
		DefaultDomain: types.StringValue(domain),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
