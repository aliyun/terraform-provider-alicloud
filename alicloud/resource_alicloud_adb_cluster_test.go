package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var adbClusterConnectionStringRegexp = "am-[a-z-A-Z-0-9]+.ads.aliyuncs.com"

func init() {
	resource.AddTestSweepers("alicloud_adb_cluster", &resource.Sweeper{
		Name: "alicloud_adb_cluster",
		F:    testSweepAdbCluster,
	})
}

func testSweepAdbCluster(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	request := adb.CreateDescribeDBClustersRequest()
	request.PageSize = requests.NewInteger(PageSizeXLarge)
	raw, err := client.WithAdbClient(func(AdbClient *adb.Client) (interface{}, error) {
		return AdbClient.DescribeDBClusters(request)
	})
	if err != nil {
		log.Printf("[ERROR] Error retrieving Adb Clusters: %s", WrapError(err))
	}
	response, _ := raw.(*adb.DescribeDBClustersResponse)

	sweeped := false
	for _, v := range response.Items.DBCluster {
		id := v.DBClusterId
		name := v.DBClusterDescription
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Abd Clusters: %s (%s)", name, id)
			continue
		}

		sweeped = true
		log.Printf("[INFO] Deleting Adb Clusters: %s (%s)", name, id)
		req := adb.CreateDeleteDBClusterRequest()
		req.DBClusterId = id
		_, err := client.WithAdbClient(func(AdbClient *adb.Client) (interface{}, error) {
			return AdbClient.DeleteDBCluster(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Adb Clusters (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		// Waiting 30 seconds to ensure these Adb Clusters have been deleted.
		time.Sleep(30 * time.Second)
	}
	return nil
}

func TestAccAlicloudAdbCluster(t *testing.T) {
	var v *adb.DBCluster
	var ips []map[string]interface{}
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sadbrecordbasic%v.abc", defaultRegionToTest, rand)
	resourceId := "alicloud_adb_cluster.default"
	var basicMap = map[string]string{
		"description":         CHECKSET,
		"vswitch_id":          CHECKSET,
		"db_cluster_category": CHECKSET,
		"db_node_class":       CHECKSET,
		"db_node_count":       CHECKSET,
		"db_node_storage":     CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAdbClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAdbClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithNoDefaultVpc(t)
			testAccPreCheckWithNoDefaultVswitch(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":         "${var.name}",
					"vswitch_id":          "${data.alicloud_vswitches.default.ids.0}",
					"db_cluster_category": "Cluster",
					"db_node_class":       "C8",
					"db_node_count":       "2",
					"db_node_storage":     "200",
					"pay_type":            "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_string": REGEXMATCH + adbClusterConnectionStringRegexp,
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
					"description": "tf-testaccadbclusterbasic",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testaccadbclusterbasic",
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
					"security_ips": []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyValueInMapsForAdb(ips, "security ip", "security_ips", "10.168.1.12,100.69.7.112"),
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
					"db_node_count": "4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_count": "4",
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
					"description":   "tf-testaccadbrecordbasic1",
					"maintain_time": "02:00Z-03:00Z",
					"security_ips":  []string{"10.168.1.13", "100.69.7.113"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":   "tf-testaccadbrecordbasic1",
						"maintain_time": "02:00Z-03:00Z",
					}),
					testAccCheckKeyValueInMapsForAdb(ips, "security ip", "security_ips", "10.168.1.13,100.69.7.113"),
				),
			},
		},
	})

}

func TestAccAlicloudAdbClusterMulti(t *testing.T) {
	var v *adb.DBCluster
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sadbrecordbasic%v.abc", defaultRegionToTest, rand)
	resourceId := "alicloud_adb_cluster.default.1"
	var basicMap = map[string]string{
		"description":         CHECKSET,
		"vswitch_id":          CHECKSET,
		"db_cluster_category": CHECKSET,
		"db_node_class":       CHECKSET,
		"db_node_count":       CHECKSET,
		"db_node_storage":     CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAdbClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAdbClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithNoDefaultVpc(t)
			testAccPreCheckWithNoDefaultVswitch(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":               "2",
					"description":         "${var.name}",
					"vswitch_id":          "${data.alicloud_vswitches.default.ids.0}",
					"db_cluster_category": "Cluster",
					"db_node_class":       "C8",
					"db_node_count":       "2",
					"db_node_storage":     "200",
					"pay_type":            "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccCheckKeyValueInMapsForAdb(ps []map[string]interface{}, propName, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, policy := range ps {
			if policy[key].(string) != value {
				return fmt.Errorf("DB %s attribute '%s' expected %#v, got %#v", propName, key, value, policy[key])
			}
		}
		return nil
	}
}

func resourceAdbClusterConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "ADB"
	}
	
	variable "name" {
		default = "%s"
	}
	`, AdbCommonTestCase, name)
}
