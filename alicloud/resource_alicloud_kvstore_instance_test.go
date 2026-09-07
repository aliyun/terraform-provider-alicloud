package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers("alicloud_kvstore_instance", &resource.Sweeper{
		Name: "alicloud_kvstore_instance",
		F:    testSweepKVStoreInstances,
	})
}

func TestAccAliCloudKVStoreRedisInstance_coverage(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudKVStoreMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreRedisInstanceCoverage%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKVStoreRedisInstanceVpcBasicDependence0)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_use_coupon":             "false",
					"backup_id":                   "backup-id",
					"business_info":               "business-info",
					"capacity":                    "1",
					"coupon_no":                   "coupon-no",
					"dedicated_host_group_id":     "dhg-id",
					"effective_time":              "Immediately",
					"enable_backup_log":           "0",
					"encryption_key":              "key",
					"encryption_name":             "AES-CTR-256",
					"engine_version":              "6.0",
					"force_upgrade":               "true",
					"global_instance":             "false",
					"global_instance_id":          "global-instance-id",
					"kms_encrypted_password":      "kms-password",
					"kms_encryption_context":      map[string]string{"name": "value"},
					"order_type":                  "UPGRADE",
					"password":                    "Test12345",
					"port":                        "6379",
					"private_ip":                  "192.168.0.10",
					"restore_time":                "2026-01-01T00:00:00Z",
					"resource_group_id":           "rg-id",
					"role_arn":                    "acs:ram::1234567890123456:role/example",
					"security_group_id":           "sg-id",
					"security_ip_group_attribute": "group-attribute",
					"srcdb_instance_id":           "r-source",
				}),
				ExpectError: regexp.MustCompile(".+"),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encryption_key":         "key-update",
					"encryption_name":        "AES-CTR-256-update",
					"engine_version":         "7.0",
					"kms_encrypted_password": "kms-password-update",
					"kms_encryption_context": map[string]string{"name": "value-update"},
					"port":                   "6380",
					"resource_group_id":      "rg-id-update",
					"role_arn":               "acs:ram::1234567890123456:role/example-update",
					"security_group_id":      "sg-id-update",
				}),
				ExpectError: regexp.MustCompile(".+"),
			},
		},
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
				ImportStateVerifyIgnore: []string{"dry_run", "business_info", "coupon_no", "effective_time", "force_upgrade", "global_instance_id", "order_type", "password", "period", "enable_public", "security_ip_group_attribute", "enable_backup_log"},
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
					"timeouts": []map[string]interface{}{
						{
							"update": "1h",
						},
					},
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
				ImportStateVerifyIgnore: []string{"dry_run", "business_info", "coupon_no", "effective_time", "force_upgrade", "global_instance_id", "order_type", "password", "period", "enable_public", "security_ip_group_attribute", "enable_backup_log"},
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
				ImportStateVerifyIgnore: []string{"dry_run", "business_info", "coupon_no", "effective_time", "force_upgrade", "global_instance_id", "order_type", "password", "period", "enable_public", "security_ip_group_attribute", "enable_backup_log"},
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
						"secondary_zone_id":    REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secondary_zone_id":           REMOVEKEY,
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
				ImportStateVerifyIgnore: []string{"dry_run", "business_info", "coupon_no", "effective_time", "force_upgrade", "global_instance_id", "order_type", "password", "period", "enable_public", "security_ip_group_attribute", "enable_backup_log"},
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
				ImportStateVerifyIgnore: []string{"dry_run", "business_info", "coupon_no", "effective_time", "force_upgrade", "global_instance_id", "order_type", "password", "period", "enable_public", "security_ip_group_attribute", "enable_backup_log"},
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

