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

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_mongodb_instance", &resource.Sweeper{
		Name: "alicloud_mongodb_instance",
		F:    testSweepMongoDBInstances,
	})
}

func testSweepMongoDBInstances(region string) error {

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
	request["ChargeType"] = "PostPaid"

	var response map[string]interface{}
	conn, err := client.NewDdsClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &runtime)
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

		resp, err := jsonpath.Get("$.DBInstances.DBInstance", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.DBInstances.DBInstance", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["DBInstanceDescription"])), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Mongodb Instance: %s", fmt.Sprint(item["DBInstanceDescription"]))
				continue
			}
			action := "DeleteDBInstance"
			request := map[string]interface{}{
				"DBInstanceId": item["DBInstanceId"],
				"RegionId":     client.RegionId,
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Mongodb Instance (%s): %s", fmt.Sprint(item["DBInstanceDescription"]), err)
			}
			log.Printf("[INFO] Delete Mongodb Instance success: %s ", fmt.Sprint(item["DBInstanceDescription"]))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudMongoDBInstance_classic(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alicloud_mongodb_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBInstanceClassicConfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMongodbInstanceClassicConfig)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.MongoDBClassicNoSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":             "${local.zone_id}",
					"engine_version":      "3.4",
					"db_instance_storage": "10",
					"db_instance_class":   "dds.mongo.mid",
					"name":                name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version":       "3.4",
						"db_instance_storage":  "10",
						"db_instance_class":    "dds.mongo.mid",
						"name":                 name,
						"storage_engine":       "WiredTiger",
						"instance_charge_type": "PostPaid",
						"replication_factor":   "3",
						"replica_sets.#":       CHECKSET,
						"ssl_status":           CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ssl_action", "order_type", "auto_renew"},
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
					"name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_storage": "30",
					"db_instance_class":   "dds.mongo.standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage": "30",
						"db_instance_class":   "dds.mongo.standard",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "YourPassword_123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "YourPassword_123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_list": []string{"10.168.1.12"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_id": "${alicloud_security_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_period": []string{"Wednesday"},
					"backup_time":   "11:00Z-12:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "1",
						"backup_time":     "11:00Z-12:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_start_time": "02:00Z",
					"maintain_end_time":   "03:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_start_time": "02:00Z",
						"maintain_end_time":   "03:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                name,
					"account_password":    "YourPassword_",
					"security_ip_list":    []string{"10.168.1.12", "10.168.1.13"},
					"db_instance_storage": "30",
					"db_instance_class":   "dds.mongo.standard",
					"backup_period":       []string{"Tuesday", "Wednesday"},
					"backup_time":         "10:00Z-11:00Z",
					"maintain_start_time": REMOVEKEY,
					"maintain_end_time":   REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                        name,
						"account_password":            "YourPassword_",
						"security_ip_list.#":          "2",
						"security_ip_list.4095458986": "10.168.1.12",
						"security_ip_list.3976237035": "10.168.1.13",
						"db_instance_storage":         "30",
						"db_instance_class":           "dds.mongo.standard",
						"backup_period.#":             "2",
						"backup_period.1592931319":    "Tuesday",
						"backup_period.1970423419":    "Wednesday",
						"backup_time":                 "10:00Z-11:00Z",
						"maintain_start_time":         REMOVEKEY,
						"maintain_end_time":           REMOVEKEY,
						"ssl_status":                  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type": "PrePaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "PrePaid",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudMongoDBInstance_classic1(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alicloud_mongodb_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBInstanceClassicConfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMongodbInstanceClassicConfig)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.MongoDBClassicNoSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":             "${local.zone_id}",
					"engine_version":      "4.0",
					"db_instance_storage": "2000",
					"db_instance_class":   "dds.mongo.2xlarge",
					"name":                name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version":       "4.0",
						"db_instance_storage":  "2000",
						"db_instance_class":    "dds.mongo.2xlarge",
						"name":                 name,
						"storage_engine":       "WiredTiger",
						"instance_charge_type": "PostPaid",
						"replication_factor":   "3",
						"replica_sets.#":       CHECKSET,
						"ssl_status":           CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ssl_action", "order_type", "auto_renew"},
			},
		},
	})
}

func resourceMongodbInstanceClassicConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_mongodb_zones" "default" {
	}

	locals {
  		zone_id = data.alicloud_mongodb_zones.default.zones[length(data.alicloud_mongodb_zones.default.zones) - 1].id
	}

	resource "alicloud_security_group" "default" {
		name = var.name
	}
