package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_emrv2_cluster", &resource.Sweeper{
		Name: "alicloud_emrv2_cluster",
		F:    testSweepEmrV2Cluster,
	})
}

func testSweepEmrV2Cluster(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	listClustersAction := "ListClusters"
	listClustersRequest := map[string]interface{}{
		"RegionId":      client.RegionId,
		"ClusterStates": []string{"STARTING, START_FAILED", "RUNNING"},
		"NextToken":     "0",
		"MaxResults":    PageSizeMedium,
	}

	for {
		conn, err := client.NewEmrClient()
		if err != nil {
			return WrapError(err)
		}
		listClustersResponse, err := conn.DoRequest(StringPointer(listClustersAction), nil, StringPointer("POST"), StringPointer("2021-03-20"), StringPointer("AK"), nil, listClustersRequest, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_emrv2_cluster", listClustersAction, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Clusters", listClustersResponse)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, listClustersAction, "$.Clusters", listClustersResponse)
		}
		if resp == nil || len(resp.([]interface{})) == 0 {
			break
		}
		for _, item := range resp.([]interface{}) {
			skip := true
			cluster := item.(map[string]interface{})
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(cluster["ClusterName"].(string), prefix) {
						skip = false
					}
				}
				if skip {
					log.Printf("[INFO] Skipping emr: %v (%v)", cluster["ClusterName"], cluster["ClusterId"])
					continue
				}
			}
			deleteClusterAction := "DeleteCluster"
			deleteClusterRequest := map[string]interface{}{
				"RegionId":  client.RegionId,
				"ClusterId": cluster["ClusterId"].(string),
			}
			_, err = conn.DoRequest(StringPointer(deleteClusterAction), nil, StringPointer("POST"), StringPointer("2021-03-20"), StringPointer("AK"), nil, deleteClusterRequest, &util.RuntimeOptions{})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "alicloud_emrv2_cluster", deleteClusterAction, AlibabaCloudSdkGoERROR)
			}
		}
		nextToken := listClustersResponse["NextToken"]
		if nextToken == nil {
			break
		}
		listClustersRequest["NextToken"] = nextToken.(string)
	}
	return nil
}

