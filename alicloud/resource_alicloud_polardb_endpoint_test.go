package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_polardb_endpoint", &resource.Sweeper{
		Name: "alicloud_polardb_endpoint",
		F:    testSweepPolarDBEndpoints,
	})
}

// testSweepPolarDBEndpoints walks every PolarDB cluster in the region and
// deletes Custom endpoints on it (Primary/Cluster endpoints cannot be
// deleted, see resourceAlicloudPolarDBEndpointDelete). A cluster whose
// description carries a do-not-delete marker (protectedSweepVpcName) is
// skipped entirely, as is any endpoint whose description carries such a
// marker; both checks run in every mode so SWEEP_ALL cannot bypass them.
// When SWEEP_ALL is unset, only endpoints on clusters (and with descriptions)
// matching the standard test prefixes are touched. Per-instance failures are
// logged and non-fatal so one bad cluster cannot block the rest of the sweep.
func testSweepPolarDBEndpoints(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{"tf-testAcc", "tf_testAcc"}

	clusters, err := listPolarDBClustersForSweep(client)
	if err != nil {
		return err
	}

	for _, c := range clusters {
		clusterName := c.DBClusterDescription
		clusterId := c.DBClusterId
		if protectedSweepVpcName(clusterName) {
			log.Printf("[INFO] Skipping PolarDB endpoints on NODELETING cluster: %s (%s)", clusterName, clusterId)
			continue
		}
		if !sweepAll() {
			match := false
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(clusterName), strings.ToLower(prefix)) {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}

		ereq := polardb.CreateDescribeDBClusterEndpointsRequest()
		ereq.RegionId = client.RegionId
		ereq.DBClusterId = clusterId
		raw, err := client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
			return polardbClient.DescribeDBClusterEndpoints(ereq)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to list PolarDB endpoints on cluster %s: %s", clusterId, err)
			continue
		}
		eresp, _ := raw.(*polardb.DescribeDBClusterEndpointsResponse)
		if eresp == nil {
			continue
		}
		for _, ep := range eresp.Items {
			// Only Custom endpoints can be deleted; Primary/Cluster endpoints
			// are managed by the cluster itself.
			if ep.EndpointType != "Custom" {
				continue
			}
			if protectedSweepVpcName(ep.DBEndpointDescription) {
				log.Printf("[INFO] Skipping PolarDB endpoint (NODELETING): %s on %s", ep.DBEndpointDescription, clusterId)
				continue
			}
			if !sweepAll() {
				match := false
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(ep.DBEndpointDescription), strings.ToLower(prefix)) {
						match = true
						break
					}
				}
				if !match {
					continue
				}
			}
			log.Printf("[INFO] Deleting PolarDB endpoint: %s (%s) on cluster %s", ep.DBEndpointDescription, ep.DBEndpointId, clusterId)
			delReq := polardb.CreateDeleteDBClusterEndpointRequest()
			delReq.RegionId = client.RegionId
			delReq.DBClusterId = clusterId
			delReq.DBEndpointId = ep.DBEndpointId
			if _, err := client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
				return polardbClient.DeleteDBClusterEndpoint(delReq)
			}); err != nil {
				log.Printf("[ERROR] Failed to delete PolarDB endpoint %s on %s: %s", ep.DBEndpointId, clusterId, err)
			}
		}
	}
	return nil
}

var privateConnectionStringRegexp = "^[a-z-A-Z-0-9]+.rwlb.([a-z-A-Z-0-9]+.){0,1}rds.aliyuncs.com"

