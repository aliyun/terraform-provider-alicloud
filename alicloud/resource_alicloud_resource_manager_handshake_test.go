package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudResourceManagerHandshake_basic(t *testing.T) {
	var v resourcemanager.Handshake
	resourceId := "alicloud_resource_manager_handshake.default"
	ra := resourceAttrInit(resourceId, ResourceManagerHandshakeMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourcemanagerService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerHandshake")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccResourceManagerHandshake%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ResourceManagerHandshakeBasicdependence)
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
					"target_entity": os.Getenv("ALICLOUD_ACCOUNT_ID"),
					"target_type":   "Account",
					"note":          "test resource manager handshake",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_entity": os.Getenv("ALICLOUD_ACCOUNT_ID"),
						"target_type":   "Account",
						"note":          "test resource manager handshake",
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

var ResourceManagerHandshakeMap = map[string]string{}

func ResourceManagerHandshakeBasicdependence(name string) string {
	return ""
}
