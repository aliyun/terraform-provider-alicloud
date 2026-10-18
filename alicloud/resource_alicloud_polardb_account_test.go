package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_polardb_account", &resource.Sweeper{
		Name: "alicloud_polardb_account",
		F:    testSweepPolarDBAccounts,
	})
}

// testSweepPolarDBAccounts walks every PolarDB cluster in the region and
// deletes the accounts on it. A cluster whose description carries a
// do-not-delete marker (protectedSweepVpcName) is skipped entirely, as is any
// account whose name carries such a marker; both checks run in every mode so
// SWEEP_ALL cannot bypass them. When SWEEP_ALL is unset, only accounts on
// clusters (and with names) matching the standard test prefixes are touched.
// DynamoDB-type accounts are skipped because DeleteAccount does not support
// them (see resourceAlicloudPolarDbAccountDelete). Per-instance failures are
// logged and non-fatal so one bad cluster cannot block the rest of the sweep.
func testSweepPolarDBAccounts(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{"tf-testAcc", "tf_testAcc"}

	clusters, err := listPolarDBClustersForSweep(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list PolarDB clusters: %s", err)
		return nil
	}

	for _, c := range clusters {
		clusterName := c.DBClusterDescription
		clusterId := c.DBClusterId
		if protectedSweepVpcName(clusterName) {
			log.Printf("[INFO] Skipping PolarDB accounts on NODELETING cluster: %s (%s)", clusterName, clusterId)
			continue
		}
		if !sweepAll() {
			match := false
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(clusterName), strings.ToLower(prefix)) {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}

		areq := polardb.CreateDescribeAccountsRequest()
		areq.RegionId = client.RegionId
		areq.DBClusterId = clusterId
		raw, err := client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
			return polardbClient.DescribeAccounts(areq)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to list PolarDB accounts on cluster %s: %s", clusterId, err)
			continue
		}
		aresp, _ := raw.(*polardb.DescribeAccountsResponse)
		if aresp == nil {
			continue
		}
		for _, acc := range aresp.Accounts {
			// DynamoDB-type accounts cannot be deleted via DeleteAccount.
			if acc.AccountType == "DynamoDB" {
				continue
			}
			if protectedSweepVpcName(acc.AccountName) {
				log.Printf("[INFO] Skipping PolarDB account (NODELETING): %s on %s", acc.AccountName, clusterId)
				continue
			}
			if !sweepAll() {
				match := false
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(acc.AccountName), strings.ToLower(prefix)) {
						match = true
						break
					}
				}
				if !match {
					continue
				}
			}
			log.Printf("[INFO] Deleting PolarDB account: %s on cluster %s", acc.AccountName, clusterId)
			dreq := polardb.CreateDeleteAccountRequest()
			dreq.RegionId = client.RegionId
			dreq.DBClusterId = clusterId
			dreq.AccountName = acc.AccountName
			if _, err := client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
				return polardbClient.DeleteAccount(dreq)
			}); err != nil {
				log.Printf("[ERROR] Failed to delete PolarDB account %s on %s: %s", acc.AccountName, clusterId, err)
			}
		}
	}
	return nil
}

// Test PolarDb Account. >>> Resource test cases, automatically generated.
// Case  Account用例_PG_高权限账号 11819
func TestAccAliCloudPolarDbAccount_basic11819(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_account.default"
	ra := resourceAttrInit(resourceId, AliCloudPolarDbAccountMap11819)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PolarDbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePolarDbAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudPolarDbAccountBasicDependence11819)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"account_password"},
			},
		},
	})
}

func TestAccAliCloudPolarDbAccount_basic11819_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_account.default"
	ra := resourceAttrInit(resourceId, AliCloudPolarDbAccountMap11819)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PolarDbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePolarDbAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudPolarDbAccountBasicDependence11819)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"account_password"},
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
	ra := resourceAttrInit(resourceId, AliCloudPolarDbAccountMap9211)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PolarDbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePolarDbAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudPolarDbAccountBasicDependence9211)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"kms_encrypted_password", "kms_encryption_context"},
			},
		},
	})
}

func TestAccAliCloudPolarDbAccount_basic9211_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_account.default"
	ra := resourceAttrInit(resourceId, AliCloudPolarDbAccountMap9211)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PolarDbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePolarDbAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudPolarDbAccountBasicDependence9211)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"kms_encrypted_password", "kms_encryption_context"},
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
	ra := resourceAttrInit(resourceId, AliCloudPolarDbAccountMapDynamoDB)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PolarDbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePolarDbAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudPolarDbAccountBasicDependenceDynamoDB)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
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
				ImportStateVerifyIgnore: []string{"account_password", "dynamodb_auth_password"},
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
	ra := resourceAttrInit(resourceId, AliCloudPolarDbAccountMap11819)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PolarDbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePolarDbAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudPolarDbAccountBasicDependence11819)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
					"account_password_wo":         REMOVEKEY,
					"account_password_wo_version": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password":            "YourPassword123!",
						"account_password_wo_version": REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "account_password_wo_version"},
			},
		},
	})
}
