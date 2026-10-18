package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/alibabacloud-go/cs-20151215/v8/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func init() {
	resource.AddTestSweepers("alicloud_cs_kubernetes_addon", &resource.Sweeper{
		Name: "alicloud_cs_kubernetes_addon",
		F:    testSweepCSKubernetesAddon,
	})
}

// testSweepCSKubernetesAddon sweeps leftover CS Kubernetes addons in the
// region. Sub-resources inherit protection from their parent cluster's name:
// NODELETING markers on the cluster name (or addon name) always skip, EVEN
// when SWEEP_ALL=1 is set (red-line #2). The prefix check on the cluster name
// is bypassed by sweepAll() (standard pattern, playbook line 55). Only
// non-default addons (InstallByDefault == false) are uninstalled —
// system-critical addons installed at cluster creation time (flannel,
// kube-proxy, coredns, ...) are always skipped. Per-addon uninstall failures
// are logged and continued (non-fatal). ACK prepaid clusters that cannot be
// deleted themselves (ForbiddenRelease, red-line #4) still have their test
// addons swept here.
func testSweepCSKubernetesAddon(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf-example",
		"tf_example",
		"terraform-example",
	}

	raw, err := aliyunClient.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
		return csClient.DescribeClusters("")
	})
	if err != nil {
		return fmt.Errorf("error retrieving CS Clusters: %s", err)
	}
	clusters, _ := raw.([]cs.ClusterType)

	for _, cluster := range clusters {
		clusterName := cluster.Name
		clusterId := cluster.ClusterID

		// NODELETING guard on the parent cluster name (always-on, NOT bypassed
		// by sweepAll — red-line #2).
		if protectedSweepVpcName(clusterName) {
			log.Printf("[INFO] Skipping CS addons in protected cluster: %s (%s)", clusterName, clusterId)
			continue
		}

		// Prefix check on the parent cluster name (bypassed by sweepAll —
		// standard pattern, playbook line 55). Sub-resources of non-test
		// clusters are spared unless sweepAll() forces it.
		if !sweepAll() {
			matched := false
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(clusterName), strings.ToLower(prefix)) {
					matched = true
					break
				}
			}
			if !matched {
				log.Printf("[INFO] Skipping CS addons in non-test cluster: %s (%s)", clusterName, clusterId)
				continue
			}
		}

		roaClient, err := aliyunClient.NewRoaCsClient()
		if err != nil {
			log.Printf("[ERROR] Failed to build ROA CS client for cluster %s (%s): %s", clusterName, clusterId, err)
			continue
		}
		csClient := CsClient{roaClient}

		availableAddons, err := csClient.ListAddons(clusterId)
		if err != nil {
			log.Printf("[ERROR] Failed to list addons for cluster %s (%s): %s", clusterName, clusterId, err)
			continue
		}

		for name, addonInfo := range availableAddons {
			// Skip system-critical addons installed by default at cluster
			// creation (flannel, kube-proxy, coredns, ...) to avoid breaking
			// the cluster.
			if addonInfo != nil && addonInfo.InstallByDefault != nil && tea.BoolValue(addonInfo.InstallByDefault) {
				continue
			}
			if protectedSweepVpcName(name) {
				log.Printf("[INFO] Skipping protected CS addon: %s in cluster %s", name, clusterId)
				continue
			}

			instance, err := csClient.GetCsKubernetesAddonInstance(clusterId, name)
			if err != nil {
				if NotFoundError(err) {
					// Addon not installed in this cluster — nothing to sweep.
					continue
				}
				log.Printf("[ERROR] Failed to describe addon %s in cluster %s (%s): %s", name, clusterName, clusterId, err)
				continue
			}
			if instance == nil || instance.Version == "" {
				continue
			}

			log.Printf("[INFO] Uninstalling CS addon: %s in cluster %s (%s)", name, clusterName, clusterId)
			body := make([]*client.UnInstallClusterAddonsRequestAddons, 0)
			body = append(body, &client.UnInstallClusterAddonsRequestAddons{
				Name: tea.String(name),
			})
			uninstallArgs := &client.UnInstallClusterAddonsRequest{
				Addons: body,
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			if err := retry.Retry(5*time.Minute, func() *retry.RetryError {
				_, err := csClient.client.UnInstallClusterAddons(&clusterId, uninstallArgs)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			}); err != nil {
				log.Printf("[ERROR] Failed to uninstall CS addon (%s in cluster %s (%s)): %s", name, clusterName, clusterId, err)
			}
		}
	}
	return nil
}

// lintignore: AT001
func TestAccAliCloudCSKubernetesAddon_apiserver(t *testing.T) {
	var v *Component

	resourceId := "alicloud_cs_kubernetes_addon.apiserver"
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		//CheckDestroy:  rac.checkResourceDestroy(),
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
						"cluster_id":   CHECKSET,
						"name":         "kube-apiserver",
						"version":      CHECKSET,
						"next_version": CHECKSET,
						"can_upgrade":  CHECKSET,
						"required":     CHECKSET,
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
						"config": CHECKSET,
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
	var v *Component

	resourceId := "alicloud_cs_kubernetes_addon.terway-eniip"
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{
				// import
				Config: testAccConfig(map[string]interface{}{
					"cluster_id": "${local.cluster_id}",
					"name":       "terway-eniip",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":   CHECKSET,
						"name":         "terway-eniip",
						"version":      CHECKSET,
						"next_version": CHECKSET,
						"can_upgrade":  CHECKSET,
						"required":     CHECKSET,
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
						"config": CHECKSET,
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
	var v *Component

	resourceId := "alicloud_cs_kubernetes_addon.logtail-ds"
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
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
						"cluster_id":   CHECKSET,
						"name":         "logtail-ds",
						"version":      "v1.8.5.0-aliyun",
						"next_version": CHECKSET,
						"can_upgrade":  CHECKSET,
						"required":     CHECKSET,
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
						"config": CHECKSET,
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
	var v *Component

	resourceId := "alicloud_cs_kubernetes_addon.ack-node-problem-detector"
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id": "${local.cluster_id}",
					"name":       "ack-node-problem-detector",
					"depends_on": []string{"alicloud_cs_kubernetes_node_pool.default"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":   CHECKSET,
						"name":         "ack-node-problem-detector",
						"next_version": CHECKSET,
						"can_upgrade":  CHECKSET,
						"required":     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config": "{\\\"npd\\\":{\\\"resources\\\":{\\\"limits\\\":{\\\"cpu\\\":\\\"2.0\\\",\\\"memory\\\":\\\"800Mi\\\"}}}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config": "{\\\"npd\\\":{\\\"resources\\\":{\\\"limits\\\":{\\\"cpu\\\":\\\"1.0\\\",\\\"memory\\\":\\\"800Mi\\\"}}}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config": CHECKSET,
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
	var v *Component

	resourceId := "alicloud_cs_kubernetes_addon.ack-virtual-node"
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id": "${local.cluster_id}",
					"name":       "ack-virtual-node",
					//"version":    "v2.10.5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id": CHECKSET,
						"name":       "ack-virtual-node",
						//"version":      "v2.10.5",
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
	"cluster_id": CHECKSET,
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
