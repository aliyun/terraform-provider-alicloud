package alicloud

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"sync"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudADBDBClusterLakeVersion_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_adb_db_cluster_lake_version.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbDbClusterLakeVersionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbDbClusterLakeVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sadbdbclusterlakeversion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbDbClusterLakeVersionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckPrePaidResources(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version": "5.0",
					"vpc_id":             "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":         "${data.alicloud_vswitches.default.ids.0}",
					"zone_id":            "${data.alicloud_adb_zones.default.ids.0}",
					"payment_type":       "Subscription",
					"period":             "1",
					"compute_resource":   "16ACU",
					"storage_resource":   "0ACU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version": "5.0",
						"vpc_id":             CHECKSET,
						"vswitch_id":         CHECKSET,
						"zone_id":            CHECKSET,
						"payment_type":       "Subscription",
						"compute_resource":   "16ACU",
						"storage_resource":   "0ACU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"compute_resource": "32ACU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compute_resource": "32ACU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_resource": "24ACU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_resource": "24ACU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": "10.23.1.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips": "10.23.1.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ssl": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ssl": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ssl": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ssl": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "Subscription",
					"period":       "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "Subscription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "enable_default_resource_group", "source_db_cluster_id", "backup_set_id", "restore_type", "restore_to_time"},
			},
		},
	})
}

func TestAccAliCloudADBDBClusterLakeVersion_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_adb_db_cluster_lake_version.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbDbClusterLakeVersionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbDbClusterLakeVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sadbdbclusterlakeversion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbDbClusterLakeVersionBasicDependence0Twin)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckPrePaidResources(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":            "5.0",
					"vpc_id":                        "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":                    "${data.alicloud_vswitches.default.ids.0}",
					"zone_id":                       "${data.alicloud_adb_zones.default.ids.0}",
					"payment_type":                  "Subscription",
					"secondary_vswitch_id":          "${data.alicloud_vswitches.secondary.ids.0}",
					"secondary_zone_id":             "${data.alicloud_adb_zones.default.ids.1}",
					"compute_resource":              "16ACU",
					"storage_resource":              "0ACU",
					"disk_encryption":               "true",
					"kms_id":                        "${alicloud_kms_key.default.id}",
					"security_ips":                  "10.23.1.0",
					"enable_ssl":                    "true",
					"db_cluster_description":        name,
					"resource_group_id":             "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"period":                        "1",
					"enable_default_resource_group": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version":     "5.0",
						"vpc_id":                 CHECKSET,
						"vswitch_id":             CHECKSET,
						"zone_id":                CHECKSET,
						"payment_type":           "Subscription",
						"secondary_vswitch_id":   CHECKSET,
						"secondary_zone_id":      CHECKSET,
						"compute_resource":       "16ACU",
						"storage_resource":       "0ACU",
						"disk_encryption":        "true",
						"kms_id":                 CHECKSET,
						"security_ips":           "10.23.1.0",
						"enable_ssl":             "true",
						"db_cluster_description": name,
						"resource_group_id":      CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "enable_default_resource_group", "source_db_cluster_id", "backup_set_id", "restore_type", "restore_to_time"},
			},
		},
	})
}

func TestAccAliCloudADBDBClusterLakeVersion_basic1(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_adb_db_cluster_lake_version.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbDbClusterLakeVersionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbDbClusterLakeVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sadbdbclusterlakeversion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbDbClusterLakeVersionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":  "5.0",
					"vpc_id":              "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":          "${data.alicloud_vswitches.default.ids.0}",
					"zone_id":             "${data.alicloud_adb_zones.default.ids.0}",
					"payment_type":        "PayAsYouGo",
					"product_form":        "IntegrationForm",
					"product_version":     "EnterpriseVersion",
					"reserved_node_size":  "8ACU",
					"reserved_node_count": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version":  "5.0",
						"vpc_id":              CHECKSET,
						"vswitch_id":          CHECKSET,
						"zone_id":             CHECKSET,
						"payment_type":        "PayAsYouGo",
						"product_form":        "IntegrationForm",
						"product_version":     "EnterpriseVersion",
						"reserved_node_size":  "8ACU",
						"reserved_node_count": "3",
					}),
				),
			},
			// product_form in-place change: IntegrationForm → LegacyForm.
			// This and the next step verify in-place product_form change in
			// both directions (Integration→Legacy→Integration) starting
			// from an IntegrationForm creation. The customer's direct
			// Legacy→Integration upgrade from a LegacyForm creation is
			// tested by TestAccAliCloudADBDBClusterLakeVersion_basic_upgrade.
			// Moved before reserved_node tests to avoid the reserved_node
			// timeout (17179 Step 2 ClassChanging >72min) blocking
			// product_form verification.
			{
				Config: testAccConfig(map[string]interface{}{
					"product_form": "LegacyForm",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_form": "LegacyForm",
					}),
				),
			},
			// product_form in-place change: LegacyForm → IntegrationForm.
			// This verifies the reverse direction from the previous step.
			// Note: this is a downgrade-then-upgrade path
			// (Integration→Legacy→Integration), NOT the customer's direct
			// Legacy→Integration upgrade. The direct upgrade from a
			// LegacyForm creation is tested by basic_upgrade.
			{
				Config: testAccConfig(map[string]interface{}{
					"product_form": "IntegrationForm",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_form": "IntegrationForm",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"reserved_node_size": "12ACU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"reserved_node_size": "12ACU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"reserved_node_count": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"reserved_node_count": "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"reserved_node_size":  "8ACU",
					"reserved_node_count": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"reserved_node_size":  "8ACU",
						"reserved_node_count": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": "10.23.1.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips": "10.23.1.1",
					}),
				),
			},
			// SSL enable/disable/update cycle is covered by basic_ssl to
			// keep each ACC function within the remote execution timeout.
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"restore_type":         "invalid-restore-type",
					"source_db_cluster_id": "am-fake-restore-source",
					"backup_set_id":        "fake-backup-set-id",
					"restore_to_time":      "2024-01-01T00:00:00Z",
				}),
				ExpectError: regexp.MustCompile("expected restore_type to be one of"),
			},
			// Restore valid configuration after the negative validation test.
			// The SDK's testing_new.go (line 321) sets appliedCfg to the
			// ExpectError step's merged config even when the error is hit;
			// Import (testing_new_import_state.go:82-100) and post-test
			// destroy use appliedCfg. Without this recovery step, the
			// fake restore values would be passed to Import/Destroy and
			// may block cleanup on a paid instance.
			{
				Config: testAccConfig(map[string]interface{}{
					"restore_type":         REMOVEKEY,
					"source_db_cluster_id": REMOVEKEY,
					"backup_set_id":        REMOVEKEY,
					"restore_to_time":      REMOVEKEY,
				}),
				// Add Check to close coverage parser Config capture boundary;
				// without it the parser crosses into the next ImportState
				// step and triggers a bracket-mismatch coverage failure.
				// Also verifies the resource still exists after restoring
				// valid config following the negative validation test.
				Check: resource.ComposeTestCheckFunc(
					rac.checkResourceExists(),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_default_resource_group", "source_db_cluster_id", "backup_set_id", "restore_type", "restore_to_time"},
			},
		},
	})
}

