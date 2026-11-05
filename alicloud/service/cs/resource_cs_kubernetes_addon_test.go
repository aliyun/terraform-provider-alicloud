package cs_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// lintignore: AT001
func TestAccAliCloudCSKubernetesAddon_apiserver(t *testing.T) {
	var v *cs.Component

	resourceId := "alicloud_cs_kubernetes_addon.apiserver"
	serviceFunc := func() interface{} {
		client, _ := acctest.Provider.Meta().(*connectivity.AliyunClient).NewRoaCsClient()
		return cs.NewCsClient(client)
	}

	ra := acctest.ResourceAttrInit(resourceId, csdKubernetesAddonBasicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rc.DescribeMethod = "DescribeCsKubernetesAddon"
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccAddon-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSAddonConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		//CheckDestroy:  rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// import
				Config: testAccConfig(map[string]interface{}{
					"cluster_id": "${local.cluster_id}",
					"name":       "kube-apiserver",
					"config":     "{\\\"requestTimeout\\\":\\\"2m0s\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":   acctest.CHECKSET,
						"name":         "kube-apiserver",
						"version":      acctest.CHECKSET,
						"next_version": acctest.CHECKSET,
						"can_upgrade":  acctest.CHECKSET,
						"required":     acctest.CHECKSET,
					}),
				),
			},
			{
				// config
				Config: testAccConfig(map[string]interface{}{
					"config": "{\\\"requestTimeout\\\":\\\"3m10s\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config"},
			},
		},
	})
}

// lintignore: AT001
func TestAccAliCloudCSKubernetesAddon_terway_eniip(t *testing.T) {
	var v *cs.Component

	resourceId := "alicloud_cs_kubernetes_addon.terway-eniip"
	serviceFunc := func() interface{} {
		client, _ := acctest.Provider.Meta().(*connectivity.AliyunClient).NewRoaCsClient()
		return cs.NewCsClient(client)
	}

	ra := acctest.ResourceAttrInit(resourceId, csdKubernetesAddonBasicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rc.DescribeMethod = "DescribeCsKubernetesAddon"
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccAddon-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSAddonConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				// import
				Config: testAccConfig(map[string]interface{}{
					"cluster_id": "${local.cluster_id}",
					"name":       "terway-eniip",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":   acctest.CHECKSET,
						"name":         "terway-eniip",
						"version":      acctest.CHECKSET,
						"next_version": acctest.CHECKSET,
						"can_upgrade":  acctest.CHECKSET,
						"required":     acctest.CHECKSET,
					}),
				),
			},
			{
				// config
				Config: testAccConfig(map[string]interface{}{
					"config": "{\\\"MaxPoolSize\\\":3}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config"},
			},
		},
	})
}

// create new addon, upgrade, update and delete
// lintignore: AT001
func TestAccAliCloudCSKubernetesAddon_logtail_ds(t *testing.T) {
	var v *cs.Component

	resourceId := "alicloud_cs_kubernetes_addon.logtail-ds"
	serviceFunc := func() interface{} {
		client, _ := acctest.Provider.Meta().(*connectivity.AliyunClient).NewRoaCsClient()
		return cs.NewCsClient(client)
	}

	ra := acctest.ResourceAttrInit(resourceId, csdKubernetesAddonBasicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rc.DescribeMethod = "DescribeCsKubernetesAddon"
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccAddon-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSAddonConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id": "${local.cluster_id}",
					"name":       "logtail-ds",
					"version":    "v1.8.5.0-aliyun",
					"depends_on": []string{"alicloud_cs_kubernetes_node_pool.default"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":   acctest.CHECKSET,
						"name":         "logtail-ds",
						"version":      "v1.8.5.0-aliyun",
						"next_version": acctest.CHECKSET,
						"can_upgrade":  acctest.CHECKSET,
						"required":     acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"version": "v1.8.6.0-aliyun",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version": "v1.8.6.0-aliyun",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config": "{\\\"LogtailDSLimitCPU\\\":\\\"3\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config"},
			},
		},
	})
}

// create new addon, upgrade, update and delete
// lintignore: AT001
func TestAccAliCloudCSKubernetesAddon_ack_node_problem_detector(t *testing.T) {
	var v *cs.Component

	resourceId := "alicloud_cs_kubernetes_addon.ack-node-problem-detector"
	serviceFunc := func() interface{} {
		client, _ := acctest.Provider.Meta().(*connectivity.AliyunClient).NewRoaCsClient()
		return cs.NewCsClient(client)
	}

	ra := acctest.ResourceAttrInit(resourceId, csdKubernetesAddonBasicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rc.DescribeMethod = "DescribeCsKubernetesAddon"
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccAddon-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSAddonConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id": "${local.cluster_id}",
					"name":       "ack-node-problem-detector",
					"depends_on": []string{"alicloud_cs_kubernetes_node_pool.default"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":   acctest.CHECKSET,
						"name":         "ack-node-problem-detector",
						"next_version": acctest.CHECKSET,
						"can_upgrade":  acctest.CHECKSET,
						"required":     acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config": "{\\\"npd\\\":{\\\"resources\\\":{\\\"limits\\\":{\\\"cpu\\\":\\\"2.0\\\",\\\"memory\\\":\\\"800Mi\\\"}}}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config": "{\\\"npd\\\":{\\\"resources\\\":{\\\"limits\\\":{\\\"cpu\\\":\\\"1.0\\\",\\\"memory\\\":\\\"800Mi\\\"}}}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config"},
			},
		},
	})
}

// create new addon, and delete with cleanup_cloud_resources
// lintignore: AT001
func TestAccAliCloudCSKubernetesAddon_vk(t *testing.T) {
	var v *cs.Component

	resourceId := "alicloud_cs_kubernetes_addon.ack-virtual-node"
	serviceFunc := func() interface{} {
		client, _ := acctest.Provider.Meta().(*connectivity.AliyunClient).NewRoaCsClient()
		return cs.NewCsClient(client)
	}

	ra := acctest.ResourceAttrInit(resourceId, csdKubernetesAddonBasicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rc.DescribeMethod = "DescribeCsKubernetesAddon"
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccAddon-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSAddonConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id": "${local.cluster_id}",
					"name":       "ack-virtual-node",
					//"version":    "v2.10.5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id": acctest.CHECKSET,
						"name":       "ack-virtual-node",
						//"version":      "v2.10.5",
						"next_version": acctest.CHECKSET,
						"can_upgrade":  acctest.CHECKSET,
						"required":     acctest.CHECKSET,
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
					"cleanup_cloud_resources": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
		},
	})
}

var csdKubernetesAddonBasicMap = map[string]string{
	"cluster_id": acctest.CHECKSET,
}

func resourceCSAddonConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_instance_types" "default" {
  cpu_core_count    = 4
  memory_size       = 8
  eni_amount        = 4
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


resource "alicloud_cs_managed_kubernetes" "default" {
  name_prefix                  = var.name
  cluster_spec                 = "ack.pro.small"
  vswitch_ids                  = [local.vswitch_id]
  pod_vswitch_ids              = [local.vswitch_id]
  new_nat_gateway              = false
  service_cidr                 = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled         = false
  is_enterprise_security_group = true
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

resource "alicloud_cs_kubernetes_node_pool" "default" {
  node_pool_name       = var.name
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [local.vswitch_id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_essd"
  system_disk_size     = 40
  desired_size         = 2
}

locals {
  vswitch_id = alicloud_vswitch.vswitch.id
  cluster_id = alicloud_cs_managed_kubernetes.default.id
}
`, name)
}
