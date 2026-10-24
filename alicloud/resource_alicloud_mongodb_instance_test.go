package alicloud

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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

// mongodbTestZoneSingleZone returns the zone used by the test that creates a
// single-zone local-disk instance. cn-beijing-h is only open for local disks in
// a multi-zone layout and rejects a single-zone CreateDBInstance with
// Zone.Closed, so it defaults to cn-beijing-l and can be overridden via the
// MONGODB_TEST_ZONE_SINGLE_ZONE environment variable.
func mongodbTestZoneSingleZone() string {
	if v := os.Getenv("MONGODB_TEST_ZONE_SINGLE_ZONE"); v != "" {
		return v
	}
	return "cn-beijing-l"
}

// mongodbTestZoneRestore returns the zone used to restore a released instance
// from its recycle-bin backup set. cn-beijing-h accepts the same single-AZ
// CreateDBInstance but rejects the restore with Zone.Closed, so it defaults to
// cn-beijing-l and can be overridden via the MONGODB_TEST_ZONE_RESTORE
// environment variable.
func mongodbTestZoneRestore() string {
	if v := os.Getenv("MONGODB_TEST_ZONE_RESTORE"); v != "" {
		return v
	}
	return "cn-beijing-l"
}

func TestUnitAliCloudMongoDBInstanceBackupIdSchema(t *testing.T) {
	backupId, ok := resourceAliCloudMongoDBInstance().Schema["backup_id"]
	if !ok {
		t.Fatal("backup_id is not defined in the alicloud_mongodb_instance schema")
	}
	if backupId.Type != schema.TypeString {
		t.Errorf("backup_id type is %v, expected schema.TypeString", backupId.Type)
	}
	if !backupId.Optional {
		t.Error("backup_id must be optional")
	}
	if !backupId.ForceNew {
		t.Error("backup_id must be ForceNew so a different backup set recreates the instance")
	}
	if backupId.Required || backupId.Computed {
		t.Errorf("backup_id must be optional-only, got required=%t computed=%t", backupId.Required, backupId.Computed)
	}
	if len(backupId.RequiredWith) != 1 || backupId.RequiredWith[0] != "src_db_instance_id" {
		t.Errorf("backup_id must require src_db_instance_id, got %v", backupId.RequiredWith)
	}

	// CreateDBInstance rejects BackupId without SrcDBInstanceId, so the provider
	// has to reject backup_id-only configurations itself instead of sending an
	// incomplete request. The constraint sits on backup_id because RequiredWith
	// is only evaluated on the attribute carrying it: putting it on
	// src_db_instance_id would instead forbid point-in-time clones.
	srcInstanceId, ok := resourceAliCloudMongoDBInstance().Schema["src_db_instance_id"]
	if !ok {
		t.Fatal("src_db_instance_id is not defined in the alicloud_mongodb_instance schema")
	}
	if !srcInstanceId.Optional || srcInstanceId.Required {
		t.Errorf("src_db_instance_id must stay optional, got required=%t optional=%t", srcInstanceId.Required, srcInstanceId.Optional)
	}
	if len(srcInstanceId.RequiredWith) != 0 {
		t.Errorf("src_db_instance_id must not be required with anything, got %v", srcInstanceId.RequiredWith)
	}

	required := map[string]interface{}{
		"engine_version":      "7.0",
		"db_instance_class":   "mdb.shard.2x.xlarge.d",
		"db_instance_storage": 20,
	}
	validate := func(extra map[string]interface{}) diag.Diagnostics {
		cfg := map[string]interface{}{}
		for k, v := range required {
			cfg[k] = v
		}
		for k, v := range extra {
			cfg[k] = v
		}
		return resourceAliCloudMongoDBInstance().Validate(sdktf.NewResourceConfigRaw(cfg))
	}

	backupOnly := validate(map[string]interface{}{"backup_id": "backup-test"})
	if len(backupOnly) == 0 {
		t.Error("backup_id without src_db_instance_id must be rejected before CreateDBInstance")
	} else {
		var detail strings.Builder
		for _, d := range backupOnly {
			detail.WriteString(d.Summary)
			detail.WriteString(d.Detail)
		}
		if !strings.Contains(detail.String(), "src_db_instance_id") {
			t.Errorf("backup_id-only rejection must point at src_db_instance_id, got %v", backupOnly)
		}
	}
	if diags := validate(map[string]interface{}{
		"src_db_instance_id": "dds-test",
		"restore_time":       "2022-01-01T12:05:05Z",
	}); len(diags) != 0 {
		t.Errorf("src_db_instance_id + restore_time (point-in-time clone) must stay valid, got %v", diags)
	}
	if diags := validate(map[string]interface{}{
		"src_db_instance_id": "dds-test",
		"backup_id":          "backup-test",
	}); len(diags) != 0 {
		t.Errorf("src_db_instance_id + backup_id must be valid, got %v", diags)
	}
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
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region(mongodbTestRegion())})
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, aliCloudMongoDBInstanceLocalDiskDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region(mongodbTestRegion())})
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
				// MigrateAvailableZone is accepted on the instance state produced by
				// the create step above (7.0, 20GB, default cloud_essd1 disk). It has
				// to run here: the same request was rejected with InsType.NotSupport
				// once the instance had been upgraded to 8.0 and resized to 60GB with
				// the disk converted to cloud_auto.
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
					"engine_version": "8.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version": "8.0",
					}),
				),
			},
			{
				// The class upgrade used to fail when it was issued right after the
				// disk-type conversion: Update returned before the new class was
				// applied. Update now waits for the class attribute, and this step
				// keeps the resize away from the disk conversion that precedes it.
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
					"backup_id":              "backup-tfacc-placeholder",
					"encryption_key":         "key-tfacc-placeholder",
					"encryptor_name":         "aes-256-cbc",
					"instance_charge_type":   "PrePaid",
					"kms_encrypted_password": "kms-tfacc-placeholder",
					// testAccConfig renders map keys unquoted, so a HCL
					// reserved word such as "for" would make the whole
					// configuration unparseable before plan runs.
					"kms_encryption_context": map[string]string{"env": "tfacc"},
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
				// The placeholder step above leaves a config the post-test
				// destroy cannot plan (no db_instance_class), so restore the
				// required arguments and drop the placeholders again.
				Config: testAccConfig(map[string]interface{}{
					"auto_renew":             REMOVEKEY,
					"auto_renew_duration":    REMOVEKEY,
					"backup_id":              REMOVEKEY,
					"encryption_key":         REMOVEKEY,
					"encryptor_name":         REMOVEKEY,
					"instance_charge_type":   REMOVEKEY,
					"kms_encrypted_password": REMOVEKEY,
					"kms_encryption_context": REMOVEKEY,
					"period":                 REMOVEKEY,
					"restore_time":           REMOVEKEY,
					"role_arn":               REMOVEKEY,
					"src_db_instance_id":     REMOVEKEY,
					"tde_status":             REMOVEKEY,
					"db_instance_class":      "mdb.shard.2x.2xlarge.d",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_class": "mdb.shard.2x.2xlarge.d",
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
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region(mongodbTestRegion())})
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
	return aliCloudMongoDBInstanceBasicDependence1InZone(name, mongodbTestZonePrimary())
}

