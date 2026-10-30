package smartag_test

import (
	"fmt"
	"testing"

	smartagsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/smartag"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudSagQosCar_basic(t *testing.T) {
	var qospy smartagsdk.QosCar
	resourceId := "alicloud_sag_qos_car.default"
	ra := acctest.ResourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return smartag.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &qospy, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testQosCarName")
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceSagQosCarDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
		},

		// module name
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"qos_id":            "${alicloud_sag_qos.default.id}",
					"priority":          "2",
					"limit_type":        "Absolute",
					"min_bandwidth_abs": "5",
					"max_bandwidth_abs": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"qos_id":            acctest.CHECKSET,
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

func TestAccAliCloudSagQosCar_multi(t *testing.T) {
	var qospy smartagsdk.QosCar
	resourceId := "alicloud_sag_qos_car.default.2"
	ra := acctest.ResourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return smartag.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &qospy, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testQosCarName")
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceSagQosCarDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
		},

		// module name
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"qos_id":            "${alicloud_sag_qos.default.id}",
					"description":       "${var.name}-${count.index}",
					"count":             "3",
					"priority":          "${count.index+1}",
					"limit_type":        "Absolute",
					"min_bandwidth_abs": "5",
					"max_bandwidth_abs": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"qos_id":            acctest.CHECKSET,
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
  name = "${var.name}"
}
`, name)
}
