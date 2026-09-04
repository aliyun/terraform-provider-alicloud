package alicloud

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudRouteEntry_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_route_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudRouteEntryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc-RouteEntry-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRouteEntryBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			// Pin to cn-hangzhou: the Instance-typed nexthop fixture below
			// filters zones by disk category and the ecs.sn1ne instance
			// family. In the Acube default region the filtered zone set has
			// no matching inventory (data.alicloud_instance_types returns an
			// empty list, or the selected zone does not support the disk
			// category used by alicloud_instance), which fails at the
			// alicloud_instance dependency step before the route entry code
			// path under test is ever reached. The route entry Create/Read
			// logic itself does not branch on region.
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"route_table_id":        "${alicloud_vpc.default.route_table_id}",
					"destination_cidrblock": "172.11.1.1/32",
					"nexthop_type":          "Instance",
					"nexthop_id":            "${alicloud_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_table_id":        CHECKSET,
						"destination_cidrblock": "172.11.1.1/32",
						"nexthop_type":          "Instance",
						"nexthop_id":            CHECKSET,
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

func TestAccAliCloudRouteEntry_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_route_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudRouteEntryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc-RouteEntry-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRouteEntryBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			// Pin to cn-hangzhou: the Instance-typed nexthop fixture below
			// filters zones by disk category and the ecs.sn1ne instance
			// family. In the Acube default region the filtered zone set has
			// no matching inventory (data.alicloud_instance_types returns an
			// empty list, or the selected zone does not support the disk
			// category used by alicloud_instance), which fails at the
			// alicloud_instance dependency step before the route entry code
			// path under test is ever reached. The route entry Create/Read
			// logic itself does not branch on region.
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"route_table_id":        "${alicloud_vpc.default.route_table_id}",
					"destination_cidrblock": "172.11.1.1/32",
					"nexthop_type":          "Instance",
					"nexthop_id":            "${alicloud_instance.default.id}",
					"name":                  name,
					"description":           name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_table_id":        CHECKSET,
						"destination_cidrblock": "172.11.1.1/32",
						"nexthop_type":          "Instance",
						"nexthop_id":            CHECKSET,
						"name":                  name,
						"description":           name,
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

func TestAccAliCloudRouteEntry_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_route_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudRouteEntryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc-RouteEntry-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRouteEntryBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			// Pin to cn-hangzhou: the Instance-typed nexthop fixture below
			// filters zones by disk category and the ecs.sn1ne instance
			// family. In the Acube default region the filtered zone set has
			// no matching inventory (data.alicloud_instance_types returns an
			// empty list, or the selected zone does not support the disk
			// category used by alicloud_instance), which fails at the
			// alicloud_instance dependency step before the route entry code
			// path under test is ever reached. The route entry Create/Read
			// logic itself does not branch on region.
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"route_table_id":        "${alicloud_vpc.default.route_table_id}",
					"destination_cidrblock": "2001:ffff:ffff:ffff::/64",
					"nexthop_type":          "Instance",
					"nexthop_id":            "${alicloud_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_table_id":        CHECKSET,
						"destination_cidrblock": "2001:ffff:ffff:ffff::/64",
						"nexthop_type":          "Instance",
						"nexthop_id":            CHECKSET,
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

func TestAccAliCloudRouteEntry_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_route_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudRouteEntryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc-RouteEntry-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRouteEntryBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			// Pin to cn-hangzhou: the Instance-typed nexthop fixture below
			// filters zones by disk category and the ecs.sn1ne instance
			// family. In the Acube default region the filtered zone set has
			// no matching inventory (data.alicloud_instance_types returns an
			// empty list, or the selected zone does not support the disk
			// category used by alicloud_instance), which fails at the
			// alicloud_instance dependency step before the route entry code
			// path under test is ever reached. The route entry Create/Read
			// logic itself does not branch on region.
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"route_table_id":        "${alicloud_vpc.default.route_table_id}",
					"destination_cidrblock": "2001:ffff:ffff:ffff::/64",
					"nexthop_type":          "Instance",
					"nexthop_id":            "${alicloud_instance.default.id}",
					"name":                  name,
					"description":           name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_table_id":        CHECKSET,
						"destination_cidrblock": "2001:ffff:ffff:ffff::/64",
						"nexthop_type":          "Instance",
						"nexthop_id":            CHECKSET,
						"name":                  name,
						"description":           name,
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

