package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_gpdb_instance",
		&resource.Sweeper{
			Name: "alicloud_gpdb_instance",
			F:    testSweepGPDBDBInstance,
		})
}

func testSweepGPDBDBInstance(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeDBInstances"
	request := map[string]interface{}{}
	request["RegionId"] = aliyunClient.RegionId

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = aliyunClient.RpcPost("gpdb", "2016-05-03", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.Items.DBInstance", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Items.DBInstance", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(item["DBInstanceDescription"].(string)), strings.ToLower(prefix)) {
						skip = false
					}
				}
				if skip {
					log.Printf("[INFO] Skipping Gpdb Instance: %s", item["DBInstanceDescription"].(string))
					continue
				}
			}
			action := "DeleteDBInstance"
			request := map[string]interface{}{
				"DBInstanceId": item["DBInstanceId"],
			}
			_, err = aliyunClient.RpcPost("gpdb", "2016-05-03", action, nil, request, true)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Gpdb Instance (%s): %s", item["DBInstanceDescription"].(string), err)
			}
			log.Printf("[INFO] Delete Gpdb Instance success: %s ", item["DBInstanceDescription"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAliCloudGPDBDBInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudGPDBDBInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGPDBDBInstanceBasicDependence0)
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
					"db_instance_category":  "HighAvailability",
					"db_instance_class":     "gpdb.group.segsdx1",
					"db_instance_mode":      "StorageElastic",
					"engine":                "gpdb",
					"engine_version":        "6.0",
					"zone_id":               "${data.alicloud_gpdb_zones.default.ids.0}",
					"instance_network_type": "VPC",
					"instance_spec":         "2C16G",
					"instance_group_count":  "2",
					"payment_type":          "PayAsYouGo",
					"seg_storage_type":      "cloud_essd",
					"seg_node_num":          "4",
					"storage_size":          "50",
					"vpc_id":                "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":            "${local.vswitch_id}",
					"create_sample_data":    "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_category":  "HighAvailability",
						"db_instance_mode":      "StorageElastic",
						"engine":                "gpdb",
						"engine_version":        "6.0",
						"zone_id":               CHECKSET,
						"instance_network_type": "VPC",
						"instance_spec":         "2C16G",
						"payment_type":          "PayAsYouGo",
						"seg_storage_type":      "cloud_essd",
						"seg_node_num":          "4",
						"storage_size":          "50",
						"vpc_id":                CHECKSET,
						"vswitch_id":            CHECKSET,
						"ip_whitelist.#":        "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
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
					"maintain_start_time": "08:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_start_time": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_end_time": "12:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_end_time": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_management_mode": "resourceGroup",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_management_mode": "resourceGroup",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_management_mode": "resourceQueue",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_management_mode": "resourceQueue",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_whitelist": []map[string]interface{}{
						{
							"ip_group_name":    "default",
							"security_ip_list": "10.0.0.1,10.0.0.2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_whitelist.#": "1",
					}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"ip_whitelist": []map[string]interface{}{
			//			{
			//				"ip_group_attribute": "attributedefault",
			//				"security_ip_list":   "10.0.0.3,10.0.0.4",
			//			},
			//		},
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"ip_whitelist.#": "1",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_whitelist": []map[string]interface{}{
						{
							"ip_group_name":    "default",
							"security_ip_list": "10.0.0.1,10.0.0.2",
						},
						{
							"ip_group_attribute": "attribute1",
							"ip_group_name":      "group1",
							"security_ip_list":   "11.0.0.1",
						},
						{
							"ip_group_attribute": "attribute2",
							"ip_group_name":      "group2",
							"security_ip_list":   "12.0.0.1,10.0.0.2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_whitelist.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"seg_node_num": "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"seg_node_num": "8",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_size": "500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_size": "500",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"seg_disk_performance_level": "pl2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"seg_disk_performance_level": "pl2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_spec": "4C32G",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_spec": "4C32G",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_management_mode": "resourceGroup",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_management_mode": "resourceGroup",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_management_mode": "resourceQueue",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_management_mode": "resourceQueue",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vector_configuration_status": "enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vector_configuration_status": "enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vector_configuration_status": "disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vector_configuration_status": "disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"master_cu": "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"master_cu": "8",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enabled": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enabled": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF2",
						"For":     "acceptance test2",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF2",
						"tags.For":     "acceptance test2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parameters": []map[string]interface{}{
						{
							"name":  "optimizer",
							"value": "on",
						},
						{
							"name":  "rds_master_mode",
							"value": "single",
						},
						{
							"name":  "statement_timeout",
							"value": "10800000",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "used_time", "instance_spec", "db_instance_class", "security_ip_list", "instance_group_count", "create_sample_data"},
			},
		},
	})
}

