package polardb_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/polardb"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// PolarDB Database resource acceptance tests (product package layout). The
// shared ACC framework (attribute map checks, HCL rendering, destroy checks) is
// the exported alicloud/acctest package.

const polarDBDatabaseTestResourceId = "alicloud_polardb_database.default"

// polarDBDatabaseTestPreCheck pins the region used by the PolarDB tests.
func polarDBDatabaseTestPreCheck(t *testing.T) {
	acctest.PreCheck(t)
}

func polarDBDatabaseTestServiceFunc() func() interface{} {
	return func() interface{} {
		return polardb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
}

// polarDBDatabaseAccountRecord returns the Accounts.Account entry of the given
// account from a DescribeDatabases response (the first entry when accountName is
// empty, mirroring the provider read semantics).
func polarDBDatabaseAccountRecord(object map[string]interface{}, accountName string) map[string]interface{} {
	accounts, _ := object["Accounts"].(map[string]interface{})
	accountList, _ := accounts["Account"].([]interface{})
	for _, item := range accountList {
		record, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if accountName == "" || fmt.Sprint(record["AccountName"]) == accountName {
			return record
		}
	}
	return nil
}

// checkDatabaseAccountPrivilege verifies through DescribeDatabases that the
// privilege reported for the expected account equals the expected value. The
// account is matched by name rather than by list position.
func checkDatabaseAccountPrivilege(resourceId, accountName, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("can't find resource by id: %s", resourceId)
		}
		object, err := polardb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient)).DescribeDatabase(rs.Primary.ID)
		if err != nil {
			return err
		}
		record := polarDBDatabaseAccountRecord(object, accountName)
		if record == nil {
			return fmt.Errorf("account %q not returned by DescribeDatabases for %s; Accounts=%v", accountName, rs.Primary.ID, object["Accounts"])
		}
		if got := fmt.Sprint(record["AccountPrivilege"]); got != expected {
			return fmt.Errorf("expected account privilege %q for account %q, got %q for %s; Accounts=%v", expected, accountName, got, rs.Primary.ID, object["Accounts"])
		}
		return nil
	}
}

// checkDatabaseAccountName verifies through DescribeDatabases that the database
// is bound to the expected account.
func checkDatabaseAccountName(resourceId, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("can't find resource by id: %s", resourceId)
		}
		object, err := polardb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient)).DescribeDatabase(rs.Primary.ID)
		if err != nil {
			return err
		}
		record := polarDBDatabaseAccountRecord(object, "")
		if record == nil {
			return fmt.Errorf("no accounts returned for %s; Accounts=%v", rs.Primary.ID, object["Accounts"])
		}
		if got := fmt.Sprint(record["AccountName"]); got != expected {
			return fmt.Errorf("expected account name %q, got %q for %s; Accounts=%v", expected, got, rs.Primary.ID, object["Accounts"])
		}
		return nil
	}
}

// checkDatabaseAccountEmpowered verifies that DescribeDatabases reports the
// database account as Empowered. PostgreSQL/Oracle do not return an
// AccountPrivilege value for the database owner (the owner holds implicit
// privileges), so the grant status is the observable API fact there.
func checkDatabaseAccountEmpowered(resourceId, accountName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("can't find resource by id: %s", resourceId)
		}
		object, err := polardb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient)).DescribeDatabase(rs.Primary.ID)
		if err != nil {
			return err
		}
		record := polarDBDatabaseAccountRecord(object, accountName)
		if record == nil {
			return fmt.Errorf("account %q not returned by DescribeDatabases for %s; Accounts=%v", accountName, rs.Primary.ID, object["Accounts"])
		}
		if got := fmt.Sprint(record["PrivilegeStatus"]); got != "Empowered" {
			return fmt.Errorf("expected privilege status Empowered for account %q, got %q for %s; Accounts=%v", accountName, got, rs.Primary.ID, object["Accounts"])
		}
		return nil
	}
}

