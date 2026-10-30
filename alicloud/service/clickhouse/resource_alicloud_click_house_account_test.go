package clickhouse_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/clickhouse"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudClickHouseAccount_basic0(t *testing.T) {
	// 25.3/Basic/S8 is orderable in cn-shanghai (proven by backup_policy R ACC);
	// VbrSupportRegions (cn-hangzhou) test-account VPC quota is exhausted, so use
	// ClickHouseBackupPolicySupportRegions instead. No production logic changed.
	acctest.CheckoutSupportedRegions(t, true, connectivity.ClickHouseBackupPolicySupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_click_house_account.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudClickHouseAccountMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return clickhouse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeClickHouseAccount")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	pwd := fmt.Sprintf("Tf-test%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudClickHouseAccountBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":    "${alicloud_click_house_db_cluster.default.id}",
					"account_name":     name,
					"account_password": pwd,
					"type":             "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id":    acctest.CHECKSET,
						"account_name":     name,
						"account_password": pwd,
						"type":             "Normal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": name + "_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description": name + "_desc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": pwd + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": pwd + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": "AccountDescription_all",
					"account_password":    pwd + "updateall",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description": "AccountDescription_all",
						"account_password":    pwd + "updateall",
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

func TestAccAliCloudClickHouseAccount_basic1(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_click_house_account.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudClickHouseAccountMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return clickhouse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeClickHouseAccount")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	pwd := fmt.Sprintf("Tf-test%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudClickHouseAccountBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":      "${alicloud_click_house_db_cluster.default.id}",
					"account_name":       name,
					"account_password":   pwd,
					"type":               "Normal",
					"ddl_authority":      "true",
					"dml_authority":      "all",
					"allow_databases":    "db1",
					"allow_dictionaries": "dt1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id":      acctest.CHECKSET,
						"account_name":       name,
						"account_password":   pwd,
						"type":               "Normal",
						"allow_databases":    "db1",
						"dml_authority":      "all",
						"allow_dictionaries": "dt1",
						"ddl_authority":      "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allow_databases": "db1,db2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allow_databases": "db1,db2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dml_authority": "readOnly,modify",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dml_authority": "readOnly,modify",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ddl_authority": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ddl_authority": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allow_dictionaries": "dt1,dt2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allow_dictionaries": "dt1,dt2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allow_databases":    "db1",
					"dml_authority":      "all",
					"allow_dictionaries": "dt1",
					"ddl_authority":      "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allow_databases":    "db1",
						"dml_authority":      "all",
						"allow_dictionaries": "dt1",
						"ddl_authority":      "false",
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

func TestAccAliCloudClickHouseAccount_CreateSuperAccount(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_click_house_account.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudClickHouseAccountMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return clickhouse.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeClickHouseAccount")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	pwd := fmt.Sprintf("Tf-test%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudClickHouseAccountBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":    "${alicloud_click_house_db_cluster.default.id}",
					"account_name":     name,
					"account_password": pwd,
					"type":             "Super",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id":    acctest.CHECKSET,
						"account_name":     name,
						"account_password": pwd,
						"type":             "Super",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type": "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type": "Normal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type": "Super",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type": "Super",
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

var AliCloudClickHouseAccountMap0 = map[string]string{
	"status":             acctest.CHECKSET,
	"type":               acctest.CHECKSET,
	"dml_authority":      acctest.CHECKSET,
	"ddl_authority":      acctest.CHECKSET,
	"allow_databases":    acctest.CHECKSET,
	"total_databases":    acctest.CHECKSET,
	"allow_dictionaries": acctest.CHECKSET,
}

func AliCloudClickHouseAccountBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_click_house_regions" "default" {
  		current = true
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = data.alicloud_click_house_regions.default.regions.0.zone_ids.0.zone_id
	}

	resource "alicloud_click_house_db_cluster" "default" {
  		db_cluster_version      = "25.3"
  		category                = "Basic"
  		db_cluster_class        = "S8"
  		db_cluster_network_type = "vpc"
  		db_cluster_description  = var.name
  		db_node_group_count     = 1
  		payment_type            = "PayAsYouGo"
  		db_node_storage         = "100"
  		storage_type            = "cloud_essd"
  		vswitch_id              = alicloud_vswitch.default.id
	}
`, name)
}
