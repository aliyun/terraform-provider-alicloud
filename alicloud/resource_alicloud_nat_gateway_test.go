package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_nat_gateway", &resource.Sweeper{
		Name: "alicloud_nat_gateway",
		F:    testSweepNatGateways,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_cs_kubernetes",
		},
	})
}

func testSweepNatGateways(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
	}

	var response map[string]interface{}
	action := "DescribeNatGateways"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	natGatewayIds := make([]string, 0)
	service := VpcService{client}
	for {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		if err != nil {
			return fmt.Errorf("Error retrieving Nat Gateways: %s", err)
		}

		resp, err := jsonpath.Get("$.NatGateways.NatGateway", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.NatGateways.NatGateway", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := fmt.Sprint(item["Name"])
			id := fmt.Sprint(item["NatGatewayId"])
			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
				// If a nat gateway name is not set successfully, it should be fetched by vpc name and deleted.
				if skip {
					if need, err := service.needSweepVpc(fmt.Sprint(item["VpcId"]), ""); err == nil {
						skip = !need
					}
				}
				if skip {
					log.Printf("[INFO] Skipping Nat Gateway: %s (%s)", name, id)
					continue
				}
			}
			natGatewayIds = append(natGatewayIds, id)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	for _, id := range natGatewayIds {
		log.Printf("[INFO] Deleting Nat Gateway:  (%s)", id)
		action := "DeleteNatGateway"
		request := map[string]interface{}{
			"NatGatewayId": id,
			"RegionId":     client.RegionId,
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"DependencyViolation.BandwidthPackages"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Nat Gateway (%s): %v", id, err)
		}
	}
	return nil
}

func TestAccAliCloudNatGateway_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nat_gateway.default"
	ra := resourceAttrInit(resourceId, AliCloudNatGatewayMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNatGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snatgateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNatGatewayBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":           "${alicloud_vpc.default.id}",
					"nat_gateway_name": "${var.name}",
					"nat_type":         "Enhanced",
					"vswitch_id":       "${alicloud_vswitch.default.id}",
					"eip_bind_mode":    "MULTI_BINDED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":           CHECKSET,
						"nat_gateway_name": name,
						"nat_type":         "Enhanced",
						"vswitch_id":       CHECKSET,
						"eip_bind_mode":    "MULTI_BINDED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"specification": "Middle",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nat_gateway_name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nat_gateway_name": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nat_type":   "Enhanced",
					"vswitch_id": "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nat_type":   "Enhanced",
						"vswitch_id": CHECKSET,
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
					"eip_bind_mode": "NAT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"eip_bind_mode": "NAT",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"icmp_reply_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"icmp_reply_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"icmp_reply_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"icmp_reply_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"specification":    "Small",
					"description":      name,
					"nat_gateway_name": name,
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
					"deletion_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"specification":       "",
						"description":         name,
						"nat_gateway_name":    name,
						"tags.%":              "2",
						"tags.Created":        "TF-update",
						"tags.For":            "Test-update",
						"deletion_protection": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "force"},
			},
		},
	})
}

func TestAccAliCloudNatGateway_NetworkType(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nat_gateway.default"
	ra := resourceAttrInit(resourceId, AliCloudNatGatewayMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNatGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snatgateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNatGatewayBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":               "${alicloud_vpc.default.id}",
					"nat_gateway_name":     "${var.name}",
					"nat_type":             "Enhanced",
					"vswitch_id":           "${alicloud_vswitch.default.id}",
					"internet_charge_type": "PayByLcu",
					"network_type":         "intranet",
					"private_link_enabled": "true",
					"access_mode": []map[string]interface{}{
						{
							"mode_value":  "tunnel",
							"tunnel_type": "geneve",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":               CHECKSET,
						"nat_gateway_name":     name,
						"nat_type":             "Enhanced",
						"network_type":         "intranet",
						"internet_charge_type": "PayByLcu",
						"vswitch_id":           CHECKSET,
						"private_link_enabled": "true",
						"access_mode.#":        "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "force"},
			},
		},
	})
}

func TestAccAliCloudNatGateway_PayByLcu(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nat_gateway.default"
	ra := resourceAttrInit(resourceId, AliCloudNatGatewayMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNatGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snatgateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNatGatewayBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":               "${alicloud_vpc.default.id}",
					"nat_gateway_name":     "${var.name}",
					"internet_charge_type": "PayByLcu",
					"nat_type":             "Enhanced",
					"vswitch_id":           "${alicloud_vswitch.default.id}",
					"icmp_reply_enabled":   "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":               CHECKSET,
						"nat_gateway_name":     name,
						"internet_charge_type": "PayByLcu",
						"nat_type":             "Enhanced",
						"vswitch_id":           CHECKSET,
						"icmp_reply_enabled":   "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"specification": "Middle",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nat_gateway_name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nat_gateway_name": name + "1",
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
					"deletion_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":      name,
					"nat_gateway_name": name,
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
					"deletion_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"specification":       "",
						"description":         name,
						"nat_gateway_name":    name,
						"tags.%":              "2",
						"tags.Created":        "TF-update",
						"tags.For":            "Test-update",
						"deletion_protection": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "force"},
			},
		},
	})
}

