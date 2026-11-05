package mongodb_test

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/mongodb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_mongodb_sharding_instance", &resource.Sweeper{
		Name: "alicloud_mongodb_sharding_instance",
		F:    testSweepMongoDBShardingInstances,
	})
}

func testSweepMongoDBShardingInstances(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
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
	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1
	request["DBInstanceType"] = "sharding"
	request["ChargeType"] = "PostPaid"

	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
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
			if !sweep.SweepAll() {
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
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAliCloudMongoDBShardingInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_sharding_instance.default"
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region(mongodbTestRegion())})
	serverFunc := func() interface{} {
		return mongodb.NewMongoDBService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	ra := acctest.ResourceAttrInit(resourceId, AliCloudMongoDBShardingInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBShardingInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBShardingInstance%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudMongoDBShardingInstanceBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version": "4.2",
					"vswitch_id":     "${alicloud_vswitch.default.id}",
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
						"vswitch_id":     acctest.CHECKSET,
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
						"security_group_id": acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
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
			// There is an OpenAPI bug
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"encryptor_name": "aes-256-cbc",
			//		"encryption_key": "${alicloud_kms_key.default.id}",
			//		"role_arn":       "acs:ram::" + os.Getenv("ALICLOUD_ACCOUNT_ID") + ":role/aliyunrdsinstanceencryptiondefaultrole",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"encryptor_name": "aes-256-cbc",
			//			"encryption_key": CHECKSET,
			//			"role_arn":       CHECKSET,
			//		}),
			//	),
			//},
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
					"parameters": []interface{}{
						map[string]interface{}{
							"name":  "operationProfiling.slowOpThresholdMs",
							"value": "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "1",
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
				ImportStateVerifyIgnore: []string{"account_password", "kms_encrypted_password", "kms_encryption_context", "ssl_action", "order_type", "parameters"},
			},
		},
	})
}

func TestAccAliCloudMongoDBShardingInstance_securityIpGroups(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_sharding_instance.default"
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region(mongodbTestRegion())})
	serverFunc := func() interface{} {
		return mongodb.NewMongoDBService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	ra := acctest.ResourceAttrInit(resourceId, AliCloudMongoDBShardingInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBShardingInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBShardingInstanceSecurityIpGroups%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudMongoDBShardingInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version":              "4.2",
					"vswitch_id":                  "${alicloud_vswitch.default.id}",
					"account_password_wo":         "YourPassword123!",
					"account_password_wo_version": 1,
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
					"security_ip_groups": []map[string]interface{}{
						{
							"security_ip_group_attribute": "test",
							"security_ip_group_name":      "test",
							"security_ip_list":            "192.168.0.1",
						},
						{
							"security_ip_group_attribute": "test1",
							"security_ip_group_name":      "test1",
							"security_ip_list":            "192.168.0.2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_groups.#":        "2",
						"account_password_wo_version": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_groups": []map[string]interface{}{
						{
							"security_ip_group_attribute": "test3",
							"security_ip_group_name":      "test3",
							"security_ip_list":            "192.168.0.3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_groups.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_groups": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_groups.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password_wo":         "YourPassword123!Update",
					"account_password_wo_version": 2,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password_wo_version": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password":            "YourPassword123!password",
					"account_password_wo":         acctest.REMOVEKEY,
					"account_password_wo_version": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password":            "YourPassword123!password",
						"account_password_wo_version": acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "account_password_wo", "account_password_wo_version", "kms_encrypted_password", "kms_encryption_context", "ssl_action", "order_type", "parameters"},
			},
		},
	})
}

