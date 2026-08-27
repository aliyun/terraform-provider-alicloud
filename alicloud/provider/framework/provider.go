package framework

import (
	"context"
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/intercept"
	providerregistry "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/registry"
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
	meta := p.primary.Meta()

	resp.ResourceData = meta
	resp.DataSourceData = meta
	resp.EphemeralResourceData = meta
	resp.ListResourceData = meta
	resp.ActionData = meta
}

// Resources and data sources go through fwadapt so their chain applies. Ephemeral
// resources, list resources and actions have no wrapper yet — each needs its own,
// capability detection being by type assertion. Functions differ in kind: no
// FunctionWithConfigure, and a function call is not one of intercept's CRUD ops.

func (p *alicloudProvider) Resources(ctx context.Context) []func() resource.Resource {
	var problems []string
	factories := collect(
		func(sp conns.ServicePackage) []conns.Resource { return sp.Resources },
		func(d conns.Resource) func() resource.Resource {
			problems = appendUnlessClientCapable(problems, "resource", d.TypeName, "ResourceBase", d.Factory())
			chain := intercept.ChainOf(d.TypeName, d.Interceptors)
			return func() resource.Resource {
				return fwadapt.WrapResource(d.TypeName, d.Factory(), chain)
			}
		})
	mustRegister(problems)
	return factories
}

func (p *alicloudProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	var problems []string
	factories := collect(
		func(sp conns.ServicePackage) []conns.DataSource { return sp.DataSources },
		func(d conns.DataSource) func() datasource.DataSource {
			problems = appendUnlessClientCapable(problems, "data source", d.TypeName, "DataSourceBase", d.Factory())
			chain := intercept.ChainOf(d.TypeName, d.Interceptors)
			return func() datasource.DataSource {
				return fwadapt.WrapDataSource(d.TypeName, d.Factory(), chain)
			}
		})
	mustRegister(problems)
	return factories
}

func (p *alicloudProvider) Functions(ctx context.Context) []func() function.Function {
	return collect(
		func(sp conns.ServicePackage) []conns.Function { return sp.Functions },
		func(d conns.Function) func() function.Function { return d.Factory })
}

func (p *alicloudProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	var problems []string
	factories := collect(
		func(sp conns.ServicePackage) []conns.EphemeralResource { return sp.EphemeralResources },
		func(d conns.EphemeralResource) func() ephemeral.EphemeralResource {
			problems = appendUnlessClientCapable(problems, "ephemeral resource", d.TypeName, "EphemeralResourceBase", d.Factory())
			problems = appendIfInterceptorsDropped(problems, "ephemeral resource", d.TypeName, intercept.ChainOf(d.TypeName, nil))
			return d.Factory
		})
	mustRegister(problems)
	return factories
}

func (p *alicloudProvider) ListResources(ctx context.Context) []func() list.ListResource {
	var problems []string
	factories := collect(
		func(sp conns.ServicePackage) []conns.ListResource { return sp.ListResources },
		func(d conns.ListResource) func() list.ListResource {
			problems = appendUnlessClientCapable(problems, "list resource", d.TypeName, "ListResourceBase", d.Factory())
			problems = appendIfInterceptorsDropped(problems, "list resource", d.TypeName, intercept.ChainOf(d.TypeName, nil))
			return d.Factory
		})
	mustRegister(problems)
	return factories
}

func (p *alicloudProvider) Actions(ctx context.Context) []func() action.Action {
	var problems []string
	factories := collect(
		func(sp conns.ServicePackage) []conns.Action { return sp.Actions },
		func(d conns.Action) func() action.Action {
			problems = appendUnlessClientCapable(problems, "action", d.TypeName, "ActionBase", d.Factory())
			problems = appendIfInterceptorsDropped(problems, "action", d.TypeName, intercept.ChainOf(d.TypeName, nil))
			return d.Factory
		})
	mustRegister(problems)
	return factories
}

func collect[D, T any](category func(conns.ServicePackage) []D, constructor func(D) func() T) []func() T {
	var result []func() T
	for _, sp := range providerregistry.ServicePackages() {
		for _, d := range category(sp) {
			result = append(result, constructor(d))
		}
	}
	return result
}

func appendUnlessClientCapable(problems []string, namespace, typeName, base string, instance interface{}) []string {
	if _, ok := instance.(fwadapt.WithClient); ok {
		return problems
	}
	return append(problems, fmt.Sprintf(
		"%s %q: %T does not implement fwadapt.WithClient — embed fwadapt.%s, and return a pointer from the factory",
		namespace, typeName, instance, base))
}

// The chain is a parameter because intercept's registry has no setter.
func appendIfInterceptorsDropped(problems []string, namespace, typeName string, chain []intercept.Interceptor) []string {
	if len(chain) == 0 {
		return problems
	}
	return append(problems, fmt.Sprintf(
		"%s %q: %d interceptor(s) resolve for this name but %ss are not wrapped, so they would never run — "+
			"write the fwadapt wrapper for this category, or drop the registration",
		namespace, typeName, len(chain), namespace))
}

// Panics: these aggregator methods return a bare slice, with no diagnostics channel.
func mustRegister(problems []string) {
	if len(problems) == 0 {
		return
	}
	panic(fmt.Sprintf("provider registration error: %d framework declaration(s) rejected:\n  %s",
		len(problems), strings.Join(problems, "\n  ")))
}
