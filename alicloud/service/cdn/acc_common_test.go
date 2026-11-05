package cdn_test

import (
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Shared acceptance-test helpers for the split CDN package.
//
// The generic resource-test core — ResourceCheck/ResourceAttr/ResourceAttrCheck
// (with their Init constructors and CheckResourceExists/CheckResourceDestroy/
// ResourceAttrMapUpdateSet methods), the resourceConfig engine
// (ResourceTestAccConfig/DataSourceTestAccConfigFunc) and the NOSET/CHECKSET/...
// sentinels — is reused directly from alicloud/acctest. Migrated cdn_test
// fixtures therefore call acctest.ResourceCheckInit/ResourceAttrInit/
// ResourceAttrCheckInit/ResourceTestAccConfig/DataSourceTestAccConfigFunc and
// the acctest.* sentinel constants instead of keeping local copies here.
//
// What remains local is the data-source acceptance runner. The public
// alicloud/acctest.DataSourceTestCheck / DataSourceTestCheckWithPreCheck hardcode
// ProviderFactories: SDKProviderFactories, while the migrated CDN data-source
// fixtures run on ProtoV5ProviderFactories. alicloud/acctest also leaves
// checkDataSourceAttr and buildDataSourceSteps unexported, so they cannot be
// reused from this external cdn_test package. The ProtoV5 runner below,
// together with its DataSourceAttr/DataSourceTestAccConfig support types and
// step-builder methods, is therefore retained as the minimal protocol
// adapter; the underlying attribute checks delegate to the public
// acctest.ResourceAttrInit / ResourceAttrMapUpdateSet so the attr machinery is
// not duplicated.

type dataSourceAttr struct {
	// IDRefreshName, like "data.alicloud_dns_records.record"
	resourceId string

	// get existMap function
	existMapFunc func(rand int) map[string]string

	// get fakeMap function
	fakeMapFunc func(rand int) map[string]string
}

// get exist and empty resourceAttrMapUpdate function.
// Delegates to the public acctest.ResourceAttrInit/ResourceAttrMapUpdateSet so
// the attribute-check machinery is not duplicated here.
func (dsa *dataSourceAttr) checkDataSourceAttr(rand int) (exist, empty acctest.ResourceAttrMapUpdate) {
	exist = acctest.ResourceAttrInit(dsa.resourceId, dsa.existMapFunc(rand)).ResourceAttrMapUpdateSet()
	empty = acctest.ResourceAttrInit(dsa.resourceId, dsa.fakeMapFunc(rand)).ResourceAttrMapUpdateSet()
	return
}

// according to configs generate step list and execute the test.
// Kept local for ProtoV5ProviderFactories; the public runner uses
// ProviderFactories: SDKProviderFactories and is therefore not substitutable.
// lintignore: AT001
func (dsa *dataSourceAttr) dataSourceTestCheck(t *testing.T, rand int, configs ...dataSourceTestAccConfig) {
	var steps []resource.TestStep
	for _, conf := range configs {
		steps = append(steps, conf.buildDataSourceSteps(t, dsa, rand)...)
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps:                    steps,
	})
}

// according to configs generate step list and execute the test with preCheck.
// Kept local for ProtoV5ProviderFactories; the public runner uses
// ProviderFactories: SDKProviderFactories and is therefore not substitutable.
// lintignore: AT001
func (dsa *dataSourceAttr) dataSourceTestCheckWithPreCheck(t *testing.T, rand int, preCheck func(), configs ...dataSourceTestAccConfig) {
	var steps []resource.TestStep
	for _, conf := range configs {
		steps = append(steps, conf.buildDataSourceSteps(t, dsa, rand)...)
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 preCheck,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps:                    steps,
	})
}

// per schema attribute test config. Kept local because the public
// buildDataSourceSteps is unexported and cannot be reused cross-package.
type dataSourceTestAccConfig struct {
	// be equal to testCase config string,but the result has only one record
	existConfig string

	// if the dataSourceAttr.existMapFunc returned map value not match we want, existChangMap can alter checkMap for existConfig
	existChangMap map[string]string

	// be equal to testCase config string,but the result is empty
	fakeConfig string

	// if the dataSourceAttr.fakeMapFunc returned map value not match we want, fakeChangMap can alter checkMap for fakeConfig
	fakeChangMap map[string]string
}

// build test cases for each attribute. Kept local because the public
// buildDataSourceSteps is unexported and cannot be reused cross-package.
func (conf *dataSourceTestAccConfig) buildDataSourceSteps(t *testing.T, info *dataSourceAttr, rand int) []resource.TestStep {
	testAccCheckExist, testAccCheckEmpty := info.checkDataSourceAttr(rand)
	var steps []resource.TestStep
	if conf.existConfig != "" {
		step := resource.TestStep{
			Config: conf.existConfig,
			Check: resource.ComposeTestCheckFunc(
				testAccCheckExist(conf.existChangMap),
			),
		}
		steps = append(steps, step)
	}
	if conf.fakeConfig != "" {
		step := resource.TestStep{
			Config: conf.fakeConfig,
			Check: resource.ComposeTestCheckFunc(
				testAccCheckEmpty(conf.fakeChangMap),
			),
		}
		steps = append(steps, step)
	}
	return steps
}

// defaultRegionToTest mirrors the root-package var of the same name so that
// test fixtures which embed the region into config strings at test-function
// entry (before PreCheck runs) resolve it from ALICLOUD_REGION, falling back
// to cn-beijing. acctest.PreCheck keeps acctest.DefaultRegionToTest in sync;
// this local var is the entry-time snapshot the migrated fixtures read.
var defaultRegionToTest = func() string {
	if v := os.Getenv("ALICLOUD_REGION"); v != "" {
		return v
	}
	return "cn-beijing"
}()
