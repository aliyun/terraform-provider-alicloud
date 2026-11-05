package sae

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Sae service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Sae",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_sae_application",
				Factory:  resourceAliCloudSaeApplication,
			},
			{
				TypeName: "alicloud_sae_application_scaling_rule",
				Factory:  ResourceAlicloudSaeApplicationScalingRule,
			},
			{
				TypeName: "alicloud_sae_config_map",
				Factory:  resourceAlicloudSaeConfigMap,
			},
			{
				TypeName: "alicloud_sae_grey_tag_route",
				Factory:  ResourceAliCloudSaeGreyTagRoute,
			},
			{
				TypeName: "alicloud_sae_ingress",
				Factory:  resourceAlicloudSaeIngress,
			},
			{
				TypeName: "alicloud_sae_load_balancer_internet",
				Factory:  resourceAlicloudSaeLoadBalancerInternet,
			},
			{
				TypeName: "alicloud_sae_load_balancer_intranet",
				Factory:  resourceAlicloudSaeLoadBalancerIntranet,
			},
			{
				TypeName: "alicloud_sae_namespace",
				Factory:  resourceAlicloudSaeNamespace,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_sae_application_groups",
				Factory:  dataSourceAliCloudSaeApplicationGroups,
			},
			{
				TypeName: "alicloud_sae_application_instances",
				Factory:  DataSourceApplicationInstances,
			},
			{
				TypeName: "alicloud_sae_application_scaling_rules",
				Factory:  dataSourceAlicloudSaeApplicationScalingRules,
			},
			{
				TypeName: "alicloud_sae_applications",
				Factory:  dataSourceAlicloudSaeApplications,
			},
			{
				TypeName: "alicloud_sae_config_maps",
				Factory:  dataSourceAlicloudSaeConfigMaps,
			},
			{
				TypeName: "alicloud_sae_grey_tag_routes",
				Factory:  dataSourceAlicloudSaeGreyTagRoutes,
			},
			{
				TypeName: "alicloud_sae_ingresses",
				Factory:  dataSourceAlicloudSaeIngresses,
			},
			{
				TypeName: "alicloud_sae_instance_specifications",
				Factory:  dataSourceAlicloudSaeInstanceSpecifications,
			},
			{
				TypeName: "alicloud_sae_namespaces",
				Factory:  dataSourceAlicloudSaeNamespaces,
			},
			{
				TypeName: "alicloud_sae_service",
				Factory:  dataSourceAlicloudSaeService,
			},
		},
	}
}
