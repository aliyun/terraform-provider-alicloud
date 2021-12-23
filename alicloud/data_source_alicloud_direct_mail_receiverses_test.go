package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDirectMailReceiversesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDirectMailReceiversesDataSourceName(rand, map[string]string{
			"ids": `[alicloud_direct_mail_receivers.default.id]`,
		}),
		fakeConfig: testAccCheckAlicloudDirectMailReceiversesDataSourceName(rand, map[string]string{
			"ids": `["fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDirectMailReceiversesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_direct_mail_receivers.default.receivers_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDirectMailReceiversesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_direct_mail_receivers.default.receivers_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDirectMailReceiversesDataSourceName(rand, map[string]string{
			"ids":    `[alicloud_direct_mail_receivers.default.id]`,
			"status": `"0"`,
		}),
		fakeConfig: testAccCheckAlicloudDirectMailReceiversesDataSourceName(rand, map[string]string{
			"ids":    `[alicloud_direct_mail_receivers.default.id]`,
			"status": `"1"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDirectMailReceiversesDataSourceName(rand, map[string]string{
			"ids":        `[alicloud_direct_mail_receivers.default.id]`,
			"name_regex": `"${alicloud_direct_mail_receivers.default.receivers_name}"`,
			"status":     `0`,
		}),
		fakeConfig: testAccCheckAlicloudDirectMailReceiversesDataSourceName(rand, map[string]string{
			"ids":        `["fake"]`,
			"name_regex": `"${alicloud_direct_mail_receivers.default.receivers_name}_fake"`,
			"status":     `1`,
		}),
	}
	var existAlicloudDirectMailReceiversesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"receiverses.#":                 "1",
			"receiverses.0.receivers_id":    CHECKSET,
			"receiverses.0.id":              CHECKSET,
			"receiverses.0.create_time":     CHECKSET,
			"receiverses.0.description":     fmt.Sprintf("tf-testAcc-%d", rand),
			"receiverses.0.receivers_alias": fmt.Sprintf("%d@onaliyun.com", rand),
			"receiverses.0.receivers_name":  fmt.Sprintf("tf-testAcc-%d", rand),
			"receiverses.0.status":          `0`,
		}
	}
	var fakeAlicloudDirectMailReceiversesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudDirectMailReceiversesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_direct_mail_receiverses.default",
		existMapFunc: existAlicloudDirectMailReceiversesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudDirectMailReceiversesDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithRegions(t, true, connectivity.DmSupportRegions)
	}
	alicloudDirectMailReceiversesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudDirectMailReceiversesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAcc-%d"
}

resource "alicloud_direct_mail_receivers" "default" {
	receivers_name = var.name
	receivers_alias = join("",["%d","@onaliyun.com"])
	description =  var.name
}

data "alicloud_direct_mail_receiverses" "default" {
	%s	
}
`, rand, rand, strings.Join(pairs, " \n "))
	return config
}
