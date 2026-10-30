package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudVPC_internationalAccount(t *testing.T) {
	var v map[string]interface{}
	id := "alicloud_vpc.default"
	ra := resourceAttrInit(id, AlicloudVpcMap)
	rc := resourceCheckInitWithDescribeMethod(id, &v, func() interface{} { return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)} }, "DescribeVpc")
	rac := resourceAttrCheckInit(rc, ra)
	name := fmt.Sprintf("tf-testAcc-intl-vpc-%d", acctest.RandIntRange(1000000, 9999999))
	config := func(description string) string {
		return fmt.Sprintf(`
provider "alicloud" {
 region = "us-west-1"
}
resource "alicloud_vpc" "default" {
 vpc_name = %q
 cidr_block = "192.168.0.0/16"
 description = %q
}`, name, description)
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckWithAccountSiteType(t, IntlSite); testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{Config: config("international acceptance test"), Check: resource.TestCheckResourceAttr(id, "vpc_name", name)},
			{ResourceName: id, ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"dry_run", "enable_ipv6"}},
			{Config: config("international acceptance test updated"), Check: resource.TestCheckResourceAttr(id, "description", "international acceptance test updated")},
		},
	})
}
