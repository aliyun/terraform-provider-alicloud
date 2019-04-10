package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudDnsRecordsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	domainNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDnsRecordsDataSourceConfig(rand, map[string]string{
			"domain_name": `"${alicloud_dns_record.default.name}"`,
		}),
	}

	hostRecordRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDnsRecordsDataSourceConfig(rand, map[string]string{
			"domain_name":       `"${alicloud_dns_record.default.name}"`,
			"host_record_regex": `"^ali"`,
		}),
		fakeConfig: testAccCheckAlicloudDnsRecordsDataSourceConfig(rand, map[string]string{
			"domain_name":       `"${alicloud_dns_record.default.name}"`,
			"host_record_regex": `"anyother"`,
		}),
	}

	typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDnsRecordsDataSourceConfig(rand, map[string]string{
			"domain_name": `"${alicloud_dns_record.default.name}"`,
			"type":        `"CNAME"`,
		}),
		fakeConfig: testAccCheckAlicloudDnsRecordsDataSourceConfig(rand, map[string]string{
			"domain_name": `"${alicloud_dns_record.default.name}"`,
			"type":        `"TXT"`,
		}),
	}

	valueRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDnsRecordsDataSourceConfig(rand, map[string]string{
			"domain_name": `"${alicloud_dns_record.default.name}"`,
			"value_regex": `"^mail"`,
		}),
		fakeConfig: testAccCheckAlicloudDnsRecordsDataSourceConfig(rand, map[string]string{
			"domain_name": `"${alicloud_dns_record.default.name}"`,
			"value_regex": `"anyother"`,
		}),
	}

	lineConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDnsRecordsDataSourceConfig(rand, map[string]string{
			"domain_name": `"${alicloud_dns_record.default.name}"`,
			"line":        `"default"`,
		}),
		fakeConfig: testAccCheckAlicloudDnsRecordsDataSourceConfig(rand, map[string]string{
			"domain_name": `"${alicloud_dns_record.default.name}"`,
			"line":        `"telecom"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDnsRecordsDataSourceConfig(rand, map[string]string{
			"domain_name": `"${alicloud_dns_record.default.name}"`,
			"status":      `"enable"`,
		}),
		fakeConfig: testAccCheckAlicloudDnsRecordsDataSourceConfig(rand, map[string]string{
			"domain_name": `"${alicloud_dns_record.default.name}"`,
			"status":      `"disable"`,
		}),
	}

	isLockConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDnsRecordsDataSourceConfig(rand, map[string]string{
			"domain_name": `"${alicloud_dns_record.default.name}"`,
			"is_locked":   `"false"`,
		}),
		fakeConfig: testAccCheckAlicloudDnsRecordsDataSourceConfig(rand, map[string]string{
			"domain_name": `"${alicloud_dns_record.default.name}"`,
			"is_locked":   `"true"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDnsRecordsDataSourceConfig(rand, map[string]string{
			"domain_name":       `"${alicloud_dns_record.default.name}"`,
			"host_record_regex": `"^ali"`,
			"value_regex":       `"^mail"`,
			"type":              `"CNAME"`,
			"line":              `"default"`,
			"status":            `"enable"`,
			"is_locked":         `"false"`,
		}),
		fakeConfig: testAccCheckAlicloudDnsRecordsDataSourceConfig(rand, map[string]string{
			"domain_name":       `"${alicloud_dns_record.default.name}"`,
			"host_record_regex": `"^ali"`,
			"value_regex":       `"^mail"`,
			"type":              `"CNAME"`,
			"line":              `"default"`,
			"status":            `"enable"`,
			"is_locked":         `"true"`,
		}),
	}

	dnsRecordsCheckInfo.dataSourceTestCheck(t, rand, domainNameConf, hostRecordRegexConf, typeConf, valueRegexConf, valueRegexConf,
		lineConf, statusConf, isLockConf, allConf)
}

func testAccCheckAlicloudDnsRecordsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
resource "alicloud_dns" "default" {
  name = "tf-testaccdns%v.abc"
}

resource "alicloud_dns_record" "default" {
  name = "${alicloud_dns.default.name}"
  host_record = "alimail"
  type = "CNAME"
  value = "mail.mxhichin.com"
  count = 1
}

data "alicloud_dns_records" "default" {
  %s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existDnsRecordsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"urls.#":                "1",
		"urls.0":                fmt.Sprintf("%v.%v", "alimail", fmt.Sprintf("tf-testaccdns%d.abc", rand)),
		"records.#":             "1",
		"domain_name":           fmt.Sprintf("tf-testaccdns%d.abc", rand),
		"records.0.locked":      "false",
		"records.0.host_record": "alimail",
		"records.0.type":        "CNAME",
		"records.0.value":       "mail.mxhichin.com",
		"records.0.record_id":   CHECKSET,
		"records.0.ttl":         "600",
		"records.0.priority":    "0",
		"records.0.line":        "default",
		"records.0.status":      "enable",
	}
}

var fakeDnsRecordsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"urls.#":    "0",
		"records.#": "0",
	}
}

var dnsRecordsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_dns_records.default",
	existMapFunc: existDnsRecordsMapFunc,
	fakeMapFunc:  fakeDnsRecordsMapFunc,
}
