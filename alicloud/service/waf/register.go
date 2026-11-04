package waf

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

// ServicePackage registers the waf service's SDK v2 resources and data sources.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Waf",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_waf_certificate", Factory: ResourceCertificate},
			{TypeName: "alicloud_waf_domain", Factory: ResourceDomain},
			{TypeName: "alicloud_waf_instance", Factory: ResourceInstance},
			{TypeName: "alicloud_waf_protection_module", Factory: ResourceProtectionModule},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_waf_certificates", Factory: DataSourceCertificates},
			{TypeName: "alicloud_waf_domains", Factory: DataSourceDomains},
			{TypeName: "alicloud_waf_instances", Factory: DataSourceInstances},
		},
	}
}
