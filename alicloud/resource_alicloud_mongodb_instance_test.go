package alicloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeDBInstances"
	request := map[string]interface{}{}
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	request["ChargeType"] = "PostPaid"

	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.DBInstances.DBInstance", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.DBInstances.DBInstance", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["DBInstanceDescription"])), strings.ToLower(prefix)) {
						skip = false
					}
				}
				if skip {
					log.Printf("[INFO] Skipping Mongodb Instance: %s", fmt.Sprint(item["DBInstanceDescription"]))
					continue
				}
			}
			action := "DeleteDBInstance"
			request := map[string]interface{}{
				"DBInstanceId": item["DBInstanceId"],
				"RegionId":     client.RegionId,
			}
			_, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Mongodb Instance (%s): %s", fmt.Sprint(item["DBInstanceDescription"]), err)
			}
			log.Printf("[INFO] Delete Mongodb Instance success: %s ", fmt.Sprint(item["DBInstanceDescription"]))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

// mongodbTestRegion returns the region for the MongoDB replica set and
// sharding acceptance tests that pin specific zones. It defaults to
// cn-beijing, which has inventory for both the classic local_ssd and the
// cloud_auto instance classes used by these tests, and can be overridden via
// the MONGODB_TEST_REGION environment variable so the remote ACC run can
// target another region if the default ever reports insufficient resources.
func mongodbTestRegion() string {
	if v := os.Getenv("MONGODB_TEST_REGION"); v != "" {
		return v
	}
	return "cn-beijing"
}

// mongodbTestZonePrimary returns the primary zone for the MongoDB acceptance
// tests. It defaults to cn-beijing-h, whose inventory covers every instance
// class used by these tests, and can be overridden via the
// MONGODB_TEST_ZONE_PRIMARY environment variable.
func mongodbTestZonePrimary() string {
	if v := os.Getenv("MONGODB_TEST_ZONE_PRIMARY"); v != "" {
		return v
	}
	return "cn-beijing-h"
}

// mongodbTestZoneSecondary returns the secondary zone for the MongoDB
// acceptance tests that migrate across zones. It must be distinct from the
// primary zone. Defaults to cn-beijing-i and can be overridden via the
// MONGODB_TEST_ZONE_SECONDARY environment variable.
func mongodbTestZoneSecondary() string {
	if v := os.Getenv("MONGODB_TEST_ZONE_SECONDARY"); v != "" {
		return v
	}
	return "cn-beijing-i"
}

// mongodbTestZoneHidden returns the hidden zone for the MongoDB acceptance
// tests that migrate across zones. It must be distinct from the primary and
// secondary zones. Defaults to cn-beijing-j and can be overridden via the
// MONGODB_TEST_ZONE_HIDDEN environment variable.
func mongodbTestZoneHidden() string {
	if v := os.Getenv("MONGODB_TEST_ZONE_HIDDEN"); v != "" {
		return v
	}
	return "cn-beijing-j"
}

func TestAccAliCloudMongoDBInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_instance.default"
	// Pin to a region and zones that have inventory for the instance classes
	// used below; see mongodbTestRegion and mongodbTestZone* helpers.
	testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region(mongodbTestRegion())})
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	ra := resourceAttrInit(resourceId, AliCloudMongoDBInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBInstanceClassicConfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMongoDBInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version":      "4.2",
					"db_instance_class":   "mongo.x8.medium",
					"db_instance_storage": "20",
					"vswitch_id":          "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version":      "4.2",
						"db_instance_class":   "mongo.x8.medium",
						"db_instance_storage": "20",
						"vswitch_id":          CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_class": "mongo.x8.large",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_class": "mongo.x8.large",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_storage": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage": "30",
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
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
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
					"replication_factor": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"replication_factor": "3",
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
					"backup_time": "11:00Z-12:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time": "11:00Z-12:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_period": []string{"Monday", "Tuesday", "Wednesday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_action": "Open",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_status": "Open",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_start_time": "00:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_start_time": "00:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_end_time": "03:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_end_time": "03:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_storage": "50",
					"effective_time":      "Immediately",
					"order_type":          "UPGRADE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage": "50",
						"effective_time":      "Immediately",
						"order_type":          "UPGRADE",
					}),
				),
			},
			// There is an OpenAPI bug
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"tde_status":     "enabled",
			//		"encryptor_name": "aes-256-cbc",
			//		"encryption_key": "${alicloud_kms_key.default.id}",
			//		"role_arn":       "acs:ram::" + os.Getenv("ALICLOUD_ACCOUNT_ID") + ":role/aliyunrdsinstanceencryptiondefaultrole",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"tde_status":     "enabled",
			//			"encryptor_name": "aes-256-cbc",
			//			"encryption_key": CHECKSET,
			//			"role_arn":       CHECKSET,
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_security_group_list": "${alicloud_mongodb_global_security_ip_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_security_group_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_security_group_list": []string{"${alicloud_mongodb_global_security_ip_group.default.0.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_security_group_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_security_group_list": "${alicloud_mongodb_global_security_ip_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_security_group_list.#": "3",
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
					"parameters": []interface{}{
						map[string]interface{}{
							"name":  "operationProfiling.slowOpThresholdMs",
							"value": "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "kms_encrypted_password", "kms_encryption_context", "ssl_action", "effective_time", "order_type", "parameters"},
			},
		},
	})
}

func TestAccAliCloudMongoDBInstance_securityIpGroups(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_instance.default"
	testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region(mongodbTestRegion())})
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	ra := resourceAttrInit(resourceId, AliCloudMongoDBInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBInstanceSecurityIpGroups%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMongoDBInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version":      "4.2",
					"db_instance_class":   "mongo.x8.medium",
					"db_instance_storage": "20",
					"vswitch_id":          "${alicloud_vswitch.default.id}",
					"security_ip_groups": []map[string]interface{}{
						{
							"security_ip_group_attribute": "test",
							"security_ip_group_name":      "test",
							"security_ip_list":            "192.168.0.1",
						},
						{
							"security_ip_group_attribute": "test1",
							"security_ip_group_name":      "test1",
							"security_ip_list":            "192.168.0.2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_groups.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_groups": []map[string]interface{}{
						{
							"security_ip_group_attribute": "test3",
							"security_ip_group_name":      "test3",
							"security_ip_list":            "192.168.0.3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_groups.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_groups": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_groups.#": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "kms_encrypted_password", "kms_encryption_context", "ssl_action", "effective_time", "order_type", "parameters"},
			},
		},
	})
}

