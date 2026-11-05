package ram_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccAliCloudRamAccessKeyEphemeral_basic(t *testing.T) {
	rand := sdkacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-ramak-%d", rand)
	dataPath := tfjsonpath.New("data")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		CheckDestroy: acctest.CheckDestroyNoop,
		Steps: []resource.TestStep{
			{
				// Create the RAM user first. Ephemeral resources open during the
				// plan phase, so a managed resource declared in the same step is
				// not yet applied when Open runs and calls CreateAccessKey; the
				// user must already exist in the cloud before the next step.
				Config: hclRamAccessKeyEphemeralUserConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_ram_user.default", tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
			{
				Config: hclRamAccessKeyEphemeralOpenConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamAccessKeyRevoked(name),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.open", dataPath.AtMapKey("user_name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open", dataPath.AtMapKey("access_key_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open", dataPath.AtMapKey("access_key_secret"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open", dataPath.AtMapKey("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open", dataPath.AtMapKey("create_date"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func hclRamAccessKeyEphemeralUserConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_ram_user" "default" {
  name  = var.name
  force = true
}
`, name)
}

func hclRamAccessKeyEphemeralOpenConfig(name string) string {
	return fmt.Sprintf(`
%[1]s

ephemeral "alicloud_ram_access_key" "default" {
  user_name = alicloud_ram_user.default.name
}

provider "echo" {
  data = ephemeral.alicloud_ram_access_key.default
}

resource "echo" "open" {}
`, hclRamAccessKeyEphemeralUserConfig(name))
}

// testAccCheckRamAccessKeyRevoked verifies that Close revoked the access key
// issued by Open. The check runs after the apply, once Terraform has closed the
// ephemeral resource, and before the framework destroys the RAM user: the
// fixture's force = true would delete any leftover key during destroy and mask
// a missing or ineffective Close.
func testAccCheckRamAccessKeyRevoked(userName string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		response, err := client.RpcPost("Ram", "2015-05-01", "ListAccessKeys",
			map[string]interface{}{}, map[string]interface{}{"UserName": userName}, true)
		if err != nil {
			return fmt.Errorf("listing access keys of RAM user %s: %s", userName, err)
		}

		remaining := 0
		if raw, ok := response["AccessKeys"]; ok && raw != nil {
			if accessKeys, ok := raw.(map[string]interface{}); ok {
				if keys, ok := accessKeys["AccessKey"].([]interface{}); ok {
					remaining = len(keys)
				} else if accessKeys["AccessKey"] != nil {
					remaining = 1
				}
			}
		}
		if remaining != 0 {
			return fmt.Errorf("expected Close to revoke the access key of RAM user %s, but %d key(s) remain", userName, remaining)
		}
		return nil
	}
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
