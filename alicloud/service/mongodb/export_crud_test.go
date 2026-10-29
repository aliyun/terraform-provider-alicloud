package mongodb

// export_crud_test.go is an internal test file (package mongodb) that re-exports
// the package's unexported resource CRUD functions as exported identifiers so
// that the external test package (mongodb_test) can reference them. This is the
// canonical Go "export_test.go" pattern.
//
// Why these are needed: mongodb_test is an EXTERNAL test package (mongodb_test,
// not mongodb) and therefore cannot reference unexported identifiers of package
// mongodb directly. The TestUnit tests in mongodb_test invoke the resource CRUD
// functions to exercise creation/update/read/delete paths without a real cloud
// round-trip. Keeping the CRUD functions unexported preserves the production
// API surface and the schema registration (which references them by their
// unexported names); this file bridges visibility only for tests.
//
// This file imports NOTHING, so it introduces no import cycle: the test-augmented
// package mongodb gains no new imports, and mongodb_test -> acctest -> alicloud
// -> registry -> mongodb (production) is not a cycle because mongodb_test is a
// distinct package path from mongodb.

var (
	// alicloud_mongodb_account
	ResourceAliCloudMongodbAccountCreate = resourceAliCloudMongodbAccountCreate
	ResourceAliCloudMongodbAccountRead   = resourceAliCloudMongodbAccountRead
	ResourceAliCloudMongodbAccountUpdate = resourceAliCloudMongodbAccountUpdate
	ResourceAliCloudMongodbAccountDelete = resourceAliCloudMongodbAccountDelete

	// alicloud_mongodb_audit_policy
	ResourceAliCloudMongodbAuditPolicyCreate = resourceAliCloudMongodbAuditPolicyCreate
	ResourceAliCloudMongodbAuditPolicyRead   = resourceAliCloudMongodbAuditPolicyRead
	ResourceAliCloudMongodbAuditPolicyUpdate = resourceAliCloudMongodbAuditPolicyUpdate
	ResourceAliCloudMongodbAuditPolicyDelete = resourceAliCloudMongodbAuditPolicyDelete

	// alicloud_mongodb_sharding_audit_filter
	ResourceAliCloudMongodbShardingAuditFilterCreate = resourceAliCloudMongodbShardingAuditFilterCreate
	ResourceAliCloudMongodbShardingAuditFilterRead   = resourceAliCloudMongodbShardingAuditFilterRead
	ResourceAliCloudMongodbShardingAuditFilterUpdate = resourceAliCloudMongodbShardingAuditFilterUpdate
	ResourceAliCloudMongodbShardingAuditFilterDelete = resourceAliCloudMongodbShardingAuditFilterDelete

	// alicloud_mongodb_sharding_network_private_address
	ResourceAliCloudMongodbShardingNetworkPrivateAddressCreate = resourceAliCloudMongodbShardingNetworkPrivateAddressCreate
	ResourceAliCloudMongodbShardingNetworkPrivateAddressRead   = resourceAliCloudMongodbShardingNetworkPrivateAddressRead
	ResourceAliCloudMongodbShardingNetworkPrivateAddressUpdate = resourceAliCloudMongodbShardingNetworkPrivateAddressUpdate
	ResourceAliCloudMongodbShardingNetworkPrivateAddressDelete = resourceAliCloudMongodbShardingNetworkPrivateAddressDelete

	// alicloud_mongodb_serverless_instance
	ResourceAlicloudMongodbServerlessInstanceCreate = resourceAlicloudMongodbServerlessInstanceCreate
	ResourceAlicloudMongodbServerlessInstanceRead   = resourceAlicloudMongodbServerlessInstanceRead
	ResourceAlicloudMongodbServerlessInstanceUpdate = resourceAlicloudMongodbServerlessInstanceUpdate
	ResourceAlicloudMongodbServerlessInstanceDelete = resourceAlicloudMongodbServerlessInstanceDelete

	// alicloud_mongodb_sharding_network_public_address (no Update function in production)
	ResourceAlicloudMongodbShardingNetworkPublicAddressCreate = resourceAlicloudMongodbShardingNetworkPublicAddressCreate
	ResourceAlicloudMongodbShardingNetworkPublicAddressRead   = resourceAlicloudMongodbShardingNetworkPublicAddressRead
	ResourceAlicloudMongodbShardingNetworkPublicAddressDelete = resourceAlicloudMongodbShardingNetworkPublicAddressDelete
)
