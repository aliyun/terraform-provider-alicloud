package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudFCFunctionsDataSourceUpdate(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_fc_functions.default"

	name := fmt.Sprintf("tf-testaccfcfunctiondsbasic-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlicloudFCFunctionsDataSourceConfig)

	serviceNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"service_name": "${alicloud_fc_function.default.service}",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"service_name": "${alicloud_fc_function.default.service}",
			"name_regex":   "${alicloud_fc_function.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"service_name": "${alicloud_fc_function.default.service}",
			"name_regex":   "${alicloud_fc_function.default.name}_fake",
		}),
	}

	var existFCFunctionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                    "1",
			"names.#":                                  "1",
			"functions.#":                              "1",
			"functions.0.name":                         name,
			"functions.0.description":                  name + "-description",
			"functions.0.runtime":                      "python2.7",
			"functions.0.handler":                      "hello.handler",
			"functions.0.timeout":                      "120",
			"functions.0.memory_size":                  "512",
			"functions.0.code_size":                    "105",
			"functions.0.code_checksum":                "5237022206872530469",
			"functions.0.creation_time":                CHECKSET,
			"functions.0.last_modification_time":       CHECKSET,
			"functions.0.environment_variables.test":   "terraform",
			"functions.0.environment_variables.prefix": "tfAcc",
		}
	}

	var fakeFCFunctionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"functions.#": "0",
			"ids.#":       "0",
			"names.#":     "0",
		}
	}

	var fcFunctionsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existFCFunctionsMapFunc,
		fakeMapFunc:  fakeFCFunctionsMapFunc,
	}

	fcFunctionsCheckInfo.dataSourceTestCheck(t, rand, serviceNameConf, allConf)
}

func testAccCheckAlicloudFCFunctionsDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_fc_service" "default" {
    name = "${var.name}"
}

resource "alicloud_oss_bucket" "default" {
	bucket = "${var.name}"
}

resource "alicloud_oss_bucket_object" "default" {
	bucket = "${alicloud_oss_bucket.default.id}"
	key = "fc/hello.zip"
	content = <<EOF
		# -*- coding: utf-8 -*-
		def handler(event, context):
			print "hello world"
			return 'hello world'
	EOF
}

resource "alicloud_fc_function" "default" {
	service = "${alicloud_fc_service.default.name}"
	name = "${var.name}"
	description = "${var.name}-description"
	oss_bucket = "${alicloud_oss_bucket.default.id}"
	oss_key = "${alicloud_oss_bucket_object.default.key}"
	memory_size = "512"
	runtime = "python2.7"
	handler = "hello.handler"
	timeout = "120"
    environment_variables {
     test = "terraform"
     prefix = "tfAcc"
  }
}
`, name)
}
