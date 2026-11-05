package dlfnext

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func catalogNotFoundErr() error {
	return &tea.SDKError{
		StatusCode: tea.Int(404),
		Code:       tea.String("NotFound"),
		Message:    tea.String("ResourceNotfound!!! the specified catalog is not found"),
	}
}

func catalogAlreadyExistsErr() error {
	return &tea.SDKError{
		StatusCode: tea.Int(409),
		Code:       tea.String("CatalogAlreadyExists"),
		Message:    tea.String("the catalog already exists"),
	}
}

func catalogValidationErr() error {
	return &tea.SDKError{
		StatusCode: tea.Int(400),
		Code:       tea.String("InvalidParameter"),
		Message:    tea.String("the request parameter is invalid"),
	}
}

func catalogThrottlingErr() error {
	return &tea.SDKError{
		StatusCode: tea.Int(503),
		Code:       tea.String("ServiceUnavailable"),
		Message:    tea.String("ServiceUnavailable: Request timed out"),
	}
}

func catalogThrottling400Err() error {
	return &tea.SDKError{
		StatusCode: tea.Int(400),
		Code:       tea.String("Throttling.User"),
		Message:    tea.String("the request was throttled"),
	}
}

func catalogInternalErr() error {
	return &tea.SDKError{
		StatusCode: tea.Int(500),
		Code:       tea.String("InternalError"),
		Message:    tea.String("internal server error"),
	}
}

func TestUnitDlfNextCatalogCreateRecovery(t *testing.T) {
	if dec, err := decideCatalogCreateRecovery(true, false, true, catalogGetFound); dec != catalogCreateFail {
		t.Fatalf("AlreadyExists must fail, got %v", dec)
	} else if err == nil {
		t.Fatalf("AlreadyExists must return an error")
	}
	if dec, err := decideCatalogCreateRecovery(false, true, false, catalogGetFound); dec != catalogCreateFail {
		t.Fatalf("validation/permission error with a same-name catalog must fail, got %v", dec)
	} else if err == nil {
		t.Fatalf("validation/permission error must return an error")
	}
	if dec, err := decideCatalogCreateRecovery(false, false, true, catalogGetFound); dec != catalogCreateFail {
		t.Fatalf("ambiguous error with a same-name catalog must fail, got %v err=%v", dec, err)
	} else if err == nil {
		t.Fatalf("ambiguous error with a found catalog must return an actionable error")
	}
	if dec, err := decideCatalogCreateRecovery(false, false, true, catalogGetAbsent); dec != catalogCreateRetry {
		t.Fatalf("retryable error with absent catalog must retry, got %v err=%v", dec, err)
	}
	if dec, err := decideCatalogCreateRecovery(false, false, false, catalogGetAbsent); dec != catalogCreateFail {
		t.Fatalf("non-retryable error with absent catalog must fail, got %v", dec)
	} else if err == nil {
		t.Fatalf("non-retryable error with absent catalog must return an error")
	}
	if dec, err := decideCatalogCreateRecovery(false, false, true, catalogGetOtherError); dec != catalogCreateFail {
		t.Fatalf("post-error Get error must fail, got %v", dec)
	} else if err == nil {
		t.Fatalf("post-error Get error must return an error")
	}
}

func scriptedGet(script []struct {
	obj map[string]interface{}
	err error
}) func(string) (map[string]interface{}, error) {
	i := 0
	return func(id string) (map[string]interface{}, error) {
		if i >= len(script) {
			i = len(script) - 1
		}
		if i < 0 {
			return nil, catalogNotFoundErr()
		}
		s := script[i]
		i++
		return s.obj, s.err
	}
}

func statusObj(status string) map[string]interface{} {
	return map[string]interface{}{"status": status, "name": "test-catalog", "type": "PAIMON"}
}

