package alidns_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudAlidnsRecordsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_alidns_records.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc%salidns%v.abc", acctest.DefaultRegionToTest, rand),
		dataSourceAlidnsRecordsConfigDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"ids":         []string{"${alicloud_alidns_record.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"ids":         []string{"${alicloud_alidns_record.default.id}-fake"},
		}),
	}

	rrRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"rr_regex":    "^ali",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"rr_regex":    "anyother",
		}),
	}

	valueRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"value_regex": "^mail",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"value_regex": "anyother",
		}),
	}

	typeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"type":        "CNAME",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"type":        "TXT",
		}),
	}

	lineConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"line":        "default",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"line":        "telecom",
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"status":      "ENABLE",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"status":      "DISABLE",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_alidns_record.default.domain_name}",
			"type":        "CNAME",
			"line":        "default",
			"status":      "ENABLE",
			"ids":         []string{"${alicloud_alidns_record.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
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
			"ids.0":               acctest.CHECKSET,
			"records.#":           "1",
			"domain_name":         fmt.Sprintf("tf-testacc%salidns%v.abc", acctest.DefaultRegionToTest, rand),
			"records.0.locked":    "false",
			"records.0.rr":        "alimail",
			"records.0.type":      "CNAME",
			"records.0.value":     "mail.mxhichin.com",
			"records.0.record_id": acctest.CHECKSET,
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

	var alidnsRecordsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existAlidnsRecordsMapFunc,
		FakeMapFunc:  fakeAlidnsRecordsMapFunc,
	}

	alidnsRecordsCheckInfo.DataSourceTestCheck(t, rand, typeConf, rrRegexConf, valueRegexConf, lineConf, statusConf, idsConf, allConf)
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
