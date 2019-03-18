package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAliyunActionTrail_update(t *testing.T) {

	num := acctest.RandInt()
	num_role := acctest.RandInt()
	num_bucket := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_actiontrail.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testActionTrailDestroy,
		Steps: []resource.TestStep{
			{
				Config: testActionTrailBasicConfig(num),
				Check: resource.ComposeTestCheckFunc(
					testActionTrailExists("alicloud_actiontrail.foo"),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "name", "tf-testacc-actiontrail"),
					resource.TestCheckNoResourceAttr(
						"alicloud_actiontrail.foo", "event_rw"),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "oss_bucket_name", fmt.Sprintf("tf-testacc-actiontrail-%v", num)),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "role_name", fmt.Sprintf("tf-testacc-actiontrail-%v", num)),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "oss_key_prefix", ""),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "sls_project_arn", ""),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "sls_write_role_arn", ""),
				),
			},
			{
				Config: testActionTrailUpdateConfig_event_rw(num),
				Check: resource.ComposeTestCheckFunc(
					testActionTrailExists("alicloud_actiontrail.foo"),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "name", "tf-testacc-actiontrail"),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "event_rw", "Write-test-123"),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "oss_bucket_name", fmt.Sprintf("tf-testacc-actiontrail-%v", num)),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "role_name", fmt.Sprintf("tf-testacc-actiontrail-%v", num)),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "oss_key_prefix", ""),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "sls_project_arn", ""),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "sls_write_role_arn", ""),
				),
			},
			{
				Config: testActionTrailUpdateConfig_oss_key_prefix(num),
				Check: resource.ComposeTestCheckFunc(
					testActionTrailExists("alicloud_actiontrail.foo"),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "name", "tf-testacc-actiontrail"),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "event_rw", "Write-test-123"),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "oss_bucket_name", fmt.Sprintf("tf-testacc-actiontrail-%v", num)),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "role_name", fmt.Sprintf("tf-testacc-actiontrail-%v", num)),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "oss_key_prefix", "at-product-account-audit-B"),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "sls_project_arn", ""),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "sls_write_role_arn", ""),
				),
			},
			{
				Config: testActionTrailUpdateConfig_role(num, num_role),
				Check: resource.ComposeTestCheckFunc(
					testActionTrailExists("alicloud_actiontrail.foo"),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "name", "tf-testacc-actiontrail"),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "event_rw", "Write-test-123"),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "oss_bucket_name", fmt.Sprintf("tf-testacc-actiontrail-%v", num)),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "role_name", fmt.Sprintf("tf-testacc-actiontrail-%v", num_role)),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "oss_key_prefix", "at-product-account-audit-B"),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "sls_project_arn", ""),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "sls_write_role_arn", ""),
				),
			},
			{
				Config: testActionTrailUpdateConfig_bucket(num_bucket, num_role),
				Check: resource.ComposeTestCheckFunc(
					testActionTrailExists("alicloud_actiontrail.foo"),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "name", "tf-testacc-actiontrail"),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "event_rw", "Write-test-123"),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "oss_bucket_name", fmt.Sprintf("tf-testacc-actiontrail-%v", num_bucket)),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "role_name", fmt.Sprintf("tf-testacc-actiontrail-%v", num_role)),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "oss_key_prefix", "at-product-account-audit-B"),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "sls_project_arn", ""),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "sls_write_role_arn", ""),
				),
			},
			{
				Config: testActionTrailUpdateConfig_sls_project_arn(num_bucket, num_role),
				Check: resource.ComposeTestCheckFunc(
					testActionTrailExists("alicloud_actiontrail.foo"),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "name", "tf-testacc-actiontrail"),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "event_rw", "Write-test-123"),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "oss_bucket_name", fmt.Sprintf("tf-testacc-actiontrail-%v", num_bucket)),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "role_name", fmt.Sprintf("tf-testacc-actiontrail-%v", num_role)),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "oss_key_prefix", "at-product-account-audit-B"),
					resource.TestCheckResourceAttrSet(
						"alicloud_actiontrail.foo", "sls_project_arn"),
					resource.TestCheckResourceAttrSet(
						"alicloud_actiontrail.foo", "sls_write_role_arn"),
				),
			},
			{
				Config: testActionTrailBasicConfig_all(num),
				Check: resource.ComposeTestCheckFunc(
					testActionTrailExists("alicloud_actiontrail.foo"),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "name", "tf-testacc-actiontrail"),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "event_rw", "Write"),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "oss_bucket_name", fmt.Sprintf("tf-testacc-actiontrail-%v", num)),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "role_name", fmt.Sprintf("tf-testacc-actiontrail-%v", num)),
					resource.TestCheckResourceAttr(
						"alicloud_actiontrail.foo", "oss_key_prefix", "at-product-account-audit-A"),
					resource.TestCheckResourceAttrSet(
						"alicloud_actiontrail.foo", "sls_project_arn"),
					resource.TestCheckResourceAttrSet(
						"alicloud_actiontrail.foo", "sls_write_role_arn"),
				),
			},
		},
	})
}

func testActionTrailDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	trailService := ActionTrailService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_actiontrail" {
			continue
		}
		_, err := trailService.DescribeActionTrail(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(fmt.Errorf("Describe Action Trail error %#v", err))
		}
	}
	return nil
}

func testActionTrailExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No action trail is set")
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		trailService := ActionTrailService{client}
		_, err := trailService.DescribeActionTrail(rs.Primary.ID)
		if err != nil {
			return WrapError(err)
		}
		return nil
	}
}

