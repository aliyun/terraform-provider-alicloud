package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ThreatDetection FileProtectClientRule. >>> Resource test cases.
func TestAccAliCloudThreatDetectionFileProtectClientRule_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_file_protect_client_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionFileProtectClientRuleMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionFileProtectClientRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionFileProtectClientRuleBasicDependence)
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
					"rule_name":     name,
					"rule_action":   "pass",
					"status":        "0",
					"platform":      "linux",
					"alert_level":   "1",
					"file_paths":    []string{"/opt/a", "/tmp/d"},
					"file_ops":      []string{"READ"},
					"proc_paths":    []string{"/usr/bin/java"},
					"file_types":    []string{"sh"},
					"exclude_users": []string{"root"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name":       name,
						"rule_action":     "pass",
						"status":          "0",
						"platform":        "linux",
						"alert_level":     "1",
						"file_paths.#":    "2",
						"file_ops.#":      "1",
						"proc_paths.#":    "1",
						"file_types.#":    "1",
						"exclude_users.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name":     name + "_update",
					"rule_action":   "monitor",
					"status":        "1",
					"platform":      "linux",
					"alert_level":   "3",
					"file_paths":    []string{"/opt/b"},
					"file_ops":      []string{"READ", "WRITE", "DELETE"},
					"proc_paths":    []string{"/usr/bin/java", "/usr/bin/python"},
					"file_types":    []string{"sh", "py"},
					"exclude_users": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name":       name + "_update",
						"rule_action":     "monitor",
						"status":          "1",
						"alert_level":     "3",
						"file_paths.#":    "1",
						"file_ops.#":      "3",
						"proc_paths.#":    "2",
						"file_types.#":    "2",
						"exclude_users.#": "0",
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

func TestAccAliCloudThreatDetectionFileProtectClientRule_basic_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_file_protect_client_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionFileProtectClientRuleMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionFileProtectClientRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionFileProtectClientRuleBasicDependence)
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
					"rule_name":   name,
					"rule_action": "block",
					"status":      "1",
					"platform":    "windows",
					"file_paths":  []string{"C:\\\\inetpub"},
					"file_ops":    []string{"WRITE", "CHMOD", "RENAME"},
					"proc_paths":  []string{"C:\\\\Windows\\\\System32"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name":    name,
						"rule_action":  "block",
						"status":       "1",
						"platform":     "windows",
						"file_paths.#": "1",
						"file_ops.#":   "3",
						"proc_paths.#": "1",
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

var AliCloudThreatDetectionFileProtectClientRuleMap = map[string]string{}

func AliCloudThreatDetectionFileProtectClientRuleBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}

// Test ThreatDetection FileProtectClientRule. <<< Resource test cases.
