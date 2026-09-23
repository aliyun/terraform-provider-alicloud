package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ThreatDetection AntiBruteForceRule. >>> Resource test cases, automatically generated.
// Case 对接Terraform 7780
func TestAccAliCloudThreatDetectionAntiBruteForceRule_basic7780(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_anti_brute_force_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionAntiBruteForceRuleMap7780)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionAntiBruteForceRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionAntiBruteForceRuleBasicDependence7780)
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
					"anti_brute_force_rule_name": name,
					"forbidden_time":             "360",
					"uuid_list": []string{
						"032b618f-b220-4a0d-bd37-fbdc6ef58b6a"},
					"fail_count": "80",
					"span":       "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"anti_brute_force_rule_name": name,
						"forbidden_time":             "360",
						"uuid_list.#":                "1",
						"fail_count":                 "80",
						"span":                       "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"anti_brute_force_rule_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"anti_brute_force_rule_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"forbidden_time": "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"forbidden_time": "300",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"fail_count": "70",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fail_count": "70",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"span": "9",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"span": "9",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"uuid_list": []string{"inet-eae024ff-c25b-4168-b839-f6ecc67e9db4", "1d343293-52db-4684-a5e1-5feb1e0670fa"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"uuid_list.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_rule": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_rule": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol_type": []map[string]interface{}{
						{
							"rdp":        "off",
							"ssh":        "off",
							"sql_server": "on",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type.#": "1",
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

func TestAccAliCloudThreatDetectionAntiBruteForceRule_basic7780_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_anti_brute_force_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionAntiBruteForceRuleMap7780)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionAntiBruteForceRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionAntiBruteForceRuleBasicDependence7780)
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
					"anti_brute_force_rule_name": name,
					"forbidden_time":             "360",
					"uuid_list": []string{
						"032b618f-b220-4a0d-bd37-fbdc6ef58b6a"},
					"fail_count":   "80",
					"span":         "10",
					"default_rule": "true",
					"protocol_type": []map[string]interface{}{
						{
							"rdp":        "off",
							"ssh":        "off",
							"sql_server": "on",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"anti_brute_force_rule_name": name,
						"forbidden_time":             "360",
						"uuid_list.#":                "1",
						"fail_count":                 "80",
						"span":                       "10",
						"default_rule":               "true",
						"protocol_type.#":            "1",
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

var AliCloudThreatDetectionAntiBruteForceRuleMap7780 = map[string]string{}

func AliCloudThreatDetectionAntiBruteForceRuleBasicDependence7780(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test ThreatDetection AntiBruteForceRule. <<< Resource test cases, automatically generated.
