package sts_test

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

func TestAccAliCloudStsTokenEphemeral_basic(t *testing.T) {
	rand := sdkacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-ststoken-%d", rand)
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
				Config: hclStsTokenEphemeralOpenConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.open", dataPath.AtMapKey("role_arn"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open", dataPath.AtMapKey("session_name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open", dataPath.AtMapKey("duration_seconds"), knownvalue.Null()),
					statecheck.ExpectKnownValue("echo.open", dataPath.AtMapKey("access_key_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open", dataPath.AtMapKey("access_key_secret"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open", dataPath.AtMapKey("security_token"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open", dataPath.AtMapKey("expiration"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open", dataPath.AtMapKey("assumed_role_arn"), knownvalue.NotNull()),
				},
			},
			{
				Config: hclStsTokenEphemeralParamsConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.open_params", dataPath.AtMapKey("duration_seconds"), knownvalue.Int64Exact(900)),
					statecheck.ExpectKnownValue("echo.open_params", dataPath.AtMapKey("access_key_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open_params", dataPath.AtMapKey("access_key_secret"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open_params", dataPath.AtMapKey("security_token"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open_params", dataPath.AtMapKey("expiration"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open_params", dataPath.AtMapKey("assumed_role_arn"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func hclStsTokenEphemeralRoleConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_caller_identity" "current" {}

resource "alicloud_ram_role" "default" {
  role_name = var.name
  assume_role_policy_document = jsonencode({
    Version   = "1"
    Statement = [{
      Effect    = "Allow"
      Action    = "sts:AssumeRole"
      Principal = {
        RAM = [data.alicloud_caller_identity.current.arn]
      }
    }]
  })
}
`, name)
}

func hclStsTokenEphemeralOpenConfig(name string) string {
	return fmt.Sprintf(`
%[1]s

ephemeral "alicloud_sts_token" "default" {
  role_arn     = alicloud_ram_role.default.arn
  session_name = "tf-testacc"
}

provider "echo" {
  data = ephemeral.alicloud_sts_token.default
}

resource "echo" "open" {}
`, hclStsTokenEphemeralRoleConfig(name))
}

func hclStsTokenEphemeralParamsConfig(name string) string {
	return fmt.Sprintf(`
%[1]s

ephemeral "alicloud_sts_token" "default" {
  role_arn         = alicloud_ram_role.default.arn
  session_name     = "tf-testacc"
  duration_seconds = 900
  policy           = jsonencode({
    Version   = "1"
    Statement = [{
      Effect   = "Allow"
      Action   = ["oss:ListBuckets"]
      Resource = ["*"]
    }]
  })
  external_id = "tf-testacc-external-id"
}

provider "echo" {
  data = ephemeral.alicloud_sts_token.default
}

resource "echo" "open_params" {}
`, hclStsTokenEphemeralRoleConfig(name))
}
