package rds_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/rds"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

// databasesStateModel mirrors the data source model for decoding the state
// after a unit-test Read; the implementation's model is unexported.
type databasesStateModel struct {
	Id         types.String                  `tfsdk:"id"`
	InstanceId types.String                  `tfsdk:"instance_id"`
	DbName     types.String                  `tfsdk:"db_name"`
	DbStatus   types.String                  `tfsdk:"db_status"`
	NameRegex  types.String                  `tfsdk:"name_regex"`
	Names      types.List                    `tfsdk:"names"`
	Databases  []databasesDatabaseStateModel `tfsdk:"databases"`
}

type databasesDatabaseStateModel struct {
	DbName           types.String                 `tfsdk:"db_name"`
	DbStatus         types.String                 `tfsdk:"db_status"`
	CharacterSetName types.String                 `tfsdk:"character_set_name"`
	DbDescription    types.String                 `tfsdk:"db_description"`
	Engine           types.String                 `tfsdk:"engine"`
	InstanceId       types.String                 `tfsdk:"instance_id"`
	Accounts         []databasesAccountStateModel `tfsdk:"accounts"`
}

type databasesAccountStateModel struct {
	AccountName            types.String `tfsdk:"account_name"`
	AccountPrivilege       types.String `tfsdk:"account_privilege"`
	AccountPrivilegeDetail types.String `tfsdk:"account_privilege_detail"`
}

// databasesNamesFromState unpacks the names state list into plain strings.
func databasesNamesFromState(t *testing.T, value types.List) []string {
	t.Helper()
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	raw, err := value.ToTerraformValue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	var entries []tftypes.Value
	if err := raw.As(&entries); err != nil {
		t.Fatal(err)
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		var name string
		if err := entry.As(&name); err != nil {
			t.Fatal(err)
		}
		names = append(names, name)
	}
	return names
}

// databasesRead drives one data source Read against the mock client and
// decodes the resulting state. It returns a nil state when the Read reported
// an error diagnostic.
func databasesRead(ctx context.Context, t *testing.T, client *connectivity.AliyunClient, config map[string]string) (*datasource.ReadResponse, *databasesStateModel) {
	t.Helper()
	ds := rds.NewDatabasesDataSource()

	configurable, ok := ds.(datasource.DataSourceWithConfigure)
	if !ok {
		t.Fatal("data source does not support Configure")
	}
	var cresp datasource.ConfigureResponse
	configurable.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: client}, &cresp)
	if cresp.Diagnostics.HasError() {
		t.Fatalf("Configure diagnostics: %+v", cresp.Diagnostics)
	}

	var sresp datasource.SchemaResponse
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &sresp)
	if sresp.Diagnostics.HasError() {
		t.Fatalf("Schema diagnostics: %+v", sresp.Diagnostics)
	}
	objType, ok := sresp.Schema.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatalf("unexpected data source schema type %T", sresp.Schema.Type().TerraformType(context.Background()))
	}

	raw := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, attrType := range objType.AttributeTypes {
		raw[name] = tftypes.NewValue(attrType, nil)
	}
	for name, value := range config {
		attrType, ok := objType.AttributeTypes[name]
		if !ok {
			t.Fatalf("config attribute %q is not in the schema", name)
		}
		raw[name] = tftypes.NewValue(attrType, value)
	}

	resp := &datasource.ReadResponse{
		State: tfsdk.State{
			Schema: sresp.Schema,
			Raw:    tftypes.NewValue(objType, nil),
		},
	}
	ds.Read(ctx, datasource.ReadRequest{
		Config: tfsdk.Config{Schema: sresp.Schema, Raw: tftypes.NewValue(objType, raw)},
	}, resp)
	if resp.Diagnostics.HasError() {
		return resp, nil
	}

	var state databasesStateModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("state decode diagnostics: %+v", diags)
	}
	return resp, &state
}

