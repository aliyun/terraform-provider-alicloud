package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ThreatDetection OssScanConfig. >>> Resource test cases, automatically generated.
// Case OssScanConfig_20250211 10184
func TestAccAliCloudThreatDetectionOssScanConfig_basic10184(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_oss_scan_config.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionOssScanConfigMap10184)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionOssScanConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionOssScanConfigBasicDependence10184)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
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
					"key_prefix_list": []string{
						"/root", "/usr", "/123"},
					"oss_scan_config_name": name,
					"end_time":             "00:00:01",
					"start_time":           "00:00:00",
					"enable":               "0",
					"key_suffix_list": []string{
						".html", ".php", ".k"},
					"scan_day_list": []string{
						"1", "2", "4", "3"},
					"all_key_prefix": "false",
					"decryption_list": []string{
						"OSS", "AES", "TEST", "TEST2"},
					"decompress_max_file_count": "2",
					"decompress_max_layer":      "10",
					"last_modified_start_time":  "1721713320000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name_list.#":        "3",
						"key_prefix_list.#":         "3",
						"oss_scan_config_name":      name,
						"end_time":                  "00:00:01",
						"start_time":                "00:00:00",
						"enable":                    "0",
						"key_suffix_list.#":         "3",
						"scan_day_list.#":           "4",
						"all_key_prefix":            "false",
						"decryption_list.#":         "4",
						"decompress_max_file_count": "2",
						"decompress_max_layer":      "10",
						"last_modified_start_time":  "1721713320000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name_list": []string{
						"gcx-test-oss-74"},
					"key_prefix_list":      []string{},
					"oss_scan_config_name": name + "_update",
					"end_time":             "00:00:02",
					"start_time":           "00:00:01",
					"enable":               "1",
					"key_suffix_list": []string{
						".jsp"},
					"scan_day_list": []string{
						"2", "5"},
					"all_key_prefix": "true",
					"decryption_list": []string{
						"KMS"},
					"decompress_max_file_count": "1",
					"decompress_max_layer":      "1",
					"last_modified_start_time":  "1721713532000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name_list.#":        "1",
						"key_prefix_list.#":         "0",
						"oss_scan_config_name":      name + "_update",
						"end_time":                  "00:00:02",
						"start_time":                "00:00:01",
						"enable":                    "1",
						"key_suffix_list.#":         "1",
						"scan_day_list.#":           "2",
						"all_key_prefix":            "true",
						"decryption_list.#":         "1",
						"decompress_max_file_count": "1",
						"decompress_max_layer":      "1",
						"last_modified_start_time":  "1721713532000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name_list": []string{
						"gcx-test-oss-74"},
					"key_prefix_list": []string{
						"/123"},
					"start_time": "00:00:08",
					"key_suffix_list": []string{
						".jsp"},
					"scan_day_list": []string{
						"2", "5"},
					"decryption_list":           []string{},
					"decompress_max_file_count": REMOVEKEY,
					"last_modified_start_time":  REMOVEKEY,
					"decompress_max_layer":      REMOVEKEY,
					"all_key_prefix":            REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name_list.#":        "1",
						"key_prefix_list.#":         "1",
						"start_time":                "00:00:08",
						"key_suffix_list.#":         "1",
						"scan_day_list.#":           "2",
						"decryption_list.#":         "0",
						"decompress_max_file_count": REMOVEKEY,
						"last_modified_start_time":  REMOVEKEY,
						"decompress_max_layer":      REMOVEKEY,
						"all_key_prefix":            REMOVEKEY,
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

var AlicloudThreatDetectionOssScanConfigMap10184 = map[string]string{}

