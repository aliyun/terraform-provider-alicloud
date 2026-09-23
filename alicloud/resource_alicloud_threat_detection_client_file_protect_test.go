package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ThreatDetection ClientFileProtect. >>> Resource test cases, automatically generated.
// Case 4514
func TestAccAliCloudThreatDetectionClientFileProtect_basic4514(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_client_file_protect.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionClientFileProtectMap4514)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionClientFileProtect")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sthreatdetectionclientfileprotect%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionClientFileProtectBasicDependence4514)
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
					"alert_level": "0",
					"file_paths": []string{
						"/usr/local"},
					"file_ops": []string{
						"CREATE"},
					"rule_action": "pass",
					"proc_paths": []string{
						"/usr/local"},
					"rule_name": "ruleTest2_843",
					"status":    "0",
					"switch_id": "FILE_PROTECT_RULE_SWITCH_TYPE_1693474122929",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_level":  "0",
						"file_paths.#": "1",
						"file_ops.#":   "1",
						"rule_action":  "pass",
						"proc_paths.#": "1",
						"rule_name":    CHECKSET,
						"status":       "0",
						"switch_id":    "FILE_PROTECT_RULE_SWITCH_TYPE_1693474122929",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proc_paths": []string{
						"/usr/local"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proc_paths.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_level": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_level": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"file_paths": []string{
						"/usr/local"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_paths.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"file_ops": []string{
						"CREATE"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_ops.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_action": "pass",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_action": "pass",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proc_paths": []string{
						"/usr/local"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proc_paths.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": "ruleTest2_662",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"file_paths": []string{
						"/tmp"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_paths.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"file_ops": []string{
						"CHMOD"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_ops.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_action": "alert",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_action": "alert",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proc_paths": []string{
						"/tmp"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proc_paths.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_level": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_level": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": "ruleTest1_927",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"file_paths": []string{
						"/tmp21"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_paths.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"file_ops": []string{
						"DELETE"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_ops.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_action": "pass",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_action": "pass",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proc_paths": []string{
						"/tmp12"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proc_paths.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": "ruleTest_206",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "1",
					"file_paths": []string{
						"/usr/local"},
					"file_ops": []string{
						"CREATE"},
					"rule_action": "pass",
					"proc_paths": []string{
						"/usr/local"},
					"alert_level": "1",
					"switch_id":   "FILE_PROTECT_RULE_SWITCH_TYPE_1693474122929",
					"rule_name":   "ruleTest2_377",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":       "1",
						"file_paths.#": "1",
						"file_ops.#":   "1",
						"rule_action":  "pass",
						"proc_paths.#": "1",
						"alert_level":  "1",
						"switch_id":    "FILE_PROTECT_RULE_SWITCH_TYPE_1693474122929",
						"rule_name":    CHECKSET,
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

var AlicloudThreatDetectionClientFileProtectMap4514 = map[string]string{
	"status": CHECKSET,
}

func AlicloudThreatDetectionClientFileProtectBasicDependence4514(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 4514  twin
func TestAccAliCloudThreatDetectionClientFileProtect_basic4514_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_client_file_protect.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionClientFileProtectMap4514)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionClientFileProtect")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sthreatdetectionclientfileprotect%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionClientFileProtectBasicDependence4514)
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
					"status": "1",
					"file_paths": []string{
						"/tmp21", "/tmp1", "/tmp2"},
					"file_ops": []string{
						"DELETE", "CREATE", "UPDATE"},
					"rule_action": "pass",
					"proc_paths": []string{
						"/tmp12", "/tmp2", "/tmp3"},
					"alert_level": "1",
					"switch_id":   "FILE_PROTECT_RULE_SWITCH_TYPE_1693474122929",
					"rule_name":   "ruleTest_772",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":       "1",
						"file_paths.#": "3",
						"file_ops.#":   "3",
						"rule_action":  "pass",
						"proc_paths.#": "3",
						"alert_level":  "1",
						"switch_id":    "FILE_PROTECT_RULE_SWITCH_TYPE_1693474122929",
						"rule_name":    CHECKSET,
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

// Test ThreatDetection ClientFileProtect. <<< Resource test cases, automatically generated.
