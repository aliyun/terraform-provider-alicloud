package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAliCloudCSKubernetesAddonMetadataDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCSKubernetesAddonMetadata-%d", rand)

	resourceId := "data.alicloud_cs_kubernetes_addon_metadata.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: dataSourceCSAddonMetadataConfigDependence(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":    CHECKSET,
						"name":          CHECKSET,
						"version":       "v0.16.6",
						"config_schema": REGEXMATCH + `[\s\S]+`,
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
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  cpu_core_count       = 4
  memory_size          = 8
  kubernetes_node_role = "Worker"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
  zone_id      = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  vswitch_name = var.name
}

locals {
  vswitch_id = alicloud_vswitch.vswitch.id
}

# Create a management cluster
resource "alicloud_cs_managed_kubernetes" "default" {
  name                 = var.name
  cluster_spec         = "ack.pro.small"
  vswitch_ids          = [local.vswitch_id]
  pod_vswitch_ids      = [local.vswitch_id]
  new_nat_gateway      = false
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = false
  addons {
    name = "terway-eniip"
  }
  delete_options {
    delete_mode   = "delete"
    resource_type = "ALB"
  }

  delete_options {
    delete_mode   = "delete"
    resource_type = "SLB"
  }

  delete_options {
    delete_mode   = "delete"
    resource_type = "SLS_Data"
  }

  delete_options {
    delete_mode   = "delete"
    resource_type = "SLS_ControlPlane"
  }

  delete_options {
    delete_mode   = "delete"
    resource_type = "PrivateZone"
  }
}

data "alicloud_cs_kubernetes_addon_metadata" "default" {
  cluster_id = alicloud_cs_managed_kubernetes.default.id
  name       = "security-inspector"
  version    = "v0.16.6"
}
`, name)
}
