// Package eflo registers the framework-native Eflo resources and data sources
// with the provider. The pre-existing SDK v2 Eflo resources live in package
// alicloud and are registered through the legacy provider.go schema map; only
// framework-native additions belong here.
package eflo

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

// ServicePackage declares the framework-native Eflo surface.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Eflo",
		Resources: []conns.Resource{
			{TypeName: "alicloud_eflo_vcc_route_entry", Factory: NewVccRouteEntryResource},
		},
		DataSources: []conns.DataSource{
			{TypeName: "alicloud_eflo_vcc_route_entries", Factory: NewVccRouteEntriesDataSource},
		},
	}
}
