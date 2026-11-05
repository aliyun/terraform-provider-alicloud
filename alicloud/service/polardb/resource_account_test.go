package polardb_test

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	polardb "github.com/aliyun/terraform-provider-alicloud/alicloud/service/polardb"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Test PolarDb Account. >>> Resource test cases, automatically generated.
// Case  Account用例_PG_高权限账号 11819
func TestAccAliCloudPolarDbAccount_basic11819(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_account.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudPolarDbAccountMap11819)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolarDbAccount")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudPolarDbAccountBasicDependence11819)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":    "${alicloud_polardb_cluster.default.id}",
					"account_name":     name,
					"account_password": "YourPassword123!",
					"account_type":     "Super",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id": CHECKSET,
						"account_name":  name,
						"account_type":  "Super",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "YourPassword123!update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "reset_permission_flag"},
			},
		},
	})
}

func TestAccAliCloudPolarDbAccount_basic11819_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_account.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudPolarDbAccountMap11819)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolarDbAccount")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudPolarDbAccountBasicDependence11819)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":       "${alicloud_polardb_cluster.default.id}",
					"account_name":        name,
					"account_password":    "YourPassword123!",
					"account_type":        "Super",
					"account_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id":       CHECKSET,
						"account_name":        name,
						"account_type":        "Super",
						"account_description": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "reset_permission_flag"},
			},
		},
	})
}

var AliCloudPolarDbAccountMap11819 = map[string]string{
	"account_lock_state":          CHECKSET,
	"account_password_valid_time": CHECKSET,
	"status":                      CHECKSET,
}

func AliCloudPolarDbAccountBasicDependence11819(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_polardb_node_classes" "default" {
  db_type    = "MySQL"
  db_version = "8.0"
  pay_type   = "PostPaid"
  category   = "Normal"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_polardb_node_classes.default.classes.0.zone_id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
}

resource "alicloud_polardb_cluster" "default" {
  db_type            = "MySQL"
  db_version         = "8.0"
  db_node_class      = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
  pay_type           = "PostPaid"
  vswitch_id         = alicloud_vswitch.default.id
  encrypt_new_tables = "OFF"
}
`, name)
}

// Case  Account用例 9211
func TestAccAliCloudPolarDbAccount_basic9211(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_account.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudPolarDbAccountMap9211)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolarDbAccount")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudPolarDbAccountBasicDependence9211)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":          "${alicloud_polardb_cluster.default.id}",
					"account_name":           name,
					"kms_encrypted_password": "${alicloud_kms_ciphertext.default.ciphertext_blob}",
					"kms_encryption_context": map[string]string{
						"name": name,
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id": CHECKSET,
						"account_name":  name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kms_encrypted_password": "${alicloud_kms_ciphertext.update.ciphertext_blob}",
					"kms_encryption_context": map[string]string{
						"name": name + "update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kms_encrypted_password":   CHECKSET,
						"kms_encryption_context.%": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_lock_state": "Lock",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_lock_state": "Lock",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password_valid_time": "2126-09-17T10:00:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password_valid_time": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_lock_state": "UnLock",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_lock_state": "UnLock",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"kms_encrypted_password", "kms_encryption_context", "reset_permission_flag"},
			},
		},
	})
}

func TestAccAliCloudPolarDbAccount_basic9211_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_account.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudPolarDbAccountMap9211)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolarDbAccount")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudPolarDbAccountBasicDependence9211)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":          "${alicloud_polardb_cluster.default.id}",
					"account_name":           name,
					"kms_encrypted_password": "${alicloud_kms_ciphertext.default.ciphertext_blob}",
					"kms_encryption_context": map[string]string{
						"name": name,
					},
					"account_type":                "Normal",
					"account_description":         name,
					"account_lock_state":          "Lock",
					"account_password_valid_time": "2126-09-17T10:00:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id":               CHECKSET,
						"account_name":                name,
						"account_type":                "Normal",
						"account_description":         name,
						"account_lock_state":          "Lock",
						"account_password_valid_time": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"kms_encrypted_password", "kms_encryption_context", "reset_permission_flag"},
			},
		},
	})
}

var AliCloudPolarDbAccountMap9211 = map[string]string{
	"account_lock_state":          CHECKSET,
	"account_password_valid_time": CHECKSET,
	"account_type":                CHECKSET,
	"status":                      CHECKSET}

func AliCloudPolarDbAccountBasicDependence9211(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_polardb_node_classes" "default" {
  db_type  = "PostgreSQL"
  pay_type = "PostPaid"
  category = "Normal"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_polardb_node_classes.default.classes.0.zone_id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
}

resource "alicloud_polardb_cluster" "default" {
  db_version         = "14"
  pay_type           = "PostPaid"
  db_node_class      = "polar.pg.x4.medium"
  db_type            = "PostgreSQL"
  vswitch_id         = alicloud_vswitch.default.id
  enable_dynamodb    = true
  encrypt_new_tables = "OFF"
}

	resource "alicloud_kms_key" "default" {
  		description            = var.name
  		status                 = "Enabled"
  		pending_window_in_days = 7
	}

	resource "alicloud_kms_ciphertext" "default" {
  		key_id    = alicloud_kms_key.default.id
  		plaintext = "YourPassword1234!"
  		encryption_context = {
    		"name" = var.name
  		}
	}

	resource "alicloud_kms_ciphertext" "update" {
  		key_id    = alicloud_kms_key.default.id
  		plaintext = "YourPassword1234!update"
  		encryption_context = {
    		"name" = "${var.name}update"
  		}
	}
`, name)
}

