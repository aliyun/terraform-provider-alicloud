package rds_test

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/rds"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Attributes that cannot be compared through an import:
//   - effective_time, effective_specific_time, upgrade_time and switch_time are
//     local scheduling settings that the proxy API does not return;
//   - db_proxy_features is write-only configuration that DescribeDBProxy does
//     not return;
//   - resource_group_id is applied only when the proxy is created: the update
//     path does not change it and Read does not restore it, so restoring it in
//     Read without an in-place update would leave a permanent diff. That is
//     pre-existing behavior unrelated to the attributes verified here.
var dbProxyImportStateVerifyIgnore = []string{"effective_time", "effective_specific_time", "upgrade_time", "switch_time", "db_proxy_features", "resource_group_id"}

// db_proxy_nodes is additionally empty after an import until it is configured
// again, which is stated in the resource documentation.
var dbProxyNodesImportStateVerifyIgnore = []string{"effective_time", "effective_specific_time", "upgrade_time", "switch_time", "db_proxy_features", "resource_group_id", "db_proxy_nodes"}

func TestAccAliCloudRdsDBProxy_MySQL(t *testing.T) {
	var connection map[string]interface{}
	var primary map[string]interface{}
	var readonly map[string]interface{}
	var proxyEndpointIdBefore string
	var proxyNodeLayoutBefore []string

	resourceId := "alicloud_rds_db_proxy.default"
	var DBProxyMap = map[string]string{
		"instance_id":           acctest.CHECKSET,
		"instance_network_type": "VPC",
		"db_proxy_instance_num": "2",
		"vswitch_id":            acctest.CHECKSET,
		"vpc_id":                acctest.CHECKSET,
	}
	ra := acctest.ResourceAttrInit(resourceId, DBProxyMap)

	rc_connection := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &connection, func() interface{} {
		return rds.NewRdsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRdsProxyEndpoint")
	rc_primary := acctest.ResourceCheckInitWithDescribeMethod("alicloud_db_instance.default", &primary, func() interface{} {
		return rds.NewRdsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDBInstance")
	rc_readonly := acctest.ResourceCheckInitWithDescribeMethod("alicloud_db_readonly_instance.default", &readonly, func() interface{} {
		return rds.NewRdsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDBReadonlyInstance")
	rand := tfacctest.RandString(5)
	// The effective specific time is only carried while EffectiveTime is Immediate,
	// because a deferred change keeps the proxy in a changing state until the
	// requested time and a wall clock dependent step is not reproducible.
	effectiveSpecificTime := time.Now().UTC().Add(45 * time.Minute).Format("2006-01-02T15:04:05Z")
	switchTime := time.Now().UTC().Add(45 * time.Minute).Format("2006-01-02T15:04:05Z")
	rac := acctest.ResourceAttrCheckInit(rc_connection, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBProxy%s", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceRdsDBProxyDependence_MySQL)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":            "${alicloud_db_instance.default.id}",
					"instance_network_type":  "VPC",
					"db_proxy_instance_num":  "2",
					"db_proxy_instance_type": "common",
					"vswitch_id":             "${alicloud_db_instance.default.vswitch_id}",
					"vpc_id":                 "${alicloud_db_instance.default.vpc_id}",
					"depends_on":             []string{"alicloud_db_readonly_instance.default"},
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					rc_primary.CheckResourceExists(),
					rc_readonly.CheckResourceExists(),
					testAccCheck(map[string]string{
						"instance_id":            acctest.CHECKSET,
						"db_proxy_instance_type": acctest.CHECKSET,
						"instance_network_type":  "VPC",
						"db_proxy_instance_num":  "2",
						"vswitch_id":             acctest.CHECKSET,
						"vpc_id":                 acctest.CHECKSET,
						"resource_group_id":      acctest.CHECKSET,
					}),
					captureDBProxyEndpointId(resourceId, &proxyEndpointIdBefore),
				),
			},
			{
				// SpecificTime requires the point in time, and the guard rejects the
				// missing value before ModifyDBProxyInstance is called. This step runs
				// before any step configures effective_specific_time: the attribute is
				// Optional+Computed, so a value planned by an earlier step survives an
				// omitted or empty configuration value.
				Config: testAccConfig(map[string]interface{}{
					"db_proxy_instance_type": "exclusive",
					"effective_time":         "SpecificTime",
				}),
				ExpectError: regexp.MustCompile(`effective_specific_time must be set when effective_time is SpecificTime`),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionUpdate),
					},
				},
				Config: testAccConfig(map[string]interface{}{
					// Restore the values set by step 1: a failed apply does not roll
					// back the cumulative configuration accumulator.
					"db_proxy_instance_type": "common",
					"effective_time":         "Immediate",
					"vswitch_id":             "${alicloud_vswitch.second.id}",
					"vpc_id":                 "${alicloud_vpc.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id": acctest.CHECKSET,
					}),
					resource.TestCheckResourceAttrPair(resourceId, "vswitch_id", "alicloud_vswitch.second", "id"),
					checkDBProxyEndpointIdUnchanged(resourceId, &proxyEndpointIdBefore),
				),
			},
			{
				// The migrated proxy has to survive an import with the same state.
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: dbProxyImportStateVerifyIgnore,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_proxy_instance_num":   "2",
					"db_proxy_instance_type":  "exclusive",
					"effective_time":          "Immediate",
					"effective_specific_time": effectiveSpecificTime,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_proxy_instance_num":   "2",
						"db_proxy_instance_type":  "exclusive",
						"effective_time":          "Immediate",
						"effective_specific_time": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"upgrade_time": "Immediate",
					"switch_time":  switchTime,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"upgrade_time": "Immediate",
						"switch_time":  acctest.CHECKSET,
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
			{
				Config: testAccConfig(map[string]interface{}{
					"db_proxy_endpoint_read_write_mode": "ReadWrite",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_proxy_endpoint_read_write_mode": "ReadWrite",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"read_only_instance_distribution_type": "Standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"read_only_instance_distribution_type": "Standard",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"read_only_instance_max_delay_time": "90",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"read_only_instance_max_delay_time": "90",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_proxy_features": "TransactionReadSqlRouteOptimizeStatus:1;ConnectionPersist:1;ReadWriteSpliting:0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_proxy_features": "TransactionReadSqlRouteOptimizeStatus:1;ConnectionPersist:1;ReadWriteSpliting:0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"read_only_instance_distribution_type": "Custom",
					"db_proxy_features":                    "TransactionReadSqlRouteOptimizeStatus:1;ConnectionPersist:1;ReadWriteSpliting:1",
					"read_only_instance_weight": []map[string]interface{}{
						{
							"instance_id": "${alicloud_db_instance.default.id}",
							"weight":      "100",
						},
						{
							"instance_id": "${alicloud_db_readonly_instance.default.id}",
							"weight":      "500",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"read_only_instance_distribution_type": "Custom",
						"db_proxy_features":                    "TransactionReadSqlRouteOptimizeStatus:1;ConnectionPersist:1;ReadWriteSpliting:1",
						"read_only_instance_weight.#":          "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_proxy_instance_num":  "3",
					"effective_time":         "Immediate",
					"db_proxy_instance_type": "exclusive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_proxy_instance_num":  "3",
						"effective_time":         "Immediate",
						"db_proxy_instance_type": "exclusive",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_proxy_connect_string_port": "3306",
					"db_proxy_connection_prefix":   fmt.Sprintf("tf-testacc%s", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_proxy_connect_string_port": "3306",
						"db_proxy_connection_prefix":   acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_proxy_ssl_enabled": "Open",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_proxy_ssl_enabled": "Open",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"persistent_connection_status": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"persistent_connection_status": "Enabled",
					}),
					checkDBProxyRemotePersistentConnectionStatus(resourceId, "Enabled"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: dbProxyImportStateVerifyIgnore,
			},
			{
				// Removing the attribute while it is Enabled stops Terraform from
				// managing the switch: the enabled setting is kept instead of being
				// reset, which is what the documented omission semantics promise.
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionNoop),
					},
				},
				Config: testAccConfig(map[string]interface{}{
					"persistent_connection_status": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "persistent_connection_status", "Enabled"),
					checkDBProxyRemotePersistentConnectionStatus(resourceId, "Enabled"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: dbProxyImportStateVerifyIgnore,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"persistent_connection_status": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"persistent_connection_status": "Disabled",
					}),
					checkDBProxyRemotePersistentConnectionStatus(resourceId, "Disabled"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: dbProxyImportStateVerifyIgnore,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_proxy_instance_num":  "2",
					"db_proxy_instance_type": "common",
					"db_proxy_nodes": []map[string]interface{}{
						{
							"cpu_cores":   "2",
							"node_counts": "2",
							"zone_id":     "${local.proxy_zone_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_proxy_instance_num":  "2",
						"db_proxy_instance_type": "common",
						"db_proxy_nodes.#":       "1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "db_proxy_nodes.*", map[string]string{
						"cpu_cores":   "2",
						"node_counts": "2",
					}),
					captureDBProxyNodeLayout(resourceId, &proxyNodeLayoutBefore),
				),
			},
			{
				// The block is not restored by an import until it is configured
				// again, which is documented behavior, so only it is ignored here.
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: dbProxyNodesImportStateVerifyIgnore,
			},
			{
				// Removing the block stops Terraform from managing the node layout
				// and leaves the existing layout unchanged.
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionUpdate),
					},
				},
				Config: testAccConfig(map[string]interface{}{
					"db_proxy_nodes": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "db_proxy_nodes.#", "0"),
					checkDBProxyNodeLayoutUnchanged(resourceId, &proxyNodeLayoutBefore),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: dbProxyImportStateVerifyIgnore,
			},
			{
				// node_counts has no default, so an incomplete block is rejected before
				// the API call instead of producing a value the configuration omits.
				Config: testAccConfig(map[string]interface{}{
					"db_proxy_nodes": []map[string]interface{}{
						{
							"cpu_cores": "2",
							"zone_id":   "${local.proxy_zone_id}",
						},
					},
				}),
				ExpectError: regexp.MustCompile(`db_proxy_nodes requires zone_id, cpu_cores and node_counts`),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"switch_time", "effective_specific_time"},
			},
		},
	})
}

