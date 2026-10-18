// Package dataworks carries the DataWorks resources and data sources that have
// been migrated to the plugin-framework native registration. DataWorks resources
// and data sources still written against SDK v2 remain in the root alicloud
// package and its provider.go map literals; only framework-native additions live
// here.
package dataworks

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the DataWorks service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "DataWorks",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_data_works_node_on_baseline",
				Factory:  NewNodeOnBaselineDataSource,
			},
			{
				TypeName: "alicloud_data_works_node_on_baselines",
				Factory:  NewNodeOnBaselinesDataSource,
			},
		},
	}
}
