package threatdetection_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/threatdetection"
	tFacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ThreatDetection OssScanConfig. >>> Resource test cases, automatically generated.
// Case OssScanConfig_20250211 10184
func TestAccAliCloudThreatDetectionOssScanConfig_basic10184(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_oss_scan_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudThreatDetectionOssScanConfigMap10184)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return threatdetection.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeThreatDetectionOssScanConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tFacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudThreatDetectionOssScanConfigBasicDependence10184)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name_list": []string{
						"${alicloud_oss_bucket.default8j4t1R.id}", "${alicloud_oss_bucket.default9HMqfT.id}", "${alicloud_oss_bucket.defaultxBXqFQ.id}"},
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
						"${alicloud_oss_bucket.defaulthZvCmR.id}"},
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
						"${alicloud_oss_bucket.defaulthZvCmR.id}"},
					"key_prefix_list": []string{
						"/123"},
					"start_time": "00:00:08",
					"key_suffix_list": []string{
						".jsp"},
					"scan_day_list": []string{
						"2", "5"},
					"decryption_list":           []string{},
					"decompress_max_file_count": acctest.REMOVEKEY,
					"last_modified_start_time":  acctest.REMOVEKEY,
					"decompress_max_layer":      acctest.REMOVEKEY,
					"all_key_prefix":            acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name_list.#":        "1",
						"key_prefix_list.#":         "1",
						"start_time":                "00:00:08",
						"key_suffix_list.#":         "1",
						"scan_day_list.#":           "2",
						"decryption_list.#":         "0",
						"decompress_max_file_count": acctest.REMOVEKEY,
						"last_modified_start_time":  acctest.REMOVEKEY,
						"decompress_max_layer":      acctest.REMOVEKEY,
						"all_key_prefix":            acctest.REMOVEKEY,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudThreatDetectionOssScanConfigMap5010)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return threatdetection.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeThreatDetectionOssScanConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tFacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudThreatDetectionOssScanConfigBasicDependence5010)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name_list": []string{
						"${alicloud_oss_bucket.default8j4t1R.id}", "${alicloud_oss_bucket.default9HMqfT.id}", "${alicloud_oss_bucket.defaultxBXqFQ.id}"},
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
						"${alicloud_oss_bucket.defaulthZvCmR.id}"},
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
						"${alicloud_oss_bucket.defaulthZvCmR.id}"},
					"key_prefix_list": []string{
						"/123"},
					"start_time": "00:00:08",
					"key_suffix_list": []string{
						".jsp"},
					"scan_day_list": []string{
						"2", "5"},
					"decryption_list":           []string{},
					"decompress_max_file_count": acctest.REMOVEKEY,
					"last_modified_start_time":  acctest.REMOVEKEY,
					"decompress_max_layer":      acctest.REMOVEKEY,
					"all_key_prefix":            acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name_list.#":        "1",
						"key_prefix_list.#":         "1",
						"start_time":                "00:00:08",
						"key_suffix_list.#":         "1",
						"scan_day_list.#":           "2",
						"decryption_list.#":         "0",
						"decompress_max_file_count": acctest.REMOVEKEY,
						"last_modified_start_time":  acctest.REMOVEKEY,
						"decompress_max_layer":      acctest.REMOVEKEY,
						"all_key_prefix":            acctest.REMOVEKEY,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudThreatDetectionOssScanConfigMap7682)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return threatdetection.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeThreatDetectionOssScanConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tFacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudThreatDetectionOssScanConfigBasicDependence7682)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name_list": []string{
						"${alicloud_oss_bucket.default8j4t1R.id}", "${alicloud_oss_bucket.default9HMqfT.id}", "${alicloud_oss_bucket.defaultxBXqFQ.id}"},
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
						"${alicloud_oss_bucket.defaulthZvCmR.id}"},
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
						"${alicloud_oss_bucket.defaulthZvCmR.id}"},
					"key_prefix_list": []string{
						"/123"},
					"start_time": "00:00:08",
					"key_suffix_list": []string{
						".jsp"},
					"scan_day_list": []string{
						"2", "5"},
					"decompress_max_file_count": acctest.REMOVEKEY,
					"last_modified_start_time":  acctest.REMOVEKEY,
					"decompress_max_layer":      acctest.REMOVEKEY,
					"all_key_prefix":            acctest.REMOVEKEY,
					"decryption_list":           acctest.CLEARLIST,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name_list.#":        "1",
						"key_prefix_list.#":         "1",
						"start_time":                "00:00:08",
						"key_suffix_list.#":         "1",
						"scan_day_list.#":           "2",
						"decompress_max_file_count": acctest.REMOVEKEY,
						"last_modified_start_time":  acctest.REMOVEKEY,
						"decompress_max_layer":      acctest.REMOVEKEY,
						"all_key_prefix":            acctest.REMOVEKEY,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudThreatDetectionOssScanConfigMap9153)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return threatdetection.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeThreatDetectionOssScanConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tFacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudThreatDetectionOssScanConfigBasicDependence9153)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name_list": []string{
						"${alicloud_oss_bucket.default8j4t1R.id}", "${alicloud_oss_bucket.default9HMqfT.id}", "${alicloud_oss_bucket.defaultxBXqFQ.id}"},
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
						"${alicloud_oss_bucket.defaulthZvCmR.id}"},
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
						"${alicloud_oss_bucket.defaulthZvCmR.id}"},
					"key_prefix_list": []string{
						"/123"},
					"start_time": "00:00:08",
					"key_suffix_list": []string{
						".jsp"},
					"scan_day_list": []string{
						"2", "5"},
					"decryption_list":           []string{},
					"decompress_max_file_count": acctest.REMOVEKEY,
					"last_modified_start_time":  acctest.REMOVEKEY,
					"decompress_max_layer":      acctest.REMOVEKEY,
					"all_key_prefix":            acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name_list.#":        "1",
						"key_prefix_list.#":         "1",
						"start_time":                "00:00:08",
						"key_suffix_list.#":         "1",
						"scan_day_list.#":           "2",
						"decryption_list.#":         "0",
						"decompress_max_file_count": acctest.REMOVEKEY,
						"last_modified_start_time":  acctest.REMOVEKEY,
						"decompress_max_layer":      acctest.REMOVEKEY,
						"all_key_prefix":            acctest.REMOVEKEY,
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
