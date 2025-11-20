package alicloud

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var clusterConnectionStringRegexp = "^[a-z-A-Z-0-9]+.rwlb.([a-z-A-Z-0-9]+.){0,1}rds.aliyuncs.com"

func init() {
	resource.AddTestSweepers("alicloud_polardb_cluster", &resource.Sweeper{
		Name: "alicloud_polardb_cluster",
		F:    testSweepPolarDBClusters,
	})
}

func testSweepPolarDBClusters(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var insts []polardb.DBCluster
	req := polardb.CreateDescribeDBClustersRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
			return polardbClient.DescribeDBClusters(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving Polardb Instances: %s", err)
		}
		resp, _ := raw.(*polardb.DescribeDBClustersResponse)
		if resp == nil || len(resp.Items.DBCluster) < 1 {
			break
		}
		insts = append(insts, resp.Items.DBCluster...)

		if len(resp.Items.DBCluster) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	vpcService := VpcService{client}
	for _, v := range insts {
		name := v.DBClusterDescription
		id := v.DBClusterId
		skip := true
		if !sweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			// If a slb name is set by other service, it should be fetched by vswitch name and deleted.
			if skip {
				if need, err := vpcService.needSweepVpc(v.VpcId, ""); err == nil {
					skip = !need
				}
			}

			if skip {
				log.Printf("[INFO] Skipping Polardb Instance: %s (%s)", name, id)
				continue
			}
		}

		log.Printf("[INFO] Deleting Polardb Instance: %s (%s)", name, id)

		req := polardb.CreateDeleteDBClusterRequest()
		req.DBClusterId = id
		_, err := client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
			return polardbClient.DeleteDBCluster(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Polardb Instance (%s (%s)): %s", name, id, err)
		}
	}

	return nil
}

func TestAccAliCloudPolarDBCluster_Update(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	var ips []map[string]interface{}
	name := "tf-testAccPolarDBClusterUpdate"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       CHECKSET,
		"db_node_class":     CHECKSET,
		"vswitch_id":        CHECKSET,
		"db_type":           CHECKSET,
		"db_version":        CHECKSET,
		"tde_status":        "Disabled",
		"connection_string": REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
		"status":            CHECKSET,
		"create_time":       CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	err := deleteTDEPolicyAndRole()
	assert.Nil(t, err)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBClusterTDEConfigDependence)

	curTime := time.Now().Add(time.Hour * 2)
	endTime := curTime.Add(time.Hour * 48)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.GDNPolarDBSupportRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":           "MySQL",
					"db_version":        "8.0",
					"pay_type":          "PostPaid",
					"db_node_count":     "2",
					"db_node_class":     "${data.alicloud_polardb_node_classes.this.classes.0.supported_engines.0.available_resources.0.db_node_class}",
					"vswitch_id":        "${local.vswitch_id}",
					"description":       "${var.name}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"zone_id":           CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "imci_switch", "sub_category"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testaccPolarDBClusterUpdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testaccPolarDBClusterUpdate",
					}),
				),
			},
			// There is an openapi bug that the maintain_time can not take effect after updating it.
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_time": "16:00Z-17:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_time": "16:00Z-17:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tde_status":         "Enabled",
					"encrypt_new_tables": "ON",
					"encryption_key":     "${alicloud_kms_key.default.id}",
					"role_arn":           "acs:ram::${data.alicloud_account.current.id}:role/aliyunrdsinstanceencryptiondefaultrole",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tde_status":         "Enabled",
						"encrypt_new_tables": "ON",
						"encryption_key":     CHECKSET,
						"role_arn":           CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_node_count": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_count": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_node_count": "2",
					"db_node_class": "${data.alicloud_polardb_node_classes.this.classes.0.supported_engines.0.available_resources.1.db_node_class}",
					"sub_category":  "Exclusive",
					"modify_type":   "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_count": "2",
						"db_node_class": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"collector_status": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"collector_status": "Enable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"upgrade_type":      "ALL",
					"from_time_service": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"upgrade_type":      "ALL",
						"from_time_service": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"upgrade_type":       "ALL",
					"from_time_service":  "false",
					"planned_start_time": curTime.Format("2006-01-02T15:04:05Z"),
					"planned_end_time":   endTime.Format("2006-01-02T15:04:05Z"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"upgrade_type":       "ALL",
						"from_time_service":  "false",
						"planned_start_time": curTime.Format("2006-01-02T15:04:05Z"),
						"planned_end_time":   endTime.Format("2006-01-02T15:04:05Z"),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyValueInMapsForPolarDB(ips, "security ip", "security_ips", "10.168.1.12,100.69.7.112"),
					testAccCheck(map[string]string{
						"security_ips.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_ip_array": []map[string]interface{}{{
						"db_cluster_ip_array_name": "default",
						"security_ips":             []string{"10.168.1.12", "100.69.7.112"},
					}, {
						"db_cluster_ip_array_name": "test_ips1",
						"security_ips":             []string{"10.168.1.13"},
					}, {
						"db_cluster_ip_array_name": "test_ips2",
						"security_ips":             []string{"100.69.7.113"},
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_ip_array.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_ip_array": []map[string]interface{}{{
						"db_cluster_ip_array_name": "default",
						"security_ips":             []string{"10.168.1.12", "100.69.7.112"},
					}, {
						"db_cluster_ip_array_name": "test_ips1",
						"modify_mode":              "Cover",
						"security_ips":             []string{"10.168.1.14"},
					}, {
						"db_cluster_ip_array_name": "test_ips2",
						"modify_mode":              "Delete",
						"security_ips":             []string{"100.69.7.113"},
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_ip_array.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       REMOVEKEY,
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testaccPolarDBClusterUpdate1",
					//"maintain_time": "02:00Z-03:00Z",
					"db_node_count": "2",
					"db_node_class": "${data.alicloud_polardb_node_classes.this.classes.0.supported_engines.0.available_resources.2.db_node_class}",
					"sub_category":  "Exclusive",
					"modify_type":   "Upgrade",
					"security_ips":  []string{"10.168.1.13", "100.69.7.113"},
					"db_cluster_ip_array": []map[string]interface{}{{
						"db_cluster_ip_array_name": "default",
						"modify_mode":              "Cover",
						"security_ips":             []string{"10.168.1.13", "100.69.7.113"},
					}, {
						"db_cluster_ip_array_name": "test_ips1",
						"modify_mode":              "Cover",
						"security_ips":             []string{"10.168.1.14"},
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testaccPolarDBClusterUpdate1",
						//"maintain_time":         "02:00Z-03:00Z",
						"db_node_class":         CHECKSET,
						"security_ips.#":        "2",
						"db_cluster_ip_array.#": "2",
					}),
					testAccCheckKeyValueInMapsForPolarDB(ips, "security ip", "security_ips", "10.168.1.13,100.69.7.113"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": "${alicloud_security_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_lock": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_lock": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_lock": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_lock": "0",
					}),
				),
			},
		},
	})

}

