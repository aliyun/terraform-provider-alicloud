package eflo

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{Name: "Eflo", Resources: []conns.Resource{{TypeName: "alicloud_eflo_vcc_grant_rule", Factory: NewVccGrantRuleResource}}}
}
