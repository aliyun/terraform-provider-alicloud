package polardb_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	polardbsvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/polardb"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkterraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var clusterConnectionStringRegexp = "^[a-z-A-Z-0-9]+.rwlb.([a-z-A-Z-0-9]+.){0,1}rds.aliyuncs.com"

func TestAccAliCloudPolarDBCluster_Distributed(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testAccClusterDistributed"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":   acctest.CHECKSET,
		"cn_node_class": acctest.CHECKSET,
		"cn_node_num":   acctest.CHECKSET,
		"cn_node_ids.#": acctest.CHECKSET,
		"dn_node_class": acctest.CHECKSET,
		"dn_node_num":   acctest.CHECKSET,
		"dn_node_ids.#": acctest.CHECKSET,
		"vswitch_id":    acctest.CHECKSET,
		"db_type":       "PostgreSQL",
		"db_version":    "16",
		"status":        acctest.CHECKSET,
		"create_time":   acctest.CHECKSET,
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceClusterDistributedConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Shanghai))
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":       "PostgreSQL",
					"db_version":    "16",
					"pay_type":      "PostPaid",
					"cn_node_class": "${local.base_node_class}",
					"dn_node_class": "${local.base_node_class}",
					"cn_node_num":   1,
					"dn_node_num":   3,
					"vswitch_id":    "${local.vswitch_id}",
					"vpc_id":        "${local.vpc_id}",
					"description":   "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cn_node_class": "polar.pg.x4.medium",
						"cn_node_num":   "1",
						"dn_node_class": "polar.pg.x4.medium",
						"dn_node_num":   "3",
					}),
				),
			},
			{
				// Declaratively add one CN and remove one DN node.
				Config: testAccConfig(map[string]interface{}{
					"db_type":       "PostgreSQL",
					"db_version":    "16",
					"pay_type":      "PostPaid",
					"cn_node_class": "${local.base_node_class}",
					"dn_node_class": "${local.base_node_class}",
					"cn_node_num":   2,
					"dn_node_num":   2,
					"vswitch_id":    "${local.vswitch_id}",
					"vpc_id":        "${local.vpc_id}",
					"description":   "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cn_node_num":   "2",
						"dn_node_num":   "2",
						"cn_node_class": "polar.pg.x4.medium",
						"dn_node_class": "polar.pg.x4.medium",
					}),
				),
			},
			{
				// Upgrade both node types and declaratively remove the added nodes.
				Config: testAccConfig(map[string]interface{}{
					"cn_node_class": "${local.target_node_class}",
					"dn_node_class": "${local.target_node_class}",
					"cn_node_num":   1,
					"dn_node_num":   2,
					"modify_type":   "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cn_node_class": "polar.pg.x4.large",
						"dn_node_class": "polar.pg.x4.large",
						"cn_node_num":   "1",
						"dn_node_num":   "2",
					}),
				),
			},
			{
				// Downgrade both node types back to the original class.
				Config: testAccConfig(map[string]interface{}{
					"cn_node_class": "${local.base_node_class}",
					"dn_node_class": "${local.base_node_class}",
					"modify_type":   "Downgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cn_node_class": "polar.pg.x4.medium",
						"dn_node_class": "polar.pg.x4.medium",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type"},
			},
		},
	})
}

func resourceClusterDistributedConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_vswitches" "default" {
		status  = "Available"
	}

	locals {
		vpc_id           = data.alicloud_vswitches.default.vswitches.0.vpc_id
		vswitch_id       = data.alicloud_vswitches.default.vswitches.0.id
		base_node_class   = "polar.pg.x4.medium"
		target_node_class = "polar.pg.x4.large"
	}
