package alicloud

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
	request.DBInstanceType = "sharding"
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
	service := VpcService{client}
	ddsService := MongoDBService{client}
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
		// If a mongoDB name is not set successfully, it should be fetched by vpc name and deleted.
		if skip {
			instance, err := ddsService.DescribeMongoDBInstance(id)
			if err != nil {
				if NotFoundError(err) {
					continue
				}
				log.Printf("[INFO] Describe MongoDB sharding instance: %s (%s) got an error: %#v\n", name, id, err)
			}
			if need, err := service.needSweepVpc(instance.VPCId, instance.VSwitchId); err == nil {
				skip = !need
			}
		}
		if skip {
			log.Printf("[INFO] Skipping MongoDB sharding instance: %s (%s)\n", name, id)
			continue
		}
		log.Printf("[INFO] Deleting MongoDB sharding instance: %s (%s)\n", name, id)

		request := dds.CreateDeleteDBInstanceRequest()
		request.DBInstanceId = id
		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.DeleteDBInstance(request)
		})

		if err != nil {
			log.Printf("[error] Failed to delete MongoDB sharding instance,ID:%v(%v)\n", id, request.GetActionName())
		} else {
			sweeped = true
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
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return err
	}
	return nil
}

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
		IDRefreshName: resourceId,
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
						"name":                      "",
						"storage_engine":            "WiredTiger",
						"instance_charge_type":      "PostPaid",
						"tags.%":                    "1",
						"tags.Created":              "TF",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
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
						"account_password": "YourPassword_",
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
				Config: testMongoDBShardingInstance_classic_backup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#":          "1",
						"backup_period.1970423419": "Wednesday",
						"backup_time":              "11:00Z-12:00Z",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_classic_together,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                        "tf-testAccMongoDBShardingInstance_test_together",
						"account_password":            "YourPassword_123",
						"security_ip_list.#":          "2",
						"security_ip_list.4095458986": "10.168.1.12",
						"security_ip_list.3976237035": "10.168.1.13",
						"backup_period.#":             "2",
						"backup_period.1592931319":    "Tuesday",
						"backup_period.1970423419":    "Wednesday",
						"backup_time":                 "10:00Z-11:00Z",
					}),
				),
			}},
	})
}

func TestAccAlicloudMongoDBShardingInstance_classicVersion4(t *testing.T) {
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckMongoDBShardingInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testMongoDBShardingInstance_classic_base4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":                   CHECKSET,
						"engine_version":            "4.0",
						"shard_list.#":              "2",
						"shard_list.0.node_class":   "dds.shard.mid",
						"shard_list.0.node_storage": "10",
						"shard_list.1.node_class":   "dds.shard.standard",
						"shard_list.1.node_storage": "20",
						"mongo_list.#":              "2",
						"mongo_list.0.node_class":   "dds.mongos.mid",
						"mongo_list.1.node_class":   "dds.mongos.mid",
						"name":                      "",
						"storage_engine":            "WiredTiger",
						"instance_charge_type":      "PostPaid",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testMongoDBShardingInstance_classic_tde,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tde_status": "enabled",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_classic_security_group_id,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id": CHECKSET,
					}),
				),
			}},
	})
}

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
			testAccPreCheckWithNoDefaultVpc(t)
		},
		IDRefreshName: resourceId,
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
						"name":                      "",
						"storage_engine":            "WiredTiger",
						"instance_charge_type":      "PostPaid",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
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
						"account_password": "YourPassword_",
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
				Config: testMongoDBShardingInstance_vpc_backup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#":          "1",
						"backup_period.1970423419": "Wednesday",
						"backup_time":              "11:00Z-12:00Z",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_vpc_together,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                        "tf-testAccMongoDBShardingInstance_test_together",
						"account_password":            "YourPassword_123",
						"security_ip_list.#":          "2",
						"security_ip_list.4095458986": "10.168.1.12",
						"security_ip_list.3976237035": "10.168.1.13",
						"backup_period.#":             "2",
						"backup_period.1592931319":    "Tuesday",
						"backup_period.1970423419":    "Wednesday",
						"backup_time":                 "10:00Z-11:00Z",
					}),
				),
			}},
	})
}