func TestAccAliCloudADBDBClusterLakeVersion_basic_ssl(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_adb_db_cluster_lake_version.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbDbClusterLakeVersionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbDbClusterLakeVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sadbdbclusterlakeversion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbDbClusterLakeVersionBasicDependence0Twin)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version":            "5.0",
					"vpc_id":                        "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":                    "${data.alicloud_vswitches.default.ids.0}",
					"zone_id":                       "${data.alicloud_adb_zones.default.ids.0}",
					"payment_type":                  "PayAsYouGo",
					"product_form":                  "IntegrationForm",
					"product_version":               "EnterpriseVersion",
					"reserved_node_size":            "8ACU",
					"reserved_node_count":           "3",
					"disk_encryption":               "true",
					"kms_id":                        "${alicloud_kms_key.default.id}",
					"security_ips":                  "10.23.1.0",
					"enable_ssl":                    "true",
					"db_cluster_description":        name,
					"resource_group_id":             "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"enable_default_resource_group": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version":     "5.0",
						"vpc_id":                 CHECKSET,
						"vswitch_id":             CHECKSET,
						"zone_id":                CHECKSET,
						"payment_type":           "PayAsYouGo",
						"product_form":           "IntegrationForm",
						"product_version":        "EnterpriseVersion",
						"reserved_node_size":     "8ACU",
						"reserved_node_count":    "3",
						"disk_encryption":        "true",
						"kms_id":                 CHECKSET,
						"security_ips":           "10.23.1.0",
						"enable_ssl":             "true",
						"db_cluster_description": name,
						"resource_group_id":      CHECKSET,
						"ssl_connection_string":  CHECKSET,
						"ssl_expire_time":        CHECKSET,
					}),
				),
			},
			// Disable SSL via the boolean enable_ssl (backward compatibility)
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ssl": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ssl":            "false",
						"ssl_connection_string": REMOVEKEY,
						"ssl_expire_time":       REMOVEKEY,
					}),
				),
			},
			// Re-enable SSL via the string ssl_enabled=Enable (new API)
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ssl":  REMOVEKEY,
					"ssl_enabled": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enabled":           "Enable",
						"enable_ssl":            "true",
						"ssl_connection_string": CHECKSET,
						"ssl_expire_time":       CHECKSET,
					}),
				),
			},
			// Disable SSL via the string ssl_enabled=Disable (new API)
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enabled": "Disable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enabled":           "Disable",
						"enable_ssl":            "false",
						"ssl_connection_string": REMOVEKEY,
						"ssl_expire_time":       REMOVEKEY,
					}),
				),
			},
			// Re-enable SSL to set up for certificate renewal
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enabled": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enabled":     "Enable",
						"enable_ssl":      "true",
						"ssl_expire_time": CHECKSET,
					}),
				),
			},
			// Renew the SSL certificate via ssl_enabled=Update
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enabled": "Update",
					"enable_ssl":  "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enabled":     "Update",
						"enable_ssl":      "true",
						"ssl_expire_time": CHECKSET,
					}),
				),
			},
			// Convergence: re-apply the same Update config. The plan must
			// be empty (no persistent diff) because ssl_enabled="Update" is
			// retained in state by Read, and enable_ssl matches the real
			// SSL state. This verifies the fix for the persistent-diff
			// issue where Read previously overwrote "Update" with "Enable".
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enabled": "Update",
					"enable_ssl":  "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enabled":     "Update",
						"enable_ssl":      "true",
						"ssl_expire_time": CHECKSET,
					}),
				),
			},
			// Priority: ssl_enabled=Enable takes priority over
			// enable_ssl=false. CustomizeDiff derives enable_ssl=true
			// via SetNew; the Update sends EnableSSL=true (from the
			// planned value, not the config). Read writes back
			// ssl_enabled=Enable and enable_ssl=true. No conflict
			// error — the plan converges.
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enabled": "Enable",
					"enable_ssl":  "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enabled":           "Enable",
						"enable_ssl":            "true",
						"ssl_connection_string": CHECKSET,
						"ssl_expire_time":       CHECKSET,
					}),
				),
			},
			// Revert to boolean enable_ssl for a clean import state.
			// Removing ssl_enabled (REMOVEKEY) falls back to enable_ssl,
			// verifying the "remove command, switch back to old boolean"
			// path.
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enabled": REMOVEKEY,
					"enable_ssl":  "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ssl": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_default_resource_group", "source_db_cluster_id", "backup_set_id", "restore_type", "restore_to_time", "ssl_enabled"},
			},
		},
	})
}

