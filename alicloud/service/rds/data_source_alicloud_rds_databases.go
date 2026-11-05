package rds

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdkretry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const databasesTypeName = "alicloud_rds_databases"

const (
	databasesAction   = "DescribeDatabases"
	databasesPageSize = 100
	// databasesMaxPages bounds a pagination loop that never shrinks. The RDS
	// quotas are far below this, so hitting it means the endpoint ignored the
	// page parameters; returning the accumulated rows would silently be partial.
	databasesMaxPages = 1000
)

// databasesRetryableErrorCodes are the transient 400/408 codes of
// DescribeDatabases: the API asks the caller to retry while the database is in
// transition, while the request concurrency limit is reached, or after a SQL
// Server query timeout.
var databasesRetryableErrorCodes = []string{
	"SqlExcutionFailed",
	"ConcurrentLimit",
	"NetworkOrSqlTimeoutError",
}

var (
	_ datasource.DataSource              = &databasesDataSource{}
	_ datasource.DataSourceWithConfigure = &databasesDataSource{}
)

func NewDatabasesDataSource() datasource.DataSource {
	return &databasesDataSource{}
}

type databasesDataSource struct {
	fwadapt.DataSourceBase
}

type databasesModel struct {
	Id         types.String             `tfsdk:"id"`
	InstanceId types.String             `tfsdk:"instance_id"`
	DbName     types.String             `tfsdk:"db_name"`
	DbStatus   types.String             `tfsdk:"db_status"`
	NameRegex  types.String             `tfsdk:"name_regex"`
	Names      types.List               `tfsdk:"names"`
	Databases  []databasesDatabaseModel `tfsdk:"databases"`
}

type databasesDatabaseModel struct {
	DbName           types.String            `tfsdk:"db_name"`
	DbStatus         types.String            `tfsdk:"db_status"`
	CharacterSetName types.String            `tfsdk:"character_set_name"`
	DbDescription    types.String            `tfsdk:"db_description"`
	Engine           types.String            `tfsdk:"engine"`
	InstanceId       types.String            `tfsdk:"instance_id"`
	Accounts         []databasesAccountModel `tfsdk:"accounts"`
}

type databasesAccountModel struct {
	AccountName            types.String `tfsdk:"account_name"`
	AccountPrivilege       types.String `tfsdk:"account_privilege"`
	AccountPrivilegeDetail types.String `tfsdk:"account_privilege_detail"`
}

var databasesAccountObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"account_name":             types.StringType,
	"account_privilege":        types.StringType,
	"account_privilege_detail": types.StringType,
}}

var databasesDatabaseObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"db_name":            types.StringType,
	"db_status":          types.StringType,
	"character_set_name": types.StringType,
	"db_description":     types.StringType,
	"engine":             types.StringType,
	"instance_id":        types.StringType,
	"accounts":           types.ListType{ElemType: databasesAccountObjectType},
}}

// databasesNameRegexValidator rejects a `name_regex` that cannot compile, so
// the mismatch surfaces at plan time and no request is sent.
type databasesNameRegexValidator struct{}

func (v databasesNameRegexValidator) Description(context.Context) string {
	return "The value must be a valid regular expression."
}

func (v databasesNameRegexValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v databasesNameRegexValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if _, err := regexp.Compile(req.ConfigValue.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid Name Regex", err.Error())
	}
}

func (d *databasesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = databasesTypeName
}

func (d *databasesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a list of databases of an ApsaraDB RDS instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, derived from `instance_id` and the filters so that it stays stable across refreshes.",
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the RDS instance.",
			},
			"db_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the database. The value is sent to the server as the `DBName` filter.",
			},
			"db_status": schema.StringAttribute{
				Optional:    true,
				Description: "The status of the database to filter. Valid values: `Creating`, `Running` and `Deleting`. The value is sent to the server as the `DBStatus` filter and is also applied on the provider side, so a server that ignores the parameter cannot widen the result.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string to filter results by database name. The filter is applied on the provider side after all pages are read.",
				Validators: []validator.String{
					databasesNameRegexValidator{},
				},
			},
			"names": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of database names that match the configured filters.",
			},
			"databases": schema.ListAttribute{
				Computed:    true,
				ElementType: databasesDatabaseObjectType,
				Description: "A list of databases that match the configured filters. Each element contains the core fields returned by DescribeDatabases. Engine-specific fields such as the SQL Server `AdvancedInfo`, `BasicInfo` and `RuntimeInfo`, the PostgreSQL `Collate`, `Ctype`, `ConnLimit` and `Tablespace`, and `DuckDBEnabled` are not exported.",
			},
		},
	}
}

