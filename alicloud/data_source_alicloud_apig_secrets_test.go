package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAliCloudApigSecretsDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_apig_secrets.default"
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceApigSecretsConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_apig_secret.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_apig_secret.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_apig_secret.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_apig_secret.default.name}_fake",
		}),
	}

	gatewayTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"gateway_type": "${alicloud_apig_secret.default.gateway_type}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"gateway_type": "AI",
		}),
	}

	nameLikeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_like": "${alicloud_apig_secret.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_like": "${alicloud_apig_secret.default.name}_fake",
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"status": "${alicloud_apig_secret.default.status}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"status": "DELETED",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":          []string{"${alicloud_apig_secret.default.id}"},
			"name_regex":   "${alicloud_apig_secret.default.name}",
			"gateway_type": "${alicloud_apig_secret.default.gateway_type}",
			"name_like":    "${alicloud_apig_secret.default.name}",
			"status":       "${alicloud_apig_secret.default.status}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
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
			"secrets.0.id":                           CHECKSET,
			"secrets.0.secret_id":                    CHECKSET,
			"secrets.0.gateway_type":                 CHECKSET,
			"secrets.0.name":                         CHECKSET,
			"secrets.0.description":                  CHECKSET,
			"secrets.0.secret_source":                CHECKSET,
			"secrets.0.reference_count":              CHECKSET,
			"secrets.0.status":                       CHECKSET,
			"secrets.0.create_timestamp":             CHECKSET,
			"secrets.0.update_timestamp":             CHECKSET,
			"secrets.0.kms_config.#":                 "1",
			"secrets.0.kms_config.0.kms_instance_id": CHECKSET,
			"secrets.0.kms_config.0.kms_key_id":      CHECKSET,
			"secrets.0.kms_config.0.kms_secret_arn":  CHECKSET,
			"secrets.0.kms_config.0.version_id":      CHECKSET,
		}
	}

	var fakeAliCloudApigSecretsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"secrets.#": "0",
		}
	}

	var aliCloudApigSecretsInfo = dataSourceAttr{
		resourceId:   "data.alicloud_apig_secrets.default",
		existMapFunc: existAliCloudApigSecretsMapFunc,
		fakeMapFunc:  fakeAliCloudApigSecretsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudApigSecretsInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, gatewayTypeConf, nameLikeConf, statusConf, allConf)
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