`, name)
}

func TestAccAliCloudPolarDBCluster_Update(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	var ips []map[string]interface{}
	name := "tf-testAccClusterUpdate"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       acctest.CHECKSET,
		"db_node_class":     acctest.CHECKSET,
		"vswitch_id":        acctest.CHECKSET,
		"db_type":           acctest.CHECKSET,
		"db_version":        acctest.CHECKSET,
		"tde_status":        "Disabled",
		"connection_string": acctest.REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
		"status":            acctest.CHECKSET,
		"create_time":       acctest.CHECKSET,
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	//err := deleteTDEPolicyAndRole()
	//assert.Nil(t, err)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceClusterTDEConfigDependence)

	curTime := time.Now().Add(time.Hour * 2)
	endTime := curTime.Add(time.Hour * 48)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":           "MySQL",
					"db_version":        "8.0",
					"pay_type":          "PostPaid",
					"db_node_count":     "2",
					"db_node_class":     "polar.mysql.x4.large",
					"vswitch_id":        "${local.vswitch_id}",
					"description":       "${var.name}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
						"zone_id":           acctest.CHECKSET,
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
					"description": "tf-testaccClusterUpdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testaccClusterUpdate",
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
					"tde_status":                "Enabled",
					"encrypt_new_tables":        "ON",
					"encryption_key":            "${alicloud_kms_key.default.id}",
					"role_arn":                  "acs:ram::${data.alicloud_account.current.id}:role/aliyunrdsinstanceencryptiondefaultrole",
					"enable_automatic_rotation": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tde_status":         "Enabled",
						"encrypt_new_tables": "ON",
						"encryption_key":     acctest.CHECKSET,
						"role_arn":           acctest.CHECKSET,
						"automatic_rotation": acctest.CHECKSET,
						"rotation_interval":  acctest.CHECKSET,
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
					"db_node_class": "polar.pg.x4.medium",
					"sub_category":  "Exclusive",
					"modify_type":   "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_count": "2",
						"db_node_class": acctest.CHECKSET,
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
					testAccCheckKeyValueInMapsForCluster(ips, "security ip", "security_ips", "10.168.1.12,100.69.7.112"),
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
				ExpectNonEmptyPlan: true,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       acctest.REMOVEKEY,
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testaccClusterUpdate1",
					//"maintain_time": "02:00Z-03:00Z",
					"db_node_count": "2",
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
						"description": "tf-testaccClusterUpdate1",
						//"maintain_time":         "02:00Z-03:00Z",
						"db_node_class":         acctest.CHECKSET,
						"security_ips.#":        "2",
						"db_cluster_ip_array.#": "2",
					}),
					testAccCheckKeyValueInMapsForCluster(ips, "security ip", "security_ips", "10.168.1.13,100.69.7.113"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": "${alicloud_security_group.default.*.id}",
					"deletion_lock":      "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "2",
						"deletion_lock":        "1",
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

// TestAccAliCloudPolarDBCluster_DeletePrePaid verifies that destroying a
// subscription cluster releases it through the billing refund API instead of
// leaving it behind.
func TestAccAliCloudPolarDBCluster_DeletePrePaid(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testAccClusterDeletePrePaid"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":   acctest.CHECKSET,
		"db_node_class": acctest.CHECKSET,
		"vswitch_id":    acctest.CHECKSET,
		"db_type":       acctest.CHECKSET,
		"db_version":    acctest.CHECKSET,
		"pay_type":      "PrePaid",
		"status":        acctest.CHECKSET,
		"create_time":   acctest.CHECKSET,
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceClusterPrePaidConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":           "MySQL",
					"db_version":        "8.0",
					"pay_type":          "PrePaid",
					"creation_category": "Normal",
					"db_node_count":     "2",
					"period":            "1",
					"db_node_class":     "polar.mysql.x4.large",
					"vswitch_id":        "${local.vswitch_id}",
					"description":       "${var.name}",
					"storage_type":      "PSL4",
					"compress_storage":  "OFF",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pay_type":         "PrePaid",
						"storage_type":     "PSL4",
						"compress_storage": "OFF",
					}),
				),
			},
		},
	})

}

func TestAccAliCloudPolarDBCluster_UpdatePrePaid(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testAccClusterUpdatePrePaid"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       acctest.CHECKSET,
		"db_node_class":     acctest.CHECKSET,
		"vswitch_id":        acctest.CHECKSET,
		"db_type":           acctest.CHECKSET,
		"db_version":        acctest.CHECKSET,
		"tde_status":        "Disabled",
		"connection_string": acctest.REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
		"status":            acctest.CHECKSET,
		"create_time":       acctest.CHECKSET,
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceClusterPrePaidConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":           "MySQL",
					"db_version":        "8.0",
					"pay_type":          "PrePaid",
					"creation_category": "Normal",
					"db_node_count":     "2",
					"period":            "1",
					"db_node_class":     "polar.mysql.x4.large",
					"vswitch_id":        "${local.vswitch_id}",
					"description":       "${var.name}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"storage_type":      "PSL4",
					"compress_storage":  "OFF",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
						"zone_id":           acctest.CHECKSET,
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
	name := "tf-testaccClusterMult"
	resourceId := "alicloud_polardb_cluster.default.2"
	var basicMap = map[string]string{
		"description":       acctest.CHECKSET,
		"db_node_class":     acctest.CHECKSET,
		"vswitch_id":        acctest.CHECKSET,
		"db_type":           acctest.CHECKSET,
		"db_version":        acctest.CHECKSET,
		"connection_string": acctest.REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":         "3",
					"db_type":       "MySQL",
					"db_version":    "8.0",
					"pay_type":      "PostPaid",
					"db_node_class": "polar.mysql.x4.large",
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
	name := "tf-testAccClusterCreateGND"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       acctest.CHECKSET,
		"db_node_class":     acctest.CHECKSET,
		"vswitch_id":        acctest.CHECKSET,
		"db_type":           acctest.CHECKSET,
		"db_version":        acctest.CHECKSET,
		"connection_string": acctest.REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceGDNDBClusterConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":           "MySQL",
					"db_version":        "8.0",
					"pay_type":          "PostPaid",
					"db_node_count":     "2",
					"db_node_class":     "polar.mysql.x4.large",
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
						"resource_group_id": acctest.CHECKSET,
						"zone_id":           acctest.CHECKSET,
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
	name := "tf-testAccClusterCreateNormal"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       acctest.CHECKSET,
		"db_node_class":     acctest.CHECKSET,
		"vswitch_id":        acctest.CHECKSET,
		"db_type":           acctest.CHECKSET,
		"db_version":        acctest.CHECKSET,
		"connection_string": acctest.REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":             "MySQL",
					"db_version":          "8.0",
					"pay_type":            "PostPaid",
					"db_node_count":       "2",
					"db_node_class":       "polar.mysql.x4.large",
					"vswitch_id":          "${local.vswitch_id}",
					"zone_id":             "${data.alicloud_polardb_node_classes.this.classes.0.zone_id}",
					"creation_category":   "Normal",
					"loose_polar_log_bin": "ON",
					"db_node_num":         "2",
					"default_time_zone":   "+1:00",
					"description":         "${var.name}",
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"parameter_group_id":  "${data.alicloud_polardb_parameter_groups.default.groups.0.id}",
					"backup_retention_policy_on_cluster_deletion": "NONE",
					"loose_xengine": "OFF",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id":      acctest.CHECKSET,
						"zone_id":                acctest.CHECKSET,
						"lower_case_table_names": "1",
						"loose_polar_log_bin":    acctest.CHECKSET,
						"default_time_zone":      acctest.CHECKSET,
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
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_ip_array.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_security_group_list": []string{"${alicloud_polardb_global_security_ip_group.default.id}", "${alicloud_polardb_global_security_ip_group.default2.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_security_group_list.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_security_group_list": []string{"${alicloud_polardb_global_security_ip_group.default2.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_security_group_list.#": "1",
					}),
				),
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
						"db_node_id":       acctest.CHECKSET,
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
						"db_node_id":       acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudPolarDBCluster_GlobalSecurityGroupOnCreate(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccDBCluGSG%d", rand)
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       acctest.CHECKSET,
		"db_node_class":     acctest.CHECKSET,
		"vswitch_id":        acctest.CHECKSET,
		"db_type":           acctest.CHECKSET,
		"db_version":        acctest.CHECKSET,
		"connection_string": acctest.REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":           "MySQL",
					"db_version":        "8.0",
					"pay_type":          "PostPaid",
					"db_node_count":     "2",
					"db_node_class":     "polar.mysql.x4.large",
					"vswitch_id":        "${local.vswitch_id}",
					"zone_id":           "${data.alicloud_polardb_node_classes.this.classes.0.zone_id}",
					"creation_category": "Normal",
					"description":       "${var.name}",
					"global_security_group_list": []string{
						"${alicloud_polardb_global_security_ip_group.default.id}",
						"${alicloud_polardb_global_security_ip_group.default2.id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_security_group_list.#": "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "creation_option", "source_resource_id", "clone_data_point", "db_node_num", "parameter_group_id"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_security_group_list": []string{
						"${alicloud_polardb_global_security_ip_group.default.id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_security_group_list.#": "1",
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

func TestAccAliCloudPolarDBCluster_CreateCloneFromPolarDB(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testAccClusterCreateCloneFromPolarDB"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       acctest.CHECKSET,
		"db_node_class":     acctest.CHECKSET,
		"vswitch_id":        acctest.CHECKSET,
		"db_type":           acctest.CHECKSET,
		"db_version":        acctest.CHECKSET,
		"connection_string": acctest.REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCloneFromClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":            "MySQL",
					"db_version":         "8.0",
					"pay_type":           "PostPaid",
					"db_node_count":      "2",
					"db_node_class":      "polar.mysql.x4.large",
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
						"resource_group_id": acctest.CHECKSET,
						"zone_id":           acctest.CHECKSET,
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
	name := "tf-testClusterCreateCloneFromRDS"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       acctest.CHECKSET,
		"db_node_class":     acctest.CHECKSET,
		"vswitch_id":        acctest.CHECKSET,
		"db_type":           acctest.CHECKSET,
		"db_version":        acctest.CHECKSET,
		"connection_string": acctest.REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCloneOrMigrationFromRDSClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":            "MySQL",
					"db_version":         "8.0",
					"pay_type":           "PostPaid",
					"db_node_count":      "2",
					"db_node_class":      "polar.mysql.x4.large",
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
						"resource_group_id": acctest.CHECKSET,
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
	name := "tf-testClusterMigrationFromRDS"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       acctest.CHECKSET,
		"db_node_class":     acctest.CHECKSET,
		"vswitch_id":        acctest.CHECKSET,
		"db_type":           acctest.CHECKSET,
		"db_version":        acctest.CHECKSET,
		"connection_string": acctest.REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCloneOrMigrationFromRDSClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":            "MySQL",
					"db_version":         "8.0",
					"pay_type":           "PostPaid",
					"db_node_count":      "2",
					"db_node_class":      "polar.mysql.x4.large",
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
						"resource_group_id": acctest.CHECKSET,
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
	name := "tf-testAccClusterUpdate"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       acctest.CHECKSET,
		"db_node_class":     acctest.CHECKSET,
		"vswitch_id":        acctest.CHECKSET,
		"db_type":           acctest.CHECKSET,
		"db_version":        acctest.CHECKSET,
		"tde_status":        "Disabled",
		"connection_string": acctest.REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":           "MySQL",
					"db_version":        "8.0",
					"pay_type":          "PostPaid",
					"db_node_count":     "2",
					"db_node_class":     "polar.mysql.x4.large",
					"vswitch_id":        "${local.vswitch_id}",
					"description":       "${var.name}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"vpc_id":            "${local.vpc_id}",
					"security_ips":      []string{"10.168.1.12", "100.69.7.112", "127.0.0.1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
						"vpc_id":            acctest.CHECKSET,
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
	name := "tf-testAccCluster"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       acctest.CHECKSET,
		"db_node_class":     acctest.CHECKSET,
		"vswitch_id":        acctest.CHECKSET,
		"db_type":           acctest.CHECKSET,
		"db_version":        acctest.CHECKSET,
		"tde_status":        "Disabled",
		"connection_string": acctest.REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceClusterNormalMultimasterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Beijing))
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
	name := "tf-testAccClusterSENormalCreate"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":   acctest.CHECKSET,
		"db_node_class": acctest.CHECKSET,
		"vswitch_id":    acctest.CHECKSET,
		"db_type":       acctest.CHECKSET,
		"db_version":    acctest.CHECKSET,
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceClusterSENormalConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"resource_group_id": acctest.CHECKSET,
						"zone_id":           acctest.CHECKSET,
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
	name := "tf-testAccClusterSENormalEqualCreateWithStandbyAz"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":   acctest.CHECKSET,
		"db_node_class": acctest.CHECKSET,
		"vswitch_id":    acctest.CHECKSET,
		"db_type":       acctest.CHECKSET,
		"db_version":    acctest.CHECKSET,
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceClusterSENormalEqualConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":             "MySQL",
					"db_version":          "8.0",
					"pay_type":            "PostPaid",
					"db_node_class":       "polar.mysql.x4.large",
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
						"resource_group_id": acctest.CHECKSET,
						"zone_id":           acctest.CHECKSET,
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
	name := "tf-testAccClusterSENormalEqualCreate"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":   acctest.CHECKSET,
		"db_node_class": acctest.CHECKSET,
		"vswitch_id":    acctest.CHECKSET,
		"db_type":       acctest.CHECKSET,
		"db_version":    acctest.CHECKSET,
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceClusterSENormalEqualConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":             "MySQL",
					"db_version":          "8.0",
					"pay_type":            "PostPaid",
					"db_node_class":       "polar.mysql.x4.large",
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
						"resource_group_id": acctest.CHECKSET,
						"zone_id":           acctest.CHECKSET,
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
	name := "tf-testAccClusterServerless"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":   acctest.CHECKSET,
		"db_node_class": acctest.CHECKSET,
		"vswitch_id":    acctest.CHECKSET,
		"db_type":       acctest.CHECKSET,
		"db_version":    acctest.CHECKSET,
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceClusterServerlessConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"resource_group_id":                     acctest.CHECKSET,
						"zone_id":                               acctest.CHECKSET,
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
	name := "tf-testAccClusterSteadyServerless"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":   acctest.CHECKSET,
		"db_node_class": acctest.CHECKSET,
		"vswitch_id":    acctest.CHECKSET,
		"db_type":       acctest.CHECKSET,
		"db_version":    acctest.CHECKSET,
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceClusterSteadyServerlessConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"resource_group_id": acctest.CHECKSET,
						"zone_id":           acctest.CHECKSET,
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
						"imci_switch":   acctest.CHECKSET,
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
	name := "tf-testAccClusterCreateNormal"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       acctest.CHECKSET,
		"db_node_class":     acctest.CHECKSET,
		"vswitch_id":        acctest.CHECKSET,
		"db_type":           acctest.CHECKSET,
		"db_version":        acctest.CHECKSET,
		"connection_string": acctest.REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":                "MySQL",
					"db_version":             "8.0",
					"pay_type":               "PostPaid",
					"db_node_count":          "2",
					"db_node_class":          "polar.mysql.x4.large",
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
						"resource_group_id":      acctest.CHECKSET,
						"zone_id":                acctest.CHECKSET,
						"lower_case_table_names": "0",
						"default_time_zone":      acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudPolarDBCluster_EnableDynamoDB(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testAccClusterEnableDynamoDB"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       acctest.CHECKSET,
		"db_node_class":     acctest.CHECKSET,
		"db_type":           acctest.CHECKSET,
		"db_version":        acctest.CHECKSET,
		"connection_string": acctest.REGEXMATCH + clusterConnectionStringRegexp,
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":         "PostgreSQL",
					"db_version":      "14",
					"pay_type":        "PostPaid",
					"db_node_count":   "2",
					"db_node_class":   "polar.pg.x4.medium",
					"enable_dynamodb": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_dynamodb": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":         "PostgreSQL",
					"db_version":      "14",
					"pay_type":        "PostPaid",
					"db_node_count":   "2",
					"db_node_class":   "polar.pg.x4.medium",
					"enable_dynamodb": false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_dynamodb": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":         "PostgreSQL",
					"db_version":      "14",
					"pay_type":        "PostPaid",
					"db_node_count":   "2",
					"db_node_class":   "polar.pg.x4.medium",
					"enable_dynamodb": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_dynamodb": "true",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudPolarDBCluster_TargetMinorVersion(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testAccClusterTargetMinorVersion"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       acctest.CHECKSET,
		"db_node_class":     acctest.CHECKSET,
		"db_type":           acctest.CHECKSET,
		"db_version":        acctest.CHECKSET,
		"connection_string": acctest.REGEXMATCH + clusterConnectionStringRegexp,
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":              "PostgreSQL",
					"db_version":           "16",
					"pay_type":             "PostPaid",
					"db_node_class":        "polar.pg.x4.medium",
					"target_minor_version": "v16.13.17.7#20260611",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_minor_version": "v16.13.17.7#20260611",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudPolarDBCluster_Xengine(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testAccCluster-x-engine"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       acctest.CHECKSET,
		"db_node_class":     acctest.CHECKSET,
		"vswitch_id":        acctest.CHECKSET,
		"db_type":           acctest.CHECKSET,
		"db_version":        acctest.CHECKSET,
		"tde_status":        "Disabled",
		"connection_string": acctest.REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceClusterXengineConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":                      "MySQL",
					"db_version":                   "8.0",
					"pay_type":                     "PostPaid",
					"db_node_class":                "polar.mysql.x4.large",
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

func TestAccAliCloudPolarDBCluster_3AZ(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testAccCluster-3az"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       acctest.CHECKSET,
		"db_node_class":     acctest.CHECKSET,
		"vswitch_id":        acctest.CHECKSET,
		"db_type":           acctest.CHECKSET,
		"db_version":        acctest.CHECKSET,
		"tde_status":        "Disabled",
		"connection_string": acctest.REGEXMATCH + clusterConnectionStringRegexp,
		"port":              "3306",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCluster3ZAConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":             "MySQL",
					"db_version":          "8.0",
					"pay_type":            "PostPaid",
					"db_node_class":       "polar.mysql.x4.large",
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

func testAccCheckKeyValueInMapsForCluster(ps []map[string]interface{}, propName, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, policy := range ps {
			if policy[key].(string) != value {
				return fmt.Errorf("DB %s attribute '%s' expected %#v, got %#v", propName, key, value, policy[key])
			}
		}
		return nil
	}
}

func resourceCluster3ZAConfigDependence(name string) string {
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

func resourceClusterXengineConfigDependence(name string) string {
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

func resourceClusterPrePaidConfigDependence(name string) string {
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

func resourceClusterTDEConfigDependence(name string) string {
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
		security_group_name   = var.name
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
      role_name   = "${var.name}-ram-role"
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
	
	// principal_name = "<roleName>@role.<userDefaultDomainName>"
    resource "alicloud_resource_manager_policy_attachment" "default" {
	    policy_name       = "AliyunRDSInstanceEncryptionRolePolicy"
	    policy_type       = "System"
	    principal_name    = "${alicloud_ram_role.default.name}@role.tfexample.onaliyunservice.com"
	    principal_type    = "ServiceRole"
	    resource_group_id = "${data.alicloud_account.current.id}"
		depends_on = [alicloud_ram_role.default]
    }`, name)
}

