package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudClickHouseAccount_basic0(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_click_house_account.default"
	ra := resourceAttrInit(resourceId, AliCloudClickHouseAccountMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ClickhouseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeClickHouseAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	pwd := fmt.Sprintf("Tf-test%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudClickHouseAccountBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"db_cluster_id":    CHECKSET,
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
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_click_house_account.default"
	ra := resourceAttrInit(resourceId, AliCloudClickHouseAccountMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ClickhouseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeClickHouseAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	pwd := fmt.Sprintf("Tf-test%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudClickHouseAccountBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"db_cluster_id":      CHECKSET,
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
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_click_house_account.default"
	ra := resourceAttrInit(resourceId, AliCloudClickHouseAccountMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ClickhouseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeClickHouseAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	pwd := fmt.Sprintf("Tf-test%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudClickHouseAccountBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"db_cluster_id":    CHECKSET,
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
	"status":             CHECKSET,
	"type":               CHECKSET,
	"dml_authority":      CHECKSET,
	"ddl_authority":      CHECKSET,
	"allow_databases":    CHECKSET,
	"total_databases":    CHECKSET,
	"allow_dictionaries": CHECKSET,
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
  		db_cluster_version      = "22.8.5.29"
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
