// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test SslCertificatesService PcaCert. >>> Resource test cases, automatically generated.
// Case 创建私有证书 12331
func TestAccAliCloudSslCertificatesServicePcaCert_basic12331(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ssl_certificates_service_pca_cert.default"
	ra := resourceAttrInit(resourceId, AliCloudSslCertificatesServicePcaCertMap12331)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SslCertificatesServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSslCertificatesServicePcaCert")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsslcertificatesservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSslCertificatesServicePcaCertBasicDependence12331)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"parent_identifier": "${alicloud_ssl_certificates_service_pca_certificate.sub.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parent_identifier": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"upload_flag": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"upload_flag": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alias_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alias_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "REVOKE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "REVOKE",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"after_time", "before_time", "enable_crl", "immediately", "months", "san_type", "san_value", "years"},
			},
		},
	})
}

func TestAccAliCloudSslCertificatesServicePcaCert_basic12331_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ssl_certificates_service_pca_cert.default"
	ra := resourceAttrInit(resourceId, AliCloudSslCertificatesServicePcaCertMap12331)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SslCertificatesServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSslCertificatesServicePcaCert")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsslcertificatesservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSslCertificatesServicePcaCertBasicDependence12331)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"immediately":       "1",
					"organization":      "terraform",
					"years":             "1",
					"upload_flag":       "1",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"locality":          "terraform",
					"months":            "1",
					"custom_identifier": "121",
					"algorithm":         "RSA_1024",
					"parent_identifier": "${alicloud_ssl_certificates_service_pca_certificate.sub.id}",
					"san_value":         "somebody@example.com",
					"enable_crl":        "1",
					"organization_unit": "aliyun",
					"state":             "Beijing",
					"before_time":       "1767948807",
					"days":              "1",
					"san_type":          "1",
					"after_time":        "1768035207",
					"country_code":      "cn",
					"common_name":       "Terraform",
					"alias_name":        name,
					"status":            "REVOKE",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"organization":      "terraform",
						"upload_flag":       "1",
						"resource_group_id": CHECKSET,
						"locality":          "terraform",
						"custom_identifier": "121",
						"algorithm":         "RSA_1024",
						"parent_identifier": CHECKSET,
						"organization_unit": "aliyun",
						"state":             "Beijing",
						"days":              "1",
						"country_code":      "cn",
						"common_name":       "Terraform",
						"alias_name":        name,
						"status":            "REVOKE",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"after_time", "before_time", "enable_crl", "immediately", "months", "san_type", "san_value", "years"},
			},
		},
	})
}

var AliCloudSslCertificatesServicePcaCertMap12331 = map[string]string{
	"algorithm":         CHECKSET,
	"days":              CHECKSET,
	"resource_group_id": CHECKSET,
	"status":            CHECKSET,
}

func AliCloudSslCertificatesServicePcaCertBasicDependence12331(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
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
`, name)
}

// Test SslCertificatesService PcaCert. <<< Resource test cases, automatically generated.