func TestAccAlicloudMongoDBShardingInstance_multi_instance(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alicloud_mongodb_sharding_instance.default.2"
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
		IDRefreshName: resourceId,
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
						"name":                      "",
						"storage_engine":            "WiredTiger",
						"instance_charge_type":      "PostPaid",
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
						"account_password": "YourPassword_",
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
				Config: testMongoDBShardingInstance_multi_instance_backup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#":          "1",
						"backup_period.1970423419": "Wednesday",
						"backup_time":              "11:00Z-12:00Z",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_multi_instance_together,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                        "tf-testAccMongoDBShardingInstance_test_together",
						"account_password":            "YourPassword_123",
						"security_ip_list.#":          "2",
						"security_ip_list.4095458986": "10.168.1.12",
						"security_ip_list.3976237035": "10.168.1.13",
						"backup_period.#":             "2",
						"backup_period.1592931319":    "Tuesday",
						"backup_period.1970423419":    "Wednesday",
						"backup_time":                 "10:00Z-11:00Z",
					}),
				),
			}},
	})
}

const testMongoDBShardingInstance_classic_base = `
data "alicloud_mongodb_zones" "default" {}
resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alicloud_mongodb_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  tags = {
    Created = "TF"
  }
}`

const testMongoDBShardingInstance_classic_base4 = `
data "alicloud_mongodb_zones" "default" {}
resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alicloud_mongodb_zones.default.zones.0.id}"
  engine_version = "4.0"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
}`

const testMongoDBShardingInstance_classic_tde = `
data "alicloud_mongodb_zones" "default" {}
resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alicloud_mongodb_zones.default.zones.0.id}"
  engine_version = "4.0"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  tde_status    = "enabled"
}`

const testMongoDBShardingInstance_classic_security_group_id = `
data "alicloud_mongodb_zones" "default" {}
data "alicloud_security_groups" "default" {
}
resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alicloud_mongodb_zones.default.zones.0.id}"
  engine_version = "4.0"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  tde_status    = "enabled"
  security_group_id    = "${data.alicloud_security_groups.default.groups.0.id}"
}`

const testMongoDBShardingInstance_classic_name = `
data "alicloud_mongodb_zones" "default" {}
resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alicloud_mongodb_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  name = "tf-testAccMongoDBShardingInstance_test"
}`

const testMongoDBShardingInstance_classic_account_password = `
data "alicloud_mongodb_zones" "default" {}
resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alicloud_mongodb_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "YourPassword_"
}`

const testMongoDBShardingInstance_classic_mongos = `
data "alicloud_mongodb_zones" "default" {}
resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alicloud_mongodb_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
   mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "YourPassword_"
}`

const testMongoDBShardingInstance_classic_shard = `
data "alicloud_mongodb_zones" "default" {}
resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alicloud_mongodb_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
   mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "YourPassword_"
}`

const testMongoDBShardingInstance_classic_backup = `
data "alicloud_mongodb_zones" "default" {}
resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alicloud_mongodb_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
   mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "YourPassword_"
  backup_period    = ["Wednesday"]
  backup_time      = "11:00Z-12:00Z"
}`

const testMongoDBShardingInstance_classic_together = `
data "alicloud_mongodb_zones" "default" {}
resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alicloud_mongodb_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
   mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test_together"
  account_password = "YourPassword_123"
  backup_period    = ["Tuesday", "Wednesday"]
  backup_time      = "10:00Z-11:00Z"
  security_ip_list = ["10.168.1.12", "10.168.1.13"]
}`

const testMongoDBShardingInstance_vpc_base = `
data "alicloud_mongodb_zones" "default" {}
data "alicloud_vpcs" "default" {
	is_default = true
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = "${data.alicloud_mongodb_zones.default.zones.0.id}"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_vpc"
}
resource "alicloud_mongodb_sharding_instance" "default" {
  vswitch_id     = data.alicloud_vswitches.default.ids.0
  zone_id        = "${data.alicloud_mongodb_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
}`

const testMongoDBShardingInstance_vpc_name = `
data "alicloud_mongodb_zones" "default" {}
data "alicloud_vpcs" "default" {
	is_default = true
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = "${data.alicloud_mongodb_zones.default.zones.0.id}"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_vpc"
}
resource "alicloud_mongodb_sharding_instance" "default" {
  vswitch_id     = data.alicloud_vswitches.default.ids.0
  zone_id        = "${data.alicloud_mongodb_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  name = "tf-testAccMongoDBShardingInstance_test"
}`

const testMongoDBShardingInstance_vpc_account_password = `
data "alicloud_mongodb_zones" "default" {}
data "alicloud_vpcs" "default" {
	is_default = true
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = "${data.alicloud_mongodb_zones.default.zones.0.id}"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_vpc"
}
resource "alicloud_mongodb_sharding_instance" "default" {
  vswitch_id     = data.alicloud_vswitches.default.ids.0
  zone_id        = "${data.alicloud_mongodb_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "YourPassword_"
}`

