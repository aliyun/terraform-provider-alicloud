package mse_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/mse"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_mse_cluster", &resource.Sweeper{
		Name: "alicloud_mse_cluster",
		F:    testSweepMSECluster,
	})
}

func testSweepMSECluster(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf-test",
	}
	request := make(map[string]interface{})
	var response map[string]interface{}
	action := "ListClusters"
	request["PageSize"] = helper.PageSizeLarge
	request["PageNum"] = 1
	for {
		response, err = client.RpcGet("mse", "2019-05-31", action, request, nil)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_mse_clusters", action, errs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Data", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Data", response)
		}
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
			action = "DeleteCluster"
			request := map[string]interface{}{
				"InstanceId": item["InstanceId"],
			}
			_, err = client.RpcPost("mse", "2019-05-31", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Mse Clusters (%s (%s)): %s", item["ClusterAliasName"].(string), item["InstanceId"].(string), err)
			}

			// Waiting 30 seconds to ensure these Mse Clusters have been deleted.
			time.Sleep(30 * time.Second)

			log.Printf("[INFO] Delete mse cluster success: %s ", item["InstanceId"].(string))
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNum"] = request["PageNum"].(int) + 1
	}
	return nil
}

func TestAccAliCloudMSECluster_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mse_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, MseClusterMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMseCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccMseCluster%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, MseClusterBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shenzhen"})
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckMseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_specification":     "MSE_SC_1_2_60_c",
					"cluster_type":              "Nacos-Ans",
					"cluster_version":           "NACOS_2_0_0",
					"instance_count":            "1",
					"net_type":                  "privatenet",
					"vswitch_id":                "${alicloud_vswitch.default.id}",
					"pub_network_flow":          "1",
					"cluster_alias_name":        name,
					"connection_type":           "slb",
					"mse_version":               "mse_dev",
					"private_slb_specification": "test",
					"pub_slb_specification":     "test",
					"disk_type":                 "test",
					"request_pars":              "",
					"vpc_id":                    "${alicloud_vpc.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_specification": "MSE_SC_1_2_60_c",
						"cluster_type":          "Nacos-Ans",
						"cluster_version":       "NACOS_2_0_0",
						"instance_count":        "1",
						"net_type":              "privatenet",
						"vswitch_id":            acctest.CHECKSET,
						"pub_network_flow":      "1",
						"cluster_alias_name":    name,
						"connection_type":       "slb",
						"mse_version":           "mse_dev",
					}),
				),
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
					"acl_entry_list":     []string{"127.0.0.0/16"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_alias_name": name,
						"acl_entry_list.#":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_specification": "MSE_SC_2_4_60_c",
					"resource_group_id":     "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_specification": "MSE_SC_2_4_60_c",
						"resource_group_id":     acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_slb_specification": "test_update",
					"pub_slb_specification":     "test_update",
					"disk_type":                 "test_update",
					"request_pars":              "{}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_alias_name": name,
						"vswitch_id":         acctest.CHECKSET,
						"vpc_id":             acctest.CHECKSET,
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl_entry_list", "private_slb_specification", "pub_slb_specification", "disk_type", "request_pars"},
			},
		},
	})
}

func TestAccAliCloudMSECluster_changeNetwork(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mse_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, MseClusterMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMseCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccMseCluster%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, MseClusterBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shenzhen"})
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckMseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_specification": "MSE_SC_1_2_60_c",
					"cluster_type":          "Nacos-Ans",
					"cluster_version":       "NACOS_2_0_0",
					"instance_count":        "1",
					"net_type":              "privatenet",
					"vswitch_id":            "${alicloud_vswitch.default.id}",
					"pub_network_flow":      "1",
					"cluster_alias_name":    name,
					"connection_type":       "slb",
					"mse_version":           "mse_dev",
					"vpc_id":                "${alicloud_vpc.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_specification": "MSE_SC_1_2_60_c",
						"cluster_type":          "Nacos-Ans",
						"cluster_version":       "NACOS_2_0_0",
						"instance_count":        "1",
						"net_type":              "privatenet",
						"vswitch_id":            acctest.CHECKSET,
						"pub_network_flow":      "1",
						"cluster_alias_name":    name,
						"connection_type":       "slb",
						"mse_version":           "mse_dev",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":           "${alicloud_vpc.default2.id}",
					"vswitch_id":       "${alicloud_vswitch.default2.id}",
					"pub_network_flow": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_alias_name": name,
						"vswitch_id":         acctest.CHECKSET,
						"vpc_id":             acctest.CHECKSET,
						"pub_network_flow":   "2",
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

func TestAccAliCloudMSECluster_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mse_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, MseClusterMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMseCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccMseCluster%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, MseClusterBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckMseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_specification": "MSE_SC_1_2_60_c",
					"cluster_type":          "ZooKeeper",
					"cluster_version":       "ZooKeeper_3_8_0",
					"instance_count":        "1",
					"net_type":              "privatenet",
					"vswitch_id":            "${alicloud_vswitch.default.id}",
					"pub_network_flow":      "1",
					"cluster_alias_name":    name,
					"connection_type":       "slb",
					"mse_version":           "mse_dev",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_specification": "MSE_SC_1_2_60_c",
						"cluster_type":          "ZooKeeper",
						"cluster_version":       "ZooKeeper_3_8_0",
						"instance_count":        "1",
						"net_type":              "privatenet",
						"vswitch_id":            acctest.CHECKSET,
						"pub_network_flow":      "1",
						"cluster_alias_name":    name,
						"connection_type":       "slb",
						"mse_version":           "mse_dev",
					}),
				),
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
					"acl_entry_list":     []string{"127.0.0.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_alias_name": name,
						"acl_entry_list.#":   "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl_entry_list"},
			},
		},
	})
}