// resourceClusterTDEPGConfigDependence is the PostgreSQL-compatible
// counterpart of resourceClusterTDEConfigDependence. The node class data
// source omits db_type on purpose: it matches the API engine string (e.g.
// "pg14") by substring, and "PostgreSQL" is not a substring of "pg14".
// AliyunRDSInstanceEncryptionDefaultRole is pre-provisioned in the shared test
// account with AliyunRDSInstanceEncryptionRolePolicy attached, so it is looked
// up instead of created (a fixed-name create would 409 EntityAlreadyExists.Role).
func resourceClusterTDEPGConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
	}

    data "alicloud_polardb_node_classes" "this" {
	  db_version = "14"
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

    resource "alicloud_kms_key" "default" {
        description            = var.name
        pending_window_in_days = 7
        status                 = "Enabled"
    }

    data "alicloud_ram_roles" "default" {
      name_regex = "^AliyunRDSInstanceEncryptionDefaultRole$"
    }`, name)
}

func resourceClusterConfigDependence(name string) string {
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

	resource "alicloud_polardb_global_security_ip_group" "default" {
		global_ip_list       = "192.168.0.1"
  		global_ip_group_name = "example_template"
	}

	resource "alicloud_polardb_global_security_ip_group" "default2" {
		global_ip_list       = "192.168.0.2"
  		global_ip_group_name = "example_template2"
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

func resourceCloneFromClusterConfigDependence(name string) string {
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

func resourceClusterNormalMultimasterConfigDependence(name string) string {
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

func resourceClusterSENormalConfigDependence(name string) string {
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

func resourceClusterSENormalEqualConfigDependence(name string) string {
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

func resourceClusterServerlessConfigDependence(name string) string {
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

func resourceClusterSteadyServerlessConfigDependence(name string) string {
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
	p := acctest.Provider.ResourcesMap
	ramPolicyExisted, _ := schema.InternalMap(p["alicloud_resource_manager_role"].Schema).Data(nil, nil)
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return err
	}
	accountId, err := rawClient.(*connectivity.AliyunClient).AccountId()
	policyId := fmt.Sprintf("AliyunRDSInstanceEncryptionRolePolicy:System:AliyunRDSInstanceEncryptionDefaultRole@role.%s.onaliyunservice.com:ServiceRole:%s", accountId, accountId)
	ramPolicyExisted.SetId(policyId)
	err = acctest.Provider.ResourcesMap["alicloud_resource_manager_policy_attachment"].Delete(ramPolicyExisted, rawClient)
	if err != nil && !errs.IsExpectedErrors(err, []string{"EntityNotExist.Role", "EntityNotExist.Policy"}) {
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

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Ram", "2015-05-01", action, query, request, true)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"DeleteConflict.Role.Policy"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.Role"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

func TestAccAliCloudPolarDBCluster_TDEOnCreateAndUpdateEncryptNewTables(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testAccClusterTDECreate"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":       acctest.CHECKSET,
		"db_node_class":     acctest.CHECKSET,
		"vswitch_id":        acctest.CHECKSET,
		"db_type":           acctest.CHECKSET,
		"db_version":        acctest.CHECKSET,
		"tde_status":        "Enabled",
		"connection_string": acctest.REGEXMATCH + clusterConnectionStringRegexp,
		"status":            acctest.CHECKSET,
		"create_time":       acctest.CHECKSET,
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceClusterKMSConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Shanghai))
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":       "MySQL",
					"db_version":    "8.0",
					"pay_type":      "PostPaid",
					"db_node_class": "polar.mysql.g1.tiny.c",
					"vswitch_id":    "${local.vswitch_id}",
					"vpc_id":        "${local.vpc_id}",
					"description":   "${var.name}",
					"tde_status":    "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tde_status": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypt_new_tables": "ON",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tde_status":         "Enabled",
						"encrypt_new_tables": "ON",
						"role_arn":           acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

// TestAccAliCloudPolarDBCluster_TDEAutomaticRotation_PG covers TDE automatic
// key rotation on a PostgreSQL-compatible PolarDB cluster. PostgreSQL is the
// engine covered by this fixture; ModifyDBClusterTDE also supports Oracle
// syntax-compatible clusters when a custom KMS key is supplied. Creation
// enables rotation, and the following steps toggle it on an already TDE-enabled
// cluster: ModifyDBClusterTDE rejects the whole request with
// InvalidTDEStatus.AlreadyEnabled when TDEStatus=Enable is sent for an enabled
// cluster, so these steps exercise the retry that drops TDEStatus and applies
// EnableAutomaticRotation.
func TestAccAliCloudPolarDBCluster_TDEAutomaticRotation_PG(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testAccClusterTDEAutoRotationPG"
	resourceId := "alicloud_polardb_cluster.default"
	var basicMap = map[string]string{
		"description":   acctest.CHECKSET,
		"db_node_class": acctest.CHECKSET,
		"db_type":       acctest.CHECKSET,
		"db_version":    acctest.CHECKSET,
		"tde_status":    "Enabled",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return polardbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeClusterAttribute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceClusterTDEPGConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":                   "PostgreSQL",
					"db_version":                "14",
					"pay_type":                  "PostPaid",
					"db_node_count":             "2",
					"db_node_class":             "polar.pg.x4.medium",
					"vswitch_id":                "${local.vswitch_id}",
					"description":               "${var.name}",
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"tde_status":                "Enabled",
					"encrypt_new_tables":        "OFF",
					"encryption_key":            "${alicloud_kms_key.default.id}",
					"role_arn":                  "${data.alicloud_ram_roles.default.roles.0.arn}",
					"enable_automatic_rotation": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encryption_key":     acctest.CHECKSET,
						"role_arn":           acctest.CHECKSET,
						"automatic_rotation": "Enabled",
					}),
					resource.TestCheckResourceAttrSet(resourceId, "rotation_interval"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_automatic_rotation": false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_automatic_rotation": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "Enabled",
					}),
					resource.TestCheckResourceAttrSet(resourceId, "rotation_interval"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_type", "imci_switch", "sub_category", "enable_automatic_rotation"},
			},
		},
	})
}

func resourceClusterKMSConfigDependence(name string) string {
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
		zone_id = "cn-shanghai-l"
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

`, name)
}

