// Package cloudapi registers the CloudAPI (API Gateway) service package for the
// framework half of the provider. The SDK v2 API Gateway resources live in
// package alicloud against the same CloudAPI/2016-07-14 product; only resources
// new on master are implemented here as framework-native, per the upstream
// baseline rule.
package cloudapi

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the CloudAPI service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "CloudApi",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_api_gateway_domain",
				Factory:  NewApiGatewayDomainDataSource,
			},
		},
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_api_gateway_domain",
				Factory:  NewApiGatewayDomainResource,
			},
		},
	}
}
