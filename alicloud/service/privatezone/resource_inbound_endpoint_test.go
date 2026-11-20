package privatezone_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea/tea"
	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/privatezone"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/vpc"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	fwvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestMain is the single package entry point: terraform-plugin-testing needs it to
// dispatch the sweeper registered in init() for this product package.
func TestMain(m *testing.M) {
	resource.TestMain(m)
}

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
				// Name is required by the list contract. Skipping a malformed item would treat a
				// partially parsed response as a non-target object and report cleanup success
				// while the endpoint is still there.
				contractErrs = append(contractErrs, fmt.Errorf("inbound endpoint %v has no string Name", item["Id"]))
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
		if itemsSeen >= totalCount {
			break
		}
		if len(result) < helper.PageSizeLarge {
			// A short page before the declared total means the response lost objects. Stopping
			// here would delete only part of the set and still report cleanup success.
			return nil, nil, fmt.Errorf("parsing PrivateZone Inbound Endpoint list: page %v returned %d items, but only %d of %d endpoints were listed", request["PageNumber"], len(result), itemsSeen, totalCount)
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

func TestUnitSweepInboundEndpointsRejectsShortPageBeforeTotal(t *testing.T) {
	api := &fakeInboundEndpointSweepAPI{listResponses: []map[string]interface{}{
		inboundEndpointSweepListResponseWithTotal(3,
			map[string]interface{}{"Id": "ep-short-1", "Name": "tf-testAcc-short-1"},
			map[string]interface{}{"Id": "ep-short-2", "Name": "tf-testAcc-short-2"},
		),
	}}
	err := sweepInboundEndpoints(context.Background(), api, []string{"tf-testAcc"})
	if err == nil || !strings.Contains(err.Error(), "only 2 of 3 endpoints were listed") {
		t.Fatalf("a short page before TotalItems must be a contract error, got %v", err)
	}
	if len(api.deleted) != 0 {
		t.Fatalf("a truncated page must abort before any delete, deleted=%v", api.deleted)
	}
}

func TestUnitSweepInboundEndpointsRejectsEmptyPageBeforeTotal(t *testing.T) {
	api := &fakeInboundEndpointSweepAPI{listResponses: []map[string]interface{}{
		{
			"TotalItems": 1,
			"InboundEndpoints": map[string]interface{}{
				"InboundEndpoint": []interface{}{},
			},
		},
	}}
	err := sweepInboundEndpoints(context.Background(), api, []string{"tf-testAcc"})
	if err == nil || !strings.Contains(err.Error(), "only 0 of 1 endpoints were listed") {
		t.Fatalf("an empty page before TotalItems must be a contract error, got %v", err)
	}
	if len(api.deleted) != 0 {
		t.Fatalf("an empty page must abort before any delete, deleted=%v", api.deleted)
	}
}

func TestUnitSweepInboundEndpointsRejectsMalformedName(t *testing.T) {
	cases := map[string]map[string]interface{}{
		"missing":    {"Id": "ep-missing-name"},
		"null":       {"Id": "ep-null-name", "Name": nil},
		"non string": {"Id": "ep-number-name", "Name": 123},
	}
	for name, item := range cases {
		t.Run(name, func(t *testing.T) {
			api := &fakeInboundEndpointSweepAPI{listResponses: []map[string]interface{}{
				inboundEndpointSweepListResponse(item),
			}}
			err := sweepInboundEndpoints(context.Background(), api, []string{"tf-testAcc"})
			if err == nil || !strings.Contains(err.Error(), "has no string Name") {
				t.Fatalf("a %s Name must be a contract error, got %v", name, err)
			}
			if len(api.deleted) != 0 {
				t.Fatalf("a malformed Name must not lead to a delete, deleted=%v", api.deleted)
			}
		})
	}
}

func TestUnitSweepInboundEndpointsReportsShortCodeDeleteFailure(t *testing.T) {
	api := &fakeInboundEndpointSweepAPI{
		listResponses: []map[string]interface{}{inboundEndpointSweepListResponse(
			map[string]interface{}{"Id": "ep-live", "Name": "tf-testAcc-live"},
		)},
		deleteErrors: map[string]error{
			"ep-live": errs.WrapErrorf(&tea.SDKError{StatusCode: tea.Int(400), Code: tea.String("InboundEndpoint")}, "delete failed"),
		},
	}
	err := sweepInboundEndpoints(context.Background(), api, []string{"tf-testAcc"})
	if err == nil || !strings.Contains(err.Error(), "delete failed") {
		t.Fatalf("a short endpoint code must not mark the endpoint as already removed, got %v", err)
	}
	if len(api.deleted) != 0 {
		t.Fatalf("a failed delete must not be recorded as deleted, deleted=%v", api.deleted)
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
			name: "correct not-exists code without data",
			err:  errs.WrapErrorf(&tea.SDKError{StatusCode: tea.Int(400), Code: tea.String("InboundEndpoint.NotExists")}, "describe failed"),
			want: true,
		},
		{
			name: "http 400 short endpoint code is not not-found",
			err:  errs.WrapErrorf(&tea.SDKError{StatusCode: tea.Int(400), Code: tea.String("InboundEndpoint")}, "describe failed"),
			want: false,
		},
		{
			name: "http 400 near-miss prefix is not not-found",
			err:  errs.WrapErrorf(&tea.SDKError{StatusCode: tea.Int(400), Code: tea.String("InboundEndpoint.Not")}, "describe failed"),
			want: false,
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
				// The minimum legal create omits both optional IP addresses. Each listener keeps
				// its configured identity and receives a non-empty address from the service.
				Config: testAccPvtzInboundEndpointConfig(name, []int{0, 1}),
				Check: resource.ComposeTestCheckFunc(
					deps.capture,
					resource.TestCheckResourceAttr(resourceId, "endpoint_name", name),
					resource.TestCheckResourceAttr(resourceId, "spec_version", "standard"),
					resource.TestCheckResourceAttrSet(resourceId, "vpc_id"),
					resource.TestCheckResourceAttrSet(resourceId, "security_group_id"),
					resource.TestCheckResourceAttrSet(resourceId, "vpc_region_id"),
					resource.TestCheckResourceAttr(resourceId, "ip_configs.#", "2"),
					testAccCheckPvtzInboundEndpointIpConfig(resourceId, 0),
					testAccCheckPvtzInboundEndpointIpConfig(resourceId, 1),
					resource.TestCheckResourceAttrSet(resourceId, "status"),
					resource.TestCheckResourceAttrSet(resourceId, "endpoint_type"),
					resource.TestCheckResourceAttrSet(resourceId, "vpc_name"),
					resource.TestCheckResourceAttrSet(resourceId, "create_time"),
					resource.TestCheckResourceAttrSet(resourceId, "update_time"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Update endpoint_name only; ip_configs unchanged.
				Config: testAccPvtzInboundEndpointConfig(name+"1", []int{0, 1}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "endpoint_name", name+"1"),
					resource.TestCheckResourceAttr(resourceId, "ip_configs.#", "2"),
					testAccCheckPvtzInboundEndpointIpConfig(resourceId, 0),
					testAccCheckPvtzInboundEndpointIpConfig(resourceId, 1),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Replace both listeners. Count alone would also pass if the old members remained,
				// so match every nested identity field against the two new VSwitch resources.
				Config: testAccPvtzInboundEndpointConfig(name+"1", []int{2, 3}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "ip_configs.#", "2"),
					testAccCheckPvtzInboundEndpointIpConfig(resourceId, 2),
					testAccCheckPvtzInboundEndpointIpConfig(resourceId, 3),
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
				// An explicitly empty address is not an omitted address. Schema validation
				// rejects it before any AddInboundEndpoint or UpdateInboundEndpoint call.
				Config:      testAccPvtzInboundEndpointConfigWithIPs(name, []int{0, 1}, map[int]string{0: ""}),
				ExpectError: regexp.MustCompile(`(?i)string length must be at least 1`),
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
				// the recovery step that follows the ExpectError steps.
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
					resource.TestCheckResourceAttrSet(resourceId, "vpc_name"),
					resource.TestCheckResourceAttrSet(resourceId, "create_time"),
					resource.TestCheckResourceAttrSet(resourceId, "update_time"),
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
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionUpdate),
					},
				},
				// New security group IDs make the dynamic collection wholly unknown during planning.
				// Apply makes it enumerable, updates the endpoint and retains both explicit addresses.
				Config: testAccPvtzInboundEndpointConfigWithUnknownDynamicIPs(name+"2", updatedFirstIp, secondIp),
				Check: resource.ComposeTestCheckFunc(
					deps.capture,
					resource.TestCheckResourceAttr(resourceId, "endpoint_name", name+"2"),
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
			{
				// Operation timeouts are configurable state, not API fields. Changing only this
				// block exercises Update's configured budget and must leave the endpoint unchanged.
				Config: testAccPvtzInboundEndpointConfigWithIPsAndTimeouts(name+"1", []int{0, 1}, map[int]string{}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "timeouts.create", "15m"),
					resource.TestCheckResourceAttr(resourceId, "timeouts.update", "10m"),
					resource.TestCheckResourceAttr(resourceId, "timeouts.delete", "15m"),
					resource.TestCheckResourceAttr(resourceId, "ip_configs.#", "2"),
				),
			},
			{
				// The API cannot return local operation timeouts; every cloud-backed attribute
				// still participates in ImportStateVerify.
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func testAccCheckPvtzInboundEndpointIpConfig(resourceId string, vswitchIndex int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok || rs.Primary == nil {
			return fmt.Errorf("%s is absent from state", resourceId)
		}
		vswitchId := fmt.Sprintf("alicloud_vswitch.default.%d", vswitchIndex)
		vs, ok := s.RootModule().Resources[vswitchId]
		if !ok || vs.Primary == nil {
			return fmt.Errorf("%s is absent from state", vswitchId)
		}
		expected := map[string]string{
			"az_id":      vs.Primary.Attributes["zone_id"],
			"cidr_block": vs.Primary.Attributes["cidr_block"],
			"vswitch_id": vs.Primary.ID,
		}
		for key, value := range rs.Primary.Attributes {
			if !strings.HasPrefix(key, "ip_configs.") || !strings.HasSuffix(key, ".vswitch_id") || value != vs.Primary.ID {
				continue
			}
			prefix := strings.TrimSuffix(key, "vswitch_id")
			for attribute, want := range expected {
				if got := rs.Primary.Attributes[prefix+attribute]; got != want || got == "" {
					return fmt.Errorf("%s%s = %q, want %q", resourceId, prefix+attribute, got, want)
				}
			}
			if ip := rs.Primary.Attributes[prefix+"ip"]; ip == "" {
				return fmt.Errorf("%s%sip is empty; the service must assign an omitted address", resourceId, prefix)
			}
			return nil
		}
		return fmt.Errorf("%s has no ip_configs member for %s (%s)", resourceId, vswitchId, vs.Primary.ID)
	}
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
		// Map presence distinguishes an omitted address from an explicitly empty string.
		if ip, ok := ips[idx]; ok {
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

func testAccPvtzInboundEndpointConfigWithUnknownDynamicIPs(name, firstIp, secondIp string) string {
	ipConfigsHCL := fmt.Sprintf(`  dynamic "ip_configs" {
    for_each = { for index, group in alicloud_security_group.listener_input : group.id => index }
    content {
      az_id      = alicloud_vswitch.default[ip_configs.value].zone_id
      cidr_block = alicloud_vswitch.default[ip_configs.value].cidr_block
      vswitch_id = alicloud_vswitch.default[ip_configs.value].id
      ip         = [%q, %q][ip_configs.value]
    }
  }`, firstIp, secondIp)
	return testAccPvtzInboundEndpointConfigWithHCL(name, ipConfigsHCL) + `
resource "alicloud_security_group" "listener_input" {
  count               = 2
  vpc_id              = alicloud_vpc.default.id
  security_group_name = "${var.name}-listener-${count.index}"
}
`
}

func testAccPvtzInboundEndpointConfigWithIPsAndTimeouts(name string, ipConfigsIndices []int, ips map[int]string) string {
	return testAccPvtzInboundEndpointConfigWithHCLAndSecurityGroupAndTimeouts(
		name,
		testAccPvtzInboundEndpointIpConfigsHCLWithIPs(ipConfigsIndices, ips),
		"alicloud_security_group.default.id",
		`timeouts {
    create = "15m"
    update = "10m"
    delete = "15m"
  }`,
	)
}

func testAccPvtzInboundEndpointConfigWithIPsAndSecurityGroup(name string, ipConfigsIndices []int, ips map[int]string, securityGroupRef string) string {
	return testAccPvtzInboundEndpointConfigWithHCLAndSecurityGroup(name, testAccPvtzInboundEndpointIpConfigsHCLWithIPs(ipConfigsIndices, ips), securityGroupRef)
}

func testAccPvtzInboundEndpointConfigWithHCL(name string, ipConfigsHCL string) string {
	return testAccPvtzInboundEndpointConfigWithHCLAndSecurityGroup(name, ipConfigsHCL, "alicloud_security_group.default.id")
}

func testAccPvtzInboundEndpointConfigWithHCLAndSecurityGroup(name string, ipConfigsHCL string, securityGroupRef string) string {
	return testAccPvtzInboundEndpointConfigWithHCLAndSecurityGroupAndTimeouts(name, ipConfigsHCL, securityGroupRef, "")
}

func testAccPvtzInboundEndpointConfigWithHCLAndSecurityGroupAndTimeouts(name string, ipConfigsHCL string, securityGroupRef string, timeoutsHCL string) string {
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
  %s
}
`, name, name, securityGroupRef, acctest.DefaultRegionToTest, ipConfigsHCL, timeoutsHCL)
}

// The unit coverage below drives the exported resource factory and the exported
// service, so it stays in the same external test package as the ACC cases. The
// unexported request, waiter and response-mapping helpers are exercised through
// Create, Update, Delete and Read instead of being reached directly.

func inboundEndpointUnitClient(t *testing.T, handler func(http.ResponseWriter, url.Values)) *connectivity.AliyunClient {
	t.Helper()
	return inboundEndpointUnitRequestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parsing request form: %s", err)
		}
		handler(w, r.Form)
	})
}

func inboundEndpointUnitRequestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	client, err := (&connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}).Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("pvtz", endpoint)
	return client
}

func inboundEndpointUnitResource(t *testing.T, client *connectivity.AliyunClient) fwresource.Resource {
	t.Helper()
	r := privatezone.NewInboundEndpointResource()
	configurable, ok := r.(fwresource.ResourceWithConfigure)
	if !ok {
		t.Fatalf("%T does not implement resource.ResourceWithConfigure", r)
	}
	var resp fwresource.ConfigureResponse
	configurable.Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: client}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("configuring the resource: %v", resp.Diagnostics)
	}
	return r
}

func inboundEndpointUnitSchema(t *testing.T, r fwresource.Resource) (fwschema.Schema, tftypes.Object) {
	t.Helper()
	ctx := context.Background()
	var schemaResp fwresource.SchemaResponse
	r.Schema(ctx, fwresource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("building schema: %v", schemaResp.Diagnostics)
	}
	objType, ok := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)
	if !ok {
		t.Fatalf("schema type is %T, want tftypes.Object", schemaResp.Schema.Type().TerraformType(ctx))
	}
	return schemaResp.Schema, objType
}

func inboundEndpointUnitValue(objType tftypes.Object, overrides map[string]tftypes.Value) tftypes.Value {
	values := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, attrType := range objType.AttributeTypes {
		values[name] = tftypes.NewValue(attrType, nil)
	}
	for name, value := range overrides {
		values[name] = value
	}
	return tftypes.NewValue(objType, values)
}

func inboundEndpointUnitString(value string) tftypes.Value {
	return tftypes.NewValue(tftypes.String, value)
}

func inboundEndpointUnitUnknownString() tftypes.Value {
	return tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
}

func inboundEndpointUnitTimeouts(t *testing.T, objType tftypes.Object, values map[string]string) tftypes.Value {
	t.Helper()
	timeoutsType, ok := objType.AttributeTypes["timeouts"].(tftypes.Object)
	if !ok {
		t.Fatalf("timeouts type is %T, want tftypes.Object", objType.AttributeTypes["timeouts"])
	}
	attrs := make(map[string]tftypes.Value, len(timeoutsType.AttributeTypes))
	for name, attrType := range timeoutsType.AttributeTypes {
		attrs[name] = tftypes.NewValue(attrType, nil)
	}
	for name, value := range values {
		attrs[name] = inboundEndpointUnitString(value)
	}
	return tftypes.NewValue(timeoutsType, attrs)
}

func TestUnitInboundEndpointIpRejectsOnlyExplicitEmptyString(t *testing.T) {
	r := inboundEndpointUnitResource(t, &connectivity.AliyunClient{})
	resourceSchema, _ := inboundEndpointUnitSchema(t, r)
	block, ok := resourceSchema.Blocks["ip_configs"].(fwschema.SetNestedBlock)
	if !ok {
		t.Fatalf("ip_configs has type %T", resourceSchema.Blocks["ip_configs"])
	}
	attribute, ok := block.NestedObject.Attributes["ip"].(fwschema.StringAttribute)
	if !ok {
		t.Fatalf("ip has type %T", block.NestedObject.Attributes["ip"])
	}
	for _, tc := range []struct {
		name      string
		value     types.String
		wantError bool
	}{
		{name: "omitted", value: types.StringNull()},
		{name: "unknown", value: types.StringUnknown()},
		{name: "configured", value: types.StringValue("172.16.0.10")},
		{name: "empty", value: types.StringValue(""), wantError: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			hasError := false
			for _, attributeValidator := range attribute.Validators {
				validatorResp := fwvalidator.StringResponse{}
				attributeValidator.ValidateString(context.Background(), fwvalidator.StringRequest{
					Path:        path.Root("ip_configs").AtName("ip"),
					ConfigValue: tc.value,
				}, &validatorResp)
				hasError = hasError || validatorResp.Diagnostics.HasError()
			}
			if hasError != tc.wantError {
				t.Fatalf("validation error = %v, want %v", hasError, tc.wantError)
			}
		})
	}
}

func TestUnitInboundEndpointTimeoutsSchema(t *testing.T) {
	r := inboundEndpointUnitResource(t, &connectivity.AliyunClient{})
	resourceSchema, _ := inboundEndpointUnitSchema(t, r)
	block, ok := resourceSchema.Blocks["timeouts"].(fwschema.SingleNestedBlock)
	if !ok {
		t.Fatalf("timeouts block has type %T", resourceSchema.Blocks["timeouts"])
	}
	if _, ok := block.CustomType.(timeouts.Type); !ok {
		t.Fatalf("timeouts CustomType = %T, want timeouts.Type", block.CustomType)
	}
	if len(block.Attributes) != 3 {
		t.Fatalf("timeouts has %d attributes, want 3", len(block.Attributes))
	}
	for _, operation := range []string{"create", "update", "delete"} {
		attribute, ok := block.Attributes[operation].(fwschema.StringAttribute)
		if !ok {
			t.Fatalf("timeouts.%s has type %T", operation, block.Attributes[operation])
		}
		for _, tc := range []struct {
			value     string
			wantError bool
		}{
			{value: "10m"},
			{value: "soon", wantError: true},
			{value: "0s"},
			{value: "-1m"},
		} {
			t.Run(operation+"/"+tc.value, func(t *testing.T) {
				hasError := false
				for _, attributeValidator := range attribute.Validators {
					validatorResp := fwvalidator.StringResponse{}
					attributeValidator.ValidateString(context.Background(), fwvalidator.StringRequest{
						Path:        path.Root("timeouts").AtName(operation),
						ConfigValue: types.StringValue(tc.value),
					}, &validatorResp)
					hasError = hasError || validatorResp.Diagnostics.HasError()
				}
				if hasError != tc.wantError {
					t.Fatalf("validation error = %v, want %v", hasError, tc.wantError)
				}
			})
		}
	}
}

func TestUnitInboundEndpointTimeoutsStateCompatibility(t *testing.T) {
	ctx := context.Background()
	r := inboundEndpointUnitResource(t, &connectivity.AliyunClient{})
	resourceSchema, objType := inboundEndpointUnitSchema(t, r)
	decode := func(t *testing.T, rawTimeouts tftypes.Value) timeouts.Value {
		t.Helper()
		overrides := inboundEndpointUnitOverrides("tf-testAcc-unit", tftypes.NewValue(objType.AttributeTypes["ip_configs"], nil))
		overrides["id"] = inboundEndpointUnitString("ep-1")
		overrides["timeouts"] = rawTimeouts
		state := tfsdk.State{Schema: resourceSchema, Raw: inboundEndpointUnitValue(objType, overrides)}
		var value timeouts.Value
		if diagnostics := state.GetAttribute(ctx, path.Root("timeouts"), &value); diagnostics.HasError() {
			t.Fatalf("decoding timeouts: %v", diagnostics)
		}
		return value
	}

	t.Run("null uses ten minute defaults", func(t *testing.T) {
		value := decode(t, tftypes.NewValue(objType.AttributeTypes["timeouts"], nil))
		if !value.IsNull() {
			t.Fatal("null timeouts did not remain null")
		}
		for operation, accessor := range map[string]func(context.Context, time.Duration) (time.Duration, diag.Diagnostics){
			"create": value.Create,
			"update": value.Update,
			"delete": value.Delete,
		} {
			got, diagnostics := accessor(ctx, 10*time.Minute)
			if diagnostics.HasError() || got != 10*time.Minute {
				t.Fatalf("%s timeout = %s, diagnostics = %v", operation, got, diagnostics)
			}
		}
	})
	t.Run("configured values survive decoding", func(t *testing.T) {
		value := decode(t, inboundEndpointUnitTimeouts(t, objType, map[string]string{
			"create": "11m", "update": "12m", "delete": "13m",
		}))
		for operation, tc := range map[string]struct {
			accessor func(context.Context, time.Duration) (time.Duration, diag.Diagnostics)
			want     time.Duration
		}{
			"create": {value.Create, 11 * time.Minute},
			"update": {value.Update, 12 * time.Minute},
			"delete": {value.Delete, 13 * time.Minute},
		} {
			got, diagnostics := tc.accessor(ctx, 10*time.Minute)
			if diagnostics.HasError() || got != tc.want {
				t.Fatalf("%s timeout = %s, want %s, diagnostics = %v", operation, got, tc.want, diagnostics)
			}
		}
	})
}

func inboundEndpointUnitIpConfigs(t *testing.T, objType tftypes.Object, items []map[string]tftypes.Value) tftypes.Value {
	t.Helper()
	setType, ok := objType.AttributeTypes["ip_configs"].(tftypes.Set)
	if !ok {
		t.Fatalf("ip_configs type is %T, want tftypes.Set", objType.AttributeTypes["ip_configs"])
	}
	elemType, ok := setType.ElementType.(tftypes.Object)
	if !ok {
		t.Fatalf("ip_configs element type is %T, want tftypes.Object", setType.ElementType)
	}
	elems := make([]tftypes.Value, 0, len(items))
	for _, item := range items {
		attrs := make(map[string]tftypes.Value, len(elemType.AttributeTypes))
		for name, attrType := range elemType.AttributeTypes {
			attrs[name] = tftypes.NewValue(attrType, nil)
		}
		for name, value := range item {
			attrs[name] = value
		}
		elems = append(elems, tftypes.NewValue(elemType, attrs))
	}
	return tftypes.NewValue(setType, elems)
}

func inboundEndpointUnitOverrides(name string, ipConfigs tftypes.Value) map[string]tftypes.Value {
	return map[string]tftypes.Value{
		"endpoint_name":     inboundEndpointUnitString(name),
		"security_group_id": inboundEndpointUnitString("sg-1"),
		"vpc_id":            inboundEndpointUnitString("vpc-1"),
		"vpc_region_id":     inboundEndpointUnitString("cn-hangzhou"),
		"spec_version":      inboundEndpointUnitString("standard"),
		"ip_configs":        ipConfigs,
	}
}

func inboundEndpointUnitRespond(t *testing.T, w http.ResponseWriter, status int, body map[string]interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		t.Errorf("encoding response: %s", err)
	}
}

func inboundEndpointUnitDescribeBody(id, name, status string) map[string]interface{} {
	return map[string]interface{}{
		"RequestId":       "req-describe",
		"Id":              id,
		"Name":            name,
		"Status":          status,
		"SecurityGroupId": "sg-1",
		"VpcId":           "vpc-1",
		"VpcRegionId":     "cn-hangzhou",
		"SpecVersion":     "standard",
		"EndpointType":    "inbound",
		"VpcName":         "vpc-unit",
		"CreateTime":      "2026-11-18T06:00:00Z",
		"UpdateTime":      "2026-11-18T06:05:00Z",
		"IpConfigs": []interface{}{
			map[string]interface{}{"AzId": "cn-hangzhou-a", "CidrBlock": "172.16.0.0/20", "VSwitchId": "vsw-explicit", "Ip": "172.16.0.10"},
			map[string]interface{}{"AzId": "cn-hangzhou-b", "CidrBlock": "172.16.16.0/20", "VSwitchId": "vsw-auto", "Ip": "172.16.0.11"},
		},
	}
}

// inboundEndpointUnitIpConfigForms groups IpConfig.<n>.<field> form entries by AzId so
// the assertions do not depend on how the Framework orders a set.
func inboundEndpointUnitIpConfigForms(form url.Values) map[string]map[string]string {
	byIndex := map[string]map[string]string{}
	for key, values := range form {
		if !strings.HasPrefix(key, "IpConfig.") || len(values) == 0 {
			continue
		}
		parts := strings.Split(key, ".")
		if len(parts) != 3 {
			continue
		}
		entry, ok := byIndex[parts[1]]
		if !ok {
			entry = map[string]string{}
			byIndex[parts[1]] = entry
		}
		entry[parts[2]] = values[0]
	}
	byAzId := make(map[string]map[string]string, len(byIndex))
	for _, entry := range byIndex {
		if azId := entry["AzId"]; azId != "" {
			byAzId[azId] = entry
		}
	}
	return byAzId
}

func inboundEndpointUnitStateValue(t *testing.T, resourceSchema fwschema.Schema, objType tftypes.Object, name string, ipConfigs tftypes.Value) (fwschema.Schema, tftypes.Value) {
	t.Helper()
	overrides := inboundEndpointUnitOverrides(name, ipConfigs)
	overrides["id"] = inboundEndpointUnitString("ep-1")
	return resourceSchema, inboundEndpointUnitValue(objType, overrides)
}

func TestUnitInboundEndpointModifyPlanPreservesUnknownIpConfigs(t *testing.T) {
	ctx := context.Background()
	r := inboundEndpointUnitResource(t, &connectivity.AliyunClient{})
	resourceSchema, objType := inboundEndpointUnitSchema(t, r)
	modifyPlan, ok := r.(fwresource.ResourceWithModifyPlan)
	if !ok {
		t.Fatalf("%T does not implement resource.ResourceWithModifyPlan", r)
	}

	stateIpConfigs := inboundEndpointUnitIpConfigs(t, objType, []map[string]tftypes.Value{
		{
			"az_id":      inboundEndpointUnitString("cn-hangzhou-a"),
			"cidr_block": inboundEndpointUnitString("172.16.0.0/20"),
			"vswitch_id": inboundEndpointUnitString("vsw-1"),
			"ip":         inboundEndpointUnitString("172.16.0.10"),
		},
		{
			"az_id":      inboundEndpointUnitString("cn-hangzhou-b"),
			"cidr_block": inboundEndpointUnitString("172.16.16.0/20"),
			"vswitch_id": inboundEndpointUnitString("vsw-2"),
			"ip":         inboundEndpointUnitString("172.16.16.10"),
		},
	})
	_, stateValue := inboundEndpointUnitStateValue(t, resourceSchema, objType, "tf-testAcc-old", stateIpConfigs)
	unknownIpConfigs := tftypes.NewValue(objType.AttributeTypes["ip_configs"], tftypes.UnknownValue)
	configValue := inboundEndpointUnitValue(objType, inboundEndpointUnitOverrides("tf-testAcc-new", unknownIpConfigs))
	planOverrides := inboundEndpointUnitOverrides("tf-testAcc-new", unknownIpConfigs)
	planOverrides["id"] = inboundEndpointUnitString("ep-1")
	planValue := inboundEndpointUnitValue(objType, planOverrides)
	plan := tfsdk.Plan{Schema: resourceSchema, Raw: planValue}
	resp := fwresource.ModifyPlanResponse{Plan: plan}

	modifyPlan.ModifyPlan(ctx, fwresource.ModifyPlanRequest{
		State:  tfsdk.State{Schema: resourceSchema, Raw: stateValue},
		Config: tfsdk.Config{Schema: resourceSchema, Raw: configValue},
		Plan:   plan,
	}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("ModifyPlan returned diagnostics: %v", resp.Diagnostics)
	}
	if !resp.Plan.Raw.Equal(planValue) {
		t.Fatalf("ModifyPlan changed an entirely unknown ip_configs collection: got %s, want %s", resp.Plan.Raw, planValue)
	}
}

// TestUnitInboundEndpointCreateSendsOnlyConfiguredIps proves on the wire that a configured
// ip reaches AddInboundEndpoint while an omitted ip stays absent for service assignment.
func TestUnitInboundEndpointCreateSendsOnlyConfiguredIps(t *testing.T) {
	ctx := context.Background()
	var mu sync.Mutex
	var createForms []url.Values

	client := inboundEndpointUnitClient(t, func(w http.ResponseWriter, form url.Values) {
		switch form.Get("Action") {
		case "AddInboundEndpoint":
			mu.Lock()
			createForms = append(createForms, form)
			mu.Unlock()
			inboundEndpointUnitRespond(t, w, http.StatusOK, map[string]interface{}{"RequestId": "req-create", "EndpointId": "ep-1"})
		case "DescribeInboundEndpoint":
			inboundEndpointUnitRespond(t, w, http.StatusOK, inboundEndpointUnitDescribeBody("ep-1", "tf-testAcc-unit", "SUCCESS"))
		default:
			t.Errorf("unexpected action %q", form.Get("Action"))
			inboundEndpointUnitRespond(t, w, http.StatusOK, map[string]interface{}{"RequestId": "req-unexpected"})
		}
	})
	r := inboundEndpointUnitResource(t, client)
	resourceSchema, objType := inboundEndpointUnitSchema(t, r)

	ipConfigs := inboundEndpointUnitIpConfigs(t, objType, []map[string]tftypes.Value{
		{
			"az_id":      inboundEndpointUnitString("cn-hangzhou-a"),
			"cidr_block": inboundEndpointUnitString("172.16.0.0/20"),
			"vswitch_id": inboundEndpointUnitString("vsw-explicit"),
			"ip":         inboundEndpointUnitString("172.16.0.10"),
		},
		{
			"az_id":      inboundEndpointUnitString("cn-hangzhou-b"),
			"cidr_block": inboundEndpointUnitString("172.16.16.0/20"),
			"vswitch_id": inboundEndpointUnitString("vsw-auto"),
		},
	})
	overrides := inboundEndpointUnitOverrides("tf-testAcc-unit", ipConfigs)
	configValue := inboundEndpointUnitValue(objType, overrides)
	planOverrides := make(map[string]tftypes.Value, len(overrides)+1)
	for name, value := range overrides {
		planOverrides[name] = value
	}
	planOverrides["id"] = inboundEndpointUnitUnknownString()
	planValue := inboundEndpointUnitValue(objType, planOverrides)

	req := fwresource.CreateRequest{
		Plan:   tfsdk.Plan{Schema: resourceSchema, Raw: planValue},
		Config: tfsdk.Config{Schema: resourceSchema, Raw: configValue},
	}
	resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}
	r.Create(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create returned diagnostics: %v", resp.Diagnostics)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(createForms) != 1 {
		t.Fatalf("expected exactly one AddInboundEndpoint call, got %d", len(createForms))
	}
	form := createForms[0]
	if got := form.Get("Name"); got != "tf-testAcc-unit" {
		t.Fatalf("expected the endpoint name on the wire, got %q", got)
	}
	if got := form.Get("ClientToken"); got == "" {
		t.Fatal("the create request must carry a ClientToken")
	}
	byAz := inboundEndpointUnitIpConfigForms(form)
	if got := byAz["cn-hangzhou-a"]["Ip"]; got != "172.16.0.10" {
		t.Fatalf("explicit ip must be sent, got %q", got)
	}
	if got := byAz["cn-hangzhou-a"]["VSwitchId"]; got != "vsw-explicit" {
		t.Fatalf("vswitch id must be sent, got %q", got)
	}
	if got := byAz["cn-hangzhou-a"]["CidrBlock"]; got != "172.16.0.0/20" {
		t.Fatalf("cidr block must be sent, got %q", got)
	}
	if _, ok := byAz["cn-hangzhou-b"]["Ip"]; ok {
		t.Fatal("an omitted ip must stay absent from the create request")
	}
	if got := byAz["cn-hangzhou-b"]["VSwitchId"]; got != "vsw-auto" {
		t.Fatalf("vswitch id must still be sent for an auto-assigned ip, got %q", got)
	}
}

// TestUnitInboundEndpointUpdateWireSendsExplicitIp proves on the wire that UpdateInboundEndpoint
// carries IpConfig.N.Ip for a configured ip and omits it for a service-assigned one.
func TestUnitInboundEndpointUpdateWireSendsExplicitIp(t *testing.T) {
	ctx := context.Background()
	var mu sync.Mutex
	var updateForms []url.Values

	client := inboundEndpointUnitClient(t, func(w http.ResponseWriter, form url.Values) {
		switch form.Get("Action") {
		case "UpdateInboundEndpoint":
			mu.Lock()
			updateForms = append(updateForms, form)
			mu.Unlock()
			inboundEndpointUnitRespond(t, w, http.StatusOK, map[string]interface{}{"RequestId": "req-update"})
		case "DescribeInboundEndpoint":
			inboundEndpointUnitRespond(t, w, http.StatusOK, inboundEndpointUnitDescribeBody("ep-1", "tf-testAcc-unit", "SUCCESS"))
		default:
			t.Errorf("unexpected action %q", form.Get("Action"))
			inboundEndpointUnitRespond(t, w, http.StatusOK, map[string]interface{}{"RequestId": "req-unexpected"})
		}
	})
	r := inboundEndpointUnitResource(t, client)
	resourceSchema, objType := inboundEndpointUnitSchema(t, r)

	stateIpConfigs := inboundEndpointUnitIpConfigs(t, objType, []map[string]tftypes.Value{
		{
			"az_id":      inboundEndpointUnitString("cn-hangzhou-a"),
			"cidr_block": inboundEndpointUnitString("172.16.0.0/20"),
			"vswitch_id": inboundEndpointUnitString("vsw-1"),
			"ip":         inboundEndpointUnitString("172.16.0.10"),
		},
		{
			"az_id":      inboundEndpointUnitString("cn-hangzhou-b"),
			"cidr_block": inboundEndpointUnitString("172.16.16.0/20"),
			"vswitch_id": inboundEndpointUnitString("vsw-2"),
		},
	})
	planIpConfigs := inboundEndpointUnitIpConfigs(t, objType, []map[string]tftypes.Value{
		{
			"az_id":      inboundEndpointUnitString("cn-hangzhou-a"),
			"cidr_block": inboundEndpointUnitString("172.16.0.0/20"),
			"vswitch_id": inboundEndpointUnitString("vsw-1"),
			"ip":         inboundEndpointUnitString("172.16.0.11"),
		},
		{
			"az_id":      inboundEndpointUnitString("cn-hangzhou-b"),
			"cidr_block": inboundEndpointUnitString("172.16.16.0/20"),
			"vswitch_id": inboundEndpointUnitString("vsw-2"),
			"ip":         inboundEndpointUnitUnknownString(),
		},
	})

	_, stateValue := inboundEndpointUnitStateValue(t, resourceSchema, objType, "tf-testAcc-old", stateIpConfigs)
	planOverrides := inboundEndpointUnitOverrides("tf-testAcc-unit", planIpConfigs)
	planOverrides["id"] = inboundEndpointUnitString("ep-1")
	planValue := inboundEndpointUnitValue(objType, planOverrides)

	req := fwresource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: resourceSchema, Raw: planValue},
		State: tfsdk.State{Schema: resourceSchema, Raw: stateValue},
	}
	resp := &fwresource.UpdateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	r.Update(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update returned diagnostics: %v", resp.Diagnostics)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(updateForms) != 1 {
		t.Fatalf("expected exactly one UpdateInboundEndpoint call, got %d", len(updateForms))
	}
	form := updateForms[0]
	if got := form.Get("EndpointId"); got != "ep-1" {
		t.Fatalf("expected EndpointId ep-1 on the wire, got %q", got)
	}
	if got := form.Get("Name"); got != "tf-testAcc-unit" {
		t.Fatalf("expected the changed name on the wire, got %q", got)
	}
	byAz := inboundEndpointUnitIpConfigForms(form)
	if got := byAz["cn-hangzhou-a"]["Ip"]; got != "172.16.0.11" {
		t.Fatalf("expected the explicit ip on the wire, got %q", got)
	}
	if _, ok := byAz["cn-hangzhou-b"]["Ip"]; ok {
		t.Fatal("a service-assigned ip must be omitted from the update request")
	}
}

// TestUnitInboundEndpointCreateRetriesReuseClientToken records every AddInboundEndpoint attempt:
// the first is throttled and the retry must carry the same ClientToken, so one logical create
// cannot turn into two endpoints.
func TestUnitInboundEndpointCreateRetriesReuseClientToken(t *testing.T) {
	ctx := context.Background()
	var mu sync.Mutex
	var tokens []string

	client := inboundEndpointUnitClient(t, func(w http.ResponseWriter, form url.Values) {
		switch form.Get("Action") {
		case "AddInboundEndpoint":
			mu.Lock()
			tokens = append(tokens, form.Get("ClientToken"))
			attempt := len(tokens)
			mu.Unlock()
			if attempt == 1 {
				inboundEndpointUnitRespond(t, w, http.StatusInternalServerError, map[string]interface{}{
					"Code": "Throttling.User", "Message": "code: 500, request throttled", "RequestId": "req-1",
				})
				return
			}
			inboundEndpointUnitRespond(t, w, http.StatusOK, map[string]interface{}{"EndpointId": "ep-1", "RequestId": "req-2"})
		case "DescribeInboundEndpoint":
			inboundEndpointUnitRespond(t, w, http.StatusOK, inboundEndpointUnitDescribeBody("ep-1", "tf-testAcc-unit", "SUCCESS"))
		default:
			t.Errorf("unexpected action %q", form.Get("Action"))
			inboundEndpointUnitRespond(t, w, http.StatusOK, map[string]interface{}{"RequestId": "req-unexpected"})
		}
	})
	r := inboundEndpointUnitResource(t, client)
	resourceSchema, objType := inboundEndpointUnitSchema(t, r)

	ipConfigs := inboundEndpointUnitIpConfigs(t, objType, []map[string]tftypes.Value{
		{
			"az_id":      inboundEndpointUnitString("cn-hangzhou-a"),
			"cidr_block": inboundEndpointUnitString("172.16.0.0/20"),
			"vswitch_id": inboundEndpointUnitString("vsw-explicit"),
			"ip":         inboundEndpointUnitString("172.16.0.10"),
		},
		{
			"az_id":      inboundEndpointUnitString("cn-hangzhou-b"),
			"cidr_block": inboundEndpointUnitString("172.16.16.0/20"),
			"vswitch_id": inboundEndpointUnitString("vsw-auto"),
		},
	})
	overrides := inboundEndpointUnitOverrides("tf-testAcc-unit", ipConfigs)
	configValue := inboundEndpointUnitValue(objType, overrides)
	planOverrides := make(map[string]tftypes.Value, len(overrides)+1)
	for name, value := range overrides {
		planOverrides[name] = value
	}
	planOverrides["id"] = inboundEndpointUnitUnknownString()
	planValue := inboundEndpointUnitValue(objType, planOverrides)

	req := fwresource.CreateRequest{
		Plan:   tfsdk.Plan{Schema: resourceSchema, Raw: planValue},
		Config: tfsdk.Config{Schema: resourceSchema, Raw: configValue},
	}
	resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}
	r.Create(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create returned diagnostics: %v", resp.Diagnostics)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(tokens) != 2 {
		t.Fatalf("expected the throttled create to be retried once, got %d attempts", len(tokens))
	}
	if tokens[0] == "" {
		t.Fatal("the create request must carry a ClientToken")
	}
	if tokens[0] != tokens[1] {
		t.Fatalf("every attempt of one create must reuse the same ClientToken, got %q and %q", tokens[0], tokens[1])
	}
}

// TestUnitInboundEndpointDeleteWaiterDistinguishesGoneFromAccepted covers the asynchronous delete
// terminal state: an accepted DeleteInboundEndpoint is not completion, and Delete must only return
// once DescribeInboundEndpoint reports the endpoint is gone.
func TestUnitInboundEndpointDeleteWaiterDistinguishesGoneFromAccepted(t *testing.T) {
	ctx := context.Background()
	var mu sync.Mutex
	describes := 0

	client := inboundEndpointUnitClient(t, func(w http.ResponseWriter, form url.Values) {
		switch form.Get("Action") {
		case "DeleteInboundEndpoint":
			inboundEndpointUnitRespond(t, w, http.StatusOK, map[string]interface{}{"RequestId": "req-delete"})
		case "DescribeInboundEndpoint":
			mu.Lock()
			describes++
			attempt := describes
			mu.Unlock()
			if attempt == 1 {
				inboundEndpointUnitRespond(t, w, http.StatusOK, map[string]interface{}{
					"RequestId": "req-describe-1", "Id": "ep-1", "Status": "CHANGE_INIT", "IpConfigs": []interface{}{},
				})
				return
			}
			inboundEndpointUnitRespond(t, w, http.StatusNotFound, map[string]interface{}{
				"Code": "InboundEndpoint.NotExsts", "Message": "code: 404, endpoint not exists", "RequestId": "req-describe-2",
			})
		default:
			t.Errorf("unexpected action %q", form.Get("Action"))
			inboundEndpointUnitRespond(t, w, http.StatusOK, map[string]interface{}{"RequestId": "req-unexpected"})
		}
	})
	r := inboundEndpointUnitResource(t, client)
	resourceSchema, objType := inboundEndpointUnitSchema(t, r)

	ipConfigs := inboundEndpointUnitIpConfigs(t, objType, []map[string]tftypes.Value{
		{
			"az_id":      inboundEndpointUnitString("cn-hangzhou-a"),
			"cidr_block": inboundEndpointUnitString("172.16.0.0/20"),
			"vswitch_id": inboundEndpointUnitString("vsw-1"),
			"ip":         inboundEndpointUnitString("172.16.0.10"),
		},
		{
			"az_id":      inboundEndpointUnitString("cn-hangzhou-b"),
			"cidr_block": inboundEndpointUnitString("172.16.16.0/20"),
			"vswitch_id": inboundEndpointUnitString("vsw-2"),
			"ip":         inboundEndpointUnitString("172.16.0.11"),
		},
	})
	_, stateValue := inboundEndpointUnitStateValue(t, resourceSchema, objType, "tf-testAcc-unit", ipConfigs)

	req := fwresource.DeleteRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	resp := &fwresource.DeleteResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	r.Delete(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Delete returned diagnostics: %v", resp.Diagnostics)
	}

	mu.Lock()
	defer mu.Unlock()
	if describes < 2 {
		t.Fatalf("the waiter must keep polling until DescribeInboundEndpoint reports not found, got %d describe calls", describes)
	}
}

func TestUnitInboundEndpointRequestsHonorContextCancellation(t *testing.T) {
	for _, tc := range []struct {
		name   string
		action string
		call   func(context.Context, *privatezone.Service) error
	}{
		{
			name: "describe", action: "DescribeInboundEndpoint",
			call: func(ctx context.Context, service *privatezone.Service) error {
				_, err := service.DescribeInboundEndpointContext(ctx, "ep-1")
				return err
			},
		},
		{
			name: "create", action: "AddInboundEndpoint",
			call: func(ctx context.Context, service *privatezone.Service) error {
				_, err := service.CreateInboundEndpoint(ctx, map[string]interface{}{"Name": "tf-testAcc-unit"})
				return err
			},
		},
		{
			name: "update", action: "UpdateInboundEndpoint",
			call: func(ctx context.Context, service *privatezone.Service) error {
				_, err := service.UpdateInboundEndpoint(ctx, map[string]interface{}{"EndpointId": "ep-1"})
				return err
			},
		},
		{
			name: "delete", action: "DeleteInboundEndpoint",
			call: func(ctx context.Context, service *privatezone.Service) error {
				_, err := service.DeleteInboundEndpoint(ctx, map[string]interface{}{"EndpointId": "ep-1"})
				return err
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			requestStarted := make(chan struct{})
			requestCanceled := make(chan error, 1)
			var started sync.Once
			client := inboundEndpointUnitRequestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parsing request form: %s", err)
				}
				if got := r.Form.Get("Action"); got != tc.action {
					t.Errorf("Action = %q, want %q", got, tc.action)
				}
				started.Do(func() { close(requestStarted) })
				select {
				case <-r.Context().Done():
					requestCanceled <- r.Context().Err()
				case <-time.After(5 * time.Second):
					t.Error("the HTTP request did not observe context cancellation")
					inboundEndpointUnitRespond(t, w, http.StatusOK, map[string]interface{}{"RequestId": "req-timeout"})
				}
			})
			service := privatezone.NewService(client)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			callResult := make(chan error, 1)
			go func() { callResult <- tc.call(ctx, service) }()

			select {
			case <-requestStarted:
			case <-time.After(5 * time.Second):
				t.Fatal("the API request did not reach the HTTP handler")
			}
			cancel()
			select {
			case err := <-requestCanceled:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("request context error = %v, want context.Canceled", err)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("the server did not observe request cancellation")
			}
			select {
			case err := <-callResult:
				if err == nil {
					t.Fatal("a cancelled API call must return an error")
				}
			case <-time.After(5 * time.Second):
				t.Fatal("the API call did not return after request cancellation")
			}
		})
	}
}

// TestUnitInboundEndpointDescribeClassifiesNotFoundOnlyForExactCodes proves the Read path does not
// turn an endpoint-scoped failure into a missing resource: only a typed 404 or the two not-exists
// spellings are NotFound, while the short code the service also returns is not.
func TestUnitInboundEndpointDescribeClassifiesNotFoundOnlyForExactCodes(t *testing.T) {
	cases := []struct {
		name     string
		status   int
		code     string
		wantGone bool
	}{
		{name: "typed 404", status: http.StatusNotFound, code: "InboundEndpoint.NotExsits", wantGone: true},
		{name: "not-exists spelling", status: http.StatusBadRequest, code: "InboundEndpoint.NotExists", wantGone: true},
		{name: "misspelled not-exists", status: http.StatusBadRequest, code: "InboundEndpoint.NotExsits", wantGone: true},
		{name: "short endpoint code", status: http.StatusBadRequest, code: "InboundEndpoint", wantGone: false},
		{name: "near-miss prefix", status: http.StatusBadRequest, code: "InboundEndpoint.Not", wantGone: false},
		{name: "permission denied", status: http.StatusForbidden, code: "Forbidden.RAM", wantGone: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client := inboundEndpointUnitClient(t, func(w http.ResponseWriter, form url.Values) {
				inboundEndpointUnitRespond(t, w, tc.status, map[string]interface{}{
					"Code": tc.code, "Message": "request failed", "RequestId": "req-describe",
				})
			})
			service := privatezone.NewService(client)
			_, err := service.DescribeInboundEndpointContext(context.Background(), "ep-1")
			if err == nil {
				t.Fatal("expected the failed describe to return an error")
			}
			// DescribeInboundEndpointContext replaces a matched not-exists response with the
			// provider NotFound error, so the caller-visible classification is the provider one;
			// the raw matcher itself is covered by TestUnitIsInboundEndpointNotFound.
			if got := errs.IsProviderNotFound(err); got != tc.wantGone {
				t.Fatalf("IsProviderNotFound() = %v, want %v (code %q status %d)", got, tc.wantGone, tc.code, tc.status)
			}
		})
	}
}

// TestUnitInboundEndpointDeleteWaiterFailsOnNonNotFoundReadError keeps a failed describe from
// ending the delete successfully: a live endpoint answering with a short endpoint code is not
// proof that it is gone.
func TestUnitInboundEndpointDeleteWaiterFailsOnNonNotFoundReadError(t *testing.T) {
	ctx := context.Background()
	client := inboundEndpointUnitClient(t, func(w http.ResponseWriter, form url.Values) {
		switch form.Get("Action") {
		case "DeleteInboundEndpoint":
			inboundEndpointUnitRespond(t, w, http.StatusOK, map[string]interface{}{"RequestId": "req-delete"})
		case "DescribeInboundEndpoint":
			inboundEndpointUnitRespond(t, w, http.StatusBadRequest, map[string]interface{}{
				"Code": "InboundEndpoint", "Message": "endpoint busy", "RequestId": "req-describe",
			})
		default:
			t.Errorf("unexpected action %q", form.Get("Action"))
			inboundEndpointUnitRespond(t, w, http.StatusOK, map[string]interface{}{"RequestId": "req-unexpected"})
		}
	})
	r := inboundEndpointUnitResource(t, client)
	resourceSchema, objType := inboundEndpointUnitSchema(t, r)

	ipConfigs := inboundEndpointUnitIpConfigs(t, objType, []map[string]tftypes.Value{
		{
			"az_id":      inboundEndpointUnitString("cn-hangzhou-a"),
			"cidr_block": inboundEndpointUnitString("172.16.0.0/20"),
			"vswitch_id": inboundEndpointUnitString("vsw-1"),
			"ip":         inboundEndpointUnitString("172.16.0.10"),
		},
		{
			"az_id":      inboundEndpointUnitString("cn-hangzhou-b"),
			"cidr_block": inboundEndpointUnitString("172.16.16.0/20"),
			"vswitch_id": inboundEndpointUnitString("vsw-2"),
			"ip":         inboundEndpointUnitString("172.16.0.11"),
		},
	})
	_, stateValue := inboundEndpointUnitStateValue(t, resourceSchema, objType, "tf-testAcc-unit", ipConfigs)

	req := fwresource.DeleteRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	resp := &fwresource.DeleteResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	r.Delete(ctx, req, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("a non-NotFound describe error must fail the delete wait")
	}
	var text strings.Builder
	for _, diagErr := range resp.Diagnostics.Errors() {
		text.WriteString(diagErr.Summary())
		text.WriteString(": ")
		text.WriteString(diagErr.Detail())
		text.WriteString("\n")
	}
	if !strings.Contains(text.String(), "endpoint busy") {
		t.Fatalf("expected the original API error to surface, got %s", text.String())
	}
}

func TestUnitInboundEndpointStateMapsComputedAttributes(t *testing.T) {
	ctx := context.Background()
	client := inboundEndpointUnitClient(t, func(w http.ResponseWriter, form url.Values) {
		inboundEndpointUnitRespond(t, w, http.StatusOK, inboundEndpointUnitDescribeBody("ep-1", "tf-testAcc-unit", "SUCCESS"))
	})
	r := inboundEndpointUnitResource(t, client)
	resourceSchema, objType := inboundEndpointUnitSchema(t, r)
	ipConfigs := inboundEndpointUnitIpConfigs(t, objType, []map[string]tftypes.Value{
		{
			"az_id": inboundEndpointUnitString("cn-hangzhou-a"), "cidr_block": inboundEndpointUnitString("172.16.0.0/20"),
			"vswitch_id": inboundEndpointUnitString("vsw-explicit"), "ip": inboundEndpointUnitString("172.16.0.10"),
		},
		{
			"az_id": inboundEndpointUnitString("cn-hangzhou-b"), "cidr_block": inboundEndpointUnitString("172.16.16.0/20"),
			"vswitch_id": inboundEndpointUnitString("vsw-auto"), "ip": inboundEndpointUnitString("172.16.0.11"),
		},
	})
	overrides := inboundEndpointUnitOverrides("tf-testAcc-unit", ipConfigs)
	overrides["id"] = inboundEndpointUnitString("ep-1")
	stateValue := inboundEndpointUnitValue(objType, overrides)
	req := fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	resp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	r.Read(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics: %v", resp.Diagnostics)
	}
	for attribute, want := range map[string]string{
		"vpc_name": "vpc-unit", "create_time": "2026-11-18T06:00:00Z", "update_time": "2026-11-18T06:05:00Z",
	} {
		var got types.String
		if diagnostics := resp.State.GetAttribute(ctx, path.Root(attribute), &got); diagnostics.HasError() {
			t.Fatalf("decoding %s: %v", attribute, diagnostics)
		}
		if got.ValueString() != want {
			t.Fatalf("%s = %q, want %q", attribute, got.ValueString(), want)
		}
	}
}

func TestUnitInboundEndpointReadPreservesTimeouts(t *testing.T) {
	ctx := context.Background()
	client := inboundEndpointUnitClient(t, func(w http.ResponseWriter, form url.Values) {
		inboundEndpointUnitRespond(t, w, http.StatusOK, inboundEndpointUnitDescribeBody("ep-1", "tf-testAcc-unit", "SUCCESS"))
	})
	r := inboundEndpointUnitResource(t, client)
	resourceSchema, objType := inboundEndpointUnitSchema(t, r)
	ipConfigs := inboundEndpointUnitIpConfigs(t, objType, []map[string]tftypes.Value{
		{
			"az_id": inboundEndpointUnitString("cn-hangzhou-a"), "cidr_block": inboundEndpointUnitString("172.16.0.0/20"),
			"vswitch_id": inboundEndpointUnitString("vsw-1"), "ip": inboundEndpointUnitString("172.16.0.10"),
		},
		{
			"az_id": inboundEndpointUnitString("cn-hangzhou-b"), "cidr_block": inboundEndpointUnitString("172.16.16.0/20"),
			"vswitch_id": inboundEndpointUnitString("vsw-2"), "ip": inboundEndpointUnitString("172.16.0.11"),
		},
	})
	overrides := inboundEndpointUnitOverrides("tf-testAcc-unit", ipConfigs)
	overrides["id"] = inboundEndpointUnitString("ep-1")
	overrides["timeouts"] = inboundEndpointUnitTimeouts(t, objType, map[string]string{
		"create": "11m", "update": "12m", "delete": "13m",
	})
	stateValue := inboundEndpointUnitValue(objType, overrides)
	req := fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	resp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	r.Read(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics: %v", resp.Diagnostics)
	}
	var value timeouts.Value
	if diagnostics := resp.State.GetAttribute(ctx, path.Root("timeouts"), &value); diagnostics.HasError() {
		t.Fatalf("decoding refreshed timeouts: %v", diagnostics)
	}
	for operation, tc := range map[string]struct {
		accessor func(context.Context, time.Duration) (time.Duration, diag.Diagnostics)
		want     time.Duration
	}{
		"create": {value.Create, 11 * time.Minute},
		"update": {value.Update, 12 * time.Minute},
		"delete": {value.Delete, 13 * time.Minute},
	} {
		got, diagnostics := tc.accessor(ctx, 10*time.Minute)
		if diagnostics.HasError() || got != tc.want {
			t.Fatalf("%s timeout = %s, want %s, diagnostics = %v", operation, got, tc.want, diagnostics)
		}
	}
}

// TestUnitInboundEndpointStateRejectsMalformedIpConfigs covers the response mapping: a malformed
// IpConfigs payload must produce a diagnostic instead of an empty or partially populated state.
func TestUnitInboundEndpointStateRejectsMalformedIpConfigs(t *testing.T) {
	ctx := context.Background()
	validBody := func() map[string]interface{} {
		return map[string]interface{}{
			"RequestId": "req-describe",
			"Id":        "ep-1",
			"Name":      "tf-testAcc-unit",
			"Status":    "SUCCESS",
			"IpConfigs": []interface{}{
				map[string]interface{}{
					"AzId":      "cn-hangzhou-a",
					"CidrBlock": "172.16.0.0/20",
					"Ip":        "172.16.0.10",
					"VSwitchId": "vsw-1",
				},
			},
		}
	}
	ipConfigs := func(t *testing.T, objType tftypes.Object) tftypes.Value {
		return inboundEndpointUnitIpConfigs(t, objType, []map[string]tftypes.Value{
			{
				"az_id":      inboundEndpointUnitString("cn-hangzhou-a"),
				"cidr_block": inboundEndpointUnitString("172.16.0.0/20"),
				"vswitch_id": inboundEndpointUnitString("vsw-1"),
				"ip":         inboundEndpointUnitString("172.16.0.10"),
			},
			{
				"az_id":      inboundEndpointUnitString("cn-hangzhou-b"),
				"cidr_block": inboundEndpointUnitString("172.16.16.0/20"),
				"vswitch_id": inboundEndpointUnitString("vsw-2"),
				"ip":         inboundEndpointUnitString("172.16.0.11"),
			},
		})
	}
	readState := func(t *testing.T, body map[string]interface{}) (*fwresource.ReadResponse, fwschema.Schema) {
		client := inboundEndpointUnitClient(t, func(w http.ResponseWriter, form url.Values) {
			inboundEndpointUnitRespond(t, w, http.StatusOK, body)
		})
		r := inboundEndpointUnitResource(t, client)
		resourceSchema, objType := inboundEndpointUnitSchema(t, r)
		_, stateValue := inboundEndpointUnitStateValue(t, resourceSchema, objType, "tf-testAcc-unit", ipConfigs(t, objType))
		req := fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
		resp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
		r.Read(ctx, req, resp)
		return resp, resourceSchema
	}

	validResp, _ := readState(t, validBody())
	if validResp.Diagnostics.HasError() {
		t.Fatalf("a valid response must not produce diagnostics: %v", validResp.Diagnostics)
	}
	var readIpConfigs types.Set
	if diags := validResp.State.GetAttribute(ctx, path.Root("ip_configs"), &readIpConfigs); diags.HasError() {
		t.Fatalf("reading ip_configs from state: %v", diags)
	}
	if got := len(readIpConfigs.Elements()); got != 1 {
		t.Fatalf("expected one ip_config, got %d", got)
	}

	cases := []struct {
		name   string
		mutate func(map[string]interface{})
	}{
		{name: "missing array", mutate: func(o map[string]interface{}) { delete(o, "IpConfigs") }},
		{name: "non array", mutate: func(o map[string]interface{}) { o["IpConfigs"] = "not-an-array" }},
		{name: "empty array", mutate: func(o map[string]interface{}) { o["IpConfigs"] = []interface{}{} }},
		{name: "non object element", mutate: func(o map[string]interface{}) {
			o["IpConfigs"] = []interface{}{"not-an-object"}
		}},
		{name: "missing az id", mutate: func(o map[string]interface{}) {
			o["IpConfigs"] = []interface{}{map[string]interface{}{
				"CidrBlock": "172.16.0.0/20",
				"VSwitchId": "vsw-1",
			}}
		}},
		{name: "non string cidr block", mutate: func(o map[string]interface{}) {
			o["IpConfigs"] = []interface{}{map[string]interface{}{
				"AzId":      "cn-hangzhou-a",
				"CidrBlock": 17216,
				"VSwitchId": "vsw-1",
			}}
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body := validBody()
			tc.mutate(body)
			resp, _ := readState(t, body)
			if !resp.Diagnostics.HasError() {
				t.Fatalf("expected a diagnostic for %s", tc.name)
			}
		})
	}
}
