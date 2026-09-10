package alicloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudECSCapacityReservation_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_capacity_reservation.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsCapacityReservationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsCapacityReservation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%secscapacityreservation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsCapacityReservationBasicDependence0)
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
					"description":               "${var.name}",
					"end_time":                  time.Now().Add(1 * time.Hour).Format("2006-01-02T15:04:05Z"),
					"platform":                  "linux",
					"capacity_reservation_name": "${var.name}",
					"end_time_type":             "Limited",
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"instance_amount":           "1",
					"instance_type":             "ecs.c5.2xlarge",
					"match_criteria":            "Open",
					"tags": map[string]string{
						"Created": "tfTestAcc0",
						"For":     "Tftestacc 0",
					},
					"zone_ids": []string{"${data.alicloud_zones.default.zones[0].id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":               name,
						"end_time":                  CHECKSET,
						"platform":                  "linux",
						"capacity_reservation_name": name,
						"end_time_type":             "Limited",
						"resource_group_id":         CHECKSET,
						"instance_amount":           "1",
						"instance_type":             "ecs.c5.2xlarge",
						"match_criteria":            "Open",
						"tags.%":                    "2",
						"zone_ids.#":                "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudEcsCapacityReservationMap0 = map[string]string{}

func AlicloudEcsCapacityReservationBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default"{
	status = "OK"
}

data "alicloud_zones" "default" {
  available_resource_creation = "Instance"
}
`, name)
}

func TestAccAlicloudECSCapacityReservation_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_capacity_reservation.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsCapacityReservationMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsCapacityReservation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%secscapacityreservation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsCapacityReservationBasicDependence1)
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
					"instance_amount": "1",
					"instance_type":   "ecs.c5.2xlarge",
					"zone_ids":        []string{"${data.alicloud_zones.default.zones[0].id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_amount": "1",
						"instance_type":   "ecs.c5.2xlarge",
						"zone_ids.#":      "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"end_time":      time.Now().Add(1 * time.Hour).Format("2006-01-02T15:04:05Z"),
					"end_time_type": "Limited",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"end_time":      CHECKSET,
						"end_time_type": "Limited",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"capacity_reservation_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"capacity_reservation_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_amount": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_amount": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "tfTestAcc8",
						"For":     "Tftestacc 8",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":               "${var.name}_update",
					"capacity_reservation_name": "${var.name}_update",
					"end_time_type":             "Unlimited",
					"instance_amount":           "1",
					"tags": map[string]string{
						"Created": "tfTestAcc9",
						"For":     "Tftestacc 9",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":               name + "_update",
						"capacity_reservation_name": name + "_update",
						"end_time_type":             "Unlimited",
						"instance_amount":           "1",
						"tags.%":                    "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudEcsCapacityReservationMap1 = map[string]string{}

func AlicloudEcsCapacityReservationBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}


data "alicloud_zones" "default" {
  available_resource_creation = "Instance"
}
`, name)
}
