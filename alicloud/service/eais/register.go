package eais

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Eais service: the
// SDK v2 resources and data sources migrated from the flat layout.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "eais",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_eais_instance", Factory: ResourceEaisInstance},
			{TypeName: "alicloud_eais_client_instance_attachment", Factory: ResourceEaisClientInstanceAttachment},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_eais_instances", Factory: DataSourceEaisInstances},
		},
	}
}
