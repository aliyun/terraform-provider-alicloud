package sts

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the STS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name:           "Sts",
		SDKDataSources: []conns.SDKDataSource{{TypeName: "alicloud_caller_identity", Factory: dataSourceAlicloudCallerIdentity}},
		EphemeralResources: []conns.EphemeralResource{
			{
				TypeName: "alicloud_sts_token",
				Factory:  NewStsTokenEphemeralResource,
			},
		},
	}
}
