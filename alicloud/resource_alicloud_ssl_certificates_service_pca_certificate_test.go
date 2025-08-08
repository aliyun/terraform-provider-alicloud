// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test SslCertificatesService PcaCertificate. >>> Resource test cases, automatically generated.
// Case PcaCertificate 11010
func TestAccAliCloudSslCertificatesServicePcaCertificate_basic11010(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ssl_certificates_service_pca_certificate.default"
	ra := resourceAttrInit(resourceId, AliCloudSslCertificatesServicePcaCertificateMap11010)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SslCertificatesServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSslCertificatesServicePcaCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsslcertificatesservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSslCertificatesServicePcaCertificateBasicDependence11010)
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
					"common_name":       "cbc.certqa.cn",
					"locality":          "a",
					"organization":      "a",
					"organization_unit": "a",
					"state":             "a",
					"years":             "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"common_name":       "cbc.certqa.cn",
						"locality":          "a",
						"organization":      "a",
						"organization_unit": "a",
						"state":             "a",
						"years":             "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudSslCertificatesServicePcaCertificate_basic11010_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ssl_certificates_service_pca_certificate.default"
	ra := resourceAttrInit(resourceId, AliCloudSslCertificatesServicePcaCertificateMap11010)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SslCertificatesServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSslCertificatesServicePcaCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsslcertificatesservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSslCertificatesServicePcaCertificateBasicDependence11010)
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
					"organization":      "a",
					"years":             "1",
					"locality":          "a",
					"organization_unit": "a",
					"state":             "a",
					"country_code":      "cn",
					"common_name":       "cbc.certqa.cn",
					"algorithm":         "RSA_1024",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"organization":      "a",
						"years":             "1",
						"locality":          "a",
						"organization_unit": "a",
						"state":             "a",
						"country_code":      "cn",
						"common_name":       "cbc.certqa.cn",
						"algorithm":         "RSA_1024",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudSslCertificatesServicePcaCertificateMap11010 = map[string]string{
	"status": CHECKSET,
}

func AliCloudSslCertificatesServicePcaCertificateBasicDependence11010(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test SslCertificatesService PcaCertificate. <<< Resource test cases, automatically generated.
