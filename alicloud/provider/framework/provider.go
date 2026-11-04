package framework

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/framework/scoped"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/registry"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// NewProvider retains the complete Framework provider for existing callers.
func NewProvider(primary *sdkschema.Provider) provider.Provider {
	return scoped.NewProvider(primary, registry.ServicePackages())
}
