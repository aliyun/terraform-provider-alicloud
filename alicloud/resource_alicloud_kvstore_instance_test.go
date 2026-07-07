package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
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
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, nil, request, true)
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
			_, err = client.RpcPost("R-kvstore", "2015-01-01", action, nil, request, false)
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
			_, err = client.RpcPost("R-kvstore", "2015-01-01", action, nil, request, false)
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

// engine_version 4.0 has been offline from July 31, 2025
func SkipTestAccAliCloudKVStoreRedisInstance_vpctest(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudKVStoreMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreRedisInstanceVpcTest%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKVStoreRedisInstanceVpcBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
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
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones.0.id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones.1.id}",
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
				ImportStateVerifyIgnore: []string{"dry_run", "business_info", "coupon_no", "effective_time", "force_upgrade", "global_instance_id", "order_type", "password", "period", "enable_public", "security_ip_group_attribute", "enable_backup_log", "port"},
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
					"instance_class": "redis.master.large.default",
					"engine_version": "5.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class": "redis.master.large.default",
						"engine_version": "5.0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones.1.id}",
					"vswitch_id":        "${data.alicloud_vswitches.update.ids.0}",
					"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones.0.id}",
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
					"payment_type":      "PrePaid",
					"period":            "1",
					"auto_renew":        "true",
					"auto_renew_period": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":      "PrePaid",
						"auto_renew":        "true",
						"auto_renew_period": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PostPaid",
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
						"maxmemory-policy":       "volatile-lru",
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"zone_id":             "${data.alicloud_kvstore_zones.default.zones.0.id}",
					"vswitch_id":          "${data.alicloud_vswitches.default.ids.0}",
					"secondary_zone_id":   REMOVEKEY,
					"maintain_start_time": "04:00Z",
					"maintain_end_time":   "06:00Z",
					// There is an OpenAPI bug
					//"backup_period":             []string{"Wednesday"},
					//"backup_time":               "11:00Z-12:00Z",
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
						"config.maxmemory-policy":       "volatile-lru",
						"tags.%":                        "2",
						"tags.Created":                  "TF",
						"tags.For":                      "acceptance test",
						"zone_id":                       CHECKSET,
						"vswitch_id":                    CHECKSET,
						"secondary_zone_id":             REMOVEKEY,
						"maintain_start_time":           "04:00Z",
						"maintain_end_time":             "06:00Z",
						//"backup_period.#":               "1",
						//"backup_time":                   "11:00Z-12:00Z",
						"private_connection_port":   "4011",
						"private_connection_prefix": CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudKVStoreRedisInstance_6_0(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudKVStoreMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreRedisInstance6_0-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKVStoreRedisInstanceVpcBasicDependence0)
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
					"payment_type":     "PostPaid",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones.0.id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":    "redis.shard.with.proxy.small.ce",
						"db_instance_name":  name,
						"instance_type":     "Redis",
						"engine_version":    "6.0",
						"shard_count":       "2",
						"payment_type":      "PostPaid",
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
				ImportStateVerifyIgnore: []string{"dry_run", "business_info", "coupon_no", "effective_time", "force_upgrade", "global_instance_id", "order_type", "password", "period", "enable_public", "security_ip_group_attribute", "enable_backup_log", "port"},
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
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones.1.id}",
					"vswitch_id":        "${data.alicloud_vswitches.update.ids.0}",
					"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones.0.id}",
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
			// There is an OpenAPI bug in eu-central-1
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"maintain_start_time": "02:00Z",
			//		"maintain_end_time":   "03:00Z",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"maintain_start_time": "02:00Z",
			//			"maintain_end_time":   "03:00Z",
			//		}),
			//	),
			//},
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
					"bandwidth": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "200",
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
					"payment_type":      "PrePaid",
					"period":            "1",
					"auto_renew":        "true",
					"auto_renew_period": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":      "PrePaid",
						"auto_renew":        "true",
						"auto_renew_period": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PostPaid",
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
					"bandwidth":                   "100",
					"config": map[string]string{
						"appendonly":             "yes",
						"lazyfree-lazy-eviction": "yes",
						"maxmemory-policy":       "volatile-lru",
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones.0.id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"secondary_zone_id": REMOVEKEY,
					// There is an OpenAPI bug in eu-central-1
					//"maintain_start_time": "04:00Z",
					//"maintain_end_time":   "06:00Z",
					// There is an OpenAPI bug
					//"backup_period":             []string{"Wednesday"},
					//"backup_time":               "11:00Z-12:00Z",
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
						"bandwidth":                     "100",
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
						//"maintain_start_time":           "04:00Z",
						//"maintain_end_time":             "06:00Z",
						//"backup_period.#":               "1",
						//"backup_time":                   "11:00Z-12:00Z",
						"private_connection_port":   "4011",
						"private_connection_prefix": CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudKVStoreRedisInstance_7_0(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudKVStoreMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreRedisInstance7_0-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKVStoreRedisInstanceVpcBasicDependence0)
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
					"bandwidth":        "200",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones.0.id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":    "redis.shard.small.ce",
						"db_instance_name":  name,
						"instance_type":     "Redis",
						"engine_version":    "7.0",
						"shard_count":       "2",
						"bandwidth":         "200",
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
				ImportStateVerifyIgnore: []string{"dry_run", "business_info", "coupon_no", "effective_time", "force_upgrade", "global_instance_id", "order_type", "password", "period", "enable_public", "security_ip_group_attribute", "enable_backup_log", "port"},
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
			//		"zone_id":           "${data.alicloud_kvstore_zones.default.zones.1.id}",
			//		"vswitch_id":        "${data.alicloud_vswitches.update.ids.0}",
			//		"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones.0.id}",
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
					"bandwidth":      "296",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class": "redis.shard.mid.ce",
						"bandwidth":      "296",
					}),
				),
			},
			// There is an OpenAPI bug in eu-central-1
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"maintain_start_time": "02:00Z",
			//		"maintain_end_time":   "03:00Z",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"maintain_start_time": "02:00Z",
			//			"maintain_end_time":   "03:00Z",
			//		}),
			//	),
			//},
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
					"instance_charge_type": "PrePaid",
					"period":               "1",
					"auto_renew":           "true",
					"auto_renew_period":    "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "PrePaid",
						"auto_renew":           "true",
						"auto_renew_period":    "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type": "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "PostPaid",
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
					"bandwidth":                   "200",
					"config": map[string]string{
						"appendonly":             "yes",
						"lazyfree-lazy-eviction": "yes",
						"maxmemory-policy":       "volatile-lru",
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					// There is an OpenAPI bug in eu-central-1
					//"maintain_start_time": "04:00Z",
					//"maintain_end_time":   "06:00Z",
					// There is an OpenAPI bug
					//"backup_period":             []string{"Wednesday"},
					//"backup_time":               "11:00Z-12:00Z",
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
						"bandwidth":                     "200",
						"vpc_auth_mode":                 "Open",
						"config.%":                      "3",
						"config.appendonly":             "yes",
						"config.lazyfree-lazy-eviction": "yes",
						"config.maxmemory-policy":       "volatile-lru",
						"tags.%":                        "2",
						"tags.Created":                  "TF",
						"tags.For":                      "acceptance test",
						//"maintain_start_time":           "04:00Z",
						//"maintain_end_time":             "06:00Z",
						//"backup_period.#":               "1",
						//"backup_time":                   "11:00Z-12:00Z",
						"private_connection_port":   "4011",
						"private_connection_prefix": CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudKVStoreRedisInstance_7_0_with_proxy_class(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudKVStoreMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreRedisInstance7_0_with_proxy-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKVStoreRedisInstanceVpcBasicDependence0)
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
					"instance_class":        "redis.shard.with.proxy.small.ce",
					"db_instance_name":      name,
					"instance_type":         "Redis",
					"engine_version":        "7.0",
					"read_only_count":       "1",
					"slave_read_only_count": "1",
					"instance_charge_type":  "PostPaid",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones.0.id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":        "redis.shard.with.proxy.small.ce",
						"db_instance_name":      name,
						"instance_type":         "Redis",
						"engine_version":        "7.0",
						"read_only_count":       "1",
						"slave_read_only_count": "1",
						"instance_charge_type":  "PostPaid",
						"tags.%":                "2",
						"tags.Created":          "TF",
						"tags.For":              "acceptance test",
						"resource_group_id":     CHECKSET,
						"zone_id":               CHECKSET,
						"vswitch_id":            CHECKSET,
						"secondary_zone_id":     CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "business_info", "coupon_no", "effective_time", "force_upgrade", "global_instance_id", "order_type", "password", "period", "enable_public", "security_ip_group_attribute", "enable_backup_log", "port"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"read_only_count": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"read_only_count": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"slave_read_only_count": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"slave_read_only_count": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"read_only_count":       "3",
					"slave_read_only_count": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"read_only_count":       "3",
						"slave_read_only_count": "3",
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
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones.1.id}",
					"vswitch_id":        "${data.alicloud_vswitches.update.ids.0}",
					"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones.0.id}",
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
			// There is an OpenAPI bug in eu-central-1
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"maintain_start_time": "02:00Z",
			//		"maintain_end_time":   "03:00Z",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"maintain_start_time": "02:00Z",
			//			"maintain_end_time":   "03:00Z",
			//		}),
			//	),
			//},
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
					"instance_charge_type": "PrePaid",
					"period":               "1",
					"auto_renew":           "true",
					"auto_renew_period":    "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "PrePaid",
						"auto_renew":           "true",
						"auto_renew_period":    "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type": "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "PostPaid",
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
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones.0.id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones.1.id}",
					// There is an OpenAPI bug in eu-central-1
					//"maintain_start_time": "04:00Z",
					//"maintain_end_time":   "06:00Z",
					// There is an OpenAPI bug
					//"backup_period":             []string{"Wednesday"},
					//"backup_time":               "11:00Z-12:00Z",
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
						"secondary_zone_id":             CHECKSET,
						//"maintain_start_time":           "04:00Z",
						//"maintain_end_time":             "06:00Z",
						//"backup_period.#":               "1",
						//"backup_time":                   "11:00Z-12:00Z",
						"private_connection_port":   "4011",
						"private_connection_prefix": CHECKSET,
					}),
				),
			},
		},
	})
}

