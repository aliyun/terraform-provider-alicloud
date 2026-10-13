package oos_test

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

func TestAccAliCloudOosSecretParameterEphemeral_basic(t *testing.T) {
	rand := sdkacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccOosSecretParameterEphemeral%d", rand)
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
				Config: hclOosSecretParameterEphemeralSecretConfig(name, "tf-testAccEphemeralParameterValue"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_oos_secret_parameter.test", tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
			{
				Config: hclOosSecretParameterEphemeralBasicConfig(name, "tf-testAccEphemeralParameterValue"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test", dataPath.AtMapKey("secret_parameter_name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("echo.test", dataPath.AtMapKey("value"), knownvalue.StringExact("tf-testAccEphemeralParameterValue")),
					statecheck.ExpectKnownValue("echo.test", dataPath.AtMapKey("parameter_version"), knownvalue.Int64Exact(1)),
				},
			},
		},
	})
}

func TestAccAliCloudOosSecretParameterEphemeral_withWriteOnly(t *testing.T) {
	rand := sdkacctest.RandIntRange(1000000, 9999999)
	secretName := fmt.Sprintf("tf-testAccOosSecretParameterEphemeralWo%d", rand)
	sinkName := fmt.Sprintf("tf-testacc-oos-param-sink-%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             acctest.CheckDestroyNoop,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		Steps: []resource.TestStep{
			{
				Config: hclOosSecretParameterEphemeralSecretConfig(secretName, "tf-testAccEphemeralWoValue"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_oos_secret_parameter.test", tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
			{
				Config: hclOosSecretParameterEphemeralWriteOnlyConfig(secretName, sinkName, "tf-testAccEphemeralWoValue", "1"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_oos_parameter.sink", tfjsonpath.New("value_wo_version"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("alicloud_oos_parameter.sink", tfjsonpath.New("value"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("data.alicloud_oos_parameters.verify", tfjsonpath.New("parameters").AtSliceIndex(0).AtMapKey("value"), knownvalue.StringExact("tf-testAccEphemeralWoValue")),
				},
			},
			{
				Config: hclOosSecretParameterEphemeralWriteOnlyConfig(secretName, sinkName, "tf-testAccEphemeralWoValue-rotated", "1"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_oos_parameter.sink", tfjsonpath.New("value_wo_version"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("data.alicloud_oos_parameters.verify", tfjsonpath.New("parameters").AtSliceIndex(0).AtMapKey("value"), knownvalue.StringExact("tf-testAccEphemeralWoValue")),
				},
			},
			{
				Config: hclOosSecretParameterEphemeralWriteOnlyConfig(secretName, sinkName, "tf-testAccEphemeralWoValue-rotated", "2"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_oos_parameter.sink", tfjsonpath.New("value_wo_version"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue("alicloud_oos_parameter.sink", tfjsonpath.New("value"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("data.alicloud_oos_parameters.verify", tfjsonpath.New("parameters").AtSliceIndex(0).AtMapKey("value"), knownvalue.StringExact("tf-testAccEphemeralWoValue-rotated")),
				},
			},
		},
	})
}

func TestAccAliCloudOosSecretParameterEphemeral_versionSelection(t *testing.T) {
	rand := sdkacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccOosSecretParameterEphemeralVer%d", rand)
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
				Config: hclOosSecretParameterEphemeralSecretConfig(name, "tf-testAccEphemeralVerV1"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_oos_secret_parameter.test", tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
			{
				Config: hclOosSecretParameterEphemeralSecretConfig(name, "tf-testAccEphemeralVerV2"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_oos_secret_parameter.test", tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
			{
				Config: hclOosSecretParameterEphemeralVersionOpenConfig(name, "pinned", "parameter_version = 1"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.pinned", dataPath.AtMapKey("value"), knownvalue.StringExact("tf-testAccEphemeralVerV1")),
					statecheck.ExpectKnownValue("echo.pinned", dataPath.AtMapKey("parameter_version"), knownvalue.Int64Exact(1)),
				},
			},
			{
				Config: hclOosSecretParameterEphemeralVersionOpenConfig(name, "latest", ""),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.latest", dataPath.AtMapKey("value"), knownvalue.StringExact("tf-testAccEphemeralVerV2")),
					statecheck.ExpectKnownValue("echo.latest", dataPath.AtMapKey("parameter_version"), knownvalue.Int64Exact(2)),
				},
			},
		},
	})
}

func hclOosSecretParameterEphemeralSecretConfig(name, value string) string {
	return fmt.Sprintf(`
resource "alicloud_oos_secret_parameter" "test" {
  secret_parameter_name = "%s"
  value                 = "%s"
}
`, name, value)
}

func hclOosSecretParameterEphemeralBasicConfig(name, value string) string {
	template := hclOosSecretParameterEphemeralSecretConfig(name, value)
	return fmt.Sprintf(`
%s

ephemeral "alicloud_oos_secret_parameter" "test" {
  secret_parameter_name = alicloud_oos_secret_parameter.test.secret_parameter_name
}

provider "echo" {
  data = ephemeral.alicloud_oos_secret_parameter.test
}

resource "echo" "test" {}
`, template)
}

func hclOosSecretParameterEphemeralVersionOpenConfig(name, echoLabel, ephemeralAttrs string) string {
	template := hclOosSecretParameterEphemeralSecretConfig(name, "tf-testAccEphemeralVerV2")
	return fmt.Sprintf(`
%s

ephemeral "alicloud_oos_secret_parameter" "test" {
  secret_parameter_name = alicloud_oos_secret_parameter.test.secret_parameter_name
  %s
}

provider "echo" {
  data = ephemeral.alicloud_oos_secret_parameter.test
}

resource "echo" "%s" {}
`, template, ephemeralAttrs, echoLabel)
}

func hclOosSecretParameterEphemeralWriteOnlyConfig(secretName, sinkName, secretValue, woVersion string) string {
	secretConfig := hclOosSecretParameterEphemeralSecretConfig(secretName, secretValue)
	return fmt.Sprintf(`
%s

ephemeral "alicloud_oos_secret_parameter" "test" {
  secret_parameter_name = alicloud_oos_secret_parameter.test.secret_parameter_name
}

resource "alicloud_oos_parameter" "sink" {
  parameter_name   = "%s"
  type             = "String"
  value_wo         = ephemeral.alicloud_oos_secret_parameter.test.value
  value_wo_version = %s
}

data "alicloud_oos_parameters" "verify" {
  parameter_name = alicloud_oos_parameter.sink.parameter_name
  enable_details = true
}
`, secretConfig, sinkName, woVersion)
}
