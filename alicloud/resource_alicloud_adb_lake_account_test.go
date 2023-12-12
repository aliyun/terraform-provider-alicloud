package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Adb LakeAccount. >>> Resource test cases, automatically generated.
// Case 5287
func TestAccAliCloudAdbLakeAccount_basic5287(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_adb_lake_account.default"
	ra := resourceAttrInit(resourceId, AlicloudAdbLakeAccountMap5287)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbLakeAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sadblakeaccount%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAdbLakeAccountBasicDependence5287)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.HologramSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":    "${alicloud_adb_db_cluster_lake_version.CreateInstance.id}",
					"account_type":     "Super",
					"account_name":     "tfnormal",
					"account_password": "normal@2022",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id":    CHECKSET,
						"account_type":     "Super",
						"account_name":     "tfnormal",
						"account_password": "normal@2022",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": "test_tf_des",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description": "test_tf_des",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_privileges": []map[string]interface{}{
						{
							"privilege_type": "Column",
							"privilege_object": []map[string]interface{}{
								{
									"database": "MYSQL",
									"table":    "COLUMNS_PRIV",
									"column":   "DB",
								},
							},
							"privileges": []string{
								"create", "select", "update"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_privileges.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "normal@2022",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "normal@2022",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": "test_tf_des2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description": "test_tf_des2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": "test_tf_des",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description": "test_tf_des",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "normal@2022",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "normal@2022",
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
				Config: testAccConfig(map[string]interface{}{
					"account_description": "test_tf_des2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description": "test_tf_des2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "normal@2023",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "normal@2023",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": "test_tf_des",
					"db_cluster_id":       "${alicloud_adb_db_cluster_lake_version.CreateInstance.id}",
					"account_type":        "Super",
					"account_name":        "tfnormal",
					"account_password":    "normal@2022",
					"account_privileges": []map[string]interface{}{
						{
							"privilege_type": "Column",
							"privilege_object": []map[string]interface{}{
								{
									"database": "MYSQL",
									"table":    "COLUMNS_PRIV",
									"column":   "DB",
								},
							},
							"privileges": []string{
								"create", "select", "update"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description":  "test_tf_des",
						"db_cluster_id":        CHECKSET,
						"account_type":         "Super",
						"account_name":         "tfnormal",
						"account_password":     "normal@2022",
						"account_privileges.#": "1",
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

var AlicloudAdbLakeAccountMap5287 = map[string]string{
	"status": CHECKSET,
}

func AlicloudAdbLakeAccountBasicDependence5287(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "VPCID" {
  vpc_name = var.name

  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "VSWITCHID" {
  vpc_id       = alicloud_vpc.VPCID.id
  zone_id      = "cn-hangzhou-k"
  vswitch_name = var.name

  cidr_block = "172.16.0.0/24"
}

resource "alicloud_adb_db_cluster_lake_version" "CreateInstance" {
  storage_resource        = "0ACU"
  zone_id                 = "cn-hangzhou-k"
  vpc_id                  = alicloud_vpc.VPCID.id
  vswitch_id              = alicloud_vswitch.VSWITCHID.id
  db_cluster_description  = "tf自动化测试-杭州-资源组"
  compute_resource        = "16ACU"
  db_cluster_version      = "5.0"
  payment_type            = "PayAsYouGo"
  security_ips            = "127.0.0.1"
}


`, name)
}

// Case 5287  twin
func TestAccAliCloudAdbLakeAccount_basic5287_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_adb_lake_account.default"
	ra := resourceAttrInit(resourceId, AlicloudAdbLakeAccountMap5287)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbLakeAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sadblakeaccount%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAdbLakeAccountBasicDependence5287)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.HologramSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": "test_tf_des2",
					"db_cluster_id":       "${alicloud_adb_db_cluster_lake_version.CreateInstance.id}",
					"account_type":        "Super",
					"account_name":        "tfnormal",
					"account_password":    "normal@2023",
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
						"account_description":  "test_tf_des2",
						"db_cluster_id":        CHECKSET,
						"account_type":         "Super",
						"account_name":         "tfnormal",
						"account_password":     "normal@2023",
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

// Test Adb LakeAccount. <<< Resource test cases, automatically generated.
