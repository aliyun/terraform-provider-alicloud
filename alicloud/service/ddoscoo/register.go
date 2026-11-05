package ddoscoo

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the ddoscoo service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "ddoscoo",
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_ddoscoo_domain_resources",
				Factory:  DataSourceDdoscooDomainResources,
			},
			{
				TypeName: "alicloud_ddoscoo_instances",
				Factory:  DataSourceDdoscooInstances,
			},
			{
				TypeName: "alicloud_ddoscoo_ports",
				Factory:  DataSourceDdoscooPorts,
			},
		},
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_ddoscoo_domain_precise_access_rule",
				Factory:  ResourceDdoscooDomainPreciseAccessRule,
			},
			{
				TypeName: "alicloud_ddoscoo_domain_resource",
				Factory:  ResourceDdoscooDomainResource,
			},
			{
				TypeName: "alicloud_ddoscoo_instance",
				Factory:  ResourceDdoscooInstance,
			},
			{
				TypeName: "alicloud_ddoscoo_port",
				Factory:  ResourceDdoscooPort,
			},
			{
				TypeName: "alicloud_ddoscoo_scheduler_rule",
				Factory:  ResourceDdoscooSchedulerRule,
			},
			{
				TypeName: "alicloud_ddoscoo_web_cc_rule",
				Factory:  ResourceDdoscooWebCcRule,
			},
		},
	}
}
