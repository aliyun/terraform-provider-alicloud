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
	resource.AddTestSweepers(
		"alicloud_graph_database_db_instance",
		&resource.Sweeper{
			Name: "alicloud_graph_database_db_instance",
			F:    testSweepGraphDatabaseDbInstance,
		})
}

func testSweepGraphDatabaseDbInstance(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeDBInstances"
	request := map[string]interface{}{}

	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	conn, err := client.NewGdsClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-03"), StringPointer("AK"), nil, request, &runtime)
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
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["DBInstanceDescription"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if item["DBInstanceStatus"].(string) != "Running" {
				skip = true
			}
			if skip {
				log.Printf("[INFO] Skipping Graph Database DbInstance: %s", item["DBInstanceDescription"].(string))
				continue
			}
			action := "DeleteDBInstance"
			deleteRequest := map[string]interface{}{
				"DBInstanceId": item["DBInstanceId"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-03"), StringPointer("AK"), nil, deleteRequest, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Graph Database DbInstance (%s): %s", item["DBInstanceDescription"].(string), err)
			}
			log.Printf("[INFO] Delete Graph Database DbInstance success: %s ", item["DBInstanceDescription"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudGraphDatabaseDbInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_graph_database_db_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudGraphDatabaseDbInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGraphDatabaseDbInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgraphdatabasedbinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGraphDatabaseDbInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.GraphDatabaseSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_node_class":            "gdb.r.xlarge",
					"db_instance_network_type": "vpc",
					"db_version":               "1.0",
					"db_instance_category":     "HA",
					"db_instance_storage_type": "cloud_essd",
					"db_node_storage":          "50",
					"payment_type":             "PayAsYouGo",
					"db_instance_description":  "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_class":            "gdb.r.xlarge",
						"db_instance_network_type": "vpc",
						"db_version":               "1.0",
						"db_instance_category":     "HA",
						"db_instance_storage_type": "cloud_essd",
						"db_node_storage":          "50",
						"payment_type":             "PayAsYouGo",
						"db_instance_description":  name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_ip_array": []map[string]interface{}{
						{
							"db_instance_ip_array_name": "default",
							"security_ips":              "127.0.0.2",
						},
						{
							"db_instance_ip_array_name": "tftest",
							"security_ips":              "192.168.0.1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_ip_array.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_description": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_node_storage": "80",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_storage": "80",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_node_class": "gdb.r.2xlarge",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_class": "gdb.r.2xlarge",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_description": "${var.name}",
					"db_node_storage":         "100",
					"db_node_class":           "gdb.r.xlarge",
					"db_instance_ip_array": []map[string]interface{}{
						{
							"db_instance_ip_array_name": "default",
							"security_ips":              "127.0.0.1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_description": name,
						"db_node_storage":         "100",
						"db_node_class":           "gdb.r.xlarge",
						"db_instance_ip_array.#":  "1",
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

func TestAccAlicloudGraphDatabaseDbInstance_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_graph_database_db_instance.default"
	checkoutSupportedRegions(t, true, connectivity.GraphDatabaseDbInstanceSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudGraphDatabaseDbInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGraphDatabaseDbInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgraphdatabasedbinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGraphDatabaseDbInstanceBasicDependence1)
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
					"db_node_class":            "gdb.r.xlarge",
					"db_instance_network_type": "vpc",
					"db_version":               "1.0",
					"db_instance_category":     "HA",
					"db_instance_storage_type": "cloud_essd",
					"db_node_storage":          "50",
					"payment_type":             "PayAsYouGo",
					"db_instance_description":  "${var.name}",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"vpc_id":                   "${data.alicloud_vpcs.default.ids.0}",
					"zone_id":                  "${local.zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_class":            "gdb.r.xlarge",
						"db_instance_network_type": "vpc",
						"db_version":               "1.0",
						"db_instance_category":     "HA",
						"db_instance_storage_type": "cloud_essd",
						"db_node_storage":          "50",
						"payment_type":             "PayAsYouGo",
						"db_instance_description":  name,
						"vswitch_id":               CHECKSET,
						"vpc_id":                   CHECKSET,
						"zone_id":                  CHECKSET,
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

func TestAccAlicloudGraphDatabaseDbInstance_single(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_graph_database_db_instance.default"
	checkoutSupportedRegions(t, true, connectivity.GraphDatabaseDbInstanceSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudGraphDatabaseDbInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGraphDatabaseDbInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgraphdatabasedbinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGraphDatabaseDbInstanceBasicDependence1)
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
					"db_node_class":            "gdb.r.xlarge_basic",
					"db_instance_network_type": "vpc",
					"db_version":               "1.0",
					"db_instance_category":     "SINGLE",
					"db_instance_storage_type": "cloud_essd",
					"db_node_storage":          "50",
					"payment_type":             "PayAsYouGo",
					"db_instance_description":  "${var.name}",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"vpc_id":                   "${data.alicloud_vpcs.default.ids.0}",
					"zone_id":                  "${local.zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_class":            "gdb.r.xlarge_basic",
						"db_instance_network_type": "vpc",
						"db_version":               "1.0",
						"db_instance_category":     "SINGLE",
						"db_instance_storage_type": "cloud_essd",
						"db_node_storage":          "50",
						"payment_type":             "PayAsYouGo",
						"db_instance_description":  name,
						"vswitch_id":               CHECKSET,
						"vpc_id":                   CHECKSET,
						"zone_id":                  CHECKSET,
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

var AlicloudGraphDatabaseDbInstanceMap0 = map[string]string{}

func AlicloudGraphDatabaseDbInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}`, name)
}

//The instance requires the same zone as the vswitch, but currently the instance does not support zone query.
func AlicloudGraphDatabaseDbInstanceBasicDependence1(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
locals {
  zone_id = "cn-hangzhou-h"
}
data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = local.zone_id
}`, name)
}
