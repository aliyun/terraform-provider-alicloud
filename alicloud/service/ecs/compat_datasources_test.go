package ecs_test

import (
	ecs "github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
)

// Bare-name aliases for the unexported data-source constructors that the
// migrated instance acceptance helper (testAccInstancePrivatePoolId) calls
// by the original root-layout names. The constructors are unexported in
// package ecs, so compat_datasource_test.go (package ecs) exports
// forwarding wrappers; these vars re-expose them under the bare names the
// test source uses, so it reads verbatim, mirroring compat_enums_test.go /
// compat_helpers_test.go. The elasticity-assurances constructor was renamed
// dataSourceAlicloudEcsElasticityAssurances -> dataSourceEcsElasticityAssurances
// during migration; the wrapper forwards to the renamed function, and this
// var keeps the original bare name the test calls.
var (
	dataSourceAlicloudEcsElasticityAssurances = ecs.DataSourceEcsElasticityAssurancesForTest
	dataSourceAliCloudInstanceTypes           = ecs.DataSourceAliCloudInstanceTypesForTest

	// Unexported security_group_rule resource constructor, its Read CRUD, and
	// the parseSecurityRuleId helper, all called bare by the
	// TestUnitECSSecurityGroupRuleInvalidImportID test. The wrappers live in
	// compat_datasource_test.go (package ecs); these bare-name aliases forward
	// to them so the test reads verbatim.
	resourceAliyunSecurityGroupRule     = ecs.ResourceAliyunSecurityGroupRuleForTest
	resourceAliyunSecurityGroupRuleRead = ecs.ResourceAliyunSecurityGroupRuleReadForTest
	parseSecurityRuleId                 = ecs.ParseSecurityRuleIdForTest
)
