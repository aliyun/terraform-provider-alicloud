// Package dlfnext implements the DlfNext (Data Lake Formation, 2025-03-10)
// Framework resources: catalog, database and table.
//
// DlfNext exposes a RESTful (ROA) API. None of its CRUD operations are backed
// by a generated Go SDK in this provider, so every call goes through the
// generic connectivity.AliyunClient ROA helpers (RoaPost/RoaGet/RoaDelete).
// Path parameters are substituted by the caller; the request body is a
// map[string]interface{} that the Darabonba ROA client JSON-encodes. The
// endpoint is resolved at runtime via the shared location service under the
// lowercased product code "dlfnext", matching the Aliyun CLI.
package dlfnext

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	dlfNextProductCode = "DlfNext"
	dlfNextApiVersion  = "2025-03-10"
)

// --- Catalog paths ---------------------------------------------------------

func catalogPath(name string) string {
	return fmt.Sprintf("/dlf/v1/catalogs/%s", name)
}

const catalogsPath = "/dlf/v1/catalogs"

// --- Database paths --------------------------------------------------------

func databasePath(catalogId, name string) string {
	return fmt.Sprintf("/dlf/v1/%s/databases/%s", catalogId, name)
}

func databasesPath(catalogId string) string {
	return fmt.Sprintf("/dlf/v1/%s/databases", catalogId)
}

// --- Table paths -----------------------------------------------------------

func tablePath(catalogId, database, name string) string {
	return fmt.Sprintf("/dlf/v1/%s/databases/%s/tables/%s", catalogId, database, name)
}

func tablesPath(catalogId, database string) string {
	return fmt.Sprintf("/dlf/v1/%s/databases/%s/tables", catalogId, database)
}

// isNotFound reports whether err is a DlfNext resource-not-found error. The
// API returns 404 (caught by errs.NotFoundError) for missing catalogs,
// databases and tables; the substring fallback covers wrapped forms whose
// code/message DlfNext may evolve without a fixed catalog.
func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	if errs.NotFoundError(err) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "notfound") ||
		strings.Contains(msg, "not found") ||
		strings.Contains(msg, "notexist") ||
		strings.Contains(msg, "not exist") ||
		strings.Contains(msg, "no such") ||
		strings.Contains(msg, "\"404\"") || strings.HasSuffix(msg, ": 404")
}

// --- Catalog ---------------------------------------------------------------

// getCatalog reads a catalog by its name (the {catalog} path parameter) and
// returns the catalog object, or nil when it no longer exists.
func getCatalog(client *connectivity.AliyunClient, name string) (map[string]interface{}, error) {
	path := catalogPath(name)
	response, err := client.RoaGet(dlfNextProductCode, dlfNextApiVersion, path, nil, nil, nil)
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, errs.WrapErrorf(err, "reading DlfNext Catalog %q via GET %s", name, path)
	}
	return response, nil
}

// isBeingCreated reports whether err is the 403 "catalog is being created, can
// not be dropped" error DropCatalog returns while a catalog is still
// provisioning (status NEW/INITIALIZING). Callers should wait for RUNNING and
// retry the drop.
func isBeingCreated(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "is being created") ||
		strings.Contains(msg, "can not be dropped") ||
		strings.Contains(msg, "please wait for the catalog to be ready")
}

// waitCatalogStatus polls getCatalog until the catalog reaches the target
// status (RUNNING for normal use) or a terminal failure, with a timeout. It
// tolerates the catalog being briefly invisible right after CreateCatalog.
func waitCatalogStatus(ctx context.Context, client *connectivity.AliyunClient, name, target string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var last string
	for {
		object, err := getCatalog(client, name)
		if err != nil {
			return err
		}
		if object != nil {
			last = stringFromResponse(object, "status")
			if last == target {
				return nil
			}
			if last == "INITIALIZE_FAILED" {
				return fmt.Errorf("catalog %q provisioning failed (status INITIALIZE_FAILED)", name)
			}
		}
		if time.Now().After(deadline) {
			if last == "" {
				return fmt.Errorf("timeout waiting for catalog %q to become visible via GetCatalog", name)
			}
			return fmt.Errorf("timeout waiting for catalog %q to reach status %s (last %q)", name, target, last)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
		}
	}
}

// --- Database --------------------------------------------------------------

// getDatabase reads a database by catalogId and database name. Returns nil
// when the database no longer exists.
func getDatabase(client *connectivity.AliyunClient, catalogId, name string) (map[string]interface{}, error) {
	path := databasePath(catalogId, name)
	response, err := client.RoaGet(dlfNextProductCode, dlfNextApiVersion, path, nil, nil, nil)
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, errs.WrapErrorf(err, "reading DlfNext Database %q via GET %s", catalogId+":"+name, path)
	}
	return response, nil
}

// --- Table -----------------------------------------------------------------

