package eventbridge

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

// ServicePackage registers the eventbridge service's SDK v2 resources and data sources.
// Source migrated as-is from the root alicloud package (handwritten SDKv2, no
// auto-generation marks); factories are exported so the registry can wire them.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "eventbridge",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_event_bridge_event_source_v2", Factory: ResourceEventBridgeEventSourceV2},
			{TypeName: "alicloud_event_bridge_event_streaming", Factory: ResourceEventBridgeEventStreaming},
			{TypeName: "alicloud_event_bridge_event_bus", Factory: ResourceEventBridgeEventBus},
			{TypeName: "alicloud_event_bridge_rule", Factory: ResourceEventBridgeRule},
			{TypeName: "alicloud_event_bridge_event_source", Factory: ResourceEventBridgeEventSource},
			{TypeName: "alicloud_event_bridge_service_linked_role", Factory: ResourceEventBridgeServiceLinkedRole},
			{TypeName: "alicloud_event_bridge_connection", Factory: ResourceEventBridgeConnection},
			{TypeName: "alicloud_event_bridge_api_destination", Factory: ResourceEventBridgeApiDestination},
			// alicloud_event_bridge_slr is a public alias for alicloud_event_bridge_service_linked_role
			// (same factory), preserved from the root package registration.
			{TypeName: "alicloud_event_bridge_slr", Factory: ResourceEventBridgeServiceLinkedRole},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_event_bridge_service", Factory: DataSourceEventBridgeService},
			{TypeName: "alicloud_event_bridge_event_buses", Factory: DataSourceEventBridgeEventBuses},
			{TypeName: "alicloud_event_bridge_rules", Factory: DataSourceEventBridgeRules},
			{TypeName: "alicloud_event_bridge_event_sources", Factory: DataSourceEventBridgeEventSources},
		},
	}
}
