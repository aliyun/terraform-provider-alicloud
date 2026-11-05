package cs_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCSKubernetesAddonsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCSKubernetesAddons-%d", rand)

	resourceId := "data.alicloud_cs_kubernetes_addons.default"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: dataSourceCSAddonsConfigDependence(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":                   acctest.CHECKSET,
						"names.#":                      acctest.CHECKSET,
						"addons.#":                     acctest.CHECKSET,
						"addons.0.name":                acctest.CHECKSET,
						"addons.0.supported_actions.#": acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudCSKubernetesAddonsDataSource_installed(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCSKubernetesAddons-%d", rand)

	resourceId := "data.alicloud_cs_kubernetes_addons.installed-metrics-server"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: dataSourceCSAddonsConfigDependence(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":                   acctest.CHECKSET,
						"ids.#":                        acctest.CHECKSET,
						"name_regex":                   acctest.CHECKSET,
						"addons.#":                     "1",
						"names.#":                      "1",
						"addons.0.name":                "metrics-server",
						"addons.0.current_config":      acctest.REGEXMATCH + "^.+$",
						"addons.0.current_version":     acctest.REGEXMATCH + "^.+$",
						"addons.0.next_version":        "",
						"addons.0.required":            acctest.CHECKSET,
						"addons.0.supported_actions.#": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: dataSourceCSAddonsConfigDependence(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id": acctest.REGEXMATCH + ".+",
						"ids.#":      acctest.REGEXMATCH + ".+",
						"name_regex": acctest.REGEXMATCH + ".+",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudCSKubernetesAddonsDataSource_notInstalled(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCSKubernetesAddons-%d", rand)

	resourceId := "data.alicloud_cs_kubernetes_addons.not-installed-migrate-controller"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: dataSourceCSAddonsConfigDependence(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":                   acctest.CHECKSET,
						"names.#":                      "1",
						"addons.#":                     "1",
						"addons.0.name":                "migrate-controller",
						"addons.0.current_config":      "",
						"addons.0.current_version":     "",
						"addons.0.next_version":        acctest.REGEXMATCH + "^.+$",
						"addons.0.required":            acctest.CHECKSET,
						"addons.0.supported_actions.#": acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudCSKubernetesAddonsDataSource_installedWithUpgradeVersion(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCSKubernetesAddons-%d", rand)

	resourceId := "data.alicloud_cs_kubernetes_addons.installed-alb-ingress-controller-with-upgradable-version"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: dataSourceCSAddonsConfigDependence(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":                   acctest.CHECKSET,
						"names.#":                      "1",
						"addons.#":                     "1",
						"addons.0.name":                "alb-ingress-controller",
						"addons.0.current_config":      "{}",
						"addons.0.current_version":     "v2.19.0",
						"addons.0.next_version":        acctest.REGEXMATCH + "^.+$",
						"addons.0.required":            acctest.CHECKSET,
						"addons.0.supported_actions.#": acctest.CHECKSET,
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

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

locals {
  vswitch_id = alicloud_vswitch.vswitch.id
}

# Create a management cluster
resource "alicloud_cs_managed_kubernetes" "default" {
  name                         = var.name
  count                        = 1
  cluster_spec                 = "ack.pro.small"
  is_enterprise_security_group = true
  deletion_protection          = false
  service_cidr                 = cidrsubnet("172.16.0.0/16", 4, 3)
  vswitch_ids                  = [local.vswitch_id]
  pod_vswitch_ids              = [local.vswitch_id]
  new_nat_gateway              = false
  slb_internet_enabled         = false
  addons {
    name = "terway-eniip"
  }
  addons {
    name = "metrics-server"
    config = jsonencode({
      MemoryRequest = "500Mi"
      CpuRequest    = "250m"
      MemoryLimit   = "8Gi"
      CpuLimit      = "4"
    })
  }
  addons {
    name    = "alb-ingress-controller"
    version = "v2.19.0"
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

data "alicloud_cs_kubernetes_addons" "default" {
  cluster_id = alicloud_cs_managed_kubernetes.default.0.id
}

data "alicloud_cs_kubernetes_addons" "installed-metrics-server" {
  cluster_id = alicloud_cs_managed_kubernetes.default.0.id
  name_regex = "^metrics-server"
}

data "alicloud_cs_kubernetes_addons" "not-installed-migrate-controller" {
  cluster_id = alicloud_cs_managed_kubernetes.default.0.id
  name_regex = "^migrate-controller"
}

data "alicloud_cs_kubernetes_addons" "installed-alb-ingress-controller-with-upgradable-version" {
  cluster_id = alicloud_cs_managed_kubernetes.default.0.id
  name_regex = "^alb-ingress-controller"
}
`, name)
}
