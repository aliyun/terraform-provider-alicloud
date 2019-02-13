package alicloud

import (
	"fmt"
	"log"
	"testing"

	"strings"
	"time"

	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

var redisInstanceConnectionDomainRegexp = regexp.MustCompile("^r-[a-z0-9]+.redis[.a-z-0-9]*.rds.aliyuncs.com")
var redisInstanceClassForTest = "redis.master.small.default"
var redisInstanceClassForTestUpdateClass = "redis.master.mid.default"
var memcacheInstanceConnectionDomainRegexp = regexp.MustCompile("^m-[a-z0-9]+.memcache[.a-z-0-9]*.rds.aliyuncs.com")
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
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_classic(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_classic"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", redisInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreRedis)),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "vswitch_id", ""),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", redisInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", ""),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},
		},
	})

}

func TestAccAlicloudKVStoreRedisInstance_classicUpdateSecurityIps(t *testing.T) {
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_classic(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_classic"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", redisInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreRedis)),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "vswitch_id", ""),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore4Dot0)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", redisInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", ""),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},
			{
				Config: testAccKVStoreInstance_classicUpdateSecuirtyIps(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_classic"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", redisInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreRedis)),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "vswitch_id", ""),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore4Dot0)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", redisInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", ""),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "2"),
				),
			},
		},
	})

}

func TestAccAlicloudKVStoreRedisInstance_classicUpdateClass(t *testing.T) {
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_classic(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_classic"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", redisInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreRedis)),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "vswitch_id", ""),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", redisInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", ""),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},
			{
				Config: testAccKVStoreInstance_classicUpdateClass(string(KVStoreRedis), redisInstanceClassForTestUpdateClass, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_classic"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", redisInstanceClassForTestUpdateClass),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreRedis)),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "vswitch_id", ""),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", redisInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", ""),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},
		},
	})

}

func TestAccAlicloudKVStoreRedisInstance_classicUpdateAttr(t *testing.T) {
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_classic(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_classic"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", redisInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreRedis)),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "vswitch_id", ""),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore4Dot0)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", redisInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", ""),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},
			{
				Config: testAccKVStoreInstance_classicUpdateAttr(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_classicUpdateAttr"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", redisInstanceClassForTest),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreRedis)),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "vswitch_id", ""),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore4Dot0)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", redisInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", ""),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},
		},
	})

}
func TestAccAlicloudKVStoreRedisInstance_classicUpdateAll(t *testing.T) {
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_classic(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_classic"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", redisInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreRedis)),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "vswitch_id", ""),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", redisInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", ""),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},
			{
				Config: testAccKVStoreInstance_classicUpdateAll(string(KVStoreRedis), redisInstanceClassForTestUpdateClass, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_classicUpdateAll"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", redisInstanceClassForTestUpdateClass),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreRedis)),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "vswitch_id", ""),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", redisInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", ""),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "2"),
				),
			},
		},
	})

}
func TestAccAlicloudKVStoreRedisInstance_vpc(t *testing.T) {
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_vpc(KVStoreCommonTestCase, redisInstanceClassForTest, string(KVStoreRedis), string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_vpc"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", redisInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreRedis)),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore4Dot0)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", redisInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", "172.16.0.10"),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},
		},
	})

}
func TestAccAlicloudKVStoreRedisInstance_vpcUpdateSecurityIps(t *testing.T) {
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_vpc(KVStoreCommonTestCase, redisInstanceClassForTest, string(KVStoreRedis), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_vpc"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", redisInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreRedis)),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", redisInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", "172.16.0.10"),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},

			{
				Config: testAccKVStoreInstance_vpcUpdateSecurityIps(KVStoreCommonTestCase, redisInstanceClassForTest, string(KVStoreRedis), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_vpc"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", redisInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreRedis)),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", redisInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", "172.16.0.10"),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "2"),
				),
			},
		},
	})

}

func TestAccAlicloudKVStoreRedisInstance_vpcUpdateVpcAuthMode(t *testing.T) {
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_vpc(KVStoreCommonTestCase, redisInstanceClassForTest, string(KVStoreRedis), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_vpc"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", redisInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreRedis)),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", redisInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", "172.16.0.10"),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "vpc_auth_mode", "Open"),
				),
			},

			{
				Config: testAccKVStoreInstance_vpcUpdateVpcAuthMode(KVStoreCommonTestCase, redisInstanceClassForTest, string(KVStoreRedis), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_vpc"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", redisInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreRedis)),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", redisInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", "172.16.0.10"),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "2"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "vpc_auth_mode", "Close"),
				),
			},
		},
	})

}

