package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudRdsDBDatabaseMySQL(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_db_database.default"

	var dbDatabaseBasicMap = map[string]string{
		"instance_id":   CHECKSET,
		"character_set": "utf8",
		"description":   "",
	}

	ra := resourceAttrInit(resourceId, dbDatabaseBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBDatabase")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testaccdbdatabase%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBDatabaseConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":    "${alicloud_db_instance.default.id}",
					"data_base_name": "${var.name}",
					"character_set":  "utf8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":    CHECKSET,
						"data_base_name": CHECKSET,
						"character_set":  "utf8",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "from terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"description": "from terraform"}),
				),
			},
		},
	})

}

func TestAccAliCloudRdsDBDatabaseMySQL1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_db_database.default"

	var dbDatabaseBasicMap = map[string]string{
		"instance_id":   CHECKSET,
		"character_set": "utf8",
		"description":   "",
	}

	ra := resourceAttrInit(resourceId, dbDatabaseBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBDatabase")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testaccdbdatabase%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBDatabaseConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":   "${alicloud_db_instance.default.id}",
					"name":          "${var.name}",
					"character_set": "utf8",
					"description":   "from terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":   CHECKSET,
						"name":          CHECKSET,
						"character_set": "utf8",
						"description":   "from terraform",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudRdsDBDatabaseMySQL2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_db_database.default"

	var dbDatabaseBasicMap = map[string]string{
		"instance_id": CHECKSET,
	}

	ra := resourceAttrInit(resourceId, dbDatabaseBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBDatabase")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testaccdbdatabase%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBDatabaseConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":    "${alicloud_db_instance.default.id}",
					"data_base_name": "${var.name}",
					"character_set":  "gbk",
					"description":    "from terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":    CHECKSET,
						"data_base_name": CHECKSET,
						"character_set":  CHECKSET,
						"description":    "from terraform",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudRdsDBDatabaseMySQL3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_db_database.default"

	var dbDatabaseBasicMap = map[string]string{
		"instance_id": CHECKSET,
		"description": "",
	}

	ra := resourceAttrInit(resourceId, dbDatabaseBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBDatabase")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testaccdbdatabase%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBDatabaseConfigAccountDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":       "${alicloud_db_instance.default.id}",
					"data_base_name":    "${var.name}",
					"character_set":     "utf8mb4",
					"collation_name":    "utf8mb4_general_ci",
					"account_name":      "${alicloud_rds_account.default.account_name}",
					"account_privilege": "ReadWrite",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":       CHECKSET,
						"data_base_name":    CHECKSET,
						"character_set":     "utf8mb4",
						"collation_name":    "utf8mb4_general_ci",
						"account_name":      CHECKSET,
						"account_privilege": "ReadWrite",
					}),
					testAccCheckRdsDatabaseAccountPrivilege(resourceId, "ReadWrite"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_name", "account_privilege", "collation_name"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "from terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"description": "from terraform"}),
				),
			},
		},
	})
}

