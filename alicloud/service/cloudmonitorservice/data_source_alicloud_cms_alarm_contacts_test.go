package cloudmonitorservice_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCmsAlarmContacts_basic(t *testing.T) {
	rand := tfacctest.RandInt()
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsAlarmContactsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cms_alarm_contact.default.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudCmsAlarmContactsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cms_alarm_contact.default.id}_fake"`,
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsAlarmContactsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cms_alarm_contact.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudCmsAlarmContactsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cms_alarm_contact.default.id}_fake"]`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsAlarmContactsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cms_alarm_contact.default.id}"`,
			"ids":        `["${alicloud_cms_alarm_contact.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudCmsAlarmContactsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cms_alarm_contact.default.id}_fake"`,
			"ids":        `["${alicloud_cms_alarm_contact.default.id}_fake"]`,
		}),
	}

	var existcmsAlarmContactsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"contacts.#":                    "1",
			"contacts.0.id":                 acctest.CHECKSET,
			"contacts.0.alarm_contact_name": acctest.CHECKSET,
			"contacts.0.describe":           "For Test",
		}
	}

	var fakecmsAlarmContactsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var cmsAlarmContactsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cms_alarm_contacts.default",
		ExistMapFunc: existcmsAlarmContactsMapFunc,
		FakeMapFunc:  fakecmsAlarmContactsMapFunc,
	}

	cmsAlarmContactsCheckInfo.DataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlicloudCmsAlarmContactsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
		variable "name" {
			default = "tf-testAccCms-%d"
		}
		resource "alicloud_cms_alarm_contact" "default" {
			alarm_contact_name = var.name
		    describe           = "For Test"
		    channels_mail      = "hello.uuuu@aaa.com"
			lifecycle {
				ignore_changes = [channels_mail]
  			}	
		}

		data "alicloud_cms_alarm_contacts" "default" {
		  %s
		}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
