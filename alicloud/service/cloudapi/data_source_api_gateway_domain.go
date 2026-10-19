package cloudapi

import (
	"context"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const apiGatewayDomainDataSourceTypeName = "alicloud_api_gateway_domain"

var (
	_ datasource.DataSource              = &apiGatewayDomainDataSource{}
	_ datasource.DataSourceWithConfigure = &apiGatewayDomainDataSource{}
)

func NewApiGatewayDomainDataSource() datasource.DataSource {
	return &apiGatewayDomainDataSource{}
}

type apiGatewayDomainDataSource struct {
	fwadapt.DataSourceBase
}

type apiGatewayDomainDataSourceModel struct {
	Id                    types.String `tfsdk:"id"`
	GroupId               types.String `tfsdk:"group_id"`
	DomainName            types.String `tfsdk:"domain_name"`
	SubDomain             types.String `tfsdk:"sub_domain"`
	DomainBindingStatus   types.String `tfsdk:"domain_binding_status"`
	DomainRemark          types.String `tfsdk:"domain_remark"`
	DomainWebSocketStatus types.String `tfsdk:"domain_web_socket_status"`
	DomainLegalStatus     types.String `tfsdk:"domain_legal_status"`
	DomainCNAMEStatus     types.String `tfsdk:"domain_cname_status"`
	CertificateId         types.String `tfsdk:"certificate_id"`
	CertificateName       types.String `tfsdk:"certificate_name"`
	CertificateValidStart types.Int64  `tfsdk:"certificate_valid_start"`
	CertificateValidEnd   types.Int64  `tfsdk:"certificate_valid_end"`
}

func (d *apiGatewayDomainDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = apiGatewayDomainDataSourceTypeName
}

func (d *apiGatewayDomainDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an API Gateway domain data source that queries the binding information of a custom domain name attached to an API Gateway group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. The value formats as `<group_id>:<domain_name>`.",
			},
			"group_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the API Gateway group that the custom domain binds to.",
			},
			"domain_name": schema.StringAttribute{
				Required:    true,
				Description: "The custom domain name to query.",
			},
			"sub_domain": schema.StringAttribute{
				Computed:    true,
				Description: "The system allocated sub domain of the API Gateway group.",
			},
			"domain_binding_status": schema.StringAttribute{
				Computed:    true,
				Description: "The binding status of the domain.",
			},
			"domain_remark": schema.StringAttribute{
				Computed:    true,
				Description: "The remark of the domain.",
			},
			"domain_web_socket_status": schema.StringAttribute{
				Computed:    true,
				Description: "The WebSocket status of the domain.",
			},
			"domain_legal_status": schema.StringAttribute{
				Computed:    true,
				Description: "The legal status of the domain.",
			},
			"domain_cname_status": schema.StringAttribute{
				Computed:    true,
				Description: "The CNAME status of the domain.",
			},
			"certificate_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the bound SSL certificate.",
			},
			"certificate_name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the bound SSL certificate.",
			},
			"certificate_valid_start": schema.Int64Attribute{
				Computed:    true,
				Description: "The valid start timestamp (milliseconds) of the bound SSL certificate.",
			},
			"certificate_valid_end": schema.Int64Attribute{
				Computed:    true,
				Description: "The valid end timestamp (milliseconds) of the bound SSL certificate.",
			},
		},
	}
}

func (d *apiGatewayDomainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config apiGatewayDomainDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId := config.GroupId.ValueString()
	domainName := config.DomainName.ValueString()

	response, err := describeDomain(d.Client(), groupId, domainName)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeDomain", apiGatewayDomainDataSourceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeDomain", apiGatewayDomainDataSourceTypeName).Error(),
		)
		return
	}
	if response == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeDomain", apiGatewayDomainDataSourceTypeName),
			fmt.Sprintf("The domain %s of group %s is not found", domainName, groupId),
		)
		return
	}

	state := apiGatewayDomainDataSourceModel{
		Id:                    types.StringValue(domainId(groupId, domainName)),
		GroupId:               strValue(response, "GroupId"),
		DomainName:            strValue(response, "DomainName"),
		SubDomain:             strValue(response, "SubDomain"),
		DomainBindingStatus:   strValue(response, "DomainBindingStatus"),
		DomainRemark:          strValue(response, "DomainRemark"),
		DomainWebSocketStatus: strValue(response, "DomainWebSocketStatus"),
		DomainLegalStatus:     strValue(response, "DomainLegalStatus"),
		DomainCNAMEStatus:     strValue(response, "DomainCNAMEStatus"),
		CertificateId:         strValue(response, "CertificateId"),
		CertificateName:       strValue(response, "CertificateName"),
		CertificateValidStart: int64Value(response, "CertificateValidStart"),
		CertificateValidEnd:   int64Value(response, "CertificateValidEnd"),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
