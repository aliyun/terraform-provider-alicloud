package cr_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccAliCloudCrEeAuthorizationTokenEphemeral_basic(t *testing.T) {
	rand := sdkacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-creetoken-%d", rand)
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
				Config: hclCrEeAuthorizationTokenEphemeralInstanceConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cr_ee_instance.default", tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
			{
				Config: hclCrEeAuthorizationTokenEphemeralOpenConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.open", dataPath.AtMapKey("instance_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open", dataPath.AtMapKey("expires_in_hours"), knownvalue.Null()),
					statecheck.ExpectKnownValue("echo.open", dataPath.AtMapKey("temp_username"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open", dataPath.AtMapKey("authorization_token"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open", dataPath.AtMapKey("expire_time"), knownvalue.NotNull()),
				},
			},
			{
				Config: hclCrEeAuthorizationTokenEphemeralExpiresConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.open_expires", dataPath.AtMapKey("instance_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open_expires", dataPath.AtMapKey("expires_in_hours"), knownvalue.Int64Exact(24)),
					statecheck.ExpectKnownValue("echo.open_expires", dataPath.AtMapKey("temp_username"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open_expires", dataPath.AtMapKey("authorization_token"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open_expires", dataPath.AtMapKey("expire_time"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func hclCrEeAuthorizationTokenEphemeralInstanceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_cr_ee_instance" "default" {
  payment_type   = "Subscription"
  period         = 1
  renewal_status = "ManualRenewal"
  instance_type  = "Basic"
  instance_name  = var.name
}
`, name)
}

func hclCrEeAuthorizationTokenEphemeralOpenConfig(name string) string {
	return fmt.Sprintf(`
%[1]s

ephemeral "alicloud_cr_ee_authorization_token" "default" {
  instance_id = alicloud_cr_ee_instance.default.id
}

provider "echo" {
  data = ephemeral.alicloud_cr_ee_authorization_token.default
}

resource "echo" "open" {}
`, hclCrEeAuthorizationTokenEphemeralInstanceConfig(name))
}

func hclCrEeAuthorizationTokenEphemeralExpiresConfig(name string) string {
	return fmt.Sprintf(`
%[1]s

ephemeral "alicloud_cr_ee_authorization_token" "default" {
  instance_id      = alicloud_cr_ee_instance.default.id
  expires_in_hours = 24
}

provider "echo" {
  data = ephemeral.alicloud_cr_ee_authorization_token.default
}

resource "echo" "open_expires" {}
`, hclCrEeAuthorizationTokenEphemeralInstanceConfig(name))
}
