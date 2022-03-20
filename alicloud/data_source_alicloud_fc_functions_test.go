package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
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

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"service_name": "${alicloud_fc_function.default.service}",
			"ids":          []string{"${alicloud_fc_function.default.function_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"service_name": "${alicloud_fc_function.default.service}",
			"ids":          []string{"${alicloud_fc_function.default.function_id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"service_name": "${alicloud_fc_function.default.service}",
			"name_regex":   "${alicloud_fc_function.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"service_name": "${alicloud_fc_function.default.service}",
			"name_regex":   "${alicloud_fc_function.default.name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"service_name": "${alicloud_fc_function.default.service}",
			"ids":          []string{"${alicloud_fc_function.default.function_id}"},
			"name_regex":   "${alicloud_fc_function.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"service_name": "${alicloud_fc_function.default.service}",
			"ids":          []string{"${alicloud_fc_function.default.function_id}"},
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
			"functions.0.code_size":                    "155",
			"functions.0.code_checksum":                "13713894526498940759",
			"functions.0.creation_time":                CHECKSET,
			"functions.0.last_modification_time":       CHECKSET,
			"functions.0.environment_variables.test":   "terraform",
			"functions.0.environment_variables.prefix": "tfAcc",
			"functions.0.initializer":                  "hello.initializer",
			"functions.0.initialization_timeout":       "30",
			"functions.0.instance_type":                "e1",
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

	fcFunctionsCheckInfo.dataSourceTestCheck(t, rand, serviceNameConf, idsConf, nameRegexConf, allConf)
}

func TestAccAlicloudFCFunctionsDataSourceCustomContainerUpdate(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_fc_functions.default"

	name := fmt.Sprintf("tf-testaccfcfunctiondscustomcontainer-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlicloudFCFunctionsDataSourceCustomContainerConfig)

	serviceNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"service_name": "${alicloud_fc_function.default.service}",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"service_name": "${alicloud_fc_function.default.service}",
			"ids":          []string{"${alicloud_fc_function.default.function_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"service_name": "${alicloud_fc_function.default.service}",
			"ids":          []string{"${alicloud_fc_function.default.function_id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"service_name": "${alicloud_fc_function.default.service}",
			"name_regex":   "${alicloud_fc_function.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"service_name": "${alicloud_fc_function.default.service}",
			"name_regex":   "${alicloud_fc_function.default.name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"service_name": "${alicloud_fc_function.default.service}",
			"ids":          []string{"${alicloud_fc_function.default.function_id}"},
			"name_regex":   "${alicloud_fc_function.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"service_name": "${alicloud_fc_function.default.service}",
			"ids":          []string{"${alicloud_fc_function.default.function_id}"},
			"name_regex":   "${alicloud_fc_function.default.name}_fake",
		}),
	}

	var existFCFunctionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"names.#":             "1",
			"functions.#":         "1",
			"functions.0.name":    name,
			"functions.0.runtime": "custom-container",
			"functions.0.custom_container_config.0.image":   fmt.Sprintf("registry.%s.aliyuncs.com/eci_open/nginx:alpine", defaultRegionToTest),
			"functions.0.custom_container_config.0.command": `["python", "server.py"]`,
			"functions.0.custom_container_config.0.args":    `["a1", "a2"]`,
			"functions.0.instance_concurrency":              "28",
			"functions.0.ca_port":                           "9527",
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

	fcFunctionsCheckInfo.dataSourceTestCheck(t, rand, serviceNameConf, idsConf, nameRegexConf, allConf)
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
		def initializer(context):
			print "hello init"

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
    environment_variables = {
     test = "terraform"
     prefix = "tfAcc"
	}
	initializer = "hello.initializer"
	initialization_timeout = "30"
	instance_type = "e1"
}
`, name)
}

func testAccCheckAlicloudFCFunctionsDataSourceCustomContainerConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_fc_service" "default" {
    name = "${var.name}"
    role = "${alicloud_ram_role.default.arn}"
    depends_on = ["alicloud_ram_role_policy_attachment.default"]
}

locals {
	container_command = "[\"python\", \"server.py\"]"
	container_args = "[\"a1\", \"a2\"]"
}

output "container_command" {
	value = "${local.container_command}"
}

output "container_args" {
	value = "${local.container_args}"
}

resource "alicloud_fc_function" "default" {
	service = "${alicloud_fc_service.default.name}"
	name = "${var.name}"
	handler = "fake"
	memory_size = "512"
	runtime = "custom-container"
	instance_concurrency = "28"
	ca_port = "9527"
	custom_container_config {
		image = "registry.%s.aliyuncs.com/eci_open/nginx:alpine"
		command = "${local.container_command}"
		args = "${local.container_args}"
	}
}

resource "alicloud_ram_role" "default" {
  name = "${var.name}"
  document = <<EOF
  %s
  EOF
  description = "this is a test"
  force = true
}
resource "alicloud_ram_role_policy_attachment" "default" {
  role_name = "${alicloud_ram_role.default.name}"
  policy_name = "AliyunContainerRegistryReadOnlyAccess"
  policy_type = "System"
}
`, name, defaultRegionToTest, testFCRoleTemplate)
}
