package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

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
		return fmt.Errorf("error getting AliCloud client: %s", err)
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
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["DBInstanceDescription"])), strings.ToLower(prefix)) {
						skip = false
					}
				}
				if skip {
					log.Printf("[INFO] Skipping Mongodb Sharding Instance: %s", fmt.Sprint(item["DBInstanceDescription"]))
					continue
				}
			}
			action := "DeleteDBInstance"
			request := map[string]interface{}{
				"DBInstanceId": item["DBInstanceId"],
				"RegionId":     client.RegionId,
			}
			_, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, false)
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

func TestAccAliCloudMongoDBShardingInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_sharding_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	ra := resourceAttrInit(resourceId, AliCloudMongoDBShardingInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBShardingInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBShardingInstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMongoDBShardingInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version": "4.2",
					"vswitch_id":     "${data.alicloud_vswitches.default.ids.0}",
					"mongo_list": []map[string]interface{}{
						{
							"node_class": "dds.mongos.mid",
						},
						{
							"node_class": "dds.mongos.mid",
						},
					},
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version": "4.2",
						"vswitch_id":     CHECKSET,
						"mongo_list.#":   "2",
						"shard_list.#":   "2",
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
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type": "PrePaid",
					"period":               "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "PrePaid",
						"period":               "1",
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
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_time": "11:00Z-12:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time": "11:00Z-12:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_period": []string{"Monday", "Tuesday", "Wednesday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "3",
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
					"global_security_group_list": "${alicloud_mongodb_global_security_ip_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_security_group_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_security_group_list": []string{"${alicloud_mongodb_global_security_ip_group.default.0.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_security_group_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_security_group_list": "${alicloud_mongodb_global_security_ip_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_security_group_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ShardingInstance",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "ShardingInstance",
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
							"node_class": "dds.mongos.standard",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mongo_list.#": "3",
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
						"shard_list.#": "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "kms_encrypted_password", "kms_encryption_context", "auto_renew", "ssl_action", "order_type"},
			},
		},
	})
}

func TestAccAliCloudMongoDBShardingInstance_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_sharding_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	ra := resourceAttrInit(resourceId, AliCloudMongoDBShardingInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBShardingInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBShardingInstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMongoDBShardingInstanceBasicDependence1)
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
					"engine_version":             "4.2",
					"storage_engine":             "WiredTiger",
					"storage_type":               "local_ssd",
					"protocol_type":              "mongodb",
					"vpc_id":                     "${alicloud_vswitch.default.vpc_id}",
					"vswitch_id":                 "${alicloud_vswitch.default.id}",
					"zone_id":                    "${data.alicloud_mongodb_zones.default.zones.0.id}",
					"security_group_id":          "${alicloud_security_group.default.id}",
					"network_type":               "VPC",
					"name":                       name,
					"instance_charge_type":       "PostPaid",
					"security_ip_list":           []string{"10.168.1.12"},
					"account_password":           "YourPassword_123",
					"resource_group_id":          "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"backup_time":                "11:00Z-12:00Z",
					"backup_period":              []string{"Monday", "Tuesday", "Wednesday"},
					"tde_status":                 "enabled",
					"global_security_group_list": "${alicloud_mongodb_global_security_ip_group.default.*.id}",
					"mongo_list": []map[string]interface{}{
						{
							"node_class": "dds.mongos.mid",
						},
						{
							"node_class": "dds.mongos.mid",
						},
					},
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
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ShardingInstance",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version":               "4.2",
						"storage_engine":               "WiredTiger",
						"storage_type":                 "local_ssd",
						"protocol_type":                "mongodb",
						"vpc_id":                       CHECKSET,
						"vswitch_id":                   CHECKSET,
						"zone_id":                      CHECKSET,
						"security_group_id":            CHECKSET,
						"network_type":                 "VPC",
						"name":                         name,
						"instance_charge_type":         "PostPaid",
						"security_ip_list.#":           "1",
						"account_password":             "YourPassword_123",
						"resource_group_id":            CHECKSET,
						"backup_time":                  "11:00Z-12:00Z",
						"backup_period.#":              "3",
						"tde_status":                   "enabled",
						"global_security_group_list.#": "3",
						"mongo_list.#":                 "2",
						"shard_list.#":                 "2",
						"tags.%":                       "2",
						"tags.Created":                 "TF",
						"tags.For":                     "ShardingInstance",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "kms_encrypted_password", "kms_encryption_context", "auto_renew", "ssl_action", "order_type"},
			},
		},
	})
}

