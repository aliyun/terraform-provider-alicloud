package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
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
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf-test",
	}
	request := make(map[string]interface{})
	var response map[string]interface{}
	action := "ListClusters"
	conn, err := client.NewMseClient()
	if err != nil {
		return WrapError(err)
	}
	request["PageSize"] = PageSizeLarge
	request["PageNum"] = 1
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-05-31"), StringPointer("AK"), request, nil, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_mse_clusters", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Data", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data", response)
		}
		sweeped := false
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["ClusterAliasName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Mse Clusters: %s (%s)", item["ClusterAliasName"], item["InstanceId"])
				continue
			}
			sweeped = true
			action = "DeleteCluster"
			request := map[string]interface{}{
				"InstanceId": item["InstanceId"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-05-31"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Mse Clusters (%s (%s)): %s", item["ClusterAliasName"].(string), item["InstanceId"].(string), err)
			}
			if sweeped {
				// Waiting 30 seconds to ensure these Mse Clusters have been deleted.
				time.Sleep(30 * time.Second)
			}
			log.Printf("[INFO] Delete mse cluster success: %s ", item["InstanceId"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNum"] = request["PageNum"].(int) + 1
	}
	return nil
}

func TestAccAlicloudMSECluster_basic(t *testing.T) {
	var v map[string]interface{}
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
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_specification": "MSE_SC_1_2_200_c",
					"cluster_type":          "Nacos-Ans",
					"cluster_version":       "NACOS_ANS_1_2_1",
					"instance_count":        "1",
					"net_type":              "privatenet",
					"vswitch_id":            "${data.alicloud_vswitches.default.ids.0}",
					"pub_network_flow":      "1",
					"cluster_alias_name":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_specification": "MSE_SC_1_2_200_c",
						"cluster_type":          "Nacos-Ans",
						"cluster_version":       "NACOS_ANS_1_2_1",
						"instance_count":        "1",
						"net_type":              "privatenet",
						"vswitch_id":            CHECKSET,
						"pub_network_flow":      "1",
						"cluster_alias_name":    name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_specification"},
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
					"cluster_alias_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_alias_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_alias_name": name,
					"acl_entry_list":     []string{"127.0.0.0/10"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_alias_name": name,
						"acl_entry_list.#":   "1",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudMSECluster_basic1(t *testing.T) {
	var v map[string]interface{}
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
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_specification": "MSE_SC_1_2_200_c",
					"cluster_type":          "ZooKeeper",
					"cluster_version":       "ZooKeeper_3_4_14",
					"instance_count":        "1",
					"net_type":              "privatenet",
					"vswitch_id":            "${data.alicloud_vswitches.default.ids.0}",
					"pub_network_flow":      "1",
					"cluster_alias_name":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_specification": "MSE_SC_1_2_200_c",
						"cluster_type":          "ZooKeeper",
						"cluster_version":       "ZooKeeper_3_4_14",
						"instance_count":        "1",
						"net_type":              "privatenet",
						"vswitch_id":            CHECKSET,
						"pub_network_flow":      "1",
						"cluster_alias_name":    name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_specification"},
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
					"cluster_alias_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_alias_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_alias_name": name,
					"acl_entry_list":     []string{"127.0.0.0/10"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_alias_name": name,
						"acl_entry_list.#":   "1",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudMSECluster_pro(t *testing.T) {
	var v map[string]interface{}
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
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_specification": "MSE_SC_2_4_200_c",
					"cluster_type":          "Nacos-Ans",
					"cluster_version":       "NACOS_2_0_0",
					"instance_count":        "3",
					"net_type":              "privatenet",
					"vswitch_id":            "${data.alicloud_vswitches.default.ids.0}",
					"pub_network_flow":      "1",
					"cluster_alias_name":    name,
					"mse_version":           "mse_pro",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_specification": "MSE_SC_2_4_200_c",
						"cluster_type":          "Nacos-Ans",
						"cluster_version":       "NACOS_2_0_0",
						"instance_count":        "3",
						"vswitch_id":            CHECKSET,
						"pub_network_flow":      "1",
						"cluster_alias_name":    name,
						"mse_version":           "mse_pro",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_specification"},
			},
		},
	})
}

var MseClusterMap = map[string]string{}

func MseClusterBasicdependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_vpcs" "default" {
	  name_regex = "default-NODELETING"
	}
	data "alicloud_vswitches" "default" {
	  vpc_id = data.alicloud_vpcs.default.ids.0
	}
`)
}
