package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ESA ClientCertificate. >>> Resource test cases, automatically generated.
// Case resource_ClientCertificate_csr_test
func TestAccAliCloudESAClientCertificateresource_ClientCertificate_csr_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_client_certificate.default"
	ra := resourceAttrInit(resourceId, AliCloudESAClientCertificateresource_ClientCertificate_csr_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaClientCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAClientCertificate%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAClientCertificateresource_ClientCertificate_csr_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":       "${data.alicloud_esa_sites.default.sites.0.id}",
					"csr":           testAccESAClientCertificateConfig,
					"validity_days": "365",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":  "revoked",
					"site_id": "${data.alicloud_esa_sites.default.sites.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"validity_days", "csr"},
			},
		},
	})
}

var AliCloudESAClientCertificateresource_ClientCertificate_csr_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAClientCertificateresource_ClientCertificate_csr_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
  site_name = "gositecdn.cn"
}

`, name)
}

// Case resource_ClientCertificate_set_test
func TestAccAliCloudESAClientCertificateresource_ClientCertificate_set_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_client_certificate.default"
	ra := resourceAttrInit(resourceId, AliCloudESAClientCertificateresource_ClientCertificate_set_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaClientCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAClientCertificate%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAClientCertificateresource_ClientCertificate_set_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":       "${data.alicloud_esa_sites.default.sites.0.id}",
					"pkey_type":     "RSA",
					"validity_days": "365",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":  "revoked",
					"site_id": "${data.alicloud_esa_sites.default.sites.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":  "active",
					"site_id": "${data.alicloud_esa_sites.default.sites.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":  "revoked",
					"site_id": "${data.alicloud_esa_sites.default.sites.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pkey_type", "validity_days"},
			},
		},
	})
}

var AliCloudESAClientCertificateresource_ClientCertificate_set_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAClientCertificateresource_ClientCertificate_set_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
  site_name = "gositecdn.cn"
}

`, name)
}

// Test ESA ClientCertificate. <<< Resource test cases, automatically generated.

const testAccESAClientCertificateConfig = `-----BEGIN CERTIFICATE REQUEST-----\nMIICyDCCAbACAQAwbjELMAkGA1UEBhMCY24xCzAJBgNVBAgMAmJqMQswCQYDVQQH\nDAJiajEMMAoGA1UECgwDYWxpMQwwCgYDVQQLDANhbGkxDDAKBgNVBAMMA2FsaTEb\nMBkGCSqGSIb3DQEJARYMdGVzdEBhbGkuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOC\nAQ8AMIIBCgKCAQEA63a84mqkV8RdsrJVpe55BSJ5tD52UA1vjjtxFHs+P9ozpAeO\n9kBOfCOm+DcrAo8YZcZqAU8pd5AjfNVxnRXaa49Q6IYFrf8AqJcxqEOb6MNOHpao\nEGP+PqrqE1K18QMPeKVtyJ3dDWyzyTVTuiQ7q62sGbebDUcPdya86vWjri/MFajY\nC1XKafIvLArgeaC6iVqvEQIS/k6xKm650v3Kr0x/SCTjrjc8Zs5sQfIkW+7zUobG\nQMZCgHGxfRx8PjCieEuLCMv2hQMApvzkGKcVKOCh9ejqcOXkeNwPsqN8PRF0QqIg\nnEvOzUnh+2sc7yL/goohaRxV0l6TdXz9U4C0GwIDAQABoBUwEwYJKoZIhvcNAQkH\nMQYMBDEwMDAwDQYJKoZIhvcNAQELBQADggEBANCbiTOPAZrWO8vByR2b/YP0gGFA\n7ORSH0Kwp2okC+pjlJ0V8zaYNEm9PL/YNyA8z+jDwU4AXEeE+ZHwFBWNGTKKDV9V\nHS40bpp+X2hQMqhx8xCyTbBW9jvc7yE9IBDfrmkOd6YHW1w6r7+M3yxed2TON9yh\nE/rHqUg3IO+RVEkCwiCBw1AqP4TCAh+ksVKBaNYujh79tW7O2G//edzXBj9lUUth\nLD43KJN+BpqVMMHJEeuzPGur8QMqyvRVwU9d8xaqOA2QEbGYGEVfxrKLWoPAth7d\nzriaUBnAWzodP/S22jm/irVKiqpMopHfq97H5KyxY01Ur4WSRX5FyxwHsg0=\n-----END CERTIFICATE REQUEST-----`