func TestAccAliCloudGPDBDBInstancePrepaid(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudGPDBDBInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGPDBDBInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_category":  "HighAvailability",
					"db_instance_class":     "gpdb.group.segsdx1",
					"db_instance_mode":      "StorageElastic",
					"description":           name,
					"engine":                "gpdb",
					"engine_version":        "6.0",
					"zone_id":               "${data.alicloud_gpdb_zones.default.ids.0}",
					"instance_network_type": "VPC",
					"instance_spec":         "2C16G",
					"payment_type":          "Subscription",
					"seg_storage_type":      "cloud_essd",
					"seg_node_num":          "4",
					"storage_size":          "50",
					"vpc_id":                "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":            "${local.vswitch_id}",
					"period":                "Month",
					"used_time":             "1",
					"create_sample_data":    "false",
					"ip_whitelist": []map[string]interface{}{
						{
							"security_ip_list": "127.0.0.1",
						},
					},
					"parameters": []map[string]interface{}{
						{
							"name":  "optimizer",
							"value": "on",
						},
						{
							"name":  "rds_master_mode",
							"value": "single",
						},
						{
							"name":  "statement_timeout",
							"value": "10800000",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_category":  "HighAvailability",
						"db_instance_mode":      "StorageElastic",
						"description":           name,
						"engine":                "gpdb",
						"engine_version":        "6.0",
						"zone_id":               CHECKSET,
						"instance_network_type": "VPC",
						"instance_spec":         "2C16G",
						"payment_type":          "Subscription",
						"seg_storage_type":      "cloud_essd",
						"seg_node_num":          "4",
						"storage_size":          "50",
						"vpc_id":                CHECKSET,
						"vswitch_id":            CHECKSET,
						"parameters.#":          "3",
						"tags.%":                "2",
						"tags.Created":          "TF",
						"tags.For":              "acceptance test",
						"ip_whitelist.#":        "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "used_time", "instance_spec", "db_instance_class", "security_ip_list", "instance_group_count", "create_sample_data"},
			},
		},
	})
}

func TestAccAliCloudGPDBDBInstanceServerless(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_instance.default"
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"ap-southeast-1"})
	ra := resourceAttrInit(resourceId, AliCloudGPDBDBInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGPDBDBInstanceBasicDependence1)
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
					"db_instance_mode":      "Serverless",
					"description":           name,
					"engine":                "gpdb",
					"engine_version":        "6.0",
					"zone_id":               "ap-southeast-1c",
					"instance_network_type": "VPC",
					"instance_spec":         "4C16G",
					"payment_type":          "PayAsYouGo",
					"seg_node_num":          "2",
					"vpc_id":                "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":            "${local.vswitch_id}",
					"serverless_mode":       "Manual",
					"create_sample_data":    "false",
					"ip_whitelist": []map[string]interface{}{
						{
							"security_ip_list": "127.0.0.1",
						},
					},
					"parameters": []map[string]interface{}{
						{
							"name":  "optimizer",
							"value": "off",
						},
						{
							"name":  "rds_master_mode",
							"value": "single",
						},
						{
							"name":  "statement_timeout",
							"value": "11800000",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_mode":      "Serverless",
						"description":           name,
						"engine":                "gpdb",
						"engine_version":        "6.0",
						"zone_id":               CHECKSET,
						"instance_network_type": "VPC",
						"instance_spec":         "4C16G",
						"payment_type":          "PayAsYouGo",
						"seg_node_num":          "2",
						"vpc_id":                CHECKSET,
						"vswitch_id":            CHECKSET,
						"serverless_mode":       "Manual",
						"ip_whitelist.#":        "1",
						"parameters.#":          "3",
						"tags.%":                "2",
						"tags.Created":          "TF",
						"tags.For":              "acceptance test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_share_status": "opened",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_share_status": "opened",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_share_status": "closed",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_share_status": "closed",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "used_time", "instance_spec", "db_instance_class", "security_ip_list", "instance_group_count", "create_sample_data"},
			},
		},
	})
}

