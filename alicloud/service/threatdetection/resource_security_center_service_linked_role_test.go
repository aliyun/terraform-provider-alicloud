package threatdetection_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/threatdetection"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// lintignore: AT001
func TestAccAliCloudSASServiceLinkedRole_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_security_center_service_linked_role.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSecurityCenterServiceLinkedRoleMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return threatdetection.NewSasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSecurityCenterServiceLinkedRole")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, "", testAccCheckAlicloudSecurityCenterServiceLinkedRoleDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
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

// lintignore: R001
