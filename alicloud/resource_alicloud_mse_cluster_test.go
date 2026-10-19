package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers("alicloud_mse_cluster", &resource.Sweeper{
		Name:         "alicloud_mse_cluster",
		Dependencies: []string{"alicloud_mse_engine_namespace", "alicloud_mse_nacos_config", "alicloud_mse_znode"},
		F:            testSweepMSECluster,
	})
}

// listMseSweepClusters lists every MSE cluster in the region. It is shared by
// the child-resource sweepers (engine_namespace, nacos_config, znode) so they
// can enumerate the clusters they need to walk. List/parse failures are logged
// and treated as an empty result (sweepers must remain non-fatal); callers may
// still guard the returned error defensively.
func listMseSweepClusters(client *connectivity.AliyunClient) ([]map[string]interface{}, error) {
	action := "ListClusters"
	request := map[string]interface{}{
		"PageSize": PageSizeLarge,
		"PageNum":  1,
	}
	var clusters []map[string]interface{}
	var response map[string]interface{}
	var err error
	for {
		response, err = client.RpcGet("mse", "2019-05-31", action, request, nil)
		if err != nil {
			log.Printf("[ERROR] Failed to list mse clusters: %s", err)
			return nil, nil
		}
		resp, err := jsonpath.Get("$.Data", response)
		if err != nil {
			log.Printf("[ERROR] Failed to parse $.Data (treat as empty): %s", err)
			return nil, nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			clusters = append(clusters, v.(map[string]interface{}))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNum"] = request["PageNum"].(int) + 1
	}
	return clusters, nil
}

// mseSweepPrefixes is the name-prefix allow-list used by the MSE child sweepers
// when ALICLOUD_SWEEP_ALL_RESOURCES is not set. When SWEEP_ALL is set the prefix
// filter is bypassed and only protectedSweepVpcName (NODELETING markers) guards
// deletion, matching the canonical sweeper contract.
var mseSweepPrefixes = []string{"tf-testAcc", "tf-test", "tf_testAcc"}

// mseSweepNameAllowed reports whether a resource name should be processed by an
// MSE child sweeper. It applies the NODELETING name-protection guard first
// (always, even under SWEEP_ALL) and then, only when SWEEP_ALL is off, the
// test-prefix allow-list.
func mseSweepNameAllowed(name string) bool {
	if protectedSweepVpcName(name) {
		return false
	}
	if sweepAll() {
		return true
	}
	for _, prefix := range mseSweepPrefixes {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			return true
		}
	}
	return false
}

// testSweepMSECluster sweeps leftover MSE clusters in the region. The cluster
// is the expensive parent resource for engine_namespace / nacos_config / znode,
// so those child sweepers are registered as Dependencies and run first. Each
// cluster is listed via ListClusters (the same API the provider Read uses),
// guarded by mseSweepNameAllowed (NODELETING name-protection always; tf-test
// prefix when ALICLOUD_SWEEP_ALL_RESOURCES is off; bypass-prefix under
// SWEEP_ALL). Deletion uses DeleteCluster mirroring resourceAlicloudMseCluster
// Delete. Per-cluster errors are logged and skipped; the sweeper always returns
// nil so siblings still run.
func testSweepMSECluster(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	clusters, err := listMseSweepClusters(client)
	if err != nil {
		log.Printf("[ERROR] ListClusters for mse_cluster sweep failed: %s", err)
		return nil
	}

	for _, cluster := range clusters {
		name := fmt.Sprint(cluster["ClusterAliasName"])
		instanceId := fmt.Sprint(cluster["InstanceId"])
		if !mseSweepNameAllowed(name) {
			log.Printf("[INFO] Skipping mse cluster: %s (%s)", name, instanceId)
			continue
		}
		deleteAction := "DeleteCluster"
		deleteRequest := map[string]interface{}{
			"InstanceId": instanceId,
		}
		if _, derr := client.RpcPost("mse", "2019-05-31", deleteAction, nil, deleteRequest, false); derr != nil {
			log.Printf("[ERROR] Failed to delete mse cluster %s (%s): %s", name, instanceId, derr)
			continue
		}
		// DeleteCluster is asynchronous; wait briefly so the next sweep pass
		// does not race a still-destroying cluster's child APIs.
		time.Sleep(30 * time.Second)
		log.Printf("[INFO] Delete mse cluster success: %s (%s)", name, instanceId)
	}
	return nil
}

