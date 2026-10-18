// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// testSweepStarRocksNodeGroup <<< Sweeper for StarRocks NodeGroup.
// node_groups are cluster-scoped: the API lists node_groups under a specific
// cluster (describeNodeGroups needs ClusterId), so this sweeper parent-walks —
// list every instance, then list node_groups under each instance — and deletes
// them with /webapi/lifecycle/releaseNodeGroup. The NODELETING guard
// (protectedSweepVpcName) is always active; NORMAL mode adds a tftestacc prefix
// filter, SWEEP_ALL mode cleans every non-NODELETING node group.
func init() {
	resource.AddTestSweepers(
		"alicloud_star_rocks_node_group",
		&resource.Sweeper{
			Name: "alicloud_star_rocks_node_group",
			F:    testSweepStarRocksNodeGroup,
		})
}

func testSweepStarRocksNodeGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	// node_groups are cluster-scoped, so walk all instances first.
	listAction := "/webapi/starrocks/describeInstances"
	listQuery := map[string]*string{
		"RegionId": StringPointer(client.RegionId),
	}
	listBody := make(map[string]interface{})
	var listResp map[string]interface{}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(2*time.Minute, func() *retry.RetryError {
		resp, e := client.RoaPost("starrocks", "2022-10-19", listAction, listQuery, nil, listBody, true)
		listResp = resp
		if e != nil {
			if NeedRetry(e) {
				wait()
				return retry.RetryableError(e)
			}
			return retry.NonRetryableError(e)
		}
		return nil
	})
	addDebug(listAction, listResp, listBody)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch StarRocks Instances for node_group walk: %s", WrapErrorf(err, DataDefaultErrorMsg, "alicloud_star_rocks_node_group", listAction, AlibabaCloudSdkGoERROR))
		return nil
	}
	v, err := jsonpath.Get("$.Data[*]", listResp)
	if err != nil {
		log.Printf("[WARN] StarRocks describeInstances returned no Data (empty region?): %v", err)
		return nil
	}
	instances, ok := v.([]interface{})
	if !ok || len(instances) == 0 {
		log.Printf("[INFO] StarRocks node_group walk: 0 instances in region %s", client.RegionId)
		return nil
	}
	for _, inst := range instances {
		instance := inst.(map[string]interface{})
		clusterId := fmt.Sprint(instance["InstanceId"])
		instanceName := fmt.Sprint(instance["InstanceName"])

		// describeNodeGroups needs ClusterId; omit nodeGroupIds to list every node group.
		ngAction := "/webapi/nodegroup/describeNodeGroups"
		ngQuery := map[string]*string{
			"ClusterId": StringPointer(clusterId),
			"RegionId":  StringPointer(client.RegionId),
		}
		ngBody := make(map[string]interface{})
		var ngResp map[string]interface{}
		ngWait := incrementalWait(3*time.Second, 5*time.Second)
		ngErr := retry.Retry(2*time.Minute, func() *retry.RetryError {
			resp, e := client.RoaPost("starrocks", "2022-10-19", ngAction, ngQuery, nil, ngBody, true)
			ngResp = resp
			if e != nil {
				if NeedRetry(e) {
					ngWait()
					return retry.RetryableError(e)
				}
				return retry.NonRetryableError(e)
			}
			return nil
		})
		addDebug(ngAction, ngResp, ngBody)
		if ngErr != nil {
			log.Printf("[ERROR] Failed to list StarRocks NodeGroups for cluster %s (%s): %s", instanceName, clusterId, ngErr)
			continue
		}
		ngV, ngPathErr := jsonpath.Get("$.Data[*]", ngResp)
		if ngPathErr != nil {
			log.Printf("[WARN] StarRocks describeNodeGroups returned no Data for cluster %s (%s): %v", instanceName, clusterId, ngPathErr)
			continue
		}
		nodeGroups, ok := ngV.([]interface{})
		if !ok || len(nodeGroups) == 0 {
			continue
		}
		for _, ngItem := range nodeGroups {
			ng := ngItem.(map[string]interface{})
			nodeGroupId := fmt.Sprint(ng["NodeGroupId"])
			// node_group names are not always present; fall back to the cluster name.
			ngName := instanceName
			if rawName, hasName := ng["NodeGroupName"]; hasName && fmt.Sprint(rawName) != "" {
				ngName = fmt.Sprint(rawName)
			}

			// NODELETING guard — always active (both NORMAL and SWEEP_ALL modes).
			if protectedSweepVpcName(ngName) {
				log.Printf("[INFO] Skipping StarRocks NodeGroup (NODELETING): %s (%s:%s)", ngName, clusterId, nodeGroupId)
				continue
			}
			// NORMAL mode: only clean resources carrying the tftestacc prefix.
			if !sweepAll() {
				skip := true
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(ngName), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
				if skip {
					log.Printf("[INFO] Skipping StarRocks NodeGroup (no prefix): %s (%s:%s)", ngName, clusterId, nodeGroupId)
					continue
				}
			}
			log.Printf("[INFO] Deleting StarRocks NodeGroup: %s (%s:%s)", ngName, clusterId, nodeGroupId)
			delAction := "/webapi/lifecycle/releaseNodeGroup"
			delQuery := map[string]*string{
				"InstanceId":  StringPointer(clusterId),
				"NodeGroupId": StringPointer(nodeGroupId),
				"RegionId":    StringPointer(client.RegionId),
			}
			delBody := make(map[string]interface{})
			if _, e := client.RoaPost("starrocks", "2022-10-19", delAction, delQuery, nil, delBody, true); e != nil {
				log.Printf("[ERROR] Failed to delete StarRocks NodeGroup (%s:%s): %s", clusterId, nodeGroupId, e)
				continue
			}
			log.Printf("[INFO] Delete StarRocks NodeGroup success: %s (%s:%s)", ngName, clusterId, nodeGroupId)
		}
	}
	return nil
}

