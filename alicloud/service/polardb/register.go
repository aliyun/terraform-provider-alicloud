// Package polardb registers the polardb product's SDKv2 resources and data
// sources with the provider assembly loop. The flat alicloud package used to
// hold every polardb type directly in provider.go; migrating alicloud_polardb_application
// out means its registration moves here too, so provider.go stops listing it.
package polardb

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

// ServicePackage declares the SDKv2 types this product package owns. Name is
// the API product code surfaced in registration diagnostics, matching the
// "PolarDB" naming used throughout the dispatch and the existing flat-layout
// resource names.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PolarDB",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_polardb_application", Factory: ResourceApplication},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_polardb_node_classes",
				Factory:  DataSourceNodeClasses,
			},
		},
	}
}
