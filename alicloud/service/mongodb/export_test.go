package mongodb

// export_test.go exposes the unexported resource constructor to the
// external mongodb_test package. The unit-test fixture for backup-policy
// request serialization (newMongoDBBackupPolicyFixtureResource) calls
// resourceAliCloudMongoDBInstance() to obtain the real *schema.Resource so
// it can exercise the Update path against a mock client. After the
// alicloud -> alicloud/service/mongodb migration the tests live in the
// external mongodb_test package and cannot reach the unexported constructor
// directly; this standard Go export_test.go bridge closes that gap.
// No production behavior changes; the production constructor and its call
// sites in provider registration are untouched.
var ResourceAliCloudMongoDBInstance = resourceAliCloudMongoDBInstance