// Test PolarDb Account. <<< Resource test cases, automatically generated.

// Test PolarDb Account DynamoDB type
func TestAccAliCloudPolarDbAccount_dynamoDB(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_account.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudPolarDbAccountMapDynamoDB)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolarDbAccount")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudPolarDbAccountBasicDependenceDynamoDB)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":    "${alicloud_polardb_cluster.default.id}",
					"account_name":     "dynamodb",
					"account_password": "YourPassword123!",
					"account_type":     "DynamoDB",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id":          CHECKSET,
						"account_name":           "dynamodb",
						"account_type":           "DynamoDB",
						"dynamodb_auth_password": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "dynamodb_auth_password", "reset_permission_flag"},
			},
		},
	})
}

var AliCloudPolarDbAccountMapDynamoDB = map[string]string{
	"account_lock_state":          CHECKSET,
	"account_password_valid_time": CHECKSET,
	"status":                      CHECKSET,
}

func AliCloudPolarDbAccountBasicDependenceDynamoDB(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.7.id
  vswitch_name = var.name
}

resource "alicloud_polardb_cluster" "default" {
  db_type         = "PostgreSQL"
  db_version      = "16"
  db_node_class   = "polar.pg.x4.medium"
  pay_type        = "PostPaid"
  vswitch_id      = alicloud_vswitch.default.id
  enable_dynamodb = true
  encrypt_new_tables = "OFF"
}
`, name)
}

func TestAccAliCloudPolarDbAccount_passwordWo(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_account.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudPolarDbAccountMap11819)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolarDbAccount")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudPolarDbAccountBasicDependence11819)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id": "${alicloud_polardb_cluster.default.id}",
					"account_name":  name,
				}),
				ExpectError: regexp.MustCompile("one\\s+of\\s+`account_password,account_password_wo,kms_encrypted_password`\\s+must\\s+be\\s+specified"),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":               "${alicloud_polardb_cluster.default.id}",
					"account_name":                name,
					"account_password_wo":         "YourPassword123!",
					"account_password_wo_version": 1,
					"account_type":                "Super",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id":               CHECKSET,
						"account_name":                name,
						"account_type":                "Super",
						"account_password_wo_version": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password_wo":         "YourPassword123!update",
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
					"account_password":            "YourPassword123!",
					"account_password_wo":         acctest.REMOVEKEY,
					"account_password_wo_version": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password":            "YourPassword123!",
						"account_password_wo_version": acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "account_password_wo_version", "reset_permission_flag"},
			},
		},
	})
}

// Test PolarDb Account property coverage for CreateAccount (db_name, priv_for_all_db, node_type) and ModifyAccountPassword (password_type).
func TestAccAliCloudPolarDbAccount_propertyCoverage(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_account.default"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{
		"account_lock_state":          CHECKSET,
		"account_password_valid_time": CHECKSET,
		"status":                      CHECKSET,
	})
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolarDbAccount")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudPolarDbAccountBasicDependencePropertyCoverage)
	testAccConfigWithPostCreationDatabase := withPolarDbAccountPostCreationDatabase(testAccConfig)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		IDRefreshName:            resourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":     "${alicloud_polardb_cluster.default.id}",
					"account_name":      name,
					"account_password":  "YourPassword123!",
					"db_name":           "${alicloud_polardb_database.default.db_name}",
					"account_privilege": "ReadWrite",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id":     CHECKSET,
						"account_name":      name,
						"account_type":      CHECKSET,
						"db_name":           "tf_testacc_db",
						"account_privilege": "ReadWrite",
					}),
					testAccCheckPolarDbAccountDatabasePrivileges(map[string]string{"tf_testacc_db": "ReadWrite"}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"account_password", "db_name", "password_type", "priv_for_all_db", "account_privilege", "reset_permission_flag",
				},
			},
			{
				// A comma-separated privilege list must pass schema validation and be applied per database.
				Config: testAccConfig(map[string]interface{}{
					"db_name":           "tf_testacc_db,tf_testacc_db2",
					"account_privilege": "ReadWrite,ReadOnly",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_name":           "tf_testacc_db,tf_testacc_db2",
						"account_privilege": "ReadWrite,ReadOnly",
					}),
					testAccCheckPolarDbAccountDatabasePrivileges(map[string]string{
						"tf_testacc_db":  "ReadWrite",
						"tf_testacc_db2": "ReadOnly",
					}),
				),
			},
			{
				// An explicit Normal account with a cluster-wide read-only privilege must not keep
				// proposing a replacement after the API reports ReadOnly. PrivForAllDB applies
				// AccountPrivilege to every database of the cluster, but DescribeAccounts returns an
				// empty DatabasePrivileges list for such an account (observed with
				// priv_for_all_db = "1"), so the per-database effect cannot be asserted here.
				Config: testAccConfig(map[string]interface{}{
					"db_name":           "tf_testacc_db",
					"account_privilege": "ReadOnly",
					"priv_for_all_db":   "1",
					"account_type":      "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_name":           "tf_testacc_db",
						"account_privilege": "ReadOnly",
						"priv_for_all_db":   "1",
					}),
					testAccCheckPolarDbAccountTypeMatchesRemote(),
					testAccCheckPolarDbAccountDatabasePrivilegesEmpty(),
				),
			},
			{
				// Create a database after the cluster-wide account exists so the scenario covers a
				// database added later, not only databases prepared before account creation.
				// DescribeAccounts does not enumerate PrivForAllDB grants, so this step records the
				// account response while the follow-up database exists and does not claim the
				// per-database effect of the grant. The extra resource is appended by a helper so
				// the step's Config text stays a plain attribute map for the coverage checker.
				Config: testAccConfigWithPostCreationDatabase(map[string]interface{}{
					"db_name": "tf_testacc_db",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_polardb_database.post_creation", "id"),
					testAccCheckPolarDbAccountDatabasePrivilegesEmpty(),
				),
			},
			{
				// Dropping the follow-up database must leave the cluster-wide account unchanged and
				// still reporting no explicit per-database entries.
				Config: testAccConfig(map[string]interface{}{
					"priv_for_all_db": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolarDbAccountDatabasePrivilegesEmpty(),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "YourPassword123!update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password_valid_time": CHECKSET,
					}),
				),
			},
			{
				// password_type only changes how a password is written. Adding it without a password
				// change must not issue a password write; this account is not a Tair account, so a
				// ModifyAccountPassword call would fail here.
				Config: testAccConfig(map[string]interface{}{
					"password_type": "Tair",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_type": "Tair",
					}),
				),
			},
			{
				// Removing password_type is not a password update either.
				Config: testAccConfig(map[string]interface{}{
					"password_type": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password_valid_time": CHECKSET,
						"password_type":               REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "YourPassword123!update2",
					"password_type":    "Tair",
				}),
				ExpectError: regexp.MustCompile("InvalidParam|tair account not exists"),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_type": "Search",
				}),
				ExpectError: regexp.MustCompile("OperationDenied|InvalidNodeType|InvalidParameter"),
			},
			{
				// account_privilege without db_name must fail before any API call.
				// The negative steps run last because they mutate the shared step attribute map.
				Config: testAccConfig(map[string]interface{}{
					"db_name": REMOVEKEY,
				}),
				ExpectError: regexp.MustCompile("all of `account_privilege,db_name` must be specified"),
			},
			{
				// priv_for_all_db = "1" without account_privilege must fail before any API call.
				Config: testAccConfig(map[string]interface{}{
					"account_privilege": REMOVEKEY,
				}),
				ExpectError: regexp.MustCompile(`"priv_for_all_db" set to "1" requires "account_privilege"`),
			},
			{
				// An unknown entry in the privilege list must be rejected locally.
				Config: testAccConfig(map[string]interface{}{
					"db_name":           "tf_testacc_db",
					"account_privilege": "ReadWrite,Unknown",
				}),
				ExpectError: regexp.MustCompile(`contains invalid value "Unknown"`),
			},
			{
				// An explicitly empty db_name must not slip past the RequiredWith key presence check
				// and reach CreateAccount without DBName while AccountPrivilege is sent.
				Config: testAccConfig(map[string]interface{}{
					"db_name":           "",
					"account_privilege": "ReadWrite",
				}),
				ExpectError: regexp.MustCompile(`"account_privilege" requires "db_name" to contain at least one database name`),
			},
			{
				// The negative steps above leave the shared attribute map and the resource state out of
				// sync (the failed Tair step configured password `update2` while the state kept `update`),
				// and the failed Search replacement may even have removed the account. This step therefore
				// restores a valid configuration by explicitly writing a known password back; depending on
				// the previous replacement it may update the password or recreate the account, so it is a
				// real write, not a no-op.
				Config: testAccConfig(map[string]interface{}{
					"db_name":           "tf_testacc_db",
					"account_privilege": "ReadOnly",
					"account_password":  "YourPassword123!restore",
					"password_type":     REMOVEKEY,
					"node_type":         REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_name":           "tf_testacc_db",
						"account_privilege": "ReadOnly",
						"priv_for_all_db":   "1",
					}),
				),
			},
		},
	})
}

// Test PolarDb Account permission reset (ResetAccount) and password reset (ResetAccountPassword) operations.
func TestAccAliCloudPolarDbAccount_permissionReset(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_account.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudPolarDbAccountMap11819)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolarDbAccount")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudPolarDbAccountBasicDependencePermissionReset)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":    "${alicloud_polardb_cluster.default.id}",
					"account_name":     name,
					"account_password": "YourPassword123!",
					"account_type":     "Super",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id": CHECKSET,
						"account_name":  name,
						"account_type":  "Super",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"reset_permission_flag": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"reset_permission_flag": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "YourPassword123!reset",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "reset_permission_flag"},
			},
		},
	})
}

func AliCloudPolarDbAccountBasicDependencePermissionReset(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_polardb_node_classes" "default" {
  db_type    = "MySQL"
  db_version = "8.0"
  pay_type   = "PostPaid"
  category   = "Normal"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "tf-testAcc-polardb-account-%s"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_polardb_node_classes.default.classes.0.zone_id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "tf-testAcc-polardb-account-%s"
}

resource "alicloud_polardb_cluster" "default" {
  db_type            = "MySQL"
  db_version         = "8.0"
  db_node_class      = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
  pay_type           = "PostPaid"
  vswitch_id         = alicloud_vswitch.default.id
  description        = "tf-testAcc-polardb-account-%s"
  encrypt_new_tables = "OFF"
}

resource "alicloud_polardb_account" "other" {
  db_cluster_id    = alicloud_polardb_cluster.default.id
  account_name     = "%sother"
  account_password = "YourPassword123!"
  account_type     = "Normal"
}
`, name, name, name, name, name)
}

