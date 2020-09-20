package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudFCCustomDomainsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_fc_custom_domains.default"
	name := fmt.Sprintf("tf-testacc-fc-custom-domains-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceFCCustomDomainsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_fc_custom_domain.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_fc_custom_domain.default.name}_fake",
		}),
	}

	var existFCCustomDomainsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"domains.#":                              "1",
			"names.#":                                "1",
			"domains.0.domain_name":                  "terraform.functioncompute.com",
			"domains.0.protocol":                     "HTTP",
			"domains.0.account_id":                   CHECKSET,
			"domains.0.api_version":                  CHECKSET,
			"domains.0.created_time":                 CHECKSET,
			"domains.0.last_modified_time":           CHECKSET,
			"domains.0.route_config.0.path":          "/*",
			"domains.0.route_config.0.service_name":  name,
			"domains.0.route_config.0.function_name": name,
			"domains.0.route_config.0.qualifier":     "v1",
			"domains.0.route_config.0.methods.0":     "GET",
			"domains.0.route_config.0.methods.1":     "POST",
		}
	}

	var fakeFCCustomDomainsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"domains.#": "0",
			"names.#":   "0",
		}
	}

	var fcCustomDomainsRecordsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existFCCustomDomainsMapFunc,
		fakeMapFunc:  fakeFCCustomDomainsMapFunc,
	}

	fcCustomDomainsRecordsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf)
}

func dataSourceFCCustomDomainsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_fc_custom_domain" "default" {
	domain_name = "terraform.functioncompute.com"
	protocol = "HTTP"
	route_config {
		path = "/*"
		service_name = "${alicloud_fc_service.default.name}"
		function_name = "${alicloud_fc_function.default.name}"
		qualifier = "v1"
		methods = ["GET","POST"]
	}
}

resource "alicloud_fc_service" "default" {
    name = "${var.name}"
    description = "${var.name}-description"
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
	oss_bucket = "${alicloud_oss_bucket.default.id}"
	oss_key = "${alicloud_oss_bucket_object.default.key}"
	memory_size = 512
	runtime = "python2.7"
	handler = "hello.handler"
}
`, name)
}