var endpoints sync.Map
var signVersion sync.Map
var regionalEndpoints sync.Map

// sharedEndpointCacheForRegion keeps SDK endpoint discovery per region; one
// shared cache would route a Hangzhou request to a cached Beijing endpoint.
func sharedEndpointCacheForRegion(region string) *sync.Map {
	cache, _ := regionalEndpoints.LoadOrStore(region, new(sync.Map))
	return cache.(*sync.Map)
}

// sharedClientForRegion returns a common AlicloudClient setup needed for the
// sweeper functions for a given region.
func sharedClientForRegion(region string) (interface{}, error) {
	var accessKey, secretKey, securityToken string
	if accessKey = os.Getenv("ALICLOUD_ACCESS_KEY"); accessKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_ACCESS_KEY")
	}

	if secretKey = os.Getenv("ALICLOUD_SECRET_KEY"); secretKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_SECRET_KEY")
	}

	securityToken = os.Getenv("ALICLOUD_SECURITY_TOKEN")

	conf := connectivity.Config{
		Region:      connectivity.Region(region),
		RegionId:    region,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		Protocol:    "HTTPS",
		Endpoints:   sharedEndpointCacheForRegion(region),
		SignVersion: &signVersion,
	}
	if securityToken != "" {
		conf.SecurityToken = securityToken
	}
	accountId := os.Getenv("ALICLOUD_ACCOUNT_ID")
	if accountId == "" {
		accountId = os.Getenv("ALIBABA_CLOUD_ACCOUNT_ID")
	}
	if accountId != "" {
		conf.AccountId = accountId
	}
	accountType := os.Getenv("ALICLOUD_ACCOUNT_TYPE")
	if accountType == "" {
		accountType = os.Getenv("ALIBABA_CLOUD_ACCOUNT_TYPE")
	}
	if accountType != "" {
		conf.AccountType = accountType
	}
	credentialConfig := new(credentials.Config).SetType("access_key").SetAccessKeyId(accessKey).SetAccessKeySecret(secretKey)
	if v := strings.TrimSpace(securityToken); v != "" {
		credentialConfig.SetType("sts").SetSecurityToken(v)
	}
	credential, err := credentials.NewCredential(credentialConfig)
	if err != nil {
		return nil, err
	}
	conf.Credential = credential

	client, err := conf.Client()
	if err != nil {
		return nil, err
	}

	return client, nil
}

const PolarDBCommonTestCase = `
data "alicloud_polardb_zones" "default"{}
data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}
resource "alicloud_vpc" "default" {
    vpc_name = var.name
}
resource "alicloud_vswitch" "default" {
 vswitch_name = "tf_testAccPolarDB"
 vpc_id  = alicloud_vpc.default.id
 zone_id = data.alicloud_polardb_zones.default.ids[length(data.alicloud_polardb_zones.default.ids)-1]
 cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
}
locals {
  vpc_id = alicloud_vpc.default.id
  vswitch_id = concat(alicloud_vswitch.default.*.id, [""])[0]
  zone_id = data.alicloud_polardb_zones.default.ids[length(data.alicloud_polardb_zones.default.ids)-1]
}
`

func TestUnitPolarDBDistributedNodeSchema(t *testing.T) {
	clusterSchema := polardbsvc.ResourceCluster().Schema
	for _, key := range []string{"cn_node_class", "dn_node_class", "cn_node_num", "dn_node_num"} {
		field := clusterSchema[key]
		if !field.Optional || !field.Computed || field.ForceNew {
			t.Fatalf("%s must be Optional and Computed without ForceNew", key)
		}
	}
	for _, key := range []string{"cn_node_num", "dn_node_num"} {
		if clusterSchema[key].Type != schema.TypeInt {
			t.Fatalf("%s must use TypeInt", key)
		}
	}
	if len(clusterSchema["db_node_count"].ConflictsWith) != 2 {
		t.Fatal("db_node_count must conflict with distributed node count fields")
	}
	for _, key := range []string{"modify_dist_nodes", "delete_sub_groups", "restart_cluster"} {
		if _, ok := clusterSchema[key]; ok {
			t.Fatalf("command-style field %s must not be part of desired state", key)
		}
	}
}