func AlicloudThreatDetectionOssScanConfigBasicDependence10184(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "default8j4t1R" {
  storage_class = "Standard"
}

resource "alicloud_oss_bucket" "default9HMqfT" {
  storage_class = "Standard"
}

resource "alicloud_oss_bucket" "defaultxBXqFQ" {
  storage_class = "Standard"
}

resource "alicloud_oss_bucket" "defaulthZvCmR" {
  storage_class = "Standard"
}


`, name)
}

// Case create 5010
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
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionOssScanConfigBasicDependence5010)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
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
					"key_prefix_list": []string{
						"/root", "/usr", "/123"},
					"oss_scan_config_name": name,
					"end_time":             "00:00:01",
					"start_time":           "00:00:00",
					"enable":               "0",
					"key_suffix_list": []string{
						".html", ".php", ".k"},
					"scan_day_list": []string{
						"1", "2", "4", "3"},
					"all_key_prefix": "false",
					"decryption_list": []string{
						"OSS", "AES", "TEST", "TEST2"},
					"decompress_max_file_count": "2",
					"decompress_max_layer":      "10",
					"last_modified_start_time":  "1721713320000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name_list.#":        "3",
						"key_prefix_list.#":         "3",
						"oss_scan_config_name":      name,
						"end_time":                  "00:00:01",
						"start_time":                "00:00:00",
						"enable":                    "0",
						"key_suffix_list.#":         "3",
						"scan_day_list.#":           "4",
						"all_key_prefix":            "false",
						"decryption_list.#":         "4",
						"decompress_max_file_count": "2",
						"decompress_max_layer":      "10",
						"last_modified_start_time":  "1721713320000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name_list": []string{
						"gcx-test-oss-74"},
					"key_prefix_list":      []string{},
					"oss_scan_config_name": name + "_update",
					"end_time":             "00:00:02",
					"start_time":           "00:00:01",
					"enable":               "1",
					"key_suffix_list": []string{
						".jsp"},
					"scan_day_list": []string{
						"2", "5"},
					"all_key_prefix": "true",
					"decryption_list": []string{
						"KMS"},
					"decompress_max_file_count": "1",
					"decompress_max_layer":      "1",
					"last_modified_start_time":  "1721713532000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name_list.#":        "1",
						"key_prefix_list.#":         "0",
						"oss_scan_config_name":      name + "_update",
						"end_time":                  "00:00:02",
						"start_time":                "00:00:01",
						"enable":                    "1",
						"key_suffix_list.#":         "1",
						"scan_day_list.#":           "2",
						"all_key_prefix":            "true",
						"decryption_list.#":         "1",
						"decompress_max_file_count": "1",
						"decompress_max_layer":      "1",
						"last_modified_start_time":  "1721713532000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name_list": []string{
						"gcx-test-oss-74"},
					"key_prefix_list": []string{
						"/123"},
					"start_time": "00:00:08",
					"key_suffix_list": []string{
						".jsp"},
					"scan_day_list": []string{
						"2", "5"},
					"decryption_list":           []string{},
					"decompress_max_file_count": REMOVEKEY,
					"last_modified_start_time":  REMOVEKEY,
					"decompress_max_layer":      REMOVEKEY,
					"all_key_prefix":            REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name_list.#":        "1",
						"key_prefix_list.#":         "1",
						"start_time":                "00:00:08",
						"key_suffix_list.#":         "1",
						"scan_day_list.#":           "2",
						"decryption_list.#":         "0",
						"decompress_max_file_count": REMOVEKEY,
						"last_modified_start_time":  REMOVEKEY,
						"decompress_max_layer":      REMOVEKEY,
						"all_key_prefix":            REMOVEKEY,
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
  storage_class = "Standard"
}

resource "alicloud_oss_bucket" "default9HMqfT" {
  storage_class = "Standard"
}

resource "alicloud_oss_bucket" "defaultxBXqFQ" {
  storage_class = "Standard"
}

resource "alicloud_oss_bucket" "defaulthZvCmR" {
  storage_class = "Standard"
}


`, name)
}