func TestAccAliCloudPolarDBCluster_UpdatePrePaid(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testAccPolarDBClusterUpdatePrePaid"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       CHECKSET,
		"db_node_class":     CHECKSET,
		"vswitch_id":        CHECKSET,
		"db_type":           CHECKSET,
		"db_version":        CHECKSET,
		"tde_status":        "Disabled",
		"connection_string": REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
		"status":            CHECKSET,
		"create_time":       CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBClusterPrePaidConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":           "MySQL",
					"db_version":        "8.0",
					"pay_type":          "PrePaid",
					"creation_category": "Normal",
					"db_node_count":     "2",
					"period":            "1",
					"db_node_class":     "${data.alicloud_polardb_node_classes.this.classes.0.supported_engines.0.available_resources.0.db_node_class}",
					"vswitch_id":        "${local.vswitch_id}",
					"description":       "${var.name}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"storage_type":      "PSL4",
					"compress_storage":  "OFF",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"zone_id":           CHECKSET,
						"storage_type":      "PSL4",
						"compress_storage":  "OFF",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "imci_switch", "sub_category", "period"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renewal_status":    "AutoRenewal",
					"auto_renew_period": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_status":    "AutoRenewal",
						"auto_renew_period": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parameters": []map[string]interface{}{{
						"name":  "wait_timeout",
						"value": "86",
					}, {
						"name":  "innodb_old_blocks_time",
						"value": "10",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pay_type": "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pay_type": "PostPaid",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"compress_storage": "ON",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compress_storage": "ON",
					}),
				),
			},
		},
	})

}

func TestAccAliCloudPolarDBCluster_Multi(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testaccPolarDBClusterMult"
	resourceId := "alicloud_polardb_cluster.default.2"
	var basicMap = map[string]string{
		"description":       CHECKSET,
		"db_node_class":     CHECKSET,
		"vswitch_id":        CHECKSET,
		"db_type":           CHECKSET,
		"db_version":        CHECKSET,
		"connection_string": REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":         "3",
					"db_type":       "MySQL",
					"db_version":    "8.0",
					"pay_type":      "PostPaid",
					"db_node_class": "${data.alicloud_polardb_node_classes.this.classes.0.supported_engines.0.available_resources.0.db_node_class}",
					"vswitch_id":    "${local.vswitch_id}",
					"description":   "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAliCloudPolarDBCluster_CreateGDN(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testAccPolarDBClusterCreateGND"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       CHECKSET,
		"db_node_class":     CHECKSET,
		"vswitch_id":        CHECKSET,
		"db_type":           CHECKSET,
		"db_version":        CHECKSET,
		"connection_string": REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceGDNDBClusterConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.GDNPolarDBSupportRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":           "MySQL",
					"db_version":        "8.0",
					"pay_type":          "PostPaid",
					"db_node_count":     "2",
					"db_node_class":     "${data.alicloud_polardb_node_classes.this.classes.0.supported_engines.0.available_resources.0.db_node_class}",
					"vswitch_id":        "${local.vswitch_id}",
					"creation_category": "Normal",
					"clone_data_point":  "LATEST",
					"creation_option":   "CreateGdnStandby",
					"description":       "${var.name}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"gdn_id":            "${alicloud_polardb_global_database_network.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"zone_id":           CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "creation_option", "gdn_id", "source_resource_id", "clone_data_point"},
			},
		},
	})
}

