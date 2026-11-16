package alicloud

// This file keeps the root-package copy of DBInstanceStatusCatcher, whose
// definition moved with the RDS product migration into alicloud/service/rds.
// The RDS sub-package has its own local copy and does not import the root
// package, so both copies coexist without a cycle.

// DBInstanceStatusCatcher moved with alicloud/service_alicloud_rds.go ->
// alicloud/service/rds/service.go. It is referenced by the adb and polardb
// service layers' DB-instance state-refresh catchers.
var DBInstanceStatusCatcher = Catcher{Reason: "OperationDenied.DBInstanceStatus", RetryCount: 60, RetryWaitSeconds: 5}
