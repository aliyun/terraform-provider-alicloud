package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudSagQosCar_basic(t *testing.T) {
	var qospy smartag.QosCar
	resourceId := "alicloud_sag_qos_car.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &SagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &qospy, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testQosCarName")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSagQosCarDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"qos_id":            alicloud_sag_qos.default.id,
					"priority":          "2",
					"limit_type":        "Absolute",
					"min_bandwidth_abs": "5",
					"max_bandwidth_abs": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"qos_id":            CHECKSET,
						"priority":          "2",
						"limit_type":        "Absolute",
						"min_bandwidth_abs": "5",
						"max_bandwidth_abs": "10",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testSagQosCarDescription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testSagQosCarDescription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_bandwidth_abs": "8",
					"max_bandwidth_abs": "12",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_bandwidth_abs": "8",
						"max_bandwidth_abs": "12",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSagQosCar_multi(t *testing.T) {
	var qospy smartag.QosCar
	resourceId := "alicloud_sag_qos_car.default.2"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &SagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &qospy, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testQosCarName")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSagQosCarDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"qos_id":            alicloud_sag_qos.default.id,
					"description":       "${var.name}-${count.index}",
					"count":             "3",
					"priority":          "${count.index+1}",
					"limit_type":        "Absolute",
					"min_bandwidth_abs": "5",
					"max_bandwidth_abs": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"qos_id":            CHECKSET,
						"description":       fmt.Sprintf("%s-2", name),
						"priority":          "3",
						"limit_type":        "Absolute",
						"min_bandwidth_abs": "5",
						"max_bandwidth_abs": "10",
					}),
				),
			},
		},
	})
}

func resourceSagQosCarDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_sag_qos" "default" {
  name = var.name
}
`, name)
}
