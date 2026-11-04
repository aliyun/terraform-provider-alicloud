package swasopen_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudSimpleApplicationServerPlansDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	bandwidhthConf := acctest.DataSourceTestAccConfig{
		ExistConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"bandwidth": `"200"`,
		}),
		FakeConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"bandwidth": `"1"`,
		}),
	}
	memoryConf := acctest.DataSourceTestAccConfig{
		ExistConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"memory": `"1"`,
		}),
		FakeConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"memory": `"-1"`,
		}),
	}
	diskSizeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"disk_size": `"40"`,
		}),
		FakeConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"disk_size": `"-1"`,
		}),
	}
	flowConf := acctest.DataSourceTestAccConfig{
		ExistConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"flow": `"1200"`,
		}),
		FakeConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"flow": `"-1"`,
		}),
	}
	coreConf := acctest.DataSourceTestAccConfig{
		ExistConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"core": `"2"`,
		}),
		FakeConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"core": `"-1"`,
		}),
	}
	platformConf := acctest.DataSourceTestAccConfig{
		ExistConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"platform": `"Linux"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"memory":    `"4"`,
			"bandwidth": `"6"`,
			"disk_size": `"80"`,
			"flow":      `"1200"`,
			"core":      `"2"`,
			"platform":  `"Linux"`,
		}),
		FakeConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"memory":    `"4"`,
			"bandwidth": `"6"`,
			"disk_size": `"80"`,
			"flow":      `"1200"`,
			"core":      `"-1"`,
		}),
	}
	var existSimpleApplicationServerPlansMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    acctest.CHECKSET,
			"plans.#":                  acctest.CHECKSET,
			"plans.0.bandwidth":        acctest.CHECKSET,
			"plans.0.core":             acctest.CHECKSET,
			"plans.0.disk_size":        acctest.CHECKSET,
			"plans.0.flow":             acctest.CHECKSET,
			"plans.0.id":               acctest.CHECKSET,
			"plans.0.plan_id":          acctest.CHECKSET,
			"plans.0.memory":           acctest.CHECKSET,
			"plans.0.support_platform": acctest.CHECKSET,
		}
	}

	var fakeSimpleApplicationServerPlansMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"plans.#": "0",
		}
	}

	var simpleApplicationServerPlansCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_simple_application_server_plans.default",
		ExistMapFunc: existSimpleApplicationServerPlansMapFunc,
		FakeMapFunc:  fakeSimpleApplicationServerPlansMapFunc,
	}
	var preCheck = func() {
		acctest.PreCheckWithRegions(t, false, connectivity.SimpleApplicationServerNotSupportRegions)
	}
	simpleApplicationServerPlansCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, bandwidhthConf, memoryConf, diskSizeConf, flowConf, coreConf, platformConf, allConf)
}

func dataSourcesimpleApplicationServerPlansConfigDependence(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
data "alicloud_simple_application_server_plans" "default" {
  %s
}

`, strings.Join(pairs, "\n   "))
	return config
}
