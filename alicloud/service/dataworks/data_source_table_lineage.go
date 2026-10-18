package dataworks

import (
	"context"
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const tableLineageTypeName = "alicloud_data_works_table_lineage"

// API contract (dataworks-public 2020-05-18, GetMetaTableLineage):
//
//	request: Direction(up/down), TableGuid?, TableName?, DatabaseName?,
//	         DataSourceType?(odps/emr), ClusterId?, PageSize, NextPrimaryKey(cursor)
//	response: Data.DataEntityList[]{TableGuid,TableName,DatabaseName,CreateTimestamp},
//	          Data.HasNext(bool), Data.NextPrimaryKey(string cursor)
const (
	tableLineageProduct   = "dataworks-public"
	tableLineageAPIVer    = "2020-05-18"
	tableLineageAction    = "GetMetaTableLineage"
	tableLineageDefaultPS = 100
	tableLineageMaxPS     = 100
)

var (
	_ datasource.DataSource              = &tableLineageDataSource{}
	_ datasource.DataSourceWithConfigure = &tableLineageDataSource{}
)

func NewTableLineageDataSource() datasource.DataSource {
	return &tableLineageDataSource{}
}

// tableLineageDataSource is the Framework-native data source for
// alicloud_data_works_table_lineage. fetch is the RPC entry point; nil means
// the default RpcPost-based fetcher with throttling retry is used. It is
// overridable in tests so Read, pagination, parsing and failure paths can be
// exercised without the cloud — the fetch receives exactly the request map
// the real Read builds (Direction, PageSize and the NextPrimaryKey cursor on
// later pages), so it covers the request shape too, not just a helper.
type tableLineageDataSource struct {
	fwadapt.DataSourceBase
	fetch func(action string, request map[string]interface{}) (map[string]interface{}, error)
}

type tableLineageModel struct {
	Direction      types.String        `tfsdk:"direction"`
	TableGuid      types.String        `tfsdk:"table_guid"`
	TableName      types.String        `tfsdk:"table_name"`
	DatabaseName   types.String        `tfsdk:"database_name"`
	DataSourceType types.String        `tfsdk:"data_source_type"`
	ClusterId      types.String        `tfsdk:"cluster_id"`
	PageSize       types.Int64         `tfsdk:"page_size"`
	Ids            []string            `tfsdk:"ids"`
	OutputFile     types.String        `tfsdk:"output_file"`
	Id             types.String        `tfsdk:"id"`
	Lineages       []tableLineageEntry `tfsdk:"lineages"`
}

type tableLineageEntry struct {
	Id              types.String `tfsdk:"id"`
	TableGuid       types.String `tfsdk:"table_guid"`
	TableName       types.String `tfsdk:"table_name"`
	DatabaseName    types.String `tfsdk:"database_name"`
	CreateTimestamp types.Int64  `tfsdk:"create_timestamp"`
}

func (d *tableLineageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = tableLineageTypeName
}

func (d *tableLineageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source queries the upstream or downstream table lineage of a table in DataWorks via the GetMetaTableLineage API (product dataworks-public, version 2020-05-18).",
		Attributes: map[string]schema.Attribute{
			"direction": schema.StringAttribute{
				Required:    true,
				Description: "The direction of the lineage. Valid values: `up` (upstream) and `down` (downstream).",
			},
			"table_guid": schema.StringAttribute{
				Optional:    true,
				Description: "The unique identifier of the table. If specified, the lineage of the table identified by this GUID is returned.",
			},
			"table_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the table. Used together with `database_name` and `data_source_type` to identify a table when `table_guid` is unknown.",
			},
			"database_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the database where the table resides.",
			},
			"data_source_type": schema.StringAttribute{
				Optional:    true,
				Description: "The data source type. Valid values: `odps` and `emr`. Required for EMR scenarios together with `cluster_id`.",
			},
			"cluster_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the EMR cluster. This parameter is required for EMR scenarios.",
			},
			"page_size": schema.Int64Attribute{
				Optional:    true,
				Description: "The number of entries per page. Default value: `100`. Maximum value: `100`. The data source auto-paginates all results using `NextPrimaryKey`.",
			},
			"ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "A list of lineage ids used to filter the returned lineages. Only lineages whose `id` matches one of the given ids are returned. Also exported as the list of matching lineage ids.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "File name where to write the result after running `terraform plan`.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The data source id, a hash of the returned lineage ids.",
			},
			"lineages": schema.ListNestedAttribute{
				Computed:    true,
				Description: "A list of table lineage entries.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The id of the lineage entry, same as `table_guid` when present.",
						},
						"table_guid": schema.StringAttribute{
							Computed:    true,
							Description: "The unique identifier of the table.",
						},
						"table_name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the table.",
						},
						"database_name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the database.",
						},
						"create_timestamp": schema.Int64Attribute{
							Computed:    true,
							Description: "The creation timestamp of the table.",
						},
					},
				},
			},
		},
	}
}

