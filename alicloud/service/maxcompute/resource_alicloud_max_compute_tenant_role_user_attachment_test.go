package maxcompute_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/maxcompute"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudMaxComputeTenantRoleUserAttachment_basic12482(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_tenant_role_user_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudMaxComputeTenantRoleUserAttachmentMap12482)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return maxcompute.NewMaxComputeServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMaxComputeTenantRoleUserAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmaxcompute%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudMaxComputeTenantRoleUserAttachmentBasicDependence12482)
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
	// the principal must already be a member of the MaxCompute/DataWorks tenant
	// (a freshly created RAM user fails with "user not exist in dw"), and it
	// permanently holds the admin tenant role in the test account, so only the
	// super_administrator attachment is exercised here
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

`, name)
}
