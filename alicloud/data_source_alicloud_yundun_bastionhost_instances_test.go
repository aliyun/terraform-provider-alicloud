package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudYundunBastionhostInstanceDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_yundun_bastionhost_instances.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, fmt.Sprintf("tf_testAcc%d", rand),
		dataSourceYundunBastionhostInstanceConfigDependency)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{alicloud_yundun_bastionhost_instance.default.id},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_yundun_bastionhost_instance.default.id}-fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": alicloud_yundun_bastionhost_instance.default.description,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_yundun_bastionhost_instance.default.description}-fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": alicloud_yundun_bastionhost_instance.default.description,
			"ids":               []string{alicloud_yundun_bastionhost_instance.default.id},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_yundun_bastionhost_instance.default.description}-fake",
			"ids":               []string{"${alicloud_yundun_bastionhost_instance.default.id}-fake"},
		}),
	}

	var existYundunBastionhostInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"descriptions.#":                    "1",
			"ids.0":                             CHECKSET,
			"descriptions.0":                    fmt.Sprintf("tf_testAcc%d", rand),
			"instances.#":                       "1",
			"instances.0.description":           fmt.Sprintf("tf_testAcc%d", rand),
			"instances.0.license_code":          "bhah_ent_50_asset",
			"instances.0.user_vswitch_id":       CHECKSET,
			"instances.0.public_network_access": "true",
			"instances.0.private_domain":        CHECKSET,
			"instances.0.instance_status":       CHECKSET,
			"instances.0.security_group_ids.#":  "1",
		}
	}
	var fakeYundunBastionhostInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"descriptions.#": "0",
		}
	}
	var yundunBastionhostInstanceCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_yundun_bastionhost_instances.default",
		existMapFunc: existYundunBastionhostInstanceMapFunc,
		fakeMapFunc:  fakeYundunBastionhostInstanceMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.YundunBastionhostSupportedRegions)
		testAccPreCheckWithAccountSiteType(t, DomesticSite)
	}

	yundunBastionhostInstanceCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, allConf)

}

func dataSourceYundunBastionhostInstanceConfigDependency(description string) string {
	return fmt.Sprintf(
		`data "alicloud_zones" "default" {
				  available_resource_creation = "VSwitch"
				}
				
				variable "name" {
				  default = "%s"
				}
				
				resource "alicloud_vpc" "default" {
				  name       = var.name
				  cidr_block = "172.16.0.0/12"
				}
				
				resource "alicloud_vswitch" "default" {
				  vpc_id            = alicloud_vpc.default.id
				  cidr_block        = "172.16.0.0/21"
				  availability_zone = data.alicloud_zones.default.zones.0.id
				  name              = var.name
				}
				
				resource "alicloud_security_group" "default" {
				  name   = var.name
				  vpc_id = alicloud_vpc.default.id
				}
				
				provider "alicloud" {
				  endpoints {
					bssopenapi = "business.aliyuncs.com"
				  }
				}
				
				resource "alicloud_yundun_bastionhost_instance" "default" {
				  description        = var.name
				  license_code       = "bhah_ent_50_asset"
				  period             = "1"
				  vswitch_id         = alicloud_vswitch.default.id
				  security_group_ids = [alicloud_security_group.default.id]
				}`, description)
}