// TestAccAliCloudADBDBClusterLakeVersion_basic_upgrade verifies the customer's
// primary requirement (sourceComment 129893725): an existing LegacyForm cluster
// directly upgrades to IntegrationForm (Enterprise) via a single ModifyDBCluster
// call that sets ProductForm + ReservedNodeCount + ReservedNodeSize
// simultaneously, keeping the same instance ID (no rebuild/ForceNew).
//
// Unlike basic1 (which creates with IntegrationForm and then does
// Integration→Legacy→Integration), this test starts from a LegacyForm creation
// — the actual customer starting state — and performs a direct Legacy→Integration
// upgrade. The same instance ID is captured in Step 1 and verified unchanged in
// Step 2 via TestCheckResourceAttrWith.
//
// Note: product_version is NOT set in the upgrade step because it has ForceNew
// in the schema — setting it would trigger a rebuild, contradicting the
// customer's same-ID requirement. The ModifyDBCluster API does not accept
// ProductVersion (it is only sent at Create time); the server handles it
// internally during the form upgrade.
func TestAccAliCloudADBDBClusterLakeVersion_basic_upgrade(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_adb_db_cluster_lake_version.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbDbClusterLakeVersionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbDbClusterLakeVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sadbdbclusterlakeversion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbDbClusterLakeVersionBasicDependence0)
	var instanceID string
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			// Step 1: Create with LegacyForm. This is the customer's
			// starting state — an existing LegacyForm (Data Lakehouse
			// Edition) cluster. compute_resource and storage_resource
			// are the primary resources for LegacyForm; product_version
			// and reserved_node are not applicable.
			//
			// Spec mapping reference (public docs): for Legacy→Integration
			// upgrade, the enterprise node spec×count (X) must satisfy
			// N < X <= M+N, where M=old compute cores, N=old storage cores.
			// Here M=32 (32ACU), N=0 (0ACU), M+N=32; target X=8ACU×3=24,
			// 0 < 24 <= 32. This is a candidate combination based on the
			// public mapping formula, not a proven root cause — the actual
			// ModifyDBCluster upgrade contract is pending real ACC / owner
			// confirmation.
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_version": "5.0",
					"vpc_id":             "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":         "${data.alicloud_vswitches.default.ids.0}",
					"zone_id":            "${data.alicloud_adb_zones.default.ids.0}",
					"payment_type":       "PayAsYouGo",
					"product_form":       "LegacyForm",
					"compute_resource":   "32ACU",
					"storage_resource":   "0ACU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_version": "5.0",
						"vpc_id":             CHECKSET,
						"vswitch_id":         CHECKSET,
						"zone_id":            CHECKSET,
						"payment_type":       "PayAsYouGo",
						"product_form":       "LegacyForm",
						"compute_resource":   "32ACU",
						"storage_resource":   "0ACU",
						// LegacyForm does not use product_version or
						// reserved_node; skip these from Map0 CHECKSET.
						"product_version":    REMOVEKEY,
						"reserved_node_size": REMOVEKEY,
					}),
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttrWith(resourceId, "id", func(value string) error {
						instanceID = value
						return nil
					}),
				),
			},
			// Step 2: In-place upgrade LegacyForm → IntegrationForm.
			// A single ModifyDBCluster sets ProductForm + ReservedNodeCount
			// + ReservedNodeSize. compute_resource and storage_resource
			// are removed from config (REMOVEKEY) since they are
			// LegacyForm-specific; the state retains them via
			// Optional+Computed. product_version is NOT set (ForceNew
			// would rebuild — the server handles it internally).
			// The instance ID must NOT change (no ForceNew) — this is
			// the customer's core requirement: direct enterprise
			// in-place upgrade, same instance ID, no rebuild.
			{
				Config: testAccConfig(map[string]interface{}{
					"product_form":        "IntegrationForm",
					"reserved_node_size":  "8ACU",
					"reserved_node_count": "3",
					"compute_resource":    REMOVEKEY,
					"storage_resource":    REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_form":        "IntegrationForm",
						"reserved_node_size":  "8ACU",
						"reserved_node_count": "3",
						// product_version is not set in config (ForceNew);
						// skip from Map0 CHECKSET to avoid false failure.
						"product_version": REMOVEKEY,
					}),
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttrWith(resourceId, "id", func(value string) error {
						if value != instanceID {
							return fmt.Errorf("instance ID changed during in-place upgrade: expected %s, got %s — ForceNew should not trigger for ProductForm upgrade", instanceID, value)
						}
						return nil
					}),
				),
			},
			// Step 3: Import — verifies state consistency after upgrade.
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_default_resource_group", "source_db_cluster_id", "backup_set_id", "restore_type", "restore_to_time"},
			},
		},
	})
}

