package alicloud

// collection-order: bypass secondary_cidr_blocks -- pre-existing schema and order coverage outside this sweeper-only change
// collection-order: bypass user_cidrs -- pre-existing schema and order coverage outside this sweeper-only change

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	roacs "github.com/alibabacloud-go/cs-20151215/v8/client"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	vpcSweepUtilV2 "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers("alicloud_vpc", &resource.Sweeper{
		Name: "alicloud_vpc",
		F:    testSweepVpcs,
		// When implemented, these should be removed firstly
		Dependencies: vpcSweepDependencies([]string{
			"alicloud_vswitch",
			"alicloud_nat_gateway",
			"alicloud_security_group",
			"alicloud_ots_instance",
			"alicloud_router_interface",
			"alicloud_route_table",
			"alicloud_cen_instance",
			"alicloud_edas_cluster",
			"alicloud_edas_k8s_cluster",
			"alicloud_network_acl",
			"alicloud_cs_kubernetes",
			"alicloud_instance",
		}),
	})
}

func testSweepVpcs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	if sweepAllVpcs() {
		rpc := throttledVpcSweepRPC(nativeVpcSweepRPC(client))
		cenRPCForVpc := newVpcCenInventory(client.RegionId, rpc)
		availability := nativeVpcServiceAvailability(client)
		return sweepVpcsWithCallbacks(client.RegionId, rpc, func(id string) error {
			return sweepSelectedVpcWithAvailability(client, id, availability, cenRPCForVpc(id))
		})
	}
	return sweepVpcsWithCallbacks(client.RegionId, client.RpcPost, func(id string) error { return sweepSelectedVpc(client, id) })
}

// Inject the list transport and selected-VPC callback so selection/protection
// can be verified without credentials, clients, network calls or monkey patches.
func sweepVpcsWithCallbacks(region string, rpc vpcSweepRPC, deleteSelected func(string) error) error {
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"testAcc",
		"tf-example",
		"terraform-example",
		"tf-test",
	}

	vpcIds := make([]string, 0)
	action := "DescribeVpcs"
	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
		"RegionId":   region,
	}
	for {
		response, err := rpc("Vpc", "2016-04-28", action, nil, request, true)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve VPC in service list: %s", err)
			return err
		}
		resp, err := jsonpath.Get("$.Vpcs.Vpc", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Vpcs.Vpc", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if skip, reason := skipVpcSweep(item, prefixes); skip {
				log.Printf("[INFO] Skipping VPC %v (%v): %s", item["VpcName"], item["VpcId"], reason)
				continue
			}
			id, ok := item["VpcId"].(string)
			if !ok || id == "" {
				return fmt.Errorf("DescribeVpcs returned selected VPC without valid VpcId; refusing deletion")
			}
			vpcIds = append(vpcIds, id)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	workers := 1
	if sweepAllVpcs() {
		workers = 4
	}
	return runVpcSweeps(vpcIds, workers, deleteSelected)
}

func sweepSelectedVpc(client *connectivity.AliyunClient, id string, scopedRPC ...vpcSweepRPC) error {
	return sweepSelectedVpcWithAvailability(client, id, nil, scopedRPC...)
}

func sweepSelectedVpcWithAvailability(client *connectivity.AliyunClient, id string, availability func(string) (bool, error), scopedRPC ...vpcSweepRPC) error {
	log.Printf("[INFO] Deleting VPC: (%s)", id)
	if sweepAllVpcs() {
		rpc := throttledVpcSweepRPC(nativeVpcSweepRPC(client))
		if len(scopedRPC) > 0 {
			rpc = scopedRPC[0]
		}
		ackAPI, ackOwnership := trackVpcACKOwnership(client.RegionId, id, nativeVpcACKAPI(client, client.RegionId))
		ackNode := vpcACKSweeper(client.RegionId, id, ackAPI)
		flinkNode := flinkVpcOwnerSweeper(client.RegionId, id, throttledVpcSweepRPC(nativeVpcFlinkRPC(client)))
		dhcpNode := vpcDHCPAttachmentSweeper(client.RegionId, id, rpc)
		hologramNode := vpcHologramSweeper(client.RegionId, id, nativeVpcHologramAPI(client))
		asmAPI := nativeVpcASMAPI(client, client.RegionId)
		asmAPI.clusterWasScoped = ackOwnership
		asmNode := vpcASMSweeper(client.RegionId, id, asmAPI)
		asmNode.Dependencies = []string{ackNode.Name}
		starrocksNode := vpcStarrocksSweeper(client.RegionId, id, nativeVpcStarrocksAPI(client))
		otsNode := nativeVpcOtsAttachmentSweeper(client.RegionId, id, client)
		graph := vpcScopedSweepers(client.RegionId, id, rpc, map[string]*resource.Sweeper{
			ackNode.Name: ackNode, flinkNode.Name: flinkNode, dhcpNode.Name: dhcpNode,
			hologramNode.Name: hologramNode, asmNode.Name: asmNode, starrocksNode.Name: starrocksNode, otsNode.Name: otsNode,
		}, serviceVpcOwnerSweepers(client.RegionId, id, rpc))
		if availability != nil {
			applyVpcServiceAvailability(graph, availability)
		}
		return runVpcScopedSweepers(client.RegionId, id, graph)
	}
	sweepVpcDependentResources(client, id)
	return deleteVpcScopedResource(client.RpcPost, client.RegionId, "Vpc", "2016-04-28", "DeleteVpc", "VpcId", id, nil)
}

// listDependentResources queries resources filtered by a parameter (typically
// VpcId, or ResourceId for FlowLogs) and returns the raw item maps. It handles
// pagination and returns an empty slice when no resources are found.
func listDependentResources(client *connectivity.AliyunClient, product, version, action, responsePath, paramKey, paramValue string) []map[string]interface{} {
	items := make([]map[string]interface{}, 0)
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		paramKey:     paramValue,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	for {
		response, err := client.RpcPost(product, version, action, nil, request, true)
		if err != nil {
			log.Printf("[ERROR] Failed to list %s for %s=%s: %v", action, paramKey, paramValue, err)
			return items
		}
		resp, err := jsonpath.Get(responsePath, response)
		if err != nil {
			return items
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				if include, reason := sweepDependentResourceScope(item, paramKey, paramValue); !include {
					log.Printf("[ERROR] Refusing %s resource outside verified %s=%s scope: %s", action, paramKey, paramValue, reason)
					continue
				}
				items = append(items, item)
			} else if sweepAllVpcs() {
				log.Printf("[ERROR] Refusing malformed %s resource with unknown VPC ownership", action)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return items
}

// deleteDependentResource deletes a single resource by ID using RpcPost with retry.
func deleteDependentResource(client *connectivity.AliyunClient, product, version, action, idField, id string) {
	log.Printf("[INFO] Deleting %s: %s", action, id)
	request := map[string]interface{}{
		idField:    id,
		"RegionId": client.RegionId,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, err := client.RpcPost(product, version, action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[ERROR] Failed to delete %s (%s): %v", action, id, err)
	}
}

// sweepVpcDependentResources proactively queries and deletes all dependent
// resources within the specified VPC before the VPC itself can be deleted.
// This ensures VPCs with leftover dependencies are properly cleaned up,
// regardless of whether the dependency sweepers caught them by name prefix.
// The deletion order follows Alibaba Cloud VPC dependency constraints:
// VSwitches, NAT Gateways, Security Groups, custom Route Tables, Network ACLs,
// VPN Gateways, HaVips, Flow Logs, and IPv6 Gateways.
func sweepVpcDependentResources(client *connectivity.AliyunClient, vpcId string) {
	vpcService := VpcService{client}
	ecsService := EcsService{client}

	// Ordinary prefix/age cleanup retains its legacy nine-resource order.
	// 1. Delete VSwitches in the VPC
	for _, item := range listDependentResources(client, "Vpc", "2016-04-28", "DescribeVSwitches", "$.VSwitches.VSwitch", "VpcId", vpcId) {
		if id := fmt.Sprint(item["VSwitchId"]); id != "" {
			if err := vpcService.sweepVSwitch(id); err != nil {
				log.Printf("[ERROR] Failed to delete VSwitch %s in VPC %s: %v", id, vpcId, err)
			}
		}
	}

	// 2. Delete NAT Gateways in the VPC (Force=true via sweepNatGateway)
	for _, item := range listDependentResources(client, "Vpc", "2016-04-28", "DescribeNatGateways", "$.NatGateways.NatGateway", "VpcId", vpcId) {
		if id := fmt.Sprint(item["NatGatewayId"]); id != "" {
			if err := vpcService.sweepNatGateway(id); err != nil {
				log.Printf("[ERROR] Failed to delete NAT Gateway %s in VPC %s: %v", id, vpcId, err)
			}
		}
	}

	// 3. Delete Security Groups in the VPC (ECS API)
	for _, item := range listDependentResources(client, "Ecs", "2014-05-26", "DescribeSecurityGroups", "$.SecurityGroups.SecurityGroup", "VpcId", vpcId) {
		if id := fmt.Sprint(item["SecurityGroupId"]); id != "" {
			if err := ecsService.sweepSecurityGroup(id); err != nil {
				log.Printf("[ERROR] Failed to delete Security Group %s in VPC %s: %v", id, vpcId, err)
			}
		}
	}

	// 4. Delete custom Route Tables in the VPC (skip System route tables)
	for _, item := range listDependentResources(client, "Vpc", "2016-04-28", "DescribeRouteTableList", "$.RouterTableList.RouterTableListType", "VpcId", vpcId) {
		if fmt.Sprint(item["RouteTableType"]) == "System" {
			continue
		}
		if id := fmt.Sprint(item["RouteTableId"]); id != "" {
			deleteDependentResource(client, "Vpc", "2016-04-28", "DeleteRouteTable", "RouteTableId", id)
		}
	}

	// 5. Delete Network ACLs in the VPC
	for _, item := range listDependentResources(client, "Vpc", "2016-04-28", "DescribeNetworkAcls", "$.NetworkAcls.NetworkAcl", "VpcId", vpcId) {
		if id := fmt.Sprint(item["NetworkAclId"]); id != "" {
			deleteDependentResource(client, "Vpc", "2016-04-28", "DeleteNetworkAcl", "NetworkAclId", id)
		}
	}

	// 6. Delete VPN Gateways in the VPC
	for _, item := range listDependentResources(client, "Vpc", "2016-04-28", "DescribeVpnGateways", "$.VpnGateways.VpnGateway", "VpcId", vpcId) {
		if id := fmt.Sprint(item["VpnGatewayId"]); id != "" {
			deleteDependentResource(client, "Vpc", "2016-04-28", "DeleteVpnGateway", "VpnGatewayId", id)
		}
	}

	// 7. Delete HaVips in the VPC
	for _, item := range listDependentResources(client, "Vpc", "2016-04-28", "DescribeHaVips", "$.HaVips.HaVip", "VpcId", vpcId) {
		if id := fmt.Sprint(item["HaVipId"]); id != "" {
			deleteDependentResource(client, "Vpc", "2016-04-28", "DeleteHaVip", "HaVipId", id)
		}
	}

	// 8. Delete Flow Logs associated with the VPC (ResourceId = VpcId)
	for _, item := range listDependentResources(client, "Vpc", "2016-04-28", "DescribeFlowLogs", "$.FlowLogs.FlowLog", "ResourceId", vpcId) {
		if id := fmt.Sprint(item["FlowLogId"]); id != "" {
			deleteDependentResource(client, "Vpc", "2016-04-28", "DeleteFlowLog", "FlowLogId", id)
		}
	}

	// 9. Delete IPv6 Gateways in the VPC
	for _, item := range listDependentResources(client, "Vpc", "2016-04-28", "DescribeIpv6Gateways", "$.Ipv6Gateways.Ipv6Gateway", "VpcId", vpcId) {
		if id := fmt.Sprint(item["Ipv6GatewayId"]); id != "" {
			deleteDependentResource(client, "Vpc", "2016-04-28", "DeleteIpv6Gateway", "Ipv6GatewayId", id)
		}
	}
}

func TestAccAliCloudVPC_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpc")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sVpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VpcIpv6SupportRegions)
		},

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"user_cidrs": []string{"106.11.62.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_cidrs.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "enable_ipv6"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cidr_block": "172.16.0.0/16",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr_block": "172.16.0.0/16",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secondary_cidr_blocks": []string{"10.0.0.0/8"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secondary_cidr_blocks.#": "1",
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
					"cidr_block":            "172.16.0.0/12",
					"vpc_name":              name + "update",
					"description":           name + "update",
					"secondary_cidr_blocks": []string{"10.0.0.0/8"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr_block":              "172.16.0.0/12",
						"vpc_name":                name + "update",
						"description":             name + "update",
						"secondary_cidr_blocks.#": "1",
					}),
				),
			},
		},
	})
}

var AlicloudVpcMap = map[string]string{
	"status":          CHECKSET,
	"router_id":       CHECKSET,
	"router_table_id": CHECKSET,
	"route_table_id":  CHECKSET,
	"ipv6_cidr_block": "",
}

func AlicloudVpcBasicDependence(name string) string {
	return fmt.Sprintf(`

data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}
`)
}

func TestAccAliCloudVPC_enableIpv6(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpc")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sVpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VpcIpv6SupportRegions)
		},

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"user_cidrs":  []string{"106.11.62.0/24"},
					"enable_ipv6": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_cidrs.#": "1",
						"enable_ipv6":  "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "enable_ipv6"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ipv6": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ipv6":     "false",
						"ipv6_cidr_block": "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ipv6": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ipv6":     "true",
						"ipv6_cidr_block": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secondary_cidr_blocks": []string{"10.0.0.0/8"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secondary_cidr_blocks.#": "1",
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
					"vpc_name":              name + "update",
					"description":           name + "update",
					"secondary_cidr_blocks": []string{"10.0.0.0/8"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_name":                name + "update",
						"description":             name + "update",
						"secondary_cidr_blocks.#": "1",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudVPC_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpc")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sVpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VpcIpv6SupportRegions)
		},

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ipv6":                    "true",
					"vpc_name":                       name,
					"description":                    name,
					"resource_group_id":              "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"dry_run":                        "false",
					"system_route_table_name":        name,
					"system_route_table_description": name,
					"user_cidrs":                     []string{"106.11.62.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ipv6":                    "true",
						"vpc_name":                       name,
						"description":                    name,
						"system_route_table_name":        name,
						"system_route_table_description": name,
						"resource_group_id":              CHECKSET,
						"dry_run":                        "false",
						"user_cidrs.#":                   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_route_table_name":        name + "_update",
					"system_route_table_description": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_route_table_name":        name + "_update",
						"system_route_table_description": name + "_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "enable_ipv6"},
			},
		},
	})
}

func TestAccAliCloudVPC_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpc")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sVpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VpcIpv6SupportRegions)
		},

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ipv6":       "true",
					"name":              name,
					"description":       name,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"dry_run":           "false",
					"user_cidrs":        []string{"106.11.62.0/24"},
					// "ipv6_cidr_block":   "240b:4000:95:5a00::/56",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ipv6":       "true",
						"name":              name,
						"description":       name,
						"resource_group_id": CHECKSET,
						"dry_run":           "false",
						"user_cidrs.#":      "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "-update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "enable_ipv6"},
			},
		},
	})
}

func TestAccAliCloudVPC_isNotDefault(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpc")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sVpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VpcIpv6SupportRegions)
		},

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ipv6":       "true",
					"is_default":        "false",
					"name":              name,
					"description":       name,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"dry_run":           "false",
					"user_cidrs":        []string{"106.11.62.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ipv6":       "true",
						"name":              name,
						"description":       name,
						"resource_group_id": CHECKSET,
						"dry_run":           "false",
						"user_cidrs.#":      "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "-update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "enable_ipv6"},
			},
		},
	})
}

func TestAccAliCloudVPC_isDefault(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpc")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sVpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-chengdu"})
		},

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ipv6":       "true",
					"is_default":        "true",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"dry_run":           "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ipv6":       "true",
						"resource_group_id": CHECKSET,
						"dry_run":           "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "-update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "enable_ipv6"},
			},
		},
	})
}

var AlicloudVpcMap1 = map[string]string{
	"status":          CHECKSET,
	"router_id":       CHECKSET,
	"router_table_id": CHECKSET,
	"route_table_id":  CHECKSET,
	"ipv6_cidr_block": CHECKSET,
}

func AlicloudVpcBasicDependence1(name string) string {
	return fmt.Sprintf(`

data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}
`)
}

// lintignore: R001
func TestUnitAlicloudVPCdsafa(t *testing.T) {
	p := Provider().ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_vpc"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_vpc"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"cidr_block":        "cidr_block",
		"description":       "description",
		"dry_run":           false,
		"enable_ipv6":       false,
		"resource_group_id": "resource_group_id",
		"vpc_name":          "vpc_name",
		"name":              "name",
		"user_cidrs":        []interface{}{"user_cidrs_1", "user_cidrs_2"},
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"Vpcs": map[string]interface{}{
			"Vpc": []interface{}{
				map[string]interface{}{
					"VpcId": "MockId",
					"UserCidrs": map[string]interface{}{
						"UserCidr": "UserCidr",
					},
					"CidrBlock":     "cidr_block",
					"Description":   "description",
					"Ipv6CidrBlock": "ipv6_cidr_block",
					"VRouterId":     "v_router_id",
					"SecondaryCidrBlocks": map[string]interface{}{
						"SecondaryCidrBlock": "secondary_cidr_blocks",
					},
					"Status": "Available",
					"Tags": map[string]interface{}{
						"key": "value",
					},
					"UserCidr": "user_cidrs",
					"VpcName":  "vpc_name",
				},
			},
		},
		//DescribeRouteTableList
		"Code": "200",
		"RouterTableList": map[string]interface{}{
			"RouterTableListType": []interface{}{
				map[string]interface{}{
					"RouteTableType": "System",
				},
			},
		},
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			result["VpcId"] = "MockId"
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"DeleteNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"ReadNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudVpcVpcCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["Normal"]("")
		})
		err := resourceAliCloudVpcVpcCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAliCloudVpcVpcCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := resourceAliCloudVpcVpcUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateMoveResourceGroupAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"resource_group_id"} {
			diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: "OldValue", New: "NewValue"}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_vpc"].Schema).Data(nil, diff)
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["Normal"]("")
		})
		err := resourceAliCloudVpcVpcUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateMoveResourceGroupNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"resource_group_id"} {
			diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: "", New: "NewValue"}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_vpc"].Schema).Data(nil, diff)
		resourceData1.SetId("MockId")
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudVpcVpcUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("UpdateModifyVpcAttributeAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"cidr_block", "description", "vpc_name", "name"} {
			diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: "OldValue", New: "NewValue"}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_vpc"].Schema).Data(nil, diff)
		resourceData1.SetId("MockId")
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["Normal"]("")
		})
		err := resourceAliCloudVpcVpcUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyVpcAttributeNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"cidr_block", "description", "vpc_name", "name", "enable_ipv6", "tags"} {
			switch p["alicloud_vpc"].Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_vpc"].Schema).Data(d.State(), diff)
		resourceData1.SetId("MockId")
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudVpcVpcUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("UpdateSetInstanceSecondaryCidrBlocksAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"cidr_block", "description", "vpc_name", "name", "enable_ipv6", "tags"} {
			switch p["alicloud_vpc"].Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_vpc"].Schema).Data(nil, diff)
		resourceData1.SetId("MockId")
		retryFlag := false
		noRetryFlag := false
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		patcheDescribeRouteTableList := gomonkey.ApplyMethod(reflect.TypeOf(&VpcService{}), "SetInstanceSecondaryCidrBlocks", func(*VpcService, *schema.ResourceData) error {
			_, err := responseMock["NoRetryError"]("NoRetryError")
			return err
		})
		err := resourceAliCloudVpcVpcUpdate(resourceData1, rawClient)
		patcheDorequest.Reset()
		patcheDescribeRouteTableList.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateSetResourceTagsAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"cidr_block", "description", "vpc_name", "name", "enable_ipv6", "tags"} {
			switch p["alicloud_vpc"].Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_vpc"].Schema).Data(nil, diff)
		resourceData1.SetId("MockId")
		retryFlag := false
		noRetryFlag := false
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		patcheSetResourceTags := gomonkey.ApplyMethod(reflect.TypeOf(&VpcService{}), "SetResourceTags", func(*VpcService, *schema.ResourceData, string) error {
			_, err := responseMock["NoRetryError"]("NoRetryError")
			return err
		})
		err := resourceAliCloudVpcVpcUpdate(resourceData1, rawClient)
		patcheDorequest.Reset()
		patcheSetResourceTags.Reset()
		assert.NotNil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudVpcVpcDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAliCloudVpcVpcDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("DeleteMockNormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Forbidden.VpcNotFound")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAliCloudVpcVpcDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("DeleteIsExpectedErrorsMockNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Forbidden.VpcNotFound")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAliCloudVpcVpcDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Read
	t.Run("ReadDescribeRouteTableListAbnormal", func(t *testing.T) {
		d.SetId("MockId")
		d.Set("cidr_block", "cidr_block")
		d.Set("description", "description")
		d.Set("dry_run", false)
		d.Set("enable_ipv6", false)
		d.Set("resource_group_id", "resource_group_id")
		d.Set("resource_group_id", []string{"resource_group_id_1", "resource_group_id_2"})
		d.Set("vpc_name", "vpc_name")
		d.Set("name", "name")
		d.Set("user_cidrs", []interface{}{"user_cidrs_1", "user_cidrs_2"})
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			retryFlag := false
			noRetryFlag := true
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})
		patcheDescribeRouteTableList := gomonkey.ApplyMethod(reflect.TypeOf(&VpcService{}), "DescribeRouteTableList", func(*VpcService, string) (map[string]interface{}, error) {
			return responseMock["NoRetryError"]("NoRetryError")
		})
		err := resourceAliCloudVpcVpcRead(d, rawClient)
		patcheDorequest.Reset()
		patcheDescribeRouteTableList.Reset()
		assert.NotNil(t, err)
	})
}

func TestAccAliCloudVPCVPC_basic3113(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcVpcMap3113)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcVpc")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcVpcBasicDependence3113)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_name":            name,
					"dns_hostname_status": "ENABLED",
					"force_delete":        "true",
					"system_route_table_route_propagation_enable": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_name":            name,
						"dns_hostname_status": "ENABLED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":         "test",
					"dns_hostname_status": "DISABLED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":         "test",
						"dns_hostname_status": "DISABLED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cidr_block": "10.0.0.0/8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr_block": "10.0.0.0/8",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cidr_block": "172.16.0.0/12",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr_block": "172.16.0.0/12",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_route_table_route_propagation_enable": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_route_table_route_propagation_enable": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"classic_link_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"classic_link_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"classic_link_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"classic_link_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6_isp":    "BGP",
					"description": "test",
					"cidr_block":  "10.0.0.0/8",
					"vpc_name":    name + "_update",
					"enable_ipv6": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_isp":    "BGP",
						"description": "test",
						"cidr_block":  "10.0.0.0/8",
						"vpc_name":    name + "_update",
						"enable_ipv6": "false",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "enable_ipv6", "ipv6_isp", "force_delete"},
			},
		},
	})
}

var AlicloudVpcVpcMap3113 = map[string]string{
	"route_table_id":    CHECKSET,
	"resource_group_id": CHECKSET,
	"cidr_block":        "172.16.0.0/12",
	"router_id":         CHECKSET,
	"status":            CHECKSET,
	"create_time":       CHECKSET,
	"system_route_table_route_propagation_enable": "true",
}

func AlicloudVpcVpcBasicDependence3113(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 3113  twin
func TestAccAliCloudVPCVPC_basic3113_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcVpcMap3113)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcVpc")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcVpcBasicDependence3113)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6_isp":             "BGP",
					"description":          "test-update",
					"cidr_block":           "10.0.0.0/16",
					"vpc_name":             name,
					"classic_link_enabled": "false",
					"system_route_table_route_propagation_enable": "false",
					"secondary_cidr_blocks": []string{
						"192.168.0.0/16"},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_isp":             "BGP",
						"description":          "test-update",
						"cidr_block":           "10.0.0.0/16",
						"vpc_name":             name,
						"classic_link_enabled": "false",
						"system_route_table_route_propagation_enable": "false",
						"secondary_cidr_blocks.#":                     "1",
						"tags.%":                                      "2",
						"tags.Created":                                "TF",
						"tags.For":                                    "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secondary_cidr_blocks": []string{"10.1.0.0/16", "172.16.0.0/16"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secondary_cidr_blocks.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secondary_cidr_blocks": []string{"10.1.0.0/16", "192.168.0.0/16", "172.16.0.0/16"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secondary_cidr_blocks.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secondary_cidr_blocks": []string{"192.168.0.0/16"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secondary_cidr_blocks.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secondary_cidr_blocks": []string{"10.1.0.0/16", "192.168.0.0/16", "172.16.0.0/16"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secondary_cidr_blocks.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secondary_cidr_blocks": []string{"192.168.0.0/16", "10.1.0.0/16", "172.16.0.0/16"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secondary_cidr_blocks.#": "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "enable_ipv6", "ipv6_isp"},
			},
		},
	})
}

// Test Vpc Vpc. >>> Resource test cases, automatically generated.
// Case 从IPAM地址池创建VPC及添加附加网段 9656
func TestAccAliCloudVPCVPC_basic9656(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcVpcMap9656)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcVpc")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcVpcBasicDependence9656)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv4_cidr_mask":    "12",
					"ipv4_ipam_pool_id": "${alicloud_vpc_ipam_ipam_pool_cidr.defaultIpamPoolCidr.ipam_pool_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv4_cidr_mask":    "12",
						"ipv4_ipam_pool_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secondary_cidr_mask": "24",
					"description":         "desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secondary_cidr_mask": "24",
						"description":         "desc",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "ipv4_cidr_mask", "ipv4_ipam_pool_id", "ipv6_isp", "secondary_cidr_mask"},
			},
		},
	})
}

var AlicloudVpcVpcMap9656 = map[string]string{
	"router_id":   CHECKSET,
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudVpcVpcBasicDependence9656(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc_ipam_ipam" "defaultIpam" {
  operating_region_list = ["cn-hangzhou"]
}

resource "alicloud_vpc_ipam_ipam_pool" "defaultIpamPool" {
  ipam_scope_id  = alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id
  pool_region_id = "cn-hangzhou"
}

resource "alicloud_vpc_ipam_ipam_pool_cidr" "defaultIpamPoolCidr" {
  cidr         = "10.0.0.0/8"
  ipam_pool_id = alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id
}


`, name)
}

// Test Vpc Vpc. <<< Resource test cases, automatically generated.

// DataWorks 2024 read actions require GET; other scoped RPC owners retain
// their native POST transport. Keep the distinction in the production adapter.
func nativeVpcSweepRPC(client *connectivity.AliyunClient) vpcSweepRPC {
	return routedVpcSweepRPC(client.RpcGet, client.RpcPost)
}

func routedVpcSweepRPC(get func(string, string, string, map[string]interface{}, map[string]interface{}) (map[string]interface{}, error), post vpcSweepRPC) vpcSweepRPC {
	dataworks := dataWorksVpcCleanupRPC(get, post)
	eci := eciVpcCleanupRPC(get, post)
	return func(product, version, action string, query, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
		if product == "dataworks-public" {
			return dataworks(product, version, action, query, body, autoRetry)
		}
		if strings.EqualFold(product, "eci") {
			return eci(product, version, action, query, body, autoRetry)
		}
		return post(product, version, action, query, body, autoRetry)
	}
}

// The VPC orchestrator reuses the established SWEEP_ALL environment switch.
// Other resources retain their existing SWEEP_ALL SDK entrypoints.
func sweepAllVpcs() bool { return sweepAll() }

func validateVpcSweepSelector(selector string) error {
	if selector != "alicloud_vpc" {
		return fmt.Errorf("the protected VPC orchestrator requires exactly -sweep-run=alicloud_vpc")
	}
	return nil
}

// The SDK registry uses substring selection. SWEEP_ALL VPC cleanup must bypass that
// selector and invoke only the same native VPC sweeper, never sibling sweepers
// such as alicloud_vpc_ipv6_gateway or alicloud_vpc_peer_connection.
func runVpcSweep(regionString, selector string, sweepVpc func(string) error) error {
	if err := validateVpcSweepSelector(selector); err != nil {
		return err
	}
	var regions []string
	seen := map[string]bool{}
	for _, value := range strings.Split(regionString, ",") {
		region := strings.TrimSpace(value)
		if region == "" {
			return fmt.Errorf("SWEEP_ALL VPC sweep requires non-empty regions; no sweepers were run")
		}
		if !seen[region] {
			regions = append(regions, region)
			seen[region] = true
		}
	}
	var failures []string
	for _, region := range regions {
		if err := sweepVpc(region); err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", region, err))
		}
	}
	if len(failures) > 0 {
		return fmt.Errorf("VPC sweep incomplete: %s", strings.Join(failures, "; "))
	}
	return nil
}

func vpcSweepDependencies(normal []string) []string {
	if sweepAllVpcs() {
		return nil
	}
	return normal
}

func protectedSweepVpcName(name string) bool {
	normalized := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return unicode.ToLower(r)
		}
		return -1
	}, name)
	for _, marker := range []string{"nodeleting", "nodelete", "donotdelete", "dontdelete", "neverdelete", "notdelete", "禁止删除", "不要删除", "勿删除", "不可删除", "不能删除", "请勿删除", "保留", "donotremove", "noremove", "donotdestroy", "不要删", "勿删", "禁止删", "不能删", "不可删"} {
		if strings.Contains(normalized, marker) {
			return true
		}
	}
	return false
}

func skipVpcSweep(item map[string]interface{}, prefixes []string) (bool, string) {
	if protectedSweepVpcName(fmt.Sprint(item["VpcName"])) {
		return true, "explicit do-not-delete VPC; all dependencies protected"
	}
	if !sweepAllVpcs() {
		if isDefault, _ := item["IsDefault"].(bool); isDefault {
			return true, "default VPC"
		}
		matches := false
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["VpcName"])), strings.ToLower(prefix)) {
				matches = true
				break
			}
		}
		if !matches {
			return true, "name does not match test prefixes"
		}
	}
	return sweepRecentCheck(item)
}

type vpcSweepRPC func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error)
type allVpcOwner struct {
	product, version, list, path, filter, vpcPath, id, remove string
	extra                                                     map[string]interface{}
}

// PolarDB's list API has no VPC filter. Its returned VpcId MUST be checked
// locally. Other owner APIs also require an explicit returned VPC match, so an
// ignored server filter can never delete another VPC's resources.
var allVpcOwners = []allVpcOwner{
	{"polardb", "2017-08-01", "DescribeDBClusters", "$.Items.DBCluster", "", "$.VpcId", "DBClusterId", "DeleteDBCluster", map[string]interface{}{"BackupRetentionPolicyOnClusterDeletion": "NONE"}},
	{"Rds", "2014-08-15", "DescribeDBInstances", "$.Items.DBInstance", "VpcId", "$.VpcId", "DBInstanceId", "DeleteDBInstance", nil},
	{"R-kvstore", "2015-01-01", "DescribeInstances", "$.Instances.KVStoreInstance", "VpcId", "$.VpcId", "InstanceId", "DeleteInstance", nil},
	{"Slb", "2014-05-15", "DescribeLoadBalancers", "$.LoadBalancers.LoadBalancer", "VpcId", "$.VpcId", "LoadBalancerId", "DeleteLoadBalancer", nil},
	{"Ecs", "2014-05-26", "DescribeInstances", "$.Instances.Instance", "VpcId", "$.VpcAttributes.VpcId", "InstanceId", "DeleteInstance", map[string]interface{}{"Force": true}},
	{"gpdb", "2016-05-03", "DescribeDBInstances", "$.Items.DBInstance", "", "$.VpcId", "DBInstanceId", "DeleteDBInstance", nil},
}

func allVpcOwnerItems(rpc vpcSweepRPC, region, vpcID string, owner allVpcOwner) ([]map[string]interface{}, error) {
	var items []map[string]interface{}
	for page := 1; ; page++ {
		request := map[string]interface{}{"RegionId": region, "PageSize": PageSizeLarge, "PageNumber": page}
		if owner.filter != "" {
			request[owner.filter] = vpcID
		}
		response, err := rpc(owner.product, owner.version, owner.list, nil, request, true)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", owner.list, err)
		}
		raw, err := jsonpath.Get(owner.path, response)
		if err != nil {
			return nil, fmt.Errorf("%s response: %w", owner.list, err)
		}
		rows, ok := raw.([]interface{})
		if !ok {
			return nil, fmt.Errorf("%s returned invalid resource list", owner.list)
		}
		for _, row := range rows {
			item, ok := row.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("%s returned invalid item", owner.list)
			}
			returnedVpc, err := jsonpath.Get(owner.vpcPath, item)
			if err != nil || returnedVpc == nil {
				return nil, fmt.Errorf("%s resource has no verifiable VPC identity; refusing deletion", owner.list)
			}
			// Classic-network resources have an explicit empty VPC and are outside this VPC cleanup.
			if fmt.Sprint(returnedVpc) == vpcID {
				items = append(items, item)
			}
		}
		if len(rows) < PageSizeLarge {
			return items, nil
		}
	}
}

func deleteAllVpcOwners(rpc vpcSweepRPC, region, vpcID string, owner allVpcOwner) error {
	items, err := allVpcOwnerItems(rpc, region, vpcID, owner)
	if err != nil {
		return err
	}
	var failures []string
	for _, item := range items {
		id, ok := item[owner.id].(string)
		if !ok || id == "" {
			failures = append(failures, fmt.Sprintf("%s has invalid %s; refusing deletion", owner.list, owner.id))
			continue
		}
		if owner.product == "polardb" {
			if err := verifyVpcPolarDBClusterExclusive(rpc, region, vpcID, id); err != nil {
				failures = append(failures, err.Error())
				continue
			}
		}
		if owner.product == "R-kvstore" && item["ChargeType"] == "PrePaid" {
			if err := refundVpcKvstoreSubscription(rpc, region, vpcID, id, owner); err != nil {
				failures = append(failures, err.Error())
			}
			continue
		}
		request := map[string]interface{}{"RegionId": region, owner.id: id}
		for key, value := range owner.extra {
			request[key] = value
		}
		if err := disableAllVpcOwnerDeletionProtection(rpc, region, owner, item, id); err != nil {
			failures = append(failures, err.Error())
			continue
		}
		log.Printf("[INFO] SWEEP_ALL VPC %s: %s %s (release owning service before its managed ENIs)", vpcID, owner.remove, id)
		if _, err := rpc(owner.product, owner.version, owner.remove, nil, request, true); err != nil {
			failures = append(failures, fmt.Sprintf("%s %s: %v", owner.remove, id, err))
		}
	}
	if len(failures) > 0 {
		return fmt.Errorf("%s deletion incomplete: %s", owner.remove, strings.Join(failures, "; "))
	}
	return nil
}

// Primary VPC identity alone is insufficient: a cluster can expose private
// endpoint addresses in other VPCs. Verify all addresses before changing the
// cluster's deletion protection or deleting it. Public addresses are not VPC
// attachments; explicit empty endpoint/address lists are valid API results.
func verifyVpcPolarDBClusterExclusive(rpc vpcSweepRPC, region, vpcID, clusterID string) error {
	detail, err := rpc("polardb", "2017-08-01", "DescribeDBClusterAttribute", nil, map[string]interface{}{"RegionId": region, "DBClusterId": clusterID}, true)
	if err != nil {
		return fmt.Errorf("PolarDB %s primary identity: %w", clusterID, err)
	}
	if detail["DBClusterId"] != clusterID || detail["VPCId"] != vpcID {
		return fmt.Errorf("PolarDB %s primary cluster/VPC identity unknown or changed; refusing shared-cluster deletion", clusterID)
	}
	response, err := rpc("polardb", "2017-08-01", "DescribeDBClusterEndpoints", nil, map[string]interface{}{"RegionId": region, "DBClusterId": clusterID}, true)
	if err != nil {
		return fmt.Errorf("PolarDB %s endpoint inventory: %w", clusterID, err)
	}
	endpoints, ok := response["Items"].([]interface{})
	if !ok {
		return fmt.Errorf("PolarDB %s missing or invalid endpoint Items; refusing deletion", clusterID)
	}
	seen := map[string]bool{}
	for _, raw := range endpoints {
		endpoint, ok := raw.(map[string]interface{})
		if !ok {
			return fmt.Errorf("PolarDB %s malformed endpoint; refusing deletion", clusterID)
		}
		endpointID, ok := endpoint["DBEndpointId"].(string)
		if !ok || endpointID == "" || seen[endpointID] {
			return fmt.Errorf("PolarDB %s unknown or duplicate endpoint ID; refusing deletion", clusterID)
		}
		seen[endpointID] = true
		addresses, ok := endpoint["AddressItems"].([]interface{})
		if !ok {
			return fmt.Errorf("PolarDB %s endpoint %s missing or invalid AddressItems; refusing deletion", clusterID, endpointID)
		}
		for _, rawAddress := range addresses {
			address, ok := rawAddress.(map[string]interface{})
			if !ok {
				return fmt.Errorf("PolarDB %s endpoint %s malformed address; refusing deletion", clusterID, endpointID)
			}
			switch address["NetType"] {
			case "Public":
				continue
			case "Private":
				addressVpc, ok := address["VPCId"].(string)
				vswitch, validSwitch := address["VSwitchId"].(string)
				if !ok || addressVpc == "" || !validSwitch || vswitch == "" {
					return fmt.Errorf("PolarDB %s endpoint %s private address ownership unknown; refusing deletion", clusterID, endpointID)
				}
				if addressVpc != vpcID {
					return fmt.Errorf("PolarDB %s endpoint %s serves another VPC; refusing shared-cluster deletion", clusterID, endpointID)
				}
			default:
				return fmt.Errorf("PolarDB %s endpoint %s unknown NetType; refusing deletion", clusterID, endpointID)
			}
		}
	}
	return nil
}

// Names on the VPC define the protection boundary in SWEEP_ALL VPC cleanup. A resource's
// API-level deletion lock is cleared only after its parent VPC is selected and
// its returned ownership has been verified.
func disableAllVpcOwnerDeletionProtection(rpc vpcSweepRPC, region string, owner allVpcOwner, item map[string]interface{}, id string) error {
	action, field, value := "", "", interface{}(false)
	switch owner.product {
	case "polardb":
		if fmt.Sprint(item["DeletionLock"]) == "1" {
			action, field = "ModifyDBClusterDeletion", "Protection"
		}
	case "Rds":
		if item["DeletionProtection"] == true {
			action, field = "ModifyDBInstanceDeletionProtection", "DeletionProtection"
		}
	case "Slb":
		if strings.EqualFold(fmt.Sprint(item["DeleteProtection"]), "on") {
			action, field, value = "SetLoadBalancerDeleteProtection", "DeleteProtection", "off"
		}
	case "Ecs":
		if item["DeletionProtection"] == true {
			action, field = "ModifyInstanceAttribute", "DeletionProtection"
		}
	}
	if action == "" {
		return nil
	}
	_, err := rpc(owner.product, owner.version, action, nil, map[string]interface{}{"RegionId": region, owner.id: id, field: value}, true)
	if err != nil {
		return fmt.Errorf("%s %s: %w", action, id, err)
	}
	return nil
}

func sweepVpcOwner(rpc vpcSweepRPC, region, vpcID string, owner allVpcOwner) error {
	if err := deleteAllVpcOwners(rpc, region, vpcID, owner); err != nil {
		return err
	}
	return retry.Retry(10*time.Minute, func() *retry.RetryError {
		items, err := allVpcOwnerItems(rpc, region, vpcID, owner)
		if err != nil {
			return retry.NonRetryableError(err)
		}
		if len(items) > 0 {
			time.Sleep(5 * time.Second)
			return retry.RetryableError(fmt.Errorf("%s still has %d resources in VPC %s", owner.list, len(items), vpcID))
		}
		return nil
	})
}

// Managed or attached interfaces must be released by their owning service.
func allVpcENISkipReason(item map[string]interface{}) string {
	if managed, ok := item["ServiceManaged"].(bool); !ok || managed {
		return "service-managed or unknown ownership; release the owning service first"
	}
	if item["Status"] != "Available" || item["Type"] != "Secondary" {
		return "interface still attached or not secondary; release the owning service first"
	}
	return ""
}

// Bound cross-VPC concurrency so a single slow dependent service does not
// serialize a whole region, without launching one goroutine per VPC.
func runVpcSweeps(ids []string, workers int, sweep func(string) error) error {
	if workers < 1 {
		workers = 1
	}
	jobs := make(chan string)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var failures []string
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range jobs {
				if err := sweep(id); err != nil {
					mu.Lock()
					failures = append(failures, err.Error())
					mu.Unlock()
				}
			}
		}()
	}
	for _, id := range ids {
		jobs <- id
	}
	close(jobs)
	wg.Wait()
	if len(failures) > 0 {
		return fmt.Errorf("VPC cleanup incomplete: %s", strings.Join(failures, "; "))
	}
	return nil
}

// SWEEP_ALL VPC cleanup must not trust server-side filtering: ignored filters must
// never expose a protected VPC's dependent resources to a delete callback.
func sweepDependentResourceScope(item map[string]interface{}, paramKey, expected string) (bool, string) {
	if !sweepAllVpcs() {
		return true, ""
	}
	return vpcDependentResourceIdentity(item, paramKey, expected)
}

func vpcDependentResourceIdentity(item map[string]interface{}, paramKey, expected string) (bool, string) {
	if expected == "" {
		return false, "empty expected VPC identity"
	}
	switch paramKey {
	case "VpcId":
		actual, ok := item["VpcId"].(string)
		if !ok || actual == "" {
			return false, "missing or invalid returned VpcId"
		}
		if actual != expected {
			return false, "returned VpcId differs from selected VPC"
		}
		return true, ""
	case "ResourceId":
		actual, ok := item["ResourceId"].(string)
		if !ok || actual != expected {
			return false, "returned ResourceId differs from selected VPC"
		}
		resourceType, ok := item["ResourceType"].(string)
		if !ok || resourceType != "VPC" {
			return false, "flow log ResourceType is not VPC"
		}
		return true, ""
	default:
		return false, "unsupported ownership filter"
	}
}

func TestUnitSweepAllVpcSelection(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	t.Setenv("ALICLOUD_SWEEP_OLDER_THAN", "")
	for _, name := range []string{"default-NODELETING", "cfw-test-no-delete2", "Do_Not_Delete", "do-not-delete", "dont-delete", "NEVER DELETE", "不要删除", "生产-请勿删除-vpc"} {
		t.Run(name, func(t *testing.T) {
			if skip, _ := skipVpcSweep(map[string]interface{}{"VpcName": name}, nil); !skip {
				t.Fatal("protected VPC selected")
			}
		})
	}
	for _, name := range []string{"default", "vpc", "manual-business-name", "tf-test-new", ""} {
		t.Run(name, func(t *testing.T) {
			item := map[string]interface{}{"VpcName": name, "IsDefault": true, "CreationTime": time.Now().UTC().Format(time.RFC3339)}
			if skip, reason := skipVpcSweep(item, []string{"tf-test"}); skip {
				t.Fatalf("SWEEP_ALL without an age filter must ignore name and default status: %s", reason)
			}
			delete(item, "CreationTime")
			if skip, reason := skipVpcSweep(item, nil); skip {
				t.Fatalf("SWEEP_ALL without an age filter must not require creation time: %s", reason)
			}
		})
	}
}

func TestUnitSweepAllVpcSelectorAndDependencies(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	if err := validateVpcSweepSelector("alicloud_vpc"); err != nil {
		t.Fatal(err)
	}
	for _, selector := range []string{"", "alicloud_vpc|alicloud_instance", "alicloud_instance", "alicloud_vpc,alicloud_instance", "alicloud_vpc.*", "^alicloud_vpc$"} {
		if err := validateVpcSweepSelector(selector); err == nil {
			t.Fatalf("unsafe VPC orchestrator selector accepted: %q", selector)
		}
	}
	if got := vpcSweepDependencies([]string{"alicloud_instance", "alicloud_vswitch"}); len(got) != 0 {
		t.Fatalf("regional sweepers could delete protected VPC dependencies: %v", got)
	}
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "false")
	if got := vpcSweepDependencies([]string{"alicloud_instance"}); len(got) != 1 {
		t.Fatal("legacy dependencies changed")
	}
}

func TestUnitSweepAllVpcNormalCompatibility(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "")
	t.Setenv("ALICLOUD_SWEEP_OLDER_THAN", "12h")
	for _, tc := range []struct {
		name      string
		isDefault bool
		age       time.Duration
		skip      bool
	}{
		{"tf-test-old", false, 24 * time.Hour, false},
		{"tf-test-new", false, time.Hour, true},
		{"tf-test-no-delete", false, 24 * time.Hour, true},
		{"other-name", false, 24 * time.Hour, true},
		{"tf-test-default", true, 24 * time.Hour, true},
	} {
		item := map[string]interface{}{"VpcName": tc.name, "IsDefault": tc.isDefault, "CreationTime": time.Now().Add(-tc.age).UTC().Format(time.RFC3339)}
		if skip, _ := skipVpcSweep(item, []string{"tf-test"}); skip != tc.skip {
			t.Errorf("%s skip=%v want %v", tc.name, skip, tc.skip)
		}
	}
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	if skip, _ := skipVpcSweep(map[string]interface{}{"VpcName": "other", "CreationTime": time.Now().Add(-24 * time.Hour).Format(time.RFC3339)}, nil); skip {
		t.Fatal("legacy SWEEP_ALL changed")
	}
}

func TestUnitSweepAllVpcPolardbOwnerScope(t *testing.T) {
	owner := allVpcOwners[0]
	var deleted []string
	rpc := vpcSweepRPC(func(product, version, action string, query, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
		if action == "DescribeDBClusterAttribute" {
			return map[string]interface{}{"DBClusterId": body["DBClusterId"], "VPCId": "vpc-target"}, nil
		}
		if action == "DescribeDBClusterEndpoints" {
			return map[string]interface{}{"Items": []interface{}{}}, nil
		}
		if action == "DescribeDBClusters" {
			if _, ok := body["VpcId"]; ok {
				t.Fatal("PolarDB list does not accept a VpcId filter")
			}
			return map[string]interface{}{"Items": map[string]interface{}{"DBCluster": []interface{}{
				map[string]interface{}{"DBClusterId": "pc-delete", "VpcId": "vpc-target", "DBClusterDescription": "any-non-test-name"},
				map[string]interface{}{"DBClusterId": "pc-protected", "VpcId": "vpc-protected"},
				map[string]interface{}{"DBClusterId": "pc-classic", "VpcId": ""},
			}}}, nil
		}
		if action != "DeleteDBCluster" {
			t.Fatalf("must delete owning service, not its managed ENI: %s", action)
		}
		if body["BackupRetentionPolicyOnClusterDeletion"] != "NONE" {
			t.Fatal("cluster backups would remain")
		}
		deleted = append(deleted, fmt.Sprint(body["DBClusterId"]))
		return map[string]interface{}{}, nil
	})
	if err := deleteAllVpcOwners(rpc, "cn-beijing", "vpc-target", owner); err != nil {
		t.Fatal(err)
	}
	if strings.Join(deleted, ",") != "pc-delete" {
		t.Fatalf("cross-VPC deletion: %v", deleted)
	}
}

func TestUnitSweepAllVpcRefusesUnverifiableOwners(t *testing.T) {
	owner := allVpcOwners[0]
	for _, rows := range []interface{}{
		"invalid", []interface{}{map[string]interface{}{"DBClusterId": "pc-unknown"}},
		[]interface{}{map[string]interface{}{"VpcId": "vpc-target"}},
	} {
		rpc := vpcSweepRPC(func(product, version, action string, query, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
			if action != "DescribeDBClusters" {
				t.Fatalf("delete attempted without verified VPC and ID: %s", action)
			}
			return map[string]interface{}{"Items": map[string]interface{}{"DBCluster": rows}}, nil
		})
		if err := deleteAllVpcOwners(rpc, "cn-beijing", "vpc-target", owner); err == nil {
			t.Fatalf("unsafe rows accepted: %v", rows)
		}
	}
}

func TestUnitSweepAllVpcOwnerErrorsPropagate(t *testing.T) {
	rpc := vpcSweepRPC(func(product, version, action string, query, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
		return nil, fmt.Errorf("access denied")
	})
	if err := deleteAllVpcOwners(rpc, "cn-beijing", "vpc-target", allVpcOwners[0]); err == nil {
		t.Fatal("list failure hidden")
	}
}

func TestUnitSweepAllVpcManagedENIProtection(t *testing.T) {
	for _, item := range []map[string]interface{}{
		{"ServiceManaged": true, "Status": "Available", "Type": "Secondary"},
		{"ServiceManaged": false, "Status": "InUse", "Type": "Secondary"},
		{"ServiceManaged": false, "Status": "Available", "Type": "Primary"},
		{"Status": "Available", "Type": "Secondary"},
	} {
		if allVpcENISkipReason(item) == "" {
			t.Fatalf("unsafe interface selected: %v", item)
		}
	}
	if reason := allVpcENISkipReason(map[string]interface{}{"ServiceManaged": false, "Status": "Available", "Type": "Secondary"}); reason != "" {
		t.Fatal(reason)
	}
}

func TestUnitSweepAllVpcOwnerDeletionProtection(t *testing.T) {
	cases := []struct {
		index         int
		item          map[string]interface{}
		action, field string
		value         interface{}
	}{
		{0, map[string]interface{}{"DeletionLock": 1}, "ModifyDBClusterDeletion", "Protection", false},
		{1, map[string]interface{}{"DeletionProtection": true}, "ModifyDBInstanceDeletionProtection", "DeletionProtection", false},
		{3, map[string]interface{}{"DeleteProtection": "on"}, "SetLoadBalancerDeleteProtection", "DeleteProtection", "off"},
		{4, map[string]interface{}{"DeletionProtection": true}, "ModifyInstanceAttribute", "DeletionProtection", false},
	}
	for _, tc := range cases {
		called := false
		rpc := vpcSweepRPC(func(product, version, action string, query, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
			called = true
			if action != tc.action || body[tc.field] != tc.value || body[allVpcOwners[tc.index].id] != "verified-id" {
				t.Fatalf("wrong scoped protection request: %s %v", action, body)
			}
			return nil, nil
		})
		if err := disableAllVpcOwnerDeletionProtection(rpc, "cn-beijing", allVpcOwners[tc.index], tc.item, "verified-id"); err != nil {
			t.Fatal(err)
		}
		if !called {
			t.Fatal("owner deletion protection was not released")
		}
	}
}

func TestUnitSweepAllVpcContinuesAfterVpcFailure(t *testing.T) {
	var mu sync.Mutex
	visited := map[string]int{}
	err := runVpcSweeps([]string{"vpc-a", "vpc-b", "vpc-c"}, 2, func(id string) error {
		mu.Lock()
		visited[id]++
		mu.Unlock()
		if id == "vpc-b" {
			return fmt.Errorf("vpc-b managed ENI remains")
		}
		return nil
	})
	if err == nil || !strings.Contains(err.Error(), "managed ENI remains") {
		t.Fatalf("missing partial failure: %v", err)
	}
	for _, id := range []string{"vpc-a", "vpc-b", "vpc-c"} {
		if visited[id] != 1 {
			t.Fatalf("%s visited %d times", id, visited[id])
		}
	}
}

func TestUnitSweepAllVpcProtectedVpcInventoryAndOwners(t *testing.T) {
	for _, ownerFailure := range []bool{false, true} {
		t.Run(fmt.Sprintf("ownerFailure=%v", ownerFailure), func(t *testing.T) {
			t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
			t.Setenv("ALICLOUD_SWEEP_OLDER_THAN", "")
			var calls, selected, deleted []string
			rpc := vpcSweepRPC(func(product, version, action string, query, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
				if body["VpcId"] == "vpc-protected" {
					t.Errorf("protected VPC dependencies queried for deletion: %s", action)
					return nil, fmt.Errorf("protected VPC queried")
				}
				calls = append(calls, product+"/"+action)
				if action == "DescribeDBClusterAttribute" {
					return map[string]interface{}{"DBClusterId": body["DBClusterId"], "VPCId": "vpc-target"}, nil
				}
				if action == "DescribeDBClusterEndpoints" {
					return map[string]interface{}{"Items": []interface{}{}}, nil
				}
				if action == "DescribeVpcs" {
					return allVpcSweepTestResponse("$.Vpcs.Vpc", []interface{}{
						map[string]interface{}{"VpcId": "vpc-protected", "VpcName": "default-NODELETING", "IsDefault": true},
						map[string]interface{}{"VpcId": "vpc-target", "VpcName": "manual-vpc", "IsDefault": true, "CreationTime": time.Now().Format(time.RFC3339)},
					}), nil
				}
				for _, owner := range allVpcOwners {
					if owner.product == product && owner.list == action {
						if ownerFailure && product == "polardb" {
							return nil, fmt.Errorf("Polardb list denied")
						}
						rows := []interface{}{map[string]interface{}{owner.id: "protected-owner", "VpcId": "vpc-protected"}}
						if product == "Ecs" {
							rows = []interface{}{map[string]interface{}{owner.id: "protected-owner", "VpcAttributes": map[string]interface{}{"VpcId": "vpc-protected"}}}
						}
						if product == "polardb" {
							rows = append(rows, map[string]interface{}{owner.id: "pc-target", "VpcId": "vpc-target"})
						}
						return allVpcSweepTestResponse(owner.path, rows), nil
					}
				}
				if action == "DeleteDBCluster" && body["DBClusterId"] == "pc-target" {
					deleted = append(deleted, "pc-target")
					return nil, nil
				}
				t.Errorf("unexpected or protected resource deletion: %s/%s %v", product, action, body)
				return nil, fmt.Errorf("unexpected RPC")
			})
			err := sweepVpcsWithCallbacks("cn-beijing", rpc, func(id string) error {
				selected = append(selected, id)
				var failures []string
				for _, owner := range allVpcOwners {
					if err := deleteAllVpcOwners(rpc, "cn-beijing", id, owner); err != nil {
						failures = append(failures, err.Error())
					}
				}
				if len(failures) > 0 {
					return fmt.Errorf("%s", strings.Join(failures, "; "))
				}
				return nil
			})
			if strings.Join(selected, ",") != "vpc-target" {
				t.Fatalf("protected VPC reached deletion callback: %v", selected)
			}
			if ownerFailure {
				if err == nil || !strings.Contains(err.Error(), "Polardb list denied") {
					t.Fatalf("owner error hidden: %v", err)
				}
				if len(deleted) > 0 {
					t.Fatalf("owner deleted after failed inventory: %v", deleted)
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				if strings.Join(deleted, ",") != "pc-target" {
					t.Fatalf("wrong scoped owner deletions: %v", deleted)
				}
			}
			if !strings.Contains(strings.Join(calls, ","), "gpdb/DescribeDBInstances") {
				t.Fatal("one failed owner blocked other scoped owners")
			}
		})
	}
}

func allVpcSweepTestResponse(path string, rows []interface{}) map[string]interface{} {
	var value interface{} = rows
	keys := strings.Split(strings.TrimPrefix(path, "$."), ".")
	for i := len(keys) - 1; i >= 0; i-- {
		value = map[string]interface{}{keys[i]: value}
	}
	return value.(map[string]interface{})
}

func TestUnitSweepAllVpcGpdbOwnerScope(t *testing.T) {
	owner := allVpcOwners[len(allVpcOwners)-1]
	deleted := false
	rpc := vpcSweepRPC(func(product, version, action string, query, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
		if product != "gpdb" || version != "2016-05-03" {
			t.Fatalf("wrong GPDB API: %s %s", product, version)
		}
		if action == "DescribeDBInstances" {
			if _, ok := body["VpcId"]; ok {
				t.Fatal("GPDB list does not support VPC filter")
			}
			return allVpcSweepTestResponse("$.Items.DBInstance", []interface{}{
				map[string]interface{}{"DBInstanceId": "gp-target", "VpcId": "vpc-target"},
				map[string]interface{}{"DBInstanceId": "gp-protected", "VpcId": "vpc-protected"},
			}), nil
		}
		if action != "DeleteDBInstance" || body["DBInstanceId"] != "gp-target" {
			t.Fatalf("unsafe GPDB deletion: %s %v", action, body)
		}
		deleted = true
		return nil, nil
	})
	if err := deleteAllVpcOwners(rpc, "cn-beijing", "vpc-target", owner); err != nil {
		t.Fatal(err)
	}
	if !deleted {
		t.Fatal("selected GPDB owner was not deleted")
	}
}

func TestUnitSweepAllVpcExactNativeCallback(t *testing.T) {
	// Register-like decoys prove that the exact VPC dispatcher only invokes its
	// explicit native VPC callback; it never iterates the SDK substring registry.
	callbacks := map[string]func(string) error{
		"alicloud_vpc":                 func(string) error { return nil },
		"alicloud_vpc_ipv6_gateway":    func(string) error { t.Fatal("sibling sweeper selected"); return nil },
		"alicloud_vpc_peer_connection": func(string) error { t.Fatal("sibling sweeper selected"); return nil },
	}
	var visited []string
	callbacks["alicloud_vpc"] = func(region string) error { visited = append(visited, region); return nil }
	if err := runVpcSweep("cn-beijing, cn-hangzhou,cn-beijing", "alicloud_vpc", callbacks["alicloud_vpc"]); err != nil {
		t.Fatal(err)
	}
	if strings.Join(visited, ",") != "cn-beijing,cn-hangzhou" {
		t.Fatalf("unexpected native VPC invocations: %v", visited)
	}
	for _, selector := range []string{"", "alicloud_vpc_ipv6_gateway", "alicloud_vpc|alicloud_vpc_ipv6_gateway", "alicloud_vpc,alicloud_instance"} {
		if err := runVpcSweep("cn-beijing", selector, func(string) error { t.Fatal("invalid selector executed callback"); return nil }); err == nil {
			t.Fatalf("accepted selector %q", selector)
		}
	}
	for _, regions := range []string{"", "cn-beijing,", ",cn-beijing", "cn-beijing, ,cn-hangzhou"} {
		if err := runVpcSweep(regions, "alicloud_vpc", func(string) error { t.Fatal("invalid region list executed callback"); return nil }); err == nil {
			t.Fatalf("accepted regions %q", regions)
		}
	}
}

func TestUnitSweepAllVpcRegionErrorsAggregate(t *testing.T) {
	var visited []string
	err := runVpcSweep("cn-beijing,cn-hangzhou", "alicloud_vpc", func(region string) error {
		visited = append(visited, region)
		return fmt.Errorf("%s incomplete", region)
	})
	if err == nil || !strings.Contains(err.Error(), "cn-beijing incomplete") || !strings.Contains(err.Error(), "cn-hangzhou incomplete") {
		t.Fatalf("region errors hidden: %v", err)
	}
	if len(visited) != 2 {
		t.Fatalf("remaining region not attempted: %v", visited)
	}
}

func TestUnitSweepAllVpcDependentScope(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	for _, action := range []string{"DescribeVSwitches", "DescribeNatGateways", "DescribeSecurityGroups", "DescribeRouteTableList", "DescribeNetworkAcls", "DescribeVpnGateways", "DescribeHaVips", "DescribeIpv6Gateways", "DescribeNetworkInterfaces"} {
		t.Run(action, func(t *testing.T) {
			for _, tc := range []struct {
				item    map[string]interface{}
				include bool
			}{
				{map[string]interface{}{"VpcId": "vpc-target"}, true},
				{map[string]interface{}{"VpcId": "vpc-protected"}, false},
				{map[string]interface{}{"VPCId": "vpc-target"}, false},
				{map[string]interface{}{"VpcId": ""}, false},
				{map[string]interface{}{"VpcId": nil}, false},
				{map[string]interface{}{}, false},
			} {
				if include, reason := sweepDependentResourceScope(tc.item, "VpcId", "vpc-target"); include != tc.include {
					t.Fatalf("%s ownership filter=%v want %v (%s): %v", action, include, tc.include, reason, tc.item)
				}
			}
		})
	}
	for _, tc := range []struct {
		item    map[string]interface{}
		include bool
	}{
		{map[string]interface{}{"ResourceId": "vpc-target", "ResourceType": "VPC"}, true},
		{map[string]interface{}{"ResourceId": "vpc-protected", "ResourceType": "VPC"}, false},
		{map[string]interface{}{"ResourceId": "vpc-target", "ResourceType": "VSwitch"}, false},
		{map[string]interface{}{"ResourceId": "vpc-target", "ResourceType": "NetworkInterface"}, false},
		{map[string]interface{}{"ResourceId": "vpc-target"}, false},
		{map[string]interface{}{}, false},
	} {
		if include, reason := sweepDependentResourceScope(tc.item, "ResourceId", "vpc-target"); include != tc.include {
			t.Fatalf("flow log ownership filter=%v want %v (%s): %v", include, tc.include, reason, tc.item)
		}
	}
	if include, _ := sweepDependentResourceScope(map[string]interface{}{"VpcId": "vpc-target"}, "Unknown", "vpc-target"); include {
		t.Fatal("unsupported scope accepted")
	}
	if include, _ := sweepDependentResourceScope(map[string]interface{}{"VpcId": ""}, "VpcId", ""); include {
		t.Fatal("empty scope accepted")
	}
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "")
	if include, reason := sweepDependentResourceScope(map[string]interface{}{}, "VpcId", "vpc-target"); !include {
		t.Fatalf("ordinary sweeper filtering changed: %s", reason)
	}
}

func TestUnitSweepAllVpcAgeCompatibility(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	t.Setenv("ALICLOUD_SWEEP_OLDER_THAN", "12h")
	for _, tc := range []struct {
		age  time.Duration
		skip bool
	}{{time.Hour, true}, {24 * time.Hour, false}} {
		item := map[string]interface{}{"VpcName": "manual-name", "IsDefault": true, "CreationTime": time.Now().Add(-tc.age).Format(time.RFC3339)}
		if skip, reason := skipVpcSweep(item, nil); skip != tc.skip {
			t.Fatalf("SWEEP_ALL with age %s skip=%v want %v: %s", tc.age, skip, tc.skip, reason)
		}
	}
	item := map[string]interface{}{"VpcName": "manual-name", "IsDefault": true}
	if skip, _ := skipVpcSweep(item, nil); !skip {
		t.Fatal("explicit OLDER_THAN must protect unknown age")
	}
	t.Setenv("ALICLOUD_SWEEP_OLDER_THAN", "")
	if skip, reason := skipVpcSweep(item, nil); skip {
		t.Fatalf("SWEEP_ALL without age should select arbitrary/default VPC: %s", reason)
	}
}

func TestUnitVpcPolarDBExclusiveEndpointOwnership(t *testing.T) {
	endpoint := func(addresses interface{}) map[string]interface{} {
		return map[string]interface{}{"Items": []interface{}{map[string]interface{}{"DBEndpointId": "pe-target", "AddressItems": addresses}}}
	}
	cases := []struct {
		name     string
		response map[string]interface{}
		valid    bool
	}{
		{"same-private-vpc", endpoint([]interface{}{map[string]interface{}{"NetType": "Private", "VPCId": "vpc-target", "VSwitchId": "vsw-target"}}), true},
		{"public-not-vpc-bound", endpoint([]interface{}{map[string]interface{}{"NetType": "Public"}}), true},
		{"public-vpc-field-not-binding", endpoint([]interface{}{map[string]interface{}{"NetType": "Public", "VPCId": "vpc-other"}}), true},
		{"legitimate-empty-addresses", endpoint([]interface{}{}), true},
		{"legitimate-empty-endpoints", map[string]interface{}{"Items": []interface{}{}}, true},
		{"other-private-vpc", endpoint([]interface{}{map[string]interface{}{"NetType": "Private", "VPCId": "vpc-protected", "VSwitchId": "vsw-protected"}}), false},
		{"private-missing-vpc", endpoint([]interface{}{map[string]interface{}{"NetType": "Private", "VSwitchId": "vsw-target"}}), false},
		{"private-missing-vswitch", endpoint([]interface{}{map[string]interface{}{"NetType": "Private", "VPCId": "vpc-target"}}), false},
		{"unknown-net-type", endpoint([]interface{}{map[string]interface{}{"NetType": "Unknown"}}), false},
		{"missing-net-type", endpoint([]interface{}{map[string]interface{}{}}), false},
		{"missing-items", map[string]interface{}{}, false},
		{"invalid-items", map[string]interface{}{"Items": map[string]interface{}{}}, false},
		{"missing-address-items", map[string]interface{}{"Items": []interface{}{map[string]interface{}{"DBEndpointId": "pe-target"}}}, false},
		{"invalid-address-items", endpoint("invalid"), false},
		{"missing-endpoint-id", map[string]interface{}{"Items": []interface{}{map[string]interface{}{"AddressItems": []interface{}{}}}}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rpc := vpcSweepRPC(func(product, version, action string, query, body map[string]interface{}, auto bool) (map[string]interface{}, error) {
				if product != "polardb" || version != "2017-08-01" || body["DBClusterId"] != "pc-target" {
					t.Fatalf("unexpected ownership API %s/%s %v", product, action, body)
				}
				if action == "DescribeDBClusterAttribute" {
					return map[string]interface{}{"DBClusterId": "pc-target", "VPCId": "vpc-target"}, nil
				}
				if action == "DescribeDBClusterEndpoints" {
					return tc.response, nil
				}
				t.Fatalf("ownership verification mutated a resource: %s", action)
				return nil, nil
			})
			err := verifyVpcPolarDBClusterExclusive(rpc, "cn-beijing", "vpc-target", "pc-target")
			if (err == nil) != tc.valid {
				t.Fatalf("valid=%v err=%v", tc.valid, err)
			}
		})
	}
}

func TestUnitVpcPolarDBSharedGuardPrecedesMutation(t *testing.T) {
	for _, scenario := range []string{"shared-endpoint", "primary-changed", "attribute-error", "endpoint-error"} {
		t.Run(scenario, func(t *testing.T) {
			rpc := vpcSweepRPC(func(product, version, action string, query, body map[string]interface{}, auto bool) (map[string]interface{}, error) {
				switch action {
				case "DescribeDBClusters":
					return allVpcSweepTestResponse("$.Items.DBCluster", []interface{}{map[string]interface{}{"DBClusterId": "pc-target", "VpcId": "vpc-target", "DeletionLock": 1}}), nil
				case "DescribeDBClusterAttribute":
					if scenario == "attribute-error" {
						return nil, fmt.Errorf("attribute denied")
					}
					vpc := "vpc-target"
					if scenario == "primary-changed" {
						vpc = "vpc-protected"
					}
					return map[string]interface{}{"DBClusterId": "pc-target", "VPCId": vpc}, nil
				case "DescribeDBClusterEndpoints":
					if scenario == "endpoint-error" {
						return nil, fmt.Errorf("endpoint list denied")
					}
					return map[string]interface{}{"Items": []interface{}{map[string]interface{}{"DBEndpointId": "pe-shared", "AddressItems": []interface{}{map[string]interface{}{"NetType": "Private", "VPCId": "vpc-protected", "VSwitchId": "vsw-protected"}}}}}, nil
				default:
					t.Fatalf("guard must run before releasing deletion protection or deleting shared cluster: %s", action)
					return nil, nil
				}
			})
			if err := deleteAllVpcOwners(rpc, "cn-beijing", "vpc-target", allVpcOwners[0]); err == nil {
				t.Fatal("unverified/shared cluster reached successful deletion")
			}
		})
	}
}

func nativeVpcFlinkRPC(client *connectivity.AliyunClient) vpcSweepRPC {
	return func(product, version, action string, query, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
		if product != "foasconsole" || version != "2019-06-01" {
			return nil, fmt.Errorf("unsupported Flink product/version")
		}
		switch action {
		case "DescribeInstances", "DescribeSupportedRegions":
			return client.RpcGet(product, version, action, query, body)
		case "DeleteInstance":
			return client.RpcPost(product, version, action, query, body, autoRetry)
		default:
			return nil, fmt.Errorf("unsupported Flink cleanup action %s", action)
		}
	}
}

// The RPC adapter must dispatch DescribeInstances through RpcGet; the official
// foasconsole-20190601 SDK uses GET with flattened DescribeInstancesRequest.*.
// DeleteInstance uses POST with flattened DeleteInstanceRequest.* body fields.
func flinkVpcOwnerSweeper(region, vpcID string, rpc vpcSweepRPC) *resource.Sweeper {
	return &resource.Sweeper{Name: "alicloud_realtime_compute_vvp_instance", F: func(actual string) error {
		if actual != region {
			return fmt.Errorf("Flink region mismatch")
		}
		return sweepFlinkVpcOwner(region, vpcID, rpc, time.Sleep, 180)
	}}
}

func flinkVpcInstances(region, instanceID string, rpc vpcSweepRPC) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	seen := map[string]bool{}
	for page := 1; page <= 1000; page++ {
		query := map[string]interface{}{"DescribeInstancesRequest.Region": region, "DescribeInstancesRequest.PageIndex": page, "DescribeInstancesRequest.PageSize": 100}
		if instanceID != "" {
			query["DescribeInstancesRequest.InstanceId"] = instanceID
		}
		response, err := rpc("foasconsole", "2019-06-01", "DescribeInstances", query, nil, true)
		if err != nil {
			return nil, err
		}
		if success, ok := response["Success"].(bool); ok && !success {
			return nil, fmt.Errorf("Flink DescribeInstances returned Success=false")
		}
		rows, ok := response["Instances"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("Flink DescribeInstances missing Instances")
		}
		for _, raw := range rows {
			item, ok := raw.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("Flink malformed instance")
			}
			id, ok := item["InstanceId"].(string)
			if !ok || id == "" {
				return nil, fmt.Errorf("Flink missing instance identity")
			}
			if seen[id] {
				return nil, fmt.Errorf("Flink repeated instance across pages")
			}
			seen[id] = true
			result = append(result, item)
		}
		// TotalPage is part of the official response. Missing/invalid pagination must
		// fail closed instead of making a truncated region inventory look complete.
		total, err := strconv.Atoi(fmt.Sprint(response["TotalPage"]))
		if err != nil || total < 0 || total > 1000 || (total == 0 && len(rows) > 0) {
			return nil, fmt.Errorf("Flink invalid TotalPage")
		}
		if page >= total {
			return result, nil
		}
		if len(rows) == 0 {
			return nil, fmt.Errorf("Flink empty nonterminal page")
		}
	}
	return nil, fmt.Errorf("Flink pagination limit exceeded")
}

func sweepFlinkVpcOwner(region, vpcID string, rpc vpcSweepRPC, sleep func(time.Duration), pollAttempts int) error {
	if region == "" || vpcID == "" || pollAttempts < 1 {
		return fmt.Errorf("Flink cleanup requires bounded VPC scope")
	}
	instances, err := flinkVpcInstances(region, "", rpc)
	if err != nil {
		return err
	}
	var failures []string
	for _, item := range instances {
		actual, ok := item["VpcId"].(string)
		if !ok || actual == "" {
			failures = append(failures, "Flink instance VPC ownership missing")
			continue
		}
		if actual != vpcID {
			continue
		}
		id := item["InstanceId"].(string)
		charge, _ := item["ChargeType"].(string)
		if charge != "POST" {
			failures = append(failures, fmt.Sprintf("Flink %s requires manual release: billing %q is not POST", id, charge))
			continue
		}
		_, err = rpc("foasconsole", "2019-06-01", "DeleteInstance", nil, map[string]interface{}{"DeleteInstanceRequest.InstanceId": id, "DeleteInstanceRequest.Region": region}, true)
		if err != nil {
			failures = append(failures, fmt.Sprintf("Flink DeleteInstance %s: %v", id, err))
			continue
		}
		released := false
		for attempt := 0; attempt < pollAttempts; attempt++ {
			current, readErr := flinkVpcInstances(region, id, rpc)
			if readErr != nil {
				err = readErr
				break
			}
			found := false
			for _, remaining := range current {
				if remaining["InstanceId"] == id {
					found = true
					if remaining["VpcId"] != vpcID {
						err = fmt.Errorf("Flink ownership changed while waiting for release")
					}
					break
				}
			}
			if err != nil {
				break
			}
			if !found {
				released = true
				break
			}
			if attempt+1 < pollAttempts {
				sleep(10 * time.Second)
			}
		}
		if !released {
			if err == nil {
				err = fmt.Errorf("release confirmation timed out")
			}
			failures = append(failures, fmt.Sprintf("Flink %s: %v", id, err))
		}
	}
	if len(failures) > 0 {
		return fmt.Errorf("Flink owner cleanup incomplete: %s", strings.Join(failures, "; "))
	}
	return nil
}

func TestUnitFlinkVpcOwnerPagedScopedAndConfirmsAsyncRelease(t *testing.T) {
	deletes := 0
	waits := 0
	polls := 0
	row := func(id, vpc, charge string) interface{} {
		return map[string]interface{}{"InstanceId": id, "VpcId": vpc, "ChargeType": charge}
	}
	rpc := vpcSweepRPC(func(product, version, action string, query, body map[string]interface{}, retry bool) (map[string]interface{}, error) {
		if product != "foasconsole" || version != "2019-06-01" {
			t.Fatal("wrong endpoint/version")
		}
		if action == "DeleteInstance" {
			deletes++
			if query != nil || body["DeleteInstanceRequest.InstanceId"] != "selected" || body["DeleteInstanceRequest.Region"] != "cn-beijing" {
				t.Fatal("unscoped deletion")
			}
			return map[string]interface{}{}, nil
		}
		if body != nil || query["DescribeInstancesRequest.Region"] != "cn-beijing" || query["DescribeInstancesRequest.PageSize"] != 100 {
			t.Fatal("invalid GET adapter parameters")
		}
		if query["DescribeInstancesRequest.InstanceId"] == "selected" {
			polls++
			if polls == 1 {
				return map[string]interface{}{"Instances": []interface{}{row("selected", "target", "POST")}, "TotalPage": 1}, nil
			}
			return map[string]interface{}{"Instances": []interface{}{}, "TotalPage": 0}, nil
		}
		if query["DescribeInstancesRequest.PageIndex"] == 1 {
			return map[string]interface{}{"Instances": []interface{}{row("protected", "other", "POST")}, "TotalPage": 2}, nil
		}
		return map[string]interface{}{"Instances": []interface{}{row("selected", "target", "POST")}, "TotalPage": 2}, nil
	})
	if err := sweepFlinkVpcOwner("cn-beijing", "target", rpc, func(time.Duration) { waits++ }, 3); err != nil {
		t.Fatal(err)
	}
	if deletes != 1 || polls != 2 || waits != 1 {
		t.Fatalf("delete=%d poll=%d wait=%d", deletes, polls, waits)
	}
	if err := flinkVpcOwnerSweeper("cn-beijing", "target", rpc).F("cn-hangzhou"); err == nil {
		t.Fatal("region mismatch accepted")
	}
}

func TestUnitFlinkVpcOwnerRefusesSubscriptionUnknownAndMalformedScope(t *testing.T) {
	for _, charge := range []string{"PRE", "Subscription", "", "UNKNOWN"} {
		writes := 0
		rpc := vpcSweepRPC(func(_, _, action string, _, _ map[string]interface{}, _ bool) (map[string]interface{}, error) {
			if action == "DeleteInstance" {
				writes++
				return nil, nil
			}
			return map[string]interface{}{"Instances": []interface{}{map[string]interface{}{"InstanceId": "one", "VpcId": "target", "ChargeType": charge}}, "TotalPage": 1}, nil
		})
		if err := sweepFlinkVpcOwner("cn-beijing", "target", rpc, func(time.Duration) {}, 2); err == nil || writes != 0 {
			t.Fatalf("billing %q writes=%d err=%v", charge, writes, err)
		}
	}
}

func TestUnitFlinkVpcOwnerPaginationAndReleaseFailuresAreBounded(t *testing.T) {
	for _, mode := range []string{"missing-pagination", "repeated-page", "delete-error", "read-error", "release-timeout"} {
		calls := 0
		writes := 0
		rpc := vpcSweepRPC(func(_, _, action string, q, b map[string]interface{}, _ bool) (map[string]interface{}, error) {
			calls++
			if action == "DeleteInstance" {
				writes++
				if mode == "delete-error" {
					return nil, fmt.Errorf("denied")
				}
				return map[string]interface{}{}, nil
			}
			if mode == "read-error" {
				return nil, fmt.Errorf("read failed")
			}
			response := map[string]interface{}{"Instances": []interface{}{map[string]interface{}{"InstanceId": "one", "VpcId": "target", "ChargeType": "POST"}}, "TotalPage": 1}
			if mode == "missing-pagination" {
				delete(response, "TotalPage")
			}
			if mode == "repeated-page" {
				response["TotalPage"] = 2
			}
			return response, nil
		})
		if err := sweepFlinkVpcOwner("cn-beijing", "target", rpc, func(time.Duration) {}, 2); err == nil {
			t.Fatal(mode, "accepted failure")
		}
		if calls > 4 || writes > 1 {
			t.Fatalf("%s unbounded calls=%d writes=%d", mode, calls, writes)
		}
	}
}

func TestUnitFlinkNativeAdapterRejectsUnrelatedActionsBeforeClientUse(t *testing.T) {
	rpc := nativeVpcFlinkRPC(nil)
	for _, item := range [][3]string{{"foasconsole", "2019-06-01", "DeleteNamespace"}, {"Vpc", "2019-06-01", "DeleteInstance"}, {"foasconsole", "2021-10-28", "DescribeInstances"}} {
		if _, err := rpc(item[0], item[1], item[2], nil, nil, false); err == nil {
			t.Fatal("unsupported adapter call accepted")
		}
	}
}

// vpcScopedSweepers builds native SDK Sweeper nodes bound to one selected VPC.
// F retains the native func(region) error contract; Dependencies express the
// service-owner -> network-interface -> VSwitch -> VPC release order. These
// nodes deliberately do not enter the global, substring-filtered SDK registry.
func vpcScopedSweepers(region, vpcID string, rpc vpcSweepRPC, extraOwnerNodes ...map[string]*resource.Sweeper) map[string]*resource.Sweeper {
	graph := map[string]*resource.Sweeper{}
	add := func(name string, deps []string, run func() error) {
		graph[name] = &resource.Sweeper{Name: name, Dependencies: deps, F: func(actualRegion string) error {
			if actualRegion != region {
				return fmt.Errorf("VPC %s sweeper region mismatch: %s != %s", vpcID, actualRegion, region)
			}
			return run()
		}}
	}
	ownerNames := map[string]string{"polardb": "alicloud_polardb_cluster", "Rds": "alicloud_db_instance", "R-kvstore": "alicloud_kvstore_instance", "Slb": "alicloud_slb_load_balancer", "Ecs": "alicloud_instance", "gpdb": "alicloud_gpdb_instance"}
	var ownerDependencies []string
	for _, owner := range allVpcOwners {
		owner := owner
		name := ownerNames[owner.product]
		ownerDependencies = append(ownerDependencies, name)
		add(name, nil, func() error { return sweepVpcOwner(rpc, region, vpcID, owner) })
	}
	networkExtras := vpcNetworkExtraSweepers(region, vpcID, rpc)
	reservationNode := networkExtras["alicloud_vpc_vswitch_cidr_reservation"]
	delete(networkExtras, "alicloud_vpc_vswitch_cidr_reservation")
	extraOwnerNodes = append([]map[string]*resource.Sweeper{additionalVpcOwnerSweepers(region, vpcID, rpc), networkExtras}, extraOwnerNodes...)
	for _, nodes := range extraOwnerNodes {
		for name, node := range nodes {
			graph[name] = node
			ownerDependencies = append(ownerDependencies, name)
		}
	}
	sort.Strings(ownerDependencies)
	// Custom routes in system tables also block VPC deletion. A router
	// interface can be a route next hop, so remove its routes before the
	// interface while releasing other service owners before their routes.
	routeNode := vpcRouteEntrySweeper(region, vpcID, rpc)
	for _, name := range ownerDependencies {
		if name != "alicloud_router_interface" && name != "alicloud_vpc_ipv4_gateway" {
			routeNode.Dependencies = append(routeNode.Dependencies, name)
		}
	}
	graph[routeNode.Name] = routeNode
	// Router interfaces and IPv4 gateways are route next hops. Their native
	// delete APIs require those routes to be removed first.
	for _, name := range []string{"alicloud_router_interface", "alicloud_vpc_ipv4_gateway"} {
		if node := graph[name]; node != nil {
			node.Dependencies = append(node.Dependencies, routeNode.Name)
		}
	}
	eniDependencies := append([]string{}, ownerDependencies...)
	eniDependencies = append(eniDependencies, "alicloud_nat_gateway", "alicloud_vpn_gateway", "alicloud_havip", "alicloud_vpc_ipv6_gateway")
	definitions := []vpcNetworkSweep{
		{"alicloud_nat_gateway", nil, "Vpc", "2016-04-28", "DescribeNatGateways", "$.NatGateways.NatGateway", "VpcId", "NatGatewayId", "DeleteNatGateway", map[string]interface{}{"Force": true}},
		{"alicloud_vpn_gateway", nil, "Vpc", "2016-04-28", "DescribeVpnGateways", "$.VpnGateways.VpnGateway", "VpcId", "VpnGatewayId", "DeleteVpnGateway", nil},
		{"alicloud_havip", nil, "Vpc", "2016-04-28", "DescribeHaVips", "$.HaVips.HaVip", "VpcId", "HaVipId", "DeleteHaVip", nil},
		{"alicloud_vpc_flow_log", nil, "Vpc", "2016-04-28", "DescribeFlowLogs", "$.FlowLogs.FlowLog", "ResourceId", "FlowLogId", "DeleteFlowLog", nil},
		{"alicloud_vpc_ipv6_gateway", nil, "Vpc", "2016-04-28", "DescribeIpv6Gateways", "$.Ipv6Gateways.Ipv6Gateway", "VpcId", "Ipv6GatewayId", "DeleteIpv6Gateway", nil},
		{"alicloud_ecs_network_interface", eniDependencies, "Ecs", "2014-05-26", "DescribeNetworkInterfaces", "$.NetworkInterfaceSets.NetworkInterfaceSet", "VpcId", "NetworkInterfaceId", "DeleteNetworkInterface", nil},
		{"alicloud_security_group", []string{"alicloud_ecs_network_interface", "alicloud_instance"}, "Ecs", "2014-05-26", "DescribeSecurityGroups", "$.SecurityGroups.SecurityGroup", "VpcId", "SecurityGroupId", "DeleteSecurityGroup", nil},
		{"alicloud_vswitch", []string{"alicloud_ecs_network_interface", "alicloud_vpc_flow_log"}, "Vpc", "2016-04-28", "DescribeVSwitches", "$.VSwitches.VSwitch", "VpcId", "VSwitchId", "DeleteVSwitch", nil},
		{"alicloud_route_table", []string{"alicloud_vswitch", "alicloud_route_entry"}, "Vpc", "2016-04-28", "DescribeRouteTableList", "$.RouterTableList.RouterTableListType", "VpcId", "RouteTableId", "DeleteRouteTable", nil},
		{"alicloud_network_acl", []string{"alicloud_vswitch"}, "Vpc", "2016-04-28", "DescribeNetworkAcls", "$.NetworkAcls.NetworkAcl", "VpcId", "NetworkAclId", "DeleteNetworkAcl", nil},
	}
	for _, definition := range definitions {
		definition := definition
		add(definition.name, definition.dependencies, func() error { return sweepVpcNetworkNode(rpc, region, vpcID, definition) })
	}
	// Prefix reservations may still be held by ENIs. Release interfaces first,
	// then reservations, and only then remove their VSwitch.
	if reservationNode != nil {
		reservationNode.Dependencies = append(reservationNode.Dependencies, "alicloud_ecs_network_interface")
		graph[reservationNode.Name] = reservationNode
		graph["alicloud_vswitch"].Dependencies = append(graph["alicloud_vswitch"].Dependencies, reservationNode.Name)
	}
	var rootDependencies []string
	for name := range graph {
		rootDependencies = append(rootDependencies, name)
	}
	sort.Strings(rootDependencies)
	add("alicloud_vpc", rootDependencies, func() error {
		return deleteVpcScopedResource(rpc, region, "Vpc", "2016-04-28", "DeleteVpc", "VpcId", vpcID, nil)
	})
	return graph
}

type vpcNetworkSweep struct {
	name                                               string
	dependencies                                       []string
	product, version, list, path, scopeKey, id, remove string
	extra                                              map[string]interface{}
}

func sweepVpcNetworkNode(rpc vpcSweepRPC, region, vpcID string, node vpcNetworkSweep) error {
	items, listErr := listVpcScopedNodeResources(rpc, region, vpcID, node)
	var failures []string
	if listErr != nil {
		failures = append(failures, listErr.Error())
	}
	for _, item := range items {
		if node.name == "alicloud_route_table" && item["RouteTableType"] == "System" {
			continue
		}
		id, ok := item[node.id].(string)
		if !ok || id == "" {
			failures = append(failures, fmt.Sprintf("%s returned invalid %s", node.list, node.id))
			continue
		}
		if node.name == "alicloud_ecs_network_interface" {
			if reason := allVpcENISkipReason(item); reason != "" {
				failures = append(failures, fmt.Sprintf("ENI %s: %s", id, reason))
				continue
			}
		}
		if err := deleteVpcScopedResource(rpc, region, node.product, node.version, node.remove, node.id, id, node.extra); err != nil {
			failures = append(failures, err.Error())
			continue
		}
		if err := waitVpcScopedNodeResourceGone(rpc, region, vpcID, node, id, time.Sleep, 120); err != nil {
			failures = append(failures, err.Error())
		}
	}
	if len(failures) > 0 {
		return fmt.Errorf("%s: %s", node.name, strings.Join(failures, "; "))
	}
	return nil
}

// Wait only for the exact identity deleted by this node. Successful empty
// scoped inventories prove release; errors, malformed IDs and unknown ownership
// cannot be interpreted as disappearance.
func waitVpcScopedNodeResourceGone(rpc vpcSweepRPC, region, vpcID string, node vpcNetworkSweep, id string, pause func(time.Duration), attempts int) error {
	if id == "" || region == "" || vpcID == "" || attempts < 1 {
		return fmt.Errorf("%s release wait requires verified identity and bounded attempts", node.name)
	}
	for attempt := 0; attempt < attempts; attempt++ {
		items, err := listVpcScopedNodeResources(rpc, region, vpcID, node, id)
		if err != nil {
			return fmt.Errorf("%s %s release verification: %w", node.name, id, err)
		}
		present := false
		for _, item := range items {
			returnedID, ok := item[node.id].(string)
			if !ok || returnedID == "" {
				return fmt.Errorf("%s malformed %s during release wait", node.name, node.id)
			}
			if returnedID == id {
				present = true
			}
		}
		if !present {
			return nil
		}
		if attempt+1 < attempts {
			pause(5 * time.Second)
		}
	}
	return fmt.Errorf("%s %s still present after bounded release wait", node.name, id)
}

func listVpcScopedNodeResources(rpc vpcSweepRPC, region, vpcID string, node vpcNetworkSweep, watchedID ...string) ([]map[string]interface{}, error) {
	if region == "" || vpcID == "" {
		return nil, fmt.Errorf("%s requires verified region and VPC identity", node.list)
	}
	var items []map[string]interface{}
	var failures []string
	for page := 1; ; page++ {
		response, err := rpc(node.product, node.version, node.list, nil, map[string]interface{}{"RegionId": region, node.scopeKey: vpcID, "PageSize": PageSizeLarge, "PageNumber": page}, true)
		if err != nil {
			return items, fmt.Errorf("%s: %w", node.list, err)
		}
		raw, err := jsonpath.Get(node.path, response)
		if err != nil {
			return items, fmt.Errorf("%s response: %w", node.list, err)
		}
		rows, ok := raw.([]interface{})
		if !ok {
			return items, fmt.Errorf("%s invalid resource list", node.list)
		}
		for _, row := range rows {
			item, ok := row.(map[string]interface{})
			if !ok {
				failures = append(failures, node.list+" malformed resource")
				continue
			}
			// Some network APIs return an unfiltered account inventory even when
			// VpcId is supplied. Explicit other-VPC rows are outside this graph,
			// not failures; absent or malformed ownership still fails closed.
			if node.scopeKey == "VpcId" {
				if actual, known := item["VpcId"].(string); known && actual != "" && actual != vpcID {
					if len(watchedID) > 0 && item[node.id] == watchedID[0] {
						return items, fmt.Errorf("%s %s ownership changed during release verification", node.name, watchedID[0])
					}
					continue
				}
			}
			if include, reason := vpcDependentResourceIdentity(item, node.scopeKey, vpcID); !include {
				failures = append(failures, node.list+": "+reason)
				continue
			}
			items = append(items, item)
		}
		if len(rows) < PageSizeLarge {
			break
		}
	}
	if len(failures) > 0 {
		return items, fmt.Errorf("refused unverified resources: %s", strings.Join(failures, "; "))
	}
	return items, nil
}

func deleteVpcScopedResource(rpc vpcSweepRPC, region, product, version, action, idKey, id string, extra map[string]interface{}) error {
	if id == "" {
		return fmt.Errorf("%s empty resource identity", action)
	}
	request := map[string]interface{}{"RegionId": region, idKey: id}
	for key, value := range extra {
		request[key] = value
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	timeout := 5 * time.Minute
	if action == "DeleteVpc" {
		timeout = 10 * time.Minute
	}
	err := retry.Retry(timeout, func() *retry.RetryError {
		_, err := rpc(product, version, action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("%s %s: %w", action, id, err)
	}
	return nil
}

// Preflight the complete native dependency graph before the first API call.
// Failed nodes do not suppress unrelated cleanup or the final DeleteVpc:
// downstream scoped APIs may safely refuse remaining service dependencies.
func runVpcScopedSweepers(region, vpcID string, graph map[string]*resource.Sweeper) error {
	state := map[string]int{}
	var order []string
	var visit func(string) error
	visit = func(name string) error {
		node, ok := graph[name]
		if !ok || node == nil || node.F == nil || node.Name != name {
			return fmt.Errorf("VPC sweeper dependency %q missing or invalid", name)
		}
		if state[name] == 1 {
			return fmt.Errorf("VPC sweeper dependency cycle at %s", name)
		}
		if state[name] == 2 {
			return nil
		}
		state[name] = 1
		for _, dependency := range node.Dependencies {
			if err := visit(dependency); err != nil {
				return err
			}
		}
		state[name] = 2
		order = append(order, name)
		return nil
	}
	if err := visit("alicloud_vpc"); err != nil {
		return err
	}
	if len(order) != len(graph) {
		return fmt.Errorf("VPC sweeper graph has nodes unreachable from alicloud_vpc")
	}
	failed := map[string]bool{}
	var failures []string
	for _, name := range order {
		node := graph[name]
		for _, dependency := range node.Dependencies {
			if failed[dependency] {
				log.Printf("[WARN] VPC %s sweeper %s dependency %s failed; attempting scoped API with native service constraints", vpcID, name, dependency)
			}
		}
		log.Printf("[INFO] VPC %s native sweeper %s starting", vpcID, name)
		if err := node.F(region); err != nil {
			failed[name] = true
			failures = append(failures, fmt.Sprintf("%s: %v", name, err))
			log.Printf("[ERROR] VPC %s native sweeper %s: %v", vpcID, name, err)
		} else {
			log.Printf("[INFO] VPC %s native sweeper %s completed", vpcID, name)
		}
	}
	if len(failures) > 0 {
		return fmt.Errorf("VPC %s cleanup incomplete: %s", vpcID, strings.Join(failures, "; "))
	}
	return nil
}

func TestUnitVpcScopedSweeperGraphPreflight(t *testing.T) {
	called := 0
	node := func(name string, deps ...string) *resource.Sweeper {
		return &resource.Sweeper{Name: name, Dependencies: deps, F: func(string) error { called++; return nil }}
	}
	cases := []map[string]*resource.Sweeper{
		{"alicloud_vpc": node("alicloud_vpc", "missing")},
		{"alicloud_vpc": node("alicloud_vpc", "owner"), "owner": node("owner", "alicloud_vpc")},
		{"alicloud_vpc": node("alicloud_vpc"), "unreachable": node("unreachable")},
		{"alicloud_vpc": node("wrong-name")},
	}
	for _, graph := range cases {
		if err := runVpcScopedSweepers("cn-beijing", "vpc-target", graph); err == nil {
			t.Fatal("invalid graph accepted")
		}
	}
	if called != 0 {
		t.Fatalf("invalid dependency graph executed %d destructive callbacks", called)
	}
}

func TestUnitVpcScopedSweeperGraphContinuesAfterFailure(t *testing.T) {
	var order []string
	node := func(name string, deps []string, fail bool) *resource.Sweeper {
		return &resource.Sweeper{Name: name, Dependencies: deps, F: func(string) error {
			order = append(order, name)
			if fail {
				return fmt.Errorf("unsupported owner API")
			}
			return nil
		}}
	}
	graph := map[string]*resource.Sweeper{
		"owner":        node("owner", nil, true),
		"network":      node("network", []string{"owner"}, false),
		"other-owner":  node("other-owner", nil, false),
		"alicloud_vpc": node("alicloud_vpc", []string{"network", "other-owner"}, false),
	}
	err := runVpcScopedSweepers("cn-beijing", "vpc-target", graph)
	if err == nil || !strings.Contains(err.Error(), "unsupported owner API") {
		t.Fatalf("dependency error hidden: %v", err)
	}
	if strings.Join(order, ",") != "owner,network,other-owner,alicloud_vpc" {
		t.Fatalf("did not attempt remaining scoped cleanup and final VPC: %v", order)
	}
}

func TestUnitVpcScopedSweeperGraphProtectedBoundary(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	t.Setenv("ALICLOUD_SWEEP_OLDER_THAN", "")
	var selected, deletions []string
	clusterDeleted := false
	eniDeleted, vswitchDeleted := false, false
	rpc := vpcSweepRPC(func(product, version, action string, query, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
		if body["VpcId"] == "vpc-protected" {
			t.Errorf("protected VPC reached scoped API: %s", action)
			return nil, fmt.Errorf("protected VPC")
		}
		if action == "DescribeVpcs" {
			return allVpcSweepTestResponse("$.Vpcs.Vpc", []interface{}{
				map[string]interface{}{"VpcId": "vpc-protected", "VpcName": "cfw-test-no-delete2"},
				map[string]interface{}{"VpcId": "vpc-target", "VpcName": "manual", "IsDefault": true, "CreationTime": time.Now().Format(time.RFC3339)},
			}), nil
		}
		for _, owner := range allVpcOwners {
			if product == owner.product && action == owner.list {
				if product == "gpdb" {
					return nil, fmt.Errorf("gpdb unsupported in region")
				}
				var rows []interface{} = []interface{}{}
				if product == "polardb" {
					rows = append(rows, map[string]interface{}{"DBClusterId": "pc-protected", "VpcId": "vpc-protected"})
					if !clusterDeleted {
						rows = append(rows, map[string]interface{}{"DBClusterId": "pc-target", "VpcId": "vpc-target"})
					}
				}
				return allVpcSweepTestResponse(owner.path, rows), nil
			}
		}
		if action == "ListLoadBalancers" {
			return allVpcSweepTestResponse("$.LoadBalancers", []interface{}{}), nil
		}
		if action == "DescribeVpcAttribute" {
			return map[string]interface{}{"VpcId": "vpc-target", "RegionId": "cn-beijing", "VRouterId": "vrt-target"}, nil
		}
		if action == "DescribeRouterInterfaces" {
			return allVpcSweepTestResponse("$.RouterInterfaceSet.RouterInterfaceType", []interface{}{}), nil
		}
		if action == "DescribeCens" {
			return allVpcSweepTestResponse("$.Cens.Cen", []interface{}{}), nil
		}
		paths := map[string]string{
			"DescribeFileSystems": "$.FileSystems.FileSystem", "ListIpv4Gateways": "$.Ipv4GatewayModels", "ListVSwitchCidrReservations": "$.VSwitchCidrReservations",
			"DescribeNatGateways": "$.NatGateways.NatGateway", "DescribeVpnGateways": "$.VpnGateways.VpnGateway", "DescribeHaVips": "$.HaVips.HaVip", "DescribeFlowLogs": "$.FlowLogs.FlowLog", "DescribeIpv6Gateways": "$.Ipv6Gateways.Ipv6Gateway", "DescribeSecurityGroups": "$.SecurityGroups.SecurityGroup", "DescribeRouteTableList": "$.RouterTableList.RouterTableListType", "DescribeNetworkAcls": "$.NetworkAcls.NetworkAcl",
		}
		if path, ok := paths[action]; ok {
			return allVpcSweepTestResponse(path, []interface{}{}), nil
		}
		if action == "DescribeDBClusterAttribute" {
			return map[string]interface{}{"DBClusterId": body["DBClusterId"], "VPCId": "vpc-target"}, nil
		}
		if action == "DescribeDBClusterEndpoints" {
			return map[string]interface{}{"Items": []interface{}{}}, nil
		}
		switch action {
		case "DescribeNetworkInterfaces":
			if eniDeleted {
				return allVpcSweepTestResponse("$.NetworkInterfaceSets.NetworkInterfaceSet", []interface{}{}), nil
			}
			return allVpcSweepTestResponse("$.NetworkInterfaceSets.NetworkInterfaceSet", []interface{}{
				map[string]interface{}{"VpcId": "vpc-protected", "NetworkInterfaceId": "eni-protected", "ServiceManaged": false, "Status": "Available", "Type": "Secondary"},
				map[string]interface{}{"VpcId": "vpc-target", "NetworkInterfaceId": "eni-target", "ServiceManaged": false, "Status": "Available", "Type": "Secondary"},
			}), nil
		case "DescribeVSwitches":
			if vswitchDeleted {
				return allVpcSweepTestResponse("$.VSwitches.VSwitch", []interface{}{}), nil
			}
			return allVpcSweepTestResponse("$.VSwitches.VSwitch", []interface{}{map[string]interface{}{"VpcId": "vpc-target", "VSwitchId": "vsw-target"}}), nil
		case "DeleteDBCluster":
			if body["DBClusterId"] != "pc-target" {
				t.Fatalf("protected owner deleted: %v", body)
			}
			clusterDeleted = true
		case "DeleteNetworkInterface":
			eniDeleted = true
			if body["NetworkInterfaceId"] != "eni-target" {
				t.Fatalf("protected ENI deleted: %v", body)
			}
		case "DeleteVSwitch":
			vswitchDeleted = true
			if body["VSwitchId"] != "vsw-target" {
				t.Fatalf("protected VSwitch deleted: %v", body)
			}
		case "DeleteVpc":
			if body["VpcId"] != "vpc-target" || body["Force"] != nil {
				t.Fatalf("unsafe root VPC request: %v", body)
			}
		default:
			t.Errorf("unexpected RPC: %s/%s", product, action)
			return nil, fmt.Errorf("unexpected RPC")
		}
		deletions = append(deletions, action)
		return nil, nil
	})
	err := sweepVpcsWithCallbacks("cn-beijing", rpc, func(id string) error {
		selected = append(selected, id)
		return runVpcScopedSweepers("cn-beijing", id, vpcScopedSweepers("cn-beijing", id, rpc))
	})
	if strings.Join(selected, ",") != "vpc-target" {
		t.Fatalf("protected VPC dispatched: %v", selected)
	}
	if err == nil || !strings.Contains(err.Error(), "gpdb unsupported") {
		t.Fatalf("real node errors hidden: %v", err)
	}
	if strings.Join(deletions, ",") != "DeleteDBCluster,DeleteNetworkInterface,DeleteVSwitch,DeleteVpc" {
		t.Fatalf("native dependency deletion order/coverage: %v", deletions)
	}
}

func TestUnitVpcScopedOwnerPartialFailureContinues(t *testing.T) {
	owner := allVpcOwners[0]
	var deleted []string
	rpc := vpcSweepRPC(func(product, version, action string, query, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
		if action == "DescribeDBClusterAttribute" {
			return map[string]interface{}{"DBClusterId": body["DBClusterId"], "VPCId": "vpc-target"}, nil
		}
		if action == "DescribeDBClusterEndpoints" {
			return map[string]interface{}{"Items": []interface{}{}}, nil
		}
		if action == owner.list {
			return allVpcSweepTestResponse(owner.path, []interface{}{
				map[string]interface{}{owner.id: "pc-prepaid", "VpcId": "vpc-target"},
				map[string]interface{}{owner.id: "pc-postpaid", "VpcId": "vpc-target"},
			}), nil
		}
		id := fmt.Sprint(body[owner.id])
		deleted = append(deleted, id)
		if id == "pc-prepaid" {
			return nil, fmt.Errorf("prepaid instance cannot be deleted")
		}
		return nil, nil
	})
	err := deleteAllVpcOwners(rpc, "cn-beijing", "vpc-target", owner)
	if err == nil || !strings.Contains(err.Error(), "prepaid") {
		t.Fatalf("partial failure hidden: %v", err)
	}
	if strings.Join(deleted, ",") != "pc-prepaid,pc-postpaid" {
		t.Fatalf("one owner failure blocked remaining instances: %v", deleted)
	}
}

func TestUnitVpcNetworkReleaseWaitScopeAndBounds(t *testing.T) {
	for _, node := range []vpcNetworkSweep{
		{name: "alicloud_nat_gateway", product: "Vpc", version: "2016-04-28", list: "DescribeNatGateways", path: "$.NatGateways.NatGateway", scopeKey: "VpcId", id: "NatGatewayId"},
		{name: "alicloud_vpn_gateway", product: "Vpc", version: "2016-04-28", list: "DescribeVpnGateways", path: "$.VpnGateways.VpnGateway", scopeKey: "VpcId", id: "VpnGatewayId"},
		{name: "alicloud_vpc_ipv6_gateway", product: "Vpc", version: "2016-04-28", list: "DescribeIpv6Gateways", path: "$.Ipv6Gateways.Ipv6Gateway", scopeKey: "VpcId", id: "Ipv6GatewayId"},
		{name: "alicloud_ecs_network_interface", product: "Ecs", version: "2014-05-26", list: "DescribeNetworkInterfaces", path: "$.NetworkInterfaceSets.NetworkInterfaceSet", scopeKey: "VpcId", id: "NetworkInterfaceId"},
		{name: "alicloud_vswitch", product: "Vpc", version: "2016-04-28", list: "DescribeVSwitches", path: "$.VSwitches.VSwitch", scopeKey: "VpcId", id: "VSwitchId"},
	} {
		t.Run(node.name, func(t *testing.T) {
			item := map[string]interface{}{"VpcId": "vpc-target", node.id: "resource-target"}
			calls, pauses := 0, 0
			rpc := vpcSweepRPC(func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
				calls++
				rows := []interface{}{}
				if calls == 1 {
					rows = append(rows, item)
				}
				return allVpcSweepTestResponse(node.path, rows), nil
			})
			if err := waitVpcScopedNodeResourceGone(rpc, "cn-beijing", "vpc-target", node, "resource-target", func(time.Duration) { pauses++ }, 2); err != nil {
				t.Fatal(err)
			}
			if calls != 2 || pauses != 1 {
				t.Fatalf("asynchronous release not awaited: calls=%d pauses=%d", calls, pauses)
			}
			for _, bad := range []map[string]interface{}{
				{"VpcId": "vpc-target"}, {"VpcId": "vpc-target", node.id: ""}, {"VpcId": "vpc-protected", node.id: "resource-target"}, {node.id: "resource-target"},
			} {
				malformed := vpcSweepRPC(func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
					return allVpcSweepTestResponse(node.path, []interface{}{bad}), nil
				})
				if err := waitVpcScopedNodeResourceGone(malformed, "cn-beijing", "vpc-target", node, "resource-target", func(time.Duration) { t.Fatal("malformed scope polled") }, 1); err == nil {
					t.Fatalf("unknown or changed ownership counted as released: %v", bad)
				}
			}
			persistent := vpcSweepRPC(func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
				return allVpcSweepTestResponse(node.path, []interface{}{item}), nil
			})
			if err := waitVpcScopedNodeResourceGone(persistent, "cn-beijing", "vpc-target", node, "resource-target", func(time.Duration) {}, 2); err == nil {
				t.Fatal("still-present resource counted as released")
			}
			failed := vpcSweepRPC(func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
				return nil, fmt.Errorf("list denied")
			})
			if err := waitVpcScopedNodeResourceGone(failed, "cn-beijing", "vpc-target", node, "resource-target", func(time.Duration) {}, 1); err == nil {
				t.Fatal("list error counted as released")
			}
		})
	}
}

func TestUnitVpcScopedRouteEntryDependencyOrder(t *testing.T) {
	graph := vpcScopedSweepers("cn-beijing", "vpc-target", nil)
	positions := map[string]int{}
	for name, node := range graph {
		name := name
		node.F = func(region string) error {
			if region != "cn-beijing" {
				t.Fatalf("unexpected region %s", region)
			}
			positions[name] = len(positions)
			return nil
		}
	}
	if err := runVpcScopedSweepers("cn-beijing", "vpc-target", graph); err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]string{
		{"alicloud_instance", "alicloud_route_entry"},
		{"alicloud_nas_mount_target", "alicloud_ecs_network_interface"},
		{"alicloud_route_entry", "alicloud_vpc_ipv4_gateway"},
		{"alicloud_ecs_network_interface", "alicloud_vpc_vswitch_cidr_reservation"},
		{"alicloud_vpc_vswitch_cidr_reservation", "alicloud_vswitch"},
		{"alicloud_route_entry", "alicloud_router_interface"},
		{"alicloud_route_entry", "alicloud_route_table"},
		{"alicloud_route_table", "alicloud_vpc"},
	} {
		before, foundBefore := positions[pair[0]]
		after, foundAfter := positions[pair[1]]
		if !foundBefore || !foundAfter || before >= after {
			t.Fatalf("required release order %s before %s was not respected: %v", pair[0], pair[1], positions)
		}
	}
}

func TestUnitVpcNetworkUnfilteredHaVipScope(t *testing.T) {
	node := vpcNetworkSweep{name: "alicloud_havip", product: "Vpc", version: "2016-04-28", list: "DescribeHaVips", path: "$.HaVips.HaVip", scopeKey: "VpcId", id: "HaVipId", remove: "DeleteHaVip"}
	for _, unknown := range []bool{false, true} {
		deleted := false
		rpc := vpcSweepRPC(func(product, version, action string, query, body map[string]interface{}, auto bool) (map[string]interface{}, error) {
			switch action {
			case "DescribeHaVips":
				if body["VpcId"] != "vpc-target" {
					t.Fatal("unscoped list")
				}
				rows := []interface{}{map[string]interface{}{"VpcId": "vpc-protected", "HaVipId": "havip-protected"}}
				if unknown {
					rows = append(rows, map[string]interface{}{"HaVipId": "havip-unknown"})
				}
				if !deleted {
					rows = append(rows, map[string]interface{}{"VpcId": "vpc-target", "HaVipId": "havip-target"})
				}
				return allVpcSweepTestResponse(node.path, rows), nil
			case "DeleteHaVip":
				if body["HaVipId"] != "havip-target" {
					t.Fatal("other or unknown VPC resource deleted")
				}
				deleted = true
				return nil, nil
			default:
				t.Fatalf("unexpected API %s", action)
				return nil, nil
			}
		})
		err := sweepVpcNetworkNode(rpc, "cn-beijing", "vpc-target", node)
		if !deleted || (err != nil) != unknown {
			t.Fatalf("unknown=%v deleted=%v err=%v", unknown, deleted, err)
		}
	}
	called := false
	_, err := listVpcScopedNodeResources(func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
		called = true
		return nil, nil
	}, "cn-beijing", "", node)
	if err == nil || called {
		t.Fatal("unknown target identity reached cloud API")
	}
}

func TestUnitVpcScopedIPv4GatewayWaitsForRouteRemoval(t *testing.T) {
	graph := vpcScopedSweepers("cn-beijing", "vpc-target", nil)
	routesRemoved := false
	gatewayDeleted := false
	for name, node := range graph {
		name := name
		node.F = func(string) error {
			switch name {
			case "alicloud_route_entry":
				routesRemoved = true
			case "alicloud_vpc_ipv4_gateway":
				if !routesRemoved {
					return fmt.Errorf("OperationDenied.RouteEntryExist")
				}
				gatewayDeleted = true
			}
			return nil
		}
	}
	if err := runVpcScopedSweepers("cn-beijing", "vpc-target", graph); err != nil {
		t.Fatal(err)
	}
	if !routesRemoved || !gatewayDeleted {
		t.Fatal("route and IPv4 gateway release must both execute")
	}
}

// Exercise the real GNUmakefile shell branches with a fake go executable. This
// test never recursively invokes Go tests or creates cloud resources.
func TestUnitVpcSweepMakeTimeoutContract(t *testing.T) {
	makePath, err := exec.LookPath("make")
	if err != nil {
		t.Fatal("make required for sweeper contract test:", err)
	}
	source, err := os.ReadFile("../GNUmakefile")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name      string
		variables []string
		timeout   string
		all       bool
	}{
		{"all-default", []string{"SWEEP_ALL=1"}, "3h", true},
		{"ordinary-default", nil, "10m", false},
		{"all-override", []string{"SWEEP_ALL=1", "SWEEP_TIMEOUT=45m"}, "45m", true},
		{"ordinary-override", []string{"SWEEP_TIMEOUT=20m"}, "20m", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			temp := t.TempDir()
			bin := filepath.Join(temp, "bin")
			if err := os.Mkdir(bin, 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(temp, "GNUmakefile"), source, 0600); err != nil {
				t.Fatal(err)
			}
			shim := `#!/bin/sh
set -eu
if [ "$1" = env ] && [ "${2:-}" = GOPATH ]; then printf '%s\n' "$SWEEP_FAKE_GOPATH"; exit 0; fi
if [ "$1" != test ]; then echo 'unexpected fake go command' >&2; exit 81; fi
printf '%s\n' "$@" > "$SWEEP_ARGS_FILE"
printf '%s\n' "${TF_ACC:-}" "${ALICLOUD_SWEEP_ALL_RESOURCES:-}" > "$SWEEP_ENV_FILE"
`
			if err := os.WriteFile(filepath.Join(bin, "go"), []byte(shim), 0700); err != nil {
				t.Fatal(err)
			}
			argsFile := filepath.Join(temp, "args")
			envFile := filepath.Join(temp, "env")
			args := append([]string{"--no-print-directory", "-f", "GNUmakefile", "sweep", "REGION=cn-beijing", "RESOURCE=alicloud_vpc"}, tc.variables...)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, makePath, args...)
			cmd.Dir = temp
			// Explicit minimal environment: never inherit real credentials, MAKEFLAGS,
			// GOPATH, a sweep override, or a real Go executable ahead of our shim.
			cmd.Env = []string{"PATH=" + bin + string(os.PathListSeparator) + os.Getenv("PATH"), "HOME=" + temp, "GOPATH=" + filepath.Join(temp, "fake-gopath"), "SWEEP_FAKE_GOPATH=" + filepath.Join(temp, "fake-gopath"), "SWEEP_ARGS_FILE=" + argsFile, "SWEEP_ENV_FILE=" + envFile}
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("make contract failed: %v: %s", err, output)
			}
			raw, err := os.ReadFile(argsFile)
			if err != nil {
				t.Fatal(err)
			}
			got := strings.Split(strings.TrimSpace(string(raw)), "\n")
			want := []string{"test", "./alicloud", "-v", "-timeout=" + tc.timeout, "-sweep=cn-beijing", "-sweep-run=alicloud_vpc"}
			if strings.Join(got, "\n") != strings.Join(want, "\n") {
				t.Fatalf("actual go argv %q; want %q", got, want)
			}
			env, err := os.ReadFile(envFile)
			if err != nil {
				t.Fatal(err)
			}
			expected := "1\n\n"
			if tc.all {
				expected = "1\ntrue\n"
			}
			if string(env) != expected {
				t.Fatalf("unexpected sweep env %q", env)
			}
		})
	}
}

func vpcNetworkExtraSweepers(region, vpcID string, rpc vpcSweepRPC) map[string]*resource.Sweeper {
	return networkExtraSweepers(region, vpcID, rpc, time.Sleep, 120)
}

func networkExtraSweepers(region, vpcID string, rpc vpcSweepRPC, pause func(time.Duration), attempts int) map[string]*resource.Sweeper {
	nodes := map[string]*resource.Sweeper{}
	add := func(name string, run func() error) {
		nodes[name] = &resource.Sweeper{Name: name, F: func(actual string) error {
			if region == "" || actual != region || vpcID == "" || attempts < 1 {
				return fmt.Errorf("%s requires selected VPC, matching region and bounded wait", name)
			}
			return run()
		}}
	}
	add("alicloud_nas_mount_target", func() error {
		files, err := networkExtraList(rpc, "NAS", "2017-06-26", "DescribeFileSystems", "$.FileSystems.FileSystem", map[string]interface{}{"RegionId": region}, false)
		if err != nil {
			return err
		}
		return serviceVpcOwnerEach(files, func(file map[string]interface{}) error {
			fs, err := serviceVpcOwnerID(file, "FileSystemId")
			if err != nil {
				return err
			}
			list := func() ([]map[string]interface{}, error) {
				return networkExtraList(rpc, "NAS", "2017-06-26", "DescribeMountTargets", "$.MountTargets.MountTarget", map[string]interface{}{"RegionId": region, "FileSystemId": fs}, false)
			}
			mounts, err := list()
			if err != nil {
				return err
			}
			return serviceVpcOwnerEach(mounts, func(mount map[string]interface{}) error {
				match, err := serviceVpcOwnerMatch(mount, "VpcId", vpcID)
				if err != nil || !match {
					return err
				}
				id, err := serviceVpcOwnerID(mount, "MountTargetDomain")
				if err != nil {
					return err
				}
				// Delete only this mount target, never its shared filesystem.
				return networkExtraDeleteWait(rpc, "NAS", "2017-06-26", "DeleteMountTarget", map[string]interface{}{"FileSystemId": fs, "MountTargetDomain": id}, list, "MountTargetDomain", id, "VpcId", vpcID, pause, attempts)
			})
		})
	})
	add("alicloud_vpc_ipv4_gateway", func() error {
		list := func() ([]map[string]interface{}, error) {
			return networkExtraList(rpc, "Vpc", "2016-04-28", "ListIpv4Gateways", "$.Ipv4GatewayModels", map[string]interface{}{"RegionId": region, "VpcId": vpcID}, true)
		}
		rows, err := list()
		if err != nil {
			return err
		}
		return serviceVpcOwnerEach(rows, func(row map[string]interface{}) error {
			match, err := serviceVpcOwnerMatch(row, "VpcId", vpcID)
			if err != nil || !match {
				return err
			}
			id, err := serviceVpcOwnerID(row, "Ipv4GatewayId")
			if err != nil {
				return err
			}
			detail, err := rpc("Vpc", "2016-04-28", "GetIpv4GatewayAttribute", nil, map[string]interface{}{"RegionId": region, "Ipv4GatewayId": id}, true)
			if err != nil {
				return err
			}
			if detail["Ipv4GatewayId"] != id || detail["VpcId"] != vpcID {
				return fmt.Errorf("IPv4 gateway %s detail ownership missing/changed", id)
			}
			return networkExtraDeleteWait(rpc, "Vpc", "2016-04-28", "DeleteIpv4Gateway", map[string]interface{}{"RegionId": region, "Ipv4GatewayId": id, "ClientToken": buildClientToken("DeleteIpv4Gateway")}, list, "Ipv4GatewayId", id, "VpcId", vpcID, pause, attempts)
		})
	})
	add("alicloud_vpc_vswitch_cidr_reservation", func() error {
		switches, err := networkExtraList(rpc, "Vpc", "2016-04-28", "DescribeVSwitches", "$.VSwitches.VSwitch", map[string]interface{}{"RegionId": region, "VpcId": vpcID}, false)
		if err != nil {
			return err
		}
		return serviceVpcOwnerEach(switches, func(vsw map[string]interface{}) error {
			match, err := serviceVpcOwnerMatch(vsw, "VpcId", vpcID)
			if err != nil || !match {
				return err
			}
			vswID, err := serviceVpcOwnerID(vsw, "VSwitchId")
			if err != nil {
				return err
			}
			list := func() ([]map[string]interface{}, error) {
				// The actual reservation response has VSwitchId, but no VpcId.
				// Reverify the parent switch before every scoped read/wait.
				parents, err := networkExtraList(rpc, "Vpc", "2016-04-28", "DescribeVSwitches", "$.VSwitches.VSwitch", map[string]interface{}{"RegionId": region, "VSwitchId": vswID}, false)
				if err != nil {
					return nil, err
				}
				if len(parents) != 1 || parents[0]["VSwitchId"] != vswID || parents[0]["VpcId"] != vpcID {
					return nil, fmt.Errorf("CIDR reservation parent VSwitch %s ownership missing/changed", vswID)
				}
				rows, err := networkExtraList(rpc, "Vpc", "2016-04-28", "ListVSwitchCidrReservations", "$.VSwitchCidrReservations", map[string]interface{}{"RegionId": region, "VSwitchId": vswID}, true)
				if err != nil {
					return nil, err
				}
				for _, row := range rows {
					if row["VSwitchId"] != vswID {
						return nil, fmt.Errorf("CIDR reservation returned wrong/unknown VSwitch scope")
					}
				}
				return rows, nil
			}
			rows, err := list()
			if err != nil {
				return err
			}
			return serviceVpcOwnerEach(rows, func(row map[string]interface{}) error {
				id, err := serviceVpcOwnerID(row, "VSwitchCidrReservationId")
				if err != nil {
					return err
				}
				return networkExtraDeleteWait(rpc, "Vpc", "2016-04-28", "DeleteVSwitchCidrReservation", map[string]interface{}{"RegionId": region, "VSwitchCidrReservationId": id}, list, "VSwitchCidrReservationId", id, "VSwitchId", vswID, pause, attempts)
			})
		})
	})
	return nodes
}

func networkExtraList(rpc vpcSweepRPC, product, version, action, path string, base map[string]interface{}, tokens bool) ([]map[string]interface{}, error) {
	request := map[string]interface{}{}
	for key, value := range base {
		request[key] = value
	}
	if tokens {
		request["MaxResults"] = 100
	} else {
		request["PageSize"], request["PageNumber"] = 50, 1
	}
	seen := map[string]bool{}
	var all []map[string]interface{}
	for page := 0; page < 10000; page++ {
		response, err := rpc(product, version, action, nil, request, true)
		if err != nil {
			return nil, err
		}
		var value interface{}
		if product == "Vpc" && version == "2016-04-28" && action == "ListVSwitchCidrReservations" && path == "$.VSwitchCidrReservations" && request["VSwitchId"] != nil && request["VSwitchId"] != "" && confirmedEmptyVpcCidrReservations(response) {
			value = []interface{}{}
		} else {
			value, err = jsonpath.Get(path, response)
			if err != nil {
				return nil, err
			}
		}
		rows, ok := value.([]interface{})
		if !ok {
			return nil, fmt.Errorf("%s unknown/malformed resource list", action)
		}
		for _, row := range rows {
			item, ok := row.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("%s malformed resource identity", action)
			}
			all = append(all, item)
		}
		if !tokens {
			if len(rows) < 50 {
				return all, nil
			}
			request["PageNumber"] = request["PageNumber"].(int) + 1
			continue
		}
		value, exists := response["NextToken"]
		if !exists || value == nil || value == "" {
			return all, nil
		}
		token, ok := value.(string)
		if !ok || seen[token] {
			return nil, fmt.Errorf("%s invalid/repeated NextToken", action)
		}
		seen[token], request["NextToken"] = true, token
	}
	return nil, fmt.Errorf("%s pagination exceeded bound", action)
}

func networkExtraDeleteWait(rpc vpcSweepRPC, product, version, action string, request map[string]interface{}, list func() ([]map[string]interface{}, error), idKey, id, scopeKey, scope string, pause func(time.Duration), attempts int) error {
	accepted := false
	var err error
	for attempt := 0; attempt < attempts; attempt++ {
		_, err = rpc(product, version, action, nil, request, true)
		if err == nil || rdsErrorHasCode(err, "InvalidIpv4GatewayId.NotFound", "InvalidMountTarget.NotFound") {
			accepted = true
			break
		}
		if !(NeedRetry(err) || rdsErrorHasCode(err, "OperationConflict", "ServiceUnavailable", "SystemBusy", "LastTokenProcessing", "IncorrectStatus.Ipv4Gateway", "IncorrectStatus.Vpc", "IncorrectStatus", "VolumeStatusForbidOperation", "OperationDenied.InvalidState")) {
			break
		}
		if attempt+1 < attempts {
			pause(5 * time.Second)
		}
	}
	if !accepted {
		return fmt.Errorf("%s %s: %v", action, id, err)
	}
	for attempt := 0; attempt < attempts; attempt++ {
		rows, err := list()
		if err != nil {
			return err
		}
		present := false
		for _, row := range rows {
			rowID, ok := row[idKey].(string)
			if !ok || rowID == "" {
				return fmt.Errorf("%s unknown resource identity during wait", action)
			}
			if rowID != id {
				continue
			}
			if row[scopeKey] != scope {
				return fmt.Errorf("%s resource ownership changed during wait", action)
			}
			present = true
		}
		if !present {
			return nil
		}
		if attempt+1 < attempts {
			pause(5 * time.Second)
		}
	}
	return fmt.Errorf("%s %s still present after bounded wait in %s %s", action, id, scopeKey, strings.TrimSpace(scope))
}

// Live ListVSwitchCidrReservations returns only TotalCount=0, MaxResults=0
// and RequestId for an empty set, omitting the optional collection/token.
// Normalize only that API's proven envelope; unknown/missing totals stay errors.
func confirmedEmptyVpcCidrReservations(response map[string]interface{}) bool {
	if value, exists := response["VSwitchCidrReservations"]; exists && value != nil {
		return false
	}
	requestID, ok := response["RequestId"].(string)
	if !ok || requestID == "" {
		return false
	}
	if !networkExtraNumericZero(response["TotalCount"]) {
		return false
	}
	maximum := response["MaxResults"]
	if !networkExtraNumericZero(maximum) && maximum != 100 && maximum != int64(100) && maximum != float64(100) && maximum != json.Number("100") {
		return false
	}
	// The proven success envelope has no error code/message. Never turn an
	// error payload that happens to contain TotalCount=0 into an empty list.
	for _, key := range []string{"Code", "ErrorCode", "ErrorMessage", "Message", "Error"} {
		if value, exists := response[key]; exists && value != nil && value != "" {
			return false
		}
	}
	if success, exists := response["Success"]; exists && success != true {
		return false
	}
	token, exists := response["NextToken"]
	return !exists || token == nil || token == ""
}

func networkExtraNumericZero(value interface{}) bool {
	switch n := value.(type) {
	case int:
		return n == 0
	case int64:
		return n == 0
	case float64:
		return n == 0
	case json.Number:
		// The RPC SDK uses ReadAsJSON/UseNumber. Only the canonical
		// zero integer token is a proven empty count, not a numeric string.
		return n == json.Number("0")
	default:
		return false
	}
}

func TestUnitVpcSweepNetworkExtrasNasPreservesSharedFilesystemAndOtherVpc(t *testing.T) {
	deletes, polls, pauses := 0, 0, 0
	rpc := func(product, version, action string, query, body map[string]interface{}, auto bool) (map[string]interface{}, error) {
		switch action {
		case "DescribeFileSystems":
			return map[string]interface{}{"FileSystems": map[string]interface{}{"FileSystem": []interface{}{map[string]interface{}{"FileSystemId": "shared"}}}}, nil
		case "DescribeMountTargets":
			if body["FileSystemId"] != "shared" {
				t.Fatal("wrong filesystem")
			}
			rows := []interface{}{map[string]interface{}{"MountTargetDomain": "protected", "VpcId": "vpc-protected"}}
			if deletes > 0 {
				polls++
			}
			if deletes == 0 || polls < 2 {
				rows = append(rows, map[string]interface{}{"MountTargetDomain": "target", "VpcId": "vpc-target"})
			}
			return map[string]interface{}{"MountTargets": map[string]interface{}{"MountTarget": rows}}, nil
		case "DeleteMountTarget":
			if body["FileSystemId"] != "shared" || body["MountTargetDomain"] != "target" {
				t.Fatalf("unsafe NAS delete: %v", body)
			}
			deletes++
			return map[string]interface{}{}, nil
		default:
			t.Fatalf("unexpected shared filesystem mutation %s", action)
		}
		return nil, nil
	}
	node := networkExtraSweepers("cn-hangzhou", "vpc-target", rpc, func(time.Duration) { pauses++ }, 3)["alicloud_nas_mount_target"]
	if err := node.F("cn-hangzhou"); err != nil || deletes != 1 || pauses != 1 {
		t.Fatalf("deletes=%d pauses=%d err=%v", deletes, pauses, err)
	}
}

func TestUnitVpcSweepNetworkExtrasIpv4DetailScopeAndAsyncWait(t *testing.T) {
	for _, wrong := range []bool{false, true} {
		deletes := 0
		rpc := func(product, version, action string, query, body map[string]interface{}, auto bool) (map[string]interface{}, error) {
			switch action {
			case "ListIpv4Gateways":
				rows := []interface{}{map[string]interface{}{"Ipv4GatewayId": "protected", "VpcId": "vpc-protected"}}
				if deletes == 0 {
					rows = append(rows, map[string]interface{}{"Ipv4GatewayId": "target", "VpcId": "vpc-target"})
				}
				return map[string]interface{}{"Ipv4GatewayModels": rows}, nil
			case "GetIpv4GatewayAttribute":
				if body["Ipv4GatewayId"] != "target" {
					t.Fatal("unrelated gateway queried")
				}
				vpc := "vpc-target"
				if wrong {
					vpc = "vpc-protected"
				}
				return map[string]interface{}{"Ipv4GatewayId": "target", "VpcId": vpc}, nil
			case "DeleteIpv4Gateway":
				if wrong || body["Ipv4GatewayId"] != "target" || body["ClientToken"] == "" {
					t.Fatalf("unsafe gateway delete: %v", body)
				}
				deletes++
				return map[string]interface{}{}, nil
			}
			return nil, fmt.Errorf("unexpected %s", action)
		}
		err := networkExtraSweepers("cn-hangzhou", "vpc-target", rpc, func(time.Duration) {}, 3)["alicloud_vpc_ipv4_gateway"].F("cn-hangzhou")
		if wrong {
			if err == nil || deletes != 0 {
				t.Fatalf("wrong scope deletes=%d err=%v", deletes, err)
			}
		} else if err != nil || deletes != 1 {
			t.Fatalf("valid scope deletes=%d err=%v", deletes, err)
		}
	}
}

func TestUnitVpcSweepNetworkExtrasReservationUsesVerifiedParentWhenVpcAbsent(t *testing.T) {
	for _, wrong := range []bool{false, true} {
		deletes := 0
		rpc := func(product, version, action string, query, body map[string]interface{}, auto bool) (map[string]interface{}, error) {
			switch action {
			case "DescribeVSwitches":
				if body["PageSize"] != 50 {
					t.Fatalf("invalid VSwitch page size: %v", body)
				}
				vpc := "vpc-target"
				if wrong && body["VSwitchId"] != nil {
					vpc = "vpc-protected"
				}
				return map[string]interface{}{"VSwitches": map[string]interface{}{"VSwitch": []interface{}{map[string]interface{}{"VSwitchId": "vsw", "VpcId": vpc}}}}, nil
			case "ListVSwitchCidrReservations":
				if body["VSwitchId"] != "vsw" {
					t.Fatal("missing reservation scope")
				}
				rows := []interface{}{}
				if deletes == 0 {
					rows = append(rows, map[string]interface{}{"VSwitchId": "vsw", "VSwitchCidrReservationId": "reservation"})
				}
				return map[string]interface{}{"VSwitchCidrReservations": rows}, nil
			case "DeleteVSwitchCidrReservation":
				if wrong || body["VSwitchCidrReservationId"] != "reservation" {
					t.Fatalf("unsafe reservation delete: %v", body)
				}
				deletes++
				return map[string]interface{}{}, nil
			}
			return nil, fmt.Errorf("unexpected %s", action)
		}
		err := networkExtraSweepers("cn-hangzhou", "vpc-target", rpc, func(time.Duration) {}, 3)["alicloud_vpc_vswitch_cidr_reservation"].F("cn-hangzhou")
		if wrong {
			if err == nil || deletes != 0 {
				t.Fatalf("wrong parent deletes=%d err=%v", deletes, err)
			}
		} else if err != nil || deletes != 1 {
			t.Fatalf("valid parent deletes=%d err=%v", deletes, err)
		}
	}
}

func TestUnitVpcSweepNetworkExtrasStrictPaginationAndUnknownLists(t *testing.T) {
	calls := 0
	rpc := func(product, version, action string, query, body map[string]interface{}, auto bool) (map[string]interface{}, error) {
		calls++
		count := 1
		if calls == 1 {
			count = 50
		}
		rows := make([]interface{}, count)
		for i := range rows {
			rows[i] = map[string]interface{}{"ID": "known"}
		}
		if body["PageNumber"] != calls {
			t.Fatalf("page did not advance: %v", body)
		}
		return map[string]interface{}{"Rows": rows}, nil
	}
	rows, err := networkExtraList(rpc, "NAS", "2017-06-26", "DescribeFileSystems", "$.Rows", nil, false)
	if err != nil || len(rows) != 51 || calls != 2 {
		t.Fatalf("rows=%d calls=%d err=%v", len(rows), calls, err)
	}
	for _, scenario := range []string{"repeated-token", "bad-token", "unknown-list", "read-error"} {
		rpc = func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
			if scenario == "read-error" {
				return nil, fmt.Errorf("permission denied")
			}
			response := map[string]interface{}{"Rows": []interface{}{}}
			switch scenario {
			case "repeated-token":
				response["NextToken"] = "same"
			case "bad-token":
				response["NextToken"] = 123
			case "unknown-list":
				response["Rows"] = nil
			}
			return response, nil
		}
		if _, err := networkExtraList(rpc, "Vpc", "2016-04-28", "ListIpv4Gateways", "$.Rows", nil, true); err == nil {
			t.Fatalf("%s accepted", scenario)
		}
	}
}

func TestUnitVpcSweepNetworkExtrasRawTeaRetryAndBoundedWait(t *testing.T) {
	for _, scenario := range []string{"retry", "denied", "wait-error", "unknown-scope", "timeout"} {
		deletes, pauses := 0, 0
		rpc := func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
			deletes++
			if scenario == "retry" && deletes == 1 {
				return nil, &tea.SDKError{Code: tea.String("IncorrectStatus.Ipv4Gateway"), Message: tea.String("retry"), StatusCode: tea.Int(400), Data: nil}
			}
			if scenario == "denied" {
				return nil, &tea.SDKError{Code: tea.String("Forbidden.RAM"), Message: tea.String("denied"), StatusCode: tea.Int(403), Data: nil}
			}
			return map[string]interface{}{}, nil
		}
		list := func() ([]map[string]interface{}, error) {
			if scenario == "wait-error" {
				return nil, fmt.Errorf("permission denied")
			}
			if scenario == "timeout" {
				return []map[string]interface{}{{"ID": "id", "VpcId": "vpc-target"}}, nil
			}
			if scenario == "unknown-scope" {
				return []map[string]interface{}{{"ID": "id"}}, nil
			}
			return []map[string]interface{}{}, nil
		}
		err := networkExtraDeleteWait(rpc, "Vpc", "2016-04-28", "DeleteIpv4Gateway", nil, list, "ID", "id", "VpcId", "vpc-target", func(time.Duration) { pauses++ }, 3)
		if scenario == "retry" {
			if err != nil || deletes != 2 || pauses != 1 {
				t.Fatalf("retry deletes=%d pauses=%d err=%v", deletes, pauses, err)
			}
		} else if err == nil {
			t.Fatalf("%s falsely succeeded", scenario)
		}
		if scenario == "timeout" && pauses != 2 {
			t.Fatalf("wait not bounded: %d", pauses)
		}
		if scenario == "denied" && (pauses != 0 || deletes != 1) {
			t.Fatal("denied retried")
		}
	}
}

func TestUnitVpcSweepNetworkExtrasCidrKnownOmittedEmptyEnvelope(t *testing.T) {
	for _, explicitNull := range []bool{false, true} {
		rpc := func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
			response := map[string]interface{}{"TotalCount": float64(0), "MaxResults": float64(0), "RequestId": "request"}
			if explicitNull {
				response["VSwitchCidrReservations"] = nil
			}
			return response, nil
		}
		rows, err := networkExtraList(rpc, "Vpc", "2016-04-28", "ListVSwitchCidrReservations", "$.VSwitchCidrReservations", map[string]interface{}{"VSwitchId": "vsw"}, true)
		if err != nil || len(rows) != 0 {
			t.Fatalf("real empty envelope rejected: %v", err)
		}
	}
	if !confirmedEmptyVpcCidrReservations(map[string]interface{}{"TotalCount": 0, "MaxResults": 100, "RequestId": "request"}) {
		t.Fatal("valid echoed page size rejected")
	}
	for _, scenario := range []string{"missing-total", "nonzero-total", "string-total", "missing-max", "unexpected-max", "missing-request", "has-token", "malformed-list", "other-api", "error-code", "code", "success-false", "error-message"} {
		rpc := func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
			response := map[string]interface{}{"TotalCount": 0, "MaxResults": 0, "RequestId": "request"}
			switch scenario {
			case "missing-total":
				delete(response, "TotalCount")
			case "nonzero-total":
				response["TotalCount"] = 1
			case "string-total":
				response["TotalCount"] = "0"
			case "missing-max":
				delete(response, "MaxResults")
			case "unexpected-max":
				response["MaxResults"] = 1
			case "missing-request":
				delete(response, "RequestId")
			case "has-token":
				response["NextToken"] = "more"
			case "error-code":
				response["ErrorCode"] = "Forbidden"
			case "code":
				response["Code"] = "403"
			case "success-false":
				response["Success"] = false
			case "error-message":
				response["ErrorMessage"] = "denied"
			case "malformed-list":
				response["VSwitchCidrReservations"] = "unknown"
			}
			return response, nil
		}
		action := "ListVSwitchCidrReservations"
		if scenario == "other-api" {
			action = "ListIpv4Gateways"
		}
		if _, err := networkExtraList(rpc, "Vpc", "2016-04-28", action, "$.VSwitchCidrReservations", map[string]interface{}{"VSwitchId": "vsw"}, true); err == nil {
			t.Fatalf("%s falsely accepted as empty", scenario)
		}
	}
}

func TestUnitVpcSweepNetworkExtrasCidrUsesProductionRpcNumberMapping(t *testing.T) {
	for _, maximum := range []string{"0", "100"} {
		decoded, err := util.ReadAsJSON(strings.NewReader(`{"TotalCount":0,"MaxResults":` + maximum + `,"RequestId":"fixture"}`))
		if err != nil {
			t.Fatal(err)
		}
		// tea-rpc ReadAsJSON -> AssertAsMap preserves json.Number in the response.
		response := util.AssertAsMap(decoded)
		if _, ok := response["TotalCount"].(json.Number); !ok {
			t.Fatalf("unexpected native count type %T", response["TotalCount"])
		}
		if _, ok := response["MaxResults"].(json.Number); !ok {
			t.Fatalf("unexpected native maximum type %T", response["MaxResults"])
		}
		rpc := func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
			return response, nil
		}
		rows, err := networkExtraList(rpc, "Vpc", "2016-04-28", "ListVSwitchCidrReservations", "$.VSwitchCidrReservations", map[string]interface{}{"VSwitchId": "vsw"}, true)
		if err != nil || len(rows) != 0 {
			t.Fatalf("SDK empty response max=%s rows=%d err=%v", maximum, len(rows), err)
		}
		response["Success"] = false
		if confirmedEmptyVpcCidrReservations(response) {
			t.Fatal("SDK error response accepted")
		}
	}
}

func TestUnitVpcSweepNetworkExtrasCidrRejectsNoncanonicalNumbers(t *testing.T) {
	for _, value := range []interface{}{json.Number("-0"), json.Number("-1"), json.Number("0.0"), json.Number("0e0"), json.Number("+0"), json.Number("00"), json.Number(""), json.Number("NaN"), "0", true, nil} {
		if networkExtraNumericZero(value) {
			t.Fatalf("invalid zero accepted: %v (%T)", value, value)
		}
	}
	for _, maximum := range []interface{}{json.Number("100.0"), json.Number("1e2"), json.Number("+100"), json.Number("0100"), json.Number("1000000"), "100", nil} {
		response := map[string]interface{}{"TotalCount": json.Number("0"), "MaxResults": maximum, "RequestId": "fixture"}
		if confirmedEmptyVpcCidrReservations(response) {
			t.Fatalf("invalid maximum accepted: %v (%T)", maximum, maximum)
		}
	}
	for _, response := range []map[string]interface{}{
		{"MaxResults": json.Number("0"), "RequestId": "fixture"},
		{"TotalCount": json.Number("0"), "RequestId": "fixture"},
		{"TotalCount": json.Number("1"), "MaxResults": json.Number("0"), "RequestId": "fixture"},
		{"TotalCount": json.Number("0"), "MaxResults": json.Number("0"), "RequestId": "fixture", "ErrorCode": "Forbidden"},
	} {
		if confirmedEmptyVpcCidrReservations(response) {
			t.Fatal("unknown/error SDK envelope accepted")
		}
	}
}

// Additional owners release their own VPC attachments, never a shared global
// CEN/transit router or resources belonging to a different (protected) VPC.
func additionalVpcOwnerSweepers(region, vpcID string, rpc vpcSweepRPC) map[string]*resource.Sweeper {
	nodes := map[string]*resource.Sweeper{}
	add := func(name string, f func() error) {
		nodes[name] = &resource.Sweeper{Name: name, F: func(actual string) error {
			if actual != region {
				return fmt.Errorf("%s region mismatch", name)
			}
			return f()
		}}
	}
	add("alicloud_alb_load_balancer", func() error { return sweepVpcModernLoadBalancers(rpc, region, vpcID, "Alb", "2020-06-16") })
	add("alicloud_nlb_load_balancer", func() error { return sweepVpcModernLoadBalancers(rpc, region, vpcID, "Nlb", "2022-04-30") })
	add("alicloud_cen_instance_attachment", func() error { return sweepVpcCenAttachments(rpc, region, vpcID, false) })
	add("alicloud_cen_transit_router_vpc_attachment", func() error { return sweepVpcCenAttachments(rpc, region, vpcID, true) })
	add("alicloud_router_interface", func() error { return sweepVpcRouterInterfaces(rpc, region, vpcID) })
	return nodes
}

func vpcOwnerError(label string, failures []string) error {
	if len(failures) > 0 {
		return fmt.Errorf("%s: %s", label, strings.Join(failures, "; "))
	}
	return nil
}

func vpcRPCList(rpc vpcSweepRPC, product, version, action, path string, base map[string]interface{}, tokens bool) ([]map[string]interface{}, error) {
	var items []map[string]interface{}
	request := map[string]interface{}{}
	for key, value := range base {
		request[key] = value
	}
	if tokens {
		request["MaxResults"] = PageSizeLarge
	} else {
		request["PageSize"] = PageSizeLarge
		request["PageNumber"] = 1
	}
	seenTokens := map[string]bool{}
	for {
		response, err := rpc(product, version, action, nil, request, true)
		if err != nil {
			return items, err
		}
		raw, err := jsonpath.Get(path, response)
		if err != nil {
			return items, fmt.Errorf("%s invalid list response: %w", action, err)
		}
		rows, ok := raw.([]interface{})
		if !ok {
			return items, fmt.Errorf("%s resource list missing", action)
		}
		for _, row := range rows {
			item, ok := row.(map[string]interface{})
			if !ok {
				return items, fmt.Errorf("%s malformed resource", action)
			}
			items = append(items, item)
		}
		if tokens {
			token, _ := response["NextToken"].(string)
			if token == "" {
				break
			}
			if seenTokens[token] {
				return items, fmt.Errorf("%s repeated NextToken", action)
			}
			seenTokens[token] = true
			request["NextToken"] = token
		} else {
			if len(rows) < PageSizeLarge {
				break
			}
			request["PageNumber"] = request["PageNumber"].(int) + 1
		}
	}
	return items, nil
}

func sweepVpcModernLoadBalancers(rpc vpcSweepRPC, region, vpcID, product, version string) error {
	rows, err := vpcRPCList(rpc, product, version, "ListLoadBalancers", "$.LoadBalancers", map[string]interface{}{"RegionId": region}, true)
	if err != nil {
		return fmt.Errorf("%s ListLoadBalancers: %w", product, err)
	}
	var failures []string
	for _, item := range rows {
		actual, ok := item["VpcId"].(string)
		if !ok {
			failures = append(failures, "load balancer missing VpcId")
			continue
		}
		if actual != vpcID {
			continue
		}
		id, ok := item["LoadBalancerId"].(string)
		if !ok || id == "" {
			failures = append(failures, "load balancer missing ID")
			continue
		}
		detail, err := rpc(product, version, "GetLoadBalancerAttribute", nil, map[string]interface{}{"RegionId": region, "LoadBalancerId": id}, true)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s detail: %v", id, err))
			continue
		}
		if detail["VpcId"] != vpcID {
			failures = append(failures, id+" detail VPC ownership mismatch")
			continue
		}
		config, ok := detail["DeletionProtectionConfig"].(map[string]interface{})
		if !ok {
			failures = append(failures, id+" deletion protection unknown")
			continue
		}
		enabled, ok := config["Enabled"].(bool)
		if !ok {
			failures = append(failures, id+" deletion protection flag unknown")
			continue
		}
		if enabled {
			action := "DisableDeletionProtection"
			request := map[string]interface{}{"RegionId": region, "ResourceId": id, "ClientToken": buildClientToken(action)}
			if product == "Nlb" {
				action = "UpdateLoadBalancerProtection"
				request = map[string]interface{}{"RegionId": region, "LoadBalancerId": id, "DeletionProtectionEnabled": false, "ClientToken": buildClientToken(action)}
			}
			if _, err := rpc(product, version, action, nil, request, true); err != nil {
				failures = append(failures, fmt.Sprintf("%s %s: %v", action, id, err))
				continue
			}
		}
		if enabled {
			err := retry.Retry(5*time.Minute, func() *retry.RetryError {
				current, err := rpc(product, version, "GetLoadBalancerAttribute", nil, map[string]interface{}{"RegionId": region, "LoadBalancerId": id}, true)
				if err != nil {
					return retry.NonRetryableError(err)
				}
				if current["VpcId"] != vpcID {
					return retry.NonRetryableError(fmt.Errorf("%s %s VPC ownership changed", product, id))
				}
				config, ok := current["DeletionProtectionConfig"].(map[string]interface{})
				if !ok {
					return retry.NonRetryableError(fmt.Errorf("%s %s deletion protection unknown", product, id))
				}
				if config["Enabled"] == false {
					return nil
				}
				time.Sleep(3 * time.Second)
				return retry.RetryableError(fmt.Errorf("%s %s deletion protection still enabled", product, id))
			})
			if err != nil {
				failures = append(failures, err.Error())
				continue
			}
		}
		if err := deleteVpcScopedResource(rpc, region, product, version, "DeleteLoadBalancer", "LoadBalancerId", id, map[string]interface{}{"ClientToken": buildClientToken("DeleteLoadBalancer")}); err != nil {
			failures = append(failures, err.Error())
			continue
		}
		err = retry.Retry(10*time.Minute, func() *retry.RetryError {
			detail, err := rpc(product, version, "GetLoadBalancerAttribute", nil, map[string]interface{}{"RegionId": region, "LoadBalancerId": id}, true)
			if err != nil {
				if NotFoundError(err) || rdsErrorHasCode(err, "ResourceNotFound.LoadBalancer", "ResourceNotFound") {
					return nil
				}
				return retry.NonRetryableError(err)
			}
			if detail["LoadBalancerStatus"] == "Deleted" {
				return nil
			}
			time.Sleep(3 * time.Second)
			return retry.RetryableError(fmt.Errorf("%s %s still releasing VPC dependencies", product, id))
		})
		if err != nil {
			failures = append(failures, err.Error())
		}
	}
	return vpcOwnerError(product, failures)
}

// Only the local VRouter can establish VPC ownership. VBR and opposite-router
// fields describe different resources and must never authorize a mutation.
func vpcRouterInterfaceOwned(item map[string]interface{}, vpcID, routerID string) (bool, error) {
	kind, ok := item["RouterType"].(string)
	if !ok || (kind != "VRouter" && kind != "VBR") {
		return false, fmt.Errorf("router interface has unknown RouterType")
	}
	router, ok := item["RouterId"].(string)
	if !ok || router == "" {
		return false, fmt.Errorf("router interface missing RouterId")
	}
	if kind == "VBR" {
		return false, nil
	}
	if router != routerID {
		if item["VpcInstanceId"] == vpcID {
			return false, fmt.Errorf("router interface VPC and router ownership disagree")
		}
		return false, nil
	}
	if item["VpcInstanceId"] != vpcID {
		return false, fmt.Errorf("router interface VPC ownership changed or missing")
	}
	return true, nil
}

func sweepVpcRouterInterfaces(rpc vpcSweepRPC, region, vpcID string) error {
	detail, err := rpc("Vpc", "2016-04-28", "DescribeVpcAttribute", nil, map[string]interface{}{"RegionId": region, "VpcId": vpcID}, false)
	if err != nil {
		return err
	}
	routerID, ok := detail["VRouterId"].(string)
	if detail["VpcId"] != vpcID || detail["RegionId"] != region || !ok || routerID == "" {
		return fmt.Errorf("router interface requires verified VPC VRouterId")
	}

	rows, err := vpcRPCList(rpc, "Vpc", "2016-04-28", "DescribeRouterInterfaces", "$.RouterInterfaceSet.RouterInterfaceType", map[string]interface{}{"RegionId": region}, false)
	if err != nil {
		return err
	}
	var failures []string
	for _, item := range rows {
		owned, scopeErr := vpcRouterInterfaceOwned(item, vpcID, routerID)
		if scopeErr != nil {
			failures = append(failures, scopeErr.Error())
			continue
		}
		if !owned {
			continue
		}
		id, ok := item["RouterInterfaceId"].(string)
		if !ok || id == "" {
			failures = append(failures, "router interface missing ID")
			continue
		}
		if item["Status"] == "Active" {
			if _, err := rpc("Vpc", "2016-04-28", "DeactivateRouterInterface", nil, map[string]interface{}{"RegionId": region, "RouterInterfaceId": id}, true); err != nil {
				failures = append(failures, fmt.Sprintf("deactivate %s: %v", id, err))
				continue
			}
		}
		present, err := waitVpcRouterInterfaceInactive(rpc, region, vpcID, routerID, id, time.Sleep, 60)
		if err != nil {
			failures = append(failures, err.Error())
			continue
		}
		if !present {
			continue
		}
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			_, err := rpc("Vpc", "2016-04-28", "DeleteRouterInterface", nil, map[string]interface{}{"RegionId": region, "RouterInterfaceId": id, "ClientToken": buildClientToken("DeleteRouterInterface")}, false)
			if err != nil {
				if retryableVpcRouterInterfaceDelete(err) {
					time.Sleep(5 * time.Second)
					return retry.RetryableError(err)
				}
				if rdsErrorHasCode(err, "InvalidInstanceId.NotFound") {
					return nil
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			failures = append(failures, fmt.Sprintf("DeleteRouterInterface %s: %v", id, err))
		}
	}
	return vpcOwnerError("router interfaces", failures)
}

func retryableVpcRouterInterfaceDelete(err error) bool {
	return NeedRetry(err) || rdsErrorHasCode(err, "IncorrectStatus", "InvalidInstance.StatusError", "DependencyViolation.RouterInterfaceReferedByRouteEntry", "OperationConflict")
}

// The native API deactivates asynchronously. Re-read the exact interface and
// its VPC before deletion; never deactivate/delete its opposite interface.
func waitVpcRouterInterfaceInactive(rpc vpcSweepRPC, region, vpcID, routerID, id string, pause func(time.Duration), attempts int) (bool, error) {
	if region == "" || vpcID == "" || routerID == "" || id == "" || attempts < 1 {
		return false, fmt.Errorf("router interface wait requires bounded verified scope")
	}
	for attempt := 0; attempt < attempts; attempt++ {
		rows, err := vpcRPCList(rpc, "Vpc", "2016-04-28", "DescribeRouterInterfaces", "$.RouterInterfaceSet.RouterInterfaceType", map[string]interface{}{"RegionId": region, "Filter.1.Key": "RouterInterfaceId", "Filter.1.Value.1": id}, false)
		if err != nil {
			return false, err
		}
		present := false
		for _, item := range rows {
			returned, ok := item["RouterInterfaceId"].(string)
			if !ok || returned == "" {
				return false, fmt.Errorf("router interface wait returned unverifiable ID")
			}
			if returned != id {
				continue
			}
			if owned, scopeErr := vpcRouterInterfaceOwned(item, vpcID, routerID); scopeErr != nil || !owned {
				return false, fmt.Errorf("router interface %s ownership changed or unverifiable: %v", id, scopeErr)
			}
			switch item["Status"] {
			case string(Inactive), string(Idle):
				return true, nil
			case string(Active), "Deactivating":
				present = true
			default:
				return false, fmt.Errorf("router interface %s unknown status %v", id, item["Status"])
			}
		}
		if !present {
			return false, nil
		}
		if attempt+1 < attempts {
			pause(5 * time.Second)
		}
	}
	return false, fmt.Errorf("router interface %s still Active/Deactivating after bounded deactivation wait", id)
}

// ListTransitRouters is the native regional inventory operation. A verified
// empty result proves there is no TR in this CEN/region; attachment-list errors
// such as IllegalParam.CenIdOrRegionId are never broadly swallowed.
func vpcCenRegionalTransitRouters(rpc vpcSweepRPC, region, cenID string) ([]map[string]interface{}, error) {
	routers, err := vpcRPCList(rpc, "Cbn", "2017-09-12", "ListTransitRouters", "$.TransitRouters", map[string]interface{}{"CenId": cenID, "RegionId": region}, false)
	if err != nil {
		return nil, err
	}
	for _, router := range routers {
		id, ok := router["TransitRouterId"].(string)
		if !ok || id == "" || router["CenId"] != cenID || router["RegionId"] != region {
			return nil, fmt.Errorf("regional transit router identity/CEN/region is unverifiable")
		}
	}
	return routers, nil
}

func vpcCenRegionalRouterKind(routers []map[string]interface{}) (string, error) {
	if len(routers) == 0 {
		return "None", nil
	}
	kind := ""
	for _, router := range routers {
		value, ok := router["Type"].(string)
		if !ok || (value != "Basic" && value != "Enterprise") {
			return "", fmt.Errorf("CEN regional transit router Type is unknown")
		}
		if kind != "" && kind != value {
			return "", fmt.Errorf("CEN regional transit router types are mixed; refusing ambiguous operation")
		}
		kind = value
	}
	return kind, nil
}

func vpcCenRouterKindForNode(rpc vpcSweepRPC, region, cenID string, cen map[string]interface{}) (string, error) {
	if kind, ok := cen["_vpcSweepRegionalTransitRouterType"].(string); ok {
		if kind == "None" || kind == "Basic" || kind == "Enterprise" {
			return kind, nil
		}
		return "", fmt.Errorf("CEN %s candidate regional router type is unknown", cenID)
	}
	if present, known := cen["_vpcSweepRegionalTransitRouterPresent"].(bool); known && !present {
		return "None", nil
	}
	routers, err := vpcCenRegionalTransitRouters(rpc, region, cenID)
	if err != nil {
		return "", err
	}
	return vpcCenRegionalRouterKind(routers)
}

func sweepVpcCenAttachments(rpc vpcSweepRPC, region, vpcID string, transit bool) error {
	cens, err := vpcRPCList(rpc, "Cbn", "2017-09-12", "DescribeCens", "$.Cens.Cen", map[string]interface{}{"RegionId": region}, false)
	if err != nil {
		return err
	}
	var failures []string
	for _, cen := range cens {
		if issue, ok := cen["_vpcSweepInventoryError"].(string); ok && issue != "" {
			failures = append(failures, issue)
		}
		cenID, ok := cen["CenId"].(string)
		if !ok || cenID == "" {
			failures = append(failures, "CEN missing ID")
			continue
		}
		if transit {
			kind, err := vpcCenRouterKindForNode(rpc, region, cenID, cen)
			if err != nil {
				failures = append(failures, err.Error())
				continue
			}
			if kind != "Enterprise" {
				continue
			}

			base := map[string]interface{}{"RegionId": region, "CenId": cenID, "VpcId": vpcID}
			rows, err := vpcRPCList(rpc, "Cbn", "2017-09-12", "ListTransitRouterVpcAttachments", "$.TransitRouterAttachments", base, true)
			if err != nil {
				failures = append(failures, fmt.Sprintf("CEN %s TR list: %v", cenID, err))
				continue
			}
			for _, item := range rows {
				if item["VpcId"] != vpcID {
					continue
				}
				if item["VpcRegionId"] != region || item["CenId"] != cenID {
					failures = append(failures, "TR attachment region/CEN ownership mismatch")
					continue
				}
				id, ok := item["TransitRouterAttachmentId"].(string)
				if !ok || id == "" {
					failures = append(failures, "TR attachment missing ID")
					continue
				}
				// Force is the native attachment's force_delete option. It applies only
				// to this verified attachment, never to the shared CEN/transit router.
				if err := deleteVpcScopedResource(rpc, region, "Cbn", "2017-09-12", "DeleteTransitRouterVpcAttachment", "TransitRouterAttachmentId", id, map[string]interface{}{"Force": true, "ClientToken": buildClientToken("DeleteTransitRouterVpcAttachment")}); err != nil {
					failures = append(failures, err.Error())
				} else if err := waitVpcCenAttachmentGone(rpc, region, vpcID, cenID, id, true); err != nil {
					failures = append(failures, err.Error())
				}
			}
		} else {
			rows, err := vpcRPCList(rpc, "Cbn", "2017-09-12", "DescribeCenAttachedChildInstances", "$.ChildInstances.ChildInstance", map[string]interface{}{"CenId": cenID, "ChildInstanceType": "VPC", "ChildInstanceRegionId": region}, false)
			if err != nil {
				failures = append(failures, fmt.Sprintf("CEN %s children: %v", cenID, err))
				continue
			}
			for _, item := range rows {
				if item["ChildInstanceId"] != vpcID {
					continue
				}
				if item["ChildInstanceType"] != "VPC" || item["ChildInstanceRegionId"] != region {
					failures = append(failures, "CEN child region/type mismatch")
					continue
				}
				kind, err := vpcCenRouterKindForNode(rpc, region, cenID, cen)
				if err != nil {
					failures = append(failures, err.Error())
					continue
				}
				// Enterprise TR child rows are released only by the separate
				// verified VPC attachment node. Classic Detach rejects them.
				if kind == "Enterprise" {
					continue
				}
				_, err = rpc("Cbn", "2017-09-12", "DetachCenChildInstance", nil, map[string]interface{}{"CenId": cenID, "ChildInstanceId": vpcID, "ChildInstanceType": "VPC", "ChildInstanceRegionId": region}, true)
				if err != nil {
					failures = append(failures, fmt.Sprintf("detach CEN %s: %v", cenID, err))
				} else if err := waitVpcCenAttachmentGone(rpc, region, vpcID, cenID, vpcID, false); err != nil {
					failures = append(failures, err.Error())
				}
			}
		}
	}
	return vpcOwnerError("selected VPC CEN attachments", failures)
}

func waitVpcCenAttachmentGone(rpc vpcSweepRPC, region, vpcID, cenID, id string, transit bool) error {
	return waitVpcCenAttachmentGoneWithWait(rpc, region, vpcID, cenID, id, transit, time.Sleep, 200)
}

func waitVpcCenAttachmentGoneWithWait(rpc vpcSweepRPC, region, vpcID, cenID, id string, transit bool, pause func(time.Duration), attempts int) error {
	if region == "" || vpcID == "" || cenID == "" || id == "" || attempts < 1 {
		return fmt.Errorf("CEN release wait requires verified scope and bounded attempts")
	}
	for attempt := 0; attempt < attempts; attempt++ {
		action, path, key := "DescribeCenAttachedChildInstances", "$.ChildInstances.ChildInstance", "ChildInstanceId"
		base := map[string]interface{}{"CenId": cenID, "ChildInstanceType": "VPC", "ChildInstanceRegionId": region}
		if transit {
			action, path, key = "ListTransitRouterVpcAttachments", "$.TransitRouterAttachments", "TransitRouterAttachmentId"
			base = map[string]interface{}{"CenId": cenID, "RegionId": region, "VpcId": vpcID}
		}
		rows, err := vpcRPCList(rpc, "Cbn", "2017-09-12", action, path, base, transit)
		if err != nil {
			return err
		}
		present := false
		for _, item := range rows {
			returnedID, ok := item[key].(string)
			if !ok || returnedID == "" {
				return fmt.Errorf("CEN %s malformed %s during release wait", cenID, key)
			}
			if returnedID != id {
				continue
			}
			if item["CenId"] != cenID {
				return fmt.Errorf("CEN %s attachment %s has unknown or changed CEN identity", cenID, id)
			}
			if transit {
				if item["VpcId"] != vpcID || item["VpcRegionId"] != region {
					return fmt.Errorf("CEN %s TR attachment %s changed VPC/region", cenID, id)
				}
			} else if item["ChildInstanceId"] != vpcID || item["ChildInstanceType"] != "VPC" || item["ChildInstanceRegionId"] != region {
				return fmt.Errorf("CEN %s child %s changed VPC/region/type", cenID, id)
			}
			present = true
		}
		if !present {
			return nil
		}
		if attempt+1 < attempts {
			pause(3 * time.Second)
		}
	}
	return fmt.Errorf("CEN %s attachment %s still present after bounded release wait", cenID, id)
}

// ACK is a ROA service; keep it separate from the RPC transport instead of
// inventing an RpcPost action. Adapters use the same SDK as the native provider.
type vpcACKAPI struct {
	list              func(int) ([]map[string]interface{}, error)
	detail            func(string) (map[string]interface{}, error)
	disableProtection func(string) error
	remove            func(string) error
	// Optional wait seams keep lifecycle regressions pure and bounded.
	pause           func(time.Duration)
	now             func() time.Time
	initialAttempts int
}

func verifyVpcACKDetail(detail map[string]interface{}, region, vpcID, id string) error {
	if detail["cluster_id"] != id || detail["vpc_id"] != vpcID || detail["region_id"] != region {
		return fmt.Errorf("ACK %s detail identity/VPC/region unknown or changed", id)
	}
	// These are the state values documented by the provider's native CS SDK.
	// Only initial is waited here; normal DeleteCluster decides other states.
	switch detail["state"] {
	case "initial", "failed", "running", "updating", "upgrading", "removing", "draining", "scaling", "inactive", "unavailable", "deleting", "deleted", "delete_failed", "waiting", "disconnected":
		return nil
	default:
		return fmt.Errorf("ACK %s unknown state %v", id, detail["state"])
	}
}

func waitVpcACKInitial(api vpcACKAPI, region, vpcID, id string) (map[string]interface{}, error) {
	pause := api.pause
	if pause == nil {
		pause = time.Sleep
	}
	attempts := api.initialAttempts
	if attempts == 0 {
		attempts = 120
	}
	if attempts < 1 || region == "" || vpcID == "" || id == "" {
		return nil, fmt.Errorf("ACK initial wait requires verified scope and bounded attempts")
	}
	now := api.now
	if now == nil {
		now = time.Now
	}
	deadline := now().Add(20 * time.Minute)
	for attempt := 0; attempt < attempts; attempt++ {
		if !now().Before(deadline) {
			break
		}
		detail, err := api.detail(id)
		if err != nil {
			if NotFoundError(err) || rdsErrorHasCode(err, "ErrClusterNotFound", "ClusterNotFound") {
				return nil, nil
			}
			return nil, fmt.Errorf("ACK %s initial wait read: %w", id, err)
		}
		if err := verifyVpcACKDetail(detail, region, vpcID, id); err != nil {
			return nil, err
		}
		if detail["state"] == "deleted" {
			return nil, nil
		}
		if detail["state"] != "initial" {
			return detail, nil
		}
		if attempt+1 < attempts {
			remaining := deadline.Sub(now())
			if remaining <= 0 {
				break
			}
			delay := 10 * time.Second
			if remaining < delay {
				delay = remaining
			}
			pause(delay)
		}
	}
	return nil, fmt.Errorf("ACK %s remains initial after bounded creation wait; cluster and VPC dependencies remain", id)
}

func vpcACKSweeper(region, vpcID string, api vpcACKAPI) *resource.Sweeper {
	return &resource.Sweeper{Name: "alicloud_cs_kubernetes", F: func(actual string) error {
		if actual != region {
			return fmt.Errorf("ACK region mismatch")
		}
		var clusters []map[string]interface{}
		for page := 1; ; page++ {
			rows, err := api.list(page)
			if err != nil {
				return err
			}
			clusters = append(clusters, rows...)
			if len(rows) < PageSizeLarge {
				break
			}
		}
		var failures []string
		for _, cluster := range clusters {
			if cluster["vpc_id"] != vpcID {
				continue
			}
			if cluster["region_id"] != region {
				failures = append(failures, "ACK cluster region mismatch")
				continue
			}
			id, ok := cluster["cluster_id"].(string)
			if !ok || id == "" {
				failures = append(failures, "ACK cluster missing ID")
				continue
			}
			detail, err := waitVpcACKInitial(api, region, vpcID, id)
			if err != nil {
				failures = append(failures, err.Error())
				continue
			}
			if detail == nil {
				continue
			}
			protection, ok := detail["deletion_protection"].(bool)
			if !ok {
				failures = append(failures, id+" ACK deletion protection unknown")
				continue
			}
			if protection {
				if err := api.disableProtection(id); err != nil {
					failures = append(failures, err.Error())
					continue
				}
				err := retry.Retry(5*time.Minute, func() *retry.RetryError {
					current, err := api.detail(id)
					if err != nil {
						if NotFoundError(err) || rdsErrorHasCode(err, "ErrClusterNotFound", "ClusterNotFound") {
							return nil
						}
						return retry.NonRetryableError(err)
					}
					if err := verifyVpcACKDetail(current, region, vpcID, id); err != nil {
						return retry.NonRetryableError(err)
					}
					if current["deletion_protection"] == false {
						return nil
					}
					time.Sleep(3 * time.Second)
					return retry.RetryableError(fmt.Errorf("ACK %s deletion protection still enabled", id))
				})
				if err != nil {
					failures = append(failures, err.Error())
					continue
				}
			}
			if protection {
				current, err := waitVpcACKInitial(api, region, vpcID, id)
				if err != nil {
					failures = append(failures, err.Error())
					continue
				}
				if current == nil {
					continue
				}
			}
			if err := api.remove(id); err != nil {
				failures = append(failures, err.Error())
				continue
			}
			err = retry.Retry(20*time.Minute, func() *retry.RetryError {
				current, err := api.detail(id)
				if err != nil {
					if NotFoundError(err) || rdsErrorHasCode(err, "ErrClusterNotFound", "ClusterNotFound") {
						return nil
					}
					return retry.NonRetryableError(err)
				}
				if err := verifyVpcACKDetail(current, region, vpcID, id); err != nil {
					return retry.NonRetryableError(err)
				}
				if current["state"] == "deleted" {
					return nil
				}
				time.Sleep(5 * time.Second)
				return retry.RetryableError(fmt.Errorf("ACK %s still releasing dependencies", id))
			})
			if err != nil {
				failures = append(failures, err.Error())
			}
		}
		return vpcOwnerError("ACK", failures)
	}}
}

func nativeVpcACKAPI(client *connectivity.AliyunClient, region string) vpcACKAPI {
	var once sync.Once
	var csClient *roacs.Client
	var initErr error
	get := func() (*roacs.Client, error) {
		once.Do(func() { csClient, initErr = client.NewRoaCsClient() })
		return csClient, initErr
	}
	bodyMap := func(body interface{}) (map[string]interface{}, error) {
		raw, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		var result map[string]interface{}
		err = json.Unmarshal(raw, &result)
		return result, err
	}
	return vpcACKAPI{
		list: func(page int) ([]map[string]interface{}, error) {
			sdk, err := get()
			if err != nil {
				return nil, err
			}
			response, err := sdk.DescribeClustersV1((&roacs.DescribeClustersV1Request{}).SetRegionId(region).SetPageNumber(int64(page)).SetPageSize(int64(PageSizeLarge)))
			if err != nil {
				return nil, err
			}
			if response == nil || response.Body == nil {
				return nil, fmt.Errorf("ACK empty cluster list response")
			}
			var rows []map[string]interface{}
			for _, cluster := range response.Body.Clusters {
				item, err := bodyMap(cluster)
				if err != nil {
					return nil, err
				}
				rows = append(rows, item)
			}
			return rows, nil
		},
		detail: func(id string) (map[string]interface{}, error) {
			sdk, err := get()
			if err != nil {
				return nil, err
			}
			response, err := sdk.DescribeClusterDetail(tea.String(id))
			if err != nil {
				return nil, err
			}
			if response == nil || response.Body == nil {
				return nil, fmt.Errorf("ACK empty cluster detail")
			}
			return bodyMap(response.Body)
		},
		disableProtection: func(id string) error {
			sdk, err := get()
			if err != nil {
				return err
			}
			_, err = sdk.ModifyCluster(tea.String(id), (&roacs.ModifyClusterRequest{}).SetDeletionProtection(false))
			return err
		},
		remove: func(id string) error {
			sdk, err := get()
			if err != nil {
				return err
			}
			_, err = sdk.DeleteCluster(tea.String(id), &roacs.DeleteClusterRequest{})
			return err
		},
	}
}

// newVpcCenInventory lives only for one region's sweepVpcs invocation and one
// credential-bound RPC. It caches read-only candidate hints, never deletion
// confirmation. Every scoped attachment list/write/wait still uses fresh RPC.
func newVpcCenInventory(region string, rpc vpcSweepRPC) func(string) vpcSweepRPC {
	var once sync.Once
	var inventoryErr error
	var cens []map[string]interface{}
	byVpc := map[string]map[string]bool{}
	fallback := map[string]bool{}
	inventoryIssues := map[string]string{}
	regionalTRPresence := map[string]bool{}
	regionalTRTypes := map[string]string{}
	record := func(vpcID, cenID string) {
		if byVpc[vpcID] == nil {
			byVpc[vpcID] = map[string]bool{}
		}
		byVpc[vpcID][cenID] = true
	}
	load := func() {
		cens, inventoryErr = vpcRPCList(rpc, "Cbn", "2017-09-12", "DescribeCens", "$.Cens.Cen", map[string]interface{}{"RegionId": region}, false)
		if inventoryErr != nil {
			return
		}
		for _, cen := range cens {
			cenID, ok := cen["CenId"].(string)
			if !ok || cenID == "" {
				inventoryErr = fmt.Errorf("CEN inventory returned missing CenId")
				return
			}
			children, err := vpcRPCList(rpc, "Cbn", "2017-09-12", "DescribeCenAttachedChildInstances", "$.ChildInstances.ChildInstance", map[string]interface{}{"CenId": cenID, "ChildInstanceType": "VPC", "ChildInstanceRegionId": region}, false)
			if err != nil {
				fallback[cenID] = true
				inventoryIssues[cenID] = fmt.Sprintf("CEN %s candidate inventory: %v", cenID, err)
				log.Printf("[WARN] CEN %s candidate inventory unavailable: %v; retaining for every selected VPC", cenID, err)
			} else {
				for _, child := range children {
					kind, ok := child["ChildInstanceType"].(string)
					if !ok || kind == "" {
						fallback[cenID] = true
						continue
					}
					if kind != "VPC" {
						continue
					}
					childRegion, ok := child["ChildInstanceRegionId"].(string)
					if !ok || childRegion == "" {
						fallback[cenID] = true
						continue
					}
					if childRegion != region {
						continue
					}
					vpcID, ok := child["ChildInstanceId"].(string)
					if !ok || vpcID == "" {
						fallback[cenID] = true
						continue
					}
					if child["CenId"] != nil && child["CenId"] != cenID {
						fallback[cenID] = true
						continue
					}
					record(vpcID, cenID)
				}
			}
			routers, err := vpcCenRegionalTransitRouters(rpc, region, cenID)
			if err != nil {
				fallback[cenID] = true
				inventoryIssues[cenID] = fmt.Sprintf("CEN %s regional TR inventory: %v", cenID, err)
				log.Printf("[WARN] %s; retaining for every selected VPC", inventoryIssues[cenID])
				continue
			}
			kind, err := vpcCenRegionalRouterKind(routers)
			if err != nil {
				fallback[cenID] = true
				inventoryIssues[cenID] = fmt.Sprintf("CEN %s regional TR type: %v", cenID, err)
				continue
			}
			regionalTRPresence[cenID] = len(routers) > 0
			regionalTRTypes[cenID] = kind
			if kind != "Enterprise" {
				if fallback[cenID] && inventoryIssues[cenID] == "" {
					inventoryIssues[cenID] = "CEN " + cenID + " candidate inventory contains unverifiable ownership"
				}
				continue
			}
			attachments, err := vpcRPCList(rpc, "Cbn", "2017-09-12", "ListTransitRouterVpcAttachments", "$.TransitRouterAttachments", map[string]interface{}{"CenId": cenID, "RegionId": region}, true)
			if err != nil {
				fallback[cenID] = true
				inventoryIssues[cenID] = fmt.Sprintf("CEN %s TR candidate inventory: %v", cenID, err)
				log.Printf("[WARN] CEN %s TR candidate inventory unavailable: %v; retaining for every selected VPC", cenID, err)
				continue
			}
			for _, attachment := range attachments {
				vpcID, ok := attachment["VpcId"].(string)
				if !ok || vpcID == "" {
					fallback[cenID] = true
					continue
				}
				attachmentRegion, ok := attachment["VpcRegionId"].(string)
				if !ok || attachmentRegion == "" {
					fallback[cenID] = true
					continue
				}
				if attachmentRegion != region {
					continue
				}
				if attachment["CenId"] != nil && attachment["CenId"] != cenID {
					fallback[cenID] = true
					continue
				}
				record(vpcID, cenID)
			}
			if fallback[cenID] {
				if inventoryIssues[cenID] == "" {
					inventoryIssues[cenID] = "CEN " + cenID + " candidate inventory contains unverifiable ownership"
				}
				log.Printf("[WARN] CEN %s has unavailable or unverifiable candidate ownership; retaining for every selected VPC", cenID)
			}
		}
	}
	return func(vpcID string) vpcSweepRPC {
		return func(product, version, action string, query, body map[string]interface{}, auto bool) (map[string]interface{}, error) {
			if product != "Cbn" || version != "2017-09-12" || action != "DescribeCens" {
				return rpc(product, version, action, query, body, auto)
			}
			if vpcID == "" {
				return nil, fmt.Errorf("CEN candidate inventory requires selected VPC")
			}
			once.Do(load)
			if inventoryErr != nil {
				return nil, inventoryErr
			}
			candidates := make([]interface{}, 0)
			for _, cen := range cens {
				cenID := cen["CenId"].(string)
				if fallback[cenID] || byVpc[vpcID][cenID] {
					candidate := map[string]interface{}{}
					for key, value := range cen {
						candidate[key] = value
					}
					if present, known := regionalTRPresence[cenID]; known {
						candidate["_vpcSweepRegionalTransitRouterPresent"] = present
					}
					if kind, known := regionalTRTypes[cenID]; known {
						candidate["_vpcSweepRegionalTransitRouterType"] = kind
					}
					if issue := inventoryIssues[cenID]; issue != "" {
						candidate["_vpcSweepInventoryError"] = issue
					}
					candidates = append(candidates, candidate)
				}
			}
			page, size := 1, PageSizeLarge
			if value, ok := body["PageNumber"].(int); ok {
				page = value
			}
			if value, ok := body["PageSize"].(int); ok {
				size = value
			}
			if page < 1 || size < 1 {
				return nil, fmt.Errorf("CEN candidate inventory invalid pagination")
			}
			start := (page - 1) * size
			if start > len(candidates) {
				start = len(candidates)
			}
			end := start + size
			if end > len(candidates) {
				end = len(candidates)
			}
			return map[string]interface{}{"Cens": map[string]interface{}{"Cen": candidates[start:end]}, "TotalCount": len(candidates), "PageNumber": page, "PageSize": size}, nil
		}
	}
}

func TestUnitVpcModernLoadBalancerOwnershipAndProtection(t *testing.T) {
	for _, product := range []string{"Alb", "Nlb"} {
		t.Run(product, func(t *testing.T) {
			version := "2020-06-16"
			if product == "Nlb" {
				version = "2022-04-30"
			}
			removed, disabled := false, false
			var mutations []string
			rpc := vpcSweepRPC(func(p, v, action string, query, body map[string]interface{}, auto bool) (map[string]interface{}, error) {
				if p != product || v != version {
					t.Fatalf("wrong API %s/%s", p, v)
				}
				switch action {
				case "ListLoadBalancers":
					if body["NextToken"] == nil {
						return map[string]interface{}{"NextToken": "next", "LoadBalancers": []interface{}{map[string]interface{}{"VpcId": "vpc-protected", "LoadBalancerId": "lb-protected"}}}, nil
					}
					return allVpcSweepTestResponse("$.LoadBalancers", []interface{}{map[string]interface{}{"VpcId": "vpc-target", "LoadBalancerId": "lb-target"}}), nil
				case "GetLoadBalancerAttribute":
					if body["LoadBalancerId"] != "lb-target" {
						t.Fatal("protected LB detail touched")
					}
					if removed {
						return map[string]interface{}{"LoadBalancerStatus": "Deleted"}, nil
					}
					return map[string]interface{}{"VpcId": "vpc-target", "DeletionProtectionConfig": map[string]interface{}{"Enabled": !disabled}}, nil
				case "DisableDeletionProtection":
					if product != "Alb" || body["ResourceId"] != "lb-target" {
						t.Fatalf("unsafe ALB protection request: %v", body)
					}
					disabled = true
				case "UpdateLoadBalancerProtection":
					if product != "Nlb" || body["LoadBalancerId"] != "lb-target" || body["DeletionProtectionEnabled"] != false {
						t.Fatalf("unsafe NLB protection request: %v", body)
					}
					disabled = true
				case "DeleteLoadBalancer":
					if !disabled || body["LoadBalancerId"] != "lb-target" {
						t.Fatalf("LB deleted before scoped protection release: %v", body)
					}
					removed = true
				default:
					t.Fatalf("unexpected API %s", action)
				}
				mutations = append(mutations, action)
				return nil, nil
			})
			if err := sweepVpcModernLoadBalancers(rpc, "cn-beijing", "vpc-target", product, version); err != nil {
				t.Fatal(err)
			}
			if len(mutations) != 2 || mutations[1] != "DeleteLoadBalancer" {
				t.Fatalf("expected protection release then deletion: %v", mutations)
			}
		})
	}
}

func TestUnitVpcModernLoadBalancerRefusesDetailMismatch(t *testing.T) {
	rpc := vpcSweepRPC(func(product, version, action string, query, body map[string]interface{}, auto bool) (map[string]interface{}, error) {
		if action == "ListLoadBalancers" {
			return allVpcSweepTestResponse("$.LoadBalancers", []interface{}{map[string]interface{}{"VpcId": "vpc-target", "LoadBalancerId": "lb-target"}}), nil
		}
		if action == "GetLoadBalancerAttribute" {
			return map[string]interface{}{"VpcId": "vpc-protected"}, nil
		}
		t.Fatalf("protected LB mutated: %s", action)
		return nil, nil
	})
	if err := sweepVpcModernLoadBalancers(rpc, "cn-beijing", "vpc-target", "Alb", "2020-06-16"); err == nil {
		t.Fatal("detail ownership mismatch accepted")
	}
}

func TestUnitVpcAckOwnerProtectionAndBoundary(t *testing.T) {
	disabled, removed := false, false
	api := vpcACKAPI{
		list: func(page int) ([]map[string]interface{}, error) {
			return []map[string]interface{}{
				{"cluster_id": "c-protected", "vpc_id": "vpc-protected", "region_id": "cn-beijing"},
				{"cluster_id": "c-target", "vpc_id": "vpc-target", "region_id": "cn-beijing"},
			}, nil
		},
		detail: func(id string) (map[string]interface{}, error) {
			if id != "c-target" {
				t.Fatal("protected ACK queried for deletion")
			}
			if removed {
				return nil, &tea.SDKError{StatusCode: tea.Int(404), Message: tea.String("cluster not found")}
			}
			return map[string]interface{}{"cluster_id": "c-target", "state": "running", "vpc_id": "vpc-target", "region_id": "cn-beijing", "deletion_protection": !disabled}, nil
		},
		disableProtection: func(id string) error {
			if id != "c-target" {
				t.Fatal("protected ACK protection changed")
			}
			disabled = true
			return nil
		},
		remove: func(id string) error {
			if id != "c-target" || !disabled {
				t.Fatal("unsafe ACK delete")
			}
			removed = true
			return nil
		},
	}
	if err := vpcACKSweeper("cn-beijing", "vpc-target", api).F("cn-beijing"); err != nil {
		t.Fatal(err)
	}
	if !removed {
		t.Fatal("selected ACK was not deleted")
	}
	api.detail = func(string) (map[string]interface{}, error) {
		return map[string]interface{}{"vpc_id": "vpc-protected", "region_id": "cn-beijing", "deletion_protection": true}, nil
	}
	api.disableProtection = func(string) error { t.Fatal("mismatched ACK protection changed"); return nil }
	api.remove = func(string) error { t.Fatal("mismatched ACK deleted"); return nil }
	if err := vpcACKSweeper("cn-beijing", "vpc-target", api).F("cn-beijing"); err == nil {
		t.Fatal("ACK detail mismatch accepted")
	}
}

func TestUnitVpcCenSharedEntitiesRemain(t *testing.T) {
	detached, transitDeleted := false, false
	transitMode := false
	var mutations []string
	rpc := vpcSweepRPC(func(product, version, action string, query, body map[string]interface{}, auto bool) (map[string]interface{}, error) {
		if product != "Cbn" || version != "2017-09-12" {
			t.Fatalf("wrong CEN API %s/%s", product, version)
		}
		switch action {
		case "DescribeCens":
			return allVpcSweepTestResponse("$.Cens.Cen", []interface{}{map[string]interface{}{"CenId": "cen-shared"}}), nil
		case "DescribeCenAttachedChildInstances":
			rows := []interface{}{map[string]interface{}{"ChildInstanceId": "vpc-protected", "ChildInstanceType": "VPC", "ChildInstanceRegionId": "cn-beijing"}}
			if !detached {
				rows = append(rows, map[string]interface{}{"ChildInstanceId": "vpc-target", "ChildInstanceType": "VPC", "ChildInstanceRegionId": "cn-beijing"})
			}
			return allVpcSweepTestResponse("$.ChildInstances.ChildInstance", rows), nil
		case "ListTransitRouters":
			kind := "Basic"
			if transitMode {
				kind = "Enterprise"
			}
			return allVpcSweepTestResponse("$.TransitRouters", []interface{}{map[string]interface{}{"CenId": body["CenId"], "RegionId": "cn-beijing", "TransitRouterId": "tr-regional", "Type": kind}}), nil
		case "ListTransitRouterVpcAttachments":
			rows := []interface{}{map[string]interface{}{"VpcId": "vpc-protected", "VpcRegionId": "cn-beijing", "CenId": "cen-shared", "TransitRouterAttachmentId": "tr-protected"}}
			if !transitDeleted {
				rows = append(rows, map[string]interface{}{"VpcId": "vpc-target", "VpcRegionId": "cn-beijing", "CenId": "cen-shared", "TransitRouterAttachmentId": "tr-target"})
			}
			return allVpcSweepTestResponse("$.TransitRouterAttachments", rows), nil
		case "DetachCenChildInstance":
			if body["ChildInstanceId"] != "vpc-target" || body["CenId"] != "cen-shared" {
				t.Fatalf("unsafe CEN detach %v", body)
			}
			detached = true
		case "DeleteTransitRouterVpcAttachment":
			if body["TransitRouterAttachmentId"] != "tr-target" {
				t.Fatalf("protected attachment deleted %v", body)
			}
			transitDeleted = true
		default:
			t.Fatalf("shared global CEN/TR entity touched: %s", action)
		}
		mutations = append(mutations, action)
		return nil, nil
	})
	if err := sweepVpcCenAttachments(rpc, "cn-beijing", "vpc-target", false); err != nil {
		t.Fatal(err)
	}
	transitMode = true
	if err := sweepVpcCenAttachments(rpc, "cn-beijing", "vpc-target", true); err != nil {
		t.Fatal(err)
	}
	if strings.Join(mutations, ",") != "DetachCenChildInstance,DeleteTransitRouterVpcAttachment" {
		t.Fatalf("wrong shared CEN mutation scope %v", mutations)
	}
}

func TestUnitVpcRouterInterfaceExactOwnership(t *testing.T) {
	var actions []string
	deactivated := false
	rpc := vpcSweepRPC(func(product, version, action string, query, body map[string]interface{}, auto bool) (map[string]interface{}, error) {
		if action == "DescribeVpcAttribute" {
			return map[string]interface{}{"VpcId": "vpc-target", "RegionId": "cn-beijing", "VRouterId": "vrt-target"}, nil
		}
		if action == "DescribeRouterInterfaces" {
			if deactivated {
				return allVpcSweepTestResponse("$.RouterInterfaceSet.RouterInterfaceType", []interface{}{map[string]interface{}{"RouterType": "VRouter", "RouterId": "vrt-target", "VpcInstanceId": "vpc-target", "RouterInterfaceId": "ri-target", "Status": "Inactive"}}), nil
			}
			return allVpcSweepTestResponse("$.RouterInterfaceSet.RouterInterfaceType", []interface{}{
				map[string]interface{}{"RouterType": "VBR", "RouterId": "vbr-other", "RouterInterfaceId": "ri-vbr", "Status": "Active"},
				map[string]interface{}{"RouterType": "VRouter", "RouterId": "vrt-protected", "VpcInstanceId": "vpc-protected", "RouterInterfaceId": "ri-protected", "Status": "Active"},
				map[string]interface{}{"RouterType": "VRouter", "RouterId": "vrt-target", "VpcInstanceId": "vpc-target", "RouterInterfaceId": "ri-target", "Status": "Active"},
			}), nil
		}
		if body["RouterInterfaceId"] != "ri-target" {
			t.Fatalf("protected RI touched %v", body)
		}
		actions = append(actions, action)
		if action == "DeactivateRouterInterface" {
			deactivated = true
		}
		return nil, nil
	})
	if err := sweepVpcRouterInterfaces(rpc, "cn-beijing", "vpc-target"); err != nil {
		t.Fatal(err)
	}
	if strings.Join(actions, ",") != "DeactivateRouterInterface,DeleteRouterInterface" {
		t.Fatalf("RI lifecycle wrong %v", actions)
	}
}

func TestUnitVpcAdditionalOwnerDependencyRegistration(t *testing.T) {
	rpc := vpcSweepRPC(func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
		t.Fatal("graph construction called API")
		return nil, nil
	})
	extra := map[string]*resource.Sweeper{"alicloud_extra_owner": {Name: "alicloud_extra_owner", F: func(string) error { return nil }}}
	graph := vpcScopedSweepers("cn-beijing", "vpc-target", rpc, extra)
	deps := strings.Join(graph["alicloud_ecs_network_interface"].Dependencies, ",")
	for _, name := range []string{"alicloud_alb_load_balancer", "alicloud_nlb_load_balancer", "alicloud_cen_instance_attachment", "alicloud_cen_transit_router_vpc_attachment", "alicloud_router_interface", "alicloud_extra_owner"} {
		if graph[name] == nil || !strings.Contains(deps, name) {
			t.Fatalf("owner not wired before ENI: %s", name)
		}
	}
	if err := graph["alicloud_extra_owner"].F("cn-beijing"); err != nil {
		t.Fatal(fmt.Sprint(err))
	}
}

func TestUnitVpcCenReleaseWaitRequiresIdentity(t *testing.T) {
	for _, transit := range []bool{false, true} {
		t.Run(fmt.Sprintf("transit=%v", transit), func(t *testing.T) {
			path, key, id := "$.ChildInstances.ChildInstance", "ChildInstanceId", "vpc-target"
			item := map[string]interface{}{"CenId": "cen-shared", "ChildInstanceId": "vpc-target", "ChildInstanceType": "VPC", "ChildInstanceRegionId": "cn-beijing"}
			if transit {
				path, key, id = "$.TransitRouterAttachments", "TransitRouterAttachmentId", "tr-target"
				item = map[string]interface{}{"CenId": "cen-shared", "VpcId": "vpc-target", "VpcRegionId": "cn-beijing", "TransitRouterAttachmentId": "tr-target"}
			}
			calls, pauses := 0, 0
			rpc := vpcSweepRPC(func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
				calls++
				rows := []interface{}{}
				if calls == 1 {
					rows = append(rows, item)
				}
				return allVpcSweepTestResponse(path, rows), nil
			})
			if err := waitVpcCenAttachmentGoneWithWait(rpc, "cn-beijing", "vpc-target", "cen-shared", id, transit, func(time.Duration) { pauses++ }, 2); err != nil {
				t.Fatal(err)
			}
			if calls != 2 || pauses != 1 {
				t.Fatalf("did not confirm asynchronous release: calls=%d pauses=%d", calls, pauses)
			}
			for _, bad := range []map[string]interface{}{
				{}, {key: ""}, {key: id, "CenId": "cen-other"}, {key: id, "CenId": "cen-shared"},
			} {
				malformed := vpcSweepRPC(func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
					return allVpcSweepTestResponse(path, []interface{}{bad}), nil
				})
				if err := waitVpcCenAttachmentGoneWithWait(malformed, "cn-beijing", "vpc-target", "cen-shared", id, transit, func(time.Duration) { t.Fatal("malformed identity polled") }, 1); err == nil {
					t.Fatalf("unknown/mismatched identity counted as gone: %v", bad)
				}
			}
			persistent := vpcSweepRPC(func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
				return allVpcSweepTestResponse(path, []interface{}{item}), nil
			})
			if err := waitVpcCenAttachmentGoneWithWait(persistent, "cn-beijing", "vpc-target", "cen-shared", id, transit, func(time.Duration) {}, 2); err == nil {
				t.Fatal("still-present attachment counted as gone")
			}
			failed := vpcSweepRPC(func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
				return nil, fmt.Errorf("list denied")
			})
			if err := waitVpcCenAttachmentGoneWithWait(failed, "cn-beijing", "vpc-target", "cen-shared", id, transit, func(time.Duration) {}, 1); err == nil {
				t.Fatal("list error counted as gone")
			}
		})
	}
}

func TestUnitVpcCenInventorySharedOnceFreshRelease(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	t.Setenv("ALICLOUD_SWEEP_OLDER_THAN", "")
	var mu sync.Mutex
	globalLists, trLists, regionalTRLists := 0, 0, 0
	childLists := map[string]int{}
	detached := map[string]bool{}
	var mutations []string
	rpc := vpcSweepRPC(func(product, version, action string, query, body map[string]interface{}, auto bool) (map[string]interface{}, error) {
		mu.Lock()
		defer mu.Unlock()
		switch action {
		case "DescribeVpcs":
			return allVpcSweepTestResponse("$.Vpcs.Vpc", []interface{}{
				map[string]interface{}{"VpcId": "vpc-a", "VpcName": "manual-a"}, map[string]interface{}{"VpcId": "vpc-b", "VpcName": "manual-b"}, map[string]interface{}{"VpcId": "vpc-protected", "VpcName": "default-NODELETING"},
			}), nil
		case "DescribeCens":
			globalLists++
			return allVpcSweepTestResponse("$.Cens.Cen", []interface{}{map[string]interface{}{"CenId": "cen-a"}, map[string]interface{}{"CenId": "cen-b"}}), nil
		case "DescribeCenAttachedChildInstances":
			cen := fmt.Sprint(body["CenId"])
			childLists[cen]++
			id := "vpc-a"
			if cen == "cen-b" {
				id = "vpc-b"
			}
			rows := []interface{}{}
			if !detached[id] {
				rows = append(rows, map[string]interface{}{"CenId": cen, "ChildInstanceId": id, "ChildInstanceType": "VPC", "ChildInstanceRegionId": "cn-beijing"})
			}
			if cen == "cen-a" {
				rows = append(rows, map[string]interface{}{"CenId": cen, "ChildInstanceId": "vpc-protected", "ChildInstanceType": "VPC", "ChildInstanceRegionId": "cn-beijing"})
			}
			return allVpcSweepTestResponse("$.ChildInstances.ChildInstance", rows), nil
		case "ListTransitRouters":
			regionalTRLists++
			return allVpcSweepTestResponse("$.TransitRouters", []interface{}{map[string]interface{}{"CenId": body["CenId"], "RegionId": "cn-beijing", "TransitRouterId": "tr-regional", "Type": "Basic"}}), nil
		case "ListTransitRouterVpcAttachments":
			trLists++
			return allVpcSweepTestResponse("$.TransitRouterAttachments", []interface{}{}), nil
		case "DetachCenChildInstance":
			id := fmt.Sprint(body["ChildInstanceId"])
			if id == "vpc-protected" {
				t.Fatal("protected VPC detached")
			}
			mutations = append(mutations, id)
			detached[id] = true
			return nil, nil
		default:
			t.Fatalf("unexpected CEN inventory API: %s", action)
			return nil, nil
		}
	})
	scoped := newVpcCenInventory("cn-beijing", rpc)
	if err := sweepVpcsWithCallbacks("cn-beijing", rpc, func(id string) error { return sweepVpcCenAttachments(scoped(id), "cn-beijing", id, false) }); err != nil {
		t.Fatal(err)
	}
	if globalLists != 1 || trLists != 0 || regionalTRLists != 2 {
		t.Fatalf("region inventory rebuilt: DescribeCens=%d TR=%d regionalTR=%d", globalLists, trLists, regionalTRLists)
	}
	if childLists["cen-a"] != 3 || childLists["cen-b"] != 3 {
		t.Fatalf("expected once inventory + fresh scoped list + fresh wait per candidate, got %v", childLists)
	}
	if len(mutations) != 2 || !detached["vpc-a"] || !detached["vpc-b"] {
		t.Fatalf("unexpected mutations %v", mutations)
	}
}

func TestUnitVpcCenInventoryErrorsRemainCandidates(t *testing.T) {
	for _, malformed := range []bool{false, true} {
		t.Run(fmt.Sprintf("malformed=%v", malformed), func(t *testing.T) {
			calls := 0
			rpc := vpcSweepRPC(func(product, version, action string, query, body map[string]interface{}, auto bool) (map[string]interface{}, error) {
				switch action {
				case "DescribeCens":
					calls++
					return allVpcSweepTestResponse("$.Cens.Cen", []interface{}{map[string]interface{}{"CenId": "cen-unknown"}}), nil
				case "DescribeCenAttachedChildInstances":
					if malformed {
						return allVpcSweepTestResponse("$.ChildInstances.ChildInstance", []interface{}{map[string]interface{}{}}), nil
					}
					return nil, fmt.Errorf("inventory child denied")
				case "ListTransitRouters":
					return allVpcSweepTestResponse("$.TransitRouters", []interface{}{map[string]interface{}{"CenId": body["CenId"], "RegionId": "cn-beijing", "TransitRouterId": "tr-regional", "Type": "Basic"}}), nil
				case "ListTransitRouterVpcAttachments":
					return allVpcSweepTestResponse("$.TransitRouterAttachments", []interface{}{}), nil
				default:
					t.Fatalf("unverified resource mutated: %s", action)
					return nil, nil
				}
			})
			scoped := newVpcCenInventory("cn-beijing", rpc)
			for _, id := range []string{"vpc-a", "vpc-b"} {
				err := sweepVpcCenAttachments(scoped(id), "cn-beijing", id, false)
				if err == nil {
					t.Fatal("unavailable/unknown inventory silently excluded a CEN")
				}
				expected := "inventory child denied"
				if malformed {
					expected = "unverifiable ownership"
				}
				if !strings.Contains(err.Error(), expected) {
					t.Fatalf("inventory error hidden: %v", err)
				}
			}
			if calls != 1 {
				t.Fatalf("error inventory rebuilt per VPC: %d", calls)
			}
		})
	}
	count := 0
	failed := vpcSweepRPC(func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
		count++
		return nil, fmt.Errorf("DescribeCens denied")
	})
	scoped := newVpcCenInventory("cn-beijing", failed)
	for _, id := range []string{"vpc-a", "vpc-b"} {
		if err := sweepVpcCenAttachments(scoped(id), "cn-beijing", id, false); err == nil || !strings.Contains(err.Error(), "DescribeCens denied") {
			t.Fatalf("global inventory error hidden: %v", err)
		}
	}
	if count != 1 {
		t.Fatalf("failed global inventory not shared: %d", count)
	}
}

func TestUnitVpcCenWithoutRegionalTransitRouterSkipsAttachments(t *testing.T) {
	regionalReads := 0
	rpc := vpcSweepRPC(func(product, version, action string, query, body map[string]interface{}, auto bool) (map[string]interface{}, error) {
		switch action {
		case "DescribeCens":
			return allVpcSweepTestResponse("$.Cens.Cen", []interface{}{map[string]interface{}{"CenId": "cen-classic-only"}}), nil
		case "DescribeCenAttachedChildInstances":
			return allVpcSweepTestResponse("$.ChildInstances.ChildInstance", []interface{}{
				map[string]interface{}{"CenId": "cen-classic-only", "ChildInstanceType": "VPC", "ChildInstanceRegionId": "cn-beijing", "ChildInstanceId": "vpc-a"},
				map[string]interface{}{"CenId": "cen-classic-only", "ChildInstanceType": "VPC", "ChildInstanceRegionId": "cn-beijing", "ChildInstanceId": "vpc-b"},
				map[string]interface{}{"CenId": "cen-classic-only", "ChildInstanceType": "VPC", "ChildInstanceRegionId": "cn-beijing", "ChildInstanceId": "vpc-protected"},
			}), nil
		case "ListTransitRouters":
			regionalReads++
			if body["CenId"] != "cen-classic-only" || body["RegionId"] != "cn-beijing" {
				t.Fatalf("regional inventory scope changed: %v", body)
			}
			return allVpcSweepTestResponse("$.TransitRouters", []interface{}{}), nil
		default:
			t.Fatalf("absent regional TR must not query attachments or mutate resources: %s", action)
			return nil, nil
		}
	})
	// Legacy/direct callers without an index also prove regional absence first.
	if err := sweepVpcCenAttachments(rpc, "cn-beijing", "vpc-a", true); err != nil {
		t.Fatal(err)
	}
	if regionalReads != 1 {
		t.Fatalf("direct regional inventory reads=%d", regionalReads)
	}
	scoped := newVpcCenInventory("cn-beijing", rpc)
	for _, id := range []string{"vpc-a", "vpc-b"} {
		if err := sweepVpcCenAttachments(scoped(id), "cn-beijing", id, true); err != nil {
			t.Fatal(err)
		}
	}
	if regionalReads != 2 {
		t.Fatalf("task-local regional TR presence was not shared: reads=%d", regionalReads)
	}
}

func TestUnitVpcCenRegionalTransitRouterIdentityAndErrors(t *testing.T) {
	for _, row := range []map[string]interface{}{
		{}, {"CenId": "cen-target", "RegionId": "cn-beijing"},
		{"CenId": "cen-other", "RegionId": "cn-beijing", "TransitRouterId": "tr-1", "Type": "Enterprise"},
		{"CenId": "cen-target", "RegionId": "cn-hangzhou", "TransitRouterId": "tr-1", "Type": "Enterprise"},
	} {
		rpc := vpcSweepRPC(func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
			return allVpcSweepTestResponse("$.TransitRouters", []interface{}{row}), nil
		})
		if _, err := vpcCenRegionalTransitRouters(rpc, "cn-beijing", "cen-target"); err == nil {
			t.Fatalf("unknown identity counted as no TR: %v", row)
		}
	}
	rpc := vpcSweepRPC(func(product, version, action string, query, body map[string]interface{}, auto bool) (map[string]interface{}, error) {
		switch action {
		case "DescribeCens":
			return allVpcSweepTestResponse("$.Cens.Cen", []interface{}{map[string]interface{}{"CenId": "cen-target"}}), nil
		case "ListTransitRouters":
			return allVpcSweepTestResponse("$.TransitRouters", []interface{}{map[string]interface{}{"CenId": "cen-target", "RegionId": "cn-beijing", "TransitRouterId": "tr-1", "Type": "Enterprise"}}), nil
		case "ListTransitRouterVpcAttachments":
			return nil, fmt.Errorf("IllegalParam.CenIdOrRegionId")
		default:
			t.Fatalf("unexpected mutation %s", action)
			return nil, nil
		}
	})
	if err := sweepVpcCenAttachments(rpc, "cn-beijing", "vpc-target", true); err == nil || !strings.Contains(err.Error(), "IllegalParam.CenIdOrRegionId") {
		t.Fatalf("attachment API error incorrectly treated as no regional TR: %v", err)
	}
}

func TestUnitVpcRouterInterfaceWaitInactiveAndScope(t *testing.T) {
	calls, pauses := 0, 0
	rpc := vpcSweepRPC(func(product, version, action string, query, body map[string]interface{}, auto bool) (map[string]interface{}, error) {
		if action != "DescribeRouterInterfaces" || body["Filter.1.Key"] != "RouterInterfaceId" || body["Filter.1.Value.1"] != "ri-target" {
			t.Fatalf("unsafe RI wait API %s %v", action, body)
		}
		calls++
		status := "Deactivating"
		if calls > 1 {
			status = "Inactive"
		}
		return allVpcSweepTestResponse("$.RouterInterfaceSet.RouterInterfaceType", []interface{}{
			map[string]interface{}{"RouterInterfaceId": "ri-protected", "VpcInstanceId": "vpc-protected", "Status": "Active"},
			map[string]interface{}{"RouterType": "VRouter", "RouterId": "vrt-target", "RouterInterfaceId": "ri-target", "VpcInstanceId": "vpc-target", "Status": status},
		}), nil
	})
	if present, err := waitVpcRouterInterfaceInactive(rpc, "cn-beijing", "vpc-target", "vrt-target", "ri-target", func(time.Duration) { pauses++ }, 2); err != nil || !present {
		t.Fatalf("inactive identity not confirmed: present=%v err=%v", present, err)
	}
	if calls != 2 || pauses != 1 {
		t.Fatalf("deactivation was not awaited: calls=%d pauses=%d", calls, pauses)
	}
	for _, item := range []map[string]interface{}{
		{"RouterType": "VRouter", "RouterId": "vrt-target", "RouterInterfaceId": "ri-target", "VpcInstanceId": "vpc-protected", "Status": "Inactive"},
		{"RouterType": "VRouter", "RouterId": "vrt-target", "RouterInterfaceId": "ri-target", "VpcInstanceId": "vpc-target", "Status": "Unknown"},
		{"VpcInstanceId": "vpc-target", "Status": "Inactive"},
	} {
		invalid := vpcSweepRPC(func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
			return allVpcSweepTestResponse("$.RouterInterfaceSet.RouterInterfaceType", []interface{}{item}), nil
		})
		if _, err := waitVpcRouterInterfaceInactive(invalid, "cn-beijing", "vpc-target", "vrt-target", "ri-target", func(time.Duration) { t.Fatal("invalid RI wait paused") }, 1); err == nil {
			t.Fatalf("unknown/changed RI identity accepted: %v", item)
		}
	}
	active := vpcSweepRPC(func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
		return allVpcSweepTestResponse("$.RouterInterfaceSet.RouterInterfaceType", []interface{}{map[string]interface{}{"RouterType": "VRouter", "RouterId": "vrt-target", "RouterInterfaceId": "ri-target", "VpcInstanceId": "vpc-target", "Status": "Active"}}), nil
	})
	if _, err := waitVpcRouterInterfaceInactive(active, "cn-beijing", "vpc-target", "vrt-target", "ri-target", func(time.Duration) {}, 2); err == nil {
		t.Fatal("still Active RI allowed through to deletion")
	}
}

func TestUnitVpcRouterInterfaceStatusErrorRetry(t *testing.T) {
	for _, code := range []string{"IncorrectStatus", "InvalidInstance.StatusError", "DependencyViolation.RouterInterfaceReferedByRouteEntry", "OperationConflict"} {
		if !retryableVpcRouterInterfaceDelete(&tea.SDKError{Code: tea.String(code), StatusCode: tea.Int(400), Message: tea.String("operation pending")}) {
			t.Fatalf("native transient code not retried: %s", code)
		}
	}
	if retryableVpcRouterInterfaceDelete(&tea.SDKError{Code: tea.String("Forbidden"), StatusCode: tea.Int(403), Message: tea.String("denied")}) {
		t.Fatal("permission error retried as pending status")
	}
}

func TestUnitVpcCenRouterGenerationUsesCorrectAPI(t *testing.T) {
	for _, scenario := range []struct {
		name    string
		types   []string
		invalid bool
	}{
		{"enterprise", []string{"Enterprise"}, false},
		{"basic", []string{"Basic"}, false},
		{"unknown", []string{"Unknown"}, true},
		{"mixed", []string{"Basic", "Enterprise"}, true},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			attachmentReads := 0
			rpc := vpcSweepRPC(func(product, version, action string, query, body map[string]interface{}, auto bool) (map[string]interface{}, error) {
				switch action {
				case "DescribeCens":
					return allVpcSweepTestResponse("$.Cens.Cen", []interface{}{map[string]interface{}{"CenId": "cen-shared"}}), nil
				case "DescribeCenAttachedChildInstances":
					return allVpcSweepTestResponse("$.ChildInstances.ChildInstance", []interface{}{
						map[string]interface{}{"CenId": "cen-shared", "ChildInstanceId": "vpc-target", "ChildInstanceType": "VPC", "ChildInstanceRegionId": "cn-beijing"},
						map[string]interface{}{"CenId": "cen-shared", "ChildInstanceId": "vpc-protected", "ChildInstanceType": "VPC", "ChildInstanceRegionId": "cn-beijing"},
					}), nil
				case "ListTransitRouters":
					rows := []interface{}{}
					for i, kind := range scenario.types {
						rows = append(rows, map[string]interface{}{"CenId": "cen-shared", "RegionId": "cn-beijing", "TransitRouterId": fmt.Sprintf("tr-%d", i), "Type": kind})
					}
					return allVpcSweepTestResponse("$.TransitRouters", rows), nil
				case "ListTransitRouterVpcAttachments":
					if scenario.name != "enterprise" {
						t.Fatalf("attachment API used for %s router", scenario.name)
					}
					attachmentReads++
					// It is valid for the original Enterprise attachment to be gone already.
					return allVpcSweepTestResponse("$.TransitRouterAttachments", []interface{}{}), nil
				default:
					t.Fatalf("classic Detach or shared entity mutation used for %s: %s", scenario.name, action)
					return nil, nil
				}
			})
			err := sweepVpcCenAttachments(rpc, "cn-beijing", "vpc-target", true)
			if (err != nil) != scenario.invalid {
				t.Fatalf("generation selection error=%v invalid=%v", err, scenario.invalid)
			}
			if scenario.name == "enterprise" {
				if err := sweepVpcCenAttachments(rpc, "cn-beijing", "vpc-target", false); err != nil {
					t.Fatalf("Enterprise child should be handed to TR node, not classic Detach: %v", err)
				}
				if attachmentReads != 1 {
					t.Fatalf("Enterprise TR node did not freshly inspect attachments: %d", attachmentReads)
				}
			} else if attachmentReads != 0 {
				t.Fatal("non-Enterprise attachment queried")
			}
		})
	}
}

func TestUnitVpcAckInitialWaitBeforeDelete(t *testing.T) {
	reads, pauses, deleted := 0, 0, false
	lastState := ""
	api := vpcACKAPI{
		list: func(int) ([]map[string]interface{}, error) {
			return []map[string]interface{}{{"cluster_id": "c-target", "vpc_id": "vpc-target", "region_id": "cn-hangzhou"}}, nil
		},
		detail: func(id string) (map[string]interface{}, error) {
			if id != "c-target" {
				t.Fatal("wrong ACK queried")
			}
			reads++
			if deleted {
				return nil, &tea.SDKError{StatusCode: tea.Int(404), Message: tea.String("cluster not found")}
			}
			lastState = "initial"
			if reads > 1 {
				lastState = "running"
			}
			return map[string]interface{}{"cluster_id": id, "vpc_id": "vpc-target", "region_id": "cn-hangzhou", "state": lastState, "deletion_protection": false}, nil
		},
		disableProtection: func(string) error { t.Fatal("initial cluster unexpectedly mutated"); return nil },
		remove: func(id string) error {
			if id != "c-target" || lastState != "running" || pauses != 1 {
				t.Fatal("DeleteCluster invoked before initial state completed")
			}
			deleted = true
			return nil
		},
		pause: func(d time.Duration) {
			if d != 10*time.Second {
				t.Fatalf("wrong initial polling interval %s", d)
			}
			pauses++
		},
		initialAttempts: 2,
	}
	if err := vpcACKSweeper("cn-hangzhou", "vpc-target", api).F("cn-hangzhou"); err != nil {
		t.Fatal(err)
	}
	if !deleted || reads != 3 || pauses != 1 {
		t.Fatalf("wrong wait/delete lifecycle: deleted=%v reads=%d pauses=%d", deleted, reads, pauses)
	}
}

func TestUnitVpcAckInitialTimeoutLeavesResource(t *testing.T) {
	reads, pauses := 0, 0
	api := vpcACKAPI{
		list: func(int) ([]map[string]interface{}, error) {
			return []map[string]interface{}{{"cluster_id": "c-target", "vpc_id": "vpc-target", "region_id": "cn-hangzhou"}}, nil
		},
		detail: func(string) (map[string]interface{}, error) {
			reads++
			return map[string]interface{}{"cluster_id": "c-target", "vpc_id": "vpc-target", "region_id": "cn-hangzhou", "state": "initial", "deletion_protection": false}, nil
		},
		disableProtection: func(string) error { t.Fatal("initial timeout changed deletion protection"); return nil },
		remove:            func(string) error { t.Fatal("initial timeout invoked DeleteCluster"); return nil },
		pause:             func(time.Duration) { pauses++ }, initialAttempts: 2,
	}
	err := vpcACKSweeper("cn-hangzhou", "vpc-target", api).F("cn-hangzhou")
	if err == nil || !strings.Contains(err.Error(), "remains initial") || !strings.Contains(err.Error(), "dependencies remain") {
		t.Fatalf("initial residue not reported: %v", err)
	}
	if reads != 2 || pauses != 1 {
		t.Fatalf("wait was not bounded: reads=%d pauses=%d", reads, pauses)
	}
	// The wall-clock deadline also bounds a large poll budget without sleeping.
	reads, pauses = 0, 0
	fakeNow := time.Unix(0, 0)
	api.now = func() time.Time { return fakeNow }
	api.pause = func(time.Duration) { pauses++; fakeNow = fakeNow.Add(20 * time.Minute) }
	api.initialAttempts = 120
	if err := vpcACKSweeper("cn-hangzhou", "vpc-target", api).F("cn-hangzhou"); err == nil {
		t.Fatal("creation wait exceeded its wall-clock deadline")
	}
	if reads != 1 || pauses != 1 {
		t.Fatalf("deadline did not stop polling: reads=%d pauses=%d", reads, pauses)
	}

}

func TestUnitVpcAckInitialWaitScopeAndGone(t *testing.T) {
	for _, scenario := range []string{"changed-vpc", "changed-region", "changed-id", "unknown-state", "missing-state", "read-error", "not-found", "deleted"} {
		t.Run(scenario, func(t *testing.T) {
			reads := 0
			api := vpcACKAPI{
				list: func(int) ([]map[string]interface{}, error) {
					return []map[string]interface{}{{"cluster_id": "c-target", "vpc_id": "vpc-target", "region_id": "cn-hangzhou"}}, nil
				},
				detail: func(string) (map[string]interface{}, error) {
					reads++
					item := map[string]interface{}{"cluster_id": "c-target", "vpc_id": "vpc-target", "region_id": "cn-hangzhou", "state": "initial", "deletion_protection": false}
					if reads == 1 {
						return item, nil
					}
					switch scenario {
					case "changed-vpc":
						item["vpc_id"] = "vpc-protected"
					case "changed-region":
						item["region_id"] = "cn-beijing"
					case "changed-id":
						item["cluster_id"] = "c-protected"
					case "unknown-state":
						item["state"] = "future-unknown"
					case "missing-state":
						delete(item, "state")
					case "read-error":
						return nil, fmt.Errorf("ACK read denied")
					case "not-found":
						return nil, &tea.SDKError{StatusCode: tea.Int(404), Message: tea.String("cluster not found")}
					case "deleted":
						item["state"] = "deleted"
					}
					return item, nil
				},
				disableProtection: func(string) error { t.Fatal("unverified or gone cluster protection changed"); return nil },
				remove:            func(string) error { t.Fatal("unverified or gone cluster deleted"); return nil },
				pause:             func(time.Duration) {}, initialAttempts: 2,
			}
			err := vpcACKSweeper("cn-hangzhou", "vpc-target", api).F("cn-hangzhou")
			gone := scenario == "not-found" || scenario == "deleted"
			if (err == nil) != gone {
				t.Fatalf("gone=%v got error %v", gone, err)
			}
			if reads != 2 {
				t.Fatalf("invalid/gone state did not stop promptly: reads=%d", reads)
			}
		})
	}
}

func TestUnitVpcRouterInterfaceNativeRouterOwnership(t *testing.T) {
	for _, tc := range []struct {
		name           string
		item           map[string]interface{}
		owned, invalid bool
	}{
		{"target", map[string]interface{}{"RouterType": "VRouter", "RouterId": "vrt-target", "VpcInstanceId": "vpc-target"}, true, false},
		{"vbr-no-vpc", map[string]interface{}{"RouterType": "VBR", "RouterId": "vbr-other"}, false, false},
		{"other-vpc", map[string]interface{}{"RouterType": "VRouter", "RouterId": "vrt-other", "VpcInstanceId": "vpc-protected"}, false, false},
		{"target-missing-vpc", map[string]interface{}{"RouterType": "VRouter", "RouterId": "vrt-target"}, false, true},
		{"target-changed-vpc", map[string]interface{}{"RouterType": "VRouter", "RouterId": "vrt-target", "VpcInstanceId": "vpc-protected"}, false, true},
		{"contradictory-router", map[string]interface{}{"RouterType": "VRouter", "RouterId": "vrt-other", "VpcInstanceId": "vpc-target"}, false, true},
		{"missing-router", map[string]interface{}{"RouterType": "VRouter", "VpcInstanceId": "vpc-target"}, false, true},
		{"unknown-type", map[string]interface{}{"RouterType": "Unknown", "RouterId": "vrt-target", "VpcInstanceId": "vpc-target"}, false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			owned, err := vpcRouterInterfaceOwned(tc.item, "vpc-target", "vrt-target")
			if owned != tc.owned || (err != nil) != tc.invalid {
				t.Fatalf("owned=%v err=%v", owned, err)
			}
		})
	}
}

func TestUnitVpcRouterInterfaceVpcDetailFailsClosed(t *testing.T) {
	for _, detail := range []map[string]interface{}{
		{}, {"VpcId": "vpc-protected", "RegionId": "cn-beijing", "VRouterId": "vrt-target"},
		{"VpcId": "vpc-target", "RegionId": "cn-hangzhou", "VRouterId": "vrt-target"},
		{"VpcId": "vpc-target", "RegionId": "cn-beijing"},
	} {
		calls := 0
		rpc := vpcSweepRPC(func(product, version, action string, query, body map[string]interface{}, auto bool) (map[string]interface{}, error) {
			calls++
			if action != "DescribeVpcAttribute" || body["VpcId"] != "vpc-target" {
				t.Fatalf("unverified VPC touched %s", action)
			}
			return detail, nil
		})
		if err := sweepVpcRouterInterfaces(rpc, "cn-beijing", "vpc-target"); err == nil || calls != 1 {
			t.Fatalf("invalid VPC accepted calls=%d err=%v", calls, err)
		}
	}
}

func TestUnitVpcRouterInterfaceWaitDeactivatingFailure(t *testing.T) {
	for _, apiError := range []bool{false, true} {
		calls := 0
		rpc := vpcSweepRPC(func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
			calls++
			if apiError {
				return nil, fmt.Errorf("read denied")
			}
			return allVpcSweepTestResponse("$.RouterInterfaceSet.RouterInterfaceType", []interface{}{map[string]interface{}{"RouterInterfaceId": "ri-target", "RouterType": "VRouter", "RouterId": "vrt-target", "VpcInstanceId": "vpc-target", "Status": "Deactivating"}}), nil
		})
		if present, err := waitVpcRouterInterfaceInactive(rpc, "cn-beijing", "vpc-target", "vrt-target", "ri-target", func(time.Duration) {}, 2); err == nil || present {
			t.Fatalf("unconfirmed transition accepted present=%v err=%v", present, err)
		}
		want := 2
		if apiError {
			want = 1
		}
		if calls != want {
			t.Fatalf("calls=%d want=%d", calls, want)
		}
	}
}

// One limiter is shared by all selected VPC workers. Individual wrappers must
// not reset the budget, otherwise parallel VPCs can still overwhelm an API.
var allVpcSweepRPCLimiter = &vpcSweepRPCLimiter{next: map[string]time.Time{}, now: time.Now, sleep: time.Sleep}

type vpcSweepRPCLimiter struct {
	mu    sync.Mutex
	next  map[string]time.Time
	now   func() time.Time
	sleep func(time.Duration)
}

func vpcSweepNetworkACLAction(product, action string) bool {
	return strings.EqualFold(product, "Vpc") && (action == "DescribeNetworkAcls" || action == "DeleteNetworkAcl")
}

func (l *vpcSweepRPCLimiter) wait(product, action string) {
	interval := 200 * time.Millisecond
	key := strings.ToLower(product) + "/" + action
	if vpcSweepNetworkACLAction(product, action) {
		// Read and delete share the account ACL quota. Three seconds also
		// leaves headroom when two independent regional FC runners overlap.
		interval = 3 * time.Second
		key = "vpc/network-acl"
	}
	l.mu.Lock()
	now := l.now()
	slot := l.next[key]
	if slot.Before(now) {
		slot = now
	}
	l.next[key] = slot.Add(interval)
	l.mu.Unlock()
	if delay := slot.Sub(now); delay > 0 {
		l.sleep(delay)
	}
}

func throttledVpcSweepRPC(rpc vpcSweepRPC) vpcSweepRPC {
	return retryingVpcSweepRPC(rpc, allVpcSweepRPCLimiter)
}

func retryingVpcSweepRPC(rpc vpcSweepRPC, limiter *vpcSweepRPCLimiter) vpcSweepRPC {
	return func(product, version, action string, query, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
		var response map[string]interface{}
		var err error
		attempts := 4
		if vpcSweepNetworkACLAction(product, action) {
			attempts = 6
		}
		for attempt := 0; attempt < attempts; attempt++ {
			limiter.wait(product, action)
			// The underlying SDK already handles 5xx when autoRetry is enabled.
			// Only structured throttling errors get this additional bounded retry.
			response, err = rpc(product, version, action, query, body, autoRetry)
			if err == nil || !rdsErrorHasCode(err, "Throttling", "Throttling.User", "Throttling.Api", "Throttling.System", "Throttling.RateQuota", "Throttling.User.Api") {
				return response, err
			}
			if attempt+1 < attempts {
				limiter.sleep(time.Second * time.Duration(1<<uint(attempt)))
			}
		}
		return response, err
	}
}

func TestUnitVpcSweepRateBudgetSharedAcrossWrappers(t *testing.T) {
	now := time.Unix(100, 0)
	var waits []time.Duration
	limiter := &vpcSweepRPCLimiter{next: map[string]time.Time{}, now: func() time.Time { return now }, sleep: func(d time.Duration) { waits = append(waits, d); now = now.Add(d) }}
	rpc := vpcSweepRPC(func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
		return nil, nil
	})
	first, second := retryingVpcSweepRPC(rpc, limiter), retryingVpcSweepRPC(rpc, limiter)
	first("Vpc", "v", "DescribeNetworkAcls", nil, nil, false)
	second("vpc", "v", "DeleteNetworkAcl", nil, nil, false)
	first("Vpc", "v", "DescribeNetworkAcls", nil, nil, false)
	if len(waits) != 2 || waits[0] != 3*time.Second || waits[1] != 3*time.Second {
		t.Fatalf("shared ACL budget: %v", waits)
	}
	// Different products/actions have independent budgets.
	second("Ecs", "v", "DescribeSecurityGroups", nil, nil, false)
	if len(waits) != 2 {
		t.Fatalf("unrelated API delayed: %v", waits)
	}
}

func TestUnitVpcSweepThrottleRetriesBounded(t *testing.T) {
	for _, permanent := range []bool{false, true} {
		calls := 0
		now := time.Unix(100, 0)
		limiter := &vpcSweepRPCLimiter{next: map[string]time.Time{}, now: func() time.Time { return now }, sleep: func(d time.Duration) { now = now.Add(d) }}
		rpc := retryingVpcSweepRPC(func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
			calls++
			if permanent {
				return nil, fmt.Errorf("permission denied")
			}
			return nil, errors.NewServerError(400, `{"Code":"Throttling.User","Message":"rate limit"}`, "")
		}, limiter)
		_, err := rpc("Vpc", "v", "DescribeNetworkAcls", nil, nil, false)
		want := 6
		if permanent {
			want = 1
		}
		if err == nil || calls != want {
			t.Fatalf("permanent=%v calls=%d error=%v", permanent, calls, err)
		}
	}
}

func TestUnitVpcSweepThrottleRecoversSDKError(t *testing.T) {
	now := time.Unix(100, 0)
	calls := 0
	limiter := &vpcSweepRPCLimiter{next: map[string]time.Time{}, now: func() time.Time { return now }, sleep: func(d time.Duration) { now = now.Add(d) }}
	rpc := retryingVpcSweepRPC(func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
		calls++
		if calls == 1 {
			return nil, tea.NewSDKError(map[string]interface{}{"code": "Throttling.User", "message": "rate limit"})
		}
		return map[string]interface{}{"deleted": true}, nil
	}, limiter)
	response, err := rpc("Vpc", "v", "DeleteNetworkAcl", nil, nil, false)
	if err != nil || calls != 2 || response["deleted"] != true {
		t.Fatalf("throttle recovery calls=%d response=%v error=%v", calls, response, err)
	}
}

func TestUnitVpcSweepNetworkACLQuotaIsolation(t *testing.T) {
	now := time.Unix(100, 0)
	var waits []time.Duration
	limiter := &vpcSweepRPCLimiter{next: map[string]time.Time{}, now: func() time.Time { return now }, sleep: func(d time.Duration) { waits = append(waits, d); now = now.Add(d) }}
	limiter.wait("Vpc", "DescribeNetworkAcls")
	limiter.wait("Vpc", "DescribeVSwitches")
	limiter.wait("Other", "DescribeNetworkAcls")
	limiter.wait("Other", "DeleteNetworkAcl")
	if len(waits) != 0 {
		t.Fatalf("unrelated action consumed ACL quota: %v", waits)
	}
	limiter.wait("Vpc", "DeleteNetworkAcl")
	if len(waits) != 1 || waits[0] != 3*time.Second {
		t.Fatalf("read/write quota not shared: %v", waits)
	}
	limiter.wait("Ecs", "DescribeSecurityGroups")
	limiter.wait("Ecs", "DescribeSecurityGroups")
	if len(waits) != 2 || waits[1] != 200*time.Millisecond {
		t.Fatalf("non-ACL interval changed: %v", waits)
	}
}

func TestUnitVpcSweepNetworkACLRetryPolicy(t *testing.T) {
	sdkThrottle := tea.NewSDKError(map[string]interface{}{"code": "Throttling.User", "message": "rate limit"})
	for _, tc := range []struct {
		name, product, action string
		err                   error
		attempts              int
	}{
		{"raw-tea", "Vpc", "DescribeNetworkAcls", sdkThrottle, 6},
		{"wrapped-tea", "Vpc", "DeleteNetworkAcl", &ComplexError{Cause: &ComplexError{Cause: sdkThrottle}}, 6},
		{"other-action", "Vpc", "DescribeVSwitches", sdkThrottle, 4},
		{"other-product", "Other", "DescribeNetworkAcls", sdkThrottle, 4},
		{"permission", "Vpc", "DeleteNetworkAcl", tea.NewSDKError(map[string]interface{}{"code": "Forbidden", "message": "Throttling.User in diagnostic text"}), 1},
		{"unstructured", "Vpc", "DescribeNetworkAcls", fmt.Errorf("Throttling.User"), 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			now := time.Unix(100, 0)
			var pauses []time.Duration
			limiter := &vpcSweepRPCLimiter{next: map[string]time.Time{}, now: func() time.Time { return now }, sleep: func(d time.Duration) { pauses = append(pauses, d); now = now.Add(d) }}
			calls := 0
			rpc := retryingVpcSweepRPC(func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
				calls++
				return nil, tc.err
			}, limiter)
			_, err := rpc(tc.product, "v", tc.action, nil, nil, false)
			if err != tc.err || calls != tc.attempts {
				t.Fatalf("calls=%d want=%d original error preserved=%v", calls, tc.attempts, err == tc.err)
			}
			if tc.attempts == 1 && len(pauses) != 0 {
				t.Fatalf("non-throttle retried: %v", pauses)
			}
			if tc.attempts == 6 {
				// Rate waits can add pauses, but each exponential retry backoff must occur
				// once, in order, and exhaustion must not add a seventh RPC attempt.
				next := 0
				for _, d := range pauses {
					if next < 5 && d == time.Second*time.Duration(1<<uint(next)) {
						next++
					}
				}
				if next != 5 {
					t.Fatalf("missing bounded ACL backoff: %v", pauses)
				}
			}
		})
	}
}

// Each factory owns its cache: never share account-specific capability across
// clients, regions, credentials, or independent sweeps.
func nativeVpcServiceAvailability(client *connectivity.AliyunClient) func(string) (bool, error) {
	region := client.RegionId
	// Snapshot before probes: RpcGet itself caches a discovered endpoint. That
	// internal discovery must not be mistaken for a new custom override later.
	custom := map[string]bool{}
	for _, product := range []string{"gpdb", "foasconsole", "dataworks_public"} {
		custom[product] = customVpcServiceEndpoint(product)
	}
	flink := nativeVpcFlinkRPC(client)
	return cachedVpcServiceAvailability(region, func(product string) bool { return custom[product] }, func(product, region string) (bool, error) {
		return probeVpcServiceAvailability(product, region, nil, flink)
	})
}

func cachedVpcServiceAvailability(region string, custom func(string) bool, probe func(string, string) (bool, error)) func(string) (bool, error) {
	type answer struct {
		available bool
		err       error
	}
	cache := map[string]answer{}
	var mu sync.Mutex
	return func(product string) (bool, error) {
		if custom(product) {
			return true, nil
		}
		mu.Lock()
		defer mu.Unlock()
		key := product + "/" + region
		if found, ok := cache[key]; ok {
			return found.available, found.err
		}
		available, err := probe(product, region)
		cache[key] = answer{available, err}
		return available, err
	}
}

func customVpcServiceEndpoint(product string) bool {
	return configuredVpcServiceEndpoint(product, os.Getenv, func(name string) bool { _, err := os.Stat(name); return !os.IsNotExist(err) }, func(key string) bool {
		value, ok := endpoints.Load(key)
		return ok && value != nil && fmt.Sprint(value) != ""
	})
}
func configuredVpcServiceEndpoint(product string, getenv func(string) string, exists func(string) bool, configured func(string) bool) bool {
	if getenv("TF_ENDPOINT_PATH") != "" || exists("./endpoints.xml") {
		return true
	}
	names := []string{product, "location"}
	if product == "dataworks_public" {
		names = append(names, "dataworks", "dide")
	}
	for _, name := range names {
		if getenv("ALIBABA_CLOUD_ENDPOINT_"+strings.ToUpper(name)) != "" || getenv(strings.ToUpper(name)+"_ENDPOINT") != "" || configured(name) {
			return true
		}
	}
	return false
}

func probeVpcServiceAvailability(product, region string, location, flink vpcSweepRPC) (bool, error) {
	if product == "foasconsole" {
		response, err := flink("foasconsole", "2019-06-01", "DescribeSupportedRegions", nil, nil, true)
		if err != nil {
			return false, err
		}
		if success, ok := response["Success"].(bool); !ok || !success {
			return false, fmt.Errorf("Flink region capability response not successful")
		}
		regions, ok := response["Regions"].([]interface{})
		if !ok {
			return false, fmt.Errorf("Flink region capability missing Regions")
		}
		found := false
		for _, raw := range regions {
			item, ok := raw.(map[string]interface{})
			if !ok {
				return false, fmt.Errorf("Flink malformed region")
			}
			id, ok := item["Region"].(string)
			if !ok || id == "" {
				return false, fmt.Errorf("Flink missing Region")
			}
			if id == region {
				found = true
			}
		}
		return found, nil
	}
	// Location is a routing registry, not a product availability catalog. Beijing
	// DataWorks returns no registered endpoint while its official fallback serves
	// real CommonV2 groups. Empty registries and InvalidRegionId must never skip
	// these native owner callbacks or erase their real business/network errors.
	if product == "gpdb" || product == "dataworks_public" {
		return true, nil
	}
	return false, fmt.Errorf("unsupported capability product")
}

func applyVpcServiceAvailability(graph map[string]*resource.Sweeper, capability func(string) (bool, error)) {
	for name, product := range map[string]string{"alicloud_gpdb_instance": "gpdb", "alicloud_realtime_compute_vvp_instance": "foasconsole", "alicloud_data_works_network": "dataworks_public"} {
		node := graph[name]
		if node == nil {
			continue
		}
		original := node.F
		product := product
		name := name
		node.F = func(region string) error {
			available, err := capability(product)
			if err != nil {
				return fmt.Errorf("%s region capability: %w", name, err)
			}
			if !available {
				log.Printf("[INFO] SKIP %s: service unavailable in selected region", name)
				return nil
			}
			if original == nil {
				return fmt.Errorf("%s has no sweeper callback", name)
			}
			return original(region)
		}
	}
}

func TestUnitVpcServiceAvailabilityRegistryCannotSkipNativeFallback(t *testing.T) {
	for _, product := range []string{"dataworks_public", "gpdb"} {
		for _, reply := range []struct {
			response map[string]interface{}
			err      error
		}{
			{map[string]interface{}{"Success": true, "Endpoints": map[string]interface{}{"Endpoint": []interface{}{}}}, nil},
			{nil, tea.NewSDKError(map[string]interface{}{"code": "InvalidRegionId", "message": "registry missing"})},
			{nil, fmt.Errorf("registry DNS failure")},
		} {
			calls := 0
			location := vpcSweepRPC(func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
				calls++
				return reply.response, reply.err
			})
			available, err := probeVpcServiceAvailability(product, "cn-beijing", location, nil)
			if !available || err != nil || calls != 0 {
				t.Fatalf("registry suppressed native product=%s available=%v calls=%d err=%v", product, available, calls, err)
			}
		}
	}
}

func TestUnitVpcServiceAvailabilityPreservesOpaqueNativeOwnerErrors(t *testing.T) {
	for _, product := range []struct{ node, product string }{{"alicloud_data_works_network", "dataworks_public"}, {"alicloud_gpdb_instance", "gpdb"}} {
		for _, originalErr := range []error{fmt.Errorf("DNS no such host"), fmt.Errorf("permission denied"), fmt.Errorf("unknown business error"), tea.NewSDKError(map[string]interface{}{"code": "InvalidRegionId", "message": "native error"})} {
			calls := 0
			graph := map[string]*resource.Sweeper{product.node: {F: func(string) error { calls++; return originalErr }}}
			applyVpcServiceAvailability(graph, func(p string) (bool, error) { return probeVpcServiceAvailability(p, "cn-beijing", nil, nil) })
			if err := graph[product.node].F("cn-beijing"); err != originalErr || calls != 1 {
				t.Fatalf("native error hidden for %s: calls=%d err=%v", product.product, calls, err)
			}
		}
	}
}

func TestUnitVpcServiceAvailabilityValidatesSuccessfulRegionLists(t *testing.T) {
	for _, region := range []string{"cn-guangzhou", "cn-beijing"} {
		rpc := vpcSweepRPC(func(product, version, action string, q, b map[string]interface{}, _ bool) (map[string]interface{}, error) {
			if action != "DescribeSupportedRegions" || q != nil || b != nil {
				t.Fatal("wrong Flink probe")
			}
			return map[string]interface{}{"Success": true, "Regions": []interface{}{map[string]interface{}{"Region": "cn-beijing"}}}, nil
		})
		available, err := probeVpcServiceAvailability("foasconsole", region, nil, rpc)
		if err != nil || available != (region == "cn-beijing") {
			t.Fatal("Flink region proof", available, err)
		}
	}
	rpc := vpcSweepRPC(func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
		return nil, tea.NewSDKError(map[string]interface{}{"code": "000000", "message": "No available regions for you"})
	})
	if _, err := probeVpcServiceAvailability("foasconsole", "cn-guangzhou", nil, rpc); err == nil {
		t.Fatal("generic Flink error swallowed")
	}
}

func TestUnitVpcServiceAvailabilityCacheIsPerFactoryAndPreservesErrors(t *testing.T) {
	calls := 0
	probe := func(string, string) (bool, error) { calls++; return false, fmt.Errorf("denied") }
	one := cachedVpcServiceAvailability("cn-beijing", func(string) bool { return false }, probe)
	for i := 0; i < 3; i++ {
		if _, err := one("gpdb"); err == nil {
			t.Fatal("lost error")
		}
	}
	if calls != 1 {
		t.Fatal("probe not cached", calls)
	}
	two := cachedVpcServiceAvailability("cn-beijing", func(string) bool { return false }, probe)
	two("gpdb")
	if calls != 2 {
		t.Fatal("capability leaked across task factories")
	}
	bypass := cachedVpcServiceAvailability("cn-beijing", func(string) bool { return true }, probe)
	if available, err := bypass("gpdb"); !available || err != nil || calls != 2 {
		t.Fatal("custom endpoint was probed")
	}
}

func TestUnitVpcServiceAvailabilityCustomEndpointSourcesDisableSkip(t *testing.T) {
	for _, envKey := range []string{"TF_ENDPOINT_PATH", "ALIBABA_CLOUD_ENDPOINT_DATAWORKS_PUBLIC", "DATAWORKS_PUBLIC_ENDPOINT", "ALIBABA_CLOUD_ENDPOINT_DATAWORKS", "DIDE_ENDPOINT", "LOCATION_ENDPOINT"} {
		if !configuredVpcServiceEndpoint("dataworks_public", func(key string) string {
			if key == envKey {
				return "configured"
			}
			return ""
		}, func(string) bool { return false }, func(string) bool { return false }) {
			t.Fatal("ignored override", envKey)
		}
	}
	if !configuredVpcServiceEndpoint("gpdb", func(string) string { return "" }, func(string) bool { return true }, func(string) bool { return false }) {
		t.Fatal("ignored endpoints.xml")
	}
	if !configuredVpcServiceEndpoint("foasconsole", func(string) string { return "" }, func(string) bool { return false }, func(key string) bool { return key == "foasconsole" }) {
		t.Fatal("ignored configured endpoint map")
	}
}

func TestUnitVpcServiceAvailabilitySkipsOnlyUnavailableOwnerNotVpcRoot(t *testing.T) {
	writes := 0
	root := 0
	other := 0
	graph := map[string]*resource.Sweeper{
		"alicloud_gpdb_instance": {F: func(string) error { writes++; return nil }},
		"alicloud_vpc":           {F: func(string) error { root++; return nil }},
		"alicloud_other":         {F: func(string) error { other++; return nil }},
	}
	applyVpcServiceAvailability(graph, func(string) (bool, error) { return false, nil })
	for _, key := range []string{"alicloud_gpdb_instance", "alicloud_vpc", "alicloud_other"} {
		if err := graph[key].F("cn-guangzhou"); err != nil {
			t.Fatal(err)
		}
	}
	if writes != 0 || root != 1 || other != 1 {
		t.Fatal("availability bypassed wrong graph nodes")
	}
	graph = map[string]*resource.Sweeper{"alicloud_gpdb_instance": {F: func(string) error { writes++; return nil }}}
	applyVpcServiceAvailability(graph, func(string) (bool, error) { return false, fmt.Errorf("permission denied") })
	if err := graph["alicloud_gpdb_instance"].F("cn-guangzhou"); err == nil || writes != 0 {
		t.Fatal("capability failure invoked owner")
	}
}

// The official dataworks-public-20240518 SDK uses GET for these list/detail
// actions and POST for deletion. Callers keep the sweeper RPC contract; this
// adapter relocates GET parameters from body to query without changing callers.
func nativeVpcDataWorksRPC(client *connectivity.AliyunClient) vpcSweepRPC {
	return dataWorksVpcCleanupRPC(client.RpcGet, client.RpcPost)
}

func eciVpcCleanupRPC(get func(string, string, string, map[string]interface{}, map[string]interface{}) (map[string]interface{}, error), post vpcSweepRPC) vpcSweepRPC {
	return func(product, version, action string, query, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
		if !strings.EqualFold(product, "eci") || version != "2018-08-08" {
			return nil, fmt.Errorf("unsupported ECI cleanup product/version")
		}
		parameters := map[string]interface{}{}
		for key, value := range query {
			parameters[key] = value
		}
		for key, value := range body {
			if _, exists := parameters[key]; exists {
				return nil, fmt.Errorf("duplicate ECI query parameter %s", key)
			}
			parameters[key] = value
		}
		switch action {
		case "DescribeVirtualNodes":
			// The service can silently return an empty inventory when these
			// query parameters are sent in a POST body. Use the verified GET.
			return get(product, version, action, parameters, nil)
		case "DeleteVirtualNode":
			return post(product, version, action, parameters, nil, autoRetry)
		default:
			return nil, fmt.Errorf("unsupported ECI cleanup action %s", action)
		}
	}
}

func dataWorksVpcCleanupRPC(get func(string, string, string, map[string]interface{}, map[string]interface{}) (map[string]interface{}, error), post vpcSweepRPC) vpcSweepRPC {
	return func(product, version, action string, query, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
		if product != "dataworks-public" || version != "2024-05-18" {
			return nil, fmt.Errorf("unsupported DataWorks cleanup product/version")
		}
		switch action {
		case "ListResourceGroups", "ListNetworks", "GetResourceGroup", "GetNetwork":
			parameters := map[string]interface{}{}
			for key, value := range query {
				parameters[key] = value
			}
			for key, value := range body {
				if _, exists := parameters[key]; exists {
					return nil, fmt.Errorf("duplicate DataWorks GET parameter %s", key)
				}
				parameters[key] = value
			}
			return get(product, version, action, parameters, nil)
		case "DeleteNetwork", "DeleteResourceGroup":
			return post(product, version, action, query, body, autoRetry)
		default:
			return nil, fmt.Errorf("unsupported DataWorks cleanup action %s", action)
		}
	}
}

// serviceVpcOwnerSweepers uses the owning service's native APIs, never deletes
// service-managed ENIs directly, and accepts only resources with verified VPC
// identity. List results are collected before deletion so paging cannot skip
// resources when deletion changes the list.
func serviceVpcOwnerSweepers(region, vpcID string, rpc vpcSweepRPC) map[string]*resource.Sweeper {
	nodes := map[string]*resource.Sweeper{}
	add := func(name string, run func() error) {
		nodes[name] = &resource.Sweeper{Name: name, F: func(actualRegion string) error {
			if actualRegion != region || vpcID == "" {
				return fmt.Errorf("%s requires the selected VPC and region", name)
			}
			return run()
		}}
	}
	add("alicloud_privatelink_vpc_endpoint", func() error {
		items, err := serviceVpcOwnerList(rpc, "Privatelink", "2020-04-15", "ListVpcEndpoints", "$.Endpoints", "", true)
		if err != nil {
			return err
		}
		return serviceVpcOwnerEach(items, func(item map[string]interface{}) error {
			match, err := serviceVpcOwnerMatch(item, "VpcId", vpcID)
			if err != nil || !match {
				return err
			}
			id, err := serviceVpcOwnerID(item, "EndpointId")
			if err != nil {
				return err
			}
			detail, err := rpc("Privatelink", "2020-04-15", "GetVpcEndpointAttribute", nil, map[string]interface{}{"EndpointId": id}, true)
			if err != nil {
				return err
			}
			if detail["EndpointId"] != id || detail["VpcId"] != vpcID {
				return fmt.Errorf("PrivateLink %s detail ownership mismatch; not deleted", id)
			}
			protected, known := detail["ProtectedEnabled"].(bool)
			if !known {
				return fmt.Errorf("PrivateLink %s managed protection unknown; not deleted", id)
			}
			// ProtectedEnabled is STS creator-managed protection, not a generic
			// deletion lock. The public UpdateVpcEndpointAttribute request has
			// no ProtectedEnabled field: do not invent an unlock operation.
			if protected {
				return fmt.Errorf("PrivateLink %s has STS creator-managed protection; no verified public disable API, requires original creator authorization; not deleted", id)
			}
			_, err = rpc("Privatelink", "2020-04-15", "DeleteVpcEndpoint", nil, map[string]interface{}{"EndpointId": id}, false)
			if err != nil {
				return err
			}
			return serviceVpcOwnerWaitGone(func() ([]map[string]interface{}, error) {
				return serviceVpcOwnerList(rpc, "Privatelink", "2020-04-15", "ListVpcEndpoints", "$.Endpoints", "", true)
			}, "EndpointId", id, vpcID, time.Sleep)
		})
	})
	add("alicloud_lindorm_instance", func() error {
		items, err := serviceVpcOwnerList(rpc, "hitsdb", "2020-06-15", "GetLindormInstanceList", "$.InstanceList", "", false)
		if err != nil {
			return err
		}
		return serviceVpcOwnerEach(items, func(item map[string]interface{}) error {
			// The list datasource exposes VpcId. Check it before describing so
			// unrelated instances cannot trigger cancellation or release.
			match, err := serviceVpcOwnerMatch(item, "VpcId", vpcID)
			if err != nil || !match {
				return err
			}
			id, err := serviceVpcOwnerID(item, "InstanceId")
			if err != nil {
				return err
			}
			detail, err := rpc("hitsdb", "2020-06-15", "GetLindormInstance", nil, map[string]interface{}{"InstanceId": id}, true)
			if err != nil {
				return err
			}
			if detail["InstanceId"] != id {
				return fmt.Errorf("Lindorm %s detail identity mismatch; not released", id)
			}
			match, err = serviceVpcOwnerMatch(detail, "VpcId", vpcID)
			if err != nil || !match {
				return fmt.Errorf("Lindorm %s changed or unverifiable VPC: %v", id, err)
			}
			if detail["PayType"] != "POSTPAY" {
				return fmt.Errorf("Lindorm %s payment type %v requires a verified unsubscribe capability; not released", id, detail["PayType"])
			}
			release := map[string]interface{}{"InstanceId": id}
			// Full VPC cleanup must release the instance immediately rather
			// than retain its data for seven days. This does not bypass state checks.
			if sweepAll() {
				release["Immediately"] = true
			}
			_, err = rpc("hitsdb", "2020-06-15", "ReleaseLindormInstance", nil, release, false)
			if err != nil {
				return err
			}
			return serviceVpcOwnerWaitGone(func() ([]map[string]interface{}, error) {
				return serviceVpcOwnerList(rpc, "hitsdb", "2020-06-15", "GetLindormInstanceList", "$.InstanceList", "", false)
			}, "InstanceId", id, vpcID, time.Sleep)
		})
	})
	add("alicloud_kms_instance", func() error {
		items, err := serviceVpcOwnerList(rpc, "Kms", "2016-01-20", "ListKmsInstances", "$.KmsInstances.KmsInstance", "", false)
		if err != nil {
			return err
		}
		return serviceVpcOwnerEach(items, func(item map[string]interface{}) error {
			id, err := serviceVpcOwnerID(item, "KmsInstanceId")
			if err != nil {
				return err
			}
			response, err := rpc("Kms", "2016-01-20", "GetKmsInstance", nil, map[string]interface{}{"KmsInstanceId": id}, true)
			if err != nil {
				return err
			}
			detail, ok := response["KmsInstance"].(map[string]interface{})
			if !ok {
				return fmt.Errorf("KMS %s has invalid detail", id)
			}
			if detail["InstanceId"] != id {
				return fmt.Errorf("KMS %s detail identity mismatch; not released", id)
			}
			match, err := serviceVpcOwnerMatch(detail, "VpcId", vpcID)
			if err != nil || !match {
				return err
			}
			// A shared KMS instance can serve protected/other VPCs. Do not
			// release it merely because its primary VPC was selected.
			bindings, err := kmsVpcOwnerBindings(detail)
			if err != nil {
				return fmt.Errorf("KMS %s: %w", id, err)
			}
			for _, bind := range bindings {
				match, err := serviceVpcOwnerMatch(bind, "VpcId", vpcID)
				if err != nil || !match {
					return fmt.Errorf("KMS %s has shared or unverifiable VPC bindings; not released", id)
				}
			}

			if detail["DeletionProtection"] == true {
				return fmt.Errorf("KMS %s deletion protection remains enabled; not released", id)
			}
			switch detail["ChargeType"] {
			case "POSTPAY":
				_, err = rpc("Kms", "2016-01-20", "ReleaseKmsInstance", nil, map[string]interface{}{"KmsInstanceId": id}, true)
			case "PREPAY":
				err = refundVpcKmsSubscription(rpc, region, vpcID, id)
			default:
				return fmt.Errorf("KMS %s payment type %v is unknown; not released", id, detail["ChargeType"])
			}
			if err != nil {
				return err
			}
			return serviceVpcOwnerWaitGone(func() ([]map[string]interface{}, error) {
				// KMS list has no VPC field; verify it through GetKmsInstance
				// for the exact still-listed ID only.
				rows, err := serviceVpcOwnerList(rpc, "Kms", "2016-01-20", "ListKmsInstances", "$.KmsInstances.KmsInstance", "", false)
				if err != nil {
					return nil, err
				}
				for _, row := range rows {
					if row["KmsInstanceId"] != id {
						continue
					}
					response, err := rpc("Kms", "2016-01-20", "GetKmsInstance", nil, map[string]interface{}{"KmsInstanceId": id}, true)
					if err != nil {
						return nil, err
					}
					detail, ok := response["KmsInstance"].(map[string]interface{})
					if !ok || detail["InstanceId"] != id {
						return nil, fmt.Errorf("KMS %s invalid identity during release wait", id)
					}
					row["VpcId"] = detail["VpcId"]
				}
				return rows, nil
			}, "KmsInstanceId", id, vpcID, time.Sleep)
		})
	})
	add("alicloud_data_works_network", func() error {
		groups, err := serviceVpcDataWorksList(rpc, region, "ListResourceGroups", "$.PagingInfo.ResourceGroupList", "")
		if err != nil {
			return err
		}
		return serviceVpcOwnerEach(groups, func(group map[string]interface{}) error {
			groupID, err := serviceVpcDataWorksID(group["Id"])
			if err != nil {
				return err
			}
			networks, err := serviceVpcDataWorksList(rpc, region, "ListNetworks", "$.PagingInfo.NetworkList", groupID)
			if err != nil {
				return err
			}
			return serviceVpcOwnerEach(networks, func(network map[string]interface{}) error {
				match, err := serviceVpcOwnerMatch(network, "VpcId", vpcID)
				if err != nil || !match {
					return err
				}
				id, err := serviceVpcDataWorksID(network["Id"])
				if err != nil {
					return err
				}
				_, err = rpc("dataworks-public", "2024-05-18", "DeleteNetwork", nil, map[string]interface{}{"Id": id, "RegionId": region}, true)
				if err != nil {
					return err
				}
				return serviceVpcOwnerWaitGone(func() ([]map[string]interface{}, error) {
					rows, err := serviceVpcDataWorksList(rpc, region, "ListNetworks", "$.PagingInfo.NetworkList", groupID)
					if err != nil {
						return nil, err
					}
					for _, row := range rows {
						normalized, err := serviceVpcDataWorksID(row["Id"])
						if err != nil {
							return nil, err
						}
						row["Id"] = normalized
					}
					return rows, nil
				}, "Id", id, vpcID, time.Sleep)
			})
		})
	})
	add("alicloud_rds_custom", func() error {
		return sweepVpcRdsCustom(rpc, region, vpcID, time.Sleep)
	})
	add("alicloud_eci_virtual_node", func() error {
		return sweepVpcEciVirtualNodes(rpc, region, vpcID, time.Sleep)
	})
	return nodes
}

func vpcEciVirtualNodes(rpc vpcSweepRPC, region string) ([]map[string]interface{}, error) {
	request := map[string]interface{}{"RegionId": region, "Limit": 20}
	var all []map[string]interface{}
	seenIDs, seenTokens := map[string]bool{}, map[string]bool{}
	for page := 0; page < 1000; page++ {
		response, err := rpc("Eci", "2018-08-08", "DescribeVirtualNodes", nil, request, true)
		if err != nil {
			return nil, err
		}
		rows, err := serviceVpcOwnerRows("$.VirtualNodes", response)
		if err != nil || len(rows) > 20 {
			return nil, fmt.Errorf("ECI virtual node malformed inventory: %v", err)
		}
		for _, row := range rows {
			nodeID, err := serviceVpcOwnerID(row, "VirtualNodeId")
			vpc, known := row["VpcId"].(string)
			if err != nil || seenIDs[nodeID] || !known || vpc == "" || row["RegionId"] != region {
				return nil, fmt.Errorf("ECI virtual node duplicate or unknown ownership")
			}
			seenIDs[nodeID] = true
			all = append(all, row)
		}
		token, known := response["NextToken"].(string)
		if response["NextToken"] != nil && !known {
			return nil, fmt.Errorf("ECI virtual node invalid pagination token")
		}
		if token == "" {
			return all, nil
		}
		if seenTokens[token] {
			return nil, fmt.Errorf("ECI virtual node repeated pagination token")
		}
		seenTokens[token], request["NextToken"] = true, token
	}
	return nil, fmt.Errorf("ECI virtual node pagination limit exceeded")
}

func sweepVpcEciVirtualNodes(rpc vpcSweepRPC, region, vpcID string, pause func(time.Duration)) error {
	if !sweepAll() || region == "" || vpcID == "" {
		return fmt.Errorf("ECI virtual node release requires scoped SWEEP_ALL")
	}
	rows, err := vpcEciVirtualNodes(rpc, region)
	if err != nil {
		return err
	}
	return serviceVpcOwnerEach(rows, func(row map[string]interface{}) error {
		if row["VpcId"] != vpcID {
			return nil
		}
		id := row["VirtualNodeId"].(string)
		current, err := vpcEciVirtualNodes(rpc, region)
		if err != nil {
			return err
		}
		// The native filtered endpoint can return an empty result for a node
		// still present in the full inventory. Only a complete unfiltered
		// snapshot can prove exact identity or disappearance.
		var fresh []map[string]interface{}
		for _, item := range current {
			if item["VirtualNodeId"] == id {
				fresh = append(fresh, item)
			}
		}
		if len(fresh) == 0 {
			return nil
		}
		if len(fresh) != 1 || fresh[0]["VirtualNodeId"] != id || fresh[0]["VpcId"] != vpcID {
			return fmt.Errorf("ECI virtual node fresh identity or VPC changed")
		}
		node := fresh[0]
		if node["Status"] != "Ready" && node["Status"] != "Failed" && node["Status"] != "Pending" {
			return fmt.Errorf("ECI virtual node unknown lifecycle state")
		}
		primary, err := serviceVpcOwnerID(node, "VirtualNodeVSwitchId")
		if err != nil {
			return err
		}
		parent, err := rpc("Vpc", "2016-04-28", "DescribeVSwitchAttributes", nil, map[string]interface{}{"RegionId": region, "VSwitchId": primary}, true)
		if err != nil {
			return err
		}
		if parent["VSwitchId"] != primary || parent["VpcId"] != vpcID {
			return fmt.Errorf("ECI virtual node primary network is shared or unverifiable")
		}
		// DeleteVirtualNode atomically refuses a node that still has ECI
		// instances (VirtualNode.DELETE.FAILED). Never force or cascade into
		// workloads, shared scheduling subnets, pools, or ACK clusters.
		response, err := rpc("Eci", "2018-08-08", "DeleteVirtualNode", nil, map[string]interface{}{"RegionId": region, "VirtualNodeId": id}, false)
		if err != nil {
			return err
		}
		if requestID, ok := response["RequestId"].(string); !ok || requestID == "" {
			return fmt.Errorf("ECI virtual node release has no verified acknowledgement")
		}
		return serviceVpcOwnerWaitGone(func() ([]map[string]interface{}, error) {
			return vpcEciVirtualNodes(rpc, region)
		}, "VirtualNodeId", id, vpcID, pause)
	})
}

// RDS Custom has a separate inventory and release API from ordinary RDS.
// Collect the complete inventory before writes, then recheck exact ownership.
var errVpcRdsCustomInventoryChanging = fmt.Errorf("RDS Custom inventory changed during enumeration")

func vpcRdsCustomList(rpc vpcSweepRPC, region string, pause func(time.Duration)) ([]map[string]interface{}, error) {
	// Other VPC workers may release instances between the service's count and
	// row queries. Retry only this read-only consistency failure, starting a
	// complete new snapshot each time. Malformed identities still fail closed.
	for attempt := 0; attempt < 5; attempt++ {
		rows, err := vpcRdsCustomSnapshot(rpc, region)
		if err != errVpcRdsCustomInventoryChanging {
			return rows, err
		}
		if attempt < 4 {
			pause(time.Duration(1<<uint(attempt)) * time.Second)
		}
	}
	return nil, fmt.Errorf("RDS Custom inventory did not stabilize after five complete snapshots")
}

func vpcRdsCustomSnapshot(rpc vpcSweepRPC, region string) ([]map[string]interface{}, error) {
	var all []map[string]interface{}
	seen := map[string]bool{}
	expected := -1
	for page := 1; page <= 10000; page++ {
		response, err := rpc("Rds", "2014-08-15", "DescribeRCInstances", nil, map[string]interface{}{"RegionId": region, "PageNumber": page, "PageSize": 100}, true)
		if err != nil {
			return nil, err
		}
		var total int
		switch response["TotalCount"].(type) {
		case int, int64, float64, json.Number:
			text := fmt.Sprint(response["TotalCount"])
			value, err := strconv.Atoi(text)
			if err != nil || value < 0 || value > 1000000 || strconv.Itoa(value) != text {
				return nil, fmt.Errorf("RDS Custom invalid total")
			}
			total = value
		default:
			return nil, fmt.Errorf("RDS Custom unknown total type")
		}
		if expected >= 0 && expected != total {
			return nil, errVpcRdsCustomInventoryChanging
		}
		expected = total
		rows, err := serviceVpcOwnerRows("$.RCInstances", response)
		if err != nil || len(rows) > 100 {
			return nil, fmt.Errorf("RDS Custom incomplete or malformed inventory: %v", err)
		}
		if len(all)+len(rows) > total {
			return nil, errVpcRdsCustomInventoryChanging
		}
		for _, row := range rows {
			id, err := serviceVpcOwnerID(row, "InstanceId")
			_, knownVpc := row["VpcId"].(string)
			if err != nil || seen[id] || !knownVpc || row["RegionId"] != region {
				return nil, fmt.Errorf("RDS Custom duplicate or unverifiable inventory identity")
			}
			seen[id] = true
			all = append(all, row)
		}
		if len(all) == total {
			return all, nil
		}
		if len(rows) == 0 {
			return nil, errVpcRdsCustomInventoryChanging
		}
	}
	return nil, fmt.Errorf("RDS Custom pagination limit exceeded")
}

func sweepVpcRdsCustom(rpc vpcSweepRPC, region, vpcID string, pause func(time.Duration)) error {
	if !sweepAll() || region == "" || vpcID == "" {
		return fmt.Errorf("RDS Custom scoped release requires explicit SWEEP_ALL")
	}
	rows, err := vpcRdsCustomList(rpc, region, pause)
	if err != nil {
		return err
	}
	return serviceVpcOwnerEach(rows, func(row map[string]interface{}) error {
		if row["VpcId"] != vpcID {
			return nil
		}
		id := row["InstanceId"].(string)
		detail, err := rpc("Rds", "2014-08-15", "DescribeRCInstanceAttribute", nil, map[string]interface{}{"RegionId": region, "InstanceId": id}, true)
		if err != nil {
			return err
		}
		network, knownNetwork := detail["VpcAttributes"].(map[string]interface{})
		if detail["InstanceId"] != id || detail["RegionId"] != region || !knownNetwork || network["VpcId"] != vpcID || detail["ClusterId"] != "" {
			return fmt.Errorf("RDS Custom identity/network/cluster association changed or unknown; not released")
		}
		charge, knownCharge := detail["InstanceChargeType"].(string)
		protected, knownProtection := detail["DeletionProtection"].(bool)
		if !knownProtection || !knownCharge || (charge != "PostPaid" && charge != "PrePaid") {
			return fmt.Errorf("RDS Custom unknown billing or deletion protection; not released")
		}
		verify := func() error {
			fresh, err := rpc("Rds", "2014-08-15", "DescribeRCInstanceAttribute", nil, map[string]interface{}{"RegionId": region, "InstanceId": id}, true)
			if err != nil {
				return err
			}
			net, ok := fresh["VpcAttributes"].(map[string]interface{})
			if !ok || fresh["InstanceId"] != id || fresh["RegionId"] != region || net["VpcId"] != vpcID || fresh["ClusterId"] != "" || fresh["InstanceChargeType"] != charge || fresh["DeletionProtection"] != false {
				return fmt.Errorf("RDS Custom fresh scope/billing/protection changed; not released")
			}
			return nil
		}
		if protected {
			// The official API supports this flag even when older CLI catalogs
			// omit the action. Never change it before exact VPC verification.
			_, err = rpc("Rds", "2014-08-15", "ModifyRCInstanceAttribute", nil, map[string]interface{}{"RegionId": region, "InstanceId": id, "DeletionProtection": false}, false)
			if err != nil {
				return err
			}
			if err := verify(); err != nil {
				return err
			}
		}
		if charge == "PrePaid" {
			err = refundVpcSubscription(rpc, region, "rds", id, verify)
		} else {
			if err := verify(); err != nil {
				return err
			}
			// The official InstanceId selector uses JSON array serialization,
			// not RPC repeated fields. Keep the request to exactly one instance.
			ids, _ := json.Marshal([]string{id})
			_, err = rpc("Rds", "2014-08-15", "DeleteRCInstances", nil, map[string]interface{}{"RegionId": region, "InstanceId": string(ids), "Force": true}, false)
		}
		if err != nil {
			return err
		}
		return serviceVpcOwnerWaitGone(func() ([]map[string]interface{}, error) {
			return vpcRdsCustomList(rpc, region, pause)
		}, "InstanceId", id, vpcID, pause)
	})
}

func serviceVpcOwnerMatch(item map[string]interface{}, field, vpcID string) (bool, error) {
	value, ok := item[field].(string)
	if !ok {
		return false, fmt.Errorf("missing or invalid %s; refusing unscoped deletion", field)
	}
	return value == vpcID, nil
}

func serviceVpcOwnerID(item map[string]interface{}, field string) (string, error) {
	value, ok := item[field].(string)
	if !ok || value == "" {
		return "", fmt.Errorf("missing or invalid %s", field)
	}
	return value, nil
}

func serviceVpcOwnerRows(path string, response map[string]interface{}) ([]map[string]interface{}, error) {
	value, err := jsonpath.Get(path, response)
	if err != nil {
		return nil, fmt.Errorf("invalid service owner list %s: %w", path, err)
	}
	rows, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid service owner list %s", path)
	}
	items := make([]map[string]interface{}, 0, len(rows))
	for _, row := range rows {
		item, ok := row.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid service owner item %s", path)
		}
		items = append(items, item)
	}
	return items, nil
}

func serviceVpcOwnerList(rpc vpcSweepRPC, product, version, action, path, region string, tokenPaging bool) ([]map[string]interface{}, error) {
	request := map[string]interface{}{}
	if region != "" {
		request["RegionId"] = region
	}
	if tokenPaging {
		request["MaxResults"] = PageSizeLarge
	} else {
		request["PageSize"], request["PageNumber"] = PageSizeLarge, 1
	}
	var all []map[string]interface{}
	seen := map[string]bool{}
	for page := 0; ; page++ {
		if page >= 1000 {
			return nil, fmt.Errorf("%s pagination limit exceeded", action)
		}
		response, err := rpc(product, version, action, nil, request, true)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", action, err)
		}
		rows, err := serviceVpcOwnerRows(path, response)
		if err != nil {
			return nil, err
		}
		if len(rows) > 0 {
			fingerprint := "rows:" + fmt.Sprint(rows)
			if seen[fingerprint] {
				return nil, fmt.Errorf("%s repeated resource page", action)
			}
			seen[fingerprint] = true
		}
		all = append(all, rows...)
		if tokenPaging {
			token, _ := response["NextToken"].(string)
			if token == "" {
				break
			}
			if seen[token] {
				return nil, fmt.Errorf("%s repeated pagination token", action)
			}
			seen[token], request["NextToken"] = true, token
		} else {
			if len(rows) < PageSizeLarge {
				break
			}
			request["PageNumber"] = request["PageNumber"].(int) + 1
		}
	}
	return all, nil
}

func serviceVpcOwnerEach(items []map[string]interface{}, run func(map[string]interface{}) error) error {
	var failures []string
	for _, item := range items {
		if err := run(item); err != nil {
			failures = append(failures, err.Error())
		}
	}
	if len(failures) > 0 {
		return fmt.Errorf("service owner cleanup incomplete: %s", strings.Join(failures, "; "))
	}
	return nil
}

// Poll only IDs actually released by this node. A successful empty list is
// evidence of disappearance; read failures, malformed IDs or changed ownership
// never count as gone. The injected pause keeps the bounded wait testable.
func serviceVpcOwnerWaitGone(list func() ([]map[string]interface{}, error), idKey, id, vpcID string, pause func(time.Duration)) error {
	for attempt := 0; attempt < 120; attempt++ {
		rows, err := list()
		if err != nil {
			return err
		}
		present := false
		for _, row := range rows {
			rowID, err := serviceVpcOwnerID(row, idKey)
			if err != nil {
				return err
			}
			if rowID != id {
				continue
			}
			match, err := serviceVpcOwnerMatch(row, "VpcId", vpcID)
			if err != nil || !match {
				return fmt.Errorf("%s %s has unverifiable or changed VPC during release wait: %v", idKey, id, err)
			}
			present = true
		}
		if !present {
			return nil
		}
		if attempt < 119 {
			pause(5 * time.Second)
		}
	}
	return fmt.Errorf("%s %s still present after bounded release wait", idKey, id)
}

// Group ownership alone is insufficient: one resource group can contain
// networks in protected VPCs. Only the exact target network is removed. Group
// subscription/cancellation metadata is not established by these list schemas,
// so this node deliberately never calls DeleteResourceGroup or billing APIs.
func serviceVpcDataWorksList(rpc vpcSweepRPC, region, action, path, groupID string) ([]map[string]interface{}, error) {
	var all []map[string]interface{}
	for page := 1; page <= 10000; page++ {
		request := map[string]interface{}{"RegionId": region, "PageNumber": page, "PageSize": PageSizeLarge}
		if groupID != "" {
			request["ResourceGroupId"] = groupID
		}
		response, err := rpc("dataworks-public", "2024-05-18", action, nil, request, true)
		if err != nil {
			return nil, err
		}
		rows, err := serviceVpcOwnerRows(path, response)
		if err != nil {
			return nil, err
		}
		all = append(all, rows...)
		if len(rows) < PageSizeLarge {
			return all, nil
		}
	}
	return nil, fmt.Errorf("%s exceeded bounded pagination", action)
}

func serviceVpcDataWorksID(value interface{}) (string, error) {
	switch id := value.(type) {
	case json.Number:
		// DataWorks network IDs are positive int64 values. SDK UseNumber must
		// preserve their exact decimal identity without float rounding or coercion.
		parsed, err := strconv.ParseInt(string(id), 10, 64)
		if err == nil && parsed > 0 && strconv.FormatInt(parsed, 10) == string(id) {
			return string(id), nil
		}
	case string:
		if id != "" {
			return id, nil
		}
	case float64:
		if id > 0 && id <= 9007199254740991 && math.Trunc(id) == id {
			return strconv.FormatFloat(id, 'f', 0, 64), nil
		}
	case int:
		if id > 0 {
			return strconv.Itoa(id), nil
		}
	case int64:
		if id > 0 {
			return strconv.FormatInt(id, 10), nil
		}
	}
	return "", fmt.Errorf("invalid DataWorks resource identity; refusing deletion")
}

// The official SDK models BindVpcs.BindVpc as []BindVpc. A returned explicit
// null is the JSON representation of that nil slice, i.e. a known empty list.
// Missing wrapper/key and all malformed non-list values remain unknown.
func kmsVpcOwnerBindings(detail map[string]interface{}) ([]map[string]interface{}, error) {
	wrapper, ok := detail["BindVpcs"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unknown VPC binding wrapper; not released")
	}
	value, exists := wrapper["BindVpc"]
	if !exists {
		return nil, fmt.Errorf("missing VPC binding list; not released")
	}
	if value == nil {
		return []map[string]interface{}{}, nil
	}
	rows, err := serviceVpcOwnerRows("$.BindVpcs.BindVpc", detail)
	if err != nil {
		return nil, fmt.Errorf("invalid VPC binding list; not released: %w", err)
	}
	return rows, nil
}

// Subscription cleanup uses the native billing refund route only after exact
// billing identity, successful inquiry and fresh service-side VPC verification.
// A refund order is acceptance, never proof that the service instance is gone.
func refundVpcKmsSubscription(rpc vpcSweepRPC, region, vpcID, id string) error {
	return refundVpcSubscription(rpc, region, "kms", id, func() error {
		fresh, err := rpc("Kms", "2016-01-20", "GetKmsInstance", nil, map[string]interface{}{"KmsInstanceId": id}, true)
		if err != nil {
			return err
		}
		detail, ok := fresh["KmsInstance"].(map[string]interface{})
		if !ok || detail["InstanceId"] != id || detail["VpcId"] != vpcID || detail["ChargeType"] != "PREPAY" || detail["DeletionProtection"] == true {
			return fmt.Errorf("KMS subscription identity/ownership changed; not refunded")
		}
		bindings, err := kmsVpcOwnerBindings(detail)
		if err != nil {
			return err
		}
		for _, binding := range bindings {
			if binding["VpcId"] != vpcID {
				return fmt.Errorf("KMS subscription has shared VPC bindings; not refunded")
			}
		}
		return nil
	})
}

func refundVpcSubscription(rpc vpcSweepRPC, region, product, id string, verify func() error) error {
	if !sweepAll() || verify == nil || (product != "kms" && product != "redisa" && product != "rds") {
		return fmt.Errorf("subscription refund requires explicit SWEEP_ALL and verified product/scope; not released")
	}
	response, err := rpc("BssOpenApi", "2017-12-14", "QueryAvailableInstances", nil, map[string]interface{}{"InstanceIDs": id, "PageNum": 1, "PageSize": 100}, true)
	if err != nil {
		return err
	}
	if response["Success"] != true || response["Code"] != "Success" {
		return fmt.Errorf("subscription billing inventory failed: %v", response["Code"])
	}
	rows, err := serviceVpcOwnerRows("$.Data.InstanceList", response)
	if err != nil || len(rows) != 1 {
		return fmt.Errorf("subscription billing inventory is not exactly one instance: %v", err)
	}
	row := rows[0]
	productType, ok := row["ProductType"].(string)
	if product == "redisa" && row["ProductType"] == nil {
		productType, ok = "", true
	}
	if row["InstanceID"] != id || row["Region"] != region || row["ProductCode"] != product || row["SubscriptionType"] != "Subscription" || !ok || (product == "kms" && productType == "") {
		return fmt.Errorf("subscription billing identity/region/subscription mismatch; not refunded")
	}
	if product == "rds" && productType != "rds_customprepaid_public_cn" {
		return fmt.Errorf("subscription billing product is not RDS Custom; not refunded")
	}
	endTime, knownEnd := row["EndTime"].(string)
	if _, err := time.Parse(time.RFC3339, endTime); !knownEnd || err != nil {
		return fmt.Errorf("subscription billing period is unknown; not refunded")
	}
	data, ok := response["Data"].(map[string]interface{})
	if !ok || fmt.Sprint(data["TotalCount"]) != "1" || fmt.Sprint(data["PageNum"]) != "1" {
		return fmt.Errorf("subscription billing pagination is incomplete; not refunded")
	}
	request := map[string]interface{}{"InstanceId": id, "ProductCode": product, "ProductType": productType}
	inquiry, err := rpc("BssOpenApi", "2017-12-14", "InquiryPriceRefundInstance", nil, request, true)
	if err != nil {
		return err
	}
	price, ok := inquiry["Data"].(map[string]interface{})
	if inquiry["Success"] != true || inquiry["Code"] != "Success" || !ok || price["InstanceId"] != id {
		return fmt.Errorf("subscription refund inquiry failed or returned another instance: %v", inquiry["Code"])
	}
	switch price["RefundAmount"].(type) {
	case json.Number, float64, int, int64:
	default:
		return fmt.Errorf("subscription refund inquiry amount has unknown native type; not refunded")
	}
	amount, validAmount := new(big.Rat).SetString(fmt.Sprint(price["RefundAmount"]))
	currency, validCurrency := price["Currency"].(string)
	for _, letter := range currency {
		if letter < 'A' || letter > 'Z' {
			validCurrency = false
		}
	}
	if !validAmount || amount.Sign() < 0 || !validCurrency || len(currency) != 3 {
		return fmt.Errorf("subscription refund inquiry amount/currency is invalid; not refunded")
	}
	if err := verify(); err != nil {
		return err
	}
	// Keep one idempotency key per billed subscription period. Ambiguous
	// responses retry the same order; a later renewal gets a new period/key.
	token := sha256.Sum256([]byte("vpc-subscription-refund|" + region + "|" + product + "|" + id + "|" + endTime))
	request["ClientToken"] = fmt.Sprintf("vpc-%s-%x", product, token[:20])
	request["ImmediatelyRelease"] = "1"
	refund, err := rpc("BssOpenApi", "2017-12-14", "RefundInstance", nil, request, false)
	if err != nil {
		return err
	}
	order, ok := refund["Data"].(map[string]interface{})
	if refund["Success"] != true || refund["Code"] != "Success" || !ok {
		return fmt.Errorf("subscription refund rejected: %v", refund["Code"])
	}
	orderID, err := serviceVpcDataWorksID(order["OrderId"])
	parsedOrder, parseErr := strconv.ParseInt(orderID, 10, 64)
	if err != nil || parseErr != nil || parsedOrder <= 0 || strconv.FormatInt(parsedOrder, 10) != orderID {
		return fmt.Errorf("subscription refund returned no verified positive order ID")
	}
	log.Printf("[INFO] SWEEP_ALL %s subscription %s refund accepted; waiting for native release confirmation", product, id)
	return nil
}

// Redis subscription cleanup reuses the billing contract. Private endpoints and
// global-distributed membership must be verified before releasing the instance.
func refundVpcKvstoreSubscription(rpc vpcSweepRPC, region, vpcID, id string, owner allVpcOwner) error {
	return refundVpcSubscription(rpc, region, "redisa", id, func() error {
		items, err := allVpcOwnerItems(rpc, region, vpcID, owner)
		if err != nil {
			return err
		}
		matches := 0
		for _, item := range items {
			if item["InstanceId"] != id {
				continue
			}
			matches++
			global, known := item["GlobalInstanceId"].(string)
			if !known || global != "" {
				return fmt.Errorf("Redis global-distributed membership is unknown or shared; not refunded")
			}
		}
		if matches != 1 {
			return fmt.Errorf("Redis fresh regional inventory has no unique target; not refunded")
		}
		response, err := rpc("R-kvstore", "2015-01-01", "DescribeInstanceAttribute", nil, map[string]interface{}{"RegionId": region, "InstanceId": id}, true)
		if err != nil {
			return err
		}
		rows, err := serviceVpcOwnerRows("$.Instances.DBInstanceAttribute", response)
		if err != nil || len(rows) != 1 {
			return fmt.Errorf("Redis subscription detail is incomplete; not refunded")
		}
		detail := rows[0]
		if detail["InstanceId"] != id || detail["VpcId"] != vpcID || detail["RegionId"] != region || detail["ChargeType"] != "PrePaid" || detail["InstanceReleaseProtection"] != false {
			return fmt.Errorf("Redis subscription ownership/payment/protection is unknown or changed; not refunded")
		}
		network, err := rpc("R-kvstore", "2015-01-01", "DescribeDBInstanceNetInfo", nil, map[string]interface{}{"RegionId": region, "InstanceId": id}, true)
		if err != nil {
			return err
		}
		endpoints, err := serviceVpcOwnerRows("$.NetInfoItems.InstanceNetInfo", network)
		if err != nil {
			return err
		}
		for _, endpoint := range endpoints {
			if endpoint["IPType"] == "Public" {
				continue
			}
			if endpoint["IPType"] != "Private" || endpoint["VPCId"] != vpcID {
				return fmt.Errorf("Redis private endpoint is shared or unknown; not refunded")
			}
		}
		return nil
	})
}

func TestUnitVpcSweepEciNativeTransportUsesQueryParameters(t *testing.T) {
	var calls []string
	get := func(product, version, action string, query, body map[string]interface{}) (map[string]interface{}, error) {
		if body != nil || query["RegionId"] != "region" {
			t.Fatal("read parameters were not placed in URL query")
		}
		switch product {
		case "Eci":
			if version != "2018-08-08" || action != "DescribeVirtualNodes" || query["Limit"] != 20 {
				t.Fatal("incorrect ECI GET contract")
			}
		case "dataworks-public":
			if version != "2024-05-18" || action != "ListResourceGroups" {
				t.Fatal("DataWorks route changed")
			}
		default:
			t.Fatal("unexpected GET route")
		}
		calls = append(calls, product+":GET")
		return map[string]interface{}{}, nil
	}
	post := func(product, version, action string, query, body map[string]interface{}, retry bool) (map[string]interface{}, error) {
		if product == "Eci" {
			if body != nil || query["RegionId"] != "region" || query["VirtualNodeId"] != "target" || len(query) != 2 || retry || action != "DeleteVirtualNode" {
				t.Fatal("ECI deletion query/retry contract changed")
			}
		} else if product != "Rds" || query != nil || body["RegionId"] != "region" || !retry || action != "DescribeDBInstances" {
			t.Fatal("ordinary RPC POST passthrough changed")
		}
		calls = append(calls, product+":POST")
		return map[string]interface{}{}, nil
	}
	router := routedVpcSweepRPC(get, post)
	for _, request := range []struct {
		product, version, action string
		body                     map[string]interface{}
		retry                    bool
	}{
		{"Eci", "2018-08-08", "DescribeVirtualNodes", map[string]interface{}{"RegionId": "region", "Limit": 20}, true},
		{"Eci", "2018-08-08", "DeleteVirtualNode", map[string]interface{}{"RegionId": "region", "VirtualNodeId": "target"}, false},
		{"dataworks-public", "2024-05-18", "ListResourceGroups", map[string]interface{}{"RegionId": "region"}, true},
		{"Rds", "2014-08-15", "DescribeDBInstances", map[string]interface{}{"RegionId": "region"}, true},
	} {
		if _, err := router(request.product, request.version, request.action, nil, request.body, request.retry); err != nil {
			t.Fatal(err)
		}
	}
	if len(calls) != 4 {
		t.Fatalf("missing transport calls: %v", calls)
	}
	if _, err := router("Eci", "2018-08-08", "DescribeVirtualNodes", map[string]interface{}{"RegionId": "region"}, map[string]interface{}{"RegionId": "other"}, true); err == nil {
		t.Fatal("ambiguous query/body parameters accepted")
	}
	if _, err := router("Eci", "wrong-version", "DescribeVirtualNodes", nil, nil, true); err == nil || len(calls) != 4 {
		t.Fatal("unknown version reached transport")
	}
}

func eciVirtualNodeFixture(mode string, writes *[]string) vpcSweepRPC {
	globalReads := 0
	makeNode := func(id, vpc string) map[string]interface{} {
		return map[string]interface{}{"VirtualNodeId": id, "VpcId": vpc, "RegionId": "region", "Status": "Ready", "VirtualNodeVSwitchId": "primary"}
	}
	return func(product, version, action string, query, request map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
		if request["RegionId"] != "region" {
			return nil, fmt.Errorf("wrong transport region")
		}
		switch action {
		case "DescribeVirtualNodes":
			if product != "Eci" || version != "2018-08-08" || request["Limit"] != 20 {
				return nil, fmt.Errorf("incorrect ECI page contract")
			}
			node := makeNode("target-id", "target")
			if request["VirtualNodeIds"] != nil {
				return nil, fmt.Errorf("filtered empty response is not authoritative inventory")
			}
			globalReads++
			if globalReads == 2 && len(*writes) == 0 {
				switch mode {
				case "fresh-vpc":
					node["VpcId"] = "protected"
				case "fresh-id":
					node["VirtualNodeId"] = "other"
				case "missing-primary":
					delete(node, "VirtualNodeVSwitchId")
				case "unknown-status":
					node["Status"] = "Unrecognized"
				}
			}
			rows := []interface{}{makeNode("other", "protected")}
			if len(*writes) == 0 || mode == "wait-never-gone" || mode == "wait-moved" || (mode == "async" && globalReads == 3) {
				if mode == "wait-moved" && len(*writes) > 0 {
					node["VpcId"] = "protected"
				}
				if mode == "missing-vpc" {
					delete(node, "VpcId")
				}
				if mode == "wrong-region" {
					node["RegionId"] = "reserved-region"
				}
				rows = append(rows, node)
				if mode == "duplicate" {
					rows = append(rows, node)
				}
			}
			response := map[string]interface{}{"VirtualNodes": rows}
			if mode == "invalid-token" {
				response["NextToken"] = 1
			}
			return response, nil
		case "DescribeVSwitchAttributes":
			if product != "Vpc" || version != "2016-04-28" || request["VSwitchId"] != "primary" {
				return nil, fmt.Errorf("wrong primary network query")
			}
			vpc := "target"
			if mode == "shared-primary" {
				vpc = "protected"
			}
			return map[string]interface{}{"VSwitchId": "primary", "VpcId": vpc}, nil
		case "DeleteVirtualNode":
			if product != "Eci" || version != "2018-08-08" || request["VirtualNodeId"] != "target-id" || len(request) != 2 || autoRetry {
				return nil, fmt.Errorf("unsafe force/batch/cascade delete parameters")
			}
			*writes = append(*writes, "target-id")
			switch mode {
			case "busy":
				return nil, &tea.SDKError{Code: tea.String("VirtualNode.DELETE.FAILED"), Message: tea.String("node still has ECI instances")}
			case "write-error":
				return nil, fmt.Errorf("ambiguous response")
			case "missing-ack":
				return map[string]interface{}{}, nil
			}
			return map[string]interface{}{"RequestId": "request"}, nil
		}
		return nil, fmt.Errorf("unexpected cascade action %s", action)
	}
}

func TestUnitVpcSweepEciVirtualNodeNativeScopeAndAsyncRelease(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	var writes []string
	waits := 0
	err := sweepVpcEciVirtualNodes(eciVirtualNodeFixture("async", &writes), "region", "target", func(time.Duration) { waits++ })
	if err != nil || len(writes) != 1 || waits != 1 {
		t.Fatalf("writes=%v waits=%d err=%v", writes, waits, err)
	}
}

func TestUnitVpcSweepEciVirtualNodeFilteredEmptyIsNotAbsenceProof(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	var writes []string
	base := eciVirtualNodeFixture("immediate", &writes)
	reads := 0
	rpc := func(product, version, action string, query, request map[string]interface{}, retry bool) (map[string]interface{}, error) {
		if action == "DescribeVirtualNodes" {
			reads++
			if request["VirtualNodeIds"] != nil {
				// Reproduce the real service: filtered lookup returns empty,
				// although its complete regional inventory still contains it.
				return map[string]interface{}{"VirtualNodes": []interface{}{}}, nil
			}
		}
		return base(product, version, action, query, request, retry)
	}
	if err := sweepVpcEciVirtualNodes(rpc, "region", "target", func(time.Duration) {}); err != nil || len(writes) != 1 || reads < 3 {
		t.Fatalf("filtered absence suppressed cleanup: reads=%d writes=%v err=%v", reads, writes, err)
	}
}

func TestUnitVpcSweepEciVirtualNodeSharedUnknownNeverDeletes(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	for _, mode := range []string{"missing-vpc", "wrong-region", "duplicate", "invalid-token", "fresh-vpc", "fresh-id", "missing-primary", "unknown-status", "shared-primary"} {
		t.Run(mode, func(t *testing.T) {
			var writes []string
			err := sweepVpcEciVirtualNodes(eciVirtualNodeFixture(mode, &writes), "region", "target", func(time.Duration) {})
			if err == nil || len(writes) != 0 {
				t.Fatalf("unsafe release writes=%v err=%v", writes, err)
			}
		})
	}
}

func TestUnitVpcSweepEciVirtualNodeBusyAndUnconfirmedNeverPass(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	for _, mode := range []string{"busy", "write-error", "missing-ack", "wait-never-gone", "wait-moved"} {
		t.Run(mode, func(t *testing.T) {
			var writes []string
			waits := 0
			err := sweepVpcEciVirtualNodes(eciVirtualNodeFixture(mode, &writes), "region", "target", func(time.Duration) { waits++ })
			if err == nil || len(writes) != 1 || waits > 119 {
				t.Fatalf("writes=%v waits=%d err=%v", writes, waits, err)
			}
		})
	}
}

func TestUnitVpcSweepEciVirtualNodeCompletePaginationAndOrdinaryMode(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "false")
	if sweepVpcEciVirtualNodes(nil, "region", "target", nil) == nil {
		t.Fatal("ordinary cleanup reached full-cleanup node")
	}
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	for _, repeated := range []bool{false, true} {
		var writes []string
		base := eciVirtualNodeFixture("immediate", &writes)
		pages := 0
		rpc := func(product, version, action string, query, request map[string]interface{}, retry bool) (map[string]interface{}, error) {
			if action == "DescribeVirtualNodes" && len(writes) == 0 && pages < 2 {
				pages++
				if pages == 1 {
					return map[string]interface{}{"VirtualNodes": []interface{}{}, "NextToken": "next"}, nil
				}
				if request["NextToken"] != "next" {
					t.Fatal("did not follow native NextToken")
				}
				if repeated {
					return map[string]interface{}{"VirtualNodes": []interface{}{}, "NextToken": "next"}, nil
				}
			}
			return base(product, version, action, query, request, retry)
		}
		err := sweepVpcEciVirtualNodes(rpc, "region", "target", func(time.Duration) {})
		if repeated {
			if err == nil || len(writes) != 0 {
				t.Fatalf("repeated token accepted writes=%v err=%v", writes, err)
			}
		} else if err != nil || len(writes) != 1 || pages != 2 {
			t.Fatalf("pages=%d writes=%v err=%v", pages, writes, err)
		}
	}
}

func rdsCustomSweepFixture(mode string, writes *[]string) vpcSweepRPC {
	lists := 0
	return func(product, version, action string, query, request map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
		if product != "Rds" || version != "2014-08-15" || request["RegionId"] != "region" {
			return nil, fmt.Errorf("unexpected RDS transport scope")
		}
		switch action {
		case "DescribeRCInstances":
			lists++
			rows := []interface{}{map[string]interface{}{"InstanceId": "other", "RegionId": "region", "VpcId": "protected"}}
			if len(*writes) == 0 || mode == "wait-never-gone" || (mode == "async" && lists == 2) || mode == "wait-moved" {
				vpc := "target"
				if mode == "wait-moved" && len(*writes) > 0 {
					vpc = "protected"
				}
				row := map[string]interface{}{"InstanceId": "target-id", "RegionId": "region", "VpcId": vpc}
				if mode == "missing-vpc" {
					delete(row, "VpcId")
				}
				rows = append(rows, row)
				if mode == "duplicate" {
					rows = append(rows, row)
				}
			}
			total := interface{}(json.Number(fmt.Sprint(len(rows))))
			if mode == "fraction-total" {
				total = json.Number("2.5")
			}
			if mode == "empty-page" {
				rows = []interface{}{}
				total = json.Number("2")
			}
			return map[string]interface{}{"TotalCount": total, "RCInstances": rows}, nil
		case "DescribeRCInstanceAttribute":
			if request["InstanceId"] != "target-id" {
				return nil, fmt.Errorf("described protected instance")
			}
			if mode == "detail-error" {
				return nil, fmt.Errorf("read denied")
			}
			d := map[string]interface{}{"InstanceId": "target-id", "RegionId": "region", "VpcAttributes": map[string]interface{}{"VpcId": "target"}, "InstanceChargeType": "PostPaid", "DeletionProtection": false, "ClusterId": ""}
			switch mode {
			case "fresh-vpc":
				d["VpcAttributes"] = map[string]interface{}{"VpcId": "protected"}
			case "fresh-region":
				d["RegionId"] = "reserved-region"
			case "fresh-id":
				d["InstanceId"] = "other"
			case "cluster":
				d["ClusterId"] = "shared-cluster"
			case "missing-cluster":
				delete(d, "ClusterId")
			case "deletion-protection":
				d["DeletionProtection"] = true
			case "missing-protection":
				delete(d, "DeletionProtection")
			case "subscription":
				d["InstanceChargeType"] = "PrePaid"
			case "unknown-payment":
				delete(d, "InstanceChargeType")
			}
			return d, nil
		case "DeleteRCInstances":
			encoded, ok := request["InstanceId"].(string)
			var ids []string
			if !ok || json.Unmarshal([]byte(encoded), &ids) != nil || len(ids) != 1 || ids[0] != "target-id" || request["Force"] != true || len(request) != 3 || autoRetry {
				return nil, fmt.Errorf("unsafe batch/retry/force release parameters")
			}
			*writes = append(*writes, "target-id")
			if mode == "write-error" {
				return nil, fmt.Errorf("ambiguous write response")
			}
			return map[string]interface{}{"RequestId": "request"}, nil
		}
		return nil, fmt.Errorf("unexpected action %s", action)
	}
}

func TestUnitVpcSweepRdsCustomNativeScopedAsyncRelease(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	var writes []string
	waits := 0
	err := sweepVpcRdsCustom(rdsCustomSweepFixture("async", &writes), "region", "target", func(time.Duration) { waits++ })
	if err != nil || len(writes) != 1 || waits != 1 {
		t.Fatalf("writes=%v waits=%d err=%v", writes, waits, err)
	}
}

func TestUnitVpcSweepRdsCustomUnknownSharedAndBillingNeverRelease(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	for _, mode := range []string{"missing-vpc", "duplicate", "fraction-total", "empty-page", "detail-error", "fresh-vpc", "fresh-region", "fresh-id", "cluster", "missing-cluster", "deletion-protection", "missing-protection", "subscription", "unknown-payment"} {
		t.Run(mode, func(t *testing.T) {
			var writes []string
			err := sweepVpcRdsCustom(rdsCustomSweepFixture(mode, &writes), "region", "target", func(time.Duration) {})
			if err == nil || len(writes) != 0 {
				t.Fatalf("unsafe release writes=%v err=%v", writes, err)
			}
		})
	}
}

func TestUnitVpcSweepRdsCustomErrorsAndOrderNotReleaseProof(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	for _, mode := range []string{"write-error", "wait-never-gone", "wait-moved"} {
		t.Run(mode, func(t *testing.T) {
			var writes []string
			waits := 0
			err := sweepVpcRdsCustom(rdsCustomSweepFixture(mode, &writes), "region", "target", func(time.Duration) { waits++ })
			if err == nil || len(writes) != 1 || waits > 119 {
				t.Fatalf("writes=%v waits=%d err=%v", writes, waits, err)
			}
		})
	}
}

func TestUnitVpcSweepRdsCustomOrdinaryCleanupUnchanged(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "false")
	if sweepVpcRdsCustom(nil, "region", "target", nil) == nil {
		t.Fatal("ordinary cleanup reached full-cleanup node")
	}
}

func TestUnitVpcSweepRdsCustomConcurrentInventoryRestartsBeforeWrites(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	for _, mode := range []string{"under-count", "empty-second-page", "changed-second-page"} {
		t.Run(mode, func(t *testing.T) {
			var writes []string
			base := rdsCustomSweepFixture("immediate", &writes)
			reads, waits := 0, 0
			rpc := func(product, version, action string, query, request map[string]interface{}, retry bool) (map[string]interface{}, error) {
				d, err := base(product, version, action, query, request, retry)
				if action != "DescribeRCInstances" || err != nil {
					return d, err
				}
				reads++
				if reads == 1 {
					if len(writes) != 0 {
						t.Fatal("mutation preceded a complete stable snapshot")
					}
					if mode == "under-count" {
						d["TotalCount"] = json.Number("1")
					} else {
						d["RCInstances"] = d["RCInstances"].([]interface{})[:1]
					}
				} else if reads == 2 && mode != "under-count" {
					if request["PageNumber"] != 2 || len(writes) != 0 {
						t.Fatal("did not finish the read-only first snapshot")
					}
					d["RCInstances"] = []interface{}{}
					if mode == "changed-second-page" {
						d["TotalCount"] = json.Number("1")
					}
				} else if len(writes) == 0 && request["PageNumber"] != 1 {
					t.Fatal("retry did not restart pagination")
				}
				return d, nil
			}
			err := sweepVpcRdsCustom(rpc, "region", "target", func(time.Duration) { waits++ })
			if err != nil || len(writes) != 1 || waits != 1 {
				t.Fatalf("reads=%d writes=%v waits=%d err=%v", reads, writes, waits, err)
			}
		})
	}
}

func TestUnitVpcSweepRdsCustomNativeUnlockAndFreshOwnership(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	for _, mode := range []string{"success", "still-protected", "moved-after-unlock", "unlock-error", "protection-appeared"} {
		t.Run(mode, func(t *testing.T) {
			var writes []string
			base := rdsCustomSweepFixture("immediate", &writes)
			unlocks, details := 0, 0
			rpc := func(product, version, action string, query, request map[string]interface{}, retry bool) (map[string]interface{}, error) {
				if action == "ModifyRCInstanceAttribute" {
					if request["InstanceId"] != "target-id" || request["RegionId"] != "region" || request["DeletionProtection"] != false || len(request) != 3 || retry {
						t.Fatal("unsafe unlock parameters")
					}
					unlocks++
					if mode == "unlock-error" {
						return nil, fmt.Errorf("unlock rejected")
					}
					return map[string]interface{}{}, nil
				}
				d, err := base(product, version, action, query, request, retry)
				if action == "DescribeRCInstanceAttribute" && err == nil {
					details++
					d["DeletionProtection"] = unlocks == 0 || mode == "still-protected"
					if mode == "protection-appeared" {
						d["DeletionProtection"] = details > 1
					}
					if mode == "moved-after-unlock" && unlocks > 0 {
						d["VpcAttributes"] = map[string]interface{}{"VpcId": "protected"}
					}
				}
				return d, err
			}
			err := sweepVpcRdsCustom(rpc, "region", "target", func(time.Duration) {})
			if mode == "success" {
				if err != nil || len(writes) != 1 || unlocks != 1 {
					t.Fatalf("writes=%v unlocks=%d err=%v", writes, unlocks, err)
				}
			} else if err == nil || len(writes) != 0 {
				t.Fatalf("unsafe release writes=%v err=%v", writes, err)
			}
		})
	}
}

func TestUnitVpcSweepRdsCustomSubscriptionRefundRequiresExactProduct(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	for _, mode := range []string{"success", "other-product-type", "inquiry-rejected", "fresh-vpc", "refund-response-lost"} {
		t.Run(mode, func(t *testing.T) {
			var writes []string
			base := rdsCustomSweepFixture("immediate", &writes)
			inquired, refunds := false, 0
			rpc := func(product, version, action string, query, request map[string]interface{}, retry bool) (map[string]interface{}, error) {
				if product == "BssOpenApi" {
					switch action {
					case "QueryAvailableInstances":
						kind := "rds_customprepaid_public_cn"
						if mode == "other-product-type" {
							kind = "rds"
						}
						return map[string]interface{}{"Success": true, "Code": "Success", "Data": map[string]interface{}{"TotalCount": json.Number("1"), "PageNum": json.Number("1"), "InstanceList": []interface{}{map[string]interface{}{"InstanceID": "target-id", "Region": "region", "ProductCode": "rds", "ProductType": kind, "SubscriptionType": "Subscription", "EndTime": "2099-01-01T00:00:00Z"}}}}, nil
					case "InquiryPriceRefundInstance":
						inquired = true
						return map[string]interface{}{"Success": mode != "inquiry-rejected", "Code": "Success", "Data": map[string]interface{}{"InstanceId": "target-id", "RefundAmount": json.Number("1.25"), "Currency": "CNY"}}, nil
					case "RefundInstance":
						if request["ProductCode"] != "rds" || request["ProductType"] != "rds_customprepaid_public_cn" || request["ImmediatelyRelease"] != "1" || retry {
							t.Fatal("unsafe subscription refund parameters")
						}
						refunds++
						writes = append(writes, "target-id")
						if mode == "refund-response-lost" {
							return nil, fmt.Errorf("ambiguous refund response")
						}
						return map[string]interface{}{"Success": true, "Code": "Success", "Data": map[string]interface{}{"OrderId": json.Number("123")}}, nil
					}
				}
				if action == "DeleteRCInstances" {
					t.Fatal("subscription reached ordinary delete API")
				}
				d, err := base(product, version, action, query, request, retry)
				if action == "DescribeRCInstanceAttribute" && err == nil {
					d["InstanceChargeType"] = "PrePaid"
					if mode == "fresh-vpc" && inquired {
						d["VpcAttributes"] = map[string]interface{}{"VpcId": "protected"}
					}
				}
				return d, err
			}
			err := sweepVpcRdsCustom(rpc, "region", "target", func(time.Duration) {})
			if mode == "success" {
				if err != nil || refunds != 1 {
					t.Fatalf("refunds=%d err=%v", refunds, err)
				}
			} else if err == nil || (mode != "refund-response-lost" && refunds != 0) || refunds > 1 {
				t.Fatalf("refunds=%d err=%v", refunds, err)
			}
		})
	}
}

func TestUnitVpcSweepServiceOwnersPrivateLinkScopeAndRegion(t *testing.T) {
	var deleted []string
	rpc := func(product, version, action string, query, request map[string]interface{}, read bool) (map[string]interface{}, error) {
		switch action {
		case "ListVpcEndpoints":
			if len(deleted) > 0 {
				return map[string]interface{}{"Endpoints": []interface{}{}}, nil
			}
			return map[string]interface{}{"Endpoints": []interface{}{
				map[string]interface{}{"EndpointId": "target", "VpcId": "vpc-target"},
				map[string]interface{}{"EndpointId": "protected", "VpcId": "vpc-protected"},
			}}, nil
		case "GetVpcEndpointAttribute":
			if request["EndpointId"] != "target" {
				t.Fatalf("described unrelated endpoint: %v", request)
			}
			return map[string]interface{}{"EndpointId": "target", "VpcId": "vpc-target", "ProtectedEnabled": false}, nil
		case "DeleteVpcEndpoint":
			deleted = append(deleted, request["EndpointId"].(string))
			return map[string]interface{}{}, nil
		}
		return nil, fmt.Errorf("unexpected action %s", action)
	}
	node := serviceVpcOwnerSweepers("cn-beijing", "vpc-target", rpc)["alicloud_privatelink_vpc_endpoint"]
	if err := node.F("cn-hangzhou"); err == nil {
		t.Fatal("region mismatch accepted")
	}
	if len(deleted) != 0 {
		t.Fatal("wrong region made mutations")
	}
	if err := node.F("cn-beijing"); err != nil {
		t.Fatal(err)
	}
	if len(deleted) != 1 || deleted[0] != "target" {
		t.Fatalf("deleted %v", deleted)
	}
}

func TestUnitVpcSweepServiceOwnersMissingVpcAndPaginationFailClosed(t *testing.T) {
	for _, malformed := range []bool{true, false} {
		calls, deleted := 0, false
		rpc := func(product, version, action string, query, request map[string]interface{}, read bool) (map[string]interface{}, error) {
			if action == "DeleteVpcEndpoint" {
				deleted = true
				return nil, nil
			}
			calls++
			item := map[string]interface{}{"EndpointId": "endpoint"}
			if !malformed {
				item["VpcId"] = "vpc-target"
			}
			response := map[string]interface{}{"Endpoints": []interface{}{item}}
			if !malformed {
				response["NextToken"] = "repeated"
			}
			return response, nil
		}
		err := serviceVpcOwnerSweepers("cn-beijing", "vpc-target", rpc)["alicloud_privatelink_vpc_endpoint"].F("cn-beijing")
		if err == nil || deleted {
			t.Fatalf("unsafe list accepted: err=%v deleted=%v", err, deleted)
		}
		if calls > 2 {
			t.Fatalf("unbounded pagination: %d", calls)
		}
	}
}

func TestUnitVpcSweepServiceOwnersLindormBillingAndOwnership(t *testing.T) {
	for _, scenario := range []struct {
		pay, detailVpc string
		wantDelete     bool
	}{
		{"POSTPAY", "vpc-target", true}, {"PREPAY", "vpc-target", false},
		{"", "vpc-target", false}, {"POSTPAY", "vpc-protected", false},
	} {
		deleted := false
		rpc := func(product, version, action string, query, request map[string]interface{}, read bool) (map[string]interface{}, error) {
			switch action {
			case "GetLindormInstanceList":
				if deleted {
					return map[string]interface{}{"InstanceList": []interface{}{}}, nil
				}
				return map[string]interface{}{"InstanceList": []interface{}{map[string]interface{}{"InstanceId": "ld-1", "VpcId": "vpc-target"}}}, nil
			case "GetLindormInstance":
				return map[string]interface{}{"InstanceId": "ld-1", "VpcId": scenario.detailVpc, "PayType": scenario.pay}, nil
			case "ReleaseLindormInstance":
				deleted = true
				return map[string]interface{}{}, nil
			}
			return nil, fmt.Errorf("unexpected action %s", action)
		}
		err := serviceVpcOwnerSweepers("cn-beijing", "vpc-target", rpc)["alicloud_lindorm_instance"].F("cn-beijing")
		if deleted != scenario.wantDelete || (err == nil) != scenario.wantDelete {
			t.Fatalf("%+v deleted=%v err=%v", scenario, deleted, err)
		}
	}
}

func TestUnitVpcSweepLindormImmediateReleaseContract(t *testing.T) {
	for _, full := range []bool{false, true} {
		t.Run(fmt.Sprintf("full=%v", full), func(t *testing.T) {
			t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", fmt.Sprint(full))
			releases, lists := 0, 0
			rpc := func(product, version, action string, query, request map[string]interface{}, read bool) (map[string]interface{}, error) {
				switch action {
				case "GetLindormInstanceList":
					lists++
					items := []interface{}{map[string]interface{}{"InstanceId": "protected", "VpcId": "vpc-protected"}}
					if releases == 0 {
						items = append(items, map[string]interface{}{"InstanceId": "ld-1", "VpcId": "vpc-target"})
					}
					return map[string]interface{}{"InstanceList": items}, nil
				case "GetLindormInstance":
					return map[string]interface{}{"InstanceId": "ld-1", "VpcId": "vpc-target", "PayType": "POSTPAY"}, nil
				case "ReleaseLindormInstance":
					releases++
					immediately, present := request["Immediately"]
					wantFields := 1
					if full {
						wantFields = 2
					}
					if request["InstanceId"] != "ld-1" || present != full || (full && immediately != true) || len(request) != wantFields || read || query != nil {
						t.Fatalf("unexpected release request: %v read=%v query=%v", request, read, query)
					}
					return map[string]interface{}{}, nil
				}
				return nil, fmt.Errorf("unexpected action %s", action)
			}
			if err := serviceVpcOwnerSweepers("cn-beijing", "vpc-target", rpc)["alicloud_lindorm_instance"].F("cn-beijing"); err != nil || releases != 1 || lists != 2 {
				t.Fatalf("err=%v releases=%d lists=%d", err, releases, lists)
			}
		})
	}
}

func TestUnitVpcSweepLindormFullScopeAndReleaseErrors(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	for _, scenario := range []struct {
		name, vpc, pay, id string
		releaseError       bool
	}{
		{"protected", "vpc-protected", "POSTPAY", "ld-1", false},
		{"unknown-scope", "", "POSTPAY", "ld-1", false},
		{"subscription", "vpc-target", "PREPAY", "ld-1", false},
		{"unknown-payment", "vpc-target", "", "ld-1", false},
		{"identity", "vpc-target", "POSTPAY", "other", false},
		{"incorrect-state", "vpc-target", "POSTPAY", "ld-1", true},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			lists, releases := 0, 0
			rpc := func(product, version, action string, query, request map[string]interface{}, read bool) (map[string]interface{}, error) {
				switch action {
				case "GetLindormInstanceList":
					lists++
					return map[string]interface{}{"InstanceList": []interface{}{map[string]interface{}{"InstanceId": "ld-1", "VpcId": "vpc-target"}}}, nil
				case "GetLindormInstance":
					return map[string]interface{}{"InstanceId": scenario.id, "VpcId": scenario.vpc, "PayType": scenario.pay}, nil
				case "ReleaseLindormInstance":
					releases++
					if !scenario.releaseError {
						t.Fatal("released unverifiable or non-POSTPAY instance")
					}
					if request["Immediately"] != true {
						t.Fatal("full cleanup omitted immediate release")
					}
					return nil, &tea.SDKError{Code: tea.String("IncorrectDBState"), Message: tea.String("instance is creating"), StatusCode: tea.Int(400)}
				}
				return nil, fmt.Errorf("unexpected action %s", action)
			}
			err := serviceVpcOwnerSweepers("cn-beijing", "vpc-target", rpc)["alicloud_lindorm_instance"].F("cn-beijing")
			wantReleases := 0
			if scenario.releaseError {
				wantReleases = 1
			}
			if err == nil || releases != wantReleases || lists != 1 {
				t.Fatalf("err=%v releases=%d lists=%d", err, releases, lists)
			}
			if scenario.releaseError && !strings.Contains(err.Error(), "IncorrectDBState") {
				t.Fatalf("lost state error: %v", err)
			}
		})
	}
}

func TestUnitVpcSweepServiceOwnersKmsSharedBindingsAndBilling(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "false")
	for _, scenario := range []struct {
		charge, bound string
		wantDelete    bool
	}{
		{"POSTPAY", "vpc-target", true}, {"PREPAY", "vpc-target", false},
		{"POSTPAY", "vpc-protected", false}, {"UNKNOWN", "vpc-target", false},
	} {
		deleted := false
		rpc := func(product, version, action string, query, request map[string]interface{}, read bool) (map[string]interface{}, error) {
			switch action {
			case "ListKmsInstances":
				if deleted {
					return map[string]interface{}{"KmsInstances": map[string]interface{}{"KmsInstance": []interface{}{}}}, nil
				}
				return map[string]interface{}{"KmsInstances": map[string]interface{}{"KmsInstance": []interface{}{map[string]interface{}{"KmsInstanceId": "kms-1"}}}}, nil
			case "GetKmsInstance":
				return map[string]interface{}{"KmsInstance": map[string]interface{}{"InstanceId": "kms-1", "VpcId": "vpc-target", "ChargeType": scenario.charge, "BindVpcs": map[string]interface{}{"BindVpc": []interface{}{map[string]interface{}{"VpcId": scenario.bound}}}}}, nil
			case "ReleaseKmsInstance":
				deleted = true
				return map[string]interface{}{}, nil
			}
			return nil, fmt.Errorf("unexpected action %s", action)
		}
		err := serviceVpcOwnerSweepers("cn-beijing", "vpc-target", rpc)["alicloud_kms_instance"].F("cn-beijing")
		if deleted != scenario.wantDelete || (err == nil) != scenario.wantDelete {
			t.Fatalf("%+v deleted=%v err=%v", scenario, deleted, err)
		}
		if !scenario.wantDelete && err != nil && !strings.Contains(err.Error(), "not released") {
			t.Fatalf("missing explicit blocker: %v", err)
		}
	}
}

func TestUnitVpcSweepServiceOwnersBoundedAsyncReleaseWait(t *testing.T) {
	calls, pauses := 0, 0
	list := func() ([]map[string]interface{}, error) {
		calls++
		if calls < 3 {
			return []map[string]interface{}{{"Id": "released", "VpcId": "vpc-target"}, {"Id": "other", "VpcId": "vpc-protected"}}, nil
		}
		return []map[string]interface{}{{"Id": "other", "VpcId": "vpc-protected"}}, nil
	}
	err := serviceVpcOwnerWaitGone(list, "Id", "released", "vpc-target", func(d time.Duration) {
		pauses++
		if d != 5*time.Second {
			t.Fatalf("unexpected pause %s", d)
		}
	})
	if err != nil || calls != 3 || pauses != 2 {
		t.Fatalf("calls=%d pauses=%d err=%v", calls, pauses, err)
	}
}

func TestUnitVpcSweepServiceOwnersWaitNeverTreatsUnknownAsGone(t *testing.T) {
	for _, list := range []func() ([]map[string]interface{}, error){
		func() ([]map[string]interface{}, error) { return nil, fmt.Errorf("permission denied") },
		func() ([]map[string]interface{}, error) {
			return []map[string]interface{}{{"VpcId": "vpc-target"}}, nil
		},
		func() ([]map[string]interface{}, error) { return []map[string]interface{}{{"Id": "released"}}, nil },
		func() ([]map[string]interface{}, error) {
			return []map[string]interface{}{{"Id": "released", "VpcId": "vpc-protected"}}, nil
		},
	} {
		if err := serviceVpcOwnerWaitGone(list, "Id", "released", "vpc-target", func(time.Duration) {}); err == nil {
			t.Fatal("unknown list accepted as gone")
		}
	}
	calls, pauses := 0, 0
	list := func() ([]map[string]interface{}, error) {
		calls++
		return []map[string]interface{}{{"Id": "released", "VpcId": "vpc-target"}}, nil
	}
	if err := serviceVpcOwnerWaitGone(list, "Id", "released", "vpc-target", func(time.Duration) { pauses++ }); err == nil || calls != 120 || pauses != 119 {
		t.Fatalf("unbounded/missing wait: calls=%d pauses=%d err=%v", calls, pauses, err)
	}
}

func TestUnitVpcSweepServiceOwnersDataWorksSharedGroupProtectsOtherVpc(t *testing.T) {
	deleted := false
	rpc := func(product, version, action string, query, request map[string]interface{}, read bool) (map[string]interface{}, error) {
		switch action {
		case "ListResourceGroups":
			return map[string]interface{}{"PagingInfo": map[string]interface{}{"ResourceGroupList": []interface{}{map[string]interface{}{"Id": float64(42), "DefaultVpcId": "vpc-target"}}}}, nil
		case "ListNetworks":
			if request["ResourceGroupId"] != "42" {
				t.Fatalf("wrong group scope: %v", request)
			}
			rows := []interface{}{map[string]interface{}{"Id": "protected", "VpcId": "vpc-protected"}}
			if !deleted {
				rows = append(rows, map[string]interface{}{"Id": "target", "VpcId": "vpc-target"})
			}
			return map[string]interface{}{"PagingInfo": map[string]interface{}{"NetworkList": rows}}, nil
		case "DeleteNetwork":
			if request["Id"] != "target" || request["RegionId"] != "cn-beijing" {
				t.Fatalf("unscoped deletion: %v", request)
			}
			deleted = true
			return map[string]interface{}{}, nil
		default:
			t.Fatalf("unexpected action, especially shared group deletion: %s", action)
		}
		return nil, nil
	}
	err := serviceVpcOwnerSweepers("cn-beijing", "vpc-target", rpc)["alicloud_data_works_network"].F("cn-beijing")
	if err != nil || !deleted {
		t.Fatalf("deleted=%v err=%v", deleted, err)
	}
}

func TestUnitVpcSweepServiceOwnersDataWorksUnknownVpcNeverDeletes(t *testing.T) {
	rpc := func(product, version, action string, query, request map[string]interface{}, read bool) (map[string]interface{}, error) {
		switch action {
		case "ListResourceGroups":
			return map[string]interface{}{"PagingInfo": map[string]interface{}{"ResourceGroupList": []interface{}{map[string]interface{}{"Id": "group"}}}}, nil
		case "ListNetworks":
			return map[string]interface{}{"PagingInfo": map[string]interface{}{"NetworkList": []interface{}{map[string]interface{}{"Id": "unknown"}}}}, nil
		default:
			t.Fatalf("unexpected write: %s", action)
		}
		return nil, nil
	}
	if err := serviceVpcOwnerSweepers("cn-beijing", "vpc-target", rpc)["alicloud_data_works_network"].F("cn-beijing"); err == nil {
		t.Fatal("unknown network ownership accepted")
	}
}

func TestUnitVpcSweepServiceOwnersRepeatedNumericPageRejected(t *testing.T) {
	calls := 0
	rows := make([]interface{}, PageSizeLarge)
	for i := range rows {
		rows[i] = map[string]interface{}{"InstanceId": fmt.Sprintf("instance-%d", i), "VpcId": "vpc-target"}
	}
	_, err := serviceVpcOwnerList(func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
		calls++
		return map[string]interface{}{"InstanceList": rows}, nil
	}, "hitsdb", "2020-06-15", "GetLindormInstanceList", "$.InstanceList", "", false)
	if err == nil || calls != 2 {
		t.Fatalf("repeated numeric page accepted calls=%d error=%v", calls, err)
	}
}

func TestUnitVpcSweepDataWorksTransportUsesOfficialHttpMethods(t *testing.T) {
	gets, posts := 0, 0
	get := func(product, version, action string, query, body map[string]interface{}) (map[string]interface{}, error) {
		gets++
		if product != "dataworks-public" || version != "2024-05-18" || body != nil || query["RegionId"] != "cn-beijing" || query["ResourceGroupId"] != "group" {
			t.Fatalf("invalid GET transport: %s %v %v", action, query, body)
		}
		return map[string]interface{}{"method": "GET"}, nil
	}
	post := func(product, version, action string, query, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
		posts++
		if body["Id"] != "network" || body["RegionId"] != "cn-beijing" || !autoRetry {
			t.Fatalf("invalid POST transport: %s %v", action, body)
		}
		return map[string]interface{}{"method": "POST"}, nil
	}
	rpc := dataWorksVpcCleanupRPC(get, post)
	for _, action := range []string{"ListResourceGroups", "ListNetworks", "GetResourceGroup", "GetNetwork"} {
		query := map[string]interface{}{"RegionId": "cn-beijing"}
		body := map[string]interface{}{"ResourceGroupId": "group"}
		response, err := rpc("dataworks-public", "2024-05-18", action, query, body, true)
		if err != nil || response["method"] != "GET" {
			t.Fatalf("%s: %v", action, err)
		}
		if len(query) != 1 || len(body) != 1 {
			t.Fatal("caller parameter maps mutated")
		}
	}
	for _, action := range []string{"DeleteNetwork", "DeleteResourceGroup"} {
		response, err := rpc("dataworks-public", "2024-05-18", action, nil, map[string]interface{}{"Id": "network", "RegionId": "cn-beijing"}, true)
		if err != nil || response["method"] != "POST" {
			t.Fatalf("%s: %v", action, err)
		}
	}
	if gets != 4 || posts != 2 {
		t.Fatalf("gets=%d posts=%d", gets, posts)
	}
}

func TestUnitVpcSweepDataWorksTransportRejectsUnknownAndAmbiguousCalls(t *testing.T) {
	get := func(string, string, string, map[string]interface{}, map[string]interface{}) (map[string]interface{}, error) {
		t.Fatal("invalid GET dispatched")
		return nil, nil
	}
	post := func(string, string, string, map[string]interface{}, map[string]interface{}, bool) (map[string]interface{}, error) {
		t.Fatal("invalid POST dispatched")
		return nil, nil
	}
	rpc := dataWorksVpcCleanupRPC(get, post)
	for _, input := range []struct{ product, version, action string }{
		{"dataworks-public", "2024-05-18", "DeleteEverything"},
		{"dataworks-public", "2020-09-18", "ListNetworks"},
		{"other", "2024-05-18", "ListNetworks"},
	} {
		if _, err := rpc(input.product, input.version, input.action, nil, nil, true); err == nil {
			t.Fatalf("invalid input accepted: %+v", input)
		}
	}
	if _, err := rpc("dataworks-public", "2024-05-18", "ListNetworks", map[string]interface{}{"RegionId": "one"}, map[string]interface{}{"RegionId": "two"}, true); err == nil {
		t.Fatal("ambiguous GET scope accepted")
	}
}

func TestUnitVpcSweepKmsBindingsOfficialNilSliceAndMalformedShapes(t *testing.T) {
	for _, value := range []interface{}{nil, []interface{}{}, []interface{}{map[string]interface{}{"VpcId": "vpc-target"}}} {
		rows, err := kmsVpcOwnerBindings(map[string]interface{}{"BindVpcs": map[string]interface{}{"BindVpc": value}})
		if err != nil {
			t.Fatalf("official shape %T rejected: %v", value, err)
		}
		if value == nil && len(rows) != 0 {
			t.Fatal("nil slice not empty")
		}
	}
	for _, detail := range []map[string]interface{}{
		{}, {"BindVpcs": nil}, {"BindVpcs": map[string]interface{}{}},
		{"BindVpcs": []interface{}{}}, {"BindVpcs": map[string]interface{}{"BindVpc": "unknown"}},
		{"BindVpcs": map[string]interface{}{"BindVpc": []interface{}{nil}}},
	} {
		if _, err := kmsVpcOwnerBindings(detail); err == nil {
			t.Fatalf("unknown shape accepted: %v", detail)
		}
	}
}

func TestUnitVpcSweepKmsRealNullBindingsStillRequireVerifiedBilling(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "false")
	rpc := func(product, version, action string, query, request map[string]interface{}, read bool) (map[string]interface{}, error) {
		switch action {
		case "ListKmsInstances":
			return map[string]interface{}{"KmsInstances": map[string]interface{}{"KmsInstance": []interface{}{map[string]interface{}{"KmsInstanceId": "kms"}}}}, nil
		case "GetKmsInstance":
			return map[string]interface{}{"KmsInstance": map[string]interface{}{"InstanceId": "kms", "VpcId": "vpc-target", "ChargeType": "PREPAY", "BindVpcs": map[string]interface{}{"BindVpc": nil}}}, nil
		default:
			t.Fatalf("subscription mutation: %s", action)
		}
		return nil, nil
	}
	err := serviceVpcOwnerSweepers("cn-beijing", "vpc-target", rpc)["alicloud_kms_instance"].F("cn-beijing")
	if err == nil || !strings.Contains(err.Error(), "explicit SWEEP_ALL") {
		t.Fatalf("wrong real-schema blocker: %v", err)
	}
}

func TestUnitVpcSweepPrivateLinkManagedProtectionAndOwnershipPrecedeDelete(t *testing.T) {
	for _, scenario := range []string{"protected", "unknown", "other-vpc", "other-id", "read-error", "delete-error", "success"} {
		var actions []string
		deleted := false
		rpc := func(product, version, action string, query, request map[string]interface{}, read bool) (map[string]interface{}, error) {
			actions = append(actions, action)
			switch action {
			case "ListVpcEndpoints":
				if deleted {
					return map[string]interface{}{"Endpoints": []interface{}{}}, nil
				}
				return map[string]interface{}{"Endpoints": []interface{}{map[string]interface{}{"EndpointId": "ep", "VpcId": "vpc-target"}}}, nil
			case "GetVpcEndpointAttribute":
				if scenario == "read-error" {
					return nil, fmt.Errorf("permission denied")
				}
				detail := map[string]interface{}{"EndpointId": "ep", "VpcId": "vpc-target", "ProtectedEnabled": false}
				switch scenario {
				case "protected":
					detail["ProtectedEnabled"] = true
				case "unknown":
					delete(detail, "ProtectedEnabled")
				case "other-vpc":
					detail["VpcId"] = "vpc-protected"
				case "other-id":
					detail["EndpointId"] = "other"
				}
				return detail, nil
			case "DeleteVpcEndpoint":
				if scenario != "success" && scenario != "delete-error" {
					t.Fatalf("%s reached deletion", scenario)
				}
				if scenario == "delete-error" {
					return nil, fmt.Errorf("permission denied")
				}
				deleted = true
				return map[string]interface{}{}, nil
			default:
				t.Fatalf("unverified protection update/action: %s", action)
			}
			return nil, nil
		}
		err := serviceVpcOwnerSweepers("cn-beijing", "vpc-target", rpc)["alicloud_privatelink_vpc_endpoint"].F("cn-beijing")
		if scenario == "success" {
			if err != nil || strings.Join(actions, ",") != "ListVpcEndpoints,GetVpcEndpointAttribute,DeleteVpcEndpoint,ListVpcEndpoints" {
				t.Fatalf("unsafe ordering: %v err=%v", actions, err)
			}
		} else if err == nil {
			t.Fatalf("%s falsely completed", scenario)
		}
	}
}

func TestUnitVpcSweepDataWorksSDKNumberIdentity(t *testing.T) {
	value, err := vpcSweepUtilV2.ReadAsJSON(strings.NewReader(`{"Id":9223372036854775807}`))
	if err != nil {
		t.Fatal(err)
	}
	id := value.(map[string]interface{})["Id"]
	if _, ok := id.(json.Number); !ok {
		t.Fatalf("production SDK numeric shape changed: %T", id)
	}
	actual, err := serviceVpcDataWorksID(id)
	if err != nil || actual != "9223372036854775807" {
		t.Fatalf("exact int64 identity lost actual=%s err=%v", actual, err)
	}
	for _, legacy := range []string{"network-id", "42"} {
		if actual, err := serviceVpcDataWorksID(legacy); err != nil || actual != legacy {
			t.Fatalf("legacy string changed %q %v", actual, err)
		}
	}
}

func TestUnitVpcSweepDataWorksSDKNumberRejectsNonCanonicalIdentity(t *testing.T) {
	for _, value := range []string{"", "0", "-0", "-1", "+1", "01", "1.0", "1e3", " 1", "1 ", "9223372036854775808", "NaN"} {
		if actual, err := serviceVpcDataWorksID(json.Number(value)); err == nil || actual != "" {
			t.Fatalf("unsafe numeric identity accepted %q as %q", value, actual)
		}
	}
	for _, value := range []string{"1", "9007199254740993", "9223372036854775807"} {
		if actual, err := serviceVpcDataWorksID(json.Number(value)); err != nil || actual != value {
			t.Fatalf("valid exact identity rejected %q: %q %v", value, actual, err)
		}
	}
}

func vpcKmsBillingFixtureRPC(t *testing.T, scenario string, tokens *[]string) vpcSweepRPC {
	t.Helper()
	return func(product, version, action string, query, request map[string]interface{}, auto bool) (map[string]interface{}, error) {
		var response map[string]interface{}
		switch action {
		case "ListKmsInstances":
			if len(*tokens) > 0 {
				if scenario == "wait-read-error" {
					return nil, fmt.Errorf("permission denied during release confirmation")
				}
				return map[string]interface{}{"KmsInstances": map[string]interface{}{"KmsInstance": []interface{}{}}}, nil
			}
			return map[string]interface{}{"KmsInstances": map[string]interface{}{"KmsInstance": []interface{}{map[string]interface{}{"KmsInstanceId": "kms-target"}}}}, nil
		case "GetKmsInstance":
			detail := map[string]interface{}{"InstanceId": "kms-target", "VpcId": "vpc-target", "ChargeType": "PREPAY", "BindVpcs": map[string]interface{}{"BindVpc": []interface{}{}}}
			switch scenario {
			case "changed-vpc":
				detail["VpcId"] = "vpc-protected"
			case "shared-binding":
				detail["BindVpcs"] = map[string]interface{}{"BindVpc": []interface{}{map[string]interface{}{"VpcId": "vpc-protected"}}}
			case "changed-charge":
				detail["ChargeType"] = "POSTPAY"
			case "deletion-protected":
				detail["DeletionProtection"] = true
			}
			return map[string]interface{}{"KmsInstance": detail}, nil
		case "QueryAvailableInstances":
			if product != "BssOpenApi" || version != "2017-12-14" || request["InstanceIDs"] != "kms-target" || request["PageNum"] != 1 || request["PageSize"] != 100 || len(request) != 3 {
				t.Fatalf("unscoped billing inventory: %s %s %v", product, version, request)
			}
			row := map[string]interface{}{"InstanceID": "kms-target", "ProductCode": "kms", "ProductType": "kms_native_type", "SubscriptionType": "Subscription", "Region": "cn-beijing", "EndTime": "2099-01-01T00:00:00Z"}
			switch scenario {
			case "wrong-id":
				row["InstanceID"] = "kms-other"
			case "wrong-region":
				row["Region"] = "cn-hangzhou"
			case "wrong-product":
				row["ProductCode"] = "ecs"
			case "wrong-subscription":
				row["SubscriptionType"] = "PayAsYouGo"
			case "period-unknown":
				delete(row, "EndTime")
			case "renewed":
				row["EndTime"] = "2100-01-01T00:00:00Z"
			}
			rows := []interface{}{row}
			total := json.Number("1")
			if scenario == "duplicate" {
				rows = append(rows, row)
			}
			if scenario == "incomplete" {
				total = json.Number("2")
			}
			response = map[string]interface{}{"Success": scenario != "inventory-failed", "Code": "Success", "Data": map[string]interface{}{"TotalCount": total, "PageNum": json.Number("1"), "InstanceList": rows}}
		case "InquiryPriceRefundInstance":
			if request["InstanceId"] != "kms-target" || request["ProductCode"] != "kms" || request["ProductType"] != "kms_native_type" || len(request) != 3 {
				t.Fatalf("refund inquiry did not use exact billing identity: %v", request)
			}
			price := map[string]interface{}{"InstanceId": "kms-target", "RefundAmount": json.Number("0.01"), "Currency": "CNY"}
			switch scenario {
			case "inquiry-other-id":
				price["InstanceId"] = "kms-other"
			case "negative-amount":
				price["RefundAmount"] = json.Number("-1")
			case "invalid-amount":
				price["RefundAmount"] = "NaN"
			case "fraction-amount":
				price["RefundAmount"] = "1/2"
			case "invalid-currency":
				price["Currency"] = "12!"
			case "missing-currency":
				delete(price, "Currency")
			}
			response = map[string]interface{}{"Success": scenario != "inquiry-failed", "Code": "Success", "Data": price}
		case "RefundInstance":
			if product != "BssOpenApi" || version != "2017-12-14" || auto || request["InstanceId"] != "kms-target" || request["ProductCode"] != "kms" || request["ProductType"] != "kms_native_type" || request["ImmediatelyRelease"] != "1" || len(request) != 5 {
				t.Fatalf("unsafe refund request or blind automatic retry: %s %s %v auto=%v", product, version, request, auto)
			}
			token, ok := request["ClientToken"].(string)
			if !ok || !strings.HasPrefix(token, "vpc-kms-") || len(token) > 64 {
				t.Fatalf("missing bounded idempotency token: %v", request)
			}
			*tokens = append(*tokens, token)
			if scenario == "refund-network-error" {
				return nil, fmt.Errorf("response lost after refund submission")
			}
			order := map[string]interface{}{"OrderId": json.Number("123456789012345678")}
			if scenario == "missing-order" {
				delete(order, "OrderId")
			} else if scenario == "invalid-order" {
				order["OrderId"] = "not-an-order"
			}
			response = map[string]interface{}{"Success": scenario != "refund-failed", "Code": "Success", "Data": order}
		default:
			t.Fatalf("unexpected KMS subscription operation %s", action)
		}
		return response, nil
	}
}

func TestUnitVpcSweepKmsSubscriptionNativeRefundAndConfirmation(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	var tokens []string
	rpc := vpcKmsBillingFixtureRPC(t, "success", &tokens)
	if err := serviceVpcOwnerSweepers("cn-beijing", "vpc-target", rpc)["alicloud_kms_instance"].F("cn-beijing"); err != nil || len(tokens) != 1 {
		t.Fatalf("refund/confirmation failed: %v writes=%d", err, len(tokens))
	}
}

func TestUnitVpcSweepKmsSubscriptionBoundariesNeverRefund(t *testing.T) {
	for _, scenario := range []string{"disabled", "wrong-id", "wrong-region", "wrong-product", "wrong-subscription", "period-unknown", "duplicate", "incomplete", "inventory-failed", "inquiry-other-id", "inquiry-failed", "negative-amount", "invalid-amount", "fraction-amount", "invalid-currency", "missing-currency", "changed-vpc", "shared-binding", "changed-charge", "deletion-protected"} {
		t.Run(scenario, func(t *testing.T) {
			t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", fmt.Sprint(scenario != "disabled"))
			var tokens []string
			err := refundVpcKmsSubscription(vpcKmsBillingFixtureRPC(t, scenario, &tokens), "cn-beijing", "vpc-target", "kms-target")
			if err == nil || len(tokens) != 0 {
				t.Fatalf("boundary accepted: %v writes=%d", err, len(tokens))
			}
		})
	}
}

func TestUnitVpcSweepKmsSubscriptionRefundErrorsAndStableIdempotency(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	for _, scenario := range []string{"refund-network-error", "missing-order", "invalid-order", "refund-failed"} {
		var tokens []string
		rpc := vpcKmsBillingFixtureRPC(t, scenario, &tokens)
		for i := 0; i < 2; i++ {
			if err := refundVpcKmsSubscription(rpc, "cn-beijing", "vpc-target", "kms-target"); err == nil {
				t.Fatalf("ambiguous/rejected refund counted as success: %s", scenario)
			}
		}
		if len(tokens) != 2 || tokens[0] != tokens[1] {
			t.Fatalf("retry did not retain the same native idempotency token: %v", tokens)
		}
	}
	var lifecycleTokens []string
	for _, scenario := range []string{"success", "renewed"} {
		if err := refundVpcKmsSubscription(vpcKmsBillingFixtureRPC(t, scenario, &lifecycleTokens), "cn-beijing", "vpc-target", "kms-target"); err != nil {
			t.Fatal(err)
		}
	}
	if len(lifecycleTokens) != 2 || lifecycleTokens[0] == lifecycleTokens[1] {
		t.Fatal("a renewed subscription reused the prior billing period's idempotency key")
	}

}

func TestUnitVpcSweepKmsSubscriptionOrderIsNotReleaseProof(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	var tokens []string
	err := serviceVpcOwnerSweepers("cn-beijing", "vpc-target", vpcKmsBillingFixtureRPC(t, "wait-read-error", &tokens))["alicloud_kms_instance"].F("cn-beijing")
	if err == nil || len(tokens) != 1 || !strings.Contains(err.Error(), "permission denied") {
		t.Fatalf("refund order counted as deletion despite failed readback: %v writes=%d", err, len(tokens))
	}
}

func vpcKvstoreBillingFixtureRPC(t *testing.T, scenario string, tokens *[]string) vpcSweepRPC {
	t.Helper()
	return func(product, version, action string, query, request map[string]interface{}, auto bool) (map[string]interface{}, error) {
		switch action {
		case "DescribeInstances":
			rows := []interface{}{}
			if len(*tokens) == 0 {
				row := map[string]interface{}{"InstanceId": "redis-target", "VpcId": "vpc-target", "GlobalInstanceId": "", "ChargeType": "PrePaid"}
				if scenario == "global-shared" {
					row["GlobalInstanceId"] = "global-shared"
				} else if scenario == "global-unknown" {
					delete(row, "GlobalInstanceId")
				}
				rows = append(rows, row)
			}
			return map[string]interface{}{"Instances": map[string]interface{}{"KVStoreInstance": rows}}, nil
		case "DescribeInstanceAttribute":
			detail := map[string]interface{}{"InstanceId": "redis-target", "VpcId": "vpc-target", "RegionId": "cn-hangzhou", "ChargeType": "PrePaid", "InstanceReleaseProtection": false}
			switch scenario {
			case "detail-protected-vpc":
				detail["VpcId"] = "vpc-protected"
			case "wrong-detail-id":
				detail["InstanceId"] = "redis-other"
			case "release-protected":
				detail["InstanceReleaseProtection"] = true
			case "protection-unknown":
				delete(detail, "InstanceReleaseProtection")
			}
			return map[string]interface{}{"Instances": map[string]interface{}{"DBInstanceAttribute": []interface{}{detail}}}, nil
		case "DescribeDBInstanceNetInfo":
			private := map[string]interface{}{"IPType": "Private", "VPCId": "vpc-target"}
			switch scenario {
			case "endpoint-protected-vpc":
				private["VPCId"] = "vpc-protected"
			case "endpoint-unknown-vpc":
				delete(private, "VPCId")
			case "endpoint-unknown-type":
				delete(private, "IPType")
			}
			return map[string]interface{}{"NetInfoItems": map[string]interface{}{"InstanceNetInfo": []interface{}{private, map[string]interface{}{"IPType": "Public"}}}}, nil
		case "QueryAvailableInstances":
			if request["InstanceIDs"] != "redis-target" || request["PageNum"] != 1 || request["PageSize"] != 100 {
				t.Fatalf("unscoped Redis billing query: %v", request)
			}
			row := map[string]interface{}{"InstanceID": "redis-target", "ProductCode": "redisa", "SubscriptionType": "Subscription", "Region": "cn-hangzhou", "EndTime": "2099-01-01T00:00:00Z"}
			if scenario == "wrong-billing-product" {
				row["ProductCode"] = "kms"
			}
			return map[string]interface{}{"Success": true, "Code": "Success", "Data": map[string]interface{}{"TotalCount": json.Number("1"), "PageNum": json.Number("1"), "InstanceList": []interface{}{row}}}, nil
		case "InquiryPriceRefundInstance":
			if request["ProductCode"] != "redisa" || request["ProductType"] != "" || request["InstanceId"] != "redis-target" {
				t.Fatalf("Redis native billing code/default empty product type changed: %v", request)
			}
			return map[string]interface{}{"Success": true, "Code": "Success", "Data": map[string]interface{}{"InstanceId": "redis-target", "RefundAmount": json.Number("0"), "Currency": "CNY"}}, nil
		case "RefundInstance":
			token, ok := request["ClientToken"].(string)
			if product != "BssOpenApi" || version != "2017-12-14" || auto || request["InstanceId"] != "redis-target" || request["ProductCode"] != "redisa" || request["ProductType"] != "" || request["ImmediatelyRelease"] != "1" || !ok || !strings.HasPrefix(token, "vpc-redisa-") || len(token) > 64 {
				t.Fatalf("unsafe Redis refund: %v auto=%v", request, auto)
			}
			*tokens = append(*tokens, token)
			return map[string]interface{}{"Success": true, "Code": "Success", "Data": map[string]interface{}{"OrderId": json.Number("123456789012345678")}}, nil
		default:
			t.Fatalf("unexpected Redis cleanup operation %s", action)
			return nil, nil
		}
	}
}

func TestUnitVpcSweepKvstoreSubscriptionNativeRefundAndExactRelease(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	var tokens []string
	for _, owner := range allVpcOwners {
		if owner.product != "R-kvstore" {
			continue
		}
		if err := sweepVpcOwner(vpcKvstoreBillingFixtureRPC(t, "success", &tokens), "cn-hangzhou", "vpc-target", owner); err != nil || len(tokens) != 1 {
			t.Fatalf("Redis refund/release failed: %v writes=%d", err, len(tokens))
		}
		return
	}
	t.Fatal("Redis owner not wired")
}

func TestUnitVpcSweepKvstoreSubscriptionSharedAndUnknownNeverRefund(t *testing.T) {
	for _, scenario := range []string{"disabled", "global-shared", "global-unknown", "detail-protected-vpc", "wrong-detail-id", "release-protected", "protection-unknown", "endpoint-protected-vpc", "endpoint-unknown-vpc", "endpoint-unknown-type", "wrong-billing-product"} {
		t.Run(scenario, func(t *testing.T) {
			t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", fmt.Sprint(scenario != "disabled"))
			for _, owner := range allVpcOwners {
				if owner.product != "R-kvstore" {
					continue
				}
				var tokens []string
				err := refundVpcKvstoreSubscription(vpcKvstoreBillingFixtureRPC(t, scenario, &tokens), "cn-hangzhou", "vpc-target", "redis-target", owner)
				if err == nil || len(tokens) != 0 {
					t.Fatalf("Redis boundary accepted: %v writes=%d", err, len(tokens))
				}
				return
			}
			t.Fatal("Redis owner not wired")
		})
	}
}

type vpcStarrocksAPI struct {
	describe func(string, int) (map[string]interface{}, error)
	release  func(string) (map[string]interface{}, error)
	vswitch  func(string) (map[string]interface{}, error)
}

type starrocksVpcRoaPost func(string, string, string, string, map[string]*string, map[string]*string, interface{}, bool) (map[string]interface{}, error)

func nativeVpcStarrocksAPI(client *connectivity.AliyunClient) vpcStarrocksAPI {
	return starrocksVpcAPIWithTransport(client.RegionId, client.RoaPostWithApiName, throttledVpcSweepRPC(client.RpcPost))
}

// Both actions are POST ROA operations with query parameters, not RPC actions.
// The existing client resolves starrocks via Location service emrstarrocks;
// no endpoint or credentials are invented or copied into this adapter.
func starrocksVpcAPIWithTransport(region string, post starrocksVpcRoaPost, vpcRPC vpcSweepRPC) vpcStarrocksAPI {
	rpc := throttledVpcSweepRPC(func(product, version, action string, query, body map[string]interface{}, auto bool) (map[string]interface{}, error) {
		if product != "starrocks" || version != "2022-10-19" || len(body) != 0 {
			return nil, fmt.Errorf("unsupported StarRocks cleanup transport")
		}
		path := ""
		switch action {
		case "DescribeInstances":
			path = "/webapi/starrocks/describeInstances"
		case "ReleaseInstance":
			path = "/webapi/cluster/release"
		default:
			return nil, fmt.Errorf("unsupported StarRocks cleanup action %s", action)
		}
		parameters := map[string]*string{}
		for key, value := range query {
			text, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("invalid StarRocks query parameter %s", key)
			}
			parameters[key] = tea.String(text)
		}
		return post(product, version, action, path, parameters, nil, nil, auto)
	})
	return vpcStarrocksAPI{
		describe: func(id string, page int) (map[string]interface{}, error) {
			query := map[string]interface{}{"RegionId": region, "PageNumber": strconv.Itoa(page), "PageSize": "100"}
			if id != "" {
				query["InstanceId"] = id
			}
			return rpc("starrocks", "2022-10-19", "DescribeInstances", query, nil, true)
		},
		release: func(id string) (map[string]interface{}, error) {
			return rpc("starrocks", "2022-10-19", "ReleaseInstance", map[string]interface{}{"InstanceId": id}, nil, true)
		},
		vswitch: func(id string) (map[string]interface{}, error) {
			response, err := vpcRPC("Vpc", "2016-04-28", "DescribeVSwitches", nil, map[string]interface{}{"RegionId": region, "VSwitchId": id, "PageSize": 50, "PageNumber": 1}, true)
			if err != nil {
				return nil, err
			}
			rows, err := serviceVpcOwnerRows("$.VSwitches.VSwitch", response)
			if err != nil {
				return nil, err
			}
			if len(rows) != 1 || rows[0]["VSwitchId"] != id {
				return nil, fmt.Errorf("StarRocks VSwitch %s missing/ambiguous identity", id)
			}
			return rows[0], nil
		},
	}
}

func vpcStarrocksSweeper(region, vpcID string, api vpcStarrocksAPI) *resource.Sweeper {
	return &resource.Sweeper{Name: "alicloud_emr_serverless_starrocks_instance", F: func(actual string) error {
		if actual != region || region == "" || vpcID == "" {
			return fmt.Errorf("StarRocks requires selected VPC and matching region")
		}
		return sweepVpcStarrocks(region, vpcID, api, time.Sleep, 120)
	}}
}

func starrocksVpcRows(response map[string]interface{}) ([]map[string]interface{}, int, error) {
	if response["Success"] != true {
		return nil, 0, fmt.Errorf("StarRocks response did not confirm success (code %v)", response["ErrCode"])
	}
	total := -1
	switch value := response["Total"].(type) {
	case int:
		total = value
	case int64:
		if value >= 0 && value <= 1000000 {
			total = int(value)
		}
	case float64:
		if value >= 0 && value <= 1000000 && value == float64(int(value)) {
			total = int(value)
		}
	case json.Number:
		// tea-roa/tea-utils decode raw API numbers with UseNumber. Accept
		// only the canonical integer form established by the API schema.
		text := value.String()
		valid := text != "" && len(text) <= 7 && (len(text) == 1 || text[0] != '0')
		for _, digit := range text {
			if digit < '0' || digit > '9' {
				valid = false
				break
			}
		}
		if valid {
			parsed, err := strconv.ParseInt(text, 10, 64)
			if err == nil && parsed <= 1000000 {
				total = int(parsed)
			}
		}
	}
	if total < 0 || total > 1000000 {
		return nil, 0, fmt.Errorf("StarRocks missing/invalid total")
	}
	rows, err := serviceVpcOwnerRows("$.Data", response)
	if err != nil {
		return nil, 0, err
	}
	if len(rows) > 100 || len(rows) > total {
		return nil, 0, fmt.Errorf("StarRocks inconsistent page count")
	}
	return rows, total, nil
}

func allVpcStarrocks(api vpcStarrocksAPI) ([]map[string]interface{}, error) {
	var all []map[string]interface{}
	expected := -1
	seen := map[string]bool{}
	for page := 1; page <= 10000; page++ {
		response, err := api.describe("", page)
		if err != nil {
			return nil, err
		}
		rows, total, err := starrocksVpcRows(response)
		if err != nil {
			return nil, err
		}
		if expected < 0 {
			expected = total
		} else if expected != total {
			return nil, fmt.Errorf("StarRocks total changed during pagination")
		}
		for _, row := range rows {
			id, err := serviceVpcOwnerID(row, "InstanceId")
			if err != nil {
				return nil, err
			}
			if seen[id] {
				return nil, fmt.Errorf("StarRocks duplicate instance across pages")
			}
			seen[id] = true
			all = append(all, row)
		}
		if len(all) == total {
			return all, nil
		}
		if len(rows) < 100 || len(all) > total {
			return nil, fmt.Errorf("StarRocks incomplete/invalid pagination")
		}
	}
	return nil, fmt.Errorf("StarRocks pagination exceeded bound")
}

func exactVpcStarrocks(region, vpcID, id string, api vpcStarrocksAPI) (map[string]interface{}, bool, error) {
	if id == "" {
		return nil, false, fmt.Errorf("StarRocks disappearance read requires a nonempty exact instance filter")
	}
	response, err := api.describe(id, 1)
	if err != nil {
		// The live exact-ID DescribeInstances API returns HTTP 400 with this
		// structured code after release, rather than an empty successful page.
		// This exception never applies to the unfiltered account inventory.
		if rdsErrorHasCode(err, "Instance.NotFound") {
			return nil, true, nil
		}
		return nil, false, err
	}
	rows, total, err := starrocksVpcRows(response)
	if err != nil {
		return nil, false, err
	}
	if total == 0 && len(rows) == 0 {
		return nil, true, nil
	}
	if total != 1 || len(rows) != 1 || rows[0]["InstanceId"] != id {
		return nil, false, fmt.Errorf("StarRocks exact instance filter ignored/ambiguous")
	}
	row := rows[0]
	if row["VpcId"] != vpcID || row["RegionId"] != region {
		return nil, false, fmt.Errorf("StarRocks %s VPC/region ownership changed or unknown", id)
	}
	return row, false, nil
}

func sweepVpcStarrocks(region, vpcID string, api vpcStarrocksAPI, pause func(time.Duration), attempts int) error {
	if attempts < 1 {
		return fmt.Errorf("StarRocks wait requires positive bound")
	}
	rows, err := allVpcStarrocks(api)
	if err != nil {
		return err
	}
	var failures []string
	for _, listed := range rows {
		match, err := serviceVpcOwnerMatch(listed, "VpcId", vpcID)
		if err != nil {
			failures = append(failures, err.Error())
			continue
		}
		if !match {
			continue
		}
		id := listed["InstanceId"].(string)
		row, gone, err := exactVpcStarrocks(region, vpcID, id, api)
		if err != nil {
			failures = append(failures, err.Error())
			continue
		}
		if gone {
			continue
		}
		pay, ok := row["PayType"].(string)
		if !ok || !strings.EqualFold(pay, "postPaid") {
			failures = append(failures, "StarRocks "+id+" unknown/subscription payment requires verified unsubscribe capability; retained")
			continue
		}
		switches, ok := row["VSwitches"].([]interface{})
		if !ok || len(switches) == 0 {
			failures = append(failures, "StarRocks "+id+" unknown VSwitch bindings; retained")
			continue
		}
		valid := true
		for _, value := range switches {
			vsw, ok := value.(map[string]interface{})
			if !ok {
				valid = false
				break
			}
			vswID, err := serviceVpcOwnerID(vsw, "VswId")
			if err != nil {
				valid = false
				break
			}
			parent, err := api.vswitch(vswID)
			if err != nil {
				failures = append(failures, err.Error())
				valid = false
				break
			}
			if parent["VSwitchId"] != vswID || parent["VpcId"] != vpcID {
				valid = false
				break
			}
		}
		if !valid {
			failures = append(failures, "StarRocks "+id+" shared/unknown VSwitch bindings; retained")
			continue
		}
		response, err := api.release(id)
		if err != nil {
			failures = append(failures, err.Error())
			continue
		}
		if response["Success"] != true {
			failures = append(failures, fmt.Sprintf("StarRocks %s release not confirmed (code %v)", id, response["ErrCode"]))
			continue
		}
		gone = false
		for attempt := 0; attempt < attempts; attempt++ {
			_, gone, err = exactVpcStarrocks(region, vpcID, id, api)
			if err != nil || gone {
				break
			}
			if attempt+1 < attempts {
				pause(5 * time.Second)
			}
		}
		if err != nil {
			failures = append(failures, err.Error())
		} else if !gone {
			failures = append(failures, "StarRocks "+id+" still present after bounded release wait")
		}
	}
	if len(failures) > 0 {
		return fmt.Errorf("StarRocks VPC owner cleanup incomplete: %s", strings.Join(failures, "; "))
	}
	return nil
}

func starrocksVpcFixture(id, vpc string) map[string]interface{} {
	return map[string]interface{}{"InstanceId": id, "VpcId": vpc, "RegionId": "cn-hangzhou", "PayType": "postPaid", "VSwitches": []interface{}{map[string]interface{}{"VswId": "vsw"}}}
}
func starrocksVpcResponse(total int, rows ...interface{}) map[string]interface{} {
	if rows == nil {
		rows = []interface{}{}
	}
	return map[string]interface{}{"Success": true, "Total": total, "Data": rows}
}

func TestUnitVpcSweepStarrocksNativeTransportUsesRoaQueryAndScopedVpcRead(t *testing.T) {
	var actions []string
	post := func(product, version, action, path string, query, headers map[string]*string, body interface{}, auto bool) (map[string]interface{}, error) {
		actions = append(actions, action)
		if product != "starrocks" || version != "2022-10-19" || body != nil || len(headers) != 0 || !auto {
			t.Fatalf("incorrect ROA request: %s %s %s", product, version, action)
		}
		switch action {
		case "DescribeInstances":
			if path != "/webapi/starrocks/describeInstances" || tea.StringValue(query["InstanceId"]) != "instance" || tea.StringValue(query["PageNumber"]) != "2" || tea.StringValue(query["PageSize"]) != "100" || tea.StringValue(query["RegionId"]) != "cn-hangzhou" {
				t.Fatalf("incorrect Describe query/path: %s %v", path, query)
			}
		case "ReleaseInstance":
			if path != "/webapi/cluster/release" || len(query) != 1 || tea.StringValue(query["InstanceId"]) != "instance" {
				t.Fatalf("incorrect Release query/path: %s %v", path, query)
			}
		default:
			t.Fatalf("unknown action %s", action)
		}
		return map[string]interface{}{"Success": true}, nil
	}
	vpcRPC := func(product, version, action string, query, body map[string]interface{}, auto bool) (map[string]interface{}, error) {
		if action != "DescribeVSwitches" || body["VSwitchId"] != "vsw" || body["RegionId"] != "cn-hangzhou" || body["PageSize"] != 50 {
			t.Fatalf("unscoped VSwitch read: %v", body)
		}
		return map[string]interface{}{"VSwitches": map[string]interface{}{"VSwitch": []interface{}{map[string]interface{}{"VSwitchId": "vsw", "VpcId": "vpc-target"}}}}, nil
	}
	api := starrocksVpcAPIWithTransport("cn-hangzhou", post, vpcRPC)
	if _, err := api.describe("instance", 2); err != nil {
		t.Fatal(err)
	}
	if _, err := api.release("instance"); err != nil {
		t.Fatal(err)
	}
	if _, err := api.vswitch("vsw"); err != nil {
		t.Fatal(err)
	}
	if len(actions) != 2 {
		t.Fatalf("unexpected actions %v", actions)
	}
}

func TestUnitVpcSweepStarrocksStrictTotalPagination(t *testing.T) {
	api := vpcStarrocksAPI{describe: func(id string, page int) (map[string]interface{}, error) {
		if id != "" {
			t.Fatal("unexpected filter")
		}
		var rows []interface{}
		start, end := 0, 100
		if page == 2 {
			start, end = 100, 101
		} else if page != 1 {
			t.Fatalf("unexpected page %d", page)
		}
		for i := start; i < end; i++ {
			rows = append(rows, starrocksVpcFixture("i"+strconv.Itoa(i), "other"))
		}
		response := starrocksVpcResponse(101, rows...)
		response["Total"] = json.Number("101")
		return response, nil
	}}
	if rows, err := allVpcStarrocks(api); err != nil || len(rows) != 101 {
		t.Fatalf("rows=%d err=%v", len(rows), err)
	}
	for _, scenario := range []string{"short-page", "missing-total", "duplicate", "failed-envelope", "unknown-data"} {
		api.describe = func(id string, page int) (map[string]interface{}, error) {
			var rows []interface{}
			if scenario == "duplicate" {
				for i := 0; i < 100; i++ {
					rows = append(rows, starrocksVpcFixture("same", "other"))
				}
			} else {
				rows = append(rows, starrocksVpcFixture("one", "other"))
			}
			response := starrocksVpcResponse(101, rows...)
			switch scenario {
			case "missing-total":
				delete(response, "Total")
			case "failed-envelope":
				response["Success"] = false
			case "unknown-data":
				response["Data"] = nil
			}
			return response, nil
		}
		if _, err := allVpcStarrocks(api); err == nil {
			t.Fatalf("%s accepted", scenario)
		}
	}
}

func TestUnitVpcSweepStarrocksExactTargetSharedBoundaryAndAsyncRelease(t *testing.T) {
	releases, polls, pauses := 0, 0, 0
	api := vpcStarrocksAPI{
		describe: func(id string, page int) (map[string]interface{}, error) {
			if id == "" {
				return starrocksVpcResponse(2, starrocksVpcFixture("target", "vpc-target"), starrocksVpcFixture("protected", "vpc-protected")), nil
			}
			if id != "target" {
				t.Fatal("unrelated instance described for mutation")
			}
			if releases > 0 {
				polls++
				if polls > 1 {
					return starrocksVpcResponse(0), nil
				}
			}
			return starrocksVpcResponse(1, starrocksVpcFixture("target", "vpc-target")), nil
		},
		release: func(id string) (map[string]interface{}, error) {
			if id != "target" {
				t.Fatal("protected instance released")
			}
			releases++
			return map[string]interface{}{"Success": true}, nil
		},
		vswitch: func(id string) (map[string]interface{}, error) {
			return map[string]interface{}{"VSwitchId": id, "VpcId": "vpc-target"}, nil
		},
	}
	if err := sweepVpcStarrocks("cn-hangzhou", "vpc-target", api, func(time.Duration) { pauses++ }, 3); err != nil || releases != 1 || pauses != 1 {
		t.Fatalf("release=%d pause=%d err=%v", releases, pauses, err)
	}
}

func TestUnitVpcSweepStarrocksUnknownSubscriptionSharedAndReadErrorsRetained(t *testing.T) {
	for _, scenario := range []string{"prepaid", "other-vpc", "filter-ignored", "unknown-vswitch", "other-vswitch", "detail-error", "release-false", "wait-error", "timeout"} {
		releases, pauses := 0, 0
		api := vpcStarrocksAPI{
			describe: func(id string, page int) (map[string]interface{}, error) {
				row := starrocksVpcFixture("target", "vpc-target")
				if id == "" {
					return starrocksVpcResponse(1, row), nil
				}
				if scenario == "detail-error" || (scenario == "wait-error" && releases > 0) {
					return nil, &tea.SDKError{Code: tea.String("Forbidden.RAM"), Message: tea.String("denied"), StatusCode: tea.Int(403), Data: nil}
				}
				switch scenario {
				case "prepaid":
					row["PayType"] = "prePaid"
				case "other-vpc":
					row["VpcId"] = "vpc-protected"
				case "filter-ignored":
					return starrocksVpcResponse(2, row, starrocksVpcFixture("other", "vpc-protected")), nil
				case "unknown-vswitch":
					delete(row, "VSwitches")
				}
				return starrocksVpcResponse(1, row), nil
			},
			release: func(id string) (map[string]interface{}, error) {
				if scenario != "release-false" && scenario != "wait-error" && scenario != "timeout" {
					t.Fatalf("%s reached release", scenario)
				}
				releases++
				return map[string]interface{}{"Success": scenario != "release-false"}, nil
			},
			vswitch: func(id string) (map[string]interface{}, error) {
				vpc := "vpc-target"
				if scenario == "other-vswitch" {
					vpc = "vpc-protected"
				}
				return map[string]interface{}{"VSwitchId": id, "VpcId": vpc}, nil
			},
		}
		if err := sweepVpcStarrocks("cn-hangzhou", "vpc-target", api, func(time.Duration) { pauses++ }, 3); err == nil {
			t.Fatalf("%s falsely completed", scenario)
		}
		if scenario == "timeout" && pauses != 2 {
			t.Fatalf("wait not bounded: %d", pauses)
		}
	}
}

func TestUnitVpcSweepStarrocksWrongRegionMakesNoApiCalls(t *testing.T) {
	api := vpcStarrocksAPI{describe: func(string, int) (map[string]interface{}, error) { return nil, fmt.Errorf("must not be called") }}
	if err := vpcStarrocksSweeper("cn-hangzhou", "vpc-target", api).F("cn-beijing"); err == nil {
		t.Fatal("wrong region accepted")
	}
}

func TestUnitVpcSweepStarrocksExactStructuredNotFoundBeforeAndAfterRelease(t *testing.T) {
	for _, alreadyGone := range []bool{true, false} {
		releases := 0
		api := vpcStarrocksAPI{
			describe: func(id string, page int) (map[string]interface{}, error) {
				if id == "" {
					return starrocksVpcResponse(1, starrocksVpcFixture("target", "vpc-target")), nil
				}
				if id != "target" {
					t.Fatalf("missing exact ID filter: %q", id)
				}
				if alreadyGone || releases > 0 {
					return nil, &tea.SDKError{Code: tea.String("Instance.NotFound"), Message: tea.String("Failed to find instance"), StatusCode: tea.Int(400), Data: nil}
				}
				return starrocksVpcResponse(1, starrocksVpcFixture("target", "vpc-target")), nil
			},
			release: func(id string) (map[string]interface{}, error) {
				releases++
				return map[string]interface{}{"Success": true}, nil
			},
			vswitch: func(id string) (map[string]interface{}, error) {
				return map[string]interface{}{"VSwitchId": id, "VpcId": "vpc-target"}, nil
			},
		}
		if err := sweepVpcStarrocks("cn-hangzhou", "vpc-target", api, func(time.Duration) { t.Fatal("known disappearance unnecessarily waited") }, 3); err != nil {
			t.Fatal(err)
		}
		want := 1
		if alreadyGone {
			want = 0
		}
		if releases != want {
			t.Fatalf("alreadyGone=%v releases=%d", alreadyGone, releases)
		}
	}
}

func TestUnitVpcSweepStarrocksNotFoundNeverCoversInventoryOpaqueOrOtherErrors(t *testing.T) {
	notFound := &tea.SDKError{Code: tea.String("Instance.NotFound"), Message: tea.String("Failed to find instance"), StatusCode: tea.Int(400), Data: nil}
	api := vpcStarrocksAPI{describe: func(id string, page int) (map[string]interface{}, error) { return nil, notFound }}
	if _, err := allVpcStarrocks(api); err == nil {
		t.Fatal("unfiltered inventory error falsely means empty account")
	}
	for _, injected := range []error{
		fmt.Errorf("Instance.NotFound"),
		&tea.SDKError{Code: tea.String("Forbidden.RAM"), Message: tea.String("Instance.NotFound denied"), StatusCode: tea.Int(403), Data: nil},
		&tea.SDKError{Code: tea.String("Other.NotFound"), Message: tea.String("not found"), StatusCode: tea.Int(400), Data: nil},
	} {
		api.describe = func(string, int) (map[string]interface{}, error) { return nil, injected }
		if _, gone, err := exactVpcStarrocks("cn-hangzhou", "vpc-target", "target", api); gone || err == nil {
			t.Fatalf("unverified error accepted as gone: %T", injected)
		}
	}
	api.describe = func(string, int) (map[string]interface{}, error) {
		t.Fatal("empty filter called API")
		return nil, notFound
	}
	if _, gone, err := exactVpcStarrocks("cn-hangzhou", "vpc-target", "", api); gone || err == nil {
		t.Fatal("empty exact filter accepted")
	}
}

func TestUnitVpcSweepStarrocksTotalMatchesProductionSdkUseNumber(t *testing.T) {
	decoded, err := util.ReadAsJSON(strings.NewReader(`{"Success":true,"Total":1,"Data":[{"InstanceId":"fixture","VpcId":"vpc-fixture"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	response, ok := decoded.(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected SDK map %T", decoded)
	}
	if _, ok := response["Total"].(json.Number); !ok {
		t.Fatalf("production SDK total type %T", response["Total"])
	}
	// tea-roa then calls tea.Convert, which also preserves UseNumber.
	converted := map[string]interface{}{}
	if err := tea.Convert(map[string]interface{}{"body": response}, &converted); err != nil {
		t.Fatal(err)
	}
	body, ok := converted["body"].(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected ROA body %T", converted["body"])
	}
	rows, total, err := starrocksVpcRows(body)
	if err != nil || total != 1 || len(rows) != 1 {
		t.Fatalf("native SDK shape rows=%d total=%d err=%v", len(rows), total, err)
	}
}

func TestUnitVpcSweepStarrocksTotalRejectsFractionOverflowAndUnknownTypes(t *testing.T) {
	for _, value := range []interface{}{json.Number("0"), json.Number("5"), json.Number("1000000"), 0, int64(5), float64(5)} {
		response := starrocksVpcResponse(0)
		response["Total"] = value
		if _, _, err := starrocksVpcRows(response); err != nil {
			t.Fatalf("valid total %v (%T): %v", value, value, err)
		}
	}
	for _, value := range []interface{}{json.Number("-1"), json.Number("-0"), json.Number("1.0"), json.Number("0.5"), json.Number("1e3"), json.Number("+1"), json.Number("01"), json.Number("1000001"), json.Number("9223372036854775808"), json.Number(""), json.Number("NaN"), "5", true, nil, float64(1.5), -1, int64(1000001)} {
		response := starrocksVpcResponse(0)
		response["Total"] = value
		if _, _, err := starrocksVpcRows(response); err == nil {
			t.Fatalf("invalid total %v (%T) accepted", value, value)
		}
	}
}

func TestUnitVpcDeleteRouterInterfaceDependency(t *testing.T) {
	for _, tc := range []struct {
		name      string
		code      string
		permanent bool
		timeout   time.Duration
		wantError bool
		wantCalls int
	}{
		{name: "dependency_clears", code: "DependencyViolation.RouterInterface", timeout: 20 * time.Second, wantCalls: 2},
		{name: "permanent_dependency", code: "DependencyViolation.RouterInterface", permanent: true, timeout: time.Second, wantError: true, wantCalls: 1},
		{name: "similar_dependency", code: "DependencyViolation.RouterInterfaceOther", timeout: 20 * time.Second, wantError: true, wantCalls: 1},
		{name: "other_dependency", code: "DependencyViolation.NatGateway", timeout: 20 * time.Second, wantError: true, wantCalls: 1},
		{name: "vpc_in_use", code: "Forbidden.VpcInUse", timeout: 20 * time.Second, wantError: true, wantCalls: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			deletes := 0
			client := routerInterfaceDeleteTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				switch r.Form.Get("Action") {
				case "DeleteVpc":
					deletes++
					if r.Form.Get("VpcId") != "vpc-test" {
						t.Error("wrong VPC")
					}
					if deletes == 1 || tc.permanent {
						routerInterfaceDeleteTestError(t, w, tc.code)
						return
					}
					routerInterfaceDeleteTestResponse(t, w, map[string]interface{}{})
				case "DescribeVpcAttribute":
					routerInterfaceDeleteTestError(t, w, "InvalidVpcId.NotFound")
				default:
					t.Errorf("unexpected action %q", r.Form.Get("Action"))
				}
			})
			res := resourceAliCloudVpcVpc()
			res.Timeouts = &schema.ResourceTimeout{Delete: schema.DefaultTimeout(tc.timeout)}
			d := res.Data(nil)
			d.SetId("vpc-test")
			err := res.Delete(d, client)
			if tc.wantError {
				if err == nil || !strings.Contains(err.Error(), tc.code) {
					t.Errorf("expected dependency error %s, got %v", tc.code, err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			if tc.permanent {
				if deletes < 1 || deletes > 3 {
					t.Errorf("unbounded dependency retries: %d", deletes)
				}
			} else if deletes != tc.wantCalls {
				t.Errorf("wanted %d deletes, got %d", tc.wantCalls, deletes)
			}
		})
	}
}
