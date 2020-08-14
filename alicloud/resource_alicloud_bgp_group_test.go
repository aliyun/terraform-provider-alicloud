package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudBgpGroup_basic(t *testing.T) {
	var v vpc.BgpGroup

	resourceId := "alicloud_bgp_group.default"
	ra := resourceAttrInit(resourceId, BgpGroupbasicMap)

	serviceFunc := func() interface{} {
		return &BgpService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 20)
	name := fmt.Sprintf("tf-testacc%sBgpNetworkbasic%v", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceBgpGroupConfigDependence)

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
					"peer_asn":    "${var.peer_asn}",
					"router_id":   "${var.router_id}",
					"description": "${var.description[0]}",
					"name":        "${var.name}",
					"is_fake_asn": "${var.is_fake_asn}",
					"auth_key":    "${var.auth_key}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(BgpGroupbasicMap),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auth_key": "fake_auth_key",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"auth_key": "fake_auth_key"}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"is_fake_asn": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"is_fake_asn": "true"}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "terraform-test-name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": "terraform-test-name"}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "terraform-test-description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"description": "terraform-test-description"}),
				),
			},
		},
	})
}

func TestAccAlicloudBgpGroup_multi(t *testing.T) {
	var v vpc.BgpGroup

	resourceId := "alicloud_bgp_group.default.1"
	ra := resourceAttrInit(resourceId, BgpGroupbasicMap)

	serviceFunc := func() interface{} {
		return &BgpService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 20)
	name := fmt.Sprintf("tf-testacc%sBgpNetworkMulti%v", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceBgpGroupConfigDependence)

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
					"peer_asn":    "${var.peer_asn}",
					"router_id":   "${var.router_id}",
					"description": "${element(var.description, count.index)}",
					"name":        "${var.name}",
					"is_fake_asn": "${var.is_fake_asn}",
					"auth_key":    "${var.auth_key}",
					"count":       "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "test_description1",
						"bgp_group_id": CHECKSET,
						"router_id":    CHECKSET,
					}),
				),
			},
		},
	})
}

func resourceBgpGroupConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
variable "router_id" {
 default = "%s"
}
variable "is_fake_asn" {
 default = true
}
variable "description" {
 default = ["test_description", "test_description1"]
}
variable "peer_asn" {
 default = 2
}
variable "auth_key" {
 default = "test_auth_key"
}
`, name, os.Getenv("ALICLOUD_VBR_ID"))
}

var BgpGroupbasicMap = map[string]string{
	"peer_asn":     "2",
	"router_id":    CHECKSET,
	"auth_key":     "test_auth_key",
	"description":  "test_description",
	"name":         CHECKSET,
	"is_fake_asn":  "true",
	"bgp_group_id": CHECKSET,
}