func TestUnitClassifyPolarDBDistributedNode(t *testing.T) {
	tests := []struct {
		name     string
		node     map[string]interface{}
		wantCN   []string
		wantDN   []string
		retained bool
		noClass  bool
	}{
		{name: "compute type", node: map[string]interface{}{"DBNodeId": "cn-1", "DBNodeType": "Compute"}, wantCN: []string{"cn-1"}},
		{name: "data type", node: map[string]interface{}{"DBNodeId": "dn-1", "DBNodeType": "Data"}, wantDN: []string{"dn-1"}},
		{name: "compute subgroup", node: map[string]interface{}{"DBNodeId": "cn-0", "SubGroupType": "Compute"}, wantCN: []string{"cn-0"}},
		{name: "data subgroup", node: map[string]interface{}{"DBNodeId": "dn-0", "SubGroupType": "Data"}, wantDN: []string{"dn-0"}},
		{name: "CN role", node: map[string]interface{}{"DBNodeId": "cn-2", "DBNodeRole": "CN"}, wantCN: []string{"cn-2"}},
		{name: "DN role", node: map[string]interface{}{"DBNodeId": "dn-2", "DBNodeRole": "DN"}, wantDN: []string{"dn-2"}},
		{name: "legacy CN", node: map[string]interface{}{"DBNodeId": "cn-3", "DBNodeClass": "polar.pg.x4.medium", "IsPrimaryCN": false}, wantCN: []string{"cn-3"}},
		{name: "legacy CN reader", node: map[string]interface{}{"DBNodeId": "cn-4", "DBNodeClass": "polar.pg.x4.medium", "DBNodeRole": "Reader", "IsPrimaryCN": false}, wantCN: []string{"cn-4"}},
		{name: "legacy DN", node: map[string]interface{}{"DBNodeId": "dn-3", "DBNodeClass": "polar.pg.x4.medium"}, wantDN: []string{"dn-3"}},
		{name: "unknown role", node: map[string]interface{}{"DBNodeId": "node-1", "DBNodeClass": "polar.pg.x4.medium", "DBNodeRole": "Reader"}, retained: true},
		{name: "missing role", node: map[string]interface{}{"DBNodeId": "node-2"}, retained: true, noClass: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, ok := tt.node["DBNodeClass"]; !ok && !tt.noClass {
				tt.node["DBNodeClass"] = "polar.pg.x4.medium"
			}
			tt.node["SubGroupName"] = "sg-" + tt.name
			api := &clusterUnitTestAPI{t: t, nodes: []map[string]interface{}{tt.node}}
			client := newClusterUnitTestClient(t, api)
			d := newClusterUnitTestResource(t, map[string]interface{}{
				"cn_node_class": "polar.pg.x4.large",
				"dn_node_class": "polar.pg.x4.large",
				"cn_node_ids":   []string{"prior-cn"},
				"dn_node_ids":   []string{"prior-dn"},
			}, nil)
			if err := polardbsvc.ResourceCluster().Read(d, client); err != nil {
				t.Fatalf("unexpected read error: %s", err)
			}
			if tt.retained {
				if got := clusterUnitTestSetValues(t, d, "cn_node_ids"); strings.Join(got, ",") != "prior-cn" {
					t.Fatalf("unclassifiable node must retain the prior CN ids, got %#v", got)
				}
				if got := d.Get("cn_node_class").(string); got != "polar.pg.x4.large" {
					t.Fatalf("unclassifiable node must retain the prior CN class, got %q", got)
				}
				return
			}
			if got := clusterUnitTestSetValues(t, d, "cn_node_ids"); strings.Join(got, ",") != strings.Join(tt.wantCN, ",") {
				t.Fatalf("expected CN ids %#v, got %#v", tt.wantCN, got)
			}
			if got := clusterUnitTestSetValues(t, d, "dn_node_ids"); strings.Join(got, ",") != strings.Join(tt.wantDN, ",") {
				t.Fatalf("expected DN ids %#v, got %#v", tt.wantDN, got)
			}
			if len(tt.wantCN) > 0 && d.Get("cn_node_class").(string) != "polar.pg.x4.medium" {
				t.Fatalf("expected the CN class to be flattened from the API response, got %q", d.Get("cn_node_class").(string))
			}
			if len(tt.wantDN) > 0 && d.Get("dn_node_class").(string) != "polar.pg.x4.medium" {
				t.Fatalf("expected the DN class to be flattened from the API response, got %q", d.Get("dn_node_class").(string))
			}
		})
	}
}

func TestUnitValidatePolarDBDistributedNodeCounts(t *testing.T) {
	newNode := func(id, subGroup string) map[string]interface{} {
		return map[string]interface{}{
			"DBNodeId": id, "SubGroupName": subGroup, "SubGroupType": "Compute",
			"DBNodeClass": "polar.pg.x4.large", "CreationTime": "2026-01-01T00:00:00Z",
		}
	}
	t.Run("matching subgroup counts", func(t *testing.T) {
		api := &clusterUnitTestAPI{t: t, cnNodeCount: 2, dnNodeCount: 1, nodes: []map[string]interface{}{
			newNode("cn-1", "sg-cn-1"),
			newNode("cn-2", "sg-cn-2"),
			func() map[string]interface{} {
				node := newNode("dn-1", "sg-dn-1")
				node["SubGroupType"] = "Data"
				return node
			}(),
		}}
		client := newClusterUnitTestClient(t, api)
		d := newClusterUnitTestResource(t, map[string]interface{}{
			"cn_node_class": "polar.pg.x4.medium",
			"dn_node_class": "polar.pg.x4.medium",
			"cn_node_ids":   []string{"prior-cn"},
			"dn_node_ids":   []string{"prior-dn"},
		}, nil)
		if err := polardbsvc.ResourceCluster().Read(d, client); err != nil {
			t.Fatalf("unexpected read error: %s", err)
		}
		if got := d.Get("cn_node_class").(string); got != "polar.pg.x4.large" {
			t.Fatalf("expected subgroup counts that match the logical counts to be flattened, got CN class %q", got)
		}
		if got := clusterUnitTestSetValues(t, d, "cn_node_ids"); strings.Join(got, ",") != "cn-1,cn-2" {
			t.Fatalf("expected both CN nodes, got %#v", got)
		}
		if got := clusterUnitTestSetValues(t, d, "dn_node_ids"); strings.Join(got, ",") != "dn-1" {
			t.Fatalf("expected the DN node, got %#v", got)
		}
	})
	t.Run("misclassified nodes are rejected", func(t *testing.T) {
		api := &clusterUnitTestAPI{t: t, cnNodeCount: 3, dnNodeCount: 0, nodes: []map[string]interface{}{
			newNode("cn-1", "sg-cn-1"),
			newNode("cn-2", "sg-cn-2"),
			func() map[string]interface{} {
				node := newNode("dn-1", "sg-dn-1")
				node["SubGroupType"] = "Data"
				return node
			}(),
		}}
		client := newClusterUnitTestClient(t, api)
		d := newClusterUnitTestResource(t, map[string]interface{}{
			"cn_node_class": "polar.pg.x4.medium",
			"dn_node_class": "polar.pg.x4.medium",
			"cn_node_ids":   []string{"prior-cn"},
			"dn_node_ids":   []string{"prior-dn"},
		}, nil)
		if err := polardbsvc.ResourceCluster().Read(d, client); err != nil {
			t.Fatalf("unexpected read error: %s", err)
		}
		if got := d.Get("cn_node_class").(string); got != "polar.pg.x4.medium" {
			t.Fatalf("expected silently misclassified nodes to keep the prior CN class, got %q", got)
		}
		if got := clusterUnitTestSetValues(t, d, "cn_node_ids"); strings.Join(got, ",") != "prior-cn" {
			t.Fatalf("expected the prior CN ids to be retained, got %#v", got)
		}
		if got := d.Get("cn_node_num").(int); got != 3 {
			t.Fatalf("expected the logical CN count to still be refreshed, got %d", got)
		}
	})
}

func TestUnitSelectPolarDBDistributedSubGroupsForDeletion(t *testing.T) {
	nodes := []map[string]interface{}{
		{"DBNodeId": "pi-primary-writer", "SubGroupName": clusterUnitTestID, "SubGroupType": "Compute", "DBNodeClass": "polar.pg.x4.large", "CreationTime": "2026-01-01T00:00:00Z"},
		{"DBNodeId": "pi-primary-reader", "SubGroupName": clusterUnitTestID, "SubGroupType": "Compute", "DBNodeClass": "polar.pg.x4.large", "CreationTime": "2026-01-01T00:00:00Z"},
		{"DBNodeId": "pi-extra-writer", "SubGroupName": "extra-cn", "SubGroupType": "Compute", "DBNodeClass": "polar.pg.x4.large", "CreationTime": "2026-01-02T00:00:00Z"},
		{"DBNodeId": "pi-extra-reader", "SubGroupName": "extra-cn", "SubGroupType": "Compute", "DBNodeClass": "polar.pg.x4.large", "CreationTime": "2026-01-02T00:00:00Z"},
		{"DBNodeId": "pi-newest", "SubGroupName": "newest-cn", "SubGroupType": "Compute", "DBNodeClass": "polar.pg.x4.large", "CreationTime": "2026-01-03T00:00:00Z"},
		{"DBNodeId": "pi-dn-primary", "SubGroupName": clusterUnitTestID, "SubGroupType": "Data", "DBNodeClass": "polar.pg.x4.large", "CreationTime": "2026-01-01T00:00:00Z"},
	}
	for _, tt := range []struct {
		name     string
		target   int
		expected []string
	}{
		{name: "delete the newest extra logical node", target: 2, expected: []string{"newest-cn"}},
		{name: "delete both extra logical nodes", target: 1, expected: []string{"newest-cn", "extra-cn"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			api := &clusterUnitTestAPI{t: t, nodes: nodes, cnNodeCount: 3, dnNodeCount: 1, deleteNodesError: true}
			client := newClusterUnitTestClient(t, api)
			d := newClusterUnitTestResource(t, map[string]interface{}{
				"db_type":           "PostgreSQL",
				"pay_type":          "Postpaid",
				"creation_category": "Normal",
				"cn_node_class":     "polar.pg.x4.large",
				"dn_node_class":     "polar.pg.x4.large",
				"cn_node_num":       3,
				"dn_node_num":       1,
			}, map[string]interface{}{"cn_node_num": tt.target})
			err := polardbsvc.ResourceCluster().Update(d, client)
			if err == nil {
				t.Fatal("expected the rejected DeleteDBNodes call to surface as an update error")
			}
			requests := api.requestsFor("DeleteDBNodes")
			if len(requests) != 1 {
				t.Fatalf("expected exactly one DeleteDBNodes call, got %d", len(requests))
			}
			if got := requests[0].Get("DBNodeType"); got != "Compute" {
				t.Fatalf("expected a Compute node deletion, got %q", got)
			}
			var deleted []string
			for i := 1; ; i++ {
				value := requests[0].Get(fmt.Sprintf("DBNodeId.%d", i))
				if value == "" {
					break
				}
				deleted = append(deleted, value)
			}
			if strings.Join(deleted, ",") != strings.Join(tt.expected, ",") {
				t.Fatalf("expected the deletion to target %#v, got %#v", tt.expected, deleted)
			}
		})
	}
}

