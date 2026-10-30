package alikafka

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Alikafka service.
// The twelve public Resource/Data Source types migrated into this product
// subpackage are declared here and wired through provider/registry; the
// matching entries in the root provider.go map literals are removed so each
// type is declared exactly once.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "alikafka",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_alikafka_consumer_group",
				Factory:  ResourceConsumerGroup,
			},
			{
				TypeName: "alicloud_alikafka_instance",
				Factory:  ResourceInstance,
			},
			{
				TypeName: "alicloud_alikafka_instance_allowed_ip_attachment",
				Factory:  ResourceInstanceAllowedIpAttachment,
			},
			{
				TypeName: "alicloud_alikafka_sasl_acl",
				Factory:  ResourceSaslAcl,
			},
			{
				TypeName: "alicloud_alikafka_sasl_user",
				Factory:  ResourceSaslUser,
			},
			{
				TypeName: "alicloud_alikafka_scheduled_scaling_rule",
				Factory:  ResourceScheduledScalingRule,
			},
			{
				TypeName: "alicloud_alikafka_topic",
				Factory:  ResourceTopic,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_alikafka_consumer_groups",
				Factory:  DataSourceConsumerGroups,
			},
			{
				TypeName: "alicloud_alikafka_instances",
				Factory:  DataSourceInstances,
			},
			{
				TypeName: "alicloud_alikafka_sasl_acls",
				Factory:  DataSourceSaslAcls,
			},
			{
				TypeName: "alicloud_alikafka_sasl_users",
				Factory:  DataSourceSaslUsers,
			},
			{
				TypeName: "alicloud_alikafka_topics",
				Factory:  DataSourceTopics,
			},
		},
	}
}
