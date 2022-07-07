package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_click_house_db_cluster", &resource.Sweeper{
		Name: "alicloud_click_house_db_cluster",
		F:    testSweepClickhouseDbCLuster,
	})
}

func testSweepClickhouseDbCLuster(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeDBClusters"
	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	request["RegionId"] = region
	var response map[string]interface{}
	conn, err := client.NewClickhouseClient()
	if err != nil {
		return WrapError(err)
	}
	ids := make([]string, 0)
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), nil, request, &runtime)
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
			ids = append(ids, fmt.Sprint(item["DBClusterId"]))
		}
		if len(result) < PageSizeLarge {
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
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Click House DBCluster (%s): %s", id, err)
		}
	}
	return nil
}
func TestAccAlicloudClickHouseDBCluster_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_click_house_db_cluster.default"
	ra := resourceAttrInit(resourceId, AlicloudClickHouseDBClusterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ClickhouseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeClickHouseDbCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sclickhousedbcluster%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudClickHouseDBClusterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ClickHouseSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":      "20.3.10.75",
					"category":                "Basic",
					"db_cluster_class":        "S8",
					"db_cluster_network_type": "vpc",
					"db_node_group_count":     "1",
					"payment_type":            "PayAsYouGo",
					"db_node_storage":         "100",
					"storage_type":            "cloud_essd",
					"vswitch_id":              "${data.alicloud_vswitches.default.vswitches.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version":      "20.3.10.75",
						"category":                "Basic",
						"db_cluster_class":        "S8",
						"db_cluster_network_type": "vpc",
						"db_node_group_count":     "1",
						"payment_type":            "PayAsYouGo",
						"db_node_storage":         "100",
						"storage_type":            "cloud_essd",
						"vswitch_id":              CHECKSET,
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
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"db_cluster_class", "db_node_group_count", "db_cluster_version"},
			},
		},
	})
}

func TestAccAlicloudClickHouseDBCluster_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_click_house_db_cluster.default"
	ra := resourceAttrInit(resourceId, AlicloudClickHouseDBClusterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ClickhouseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeClickHouseDbCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sclickhousedbcluster%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudClickHouseDBClusterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ClickHouseSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":      "20.3.10.75",
					"category":                "HighAvailability",
					"db_cluster_class":        "C16",
					"db_cluster_network_type": "vpc",
					"db_node_group_count":     "1",
					"payment_type":            "PayAsYouGo",
					"db_node_storage":         "500",
					"storage_type":            "cloud_essd_pl2",
					"vswitch_id":              "${data.alicloud_vswitches.default.vswitches.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version":      "20.3.10.75",
						"category":                "HighAvailability",
						"db_cluster_class":        "C16",
						"db_cluster_network_type": "vpc",
						"db_node_group_count":     "1",
						"payment_type":            "PayAsYouGo",
						"db_node_storage":         "500",
						"storage_type":            "cloud_essd_pl2",
						"vswitch_id":              CHECKSET,
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
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"db_cluster_class", "db_node_group_count", "db_cluster_version"},
			},
		},
	})
}

func TestAccAlicloudClickHouseDBCluster_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_click_house_db_cluster.default"
	ra := resourceAttrInit(resourceId, AlicloudClickHouseDBClusterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ClickhouseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeClickHouseDbCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sclickhousedbcluster%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudClickHouseDBClusterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ClickHouseSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":      "20.3.10.75",
					"category":                "Basic",
					"db_cluster_class":        "S8",
					"db_cluster_network_type": "vpc",
					"db_node_group_count":     "1",
					"payment_type":            "PayAsYouGo",
					"db_node_storage":         "100",
					"storage_type":            "cloud_essd",
					"vswitch_id":              "${data.alicloud_vswitches.default.vswitches.0.id}",
					"db_cluster_access_white_list": []map[string]interface{}{
						{
							"db_cluster_ip_array_name": "test1",
							"security_ip_list":         "192.168.0.1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version":             "20.3.10.75",
						"category":                       "Basic",
						"db_cluster_class":               "S8",
						"db_cluster_network_type":        "vpc",
						"db_node_group_count":            "1",
						"payment_type":                   "PayAsYouGo",
						"db_node_storage":                "100",
						"storage_type":                   "cloud_essd",
						"vswitch_id":                     CHECKSET,
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"db_cluster_class", "db_node_group_count", "db_cluster_version"},
			},
		},
	})
}

func TestAccAlicloudClickHouseDBCluster_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_click_house_db_cluster.default"
	ra := resourceAttrInit(resourceId, AlicloudClickHouseDBClusterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ClickhouseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeClickHouseDbCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sclic  khousedbcluster%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudClickHouseDBClusterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
			testAccPreCheckWithRegions(t, true, connectivity.ClickHouseSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":      "20.3.10.75",
					"category":                "Basic",
					"db_cluster_class":        "S8",
					"db_cluster_network_type": "vpc",
					"db_node_group_count":     "1",
					"payment_type":            "Subscription",
					"period":                  "Month",
					"used_time":               "1",
					"db_node_storage":         "100",
					"storage_type":            "cloud_essd",
					"vswitch_id":              "${data.alicloud_vswitches.default.vswitches.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version":      "20.3.10.75",
						"category":                "Basic",
						"db_cluster_class":        "S8",
						"db_cluster_network_type": "vpc",
						"db_node_group_count":     "1",
						"payment_type":            "Subscription",
						"period":                  "Month",
						"used_time":               "1",
						"db_node_storage":         "100",
						"storage_type":            "cloud_essd",
						"vswitch_id":              CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"db_cluster_class", "db_node_group_count", "db_cluster_version", "period", "used_time"},
			},
		},
	})
}

var AlicloudClickHouseDBClusterMap0 = map[string]string{
	"db_cluster_version":      CHECKSET,
	"category":                CHECKSET,
	"db_cluster_class":        CHECKSET,
	"db_cluster_network_type": CHECKSET,
	"db_node_group_count":     CHECKSET,
	"payment_type":            CHECKSET,
	"db_node_storage":         CHECKSET,
	"storage_type":            CHECKSET,
}

func AlicloudClickHouseDBClusterBasicDependence0(name string) string {
	return fmt.Sprintf(`
data "alicloud_click_house_regions" "default" {	
  current = true
}

data "alicloud_vpcs" "default"	{
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
  zone_id = data.alicloud_click_house_regions.default.regions.0.zone_ids.0.zone_id
}
variable "name" {
  default = "%s"
}
`, name)
}