// replica_count / slave_replica_count are mutually exclusive with read_only_count /
// slave_read_only_count (the API rejects "both replicas and read-only nodes at the same
// time"), so replicas get their own case here rather than being mixed into the read/write
// splitting case above.
func TestAccAliCloudKVStoreRedisInstance_7_0_with_proxy_class_replica(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudKVStoreMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreRedisInstance7_0_proxy_replica-%d", rand)
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
					"instance_class":       "redis.shard.small.ce",
					"db_instance_name":     name,
					"instance_type":        "Redis",
					"engine_version":       "7.0",
					"shard_count":          "2",
					"replica_count":        "1",
					"slave_replica_count":  "1",
					"instance_charge_type": "PostPaid",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"zone_id":              "${data.alicloud_kvstore_zones.default.zones.0.id}",
					"vswitch_id":           "${data.alicloud_vswitches.default.ids.0}",
					"secondary_zone_id":    "${data.alicloud_kvstore_zones.default.zones.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":       "redis.shard.small.ce",
						"db_instance_name":     name,
						"instance_type":        "Redis",
						"engine_version":       "7.0",
						"shard_count":          "2",
						"replica_count":        "1",
						"slave_replica_count":  "1",
						"instance_charge_type": "PostPaid",
						"resource_group_id":    CHECKSET,
						"zone_id":              CHECKSET,
						"vswitch_id":           CHECKSET,
						"secondary_zone_id":    CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "business_info", "coupon_no", "effective_time", "force_upgrade", "global_instance_id", "order_type", "password", "period", "enable_public", "security_ip_group_attribute", "enable_backup_log"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"replica_count":       "2",
					"slave_replica_count": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"replica_count":       "2",
						"slave_replica_count": "2",
					}),
				),
			},
			// Dual-zone -> single-zone conversion by clearing secondary_zone_id. The
			// provider zeros SlaveReplicaCount and absorbs the former secondary replicas
			// into ReplicaCount so the sum stays invariant (2 + 2 -> 4 primary, 0 slave),
			// then omits SecondaryZoneId so MigrateToOtherZone converts the instance to
			// single-zone. slave_replica_count/replica_count are dropped from the config in
			// this step: a single-zone instance cannot hold secondary-zone replicas, so the
			// counts are computed from the absorbed topology (both are Optional+Computed)
			// rather than declared, which also avoids a concurrent ModifyInstanceSpec.
			{
				Config: testAccConfig(map[string]interface{}{
					"replica_count":       REMOVEKEY,
					"slave_replica_count": REMOVEKEY,
					"secondary_zone_id":   REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"replica_count":       "4",
						"slave_replica_count": "0",
						"secondary_zone_id":   REMOVEKEY,
					}),
				),
			},
			// After the conversion the state matches the config: secondary_zone_id is
			// empty and replica_count/slave_replica_count are computed, so the plan must
			// be empty (no spurious diff from the former no-op REMOVEKEY behavior). The
			// keys were already dropped in the previous step, so re-plan the cumulative
			// config as-is (an empty change map) rather than re-applying REMOVEKEY.
			{
				Config:             testAccConfig(map[string]interface{}{}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccAliCloudKVStoreRedisInstance_single_copy(t *testing.T) {
	testAccKvstoreNodeType(t, "single", "redis.shard.small.2.ce", "OnECS", "")
}

func TestAccAliCloudKVStoreRedisInstance_node_type_double(t *testing.T) {
	testAccKvstoreNodeType(t, "double", "redis.amber.master.small.multithread", "", "double")
}

func testAccKvstoreNodeType(t *testing.T, nodeType, instanceClass, productType, inventoryNodeType string) {
	t.Helper()
	var v r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, resourceAttrInit(resourceId, nil))
	name := fmt.Sprintf("tf-testAccKvstoreNodeType%d", acctest.RandIntRange(1000000, 9999999))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		productFilter, nodeFilter := "", ""
		if productType != "" {
			productFilter = fmt.Sprintf("product_type = %q", productType)
		}
		// OnECS inventory reports numeric values instead of classic node-type values.
		if inventoryNodeType != "" {
			nodeFilter = fmt.Sprintf("node_type = %q", inventoryNodeType)
		}
		return fmt.Sprintf(`
data "alicloud_kvstore_zones" "default" {
  instance_charge_type = "PostPaid"
  %s
}

data "alicloud_kvstore_instance_classes" "default" {
  count                = length(data.alicloud_kvstore_zones.default.ids)
  zone_id              = data.alicloud_kvstore_zones.default.ids[count.index]
  instance_charge_type = "PostPaid"
  engine               = "Redis"
  engine_version       = "5.0"
  %s
  %s
}

locals {
  class_supported_zones = [for i, classes in data.alicloud_kvstore_instance_classes.default : data.alicloud_kvstore_zones.default.ids[i] if contains(classes.instance_classes, %q)]
  # The catalog also includes zones that no longer support instance creation.
  supported_zones = [for zone in ["cn-hangzhou-i", "cn-hangzhou-j", "cn-hangzhou-k"] : zone if contains(local.class_supported_zones, zone)]
}

resource "alicloud_vpc" "default" {
  vpc_name   = %q
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = %q
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.0.0/24"
  zone_id      = local.supported_zones[0]
}
`, productFilter, productFilter, nodeFilter, instanceClass, name, name)
	})
	var instanceId string
	check := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceId, "instance_class", instanceClass),
		resource.TestCheckResourceAttr(resourceId, "node_type", nodeType),
		func(s *terraform.State) error {
			rs, ok := s.RootModule().Resources[resourceId]
			if !ok || rs.Primary == nil || rs.Primary.ID == "" {
				return fmt.Errorf("instance is missing from state")
			}
			id := rs.Primary.ID
			if instanceId != "" && id != instanceId {
				return fmt.Errorf("instance was unexpectedly replaced")
			}
			service := R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
			instance, err := service.DescribeKvstoreInstance(id)
			if err != nil {
				return err
			}
			if instance["NodeType"] != nodeType {
				return fmt.Errorf("expected remote NodeType %q, got %v", nodeType, instance["NodeType"])
			}
			if instance["InstanceClass"] != instanceClass {
				return fmt.Errorf("expected remote InstanceClass %q, got %v", instanceClass, instance["InstanceClass"])
			}
			instanceId = id
			return nil
		},
	)
	importCheck := func(states []*terraform.InstanceState) error {
		if len(states) != 1 || states[0].ID != instanceId {
			return fmt.Errorf("expected the existing instance to be imported")
		}
		if states[0].Attributes["node_type"] != nodeType {
			return fmt.Errorf("expected imported node_type %q, got %q", nodeType, states[0].Attributes["node_type"])
		}
		for _, configureNodeType := range []bool{false, true} {
			config := map[string]interface{}{
				"instance_class":   instanceClass,
				"instance_type":    states[0].Attributes["instance_type"],
				"engine_version":   "5.0",
				"payment_type":     "PostPaid",
				"db_instance_name": states[0].Attributes["db_instance_name"],
				"vswitch_id":       states[0].Attributes["vswitch_id"],
			}
			if configureNodeType {
				config["node_type"] = nodeType
			}
			diff, err := resourceAliCloudKvstoreInstance().Diff(states[0], terraform.NewResourceConfigRaw(config), testAccProvider.Meta())
			if err != nil {
				return err
			}
			if diff != nil && diff.Attributes["node_type"] != nil {
				return fmt.Errorf("imported node_type changed with configuration %v: %v", config, diff.Attributes["node_type"])
			}
			if diff != nil && diff.RequiresNew() {
				return fmt.Errorf("imported instance unexpectedly requires replacement")
			}
		}
		return nil
	}
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
					"instance_class":   instanceClass,
					"node_type":        nodeType,
					"engine_version":   "5.0",
					"payment_type":     "PostPaid",
					"db_instance_name": name,
					"vswitch_id":       "${alicloud_vswitch.default.id}",
				}),
				Check: check,
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateCheck:        importCheck,
				ImportStateVerifyIgnore: []string{"dry_run", "business_info", "coupon_no", "effective_time", "force_upgrade", "global_instance_id", "order_type", "password", "period", "enable_public", "security_ip_group_attribute", "enable_backup_log"},
			},
			{
				Config:             testAccConfig(map[string]interface{}{}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_type":        REMOVEKEY,
					"db_instance_name": name + "-updated",
				}),
				Check: resource.ComposeTestCheckFunc(check,
					resource.TestCheckResourceAttr(resourceId, "db_instance_name", name+"-updated")),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateCheck:        importCheck,
				ImportStateVerifyIgnore: []string{"dry_run", "business_info", "coupon_no", "effective_time", "force_upgrade", "global_instance_id", "order_type", "password", "period", "enable_public", "security_ip_group_attribute", "enable_backup_log"},
			},
			{
				Config:             testAccConfig(map[string]interface{}{}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccKvstoreClassicDependence("Redis", "5.0", "PrePaid", "redis.amber.master.small.multithread", "redis.amber.master.mid.multithread"))
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
					"node_type":            "double",
					"db_instance_name":     name,
					"instance_type":        "Redis",
					"engine_version":       "5.0",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"zone_id":              "${alicloud_vswitch.default.zone_id}",
					"vswitch_id":           "${alicloud_vswitch.default.id}",
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
						"node_type":            "double",
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
				ImportStateVerifyIgnore: []string{"dry_run", "business_info", "coupon_no", "effective_time", "force_upgrade", "global_instance_id", "order_type", "password", "period", "enable_public", "security_ip_group_attribute", "enable_backup_log"},
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
					"zone_id":           "${alicloud_vswitch.update.zone_id}",
					"vswitch_id":        "${alicloud_vswitch.update.id}",
					"secondary_zone_id": "${alicloud_vswitch.default.zone_id}",
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
					"zone_id":           "${alicloud_vswitch.default.zone_id}",
					"vswitch_id":        "${alicloud_vswitch.default.id}",
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccKvstoreClassicDependence("Redis", "5.0", "PostPaid", "redis.amber.logic.sharding.1g.2db.0rodb.6proxy.multithread", "redis.amber.logic.sharding.2g.2db.0rodb.6proxy.multithread"))
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
					"zone_id":           "${alicloud_vswitch.default.zone_id}",
					"vswitch_id":        "${alicloud_vswitch.default.id}",
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
				ImportStateVerifyIgnore: []string{"dry_run", "business_info", "coupon_no", "effective_time", "force_upgrade", "global_instance_id", "order_type", "password", "period", "enable_public", "security_ip_group_attribute", "enable_backup_log"},
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
			//		"zone_id":           "${alicloud_vswitch.update.zone_id}",
			//		"vswitch_id":        "${alicloud_vswitch.update.id}",
			//		"secondary_zone_id": "${alicloud_vswitch.default.zone_id}",
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
					"zone_id":    "${alicloud_vswitch.default.zone_id}",
					"vswitch_id": "${alicloud_vswitch.default.id}",
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
	// Memcache inventory leaves Version empty even though instances use engine 4.0.
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccKvstoreClassicDependence("Memcache", "", "PostPaid", "memcache.master.small.default", "memcache.master.mid.default"))
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
					"zone_id":           "${alicloud_vswitch.default.zone_id}",
					"vswitch_id":        "${alicloud_vswitch.default.id}",
					"secondary_zone_id": "${alicloud_vswitch.update.zone_id}",
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
				ImportStateVerifyIgnore: []string{"dry_run", "business_info", "coupon_no", "effective_time", "force_upgrade", "global_instance_id", "order_type", "password", "period", "enable_public", "security_ip_group_attribute", "enable_backup_log"},
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
					"zone_id":           "${alicloud_vswitch.update.zone_id}",
					"vswitch_id":        "${alicloud_vswitch.update.id}",
					"secondary_zone_id": "${alicloud_vswitch.default.zone_id}",
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

