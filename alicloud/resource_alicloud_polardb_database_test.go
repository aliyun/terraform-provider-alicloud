package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_polardb_database", &resource.Sweeper{
		Name: "alicloud_polardb_database",
		F:    testSweepPolarDBDatabases,
	})
}

// testSweepPolarDBDatabases walks every PolarDB cluster in the region and
// deletes the databases on it. A cluster whose description carries a
// do-not-delete marker (protectedSweepVpcName) is skipped entirely, as is any
// database whose name carries such a marker; both checks run in every mode so
// SWEEP_ALL cannot bypass them. When SWEEP_ALL is unset, only databases on
// clusters (and with names) matching the standard test prefixes are touched.
// Per-instance failures are logged and non-fatal so one bad cluster cannot
// block the rest of the sweep.
func testSweepPolarDBDatabases(region string) error {
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
			log.Printf("[INFO] Skipping PolarDB databases on NODELETING cluster: %s (%s)", clusterName, clusterId)
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

		dreq := polardb.CreateDescribeDatabasesRequest()
		dreq.RegionId = client.RegionId
		dreq.DBClusterId = clusterId
		raw, err := client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
			return polardbClient.DescribeDatabases(dreq)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to list PolarDB databases on cluster %s: %s", clusterId, err)
			continue
		}
		dresp, _ := raw.(*polardb.DescribeDatabasesResponse)
		if dresp == nil {
			continue
		}
		for _, db := range dresp.Databases.Database {
			if protectedSweepVpcName(db.DBName) {
				log.Printf("[INFO] Skipping PolarDB database (NODELETING): %s on %s", db.DBName, clusterId)
				continue
			}
			if !sweepAll() {
				match := false
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(db.DBName), strings.ToLower(prefix)) {
						match = true
						break
					}
				}
				if !match {
					continue
				}
			}
			log.Printf("[INFO] Deleting PolarDB database: %s on cluster %s", db.DBName, clusterId)
			delReq := polardb.CreateDeleteDatabaseRequest()
			delReq.RegionId = client.RegionId
			delReq.DBClusterId = clusterId
			delReq.DBName = db.DBName
			if _, err := client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
				return polardbClient.DeleteDatabase(delReq)
			}); err != nil {
				log.Printf("[ERROR] Failed to delete PolarDB database %s on %s: %s", db.DBName, clusterId, err)
			}
		}
	}
	return nil
}

// Test PolarDb Database. >>> Resource test cases, automatically generated.
// Case Database用例_mysql 11854
func TestAccAliCloudPolarDbDatabase_basic11854(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_database.default"
	ra := resourceAttrInit(resourceId, AliCloudPolarDbDatabaseMap11854)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PolarDbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePolarDbDatabase")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudPolarDbDatabaseBasicDependence11854)
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
					"db_cluster_id": "${alicloud_polardb_account.default.db_cluster_id}",
					"db_name":       name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id": CHECKSET,
						"db_name":       name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_description": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"collate", "ctype"},
			},
		},
	})
}

func TestAccAliCloudPolarDbDatabase_basic11854_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_database.default"
	ra := resourceAttrInit(resourceId, AliCloudPolarDbDatabaseMap11854)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PolarDbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePolarDbDatabase")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudPolarDbDatabaseBasicDependence11854)
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
					"db_cluster_id":      "${alicloud_polardb_account.default.db_cluster_id}",
					"db_name":            name,
					"character_set_name": "utf8mb4",
					"account_name":       "${alicloud_polardb_account.default.account_name}",
					"db_description":     name,
					"collate":            "utf8mb4_bin",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id":      CHECKSET,
						"db_name":            name,
						"character_set_name": "utf8mb4",
						"account_name":       CHECKSET,
						"db_description":     name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"collate", "ctype"},
			},
		},
	})
}

var AliCloudPolarDbDatabaseMap11854 = map[string]string{
	"character_set_name": CHECKSET,
	"status":             CHECKSET,
}

func AliCloudPolarDbDatabaseBasicDependence11854(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_polardb_zones" "default" {}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_polardb_zones.default.ids[length(data.alicloud_polardb_zones.default.ids)-1]
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "8.0"
  db_node_class = "polar.mysql.x4.medium"
  pay_type      = "PostPaid"
  vswitch_id    = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_polardb_account" "default" {
  db_cluster_id       = alicloud_polardb_cluster.default.id
  account_name        = var.name
  account_password    = "Example1234"
}
`, name)
}

// Case  Database用例1 9535
func TestAccAliCloudPolarDbDatabase_basic9535(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_database.default"
	ra := resourceAttrInit(resourceId, AliCloudPolarDbDatabaseMap9535)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PolarDbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePolarDbDatabase")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudPolarDbDatabaseBasicDependence9535)
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
					"db_cluster_id": "${alicloud_polardb_account.default.db_cluster_id}",
					"db_name":       name,
					"account_name":  "${alicloud_polardb_account.default.account_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id": CHECKSET,
						"db_name":       name,
						"account_name":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_name": "${alicloud_polardb_account.update.account_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_name": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"collate", "ctype"},
			},
		},
	})
}

func TestAccAliCloudPolarDbDatabase_basic9535_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_database.default"
	ra := resourceAttrInit(resourceId, AliCloudPolarDbDatabaseMap9535)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PolarDbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePolarDbDatabase")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudPolarDbDatabaseBasicDependence9535)
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
					"db_cluster_id":      "${alicloud_polardb_account.default.db_cluster_id}",
					"db_name":            name,
					"character_set_name": "ISO_8859_5",
					"account_name":       "${alicloud_polardb_account.default.account_name}",
					"db_description":     name,
					"collate":            "POSIX",
					"ctype":              "POSIX",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id":      CHECKSET,
						"db_name":            name,
						"character_set_name": "ISO_8859_5",
						"account_name":       CHECKSET,
						"db_description":     name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"collate", "ctype"},
			},
		},
	})
}

var AliCloudPolarDbDatabaseMap9535 = map[string]string{
	"character_set_name": CHECKSET,
	"db_description":     CHECKSET,
	"status":             CHECKSET,
}

func AliCloudPolarDbDatabaseBasicDependence9535(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_polardb_node_classes" "default" {
  db_type  = "PostgreSQL"
  pay_type = "PostPaid"
  category = "Normal"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_polardb_node_classes.default.classes.0.zone_id
}

resource "alicloud_polardb_cluster" "default" {
  db_version    = "14"
  pay_type      = "PostPaid"
  db_node_class = "polar.pg.x4.medium"
  db_type       = "PostgreSQL"
  vswitch_id    = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_polardb_account" "default" {
  db_cluster_id       = alicloud_polardb_cluster.default.id
  account_name        = var.name
  account_password    = "Example1234"
}

resource "alicloud_polardb_account" "update" {
  db_cluster_id       = alicloud_polardb_cluster.default.id
  account_name        = "${var.name}update"
  account_password    = "Example1234"
}
`, name)
}

// Test PolarDb Database. <<< Resource test cases, automatically generated.
