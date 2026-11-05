package cs

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the CS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Cs",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_cs_kubernetes",
				Factory:  ResourceCSKubernetes,
			},
			{
				TypeName: "alicloud_cs_kubernetes_policy_instance",
				Factory:  ResourceCSKubernetesPolicyInstance,
			},
			{
				TypeName: "alicloud_cs_kubernetes_addon",
				Factory:  ResourceCSKubernetesAddon,
			},
			{
				TypeName: "alicloud_cs_managed_kubernetes",
				Factory:  ResourceCSManagedKubernetes,
			},
			{
				TypeName: "alicloud_cs_edge_kubernetes",
				Factory:  ResourceCSEdgeKubernetes,
			},
			{
				TypeName: "alicloud_cs_serverless_kubernetes",
				Factory:  ResourceCSServerlessKubernetes,
			},
			{
				TypeName: "alicloud_cs_kubernetes_autoscaler",
				Factory:  ResourceCSKubernetesAutoscaler,
			},
			{
				TypeName: "alicloud_cs_kubernetes_node_pool",
				Factory:  ResourceAckNodepool,
			},
			{
				TypeName: "alicloud_cs_kubernetes_permissions",
				Factory:  ResourceCSKubernetesPermissions,
			},
			{
				TypeName: "alicloud_cs_autoscaling_config",
				Factory:  ResourceCSAutoscalingConfig,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_ack_service",
				Factory:  DataSourceAckService,
			},
			{
				TypeName: "alicloud_cs_clusters",
				Factory:  DataSourceAckClusters,
			},
			{
				TypeName: "alicloud_cs_kubernetes_node_pools",
				Factory:  DataSourceAckNodepools,
			},
			{
				TypeName: "alicloud_cs_kubernetes_clusters",
				Factory:  DataSourceCSKubernetesClusters,
			},
			{
				TypeName: "alicloud_cs_managed_kubernetes_clusters",
				Factory:  DataSourceCSManagerKubernetesClusters,
			},
			{
				TypeName: "alicloud_cs_edge_kubernetes_clusters",
				Factory:  DataSourceCSEdgeKubernetesClusters,
			},
			{
				TypeName: "alicloud_cs_serverless_kubernetes_clusters",
				Factory:  DataSourceCSServerlessKubernetesClusters,
			},
			{
				TypeName: "alicloud_cs_kubernetes_permissions",
				Factory:  DataSourceCSKubernetesPermissions,
			},
			{
				TypeName: "alicloud_cs_kubernetes_addons",
				Factory:  DataSourceCSKubernetesAddons,
			},
			{
				TypeName: "alicloud_cs_kubernetes_version",
				Factory:  DataSourceCSKubernetesVersion,
			},
			{
				TypeName: "alicloud_cs_kubernetes_addon_metadata",
				Factory:  DataSourceCSKubernetesAddonMetadata,
			},
			{
				TypeName: "alicloud_cs_cluster_credential",
				Factory:  DataSourceCSClusterCredential,
			},
		},
		EphemeralResources: []conns.EphemeralResource{
			{
				TypeName: "alicloud_cs_cluster_credential",
				Factory:  NewClusterCredentialEphemeralResource,
			},
		},
	}
}
