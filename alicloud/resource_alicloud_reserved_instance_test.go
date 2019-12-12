package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
			testAccPreCheckWithTime(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{

				Config: testAccConfig(map[string]interface{}{
					"instance_type":   "ecs.g6.large",
					"instance_amount": "1",
					"period_unit":     "Year",
					"offering_type":   "All Upfront",
					"name":            name,
					"description":     "ReservedInstance",
					"zone_id":         "cn-shanghai-g",
					"scope":           "Zone",
					"period":          "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":   "ecs.g6.large",
						"instance_amount": "1",
						"period_unit":     "Year",
						"offering_type":   "All Upfront",
						"name":            name,
						"description":     "ReservedInstance",
						"zone_id":         "cn-shanghai-g",
						"scope":           "Zone",
						"period":          "1",
					}),
				),
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
		},
	})
}

var testAccReservedInstanceCheckMap = map[string]string{}

func resourceReservedInstanceBasicConfigDependence(name string) string {
	return ""
}