func TestAccAliCloudNatGateway_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nat_gateway.default"
	ra := resourceAttrInit(resourceId, AliCloudNatGatewayMap2)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNatGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snatgateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNatGatewayBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":               "${data.alicloud_vpcs.default.ids.0}",
					"name":                 name,
					"nat_type":             "Enhanced",
					"vswitch_id":           "${data.alicloud_vswitches.default.vswitches.0.id}",
					"description":          name,
					"network_type":         "internet",
					"payment_type":         "Subscription",
					"period":               "12",
					"internet_charge_type": "PayBySpec",
					"specification":        "Middle",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":               CHECKSET,
						"name":                 name,
						"nat_type":             "Enhanced",
						"vswitch_id":           CHECKSET,
						"description":          name,
						"network_type":         "internet",
						"payment_type":         CHECKSET,
						"internet_charge_type": "PayBySpec",
						"specification":        "Middle",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "force"},
			},
		},
	})
}

func TestAccAliCloudNatGateway_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nat_gateway.default"
	ra := resourceAttrInit(resourceId, AliCloudNatGatewayMap2)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNatGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snatgateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNatGatewayBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":               "${data.alicloud_vpcs.default.ids.0}",
					"name":                 name,
					"nat_type":             "Enhanced",
					"vswitch_id":           "${data.alicloud_vswitches.default.vswitches.0.id}",
					"description":          name,
					"network_type":         "internet",
					"instance_charge_type": "PrePaid",
					"period":               "5",
					"internet_charge_type": "PayBySpec",
					"specification":        "Middle",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":               CHECKSET,
						"name":                 name,
						"nat_type":             "Enhanced",
						"vswitch_id":           CHECKSET,
						"description":          name,
						"network_type":         "internet",
						"instance_charge_type": CHECKSET,
						"internet_charge_type": "PayBySpec",
						"specification":        "Middle",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "force", "period"},
			},
		},
	})
}

var AliCloudNatGatewayMap0 = map[string]string{
	"description":          "",
	"dry_run":              NOSET,
	"force":                NOSET,
	"forward_table_ids":    CHECKSET,
	"internet_charge_type": "PayByLcu",
	"nat_type":             "Normal",
	"payment_type":         "PayAsYouGo",
	"period":               NOSET,
	"snat_table_ids":       CHECKSET,
	"specification":        "",
	"status":               "Available",
	"tags.%":               "0",
	"vswitch_id":           "",
	"deletion_protection":  "false",
}

func AliCloudNatGatewayBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
  		available_resource_creation = "VSwitch"
	}

	resource "alicloud_vpc" "default" {
		vpc_name   = var.name
  		cidr_block = "172.16.0.0/12"
	}

	resource "alicloud_vswitch" "default" {
  		vpc_id       = alicloud_vpc.default.id
		cidr_block   = "172.16.0.0/21"
  		zone_id      = data.alicloud_zones.default.zones.0.id
  		vswitch_name = var.name
	}
`, name)
}

var AliCloudNatGatewayMap1 = map[string]string{
	"description":         "",
	"dry_run":             NOSET,
	"force":               NOSET,
	"forward_table_ids":   CHECKSET,
	"nat_type":            "Enhanced",
	"payment_type":        "PayAsYouGo",
	"period":              NOSET,
	"snat_table_ids":      CHECKSET,
	"status":              "Available",
	"tags.%":              "0",
	"deletion_protection": "false",
}

var AliCloudNatGatewayMap2 = map[string]string{
	"description":         CHECKSET,
	"dry_run":             NOSET,
	"force":               NOSET,
	"forward_table_ids":   CHECKSET,
	"nat_type":            CHECKSET,
	"payment_type":        CHECKSET,
	"snat_table_ids":      CHECKSET,
	"status":              "Available",
	"tags.%":              "0",
	"vswitch_id":          "",
	"deletion_protection": "false",
}

func AliCloudNatGatewayBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_vpcs" "default"	{
		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
		vpc_id = "${data.alicloud_vpcs.default.ids.0}"
	}
`, name)
}
