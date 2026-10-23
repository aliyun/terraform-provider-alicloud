package polardb_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	polardb "github.com/aliyun/terraform-provider-alicloud/alicloud/service/polardb"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Test PolarDb Database. >>> Resource test cases, automatically generated.
// Case Database用例_mysql 11854
func TestAccAliCloudPolarDbDatabase_basic11854(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_database.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudPolarDbDatabaseMap11854)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolarDbDatabase")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudPolarDbDatabaseBasicDependence11854)
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
				ImportStateVerifyIgnore: []string{"collate", "ctype", "account_name", "account_privilege"},
			},
		},
	})
}

func TestAccAliCloudPolarDbDatabase_basic11854_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_database.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudPolarDbDatabaseMap11854)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolarDbDatabase")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudPolarDbDatabaseBasicDependence11854)
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
				ImportStateVerifyIgnore: []string{"collate", "ctype", "account_name", "account_privilege"},
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

resource "alicloud_vpc" "default" {
  vpc_name   = "tf-testAcc-polardb-db-basic-${var.name}"
  cidr_block = "172.19.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.19.0.0/24"
  zone_id      = data.alicloud_polardb_node_classes.default.classes.0.zone_id
  vswitch_name = "tf-testAcc-polardb-db-basic-${var.name}"
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "8.0"
  db_node_class = "polar.mysql.x4.medium"
  pay_type      = "PostPaid"
  vswitch_id    = alicloud_vswitch.default.id
  lifecycle {
    ignore_changes = [encrypt_new_tables]
  }
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudPolarDbDatabaseMap9535)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolarDbDatabase")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudPolarDbDatabaseBasicDependence9535)
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
						"account_name": name + "update",
					}),
					checkDatabaseAccountName(resourceId, name+"update"),
					checkDatabaseAccountEmpowered(resourceId, name+"update"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"collate", "ctype", "account_name", "account_privilege"},
			},
		},
	})
}

func TestAccAliCloudPolarDbDatabase_basic9535_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_database.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudPolarDbDatabaseMap9535)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolarDbDatabase")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudPolarDbDatabaseBasicDependence9535)
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
				ImportStateVerifyIgnore: []string{"collate", "ctype", "account_name", "account_privilege"},
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

// PolarDB database and account names allow only lower-case letters, digits and
// underscores, and account names are limited to 2-16 characters. The standard
// tf-testAcc prefix therefore cannot be used literally: new cases use the
// compatible lower-case form tf_testacc and keep derived account names short.
const polarDBDatabaseTestResourceId = "alicloud_polardb_database.default"

func polarDBDatabaseTestPreCheck(t *testing.T) {
	acctest.PreCheck(t)
}

func polarDBDatabaseTestServiceFunc() func() interface{} {
	return func() interface{} {
		return polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
}

// polarDBDatabaseAccountRecord returns the Accounts.Account entry of the given
// account from a DescribeDatabases response. The account is matched by name
// only: the API does not guarantee list order, so a first-entry fallback could
// observe another account's records.
func polarDBDatabaseAccountRecord(object map[string]interface{}, accountName string) map[string]interface{} {
	if accountName == "" {
		return nil
	}
	accounts, _ := object["Accounts"].(map[string]interface{})
	accountList, _ := accounts["Account"].([]interface{})
	for _, item := range accountList {
		record, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if fmt.Sprint(record["AccountName"]) == accountName {
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
		object, err := polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient)).DescribePolarDbDatabase(rs.Primary.ID)
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
		object, err := polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient)).DescribePolarDbDatabase(rs.Primary.ID)
		if err != nil {
			return err
		}
		record := polarDBDatabaseAccountRecord(object, expected)
		if record == nil {
			return fmt.Errorf("account %q not returned by DescribeDatabases for %s; Accounts=%v", expected, rs.Primary.ID, object["Accounts"])
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
		object, err := polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient)).DescribePolarDbDatabase(rs.Primary.ID)
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

func polarDBDatabaseCheck(resourceId string, checkMap map[string]string) acctest.ResourceAttrMapUpdate {
	rac := acctest.ResourceAttrCheckInit(
		acctest.ResourceCheckInitWithDescribeMethod(resourceId, new(map[string]interface{}), polarDBDatabaseTestServiceFunc(), "DescribePolarDbDatabase"),
		acctest.ResourceAttrInit(resourceId, checkMap),
	)
	return rac.ResourceAttrMapUpdateSet()
}