func TestAccAliCloudEmrV2Cluster_basic(t *testing.T) {
	v := map[string]interface{}{}
	resourceId := "alicloud_emrv2_cluster.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EmrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "GetEmrV2Cluster")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc%sEmrV2ClusterConfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEmrV2ClusterCommonConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"payment_type":      "PayAsYouGo",
					"cluster_type":      "DATAFLOW",
					"release_version":   "EMR-5.10.0",
					"cluster_name":      name,
					"deploy_mode":       "NORMAL",
					"applications":      []string{"HADOOP-COMMON", "HDFS", "YARN"},
					"node_attributes": []map[string]interface{}{
						{
							"vpc_id":            "${alicloud_vpc.default.id}",
							"ram_role":          "${alicloud_ram_role.default.name}",
							"security_group_id": "${alicloud_security_group.default.id}",
							"zone_id":           "${data.alicloud_zones.default.zones.0.id}",
							"key_pair_name":     "${alicloud_ecs_key_pair.default.id}",
						},
					},
					"node_groups": []map[string]interface{}{
						{
							"node_group_type":      "MASTER",
							"node_group_name":      "emr-master",
							"payment_type":         "PayAsYouGo",
							"vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
							"instance_types":       []string{"ecs.g7.xlarge"},
							"node_count":           "1",
							"with_public_ip":       "false",
							"graceful_shutdown":    "false",
							"spot_instance_remedy": "false",
							"system_disk": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"performance_level": "PL0",
									"count":             "1",
								},
							},
							"data_disks": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"count":             "3",
									"performance_level": "PL0",
								},
							},
						},
						{
							"node_group_type":      "CORE",
							"node_group_name":      "emr-core",
							"payment_type":         "PayAsYouGo",
							"vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
							"instance_types":       []string{"ecs.g7.xlarge"},
							"node_count":           "2",
							"with_public_ip":       "false",
							"graceful_shutdown":    "false",
							"spot_instance_remedy": "false",
							"system_disk": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"performance_level": "PL0",
									"count":             "1",
								},
							},
							"data_disks": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"count":             "3",
									"performance_level": "PL0",
								},
							},
						},
					},
					"tags": map[string]interface{}{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name":      name,
						"cluster_type":      "DATAFLOW",
						"payment_type":      "PayAsYouGo",
						"release_version":   "EMR-5.10.0",
						"deploy_mode":       "NORMAL",
						"security_mode":     "NORMAL",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "acceptance test",
						"node_attributes.#": "1",
						"applications.#":    "3",
						"node_groups.#":     "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_name": name + "v2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name": name + "v2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_groups": []map[string]interface{}{
						{
							"node_group_type":      "MASTER",
							"node_group_name":      "emr-master",
							"payment_type":         "PayAsYouGo",
							"vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
							"instance_types":       []string{"ecs.g7.xlarge"},
							"node_count":           "1",
							"with_public_ip":       "false",
							"graceful_shutdown":    "false",
							"spot_instance_remedy": "false",
							"system_disk": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"performance_level": "PL0",
									"count":             "1",
								},
							},
							"data_disks": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"count":             "3",
									"performance_level": "PL0",
								},
							},
						},
						{
							"node_group_type":      "CORE",
							"node_group_name":      "emr-core",
							"payment_type":         "PayAsYouGo",
							"vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
							"instance_types":       []string{"ecs.g7.xlarge"},
							"node_count":           "3",
							"with_public_ip":       "false",
							"graceful_shutdown":    "false",
							"spot_instance_remedy": "false",
							"system_disk": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"performance_level": "PL0",
									"count":             "1",
								},
							},
							"data_disks": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"count":             "3",
									"performance_level": "PL0",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_groups.#": "2",
						"force_sleep":   "240",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_groups": []map[string]interface{}{
						{
							"node_group_type":      "MASTER",
							"node_group_name":      "emr-master",
							"payment_type":         "PayAsYouGo",
							"vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
							"instance_types":       []string{"ecs.g7.xlarge"},
							"node_count":           "1",
							"with_public_ip":       "false",
							"graceful_shutdown":    "false",
							"spot_instance_remedy": "false",
							"system_disk": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"performance_level": "PL0",
									"count":             "1",
								},
							},
							"data_disks": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "90",
									"count":             "3",
									"performance_level": "PL0",
								},
							},
						},
						{
							"node_group_type":      "CORE",
							"node_group_name":      "emr-core",
							"payment_type":         "PayAsYouGo",
							"vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
							"instance_types":       []string{"ecs.g7.xlarge"},
							"node_count":           "3",
							"with_public_ip":       "false",
							"graceful_shutdown":    "false",
							"spot_instance_remedy": "false",
							"system_disk": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"performance_level": "PL0",
									"count":             "1",
								},
							},
							"data_disks": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"count":             "3",
									"performance_level": "PL0",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_groups.#": "2",
						"force_sleep":   "120",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_groups": []map[string]interface{}{
						{
							"node_group_type":      "MASTER",
							"node_group_name":      "emr-master",
							"payment_type":         "PayAsYouGo",
							"vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
							"instance_types":       []string{"ecs.g7.xlarge"},
							"node_count":           "1",
							"with_public_ip":       "false",
							"graceful_shutdown":    "false",
							"spot_instance_remedy": "false",
							"system_disk": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"performance_level": "PL0",
									"count":             "1",
								},
							},
							"data_disks": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "90",
									"count":             "3",
									"performance_level": "PL0",
								},
							},
						},
						{
							"node_group_type":      "CORE",
							"node_group_name":      "emr-core",
							"payment_type":         "PayAsYouGo",
							"vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
							"instance_types":       []string{"ecs.g7.xlarge"},
							"node_count":           "3",
							"with_public_ip":       "false",
							"graceful_shutdown":    "false",
							"spot_instance_remedy": "false",
							"system_disk": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"performance_level": "PL0",
									"count":             "1",
								},
							},
							"data_disks": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"count":             "3",
									"performance_level": "PL0",
								},
							},
						},
						{
							"node_group_type":      "TASK",
							"node_group_name":      "emr-task",
							"payment_type":         "PayAsYouGo",
							"vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
							"instance_types":       []string{"ecs.g7.xlarge"},
							"node_count":           "1",
							"with_public_ip":       "false",
							"graceful_shutdown":    "false",
							"spot_instance_remedy": "false",
							"node_resize_strategy": "COST_OPTIMIZED",
							"cost_optimized_config": []map[string]interface{}{
								{
									"on_demand_base_capacity":                  "1",
									"on_demand_percentage_above_base_capacity": "10",
									"spot_instance_pools":                      "1",
								},
							},
							"spot_bid_prices": []map[string]interface{}{
								{
									"instance_type": "ecs.g7.xlarge",
									"bid_price":     "1",
								},
							},
							"system_disk": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"performance_level": "PL0",
									"count":             "1",
								},
							},
							"data_disks": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"count":             "3",
									"performance_level": "PL0",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_groups.#": "3",
						"force_sleep":   "240",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"node_groups"},
			},
		},
	})
}

