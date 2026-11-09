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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/privatezone"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/vpc"
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

// sweepInboundEndpoints deletes every inbound endpoint whose Name matches a test prefix. IDs are
// collected from all pages first and deleted afterwards: DescribeInboundEndpoints pages by
// PageNumber/PageSize, so deleting while paging shifts later objects into an already-read page and
// silently skips them. List/parse failures abort with an error and delete failures are aggregated,
// so the sweeper never reports success for a partially failed cleanup.
func sweepInboundEndpoints(ctx context.Context, client inboundEndpointSweepAPI, prefixes []string) error {
	ids, contractErrs, err := collectInboundEndpointSweepIDs(ctx, client, prefixes)
	if err != nil {
		return err
	}
	deleteErrs := append([]error{}, contractErrs...)
	for _, id := range ids {
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
	return errors.Join(deleteErrs...)
}

// collectInboundEndpointSweepIDs pages through DescribeInboundEndpoints and returns the IDs of
// every endpoint whose Name matches a test prefix. The list response nests the array under
// $.InboundEndpoints.InboundEndpoint (unlike the outbound $.Endpoints flat array) and reports the
// total via TotalItems.
func collectInboundEndpointSweepIDs(ctx context.Context, client inboundEndpointSweepAPI, prefixes []string) ([]string, []error, error) {
	action := "DescribeInboundEndpoints"
	request := map[string]interface{}{
		"Lang":       "en",
		"PageSize":   helper.PageSizeLarge,
		"PageNumber": 1,
	}

	var ids []string
	var contractErrs []error
	seen := make(map[string]struct{})
	itemsSeen := 0
	for {
		response, err := helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
			return client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
		}, 5*time.Minute, helper.NeedRetry)
		logs.AddDebug(action, response, request)
		if err != nil {
			return nil, nil, fmt.Errorf("listing PrivateZone Inbound Endpoints: %w", err)
		}

		totalCount, err := convert.FormatInt(response["TotalItems"])
		if err != nil {
			return nil, nil, fmt.Errorf("parsing PrivateZone Inbound Endpoint TotalItems: %w", err)
		}
		raw, pathErr := jsonpath.Get("$.InboundEndpoints.InboundEndpoint", response)
		if totalCount != 0 && pathErr != nil {
			return nil, nil, fmt.Errorf("parsing PrivateZone Inbound Endpoint list: %w", pathErr)
		}
		result, ok := raw.([]interface{})
		if !ok {
			// A response that reports items must carry the array: an object, scalar or null at the
			// path is a contract violation, and treating it as an empty page would report cleanup
			// success while leaving endpoints behind.
			if totalCount != 0 {
				return nil, nil, fmt.Errorf("parsing PrivateZone Inbound Endpoint list: unexpected %T at $.InboundEndpoints.InboundEndpoint", raw)
			}
			result = nil
		}
		itemsSeen += len(result)
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				return nil, nil, fmt.Errorf("listing PrivateZone Inbound Endpoints: unexpected item type %T", v)
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
				contractErrs = append(contractErrs, fmt.Errorf("inbound endpoint %q has no string Id", name))
				continue
			}
			if _, duplicate := seen[id]; duplicate {
				continue
			}
			seen[id] = struct{}{}
			ids = append(ids, id)
		}

		// TotalItems is authoritative: when the total is an exact multiple of the page size a
		// further request would ask for an empty page that can look like a parse failure.
		if itemsSeen >= totalCount || len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return ids, contractErrs, nil
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
	// calls records the API order so a test can prove every list page is read before the
	// first delete; deleting while paging would shift the page window.
	calls       []string
	pageNumbers []int
}

func (f *fakeInboundEndpointSweepAPI) RpcPost(apiProductCode string, apiVersion string, apiName string, query map[string]interface{}, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
	switch apiName {
	case "DescribeInboundEndpoints":
		f.calls = append(f.calls, "list")
		if page, ok := body["PageNumber"].(int); ok {
			f.pageNumbers = append(f.pageNumbers, page)
		}
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
		f.calls = append(f.calls, "delete:"+id)
		if err, ok := f.deleteErrors[id]; ok {
			return nil, err
		}
		f.deleted = append(f.deleted, id)
		return map[string]interface{}{}, nil
	}
	return nil, fmt.Errorf("unexpected action %s", apiName)
}

func inboundEndpointSweepListResponse(items ...map[string]interface{}) map[string]interface{} {
	return inboundEndpointSweepListResponseWithTotal(len(items), items...)
}

// inboundEndpointSweepListResponseWithTotal models a page whose TotalItems is the total across all
// pages, which is what DescribeInboundEndpoints returns.
func inboundEndpointSweepListResponseWithTotal(total int, items ...map[string]interface{}) map[string]interface{} {
	values := make([]interface{}, 0, len(items))
	for _, item := range items {
		values = append(values, item)
	}
	return map[string]interface{}{
		"TotalItems": total,
		"InboundEndpoints": map[string]interface{}{
			"InboundEndpoint": values,
		},
	}
}