func TestAccAliCloudMSECluster_pro(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mse_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, MseClusterMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMseCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccMseCluster%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, MseClusterBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckMseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_specification": "MSE_SC_4_8_60_c",
					"cluster_type":          "Nacos-Ans",
					"cluster_version":       "NACOS_2_0_0",
					"instance_count":        "3",
					"net_type":              "privatenet",
					"vswitch_id":            "${alicloud_vswitch.default.id}",
					"pub_network_flow":      "1",
					"cluster_alias_name":    name,
					"mse_version":           "mse_pro",
					"connection_type":       "slb",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_specification": "MSE_SC_4_8_60_c",
						"cluster_type":          "Nacos-Ans",
						"cluster_version":       "NACOS_2_0_0",
						"instance_count":        "3",
						"vswitch_id":            acctest.CHECKSET,
						"pub_network_flow":      "1",
						"cluster_alias_name":    name,
						"mse_version":           "mse_pro",
						"connection_type":       "slb",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_specification": "MSE_SC_8_16_60_c",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_specification": "MSE_SC_8_16_60_c",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_count": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_count": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"version_code": "LATEST",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version_code": acctest.CHECKSET,
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

func TestAccAliCloudMSECluster_serverless(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mse_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, MseClusterMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMseCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccMseCluster%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, MseClusterBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckMseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_specification": "MSE_SC_SERVERLESS",
					"cluster_type":          "Nacos-Ans",
					"cluster_version":       "NACOS_2_0_0",
					"instance_count":        "3",
					"net_type":              "both",
					"vswitch_id":            "${alicloud_vswitch.default.id}",
					"pub_network_flow":      "1",
					"cluster_alias_name":    name,
					"mse_version":           "mse_serverless",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_specification": "MSE_SC_SERVERLESS",
						"cluster_type":          "Nacos-Ans",
						"cluster_version":       "NACOS_2_0_0",
						"instance_count":        "3",
						"vswitch_id":            acctest.CHECKSET,
						"pub_network_flow":      "1",
						"cluster_alias_name":    name,
						"mse_version":           "mse_serverless",
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

func TestAccAliCloudMSECluster_VpcId(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mse_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, MseClusterMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMseCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccMseCluster%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, MseClusterBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckMseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_specification": "MSE_SC_1_2_60_c",
					"cluster_type":          "Nacos-Ans",
					"cluster_version":       "NACOS_2_0_0",
					"instance_count":        "1",
					"net_type":              "privatenet",
					"vswitch_id":            "${alicloud_vswitch.default.id}",
					"pub_network_flow":      "1",
					"cluster_alias_name":    name,
					"connection_type":       "slb",
					"mse_version":           "mse_dev",
					"vpc_id":                "${alicloud_vpc.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_specification": "MSE_SC_1_2_60_c",
						"cluster_type":          "Nacos-Ans",
						"cluster_version":       "NACOS_2_0_0",
						"instance_count":        "1",
						"net_type":              "privatenet",
						"vswitch_id":            acctest.CHECKSET,
						"pub_network_flow":      "1",
						"cluster_alias_name":    name,
						"connection_type":       "slb",
						"mse_version":           "mse_dev",
						"vpc_id":                acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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

func TestAccAliCloudMSECluster_PrePaid(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mse_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, MseClusterMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMseCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccMseCluster%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, MseClusterBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPrePaidResources(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheckWithTime(t, []int{22})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckMseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_specification": "MSE_SC_1_2_60_c",
					"cluster_type":          "Nacos-Ans",
					"cluster_version":       "NACOS_2_0_0",
					"instance_count":        "1",
					"net_type":              "privatenet",
					"payment_type":          "Subscription",
					"vswitch_id":            "${alicloud_vswitch.default.id}",
					"pub_network_flow":      "1",
					"cluster_alias_name":    name,
					"connection_type":       "slb",
					"mse_version":           "mse_dev",
					"vpc_id":                "${alicloud_vpc.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_specification": "MSE_SC_1_2_60_c",
						"cluster_type":          "Nacos-Ans",
						"cluster_version":       "NACOS_2_0_0",
						"instance_count":        "1",
						"net_type":              "privatenet",
						"vswitch_id":            acctest.CHECKSET,
						"pub_network_flow":      "1",
						"cluster_alias_name":    name,
						"connection_type":       "slb",
						"mse_version":           "mse_dev",
						"vpc_id":                acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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

func TestAccAliCloudMSECluster_single_eni(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mse_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, MseClusterMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMseCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccMseCluster%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, MseClusterBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"ap-southeast-1"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckMseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_specification": "MSE_SC_1_2_60_c",
					"cluster_type":          "Nacos-Ans",
					"cluster_version":       "NACOS_2_0_0",
					"instance_count":        "1",
					"net_type":              "privatenet",
					"vswitch_id":            "${alicloud_vswitch.default.id}",
					"pub_network_flow":      "1",
					"cluster_alias_name":    name,
					"connection_type":       "single_eni",
					"mse_version":           "mse_dev",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_specification": "MSE_SC_1_2_60_c",
						"cluster_type":          "Nacos-Ans",
						"cluster_version":       "NACOS_2_0_0",
						"instance_count":        "1",
						"net_type":              "privatenet",
						"vswitch_id":            acctest.CHECKSET,
						"pub_network_flow":      "1",
						"cluster_alias_name":    name,
						"connection_type":       "single_eni",
						"mse_version":           "mse_dev",
					}),
				),
			},
		},
	})
}

