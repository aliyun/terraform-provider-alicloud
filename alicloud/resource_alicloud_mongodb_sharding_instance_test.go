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
	resource.AddTestSweepers("alicloud_mongodb_sharding_instance", &resource.Sweeper{
		Name: "alicloud_mongodb_sharding_instance",
		F:    testSweepMongoDBShardingInstances,
	})
}

func testSweepMongoDBShardingInstances(region string) error {

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
	request["DBInstanceType"] = "sharding"
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
				log.Printf("[INFO] Skipping Mongodb Sharding Instance: %s", fmt.Sprint(item["DBInstanceDescription"]))
				continue
			}
			action := "DeleteDBInstance"
			request := map[string]interface{}{
				"DBInstanceId": item["DBInstanceId"],
				"RegionId":     client.RegionId,
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Mongodb Sharding Instance (%s): %s", fmt.Sprint(item["DBInstanceDescription"]), err)
			}
			log.Printf("[INFO] Delete Mongodb Sharding Instance success: %s ", fmt.Sprint(item["DBInstanceDescription"]))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func resourceMongodbShardingInstanceClassicConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_mongodb_zones" "default" {}

	resource "alicloud_security_group" "default" {
		name = var.name
	}
`, name)
}

func TestAccAlicloudMongoDBShardingInstance_classic(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alicloud_mongodb_sharding_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBShardingInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBShardingInstanceClassicConfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMongodbShardingInstanceClassicConfig)
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
					"zone_id":        "${data.alicloud_mongodb_zones.default.zones.0.id}",
					"engine_version": "3.4",
					"name":           name,
					"shard_list": []map[string]interface{}{
						{
							"node_class":   "dds.shard.mid",
							"node_storage": "10",
						},
						{
							"node_class":        "dds.shard.standard",
							"node_storage":      "20",
							"readonly_replicas": "1",
						},
					},
					"mongo_list": []map[string]interface{}{
						{
							"node_class": "dds.mongos.mid",
						},
						{
							"node_class": "dds.mongos.mid",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":                        CHECKSET,
						"engine_version":                 "3.4",
						"shard_list.#":                   "2",
						"shard_list.0.node_class":        "dds.shard.mid",
						"shard_list.0.node_storage":      "10",
						"shard_list.0.readonly_replicas": "0",
						"shard_list.1.node_class":        "dds.shard.standard",
						"shard_list.1.node_storage":      "20",
						"shard_list.1.readonly_replicas": "1",
						"mongo_list.#":                   "2",
						"mongo_list.0.node_class":        "dds.mongos.mid",
						"mongo_list.1.node_class":        "dds.mongos.mid",
						"name":                           name,
						"storage_engine":                 "WiredTiger",
						"instance_charge_type":           "PostPaid",
						"tags.%":                         "0",
						"config_server_list.#":           CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"order_type", "auto_renew"},
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
					"account_password": "YourPassword_",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "YourPassword_",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mongo_list": []map[string]interface{}{
						{
							"node_class": "dds.mongos.mid",
						},
						{
							"node_class": "dds.mongos.mid",
						},
						{
							"node_class": "dds.mongos.mid",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mongo_list.#":            "3",
						"mongo_list.0.node_class": "dds.mongos.mid",
						"mongo_list.1.node_class": "dds.mongos.mid",
						"mongo_list.2.node_class": "dds.mongos.mid",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"shard_list": []map[string]interface{}{
						{
							"node_class":   "dds.shard.mid",
							"node_storage": "10",
						},
						{
							"node_class":        "dds.shard.standard",
							"node_storage":      "20",
							"readonly_replicas": "1",
						},
						{
							"node_class":        "dds.shard.standard",
							"node_storage":      "20",
							"readonly_replicas": "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"shard_list.#":                   "3",
						"shard_list.0.node_class":        "dds.shard.mid",
						"shard_list.0.node_storage":      "10",
						"shard_list.0.readonly_replicas": "0",
						"shard_list.1.node_class":        "dds.shard.standard",
						"shard_list.1.node_storage":      "20",
						"shard_list.1.readonly_replicas": "1",
						"shard_list.2.node_class":        "dds.shard.standard",
						"shard_list.2.node_storage":      "20",
						"shard_list.2.readonly_replicas": "2",
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
					"name":             name,
					"account_password": "YourPassword_123",
					"security_ip_list": []string{"10.168.1.12", "10.168.1.13"},
					"backup_period":    []string{"Tuesday", "Wednesday"},
					"backup_time":      "10:00Z-11:00Z",
					"shard_list": []map[string]interface{}{
						{
							"node_class":   "dds.shard.mid",
							"node_storage": "10",
						},
						{
							"node_class":        "dds.shard.standard",
							"node_storage":      "20",
							"readonly_replicas": "1",
						},
					},
					"mongo_list": []map[string]interface{}{
						{
							"node_class": "dds.mongos.mid",
						},
						{
							"node_class": "dds.mongos.mid",
						},
					},
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                           name,
						"account_password":               "YourPassword_123",
						"security_ip_list.#":             "2",
						"backup_period.#":                "2",
						"backup_time":                    "10:00Z-11:00Z",
						"shard_list.#":                   "2",
						"shard_list.0.node_class":        "dds.shard.mid",
						"shard_list.0.node_storage":      "10",
						"shard_list.1.node_class":        "dds.shard.standard",
						"shard_list.1.node_storage":      "20",
						"shard_list.1.readonly_replicas": "1",
						"shard_list.2.node_class":        REMOVEKEY,
						"shard_list.2.readonly_replicas": REMOVEKEY,
						"shard_list.2.node_storage":      REMOVEKEY,
						"mongo_list.#":                   "2",
						"mongo_list.0.node_class":        "dds.mongos.mid",
						"mongo_list.1.node_class":        "dds.mongos.mid",
						"mongo_list.2.node_class":        REMOVEKEY,
						"tags.%":                         "0",
						"tags.For":                       REMOVEKEY,
						"tags.Created":                   REMOVEKEY,
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

func TestAccAlicloudMongoDBShardingInstance_vpc(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alicloud_mongodb_sharding_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBShardingInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBShardingInstanceVpcConfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMongodbShardingInstanceVpcConfig)
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
					"zone_id":        "${data.alicloud_mongodb_zones.default.zones.0.id}",
					"vswitch_id":     "${data.alicloud_vswitches.default.ids.0}",
					"engine_version": "4.0",
					"name":           name,
					"shard_list": []map[string]interface{}{
						{
							"node_class":   "dds.shard.mid",
							"node_storage": "10",
						},
						{
							"node_class":        "dds.shard.standard",
							"node_storage":      "20",
							"readonly_replicas": "1",
						},
					},
					"mongo_list": []map[string]interface{}{
						{
							"node_class": "dds.mongos.mid",
						},
						{
							"node_class": "dds.mongos.mid",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":                        CHECKSET,
						"vswitch_id":                     CHECKSET,
						"engine_version":                 "4.0",
						"shard_list.#":                   "2",
						"shard_list.0.node_class":        "dds.shard.mid",
						"shard_list.0.node_storage":      "10",
						"shard_list.0.readonly_replicas": "0",
						"shard_list.1.node_class":        "dds.shard.standard",
						"shard_list.1.node_storage":      "20",
						"shard_list.1.readonly_replicas": "1",
						"mongo_list.#":                   "2",
						"mongo_list.0.node_class":        "dds.mongos.mid",
						"mongo_list.1.node_class":        "dds.mongos.mid",
						"name":                           name,
						"storage_engine":                 "WiredTiger",
						"instance_charge_type":           "PostPaid",
						"tags.%":                         "0",
						"config_server_list.#":           CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"order_type", "auto_renew"},
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
					"tde_status": "enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tde_status": "enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "YourPassword_",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "YourPassword_",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mongo_list": []map[string]interface{}{
						{
							"node_class": "dds.mongos.mid",
						},
						{
							"node_class": "dds.mongos.mid",
						},
						{
							"node_class": "dds.mongos.mid",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mongo_list.#":            "3",
						"mongo_list.0.node_class": "dds.mongos.mid",
						"mongo_list.1.node_class": "dds.mongos.mid",
						"mongo_list.2.node_class": "dds.mongos.mid",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"shard_list": []map[string]interface{}{
						{
							"node_class":   "dds.shard.mid",
							"node_storage": "10",
						},
						{
							"node_class":        "dds.shard.standard",
							"node_storage":      "20",
							"readonly_replicas": "1",
						},
						{
							"node_class":        "dds.shard.standard",
							"node_storage":      "20",
							"readonly_replicas": "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"shard_list.#":                   "3",
						"shard_list.0.node_class":        "dds.shard.mid",
						"shard_list.0.node_storage":      "10",
						"shard_list.0.readonly_replicas": "0",
						"shard_list.1.node_class":        "dds.shard.standard",
						"shard_list.1.node_storage":      "20",
						"shard_list.1.readonly_replicas": "1",
						"shard_list.2.node_class":        "dds.shard.standard",
						"shard_list.2.node_storage":      "20",
						"shard_list.2.readonly_replicas": "2",
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
					"name":             name,
					"account_password": "YourPassword_123",
					"security_ip_list": []string{"10.168.1.12", "10.168.1.13"},
					"backup_period":    []string{"Tuesday", "Wednesday"},
					"backup_time":      "10:00Z-11:00Z",
					"shard_list": []map[string]interface{}{
						{
							"node_class":   "dds.shard.mid",
							"node_storage": "10",
						},
						{
							"node_class":        "dds.shard.standard",
							"node_storage":      "20",
							"readonly_replicas": "1",
						},
					},
					"mongo_list": []map[string]interface{}{
						{
							"node_class": "dds.mongos.mid",
						},
						{
							"node_class": "dds.mongos.mid",
						},
					},
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                           name,
						"account_password":               "YourPassword_123",
						"security_ip_list.#":             "2",
						"backup_period.#":                "2",
						"backup_time":                    "10:00Z-11:00Z",
						"shard_list.#":                   "2",
						"shard_list.0.node_class":        "dds.shard.mid",
						"shard_list.0.node_storage":      "10",
						"shard_list.1.node_class":        "dds.shard.standard",
						"shard_list.1.node_storage":      "20",
						"shard_list.1.readonly_replicas": "1",
						"shard_list.2.node_class":        REMOVEKEY,
						"shard_list.2.readonly_replicas": REMOVEKEY,
						"shard_list.2.node_storage":      REMOVEKEY,
						"mongo_list.#":                   "2",
						"mongo_list.0.node_class":        "dds.mongos.mid",
						"mongo_list.1.node_class":        "dds.mongos.mid",
						"mongo_list.2.node_class":        REMOVEKEY,
						"tags.%":                         "0",
						"tags.For":                       REMOVEKEY,
						"tags.Created":                   REMOVEKEY,
					}),
				),
			},
		},
	})
}

func resourceMongodbShardingInstanceVpcConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	data "alicloud_resource_manager_resource_groups" "default" {
  		status = "OK"
	}
	data "alicloud_mongodb_zones" "default" {}
	data "alicloud_vpcs" "default" {
		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
	  vpc_id = data.alicloud_vpcs.default.ids.0
	  zone_id = "${data.alicloud_mongodb_zones.default.zones.0.id}"
	}
`, name)
}