func TestAccAlicloudKVStoreRedisInstance_vpcUpdateClass(t *testing.T) {
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_vpc(KVStoreCommonTestCase, redisInstanceClassForTest, string(KVStoreRedis), string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_vpc"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", redisInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreRedis)),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore4Dot0)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", redisInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", "172.16.0.10"),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},

			{
				Config: testAccKVStoreInstance_vpcUpdateClass(KVStoreCommonTestCase, redisInstanceClassForTestUpdateClass, string(KVStoreRedis), string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_vpc"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", redisInstanceClassForTestUpdateClass),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreRedis)),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore4Dot0)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", redisInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", "172.16.0.10"),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},
		},
	})

}
func TestAccAlicloudKVStoreRedisInstance_vpcUpdateAttr(t *testing.T) {
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_vpc(KVStoreCommonTestCase, redisInstanceClassForTest, string(KVStoreRedis), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_vpc"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", redisInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreRedis)),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", redisInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", "172.16.0.10"),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},

			{
				Config: testAccKVStoreInstance_vpcUpdateAttr(KVStoreCommonTestCase, redisInstanceClassForTest, string(KVStoreRedis), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_vpcUpdateAttr"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", redisInstanceClassForTest),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreRedis)),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", redisInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", "172.16.0.10"),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},
		},
	})

}

func TestAccAlicloudKVStoreRedisInstance_vpcUpdateAll(t *testing.T) {
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_vpc(KVStoreCommonTestCase, redisInstanceClassForTest, string(KVStoreRedis), string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_vpc"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", redisInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreRedis)),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore4Dot0)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", redisInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", "172.16.0.10"),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},

			{
				Config: testAccKVStoreInstance_vpcUpdateAll(KVStoreCommonTestCase, redisInstanceClassForTestUpdateClass, string(KVStoreRedis), string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_vpcUpdateAll"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", redisInstanceClassForTestUpdateClass),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreRedis)),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore4Dot0)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", redisInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", "172.16.0.10"),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "2"),
				),
			},
		},
	})

}

// Test Memcache
func TestAccAlicloudKVStoreMemcacheInstance_classic(t *testing.T) {
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_classic(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_classic"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", memcacheInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreMemcache)),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "vswitch_id", ""),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", memcacheInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", ""),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},
		},
	})

}

func TestAccAlicloudKVStoreMemcacheInstance_classicUpdateSecurityIps(t *testing.T) {
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_classic(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_classic"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", memcacheInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreMemcache)),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "vswitch_id", ""),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", memcacheInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", ""),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},
			{
				Config: testAccKVStoreInstance_classicUpdateSecuirtyIps(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_classic"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", memcacheInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreMemcache)),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "vswitch_id", ""),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", memcacheInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", ""),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "2"),
				),
			},
		},
	})

}

func TestAccAlicloudKVStoreMemcacheInstance_classicUpdateClass(t *testing.T) {
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_classic(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_classic"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", memcacheInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreMemcache)),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "vswitch_id", ""),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", memcacheInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", ""),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},
			{
				Config: testAccKVStoreInstance_classicUpdateClass(string(KVStoreMemcache), memcacheInstanceClassForTestUpdateClass, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_classic"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", memcacheInstanceClassForTestUpdateClass),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreMemcache)),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "vswitch_id", ""),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", memcacheInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", ""),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},
		},
	})

}

func TestAccAlicloudKVStoreMemcacheInstance_classicUpdateAttr(t *testing.T) {
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_classic(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_classic"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", memcacheInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreMemcache)),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "vswitch_id", ""),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", memcacheInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", ""),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},
			{
				Config: testAccKVStoreInstance_classicUpdateAttr(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_classicUpdateAttr"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", memcacheInstanceClassForTest),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreMemcache)),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "vswitch_id", ""),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", memcacheInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", ""),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},
		},
	})

}
func TestAccAlicloudKVStoreMemcacheInstance_classicUpdateAll(t *testing.T) {
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_classic(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_classic"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", memcacheInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreMemcache)),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "vswitch_id", ""),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", memcacheInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", ""),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},
			{
				Config: testAccKVStoreInstance_classicUpdateAll(string(KVStoreMemcache), memcacheInstanceClassForTestUpdateClass, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_classicUpdateAll"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", memcacheInstanceClassForTestUpdateClass),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreMemcache)),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "vswitch_id", ""),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", memcacheInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", ""),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "2"),
				),
			},
		},
	})

}
func TestAccAlicloudKVStoreMemcacheInstance_vpc(t *testing.T) {
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_vpc(KVStoreCommonTestCase, memcacheInstanceClassForTest, string(KVStoreMemcache), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_vpc"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", memcacheInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreMemcache)),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", memcacheInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", "172.16.0.10"),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},
		},
	})

}
func TestAccAlicloudKVStoreMemcacheInstance_vpcUpdateSecurityIps(t *testing.T) {
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_vpc(KVStoreCommonTestCase, memcacheInstanceClassForTest, string(KVStoreMemcache), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_vpc"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", memcacheInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreMemcache)),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", memcacheInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", "172.16.0.10"),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},

			{
				Config: testAccKVStoreInstance_vpcUpdateSecurityIps(KVStoreCommonTestCase, memcacheInstanceClassForTest, string(KVStoreMemcache), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_vpc"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", memcacheInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreMemcache)),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", memcacheInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", "172.16.0.10"),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "2"),
				),
			},
		},
	})

}
func TestAccAlicloudKVStoreMemcacheInstance_vpcUpdateClass(t *testing.T) {
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_vpc(KVStoreCommonTestCase, memcacheInstanceClassForTest, string(KVStoreMemcache), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_vpc"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", memcacheInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreMemcache)),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", memcacheInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", "172.16.0.10"),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},

			{
				Config: testAccKVStoreInstance_vpcUpdateClass(KVStoreCommonTestCase, memcacheInstanceClassForTestUpdateClass, string(KVStoreMemcache), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_vpc"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", memcacheInstanceClassForTestUpdateClass),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreMemcache)),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", memcacheInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", "172.16.0.10"),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},
		},
	})

}
func TestAccAlicloudKVStoreMemcacheInstance_vpcUpdateAttr(t *testing.T) {
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_vpc(KVStoreCommonTestCase, memcacheInstanceClassForTest, string(KVStoreMemcache), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_vpc"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", memcacheInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreMemcache)),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", memcacheInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", "172.16.0.10"),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},

			{
				Config: testAccKVStoreInstance_vpcUpdateAttr(KVStoreCommonTestCase, memcacheInstanceClassForTest, string(KVStoreMemcache), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_vpcUpdateAttr"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", memcacheInstanceClassForTest),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreMemcache)),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", memcacheInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", "172.16.0.10"),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},
		},
	})

}

