package alicloud

import (
	"fmt"
	"log"
	"testing"

	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

var redisInstanceConnectionDomainRegexp = "^r-[a-z0-9]+.redis[.a-z-0-9]*.rds.aliyuncs.com"
var redisInstanceClassForTest = "redis.master.small.default"
var redisInstanceClassForTestUpdateClass = "redis.master.mid.default"
var memcacheInstanceConnectionDomainRegexp = "^m-[a-z0-9]+.memcache[.a-z-0-9]*.rds.aliyuncs.com"
var memcacheInstanceClassForTest = "memcache.master.small.default"
var memcacheInstanceClassForTestUpdateClass = "memcache.master.mid.default"

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
		"testAcc",
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

			if page, err := getNextpageNumber(req.PageNumber); err != nil {
				return err
			} else {
				req.PageNumber = page
			}
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
		// Waiting 30 seconds to eusure these KVStore instances have been deleted.
		time.Sleep(30 * time.Second)
	}
	return nil
}

func TestAccAlicloudKVStoreRedisInstance_classictest(t *testing.T) {
	var instance *r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKVstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_classic(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":        "tf-testAccKVStoreInstance_classic",
						"instance_class":       redisInstanceClassForTest,
						"password":             NOSET,
						"availability_zone":    CHECKSET,
						"instance_charge_type": string(PostPaid),
						"period":               NOSET,
						"instance_type":        string(KVStoreRedis),
						"vswitch_id":           "",
						"engine_version":       string(KVStore2Dot8),
						"connection_domain":    REGEXMATCH + redisInstanceConnectionDomainRegexp,
						"private_ip":           "",
						"backup_id":            NOSET,
						"security_ips.#":       "1",
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_classicUpdateParameter(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "1",
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_classicAddParameter(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "2",
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_classicDeleteParameter(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "1",
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_classicUpdateSecuirtyIps(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips.#": "2",
						"parameters.#":   REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_classicUpdateClass(string(KVStoreRedis), redisInstanceClassForTestUpdateClass, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class": redisInstanceClassForTestUpdateClass,
						"security_ips.#": "1",
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_classicUpdateAttr(string(KVStoreRedis), redisInstanceClassForTestUpdateClass, string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password":       CHECKSET,
						"engine_version": string(KVStore4Dot0),
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_classicUpdateAll(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":  "tf-testAccKVStoreInstance_classicUpdateAll",
						"instance_class": redisInstanceClassForTest,
						"engine_version": string(KVStore4Dot0),
						"security_ips.#": "2",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudKVStoreMemcacheInstance_classictest(t *testing.T) {
	var instance *r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKVstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_classic(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":        "tf-testAccKVStoreInstance_classic",
						"instance_class":       memcacheInstanceClassForTest,
						"password":             NOSET,
						"availability_zone":    CHECKSET,
						"instance_charge_type": string(PostPaid),
						"period":               NOSET,
						"instance_type":        string(KVStoreMemcache),
						"vswitch_id":           "",
						"engine_version":       string(KVStore2Dot8),
						"connection_domain":    REGEXMATCH + memcacheInstanceConnectionDomainRegexp,
						"private_ip":           "",
						"backup_id":            NOSET,
						"security_ips.#":       "1",
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_classicUpdateParameter(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "1",
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_classicAddParameter(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "2",
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_classicDeleteParameter(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "1",
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_classicUpdateSecuirtyIps(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips.#": "2",
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_classicUpdateClass(string(KVStoreMemcache), memcacheInstanceClassForTestUpdateClass, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class": memcacheInstanceClassForTestUpdateClass,
						"security_ips.#": "1",
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_classicUpdateAttr(string(KVStoreMemcache), memcacheInstanceClassForTestUpdateClass, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password":       CHECKSET,
						"engine_version": string(KVStore2Dot8),
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_classicUpdateAll(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":  "tf-testAccKVStoreInstance_classicUpdateAll",
						"instance_class": memcacheInstanceClassForTest,
						"engine_version": string(KVStore2Dot8),
						"security_ips.#": "2",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudKVStoreRedisInstance_vpctest(t *testing.T) {
	var instance *r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKVstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_vpc(KVStoreCommonTestCase, redisInstanceClassForTest, string(KVStoreRedis), string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":        "tf-testAccKVStoreInstance_vpc",
						"instance_class":       redisInstanceClassForTest,
						"password":             NOSET,
						"availability_zone":    CHECKSET,
						"instance_charge_type": string(PostPaid),
						"period":               NOSET,
						"instance_type":        string(KVStoreRedis),
						"vswitch_id":           CHECKSET,
						"engine_version":       string(KVStore4Dot0),
						"connection_domain":    REGEXMATCH + redisInstanceConnectionDomainRegexp,
						"private_ip":           "172.16.0.10",
						"backup_id":            NOSET,
						"security_ips.#":       "1",
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_vpcUpdateSecurityIps(KVStoreCommonTestCase, redisInstanceClassForTest, string(KVStoreRedis), string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips.#": "2",
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_vpcUpdateClass(KVStoreCommonTestCase, redisInstanceClassForTestUpdateClass, string(KVStoreRedis), string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class": redisInstanceClassForTestUpdateClass,
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_vpcUpdateVpcAuthMode(KVStoreCommonTestCase, redisInstanceClassForTestUpdateClass, string(KVStoreRedis), string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_auth_mode": "Close",
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_vpcUpdateParameter(KVStoreCommonTestCase, redisInstanceClassForTestUpdateClass, string(KVStoreRedis), string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "1",
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_vpcAddParameter(KVStoreCommonTestCase, redisInstanceClassForTestUpdateClass, string(KVStoreRedis), string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "2",
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_vpcDeleteParameter(KVStoreCommonTestCase, redisInstanceClassForTestUpdateClass, string(KVStoreRedis), string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "1",
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_vpcUpdateAll(KVStoreCommonTestCase, redisInstanceClassForTest, string(KVStoreRedis), string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":  "tf-testAccKVStoreInstance_vpcUpdateAll",
						"password":       CHECKSET,
						"instance_class": redisInstanceClassForTest,
						"security_ips.#": "1",
					}),
				),
			},
		},
	})
}

// Currently Memcache instance only supports engine version 2.8.
func TestAccAlicloudKVStoreMemcacheInstance_vpctest(t *testing.T) {
	var instance *r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKVstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_vpc(KVStoreCommonTestCase, memcacheInstanceClassForTest, string(KVStoreMemcache), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":        "tf-testAccKVStoreInstance_vpc",
						"instance_class":       memcacheInstanceClassForTest,
						"password":             NOSET,
						"availability_zone":    CHECKSET,
						"instance_charge_type": string(PostPaid),
						"period":               NOSET,
						"instance_type":        string(KVStoreMemcache),
						"vswitch_id":           CHECKSET,
						"engine_version":       string(KVStore2Dot8),
						"connection_domain":    REGEXMATCH + memcacheInstanceConnectionDomainRegexp,
						"private_ip":           "172.16.0.10",
						"backup_id":            NOSET,
						"security_ips.#":       "1",
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_vpcUpdateSecurityIps(KVStoreCommonTestCase, memcacheInstanceClassForTest, string(KVStoreMemcache), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips.#": "2",
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_vpcUpdateClass(KVStoreCommonTestCase, memcacheInstanceClassForTestUpdateClass, string(KVStoreMemcache), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class": memcacheInstanceClassForTestUpdateClass,
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_vpcUpdateParameter(KVStoreCommonTestCase, memcacheInstanceClassForTestUpdateClass, string(KVStoreMemcache), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "1",
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_vpcAddParameter(KVStoreCommonTestCase, memcacheInstanceClassForTestUpdateClass, string(KVStoreMemcache), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "2",
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_vpcDeleteParameter(KVStoreCommonTestCase, memcacheInstanceClassForTestUpdateClass, string(KVStoreMemcache), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "1",
					}),
				),
			},
			{
				Config: testAccKVStoreInstance_vpcUpdateAll(KVStoreCommonTestCase, memcacheInstanceClassForTest, string(KVStoreMemcache), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":  "tf-testAccKVStoreInstance_vpcUpdateAll",
						"instance_class": memcacheInstanceClassForTest,
						"parameters.#":   "1",
						"security_ips.#": "1",
						"password":       CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudKVStoreRedisInstance_vpcmulti(t *testing.T) {
	var instance *r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default.9"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKVstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_vpcmulti(KVStoreCommonTestCase, redisInstanceClassForTest, string(KVStoreRedis), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":        "tf-testAccKVStoreInstance_vpc",
						"instance_class":       redisInstanceClassForTest,
						"password":             NOSET,
						"availability_zone":    CHECKSET,
						"instance_charge_type": string(PostPaid),
						"period":               NOSET,
						"instance_type":        string(KVStoreRedis),
						"vswitch_id":           CHECKSET,
						"engine_version":       string(KVStore2Dot8),
						"connection_domain":    REGEXMATCH + redisInstanceConnectionDomainRegexp,
						"private_ip":           "172.16.0.10",
						"backup_id":            NOSET,
						"security_ips.#":       "1",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudKVStoreRedisInstance_classicmulti(t *testing.T) {
	var instance *r_kvstore.DBInstanceAttribute
	resourceId := "alicloud_kvstore_instance.default.9"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKVstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_classicmulti(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":        "tf-testAccKVStoreInstance_classic",
						"instance_class":       redisInstanceClassForTest,
						"password":             NOSET,
						"availability_zone":    CHECKSET,
						"instance_charge_type": string(PostPaid),
						"period":               NOSET,
						"instance_type":        string(KVStoreRedis),
						"vswitch_id":           "",
						"engine_version":       string(KVStore2Dot8),
						"connection_domain":    REGEXMATCH + redisInstanceConnectionDomainRegexp,
						"private_ip":           "",
						"backup_id":            NOSET,
						"security_ips.#":       "1",
					}),
				),
			},
		},
	})
}

func testAccCheckKVStoreInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	kvstoreService := KvstoreService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_kvstore_instance" {
			continue
		}

		_, err := kvstoreService.DescribeKVstoreInstance(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
	}

	return nil
}

func testAccKVStoreInstance_classic(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alicloud_zones" "default" {
		available_resource_creation = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_classic"
	}

	resource "alicloud_kvstore_instance" "default" {
		availability_zone = "${lookup(data.alicloud_zones.default.zones[(length(data.alicloud_zones.default.zones)-1)%%length(data.alicloud_zones.default.zones)], "id")}"
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
	}
	`, instanceType, instanceClass, engineVersion)
}

func testAccKVStoreInstance_classicUpdateParameter(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alicloud_zones" "default" {
		available_resource_creation = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_classic"
	}

	resource "alicloud_kvstore_instance" "default" {
		availability_zone = "${lookup(data.alicloud_zones.default.zones[(length(data.alicloud_zones.default.zones)-1)%%length(data.alicloud_zones.default.zones)], "id")}"
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
		parameters = [
			{
			  name = "maxmemory-policy"
			  value = "volatile-ttl"
			}
		]
	}
	`, instanceType, instanceClass, engineVersion)
}

func testAccKVStoreInstance_classicAddParameter(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alicloud_zones" "default" {
		available_resource_creation = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_classic"
	}

	resource "alicloud_kvstore_instance" "default" {
		availability_zone = "${lookup(data.alicloud_zones.default.zones[(length(data.alicloud_zones.default.zones)-1)%%length(data.alicloud_zones.default.zones)], "id")}"
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
		parameters = [
			{
				name = "maxmemory-policy"
				value = "volatile-ttl"
			  },
			  {
				  name = "slowlog-max-len"
				  value = "1111"
			  }
		]
	}
	`, instanceType, instanceClass, engineVersion)
}

func testAccKVStoreInstance_classicDeleteParameter(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alicloud_zones" "default" {
		available_resource_creation = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_classic"
	}

	resource "alicloud_kvstore_instance" "default" {
		availability_zone = "${lookup(data.alicloud_zones.default.zones[(length(data.alicloud_zones.default.zones)-1)%%length(data.alicloud_zones.default.zones)], "id")}"
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
		parameters = [
			{
				name = "slowlog-max-len"
				value = "1111"
			}
		]
	}
	`, instanceType, instanceClass, engineVersion)
}

func testAccKVStoreInstance_classicUpdateSecuirtyIps(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alicloud_zones" "default" {
		available_resource_creation = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_classic"
	}

	resource "alicloud_kvstore_instance" "default" {
		availability_zone = "${lookup(data.alicloud_zones.default.zones[(length(data.alicloud_zones.default.zones)-1)%%length(data.alicloud_zones.default.zones)], "id")}"
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.3", "10.0.0.2"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
	}
	`, instanceType, instanceClass, engineVersion)
}
func testAccKVStoreInstance_classicUpdateClass(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alicloud_zones" "default" {
		available_resource_creation = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_classic"
	}

	resource "alicloud_kvstore_instance" "default" {
		availability_zone = "${lookup(data.alicloud_zones.default.zones[(length(data.alicloud_zones.default.zones)-1)%%length(data.alicloud_zones.default.zones)], "id")}"
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
	}
	`, instanceType, instanceClass, engineVersion)
}
func testAccKVStoreInstance_classicUpdateAttr(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alicloud_zones" "default" {
		available_resource_creation = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_classic"
	}

	resource "alicloud_kvstore_instance" "default" {
		availability_zone = "${lookup(data.alicloud_zones.default.zones[(length(data.alicloud_zones.default.zones)-1)%%length(data.alicloud_zones.default.zones)], "id")}"
		password = "Abc123456"
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
	}
	`, instanceType, instanceClass, engineVersion)
}
func testAccKVStoreInstance_classicUpdateAll(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alicloud_zones" "default" {
		available_resource_creation = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_classicUpdateAll"
	}

	resource "alicloud_kvstore_instance" "default" {
		availability_zone = "${lookup(data.alicloud_zones.default.zones[(length(data.alicloud_zones.default.zones)-1)%%length(data.alicloud_zones.default.zones)], "id")}"
		password = "Abc123456"
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.2","10.0.0.3"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
	}
	`, instanceType, instanceClass, engineVersion)
}

func testAccKVStoreInstance_vpc(common, instanceClass, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_vpc"
	}
	resource "alicloud_kvstore_instance" "default" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = "${alicloud_vswitch.default.id}"
		private_ip     = "172.16.0.10"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		engine_version = "%s"
	}
	`, common, instanceClass, instanceType, engineVersion)
}
func testAccKVStoreInstance_vpcUpdateSecurityIps(common, instanceClass, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_vpc"
	}
	resource "alicloud_kvstore_instance" "default" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = "${alicloud_vswitch.default.id}"
		private_ip     = "172.16.0.10"
		security_ips = ["10.0.0.3", "10.0.0.2"]
		instance_type = "%s"
		engine_version = "%s"
	}
	`, common, instanceClass, instanceType, engineVersion)
}

func testAccKVStoreInstance_vpcUpdateVpcAuthMode(common, instanceClass, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_vpc"
	}
	resource "alicloud_kvstore_instance" "default" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = "${alicloud_vswitch.default.id}"
		vpc_auth_mode = "Close"
		private_ip     = "172.16.0.10"
		security_ips = ["10.0.0.3", "10.0.0.2"]
		instance_type = "%s"
		engine_version = "%s"
	}
	`, common, instanceClass, instanceType, engineVersion)
}

func testAccKVStoreInstance_vpcUpdateParameter(common, instanceClass, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_vpc"
	}
	resource "alicloud_kvstore_instance" "default" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = "${alicloud_vswitch.default.id}"
		private_ip     = "172.16.0.10"
		security_ips = ["10.0.0.3", "10.0.0.2"]
		parameters = [
			{
			  name = "maxmemory-policy"
			  value = "volatile-ttl"
			}
		]
		instance_type = "%s"
		engine_version = "%s"
	}
	`, common, instanceClass, instanceType, engineVersion)
}

func testAccKVStoreInstance_vpcAddParameter(common, instanceClass, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_vpc"
	}
	resource "alicloud_kvstore_instance" "default" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = "${alicloud_vswitch.default.id}"
		private_ip     = "172.16.0.10"
		security_ips = ["10.0.0.3", "10.0.0.2"]
		parameters = [
			{
			  name = "maxmemory-policy"
			  value = "volatile-ttl"
			},
			{
				name = "slowlog-max-len"
				value = "1111"
			}
		]
		instance_type = "%s"
		engine_version = "%s"
	}
	`, common, instanceClass, instanceType, engineVersion)
}

func testAccKVStoreInstance_vpcDeleteParameter(common, instanceClass, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_vpc"
	}
	resource "alicloud_kvstore_instance" "default" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = "${alicloud_vswitch.default.id}"
		private_ip     = "172.16.0.10"
		security_ips = ["10.0.0.3", "10.0.0.2"]
		parameters = [
			{
				name = "slowlog-max-len"
				value = "1111"
			}
		]
		instance_type = "%s"
		engine_version = "%s"
	}
	`, common, instanceClass, instanceType, engineVersion)
}

