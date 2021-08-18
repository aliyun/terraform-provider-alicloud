package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudLindormInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_lindorm_instance.default_0"
	ra := resourceAttrInit(resourceId, AlicloudLindormInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HitsdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccLindorminstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudLindormInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_category":             "cloud_efficiency",
					"payment_type":              "PayAsYouGo",
					"zone_id":                   "${local.zone_id}",
					"vswitch_id":                "${local.vswitch_id}",
					"instance_name":             "${var.name}",
					"file_engine_specification": "lindorm.c.xlarge",
					"file_engine_node_count":    "2",
					"instance_storage":          "480",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_category":             "cloud_efficiency",
						"payment_type":              "PayAsYouGo",
						"instance_name":             name,
						"file_engine_specification": "lindorm.c.xlarge",
						"file_engine_node_count":    "2",
						"instance_storage":          "480",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_white_list": []string{"118.118.118.118"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_white_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_white_list": []string{"117.117.117.117", "116.116.116.116"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_white_list.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"upgrade_type":     "upgrade-disk-size",
					"instance_storage": "2400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_storage": "2400",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_proection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_proection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":      name,
					"deletion_proection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":      name,
						"deletion_proection": "false",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"upgrade_type", "core_num", "group_name", "core_spec", "pricing_cycle", "duration"},
			},
		},
	})
}

func TestAccAlicloudLindormInstance_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_lindorm_instance.default_1"
	ra := resourceAttrInit(resourceId, AlicloudLindormInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HitsdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccLindorminstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudLindormInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_category":             "cloud_efficiency",
					"payment_type":              "PayAsYouGo",
					"zone_id":                   "${local.zone_id}",
					"vswitch_id":                "${local.vswitch_id}",
					"instance_name":             "${var.name}",
					"file_engine_specification": "lindorm.c.xlarge",
					"file_engine_node_count":    "2",
					"instance_storage":          "480",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_category":             "cloud_efficiency",
						"payment_type":              "PayAsYouGo",
						"instance_name":             name,
						"file_engine_specification": "lindorm.c.xlarge",
						"file_engine_node_count":    "2",
						"instance_storage":          "480",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_serires_engine_specification": "lindorm.g.2xlarge",
					"time_series_engine_node_count":     "2",
					"instance_storage":                  "960",
					"upgrade_type":                      "open-tsdb-engine",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_serires_engine_specification": "lindorm.g.2xlarge",
						"time_series_engine_node_count":     "2",
						"instance_storage":                  "960",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_serires_engine_specification": "lindorm.g.4xlarge",
					"upgrade_type":                      "upgrade-tsdb-engine",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_serires_engine_specification": "lindorm.g.4xlarge",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"upgrade_type":                  "upgrade-tsdb-core-num",
					"time_series_engine_node_count": "3",
					"instance_storage":              "1200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_series_engine_node_count": "3",
						"instance_storage":              "1200",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"upgrade_type", "core_num", "group_name", "core_spec", "pricing_cycle", "duration"},
			},
		},
	})
}

func TestAccAlicloudLindormInstance_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_lindorm_instance.default_2"
	ra := resourceAttrInit(resourceId, AlicloudLindormInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HitsdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccLindorminstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudLindormInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_category":              "cloud_efficiency",
					"payment_type":               "PayAsYouGo",
					"zone_id":                    "${local.zone_id}",
					"vswitch_id":                 "${local.vswitch_id}",
					"instance_name":              "${var.name}",
					"table_engine_specification": "lindorm.c.xlarge",
					"table_engine_node_count":    "2",
					"instance_storage":           "480",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_category":              "cloud_efficiency",
						"payment_type":               "PayAsYouGo",
						"instance_name":              name,
						"table_engine_specification": "lindorm.c.xlarge",
						"table_engine_node_count":    "2",
						"instance_storage":           "480",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"table_engine_specification": "lindorm.g.4xlarge",
					"upgrade_type":               "upgrade-lindorm-engine",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"table_engine_specification": "lindorm.g.4xlarge",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"upgrade_type":            "upgrade-lindorm-core-num",
					"table_engine_node_count": "3",
					"instance_storage":        "800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"table_engine_node_count": "3",
						"instance_storage":        "800",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"upgrade_type", "core_num", "group_name", "core_spec", "pricing_cycle", "duration"},
			},
		},
	})
}

func TestAccAlicloudLindormInstance_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_lindorm_instance.default_3"
	ra := resourceAttrInit(resourceId, AlicloudLindormInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HitsdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccLindorminstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudLindormInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_category":             "cloud_efficiency",
					"payment_type":              "PayAsYouGo",
					"zone_id":                   "${local.zone_id}",
					"vswitch_id":                "${local.vswitch_id}",
					"instance_name":             "${var.name}",
					"file_engine_specification": "lindorm.c.xlarge",
					"file_engine_node_count":    "2",
					"instance_storage":          "480",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_category":             "cloud_efficiency",
						"payment_type":              "PayAsYouGo",
						"instance_name":             name,
						"file_engine_specification": "lindorm.c.xlarge",
						"file_engine_node_count":    "2",
						"instance_storage":          "480",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"search_engine_specification": "lindorm.g.2xlarge",
					"search_engine_node_count":    "2",
					"instance_storage":            "960",
					"upgrade_type":                "open-search-engine",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"search_engine_specification": "lindorm.g.2xlarge",
						"search_engine_node_count":    "2",
						"instance_storage":            "960",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"search_engine_specification": "lindorm.g.4xlarge",
					"upgrade_type":                "upgrade-search-engine",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"search_engine_specification": "lindorm.g.4xlarge",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"upgrade_type":             "upgrade-search-core-num",
					"search_engine_node_count": "3",
					"instance_storage":         "1200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"search_engine_node_count": "3",
						"instance_storage":         "1200",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"upgrade_type", "core_num", "group_name", "core_spec", "pricing_cycle", "duration"},
			},
		},
	})
}

var AlicloudLindormInstanceMap0 = map[string]string{
	"cold_storage":                      CHECKSET,
	"search_engine_specification":       CHECKSET,
	"duration":                          NOSET,
	"deletion_proection":                CHECKSET,
	"file_engine_specification":         CHECKSET,
	"status":                            CHECKSET,
	"core_num":                          NOSET,
	"phoenix_node_count":                CHECKSET,
	"phoenix_node_specification":        CHECKSET,
	"group_name":                        NOSET,
	"lts_node_specification":            CHECKSET,
	"time_series_engine_node_count":     CHECKSET,
	"time_serires_engine_specification": CHECKSET,
	"file_engine_node_count":            CHECKSET,
	"lts_node_count":                    CHECKSET,
	"search_engine_node_count":          CHECKSET,
	"core_spec":                         NOSET,
	"pricing_cycle":                     NOSET,
	"table_engine_node_count":           CHECKSET,
	"instance_storage":                  "480",
	"zone_id":                           CHECKSET,
	"disk_category":                     "cloud_efficiency",
	"payment_type":                      "PayAsYouGo",
	"vswitch_id":                        CHECKSET,
	"instance_name":                     CHECKSET,
	"table_engine_specification":        CHECKSET,
}

func AlicloudLindormInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(`
provider "alicloud" {
  region = "cn-shenzhen"
}
variable "name" {
  default = "%s"
}
locals {
  zone_id = "cn-shenzhen-e"
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}
data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = local.zone_id
}
resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id = local.zone_id
  vswitch_name              = var.name
}
`, name)
}
