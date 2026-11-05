package bastionhost

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the bastionhost service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Yundun-bastionhost",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_bastionhost_instance", Factory: ResourceInstance},
			{TypeName: "alicloud_bastionhost_user_group", Factory: ResourceUserGroup},
			{TypeName: "alicloud_bastionhost_user", Factory: ResourceUser},
			{TypeName: "alicloud_bastionhost_host_group", Factory: ResourceHostGroup},
			{TypeName: "alicloud_bastionhost_user_attachment", Factory: ResourceUserAttachment},
			{TypeName: "alicloud_bastionhost_host", Factory: ResourceHost},
			{TypeName: "alicloud_bastionhost_host_account", Factory: ResourceHostAccount},
			{TypeName: "alicloud_bastionhost_host_attachment", Factory: ResourceHostAttachment},
			{TypeName: "alicloud_bastionhost_host_account_user_group_attachment", Factory: ResourceHostAccountUserGroupAttachment},
			{TypeName: "alicloud_bastionhost_host_account_user_attachment", Factory: ResourceHostAccountUserAttachment},
			{TypeName: "alicloud_bastionhost_host_group_account_user_attachment", Factory: ResourceHostGroupAccountUserAttachment},
			{TypeName: "alicloud_bastionhost_host_group_account_user_group_attachment", Factory: ResourceHostGroupAccountUserGroupAttachment},
			{TypeName: "alicloud_bastionhost_host_share_key", Factory: ResourceHostShareKey},
			{TypeName: "alicloud_bastionhost_host_account_share_key_attachment", Factory: ResourceHostAccountShareKeyAttachment},
			// yundun alias
			{TypeName: "alicloud_yundun_bastionhost_instance", Factory: ResourceInstance},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_bastionhost_instances", Factory: DataSourceInstances},
			{TypeName: "alicloud_bastionhost_user_groups", Factory: DataSourceUserGroups},
			{TypeName: "alicloud_bastionhost_users", Factory: DataSourceUsers},
			{TypeName: "alicloud_bastionhost_host_groups", Factory: DataSourceHostGroups},
			{TypeName: "alicloud_bastionhost_hosts", Factory: DataSourceHosts},
			{TypeName: "alicloud_bastionhost_host_accounts", Factory: DataSourceHostAccounts},
			{TypeName: "alicloud_bastionhost_host_share_keys", Factory: DataSourceHostShareKeys},
			// yundun alias
			{TypeName: "alicloud_yundun_bastionhost_instances", Factory: DataSourceInstances},
		},
	}
}
