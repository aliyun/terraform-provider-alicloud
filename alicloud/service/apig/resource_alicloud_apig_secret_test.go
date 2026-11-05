// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package apig_test

import (
	"fmt"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	apig "github.com/aliyun/terraform-provider-alicloud/alicloud/service/apig"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ApigSecret. >>> Resource test cases, automatically generated.
// Case ApigSecret 1668
func TestAccAliCloudApigSecret_basic1668(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_secret.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudApigSecretMap1668)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apig.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApigSecret")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudApigSecretBasicDependence1668)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_type":  "API",
					"name":          name,
					"secret_source": "KMS",
					"secret_data":   "${alicloud_kms_secret.default.secret_data}",
					"kms_config": []map[string]interface{}{
						{
							"kms_instance_id": "${alicloud_kms_secret.default.dkms_instance_id}",
							"kms_key_id":      "${alicloud_kms_secret.default.encryption_key_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_type":  "API",
						"name":          name,
						"secret_source": "KMS",
						"kms_config.#":  "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secret_data": "${alicloud_kms_secret.update.secret_data}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_data"},
			},
		},
	})
}

func TestAccAliCloudApigSecret_basic1668_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_secret.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudApigSecretMap1668)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apig.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApigSecret")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudApigSecretBasicDependence1668)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_type":  "API",
					"name":          name,
					"secret_source": "KMS",
					"secret_data":   "${alicloud_kms_secret.default.secret_data}",
					"description":   name,
					"kms_config": []map[string]interface{}{
						{
							"kms_instance_id": "${alicloud_kms_secret.default.dkms_instance_id}",
							"kms_key_id":      "${alicloud_kms_secret.default.encryption_key_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_type":  "API",
						"name":          name,
						"secret_source": "KMS",
						"description":   name,
						"kms_config.#":  "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_data"},
			},
		},
	})
}

var AliCloudApigSecretMap1668 = map[string]string{
	"reference_count":  acctest.CHECKSET,
	"status":           acctest.CHECKSET,
	"create_timestamp": acctest.CHECKSET,
	"update_timestamp": acctest.CHECKSET,
}

func AliCloudApigSecretBasicDependence1668(name string) string {
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

resource "alicloud_kms_secret" "update" {
  secret_data                   = "${var.name}update"
  secret_name                   = "${var.name}update"
  version_id                    = "v1"
  dkms_instance_id              = alicloud_kms_key.default.dkms_instance_id
  encryption_key_id             = alicloud_kms_key.default.id
  force_delete_without_recovery = true
}
`, name)
}

// Test ApigSecret. <<< Resource test cases, automatically generated.