func databasesMockItem(name string) map[string]interface{} {
	return map[string]interface{}{
		"DBName":           name,
		"DBStatus":         "Running",
		"CharacterSetName": "utf8",
		"DBDescription":    "description of " + name,
		"Engine":           "MySQL",
		"DBInstanceId":     "rm-test",
		"Accounts": map[string]interface{}{
			"AccountPrivilegeInfo": []map[string]interface{}{
				{
					"Account":                "test_account",
					"AccountPrivilege":       "ReadWrite",
					"AccountPrivilegeDetail": "SELECT,INSERT",
				},
			},
		},
	}
}

func databasesWriteResponse(t *testing.T, w http.ResponseWriter, items []map[string]interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if items == nil {
		items = []map[string]interface{}{}
	}
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"RequestId": "mock-request-id",
		"Databases": map[string]interface{}{"Database": items},
	}); err != nil {
		t.Error(err)
	}
}

// Two pages are merged in order, the request carries the documented PageSize,
// and the loop stops on the short second page.
func TestUnitRdsDatabasesReadPagination(t *testing.T) {
	var calls int32
	var pages []string
	var pagesMu sync.Mutex
	client := rdsPerformanceTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		assert.Equal(t, "DescribeDatabases", r.Form.Get("Action"))
		assert.Equal(t, "100", r.Form.Get("PageSize"))
		pageNumber := r.Form.Get("PageNumber")
		pagesMu.Lock()
		pages = append(pages, pageNumber)
		pagesMu.Unlock()

		items := make([]map[string]interface{}, 0)
		switch pageNumber {
		case "1":
			for i := 0; i < 100; i++ {
				items = append(items, databasesMockItem(fmt.Sprintf("db_%03d", i)))
			}
		case "2":
			for i := 100; i < 102; i++ {
				items = append(items, databasesMockItem(fmt.Sprintf("db_%03d", i)))
			}
		default:
			assert.Failf(t, "unexpected page", "PageNumber=%s", pageNumber)
		}
		databasesWriteResponse(t, w, items)
	})

	resp, state := databasesRead(context.Background(), t, client, map[string]string{
		"instance_id": "rm-test",
	})
	assert.False(t, resp.Diagnostics.HasError(), "diagnostics: %+v", resp.Diagnostics)
	if !assert.NotNil(t, state) {
		return
	}
	assert.Equal(t, int32(2), atomic.LoadInt32(&calls))
	pagesMu.Lock()
	assert.Equal(t, []string{"1", "2"}, pages)
	pagesMu.Unlock()

	if !assert.Len(t, state.Databases, 102) {
		return
	}
	names := databasesNamesFromState(t, state.Names)
	if assert.Len(t, names, 102) {
		seen := make(map[string]struct{}, len(names))
		for i, name := range names {
			if _, ok := seen[name]; ok {
				t.Fatalf("duplicate database name %q at index %d", name, i)
			}
			seen[name] = struct{}{}
			assert.Equal(t, fmt.Sprintf("db_%03d", i), name)
		}
	}
	assert.Equal(t, "db_000", state.Databases[0].DbName.ValueString())
	assert.Equal(t, "db_101", state.Databases[101].DbName.ValueString())
	assert.NotEmpty(t, state.Id.ValueString())
}

// A result set that is exactly one page long ends with one extra empty page.
func TestUnitRdsDatabasesReadExactPageBoundary(t *testing.T) {
	var calls int32
	client := rdsPerformanceTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		items := make([]map[string]interface{}, 0)
		if r.Form.Get("PageNumber") == "1" {
			for i := 0; i < 100; i++ {
				items = append(items, databasesMockItem(fmt.Sprintf("db_%03d", i)))
			}
		}
		databasesWriteResponse(t, w, items)
	})

	resp, state := databasesRead(context.Background(), t, client, map[string]string{
		"instance_id": "rm-test",
	})
	assert.False(t, resp.Diagnostics.HasError(), "diagnostics: %+v", resp.Diagnostics)
	if !assert.NotNil(t, state) {
		return
	}
	assert.Equal(t, int32(2), atomic.LoadInt32(&calls), "an exact full page must be followed by one empty page")
	assert.Len(t, state.Databases, 100)
	assert.Len(t, databasesNamesFromState(t, state.Names), 100)
}

