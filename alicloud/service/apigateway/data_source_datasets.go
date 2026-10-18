package apigateway

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"regexp"
	"sort"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const datasetsTypeName = "alicloud_api_gateway_datasets"

// PageSizeLarge mirrors the SDK v2 page size used by the legacy data source; the
// DescribeDatasetList endpoint caps PageSize, and 50 keeps pagination short.
const datasetsPageSize = 50

var (
	_ datasource.DataSource              = &datasetsDataSource{}
	_ datasource.DataSourceWithConfigure = &datasetsDataSource{}
)

func NewDatasetsDataSource() datasource.DataSource {
	return &datasetsDataSource{}
}

type datasetsDataSource struct {
	fwadapt.DataSourceBase
}

type datasetsDataSourceModel struct {
	Id         types.String `tfsdk:"id"`
	Ids        types.List   `tfsdk:"ids"`
	NameRegex  types.String `tfsdk:"name_regex"`
	OutputFile types.String `tfsdk:"output_file"`
	Datasets   types.List   `tfsdk:"datasets"`
}

// datasetItemAttrTypes is the object type of one element of the datasets list.
// It must line up, attribute name for attribute name, with the ListNestedAttribute
// below, or State.Set would reject the values.
var datasetItemAttrTypes = map[string]attr.Type{
	"id":            types.StringType,
	"dataset_id":    types.StringType,
	"dataset_name":  types.StringType,
	"dataset_type":  types.StringType,
	"description":   types.StringType,
	"create_time":   types.StringType,
	"modified_time": types.StringType,
}

func (d *datasetsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = datasetsTypeName
}

func (d *datasetsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides the list of API Gateway datasets and can filter by id and name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The data source ID, a stable hash of the returned dataset IDs.",
			},
			"ids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of dataset IDs. When set, only those datasets are returned.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex used to filter datasets by name.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "Write the result to a file in addition to state.",
			},
			// ListAttribute with an ObjectType element, not ListNestedAttribute:
			// the provider is muxed behind protocol v5, and v5 cannot express a
			// NestedAttributeObject's Attributes. An ObjectType element type
			// routes through AttributeType (cty), which v5 supports. The element
			// shape must line up with datasetItemAttrTypes and the ObjectValue
			// built in Read.
			"datasets": schema.ListAttribute{
				Computed:    true,
				ElementType: types.ObjectType{AttrTypes: datasetItemAttrTypes},
				Description: "The list of API Gateway datasets.",
			},
		},
	}
}

func (d *datasetsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config datasetsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := d.Client()

	idsFilter := map[string]bool{}
	for _, v := range config.Ids.Elements() {
		if s, ok := v.(types.String); ok && s.ValueString() != "" {
			idsFilter[s.ValueString()] = true
		}
	}

	var nameRegex *regexp.Regexp
	if !config.NameRegex.IsNull() && config.NameRegex.ValueString() != "" {
		r, err := regexp.Compile(config.NameRegex.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: compiling name_regex", datasetsTypeName),
				err.Error(),
			)
			return
		}
		nameRegex = r
	}

	objects, err := listDatasets(client, idsFilter, nameRegex)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeDatasetList", datasetsTypeName),
			err.Error(),
		)
		return
	}

	var datasetElems []attr.Value
	var idsElems []attr.Value
	for _, obj := range objects {
		id := fmt.Sprint(obj["DatasetId"])
		idsElems = append(idsElems, types.StringValue(id))
		objValue, diags := types.ObjectValue(datasetItemAttrTypes, map[string]attr.Value{
			"id":            types.StringValue(id),
			"dataset_id":    strVal(obj["DatasetId"], id),
			"dataset_name":  strVal(obj["DatasetName"], ""),
			"dataset_type":  strVal(obj["DatasetType"], ""),
			"description":   strVal(obj["Description"], ""),
			"create_time":   strVal(obj["CreatedTime"], ""),
			"modified_time": strVal(obj["ModifiedTime"], ""),
		})
		resp.Diagnostics.Append(diags...)
		datasetElems = append(datasetElems, objValue)
	}

	datasetsList, diags := types.ListValue(types.ObjectType{AttrTypes: datasetItemAttrTypes}, datasetElems)
	resp.Diagnostics.Append(diags...)
	idsList, diags := types.ListValue(types.StringType, idsElems)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Ids = idsList
	config.Datasets = datasetsList
	config.Id = types.StringValue(datasetIDHash(idsElems))

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// listDatasets pages through DescribeDatasetList, applying the id and name_regex
// filters. Pagination ends when the page covers TotalCount, or when every id in
// the filter has been collected.
func listDatasets(client *connectivity.AliyunClient, idsFilter map[string]bool, nameRegex *regexp.Regexp) ([]map[string]interface{}, error) {
	var objects []map[string]interface{}
	pageNumber := 1
	for {
		query := map[string]interface{}{
			"PageNumber": pageNumber,
			"PageSize":   datasetsPageSize,
		}
		response, err := client.RpcPost(datasetAPIProduct, datasetAPIVersion, "DescribeDatasetList", query, map[string]interface{}{}, true)
		if err != nil {
			return nil, errs.WrapErrorf(err, "Reading %s: calling DescribeDatasetList", datasetsTypeName)
		}
		raw, _ := response["DatasetInfoList"].([]interface{})
		for _, item := range raw {
			obj, _ := item.(map[string]interface{})
			if obj == nil {
				continue
			}
			id := fmt.Sprint(obj["DatasetId"])
			if len(idsFilter) > 0 && !idsFilter[id] {
				continue
			}
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(obj["DatasetName"])) {
				continue
			}
			objects = append(objects, obj)
		}

		totalCount := 0
		switch tc := response["TotalCount"].(type) {
		case float64:
			totalCount = int(tc)
		case json.Number:
			n, _ := tc.Int64()
			totalCount = int(n)
		}
		if pageNumber*datasetsPageSize >= totalCount {
			break
		}
		if len(idsFilter) > 0 && len(objects) >= len(idsFilter) {
			break
		}
		pageNumber++
	}
	return objects, nil
}

// datasetIDHash is a stable hash of the returned ids, used as the data source ID.
func datasetIDHash(elems []attr.Value) string {
	ids := make([]string, 0, len(elems))
	for _, v := range elems {
		if s, ok := v.(types.String); ok {
			ids = append(ids, s.ValueString())
		}
	}
	sort.Strings(ids)
	h := fnv.New32a()
	for _, id := range ids {
		_, _ = h.Write([]byte(id))
		_, _ = h.Write([]byte{0})
	}
	return fmt.Sprintf("%d", h.Sum32())
}
