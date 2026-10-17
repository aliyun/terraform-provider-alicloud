package ecs

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the ECS service.
// Saving Plan resource and data source are Framework-native implementations
// registered here; the rest of ECS remains SDKv2 in the flat alicloud package.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Ecs",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_ecs_saving_plan",
				Factory:  NewSavingPlanResource,
			},
		},
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_ecs_saving_plans",
				Factory:  NewSavingPlansDataSource,
			},
		},
	}
}