func TestAccAlicloudMongoDBShardingInstance_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_sharding_instance.default"
	checkoutSupportedRegions(t, false, connectivity.MongoDBClassicNoSupportedRegions)
	ra := resourceAttrInit(resourceId, AlicloudMongoDBShardingInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMongoDBShardingInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-mongodbshardinginstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMongodbShardingInstanceVpcConfig)
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
					"name":              "${var.name}",
					"zone_id":           "${data.alicloud_mongodb_zones.default.zones.0.id}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"engine_version":    "3.4",
					"protocol_type":     "mongodb",
					"network_type":      "Classic",
					"shard_list": []map[string]interface{}{
						{
							"node_class":   "dds.shard.mid",
							"node_storage": "10",
						},
						{
							"node_class":        "dds.shard.standard",
							"node_storage":      "20",
							"readonly_replicas": "1",
						},
					},
					"mongo_list": []map[string]interface{}{
						{
							"node_class": "dds.mongos.mid",
						},
						{
							"node_class": "dds.mongos.mid",
						},
					},

					"storage_engine":       "WiredTiger",
					"instance_charge_type": "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":                        CHECKSET,
						"engine_version":                 "3.4",
						"shard_list.#":                   "2",
						"shard_list.0.node_class":        "dds.shard.mid",
						"shard_list.0.node_storage":      "10",
						"shard_list.0.readonly_replicas": "0",
						"shard_list.1.node_class":        "dds.shard.standard",
						"shard_list.1.node_storage":      "20",
						"shard_list.1.readonly_replicas": "1",
						"mongo_list.#":                   "2",
						"mongo_list.0.node_class":        "dds.mongos.mid",
						"mongo_list.1.node_class":        "dds.mongos.mid",
						"name":                           name,
						"storage_engine":                 "WiredTiger",
						"instance_charge_type":           "PostPaid",
						"tags.%":                         "0",
						"protocol_type":                  "mongodb",
						"network_type":                   "Classic",
						"config_server_list.#":           CHECKSET,
						"resource_group_id":              CHECKSET,
					}),
				),
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
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"order_type", "auto_renew"},
			},
		},
	})
}