// TestAccAliCloudKVStoreRedisInstance_classic_cluster_instance_class isolates the local-disk
// (classic architecture) cluster instance_class change path. Changing the spec of a local-disk
// instance via ModifyInstanceSpec requires MajorVersion to be carried alongside InstanceClass;
// cloud-disk (.ce/.ee/tair.*) specs do not. This case creates a local-disk cluster, reimports,
// then modifies instance_class both up and back down to exercise that path in isolation.
func TestAccAliCloudKVStoreRedisInstance_classic_cluster_instance_class(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	checkoutSupportedRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudKVStoreMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreRedisClassicClusterSpec%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccKvstoreClassicDependence("Redis", "5.0", "PostPaid", "redis.amber.logic.sharding.1g.2db.0rodb.6proxy.multithread", "redis.amber.logic.sharding.2g.2db.0rodb.6proxy.multithread"))
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
					"instance_class":   "redis.amber.logic.sharding.1g.2db.0rodb.6proxy.multithread",
					"db_instance_name": name,
					"instance_type":    "Redis",
					"engine_version":   "5.0",
					"zone_id":          "${alicloud_vswitch.default.zone_id}",
					"vswitch_id":       "${alicloud_vswitch.default.id}",
					"shard_count":      "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":   "redis.amber.logic.sharding.1g.2db.0rodb.6proxy.multithread",
						"db_instance_name": name,
						"instance_type":    "Redis",
						"engine_version":   "5.0",
						"zone_id":          CHECKSET,
						"vswitch_id":       CHECKSET,
						"shard_count":      "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "business_info", "coupon_no", "effective_time", "force_upgrade", "global_instance_id", "order_type", "password", "period", "enable_public", "security_ip_group_attribute", "enable_backup_log"},
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
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class": "redis.amber.logic.sharding.1g.2db.0rodb.6proxy.multithread",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class": "redis.amber.logic.sharding.1g.2db.0rodb.6proxy.multithread",
					}),
				),
			},
		},
	})
}

