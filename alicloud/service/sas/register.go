package sas

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Sas (Threat
// Detection / Security Center) service. The custom check standard policy
// resource and data source are framework-native, registered here so the
// provider assembly loop picks them up alongside the SDKv2 Sas resources
// that still live in package alicloud.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Sas",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_threat_detection_custom_check_standard_policy",
				Factory:  NewCustomCheckStandardPolicyResource,
			},
		},
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_threat_detection_custom_check_standard_policies",
				Factory:  NewCustomCheckStandardPoliciesDataSource,
			},
		},
	}
}