// readTableLineage is the cloud-independent core of Read: it builds the request
// from the model, paginates via HasNext + NextPrimaryKey, parses each lineage
// entry, applies the optional `ids` filter, and computes the data source id.
// It is exercised directly by the unit tests with a fake fetch so Read,
// pagination, parsing, the request shape and failure paths are covered
// without the cloud — not just a helper.
func (d *tableLineageDataSource) readTableLineage(ctx context.Context, model tableLineageModel) ([]tableLineageEntry, []string, error) {
	fetch := d.fetch
	if fetch == nil {
		fetch = d.defaultFetch
	}

	direction := model.Direction.ValueString()
	pageSize := int(model.PageSize.ValueInt64())
	if pageSize < 1 {
		pageSize = tableLineageDefaultPS
	}
	if pageSize > tableLineageMaxPS {
		pageSize = tableLineageMaxPS
	}

	idsFilter := make(map[string]struct{}, len(model.Ids))
	for _, id := range model.Ids {
		idsFilter[id] = struct{}{}
	}

	request := map[string]interface{}{
		"Direction": direction,
		"PageSize":  pageSize,
	}
	if v := model.TableGuid.ValueString(); v != "" {
		request["TableGuid"] = v
	}
	if v := model.TableName.ValueString(); v != "" {
		request["TableName"] = v
	}
	if v := model.DatabaseName.ValueString(); v != "" {
		request["DatabaseName"] = v
	}
	if v := model.DataSourceType.ValueString(); v != "" {
		request["DataSourceType"] = v
	}
	if v := model.ClusterId.ValueString(); v != "" {
		request["ClusterId"] = v
	}

	var entries []tableLineageEntry
	var ids []string
	for {
		response, err := fetch(tableLineageAction, request)
		if err != nil {
			return nil, nil, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, tableLineageTypeName, tableLineageAction, errs.AlibabaCloudSdkGoERROR)
		}

		rawList, perr := jsonpath.Get("$.Data.DataEntityList", response)
		if perr != nil {
			return nil, nil, errs.WrapErrorf(perr, errs.FailedGetAttributeMsg, tableLineageAction, "$.Data.DataEntityList", response)
		}
		list, _ := rawList.([]interface{})

		for _, v := range list {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			tableGuid := asString(item["TableGuid"])
			tableName := asString(item["TableName"])
			databaseName := asString(item["DatabaseName"])
			// CreateTimestamp may arrive as json.Number, float64 or int64
			// depending on the SDK decode path; formatInt tolerates all.
			createTimestamp := formatInt(item["CreateTimestamp"])

			id := tableGuid
			if id == "" {
				// Fall back to a composite key when the API does not return a
				// TableGuid (e.g. EMR tables), so the element id stays unique
				// within the data source.
				id = fmt.Sprintf("%s:%s:%s", tableName, databaseName, direction)
			}
			if len(idsFilter) > 0 {
				if _, ok := idsFilter[id]; !ok {
					continue
				}
			}
			entries = append(entries, tableLineageEntry{
				Id:              types.StringValue(id),
				TableGuid:       types.StringValue(tableGuid),
				TableName:       types.StringValue(tableName),
				DatabaseName:    types.StringValue(databaseName),
				CreateTimestamp: types.Int64Value(int64(createTimestamp)),
			})
			ids = append(ids, id)
		}

		// GetMetaTableLineage paginates with HasNext + NextPrimaryKey (cursor
		// pagination), not PageNumber. Stop when HasNext is false/missing,
		// NextPrimaryKey is empty, or the current page returned no items.
		hasNext := false
		if v, perr := jsonpath.Get("$.Data.HasNext", response); perr == nil {
			if hb, ok := v.(bool); ok {
				hasNext = hb
			}
		}
		nextPrimaryKey := ""
		if v, perr := jsonpath.Get("$.Data.NextPrimaryKey", response); perr == nil {
			if s, ok := v.(string); ok {
				nextPrimaryKey = s
			}
		}
		if !hasNext || nextPrimaryKey == "" || len(list) == 0 {
			break
		}
		request["NextPrimaryKey"] = nextPrimaryKey
	}

	return entries, ids, nil
}

// defaultFetch is the real RPC entry point: RpcPost against dataworks-public
// with throttling retry, mirroring the SDKv2 data sources' behaviour.
func (d *tableLineageDataSource) defaultFetch(action string, request map[string]interface{}) (map[string]interface{}, error) {
	client := d.Client()
	if client == nil {
		return nil, fmt.Errorf("provider client is not configured for %s", tableLineageTypeName)
	}
	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost(tableLineageProduct, tableLineageAPIVer, action, nil, request, true)
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
	return response, err
}

func (d *tableLineageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model tableLineageModel
	diags := req.Config.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Schema-level enums are validated here instead of via
	// terraform-plugin-framework-validators so the data source does not pull in
	// a module the repository does not otherwise build against.
	direction := model.Direction.ValueString()
	if direction != "up" && direction != "down" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Invalid %s direction", tableLineageTypeName),
			fmt.Sprintf("direction must be one of `up`, `down`, got %q.", direction),
		)
	}
	if v := model.DataSourceType.ValueString(); v != "" && v != "odps" && v != "emr" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Invalid %s data_source_type", tableLineageTypeName),
			fmt.Sprintf("data_source_type must be one of `odps`, `emr`, got %q.", v),
		)
	}
	if !model.PageSize.IsNull() && !model.PageSize.IsUnknown() {
		ps := model.PageSize.ValueInt64()
		if ps < 1 || ps > tableLineageMaxPS {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Invalid %s page_size", tableLineageTypeName),
				fmt.Sprintf("page_size must be between 1 and %d, got %d.", tableLineageMaxPS, ps),
			)
		}
	}
	if resp.Diagnostics.HasError() {
		return
	}

	entries, ids, err := d.readTableLineage(ctx, model)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Reading %s: calling %s", tableLineageTypeName, tableLineageAction), err.Error())
		return
	}

	model.Id = types.StringValue(dataResourceIdHash(ids))
	model.Ids = ids
	model.Lineages = entries

	if output := model.OutputFile.ValueString(); output != "" {
		if werr := writeToFile(output, entries); werr != nil {
			resp.Diagnostics.AddWarning(fmt.Sprintf("Writing %s output_file", tableLineageTypeName), werr.Error())
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
