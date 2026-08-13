package framework

import (
	"context"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ims"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	_ provider.Provider                       = &alicloudProvider{}
	_ provider.ProviderWithEphemeralResources = &alicloudProvider{}
	_ provider.ProviderWithFunctions          = &alicloudProvider{}
	_ provider.ProviderWithListResources      = &alicloudProvider{}
	_ provider.ProviderWithActions            = &alicloudProvider{}
)

type alicloudProvider struct {
	primary *sdkschema.Provider
}

// NewProvider returns the framework provider to be muxed with the SDK v2 provider.
func NewProvider(primary *sdkschema.Provider) provider.Provider {
	return &alicloudProvider{
		primary: primary,
	}
}

func (p *alicloudProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "alicloud"
}

func (p *alicloudProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	schema, diags := ProviderSchema(p.primary.Schema)

	resp.Schema = schema
	resp.Diagnostics.Append(diags...)
}

func (p *alicloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// The provider's parsed configuration is available through the primary provider's
	// Meta method: tf5muxserver configures its servers in the order they were passed to
	// NewMuxServer, and the SDK v2 one goes first. Credential resolution lives in
	// alicloud.providerConfigure and is not duplicated here.
	meta := p.primary.Meta()

	resp.ResourceData = meta
	resp.DataSourceData = meta
	resp.EphemeralResourceData = meta
	resp.ListResourceData = meta
	resp.ActionData = meta
}

func (p *alicloudProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		ims.NewDefaultDomainDataSource,
	}
}

func (p *alicloudProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *alicloudProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *alicloudProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func (p *alicloudProvider) ListResources(ctx context.Context) []func() list.ListResource {
	return []func() list.ListResource{}
}

func (p *alicloudProvider) Actions(ctx context.Context) []func() action.Action {
	return []func() action.Action{}
}
