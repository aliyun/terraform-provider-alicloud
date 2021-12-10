package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudMongoDBServerlessInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_serverless_instance.default"
	checkoutSupportedRegions(t, true, connectivity.MongoDBServerlessSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudMongoDBServerlessInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMongodbServerlessInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-mongodbserverlessinstance-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMongoDBServerlessInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password":    "Abc12345",
					"db_instance_storage": "5",
					"capacity_unit":       "100",
					"engine_version":      "4.2",
					"vswitch_id":          "${local.vswitch_id}",
					"vpc_id":              "${data.alicloud_vpcs.default.ids.0}",
					"zone_id":             "${data.alicloud_mongodb_zones.default.zones.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage": "5",
						"capacity_unit":       "100",
						"engine_version":      "4.2",
						"vswitch_id":          CHECKSET,
						"vpc_id":              CHECKSET,
						"zone_id":             CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "MongodbServerlessInstance",
						"For":     "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "MongodbServerlessInstance",
						"tags.For":     "TF",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_description": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_start_time": "02:00Z",
					"maintain_end_time":   "03:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_start_time": "02:00Z",
						"maintain_end_time":   "03:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_groups": []map[string]interface{}{
						{
							"security_ip_group_attribute": "test",
							"security_ip_group_name":      "test",
							"security_ip_list":            "192.168.0.1",
						},
						{
							"security_ip_group_attribute": "test1",
							"security_ip_group_name":      "test1",
							"security_ip_list":            "192.168.0.2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_groups.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"capacity_unit":       "2000",
					"db_instance_storage": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"capacity_unit":       "2000",
						"db_instance_storage": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"capacity_unit":           "100",
					"db_instance_storage":     "5",
					"db_instance_description": "${var.name}_update",
					"maintain_start_time":     "01:00Z",
					"maintain_end_time":       "02:00Z",
					"tags": map[string]string{
						"Created": "MongodbServerlessInstance1",
						"For":     "TF1",
					},
					"security_ip_groups": []map[string]interface{}{
						{
							"security_ip_group_attribute": "test3",
							"security_ip_group_name":      "test3",
							"security_ip_list":            "192.168.0.3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"capacity_unit":           "100",
						"db_instance_storage":     "5",
						"db_instance_description": name + "_update",
						"maintain_start_time":     "01:00Z",
						"maintain_end_time":       "02:00Z",
						"tags.%":                  "2",
						"tags.Created":            "MongodbServerlessInstance1",
						"tags.For":                "TF1",
						"security_ip_groups.#":    "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "period", "period_price_type", "account_password"},
			},
		},
	})
}

func TestAccAlicloudMongoDBServerlessInstance_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_serverless_instance.default"
	checkoutSupportedRegions(t, true, connectivity.MongoDBServerlessSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudMongoDBServerlessInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMongodbServerlessInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-mongodbserverlessinstance-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMongoDBServerlessInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew":              "false",
					"account_password":        "Abc12345",
					"capacity_unit":           "100",
					"db_instance_storage":     "5",
					"storage_engine":          "WiredTiger",
					"engine":                  "MongoDB",
					"resource_group_id":       "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"engine_version":          "4.2",
					"db_instance_description": "${var.name}",
					"vpc_id":                  "${data.alicloud_vpcs.default.ids.0}",
					"zone_id":                 "${data.alicloud_mongodb_zones.default.zones.0.id}",
					"vswitch_id":              "${local.vswitch_id}",
					"period":                  "1",
					"period_price_type":       "Day",
					"tags": map[string]string{
						"Created": "MongodbServerlessInstance",
						"For":     "TF",
					},
					"security_ip_groups": []map[string]interface{}{
						{
							"security_ip_group_attribute": "test",
							"security_ip_group_name":      "test",
							"security_ip_list":            "192.168.0.1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"capacity_unit":           "100",
						"db_instance_storage":     "5",
						"engine":                  "MongoDB",
						"engine_version":          "4.2",
						"storage_engine":          "WiredTiger",
						"db_instance_description": name,
						"vswitch_id":              CHECKSET,
						"tags.%":                  "2",
						"tags.Created":            "MongodbServerlessInstance",
						"tags.For":                "TF",
						"security_ip_groups.#":    "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "period", "period_price_type", "account_password"},
			},
		},
	})
}

var AlicloudMongoDBServerlessInstanceMap0 = map[string]string{
	"auto_pay": NOSET,
	"status":   CHECKSET,
}

func AlicloudMongoDBServerlessInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_mongodb_zones" "default" {}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_mongodb_zones.default.zones.0.id
}

resource "alicloud_vswitch" "default" {
	count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
	vswitch_name = var.name
	vpc_id       = data.alicloud_vpcs.default.ids.0
    zone_id      = data.alicloud_mongodb_zones.default.zones.0.id
    cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
}

locals {
  vswitch_id  = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.default.*.id, [""])[0]
}
`, name)
}
