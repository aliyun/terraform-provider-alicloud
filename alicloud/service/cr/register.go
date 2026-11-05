package cr

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the CR (Container
// Registry) service. It is the single registration entry point consumed by
// provider/registry; the SDKv2 resources and data sources migrated here from
// the root alicloud package keep their original factories and behavior.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Cr",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_cr_namespace", Factory: resourceAlicloudCRNamespace},
			{TypeName: "alicloud_cr_repo", Factory: resourceAlicloudCRRepo},
			{TypeName: "alicloud_cr_ee_instance", Factory: resourceAliCloudCrInstance},
			{TypeName: "alicloud_cr_ee_namespace", Factory: resourceAliCloudCrEENamespace},
			{TypeName: "alicloud_cr_ee_repo", Factory: resourceAliCloudCrEERepo},
			{TypeName: "alicloud_cr_ee_sync_rule", Factory: resourceAliCloudCrRepoSyncRule},
			{TypeName: "alicloud_cr_endpoint_acl_policy", Factory: resourceAlicloudCrEndpointAclPolicy},
			{TypeName: "alicloud_cr_chart_namespace", Factory: resourceAlicloudCrChartNamespace},
			{TypeName: "alicloud_cr_chart_repository", Factory: resourceAlicloudCrChartRepository},
			{TypeName: "alicloud_cr_chain", Factory: resourceAlicloudCrChain},
			{TypeName: "alicloud_cr_vpc_endpoint_linked_vpc", Factory: resourceAlicloudCrVpcEndpointLinkedVpc},
			{TypeName: "alicloud_cr_internet_endpoint", Factory: resourceAlicloudCrInternetEndpoint},
			{TypeName: "alicloud_cr_instance_customized_domain", Factory: resourceAliCloudCrInstanceCustomizedDomain},
			{TypeName: "alicloud_cr_artifact_subscription_rule", Factory: resourceAliCloudCrArtifactSubscriptionRule},
			{TypeName: "alicloud_cr_artifact_lifecycle_rule", Factory: resourceAliCloudCrArtifactLifecycleRule},
			{TypeName: "alicloud_cr_storage_domain_routing_rule", Factory: resourceAliCloudCrStorageDomainRoutingRule},
			{TypeName: "alicloud_cr_scan_rule", Factory: resourceAliCloudCrScanRule},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_cr_namespaces", Factory: dataSourceAlicloudCRNamespaces},
			{TypeName: "alicloud_cr_repos", Factory: dataSourceAlicloudCRRepos},
			{TypeName: "alicloud_cr_ee_instances", Factory: dataSourceAlicloudCrEEInstances},
			{TypeName: "alicloud_cr_ee_namespaces", Factory: dataSourceAlicloudCrEENamespaces},
			{TypeName: "alicloud_cr_ee_repos", Factory: dataSourceAliCloudCrEERepos},
			{TypeName: "alicloud_cr_ee_sync_rules", Factory: dataSourceAlicloudCrEESyncRules},
			{TypeName: "alicloud_cr_service", Factory: dataSourceAlicloudCrService},
			{TypeName: "alicloud_cr_endpoint_acl_policies", Factory: dataSourceAlicloudCrEndpointAclPolicies},
			{TypeName: "alicloud_cr_endpoint_acl_service", Factory: dataSourceAlicloudCrEndpointAclService},
			{TypeName: "alicloud_cr_internet_endpoint", Factory: dataSourceAlicloudCrInternetEndpoint},
			{TypeName: "alicloud_cr_chart_namespaces", Factory: dataSourceAlicloudCrChartNamespaces},
			{TypeName: "alicloud_cr_chart_repositories", Factory: dataSourceAlicloudCrChartRepositories},
			{TypeName: "alicloud_cr_chains", Factory: dataSourceAlicloudCrChains},
			{TypeName: "alicloud_cr_vpc_endpoint_linked_vpcs", Factory: dataSourceAlicloudCrVpcEndpointLinkedVpcs},
			{TypeName: "alicloud_cr_artifact_subscription_rules", Factory: dataSourceAliCloudCrArtifactSubscriptionRules},
			{TypeName: "alicloud_cr_artifact_lifecycle_rules", Factory: dataSourceAliCloudCrArtifactLifecycleRules},
		},
		EphemeralResources: []conns.EphemeralResource{
			{
				TypeName: "alicloud_cr_ee_authorization_token",
				Factory:  NewAuthorizationTokenEphemeralResource,
			},
		},
	}
}
