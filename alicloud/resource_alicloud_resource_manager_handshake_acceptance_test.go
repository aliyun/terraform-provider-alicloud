package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var ResourceManagerHandshakeAcceptanceMap = map[string]string{
	"status": CHECKSET,
}

// Accepting a handshake must be performed by the invited account, so this test needs two accounts.
// The management (master) account is the default provider (ALICLOUD_ACCESS_KEY/SECRET_KEY) and sends
// the invitation; the invited account is supplied via ALICLOUD_ACCESS_KEY_2/SECRET_KEY_2 (uid in
// ALICLOUD_ACCOUNT_ID_2) and accepts it. handshake_id is ForceNew, so no second/update step is needed.
func ResourceManagerHandshakeAcceptanceBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

# Invited account (accepts the invitation), explicit credentials so it overrides the default keys.
provider "alicloud" {
  alias      = "invited"
  access_key = "%s"
  secret_key = "%s"
}

# Management account (default provider) sends the invitation.
resource "alicloud_resource_manager_handshake" "default" {
  target_entity = "%s"
  target_type   = "Account"
  note          = "tf-test-acceptance"
}
`, name, os.Getenv("ALICLOUD_ACCESS_KEY_2"), os.Getenv("ALICLOUD_SECRET_KEY_2"), os.Getenv("ALICLOUD_ACCOUNT_ID_2"))
}

func testAccPreCheckHandshakeAcceptanceProfiles(t *testing.T) {
	if os.Getenv("ALICLOUD_ACCESS_KEY_2") == "" || os.Getenv("ALICLOUD_SECRET_KEY_2") == "" ||
		os.Getenv("ALICLOUD_ACCOUNT_ID_2") == "" {
		t.Skipf("Skipping: set ALICLOUD_ACCESS_KEY_2/SECRET_KEY_2/ALICLOUD_ACCOUNT_ID_2 for the invited account")
	}
}

func TestAccAliCloudResourceManagerHandshakeAcceptance_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_handshake_acceptance.default"
	ra := resourceAttrInit(resourceId, ResourceManagerHandshakeAcceptanceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourceManagerServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerHandshake")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccResourceManagerHandshakeAcceptance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ResourceManagerHandshakeAcceptanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckEnterpriseAccountEnabled(t)
			testAccPreCheckHandshakeAcceptanceProfiles(t)
		},

		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"handshake_id": "${alicloud_resource_manager_handshake.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"handshake_id": CHECKSET,
						"status":       "Accepted",
					}),
				),
			},
		},
	})
}
