package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAlidnsRecordsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_alidns_records.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc%salidns%v.abc", defaultRegionToTest, rand),
		dataSourceAlidnsRecordsConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"ids":         []string{"${alicloud_alidns_record.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"ids":         []string{"${alicloud_alidns_record.default.id}-fake"},
		}),
	}

	rrRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"rr_regex":    "^ali",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"rr_regex":    "anyother",
		}),
	}

	valueRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"value_regex": "^mail",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"value_regex": "anyother",
		}),
	}

	typeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"type":        "CNAME",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"type":        "TXT",
		}),
	}

	lineConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"line":        "default",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"line":        "telecom",
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"status":      "ENABLE",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"status":      "DISABLE",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"type":        "CNAME",
			"line":        "default",
			"status":      "ENABLE",
			"ids":         []string{"${alicloud_alidns_record.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"type":        "CNAME",
			"line":        "default",
			"status":      "ENABLE",
			"ids":         []string{"${alicloud_alidns_record.default.id}-fake"},
		}),
	}

	var existAlidnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"ids.0":               CHECKSET,
			"records.#":           "1",
			"domain_name":         fmt.Sprintf("tf-testacc%salidns%v.abc", defaultRegionToTest, rand),
			"records.0.locked":    "false",
			"records.0.rr":        "alimail",
			"records.0.type":      "CNAME",
			"records.0.value":     "mail.mxhichin.com",
			"records.0.record_id": CHECKSET,
			"records.0.ttl":       "600",
			"records.0.priority":  "0",
			"records.0.line":      "default",
			"records.0.status":    "ENABLE",
			"records.0.remark":    "tf-testacc",
		}
	}

	var fakeAlidnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"records.#": "0",
		}
	}

	var alidnsRecordsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAlidnsRecordsMapFunc,
		fakeMapFunc:  fakeAlidnsRecordsMapFunc,
	}

	alidnsRecordsCheckInfo.dataSourceTestCheck(t, rand, typeConf, rrRegexConf, valueRegexConf, lineConf, statusConf, idsConf, allConf)
}

func dataSourceAlidnsRecordsConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_dns_domain" "default" {
  domain_name = "%s"
}

resource "alicloud_alidns_record" "default" {
  domain_name = "${alicloud_dns_domain.default.domain_name}"
  rr = "alimail"
  type = "CNAME"
  remark = "tf-testacc"
  value = "mail.mxhichin.com"
}
`, name)
}
