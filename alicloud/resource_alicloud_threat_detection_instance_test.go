package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Skip Test,Because each account can only be opened once
func TestAccAliCloudThreatDetectionInstance_basic1826(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionInstanceMap1826)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sThreatDetectionInstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionInstanceBasicDependence1826)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckForCleanUpInstances(t, "", "sas", "sas", "sas", "")
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
					"buy_number":             "55",
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
						"buy_number":             "55",
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
					"container_image_scan":   "200",
					"sas_sc":                 "true",
					"version_code":           "level2",
					"buy_number":             "60",
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
						"container_image_scan":   "200",
						"sas_sc":                 "true",
						"version_code":           "level2",
						"buy_number":             "60",
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

var AliCloudThreatDetectionInstanceMap1826 = map[string]string{}

func AliCloudThreatDetectionInstanceBasicDependence1826(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

`, name)
}

func TestAccAliCloudThreatDetectionInstance_basic4253(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionInstanceMap4253)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sthreatdetectioninstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionInstanceBasicDependence4253)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckForCleanUpInstances(t, "", "sas", "sas", "sas", "")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":   "Subscription",
					"version_code":   "level3",
					"period":         "12",
					"buy_number":     "55",
					"renewal_status": "ManualRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":   "Subscription",
						"version_code":   "level3",
						"period":         "12",
						"buy_number":     "55",
						"renewal_status": "ManualRenewal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vul_switch":  "0",
					"modify_type": "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vul_switch": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sas_sls_storage": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sas_sls_storage": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"threat_analysis_switch": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"threat_analysis_switch": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"v_core": "55",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"v_core": "55",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sas_sc": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sas_cspm_switch": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sas_cspm_switch": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sas_webguard_boolean": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sas_webguard_boolean": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"honeypot_switch": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"honeypot_switch": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sas_sdk": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sas_sdk": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sas_anti_ransomware": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sas_anti_ransomware": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sas_webguard_order_num": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sas_webguard_order_num": "0",
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
			{
				Config: testAccConfig(map[string]interface{}{
					"rasp_count": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rasp_count": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vul_count": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vul_count": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sas_cspm": "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sas_cspm": "1000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sas_sdk_switch": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sas_sdk_switch": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renewal_period_unit": "M",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_period_unit": "M",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"container_image_scan_new": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"container_image_scan_new": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"honeypot": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"honeypot": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"version_code": "level3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version_code": "level3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"threat_analysis": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"threat_analysis": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sas_sls_storage": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sas_sls_storage": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"threat_analysis_switch": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"threat_analysis_switch": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"v_core": "58",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"v_core": "58",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sas_sc":          "true",
					"sas_cspm_switch": "1",
					"sas_cspm":        "15000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sas_sc":          "true",
						"sas_cspm_switch": "1",
						"sas_cspm":        "15000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"buy_number": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"buy_number": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sas_webguard_boolean": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sas_webguard_boolean": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"honeypot_switch": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"honeypot_switch": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sas_sdk": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sas_sdk": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sas_anti_ransomware": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sas_anti_ransomware": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sas_webguard_order_num": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sas_webguard_order_num": "20",
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
					"rasp_count": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rasp_count": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vul_count": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vul_count": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"version_code": "level2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version_code": "level2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sas_cspm": "70000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sas_cspm": "70000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sas_sdk_switch": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sas_sdk_switch": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renewal_period_unit": "Y",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_period_unit": "Y",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"container_image_scan_new": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"container_image_scan_new": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"honeypot": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"honeypot": "30",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period"},
			},
		},
	})
}

var AliCloudThreatDetectionInstanceMap4253 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AliCloudThreatDetectionInstanceBasicDependence4253(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 4253  twin
func TestAccAliCloudThreatDetectionInstance_basic4253_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionInstanceMap4253)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sthreatdetectioninstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionInstanceBasicDependence4253)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			// testAccPreCheckForCleanUpInstances(t, "", "sas", "sas", "sas", "")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"threat_analysis":          "20",
					"sas_sls_storage":          "20",
					"threat_analysis_switch":   "1",
					"v_core":                   "20",
					"sas_sc":                   "true",
					"sas_cspm_switch":          "1",
					"buy_number":               "60",
					"sas_webguard_boolean":     "1",
					"honeypot_switch":          "1",
					"payment_type":             "Subscription",
					"sas_sdk":                  "20",
					"sas_anti_ransomware":      "20",
					"sas_webguard_order_num":   "20",
					"renewal_status":           "AutoRenewal",
					"period":                   "1",
					"vul_switch":               "0",
					"rasp_count":               "2",
					"vul_count":                "30",
					"version_code":             "level2",
					"sas_cspm":                 "15000",
					"sas_sdk_switch":           "1",
					"renewal_period_unit":      "Y",
					"container_image_scan_new": "200",
					"honeypot":                 "30",
					"renew_period":             "2",
					"modify_type":              "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"threat_analysis":          "20",
						"sas_sls_storage":          "20",
						"threat_analysis_switch":   "1",
						"v_core":                   "20",
						"sas_sc":                   "true",
						"sas_cspm_switch":          "1",
						"buy_number":               "60",
						"sas_webguard_boolean":     "1",
						"honeypot_switch":          "1",
						"payment_type":             "Subscription",
						"sas_sdk":                  "20",
						"sas_anti_ransomware":      "20",
						"sas_webguard_order_num":   "20",
						"renewal_status":           "AutoRenewal",
						"period":                   "1",
						"vul_switch":               "0",
						"rasp_count":               "2",
						"vul_count":                "30",
						"version_code":             "level2",
						"sas_cspm":                 "15000",
						"sas_sdk_switch":           "1",
						"renewal_period_unit":      "Y",
						"container_image_scan_new": "200",
						"honeypot":                 "30",
						"renew_period":             "2",
						"modify_type":              "Upgrade",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period"},
			},
		},
	})
}

func TestAccAliCloudThreatDetectionInstance_basic4253_intl(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionInstanceMap4253)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sthreatdetectioninstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionInstanceBasicDependence4253)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			//testAccPreCheckWithAccountSiteType(t, IntlSite)
			//testAccPreCheckForCleanUpInstances(t, "", "sas", "sas", "sas", "")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"sas_sls_storage":             "20",
					"sas_cspm_switch":             "1",
					"buy_number":                  "40",
					"sas_webguard_boolean":        "1",
					"honeypot_switch":             "1",
					"payment_type":                "Subscription",
					"sas_sdk":                     "20",
					"sas_anti_ransomware":         "20",
					"sas_webguard_order_num":      "20",
					"renewal_status":              "AutoRenewal",
					"period":                      "1",
					"vul_switch":                  "0",
					"rasp_count":                  "2",
					"vul_count":                   "30",
					"version_code":                "level2",
					"sas_cspm":                    "15000",
					"sas_sdk_switch":              "1",
					"renewal_period_unit":         "Y",
					"container_image_scan_new":    "200",
					"honeypot":                    "20",
					"renew_period":                "2",
					"modify_type":                 "Upgrade",
					"threat_analysis_flow":        "10",
					"threat_analysis_sls_storage": "20",
					"threat_analysis_switch1":     "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sas_sls_storage":             "20",
						"sas_cspm_switch":             "1",
						"buy_number":                  "40",
						"sas_webguard_boolean":        "1",
						"honeypot_switch":             "1",
						"payment_type":                "Subscription",
						"sas_sdk":                     "20",
						"sas_anti_ransomware":         "20",
						"sas_webguard_order_num":      "20",
						"renewal_status":              "AutoRenewal",
						"period":                      "1",
						"vul_switch":                  "0",
						"rasp_count":                  "2",
						"vul_count":                   "30",
						"version_code":                "level2",
						"sas_cspm":                    "15000",
						"sas_sdk_switch":              "1",
						"renewal_period_unit":         "Y",
						"container_image_scan_new":    "200",
						"honeypot":                    "20",
						"renew_period":                "2",
						"modify_type":                 "Upgrade",
						"threat_analysis_flow":        "10",
						"threat_analysis_sls_storage": "20",
						"threat_analysis_switch1":     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"honeypot": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"honeypot": "30",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period"},
			},
		},
	})
}

func TestAccAliCloudThreatDetectionInstance_basic4253_twin_fix(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionInstanceMap4253)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sthreatdetectioninstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionInstanceBasicDependence4253)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckForCleanUpInstances(t, "", "sas", "sas", "sas", "")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":   "Subscription",
					"period":         "12",
					"buy_number":     "40",
					"renewal_status": "ManualRenewal",
					"vul_switch":     "1",
					"vul_count":      "20",
					"version_code":   "level7",
					"v_core":         "52",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":   "Subscription",
						"period":         "12",
						"buy_number":     "40",
						"renewal_status": "ManualRenewal",
						"vul_switch":     "1",
						"vul_count":      "20",
						"version_code":   "level7",
						"v_core":         "52",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vul_switch":   "0",
					"vul_count":    "0",
					"version_code": "level3",
					"modify_type":  "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vul_count":    "0",
						"version_code": "level3",
						"vul_switch":   "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period"},
			},
		},
	})
}

// Test ThreatDetection Instance. >>> Resource test cases, automatically generated.
// Case 中国站资源测试用例_20250217_后付费 10249 适配废弃字段post_pay_module_switch
func TestAccAliCloudThreatDetectionInstance_basic10249(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionInstanceMap10249)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionInstanceBasicDependence10249)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":   "PayAsYouGo",
					"post_paid_flag": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":   "PayAsYouGo",
						"post_paid_flag": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"subscription_type":      "PayAsYouGo",
					"post_pay_module_switch": "{\\\"VUL\\\":1}",
					"modify_type":            "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"subscription_type":      "PayAsYouGo",
						"post_pay_module_switch": CHECKSET,
						"modify_type":            "Upgrade",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"post_paid_host_auto_bind":         "1",
					"post_paid_host_auto_bind_version": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"post_paid_host_auto_bind":         "1",
						"post_paid_host_auto_bind_version": "7",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period", "post_paid_flag", "post_pay_module_switch", "product_type", "subscription_type"},
			},
		},
	})
}

func TestAccAliCloudThreatDetectionInstance_basic10249_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionInstanceMap10249)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionInstanceBasicDependence10249)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":                     "PayAsYouGo",
					"post_paid_flag":                   "1",
					"post_paid_host_auto_bind":         "1",
					"post_paid_host_auto_bind_version": "7",
					"post_pay_module_switch":           "{\\\"VUL\\\":1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":                     "PayAsYouGo",
						"post_paid_flag":                   "1",
						"post_paid_host_auto_bind":         "1",
						"post_paid_host_auto_bind_version": "7",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period", "post_paid_flag", "post_pay_module_switch", "product_type", "subscription_type"},
			},
		},
	})
}

var AliCloudThreatDetectionInstanceMap10249 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AliCloudThreatDetectionInstanceBasicDependence10249(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 中国站资源测试用例_20250217_后付费New 12300
func TestAccAliCloudThreatDetectionInstance_basic12300(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionInstanceMap12300)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionInstanceBasicDependence12300)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":   "PayAsYouGo",
					"post_paid_flag": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":   "PayAsYouGo",
						"post_paid_flag": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"post_pay_module_switch_obj": []map[string]interface{}{
						{
							"vul":             "0",
							"cspm":            "0",
							"agentless":       "0",
							"serverless":      "0",
							"ctdr":            "0",
							"rasp":            "0",
							"sdk":             "1",
							"post_host":       "0",
							"ctdr_storage":    "0",
							"anti_ransomware": "0",
							"basic_service":   "1",
							"web_lock":        "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"post_pay_module_switch_obj.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"post_pay_module_switch_obj": []map[string]interface{}{
						{
							"vul":             "1",
							"cspm":            "1",
							"agentless":       "1",
							"serverless":      "1",
							"ctdr":            "1",
							"rasp":            "1",
							"sdk":             "0",
							"post_host":       "1",
							"ctdr_storage":    "1",
							"anti_ransomware": "1",
							"basic_service":   "1",
							"web_lock":        "1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"post_pay_module_switch_obj.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"post_paid_host_auto_bind":         "1",
					"post_paid_host_auto_bind_version": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"post_paid_host_auto_bind":         "1",
						"post_paid_host_auto_bind_version": "7",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period", "post_paid_flag", "post_pay_module_switch", "subscription_type"},
			},
		},
	})
}

func TestAccAliCloudThreatDetectionInstance_basic12300_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionInstanceMap12300)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionInstanceBasicDependence12300)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":                     "PayAsYouGo",
					"post_paid_flag":                   "1",
					"post_paid_host_auto_bind":         "1",
					"post_paid_host_auto_bind_version": "7",
					"post_pay_module_switch_obj": []map[string]interface{}{
						{
							"vul":             "0",
							"cspm":            "0",
							"agentless":       "0",
							"serverless":      "0",
							"ctdr":            "0",
							"rasp":            "0",
							"sdk":             "1",
							"post_host":       "1",
							"ctdr_storage":    "0",
							"anti_ransomware": "0",
							"basic_service":   "1",
							"web_lock":        "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":                     "PayAsYouGo",
						"post_paid_flag":                   "1",
						"post_paid_host_auto_bind":         "1",
						"post_paid_host_auto_bind_version": "7",
						"post_pay_module_switch_obj.#":     "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period", "post_paid_flag", "post_pay_module_switch", "subscription_type"},
			},
		},
	})
}

// Case 国际站资源测试用例_20250217_后付费Intl 12300
func TestAccAliCloudThreatDetectionInstance_basic12300_intl(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionInstanceMap12300)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionInstanceBasicDependence12300)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":                     "PayAsYouGo",
					"post_paid_flag":                   "1",
					"post_paid_host_auto_bind":         "1",
					"post_paid_host_auto_bind_version": "7",
					"post_pay_module_switch_obj": []map[string]interface{}{
						{
							"vul":             "0",
							"cspm":            "0",
							"agentless":       "0",
							"serverless":      "0",
							"ctdr":            "0",
							"rasp":            "0",
							"sdk":             "1",
							"post_host":       "1",
							"ctdr_storage":    "0",
							"anti_ransomware": "0",
							"basic_service":   "1",
							"web_lock":        "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":                     "PayAsYouGo",
						"post_paid_flag":                   "1",
						"post_paid_host_auto_bind":         "1",
						"post_paid_host_auto_bind_version": "7",
						"post_pay_module_switch_obj.#":     "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "period", "post_paid_flag", "post_pay_module_switch", "subscription_type"},
			},
		},
	})
}

var AliCloudThreatDetectionInstanceMap12300 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AliCloudThreatDetectionInstanceBasicDependence12300(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test ThreatDetection Instance. <<< Resource test cases, automatically generated.
