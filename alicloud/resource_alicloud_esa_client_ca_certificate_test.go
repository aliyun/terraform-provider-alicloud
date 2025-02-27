package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ESA ClientCaCertificate. >>> Resource test cases, automatically generated.
// Case resource_ClientCaCertificate_set_test
func TestAccAliCloudESAClientCaCertificateresource_ClientCaCertificate_set_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_client_ca_certificate.default"
	ra := resourceAttrInit(resourceId, AliCloudESAClientCaCertificateresource_ClientCaCertificate_set_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaClientCaCertificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAClientCaCertificate%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAClientCaCertificateresource_ClientCaCertificate_set_testBasicDependence)
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
					"client_ca_cert_name": "test",
					"site_id":             "${data.alicloud_esa_sites.default.sites.0.id}",
					"certificate":         testEsaClientCaCertificate,
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

var AliCloudESAClientCaCertificateresource_ClientCaCertificate_set_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAClientCaCertificateresource_ClientCaCertificate_set_testBasicDependence(name string) string {
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

// Test ESA ClientCaCertificate. <<< Resource test cases, automatically generated.

const testEsaClientCaCertificate = `-----BEGIN CERTIFICATE-----\nMIIDQDCCAiigAwIBAgIUBRXhIKYZoD503x0GIRcJjRT1ik8wDQYJKoZIhvcNAQEF\nBQAwMjEVMBMGA1UEAwwMZ29zaXRlY2RuLmNuMQwwCgYDVQQKDANhbGkxCzAJBgNV\nBAYTAkNOMB4XDTI1MDIxODA2NDIwMFoXDTI1MDczMDA2MDAwMFowYDELMAkGA1UE\nBhMCQ04xCzAJBgNVBAgMAmJqMQswCQYDVQQHDAJiajENMAsGA1UECgwEdGVzdDER\nMA8GA1UECwwIYWxpYmFiYWIxFTATBgNVBAMMDGdvc2l0ZWNkbi5jbjCCASIwDQYJ\nKoZIhvcNAQEBBQADggEPADCCAQoCggEBALAazzh1i7SBJI+kDntmbRDFMYquOAU/\nEx3XKP5UHiTkDczkHyI6h6dHMo2QD4U9hYRLnbtOKXLLGCU9AtwKUNBUQppTC/Oj\nLnOsxRmei5g2wtuHNmBTfivY6QfdJDppomoriv2EAG3aC4moJg+uFoeY350bohqh\n18c9JrpCuEHoyqf34aYUEETRtfW1PicnvW9ZhnlDe8YVEpRyDxTJBZSc8Hhr8vHE\nZqqTM991yI7p4koH1x4lXnNPu6T7WrdUo2AxoEF2ODZpurJwCoKsezrCV91DSYtX\n75mMvcifD/arDflzF66s8bKHefWdr7wvT5zLVFhEAXwxXW4pWg/Tj68CAwEAAaMg\nMB4wCwYDVR0RBAQwAoIAMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQEFBQAD\nggEBAKu40aIYYCHbOMARUArRlIhTtRaXWqaPeyJEWT0n7fW85t/i5gU9McKK5LlX\nbsj2zXjRxexW15WisG+ZTSGBi5nG/BhS6lrqeCH8gAj1Ei7n1mO+aTO3R/o6cdIf\nDkoR2OyHdMZGYflqKeePBulA5yn8ZeUr0E0Kci4a/Qsj+LZ5mdB8R05j8v+OwFuJ\nU02pA2oHptVZxGrJLY+hyymyubIK+vQZXsPYwYBO96A8F+JJEhBScfG+JpyGCjuT\nR4RlhRpxtKbucz+3ebgaJsNC5Kit1usMfohbgM4LvvS+ODlkpwdYVZNVw+2IhC4L\nliiXejj/9o7w7270QQdY7e0wBfI=\n-----END CERTIFICATE-----`
