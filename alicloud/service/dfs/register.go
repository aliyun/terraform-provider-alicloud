package dfs

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Dfs service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "DFS",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_dfs_access_group",
				Factory:  ResourceAccessGroup,
			},
			{
				TypeName: "alicloud_dfs_access_rule",
				Factory:  ResourceAccessRule,
			},
			{
				TypeName: "alicloud_dfs_file_system",
				Factory:  ResourceFileSystem,
			},
			{
				TypeName: "alicloud_dfs_mount_point",
				Factory:  ResourceMountPoint,
			},
			{
				TypeName: "alicloud_dfs_vsc_mount_point",
				Factory:  ResourceVscMountPoint,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_dfs_access_groups",
				Factory:  DataSourceAccessGroups,
			},
			{
				TypeName: "alicloud_dfs_access_rules",
				Factory:  DataSourceAccessRules,
			},
			{
				TypeName: "alicloud_dfs_file_systems",
				Factory:  DataSourceFileSystems,
			},
			{
				TypeName: "alicloud_dfs_mount_points",
				Factory:  DataSourceMountPoints,
			},
			{
				TypeName: "alicloud_dfs_zones",
				Factory:  DataSourceZones,
			},
		},
	}
}
