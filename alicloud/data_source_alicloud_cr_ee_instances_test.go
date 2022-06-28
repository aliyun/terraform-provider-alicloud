package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCREEInstancesDataSource(t *testing.T) {
	resourceId := "data.alicloud_cr_ee_instances.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-basic-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceCrEEInstancesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_cr_ee_instance.default.instance_name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_cr_ee_instance.default.instance_name}-fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_cr_ee_instance.default.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"test-id-fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     name,
			"ids":            []string{"${alicloud_cr_ee_instance.default.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_cr_ee_instance.default.id}"},
			"name_regex": "test-fake.*",
		}),
	}

	var existCrEEInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                         CHECKSET,
			"names.0":                         CHECKSET,
			"instances.#":                     CHECKSET,
			"instances.0.id":                  CHECKSET,
			"instances.0.name":                CHECKSET,
			"instances.0.region":              CHECKSET,
			"instances.0.specification":       CHECKSET,
			"instances.0.namespace_quota":     CHECKSET,
			"instances.0.namespace_usage":     CHECKSET,
			"instances.0.repo_quota":          CHECKSET,
			"instances.0.repo_usage":          CHECKSET,
			"instances.0.vpc_endpoints.#":     CHECKSET,
			"instances.0.public_endpoints.#":  CHECKSET,
			"instances.0.authorization_token": CHECKSET,
			"instances.0.temp_username":       CHECKSET,
		}
	}

	var fakeCrEEInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"instances.#": "0",
		}
	}

	var crEEInstancesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCrEEInstancesMapFunc,
		fakeMapFunc:  fakeCrEEInstancesMapFunc,
	}

	crEEInstancesCheckInfo.dataSourceTestCheck(t, 0, nameRegexConf, idsConf, allConf)
}

func dataSourceCrEEInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
resource "alicloud_cr_ee_instance" "default" {
  payment_type        = "Subscription"
  period              = 1
  renewal_status      = "ManualRenewal"
  instance_type       = "Advanced"
  instance_name       = var.name
}
`, name)
}