// Case create_副本1724985158441 7682
func TestAccAliCloudThreatDetectionOssScanConfig_basic7682(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_oss_scan_config.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionOssScanConfigMap7682)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionOssScanConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionOssScanConfigBasicDependence7682)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
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
					"key_prefix_list": []string{
						"/root", "/usr", "/123"},
					"oss_scan_config_name": name,
					"end_time":             "00:00:01",
					"start_time":           "00:00:00",
					"enable":               "0",
					"key_suffix_list": []string{
						".html", ".php", ".k"},
					"scan_day_list": []string{
						"1", "2", "4", "3"},
					"all_key_prefix": "false",
					"decryption_list": []string{
						"OSS", "AES", "TEST", "TEST2"},
					"decompress_max_file_count": "2",
					"decompress_max_layer":      "10",
					"last_modified_start_time":  "1721713320000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name_list.#":        "3",
						"key_prefix_list.#":         "3",
						"oss_scan_config_name":      name,
						"end_time":                  "00:00:01",
						"start_time":                "00:00:00",
						"enable":                    "0",
						"key_suffix_list.#":         "3",
						"scan_day_list.#":           "4",
						"all_key_prefix":            "false",
						"decryption_list.#":         "4",
						"decompress_max_file_count": "2",
						"decompress_max_layer":      "10",
						"last_modified_start_time":  "1721713320000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name_list": []string{
						"gcx-test-oss-74"},
					"key_prefix_list":      []string{},
					"oss_scan_config_name": name + "_update",
					"end_time":             "00:00:02",
					"start_time":           "00:00:01",
					"enable":               "1",
					"key_suffix_list": []string{
						".jsp"},
					"scan_day_list": []string{
						"2", "5"},
					"all_key_prefix": "true",
					"decryption_list": []string{
						"KMS"},
					"decompress_max_file_count": "1",
					"decompress_max_layer":      "1",
					"last_modified_start_time":  "1721713532000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name_list.#":        "1",
						"key_prefix_list.#":         "0",
						"oss_scan_config_name":      name + "_update",
						"end_time":                  "00:00:02",
						"start_time":                "00:00:01",
						"enable":                    "1",
						"key_suffix_list.#":         "1",
						"scan_day_list.#":           "2",
						"all_key_prefix":            "true",
						"decryption_list.#":         "1",
						"decompress_max_file_count": "1",
						"decompress_max_layer":      "1",
						"last_modified_start_time":  "1721713532000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name_list": []string{
						"gcx-test-oss-74"},
					"key_prefix_list": []string{
						"/123"},
					"start_time": "00:00:08",
					"key_suffix_list": []string{
						".jsp"},
					"scan_day_list": []string{
						"2", "5"},
					"decompress_max_file_count": REMOVEKEY,
					"last_modified_start_time":  REMOVEKEY,
					"decompress_max_layer":      REMOVEKEY,
					"all_key_prefix":            REMOVEKEY,
					"decryption_list":           CLEARLIST,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name_list.#":        "1",
						"key_prefix_list.#":         "1",
						"start_time":                "00:00:08",
						"key_suffix_list.#":         "1",
						"scan_day_list.#":           "2",
						"decompress_max_file_count": REMOVEKEY,
						"last_modified_start_time":  REMOVEKEY,
						"decompress_max_layer":      REMOVEKEY,
						"all_key_prefix":            REMOVEKEY,
						"decryption_list.#":         "0",
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

var AlicloudThreatDetectionOssScanConfigMap7682 = map[string]string{}