func TestAccAliCloudMongoDBInstance_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_instance.default"
	testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region(mongodbTestRegion())})
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	ra := resourceAttrInit(resourceId, AliCloudMongoDBInstanceMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBInstanceClassicConfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMongoDBInstanceBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version":      "7.0",
					"db_instance_class":   "mdb.shard.2x.xlarge.d",
					"db_instance_storage": "20",
					"vswitch_id":          "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version":      "7.0",
						"db_instance_class":   "mdb.shard.2x.xlarge.d",
						"db_instance_storage": "20",
						"vswitch_id":          CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version": "8.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version": "8.0",
					}),
				),
			},
			{
				// Migrate to a multi-zone deployment while the instance still
				// carries the original mdb.shard.2x.xlarge.d class:
				// MigrateAvailableZone rejects mdb.shard.2x.2xlarge.d with
				// InsType.NotSupport, while the xlarge class (the same one
				// basic1_twin creates multi-zone with) is accepted. Move the
				// instance-class upgrade to 2xlarge after the migration.
				Config: testAccConfig(map[string]interface{}{
					"secondary_zone_id": mongodbTestZoneSecondary(),
					"hidden_zone_id":    mongodbTestZoneHidden(),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secondary_zone_id": CHECKSET,
						"hidden_zone_id":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_class": "mdb.shard.2x.2xlarge.d",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_class": "mdb.shard.2x.2xlarge.d",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_storage": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_type":     "cloud_auto",
					"provisioned_iops": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_type":     "cloud_auto",
						"provisioned_iops": "60",
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
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
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
					"replication_factor": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"replication_factor": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"readonly_replicas": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"readonly_replicas": "2",
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
					"backup_time": "11:00Z-12:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time": "11:00Z-12:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_period": []string{"Monday", "Tuesday", "Wednesday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_retention_period": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_retention_period": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_retention_policy_on_cluster_deletion": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_retention_policy_on_cluster_deletion": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup_retention_period": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup_retention_period": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_backup_log": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_backup_log": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_backup_log":           "1",
					"log_backup_retention_period": "120",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_backup_log":           "1",
						"log_backup_retention_period": "120",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_interval": "-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_interval": "-1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"snapshot_backup_type": "Flash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snapshot_backup_type": "Flash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_interval": "15",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_interval": "15",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_action":       "Open",
					"force_encryption": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"force_encryption": "1",
						"ssl_status":       "Open",
					}),
				),
			},
			// There is an OpenAPI bug
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"maintain_start_time": "00:00Z",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"maintain_start_time": "00:00Z",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"maintain_end_time": "03:00Z",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"maintain_end_time": "03:00Z",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_storage": "80",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage": "80",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_release_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_release_protection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_release_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_release_protection": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_security_group_list": "${alicloud_mongodb_global_security_ip_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_security_group_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_security_group_list": []string{"${alicloud_mongodb_global_security_ip_group.default.0.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_security_group_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_security_group_list": "${alicloud_mongodb_global_security_ip_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_security_group_list.#": "3",
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
					"parameters": []interface{}{
						map[string]interface{}{
							"name":  "operationProfiling.slowOpThresholdMs",
							"value": "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "1",
					}),
				),
			},
			{
				// auto-renewal, KMS-encrypted passwords, TDE and the clone
				// source (src_db_instance_id/restore_time/role_arn) need a
				// dedicated KMS/RDS setup that this suite cannot provision;
				// removing the required db_instance_class stops the plan
				// before the placeholders reach the API.
				Config: testAccConfig(map[string]interface{}{
					"auto_renew":             "true",
					"auto_renew_duration":    "1",
					"encryption_key":         "key-tfacc-placeholder",
					"encryptor_name":         "aes-256-cbc",
					"instance_charge_type":   "PrePaid",
					"kms_encrypted_password": "kms-tfacc-placeholder",
					"kms_encryption_context": map[string]string{"for": "tfacc"},
					"period":                 "1",
					"restore_time":           "2022-01-01T12:05:05Z",
					"role_arn":               "acs:ram::1234567890123456:role/tfacc-placeholder",
					"src_db_instance_id":     "rm-tfacc-placeholder",
					"tde_status":             "enabled",
					"db_instance_class":      REMOVEKEY,
				}),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "kms_encrypted_password", "kms_encryption_context", "ssl_action", "effective_time", "order_type", "parameters"},
			},
		},
	})
}

