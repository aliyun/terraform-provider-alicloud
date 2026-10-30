package alicloud

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	sdkendpoints "github.com/aliyun/alibaba-cloud-sdk-go/sdk/endpoints"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKVStoreRedisInstanceVpcSecurityGroupDependence0)
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
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
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
						"ssl_enable":           "Enable",
						"ca_cert_download_url": CHECKSET,
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
			{
				Config: testAccConfig(map[string]interface{}{
					"global_security_group_list": "${alicloud_redis_global_security_ip_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_security_group_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_security_group_list": []string{"${alicloud_redis_global_security_ip_group.default.0.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_security_group_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_security_group_list": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_security_group_list.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_security_group_list": "${alicloud_redis_global_security_ip_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_security_group_list.#": "3",
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
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
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
					"resource_group_id":          "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"zone_id":                    "${data.alicloud_kvstore_zones.default.zones.0.id}",
					"vswitch_id":                 "${data.alicloud_vswitches.default.ids.0}",
					"secondary_zone_id":          "${data.alicloud_kvstore_zones.default.zones.1.id}",
					"global_security_group_list": "${alicloud_redis_global_security_ip_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":               "redis.shard.small.ce",
						"db_instance_name":             name,
						"instance_type":                "Redis",
						"engine_version":               "7.0",
						"shard_count":                  "2",
						"bandwidth":                    "200",
						"tags.%":                       "2",
						"tags.Created":                 "TF",
						"tags.For":                     "acceptance test",
						"resource_group_id":            CHECKSET,
						"zone_id":                      CHECKSET,
						"vswitch_id":                   CHECKSET,
						"secondary_zone_id":            CHECKSET,
						"global_security_group_list.#": "3",
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

func TestAccAliCloudKVStoreRedisInstance_secondaryZoneOnly(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudKVStoreMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreRedisInstanceSecondaryZoneOnly%d", rand)
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
					"instance_class":    "redis.shard.small.ce",
					"db_instance_name":  name,
					"instance_type":     "Redis",
					"engine_version":    "7.0",
					"shard_count":       "2",
					"payment_type":      "PostPaid",
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
						"payment_type":      "PostPaid",
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
				// Changing only secondary_zone_id (zone_id unchanged) must still send ZoneId to MigrateToOtherZone.
				Config: testAccConfig(map[string]interface{}{
					"secondary_zone_id": "${data.alicloud_kvstore_zones.default.zones.2.id}",
					"timeouts": []map[string]interface{}{
						{
							"update": "1h",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":           CHECKSET,
						"secondary_zone_id": CHECKSET,
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
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
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

// Destroying a PrePaid instance must fail with an error instead of silently
// removing the instance from the state file (behavior changed in v1.294.0).
func TestAccAliCloudKVStoreRedisInstance_prePaidDestroyError(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudKVStoreMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreRedisInstancePrePaidDestroyError%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKVStoreRedisInstanceVpcBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class":   "redis.shard.with.proxy.small.ce",
					"db_instance_name": name,
					"instance_type":    "Redis",
					"engine_version":   "6.0",
					"shard_count":      "2",
					"payment_type":     "PostPaid",
					"zone_id":          "${data.alicloud_kvstore_zones.default.zones.0.id}",
					"vswitch_id":       "${data.alicloud_vswitches.default.ids.0}",
					"timeouts": []map[string]interface{}{
						{
							"update": "1h",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PostPaid",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PrePaid",
					"period":       "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PrePaid",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PrePaid",
					"period":       "1",
				}),
				Destroy:     true,
				ExpectError: regexp.MustCompile(`'PrePaid' instance cannot be destroyed`),
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
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
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
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
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

func TestAccAliCloudKVStoreRedisInstance_node_type_stand_alone(t *testing.T) {
	testAccKvstoreNodeType(t, "STAND_ALONE", "single", "redis.shard.small.2.ce", "OnECS", "")
}

func TestAccAliCloudKVStoreRedisInstance_node_type_master_slave(t *testing.T) {
	testAccKvstoreNodeType(t, "MASTER_SLAVE", "double", "redis.shard.small.2.ce", "OnECS", "")
}

func TestAccAliCloudKVStoreRedisInstance_node_type_double(t *testing.T) {
	testAccKvstoreNodeType(t, "double", "double", "redis.amber.master.small.multithread", "", "double")
}

func testAccKvstoreNodeType(t *testing.T, nodeType, expectedNodeType, instanceClass, productType, inventoryNodeType string) {
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
			t.Logf("Requested InstanceClass=%s NodeType=%s; returned InstanceClass=%v RealInstanceClass=%v NodeType=%v ArchitectureType=%v", instanceClass, nodeType, instance["InstanceClass"], instance["RealInstanceClass"], instance["NodeType"], instance["ArchitectureType"])
			if instance["NodeType"] != expectedNodeType {
				return fmt.Errorf("expected remote NodeType %q, got %v", expectedNodeType, instance["NodeType"])
			}
			if instance["InstanceClass"] != instanceClass {
				return fmt.Errorf("expected remote InstanceClass %q, got %v", instanceClass, instance["InstanceClass"])
			}
			instanceId = id
			return nil
		},
		resource.TestCheckResourceAttr(resourceId, "instance_class", instanceClass),
		resource.TestCheckResourceAttr(resourceId, "node_type", expectedNodeType),
	)
	importCheck := func(states []*terraform.InstanceState) error {
		if len(states) != 1 || states[0].ID != instanceId {
			return fmt.Errorf("expected the existing instance to be imported")
		}
		if states[0].Attributes["node_type"] != expectedNodeType {
			return fmt.Errorf("expected imported node_type %q, got %q", expectedNodeType, states[0].Attributes["node_type"])
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
			// The import runner hands plugin-testing states to this closure;
			// the resource's Diff is an SDKv2 helper/schema API, so the state
			// crosses forks through a shallow copy (structurally identical).
			prior := &sdktf.InstanceState{
				ID:         states[0].ID,
				Attributes: states[0].Attributes,
				Ephemeral: sdktf.EphemeralState{
					ConnInfo: states[0].Ephemeral.ConnInfo,
					Type:     states[0].Ephemeral.Type,
				},
				Meta: states[0].Meta,
			}
			diff, err := resourceAliCloudKvstoreInstance().Diff(context.Background(), prior, sdktf.NewResourceConfigRaw(config), testAccProvider.Meta())
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
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
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
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
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
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
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
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
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
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
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

func TestUnitKvstoreCreateLockRetry(t *testing.T) {
	vpcEndpoint := sdkendpoints.GetEndpointFromMap("cn-hangzhou", "vpc")
	t.Cleanup(func() { sdkendpoints.AddEndpointMapping("cn-hangzhou", "vpc", vpcEndpoint) })
	for _, variable := range []string{"HTTP_PROXY", "HTTPS_PROXY", "http_proxy", "https_proxy", "ALL_PROXY", "all_proxy", "TF_ENDPOINT_PATH"} {
		t.Setenv(variable, "")
	}
	t.Setenv("NO_PROXY", "*")
	t.Setenv("no_proxy", "*")
	for _, test := range []struct {
		name      string
		firstCode string
		calls     int
		zoneField string
		zone      string
		vswitch   bool
	}{
		{name: "lock conflict", firstCode: "CanNotAcquireLock", calls: 2},
		{name: "invalid parameter", firstCode: "InvalidParameter", calls: 1},
		{name: "VSwitch zone fallback", firstCode: "InvalidParameter", calls: 1, zone: "cn-hangzhou-j", vswitch: true},
		{name: "explicit zone", firstCode: "InvalidParameter", calls: 1, zoneField: "zone_id", zone: "cn-hangzhou-i", vswitch: true},
		{name: "legacy zone", firstCode: "InvalidParameter", calls: 1, zoneField: "availability_zone", zone: "cn-hangzhou-i", vswitch: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			var tokens []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse request: %v", err)
				}
				if r.Form.Get("Action") == "DescribeVSwitchAttributes" {
					w.Header().Set("Content-Type", "application/json")
					fmt.Fprint(w, `{"VSwitchId":"vsw-test","VpcId":"vpc-test","ZoneId":"cn-hangzhou-j"}`)
					return
				}
				if r.Form.Get("Action") != "CreateInstance" {
					t.Errorf("unexpected action: %s", r.Form.Get("Action"))
				}
				if got := r.Form.Get("ZoneId"); got != test.zone {
					t.Errorf("CreateInstance ZoneId = %q, want %q", got, test.zone)
				}
				tokens = append(tokens, r.Form.Get("Token"))
				t.Logf("received create request %d", len(tokens))
				code := "InvalidParameter"
				if len(tokens) == 1 {
					code = test.firstCode
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"Code":%q,"Message":"test failure","RequestId":"test"}`, code)
			}))
			defer server.Close()

			credential, err := credentials.NewCredential(new(credentials.Config).
				SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
			if err != nil {
				t.Fatal(err)
			}
			var endpoints, signVersion sync.Map
			endpoints.Store("r_kvstore", strings.TrimPrefix(server.URL, "http://"))
			endpoints.Store("vpc", strings.TrimPrefix(server.URL, "http://"))
			config := &connectivity.Config{
				Region: connectivity.Hangzhou, RegionId: "cn-hangzhou", Protocol: "HTTP",
				AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
				Endpoints: &endpoints, SignVersion: &signVersion,
				AccountType: "test", SkipRegionValidation: true,
				MaxRetryTimeout: 15, ClientReadTimeout: 1000, ClientConnectTimeout: 1000,
			}
			client, err := config.Client()
			if err != nil {
				t.Fatal(err)
			}
			d := schema.TestResourceDataRaw(t, resourceAliCloudKvstoreInstance().Schema, map[string]interface{}{
				"instance_class": "redis.shard.with.proxy.small.ce",
				"engine_version": "6.0",
			})
			if test.vswitch {
				d.Set("vswitch_id", "vsw-test")
			}
			if test.zoneField == "zone_id" {
				d.Set("zone_id", test.zone)
			} else if test.zoneField == "availability_zone" {
				d.Set("availability_zone", test.zone)
			}
			// A terminal response after the conflict exercises the create retry loop
			// without creating an instance or entering the instance-state waiter.
			err = resourceAliCloudKvstoreInstanceCreate(d, client)
			if err == nil || !IsExpectedErrors(err, []string{"InvalidParameter"}) {
				t.Fatalf("expected the terminal API error, got %v", err)
			}
			if len(tokens) != test.calls {
				t.Fatalf("sent %d requests, want %d", len(tokens), test.calls)
			}
			for _, token := range tokens {
				if token == "" || token != tokens[0] {
					t.Fatal("create retries must reuse the nonempty idempotency token")
				}
			}
			if d.Id() != "" {
				t.Fatal("failed creation must not set an instance ID")
			}
		})
	}
}

func TestUnitKvstoreNodeTypeCompatibility(t *testing.T) {
	r := resourceAliCloudKvstoreInstance()
	for _, test := range []struct {
		nodeType string
		alias    string
		other    string
	}{
		{nodeType: "single", alias: "STAND_ALONE", other: "double"},
		{nodeType: "double", alias: "MASTER_SLAVE", other: "single"},
		{nodeType: "MASTER_SLAVE", alias: "double", other: "STAND_ALONE"},
		{nodeType: "STAND_ALONE", alias: "single", other: "MASTER_SLAVE"},
	} {
		t.Run(test.nodeType, func(t *testing.T) {
			config := sdktf.NewResourceConfigRaw(map[string]interface{}{"node_type": test.nodeType})
			diags := r.Validate(config)
			assert.Empty(t, diags, "supported node types must not produce a deprecation warning")

			data := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"node_type": test.nodeType})
			data.SetId("test-instance")
			state := data.State()
			for _, config := range []map[string]interface{}{{}, {"node_type": test.nodeType}, {"node_type": test.alias}} {
				diff, err := r.Diff(context.Background(), state, sdktf.NewResourceConfigRaw(config), nil)
				if !assert.NoError(t, err) || diff == nil {
					continue
				}
				assert.NotContains(t, diff.Attributes, "node_type", "existing and imported node types must remain stable")
				assert.False(t, diff.RequiresNew(), "equivalent node-type aliases must not replace the instance")
			}
			diff, err := r.Diff(context.Background(), state, sdktf.NewResourceConfigRaw(map[string]interface{}{"node_type": test.other}), nil)
			if assert.NoError(t, err) && assert.NotNil(t, diff) && assert.Contains(t, diff.Attributes, "node_type") {
				assert.True(t, diff.Attributes["node_type"].RequiresNew)
			}
		})
	}
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{})
	_, configured := d.GetOk("node_type")
	assert.False(t, configured, "omission must preserve service-side default selection")
	for _, values := range [][2]string{{"", "single"}, {"single", ""}, {"", ""}, {"unknown", "single"}, {"single", "unknown"}, {"unknown", "unknown"}, {"master_slave", "double"}, {"single", "stand_alone"}} {
		if suppress := r.Schema["node_type"].DiffSuppressFunc; suppress != nil {
			assert.False(t, suppress("node_type", values[0], values[1], d), "empty, unknown and differently cased values must not be suppressed")
		}
	}
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
		{
			name:         "fallback returns empty when response is nil",
			describeFunc: func(string) (map[string]interface{}, error) { return nil, nil },
		},
		{
			name: "fallback returns empty when EngineVersion is nil",
			describeFunc: func(string) (map[string]interface{}, error) {
				return map[string]interface{}{"EngineVersion": nil}, nil
			},
		},
		{
			name: "fallback preserves empty EngineVersion",
			describeFunc: func(string) (map[string]interface{}, error) {
				return map[string]interface{}{"EngineVersion": ""}, nil
			},
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

func TestUnitKvstoreVpcBasicDependenceSecurityGroupIsolation(t *testing.T) {
	name := "tf-testAccKvstoreVpcFixture"
	networkOnly := AliCloudKVStoreRedisInstanceVpcBasicDependence0(name)
	assert.NotContains(t, networkOnly, "alicloud_security_group")
	assert.Contains(t, networkOnly, `data "alicloud_vpcs" "default"`)
	assert.Contains(t, networkOnly, `data "alicloud_vswitches" "default"`)
	assert.Contains(t, networkOnly, `data "alicloud_vswitches" "update"`)

	withSecurityGroup := AliCloudKVStoreRedisInstanceVpcSecurityGroupDependence0(name)
	assert.True(t, strings.HasPrefix(withSecurityGroup, networkOnly))
	securityGroup := strings.TrimPrefix(withSecurityGroup, networkOnly)
	assert.Equal(t, 1, strings.Count(securityGroup, `resource "alicloud_security_group" "default"`))
	assert.Contains(t, securityGroup, fmt.Sprintf("name                = %q", name))
	assert.Contains(t, securityGroup, `vpc_id              = data.alicloud_vpcs.default.ids.0`)
	assert.Contains(t, securityGroup, `inner_access_policy = "Accept"`)
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

resource "alicloud_redis_global_security_ip_group" "default" {
  count                = 3
  global_ip_group_name = "tfacc${count.index}"
  global_ip_list       = "192.168.0.1,10.10.10.10,172.16.0.1"
}
	`)
}

func AliCloudKVStoreRedisInstanceVpcSecurityGroupDependence0(name string) string {
	return AliCloudKVStoreRedisInstanceVpcBasicDependence0(name) + fmt.Sprintf(`
resource "alicloud_security_group" "default" {
  inner_access_policy = "Accept"
  name                = %q
  vpc_id              = data.alicloud_vpcs.default.ids.0
}
`, name)
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

func TestAccAliCloudKVStoreRedisInstance_elastic_burst(t *testing.T) {
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
	name := fmt.Sprintf("tf-testAccKvstoreRedisInstanceElasticBurst%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccKvstoreClassicDependence("Redis", "5.0", "PostPaid", "redis.amber.logic.sharding.1g.2db.0rodb.6proxy.multithread", "redis.amber.logic.sharding.2g.2db.0rodb.6proxy.multithread"))
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
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
					"bandwidth_burst":  "true",
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
						"bandwidth_burst":  "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_burst": "false",
					"additional_bandwidth": []map[string]interface{}{
						{"node_id": "db-0", "bandwidth": 20},
						{"node_id": "db-1", "bandwidth": 20},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_burst": "false",
					}),
					resource.TestCheckResourceAttr(resourceId, "additional_bandwidth.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "additional_bandwidth.*", map[string]string{
						"node_id":   "db-0",
						"bandwidth": "20",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "additional_bandwidth.*", map[string]string{
						"node_id":   "db-1",
						"bandwidth": "20",
					}),
					testAccCheckAdditionalBandwidthShardsFromApi(resourceId, map[string]int{"db-0": 20, "db-1": 20}),
				),
			},
			{
				// Non-uniform per-shard values: each block carries its own value.
				Config: testAccConfig(map[string]interface{}{
					"additional_bandwidth": []map[string]interface{}{
						{"node_id": "db-0", "bandwidth": 30},
						{"node_id": "db-1", "bandwidth": 40},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "additional_bandwidth.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "additional_bandwidth.*", map[string]string{
						"node_id":   "db-0",
						"bandwidth": "30",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "additional_bandwidth.*", map[string]string{
						"node_id":   "db-1",
						"bandwidth": "40",
					}),
					testAccCheckAdditionalBandwidthShardsFromApi(resourceId, map[string]int{"db-0": 30, "db-1": 40}),
				),
			},
			{
				// An import has no configuration to key the read-back on, so every shard with
				// purchased bandwidth is reported under its node name.
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "business_info", "coupon_no", "effective_time", "force_upgrade", "global_instance_id", "order_type", "password", "period", "enable_public", "security_ip_group_attribute", "enable_backup_log", "bandwidth_burst", "additional_bandwidth"},
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected one imported state, got %d", len(states))
					}
					instanceId := states[0].Attributes["id"]
					want := map[string]int{instanceId + "-db-0": 30, instanceId + "-db-1": 40}
					if got := importedAdditionalBandwidth(states[0]); !reflect.DeepEqual(want, got) {
						return fmt.Errorf("expected the imported per-shard additional bandwidth %v, got %v", want, got)
					}
					return nil
				},
			},
			{
				// Dropping a block leaves the bandwidth purchased for that shard in place; only
				// the remaining shard is managed and reported.
				Config: testAccConfig(map[string]interface{}{
					"additional_bandwidth": []map[string]interface{}{
						{"node_id": "db-0", "bandwidth": 20},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "additional_bandwidth.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "additional_bandwidth.*", map[string]string{
						"node_id":   "db-0",
						"bandwidth": "20",
					}),
					testAccCheckAdditionalBandwidthShardsFromApi(resourceId, map[string]int{"db-0": 20, "db-1": 40}),
				),
			},
			{
				// A node_id that does not name exactly one data shard is rejected before the
				// paid request is sent.
				Config: testAccConfig(map[string]interface{}{
					"additional_bandwidth": []map[string]interface{}{
						{"node_id": "db-9", "bandwidth": 10},
					},
				}),
				ExpectError: regexp.MustCompile("does not match exactly one data shard"),
			},
			{
				// A block without bandwidth has no meaning and is rejected instead of being
				// sent as zero (which would release the purchased bandwidth).
				Config: testAccConfig(map[string]interface{}{
					"additional_bandwidth": []map[string]interface{}{
						{"node_id": "db-0"},
					},
				}),
				ExpectError: regexp.MustCompile("additional_bandwidth requires bandwidth in every block"),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"additional_bandwidth": []map[string]interface{}{
						{"node_id": "db-0", "bandwidth": 30},
						{"node_id": "db-1", "bandwidth": 30},
					},
				}),
				// The out-of-band change made by the check below leaves the plan non-empty.
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "additional_bandwidth.#", "2"),
					testAccCheckAdditionalBandwidthShardsFromApi(resourceId, map[string]int{"db-0": 30, "db-1": 30}),
					// Change the purchased bandwidth of one shard out of band; the next step's
					// refresh must see the drift and reconcile it back to the configured value.
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceId]
						if !ok {
							return fmt.Errorf("can't find resource by id: %s", resourceId)
						}
						return testAccSetAdditionalBandwidthOutOfBand(rs.Primary.ID, "db-1", 50)
					},
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"additional_bandwidth": []map[string]interface{}{
						{"node_id": "db-0", "bandwidth": 30},
						{"node_id": "db-1", "bandwidth": 30},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "additional_bandwidth.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "additional_bandwidth.*", map[string]string{
						"node_id":   "db-1",
						"bandwidth": "30",
					}),
					testAccCheckAdditionalBandwidthShardsFromApi(resourceId, map[string]int{"db-0": 30, "db-1": 30}),
				),
			},
			{
				// Explicit "0" is the documented release path for the purchased bandwidth.
				Config: testAccConfig(map[string]interface{}{
					"additional_bandwidth": []map[string]interface{}{
						{"node_id": "db-0", "bandwidth": 0},
						{"node_id": "db-1", "bandwidth": 0},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "additional_bandwidth.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "additional_bandwidth.*", map[string]string{
						"node_id":   "db-0",
						"bandwidth": "0",
					}),
					testAccCheckAdditionalBandwidthShardsFromApi(resourceId, map[string]int{"db-0": 0, "db-1": 0}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_burst":      REMOVEKEY,
					"additional_bandwidth": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_burst": REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "business_info", "coupon_no", "effective_time", "force_upgrade", "global_instance_id", "order_type", "password", "period", "enable_public", "security_ip_group_attribute", "enable_backup_log", "bandwidth_burst", "additional_bandwidth"},
			},
		},
	})
}

// testAccSetAdditionalBandwidthOutOfBand changes the purchased bandwidth of one shard outside
// Terraform and verifies the change through DescribeRoleZoneInfo, so a later Terraform run can be
// shown to detect and reconcile the drift.
func testAccSetAdditionalBandwidthOutOfBand(instanceId, shardToken string, bandwidth int) error {
	extras, err := testAccAdditionalBandwidthExtras(instanceId)
	if err != nil {
		return err
	}
	shard, err := testAccAdditionalBandwidthShardName(extras, shardToken)
	if err != nil {
		return err
	}

	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"InstanceId": instanceId,
		"ChargeType": "PostPaid",
		"Bandwidth":  strconv.Itoa(bandwidth),
		"NodeId":     shard,
	}
	if _, err := client.RpcPost("R-kvstore", "2015-01-01", "EnableAdditionalBandwidth", nil, request, false); err != nil {
		return fmt.Errorf("out-of-band EnableAdditionalBandwidth(%d) on %s failed: %s", bandwidth, shard, err)
	}

	service := R_kvstoreService{client}
	// The instance must be back to Normal before another bandwidth request is accepted;
	// waiting only for the order flag leaves it in a transient status.
	statusConf := BuildStateConf([]string{}, []string{"Normal"}, 15*time.Minute, 10*time.Second, service.KvstoreInstanceAttributeRefreshFunc(instanceId, "InstanceStatus"))
	if _, err := statusConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for the instance to become Normal after out-of-band EnableAdditionalBandwidth(%d) failed: %s", bandwidth, err)
	}
	stateConf := BuildStateConf([]string{}, []string{"true"}, 10*time.Minute, 10*time.Second, service.KvstoreInstanceAttributeRefreshFunc(instanceId, "IsOrderCompleted"))
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for out-of-band EnableAdditionalBandwidth(%d) failed: %s", bandwidth, err)
	}

	actual, err := testAccAdditionalBandwidthExtras(instanceId)
	if err != nil {
		return err
	}
	if got := int(actual[shard]); got != bandwidth {
		return fmt.Errorf("expected out-of-band server-side additional bandwidth %d on shard %s, got %d", bandwidth, shard, got)
	}
	return nil
}

