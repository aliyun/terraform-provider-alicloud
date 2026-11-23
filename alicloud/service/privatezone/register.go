package privatezone

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PrivateZone
// (pvtz) service. The product code is "pvtz" (the OpenAPI endpoint key);
// the package name "privatezone" is the Go import path only.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "pvtz",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_pvtz_zone", Factory: ResourceZone},
			{TypeName: "alicloud_pvtz_zone_attachment", Factory: ResourceZoneAttachment},
			{TypeName: "alicloud_pvtz_zone_record", Factory: ResourceZoneRecord},
			{TypeName: "alicloud_pvtz_user_vpc_authorization", Factory: ResourceUserVpcAuthorization},
			{TypeName: "alicloud_pvtz_endpoint", Factory: ResourceEndpoint},
			{TypeName: "alicloud_pvtz_rule", Factory: ResourceRule},
			{TypeName: "alicloud_pvtz_rule_attachment", Factory: ResourceRuleAttachment},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_pvtz_zones", Factory: DataSourceZones},
			{TypeName: "alicloud_pvtz_zone_records", Factory: DataSourceZoneRecords},
			{TypeName: "alicloud_pvtz_service", Factory: DataSourceService},
			{TypeName: "alicloud_pvtz_endpoints", Factory: DataSourceEndpoints},
			{TypeName: "alicloud_pvtz_resolver_zones", Factory: DataSourceResolverZones},
			{TypeName: "alicloud_pvtz_rules", Factory: DataSourceRules},
		},
		Resources: []conns.Resource{
			// alicloud_pvtz_inbound_endpoint is a Framework-native resource (the platform v2
			// policy for new types). It lives alongside the SDKv2 resources in the same
			// ServicePackage; fwadapt.WrapResource adapts it into the provider.
			{TypeName: "alicloud_pvtz_inbound_endpoint", Factory: NewInboundEndpointResource},
		},
	}
}
