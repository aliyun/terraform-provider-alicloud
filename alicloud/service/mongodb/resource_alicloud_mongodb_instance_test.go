package mongodb_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/mongodb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_mongodb_instance", &resource.Sweeper{
		Name: "alicloud_mongodb_instance",
		F:    testSweepMongoDBInstances,
	})
}

// mongodbInstanceSweeperPrefixes are the test-name prefixes the sweeper owns. The
// resource documentation example uses a fixed name that does not carry a prefix, so it
// is matched separately and exactly: matching it as a prefix would also delete
// unrelated user instances whose name only starts with the same text.
var mongodbInstanceSweeperPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
}

const mongodbInstanceSweeperDocumentedName = "terraform-example"

// mongodbInstanceSweeperMatches reports whether the sweeper owns the instance with this
// description. It must recognize both the test names and the documentation example name
// so a failed run is still cleanable, and it must leave every other instance alone.
func mongodbInstanceSweeperMatches(description string) bool {
	description = strings.ToLower(strings.TrimSpace(description))
	for _, prefix := range mongodbInstanceSweeperPrefixes {
		if strings.HasPrefix(description, strings.ToLower(prefix)) {
			return true
		}
	}
	return description == mongodbInstanceSweeperDocumentedName
}

func testSweepMongoDBInstances(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	action := "DescribeDBInstances"
	request := map[string]interface{}{}
	request["RegionId"] = client.RegionId
	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1
	request["ChargeType"] = "PostPaid"

	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
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
			if !sweep.SweepAll() {
				if mongodbInstanceSweeperMatches(fmt.Sprint(item["DBInstanceDescription"])) {
					skip = false
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
		if len(result) < helper.PageSizeLarge {
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

// mongodbCrossBackupDestRegion returns a geo-redundant backup destination that is
// guaranteed to differ from the test region. cn-shanghai and cn-beijing are the
// regions used by the existing cross-region backup runs; the fallback keeps the pair
// valid if the test region itself is overridden to cn-shanghai.
func mongodbCrossBackupDestRegion() string {
	if mongodbTestRegion() == "cn-shanghai" {
		return "cn-beijing"
	}
	return "cn-shanghai"
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
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region(mongodbTestRegion())})
	serverFunc := func() interface{} {
		return mongodb.NewMongoDBService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	ra := acctest.ResourceAttrInit(resourceId, AliCloudMongoDBInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBInstanceClassicConfig%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudMongoDBInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
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
						"vswitch_id":          acctest.CHECKSET,
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
						"security_group_id": acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
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
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region(mongodbTestRegion())})
	serverFunc := func() interface{} {
		return mongodb.NewMongoDBService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	ra := acctest.ResourceAttrInit(resourceId, AliCloudMongoDBInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBInstanceSecurityIpGroups%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudMongoDBInstanceNetworkDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version":              "4.2",
					"db_instance_class":           "mongo.x8.medium",
					"db_instance_storage":         "20",
					"vswitch_id":                  "${alicloud_vswitch.default.id}",
					"account_password_wo":         "YourPassword123!",
					"account_password_wo_version": 1,
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
						"account_password_wo_version": "1",
						"security_ip_groups.#":        "2",
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
					"security_ip_groups": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_groups.#": "0",
					}),
				),
			},
			{
				// The shared Create/Update backup-policy builder changed in this branch, so
				// the pre-existing Parameters path also needs same-head execution evidence
				// (ACC-SCOPE-DIFFERENCE). The element-level assertion pins the exact
				// name/value pair rather than only the set size, and the independent
				// DescribeParameters check proves the modify reached the service instead of
				// reading a stale configuration value back from the state.
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
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "parameters.*", map[string]string{
						"name":  "operationProfiling.slowOpThresholdMs",
						"value": "80",
					}),
					testAccCheckMongoDBInstanceParameterValue(resourceId, "operationProfiling.slowOpThresholdMs", "80"),
				),
			},
			{
				// Each successful configuration step is followed by its own import. The
				// parameters set is exempt because Read refreshes only the parameters that
				// are present in the configuration, so an imported state (which has no
				// configuration) cannot reproduce them; the exact value is proven by the
				// element assertion and the independent DescribeParameters check above.
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "account_password_wo", "account_password_wo_version", "kms_encrypted_password", "kms_encryption_context", "ssl_action", "effective_time", "order_type", "parameters"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password_wo":         "YourPassword123!Update",
					"account_password_wo_version": 2,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password_wo_version": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password":            "YourPassword123!password",
					"account_password_wo":         acctest.REMOVEKEY,
					"account_password_wo_version": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password":            "YourPassword123!password",
						"account_password_wo_version": acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "account_password_wo", "account_password_wo_version", "kms_encrypted_password", "kms_encryption_context", "ssl_action", "effective_time", "order_type", "parameters"},
			},
		},
	})
}

