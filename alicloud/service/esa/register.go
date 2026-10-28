package esa

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the ESA service.
// 59 public SDKv2 types (51 resources + 8 data sources) migrated from the
// alicloud root package; generation markers and original SDKv2 framework preserved.
// alicloud_esa_site_environment is a native Framework-style resource already
// living in this package; it stays registered here alongside the migrated types.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "ESA",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_esa_site_environment", Factory: ResourceSiteEnvironment},
			{TypeName: "alicloud_esa_cache_reserve_instance", Factory: resourceAliCloudEsaCacheReserveInstance},
			{TypeName: "alicloud_esa_cache_rule", Factory: resourceAliCloudEsaCacheRule},
			{TypeName: "alicloud_esa_certificate", Factory: resourceAliCloudEsaCertificate},
			{TypeName: "alicloud_esa_client_ca_certificate", Factory: resourceAliCloudEsaClientCaCertificate},
			{TypeName: "alicloud_esa_client_certificate", Factory: resourceAliCloudEsaClientCertificate},
			{TypeName: "alicloud_esa_compression_rule", Factory: resourceAliCloudEsaCompressionRule},
			{TypeName: "alicloud_esa_custom_response_code_rule", Factory: resourceAliCloudEsaCustomResponseCodeRule},
			{TypeName: "alicloud_esa_custom_scene_policy", Factory: resourceAliCloudEsaCustomScenePolicy},
			{TypeName: "alicloud_esa_edge_container_app_record", Factory: resourceAliCloudEsaEdgeContainerAppRecord},
			{TypeName: "alicloud_esa_edge_container_app", Factory: resourceAliCloudEsaEdgeContainerApp},
			{TypeName: "alicloud_esa_http_incoming_request_header_modification_rule", Factory: resourceAliCloudEsaHttpIncomingRequestHeaderModificationRule},
			{TypeName: "alicloud_esa_http_incoming_response_header_modification_rule", Factory: resourceAliCloudEsaHttpIncomingResponseHeaderModificationRule},
			{TypeName: "alicloud_esa_http_request_header_modification_rule", Factory: resourceAliCloudEsaHttpRequestHeaderModificationRule},
			{TypeName: "alicloud_esa_http_response_header_modification_rule", Factory: resourceAliCloudEsaHttpResponseHeaderModificationRule},
			{TypeName: "alicloud_esa_https_application_configuration", Factory: resourceAliCloudEsaHttpsApplicationConfiguration},
			{TypeName: "alicloud_esa_https_basic_configuration", Factory: resourceAliCloudEsaHttpsBasicConfiguration},
			{TypeName: "alicloud_esa_image_transform", Factory: resourceAliCloudEsaImageTransform},
			{TypeName: "alicloud_esa_kv_account", Factory: resourceAliCloudEsaKvAccount},
			{TypeName: "alicloud_esa_kv_namespace", Factory: resourceAliCloudEsaKvNamespace},
			{TypeName: "alicloud_esa_kv", Factory: resourceAliCloudEsaKv},
			{TypeName: "alicloud_esa_list", Factory: resourceAliCloudEsaList},
			{TypeName: "alicloud_esa_load_balancer", Factory: resourceAliCloudEsaLoadBalancer},
			{TypeName: "alicloud_esa_network_optimization", Factory: resourceAliCloudEsaNetworkOptimization},
			{TypeName: "alicloud_esa_origin_ca_certificate", Factory: resourceAliCloudEsaOriginCaCertificate},
			{TypeName: "alicloud_esa_origin_client_certificate", Factory: resourceAliCloudEsaOriginClientCertificate},
			{TypeName: "alicloud_esa_origin_pool", Factory: resourceAliCloudEsaOriginPool},
			{TypeName: "alicloud_esa_origin_protection", Factory: resourceAliCloudEsaOriginProtection},
			{TypeName: "alicloud_esa_origin_rule", Factory: resourceAliCloudEsaOriginRule},
			{TypeName: "alicloud_esa_page", Factory: resourceAliCloudEsaPage},
			{TypeName: "alicloud_esa_rate_plan_instance", Factory: resourceAliCloudEsaRatePlanInstance},
			{TypeName: "alicloud_esa_record", Factory: resourceAliCloudEsaRecord},
			{TypeName: "alicloud_esa_redirect_rule", Factory: resourceAliCloudEsaRedirectRule},
			{TypeName: "alicloud_esa_rewrite_url_rule", Factory: resourceAliCloudEsaRewriteUrlRule},
			{TypeName: "alicloud_esa_routine_code_deployment", Factory: resourceAliCloudEsaRoutineCodeDeployment},
			{TypeName: "alicloud_esa_routine_related_record", Factory: resourceAliCloudEsaRoutineRelatedRecord},
			{TypeName: "alicloud_esa_routine", Factory: resourceAliCloudEsaRoutine},
			{TypeName: "alicloud_esa_routine_route", Factory: resourceAliCloudEsaRoutineRoute},
			{TypeName: "alicloud_esa_scheduled_preload_execution", Factory: resourceAliCloudEsaScheduledPreloadExecution},
			{TypeName: "alicloud_esa_scheduled_preload_job", Factory: resourceAliCloudEsaScheduledPreloadJob},
			{TypeName: "alicloud_esa_site_delivery_task", Factory: resourceAliCloudEsaSiteDeliveryTask},
			{TypeName: "alicloud_esa_site_origin_client_certificate", Factory: resourceAliCloudEsaSiteOriginClientCertificate},
			{TypeName: "alicloud_esa_site", Factory: resourceAliCloudEsaSite},
			{TypeName: "alicloud_esa_transport_layer_application", Factory: resourceAliCloudEsaTransportLayerApplication},
			{TypeName: "alicloud_esa_url_observation", Factory: resourceAliCloudEsaUrlObservation},
			{TypeName: "alicloud_esa_version", Factory: resourceAliCloudEsaVersion},
			{TypeName: "alicloud_esa_video_processing", Factory: resourceAliCloudEsaVideoProcessing},
			{TypeName: "alicloud_esa_waf_rule", Factory: resourceAliCloudEsaWafRule},
			{TypeName: "alicloud_esa_waf_ruleset", Factory: resourceAliCloudEsaWafRuleset},
			{TypeName: "alicloud_esa_waiting_room_event", Factory: resourceAliCloudEsaWaitingRoomEvent},
			{TypeName: "alicloud_esa_waiting_room", Factory: resourceAliCloudEsaWaitingRoom},
			{TypeName: "alicloud_esa_waiting_room_rule", Factory: resourceAliCloudEsaWaitingRoomRule},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_esa_cache_reserve_instances", Factory: dataSourceAliCloudEsaCacheReserveInstances},
			{TypeName: "alicloud_esa_https_basic_configurations", Factory: dataSourceAliCloudEsaHttpsBasicConfigurations},
			{TypeName: "alicloud_esa_lists", Factory: dataSourceAliCloudEsaLists},
			{TypeName: "alicloud_esa_network_optimizations", Factory: dataSourceAliCloudEsaNetworkOptimizations},
			{TypeName: "alicloud_esa_origin_rules", Factory: dataSourceAliCloudEsaOriginRules},
			{TypeName: "alicloud_esa_routine_code_versions", Factory: dataSourceAliCloudEsaRoutineCodeVersions},
			{TypeName: "alicloud_esa_sites", Factory: dataSourceAliCloudEsaSites},
			{TypeName: "alicloud_esa_waf_rulesets", Factory: dataSourceAliCloudEsaWafRuleSets},
		},
	}
}
