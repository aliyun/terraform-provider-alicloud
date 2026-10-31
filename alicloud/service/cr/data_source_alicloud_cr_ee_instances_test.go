package cr_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCREEInstancesDataSource(t *testing.T) {
	resourceId := "data.alicloud_cr_ee_instances.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-basic-%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceCrEEInstancesConfigDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_cr_ee_instance.default.instance_name}",
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_cr_ee_instance.default.instance_name}-fake",
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_cr_ee_instance.default.id}"},
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"test-id-fake"},
		}),
	}

	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_cr_ee_instance.default.id}"},
			"enable_details": "true",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Test",
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_cr_ee_instance.default.id}"},
			"tags": map[string]string{
				"Created": "TF-fake",
			},
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":     name,
			"ids":            []string{"${alicloud_cr_ee_instance.default.id}"},
			"enable_details": "true",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Test",
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_cr_ee_instance.default.id}"},
			"name_regex": "test-fake.*",
		}),
	}

	var existCrEEInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                         acctest.CHECKSET,
			"names.0":                         acctest.CHECKSET,
			"instances.#":                     acctest.CHECKSET,
			"instances.0.id":                  acctest.CHECKSET,
			"instances.0.name":                acctest.CHECKSET,
			"instances.0.region":              acctest.CHECKSET,
			"instances.0.specification":       acctest.CHECKSET,
			"instances.0.namespace_quota":     acctest.CHECKSET,
			"instances.0.namespace_usage":     acctest.CHECKSET,
			"instances.0.repo_quota":          acctest.CHECKSET,
			"instances.0.repo_usage":          acctest.CHECKSET,
			"instances.0.vpc_endpoints.#":     acctest.CHECKSET,
			"instances.0.public_endpoints.#":  acctest.CHECKSET,
			"instances.0.authorization_token": acctest.CHECKSET,
			"instances.0.temp_username":       acctest.CHECKSET,
			"instances.0.tags.%":              "2",
			"instances.0.tags.Created":        "TF",
			"instances.0.tags.For":            "Test",
		}
	}

	var fakeCrEEInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"instances.#": "0",
		}
	}

	var crEEInstancesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existCrEEInstancesMapFunc,
		FakeMapFunc:  fakeCrEEInstancesMapFunc,
	}

	crEEInstancesCheckInfo.DataSourceTestCheck(t, 0, nameRegexConf, idsConf, tagsConf, allConf)
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
  tags = {
    Created = "TF"
    For     = "Test"
  }
}
`, name)
}
