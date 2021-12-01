package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudOOSApplicationGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oos_application_group.default"
	checkoutSupportedRegions(t, true, connectivity.OOSApplicationSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudOOSApplicationGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOosApplicationGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%soosapplicationgroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOOSApplicationGroupBasicDependence0)
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
					"deploy_region_id":       os.Getenv("ALICLOUD_REGION"),
					"application_name":       "${alicloud_oos_application.default.id}",
					"application_group_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deploy_region_id":       CHECKSET,
						"application_name":       CHECKSET,
						"application_group_name": name,
						"import_tag_key":         "app-" + name,
						"import_tag_value":       name,
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

func TestAccAlicloudOOSApplicationGroup_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oos_application_group.default"
	checkoutSupportedRegions(t, true, connectivity.OOSApplicationSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudOOSApplicationGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOosApplicationGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%soosapplicationgroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOOSApplicationGroupBasicDependence0)
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

					"deploy_region_id":       os.Getenv("ALICLOUD_REGION"),
					"description":            "${var.name}",
					"application_name":       "${alicloud_oos_application.default.id}",
					"application_group_name": "${var.name}",
					"import_tag_key":         "${var.name}",
					"import_tag_value":       "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"deploy_region_id":       CHECKSET,
						"description":            name,
						"application_name":       CHECKSET,
						"application_group_name": name,
						"import_tag_key":         name,
						"import_tag_value":       name,
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

var AlicloudOOSApplicationGroupMap0 = map[string]string{
	"application_group_name": CHECKSET,
	"application_name":       CHECKSET,
}

func AlicloudOOSApplicationGroupBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}


data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_oos_application" "default" {
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  application_name  = var.name
  description       = var.name
  tags = {
    Created = "TF"
  }
}
`, name)
}
