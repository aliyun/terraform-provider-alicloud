package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudBgpNetwork_basic(t *testing.T) {
	var v vpc.BgpNetwork

	resourceId := "alicloud_bgp_network.default"
	ra := resourceAttrInit(resourceId, BgpNetworkbasicMap)

	serviceFunc := func() interface{} {
		return &BgpService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sBgpNetworkbasic%v", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceBgpNetworkConfigDependence)

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
					"dst_cidr_block": "${var.dst_cidr_blocks.0}",
					"router_id":      "${var.router_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dst_cidr_block": "192.168.2.11",
						"router_id":      CHECKSET,
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

func TestAccAlicloudBgpNetwork_multi(t *testing.T) {
	var v vpc.BgpNetwork

	resourceId := "alicloud_bgp_network.default.1"
	ra := resourceAttrInit(resourceId, BgpNetworkbasicMap)

	serviceFunc := func() interface{} {
		return &BgpService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sBgpNetworkMulti%v", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceBgpNetworkConfigDependence)

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
					"dst_cidr_block": "${element(var.dst_cidr_blocks, count.index)}",
					"router_id":      "${var.router_id}",
					"count":          "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dst_cidr_block": "192.168.2.10",
						"router_id":      CHECKSET,
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

func resourceBgpNetworkConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
variable "router_id" {
 default = "%s"
}
variable "dst_cidr_blocks" {
 default = ["192.168.2.11", "192.168.2.10"]
}
`, name, os.Getenv("ALICLOUD_VBR_ID"))
}

var BgpNetworkbasicMap = map[string]string{
	"dst_cidr_block": "192.168.2.11",
	"router_id":      CHECKSET,
}
