package eflo

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the eflo service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "eflo",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_eflo_cluster",
				Factory:  ResourceCluster,
			},
			{
				TypeName: "alicloud_eflo_er",
				Factory:  ResourceEr,
			},
			{
				TypeName: "alicloud_eflo_experiment_plan",
				Factory:  ResourceExperimentPlan,
			},
			{
				TypeName: "alicloud_eflo_experiment_plan_template",
				Factory:  ResourceExperimentPlanTemplate,
			},
			{
				TypeName: "alicloud_eflo_hyper_node",
				Factory:  ResourceHyperNode,
			},
			{
				TypeName: "alicloud_eflo_invocation",
				Factory:  ResourceInvocation,
			},
			{
				TypeName: "alicloud_eflo_node",
				Factory:  ResourceNode,
			},
			{
				TypeName: "alicloud_eflo_node_group",
				Factory:  ResourceNodeGroup,
			},
			{
				TypeName: "alicloud_eflo_node_group_attachment",
				Factory:  ResourceNodeGroupAttachment,
			},
			{
				TypeName: "alicloud_eflo_resource",
				Factory:  ResourceResource,
			},
			{
				TypeName: "alicloud_eflo_subnet",
				Factory:  ResourceSubnet,
			},
			{
				TypeName: "alicloud_eflo_vpd",
				Factory:  ResourceVpd,
			},
			{
				TypeName: "alicloud_eflo_vpd_grant_rule",
				Factory:  ResourceVpdGrantRule,
			},
			{
				TypeName: "alicloud_eflo_vsc",
				Factory:  ResourceVsc,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_eflo_subnets",
				Factory:  DataSourceSubnets,
			},
			{
				TypeName: "alicloud_eflo_vpds",
				Factory:  DataSourceVpds,
			},
		},
	}
}
