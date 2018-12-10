package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudKVStoreParameter_basic(t *testing.T) {
	var parameter r_kvstore.DescribeParametersResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_parameter.compat",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreParameterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccKVStoreParameter_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreParameterExist(
						"alicloud_kvstore_parameter.compat", &parameter),
					resource.TestCheckResourceAttr("alicloud_kvstore_parameter.compat", "name", "cluster_compat_enable"),
					resource.TestCheckResourceAttr("alicloud_kvstore_parameter.compat", "value", "1"),
				),
			},
		},
	})

}

func testAccCheckKVStoreParameterExist(n string, d *r_kvstore.DescribeParametersResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No KVStore Instance parameter ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		kvstoreService := KvstoreService{client}
		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		response, err := kvstoreService.DescribeRKVInstanceParameter(parts[0])
		if err != nil {
			return fmt.Errorf("Error Describe KVStore Instance parameter: %#v", err)
		}

		*d = *response
		return nil
	}
}

func testAccCheckKVStoreParameterDestroy(s *terraform.State) error {

	return nil
}

const testAccKVStoreParameter_basic = `
data "alicloud_zones" "default" {
	available_resource_creation = "KVStore"
}
variable "name" {
	default = "tf-testacckvstoreparameter_basic"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
 	vpc_id = "${alicloud_vpc.foo.id}"
 	cidr_block = "172.16.0.0/21"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
 	name = "${var.name}"
}

resource "alicloud_kvstore_instance" "foo" {
	instance_class = "redis.logic.sharding.2g.2db.0rodb.4proxy.default"
	instance_name  = "${var.name}"
	password       = "Test12345"
	vswitch_id     = "${alicloud_vswitch.foo.id}"
}

resource "alicloud_kvstore_parameter" "compat" {
	instance_id = "${alicloud_kvstore_instance.foo.id}"
	name = "cluster_compat_enable"
	value = "1"
}
`