func TestAccAliCloudMongoDBShardingInstance_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_sharding_instance.default"
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region(mongodbTestRegion())})
	serverFunc := func() interface{} {
		return mongodb.NewMongoDBService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	ra := acctest.ResourceAttrInit(resourceId, AliCloudMongoDBShardingInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBShardingInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBShardingInstance%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudMongoDBShardingInstanceBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version":             "4.2",
					"storage_engine":             "WiredTiger",
					"storage_type":               "local_ssd",
					"protocol_type":              "mongodb",
					"vpc_id":                     "${alicloud_vswitch.default.vpc_id}",
					"vswitch_id":                 "${alicloud_vswitch.default.id}",
					"zone_id":                    "${alicloud_vswitch.default.zone_id}",
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
					"parameters": []interface{}{
						map[string]interface{}{
							"name":  "operationProfiling.slowOpThresholdMs",
							"value": "80",
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
						"vpc_id":                       acctest.CHECKSET,
						"vswitch_id":                   acctest.CHECKSET,
						"zone_id":                      acctest.CHECKSET,
						"security_group_id":            acctest.CHECKSET,
						"network_type":                 "VPC",
						"name":                         name,
						"instance_charge_type":         "PostPaid",
						"security_ip_list.#":           "1",
						"account_password":             "YourPassword_123",
						"resource_group_id":            acctest.CHECKSET,
						"backup_time":                  "11:00Z-12:00Z",
						"backup_period.#":              "3",
						"tde_status":                   "enabled",
						"global_security_group_list.#": "3",
						"parameters.#":                 "1",
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
				ImportStateVerifyIgnore: []string{"account_password", "kms_encrypted_password", "kms_encryption_context", "ssl_action", "order_type", "parameters"},
			},
		},
	})
}

func TestAccAliCloudMongoDBShardingInstance_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_sharding_instance.default"
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region(mongodbTestRegion())})
	serverFunc := func() interface{} {
		return mongodb.NewMongoDBService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	ra := acctest.ResourceAttrInit(resourceId, AliCloudMongoDBShardingInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBShardingInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBShardingInstance%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudMongoDBShardingInstanceBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version": "7.0",
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
							"node_storage": "60",
						},
						{
							"node_class":        "mdb.shard.8x.large.d",
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
						"engine_version":       "7.0",
						"vswitch_id":           acctest.CHECKSET,
						"mongo_list.#":         "2",
						"shard_list.#":         "2",
						"config_server_list.#": "1",
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
					"engine_version": "8.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version": "8.0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secondary_zone_id": mongodbTestZoneSecondary(),
					"hidden_zone_id":    mongodbTestZoneHidden(),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secondary_zone_id": acctest.CHECKSET,
						"hidden_zone_id":    acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_id": "${alicloud_security_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id": acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
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
					"ssl_action":       "Update",
					"force_encryption": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"force_encryption": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"force_encryption": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"force_encryption": "0",
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
			// There is an OpenAPI bug
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"maintain_start_time": "00:00Z",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"maintain_start_time": "00:00Z",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"maintain_end_time": "03:00Z",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"maintain_end_time": "03:00Z",
			//		}),
			//	),
			//},
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
					"parameters": []interface{}{
						map[string]interface{}{
							"name":  "operationProfiling.slowOpThresholdMs",
							"value": "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "1",
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
							"node_class": "mdb.shard.2x.2xlarge.d",
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
							"node_class":        "mdb.shard.8x.xlarge.d",
							"node_storage":      "100",
							"readonly_replicas": "2",
						},
						{
							"node_class":        "mdb.shard.8x.xlarge.d",
							"node_storage":      "100",
							"readonly_replicas": "2",
						},
						// There is an api bug that does not support to update readonly_replicas
						//{
						//	"node_class":        "mdb.shard.8x.xlarge.d",
						//	"node_storage":      "100",
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
				// auto-renewal, KMS-encrypted passwords, order changes and
				// the clone source (src_db_instance_id/restore_time/
				// role_arn) need a dedicated KMS/RDS setup that this suite
				// cannot provision; removing the required engine_version
				// stops the plan before the placeholders reach the API.
				Config: testAccConfig(map[string]interface{}{
					"auto_renew":             "true",
					"auto_renew_duration":    "1",
					"encryption_key":         "key-tfacc-placeholder",
					"encryptor_name":         "aes-256-cbc",
					"instance_charge_type":   "PrePaid",
					"kms_encrypted_password": "kms-tfacc-placeholder",
					"kms_encryption_context": map[string]string{"for": "tfacc"},
					"maintain_start_time":    "01:00Z",
					"maintain_end_time":      "02:00Z",
					"order_type":             "UPGRADE",
					"period":                 "1",
					"restore_time":           "2022-01-01T12:05:05Z",
					"role_arn":               "acs:ram::1234567890123456:role/tfacc-placeholder",
					"src_db_instance_id":     "rm-tfacc-placeholder",
					"engine_version":         acctest.REMOVEKEY,
				}),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "kms_encrypted_password", "kms_encryption_context", "ssl_action", "order_type", "parameters"},
			},
		},
	})
}

