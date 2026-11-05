package actiontrail

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

// ServicePackage registers the actiontrail service's SDK v2 resources and data sources.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Actiontrail",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_actiontrail_trail", Factory: ResourceTrail},
			{TypeName: "alicloud_actiontrail", Factory: ResourceTrail},
			{TypeName: "alicloud_actiontrail_advanced_query_template", Factory: ResourceAdvancedQueryTemplate},
			{TypeName: "alicloud_actiontrail_history_delivery_job", Factory: ResourceHistoryDeliveryJob},
			{TypeName: "alicloud_actiontrail_global_events_storage_region", Factory: ResourceGlobalEventsStorageRegion},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_actiontrail_trails", Factory: DataSourceTrails},
			{TypeName: "alicloud_actiontrails", Factory: DataSourceTrails},
			{TypeName: "alicloud_actiontrail_history_delivery_jobs", Factory: DataSourceHistoryDeliveryJobs},
			{TypeName: "alicloud_actiontrail_global_events_storage_region", Factory: DataSourceGlobalEventsStorageRegion},
		},
	}
}
