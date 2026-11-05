package kms_test

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

func TestAccAliCloudKmsSecretEphemeral_basic(t *testing.T) {
	rand := sdkacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKmsSecretEphemeral%d", rand)
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
				Config: hclKmsSecretEphemeralSecretConfig(name, "tf-testAccEphemeralSecretValue", "v1"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_kms_secret.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_kms_secret.test", tfjsonpath.New("version_id"), knownvalue.StringExact("v1")),
				},
			},
			{
				Config: hclKmsSecretEphemeralConfigBasic(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test", dataPath.AtMapKey("secret_name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("echo.test", dataPath.AtMapKey("secret_data"), knownvalue.StringExact("tf-testAccEphemeralSecretValue")),
					statecheck.ExpectKnownValue("echo.test", dataPath.AtMapKey("version_id"), knownvalue.StringExact("v1")),
					statecheck.ExpectKnownValue("echo.test", dataPath.AtMapKey("version_stages"), knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("ACSCurrent")})),
					statecheck.ExpectKnownValue("echo.test", dataPath.AtMapKey("secret_data_type"), knownvalue.StringExact("text")),
					statecheck.ExpectKnownValue("echo.test", dataPath.AtMapKey("secret_type"), knownvalue.StringExact("Generic")),
				},
			},
		},
	})
}

func TestAccAliCloudKmsSecretEphemeral_withWriteOnly(t *testing.T) {
	name := sdkacctest.RandomWithPrefix("tf-testAccKmsSecretEphemeralWo")
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
				Config: hclKmsSecretEphemeralWriteOnlySecretConfig(name, "tf-testAccEphemeralWoValue", "v1"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_kms_secret.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_kms_secret.test", tfjsonpath.New("version_id"), knownvalue.StringExact("v1")),
				},
			},
			{
				Config: hclKmsSecretEphemeralWriteOnlyConfig(name, "tf-testAccEphemeralWoValue", "v1", "1"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_kms_ciphertext.test", tfjsonpath.New("ciphertext_blob"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_kms_ciphertext.test", tfjsonpath.New("plaintext_wo_version"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("alicloud_kms_ciphertext.test", tfjsonpath.New("plaintext"), knownvalue.Null()),
					statecheck.ExpectKnownValue("data.alicloud_kms_plaintext.verify", tfjsonpath.New("plaintext"), knownvalue.StringExact("tf-testAccEphemeralWoValue")),
				},
			},
			{
				Config: hclKmsSecretEphemeralWriteOnlyConfig(name, "tf-testAccEphemeralWoValue-rotated", "v2", "1"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_kms_ciphertext.test", tfjsonpath.New("plaintext_wo_version"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("data.alicloud_kms_plaintext.verify", tfjsonpath.New("plaintext"), knownvalue.StringExact("tf-testAccEphemeralWoValue")),
				},
			},
			{
				Config: hclKmsSecretEphemeralWriteOnlyConfig(name, "tf-testAccEphemeralWoValue-rotated", "v2", "2"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_kms_ciphertext.test", tfjsonpath.New("ciphertext_blob"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_kms_ciphertext.test", tfjsonpath.New("plaintext_wo_version"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue("alicloud_kms_ciphertext.test", tfjsonpath.New("plaintext"), knownvalue.Null()),
					statecheck.ExpectKnownValue("data.alicloud_kms_plaintext.verify", tfjsonpath.New("plaintext"), knownvalue.StringExact("tf-testAccEphemeralWoValue-rotated")),
				},
			},
		},
	})
}

func TestAccAliCloudKmsSecretEphemeral_versionSelection(t *testing.T) {
	name := sdkacctest.RandomWithPrefix("tf-testAccKmsSecretEphemeralVer")
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
				Config: hclKmsSecretEphemeralSecretConfig(name, "tf-testAccEphemeralVerV1", "v1"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_kms_secret.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_kms_secret.test", tfjsonpath.New("version_id"), knownvalue.StringExact("v1")),
				},
			},
			{
				Config: hclKmsSecretEphemeralSecretConfig(name, "tf-testAccEphemeralVerV2", "v2"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_kms_secret.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_kms_secret.test", tfjsonpath.New("version_id"), knownvalue.StringExact("v2")),
				},
			},
			{
				Config: hclKmsSecretEphemeralVersionOpenConfig(name, "openv1", `version_id = "v1"`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.openv1", dataPath.AtMapKey("secret_data"), knownvalue.StringExact("tf-testAccEphemeralVerV1")),
					statecheck.ExpectKnownValue("echo.openv1", dataPath.AtMapKey("version_id"), knownvalue.StringExact("v1")),
				},
			},
			{
				Config: hclKmsSecretEphemeralVersionOpenConfig(name, "openv2", `version_id = "v2"`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.openv2", dataPath.AtMapKey("secret_data"), knownvalue.StringExact("tf-testAccEphemeralVerV2")),
					statecheck.ExpectKnownValue("echo.openv2", dataPath.AtMapKey("version_id"), knownvalue.StringExact("v2")),
				},
			},
			{
				Config: hclKmsSecretEphemeralVersionOpenConfig(name, "openprev", `version_stage = "ACSPrevious"`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.openprev", dataPath.AtMapKey("secret_data"), knownvalue.StringExact("tf-testAccEphemeralVerV1")),
					statecheck.ExpectKnownValue("echo.openprev", dataPath.AtMapKey("version_id"), knownvalue.StringExact("v1")),
					statecheck.ExpectKnownValue("echo.openprev", dataPath.AtMapKey("version_stages"), knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("ACSPrevious")})),
				},
			},
		},
	})
}