// engine_version 4.0 has been offline from July 31, 2025
func SkipTestAccAliCloudKVStoreRedisInstance_prepaid(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudKVStoreMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreRedisInstancePrePaid%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKVStoreRedisInstancePrePaidBasicDependence0)
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
					"auto_renew":        "true",
					"auto_renew_period": "1",
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
						"auto_renew":        "true",
						"auto_renew_period": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "business_info", "coupon_no", "effective_time", "force_upgrade", "global_instance_id", "order_type", "password", "period", "enable_public", "security_ip_group_attribute", "enable_backup_log", "port"},
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
					"payment_type": "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PostPaid",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":      "PrePaid",
					"period":            "1",
					"auto_renew":        "true",
					"auto_renew_period": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":      "PrePaid",
						"auto_renew":        "true",
						"auto_renew_period": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PostPaid",
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
			// There is an OpenAPI bug in eu-central-1
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"maintain_start_time": "02:00Z",
			//		"maintain_end_time":   "03:00Z",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"maintain_start_time": "02:00Z",
			//			"maintain_end_time":   "03:00Z",
			//		}),
			//	),
			//},
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
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones[0].id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"secondary_zone_id": REMOVEKEY,
					// There is an OpenAPI bug in eu-central-1
					//"maintain_start_time": "04:00Z",
					//"maintain_end_time":   "06:00Z",
					// There is an OpenAPI bug
					//"backup_period":             []string{"Wednesday"},
					//"backup_time":               "11:00Z-12:00Z",
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
						//"maintain_start_time":           "04:00Z",
						//"maintain_end_time":             "06:00Z",
						//"backup_period.#":               "1",
						//"backup_time":                   "11:00Z-12:00Z",
						"private_connection_port":   "4011",
						"private_connection_prefix": CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudKVStoreRedisInstance_5_0_memory_classic_standard(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	// en-central-1 has no enough quota for this class
	checkoutSupportedRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudKVStoreMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreRedisInstanceVpcMultiTest%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKVStoreRedisInstanceVpcBasicDependence0)
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
					"instance_class":       "redis.amber.master.small.multithread",
					"db_instance_name":     name,
					"instance_type":        "Redis",
					"engine_version":       "5.0",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"zone_id":              "${data.alicloud_kvstore_zones.default.zones.0.id}",
					"vswitch_id":           "${data.alicloud_vswitches.default.ids.0}",
					"instance_charge_type": "PrePaid",
					"period":               "1",
					"is_auto_upgrade_open": "1",
					//"shard_count":       "2",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":       "redis.amber.master.small.multithread",
						"db_instance_name":     name,
						"instance_type":        "Redis",
						"engine_version":       "5.0",
						"resource_group_id":    CHECKSET,
						"zone_id":              CHECKSET,
						"vswitch_id":           CHECKSET,
						"instance_charge_type": "PrePaid",
						"period":               "1",
						"is_auto_upgrade_open": "1",
						//"shard_count":       "2",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "business_info", "coupon_no", "effective_time", "force_upgrade", "global_instance_id", "order_type", "password", "period", "enable_public", "security_ip_group_attribute", "enable_backup_log", "port"},
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
					"instance_charge_type": "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "PostPaid",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type": "PrePaid",
					"period":               "1",
					"auto_renew":           "true",
					"auto_renew_period":    "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "PrePaid",
						"auto_renew":           "true",
						"auto_renew_period":    "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type": "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "PostPaid",
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
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones.1.id}",
					"vswitch_id":        "${data.alicloud_vswitches.update.ids.0}",
					"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones.0.id}",
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
			// There is an OpenAPI bug in eu-central-1
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"maintain_start_time": "02:00Z",
			//		"maintain_end_time":   "03:00Z",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"maintain_start_time": "02:00Z",
			//			"maintain_end_time":   "03:00Z",
			//		}),
			//	),
			//},
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
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones.0.id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"secondary_zone_id": REMOVEKEY,
					// There is an OpenAPI bug in eu-central-1
					//"maintain_start_time":       "04:00Z",
					//"maintain_end_time":         "06:00Z",
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
						//"maintain_start_time":           "04:00Z",
						//"maintain_end_time":             "06:00Z",
						"backup_period.#":           "1",
						"backup_time":               "11:00Z-12:00Z",
						"private_connection_port":   "4011",
						"private_connection_prefix": CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudKVStoreRedisInstance_5_0_memory_classic_cluster(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	// en-central-1 has no enough quota for this class
	checkoutSupportedRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudKVStoreMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreRedisInstanceVpcMultiTest%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKVStoreRedisInstanceVpcBasicDependence0)
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
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones.0.id}",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "business_info", "coupon_no", "effective_time", "force_upgrade", "global_instance_id", "order_type", "password", "period", "enable_public", "security_ip_group_attribute", "enable_backup_log", "port"},
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
			//		"zone_id":           "${data.alicloud_kvstore_zones.default.zones.1.id}",
			//		"vswitch_id":        "${data.alicloud_vswitches.update.ids.0}",
			//		"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones.0.id}",
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
			// There is an OpenAPI bug in eu-central-1
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"maintain_start_time": "02:00Z",
			//		"maintain_end_time":   "03:00Z",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"maintain_start_time": "02:00Z",
			//			"maintain_end_time":   "03:00Z",
			//		}),
			//	),
			//},
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
					"ssl_enable":                "Disable",
					"private_connection_prefix": fmt.Sprintf("privateprefix%d", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enable":                "Disable",
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
					"is_auto_upgrade_open": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"is_auto_upgrade_open": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"is_auto_upgrade_open": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"is_auto_upgrade_open": "0",
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
					"zone_id":    "${data.alicloud_kvstore_zones.default.zones.0.id}",
					"vswitch_id": "${data.alicloud_vswitches.default.ids.0}",
					// There is an OpenAPI bug in eu-central-1
					//"maintain_start_time":       "04:00Z",
					//"maintain_end_time":         "06:00Z",
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
						//"maintain_start_time":           "04:00Z",
						//"maintain_end_time":             "06:00Z",
						"backup_period.#":           "1",
						"backup_time":               "11:00Z-12:00Z",
						"private_connection_port":   "4011",
						"private_connection_prefix": CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudKVStoreMemcacheInstance_vpctest(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudKVStoreMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreMemcacheInstanceVpcTest%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKVStoreMemcacheInstanceVpcBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
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
					"zone_id":           "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"secondary_zone_id": "${data.alicloud_vswitches.slave.vswitches.0.zone_id}",
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
				ImportStateVerifyIgnore: []string{"dry_run", "business_info", "coupon_no", "effective_time", "force_upgrade", "global_instance_id", "order_type", "password", "period", "enable_public", "security_ip_group_attribute", "enable_backup_log", "port"},
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
					"zone_id":           "${data.alicloud_vswitches.slave.vswitches.0.zone_id}",
					"vswitch_id":        "${data.alicloud_vswitches.slave.ids.0}",
					"secondary_zone_id": "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
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
					"instance_release_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_release_protection": "false",
					}),
				),
			},
		},
	})
}

