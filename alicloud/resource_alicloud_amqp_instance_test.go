package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudAmqpInstance_professional(t *testing.T) {

	var v map[string]interface{}
	resourceId := "alicloud_amqp_instance.default"
	ra := resourceAttrInit(resourceId, AmqpInstanceBasicMap)
	serviceFunc := func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-AmqpInstanceprofessional%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAmqpInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type":  "professional",
					"max_tps":        "1000",
					"payment_type":   "Subscription",
					"period":         "1",
					"queue_capacity": "50",
					"support_eip":    "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":         "professional",
						"max_tps":               "1000",
						"payment_type":          "Subscription",
						"queue_capacity":        "50",
						"support_eip":           "false",
						"renewal_duration":      NOSET,
						"renewal_duration_unit": NOSET,
						"renewal_status":        "ManualRenewal",
						"status":                "SERVING",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modify_type": "Upgrade",
					"max_tps":     "1500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_tps": "1500",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"queue_capacity": "55",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"queue_capacity": "55",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"support_eip": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"support_eip": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_eip_tps": "256",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_eip_tps": "256",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renewal_duration":      "1",
					"renewal_duration_unit": "Month",
					"renewal_status":        "AutoRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_duration":      "1",
						"renewal_duration_unit": "Month",
						"renewal_status":        "AutoRenewal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modify_type":    "Downgrade",
					"max_tps":        "1000",
					"queue_capacity": "50",
					"support_eip":    "false",
					"renewal_status": "NotRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_tps":               "1000",
						"queue_capacity":        "50",
						"support_eip":           "false",
						"renewal_status":        "NotRenewal",
						"renewal_duration":      NOSET,
						"renewal_duration_unit": NOSET,
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period", "max_eip_tps", "queue_capacity"},
			},
		},
	})
}
func TestAccAlicloudAmqpInstance_vip(t *testing.T) {

	var v map[string]interface{}
	resourceId := "alicloud_amqp_instance.default"
	ra := resourceAttrInit(resourceId, AmqpInstanceBasicMap)
	serviceFunc := func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-AmqpInstancevip%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAmqpInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type":  "vip",
					"max_tps":        "5000",
					"payment_type":   "Subscription",
					"period":         "1",
					"queue_capacity": "50",
					"storage_size":   "700",
					"support_eip":    "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":         "vip",
						"max_tps":               "5000",
						"payment_type":          "Subscription",
						"queue_capacity":        "50",
						"storage_size":          "700",
						"support_eip":           "false",
						"renewal_duration":      NOSET,
						"renewal_duration_unit": NOSET,
						"renewal_status":        "ManualRenewal",
						"status":                "SERVING",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modify_type": "Upgrade",
					"max_tps":     "10000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_tps": "10000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"queue_capacity": "55",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"queue_capacity": "55",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_size": "800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_size": "800",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"support_eip": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"support_eip": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_eip_tps": "256",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_eip_tps": "256",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renewal_duration":      "1",
					"renewal_duration_unit": "Month",
					"renewal_status":        "AutoRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_duration":      "1",
						"renewal_duration_unit": "Month",
						"renewal_status":        "AutoRenewal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modify_type":    "Downgrade",
					"max_tps":        "5000",
					"queue_capacity": "50",
					"storage_size":   "700",
					"support_eip":    "false",
					"renewal_status": "NotRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_tps":               "5000",
						"queue_capacity":        "50",
						"storage_size":          "700",
						"support_eip":           "false",
						"renewal_status":        "NotRenewal",
						"renewal_duration":      NOSET,
						"renewal_duration_unit": NOSET,
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period", "max_eip_tps", "queue_capacity", "storage_size"},
			},
		},
	})
}

func resourceAmqpInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
 			default = "%v"
		}
		`, name)
}

var AmqpInstanceBasicMap = map[string]string{}
