package edas_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	edassvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_edas_k8s_cluster", &resource.Sweeper{
		Name: "alicloud_edas_k8s_cluster",
		F:    testSweepEdasK8sCluster,
	})
}

func testSweepEdasK8sCluster(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}

	clusterListRq := edas.CreateListClusterRequest()
	clusterListRq.RegionId = region

	raw, err := client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
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
		if !sweep.SweepAll() {
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

		wait := helper.IncrementalWait(1*time.Second, 2*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			raw, err := client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
				return edasClient.DeleteCluster(deleteClusterRq)
			})
			if err != nil {
				if errs.IsExpectedErrors(err, []string{errs.ThrottlingUser}) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(deleteClusterRq.GetActionName(), raw, deleteClusterRq.RoaRequest, deleteClusterRq)
			rsp := raw.(*edas.DeleteClusterResponse)
			if rsp.Code == 601 && strings.Contains(rsp.Message, "Operation cannot be processed because there are running instances.") {
				err = errs.Error("Operation cannot be processed because there are running instances.")
				return retry.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete edas cluster (%s): %s", name, err)
		}
	}

	return nil
}

func TestAccAliCloudEdasK8sCluster_basic(t *testing.T) {
	var v *edas.Cluster
	resourceId := "alicloud_edas_k8s_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, edasK8sClusterBasicMap)
	serviceFunc := func() interface{} {
		return edassvc.NewEdasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	rand := tfacctest.RandIntRange(0, 1000)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edask8sclusterbasic%v", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEdasK8sClusterConfigDependence)
	region := "cn-beijing"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEdasK8sClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cs_cluster_id": "${alicloud_cs_managed_kubernetes.default.id}",
					"namespace_id":  region,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name":  acctest.CHECKSET,
						"cluster_type":  "5",
						"cs_cluster_id": acctest.CHECKSET,
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
	"cluster_name":          acctest.CHECKSET,
	"cluster_type":          acctest.CHECKSET,
	"network_mode":          acctest.CHECKSET,
	"cluster_import_status": acctest.CHECKSET,
}

func testAccCheckEdasK8sClusterDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)

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

		raw, err := client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.GetCluster(request)
		})

		if err != nil {
			if errs.NotFoundError(err) {
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
		  cpu_core_count       = 2
		  memory_size          = 4
		  kubernetes_node_role = "Worker"
		  instance_type_family = "ecs.c7a"
		  system_disk_category = "cloud_essd"
		}

		data "alicloud_cs_kubernetes_version" "default" {
		  cluster_type       = "ManagedKubernetes"
		  profile            = "Default"
		  kubernetes_version = "1.34"
		}

		resource "alicloud_vpc" "default" {
		  name       = var.name
		  cidr_block = "172.16.0.0/12"
		}

		resource "alicloud_vswitch" "default" {
		  name       = var.name
		  vpc_id     = alicloud_vpc.default.id
		  cidr_block = "172.16.0.0/24"
		  zone_id    = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
		}

		resource "alicloud_cs_managed_kubernetes" "default" {
		  name_prefix          = var.name
		  cluster_spec         = "ack.pro.small"
		  version              = data.alicloud_cs_kubernetes_version.default.metadata.0.version
		  vswitch_ids          = [alicloud_vswitch.default.id]
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
		  vswitch_ids          = [alicloud_vswitch.default.id]
		  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
		  system_disk_category = "cloud_essd"
		  system_disk_size     = 40
		  password             = "Test12345"
		  desired_size         = 2
		}

		`, name)
}
