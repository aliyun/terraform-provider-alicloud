package dataworks_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudDataWorksCertificate_basic0 covers create, read, import and delete. The
// certificate file is a URL: DataWorks downloads the object itself, so the test uploads a
// self-signed certificate to OSS and passes that object URL to ImportCertificate.
func TestAccAliCloudDataWorksCertificate_basic0(t *testing.T) {
	r := tfacctest.RandIntRange(100000, 999999)
	name := fmt.Sprintf("tf-acc-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudDataWorksCertificateConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_data_works_certificate.default", "id"),
					resource.TestCheckResourceAttr("alicloud_data_works_certificate.default", "name", name),
					resource.TestCheckResourceAttr("alicloud_data_works_certificate.default", "description", "tf-acc-cert"),
					resource.TestCheckResourceAttrSet("alicloud_data_works_certificate.default", "project_id"),
					resource.TestMatchResourceAttr("alicloud_data_works_certificate.default", "create_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T`)),
					resource.TestCheckResourceAttrSet("alicloud_data_works_certificate.default", "create_user"),
					resource.TestCheckResourceAttrSet("alicloud_data_works_certificate.default", "file_size_in_bytes"),
				),
			},
			{
				ResourceName:            "alicloud_data_works_certificate.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate_file"},
			},
		},
	})
}

const testAccDataWorksCertificateFile = `-----BEGIN CERTIFICATE-----
MIICqDCCAZACCQCbNs38rCeB6TANBgkqhkiG9w0BAQsFADAWMRQwEgYDVQQDDAt0
Zi1hY2MtdGVzdDAeFw0yNjA5MDMwMDE4MDFaFw0zNjA4MzEwMDE4MDFaMBYxFDAS
BgNVBAMMC3RmLWFjYy10ZXN0MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEAwyKgyXSNG+NfmE6Ukz7oX4AgDbTaHAyJ3BWHfTK/DJNY0+0tgYE5Zle8bNui
rksdOpydcGAJeaoFNpCGwvd5FO9HklcYwLGdzASRXEHEN0XtG2oeHNnyH82yE2SH
hINbhYcM98WCEuIU5cIsohaQUXhdgBzQcUdc7q5mG/U87LQk4ZGq+JuqUhOAB2xG
8phoguTZ3ODk4+x9MEycpXfpTVkzz9y4JmL6jh3dEZvpyMmC94NYjrGRbioI9CDe
mpY1BsHRtSgxJ7Pvl5yYy7Tj2KC3jRv4FEP7C8vg9tqqV37yLtZsXw/UeVhn9nOm
tOQcqaS3o4BTp//Db7/ReoE2fwIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQBn+Hi1
Lq6pHiV0PPLN+4cQK8sLFSk1nHvJ1ONk9FsBmsnkrbky6I5fNQiBnzOm/rNl72Ao
JvsyROYVyyxX/M5LnZQwahUwDTaz4cNasUiPt33lNtueiLLhzC7UpB8NzH/v29MT
YnEmgc7msKrSnj8SfvZu6y9vbmLhJ0CbrpBTmcG6cXr/Em+l781p4hFdJmpOA7o8
oX+vfX/riPSP8bGgYlmsr2Iv4Uo/b8q6ZpF5p76M/MJirligYjc5fZG8UmsBu6eX
d5LC+4/r6W5uyIS/7fF7FLdzHfRhZknuu8Wzx+HXI+oKv0r/z8MGa1DsvJ6DHQIE
qWre020mtWvGRU1t
-----END CERTIFICATE-----`

func testAccAliCloudDataWorksCertificateConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%[1]s"
}

resource "alicloud_data_works_project" "default" {
  description      = "tf-acc-test"
  project_name     = var.name
  pai_task_enabled = false
  display_name     = var.name
}

resource "alicloud_oss_bucket" "default" {
  bucket = "${var.name}-cert"
  acl    = "public-read"
}

resource "alicloud_oss_bucket_object" "default" {
  bucket  = alicloud_oss_bucket.default.id
  key     = "tf-acc-cert.pem"
  acl     = "public-read"
  content = <<-EOT
%[2]s
EOT
}

resource "alicloud_data_works_certificate" "default" {
  project_id       = alicloud_data_works_project.default.id
  name             = var.name
  description      = "tf-acc-cert"
  certificate_file = "https://${alicloud_oss_bucket.default.id}.${alicloud_oss_bucket.default.extranet_endpoint}/${alicloud_oss_bucket_object.default.key}"
}
`, name, testAccDataWorksCertificateFile)
}