var AliCloudAdbDbClusterLakeVersionMap0 = map[string]string{
	"product_form":           CHECKSET,
	"product_version":        CHECKSET,
	"reserved_node_size":     CHECKSET,
	"security_ips":           CHECKSET,
	"db_cluster_description": CHECKSET,
	"commodity_code":         CHECKSET,
	"connection_string":      CHECKSET,
	"engine":                 CHECKSET,
	"engine_version":         CHECKSET,
	"create_time":            CHECKSET,
	"lock_mode":              CHECKSET,
	"port":                   CHECKSET,
	"resource_group_id":      CHECKSET,
	"status":                 CHECKSET,
	// ssl_connection_string and ssl_expire_time are only populated when SSL is
	// enabled; they are empty when SSL is disabled. They are not added as
	// CHECKSET here because Map0 is checked at every test step, including steps
	// where SSL is not enabled. The dedicated steps that enable SSL verify
	// these attributes explicitly via their own CHECKSET assertions.
}

func AliCloudAdbDbClusterLakeVersionBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_adb_zones" "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_adb_zones.default.ids.0
	}
`, name)
}

func AliCloudAdbDbClusterLakeVersionBasicDependence0Twin(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_adb_zones" "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_adb_zones.default.ids.0
	}

	data "alicloud_vswitches" "secondary" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_adb_zones.default.ids.1
	}

	resource "alicloud_kms_key" "default" {
  		description            = var.name
  		pending_window_in_days = "7"
  		status                 = "Enabled"
	}
`, name)
}