func TestUnitKvstoreNodeTypeCompatibility(t *testing.T) {
	r := resourceAliCloudKvstoreInstance()
	for _, nodeType := range []string{"single", "double", "MASTER_SLAVE", "STAND_ALONE"} {
		t.Run(nodeType, func(t *testing.T) {
			config := terraform.NewResourceConfigRaw(map[string]interface{}{"node_type": nodeType})
			warnings, errors := r.Validate(config)
			assert.Empty(t, errors)
			assert.Empty(t, warnings, "supported node types must not produce a deprecation warning")

			data := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"node_type": nodeType})
			data.SetId("test-instance")
			state := data.State()
			for _, config := range []map[string]interface{}{{}, {"node_type": nodeType}} {
				diff, err := r.Diff(state, terraform.NewResourceConfigRaw(config), nil)
				if !assert.NoError(t, err) || diff == nil {
					continue
				}
				assert.NotContains(t, diff.Attributes, "node_type", "existing and imported node types must remain stable")
			}
			other := "single"
			if nodeType == other {
				other = "double"
			}
			diff, err := r.Diff(state, terraform.NewResourceConfigRaw(map[string]interface{}{"node_type": other}), nil)
			if assert.NoError(t, err) && assert.NotNil(t, diff) && assert.Contains(t, diff.Attributes, "node_type") {
				assert.True(t, diff.Attributes["node_type"].RequiresNew)
			}
		})
	}
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{})
	_, configured := d.GetOk("node_type")
	assert.False(t, configured, "omission must preserve service-side default selection")
}