// The server-side filters are sent on every page, the client-side regex filters
// the merged result, and the account details are mapped into the output.
func TestUnitRdsDatabasesReadFilters(t *testing.T) {
	client := rdsPerformanceTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DescribeDatabases", r.Form.Get("Action"))
		assert.Equal(t, "rm-test", r.Form.Get("DBInstanceId"))
		assert.Equal(t, "alpha", r.Form.Get("DBName"))
		assert.Equal(t, "Running", r.Form.Get("DBStatus"))
		assert.Equal(t, "1", r.Form.Get("PageNumber"))
		databasesWriteResponse(t, w, []map[string]interface{}{
			databasesMockItem("alpha"),
			databasesMockItem("beta"),
			databasesMockItem("gamma"),
		})
	})

	resp, state := databasesRead(context.Background(), t, client, map[string]string{
		"instance_id": "rm-test",
		"db_name":     "alpha",
		"db_status":   "Running",
		"name_regex":  "^a",
	})
	assert.False(t, resp.Diagnostics.HasError(), "diagnostics: %+v", resp.Diagnostics)
	if !assert.NotNil(t, state) {
		return
	}
	assert.Equal(t, []string{"alpha"}, databasesNamesFromState(t, state.Names))
	if !assert.Len(t, state.Databases, 1) {
		return
	}
	database := state.Databases[0]
	assert.Equal(t, "alpha", database.DbName.ValueString())
	assert.Equal(t, "Running", database.DbStatus.ValueString())
	assert.Equal(t, "utf8", database.CharacterSetName.ValueString())
	assert.Equal(t, "description of alpha", database.DbDescription.ValueString())
	assert.Equal(t, "MySQL", database.Engine.ValueString())
	assert.Equal(t, "rm-test", database.InstanceId.ValueString())
	if assert.Len(t, database.Accounts, 1) {
		assert.Equal(t, "test_account", database.Accounts[0].AccountName.ValueString())
		assert.Equal(t, "ReadWrite", database.Accounts[0].AccountPrivilege.ValueString())
		assert.Equal(t, "SELECT,INSERT", database.Accounts[0].AccountPrivilegeDetail.ValueString())
	}
}

// The real endpoint can ignore DBStatus, so the provider must enforce the
// requested status on the returned items as well.
func TestUnitRdsDatabasesReadFiltersByStatus(t *testing.T) {
	client := rdsPerformanceTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Running", r.Form.Get("DBStatus"))
		creating := databasesMockItem("beta")
		creating["DBStatus"] = "Creating"
		deleting := databasesMockItem("gamma")
		deleting["DBStatus"] = "Deleting"
		databasesWriteResponse(t, w, []map[string]interface{}{
			databasesMockItem("alpha"),
			creating,
			deleting,
		})
	})
	resp, state := databasesRead(context.Background(), t, client, map[string]string{
		"instance_id": "rm-test",
		"db_status":   "Running",
	})
	assert.False(t, resp.Diagnostics.HasError(), "diagnostics: %+v", resp.Diagnostics)
	if assert.NotNil(t, state) {
		assert.Equal(t, []string{"alpha"}, databasesNamesFromState(t, state.Names))
	}
}

// An instance without databases is a successful empty result, not an error.
func TestUnitRdsDatabasesReadEmpty(t *testing.T) {
	client := rdsPerformanceTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		databasesWriteResponse(t, w, nil)
	})

	resp, state := databasesRead(context.Background(), t, client, map[string]string{
		"instance_id": "rm-test",
	})
	assert.False(t, resp.Diagnostics.HasError(), "diagnostics: %+v", resp.Diagnostics)
	if !assert.NotNil(t, state) {
		return
	}
	assert.Empty(t, state.Databases)
	assert.Empty(t, databasesNamesFromState(t, state.Names))
	assert.NotEmpty(t, state.Id.ValueString(), "an empty result must still carry a stable id")
}

