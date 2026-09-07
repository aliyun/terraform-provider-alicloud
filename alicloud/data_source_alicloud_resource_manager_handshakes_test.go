package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudResourceManagerHandshakesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	// InvitedAccountId is required when creating test dependent resources. If not set, the test will be skipped.
	invitedAccountId := os.Getenv("INVITED_ALICLOUD_ACCOUNT_ID")
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerHandshakesSourceConfig(rand, invitedAccountId, map[string]string{
			"ids": `["${alicloud_resource_manager_handshake.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerHandshakesSourceConfig(rand, invitedAccountId, map[string]string{
			"ids": `["${alicloud_resource_manager_handshake.default.id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerHandshakesSourceConfig(rand, invitedAccountId, map[string]string{
			"status": `"Accepted"`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerHandshakesSourceConfig(rand, invitedAccountId, map[string]string{
			"status": `"Pending"`,
		}),
	}

	var existResourceManagerHandshakesRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"handshakes.#":                           "1",
			"ids.#":                                  "1",
			"handshakes.0.id":                        CHECKSET,
			"handshakes.0.handshake_id":              CHECKSET,
			"handshakes.0.expire_time":               CHECKSET,
			"handshakes.0.master_account_id":         CHECKSET,
			"handshakes.0.master_account_name":       CHECKSET,
			"handshakes.0.modify_time":               CHECKSET,
			"handshakes.0.note":                      fmt.Sprintf("tftest_%d", rand),
			"handshakes.0.resource_directory_id":     CHECKSET,
			"handshakes.0.status":                    CHECKSET,
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

	var handshakesRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_resource_manager_handshakes.default",
		existMapFunc: existResourceManagerHandshakesRecordsMapFunc,
		fakeMapFunc:  fakeResourceManagerHandshakesRecordsMapFunc,
	}

	var preCheck = func() {
		testAccPreCheckEnterpriseAccountEnabled(t)
		testAccPreCheck(t)
		testAccPreCheckWithResourceManagerHandshakesSetting(t)
	}

	handshakesRecordsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf)

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
