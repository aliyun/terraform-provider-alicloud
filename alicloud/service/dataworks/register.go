package dataworks

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the DataWorks service.
// The public type name alicloud_data_works_service and its historical alias
// alicloud_dataworks_service share the same factory (DataSourceService).
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "dataworks-public",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_data_works_data_source_shared_rule",
				Factory:  ResourceDataSourceSharedRule,
			},
			{
				TypeName: "alicloud_data_works_data_source",
				Factory:  ResourceDataSource,
			},
			{
				TypeName: "alicloud_data_works_di_alarm_rule",
				Factory:  ResourceDiAlarmRule,
			},
			{
				TypeName: "alicloud_data_works_di_job",
				Factory:  ResourceDiJob,
			},
			{
				TypeName: "alicloud_data_works_dw_resource_group",
				Factory:  ResourceDwResourceGroup,
			},
			{
				TypeName: "alicloud_data_works_folder",
				Factory:  ResourceFolder,
			},
			{
				TypeName: "alicloud_data_works_network",
				Factory:  ResourceNetwork,
			},
			{
				TypeName: "alicloud_data_works_project_member",
				Factory:  ResourceProjectMember,
			},
			{
				TypeName: "alicloud_data_works_project",
				Factory:  ResourceProject,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_data_works_folders",
				Factory:  DataSourceFolders,
			},
			{
				TypeName: "alicloud_data_works_service",
				Factory:  DataSourceService,
			},
			{
				TypeName: "alicloud_dataworks_service",
				Factory:  DataSourceService,
			},
		},
	}
}
