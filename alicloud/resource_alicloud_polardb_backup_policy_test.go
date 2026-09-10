package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudPolarDBBackupPolicy(t *testing.T) {
	var v *polardb.DescribeBackupPolicyResponse
	resourceId := "alicloud_polardb_backup_policy.default"
	serverFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeBackupPolicy")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccPolarDBbackuppolicy"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBBackupPolicyConfigDependence)
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
					"db_cluster_id":           "${alicloud_polardb_cluster.default.id}",
					"preferred_backup_period": []string{"Tuesday", "Wednesday"},
					"preferred_backup_time":   "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "2",
						"preferred_backup_time":     "10:00Z-11:00Z",
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
					"preferred_backup_period": []string{"Wednesday", "Monday", "Saturday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_time": "15:00Z-16:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_time": "15:00Z-16:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_period": []string{"Tuesday", "Thursday", "Friday", "Sunday"},
					"preferred_backup_time":   "17:00Z-18:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "4",
						"preferred_backup_time":     "17:00Z-18:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_retention_policy_on_cluster_deletion": "LATEST",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_retention_policy_on_cluster_deletion": "LATEST",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudPolarDBNewBackupPolicy(t *testing.T) {
	var v *polardb.DescribeBackupPolicyResponse
	resourceId := "alicloud_polardb_backup_policy.default"
	serverFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	var basicMap = map[string]string{
		"preferred_backup_time":                              CHECKSET,
		"backup_retention_policy_on_cluster_deletion":        CHECKSET,
		"backup_frequency":                                   CHECKSET,
		"data_level2_backup_another_region_retention_period": CHECKSET,
		"data_level1_backup_frequency":                       CHECKSET,
		"backup_retention_period":                            CHECKSET,
		"data_level1_backup_time":                            CHECKSET,
		"data_level1_backup_retention_period":                CHECKSET,
		"data_level2_backup_retention_period":                CHECKSET,
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeBackupPolicy")
	ra := resourceAttrInit(resourceId, basicMap)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccPolarDBNewBackupPolicy"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBBackupPolicyConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.PolarDBNewBackupPolicySupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":           "${alicloud_polardb_cluster.default.id}",
					"preferred_backup_period": []string{"Tuesday", "Wednesday"},
					"preferred_backup_time":   "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "2",
						"preferred_backup_time":     "10:00Z-11:00Z",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"log_backup_another_region_region", "log_backup_another_region_retention_period"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_level1_backup_retention_period": "7",
					"backup_retention_period":             "7",
					"data_level1_backup_time":             "10:00Z-11:00Z",
					"data_level1_backup_period":           []string{"Tuesday", "Wednesday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_level1_backup_retention_period": "7",
						"backup_retention_period":             "7",
						"data_level1_backup_time":             "10:00Z-11:00Z",
						"data_level1_backup_period.#":         "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_level1_backup_retention_period":                "7",
					"data_level1_backup_time":                            "10:00Z-11:00Z",
					"data_level1_backup_period":                          []string{"Tuesday", "Wednesday"},
					"data_level1_backup_frequency":                       "Normal",
					"data_level2_backup_retention_period":                "30",
					"data_level2_backup_period":                          []string{"Tuesday", "Wednesday"},
					"data_level2_backup_another_region_region":           "cn-beijing",
					"data_level2_backup_another_region_retention_period": "30",
					"backup_retention_policy_on_cluster_deletion":        "NONE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_level1_backup_retention_period":                "7",
						"data_level1_backup_time":                            "10:00Z-11:00Z",
						"data_level1_backup_period.#":                        "2",
						"data_level1_backup_frequency":                       "Normal",
						"data_level2_backup_retention_period":                "30",
						"data_level2_backup_period.#":                        "2",
						"data_level2_backup_another_region_region":           "cn-beijing",
						"data_level2_backup_another_region_retention_period": "30",
						"backup_retention_policy_on_cluster_deletion":        "NONE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_retention_policy_on_cluster_deletion": "LATEST",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_retention_policy_on_cluster_deletion": "LATEST",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_frequency": "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_frequency": "Normal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_level1_backup_frequency": "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_level1_backup_frequency": "Normal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup_retention_period": "4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup_retention_period": "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup_another_region_region": "cn-beijing",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup_another_region_region": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup_another_region_retention_period": "32",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup_another_region_retention_period": CHECKSET,
					}),
				),
			},
		},
	})
}

func resourcePolarDBBackupPolicyConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
	}
	resource "alicloud_vpc" "default" {
		vpc_name = var.name
	}
	resource "alicloud_vswitch" "default" {
		zone_id = data.alicloud_polardb_node_classes.default.classes.0.zone_id
		vpc_id = alicloud_vpc.default.id
		cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
	}
	data "alicloud_polardb_node_classes" "default" {
		pay_type   = "PostPaid"
		db_type    = "MySQL"
		db_version = "8.0"
		category   = "Normal"
	}
	resource "alicloud_polardb_cluster" "default" {
		db_type       = "MySQL"
		db_version    = "8.0"
		pay_type      = "PostPaid"
		db_node_class = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.2.db_node_class
		vswitch_id    = alicloud_vswitch.default.id
		description   = "${var.name}"
	}
	`, name)
}
