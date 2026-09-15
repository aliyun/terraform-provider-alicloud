package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudSslCertificatesServicePcaCertSync_basic(t *testing.T) {
	resourceId := "alicloud_ssl_certificates_service_pca_cert_sync.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsslcertificatesservicepcacertsync%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(s *terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudSslCertificatesServicePcaCertSyncBasicConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "ids.#", "2"),
					resource.TestCheckResourceAttrSet(resourceId, "ids.0"),
					resource.TestCheckResourceAttrSet(resourceId, "ids.1"),
				),
			},
		},
	})
}

func testAccAliCloudSslCertificatesServicePcaCertSyncBasicConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

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

resource "alicloud_ssl_certificates_service_pca_cert" "default" {
  count             = 2
  days              = "1"
  parent_identifier = alicloud_ssl_certificates_service_pca_certificate.sub.id
  algorithm         = "RSA_2048"
  common_name       = "${var.name}-${count.index}"
  organization      = "terraform"
  state             = "Beijing"
  country_code      = "cn"
  upload_flag       = 1
  status            = "REVOKE"
}

resource "alicloud_ssl_certificates_service_pca_cert_sync" "default" {
  ids = alicloud_ssl_certificates_service_pca_cert.default[*].id
}
`, name)
}
