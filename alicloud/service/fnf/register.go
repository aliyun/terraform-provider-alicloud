package fnf

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the fnf service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "fnf",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_fnf_flow",
				Factory:  ResourceFnfFlow,
			},
			{
				TypeName: "alicloud_fnf_execution",
				Factory:  ResourceFnFExecution,
			},
			{
				TypeName: "alicloud_fnf_schedule",
				Factory:  ResourceFnfSchedule,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_fnf_service",
				Factory:  DataSourceFnfService,
			},
			{
				TypeName: "alicloud_fnf_flows",
				Factory:  DataSourceFnfFlows,
			},
			{
				TypeName: "alicloud_fnf_schedules",
				Factory:  DataSourceFnfSchedules,
			},
			{
				TypeName: "alicloud_fnf_executions",
				Factory:  DataSourceFnFExecutions,
			},
		},
	}
}