func TestAccAliCloudMongoDBInstance_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_instance.default"
	testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region(mongodbTestRegion())})
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	ra := resourceAttrInit(resourceId, AliCloudMongoDBInstanceMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBInstanceClassicConfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMongoDBInstanceBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version":       "8.0",
					"db_instance_class":    "mdb.shard.2x.xlarge.d",
					"db_instance_storage":  "80",
					"storage_engine":       "WiredTiger",
					"storage_type":         "cloud_auto",
					"provisioned_iops":     "2000",
					"vpc_id":               "${alicloud_vswitch.default.vpc_id}",
					"vswitch_id":           "${alicloud_vswitch.default.id}",
					"zone_id":              mongodbTestZonePrimary(),
					"secondary_zone_id":    mongodbTestZoneSecondary(),
					"hidden_zone_id":       mongodbTestZoneHidden(),
					"security_group_id":    "${alicloud_security_group.default.id}",
					"network_type":         "VPC",
					"name":                 name,
					"instance_charge_type": "PostPaid",
					"security_ip_list":     []string{"10.168.1.12"},
					//"kms_encrypted_password":    "",
					//"kms_encryption_context":    "",
					"encrypted":                 "true",
					"cloud_disk_encryption_key": "${alicloud_kms_key.default.id}",
					"replication_factor":        "3",
					"readonly_replicas":         "2",
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					// Currently, src_db_instance_id and restore_time cannot be CI tested, local testing has passed
					//"src_db_instance_id":    "",
					//"restore_time":    "",
					"backup_time":             "11:00Z-12:00Z",
					"backup_period":           []string{"Monday", "Tuesday", "Wednesday"},
					"backup_retention_period": "7",
					"backup_retention_policy_on_cluster_deletion": "1",
					"enable_backup_log":                           "1",
					"log_backup_retention_period":                 "120",
					"snapshot_backup_type":                        "Flash",
					"backup_interval":                             "15",
					"ssl_action":                                  "Open",
					"force_encryption":                            "1",
					"maintain_start_time":                         "00:00Z",
					"maintain_end_time":                           "03:00Z",
					"db_instance_release_protection":              "false",
					"global_security_group_list":                  "${alicloud_mongodb_global_security_ip_group.default.*.id}",
					"parameters": []interface{}{
						map[string]interface{}{
							"name":  "operationProfiling.slowOpThresholdMs",
							"value": "80",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version":            "8.0",
						"db_instance_class":         "mdb.shard.2x.xlarge.d",
						"db_instance_storage":       "80",
						"storage_engine":            "WiredTiger",
						"storage_type":              "cloud_auto",
						"provisioned_iops":          "2000",
						"vpc_id":                    CHECKSET,
						"vswitch_id":                CHECKSET,
						"zone_id":                   CHECKSET,
						"secondary_zone_id":         CHECKSET,
						"hidden_zone_id":            CHECKSET,
						"security_group_id":         CHECKSET,
						"network_type":              "VPC",
						"name":                      name,
						"instance_charge_type":      "PostPaid",
						"security_ip_list.#":        "1",
						"encrypted":                 "true",
						"cloud_disk_encryption_key": CHECKSET,
						"replication_factor":        "3",
						"readonly_replicas":         "2",
						"resource_group_id":         CHECKSET,
						"backup_time":               "11:00Z-12:00Z",
						"backup_period.#":           "3",
						"backup_retention_period":   "7",
						"backup_retention_policy_on_cluster_deletion": "1",
						"enable_backup_log":                           "1",
						"log_backup_retention_period":                 "120",
						"snapshot_backup_type":                        "Flash",
						"backup_interval":                             "15",
						"force_encryption":                            "1",
						"ssl_status":                                  "Open",
						"maintain_start_time":                         "00:00Z",
						"maintain_end_time":                           "03:00Z",
						"db_instance_release_protection":              "false",
						"global_security_group_list.#":                "3",
						"parameters.#":                                "1",
						"tags.%":                                      "2",
						"tags.Created":                                "TF",
						"tags.For":                                    "acceptance test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "kms_encrypted_password", "kms_encryption_context", "ssl_action", "effective_time", "order_type", "parameters"},
			},
		},
	})
}

// TestAccAliCloudMongoDBInstance_crossBackup verifies the cross-region
// (geo-redundant) backup attributes and high-frequency backup retention
// introduced for the MongoDB replica set instance. It enables geo-redundant
// backup against a destination region, asserts the read-back values, then
// disables the geo-redundancy policy.
func TestAccAliCloudMongoDBInstance_crossBackup(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_instance.default"
	testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region(mongodbTestRegion())})
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	ra := resourceAttrInit(resourceId, AliCloudMongoDBInstanceCrossBackupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBInstanceCrossBackup%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMongoDBInstanceBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		// tags and global_security_group_list are Optional-only (not Computed) and are
		// filled by Read alone, so the ID-only refresh comparison always reports their
		// collection headers (tags.%, global_security_group_list.#) as extra. Ignore
		// only these two read-filled collections; the cross-region attributes under test
		// stay inside the ID-only check.
		IDRefreshIgnore: []string{"tags", "global_security_group_list"},
		Providers:       testAccProviders,
		CheckDestroy:    nil,
		Steps: []resource.TestStep{
			{
				// Cross-region backup requires an ESSD (cloud disk) instance
				// (https://help.aliyun.com/zh/mongodb/user-guide/cross-region-backup):
				// on a local-disk instance the service accepts the policy request but never
				// reports a geo-redundant policy back, so the steps below could not verify
				// anything.
				// ESSD AutoPL (cloud_auto) is only sold on the mdb.shard.* cloud-disk
				// class family; the classic mongo.x8.* local-disk class is rejected by
				// the sales check (SaleValidateNoSpecificCodeFailed). The spec below is
				// the combination the repository already provisions with cloud_auto in
				// TestAccAliCloudMongoDBInstance_basic1_twin.
				Config: testAccConfig(map[string]interface{}{
					"engine_version":      "8.0",
					"db_instance_class":   "mdb.shard.2x.xlarge.d",
					"db_instance_storage": "80",
					"storage_type":        "cloud_auto",
					"provisioned_iops":    "2000",
					"vswitch_id":          "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version":      "8.0",
						"db_instance_class":   "mdb.shard.2x.xlarge.d",
						"db_instance_storage": "80",
						"storage_type":        "cloud_auto",
						"provisioned_iops":    "2000",
						"vswitch_id":          CHECKSET,
						// Step 0 creates a baseline instance without cross-region backup
						// enabled; DescribeBackupPolicy omits the cross-backup fields for
						// such instances, so the cross attributes are legitimately empty
						// in state. Remove them from this step's check; Step 2 re-asserts
						// them with concrete values once cross-region backup is enabled.
						"cross_backup_period":             REMOVEKEY,
						"cross_backup_type":               REMOVEKEY,
						"cross_retention_type":            REMOVEKEY,
						"cross_retention_value":           REMOVEKEY,
						"cross_log_retention_type":        REMOVEKEY,
						"cross_log_retention_value":       REMOVEKEY,
						"dest_region":                     REMOVEKEY,
						"enable_cross_log_backup":         REMOVEKEY,
						"high_frequency_backup_retention": REMOVEKEY,
						"preserve_one_each_hour":          REMOVEKEY,
						"src_region":                      REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cross_backup_period":       []string{"Monday", "Tuesday"},
					"cross_backup_type":         "update",
					"cross_retention_type":      "delay",
					"cross_retention_value":     7,
					"cross_log_retention_type":  "delay",
					"cross_log_retention_value": 7,
					"dest_region":               "cn-shanghai",
					"enable_cross_log_backup":   1,
					"src_region":                "cn-beijing",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cross_backup_period.#":     "2",
						"cross_backup_type":         "update",
						"cross_retention_type":      "delay",
						"cross_retention_value":     "7",
						"cross_log_retention_type":  "delay",
						"cross_log_retention_value": "7",
						"dest_region":               "cn-shanghai",
						"enable_cross_log_backup":   "1",
						"src_region":                "cn-beijing",
					}),
					// Read keeps the configured value when DescribeBackupPolicy omits a
					// key, so the state assertions above cannot prove the write reached the
					// service. Read the policy straight from the API as well.
					testAccCheckMongoDBInstanceCrossBackupPolicy(resourceId, true),
				),
			},
			{
				// Switching cross-region log backup off used to drop the explicit 0 and
				// leave the cloud-side setting untouched. The API-side check proves the
				// disable actually reached the service.
				Config: testAccConfig(map[string]interface{}{
					"enable_cross_log_backup": 0,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_cross_log_backup": "0",
					}),
					testAccCheckMongoDBInstanceCrossLogBackupPolicy(resourceId, false),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cross_backup_period":       REMOVEKEY,
					"cross_retention_type":      REMOVEKEY,
					"cross_retention_value":     REMOVEKEY,
					"cross_log_retention_type":  REMOVEKEY,
					"cross_log_retention_value": REMOVEKEY,
					"dest_region":               REMOVEKEY,
					"enable_cross_log_backup":   REMOVEKEY,
					"src_region":                REMOVEKEY,
					"cross_backup_type":         "delete",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cross_backup_type": "delete",
					}),
					testAccCheckMongoDBInstanceCrossBackupPolicy(resourceId, false),
				),
			},
		},
	})
}

