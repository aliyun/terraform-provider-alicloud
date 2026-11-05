package ots

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the OTS (Tablestore) service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Ots",
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_ots_service", Factory: DataSourceOtsService},
			{TypeName: "alicloud_ots_instances", Factory: DataSourceOtsInstances},
			{TypeName: "alicloud_ots_instance_attachments", Factory: DataSourceOtsInstanceAttachments},
			{TypeName: "alicloud_ots_tables", Factory: DataSourceOtsTables},
			{TypeName: "alicloud_ots_tunnels", Factory: DataSourceOtsTunnels},
			{TypeName: "alicloud_ots_secondary_indexes", Factory: DataSourceOtsSecondaryIndexes},
			{TypeName: "alicloud_ots_search_indexes", Factory: DataSourceOtsSearchIndexes},
		},
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_ots_instance", Factory: ResourceOtsInstance},
			{TypeName: "alicloud_ots_instance_attachment", Factory: ResourceOtsInstanceAttachment},
			{TypeName: "alicloud_ots_search_index", Factory: ResourceOtsSearchIndex},
			{TypeName: "alicloud_ots_secondary_index", Factory: ResourceOtsSecondaryIndex},
			{TypeName: "alicloud_ots_table", Factory: ResourceOtsTable},
			{TypeName: "alicloud_ots_tunnel", Factory: ResourceOtsTunnel},
		},
	}
}