func TestUnitDlfNextCatalogCreateWaiter(t *testing.T) {
	svc := &Service{}
	ctx := context.Background()

	script := []struct {
		obj map[string]interface{}
		err error
	}{
		{nil, catalogNotFoundErr()},
		{statusObj("NEW"), nil},
		{statusObj("INITIALIZING"), nil},
		{statusObj("RUNNING"), nil},
	}
	refresh := svc.CatalogStateRefreshFuncWithApi("test-catalog", "status", []string{"INITIALIZE_FAILED"}, scriptedGet(script))
	conf := &retry.StateChangeConf{Pending: []string{"NOT_FOUND", "NEW", "INITIALIZING"}, Target: []string{"RUNNING"}, Timeout: 10 * time.Second, Refresh: refresh}
	if _, err := conf.WaitForStateContext(ctx); err != nil {
		t.Fatalf("expected waiter to reach RUNNING, got error: %v", err)
	}

	script = []struct {
		obj map[string]interface{}
		err error
	}{{statusObj("INITIALIZE_FAILED"), nil}}
	refresh = svc.CatalogStateRefreshFuncWithApi("test-catalog", "status", []string{"INITIALIZE_FAILED"}, scriptedGet(script))
	conf = &retry.StateChangeConf{Pending: []string{"NOT_FOUND", "NEW", "INITIALIZING"}, Target: []string{"RUNNING"}, Timeout: 10 * time.Second, Refresh: refresh}
	if _, err := conf.WaitForStateContext(ctx); err == nil {
		t.Fatalf("expected INITIALIZE_FAILED to fail the waiter, got nil")
	}

	script = []struct {
		obj map[string]interface{}
		err error
	}{{nil, catalogNotFoundErr()}}
	refresh = svc.CatalogStateRefreshFuncWithApi("test-catalog", "status", []string{"INITIALIZE_FAILED"}, scriptedGet(script))
	conf = &retry.StateChangeConf{Pending: []string{"NOT_FOUND", "NEW", "INITIALIZING"}, Target: []string{"RUNNING"}, Timeout: 1 * time.Second, Refresh: refresh}
	start := time.Now()
	if _, err := conf.WaitForStateContext(ctx); err == nil {
		t.Fatalf("expected perpetual NotFound to time out, got nil")
	}
	if elapsed := time.Since(start); elapsed > 10*time.Second {
		t.Fatalf("waiter with 1s timeout took %v", elapsed)
	}
}

func TestUnitDlfNextCatalogDeleteWaiter(t *testing.T) {
	svc := &Service{}
	ctx := context.Background()

	script := []struct {
		obj map[string]interface{}
		err error
	}{
		{statusObj("DELETING"), nil},
		{statusObj("DELETED"), nil},
		{nil, catalogNotFoundErr()},
	}
	if err := svc.WaitForCatalogGoneWithApi(ctx, "test-catalog", 10*time.Second, scriptedGet(script)); err != nil {
		t.Fatalf("expected delete waiter to reach gone, got error: %v", err)
	}

	script = []struct {
		obj map[string]interface{}
		err error
	}{
		{nil, catalogThrottlingErr()},
		{nil, catalogNotFoundErr()},
	}
	if err := svc.WaitForCatalogGoneWithApi(ctx, "test-catalog", 20*time.Second, scriptedGet(script)); err != nil {
		t.Fatalf("expected delete waiter to retry throttling then reach gone, got error: %v", err)
	}

	script = []struct {
		obj map[string]interface{}
		err error
	}{{nil, catalogInternalErr()}}
	if err := svc.WaitForCatalogGoneWithApi(ctx, "test-catalog", 5*time.Second, scriptedGet(script)); err == nil {
		t.Fatalf("expected a non-retryable Get error to fail the delete waiter, got nil")
	}
}

func TestUnitDlfNextCatalogStateRefreshMissingStatus(t *testing.T) {
	svc := &Service{}
	// A malformed object without "status" must surface the extraction error
	// instead of becoming the string "<nil>".
	refresh := svc.CatalogStateRefreshFuncWithApi("test-catalog", "status", nil, func(id string) (map[string]interface{}, error) {
		return map[string]interface{}{"name": "test-catalog"}, nil
	})
	if _, _, err := refresh(); err == nil {
		t.Fatalf("expected the missing status to surface an error")
	}
}

