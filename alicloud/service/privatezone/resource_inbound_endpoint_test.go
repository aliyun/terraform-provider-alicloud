package privatezone_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/privatezone"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
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

// inboundEndpointSweepAPI is the slice of the provider client the sweeper uses, so the error
// paths can be covered by unit tests without a cloud account.
type inboundEndpointSweepAPI interface {
	RpcPost(apiProductCode string, apiVersion string, apiName string, query map[string]interface{}, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error)
}

// inboundEndpointSweepPrefixes covers ACC names (tf-testAcc and legacy spellings) and the
// documentation example name (terraform-example), matching the naming/cleanup standard.
var inboundEndpointSweepPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf-testacc",
	"tf_testacc",
	"terraform-example",
}

func testSweepPrivateZoneInboundEndpoint(region string) error {
	client, err := privatezoneTestClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	return sweepInboundEndpoints(ctx, client, inboundEndpointSweepPrefixes)
}

// sweepInboundEndpoints lists every inbound endpoint via DescribeInboundEndpoints and deletes
// those whose Name matches a test prefix. The list response nests the array under
// $.InboundEndpoints.InboundEndpoint (unlike the outbound $.Endpoints flat array) and reports the
// total via TotalItems. List/parse failures abort with an error and delete failures are aggregated,
// so the sweeper never reports success for a partially failed cleanup.
func sweepInboundEndpoints(ctx context.Context, client inboundEndpointSweepAPI, prefixes []string) error {
	action := "DescribeInboundEndpoints"
	request := map[string]interface{}{
		"Lang":       "en",
		"PageSize":   helper.PageSizeLarge,
		"PageNumber": 1,
	}

	var deleteErrs []error
	for {
		response, err := helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
			return client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
		}, 5*time.Minute, helper.NeedRetry)
		logs.AddDebug(action, response, request)
		if err != nil {
			return fmt.Errorf("listing PrivateZone Inbound Endpoints: %w", err)
		}

		totalCount, err := convert.FormatInt(response["TotalItems"])
		if err != nil {
			return fmt.Errorf("parsing PrivateZone Inbound Endpoint TotalItems: %w", err)
		}
		raw, pathErr := jsonpath.Get("$.InboundEndpoints.InboundEndpoint", response)
		if totalCount != 0 && pathErr != nil {
			return fmt.Errorf("parsing PrivateZone Inbound Endpoint list: %w", pathErr)
		}
		result, _ := raw.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				return fmt.Errorf("listing PrivateZone Inbound Endpoints: unexpected item type %T", v)
			}
			name, ok := item["Name"].(string)
			if !ok {
				continue
			}
			if !inboundEndpointSweepMatch(name, prefixes) {
				log.Printf("[INFO] Skipping PrivateZone Inbound Endpoint: %s", name)
				continue
			}
			id, ok := item["Id"].(string)
			if !ok || id == "" {
				// DescribeInboundEndpoints declares Id as a string (unlike
				// DescribeInboundEndpoint, whose Id can arrive as a JSON number), so a
				// non-string or empty Id is a contract violation rather than a numeric ID.
				deleteErrs = append(deleteErrs, fmt.Errorf("inbound endpoint %q has no string Id", name))
				continue
			}
			deleteErr := sweepInboundEndpointDelete(ctx, client, id)
			if deleteErr != nil {
				if privatezone.IsInboundEndpointNotFound(deleteErr) {
					log.Printf("[INFO] PrivateZone Inbound Endpoint already removed: %s", id)
					continue
				}
				deleteErrs = append(deleteErrs, fmt.Errorf("deleting PrivateZone Inbound Endpoint %s: %w", id, deleteErr))
				continue
			}
			log.Printf("[INFO] Delete PrivateZone Inbound Endpoint success: %s", id)
		}

		if len(result) < helper.PageSizeLarge {
			break
		}
		page, _ := request["PageNumber"].(int)
		request["PageNumber"] = page + 1
	}
	return errors.Join(deleteErrs...)
}

func sweepInboundEndpointDelete(ctx context.Context, client inboundEndpointSweepAPI, id string) error {
	action := "DeleteInboundEndpoint"
	request := map[string]interface{}{
		"EndpointId":  id,
		"Lang":        "en",
		"ClientToken": helper.BuildClientToken(action),
	}
	_, err := helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
		return client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
	}, 5*time.Minute, helper.NeedRetry)
	return err
}

func inboundEndpointSweepMatch(name string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			return true
		}
	}
	return false
}