// TestUnitKvstoreIsCloudDiskSpec verifies the architecture classification helper
// that decides whether MajorVersion must accompany InstanceClass on a
// ModifyInstanceSpec call. Cloud-disk specs (.ce/.ee suffixes and the tair.*
// prefix) do not require MajorVersion; everything else is treated as a
// local-disk (classic) spec and does require it.
func TestUnitKvstoreIsCloudDiskSpec(t *testing.T) {
	tests := []struct {
		name          string
		instanceClass string
		expected      bool
	}{
		{"cloud-disk ce suffix", "redis.master.small.ce", true},
		{"cloud-disk ee suffix", "redis.master.large.ee", true},
		{"tair prefix", "tair.rdb.1g", true},
		{"tair prefix mixed case", "tair.Redis.2g", true},
		{"local-disk classic amber sharding", "redis.amber.logic.sharding.1g.2db.0rodb.6proxy.multithread", false},
		{"local-disk classic master", "redis.master.small.default", false},
		{"empty string", "", false},
		{"tair substring not at prefix must not match", "redis.tair.1g", false},
		{"ce substring not at suffix must not match", "redis.ce.master", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isCloudDiskSpec(tt.instanceClass))
		})
	}
}

// TestUnitKvstoreResolveMajorVersionForSpecChange verifies MajorVersion
// resolution for the local-disk ModifyInstanceSpec path: the engine_version
// stored in state is preferred, and when it is empty the helper falls back to
// DescribeKvstoreInstance (whose EngineVersion reflects the instance's current
// major version).
func TestUnitKvstoreResolveMajorVersionForSpecChange(t *testing.T) {
	tests := []struct {
		name               string
		engineVersion      string
		instanceId         string
		describeFunc       func(id string) (map[string]interface{}, error)
		expectedVersion    string
		expectErr          bool
		expectFallbackCall bool
	}{
		{
			name:          "state engine_version present, no fallback query",
			engineVersion: "5.0",
			instanceId:    "r-xxx",
			describeFunc: func(id string) (map[string]interface{}, error) {
				t.Fatalf("fallback should not be called when engine_version is in state")
				return nil, nil
			},
			expectedVersion:    "5.0",
			expectFallbackCall: false,
		},
		{
			name:          "empty state engine_version falls back to DescribeKvstoreInstance",
			engineVersion: "",
			instanceId:    "r-abc",
			describeFunc: func(id string) (map[string]interface{}, error) {
				assert.Equal(t, "r-abc", id)
				return map[string]interface{}{"EngineVersion": "4.0"}, nil
			},
			expectedVersion:    "4.0",
			expectFallbackCall: true,
		},
		{
			name:          "fallback query error propagates",
			engineVersion: "",
			instanceId:    "r-err",
			describeFunc: func(id string) (map[string]interface{}, error) {
				return nil, fmt.Errorf("NotFound")
			},
			expectedVersion:    "",
			expectErr:          true,
			expectFallbackCall: true,
		},
		{
			name:          "fallback returns empty when EngineVersion missing from response",
			engineVersion: "",
			instanceId:    "r-missing",
			describeFunc: func(id string) (map[string]interface{}, error) {
				return map[string]interface{}{}, nil
			},
			expectedVersion:    "",
			expectFallbackCall: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resolveMajorVersionForSpecChange(tt.engineVersion, tt.describeFunc, tt.instanceId)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedVersion, result)
		})
	}
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
		product_type         = "OnECS"
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

