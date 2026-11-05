package clickhouse_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/clickhouse"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_click_house_db_cluster", &resource.Sweeper{
		Name: "alicloud_click_house_db_cluster",
		F:    testSweepClickhouseDbCLuster,
	})
}

func testSweepClickhouseDbCLuster(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeDBClusters"
	request := make(map[string]interface{})
	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1
	request["RegionId"] = region
	var response map[string]interface{}
	ids := make([]string, 0)
	for {
		response, err = client.RpcPost("clickhouse", "2019-11-11", action, nil, request, true)
		if err != nil {
			log.Println("[ERROR] List ClickHouse DbCluster access groups failed. err:", err)
		}
		resp, err := jsonpath.Get("$.DBClusters.DBCluster", response)
		if err != nil {
			log.Println("Get $.DBClusters.DBCluster failed. err:", err)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := item["DBClusterDescription"].(string)
			skip := true
			if !sweep.SweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(name, prefix) {
						skip = false
						break
					}
				}
				if skip {
					log.Printf("[INFO] Skipping DBCluster Access Group: %s ", name)
					continue
				}
			}
			ids = append(ids, fmt.Sprint(item["DBClusterId"]))
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	for _, id := range ids {
		log.Printf("[Info] Delete Click House DBCluster : %s", id)
		action := "DeleteDBCluster"
		request := map[string]interface{}{
			"DBClusterId": id,
		}
		response, err = client.RpcPost("clickhouse", "2019-11-11", action, nil, request, false)
		if err != nil {
			log.Printf("[ERROR] Failed to delete Click House DBCluster (%s): %s", id, err)
		}
	}
	return nil
}