var aliCloudPolarDbDatabaseMapBase = map[string]string{
	"character_set_name": acctest.CHECKSET,
	"status":             acctest.CHECKSET,
}

var aliCloudPolarDbDatabaseMap9535 = map[string]string{
	"character_set_name": acctest.CHECKSET,
	"db_description":     acctest.CHECKSET,
	"status":             acctest.CHECKSET,
}

func polarDBDatabaseCheck(resourceId string, checkMap map[string]string) acctest.ResourceAttrMapUpdate {
	rac := acctest.ResourceAttrCheckInit(
		acctest.ResourceCheckInitWithDescribeMethod(resourceId, new(map[string]interface{}), polarDBDatabaseTestServiceFunc(), "DescribeDatabase"),
		acctest.ResourceAttrInit(resourceId, checkMap),
	)
	return rac.ResourceAttrMapUpdateSet()
}

func polarDBDatabaseDestroyCheck(resourceId string, checkMap map[string]string) resource.TestCheckFunc {
	rac := acctest.ResourceAttrCheckInit(
		acctest.ResourceCheckInitWithDescribeMethod(resourceId, new(map[string]interface{}), polarDBDatabaseTestServiceFunc(), "DescribeDatabase"),
		acctest.ResourceAttrInit(resourceId, checkMap),
	)
	return rac.CheckResourceDestroy()
}

// TestAccAliCloudPolarDbDatabase_mysqlLifecycle11854 covers the baseline MySQL
// lifecycle: create, db_description update and import.
func TestAccAliCloudPolarDbDatabase_mysqlLifecycle11854(t *testing.T) {
	resourceId := polarDBDatabaseTestResourceId
	testAccCheck := polarDBDatabaseCheck(resourceId, aliCloudPolarDbDatabaseMapBase)
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, aliCloudPolarDbDatabaseMySQLDependence)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { polarDBDatabaseTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             polarDBDatabaseDestroyCheck(resourceId, aliCloudPolarDbDatabaseMapBase),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id": "${alicloud_polardb_account.default.db_cluster_id}",
					"db_name":       name,
				}),
				Check: testAccCheck(map[string]string{
					"db_cluster_id": acctest.CHECKSET,
					"db_name":       name,
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_description": name,
				}),
				Check: testAccCheck(map[string]string{
					"db_description": name,
				}),
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