func polarDBDatabaseDestroyCheck(resourceId string, checkMap map[string]string) resource.TestCheckFunc {
	rac := acctest.ResourceAttrCheckInit(
		acctest.ResourceCheckInitWithDescribeMethod(resourceId, new(map[string]interface{}), polarDBDatabaseTestServiceFunc(), "DescribePolarDbDatabase"),
		acctest.ResourceAttrInit(resourceId, checkMap),
	)
	return rac.CheckResourceDestroy()
}

// grantDatabaseAccountPrivilege grants a privilege to an additional account of
// the database managed by the test, so a refresh observes a multi-account
// DescribeDatabases response.
func grantDatabaseAccountPrivilege(resourceId, accountName, privilege string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("can't find resource by id: %s", resourceId)
		}
		request := map[string]interface{}{
			"DBClusterId":      rs.Primary.Attributes["db_cluster_id"],
			"DBName":           rs.Primary.Attributes["db_name"],
			"AccountName":      accountName,
			"AccountPrivilege": privilege,
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		_, err := client.RpcPost("polardb", "2017-08-01", "GrantAccountPrivilege", map[string]interface{}{}, request, true)
		return err
	}
}

// revokeDatabaseAccountPrivilege removes the configured account's privilege so
// the next refresh can prove the drift is surfaced instead of keeping a stale
// privilege bound to another account.
func revokeDatabaseAccountPrivilege(resourceId, accountName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("can't find resource by id: %s", resourceId)
		}
		request := map[string]interface{}{
			"DBClusterId": rs.Primary.Attributes["db_cluster_id"],
			"DBName":      rs.Primary.Attributes["db_name"],
			"AccountName": accountName,
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		_, err := client.RpcPost("polardb", "2017-08-01", "RevokeAccountPrivilege", map[string]interface{}{}, request, true)
		return err
	}
}

// waitDatabaseAccountPrivilege polls DescribeDatabases until the account reports
// the expected privilege (the grant/revoke APIs are asynchronous).
func waitDatabaseAccountPrivilege(resourceId, accountName, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("can't find resource by id: %s", resourceId)
		}
		deadline := time.Now().Add(3 * time.Minute)
		last := "<unread>"
		for {
			object, err := polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient)).DescribePolarDbDatabase(rs.Primary.ID)
			if err == nil {
				if record := polarDBDatabaseAccountRecord(object, accountName); record != nil {
					last = fmt.Sprint(record["AccountPrivilege"])
					if last == expected {
						return nil
					}
				} else {
					last = "<absent>"
				}
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("timed out waiting for account %q privilege %q, last %q for %s", accountName, expected, last, rs.Primary.ID)
			}
			time.Sleep(10 * time.Second)
		}
	}
}

// waitDatabaseAccountAbsent polls DescribeDatabases until the account is no
// longer authorized for the database (the revoke API is asynchronous).
func waitDatabaseAccountAbsent(resourceId, accountName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("can't find resource by id: %s", resourceId)
		}
		deadline := time.Now().Add(3 * time.Minute)
		last := "<unread>"
		for {
			object, err := polardb.NewPolarDbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient)).DescribePolarDbDatabase(rs.Primary.ID)
			if err == nil {
				record := polarDBDatabaseAccountRecord(object, accountName)
				if record == nil {
					return nil
				}
				last = fmt.Sprint(record["AccountPrivilege"])
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("timed out waiting for account %q to be revoked, last privilege %q for %s", accountName, last, rs.Primary.ID)
			}
			time.Sleep(10 * time.Second)
		}
	}
}

