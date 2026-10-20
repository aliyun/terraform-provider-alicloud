package acctest

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// in most cases, the TestCheckFunc list of dataSource test case is repeated，so we make an abstract in
// order to reduce redundant code.
// DataSourceAttr has 3 field ,incloud ResourceId  ExistMapFunc FakeMapFunc, every dataSource test can use only one
type DataSourceAttr struct {
	// IDRefreshName, like "data.alicloud_dns_records.record"
	ResourceId string

	// get existMap function
	ExistMapFunc func(rand int) map[string]string

	// get fakeMap function
	FakeMapFunc func(rand int) map[string]string
}

// get exist and empty ResourceAttrMapUpdate function
func (dsa *DataSourceAttr) checkDataSourceAttr(rand int) (exist, empty ResourceAttrMapUpdate) {
	exist = ResourceAttrInit(dsa.ResourceId, dsa.ExistMapFunc(rand)).ResourceAttrMapUpdateSet()
	empty = ResourceAttrInit(dsa.ResourceId, dsa.FakeMapFunc(rand)).ResourceAttrMapUpdateSet()
	return
}

// according to configs generate step list and execute the test
// lintignore: AT001
func (dsa *DataSourceAttr) DataSourceTestCheck(t *testing.T, rand int, configs ...DataSourceTestAccConfig) {
	var steps []resource.TestStep
	for _, conf := range configs {
		steps = append(steps, conf.buildDataSourceSteps(t, dsa, rand)...)
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			PreCheck(t)
		},
		ProviderFactories: SDKProviderFactories,
		Steps:             steps,
	})
}

// according to configs generate step list and execute the test with preCheck
// lintignore: AT001
func (dsa *DataSourceAttr) DataSourceTestCheckWithPreCheck(t *testing.T, rand int, preCheck func(), configs ...DataSourceTestAccConfig) {
	var steps []resource.TestStep
	for _, conf := range configs {
		steps = append(steps, conf.buildDataSourceSteps(t, dsa, rand)...)
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          preCheck,
		ProviderFactories: SDKProviderFactories,
		Steps:             steps,
	})
}

// per schema attribute test config
type DataSourceTestAccConfig struct {
	// be equal to testCase config string,but the result has only one record
	ExistConfig string

	// if the DataSourceAttr.ExistMapFunc returned map value not match we want, ExistChangMap can alter checkMap for ExistConfig
	ExistChangMap map[string]string

	// be equal to testCase config string,but the result is empty
	FakeConfig string

	// if the DataSourceAttr.FakeMapFunc returned map value not match we want, FakeChangMap can alter checkMap for FakeConfig
	FakeChangMap map[string]string
}

// build test cases for each attribute
func (conf *DataSourceTestAccConfig) buildDataSourceSteps(t *testing.T, info *DataSourceAttr, rand int) []resource.TestStep {
	testAccCheckExist, testAccCheckEmpty := info.checkDataSourceAttr(rand)
	var steps []resource.TestStep
	if conf.ExistConfig != "" {
		step := resource.TestStep{
			Config: conf.ExistConfig,
			Check: resource.ComposeTestCheckFunc(
				testAccCheckExist(conf.ExistChangMap),
			),
		}
		steps = append(steps, step)
	}
	if conf.FakeConfig != "" {
		step := resource.TestStep{
			Config: conf.FakeConfig,
			Check: resource.ComposeTestCheckFunc(
				testAccCheckEmpty(conf.FakeChangMap),
			),
		}
		steps = append(steps, step)
	}
	return steps
}
