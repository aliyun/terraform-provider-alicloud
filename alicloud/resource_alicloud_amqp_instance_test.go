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
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":  "${var.name}",
					"instance_type":  "professional",
					"max_tps":        "1000",
					"payment_type":   "Subscription",
					"period":         "1",
					"queue_capacity": "50",
					"support_eip":    "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":         name,
						"instance_type":         "professional",
						"max_tps":               "1000",
						"payment_type":          "Subscription",
						"queue_capacity":        "50",
						"support_eip":           "false",
						"renewal_status":        "ManualRenewal",
						"renewal_duration_unit": "",
						"status":                "SERVING",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "-update",
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
					"max_eip_tps": "128",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"support_eip": "true",
						"max_eip_tps": "128",
					}),
				),
			},
			// There is an OpenAPI bug that the api return renewal_duration_unit is ""
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"renewal_duration":      "1",
			//		"renewal_duration_unit": "Month",
			//		"renewal_status":        "AutoRenewal",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"renewal_duration":      "1",
			//			"renewal_duration_unit": "Month",
			//			"renewal_status":        "AutoRenewal",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":  name,
					"modify_type":    "Downgrade",
					"max_tps":        "1000",
					"queue_capacity": "50",
					"support_eip":    "false",
					"renewal_status": "NotRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":         name,
						"max_tps":               "1000",
						"queue_capacity":        "50",
						"support_eip":           "false",
						"renewal_status":        "NotRenewal",
						"renewal_duration":      NOSET,
						"renewal_duration_unit": "",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period", "max_tps", "max_eip_tps", "queue_capacity"},
			},
		},
	})
}

// Currently, the test account does not support the vip
func SkipTestAccAlicloudAmqpInstance_vip(t *testing.T) {

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
		CheckDestroy:  nil,
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
						"renewal_status":        "ManualRenewal",
						"renewal_duration":      "0",
						"renewal_duration_unit": "",
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
					"max_eip_tps": "128",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"support_eip": "true",
						"max_eip_tps": "128",
					}),
				),
			},
			// There is an OpenAPI bug that the api return renewal_duration_unit is ""
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"renewal_duration":      "1",
			//		"renewal_duration_unit": "Month",
			//		"renewal_status":        "AutoRenewal",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"renewal_duration":      "1",
			//			"renewal_duration_unit": "Month",
			//			"renewal_status":        "AutoRenewal",
			//		}),
			//	),
			//},
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
						"renewal_duration":      "0",
						"renewal_duration_unit": "",
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
