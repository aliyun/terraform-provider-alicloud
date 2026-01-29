package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudMaxComputeTenantRoleUserAttachment_basic12482(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_tenant_role_user_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxComputeTenantRoleUserAttachmentMap12482)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeTenantRoleUserAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmaxcompute%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxComputeTenantRoleUserAttachmentBasicDependence12482)
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
					"account_id":  "p4_200053869413670560",
					"tenant_role": "super_administrator",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_id":  "p4_200053869413670560",
						"tenant_role": "super_administrator",
					}),
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

var AlicloudMaxComputeTenantRoleUserAttachmentMap12482 = map[string]string{}

func AlicloudMaxComputeTenantRoleUserAttachmentBasicDependence12482(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_max_compute_tenant_role_user_attachment" "default0" {
  account_id = "p4_200053869413670560"
  tenant_role = "admin"
} 

`, name)
}