func AlicloudThreatDetectionOssScanConfigBasicDependence7682(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "default8j4t1R" {
  storage_class = "Standard"
}

resource "alicloud_oss_bucket" "default9HMqfT" {
  storage_class = "Standard"
}

resource "alicloud_oss_bucket" "defaultxBXqFQ" {
  storage_class = "Standard"
}

resource "alicloud_oss_bucket" "defaulthZvCmR" {
  storage_class = "Standard"
}


`, name)
}

// Case OssScanConfig 9153
func TestAccAliCloudThreatDetectionOssScanConfig_basic9153(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_oss_scan_config.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionOssScanConfigMap9153)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionOssScanConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionOssScanConfigBasicDependence9153)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
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
					"key_prefix_list": []string{
						"/root", "/usr", "/123"},
					"oss_scan_config_name": name,
					"end_time":             "00:00:01",
					"start_time":           "00:00:00",
					"enable":               "0",
					"key_suffix_list": []string{
						".html", ".php", ".k"},
					"scan_day_list": []string{
						"1", "2", "4", "3"},
					"all_key_prefix": "false",
					"decryption_list": []string{
						"OSS", "AES", "TEST", "TEST2"},
					"decompress_max_file_count": "2",
					"decompress_max_layer":      "10",
					"last_modified_start_time":  "1721713320000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name_list.#":        "3",
						"key_prefix_list.#":         "3",
						"oss_scan_config_name":      name,
						"end_time":                  "00:00:01",
						"start_time":                "00:00:00",
						"enable":                    "0",
						"key_suffix_list.#":         "3",
						"scan_day_list.#":           "4",
						"all_key_prefix":            "false",
						"decryption_list.#":         "4",
						"decompress_max_file_count": "2",
						"decompress_max_layer":      "10",
						"last_modified_start_time":  "1721713320000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name_list": []string{
						"gcx-test-oss-74"},
					"key_prefix_list":      []string{},
					"oss_scan_config_name": name + "_update",
					"end_time":             "00:00:02",
					"start_time":           "00:00:01",
					"enable":               "1",
					"key_suffix_list": []string{
						".jsp"},
					"scan_day_list": []string{
						"2", "5"},
					"all_key_prefix": "true",
					"decryption_list": []string{
						"KMS"},
					"decompress_max_file_count": "1",
					"decompress_max_layer":      "1",
					"last_modified_start_time":  "1721713532000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name_list.#":        "1",
						"key_prefix_list.#":         "0",
						"oss_scan_config_name":      name + "_update",
						"end_time":                  "00:00:02",
						"start_time":                "00:00:01",
						"enable":                    "1",
						"key_suffix_list.#":         "1",
						"scan_day_list.#":           "2",
						"all_key_prefix":            "true",
						"decryption_list.#":         "1",
						"decompress_max_file_count": "1",
						"decompress_max_layer":      "1",
						"last_modified_start_time":  "1721713532000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name_list": []string{
						"gcx-test-oss-74"},
					"key_prefix_list": []string{
						"/123"},
					"start_time": "00:00:08",
					"key_suffix_list": []string{
						".jsp"},
					"scan_day_list": []string{
						"2", "5"},
					"decryption_list":           []string{},
					"decompress_max_file_count": REMOVEKEY,
					"last_modified_start_time":  REMOVEKEY,
					"decompress_max_layer":      REMOVEKEY,
					"all_key_prefix":            REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name_list.#":        "1",
						"key_prefix_list.#":         "1",
						"start_time":                "00:00:08",
						"key_suffix_list.#":         "1",
						"scan_day_list.#":           "2",
						"decryption_list.#":         "0",
						"decompress_max_file_count": REMOVEKEY,
						"last_modified_start_time":  REMOVEKEY,
						"decompress_max_layer":      REMOVEKEY,
						"all_key_prefix":            REMOVEKEY,
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

var AlicloudThreatDetectionOssScanConfigMap9153 = map[string]string{}

func AlicloudThreatDetectionOssScanConfigBasicDependence9153(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "default8j4t1R" {
  storage_class = "Standard"
}

resource "alicloud_oss_bucket" "default9HMqfT" {
  storage_class = "Standard"
}

resource "alicloud_oss_bucket" "defaultxBXqFQ" {
  storage_class = "Standard"
}

resource "alicloud_oss_bucket" "defaulthZvCmR" {
  storage_class = "Standard"
}


`, name)
}

// Test ThreatDetection OssScanConfig. <<< Resource test cases, automatically generated.
