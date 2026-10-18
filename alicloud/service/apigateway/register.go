package apigateway

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

// ServicePackage returns the registration declaration of the API Gateway
// (CloudAPI 2016-07-14) service. The Framework provider serves these types;
// the legacy SDK v2 api_gateway resources stay on the SDK v2 half.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "ApiGateway",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_api_gateway_datasets",
				Factory:  NewDatasetsDataSource,
			},
		},
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_api_gateway_dataset",
				Factory:  NewDatasetResource,
			},
		},
	}
}