// TestAccAliCloudKVStoreRedisInstance_privateConnectionCreate is a regression for
// a create-time config that pins the connection string fields to the instance's
// default private endpoint values (e.g. private_connection_port = "6379"). The
// Create -> Update chain marks every explicitly configured attribute as changed
// on a brand new resource, and ModifyDBInstanceConnectionString used to reject
// same-value modifications with InvalidConnectionStringOrPort.Duplicate, tainting
// the resource and forcing an unrecoverable replace on the next apply.
func TestAccAliCloudKVStoreRedisInstance_privateConnectionCreate(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudKVStoreMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreRedisPrivConnCreate%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKVStoreRedisInstanceVpcBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Create step carries both connection string fields, with port pinned to
				// the instance default (6379). Historical bug: create failed here.
				// Uses redis.master.small.default (classic master/replica, engine 5.0)
				// since it is broadly sellable; the sharded proxy class is not.
				// password is set here (rather than in the attributesCoverage test's
				// step 1) because this test does not exercise the plain-password ->
				// KMS migration path, so it steers clear of the kmsDiffSuppressFunc
				// misread that would otherwise fail on that migration step.
				Config: testAccConfig(map[string]interface{}{
					"instance_class":            "redis.master.small.default",
					"db_instance_name":          name,
					"instance_type":             "Redis",
					"engine_version":            "5.0",
					"payment_type":              "PostPaid",
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"zone_id":                   "${data.alicloud_kvstore_zones.default.zones.0.id}",
					"vswitch_id":                "${data.alicloud_vswitches.default.ids.0}",
					"password":                  "YourPassword_123",
					"private_connection_prefix": fmt.Sprintf("privateprefixcreate%d", rand),
					"private_connection_port":   "6379",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":            "redis.master.small.default",
						"db_instance_name":          name,
						"instance_type":             "Redis",
						"engine_version":            "5.0",
						"payment_type":              "PostPaid",
						"resource_group_id":         CHECKSET,
						"zone_id":                   CHECKSET,
						"vswitch_id":                CHECKSET,
						"private_connection_port":   "6379",
						"private_connection_prefix": CHECKSET,
					}),
				),
			},
			{
				// In-place prefix change: verifies the diff filter still forwards a real change.
				Config: testAccConfig(map[string]interface{}{
					"private_connection_prefix": fmt.Sprintf("privateprefixmodify%d", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_connection_prefix": CHECKSET,
						"private_connection_port":   "6379",
					}),
				),
			},
			{
				// Unrelated update while connection fields stay put: the diff filter must
				// collapse to a no-op instead of re-sending the current prefix/port.
				Config: testAccConfig(map[string]interface{}{
					"db_instance_name": name + "_upd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_name":        name + "_upd",
						"private_connection_port": "6379",
					}),
				),
			},
		},
	})
}