func TestUnitDlfNextCatalogsStableID(t *testing.T) {
	id1 := catalogsID([]string{"alpha", "bravo", "charlie"})
	id2 := catalogsID([]string{"charlie", "alpha", "bravo"})
	if id1 != id2 || id1 == "" {
		t.Fatalf("expected stable non-empty id, got %q and %q", id1, id2)
	}
	if catalogsID([]string{"a", "b"}) == catalogsID([]string{"a", "c"}) {
		t.Fatalf("expected different sets to yield different ids")
	}
}

func TestUnitDlfNextCatalogsStableOrder(t *testing.T) {
	order1 := []catalogItem{{Name: "charlie"}, {Name: "alpha"}, {Name: "bravo"}}
	sorted := sortCatalogItems(order1)
	want := []string{"alpha", "bravo", "charlie"}
	for i, w := range want {
		if sorted[i].Name != w {
			t.Fatalf("sorted[%d] = %q, want %q", i, sorted[i].Name, w)
		}
	}
	if order1[0].Name != "charlie" {
		t.Fatalf("sortCatalogItems must not mutate its input")
	}
}

func TestUnitDlfNextCatalogErrorClassifiersNilSafe(t *testing.T) {
	if isCatalogAlreadyExistsError(catalogValidationErr()) {
		t.Fatalf("validation error must not be classified as already-exists")
	}
	if !isCatalogExplicitClientError(catalogValidationErr()) {
		t.Fatalf("validation error (HTTP 400) must be classified as explicit client error")
	}
	if !isCatalogAlreadyExistsError(catalogAlreadyExistsErr()) {
		t.Fatalf("already-exists error must be classified as already-exists")
	}

	allNilErr := &tea.SDKError{}
	if errs.IsExpectedErrors(allNilErr, []string{"AnyCode"}) {
		t.Fatalf("all-nil SDKError must not match any expected code")
	}
	if isCatalogAlreadyExistsError(allNilErr) || isCatalogExplicitClientError(allNilErr) {
		t.Fatalf("all-nil SDKError must not match a classifier")
	}
	if helper.NeedRetry(allNilErr) {
		t.Fatalf("all-nil SDKError must not be classified as retryable")
	}

	nilDataErr := &tea.SDKError{Code: tea.String("SomeOtherCode"), Data: nil, Message: tea.String("msg"), StatusCode: tea.Int(500)}
	if errs.IsExpectedErrors(nilDataErr, []string{"NotFound"}) {
		t.Fatalf("non-matching code with nil Data must not match")
	}

	for _, code := range []int{400, 401, 403} {
		err := &tea.SDKError{Code: tea.String("SomeCode"), Data: nil, Message: tea.String("msg"), StatusCode: tea.Int(code)}
		if !isCatalogExplicitClientError(err) {
			t.Fatalf("HTTP %d must be classified as explicit client error", code)
		}
	}
	if isCatalogExplicitClientError(catalogThrottlingErr()) {
		t.Fatalf("HTTP 503 must not be classified as explicit client error")
	}
}

func TestUnitDlfNextCatalogExplicitClientThrottling400(t *testing.T) {
	throttling400 := catalogThrottling400Err()
	if !helper.NeedRetry(throttling400) {
		t.Fatalf("expected the throttling 400 error to be classified as retryable")
	}
	if isCatalogExplicitClientError(throttling400) {
		t.Fatalf("throttling 400 must not be classified as an explicit client error")
	}
}

func TestUnitDlfNextCatalogTimeoutPolicy(t *testing.T) {
	if got := remainingUntil(time.Now().Add(-1 * time.Second)); got != 0 {
		t.Fatalf("remainingUntil of a passed deadline must be 0, got %v", got)
	}
	if got := remainingUntil(time.Now().Add(2 * time.Second)); got <= 0 || got > 2*time.Second {
		t.Fatalf("remainingUntil of a 2s-future deadline must be in (0, 2s], got %v", got)
	}
}

