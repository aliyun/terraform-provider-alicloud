package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_cas_certificate", &resource.Sweeper{
		Name: "alicloud_cas_certificate",
		F:    testSweepCasCertificate,
	})
}

func testSweepCasCertificate(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var allcerts []cas.Certificate
	req := cas.CreateDescribeUserCertificateListRequest()
	req.RegionId = client.RegionId
	req.ShowSize = requests.NewInteger(PageSizeLarge)
	req.CurrentPage = requests.NewInteger(1)
	invoker := NewInvoker()
	for {
		var raw interface{}
		if err := invoker.Run(func() error {
			rsp, err := client.WithCasClient(func(casClient *cas.Client) (interface{}, error) {
				return casClient.DescribeUserCertificateList(req)
			})
			raw = rsp
			return err
		}); err != nil {
			log.Printf("[ERROR] Error retrieving Certificates: %s", WrapError(err))
		}
		resp, _ := raw.(*cas.DescribeUserCertificateListResponse)
		if resp == nil || len(resp.CertificateList) < 1 {
			break
		}
		allcerts = append(allcerts, resp.CertificateList...)

		if len(resp.CertificateList) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.CurrentPage); err != nil {
			log.Printf("[ERROR] %s", WrapError(err))
		} else {
			req.CurrentPage = page
		}
	}

	for _, v := range allcerts {
		name := v.Name
		id := v.Id
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
			}

			if skip {
				log.Printf("[INFO] Skipping Certificate: %s (%d)", name, id)
				continue
			}
			log.Printf("[INFO] Deleting Certificate: %s (%d)", name, id)
			req := cas.CreateDeleteUserCertificateRequest()
			req.CertId = requests.NewInteger(int(id))
			_, err := client.WithCasClient(func(casClient *cas.Client) (interface{}, error) {
				return casClient.DeleteUserCertificate(req)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Certificate (%s (%d)): %s", name, id, err)
			}
		}
	}

	return nil
}

func TestAccAlicloudCasCertificate_basic(t *testing.T) {
	var v *cas.Certificate

	randInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.CasClassicSupportedRegions)
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_cas_certificate.cert",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCasCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCasCertificateConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCasExists("alicloud_cas_certificate.cert", v),
					resource.TestCheckResourceAttr("alicloud_cas_certificate.cert", "name", fmt.Sprintf("tf_testAcc_%v", randInt)),
				),
			},
		},
	})
}

func testAccCheckCasExists(n string, cert *cas.Certificate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No Domain ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		casService := &CasService{client: client}
		certInfo, err := casService.DescribeCas(rs.Primary.Attributes["id"])
		log.Printf("[WARN] Cert id %#v", rs.Primary.ID)

		if err == nil {
			cert = certInfo
			return nil
		}
		return WrapError(err)
	}
}

func testAccCheckCasCertificateDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cas_certificate" {
			continue
		}

		// Try to find the cert
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		casService := &CasService{client: client}
		_, err := casService.DescribeCas(rs.Primary.Attributes["id"])
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return WrapError(err)
		}
	}
	return nil
}

func testAccCasCertificateConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_cas_certificate" "cert" {
  name = "tf_testAcc_%d"
  cert = <<EOF
-----BEGIN CERTIFICATE-----
MIICWDCCAcGgAwIBAgIJAP7vOtjPtQIjMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV
BAYTAkNOMRMwEQYDVQQIDApjbi1iZWlqaW5nMSEwHwYDVQQKDBhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMjAxMDIwMDYxOTUxWhcNMjAxMTE5MDYxOTUxWjBF
MQswCQYDVQQGEwJDTjETMBEGA1UECAwKY24tYmVpamluZzEhMB8GA1UECgwYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKB
gQDEdoyaJ0kdtjtbLRx5X9qwI7FblhJPRcScvhQSE8P5y/b/T8J9BVuFIBoU8nrP
Y9ABz4JFklZ6SznxLbFBqtXoJTmzV6ixyjjH+AGEw6hCiA8Pqy2CNIzxr9DjCzN5
tWruiHqO60O3Bve6cHipH0VyLAhrB85mflvOZSH4xGsJkwIDAQABo1AwTjAdBgNV
HQ4EFgQUYDwuuqC2a2UPrfm1v31vE7+GRM4wHwYDVR0jBBgwFoAUYDwuuqC2a2UP
rfm1v31vE7+GRM4wDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOBgQAovSB0
5JRKrg7lYR/KlTuKHmozfyL9UER0/dpTSoqsCyt8yc1BbtAKUJWh09BujBE1H22f
lKvCAjhPmnNdfd/l9GrmAWNDWEDPLdUTkGSkKAScMpdS+mLmOBuYWgdnOtq3eQGf
t07tlBL+dtzrrohHpfLeuNyYb40g8VQdp3RRRQ==
-----END CERTIFICATE-----
EOF
  key = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDEdoyaJ0kdtjtbLRx5X9qwI7FblhJPRcScvhQSE8P5y/b/T8J9
BVuFIBoU8nrPY9ABz4JFklZ6SznxLbFBqtXoJTmzV6ixyjjH+AGEw6hCiA8Pqy2C
NIzxr9DjCzN5tWruiHqO60O3Bve6cHipH0VyLAhrB85mflvOZSH4xGsJkwIDAQAB
AoGARe2oaCo5lTDK+c4Zx3392hoqQ94r0DmWHPBvNmwAooYd+YxLPrLMe5sMjY4t
dmohnLNevCK1Uzw5eIX6BNSo5CORBcIDRmiAgwiYiS3WOv2+qi9g5uIdMiDr+EED
K8wZJjB5E2WyfxL507vtW4T5L36yfr8SkmqH3GvzpI2jCqECQQDsy0AmBzyfK0tG
Nw1+iF9SReJWgb1f5iHvz+6Dt5ueVQngrl/5++Gp5bNoaQMkLEDsy0iHIj9j43ji
0DON05uDAkEA1GXgGn8MXXKyuzYuoyYXCBH7aF579d7KEGET/jjnXx9DHcfRJZBY
B9ghMnnonSOGboF04Zsdd3xwYF/3OHYssQJAekd/SeQEzyE5TvoQ8t2Tc9X4yrlW
xNX/gmp6/fPr3biGUEtb7qi+4NBodCt+XsingmB7hKUP3RJTk7T2WnAC5wJAMqHi
jY5x3SkFkHl3Hq9q2CKpQxUbCd7FXqg1wum/xj5GmqfSpNjHE3+jUkwbdrJMTrWP
rmRy3tQMWf0mixAo0QJBAN4IcZChanq8cZyNqqoNbxGm4hkxUmE0W4hxHmLC2CYZ
V4JpNm8dpi4CiMWLasF6TYlVMgX+aPxYRUWc/qqf1/Q=
-----END RSA PRIVATE KEY-----
EOF
}
`, randInt)
}