func TestAccAliCloudGPDBDBInstanceServerless_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_instance.default"
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"ap-southeast-1"})
	ra := resourceAttrInit(resourceId, AliCloudGPDBDBInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGPDBDBInstanceBasicDependence1)
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
					"db_instance_mode":      "Serverless",
					"description":           name,
					"engine":                "gpdb",
					"engine_version":        "6.0",
					"zone_id":               "ap-southeast-1c",
					"instance_network_type": "VPC",
					"instance_spec":         "4C16G",
					"payment_type":          "PayAsYouGo",
					"seg_node_num":          "2",
					"vpc_id":                "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":            "${local.vswitch_id}",
					"serverless_mode":       "Manual",
					"data_share_status":     "opened",
					"create_sample_data":    "false",
					"ip_whitelist": []map[string]interface{}{
						{
							"security_ip_list": "127.0.0.1",
						},
					},
					"parameters": []map[string]interface{}{
						{
							"name":  "optimizer",
							"value": "off",
						},
						{
							"name":  "rds_master_mode",
							"value": "single",
						},
						{
							"name":  "statement_timeout",
							"value": "11800000",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_mode":      "Serverless",
						"description":           name,
						"engine":                "gpdb",
						"engine_version":        "6.0",
						"zone_id":               CHECKSET,
						"instance_network_type": "VPC",
						"instance_spec":         "4C16G",
						"payment_type":          "PayAsYouGo",
						"seg_node_num":          "2",
						"vpc_id":                CHECKSET,
						"vswitch_id":            CHECKSET,
						"serverless_mode":       "Manual",
						"data_share_status":     "opened",
						"ip_whitelist.#":        "1",
						"parameters.#":          "3",
						"tags.%":                "2",
						"tags.Created":          "TF",
						"tags.For":              "acceptance test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "used_time", "instance_spec", "db_instance_class", "security_ip_list", "instance_group_count", "create_sample_data"},
			},
		},
	})
}

func TestAccAliCloudGPDBDBInstance_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_instance.default"
	testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	ra := resourceAttrInit(resourceId, AliCloudGPDBDBInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGPDBDBInstanceBasicDependence2)
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
					"db_instance_category":       "HighAvailability",
					"db_instance_mode":           "StorageElastic",
					"description":                name,
					"engine":                     "gpdb",
					"engine_version":             "7.0",
					"availability_zone":          "cn-hangzhou-j",
					"instance_network_type":      "VPC",
					"instance_spec":              "2C16G",
					"master_cu":                  "4",
					"resource_group_id":          "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"instance_charge_type":       "Postpaid",
					"seg_storage_type":           "cloud_essd",
					"seg_disk_performance_level": "pl0",
					"seg_node_num":               "4",
					"storage_size":               "50",
					"vpc_id":                     "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":                 "${local.vswitch_id}",
					"prod_type":                  "cost-effective",
					"create_sample_data":         "false",
					"encryption_type":            "CloudDisk",
					"encryption_key":             "${alicloud_kms_key.key.id}",
					"ssl_enabled":                "1",
					"security_ip_list":           []string{"10.0.0.1,10.0.0.2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_category":       "HighAvailability",
						"db_instance_mode":           "StorageElastic",
						"description":                name,
						"engine":                     "gpdb",
						"engine_version":             "7.0",
						"availability_zone":          CHECKSET,
						"instance_network_type":      "VPC",
						"instance_spec":              "2C16G",
						"master_cu":                  "4",
						"resource_group_id":          CHECKSET,
						"instance_charge_type":       "Postpaid",
						"seg_storage_type":           "cloud_essd",
						"seg_disk_performance_level": "pl0",
						"seg_node_num":               "4",
						"storage_size":               "50",
						"vpc_id":                     CHECKSET,
						"vswitch_id":                 CHECKSET,
						"prod_type":                  "cost-effective",
						"ip_whitelist.#":             "1",
						"encryption_type":            "CloudDisk",
						"encryption_key":             CHECKSET,
						"ssl_enabled":                "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "used_time", "instance_spec", "db_instance_class", "security_ip_list", "instance_group_count", "create_sample_data"},
			},
		},
	})
}

