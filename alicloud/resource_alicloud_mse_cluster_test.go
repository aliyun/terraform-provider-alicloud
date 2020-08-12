package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/mse"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_mse_cluster", &resource.Sweeper{
		Name: "alicloud_mse_cluster",
		F:    testSweepMSECluster,
	})
}

func testSweepMSECluster(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	request := mse.CreateListClustersRequest()
	request.PageSize = requests.NewInteger(PageSizeXLarge)
	request.PageNum = requests.NewInteger(1)
	raw, err := client.WithMseClient(func(MseClient *mse.Client) (interface{}, error) {
		return MseClient.ListClusters(request)
	})
	if err != nil {
		log.Printf("[ERROR] Error retrieving Mse Clusters: %s", WrapError(err))
	}
	response, _ := raw.(*mse.ListClustersResponse)

	sweeped := false
	for _, v := range response.Data {
		id := v.InstanceId
		name := v.ClusterAliasName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Mse Clusters: %s (%s)", name, id)
			continue
		}

		sweeped = true
		log.Printf("[INFO] Deleting Mse Clusters: %s (%s)", name, id)
		req := mse.CreateDeleteClusterRequest()
		req.InstanceId = id
		_, err := client.WithMseClient(func(MseClient *mse.Client) (interface{}, error) {
			return MseClient.DeleteCluster(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Mse Clusters (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		// Waiting 30 seconds to ensure these Mse Clusters have been deleted.
		time.Sleep(30 * time.Second)
	}
	return nil
}

func TestAccAlicloudMSECluster_basic(t *testing.T) {
	var v mse.Data
	resourceId := "alicloud_mse_cluster.default"
	ra := resourceAttrInit(resourceId, MseClusterMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMseCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccMseCluster%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, MseClusterBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithNoDefaultVpc(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_specification": "MSE_SC_1_2_200_c",
					"cluster_type":          "Eureka",
					"cluster_version":       "EUREKA_1_9_3",
					"instance_count":        "1",
					"net_type":              "privatenet",
					"vswitch_id":            "${data.alicloud_vswitches.default.ids.0}",
					"pub_network_flow":      "1",
					"cluster_alias_name":    "tf-mse",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_specification": "MSE_SC_1_2_200_c",
						"cluster_type":          "Eureka",
						"cluster_version":       "EUREKA_1_9_3",
						"instance_count":        "1",
						"net_type":              "privatenet",
						"vswitch_id":            CHECKSET,
						"pub_network_flow":      "1",
						"cluster_alias_name":    "tf-mse",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_specification", "cluster_version", "net_type", "vswitch_id", "cluster_alias_name"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_entry_list": []string{"127.0.0.1/32"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_entry_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_alias_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_alias_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_alias_name": name + "update",
					"acl_entry_list":     []string{"127.0.0.1/10"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_alias_name": name + "update",
						"acl_entry_list.#":   "1",
					}),
				),
			},
		},
	})
}

var MseClusterMap = map[string]string{}

func MseClusterBasicdependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_vpcs" "default" {
	  is_default = true
	}
	data "alicloud_vswitches" "default" {
	  vpc_id = data.alicloud_vpcs.default.ids.0
	}
`)
}