func TestAccAliCloudRdsDBDatabasePostgreSQL(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_db_database.default"

	var dbDatabaseBasicMap = map[string]string{
		"instance_id": CHECKSET,
		"description": "",
	}

	ra := resourceAttrInit(resourceId, dbDatabaseBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBDatabase")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testaccdbdatabase%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBDatabaseConfigPostgreSQLDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},

		// module name

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":    "${alicloud_db_instance.default.id}",
					"data_base_name": "${var.name}",
					"description":    "from terraform",
					"character_set":  "UTF8,C,en_US.utf8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":    CHECKSET,
						"data_base_name": CHECKSET,
						"character_set":  "UTF8,C,en_US.utf8",
						"description":    "from terraform",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudRdsDBDatabasePostgreSQL1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_db_database.default"

	var dbDatabaseBasicMap = map[string]string{
		"instance_id": CHECKSET,
		"description": "",
	}

	ra := resourceAttrInit(resourceId, dbDatabaseBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBDatabase")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testaccdbdatabase%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBDatabaseConfigPostgreSQLDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":   "${alicloud_db_instance.default.id}",
					"name":          "${var.name}",
					"character_set": "UTF8",
					"description":   "from terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":   CHECKSET,
						"name":          CHECKSET,
						"character_set": CHECKSET,
						"description":   "from terraform",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudRdsDBDatabasePostgreSQL2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_db_database.default"

	var dbDatabaseBasicMap = map[string]string{
		"instance_id": CHECKSET,
		"description": "",
	}

	ra := resourceAttrInit(resourceId, dbDatabaseBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBDatabase")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testaccdbdatabase%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBDatabaseConfigPostgreSQLDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},

		// module name

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":    "${alicloud_db_instance.default.id}",
					"data_base_name": "${var.name}",
					"character_set":  "UTF8,fo_FO.utf8,bo_CN",
					"description":    "from terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":    CHECKSET,
						"data_base_name": CHECKSET,
						"character_set":  CHECKSET,
						"description":    "from terraform",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudRdsDBDatabasePostgreSQL3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_db_database.default"

	var dbDatabaseBasicMap = map[string]string{
		"instance_id": CHECKSET,
	}

	ra := resourceAttrInit(resourceId, dbDatabaseBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBDatabase")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testaccdbdatabase%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBDatabaseConfigPostgreSQLDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},

		// module name

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":    "${alicloud_db_instance.default.id}",
					"data_base_name": "${var.name}",
					"character_set":  "KOI8U,ru_UA.koi8u,ru_UA.koi8u",
					"description":    "from terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":    CHECKSET,
						"data_base_name": CHECKSET,
						"character_set":  CHECKSET,
						"description":    "from terraform",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudRdsDBDatabaseSQLServer(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_db_database.default"

	var dbDatabaseBasicMap = map[string]string{
		"instance_id": CHECKSET,
		"description": "",
	}

	ra := resourceAttrInit(resourceId, dbDatabaseBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBDatabase")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testaccdbdatabase%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBDatabaseConfigSQLServerDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},

		// module name

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":    "${alicloud_db_instance.default.id}",
					"data_base_name": "${var.name}",
					"character_set":  "Chinese_PRC_CI_AS",
					"description":    "from terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":    CHECKSET,
						"data_base_name": CHECKSET,
						"character_set":  CHECKSET,
						"description":    "from terraform",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func resourceDBDatabaseConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "8.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones[0].id
	engine = "MySQL"
	engine_version = "8.0"
    category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids[0]
  zone_id = data.alicloud_db_zones.default.zones[0].id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids[0]
 zone_id = data.alicloud_db_zones.default.ids[0]
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_security_group" "default" {
	name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids[0]
}

resource "alicloud_db_instance" "default" {
    engine = "MySQL"
	engine_version = "8.0"
 	db_instance_storage_type = "cloud_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes[0].instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].min
	vswitch_id = local.vswitch_id
	instance_name = var.name
	instance_charge_type = "Postpaid"
}
`, name)
}

func testAccCheckRdsDatabaseAccountPrivilege(resourceId, accountPrivilege string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		accountRs, ok := s.RootModule().Resources["alicloud_rds_account.default"]
		if !ok {
			return fmt.Errorf("resource alicloud_rds_account.default not found in state")
		}
		accountName := accountRs.Primary.Attributes["account_name"]

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		parts := strings.Split(rs.Primary.ID, ":")
		if len(parts) != 2 {
			return fmt.Errorf("unexpected database id %s", rs.Primary.ID)
		}
		response, err := client.RpcPost("Rds", "2014-08-15", "DescribeDatabases", nil, map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": parts[0],
			"DBName":       parts[1],
		}, true)
		if err != nil {
			return err
		}
		accounts, err := jsonpath.Get("$.Databases.Database[0].Accounts.AccountPrivilegeInfo[*]", response)
		if err != nil {
			return err
		}
		for _, item := range accounts.([]interface{}) {
			account, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if account["Account"] == accountName && account["AccountPrivilege"] == accountPrivilege {
				return nil
			}
		}
		return fmt.Errorf("account %s with privilege %s not found on database %s", accountName, accountPrivilege, rs.Primary.ID)
	}
}

func resourceDBDatabaseConfigAccountDependence(name string) string {
	return resourceDBDatabaseConfigDependence(name) + `
resource "alicloud_rds_account" "default" {
  db_instance_id   = alicloud_db_instance.default.id
  account_name     = replace(var.name, "-", "_")
  account_password = "Test12345"
}
`
}

func resourceDBDatabaseConfigPostgreSQLDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "PostgreSQL"
	engine_version = "14.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids[0]
  zone_id = data.alicloud_db_zones.default.zones[0].id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids[0]
 zone_id = data.alicloud_db_zones.default.zones[0].id
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
}

// RDS PostgreSQL on ECS requires this service-linked role before CreateDBInstance
resource "alicloud_rds_service_linked_role" "default" {
	service_name = "AliyunServiceRoleForRdsPgsqlOnEcs"
}

resource "alicloud_db_instance" "default" {
    engine = "PostgreSQL"
	engine_version = "14.0"
 	db_instance_storage_type = "cloud_essd"
instance_type = "pg.n2.2c.2m"
	instance_storage = 40
	vswitch_id = local.vswitch_id
	instance_name = var.name
	instance_charge_type = "Postpaid"
	depends_on = [alicloud_rds_service_linked_role.default]
}
`, name)
}

func resourceDBDatabaseConfigSQLServerDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "SQLServer"
	engine_version = "2019_ent"
	instance_charge_type = "PostPaid"
	category = "AlwaysOn"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones[0].id
	engine = "SQLServer"
	engine_version = "2019_ent"
    category = "AlwaysOn"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
}
data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids[0]
  zone_id = data.alicloud_db_zones.default.zones[0].id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids[0]
 zone_id = data.alicloud_db_zones.default.zones[0].id
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
}

resource "alicloud_db_instance" "default" {
    engine = "SQLServer"
	engine_version = "2019_ent"
 	db_instance_storage_type = "cloud_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes[0].instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].min
	vswitch_id = local.vswitch_id
	instance_name = var.name
	category = "AlwaysOn"
	instance_charge_type = "Postpaid"
}
`, name)
}
