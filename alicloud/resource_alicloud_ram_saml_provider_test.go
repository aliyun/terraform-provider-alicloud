package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

//  The test parameter encodedsaml_metadata_document should not be exposed
func SkipTestAccAlicloudRamSamlProvider_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_saml_provider.default"
	ra := resourceAttrInit(resourceId, AlicloudRamSamlProviderMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ImsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamSamlProvider")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudRamSamlProvider%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRamSamlProviderBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"saml_provider_name":            name,
					"encodedsaml_metadata_document": "your encodedsaml metadata document",
					"description":                   "For Terraform Test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"saml_provider_name":            name,
						"encodedsaml_metadata_document": CHECKSET,
						"description":                   "For Terraform Test",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "For Terraform Test Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "For Terraform Test Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encodedsaml_metadata_document": "your encodedsaml metadata document update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encodedsaml_metadata_document": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encodedsaml_metadata_document": "your encodedsaml metadata document",
					"description":                   "For Terraform Test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encodedsaml_metadata_document": CHECKSET,
						"description":                   "For Terraform Test",
					}),
				),
			},
		},
	})
}

var AlicloudRamSamlProviderMap = map[string]string{
	"arn":         CHECKSET,
	"update_date": CHECKSET,
}

func AlicloudRamSamlProviderBasicDependence(name string) string {
	return ""
}