func TestAccAliCloudPolarDBCluster_CreateNormal(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testAccPolarDBClusterCreateNormal"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       CHECKSET,
		"db_node_class":     CHECKSET,
		"vswitch_id":        CHECKSET,
		"db_type":           CHECKSET,
		"db_version":        CHECKSET,
		"connection_string": REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":                "MySQL",
					"db_version":             "8.0",
					"pay_type":               "PostPaid",
					"db_node_count":          "2",
					"db_node_class":          "${data.alicloud_polardb_node_classes.this.classes.0.supported_engines.0.available_resources.0.db_node_class}",
					"vswitch_id":             "${local.vswitch_id}",
					"zone_id":                "${data.alicloud_polardb_node_classes.this.classes.0.zone_id}",
					"creation_category":      "Normal",
					"loose_polar_log_bin":    "ON",
					"db_node_num":            "2",
					"default_time_zone":      "+1:00",
					"description":            "${var.name}",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"parameter_group_id":     "${data.alicloud_polardb_parameter_groups.default.groups.0.id}",
					"lower_case_table_names": "0",
					"backup_retention_policy_on_cluster_deletion": "NONE",
					"loose_xengine": "OFF",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id":      CHECKSET,
						"zone_id":                CHECKSET,
						"lower_case_table_names": CHECKSET,
						"loose_polar_log_bin":    CHECKSET,
						"default_time_zone":      CHECKSET,
						"loose_xengine":          "OFF",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "creation_option", "db_node_num", "parameter_group_id", "backup_retention_policy_on_cluster_deletion"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"loose_polar_log_bin": "OFF",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"loose_polar_log_bin": "OFF",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"loose_xengine":                "ON",
					"loose_xengine_use_memory_pct": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"loose_xengine":                "ON",
						"loose_xengine_use_memory_pct": "50",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"loose_xengine_use_memory_pct": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"loose_xengine_use_memory_pct": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_time_zone": "+2:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_time_zone": "+2:00",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hot_replica_mode": "OFF",
					"db_node_id":       "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"hot_replica_mode": "OFF",
						"db_node_id":       CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hot_replica_mode": "ON",
					"db_node_id":       "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"hot_replica_mode": "ON",
						"db_node_id":       CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudPolarDBCluster_CreateCloneFromPolarDB(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testAccPolarDBClusterCreateCloneFromPolarDB"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       CHECKSET,
		"db_node_class":     CHECKSET,
		"vswitch_id":        CHECKSET,
		"db_type":           CHECKSET,
		"db_version":        CHECKSET,
		"connection_string": REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCloneFromPolarDBClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.PolarDBCloneFromRdsSupportRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":            "MySQL",
					"db_version":         "8.0",
					"pay_type":           "PostPaid",
					"db_node_count":      "2",
					"db_node_class":      "${data.alicloud_polardb_node_classes.this.classes.0.supported_engines.0.available_resources.0.db_node_class}",
					"vswitch_id":         "${local.vswitch_id}",
					"creation_category":  "Normal",
					"clone_data_point":   "LATEST",
					"db_node_num":        "2",
					"creation_option":    "CloneFromPolarDB",
					"source_resource_id": "${alicloud_polardb_cluster.cluster.id}",
					"description":        "${var.name}",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"parameter_group_id": "${data.alicloud_polardb_parameter_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"zone_id":           CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"creation_category": "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"creation_category": "Normal",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "creation_option", "source_resource_id", "clone_data_point", "db_node_num", "parameter_group_id"},
			},
		},
	})
}

func TestAccAliCloudPolarDBCluster_CreateCloneFromRDS(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testPolarDBClusterCreateCloneFromRDS"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       CHECKSET,
		"db_node_class":     CHECKSET,
		"vswitch_id":        CHECKSET,
		"db_type":           CHECKSET,
		"db_version":        CHECKSET,
		"connection_string": REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCloneOrMigrationFromRDSClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.PolarDBCloneFromRdsSupportRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":            "MySQL",
					"db_version":         "8.0",
					"pay_type":           "PostPaid",
					"db_node_count":      "2",
					"db_node_class":      "${data.alicloud_polardb_node_classes.this.classes.0.supported_engines.0.available_resources.0.db_node_class}",
					"vswitch_id":         "${local.vswitch_id}",
					"description":        "${var.name}",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"clone_data_point":   "LATEST",
					"creation_option":    "CloneFromRDS",
					"creation_category":  "Normal",
					"source_resource_id": "${alicloud_db_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"creation_category": "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"creation_category": "Normal",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "creation_option", "gdn_id", "source_resource_id", "clone_data_point"},
			},
		},
	})
}

func TestAccAliCloudPolarDBCluster_CreateMigrationFromRDS(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testPolarDBClusterMigrationFromRDS"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       CHECKSET,
		"db_node_class":     CHECKSET,
		"vswitch_id":        CHECKSET,
		"db_type":           CHECKSET,
		"db_version":        CHECKSET,
		"connection_string": REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCloneOrMigrationFromRDSClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.PolarDBCloneFromRdsSupportRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":            "MySQL",
					"db_version":         "8.0",
					"pay_type":           "PostPaid",
					"db_node_count":      "2",
					"db_node_class":      "${data.alicloud_polardb_node_classes.this.classes.0.supported_engines.0.available_resources.0.db_node_class}",
					"vswitch_id":         "${local.vswitch_id}",
					"description":        "${var.name}",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"clone_data_point":   "LATEST",
					"creation_option":    "MigrationFromRDS",
					"creation_category":  "Normal",
					"source_resource_id": "${alicloud_db_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"creation_category": "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"creation_category": "Normal",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "creation_option", "gdn_id", "source_resource_id", "clone_data_point"},
			},
		},
	})
}

func TestAccAliCloudPolarDBCluster_VpcId(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testAccPolarDBClusterUpdate"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       CHECKSET,
		"db_node_class":     CHECKSET,
		"vswitch_id":        CHECKSET,
		"db_type":           CHECKSET,
		"db_version":        CHECKSET,
		"tde_status":        "Disabled",
		"connection_string": REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":           "MySQL",
					"db_version":        "8.0",
					"pay_type":          "PostPaid",
					"db_node_count":     "2",
					"db_node_class":     "${data.alicloud_polardb_node_classes.this.classes.0.supported_engines.0.available_resources.0.db_node_class}",
					"vswitch_id":        "${local.vswitch_id}",
					"description":       "${var.name}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"vpc_id":            "${local.vpc_id}",
					"security_ips":      []string{"10.168.1.12", "100.69.7.112", "127.0.0.1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"vpc_id":            CHECKSET,
						"security_ips.#":    "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "imci_switch", "sub_category"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"target_db_revision_version_code": "20230707",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_db_revision_version_code": "20230707",
					}),
				),
			},
		},
	})

}

