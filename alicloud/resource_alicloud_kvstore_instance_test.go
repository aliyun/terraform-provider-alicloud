package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"strings"
	"time"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_kvstore_instance", &resource.Sweeper{
		Name: "alicloud_kvstore_instance",
		F:    testSweepKVStoreInstances,
	})
}

func testSweepKVStoreInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	conn, err := client.NewRedisaClient()
	if err != nil {
		return WrapError(err)
	}

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeInstances"
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeXLarge,
		"PageNumber": 1,
	}

	kvstoreInstanceIds := make([]string, 0)
	var response map[string]interface{}
	for _, instanceType := range []string{string(KVStoreRedis), string(KVStoreMemcache)} {
		request["InstanceType"] = instanceType
		for {
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				log.Printf("[ERROR] Failed to retrieve VPC in service list: %s", err)
				return nil
			}
			resp, err := jsonpath.Get("$.Instances.KVStoreInstance", response)
			if err != nil {
				return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Instances.KVStoreInstance", response)
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				skip := true
				item := v.(map[string]interface{})
				if !sweepAll() {
					for _, prefix := range prefixes {
						if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["InstanceName"])), strings.ToLower(prefix)) {
							skip = false
							break
						}
					}
					if skip {
						log.Printf("[INFO] Skipping KVStore Instance: %v (%v)", item["InstanceName"], item["InstanceId"])
						continue
					}
				}
				kvstoreInstanceIds = append(kvstoreInstanceIds, fmt.Sprint(item["InstanceId"]))
			}
			if len(result) < PageSizeXLarge {
				break
			}
			request["PageNumber"] = request["PageNumber"].(int) + 1
		}
	}

	for _, id := range kvstoreInstanceIds {
		log.Printf("[INFO] Deleting KVStore Instance: %s", id)
		action := "ModifyInstanceAttribute"
		request := map[string]interface{}{
			"InstanceId":                id,
			"InstanceReleaseProtection": false,
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			log.Printf("[ERROR] Failed to modify KVStore Instance release protection (%s): %s", id, err)
		}
		action = "DeleteInstance"
		request = map[string]interface{}{
			"InstanceId": id,
		}
		wait = incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			log.Printf("[ERROR] Failed to delete KVStore Instance (%s): %s", id, err)
		}
	}
	return nil
}