func TestAccAliCloudMongoDBInstance_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_instance.default"
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region(mongodbTestRegion())})
	serverFunc := func() interface{} {
		return mongodb.NewMongoDBService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	ra := acctest.ResourceAttrInit(resourceId, AliCloudMongoDBInstanceMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBInstanceClassicConfig%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudMongoDBInstanceBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"vswitch_id":          acctest.CHECKSET,
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
						"secondary_zone_id": acctest.CHECKSET,
						"hidden_zone_id":    acctest.CHECKSET,
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
						"security_group_id": acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
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
					"db_instance_class":      acctest.REMOVEKEY,
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
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region(mongodbTestRegion())})
	serverFunc := func() interface{} {
		return mongodb.NewMongoDBService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	ra := acctest.ResourceAttrInit(resourceId, AliCloudMongoDBInstanceMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBInstanceClassicConfig%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudMongoDBInstanceBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"vpc_id":                    acctest.CHECKSET,
						"vswitch_id":                acctest.CHECKSET,
						"zone_id":                   acctest.CHECKSET,
						"secondary_zone_id":         acctest.CHECKSET,
						"hidden_zone_id":            acctest.CHECKSET,
						"security_group_id":         acctest.CHECKSET,
						"network_type":              "VPC",
						"name":                      name,
						"instance_charge_type":      "PostPaid",
						"security_ip_list.#":        "1",
						"encrypted":                 "true",
						"cloud_disk_encryption_key": acctest.CHECKSET,
						"replication_factor":        "3",
						"readonly_replicas":         "2",
						"resource_group_id":         acctest.CHECKSET,
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

// mongodbInstanceImportStateVerifyIgnore lists the input-only attributes the instance
// tests exempt from import verification: the backup and policy APIs never return them,
// so an imported state cannot reproduce the configured value.
var mongodbInstanceImportStateVerifyIgnore = []string{
	"account_password",
	"kms_encrypted_password",
	"kms_encryption_context",
	"ssl_action",
	"effective_time",
	"order_type",
	"parameters",
}

// TestAccAliCloudMongoDBInstance_crossBackup verifies the cross-region
// (geo-redundant) backup attributes and high-frequency backup retention
// introduced for the MongoDB replica set instance. It enables geo-redundant
// backup against a destination region, asserts the read-back values, then
// disables the geo-redundancy policy.
func TestAccAliCloudMongoDBInstance_crossBackup(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_instance.default"
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region(mongodbTestRegion())})
	serverFunc := func() interface{} {
		return mongodb.NewMongoDBService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	ra := acctest.ResourceAttrInit(resourceId, AliCloudMongoDBInstanceCrossBackupMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBInstanceCrossBackup%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudMongoDBInstanceNetworkDependence)
	srcRegion := mongodbTestRegion()
	destRegion := mongodbCrossBackupDestRegion()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		IDRefreshName: resourceId,
		// tags and global_security_group_list are Optional-only (not Computed) and are
		// filled by Read alone, so the ID-only refresh comparison always reports their
		// collection headers (tags.%, global_security_group_list.#) as extra. Ignore
		// only these two read-filled collections; the cross-region attributes under test
		// stay inside the ID-only check.
		IDRefreshIgnore:   []string{"tags", "global_security_group_list"},
		ProviderFactories: acctest.SDKProviderFactories,
		// The case creates a real instance and every dependency below it, so the
		// instance deletion must be asserted instead of relying on one external cleanup
		// (QA-N2). The fixture dependency names carry the per-run name, which makes the
		// post-run dependency readback attributable.
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Create the instance with the complete geo-redundant policy in the
				// initial configuration: Create finishes with the shared Update path, so
				// this legally exercises the new attributes on a brand-new instance
				// instead of only through later updates (QA-N2).
				//
				// Cross-region backup requires an ESSD (cloud disk) instance
				// (https://help.aliyun.com/zh/mongodb/user-guide/cross-region-backup):
				// on a local-disk instance the service accepts the policy request but never
				// reports a geo-redundant policy back, so the steps below could not verify
				// anything. ESSD AutoPL (cloud_auto) is only sold on the mdb.shard.*
				// cloud-disk class family; the spec below is the combination the
				// repository already provisions with cloud_auto in
				// TestAccAliCloudMongoDBInstance_basic1_twin.
				Config: testAccConfig(map[string]interface{}{
					"engine_version":            "8.0",
					"db_instance_class":         "mdb.shard.2x.xlarge.d",
					"db_instance_storage":       "80",
					"storage_type":              "cloud_auto",
					"provisioned_iops":          "2000",
					"vswitch_id":                "${alicloud_vswitch.default.id}",
					"backup_time":               "02:00Z-03:00Z",
					"backup_period":             []string{"Monday", "Tuesday"},
					"cross_backup_period":       []string{"Monday", "Tuesday"},
					"cross_backup_type":         "update",
					"cross_retention_type":      "delay",
					"cross_retention_value":     7,
					"cross_log_retention_type":  "delay",
					"cross_log_retention_value": 7,
					"dest_region":               destRegion,
					"enable_cross_log_backup":   1,
					"src_region":                srcRegion,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version":      "8.0",
						"db_instance_class":   "mdb.shard.2x.xlarge.d",
						"db_instance_storage": "80",
						"storage_type":        "cloud_auto",
						"provisioned_iops":    "2000",
						"vswitch_id":          acctest.CHECKSET,
						"backup_time":         "02:00Z-03:00Z",
						"backup_period.#":     "2",
						// The shared attribute map seeds cross_backup_period with CHECKSET;
						// a TypeSet must be checked by its element count, so replace the
						// seeded value instead of keeping both keys.
						"cross_backup_period":             acctest.REMOVEKEY,
						"cross_backup_period.#":           "2",
						"cross_backup_type":               "update",
						"cross_retention_type":            "delay",
						"cross_retention_value":           "7",
						"cross_log_retention_type":        "delay",
						"cross_log_retention_value":       "7",
						"dest_region":                     destRegion,
						"enable_cross_log_backup":         "1",
						"src_region":                      srcRegion,
						"high_frequency_backup_retention": acctest.REMOVEKEY,
						"preserve_one_each_hour":          acctest.REMOVEKEY,
					}),
					// Read keeps the configured value when DescribeBackupPolicy omits a
					// key, so the state assertions above cannot prove the create-time write
					// reached the service. Compare every geo-redundant field against the
					// API response.
					testAccCheckMongoDBInstanceCrossBackupPolicyValues(resourceId, []string{"Monday", "Tuesday"}, map[string]string{
						"CrossRetentionType":     "delay",
						"CrossRetentionValue":    "7",
						"CrossLogRetentionType":  "delay",
						"CrossLogRetentionValue": "7",
						"DestRegion":             destRegion,
						"SrcRegion":              srcRegion,
						"EnableCrossLogBackup":   "1",
					}),
				),
			},
			{
				// Import of the active policy must rebuild the same managed state;
				// cross_backup_type is the write-only selector the API does not return.
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: append([]string{"cross_backup_type"}, mongodbInstanceImportStateVerifyIgnore...),
			},
			{
				// Second legal geo-redundant policy: another weekday pair and larger
				// retention values. cross_backup_period must stay a subset of
				// backup_period, so the conventional schedule grows in the same apply
				// (QA-N4). The API readback proves the update converged instead of
				// leaving the initial Monday/Tuesday and 7/7 values.
				Config: testAccConfig(map[string]interface{}{
					"backup_period":             []string{"Monday", "Tuesday", "Wednesday", "Thursday"},
					"cross_backup_period":       []string{"Wednesday", "Thursday"},
					"cross_retention_value":     14,
					"cross_log_retention_value": 10,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#":           "4",
						"cross_backup_period.#":     "2",
						"cross_retention_value":     "14",
						"cross_log_retention_value": "10",
					}),
					testAccCheckMongoDBInstanceCrossBackupPolicyValues(resourceId, []string{"Wednesday", "Thursday"}, map[string]string{
						"CrossRetentionType":     "delay",
						"CrossRetentionValue":    "14",
						"CrossLogRetentionType":  "delay",
						"CrossLogRetentionValue": "10",
						"EnableCrossLogBackup":   "1",
					}),
				),
			},
			{
				// The changed period and retention values are readable, so the import
				// after them keeps every attribute except the write-only selector under
				// verification.
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: append([]string{"cross_backup_type"}, mongodbInstanceImportStateVerifyIgnore...),
			},
			{
				// The other legal retention types are "never" (no expiry). Both types
				// change together and the API readback proves the service converges on the
				// non-delay values (QA-N4).
				Config: testAccConfig(map[string]interface{}{
					"cross_retention_type":     "never",
					"cross_log_retention_type": "never",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						// The service rejects a retention value under a "never" type
						// ("never retention type do not support retention value"), so the
						// value configured by an earlier step is no longer part of the
						// policy and must not be asserted from the inherited check map. The
						// type readback below is the assertion for this transition.
						"cross_retention_value":     acctest.REMOVEKEY,
						"cross_log_retention_value": acctest.REMOVEKEY,
						"cross_retention_type":      "never",
						"cross_log_retention_type":  "never",
					}),
					testAccCheckMongoDBInstanceCrossBackupPolicyValues(resourceId, []string{"Wednesday", "Thursday"}, map[string]string{
						"CrossRetentionType":    "never",
						"CrossLogRetentionType": "never",
						"EnableCrossLogBackup":  "1",
					}),
				),
			},
			{
				// The retention types are read back from the API. The retention values are
				// not part of a "never" policy (the service rejects them for that type), so
				// they cannot be compared on import; the write-only selector is skipped for
				// the same structural reason as above.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: append([]string{
					"cross_backup_type",
					"cross_retention_value",
					"cross_log_retention_value",
				}, mongodbInstanceImportStateVerifyIgnore...),
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
				// The explicit zero must survive an import: the API-backed switch is
				// readable. The cross log retention pair is exempted because the service
				// stops reporting it once cross-region log backup is disabled (observed on
				// this head: the imported state has no CrossLogRetentionType/Value), while
				// the resource keeps the configured values by design. The cross retention
				// value is exempted for the same structural reason as the "never" import
				// above: the retention type is "never", so the service does not report a
				// retention value.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: append([]string{
					"cross_backup_type",
					"cross_retention_value",
					"cross_log_retention_type",
					"cross_log_retention_value",
				}, mongodbInstanceImportStateVerifyIgnore...),
			},
			{
				// The delete branch of ModifyBackupPolicy locates the policy by its source
				// and destination regions, so both are kept here; the retention and period
				// attributes belong to the update branch and must not accompany a delete.
				Config: testAccConfig(map[string]interface{}{
					"cross_backup_period":       acctest.REMOVEKEY,
					"cross_retention_type":      acctest.REMOVEKEY,
					"cross_retention_value":     acctest.REMOVEKEY,
					"cross_log_retention_type":  acctest.REMOVEKEY,
					"cross_log_retention_value": acctest.REMOVEKEY,
					"enable_cross_log_backup":   acctest.REMOVEKEY,
					"src_region":                srcRegion,
					"dest_region":               destRegion,
					"cross_backup_type":         "delete",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cross_backup_type": "delete",
						"src_region":        srcRegion,
						"dest_region":       destRegion,
					}),
					// Assert every non-selector cross field is absent from the real policy
					// response instead of only checking the period; the Import exemption
					// below is derived from this response (R-2).
					testAccCheckMongoDBInstanceCrossBackupPolicyCleared(resourceId),
				),
			},
			{
				// Import after the policy is deleted still has to rebuild a usable
				// resource state. Each exempted attribute has its own reason:
				// - cross_backup_type: write-only selector, never returned by the API.
				// - cross_backup_period, cross_retention_type/value, cross_log_retention_type/value
				//   and enable_cross_log_backup: inputs of the policy that the delete removed;
				//   the check above asserted every one of them is absent from the response,
				//   and the resource retains the configured values by design.
				// - src_region / dest_region: the policy locator, retained from the
				//   configuration for a subsequent delete and not reported once the policy
				//   is gone.
				// The enable/zero imports above keep the same attributes under full
				// verification while the policy still exists.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: append([]string{
					"cross_backup_type",
					"cross_backup_period",
					"cross_retention_type",
					"cross_retention_value",
					"cross_log_retention_type",
					"cross_log_retention_value",
					"src_region",
					"dest_region",
					"enable_cross_log_backup",
				}, mongodbInstanceImportStateVerifyIgnore...),
			},
			{
				// Change the snapshot interval to Flash/15 and set both high-frequency
				// attributes in the same apply. The provider sends the interval and the
				// retention in one ModifyBackupPolicy request; this step proves that the
				// service converges instead of staying on the default 1/false (R-4).
				Config: testAccConfig(map[string]interface{}{
					"snapshot_backup_type":            "Flash",
					"backup_interval":                 "15",
					"high_frequency_backup_retention": 7,
					"preserve_one_each_hour":          true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snapshot_backup_type":            "Flash",
						"backup_interval":                 "15",
						"high_frequency_backup_retention": "7",
						"preserve_one_each_hour":          "true",
					}),
					testAccCheckMongoDBInstanceCrossBackupPolicyValues(resourceId, nil, map[string]string{
						"HighFrequencyBackupRetention": "7",
						"PreserveOneEachHour":          "true",
						"BackupInterval":               "15",
						"SnapshotBackupType":           "Flash",
					}),
				),
			},
			{
				// The high-frequency values are readable, so an import after enabling them
				// keeps them under verification; only the policy attributes exempted for
				// the deleted policy are skipped.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: append([]string{
					"cross_backup_type",
					"cross_backup_period",
					"cross_retention_type",
					"cross_retention_value",
					"cross_log_retention_type",
					"cross_log_retention_value",
					"src_region",
					"dest_region",
					"enable_cross_log_backup",
				}, mongodbInstanceImportStateVerifyIgnore...),
			},
			{
				// The contract enables high-frequency backup for Flash snapshots with any
				// positive interval, not only for 15. Changing the interval to 60 while both
				// high-frequency values stay unchanged must still resend them in the full
				// policy request instead of leaving the retained values to the service
				// default (QA-N6).
				Config: testAccConfig(map[string]interface{}{
					"backup_interval": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_interval":                 "60",
						"high_frequency_backup_retention": "7",
						"preserve_one_each_hour":          "true",
					}),
					testAccCheckMongoDBInstanceCrossBackupPolicyValues(resourceId, nil, map[string]string{
						"HighFrequencyBackupRetention": "7",
						"PreserveOneEachHour":          "true",
						"BackupInterval":               "60",
						"SnapshotBackupType":           "Flash",
					}),
				),
			},
			{
				// Keep the same per-step import contract for the interval change: the
				// high-frequency values and the 60-minute interval are readable, so the
				// import verifies them; only the policy attributes exempted for the
				// deleted policy are skipped.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: append([]string{
					"cross_backup_type",
					"cross_backup_period",
					"cross_retention_type",
					"cross_retention_value",
					"cross_log_retention_type",
					"cross_log_retention_value",
					"src_region",
					"dest_region",
					"enable_cross_log_backup",
				}, mongodbInstanceImportStateVerifyIgnore...),
			},
			{
				// A second high-frequency combination: the default one-day retention and
				// the hourly sparse switch turned off. The API readback proves both
				// transitions reach the service instead of being replayed from state
				// (QA-N4).
				Config: testAccConfig(map[string]interface{}{
					"high_frequency_backup_retention": 1,
					"preserve_one_each_hour":          false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"high_frequency_backup_retention": "1",
						"preserve_one_each_hour":          "false",
					}),
					testAccCheckMongoDBInstanceCrossBackupPolicyValues(resourceId, nil, map[string]string{
						"HighFrequencyBackupRetention": "1",
						"PreserveOneEachHour":          "false",
					}),
				),
			},
			{
				// Both high-frequency values are readable, so the import after the second
				// combination keeps them under verification.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: append([]string{
					"cross_backup_type",
					"cross_backup_period",
					"cross_retention_type",
					"cross_retention_value",
					"cross_log_retention_type",
					"cross_log_retention_value",
					"src_region",
					"dest_region",
					"enable_cross_log_backup",
				}, mongodbInstanceImportStateVerifyIgnore...),
			},
			{
				// Removing the two regions from the configuration while the selector stays
				// "delete" is a no-op: both attributes are Optional+Computed, so the last
				// values stay in state and no ModifyBackupPolicy request is sent. Terraform
				// must therefore keep reporting the same values, and the documentation
				// states that removing them does not locate a different policy (R-1).
				Config: testAccConfig(map[string]interface{}{
					"src_region":  acctest.REMOVEKEY,
					"dest_region": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"src_region":  srcRegion,
						"dest_region": destRegion,
					}),
				),
			},
			{
				// The import after the no-op removal must still rebuild the same managed
				// state: the retained regions and the deleted policy attributes stay under
				// the same exemptions as the post-delete import (T03/R7.1).
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: append([]string{
					"cross_backup_type",
					"cross_backup_period",
					"cross_retention_type",
					"cross_retention_value",
					"cross_log_retention_type",
					"cross_log_retention_value",
					"src_region",
					"dest_region",
					"enable_cross_log_backup",
				}, mongodbInstanceImportStateVerifyIgnore...),
			},
		},
	})
}