func TestAccAliCloudPolarDBCluster_NormalMultimaster(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testAccPolarDBCluster"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       CHECKSET,
		"db_node_class":     CHECKSET,
		"vswitch_id":        CHECKSET,
		"db_type":           CHECKSET,
		"db_version":        CHECKSET,
		"tde_status":        "Disabled",
		"connection_string": REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBClusterNormalMultimasterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.PolarDBCloneFromRdsSupportRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":           "MySQL",
					"db_version":        "8.0",
					"pay_type":          "PostPaid",
					"db_node_class":     "polar.mysql.mmx4.medium",
					"vswitch_id":        "${local.vswitch_id}",
					"creation_category": "NormalMultimaster",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pay_type":          "PostPaid",
						"creation_category": "NormalMultimaster",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "creation_option", "gdn_id", "source_resource_id", "clone_data_point"},
			},
		},
	})
}

// Currently, there is no region support creating SENormal
func TestAccAliCloudPolarDBClusterSENormalCreate(t *testing.T) {
	var v map[string]interface{}
	name := "tf-testAccPolarDBClusterSENormalCreate"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":   CHECKSET,
		"db_node_class": CHECKSET,
		"vswitch_id":    CHECKSET,
		"db_type":       CHECKSET,
		"db_version":    CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBClusterSENormalConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.SENormalPolarDBSupportRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":             "MySQL",
					"db_version":          "8.0",
					"pay_type":            "PostPaid",
					"db_node_class":       "polar.mysql.x2.large.c",
					"vswitch_id":          "${local.vswitch_id}",
					"description":         "${var.name}",
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"creation_category":   "SENormal",
					"storage_type":        "ESSDAUTOPL",
					"provisioned_iops":    "1000",
					"storage_space":       "40",
					"db_node_num":         "2",
					"hot_standby_cluster": "ON",
					"storage_pay_type":    "PostPaid",
					"proxy_type":          "EXCLUSIVE",
					"proxy_class":         "polar.maxscale.g2.medium.c",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"zone_id":           CHECKSET,
						"storage_type":      "ESSDAUTOPL",
						"provisioned_iops":  "1000",
						"storage_space":     "40",
						"storage_pay_type":  "PostPaid",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "proxy_type", "proxy_class", "db_node_num"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_space": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_space": "50",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"provisioned_iops": "1200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"provisioned_iops": "1200",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudPolarDBClusterSENormalEqualCreateWithStandbyAz(t *testing.T) {
	var v map[string]interface{}
	name := "tf-testAccPolarDBClusterSENormalEqualCreateWithStandbyAz"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":   CHECKSET,
		"db_node_class": CHECKSET,
		"vswitch_id":    CHECKSET,
		"db_type":       CHECKSET,
		"db_version":    CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBClusterSENormalEqualConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SENormalPolarDBEqualTestRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":             "MySQL",
					"db_version":          "8.0",
					"pay_type":            "PostPaid",
					"db_node_class":       "polar.mysql.x2.large.c",
					"vswitch_id":          "${local.vswitch_id}",
					"description":         "${var.name}",
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"creation_category":   "SENormal",
					"storage_type":        "PSL5",
					"storage_space":       "40",
					"db_node_num":         "1",
					"hot_standby_cluster": "EQUAL",
					"storage_pay_type":    "PostPaid",
					"proxy_type":          "EXCLUSIVE",
					"proxy_class":         "polar.maxscale.g2.medium.c",
					"standby_az":          "cn-hangzhou-k",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"zone_id":           CHECKSET,
						"storage_type":      "PSL5",
						"storage_pay_type":  "PostPaid",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "proxy_type", "proxy_class", "db_node_num"},
			},
		},
	})
}