const testMongoDBShardingInstance_vpc_mongos = `
data "alicloud_mongodb_zones" "default" {}
data "alicloud_vpcs" "default" {
	is_default = true
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = "${data.alicloud_mongodb_zones.default.zones.0.id}"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_vpc"
}
resource "alicloud_mongodb_sharding_instance" "default" {
  vswitch_id     = data.alicloud_vswitches.default.ids.0
  zone_id        = "${data.alicloud_mongodb_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
   mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "YourPassword_"
}`

const testMongoDBShardingInstance_vpc_shard = `
data "alicloud_mongodb_zones" "default" {}
data "alicloud_vpcs" "default" {
	is_default = true
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = "${data.alicloud_mongodb_zones.default.zones.0.id}"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_vpc"
}
resource "alicloud_mongodb_sharding_instance" "default" {
  vswitch_id     = data.alicloud_vswitches.default.ids.0
  zone_id        = "${data.alicloud_mongodb_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
   mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "YourPassword_"
}`

const testMongoDBShardingInstance_vpc_backup = `
data "alicloud_mongodb_zones" "default" {}
data "alicloud_vpcs" "default" {
	is_default = true
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = "${data.alicloud_mongodb_zones.default.zones.0.id}"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_vpc"
}
resource "alicloud_mongodb_sharding_instance" "default" {
  vswitch_id     = data.alicloud_vswitches.default.ids.0
  zone_id        = "${data.alicloud_mongodb_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
   mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "YourPassword_"
  backup_period    = ["Wednesday"]
  backup_time      = "11:00Z-12:00Z"
}`

const testMongoDBShardingInstance_vpc_together = `
data "alicloud_mongodb_zones" "default" {}
data "alicloud_vpcs" "default" {
	is_default = true
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = "${data.alicloud_mongodb_zones.default.zones.0.id}"
}

variable "name" {
  default = "tf-testAccMongoDBShardingInstance_vpc"
}
resource "alicloud_mongodb_sharding_instance" "default" {
  vswitch_id     = data.alicloud_vswitches.default.ids.0
  zone_id        = "${data.alicloud_mongodb_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
   mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test_together"
  account_password = "YourPassword_123"
  backup_period    = ["Tuesday", "Wednesday"]
  backup_time      = "10:00Z-11:00Z"
  security_ip_list = ["10.168.1.12", "10.168.1.13"]
}`

const testMongoDBShardingInstance_multi_instance_base = `
data "alicloud_mongodb_zones" "default" {}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_multi_instance"
}

resource "alicloud_mongodb_sharding_instance" "default" {
  count          = 3
  zone_id        = "${data.alicloud_mongodb_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
}`

const testMongoDBShardingInstance_multi_instance_name = `
data "alicloud_mongodb_zones" "default" {}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_multi_instance"
}

resource "alicloud_mongodb_sharding_instance" "default" {
  count          = 3
  zone_id        = "${data.alicloud_mongodb_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  name = "tf-testAccMongoDBShardingInstance_test"
}`

const testMongoDBShardingInstance_multi_instance_account_password = `
data "alicloud_mongodb_zones" "default" {}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_multi_instance"
}

resource "alicloud_mongodb_sharding_instance" "default" {
  count          = 3
  zone_id        = "${data.alicloud_mongodb_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "YourPassword_"
}`

const testMongoDBShardingInstance_multi_instance_mongos = `
data "alicloud_mongodb_zones" "default" {}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_multi_instance"
}

resource "alicloud_mongodb_sharding_instance" "default" {
  count          = 3
  zone_id        = "${data.alicloud_mongodb_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "YourPassword_"
}`

const testMongoDBShardingInstance_multi_instance_shard = `
data "alicloud_mongodb_zones" "default" {}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_multi_instance"
}

resource "alicloud_mongodb_sharding_instance" "default" {
  count          = 3
  zone_id        = "${data.alicloud_mongodb_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "YourPassword_"
}`

const testMongoDBShardingInstance_multi_instance_backup = `
data "alicloud_mongodb_zones" "default" {}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_multi_instance"
}

resource "alicloud_mongodb_sharding_instance" "default" {
  count          = 3
  zone_id        = "${data.alicloud_mongodb_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
   mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "YourPassword_"
  backup_period    = ["Wednesday"]
  backup_time      = "11:00Z-12:00Z"
}`

const testMongoDBShardingInstance_multi_instance_together = `
data "alicloud_mongodb_zones" "default" {}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_multi_instance"
}

resource "alicloud_mongodb_sharding_instance" "default" {
  count          = 3
  zone_id        = "${data.alicloud_mongodb_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
   mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test_together"
  account_password = "YourPassword_123"
  backup_period    = ["Tuesday", "Wednesday"]
  backup_time      = "10:00Z-11:00Z"
  security_ip_list = ["10.168.1.12", "10.168.1.13"]
}`
