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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	edassvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_edas_cluster", &resource.Sweeper{
		Name: "alicloud_edas_cluster",
		F:    testSweepEdasCluster,
		Dependencies: []string{
			"alicloud_edas_application",
		},
	})
}

func testSweepEdasCluster(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
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
		log.Printf("[ERROR] Failed to retrieve edas cluster in service list: %s", listClusterResponse)
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

		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			raw, err := client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
				return edasClient.DeleteCluster(deleteClusterRq)
			})
			if err != nil {
				if errs.IsExpectedErrors(err, []string{errs.ThrottlingUser}) {
					time.Sleep(5 * time.Second)
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

func TestAccAliCloudEdasCluster_basic(t *testing.T) {
	var v *edas.Cluster
	resourceId := "alicloud_edas_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, edasClusterBasicMap)
	serviceFunc := func() interface{} {
		return edassvc.NewEdasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	rand := tfacctest.RandInt()
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edasclusterbasic%v", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEdasClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEdasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_name": "${var.name}",
					"cluster_type": "2",
					"network_mode": "2",
					"vpc_id":       "${data.alicloud_vpcs.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name": name,
					}),
				),
			},

			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_name": fmt.Sprintf("tf-testacc-edasclusterchange%v", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name": fmt.Sprintf("tf-testacc-edasclusterchange%v", rand)}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_type": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_type": "3"}),
				),
			},
		},
	})
}

func TestAccAliCloudEdasCluster_multi(t *testing.T) {
	var v *edas.Cluster
	resourceId := "alicloud_edas_cluster.default.1"
	ra := acctest.ResourceAttrInit(resourceId, edasClusterBasicMap)
	serviceFunc := func() interface{} {
		return edassvc.NewEdasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	rand := tfacctest.RandInt()
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edasclustermulti%v", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEdasClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEdasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":        "2",
					"cluster_name": "${var.name}-${count.index}",
					"cluster_type": "2",
					"network_mode": "2",
					"vpc_id":       "${data.alicloud_vpcs.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

var edasClusterBasicMap = map[string]string{
	"cluster_name": acctest.CHECKSET,
	"cluster_type": acctest.CHECKSET,
	"network_mode": acctest.CHECKSET,
}

func testAccCheckEdasClusterDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_edas_cluster" {
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
		if rsp.Cluster.ClusterId != "" {
			return fmt.Errorf("cluster %s still exist", rsp.Cluster.ClusterId)
		}
	}

	return nil
}

func resourceEdasClusterConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		  default = "%v"
		}

		data "alicloud_vpcs" "default" {
			name_regex = "^default-NODELETING$"
		}
		`, name)
}