func TestAccAliCloudPolarDBClusterSENormalEqualCreate(t *testing.T) {
	var v map[string]interface{}
	name := "tf-testAccPolarDBClusterSENormalEqualCreate"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":   CHECKSET,
		"db_node_class": CHECKSET,
		"vswitch_id":    CHECKSET,
		"db_type":       CHECKSET,
		"db_version":    CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBClusterSENormalEqualConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SENormalPolarDBEqualTestRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":             "MySQL",
					"db_version":          "8.0",
					"pay_type":            "PostPaid",
					"db_node_class":       "polar.mysql.x2.large.c",
					"vswitch_id":          "${local.vswitch_id}",
					"description":         "${var.name}",
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"creation_category":   "SENormal",
					"storage_type":        "PSL5",
					"storage_space":       "40",
					"db_node_num":         "2",
					"hot_standby_cluster": "ON",
					"storage_pay_type":    "PostPaid",
					"proxy_type":          "EXCLUSIVE",
					"proxy_class":         "polar.maxscale.g2.medium.c",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"zone_id":           CHECKSET,
						"storage_type":      "PSL5",
						"storage_pay_type":  "PostPaid",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "proxy_type", "proxy_class", "db_node_num"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"standby_az": "cn-hangzhou-i",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"standby_az": "cn-hangzhou-i",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"standby_az": "cn-hangzhou-k",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"standby_az": "cn-hangzhou-k",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudPolarDBCluster_Serverless(t *testing.T) {
	var v map[string]interface{}
	name := "tf-testAccPolarDBClusterServerless"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":   CHECKSET,
		"db_node_class": CHECKSET,
		"vswitch_id":    CHECKSET,
		"db_type":       CHECKSET,
		"db_version":    CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBClusterServerlessConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":                               "MySQL",
					"db_version":                            "8.0",
					"pay_type":                              "PostPaid",
					"db_node_class":                         "polar.mysql.sl.small",
					"vswitch_id":                            "${local.vswitch_id}",
					"description":                           "${var.name}",
					"resource_group_id":                     "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"creation_category":                     "Normal",
					"serverless_type":                       "AgileServerless",
					"scale_min":                             "1",
					"scale_max":                             "2",
					"scale_ro_num_min":                      "1",
					"scale_ro_num_max":                      "2",
					"allow_shut_down":                       "false",
					"serverless_rule_mode":                  "normal",
					"serverless_rule_cpu_shrink_threshold":  "40",
					"serverless_rule_cpu_enlarge_threshold": "70",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id":                     CHECKSET,
						"zone_id":                               CHECKSET,
						"serverless_type":                       "AgileServerless",
						"scale_min":                             "1",
						"scale_max":                             "2",
						"scale_ro_num_min":                      "1",
						"scale_ro_num_max":                      "2",
						"allow_shut_down":                       "false",
						"db_node_class":                         "polar.mysql.sl.small",
						"serverless_rule_mode":                  "normal",
						"serverless_rule_cpu_shrink_threshold":  "40",
						"serverless_rule_cpu_enlarge_threshold": "70",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scale_max": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scale_max": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scale_min": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scale_min": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scale_ro_num_max": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scale_ro_num_max": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scale_ro_num_min": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scale_ro_num_min": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allow_shut_down": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allow_shut_down": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"seconds_until_auto_pause": "3660",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"seconds_until_auto_pause": "3660",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudPolarDBCluster_SteadyServerless(t *testing.T) {
	var v map[string]interface{}
	name := "tf-testAccPolarDBClusterSteadyServerless"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":   CHECKSET,
		"db_node_class": CHECKSET,
		"vswitch_id":    CHECKSET,
		"db_type":       CHECKSET,
		"db_version":    CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBClusterSteadyServerlessConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":           "MySQL",
					"db_version":        "8.0",
					"pay_type":          "PostPaid",
					"db_node_class":     "polar.mysql.x4.medium",
					"vswitch_id":        "${local.vswitch_id}",
					"description":       "${var.name}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"creation_category": "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"zone_id":           CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_node_count": "3",
					"imci_switch":   "ON",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_count": "3",
						"imci_switch":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"serverless_type":          "SteadyServerless",
					"serverless_steady_switch": "ON",
					"scale_min":                "1",
					"scale_max":                "2",
					"scale_ro_num_min":         "1",
					"scale_ro_num_max":         "2",
					"allow_shut_down":          "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"serverless_type":          "SteadyServerless",
						"serverless_steady_switch": "ON",
						"scale_min":                "1",
						"scale_max":                "2",
						"scale_ro_num_min":         "1",
						"scale_ro_num_max":         "2",
						"allow_shut_down":          "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scale_max": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scale_max": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scale_min": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scale_min": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scale_min": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scale_min": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scale_max": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scale_max": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scale_ro_num_max": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scale_ro_num_max": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scale_ro_num_min": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scale_ro_num_min": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scale_ap_ro_num_max": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scale_ap_ro_num_max": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scale_ap_ro_num_min": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scale_ap_ro_num_min": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"serverless_rule_mode": "flexible",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"serverless_rule_mode": "flexible",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"serverless_rule_cpu_shrink_threshold": "40",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"serverless_rule_cpu_shrink_threshold": "40",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"serverless_rule_cpu_enlarge_threshold": "70",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"serverless_rule_cpu_enlarge_threshold": "70",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"serverless_steady_switch": "OFF",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"serverless_steady_switch": "OFF",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudPolarDBCluster_CreateDBCluster(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testAccPolarDBClusterCreateNormal"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       CHECKSET,
		"db_node_class":     CHECKSET,
		"vswitch_id":        CHECKSET,
		"db_type":           CHECKSET,
		"db_version":        CHECKSET,
		"connection_string": REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":                "MySQL",
					"db_version":             "8.0",
					"pay_type":               "PostPaid",
					"db_node_count":          "2",
					"db_node_class":          "${data.alicloud_polardb_node_classes.this.classes.0.supported_engines.0.available_resources.0.db_node_class}",
					"vswitch_id":             "${local.vswitch_id}",
					"zone_id":                "${data.alicloud_polardb_node_classes.this.classes.0.zone_id}",
					"creation_category":      "Normal",
					"db_node_num":            "2",
					"default_time_zone":      "+1:00",
					"description":            "${var.name}",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"parameter_group_id":     "${data.alicloud_polardb_parameter_groups.default.groups.0.id}",
					"lower_case_table_names": "0",
					"backup_retention_policy_on_cluster_deletion": "NONE",
					"storage_type":     "PSL4",
					"db_minor_version": "8.0.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id":      CHECKSET,
						"zone_id":                CHECKSET,
						"lower_case_table_names": CHECKSET,
						"default_time_zone":      CHECKSET,
						"storage_type":           "PSL4",
						"db_minor_version":       "8.0.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"creation_category": "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"creation_category": "Normal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_type": "PSL5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_type": "PSL5",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "creation_option", "db_node_num", "parameter_group_id", "backup_retention_policy_on_cluster_deletion"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_time_zone": "+2:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_time_zone": "+2:00",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudPolarDBCluster_Xengine(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testAccPolarDBCluster-x-engine"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       CHECKSET,
		"db_node_class":     CHECKSET,
		"vswitch_id":        CHECKSET,
		"db_type":           CHECKSET,
		"db_version":        CHECKSET,
		"tde_status":        "Disabled",
		"connection_string": REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBClusterXengineConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":                      "MySQL",
					"db_version":                   "8.0",
					"pay_type":                     "PostPaid",
					"db_node_class":                "${data.alicloud_polardb_node_classes.this.classes.0.supported_engines.0.available_resources.0.db_node_class}",
					"vswitch_id":                   "${local.vswitch_id}",
					"description":                  "${var.name}",
					"creation_category":            "Normal",
					"loose_xengine":                "ON",
					"loose_xengine_use_memory_pct": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pay_type":                     "PostPaid",
						"creation_category":            "Normal",
						"loose_xengine":                "ON",
						"loose_xengine_use_memory_pct": "60",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "creation_option", "gdn_id", "source_resource_id", "clone_data_point"},
			},
		},
	})

}