func TestUnitPolarDBParameterModificationPendingStates(t *testing.T) {
	api := &clusterUnitTestAPI{t: t}
	client := newClusterUnitTestClient(t, api)
	d := newClusterUnitTestResource(t, map[string]interface{}{
		"db_type":           "PostgreSQL",
		"pay_type":          "Postpaid",
		"creation_category": "Normal",
	}, map[string]interface{}{
		"parameters": []interface{}{map[string]interface{}{"name": "wal_level", "value": "logical"}},
	})
	if err := polardbsvc.ResourceCluster().Update(d, client); err != nil {
		t.Fatalf("expected Rebooting to be treated as a pending state: %s", err)
	}
	if got := api.rebootingCount(); got < 2 {
		t.Fatalf("expected the waiter to observe Rebooting twice before Running, observed %d", got)
	}
	requests := api.requestsFor("ModifyDBClusterParameters")
	if len(requests) != 1 {
		t.Fatalf("expected exactly one ModifyDBClusterParameters call, got %d", len(requests))
	}
	if got := requests[0].Get("Parameters"); !strings.Contains(got, `"wal_level":"logical"`) {
		t.Fatalf("expected the changed parameter to be submitted, got %q", got)
	}
}

func TestUnitPolarDBClusterTDEAutomaticRotationRequest(t *testing.T) {
	rotation := func(value bool) *string {
		text := fmt.Sprintf("%t", value)
		return &text
	}
	baseAttributes := func() map[string]interface{} {
		return map[string]interface{}{
			"db_type":           "PostgreSQL",
			"pay_type":          "Postpaid",
			"creation_category": "Normal",
			"tde_status":        "Enabled",
		}
	}
	alreadyEnabled := clusterUnitTestTdeResponse{code: "InvalidTDEStatus.AlreadyEnabled", message: "TDE is already enabled"}
	rejected := clusterUnitTestTdeResponse{code: "InvalidParameter", message: "unit test rejects the TDE change"}

	tests := []struct {
		name         string
		attributes   map[string]interface{}
		changes      map[string]interface{}
		rawConfig    map[string]cty.Value
		responses    []clusterUnitTestTdeResponse
		wantRequests int
		wantRotation *string
		wantParams   []map[string]string
	}{
		{
			name:         "rotation not configured",
			changes:      map[string]interface{}{"encrypt_new_tables": "ON"},
			rawConfig:    map[string]cty.Value{},
			responses:    []clusterUnitTestTdeResponse{rejected},
			wantRequests: 1,
			wantParams:   []map[string]string{{"TDEStatus": "Enable", "EncryptNewTables": "ON"}},
		},
		{
			name:         "explicit true",
			changes:      map[string]interface{}{"enable_automatic_rotation": true},
			rawConfig:    map[string]cty.Value{"enable_automatic_rotation": cty.True},
			responses:    []clusterUnitTestTdeResponse{rejected},
			wantRequests: 1,
			wantRotation: rotation(true),
			wantParams:   []map[string]string{{"TDEStatus": "Enable"}},
		},
		{
			name:         "explicit false",
			attributes:   map[string]interface{}{"enable_automatic_rotation": true},
			changes:      map[string]interface{}{"enable_automatic_rotation": false},
			rawConfig:    map[string]cty.Value{"enable_automatic_rotation": cty.False},
			responses:    []clusterUnitTestTdeResponse{rejected},
			wantRequests: 1,
			wantRotation: rotation(false),
			wantParams:   []map[string]string{{"TDEStatus": "Enable"}},
		},
		{
			name:         "removed from configuration",
			attributes:   map[string]interface{}{"enable_automatic_rotation": true},
			changes:      map[string]interface{}{"enable_automatic_rotation": false},
			rawConfig:    map[string]cty.Value{},
			responses:    []clusterUnitTestTdeResponse{rejected},
			wantRequests: 1,
			wantRotation: rotation(false),
			wantParams:   []map[string]string{{"TDEStatus": "Enable"}},
		},
		{
			name:       "other TDE update keeps the explicit value",
			attributes: map[string]interface{}{"enable_automatic_rotation": true},
			changes: map[string]interface{}{
				"encryption_key": "new-key",
				"role_arn":       "acs:ram::1234567890:role/new",
			},
			rawConfig: map[string]cty.Value{
				"enable_automatic_rotation": cty.True,
				"encryption_key":            cty.StringVal("new-key"),
				"role_arn":                  cty.StringVal("acs:ram::1234567890:role/new"),
			},
			responses:    []clusterUnitTestTdeResponse{rejected},
			wantRequests: 1,
			wantRotation: rotation(true),
			wantParams: []map[string]string{{
				"TDEStatus":     "Enable",
				"EncryptionKey": "new-key",
				"RoleArn":       "acs:ram::1234567890:role/new",
			}},
		},
		{
			name:         "fallback keeps the explicit value",
			changes:      map[string]interface{}{"encrypt_new_tables": "ON"},
			rawConfig:    map[string]cty.Value{"enable_automatic_rotation": cty.True},
			responses:    []clusterUnitTestTdeResponse{alreadyEnabled, rejected},
			wantRequests: 2,
			wantRotation: rotation(true),
			wantParams: []map[string]string{
				{"TDEStatus": "Enable", "EncryptNewTables": "ON"},
				{"EncryptNewTables": "ON"},
			},
		},
		{
			name:         "fallback does not carry the unconfigured switch",
			changes:      map[string]interface{}{"encrypt_new_tables": "ON"},
			rawConfig:    map[string]cty.Value{},
			responses:    []clusterUnitTestTdeResponse{alreadyEnabled, rejected},
			wantRequests: 2,
			wantParams: []map[string]string{
				{"TDEStatus": "Enable", "EncryptNewTables": "ON"},
				{"EncryptNewTables": "ON"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &clusterUnitTestAPI{t: t, tdeResponses: tt.responses}
			client := newClusterUnitTestClient(t, api)
			attributes := baseAttributes()
			for key, value := range tt.attributes {
				attributes[key] = value
			}
			d := newClusterUnitTestResourceWithRawConfig(t, attributes, tt.changes, tt.rawConfig)
			if err := polardbsvc.ResourceCluster().Update(d, client); err == nil {
				t.Fatal("expected the scripted ModifyDBClusterTDE failure to abort the update")
			}
			requests := api.requestsFor("ModifyDBClusterTDE")
			if len(requests) != tt.wantRequests {
				t.Fatalf("expected %d ModifyDBClusterTDE calls, got %d: %v", tt.wantRequests, len(requests), requests)
			}
			for index, values := range requests {
				t.Logf("ModifyDBClusterTDE wire request %d: %s", index+1, describeTdeWireRequest(values))
				if tt.wantRotation == nil {
					if _, ok := values["EnableAutomaticRotation"]; ok {
						t.Fatalf("request %d must not carry EnableAutomaticRotation when it is not configured: %v", index+1, values)
					}
				} else if got := values.Get("EnableAutomaticRotation"); got != *tt.wantRotation {
					t.Fatalf("request %d: expected EnableAutomaticRotation=%s, got %q", index+1, *tt.wantRotation, got)
				}
				for key, want := range tt.wantParams[index] {
					if got := values.Get(key); got != want {
						t.Fatalf("request %d: expected %s=%q, got %q", index+1, key, want, got)
					}
				}
				for _, key := range []string{"TDEStatus", "EncryptNewTables", "EncryptionKey", "RoleArn"} {
					if _, expected := tt.wantParams[index][key]; expected {
						continue
					}
					if _, ok := values[key]; ok {
						t.Fatalf("request %d must not carry %s: %v", index+1, key, values)
					}
				}
			}
		})
	}
}

// describeTdeWireRequest renders the TDE parameters of one ModifyDBClusterTDE
// call as they were sent over the wire, so the request sample can be quoted
// without exposing the signed/credential parameters of the form.
func describeTdeWireRequest(values url.Values) string {
	keys := []string{"Action", "DBClusterId", "TDEStatus", "EncryptNewTables", "EncryptionKey", "RoleArn", "EnableAutomaticRotation"}
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		if value, ok := values[key]; ok {
			parts = append(parts, fmt.Sprintf("%s=%s", key, strings.Join(value, ",")))
			continue
		}
		parts = append(parts, fmt.Sprintf("%s=<omitted>", key))
	}
	return strings.Join(parts, " ")
}

