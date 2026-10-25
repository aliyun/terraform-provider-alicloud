package eflo

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Eflo service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "eflo",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_eflo_vcc_grant_rule",
				Factory:  NewVccGrantRuleResource,
			},
		},
	}
}