func TestAccAliCloudMongoDBShardingInstance_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_sharding_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	ra := resourceAttrInit(resourceId, AliCloudMongoDBShardingInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBShardingInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBShardingInstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMongoDBShardingInstanceBasicDependence1)
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
					"engine_version": "4.4",
					"vswitch_id":     "${alicloud_vswitch.default.id}",
					"mongo_list": []map[string]interface{}{
						{
							"node_class": "mdb.shard.8x.large.d",
						},
						{
							"node_class": "mdb.shard.8x.large.d",
						},
					},
					"shard_list": []map[string]interface{}{
						{
							"node_class":   "mdb.shard.8x.large.d",
							"node_storage": "50",
						},
						{
							"node_class":        "mdb.shard.8x.xlarge.d",
							"node_storage":      "60",
							"readonly_replicas": "1",
						},
					},
					"config_server_list": []map[string]interface{}{
						{
							"node_class":   "mdb.shard.2x.xlarge.d",
							"node_storage": "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version":       "4.4",
						"vswitch_id":           CHECKSET,
						"mongo_list.#":         "2",
						"shard_list.#":         "2",
						"config_server_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version": "5.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version": "5.0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_type":     "cloud_auto",
					"provisioned_iops": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_type":     "cloud_auto",
						"provisioned_iops": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secondary_zone_id": "${data.alicloud_mongodb_zones.default.zones.1.id}",
					"hidden_zone_id":    "${data.alicloud_mongodb_zones.default.zones.2.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secondary_zone_id": CHECKSET,
						"hidden_zone_id":    CHECKSET,
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
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
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
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_time": "11:00Z-12:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time": "11:00Z-12:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_period": []string{"Monday", "Tuesday", "Wednesday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_retention_period": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_retention_period": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_retention_policy_on_cluster_deletion": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_retention_policy_on_cluster_deletion": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup_retention_period": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup_retention_period": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_backup_log":           "1",
					"log_backup_retention_period": "120",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_backup_log":           "1",
						"log_backup_retention_period": "120",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"snapshot_backup_type": "Flash",
					"backup_interval":      "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snapshot_backup_type": "Flash",
						"backup_interval":      "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_action": "Open",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_status": "Open",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_action": "Close",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_status": "Closed",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_start_time": "00:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_start_time": "00:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_end_time": "03:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_end_time": "03:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_release_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_release_protection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_release_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_release_protection": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_security_group_list": "${alicloud_mongodb_global_security_ip_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_security_group_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_security_group_list": []string{"${alicloud_mongodb_global_security_ip_group.default.0.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_security_group_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_security_group_list": "${alicloud_mongodb_global_security_ip_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_security_group_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ShardingInstance",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "ShardingInstance",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mongo_list": []map[string]interface{}{
						{
							"node_class": "mdb.shard.8x.large.d",
						},
						{
							"node_class": "mdb.shard.8x.large.d",
						},
						{
							"node_class": "mdb.shard.8x.xlarge.d",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mongo_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"shard_list": []map[string]interface{}{
						{
							"node_class":   "mdb.shard.8x.large.d",
							"node_storage": "60",
						},
						{
							"node_class":        "mdb.shard.8x.xlarge.d",
							"node_storage":      "80",
							"readonly_replicas": "1",
						},
						// There is an api bug that does not support to update readonly_replicas
						//{
						//	"node_class":        "mdb.shard.8x.xlarge.d",
						//	"node_storage":      "50",
						//	"readonly_replicas": "2",
						//},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"shard_list.#": "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "kms_encrypted_password", "kms_encryption_context", "auto_renew", "ssl_action", "order_type"},
			},
		},
	})
}

func TestAccAliCloudMongoDBShardingInstance_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_sharding_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	ra := resourceAttrInit(resourceId, AliCloudMongoDBShardingInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBShardingInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBShardingInstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMongoDBShardingInstanceBasicDependence1)
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
					"engine_version":          "4.4",
					"storage_engine":          "WiredTiger",
					"storage_type":            "cloud_auto",
					"provisioned_iops":        "2000",
					"protocol_type":           "mongodb",
					"vpc_id":                  "${alicloud_vswitch.default.vpc_id}",
					"vswitch_id":              "${alicloud_vswitch.default.id}",
					"zone_id":                 "${data.alicloud_mongodb_zones.default.zones.0.id}",
					"secondary_zone_id":       "${data.alicloud_mongodb_zones.default.zones.1.id}",
					"hidden_zone_id":          "${data.alicloud_mongodb_zones.default.zones.2.id}",
					"security_group_id":       "${alicloud_security_group.default.id}",
					"network_type":            "VPC",
					"name":                    name,
					"instance_charge_type":    "PostPaid",
					"security_ip_list":        []string{"10.168.1.12"},
					"account_password":        "YourPassword_123",
					"resource_group_id":       "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"backup_time":             "11:00Z-12:00Z",
					"backup_period":           []string{"Monday", "Tuesday", "Wednesday"},
					"backup_retention_period": "7",
					"backup_retention_policy_on_cluster_deletion": "1",
					"enable_backup_log":                           "1",
					"log_backup_retention_period":                 "120",
					"snapshot_backup_type":                        "Flash",
					"backup_interval":                             "60",
					"ssl_action":                                  "Open",
					"maintain_start_time":                         "00:00Z",
					"maintain_end_time":                           "03:00Z",
					"db_instance_release_protection":              "false",
					"global_security_group_list":                  "${alicloud_mongodb_global_security_ip_group.default.*.id}",
					"mongo_list": []map[string]interface{}{
						{
							"node_class": "mdb.shard.8x.large.d",
						},
						{
							"node_class": "mdb.shard.8x.large.d",
						},
					},
					"shard_list": []map[string]interface{}{
						{
							"node_class":   "mdb.shard.8x.large.d",
							"node_storage": "500",
						},
						{
							"node_class":        "mdb.shard.8x.xlarge.d",
							"node_storage":      "510",
							"readonly_replicas": "1",
						},
					},
					"config_server_list": []map[string]interface{}{
						{
							"node_class":   "mdb.shard.2x.xlarge.d",
							"node_storage": "500",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ShardingInstance",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version":          "4.4",
						"storage_engine":          "WiredTiger",
						"storage_type":            "cloud_auto",
						"provisioned_iops":        "2000",
						"protocol_type":           "mongodb",
						"vpc_id":                  CHECKSET,
						"vswitch_id":              CHECKSET,
						"zone_id":                 CHECKSET,
						"secondary_zone_id":       CHECKSET,
						"hidden_zone_id":          CHECKSET,
						"security_group_id":       CHECKSET,
						"network_type":            "VPC",
						"name":                    name,
						"instance_charge_type":    "PostPaid",
						"security_ip_list.#":      "1",
						"account_password":        "YourPassword_123",
						"resource_group_id":       CHECKSET,
						"backup_time":             "11:00Z-12:00Z",
						"backup_period.#":         "3",
						"backup_retention_period": "7",
						"backup_retention_policy_on_cluster_deletion": "1",
						"enable_backup_log":                           "1",
						"log_backup_retention_period":                 "120",
						"snapshot_backup_type":                        "Flash",
						"backup_interval":                             "60",
						"ssl_status":                                  "Open",
						"maintain_start_time":                         "00:00Z",
						"maintain_end_time":                           "03:00Z",
						"db_instance_release_protection":              "false",
						"global_security_group_list.#":                "3",
						"mongo_list.#":                                "2",
						"shard_list.#":                                "2",
						"config_server_list.#":                        "1",
						"tags.%":                                      "2",
						"tags.Created":                                "TF",
						"tags.For":                                    "ShardingInstance",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "kms_encrypted_password", "kms_encryption_context", "auto_renew", "ssl_action", "order_type"},
			},
		},
	})
}