func (d *databasesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model databasesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var nameRegex *regexp.Regexp
	if v := model.NameRegex.ValueString(); !model.NameRegex.IsNull() && v != "" {
		compiled, err := regexp.Compile(v)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: invalid name_regex", databasesTypeName),
				err.Error(),
			)
			return
		}
		nameRegex = compiled
	}

	items, err := d.describeAllDatabases(ctx, &model)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", databasesTypeName, databasesAction),
			err.Error(),
		)
		return
	}

	names := make([]string, 0, len(items))
	databases := make([]databasesDatabaseModel, 0, len(items))
	statusFilter := ""
	if !model.DbStatus.IsNull() {
		statusFilter = model.DbStatus.ValueString()
	}
	for _, item := range items {
		name := databasesString(item, "DBName")
		if statusFilter != "" && databasesString(item, "DBStatus") != statusFilter {
			continue
		}
		if nameRegex != nil && !nameRegex.MatchString(name) {
			continue
		}
		names = append(names, name)
		databases = append(databases, databasesDatabaseModel{
			DbName:           types.StringValue(name),
			DbStatus:         types.StringValue(databasesString(item, "DBStatus")),
			CharacterSetName: types.StringValue(databasesString(item, "CharacterSetName")),
			DbDescription:    types.StringValue(databasesString(item, "DBDescription")),
			Engine:           types.StringValue(databasesString(item, "Engine")),
			InstanceId:       types.StringValue(databasesString(item, "DBInstanceId")),
			Accounts:         databasesAccounts(item),
		})
	}

	namesList, diags := types.ListValueFrom(ctx, types.StringType, names)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	model.Id = types.StringValue(databasesId(&model))
	model.Names = namesList
	model.Databases = databases
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// describeAllDatabases reads every page. DescribeDatabases has no top-level
// total count, so the loop stops at the first page shorter than PageSize; an
// exact multiple of PageSize therefore costs one extra, empty page.
func (d *databasesDataSource) describeAllDatabases(ctx context.Context, model *databasesModel) ([]map[string]interface{}, error) {
	return databasesCollectPages(databasesMaxPages, databasesPageSize, func(pageNumber int) ([]map[string]interface{}, error) {
		return d.describeDatabasesPage(ctx, model, pageNumber)
	})
}

func databasesCollectPages(maxPages, pageSize int, fetch func(pageNumber int) ([]map[string]interface{}, error)) ([]map[string]interface{}, error) {
	items := make([]map[string]interface{}, 0)
	for pageNumber := 1; pageNumber <= maxPages; pageNumber++ {
		pageItems, err := fetch(pageNumber)
		if err != nil {
			return nil, err
		}
		items = append(items, pageItems...)
		if len(pageItems) < pageSize {
			return items, nil
		}
	}
	return nil, fmt.Errorf("%s returned a full page for %d consecutive requests (PageSize=%d), so pagination did not terminate; the result was not returned because it would be incomplete", databasesAction, maxPages, pageSize)
}

func (d *databasesDataSource) describeDatabasesPage(ctx context.Context, model *databasesModel, pageNumber int) ([]map[string]interface{}, error) {
	client := d.Meta()
	request := map[string]interface{}{
		"DBInstanceId": model.InstanceId.ValueString(),
		"PageNumber":   pageNumber,
		"PageSize":     databasesPageSize,
	}
	if v := model.DbName.ValueString(); !model.DbName.IsNull() && v != "" {
		request["DBName"] = v
	}
	if v := model.DbStatus.ValueString(); !model.DbStatus.IsNull() && v != "" {
		request["DBStatus"] = v
	}

	var response map[string]interface{}
	// Cancellation can return before the RPC callback finishes.
	var responseMu sync.Mutex
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := sdkretry.RetryContext(ctx, 5*time.Minute, func() *sdkretry.RetryError {
		result, rpcErr := client.RpcPost("Rds", "2014-08-15", databasesAction, nil, request, true)
		responseMu.Lock()
		response = result
		responseMu.Unlock()
		if rpcErr != nil {
			if helper.NeedRetry(rpcErr) || IsExpectedErrors(rpcErr, databasesRetryableErrorCodes) {
				wait()
				return sdkretry.RetryableError(rpcErr)
			}
			return sdkretry.NonRetryableError(rpcErr)
		}
		return nil
	})
	responseMu.Lock()
	defer responseMu.Unlock()
	logs.AddDebug(databasesAction, response, request)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, databasesTypeName, databasesAction, errs.AlibabaCloudSdkGoERROR)
	}
	return databasesItems(response), nil
}

// databasesItems parses the Databases.Database array. A response without the
// key or with an empty array yields an empty slice; it is not an error.
func databasesItems(response map[string]interface{}) []map[string]interface{} {
	items := make([]map[string]interface{}, 0)
	if response == nil {
		return items
	}
	raw, err := jsonpath.Get("$.Databases.Database", response)
	if err != nil {
		return items
	}
	list, ok := raw.([]interface{})
	if !ok {
		return items
	}
	for _, v := range list {
		if item, ok := v.(map[string]interface{}); ok {
			items = append(items, item)
		}
	}
	return items
}

func databasesAccounts(item map[string]interface{}) []databasesAccountModel {
	accounts := make([]databasesAccountModel, 0)
	raw, err := jsonpath.Get("$.Accounts.AccountPrivilegeInfo", item)
	if err != nil {
		return accounts
	}
	list, ok := raw.([]interface{})
	if !ok {
		return accounts
	}
	for _, v := range list {
		account, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		accounts = append(accounts, databasesAccountModel{
			AccountName:            types.StringValue(databasesString(account, "Account")),
			AccountPrivilege:       types.StringValue(databasesString(account, "AccountPrivilege")),
			AccountPrivilegeDetail: types.StringValue(databasesString(account, "AccountPrivilegeDetail")),
		})
	}
	return accounts
}

func databasesString(item map[string]interface{}, key string) string {
	if v, ok := item[key].(string); ok {
		return v
	}
	return ""
}

func databasesId(model *databasesModel) string {
	return dataResourceIdHash([]string{
		model.InstanceId.ValueString(),
		model.DbName.ValueString(),
		model.DbStatus.ValueString(),
		model.NameRegex.ValueString(),
	})
}