// TestUnitPolarDBClusterRefundPrePaidRequest exercises the billing calls the
// delete path makes for a subscription cluster: the commodity ProductType is
// read from the account's billing record and handed to RefundInstance with the
// immediate release flag.
func TestUnitPolarDBClusterRefundPrePaidRequest(t *testing.T) {
	tests := []struct {
		name         string
		api          *clusterUnitTestAPI
		wantOrderId  int64
		wantNotFound bool
		wantError    string
		wantRejected bool
	}{
		{
			name:        "refund reuses the billing product type",
			api:         &clusterUnitTestAPI{bssProductType: "polardb_subscription_public_cn"},
			wantOrderId: 2100000000001,
		},
		{
			name:         "missing billing record is not found",
			api:          &clusterUnitTestAPI{bssInstanceMiss: true},
			wantNotFound: true,
		},
		{
			name:      "refund failure is propagated",
			api:       &clusterUnitTestAPI{refundErrorCode: "ExistRefundingOrderError"},
			wantError: "ExistRefundingOrderError",
		},
		{
			name:      "business rejection keeps the response code",
			api:       &clusterUnitTestAPI{refundFalseCode: "ResourceNotExists"},
			wantError: "ResourceNotExists",
		},
		{
			name:         "missing success flag is rejected",
			api:          &clusterUnitTestAPI{refundNoSuccess: true},
			wantRejected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.api.t = t
			client := newClusterUnitTestClient(t, tt.api)
			orderId, err := polardbsvc.NewService(client).RefundPrePaidCluster(clusterUnitTestID, 30*time.Second)

			switch {
			case tt.wantNotFound:
				if err == nil || !errs.NotFoundError(err) {
					t.Fatalf("expected a not found error, got %v", err)
				}
				if calls := tt.api.requestsFor("RefundInstance"); len(calls) != 0 {
					t.Fatalf("a missing billing record must not be refunded: %v", calls)
				}
				return
			case tt.wantRejected:
				if err == nil {
					t.Fatal("expected a response without Success to be rejected")
				}
				if errs.IsExpectedErrors(err, []string{"ResourceNotExists", "ExistRefundingOrderError"}) {
					t.Fatalf("a response without Success must not look like an idempotent rejection: %v", err)
				}
				return
			case tt.wantError != "":
				if err == nil {
					t.Fatal("expected the refund failure to be returned")
				}
				if !errs.IsExpectedErrors(err, []string{tt.wantError}) {
					t.Fatalf("expected error code %s, got %v", tt.wantError, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected refund error: %v", err)
			}
			if orderId != tt.wantOrderId {
				t.Fatalf("expected refund order id %d, got %d", tt.wantOrderId, orderId)
			}
			queries := tt.api.requestsFor("QueryAvailableInstances")
			if len(queries) != 1 {
				t.Fatalf("expected one billing lookup, got %d", len(queries))
			}
			if got := queries[0].Get("ProductCode"); got != "polardb" {
				t.Fatalf("expected ProductCode=polardb, got %q", got)
			}
			if got := queries[0].Get("InstanceIDs"); got != clusterUnitTestID {
				t.Fatalf("expected InstanceIDs=%s, got %q", clusterUnitTestID, got)
			}
			refunds := tt.api.requestsFor("RefundInstance")
			if len(refunds) != 1 {
				t.Fatalf("expected one refund call, got %d", len(refunds))
			}
			refund := refunds[0]
			if got := refund.Get("ProductCode"); got != "polardb" {
				t.Fatalf("expected ProductCode=polardb, got %q", got)
			}
			if got := refund.Get("ProductType"); got != "polardb_subscription_public_cn" {
				t.Fatalf("expected the billing ProductType to be reused, got %q", got)
			}
			if got := refund.Get("InstanceId"); got != clusterUnitTestID {
				t.Fatalf("expected InstanceId=%s, got %q", clusterUnitTestID, got)
			}
			if got := refund.Get("ImmediatelyRelease"); got != "1" {
				t.Fatalf("expected ImmediatelyRelease=1, got %q", got)
			}
			if got := refund.Get("ClientToken"); got == "" {
				t.Fatal("expected a ClientToken on the refund call")
			}
		})
	}
}

const clusterUnitTestID = "pc-unit-test"

// clusterUnitTestTdeResponse scripts one ModifyDBClusterTDE reply so the TDE
// update flow can be observed without reaching the cloud or its 3 minute wait.
type clusterUnitTestTdeResponse struct {
	code    string
	message string
}

// clusterUnitTestAPI serves the hermetic PolarDB responses used by the resource
// unit tests so the resource factory can be exercised without real cloud calls.
type clusterUnitTestAPI struct {
	t                *testing.T
	mu               sync.Mutex
	nodes            []map[string]interface{}
	cnNodeCount      interface{}
	dnNodeCount      interface{}
	statuses         []string
	requests         map[string][]url.Values
	deleteNodesError bool
	rebootingServed  int
	tdeStatus        string
	tdeResponses     []clusterUnitTestTdeResponse
	bssInstanceMiss  bool
	bssProductType   string
	refundErrorCode  string
	refundFalseCode  string
	refundNoSuccess  bool
}

func (api *clusterUnitTestAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		api.t.Errorf("failed to parse request form: %s", err)
	}
	action := r.Form.Get("Action")
	api.mu.Lock()
	if api.requests == nil {
		api.requests = make(map[string][]url.Values)
	}
	api.requests[action] = append(api.requests[action], r.Form)
	api.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")

	var response interface{}
	switch action {
	case "DescribeDBClusterAttribute":
		attribute := map[string]interface{}{
			"DBClusterId":       clusterUnitTestID,
			"DBClusterStatus":   api.nextStatus(),
			"Category":          "Normal",
			"DBType":            "PostgreSQL",
			"DBVersion":         "14",
			"PayType":           "Postpaid",
			"HotStandbyCluster": "StandbyClusterOFF",
			"VSwitchId":         "vsw-unit-test",
			"VPCId":             "vpc-unit-test",
			"ZoneIds":           "cn-hangzhou-h",
			"DBNodes":           api.nodes,
		}
		if api.cnNodeCount != nil {
			attribute["CnNodeCount"] = api.cnNodeCount
		}
		if api.dnNodeCount != nil {
			attribute["DnNodeCount"] = api.dnNodeCount
		}
		response = attribute
	case "DescribeDBClusters":
		response = map[string]interface{}{
			"Items": map[string]interface{}{
				"DBCluster": []interface{}{
					map[string]interface{}{"DBClusterId": clusterUnitTestID, "DBNodeClass": "polar.pg.x4.medium", "PayType": "Postpaid"},
				},
			},
		}
	case "DescribeDBClusterAccessWhitelist":
		response = map[string]interface{}{"Items": map[string]interface{}{"DBClusterIPArray": []interface{}{}}}
	case "DescribeDBClusterEndpoints":
		response = map[string]interface{}{
			"Items": []interface{}{
				map[string]interface{}{
					"EndpointType": "Cluster",
					"AddressItems": []interface{}{
						map[string]interface{}{"NetType": "Private", "ConnectionString": clusterUnitTestID + ".rwlb.rds.aliyuncs.com", "Port": "1921"},
					},
				},
			},
		}
	case "ListTagResources":
		response = map[string]interface{}{"TagResources": map[string]interface{}{"TagResource": []interface{}{}}}
	case "DescribeDBClusterParameters":
		response = map[string]interface{}{
			"RunningParameters": map[string]interface{}{
				"Parameter": []interface{}{
					map[string]interface{}{"ParameterName": "wal_level", "ParameterValue": "logical"},
				},
			},
		}
	case "DescribeDBClusterAuditLogCollector":
		response = map[string]interface{}{"CollectorStatus": "Disabled"}
	case "DescribeDBClusterTDE":
		api.mu.Lock()
		tdeStatus := api.tdeStatus
		api.mu.Unlock()
		if tdeStatus == "" {
			tdeStatus = "Disabled"
		}
		response = map[string]interface{}{"TDEStatus": tdeStatus}
	case "ModifyDBClusterTDE":
		api.mu.Lock()
		var next clusterUnitTestTdeResponse
		if len(api.tdeResponses) > 0 {
			next = api.tdeResponses[0]
			api.tdeResponses = api.tdeResponses[1:]
		}
		api.mu.Unlock()
		if next.code != "" {
			w.WriteHeader(http.StatusBadRequest)
			response = map[string]interface{}{"Code": next.code, "Message": next.message}
			break
		}
		response = map[string]interface{}{}
	case "DescribeDBClusterVersion":
		response = map[string]interface{}{"DBMinorVersion": ""}
	case "DescribeGlobalSecurityIPGroupRelation":
		response = map[string]interface{}{"GlobalSecurityIPGroupRel": []interface{}{}}
	case "QueryAvailableInstances":
		if api.bssInstanceMiss {
			response = map[string]interface{}{"Success": true, "Data": map[string]interface{}{"InstanceList": []interface{}{}}}
			break
		}
		productType := api.bssProductType
		if productType == "" {
			productType = "online"
		}
		response = map[string]interface{}{
			"Success": true,
			"Data": map[string]interface{}{
				"InstanceList": []interface{}{
					map[string]interface{}{
						"InstanceID":  clusterUnitTestID,
						"ProductCode": "polardb",
						"ProductType": productType,
					},
				},
			},
		}
	case "RefundInstance":
		if api.refundErrorCode != "" {
			w.WriteHeader(http.StatusBadRequest)
			response = map[string]interface{}{"Code": api.refundErrorCode, "Message": "unit test rejects the refund"}
			break
		}
		if api.refundFalseCode != "" {
			// HTTP 200 with Success=false and an empty message is the shape that
			// carries the business code without an SDK error.
			response = map[string]interface{}{"Success": false, "Code": api.refundFalseCode, "Message": ""}
			break
		}
		if api.refundNoSuccess {
			// A response without a boolean Success must not be read as an
			// accepted refund.
			response = map[string]interface{}{"Data": map[string]interface{}{"OrderId": 0}}
			break
		}
		response = map[string]interface{}{"Success": true, "Data": map[string]interface{}{"OrderId": 2100000000001}}
	case "ModifyDBClusterParameters":
		api.mu.Lock()
		api.statuses = append(api.statuses, "Rebooting", "Rebooting", "Running")
		api.mu.Unlock()
		response = map[string]interface{}{}
	case "DeleteDBNodes":
		if !api.deleteNodesError {
			api.mu.Lock()
			api.statuses = append(api.statuses, "DBNodeDeleting", "Running")
			api.mu.Unlock()
			response = map[string]interface{}{}
			break
		}
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]interface{}{"Code": "InvalidParameter", "Message": "unit test rejects node deletion"}
	default:
		api.t.Errorf("unexpected PolarDB action %q", action)
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]interface{}{"Code": "UnexpectedAction", "Message": action}
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		api.t.Errorf("failed to encode response: %s", err)
	}
}