var AliCloudGPDBDBInstanceMap0 = map[string]string{
	"db_instance_category": CHECKSET,
	"resource_group_id":    CHECKSET,
	"prod_type":            CHECKSET,
	"description":          CHECKSET,
	"connection_string":    CHECKSET,
	"port":                 CHECKSET,
	"status":               CHECKSET,
}

func AliCloudGPDBDBInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_gpdb_zones" "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_gpdb_zones.default.ids.0
	}

	resource "alicloud_vswitch" "vswitch" {
  		count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  		vpc_id       = data.alicloud_vpcs.default.ids.0
  		cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  		zone_id      = data.alicloud_gpdb_zones.default.ids.0
  		vswitch_name = var.name
	}

	locals {
  		vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
	}
`, name)
}

func AliCloudGPDBDBInstanceBasicDependence1(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = "ap-southeast-1c"
	}

	resource "alicloud_vswitch" "vswitch" {
  		count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  		vpc_id       = data.alicloud_vpcs.default.ids.0
  		cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  		zone_id      = "ap-southeast-1c"
  		vswitch_name = var.name
	}

	resource "alicloud_kms_key" "key" {
  		pending_window_in_days = "7"
  		key_state              = "Enabled"
  		description            = var.name
	}

	locals {
  		vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
	}
`, name)
}

func AliCloudGPDBDBInstanceBasicDependence2(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = "cn-hangzhou-j"
	}

	resource "alicloud_vswitch" "vswitch" {
  		count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  		vpc_id       = data.alicloud_vpcs.default.ids.0
  		cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  		zone_id      = "cn-hangzhou-j"
  		vswitch_name = var.name
	}

	resource "alicloud_kms_key" "key" {
  		pending_window_in_days = "7"
  		key_state              = "Enabled"
  		description            = var.name
	}

	locals {
  		vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
	}
