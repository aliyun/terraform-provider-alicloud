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

// Test NATGateway SnatEntry. >>> Resource test cases, automatically generated.
// Case 全生命周期_SnatEntry source_cidr 8016
func TestAccAliCloudNATGatewaySnatEntry_basic8016(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_snat_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudNATGatewaySnatEntryMap8016)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NATGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNATGatewaySnatEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snatgatewaysnatentry%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNATGatewaySnatEntryBasicDependence8016)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"snat_ip":       "${alicloud_vpc_nat_ip.default.nat_ip}" + "," + "${alicloud_vpc_nat_ip.update.nat_ip}",
					"snat_table_id": "${alicloud_nat_gateway.default.snat_table_ids}",
					"source_cidr":   "${cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_ip":       CHECKSET,
						"snat_table_id": CHECKSET,
						"source_cidr":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"eip_affinity": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"eip_affinity": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"eip_affinity": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"eip_affinity": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"snat_entry_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_entry_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"snat_ip": "${alicloud_vpc_nat_ip.default.nat_ip}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_ip": CHECKSET,
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

func TestAccAliCloudNATGatewaySnatEntry_basic8016_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_snat_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudNATGatewaySnatEntryMap8016)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NATGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNATGatewaySnatEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snatgatewaysnatentry%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNATGatewaySnatEntryBasicDependence8016)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"snat_ip":         "${alicloud_vpc_nat_ip.default.nat_ip}" + "," + "${alicloud_vpc_nat_ip.update.nat_ip}",
					"snat_table_id":   "${alicloud_nat_gateway.default.snat_table_ids}",
					"source_cidr":     "${cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)}",
					"snat_entry_name": name,
					"eip_affinity":    "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_ip":         CHECKSET,
						"snat_table_id":   CHECKSET,
						"source_cidr":     CHECKSET,
						"snat_entry_name": name,
						"eip_affinity":    "1",
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

// Case 全生命周期_SnatEntry source_vswitch_id 8018
func TestAccAliCloudNATGatewaySnatEntry_basic8018(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_snat_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudNATGatewaySnatEntryMap8016)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NATGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNATGatewaySnatEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snatgatewaysnatentry%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNATGatewaySnatEntryBasicDependence8016)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"snat_ip":           "${alicloud_vpc_nat_ip.default.nat_ip}" + "," + "${alicloud_vpc_nat_ip.update.nat_ip}",
					"snat_table_id":     "${alicloud_nat_gateway.default.snat_table_ids}",
					"source_vswitch_id": "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_ip":           CHECKSET,
						"snat_table_id":     CHECKSET,
						"source_vswitch_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"eip_affinity": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"eip_affinity": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"eip_affinity": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"eip_affinity": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"snat_entry_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_entry_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"snat_ip": "${alicloud_vpc_nat_ip.default.nat_ip}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_ip": CHECKSET,
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

func TestAccAliCloudNATGatewaySnatEntry_basic8018_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_snat_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudNATGatewaySnatEntryMap8016)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NATGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNATGatewaySnatEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snatgatewaysnatentry%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNATGatewaySnatEntryBasicDependence8016)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"snat_ip":           "${alicloud_vpc_nat_ip.default.nat_ip}" + "," + "${alicloud_vpc_nat_ip.update.nat_ip}",
					"snat_table_id":     "${alicloud_nat_gateway.default.snat_table_ids}",
					"source_vswitch_id": "${alicloud_vswitch.default.id}",
					"snat_entry_name":   name,
					"eip_affinity":      "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_ip":           CHECKSET,
						"snat_table_id":     CHECKSET,
						"source_vswitch_id": CHECKSET,
						"snat_entry_name":   name,
						"eip_affinity":      "1",
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

var AliCloudNATGatewaySnatEntryMap8016 = map[string]string{
	"snat_entry_id": CHECKSET,
	"status":        CHECKSET,
}

func AliCloudNATGatewaySnatEntryBasicDependence8016(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "172.16.0.0/12"
	}

	resource "alicloud_vswitch" "default" {
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "172.16.0.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
  		vswitch_name = var.name
	}

	resource "alicloud_nat_gateway" "default" {
  		vpc_id               = alicloud_vpc.default.id
  		network_type         = "intranet"
  		nat_gateway_name     = var.name
  		vswitch_id           = alicloud_vswitch.default.id
  		nat_type             = "Enhanced"
  		internet_charge_type = "PayByLcu"
	}

	resource "alicloud_vpc_nat_ip" "default" {
  		nat_ip         = "172.16.0.66"
  		nat_ip_cidr    = alicloud_vswitch.default.cidr_block
  		nat_gateway_id = alicloud_nat_gateway.default.id
	}

	resource "alicloud_vpc_nat_ip" "update" {
  		nat_ip         = "172.16.0.88"
  		nat_ip_cidr    = alicloud_vswitch.default.cidr_block
  		nat_gateway_id = alicloud_nat_gateway.default.id
	}