// testAccAdditionalBandwidthExtras returns the additional bandwidth reported by the read-only
// DescribeRoleZoneInfo API for every data shard, keyed by shard node name.
func testAccAdditionalBandwidthExtras(instanceId string) (map[string]int64, error) {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	service := R_kvstoreService{client}
	nodes, err := service.DescribeKvstoreDataNodes(instanceId)
	if err != nil {
		return nil, err
	}

	extras := make(map[string]int64)
	for _, node := range nodes {
		if fmt.Sprint(node["InsType"]) == "3" {
			continue
		}
		if !strings.EqualFold(fmt.Sprint(node["Role"]), "master") {
			continue
		}
		defaultBandWidth, _ := strconv.ParseInt(fmt.Sprint(node["DefaultBandWidth"]), 10, 64)
		currentBandWidth, _ := strconv.ParseInt(fmt.Sprint(node["CurrentBandWidth"]), 10, 64)
		extra := currentBandWidth - defaultBandWidth
		if extra < 0 {
			extra = 0
		}
		extras[fmt.Sprint(node["InsName"])] = extra
	}
	return extras, nil
}

// testAccAdditionalBandwidthShardName resolves a configured node_id the way an operator reads the
// topology: the shard node name (or its full topology ID) is used as is, and a suffix such as
// "db-0" is matched against the names reported by the API. The production resolver is deliberately
// not reused, so the assertion can detect a resolution defect.
func testAccAdditionalBandwidthShardName(extras map[string]int64, token string) (string, error) {
	name := token
	if idx := strings.Index(name, "#"); idx >= 0 {
		name = name[:idx]
	}
	if _, ok := extras[name]; ok {
		return name, nil
	}

	matched := make([]string, 0, 1)
	for shard := range extras {
		if strings.HasSuffix(shard, "-"+name) {
			matched = append(matched, shard)
		}
	}
	if len(matched) == 1 {
		return matched[0], nil
	}
	names := make([]string, 0, len(extras))
	for shard := range extras {
		names = append(names, shard)
	}
	sort.Strings(names)
	return "", fmt.Errorf("node_id %q does not match exactly one data shard, data shards: %v", token, names)
}

