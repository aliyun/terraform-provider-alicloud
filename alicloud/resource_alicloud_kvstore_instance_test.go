package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
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

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var insts []r_kvstore.KVStoreInstance
	req := r_kvstore.CreateDescribeInstancesRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for _, instanceType := range []string{string(KVStoreRedis), string(KVStoreMemcache)} {
		req.InstanceType = instanceType
		for {
			raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
				return rkvClient.DescribeInstances(req)
			})
			if err != nil {
				return fmt.Errorf("Error retrieving KVStore Instances: %s", err)
			}
			resp, _ := raw.(*r_kvstore.DescribeInstancesResponse)
			if resp == nil || len(resp.Instances.KVStoreInstance) < 1 {
				break
			}
			insts = append(insts, resp.Instances.KVStoreInstance...)

			if len(resp.Instances.KVStoreInstance) < PageSizeLarge {
				break
			}

			page, err := getNextpageNumber(req.PageNumber)
			if err != nil {
				return err
			}
			req.PageNumber = page
		}
	}

	sweeped := false
	for _, v := range insts {
		name := v.InstanceName
		id := v.InstanceId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping KVStore Instance: %s (%s)", name, id)
			continue
		}

		sweeped = true
		log.Printf("[INFO] Deleting KVStore Instance: %s (%s)", name, id)
		req := r_kvstore.CreateDeleteInstanceRequest()
		req.InstanceId = id
		_, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.DeleteInstance(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete KVStore Instance (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		// Waiting 30 seconds to ensure these KVStore instances have been deleted.
		time.Sleep(30 * time.Second)
	}
	return nil
}

