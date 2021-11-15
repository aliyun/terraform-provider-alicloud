package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudECDCommand_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecd_command.default"
	ra := resourceAttrInit(resourceId, AlicloudECDCommandMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcdService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcdCommand")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdcommand%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECDCommandBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"command_content": "ipconfig",
					"command_type":    "RunBatScript",
					"desktop_id":      "ecd-dpteqvk58v8b80ujo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"command_content": "ipconfig",
						"command_type":    "RunBatScript",
						"desktop_id":      "ecd-dpteqvk58v8b80ujo",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeout", "content_encoding"},
			},
		},
	})
}

var AlicloudECDCommandMap0 = map[string]string{
	"content_encoding": NOSET,
	"status":           CHECKSET,
	"timeout":          NOSET,
	"desktop_id":       CHECKSET,
}

func AlicloudECDCommandBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
