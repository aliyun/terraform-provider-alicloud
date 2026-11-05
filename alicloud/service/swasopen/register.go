package swasopen

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage registers the swasopen service's SDK v2 resources and data sources.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "SWAS-OPEN",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_simple_application_server_custom_image", Factory: ResourceCustomImage},
			{TypeName: "alicloud_simple_application_server_disk", Factory: ResourceDisk},
			{TypeName: "alicloud_simple_application_server_firewall_rule", Factory: ResourceFirewallRule},
			{TypeName: "alicloud_simple_application_server_instance", Factory: ResourceInstance},
			{TypeName: "alicloud_simple_application_server_snapshot", Factory: ResourceSnapshot},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_simple_application_server_custom_images", Factory: DataSourceCustomImages},
			{TypeName: "alicloud_simple_application_server_disks", Factory: DataSourceDisks},
			{TypeName: "alicloud_simple_application_server_firewall_rules", Factory: DataSourceFirewallRules},
			{TypeName: "alicloud_simple_application_server_images", Factory: DataSourceImages},
			{TypeName: "alicloud_simple_application_server_instances", Factory: DataSourceInstances},
			{TypeName: "alicloud_simple_application_server_plans", Factory: DataSourcePlans},
			{TypeName: "alicloud_simple_application_server_snapshots", Factory: DataSourceSnapshots},
		},
	}
}
