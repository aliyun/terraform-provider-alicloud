package expressconnect

// export_test.go exposes the unexported resource CRUD helpers to the
// external expressconnect_test package. This file is compiled only as part of
// the expressconnect internal test variant and imports nothing (in particular
// not alicloud/acctest), so the internal test variant does not join the
// acctest -> alicloud -> provider/registry -> service/expressconnect import
// cycle. That cycle forbids an internal expressconnect test package here, so
// the ACC/unit tests live in the external expressconnect_test package and
// reach these unexported helpers through the exported aliases below.
var (
	ResourceAliCloudExpressConnectBgpNetworkCreate          = resourceAliCloudExpressConnectBgpNetworkCreate
	ResourceAliCloudExpressConnectBgpNetworkRead            = resourceAliCloudExpressConnectBgpNetworkRead
	ResourceAliCloudExpressConnectBgpNetworkDelete          = resourceAliCloudExpressConnectBgpNetworkDelete
	ResourceAliCloudExpressConnectPhysicalConnectionCreate  = resourceAliCloudExpressConnectPhysicalConnectionCreate
	ResourceAliCloudExpressConnectPhysicalConnectionRead    = resourceAliCloudExpressConnectPhysicalConnectionRead
	ResourceAliCloudExpressConnectPhysicalConnectionUpdate  = resourceAliCloudExpressConnectPhysicalConnectionUpdate
	ResourceAliCloudExpressConnectPhysicalConnectionDelete  = resourceAliCloudExpressConnectPhysicalConnectionDelete
	ResourceAliCloudExpressConnectRouterInterface           = resourceAliCloudExpressConnectRouterInterface
	ResourceAliCloudExpressConnectVirtualBorderRouterCreate = resourceAliCloudExpressConnectVirtualBorderRouterCreate
	ResourceAliCloudExpressConnectVirtualBorderRouterRead   = resourceAliCloudExpressConnectVirtualBorderRouterRead
	ResourceAliCloudExpressConnectVirtualBorderRouterUpdate = resourceAliCloudExpressConnectVirtualBorderRouterUpdate
	ResourceAliCloudExpressConnectVirtualBorderRouterDelete = resourceAliCloudExpressConnectVirtualBorderRouterDelete
	// Legacy alicloud_router_interface resource (separate from
	// alicloud_express_connect_router_interface). Its unit test reaches the
	// unexported constructor and delete CRUD through these aliases.
	ResourceAlicloudRouterInterface       = resourceAlicloudRouterInterface
	ResourceAlicloudRouterInterfaceDelete = resourceAlicloudRouterInterfaceDelete
)
