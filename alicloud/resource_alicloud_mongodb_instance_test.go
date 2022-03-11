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
	resource.AddTestSweepers("alicloud_mongodb_instance", &resource.Sweeper{
		Name: "alicloud_mongodb_instance",
		F:    testSweepMongoDBInstances,
	})
}

func testSweepMongoDBInstances(region string) error {
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
			return WrapErrorf(err, DefaultErrorMsg, "testSweepMongoDBInstances", request.GetActionName(), AlibabaCloudSdkGoERROR)
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
		request := dds.CreateDeleteDBInstanceRequest()
		request.DBInstanceId = id
		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.DeleteDBInstance(request)
		})

		if err != nil {
			log.Printf("[error] Failed to delete MongoDB instance,ID:%v(%v)\n", id, request.GetActionName())
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

func TestAccAlicloudMongoDBInstance_classic(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alicloud_mongodb_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBInstanceClassicConfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMongodbInstanceClassicConfig)
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
					"zone_id":             "${data.alicloud_mongodb_zones.default.zones.0.id}",
					"engine_version":      "3.4",
					"db_instance_storage": "10",
					"db_instance_class":   "dds.mongo.mid",
					"name":                name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version":       "3.4",
						"db_instance_storage":  "10",
						"db_instance_class":    "dds.mongo.mid",
						"name":                 name,
						"storage_engine":       "WiredTiger",
						"instance_charge_type": "PostPaid",
						"replication_factor":   "3",
						"replica_sets.#":       CHECKSET,
						"ssl_status":           CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ssl_action", "order_type", "auto_renew"},
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
					"db_instance_storage": "30",
					"db_instance_class":   "dds.mongo.standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage": "30",
						"db_instance_class":   "dds.mongo.standard",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "YourPassword_123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "YourPassword_123",
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
					"name":                name,
					"account_password":    "YourPassword_",
					"security_ip_list":    []string{"10.168.1.12", "10.168.1.13"},
					"db_instance_storage": "30",
					"db_instance_class":   "dds.mongo.standard",
					"backup_period":       []string{"Tuesday", "Wednesday"},
					"backup_time":         "10:00Z-11:00Z",
					"maintain_start_time": REMOVEKEY,
					"maintain_end_time":   REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                        name,
						"account_password":            "YourPassword_",
						"security_ip_list.#":          "2",
						"security_ip_list.4095458986": "10.168.1.12",
						"security_ip_list.3976237035": "10.168.1.13",
						"db_instance_storage":         "30",
						"db_instance_class":           "dds.mongo.standard",
						"backup_period.#":             "2",
						"backup_period.1592931319":    "Tuesday",
						"backup_period.1970423419":    "Wednesday",
						"backup_time":                 "10:00Z-11:00Z",
						"maintain_start_time":         REMOVEKEY,
						"maintain_end_time":           REMOVEKEY,
						"ssl_status":                  CHECKSET,
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

func resourceMongodbInstanceClassicConfig(name string) string {
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

func TestAccAlicloudMongoDBInstance_vpc(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alicloud_mongodb_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBInstanceVpcConfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMongodbInstanceVpcConfig)
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
					"vswitch_id":          "${data.alicloud_vswitches.default.ids.0}",
					"engine_version":      "4.0",
					"db_instance_storage": "10",
					"db_instance_class":   "dds.mongo.mid",
					"name":                name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version":       "4.0",
						"db_instance_storage":  "10",
						"db_instance_class":    "dds.mongo.mid",
						"name":                 name,
						"storage_engine":       "WiredTiger",
						"instance_charge_type": "PostPaid",
						"replication_factor":   "3",
						"replica_sets.#":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ssl_action", "order_type", "auto_renew"},
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
					"db_instance_storage": "30",
					"db_instance_class":   "dds.mongo.standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage": "30",
						"db_instance_class":   "dds.mongo.standard",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "YourPassword_123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "YourPassword_123",
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
					"name":                name,
					"account_password":    "YourPassword_",
					"security_ip_list":    []string{"10.168.1.12", "10.168.1.13"},
					"db_instance_storage": "30",
					"db_instance_class":   "dds.mongo.standard",
					"backup_period":       []string{"Tuesday", "Wednesday"},
					"backup_time":         "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                        name,
						"account_password":            "YourPassword_",
						"security_ip_list.#":          "2",
						"security_ip_list.4095458986": "10.168.1.12",
						"security_ip_list.3976237035": "10.168.1.13",
						"db_instance_storage":         "30",
						"db_instance_class":           "dds.mongo.standard",
						"backup_period.#":             "2",
						"backup_time":                 "10:00Z-11:00Z",
					}),
				),
			}},
	})
}

func resourceMongodbInstanceVpcConfig(name string) string {
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
	  zone_id = data.alicloud_mongodb_zones.default.zones.0.id
	}
`, name)
}

func TestAccAlicloudMongoDBInstance_multiAZ(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alicloud_mongodb_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBInstanceMultiAZConfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMongodbInstanceMultiAZConfig)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.MongoDBMultiAzSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":             "${data.alicloud_mongodb_zones.default.zones.0.id}",
					"vswitch_id":          "${data.alicloud_vswitches.default.ids.0}",
					"engine_version":      "3.4",
					"db_instance_storage": "10",
					"db_instance_class":   "dds.mongo.mid",
					"name":                name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version":       "3.4",
						"db_instance_storage":  "10",
						"db_instance_class":    "dds.mongo.mid",
						"name":                 name,
						"storage_engine":       "WiredTiger",
						"instance_charge_type": "PostPaid",
						"replication_factor":   "3",
						"replica_sets.#":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ssl_action", "order_type", "auto_renew"},
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
					"db_instance_storage": "30",
					"db_instance_class":   "dds.mongo.standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage": "30",
						"db_instance_class":   "dds.mongo.standard",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "YourPassword_123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "YourPassword_123",
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
					"name":                name,
					"account_password":    "YourPassword_",
					"security_ip_list":    []string{"10.168.1.12", "10.168.1.13"},
					"db_instance_storage": "30",
					"db_instance_class":   "dds.mongo.standard",
					"backup_period":       []string{"Tuesday", "Wednesday"},
					"backup_time":         "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                        name,
						"account_password":            "YourPassword_",
						"security_ip_list.#":          "2",
						"security_ip_list.4095458986": "10.168.1.12",
						"security_ip_list.3976237035": "10.168.1.13",
						"db_instance_storage":         "30",
						"db_instance_class":           "dds.mongo.standard",
						"backup_period.#":             "2",
						"backup_time":                 "10:00Z-11:00Z",
					}),
				),
			},
		},
	})
}
func resourceMongodbInstanceMultiAZConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_mongodb_zones" "default" {
		multi = true
	}
	data "alicloud_vpcs" "default" {
		name_regex = "default-NODELETING"
	}
	
	data "alicloud_vswitches" "default" {
	  vpc_id = data.alicloud_vpcs.default.ids.0
	  zone_id = data.alicloud_mongodb_zones.default.zones.0.multi_zone_ids.0
	}

	
`, name)
}