func testActionTrailBasicConfig(randInt int) string {
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
	name = "tf-testacc-actiontrail"
	oss_bucket_name = "${alicloud_oss_bucket.bucket.id}"
	role_name = "${alicloud_ram_role_policy_attachment.attach.role_name}"
}
`, randInt)
}

func testActionTrailUpdateConfig_event_rw(randInt int) string {

	return fmt.Sprintf(`
variable "name" {
    default = "tf-testacc-actiontrail-%v"
}

	resource "alicloud_ram_role" "role" {
	  name = "${var.name}"
	  services = ["actiontrail.aliyuncs.com", "oss.aliyuncs.com"]
	  description = "this is a test"
	  force = true
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
		name = "tf-testacc-actiontrail"
		event_rw = "Write-test-123"
		oss_bucket_name = "${alicloud_oss_bucket.bucket.id}"
		role_name = "${alicloud_ram_role_policy_attachment.attach.role_name}"
	}
`, randInt)

}

func testActionTrailUpdateConfig_oss_key_prefix(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
    default = "tf-testacc-actiontrail-%v"
}

	resource "alicloud_ram_role" "role" {
	  name = "${var.name}"
	  services = ["actiontrail.aliyuncs.com", "oss.aliyuncs.com"]
	  description = "this is a test"
	  force = true
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
		name = "tf-testacc-actiontrail"
		event_rw = "Write-test-123"
		oss_bucket_name = "${alicloud_oss_bucket.bucket.id}"
		role_name = "${alicloud_ram_role_policy_attachment.attach.role_name}"
		oss_key_prefix = "at-product-account-audit-B"
	}
`, randInt)
}

func testActionTrailUpdateConfig_role(rand_1, rand_2 int) string {
	return fmt.Sprintf(`
variable "name" {
    default = "tf-testacc-actiontrail-%v"
}

	resource "alicloud_ram_role" "role" {
	  name = "tf-testacc-actiontrail-%v"
	  services = ["actiontrail.aliyuncs.com", "oss.aliyuncs.com"]
	  description = "this is a test"
	  force = true
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
		name = "tf-testacc-actiontrail"
		event_rw = "Write-test-123"
		oss_bucket_name = "${alicloud_oss_bucket.bucket.id}"
		role_name = "${alicloud_ram_role_policy_attachment.attach.role_name}"
		oss_key_prefix = "at-product-account-audit-B"
	}
`, rand_1, rand_2)
}

func testActionTrailUpdateConfig_bucket(rand_1, rand_2 int) string {
	return fmt.Sprintf(`
variable "name" {
    default = "tf-testacc-actiontrail-%v"
}

	resource "alicloud_ram_role" "role" {
	  name = "tf-testacc-actiontrail-%v"
	  services = ["actiontrail.aliyuncs.com", "oss.aliyuncs.com"]
	  description = "this is a test"
	  force = true
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
		name = "tf-testacc-actiontrail"
		event_rw = "Write-test-123"
		oss_bucket_name = "${alicloud_oss_bucket.bucket.id}"
		role_name = "${alicloud_ram_role_policy_attachment.attach.role_name}"
		oss_key_prefix = "at-product-account-audit-B"
	}
`, rand_1, rand_2)
}

func testActionTrailUpdateConfig_sls_project_arn(rand_1, rand_2 int) string {
	return fmt.Sprintf(`
variable "name" {
    default = "tf-testacc-actiontrail-%v"
}

data "alicloud_regions" "current_region" {
  current = true
}
data "alicloud_account" "current" {
}

resource "alicloud_log_project" "foo" {
  name = "${var.name}"
  description = "tf unit test"
}

resource "alicloud_ram_role" "role" {
	  name = "tf-testacc-actiontrail-%v"
	  services = ["actiontrail.aliyuncs.com", "oss.aliyuncs.com"]
	  description = "this is a test"
	  force = true
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
		name = "tf-testacc-actiontrail"
		event_rw = "Write-test-123"
		oss_bucket_name = "${alicloud_oss_bucket.bucket.id}"
		role_name = "${alicloud_ram_role_policy_attachment.attach.role_name}"
		oss_key_prefix = "at-product-account-audit-B"
	    sls_project_arn = "acs:log:${data.alicloud_regions.current_region.regions.0.id}:${data.alicloud_account.current.id}:project/${alicloud_log_project.foo.name}"
		sls_write_role_arn = "${alicloud_ram_role_policy_attachment.attach.role_name}"
	}
`, rand_1, rand_2)
}

func testActionTrailBasicConfig_all(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
    default = "tf-testacc-actiontrail-%v"
}

data "alicloud_regions" "current_region" {
  current = true
}

data "alicloud_account" "current" {
}

resource "alicloud_log_project" "foo" {
  name = "${var.name}"
  description = "tf unit test"
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
	name = "tf-testacc-actiontrail"
	event_rw = "Write"
	oss_bucket_name = "${alicloud_oss_bucket.bucket.id}"
	role_name = "${alicloud_ram_role_policy_attachment.attach.role_name}"
	oss_key_prefix = "at-product-account-audit-A"
	sls_project_arn = "acs:log:${data.alicloud_regions.current_region.regions.0.id}:${data.alicloud_account.current.id}:project/${alicloud_log_project.foo.name}"
	sls_write_role_arn = "${alicloud_ram_role_policy_attachment.attach.role_name}"
}
`, randInt)
}
