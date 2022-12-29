package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Skip Test,Because each account can only be opened once
func SkipTestAccAlicloudThreatDetectionInstance_basic1826(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionInstanceMap1826)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sThreatDetectionInstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionInstanceBasicDependence1826)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":           "Subscription",
					"period":                 "12",
					"renewal_status":         "ManualRenewal",
					"sas_sls_storage":        "100",
					"sas_anti_ransomware":    "100",
					"container_image_scan":   "100",
					"sas_webguard_order_num": "100",
					"sas_sc":                 "true",
					"version_code":           "level2",
					"buy_number":             "30",
					"honeypot_switch":        "1",
					"sas_sdk_switch":         "1",
					"sas_sdk":                "1000",
					"honeypot":               "32",
					"v_core":                 "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":           "Subscription",
						"period":                 "12",
						"renewal_status":         "ManualRenewal",
						"sas_sls_storage":        "100",
						"sas_anti_ransomware":    "100",
						"container_image_scan":   "100",
						"sas_webguard_order_num": "100",
						"sas_sc":                 "true",
						"version_code":           "level2",
						"buy_number":             "30",
						"honeypot_switch":        "1",
						"sas_sdk_switch":         "1",
						"sas_sdk":                "1000",
						"honeypot":               "32",
						"v_core":                 "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sas_sls_storage":        "120",
					"sas_anti_ransomware":    "120",
					"sas_webguard_order_num": "120",
					"sas_sc":                 "true",
					"version_code":           "level2",
					"buy_number":             "30",
					"modify_type":            "Upgrade",
					"honeypot_switch":        "1",
					"sas_sdk_switch":         "1",
					"payment_type":           "Subscription",
					"sas_sdk":                "1200",
					"honeypot":               "32",
					"v_core":                 "120",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sas_sls_storage":        "120",
						"sas_anti_ransomware":    "120",
						"sas_webguard_order_num": "120",
						"sas_sc":                 "true",
						"version_code":           "level2",
						"buy_number":             "30",
						"modify_type":            "Upgrade",
						"honeypot_switch":        "1",
						"sas_sdk_switch":         "1",
						"payment_type":           "Subscription",
						"sas_sdk":                "1200",
						"honeypot":               "32",
						"v_core":                 "120",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renew_period":        "1",
					"renewal_period_unit": "M",
					"renewal_status":      "AutoRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renew_period":        "1",
						"renewal_period_unit": "M",
						"renewal_status":      "AutoRenewal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renewal_status": "ManualRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_status": "ManualRenewal",
					}),
				),
			},
		},
	})
}

var AlicloudThreatDetectionInstanceMap1826 = map[string]string{}

func AlicloudThreatDetectionInstanceBasicDependence1826(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

`, name)
}
