// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Cr StorageDomainRoutingRule. >>> Resource test cases, automatically generated.
// Case storage-domain-routing-rule-test 11834
func TestAccAliCloudCrStorageDomainRoutingRule_basic11834(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_storage_domain_routing_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudCrStorageDomainRoutingRuleMap11834)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrStorageDomainRoutingRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCrStorageDomainRoutingRuleBasicDependence11834)
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
					"routes": []map[string]interface{}{
						{
							"instance_domain": "${alicloud_cr_ee_instance.default.instance_name}-registry-vpc.cn-hangzhou.cr.aliyuncs.com",
							"storage_domain":  "https://${alicloud_cr_ee_instance.default.id}-registry.oss-cn-hangzhou-internal.aliyuncs.com",
							"endpoint_type":   "Internet",
						},
					},
					"instance_id": "${alicloud_cr_ee_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"routes.#":    "1",
						"instance_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"routes": []map[string]interface{}{
						{
							"instance_domain": "${alicloud_cr_ee_instance.default.instance_name}-registry-vpc.cn-hangzhou.cr.aliyuncs.com",
							"storage_domain":  "https://${alicloud_cr_ee_instance.default.id}-registry.oss-cn-hangzhou-internal.aliyuncs.com",
							"endpoint_type":   "VPC",
						},
						{
							"instance_domain": "${alicloud_cr_ee_instance.default.instance_name}-registry.cn-hangzhou.cr.aliyuncs.com",
							"storage_domain":  "https://${alicloud_cr_ee_instance.default.id}-registry.oss-cn-hangzhou.aliyuncs.com",
							"endpoint_type":   "Internet",
						},
						{
							"instance_domain": "${alicloud_cr_ee_instance.default.instance_name}-registry-intranet.cn-hangzhou.cr.aliyuncs.com",
							"storage_domain":  "https://${alicloud_cr_ee_instance.default.id}-registry.oss-cn-hangzhou.aliyuncs.com",
							"endpoint_type":   "VPC",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"routes.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"routes": []map[string]interface{}{
						{
							"instance_domain": "${alicloud_cr_ee_instance.default.instance_name}-registry-vpc.cn-hangzhou.cr.aliyuncs.com",
							"storage_domain":  "https://${alicloud_cr_ee_instance.default.id}-registry.oss-cn-hangzhou-internal.aliyuncs.com",
							"endpoint_type":   "VPC",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"routes.#": "1",
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

func TestAccAliCloudCrStorageDomainRoutingRule_basic11834_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_storage_domain_routing_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudCrStorageDomainRoutingRuleMap11834)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrStorageDomainRoutingRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCrStorageDomainRoutingRuleBasicDependence11834)
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
					"routes": []map[string]interface{}{
						{
							"instance_domain": "${alicloud_cr_ee_instance.default.instance_name}-registry-vpc.cn-hangzhou.cr.aliyuncs.com",
							"storage_domain":  "https://${alicloud_cr_ee_instance.default.id}-registry.oss-cn-hangzhou-internal.aliyuncs.com",
							"endpoint_type":   "VPC",
						},
						{
							"instance_domain": "${alicloud_cr_ee_instance.default.instance_name}-registry.cn-hangzhou.cr.aliyuncs.com",
							"storage_domain":  "https://${alicloud_cr_ee_instance.default.id}-registry.oss-cn-hangzhou.aliyuncs.com",
							"endpoint_type":   "Internet",
						},
						{
							"instance_domain": "${alicloud_cr_ee_instance.default.instance_name}-registry-intranet.cn-hangzhou.cr.aliyuncs.com",
							"storage_domain":  "https://${alicloud_cr_ee_instance.default.id}-registry.oss-cn-hangzhou.aliyuncs.com",
							"endpoint_type":   "VPC",
						},
					},
					"instance_id": "${alicloud_cr_ee_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"routes.#":    "3",
						"instance_id": CHECKSET,
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

var AliCloudCrStorageDomainRoutingRuleMap11834 = map[string]string{
	"rule_id":     CHECKSET,
	"create_time": CHECKSET,
}

func AliCloudCrStorageDomainRoutingRuleBasicDependence11834(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_cr_ee_instance" "default" {
  payment_type   = "Subscription"
  period         = 1
  renew_period   = 1
  renewal_status = "AutoRenewal"
  instance_type  = "Advanced"
  instance_name  = var.name
}
`, name)
}

// Test Cr StorageDomainRoutingRule. <<< Resource test cases, automatically generated.