func TestUnitSweepInboundEndpointsReadsAllPagesBeforeDeleting(t *testing.T) {
	firstPage := make([]map[string]interface{}, 0, helper.PageSizeLarge)
	for i := 0; i < helper.PageSizeLarge; i++ {
		firstPage = append(firstPage, map[string]interface{}{
			"Id":   fmt.Sprintf("ep-page1-%d", i),
			"Name": fmt.Sprintf("tf-testAcc-p1-%d", i),
		})
	}
	secondPage := []map[string]interface{}{{"Id": "ep-page2", "Name": "tf-testAcc-p2"}}
	api := &fakeInboundEndpointSweepAPI{listResponses: []map[string]interface{}{
		inboundEndpointSweepListResponseWithTotal(helper.PageSizeLarge+1, firstPage...),
		inboundEndpointSweepListResponseWithTotal(helper.PageSizeLarge+1, secondPage...),
	}}

	if err := sweepInboundEndpoints(context.Background(), api, []string{"tf-testAcc"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if api.listCalls != 2 {
		t.Fatalf("expected two list pages, got %d", api.listCalls)
	}
	if len(api.deleted) != helper.PageSizeLarge+1 {
		t.Fatalf("expected every matching endpoint to be deleted, got %d", len(api.deleted))
	}
	if len(api.pageNumbers) != 2 || api.pageNumbers[0] != 1 || api.pageNumbers[1] != 2 {
		t.Fatalf("expected PageNumber 1 then 2, got %v", api.pageNumbers)
	}
	// Deleting while paging shifts later objects into an already-read page. Every list call must
	// finish before the first delete so the collected set cannot miss a shifted object.
	for i, call := range api.calls {
		if strings.HasPrefix(call, "delete:") {
			if i < 2 {
				t.Fatalf("delete %q happened before all list pages were read: %v", call, api.calls)
			}
			break
		}
	}
}

func TestUnitSweepInboundEndpointsStopsAtTotalItems(t *testing.T) {
	page := make([]map[string]interface{}, 0, helper.PageSizeLarge)
	for i := 0; i < helper.PageSizeLarge; i++ {
		page = append(page, map[string]interface{}{
			"Id":   fmt.Sprintf("ep-full-%d", i),
			"Name": fmt.Sprintf("tf-testAcc-full-%d", i),
		})
	}
	api := &fakeInboundEndpointSweepAPI{listResponses: []map[string]interface{}{
		inboundEndpointSweepListResponse(page...),
	}}

	if err := sweepInboundEndpoints(context.Background(), api, []string{"tf-testAcc"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if api.listCalls != 1 {
		t.Fatalf("a full page that already reaches TotalItems must not trigger another request, got %d list calls", api.listCalls)
	}
	if len(api.deleted) != helper.PageSizeLarge {
		t.Fatalf("expected every matching endpoint to be deleted, got %d", len(api.deleted))
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

func TestUnitSweepInboundEndpointsRejectsNonArrayList(t *testing.T) {
	cases := map[string]interface{}{
		"object": map[string]interface{}{"Id": "ep-object", "Name": "tf-testAcc-object"},
		"string": "tf-testAcc-string",
		"null":   nil,
	}
	for name, value := range cases {
		t.Run(name, func(t *testing.T) {
			api := &fakeInboundEndpointSweepAPI{listResponses: []map[string]interface{}{{
				"TotalItems": 1,
				"InboundEndpoints": map[string]interface{}{
					"InboundEndpoint": value,
				},
			}}}
			err := sweepInboundEndpoints(context.Background(), api, []string{"tf-testAcc"})
			if err == nil || !strings.Contains(err.Error(), "parsing PrivateZone Inbound Endpoint list") {
				t.Fatalf("a non-array path value must be reported as a parse error, got %v", err)
			}
			// A false success would skip the delete phase and leave endpoints behind.
			if api.listCalls != 1 || len(api.deleted) != 0 {
				t.Fatalf("a non-array path value must abort before any delete, list calls=%d deleted=%v", api.listCalls, api.deleted)
			}
		})
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
	deps := &inboundEndpointTestDependencies{}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPrivateZoneInboundEndpointDestroyWithDependencies(deps),
		Steps: []resource.TestStep{
			{
				// Create with 3 ip_configs (vswitches 0,1,2). ip is Optional+Computed
				// and service-assigned, omitted as in the outbound endpoint test.
				Config: testAccPvtzInboundEndpointConfig(name, []int{0, 1, 2}),
				Check: resource.ComposeTestCheckFunc(
					deps.capture,
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
			{
				// Every successful step in this new test file re-imports the resource, including
				// the recovery step that follows the two ExpectError steps.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
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
	deps := &inboundEndpointTestDependencies{}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPrivateZoneInboundEndpointDestroyWithDependencies(deps),
		Steps: []resource.TestStep{
			{
				// Create with two explicitly configured ips; Create must send them.
				Config: testAccPvtzInboundEndpointConfigWithIPs(name, []int{0, 1}, map[int]string{0: firstIp, 1: secondIp}),
				Check: resource.ComposeTestCheckFunc(
					deps.capture,
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
					deps.capture,
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
			{
				// Replacing a ForceNew field with ip omitted must let the service assign the new
				// addresses. Replaying the prior assigned ip would pin addresses the configuration
				// never asked for, and collide when create_before_destroy keeps the old endpoint.
				Config: testAccPvtzInboundEndpointConfig(name+"1", []int{0, 1}),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					deps.capture,
					resource.TestCheckResourceAttr(resourceId, "ip_configs.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "endpoint_name", name+"1"),
					resource.TestCheckResourceAttrSet(resourceId, "status"),
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

// inboundEndpointTestDependencies holds the VPC, VSwitch and security group IDs created by the
// test configuration. CheckDestroy only receives the post-destroy state, which no longer contains
// them, so they are captured while the state still has them.
type inboundEndpointTestDependencies struct {
	vpcIds           []string
	vswitchIds       []string
	securityGroupIds []string
}

func (d *inboundEndpointTestDependencies) capture(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "alicloud_vpc":
			d.vpcIds = appendUniqueID(d.vpcIds, rs.Primary.ID)
		case "alicloud_vswitch":
			d.vswitchIds = appendUniqueID(d.vswitchIds, rs.Primary.ID)
		case "alicloud_security_group":
			d.securityGroupIds = appendUniqueID(d.securityGroupIds, rs.Primary.ID)
		}
	}
	return nil
}

func appendUniqueID(ids []string, id string) []string {
	if id == "" {
		return ids
	}
	for _, existing := range ids {
		if existing == id {
			return ids
		}
	}
	return append(ids, id)
}

// verifyDestroyed confirms every dependency created by the test is gone. A failed terraform
// destroy already fails the test, but the cleanup evidence needs an authoritative per-object
// readback instead of inferring deletion from an empty workspace state.
func (d *inboundEndpointTestDependencies) verifyDestroyed(client *connectivity.AliyunClient) error {
	vpcConn := vpc.NewVpcService(client)
	for _, id := range d.vpcIds {
		if _, err := vpcConn.DescribeVpc(id); err == nil {
			return fmt.Errorf("test dependency VPC %s still exists", id)
		} else if !errs.IsProviderNotFound(err) {
			return fmt.Errorf("describing test dependency VPC %s failed: %s", id, err)
		}
		log.Printf("[INFO] Test dependency VPC deleted: %s", id)
	}
	for _, id := range d.vswitchIds {
		if _, err := vpcConn.DescribeVswitch(id); err == nil {
			return fmt.Errorf("test dependency VSwitch %s still exists", id)
		} else if !errs.IsProviderNotFound(err) {
			return fmt.Errorf("describing test dependency VSwitch %s failed: %s", id, err)
		}
		log.Printf("[INFO] Test dependency VSwitch deleted: %s", id)
	}
	ecsConn := ecs.NewService(client)
	for _, id := range d.securityGroupIds {
		if _, err := ecsConn.DescribeSecurityGroupAttribute(id); err == nil {
			return fmt.Errorf("test dependency security group %s still exists", id)
		} else if !errs.IsProviderNotFound(err) {
			return fmt.Errorf("describing test dependency security group %s failed: %s", id, err)
		}
		log.Printf("[INFO] Test dependency security group deleted: %s", id)
	}
	return nil
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

func testAccCheckPrivateZoneInboundEndpointDestroyWithDependencies(deps *inboundEndpointTestDependencies) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if err := testAccCheckPrivateZoneInboundEndpointDestroy(s); err != nil {
			return err
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		return deps.verifyDestroyed(client)
	}
}

// testAccPvtzInboundEndpointIpConfigsHCL renders ip_configs in block syntax
// (SetNestedBlock): ip_configs { ... }. The inbound endpoint uses a Framework
// SetNestedBlock (served at protocol v5 via fwadapt), which accepts block syntax,
// matching the SDKv2 outbound endpoint's TypeSet block.
func testAccPvtzInboundEndpointIpConfigsHCL(indices []int) string {
	configs := make([]string, len(indices))
	for i, idx := range indices {
		configs[i] = fmt.Sprintf(`  ip_configs {
    az_id      = alicloud_vswitch.default[%[1]d].zone_id
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
			fmt.Sprintf("    az_id      = alicloud_vswitch.default[%d].zone_id", idx),
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
  vpc_id              = alicloud_vpc.default.id
  security_group_name = var.name
}

resource "alicloud_security_group" "replacement" {
  vpc_id              = alicloud_vpc.default.id
  security_group_name = "${var.name}-replacement"
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
