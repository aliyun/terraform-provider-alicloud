package alicloud

import (
	"fmt"
	_ "github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAmqpStaticAccountDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAmqpStaticAccountSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_amqp_static_account.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAmqpStaticAccountSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_amqp_static_account.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAmqpStaticAccountSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_static_account.default.id}"]`,
			"instance_id": `"${data.alicloud_amqp_instances.default.ids.0}"`,
		}),
		fakeConfig: testAccCheckAlicloudAmqpStaticAccountSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_static_account.default.id}_fake"]`,
			"instance_id": `"${data.alicloud_amqp_instances.default.ids.0}_fake"`,
		}),
	}

	AmqpStaticAccountCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

var existAmqpStaticAccountMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"accounts.#":             "1",
		"accounts.0.id":          CHECKSET,
		"accounts.0.access_key":  CHECKSET,
		"accounts.0.create_time": CHECKSET,
		"accounts.0.instance_id": CHECKSET,
		"accounts.0.master_uid":  CHECKSET,
		"accounts.0.password":    CHECKSET,
		"accounts.0.user_name":   CHECKSET,
	}
}

var fakeAmqpStaticAccountMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"accounts.#": "0",
	}
}

var AmqpStaticAccountCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_amqp_static_accounts.default",
	existMapFunc: existAmqpStaticAccountMapFunc,
	fakeMapFunc:  fakeAmqpStaticAccountMapFunc,
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