`, name)
}

func TestAccAlicloudMongoDBInstance_vpc(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alicloud_mongodb_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBInstanceVpcConfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMongodbInstanceVpcConfig)
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
					"vswitch_id":          "${local.vswitch_id}",
					"engine_version":      "4.0",
					"db_instance_storage": "10",
					"db_instance_class":   "dds.mongo.mid",
					"name":                name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version":       "4.0",
						"db_instance_storage":  "10",
						"db_instance_class":    "dds.mongo.mid",
						"name":                 name,
						"storage_engine":       "WiredTiger",
						"instance_charge_type": "PostPaid",
						"replication_factor":   "3",
						"replica_sets.#":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ssl_action", "order_type", "auto_renew"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_storage": "30",
					"db_instance_class":   "dds.mongo.standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage": "30",
						"db_instance_class":   "dds.mongo.standard",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "YourPassword_123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "YourPassword_123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_list": []string{"10.168.1.12"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_period": []string{"Wednesday"},
					"backup_time":   "11:00Z-12:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "1",
						"backup_time":     "11:00Z-12:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                name,
					"account_password":    "YourPassword_",
					"security_ip_list":    []string{"10.168.1.12", "10.168.1.13"},
					"db_instance_storage": "30",
					"db_instance_class":   "dds.mongo.standard",
					"backup_period":       []string{"Tuesday", "Wednesday"},
					"backup_time":         "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                        name,
						"account_password":            "YourPassword_",
						"security_ip_list.#":          "2",
						"security_ip_list.4095458986": "10.168.1.12",
						"security_ip_list.3976237035": "10.168.1.13",
						"db_instance_storage":         "30",
						"db_instance_class":           "dds.mongo.standard",
						"backup_period.#":             "2",
						"backup_time":                 "10:00Z-11:00Z",
					}),
				),
			}},
	})
}

func TestAccAlicloudMongoDBInstance_vpc1(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alicloud_mongodb_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBInstanceVpcConfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMongodbInstanceVpcConfig)
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
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"vswitch_id":          "${local.vswitch_id}",
					"engine_version":      "4.0",
					"network_type":        "VPC",
					"vpc_id":              "${data.alicloud_vpcs.default.ids.0}",
					"zone_id":             "${local.zone_id}",
					"db_instance_storage": "10",
					"db_instance_class":   "dds.mongo.mid",
					"name":                name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version":       "4.0",
						"network_type":         "VPC",
						"db_instance_storage":  "10",
						"db_instance_class":    "dds.mongo.mid",
						"name":                 name,
						"storage_engine":       "WiredTiger",
						"instance_charge_type": "PostPaid",
						"replication_factor":   "3",
						"replica_sets.#":       CHECKSET,
						"resource_group_id":    CHECKSET,
						"vpc_id":               CHECKSET,
						"vswitch_id":           CHECKSET,
						"zone_id":              CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ssl_action", "order_type", "auto_renew"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
		},
	})
}

func resourceMongodbInstanceVpcConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	
	data "alicloud_resource_manager_resource_groups" "default" {
  		status = "OK"
	}
	
	data "alicloud_mongodb_zones" "default" {
	}
	
	data "alicloud_vpcs" "default" {
	  	name_regex = "default-NODELETING"
	}
	
	data "alicloud_vswitches" "default" {
	  	vpc_id = data.alicloud_vpcs.default.ids.0
	}
	
	resource "alicloud_vswitch" "this" {
	  	count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
	  	vswitch_name = var.name
	  	vpc_id       = data.alicloud_vpcs.default.ids.0
	  	zone_id      = data.alicloud_mongodb_zones.default.ids.0
	  	cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
	}
	
	locals {
	  	zone_id    = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.vswitches.0.zone_id : data.alicloud_mongodb_zones.default.zones.0.id
	  	vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.this.*.id, [""])[0]
	}
`, name)
}

func TestAccAlicloudMongoDBInstance_multiAZ(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alicloud_mongodb_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBInstanceMultiAZConfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMongodbInstanceMultiAZConfig)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.MongoDBMultiAzSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":             "${local.zone_id}",
					"vswitch_id":          "${local.vswitch_id}",
					"engine_version":      "3.4",
					"db_instance_storage": "10",
					"db_instance_class":   "dds.mongo.mid",
					"name":                name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version":       "3.4",
						"db_instance_storage":  "10",
						"db_instance_class":    "dds.mongo.mid",
						"name":                 name,
						"storage_engine":       "WiredTiger",
						"instance_charge_type": "PostPaid",
						"replication_factor":   "3",
						"replica_sets.#":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ssl_action", "order_type", "auto_renew"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_storage": "30",
					"db_instance_class":   "dds.mongo.standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage": "30",
						"db_instance_class":   "dds.mongo.standard",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "YourPassword_123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "YourPassword_123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_list": []string{"10.168.1.12"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_period": []string{"Wednesday"},
					"backup_time":   "11:00Z-12:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "1",
						"backup_time":     "11:00Z-12:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                name,
					"account_password":    "YourPassword_",
					"security_ip_list":    []string{"10.168.1.12", "10.168.1.13"},
					"db_instance_storage": "30",
					"db_instance_class":   "dds.mongo.standard",
					"backup_period":       []string{"Tuesday", "Wednesday"},
					"backup_time":         "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                        name,
						"account_password":            "YourPassword_",
						"security_ip_list.#":          "2",
						"security_ip_list.4095458986": "10.168.1.12",
						"security_ip_list.3976237035": "10.168.1.13",
						"db_instance_storage":         "30",
						"db_instance_class":           "dds.mongo.standard",
						"backup_period.#":             "2",
						"backup_time":                 "10:00Z-11:00Z",
					}),
				),
			},
		},
	})
}

func resourceMongodbInstanceMultiAZConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_mongodb_zones" "default" {
		multi = true
	}
	data "alicloud_vpcs" "default" {
		name_regex = "default-NODELETING"
	}

	locals {
  		zone_id = data.alicloud_mongodb_zones.default.zones.0.multi_zone_ids[length(data.alicloud_mongodb_zones.default.zones.0.multi_zone_ids) - 1]
	}
	
	data "alicloud_vswitches" "default" {
	  vpc_id = data.alicloud_vpcs.default.ids.0
	  zone_id = local.zone_id
	}
	
	locals {
	  vswitch_id = data.alicloud_vswitches.default.ids[0]
	}
`, name)
}
