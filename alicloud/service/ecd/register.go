package ecd

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the ecd service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Ecd",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_ecd_ad_connector_directory", Factory: ResourceAdConnectorDirectory},
			{TypeName: "alicloud_ecd_ad_connector_office_site", Factory: ResourceAdConnectorOfficeSite},
			{TypeName: "alicloud_ecd_bundle", Factory: ResourceBundle},
			{TypeName: "alicloud_ecd_command", Factory: ResourceCommand},
			{TypeName: "alicloud_ecd_custom_property", Factory: ResourceCustomProperty},
			{TypeName: "alicloud_ecd_desktop", Factory: ResourceDesktop},
			{TypeName: "alicloud_ecd_desktop_group", Factory: ResourceDesktopGroup},
			{TypeName: "alicloud_ecd_image", Factory: ResourceImage},
			{TypeName: "alicloud_ecd_nas_file_system", Factory: ResourceNasFileSystem},
			{TypeName: "alicloud_ecd_network_package", Factory: ResourceNetworkPackage},
			{TypeName: "alicloud_ecd_policy_group", Factory: ResourcePolicyGroup},
			{TypeName: "alicloud_ecd_ram_directory", Factory: ResourceRamDirectory},
			{TypeName: "alicloud_ecd_simple_office_site", Factory: ResourceSimpleOfficeSite},
			{TypeName: "alicloud_ecd_snapshot", Factory: ResourceSnapshot},
			{TypeName: "alicloud_ecd_user", Factory: ResourceUser},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_ecd_ad_connector_directories", Factory: DataSourceAdConnectorDirectories},
			{TypeName: "alicloud_ecd_ad_connector_office_sites", Factory: DataSourceAdConnectorOfficeSites},
			{TypeName: "alicloud_ecd_bundles", Factory: DataSourceBundles},
			{TypeName: "alicloud_ecd_commands", Factory: DataSourceCommands},
			{TypeName: "alicloud_ecd_custom_properties", Factory: DataSourceCustomProperties},
			{TypeName: "alicloud_ecd_desktop_groups", Factory: DataSourceDesktopGroups},
			{TypeName: "alicloud_ecd_desktop_types", Factory: DataSourceDesktopTypes},
			{TypeName: "alicloud_ecd_desktops", Factory: DataSourceDesktops},
			{TypeName: "alicloud_ecd_images", Factory: DataSourceImages},
			{TypeName: "alicloud_ecd_nas_file_systems", Factory: DataSourceNasFileSystems},
			{TypeName: "alicloud_ecd_network_packages", Factory: DataSourceNetworkPackages},
			{TypeName: "alicloud_ecd_policy_groups", Factory: DataSourcePolicyGroups},
			{TypeName: "alicloud_ecd_ram_directories", Factory: DataSourceRamDirectories},
			{TypeName: "alicloud_ecd_simple_office_sites", Factory: DataSourceSimpleOfficeSites},
			{TypeName: "alicloud_ecd_snapshots", Factory: DataSourceSnapshots},
			{TypeName: "alicloud_ecd_users", Factory: DataSourceUsers},
			{TypeName: "alicloud_ecd_zones", Factory: DataSourceZones},
		},
	}
}
