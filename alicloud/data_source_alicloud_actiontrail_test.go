package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudActiontrailDataSource_name(t *testing.T) {

	num := acctest.RandIntRange(10000, 99999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudActiontrailDataSource_name(num),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_actiontrails.trails"),
					resource.TestCheckResourceAttr("data.alicloud_actiontrails.trails", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_actiontrails.trails", "names.0", fmt.Sprintf("tf-testacc-actiontrail-%v", num)),
					resource.TestCheckResourceAttr("data.alicloud_actiontrails.trails", "actiontrails.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_actiontrails.trails", "actiontrails.0.name", fmt.Sprintf("tf-testacc-actiontrail-%v", num)),
					resource.TestCheckResourceAttr("data.alicloud_actiontrails.trails", "actiontrails.0.event_rw", "Write"),
					resource.TestCheckResourceAttr("data.alicloud_actiontrails.trails", "actiontrails.0.oss_bucket_name", fmt.Sprintf("tf-testacc-actiontrail-%v", num)),
					resource.TestCheckResourceAttr("data.alicloud_actiontrails.trails", "actiontrails.0.role_name", fmt.Sprintf("tf-testacc-actiontrail-%v", num)),
					resource.TestCheckResourceAttr("data.alicloud_actiontrails.trails", "actiontrails.0.oss_key_prefix", "at-product-account-audit-B"),
					resource.TestCheckResourceAttr("data.alicloud_actiontrails.trails", "actiontrails.0.sls_project_arn", ""),
					resource.TestCheckResourceAttr("data.alicloud_actiontrails.trails", "actiontrails.0.sls_write_role_arn", ""),
				),
			},
			{
				Config: testAccCheckAlicloudActiontrailDataSource_name_empty(num),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_actiontrails.trails"),
					resource.TestCheckResourceAttr("data.alicloud_actiontrails.trails", "actiontrails.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_actiontrails.trails", "names.#", "0"),
				),
			},
		},
	})
}

func testAccCheckAlicloudActiontrailDataSource_name(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
    default = "tf-testacc-actiontrail-%v"
}

resource "alicloud_ram_role" "role" {
	  name = "${var.name}"
	  services = ["actiontrail.aliyuncs.com", "oss.aliyuncs.com"]
	  description = "this is a test"
	  force = "true"
}

resource "alicloud_oss_bucket" "bucket" {
    bucket  = "${var.name}"
}

resource "alicloud_ram_policy" "policy" {
	  name = "${var.name}"
	  statement = [
	    {
	      effect = "Allow"
	      action = ["*"]
	      resource = [
		"acs:oss:*:*:${alicloud_oss_bucket.bucket.id}",
		"acs:oss:*:*:${alicloud_oss_bucket.bucket.id}"]
	    }]
	  description = "this is a policy test"
	  force = true
	}

	resource "alicloud_ram_role_policy_attachment" "attach" {
    	  policy_name = "${alicloud_ram_policy.policy.name}"
    	  role_name = "${alicloud_ram_role.role.name}"
    	  policy_type = "${alicloud_ram_policy.policy.type}"
	}
	
resource "alicloud_actiontrail" "foo" {
	name = "${var.name}"
	event_rw = "Write"
	oss_bucket_name = "${alicloud_oss_bucket.bucket.id}"
	role_name = "${alicloud_ram_role_policy_attachment.attach.role_name}"
	oss_key_prefix = "at-product-account-audit-B"
}

data "alicloud_actiontrails" "trails"{
	  name_regex = "${alicloud_actiontrail.foo.name}"
	}

`, randInt)
}

func testAccCheckAlicloudActiontrailDataSource_name_empty(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
    default = "tf-testacc-actiontrail-%v"
}

resource "alicloud_ram_role" "role" {
	  name = "${var.name}"
	  services = ["actiontrail.aliyuncs.com", "oss.aliyuncs.com"]
	  description = "this is a test"
	  force = "true"
}

resource "alicloud_oss_bucket" "bucket" {
    bucket  = "${var.name}"
}

resource "alicloud_ram_policy" "policy" {
	  name = "${var.name}"
	  statement = [
	    {
	      effect = "Allow"
	      action = ["*"]
	      resource = [
		"acs:oss:*:*:${alicloud_oss_bucket.bucket.id}",
		"acs:oss:*:*:${alicloud_oss_bucket.bucket.id}"]
	    }]
	  description = "this is a policy test"
	  force = true
	}

	resource "alicloud_ram_role_policy_attachment" "attach" {
    	  policy_name = "${alicloud_ram_policy.policy.name}"
    	  role_name = "${alicloud_ram_role.role.name}"
    	  policy_type = "${alicloud_ram_policy.policy.type}"
	}
	
resource "alicloud_actiontrail" "foo" {
	name = "${var.name}"
	event_rw = "Write-test"
	oss_bucket_name = "${alicloud_oss_bucket.bucket.id}"
	role_name = "${alicloud_ram_role_policy_attachment.attach.role_name}"
	oss_key_prefix = "at-product-account-audit-B"
}

data "alicloud_actiontrails" "trails"{
	  name_regex = "${alicloud_actiontrail.foo.name}-fake"
	}

`, randInt)
}
