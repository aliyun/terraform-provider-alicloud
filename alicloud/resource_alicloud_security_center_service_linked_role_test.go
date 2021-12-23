package alicloud

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudSecurityCenterServiceLinkedRole_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_security_center_service_linked_role.default"
	ra := resourceAttrInit(resourceId, AlicloudSecurityCenterServiceLinkedRoleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSecurityCenterServiceLinkedRole")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", testAccCheckAlicloudSecurityCenterServiceLinkedRoleDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "true",
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

var AlicloudSecurityCenterServiceLinkedRoleMap0 = map[string]string{}

func testAccCheckAlicloudSecurityCenterServiceLinkedRoleDependence(name string) string {
	return ""
}
