package alicloud

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/registry"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/scoped"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func assembleProvider(p *schema.Provider) diag.Diagnostics {
	return scoped.Assemble(p, registry.ServicePackages())
}

var formatRegistrationDiags = scoped.FormatRegistrationDiags
