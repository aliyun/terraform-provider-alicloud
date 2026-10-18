package apigateway

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the API Gateway (CloudAPI)
// service. The POP product code is CloudAPI and the API version is 2016-07-14.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "ApiGateway",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_api_gateway_signatures",
				Factory:  NewSignaturesDataSource,
			},
		},
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_api_gateway_signature",
				Factory:  NewSignatureResource,
			},
		},
	}
}
