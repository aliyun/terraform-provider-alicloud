package scdn

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage registers the scdn service's SDK v2 resources and data sources.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "scdn",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_scdn_domain", Factory: ResourceScdnDomain},
			{TypeName: "alicloud_scdn_domain_config", Factory: ResourceScdnDomainConfig},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_scdn_domains", Factory: DataSourceScdnDomains},
		},
	}
}
