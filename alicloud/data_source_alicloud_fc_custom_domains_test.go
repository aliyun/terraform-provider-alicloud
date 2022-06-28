package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudFCCustomDomainsDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.FCCustomDomainSupportRegions)
	rand := acctest.RandInt()
	resourceId := "data.alicloud_fc_custom_domains.default"
	name := fmt.Sprintf("tf-testacc-fc-custom-domains-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceFCCustomDomainsConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_fc_custom_domain.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_fc_custom_domain.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_fc_custom_domain.default.domain_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_fc_custom_domain.default.domain_name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_fc_custom_domain.default.domain_name}",
			"ids":        []string{"${alicloud_fc_custom_domain.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_fc_custom_domain.default.domain_name}_fake",
			"ids":        []string{"${alicloud_fc_custom_domain.default.id}_fake"},
		}),
	}

	var existFCCustomDomainsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"domains.#":                              "1",
			"ids.#":                                  "1",
			"names.#":                                "1",
			"domains.0.id":                           "alicloud-provider.shop",
			"domains.0.domain_name":                  "alicloud-provider.shop",
			"domains.0.protocol":                     "HTTP,HTTPS",
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
			"domains.0.cert_config.0.cert_name":      "test",
			"domains.0.cert_config.0.certificate":    strings.Replace(testFcCertificate, `\n`, "\n", -1),
		}
	}

	var fakeFCCustomDomainsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"domains.#": "0",
			"ids.#":     "0",
			"names.#":   "0",
		}
	}

	var fcCustomDomainsRecordsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existFCCustomDomainsMapFunc,
		fakeMapFunc:  fakeFCCustomDomainsMapFunc,
	}

	fcCustomDomainsRecordsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

func dataSourceFCCustomDomainsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_fc_custom_domain" "default" {
	domain_name = "alicloud-provider.shop"
	protocol = "HTTP,HTTPS"
	route_config {
		path = "/*"
		service_name = "${alicloud_fc_service.default.name}"
		function_name = "${alicloud_fc_function.default.name}"
		qualifier = "v1"
		methods = ["GET","POST"]
	}
	cert_config {
		cert_name = "test"
		private_key = "%s"
		certificate = "%s"
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
`, name, testFcPrivateKey, testFcCertificate)
}