// A non-existent instance must surface the API error instead of an empty list.
func TestUnitRdsDatabasesReadNotFound(t *testing.T) {
	var calls int32
	client := rdsPerformanceTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		rdsPerformanceWriteError(t, w, 404, "InvalidDBInstanceId.NotFound")
	})
	resp, state := databasesRead(context.Background(), t, client, map[string]string{
		"instance_id": "rm-missing",
	})
	assert.True(t, resp.Diagnostics.HasError())
	assert.Nil(t, state)
	assert.Equal(t, int32(1), atomic.LoadInt32(&calls), "a non-retryable InvalidDBInstanceId.NotFound must not be retried")
	found := false
	for _, e := range resp.Diagnostics.Errors() {
		if strings.Contains(e.Summary(), "InvalidDBInstanceId.NotFound") || strings.Contains(e.Detail(), "InvalidDBInstanceId.NotFound") {
			found = true
		}
	}
	assert.True(t, found, "diagnostics: %+v", resp.Diagnostics)
}

// SqlExcutionFailed is the documented "database is in transition" code; the
// Read must send the request again instead of failing immediately.
func TestUnitRdsDatabasesReadRetriesTransition(t *testing.T) {
	var calls int32
	client := rdsPerformanceTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			rdsPerformanceWriteError(t, w, 400, "SqlExcutionFailed")
			return
		}
		assert.Equal(t, "DescribeDatabases", r.Form.Get("Action"))
		databasesWriteResponse(t, w, []map[string]interface{}{databasesMockItem("alpha")})
	})
	resp, state := databasesRead(context.Background(), t, client, map[string]string{
		"instance_id": "rm-test",
	})
	assert.False(t, resp.Diagnostics.HasError(), "diagnostics: %+v", resp.Diagnostics)
	assert.Equal(t, int32(2), atomic.LoadInt32(&calls), "SqlExcutionFailed must be retried once and then succeed")
	if assert.NotNil(t, state) {
		assert.Equal(t, []string{"alpha"}, databasesNamesFromState(t, state.Names))
	}
}

// databasesNameRegexValidators returns the plan-stage validators declared on
// `name_regex`, read through the public Schema method of the factory.
func databasesNameRegexValidators(t *testing.T) []validator.String {
	t.Helper()
	response := datasource.SchemaResponse{}
	rds.NewDatabasesDataSource().Schema(context.Background(), datasource.SchemaRequest{}, &response)
	attribute, ok := response.Schema.Attributes["name_regex"].(schema.StringAttribute)
	if !assert.True(t, ok, "name_regex must be a string attribute") {
		return nil
	}
	if !assert.NotEmpty(t, attribute.Validators, "name_regex must declare a plan-stage validator") {
		return nil
	}
	return attribute.Validators
}

func TestUnitRdsDatabasesNameRegexValidator(t *testing.T) {
	for _, tc := range []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "anchored name", value: "^db_01$"},
		{name: "character class", value: `^db_[0-9]+`},
		{name: "empty value is not validated", value: ""},
		{name: "unclosed group", value: "^(db_01", wantErr: true},
		{name: "unclosed class", value: "[a-z", wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var diagnostics diag.Diagnostics
			for _, nameRegexValidator := range databasesNameRegexValidators(t) {
				response := validator.StringResponse{}
				nameRegexValidator.ValidateString(context.Background(), validator.StringRequest{
					Path:        path.Root("name_regex"),
					ConfigValue: types.StringValue(tc.value),
				}, &response)
				diagnostics = append(diagnostics, response.Diagnostics...)
			}
			if !tc.wantErr {
				assert.False(t, diagnostics.HasError(), "diagnostics: %+v", diagnostics)
				return
			}
			if !assert.True(t, diagnostics.HasError(), "expected the validator to reject %q", tc.value) {
				return
			}
			withPath, ok := diagnostics.Errors()[0].(diag.DiagnosticWithPath)
			if assert.True(t, ok, "diagnostic must carry an attribute path") {
				assert.True(t, withPath.Path().Equal(path.Root("name_regex")))
			}
		})
	}
}

