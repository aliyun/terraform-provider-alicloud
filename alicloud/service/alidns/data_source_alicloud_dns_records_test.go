package alidns_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudAlidnsRecordsDataSourceOld(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_dns_records.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc%sdns%v.abc", acctest.DefaultRegionToTest, rand),
		dataSourceDnsRecordsConfigDependence)

	domainNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_dns_record.default.name}",
		}),
	}

	hostRecordRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"domain_name":       "${alicloud_dns_record.default.name}",
			"host_record_regex": "^ali",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"domain_name":       "${alicloud_dns_record.default.name}",
			"host_record_regex": "anyother",
		}),
	}

	typeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_dns_record.default.name}",
			"type":        "CNAME",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_dns_record.default.name}",
			"type":        "TXT",
		}),
	}

	valueRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_dns_record.default.name}",
			"value_regex": "^mail",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_dns_record.default.name}",
			"value_regex": "anyother",
		}),
	}

	lineConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_dns_record.default.name}",
			"line":        "default",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_dns_record.default.name}",
			"line":        "telecom",
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_dns_record.default.name}",
			"status":      "enable",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_dns_record.default.name}",
			"status":      "disable",
		}),
	}

	isLockConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_dns_record.default.name}",
			"is_locked":   "false",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_dns_record.default.name}",
			"is_locked":   "true",
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_dns_record.default.name}",
			"ids":         []string{"${alicloud_dns_record.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_dns_record.default.name}",
			"ids":         []string{"${alicloud_dns_record.default.id}-fake"},
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"domain_name":       "${alicloud_dns_record.default.name}",
			"host_record_regex": "^ali",
			"value_regex":       "^mail",
			"type":              "CNAME",
			"line":              "default",
			"status":            "enable",
			"is_locked":         "false",
			"ids":               []string{"${alicloud_dns_record.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"domain_name":       "${alicloud_dns_record.default.name}",
			"host_record_regex": "^ali",
			"value_regex":       "^mail",
			"type":              "CNAME",
			"line":              "default",
			"status":            "enable",
			"is_locked":         "true",
			"ids":               []string{"${alicloud_dns_record.default.id}"},
		}),
	}

	var existDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 "1",
			"ids.0":                 acctest.CHECKSET,
			"urls.#":                "1",
			"urls.0":                fmt.Sprintf("alimail.tf-testacc%sdns%d.abc", acctest.DefaultRegionToTest, rand),
			"records.#":             "1",
			"domain_name":           fmt.Sprintf("tf-testacc%sdns%d.abc", acctest.DefaultRegionToTest, rand),
			"records.0.locked":      "false",
			"records.0.host_record": "alimail",
			"records.0.type":        "CNAME",
			"records.0.value":       "mail.mxhichin.com",
			"records.0.record_id":   acctest.CHECKSET,
			"records.0.ttl":         "600",
			"records.0.priority":    "0",
			"records.0.line":        "default",
			"records.0.status":      "enable",
		}
	}

	var fakeDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"urls.#":    "0",
			"records.#": "0",
		}
	}

	var dnsRecordsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existDnsRecordsMapFunc,
		FakeMapFunc:  fakeDnsRecordsMapFunc,
	}

	dnsRecordsCheckInfo.DataSourceTestCheck(t, rand, domainNameConf, hostRecordRegexConf, typeConf, valueRegexConf, valueRegexConf,
		lineConf, statusConf, isLockConf, idsConf, allConf)
}

func dataSourceDnsRecordsConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "default" {
  name = "%s"
}

resource "alicloud_dns_record" "default" {
  name = "${alicloud_dns.default.name}"
  host_record = "alimail"
  type = "CNAME"
  value = "mail.mxhichin.com"
}
`, name)
}