func TestUnitDlfNextCatalogOptionChanges(t *testing.T) {
	u, r := buildOptionChanges(map[string]string{}, map[string]string{"a": "1", "b": "2"})
	if len(u) != 2 || u["a"] != "1" || len(r) != 0 {
		t.Fatalf("add: updates=%v removals=%v", u, r)
	}
	u, r = buildOptionChanges(map[string]string{"a": "1", "b": "2"}, map[string]string{"a": "1"})
	if len(u) != 1 || len(r) != 1 || fmt.Sprint(r[0]) != "b" {
		t.Fatalf("remove: updates=%v removals=%v", u, r)
	}
	u, r = buildOptionChanges(map[string]string{"a": "1", "b": "2"}, map[string]string{})
	if len(u) != 0 || len(r) != 2 {
		t.Fatalf("clear: updates=%v removals=%v", u, r)
	}
	u, r = buildOptionChanges(nil, nil)
	if len(u) != 0 || len(r) != 0 {
		t.Fatalf("nil/nil: updates=%v removals=%v", u, r)
	}
}

func TestUnitDlfNextCatalogOptionChangesTwoState(t *testing.T) {
	oldManaged := map[string]string{"comment": "v1", "description": "d", "dlf.discovery-query-results-retained-days": "7"}
	_, rOmit := buildOptionChanges(oldManaged, nil)
	_, rEmpty := buildOptionChanges(oldManaged, map[string]string{})
	if len(rOmit) != len(oldManaged) || len(rEmpty) != len(oldManaged) {
		t.Fatalf("two-state clear: omit=%d empty=%d want %d", len(rOmit), len(rEmpty), len(oldManaged))
	}
	setOmit, setEmpty := map[string]bool{}, map[string]bool{}
	for _, k := range rOmit {
		setOmit[fmt.Sprint(k)] = true
	}
	for _, k := range rEmpty {
		setEmpty[fmt.Sprint(k)] = true
	}
	for k := range oldManaged {
		if !setOmit[k] || !setEmpty[k] {
			t.Fatalf("key %q must be removed by both omit and {}", k)
		}
	}
}

func TestUnitDlfNextCatalogManagedOptionsReconcile(t *testing.T) {
	managed := map[string]string{"comment": "v1", "description": "d", "description2": "x"}
	api := map[string]interface{}{"comment": "v2", "description": "d", "dlf.trashed-file-retained-days": 7}
	got := reconcileManagedOptions(managed, api)
	if got["comment"] != "v2" || got["description"] != "d" {
		t.Fatalf("unexpected reconcile: %v", got)
	}
	if _, ok := got["description2"]; ok {
		t.Fatalf("managed key absent from API must be dropped")
	}
	if _, ok := got["dlf.trashed-file-retained-days"]; ok {
		t.Fatalf("unmanaged key must not be promoted")
	}
}

func TestUnitDlfNextCatalogRemoteOptionsMirror(t *testing.T) {
	api := map[string]interface{}{"comment": "v1", "dlf.trashed-file-retained-days": 7, "enable.openapi": false}
	got := reconcileRemoteOptions(api)
	if got["comment"] != "v1" || got["dlf.trashed-file-retained-days"] != "7" || got["enable.openapi"] != "false" {
		t.Fatalf("unexpected mirror: %v", got)
	}
}

