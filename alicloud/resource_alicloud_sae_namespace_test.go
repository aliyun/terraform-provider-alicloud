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
					"namespace_name":            name,
					"namespace_id":              fmt.Sprintf("%s:tftest%d", defaultRegionToTest, rand),
					"namespace_description":     "tftestaccdescription",
					"enable_micro_registration": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_name":            name,
						"namespace_id":              fmt.Sprintf("%s:tftest%d", defaultRegionToTest, rand),
						"namespace_description":     "tftestaccdescription",
						"enable_micro_registration": "false",
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
				Config: testAccConfig(map[string]interface{}{
					"enable_micro_registration": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_micro_registration": "true",
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

func TestAccAlicloudSAENamespace_basic1(t *testing.T) {
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
					"namespace_name":            name,
					"namespace_short_id":        fmt.Sprintf("tftest%d", rand),
					"namespace_description":     "tftestaccdescription",
					"enable_micro_registration": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_name":            name,
						"namespace_short_id":        CHECKSET,
						"namespace_description":     "tftestaccdescription",
						"enable_micro_registration": "false",
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
				Config: testAccConfig(map[string]interface{}{
					"enable_micro_registration": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_micro_registration": "true",
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
	"namespace_id":              CHECKSET,
	"namespace_short_id":        CHECKSET,
	"enable_micro_registration": CHECKSET,
}

func AlicloudSAENamespaceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}
`, name)
}