// testSweepStarRocksNodeGroup >>> Sweeper.

// Test StarRocks NodeGroup. >>> Resource test cases, automatically generated.
// Case NodeGroup_副本1760322662795 11601
func TestAccAliCloudStarRocksNodeGroup_basic11601(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_star_rocks_node_group.default"
	ra := resourceAttrInit(resourceId, AlicloudStarRocksNodeGroupMap11601)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &StarRocksServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeStarRocksNodeGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudStarRocksNodeGroupBasicDependence11601)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":                 "test_desc",
					"node_group_name":             name,
					"instance_id":                 "${alicloud_star_rocks_instance.defaultvjnpM0.id}",
					"spec_type":                   "standard",
					"storage_performance_level":   "pl1",
					"pricing_cycle":               "1",
					"auto_renew":                  "false",
					"storage_size":                "200",
					"duration":                    "1",
					"pay_type":                    "postPaid",
					"cu":                          "8",
					"disk_number":                 "1",
					"resident_node_number":        "1",
					"local_storage_instance_type": "non_local_storage",
					"promotion_option_no":         "blank",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                 "test_desc",
						"node_group_name":             name,
						"instance_id":                 CHECKSET,
						"spec_type":                   "standard",
						"storage_performance_level":   "pl1",
						"pricing_cycle":               CHECKSET,
						"auto_renew":                  "false",
						"storage_size":                "200",
						"duration":                    "1",
						"pay_type":                    "postPaid",
						"cu":                          "8",
						"disk_number":                 "1",
						"resident_node_number":        "1",
						"local_storage_instance_type": "non_local_storage",
						"promotion_option_no":         "blank",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cu":                  "16",
					"promotion_option_no": "2345",
					"fast_mode":           "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cu":                  "16",
						"promotion_option_no": CHECKSET,
						"fast_mode":           "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_size": "500",
					"spec_type":    "ramEnhanced",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_size": "500",
						"spec_type":    "ramEnhanced",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_performance_level": "pl2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_performance_level": "pl2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_number": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_number": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resident_node_number": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resident_node_number": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "duration", "fast_mode", "pricing_cycle", "promotion_option_no"},
			},
		},
	})
}

var AlicloudStarRocksNodeGroupMap11601 = map[string]string{
	"status":        CHECKSET,
	"create_time":   CHECKSET,
	"node_group_id": CHECKSET,
	"region_id":     CHECKSET,
}

func AlicloudStarRocksNodeGroupBasicDependence11601(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultq6pcFe" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "test-vpc-511"
}

resource "alicloud_vswitch" "defaultujlpyG" {
  vpc_id       = alicloud_vpc.defaultq6pcFe.id
  zone_id      = "cn-hangzhou-i"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "sr-test-ng"
}

