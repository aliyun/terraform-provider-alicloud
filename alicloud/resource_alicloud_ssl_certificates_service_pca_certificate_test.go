package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test SslCertificatesService PcaCertificate. >>> Resource test cases, automatically generated.
// Case PcaCertificate_sub 12329
func TestAccAliCloudSslCertificatesServicePcaCertificate_basic12329(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ssl_certificates_service_pca_certificate.default"
	ra := resourceAttrInit(resourceId, AliCloudSslCertificatesServicePcaCertificateMap12329)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SslCertificatesServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSslCertificatesServicePcaCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsslcertificatesservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSslCertificatesServicePcaCertificateBasicDependence12329)
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
					"parent_identifier": "${alicloud_ssl_certificates_service_pca_certificate.root.id}",
					"common_name":       "cbc.certqa.cn",
					"locality":          "a",
					"organization":      "a",
					"organization_unit": "a",
					"state":             "a",
					"years":             "1",
					"algorithm":         "RSA_2048",
					"certificate_type":  "SUB_ROOT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parent_identifier": CHECKSET,
						"common_name":       "cbc.certqa.cn",
						"locality":          "a",
						"organization":      "a",
						"organization_unit": "a",
						"state":             "a",
						"years":             "1",
						"algorithm":         "RSA_2048",
						"certificate_type":  "SUB_ROOT",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alias_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alias_name", "enable_crl", "path_len_constraint", "extended_key_usages"},
			},
		},
	})
}

func TestAccAliCloudSslCertificatesServicePcaCertificate_basic12329_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ssl_certificates_service_pca_certificate.default"
	ra := resourceAttrInit(resourceId, AliCloudSslCertificatesServicePcaCertificateMap12329)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SslCertificatesServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSslCertificatesServicePcaCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsslcertificatesservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSslCertificatesServicePcaCertificateBasicDependence12329)
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
					"parent_identifier":   "${alicloud_ssl_certificates_service_pca_certificate.root.id}",
					"organization":        "a",
					"years":               "1",
					"locality":            "a",
					"organization_unit":   "a",
					"state":               "a",
					"country_code":        "cn",
					"common_name":         "cbc.certqa.cn",
					"algorithm":           "RSA_2048",
					"certificate_type":    "SUB_ROOT",
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"crl_day":             "2",
					"enable_crl":          "true",
					"alias_name":          name,
					"path_len_constraint": "1",
					"extended_key_usages": []string{
						"serverAuth"},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parent_identifier": CHECKSET,
						"organization":      "a",
						"years":             "1",
						"locality":          "a",
						"organization_unit": "a",
						"state":             "a",
						"country_code":      "cn",
						"common_name":       "cbc.certqa.cn",
						"algorithm":         "RSA_2048",
						"certificate_type":  "SUB_ROOT",
						"resource_group_id": CHECKSET,
						"crl_day":           "2",
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
				ImportStateVerifyIgnore: []string{"alias_name", "enable_crl", "path_len_constraint", "extended_key_usages"},
			},
		},
	})
}

var AliCloudSslCertificatesServicePcaCertificateMap12329 = map[string]string{
	"status": CHECKSET,
}

func AliCloudSslCertificatesServicePcaCertificateBasicDependence12329(name string) string {
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
`, name)
}

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
				Config: testAccConfig(map[string]interface{}{
					"alias_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alias_name"},
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
					"certificate_type":  "ROOT",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"alias_name":        name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
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
						"certificate_type":  "ROOT",
						"resource_group_id": CHECKSET,
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
				ImportStateVerifyIgnore: []string{"alias_name"},
			},
		},
	})
}

var AliCloudSslCertificatesServicePcaCertificateMap11010 = map[string]string{
	"algorithm":         CHECKSET,
	"resource_group_id": CHECKSET,
	"status":            CHECKSET,
	"certificate_type":  CHECKSET,
}

func AliCloudSslCertificatesServicePcaCertificateBasicDependence11010(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

`, name)
}

// Test SslCertificatesService PcaCertificate. <<< Resource test cases, automatically generated.
