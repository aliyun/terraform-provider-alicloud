package alidns_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudAlidnsResolutionLinesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sdns%v.abc", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc("data.alicloud_dns_resolution_lines.default", name, dataSourceDnsResolutionLinesConfigDependence)
	lineCodesConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"line_codes": []string{"${alicloud_dns_record.default.routing}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"line_codes": []string{"${alicloud_dns_record.default.routing}_fake"},
		}),
	}
	lineDisplayNamesConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"line_display_names": []string{"中国联通"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"line_display_names": []string{"中国联通_fake"},
		}),
	}
	lineNamesConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"line_codes": []string{"${alicloud_dns_record.default.routing}"},
			"line_names": []string{"中国联通"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"line_codes": []string{"${alicloud_dns_record.default.routing}"},
			"line_names": []string{"中国联通-fake"},
		}),
	}

	domainNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"line_codes":  []string{"${alicloud_dns_record.default.routing}"},
			"domain_name": "${alicloud_dns_record.default.name}",
		}),
		//FakeConfig: testAccConfig(map[string]interface{}{
		//	"line_codes":  []string{"${alicloud_dns_record.default.routing}"},
		//	"domain_name": "${alicloud_dns_record.default.name}_fake",
		//}),
	}

	userClientIpConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"line_codes":     []string{"${alicloud_dns_record.default.routing}"},
			"user_client_ip": "205.204.117.106",
		}),
		/*FakeConfig: testAccConfig(map[string]interface{}{
			"line_codes":     []string{"cn_unicom_shanxi"},
			"user_client_ip": "",
		}),*/
	}
	langConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"line_codes": []string{"${alicloud_dns_record.default.routing}"},
			"lang":       "zh",
		}),
		// if the lang were set fake, it will be reset as default en(english)
		/*FakeConfig: testAccConfig(map[string]interface{}{
			"line_codes": []string{"cn_unicom_shanxi"},
			"lang":       "zh_fake",
		}),*/
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"line_codes":         []string{"${alicloud_dns_record.default.routing}"},
			"line_display_names": []string{"中国联通"},
			"domain_name":        "${alicloud_dns_record.default.name}",
			"user_client_ip":     "205.204.117.106",
			"lang":               "zh",
			"line_names":         []string{"中国联通"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"line_codes":         []string{"${alicloud_dns_record.default.routing}"},
			"line_display_names": []string{"中国联通"},
			"domain_name":        "${alicloud_dns_record.default.name}",
			"user_client_ip":     "205.204.117.106",
			"lang":               "zh",
			"line_names":         []string{"中国联通_fake"},
		}),
	}
	dnsSupportLinesCheckInfo.DataSourceTestCheck(t, rand, lineCodesConf, lineDisplayNamesConf, lineNamesConf, domainNameConf, userClientIpConf, langConf, allConf)
}

func dataSourceDnsResolutionLinesConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "default" {
  name = "%s"
}

resource "alicloud_dns_record" "default" {
  name = "${alicloud_dns.default.name}"
  host_record = "alimail"
  routing = "unicom"
  type = "CNAME"
  value = "mail.mxhichin.com"
}
`, name)
}

var existDnsResolutionLinesMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"lines.#":                   "1",
		"lines.0.line_code":         "unicom",
		"lines.0.line_display_name": "中国联通",
		"lines.0.line_name":         "中国联通",
		"line_codes.#":              "1",
		"line_codes.0":              "unicom",
		"line_display_names.#":      "1",
		"line_display_names.0":      "中国联通",
	}
}

var fakeDnsResolutionLinesMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"lines.#":              "0",
		"line_codes.#":         "0",
		"line_display_names.#": "0",
	}
}

var dnsSupportLinesCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_dns_resolution_lines.default",
	ExistMapFunc: existDnsResolutionLinesMapCheck,
	FakeMapFunc:  fakeDnsResolutionLinesMapCheck,
}