func TestAccAliCloudPolarDBCluster_CreateRecoverFromRecyclebin(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testPolarDBClusterCreateRecoverFromRecyclebin"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       CHECKSET,
		"db_node_class":     CHECKSET,
		"vswitch_id":        CHECKSET,
		"db_type":           CHECKSET,
		"db_version":        CHECKSET,
		"connection_string": REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceRecoverFromRecyclebinClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":            "MySQL",
					"db_version":         "8.0",
					"pay_type":           "PostPaid",
					"db_node_count":      "2",
					"db_node_class":      "${data.alicloud_polardb_node_classes.this.classes.0.supported_engines.0.available_resources.0.db_node_class}",
					"vswitch_id":         "${local.vswitch_id}",
					"description":        "${var.name}",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"creation_option":    "RecoverFromRecyclebin",
					"creation_category":  "Normal",
					"source_resource_id": "pc-bp1xzqe7d4n1s0ge9",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "creation_option", "gdn_id", "source_resource_id"},
			},
		},
	})
}

func TestAccAliCloudPolarDBCluster_3AZ(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testAccPolarDBCluster-3az"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       CHECKSET,
		"db_node_class":     CHECKSET,
		"vswitch_id":        CHECKSET,
		"db_type":           CHECKSET,
		"db_version":        CHECKSET,
		"tde_status":        "Disabled",
		"connection_string": REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBCluster3ZAConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.Cluster3AZPolarDBSupportRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":             "MySQL",
					"db_version":          "8.0",
					"pay_type":            "PostPaid",
					"db_node_class":       "${data.alicloud_polardb_node_classes.this.classes.0.supported_engines.0.available_resources.0.db_node_class}",
					"vswitch_id":          "${local.vswitch_id}",
					"description":         "${var.name}",
					"creation_category":   "Normal",
					"hot_standby_cluster": "ON",
					"strict_consistency":  "ON",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pay_type":            "PostPaid",
						"creation_category":   "Normal",
						"hot_standby_cluster": "ON",
						"strict_consistency":  "ON",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "creation_option", "gdn_id", "source_resource_id", "clone_data_point"},
			},
		},
	})
}

func testAccCheckKeyValueInMapsForPolarDB(ps []map[string]interface{}, propName, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, policy := range ps {
			if policy[key].(string) != value {
				return fmt.Errorf("DB %s attribute '%s' expected %#v, got %#v", propName, key, value, policy[key])
			}
		}
		return nil
	}
}

func resourcePolarDBCluster3ZAConfigDependence(name string) string {
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
		zone_id = data.alicloud_polardb_node_classes.this.classes.0.zone_id
		vpc_id = alicloud_vpc.default.id
		cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
	}
	locals {
		vpc_id = alicloud_vpc.default.id
		vswitch_id = concat(alicloud_vswitch.default.*.id, [""])[0]
	}

	data "alicloud_polardb_node_classes" "this" {
	  db_type    = "MySQL"
	  db_version = "8.0"
      pay_type   = "PostPaid"
	  category   = "Normal"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
		status = "OK"
	}

	resource "alicloud_security_group" "default" {
		count = 2
		name   = var.name
		vpc_id = alicloud_vpc.default.id
	}
`, name)
}

func resourcePolarDBClusterXengineConfigDependence(name string) string {
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
		zone_id = data.alicloud_polardb_node_classes.this.classes.0.zone_id
		vpc_id = alicloud_vpc.default.id
		cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
	}
	locals {
		vpc_id = alicloud_vpc.default.id
		vswitch_id = concat(alicloud_vswitch.default.*.id, [""])[0]
	}

	data "alicloud_polardb_node_classes" "this" {
	  db_type    = "MySQL"
	  db_version = "8.0"
      pay_type   = "PostPaid"
	  category   = "Normal"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
		status = "OK"
	}

	resource "alicloud_security_group" "default" {
		count = 2
		name   = var.name
		vpc_id = alicloud_vpc.default.id
	}
`, name)
}

func resourcePolarDBClusterPrePaidConfigDependence(name string) string {
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
		zone_id = data.alicloud_polardb_node_classes.this.classes.0.zone_id
		vpc_id = alicloud_vpc.default.id
		cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
	}

	data "alicloud_polardb_node_classes" "this" {
		  db_type    = "MySQL"
		  db_version = "8.0"
		  pay_type   = "PrePaid"
		  category   = "Normal"
	}

	locals {
		vpc_id = alicloud_vpc.default.id
		vswitch_id = concat(alicloud_vswitch.default.*.id, [""])[0]
	}

	data "alicloud_resource_manager_resource_groups" "default" {
		status = "OK"
	}