func TestAccAliCloudMongoDBShardingInstance_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_sharding_instance.default"
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region(mongodbTestRegion())})
	serverFunc := func() interface{} {
		return mongodb.NewMongoDBService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	ra := acctest.ResourceAttrInit(resourceId, AliCloudMongoDBShardingInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBShardingInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBShardingInstance%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudMongoDBShardingInstanceBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version":            "8.0",
					"storage_engine":            "WiredTiger",
					"storage_type":              "cloud_auto",
					"provisioned_iops":          "2000",
					"protocol_type":             "mongodb",
					"vpc_id":                    "${alicloud_vswitch.default.vpc_id}",
					"vswitch_id":                "${alicloud_vswitch.default.id}",
					"zone_id":                   "${alicloud_vswitch.default.zone_id}",
					"secondary_zone_id":         mongodbTestZoneSecondary(),
					"hidden_zone_id":            mongodbTestZoneHidden(),
					"security_group_id":         "${alicloud_security_group.default.id}",
					"network_type":              "VPC",
					"name":                      name,
					"instance_charge_type":      "PostPaid",
					"security_ip_list":          []string{"10.168.1.12"},
					"encrypted":                 "true",
					"cloud_disk_encryption_key": "${alicloud_kms_key.default.id}",
					"account_password":          "YourPassword_123",
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					// Currently, src_db_instance_id and restore_time cannot be CI tested, local testing has passed
					//"src_db_instance_id":      "",
					//"restore_time":            "",
					"backup_time":             "11:00Z-12:00Z",
					"backup_period":           []string{"Monday", "Tuesday", "Wednesday"},
					"backup_retention_period": "7",
					"backup_retention_policy_on_cluster_deletion": "1",
					"enable_backup_log":                           "1",
					"log_backup_retention_period":                 "120",
					"snapshot_backup_type":                        "Flash",
					"backup_interval":                             "60",
					"ssl_action":                                  "Open",
					"force_encryption":                            "1",
					"maintain_start_time":                         "00:00Z",
					"maintain_end_time":                           "03:00Z",
					"db_instance_release_protection":              "false",
					"global_security_group_list":                  "${alicloud_mongodb_global_security_ip_group.default.*.id}",
					"parameters": []interface{}{
						map[string]interface{}{
							"name":  "operationProfiling.slowOpThresholdMs",
							"value": "80",
						},
					},
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
							"node_class":        "mdb.shard.8x.large.d",
							"node_storage":      "500",
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
						"engine_version":            "8.0",
						"storage_engine":            "WiredTiger",
						"storage_type":              "cloud_auto",
						"provisioned_iops":          "2000",
						"protocol_type":             "mongodb",
						"vpc_id":                    acctest.CHECKSET,
						"vswitch_id":                acctest.CHECKSET,
						"zone_id":                   acctest.CHECKSET,
						"secondary_zone_id":         acctest.CHECKSET,
						"hidden_zone_id":            acctest.CHECKSET,
						"security_group_id":         acctest.CHECKSET,
						"network_type":              "VPC",
						"name":                      name,
						"instance_charge_type":      "PostPaid",
						"security_ip_list.#":        "1",
						"encrypted":                 "true",
						"cloud_disk_encryption_key": acctest.CHECKSET,
						"account_password":          "YourPassword_123",
						"resource_group_id":         acctest.CHECKSET,
						"backup_time":               "11:00Z-12:00Z",
						"backup_period.#":           "3",
						"backup_retention_period":   "7",
						"backup_retention_policy_on_cluster_deletion": "1",
						"enable_backup_log":                           "1",
						"log_backup_retention_period":                 "120",
						"snapshot_backup_type":                        "Flash",
						"backup_interval":                             "60",
						"force_encryption":                            "1",
						"ssl_status":                                  "Open",
						"maintain_start_time":                         "00:00Z",
						"maintain_end_time":                           "03:00Z",
						"db_instance_release_protection":              "false",
						"global_security_group_list.#":                "3",
						"parameters.#":                                "1",
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
				ImportStateVerifyIgnore: []string{"account_password", "kms_encrypted_password", "kms_encryption_context", "ssl_action", "order_type", "parameters"},
			},
		},
	})
}

