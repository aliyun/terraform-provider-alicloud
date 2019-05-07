package alicloud

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var insts []dds.DBInstance
	request := dds.CreateDescribeDBInstancesRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.DescribeDBInstances(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "testSweepMongoDBShardingInstances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		response, _ := raw.(*dds.DescribeDBInstancesResponse)
		addDebug(request.GetActionName(), response)

		if response == nil || len(response.DBInstances.DBInstance) < 1 {
			break
		}
		insts = append(insts, response.DBInstances.DBInstance...)

		if len(response.DBInstances.DBInstance) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	sweeped := false
	for _, v := range insts {
		name := v.DBInstanceDescription
		id := v.DBInstanceId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping MongoDB instance: %s (%s)\n", name, id)
			continue
		}
		log.Printf("[INFO] Deleting MongoDB instance: %s (%s)\n", name, id)

		sweeped = true

		request := dds.CreateDeleteDBInstanceRequest()
		request.DBInstanceId = id
		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.DeleteDBInstance(request)
		})

		if err != nil {
			log.Printf("[error] Failed to delete MongoDB instance,ID:%v(%v)\n", id, request.GetActionName())
		}
		addDebug(request.GetActionName(), raw)
	}
	if sweeped {
		// Waiting 30 seconds to eusure these DB instances have been deleted.
		time.Sleep(30 * time.Second)
	}
	return nil
}

func testAccCheckMongoDBShardingInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_mongodb_sharding_instance" {
			continue
		}
		_, err := ddsService.DescribeMongoDBInstance(rs.Primary.ID)
		if err != nil {
			if ddsService.NotFoundMongoDBInstance(err) {
				continue
			}
			return WrapError(err)
		}
		return err
	}
	return nil
}

const testMongoDBShardingInstance_classic_base = `
data "alicloud_zones" "default" {
  available_resource_creation = "MongoDB"
}
resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list = [{
    node_class   = "dds.shard.mid"
    node_storage = 10
    }, {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }]
  mongo_list = [{
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
  }]
}`

const testMongoDBShardingInstance_classic_name = `
data "alicloud_zones" "default" {
  available_resource_creation = "MongoDB"
}
resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list = [{
    node_class   = "dds.shard.mid"
    node_storage = 10
    }, {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }]
  mongo_list = [{
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
  }]
  name = "tf-testAccMongoDBShardingInstance_test"
}`

const testMongoDBShardingInstance_classic_account_password = `
data "alicloud_zones" "default" {
  available_resource_creation = "MongoDB"
}
resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list = [{
    node_class   = "dds.shard.mid"
    node_storage = 10
    }, {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }]
  mongo_list = [{
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
  }]
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "1234@test"
}`

const testMongoDBShardingInstance_classic_mongos = `
data "alicloud_zones" "default" {
  available_resource_creation = "MongoDB"
}
resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list = [{
    node_class   = "dds.shard.mid"
    node_storage = 10
    }, {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }]
  mongo_list = [{
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
  }]
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "1234@test"
}`

const testMongoDBShardingInstance_classic_shard = `
data "alicloud_zones" "default" {
  available_resource_creation = "MongoDB"
}
resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list = [{
    node_class   = "dds.shard.mid"
    node_storage = 10
    }, {
    node_class   = "dds.shard.standard"
    node_storage = 20
    }, {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }]
  mongo_list = [{
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
  }]
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "1234@test"
}`

const testMongoDBShardingInstance_classic_together = `
data "alicloud_zones" "default" {
  available_resource_creation = "MongoDB"
}
resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list = [{
    node_class   = "dds.shard.mid"
    node_storage = 10
    }, {
    node_class   = "dds.shard.standard"
    node_storage = 20
    }, {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }]
  mongo_list = [{
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
  }]
  name             = "tf-testAccMongoDBShardingInstance_test_together"
  account_password = "1234@test"
  security_ip_list = ["10.168.1.12", "10.168.1.13"]
}`

func TestAccAlicloudMongoDBShardingInstance_classic(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alicloud_mongodb_sharding_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.MongoDBClassicNoSupportedRegions)
		},
		IDRefreshName: "alicloud_mongodb_sharding_instance.default",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckMongoDBShardingInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testMongoDBShardingInstance_classic_base,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":                   CHECKSET,
						"engine_version":            "3.4",
						"shard_list.#":              "2",
						"shard_list.0.node_class":   "dds.shard.mid",
						"shard_list.0.node_storage": "10",
						"shard_list.1.node_class":   "dds.shard.standard",
						"shard_list.1.node_storage": "20",
						"mongo_list.#":              "2",
						"mongo_list.0.node_class":   "dds.mongos.mid",
						"mongo_list.1.node_class":   "dds.mongos.mid",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_classic_name,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccMongoDBShardingInstance_test",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_classic_account_password,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "1234@test",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_classic_mongos,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mongo_list.#":            "3",
						"mongo_list.0.node_class": "dds.mongos.mid",
						"mongo_list.1.node_class": "dds.mongos.mid",
						"mongo_list.2.node_class": "dds.mongos.mid",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_classic_shard,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"shard_list.#":              "3",
						"shard_list.0.node_class":   "dds.shard.mid",
						"shard_list.0.node_storage": "10",
						"shard_list.1.node_class":   "dds.shard.standard",
						"shard_list.1.node_storage": "20",
						"shard_list.2.node_class":   "dds.shard.standard",
						"shard_list.2.node_storage": "20",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_classic_together,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                        "tf-testAccMongoDBShardingInstance_test_together",
						"account_password":            "1234@test",
						"security_ip_list.#":          "2",
						"security_ip_list.4095458986": "10.168.1.12",
						"security_ip_list.3976237035": "10.168.1.13",
					}),
				),
			}},
	})
}