// After DeleteCluster the cluster reaches DESTROY_SUCCESS (the API's terminal
// destroyed state, the same target the production Delete stateConf waits for)
// but remains queryable via QueryClusterDetail for an extended window before
// the API eventually returns mse-200-021 (NotFound). The shared
// CheckResourceDestroy performs a single Describe and treats the still-present
// cluster as not destroyed. testAccCheckMseClusterDestroy bounded-retries
// DescribeMseCluster until the cluster is NotFound (mse-200-021) or reaches
// DESTROY_SUCCESS — both are real destroy assertions; it fails if the cluster
// lingers in a non-destroyed state. No sleep-without-check, no error
// swallowing, no CheckDestroy removal.
func testAccCheckMseClusterDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	mseService := mse.NewService(client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_mse_cluster" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}
		wait := helper.IncrementalWait(5*time.Second, 10*time.Second)
		if err := retry.Retry(10*time.Minute, func() *retry.RetryError {
			object, err := mseService.DescribeMseCluster(rs.Primary.ID)
			if err != nil {
				if errs.NotFoundError(err) {
					return nil
				}
				return retry.NonRetryableError(err)
			}
			if fmt.Sprint(object["InitStatus"]) == "DESTROY_SUCCESS" {
				return nil
			}
			wait()
			return retry.RetryableError(fmt.Errorf("mse cluster %s still in status %q after delete", rs.Primary.ID, fmt.Sprint(object["InitStatus"])))
		}); err != nil {
			return fmt.Errorf("mse cluster %s was not destroyed: %s", rs.Primary.ID, err)
		}
	}
	return nil
}

var MseClusterMap = map[string]string{
	"cluster_id":  acctest.CHECKSET,
	"app_version": acctest.CHECKSET,
	"status":      acctest.CHECKSET,
}

func MseClusterBasicdependence(name string) string {
	return fmt.Sprintf(`

	variable "name" {
		 default = "%v"
		}

	data "alicloud_zones" "default" {
	  available_resource_creation = "VSwitch"
	}

	resource "alicloud_vpc" "default" {
	  vpc_name       = "${var.name}"
      cidr_block = "172.17.3.0/24"
	}

	resource "alicloud_vswitch" "default" {
	  vswitch_name = "${var.name}"
	  cidr_block   = "172.17.3.0/24"
	  vpc_id       = alicloud_vpc.default.id
	  zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_vpc" "default2" {
	  vpc_name       = "tftest"
      cidr_block = "172.17.3.0/24"
	}

	resource "alicloud_vswitch" "default2" {
	  vswitch_name = "tftest"
	  cidr_block   = "172.17.3.0/24"
	  vpc_id       = alicloud_vpc.default2.id
	  zone_id      = data.alicloud_zones.default.zones.0.id
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}
`, name)
}
