package amqpopen_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	_ "github.com/stretchr/testify/assert"
)

func TestAccAlicloudAmqpStaticAccountDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAmqpStaticAccountSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_amqp_static_account.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudAmqpStaticAccountSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_amqp_static_account.default.id}_fake"]`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAmqpStaticAccountSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_static_account.default.id}"]`,
			"instance_id": `"${data.alicloud_amqp_instances.default.ids.0}"`,
		}),
		FakeConfig: testAccCheckAlicloudAmqpStaticAccountSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_static_account.default.id}_fake"]`,
			"instance_id": `"${data.alicloud_amqp_instances.default.ids.0}_fake"`,
		}),
	}

	AmqpStaticAccountCheckInfo.DataSourceTestCheck(t, rand, idsConf, allConf)
}

var existAmqpStaticAccountMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"accounts.#":             "1",
		"accounts.0.id":          acctest.CHECKSET,
		"accounts.0.access_key":  acctest.CHECKSET,
		"accounts.0.create_time": acctest.CHECKSET,
		"accounts.0.instance_id": acctest.CHECKSET,
		"accounts.0.master_uid":  acctest.CHECKSET,
		"accounts.0.password":    acctest.CHECKSET,
		"accounts.0.user_name":   acctest.CHECKSET,
	}
}

var fakeAmqpStaticAccountMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"accounts.#": "0",
	}
}

var AmqpStaticAccountCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_amqp_static_accounts.default",
	ExistMapFunc: existAmqpStaticAccountMapFunc,
	FakeMapFunc:  fakeAmqpStaticAccountMapFunc,
}

func testAccCheckAlicloudAmqpStaticAccountSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccAmqpStaticAccount%d"
}
data "alicloud_amqp_instances" "default" {
	status = "SERVING"
}

resource "alicloud_amqp_static_account" "default" {
  instance_id = "${data.alicloud_amqp_instances.default.ids.0}"
  access_key  = "%s"
  secret_key  = "%s"
}

data "alicloud_amqp_static_accounts" "default" {
%s
}
`, rand, os.Getenv("ALICLOUD_ACCESS_KEY"), os.Getenv("ALICLOUD_SECRET_KEY"), strings.Join(pairs, "\n   "))
	return config
}
