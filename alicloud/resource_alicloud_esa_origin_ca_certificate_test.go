package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ESA OriginCaCertificate. >>> Resource test cases, automatically generated.
// Case resource_OriginCaCertificate_set_test
func TestAccAliCloudESAOriginCaCertificateresource_OriginCaCertificate_set_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_origin_ca_certificate.default"
	ra := resourceAttrInit(resourceId, AliCloudESAOriginCaCertificateresource_OriginCaCertificate_set_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaOriginCaCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAOriginCaCertificate%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAOriginCaCertificateresource_OriginCaCertificate_set_testBasicDependence)
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
					"site_id":     "${alicloud_esa_site.resource_Site_OriginCaCertificateset_test.id}",
					"certificate": testEsaClientCaCertificate,
					"name":        "test",
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudESAOriginCaCertificateresource_OriginCaCertificate_set_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAOriginCaCertificateresource_OriginCaCertificate_set_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "resource_Site_OriginCaCertificateset_test" {
  site_name   = "bcd.com"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Test ESA OriginCaCertificate. <<< Resource test cases, automatically generated.