const testMongoDBShardingInstance_vpc_base = `
data "alicloud_zones" "default" {
  available_resource_creation = "MongoDB"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_vpc"
}
resource "alicloud_mongodb_sharding_instance" "default" {
  vswitch_id     = "${alicloud_vswitch.default.id}"
  zone_id        = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list = [{
    node_class   = "dds.shard.mid"
    node_storage = 10
    }, {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }]
  mongo_list = [{
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
  }]
}`

const testMongoDBShardingInstance_vpc_name = `
data "alicloud_zones" "default" {
  available_resource_creation = "MongoDB"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_vpc"
}
resource "alicloud_mongodb_sharding_instance" "default" {
  vswitch_id     = "${alicloud_vswitch.default.id}"
  zone_id        = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list = [{
    node_class   = "dds.shard.mid"
    node_storage = 10
    }, {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }]
  mongo_list = [{
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
  }]
  name = "tf-testAccMongoDBShardingInstance_test"
}`

const testMongoDBShardingInstance_vpc_account_password = `
data "alicloud_zones" "default" {
  available_resource_creation = "MongoDB"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_vpc"
}
resource "alicloud_mongodb_sharding_instance" "default" {
  vswitch_id     = "${alicloud_vswitch.default.id}"
  zone_id        = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list = [{
    node_class   = "dds.shard.mid"
    node_storage = 10
    }, {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }]
  mongo_list = [{
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
  }]
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "1234@test"
}`

const testMongoDBShardingInstance_vpc_mongos = `
data "alicloud_zones" "default" {
  available_resource_creation = "MongoDB"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_vpc"
}
resource "alicloud_mongodb_sharding_instance" "default" {
  vswitch_id     = "${alicloud_vswitch.default.id}"
  zone_id        = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list = [{
    node_class   = "dds.shard.mid"
    node_storage = 10
    }, {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }]
  mongo_list = [{
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
  }]
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "1234@test"
}`

const testMongoDBShardingInstance_vpc_shard = `
data "alicloud_zones" "default" {
  available_resource_creation = "MongoDB"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_vpc"
}
resource "alicloud_mongodb_sharding_instance" "default" {
  vswitch_id     = "${alicloud_vswitch.default.id}"
  zone_id        = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list = [{
    node_class   = "dds.shard.mid"
    node_storage = 10
    }, {
    node_class   = "dds.shard.standard"
    node_storage = 20
    }, {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }]
  mongo_list = [{
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
  }]
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "1234@test"
}`

const testMongoDBShardingInstance_vpc_together = `
data "alicloud_zones" "default" {
  available_resource_creation = "MongoDB"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_vpc"
}
resource "alicloud_mongodb_sharding_instance" "default" {
  vswitch_id     = "${alicloud_vswitch.default.id}"
  zone_id        = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list = [{
    node_class   = "dds.shard.mid"
    node_storage = 10
    }, {
    node_class   = "dds.shard.standard"
    node_storage = 20
    }, {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }]
  mongo_list = [{
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
  }]
  name             = "tf-testAccMongoDBShardingInstance_test_together"
  account_password = "1234@test"
  security_ip_list = ["10.168.1.12", "10.168.1.13"]
}`

func TestAccAlicloudMongoDBShardingInstance_vpc(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alicloud_mongodb_sharding_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_mongodb_sharding_instance.default",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckMongoDBShardingInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testMongoDBShardingInstance_vpc_base,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":                CHECKSET,
						"zone_id":                   CHECKSET,
						"engine_version":            "3.4",
						"shard_list.#":              "2",
						"shard_list.0.node_class":   "dds.shard.mid",
						"shard_list.0.node_storage": "10",
						"shard_list.1.node_class":   "dds.shard.standard",
						"shard_list.1.node_storage": "20",
						"mongo_list.#":              "2",
						"mongo_list.0.node_class":   "dds.mongos.mid",
						"mongo_list.1.node_class":   "dds.mongos.mid",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_vpc_name,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccMongoDBShardingInstance_test",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_vpc_account_password,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "1234@test",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_vpc_mongos,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mongo_list.#":            "3",
						"mongo_list.0.node_class": "dds.mongos.mid",
						"mongo_list.1.node_class": "dds.mongos.mid",
						"mongo_list.2.node_class": "dds.mongos.mid",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_vpc_shard,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"shard_list.#":              "3",
						"shard_list.0.node_class":   "dds.shard.mid",
						"shard_list.0.node_storage": "10",
						"shard_list.1.node_class":   "dds.shard.standard",
						"shard_list.1.node_storage": "20",
						"shard_list.2.node_class":   "dds.shard.standard",
						"shard_list.2.node_storage": "20",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_vpc_together,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                        "tf-testAccMongoDBShardingInstance_test_together",
						"account_password":            "1234@test",
						"security_ip_list.#":          "2",
						"security_ip_list.4095458986": "10.168.1.12",
						"security_ip_list.3976237035": "10.168.1.13",
					}),
				),
			}},
	})
}

