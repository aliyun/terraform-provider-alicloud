package resourcemanager_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudResourceManagerHandshakesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	// InvitedAccountId is required when creating test dependent resources. If not set, the test will be skipped.
	invitedAccountId := os.Getenv("INVITED_ALICLOUD_ACCOUNT_ID")
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudResourceManagerHandshakesSourceConfig(rand, invitedAccountId, map[string]string{
			"ids": `["${alicloud_resource_manager_handshake.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudResourceManagerHandshakesSourceConfig(rand, invitedAccountId, map[string]string{
			"ids": `["${alicloud_resource_manager_handshake.default.id}_fake"]`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudResourceManagerHandshakesSourceConfig(rand, invitedAccountId, map[string]string{
			"status": `"Accepted"`,
		}),
		FakeConfig: testAccCheckAlicloudResourceManagerHandshakesSourceConfig(rand, invitedAccountId, map[string]string{
			"status": `"Pending"`,
		}),
	}

	var existResourceManagerHandshakesRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"handshakes.#":                           "1",
			"ids.#":                                  "1",
			"handshakes.0.id":                        acctest.CHECKSET,
			"handshakes.0.handshake_id":              acctest.CHECKSET,
			"handshakes.0.expire_time":               acctest.CHECKSET,
			"handshakes.0.master_account_id":         acctest.CHECKSET,
			"handshakes.0.master_account_name":       acctest.CHECKSET,
			"handshakes.0.modify_time":               acctest.CHECKSET,
			"handshakes.0.note":                      fmt.Sprintf("tftest_%d", rand),
			"handshakes.0.resource_directory_id":     acctest.CHECKSET,
			"handshakes.0.status":                    acctest.CHECKSET,
			"handshakes.0.target_entity":             invitedAccountId,
			"handshakes.0.target_type":               "Account",
			"handshakes.0.invited_account_real_name": "",
			"handshakes.0.master_account_real_name":  "",
		}
	}

	var fakeResourceManagerHandshakesRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"handshakes.#": "0",
			"ids.#":        "0",
		}
	}

	var handshakesRecordsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_resource_manager_handshakes.default",
		ExistMapFunc: existResourceManagerHandshakesRecordsMapFunc,
		FakeMapFunc:  fakeResourceManagerHandshakesRecordsMapFunc,
	}

	var preCheck = func() {
		acctest.PreCheckEnterpriseAccountEnabled(t)
		acctest.PreCheck(t)
		testAccPreCheckWithResourceManagerHandshakesSetting(t)
	}

	handshakesRecordsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf)

}

func testAccCheckAlicloudResourceManagerHandshakesSourceConfig(rand int, invitedAccountId string, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
resource "alicloud_resource_manager_handshake" "default" {
  target_entity = "%s"
  target_type = "Account"
  note = "tftest_%d"
}

data "alicloud_resource_manager_handshakes" "default"{
	enable_details = true
%s
}

`, invitedAccountId, rand, strings.Join(pairs, "\n   "))
	return config
}