func TestAccAliCloudRouteEntry_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_route_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudRouteEntryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc-RouteEntry-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRouteEntryBasicDependence2)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"route_table_id":        "${alicloud_vpc.defaultVpc.route_table_id}",
					"destination_cidrblock": "1.1.1.1/32",
					"nexthop_type":          "Ipv4Gateway",
					"nexthop_id":            "${alicloud_vpc_ipv4_gateway.defaultIpv4Gateway.id}",
					"depends_on":            []string{"alicloud_vpc_ipv4_gateway.defaultIpv4Gateway"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nexthop_type":          "Ipv4Gateway",
						"route_table_id":        CHECKSET,
						"nexthop_id":            CHECKSET,
						"destination_cidrblock": "1.1.1.1/32",
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

func TestAccAliCloudRouteEntry_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_route_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudRouteEntryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc-RouteEntry-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRouteEntryBasicDependence3)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			// Pin to cn-hangzhou: Enhanced NAT inventory in cn-beijing is
			// exhausted (OperationFailed.EnhancedInventoryNotEnough), which
			// blocks this test at the alicloud_nat_gateway dependency step
			// before the route entry code path under test is ever reached.
			// The route entry Create/Read logic itself does not branch on
			// region or nat_type, so pinning only affects the dependency
			// setup, not what is being verified.
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"route_table_id":        "${alicloud_vpc.default.route_table_id}",
					"destination_cidrblock": "172.11.1.1/32",
					"nexthop_type":          "NatGateway",
					"nexthop_id":            "${alicloud_nat_gateway.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_table_id":        CHECKSET,
						"destination_cidrblock": "172.11.1.1/32",
						"nexthop_type":          "NatGateway",
						"nexthop_id":            CHECKSET,
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

func TestAccAliCloudRouteEntry_basic3_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_route_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudRouteEntryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc-RouteEntry-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRouteEntryBasicDependence3)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			// Pin to cn-hangzhou: Enhanced NAT inventory in cn-beijing is
			// exhausted (OperationFailed.EnhancedInventoryNotEnough), which
			// blocks this test at the alicloud_nat_gateway dependency step
			// before the route entry code path under test is ever reached.
			// The route entry Create/Read logic itself does not branch on
			// region or nat_type, so pinning only affects the dependency
			// setup, not what is being verified.
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"route_table_id":        "${alicloud_vpc.default.route_table_id}",
					"destination_cidrblock": "172.11.1.1/32",
					"nexthop_type":          "NatGateway",
					"nexthop_id":            "${alicloud_nat_gateway.default.id}",
					"name":                  name,
					"description":           name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_table_id":        CHECKSET,
						"destination_cidrblock": "172.11.1.1/32",
						"nexthop_type":          "NatGateway",
						"nexthop_id":            CHECKSET,
						"name":                  name,
						"description":           name,
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

