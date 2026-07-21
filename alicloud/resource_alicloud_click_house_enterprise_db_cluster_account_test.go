package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ClickHouse EnterpriseDBClusterAccount. >>> Resource test cases, automatically generated.
// Case CK企业版账号3-线上 10563
func TestAccAliCloudClickHouseEnterpriseDBClusterAccount_basic10563(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_click_house_enterprise_db_cluster_account.default"
	ra := resourceAttrInit(resourceId, AlicloudClickHouseEnterpriseDBClusterAccountMap10563)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ClickHouseServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeClickHouseEnterpriseDBClusterAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccclickhouse%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudClickHouseEnterpriseDBClusterAccountBasicDependence10563)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"account":        "abc",
					"description":    "test_desc",
					"db_instance_id": "${alicloud_click_house_enterprise_db_cluster.defaultWrovOd.id}",
					"account_type":   "NormalAccount",
					"password":       "abc123456!",
					"dml_auth_setting": []map[string]interface{}{
						{
							"dml_authority": "0",
							"ddl_authority": "true",
							"allow_dictionaries": []string{
								"*"},
							"allow_databases": []string{
								"*"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account":        "abc",
						"description":    "test_desc",
						"db_instance_id": CHECKSET,
						"account_type":   "NormalAccount",
						"password":       "abc123456!",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test_desc_new",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test_desc_new",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "abc1234567!",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "abc1234567!",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dml_auth_setting": []map[string]interface{}{
						{
							"dml_authority": "1",
							"ddl_authority": "true",
							"allow_dictionaries": []string{
								"*"},
							"allow_databases": []string{
								"*"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dml_auth_setting": []map[string]interface{}{
						{
							"dml_authority": "1",
							"ddl_authority": "false",
							"allow_dictionaries": []string{
								"*"},
							"allow_databases": []string{
								"*"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dml_auth_setting": []map[string]interface{}{
						{
							"dml_authority": "1",
							"ddl_authority": "false",
							"allow_databases": []string{
								"default"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dml_auth_setting": []map[string]interface{}{
						{
							"dml_authority": "1",
							"ddl_authority": "false",
							"allow_dictionaries": []string{
								"system.test"},
							"allow_databases": []string{
								"default"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

var AlicloudClickHouseEnterpriseDBClusterAccountMap10563 = map[string]string{}

func AlicloudClickHouseEnterpriseDBClusterAccountBasicDependence10563(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-beijing"
}

variable "vsw_ip_range_i" {
  default = "172.16.1.0/24"
}

variable "vpc_ip_range" {
  default = "172.16.0.0/12"
}

variable "zone_id_i" {
  default = "cn-beijing-i"
}

resource "alicloud_vpc" "defaultktKLuM" {
  cidr_block = var.vpc_ip_range
}

resource "alicloud_vswitch" "defaultTQWN3k" {
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  zone_id    = var.zone_id_i
  cidr_block = var.vsw_ip_range_i
}

resource "alicloud_click_house_enterprise_db_cluster" "defaultWrovOd" {
  zone_id    = var.zone_id_i
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  scale_min  = "8"
  scale_max  = "16"
  vswitch_id = alicloud_vswitch.defaultTQWN3k.id
}


`, name)
}

// Test ClickHouse EnterpriseDBClusterAccount. <<< Resource test cases, automatically generated.
