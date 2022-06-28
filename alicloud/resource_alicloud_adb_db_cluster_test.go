package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_adb_db_instance", &resource.Sweeper{
		Name: "alicloud_adb_db_instance",
		F:    testSweepAdbDbInstances,
	})
}

func testSweepAdbDbInstances(region string) error {
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
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var response map[string]interface{}
	conn, err := client.NewAdsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			log.Println(WrapErrorf(err, DataDefaultErrorMsg, "alicloud_adb_db_clusters", action, AlibabaCloudSdkGoERROR))
			break
		}

		resp, err := jsonpath.Get("$.Items.DBCluster", response)
		if err != nil {
			log.Println(WrapErrorf(err, FailedGetAttributeMsg, action, "$.Items.DBCluster", response))
			break
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := fmt.Sprint(item["DBClusterDescription"])
			id := fmt.Sprint(item["DBClusterId"])
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping ADB Instance: %s (%s)", name, id)
				continue
			}
			log.Printf("[INFO] Deleting adb Instance: %s (%s)", name, id)
			action := "DeleteDBCluster"
			conn, err := client.NewAdsClient()
			if err != nil {
				log.Println(WrapError(err))
				break
			}
			request := map[string]interface{}{
				"DBClusterId": id,
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				log.Printf("[ERROR] Deleting ADB cluster failed with error: %#v", err)
				return nil
			})
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudADBDbCluster_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_adb_db_cluster.default"
	ra := resourceAttrInit(resourceId, AlicloudAdbDbClusterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbDbCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sadbCluster%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAdbDbClusterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.AdbReserverUnSupportRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_category": "Cluster",
					"db_node_class":       "C8",
					"description":         "${var.name}",
					"db_node_count":       "1",
					"db_node_storage":     "100",
					"mode":                "reserver",
					"vswitch_id":          "${local.vswitch_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_category": "Cluster",
						"db_node_class":       "C8",
						"description":         name,
						"db_node_count":       "1",
						"db_node_storage":     "100",
						"mode":                "reserver",
						"vswitch_id":          CHECKSET,
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
					"db_node_class": "C32",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_class": "C32",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_node_count": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_count": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_node_storage": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_storage": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_time": "23:00Z-00:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_time": "23:00Z-00:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_node_class":   "C8",
					"db_node_count":   "1",
					"db_node_storage": "200",
					"description":     name,
					"maintain_time":   "01:00Z-02:00Z",
					"security_ips":    []string{"10.168.1.13"},
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_class":   "C8",
						"db_node_count":   "1",
						"db_node_storage": "200",
						"description":     name,
						"maintain_time":   "01:00Z-02:00Z",
						"security_ips.#":  "1",
						"tags.%":          "2",
						"tags.Created":    "TF-update",
						"tags.For":        "test-update",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudADBDbCluster_flexible8C(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_adb_db_cluster.default"
	ra := resourceAttrInit(resourceId, AlicloudAdbDbClusterMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbDbCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sadbCluster%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAdbDbClusterBasicDependence1)
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
					"db_cluster_category": "MixedStorage",
					"description":         "${var.name}",
					"mode":                "flexible",
					"compute_resource":    "8Core32GB",
					"vswitch_id":          "${local.vswitch_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_category": "MixedStorage",
						"description":         name,
						"mode":                "flexible",
						"compute_resource":    "8Core32GB",
						"vswitch_id":          CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			//// API does not support to updating the compute_resource
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"compute_resource": "16Core64GB",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"compute_resource": "16Core64GB",
			//		}),
			//	),
			//},
			//// API does not support updating elastic_io_resource when compute_resource is 8Core32GB or 16Core64GB
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"elastic_io_resource": "1",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"elastic_io_resource": "1",
			//		}),
			//	),
			//},

			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_time": "23:00Z-00:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_time": "23:00Z-00:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"compute_resource": "8Core32GB",
					// API does not support updating elastic_io_resource when compute_resource is 8Core32GB or 16Core64GB
					//"elastic_io_resource": "1",
					"description":   name,
					"maintain_time": "01:00Z-02:00Z",
					"security_ips":  []string{"10.168.1.13"},
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compute_resource": "8Core32GB",
						//"elastic_io_resource": "1",
						"description":    name,
						"maintain_time":  "01:00Z-02:00Z",
						"security_ips.#": "1",
						"tags.%":         "2",
						"tags.Created":   "TF-update",
						"tags.For":       "test-update",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudADBDbCluster_flexible32C(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_adb_db_cluster.default"
	ra := resourceAttrInit(resourceId, AlicloudAdbDbClusterMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbDbCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sadbCluster%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAdbDbClusterBasicDependence1)
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
					"db_cluster_category": "MixedStorage",
					"description":         "${var.name}",
					"mode":                "flexible",
					"compute_resource":    "32Core128GB",
					"vswitch_id":          "${local.vswitch_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_category": "MixedStorage",
						"description":         name,
						"mode":                "flexible",
						"compute_resource":    "32Core128GB",
						"vswitch_id":          CHECKSET,
						"db_node_class":       "E32",
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
					"compute_resource": "48Core192GB",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compute_resource": "48Core192GB",
						"db_node_count":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"elastic_io_resource": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"elastic_io_resource": "1",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_time": "23:00Z-00:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_time": "23:00Z-00:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"compute_resource":    "64Core256GB",
					"elastic_io_resource": "2",
					"description":         name,
					"maintain_time":       "01:00Z-02:00Z",
					"security_ips":        []string{"10.168.1.13"},
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compute_resource":    "64Core256GB",
						"elastic_io_resource": "2",
						"description":         name,
						"maintain_time":       "01:00Z-02:00Z",
						"security_ips.#":      "1",
						"tags.%":              "2",
						"tags.Created":        "TF-update",
						"tags.For":            "test-update",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudADBDbCluster_modifyPayType(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_adb_db_cluster.default"
	ra := resourceAttrInit(resourceId, AlicloudAdbDbClusterMap2)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbDbCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sadbCluster%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAdbDbClusterBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.AdbReserverUnSupportRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_category": "Cluster",
					"db_node_class":       "C8",
					"description":         "${var.name}",
					"db_node_count":       "1",
					"db_node_storage":     "100",
					"mode":                "reserver",
					"vswitch_id":          "${local.vswitch_id}",
					"payment_type":        "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_category": "Cluster",
						"db_node_class":       "C8",
						"description":         name,
						"db_node_count":       "1",
						"db_node_storage":     "100",
						"mode":                "reserver",
						"vswitch_id":          CHECKSET,
						"payment_type":        "PayAsYouGo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "Subscription",
					"period":       "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "Subscription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew_period", "period", "renewal_status"},
			},
		},
	})
}