func TestAccAliCloudRdsDatabasesDataSource(t *testing.T) {
	name := fmt.Sprintf("tf-testAcc%d", tfacctest.RandInt())
	resourceId := "data.alicloud_rds_databases.default"
	baseConfig := rdsDbInstancePerformanceBaseConfig(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: baseConfig + databasesResourcesConfig() + databasesDataBlock(""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "databases.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "names.#", "2"),
					checkDatabasesStringSet("names.", "", "tf_testacc0", "tf_testacc1"),
					checkDatabasesStringSet("databases.", ".db_name", "tf_testacc0", "tf_testacc1"),
					checkDatabasesStringSet("databases.", ".db_description", "terraform test database 0", "terraform test database 1"),
					resource.TestMatchResourceAttr(resourceId, "databases.0.db_status", regexp.MustCompile("^(Creating|Running)$")),
					resource.TestCheckResourceAttr(resourceId, "databases.0.character_set_name", "utf8"),
					resource.TestCheckResourceAttr(resourceId, "databases.0.engine", "MySQL"),
					resource.TestCheckResourceAttrPair(resourceId, "databases.0.instance_id", "alicloud_db_instance.default", "id"),
				),
			},
			{
				Config: baseConfig + databasesResourcesConfig() + databasesDataBlock("  db_name = \"tf_testacc0\"\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "databases.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "names.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "names.0", "tf_testacc0"),
					resource.TestCheckResourceAttr(resourceId, "databases.0.db_name", "tf_testacc0"),
				),
			},
			{
				Config: baseConfig + databasesResourcesConfig() + databasesDataBlock("  name_regex = \"^tf_testacc1$\"\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "databases.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "names.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "names.0", "tf_testacc1"),
					resource.TestCheckResourceAttr(resourceId, "databases.0.db_name", "tf_testacc1"),
				),
			},
			{
				// No database is being deleted during the test, so a Deleting
				// filter must return nothing; if the server ignores DBStatus the
				// two existing databases would leak through.
				Config: baseConfig + databasesResourcesConfig() + databasesDataBlock("  db_status = \"Deleting\"\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "databases.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "names.#", "0"),
				),
			},
		},
	})
}

func TestAccAliCloudRdsDatabasesDataSourceNotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "alicloud_rds_databases" "default" {
  instance_id = "rm-bp1%s"
}`, tfacctest.RandString(14)),
				ExpectError: regexp.MustCompile("InvalidDBInstanceId\\.NotFound"),
			},
		},
	})
}

func databasesResourcesConfig() string {
	return `
resource "alicloud_db_database" "default" {
  count          = 2
  instance_id    = alicloud_db_instance.default.id
  data_base_name = "tf_testacc${count.index}"
  character_set  = "utf8"
  description    = "terraform test database ${count.index}"
}
`
}

func databasesDataBlock(filters string) string {
	return fmt.Sprintf(`
data "alicloud_rds_databases" "default" {
  # The databases are created after the instance. The splat reads the instance
  # id from every database instance, so the data source waits for all of them
  # instead of reading alongside their creation.
  instance_id = element(alicloud_db_database.default[*].instance_id, 0)
%s}`, filters)
}

// checkDatabasesStringSet compares the values under prefix/suffix as a set,
// because the API does not guarantee the order of the returned databases.
func checkDatabasesStringSet(prefix, suffix string, want ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["data.alicloud_rds_databases.default"]
		if !ok {
			return fmt.Errorf("data source data.alicloud_rds_databases.default was not found in the state")
		}
		actual := make([]string, 0)
		for key, value := range rs.Primary.Attributes {
			if !strings.HasPrefix(key, prefix) || !strings.HasSuffix(key, suffix) || strings.HasSuffix(key, ".#") {
				continue
			}
			actual = append(actual, value)
		}
		sort.Strings(actual)
		expected := append([]string(nil), want...)
		sort.Strings(expected)
		if !reflect.DeepEqual(actual, expected) {
			return fmt.Errorf("unexpected %s*%s values: got %v, want %v", prefix, suffix, actual, expected)
		}
		return nil
	}
}
