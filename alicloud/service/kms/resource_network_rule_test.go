package kms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/kms"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// collection-order: bypass source_private_ip -- pre-existing stock TypeList; schema and member-order coverage unchanged by this pure directory migration (route provider_handwritten)
// Test Kms NetworkRule. >>> Resource test cases, automatically generated.
// Case 3872
func TestAccAliCloudKmsNetworkRule_basic3872(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_network_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudKmsNetworkRuleMap3872)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return kms.NewKmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeKmsNetworkRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmsnetworkrule%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudKmsNetworkRuleBasicDependence3872)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
	"network_rule_name": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudKmsNetworkRuleMap3872)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return kms.NewKmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeKmsNetworkRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmsnetworkrule%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudKmsNetworkRuleBasicDependence3872)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
