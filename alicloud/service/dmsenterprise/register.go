package dmsenterprise

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage declares the resources and data sources owned by this package.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "dms-enterprise",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_dms_enterprise_instance", Factory: resourceAlicloudDmsEnterpriseInstance},
			{TypeName: "alicloud_dms_enterprise_proxy", Factory: resourceAlicloudDmsEnterpriseProxy},
			{TypeName: "alicloud_dms_enterprise_logic_database", Factory: resourceAlicloudDmsEnterpriseLogicDatabase},
			{TypeName: "alicloud_dms_enterprise_authority_template", Factory: resourceAliCloudDMSEnterpriseAuthorityTemplate},
			{TypeName: "alicloud_dms_enterprise_workspace", Factory: resourceAliCloudDmsEnterpriseWorkspace},
			{TypeName: "alicloud_dms_enterprise_user", Factory: resourceAlicloudDmsEnterpriseUser},
			{TypeName: "alicloud_dms_enterprise_proxy_access", Factory: resourceAlicloudDmsEnterpriseProxyAccess},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_dms_enterprise_proxy_accesses", Factory: dataSourceAlicloudDmsEnterpriseProxyAccesses},
			{TypeName: "alicloud_dms_enterprise_proxies", Factory: dataSourceAlicloudDmsEnterpriseProxies},
			{TypeName: "alicloud_dms_user_tenants", Factory: dataSourceAlicloudDmsUserTenants},
			{TypeName: "alicloud_dms_enterprise_instances", Factory: dataSourceAlicloudDmsEnterpriseInstances},
			{TypeName: "alicloud_dms_enterprise_logic_databases", Factory: dataSourceAlicloudDmsEnterpriseLogicDatabases},
			{TypeName: "alicloud_dms_enterprise_users", Factory: dataSourceAlicloudDmsEnterpriseUsers},
			{TypeName: "alicloud_dms_enterprise_databases", Factory: dataSourceAlicloudDmsEnterpriseDatabases},
		},
	}
}