func AliCloudPolarDbAccountBasicDependencePropertyCoverage(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_polardb_node_classes" "default" {
  db_type    = "MySQL"
  db_version = "8.0"
  pay_type   = "PostPaid"
  category   = "Normal"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "tf-testAcc-polardb-account-%s"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_polardb_node_classes.default.classes.0.zone_id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "tf-testAcc-polardb-account-%s"
}

resource "alicloud_polardb_cluster" "default" {
  db_type            = "MySQL"
  db_version         = "8.0"
  db_node_class      = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
  pay_type           = "PostPaid"
  vswitch_id         = alicloud_vswitch.default.id
  description        = "tf-testAcc-polardb-account-%s"
  encrypt_new_tables = "OFF"
}

resource "alicloud_polardb_database" "default" {
  db_cluster_id = alicloud_polardb_cluster.default.id
  db_name       = "tf_testacc_db"
}

resource "alicloud_polardb_database" "second" {
  db_cluster_id = alicloud_polardb_cluster.default.id
  db_name       = "tf_testacc_db2"
}
`, name, name, name, name)
}

// polarDbAccountDatabasePrivilegesFromObject extracts the explicit per-database privileges
// reported by DescribeAccounts. A cluster-wide grant (PrivForAllDB) is not enumerated in this
// field, so callers must distinguish a field the service omitted from an explicitly empty list;
// silently treating either an omitted field or malformed entries as "no privileges" would make
// the assertion unable to tell a parse failure from a real empty response.
func polarDbAccountDatabasePrivilegesFromObject(account map[string]interface{}) (map[string]string, bool, error) {
	raw, reported := account["DatabasePrivileges"]
	if !reported || raw == nil {
		return map[string]string{}, false, nil
	}
	items, ok := raw.([]interface{})
	if !ok {
		return nil, true, fmt.Errorf("DescribeAccounts returned DatabasePrivileges with unexpected type %T", raw)
	}
	privileges := make(map[string]string, len(items))
	for _, item := range items {
		privilege, ok := item.(map[string]interface{})
		if !ok {
			return nil, true, fmt.Errorf("DescribeAccounts returned a DatabasePrivileges entry with unexpected type %T", item)
		}
		dbName, _ := privilege["DBName"].(string)
		accountPrivilege, _ := privilege["AccountPrivilege"].(string)
		if dbName == "" {
			return nil, true, fmt.Errorf("DescribeAccounts returned a DatabasePrivileges entry without DBName: %v", privilege)
		}
		privileges[dbName] = accountPrivilege
	}
	return privileges, true, nil
}

func polarDbAccountDatabasePrivilegesMatch(actual, expected map[string]string) bool {
	if len(actual) != len(expected) {
		return false
	}
	for dbName, privilege := range expected {
		if actual[dbName] != privilege {
			return false
		}
	}
	return true
}

func polarDbAccountDatabasePrivilegeList(privileges map[string]string) []string {
	list := make([]string, 0, len(privileges))
	for dbName, privilege := range privileges {
		list = append(list, dbName+":"+privilege)
	}
	sort.Strings(list)
	return list
}

const polarDbAccountPostCreationDatabase = `
resource "alicloud_polardb_database" "post_creation" {
  db_cluster_id = alicloud_polardb_cluster.default.id
  db_name       = "tf_testacc_db3"
}
`

// withPolarDbAccountPostCreationDatabase appends a database resource to the generated account
// config. Keeping the literal here instead of inline in the test step keeps each step's Config a
// plain attribute map so the testing-coverage checker can still parse the later steps.
func withPolarDbAccountPostCreationDatabase(config acctest.ResourceTestAccConfigFunc) acctest.ResourceTestAccConfigFunc {
	return func(attributes map[string]interface{}) string {
		return config(attributes) + polarDbAccountPostCreationDatabase
	}
}

func testAccCheckPolarDbAccountDatabasePrivileges(expected map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["alicloud_polardb_account.default"]
		if !ok {
			return fmt.Errorf("resource alicloud_polardb_account.default not found in state")
		}
		accountName := rs.Primary.Attributes["account_name"]
		service := polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
		var got map[string]string
		for i := 0; i < 24; i++ {
			account, err := service.DescribePolarDbAccount(rs.Primary.ID)
			if err != nil {
				return err
			}
			if respondedName, _ := account["AccountName"].(string); respondedName != accountName {
				return fmt.Errorf("DescribeAccounts returned account %q instead of %q", respondedName, accountName)
			}
			privileges, _, err := polarDbAccountDatabasePrivilegesFromObject(account)
			if err != nil {
				return err
			}
			got = privileges
			// Compare the full set, not only the expected keys, so an unintended grant to
			// another database cannot satisfy the assertion.
			if polarDbAccountDatabasePrivilegesMatch(privileges, expected) {
				return nil
			}
			time.Sleep(5 * time.Second)
		}
		return fmt.Errorf("expected exactly database privileges %v on account %s, got %v", polarDbAccountDatabasePrivilegeList(expected), rs.Primary.ID, polarDbAccountDatabasePrivilegeList(got))
	}
}

func testAccCheckPolarDbAccountDatabasePrivilegesEmpty() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["alicloud_polardb_account.default"]
		if !ok {
			return fmt.Errorf("resource alicloud_polardb_account.default not found in state")
		}
		accountName := rs.Primary.Attributes["account_name"]
		service := polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
		var got []string
		for i := 0; i < 24; i++ {
			account, err := service.DescribePolarDbAccount(rs.Primary.ID)
			if err != nil {
				return err
			}
			if respondedName, _ := account["AccountName"].(string); respondedName != accountName {
				return fmt.Errorf("DescribeAccounts returned account %q instead of %q", respondedName, accountName)
			}
			privileges, reported, err := polarDbAccountDatabasePrivilegesFromObject(account)
			if err != nil {
				return err
			}
			got = polarDbAccountDatabasePrivilegeList(privileges)
			if len(got) == 0 {
				// Record whether the service omitted the field or returned an explicitly empty
				// list. DescribeAccounts has no read API for cluster-wide grants, so this is the
				// strongest control-plane fact available for a PrivForAllDB account.
				log.Printf("[INFO] alicloud_polardb_account %s: cluster-wide account reports no explicit DatabasePrivileges (fieldReported=%t, accountType=%v)", rs.Primary.ID, reported, account["AccountType"])
				return nil
			}
			time.Sleep(5 * time.Second)
		}
		return fmt.Errorf("expected no explicit database privileges for a cluster-wide grant on account %s, got %v", rs.Primary.ID, got)
	}
}

func testAccCheckPolarDbAccountTypeMatchesRemote() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["alicloud_polardb_account.default"]
		if !ok {
			return fmt.Errorf("resource alicloud_polardb_account.default not found in state")
		}
		service := polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
		objectRaw, err := service.DescribePolarDbAccount(rs.Primary.ID)
		if err != nil {
			return err
		}
		remote, _ := objectRaw["AccountType"].(string)
		local := rs.Primary.Attributes["account_type"]
		log.Printf("[INFO] alicloud_polardb_account account_type state=%q remote=%q", local, remote)
		if remote != local {
			return fmt.Errorf("expected state account_type %q to match DescribeAccounts AccountType %q", local, remote)
		}
		return nil
	}
}

// TestUnitPolarDbAccountDatabasePrivilegesParse covers the assertion helper's ability to tell an
// omitted DatabasePrivileges field, an explicitly empty list, a populated list and malformed
// entries apart. Without that distinction the cluster-wide (PrivForAllDB) step could pass on a
// parse failure instead of an actual empty response.
func TestUnitPolarDbAccountDatabasePrivilegesParse(t *testing.T) {
	cases := []struct {
		name         string
		account      map[string]interface{}
		want         map[string]string
		wantReported bool
		wantErr      bool
	}{
		{
			name:         "field omitted",
			account:      map[string]interface{}{"AccountName": "acc"},
			want:         map[string]string{},
			wantReported: false,
		},
		{
			name:         "field nil",
			account:      map[string]interface{}{"DatabasePrivileges": nil},
			want:         map[string]string{},
			wantReported: false,
		},
		{
			name:         "explicitly empty list",
			account:      map[string]interface{}{"DatabasePrivileges": []interface{}{}},
			want:         map[string]string{},
			wantReported: true,
		},
		{
			name: "populated list",
			account: map[string]interface{}{"DatabasePrivileges": []interface{}{
				map[string]interface{}{"DBName": "db1", "AccountPrivilege": "ReadWrite"},
				map[string]interface{}{"DBName": "db2", "AccountPrivilege": "ReadOnly"},
			}},
			want:         map[string]string{"db1": "ReadWrite", "db2": "ReadOnly"},
			wantReported: true,
		},
		{
			name:         "field of unexpected type",
			account:      map[string]interface{}{"DatabasePrivileges": "db1:ReadWrite"},
			wantReported: true,
			wantErr:      true,
		},
		{
			name:         "entry of unexpected type",
			account:      map[string]interface{}{"DatabasePrivileges": []interface{}{"db1:ReadWrite"}},
			wantReported: true,
			wantErr:      true,
		},
		{
			name: "entry without database name",
			account: map[string]interface{}{"DatabasePrivileges": []interface{}{
				map[string]interface{}{"AccountPrivilege": "ReadWrite"},
			}},
			wantReported: true,
			wantErr:      true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, reported, err := polarDbAccountDatabasePrivilegesFromObject(tc.account)
			if tc.wantErr != (err != nil) {
				t.Fatalf("%s: wantErr=%v, got %v", tc.name, tc.wantErr, err)
			}
			if err != nil {
				return
			}
			if reported != tc.wantReported {
				t.Fatalf("%s: want reported=%v, got %v", tc.name, tc.wantReported, reported)
			}
			if !polarDbAccountDatabasePrivilegesMatch(got, tc.want) {
				t.Fatalf("%s: want %v, got %v", tc.name, tc.want, got)
			}
		})
	}
}

func TestUnitPolarDbAccountDatabasePrivilegesMatch(t *testing.T) {
	cases := []struct {
		name     string
		actual   map[string]string
		expected map[string]string
		want     bool
	}{
		{"equal", map[string]string{"db1": "ReadWrite"}, map[string]string{"db1": "ReadWrite"}, true},
		{"both empty", map[string]string{}, map[string]string{}, true},
		{"unexpected extra grant", map[string]string{"db1": "ReadWrite", "db2": "ReadOnly"}, map[string]string{"db1": "ReadWrite"}, false},
		{"missing grant", map[string]string{"db1": "ReadWrite"}, map[string]string{"db1": "ReadWrite", "db2": "ReadOnly"}, false},
		{"different privilege", map[string]string{"db1": "ReadOnly"}, map[string]string{"db1": "ReadWrite"}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := polarDbAccountDatabasePrivilegesMatch(tc.actual, tc.expected); got != tc.want {
				t.Fatalf("%s: want %v, got %v", tc.name, tc.want, got)
			}
		})
	}
}