func TestAccAliCloudPolarDBEndpointConfigUpdate(t *testing.T) {
	var v *polardb.DBEndpoint
	rand := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBendpoint-%s", rand)
	var basicMap = map[string]string{
		"db_cluster_id": CHECKSET,
	}
	resourceId := "alicloud_polardb_endpoint.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBEndpointConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id": "${alicloud_polardb_cluster.cluster.id}",
					"endpoint_type": "Custom",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_type": "Custom",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_add_new_nodes": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_add_new_nodes": "Enable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"read_write_mode": "ReadWrite",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"read_write_mode": "ReadWrite",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nodes": []string{"${data.alicloud_polardb_clusters.default.clusters.0.db_nodes.0.db_node_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nodes.#": "1",
					}),
				),
			},
			//todo: After resource polardb_node is supported, it is necessary to add a modification check on the “nodes” parameter
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enabled":     "Enable",
					"net_type":        "Private",
					"ssl_auto_rotate": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enabled":           "Enable",
						"net_type":              "Private",
						"ssl_connection_string": REGEXMATCH + privateConnectionStringRegexp,
						"ssl_auto_rotate":       "Enable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_endpoint_description": "polar_db_endpoint_test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_endpoint_description": "polar_db_endpoint_test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_prefix": "tf-privatetestacc-3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_prefix": "tf-privatetestacc-3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port": "3307",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "3307",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_config": map[string]string{
						"ConsistLevel": "1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_config.ConsistLevel": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"net_type", "endpoint_config"},
			},
		},
	})
}

func TestAccAliCloudPolarDBEndpointConfigUpdate_SslConnectionStringAndConnectionPrefix(t *testing.T) {
	var v *polardb.DBEndpoint
	rand := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBendpoint-%s", rand)
	var basicMap = map[string]string{
		"db_cluster_id": CHECKSET,
	}
	resourceId := "alicloud_polardb_endpoint.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBEndpointConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id": "${alicloud_polardb_cluster.cluster.id}",
					"endpoint_type": "Custom",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_type": "Custom",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_add_new_nodes": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_add_new_nodes": "Enable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"read_write_mode": "ReadWrite",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"read_write_mode": "ReadWrite",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nodes": []string{"${data.alicloud_polardb_clusters.default.clusters.0.db_nodes.0.db_node_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nodes.#": "1",
					}),
				),
			},
			//todo: After resource polardb_node is supported, it is necessary to add a modification check on the “nodes” parameter
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_prefix": "tf-privatetestacc-4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_prefix": "tf-privatetestacc-4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enabled":     "Enable",
					"net_type":        "Private",
					"ssl_auto_rotate": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enabled":           "Enable",
						"net_type":              "Private",
						"ssl_connection_string": REGEXMATCH + privateConnectionStringRegexp,
						"ssl_auto_rotate":       "Enable",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"net_type", "endpoint_config"},
			},
		},
	})
}

func resourcePolarDBEndpointConfigDependence(name string) string {
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

        resource "alicloud_polardb_cluster" "cluster" {
                db_type = "MySQL"
                db_version = "8.0"
                pay_type = "PostPaid"
                db_node_class = "polar.mysql.x4.medium"
                vswitch_id = alicloud_vswitch.default.id
                description = "${var.name}"
                lifecycle {
                        ignore_changes = [encrypt_new_tables]
                }
        }

		data "alicloud_polardb_clusters" "default" {
		  	ids = [alicloud_polardb_cluster.cluster.id]
		}
        `, name)
}

// Test PolarDB Endpoint DynamoDB type
func TestAccAliCloudPolarDBEndpoint_dynamoDB(t *testing.T) {
	var v *polardb.DBEndpoint
	rand := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBendpoint-dynamodb-%s", rand)
	var basicMap = map[string]string{
		"db_cluster_id": CHECKSET,
	}
	resourceId := "alicloud_polardb_endpoint.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBEndpointDynamoDBDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id": "${alicloud_polardb_cluster.default.id}",
					"endpoint_type": "DynamoDB",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_id": CHECKSET,
						"endpoint_type": "DynamoDB",
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

func resourcePolarDBEndpointDynamoDBDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.7.id
  vswitch_name = var.name
}

resource "alicloud_polardb_cluster" "default" {
  db_type         = "PostgreSQL"
  db_version      = "14"
  db_node_class   = "polar.pg.x4.medium"
  pay_type        = "PostPaid"
  vswitch_id      = alicloud_vswitch.default.id
  enable_dynamodb = true
  lifecycle {
    ignore_changes = [encrypt_new_tables]
  }
}
`, name)
}