`, name)
}

func resourcePolarDBClusterTDEConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
	}

    data "alicloud_polardb_node_classes" "this" {
	  db_type    = "MySQL"
	  db_version = "8.0"
      pay_type   = "PostPaid"
	  category   = "Normal"
	}
	
	resource "alicloud_vpc" "default" {
    	vpc_name = var.name
	}

	resource "alicloud_vswitch" "default" {
		zone_id = data.alicloud_polardb_node_classes.this.classes.0.zone_id
		vpc_id = alicloud_vpc.default.id
		cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
	}
	
	locals {
		vpc_id = alicloud_vpc.default.id
		vswitch_id = concat(alicloud_vswitch.default.*.id, [""])[0]
	}

	data "alicloud_resource_manager_resource_groups" "default" {
		status = "OK"
	}

	resource "alicloud_security_group" "default" {
		count = 2
		name   = var.name
		vpc_id = alicloud_vpc.default.id
	}

    data "alicloud_account" "current" {
    }

    resource "alicloud_kms_key" "default" {
        description             =  var.name
        pending_window_in_days =  7
        status                  = "Enabled"
    }
    resource "alicloud_ram_role" "default" {
      name        = "AliyunRDSInstanceEncryptionDefaultRole"
	  document    = <<DEFINITION
		{
			"Statement": [
				{
					"Action": "sts:AssumeRole",
					"Effect": "Allow",
					"Principal": {
						"Service": [
							"rds.aliyuncs.com"
						]
					}
				}
			],
			"Version": "1"
		}
		DEFINITION
	  description = "RDS使用此角色来访问您在其他云产品中的资源"
    }
    resource "alicloud_resource_manager_policy_attachment" "default" {
	    policy_name       = "AliyunRDSInstanceEncryptionRolePolicy"
	    policy_type       = "System"
	    principal_name    = "${alicloud_ram_role.default.name}@role.${data.alicloud_account.current.id}.onaliyunservice.com"
	    principal_type    = "ServiceRole"
	    resource_group_id = "${data.alicloud_account.current.id}"
    }

`, name)
}

func resourcePolarDBClusterConfigDependence(name string) string {
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
		zone_id = data.alicloud_polardb_node_classes.this.classes.0.zone_id
		vpc_id = alicloud_vpc.default.id
		cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
	}
	
	locals {
		vpc_id = alicloud_vpc.default.id
		vswitch_id = concat(alicloud_vswitch.default.*.id, [""])[0]
	}

    data "alicloud_polardb_parameter_groups" "default" {
          db_type = "MySQL"
          db_version = "8.0"
    }

	data "alicloud_polardb_node_classes" "this" {
	  db_type    = "MySQL"
	  db_version = "8.0"
      pay_type   = "PostPaid"
	  category   = "Normal"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
		status = "OK"
	}

	resource "alicloud_security_group" "default" {
		count = 2
		name   = var.name
		vpc_id = alicloud_vpc.default.id
	}

`, name)
}

func resourceGDNDBClusterConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	data "alicloud_polardb_node_classes" "this" {
	  db_type    = "MySQL"
	  db_version = "8.0"
      pay_type   = "PostPaid"
	  category   = "Normal"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
		status = "OK"
	}

	resource "alicloud_security_group" "default" {
		count = 2
		name   = "${var.name}"
	    vpc_id = "${local.vpc_id}"
	}
	
	resource "alicloud_polardb_cluster" "cluster" {
		db_type = "MySQL"
		db_version = "8.0"
		pay_type = "PostPaid"
        db_node_count = "2"
		db_node_class = "polar.mysql.x4.large"
		vswitch_id = "${local.vswitch_id}"
		description = "${var.name}"
	}

	resource "alicloud_polardb_global_database_network" "default" {
		db_cluster_id = "${alicloud_polardb_cluster.cluster.id}"
		description   = var.name
	}
`, PolarDBCommonTestCase, name)
}

func resourceCloneFromPolarDBClusterConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	data "alicloud_polardb_node_classes" "this" {
	  db_type    = "MySQL"
	  db_version = "8.0"
      pay_type   = "PostPaid"
	  category   = "Normal"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
		status = "OK"
	}

    data "alicloud_polardb_parameter_groups" "default" {
          db_type = "MySQL"
          db_version = "8.0"
    }

	resource "alicloud_security_group" "default" {
		count = 2
		name   = "${var.name}"
	    vpc_id = "${local.vpc_id}"
	}
	
	resource "alicloud_polardb_cluster" "cluster" {
		db_type = "MySQL"
		db_version = "8.0"
		pay_type = "PostPaid"
        db_node_count = "2"
		db_node_class = "polar.mysql.x4.large"
		vswitch_id = "${local.vswitch_id}"
		description = "${var.name}"
	}
`, PolarDBCommonTestCase, name)
}

func resourceCloneOrMigrationFromRDSClusterConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

    data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
	}

	data "alicloud_polardb_node_classes" "this" {
	  db_type    = "MySQL"
	  db_version = "8.0"
      pay_type   = "PostPaid"
      category   = "Normal"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
		status = "OK"
	}

	resource "alicloud_security_group" "default" {
		count = 2
		name   = "${var.name}"
	    vpc_id = alicloud_vpc.default.id
	}
	
	resource "alicloud_vpc" "default" {
    	vpc_name = var.name
	}

	resource "alicloud_vswitch" "default" {
		vpc_id = alicloud_vpc.default.id
	  	zone_id = "cn-beijing-k"
		cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
	}

	locals {
		vpc_id = alicloud_vpc.default.id
		vswitch_id = concat(alicloud_vswitch.default.*.id, [""])[0]
	}

	resource "alicloud_db_database" "default" {
		instance_id = alicloud_db_instance.default.id
  		name        = "testdb01"
	}

	resource "alicloud_db_instance" "default" {
		engine = "MySQL"
		engine_version = "8.0"
		db_instance_storage_type = "local_ssd"
		instance_charge_type = "Postpaid"
		instance_type = "mysql.x8.medium.2"
		instance_storage = "20"
		vswitch_id = "${local.vswitch_id}"
		instance_name = "tf-testAccDBInstance"
        zone_id = "cn-beijing-k"
    }
`, name)
}

func resourcePolarDBClusterNormalMultimasterConfigDependence(name string) string {
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
		zone_id = "cn-beijing-k"
		vpc_id = alicloud_vpc.default.id
		cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
	}
	locals {
		vpc_id = alicloud_vpc.default.id
		vswitch_id = concat(alicloud_vswitch.default.*.id, [""])[0]
	}

	data "alicloud_resource_manager_resource_groups" "default" {
		status = "OK"
	}

	resource "alicloud_security_group" "default" {
		count = 2
		name   = var.name
		vpc_id = alicloud_vpc.default.id
	}