// testAccCheckAdditionalBandwidthShardsFromApi asserts the additional bandwidth reported by the
// read-only DescribeRoleZoneInfo API for each configured shard. The expected map is keyed by the
// node_id used in the configuration.
func testAccCheckAdditionalBandwidthShardsFromApi(resourceId string, expected map[string]int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("can't find resource by id: %s", resourceId)
		}
		extras, err := testAccAdditionalBandwidthExtras(rs.Primary.ID)
		if err != nil {
			return err
		}
		for token, want := range expected {
			shard, err := testAccAdditionalBandwidthShardName(extras, token)
			if err != nil {
				return err
			}
			if got := int(extras[shard]); got != want {
				return fmt.Errorf("expected server-side additional bandwidth %d on shard %s (%s), got %d", want, token, shard, got)
			}
		}
		return nil
	}
}

// importedAdditionalBandwidth collects the additional_bandwidth blocks of an imported state, keyed
// by node_id.
func importedAdditionalBandwidth(state *terraform.InstanceState) map[string]int {
	blocks := make(map[string]int)
	for key, value := range state.Attributes {
		if !strings.HasPrefix(key, "additional_bandwidth.") || !strings.HasSuffix(key, ".node_id") {
			continue
		}
		bandwidth, _ := strconv.Atoi(state.Attributes[strings.TrimSuffix(key, ".node_id")+".bandwidth"])
		blocks[value] = bandwidth
	}
	return blocks
}