// TestAccAliCloudPolarDbDatabase_accountPrivilegeMySQL covers account_privilege
// on a MySQL cluster: the value is sent by CreateDatabase, verified against
// DescribeDatabases, and read back into state / import. The extra steps cover
// the multi-account responses that must not bind another account's privilege to
// the ForceNew attribute, and the drift signalled when the configured account's
// grant disappears.
func TestAccAliCloudPolarDbDatabase_accountPrivilegeMySQL(t *testing.T) {
	resourceId := polarDBDatabaseTestResourceId
	testAccCheck := polarDBDatabaseCheck(resourceId, AliCloudPolarDbDatabaseMap11854)
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, aliCloudPolarDbDatabaseMySQLDependence)
	privilegeConfig := map[string]interface{}{
		"db_cluster_id":     "${alicloud_polardb_account.default.db_cluster_id}",
		"db_name":           name,
		"account_name":      "${alicloud_polardb_account.default.account_name}",
		"account_privilege": "ReadOnly",
		"db_description":    name,
	}
	// Omitting account_privilege lets a step re-grant the account without asking
	// Terraform to reconcile the cleared state against the configured value.
	privilegeConfigWithoutPrivilege := map[string]interface{}{
		"db_cluster_id":  "${alicloud_polardb_account.default.db_cluster_id}",
		"db_name":        name,
		"account_name":   "${alicloud_polardb_account.default.account_name}",
		"db_description": name,
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { polarDBDatabaseTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             polarDBDatabaseDestroyCheck(resourceId, AliCloudPolarDbDatabaseMap11854),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(privilegeConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id":     CHECKSET,
						"db_name":           name,
						"account_name":      CHECKSET,
						"account_privilege": "ReadOnly",
						"db_description":    name,
						"engine":            "MySQL",
					}),
					checkDatabaseAccountPrivilege(resourceId, name, "ReadOnly"),
				),
			},
			{
				// Baseline refresh: the single authorized account is read back.
				RefreshState: true,
				Check: testAccCheck(map[string]string{
					"account_name":      name,
					"account_privilege": "ReadOnly",
				}),
			},
			{
				// Revoke the only grant. DescribeDatabases then returns a legal
				// empty Accounts.Account list, which the post-apply plan must
				// surface as drift instead of keeping the stale privilege.
				Config:             testAccConfig(privilegeConfig),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					revokeDatabaseAccountPrivilege(resourceId, name),
					waitDatabaseAccountAbsent(resourceId, name),
				),
			},
			{
				// The legal empty list clears the stale privilege in state.
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
				Check: testAccCheck(map[string]string{
					"account_name":      name,
					"account_privilege": "",
				}),
			},
			{
				// Restore the grant without configuring the privilege, so the
				// cleared state does not plan a replacement.
				Config: testAccConfig(privilegeConfigWithoutPrivilege),
				Check: resource.ComposeTestCheckFunc(
					grantDatabaseAccountPrivilege(resourceId, name, "ReadOnly"),
					waitDatabaseAccountPrivilege(resourceId, name, "ReadOnly"),
				),
			},
			{
				// A refresh reads the restored privilege back into state.
				RefreshState: true,
				Check: testAccCheck(map[string]string{
					"account_name":      name,
					"account_privilege": "ReadOnly",
				}),
			},
			{
				// Grant a second account ReadWrite on the same database so the
				// following refresh observes two authorized accounts.
				Config: testAccConfig(privilegeConfig),
				Check: resource.ComposeTestCheckFunc(
					grantDatabaseAccountPrivilege(resourceId, name+"b", "ReadWrite"),
					waitDatabaseAccountPrivilege(resourceId, name+"b", "ReadWrite"),
				),
			},
			{
				// With two authorized accounts the configured account must keep its
				// own ReadOnly privilege regardless of the response order.
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_name":      name,
						"account_privilege": "ReadOnly",
					}),
				),
			},
			{
				// Importing a database with several accounts cannot tell which account
				// the configuration manages, so neither account field is inferred from
				// the list position.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected 1 imported state, got %d", len(states))
					}
					attributes := states[0].Attributes
					if got := attributes["account_name"]; got != "" {
						return fmt.Errorf("expected imported account_name to stay empty for a multi-account database, got %q", got)
					}
					if got := attributes["account_privilege"]; got != "" {
						return fmt.Errorf("expected imported account_privilege to stay empty for a multi-account database, got %q", got)
					}
					return nil
				},
			},
			{
				// Remove the configured account's grant while the second account keeps
				// its own; the next refresh must surface the drift. The revocation is
				// performed in the check, so this step intentionally ends with a
				// non-empty plan.
				Config:             testAccConfig(privilegeConfig),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					revokeDatabaseAccountPrivilege(resourceId, name),
					waitDatabaseAccountAbsent(resourceId, name),
				),
			},
			{
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_name":      name,
						"account_privilege": "",
					}),
				),
			},
		},
	})
}