`, name)
}

func resourcePolarDBClusterSENormalConfigDependence(name string) string {
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
		zone_id = data.alicloud_polardb_node_classes.this.classes.0.zone_id
		vpc_id = alicloud_vpc.default.id
		cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
	}

	locals {
		vpc_id = alicloud_vpc.default.id
		vswitch_id = concat(alicloud_vswitch.default.*.id, [""])[0]
	}

	data "alicloud_polardb_node_classes" "this" {
	  db_type    = "MySQL"
	  db_version = "8.0"
      pay_type   = "PostPaid"
	  category   = "SENormal"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
		status = "OK"
	}

	resource "alicloud_security_group" "default" {
		count = 2
		name   = var.name
		vpc_id = alicloud_vpc.default.id
	}

`, name)
}

func resourcePolarDBClusterSENormalEqualConfigDependence(name string) string {
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
		zone_id = data.alicloud_polardb_node_classes.this.classes.0.zone_id
		vpc_id = alicloud_vpc.default.id
		cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
	}

	locals {
		vpc_id = alicloud_vpc.default.id
		vswitch_id = concat(alicloud_vswitch.default.*.id, [""])[0]
	}

	data "alicloud_polardb_node_classes" "this" {
	  db_type    = "MySQL"
	  db_version = "8.0"
      pay_type   = "PostPaid"
	  category   = "SENormal"
	  zone_id    = "cn-hangzhou-j"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
		status = "OK"
	}

	resource "alicloud_security_group" "default" {
		count = 2
		name   = var.name
		vpc_id = alicloud_vpc.default.id
	}

`, name)
}

func resourcePolarDBClusterServerlessConfigDependence(name string) string {
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
		zone_id = data.alicloud_polardb_node_classes.this.classes.0.zone_id
		vpc_id = alicloud_vpc.default.id
		cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
	}
	
	locals {
		vpc_id = alicloud_vpc.default.id
		vswitch_id = concat(alicloud_vswitch.default.*.id, [""])[0]
	}
	
	data "alicloud_polardb_node_classes" "this" {
		db_type    = "MySQL"
		db_version = "8.0"
		pay_type   = "PostPaid"
		category   = "Normal"
	}
	
	data "alicloud_resource_manager_resource_groups" "default" {
		status = "OK"
	}
	
	resource "alicloud_security_group" "default" {
		count = 2
		name   = var.name
		vpc_id = alicloud_vpc.default.id
	}

`, name)
}

func resourcePolarDBClusterSteadyServerlessConfigDependence(name string) string {
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
		zone_id = data.alicloud_polardb_node_classes.this.classes.0.zone_id
		vpc_id = alicloud_vpc.default.id
		cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
	}
	
	locals {
		vpc_id = alicloud_vpc.default.id
		vswitch_id = concat(alicloud_vswitch.default.*.id, [""])[0]
	}

	data "alicloud_polardb_node_classes" "this" {
		db_type    = "MySQL"
		db_version = "8.0"
		pay_type   = "PostPaid"
		category   = "Normal"
	}
	
	data "alicloud_resource_manager_resource_groups" "default" {
		status = "OK"
	}
	
	resource "alicloud_security_group" "default" {
		count = 2
		name   = var.name
		vpc_id = alicloud_vpc.default.id
	}

`, name)
}

func resourceRecoverFromRecyclebinClusterConfigDependence(name string) string {
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
		zone_id = data.alicloud_polardb_node_classes.this.classes.0.zone_id
		vpc_id = alicloud_vpc.default.id
		cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
	}

    locals {
		vpc_id = alicloud_vpc.default.id
		vswitch_id = concat(alicloud_vswitch.default.*.id, [""])[0]
	}

	data "alicloud_polardb_node_classes" "this" {
	    db_type    = "MySQL"
	    db_version = "8.0"
        pay_type   = "PostPaid"
        category   = "Normal"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
		status = "OK"
	}

	resource "alicloud_security_group" "default" {
		count = 2
		name   = "${var.name}"
	    vpc_id = alicloud_vpc.default.id
	}
`, name)
}

func deleteTDEPolicyAndRole() error {
	// 删除策略
	p := Provider().(*schema.Provider).ResourcesMap
	ramPolicyExisted, _ := schema.InternalMap(p["alicloud_resource_manager_role"].Schema).Data(nil, nil)
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return err
	}
	accountId, err := rawClient.(*connectivity.AliyunClient).AccountId()
	policyId := fmt.Sprintf("AliyunRDSInstanceEncryptionRolePolicy:System:AliyunRDSInstanceEncryptionDefaultRole@role.%s.onaliyunservice.com:ServiceRole:%s", accountId, accountId)
	ramPolicyExisted.SetId(policyId)
	err = resourceAlicloudResourceManagerPolicyAttachmentDelete(ramPolicyExisted, rawClient)
	if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist.Role", "EntityNotExist.Policy"}) {
		return err
	}
	ramRoleExisted, _ := schema.InternalMap(p["alicloud_ram_role"].Schema).Data(nil, nil)
	ramRoleExisted.SetId("AliyunRDSInstanceEncryptionDefaultRole")
	err = deleteRamRole(ramRoleExisted, rawClient)
	if err != nil {
		return err
	}
	return nil
}

func deleteRamRole(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteRole"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RoleName"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Ram", "2015-05-01", action, query, request, true)

		if err != nil {
			if IsExpectedErrors(err, []string{"DeleteConflict.Role.Policy"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Role"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
