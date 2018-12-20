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

func TestAccAlicloudKVStoreInstance_vpc(t *testing.T) {
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
			resource.TestStep{
				Config: testAccKVStoreInstance_vpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists(
						"alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_kvstore_instance.foo",
						"instance_class",
						"redis.master.small.default"),
					resource.TestCheckResourceAttr(
						"alicloud_kvstore_instance.foo",
						"engine_version",
						"2.8"),
					resource.TestCheckResourceAttr(
						"alicloud_kvstore_instance.foo",
						"instance_type",
						"Redis"),
					resource.TestCheckResourceAttr(
						"alicloud_kvstore_instance.foo",
						"instance_name",
						"tf-testAccKVStoreInstance_vpc"),
					resource.TestCheckResourceAttr(
						"alicloud_kvstore_instance.foo",
						"private_ip",
						"172.16.0.10"),
				),
			},
		},
	})

}

func TestAccAlicloudKVStoreInstance_upgradeClass(t *testing.T) {
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
			resource.TestStep{
				Config: testAccKVStoreInstance_vpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists(
						"alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", "redis.master.small.default"),
				),
			},

			resource.TestStep{
				Config: testAccKVStoreInstance_classUpgrade,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists(
						"alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", "redis.master.mid.default"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_classUpgrade"),
				),
			},
		},
	})

}

func TestAccAlicloudKVStoreInstance_updateSecurityIps(t *testing.T) {
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
			resource.TestStep{
				Config: testAccKVStoreInstance_vpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists(
						"alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", "redis.master.small.default"),
				),
			},

			resource.TestStep{
				Config: testAccKVStoreInstance_updateSecurityIps,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists(
						"alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", "redis.master.mid.default"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_updateSecurityIps"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "security_ips.#", "2"),
				),
			},
		},
	})

}

func TestAccAlicloudKVStoreInstance_changeChargeType(t *testing.T) {
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceNopDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccKVStoreInstance_vpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists(
						"alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", "redis.master.small.default"),
				),
			},

			resource.TestStep{
				Config: testAccKVStoreInstance_changeChargeType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists(
						"alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_charge_type", "PrePaid"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "period", "1"),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_name", "tf-testAccKVStoreInstance_changeChargeType"),
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

		if _, err := kvstoreService.DescribeRKVInstanceById(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("Error KVStore Instance still exist")
	}

	return nil
}

func testAccCheckKVStoreInstanceNopDestroy(s *terraform.State) error {
	return nil
}

const testAccKVStoreInstance_vpc = `
data "alicloud_zones" "default" {
	available_resource_creation = "KVStore"
}
variable "name" {
	default = "tf-testAccKVStoreInstance_vpc"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "foo" {
 	vpc_id = "${alicloud_vpc.foo.id}"
 	cidr_block = "172.16.0.0/24"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
 	name = "${var.name}"
}

resource "alicloud_kvstore_instance" "foo" {
	instance_class = "redis.master.small.default"
	instance_name  = "${var.name}"
	password       = "Test12345"
	vswitch_id     = "${alicloud_vswitch.foo.id}"
	private_ip     = "172.16.0.10"
}
`
const testAccKVStoreInstance_classUpgrade = `
data "alicloud_zones" "default" {
	available_resource_creation = "KVStore"
}
variable "name" {
	default = "tf-testAccKVStoreInstance_classUpgrade"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "foo" {
 	vpc_id = "${alicloud_vpc.foo.id}"
 	cidr_block = "172.16.0.0/24"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
 	name = "${var.name}"
}

resource "alicloud_kvstore_instance" "foo" {
	instance_class = "redis.master.mid.default"
	instance_name  = "${var.name}"
	password       = "Test12345"
	vswitch_id     = "${alicloud_vswitch.foo.id}"
	private_ip     = "172.16.0.10"
}
`

const testAccKVStoreInstance_updateSecurityIps = `
data "alicloud_zones" "default" {
	available_resource_creation = "KVStore"
}
variable "name" {
	default = "tf-testAccKVStoreInstance_updateSecurityIps"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "foo" {
 	vpc_id = "${alicloud_vpc.foo.id}"
 	cidr_block = "172.16.0.0/24"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
 	name = "${var.name}"
}

resource "alicloud_kvstore_instance" "foo" {
	instance_class = "redis.master.mid.default"
	instance_name  = "${var.name}"
	password       = "Test12345"
	vswitch_id     = "${alicloud_vswitch.foo.id}"
	private_ip     = "172.16.0.10"
	security_ips   = ["10.110.10.10", "10.110.10.20"]
}
`
const testAccKVStoreInstance_changeChargeType = `
data "alicloud_zones" "default" {
	available_resource_creation = "KVStore"
}
variable "name" {
	default = "tf-testAccKVStoreInstance_changeChargeType"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "foo" {
 	vpc_id = "${alicloud_vpc.foo.id}"
 	cidr_block = "172.16.0.0/24"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
 	name = "${var.name}"
}

resource "alicloud_kvstore_instance" "foo" {
	instance_class = "redis.master.small.default"
	instance_name  = "${var.name}"
	password       = "Test12345"
	vswitch_id     = "${alicloud_vswitch.foo.id}"
	private_ip     = "172.16.0.10"
    instance_charge_type = "PrePaid"
    period = 1
}
`