func testAccKvstoreClassicDependence(engine, version, paymentType string, instanceClasses ...string) func(string) string {
	return func(name string) string {
		versionFilter := ""
		if version != "" {
			versionFilter = fmt.Sprintf("engine_version       = %q", version)
		}
		classes := make([]string, len(instanceClasses))
		for i, class := range instanceClasses {
			classes[i] = fmt.Sprintf("%q", class)
		}
		zones := `"cn-hangzhou-i", "cn-hangzhou-j", "cn-hangzhou-k"`
		if engine == "Memcache" {
			zones = `"cn-hangzhou-j", "cn-hangzhou-k"`
		}
		return fmt.Sprintf(`
data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_kvstore_zones" "default" {
  engine               = %q
  instance_charge_type = %q
}

data "alicloud_kvstore_instance_classes" "default" {
  count                = length(data.alicloud_kvstore_zones.default.ids)
  zone_id              = data.alicloud_kvstore_zones.default.ids[count.index]
  engine               = %q
  %s
  instance_charge_type = %q
}

locals {
  class_supported_zones = [for i, classes in data.alicloud_kvstore_instance_classes.default : data.alicloud_kvstore_zones.default.ids[i] if length([for class in [%s] : class if !contains(classes.instance_classes, class)]) == 0]
  # Older catalog entries can refer to zones closed for new instance creation.
  supported_zones = [for zone in [%s] : zone if contains(local.class_supported_zones, zone)]
}

resource "alicloud_vpc" "default" {
  vpc_name   = %q
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = %q
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.0.0/24"
  zone_id      = local.supported_zones[0]
}

resource "alicloud_vswitch" "update" {
  vswitch_name = %q
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.1.0/24"
  zone_id      = local.supported_zones[1]
}

resource "alicloud_security_group" "default" {
  name                = %q
  vpc_id              = alicloud_vpc.default.id
  inner_access_policy = "Accept"
}
`, engine, paymentType, engine, versionFilter, paymentType, strings.Join(classes, ", "), zones, name, name, name+"-update", name)
	}
}
