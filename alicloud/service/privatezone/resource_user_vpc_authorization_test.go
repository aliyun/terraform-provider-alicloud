package privatezone_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/privatezone"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudPvtzUserVpcAuthorization_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pvtz_user_vpc_authorization.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPrivateZoneUserVpcAuthorizationMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return privatezone.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeUserVpcAuthorization")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sprivatezoneuservpcauthorization%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudPrivateZoneUserVpcAuthorizationBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			privatezoneTestPreCheckEnterpriseAccountEnabled(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"authorized_user_id": "1359528198477648",
					"auth_channel":       "RESOURCE_DIRECTORY",
					"auth_type":          "NORMAL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"authorized_user_id": acctest.CHECKSET,
						"auth_type":          "NORMAL",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"auth_channel"},
			},
		},
	})
}

// defaultAuthType exercises the bug-fix path: auth_type is omitted so the
// resource ID must default to "<authorized_user_id>:NORMAL" instead of the
// previous buggy "<authorized_user_id>:<nil>".
func TestAccAliCloudPvtzUserVpcAuthorization_defaultAuthType(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pvtz_user_vpc_authorization.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPrivateZoneUserVpcAuthorizationMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return privatezone.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeUserVpcAuthorization")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sprivatezoneuservpcauthorization%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudPrivateZoneUserVpcAuthorizationBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			privatezoneTestPreCheckEnterpriseAccountEnabled(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"authorized_user_id": "1359528198477648",
					"auth_channel":       "RESOURCE_DIRECTORY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"authorized_user_id": acctest.CHECKSET,
						"auth_type":          "NORMAL",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"auth_channel"},
			},
		},
	})
}

var AlicloudPrivateZoneUserVpcAuthorizationMap0 = map[string]string{}

func AlicloudPrivateZoneUserVpcAuthorizationBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

// lintignore: R001
