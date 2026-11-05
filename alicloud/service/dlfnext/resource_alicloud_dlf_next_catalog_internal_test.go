package dlfnext

import (
	"fmt"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// catalogNotFoundErr builds a tea.SDKError that errs.NotFoundError recognises
// as NotFound (HTTP 404), for unit tests that script GetCatalog responses
// without a live backend.
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

func TestUnitDlfNextCatalogSchemaType(t *testing.T) {
	r := ResourceCatalog()
	typeSchema := r.Schema["type"]

	if typeSchema.Default != "PAIMON" {
		t.Fatalf("expected type Default to be PAIMON, got %v", typeSchema.Default)
	}

	// CreateCatalog documents enum ["PAIMON", "ICEBERG"]; the schema must
	// accept both contract values.
	for _, valid := range []string{"PAIMON", "ICEBERG"} {
		if _, errs := typeSchema.ValidateFunc(valid, "type"); len(errs) != 0 {
			t.Fatalf("expected %s to be accepted by the type ValidateFunc, got %v", valid, errs)
		}
	}

	if _, errs := typeSchema.ValidateFunc("DELTA", "type"); len(errs) == 0 {
		t.Fatal("expected an undocumented type to be rejected by the type ValidateFunc")
	}
}

func TestUnitDlfNextCatalogTimestampString(t *testing.T) {
	// GetCatalog returns int64 milliseconds; the string form must keep the
	// full value (32-bit int would overflow).
	cases := []struct {
		in   interface{}
		want string
	}{
		{float64(1744970111419), "1744970111419"},
		{int64(1744970111419), "1744970111419"},
		{int(1744970111), "1744970111"},
		{"1744970111419", "1744970111419"},
		{nil, ""},
	}
	for _, c := range cases {
		if got := catalogTimestampString(c.in); got != c.want {
			t.Fatalf("catalogTimestampString(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestUnitDlfNextCatalogPreCreate(t *testing.T) {
	// 1. Pre-create: catalog already exists -> refuse, do not take over.
	if proceed, err := decideCatalogPreCreate(nil); proceed {
		t.Fatalf("pre-create with an existing catalog must not proceed")
	} else if err == nil {
		t.Fatalf("pre-create with an existing catalog must return a conflict error")
	}

	// 2. Pre-create: NotFound -> proceed.
	if proceed, err := decideCatalogPreCreate(catalogNotFoundErr()); !proceed || err != nil {
		t.Fatalf("pre-create NotFound must proceed, got proceed=%v err=%v", proceed, err)
	}

	// 3. Pre-create: other Get error -> fail (do not proceed).
	if proceed, err := decideCatalogPreCreate(catalogInternalErr()); proceed {
		t.Fatalf("pre-create with a non-NotFound Get error must not proceed")
	} else if err == nil {
		t.Fatalf("pre-create with a non-NotFound Get error must surface the error")
	}
}

func TestUnitDlfNextCatalogCreateRecovery(t *testing.T) {
	// AlreadyExists -> fail, never import.
	if dec, err := decideCatalogCreateRecovery(true, false, true, true, catalogGetFound); dec != catalogCreateFail {
		t.Fatalf("AlreadyExists must fail, got %v", dec)
	} else if err == nil {
		t.Fatalf("AlreadyExists must return an error")
	}

	// Validation/permission error after a same-name catalog exists -> fail.
	if dec, err := decideCatalogCreateRecovery(false, true, false, true, catalogGetFound); dec != catalogCreateFail {
		t.Fatalf("validation/permission error with a same-name catalog must fail, got %v", dec)
	} else if err == nil {
		t.Fatalf("validation/permission error must return an error")
	}

	// Ambiguous error and Get finds a same-name catalog -> MUST fail.
	dec, err := decideCatalogCreateRecovery(false, false, true, true, catalogGetFound)
	if dec != catalogCreateFail {
		t.Fatalf("ambiguous error with a same-name catalog must fail (not recover), got %v err=%v", dec, err)
	}
	if err == nil {
		t.Fatalf("ambiguous error with a found catalog must return an actionable error")
	}

	// Retryable error and resource was not created (Get NotFound) -> retry.
	if dec, err := decideCatalogCreateRecovery(false, false, true, true, catalogGetAbsent); dec != catalogCreateRetry {
		t.Fatalf("retryable error with absent catalog must retry, got %v err=%v", dec, err)
	}

	// Non-retryable error and catalog absent -> fail.
	if dec, err := decideCatalogCreateRecovery(false, false, false, true, catalogGetAbsent); dec != catalogCreateFail {
		t.Fatalf("non-retryable error with absent catalog must fail, got %v", dec)
	} else if err == nil {
		t.Fatalf("non-retryable error with absent catalog must return an error")
	}

	// pre-create did NOT prove absence -> fail.
	if dec, err := decideCatalogCreateRecovery(false, false, true, false, catalogGetFound); dec != catalogCreateFail {
		t.Fatalf("pre-create not confirmed absent must fail, got %v", dec)
	} else if err == nil {
		t.Fatalf("pre-create not confirmed absent must return an error")
	}

	// Post-error Get itself errored -> fail.
	if dec, err := decideCatalogCreateRecovery(false, false, true, true, catalogGetOtherError); dec != catalogCreateFail {
		t.Fatalf("post-error Get error must fail, got %v", dec)
	} else if err == nil {
		t.Fatalf("post-error Get error must return an error")
	}
}

// scriptedGet plays back a sequence of (object, error) responses, holding on
// the last response once the script is exhausted.
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

	// NOT_FOUND -> NEW -> INITIALIZING -> RUNNING reaches RUNNING.
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
	conf := &retry.StateChangeConf{
		Pending: []string{"NOT_FOUND", "NEW", "INITIALIZING"},
		Target:  []string{"RUNNING"},
		Timeout: 10 * time.Second,
		Refresh: refresh,
	}
	if _, err := conf.WaitForState(); err != nil {
		t.Fatalf("expected waiter to reach RUNNING, got error: %v", err)
	}

	// INITIALIZE_FAILED is a terminal failure.
	script = []struct {
		obj map[string]interface{}
		err error
	}{{statusObj("INITIALIZE_FAILED"), nil}}
	refresh = svc.CatalogStateRefreshFuncWithApi("test-catalog", "status", []string{"INITIALIZE_FAILED"}, scriptedGet(script))
	conf = &retry.StateChangeConf{
		Pending: []string{"NOT_FOUND", "NEW", "INITIALIZING"},
		Target:  []string{"RUNNING"},
		Timeout: 10 * time.Second,
		Refresh: refresh,
	}
	if _, err := conf.WaitForState(); err == nil {
		t.Fatalf("expected INITIALIZE_FAILED to fail the waiter, got nil")
	}

	// Perpetual NotFound times out (NOT_FOUND is Pending).
	script = []struct {
		obj map[string]interface{}
		err error
	}{{nil, catalogNotFoundErr()}}
	refresh = svc.CatalogStateRefreshFuncWithApi("test-catalog", "status", []string{"INITIALIZE_FAILED"}, scriptedGet(script))
	conf = &retry.StateChangeConf{
		Pending: []string{"NOT_FOUND", "NEW", "INITIALIZING"},
		Target:  []string{"RUNNING"},
		Timeout: 2 * time.Second,
		Refresh: refresh,
	}
	if _, err := conf.WaitForState(); err == nil {
		t.Fatalf("expected perpetual NotFound to time out, got nil")
	}

	// A non-NotFound Get error during polling surfaces.
	script = []struct {
		obj map[string]interface{}
		err error
	}{{nil, catalogInternalErr()}}
	refresh = svc.CatalogStateRefreshFuncWithApi("test-catalog", "status", []string{"INITIALIZE_FAILED"}, scriptedGet(script))
	conf = &retry.StateChangeConf{
		Pending: []string{"NOT_FOUND", "NEW", "INITIALIZING"},
		Target:  []string{"RUNNING"},
		Timeout: 5 * time.Second,
		Refresh: refresh,
	}
	if _, err := conf.WaitForState(); err == nil {
		t.Fatalf("expected a non-NotFound Get error to surface, got nil")
	}
}

func TestUnitDlfNextCatalogDeleteWaiter(t *testing.T) {
	svc := &Service{}

	// DELETING -> DELETED -> NotFound reaches gone.
	script := []struct {
		obj map[string]interface{}
		err error
	}{
		{statusObj("DELETING"), nil},
		{statusObj("DELETED"), nil},
		{nil, catalogNotFoundErr()},
	}
	if err := svc.WaitForCatalogGoneWithApi("test-catalog", 10*time.Second, scriptedGet(script)); err != nil {
		t.Fatalf("expected delete waiter to reach gone, got error: %v", err)
	}

	// A transient retryable Get error is retried until NotFound.
	script = []struct {
		obj map[string]interface{}
		err error
	}{
		{nil, catalogThrottlingErr()},
		{nil, catalogNotFoundErr()},
	}
	if err := svc.WaitForCatalogGoneWithApi("test-catalog", 20*time.Second, scriptedGet(script)); err != nil {
		t.Fatalf("expected delete waiter to retry a throttling error then reach gone, got error: %v", err)
	}

	// A non-retryable Get error fails instead of looping.
	script = []struct {
		obj map[string]interface{}
		err error
	}{{nil, catalogInternalErr()}}
	if err := svc.WaitForCatalogGoneWithApi("test-catalog", 5*time.Second, scriptedGet(script)); err == nil {
		t.Fatalf("expected a non-retryable Get error to fail the delete waiter, got nil")
	}
}

func TestUnitDlfNextCatalogsStableID(t *testing.T) {
	id1 := catalogsID([]string{"alpha", "bravo", "charlie"})
	id2 := catalogsID([]string{"charlie", "alpha", "bravo"})
	if id1 != id2 {
		t.Fatalf("expected same catalog set in different order to yield the same id, got %q and %q", id1, id2)
	}
	if id1 == "" {
		t.Fatalf("expected a non-empty data source id")
	}
	if catalogsID([]string{"a", "b"}) != catalogsID([]string{"a", "b"}) {
		t.Fatalf("expected catalogsID to be deterministic")
	}
	if catalogsID([]string{"a", "b"}) == catalogsID([]string{"a", "c"}) {
		t.Fatalf("expected different catalog sets to yield different ids")
	}
	if catalogsID(nil) != catalogsID([]string{}) {
		t.Fatalf("expected empty input to be stable")
	}
}

func TestUnitDlfNextCatalogsStableOrder(t *testing.T) {
	cat := func(name string) map[string]interface{} {
		return map[string]interface{}{"name": name, "type": "PAIMON", "status": "RUNNING"}
	}
	order1 := []map[string]interface{}{cat("charlie"), cat("alpha"), cat("bravo")}
	order2 := []map[string]interface{}{cat("bravo"), cat("charlie"), cat("alpha")}

	s1 := sortCatalogsByName(order1)
	s2 := sortCatalogsByName(order2)

	want := []string{"alpha", "bravo", "charlie"}
	for i, w := range want {
		if got := fmt.Sprint(s1[i]["name"]); got != w {
			t.Fatalf("order1[%d]: got %q, want %q", i, got, w)
		}
		if got := fmt.Sprint(s2[i]["name"]); got != w {
			t.Fatalf("order2[%d]: got %q, want %q", i, got, w)
		}
	}

	if got := fmt.Sprint(order1[0]["name"]); got != "charlie" {
		t.Fatalf("sortCatalogsByName must not mutate its input, first element got %q", got)
	}
	if got := sortCatalogsByName(nil); len(got) != 0 {
		t.Fatalf("expected nil input to sort to an empty slice, got %d elements", len(got))
	}
}

func TestUnitDlfNextCatalogErrorClassifiersNilSafe(t *testing.T) {
	// A validation error has Code=InvalidParameter, Data=nil. It must not
	// panic and must not be classified as already-exists.
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
	if isCatalogAlreadyExistsError(allNilErr) {
		t.Fatalf("all-nil SDKError must not be classified as already-exists")
	}
	if isCatalogExplicitClientError(allNilErr) {
		t.Fatalf("all-nil SDKError must not be classified as explicit client error")
	}
	if helper.NeedRetry(allNilErr) {
		t.Fatalf("all-nil SDKError must not be classified as retryable")
	}

	nilCodeErr := &tea.SDKError{
		Code:       nil,
		Data:       tea.String("some data"),
		Message:    tea.String("msg"),
		StatusCode: tea.Int(500),
	}
	if errs.IsExpectedErrors(nilCodeErr, []string{"NotFound"}) {
		t.Fatalf("nil-code error must not match any expected code")
	}

	nilDataErr := &tea.SDKError{
		Code:       tea.String("SomeOtherCode"),
		Data:       nil,
		Message:    tea.String("msg"),
		StatusCode: tea.Int(500),
	}
	if errs.IsExpectedErrors(nilDataErr, []string{"NotFound"}) {
		t.Fatalf("non-matching code with nil Data must not match")
	}

	for _, code := range []int{400, 401, 403} {
		err := &tea.SDKError{
			Code:       tea.String("SomeCode"),
			Data:       nil,
			Message:    tea.String("msg"),
			StatusCode: tea.Int(code),
		}
		if !isCatalogExplicitClientError(err) {
			t.Fatalf("HTTP %d must be classified as explicit client error", code)
		}
	}

	if isCatalogExplicitClientError(catalogThrottlingErr()) {
		t.Fatalf("HTTP 503 must not be classified as explicit client error")
	}

	wrappedErr := &errs.ComplexError{
		Cause: &tea.SDKError{
			Code:       tea.String("InvalidParameter"),
			Data:       nil,
			Message:    tea.String("msg"),
			StatusCode: tea.Int(400),
		},
	}
	if errs.IsExpectedErrors(wrappedErr, []string{"NotFound"}) {
		t.Fatalf("ComplexError wrapping non-matching SDKError must not match")
	}
}

// TestUnitDlfNextCatalogExplicitClientThrottling400 pins the classifier split
// requested in review: throttling is matched by helper.NeedRetry regardless of
// HTTP status, so a throttling error delivered with HTTP 400 must not be
// treated as a validation/permission error (which would fail the create
// instead of retrying it).
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

	svc := &Service{}
	script := []struct {
		obj map[string]interface{}
		err error
	}{{statusObj("INITIALIZING"), nil}}
	refresh := svc.CatalogStateRefreshFuncWithApi("test-catalog", "status", []string{"INITIALIZE_FAILED"}, scriptedGet(script))
	conf := &retry.StateChangeConf{
		Pending: []string{"NOT_FOUND", "NEW", "INITIALIZING"},
		Target:  []string{"RUNNING"},
		Timeout: 1 * time.Second,
		Refresh: refresh,
	}
	start := time.Now()
	_, err := conf.WaitForState()
	elapsed := time.Since(start)
	if err == nil {
		t.Fatalf("expected perpetual INITIALIZING to time out, got nil")
	}
	if elapsed > 10*time.Second {
		t.Fatalf("waiter with 1s timeout took %v, expected well under 1m (no full-timeout re-acquisition)", elapsed)
	}
}

func TestUnitDlfNextCatalogOptionChanges(t *testing.T) {
	// Add only: old empty, new {a,b} -> updates {a,b}, removals [].
	u, r := buildOptionChanges(map[string]interface{}{}, map[string]interface{}{"a": "1", "b": "2"})
	if len(u) != 2 || u["a"] != "1" || u["b"] != "2" {
		t.Fatalf("add: updates = %v, want {a:1,b:2}", u)
	}
	if len(r) != 0 {
		t.Fatalf("add: removals = %v, want empty", r)
	}

	// Modify + add.
	u, r = buildOptionChanges(map[string]interface{}{"a": "1"}, map[string]interface{}{"a": "9", "c": "3"})
	if u["a"] != "9" || u["c"] != "3" || len(u) != 2 {
		t.Fatalf("modify+add: updates = %v, want {a:9,c:3}", u)
	}
	if len(r) != 0 {
		t.Fatalf("modify+add: removals = %v, want empty", r)
	}

	// Remove one.
	u, r = buildOptionChanges(map[string]interface{}{"a": "1", "b": "2"}, map[string]interface{}{"a": "1"})
	if len(u) != 1 || u["a"] != "1" {
		t.Fatalf("remove: updates = %v, want {a:1}", u)
	}
	if len(r) != 1 || fmt.Sprint(r[0]) != "b" {
		t.Fatalf("remove: removals = %v, want [b]", r)
	}

	// Clear all: updates {}, removals {a,b}.
	u, r = buildOptionChanges(map[string]interface{}{"a": "1", "b": "2"}, map[string]interface{}{})
	if len(u) != 0 {
		t.Fatalf("clear: updates = %v, want empty", u)
	}
	if len(r) != 2 {
		t.Fatalf("clear: removals = %v, want 2 items", r)
	}
	set := map[string]bool{}
	for _, k := range r {
		set[fmt.Sprint(k)] = true
	}
	if !set["a"] || !set["b"] {
		t.Fatalf("clear: removals = %v, want {a,b}", r)
	}

	// Boundary: server defaults never appear in updates or removals.
	oldManaged := map[string]interface{}{"comment": "v1", "description": "d"}
	u, r = buildOptionChanges(oldManaged, map[string]interface{}{})
	for _, dflt := range []string{"dlf.trashed-file-retained-days", "storage.data.redundancy.type"} {
		for _, k := range r {
			if fmt.Sprint(k) == dflt {
				t.Fatalf("clear: server default %q must NEVER appear in removals, got %v", dflt, r)
			}
		}
		if _, ok := u[dflt]; ok {
			t.Fatalf("clear: server default %q must NEVER appear in updates", dflt)
		}
	}

	// No-op and nil inputs.
	u, r = buildOptionChanges(map[string]interface{}{}, map[string]interface{}{})
	if len(u) != 0 || len(r) != 0 {
		t.Fatalf("noop: updates=%v removals=%v, want both empty", u, r)
	}
	u, r = buildOptionChanges(nil, nil)
	if len(u) != 0 || len(r) != 0 {
		t.Fatalf("nil/nil: updates=%v removals=%v, want both empty", u, r)
	}
	u, r = buildOptionChanges(map[string]interface{}{"a": "1"}, nil)
	if len(u) != 0 || len(r) != 1 || fmt.Sprint(r[0]) != "a" {
		t.Fatalf("old/nil: updates=%v removals=%v, want removals [a]", u, r)
	}
}

func TestUnitDlfNextCatalogOptionChangesTwoState(t *testing.T) {
	oldManaged := map[string]interface{}{
		"comment":     "v1",
		"description": "d",
		"dlf.discovery-query-results-retained-days": "7",
	}

	uOmit, rOmit := buildOptionChanges(oldManaged, nil)
	uEmpty, rEmpty := buildOptionChanges(oldManaged, map[string]interface{}{})

	if len(uOmit) != 0 || len(uEmpty) != 0 {
		t.Fatalf("two-state clear: omit updates=%v, {} updates=%v, want both empty", uOmit, uEmpty)
	}
	if len(rOmit) != len(oldManaged) || len(rEmpty) != len(oldManaged) {
		t.Fatalf("two-state clear: omit removals=%d, {} removals=%d, want %d (all managed keys)", len(rOmit), len(rEmpty), len(oldManaged))
	}
	setOmit := map[string]bool{}
	setEmpty := map[string]bool{}
	for _, k := range rOmit {
		setOmit[fmt.Sprint(k)] = true
	}
	for _, k := range rEmpty {
		setEmpty[fmt.Sprint(k)] = true
	}
	for k := range oldManaged {
		if !setOmit[k] {
			t.Fatalf("two-state clear (omit): key %q must be in removals, got %v", k, rOmit)
		}
		if !setEmpty[k] {
			t.Fatalf("two-state clear ({}): key %q must be in removals, got %v", k, rEmpty)
		}
	}
	if len(setOmit) != len(setEmpty) {
		t.Fatalf("two-state: omit and {} removal sets differ in size")
	}
	for k := range setOmit {
		if !setEmpty[k] {
			t.Fatalf("two-state: key %q in omit removals but not in {} removals — contract broken", k)
		}
	}
	for _, dflt := range []string{"dlf.trashed-file-retained-days", "storage.data.redundancy.type"} {
		if setOmit[dflt] || setEmpty[dflt] {
			t.Fatalf("two-state clear: server default %q must NEVER be in removals", dflt)
		}
	}
}

func TestUnitDlfNextCatalogManagedOptionsReconcile(t *testing.T) {
	managed := map[string]interface{}{"comment": "v1", "description": "d", "description2": "x"}
	api := map[string]interface{}{
		"comment":                        "v2",
		"description":                    "d",
		"dlf.trashed-file-retained-days": 7,
		"storage.data.redundancy.type":   "LRS",
		"enable.openapi":                 false,
	}
	got := reconcileManagedOptions(managed, api)
	if got["comment"] != "v2" {
		t.Fatalf("managed key value change must surface API value, got comment=%v want v2", got["comment"])
	}
	if got["description"] != "d" {
		t.Fatalf("unchanged managed key must keep API value, got description=%v", got["description"])
	}
	if _, ok := got["description2"]; ok {
		t.Fatalf("managed key absent from API must be dropped (drift), but description2 present")
	}
	for _, k := range []string{"dlf.trashed-file-retained-days", "storage.data.redundancy.type", "enable.openapi"} {
		if _, ok := got[k]; ok {
			t.Fatalf("unmanaged remote key %q must not be promoted into managed options", k)
		}
	}

	managed2 := map[string]interface{}{"comment": "v1", "description": "d"}
	api2 := map[string]interface{}{"comment": "v1"}
	got2 := reconcileManagedOptions(managed2, api2)
	if _, ok := got2["description"]; ok {
		t.Fatalf("managed key deleted by API must be dropped from state (drift), got description present")
	}
	if got2["comment"] != "v1" {
		t.Fatalf("managed key still in API must keep its value")
	}
	if len(got2) != 1 {
		t.Fatalf("reconciled managed set must have exactly 1 key, got %v", got2)
	}

	got3 := reconcileManagedOptions(map[string]interface{}{}, api)
	if len(got3) != 0 {
		t.Fatalf("empty managed state must yield empty reconcile, got %v", got3)
	}
}

func TestUnitDlfNextCatalogRemoteOptionsMirror(t *testing.T) {
	api := map[string]interface{}{
		"comment":                        "v1",
		"dlf.trashed-file-retained-days": 7,
		"storage.data.redundancy.type":   "LRS",
		"enable.openapi":                 false,
	}
	got := reconcileRemoteOptions(api)
	if got["comment"] != "v1" {
		t.Fatalf("comment = %v, want v1", got["comment"])
	}
	if got["dlf.trashed-file-retained-days"] != "7" {
		t.Fatalf("server default int 7 must stringify to \"7\", got %v", got["dlf.trashed-file-retained-days"])
	}
	if got["storage.data.redundancy.type"] != "LRS" {
		t.Fatalf("server default LRS wrong/missing, got %v", got["storage.data.redundancy.type"])
	}
	if got["enable.openapi"] != "false" {
		t.Fatalf("bool false must stringify to \"false\", got %v", got["enable.openapi"])
	}
	if len(got) != 4 {
		t.Fatalf("remote_options must mirror ALL 4 remote keys, got %d: %v", len(got), got)
	}
	if g := reconcileRemoteOptions(map[string]interface{}{}); len(g) != 0 {
		t.Fatalf("empty API must yield empty remote_options, got %v", g)
	}
}