// TestAccAliCloudPolarDbDatabase_mysqlFullAttributes11854Twin covers a MySQL
// database created with every optional argument set explicitly.
func TestAccAliCloudPolarDbDatabase_mysqlFullAttributes11854Twin(t *testing.T) {
	resourceId := polarDBDatabaseTestResourceId
	testAccCheck := polarDBDatabaseCheck(resourceId, aliCloudPolarDbDatabaseMapBase)
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, aliCloudPolarDbDatabaseMySQLDependence)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { polarDBDatabaseTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             polarDBDatabaseDestroyCheck(resourceId, aliCloudPolarDbDatabaseMapBase),
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
				Check: testAccCheck(map[string]string{
					"db_cluster_id":      acctest.CHECKSET,
					"db_name":            name,
					"character_set_name": "utf8mb4",
					"account_name":       acctest.CHECKSET,
					"db_description":     name,
				}),
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

// TestAccAliCloudPolarDbDatabase_postgreSQLLifecycle9535 covers the PostgreSQL
// lifecycle including the account_name update, which is applied through
// GrantAccountPrivilege with the DBOwner constant.
func TestAccAliCloudPolarDbDatabase_postgreSQLLifecycle9535(t *testing.T) {
	resourceId := polarDBDatabaseTestResourceId
	testAccCheck := polarDBDatabaseCheck(resourceId, aliCloudPolarDbDatabaseMapBase)
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, aliCloudPolarDbDatabasePostgreSQLDependence)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { polarDBDatabaseTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             polarDBDatabaseDestroyCheck(resourceId, aliCloudPolarDbDatabaseMapBase),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id": "${alicloud_polardb_account.default.db_cluster_id}",
					"db_name":       name,
					"account_name":  "${alicloud_polardb_account.default.account_name}",
				}),
				Check: testAccCheck(map[string]string{
					"db_cluster_id": acctest.CHECKSET,
					"db_name":       name,
					"account_name":  acctest.CHECKSET,
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_name": "${alicloud_polardb_account.update.account_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_name": name + "update",
					}),
					checkDatabaseAccountName(resourceId, name+"update"),
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

// TestAccAliCloudPolarDbDatabase_postgreSQLFullAttributes9535Twin covers a
// PostgreSQL database created with explicit character set, collate and ctype.
func TestAccAliCloudPolarDbDatabase_postgreSQLFullAttributes9535Twin(t *testing.T) {
	resourceId := polarDBDatabaseTestResourceId
	testAccCheck := polarDBDatabaseCheck(resourceId, aliCloudPolarDbDatabaseMap9535)
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, aliCloudPolarDbDatabasePostgreSQLDependence)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { polarDBDatabaseTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             polarDBDatabaseDestroyCheck(resourceId, aliCloudPolarDbDatabaseMap9535),
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
				Check: testAccCheck(map[string]string{
					"db_cluster_id":      acctest.CHECKSET,
					"db_name":            name,
					"character_set_name": "ISO_8859_5",
					"account_name":       acctest.CHECKSET,
					"db_description":     name,
				}),
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

// TestAccAliCloudPolarDbDatabase_accountPrivilegeMySQL covers account_privilege
// on a MySQL cluster: the value is sent by CreateDatabase, verified against
// DescribeDatabases, and read back into state / import.
func TestAccAliCloudPolarDbDatabase_accountPrivilegeMySQL(t *testing.T) {
	resourceId := polarDBDatabaseTestResourceId
	testAccCheck := polarDBDatabaseCheck(resourceId, aliCloudPolarDbDatabaseMapBase)
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, aliCloudPolarDbDatabaseMySQLDependence)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { polarDBDatabaseTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             polarDBDatabaseDestroyCheck(resourceId, aliCloudPolarDbDatabaseMapBase),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":     "${alicloud_polardb_account.default.db_cluster_id}",
					"db_name":           name,
					"account_name":      "${alicloud_polardb_account.default.account_name}",
					"account_privilege": "ReadOnly",
					"db_description":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id":     acctest.CHECKSET,
						"db_name":           name,
						"account_name":      acctest.CHECKSET,
						"account_privilege": "ReadOnly",
						"db_description":    name,
					}),
					checkDatabaseAccountPrivilege(resourceId, name, "ReadOnly"),
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

// TestAccAliCloudPolarDbDatabase_accountPrivilegePostgreSQL covers
// account_privilege on a PostgreSQL cluster. The database is created without
// collate/ctype, which exercises the PostgreSQL defaults ("C"). PostgreSQL does
// not report an AccountPrivilege value for the database owner, so the API fact
// asserted here is the grant status while the state keeps the configured value.
// The import proves the value is unrecoverable, and a plan-only step proves that
// keeping it in the configuration plans a replacement.
func TestAccAliCloudPolarDbDatabase_accountPrivilegePostgreSQL(t *testing.T) {
	resourceId := polarDBDatabaseTestResourceId
	testAccCheck := polarDBDatabaseCheck(resourceId, aliCloudPolarDbDatabaseMapBase)
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, aliCloudPolarDbDatabasePostgreSQLDependence)
	privilegeConfig := map[string]interface{}{
		"db_cluster_id":     "${alicloud_polardb_account.default.db_cluster_id}",
		"db_name":           name,
		"account_name":      "${alicloud_polardb_account.default.account_name}",
		"account_privilege": "ReadOnly",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { polarDBDatabaseTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             polarDBDatabaseDestroyCheck(resourceId, aliCloudPolarDbDatabaseMapBase),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(privilegeConfig),
				Check: resource.ComposeTestCheckFunc(
					checkDatabaseAccountEmpowered(resourceId, name),
					testAccCheck(map[string]string{
						"db_cluster_id":     acctest.CHECKSET,
						"db_name":           name,
						"account_name":      acctest.CHECKSET,
						"account_privilege": "ReadOnly",
					}),
				),
			},
			{
				// The API does not return AccountPrivilege for the PostgreSQL
				// owner, so an import cannot recover the configured value.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected 1 imported state, got %d", len(states))
					}
					attributes := states[0].Attributes
					if got := attributes["account_name"]; got != name {
						return fmt.Errorf("expected imported account_name %q, got %q", name, got)
					}
					if got := attributes["account_privilege"]; got != "" {
						return fmt.Errorf("expected account_privilege to be unrecoverable after import, got %q", got)
					}
					return nil
				},
			},
			{
				// Keeping account_privilege in the configuration after the import
				// must replace this database (ForceNew). The replacement action is
				// asserted on the plan before it is applied, and the recreated
				// database is then verified.
				Config: testAccConfig(privilegeConfig),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					checkDatabaseAccountEmpowered(resourceId, name),
					testAccCheck(map[string]string{
						"account_privilege": "ReadOnly",
					}),
				),
			},
		},
	})
}