func TestAccAliCloudKmsSecretEphemeral_arnAddressing(t *testing.T) {
	rand := sdkacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKmsSecretEphemeralArn%d", rand)
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
				Config: hclKmsSecretEphemeralArnConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_kms_secret.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.test", dataPath.AtMapKey("secret_name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.test", dataPath.AtMapKey("secret_data"), knownvalue.StringExact("tf-testAccEphemeralArnValue")),
					statecheck.ExpectKnownValue("echo.test", dataPath.AtMapKey("version_id"), knownvalue.StringExact("v1")),
				},
			},
		},
	})
}

func hclKmsSecretEphemeralSecretConfig(name, secretData, versionId string) string {
	return fmt.Sprintf(`
resource "alicloud_kms_secret" "test" {
  secret_name                   = "%s"
  secret_data                   = "%s"
  version_id                    = "%s"
  force_delete_without_recovery = true
}
`, name, secretData, versionId)
}

func hclKmsSecretEphemeralKeyConfig(name string) string {
	return fmt.Sprintf(`
resource "alicloud_kms_key" "test" {
  description            = "%s"
  pending_window_in_days = 7
}
`, name)
}

func hclKmsSecretEphemeralConfigBasic(name string) string {
	template := hclKmsSecretEphemeralSecretConfig(name, "tf-testAccEphemeralSecretValue", "v1")
	return fmt.Sprintf(`
%s

ephemeral "alicloud_kms_secret" "test" {
  secret_name = alicloud_kms_secret.test.secret_name
}

provider "echo" {
  data = ephemeral.alicloud_kms_secret.test
}

resource "echo" "test" {}
`, template)
}

func hclKmsSecretEphemeralWriteOnlySecretConfig(name, secretData, versionId string) string {
	keyConfig := hclKmsSecretEphemeralKeyConfig(name)
	secretConfig := hclKmsSecretEphemeralSecretConfig(name, secretData, versionId)
	return fmt.Sprintf(`
%s
%s`, keyConfig, secretConfig)
}

func hclKmsSecretEphemeralVersionOpenConfig(name, echoLabel, ephemeralAttrs string) string {
	template := hclKmsSecretEphemeralSecretConfig(name, "tf-testAccEphemeralVerV2", "v2")
	return fmt.Sprintf(`
%s

ephemeral "alicloud_kms_secret" "test" {
  secret_name = alicloud_kms_secret.test.secret_name
  %s
}

provider "echo" {
  data = ephemeral.alicloud_kms_secret.test
}

resource "echo" "%s" {}
`, template, ephemeralAttrs, echoLabel)
}

func hclKmsSecretEphemeralWriteOnlyConfig(name, secretData, versionId, woVersion string) string {
	keyConfig := hclKmsSecretEphemeralKeyConfig(name)
	secretConfig := hclKmsSecretEphemeralSecretConfig(name, secretData, versionId)
	return fmt.Sprintf(`
%s

%s

ephemeral "alicloud_kms_secret" "test" {
  secret_name = alicloud_kms_secret.test.secret_name
}

resource "alicloud_kms_ciphertext" "test" {
  key_id               = alicloud_kms_key.test.id
  plaintext_wo         = ephemeral.alicloud_kms_secret.test.secret_data
  plaintext_wo_version = %s
}

data "alicloud_kms_plaintext" "verify" {
  ciphertext_blob = alicloud_kms_ciphertext.test.ciphertext_blob
}
`, keyConfig, secretConfig, woVersion)
}

func hclKmsSecretEphemeralArnConfig(name string) string {
	template := hclKmsSecretEphemeralSecretConfig(name, "tf-testAccEphemeralArnValue", "v1")
	return fmt.Sprintf(`
%s

ephemeral "alicloud_kms_secret" "test" {
  secret_name = alicloud_kms_secret.test.arn
}

provider "echo" {
  data = ephemeral.alicloud_kms_secret.test
}

resource "echo" "test" {}
`, template)
}
