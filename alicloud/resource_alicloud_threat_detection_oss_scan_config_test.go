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
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionOssScanConfigMap5010)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionOssScanConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionOssScanConfigBasicDependence5010)
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
					"end_time":   "00:00:01",
					"start_time": "00:00:00",
					"enable":     "0",
					"key_suffix_list": []string{
						".html", ".php", ".k"},
					"scan_day_list": []string{
						"1", "2", "4", "3"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name_list.#": "3",
						"end_time":           "00:00:01",
						"start_time":         "00:00:00",
						"enable":             "0",
						"key_suffix_list.#":  "3",
						"scan_day_list.#":    "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_prefix_list": []string{
						"/root", "/usr", "/123"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_prefix_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"all_key_prefix": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"all_key_prefix": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name_list": []string{"gcx-test-oss-71"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"decompress_max_layer":      "10",
					"decompress_max_file_count": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"decompress_max_layer":      "10",
						"decompress_max_file_count": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"decompress_max_layer": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"decompress_max_layer": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"decompress_max_file_count": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"decompress_max_file_count": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"decryption_list": []string{
						"OSS", "AES", "TEST", "TEST2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"decryption_list.#": "4",
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
					"last_modified_start_time": "1721713532000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"last_modified_start_time": "1721713532000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"oss_scan_config_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oss_scan_config_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scan_day_list": []string{
						"1", "2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scan_day_list.#": "2",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudThreatDetectionOssScanConfigMap5010 = map[string]string{
	"all_key_prefix": CHECKSET,
}

func AliCloudThreatDetectionOssScanConfigBasicDependence5010(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}

// Case 5010  twin

func TestAccAliCloudThreatDetectionOssScanConfig_basic5010_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_oss_scan_config.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionOssScanConfigMap5010)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionOssScanConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionOssScanConfigBasicDependence5010)
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
					"end_time":                  "00:00:01",
					"start_time":                "00:00:00",
					"enable":                    "0",
					"decompress_max_layer":      "10",
					"decompress_max_file_count": "2",
					"last_modified_start_time":  "1721713532000",
					"all_key_prefix":            "false",
					"oss_scan_config_name":      name,
					"key_suffix_list": []string{
						".html", ".php", ".k"},
					"key_prefix_list": []string{
						"/root", "/usr", "/123"},
					"scan_day_list": []string{
						"1", "2", "4", "3"},
					"decryption_list": []string{
						"OSS", "AES", "TEST", "TEST2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name_list.#":        "3",
						"end_time":                  "00:00:01",
						"start_time":                "00:00:00",
						"enable":                    "0",
						"decompress_max_layer":      "10",
						"decompress_max_file_count": "2",
						"last_modified_start_time":  "1721713532000",
						"all_key_prefix":            "false",
						"oss_scan_config_name":      name,
						"key_suffix_list.#":         "3",
						"key_prefix_list.#":         "3",
						"scan_day_list.#":           "4",
						"decryption_list.#":         "4",
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
