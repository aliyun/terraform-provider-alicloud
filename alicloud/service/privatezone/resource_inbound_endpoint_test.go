package privatezone_test

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/privatezone"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_pvtz_inbound_endpoint",
		&resource.Sweeper{
			Name: "alicloud_pvtz_inbound_endpoint",
			F:    testSweepPrivateZoneInboundEndpoint,
		})
}

// testSweepPrivateZoneInboundEndpoint lists every inbound endpoint via
// DescribeInboundEndpoints and deletes those whose Name matches the test
// prefix. The list response nests the array under
// $.InboundEndpoints.InboundEndpoint (unlike the outbound $.Endpoints flat
// array) and reports the total via TotalItems.
func testSweepPrivateZoneInboundEndpoint(region string) error {
	client, err := privatezoneTestClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}
	action := "DescribeInboundEndpoints"
	request := map[string]interface{}{}

	request["Lang"] = "en"
	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.InboundEndpoints.InboundEndpoint", response)
		totalCount, formatErr := convert.FormatInt(response["TotalItems"])
		if formatErr != nil {
			return formatErr
		}
		if totalCount != 0 && err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.InboundEndpoints.InboundEndpoint", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			if _, ok := item["Name"]; !ok {
				continue
			}
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["Name"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping PrivateZone Inbound Endpoint: %s", item["Name"].(string))
				continue
			}
			action := "DeleteInboundEndpoint"
			request := map[string]interface{}{
				"EndpointId": item["Id"],
			}
			request["ClientToken"] = helper.BuildClientToken(action)
			_, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
			if err != nil {
				log.Printf("[ERROR] Failed to delete PrivateZone Inbound Endpoint (%s): %s", item["Id"].(string), err)
			}
			log.Printf("[INFO] Delete PrivateZone Inbound Endpoint success: %s ", item["Id"].(string))
		}

		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAliCloudPvtzInboundEndpoint_basic0(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, connectivity.TestPvtzRegions)
	resourceId := "alicloud_pvtz_inbound_endpoint.default"
	rand := tfacctest.RandIntRange(1000, 9999999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPrivateZoneInboundEndpointDestroy,
		Steps: []resource.TestStep{
			{
				// Create with 3 ip_configs (vswitches 0,1,2). ip is Optional+Computed
				// and service-assigned, omitted as in the outbound endpoint test.
				Config: testAccPvtzInboundEndpointConfig(name, []int{0, 1, 2}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "endpoint_name", name),
					resource.TestCheckResourceAttr(resourceId, "spec_version", "standard"),
					resource.TestCheckResourceAttrSet(resourceId, "vpc_id"),
					resource.TestCheckResourceAttrSet(resourceId, "security_group_id"),
					resource.TestCheckResourceAttrSet(resourceId, "vpc_region_id"),
					resource.TestCheckResourceAttr(resourceId, "ip_configs.#", "3"),
					resource.TestCheckResourceAttrSet(resourceId, "status"),
					resource.TestCheckResourceAttrSet(resourceId, "endpoint_type"),
				),
			},
			{
				// Update endpoint_name only; ip_configs unchanged.
				Config: testAccPvtzInboundEndpointConfig(name+"1", []int{0, 1, 2}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "endpoint_name", name+"1"),
					resource.TestCheckResourceAttr(resourceId, "ip_configs.#", "3"),
				),
			},
			{
				// Update ip_configs: 3 -> 2 (vswitches 2,3).
				Config: testAccPvtzInboundEndpointConfig(name+"1", []int{2, 3}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "ip_configs.#", "2"),
				),
			},
			{
				// Restore endpoint_name; ip_configs unchanged.
				Config: testAccPvtzInboundEndpointConfig(name, []int{2, 3}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "endpoint_name", name),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				// A single ip_config violates the SizeAtLeast(2) validator that
				// enforces the high-availability minimum documented for ip_configs.
				Config:      testAccPvtzInboundEndpointConfig(name, []int{0}),
				ExpectError: regexp.MustCompile(`(?i)at least 2`),
			},
			{
				// The ExpectError step above leaves an invalid config as the last
				// written one. Restore a valid Config (matching the Step 4 state)
				// so the post-test destroy runs against a valid plan. This apply
				// is a no-op against the existing resource.
				Config: testAccPvtzInboundEndpointConfig(name, []int{2, 3}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "endpoint_name", name),
				),
			},
		},
	})
}

// testAccCheckPrivateZoneInboundEndpointDestroy accepts only the provider
// not-found that DescribeInboundEndpointContext constructs itself. A permission
// or other non-404 failure must surface as a cleanup failure instead of
// pretending the endpoint was deleted.
func testAccCheckPrivateZoneInboundEndpointDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	pvtzService := privatezone.NewService(client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_pvtz_inbound_endpoint" {
			continue
		}
		_, err := pvtzService.DescribeInboundEndpointContext(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("PrivateZone Inbound Endpoint %s still exists", rs.Primary.ID)
		}
		if errs.IsProviderNotFound(err) {
			continue
		}
		return fmt.Errorf("describing PrivateZone Inbound Endpoint %s failed: %s", rs.Primary.ID, err)
	}
	return nil
}

// testAccPvtzInboundEndpointIpConfigsHCL renders ip_configs in block syntax
// (SetNestedBlock): ip_configs { ... }. The inbound endpoint uses a Framework
// SetNestedBlock (served at protocol v5 via fwadapt), which accepts block syntax,
// matching the SDKv2 outbound endpoint's TypeSet block.
func testAccPvtzInboundEndpointIpConfigsHCL(indices []int) string {
	configs := make([]string, len(indices))
	for i, idx := range indices {
		configs[i] = fmt.Sprintf(`  ip_configs {
    zone_id    = alicloud_vswitch.default[%[1]d].zone_id
    cidr_block = alicloud_vswitch.default[%[1]d].cidr_block
    vswitch_id = alicloud_vswitch.default[%[1]d].id
  }`, idx)
	}
	return strings.Join(configs, "\n")
}

func testAccPvtzInboundEndpointConfig(name string, ipConfigsIndices []int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_pvtz_resolver_zones" "default" {
  status = "NORMAL"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_security_group" "default" {
  vpc_id      = alicloud_vpc.default.id
  name        = var.name
}

resource "alicloud_vswitch" "default" {
  count      = 4
  vpc_id     = alicloud_vpc.default.id
  cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, count.index)
  zone_id    = data.alicloud_pvtz_resolver_zones.default.zones[count.index].zone_id
}

resource "alicloud_pvtz_inbound_endpoint" "default" {
  endpoint_name     = "%s"
  security_group_id = alicloud_security_group.default.id
  vpc_id            = alicloud_vpc.default.id
  vpc_region_id     = "%s"
  spec_version      = "standard"
  %s
}
`, name, name, acctest.DefaultRegionToTest, testAccPvtzInboundEndpointIpConfigsHCL(ipConfigsIndices))
}

// lintignore: R001
