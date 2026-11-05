package amqpopen

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the amqp-open service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "amqp-open",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_amqp_instance",
				Factory:  ResourceAmqpInstance,
			},
			{
				TypeName: "alicloud_amqp_binding",
				Factory:  ResourceAmqpBinding,
			},
			{
				TypeName: "alicloud_amqp_exchange",
				Factory:  ResourceAmqpExchange,
			},
			{
				TypeName: "alicloud_amqp_queue",
				Factory:  ResourceAmqpQueue,
			},
			{
				TypeName: "alicloud_amqp_static_account",
				Factory:  ResourceAmqpStaticAccount,
			},
			{
				TypeName: "alicloud_amqp_open_source_account",
				Factory:  ResourceAmqpOpenSourceAccount,
			},
			{
				TypeName: "alicloud_amqp_open_source_permission",
				Factory:  ResourceAmqpOpenSourcePermission,
			},
			{
				TypeName: "alicloud_amqp_virtual_host",
				Factory:  ResourceAmqpVirtualHost,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_amqp_instances",
				Factory:  DataSourceAmqpInstances,
			},
			{
				TypeName: "alicloud_amqp_bindings",
				Factory:  DataSourceAmqpBindings,
			},
			{
				TypeName: "alicloud_amqp_exchanges",
				Factory:  DataSourceAmqpExchanges,
			},
			{
				TypeName: "alicloud_amqp_queues",
				Factory:  DataSourceAmqpQueues,
			},
			{
				TypeName: "alicloud_amqp_static_accounts",
				Factory:  DataSourceAmqpStaticAccounts,
			},
			{
				TypeName: "alicloud_amqp_virtual_hosts",
				Factory:  DataSourceAmqpVirtualHosts,
			},
			{
				TypeName: "alicloud_amqp_open_source_accounts",
				Factory:  DataSourceAmqpOpenSourceAccounts,
			},
			{
				TypeName: "alicloud_amqp_open_source_permissions",
				Factory:  DataSourceAmqpOpenSourcePermissions,
			},
		},
	}
}