// testAccCheckMongoDBInstanceCrossBackupPolicyValues reads the backup policy from the
// DDS API and compares the exact values the configuration asked for. SrcRegion is
// supplied when the state has it so the service can resolve the geo-redundant entry
// without relying on the instance's own region; the policy fields are also reported
// without it while the policy is active. Terraform state is not consulted: Read keeps
// the configured value when the response omits a key, so a state-only assertion cannot
// prove that the write reached the service.
func testAccCheckMongoDBInstanceCrossBackupPolicyValues(resourceId string, expectPeriod []string, expected map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s is not present in the state", resourceId)
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		request := map[string]interface{}{"DBInstanceId": rs.Primary.ID}
		if v, ok := rs.Primary.Attributes["src_region"]; ok && v != "" {
			request["SrcRegion"] = v
		}
		policy, err := client.RpcPost("Dds", "2015-12-01", "DescribeBackupPolicy", nil, request, true)
		if err != nil {
			return fmt.Errorf("failed to describe the backup policy of %s: %s", rs.Primary.ID, err)
		}
		rawPeriod := fmt.Sprint(policy["CrossBackupPeriod"])
		if expectPeriod != nil {
			gotPeriod := strings.Split(rawPeriod, ",")
			sort.Strings(gotPeriod)
			wantPeriod := append([]string(nil), expectPeriod...)
			sort.Strings(wantPeriod)
			if rawPeriod == "" || strings.Join(gotPeriod, ",") != strings.Join(wantPeriod, ",") {
				return fmt.Errorf("DescribeBackupPolicy must report CrossBackupPeriod=%v, got %q (policy %#v)", wantPeriod, rawPeriod, policy)
			}
		}
		for key, want := range expected {
			got := ""
			if v, ok := policy[key]; ok && v != nil {
				got = fmt.Sprint(v)
			}
			if got != want {
				return fmt.Errorf("DescribeBackupPolicy must report %s=%q, got %q (policy %#v)", key, want, got, policy)
			}
		}
		return nil
	}
}