func TestUnitKvstoreAdditionalBandwidthState(t *testing.T) {
	node := func(insName, role string, defaultBandWidth, currentBandWidth int64) map[string]interface{} {
		return map[string]interface{}{
			"CustinsId":        insName + "-shard",
			"InsName":          insName,
			"Role":             role,
			"DefaultBandWidth": defaultBandWidth,
			"CurrentBandWidth": currentBandWidth,
		}
	}
	twoShards := []map[string]interface{}{
		node("r-x-db-0", "MASTER", 96, 116),
		node("r-x-db-1", "MASTER", 96, 126),
	}

	state, ok := buildAdditionalBandwidthState(nil, nil)
	assert.False(t, ok)
	assert.Nil(t, state)

	// An import (no managed shards) reports every shard with purchased bandwidth, keyed by name.
	state, ok = buildAdditionalBandwidthState(twoShards, nil)
	assert.True(t, ok)
	assert.Equal(t, []map[string]interface{}{
		{"node_id": "r-x-db-0", "bandwidth": 20},
		{"node_id": "r-x-db-1", "bandwidth": 30},
	}, state)

	// A shard without purchased bandwidth is not reported by an import.
	state, ok = buildAdditionalBandwidthState([]map[string]interface{}{
		node("r-x-db-0", "MASTER", 96, 116),
		node("r-x-db-1", "MASTER", 96, 96),
	}, nil)
	assert.True(t, ok)
	assert.Equal(t, []map[string]interface{}{{"node_id": "r-x-db-0", "bandwidth": 20}}, state)

	// Managed shards are reported with their own value and keep the configured node_id, so a
	// shard that the configuration does not manage does not appear.
	state, ok = buildAdditionalBandwidthState(twoShards, []kvstoreAdditionalBandwidth{
		{NodeId: "db-1", Bandwidth: 0},
		{NodeId: "r-x-db-0#100", Bandwidth: 0},
	})
	assert.True(t, ok)
	assert.Equal(t, []map[string]interface{}{
		{"node_id": "db-1", "bandwidth": 30},
		{"node_id": "r-x-db-0#100", "bandwidth": 20},
	}, state)

	// A shard reports a master and a slave row; the master row wins, and read-only nodes
	// (InsType 3) are not bandwidth targets.
	readOnly := node("r-x-db-1", "MASTER", 96, 200)
	readOnly["InsType"] = 3
	state, ok = buildAdditionalBandwidthState([]map[string]interface{}{
		node("r-x-db-0", "SLAVE", 96, 140),
		node("r-x-db-0", "MASTER", 96, 116),
		readOnly,
	}, nil)
	assert.True(t, ok)
	assert.Equal(t, []map[string]interface{}{{"node_id": "r-x-db-0", "bandwidth": 20}}, state)

	// Bandwidth below the specification default is reported as no additional bandwidth.
	state, ok = buildAdditionalBandwidthState([]map[string]interface{}{
		node("r-x-db-0", "MASTER", 96, 90),
	}, nil)
	assert.True(t, ok)
	assert.Equal(t, []map[string]interface{}{}, state)

	// A managed shard without a bandwidth row means the listing is incomplete: the caller keeps
	// the previous state instead of reporting a smaller value.
	state, ok = buildAdditionalBandwidthState(twoShards[:1], []kvstoreAdditionalBandwidth{
		{NodeId: "db-1", Bandwidth: 30},
	})
	assert.False(t, ok)
	assert.Nil(t, state)

	// A managed node_id that no longer maps to exactly one shard cannot be refreshed.
	state, ok = buildAdditionalBandwidthState([]map[string]interface{}{
		node("r-x-db-0", "MASTER", 96, 116),
		node("r-y-db-0", "MASTER", 96, 116),
	}, []kvstoreAdditionalBandwidth{{NodeId: "db-0", Bandwidth: 20}})
	assert.False(t, ok)
	assert.Nil(t, state)

	// An unambiguous suffix is resolved to its shard.
	state, ok = buildAdditionalBandwidthState(twoShards, []kvstoreAdditionalBandwidth{{NodeId: "db-0", Bandwidth: 20}})
	assert.True(t, ok)
	assert.Equal(t, []map[string]interface{}{{"node_id": "db-0", "bandwidth": 20}}, state)

	// A standard instance configures the bandwidth without a node_id; its single data shard is
	// the target, including for a refreshed value.
	state, ok = buildAdditionalBandwidthState([]map[string]interface{}{node("r-x", "MASTER", 96, 116)}, []kvstoreAdditionalBandwidth{{Bandwidth: 20}})
	assert.True(t, ok)
	assert.Equal(t, []map[string]interface{}{{"node_id": "", "bandwidth": 20}}, state)

	state, ok = buildAdditionalBandwidthState([]map[string]interface{}{node("r-x", "MASTER", 96, 136)}, []kvstoreAdditionalBandwidth{{Bandwidth: 20}})
	assert.True(t, ok)
	assert.Equal(t, []map[string]interface{}{{"node_id": "", "bandwidth": 40}}, state)

	// A node_id-less block cannot be refreshed when the listing has several data shards.
	state, ok = buildAdditionalBandwidthState(twoShards, []kvstoreAdditionalBandwidth{{Bandwidth: 20}})
	assert.False(t, ok)
	assert.Nil(t, state)
}