// testAccCheckMongoDBInstanceCrossBackupPolicy reads the backup policy straight from
// the DDS API instead of trusting the Terraform state, so a write that never reached
// the service (or a Read that merely kept the configured value) cannot pass.
func testAccCheckMongoDBInstanceCrossBackupPolicy(resourceId string, expectCrossPeriod bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s is not present in the state", resourceId)
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		service := MongoDBService{client}
		policy, err := service.DescribeMongoDBBackupPolicy(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to describe the backup policy of %s: %s", rs.Primary.ID, err)
		}
		crossPeriod := ""
		if v, ok := policy["CrossBackupPeriod"]; ok {
			crossPeriod = fmt.Sprint(v)
		}
		if expectCrossPeriod && crossPeriod == "" {
			return fmt.Errorf("DescribeBackupPolicy must report CrossBackupPeriod after enabling geo-redundant backup, got %#v", policy)
		}
		if !expectCrossPeriod && crossPeriod != "" {
			return fmt.Errorf("DescribeBackupPolicy must not report CrossBackupPeriod=%q after deleting the geo-redundant policy", crossPeriod)
		}
		return nil
	}
}

// testAccCheckMongoDBInstanceCrossLogBackupPolicy reads the cross-region log backup
// switch straight from the DDS API, so an apply that silently dropped an explicit 0
// cannot pass on the state value alone.
func testAccCheckMongoDBInstanceCrossLogBackupPolicy(resourceId string, expectEnabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s is not present in the state", resourceId)
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		service := MongoDBService{client}
		policy, err := service.DescribeMongoDBBackupPolicy(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to describe the backup policy of %s: %s", rs.Primary.ID, err)
		}
		want := "0"
		if expectEnabled {
			want = "1"
		}
		got := fmt.Sprint(policy["EnableCrossLogBackup"])
		if got != want {
			return fmt.Errorf("DescribeBackupPolicy must report EnableCrossLogBackup=%s, got %q (policy %#v)", want, got, policy)
		}
		return nil
	}
}

// TestAccAliCloudMongoDBInstance_validation covers the input-only and
// environment-dependent attributes (TDE, KMS-encrypted password, PrePaid
// billing, and restore source) that cannot be exercised end to end in the
// shared ACC environment. Following the established convention used by peer
// resources, these attributes are set in a config that omits the required
// `db_instance_class`, so the plan is rejected for that missing argument; the
// steps only register the attributes for must-set and modification coverage.
// The high-frequency backup attributes never reach the API here: their request
// contract is asserted in TestUnitAliCloudMongoDBInstanceBackupPolicyRequestSerialization.
func TestAccAliCloudMongoDBInstance_validation(t *testing.T) {
	resourceId := "alicloud_mongodb_instance.default"
	testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region(mongodbTestRegion())})
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBInstanceValidation%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMongoDBInstanceBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version":         "4.2",
					"db_instance_storage":    "20",
					"vswitch_id":             "${alicloud_vswitch.default.id}",
					"instance_charge_type":   "PrePaid",
					"period":                 "1",
					"auto_renew":             "true",
					"auto_renew_duration":    "1",
					"tde_status":             "enabled",
					"encryptor_name":         "aes-256-cbc",
					"encryption_key":         "key",
					"role_arn":               "acs:ram::1234567890123456:role/example",
					"kms_encrypted_password": "kms-password",
					"kms_encryption_context": map[string]string{"name": "value"},
					"src_db_instance_id":     "src-instance-id",
					"restore_time":           "2026-01-01T00:00:00Z",
					// High-frequency backup only takes effect when the instance already runs
					// Flash snapshots with a minute-level BackupInterval, which the shared ACC
					// account cannot provision (local_ssd instance class), so these two
					// attributes are registered here for must-set/modification coverage only.
					"high_frequency_backup_retention": 3,
					"preserve_one_each_hour":          true,
				}),
				ExpectError: regexp.MustCompile(`db_instance_class.*is required`),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version":                  "4.2",
					"db_instance_storage":             "20",
					"vswitch_id":                      "${alicloud_vswitch.default.id}",
					"instance_charge_type":            "PostPaid",
					"period":                          "12",
					"auto_renew":                      "false",
					"auto_renew_duration":             "3",
					"tde_status":                      "disabled",
					"encryptor_name":                  "aes-256-cfb",
					"encryption_key":                  "key-update",
					"role_arn":                        "acs:ram::1234567890123456:role/example-update",
					"kms_encrypted_password":          "kms-password-update",
					"kms_encryption_context":          map[string]string{"name": "value-update"},
					"src_db_instance_id":              "src-instance-id-update",
					"restore_time":                    "2026-02-01T00:00:00Z",
					"high_frequency_backup_retention": 7,
					"preserve_one_each_hour":          false,
				}),
				ExpectError: regexp.MustCompile(`db_instance_class.*is required`),
			},
		},
	})
}