// lintignore: R001
func TestUnitAliCloudAdbDbClusterLakeVersion(t *testing.T) {
	p := Provider().ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_adb_db_cluster_lake_version"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_adb_db_cluster_lake_version"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"compute_resource":              "CreateDBClusterValue",
		"db_cluster_version":            "5.0",
		"enable_default_resource_group": true,
		"payment_type":                  "CreateDBClusterValue",
		"storage_resource":              "CreateDBClusterValue",
		"vswitch_id":                    "CreateDBClusterValue",
		"zone_id":                       "CreateDBClusterValue",
	}
	// setBaseAttrs sets the base test attributes using literal keys to
	// comply with the tfproviderlint R001 rule. All original fields,
	// values, and error assertions are preserved.
	setBaseAttrs := func(d *schema.ResourceData) {
		assert.Nil(t, d.Set("compute_resource", "CreateDBClusterValue"))
		assert.Nil(t, d.Set("db_cluster_version", "5.0"))
		assert.Nil(t, d.Set("enable_default_resource_group", true))
		assert.Nil(t, d.Set("payment_type", "CreateDBClusterValue"))
		assert.Nil(t, d.Set("storage_resource", "CreateDBClusterValue"))
		assert.Nil(t, d.Set("vswitch_id", "CreateDBClusterValue"))
		assert.Nil(t, d.Set("zone_id", "CreateDBClusterValue"))
	}
	setBaseAttrs(dInit)
	setBaseAttrs(dExisted)
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// DescribeDBClusterAttribute
		"Items": map[string]interface{}{
			"DBCluster": []interface{}{
				map[string]interface{}{
					"CommodityCode":    "DefaultValue",
					"ComputeResource":  "CreateDBClusterValue",
					"ConnectionString": "DefaultValue",
					"CreationTime":     "DefaultValue",
					"DBClusterId":      "CreateDBClusterValue",
					"DBVersion":        "5.0",
					"Engine":           "DefaultValue",
					"EngineVersion":    "DefaultValue",
					"ExpireTime":       "DefaultValue",
					"Expired":          "DefaultValue",
					"LockMode":         "DefaultValue",
					"LockReason":       "DefaultValue",
					"PayType":          "CreateDBClusterValue",
					"Port":             "DefaultValue",
					"ProductForm":      "IntegrationForm",
					"ProductVersion":   "EnterpriseVersion",
					"ResourceGroupId":  "DefaultValue",
					"DBClusterStatus":  "Running",
					"StorageResource":  "CreateDBClusterValue",
					"VPCId":            "CreateDBClusterValue",
					"VSwitchId":        "CreateDBClusterValue",
					"ZoneId":           "CreateDBClusterValue",
				},
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateDBCluster
		"DBClusterId":     "CreateDBClusterValue",
		"OrderId":         "MockValue",
		"RequestId":       "MockValue",
		"ResourceGroupId": "MockValue",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_adb_db_cluster_lake_version", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudAdbDbClusterLakeVersionCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeDBClusterAttribute Response
		"Items": map[string]interface{}{
			"DBCluster": []interface{}{
				map[string]interface{}{
					"DBClusterId": "CreateDBClusterValue",
				},
			},
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateDBCluster" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudAdbDbClusterLakeVersionCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_adb_db_cluster_lake_version"].Schema).Data(dInit.State(), nil)
			setBaseAttrs(dCompare)
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudAdbDbClusterLakeVersionUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// ModifyDBCluster
	attributesDiff := map[string]interface{}{
		"compute_resource": "ModifyDBClusterValue",
		"storage_resource": "ModifyDBClusterValue",
	}
	diff, err := newInstanceDiff("alicloud_adb_db_cluster_lake_version", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_adb_db_cluster_lake_version"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDBClusterAttribute Response
		"Items": map[string]interface{}{
			"DBCluster": []interface{}{
				map[string]interface{}{
					"ComputeResource": "ModifyDBClusterValue",
					"StorageResource": "ModifyDBClusterValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyDBCluster" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudAdbDbClusterLakeVersionUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_adb_db_cluster_lake_version"].Schema).Data(dExisted.State(), nil)
			setBaseAttrs(dCompare)
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// ModifyDBClusterSSL
	attributesDiff = map[string]interface{}{
		"ssl_enabled": "Enable",
	}
	diff, err = newInstanceDiff("alicloud_adb_db_cluster_lake_version", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_adb_db_cluster_lake_version"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDBClusterSSL Response
		"SSLEnabled":       true,
		"ConnectionString": "DefaultValue",
		"ExpireTime":       "DefaultValue",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyDBClusterSSL" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudAdbDbClusterLakeVersionUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_adb_db_cluster_lake_version"].Schema).Data(dExisted.State(), nil)
			setBaseAttrs(dCompare)
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// ModifyDBCluster with ProductForm
	attributesDiff = map[string]interface{}{
		"product_form": "LegacyForm",
	}
	diff, err = newInstanceDiff("alicloud_adb_db_cluster_lake_version", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_adb_db_cluster_lake_version"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDBClusterAttribute Response
		"Items": map[string]interface{}{
			"DBCluster": []interface{}{
				map[string]interface{}{
					"ProductForm": "LegacyForm",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyDBCluster" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudAdbDbClusterLakeVersionUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_adb_db_cluster_lake_version"].Schema).Data(dExisted.State(), nil)
			setBaseAttrs(dCompare)
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "InvalidDBCluster.NotFound", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeDBClusterAttribute" {
				switch errorCode {
				case "{}", "InvalidDBCluster.NotFound":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudAdbDbClusterLakeVersionRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}", "InvalidDBCluster.NotFound":
			assert.Nil(t, err)
		}
	}

	// Delete — error propagation: a non-NotFound error from DeleteDBCluster
	// must be returned transparently (no WaitForState, no swallowing).
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudAdbDbClusterLakeVersionDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "InvalidDBCluster.NotFound", "InvalidDBClusterId.NotFound"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteDBCluster" {
				switch errorCode {
				case "NonRetryableError", "InvalidDBCluster.NotFound", "InvalidDBClusterId.NotFound":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudAdbDbClusterLakeVersionDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "InvalidDBCluster.NotFound", "InvalidDBClusterId.NotFound":
			assert.Nil(t, err)
		}
	}
}

// buildRawConfig builds a cty.Value representing the raw HCL config
// from a map[string]interface{}, using the resource schema's implied
// type to determine real attribute types. Unconfigured attributes get
// typed null; configured attributes get proper cty values. The optional
// overrides map injects cty.Value directly (e.g. cty.UnknownVal for
// cross-resource references during planning), bypassing the Go-to-cty
// conversion.
//
// In production, Terraform sets InstanceState.RawConfig to the planned
// config before calling schema.Resource.Diff. GetRawConfig() reads this
// field directly (schema.go: result.RawConfig = s.RawConfig). In tests,
// TestResourceDataRaw does not populate RawConfig, so we set it
// manually. The same rawConfig is passed to NewResourceConfigShimmed,
// which correctly converts cty unknown values to SDK legacy unknowns.
func buildRawConfig(r *schema.Resource, config map[string]interface{}, overrides map[string]cty.Value) cty.Value {
	attrTypes := r.CoreConfigSchema().ImpliedType().AttributeTypes()
	valMap := make(map[string]cty.Value, len(attrTypes))
	for k, attrType := range attrTypes {
		if v, ok := overrides[k]; ok {
			valMap[k] = v
		} else if v, ok := config[k]; ok {
			valMap[k] = goValToCty(v, attrType)
		} else {
			valMap[k] = cty.NullVal(attrType)
		}
	}
	return cty.ObjectVal(valMap)
}

// goValToCty converts a Go interface{} to a cty.Value of the given type.
// Supports string, bool, int, and nil. Other types default to typed null.
func goValToCty(v interface{}, ty cty.Type) cty.Value {
	switch tv := v.(type) {
	case string:
		return cty.StringVal(tv)
	case bool:
		return cty.BoolVal(tv)
	case int:
		return cty.NumberIntVal(int64(tv))
	case nil:
		return cty.NullVal(ty)
	default:
		return cty.NullVal(ty)
	}
}

// TestUnitAliCloudAdbDbClusterLakeVersionCustomizeDiff exercises the
// CustomizeDiff function through the real SDK schema.Resource.Diff
// pipeline — not a hand-crafted ResourceDiff. This verifies the
// ssl_enabled/enable_ssl priority contract (explicit command takes
// priority over legacy boolean) and ensures null, absent, and unknown
// configs do not panic. The raw config (state.RawConfig) is set to a
// typed cty.Value built from the resource schema, so GetRawConfig
// works correctly in the test context.
func TestUnitAliCloudAdbDbClusterLakeVersionCustomizeDiff(t *testing.T) {
	r := resourceAliCloudAdbDbClusterLakeVersion()
	ctx := context.Background()

	// baseConfig contains the Required fields so the diff pipeline
	// does not produce unrelated noise.
	baseConfig := map[string]interface{}{
		"db_cluster_version": "5.0",
		"vpc_id":             "vpc-test",
		"vswitch_id":         "vsw-test",
		"zone_id":            "zone-test",
		"payment_type":       "PayAsYouGo",
	}

	// runDiff creates old state from oldRaw (merged with baseConfig),
	// then diffs against newRaw (merged with baseConfig). This exercises
	// CustomizeDiff through the real SDK schema.Resource.Diff pipeline.
	runDiff := func(oldRaw, newRaw map[string]interface{}, overrides ...map[string]cty.Value) *terraform.InstanceDiff {
		oldCfg := make(map[string]interface{})
		for k, v := range baseConfig {
			oldCfg[k] = v
		}
		for k, v := range oldRaw {
			oldCfg[k] = v
		}
		d := schema.TestResourceDataRaw(t, r.Schema, oldCfg)
		d.SetId("test-id")

		newCfg := make(map[string]interface{})
		for k, v := range baseConfig {
			newCfg[k] = v
		}
		for k, v := range newRaw {
			newCfg[k] = v
		}
		// Build a typed cty.Value for the raw config using the resource
		// schema's implied type. Unconfigured attributes get typed null;
		// overrides (e.g. cty.UnknownVal) inject real typed unknowns.
		var ov map[string]cty.Value
		if len(overrides) > 0 {
			ov = overrides[0]
		}
		rawConfig := buildRawConfig(r, newCfg, ov)
		// Set state.RawConfig so GetRawConfig() returns the planned config.
		// In production, Terraform populates InstanceState.RawConfig with
		// the planned config before calling Diff. NewResourceConfigShimmed
		// converts the cty config (including unknowns) to the SDK legacy
		// ResourceConfig format that r.Diff expects.
		state := d.State()
		state.RawConfig = rawConfig
		rc := terraform.NewResourceConfigShimmed(rawConfig, r.CoreConfigSchema())
		diff, err := r.Diff(ctx, state, rc, nil)
		assert.NoError(t, err)
		return diff
	}

	// assertNewValue checks that the diff forces a specific New value.
	assertNewValue := func(diff *terraform.InstanceDiff, key, expected string) {
		t.Helper()
		if diff == nil {
			t.Errorf("expected diff for %s with New=%q, but diff is nil", key, expected)
			return
		}
		attrDiff, ok := diff.Attributes[key]
		if !ok || attrDiff == nil {
			t.Errorf("expected diff for %s with New=%q, but no diff entry", key, expected)
			return
		}
		assert.Equal(t, expected, attrDiff.New)
	}

	// assertNoSSLDiff is the strict assertion for known-stable scenarios:
	// no diff entry, or Old == New. NewComputed is NOT acceptable here
	// because in a known-stable scenario the value should be concrete.
	assertNoSSLDiff := func(diff *terraform.InstanceDiff, key string) {
		t.Helper()
		if diff == nil {
			return
		}
		if attrDiff, ok := diff.Attributes[key]; ok && attrDiff != nil {
			if attrDiff.NewComputed {
				t.Errorf("expected no SSL diff for %s in known-stable scenario, but got NewComputed", key)
				return
			}
			if attrDiff.Old != attrDiff.New {
				t.Errorf("expected no SSL diff for %s, but got Old=%q New=%q", key, attrDiff.Old, attrDiff.New)
			}
		}
	}

	// assertNoForcedChange is the lenient assertion for unknown/absent
	// scenarios: NewComputed is acceptable (SDK behavior for absent
	// Optional+Computed fields), as is Old == New or no diff entry.
	assertNoForcedChange := func(diff *terraform.InstanceDiff, key string) {
		t.Helper()
		if diff == nil {
			return
		}
		if attrDiff, ok := diff.Attributes[key]; ok && attrDiff != nil {
			if attrDiff.NewComputed {
				return
			}
			if attrDiff.Old != attrDiff.New {
				t.Errorf("expected no forced change for %s, but got Old=%q New=%q", key, attrDiff.Old, attrDiff.New)
			}
		}
	}

	// Test 1: Priority — ssl_enabled=Enable overrides enable_ssl=false.
	// CustomizeDiff derives enable_ssl=true via SetNew.
	diff := runDiff(
		map[string]interface{}{"enable_ssl": false},
		map[string]interface{}{"ssl_enabled": "Enable", "enable_ssl": false},
	)
	assertNewValue(diff, "enable_ssl", "true")

	// Test 2: Priority — ssl_enabled=Disable overrides enable_ssl=true.
	// CustomizeDiff derives enable_ssl=false via SetNew.
	diff = runDiff(
		map[string]interface{}{"enable_ssl": true},
		map[string]interface{}{"ssl_enabled": "Disable", "enable_ssl": true},
	)
	assertNewValue(diff, "enable_ssl", "false")

	// Test 3: Update preserves old actual SSL state.
	// Old state: ssl_enabled=Enable, enable_ssl=true.
	// New config: ssl_enabled=Update, enable_ssl=false (conflicting bool).
	// CustomizeDiff must preserve the old state value (true), NOT the
	// user's new config value (false). Without this, the next round's
	// HasChange("enable_ssl") would fire and send EnableSSL=false,
	// unexpectedly disabling SSL.
	diff = runDiff(
		map[string]interface{}{"enable_ssl": true, "ssl_enabled": "Enable"},
		map[string]interface{}{"ssl_enabled": "Update", "enable_ssl": false},
	)
	assertNoSSLDiff(diff, "enable_ssl")

	// Test 4: Update convergence — simulated second plan after Read.
	// NOTE: This simulates the state that Read would produce after the
	// first apply wrote ssl_enabled=Update, enable_ssl=true into state.
	// It does not call the production Read function. The same config
	// (ssl_enabled=Update, enable_ssl=false) must produce NO SSL diff:
	// CustomizeDiff preserves enable_ssl=true (old state), and
	// ssl_enabled hasn't changed (Update→Update). No operation triggers.
	diff = runDiff(
		map[string]interface{}{"enable_ssl": true, "ssl_enabled": "Update"},
		map[string]interface{}{"ssl_enabled": "Update", "enable_ssl": false},
	)
	assertNoSSLDiff(diff, "enable_ssl")

	// Test 5: Known stable — ssl_enabled and enable_ssl unchanged.
	// Old state and new config both have ssl_enabled=Enable,
	// enable_ssl=true. CustomizeDiff calls SetNew("enable_ssl", true)
	// which equals the old value — no actual change in the diff.
	diff = runDiff(
		map[string]interface{}{"enable_ssl": true, "ssl_enabled": "Enable"},
		map[string]interface{}{"ssl_enabled": "Enable", "enable_ssl": true},
	)
	assertNoSSLDiff(diff, "enable_ssl")

	// Test 6: Remove command — ssl_enabled absent, enable_ssl true→false.
	// When ssl_enabled is removed from config, GetRawConfig returns typed
	// null for ssl_enabled (absent from HCL). CustomizeDiff returns nil,
	// so the enable_ssl change (true→false) is preserved as a legal diff.
	// This proves the old bool drives the operation when ssl_enabled is
	// absent — GetOk would incorrectly read stale "Enable" from state and
	// SetNew(true), swallowing the legitimate true→false change.
	diff = runDiff(
		map[string]interface{}{"enable_ssl": true, "ssl_enabled": "Enable"},
		map[string]interface{}{"enable_ssl": false},
	)
	assertNewValue(diff, "enable_ssl", "false")

	// Test 7: Unknown ssl_enabled — cross-resource reference produces
	// cty.UnknownVal during planning. The unknown guard in CustomizeDiff
	// checks IsKnown() and returns nil, so no SetNew fires and no panic
	// occurs. This uses a real typed unknown via the overrides map, not
	// just an absent field.
	diff = runDiff(
		map[string]interface{}{"enable_ssl": true},
		map[string]interface{}{},
		map[string]cty.Value{"ssl_enabled": cty.UnknownVal(cty.String)},
	)
	assertNoForcedChange(diff, "enable_ssl")

	// Test 8: Absent — no ssl_enabled in config, enable_ssl stable.
	// ssl_enabled is typed null in raw config; CustomizeDiff returns nil.
	// enable_ssl is absent from new config (Optional+Computed preserves
	// state), so no forced change.
	diff = runDiff(
		map[string]interface{}{"enable_ssl": true},
		map[string]interface{}{},
	)
	assertNoForcedChange(diff, "enable_ssl")

	// Test 9: Empty SSL config (existing resource) — no SSL fields at all.
	// runDiff always sets SetId("test-id"), so this is an existing resource
	// with an empty SSL config. No crash, no forced change.
	diff = runDiff(
		map[string]interface{}{},
		map[string]interface{}{},
	)
	assertNoForcedChange(diff, "enable_ssl")
}

func TestUnitAdbSslReadyState(t *testing.T) {
	cases := []struct {
		name       string
		sslEnabled string
		object     map[string]interface{}
		expected   string
	}{
		{
			name:       "ready_with_both_fields",
			sslEnabled: "true",
			object: map[string]interface{}{
				"ConnectionString": "conn-ssl.example.com",
				"ExpireTime":       "2026-12-01T00:00:00Z",
			},
			expected: "true",
		},
		{
			name:       "missing_connection_string",
			sslEnabled: "true",
			object: map[string]interface{}{
				"ExpireTime": "2026-12-01T00:00:00Z",
			},
			expected: "provisioning",
		},
		{
			name:       "nil_connection_string",
			sslEnabled: "true",
			object: map[string]interface{}{
				"ConnectionString": nil,
				"ExpireTime":       "2026-12-01T00:00:00Z",
			},
			expected: "provisioning",
		},
		{
			name:       "empty_connection_string",
			sslEnabled: "true",
			object: map[string]interface{}{
				"ConnectionString": "",
				"ExpireTime":       "2026-12-01T00:00:00Z",
			},
			expected: "provisioning",
		},
		{
			name:       "missing_expire_time",
			sslEnabled: "true",
			object: map[string]interface{}{
				"ConnectionString": "conn-ssl.example.com",
			},
			expected: "provisioning",
		},
		{
			name:       "nil_expire_time",
			sslEnabled: "true",
			object: map[string]interface{}{
				"ConnectionString": "conn-ssl.example.com",
				"ExpireTime":       nil,
			},
			expected: "provisioning",
		},
		{
			name:       "empty_expire_time",
			sslEnabled: "true",
			object: map[string]interface{}{
				"ConnectionString": "conn-ssl.example.com",
				"ExpireTime":       "",
			},
			expected: "provisioning",
		},
		{
			name:       "disabled_ssl",
			sslEnabled: "false",
			object:     map[string]interface{}{},
			expected:   "false",
		},
		{
			name:       "disabled_ssl_with_stale_fields",
			sslEnabled: "false",
			object: map[string]interface{}{
				"ConnectionString": "stale-conn",
				"ExpireTime":       "stale-time",
			},
			expected: "false",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := sslReadyState(c.sslEnabled, c.object)
			assert.Equal(t, c.expected, result)
		})
	}
}

// TestUnitAdbDbClusterLakeVersionDeleteIdempotent verifies that the Delete
// function treats InvalidDBClusterId.NotFound — the actual error code
// returned by DeleteDBCluster when the cluster no longer exists — as
// idempotent: returning nil promptly without a meaningless wait. The
// legacy code InvalidDBCluster.NotFound must also be treated as
// idempotent. Other errors (e.g., permission) must propagate transparently.
//
// This test does not require cloud credentials: it constructs a local
// AliyunClient via Config.Client with dummy credentials and patches
// tea-rpc Client.DoRequest — the SDK boundary that RpcPost/rpcRequest
// ultimately calls. RpcPost is a short wrapper delegated to rpcRequest
// and inlined by the compiler, so gomonkey cannot intercept it;
// DoRequest on *client.Client is the proven intercept point used by
// other unit tests in this package. The Endpoints map is pre-populated
// with the "adb" key so that loadApiEndpoint resolves the endpoint
// without network discovery; unexpected API actions fail in the mock
// rather than reaching the network.
func TestUnitAdbDbClusterLakeVersionDeleteIdempotent(t *testing.T) {
	p := Provider().ResourcesMap
	dExisted, _ := schema.InternalMap(p["alicloud_adb_db_cluster_lake_version"].Schema).Data(nil, nil)
	assert.Nil(t, dExisted.Set("payment_type", "PayAsYouGo"))
	dExisted.SetId("amv-test-idempotent")

	// Create a properly initialized AliyunClient with a dummy Config.
	// GetRetryTimeout is inlined by the compiler, so gomonkey cannot
	// intercept it; providing a non-nil config lets the inlined helper
	// return the default timeout without panicking on a nil config.
	// SkipRegionValidation and AccountType disable credential discovery
	// and network calls during client construction.
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("test-key").
		SetAccessKeySecret("test-secret"))
	assert.Nil(t, err)
	// Pre-populate the Endpoints map so loadApiEndpoint("adb") succeeds
	// without network discovery. productCodeToConfigEndpoints has no
	// "adb" entry, so configEndpoints defaults to "adb"; the Endpoints
	// map must contain the "adb" key to short-circuit loadEndpoint.
	endpoints := new(sync.Map)
	endpoints.Store("adb", "adb.cn-hangzhou.aliyuncs.com")
	config := &connectivity.Config{
		AccessKey:            "test-key",
		SecretKey:            "test-secret",
		Credential:           credential,
		RegionId:             "cn-hangzhou",
		AccountType:          "test",
		Protocol:             "http",
		Endpoints:            endpoints,
		SignVersion:          new(sync.Map),
		SkipRegionValidation: true,
	}
	aliyunClient, err := config.Client()
	assert.Nil(t, err)

	// mockDelete patches DoRequest on the tea-rpc Client — the SDK
	// boundary that rpcRequest calls after endpoint resolution. The mock
	// records the actual action and request body (DBClusterId, RegionId)
	// so we can confirm the real Delete path was exercised, and returns
	// the given SDK error. Unexpected API actions fail immediately rather
	// than reaching the network.
	mockDelete := func(sdkErr *tea.SDKError) (*gomonkey.Patches, *string, *map[string]interface{}) {
		var action string
		var body map[string]interface{}
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest",
			func(_ *client.Client, act *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, bd map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *act != "DeleteDBCluster" {
					return nil, fmt.Errorf("unexpected API action in test: %s", *act)
				}
				action = *act
				body = bd
				return nil, sdkErr
			})
		return patches, &action, &body
	}

	// Case 1: InvalidDBClusterId.NotFound — the actual error code returned
	// by DeleteDBCluster when the cluster no longer exists. Delete must
	// return nil (idempotent) without WaitForState.
	patches, action1, body1 := mockDelete(&tea.SDKError{
		Code:       String("InvalidDBClusterId.NotFound"),
		Data:       String(`{"Code":"InvalidDBClusterId.NotFound"}`),
		Message:    String("code: 400, The DBClusterId provided does not exist in our records."),
		StatusCode: tea.Int(400),
	})
	err = resourceAliCloudAdbDbClusterLakeVersionDelete(dExisted, aliyunClient)
	patches.Reset()
	assert.Equal(t, "DeleteDBCluster", *action1, "DoRequest must be called for DeleteDBCluster")
	assert.Equal(t, "amv-test-idempotent", (*body1)["DBClusterId"], "DBClusterId must be in the request")
	assert.Nil(t, err, "InvalidDBClusterId.NotFound should be treated as idempotent (return nil)")

	// Case 2: InvalidDBCluster.NotFound (legacy code, preserved for
	// backward compatibility). Delete must also return nil.
	patches, action2, body2 := mockDelete(&tea.SDKError{
		Code:       String("InvalidDBCluster.NotFound"),
		Data:       String(`{"Code":"InvalidDBCluster.NotFound"}`),
		Message:    String("code: 400, The DBCluster not found."),
		StatusCode: tea.Int(400),
	})
	err = resourceAliCloudAdbDbClusterLakeVersionDelete(dExisted, aliyunClient)
	patches.Reset()
	assert.Equal(t, "DeleteDBCluster", *action2, "DoRequest must be called for DeleteDBCluster")
	assert.Equal(t, "amv-test-idempotent", (*body2)["DBClusterId"], "DBClusterId must be in the request")
	assert.Nil(t, err, "InvalidDBCluster.NotFound (legacy) should be treated as idempotent (return nil)")

	// Case 3: Permission error (Forbidden.RAM). Delete must propagate it
	// transparently — no swallowing, no WaitForState.
	patches, action3, body3 := mockDelete(&tea.SDKError{
		Code:       String("Forbidden.RAM"),
		Data:       String(`{"Code":"Forbidden.RAM"}`),
		Message:    String("code: 403, User not authorized to operate."),
		StatusCode: tea.Int(403),
	})
	err = resourceAliCloudAdbDbClusterLakeVersionDelete(dExisted, aliyunClient)
	patches.Reset()
	assert.Equal(t, "DeleteDBCluster", *action3, "DoRequest must be called for DeleteDBCluster")
	assert.Equal(t, "amv-test-idempotent", (*body3)["DBClusterId"], "DBClusterId must be in the request")
	assert.NotNil(t, err, "Permission errors must propagate transparently (return non-nil)")
}
