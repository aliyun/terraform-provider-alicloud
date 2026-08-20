// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test DataWorks Certificate. >>> Resource test cases, automatically generated.
// Case 测试Certificate1（有前置步骤，上海预发执行） 9054
func TestAccAliCloudDataWorksCertificate_basic9054(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_certificate.default"
	ra := resourceAttrInit(resourceId, AlicloudDataWorksCertificateMap9054)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataWorksServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDataWorksCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdataworks%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDataWorksCertificateBasicDependence9054)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shenzhen"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project_id":       "${alicloud_data_works_project.defaultxWQGsP.id}",
					"name":             name,
					"certificate_file": "http://${alicloud_oss_bucket.cert.id}.${alicloud_oss_bucket.cert.extranet_endpoint}/${alicloud_oss_bucket_object.cert.key}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_id":       CHECKSET,
						"name":             CHECKSET,
						"certificate_file": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate_file"},
			},
		},
	})
}

var AlicloudDataWorksCertificateMap9054 = map[string]string{
	"create_time":        CHECKSET,
	"create_user":        CHECKSET,
	"file_size_in_bytes": CHECKSET,
}

func AlicloudDataWorksCertificateBasicDependence9054(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_data_works_project" "defaultxWQGsP" {
  description      = "tf-acc-test"
  project_name     = var.name
  pai_task_enabled = false
  display_name     = var.name
}

resource "alicloud_oss_bucket" "cert" {
  bucket = "${var.name}-cert"
  acl    = "public-read"
}

resource "alicloud_oss_bucket_object" "cert" {
  bucket  = alicloud_oss_bucket.cert.id
  key     = "tf-acc-cert.pem"
  acl     = "public-read"
  content = <<-EOT
-----BEGIN CERTIFICATE-----
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
-----END CERTIFICATE-----
EOT
}

`, name)
}

// Case 测试Certificate2（没有前置步骤，能够运行就行） 9055
func TestAccAliCloudDataWorksCertificate_basic9055(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_certificate.default"
	ra := resourceAttrInit(resourceId, AlicloudDataWorksCertificateMap9055)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataWorksServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDataWorksCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdataworks%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDataWorksCertificateBasicDependence9055)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shenzhen"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":      "tf-acc-cert",
					"project_id":       "${alicloud_data_works_project.defaultxWQGsP.id}",
					"name":             name,
					"certificate_file": "http://${alicloud_oss_bucket.cert.id}.${alicloud_oss_bucket.cert.extranet_endpoint}/${alicloud_oss_bucket_object.cert.key}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":      "tf-acc-cert",
						"project_id":       CHECKSET,
						"name":             CHECKSET,
						"certificate_file": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate_file"},
			},
		},
	})
}

var AlicloudDataWorksCertificateMap9055 = map[string]string{
	"create_time":        CHECKSET,
	"create_user":        CHECKSET,
	"file_size_in_bytes": CHECKSET,
}

func AlicloudDataWorksCertificateBasicDependence9055(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_data_works_project" "defaultxWQGsP" {
  description      = "tf-acc-test"
  project_name     = var.name
  pai_task_enabled = false
  display_name     = var.name
}

resource "alicloud_oss_bucket" "cert" {
  bucket = "${var.name}-cert"
  acl    = "public-read"
}

resource "alicloud_oss_bucket_object" "cert" {
  bucket  = alicloud_oss_bucket.cert.id
  key     = "tf-acc-cert.pem"
  acl     = "public-read"
  content = <<-EOT
-----BEGIN CERTIFICATE-----
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
-----END CERTIFICATE-----
EOT
}

`, name)
}

// Test DataWorks Certificate. <<< Resource test cases, automatically generated.