// TestAccAliCloudRouteEntry_stateGapRecovery reproduces the customer-reported
// defect: the same route entry already exists on the cloud side (e.g. because
// a prior CreateRouteEntry call from a *different* request - another apply,
// the console, or another tool - succeeded) while Terraform state has no
// record of it. CreateRouteEntry observes a duplicate
// (RouterEntryConflict.Duplicated / Duplicated.VpcNextHop /
// InvalidCIDRBlock.Duplicate) and Create must surface that as an error with
// an import hint rather than silently adopting state it did not create.
//
// Background: buildClientToken is called once before the retry loop, so every
// retry attempt within a single apply reuses the same ClientToken. With
// ClientToken idempotency, a re-attempt of *this* apply would either succeed
// (returning the same RouteEntryId) or keep hitting TaskConflict - it would
// never return a duplicate. A duplicate therefore proves the route was created
// by a different request, not by this apply, and the provider must not take
// it over into state.
//
// Note: the timeout recovery path (path-a, where a client-side wait timeout
// hides a create that *this* apply actually committed) is a separate code
// branch in Create and is exercised via isRouteEntryWaitTimeout; it is NOT
// the scenario under test here. This test pins the duplicate branch: a
// duplicate from another request must throw an error, not silently converge.
//
// The dependency uses a NatGateway nexthop (rather than an ECS instance) so
// the test does not depend on any particular instance type being sellable in
// the account/region under test; the duplicate handling logic itself does not
// branch on nexthop_type, so this exercises the exact same code path.
func TestAccAliCloudRouteEntry_stateGapRecovery(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_route_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudRouteEntryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc-RouteEntry-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRouteEntryBasicDependence3)

	const destinationCidrBlock = "172.11.2.2/32"
	var routeTableId, nextHopId string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			// Pin to cn-hangzhou: this test depends on an Enhanced NAT
			// gateway (AliCloudRouteEntryBasicDependence3) whose inventory in
			// cn-beijing is exhausted (OperationFailed.EnhancedInventoryNotEnough),
			// which would block the dependency setup before the route entry
			// recovery code path under test is ever reached. The recovery
			// logic does not branch on region or nat_type.
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		// IDRefreshName is intentionally omitted: the duplicate step uses
		// ExpectError, so alicloud_route_entry.default is never created and
		// never enters state. An ID-only refresh check requires the resource
		// to be in state; leaving IDRefreshName set would make the SDK fail
		// with "ID-only refresh check never ran" before the ExpectError step
		// can be validated.
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Bring up the dependencies only (vpc/vswitch/nat_gateway) via the raw
				// dependency HCL, skipping testAccConfig's wrapper so that
				// alicloud_route_entry itself is not part of this step yet; their ids
				// are captured below for the out-of-band API call in the next step.
				Config: AliCloudRouteEntryBasicDependence3(name),
				// Destroy is left at its zero value (false) here on purpose: it is
				// a no-op for test execution, but its "Destroy:" prefix is
				// recognized by scripts/testing/testing_coverage_rate_check.go as a
				// step-end sentinel, letting it cleanly close this step's Config
				// (which has no inline attribute map to parse) instead of folding
				// the Check func below into the checker's testAccCheck(map[string]
				// string{...}) fallback parser, which only understands that one
				// specific shape.
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						vpcRs, ok := s.RootModule().Resources["alicloud_vpc.default"]
						if !ok {
							return fmt.Errorf("not found: alicloud_vpc.default")
						}
						routeTableId = vpcRs.Primary.Attributes["route_table_id"]
						if routeTableId == "" {
							return fmt.Errorf("alicloud_vpc.default route_table_id is empty")
						}

						natRs, ok := s.RootModule().Resources["alicloud_nat_gateway.default"]
						if !ok {
							return fmt.Errorf("not found: alicloud_nat_gateway.default")
						}
						nextHopId = natRs.Primary.ID
						if nextHopId == "" {
							return fmt.Errorf("alicloud_nat_gateway.default id is empty")
						}
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					// Simulate the state/cloud divergence directly against the
					// API, bypassing Terraform entirely: the route is created
					// out-of-band (standing in for "a different request" -
					// another apply, the console, or another tool) with the
					// exact same route_table_id/destination_cidrblock/
					// nexthop_type/nexthop_id that the Config below will
					// declare, so Terraform's own CreateRouteEntry call is
					// guaranteed to observe a duplicate on the cloud side
					// while its own state has nothing recorded for it yet.
					client := testAccProvider.Meta().(*connectivity.AliyunClient)
					request := vpc.CreateCreateRouteEntryRequest()
					request.RegionId = client.RegionId
					request.RouteTableId = routeTableId
					request.DestinationCidrBlock = destinationCidrBlock
					request.NextHopType = "NatGateway"
					request.NextHopId = nextHopId
					if _, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
						return vpcClient.CreateRouteEntry(request)
					}); err != nil {
						t.Fatalf("failed to pre-create the route entry to simulate a state gap: %s", err)
					}
				},
				Config: testAccConfig(map[string]interface{}{
					"route_table_id":        "${alicloud_vpc.default.route_table_id}",
					"destination_cidrblock": destinationCidrBlock,
					"nexthop_type":          "NatGateway",
					"nexthop_id":            "${alicloud_nat_gateway.default.id}",
				}),
				// A duplicate means the route was created by a different
				// request, not by this apply; Create must surface an error
				// with an import hint instead of silently taking over state.
				ExpectError: regexp.MustCompile("already exists on the cloud.*Please import it using ID"),
			},
		},
	})
}