func captureDBProxyEndpointId(resourceId string, target *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		endpointId := rs.Primary.Attributes["db_proxy_endpoint_id"]
		if endpointId == "" {
			return fmt.Errorf("resource %s has no proxy endpoint id", resourceId)
		}
		*target = endpointId
		return nil
	}
}

// A vSwitch update is applied as an in-place migration, so the proxy endpoint
// survives the change; a replacement would create a new proxy endpoint.
func checkDBProxyEndpointIdUnchanged(resourceId string, expected *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		if *expected == "" {
			return fmt.Errorf("proxy endpoint id was not captured before the migration")
		}
		if endpointId := rs.Primary.Attributes["db_proxy_endpoint_id"]; endpointId != *expected {
			return fmt.Errorf("proxy endpoint id changed from %q to %q, the proxy was replaced instead of migrated", *expected, endpointId)
		}
		return nil
	}
}

// describeDBProxyForCheck queries the proxy directly so the behavior checks do
// not reuse the resource Read path they verify.
func describeDBProxyForCheck(resourceId string, state *terraform.State) (map[string]interface{}, error) {
	rs, ok := state.RootModule().Resources[resourceId]
	if !ok {
		return nil, fmt.Errorf("resource %s not found in state", resourceId)
	}
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	response, err := client.RpcPost("Rds", "2014-08-15", "DescribeDBProxy", nil, map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": rs.Primary.ID,
		"SourceIp":     client.SourceIp,
	}, true)
	if err != nil {
		return nil, fmt.Errorf("DescribeDBProxy %s: %w", rs.Primary.ID, err)
	}
	return response, nil
}

