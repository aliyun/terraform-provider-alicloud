package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudFcFunctionsDataSource_basic(t *testing.T) {
	randInt := acctest.RandInt()
	functionName := fmt.Sprintf("tf-testacc-fc-function-ds-basic-%d", randInt)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudFcFunctionsDataSourceBasic(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_fc_functions.functions"),
					resource.TestCheckResourceAttr("data.alicloud_fc_functions.functions", "functions.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_fc_functions.functions", "functions.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_fc_functions.functions", "functions.0.name", functionName),
					resource.TestCheckResourceAttr("data.alicloud_fc_functions.functions", "functions.0.description", functionName+"-description"),
					resource.TestCheckResourceAttr("data.alicloud_fc_functions.functions", "functions.0.runtime", "python2.7"),
					resource.TestCheckResourceAttr("data.alicloud_fc_functions.functions", "functions.0.handler", "hello.handler"),
					resource.TestCheckResourceAttr("data.alicloud_fc_functions.functions", "functions.0.timeout", "120"),
					resource.TestCheckResourceAttr("data.alicloud_fc_functions.functions", "functions.0.memory_size", "512"),
					resource.TestCheckResourceAttr("data.alicloud_fc_functions.functions", "functions.0.code_size", "105"),
					resource.TestCheckResourceAttr("data.alicloud_fc_functions.functions", "functions.0.code_checksum", "5237022206872530469"),
					resource.TestCheckResourceAttrSet("data.alicloud_fc_functions.functions", "functions.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_fc_functions.functions", "functions.0.last_modification_time"),
				),
			},
		},
	})
}

func TestAccAlicloudFcFunctionsDataSource_empty(t *testing.T) {
	randInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudFcFunctionsDataSourceEmpty(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_fc_functions.functions"),
					resource.TestCheckResourceAttr("data.alicloud_fc_functions.functions", "functions.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_fc_functions.functions", "functions.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_fc_functions.functions", "functions.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_fc_functions.functions", "functions.0.description"),
					resource.TestCheckNoResourceAttr("data.alicloud_fc_functions.functions", "functions.0.runtime"),
					resource.TestCheckNoResourceAttr("data.alicloud_fc_functions.functions", "functions.0.handler"),
					resource.TestCheckNoResourceAttr("data.alicloud_fc_functions.functions", "functions.0.timeout"),
					resource.TestCheckNoResourceAttr("data.alicloud_fc_functions.functions", "functions.0.memory_size"),
					resource.TestCheckNoResourceAttr("data.alicloud_fc_functions.functions", "functions.0.code_size"),
					resource.TestCheckNoResourceAttr("data.alicloud_fc_functions.functions", "functions.0.code_checksum"),
					resource.TestCheckNoResourceAttr("data.alicloud_fc_functions.functions", "functions.0.creation_time"),
					resource.TestCheckNoResourceAttr("data.alicloud_fc_functions.functions", "functions.0.last_modification_time"),
				),
			},
		},
	})
}

func testAccCheckAlicloudFcFunctionsDataSourceBasic(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testacc-fc-function-ds-basic-%d"
}

resource "alicloud_fc_service" "sample_service" {
    name = "${var.name}"
}

resource "alicloud_oss_bucket" "sample_bucket" {
	bucket = "${var.name}"
}

resource "alicloud_oss_bucket_object" "sample_object" {
	bucket = "${alicloud_oss_bucket.sample_bucket.id}"
	key = "fc/hello.zip"
	content = <<EOF
		# -*- coding: utf-8 -*-
		def handler(event, context):
			print "hello world"
			return 'hello world'
	EOF
}

resource "alicloud_fc_function" "sample_function" {
	service = "${alicloud_fc_service.sample_service.name}"
	name = "${var.name}"
	description = "${var.name}-description"
	oss_bucket = "${alicloud_oss_bucket.sample_bucket.id}"
	oss_key = "${alicloud_oss_bucket_object.sample_object.key}"
	memory_size = "512"
	runtime = "python2.7"
	handler = "hello.handler"
	timeout = "120"
}

data "alicloud_fc_functions" "functions" {
	service_name = "${alicloud_fc_service.sample_service.name}"
    name_regex = "${alicloud_fc_function.sample_function.name}"
}
`, randInt)
}

func testAccCheckAlicloudFcFunctionsDataSourceEmpty(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testacc-fc-function-ds-basic-%d"
}

resource "alicloud_fc_service" "sample_service" {
    name = "${var.name}"
}

data "alicloud_fc_functions" "functions" {
    service_name = "${alicloud_fc_service.sample_service.name}"
    name_regex = "^tf-testacc-fake-name*"
}
`, randInt)
}