func TestAccAliCloudKVStoreRedisInstance_vpctest(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, RedisDbInstanceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreRedisInstanceVpcTest%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, KvstoreInstanceVpcTestdependence)
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
					"instance_class":   "redis.master.small.default",
					"db_instance_name": name,
					"instance_type":    "Redis",
					"engine_version":   "4.0",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 2].id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":    "redis.master.small.default",
						"db_instance_name":  name,
						"instance_type":     "Redis",
						"engine_version":    "4.0",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "acceptance test",
						"resource_group_id": CHECKSET,
						"zone_id":           CHECKSET,
						"vswitch_id":        CHECKSET,
						"secondary_zone_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "auto_use_coupon", "backup_id", "business_info", "coupon_no", "dedicated_host_group_id", "effective_time", "force_upgrade", "global_instance", "global_instance_id", "order_type", "password", "period", "restore_time", "src_db_instance_id", "enable_public", "security_ip_group_attribute", "enable_backup_log"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "kvstore",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "kvstore",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config": map[string]string{
						"appendonly":             "no",
						"lazyfree-lazy-eviction": "no",
						"EvictionPolicy":         "volatile-lru",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config.%":                      "3",
						"config.appendonly":             "no",
						"config.lazyfree-lazy-eviction": "no",
						"config.EvictionPolicy":         "volatile-lru",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_group_name": "tf",
					"security_ips":           []string{"10.23.12.24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_group_name": "tf",
						"security_ips.#":         "1",
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
					"db_instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class": "redis.master.mid.default",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class": "redis.master.mid.default",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 2].id}",
					"vswitch_id":        "${data.alicloud_vswitches.update.ids.0}",
					"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id}",
					"timeouts": []map[string]interface{}{
						{
							"update": "1h",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":           CHECKSET,
						"vswitch_id":        CHECKSET,
						"secondary_zone_id": CHECKSET,
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
					"backup_period": []string{"Tuesday", "Wednesday"},
					"backup_time":   "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "2",
						"backup_time":     "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_connection_prefix": fmt.Sprintf("privateprefix%d", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_connection_prefix": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_connection_port": "4010",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_connection_port": "4010",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_release_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_release_protection": "true",
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
					"instance_class":              "redis.master.small.default",
					"instance_release_protection": "false",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"security_ips":                []string{"10.0.0.1"},
					"db_instance_name":            name,
					"vpc_auth_mode":               "Open",
					"config": map[string]string{
						"appendonly":             "yes",
						"lazyfree-lazy-eviction": "yes",
						"EvictionPolicy":         "volatile-lru",
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"zone_id":                   "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id}",
					"vswitch_id":                "${data.alicloud_vswitches.default.ids.0}",
					"secondary_zone_id":         REMOVEKEY,
					"maintain_start_time":       "04:00Z",
					"maintain_end_time":         "06:00Z",
					"backup_period":             []string{"Wednesday"},
					"backup_time":               "11:00Z-12:00Z",
					"private_connection_prefix": fmt.Sprintf("privateprefixupdate%d", rand),
					"private_connection_port":   "4011",
					"timeouts": []map[string]interface{}{
						{
							"update": "1h",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":                "redis.master.small.default",
						"instance_release_protection":   "false",
						"resource_group_id":             CHECKSET,
						"security_ips.#":                "1",
						"db_instance_name":              name,
						"vpc_auth_mode":                 "Open",
						"config.%":                      "3",
						"config.appendonly":             "yes",
						"config.lazyfree-lazy-eviction": "yes",
						"config.EvictionPolicy":         "volatile-lru",
						"tags.%":                        "2",
						"tags.Created":                  "TF",
						"tags.For":                      "acceptance test",
						"zone_id":                       CHECKSET,
						"vswitch_id":                    CHECKSET,
						"secondary_zone_id":             REMOVEKEY,
						"maintain_start_time":           "04:00Z",
						"maintain_end_time":             "06:00Z",
						"backup_period.#":               "1",
						"backup_time":                   "11:00Z-12:00Z",
						"private_connection_port":       "4011",
						"private_connection_prefix":     CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudKVStoreRedisInstance_6_0(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, RedisDbInstanceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreRedisInstance6_0-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, KvstoreInstanceVpcTestdependence)
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
					"instance_class":   "redis.shard.with.proxy.small.ce",
					"db_instance_name": name,
					"instance_type":    "Redis",
					"engine_version":   "6.0",
					"shard_count":      "2",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 2].id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":    "redis.shard.with.proxy.small.ce",
						"db_instance_name":  name,
						"instance_type":     "Redis",
						"engine_version":    "6.0",
						"shard_count":       "2",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "acceptance test",
						"resource_group_id": CHECKSET,
						"zone_id":           CHECKSET,
						"vswitch_id":        CHECKSET,
						"secondary_zone_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"shard_count": "4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"shard_count": "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"shard_count": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"shard_count": "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "auto_use_coupon", "backup_id", "business_info", "coupon_no", "dedicated_host_group_id", "effective_time", "force_upgrade", "global_instance", "global_instance_id", "order_type", "password", "period", "restore_time", "src_db_instance_id", "enable_public", "security_ip_group_attribute", "enable_backup_log"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "kvstore",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "kvstore",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config": map[string]string{
						"appendonly":             "no",
						"lazyfree-lazy-eviction": "no",
						"maxmemory-policy":       "volatile-lru",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config.%":                      "3",
						"config.appendonly":             "no",
						"config.lazyfree-lazy-eviction": "no",
						"config.maxmemory-policy":       "volatile-lru",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.23.12.24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips.#": "1",
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
					"db_instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class": "redis.shard.with.proxy.mid.ce",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class": "redis.shard.with.proxy.mid.ce",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 2].id}",
					"vswitch_id":        "${data.alicloud_vswitches.update.ids.0}",
					"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id}",
					"timeouts": []map[string]interface{}{
						{
							"update": "1h",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":           CHECKSET,
						"vswitch_id":        CHECKSET,
						"secondary_zone_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enable": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enable": "Enable",
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
					"backup_period": []string{"Tuesday", "Wednesday"},
					"backup_time":   "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "2",
						"backup_time":     "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_connection_prefix": fmt.Sprintf("privateprefix%d", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_connection_prefix": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enable": "Disable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enable": "Disable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_connection_port": "4010",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_connection_port": "4010",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_release_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_release_protection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class":              "redis.shard.with.proxy.small.ce",
					"instance_release_protection": "false",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"security_ips":                []string{"10.0.0.1"},
					"db_instance_name":            name,
					"vpc_auth_mode":               "Open",
					"config": map[string]string{
						"appendonly":             "yes",
						"lazyfree-lazy-eviction": "yes",
						"maxmemory-policy":       "volatile-lru",
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"zone_id":                   "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id}",
					"vswitch_id":                "${data.alicloud_vswitches.default.ids.0}",
					"secondary_zone_id":         REMOVEKEY,
					"maintain_start_time":       "04:00Z",
					"maintain_end_time":         "06:00Z",
					"backup_period":             []string{"Wednesday"},
					"backup_time":               "11:00Z-12:00Z",
					"private_connection_prefix": fmt.Sprintf("privateprefixupdate%d", rand),
					"private_connection_port":   "4011",
					"timeouts": []map[string]interface{}{
						{
							"update": "1h",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":                "redis.shard.with.proxy.small.ce",
						"instance_release_protection":   "false",
						"resource_group_id":             CHECKSET,
						"security_ips.#":                "1",
						"db_instance_name":              name,
						"vpc_auth_mode":                 "Open",
						"config.%":                      "3",
						"config.appendonly":             "yes",
						"config.lazyfree-lazy-eviction": "yes",
						"config.maxmemory-policy":       "volatile-lru",
						"tags.%":                        "2",
						"tags.Created":                  "TF",
						"tags.For":                      "acceptance test",
						"zone_id":                       CHECKSET,
						"vswitch_id":                    CHECKSET,
						"secondary_zone_id":             REMOVEKEY,
						"maintain_start_time":           "04:00Z",
						"maintain_end_time":             "06:00Z",
						"backup_period.#":               "1",
						"backup_time":                   "11:00Z-12:00Z",
						"private_connection_port":       "4011",
						"private_connection_prefix":     CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudKVStoreRedisInstance_7_0(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, RedisDbInstanceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreRedisInstance7_0-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, KvstoreInstanceVpcTestdependence)
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
					"instance_class":   "redis.shard.small.ce",
					"db_instance_name": name,
					"instance_type":    "Redis",
					"engine_version":   "7.0",
					"shard_count":      "2",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 2].id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":    "redis.shard.small.ce",
						"db_instance_name":  name,
						"instance_type":     "Redis",
						"engine_version":    "7.0",
						"shard_count":       "2",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "acceptance test",
						"resource_group_id": CHECKSET,
						"zone_id":           CHECKSET,
						"vswitch_id":        CHECKSET,
						"secondary_zone_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "auto_use_coupon", "backup_id", "business_info", "coupon_no", "dedicated_host_group_id", "effective_time", "force_upgrade", "global_instance", "global_instance_id", "order_type", "password", "period", "restore_time", "src_db_instance_id", "enable_public", "security_ip_group_attribute", "enable_backup_log"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "kvstore",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "kvstore",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config": map[string]string{
						"appendonly":             "no",
						"lazyfree-lazy-eviction": "no",
						"maxmemory-policy":       "volatile-lru",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config.%":                      "3",
						"config.appendonly":             "no",
						"config.lazyfree-lazy-eviction": "no",
						"config.maxmemory-policy":       "volatile-lru",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.23.12.24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips.#": "1",
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
					"db_instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_name": name + "_update",
					}),
				),
			},
			// there is no more quota for this class on multi-zone
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"zone_id":           "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 2].id}",
			//		"vswitch_id":        "${data.alicloud_vswitches.update.ids.0}",
			//		"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id}",
			//		"timeouts": []map[string]interface{}{
			//			{
			//				"update": "1h",
			//			},
			//		},
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"zone_id":           CHECKSET,
			//			"vswitch_id":        CHECKSET,
			//			"secondary_zone_id": CHECKSET,
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class": "redis.shard.mid.ce",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class": "redis.shard.mid.ce",
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
					"backup_period": []string{"Tuesday", "Wednesday"},
					"backup_time":   "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "2",
						"backup_time":     "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_connection_prefix": fmt.Sprintf("privateprefix%d", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_connection_prefix": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_connection_port": "4010",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_connection_port": "4010",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_release_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_release_protection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class":              "redis.shard.small.ce",
					"instance_release_protection": "false",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"security_ip_group_name":      "tf",
					"security_ips":                []string{"10.0.0.1"},
					"db_instance_name":            name,
					"vpc_auth_mode":               "Open",
					"config": map[string]string{
						"appendonly":             "yes",
						"lazyfree-lazy-eviction": "yes",
						"maxmemory-policy":       "volatile-lru",
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"maintain_start_time":       "04:00Z",
					"maintain_end_time":         "06:00Z",
					"backup_period":             []string{"Wednesday"},
					"backup_time":               "11:00Z-12:00Z",
					"private_connection_prefix": fmt.Sprintf("privateprefixupdate%d", rand),
					"private_connection_port":   "4011",
					"timeouts": []map[string]interface{}{
						{
							"update": "1h",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":                "redis.shard.small.ce",
						"instance_release_protection":   "false",
						"resource_group_id":             CHECKSET,
						"security_ip_group_name":        "tf",
						"security_ips.#":                "1",
						"db_instance_name":              name,
						"vpc_auth_mode":                 "Open",
						"config.%":                      "3",
						"config.appendonly":             "yes",
						"config.lazyfree-lazy-eviction": "yes",
						"config.maxmemory-policy":       "volatile-lru",
						"tags.%":                        "2",
						"tags.Created":                  "TF",
						"tags.For":                      "acceptance test",
						"maintain_start_time":           "04:00Z",
						"maintain_end_time":             "06:00Z",
						"backup_period.#":               "1",
						"backup_time":                   "11:00Z-12:00Z",
						"private_connection_port":       "4011",
						"private_connection_prefix":     CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudKVStoreRedisInstance_7_0_with_proxy_class(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, RedisDbInstanceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreRedisInstance7_0_with_proxy-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, KvstoreInstanceVpcTestdependence)
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
					"instance_class":   "redis.shard.with.proxy.small.ce",
					"db_instance_name": name,
					"instance_type":    "Redis",
					"engine_version":   "7.0",
					"shard_count":      "2",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 2].id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":    "redis.shard.with.proxy.small.ce",
						"db_instance_name":  name,
						"instance_type":     "Redis",
						"engine_version":    "7.0",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "acceptance test",
						"shard_count":       "2",
						"resource_group_id": CHECKSET,
						"zone_id":           CHECKSET,
						"vswitch_id":        CHECKSET,
						"secondary_zone_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "auto_use_coupon", "backup_id", "business_info", "coupon_no", "dedicated_host_group_id", "effective_time", "force_upgrade", "global_instance", "global_instance_id", "order_type", "password", "period", "restore_time", "src_db_instance_id", "enable_public", "security_ip_group_attribute", "enable_backup_log"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "kvstore",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "kvstore",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config": map[string]string{
						"appendonly":             "no",
						"lazyfree-lazy-eviction": "no",
						"maxmemory-policy":       "volatile-lru",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config.%":                      "3",
						"config.appendonly":             "no",
						"config.lazyfree-lazy-eviction": "no",
						"config.maxmemory-policy":       "volatile-lru",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.23.12.24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips.#": "1",
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
					"db_instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class": "redis.shard.with.proxy.mid.ce",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class": "redis.shard.with.proxy.mid.ce",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 2].id}",
					"vswitch_id":        "${data.alicloud_vswitches.update.ids.0}",
					"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id}",
					"timeouts": []map[string]interface{}{
						{
							"update": "1h",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":           CHECKSET,
						"vswitch_id":        CHECKSET,
						"secondary_zone_id": CHECKSET,
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
					"backup_period": []string{"Tuesday", "Wednesday"},
					"backup_time":   "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "2",
						"backup_time":     "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_connection_prefix": fmt.Sprintf("privateprefix%d", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_connection_prefix": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_connection_port": "4010",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_connection_port": "4010",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_release_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_release_protection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class":              "redis.shard.with.proxy.small.ce",
					"instance_release_protection": "false",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"security_ips":                []string{"10.0.0.1"},
					"db_instance_name":            name,
					"vpc_auth_mode":               "Open",
					"config": map[string]string{
						"appendonly":             "yes",
						"lazyfree-lazy-eviction": "yes",
						"maxmemory-policy":       "volatile-lru",
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"zone_id":                   "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id}",
					"vswitch_id":                "${data.alicloud_vswitches.default.ids.0}",
					"secondary_zone_id":         REMOVEKEY,
					"maintain_start_time":       "04:00Z",
					"maintain_end_time":         "06:00Z",
					"backup_period":             []string{"Wednesday"},
					"backup_time":               "11:00Z-12:00Z",
					"private_connection_prefix": fmt.Sprintf("privateprefixupdate%d", rand),
					"private_connection_port":   "4011",
					"timeouts": []map[string]interface{}{
						{
							"update": "1h",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":                "redis.shard.with.proxy.small.ce",
						"instance_release_protection":   "false",
						"resource_group_id":             CHECKSET,
						"security_ips.#":                "1",
						"db_instance_name":              name,
						"vpc_auth_mode":                 "Open",
						"config.%":                      "3",
						"config.appendonly":             "yes",
						"config.lazyfree-lazy-eviction": "yes",
						"config.maxmemory-policy":       "volatile-lru",
						"tags.%":                        "2",
						"tags.Created":                  "TF",
						"tags.For":                      "acceptance test",
						"zone_id":                       CHECKSET,
						"vswitch_id":                    CHECKSET,
						"secondary_zone_id":             REMOVEKEY,
						"maintain_start_time":           "04:00Z",
						"maintain_end_time":             "06:00Z",
						"backup_period.#":               "1",
						"backup_time":                   "11:00Z-12:00Z",
						"private_connection_port":       "4011",
						"private_connection_prefix":     CHECKSET,
					}),
				),
			},
		},
	})
}

// Memcache Instance has been offline since 2023-08-13
func SkipTestAccAliCloudKVStoreMemcacheInstance_vpctest(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, RedisDbInstanceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreMemcacheInstanceVpcTest%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, KvstoreMemcacheInstanceVpcTestdependence)
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
					"instance_class":   "memcache.master.small.default",
					"db_instance_name": name,
					"instance_type":    "Memcache",
					"engine_version":   "4.0",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 2].id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":    "memcache.master.small.default",
						"db_instance_name":  name,
						"instance_type":     "Memcache",
						"engine_version":    "4.0",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "acceptance test",
						"resource_group_id": CHECKSET,
						"zone_id":           CHECKSET,
						"vswitch_id":        CHECKSET,
						"secondary_zone_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "auto_use_coupon", "backup_id", "business_info", "coupon_no", "dedicated_host_group_id", "effective_time", "force_upgrade", "global_instance", "global_instance_id", "order_type", "password", "period", "restore_time", "src_db_instance_id", "enable_public", "security_ip_group_attribute", "enable_backup_log"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "kvstore",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "kvstore",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.23.12.24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips.#": "1",
					}),
				),
			},
			// Currently, the memcache only support version 4.0
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"engine_version": "2.8",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"engine_version": "2.8",
			//		}),
			//	),
			//},
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
					"db_instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class": "memcache.master.mid.default",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class": "memcache.master.mid.default",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 2].id}",
					"vswitch_id":        "${data.alicloud_vswitches.update.ids.0}",
					"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id}",
					"timeouts": []map[string]interface{}{
						{
							"update": "1h",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":           CHECKSET,
						"vswitch_id":        CHECKSET,
						"secondary_zone_id": CHECKSET,
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
					"backup_period": []string{"Tuesday", "Wednesday"},
					"backup_time":   "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "2",
						"backup_time":     "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_connection_prefix": fmt.Sprintf("privateconnection%d", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_connection_prefix": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_release_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_release_protection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class":              "memcache.master.small.default",
					"instance_release_protection": "false",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"security_ips":                []string{"10.0.0.1"},
					"db_instance_name":            name,
					"vpc_auth_mode":               "Open",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"zone_id":                   "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id}",
					"vswitch_id":                "${data.alicloud_vswitches.default.ids.0}",
					"secondary_zone_id":         REMOVEKEY,
					"maintain_start_time":       "04:00Z",
					"maintain_end_time":         "06:00Z",
					"backup_period":             []string{"Wednesday"},
					"backup_time":               "11:00Z-12:00Z",
					"private_connection_prefix": fmt.Sprintf("privateconnectionupdate%d", rand),
					"timeouts": []map[string]interface{}{
						{
							"update": "1h",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":              "memcache.master.small.default",
						"instance_release_protection": "false",
						"resource_group_id":           CHECKSET,
						"security_ips.#":              "1",
						"db_instance_name":            name,
						"vpc_auth_mode":               "Open",
						"tags.%":                      "2",
						"tags.Created":                "TF",
						"tags.For":                    "acceptance test",
						"zone_id":                     CHECKSET,
						"vswitch_id":                  CHECKSET,
						"secondary_zone_id":           REMOVEKEY,
						"maintain_start_time":         "04:00Z",
						"maintain_end_time":           "06:00Z",
						"backup_period.#":             "1",
						"backup_time":                 "11:00Z-12:00Z",
						"private_connection_prefix":   CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudKVStoreRedisInstance_prepaid(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, RedisDbInstanceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreRedisInstancePrePaid%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, KvstoreInstancePrePaidDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class":   "redis.master.small.default",
					"payment_type":     "PrePaid",
					"period":           "1",
					"db_instance_name": name,
					"instance_type":    "Redis",
					"engine_version":   "4.0",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones[0].id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones[1].id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":    "redis.master.small.default",
						"db_instance_name":  name,
						"instance_type":     "Redis",
						"engine_version":    "4.0",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "acceptance test",
						"resource_group_id": CHECKSET,
						"zone_id":           CHECKSET,
						"vswitch_id":        CHECKSET,
						"secondary_zone_id": CHECKSET,
						"payment_type":      "PrePaid",
						"auto_renew":        "false",
						"auto_renew_period": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "auto_use_coupon", "backup_id", "business_info", "coupon_no", "dedicated_host_group_id", "effective_time", "force_upgrade", "global_instance", "global_instance_id", "order_type", "password", "period", "restore_time", "src_db_instance_id", "enable_public", "security_ip_group_attribute", "enable_backup_log"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew":        "true",
						"auto_renew_period": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew_period": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew_period": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew":        "false",
						"auto_renew_period": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "kvstore",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "kvstore",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config": map[string]string{
						"appendonly":             "no",
						"lazyfree-lazy-eviction": "no",
						"maxmemory-policy":       "volatile-lru",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config.%":                      "3",
						"config.appendonly":             "no",
						"config.lazyfree-lazy-eviction": "no",
						"config.maxmemory-policy":       "volatile-lru",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.23.12.24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips.#": "1",
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
					"db_instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class": "redis.master.mid.default",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class": "redis.master.mid.default",
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
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones[1].id}",
					"vswitch_id":        "${data.alicloud_vswitches.update.ids.0}",
					"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones[0].id}",
					"timeouts": []map[string]interface{}{
						{
							"update": "1h",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":           CHECKSET,
						"vswitch_id":        CHECKSET,
						"secondary_zone_id": CHECKSET,
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
					"backup_period": []string{"Tuesday", "Wednesday"},
					"backup_time":   "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "2",
						"backup_time":     "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_connection_prefix": fmt.Sprintf("privateprefix%d", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_connection_prefix": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_connection_port": "4010",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_connection_port": "4010",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_release_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_release_protection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_release_protection": "false",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"security_ips":                []string{"10.0.0.1"},
					"db_instance_name":            name,
					"vpc_auth_mode":               "Open",
					"config": map[string]string{
						"appendonly":             "yes",
						"lazyfree-lazy-eviction": "yes",
						"maxmemory-policy":       "volatile-lru",
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"zone_id":                   "${data.alicloud_kvstore_zones.default.zones[0].id}",
					"vswitch_id":                "${data.alicloud_vswitches.default.ids.0}",
					"secondary_zone_id":         REMOVEKEY,
					"maintain_start_time":       "04:00Z",
					"maintain_end_time":         "06:00Z",
					"backup_period":             []string{"Wednesday"},
					"backup_time":               "11:00Z-12:00Z",
					"private_connection_prefix": fmt.Sprintf("privateprefixupdate%d", rand),
					"private_connection_port":   "4011",
					"timeouts": []map[string]interface{}{
						{
							"update": "1h",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_release_protection":   "false",
						"resource_group_id":             CHECKSET,
						"security_ips.#":                "1",
						"db_instance_name":              name,
						"vpc_auth_mode":                 "Open",
						"config.%":                      "3",
						"config.appendonly":             "yes",
						"config.lazyfree-lazy-eviction": "yes",
						"config.maxmemory-policy":       "volatile-lru",
						"tags.%":                        "2",
						"tags.Created":                  "TF",
						"tags.For":                      "acceptance test",
						"zone_id":                       CHECKSET,
						"vswitch_id":                    CHECKSET,
						"secondary_zone_id":             REMOVEKEY,
						"maintain_start_time":           "04:00Z",
						"maintain_end_time":             "06:00Z",
						"backup_period.#":               "1",
						"backup_time":                   "11:00Z-12:00Z",
						"private_connection_port":       "4011",
						"private_connection_prefix":     CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudKVStoreRedisInstance_5_0_memory_classic_standard(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	// en-central-1 has no enough quota for this class
	checkoutSupportedRegions(t, true, []connectivity.Region{connectivity.APSouthEast1, connectivity.Hangzhou})
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, RedisDbInstanceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreRedisInstanceVpcMultiTest%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, KvstoreInstanceVpcTestdependence)
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
					"instance_class":    "redis.amber.master.small.multithread",
					"db_instance_name":  name,
					"instance_type":     "Redis",
					"engine_version":    "5.0",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					//"shard_count":       "2",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":    "redis.amber.master.small.multithread",
						"db_instance_name":  name,
						"instance_type":     "Redis",
						"engine_version":    "5.0",
						"resource_group_id": CHECKSET,
						"zone_id":           CHECKSET,
						"vswitch_id":        CHECKSET,
						//"shard_count":       "2",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tde_status": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tde_status": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "kvstore",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "kvstore",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config": map[string]string{
						"appendonly":             "no",
						"lazyfree-lazy-eviction": "no",
						"maxmemory-policy":       "volatile-lru",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config.%":                      "3",
						"config.appendonly":             "no",
						"config.lazyfree-lazy-eviction": "no",
						"config.maxmemory-policy":       "volatile-lru",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.23.12.24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips.#": "1",
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
					"db_instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class": "redis.amber.master.mid.multithread",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class": "redis.amber.master.mid.multithread",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 2].id}",
					"vswitch_id":        "${data.alicloud_vswitches.update.ids.0}",
					"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id}",
					"timeouts": []map[string]interface{}{
						{
							"update": "1h",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":           CHECKSET,
						"vswitch_id":        CHECKSET,
						"secondary_zone_id": CHECKSET,
					}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"ssl_enable": "Enable",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"ssl_enable": "Enable",
			//		}),
			//	),
			//},
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
					"backup_period": []string{"Tuesday", "Wednesday"},
					"backup_time":   "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "2",
						"backup_time":     "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_connection_prefix": fmt.Sprintf("privateprefix%d", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_connection_prefix": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_connection_port": "4010",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_connection_port": "4010",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_release_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_release_protection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class":              "redis.amber.master.small.multithread",
					"instance_release_protection": "false",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"security_ips":                []string{"10.0.0.1"},
					"db_instance_name":            name,
					"vpc_auth_mode":               "Open",
					//"ssl_enable":                  "Disable",
					"config": map[string]string{
						"appendonly":             "yes",
						"lazyfree-lazy-eviction": "yes",
						"maxmemory-policy":       "volatile-lru",
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"zone_id":                   "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id}",
					"vswitch_id":                "${data.alicloud_vswitches.default.ids.0}",
					"secondary_zone_id":         REMOVEKEY,
					"maintain_start_time":       "04:00Z",
					"maintain_end_time":         "06:00Z",
					"backup_period":             []string{"Wednesday"},
					"backup_time":               "11:00Z-12:00Z",
					"private_connection_prefix": fmt.Sprintf("privateprefixupdate%d", rand),
					"private_connection_port":   "4011",
					"timeouts": []map[string]interface{}{
						{
							"update": "1h",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":              "redis.amber.master.small.multithread",
						"instance_release_protection": "false",
						"resource_group_id":           CHECKSET,
						"security_ips.#":              "1",
						"db_instance_name":            name,
						"vpc_auth_mode":               "Open",
						//"ssl_enable":                    "Disable",
						"config.%":                      "3",
						"config.appendonly":             "yes",
						"config.lazyfree-lazy-eviction": "yes",
						"config.maxmemory-policy":       "volatile-lru",
						"tags.%":                        "2",
						"tags.Created":                  "TF",
						"tags.For":                      "acceptance test",
						"zone_id":                       CHECKSET,
						"vswitch_id":                    CHECKSET,
						"secondary_zone_id":             REMOVEKEY,
						"maintain_start_time":           "04:00Z",
						"maintain_end_time":             "06:00Z",
						"backup_period.#":               "1",
						"backup_time":                   "11:00Z-12:00Z",
						"private_connection_port":       "4011",
						"private_connection_prefix":     CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudKVStoreRedisInstance_5_0_memory_classic_cluster(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	// en-central-1 has no enough quota for this class
	checkoutSupportedRegions(t, true, []connectivity.Region{connectivity.APSouthEast1, connectivity.Hangzhou})
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, RedisDbInstanceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreRedisInstanceVpcMultiTest%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, KvstoreInstanceVpcTestdependence)
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
					"instance_class":    "redis.amber.logic.sharding.1g.2db.0rodb.6proxy.multithread",
					"db_instance_name":  name,
					"instance_type":     "Redis",
					"engine_version":    "5.0",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"shard_count":       "2",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":    "redis.amber.logic.sharding.1g.2db.0rodb.6proxy.multithread",
						"db_instance_name":  name,
						"instance_type":     "Redis",
						"engine_version":    "5.0",
						"resource_group_id": CHECKSET,
						"zone_id":           CHECKSET,
						"vswitch_id":        CHECKSET,
						"shard_count":       "2",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "acceptance test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enable": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enable": "Enable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tde_status": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tde_status": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "kvstore",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "kvstore",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config": map[string]string{
						"appendonly":             "no",
						"lazyfree-lazy-eviction": "no",
						"maxmemory-policy":       "volatile-lru",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config.%":                      "3",
						"config.appendonly":             "no",
						"config.lazyfree-lazy-eviction": "no",
						"config.maxmemory-policy":       "volatile-lru",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.23.12.24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips.#": "1",
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
					"db_instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class": "redis.amber.logic.sharding.2g.2db.0rodb.6proxy.multithread",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class": "redis.amber.logic.sharding.2g.2db.0rodb.6proxy.multithread",
					}),
				),
			},
			// there is no more quota for this class on multi-zone
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"zone_id":           "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 2].id}",
			//		"vswitch_id":        "${data.alicloud_vswitches.update.ids.0}",
			//		"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id}",
			//		"timeouts": []map[string]interface{}{
			//			{
			//				"update": "1h",
			//			},
			//		},
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"zone_id":           CHECKSET,
			//			"vswitch_id":        CHECKSET,
			//			"secondary_zone_id": CHECKSET,
			//		}),
			//	),
			//},
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
					"backup_period": []string{"Tuesday", "Wednesday"},
					"backup_time":   "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "2",
						"backup_time":     "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_connection_prefix": fmt.Sprintf("privateprefix%d", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_connection_prefix": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_connection_port": "4010",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_connection_port": "4010",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_release_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_release_protection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class":              "redis.amber.logic.sharding.1g.2db.0rodb.6proxy.multithread",
					"instance_release_protection": "false",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"security_ips":                []string{"10.0.0.1"},
					"ssl_enable":                  "Disable",
					"db_instance_name":            name,
					"vpc_auth_mode":               "Open",
					"config": map[string]string{
						"appendonly":             "yes",
						"lazyfree-lazy-eviction": "yes",
						"maxmemory-policy":       "volatile-lru",
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"zone_id":                   "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id}",
					"vswitch_id":                "${data.alicloud_vswitches.default.ids.0}",
					"maintain_start_time":       "04:00Z",
					"maintain_end_time":         "06:00Z",
					"backup_period":             []string{"Wednesday"},
					"backup_time":               "11:00Z-12:00Z",
					"private_connection_prefix": fmt.Sprintf("privateprefixupdate%d", rand),
					"private_connection_port":   "4011",
					"timeouts": []map[string]interface{}{
						{
							"update": "1h",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":                "redis.amber.logic.sharding.1g.2db.0rodb.6proxy.multithread",
						"instance_release_protection":   "false",
						"resource_group_id":             CHECKSET,
						"security_ips.#":                "1",
						"db_instance_name":              name,
						"vpc_auth_mode":                 "Open",
						"ssl_enable":                    "Disable",
						"config.%":                      "3",
						"config.appendonly":             "yes",
						"config.lazyfree-lazy-eviction": "yes",
						"config.maxmemory-policy":       "volatile-lru",
						"tags.%":                        "2",
						"tags.Created":                  "TF",
						"tags.For":                      "acceptance test",
						"zone_id":                       CHECKSET,
						"vswitch_id":                    CHECKSET,
						"maintain_start_time":           "04:00Z",
						"maintain_end_time":             "06:00Z",
						"backup_period.#":               "1",
						"backup_time":                   "11:00Z-12:00Z",
						"private_connection_port":       "4011",
						"private_connection_prefix":     CHECKSET,
					}),
				),
			},
		},
	})
}

var RedisDbInstanceMap = map[string]string{
	"connection_domain": CHECKSET,
	"ssl_enable":        "Disable",
}

func KvstoreInstanceVpcTestdependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_kvstore_zones" "default"{
		instance_charge_type = "PostPaid"
	}
	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
	}
	data "alicloud_vswitches" "default" {
  		zone_id = data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}
	data "alicloud_resource_manager_resource_groups" "default"{
		status = "OK"
	}

	data "alicloud_vswitches" "update" {
  		zone_id = data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 2].id
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}
	resource "alicloud_security_group" "default" {
	 inner_access_policy = "Accept"
	 name                = "tf-example"
	 vpc_id              = data.alicloud_vpcs.default.ids.0
	}
	`)
}

func KvstoreInstancePrePaidDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_kvstore_zones" "default"{
		instance_charge_type = "PrePaid"
	}
	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
	}
	data "alicloud_vswitches" "default" {
  		zone_id = data.alicloud_kvstore_zones.default.zones[0].id
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}
	data "alicloud_resource_manager_resource_groups" "default"{
		status = "OK"
	}

	data "alicloud_vswitches" "update" {
  		zone_id = data.alicloud_kvstore_zones.default.zones[1].id
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}
	`)
}

func KvstoreMemcacheInstanceVpcTestdependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_kvstore_zones" "default"{
		instance_charge_type = "PostPaid"
		engine = "memcache"
	}
	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
	}
	data "alicloud_vswitches" "default" {
  		zone_id = data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}
	data "alicloud_resource_manager_resource_groups" "default"{
		status = "OK"
	}

	data "alicloud_vswitches" "update" {
  		zone_id = data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 2].id
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}
	`)
}
