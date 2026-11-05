package natgateway_test

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	natgateway "github.com/aliyun/terraform-provider-alicloud/alicloud/service/natgateway"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// The resource has a single optional-less configuration (nat_gateway_id and
// capacity_units are both required), so the basic and complete cases would use
// the same first-step config and are merged into one function.
func TestAccAliCloudNatGatewayCapacityReservation_basic(t *testing.T) {
	resourceId := "alicloud_nat_gateway_capacity_reservation.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "eu-central-1"
	}
	name := fmt.Sprintf("tf-testAcc%sNatGwCr%d", region, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckNatGatewayCapacityReservationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudNatGatewayCapacityReservationBasic(name, 9000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "nat_gateway_id"),
					resource.TestCheckResourceAttr(resourceId, "capacity_units", "9000"),
					resource.TestCheckResourceAttr(resourceId, "status", "provisioned"),
					// The API derives the guaranteed performance values from the
					// reserved CU, so they are read-only outputs of the resource.
					resource.TestCheckResourceAttrSet(resourceId, "bandwidth"),
					resource.TestCheckResourceAttrSet(resourceId, "cps"),
					resource.TestCheckResourceAttrSet(resourceId, "pps"),
					resource.TestCheckResourceAttrSet(resourceId, "session_count"),
					resource.TestCheckResourceAttrSet(resourceId, "decrease_requests_remaining"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Scaling up never consumes the daily decrease quota, so the
				// update step stays repeatable within one day.
				Config: testAccAliCloudNatGatewayCapacityReservationBasic(name, 9500),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "capacity_units", "9500"),
					resource.TestCheckResourceAttr(resourceId, "status", "provisioned"),
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

// TestUnitAliCloudNatGatewayCapacityReservation covers the Describe request
// shape and the response parsing without network access.
func TestUnitAliCloudNatGatewayCapacityReservation(t *testing.T) {
	rawClient := unitTestMockClient(t)
	var capturedQuery map[string]interface{}
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action, protocol, method, version, authType *string, query map[string]interface{}, body map[string]interface{}, runtime *util.RuntimeOptions) (map[string]interface{}, error) {
		capturedQuery = query
		return map[string]interface{}{
			"NatGatewayId":              "ngw-test",
			"CapacityReservationStatus": map[string]interface{}{"Status": "provisioned"},
			"DecreaseRequestsRemaining": "2",
			"LastModifiedTime":          "2026-07-02 10:07:51",
			"NatGatewayCapacity":        map[string]interface{}{"CapacityUnits": float64(9000)},
			"NatGatewayPerformanceValues": map[string]interface{}{
				"Bandwidth":    20,
				"Cps":          100000,
				"Pps":          3200000,
				"SessionCount": 2000000,
			},
		}, nil
	})
	defer patches.Reset()

	object, err := natgateway.DescribeNatGatewayCapacityReservation(context.Background(), rawClient, "ngw-test")
	if err != nil {
		t.Fatalf("describing capacity reservation failed: %s", err)
	}
	if got := fmt.Sprint(capturedQuery["Action"]); got != "DescribeNatGatewayCapacityReservation" {
		t.Fatalf("unexpected action: %s", got)
	}
	if got := fmt.Sprint(capturedQuery["NatGatewayId"]); got != "ngw-test" {
		t.Fatalf("unexpected NatGatewayId: %s", got)
	}
	capacity, ok := object["NatGatewayCapacity"].(map[string]interface{})
	if !ok {
		t.Fatalf("response lost the NatGatewayCapacity object: %v", object)
	}
	if got := fmt.Sprint(capacity["CapacityUnits"]); got != "9000" {
		t.Fatalf("unexpected CapacityUnits: %s", got)
	}
}

func testAccCheckNatGatewayCapacityReservationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_nat_gateway_capacity_reservation" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		object, err := natgateway.DescribeNatGatewayCapacityReservation(context.Background(), client, rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing capacity reservation of NAT gateway %s failed: %s", rs.Primary.ID, err)
		}
		status := ""
		if v, ok := object["CapacityReservationStatus"].(map[string]interface{}); ok {
			status, _ = v["Status"].(string)
		}
		if status == "" || status == "deleted" {
			continue
		}
		return fmt.Errorf("capacity reservation of NAT gateway %s still exists in status %q", rs.Primary.ID, status)
	}
	return nil
}

func testAccAliCloudNatGatewayCapacityReservationBasic(name string, capacityUnits int) string {
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

	resource "alicloud_nat_gateway" "default" {
		vpc_id           = alicloud_vpc.default.id
		nat_gateway_name = var.name
		nat_type         = "Enhanced"
		vswitch_id       = alicloud_vswitch.default.id
		eip_bind_mode    = "MULTI_BINDED"
	}

	resource "alicloud_nat_gateway_capacity_reservation" "default" {
		nat_gateway_id = alicloud_nat_gateway.default.id
		capacity_units = %d
	}
`, name, capacityUnits)
}
