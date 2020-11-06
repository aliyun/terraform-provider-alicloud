package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCSServerlessKubernetes_basic(t *testing.T) {
	var v *cs.ServerlessClusterResponse

	resourceId := "alicloud_cs_serverless_kubernetes.default"
	ra := resourceAttrInit(resourceId, csServerlessKubernetesBasicMap)

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccserverlesskubernetes-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSServerlessKubernetesConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ServerlessKubernetesSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                           name,
					"vpc_id":                         "${alicloud_vpc.default.id}",
					"vswitch_id":                     "${alicloud_vswitch.default.id}",
					"new_nat_gateway":                "true",
					"deletion_protection":            "false",
					"endpoint_public_access_enabled": "true",
					"resource_group_id":              "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                           name,
						"new_nat_gateway":                "true",
						"deletion_protection":            "false",
						"endpoint_public_access_enabled": "true",
						"resource_group_id":              CHECKSET,
					}),
				),
			},
		},
	})
}

func resourceCSServerlessKubernetesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_zones" default {
  available_resource_creation = "VSwitch"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
`, name)
}

var csServerlessKubernetesBasicMap = map[string]string{
	"new_nat_gateway":                "true",
	"deletion_protection":            "false",
	"endpoint_public_access_enabled": "true",
	"force_update":                   "false",
}
