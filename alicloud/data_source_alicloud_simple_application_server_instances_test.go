package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSimpleApplicationServerInstanceDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSimpleApplicationServerInstanceDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_simple_application_server_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSimpleApplicationServerInstanceDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_simple_application_server_instance.default.id}_fake"]`,
		}),
	}

	paymentTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSimpleApplicationServerInstanceDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_simple_application_server_instance.default.id}"]`,
			"payment_type": `"Subscription"`,
		}),
		fakeConfig: testAccCheckAlicloudSimpleApplicationServerInstanceDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_simple_application_server_instance.default.id}_fake"]`,
			"payment_type": `"Subscription"`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSimpleApplicationServerInstanceDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_simple_application_server_instance.default.instance_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSimpleApplicationServerInstanceDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_simple_application_server_instance.default.instance_name}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSimpleApplicationServerInstanceDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_simple_application_server_instance.default.id}"]`,
			"status": `"Running"`,
		}),
		fakeConfig: testAccCheckAlicloudSimpleApplicationServerInstanceDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_simple_application_server_instance.default.id}"]`,
			"status": `"Disabled"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSimpleApplicationServerInstanceDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_simple_application_server_instance.default.id}"]`,
			"payment_type": `"Subscription"`,
			"name_regex":   `"${alicloud_simple_application_server_instance.default.instance_name}"`,
			"status":       `"Running"`,
		}),
		fakeConfig: testAccCheckAlicloudSimpleApplicationServerInstanceDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_simple_application_server_instance.default.id}_fake"]`,
			"payment_type": `"Subscription"`,
			"status":       `"Disabled"`,
			"name_regex":   `"${alicloud_simple_application_server_instance.default.instance_name}_fake"`,
		}),
	}

	var existDataAlicloudSimpleApplicationServerInstancesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"instances.#":               "1",
			"instances.0.instance_name": fmt.Sprintf("tf-testaccswas%d", rand),
			"instances.0.status":        "Running",
		}
	}
	var fakeDataAlicloudSimpleApplicationServerInstancesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}
	var alicloudSimpleApplicationServerInstanceCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_simple_application_server_instances.default",
		existMapFunc: existDataAlicloudSimpleApplicationServerInstancesSourceNameMapFunc,
		fakeMapFunc:  fakeDataAlicloudSimpleApplicationServerInstancesSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, false, connectivity.SimpleApplicationServerNotSupportRegions)
	}
	alicloudSimpleApplicationServerInstanceCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, paymentTypeConf, statusConf, allConf)
}
func testAccCheckAlicloudSimpleApplicationServerInstanceDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testaccswas%d"
}

data "alicloud_simple_application_server_images" "default" {
	platform = "Linux"
}
data "alicloud_simple_application_server_plans" "default" {
	platform = "Linux"
}

resource "alicloud_simple_application_server_instance" "default" {
  payment_type   = "Subscription"
  plan_id        = data.alicloud_simple_application_server_plans.default.plans.0.id
  instance_name  = var.name
  image_id       = data.alicloud_simple_application_server_images.default.images.0.id
  period         = 1
  data_disk_size = 100
}

data "alicloud_simple_application_server_instances" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