`, name)
}

func TestUnitAliCloudGPDBDBInstance(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"db_instance_category":  "CreateDBInstanceValue",
		"db_instance_class":     "CreateDBInstanceValue",
		"db_instance_mode":      "CreateDBInstanceValue",
		"description":           "CreateDBInstanceValue",
		"engine":                "CreateDBInstanceValue",
		"engine_version":        "CreateDBInstanceValue",
		"zone_id":               "CreateDBInstanceValue",
		"instance_network_type": "CreateDBInstanceValue",
		"instance_spec":         "CreateDBInstanceValue",
		"payment_type":          "PayAsYouGo",
		"seg_storage_type":      "CreateDBInstanceValue",
		"seg_node_num":          4,
		"storage_size":          50,
		"vpc_id":                "CreateDBInstanceValue",
		"vswitch_id":            "CreateDBInstanceValue",
		"resource_group_id":     "CreateDBInstanceValue",
		"period":                "Month",
		"instance_group_count":  2,
		"create_sample_data":    false,
		"used_time":             "1",
		"ip_whitelist": []map[string]interface{}{
			{
				"security_ip_list": "127.0.0.1",
			},
		},
		"tags": map[string]string{
			"Created": "TF",
			"For":     "acceptance test",
		},
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
		if err != nil {
			log.Printf("[ERROR] the field %s setting error", key)
		}
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// DescribeDBInstanceAttribute
		"Items": map[string]interface{}{
			"DBInstanceAttribute": []interface{}{
				map[string]interface{}{
					"DBInstanceId":          "CreateDBInstanceValue",
					"DBInstanceCategory":    "CreateDBInstanceValue",
					"DBInstanceMode":        "CreateDBInstanceValue",
					"DBInstanceDescription": "CreateDBInstanceValue",
					"Engine":                "CreateDBInstanceValue",
					"EngineVersion":         "CreateDBInstanceValue",
					"InstanceNetworkType":   "CreateDBInstanceValue",
					"MaintainEndTime":       "CreateDBInstanceValue",
					"MaintainStartTime":     "CreateDBInstanceValue",
					"MasterNodeNum":         1,
					"SegmentCounts":         0,
					"PayType":               "Postpaid",
					"SegNodeNum":            4,
					"DBInstanceStatus":      "Running",
					"StorageSize":           50,
					"VSwitchId":             "CreateDBInstanceValue",
					"VpcId":                 "CreateDBInstanceValue",
					"ZoneId":                "CreateDBInstanceValue",
					"TagResources": map[string]interface{}{
						"TagResource": []interface{}{
							map[string]interface{}{
								"TagKey":   "Created",
								"TagValue": "TF",
							},
							map[string]interface{}{
								"TagKey":   "For",
								"TagValue": "acceptance test123",
							},
						},
					},
				},
			},
			"DBInstanceIPArray": []interface{}{
				map[string]interface{}{
					"DBInstanceIPArrayAttribute": "",
					"DBInstanceIPArrayName":      "",
					"SecurityIPList":             "127.0.0.1",
				},
			},
		},
		"DBInstanceId": "CreateDBInstanceValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateDBInstance
		"DBInstanceId": "CreateDBInstanceValue",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_gpdb_instance", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGpdbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudGpdbDbInstanceCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeDBInstanceAttribute Response
		"DBInstanceId": "CreateDBInstanceValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateDBInstance" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGpdbDbInstanceCreate(dInit, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGpdbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudGpdbDbInstanceUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// ModifyDBInstanceDescription
	attributesDiff := map[string]interface{}{
		"description": "ModifyDBInstanceDescriptionValue",
	}
	diff, err := newInstanceDiff("alicloud_gpdb_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDBInstanceAttribute Response
		"Items": map[string]interface{}{
			"DBInstanceAttribute": []interface{}{
				map[string]interface{}{
					"DBInstanceDescription": "ModifyDBInstanceDescriptionValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyDBInstanceDescription" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGpdbDbInstanceUpdate(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// ModifyDBInstanceMaintainTime
	attributesDiff = map[string]interface{}{
		"maintain_end_time":   "ModifyDBInstanceMaintainTimeValue",
		"maintain_start_time": "ModifyDBInstanceMaintainTimeValue",
	}
	diff, err = newInstanceDiff("alicloud_gpdb_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDBInstanceAttribute Response
		"Items": map[string]interface{}{
			"DBInstanceAttribute": []interface{}{
				map[string]interface{}{
					"MaintainEndTime":   "ModifyDBInstanceMaintainTimeValue",
					"MaintainStartTime": "ModifyDBInstanceMaintainTimeValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyDBInstanceMaintainTime" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGpdbDbInstanceUpdate(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// ModifySecurityIps
	attributesDiff = map[string]interface{}{
		"ip_whitelist": []map[string]interface{}{
			{
				"security_ip_list": "1.1.1.1",
			},
		},
	}
	diff, err = newInstanceDiff("alicloud_gpdb_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDBInstanceAttribute Response
		"Items": map[string]interface{}{
			"DBInstanceIPArray": []interface{}{
				map[string]interface{}{
					"DBInstanceIPArrayAttribute": "",
					"DBInstanceIPArrayName":      "",
					"SecurityIPList":             "1.1.1.1",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifySecurityIps" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGpdbDbInstanceUpdate(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// UpgradeDBInstance
	attributesDiff = map[string]interface{}{
		"seg_node_num": 8,
	}
	diff, err = newInstanceDiff("alicloud_gpdb_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDBInstanceAttribute Response
		"Items": map[string]interface{}{
			"DBInstanceAttribute": []interface{}{
				map[string]interface{}{
					"SegNodeNum": 8,
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpgradeDBInstance" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGpdbDbInstanceUpdate(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// UpgradeDBInstance
	attributesDiff = map[string]interface{}{
		"storage_size": 100,
	}
	diff, err = newInstanceDiff("alicloud_gpdb_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDBInstanceAttribute Response
		"Items": map[string]interface{}{
			"DBInstanceAttribute": []interface{}{
				map[string]interface{}{
					"StorageSize": 100,
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpgradeDBInstance" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGpdbDbInstanceUpdate(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeDBInstanceAttribute" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGpdbDbInstanceRead(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGpdbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudGpdbDbInstanceDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "AclNotExist", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteDBInstance" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGpdbDbInstanceDelete(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}