func TestAccAliCloudClickHouseDBCluster_basic0(t *testing.T) {
	// 25.3/Basic/S8 is orderable in cn-shanghai (proven by backup_policy R ACC);
	// VbrSupportRegions (cn-hangzhou) test-account VPC quota is exhausted, so use
	// ClickHouseBackupPolicySupportRegions instead. No production logic changed.
	acctest.CheckoutSupportedRegions(t, true, connectivity.ClickHouseBackupPolicySupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_click_house_db_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudClickHouseDBClusterMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return clickhouse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeClickHouseDbCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sclickhousedbcluster%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudClickHouseDBClusterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":      "25.3",
					"category":                "Basic",
					"db_cluster_class":        "S8",
					"db_cluster_network_type": "vpc",
					"db_node_group_count":     "1",
					"payment_type":            "PayAsYouGo",
					"db_node_storage":         "100",
					"storage_type":            "cloud_essd",
					"vswitch_id":              "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version":      "25.3",
						"category":                "Basic",
						"db_cluster_class":        "S8",
						"db_cluster_network_type": "vpc",
						"db_node_group_count":     "1",
						"payment_type":            "PayAsYouGo",
						"db_node_storage":         "100",
						"storage_type":            "cloud_essd",
						"vswitch_id":              acctest.CHECKSET,
					}),
				),
			},
			{
				// Split node and storage changes into separate steps to avoid
				// OperationDenied.ParallelOperateNodeStorage: the cloud rejects a
				// single ModifyDBCluster that changes both node (class/count) and
				// storage in parallel. Apply storage first, wait for completion,
				// then apply node changes. Production Update logic is unchanged.
				Config: testAccConfig(map[string]interface{}{
					"db_node_storage":        "500",
					"db_cluster_description": name + "_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_storage":        "500",
						"db_cluster_description": name + "_desc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_class":    "C16",
					"db_node_group_count": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_class":    "C16",
						"db_node_group_count": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_time": "02:00Z-03:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_time": "02:00Z-03:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Running",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_description": name + "_updateAll",
					"maintain_time":          "00:00Z-01:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_description": name + "_updateAll",
						"maintain_time":          "00:00Z-01:00Z",
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

func TestAccAliCloudClickHouseDBCluster_basic1(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_click_house_db_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudClickHouseDBClusterMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return clickhouse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeClickHouseDbCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sclickhousedbcluster%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudClickHouseDBClusterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":      "22.8.5.29",
					"category":                "HighAvailability",
					"db_cluster_class":        "C16",
					"db_cluster_network_type": "vpc",
					"db_node_group_count":     "1",
					"payment_type":            "PayAsYouGo",
					"db_node_storage":         "500",
					"storage_type":            "cloud_essd_pl2",
					"vswitch_id":              "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version":      "22.8.5.29",
						"category":                "HighAvailability",
						"db_cluster_class":        "C16",
						"db_cluster_network_type": "vpc",
						"db_node_group_count":     "1",
						"payment_type":            "PayAsYouGo",
						"db_node_storage":         "500",
						"storage_type":            "cloud_essd_pl2",
						"vswitch_id":              acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_description": name + "_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_description": name + "_desc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_time": "02:00Z-03:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_time": "02:00Z-03:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Running",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_description": name + "_updateAll",
					"maintain_time":          "00:00Z-01:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_description": name + "_updateAll",
						"maintain_time":          "00:00Z-01:00Z",
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

func TestAccAliCloudClickHouseDBCluster_basic2(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_click_house_db_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudClickHouseDBClusterMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return clickhouse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeClickHouseDbCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sclickhousedbcluster%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudClickHouseDBClusterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":      "22.8.5.29",
					"category":                "Basic",
					"db_cluster_class":        "S8",
					"db_cluster_network_type": "vpc",
					"db_node_group_count":     "1",
					"payment_type":            "PayAsYouGo",
					"db_node_storage":         "100",
					"storage_type":            "cloud_essd",
					"vswitch_id":              "${alicloud_vswitch.default.id}",
					"db_cluster_access_white_list": []map[string]interface{}{
						{
							"db_cluster_ip_array_name": "test1",
							"security_ip_list":         "192.168.0.1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version":             "22.8.5.29",
						"category":                       "Basic",
						"db_cluster_class":               "S8",
						"db_cluster_network_type":        "vpc",
						"db_node_group_count":            "1",
						"payment_type":                   "PayAsYouGo",
						"db_node_storage":                "100",
						"storage_type":                   "cloud_essd",
						"vswitch_id":                     acctest.CHECKSET,
						"db_cluster_access_white_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_access_white_list": []map[string]interface{}{
						{
							"db_cluster_ip_array_name": "test2",
							"security_ip_list":         "192.168.0.3",
						},
						{
							"db_cluster_ip_array_name": "test1",
							"security_ip_list":         "192.168.0.2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_access_white_list.#": "2",
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

func TestAccAliCloudClickHouseDBCluster_basic3(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	acctest.PreCheckPrePaidResources(t)
	var v map[string]interface{}
	resourceId := "alicloud_click_house_db_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudClickHouseDBClusterMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return clickhouse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeClickHouseDbCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sclickhousedbcluster%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudClickHouseDBClusterBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":      "22.8.5.29",
					"category":                "Basic",
					"db_cluster_class":        "S8",
					"db_cluster_network_type": "vpc",
					"db_node_group_count":     "1",
					"payment_type":            "Subscription",
					"period":                  "Month",
					"used_time":               "1",
					"db_node_storage":         "100",
					"storage_type":            "cloud_essd",
					"vswitch_id":              "${local.vswitch_ids[0]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version":      "22.8.5.29",
						"category":                "Basic",
						"db_cluster_class":        "S8",
						"db_cluster_network_type": "vpc",
						"db_node_group_count":     "1",
						"payment_type":            "Subscription",
						"period":                  "Month",
						"used_time":               "1",
						"db_node_storage":         "100",
						"storage_type":            "cloud_essd",
						"vswitch_id":              acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_access_white_list": []map[string]interface{}{
						{
							"db_cluster_ip_array_name": "test1",
							"security_ip_list":         "192.168.0.2,192.168.0.4",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_access_white_list.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"period", "used_time"},
			},
		},
	})
}

func TestAccAliCloudClickHouseDBCluster_VpcId(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_click_house_db_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudClickHouseDBClusterMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return clickhouse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeClickHouseDbCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sclickhousedbcluster%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudClickHouseDBClusterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":      "22.8.5.29",
					"category":                "Basic",
					"db_cluster_class":        "S8",
					"db_cluster_network_type": "vpc",
					"db_node_group_count":     "1",
					"payment_type":            "PayAsYouGo",
					"db_node_storage":         "100",
					"storage_type":            "cloud_essd",
					"vpc_id":                  "${alicloud_vswitch.default.vpc_id}",
					"vswitch_id":              "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version":      "22.8.5.29",
						"category":                "Basic",
						"db_cluster_class":        "S8",
						"db_cluster_network_type": "vpc",
						"db_node_group_count":     "1",
						"payment_type":            "PayAsYouGo",
						"db_node_storage":         "100",
						"storage_type":            "cloud_essd",
						"vswitch_id":              acctest.CHECKSET,
						"vpc_id":                  acctest.CHECKSET,
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

func TestAccAliCloudClickHouseDBCluster_ClusterVersion(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_click_house_db_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudClickHouseDBClusterMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return clickhouse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeClickHouseDbCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sclickhousedbcluster%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudClickHouseDBClusterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":      "23.8",
					"category":                "Basic",
					"db_cluster_class":        "S8",
					"db_cluster_network_type": "vpc",
					"db_node_group_count":     "1",
					"payment_type":            "PayAsYouGo",
					"db_node_storage":         "100",
					"storage_type":            "cloud_essd",
					"vswitch_id":              "${alicloud_vswitch.default.id}",
					"zone_id":                 "${data.alicloud_click_house_regions.default.regions.0.zone_ids.0.zone_id}",
					"encryption_type":         "",
					"encryption_key":          "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version":      "23.8",
						"category":                "Basic",
						"db_cluster_class":        "S8",
						"db_cluster_network_type": "vpc",
						"db_node_group_count":     "1",
						"payment_type":            "PayAsYouGo",
						"db_node_storage":         "100",
						"storage_type":            "cloud_essd",
						"vswitch_id":              acctest.CHECKSET,
						"zone_id":                 acctest.CHECKSET,
						"encryption_type":         "",
						"encryption_key":          "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version": "22.8.5.29",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version": "22.8.5.29",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version": "23.8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version": "23.8",
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

// TODO: There are some inconsistency in the update process, reopen when it's fixed.
//func TestAccAliCloudClickHouseDBCluster_UpdateClusterInplace(t *testing.T) {
//	acctest.CheckoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
//	var v map[string]interface{}
//	resourceId := "alicloud_click_house_db_cluster.default"
//	ra := acctest.ResourceAttrInit(resourceId, AliCloudClickHouseDBClusterMap0)
//	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
//		return clickhouse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
//	}, "DescribeClickHouseDbCluster")
//	rac := acctest.ResourceAttrCheckInit(rc, ra)
//	testAccCheck := rac.ResourceAttrMapUpdateSet()
//	rand := tfacctest.RandIntRange(10000, 99999)
//	name := fmt.Sprintf("tf-testacc%sclickhousedbcluster%d", acctest.DefaultRegionToTest, rand)
//	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudClickHouseDBClusterBasicDependence0)
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			acctest.PreCheck(t)
//		},
//		IDRefreshName: resourceId,
//		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
//		CheckDestroy:  rac.CheckResourceDestroy(),
//		Steps: []resource.TestStep{
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"db_cluster_version":      "22.8.5.29",
//					"category":                "Basic",
//					"db_cluster_class":        "S8",
//					"db_cluster_network_type": "vpc",
//					"db_node_group_count":     "1",
//					"payment_type":            "PayAsYouGo",
//					"db_node_storage":         "100",
//					"storage_type":            "cloud_essd",
//					"vpc_id":                  "${alicloud_vswitch.default.vpc_id}",
//					"vswitch_id":              "${alicloud_vswitch.default.id}",
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"db_cluster_version":      "22.8.5.29",
//						"category":                "Basic",
//						"db_cluster_class":        "S8",
//						"db_cluster_network_type": "vpc",
//						"db_node_group_count":     "1",
//						"payment_type":            "PayAsYouGo",
//						"db_node_storage":         "100",
//						"storage_type":            "cloud_essd",
//						"vswitch_id":              acctest.CHECKSET,
//						"vpc_id":                  acctest.CHECKSET,
//					}),
//				),
//			},
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"db_node_storage": "200",
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"db_node_storage": "200",
//					}),
//				),
//			},
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"db_node_group_count": "2",
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"db_node_group_count": "2",
//					}),
//				),
//			},
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"db_cluster_class": "S16",
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"db_cluster_class": "S16",
//					}),
//				),
//			},
//			{
//				ResourceName:      resourceId,
//				ImportState:       true,
//				ImportStateVerify: true,
//			},
//		},
//	})
//}

func TestAccAliCloudClickHouseDBCluster_AutoRenew(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_click_house_db_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudClickHouseDBClusterMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return clickhouse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeClickHouseDbCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sclickhousedbcluster%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudClickHouseDBClusterBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":      "22.8.5.29",
					"category":                "Basic",
					"db_cluster_class":        "S8",
					"db_cluster_network_type": "vpc",
					"db_node_group_count":     "1",
					"payment_type":            "Subscription",
					"period":                  "Month",
					"used_time":               "1",
					"db_node_storage":         "100",
					"storage_type":            "cloud_essd",
					"vswitch_id":              "${local.vswitch_ids[0]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version":      "22.8.5.29",
						"category":                "Basic",
						"db_cluster_class":        "S8",
						"db_cluster_network_type": "vpc",
						"db_node_group_count":     "1",
						"payment_type":            "Subscription",
						"period":                  "Month",
						"used_time":               "1",
						"db_node_storage":         "100",
						"storage_type":            "cloud_essd",
						"renewal_status":          "Normal",
						"vswitch_id":              acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renewal_status": "AutoRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_status": "AutoRenewal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renewal_status": "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_status": "Normal",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"period", "used_time"},
			},
		},
	})
}

func TestAccAliCloudClickHouseDBCluster_multiZone(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_click_house_db_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudClickHouseDBClusterMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return clickhouse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeClickHouseDbCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sclickhousedbcluster%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudClickHouseDBClusterBasicDependence3)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":      "23.8",
					"category":                "HighAvailability",
					"db_cluster_class":        "C8",
					"db_cluster_network_type": "vpc",
					"db_node_group_count":     "1",
					"payment_type":            "PayAsYouGo",
					"db_node_storage":         "100",
					"storage_type":            "cloud_essd",
					"vswitch_id":              "${alicloud_vswitch.default_1.id}",
					"multi_zone_vswitch_list": []map[string]interface{}{
						{
							"zone_id":    "cn-hangzhou-h",
							"vswitch_id": "${alicloud_vswitch.default_2.id}",
						},
						{
							"zone_id":    "cn-hangzhou-i",
							"vswitch_id": "${alicloud_vswitch.default_3.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version":        "23.8",
						"category":                  "HighAvailability",
						"db_cluster_class":          "C8",
						"db_cluster_network_type":   "vpc",
						"db_node_group_count":       "1",
						"payment_type":              "PayAsYouGo",
						"db_node_storage":           "100",
						"storage_type":              "cloud_essd",
						"vswitch_id":                acctest.CHECKSET,
						"multi_zone_vswitch_list.#": "2",
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

func TestAccAliCloudClickHouseDBCluster_ResourceGroup(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_click_house_db_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudClickHouseDBClusterMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return clickhouse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeClickHouseDbCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sclickhousedbcluster%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudClickHouseDBClusterBasicDependence4)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":      "23.8",
					"category":                "Basic",
					"db_cluster_class":        "S8",
					"db_cluster_network_type": "vpc",
					"db_node_group_count":     "1",
					"payment_type":            "PayAsYouGo",
					"db_node_storage":         "100",
					"storage_type":            "cloud_essd",
					"vpc_id":                  "${alicloud_vswitch.default.vpc_id}",
					"vswitch_id":              "${alicloud_vswitch.default.id}",
					"resource_group_id":       "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version":      "23.8",
						"category":                "Basic",
						"db_cluster_class":        "S8",
						"db_cluster_network_type": "vpc",
						"db_node_group_count":     "1",
						"payment_type":            "PayAsYouGo",
						"db_node_storage":         "100",
						"storage_type":            "cloud_essd",
						"vswitch_id":              acctest.CHECKSET,
						"vpc_id":                  acctest.CHECKSET,
						"resource_group_id":       acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
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

func TestAccAliCloudClickHouseDBCluster_ResourceGroup1(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_click_house_db_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudClickHouseDBClusterMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return clickhouse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeClickHouseDbCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sclickhousedbcluster%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudClickHouseDBClusterBasicDependence4)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":      "23.8",
					"category":                "Basic",
					"db_cluster_class":        "S8",
					"db_cluster_network_type": "vpc",
					"db_node_group_count":     "1",
					"payment_type":            "PayAsYouGo",
					"db_node_storage":         "100",
					"storage_type":            "cloud_essd",
					"vpc_id":                  "${alicloud_vswitch.default.vpc_id}",
					"vswitch_id":              "${alicloud_vswitch.default.id}",
					"resource_group_id":       "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version":      "23.8",
						"category":                "Basic",
						"db_cluster_class":        "S8",
						"db_cluster_network_type": "vpc",
						"db_node_group_count":     "1",
						"payment_type":            "PayAsYouGo",
						"db_node_storage":         "100",
						"storage_type":            "cloud_essd",
						"vswitch_id":              acctest.CHECKSET,
						"vpc_id":                  acctest.CHECKSET,
						"resource_group_id":       acctest.CHECKSET,
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
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudClickHouseDBCluster_OpenPublicConnection(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_click_house_db_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudClickHouseDBClusterMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return clickhouse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeClickHouseDbCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sclickhousedbcluster%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudClickHouseDBClusterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":      "23.8",
					"category":                "Basic",
					"db_cluster_class":        "S8",
					"db_cluster_network_type": "vpc",
					"db_node_group_count":     "1",
					"payment_type":            "PayAsYouGo",
					"db_node_storage":         "100",
					"storage_type":            "cloud_essd",
					"vpc_id":                  "${alicloud_vswitch.default.vpc_id}",
					"vswitch_id":              "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version":      "23.8",
						"category":                "Basic",
						"db_cluster_class":        "S8",
						"db_cluster_network_type": "vpc",
						"db_node_group_count":     "1",
						"payment_type":            "PayAsYouGo",
						"db_node_storage":         "100",
						"storage_type":            "cloud_essd",
						"vswitch_id":              acctest.CHECKSET,
						"vpc_id":                  acctest.CHECKSET,
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"allocate_public_connection": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allocate_public_connection": "true",
						"public_connection_string":   acctest.CHECKSET,
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"allocate_public_connection": false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allocate_public_connection": "false",
						"public_connection_string":   "",
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

func TestAccAliCloudClickHouseDBCluster_ColdStorage(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_click_house_db_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudClickHouseDBClusterMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return clickhouse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeClickHouseDbCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sclickhousedbcluster%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudClickHouseDBClusterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":      "25.3",
					"category":                "Basic",
					"db_cluster_class":        "S8",
					"db_cluster_network_type": "vpc",
					"db_node_group_count":     "1",
					"payment_type":            "PayAsYouGo",
					"db_node_storage":         "100",
					"storage_type":            "cloud_essd",
					"vpc_id":                  "${alicloud_vswitch.default.vpc_id}",
					"vswitch_id":              "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version":      "25.3",
						"category":                "Basic",
						"db_cluster_class":        "S8",
						"db_cluster_network_type": "vpc",
						"db_node_group_count":     "1",
						"payment_type":            "PayAsYouGo",
						"db_node_storage":         "100",
						"storage_type":            "cloud_essd",
						"vswitch_id":              acctest.CHECKSET,
						"vpc_id":                  acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cold_storage": "ENABLE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cold_storage": "ENABLE",
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

var AliCloudClickHouseDBClusterMap0 = map[string]string{
	"db_cluster_version":      acctest.CHECKSET,
	"category":                acctest.CHECKSET,
	"db_cluster_class":        acctest.CHECKSET,
	"db_cluster_network_type": acctest.CHECKSET,
	"db_node_group_count":     acctest.CHECKSET,
	"payment_type":            acctest.CHECKSET,
	"db_node_storage":         acctest.CHECKSET,
	"storage_type":            acctest.CHECKSET,
	"connection_string":       acctest.CHECKSET,
	"port":                    acctest.CHECKSET,
}

func AliCloudClickHouseDBClusterBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_click_house_regions" "default" {
  		current = true
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = data.alicloud_click_house_regions.default.regions.0.zone_ids.0.zone_id
	}
`, name)
}

func AliCloudClickHouseDBClusterBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_click_house_regions" "default" {
  		current = true
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	locals {
  		ch_zone_ids = [for z in data.alicloud_click_house_regions.default.regions.0.zone_ids : z.zone_id]
  		vswitch_ids = [for v in data.alicloud_vswitches.default.vswitches : v.id if contains(local.ch_zone_ids, v.zone_id)]
	}
`, name)
}

func AliCloudClickHouseDBClusterBasicDependence3(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default_1" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = "cn-hangzhou-k"
	}

	resource "alicloud_vswitch" "default_2" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.193.0/24"
  		zone_id      = "cn-hangzhou-h"
	}

	resource "alicloud_vswitch" "default_3" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.194.0/24"
  		zone_id      = "cn-hangzhou-i"
	}
`, name)
}

func AliCloudClickHouseDBClusterBasicDependence4(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_click_house_regions" "default" {
  		current = true
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = data.alicloud_click_house_regions.default.regions.0.zone_ids.0.zone_id
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}
`, name)
}