func TestUnitKvstoreAdditionalBandwidthStandardRefresh(t *testing.T) {
	node := func(insName string, defaultBandWidth, currentBandWidth int64) map[string]interface{} {
		return map[string]interface{}{
			"InsName":          insName,
			"Role":             "MASTER",
			"DefaultBandWidth": defaultBandWidth,
			"CurrentBandWidth": currentBandWidth,
		}
	}
	stateOf := func(t *testing.T, d *schema.ResourceData) map[string]int {
		value, ok := d.Get("additional_bandwidth").(*schema.Set)
		assert.True(t, ok)
		blocks := make(map[string]int)
		for _, item := range value.List() {
			entry := item.(map[string]interface{})
			blocks[fmt.Sprint(entry["node_id"])] = entry["bandwidth"].(int)
		}
		return blocks
	}
	newResourceData := func(raw map[string]interface{}) *schema.ResourceData {
		d := schema.TestResourceDataRaw(t, resourceAliCloudKvstoreInstance().Schema, raw)
		d.SetId("r-x")
		return d
	}

	// A standard instance configures the bandwidth without a node_id. The refresh keeps that block
	// instead of treating the instance as an import, so the state still matches the configuration.
	d := newResourceData(map[string]interface{}{
		"additional_bandwidth": []interface{}{map[string]interface{}{"bandwidth": 20}},
	})
	assert.True(t, refreshAdditionalBandwidthState(d, []map[string]interface{}{node("r-x", 96, 116)}))
	assert.Equal(t, map[string]int{"": 20}, stateOf(t, d))

	// A second refresh of the same listing produces the same state, so the next plan is empty.
	assert.True(t, refreshAdditionalBandwidthState(d, []map[string]interface{}{node("r-x", 96, 116)}))
	assert.Equal(t, map[string]int{"": 20}, stateOf(t, d))

	// A bandwidth changed outside Terraform reaches the state, so the drift shows up in the plan.
	assert.True(t, refreshAdditionalBandwidthState(d, []map[string]interface{}{node("r-x", 96, 136)}))
	assert.Equal(t, map[string]int{"": 40}, stateOf(t, d))

	// An explicit zero releases the purchased bandwidth but keeps the block, so the release stays
	// in the state and does not turn into a pending creation.
	released := newResourceData(map[string]interface{}{
		"additional_bandwidth": []interface{}{map[string]interface{}{"bandwidth": 0}},
	})
	assert.True(t, refreshAdditionalBandwidthState(released, []map[string]interface{}{node("r-x", 96, 96)}))
	assert.Equal(t, map[string]int{"": 0}, stateOf(t, released))

	// Only a state without a managed block (an import) rebuilds the blocks from the node names.
	imported := newResourceData(map[string]interface{}{})
	assert.True(t, refreshAdditionalBandwidthState(imported, []map[string]interface{}{node("r-x", 96, 116)}))
	assert.Equal(t, map[string]int{"r-x": 20}, stateOf(t, imported))

	// A node_id-less block cannot be matched to a listing with several data shards, so the refresh
	// keeps the previous state rather than writing another shard's value.
	ambiguous := newResourceData(map[string]interface{}{
		"additional_bandwidth": []interface{}{map[string]interface{}{"bandwidth": 20}},
	})
	assert.False(t, refreshAdditionalBandwidthState(ambiguous, []map[string]interface{}{
		node("r-x-db-0", 96, 116),
		node("r-x-db-1", 96, 116),
	}))
	assert.Equal(t, map[string]int{"": 20}, stateOf(t, ambiguous))
}