func aliCloudMongoDBInstanceBasicDependence1InZone(name, zone string) string {
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
`, name, zone)
}

func aliCloudMongoDBInstanceLocalDiskDependence(name string) string {
	return aliCloudMongoDBInstanceBasicDependence1InZone(name, mongodbTestZoneSingleZone())
}

const (
	mongodbRecycleTestStorage     = 80
	mongodbRecycleTestStorageType = "cloud_essd1"
)

type mongodbRecycleTestSpec struct {
	zoneId        string
	engineVersion string
	instanceClass string
}

// testAccMongoDBInstanceRecycleSpec picks an ESSD cloud disk replica set
// inventory that is actually available in the test region, restricted to the
// zone that accepts restoring a released instance from its recycle-bin backup
// set.
func testAccMongoDBInstanceRecycleSpec(client *connectivity.AliyunClient) (mongodbRecycleTestSpec, error) {
	restoreZone := mongodbTestZoneRestore()
	candidates := []mongodbRecycleTestSpec{
		{engineVersion: "7.0", instanceClass: "mdb.shard.4x.large.d"},
		{engineVersion: "8.0", instanceClass: "mdb.shard.4x.large.d"},
		{engineVersion: "7.0", instanceClass: "mdb.shard.2x.xlarge.d"},
		{engineVersion: "8.0", instanceClass: "mdb.shard.2x.xlarge.d"},
	}
	for _, candidate := range candidates {
		response, err := client.RpcPost("Dds", "2015-12-01", "DescribeAvailabilityZones", nil, map[string]interface{}{
			"RegionId":           client.RegionId,
			"DbType":             "normal",
			"InstanceChargeType": "PostPaid",
			"StorageSupport":     "cloud",
			"StorageType":        mongodbRecycleTestStorageType,
			"ReplicationFactor":  "3",
			"EngineVersion":      candidate.engineVersion,
			"DBInstanceClass":    candidate.instanceClass,
		}, true)
		if err != nil {
			return mongodbRecycleTestSpec{}, WrapError(err)
		}
		rawZones, err := jsonpath.Get("$.AvailableZones", response)
		if err != nil {
			return mongodbRecycleTestSpec{}, WrapError(err)
		}
		zones, ok := rawZones.([]interface{})
		if !ok {
			continue
		}
		for _, rawZone := range zones {
			zone, ok := rawZone.(map[string]interface{})
			if !ok {
				continue
			}
			zoneId := fmt.Sprint(zone["ZoneId"])
			if zoneId != restoreZone || strings.Contains(zoneId, MULTI_IZ_SYMBOL) {
				continue
			}
			candidate.zoneId = zoneId
			return candidate, nil
		}
	}
	return mongodbRecycleTestSpec{}, fmt.Errorf("no cloud disk replica set inventory is available in zone %s of region %s for classes mdb.shard.4x.large.d and mdb.shard.2x.xlarge.d", restoreZone, client.RegionId)
}

// mongodbRecycleNetworkOwner records the network resources created for the
// recycle bin restore test. Its cleanup is registered before the first
// CreateVpc call, so each ID is recorded as soon as it exists: a failure after
// that point still has to let the cleanup delete what was already created.
type mongodbRecycleNetworkOwner struct {
	VpcId     string
	VSwitchId string
}

// mongodbRecycleNetworkApi is the slice of the VPC API used while preparing the
// network. It is an interface so the partial-failure paths can be covered
// without creating cloud resources.
type mongodbRecycleNetworkApi interface {
	CreateVpc(name string) (string, error)
	WaitVpcAvailable(vpcId string) error
	CreateVSwitch(name, vpcId, zoneId string) (string, error)
	WaitVSwitchAvailable(vswitchId string) error
}

func (owner *mongodbRecycleNetworkOwner) prepare(api mongodbRecycleNetworkApi, name, zoneId string) error {
	vpcId, err := api.CreateVpc(name)
	if err != nil {
		return WrapError(err)
	}
	owner.VpcId = vpcId
	if err := api.WaitVpcAvailable(vpcId); err != nil {
		return WrapError(err)
	}

	vswitchId, err := api.CreateVSwitch(name, vpcId, zoneId)
	if err != nil {
		return WrapError(err)
	}
	owner.VSwitchId = vswitchId
	if err := api.WaitVSwitchAvailable(vswitchId); err != nil {
		return WrapError(err)
	}
	return nil
}

type mongodbRecycleNetworkClient struct {
	client *connectivity.AliyunClient
}

func (api mongodbRecycleNetworkClient) CreateVpc(name string) (string, error) {
	response, err := api.client.RpcPost("Vpc", "2016-04-28", "CreateVpc", nil, map[string]interface{}{
		"RegionId":    api.client.RegionId,
		"ClientToken": buildClientToken("CreateVpc"),
		"CidrBlock":   "172.16.0.0/16",
		"VpcName":     name,
	}, true)
	if err != nil {
		return "", err
	}
	vpcId := fmt.Sprint(response["VpcId"])
	if vpcId == "" || vpcId == "<nil>" {
		return "", fmt.Errorf("CreateVpc for %s returned no VpcId", name)
	}
	return vpcId, nil
}

func (api mongodbRecycleNetworkClient) WaitVpcAvailable(vpcId string) error {
	vpcStateConf := BuildStateConf([]string{"Pending"}, []string{"Available"}, 5*time.Minute, 5*time.Second, func() (interface{}, string, error) {
		object, err := (&VpcService{api.client}).DescribeVpc(vpcId)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		return object, fmt.Sprint(object["Status"]), nil
	})
	_, err := vpcStateConf.WaitForState()
	return err
}

func (api mongodbRecycleNetworkClient) CreateVSwitch(name, vpcId, zoneId string) (string, error) {
	response, err := api.client.RpcPost("Vpc", "2016-04-28", "CreateVSwitch", nil, map[string]interface{}{
		"RegionId":    api.client.RegionId,
		"ClientToken": buildClientToken("CreateVSwitch"),
		"CidrBlock":   "172.16.0.0/24",
		"VpcId":       vpcId,
		"VSwitchName": name,
		"ZoneId":      zoneId,
	}, true)
	if err != nil {
		return "", err
	}
	vswitchId := fmt.Sprint(response["VSwitchId"])
	if vswitchId == "" || vswitchId == "<nil>" {
		return "", fmt.Errorf("CreateVSwitch for %s returned no VSwitchId", name)
	}
	return vswitchId, nil
}

func (api mongodbRecycleNetworkClient) WaitVSwitchAvailable(vswitchId string) error {
	vswitchStateConf := BuildStateConf([]string{"Pending"}, []string{"Available"}, 5*time.Minute, 5*time.Second, func() (interface{}, string, error) {
		object, err := (&VpcService{api.client}).DescribeVSwitchWithTeadsl(vswitchId)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		return object, fmt.Sprint(object["Status"]), nil
	})
	_, err := vswitchStateConf.WaitForState()
	return err
}

// TestUnitMongoDBInstanceRecycleNetworkCleanupOwnership covers the controlled
// failures of the network preparation: whatever was created before the failure
// must be recorded on the owner so the cleanup that the acceptance test
// registers ahead of the preparation can delete it.
func TestUnitMongoDBInstanceRecycleNetworkCleanupOwnership(t *testing.T) {
	testCases := []struct {
		name            string
		failOn          string
		expectedVpc     string
		expectedVSwitch string
	}{
		{name: "create vpc fails", failOn: "CreateVpc"},
		{name: "wait vpc available fails", failOn: "WaitVpcAvailable", expectedVpc: "vpc-unit"},
		{name: "create vswitch fails", failOn: "CreateVSwitch", expectedVpc: "vpc-unit"},
		{name: "wait vswitch available fails", failOn: "WaitVSwitchAvailable", expectedVpc: "vpc-unit", expectedVSwitch: "vsw-unit"},
		{name: "succeeds", expectedVpc: "vpc-unit", expectedVSwitch: "vsw-unit"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			owner := &mongodbRecycleNetworkOwner{}
			err := owner.prepare(&fakeMongoDBRecycleNetworkApi{failOn: testCase.failOn}, "tf-unit", "cn-hangzhou-h")
			if testCase.failOn != "" && err == nil {
				t.Fatalf("expected the preparation to fail on %s", testCase.failOn)
			}
			if testCase.failOn == "" && err != nil {
				t.Fatalf("expected the preparation to succeed, got %s", err)
			}
			if owner.VpcId != testCase.expectedVpc {
				t.Errorf("owner.VpcId is %q, expected %q", owner.VpcId, testCase.expectedVpc)
			}
			if owner.VSwitchId != testCase.expectedVSwitch {
				t.Errorf("owner.VSwitchId is %q, expected %q", owner.VSwitchId, testCase.expectedVSwitch)
			}
		})
	}
}

type fakeMongoDBRecycleNetworkApi struct {
	failOn string
}

func (api *fakeMongoDBRecycleNetworkApi) fail(step string) error {
	if api.failOn == step {
		return fmt.Errorf("injected %s failure", step)
	}
	return nil
}

func (api *fakeMongoDBRecycleNetworkApi) CreateVpc(string) (string, error) {
	if err := api.fail("CreateVpc"); err != nil {
		return "", err
	}
	return "vpc-unit", nil
}

func (api *fakeMongoDBRecycleNetworkApi) WaitVpcAvailable(string) error {
	return api.fail("WaitVpcAvailable")
}

func (api *fakeMongoDBRecycleNetworkApi) CreateVSwitch(string, string, string) (string, error) {
	if err := api.fail("CreateVSwitch"); err != nil {
		return "", err
	}
	return "vsw-unit", nil
}

func (api *fakeMongoDBRecycleNetworkApi) WaitVSwitchAvailable(string) error {
	return api.fail("WaitVSwitchAvailable")
}

func testAccMongoDBInstanceRecycleSourceInstance(client *connectivity.AliyunClient, name, vpcId, vswitchId string, spec mongodbRecycleTestSpec) (string, error) {
	response, err := client.RpcPost("Dds", "2015-12-01", "CreateDBInstance", nil, map[string]interface{}{
		"RegionId":              client.RegionId,
		"ClientToken":           buildClientToken("CreateDBInstance"),
		"Engine":                "MongoDB",
		"EngineVersion":         spec.engineVersion,
		"DBInstanceClass":       spec.instanceClass,
		"DBInstanceStorage":     mongodbRecycleTestStorage,
		"StorageType":           mongodbRecycleTestStorageType,
		"VpcId":                 vpcId,
		"VSwitchId":             vswitchId,
		"ZoneId":                spec.zoneId,
		"NetworkType":           "VPC",
		"ChargeType":            "PostPaid",
		"DBInstanceDescription": name,
		"SecurityIPList":        "127.0.0.1",
		"ReplicationFactor":     "3",
	}, true)
	if err != nil {
		return "", WrapError(err)
	}
	instanceId := fmt.Sprint(response["DBInstanceId"])
	if instanceId == "" || instanceId == "<nil>" {
		return "", fmt.Errorf("CreateDBInstance for %s returned no DBInstanceId", name)
	}
	service := &MongoDBService{client}
	stateConf := BuildStateConf([]string{"Creating", "Initializing", "Restoring"}, []string{"Running"}, 60*time.Minute, 10*time.Second, service.RdsMongodbDBInstanceStateRefreshFunc(instanceId, []string{"Failed", "DeleteFailed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return instanceId, WrapError(err)
	}
	return instanceId, nil
}

// testAccMongoDBInstanceRecycleRetention keeps the backup set that the service
// creates when the instance is released (retention policy 1) so the recycle bin
// can rebuild data from it.
func testAccMongoDBInstanceRecycleRetention(client *connectivity.AliyunClient, instanceId string) error {
	if _, err := client.RpcPost("Dds", "2015-12-01", "ModifyBackupPolicy", nil, map[string]interface{}{
		"RegionId":                               client.RegionId,
		"DBInstanceId":                           instanceId,
		"PreferredBackupTime":                    "02:00Z-03:00Z",
		"PreferredBackupPeriod":                  "Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday",
		"BackupRetentionPolicyOnClusterDeletion": "1",
	}, true); err != nil {
		return WrapError(err)
	}
	response, err := client.RpcPost("Dds", "2015-12-01", "DescribeBackupPolicy", nil, map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": instanceId,
	}, true)
	if err != nil {
		return WrapError(err)
	}
	if policy := fmt.Sprint(response["BackupRetentionPolicyOnClusterDeletion"]); policy != "1" {
		return fmt.Errorf("the backup retention policy of instance %s is %s, expected 1", instanceId, policy)
	}
	return nil
}

// testAccMongoDBInstanceRecycleBinBackup returns the backup set that the recycle
// bin exposes for the released instance, which is the backup point the restore
// has to use.
func testAccMongoDBInstanceRecycleBinBackup(client *connectivity.AliyunClient, instanceId string) (string, error) {
	deadline := time.Now().Add(30 * time.Minute)
	for {
		response, err := client.RpcPost("Dds", "2015-12-01", "DescribeBackups", nil, map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": instanceId,
			"SrcRegion":    client.RegionId,
			"PageSize":     PageSizeLarge,
			"PageNumber":   1,
		}, true)
		if err == nil {
			latestBackupId := ""
			latestEndTime := ""
			if rawBackups, e := jsonpath.Get("$.Backups.Backup", response); e == nil {
				if backups, ok := rawBackups.([]interface{}); ok {
					for _, rawBackup := range backups {
						backup, ok := rawBackup.(map[string]interface{})
						if !ok || fmt.Sprint(backup["BackupStatus"]) != "Success" {
							continue
						}
						backupId := fmt.Sprint(backup["BackupId"])
						if backupId == "" || backupId == "<nil>" {
							continue
						}
						if endTime := fmt.Sprint(backup["BackupEndTime"]); endTime >= latestEndTime {
							latestEndTime = endTime
							latestBackupId = backupId
						}
					}
				}
			}
			if latestBackupId != "" {
				return latestBackupId, nil
			}
		} else {
			log.Printf("[WARN] DescribeBackups of the released instance %s failed: %s", instanceId, err)
		}
		if time.Now().After(deadline) {
			return "", fmt.Errorf("no successful backup set of the released instance %s is visible in the recycle bin within 30 minutes", instanceId)
		}
		time.Sleep(15 * time.Second)
	}
}

func testAccMongoDBInstanceRecycleRelease(client *connectivity.AliyunClient, instanceId string) error {
	if _, err := client.RpcPost("Dds", "2015-12-01", "DeleteDBInstance", nil, map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": instanceId,
		"ClientToken":  buildClientToken("DeleteDBInstance"),
	}, true); err != nil && !IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
		return WrapError(err)
	}
	service := &MongoDBService{client}
	stateConf := BuildStateConf([]string{"Creating", "Running", "Deleting", "Restoring", "Initializing", "Maintaining"}, []string{}, 45*time.Minute, 15*time.Second, service.RdsMongodbDBInstanceStateRefreshFunc(instanceId, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}
	return nil
}

func testAccMongoDBInstanceRecycleCleanup(client *connectivity.AliyunClient, instanceDescription string, instanceIds []string, vswitchId, vpcId string) {
	for i := 0; i < 60; i++ {
		response, err := client.RpcPost("Dds", "2015-12-01", "DescribeDBInstances", nil, map[string]interface{}{
			"RegionId":              client.RegionId,
			"DBInstanceDescription": instanceDescription,
			"PageSize":              PageSizeLarge,
			"PageNumber":            1,
		}, true)
		if err != nil {
			log.Printf("[WARN] cleanup: DescribeDBInstances for %s failed: %s", instanceDescription, err)
			break
		}
		rawInstances, err := jsonpath.Get("$.DBInstances.DBInstance", response)
		if err != nil {
			break
		}
		instances, ok := rawInstances.([]interface{})
		if !ok || len(instances) == 0 {
			break
		}
		for _, rawInstance := range instances {
			instance, ok := rawInstance.(map[string]interface{})
			if !ok {
				continue
			}
			instanceId := fmt.Sprint(instance["DBInstanceId"])
			if _, err := client.RpcPost("Dds", "2015-12-01", "DeleteDBInstance", nil, map[string]interface{}{
				"RegionId":     client.RegionId,
				"DBInstanceId": instanceId,
				"ClientToken":  buildClientToken("DeleteDBInstance"),
			}, true); err != nil && !IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
				log.Printf("[WARN] cleanup: DeleteDBInstance %s failed: %s", instanceId, err)
			}
		}
		time.Sleep(15 * time.Second)
	}

	// The released instances keep their recycle bin backups after the test, so
	// delete them explicitly instead of leaving a long-term backup behind.
	for _, instanceId := range instanceIds {
		if instanceId == "" || instanceId == "<nil>" {
			continue
		}
		response, err := client.RpcPost("Dds", "2015-12-01", "DescribeBackups", nil, map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": instanceId,
			"SrcRegion":    client.RegionId,
			"PageSize":     PageSizeLarge,
			"PageNumber":   1,
		}, true)
		if err != nil {
			log.Printf("[WARN] cleanup: DescribeBackups for the released instance %s failed: %s", instanceId, err)
			continue
		}
		rawBackups, err := jsonpath.Get("$.Backups.Backup", response)
		if err != nil {
			continue
		}
		backups, ok := rawBackups.([]interface{})
		if !ok {
			continue
		}
		for _, rawBackup := range backups {
			backup, ok := rawBackup.(map[string]interface{})
			if !ok {
				continue
			}
			backupId := fmt.Sprint(backup["BackupId"])
			if backupId == "" || backupId == "<nil>" {
				continue
			}
			if _, err := client.RpcPost("Dds", "2015-12-01", "DeleteBackup", nil, map[string]interface{}{
				"RegionId":     client.RegionId,
				"DBInstanceId": instanceId,
				"BackupId":     backupId,
			}, false); err != nil {
				log.Printf("[WARN] cleanup: DeleteBackup %s of the released instance %s failed: %s", backupId, instanceId, err)
			}
		}
	}

	if vswitchId != "" {
		for i := 0; i < 40; i++ {
			_, err := client.RpcPost("Vpc", "2016-04-28", "DeleteVSwitch", nil, map[string]interface{}{
				"RegionId":  client.RegionId,
				"VSwitchId": vswitchId,
			}, false)
			if err == nil || IsExpectedErrors(err, []string{"InvalidVswitchID.NotFound", "InvalidVSwitchId.NotFound", "Forbidden.VSwitchNotFound"}) {
				break
			}
			if NeedRetry(err) || IsExpectedErrors(err, []string{"DependencyViolation", "OperationConflict", "IncorrectVSwitchStatus"}) {
				time.Sleep(15 * time.Second)
				continue
			}
			log.Printf("[WARN] cleanup: DeleteVSwitch %s failed: %s", vswitchId, err)
			break
		}
	}
	if vpcId != "" {
		for i := 0; i < 40; i++ {
			_, err := client.RpcPost("Vpc", "2016-04-28", "DeleteVpc", nil, map[string]interface{}{
				"RegionId": client.RegionId,
				"VpcId":    vpcId,
			}, false)
			if err == nil || IsExpectedErrors(err, []string{"InvalidVpcID.NotFound", "Forbidden.VpcNotFound"}) {
				break
			}
			if NeedRetry(err) || IsExpectedErrors(err, []string{"DependencyViolation.VSwitch", "DependencyViolation.SecurityGroup", "IncorrectVpcStatus"}) {
				time.Sleep(15 * time.Second)
				continue
			}
			log.Printf("[WARN] cleanup: DeleteVpc %s failed: %s", vpcId, err)
			break
		}
	}
}

func testAccMongoDBInstanceRecycleRestoreConfig(name string, spec mongodbRecycleTestSpec, vswitchId, sourceInstanceId, backupId string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_mongodb_instance" "default" {
  name                 = var.name
  engine_version       = "%s"
  db_instance_class    = "%s"
  db_instance_storage  = %d
  storage_type         = "%s"
  instance_charge_type = "PostPaid"
  replication_factor   = 3
  vswitch_id           = "%s"
  src_db_instance_id   = "%s"
  backup_id            = "%s"
}
`, name, spec.engineVersion, spec.instanceClass, mongodbRecycleTestStorage, mongodbRecycleTestStorageType, vswitchId, sourceInstanceId, backupId)
}

