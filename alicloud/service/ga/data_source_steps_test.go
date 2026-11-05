package ga_test

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// gaBuildDataSourceSteps mirrors the unexported
// acctest.DataSourceTestAccConfig.buildDataSourceSteps method, which the legacy
// flat-suite data-source tests invoked while still compiling in the root
// alicloud package. Now that these tests live in the external ga_test package,
// the unexported method (and its helper checkDataSourceAttr) are no longer
// reachable, so the step-construction logic is reproduced here verbatim using
// the exported acctest.ResourceAttrInit / ResourceAttrMapUpdateSet helpers.
//
// This preserves the original test structure exactly, including the custom
// PreCheck, ProviderFactories and CheckDestroy wired up by each call site; only
// the []resource.TestStep construction is delegated here, and it is byte-for-byte
// equivalent to the acctest implementation.
func gaBuildDataSourceSteps(conf acctest.DataSourceTestAccConfig, info *acctest.DataSourceAttr, rand int) []resource.TestStep {
	existCheck := acctest.ResourceAttrInit(info.ResourceId, info.ExistMapFunc(rand)).ResourceAttrMapUpdateSet()
	emptyCheck := acctest.ResourceAttrInit(info.ResourceId, info.FakeMapFunc(rand)).ResourceAttrMapUpdateSet()
	var steps []resource.TestStep
	if conf.ExistConfig != "" {
		steps = append(steps, resource.TestStep{
			Config: conf.ExistConfig,
			Check: resource.ComposeTestCheckFunc(
				existCheck(conf.ExistChangMap),
			),
		})
	}
	if conf.FakeConfig != "" {
		steps = append(steps, resource.TestStep{
			Config: conf.FakeConfig,
			Check: resource.ComposeTestCheckFunc(
				emptyCheck(conf.FakeChangMap),
			),
		})
	}
	return steps
}
