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

func TestAccAliCloudKmsPlaintextEphemeral_basic(t *testing.T) {
	rand := sdkacctest.RandIntRange(1000000, 9999999)
	value := fmt.Sprintf("tf-testAccEphemeralPlaintext%d", rand)
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
				Config: hclKmsPlaintextEphemeralCiphertextConfig(value),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_kms_ciphertext.test", tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
			{
				Config: hclKmsPlaintextEphemeralBasicConfig(value),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test", dataPath.AtMapKey("plaintext"), knownvalue.StringExact(value)),
					statecheck.ExpectKnownValue("echo.test", dataPath.AtMapKey("key_id"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccAliCloudKmsPlaintextEphemeral_withEncryptionContext(t *testing.T) {
	rand := sdkacctest.RandIntRange(1000000, 9999999)
	value := fmt.Sprintf("tf-testAccEphemeralPlaintextCtx%d", rand)
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
				Config: hclKmsPlaintextEphemeralEncryptionContextConfig(value, "test"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test", dataPath.AtMapKey("plaintext"), knownvalue.StringExact(value)),
					statecheck.ExpectKnownValue("echo.test", dataPath.AtMapKey("key_id"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccAliCloudKmsPlaintextEphemeral_withWriteOnly(t *testing.T) {
	rand := sdkacctest.RandIntRange(1000000, 9999999)
	value := fmt.Sprintf("tf-testAccEphemeralPlaintextWo%d", rand)
	dataPath := tfjsonpath.New("data")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
		CheckDestroy:             acctest.CheckDestroyNoop,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		Steps: []resource.TestStep{
			{
				Config: hclKmsPlaintextEphemeralCiphertextConfig(value),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_kms_ciphertext.test", tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
			{
				Config: hclKmsPlaintextEphemeralWriteOnlyConfig(value, "1", "verify2"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_kms_ciphertext.sink", tfjsonpath.New("plaintext_wo_version"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("alicloud_kms_ciphertext.sink", tfjsonpath.New("plaintext"), knownvalue.Null()),
					statecheck.ExpectKnownValue("echo.verify2", dataPath.AtMapKey("plaintext"), knownvalue.StringExact(value)),
				},
			},
			{
				Config: hclKmsPlaintextEphemeralWriteOnlyConfig(value, "2", "verify3"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_kms_ciphertext.sink", tfjsonpath.New("plaintext_wo_version"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue("echo.verify3", dataPath.AtMapKey("plaintext"), knownvalue.StringExact(value)),
				},
			},
		},
	})
}

func hclKmsPlaintextEphemeralKeyConfig() string {
	return `
resource "alicloud_kms_key" "default" {
  description            = "tf-testAccKmsPlaintextEphemeral"
  is_enabled             = true
  pending_window_in_days = 7
}
`
}

func hclKmsPlaintextEphemeralCiphertextConfig(value string) string {
	return fmt.Sprintf(`
%s

resource "alicloud_kms_ciphertext" "test" {
  key_id    = alicloud_kms_key.default.id
  plaintext = "%s"
}
`, hclKmsPlaintextEphemeralKeyConfig(), value)
}

func hclKmsPlaintextEphemeralBasicConfig(value string) string {
	return fmt.Sprintf(`
%s

ephemeral "alicloud_kms_plaintext" "test" {
  ciphertext_blob = alicloud_kms_ciphertext.test.ciphertext_blob
}

provider "echo" {
  data = ephemeral.alicloud_kms_plaintext.test
}

resource "echo" "test" {}
`, hclKmsPlaintextEphemeralCiphertextConfig(value))
}

func hclKmsPlaintextEphemeralEncryptionContextConfig(value, stage string) string {
	return fmt.Sprintf(`
%s

resource "alicloud_kms_ciphertext" "test" {
  key_id             = alicloud_kms_key.default.id
  plaintext          = "%s"
  encryption_context = {
    stage = "%s"
  }
}

ephemeral "alicloud_kms_plaintext" "test" {
  ciphertext_blob    = alicloud_kms_ciphertext.test.ciphertext_blob
  encryption_context = {
    stage = "%s"
  }
}

provider "echo" {
  data = ephemeral.alicloud_kms_plaintext.test
}

resource "echo" "test" {}
`, hclKmsPlaintextEphemeralKeyConfig(), value, stage, stage)
}

func hclKmsPlaintextEphemeralWriteOnlyConfig(value, woVersion, echoLabel string) string {
	return fmt.Sprintf(`
%s

ephemeral "alicloud_kms_plaintext" "source" {
  ciphertext_blob = alicloud_kms_ciphertext.test.ciphertext_blob
}

resource "alicloud_kms_ciphertext" "sink" {
  key_id              = alicloud_kms_key.default.id
  plaintext_wo        = ephemeral.alicloud_kms_plaintext.source.plaintext
  plaintext_wo_version = %s
}

ephemeral "alicloud_kms_plaintext" "verify" {
  ciphertext_blob = alicloud_kms_ciphertext.sink.ciphertext_blob
}

provider "echo" {
  data = ephemeral.alicloud_kms_plaintext.verify
}

resource "echo" "%s" {}
`, hclKmsPlaintextEphemeralCiphertextConfig(value), woVersion, echoLabel)
}
