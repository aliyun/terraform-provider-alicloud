package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/resourcemanager"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ResourceManager ResourceDirectorySharing. >>> Resource test cases, automatically generated.
// Case resource_directory_sharing 11900
func TestAccAliCloudResourceManagerResourceDirectorySharing_basic11900(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_resource_directory_sharing.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudResourceManagerResourceDirectorySharingMap11900)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return resourcemanager.NewResourceManagerServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeResourceManagerResourceDirectorySharing")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccresourcemanager%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudResourceManagerResourceDirectorySharingBasicDependence11900)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
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
	"enable_sharing_with_rd": acctest.CHECKSET,
}

func AlicloudResourceManagerResourceDirectorySharingBasicDependence11900(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test ResourceManager ResourceDirectorySharing. <<< Resource test cases, automatically generated.
