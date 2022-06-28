package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCSKubernetesAddonsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCSKubernetesAddons-%d", rand)

	resourceId := "data.alicloud_cs_kubernetes_addons.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceCSAddonsConfigDependence(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id": CHECKSET,
					}),
				),
			},
		},
	})
}

func dataSourceCSAddonsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_zones" default {
  available_resource_creation  = "VSwitch"
}

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

# Create a management cluster
resource "alicloud_cs_managed_kubernetes" "default" {
  name                         = var.name
  count                        = 1
  cluster_spec                 = "ack.pro.small"
  is_enterprise_security_group = true
  worker_number                = 2
  deletion_protection          = false
  password                     = "Hello1234"
  pod_cidr                     = "172.20.0.0/16"
  service_cidr                 = "172.21.0.0/20"
  worker_vswitch_ids           = [alicloud_vswitch.default.id]
  worker_instance_types        = [data.alicloud_instance_types.default.instance_types.0.id]
}

data "alicloud_cs_kubernetes_addons" "default" {
  cluster_id = alicloud_cs_managed_kubernetes.default.0.id
}
`, name)
}