// testAccCheckMongoDBInstanceCrossBackupPolicyCleared reads the backup policy straight
// from the DDS API after a delete and requires every non-selector geo-redundant field
// to be absent. The Import exemption for the deleted policy is based on this exact
// response instead of an assumption about the whole block (R-2).
func testAccCheckMongoDBInstanceCrossBackupPolicyCleared(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s is not present in the state", resourceId)
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		policy, err := client.RpcPost("Dds", "2015-12-01", "DescribeBackupPolicy", nil, map[string]interface{}{"DBInstanceId": rs.Primary.ID}, true)
		if err != nil {
			return fmt.Errorf("failed to describe the backup policy of %s: %s", rs.Primary.ID, err)
		}
		for _, key := range []string{
			"CrossBackupPeriod",
			"CrossRetentionType",
			"CrossRetentionValue",
			"CrossLogRetentionType",
			"CrossLogRetentionValue",
			"EnableCrossLogBackup",
			"SrcRegion",
			"DestRegion",
		} {
			if v, ok := policy[key]; ok && fmt.Sprint(v) != "" {
				return fmt.Errorf("DescribeBackupPolicy must not report %s=%q after deleting the geo-redundant policy (policy %#v)", key, fmt.Sprint(v), policy)
			}
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
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		service := mongodb.NewMongoDBService(client)
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

// testAccCheckMongoDBInstanceParameterValue reads the running parameter straight from
// the DDS API and compares the exact value the configuration asked for. Read refreshes
// the parameters set from DescribeParameters only for the entries that are present in
// the configuration, so a state-only assertion cannot distinguish a write that reached
// the service from a stale value left in the state.
func testAccCheckMongoDBInstanceParameterValue(resourceId, name, want string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s is not present in the state", resourceId)
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		object, err := client.RpcPost("Dds", "2015-12-01", "DescribeParameters", nil, map[string]interface{}{
			"DBInstanceId": rs.Primary.ID,
			"ExtraParam":   "terraform",
		}, true)
		if err != nil {
			return fmt.Errorf("failed to describe the parameters of %s: %s", rs.Primary.ID, err)
		}
		running, ok := object["RunningParameters"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("DescribeParameters did not return RunningParameters for %s (response %#v)", rs.Primary.ID, object)
		}
		items, ok := running["Parameter"].([]interface{})
		if !ok {
			return fmt.Errorf("DescribeParameters did not return RunningParameters.Parameter for %s (response %#v)", rs.Primary.ID, object)
		}
		for _, item := range items {
			parameter, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if fmt.Sprint(parameter["ParameterName"]) != name {
				continue
			}
			if got := fmt.Sprint(parameter["ParameterValue"]); got != want {
				return fmt.Errorf("DescribeParameters must report %s=%q, got %q (response %#v)", name, want, got, object)
			}
			return nil
		}
		return fmt.Errorf("DescribeParameters did not report parameter %q for %s (response %#v)", name, rs.Primary.ID, object)
	}
}

// TestAccAliCloudMongoDBInstance_validation covers the input-only and
// environment-dependent attributes (TDE, KMS-encrypted password, PrePaid
// billing, and restore source) that cannot be exercised end to end in the
// shared ACC environment. Following the established convention used by peer
// resources, these attributes are set in a config that omits the required
// `db_instance_class`, so the plan is rejected for that missing argument; the
// steps prove only that rejection and are not behavior coverage for any of the
// attributes they list.
//
// The high-frequency backup attributes are not part of these configs: their
// behavior is covered by TestAccAliCloudMongoDBInstance_crossBackup, which
// switches the instance to Flash snapshots with a positive interval and
// reads the values back from ModifyBackupPolicy/DescribeBackupPolicy.
func TestAccAliCloudMongoDBInstance_validation(t *testing.T) {
	resourceId := "alicloud_mongodb_instance.default"
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region(mongodbTestRegion())})
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccMongoDBInstanceValidation%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudMongoDBInstanceNetworkDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
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
				}),
				ExpectError: regexp.MustCompile(`db_instance_class.*is required`),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version":         "4.2",
					"db_instance_storage":    "20",
					"vswitch_id":             "${alicloud_vswitch.default.id}",
					"instance_charge_type":   "PostPaid",
					"period":                 "12",
					"auto_renew":             "false",
					"auto_renew_duration":    "3",
					"tde_status":             "disabled",
					"encryptor_name":         "aes-256-cfb",
					"encryption_key":         "key-update",
					"role_arn":               "acs:ram::1234567890123456:role/example-update",
					"kms_encrypted_password": "kms-password-update",
					"kms_encryption_context": map[string]string{"name": "value-update"},
					"src_db_instance_id":     "src-instance-id-update",
					"restore_time":           "2026-02-01T00:00:00Z",
				}),
				ExpectError: regexp.MustCompile(`db_instance_class.*is required`),
			},
		},
	})
}

var AliCloudMongoDBInstanceCrossBackupMap = map[string]string{
	"cross_backup_period":             acctest.CHECKSET,
	"cross_backup_type":               acctest.CHECKSET,
	"cross_retention_type":            acctest.CHECKSET,
	"cross_retention_value":           acctest.CHECKSET,
	"cross_log_retention_type":        acctest.CHECKSET,
	"cross_log_retention_value":       acctest.CHECKSET,
	"dest_region":                     acctest.CHECKSET,
	"enable_cross_log_backup":         acctest.CHECKSET,
	"high_frequency_backup_retention": acctest.CHECKSET,
	"preserve_one_each_hour":          acctest.CHECKSET,
	"src_region":                      acctest.CHECKSET,
}

var AliCloudMongoDBInstanceMap0 = map[string]string{
	"storage_engine":      acctest.CHECKSET,
	"storage_type":        acctest.CHECKSET,
	"vpc_id":              acctest.CHECKSET,
	"vswitch_id":          acctest.CHECKSET,
	"zone_id":             acctest.CHECKSET,
	"network_type":        acctest.CHECKSET,
	"replication_factor":  acctest.CHECKSET,
	"readonly_replicas":   acctest.CHECKSET,
	"resource_group_id":   acctest.CHECKSET,
	"maintain_start_time": acctest.CHECKSET,
	"maintain_end_time":   acctest.CHECKSET,
	"retention_period":    acctest.CHECKSET,
	"replica_set_name":    acctest.CHECKSET,
	"ssl_status":          acctest.CHECKSET,
	"zone_infos.#":        acctest.CHECKSET,
	"replica_sets.#":      acctest.CHECKSET,
}

