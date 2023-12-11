package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ThreatDetection OssScanConfig. >>> Resource test cases, automatically generated.
// Case 5010
func TestAccAliCloudThreatDetectionOssScanConfig_basic5010(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_oss_scan_config.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionOssScanConfigMap5010)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionOssScanConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sthreatdetectionossscanconfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionOssScanConfigBasicDependence5010)
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
					"bucket_name_list": []string{
						"gcx-test-oss-71", "gcx-test-oss-72", "gcx-test-oss-73"},
					"end_time":   "00:00:01",
					"start_time": "00:00:00",
					"enable":     "0",
					"key_suffix_list": []string{
						".html", ".php", ".k"},
					"scan_day_list": []string{
						"1", "2", "4", "3"},
					"oss_scan_config_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name_list.#":   "3",
						"end_time":             "00:00:01",
						"start_time":           "00:00:00",
						"enable":               "0",
						"key_suffix_list.#":    "3",
						"scan_day_list.#":      "4",
						"oss_scan_config_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"all_key_prefix": "false",
					"key_prefix_list": []string{
						"/root", "/usr", "/123"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"all_key_prefix":    "false",
						"key_prefix_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name_list": []string{
						"gcx-test-oss-71", "gcx-test-oss-72", "gcx-test-oss-73"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"end_time": "00:00:01",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"end_time": "00:00:01",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"start_time": "00:00:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"start_time": "00:00:00",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_suffix_list": []string{
						".html", ".php", ".k"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_suffix_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scan_day_list": []string{
						"1", "2", "4", "3"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scan_day_list.#": "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"oss_scan_config_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oss_scan_config_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name_list": []string{
						"gcx-test-oss-74"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"oss_scan_config_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oss_scan_config_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"end_time": "00:00:02",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"end_time": "00:00:02",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"start_time": "00:00:01",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"start_time": "00:00:01",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_suffix_list": []string{
						".jsp"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_suffix_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scan_day_list": []string{
						"2", "5"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scan_day_list.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"all_key_prefix":  "true",
					"key_prefix_list": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"all_key_prefix":    "true",
						"key_prefix_list.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name_list": []string{
						"gcx-test-oss-71", "gcx-test-oss-72", "gcx-test-oss-73"},
					"key_prefix_list": []string{
						"/root", "/usr", "/123"},
					"oss_scan_config_name": name + "_update",
					"end_time":             "00:00:01",
					"start_time":           "00:00:00",
					"enable":               "0",
					"key_suffix_list": []string{
						".html", ".php", ".k"},
					"scan_day_list": []string{
						"1", "2", "4", "3"},
					"all_key_prefix": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name_list.#":   "3",
						"key_prefix_list.#":    "3",
						"oss_scan_config_name": name + "_update",
						"end_time":             "00:00:01",
						"start_time":           "00:00:00",
						"enable":               "0",
						"key_suffix_list.#":    "3",
						"scan_day_list.#":      "4",
						"all_key_prefix":       "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudThreatDetectionOssScanConfigMap5010 = map[string]string{}

func AlicloudThreatDetectionOssScanConfigBasicDependence5010(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "default8j4t1R" {
  bucket = "${var.name}-1"

  storage_class = "Standard"
}

resource "alicloud_oss_bucket" "default9HMqfT" {
  bucket = "${var.name}-2"

  storage_class = "Standard"
}

resource "alicloud_oss_bucket" "defaultxBXqFQ" {
  bucket = "${var.name}-3"
  storage_class = "Standard"
}

resource "alicloud_oss_bucket" "defaulthZvCmR" {
  bucket = "${var.name}-4"

  storage_class = "Standard"
}


`, name)
}

// Case 5010  twin
func TestAccAliCloudThreatDetectionOssScanConfig_basic5010_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_oss_scan_config.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionOssScanConfigMap5010)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionOssScanConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sthreatdetectionossscanconfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionOssScanConfigBasicDependence5010)
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
					"bucket_name_list": []string{
						"gcx-test-oss-74", "gcx-test-oss-72", "gcx-test-oss-73"},
					"key_prefix_list": []string{
						"/root", "/usr", "/123"},
					"oss_scan_config_name": name,
					"end_time":             "00:00:02",
					"start_time":           "00:00:01",
					"enable":               "1",
					"key_suffix_list": []string{
						".jsp", ".php", ".k"},
					"scan_day_list": []string{
						"2", "5", "4", "3"},
					"all_key_prefix": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name_list.#":   "3",
						"key_prefix_list.#":    "3",
						"oss_scan_config_name": name,
						"end_time":             "00:00:02",
						"start_time":           "00:00:01",
						"enable":               "1",
						"key_suffix_list.#":    "3",
						"scan_day_list.#":      "4",
						"all_key_prefix":       "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Test ThreatDetection OssScanConfig. <<< Resource test cases, automatically generated.
