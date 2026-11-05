package sslcertificateservice_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sslcertificateservice"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudSslCertificatesServiceCloudAccess_basic(t *testing.T) {
	resourceId := "alicloud_ssl_certificates_service_cloud_access.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	secretId := fmt.Sprintf("tf-testacc-cloud-access-secret-id-%d", rand)
	secretKey := fmt.Sprintf("tf-testacc-cloud-access-secret-key-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckSslCertificatesServiceCloudAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudSslCertificatesServiceCloudAccessBasic(secretId, secretKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttrSet(resourceId, "access_id"),
					resource.TestCheckResourceAttr(resourceId, "cloud_name", "Tencent"),
					resource.TestCheckResourceAttr(resourceId, "secret_id", secretId),
					resource.TestCheckResourceAttr(resourceId, "secret_key", secretKey),
					resource.TestCheckResourceAttrSet(resourceId, "service_status"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_key"},
			},
		},
	})
}

func testAccCheckSslCertificatesServiceCloudAccessDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ssl_certificates_service_cloud_access" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		_, err := sslcertificateservice.NewCasService(client).DescribeCloudAccess(rs.Primary.ID, "", "")
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing cloud access %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("cloud access %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccAliCloudSslCertificatesServiceCloudAccessBasic(secretId, secretKey string) string {
	return fmt.Sprintf(`
resource "alicloud_ssl_certificates_service_cloud_access" "default" {
  cloud_name = "Tencent"
  secret_id  = "%s"
  secret_key = "%s"
}
`, secretId, secretKey)
}
