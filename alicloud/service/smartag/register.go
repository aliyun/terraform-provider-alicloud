package smartag

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Smartag service.
// The public type names, schemas, state, CRUD, ID/import, error, retry, timeout
// and cleanup semantics are preserved verbatim from the flat-layout root
// package; only the package/import wiring changed during the split-layout
// migration.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Smartag",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_cloud_connect_network",
				Factory:  resourceAlicloudCloudConnectNetwork,
			},
			{
				TypeName: "alicloud_cloud_connect_network_attachment",
				Factory:  resourceAlicloudCloudConnectNetworkAttachment,
			},
			{
				TypeName: "alicloud_cloud_connect_network_grant",
				Factory:  resourceAlicloudCloudConnectNetworkGrant,
			},
			{
				TypeName: "alicloud_sag_acl",
				Factory:  resourceAlicloudSagAcl,
			},
			{
				TypeName: "alicloud_sag_acl_rule",
				Factory:  resourceAlicloudSagAclRule,
			},
			{
				TypeName: "alicloud_sag_qos",
				Factory:  resourceAlicloudSagQos,
			},
			{
				TypeName: "alicloud_sag_qos_policy",
				Factory:  resourceAlicloudSagQosPolicy,
			},
			{
				TypeName: "alicloud_sag_qos_car",
				Factory:  resourceAlicloudSagQosCar,
			},
			{
				TypeName: "alicloud_sag_snat_entry",
				Factory:  resourceAlicloudSagSnatEntry,
			},
			{
				TypeName: "alicloud_sag_dnat_entry",
				Factory:  resourceAlicloudSagDnatEntry,
			},
			{
				TypeName: "alicloud_sag_client_user",
				Factory:  resourceAlicloudSagClientUser,
			},
			{
				TypeName: "alicloud_smartag_flow_log",
				Factory:  resourceAlicloudSmartagFlowLog,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_cloud_connect_networks",
				Factory:  dataSourceAlicloudCloudConnectNetworks,
			},
			{
				TypeName: "alicloud_sag_acls",
				Factory:  dataSourceAlicloudSagAcls,
			},
			{
				TypeName: "alicloud_smartag_flow_logs",
				Factory:  dataSourceAlicloudSmartagFlowLogs,
			},
		},
	}
}