type fakeInboundEndpointSweepAPI struct {
	listResponses []map[string]interface{}
	listErrors    []error
	deleteErrors  map[string]error
	deleted       []string
	listCalls     int
}

func (f *fakeInboundEndpointSweepAPI) RpcPost(apiProductCode string, apiVersion string, apiName string, query map[string]interface{}, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
	switch apiName {
	case "DescribeInboundEndpoints":
		f.listCalls++
		if len(f.listErrors) >= f.listCalls {
			return nil, f.listErrors[f.listCalls-1]
		}
		if len(f.listResponses) >= f.listCalls {
			return f.listResponses[f.listCalls-1], nil
		}
		return map[string]interface{}{"TotalItems": 0}, nil
	case "DeleteInboundEndpoint":
		id, _ := body["EndpointId"].(string)
		if err, ok := f.deleteErrors[id]; ok {
			return nil, err
		}
		f.deleted = append(f.deleted, id)
		return map[string]interface{}{}, nil
	}
	return nil, fmt.Errorf("unexpected action %s", apiName)
}

func inboundEndpointSweepListResponse(items ...map[string]interface{}) map[string]interface{} {
	values := make([]interface{}, 0, len(items))
	for _, item := range items {
		values = append(values, item)
	}
	return map[string]interface{}{
		"TotalItems": len(items),
		"InboundEndpoints": map[string]interface{}{
			"InboundEndpoint": values,
		},
	}
}

func TestUnitSweepInboundEndpointsReturnsListError(t *testing.T) {
	api := &fakeInboundEndpointSweepAPI{listErrors: []error{fmt.Errorf("list failed")}}
	err := sweepInboundEndpoints(context.Background(), api, []string{"tf-testAcc"})
	if err == nil || !strings.Contains(err.Error(), "list failed") {
		t.Fatalf("expected the list error to be returned, got %v", err)
	}
}

func TestUnitSweepInboundEndpointsReturnsParseError(t *testing.T) {
	api := &fakeInboundEndpointSweepAPI{
		listResponses: []map[string]interface{}{{"TotalItems": 1}},
	}
	err := sweepInboundEndpoints(context.Background(), api, []string{"tf-testAcc"})
	if err == nil || !strings.Contains(err.Error(), "parsing PrivateZone Inbound Endpoint list") {
		t.Fatalf("expected the list parse error to be returned, got %v", err)
	}
}

func TestUnitSweepInboundEndpointsAggregatesDeleteErrors(t *testing.T) {
	api := &fakeInboundEndpointSweepAPI{
		listResponses: []map[string]interface{}{inboundEndpointSweepListResponse(
			map[string]interface{}{"Id": "ep-fail", "Name": "tf-testAcc-fail"},
			map[string]interface{}{"Id": "ep-ok", "Name": "tf-testAcc-ok"},
		)},
		deleteErrors: map[string]error{"ep-fail": fmt.Errorf("delete failed")},
	}
	err := sweepInboundEndpoints(context.Background(), api, []string{"tf-testAcc"})
	if err == nil || !strings.Contains(err.Error(), "delete failed") {
		t.Fatalf("expected the delete error to be returned, got %v", err)
	}
	if len(api.deleted) != 1 || api.deleted[0] != "ep-ok" {
		t.Fatalf("expected only ep-ok to be deleted, got %v", api.deleted)
	}
}

func TestUnitSweepInboundEndpointsIgnoresNotFoundDelete(t *testing.T) {
	api := &fakeInboundEndpointSweepAPI{
		listResponses: []map[string]interface{}{inboundEndpointSweepListResponse(
			map[string]interface{}{"Id": "ep-gone", "Name": "tf-testAcc-gone"},
		)},
		deleteErrors: map[string]error{
			"ep-gone": errs.WrapErrorf(&tea.SDKError{StatusCode: tea.Int(400), Code: tea.String("InboundEndpoint.NotExsits")}, "delete failed"),
		},
	}
	if err := sweepInboundEndpoints(context.Background(), api, []string{"tf-testAcc"}); err != nil {
		t.Fatalf("expected not-found delete to be ignored, got %v", err)
	}
}

