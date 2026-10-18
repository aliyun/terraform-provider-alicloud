package dataworks

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// fakeLineageFetch serves canned GetMetaTableLineage responses in sequence and
// records each request so the unit tests can assert the request shape the real
// Read builds — including the NextPrimaryKey cursor carried on later pages —
// without touching the cloud. This covers Read, pagination, parsing, the
// request shape and the failure path; the framework state-setting layer is
// exercised only by the environment-guarded acceptance test.
type fakeLineageFetch struct {
	pages []map[string]interface{}
	calls []map[string]interface{}
	err   error
}

func (f *fakeLineageFetch) fetch(action string, request map[string]interface{}) (map[string]interface{}, error) {
	// Record a copy so later mutation of the request (NextPrimaryKey) does
	// not retroactively edit the history.
	cp := make(map[string]interface{}, len(request))
	for k, v := range request {
		cp[k] = v
	}
	f.calls = append(f.calls, cp)
	if f.err != nil {
		return nil, f.err
	}
	if len(f.pages) == 0 {
		return map[string]interface{}{"Data": map[string]interface{}{}}, nil
	}
	resp := f.pages[0]
	f.pages = f.pages[1:]
	return resp, nil
}

func lineageResponse(entries []map[string]interface{}, hasNext bool, nextKey string) map[string]interface{} {
	// RpcPost decodes JSON arrays as []interface{} (each element a
	// map[string]interface{}), so the fixture must mirror that shape: a typed
	// []map[string]interface{} would fail the .([]interface{}) assertion in
	// readTableLineage and silently yield zero entries.
	list := make([]interface{}, len(entries))
	for i, e := range entries {
		list[i] = e
	}
	return map[string]interface{}{
		"Data": map[string]interface{}{
			"DataEntityList": list,
			"HasNext":        hasNext,
			"NextPrimaryKey": nextKey,
		},
	}
}

func newLineageModel(direction string, ids []string) tableLineageModel {
	m := tableLineageModel{
		Direction: types.StringValue(direction),
		PageSize:  types.Int64Value(100),
	}
	if ids != nil {
		m.Ids = ids
	}
	return m
}

