package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ThreatDetection MaliciousFileWhitelistConfig. >>> Resource test cases, automatically generated.
// Case 4979
func TestAccAliCloudThreatDetectionMaliciousFileWhitelistConfig_basic4979(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_malicious_file_whitelist_config.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionMaliciousFileWhitelistConfigMap4979)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionMaliciousFileWhitelistConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sthreatdetectionmaliciousfilewhitelistconfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionMaliciousFileWhitelistConfigBasicDependence4979)
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
					"operator":     "strEqual",
					"field":        "fileMd5",
					"target_value": "ALL",
					"target_type":  "ALL",
					"event_name":   "ALL",
					"source":       "agentless",
					"field_value":  "714",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"operator":     "strEqual",
						"field":        "fileMd5",
						"target_value": "ALL",
						"target_type":  "ALL",
						"event_name":   "ALL",
						"source":       "agentless",
						"field_value":  "714",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"field": "fileMd5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"field": "fileMd5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"target_value": "ALL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_value": "ALL",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"target_type": "ALL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_type": "ALL",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"event_name": "ALL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_name": "ALL",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"field_value": "186",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"field_value": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"operator": "strEquals",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"operator": "strEquals",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"field": "fileMd6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"field": "fileMd6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"target_value": "123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_value": "123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"target_type": "SELECTION_KEY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_type": "SELECTION_KEY",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"event_name": "123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_name": "123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"field_value": "sadfas",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"field_value": "sadfas",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"operator": "strEqual",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"operator": "strEqual",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"field": "fileMd5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"field": "fileMd5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"target_value": "ALL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_value": "ALL",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"target_type": "ALL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_type": "ALL",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"event_name": "ALL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_name": "ALL",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"field_value": "551",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"field_value": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"operator": "strEquals",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"operator": "strEquals",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"field": "fileMd6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"field": "fileMd6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"target_value": "123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_value": "123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"target_type": "SELECTION_KEY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_type": "SELECTION_KEY",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"event_name": "123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_name": "123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"field_value": "sadfas",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"field_value": "sadfas",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"operator":     "strEqual",
					"field":        "fileMd5",
					"target_value": "ALL",
					"target_type":  "ALL",
					"event_name":   "ALL",
					"source":       "agentless",
					"field_value":  "714",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"operator":     "strEqual",
						"field":        "fileMd5",
						"target_value": "ALL",
						"target_type":  "ALL",
						"event_name":   "ALL",
						"source":       "agentless",
						"field_value":  CHECKSET,
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

var AlicloudThreatDetectionMaliciousFileWhitelistConfigMap4979 = map[string]string{}

func AlicloudThreatDetectionMaliciousFileWhitelistConfigBasicDependence4979(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 4979  twin
func TestAccAliCloudThreatDetectionMaliciousFileWhitelistConfig_basic4979_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_malicious_file_whitelist_config.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionMaliciousFileWhitelistConfigMap4979)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionMaliciousFileWhitelistConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sthreatdetectionmaliciousfilewhitelistconfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionMaliciousFileWhitelistConfigBasicDependence4979)
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
					"operator":     "strEquals",
					"field":        "fileMd6",
					"target_value": "123",
					"target_type":  "SELECTION_KEY",
					"event_name":   "123",
					"source":       "agentless",
					"field_value":  "sadfas",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"operator":     "strEquals",
						"field":        "fileMd6",
						"target_value": "123",
						"target_type":  "SELECTION_KEY",
						"event_name":   "123",
						"source":       "agentless",
						"field_value":  "sadfas",
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

// Test ThreatDetection MaliciousFileWhitelistConfig. <<< Resource test cases, automatically generated.