func TestAccAlicloudMongoDBShardingInstance_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_sharding_instance.default"
	checkoutSupportedRegions(t, false, connectivity.MongoDBClassicNoSupportedRegions)
	ra := resourceAttrInit(resourceId, AlicloudMongoDBShardingInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMongoDBShardingInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-mongodbshardinginstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMongodbShardingInstanceVpcConfig)
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
					"name":              "${var.name}",
					"zone_id":           "${data.alicloud_mongodb_zones.default.zones.0.id}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"engine_version":    "3.4",
					"protocol_type":     "mongodb",
					"network_type":      "VPC",
					"vpc_id":            "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"shard_list": []map[string]interface{}{
						{
							"node_class":   "dds.shard.mid",
							"node_storage": "10",
						},
						{
							"node_class":        "dds.shard.standard",
							"node_storage":      "20",
							"readonly_replicas": "1",
						},
					},
					"mongo_list": []map[string]interface{}{
						{
							"node_class": "dds.mongos.mid",
						},
						{
							"node_class": "dds.mongos.mid",
						},
					},

					"storage_engine":       "WiredTiger",
					"instance_charge_type": "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":                        CHECKSET,
						"engine_version":                 "3.4",
						"shard_list.#":                   "2",
						"shard_list.0.node_class":        "dds.shard.mid",
						"shard_list.0.node_storage":      "10",
						"shard_list.0.readonly_replicas": "0",
						"shard_list.1.node_class":        "dds.shard.standard",
						"shard_list.1.node_storage":      "20",
						"shard_list.1.readonly_replicas": "1",
						"mongo_list.#":                   "2",
						"mongo_list.0.node_class":        "dds.mongos.mid",
						"mongo_list.1.node_class":        "dds.mongos.mid",
						"name":                           name,
						"storage_engine":                 "WiredTiger",
						"instance_charge_type":           "PostPaid",
						"tags.%":                         "0",
						"protocol_type":                  "mongodb",
						"network_type":                   "VPC",
						"vpc_id":                         CHECKSET,
						"vswitch_id":                     CHECKSET,
						"config_server_list.#":           CHECKSET,
						"resource_group_id":              CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"order_type", "auto_renew"},
			},
		},
	})
}

var AlicloudMongoDBShardingInstanceMap0 = map[string]string{}
