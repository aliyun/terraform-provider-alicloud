package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
