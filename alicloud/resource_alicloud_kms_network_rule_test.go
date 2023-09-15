package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Kms NetworkRule. >>> Resource test cases, automatically generated.
// Case 3872
func TestAccAliCloudKmsNetworkRule_basic3872(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_network_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsNetworkRuleMap3872)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsNetworkRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmsnetworkrule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsNetworkRuleBasicDependence3872)
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
					"source_private_ip": []string{
						"10.10.10.10/24", "192.168.17.13", "100.177.24.254"},
					"network_rule_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_private_ip.#": "3",
						"network_rule_name":   name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test-description",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-description-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test-description-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_private_ip": []string{
						"71.71.71.71"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_private_ip.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-description",
					"source_private_ip": []string{
						"10.10.10.10/24", "192.168.17.13", "100.177.24.254"},
					"network_rule_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":         "test-description",
						"source_private_ip.#": "3",
						"network_rule_name":   name + "_update",
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

var AlicloudKmsNetworkRuleMap3872 = map[string]string{
	"network_rule_name": CHECKSET,
}

func AlicloudKmsNetworkRuleBasicDependence3872(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 3872  twin
func TestAccAliCloudKmsNetworkRule_basic3872_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_network_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsNetworkRuleMap3872)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsNetworkRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmsnetworkrule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsNetworkRuleBasicDependence3872)
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
					"description": "test-description",
					"source_private_ip": []string{
						"10.10.10.10/24", "192.168.17.13", "100.177.24.254"},
					"network_rule_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":         "test-description",
						"source_private_ip.#": "3",
						"network_rule_name":   name,
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

// Test Kms NetworkRule. <<< Resource test cases, automatically generated.
