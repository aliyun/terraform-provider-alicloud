package directmail_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudDirectMailReceiversesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDirectMailReceiversesDataSourceName(rand, map[string]string{
			"ids": `[alicloud_direct_mail_receivers.default.id]`,
		}),
		FakeConfig: testAccCheckAlicloudDirectMailReceiversesDataSourceName(rand, map[string]string{
			"ids": `["fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDirectMailReceiversesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_direct_mail_receivers.default.receivers_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudDirectMailReceiversesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_direct_mail_receivers.default.receivers_name}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDirectMailReceiversesDataSourceName(rand, map[string]string{
			"ids":    `[alicloud_direct_mail_receivers.default.id]`,
			"status": `"0"`,
		}),
		FakeConfig: testAccCheckAlicloudDirectMailReceiversesDataSourceName(rand, map[string]string{
			"ids":    `[alicloud_direct_mail_receivers.default.id]`,
			"status": `"1"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDirectMailReceiversesDataSourceName(rand, map[string]string{
			"ids":        `[alicloud_direct_mail_receivers.default.id]`,
			"name_regex": `"${alicloud_direct_mail_receivers.default.receivers_name}"`,
			"status":     `0`,
		}),
		FakeConfig: testAccCheckAlicloudDirectMailReceiversesDataSourceName(rand, map[string]string{
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
			"receiverses.0.receivers_id":    acctest.CHECKSET,
			"receiverses.0.id":              acctest.CHECKSET,
			"receiverses.0.create_time":     acctest.CHECKSET,
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
	var alicloudDirectMailReceiversesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_direct_mail_receiverses.default",
		ExistMapFunc: existAlicloudDirectMailReceiversesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudDirectMailReceiversesDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithRegions(t, true, connectivity.DmSupportRegions)
	}
	alicloudDirectMailReceiversesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
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
