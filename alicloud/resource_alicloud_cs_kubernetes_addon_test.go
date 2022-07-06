package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCSKubernetesAddon_basic(t *testing.T) {
	var v *Component

	resourceId := "alicloud_cs_kubernetes_addon.default"
	serviceFunc := func() interface{} {
		client, _ := testAccProvider.Meta().(*connectivity.AliyunClient).NewRoaCsClient()
		return &CsClient{client}
	}

	ra := resourceAttrInit(resourceId, csdKubernetesAddonBasicMap)
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rc.describeMethod = "DescribeCsKubernetesAddon"
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccAddon-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSAddonConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id": "${alicloud_cs_managed_kubernetes.default.0.id}",
					"name":       "arms-prometheus",
					"version":    "1.1.6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":   CHECKSET,
						"name":         "arms-prometheus",
						"version":      "1.1.6",
						"next_version": CHECKSET,
						"can_upgrade":  CHECKSET,
						"required":     CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"version": "1.1.7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version": "1.1.7",
					}),
				),
			},
		},
	})
}

var csdKubernetesAddonBasicMap = map[string]string{
	"cluster_id": CHECKSET,
}

func resourceCSAddonConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation  = "VSwitch"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_instance_types" "default" {
	availability_zone          = data.alicloud_zones.default.zones.0.id
	cpu_core_count             = 4
	memory_size                = 8
	kubernetes_node_role       = "Worker"
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
	count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
	vswitch_name = var.name
	vpc_id       = data.alicloud_vpcs.default.ids.0
	cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
	zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_key_pair" "default" {
	key_pair_name = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name                         = var.name
  count                        = 1
  cluster_spec                 = "ack.pro.small"
  is_enterprise_security_group = true
  worker_number                = 2
  password                     = "Hello1234"
  pod_cidr                     = "172.20.0.0/16"
  service_cidr                 = "172.21.0.0/20"
  worker_vswitch_ids           = [local.vswitch_id]
  worker_instance_types        = [data.alicloud_instance_types.default.instance_types.0.id]
  
  maintenance_window {
    enable            = true
    maintenance_time  = "03:00:00Z"
    duration          = "3h"
    weekly_period     = "Thursday"
  }
}
`, name)
}
