package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCSKubernetesAddonMetadataDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCSKubernetesAddonMetadata-%d", rand)

	resourceId := "data.alicloud_cs_kubernetes_addon_metadata.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceCSAddonMetadataConfigDependence(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":    CHECKSET,
						"config_schema": CHECKSET,
					}),
				),
			},
		},
	})
}

func dataSourceCSAddonMetadataConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation  = "VSwitch"
}

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
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
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
  pod_cidr                     = "10.99.0.0/16"
  service_cidr                 = "192.168.0.0/16"
  worker_vswitch_ids           = [local.vswitch_id]
  worker_instance_types        = [data.alicloud_instance_types.default.instance_types.0.id]
}

data "alicloud_cs_kubernetes_addon_metadata" "default" {
  cluster_id = alicloud_cs_managed_kubernetes.default.0.id
  name = "migrate-controller"
  version = "v1.6.3-6fd55d8-aliyun"
}
`, name)
}
