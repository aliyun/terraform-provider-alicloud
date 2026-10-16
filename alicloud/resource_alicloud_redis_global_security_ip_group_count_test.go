package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// DescribeGlobalSecurityIPGroup is a list query that returns every global
// security IP group, so Read must locate the matching group by
// GlobalSecurityGroupId rather than blindly taking list element [0]. This
// regression test creates several groups with count > 1 and then asserts a
// clean plan, which fails when Read writes the first list entry into every
// instance's state and produces a persistent diff.
func TestAccAliCloudRedisGlobalSecurityIpGroup_count(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_global_security_ip_group.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisGlobalSecurityIpGroupMap12846)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisGlobalSecurityIpGroup")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccredis%d", rand)
	testConfig := fmt.Sprintf(`
resource "alicloud_redis_global_security_ip_group" "default" {
  count                = 3
  global_ip_group_name = "%s${count.index}"
  global_ip_list       = "192.168.${count.index}.1,10.10.10.1${count.index}"
}
`, name)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_redis_global_security_ip_group.default.0", "global_ip_group_name", name+"0"),
					resource.TestCheckResourceAttr("alicloud_redis_global_security_ip_group.default.0", "global_ip_list", "192.168.0.1,10.10.10.10"),
					resource.TestCheckResourceAttr("alicloud_redis_global_security_ip_group.default.1", "global_ip_group_name", name+"1"),
					resource.TestCheckResourceAttr("alicloud_redis_global_security_ip_group.default.1", "global_ip_list", "192.168.1.1,10.10.10.11"),
					resource.TestCheckResourceAttr("alicloud_redis_global_security_ip_group.default.2", "global_ip_group_name", name+"2"),
					resource.TestCheckResourceAttr("alicloud_redis_global_security_ip_group.default.2", "global_ip_list", "192.168.2.1,10.10.10.12"),
				),
			},
			{
				Config:             testConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