// TestAccAliCloudPolarDbDatabase_oracleAccountPrivilege covers the
// Oracle-compatible engine: the cluster node class and zone are resolved from the
// node-class capability query for Oracle 14, and the test verifies database
// creation without collate/ctype (default "C"), the account_privilege grant and
// the account_name update served by GrantAccountPrivilege.
func TestAccAliCloudPolarDbDatabase_oracleAccountPrivilege(t *testing.T) {
	resourceId := polarDBDatabaseTestResourceId
	testAccCheck := polarDBDatabaseCheck(resourceId, aliCloudPolarDbDatabaseMapBase)
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, aliCloudPolarDbDatabaseOracleDependence)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { polarDBDatabaseTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             polarDBDatabaseDestroyCheck(resourceId, aliCloudPolarDbDatabaseMapBase),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":     "${alicloud_polardb_account.default.db_cluster_id}",
					"db_name":           name,
					"account_name":      "${alicloud_polardb_account.default.account_name}",
					"account_privilege": "ReadOnly",
				}),
				Check: resource.ComposeTestCheckFunc(
					checkDatabaseAccountEmpowered(resourceId, name),
					testAccCheck(map[string]string{
						"db_cluster_id":     acctest.CHECKSET,
						"db_name":           name,
						"account_name":      acctest.CHECKSET,
						"account_privilege": "ReadOnly",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_name": "${alicloud_polardb_account.update.account_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_name": name + "update",
					}),
					checkDatabaseAccountName(resourceId, name+"update"),
				),
			},
		},
	})
}