// TestAccAliCloudPolarDbDatabase_accountPrivilegePostgreSQL covers
// account_privilege on a PostgreSQL cluster. PostgreSQL does not report an
// AccountPrivilege value for the database owner, so the API fact asserted here
// is the grant status while the state keeps the configured value. The import
// proves the value is unrecoverable, and a plan-only step proves that keeping it
// in the configuration plans a replacement.
func TestAccAliCloudPolarDbDatabase_accountPrivilegePostgreSQL(t *testing.T) {
	resourceId := polarDBDatabaseTestResourceId
	testAccCheck := polarDBDatabaseCheck(resourceId, AliCloudPolarDbDatabaseMap11854)
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
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
		CheckDestroy:             polarDBDatabaseDestroyCheck(resourceId, AliCloudPolarDbDatabaseMap11854),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(privilegeConfig),
				Check: resource.ComposeTestCheckFunc(
					checkDatabaseAccountEmpowered(resourceId, name),
					testAccCheck(map[string]string{
						"db_cluster_id":     CHECKSET,
						"db_name":           name,
						"account_name":      CHECKSET,
						"account_privilege": "ReadOnly",
					}),
				),
			},
			{
				// The API does not return AccountPrivilege for the PostgreSQL
				// owner, so an import cannot recover the configured value. The
				// imported state is persisted so the next step plans against it.
				ResourceName:       resourceId,
				ImportState:        true,
				ImportStateVerify:  false,
				ImportStatePersist: true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected 1 imported state, got %d", len(states))
					}
					attributes := states[0].Attributes
					if got := attributes["account_name"]; got != "" {
						return fmt.Errorf("expected imported account_name to stay empty because the import has no account identity, got %q", got)
					}
					if got := attributes["account_privilege"]; got != "" {
						return fmt.Errorf("expected account_privilege to be unrecoverable after import, got %q", got)
					}
					return nil
				},
			},
			{
				// Keeping account_privilege in the configuration after the import
				// must replace this database (ForceNew). This step only plans
				// against the persisted imported state: PlanOnly skips the
				// apply path, so the exact replacement action is asserted on the
				// plan-only follow-up plan and the expected non-empty plan is
				// allowed instead of being applied.
				Config:             testAccConfig(privilegeConfig),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionReplace),
					},
				},
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
	testAccCheck := polarDBDatabaseCheck(resourceId, AliCloudPolarDbDatabaseMap11854)
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, aliCloudPolarDbDatabaseOracleDependence)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { polarDBDatabaseTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             polarDBDatabaseDestroyCheck(resourceId, AliCloudPolarDbDatabaseMap11854),
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
						"db_cluster_id":     CHECKSET,
						"db_name":           name,
						"account_name":      CHECKSET,
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
						"account_name": name + "u",
					}),
					checkDatabaseAccountName(resourceId, name+"u"),
					checkDatabaseAccountEmpowered(resourceId, name+"u"),
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
  db_name           = "tf_testacc_privilege_without_account"
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

// TestAccAliCloudPolarDbDatabase_accountPrivilegeUnknownAccountName covers the
// apply-time guard: an account_name expression that is unknown at plan time can
// resolve to an empty string when Create runs. CreateDatabase only applies
// AccountPrivilege together with AccountName, so the request must be rejected
// before it is sent. The VPC id keeps the value unknown during plan and empty
// during apply without creating a PolarDB cluster.
func TestAccAliCloudPolarDbDatabase_accountPrivilegeUnknownAccountName(t *testing.T) {
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	config := fmt.Sprintf(`
resource "alicloud_vpc" "default" {
  vpc_name   = "tf-testAcc-polardb-db-guard-%[1]d"
  cidr_block = "172.19.0.0/16"
}

resource "alicloud_polardb_database" "default" {
  db_cluster_id     = "pc-unknown-account-guard"
  db_name           = "%[2]s"
  account_name      = substr(alicloud_vpc.default.id, 0, 0)
  account_privilege = "ReadOnly"
}
`, rand, name)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { polarDBDatabaseTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				// The account_name expression is unknown at plan time, so the plan
				// must succeed without the CustomizeDiff rejection.
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				Check:              func(_ *terraform.State) error { return nil },
			},
			{
				// The value resolves to an empty string at apply time. The rejection
				// happens before CreateDatabase is sent, either in the apply-time
				// CustomizeDiff or in the create guard.
				Config:      config,
				ExpectError: regexp.MustCompile(`"account_privilege" requires "account_name"`),
				// The failed step creates no database state; the no-op check keeps
				// the step shape for the coverage parser.
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
  db_node_class = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
  pay_type      = "PostPaid"
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

resource "alicloud_polardb_account" "second" {
  db_cluster_id    = alicloud_polardb_cluster.default.id
  account_name     = "${var.name}b"
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
  account_name     = "${var.name}u"
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
  db_type            = "Oracle"
  db_version         = "14"
  db_node_class      = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
  db_node_count      = "2"
  pay_type           = "PostPaid"
  vswitch_id         = alicloud_vswitch.default.id
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
  account_name     = "${var.name}u"
  account_password = "Example1234"
}
`, name)
}
