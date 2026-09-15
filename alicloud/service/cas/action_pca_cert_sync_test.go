package cas_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TestAccAliCloudSslCertificatesServicePcaCertSyncAction exercises the framework Action via
// HCL: the after_create action_trigger of the second client certificate invokes
// alicloud_ssl_certificates_service_pca_cert_sync with both certificate identifiers, and the
// check polls the API until the server-side UploadFlag flips. Requires Terraform >= 1.14.
func TestAccAliCloudSslCertificatesServicePcaCertSyncAction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudSslCertificatesServicePcaCertSyncActionConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckPcaCertsUploaded("alicloud_ssl_certificates_service_pca_cert.one"),
					testCheckPcaCertsUploaded("alicloud_ssl_certificates_service_pca_cert.two"),
				),
			},
		},
	})
}

func testAccAliCloudSslCertificatesServicePcaCertSyncActionConfig() string {
	return `
resource "alicloud_ssl_certificates_service_pca_certificate" "root" {
  organization      = "a"
  years             = "1"
  locality          = "a"
  organization_unit = "a"
  state             = "a"
  common_name       = "cbc.certqa.cn"
}

resource "alicloud_ssl_certificates_service_pca_certificate" "sub" {
  parent_identifier = alicloud_ssl_certificates_service_pca_certificate.root.id
  organization      = "a"
  years             = "1"
  locality          = "a"
  organization_unit = "a"
  state             = "a"
  common_name       = "cbc.certqa.cn"
  algorithm         = "RSA_2048"
  certificate_type  = "SUB_ROOT"
  enable_crl        = true
}

resource "alicloud_ssl_certificates_service_pca_cert" "one" {
  days              = "1"
  parent_identifier = alicloud_ssl_certificates_service_pca_certificate.sub.id
  algorithm         = "RSA_2048"
  common_name       = "tfaccpcasync-one.example.com"
  organization      = "terraform"
  state             = "Beijing"
  country_code      = "cn"
  status            = "REVOKE"
}

resource "alicloud_ssl_certificates_service_pca_cert" "two" {
  days              = "1"
  parent_identifier = alicloud_ssl_certificates_service_pca_certificate.sub.id
  algorithm         = "RSA_2048"
  common_name       = "tfaccpcasync-two.example.com"
  organization      = "terraform"
  state             = "Beijing"
  country_code      = "cn"
  status            = "REVOKE"

  depends_on = [alicloud_ssl_certificates_service_pca_cert.one]

  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.alicloud_ssl_certificates_service_pca_cert_sync.sync]
    }
  }
}

action "alicloud_ssl_certificates_service_pca_cert_sync" "sync" {
  config {
    ids = [
      alicloud_ssl_certificates_service_pca_cert.one.id,
      alicloud_ssl_certificates_service_pca_cert.two.id,
    ]
  }
}
`
}

// testCheckPcaCertsUploaded polls the certificate until the server-side UploadFlag is 1,
// proving the action reached the service. The processing is asynchronous, so the poll
// waits up to three minutes.
func testCheckPcaCertsUploaded(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource %s not found in state", name)
		}

		id := rs.Primary.Attributes["id"]
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)

		return retry.RetryContext(context.Background(), 3*time.Minute, func() *retry.RetryError {
			response, err := client.RpcPost("cas", "2020-06-30", "DescribeClientCertificate", map[string]interface{}{}, map[string]interface{}{
				"Identifier": id,
			}, true)
			if err != nil {
				return retry.NonRetryableError(fmt.Errorf("describing client certificate %s: %w", id, err))
			}

			certObject, ok := response["Certificate"].(map[string]interface{})
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("describing client certificate %s: no Certificate object in response: %v", id, response))
			}

			if uploadFlag, ok := certObject["UploadFlag"].(json.Number); ok && uploadFlag.String() == "1" {
				return nil
			}

			return retry.RetryableError(fmt.Errorf("certificate %s not yet uploaded (UploadFlag=%v)", id, certObject["UploadFlag"]))
		})
	}
}
