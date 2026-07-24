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
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
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
				ImportStateVerifyIgnore: []string{"period", "used_time", "db_instance_class", "security_ip_list", "instance_group_count", "create_sample_data", "parameters"},
			},
		},
	})
}

func TestAccAliCloudGPDBDBInstancePrepaid(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_instance.default"
	// Pin to cn-beijing, which has StorageElastic subscription inventory; the
	// data.alicloud_gpdb_zones data source then auto-selects an available zone.
	testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region("cn-beijing")})
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
				ImportStateVerifyIgnore: []string{"period", "used_time", "db_instance_class", "security_ip_list", "instance_group_count", "create_sample_data", "parameters"},
			},
		},
	})
}

// gpdbServerlessTestRegion returns the region for the GPDB serverless acceptance
// tests. It defaults to cn-beijing, which has serverless inventory, and can be
// overridden via the GPDB_SERVERLESS_REGION environment variable so the remote
// ACC run can target another region if the default ever reports
// OperationDenied.InsufficientResourceCapacity.
func gpdbServerlessTestRegion() string {
	if v := os.Getenv("GPDB_SERVERLESS_REGION"); v != "" {
		return v
	}
	return "cn-beijing"
}

// gpdbServerlessTestZone returns the zone for the GPDB serverless acceptance
// tests, paired with gpdbServerlessTestRegion. Defaults to cn-beijing-h and can
// be overridden via the GPDB_SERVERLESS_ZONE environment variable so the instance
// and its vswitch stay in the same zone.
func gpdbServerlessTestZone() string {
	if v := os.Getenv("GPDB_SERVERLESS_ZONE"); v != "" {
		return v
	}
	return "cn-beijing-h"
}

func TestAccAliCloudGPDBDBInstanceServerless(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_instance.default"
	testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region(gpdbServerlessTestRegion())})
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
					"zone_id":               gpdbServerlessTestZone(),
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
				ImportStateVerifyIgnore: []string{"period", "used_time", "db_instance_class", "security_ip_list", "instance_group_count", "create_sample_data", "parameters"},
			},
		},
	})
}

