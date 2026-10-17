package paiplugin

import (
	"context"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const algorithmDataSourceTypeName = "alicloud_pai_plugin_algorithm"

var _ datasource.DataSource = &algorithmDataSource{}

func NewAlgorithmDataSource() datasource.DataSource {
	return &algorithmDataSource{}
}

type algorithmDataSource struct {
	fwadapt.DataSourceBase
}

// algorithmModel flattens the algorithm collection into parallel arrays so the
// framework state stays free of nested-object lists (which carry a key-order
// hazard without compile feedback). ListAlgorithms returns only Id+Name, so
// the detail fields (description, infer/train user config map) are populated
// only when the caller pins a single algorithm via algorithm_id, by following
// up GetAlgorithm.
type algorithmModel struct {
	Id                 types.String `tfsdk:"id"`
	AlgorithmId        types.String `tfsdk:"algorithm_id"`
	AlgorithmName      types.String `tfsdk:"algorithm_name"`
	AlgorithmIdOut     types.String `tfsdk:"algorithm_id_out"`
	AlgorithmNameOut   types.String `tfsdk:"algorithm_name_out"`
	Description        types.String `tfsdk:"description"`
	InferUserConfigMap types.String `tfsdk:"infer_user_config_map"`
	TrainUserConfigMap types.String `tfsdk:"train_user_config_map"`
	Ids                types.List   `tfsdk:"ids"`
	Names              types.List   `tfsdk:"names"`
}

func (d *algorithmDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = algorithmDataSourceTypeName
}

func (d *algorithmDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides PAI Plugin algorithms. With `algorithm_id` set, returns the single algorithm's full detail; otherwise lists all algorithms (id and name only).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source.",
			},
			"algorithm_id": schema.StringAttribute{
				Optional:    true,
				Description: "The algorithm id. When set, the data source reads the single algorithm via GetAlgorithm and populates the detail fields.",
			},
			"algorithm_name": schema.StringAttribute{
				Optional:    true,
				Description: "The algorithm name. When set without algorithm_id, ListAlgorithms is filtered by name.",
			},
			"algorithm_id_out": schema.StringAttribute{
				Computed:    true,
				Description: "The algorithm id of the first matching algorithm.",
			},
			"algorithm_name_out": schema.StringAttribute{
				Computed:    true,
				Description: "The algorithm name of the first matching algorithm.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "The description. Populated only when algorithm_id is set (ListAlgorithms does not return it).",
			},
			"infer_user_config_map": schema.StringAttribute{
				Computed:    true,
				Description: "The infer user config map (JSON string). Populated only when algorithm_id is set.",
			},
			"train_user_config_map": schema.StringAttribute{
				Computed:    true,
				Description: "The train user config map (JSON string). Populated only when algorithm_id is set.",
			},
			"ids": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "The ids of all matching algorithms, in the order returned by ListAlgorithms.",
			},
			"names": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "The names of all matching algorithms, in the order returned by ListAlgorithms.",
			},
		},
	}
}

func (d *algorithmDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config algorithmModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := d.Client()

	// When a single algorithm id is given, read its full detail via GetAlgorithm
	// so description and the user config maps are available; otherwise enumerate
	// via ListAlgorithms, which only carries Id and Name.
	if id := config.AlgorithmId.ValueString(); id != "" {
		object, err := getAlgorithm(client, id)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: calling GetAlgorithm", algorithmDataSourceTypeName),
				errs.WrapErrorf(err, "Reading %s: calling GetAlgorithm", algorithmDataSourceTypeName).Error(),
			)
			return
		}
		if object == nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: calling GetAlgorithm", algorithmDataSourceTypeName),
				fmt.Sprintf("The algorithm %q was not found", id),
			)
			return
		}
		name := str(object, "Name")
		state := algorithmModel{
			Id:                 types.StringValue(algorithmDataSourceTypeName),
			AlgorithmId:        types.StringValue(id),
			AlgorithmName:      strValue(object, "Name"),
			AlgorithmIdOut:     types.StringValue(id),
			AlgorithmNameOut:   strValue(object, "Name"),
			Description:        strValue(object, "Description"),
			InferUserConfigMap: strValue(object, "InferUserConfigMap"),
			TrainUserConfigMap: strValue(object, "TrainUserConfigMap"),
			Ids:                stringList([]string{id}),
			Names:              stringList([]string{name}),
		}
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}

	name := config.AlgorithmName.ValueString()
	response, err := listAlgorithms(client, name)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListAlgorithms", algorithmDataSourceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling ListAlgorithms", algorithmDataSourceTypeName).Error(),
		)
		return
	}

	data := dataObject(response)
	algorithms := objectSlice(data, "Algorithms")

	ids := make([]string, 0, len(algorithms))
	names := make([]string, 0, len(algorithms))
	for _, a := range algorithms {
		ids = append(ids, str(a, "Id"))
		names = append(names, str(a, "Name"))
	}

	state := algorithmModel{
		Id:    types.StringValue(algorithmDataSourceTypeName),
		Ids:   stringList(ids),
		Names: stringList(names),
	}
	if len(algorithms) > 0 {
		first := algorithms[0]
		state.AlgorithmIdOut = strValue(first, "Id")
		state.AlgorithmNameOut = strValue(first, "Name")
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// getAlgorithm reads a single algorithm via GetAlgorithm.
func getAlgorithm(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	response, err := roaGet(client, pathAlgorithms+"/"+id, nil)
	if err != nil {
		if isNotFound(err, response) {
			return nil, nil
		}
		return nil, err
	}
	return dataObject(response), nil
}

// listAlgorithms enumerates algorithms, optionally filtered by name.
func listAlgorithms(client *connectivity.AliyunClient, name string) (map[string]interface{}, error) {
	query := queryFrom("Name", name)
	return roaGet(client, pathAlgorithms, query)
}

// objectSlice extracts a []map[string]interface{} array field from a response
// object, tolerating non-object elements by skipping them.
func objectSlice(obj map[string]interface{}, key string) []map[string]interface{} {
	if obj == nil {
		return nil
	}
	arr, ok := obj[key].([]interface{})
	if !ok {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(arr))
	for _, e := range arr {
		if m, ok := e.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out
}

// stringList builds a framework String list from a slice, returning null for nil.
func stringList(values []string) types.List {
	if values == nil {
		return types.ListNull(types.StringType)
	}
	elements := make([]attr.Value, 0, len(values))
	for _, v := range values {
		elements = append(elements, types.StringValue(v))
	}
	list, diags := types.ListValue(types.StringType, elements)
	_ = diags
	return list
}