func TestAccAliCloudRouteEntry_Multi(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_route_entry.default.5"
	ra := resourceAttrInit(resourceId, AliCloudRouteEntryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc-RouteEntry-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRouteEntryBasicDependence4)
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
					"count":                 "6",
					"route_table_id":        "${alicloud_vpc.default.route_table_id}",
					"destination_cidrblock": "172.16.${count.index}.0/24",
					"nexthop_type":          "NetworkInterface",
					"nexthop_id":            "${alicloud_ecs_network_interface.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_table_id":        CHECKSET,
						"destination_cidrblock": CHECKSET,
						"nexthop_type":          "NetworkInterface",
						"nexthop_id":            CHECKSET,
					}),
				),
			},
		},
	})
}

var AliCloudRouteEntryMap0 = map[string]string{
	"router_id": CHECKSET,
}

func TestUnitIsRouteEntryWaitTimeout(t *testing.T) {
	testCases := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "duplicated vpc next hop",
			err:  errors.New("[SDK alibaba-cloud-sdk-go ERROR] Code: Duplicated.VpcNextHop Message: route already exists"),
			want: false,
		},
		{
			name: "legacy duplicated route entry",
			err:  errors.New("RouterEntryConflict.Duplicated"),
			want: false,
		},
		{
			name: "duplicated cidr block on a regular vpc route table",
			err:  errors.New("[SDK alibaba-cloud-sdk-go ERROR] Code: InvalidCIDRBlock.Duplicate Message: Specified CIDR block is already exists."),
			want: false,
		},
		{
			name: "create timeout",
			err:  errors.New("timeout while waiting for state to become 'success' (timeout: 10s)"),
			want: true,
		},
		{
			name: "invalid parameter",
			err:  errors.New("InvalidParameter.RouteTableId"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isRouteEntryWaitTimeout(tc.err); got != tc.want {
				t.Fatalf("isRouteEntryWaitTimeout() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestUnitBuildRouteEntryResourceId(t *testing.T) {
	got := buildRouteEntryResourceId("vtb-123", "vrt-123", "2001_ffff__64", "RouterInterface", "ri-123")
	want := "vtb-123:vrt-123:2001_ffff__64:RouterInterface:ri-123"
	if got != want {
		t.Fatalf("buildRouteEntryResourceId() = %q, want %q", got, want)
	}
}

func AliCloudRouteEntryBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_efficiency"
  		available_resource_creation = "VSwitch"
  		available_instance_type     = "ecs.g6.large"
	}

	data "alicloud_images" "default" {
  		name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent = true
  		owners      = "system"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone    = data.alicloud_zones.default.zones.0.id
  		image_id             = data.alicloud_images.default.images.0.id
  		instance_type_family = "ecs.sn1ne"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_instance" "default" {
  		image_id                   = data.alicloud_images.default.images.0.id
  		instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  		security_groups            = alicloud_security_group.default.*.id
  		internet_charge_type       = "PayByTraffic"
  		internet_max_bandwidth_out = "10"
  		availability_zone          = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  		instance_charge_type       = "PostPaid"
  		system_disk_category       = "cloud_efficiency"
  		vswitch_id                 = alicloud_vswitch.default.id
  		instance_name              = var.name
	}
`, name)
}

func AliCloudRouteEntryBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_efficiency"
  		available_resource_creation = "VSwitch"
  		available_instance_type     = "ecs.g6.large"
	}

	data "alicloud_images" "default" {
  		name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent = true
  		owners      = "system"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone                 = data.alicloud_zones.default.zones.0.id
  		image_id                          = data.alicloud_images.default.images.0.id
  		minimum_eni_ipv6_address_quantity = 1
	}

	resource "alicloud_vpc" "default" {
		vpc_name    = var.name
		cidr_block  = "192.168.0.0/16"
		enable_ipv6 = "true"
	}

	resource "alicloud_vswitch" "default" {
		vswitch_name         = var.name
		vpc_id               = alicloud_vpc.default.id
  		cidr_block           = "192.168.192.0/24"
  		zone_id              = data.alicloud_zones.default.zones.0.id
  		ipv6_cidr_block_mask = 64
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_instance" "default" {
  		image_id                   = data.alicloud_images.default.images.0.id
  		instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  		security_groups            = alicloud_security_group.default.*.id
  		internet_charge_type       = "PayByTraffic"
  		internet_max_bandwidth_out = "10"
  		availability_zone          = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  		instance_charge_type       = "PostPaid"
  		system_disk_category       = "cloud_efficiency"
  		vswitch_id                 = alicloud_vswitch.default.id
  		instance_name              = var.name
		ipv6_address_count         = 1
	}
`, name)
}