func TestAccAlicloudKVStoreMemcacheInstance_vpcUpdateAll(t *testing.T) {
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_vpc(KVStoreCommonTestCase, memcacheInstanceClassForTest, string(KVStoreMemcache), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_vpc"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", memcacheInstanceClassForTest),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreMemcache)),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", memcacheInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", "172.16.0.10"),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "1"),
				),
			},

			{
				Config: testAccKVStoreInstance_vpcUpdateAll(KVStoreCommonTestCase, memcacheInstanceClassForTestUpdateClass, string(KVStoreMemcache), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists("alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_vpcUpdateAll"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", memcacheInstanceClassForTestUpdateClass),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "password"),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "period"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_type", string(KVStoreMemcache)),
					resource.TestCheckResourceAttrSet("alicloud_kvstore_instance.foo", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "engine_version", string(KVStore2Dot8)),
					resource.TestMatchResourceAttr("alicloud_kvstore_instance.foo", "connection_domain", memcacheInstanceConnectionDomainRegexp),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "private_ip", "172.16.0.10"),
					resource.TestCheckNoResourceAttr("alicloud_kvstore_instance.foo", "backup_id"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "2"),
				),
			},
		},
	})

}

func testAccCheckKVStoreInstanceExists(n string, d *r_kvstore.DBInstanceAttribute) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No KVStore Instance ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		kvstoreService := KvstoreService{client}
		attr, err := kvstoreService.DescribeRKVInstanceById(rs.Primary.ID)
		log.Printf("[DEBUG] check instance %s attribute %#v", rs.Primary.ID, attr)

		if err != nil {
			return err
		}

		*d = *attr
		return nil
	}
}

func testAccCheckKVStoreInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	kvstoreService := KvstoreService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_kvstore_instance" {
			continue
		}

		_, err := kvstoreService.DescribeRKVInstanceById(rs.Primary.ID)
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

	resource "alicloud_kvstore_instance" "foo" {
		availability_zone = "${lookup(data.alicloud_zones.default.zones[(length(data.alicloud_zones.default.zones)-1)%%length(data.alicloud_zones.default.zones)], "id")}"
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
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

	resource "alicloud_kvstore_instance" "foo" {
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

	resource "alicloud_kvstore_instance" "foo" {
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
		default = "tf-testAccKVStoreInstance_classicUpdateAttr"
	}

	resource "alicloud_kvstore_instance" "foo" {
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

	resource "alicloud_kvstore_instance" "foo" {
		availability_zone = "${lookup(data.alicloud_zones.default.zones[(length(data.alicloud_zones.default.zones)-1)%%length(data.alicloud_zones.default.zones)], "id")}"
		password = "Abc123456"
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.3", "10.0.0.2"]
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
	resource "alicloud_kvstore_instance" "foo" {
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
	resource "alicloud_kvstore_instance" "foo" {
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
	resource "alicloud_kvstore_instance" "foo" {
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

func testAccKVStoreInstance_vpcUpdateClass(common, instanceClass, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_vpc"
	}
	resource "alicloud_kvstore_instance" "foo" {
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
func testAccKVStoreInstance_vpcUpdateAttr(common, instanceClass, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_vpcUpdateAttr"
	}
	resource "alicloud_kvstore_instance" "foo" {
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

func testAccKVStoreInstance_vpcUpdateAll(common, instanceClass, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_vpcUpdateAll"
	}
	resource "alicloud_kvstore_instance" "foo" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		password       = "Abc12345"
		vswitch_id     = "${alicloud_vswitch.default.id}"
		private_ip     = "172.16.0.10"
		security_ips = ["10.0.0.3", "10.0.0.2"]
		instance_type = "%s"
		engine_version = "%s"
	}
	`, common, instanceClass, instanceType, engineVersion)
}
