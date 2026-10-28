package edas

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

// ServicePackage registers the edas service's SDK v2 resources and data sources.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Edas",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_edas_application", Factory: resourceAlicloudEdasApplication},
			{TypeName: "alicloud_edas_deploy_group", Factory: resourceAlicloudEdasDeployGroup},
			{TypeName: "alicloud_edas_application_scale", Factory: resourceAlicloudEdasApplicationScale},
			{TypeName: "alicloud_edas_slb_attachment", Factory: resourceAlicloudEdasSlbAttachment},
			{TypeName: "alicloud_edas_cluster", Factory: resourceAlicloudEdasCluster},
			{TypeName: "alicloud_edas_instance_cluster_attachment", Factory: resourceAlicloudEdasInstanceClusterAttachment},
			{TypeName: "alicloud_edas_application_deployment", Factory: resourceAlicloudEdasApplicationPackageAttachment},
			{TypeName: "alicloud_edas_k8s_cluster", Factory: resourceAlicloudEdasK8sCluster},
			{TypeName: "alicloud_edas_k8s_application", Factory: resourceAlicloudEdasK8sApplication},
			{TypeName: "alicloud_edas_k8s_slb_attachment", Factory: resourceAlicloudEdasK8sSlbAttachment},
			{TypeName: "alicloud_edas_namespace", Factory: resourceAlicloudEdasNamespace},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_edas_service", Factory: dataSourceAlicloudEdasService},
			{TypeName: "alicloud_edas_applications", Factory: dataSourceAlicloudEdasApplications},
			{TypeName: "alicloud_edas_deploy_groups", Factory: dataSourceAlicloudEdasDeployGroups},
			{TypeName: "alicloud_edas_clusters", Factory: dataSourceAlicloudEdasClusters},
			{TypeName: "alicloud_edas_namespaces", Factory: dataSourceAlicloudEdasNamespaces},
		},
	}
}
