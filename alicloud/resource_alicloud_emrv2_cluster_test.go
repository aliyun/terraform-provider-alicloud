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
							"node_resize_strategy": "PRIORITY",
							"auto_scaling_policy": []map[string]interface{}{
								{
									"constraints": []map[string]interface{}{
										{
											"max_capacity": "999",
											"min_capacity": "1",
										},
									},
									"scaling_rules": []map[string]interface{}{
										{
											"rule_name":            "scalingRule01",
											"trigger_type":         "METRICS_TRIGGER",
											"activity_type":        "SCALE_OUT",
											"adjustment_type":      "CHANGE_IN_CAPACITY",
											"adjustment_value":     "1",
											"min_adjustment_value": "1",
											"metrics_trigger": []map[string]interface{}{
												{
													"time_window":              "120",
													"evaluation_count":         "1",
													"cool_down_interval":       "120",
													"condition_logic_operator": "And",
													"time_constraints": []map[string]interface{}{
														{
															"start_time": "00:00",
															"end_time":   "23:59",
														},
													},
													"conditions": []map[string]interface{}{
														{
															"metric_name":         "yarn_resourcemanager_queue_AvailableMBPercentage",
															"statistics":          "AVG",
															"comparison_operator": "LE",
															"threshold":           "10",
															"tags": []map[string]interface{}{
																{
																	"key":   "app",
																	"value": "emr",
																},
															},
														},
													},
												},
											},
										},
										{
											"rule_name":            "scalingRule02",
											"trigger_type":         "TIME_TRIGGER",
											"activity_type":        "SCALE_OUT",
											"adjustment_type":      "CHANGE_IN_CAPACITY",
											"adjustment_value":     "1",
											"min_adjustment_value": "1",
											"time_trigger": []map[string]interface{}{
												{
													"launch_time":            "16:00:00",
													"start_time":             "1745739800000",
													"end_time":               "1745744400000",
													"launch_expiration_time": "3600",
													"recurrence_type":        "DAILY",
													"recurrence_value":       "3",
												},
											},
										},
									},
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
					"bootstrap_scripts": []map[string]interface{}{
						{
							"script_name":             "bssName01",
							"script_path":             "oss://emr/tf-test/ts.sh",
							"script_args":             "--a=b",
							"execution_moment":        "AFTER_STARTED",
							"execution_fail_strategy": "FAILED_CONTINUE",
							"node_selector": []map[string]interface{}{
								{
									"node_select_type": "CLUSTER",
								},
							},
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
				Config: testAccConfig(map[string]interface{}{
					"bootstrap_scripts": []map[string]interface{}{
						{
							"script_name":             "bssName02",
							"script_path":             "oss://emr/tf-test/ts2.sh",
							"script_args":             "--b=a",
							"execution_moment":        "BEFORE_INSTALL",
							"execution_fail_strategy": "FAILED_CONTINUE",
							"node_selector": []map[string]interface{}{
								{
									"node_select_type": "CLUSTER",
								},
							},
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bootstrap_scripts.#":             "1",
						"bootstrap_scripts.0.script_name": "bssName02",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bootstrap_scripts": []map[string]interface{}{
						{
							"script_name":             "bssName02",
							"script_path":             "oss://emr/tf-test/ts2.sh",
							"script_args":             "--b=a",
							"execution_moment":        "BEFORE_INSTALL",
							"execution_fail_strategy": "FAILED_CONTINUE",
							"node_selector": []map[string]interface{}{
								{
									"node_select_type": "NODE_GROUP",
									"node_group_names": []string{"emr-core"},
									"node_group_name":  "emr-core",
								},
							},
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bootstrap_scripts.#":                  "1",
						"bootstrap_scripts.0.script_name":      "bssName02",
						"bootstrap_scripts.0.execution_moment": "BEFORE_INSTALL",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bootstrap_scripts": []map[string]interface{}{
						{
							"script_name":             "bssName02",
							"script_path":             "oss://emr/tf-test/ts2.sh",
							"script_args":             "--b=a",
							"execution_moment":        "BEFORE_INSTALL",
							"execution_fail_strategy": "FAILED_CONTINUE",
							"node_selector": []map[string]interface{}{
								{
									"node_select_type": "NODE_GROUP",
									"node_group_types": []string{"CORE"},
								},
							},
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bootstrap_scripts.#":                  "1",
						"bootstrap_scripts.0.script_name":      "bssName02",
						"bootstrap_scripts.0.execution_moment": "BEFORE_INSTALL",
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
						"node_groups.#": "4",
						"force_sleep":   "60",
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

func TestAccAliCloudEmrV2Cluster_basic3(t *testing.T) {
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
						{
							"node_group_type":         "TASK",
							"node_group_name":         "emr-task",
							"payment_type":            "PayAsYouGo",
							"vswitch_ids":             []string{"${alicloud_vswitch.default.id}"},
							"instance_types":          []string{"ecs.g7.xlarge"},
							"node_count":              "1",
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
						"node_groups.#":     "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "Subscription",
					"subscription_config": []map[string]interface{}{
						{
							"payment_duration_unit":    "Month",
							"payment_duration":         "1",
							"auto_renew":               "true",
							"auto_pay_order":           "true",
							"auto_renew_duration_unit": "Month",
							"auto_renew_duration":      "1",
						},
					},
					"node_groups": []map[string]interface{}{
						{
							"node_group_type":      "MASTER",
							"node_group_name":      "emr-master",
							"payment_type":         "Subscription",
							"vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
							"instance_types":       []string{"ecs.g7.xlarge"},
							"node_count":           "1",
							"with_public_ip":       "false",
							"graceful_shutdown":    "false",
							"spot_instance_remedy": "false",
							"subscription_config": []map[string]interface{}{
								{
									"payment_duration_unit":    "Month",
									"payment_duration":         "1",
									"auto_renew":               "true",
									"auto_pay_order":           "true",
									"auto_renew_duration_unit": "Month",
									"auto_renew_duration":      "1",
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
									"performance_level": "PL0",
									"count":             "3",
								},
							},
						},
						{
							"node_group_type":         "CORE",
							"node_group_name":         "emr-core",
							"payment_type":            "Subscription",
							"vswitch_ids":             []string{"${alicloud_vswitch.default.id}"},
							"instance_types":          []string{"ecs.g7.xlarge"},
							"node_count":              "2",
							"with_public_ip":          "false",
							"deployment_set_strategy": "CLUSTER",
							"graceful_shutdown":       "false",
							"spot_instance_remedy":    "false",
							"subscription_config": []map[string]interface{}{
								{
									"payment_duration_unit":    "Month",
									"payment_duration":         "1",
									"auto_renew":               "true",
									"auto_pay_order":           "true",
									"auto_renew_duration_unit": "Month",
									"auto_renew_duration":      "1",
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
						{
							"node_group_type":         "TASK",
							"node_group_name":         "emr-task",
							"payment_type":            "PayAsYouGo",
							"vswitch_ids":             []string{"${alicloud_vswitch.default.id}"},
							"instance_types":          []string{"ecs.g7.xlarge"},
							"node_count":              "1",
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
						"payment_type":               "Subscription",
						"node_groups.#":              "3",
						"node_groups.0.payment_type": "Subscription",
						"node_groups.1.payment_type": "Subscription",
						"node_groups.2.payment_type": "PayAsYouGo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "Subscription",
					"subscription_config": []map[string]interface{}{
						{
							"payment_duration_unit":    "Month",
							"payment_duration":         "1",
							"auto_renew":               "true",
							"auto_pay_order":           "true",
							"auto_renew_duration_unit": "Month",
							"auto_renew_duration":      "1",
						},
					},
					"node_groups": []map[string]interface{}{
						{
							"node_group_type":      "MASTER",
							"node_group_name":      "emr-master",
							"payment_type":         "Subscription",
							"vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
							"instance_types":       []string{"ecs.g7.xlarge"},
							"node_count":           "1",
							"with_public_ip":       "false",
							"graceful_shutdown":    "false",
							"spot_instance_remedy": "false",
							"subscription_config": []map[string]interface{}{
								{
									"payment_duration_unit":    "Month",
									"payment_duration":         "1",
									"auto_renew":               "true",
									"auto_pay_order":           "true",
									"auto_renew_duration_unit": "Month",
									"auto_renew_duration":      "1",
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
									"performance_level": "PL0",
									"count":             "3",
								},
							},
						},
						{
							"node_group_type":         "CORE",
							"node_group_name":         "emr-core",
							"payment_type":            "Subscription",
							"vswitch_ids":             []string{"${alicloud_vswitch.default.id}"},
							"instance_types":          []string{"ecs.g7.xlarge"},
							"node_count":              "2",
							"with_public_ip":          "false",
							"deployment_set_strategy": "CLUSTER",
							"graceful_shutdown":       "false",
							"spot_instance_remedy":    "false",
							"subscription_config": []map[string]interface{}{
								{
									"payment_duration_unit":    "Month",
									"payment_duration":         "1",
									"auto_renew":               "true",
									"auto_pay_order":           "true",
									"auto_renew_duration_unit": "Month",
									"auto_renew_duration":      "1",
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
						{
							"node_group_type":         "TASK",
							"node_group_name":         "emr-task",
							"payment_type":            "Subscription",
							"vswitch_ids":             []string{"${alicloud_vswitch.default.id}"},
							"instance_types":          []string{"ecs.g7.xlarge"},
							"node_count":              "1",
							"with_public_ip":          "false",
							"deployment_set_strategy": "CLUSTER",
							"graceful_shutdown":       "false",
							"spot_instance_remedy":    "false",
							"subscription_config": []map[string]interface{}{
								{
									"payment_duration_unit":    "Month",
									"payment_duration":         "1",
									"auto_renew":               "true",
									"auto_pay_order":           "true",
									"auto_renew_duration_unit": "Month",
									"auto_renew_duration":      "1",
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
					"tags": map[string]interface{}{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":               "Subscription",
						"node_groups.#":              "3",
						"node_groups.0.payment_type": "Subscription",
						"node_groups.1.payment_type": "Subscription",
						"node_groups.2.payment_type": "Subscription",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"node_groups", "subscription_config"},
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
