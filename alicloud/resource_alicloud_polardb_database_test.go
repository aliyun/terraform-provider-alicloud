package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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

data "alicloud_polardb_node_classes" "default" {
  db_type    = "MySQL"
  db_version = "8.0"
  pay_type   = "PostPaid"
  category   = "Normal"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_polardb_node_classes.default.classes.0.zone_id
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "8.0"
  db_node_class = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
