package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudReservedInstanceBasic(t *testing.T) {
	var v ecs.ReservedInstance

	resourceId := "alicloud_reserved_instance.default"
	ra := resourceAttrInit(resourceId, testAccReservedInstanceCheckMap)

	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeReservedInstance")
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEcsReservedInstanceConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceReservedInstanceBasicConfigDependence)
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
					"instance_type":     "ecs.g6.large",
					"instance_amount":   "1",
					"period_unit":       "Month",
					"offering_type":     "All Upfront",
					"name":              name,
					"description":       "ReservedInstance",
					"zone_id":           "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}",
					"scope":             "Zone",
					"period":            "1",
					"renewal_status":    "AutoRenewal",
					"auto_renew_period": "36",
					"tags":              map[string]interface{}{"Created": "TF", "Foo": "Bar"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":     "ecs.g6.large",
						"instance_amount":   "1",
						"period_unit":       "Month",
						"offering_type":     "All Upfront",
						"name":              name,
						"description":       "ReservedInstance",
						"zone_id":           CHECKSET,
						"scope":             "Zone",
						"renewal_status":    "AutoRenewal",
						"auto_renew_period": "36",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.Foo":          "Bar",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "period_unit"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "change",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "ReservedInstanceChange",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "ReservedInstanceChange",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        name,
					"description": "ReservedInstance",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name,
						"description": "ReservedInstance",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]interface{}{"Created": "TF1", "Foo": "Bar1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF1",
						"tags.Foo":     "Bar1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew_period": "12",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew_period": "12",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renewal_status": "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_status": "Normal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew_period": "36",
					"renewal_status":    "AutoRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew_period": "36",
						"renewal_status":    "AutoRenewal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renewal_status": "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_status": "Normal",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "period_unit", "auto_renew_period"},
			},
		},
	})
}

func TestAccAliCloudReservedInstanceBasic1(t *testing.T) {
	var v ecs.ReservedInstance

	resourceId := "alicloud_reserved_instance.default"
	ra := resourceAttrInit(resourceId, testAccReservedInstanceCheckMap)

	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeReservedInstance")
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEcsReservedInstanceConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceReservedInstanceBasicConfigDependence)
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
					"instance_type":          "${data.alicloud_instance_types.default.instance_types.0.id}",
					"instance_amount":        "1",
					"period":                 "1",
					"period_unit":            "Month",
					"offering_type":          "All Upfront",
					"reserved_instance_name": name,
					"description":            "ReservedInstance",
					"zone_id":                "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}",
					"scope":                  "Zone",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":          CHECKSET,
						"instance_amount":        "1",
						"period":                 "1",
						"period_unit":            "Month",
						"offering_type":          "All Upfront",
						"reserved_instance_name": name,
						"description":            "ReservedInstance",
						"zone_id":                CHECKSET,
						"scope":                  "Zone",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"reserved_instance_name": name + "change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"reserved_instance_name": name + "change",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "period_unit"},
			},
		},
	})
}

var testAccReservedInstanceCheckMap = map[string]string{
	"allocation_status": "",
	"create_time":       CHECKSET,
	"expired_time":      CHECKSET,
	"start_time":        CHECKSET,
	"status":            CHECKSET,
	"operation_locks.#": "0",
}

func resourceReservedInstanceBasicConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	
	data "alicloud_instance_types" "default" {
		instance_type_family = "ecs.g7"
	}
`, name)
}
