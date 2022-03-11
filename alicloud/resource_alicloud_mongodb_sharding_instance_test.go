package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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

func resourceMongodbShardingInstanceClassicConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_mongodb_zones" "default" {}

	resource "alicloud_security_group" "default" {
		name = var.name
	}
`, name)
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
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBShardingInstanceClassicConfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMongodbShardingInstanceClassicConfig)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.MongoDBClassicNoSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":        "${data.alicloud_mongodb_zones.default.zones.0.id}",
					"engine_version": "3.4",
					"name":           name,
					"shard_list": []map[string]interface{}{
						{
							"node_class":   "dds.shard.mid",
							"node_storage": "10",
						},
						{
							"node_class":        "dds.shard.standard",
							"node_storage":      "20",
							"readonly_replicas": "1",
						},
					},
					"mongo_list": []map[string]interface{}{
						{
							"node_class": "dds.mongos.mid",
						},
						{
							"node_class": "dds.mongos.mid",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":                        CHECKSET,
						"engine_version":                 "3.4",
						"shard_list.#":                   "2",
						"shard_list.0.node_class":        "dds.shard.mid",
						"shard_list.0.node_storage":      "10",
						"shard_list.0.readonly_replicas": "0",
						"shard_list.1.node_class":        "dds.shard.standard",
						"shard_list.1.node_storage":      "20",
						"shard_list.1.readonly_replicas": "1",
						"mongo_list.#":                   "2",
						"mongo_list.0.node_class":        "dds.mongos.mid",
						"mongo_list.1.node_class":        "dds.mongos.mid",
						"name":                           name,
						"storage_engine":                 "WiredTiger",
						"instance_charge_type":           "PostPaid",
						"tags.%":                         "0",
						"config_server_list.#":           CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"order_type", "auto_renew"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "YourPassword_",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "YourPassword_",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mongo_list": []map[string]interface{}{
						{
							"node_class": "dds.mongos.mid",
						},
						{
							"node_class": "dds.mongos.mid",
						},
						{
							"node_class": "dds.mongos.mid",
						},
					},
				}),
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
				Config: testAccConfig(map[string]interface{}{
					"shard_list": []map[string]interface{}{
						{
							"node_class":   "dds.shard.mid",
							"node_storage": "10",
						},
						{
							"node_class":        "dds.shard.standard",
							"node_storage":      "20",
							"readonly_replicas": "1",
						},
						{
							"node_class":        "dds.shard.standard",
							"node_storage":      "20",
							"readonly_replicas": "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"shard_list.#":                   "3",
						"shard_list.0.node_class":        "dds.shard.mid",
						"shard_list.0.node_storage":      "10",
						"shard_list.0.readonly_replicas": "0",
						"shard_list.1.node_class":        "dds.shard.standard",
						"shard_list.1.node_storage":      "20",
						"shard_list.1.readonly_replicas": "1",
						"shard_list.2.node_class":        "dds.shard.standard",
						"shard_list.2.node_storage":      "20",
						"shard_list.2.readonly_replicas": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_list": []string{"10.168.1.12"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_period": []string{"Wednesday"},
					"backup_time":   "11:00Z-12:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "1",
						"backup_time":     "11:00Z-12:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":             name,
					"account_password": "YourPassword_123",
					"security_ip_list": []string{"10.168.1.12", "10.168.1.13"},
					"backup_period":    []string{"Tuesday", "Wednesday"},
					"backup_time":      "10:00Z-11:00Z",
					"shard_list": []map[string]interface{}{
						{
							"node_class":   "dds.shard.mid",
							"node_storage": "10",
						},
						{
							"node_class":        "dds.shard.standard",
							"node_storage":      "20",
							"readonly_replicas": "1",
						},
					},
					"mongo_list": []map[string]interface{}{
						{
							"node_class": "dds.mongos.mid",
						},
						{
							"node_class": "dds.mongos.mid",
						},
					},
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                           name,
						"account_password":               "YourPassword_123",
						"security_ip_list.#":             "2",
						"backup_period.#":                "2",
						"backup_time":                    "10:00Z-11:00Z",
						"shard_list.#":                   "2",
						"shard_list.0.node_class":        "dds.shard.mid",
						"shard_list.0.node_storage":      "10",
						"shard_list.1.node_class":        "dds.shard.standard",
						"shard_list.1.node_storage":      "20",
						"shard_list.1.readonly_replicas": "1",
						"shard_list.2.node_class":        REMOVEKEY,
						"shard_list.2.readonly_replicas": REMOVEKEY,
						"shard_list.2.node_storage":      REMOVEKEY,
						"mongo_list.#":                   "2",
						"mongo_list.0.node_class":        "dds.mongos.mid",
						"mongo_list.1.node_class":        "dds.mongos.mid",
						"mongo_list.2.node_class":        REMOVEKEY,
						"tags.%":                         "0",
						"tags.For":                       REMOVEKEY,
						"tags.Created":                   REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type": "PrePaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "PrePaid",
					}),
				),
			},
		},
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
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBShardingInstanceVpcConfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMongodbShardingInstanceVpcConfig)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.MongoDBClassicNoSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":        "${data.alicloud_mongodb_zones.default.zones.0.id}",
					"vswitch_id":     "${data.alicloud_vswitches.default.ids.0}",
					"engine_version": "4.0",
					"name":           name,
					"shard_list": []map[string]interface{}{
						{
							"node_class":   "dds.shard.mid",
							"node_storage": "10",
						},
						{
							"node_class":        "dds.shard.standard",
							"node_storage":      "20",
							"readonly_replicas": "1",
						},
					},
					"mongo_list": []map[string]interface{}{
						{
							"node_class": "dds.mongos.mid",
						},
						{
							"node_class": "dds.mongos.mid",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":                        CHECKSET,
						"vswitch_id":                     CHECKSET,
						"engine_version":                 "4.0",
						"shard_list.#":                   "2",
						"shard_list.0.node_class":        "dds.shard.mid",
						"shard_list.0.node_storage":      "10",
						"shard_list.0.readonly_replicas": "0",
						"shard_list.1.node_class":        "dds.shard.standard",
						"shard_list.1.node_storage":      "20",
						"shard_list.1.readonly_replicas": "1",
						"mongo_list.#":                   "2",
						"mongo_list.0.node_class":        "dds.mongos.mid",
						"mongo_list.1.node_class":        "dds.mongos.mid",
						"name":                           name,
						"storage_engine":                 "WiredTiger",
						"instance_charge_type":           "PostPaid",
						"tags.%":                         "0",
						"config_server_list.#":           CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"order_type", "auto_renew"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tde_status": "enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tde_status": "enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "YourPassword_",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "YourPassword_",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mongo_list": []map[string]interface{}{
						{
							"node_class": "dds.mongos.mid",
						},
						{
							"node_class": "dds.mongos.mid",
						},
						{
							"node_class": "dds.mongos.mid",
						},
					},
				}),
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
				Config: testAccConfig(map[string]interface{}{
					"shard_list": []map[string]interface{}{
						{
							"node_class":   "dds.shard.mid",
							"node_storage": "10",
						},
						{
							"node_class":        "dds.shard.standard",
							"node_storage":      "20",
							"readonly_replicas": "1",
						},
						{
							"node_class":        "dds.shard.standard",
							"node_storage":      "20",
							"readonly_replicas": "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"shard_list.#":                   "3",
						"shard_list.0.node_class":        "dds.shard.mid",
						"shard_list.0.node_storage":      "10",
						"shard_list.0.readonly_replicas": "0",
						"shard_list.1.node_class":        "dds.shard.standard",
						"shard_list.1.node_storage":      "20",
						"shard_list.1.readonly_replicas": "1",
						"shard_list.2.node_class":        "dds.shard.standard",
						"shard_list.2.node_storage":      "20",
						"shard_list.2.readonly_replicas": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_list": []string{"10.168.1.12"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_period": []string{"Wednesday"},
					"backup_time":   "11:00Z-12:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "1",
						"backup_time":     "11:00Z-12:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":             name,
					"account_password": "YourPassword_123",
					"security_ip_list": []string{"10.168.1.12", "10.168.1.13"},
					"backup_period":    []string{"Tuesday", "Wednesday"},
					"backup_time":      "10:00Z-11:00Z",
					"shard_list": []map[string]interface{}{
						{
							"node_class":   "dds.shard.mid",
							"node_storage": "10",
						},
						{
							"node_class":        "dds.shard.standard",
							"node_storage":      "20",
							"readonly_replicas": "1",
						},
					},
					"mongo_list": []map[string]interface{}{
						{
							"node_class": "dds.mongos.mid",
						},
						{
							"node_class": "dds.mongos.mid",
						},
					},
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                           name,
						"account_password":               "YourPassword_123",
						"security_ip_list.#":             "2",
						"backup_period.#":                "2",
						"backup_time":                    "10:00Z-11:00Z",
						"shard_list.#":                   "2",
						"shard_list.0.node_class":        "dds.shard.mid",
						"shard_list.0.node_storage":      "10",
						"shard_list.1.node_class":        "dds.shard.standard",
						"shard_list.1.node_storage":      "20",
						"shard_list.1.readonly_replicas": "1",
						"shard_list.2.node_class":        REMOVEKEY,
						"shard_list.2.readonly_replicas": REMOVEKEY,
						"shard_list.2.node_storage":      REMOVEKEY,
						"mongo_list.#":                   "2",
						"mongo_list.0.node_class":        "dds.mongos.mid",
						"mongo_list.1.node_class":        "dds.mongos.mid",
						"mongo_list.2.node_class":        REMOVEKEY,
						"tags.%":                         "0",
						"tags.For":                       REMOVEKEY,
						"tags.Created":                   REMOVEKEY,
					}),
				),
			},
		},
	})
}
func resourceMongodbShardingInstanceVpcConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_mongodb_zones" "default" {}
	data "alicloud_vpcs" "default" {
		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
	  vpc_id = data.alicloud_vpcs.default.ids.0
	  zone_id = "${data.alicloud_mongodb_zones.default.zones.0.id}"
	}
`, name)
}
