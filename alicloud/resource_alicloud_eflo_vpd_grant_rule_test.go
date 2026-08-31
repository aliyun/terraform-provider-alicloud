// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Eflo VpdGrantRule. >>> Resource test cases, automatically generated.
// Case vpd_grant_rule_tf 11414
func TestAccAliCloudEfloVpdGrantRule_basic11414(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_vpd_grant_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudEfloVpdGrantRuleMap11414)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloVpdGrantRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEfloVpdGrantRuleBasicDependence11414)
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
					"grant_tenant_id": "${data.alicloud_account.default.id}",
					"er_id":           "${alicloud_eflo_er.default.id}",
					"instance_id":     "${alicloud_eflo_vpd.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grant_tenant_id": CHECKSET,
						"er_id":           CHECKSET,
						"instance_id":     CHECKSET,
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

var AliCloudEfloVpdGrantRuleMap11414 = map[string]string{
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AliCloudEfloVpdGrantRuleBasicDependence11414(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "default" {
}

resource "alicloud_eflo_er" "default" {
  er_name        = var.name
  master_zone_id = "cn-hangzhou-a"
}

resource "alicloud_eflo_vpd" "default" {
  cidr     = "10.0.0.0/8"
  vpd_name = var.name
}
`, name)
}

// Test Eflo VpdGrantRule. <<< Resource test cases, automatically generated.