func TestUnitKvstoreAdditionalBandwidthRequest(t *testing.T) {
	shards := []string{"r-x-db-0", "r-x-db-1", "r-x-db-10"}

	nodeIds, bandwidths, err := buildAdditionalBandwidthRequest("cluster", 3, shards, nil)
	assert.NoError(t, err)
	assert.Equal(t, "", nodeIds)
	assert.Equal(t, "", bandwidths)

	// Every bandwidth value is emitted next to its own shard, and the vectors are ordered by
	// shard name so a plan produces a stable request.
	nodeIds, bandwidths, err = buildAdditionalBandwidthRequest("cluster", 3, shards, []kvstoreAdditionalBandwidth{
		{NodeId: "r-x-db-1", Bandwidth: 30},
		{NodeId: "r-x-db-0", Bandwidth: 20},
		{NodeId: "r-x-db-10", Bandwidth: 40},
	})
	assert.NoError(t, err)
	assert.Equal(t, "r-x-db-0,r-x-db-1,r-x-db-10", nodeIds)
	assert.Equal(t, "20,30,40", bandwidths)

	// A suffix addresses one shard without repeating the instance ID.
	nodeIds, bandwidths, err = buildAdditionalBandwidthRequest("cluster", 3, shards, []kvstoreAdditionalBandwidth{
		{NodeId: "db-10", Bandwidth: 50},
	})
	assert.NoError(t, err)
	assert.Equal(t, "r-x-db-10", nodeIds)
	assert.Equal(t, "50", bandwidths)

	// The topology form is accepted, but only the shard node name is sent.
	nodeIds, bandwidths, err = buildAdditionalBandwidthRequest("cluster", 3, shards, []kvstoreAdditionalBandwidth{
		{NodeId: "r-x-db-0#123456", Bandwidth: 20},
	})
	assert.NoError(t, err)
	assert.Equal(t, "r-x-db-0", nodeIds)
	assert.Equal(t, "20", bandwidths)

	// Releasing the purchased bandwidth is an explicit zero.
	nodeIds, bandwidths, err = buildAdditionalBandwidthRequest("cluster", 3, shards, []kvstoreAdditionalBandwidth{
		{NodeId: "db-0", Bandwidth: 0},
	})
	assert.NoError(t, err)
	assert.Equal(t, "r-x-db-0", nodeIds)
	assert.Equal(t, "0", bandwidths)

	// A standard instance takes Bandwidth alone; the API does not accept NodeId for it.
	nodeIds, bandwidths, err = buildAdditionalBandwidthRequest("standard", 1, nil, []kvstoreAdditionalBandwidth{
		{Bandwidth: 20},
	})
	assert.NoError(t, err)
	assert.Equal(t, "", nodeIds)
	assert.Equal(t, "20", bandwidths)

	// A read/write splitting architecture identified only by its shard count still requires the
	// shard name, and a standard instance rejects a block that names one.
	nodeIds, bandwidths, err = buildAdditionalBandwidthRequest("", 2, []string{"r-x-db-0", "r-x-db-1"}, []kvstoreAdditionalBandwidth{
		{NodeId: "db-1", Bandwidth: 20},
	})
	assert.NoError(t, err)
	assert.Equal(t, "r-x-db-1", nodeIds)
	assert.Equal(t, "20", bandwidths)

	_, _, err = buildAdditionalBandwidthRequest("standard", 1, nil, []kvstoreAdditionalBandwidth{
		{NodeId: "r-x", Bandwidth: 20},
	})
	assert.Error(t, err)

	// A cluster instance rejects a block without a shard name, and a standard instance accepts a
	// single block only.
	_, _, err = buildAdditionalBandwidthRequest("cluster", 2, shards, []kvstoreAdditionalBandwidth{
		{Bandwidth: 20},
	})
	assert.Error(t, err)

	_, _, err = buildAdditionalBandwidthRequest("standard", 1, nil, []kvstoreAdditionalBandwidth{
		{Bandwidth: 20},
		{Bandwidth: 30},
	})
	assert.Error(t, err)

	// The same shard spelled two ways is rejected instead of being sent twice.
	_, _, err = buildAdditionalBandwidthRequest("cluster", 3, shards, []kvstoreAdditionalBandwidth{
		{NodeId: "db-0", Bandwidth: 20},
		{NodeId: "r-x-db-0#123456", Bandwidth: 30},
	})
	assert.Error(t, err)

	// A node_id that matches no shard is rejected before the API call.
	_, _, err = buildAdditionalBandwidthRequest("cluster", 3, shards, []kvstoreAdditionalBandwidth{
		{NodeId: "db-9", Bandwidth: 20},
	})
	assert.Error(t, err)

	// A suffix shared by two shards is ambiguous and rejected.
	_, _, err = buildAdditionalBandwidthRequest("cluster", 2, []string{"r-x-db-0", "r-y-db-0"}, []kvstoreAdditionalBandwidth{
		{NodeId: "db-0", Bandwidth: 20},
	})
	assert.Error(t, err)

	// A negative bandwidth value is rejected.
	_, _, err = buildAdditionalBandwidthRequest("cluster", 3, shards, []kvstoreAdditionalBandwidth{
		{NodeId: "db-0", Bandwidth: -1},
	})
	assert.Error(t, err)
}