var AlicloudAdbDbClusterMap0 = map[string]string{
	"auto_renew_period":   NOSET,
	"compute_resource":    "",
	"connection_string":   CHECKSET,
	"db_cluster_version":  "3.0",
	"db_node_storage":     "0",
	"elastic_io_resource": "0",
	"maintain_time":       CHECKSET,
	"modify_type":         NOSET,
	"payment_type":        "PayAsYouGo",
	"pay_type":            "PostPaid",
	"period":              NOSET,
	"renewal_status":      NOSET,
	"resource_group_id":   CHECKSET,
	"security_ips.#":      "1",
	"status":              "Running",
	"tags.%":              "0",
	"zone_id":             CHECKSET,
}

func AlicloudAdbDbClusterBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "terraformci"
}
%s
`, name, AdbCommonTestCase)
}

var AlicloudAdbDbClusterMap1 = map[string]string{
	"auto_renew_period": NOSET,
	//"compute_resource": "8c16g",
	"connection_string":   CHECKSET,
	"db_cluster_version":  "3.0",
	"db_node_class":       "E8",
	"db_node_count":       "1",
	"db_node_storage":     "100",
	"elastic_io_resource": "0",
	"maintain_time":       CHECKSET,
	"modify_type":         NOSET,
	"payment_type":        "PayAsYouGo",
	"pay_type":            "PostPaid",
	"period":              NOSET,
	"renewal_status":      NOSET,
	"resource_group_id":   CHECKSET,
	"security_ips.#":      "1",
	"status":              "Running",
	"tags.%":              "0",
	"zone_id":             CHECKSET,
}
var AlicloudAdbDbClusterMap2 = map[string]string{}

func AlicloudAdbDbClusterBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}
%s
`, name, AdbCommonTestCase)
}