func TestAccAliCloudEmrV2Cluster_basic1(t *testing.T) {
	v := map[string]interface{}{}
	resourceId := "alicloud_emrv2_cluster.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EmrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "GetEmrV2Cluster")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc%sEmrV2ClusterConfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEmrV2ClusterCommonConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"payment_type":      "PayAsYouGo",
					"cluster_type":      "DATAFLOW",
					"release_version":   "EMR-5.10.0",
					"cluster_name":      name,
					"deploy_mode":       "NORMAL",
					"applications":      []string{"HADOOP-COMMON", "HDFS", "YARN"},
					"node_attributes": []map[string]interface{}{
						{
							"vpc_id":               "${alicloud_vpc.default.id}",
							"ram_role":             "${alicloud_ram_role.default.name}",
							"security_group_id":    "${alicloud_security_group.default.id}",
							"zone_id":              "${data.alicloud_zones.default.zones.0.id}",
							"key_pair_name":        "${alicloud_ecs_key_pair.default.id}",
							"data_disk_encrypted":  "true",
							"data_disk_kms_key_id": "${data.alicloud_kms_keys.default.ids.0}",
						},
					},
					"node_groups": []map[string]interface{}{
						{
							"node_group_type":      "MASTER",
							"node_group_name":      "emr-master",
							"payment_type":         "PayAsYouGo",
							"vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
							"instance_types":       []string{"ecs.g7.xlarge"},
							"node_count":           "1",
							"with_public_ip":       "false",
							"graceful_shutdown":    "false",
							"spot_instance_remedy": "false",
							"system_disk": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"performance_level": "PL0",
									"count":             "1",
								},
							},
							"data_disks": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"performance_level": "PL0",
									"count":             "3",
								},
							},
						},
						{
							"node_group_type":         "CORE",
							"node_group_name":         "emr-core",
							"payment_type":            "PayAsYouGo",
							"vswitch_ids":             []string{"${alicloud_vswitch.default.id}"},
							"instance_types":          []string{"ecs.g7.xlarge"},
							"node_count":              "2",
							"with_public_ip":          "false",
							"deployment_set_strategy": "CLUSTER",
							"graceful_shutdown":       "false",
							"spot_instance_remedy":    "false",
							"system_disk": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"performance_level": "PL0",
									"count":             "1",
								},
							},
							"data_disks": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"count":             "3",
									"performance_level": "PL0",
								},
							},
						},
					},
					"tags": map[string]interface{}{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name":      name,
						"cluster_type":      "DATAFLOW",
						"payment_type":      "PayAsYouGo",
						"release_version":   "EMR-5.10.0",
						"deploy_mode":       "NORMAL",
						"security_mode":     "NORMAL",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "acceptance test",
						"node_attributes.#": "1",
						"applications.#":    "3",
						"node_groups.#":     "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"node_groups"},
			},
		},
	})
}

