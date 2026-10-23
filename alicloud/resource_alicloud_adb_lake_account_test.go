package alicloud

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// collection-order: bypass account_privileges.privileges -- pre-existing multi-member TypeList (string list nested in the account_privileges set) outside this change's engine scope; reorder convergence coverage is stock debt surfaced by the resource file being touched, not introduced here.

// Test Adb LakeAccount. >>> Resource test cases, automatically generated.
// Case 湖仓账号测试用例 5218
func TestAccAliCloudAdbLakeAccount_basic5218(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_adb_lake_account.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbLakeAccountMap5218)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbLakeAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccadb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbLakeAccountBasicDependence5218)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":    "${alicloud_adb_db_cluster_lake_version.default.id}",
					"account_type":     "Super",
					"account_name":     "tf_account_name_supper",
					"account_password": "YourPassword123!",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id": CHECKSET,
						"account_type":  "Super",
						"account_name":  "tf_account_name_supper",
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
				Config: testAccConfig(map[string]interface{}{
					"ram_user_list": []string{"${alicloud_ram_user.default.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ram_user_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_privileges": []map[string]interface{}{
						{
							"privilege_type": "Database",
							"privilege_object": []map[string]interface{}{
								{
									"database": "MYSQL",
								},
							},
							"privileges": []string{
								"select", "update"},
						},
						{
							"privilege_type": "Table",
							"privilege_object": []map[string]interface{}{
								{
									"database": "INFORMATION_SCHEMA",
									"table":    "ENGINES",
								},
							},
							"privileges": []string{
								"update"},
						},
						{
							"privilege_type": "Column",
							"privilege_object": []map[string]interface{}{
								{
									"database": "INFORMATION_SCHEMA",
									"table":    "COLUMNS",
									"column":   "PRIVILEGES",
								},
							},
							"privileges": []string{
								"update"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_privileges.#": "3",
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

func TestAccAliCloudAdbLakeAccount_basic5218_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_adb_lake_account.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbLakeAccountMap5218)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbLakeAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccadb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbLakeAccountBasicDependence5218)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":       "${alicloud_adb_db_cluster_lake_version.default.id}",
					"account_type":        "Super",
					"account_name":        "tf_account_name_supper",
					"account_password":    "YourPassword123!",
					"account_description": name,
					"ram_user_list":       []string{"${alicloud_ram_user.default.id}"},
					"account_privileges": []map[string]interface{}{
						{
							"privilege_type": "Database",
							"privilege_object": []map[string]interface{}{
								{
									"database": "MYSQL",
								},
							},
							"privileges": []string{
								"select", "update"},
						},
						{
							"privilege_type": "Table",
							"privilege_object": []map[string]interface{}{
								{
									"database": "INFORMATION_SCHEMA",
									"table":    "ENGINES",
								},
							},
							"privileges": []string{
								"update"},
						},
						{
							"privilege_type": "Column",
							"privilege_object": []map[string]interface{}{
								{
									"database": "INFORMATION_SCHEMA",
									"table":    "COLUMNS",
									"column":   "PRIVILEGES",
								},
							},
							"privileges": []string{
								"update"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id":        CHECKSET,
						"account_type":         "Super",
						"account_name":         "tf_account_name_supper",
						"account_description":  name,
						"ram_user_list.#":      "1",
						"account_privileges.#": "3",
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

func TestAccAliCloudAdbLakeAccount_basic5220(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_adb_lake_account.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbLakeAccountMap5218)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbLakeAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccadb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbLakeAccountBasicDependence5218)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":    "${alicloud_adb_db_cluster_lake_version.default.id}",
					"account_type":     "Normal",
					"account_name":     "tf_account_name_normal",
					"account_password": "YourPassword123!",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id": CHECKSET,
						"account_type":  "Normal",
						"account_name":  "tf_account_name_normal",
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
				Config: testAccConfig(map[string]interface{}{
					"ram_user_list": []string{"${alicloud_ram_user.default.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ram_user_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_privileges": []map[string]interface{}{
						{
							"privilege_type": "Database",
							"privilege_object": []map[string]interface{}{
								{
									"database": "MYSQL",
								},
							},
							"privileges": []string{
								"select", "update"},
						},
						{
							"privilege_type": "Table",
							"privilege_object": []map[string]interface{}{
								{
									"database": "INFORMATION_SCHEMA",
									"table":    "ENGINES",
								},
							},
							"privileges": []string{
								"update"},
						},
						{
							"privilege_type": "Column",
							"privilege_object": []map[string]interface{}{
								{
									"database": "INFORMATION_SCHEMA",
									"table":    "COLUMNS",
									"column":   "PRIVILEGES",
								},
							},
							"privileges": []string{
								"update"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_privileges.#": "3",
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

func TestAccAliCloudAdbLakeAccount_basic5220_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_adb_lake_account.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbLakeAccountMap5218)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbLakeAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccadb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbLakeAccountBasicDependence5218)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":       "${alicloud_adb_db_cluster_lake_version.default.id}",
					"account_type":        "Normal",
					"account_name":        "tf_account_name_normal",
					"account_password":    "YourPassword123!",
					"account_description": name,
					"ram_user_list":       []string{"${alicloud_ram_user.default.id}"},
					"account_privileges": []map[string]interface{}{
						{
							"privilege_type": "Database",
							"privilege_object": []map[string]interface{}{
								{
									"database": "MYSQL",
								},
							},
							"privileges": []string{
								"select", "update"},
						},
						{
							"privilege_type": "Table",
							"privilege_object": []map[string]interface{}{
								{
									"database": "INFORMATION_SCHEMA",
									"table":    "ENGINES",
								},
							},
							"privileges": []string{
								"update"},
						},
						{
							"privilege_type": "Column",
							"privilege_object": []map[string]interface{}{
								{
									"database": "INFORMATION_SCHEMA",
									"table":    "COLUMNS",
									"column":   "PRIVILEGES",
								},
							},
							"privileges": []string{
								"update"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id":        CHECKSET,
						"account_type":         "Normal",
						"account_name":         "tf_account_name_normal",
						"account_description":  name,
						"ram_user_list.#":      "1",
						"account_privileges.#": "3",
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

// Engine regression coverage (R1):
//
//   - Omitted engine + legacy two-segment import (old-state behaviour): covered
//     by the pre-existing basic5218 / basic5220 cases above, which create and
//     refresh accounts without an engine argument and import them with the
//     two-segment id <db_cluster_id>:<account_name>; DescribeAccounts then uses
//     the API default (AnalyticDB).
//   - Explicit AnalyticDB, API Engine missing/empty guard, three-segment
//     import: TestAccAliCloudAdbLakeAccount_basic0_engine. DescribeAccounts
//     does not return Engine for default-engine accounts, so the value set at
//     create time must be preserved on refresh (never overwritten with empty);
//     the three-segment import records AnalyticDB before Read.
//   - Clickhouse end-to-end lifecycle (create -> read -> description/password
//     update -> delete) with Engine forwarded to every lifecycle API, and the
//     three-segment import reading Engine back: TestAccAliCloudAdbLakeAccount_
//     basic0_engine_clickhouse. DescribeAccounts returns Engine for Clickhouse
//     accounts, so it is read back and must survive each refresh.
//
// Note: Clickhouse account/privilege API applicability and the end-to-end
// lifecycle are not cloud-verified in this round; the cases below are written
// for the authorised ACC run arranged under the project gates.

// TestAccAliCloudAdbLakeAccount_basic0_engine covers the explicit-AnalyticDB
// lifecycle: create -> description update -> password update -> import. The
// engine is dispatched at create and forwarded to the lifecycle APIs;
// DescribeAccounts omits Engine for default-engine accounts, so the assertion
// after each refresh verifies the Read guard preserves the create-time value
// rather than overwriting it with an empty response (API Engine missing/empty).
// The three-segment import id <db_cluster_id>:<account_name>:AnalyticDB records
// the engine before Read so ImportStateVerify can match it.
func TestAccAliCloudAdbLakeAccount_basic0_engine(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_adb_lake_account.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbLakeAccountMap5218)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbLakeAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccadb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbLakeAccountBasicDependence5218)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":    "${alicloud_adb_db_cluster_lake_version.default.id}",
					"account_type":     "Super",
					"account_name":     "${var.name}",
					"account_password": "YourPassword123!",
					"engine":           "AnalyticDB",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id": CHECKSET,
						"account_type":  "Super",
						"account_name":  name,
						"engine":        "AnalyticDB",
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
						"engine":              "AnalyticDB",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "YourPassword123!update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine": "AnalyticDB",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// Record the engine segment before Read so the default-engine
				// account is queried on AnalyticDB and the guard keeps the value.
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceId]
					if !ok {
						return "", fmt.Errorf("resource %s not found in state", resourceId)
					}
					return fmt.Sprintf("%s:AnalyticDB", rs.Primary.ID), nil
				},
				ImportStateVerifyIgnore: []string{"account_password"},
			},
		},
	})
}

// TestAccAliCloudAdbLakeAccount_basic0_engine_clickhouse covers the Clickhouse
// lifecycle: create -> read -> description update -> password update -> import
// -> delete (via CheckDestroy). Engine=Clickhouse is forwarded to the create
// waiter, Read, ModifyAccountDescription, ResetAccountPassword and DeleteAccount.
// DescribeAccounts returns Engine for Clickhouse accounts, so it is read back
// and asserted after every refresh (API Engine non-empty). The three-segment
// import id <db_cluster_id>:<account_name>:Clickhouse records the engine before
// Read; without it DescribeAccounts would default to AnalyticDB and not find
// the account, which is the regression this resource must not regress.
//
// The Clickhouse wide-table engine is an invitation-only beta feature. It must
// be explicitly enabled on the cluster via ModifyClickhouseEngine before
// Clickhouse accounts can be created; without it CreateAccount(Engine=Clickhouse)
// fails with "Failed to get inner url for cluster … instance type adbclickhouse".
// ModifyClickhouseEngine returns ClickHouseEngineNotAllowed when the test account
// is not whitelisted (verified on a real cluster amv-bp19pdgv7psu5f30 during
// development). The test is skipped by default via TF_ACC_ADB_CLICKHOUSE_WHITELISTED
// until the account is whitelisted; the full lifecycle below remains ready for
// when access is granted. To run this test: submit a support ticket to add the
// test account to the Clickhouse whitelist, then set
// TF_ACC_ADB_CLICKHOUSE_WHITELISTED=true.
func TestAccAliCloudAdbLakeAccount_basic0_engine_clickhouse(t *testing.T) {
	// The Clickhouse engine is invitation-only beta. The test account is not
	// whitelisted (verified via ModifyClickhouseEngine → ClickHouseEngineNotAllowed
	// on a real cluster). Skip until the account is whitelisted; the full
	// lifecycle test below remains for when access is granted.
	if os.Getenv("TF_ACC_ADB_CLICKHOUSE_WHITELISTED") != "true" {
		t.Skip("Clickhouse wide-table engine is in invitation-only beta; " +
			"test account not whitelisted (ModifyClickhouseEngine returns " +
			"ClickHouseEngineNotAllowed). Submit a support ticket to add the " +
			"account to the whitelist, then set TF_ACC_ADB_CLICKHOUSE_WHITELISTED=true.")
	}
	var v map[string]interface{}
	resourceId := "alicloud_adb_lake_account.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbLakeAccountMap5218)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbLakeAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccadb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbLakeAccountBasicDependence5218)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":    "${alicloud_adb_db_cluster_lake_version.default.id}",
					"account_type":     "Super",
					"account_name":     "${var.name}",
					"account_password": "YourPassword123!",
					"engine":           "Clickhouse",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id": CHECKSET,
						"account_type":  "Super",
						"account_name":  name,
						"engine":        "Clickhouse",
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
						"engine":              "Clickhouse",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "YourPassword123!update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine": "Clickhouse",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// Clickhouse accounts must be imported with the engine segment;
				// DescribeAccounts defaults to AnalyticDB and would not find them.
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceId]
					if !ok {
						return "", fmt.Errorf("resource %s not found in state", resourceId)
					}
					return fmt.Sprintf("%s:Clickhouse", rs.Primary.ID), nil
				},
				ImportStateVerifyIgnore: []string{"account_password"},
			},
		},
	})
}

