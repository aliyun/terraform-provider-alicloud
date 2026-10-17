package paifs

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PaiFeatureStore
// service. Both the managed resource and the list data source are Terraform
// Plugin Framework native implementations.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PaiFeatureStore",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_pai_feature_store_datasource",
				Factory:  NewDatasourceResource,
			},
		},
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_pai_feature_store_datasources",
				Factory:  NewDatasourcesDataSource,
			},
		},
	}
}
