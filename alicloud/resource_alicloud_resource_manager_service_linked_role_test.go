package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudResourceManagerServiceLinkedRole_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_service_linked_role.default"
	ra := resourceAttrInit(resourceId, AlicloudResourceManagerServiceLinkedRoleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamServiceLinkedRole")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sresourcemanagerservicelinkedrole%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudResourceManagerServiceLinkedRoleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_name": "csb.aliyuncs.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_id":   CHECKSET,
						"role_name": CHECKSET,
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description", "custom_suffix"},
			},
		},
	})
}

var AlicloudResourceManagerServiceLinkedRoleMap0 = map[string]string{}

func AlicloudResourceManagerServiceLinkedRoleBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
default = "%s"
}
`, name)
}