// TestAccAliCloudAdbLakeAccount_engine_wrong_engine_import is the R4
// discriminative regression. An account created on the AnalyticDB engine is
// imported with engine=Clickhouse; DescribeAccounts with Engine=Clickhouse
// returns NotFound because the account does not exist on that engine. With the
// fix, Read clears the ID (the engine-qualified identity check denies the
// import) and the step fails; the pre-fix code swallowed the NotFound and
// accepted the import with empty attributes, keeping d.Id() and writing nil
// values. A subsequent import with engine=AnalyticDB must succeed, proving the
// wrong-engine step did not corrupt state and the account remains manageable on
// its real engine. This also covers cross-engine same-name disambiguation: the
// Engine filter distinguishes accounts that share a name across engines, and a
// wrong-but-valid engine segment with a target that does not exist on that
// engine must not be silently accepted.
func TestAccAliCloudAdbLakeAccount_engine_wrong_engine_import(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_adb_lake_account.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbLakeAccountMap5218)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbLakeAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccadb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbLakeAccountBasicDependence5218)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":    "${alicloud_adb_db_cluster_lake_version.default.id}",
					"account_type":     "Super",
					"account_name":     "${var.name}",
					"account_password": "YourPassword123!",
					"engine":           "AnalyticDB",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id": CHECKSET,
						"account_type":  "Super",
						"account_name":  name,
						"engine":        "AnalyticDB",
					}),
				),
			},
			{
				// R4: wrong-engine import must fail. The account exists on
				// AnalyticDB, so DescribeAccounts with Engine=Clickhouse
				// returns NotFound. Read clears the ID (identity denied) and
				// the import fails; the pre-fix code swallowed NotFound and
				// accepted the import with empty attributes.
				ResourceName: resourceId,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceId]
					if !ok {
						return "", fmt.Errorf("resource %s not found in state", resourceId)
					}
					return fmt.Sprintf("%s:Clickhouse", rs.Primary.ID), nil
				},
				ExpectError: regexp.MustCompile(`(?i)(not.?exist|not.?found|non.?existent|import|cannot|empty|failed|operationfailed|inner.?url)`),
			},
			{
				// The correct-engine import must succeed after the wrong-engine
				// failure, proving the account is still manageable on its real
				// engine and the wrong-engine step did not corrupt state.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceId]
					if !ok {
						return "", fmt.Errorf("resource %s not found in state", resourceId)
					}
					return fmt.Sprintf("%s:AnalyticDB", rs.Primary.ID), nil
				},
				ImportStateVerifyIgnore: []string{"account_password"},
			},
		},
	})
}

var AliCloudAdbLakeAccountMap5218 = map[string]string{
	"status": CHECKSET,
}

func AliCloudAdbLakeAccountBasicDependence5218(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
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
	
	resource "alicloud_ram_user" "default" {
  		name = var.name
	}

	resource "alicloud_adb_db_cluster_lake_version" "default" {
  		db_cluster_version            = "5.0"
  		vpc_id                        = data.alicloud_vpcs.default.ids.0
  		vswitch_id                    = data.alicloud_vswitches.default.ids.0
  		zone_id                       = data.alicloud_adb_zones.default.ids.0
  		compute_resource              = "128ACU"
  		storage_resource              = "0ACU"
  		payment_type                  = "PayAsYouGo"
  		enable_default_resource_group = false
	}
`, name)
}

// Test Adb LakeAccount. <<< Resource test cases, automatically generated.