func TestUnitReadTableLineageSinglePage(t *testing.T) {
	f := &fakeLineageFetch{
		pages: []map[string]interface{}{
			lineageResponse([]map[string]interface{}{
				{"TableGuid": "t-guid-1", "TableName": "t1", "DatabaseName": "db1", "CreateTimestamp": json.Number("1700000000")},
				{"TableGuid": "t-guid-2", "TableName": "t2", "DatabaseName": "db2", "CreateTimestamp": float64(1700000001)},
			}, false, ""),
		},
	}
	ds := &tableLineageDataSource{fetch: f.fetch}

	entries, ids, err := ds.readTableLineage(nil, newLineageModel("down", nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if !reflect.DeepEqual(ids, []string{"t-guid-1", "t-guid-2"}) {
		t.Fatalf("unexpected ids: %v", ids)
	}
	if entries[0].CreateTimestamp.ValueInt64() != 1700000000 {
		t.Errorf("json.Number create_timestamp not parsed, got %d", entries[0].CreateTimestamp.ValueInt64())
	}
	if entries[1].CreateTimestamp.ValueInt64() != 1700000001 {
		t.Errorf("float64 create_timestamp not parsed, got %d", entries[1].CreateTimestamp.ValueInt64())
	}
	if len(f.calls) != 1 {
		t.Errorf("expected 1 fetch call, got %d", len(f.calls))
	}
	if f.calls[0]["Direction"] != "down" || f.calls[0]["PageSize"] != 100 {
		t.Errorf("unexpected request shape: %v", f.calls[0])
	}
	if _, ok := f.calls[0]["NextPrimaryKey"]; ok {
		t.Errorf("first page must not carry NextPrimaryKey")
	}
}

func TestUnitReadTableLineagePaginationCarriesCursor(t *testing.T) {
	f := &fakeLineageFetch{
		pages: []map[string]interface{}{
			lineageResponse([]map[string]interface{}{
				{"TableGuid": "p1-1", "TableName": "a", "DatabaseName": "db", "CreateTimestamp": int64(1)},
			}, true, "cursor-page-2"),
			lineageResponse([]map[string]interface{}{
				{"TableGuid": "p2-1", "TableName": "b", "DatabaseName": "db", "CreateTimestamp": int64(2)},
			}, false, ""),
		},
	}
	ds := &tableLineageDataSource{fetch: f.fetch}

	entries, ids, err := ds.readTableLineage(nil, newLineageModel("up", nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(ids, []string{"p1-1", "p2-1"}) {
		t.Fatalf("unexpected ids across pages: %v", ids)
	}
	if len(f.calls) != 2 {
		t.Fatalf("expected 2 fetch calls, got %d", len(f.calls))
	}
	if f.calls[1]["NextPrimaryKey"] != "cursor-page-2" {
		t.Errorf("second page must carry the cursor from page 1, got %v", f.calls[1]["NextPrimaryKey"])
	}
	if entries[1].TableName.ValueString() != "b" {
		t.Errorf("second entry not parsed, got %v", entries[1])
	}
}

func TestUnitReadTableLineageEmptyResult(t *testing.T) {
	f := &fakeLineageFetch{
		pages: []map[string]interface{}{
			lineageResponse(nil, false, ""),
		},
	}
	ds := &tableLineageDataSource{fetch: f.fetch}

	entries, ids, err := ds.readTableLineage(nil, newLineageModel("up", nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 || len(ids) != 0 {
		t.Fatalf("expected empty result, got entries=%v ids=%v", entries, ids)
	}
}

func TestUnitReadTableLineageFetchError(t *testing.T) {
	f := &fakeLineageFetch{err: errors.New("ResourceNotAuthorized.Api")}
	ds := &tableLineageDataSource{fetch: f.fetch}

	if _, _, err := ds.readTableLineage(nil, newLineageModel("up", nil)); err == nil {
		t.Fatal("expected an error when fetch fails, got nil")
	}
}

func TestUnitReadTableLineageMissingDataEntityList(t *testing.T) {
	// A malformed response with no Data.DataEntityList must surface an error,
	// not silently return an empty result.
	f := &fakeLineageFetch{
		pages: []map[string]interface{}{
			{"Data": map[string]interface{}{}},
		},
	}
	ds := &tableLineageDataSource{fetch: f.fetch}

	if _, _, err := ds.readTableLineage(nil, newLineageModel("up", nil)); err == nil {
		t.Fatal("expected an error when Data.DataEntityList is missing, got nil")
	}
}

func TestUnitReadTableLineageIdsFilter(t *testing.T) {
	f := &fakeLineageFetch{
		pages: []map[string]interface{}{
			lineageResponse([]map[string]interface{}{
				{"TableGuid": "keep", "TableName": "a", "DatabaseName": "db", "CreateTimestamp": int64(1)},
				{"TableGuid": "drop", "TableName": "b", "DatabaseName": "db", "CreateTimestamp": int64(2)},
			}, false, ""),
		},
	}
	ds := &tableLineageDataSource{fetch: f.fetch}

	entries, ids, err := ds.readTableLineage(nil, newLineageModel("up", []string{"keep"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(ids, []string{"keep"}) {
		t.Fatalf("ids filter did not narrow the result: %v", ids)
	}
	if len(entries) != 1 || entries[0].Id.ValueString() != "keep" {
		t.Fatalf("unexpected filtered entry: %v", entries)
	}
}

func TestUnitReadTableLineageCompositeIdFallback(t *testing.T) {
	// EMR-style entries may not carry a TableGuid; the element id must fall
	// back to a composite key so it stays unique within the data source.
	f := &fakeLineageFetch{
		pages: []map[string]interface{}{
			lineageResponse([]map[string]interface{}{
				{"TableGuid": "", "TableName": "emr_t", "DatabaseName": "emr_db", "CreateTimestamp": int64(3)},
			}, false, ""),
		},
	}
	ds := &tableLineageDataSource{fetch: f.fetch}

	entries, ids, err := ds.readTableLineage(nil, newLineageModel("down", nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "emr_t:emr_db:down"
	if len(ids) != 1 || ids[0] != want {
		t.Fatalf("expected composite id %q, got %v", want, ids)
	}
	if entries[0].Id.ValueString() != want {
		t.Errorf("entry id not the composite key, got %q", entries[0].Id.ValueString())
	}
}