const testMongoDBShardingInstance_multi_instance_base = `
data "alicloud_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_multi_instance"
}

resource "alicloud_mongodb_sharding_instance" "default" {
  count          = 5
  zone_id        = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list = [{
    node_class   = "dds.shard.mid"
    node_storage = 10
    }, {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }]
  mongo_list = [{
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
  }]
}`

const testMongoDBShardingInstance_multi_instance_name = `
data "alicloud_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_multi_instance"
}

resource "alicloud_mongodb_sharding_instance" "default" {
  count          = 5
  zone_id        = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list = [{
    node_class   = "dds.shard.mid"
    node_storage = 10
    }, {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }]
  mongo_list = [{
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
  }]
  name = "tf-testAccMongoDBShardingInstance_test"
}`

const testMongoDBShardingInstance_multi_instance_account_password = `
data "alicloud_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_multi_instance"
}

resource "alicloud_mongodb_sharding_instance" "default" {
  count          = 5
  zone_id        = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list = [{
    node_class   = "dds.shard.mid"
    node_storage = 10
    }, {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }]
  mongo_list = [{
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
  }]
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "1234@test"
}`

const testMongoDBShardingInstance_multi_instance_mongos = `
data "alicloud_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_multi_instance"
}

resource "alicloud_mongodb_sharding_instance" "default" {
  count          = 5
  zone_id        = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list = [{
    node_class   = "dds.shard.mid"
    node_storage = 10
    }, {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }]
  mongo_list = [{
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
  }]
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "1234@test"
}`

const testMongoDBShardingInstance_multi_instance_shard = `
data "alicloud_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_multi_instance"
}

resource "alicloud_mongodb_sharding_instance" "default" {
  count          = 5
  zone_id        = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list = [{
    node_class   = "dds.shard.mid"
    node_storage = 10
    }, {
    node_class   = "dds.shard.standard"
    node_storage = 20
    }, {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }]
  mongo_list = [{
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
  }]
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "1234@test"
}`

const testMongoDBShardingInstance_multi_instance_together = `
data "alicloud_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_multi_instance"
}

resource "alicloud_mongodb_sharding_instance" "default" {
  count          = 5
  zone_id        = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list = [{
    node_class   = "dds.shard.mid"
    node_storage = 10
    }, {
    node_class   = "dds.shard.standard"
    node_storage = 20
    }, {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }]
  mongo_list = [{
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
    }, {
    node_class = "dds.mongos.mid"
  }]
  name             = "tf-testAccMongoDBShardingInstance_test_together"
  account_password = "1234@test"
  security_ip_list = ["10.168.1.12", "10.168.1.13"]
}`

func TestAccAlicloudMongoDBShardingInstance_multi_instance(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alicloud_mongodb_sharding_instance.default.4"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

		},
		IDRefreshName: "alicloud_mongodb_sharding_instance.default.4",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckMongoDBShardingInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testMongoDBShardingInstance_multi_instance_base,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":                   CHECKSET,
						"engine_version":            "3.4",
						"shard_list.#":              "2",
						"shard_list.0.node_class":   "dds.shard.mid",
						"shard_list.0.node_storage": "10",
						"shard_list.1.node_class":   "dds.shard.standard",
						"shard_list.1.node_storage": "20",
						"mongo_list.#":              "2",
						"mongo_list.0.node_class":   "dds.mongos.mid",
						"mongo_list.1.node_class":   "dds.mongos.mid",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_multi_instance_name,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccMongoDBShardingInstance_test",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_multi_instance_account_password,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "1234@test",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_multi_instance_mongos,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mongo_list.#":            "3",
						"mongo_list.0.node_class": "dds.mongos.mid",
						"mongo_list.1.node_class": "dds.mongos.mid",
						"mongo_list.2.node_class": "dds.mongos.mid",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_multi_instance_shard,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"shard_list.#":              "3",
						"shard_list.0.node_class":   "dds.shard.mid",
						"shard_list.0.node_storage": "10",
						"shard_list.1.node_class":   "dds.shard.standard",
						"shard_list.1.node_storage": "20",
						"shard_list.2.node_class":   "dds.shard.standard",
						"shard_list.2.node_storage": "20",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_multi_instance_together,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                        "tf-testAccMongoDBShardingInstance_test_together",
						"account_password":            "1234@test",
						"security_ip_list.#":          "2",
						"security_ip_list.4095458986": "10.168.1.12",
						"security_ip_list.3976237035": "10.168.1.13",
					}),
				),
			}},
	})
}
