package alicloud

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns the complete SDK provider used by the executable and legacy tests.
func Provider() *schema.Provider {
	p := &schema.Provider{
		ResourcesMap:   map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{},
	}
	config.Configure(p)
	if diags := assembleProvider(p); diags.HasError() {
		panic("provider registration error: " + formatRegistrationDiags(diags))
	}
	return p
}

// CredentialsURIResponse preserves the response type exposed by the root package.
type CredentialsURIResponse = config.CredentialsURIResponse
