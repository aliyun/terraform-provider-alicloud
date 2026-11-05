package ros

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the ROS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "ROS",
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_ros_change_sets",
				Factory:  dataSourceAlicloudRosChangeSets,
			},
			{
				TypeName: "alicloud_ros_regions",
				Factory:  dataSourceAlicloudRosRegions,
			},
			{
				TypeName: "alicloud_ros_stack_groups",
				Factory:  dataSourceAlicloudRosStackGroups,
			},
			{
				TypeName: "alicloud_ros_stack_instances",
				Factory:  dataSourceAlicloudRosStackInstances,
			},
			{
				TypeName: "alicloud_ros_stacks",
				Factory:  dataSourceAlicloudRosStacks,
			},
			{
				TypeName: "alicloud_ros_template_scratches",
				Factory:  dataSourceAlicloudRosTemplateScratches,
			},
			{
				TypeName: "alicloud_ros_templates",
				Factory:  dataSourceAlicloudRosTemplates,
			},
		},
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_ros_change_set",
				Factory:  resourceAlicloudRosChangeSet,
			},
			{
				TypeName: "alicloud_ros_stack",
				Factory:  resourceAlicloudRosStack,
			},
			{
				TypeName: "alicloud_ros_stack_group",
				Factory:  resourceAliCloudRosStackGroup,
			},
			{
				TypeName: "alicloud_ros_stack_instance",
				Factory:  resourceAlicloudRosStackInstance,
			},
			{
				TypeName: "alicloud_ros_stack_instances",
				Factory:  resourceAlicloudRosStackInstances,
			},
			{
				TypeName: "alicloud_ros_template",
				Factory:  resourceAlicloudRosTemplate,
			},
			{
				TypeName: "alicloud_ros_template_scratch",
				Factory:  resourceAlicloudRosTemplateScratch,
			},
		},
	}
}
