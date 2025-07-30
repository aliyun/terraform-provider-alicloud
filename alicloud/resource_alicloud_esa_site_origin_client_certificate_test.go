package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ESA SiteOriginClientCertificate. >>> Resource test cases, automatically generated.
// Case siteoriginclientcertificate_test
func TestAccAliCloudESASiteOriginClientCertificatesiteoriginclientcertificate_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_site_origin_client_certificate.default"
	ra := resourceAttrInit(resourceId, AliCloudESASiteOriginClientCertificatesiteoriginclientcertificate_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaSiteOriginClientCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESASiteOriginClientCertificate%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESASiteOriginClientCertificatesiteoriginclientcertificate_testBasicDependence)
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
					"site_id":                             "${alicloud_esa_site.resource_Site_SiteOriginClientCertificate_test.id}",
					"private_key":                         testFcPrivateKey,
					"site_origin_client_certificate_name": "testCertificate",
					"certificate":                         testFcCertificate,
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
				ImportStateVerifyIgnore: []string{"private_key"},
			},
		},
	})
}

var AliCloudESASiteOriginClientCertificatesiteoriginclientcertificate_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESASiteOriginClientCertificatesiteoriginclientcertificate_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "resource_Site_SiteOriginClientCertificate_test" {
  site_name   = "chenxin0116.site"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Test ESA SiteOriginClientCertificate. <<< Resource test cases, automatically generated.