func TestUnitDlfNextCatalogTimestampString(t *testing.T) {
	cases := []struct {
		in   interface{}
		want string
	}{
		{float64(1744970111419), "1744970111419"},
		{int64(1744970111419), "1744970111419"},
		{"1744970111419", "1744970111419"},
		{nil, ""},
	}
	for _, c := range cases {
		if got := catalogTimestampString(c.in); got != c.want {
			t.Fatalf("catalogTimestampString(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestUnitDlfNextCatalogCreateBody(t *testing.T) {
	ctx := context.Background()
	options, _ := types.MapValueFrom(ctx, types.StringType, map[string]string{"comment": "c1"})
	plan := catalogResourceModel{
		Name:     types.StringValue("cat"),
		Type:     types.StringValue("PAIMON"),
		Options:  options,
		IsShared: types.BoolValue(false),
		ShareId:  types.StringValue("share-1"),
	}
	body, err := buildCatalogCreateBody(ctx, plan)
	if err != nil {
		t.Fatalf("buildCatalogCreateBody: %v", err)
	}
	if v, ok := body["isShared"]; !ok || v != false {
		t.Fatalf("explicit is_shared=false must be sent as false, body=%v", body)
	}
	if body["shareId"] != "share-1" {
		t.Fatalf("shareId missing: %v", body)
	}

	omitted := catalogResourceModel{
		Name:     types.StringValue("cat"),
		Type:     types.StringValue("PAIMON"),
		IsShared: types.BoolNull(),
		ShareId:  types.StringNull(),
	}
	body, err = buildCatalogCreateBody(ctx, omitted)
	if err != nil {
		t.Fatalf("buildCatalogCreateBody: %v", err)
	}
	if _, ok := body["isShared"]; ok {
		t.Fatalf("omitted is_shared must not be sent, body=%v", body)
	}
	if _, ok := body["shareId"]; ok {
		t.Fatalf("empty share_id must not be sent, body=%v", body)
	}
}

func TestUnitDlfNextCatalogFilterDiscrimination(t *testing.T) {
	objects := []interface{}{
		map[string]interface{}{"name": "target", "type": "PAIMON", "createdAt": float64(1744970111419)},
		map[string]interface{}{"name": "other", "type": "PAIMON"},
	}
	names, items := filterCatalogs(objects, map[string]bool{"target": true}, nil)
	if len(names) != 1 || names[0] != "target" || len(items) != 1 {
		t.Fatalf("ids filter must return only the target, got %v", names)
	}
	if items[0].Id != "target" || items[0].CreatedAt != "1744970111419" {
		t.Fatalf("unexpected item projection: %+v", items[0])
	}
	names, _ = filterCatalogs(objects, map[string]bool{"missing": true}, nil)
	if len(names) != 0 {
		t.Fatalf("absent ids target must return no objects, got %v", names)
	}
}

func TestUnitDlfNextCatalogPagination(t *testing.T) {
	ctx := context.Background()
	pages := map[string]struct {
		objects []interface{}
		next    string
	}{
		"":   {objects: []interface{}{map[string]interface{}{"name": "a"}}, next: "t1"},
		"t1": {objects: []interface{}{map[string]interface{}{"name": "b"}}, next: ""},
	}
	calls := 0
	objects, err := listAllCatalogs(ctx, func(pattern, token string, max int) ([]interface{}, string, error) {
		calls++
		page := pages[token]
		return page.objects, page.next, nil
	}, "pattern")
	if err != nil {
		t.Fatalf("listAllCatalogs: %v", err)
	}
	if calls != 2 || len(objects) != 2 {
		t.Fatalf("expected 2 pages and 2 objects, got calls=%d objects=%d", calls, len(objects))
	}
}

func TestUnitDlfNextCatalogOutputFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "catalogs.json")
	if err := writeCatalogOutput(path, []catalogItem{{Name: "a", Id: "a"}}); err != nil {
		t.Fatalf("writeCatalogOutput: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}
	if len(data) == 0 || string(data) == "[]" {
		t.Fatalf("output file must contain the matched catalog, got %q", string(data))
	}
}

func TestUnitDlfNextCatalogBackoffCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	start := time.Now()
	backoff := newCatalogBackoff(time.Hour, time.Hour)
	if err := backoff.Wait(ctx, time.Hour); err == nil {
		t.Fatalf("expected a cancelled context to stop the wait")
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("cancelled wait took %v, expected immediate return", elapsed)
	}
}

func TestUnitDlfNextCatalogPathEscaping(t *testing.T) {
	if got := catalogAction("plain-catalog"); got != "/dlf/v1/catalogs/plain-catalog" {
		t.Fatalf("unexpected action for a plain name: %q", got)
	}
	if got := catalogAction("a b/c"); got != "/dlf/v1/catalogs/a%20b%2Fc" {
		t.Fatalf("reserved characters must be escaped as one path segment, got %q", got)
	}
}