// TestAccAliCloudKVStoreRedisInstance_attributesCoverage is a create-only
// coverage-focused test: it wires up in step 1 every Optional attribute the
// resource is capable of exercising on redis.master.small.default (community
// Redis, engine 5.0, PostPaid), including create-only fields such as capacity
// / private_ip / auto_use_coupon / global_instance / business_info /
// coupon_no / effective_time / force_upgrade / order_type. It has no
// ImportState step on purpose: several of the attributes we need to cover
// are ForceNew and would necessarily diff on import. Subsequent steps
// exercise the real Update paths (resource_group_id ModifyResourceGroup,
// security_group_id ModifyInstanceSpec, kms_encrypted_password / context via
// ModifyInstanceAttribute's KMS Decrypt path), giving modify coverage for
// attributes not otherwise exercised.
//
// The KMS credential (kms_encrypted_password + kms_encryption_context) is
// introduced at create time rather than migrated from a plain password in a
// later step. The migration path (REMOVEKEY password + set KMS in the same
// step) hits a kmsDiffSuppressFunc bug: the suppress reads the still-set
// password out of state and suppresses the KMS diff, so the migration step
// arrives at ModifyInstanceAttribute with neither password nor KMS set and
// the provider bails at :1144 ("One of the 'password' and
// 'kms_encrypted_password' should be set"). Sidestepping the migration
// still gives full must-set + modify coverage for the KMS pair; the
// plain-password path stays covered by
// TestAccAliCloudKVStoreRedisInstance_privateConnectionCreate.
//
// Attributes deliberately NOT covered here (require external prerequisites the
// provider cannot synthesize, or the API path is systemically unreachable in
// the test environment):
//   - backup_id / restore_time / srcdb_instance_id: need an existing backup.
//   - dedicated_host_group_id: needs a dedicated host cluster.
//   - global_instance_id: needs an existing global-instance to join.
//   - tde_status / encryption_name / encryption_key / role_arn: TDE update
//     path (see resource_alicloud_kvstore_instance.go :1514) fires whenever
//     any of the four HasChange, calls ModifyInstanceTDE, and community
//     Redis (redis.master.small.default) returns IsSupportTDE=false at
//     Read (:764), so the API rejects the call on this class.
//   - port: only reachable together with enable_public +
//     connection_string_prefix via AllocateInstancePublicConnection. That
//     API returned InternalFailure for two independent instances (both at
//     create time and as a follow-up step) in the test environment, so
//     the public-endpoint path is skipped end-to-end here. port stays in
//     ImportStateVerifyIgnore for import safety of any future test that
//     does set it.
//   - engine_version modify (create-time set is still covered): the
//     classic redis.master.*.default series only has 5.0 available for
//     purchase in the test region (4.0 retired, 6.0 / 7.0 not offered on
//     this series), so ModifyInstanceMajorVersion has no reachable
//     in-place upgrade target and returns InvalidDBInstanceClass.NotFound.
func TestAccAliCloudKVStoreRedisInstance_attributesCoverage(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudKVStoreMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreRedisAttrCoverage%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKVStoreRedisInstanceAttributesCoverageDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Full create-time set. See function-level comment for why the
				// ForceNew and deprecated fields have to land in the create step.
				// The public-endpoint knobs (enable_public /
				// connection_string_prefix / port) are intentionally NOT set:
				// see the function-level comment for why.
				Config: testAccConfig(map[string]interface{}{
					"instance_type":    "Redis",
					"instance_class":   "redis.master.small.default",
					"engine_version":   "5.0",
					"payment_type":     "PostPaid",
					"db_instance_name": name,
					"zone_id":          "${data.alicloud_kvstore_zones.default.zones.0.id}",
					"vswitch_id":       "${alicloud_vswitch.default.id}",
					// kms_encrypted_password + kms_encryption_context land at
					// create time (rather than migrating from plain password in
					// a later step) because kmsDiffSuppressFunc misreads the
					// still-populated password state and suppresses the KMS
					// diff, leaving the migration step with neither value set
					// (kvstore rejects that at :1144). Starting with KMS avoids
					// the suppress condition; the plain-password path is still
					// covered by TestAccAliCloudKVStoreRedisInstance_privateConnectionCreate.
					"kms_encrypted_password": "${alicloud_kms_ciphertext.default.ciphertext_blob}",
					"kms_encryption_context": map[string]string{"kv": "1"},
					"capacity":               "1024",
					"private_ip":             "172.16.0.100",
					"business_info":          "terraform_test",
					"coupon_no":              "youhuiquan_promotion_option_id_for_blank",
					// auto_use_coupon and global_instance have DiffSuppressFunc
					// !d.IsNewResource() and no Read setter, so their post-create
					// state is empty; keep them in config (for must-set coverage
					// and for the real API round-trip) but skip Check assertions
					// below to avoid TestCheckResourceAttr "not found" failures.
					"auto_use_coupon":             "false",
					"global_instance":             "false",
					"effective_time":              "Immediately",
					"force_upgrade":               "true",
					"order_type":                  "UPGRADE",
					"security_group_id":           "${alicloud_security_group.default.id}",
					"security_ip_group_attribute": "tftest",
					"security_ip_group_name":      "coveragegrp1",
					"security_ips":                []string{"10.23.12.24"},
					"enable_backup_log":           "0",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"tags": map[string]string{
						"coverage": "yes",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":            "redis.master.small.default",
						"engine_version":            "5.0",
						"payment_type":              "PostPaid",
						"db_instance_name":          name,
						"private_ip":                "172.16.0.100",
						"business_info":             "terraform_test",
						"coupon_no":                 "youhuiquan_promotion_option_id_for_blank",
						"kms_encrypted_password":    CHECKSET,
						"kms_encryption_context.%":  "1",
						"kms_encryption_context.kv": "1",
						// auto_use_coupon and global_instance intentionally omitted here
						// (see config-side comment above; DiffSuppressFunc + no Read setter
						// means the attribute never lands in state after Read).
						"effective_time":              "Immediately",
						"force_upgrade":               "true",
						"order_type":                  "UPGRADE",
						"security_group_id":           CHECKSET,
						"security_ip_group_attribute": "tftest",
						"security_ip_group_name":      "coveragegrp1",
						"enable_backup_log":           "0",
						"resource_group_id":           CHECKSET,
						"tags.%":                      "1",
						"tags.coverage":               "yes",
					}),
				),
			},
			{
				// resource_group_id modify -> ModifyResourceGroup.
				// security_ip_group_name modify -> ModifySecurityIps (piggybacked
				// so the same round-trip also covers this attribute's modify).
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"security_ip_group_name": "coveragegrp2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id":      CHECKSET,
						"security_ip_group_name": "coveragegrp2",
					}),
				),
			},
			{
				// security_group_id modify -> ModifyInstanceSpec / security-group swap.
				Config: testAccConfig(map[string]interface{}{
					"security_group_id": "${alicloud_security_group.update.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id": CHECKSET,
					}),
				),
			},
			{
				// Swap to the second KMS ciphertext with a different
				// encryption_context: gives modify coverage for
				// kms_encrypted_password AND kms_encryption_context and
				// exercises ModifyInstanceAttribute's KMS Decrypt path
				// (:1138). password stays unset throughout this test, so the
				// kmsDiffSuppressFunc suppress condition (see step 1 comment)
				// is not triggered.
				Config: testAccConfig(map[string]interface{}{
					"kms_encrypted_password": "${alicloud_kms_ciphertext.update.ciphertext_blob}",
					"kms_encryption_context": map[string]string{
						"kv": "2",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kms_encrypted_password":    CHECKSET,
						"kms_encryption_context.%":  "1",
						"kms_encryption_context.kv": "2",
					}),
				),
			},
		},
	})
}

