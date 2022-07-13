package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEcsActivationsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.EcsActivationsSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsActivationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_activation.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcsActivationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_activation.default.id}_fake"]`,
		}),
	}
	instanceNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsActivationsDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_ecs_activation.default.id}"]`,
			"instance_name": `"${alicloud_ecs_activation.default.instance_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsActivationsDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_ecs_activation.default.id}"]`,
			"instance_name": `"${alicloud_ecs_activation.default.instance_name}_fake"`,
		}),
	}
	var existAlicloudEcsActivationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"activations.#":                       "1",
			"activations.0.description":           fmt.Sprintf("tf-testAccActivation-%d", rand),
			"activations.0.instance_count":        "10",
			"activations.0.instance_name":         fmt.Sprintf("tf-testAccActivation-%d", rand),
			"activations.0.ip_address_range":      "0.0.0.0/0",
			"activations.0.time_to_live_in_hours": "4",
			"activations.0.id":                    CHECKSET,
			"activations.0.activation_id":         CHECKSET,
			"activations.0.create_time":           CHECKSET,
			"activations.0.deregistered_count":    CHECKSET,
			"activations.0.disabled":              CHECKSET,
			"activations.0.registered_count":      CHECKSET,
		}
	}
	var fakeAlicloudEcsActivationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudEcsActivationsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecs_activations.default",
		existMapFunc: existAlicloudEcsActivationsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcsActivationsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudEcsActivationsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, instanceNameConf)
}
func testAccCheckAlicloudEcsActivationsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccActivation-%d"
}

resource "alicloud_ecs_activation" "default" {
	description = var.name
	instance_count = "10"
	instance_name = var.name
	ip_address_range = "0.0.0.0/0"
	time_to_live_in_hours = "4"
}

data "alicloud_ecs_activations" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
