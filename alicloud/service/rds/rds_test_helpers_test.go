package rds_test

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/rds"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider builds the SDKv2 test provider from the rds ServicePackage
// registration. It mirrors the package-rds internal test Provider()
// (alicloud/service/rds/rds_account_unit_test.go) but lives in the external
// rds_test package so the migrated gomonkey/unit tests (TestUnitAlicloudRdsBackup,
// TestUnitAlicloudRdsInstanceCrossBackupPolicy, TestUnitAlicloudRdsParameterGroup)
// and the readonly-instance schema-validation test can reach it. The external
// package cannot construct the resource factories directly (they are
// unexported) nor see the internal-test Provider(); this wrapper reconstructs
// the same ResourcesMap so p["<type>"] resolves identically.
func Provider() *schema.Provider {
	sp := rds.ServicePackage()
	p := &schema.Provider{ResourcesMap: map[string]*schema.Resource{}}
	for _, r := range sp.SDKResources {
		p.ResourcesMap[r.TypeName] = r.Factory()
	}
	return p
}

// String mirrors the root-package `var String = helper.String` alias
// (alicloud/tags.go) so the migrated gomonkey tests that build tea.SDKError
// mocks can keep using the bare `String(...)` call shape from the baseline
// root package. No behavior change.
var String = helper.String

// mapMerge mirrors the root-package `mapMerge = helper.MapMerge` alias
// (alicloud/common.go) used by the migrated gomonkey unit tests to merge
// mock response maps. No behavior change.
var mapMerge = helper.MapMerge