func TestAccAliCloudMSECluster_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mse_cluster.default"
	ra := resourceAttrInit(resourceId, MseClusterMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMseCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccMseCluster%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, MseClusterBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shenzhen"})
		},

		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_specification":     "MSE_SC_1_2_60_c",
					"cluster_type":              "Nacos-Ans",
					"cluster_version":           "NACOS_2_0_0",
					"instance_count":            "1",
					"net_type":                  "privatenet",
					"vswitch_id":                "${alicloud_vswitch.default.id}",
					"pub_network_flow":          "1",
					"cluster_alias_name":        name,
					"connection_type":           "slb",
					"mse_version":               "mse_dev",
					"private_slb_specification": "test",
					"pub_slb_specification":     "test",
					"disk_type":                 "test",
					"request_pars":              "",
					"vpc_id":                    "${alicloud_vpc.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_specification": "MSE_SC_1_2_60_c",
						"cluster_type":          "Nacos-Ans",
						"cluster_version":       "NACOS_2_0_0",
						"instance_count":        "1",
						"net_type":              "privatenet",
						"vswitch_id":            CHECKSET,
						"pub_network_flow":      "1",
						"cluster_alias_name":    name,
						"connection_type":       "slb",
						"mse_version":           "mse_dev",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_entry_list": []string{"127.0.0.1/32"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_entry_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_alias_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_alias_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_alias_name": name,
					"acl_entry_list":     []string{"127.0.0.0/16"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_alias_name": name,
						"acl_entry_list.#":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_specification": "MSE_SC_2_4_60_c",
					"resource_group_id":     "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_specification": "MSE_SC_2_4_60_c",
						"resource_group_id":     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_slb_specification": "test_update",
					"pub_slb_specification":     "test_update",
					"disk_type":                 "test_update",
					"request_pars":              "{}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_alias_name": name,
						"vswitch_id":         CHECKSET,
						"vpc_id":             CHECKSET,
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl_entry_list", "private_slb_specification", "pub_slb_specification", "disk_type", "request_pars"},
			},
		},
	})
}

func TestAccAliCloudMSECluster_changeNetwork(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mse_cluster.default"
	ra := resourceAttrInit(resourceId, MseClusterMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMseCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccMseCluster%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, MseClusterBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shenzhen"})
		},

		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_specification": "MSE_SC_1_2_60_c",
					"cluster_type":          "Nacos-Ans",
					"cluster_version":       "NACOS_2_0_0",
					"instance_count":        "1",
					"net_type":              "privatenet",
					"vswitch_id":            "${alicloud_vswitch.default.id}",
					"pub_network_flow":      "1",
					"cluster_alias_name":    name,
					"connection_type":       "slb",
					"mse_version":           "mse_dev",
					"vpc_id":                "${alicloud_vpc.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_specification": "MSE_SC_1_2_60_c",
						"cluster_type":          "Nacos-Ans",
						"cluster_version":       "NACOS_2_0_0",
						"instance_count":        "1",
						"net_type":              "privatenet",
						"vswitch_id":            CHECKSET,
						"pub_network_flow":      "1",
						"cluster_alias_name":    name,
						"connection_type":       "slb",
						"mse_version":           "mse_dev",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":           "${alicloud_vpc.default2.id}",
					"vswitch_id":       "${alicloud_vswitch.default2.id}",
					"pub_network_flow": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_alias_name": name,
						"vswitch_id":         CHECKSET,
						"vpc_id":             CHECKSET,
						"pub_network_flow":   "2",
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

func TestAccAliCloudMSECluster_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mse_cluster.default"
	ra := resourceAttrInit(resourceId, MseClusterMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMseCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccMseCluster%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, MseClusterBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},

		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_specification": "MSE_SC_1_2_60_c",
					"cluster_type":          "ZooKeeper",
					"cluster_version":       "ZooKeeper_3_8_0",
					"instance_count":        "1",
					"net_type":              "privatenet",
					"vswitch_id":            "${alicloud_vswitch.default.id}",
					"pub_network_flow":      "1",
					"cluster_alias_name":    name,
					"connection_type":       "slb",
					"mse_version":           "mse_dev",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_specification": "MSE_SC_1_2_60_c",
						"cluster_type":          "ZooKeeper",
						"cluster_version":       "ZooKeeper_3_8_0",
						"instance_count":        "1",
						"net_type":              "privatenet",
						"vswitch_id":            CHECKSET,
						"pub_network_flow":      "1",
						"cluster_alias_name":    name,
						"connection_type":       "slb",
						"mse_version":           "mse_dev",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_entry_list": []string{"127.0.0.1/32"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_entry_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_alias_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_alias_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_alias_name": name,
					"acl_entry_list":     []string{"127.0.0.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_alias_name": name,
						"acl_entry_list.#":   "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl_entry_list"},
			},
		},
	})
}

