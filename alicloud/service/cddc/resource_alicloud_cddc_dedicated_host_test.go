package cddc_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cddc"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCddcDedicatedHost_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cddc_dedicated_host.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCDDCDedicatedHostMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cddc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDedicatedHost")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-cddcdedicatedhost%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCDDCDedicatedHostBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.CDDCSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dedicated_host_group_id": "${data.alicloud_cddc_dedicated_host_groups.default.ids.0}",
					"host_class":              "${data.alicloud_cddc_host_ecs_level_infos.default.infos.0.res_class_code}",
					"zone_id":                 "${data.alicloud_cddc_zones.default.ids.0}",
					"vswitch_id":              "${data.alicloud_vswitches.default.ids.0}",
					"payment_type":            "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_host_group_id": acctest.CHECKSET,
						"host_class":              acctest.CHECKSET,
						"zone_id":                 acctest.CHECKSET,
						"vswitch_id":              acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allocation_status": "Allocatable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allocation_status": "Allocatable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_class": "${data.alicloud_cddc_host_ecs_level_infos.default.infos.1.res_class_code}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_class": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "CDDC_DEDICATED",
						"For":     "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "CDDC_DEDICATED",
						"tags.For":     "TF",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_name":         "${var.name}_update",
					"allocation_status": "Suspended",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "CDDC_DEDICATED",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_name":         name + "_update",
						"allocation_status": "Suspended",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "CDDC_DEDICATED",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_category", "payment_type", "used_time", "period", "auto_renew", "os_password"},
			},
		},
	})
}

func TestAccAliCloudCddcDedicatedHost_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cddc_dedicated_host.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCDDCDedicatedHostMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cddc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDedicatedHost")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-cddcdedicatedhost%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCDDCDedicatedHostBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.CDDCSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dedicated_host_group_id": "${data.alicloud_cddc_dedicated_host_groups.default.ids.0}",
					"host_class":              "${data.alicloud_cddc_host_ecs_level_infos.default.infos.0.res_class_code}",
					"auto_renew":              "false",
					"zone_id":                 "${data.alicloud_cddc_zones.default.ids.0}",
					"vswitch_id":              "${data.alicloud_vswitches.default.ids.0}",
					"payment_type":            "Subscription",
					"host_name":               "${var.name}",
					"period":                  "Month",
					"used_time":               "1",
					"image_category":          "AliLinux",
					"os_password":             "Password1234.",
					"allocation_status":       "Allocatable",
					"tags": map[string]string{
						"Created": "CDDC_DEDICATED",
						"For":     "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_host_group_id": acctest.CHECKSET,
						"host_name":               name,
						"host_class":              acctest.CHECKSET,
						"zone_id":                 acctest.CHECKSET,
						"vswitch_id":              acctest.CHECKSET,
						"allocation_status":       "Allocatable",
						"tags.%":                  "2",
						"tags.Created":            "CDDC_DEDICATED",
						"tags.For":                "TF",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_category", "payment_type", "used_time", "period", "auto_renew", "os_password"},
			},
		},
	})
}

var AlicloudCDDCDedicatedHostMap0 = map[string]string{
	"status": acctest.CHECKSET,
}

func AlicloudCDDCDedicatedHostBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_cddc_zones" "default" {}

data "alicloud_cddc_host_ecs_level_infos" "default" {
  db_type      = "mysql"
  zone_id      = data.alicloud_cddc_zones.default.ids.0
  storage_type = "cloud_essd"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_cddc_dedicated_host_groups.default.groups.0.vpc_id
  zone_id = data.alicloud_cddc_zones.default.ids.0
}

data "alicloud_cddc_dedicated_host_groups" "default" {
  name_regex = "^NO-DELETING"
  engine     = "MySQL"
}

`, name)
}

// lintignore: R001
