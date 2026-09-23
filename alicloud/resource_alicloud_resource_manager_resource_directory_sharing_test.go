package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ResourceManager ResourceDirectorySharing. >>> Resource test cases, automatically generated.
// Case resource_directory_sharing 11900
func TestAccAliCloudResourceManagerResourceDirectorySharing_basic11900(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_resource_directory_sharing.default"
	ra := resourceAttrInit(resourceId, AlicloudResourceManagerResourceDirectorySharingMap11900)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourceManagerServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerResourceDirectorySharing")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccresourcemanager%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudResourceManagerResourceDirectorySharingBasicDependence11900)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_sharing_with_rd": "true",
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

var AlicloudResourceManagerResourceDirectorySharingMap11900 = map[string]string{
	"enable_sharing_with_rd": CHECKSET,
}

func AlicloudResourceManagerResourceDirectorySharingBasicDependence11900(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test ResourceManager ResourceDirectorySharing. <<< Resource test cases, automatically generated.
