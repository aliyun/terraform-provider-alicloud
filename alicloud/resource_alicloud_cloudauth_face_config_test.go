package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCloudauthFaceConfig_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloudauth_face_config.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudauthFaceConfigMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudauthService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudauthFaceConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testaccCAfaceconfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudauthFaceConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.CloudAuthSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"biz_name": name + "-biz",
					"biz_type": name + "-type",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"biz_name": name + "-biz",
						"biz_type": name + "-type",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"biz_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"biz_name": name + "update",
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

var AlicloudCloudauthFaceConfigMap0 = map[string]string{}

func AlicloudCloudauthFaceConfigBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