// TestAccAliCloudPolarDbDatabase_accountPrivilegeRequiresAccountName verifies
// that configuring account_privilege without account_name is rejected during
// plan: CreateDatabase only applies AccountPrivilege together with AccountName,
// so accepting the combination would silently ignore it.
func TestAccAliCloudPolarDbDatabase_accountPrivilegeRequiresAccountName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { polarDBDatabaseTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
resource "alicloud_polardb_database" "default" {
  db_cluster_id     = "pc-plan-check-only"
  db_name           = "tfacc_privilege_without_account"
  account_privilege = "ReadOnly"
}
`,
				ExpectError: regexp.MustCompile(`"account_privilege" requires "account_name"`),
				// The step fails during plan, so no resource state exists; the no-op
				// check keeps the step shape consistent for the coverage parser.
				Check: func(_ *terraform.State) error { return nil },
			},
		},
	})
}

// aliCloudPolarDbDatabaseMySQLDependence creates a self-contained MySQL cluster
// (VPC, vswitch, cluster and account) so the case does not depend on
// pre-existing account fixtures.
func aliCloudPolarDbDatabaseMySQLDependence(name string) string {
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
  vpc_name   = "tf-testAcc-polardb-db-${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_polardb_node_classes.default.classes.0.zone_id
  vswitch_name = "tf-testAcc-polardb-db-${var.name}"
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "8.0"
  db_node_class = "polar.mysql.x4.medium"
  pay_type      = "PostPaid"
  vswitch_id    = alicloud_vswitch.default.id
  # encrypt_new_tables is read back as "OFF" by the cluster resource (TDE default)
  # while the schema has no Computed flag, so an unconfigured value drifts to null
  # and produces a non-empty plan. Pre-existing cluster-resource behaviour,
  # unrelated to the database resource under test.
  lifecycle {
    ignore_changes = [encrypt_new_tables]
  }
}

resource "alicloud_polardb_account" "default" {
  db_cluster_id    = alicloud_polardb_cluster.default.id
  account_name     = var.name
  account_password = "Example1234"
}
`, name)
}

// aliCloudPolarDbDatabasePostgreSQLDependence creates a self-contained
// PostgreSQL cluster plus two accounts (the second one is used by the
// account_name update step).
func aliCloudPolarDbDatabasePostgreSQLDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_polardb_node_classes" "default" {
  db_type    = "PostgreSQL"
  db_version = "14"
  pay_type   = "PostPaid"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "tf-testAcc-polardb-db-pg-${var.name}"
  cidr_block = "172.17.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.17.0.0/24"
  zone_id      = data.alicloud_polardb_node_classes.default.classes.0.zone_id
  vswitch_name = "tf-testAcc-polardb-db-pg-${var.name}"
}

resource "alicloud_polardb_cluster" "default" {
  db_version    = "14"
  pay_type      = "PostPaid"
  db_node_class = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
  db_type       = "PostgreSQL"
  vswitch_id    = alicloud_vswitch.default.id
  lifecycle {
    ignore_changes = [encrypt_new_tables]
  }
}

resource "alicloud_polardb_account" "default" {
  db_cluster_id    = alicloud_polardb_cluster.default.id
  account_name     = var.name
  account_password = "Example1234"
}

resource "alicloud_polardb_account" "update" {
  db_cluster_id    = alicloud_polardb_cluster.default.id
  account_name     = "${var.name}update"
  account_password = "Example1234"
}
`, name)
}

// aliCloudPolarDbDatabaseOracleDependence creates a self-contained
// Oracle-compatible cluster. The node class and zone are resolved from the
// node-class capability query for the Oracle 14 engine instead of being
// hardcoded, so the fixture adapts to the capabilities of the account/region.
func aliCloudPolarDbDatabaseOracleDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_polardb_node_classes" "default" {
  db_type    = "Oracle"
  db_version = "14"
  pay_type   = "PostPaid"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "tf-testAcc-polardb-db-oracle-${var.name}"
  cidr_block = "172.18.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.18.0.0/24"
  zone_id      = data.alicloud_polardb_node_classes.default.classes.0.zone_id
  vswitch_name = "tf-testAcc-polardb-db-oracle-${var.name}"
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "Oracle"
  db_version    = "14"
  db_node_class = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
  db_node_count = "2"
  pay_type      = "PostPaid"
  vswitch_id    = alicloud_vswitch.default.id
  encrypt_new_tables = "OFF"
  lifecycle {
    ignore_changes = [encrypt_new_tables]
  }
}

resource "alicloud_polardb_account" "default" {
  db_cluster_id    = alicloud_polardb_cluster.default.id
  account_name     = var.name
  account_password = "Example1234"
}

resource "alicloud_polardb_account" "update" {
  db_cluster_id    = alicloud_polardb_cluster.default.id
  account_name     = "${var.name}update"
  account_password = "Example1234"
}
`, name)
}