func TestAccAliCloudEmrV2Cluster_basic2(t *testing.T) {
	v := map[string]interface{}{}
	resourceId := "alicloud_emrv2_cluster.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EmrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "GetEmrV2Cluster")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc%sEmrV2ClusterConfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEmrV2ClusterCommonConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"payment_type":      "PayAsYouGo",
					"cluster_type":      "DATAFLOW",
					"release_version":   "EMR-5.16.0",
					"cluster_name":      name,
					"deploy_mode":       "NORMAL",
					"security_mode":     "NORMAL",
					"applications":         []string{"HADOOP-COMMON", "HDFS", "YARN"},
					"node_attributes": []map[string]interface{}{
						{
							"vpc_id":            "${alicloud_vpc.default.id}",
							"ram_role":          "${alicloud_ram_role.default.name}",
							"security_group_id": "${alicloud_security_group.default.id}",
							"zone_id":           "${data.alicloud_zones.default.zones.0.id}",
							"key_pair_name":     "${alicloud_ecs_key_pair.default.id}",
						},
					},
					"node_groups": []map[string]interface{}{
						{
							"node_group_type":      "MASTER",
							"node_group_name":      "emr-master",
							"payment_type":         "PayAsYouGo",
							"vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
							"instance_types":       []string{"ecs.g7.xlarge"},
							"node_count":           "1",
							"with_public_ip":       "false",
							"graceful_shutdown":    "false",
							"spot_instance_remedy": "false",
							"system_disk": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"performance_level": "PL0",
									"count":             "1",
								},
							},
							"data_disks": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"count":             "3",
									"performance_level": "PL0",
								},
							},
						},
						{
							"node_group_type":               "CORE",
							"node_group_name":               "emr-core",
							"payment_type":                  "PayAsYouGo",
							"vswitch_ids":                   []string{"${alicloud_vswitch.default.id}"},
							"instance_types":                []string{"ecs.g7.xlarge"},
							"node_count":                    "2",
							"with_public_ip":                "false",
							"graceful_shutdown":             "false",
							"spot_instance_remedy":          "false",
							"system_disk": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"performance_level": "PL0",
									"count":             "1",
								},
							},
							"data_disks": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"count":             "3",
									"performance_level": "PL0",
								},
							},
						},
						{
							"node_group_type":      "GATEWAY",
							"node_group_name":      "emr-gateway",
							"payment_type":         "PayAsYouGo",
							"vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
							"instance_types":       []string{"ecs.g7.xlarge"},
							"node_count":           "1",
							"with_public_ip":       "false",
							"graceful_shutdown":    "false",
							"spot_instance_remedy": "false",
							"system_disk": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"performance_level": "PL0",
									"count":             "1",
								},
							},
							"data_disks": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"count":             "3",
									"performance_level": "PL0",
								},
							},
						},
					},
					"tags": map[string]interface{}{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name":      name,
						"cluster_type":      "DATAFLOW",
						"payment_type":      "PayAsYouGo",
						"release_version":   "EMR-5.16.0",
						"deploy_mode":       "NORMAL",
						"security_mode":     "NORMAL",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "acceptance test",
						"node_attributes.#": "1",
						"applications.#":    "3",
						"node_groups.#":     "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_name": name + "v2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name": name + "v2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_groups": []map[string]interface{}{
						{
							"node_group_type":      "MASTER",
							"node_group_name":      "emr-master",
							"payment_type":         "PayAsYouGo",
							"vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
							"instance_types":       []string{"ecs.g7.xlarge"},
							"node_count":           "1",
							"with_public_ip":       "false",
							"graceful_shutdown":    "false",
							"spot_instance_remedy": "false",
							"system_disk": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"performance_level": "PL0",
									"count":             "1",
								},
							},
							"data_disks": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"count":             "3",
									"performance_level": "PL0",
								},
							},
						},
						{
							"node_group_type":      "CORE",
							"node_group_name":      "emr-core",
							"payment_type":         "PayAsYouGo",
							"vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
							"instance_types":       []string{"ecs.g7.xlarge"},
							"node_count":           "2",
							"with_public_ip":       "false",
							"graceful_shutdown":    "false",
							"spot_instance_remedy": "false",
							"system_disk": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"performance_level": "PL0",
									"count":             "1",
								},
							},
							"data_disks": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"count":             "3",
									"performance_level": "PL0",
								},
							},
						},
						{
							"node_group_type":      "GATEWAY",
							"node_group_name":      "emr-gateway",
							"payment_type":         "PayAsYouGo",
							"vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
							"instance_types":       []string{"ecs.g7.xlarge"},
							"node_count":           "0",
							"with_public_ip":       "false",
							"graceful_shutdown":    "false",
							"spot_instance_remedy": "false",
							"system_disk": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"performance_level": "PL0",
									"count":             "1",
								},
							},
							"data_disks": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"count":             "3",
									"performance_level": "PL0",
								},
							},
						},
						{
							"node_group_type":      "GATEWAY",
							"node_group_name":      "emr-gateway-1",
							"payment_type":         "PayAsYouGo",
							"vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
							"instance_types":       []string{"ecs.g7.xlarge"},
							"node_count":           "0",
							"with_public_ip":       "false",
							"graceful_shutdown":    "false",
							"spot_instance_remedy": "false",
							"system_disk": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"performance_level": "PL0",
									"count":             "1",
								},
							},
							"data_disks": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"count":             "3",
									"performance_level": "PL0",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_groups.#":        "4",
						"force_sleep":          "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_groups": []map[string]interface{}{
						{
							"node_group_type":      "MASTER",
							"node_group_name":      "emr-master",
							"payment_type":         "PayAsYouGo",
							"vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
							"instance_types":       []string{"ecs.g7.xlarge"},
							"node_count":           "1",
							"with_public_ip":       "false",
							"graceful_shutdown":    "false",
							"spot_instance_remedy": "false",
							"system_disk": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"performance_level": "PL0",
									"count":             "1",
								},
							},
							"data_disks": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"count":             "3",
									"performance_level": "PL0",
								},
							},
						},
						{
							"node_group_type":      "CORE",
							"node_group_name":      "emr-core",
							"payment_type":         "PayAsYouGo",
							"vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
							"instance_types":       []string{"ecs.g7.xlarge"},
							"node_count":           "3",
							"with_public_ip":       "false",
							"graceful_shutdown":    "false",
							"spot_instance_remedy": "false",
							"system_disk": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"performance_level": "PL0",
									"count":             "1",
								},
							},
							"data_disks": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"count":             "3",
									"performance_level": "PL0",
								},
							},
						},
						{
							"node_group_type":      "GATEWAY",
							"node_group_name":      "emr-gateway",
							"payment_type":         "PayAsYouGo",
							"vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
							"instance_types":       []string{"ecs.g7.xlarge"},
							"node_count":           "0",
							"with_public_ip":       "false",
							"graceful_shutdown":    "false",
							"spot_instance_remedy": "false",
							"system_disk": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"performance_level": "PL0",
									"count":             "1",
								},
							},
							"data_disks": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"count":             "3",
									"performance_level": "PL0",
								},
							},
						},
						{
							"node_group_type":      "GATEWAY",
							"node_group_name":      "emr-gateway-1",
							"payment_type":         "PayAsYouGo",
							"vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
							"instance_types":       []string{"ecs.g7.xlarge"},
							"node_count":           "1",
							"with_public_ip":       "false",
							"graceful_shutdown":    "false",
							"spot_instance_remedy": "false",
							"system_disk": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"performance_level": "PL0",
									"count":             "1",
								},
							},
							"data_disks": []map[string]interface{}{
								{
									"category":          "cloud_essd",
									"size":              "80",
									"count":             "3",
									"performance_level": "PL0",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_groups.#": "4",
						"force_sleep":   "240",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"node_groups"},
			},
		},
	})
}

func resourceEmrV2ClusterCommonConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	`, EmrV2CommonTestCase, name)
}
