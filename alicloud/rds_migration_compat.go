package alicloud

import "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"

// This file restores root-package symbols whose definitions moved with the
// RDS product migration into alicloud/service/rds but are still referenced by
// other (non-RDS) root files. Per stock-method-disposition rule #4 the root
// keeps a behavior-preserving copy; the bodies are verbatim from the moved
// files. The RDS sub-package has its own local copies and does not import the
// root package, so both copies coexist without a cycle.

// parameterToHash moved with alicloud/resource_alicloud_db_instance.go ->
// alicloud/service/rds/resource_alicloud_db_instance.go. It is the SchemaSet
// hash func for the shared "parameters" set used by the polardb, mongodb,
// gpdb and db instance resources.
func parameterToHash(v interface{}) int {
	m := v.(map[string]interface{})
	return helper.Hashcode(m["name"].(string) + "|" + m["value"].(string))
}

// DBInstanceStatusCatcher moved with alicloud/service_alicloud_rds.go ->
// alicloud/service/rds/service.go. It is referenced by the adb and polardb
// service layers' DB-instance state-refresh catchers.
var DBInstanceStatusCatcher = Catcher{Reason: "OperationDenied.DBInstanceStatus", RetryCount: 60, RetryWaitSeconds: 5}