var AliCloudKVStoreMap0 = map[string]string{
	"connection_domain": CHECKSET,
	"bandwidth":         CHECKSET,
	"ssl_enable":        "Disable",
}

func AliCloudKVStoreRedisInstanceVpcBasicDependence0(name string) string {
	return fmt.Sprintf(`
	data "alicloud_resource_manager_resource_groups" "default" {
  		status = "OK"
	}

	data "alicloud_kvstore_zones" "default" {
  		instance_charge_type = "PostPaid"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_kvstore_zones.default.zones.0.id
	}

	data "alicloud_vswitches" "update" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_kvstore_zones.default.zones.1.id
	}

	resource "alicloud_security_group" "default" {
  		inner_access_policy = "Accept"
  		name                = "tf-example"
  		vpc_id              = data.alicloud_vpcs.default.ids.0
	}
	`)
}

func AliCloudKVStoreRedisInstancePrePaidBasicDependence0(name string) string {
	return fmt.Sprintf(`
	data "alicloud_resource_manager_resource_groups" "default" {
  		status = "OK"
	}

	data "alicloud_kvstore_zones" "default" {
  		instance_charge_type = "PrePaid"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_kvstore_zones.default.zones.0.id
	}

	data "alicloud_vswitches" "update" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_kvstore_zones.default.zones.1.id
	}
	`)
}