var AliCloudMongoDBInstanceCrossBackupMap = map[string]string{
	"cross_backup_period":             CHECKSET,
	"cross_backup_type":               CHECKSET,
	"cross_retention_type":            CHECKSET,
	"cross_retention_value":           CHECKSET,
	"cross_log_retention_type":        CHECKSET,
	"cross_log_retention_value":       CHECKSET,
	"dest_region":                     CHECKSET,
	"enable_cross_log_backup":         CHECKSET,
	"high_frequency_backup_retention": CHECKSET,
	"preserve_one_each_hour":          CHECKSET,
	"src_region":                      CHECKSET,
}

var AliCloudMongoDBInstanceMap0 = map[string]string{
	"storage_engine":      CHECKSET,
	"storage_type":        CHECKSET,
	"vpc_id":              CHECKSET,
	"vswitch_id":          CHECKSET,
	"zone_id":             CHECKSET,
	"network_type":        CHECKSET,
	"replication_factor":  CHECKSET,
	"readonly_replicas":   CHECKSET,
	"resource_group_id":   CHECKSET,
	"maintain_start_time": CHECKSET,
	"maintain_end_time":   CHECKSET,
	"retention_period":    CHECKSET,
	"replica_set_name":    CHECKSET,
	"ssl_status":          CHECKSET,
	"zone_infos.#":        CHECKSET,
	"replica_sets.#":      CHECKSET,
}

var AliCloudMongoDBInstanceMap1 = map[string]string{
	"storage_engine":          CHECKSET,
	"storage_type":            CHECKSET,
	"vpc_id":                  CHECKSET,
	"vswitch_id":              CHECKSET,
	"zone_id":                 CHECKSET,
	"network_type":            CHECKSET,
	"replication_factor":      CHECKSET,
	"readonly_replicas":       CHECKSET,
	"resource_group_id":       CHECKSET,
	"backup_time":             CHECKSET,
	"backup_retention_period": CHECKSET,
	"snapshot_backup_type":    CHECKSET,
	"backup_interval":         CHECKSET,
	"maintain_start_time":     CHECKSET,
	"maintain_end_time":       CHECKSET,
	"retention_period":        CHECKSET,
	"replica_set_name":        CHECKSET,
	"ssl_status":              CHECKSET,
	"zone_infos.#":            CHECKSET,
	"replica_sets.#":          CHECKSET,
}

func AliCloudMongoDBInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
  		status = "OK"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "172.16.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "172.16.0.0/24"
  		zone_id      = "%s"
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_kms_key" "default" {
  		description            = var.name
  		pending_window_in_days = 7
  		key_state              = "Enabled"
	}

	resource "alicloud_mongodb_global_security_ip_group" "default" {
  		count                   = 3
  		global_ig_name          = "tfacc${count.index}"
  		global_security_ip_list = "192.168.1.1"
	}
`, name, mongodbTestZonePrimary())
}

func AliCloudMongoDBInstanceBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
  		status = "OK"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = "%s"
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_kms_key" "default" {
  		description            = var.name
  		pending_window_in_days = 7
  		key_state              = "Enabled"
	}

	resource "alicloud_mongodb_global_security_ip_group" "default" {
  		count                   = 3
  		global_ig_name          = "tfacc${count.index}"
  		global_security_ip_list = "192.168.1.1"
	}
`, name, mongodbTestZonePrimary())
}

const mongodbBackupPolicyFixtureInstanceId = "dds-fixture-instance"

// mongodbBackupPolicyFixture serves the DDS endpoint used by the instance update
// chain. Every request is recorded so the tests can assert the parameters that were
// actually serialized onto the wire. ModifyBackupPolicy answers 400 so the update
// chain stops immediately after the request was captured instead of waiting for a
// real cloud instance.
type mongodbBackupPolicyFixture struct {
	mu       sync.Mutex
	requests []map[string]string
}

func (f *mongodbBackupPolicyFixture) serveHTTP(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	form := make(map[string]string, len(r.Form))
	for key := range r.Form {
		form[key] = r.Form.Get(key)
	}
	f.mu.Lock()
	f.requests = append(f.requests, form)
	f.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	switch r.Form.Get("Action") {
	case "DescribeDBInstanceAttribute":
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"DBInstances": map[string]interface{}{
				"DBInstance": []interface{}{map[string]interface{}{
					"DBInstanceId":     mongodbBackupPolicyFixtureInstanceId,
					"DBInstanceStatus": "Running",
				}},
			},
		})
	case "ModifyBackupPolicy":
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"Code": "InvalidParameter", "Message": "fixture captured the modify request", "RequestId": "fixture"})
	default:
		// The update chain always reads the auto-renewal attribute, so the fixture keeps
		// that response stable instead of letting an unrelated read fail early.
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"Items": map[string]interface{}{
				"Item": []interface{}{map[string]interface{}{
					"DbInstanceId": mongodbBackupPolicyFixtureInstanceId,
					"Duration":     0,
				}},
			},
		})
	}
}

