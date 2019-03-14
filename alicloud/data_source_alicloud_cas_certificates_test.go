package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestCasAlicloudAccountDataSource_certificates(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCasDataSourceCertificates,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cas_certificates.certs"),
					resource.TestCheckResourceAttrSet("data.alicloud_cas_certificates.certs", "id"),
					resource.TestCheckResourceAttrSet("data.alicloud_cas_certificates.certs", "name"),
					resource.TestCheckResourceAttrSet("data.alicloud_cas_certificates.certs", "common"),
					resource.TestCheckResourceAttrSet("data.alicloud_cas_certificates.certs", "finger_print"),
					resource.TestCheckResourceAttrSet("data.alicloud_cas_certificates.certs", "issuer"),
					resource.TestCheckResourceAttrSet("data.alicloud_cas_certificates.certs", "org_name"),
					resource.TestCheckResourceAttrSet("data.alicloud_cas_certificates.certs", "province"),
					resource.TestCheckResourceAttrSet("data.alicloud_cas_certificates.certs", "city"),
					resource.TestCheckResourceAttrSet("data.alicloud_cas_certificates.certs", "country"),
					resource.TestCheckResourceAttrSet("data.alicloud_cas_certificates.certs", "start_date"),
					resource.TestCheckResourceAttrSet("data.alicloud_cas_certificates.certs", "end_date"),
					resource.TestCheckResourceAttrSet("data.alicloud_cas_certificates.certs", "sans"),
					resource.TestCheckResourceAttrSet("data.alicloud_cas_certificates.certs", "expired"),
					resource.TestCheckResourceAttrSet("data.alicloud_cas_certificates.certs", "buy_in_aliyun"),
				),
			},
		},
	})
}

const testAccCheckAlicloudCasDataSourceCertificates = `
resource "alicloud_cas_certificate" "cert" {
   name = "tf_testAcc"
   cert = "./test.crt"
   key = "./test.key"
}
data "alicloud_cas_certificates" "certs" {
  show_size	= "50"
  current_page	= "1"
  lang  = "zh"
  output_file = "./tmp.txt"
}
`
