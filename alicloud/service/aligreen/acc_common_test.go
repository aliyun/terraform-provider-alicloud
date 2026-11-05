package aligreen_test

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
)

// resourceTestAccConfigFunc is the subpackage-internal builder spelling that the
// collection-order checker (scripts/collection-order) discovers via AST: it
// matches the bare resourceTestAccConfigFunc ident with three arguments, like
// the flat alicloud suite's helper and the peer subpackages
// (alicloud/service/ecs, alicloud/service/agentloop). It is a thin wrapper over
// the shared acctest.ResourceTestAccConfig so migrated aligreen tests keep the
// flat suite's cumulative config-map semantics while remaining recognizable to
// the checker. Existing generated tests continue to use
// acctest.ResourceTestAccConfig directly; only the order-coverage cases added
// for the migration use this builder.
func resourceTestAccConfigFunc(resourceId string, name string, configDependence func(name string) string) func(map[string]interface{}) string {
	return acctest.ResourceTestAccConfig(resourceId, name, configDependence)
}
