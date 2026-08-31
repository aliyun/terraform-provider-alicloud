package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudAdbResourceGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_adb_resource_group.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbResourceGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbResourceGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestaccAdbResourceGroup%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbResourceGroupBasicDependence0)
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
					"db_cluster_id": "${alicloud_adb_db_cluster.default.id}",
					"group_name":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id": CHECKSET,
						"group_name":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_type": "batch",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_type": "batch",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_num": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_num": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"users": []string{"${alicloud_adb_account.default.account_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"users.#": "1",
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

func TestAccAliCloudAdbResourceGroup_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_adb_resource_group.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbResourceGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbResourceGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-AdbResourceGroup%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbResourceGroupBasicDependence0)
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
					"db_cluster_id": "${alicloud_adb_db_cluster.default.id}",
					"group_name":    name,
					"group_type":    "batch",
					"node_num":      "1",
					"users":         []string{"${alicloud_adb_account.default.account_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id": CHECKSET,
						"group_name":    CHECKSET,
						"group_type":    "batch",
						"node_num":      "1",
						"users.#":       "1",
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

func TestAccAliCloudAdbResourceGroup_basic1(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_adb_resource_group.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbResourceGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbResourceGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("TFADBRG%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbResourceGroupBasicDependence1)
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
					"db_cluster_id":        "${alicloud_adb_lake_account.default.db_cluster_id}",
					"group_name":           name,
					"min_compute_resource": "16ACU",
					"max_compute_resource": "16ACU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id":        CHECKSET,
						"group_name":           CHECKSET,
						"min_compute_resource": "16ACU",
						"max_compute_resource": "16ACU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_compute_resource": "128ACU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_compute_resource": "128ACU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_compute_resource": "128ACU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_compute_resource": "128ACU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_mode":          "AutoScale",
					"cluster_size_resource": "16ACU",
					"max_cluster_count":     "2",
					"min_cluster_count":     "1",
					"min_compute_resource":  REMOVEKEY,
					"max_compute_resource":  REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_mode":          "AutoScale",
						"cluster_size_resource": "16ACU",
						"max_cluster_count":     "2",
						"min_cluster_count":     "1",
						"min_compute_resource":  CHECKSET,
						"max_compute_resource":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_size_resource": "32ACU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_size_resource": "32ACU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_cluster_count": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_cluster_count": "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_cluster_count": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_cluster_count": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"users": []string{"${alicloud_adb_lake_account.default.account_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"users.#": "1",
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

func TestAccAliCloudAdbResourceGroup_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_adb_resource_group.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbResourceGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbResourceGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("TFADBRG%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbResourceGroupBasicDependence1)
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
					"db_cluster_id":        "${alicloud_adb_lake_account.default.db_cluster_id}",
					"group_name":           name,
					"group_type":           "job",
					"cluster_mode":         "Disable",
					"min_compute_resource": "16ACU",
					"max_compute_resource": "16ACU",
					"users":                []string{"${alicloud_adb_lake_account.default.account_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id":        CHECKSET,
						"group_name":           CHECKSET,
						"group_type":           "job",
						"cluster_mode":         "Disable",
						"min_compute_resource": "16ACU",
						"max_compute_resource": "16ACU",
						"users.#":              "1",
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

func TestAccAliCloudAdbResourceGroup_basic2(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_adb_resource_group.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbResourceGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbResourceGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("TFADBRG%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbResourceGroupBasicDependence1)
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
					"db_cluster_id":         "${alicloud_adb_lake_account.default.db_cluster_id}",
					"group_name":            name,
					"cluster_mode":          "AutoScale",
					"cluster_size_resource": "16ACU",
					"max_cluster_count":     "2",
					"min_cluster_count":     "1",
					"engine":                "SparkWarehouse",
					"engine_params": map[string]interface{}{
						"\"spark.adb.version\"":      "3.5",
						"\"spark.app.log.rootPath\"": "oss://" + "${data.alicloud_oss_buckets.default.buckets.0.name}" + "/",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id":         CHECKSET,
						"group_name":            CHECKSET,
						"cluster_mode":          "AutoScale",
						"cluster_size_resource": "16ACU",
						"max_cluster_count":     "2",
						"min_cluster_count":     "1",
						"engine":                "SparkWarehouse",
						"engine_params.%":       "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_size_resource": "36ACU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_size_resource": "36ACU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_cluster_count": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_cluster_count": "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_cluster_count": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_cluster_count": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_params": map[string]interface{}{
						"\"spark.adb.version\"":                 "3.5",
						"\"spark.app.log.rootPath\"":            "oss://" + "${data.alicloud_oss_buckets.default.buckets.0.name}" + "/",
						"\"spark.driver.memoryOverheadFactor\"": "0.5",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_params.%": "3",
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

func TestAccAliCloudAdbResourceGroup_basic2_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_adb_resource_group.default"
	ra := resourceAttrInit(resourceId, AliCloudAdbResourceGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbResourceGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("TFADBRG%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAdbResourceGroupBasicDependence1)
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
					"db_cluster_id":         "${alicloud_adb_lake_account.default.db_cluster_id}",
					"group_name":            name,
					"group_type":            "interactive",
					"cluster_mode":          "AutoScale",
					"cluster_size_resource": "16ACU",
					"max_cluster_count":     "2",
					"min_cluster_count":     "1",
					"engine":                "SparkWarehouse",
					"engine_params": map[string]interface{}{
						"\"spark.adb.version\"":      "3.5",
						"\"spark.app.log.rootPath\"": "oss://" + "${data.alicloud_oss_buckets.default.buckets.0.name}" + "/",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id":         CHECKSET,
						"group_name":            name,
						"group_type":            "interactive",
						"cluster_mode":          "AutoScale",
						"cluster_size_resource": "16ACU",
						"max_cluster_count":     "2",
						"min_cluster_count":     "1",
						"engine":                "SparkWarehouse",
						"engine_params.%":       "2",
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

var AliCloudAdbResourceGroupMap0 = map[string]string{
	"group_type": CHECKSET,
	"status":     CHECKSET,
}

func AliCloudAdbResourceGroupBasicDependence0(name string) string {
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
	
	resource "alicloud_adb_db_cluster" "default" {
  		compute_resource    = "32Core128GBNEW"
  		db_cluster_category = "MixedStorage"
  		description         = var.name
  		elastic_io_resource = 1
  		mode                = "flexible"
  		payment_type        = "PayAsYouGo"
  		vpc_id              = data.alicloud_vpcs.default.ids.0
  		vswitch_id          = data.alicloud_vswitches.default.ids.0
  		zone_id             = data.alicloud_adb_zones.default.zones.0.id
	}
	
	resource "alicloud_adb_account" "default" {
  		db_cluster_id    = alicloud_adb_db_cluster.default.id
  		account_name     = "tf_account_name"
  		account_password = "YourPassword123!"
	}
`, name)
}

func AliCloudAdbResourceGroupBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_oss_buckets" "default" {
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
	
	resource "alicloud_adb_lake_account" "default" {
  		db_cluster_id    = alicloud_adb_db_cluster_lake_version.default.id
		account_type     = "Super"
  		account_name     = "tf_account_name"
  		account_password = "YourPassword123!"
	}
`, name)
}
