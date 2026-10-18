// Package dataworks holds the Framework-native (terraform-plugin-framework)
// implementations of the DataWorks resource and data source types that the
// flat alicloud layout registered through the legacy Plugin SDK.
package dataworks

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the DataWorks service.
// Both factories are plain functions so the Framework aggregator can wrap them
// with the interceptor chain (empty for now) and inject the provider client.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "DataWorks",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_data_works_resource",
				Factory:  NewResourceDataSource,
			},
		},
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_data_works_resource",
				Factory:  NewResourceResource,
			},
		},
	}
}
