package apig_test

import (
	"fmt"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudApigSecretsDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_apig_secrets.default"
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceApigSecretsConfig)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_apig_secret.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_apig_secret.default.id}_fake"},
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_apig_secret.default.name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_apig_secret.default.name}_fake",
		}),
	}

	gatewayTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"gateway_type": "${alicloud_apig_secret.default.gateway_type}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"gateway_type": "AI",
		}),
	}

	nameLikeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_like": "${alicloud_apig_secret.default.name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_like": "${alicloud_apig_secret.default.name}_fake",
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"status": "${alicloud_apig_secret.default.status}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"status": "DELETED",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":          []string{"${alicloud_apig_secret.default.id}"},
			"name_regex":   "${alicloud_apig_secret.default.name}",
			"gateway_type": "${alicloud_apig_secret.default.gateway_type}",
			"name_like":    "${alicloud_apig_secret.default.name}",
			"status":       "${alicloud_apig_secret.default.status}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":          []string{"${alicloud_apig_secret.default.id}_fake"},
			"name_regex":   "${alicloud_apig_secret.default.name}_fake",
			"gateway_type": "AI",
			"name_like":    "${alicloud_apig_secret.default.name}_fake",
			"status":       "DELETED",
		}),
	}

	var existAliCloudApigSecretsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"names.#":                                "1",
			"secrets.#":                              "1",
			"secrets.0.id":                           acctest.CHECKSET,
			"secrets.0.secret_id":                    acctest.CHECKSET,
			"secrets.0.gateway_type":                 acctest.CHECKSET,
			"secrets.0.name":                         acctest.CHECKSET,
			"secrets.0.description":                  acctest.CHECKSET,
			"secrets.0.secret_source":                acctest.CHECKSET,
			"secrets.0.reference_count":              acctest.CHECKSET,
			"secrets.0.status":                       acctest.CHECKSET,
			"secrets.0.create_timestamp":             acctest.CHECKSET,
			"secrets.0.update_timestamp":             acctest.CHECKSET,
			"secrets.0.kms_config.#":                 "1",
			"secrets.0.kms_config.0.kms_instance_id": acctest.CHECKSET,
			"secrets.0.kms_config.0.kms_key_id":      acctest.CHECKSET,
			"secrets.0.kms_config.0.kms_secret_arn":  acctest.CHECKSET,
			"secrets.0.kms_config.0.version_id":      acctest.CHECKSET,
		}
	}

	var fakeAliCloudApigSecretsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"secrets.#": "0",
		}
	}

	var aliCloudApigSecretsInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_apig_secrets.default",
		ExistMapFunc: existAliCloudApigSecretsMapFunc,
		FakeMapFunc:  fakeAliCloudApigSecretsMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudApigSecretsInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, gatewayTypeConf, nameLikeConf, statusConf, allConf)
}

func dataSourceApigSecretsConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-hangzhou-j"
}

resource "alicloud_kms_instance" "default" {
  product_version = "3"
  vpc_num         = "1"
  key_num         = "1000"
  secret_num      = "1000"
  spec            = "1000"
  vpc_id          = data.alicloud_vpcs.default.ids.0
  vswitch_ids = [
    data.alicloud_vswitches.default.ids.0
  ]
  zone_ids = [
    "cn-hangzhou-k",
    "cn-hangzhou-j"
  ]
}

resource "alicloud_kms_key" "default" {
  dkms_instance_id       = alicloud_kms_instance.default.id
  pending_window_in_days = 7
}

resource "alicloud_kms_secret" "default" {
  secret_data                   = var.name
  secret_name                   = var.name
  version_id                    = "v1"
  dkms_instance_id              = alicloud_kms_key.default.dkms_instance_id
  encryption_key_id             = alicloud_kms_key.default.id
  force_delete_without_recovery = true
}

resource "alicloud_apig_secret" "default" {
  gateway_type  = "API"
  name          = var.name
  secret_source = "KMS"
  secret_data   = alicloud_kms_secret.default.secret_data
  description   = var.name
  kms_config {
    kms_instance_id = alicloud_kms_secret.default.dkms_instance_id
    kms_key_id      = alicloud_kms_secret.default.encryption_key_id
  }
}
`, name)
}
