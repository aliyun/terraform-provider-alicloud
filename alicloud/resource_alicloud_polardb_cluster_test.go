package alicloud

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
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

func TestAccAlicloudPolarDBCluster_Update(t *testing.T) {
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
					"vswitch_id":        "${data.alicloud_vswitches.default.vswitches.0.id}",
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
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"maintain_time": "16:00Z-17:00Z",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"maintain_time": "16:00Z-17:00Z",
			//		}),
			//	),
			//},
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
					"imci_switch":   "ON",
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
						"security_ips":             []string{"10.168.1.13", "100.69.7.113"},
					}, {
						"db_cluster_ip_array_name": "test_ips1",
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

func TestAccAlicloudPolarDBCluster_Multi(t *testing.T) {
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
					"vswitch_id":    "${data.alicloud_vswitches.default.vswitches.0.id}",
					"description":   "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func TestAccAlicloudPolarDBCluster_Create(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testAccPolarDBClusterCreate"
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
					"creation_option":    "CloneFromPolarDB",
					"source_resource_id": "${alicloud_polardb_cluster.cluster.id}",
					"description":        "${var.name}",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
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
				ImportStateVerifyIgnore: []string{"modify_type", "creation_option", "gdn_id", "source_resource_id", "clone_data_point"},
			},
		},
	})
}

func TestAccAlicloudPolarDBCluster_CreateCloneFromRDS(t *testing.T) {
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

func TestAccAlicloudPolarDBCluster_CreateMigrationFromRDS(t *testing.T) {
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

func TestAccAlicloudPolarDBCluster_VpcId(t *testing.T) {
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
					"vswitch_id":        "${data.alicloud_vswitches.default.vswitches.0.id}",
					"description":       "${var.name}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"vpc_id":            "${data.alicloud_vpcs.default.ids.0}",
					"security_ips":      []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"vpc_id":            CHECKSET,
						"security_ips.#":    "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "imci_switch", "sub_category"},
			},
		},
	})

}

func TestAccAlicloudPolarDBCluster_NormalMultimaster(t *testing.T) {
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
					"db_node_class":     "${data.alicloud_polardb_node_classes.this.classes.0.supported_engines.0.available_resources.0.db_node_class}",
					"vswitch_id":        "${data.alicloud_vswitches.default.vswitches.0.id}",
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
func SkipTestAccAlicloudPolarDBClusterSENormalCreate(t *testing.T) {
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
					"db_node_class":       "${data.alicloud_polardb_node_classes.this.classes.0.supported_engines.0.available_resources.0.db_node_class}",
					"vswitch_id":          "${data.alicloud_vswitches.default.vswitches.0.id}",
					"description":         "${var.name}",
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"creation_category":   "SENormal",
					"storage_type":        "ESSDPL1",
					"storage_space":       "20",
					"hot_standby_cluster": "ON",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"zone_id":           CHECKSET,
						"storage_type":      "ESSDPL1",
						"storage_space":     "20",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "hot_standby_cluster"},
			},
		},
	})
}

func TestAccAlicloudPolarDBCluster_Serverless(t *testing.T) {
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
					"db_type":           "MySQL",
					"db_version":        "8.0",
					"pay_type":          "PostPaid",
					"db_node_class":     "polar.mysql.sl.small",
					"vswitch_id":        "${data.alicloud_vswitches.default.vswitches.0.id}",
					"description":       "${var.name}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"creation_category": "Normal",
					"serverless_type":   "AgileServerless",
					"scale_min":         "1",
					"scale_max":         "2",
					"scale_ro_num_min":  "1",
					"scale_ro_num_max":  "2",
					"allow_shut_down":   "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"zone_id":           CHECKSET,
						"serverless_type":   "AgileServerless",
						"scale_min":         "1",
						"scale_max":         "2",
						"scale_ro_num_min":  "1",
						"scale_ro_num_max":  "2",
						"allow_shut_down":   "false",
						"db_node_class":     "polar.mysql.sl.small",
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

func resourcePolarDBClusterTDEConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
	}
	data "alicloud_vswitches" "default" {
		zone_id = data.alicloud_polardb_node_classes.this.classes.0.zone_id
		vpc_id = data.alicloud_vpcs.default.ids.0
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
		vpc_id = data.alicloud_vpcs.default.ids.0
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
	data "alicloud_vswitches" "default" {
		zone_id = data.alicloud_polardb_node_classes.this.classes.0.zone_id
		vpc_id = data.alicloud_vpcs.default.ids.0
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
		vpc_id = data.alicloud_vpcs.default.ids.0
	}

`, name)
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
	  zone_id    = local.zone_id
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
`, PolarDBCommonTestCase, name)
}

func resourceCloneOrMigrationFromRDSClusterConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	data "alicloud_polardb_node_classes" "this" {
	  db_type    = "MySQL"
	  db_version = "8.0"
      pay_type   = "PostPaid"
	  zone_id    = local.zone_id
	}

	data "alicloud_resource_manager_resource_groups" "default" {
		status = "OK"
	}

	resource "alicloud_security_group" "default" {
		count = 2
		name   = "${var.name}"
	    vpc_id = "${local.vpc_id}"
	}

	resource "alicloud_db_instance" "default" {
		engine = "MySQL"
		engine_version = "8.0"
		db_instance_storage_type = "local_ssd"
		instance_charge_type = "Postpaid"
		instance_type = "mysql.x8.medium.2"
		instance_storage = "20"
		vswitch_id = local.vswitch_id
		instance_name = "tf-testAccDBInstance"
    }
`, PolarDBCommonTestCase, name)
}

func resourcePolarDBClusterNormalMultimasterConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
	}
	data "alicloud_vswitches" "default" {
		zone_id = data.alicloud_polardb_node_classes.this.classes.0.zone_id
		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	data "alicloud_polardb_node_classes" "this" {
	  db_type    = "MySQL"
	  db_version = "8.0"
      pay_type   = "PostPaid"
	  category   = "NormalMultimaster"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
		status = "OK"
	}

	resource "alicloud_security_group" "default" {
		count = 2
		name   = var.name
		vpc_id = data.alicloud_vpcs.default.ids.0
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
	data "alicloud_vswitches" "default" {
		zone_id = data.alicloud_polardb_node_classes.this.classes.0.zone_id
		vpc_id = data.alicloud_vpcs.default.ids.0
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
		vpc_id = data.alicloud_vpcs.default.ids.0
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

	data "alicloud_vswitches" "default" {
		zone_id = data.alicloud_polardb_node_classes.this.classes.0.zone_id
		vpc_id = data.alicloud_vpcs.default.ids.0
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
		vpc_id = data.alicloud_vpcs.default.ids.0
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
	err = resourceAlicloudRamRoleDelete(ramRoleExisted, rawClient)
	if err != nil {
		return err
	}
	return nil
}
