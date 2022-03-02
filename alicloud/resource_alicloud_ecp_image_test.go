package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudECPImage_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecp_image.default"
	checkoutSupportedRegions(t, true, connectivity.ECPSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudECPImageMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudphoneService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcpImage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secpimage%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECPImageBasicDependence0)
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
					"image_name":  name,
					"description": name,
					"instance_id": "${local.instance_id}",
					"force":       "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name":  name,
						"description": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"instance_id", "force"},
			},
		},
	})
}

var AlicloudECPImageMap0 = map[string]string{
	"force":       CHECKSET,
	"instance_id": CHECKSET,
	"status":      CHECKSET,
}

func AlicloudECPImageBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_ecp_instances" "default" {
}

locals {
  instance_id = data.alicloud_ecp_instances.default.instances[0].instance_id
}
`, name)
}
func TestAccAlicloudECPImage_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecp_image.default"
	checkoutSupportedRegions(t, true, connectivity.ECPSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudECPImageMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudphoneService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcpImage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secpimage%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECPImageBasicDependence1)
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
					"image_name":  name,
					"description": name,
					"force":       "true",
					"instance_id": "${local.instance_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name":  name,
						"description": name,
						"instance_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_name":  name,
					"description": name,
					"force":       "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name":  name,
						"description": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"instance_id", "force"},
			},
		},
	})
}

var AlicloudECPImageMap1 = map[string]string{
	"instance_id": CHECKSET,
	"force":       CHECKSET,
	"status":      CHECKSET,
}

func AlicloudECPImageBasicDependence1(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_ecp_instances" "default" {
}

locals {
  instance_id = data.alicloud_ecp_instances.default.instances[0].instance_id
}
`, name)
}