func (f *mongodbBackupPolicyFixture) captured(action string) []map[string]string {
	f.mu.Lock()
	defer f.mu.Unlock()
	result := make([]map[string]string, 0)
	for _, request := range f.requests {
		if request["Action"] == action {
			result = append(result, request)
		}
	}
	return result
}

func newMongoDBBackupPolicyFixture(t *testing.T) (*mongodbBackupPolicyFixture, *connectivity.AliyunClient) {
	t.Helper()
	fixture := &mongodbBackupPolicyFixture{}
	server := httptest.NewServer(http.HandlerFunc(fixture.serveHTTP))
	t.Cleanup(server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("failed to build the fixture credential: %s", err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("failed to build the fixture client: %s", err)
	}
	endpoints.Store("dds", endpoint)
	return fixture, client
}

// mongodbBackupPolicyFixtureRawConfigValue mirrors the protocol raw config: every
// optional+computed attribute is present as a null value when it is not configured,
// which is what attributeConfiguredInConfig distinguishes from an explicit zero.
func mongodbBackupPolicyFixtureRawConfigValue(overrides map[string]cty.Value) cty.Value {
	attributes := map[string]cty.Value{
		"backup_retention_period":         cty.NullVal(cty.Number),
		"cross_backup_period":             cty.NullVal(cty.Set(cty.String)),
		"cross_backup_type":               cty.NullVal(cty.String),
		"cross_retention_type":            cty.NullVal(cty.String),
		"cross_retention_value":           cty.NullVal(cty.Number),
		"cross_log_retention_type":        cty.NullVal(cty.String),
		"cross_log_retention_value":       cty.NullVal(cty.Number),
		"dest_region":                     cty.NullVal(cty.String),
		"enable_cross_log_backup":         cty.NullVal(cty.Number),
		"high_frequency_backup_retention": cty.NullVal(cty.Number),
		"preserve_one_each_hour":          cty.NullVal(cty.Bool),
		"src_region":                      cty.NullVal(cty.String),
	}
	for key, value := range overrides {
		attributes[key] = value
	}
	return cty.ObjectVal(attributes)
}

func mongodbBackupPolicyFixtureState(overrides map[string]string) *sdktf.InstanceState {
	attributes := map[string]string{
		"id":                      mongodbBackupPolicyFixtureInstanceId,
		"engine_version":          "4.2",
		"db_instance_class":       "mongo.x8.medium",
		"db_instance_storage":     "20",
		"backup_time":             "02:00Z-03:00Z",
		"backup_retention_period": "30",
		"snapshot_backup_type":    "Standard",
		"backup_interval":         "-1",
	}
	for key, value := range overrides {
		attributes[key] = value
	}
	return &sdktf.InstanceState{ID: mongodbBackupPolicyFixtureInstanceId, Attributes: attributes}
}

func mongodbBackupPolicyFixtureConfig(overrides map[string]interface{}) map[string]interface{} {
	raw := map[string]interface{}{
		"engine_version":      "4.2",
		"db_instance_class":   "mongo.x8.medium",
		"db_instance_storage": 20,
	}
	for key, value := range overrides {
		raw[key] = value
	}
	return raw
}

func mongodbBackupPolicyFixtureUpdateData(t *testing.T, resourceInstance *schema.Resource, state *sdktf.InstanceState, raw map[string]interface{}, rawConfig map[string]cty.Value) *schema.ResourceData {
	t.Helper()
	diff, err := resourceInstance.Diff(context.Background(), state, sdktf.NewResourceConfigRaw(raw), nil)
	if err != nil {
		t.Fatalf("failed to build the update diff: %s", err)
	}
	diff.RawConfig = mongodbBackupPolicyFixtureRawConfigValue(rawConfig)
	d, err := schema.InternalMap(resourceInstance.Schema).Data(state, diff)
	if err != nil {
		t.Fatalf("failed to build the update resource data: %s", err)
	}
	d.SetId(mongodbBackupPolicyFixtureInstanceId)
	return d
}

func newMongoDBBackupPolicyFixtureResource() *schema.Resource {
	resourceInstance := resourceAliCloudMongoDBInstance()
	resourceInstance.Timeouts = &schema.ResourceTimeout{Update: schema.DefaultTimeout(60 * time.Second)}
	return resourceInstance
}

// mongodbBackupPolicyFixtureRunUpdate drives the public resource update entry and
// returns the single ModifyBackupPolicy request it serialized.
func mongodbBackupPolicyFixtureRunUpdate(t *testing.T, state *sdktf.InstanceState, raw map[string]interface{}, rawConfig map[string]cty.Value) map[string]string {
	t.Helper()
	fixture, client := newMongoDBBackupPolicyFixture(t)
	resourceInstance := newMongoDBBackupPolicyFixtureResource()
	d := mongodbBackupPolicyFixtureUpdateData(t, resourceInstance, state, raw, rawConfig)
	if err := resourceInstance.Update(d, client); err == nil {
		t.Fatal("the fixture answers ModifyBackupPolicy with 400, so Update must return the API error")
	}
	requests := fixture.captured("ModifyBackupPolicy")
	if len(requests) != 1 {
		t.Fatalf("expected exactly 1 ModifyBackupPolicy request, got %d", len(requests))
	}
	return requests[0]
}

func mongodbBackupPolicyFixtureRunUpdateError(t *testing.T, state *sdktf.InstanceState, raw map[string]interface{}, rawConfig map[string]cty.Value) (string, []map[string]string) {
	t.Helper()
	fixture, client := newMongoDBBackupPolicyFixture(t)
	resourceInstance := newMongoDBBackupPolicyFixtureResource()
	d := mongodbBackupPolicyFixtureUpdateData(t, resourceInstance, state, raw, rawConfig)
	err := resourceInstance.Update(d, client)
	if err == nil {
		t.Fatal("the update must be rejected before any request is sent")
	}
	return err.Error(), fixture.captured("ModifyBackupPolicy")
}

func mongodbBackupPolicyFixtureAssertOmitted(t *testing.T, request map[string]string, keys ...string) {
	t.Helper()
	for _, key := range keys {
		if value, ok := request[key]; ok {
			t.Fatalf("the ModifyBackupPolicy request must omit %s, got %q", key, value)
		}
	}
}

func mongodbBackupPolicyFixtureAssertSent(t *testing.T, request map[string]string, key, want string) {
	t.Helper()
	value, ok := request[key]
	if !ok {
		t.Fatalf("the ModifyBackupPolicy request must carry %s=%q, got keys %v", key, want, request)
	}
	if value != want {
		t.Fatalf("the ModifyBackupPolicy request must carry %s=%q, got %q", key, want, value)
	}
}

func mongodbBackupPolicyFixtureCrossConfig(t *testing.T, enableCrossLogBackup int) (map[string]interface{}, map[string]cty.Value) {
	t.Helper()
	raw := mongodbBackupPolicyFixtureConfig(map[string]interface{}{
		"cross_backup_period":       []interface{}{"Monday"},
		"cross_backup_type":         "update",
		"cross_retention_type":      "never",
		"cross_log_retention_type":  "never",
		"cross_log_retention_value": 7,
		"dest_region":               "cn-shanghai",
		"src_region":                "cn-beijing",
		"enable_cross_log_backup":   enableCrossLogBackup,
	})
	values := map[string]cty.Value{
		"cross_backup_period":       cty.SetVal([]cty.Value{cty.StringVal("Monday")}),
		"cross_backup_type":         cty.StringVal("update"),
		"cross_retention_type":      cty.StringVal("never"),
		"cross_log_retention_type":  cty.StringVal("never"),
		"cross_log_retention_value": cty.NumberIntVal(7),
		"dest_region":               cty.StringVal("cn-shanghai"),
		"src_region":                cty.StringVal("cn-beijing"),
		"enable_cross_log_backup":   cty.NumberIntVal(int64(enableCrossLogBackup)),
	}
	return raw, values
}

// TestUnitAliCloudMongoDBInstanceBackupPolicyRequestSerialization locks the wire
// contract of the geo-redundant backup attributes: legitimate zero values are sent
// when the configuration sets them, and the high-frequency backup attributes are only
// replayed while the snapshot backup type and interval accept them. The pre-fix
// implementation used GetOk, so the explicit 0/false assertions captured no key, and
// it replayed the high-frequency attributes for Standard instances.
func TestUnitAliCloudMongoDBInstanceBackupPolicyRequestSerialization(t *testing.T) {
	t.Run("explicit zero disables cross log backup", func(t *testing.T) {
		raw, values := mongodbBackupPolicyFixtureCrossConfig(t, 0)
		state := mongodbBackupPolicyFixtureState(map[string]string{"enable_cross_log_backup": "1"})
		request := mongodbBackupPolicyFixtureRunUpdate(t, state, raw, values)
		mongodbBackupPolicyFixtureAssertSent(t, request, "EnableCrossLogBackup", "0")
	})

	t.Run("explicit false disables hourly sparse backup", func(t *testing.T) {
		raw := mongodbBackupPolicyFixtureConfig(map[string]interface{}{"preserve_one_each_hour": false})
		values := map[string]cty.Value{"preserve_one_each_hour": cty.False}
		state := mongodbBackupPolicyFixtureState(map[string]string{
			"snapshot_backup_type":   "Flash",
			"backup_interval":        "15",
			"preserve_one_each_hour": "true",
		})
		request := mongodbBackupPolicyFixtureRunUpdate(t, state, raw, values)
		mongodbBackupPolicyFixtureAssertSent(t, request, "PreserveOneEachHour", "false")
	})

	t.Run("unrelated update omits the unconfigured cross log switch", func(t *testing.T) {
		raw := mongodbBackupPolicyFixtureConfig(map[string]interface{}{"backup_retention_period": 35})
		values := map[string]cty.Value{"backup_retention_period": cty.NumberIntVal(35)}
		state := mongodbBackupPolicyFixtureState(map[string]string{"enable_cross_log_backup": "1"})
		request := mongodbBackupPolicyFixtureRunUpdate(t, state, raw, values)
		mongodbBackupPolicyFixtureAssertSent(t, request, "BackupRetentionPeriod", "35")
		mongodbBackupPolicyFixtureAssertOmitted(t, request, "EnableCrossLogBackup", "PreserveOneEachHour", "HighFrequencyBackupRetention")
	})

	t.Run("standard snapshot type omits high frequency attributes", func(t *testing.T) {
		raw := mongodbBackupPolicyFixtureConfig(map[string]interface{}{"backup_retention_period": 35})
		values := map[string]cty.Value{"backup_retention_period": cty.NumberIntVal(35)}
		state := mongodbBackupPolicyFixtureState(map[string]string{
			"high_frequency_backup_retention": "7",
			"preserve_one_each_hour":          "true",
		})
		request := mongodbBackupPolicyFixtureRunUpdate(t, state, raw, values)
		mongodbBackupPolicyFixtureAssertOmitted(t, request, "HighFrequencyBackupRetention", "PreserveOneEachHour")
	})

	t.Run("flash minute interval replays high frequency attributes", func(t *testing.T) {
		raw := mongodbBackupPolicyFixtureConfig(map[string]interface{}{"backup_retention_period": 35})
		values := map[string]cty.Value{"backup_retention_period": cty.NumberIntVal(35)}
		state := mongodbBackupPolicyFixtureState(map[string]string{
			"snapshot_backup_type":            "Flash",
			"backup_interval":                 "15",
			"high_frequency_backup_retention": "7",
			"preserve_one_each_hour":          "false",
		})
		request := mongodbBackupPolicyFixtureRunUpdate(t, state, raw, values)
		mongodbBackupPolicyFixtureAssertSent(t, request, "HighFrequencyBackupRetention", "7")
		mongodbBackupPolicyFixtureAssertSent(t, request, "PreserveOneEachHour", "false")
	})

	t.Run("complete update mode is serialized", func(t *testing.T) {
		raw, values := mongodbBackupPolicyFixtureCrossConfig(t, 1)
		state := mongodbBackupPolicyFixtureState(map[string]string{"enable_cross_log_backup": "0"})
		request := mongodbBackupPolicyFixtureRunUpdate(t, state, raw, values)
		mongodbBackupPolicyFixtureAssertSent(t, request, "CrossBackupType", "update")
		mongodbBackupPolicyFixtureAssertSent(t, request, "CrossBackupPeriod", "Monday")
		mongodbBackupPolicyFixtureAssertSent(t, request, "CrossRetentionType", "never")
		mongodbBackupPolicyFixtureAssertSent(t, request, "CrossLogRetentionType", "never")
		mongodbBackupPolicyFixtureAssertSent(t, request, "DestRegion", "cn-shanghai")
		mongodbBackupPolicyFixtureAssertSent(t, request, "SrcRegion", "cn-beijing")
		mongodbBackupPolicyFixtureAssertSent(t, request, "EnableCrossLogBackup", "1")
	})
}

// TestUnitAliCloudMongoDBInstanceCrossBackupValidation locks the local contract check
// that runs before the ModifyBackupPolicy request is built: incomplete geo-redundant
// combinations are rejected with the missing attribute named, and no request leaves
// the provider.
func TestUnitAliCloudMongoDBInstanceCrossBackupValidation(t *testing.T) {
	t.Run("enabling without the cross log switch is rejected", func(t *testing.T) {
		// Changing the operation selector to "update" requires the user to state the log
		// backup switch explicitly; a value that only exists in state must not silently
		// decide it.
		raw, values := mongodbBackupPolicyFixtureCrossConfig(t, 1)
		delete(raw, "enable_cross_log_backup")
		delete(values, "enable_cross_log_backup")
		state := mongodbBackupPolicyFixtureState(map[string]string{"enable_cross_log_backup": "1"})
		errText, requests := mongodbBackupPolicyFixtureRunUpdateError(t, state, raw, values)
		if !strings.Contains(errText, "enable_cross_log_backup") {
			t.Fatalf("the error must name the missing enable_cross_log_backup, got %s", errText)
		}
		if len(requests) != 0 {
			t.Fatalf("an incomplete combination must not reach the API, got %d requests", len(requests))
		}
	})

	t.Run("geo attributes without the operation selector are rejected", func(t *testing.T) {
		raw := mongodbBackupPolicyFixtureConfig(map[string]interface{}{"dest_region": "cn-shanghai"})
		values := map[string]cty.Value{"dest_region": cty.StringVal("cn-shanghai")}
		errText, requests := mongodbBackupPolicyFixtureRunUpdateError(t, mongodbBackupPolicyFixtureState(nil), raw, values)
		if !strings.Contains(errText, "cross_backup_type") {
			t.Fatalf("the error must name the missing cross_backup_type, got %s", errText)
		}
		if len(requests) != 0 {
			t.Fatalf("a partial geo-redundant configuration must not reach the API, got %d requests", len(requests))
		}
	})

	t.Run("removed cross log switch on an unrelated update replays the state value", func(t *testing.T) {
		// Optional+Computed keeps the state value when the attribute is removed from the
		// configuration, so an unrelated backup update must not fail and must carry the
		// effective value.
		raw, values := mongodbBackupPolicyFixtureCrossConfig(t, 1)
		delete(raw, "enable_cross_log_backup")
		delete(values, "enable_cross_log_backup")
		raw["backup_retention_period"] = 35
		values["backup_retention_period"] = cty.NumberIntVal(35)
		state := mongodbBackupPolicyFixtureState(map[string]string{
			"cross_backup_type":       "update",
			"enable_cross_log_backup": "1",
		})
		request := mongodbBackupPolicyFixtureRunUpdate(t, state, raw, values)
		mongodbBackupPolicyFixtureAssertSent(t, request, "EnableCrossLogBackup", "1")
		mongodbBackupPolicyFixtureAssertSent(t, request, "BackupRetentionPeriod", "35")
	})

	t.Run("update mode without the destination region is rejected", func(t *testing.T) {
		raw, values := mongodbBackupPolicyFixtureCrossConfig(t, 1)
		delete(raw, "dest_region")
		delete(values, "dest_region")
		errText, requests := mongodbBackupPolicyFixtureRunUpdateError(t, mongodbBackupPolicyFixtureState(nil), raw, values)
		if !strings.Contains(errText, "dest_region") {
			t.Fatalf("the error must name the missing dest_region, got %s", errText)
		}
		if len(requests) != 0 {
			t.Fatalf("an incomplete combination must not reach the API, got %d requests", len(requests))
		}
	})

	t.Run("delay retention without a value is rejected", func(t *testing.T) {
		raw, values := mongodbBackupPolicyFixtureCrossConfig(t, 1)
		raw["cross_retention_type"] = "delay"
		values["cross_retention_type"] = cty.StringVal("delay")
		errText, requests := mongodbBackupPolicyFixtureRunUpdateError(t, mongodbBackupPolicyFixtureState(nil), raw, values)
		if !strings.Contains(errText, "cross_retention_value") {
			t.Fatalf("the error must name the missing cross_retention_value, got %s", errText)
		}
		if len(requests) != 0 {
			t.Fatalf("an incomplete combination must not reach the API, got %d requests", len(requests))
		}
	})
}
