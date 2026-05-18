package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_edas_k8s_cluster", &resource.Sweeper{
		Name: "alicloud_edas_k8s_cluster",
		F:    testSweepEdasK8sCluster,
	})
}

func testSweepEdasK8sCluster(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}

	clusterListRq := edas.CreateListClusterRequest()
	clusterListRq.RegionId = region

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListCluster(clusterListRq)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve edas cluster in service list: %s", err)
		return nil
	}

	listClusterResponse, _ := raw.(*edas.ListClusterResponse)
	if listClusterResponse.Code != 200 {
		log.Printf("[ERROR] Failed to retrieve edas cluster in service list: %s", listClusterResponse.Message)
		return nil
	}

	for _, v := range listClusterResponse.ClusterList.Cluster {
		name := v.ClusterName
		skip := true
		if !sweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}

			if skip {
				log.Printf("[INFO] Skipping edas cluster: %s", name)
				continue
			}
		}
		log.Printf("[INFO] delete edas cluster: %s", name)

		deleteClusterRq := edas.CreateDeleteClusterRequest()
		deleteClusterRq.RegionId = region
		deleteClusterRq.ClusterId = v.ClusterId

		wait := incrementalWait(1*time.Second, 2*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
				return edasClient.DeleteCluster(deleteClusterRq)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(deleteClusterRq.GetActionName(), raw, deleteClusterRq.RoaRequest, deleteClusterRq)
			rsp := raw.(*edas.DeleteClusterResponse)
			if rsp.Code == 601 && strings.Contains(rsp.Message, "Operation cannot be processed because there are running instances.") {
				err = Error("Operation cannot be processed because there are running instances.")
				return resource.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete edas cluster (%s): %s", name, err)
		}
	}

	return nil
}

func TestAccAlicloudEdasK8sCluster_basic(t *testing.T) {
	var v *edas.Cluster
	resourceId := "alicloud_edas_k8s_cluster.default"
	ra := resourceAttrInit(resourceId, edasK8sClusterBasicMap)
	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(0, 1000)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edask8sclusterbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasK8sClusterConfigDependence)
	region := "cn-beijing"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  testAccCheckEdasK8sClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cs_cluster_id": "${alicloud_cs_managed_kubernetes.default.id}",
					"namespace_id":  region,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name":  CHECKSET,
						"cluster_type":  "5",
						"cs_cluster_id": CHECKSET,
						"namespace_id":  region,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var edasK8sClusterBasicMap = map[string]string{
	"cluster_name":          CHECKSET,
	"cluster_type":          CHECKSET,
	"network_mode":          CHECKSET,
	"cluster_import_status": CHECKSET,
}

func testAccCheckEdasK8sClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	edasService := EdasService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_edas_k8s_cluster" {
			continue
		}

		// Try to find the cluster
		clusterId := rs.Primary.ID
		regionId := client.RegionId

		request := edas.CreateGetClusterRequest()
		request.RegionId = regionId
		request.ClusterId = clusterId

		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.GetCluster(request)
		})

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		rsp := raw.(*edas.GetClusterResponse)
		// ClusterImportStatus 4 indicates the EDAS-side registration has been marked
		// for deletion (record may linger in GetCluster while ACK tears down).
		if rsp.Cluster.ClusterId != "" && rsp.Cluster.ClusterImportStatus != 4 {
			return fmt.Errorf("cluster %s still exist", rsp.Cluster.ClusterId)
		}
	}

	return nil
}

func resourceEdasK8sClusterConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		  default = "%v"
		}

		data "alicloud_zones" default {
		  available_resource_creation = "VSwitch"
		}

		data "alicloud_instance_types" "default" {
		  availability_zone    = data.alicloud_zones.default.zones.0.id
		  cpu_core_count       = 2
		  memory_size          = 4
		  kubernetes_node_role = "Worker"
		}

		data "alicloud_cs_kubernetes_version" "default" {
		  cluster_type       = "ManagedKubernetes"
		  profile            = "Default"
		  kubernetes_version = "1.34"
		}

		data "alicloud_vpcs" "default" {
		  name_regex = "^default-NODELETING$"
		}

		data "alicloud_vswitches" "default" {
		  vpc_id  = data.alicloud_vpcs.default.ids.0
		  zone_id = data.alicloud_zones.default.zones.0.id
		}

		resource "alicloud_vswitch" "vswitch" {
		  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
		  vpc_id       = data.alicloud_vpcs.default.ids.0
		  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
		  zone_id      = data.alicloud_zones.default.zones.0.id
		  vswitch_name = var.name
		}

		locals {
		  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
		}

		resource "alicloud_cs_managed_kubernetes" "default" {
		  name_prefix          = var.name
		  cluster_spec         = "ack.pro.small"
		  version              = data.alicloud_cs_kubernetes_version.default.metadata.0.version
		  vswitch_ids          = [local.vswitch_id]
		  new_nat_gateway      = false
		  pod_cidr             = "10.123.0.0/16"
		  service_cidr         = "192.168.0.0/20"
		  slb_internet_enabled = true
		  addons {
		    name = "storage-operator"
		  }
		  addons {
		    name = "csi-plugin"
		  }
		  addons {
		    name = "csi-provisioner"
		  }
		}

		resource "alicloud_cs_kubernetes_node_pool" "default" {
		  name                 = "desired_size"
		  cluster_id           = alicloud_cs_managed_kubernetes.default.id
		  vswitch_ids          = [local.vswitch_id]
		  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
		  system_disk_category = "cloud_efficiency"
		  system_disk_size     = 40
		  password             = "Test12345"
		  desired_size         = 2
		}

		`, name)
}