func checkDBProxyRemotePersistentConnectionStatus(resourceId, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		response, err := describeDBProxyForCheck(resourceId, s)
		if err != nil {
			return err
		}
		if status := fmt.Sprint(response["DBProxyPersistentConnectionStatus"]); status != expected {
			return fmt.Errorf("remote persistent connection status is %q, want %q", status, expected)
		}
		return nil
	}
}

func dbProxyNodeLayoutOfResponse(response map[string]interface{}) ([]string, error) {
	wrapper, ok := response["DBProxyNodes"].(map[string]interface{})
	if !ok {
		return nil, nil
	}
	rawNodes, _ := wrapper["DBProxyNodes"].([]interface{})
	layout := make([]string, 0, len(rawNodes))
	for _, raw := range rawNodes {
		node, ok := raw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected proxy node entry in the response: %T", raw)
		}
		layout = append(layout, fmt.Sprintf("%s/%s", fmt.Sprint(node["zoneId"]), fmt.Sprint(node["cpuCores"])))
	}
	sort.Strings(layout)
	return layout, nil
}

func captureDBProxyNodeLayout(resourceId string, target *[]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		response, err := describeDBProxyForCheck(resourceId, s)
		if err != nil {
			return err
		}
		layout, err := dbProxyNodeLayoutOfResponse(response)
		if err != nil {
			return err
		}
		if len(layout) == 0 {
			return fmt.Errorf("the proxy did not report any proxy node")
		}
		*target = layout
		return nil
	}
}

func checkDBProxyNodeLayoutUnchanged(resourceId string, expected *[]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(*expected) == 0 {
			return fmt.Errorf("the proxy node layout was not captured before db_proxy_nodes was removed")
		}
		response, err := describeDBProxyForCheck(resourceId, s)
		if err != nil {
			return err
		}
		layout, err := dbProxyNodeLayoutOfResponse(response)
		if err != nil {
			return err
		}
		if strings.Join(layout, ",") != strings.Join(*expected, ",") {
			return fmt.Errorf("the remote proxy node layout changed from %v to %v after db_proxy_nodes was removed", *expected, layout)
		}
		return nil
	}
}

func resourceRdsDBProxyDependence_MySQL(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "8.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "8.0"
    category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_db_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_vswitch" "second" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.1.0/24"
  zone_id      = local.proxy_zone_id
  vswitch_name = "${var.name}-second"
}

locals {
  proxy_zone_id = length(data.alicloud_db_zones.default.zones) > 1 ? data.alicloud_db_zones.default.zones.1.id : data.alicloud_db_zones.default.zones.0.id
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_db_instance" "default" {
    engine = "MySQL"
	engine_version = "8.0"
 	db_instance_storage_type = "cloud_essd"
	instance_type = "mysql.x2.large.2c"
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
	vswitch_id = alicloud_vswitch.default.id
	instance_name = var.name
}

resource "alicloud_db_readonly_instance" "default" {
	master_db_instance_id = alicloud_db_instance.default.id
	zone_id = alicloud_db_instance.default.zone_id
	engine_version = alicloud_db_instance.default.engine_version
	instance_type = "mysqlro.n4.medium.1c"
	instance_storage = alicloud_db_instance.default.instance_storage
	instance_name = "${var.name}_ro"
	vswitch_id = alicloud_db_instance.default.vswitch_id
}

`, name)
}
