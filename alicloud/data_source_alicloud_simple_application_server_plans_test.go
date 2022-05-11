package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSimpleApplicationServerPlansDataSource(t *testing.T) {
	rand := acctest.RandInt()
	bandwidhthConf := dataSourceTestAccConfig{
		existConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"bandwidth": `"30"`,
		}),
		fakeConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"bandwidth": `"1"`,
		}),
	}
	memoryConf := dataSourceTestAccConfig{
		existConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"memory": `"1"`,
		}),
		fakeConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"memory": `"-1"`,
		}),
	}
	diskSizeConf := dataSourceTestAccConfig{
		existConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"disk_size": `"40"`,
		}),
		fakeConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"disk_size": `"-1"`,
		}),
	}
	flowConf := dataSourceTestAccConfig{
		existConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"flow": `"3072"`,
		}),
		fakeConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"flow": `"-1"`,
		}),
	}
	coreConf := dataSourceTestAccConfig{
		existConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"core": `"2"`,
		}),
		fakeConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"core": `"-1"`,
		}),
	}
	platformConf := dataSourceTestAccConfig{
		existConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"platform": `"Linux"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"memory":    `"1"`,
			"bandwidth": `"30"`,
			"disk_size": `"40"`,
			"flow":      `"1024"`,
			"core":      `"2"`,
			"platform":  `"Linux"`,
		}),
		fakeConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{
			"memory":    `"1"`,
			"bandwidth": `"30"`,
			"disk_size": `"40"`,
			"flow":      `"3072"`,
			"core":      `"-1"`,
		}),
	}
	var existSimpleApplicationServerPlansMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    CHECKSET,
			"plans.#":                  CHECKSET,
			"plans.0.bandwidth":        CHECKSET,
			"plans.0.core":             CHECKSET,
			"plans.0.disk_size":        CHECKSET,
			"plans.0.flow":             CHECKSET,
			"plans.0.id":               CHECKSET,
			"plans.0.plan_id":          CHECKSET,
			"plans.0.memory":           CHECKSET,
			"plans.0.support_platform": CHECKSET,
		}
	}

	var fakeSimpleApplicationServerPlansMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"plans.#": "0",
		}
	}

	var simpleApplicationServerPlansCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_simple_application_server_plans.default",
		existMapFunc: existSimpleApplicationServerPlansMapFunc,
		fakeMapFunc:  fakeSimpleApplicationServerPlansMapFunc,
	}
	var preCheck = func() {
		testAccPreCheckWithRegions(t, false, connectivity.SimpleApplicationServerNotSupportRegions)
	}
	simpleApplicationServerPlansCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, bandwidhthConf, memoryConf, diskSizeConf, flowConf, coreConf, platformConf, allConf)
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