var AliCloudMongoDBShardingInstanceMap0 = map[string]string{
	"storage_engine":       acctest.CHECKSET,
	"storage_type":         acctest.CHECKSET,
	"protocol_type":        acctest.CHECKSET,
	"vpc_id":               acctest.CHECKSET,
	"vswitch_id":           acctest.CHECKSET,
	"zone_id":              acctest.CHECKSET,
	"network_type":         acctest.CHECKSET,
	"instance_charge_type": acctest.CHECKSET,
	"security_ip_list.#":   acctest.CHECKSET,
	"resource_group_id":    acctest.CHECKSET,
	"maintain_start_time":  acctest.CHECKSET,
	"maintain_end_time":    acctest.CHECKSET,
	"backup_time":          acctest.CHECKSET,
	"config_server_list.#": acctest.CHECKSET,
	"retention_period":     acctest.CHECKSET,
	"ssl_status":           acctest.CHECKSET,
	"zone_infos.#":         acctest.CHECKSET,
}

func AliCloudMongoDBShardingInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
  		status = "OK"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "172.16.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "172.16.0.0/24"
  		zone_id      = "%s"
	}

	resource "alicloud_security_group" "default" {
  		vpc_id = alicloud_vpc.default.id
        security_group_name = var.name
	}

	resource "alicloud_kms_key" "default" {
  		description            = var.name
  		pending_window_in_days = 7
  		key_state              = "Enabled"
	}

	resource "alicloud_mongodb_global_security_ip_group" "default" {
  		count                   = 3
  		global_ig_name          = "tfacc${count.index}"
  		global_security_ip_list = "192.168.1.1"
	}
`, name, mongodbTestZonePrimary())
}

func AliCloudMongoDBShardingInstanceBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
  		status = "OK"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = "%s"
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_kms_key" "default" {
  		description            = var.name
  		pending_window_in_days = 7
  		key_state              = "Enabled"
	}

	resource "alicloud_mongodb_global_security_ip_group" "default" {
  		count                   = 3
  		global_ig_name          = "tfacc${count.index}"
  		global_security_ip_list = "192.168.1.1"
	}
`, name, mongodbTestZonePrimary())
}

// TestAccAliCloudMongoDBShardingInstance_mongoListReorder exercises a convergent
// reorder of the multi-member TypeList mongo_list: apply the original member
// order, then apply a permutation of the same members and require the
// post-apply plan to be empty. See scripts/collection-order/README.md.
func TestAccAliCloudMongoDBShardingInstance_mongoListReorder(t *testing.T) {
	resourceId := "alicloud_mongodb_sharding_instance.default"
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBShardingInstance%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudMongoDBShardingInstanceBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version": "4.2",
					"vswitch_id":     "${alicloud_vswitch.default.id}",
					"mongo_list": []map[string]interface{}{
						{"node_class": "dds.mongos.mid"},
						{"node_class": "dds.mongos.standard"},
					},
					"shard_list": []map[string]interface{}{
						{"node_class": "dds.shard.mid", "node_storage": "10"},
						{"node_class": "dds.shard.standard", "node_storage": "20"},
					},
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mongo_list": []map[string]interface{}{
						{"node_class": "dds.mongos.standard"},
						{"node_class": "dds.mongos.mid"},
					},
				}),
			},
		},
	})
}

// TestAccAliCloudMongoDBShardingInstance_shardListReorder exercises a convergent
// reorder of the multi-member TypeList shard_list: apply the original member
// order, then apply a permutation of the same members and require the
// post-apply plan to be empty. See scripts/collection-order/README.md.
func TestAccAliCloudMongoDBShardingInstance_shardListReorder(t *testing.T) {
	resourceId := "alicloud_mongodb_sharding_instance.default"
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBShardingInstance%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudMongoDBShardingInstanceBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version": "4.2",
					"vswitch_id":     "${alicloud_vswitch.default.id}",
					"mongo_list": []map[string]interface{}{
						{"node_class": "dds.mongos.mid"},
						{"node_class": "dds.mongos.standard"},
					},
					"shard_list": []map[string]interface{}{
						{"node_class": "dds.shard.mid", "node_storage": "10"},
						{"node_class": "dds.shard.standard", "node_storage": "20"},
					},
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"shard_list": []map[string]interface{}{
						{"node_class": "dds.shard.standard", "node_storage": "20"},
						{"node_class": "dds.shard.mid", "node_storage": "10"},
					},
				}),
			},
		},
	})
}
