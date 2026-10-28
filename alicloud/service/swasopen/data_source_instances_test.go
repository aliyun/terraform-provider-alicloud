package swasopen_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudSimpleApplicationServerInstanceDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSimpleApplicationServerInstanceDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_simple_application_server_instance.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudSimpleApplicationServerInstanceDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_simple_application_server_instance.default.id}_fake"]`,
		}),
	}

	paymentTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSimpleApplicationServerInstanceDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_simple_application_server_instance.default.id}"]`,
			"payment_type": `"Subscription"`,
		}),
		FakeConfig: testAccCheckAlicloudSimpleApplicationServerInstanceDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_simple_application_server_instance.default.id}_fake"]`,
			"payment_type": `"Subscription"`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSimpleApplicationServerInstanceDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_simple_application_server_instance.default.instance_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudSimpleApplicationServerInstanceDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_simple_application_server_instance.default.instance_name}_fake"`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSimpleApplicationServerInstanceDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_simple_application_server_instance.default.id}"]`,
			"status": `"Running"`,
		}),
		FakeConfig: testAccCheckAlicloudSimpleApplicationServerInstanceDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_simple_application_server_instance.default.id}"]`,
			"status": `"Disabled"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSimpleApplicationServerInstanceDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_simple_application_server_instance.default.id}"]`,
			"payment_type": `"Subscription"`,
			"name_regex":   `"${alicloud_simple_application_server_instance.default.instance_name}"`,
			"status":       `"Running"`,
		}),
		FakeConfig: testAccCheckAlicloudSimpleApplicationServerInstanceDataSourceName(rand, map[string]string{
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
	var alicloudSimpleApplicationServerInstanceCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_simple_application_server_instances.default",
		ExistMapFunc: existDataAlicloudSimpleApplicationServerInstancesSourceNameMapFunc,
		FakeMapFunc:  fakeDataAlicloudSimpleApplicationServerInstancesSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckPrePaidResources(t)
		acctest.PreCheckWithRegions(t, false, connectivity.SimpleApplicationServerNotSupportRegions)
	}
	alicloudSimpleApplicationServerInstanceCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, paymentTypeConf, statusConf, allConf)
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
