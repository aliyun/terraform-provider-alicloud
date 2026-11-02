// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package eflo_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eflo"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Eflo VpdGrantRule. >>> Resource test cases, automatically generated.
// Case vpd_grant_rule_tf 11414
func TestAccAliCloudEfloVpdGrantRule_basic11414(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_vpd_grant_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEfloVpdGrantRuleMap11414)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eflo.NewEfloServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEfloVpdGrantRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEfloVpdGrantRuleBasicDependence11414)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"grant_tenant_id": "${data.alicloud_account.default.id}",
					"er_id":           "${alicloud_eflo_er.default.id}",
					"instance_id":     "${alicloud_eflo_vpd.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grant_tenant_id": acctest.CHECKSET,
						"er_id":           acctest.CHECKSET,
						"instance_id":     acctest.CHECKSET,
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
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
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
