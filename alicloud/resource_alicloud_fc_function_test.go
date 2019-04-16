package alicloud

import (
	"fmt"
	"testing"

	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/fc-go-sdk"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudFCFunction_basic(t *testing.T) {
	var service fc.GetServiceOutput
	var function fc.GetFunctionOutput
	var bucket oss.BucketInfo

	randInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudFCFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAlicloudFCFunctionBasic(testFCRoleTemplate, randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudFCServiceExists("alicloud_fc_service.foo", &service),
					testAccCheckAlicloudFCFunctionExists("alicloud_fc_function.foo", &function),
					resource.TestCheckResourceAttr("alicloud_fc_function.foo", "name", fmt.Sprintf("tf-testacc-alicloud-fc-function-basic-%v", randInt)),
					resource.TestCheckResourceAttr("alicloud_fc_function.foo", "description", "tf unit test"),
					resource.TestCheckResourceAttr("alicloud_fc_function.foo", "runtime", "python2.7"),
					resource.TestCheckResourceAttr("alicloud_fc_function.foo", "memory_size", "512"),
					resource.TestCheckResourceAttr("alicloud_fc_function.foo", "environment_variables.test", `terraform`),
					resource.TestCheckResourceAttr("alicloud_fc_function.foo", "environment_variables.prefix", `tfAcc`),
				),
			},
			{
				Config: testAlicloudFCFunctionUpdate(testFCRoleTemplate, randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssBucketExists("alicloud_oss_bucket.foo", &bucket),
					testAccCheckAlicloudFCFunctionExists("alicloud_fc_function.foo", &function),
					resource.TestCheckResourceAttr("alicloud_fc_function.foo", "name", fmt.Sprintf("tf-testacc-alicloud-fc-function-basic-%v", randInt)),
					resource.TestCheckResourceAttr("alicloud_fc_function.foo", "description", "tf unit test"),
					resource.TestCheckResourceAttr("alicloud_fc_function.foo", "runtime", "nodejs6"),
					resource.TestCheckResourceAttr("alicloud_fc_function.foo", "memory_size", "128"),
					resource.TestCheckResourceAttr("alicloud_fc_function.foo", "environment_variables.test", `success`),
					resource.TestCheckResourceAttr("alicloud_fc_function.foo", "environment_variables.prefix", `tfAcc-test`),
				),
			},
		},
	})
}

func testAccCheckAlicloudFCFunctionExists(name string, service *fc.GetFunctionOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", name))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No Log store ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		fcService := FcService{client}
		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		ser, err := fcService.DescribeFcFunction(split[0], split[1])
		if err != nil {
			return WrapError(err)
		}

		service = ser

		return nil
	}
}

func testAccCheckAlicloudFCFunctionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	fcService := FcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_fc_function" {
			continue
		}

		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		if _, err := fcService.DescribeFcFunction(split[0], split[1]); err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Check fc service got an error: %#v.", err)
		}

		return fmt.Errorf("FC service %s still exists.", rs.Primary.ID)
	}

	return nil
}

func testAlicloudFCFunctionBasic(role string, randInt int) string {
	return fmt.Sprintf(`
variable "name" {
    default = "tf-testacc-alicloud-fc-function-basic-%v"
}
resource "alicloud_log_project" "foo" {
  name = "${var.name}"
  description = "tf unit test"
}

resource "alicloud_log_store" "foo" {
  project = "${alicloud_log_project.foo.name}"
  name = "${var.name}"
  retention_period = "3000"
  shard_count = 1
}
resource "alicloud_fc_service" "foo" {
    name = "${var.name}"
    description = "tf unit test"
    log_config {
	project = "${alicloud_log_project.foo.name}"
	logstore = "${alicloud_log_store.foo.name}"
    }
    role = "${alicloud_ram_role.foo.arn}"
    depends_on = ["alicloud_ram_role_policy_attachment.foo"]
}
resource "alicloud_oss_bucket" "foo" {
  bucket = "${var.name}"
}

resource "alicloud_oss_bucket_object" "foo" {
  bucket = "${alicloud_oss_bucket.foo.id}"
  key = "fc/hello.zip"
  content = <<EOF
  	# -*- coding: utf-8 -*-
	def handler(event, context):
	    print "hello world"
	    return 'hello world'
  EOF
}

resource "alicloud_fc_function" "foo" {
  service = "${alicloud_fc_service.foo.name}"
  name = "${var.name}"
  description = "tf unit test"
  oss_bucket = "${alicloud_oss_bucket.foo.id}"
  oss_key = "${alicloud_oss_bucket_object.foo.key}"
  memory_size = "512"
  runtime = "python2.7"
  handler = "hello.handler"
  environment_variables {
     test = "terraform"
     prefix = "tfAcc"
  }
}

resource "alicloud_ram_role" "foo" {
  name = "${var.name}"
  document = <<EOF
  %s
  EOF
  description = "this is a test"
  force = true
}
resource "alicloud_ram_role_policy_attachment" "foo" {
  role_name = "${alicloud_ram_role.foo.name}"
  policy_name = "AliyunLogFullAccess"
  policy_type = "System"
}
`, randInt, role)
}

func testAlicloudFCFunctionUpdate(role string, randInt int) string {
	return fmt.Sprintf(`
variable "name" {
    default = "tf-testacc-alicloud-fc-function-basic-%v"
}
resource "alicloud_log_project" "foo" {
  name = "${var.name}"
  description = "tf unit test"
}

resource "alicloud_log_store" "foo" {
  project = "${alicloud_log_project.foo.name}"
  name = "${var.name}"
  retention_period = "3000"
  shard_count = 1
}
resource "alicloud_fc_service" "foo" {
    name = "${var.name}"
    description = "tf unit test"
    log_config {
	project = "${alicloud_log_project.foo.name}"
	logstore = "${alicloud_log_store.foo.name}"
    }
    role = "${alicloud_ram_role.foo.arn}"
    depends_on = ["alicloud_ram_role_policy_attachment.foo"]
}

resource "alicloud_oss_bucket" "foo" {
  bucket = "${var.name}"
}

resource "alicloud_oss_bucket_object" "foo" {
  bucket = "${alicloud_oss_bucket.foo.id}"
  key = "fc/hello.zip"
  content = <<EOF
  	# -*- coding: utf-8 -*-
	def handler(event, context):
	    print "hello world"
	    return 'hello world'
  EOF
}

resource "alicloud_fc_function" "foo" {
  service = "${alicloud_fc_service.foo.name}"
  description = "tf unit test"
  name = "${var.name}"
  oss_bucket = "${alicloud_oss_bucket.foo.id}"
  oss_key = "${alicloud_oss_bucket_object.foo.key}"
  runtime = "nodejs6"
  handler = "hello.handler"
  environment_variables {
     test = "success"
     prefix = "tfAcc-test"
  }
}
resource "alicloud_ram_role" "foo" {
  name = "${var.name}"
  document = <<EOF
  %s
  EOF
  description = "this is a test"
  force = true
}

resource "alicloud_ram_role_policy_attachment" "foo" {
  role_name = "${alicloud_ram_role.foo.name}"
  policy_name = "AliyunLogFullAccess"
  policy_type = "System"
}
`, randInt, role)
}