`, name)
}

// Test NATGateway SnatEntry. <<< Resource test cases, automatically generated.

func init() {
	resource.AddTestSweepers("alicloud_snat_entry", &resource.Sweeper{
		Name: "alicloud_snat_entry",
		F:    testSweepSnatEntries,
	})
}

// sweepListNatGateways lists every NAT gateway in the client's region. SNAT/DNAT
// entries do not have a region-wide list API of their own; they belong to a NAT
// gateway's SnatTable/ForwardTable, so the snat_entry and forward_entry sweepers
// must first enumerate NAT gateways to discover those table ids. Shared by both.
func sweepListNatGateways(client *connectivity.AliyunClient) ([]map[string]interface{}, error) {
	items := make([]map[string]interface{}, 0)
	action := "DescribeNatGateways"
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	for {
		response, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		if err != nil {
			return nil, fmt.Errorf("error retrieving Nat Gateways: %s", err)
		}
		resp, err := jsonpath.Get("$.NatGateways.NatGateway", response)
		if err != nil {
			return nil, WrapErrorf(err, FailedGetAttributeMsg, action, "$.NatGateways.NatGateway", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				items = append(items, item)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return items, nil
}

// sweepNatTableIds extracts the table id list from a DescribeNatGateways item.
// The API nests ids as e.g. {"SnatTableIds": {"SnatTableId": [...]}}.
func sweepNatTableIds(nat map[string]interface{}, key, subkey string) []string {
	ids := []string{}
	raw, ok := nat[key].(map[string]interface{})
	if !ok {
		return ids
	}
	arr, ok := raw[subkey].([]interface{})
	if !ok {
		return ids
	}
	for _, id := range arr {
		ids = append(ids, fmt.Sprint(id))
	}
	return ids
}

func testSweepSnatEntries(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf-test-",
		"tf_test_",
		"tf-example",
		"tf_example",
	}

	natGateways, err := sweepListNatGateways(client)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return nil
	}

	for _, nat := range natGateways {
		natName := fmt.Sprint(nat["Name"])
		natID := fmt.Sprint(nat["NatGatewayId"])
		// If the parent NAT gateway is explicitly do-not-delete, protect its entries too.
		if protectedSweepVpcName(natName) {
			log.Printf("[INFO] Skipping SnatEntries of NAT Gateway %s (%s): do-not-delete parent", natName, natID)
			continue
		}
		for _, snatTableId := range sweepNatTableIds(nat, "SnatTableIds", "SnatTableId") {
			entries, err := sweepListSnatEntries(client, snatTableId)
			if err != nil {
				log.Printf("[ERROR] Failed to list SnatEntries for table %s: %s", snatTableId, err)
				continue
			}
			for _, entry := range entries {
				name := fmt.Sprint(entry["SnatEntryName"])
				snatEntryId := fmt.Sprint(entry["SnatEntryId"])
				if protectedSweepVpcName(name) {
					log.Printf("[INFO] Skipping SnatEntry %s (%s): do-not-delete name", name, snatEntryId)
					continue
				}
				if !sweepAll() {
					skip := true
					for _, prefix := range prefixes {
						if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
							skip = false
							break
						}
					}
					if skip {
						log.Printf("[INFO] Skipping SnatEntry: %s (%s)", name, snatEntryId)
						continue
					}
				}
				log.Printf("[INFO] Deleting SnatEntry: %s (%s) in table %s", name, snatEntryId, snatTableId)
				action := "DeleteSnatEntry"
				request := map[string]interface{}{
					"SnatTableId": snatTableId,
					"SnatEntryId": snatEntryId,
					"RegionId":    client.RegionId,
					"ClientToken": buildClientToken(action),
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(5*time.Minute, func() *retry.RetryError {
					_, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)
					if err != nil {
						if IsExpectedErrors(err, []string{"IncorretSnatEntryStatus", "IncorrectStatus.NATGW", "OperationConflict"}) || NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				if err != nil {
					if IsExpectedErrors(err, []string{"InvalidSnatTableId.NotFound", "InvalidSnatEntryId.NotFound"}) || NotFoundError(err) {
						continue
					}
					log.Printf("[ERROR] Failed to delete SnatEntry %s (%s): %s", name, snatEntryId, err)
				}
			}
		}
	}
	return nil
}

func sweepListSnatEntries(client *connectivity.AliyunClient, snatTableId string) ([]map[string]interface{}, error) {
	items := make([]map[string]interface{}, 0)
	action := "DescribeSnatTableEntries"
	request := map[string]interface{}{
		"SnatTableId": snatTableId,
		"RegionId":    client.RegionId,
		"PageSize":    PageSizeLarge,
		"PageNumber":  1,
	}
	for {
		response, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		if err != nil {
			return nil, err
		}
		resp, err := jsonpath.Get("$.SnatTableEntries.SnatTableEntry", response)
		if err != nil {
			return nil, err
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				items = append(items, item)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return items, nil
}