func TestUnitKvstoreAdditionalBandwidthPagination(t *testing.T) {
	pages := map[int][]map[string]interface{}{
		1: {{"InsName": "r-x-db-0"}, {"InsName": "r-x-db-1"}},
		2: {{"InsName": "r-x-db-2"}},
	}
	var requested []int
	nodes, err := collectAdditionalBandwidthNodes(func(pageNumber, pageSize int) ([]map[string]interface{}, int, error) {
		requested = append(requested, pageNumber)
		return pages[pageNumber], 3, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2}, requested)
	assert.Len(t, nodes, 3)

	// A late page is still collected when the earlier pages are full.
	requested = nil
	nodes, err = collectAdditionalBandwidthNodes(func(pageNumber, pageSize int) ([]map[string]interface{}, int, error) {
		requested = append(requested, pageNumber)
		if pageNumber == 1 {
			return []map[string]interface{}{{"InsName": "r-x-db-0"}}, 2, nil
		}
		if pageNumber == 2 {
			return []map[string]interface{}{{"InsName": "r-x-db-1"}}, 2, nil
		}
		return nil, 2, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2}, requested)
	assert.Len(t, nodes, 2)

	// An empty page before TotalCount is reached is an incomplete listing.
	_, err = collectAdditionalBandwidthNodes(func(pageNumber, pageSize int) ([]map[string]interface{}, int, error) {
		if pageNumber == 1 {
			return []map[string]interface{}{{"InsName": "r-x-db-0"}}, 3, nil
		}
		return nil, 3, nil
	})
	assert.ErrorIs(t, err, errIncompleteAdditionalBandwidthNodes)

	// Without TotalCount a short page ends the walk.
	requested = nil
	nodes, err = collectAdditionalBandwidthNodes(func(pageNumber, pageSize int) ([]map[string]interface{}, int, error) {
		requested = append(requested, pageNumber)
		if pageNumber == 1 {
			return []map[string]interface{}{{"InsName": "r-x-db-0"}}, 0, nil
		}
		return nil, 0, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, []int{1}, requested)
	assert.Len(t, nodes, 1)

	// A page error is surfaced with the nodes collected so far.
	_, err = collectAdditionalBandwidthNodes(func(pageNumber, pageSize int) ([]map[string]interface{}, int, error) {
		if pageNumber == 2 {
			return nil, 0, fmt.Errorf("boom")
		}
		return []map[string]interface{}{{"InsName": "r-x-db-0"}}, 2, nil
	})
	assert.Error(t, err)
}