func (api *clusterUnitTestAPI) nextStatus() string {
	api.mu.Lock()
	defer api.mu.Unlock()
	if len(api.statuses) == 0 {
		return "Running"
	}
	status := api.statuses[0]
	api.statuses = api.statuses[1:]
	if status == "Rebooting" {
		api.rebootingServed++
	}
	return status
}

func (api *clusterUnitTestAPI) requestsFor(action string) []url.Values {
	api.mu.Lock()
	defer api.mu.Unlock()
	return append([]url.Values(nil), api.requests[action]...)
}

func (api *clusterUnitTestAPI) rebootingCount() int {
	api.mu.Lock()
	defer api.mu.Unlock()
	return api.rebootingServed
}

func newClusterUnitTestClient(t *testing.T, api *clusterUnitTestAPI) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(api)
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("polardb", endpoint)
	// The subscription release path calls the billing API, which the same
	// hermetic server also serves.
	endpoints.Store("bssopenapi", endpoint)
	return client
}

func newClusterUnitTestResource(t *testing.T, attributes, changes map[string]interface{}) *schema.ResourceData {
	t.Helper()
	return newClusterUnitTestResourceWithRawConfig(t, attributes, changes, nil)
}

// newClusterUnitTestResourceWithRawConfig builds the same merged update state as
// newClusterUnitTestResource and additionally attaches a raw configuration so
// the GetRawConfig based presence checks can be exercised. A nil map means the
// raw configuration is unknown, like Update calls triggered from Create.
func newClusterUnitTestResourceWithRawConfig(t *testing.T, attributes, changes map[string]interface{}, configured map[string]cty.Value) *schema.ResourceData {
	t.Helper()
	resourceSchema := polardbsvc.ResourceCluster().Schema
	stateData, err := schema.InternalMap(resourceSchema).Data(&sdkterraform.InstanceState{
		ID:         clusterUnitTestID,
		Attributes: map[string]string{"id": clusterUnitTestID},
	}, nil)
	if err != nil {
		t.Fatalf("failed to build the unit test state: %s", err)
	}
	for key, value := range attributes {
		clusterUnitTestSetAttribute(t, stateData, key, value)
	}
	state := stateData.State()
	if state == nil {
		t.Fatal("failed to build the unit test instance state")
	}
	if changes == nil {
		d, err := schema.InternalMap(resourceSchema).Data(state, nil)
		if err != nil {
			t.Fatalf("failed to build the unit test resource data: %s", err)
		}
		return d
	}
	oldData, err := schema.InternalMap(resourceSchema).Data(state, nil)
	if err != nil {
		t.Fatalf("failed to build the prior unit test resource data: %s", err)
	}
	newData, err := schema.InternalMap(resourceSchema).Data(state, nil)
	if err != nil {
		t.Fatalf("failed to build the changed unit test resource data: %s", err)
	}
	for key, value := range changes {
		clusterUnitTestSetAttribute(t, newData, key, value)
	}
	diff := sdkterraform.NewInstanceDiff()
	for _, key := range helper.MapSort(newData.State().Attributes) {
		newValue := newData.State().Attributes[key]
		oldValue, ok := oldData.State().Attributes[key]
		if ok && oldValue == newValue {
			continue
		}
		if oldValue == "" && (strings.HasSuffix(key, "#") || strings.HasSuffix(key, "%")) {
			oldValue = "0"
		}
		diff.Attributes[key] = &sdkterraform.ResourceAttrDiff{Old: oldValue, New: newValue}
	}
	if configured != nil {
		diff.RawConfig = clusterUnitTestRawConfig(t, configured)
	}
	d, err := schema.InternalMap(resourceSchema).Data(state, diff)
	if err != nil {
		t.Fatalf("failed to build the diffed unit test resource data: %s", err)
	}
	return d
}

// clusterUnitTestRawConfig builds the resource configuration Terraform sends to
// the provider: every attribute has its schema type and only the listed ones
// carry a value, so unlisted attributes are null exactly like omitted arguments.
func clusterUnitTestRawConfig(t *testing.T, configured map[string]cty.Value) cty.Value {
	t.Helper()
	impliedType := polardbsvc.ResourceCluster().CoreConfigSchema().ImpliedType()
	values := make(map[string]cty.Value, len(impliedType.AttributeTypes()))
	for name, attributeType := range impliedType.AttributeTypes() {
		values[name] = cty.NullVal(attributeType)
	}
	for name, value := range configured {
		if _, ok := values[name]; !ok {
			t.Fatalf("raw config attribute %q is not part of the resource schema", name)
		}
		values[name] = value
	}
	return cty.ObjectVal(values)
}

// tfproviderlint R001 requires literal ResourceData.Set keys, so the fixture maps are applied through an explicit key switch.
func clusterUnitTestSetAttribute(t *testing.T, d *schema.ResourceData, key string, value interface{}) {
	t.Helper()
	var err error
	switch key {
	case "cn_node_class":
		err = d.Set("cn_node_class", value)
	case "dn_node_class":
		err = d.Set("dn_node_class", value)
	case "cn_node_ids":
		err = d.Set("cn_node_ids", value)
	case "dn_node_ids":
		err = d.Set("dn_node_ids", value)
	case "db_type":
		err = d.Set("db_type", value)
	case "pay_type":
		err = d.Set("pay_type", value)
	case "creation_category":
		err = d.Set("creation_category", value)
	case "cn_node_num":
		err = d.Set("cn_node_num", value)
	case "dn_node_num":
		err = d.Set("dn_node_num", value)
	case "parameters":
		err = d.Set("parameters", value)
	case "tde_status":
		err = d.Set("tde_status", value)
	case "encrypt_new_tables":
		err = d.Set("encrypt_new_tables", value)
	case "encryption_key":
		err = d.Set("encryption_key", value)
	case "role_arn":
		err = d.Set("role_arn", value)
	case "enable_automatic_rotation":
		err = d.Set("enable_automatic_rotation", value)
	default:
		t.Fatalf("unsupported unit test attribute %q", key)
	}
	if err != nil {
		t.Fatalf("failed to set %s in the unit test state: %s", key, err)
	}
}

func clusterUnitTestSetValues(t *testing.T, d *schema.ResourceData, key string) []string {
	t.Helper()
	raw, ok := d.Get(key).(*schema.Set)
	if !ok {
		t.Fatalf("%s is not a set", key)
	}
	values := make([]string, 0, raw.Len())
	for _, item := range raw.List() {
		values = append(values, fmt.Sprint(item))
	}
	sort.Strings(values)
	return values
}