// getTable reads a table by catalogId, database and table name. Returns nil
// when the table no longer exists.
func getTable(client *connectivity.AliyunClient, catalogId, database, name string) (map[string]interface{}, error) {
	path := tablePath(catalogId, database, name)
	response, err := client.RoaGet(dlfNextProductCode, dlfNextApiVersion, path, nil, nil, nil)
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, errs.WrapErrorf(err, "reading DlfNext Table %q via GET %s", catalogId+":"+database+":"+name, path)
	}
	return response, nil
}

// --- Shared helpers --------------------------------------------------------

// optionsToAPI converts a map[string]string options block into the
// map[string]interface{} the DlfNext body expects. nil/empty maps produce nil
// so the field is omitted from the JSON body.
func optionsToAPI(m map[string]string) map[string]interface{} {
	if len(m) == 0 {
		return nil
	}
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// optionsFromAPI converts the options object of a Get response back to a
// map[string]string. Nested map values are stringified; absent options yield
// nil.
func optionsFromAPI(v interface{}) map[string]string {
	if v == nil {
		return nil
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, val := range m {
		out[k] = fmt.Sprint(val)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// diffOptions computes the Alter{Catalog,Database} body (updates/removals)
// from the old and new options maps. When both are empty no change is needed
// and both returned maps are nil.
func diffOptions(old, new_ map[string]string) (updates map[string]interface{}, removals []string) {
	updates = nil
	removals = nil
	for k, v := range new_ {
		prev, had := old[k]
		if !had || prev != v {
			if updates == nil {
				updates = make(map[string]interface{})
			}
			updates[k] = v
		}
	}
	for k := range old {
		if _, has := new_[k]; !has {
			removals = append(removals, k)
		}
	}
	return
}

// --- Table column type parser ---------------------------------------------

// parseTypeString parses a SQL-like column type expression into the recursive
// FullDataType object DlfNext expects in the CreateTable body:
//
//	STRING                         -> {type: STRING}
//	ARRAY<INT>                     -> {type: ARRAY, element: {type: INT}}
//	MULTISET<STRING>               -> {type: MULTISET, element: {type: STRING}}
//	MAP<STRING,BIGINT>             -> {type: MAP, key: {type: STRING}, value: {type: BIGINT}}
//	ROW<id INT, name STRING>       -> {type: ROW, fields: [{name: id, type: {type: INT}}, ...]}
//
// Primitive names are upper-cased so users may write them in any case.
func parseTypeString(s string) (map[string]interface{}, error) {
	p := &typeParser{input: s}
	t, err := p.parseType()
	if err != nil {
		return nil, err
	}
	p.skipSpace()
	if p.pos < len(p.input) {
		return nil, fmt.Errorf("unexpected trailing characters in column type %q at position %d", s, p.pos)
	}
	return t, nil
}

type typeParser struct {
	input string
	pos   int
}

func (p *typeParser) skipSpace() {
	for p.pos < len(p.input) && (p.input[p.pos] == ' ' || p.input[p.pos] == '\t') {
		p.pos++
	}
}

func (p *typeParser) expect(c byte) error {
	p.skipSpace()
	if p.pos >= len(p.input) || p.input[p.pos] != c {
		return fmt.Errorf("expected %q in column type at position %d", string(c), p.pos)
	}
	p.pos++
	return nil
}

func (p *typeParser) parseIdent() (string, error) {
	p.skipSpace()
	start := p.pos
	for p.pos < len(p.input) {
		c := p.input[p.pos]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
			p.pos++
			continue
		}
		break
	}
	if p.pos == start {
		return "", fmt.Errorf("expected an identifier in column type at position %d", p.pos)
	}
	return p.input[start:p.pos], nil
}

func (p *typeParser) parseType() (map[string]interface{}, error) {
	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	upper := strings.ToUpper(name)
	p.skipSpace()
	if p.pos >= len(p.input) || p.input[p.pos] != '<' {
		return map[string]interface{}{"type": upper}, nil
	}
	p.pos++ // consume '<'
	switch upper {
	case "ARRAY", "MULTISET":
		elem, err := p.parseType()
		if err != nil {
			return nil, err
		}
		if err := p.expect('>'); err != nil {
			return nil, err
		}
		return map[string]interface{}{"type": upper, "element": elem}, nil
	case "MAP":
		key, err := p.parseType()
		if err != nil {
			return nil, err
		}
		if err := p.expect(','); err != nil {
			return nil, err
		}
		val, err := p.parseType()
		if err != nil {
			return nil, err
		}
		if err := p.expect('>'); err != nil {
			return nil, err
		}
		return map[string]interface{}{"type": upper, "key": key, "value": val}, nil
	case "ROW":
		fields, err := p.parseRowFields()
		if err != nil {
			return nil, err
		}
		if err := p.expect('>'); err != nil {
			return nil, err
		}
		out := map[string]interface{}{"type": upper}
		if len(fields) > 0 {
			out["fields"] = fields
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unsupported parametric column type %q", upper)
	}
}

func (p *typeParser) parseRowFields() ([]map[string]interface{}, error) {
	var fields []map[string]interface{}
	for {
		p.skipSpace()
		if p.pos < len(p.input) && p.input[p.pos] == '>' {
			break
		}
		fname, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		ftype, err := p.parseType()
		if err != nil {
			return nil, err
		}
		fields = append(fields, map[string]interface{}{"name": fname, "type": ftype})
		p.skipSpace()
		if p.pos < len(p.input) && p.input[p.pos] == ',' {
			p.pos++
			continue
		}
		break
	}
	return fields, nil
}

// formatTypeString renders a FullDataType object (from a GetTable response)
// back into the SQL-like column type expression used by the schema.
func formatTypeString(t map[string]interface{}) string {
	if t == nil {
		return ""
	}
	name, _ := t["type"].(string)
	// The API appends " NOT NULL" to the type of primary-key columns
	// (e.g. "BIGINT NOT NULL"). The NOT NULL constraint is implied by
	// primary_keys and cannot be expressed in the user-facing type
	// expression, so strip it to keep the type round-tripping cleanly
	// (otherwise every refresh would drift config vs state and force a
	// needless replace).
	name = strings.TrimSuffix(name, " NOT NULL")
	switch name {
	case "ARRAY", "MULTISET":
		return fmt.Sprintf("%s<%s>", name, formatTypeString(asMap(t["element"])))
	case "MAP":
		return fmt.Sprintf("MAP<%s,%s>", formatTypeString(asMap(t["key"])), formatTypeString(asMap(t["value"])))
	case "ROW":
		var parts []string
		if fields, ok := t["fields"].([]interface{}); ok {
			for _, f := range fields {
				fm := asMap(f)
				fname, _ := fm["name"].(string)
				parts = append(parts, fmt.Sprintf("%s %s", fname, formatTypeString(asMap(fm["type"]))))
			}
		}
		return fmt.Sprintf("ROW<%s>", strings.Join(parts, ","))
	default:
		return name
	}
}

func asMap(v interface{}) map[string]interface{} {
	if m, ok := v.(map[string]interface{}); ok {
		return m
	}
	return nil
}

// stringSlice converts []interface{} (from a JSON response) to []string,
// stringifying each element. nil/empty input yields nil.
func stringSlice(v interface{}) []string {
	if v == nil {
		return nil
	}
	arr, ok := v.([]interface{})
	if !ok || len(arr) == 0 {
		return nil
	}
	out := make([]string, 0, len(arr))
	for _, e := range arr {
		out = append(out, fmt.Sprint(e))
	}
	return out
}

// stringFromResponse returns the string value of a response field, or "" when
// absent or non-string.
func stringFromResponse(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

// boolFromResponse returns the bool value of a response field, or false when
// absent or non-bool.
func boolFromResponse(m map[string]interface{}, key string) bool {
	if m == nil {
		return false
	}
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
}

// stringAttr maps a Go string to a types.String, turning the empty string into
// null so an absent response field reads as unset rather than "".
func stringAttr(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

// int64FromResponse returns the int64 value of a response field, or 0 when
// absent. JSON numbers decode as float64.
func int64FromResponse(m map[string]interface{}, key string) int64 {
	if m == nil {
		return 0
	}
	switch v := m[key].(type) {
	case float64:
		return int64(v)
	case int64:
		return v
	case int:
		return int64(v)
	case string:
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return n
		}
	}
	return 0
}

// --- Framework type conversions -------------------------------------------

// stringMapToTerraform builds a types.Map from a Go map[string]string. An
// empty/nil map maps to a null Map so the attribute reads as absent.
func stringMapToTerraform(m map[string]string) types.Map {
	if len(m) == 0 {
		return types.MapNull(types.StringType)
	}
	elems := make(map[string]attr.Value, len(m))
	for k, v := range m {
		elems[k] = types.StringValue(v)
	}
	return types.MapValueMust(types.StringType, elems)
}

// terraformMapToStringMap extracts a Go map[string]string from a types.Map of
// string elements. A null/unknown map yields nil.
func terraformMapToStringMap(m types.Map) map[string]string {
	if m.IsNull() || m.IsUnknown() {
		return nil
	}
	elements := m.Elements()
	out := make(map[string]string, len(elements))
	for k, v := range elements {
		if sv, ok := v.(types.String); ok {
			out[k] = sv.ValueString()
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// stringSliceToTerraform builds a types.List of strings from a Go []string.
// An empty/nil slice maps to a null List.
func stringSliceToTerraform(s []string) types.List {
	if len(s) == 0 {
		return types.ListNull(types.StringType)
	}
	elems := make([]attr.Value, len(s))
	for i, v := range s {
		elems[i] = types.StringValue(v)
	}
	return types.ListValueMust(types.StringType, elems)
}

// terraformListToStringSlice extracts a Go []string from a types.List of
// string elements. A null/unknown list yields nil.
func terraformListToStringSlice(l types.List) []string {
	if l.IsNull() || l.IsUnknown() {
		return nil
	}
	elements := l.Elements()
	out := make([]string, 0, len(elements))
	for _, v := range elements {
		if sv, ok := v.(types.String); ok {
			out = append(out, sv.ValueString())
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
