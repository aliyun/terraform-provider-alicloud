package brainindustrial

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Brain Industrial service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "brain-industrial",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_brain_industrial_pid_organization",
				Factory:  ResourcePidOrganization,
			},
			{
				TypeName: "alicloud_brain_industrial_pid_project",
				Factory:  ResourcePidProject,
			},
			{
				TypeName: "alicloud_brain_industrial_pid_loop",
				Factory:  ResourcePidLoop,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_brain_industrial_service",
				Factory:  DataSourceService,
			},
			{
				TypeName: "alicloud_brain_industrial_pid_organizations",
				Factory:  DataSourcePidOrganizations,
			},
			{
				TypeName: "alicloud_brain_industrial_pid_projects",
				Factory:  DataSourcePidProjects,
			},
			{
				TypeName: "alicloud_brain_industrial_pid_loops",
				Factory:  DataSourcePidLoops,
			},
		},
	}
}
