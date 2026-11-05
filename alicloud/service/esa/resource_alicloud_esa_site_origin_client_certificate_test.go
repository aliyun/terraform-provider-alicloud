package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	esa "github.com/aliyun/terraform-provider-alicloud/alicloud/service/esa"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ESA SiteOriginClientCertificate. >>> Resource test cases, automatically generated.
// Case siteoriginclientcertificate_test
func TestAccAliCloudESASiteOriginClientCertificatesiteoriginclientcertificate_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_site_origin_client_certificate.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudESASiteOriginClientCertificatesiteoriginclientcertificate_testMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaSiteOriginClientCertificate")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESASiteOriginClientCertificate%d", acctest.DefaultRegionToTest, rand)

	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudESASiteOriginClientCertificatesiteoriginclientcertificate_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
	"id": acctest.CHECKSET,
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
  site_name   = "${lower(var.name)}.com"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Test ESA SiteOriginClientCertificate. <<< Resource test cases, automatically generated.
