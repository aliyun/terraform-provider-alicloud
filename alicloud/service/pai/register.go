// Package pai registers the PAI (PaiStudio) Framework data sources.
//
// The PAI product family (EAS, PAIFlow, AIWorkSpace) keeps its existing SDKv2
// resources in the root alicloud package; new Framework data sources for the
// PaiStudio/2021-02-02 product live here so they do not require a framework
// conversion of the whole product.
package pai

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PAI service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PAIStudio",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_pai_node_output",
				Factory:  NewNodeOutputDataSource,
			},
			{
				TypeName: "alicloud_pai_node_outputs",
				Factory:  NewNodeOutputsDataSource,
			},
		},
	}
}
