package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
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
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			_, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"DependencyViolation.BandwidthPackages"}) || NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
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
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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

func TestAccAliCloudNatGateway_AvailabilityMode(t *testing.T) {
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
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// availability_mode is set to CrossAZ (cross-zone DR), which is also the
				// API default, so this exercises the CrossAZ branch and confirms it is
				// read back correctly along with the Read wiring.
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":            "${alicloud_vpc.default.id}",
					"nat_gateway_name":  "${var.name}",
					"nat_type":          "Enhanced",
					"vswitch_id":        "${alicloud_vswitch.default.id}",
					"availability_mode": "CrossAZ",
					"force":             "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":            CHECKSET,
						"nat_gateway_name":  name,
						"nat_type":          "Enhanced",
						"vswitch_id":        CHECKSET,
						"availability_mode": "CrossAZ",
						"force":             "true",
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

// TestAccAliCloudNatGateway_AvailabilityModeSingleAZ covers explicitly setting
// availability_mode = SingleAZ (single-zone DR).
//
// NOTE: single-zone NAT is an account-level entitlement. On accounts without it,
// CreateNatGateway returns InvalidParameter.SingleZone ("current user not support
// create single zone nat"). This test pins the entitled account via
// ALICLOUD_ACCESS_KEY_1/SECRET_KEY_1: the keys are embedded in an explicit
// provider block so the default ALICLOUD_ACCESS_KEY cannot override them, and
// the env vars are also overridden so every SDK client init path uses the
// entitled account (mirrors resource_manager_handshake_acceptance_test.go).
func TestAccAliCloudNatGateway_AvailabilityModeSingleAZ(t *testing.T) {
	ak1 := strings.TrimSpace(os.Getenv("ALICLOUD_ACCESS_KEY_1"))
	sk1 := strings.TrimSpace(os.Getenv("ALICLOUD_SECRET_KEY_1"))
	if ak1 == "" || sk1 == "" {
		t.Skipf("Skipping: set ALICLOUD_ACCESS_KEY_1/SECRET_KEY_1 to run single-zone NAT on an entitled account")
	}
	// Override the default credential env so every code path that reads
	// ALICLOUD_ACCESS_KEY (testAccPreCheck, SDK client init) uses the
	// single-zone-NAT-entitled account.
	os.Setenv("ALICLOUD_ACCESS_KEY", ak1)
	os.Setenv("ALICLOUD_SECRET_KEY", sk1)

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
	// Explicit provider credentials so the default ALICLOUD_ACCESS_KEY cannot
	// override the single-zone-NAT-entitled account pinned by ALICLOUD_ACCESS_KEY_1.
	singleZoneNatDependence := func(n string) string {
		return fmt.Sprintf(`
	provider "alicloud" {
	  access_key = "%s"
	  secret_key = "%s"
	}

%s
`, ak1, sk1, AliCloudNatGatewayBasicDependence0(n))
	}
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, singleZoneNatDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":            "${alicloud_vpc.default.id}",
					"nat_gateway_name":  "${var.name}",
					"nat_type":          "Enhanced",
					"vswitch_id":        "${alicloud_vswitch.default.id}",
					"availability_mode": "SingleAZ",
					"force":             "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":            CHECKSET,
						"nat_gateway_name":  name,
						"nat_type":          "Enhanced",
						"vswitch_id":        CHECKSET,
						"availability_mode": "SingleAZ",
						"force":             "true",
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
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
	t.Skipf("The PrePaid Nat Gateway has not supported any more.")
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
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
	t.Skipf("The PrePaid Nat Gateway has not supported any more.")
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
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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

// TestUnitAliCloudNatGatewayDelete covers the delete error-branch routing that
// acceptance tests cannot reproduce: the IncorrectStatus.NatGateway race where
// DeleteNatGateway reports the gateway is already being deleted. The fix keeps
// idempotent success for not-found errors but falls through to the state waiter
// for IncorrectStatus.NatGateway so deletion is confirmed before reporting done.
//
// A no-credentials connectivity.AliyunClient is built locally so the test runs
// (rather than skips) in CI envs without ALICLOUD_ACCESS_KEY: Config.Client with
// SkipRegionValidation, a dummy StaticAK credential and a pre-populated "vpc"
// endpoint cache initializes client.config without any network call, so the
// inlined GetRetryTimeout and rpcRequest do not panic. The tea-rpc client.DoRequest
// is stubbed to inject the DeleteNatGateway error under test and to answer the
// waiter's DescribeNatGateways poll with a not-found error, driving the real
// NatGatewayStateRefreshFunc into its production (nil, "", nil) NotFound idiom so the
// empty-Target waiter confirms disappearance; reachability of the waiter is observed
// through that DescribeNatGateways call rather than by stubbing the refresh itself
// (which Go 1.25 inlines and cannot be patched).
func TestUnitAliCloudNatGatewayDelete(t *testing.T) {
	endpoints := sharedEndpointCacheForRegion("cn-hangzhou")
	endpoints.Store("vpc", "vpc.cn-hangzhou.aliyuncs.com")
	conf := &connectivity.Config{
		Region:               connectivity.Region("cn-hangzhou"),
		RegionId:             "cn-hangzhou",
		AccessKey:            "dummy-access-key",
		SecretKey:            "dummy-secret-key",
		Protocol:             "HTTPS",
		Endpoints:            endpoints,
		SignVersion:          &signVersion,
		SkipRegionValidation: true,
		AccountType:          "Domestic",
	}
	cred, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId(conf.AccessKey).
		SetAccessKeySecret(conf.SecretKey))
	if err != nil {
		t.Fatalf("failed to build dummy credential: %s", err)
	}
	conf.Credential = cred
	aliClient, err := conf.Client()
	if err != nil {
		t.Fatalf("failed to build dummy client: %s", err)
	}

	p := Provider().ResourcesMap
	newData := func() *schema.ResourceData {
		d, _ := schema.InternalMap(p["alicloud_nat_gateway"].Schema).Data(nil, nil)
		d.SetId("ngw-mock-id")
		return d
	}

	sdkErr := func(code string) error {
		return &tea.SDKError{
			Code:       tea.String(code),
			Data:       tea.String(code),
			Message:    tea.String(code),
			StatusCode: tea.Int(400),
		}
	}

	cases := []struct {
		name     string
		resp     map[string]interface{}
		err      error
		wantErr  bool
		wantWait bool
	}{
		{"DeleteSucceedsAndWaitsForConfirmation", map[string]interface{}{}, nil, false, true},
		{"IncorrectStatusFallsThroughToWaiter", nil, sdkErr("IncorrectStatus.NatGateway"), false, true},
		{"InvalidNatGatewayIdNotFoundIsIdempotent", nil, sdkErr("InvalidNatGatewayId.NotFound"), false, false},
		{"InstanceNotExistsIsIdempotent", nil, sdkErr("INSTANCE_NOT_EXISTS"), false, false},
		{"InvalidRegionIdNotFoundIsIdempotent", nil, sdkErr("InvalidRegionId.NotFound"), false, false},
		{"UnexpectedErrorIsWrappedNotWaited", nil, sdkErr("Forbidden.RAM"), true, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := newData()
			var describeCalled int32

			// The real VpcService.NatGatewayStateRefreshFunc runs unchanged; reachability
			// of the state waiter is observed through the DescribeNatGateways poll it
			// issues. Returning a not-found error there drives the refresh into its
			// production (nil, "", nil) NotFound idiom, which the empty-Target waiter
			// treats as "gone" and succeeds — so the delete reports done only after the
			// waiter confirms disappearance, exactly as on the live API.
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest",
				func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
					if action != nil {
						switch *action {
						case "DeleteNatGateway":
							return tc.resp, tc.err
						case "DescribeNatGateways":
							atomic.AddInt32(&describeCalled, 1)
							return nil, sdkErr("InvalidNatGatewayId.NotFound")
						}
					}
					return nil, nil
				})
			defer patches.Reset()

			err := resourceAliCloudNatGatewayDelete(d, aliClient)

			if tc.wantErr {
				assert.NotNil(t, err, "expected the delete error to be surfaced")
			} else {
				assert.Nil(t, err, "expected the delete to succeed")
			}
			assert.Equal(t, tc.wantWait, atomic.LoadInt32(&describeCalled) > 0,
				"waiter reachability mismatch for %s", tc.name)
		})
	}
}