func TestAccAliCloudMSECluster_pro(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mse_cluster.default"
	ra := resourceAttrInit(resourceId, MseClusterMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMseCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccMseCluster%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, MseClusterBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_specification": "MSE_SC_4_8_60_c",
					"cluster_type":          "Nacos-Ans",
					"cluster_version":       "NACOS_2_0_0",
					"instance_count":        "3",
					"net_type":              "privatenet",
					"vswitch_id":            "${alicloud_vswitch.default.id}",
					"pub_network_flow":      "1",
					"cluster_alias_name":    name,
					"mse_version":           "mse_pro",
					"connection_type":       "slb",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_specification": "MSE_SC_4_8_60_c",
						"cluster_type":          "Nacos-Ans",
						"cluster_version":       "NACOS_2_0_0",
						"instance_count":        "3",
						"vswitch_id":            CHECKSET,
						"pub_network_flow":      "1",
						"cluster_alias_name":    name,
						"mse_version":           "mse_pro",
						"connection_type":       "slb",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_specification": "MSE_SC_8_16_60_c",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_specification": "MSE_SC_8_16_60_c",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_count": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_count": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"version_code": "LATEST",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version_code": CHECKSET,
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

func TestAccAliCloudMSECluster_serverless(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mse_cluster.default"
	ra := resourceAttrInit(resourceId, MseClusterMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMseCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccMseCluster%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, MseClusterBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_specification": "MSE_SC_SERVERLESS",
					"cluster_type":          "Nacos-Ans",
					"cluster_version":       "NACOS_2_0_0",
					"instance_count":        "3",
					"net_type":              "both",
					"vswitch_id":            "${alicloud_vswitch.default.id}",
					"pub_network_flow":      "1",
					"cluster_alias_name":    name,
					"mse_version":           "mse_serverless",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_specification": "MSE_SC_SERVERLESS",
						"cluster_type":          "Nacos-Ans",
						"cluster_version":       "NACOS_2_0_0",
						"instance_count":        "3",
						"vswitch_id":            CHECKSET,
						"pub_network_flow":      "1",
						"cluster_alias_name":    name,
						"mse_version":           "mse_serverless",
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

func TestAccAliCloudMSECluster_VpcId(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mse_cluster.default"
	ra := resourceAttrInit(resourceId, MseClusterMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMseCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccMseCluster%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, MseClusterBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_specification": "MSE_SC_1_2_60_c",
					"cluster_type":          "Nacos-Ans",
					"cluster_version":       "NACOS_2_0_0",
					"instance_count":        "1",
					"net_type":              "privatenet",
					"vswitch_id":            "${alicloud_vswitch.default.id}",
					"pub_network_flow":      "1",
					"cluster_alias_name":    name,
					"connection_type":       "slb",
					"mse_version":           "mse_dev",
					"vpc_id":                "${alicloud_vpc.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_specification": "MSE_SC_1_2_60_c",
						"cluster_type":          "Nacos-Ans",
						"cluster_version":       "NACOS_2_0_0",
						"instance_count":        "1",
						"net_type":              "privatenet",
						"vswitch_id":            CHECKSET,
						"pub_network_flow":      "1",
						"cluster_alias_name":    name,
						"connection_type":       "slb",
						"mse_version":           "mse_dev",
						"vpc_id":                CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

func TestAccAliCloudMSECluster_PrePaid(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mse_cluster.default"
	ra := resourceAttrInit(resourceId, MseClusterMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMseCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccMseCluster%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, MseClusterBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckPrePaidResources(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheckWithTime(t, []int{22})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_specification": "MSE_SC_1_2_60_c",
					"cluster_type":          "Nacos-Ans",
					"cluster_version":       "NACOS_2_0_0",
					"instance_count":        "1",
					"net_type":              "privatenet",
					"payment_type":          "Subscription",
					"vswitch_id":            "${alicloud_vswitch.default.id}",
					"pub_network_flow":      "1",
					"cluster_alias_name":    name,
					"connection_type":       "slb",
					"mse_version":           "mse_dev",
					"vpc_id":                "${alicloud_vpc.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_specification": "MSE_SC_1_2_60_c",
						"cluster_type":          "Nacos-Ans",
						"cluster_version":       "NACOS_2_0_0",
						"instance_count":        "1",
						"net_type":              "privatenet",
						"vswitch_id":            CHECKSET,
						"pub_network_flow":      "1",
						"cluster_alias_name":    name,
						"connection_type":       "slb",
						"mse_version":           "mse_dev",
						"vpc_id":                CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

func TestAccAliCloudMSECluster_single_eni(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mse_cluster.default"
	ra := resourceAttrInit(resourceId, MseClusterMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMseCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccMseCluster%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, MseClusterBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"ap-southeast-1"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_specification": "MSE_SC_1_2_60_c",
					"cluster_type":          "Nacos-Ans",
					"cluster_version":       "NACOS_2_0_0",
					"instance_count":        "1",
					"net_type":              "privatenet",
					"vswitch_id":            "${alicloud_vswitch.default.id}",
					"pub_network_flow":      "1",
					"cluster_alias_name":    name,
					"connection_type":       "single_eni",
					"mse_version":           "mse_dev",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_specification": "MSE_SC_1_2_60_c",
						"cluster_type":          "Nacos-Ans",
						"cluster_version":       "NACOS_2_0_0",
						"instance_count":        "1",
						"net_type":              "privatenet",
						"vswitch_id":            CHECKSET,
						"pub_network_flow":      "1",
						"cluster_alias_name":    name,
						"connection_type":       "single_eni",
						"mse_version":           "mse_dev",
					}),
				),
			},
		},
	})
}

var MseClusterMap = map[string]string{
	"cluster_id":  CHECKSET,
	"app_version": CHECKSET,
	"status":      CHECKSET,
}

func MseClusterBasicdependence(name string) string {
	return fmt.Sprintf(`

	variable "name" {
		 default = "%v"
		}

	data "alicloud_zones" "default" {
	  available_resource_creation = "VSwitch"
	}

	resource "alicloud_vpc" "default" {
	  vpc_name       = "${var.name}"
      cidr_block = "172.17.3.0/24"
	}

	resource "alicloud_vswitch" "default" {
	  vswitch_name = "${var.name}"
	  cidr_block   = "172.17.3.0/24"
	  vpc_id       = alicloud_vpc.default.id
	  zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_vpc" "default2" {
	  vpc_name       = "tftest"
      cidr_block = "172.17.3.0/24"
	}

	resource "alicloud_vswitch" "default2" {
	  vswitch_name = "tftest"
	  cidr_block   = "172.17.3.0/24"
	  vpc_id       = alicloud_vpc.default2.id
	  zone_id      = data.alicloud_zones.default.zones.0.id
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}
`, name)
}

// lintignore: R001
func TestUnitAlicloudMSECluster(t *testing.T) {
	p := Provider().ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_mse_cluster"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_mse_cluster"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"cluster_specification": "CreateClusterValue",
		"cluster_type":          "CreateClusterValue",
		"cluster_version":       "CreateClusterValue",
		"instance_count":        3,
		"net_type":              "CreateClusterValue",
		"vswitch_id":            "CreateClusterValue",
		"pub_network_flow":      "1",
		"cluster_alias_name":    "CreateClusterValue",
		"mse_version":           "CreateClusterValue",
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
		if err != nil {
			log.Printf("[ERROR] the field %s setting error", key)
		}
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// QueryClusterDetail
		"Data": map[string]interface{}{
			"ClusterType":         "CreateClusterValue",
			"InstanceCount":       3,
			"PubNetworkFlow":      "1",
			"InitStatus":          "INIT_SUCCESS",
			"ClusterId":           "CreateClusterValue",
			"MseVersion":          "CreateClusterValue",
			"NetType":             "CreateClusterValue",
			"VSwitchId":           "CreateClusterValue",
			"OrderClusterVersion": "CreateClusterValue",
			"ClusterAliasName":    "CreateClusterValue",
		},
		"VSwitchId":  "CreateClusterValue",
		"VpcId":      "CreateClusterValue",
		"InstanceId": "CreateClusterValue",
		"Success":    "true",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateCluster
		"VSwitchId":  "CreateClusterValue",
		"VpcId":      "CreateClusterValue",
		"InstanceId": "CreateClusterValue",
		"Success":    "true",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_mse_cluster", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewMseClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudMseClusterCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// QueryClusterDetail Response
		"VSwitchId":  "CreateClusterValue",
		"VpcId":      "CreateClusterValue",
		"InstanceId": "CreateClusterValue",
		"Success":    "true",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateCluster" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudMseClusterCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_mse_cluster"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewMseClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudMseClusterUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateAcl
	attributesDiff := map[string]interface{}{
		"acl_entry_list": []string{"UpdateAclValue"},
	}
	diff, err := newInstanceDiff("alicloud_mse_cluster", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_mse_cluster"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// QueryClusterDetail
		"Success": "true",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateAcl" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudMseClusterUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_mse_cluster"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// UpdateCluster
	attributesDiff = map[string]interface{}{
		"cluster_alias_name": "UpdateCluster",
	}
	diff, err = newInstanceDiff("alicloud_mse_cluster", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_mse_cluster"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// QueryClusterDetail
		"Data": map[string]interface{}{
			"ClusterAliasName": "UpdateCluster",
		},
		"Success": "true",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateCluster" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudMseClusterUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_mse_cluster"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "QueryClusterDetail" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudMseClusterRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewMseClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudMseClusterDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteCluster" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{
							"Data": map[string]interface{}{
								"InitStatus": "DESTROY_SUCCESS",
							},
							"Success": "true",
						}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudMseClusterDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