func testAccKVStoreInstance_vpcUpdateClass(common, instanceClass, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_vpc"
	}
	resource "alicloud_kvstore_instance" "default" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = "${alicloud_vswitch.default.id}"
		private_ip     = "172.16.0.10"
		security_ips = ["10.0.0.3", "10.0.0.2"]
		instance_type = "%s"
		engine_version = "%s"
	}
	`, common, instanceClass, instanceType, engineVersion)
}
func testAccKVStoreInstance_vpcUpdateAll(common, instanceClass, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_vpcUpdateAll"
	}
	resource "alicloud_kvstore_instance" "default" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		password       = "Abc12345"
		vswitch_id     = "${alicloud_vswitch.default.id}"
		private_ip     = "172.16.0.10"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		engine_version = "%s"
	}
	`, common, instanceClass, instanceType, engineVersion)
}

func testAccKVStoreInstance_vpcmulti(common, instanceClass, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_vpc"
	}
	resource "alicloud_kvstore_instance" "default" {
		count		   = 10
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = "${alicloud_vswitch.default.id}"
		private_ip     = "172.16.0.10"
		security_ips   = ["10.0.0.1"]
		instance_type  = "%s"
		engine_version = "%s"
	}
	`, common, instanceClass, instanceType, engineVersion)
}

func testAccKVStoreInstance_classicmulti(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alicloud_zones" "default" {
		available_resource_creation = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_classic"
	}

	resource "alicloud_kvstore_instance" "default" {
		count = 10
		availability_zone = "${lookup(data.alicloud_zones.default.zones[(length(data.alicloud_zones.default.zones)-1)%%length(data.alicloud_zones.default.zones)], "id")}"
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
	}
	`, instanceType, instanceClass, engineVersion)
}
