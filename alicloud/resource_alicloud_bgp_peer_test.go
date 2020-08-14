package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudBgpPeer_basic(t *testing.T) {
	var v vpc.BgpPeer

	resourceId := "alicloud_bgp_peer.default"
	ra := resourceAttrInit(resourceId, BgpPeerbasicMap)

	serviceFunc := func() interface{} {
		return &BgpService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sBgpPeerbasic%v", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceBgpPeerConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bgp_group_id":    "${alicloud_bgp_group.foo.id}",
					"peer_ip_address": "192.168.8.4",
					"enable_bfd":      "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bgp_group_id":    CHECKSET,
						"peer_ip_address": "192.168.8.4",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"peer_ip_address": "192.168.8.5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"peer_ip_address": "192.168.8.5"}),
				),
			},
		},
	})
}

func resourceBgpPeerConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_bgp_group" "foo" {
    peer_asn = 2
    name = "%s"
    router_id = "%s"
    description = "test-description11"
    is_fake_asn = true
    auth_key= "dasdasda"
}
`, name, os.Getenv("ALICLOUD_VBR_ID"))
}

var BgpPeerbasicMap = map[string]string{
	"peer_ip_address": "192.168.4.0",
	"enable_bfd":      "false",
	"bgp_group_id":    CHECKSET,
}
