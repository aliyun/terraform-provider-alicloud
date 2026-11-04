package selectdb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/selectdb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_selectdb_db_cluster", &resource.Sweeper{
		Name: "alicloud_selectdb_db_cluster",
		F:    testSweepSelectDBDbClusterAndDbInstance,
	})
}

func testSweepSelectDBDbClusterAndDbInstance(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	selectDBService := selectdb.NewService(client)
	instanceResp, err := selectDBService.DescribeDbInstances("", nil)
	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_selectdb_db_instances", errs.AlibabaCloudSdkGoERROR)
	}

	var objects []map[string]interface{}

	for _, item := range instanceResp {
		name := item["Description"].(string)
		skip := true
		if !sweep.SweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(name, prefix) {
					skip = false
					break
				}
			}
			if skip {
				continue
			}
		}
		objects = append(objects, item)
	}

	for _, id := range objects {
		_, err := selectDBService.DeleteInstance(id["DBInstanceId"].(string))
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_selectdb_db_instances_clusters", errs.AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

func TestAccAliCloudSelectDBDbCluster_basic_info(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_selectdb_db_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudSelectDBDbClusterMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return selectdb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDbCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testAccselectdbcluster%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudSelectDBDbClusterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_class":       "selectdb.g.2xlarge",
					"cache_size":             "200",
					"payment_type":           "PayAsYouGo",
					"db_cluster_description": name,
					"db_instance_id":         "${alicloud_selectdb_db_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_class":       "selectdb.g.2xlarge",
						"cache_size":             "200",
						"payment_type":           "PayAsYouGo",
						"db_cluster_description": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_description": name + "_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_description": name + "_desc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desired_params": []map[string]interface{}{
						{
							"name":  "es_http_timeout_ms",
							"value": "6000",
						},
						{
							"name":  "sys_log_roll_num",
							"value": "12",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desired_params.#": "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"db_cluster_description", "cache_size", "db_cluster_class", "desired_params"},
			},
		},
	})
}

func TestAccAliCloudSelectDBDbCluster_status(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_selectdb_db_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudSelectDBDbClusterMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return selectdb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDbCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testAccselectdbcluster%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudSelectDBDbClusterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_class":       "selectdb.g.2xlarge",
					"cache_size":             "200",
					"payment_type":           "PayAsYouGo",
					"db_cluster_description": name,
					"db_instance_id":         "${alicloud_selectdb_db_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_class":       "selectdb.g.2xlarge",
						"cache_size":             "200",
						"payment_type":           "PayAsYouGo",
						"db_cluster_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desired_status": "STOPPING",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desired_status": "STOPPING",
						"status":         "STOPPED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desired_status": "STARTING",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desired_status": "STARTING",
						"status":         "ACTIVATION",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desired_status": "RESTART",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desired_status": "RESTART",
						"status":         "ACTIVATION",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"status", "cache_size", "desired_status"},
			},
		},
	})
}

func TestAccAliCloudSelectDBDbCluster_modify(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_selectdb_db_cluster.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudSelectDBDbClusterMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return selectdb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDbCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testAccselectdbcluster%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudSelectDBDbClusterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_class":       "selectdb.g.2xlarge",
					"cache_size":             "200",
					"payment_type":           "PayAsYouGo",
					"db_cluster_description": name,
					"db_instance_id":         "${alicloud_selectdb_db_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_class":       "selectdb.g.2xlarge",
						"cache_size":             "200",
						"payment_type":           "PayAsYouGo",
						"db_cluster_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cache_size": "400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cache_size": "400",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_class": "selectdb.g.2xlarge",
					"cache_size":       "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_class": "selectdb.g.2xlarge",
						"cache_size":       "1000",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"db_cluster_class", "cache_size"},
			},
		},
	})
}

var AliCloudSelectDBDbClusterMap0 = map[string]string{
	"db_instance_id":         acctest.CHECKSET,
	"db_cluster_description": acctest.CHECKSET,
	"payment_type":           acctest.CHECKSET,
	"db_cluster_class":       acctest.CHECKSET,
	"cache_size":             acctest.CHECKSET,
}

func AliCloudSelectDBDbClusterBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.ids.0
}
resource "alicloud_selectdb_db_instance" "default" {
  db_instance_class       = "selectdb.g.2xlarge"
  db_instance_description = var.name
  cache_size              = "200"
  engine_minor_version    = "4.0"
  payment_type            = "PayAsYouGo"
  vpc_id                  = data.alicloud_vswitches.default.vswitches.0.vpc_id
  zone_id                 = data.alicloud_vswitches.default.vswitches.0.zone_id
  vswitch_id              = data.alicloud_vswitches.default.vswitches.0.id
}
`, name)
}
