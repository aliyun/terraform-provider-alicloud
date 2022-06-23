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
					"name":       "ack-virtual-node",
					"version":    "v2.2.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":   CHECKSET,
						"name":         "ack-virtual-node",
						"version":      "v2.2.0",
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
					"version": "v2.3.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version": "v2.3.0",
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
	cpu_core_count             = 2
	memory_size                = 4
	kubernetes_node_role       = "Worker"
}

resource "alicloud_vpc" "default" {
  vpc_name                     = var.name
  cidr_block                   = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
  vswitch_name                 = var.name
  vpc_id                       = alicloud_vpc.default.id
  cidr_block                   = "10.1.1.0/24"
  availability_zone            = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_key_pair" "default" {
	key_name                   = var.name
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
  worker_vswitch_ids           = [alicloud_vswitch.default.id]
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
