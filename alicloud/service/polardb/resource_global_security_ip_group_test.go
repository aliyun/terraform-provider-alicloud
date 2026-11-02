// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package polardb_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	polardb "github.com/aliyun/terraform-provider-alicloud/alicloud/service/polardb"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Polardb GlobalSecurityIpGroup. >>> Resource test cases, automatically generated.
// Case IP白名单模板用例 11862
func TestAccAliCloudPolardbGlobalSecurityIpGroup_basic11862(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_global_security_ip_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPolardbGlobalSecurityIpGroupMap11862)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return polardb.NewPolardbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolardbGlobalSecurityIpGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccpolardb%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudPolardbGlobalSecurityIpGroupBasicDependence11862)
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
					"global_ip_list":       "192.168.0.1",
					"global_ip_group_name": "test_template",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_ip_list":       "192.168.0.1",
						"global_ip_group_name": "test_template",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_ip_list":       "192.168.2.3",
					"global_ip_group_name": "test_name_new",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_ip_list":       "192.168.2.3",
						"global_ip_group_name": "test_name_new",
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

var AlicloudPolardbGlobalSecurityIpGroupMap11862 = map[string]string{
	"region_id": CHECKSET,
}

func AlicloudPolardbGlobalSecurityIpGroupBasicDependence11862(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Polardb GlobalSecurityIpGroup. <<< Resource test cases, automatically generated.