func TestAccAlicloudKVStoreRedisInstance_vpctest(t *testing.T) {
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
			testAccPreCheckWithNoDefaultVpc(t)
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
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "auto_use_coupon", "backup_id", "business_info", "coupon_no", "dedicated_host_group_id", "effect_time", "effective_time", "force_upgrade", "global_instance", "global_instance_id", "instance_release_protection", "major_version", "modify_mode", "order_type", "password", "period", "restore_time", "src_db_instance_id", "enable_public", "security_ip_group_attribute", "security_ip_group_name", "security_ips"},
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
					"zone_id":    "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 2].id}",
					"vswitch_id": "${data.alicloud_vswitches.update.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":    CHECKSET,
						"vswitch_id": CHECKSET,
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
					"zone_id":             "${data.alicloud_zones.default.zones[length(data.alicloud_zones.default.ids) - 1].id}",
					"vswitch_id":          "${data.alicloud_vswitches.default.ids.0}",
					"maintain_start_time": "04:00Z",
					"maintain_end_time":   "06:00Z",
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
						"maintain_start_time":           "04:00Z",
						"maintain_end_time":             "06:00Z",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudKVStoreMemcacheInstance_vpctest(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, KvstoreInstanceVpcTestdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithNoDefaultVpc(t)
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
					"engine_version":   "2.8",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":    "memcache.master.small.default",
						"db_instance_name":  name,
						"instance_type":     "Memcache",
						"engine_version":    "2.8",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "acceptance test",
						"resource_group_id": CHECKSET,
						"zone_id":           CHECKSET,
						"vswitch_id":        CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "auto_use_coupon", "backup_id", "business_info", "coupon_no", "dedicated_host_group_id", "effect_time", "effective_time", "force_upgrade", "global_instance", "global_instance_id", "instance_release_protection", "major_version", "modify_mode", "order_type", "password", "period", "restore_time", "src_db_instance_id", "enable_public", "security_ip_group_attribute", "security_ip_group_name", "security_ips"},
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
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version": "4.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version": "4.0",
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
					"zone_id":    "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 2].id}",
					"vswitch_id": "${data.alicloud_vswitches.update.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":    CHECKSET,
						"vswitch_id": CHECKSET,
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
					"zone_id":             "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id}",
					"vswitch_id":          "${data.alicloud_vswitches.default.ids.0}",
					"maintain_start_time": "04:00Z",
					"maintain_end_time":   "06:00Z",
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
						"maintain_start_time":         "04:00Z",
						"maintain_end_time":           "06:00Z",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudKVStoreRedisInstance_classictest(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, RedisDbInstanceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreRedisInstanceClassicTest%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, KvstoreInstanceClassicTestdependence)
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
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "auto_use_coupon", "backup_id", "business_info", "coupon_no", "dedicated_host_group_id", "effect_time", "effective_time", "force_upgrade", "global_instance", "global_instance_id", "instance_release_protection", "major_version", "modify_mode", "order_type", "password", "period", "restore_time", "src_db_instance_id", "enable_public", "security_ip_group_attribute", "security_ip_group_name", "security_ips"},
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
					"instance_class":              "redis.master.small.default",
					"instance_release_protection": "false",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"security_ips":                []string{"10.0.0.1"},
					"db_instance_name":            name,
					"config": map[string]string{
						"appendonly":             "yes",
						"lazyfree-lazy-eviction": "yes",
						"EvictionPolicy":         "volatile-lru",
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"maintain_start_time": "04:00Z",
					"maintain_end_time":   "06:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":                "redis.master.small.default",
						"instance_release_protection":   "false",
						"resource_group_id":             CHECKSET,
						"security_ips.#":                "1",
						"db_instance_name":              name,
						"config.%":                      "3",
						"config.appendonly":             "yes",
						"config.lazyfree-lazy-eviction": "yes",
						"config.EvictionPolicy":         "volatile-lru",
						"tags.%":                        "2",
						"tags.Created":                  "TF",
						"tags.For":                      "acceptance test",
						"maintain_start_time":           "04:00Z",
						"maintain_end_time":             "06:00Z",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudKVStoreMemcacheInstance_classictest(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, RedisDbInstanceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreMemcacheInstanceClassicTest%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, KvstoreInstanceClassicTestdependence)
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
					"engine_version":   "2.8",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"zone_id":           "${data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":    "memcache.master.small.default",
						"db_instance_name":  name,
						"instance_type":     "Memcache",
						"engine_version":    "2.8",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "acceptance test",
						"resource_group_id": CHECKSET,
						"zone_id":           CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "auto_use_coupon", "backup_id", "business_info", "coupon_no", "dedicated_host_group_id", "effect_time", "effective_time", "force_upgrade", "global_instance", "global_instance_id", "instance_release_protection", "major_version", "modify_mode", "order_type", "password", "period", "restore_time", "src_db_instance_id", "enable_public", "security_ip_group_attribute", "security_ip_group_name", "security_ips"},
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
			// Currently, the memcache does not support upgrade version to 4.0 from 2.8
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"engine_version": "4.0",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"engine_version": "4.0",
			//		}),
			//	),
			//},
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
					"instance_class":              "memcache.master.small.default",
					"instance_release_protection": "false",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"security_ips":                []string{"10.0.0.1"},
					"db_instance_name":            name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"maintain_start_time": "04:00Z",
					"maintain_end_time":   "06:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class":              "memcache.master.small.default",
						"instance_release_protection": "false",
						"resource_group_id":           CHECKSET,
						"security_ips.#":              "1",
						"db_instance_name":            name,
						"tags.%":                      "2",
						"tags.Created":                "TF",
						"tags.For":                    "acceptance test",
						"maintain_start_time":         "04:00Z",
						"maintain_end_time":           "06:00Z",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudKVStoreRedisInstance_vpcmulti(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default.1"
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
			testAccPreCheckWithNoDefaultVpc(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":            "2",
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
					}),
				),
			},
		},
	})
}

func TestAccAlicloudKVStoreRedisInstance_classicmulti(t *testing.T) {
	var v r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default.1"
	ra := resourceAttrInit(resourceId, RedisDbInstanceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreRedisInstanceClassicMultiTest%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, KvstoreInstanceClassicTestdependence)
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
					"count":            "2",
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
					}),
				),
			},
		},
	})
}

var RedisDbInstanceMap = map[string]string{
	"connection_domain": CHECKSET,
}

func KvstoreInstanceVpcTestdependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_kvstore_zones" "default"{
		instance_charge_type = "PostPaid"
	}
	data "alicloud_vpcs" "default" {
		is_default = true
	}
	data "alicloud_vswitches" "default" {
  		zone_id = data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}
	data "alicloud_resource_manager_resource_groups" "default"{
	}

	data "alicloud_vswitches" "update" {
  		zone_id = data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 2].id
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}
	`)
}

func KvstoreInstanceClassicTestdependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_kvstore_zones" "default"{
		instance_charge_type = "PostPaid"
	}
	data "alicloud_resource_manager_resource_groups" "default"{
	}
	`)
}