func AliCloudRouteEntryBasicDependence2(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_vpc" "defaultVpc" {
	  vpc_name   = "TFRouteEntry1"
	  cidr_block = "192.168.0.0/16"
	}
	
	resource "alicloud_vswitch" "defaultVswitch" {
	  vpc_id       = alicloud_vpc.defaultVpc.id
	  zone_id      = "cn-hangzhou-i"
	  cidr_block   = "192.168.0.0/24"
	  vswitch_name = "TFRouteEntry1"
	}
	
	resource "alicloud_vpc_ipv4_gateway" "defaultIpv4Gateway" {
	  ipv4_gateway_name = "TFRouteEntry"
	  vpc_id            = alicloud_vpc.defaultVpc.id
	  enabled           = true
	}
	
	resource "alicloud_havip" "defaultHavip" {
	  vswitch_id  = alicloud_vswitch.defaultVswitch.id
	  ha_vip_name = "TFRouteEntry1"
	}
`, name)
}

func AliCloudRouteEntryBasicDependence3(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
  		available_resource_creation = "VSwitch"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_nat_gateway" "default" {
  		vpc_id           = alicloud_vpc.default.id
  		nat_type         = "Enhanced"
  		vswitch_id       = alicloud_vswitch.default.id
		nat_gateway_name = var.name
	}
`, name)
}

func AliCloudRouteEntryBasicDependence4(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_efficiency"
  		available_resource_creation = "VSwitch"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_ecs_network_interface" "default" {
		network_interface_name = var.name
		vswitch_id             = alicloud_vswitch.default.id
  		security_group_ids     = alicloud_security_group.default.*.id
	}