resource "alicloud_star_rocks_instance" "defaultvjnpM0" {
  cluster_zone_id = "cn-hangzhou-i"
  encrypted       = false
  auto_renew      = false
  pay_type        = "postPaid"
  frontend_node_groups {
    cu                        = "8"
    storage_size              = "100"
    storage_performance_level = "pl1"
    disk_number               = "1"
    zone_id                   = "cn-hangzhou-i"
    spec_type                 = "standard"
    resident_node_number      = "1"
  }
  instance_name = "t1"
  vswitches {
    zone_id    = "cn-hangzhou-i"
    vswitch_id = alicloud_vswitch.defaultujlpyG.id
  }
  vpc_id                  = alicloud_vpc.defaultq6pcFe.id
  version                 = "3.3"
  run_mode                = "shared_data"
  package_type            = "official"
  oss_accessing_role_name = "AliyunEMRStarRocksAccessingOSSRolecn"
  admin_password          = "1qaz@QAZ"
  backend_node_groups {
    cu                        = "8"
    storage_size              = "200"
    zone_id                   = "cn-hangzhou-i"
    spec_type                 = "standard"
    resident_node_number      = "3"
    disk_number               = "1"
    storage_performance_level = "pl1"
  }
}


`, name)
}

// Case NodeGroup 11074
func TestAccAliCloudStarRocksNodeGroup_basic11074(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_star_rocks_node_group.default"
	ra := resourceAttrInit(resourceId, AlicloudStarRocksNodeGroupMap11074)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &StarRocksServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeStarRocksNodeGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudStarRocksNodeGroupBasicDependence11074)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":                 "test_desc",
					"node_group_name":             name,
					"instance_id":                 "${alicloud_star_rocks_instance.defaultvjnpM0.id}",
					"spec_type":                   "standard",
					"storage_performance_level":   "pl1",
					"pricing_cycle":               "1",
					"auto_renew":                  "false",
					"storage_size":                "200",
					"duration":                    "1",
					"pay_type":                    "postPaid",
					"cu":                          "8",
					"disk_number":                 "1",
					"resident_node_number":        "1",
					"local_storage_instance_type": "non_local_storage",
					"promotion_option_no":         "blank",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                 "test_desc",
						"node_group_name":             name,
						"instance_id":                 CHECKSET,
						"spec_type":                   "standard",
						"storage_performance_level":   "pl1",
						"pricing_cycle":               CHECKSET,
						"auto_renew":                  "false",
						"storage_size":                "200",
						"duration":                    "1",
						"pay_type":                    "postPaid",
						"cu":                          "8",
						"disk_number":                 "1",
						"resident_node_number":        "1",
						"local_storage_instance_type": "non_local_storage",
						"promotion_option_no":         "blank",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cu":                  "16",
					"promotion_option_no": "2345",
					"fast_mode":           "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cu":                  "16",
						"promotion_option_no": CHECKSET,
						"fast_mode":           "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_size": "500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_size": "500",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_performance_level": "pl2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_performance_level": "pl2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_number": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_number": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resident_node_number": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resident_node_number": "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "duration", "fast_mode", "pricing_cycle", "promotion_option_no"},
			},
		},
	})
}

var AlicloudStarRocksNodeGroupMap11074 = map[string]string{
	"status":        CHECKSET,
	"create_time":   CHECKSET,
	"node_group_id": CHECKSET,
	"region_id":     CHECKSET,
}

func AlicloudStarRocksNodeGroupBasicDependence11074(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultq6pcFe" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "defaultujlpyG" {
  vpc_id       = alicloud_vpc.defaultq6pcFe.id
  zone_id      = "cn-hangzhou-i"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "sr-test-ng"
}

resource "alicloud_star_rocks_instance" "defaultvjnpM0" {
  cluster_zone_id = "cn-hangzhou-i"
  encrypted       = false
  auto_renew      = false
  pay_type        = "postPaid"
  frontend_node_groups {
    cu                        = "8"
    storage_size              = "100"
    storage_performance_level = "pl1"
    disk_number               = "1"
    zone_id                   = "cn-hangzhou-i"
    spec_type                 = "standard"
    resident_node_number      = "1"
  }
  instance_name = "t1"
  vswitches {
    zone_id    = "cn-hangzhou-i"
    vswitch_id = alicloud_vswitch.defaultujlpyG.id
  }
  vpc_id                  = alicloud_vpc.defaultq6pcFe.id
  version                 = "3.3"
  run_mode                = "shared_data"
  package_type            = "official"
  oss_accessing_role_name = "AliyunEMRStarRocksAccessingOSSRole"
  admin_password          = "1qaz@QAZ"
  backend_node_groups {
    cu                        = "8"
    storage_size              = "200"
    zone_id                   = "cn-hangzhou-i"
    spec_type                 = "standard"
    resident_node_number      = "3"
    disk_number               = "1"
    storage_performance_level = "pl1"
  }
}


`, name)
}

// Test StarRocks NodeGroup. <<< Resource test cases, automatically generated.
