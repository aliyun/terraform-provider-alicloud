package scdn

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the SCDN service.
// Both resources are deprecated since v1.219.0 but keep their original public
// type names and deprecation markers, so existing configurations continue to
// import and plan unchanged.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "scdn",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_scdn_domain",
				Factory:  ResourceScdnDomain,
			},
			{
				TypeName: "alicloud_scdn_domain_config",
				Factory:  ResourceScdnDomainConfig,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_scdn_domains",
				Factory:  DataSourceScdnDomains,
			},
		},
	}
}
