package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Kms Policy. >>> Resource test cases, automatically generated.
// Case 3883
func TestAccAliCloudKmsPolicy_basic3883(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsPolicyMap3883)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmspolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsPolicyBasicDependence3883)
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
					"permissions": []string{
						"RbacPermission/Template/CryptoServiceKeyUser", "RbacPermission/Template/CryptoServiceSecretUser"},
					"resources": []string{
						"secret/*", "key/*"},
					"policy_name":          name,
					"kms_instance_id":      "shared",
					"access_control_rules": "{\\\"NetworkRules\\\":[\\\"alicloud_kms_network_rule.networkRule1.network_rule_name\\\"]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"permissions.#":        "2",
						"resources.#":          "2",
						"policy_name":          name,
						"kms_instance_id":      "shared",
						"access_control_rules": "{\"NetworkRules\":[\"alicloud_kms_network_rule.networkRule1.network_rule_name\"]}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "terraform测试policy",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "terraform测试policy",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"permissions": []string{
						"RbacPermission/Template/CryptoServiceKeyUser", "RbacPermission/Template/CryptoServiceSecretUser"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"permissions.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resources": []string{
						"secret/*", "key/*"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resources.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_control_rules": "{\\\"NetworkRules\\\":[\\\"alicloud_kms_network_rule.networkRule1.network_rule_name\\\"]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_control_rules": "{\"NetworkRules\":[\"alicloud_kms_network_rule.networkRule1.network_rule_name\"]}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "asdfasdf",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "asdfasdf",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"permissions": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"permissions.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resources": []string{
						"key/*"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resources.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_control_rules": "{\\\"NetworkRules\\\":[\\\"alicloud_kms_network_rule.networkRule2.network_rule_name\\\"]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_control_rules": "{\"NetworkRules\":[\"alicloud_kms_network_rule.networkRule2.network_rule_name\"]}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "sadfadsfasdf",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "sadfadsfasdf",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"permissions": []string{
						"RbacPermission/Template/CryptoServiceSecretUser"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"permissions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resources": []string{
						"key/a*", "key/b*"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resources.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_control_rules": "{\\\"NetworkRules\\\":[\\\"alicloud_kms_network_rule.networkRule3.network_rule_name\\\"]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_control_rules": "{\"NetworkRules\":[\"alicloud_kms_network_rule.networkRule3.network_rule_name\"]}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "terraform测试policy",
					"permissions": []string{
						"RbacPermission/Template/CryptoServiceKeyUser", "RbacPermission/Template/CryptoServiceSecretUser"},
					"resources": []string{
						"secret/*", "key/*"},
					"policy_name":          name + "_update",
					"kms_instance_id":      "shared",
					"access_control_rules": "{\\\"NetworkRules\\\":[\\\"alicloud_kms_network_rule.networkRule1.network_rule_name\\\"]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":          "terraform测试policy",
						"permissions.#":        "2",
						"resources.#":          "2",
						"policy_name":          name + "_update",
						"kms_instance_id":      "shared",
						"access_control_rules": "{\"NetworkRules\":[\"alicloud_kms_network_rule.networkRule1.network_rule_name\"]}",
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

var AlicloudKmsPolicyMap3883 = map[string]string{}

func AlicloudKmsPolicyBasicDependence3883(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_kms_network_rule" "networkRule1" {
  description       = "dummy"
  source_private_ip = ["10.10.10.10"]
  network_rule_name = var.name
}

resource "alicloud_kms_network_rule" "networkRule2" {
  description       = "dummy"
  source_private_ip = ["10.10.10.10"]
  network_rule_name = "${var.name}1"
}

resource "alicloud_kms_network_rule" "networkRule3" {
  description       = "dummy"
  source_private_ip = ["10.10.10.10"]
  network_rule_name = "${var.name}2"
}


`, name)
}

// Case 3883  twin
func TestAccAliCloudKmsPolicy_basic3883_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsPolicyMap3883)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmspolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsPolicyBasicDependence3883)
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
					"description": "sadfadsfasdf",
					"permissions": []string{
						"RbacPermission/Template/CryptoServiceSecretUser", "RbacPermission/Template/CryptoServiceSecretUser"},
					"resources": []string{
						"key/a*", "key/b*", "key/c*"},
					"policy_name":          name,
					"kms_instance_id":      "shared",
					"access_control_rules": "{\\\"NetworkRules\\\":[\\\"alicloud_kms_network_rule.networkRule3.network_rule_name\\\"]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":          "sadfadsfasdf",
						"permissions.#":        "2",
						"resources.#":          "3",
						"policy_name":          name,
						"kms_instance_id":      "shared",
						"access_control_rules": "{\"NetworkRules\":[\"alicloud_kms_network_rule.networkRule3.network_rule_name\"]}",
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

// Test Kms Policy. <<< Resource test cases, automatically generated.