`, name)
}

// AliCloudRouteEntryVbrDependence provisions the exact customer topology that
// originally surfaced Duplicated.VpcNextHop in the field: a VBR bound to a
// real physical connection, wired to a preserved test VPC via a pair of
// InitiatingSide/AcceptingSide router interfaces, so that a route entry can
// be declared on the VBR's own route table with nexthop_type=RouterInterface
// (the nexthop/route-table combination that returns Duplicated.VpcNextHop on
// a duplicate create, as opposed to InvalidCIDRBlock.Duplicate on a regular
// VPC route table). It relies on the "preserved-NODELETING" physical
// connection that CI keeps alive in cn-hangzhou specifically for VBR-based
// acceptance tests; TestAccAliCloudRouteEntry_stateGapRecoveryVbr forces the
// test region to cn-hangzhou via testAccPreCheckWithRegions so this
// datasource lookup resolves against that connection regardless of the
// account's default test region. The VPC side intentionally uses the existing
// default-NODELETING VPC instead of creating a temporary one because RI
// teardown can remain eventually consistent long enough to make VPC deletion
// flaky, while using a preserved VPC keeps the topology equivalent and the
// test cleanup bounded to the VBR/RI/route resources it creates.
func AliCloudRouteEntryVbrDependence(name string) string {
	vlanSeed := 0
	for _, r := range name {
		if r >= '0' && r <= '9' {
			vlanSeed = (vlanSeed*10 + int(r-'0')) % 1500
		}
	}
	vlanId := 500 + vlanSeed
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_express_connect_physical_connections" "default" {
	name_regex = "^preserved-NODELETING"
	status     = "Enabled"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING"
}

resource "alicloud_express_connect_virtual_border_router" "default" {
	local_gateway_ip           = "10.0.0.1"
	peer_gateway_ip            = "10.0.0.2"
	peering_subnet_mask        = "255.255.255.252"
	physical_connection_id     = data.alicloud_express_connect_physical_connections.default.ids.0
	virtual_border_router_name = var.name
	vlan_id                    = %d
}

// VBR-side router interface (InitiatingSide). The route entry under test
// will use this interface's id as nexthop_id on the VBR's route table.
resource "alicloud_express_connect_router_interface" "vbr" {
	router_type          = "VBR"
	router_id            = alicloud_express_connect_virtual_border_router.default.id
	opposite_region_id   = "cn-hangzhou"
	role                 = "InitiatingSide"
	spec                 = "Small.1"
	payment_type         = "PayAsYouGo"
	access_point_id      = data.alicloud_express_connect_physical_connections.default.connections.0.access_point_id
	opposite_router_type = "VRouter"

	lifecycle {
		ignore_changes = [opposite_interface_id, opposite_interface_owner_id, opposite_router_id, hc_rate]
	}
}

// VPC-side router interface (AcceptingSide), the counterpart the VBR side
// connects to.
resource "alicloud_express_connect_router_interface" "vpc" {
	router_type              = "VRouter"
	router_id                = data.alicloud_vpcs.default.vpcs.0.router_id
	opposite_region_id       = "cn-hangzhou"
	role                     = "AcceptingSide"
	spec                     = "Negative"
	payment_type             = "PayAsYouGo"
	opposite_access_point_id = data.alicloud_express_connect_physical_connections.default.connections.0.access_point_id
	opposite_router_type     = "VBR"
	depends_on               = [alicloud_express_connect_router_interface.vbr]

	lifecycle {
		ignore_changes = [opposite_interface_id, opposite_interface_owner_id, opposite_router_id, hc_rate]
	}
}

// Connections must be established accepting-side-first, then
// initiating-side, mirroring the working pattern used elsewhere in this
// provider's VBR acceptance tests.
resource "alicloud_router_interface_connection" "vpc_to_vbr" {
	interface_id          = alicloud_express_connect_router_interface.vpc.id
	opposite_interface_id = alicloud_express_connect_router_interface.vbr.id
	opposite_router_id    = alicloud_express_connect_virtual_border_router.default.id
	depends_on            = [alicloud_express_connect_router_interface.vpc]

	lifecycle {
		ignore_changes = [opposite_interface_id, opposite_interface_owner_id, opposite_router_id, opposite_router_type, interface_id]
	}
}

resource "alicloud_router_interface_connection" "vbr_to_vpc" {
	interface_id          = alicloud_express_connect_router_interface.vbr.id
	opposite_interface_id = alicloud_express_connect_router_interface.vpc.id
	opposite_router_id    = data.alicloud_vpcs.default.vpcs.0.router_id
	depends_on            = [alicloud_router_interface_connection.vpc_to_vbr]

	lifecycle {
		ignore_changes = [opposite_interface_id, opposite_interface_owner_id, opposite_router_id, opposite_router_type, interface_id]
	}
}
`, name, vlanId)
}