func TestUnitIsInboundEndpointNotFound(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "typed 404",
			err:  errs.WrapErrorf(&tea.SDKError{StatusCode: tea.Int(404), Code: tea.String("404"), Data: tea.String("404")}, "describe failed"),
			want: true,
		},
		{
			name: "misspelled not-exists code without data",
			err:  errs.WrapErrorf(&tea.SDKError{StatusCode: tea.Int(400), Code: tea.String("InboundEndpoint.NotExsits")}, "describe failed"),
			want: true,
		},
		{
			name: "tea error without code or data",
			err:  errs.WrapErrorf(&tea.SDKError{}, "describe failed"),
			want: false,
		},
		{
			name: "unrelated error",
			err:  errs.WrapErrorf(fmt.Errorf("permission denied"), "describe failed"),
			want: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := privatezone.IsInboundEndpointNotFound(tc.err); got != tc.want {
				t.Fatalf("IsInboundEndpointNotFound() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestUnitSweepInboundEndpointsMatchesPrefixes(t *testing.T) {
	api := &fakeInboundEndpointSweepAPI{
		listResponses: []map[string]interface{}{inboundEndpointSweepListResponse(
			map[string]interface{}{"Id": "ep-test", "Name": "tf-testAcc-1"},
			map[string]interface{}{"Id": "ep-doc", "Name": "terraform-example-1"},
			map[string]interface{}{"Id": "ep-prod", "Name": "prod-endpoint"},
		)},
	}
	if err := sweepInboundEndpoints(context.Background(), api, inboundEndpointSweepPrefixes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(api.deleted) != 2 || api.deleted[0] != "ep-test" || api.deleted[1] != "ep-doc" {
		t.Fatalf("expected the ACC and documentation example endpoints to be deleted, got %v", api.deleted)
	}
}

func TestAccAliCloudPvtzInboundEndpoint_basic0(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, connectivity.TestPvtzRegions)
	resourceId := "alicloud_pvtz_inbound_endpoint.default"
	rand := tfacctest.RandIntRange(1000, 9999999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
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
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
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
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Update ip_configs: 3 -> 2 (vswitches 2,3).
				Config: testAccPvtzInboundEndpointConfig(name+"1", []int{2, 3}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "ip_configs.#", "2"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Restore endpoint_name; ip_configs unchanged.
				Config: testAccPvtzInboundEndpointConfig(name, []int{2, 3}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "endpoint_name", name),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// A single ip_config violates the SizeAtLeast(2) validator that
				// enforces the high-availability minimum documented for ip_configs.
				Config:      testAccPvtzInboundEndpointConfig(name, []int{0}),
				ExpectError: regexp.MustCompile(`(?i)at least 2`),
			},
			{
				// The contract requires IpConfig, and SetNestedBlock has no Required flag,
				// so the explicit validator must reject a configuration that omits the block.
				Config:      testAccPvtzInboundEndpointConfigWithoutIpConfigs(name),
				ExpectError: regexp.MustCompile(`(?i)ip_configs`),
			},
			{
				// The ExpectError steps above leave an invalid config as the last
				// written one. Restore a valid Config (matching the previous state)
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

func TestAccAliCloudPvtzInboundEndpoint_complete(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, connectivity.TestPvtzRegions)
	resourceId := "alicloud_pvtz_inbound_endpoint.default"
	rand := tfacctest.RandIntRange(1000, 9999999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	// The vswitches use cidrsubnet("172.16.0.0/12", 8, index), i.e. /20 blocks at
	// 172.16.0.0 and 172.16.16.0. Host 10 is inside the block and outside the addresses
	// Alibaba Cloud reserves, so it can be assigned explicitly.
	firstIp := "172.16.0.10"
	secondIp := "172.16.16.10"
	updatedFirstIp := "172.16.0.11"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPrivateZoneInboundEndpointDestroy,
		Steps: []resource.TestStep{
			{
				// Create with two explicitly configured ips; Create must send them.
				Config: testAccPvtzInboundEndpointConfigWithIPs(name, []int{0, 1}, map[int]string{0: firstIp, 1: secondIp}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "ip_configs.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "ip_configs.*", map[string]string{"ip": firstIp}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "ip_configs.*", map[string]string{"ip": secondIp}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Changing endpoint_name must not plan an ip change. The framework's post-apply
				// empty plan check fails if ModifyPlan replaces the explicit ips with unknowns.
				Config: testAccPvtzInboundEndpointConfigWithIPs(name+"1", []int{0, 1}, map[int]string{0: firstIp, 1: secondIp}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "endpoint_name", name+"1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "ip_configs.*", map[string]string{"ip": firstIp}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "ip_configs.*", map[string]string{"ip": secondIp}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Changing one explicit ip must reach UpdateInboundEndpoint. If Ip were
				// omitted, the service would assign a different address and this check fails.
				Config: testAccPvtzInboundEndpointConfigWithIPs(name+"1", []int{0, 1}, map[int]string{0: updatedFirstIp, 1: secondIp}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "ip_configs.*", map[string]string{"ip": updatedFirstIp}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "ip_configs.*", map[string]string{"ip": secondIp}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Removing ip from the configuration is a no-op for this Optional+Computed
				// attribute: the framework keeps the state value and the API assigns an IP
				// only when Ip is absent from the request. The prior value must remain.
				Config: testAccPvtzInboundEndpointConfigWithIPs(name+"1", []int{0, 1}, map[int]string{1: secondIp}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "ip_configs.*", map[string]string{"ip": updatedFirstIp}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "ip_configs.*", map[string]string{"ip": secondIp}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// security_group_id is ForceNew: replacing it must plan a replacement and
				// recreate the endpoint with the explicitly configured ips.
				Config: testAccPvtzInboundEndpointConfigWithIPsAndSecurityGroup(name+"1", []int{0, 1}, map[int]string{0: updatedFirstIp, 1: secondIp}, "alicloud_security_group.replacement.id"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "security_group_id"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "ip_configs.*", map[string]string{"ip": updatedFirstIp}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "ip_configs.*", map[string]string{"ip": secondIp}),
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

func testAccPvtzInboundEndpointIpConfigsHCLWithIPs(indices []int, ips map[int]string) string {
	configs := make([]string, len(indices))
	for i, idx := range indices {
		lines := []string{
			"  ip_configs {",
			fmt.Sprintf("    zone_id    = alicloud_vswitch.default[%d].zone_id", idx),
			fmt.Sprintf("    cidr_block = alicloud_vswitch.default[%d].cidr_block", idx),
			fmt.Sprintf("    vswitch_id = alicloud_vswitch.default[%d].id", idx),
		}
		// An empty value means the attribute is omitted from the configuration, not set to "":
		// a known empty string would defeat UseStateForUnknown and create a real diff.
		if ip := ips[idx]; ip != "" {
			lines = append(lines, fmt.Sprintf("    ip         = %q", ip))
		}
		lines = append(lines, "  }")
		configs[i] = strings.Join(lines, "\n")
	}
	return strings.Join(configs, "\n")
}

func testAccPvtzInboundEndpointConfig(name string, ipConfigsIndices []int) string {
	return testAccPvtzInboundEndpointConfigWithHCL(name, testAccPvtzInboundEndpointIpConfigsHCL(ipConfigsIndices))
}

func testAccPvtzInboundEndpointConfigWithoutIpConfigs(name string) string {
	return testAccPvtzInboundEndpointConfigWithHCL(name, "")
}

func testAccPvtzInboundEndpointConfigWithIPs(name string, ipConfigsIndices []int, ips map[int]string) string {
	return testAccPvtzInboundEndpointConfigWithHCL(name, testAccPvtzInboundEndpointIpConfigsHCLWithIPs(ipConfigsIndices, ips))
}

func testAccPvtzInboundEndpointConfigWithIPsAndSecurityGroup(name string, ipConfigsIndices []int, ips map[int]string, securityGroupRef string) string {
	return testAccPvtzInboundEndpointConfigWithHCLAndSecurityGroup(name, testAccPvtzInboundEndpointIpConfigsHCLWithIPs(ipConfigsIndices, ips), securityGroupRef)
}

func testAccPvtzInboundEndpointConfigWithHCL(name string, ipConfigsHCL string) string {
	return testAccPvtzInboundEndpointConfigWithHCLAndSecurityGroup(name, ipConfigsHCL, "alicloud_security_group.default.id")
}

func testAccPvtzInboundEndpointConfigWithHCLAndSecurityGroup(name string, ipConfigsHCL string, securityGroupRef string) string {
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

resource "alicloud_security_group" "replacement" {
  vpc_id = alicloud_vpc.default.id
  name   = "${var.name}-replacement"
}

resource "alicloud_vswitch" "default" {
  count      = 4
  vpc_id     = alicloud_vpc.default.id
  cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, count.index)
  zone_id    = data.alicloud_pvtz_resolver_zones.default.zones[count.index].zone_id
}

resource "alicloud_pvtz_inbound_endpoint" "default" {
  endpoint_name     = "%s"
  security_group_id = %s
  vpc_id            = alicloud_vpc.default.id
  vpc_region_id     = "%s"
  spec_version      = "standard"
  %s
}
`, name, name, securityGroupRef, acctest.DefaultRegionToTest, ipConfigsHCL)
}