// TestAccAliCloudMongoDBInstance_backupIdRecycleBinRestore restores a released
// ESSD cloud disk replica set instance from its recycle-bin backup set with
// src_db_instance_id + backup_id. The precondition (source instance, backup set
// and release) is built with the API because the HCL steps cannot create a
// backup set dynamically nor release the source between steps.
func TestAccAliCloudMongoDBInstance_backupIdRecycleBinRestore(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_instance.default"
	testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Region(mongodbTestRegion())})
	rawClient, err := sharedClientForRegion(defaultRegionToTest)
	if err != nil {
		t.Fatalf("failed to build the client for region %s: %s", defaultRegionToTest, err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccMongoDBInstanceRecycle%d", rand)

	spec, err := testAccMongoDBInstanceRecycleSpec(client)
	if err != nil {
		t.Fatalf("failed to discover an ESSD cloud disk replica set inventory: %s", err)
	}
	log.Printf("[INFO] recycle bin restore test uses region %s, zone %s, engine %s, class %s", client.RegionId, spec.zoneId, spec.engineVersion, spec.instanceClass)

	// The cleanup has to be registered before the first CreateVpc call:
	// preparing the network can fail after the VPC or the vSwitch already
	// exists, and this test creates both of them itself.
	network := &mongodbRecycleNetworkOwner{}
	createdInstanceIds := make([]string, 0, 2)
	t.Cleanup(func() {
		testAccMongoDBInstanceRecycleCleanup(client, name, createdInstanceIds, network.VSwitchId, network.VpcId)
	})
	if err := network.prepare(mongodbRecycleNetworkClient{client: client}, name, spec.zoneId); err != nil {
		t.Fatalf("failed to prepare the VPC network: %s", err)
	}
	vpcId, vswitchId := network.VpcId, network.VSwitchId

	sourceInstanceId, err := testAccMongoDBInstanceRecycleSourceInstance(client, name, vpcId, vswitchId, spec)
	if err != nil {
		t.Fatalf("failed to create the source instance: %s", err)
	}
	createdInstanceIds = append(createdInstanceIds, sourceInstanceId)
	if err := testAccMongoDBInstanceRecycleRetention(client, sourceInstanceId); err != nil {
		t.Fatalf("failed to set the backup retention policy of the source instance: %s", err)
	}
	if err := testAccMongoDBInstanceRecycleRelease(client, sourceInstanceId); err != nil {
		t.Fatalf("failed to release the source instance %s: %s", sourceInstanceId, err)
	}
	backupId, err := testAccMongoDBInstanceRecycleBinBackup(client, sourceInstanceId)
	if err != nil {
		t.Fatalf("the recycle bin backup set is not available: %s", err)
	}

	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBInstanceRecycleRestoreConfig(name, spec, vswitchId, sourceInstanceId, backupId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_id":           backupId,
						"src_db_instance_id":  sourceInstanceId,
						"db_instance_class":   spec.instanceClass,
						"db_instance_storage": fmt.Sprint(mongodbRecycleTestStorage),
						"engine_version":      spec.engineVersion,
						"storage_type":        mongodbRecycleTestStorageType,
					}),
					func(state *terraform.State) error {
						resourceState, ok := state.RootModule().Resources[resourceId]
						if !ok {
							return fmt.Errorf("%s is not found in the state", resourceId)
						}
						createdInstanceIds = append(createdInstanceIds, resourceState.Primary.ID)
						return nil
					},
				),
			},
		},
	})
}