var AliCloudMongoDBShardingInstanceMap0 = map[string]string{
	"storage_engine":       CHECKSET,
	"storage_type":         CHECKSET,
	"protocol_type":        CHECKSET,
	"vpc_id":               CHECKSET,
	"vswitch_id":           CHECKSET,
	"zone_id":              CHECKSET,
	"network_type":         CHECKSET,
	"instance_charge_type": CHECKSET,
	"security_ip_list.#":   CHECKSET,
	"resource_group_id":    CHECKSET,
	"maintain_start_time":  CHECKSET,
	"maintain_end_time":    CHECKSET,
	"backup_time":          CHECKSET,
	"config_server_list.#": CHECKSET,
	"retention_period":     CHECKSET,
	"ssl_status":           CHECKSET,
}

func AliCloudMongoDBShardingInstanceBasicDependence0(name string) string {
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
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_mongodb_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
        security_group_name = var.name
	}

	resource "alicloud_mongodb_global_security_ip_group" "default" {
  		count                   = 3
  		global_ig_name          = "tfacc${count.index}"
  		global_security_ip_list = "192.168.1.1"
	}
`, name)
}

func AliCloudMongoDBShardingInstanceBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
  		status = "OK"
	}

	data "alicloud_mongodb_zones" "default" {
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = data.alicloud_mongodb_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_mongodb_global_security_ip_group" "default" {
  		count                   = 3
  		global_ig_name          = "tfacc${count.index}"
  		global_security_ip_list = "192.168.1.1"
	}
`, name)
}
