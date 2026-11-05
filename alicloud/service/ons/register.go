package ons

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Ons",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_ons_instance", Factory: ResourceOnsInstance},
			{TypeName: "alicloud_ons_topic", Factory: ResourceOnsTopic},
			{TypeName: "alicloud_ons_group", Factory: ResourceOnsGroup},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_ons_instances", Factory: DataSourceOnsInstances},
			{TypeName: "alicloud_ons_topics", Factory: DataSourceOnsTopics},
			{TypeName: "alicloud_ons_groups", Factory: DataSourceOnsGroups},
			{TypeName: "alicloud_ons_service", Factory: DataSourceOnsService},
		},
	}
}