var AliCloudMongoDBInstanceMap1 = map[string]string{
	"storage_engine":          acctest.CHECKSET,
	"storage_type":            acctest.CHECKSET,
	"vpc_id":                  acctest.CHECKSET,
	"vswitch_id":              acctest.CHECKSET,
	"zone_id":                 acctest.CHECKSET,
	"network_type":            acctest.CHECKSET,
	"replication_factor":      acctest.CHECKSET,
	"readonly_replicas":       acctest.CHECKSET,
	"resource_group_id":       acctest.CHECKSET,
	"backup_time":             acctest.CHECKSET,
	"backup_retention_period": acctest.CHECKSET,
	"snapshot_backup_type":    acctest.CHECKSET,
	"backup_interval":         acctest.CHECKSET,
	"maintain_start_time":     acctest.CHECKSET,
	"maintain_end_time":       acctest.CHECKSET,
	"retention_period":        acctest.CHECKSET,
	"replica_set_name":        acctest.CHECKSET,
	"ssl_status":              acctest.CHECKSET,
	"zone_infos.#":            acctest.CHECKSET,
	"replica_sets.#":          acctest.CHECKSET,
}

// AliCloudMongoDBInstanceNetworkDependence provides only the network
// dependencies used by the cross-backup, validation and security IP group cases.
// Those cases do not reference the security group, KMS key or global security IP
// groups, so creating them would only add resources that are not exercised and can
// be left behind when the test case aborts (NEW-2). Keeping the fixture minimal
// makes every created resource part of the assertion and cleanup surface.
func AliCloudMongoDBInstanceNetworkDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
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
`, name, mongodbTestZonePrimary())
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
				global_ig_name          = "tfacc${count.index}${substr(md5(var.name), 0, 8)}"
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
				global_ig_name          = "tfacc${count.index}${substr(md5(var.name), 0, 8)}"
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

	// describeBackupPolicyStatus makes DescribeBackupPolicy answer with an error, so a
	// test can prove the update chain stops before ModifyBackupPolicy. Zero keeps the
	// normal successful response. describeBackupPolicyEmpty makes the successful
	// response carry no schedule, which exercises the omitted-value path.
	describeBackupPolicyStatus int
	describeBackupPolicyCode   string
	describeBackupPolicyEmpty  bool
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
	case "DescribeBackupPolicy":
		if f.describeBackupPolicyStatus != 0 {
			w.WriteHeader(f.describeBackupPolicyStatus)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Code":      f.describeBackupPolicyCode,
				"Message":   "fixture rejected the describe request",
				"RequestId": "fixture",
			})
			return
		}
		if f.describeBackupPolicyEmpty {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"PreferredBackupTime":   "23:00Z-24:00Z",
			"PreferredBackupPeriod": "Monday,Tuesday",
		})
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
	return newMongoDBBackupPolicyFixtureWithDescribeResponse(t, 0, "", false)
}

func newMongoDBBackupPolicyFixtureWithDescribeFailure(t *testing.T, status int, code string) (*mongodbBackupPolicyFixture, *connectivity.AliyunClient) {
	t.Helper()
	return newMongoDBBackupPolicyFixtureWithDescribeResponse(t, status, code, false)
}

func newMongoDBBackupPolicyFixtureWithDescribeResponse(t *testing.T, status int, code string, empty bool) (*mongodbBackupPolicyFixture, *connectivity.AliyunClient) {
	t.Helper()
	fixture := &mongodbBackupPolicyFixture{
		describeBackupPolicyStatus: status,
		describeBackupPolicyCode:   code,
		describeBackupPolicyEmpty:  empty,
	}
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
	resourceInstance := mongodb.ResourceAliCloudMongoDBInstance()
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

// mongodbBackupPolicyFixtureAssertNoPolicyChange builds the update data for a plan
// that removes only Optional+Computed attributes and asserts that no backup-policy
// attribute reports a change, so Update cannot issue a ModifyBackupPolicy request.
func mongodbBackupPolicyFixtureAssertNoPolicyChange(t *testing.T, state *sdktf.InstanceState, raw map[string]interface{}, rawConfig map[string]cty.Value) {
	t.Helper()
	resourceInstance := newMongoDBBackupPolicyFixtureResource()
	d := mongodbBackupPolicyFixtureUpdateData(t, resourceInstance, state, raw, rawConfig)
	for _, key := range []string{
		"cross_backup_period",
		"cross_backup_type",
		"cross_retention_type",
		"cross_retention_value",
		"cross_log_retention_type",
		"cross_log_retention_value",
		"dest_region",
		"enable_cross_log_backup",
		"src_region",
	} {
		if d.HasChange(key) {
			t.Fatalf("removing the regions after a delete must not change %s", key)
		}
	}
}

// mongodbBackupPolicyFixtureRunPlanValidation runs the schema validation Terraform
// performs while planning, where RequiredWith is enforced, without entering Update.
func mongodbBackupPolicyFixtureRunPlanValidation(t *testing.T, raw map[string]interface{}) error {
	t.Helper()
	resourceInstance := newMongoDBBackupPolicyFixtureResource()
	diags := resourceInstance.Validate(sdktf.NewResourceConfigRaw(raw))
	if diags.HasError() {
		for _, diagnostic := range diags {
			if diagnostic.Severity == diag.Error {
				return fmt.Errorf("%s: %s", diagnostic.Summary, diagnostic.Detail)
			}
		}
	}
	return nil
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
	// Both retention types are "never", so the configuration legitimately omits both
	// retention values: the builder never sends a value for a "never" type (QA-N5).
	raw := mongodbBackupPolicyFixtureConfig(map[string]interface{}{
		"cross_backup_period":      []interface{}{"Monday"},
		"cross_backup_type":        "update",
		"cross_retention_type":     "never",
		"cross_log_retention_type": "never",
		"dest_region":              "cn-shanghai",
		"src_region":               "cn-beijing",
		"enable_cross_log_backup":  enableCrossLogBackup,
	})
	values := map[string]cty.Value{
		"cross_backup_period":      cty.SetVal([]cty.Value{cty.StringVal("Monday")}),
		"cross_backup_type":        cty.StringVal("update"),
		"cross_retention_type":     cty.StringVal("never"),
		"cross_log_retention_type": cty.StringVal("never"),
		"dest_region":              cty.StringVal("cn-shanghai"),
		"src_region":               cty.StringVal("cn-beijing"),
		"enable_cross_log_backup":  cty.NumberIntVal(int64(enableCrossLogBackup)),
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

	t.Run("flash non-minute interval replays high frequency attributes", func(t *testing.T) {
		// The contract enables high-frequency backup for Flash snapshots with any
		// positive BackupInterval; 60 is a contract-valid interval and must replay the
		// retained values exactly like 15 does (QA-N6).
		raw := mongodbBackupPolicyFixtureConfig(map[string]interface{}{"backup_retention_period": 35})
		values := map[string]cty.Value{"backup_retention_period": cty.NumberIntVal(35)}
		state := mongodbBackupPolicyFixtureState(map[string]string{
			"snapshot_backup_type":            "Flash",
			"backup_interval":                 "60",
			"high_frequency_backup_retention": "7",
			"preserve_one_each_hour":          "true",
		})
		request := mongodbBackupPolicyFixtureRunUpdate(t, state, raw, values)
		mongodbBackupPolicyFixtureAssertSent(t, request, "HighFrequencyBackupRetention", "7")
		mongodbBackupPolicyFixtureAssertSent(t, request, "PreserveOneEachHour", "true")
	})

	t.Run("interval change replays the unchanged high frequency attributes", func(t *testing.T) {
		// Changing the interval from 15 to 60 while the high-frequency values stay in
		// the configuration changes neither of those attributes, so only the applicable
		// predicate can keep them in the full policy request (QA-N6).
		raw := mongodbBackupPolicyFixtureConfig(map[string]interface{}{
			"backup_interval":                 "60",
			"high_frequency_backup_retention": 7,
			"preserve_one_each_hour":          true,
		})
		values := map[string]cty.Value{
			"backup_interval":                 cty.StringVal("60"),
			"high_frequency_backup_retention": cty.NumberIntVal(7),
			"preserve_one_each_hour":          cty.True,
		}
		state := mongodbBackupPolicyFixtureState(map[string]string{
			"snapshot_backup_type":            "Flash",
			"backup_interval":                 "15",
			"high_frequency_backup_retention": "7",
			"preserve_one_each_hour":          "true",
		})
		request := mongodbBackupPolicyFixtureRunUpdate(t, state, raw, values)
		mongodbBackupPolicyFixtureAssertSent(t, request, "BackupInterval", "60")
		mongodbBackupPolicyFixtureAssertSent(t, request, "HighFrequencyBackupRetention", "7")
		mongodbBackupPolicyFixtureAssertSent(t, request, "PreserveOneEachHour", "true")
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

	t.Run("update mode serializes both retention values", func(t *testing.T) {
		// The retention values only accompany a "delay" retention type, so the
		// never-typed fixture above cannot prove they are serialized at all.
		raw := mongodbBackupPolicyFixtureConfig(map[string]interface{}{
			"cross_backup_period":       []interface{}{"Monday"},
			"cross_backup_type":         "update",
			"cross_retention_type":      "delay",
			"cross_retention_value":     7,
			"cross_log_retention_type":  "delay",
			"cross_log_retention_value": 3,
			"dest_region":               "cn-shanghai",
			"src_region":                "cn-beijing",
			"enable_cross_log_backup":   1,
		})
		values := map[string]cty.Value{
			"cross_backup_period":       cty.SetVal([]cty.Value{cty.StringVal("Monday")}),
			"cross_backup_type":         cty.StringVal("update"),
			"cross_retention_type":      cty.StringVal("delay"),
			"cross_retention_value":     cty.NumberIntVal(7),
			"cross_log_retention_type":  cty.StringVal("delay"),
			"cross_log_retention_value": cty.NumberIntVal(3),
			"dest_region":               cty.StringVal("cn-shanghai"),
			"src_region":                cty.StringVal("cn-beijing"),
			"enable_cross_log_backup":   cty.NumberIntVal(1),
		}
		state := mongodbBackupPolicyFixtureState(map[string]string{"enable_cross_log_backup": "0"})
		request := mongodbBackupPolicyFixtureRunUpdate(t, state, raw, values)
		mongodbBackupPolicyFixtureAssertSent(t, request, "CrossRetentionType", "delay")
		mongodbBackupPolicyFixtureAssertSent(t, request, "CrossRetentionValue", "7")
		mongodbBackupPolicyFixtureAssertSent(t, request, "CrossLogRetentionType", "delay")
		mongodbBackupPolicyFixtureAssertSent(t, request, "CrossLogRetentionValue", "3")
	})

	t.Run("never retention types omit the retention values", func(t *testing.T) {
		// The service rejects a retention value when its type is "never" ("never
		// retention type do not support retention value"), and the contract only uses the
		// value for "delay". A value left in state by a delay -> never transition must
		// therefore not be serialized.
		raw := mongodbBackupPolicyFixtureConfig(map[string]interface{}{
			"cross_backup_period":       []interface{}{"Wednesday"},
			"cross_backup_type":         "update",
			"cross_retention_type":      "never",
			"cross_retention_value":     14,
			"cross_log_retention_type":  "never",
			"cross_log_retention_value": 10,
			"dest_region":               "cn-shanghai",
			"src_region":                "cn-beijing",
			"enable_cross_log_backup":   1,
		})
		values := map[string]cty.Value{
			"cross_backup_period":       cty.SetVal([]cty.Value{cty.StringVal("Wednesday")}),
			"cross_backup_type":         cty.StringVal("update"),
			"cross_retention_type":      cty.StringVal("never"),
			"cross_retention_value":     cty.NumberIntVal(14),
			"cross_log_retention_type":  cty.StringVal("never"),
			"cross_log_retention_value": cty.NumberIntVal(10),
			"dest_region":               cty.StringVal("cn-shanghai"),
			"src_region":                cty.StringVal("cn-beijing"),
			"enable_cross_log_backup":   cty.NumberIntVal(1),
		}
		state := mongodbBackupPolicyFixtureState(map[string]string{
			"cross_retention_type":      "delay",
			"cross_retention_value":     "14",
			"cross_log_retention_type":  "delay",
			"cross_log_retention_value": "10",
		})
		request := mongodbBackupPolicyFixtureRunUpdate(t, state, raw, values)
		mongodbBackupPolicyFixtureAssertSent(t, request, "CrossRetentionType", "never")
		mongodbBackupPolicyFixtureAssertSent(t, request, "CrossLogRetentionType", "never")
		mongodbBackupPolicyFixtureAssertOmitted(t, request, "CrossRetentionValue", "CrossLogRetentionValue")
	})

	t.Run("fresh never accepts missing retention values and omits them", func(t *testing.T) {
		// A brand-new configuration with both retention types set to "never" does not
		// provide a retention value at all. The values are never serialized for "never",
		// so requiring one locally forces the user to configure a value the provider
		// then drops. Reaching ModifyBackupPolicy proves the local validation accepted
		// the configuration; the wire assertions prove both values stay out (QA-N5).
		raw, values := mongodbBackupPolicyFixtureCrossConfig(t, 1)
		state := mongodbBackupPolicyFixtureState(map[string]string{"cross_backup_type": "update"})
		request := mongodbBackupPolicyFixtureRunUpdate(t, state, raw, values)
		mongodbBackupPolicyFixtureAssertSent(t, request, "CrossRetentionType", "never")
		mongodbBackupPolicyFixtureAssertSent(t, request, "CrossLogRetentionType", "never")
		mongodbBackupPolicyFixtureAssertOmitted(t, request, "CrossRetentionValue", "CrossLogRetentionValue")
	})

	t.Run("create-time policy reuses the service backup schedule", func(t *testing.T) {
		// A resource created with geo-redundant attributes reaches the policy builder
		// before a Read populated backup_time/backup_period. The effective schedule is
		// read from the service instead of sending empty values that the API rejects.
		raw, values := mongodbBackupPolicyFixtureCrossConfig(t, 1)
		state := mongodbBackupPolicyFixtureState(map[string]string{
			"backup_time": "",
		})
		request := mongodbBackupPolicyFixtureRunUpdate(t, state, raw, values)
		mongodbBackupPolicyFixtureAssertSent(t, request, "PreferredBackupTime", "23:00Z-24:00Z")
		mongodbBackupPolicyFixtureAssertSent(t, request, "PreferredBackupPeriod", "Monday,Tuesday")
	})

	t.Run("create-time policy read failure stops before the modify request", func(t *testing.T) {
		// The create-time fallback reads the effective schedule from the service. A
		// failed read must surface the original error and must not be replaced by a
		// ModifyBackupPolicy request built from empty values (QA-N3).
		fixture, client := newMongoDBBackupPolicyFixtureWithDescribeFailure(t, http.StatusForbidden, "Forbidden.RAM")
		resourceInstance := newMongoDBBackupPolicyFixtureResource()
		raw, values := mongodbBackupPolicyFixtureCrossConfig(t, 1)
		state := mongodbBackupPolicyFixtureState(map[string]string{"backup_time": ""})
		d := mongodbBackupPolicyFixtureUpdateData(t, resourceInstance, state, raw, values)
		err := resourceInstance.Update(d, client)
		if err == nil {
			t.Fatal("a failed DescribeBackupPolicy must fail the update")
		}
		if !strings.Contains(err.Error(), "Forbidden.RAM") {
			t.Fatalf("the update must surface the original DescribeBackupPolicy error, got %s", err)
		}
		if requests := fixture.captured("ModifyBackupPolicy"); len(requests) != 0 {
			t.Fatalf("a failed DescribeBackupPolicy must not reach ModifyBackupPolicy, got %d requests", len(requests))
		}
	})

	t.Run("create-time policy omits an undetermined backup schedule", func(t *testing.T) {
		// DescribeBackupPolicy can succeed without reporting a schedule. An empty value
		// must be omitted instead of being sent as an empty string the service rejects
		// with InvalidPreferredBackupTime/Period (QA-N3).
		fixture, client := newMongoDBBackupPolicyFixtureWithDescribeResponse(t, 0, "", true)
		resourceInstance := newMongoDBBackupPolicyFixtureResource()
		raw, values := mongodbBackupPolicyFixtureCrossConfig(t, 1)
		state := mongodbBackupPolicyFixtureState(map[string]string{
			"backup_time": "",
		})
		d := mongodbBackupPolicyFixtureUpdateData(t, resourceInstance, state, raw, values)
		if err := resourceInstance.Update(d, client); err == nil {
			t.Fatal("the fixture answers ModifyBackupPolicy with 400, so Update must return the API error")
		}
		requests := fixture.captured("ModifyBackupPolicy")
		if len(requests) != 1 {
			t.Fatalf("expected exactly 1 ModifyBackupPolicy request, got %d", len(requests))
		}
		mongodbBackupPolicyFixtureAssertOmitted(t, requests[0], "PreferredBackupTime", "PreferredBackupPeriod")
	})

	t.Run("delete mode carries the regions and omits the retention attributes", func(t *testing.T) {
		// The service assembles the delete entry from the region and type fields only:
		// the policy cannot be located without both regions, and the retention/period
		// fields belong to the update branch. Sending them here used to fail the delete
		// with an opaque backend error.
		raw := mongodbBackupPolicyFixtureConfig(map[string]interface{}{
			"cross_backup_type":         "delete",
			"cross_backup_period":       []interface{}{"Monday"},
			"cross_retention_type":      "never",
			"cross_retention_value":     7,
			"cross_log_retention_type":  "never",
			"cross_log_retention_value": 7,
			"enable_cross_log_backup":   1,
			"src_region":                "cn-beijing",
			"dest_region":               "cn-shanghai",
		})
		values := map[string]cty.Value{
			"cross_backup_type":         cty.StringVal("delete"),
			"cross_backup_period":       cty.SetVal([]cty.Value{cty.StringVal("Monday")}),
			"cross_retention_type":      cty.StringVal("never"),
			"cross_retention_value":     cty.NumberIntVal(7),
			"cross_log_retention_type":  cty.StringVal("never"),
			"cross_log_retention_value": cty.NumberIntVal(7),
			"enable_cross_log_backup":   cty.NumberIntVal(1),
			"src_region":                cty.StringVal("cn-beijing"),
			"dest_region":               cty.StringVal("cn-shanghai"),
		}
		state := mongodbBackupPolicyFixtureState(map[string]string{
			"cross_backup_type": "update",
			"src_region":        "cn-beijing",
			"dest_region":       "cn-shanghai",
		})
		request := mongodbBackupPolicyFixtureRunUpdate(t, state, raw, values)
		mongodbBackupPolicyFixtureAssertSent(t, request, "CrossBackupType", "delete")
		mongodbBackupPolicyFixtureAssertSent(t, request, "SrcRegion", "cn-beijing")
		mongodbBackupPolicyFixtureAssertSent(t, request, "DestRegion", "cn-shanghai")
		mongodbBackupPolicyFixtureAssertOmitted(t, request, "CrossBackupPeriod", "CrossRetentionType", "CrossRetentionValue", "CrossLogRetentionType", "CrossLogRetentionValue", "EnableCrossLogBackup")
	})

	t.Run("delete mode re-issues a delete request when only a region changes", func(t *testing.T) {
		// cross_backup_type is write-only and stays "delete" in both the state and the
		// configuration, so the selector itself reports no change. A region change is
		// still a delete that has to locate the policy, which needs the selector and
		// both regions; gating the delete shape on the selector change alone sent a
		// partial request that could not resolve to any policy.
		raw := mongodbBackupPolicyFixtureConfig(map[string]interface{}{
			"cross_backup_type": "delete",
			"src_region":        "cn-beijing",
			"dest_region":       "cn-hangzhou",
		})
		values := map[string]cty.Value{
			"cross_backup_type": cty.StringVal("delete"),
			"src_region":        cty.StringVal("cn-beijing"),
			"dest_region":       cty.StringVal("cn-hangzhou"),
		}
		state := mongodbBackupPolicyFixtureState(map[string]string{
			"cross_backup_type": "delete",
			"src_region":        "cn-beijing",
			"dest_region":       "cn-shanghai",
		})
		request := mongodbBackupPolicyFixtureRunUpdate(t, state, raw, values)
		mongodbBackupPolicyFixtureAssertSent(t, request, "CrossBackupType", "delete")
		mongodbBackupPolicyFixtureAssertSent(t, request, "SrcRegion", "cn-beijing")
		mongodbBackupPolicyFixtureAssertSent(t, request, "DestRegion", "cn-hangzhou")
		mongodbBackupPolicyFixtureAssertOmitted(t, request, "CrossBackupPeriod", "CrossRetentionType", "CrossRetentionValue", "CrossLogRetentionType", "CrossLogRetentionValue", "EnableCrossLogBackup")
	})

	t.Run("delete mode re-issues a delete request when only the source region changes", func(t *testing.T) {
		raw := mongodbBackupPolicyFixtureConfig(map[string]interface{}{
			"cross_backup_type": "delete",
			"src_region":        "cn-hangzhou",
			"dest_region":       "cn-shanghai",
		})
		values := map[string]cty.Value{
			"cross_backup_type": cty.StringVal("delete"),
			"src_region":        cty.StringVal("cn-hangzhou"),
			"dest_region":       cty.StringVal("cn-shanghai"),
		}
		state := mongodbBackupPolicyFixtureState(map[string]string{
			"cross_backup_type": "delete",
			"src_region":        "cn-beijing",
			"dest_region":       "cn-shanghai",
		})
		request := mongodbBackupPolicyFixtureRunUpdate(t, state, raw, values)
		mongodbBackupPolicyFixtureAssertSent(t, request, "CrossBackupType", "delete")
		mongodbBackupPolicyFixtureAssertSent(t, request, "SrcRegion", "cn-hangzhou")
		mongodbBackupPolicyFixtureAssertSent(t, request, "DestRegion", "cn-shanghai")
	})

	t.Run("delete after update carries the state regions when they leave the configuration", func(t *testing.T) {
		// R-1 transition 1: update -> delete in one apply while both regions are removed
		// from the configuration. They are Optional+Computed, so the values already in
		// state stay effective and the delete must still carry the complete shape.
		raw := mongodbBackupPolicyFixtureConfig(map[string]interface{}{
			"cross_backup_type": "delete",
		})
		values := map[string]cty.Value{
			"cross_backup_type": cty.StringVal("delete"),
		}
		state := mongodbBackupPolicyFixtureState(map[string]string{
			"cross_backup_type": "update",
			"src_region":        "cn-beijing",
			"dest_region":       "cn-shanghai",
		})
		request := mongodbBackupPolicyFixtureRunUpdate(t, state, raw, values)
		mongodbBackupPolicyFixtureAssertSent(t, request, "CrossBackupType", "delete")
		mongodbBackupPolicyFixtureAssertSent(t, request, "SrcRegion", "cn-beijing")
		mongodbBackupPolicyFixtureAssertSent(t, request, "DestRegion", "cn-shanghai")
		mongodbBackupPolicyFixtureAssertOmitted(t, request, "CrossBackupPeriod", "CrossRetentionType", "CrossRetentionValue", "CrossLogRetentionType", "CrossLogRetentionValue", "EnableCrossLogBackup")
	})

	t.Run("removing the regions after delete sends no request", func(t *testing.T) {
		// R-1 transition 2: state and configuration both read "delete" and the regions
		// leave the configuration. Optional+Computed keeps the last values, so nothing
		// changes and no ModifyBackupPolicy request is sent. The documentation states
		// that this removal does not locate a different policy.
		raw := mongodbBackupPolicyFixtureConfig(map[string]interface{}{
			"cross_backup_type": "delete",
		})
		values := map[string]cty.Value{
			"cross_backup_type": cty.StringVal("delete"),
		}
		state := mongodbBackupPolicyFixtureState(map[string]string{
			"cross_backup_type": "delete",
			"src_region":        "cn-beijing",
			"dest_region":       "cn-shanghai",
		})
		mongodbBackupPolicyFixtureAssertNoPolicyChange(t, state, raw, values)
	})

	t.Run("delete mode rejects a missing region on a region-driven delete", func(t *testing.T) {
		// The delete shape needs both regions even when the selector did not change; a
		// missing location must be named before any request leaves the provider.
		raw := mongodbBackupPolicyFixtureConfig(map[string]interface{}{
			"cross_backup_type":       "delete",
			"src_region":              "cn-beijing",
			"enable_cross_log_backup": 0,
		})
		values := map[string]cty.Value{
			"cross_backup_type":       cty.StringVal("delete"),
			"src_region":              cty.StringVal("cn-beijing"),
			"enable_cross_log_backup": cty.NumberIntVal(0),
		}
		// The selector already reads "delete" and only a policy field changes, so the
		// delete branch is selected; the destination region is absent and the policy
		// cannot be located, which must be reported before the request is built.
		state := mongodbBackupPolicyFixtureState(map[string]string{
			"cross_backup_type":       "delete",
			"src_region":              "cn-beijing",
			"enable_cross_log_backup": "1",
		})
		errText, requests := mongodbBackupPolicyFixtureRunUpdateError(t, state, raw, values)
		if !strings.Contains(errText, "dest_region") {
			t.Fatalf("the error must name the missing dest_region, got %s", errText)
		}
		if len(requests) != 0 {
			t.Fatalf("an incomplete delete must not reach the API, got %d requests", len(requests))
		}
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

	t.Run("geo attributes without the operation selector are rejected at plan time", func(t *testing.T) {
		// The selector requirement is structural (RequiredWith on every geo attribute),
		// so the configuration is rejected while planning instead of after apply starts.
		raw := mongodbBackupPolicyFixtureConfig(map[string]interface{}{
			"dest_region":             "cn-shanghai",
			"enable_cross_log_backup": 0,
		})
		err := mongodbBackupPolicyFixtureRunPlanValidation(t, raw)
		if err == nil || !strings.Contains(err.Error(), "cross_backup_type") {
			t.Fatalf("the plan must reject a geo attribute without cross_backup_type, got %v", err)
		}
	})

	t.Run("an explicit zero geo attribute alone still requires the selector", func(t *testing.T) {
		// The only geo input is an explicit 0 and no other geo attribute is present, so
		// nothing except the explicit zero can trigger RequiredWith. The plan must name
		// cross_backup_type, which proves an explicit zero counts as configured instead
		// of being treated as absent (QA-N1). Plan validation runs before Update, so no
		// cloud request can be issued by this case.
		raw := mongodbBackupPolicyFixtureConfig(map[string]interface{}{
			"enable_cross_log_backup": 0,
		})
		err := mongodbBackupPolicyFixtureRunPlanValidation(t, raw)
		if err == nil || !strings.Contains(err.Error(), "cross_backup_type") {
			t.Fatalf("the plan must reject an explicit zero geo attribute without cross_backup_type, got %v", err)
		}
	})

	t.Run("explicit zero counts as configured at plan time", func(t *testing.T) {
		// An explicit 0 must be treated as a configured attribute, not as absent.
		raw := mongodbBackupPolicyFixtureConfig(map[string]interface{}{
			"cross_backup_type":       "delete",
			"enable_cross_log_backup": 0,
		})
		if err := mongodbBackupPolicyFixtureRunPlanValidation(t, raw); err != nil {
			t.Fatalf("an explicit zero with the selector must pass the plan, got %s", err)
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

	t.Run("delete mode without the regions is rejected", func(t *testing.T) {
		raw := mongodbBackupPolicyFixtureConfig(map[string]interface{}{"cross_backup_type": "delete"})
		values := map[string]cty.Value{"cross_backup_type": cty.StringVal("delete")}
		state := mongodbBackupPolicyFixtureState(map[string]string{"cross_backup_type": "update"})
		errText, requests := mongodbBackupPolicyFixtureRunUpdateError(t, state, raw, values)
		if !strings.Contains(errText, "src_region") || !strings.Contains(errText, "dest_region") {
			t.Fatalf("the error must name the missing regions, got %s", errText)
		}
		if len(requests) != 0 {
			t.Fatalf("an incomplete delete must not reach the API, got %d requests", len(requests))
		}
	})

	t.Run("log delay retention without a value is rejected", func(t *testing.T) {
		// The log retention value keeps its own delay-only requirement; the
		// never-without-value test above proves that requirement does not leak into
		// "never" (QA-N5).
		raw, values := mongodbBackupPolicyFixtureCrossConfig(t, 1)
		raw["cross_log_retention_type"] = "delay"
		values["cross_log_retention_type"] = cty.StringVal("delay")
		errText, requests := mongodbBackupPolicyFixtureRunUpdateError(t, mongodbBackupPolicyFixtureState(nil), raw, values)
		if !strings.Contains(errText, "cross_log_retention_value") {
			t.Fatalf("the error must name the missing cross_log_retention_value, got %s", errText)
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

// TestUnitMongoDBInstanceSweeperMatches locks the instance names the sweeper owns. The
// documentation example uses a fixed name, so a failed example run must stay cleanable;
// a user instance that only resembles that name must not be deleted (QA-N7).
func TestUnitMongoDBInstanceSweeperMatches(t *testing.T) {
	cases := []struct {
		description string
		want        bool
	}{
		{"tf-testAccMongoDBInstanceClassicConfig1234", true},
		{"tf_testAccMongoDBInstanceCrossBackup", true},
		{"Tf-TestAccMongoDBInstance", true},
		{"terraform-example", true},
		{" Terraform-Example ", true},
		{"terraform-example-prod", false},
		{"my-terraform-example", false},
		{"tf-other-name", false},
		{"mydb", false},
		{"", false},
	}
	for _, c := range cases {
		if got := mongodbInstanceSweeperMatches(c.description); got != c.want {
			t.Fatalf("mongodbInstanceSweeperMatches(%q) = %v, want %v", c.description, got, c.want)
		}
	}
}
