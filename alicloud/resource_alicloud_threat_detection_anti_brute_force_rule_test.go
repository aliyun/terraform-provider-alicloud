package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ThreatDetection AntiBruteForceRule. >>> Resource test cases, automatically generated.
// Case 对接Terraform 4056
func TestAccAliCloudThreatDetectionAntiBruteForceRule_basic4056(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_anti_brute_force_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionAntiBruteForceRuleMap4056)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionAntiBruteForceRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionAntiBruteForceRuleBasicDependence4056)
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
					"default_rule":               "false",
					"anti_brute_force_rule_name": name,
					"forbidden_time":             "360",
					"uuid_list": []string{
						"032b618f-b220-4a0d-bd37-fbdc6ef58b6a"},
					"fail_count": "80",
					"span":       "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_rule":               "false",
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
					"forbidden_time":             "300",
					"fail_count":                 "70",
					"span":                       "9",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"anti_brute_force_rule_name": name + "_update",
						"forbidden_time":             "300",
						"fail_count":                 "70",
						"span":                       "9",
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

var AlicloudThreatDetectionAntiBruteForceRuleMap4056 = map[string]string{}

func AlicloudThreatDetectionAntiBruteForceRuleBasicDependence4056(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 对接Terraform_副本1725242231076 7709
func TestAccAliCloudThreatDetectionAntiBruteForceRule_basic7709(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_anti_brute_force_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionAntiBruteForceRuleMap7709)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionAntiBruteForceRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionAntiBruteForceRuleBasicDependence7709)
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
					"default_rule":               "false",
					"anti_brute_force_rule_name": name,
					"forbidden_time":             "360",
					"uuid_list": []string{
						"032b618f-b220-4a0d-bd37-fbdc6ef58b6a"},
					"fail_count": "80",
					"span":       "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_rule":               "false",
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
					"forbidden_time":             "300",
					"fail_count":                 "70",
					"span":                       "9",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"anti_brute_force_rule_name": name + "_update",
						"forbidden_time":             "300",
						"fail_count":                 "70",
						"span":                       "9",
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

var AlicloudThreatDetectionAntiBruteForceRuleMap7709 = map[string]string{}

func AlicloudThreatDetectionAntiBruteForceRuleBasicDependence7709(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 对接Terraform_副本1725500044486 7780
func TestAccAliCloudThreatDetectionAntiBruteForceRule_basic7780(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_anti_brute_force_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionAntiBruteForceRuleMap7780)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionAntiBruteForceRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionAntiBruteForceRuleBasicDependence7780)
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
					"default_rule":               "false",
					"anti_brute_force_rule_name": name,
					"forbidden_time":             "360",
					"uuid_list": []string{
						"inet-eae024ff-c25b-4168-b839-f6ecc67e9db4", "1d343293-52db-4684-a5e1-5feb1e0670fa"},
					"fail_count": "80",
					"span":       "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_rule":               "false",
						"anti_brute_force_rule_name": name,
						"forbidden_time":             "360",
						"uuid_list.#":                "2",
						"fail_count":                 "80",
						"span":                       "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_rule":               "true",
					"anti_brute_force_rule_name": name + "_update",
					"forbidden_time":             "300",
					"uuid_list": []string{
						"1d343293-52db-4684-a5e1-5feb1e0670fa"},
					"fail_count": "70",
					"span":       "9",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_rule":               "true",
						"anti_brute_force_rule_name": name + "_update",
						"forbidden_time":             "300",
						"uuid_list.#":                "1",
						"fail_count":                 "70",
						"span":                       "9",
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

var AlicloudThreatDetectionAntiBruteForceRuleMap7780 = map[string]string{}

func AlicloudThreatDetectionAntiBruteForceRuleBasicDependence7780(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_threat_detection_assets" "default" {
    machine_types = "ecs"
}


`, name)
}

// Test ThreatDetection AntiBruteForceRule. <<< Resource test cases, automatically generated.