// TestAccAliCloudRouteEntry_stateGapRecoveryVbr precisely reproduces the
// customer's original field scenario, error code included: a route entry on
// a VBR's own route table with nexthop_type=RouterInterface. Unlike
// TestAccAliCloudRouteEntry_stateGapRecovery (which uses a NatGateway
// nexthop on a regular VPC route table and observes
// InvalidCIDRBlock.Duplicate), this test's PreConfig pre-creates the route
// directly against the VBR route table so Terraform's own create call
// observes the exact ErrorCode the customer hit: Duplicated.VpcNextHop.
//
// This requires a real, already-provisioned physical connection (VBRs
// cannot be created without one, and physical connections require carrier
// provisioning so they cannot be spun up by the test itself); it uses the
// "preserved-NODELETING" physical connection that CI keeps available in
// cn-hangzhou for exactly this purpose, and forces the test region there via
// testAccPreCheckWithRegions regardless of the account's default test
// region.
func TestAccAliCloudRouteEntry_stateGapRecoveryVbr(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_route_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudRouteEntryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc-RouteEntryVbr-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRouteEntryVbrDependence)

	const destinationCidrBlock = "10.4.0.0/16"
	var routeTableId, routerId, nextHopId string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
			// This test creates two new RouterInterfaces (VBR InitiatingSide +
			// VRouter AcceptingSide) in cn-hangzhou. RiPerUser defaults to 5
			// per region per user; if existing RIs from prior failed runs
			// already saturate the quota, CreateRouterInterface returns
			// QuotaExceeded before the route entry code path under test is
			// ever reached. Skip the test in that case rather than reporting
			// a false failure that does not exercise the fix. This is an
			// environmental gate, not a permanent skip: the test runs
			// normally whenever the quota can accommodate two new RIs.
			// Build the client via sharedClientForRegion instead of
			// testAccProvider.Meta(): PreCheck runs before any Step has
			// configured the provider, so Meta() is still nil here and a plain
			// type assertion panics. sharedClientForRegion builds the client
			// directly from the ALICLOUD_ACCESS_KEY / ALICLOUD_SECRET_KEY env
			// vars (same pattern as testAccPreCheckForCleanUpInstances), so
			// the quota gate actually runs instead of the test silently
			// skipping every time.
			rawClient, scErr := sharedClientForRegion("cn-hangzhou")
			if scErr != nil {
				t.Skipf("Skipping: cannot build shared client for cn-hangzhou in PreCheck: %v", scErr)
			}
			client := rawClient.(*connectivity.AliyunClient)
			listReq := vpc.CreateDescribeRouterInterfacesRequest()
			listReq.RegionId = "cn-hangzhou"
			if raw, qerr := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeRouterInterfaces(listReq)
			}); qerr == nil {
				if resp, ok := raw.(*vpc.DescribeRouterInterfacesResponse); ok && len(resp.RouterInterfaceSet.RouterInterfaceType) >= 4 {
					t.Skipf("Skipping: %d existing RouterInterfaces in cn-hangzhou would overflow RiPerUser quota when this test creates 2 more; clean up stale RIs to re-enable", len(resp.RouterInterfaceSet.RouterInterfaceType))
				}
			}
		},
		// IDRefreshName intentionally omitted: the duplicate step uses
		// ExpectError, so alicloud_route_entry.default is never created.
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Bring up the VBR + router interface topology only, skipping
				// testAccConfig's wrapper so that alicloud_route_entry itself is
				// not part of this step yet; the VBR's route table id and the
				// VBR-side router interface id are captured below for the
				// out-of-band API call in the next step.
				Config: AliCloudRouteEntryVbrDependence(name),
				// Destroy is left at its zero value (false) here on purpose: it
				// is a no-op for test execution, but its "Destroy:" prefix is
				// recognized by scripts/testing/testing_coverage_rate_check.go
				// as a step-end sentinel, letting it cleanly close this step's
				// Config (which has no inline attribute map to parse) instead
				// of folding the Check func below into the checker's
				// testAccCheck(map[string]string{...}) fallback parser, which
				// only understands that one specific shape.
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						vbrRs, ok := s.RootModule().Resources["alicloud_express_connect_virtual_border_router.default"]
						if !ok {
							return fmt.Errorf("not found: alicloud_express_connect_virtual_border_router.default")
						}
						routeTableId = vbrRs.Primary.Attributes["route_table_id"]
						if routeTableId == "" {
							return fmt.Errorf("alicloud_express_connect_virtual_border_router.default route_table_id is empty")
						}
						routerId = vbrRs.Primary.ID
						if routerId == "" {
							return fmt.Errorf("alicloud_express_connect_virtual_border_router.default id is empty")
						}

						riRs, ok := s.RootModule().Resources["alicloud_express_connect_router_interface.vbr"]
						if !ok {
							return fmt.Errorf("not found: alicloud_express_connect_router_interface.vbr")
						}
						nextHopId = riRs.Primary.ID
						if nextHopId == "" {
							return fmt.Errorf("alicloud_express_connect_router_interface.vbr id is empty")
						}
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					// Simulate the state/cloud divergence directly against the
					// API, bypassing Terraform entirely: the route is created
					// out-of-band on the VBR's own route table with the exact
					// same destination_cidrblock/nexthop_type/nexthop_id that the
					// Config below will declare, so Terraform's own
					// CreateRouteEntry call is guaranteed to observe a duplicate
					// nexthop on the cloud side while its own state has nothing
					// recorded for it yet. On a VBR route table this duplicate
					// surfaces as Duplicated.VpcNextHop, matching the customer's
					// original report verbatim.
					client := testAccProvider.Meta().(*connectivity.AliyunClient)
					request := vpc.CreateCreateRouteEntryRequest()
					request.RegionId = client.RegionId
					request.RouteTableId = routeTableId
					request.DestinationCidrBlock = destinationCidrBlock
					request.NextHopType = "RouterInterface"
					request.NextHopId = nextHopId
					if _, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
						return vpcClient.CreateRouteEntry(request)
					}); err != nil {
						t.Fatalf("failed to pre-create the route entry to simulate a state gap: %s", err)
					}

					vpcService := VpcService{client}
					routeEntryId := buildRouteEntryResourceId(routeTableId, routerId, destinationCidrBlock, "RouterInterface", nextHopId)
					if err := vpcService.WaitForRouteEntry(routeEntryId, Available, DefaultTimeout); err != nil {
						t.Fatalf("failed waiting for the pre-created VBR route entry to become available: %s", err)
					}

					if _, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
						return vpcClient.CreateRouteEntry(request)
					}); err == nil {
						t.Fatalf("expected Duplicated.VpcNextHop when creating the same VBR route entry twice, got nil")
					} else if !IsExpectedErrors(err, []string{"Duplicated.VpcNextHop"}) {
						t.Fatalf("expected Duplicated.VpcNextHop when creating the same VBR route entry twice, got: %s", err)
					}
				},
				Config: testAccConfig(map[string]interface{}{
					"route_table_id":        "${alicloud_express_connect_virtual_border_router.default.route_table_id}",
					"destination_cidrblock": destinationCidrBlock,
					"nexthop_type":          "RouterInterface",
					"nexthop_id":            "${alicloud_express_connect_router_interface.vbr.id}",
				}),
				// A duplicate means the route was created by a different
				// request, not by this apply; Create must surface an error
				// with an import hint instead of silently taking over state.
				ExpectError: regexp.MustCompile("already exists on the cloud.*Please import it using ID"),
			},
		},
	})
}
