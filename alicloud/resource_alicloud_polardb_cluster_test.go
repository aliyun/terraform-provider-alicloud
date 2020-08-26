package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var clusterConnectionStringRegexp = "^[a-z-A-Z-0-9]+.rwlb.[a-z]+.rds.aliyuncs.com"

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

func TestAccAlicloudPolarDBClusterUpdate(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	var ips []map[string]interface{}
	name := "tf-testAccPolarDBClusterUpdate"
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
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithNoDefaultVpc(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_type":       "MySQL",
					"db_version":    "8.0",
					"pay_type":      "PostPaid",
					"db_node_class": "polar.mysql.x4.medium",
					"vswitch_id":    "${data.alicloud_vswitches.default.ids.0}",
					"description":   "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_string": "",
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
					"description": "tf-testaccPolarDBClusterUpdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testaccPolarDBClusterUpdate",
					}),
				),
			},
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
					"db_node_class": "polar.mysql.x4.large",
					"modify_type":   "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_count": "2",
						"db_node_class": "polar.mysql.x4.large",
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
						"connection_string": REGEXMATCH + clusterConnectionStringRegexp,
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
					"description":   "tf-testaccPolarDBClusterUpdate1",
					"maintain_time": "02:00Z-03:00Z",
					"db_node_class": "polar.mysql.x8.xlarge",
					"modify_type":   "Upgrade",
					"security_ips":  []string{"10.168.1.13", "100.69.7.113"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":   "tf-testaccPolarDBClusterUpdate1",
						"maintain_time": "02:00Z-03:00Z",
						"db_node_class": "polar.mysql.x8.xlarge",
					}),
					testAccCheckKeyValueInMapsForPolarDB(ips, "security ip", "security_ips", "10.168.1.13,100.69.7.113"),
				),
			},
		},
	})

}

func TestAccAlicloudPolarDBClusterMulti(t *testing.T) {
	var v *polardb.DescribeDBClusterAttributeResponse
	name := "tf-testaccPolarDBClusterMult"
	resourceId := "alicloud_polardb_cluster.default.2"
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
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithNoDefaultVpc(t)
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
					"db_node_class": "polar.mysql.x4.large",
					"vswitch_id":    "${data.alicloud_vswitches.default.ids.0}",
					"description":   "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
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

func resourcePolarDBClusterConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "PolarDB"
	}

	variable "name" {
		default = "%s"
	}

`, PolarDBCommonTestCase, name)
}
