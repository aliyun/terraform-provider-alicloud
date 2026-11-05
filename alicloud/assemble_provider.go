package alicloud

import (
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/intercept"
	providerregistry "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/registry"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/sdkv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Wraps what the provider.go map literals already hold with its per-name chain,
// then appends the registry's declarations. A declaration shadowing a map literal
// entry is a startup hard error, detectable only here.
func assembleProvider(p *schema.Provider) diag.Diagnostics {
	var diags diag.Diagnostics

	for name, r := range p.ResourcesMap {
		p.ResourcesMap[name] = sdkv2.WrapResource(name, r, intercept.ChainOf(name, nil))
	}
	for name, ds := range p.DataSourcesMap {
		p.DataSourcesMap[name] = sdkv2.WrapDataSource(name, ds, intercept.ChainOf(name, nil))
	}

	for _, sp := range providerregistry.ServicePackages() {
		for _, reg := range sp.SDKResources {
			name := reg.TypeName
			if _, dup := p.ResourcesMap[name]; dup {
				diags = append(diags, duplicateDeclaration("resource", name, sp.Name))
				continue
			}
			p.ResourcesMap[name] = sdkv2.WrapResource(name, reg.Factory(), intercept.ChainOf(name, reg.Interceptors))
		}
		for _, reg := range sp.SDKDataSources {
			name := reg.TypeName
			if _, dup := p.DataSourcesMap[name]; dup {
				diags = append(diags, duplicateDeclaration("data source", name, sp.Name))
				continue
			}
			p.DataSourcesMap[name] = sdkv2.WrapDataSource(name, reg.Factory(), intercept.ChainOf(name, reg.Interceptors))
		}
	}

	return diags
}

func duplicateDeclaration(noun, typeName, service string) diag.Diagnostic {
	return diag.Diagnostic{
		Severity: diag.Error,
		Summary:  fmt.Sprintf("duplicate %s %q", noun, typeName),
		Detail: fmt.Sprintf("declared both in the provider.go map literal and in the %s ServicePackage(); remove the map literal entry",
			service),
	}
}

func formatRegistrationDiags(diags diag.Diagnostics) string {
	lines := make([]string, 0, len(diags))
	for _, d := range diags {
		if d.Severity != diag.Error {
			continue
		}
		if d.Detail != "" {
			lines = append(lines, d.Summary+": "+d.Detail)
			continue
		}
		lines = append(lines, d.Summary)
	}
	return strings.Join(lines, "\n  ")
}
