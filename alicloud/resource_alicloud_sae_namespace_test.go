package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudSAENamespace_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_namespace.default"
	ra := resourceAttrInit(resourceId, AlicloudSAENamespaceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSaeNamespace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%ssaenamespace%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSAENamespaceBasicDependence0)
	resource.Test(t, resource.TestCase{
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_name":        name,
					"namespace_id":          fmt.Sprintf("%s:tftest%d", defaultRegionToTest, rand),
					"namespace_description": "tftestaccdescription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_name":        name,
						"namespace_id":          fmt.Sprintf("%s:tftest%d", defaultRegionToTest, rand),
						"namespace_description": "tftestaccdescription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_name":        name + "update",
						"namespace_description": "tftestaccdescription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_description": "tftestaccdescriptionupdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_name":        name + "update",
						"namespace_description": "tftestaccdescriptionupdate",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_name":        name + "updateall",
					"namespace_description": "tftestaccdescriptionupdateall",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_name":        name + "updateall",
						"namespace_description": "tftestaccdescriptionupdateall",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudSAENamespaceMap0 = map[string]string{
	"namespace_name":        CHECKSET,
	"namespace_id":          CHECKSET,
	"namespace_description": CHECKSET,
}

func AlicloudSAENamespaceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