func AliCloudKVStoreMemcacheInstanceVpcBasicDependence0(name string) string {
	return fmt.Sprintf(`
	data "alicloud_resource_manager_resource_groups" "default" {
  		status = "OK"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = "cn-hangzhou-h"
	}

	data "alicloud_vswitches" "slave" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = "cn-hangzhou-i"
	}
	`)
}

// AliCloudKVStoreRedisInstanceAttributesCoverageDependence provisions the
// full support cast for TestAccAliCloudKVStoreRedisInstance_attributesCoverage:
// a self-owned VPC + vswitch (predictable CIDR so private_ip can be pinned),
// two security groups (for a modify step), and a KMS key with two ciphertexts
// (for the kms_encrypted_password / kms_encryption_context modify pair).
func AliCloudKVStoreRedisInstanceAttributesCoverageDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_kvstore_zones" "default" {
		instance_charge_type = "PostPaid"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
		status = "OK"
	}

	resource "alicloud_vpc" "default" {
		vpc_name   = var.name
		cidr_block = "172.16.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
		vpc_id       = alicloud_vpc.default.id
		cidr_block   = "172.16.0.0/24"
		zone_id      = data.alicloud_kvstore_zones.default.zones.0.id
		vswitch_name = var.name
	}

	resource "alicloud_security_group" "default" {
		name   = "${var.name}-1"
		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_security_group" "update" {
		name   = "${var.name}-2"
		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_kms_key" "default" {
		description            = var.name
		pending_window_in_days = 7
		status                 = "Enabled"
	}

	resource "alicloud_kms_ciphertext" "default" {
		key_id    = alicloud_kms_key.default.id
		plaintext = "YourPassword_kms1"
		encryption_context = {
			kv = "1"
		}
	}

	resource "alicloud_kms_ciphertext" "update" {
		key_id    = alicloud_kms_key.default.id
		plaintext = "YourPassword_kms2"
		encryption_context = {
			kv = "2"
		}
	}
	`, name)
}