func TestAccAliCloudGPDBDBInstanceServerless_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_instance.default"
	testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region(gpdbServerlessTestRegion())})
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
					"zone_id":               gpdbServerlessTestZone(),
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
				ImportStateVerifyIgnore: []string{"period", "used_time", "db_instance_class", "security_ip_list", "instance_group_count", "create_sample_data", "parameters"},
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
				ImportStateVerifyIgnore: []string{"period", "used_time", "db_instance_class", "security_ip_list", "instance_group_count", "create_sample_data", "parameters"},
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
	zone := gpdbServerlessTestZone()
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
  		zone_id = "%s"
	}

	resource "alicloud_vswitch" "vswitch" {
  		count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  		vpc_id       = data.alicloud_vpcs.default.ids.0
  		cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  		zone_id      = "%s"
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
`, name, zone, zone)
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

// lintignore: R001
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

// TestAccAliCloudGPDBDBInstance_backupId covers the CreateDBInstance.BackupId create path:
// a new instance is created from a data backup set. The test builds its own prerequisites
// end-to-end instead of relying on any pre-existing shared instance or an environment
// variable:
//   - step 0 provisions a real source GPDB instance (Terraform-managed, so it is cleaned up
//     automatically on destroy);
//   - step 1 triggers a manual backup on that source instance via the SDK (there is no
//     Terraform resource for CreateBackup) in PreConfig, waits until it is a completed manual
//     data backup, then creates the target instance whose backup_id is resolved from the
//     alicloud_gpdb_data_backups data source pointed at the source instance.
//
// Both instances live in the same region (testAccPreCheck default) so no region/fixture
// mismatch is possible.
func TestAccAliCloudGPDBDBInstance_backupId(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudGPDBDBInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbbackup%d", defaultRegionToTest, rand)
	// Unique description used to locate the source instance from PreConfig via the SDK.
	sourceDesc := fmt.Sprintf("%s-src", name)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Build the real source instance we control end-to-end.
				Config: testAccGpdbBackupSourceConfig(name, sourceDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_gpdb_instance.source", "id"),
					testAccGpdbEnsureSourceBackupCheck(t, sourceDesc),
				),
			},
			{
				// Create the target instance from the backup set id resolved via the data source.
				Config: testAccGpdbBackupTargetConfig(name, sourceDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_id":            CHECKSET,
						"src_db_instance_name": CHECKSET,
					}),
				),
			},
		},
	})
}

// testAccGpdbBackupSourceConfig builds the Terraform-managed source GPDB instance used as the
// backup origin. Its description is set to a known unique value so PreConfig can locate it.
func testAccGpdbBackupSourceConfig(name, sourceDesc string) string {
	return AliCloudGPDBDBInstanceBasicDependence0(name) + fmt.Sprintf(`
	resource "alicloud_gpdb_instance" "source" {
	  db_instance_category  = "HighAvailability"
	  db_instance_class     = "gpdb.group.segsdx1"
	  db_instance_mode      = "StorageElastic"
	  description           = "%s"
	  engine                = "gpdb"
	  engine_version        = "6.0"
	  zone_id               = data.alicloud_gpdb_zones.default.ids.0
	  instance_network_type = "VPC"
	  instance_spec         = "2C16G"
	  instance_group_count  = "2"
	  payment_type          = "PayAsYouGo"
	  seg_storage_type      = "cloud_essd"
	  seg_node_num          = "4"
	  storage_size          = "50"
	  vpc_id                = data.alicloud_vpcs.default.ids.0
	  vswitch_id            = local.vswitch_id
	  create_sample_data    = "false"
	}
`, sourceDesc)
}

// testAccGpdbBackupTargetConfig keeps the source instance and adds a data source that resolves
// the source instance's completed manual backup set id, then creates the target instance from
// that backup set id.
func testAccGpdbBackupTargetConfig(name, sourceDesc string) string {
	return testAccGpdbBackupSourceConfig(name, sourceDesc) + `
	data "alicloud_gpdb_data_backups" "source" {
	  db_instance_id = alicloud_gpdb_instance.source.id
	  backup_mode    = "Manual"
	  status         = "Success"
	  data_type      = "DATA"
	}

	resource "alicloud_gpdb_instance" "default" {
	  db_instance_category  = "HighAvailability"
	  db_instance_class     = "gpdb.group.segsdx1"
	  db_instance_mode      = "StorageElastic"
	  engine                = "gpdb"
	  engine_version        = "6.0"
	  zone_id               = data.alicloud_gpdb_zones.default.ids.0
	  instance_network_type = "VPC"
	  instance_spec         = "2C16G"
	  instance_group_count  = "2"
	  payment_type          = "PayAsYouGo"
	  seg_storage_type      = "cloud_essd"
	  seg_node_num          = "4"
	  storage_size          = "50"
	  vpc_id                = data.alicloud_vpcs.default.ids.0
	  vswitch_id            = local.vswitch_id
	  create_sample_data    = "false"
	  backup_id             = data.alicloud_gpdb_data_backups.source.backups.0.backup_set_id
	  src_db_instance_name  = alicloud_gpdb_instance.source.id
	}
`
}

// testAccGpdbEnsureSourceBackupCheck wraps the backup-creation side effect as a TestCheckFunc
// so it runs inside a step Check (after the source instance is applied) via a plain function
// call instead of a func literal in the TestStep, which the testing-coverage parser cannot
// marshal.
func testAccGpdbEnsureSourceBackupCheck(t *testing.T, sourceDesc string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		testAccGpdbEnsureSourceBackup(t, sourceDesc)
		return nil
	}
}

// testAccGpdbEnsureSourceBackup ensures the source instance (located by its description) has a
// completed manual data backup visible to DescribeDataBackups, so the data source in the
// target config resolves a backup set id. It is idempotent across step retries.
func testAccGpdbEnsureSourceBackup(t *testing.T, sourceDesc string) {
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-beijing"
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Fatalf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	sourceId := gpdbFindInstanceIdByDescription(t, client, sourceDesc)
	if sourceId == "" {
		t.Fatalf("source GPDB instance with description %q not found", sourceDesc)
	}

	// Idempotent: a completed manual backup may already exist from a prior step attempt.
	if gpdbLatestManualDataBackupId(client, sourceId) != "" {
		return
	}

	backupJobId := gpdbTriggerBackup(t, client, sourceId)
	gpdbWaitBackupJobFinished(t, client, sourceId, backupJobId)

	// Ensure the finished backup is now listed by DescribeDataBackups (what the data source
	// queries) before the target config is planned.
	for i := 0; i < 30; i++ {
		if gpdbLatestManualDataBackupId(client, sourceId) != "" {
			return
		}
		time.Sleep(20 * time.Second)
	}
	t.Fatalf("manual data backup for source instance %s did not become visible in DescribeDataBackups", sourceId)
}

// gpdbFindInstanceIdByDescription locates a GPDB instance id by its DBInstanceDescription.
func gpdbFindInstanceIdByDescription(t *testing.T, client *connectivity.AliyunClient, description string) string {
	action := "DescribeDBInstances"
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	for {
		var response map[string]interface{}
		var err error
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("gpdb", "2016-05-03", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("DescribeDBInstances failed: %s", err)
		}
		resp, _ := jsonpath.Get("$.Items.DBInstance", response)
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if fmt.Sprint(item["DBInstanceDescription"]) == description {
				return fmt.Sprint(item["DBInstanceId"])
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return ""
}

// gpdbLatestManualDataBackupId returns the backup set id of a completed manual data backup on
// the given instance, or an empty string when none exists yet.
func gpdbLatestManualDataBackupId(client *connectivity.AliyunClient, instanceId string) string {
	action := "DescribeDataBackups"
	request := map[string]interface{}{
		"DBInstanceId": instanceId,
		"BackupMode":   "Manual",
		"BackupStatus": "Success",
		"DataType":     "DATA",
		"PageSize":     PageSizeLarge,
		"PageNumber":   1,
	}
	response, err := client.RpcPost("gpdb", "2016-05-03", action, nil, request, true)
	if err != nil {
		return ""
	}
	resp, _ := jsonpath.Get("$.Items[*]", response)
	result, _ := resp.([]interface{})
	if len(result) == 0 {
		return ""
	}
	item := result[0].(map[string]interface{})
	return fmt.Sprint(item["BackupSetId"])
}

// gpdbTriggerBackup starts a manual backup on the instance and returns the backup job id as a
// plain integer string (the API returns it as a JSON number).
func gpdbTriggerBackup(t *testing.T, client *connectivity.AliyunClient, instanceId string) string {
	action := "CreateBackup"
	request := map[string]interface{}{
		"DBInstanceId": instanceId,
	}
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("gpdb", "2016-05-03", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		return nil
	})
	if err != nil {
		t.Fatalf("CreateBackup on instance %s failed: %s", instanceId, err)
	}
	jobId, err := jsonpath.Get("$.BackupJobId", response)
	if err != nil || jobId == nil {
		t.Fatalf("CreateBackup on instance %s returned no BackupJobId: %v", instanceId, response)
	}
	// The API returns BackupJobId as a JSON number; format it as a plain integer string to
	// avoid float scientific notation when it is passed back as a query parameter.
	if f, ok := jobId.(float64); ok {
		return fmt.Sprintf("%.0f", f)
	}
	return fmt.Sprint(jobId)
}

// gpdbWaitBackupJobFinished polls DescribeBackupJob until the job reaches the finish state and
// returns the resulting backup set id.
func gpdbWaitBackupJobFinished(t *testing.T, client *connectivity.AliyunClient, instanceId string, backupJobId string) string {
	action := "DescribeBackupJob"
	request := map[string]interface{}{
		"DBInstanceId": instanceId,
		"BackupJobId":  backupJobId,
	}
	var backupId string
	err := resource.Retry(30*time.Minute, func() *resource.RetryError {
		response, err := client.RpcPost("gpdb", "2016-05-03", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		status, _ := jsonpath.Get("$.BackupStatus", response)
		if fmt.Sprint(status) != "finish" {
			return resource.RetryableError(fmt.Errorf("backup job %v on instance %s not finished, current status: %v", backupJobId, instanceId, status))
		}
		if id, e := jsonpath.Get("$.BackupId", response); e == nil && id != nil {
			backupId = fmt.Sprint(id)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("waiting for backup job %v on instance %s to finish failed: %s", backupJobId, instanceId, err)
	}
	if backupId == "" {
		// DescribeBackupJob does not always echo the backup set id; fall back to the latest
		// completed manual data backup on the source instance (DescribeDataBackups).
		backupId = gpdbLatestManualDataBackupId(client, instanceId)
	}
	return backupId
}
