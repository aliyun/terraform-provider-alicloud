package ecs

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// Exported test-only wrappers for unexported data-source constructors that
// the external ecs_test package's instance acceptance helper
// (testAccInstancePrivatePoolId) calls by bare name. Under the flat root
// layout these were same-package calls; in the split package the
// constructors are unexported in package ecs and thus unreachable
// cross-package, so these wrappers expose them to the test binary (visible
// to both package ecs and package ecs_test tests). They forward verbatim to
// the unexported constructors, preserving behavior.
func DataSourceEcsElasticityAssurancesForTest() *schema.Resource {
	return dataSourceEcsElasticityAssurances()
}

func DataSourceAliCloudInstanceTypesForTest() *schema.Resource {
	return dataSourceAliCloudInstanceTypes()
}

// The migrated security_group_rule TestUnit (TestUnitECSSecurityGroupRuleInvalidImportID,
// in package ecs_test) calls these unexported resource/CRUD/helpers by bare name:
// the resource constructor resourceAliyunSecurityGroupRule, its Read CRUD
// resourceAliyunSecurityGroupRuleRead, and the parseSecurityRuleId helper. They
// are unexported in package ecs and unreachable cross-package, so these
// wrappers expose them to the test binary. Forward verbatim, preserving behavior.
func ResourceAliyunSecurityGroupRuleForTest() *schema.Resource {
	return resourceAliyunSecurityGroupRule()
}

func ResourceAliyunSecurityGroupRuleReadForTest(d *schema.ResourceData, meta interface{}) error {
	return resourceAliyunSecurityGroupRuleRead(d, meta)
}

func ParseSecurityRuleIdForTest(d *schema.ResourceData, meta interface{}, index int) string {
	return parseSecurityRuleId(d, meta, index)
}
