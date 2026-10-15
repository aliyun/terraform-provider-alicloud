package pai

import (
	"context"
	"fmt"
	"regexp"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const nodeOutputsTypeName = "alicloud_pai_node_outputs"

var (
	_ datasource.DataSource              = (*nodeOutputsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*nodeOutputsDataSource)(nil)
)

func NewNodeOutputsDataSource() datasource.DataSource {
	return &nodeOutputsDataSource{}
}

type nodeOutputsDataSource struct {
	fwadapt.DataSourceBase
}

// nodeOutputItemModel is one element of the outputs list. It mirrors the
// fields ListNodeOutputs returns per element.
type nodeOutputItemModel struct {
	OutputId     types.String `tfsdk:"output_id"`
	OutputIndex  types.String `tfsdk:"output_index"`
	NodeName     types.String `tfsdk:"node_name"`
	AlgoName     types.String `tfsdk:"algo_name"`
	DisplayName  types.String `tfsdk:"display_name"`
	Type         types.String `tfsdk:"type"`
	LocationType types.String `tfsdk:"location_type"`
	Value        types.String `tfsdk:"value"`
}

type nodeOutputsModel struct {
	Id            types.String          `tfsdk:"id"`
	ExperimentId  types.String          `tfsdk:"experiment_id"`
	NodeId        types.String          `tfsdk:"node_id"`
	OutputIdRegex types.String          `tfsdk:"output_id_regex"`
	TypeFilter    types.String          `tfsdk:"type"`
	Ids           types.List            `tfsdk:"ids"`
	Outputs       []nodeOutputItemModel `tfsdk:"outputs"`
}

func (d *nodeOutputsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = nodeOutputsTypeName
}

func (d *nodeOutputsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides the list of outputs of a node in a PAI Studio experiment.\n\n" +
			"This data source calls the PaiStudio ListNodeOutputs API to enumerate every output of a node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, formatted as `<experiment_id>:<node_id>`.",
			},
			"experiment_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the PAI Studio experiment.",
			},
			"node_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the node in the experiment.",
			},
			"output_id_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string to filter results by OutputId. The filter is applied client-side.",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "Filter results by node output Type, e.g. `Model`. The filter is applied client-side.",
			},
			"ids": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "A list of OutputIds of the matched node outputs.",
			},
			"outputs": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"output_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the node output.",
						},
						"output_index": schema.StringAttribute{
							Computed:    true,
							Description: "The index of the node output.",
						},
						"node_name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the node.",
						},
						"algo_name": schema.StringAttribute{
							Computed:    true,
							Description: "The algorithm name of the node.",
						},
						"display_name": schema.StringAttribute{
							Computed:    true,
							Description: "The display name of the node output.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the node output, e.g. `Model`, `MaxComputeTable`.",
						},
						"location_type": schema.StringAttribute{
							Computed:    true,
							Description: "The location type of the node output, e.g. `MaxComputeTable`.",
						},
						"value": schema.StringAttribute{
							Computed:    true,
							Description: "The value payload of the node output, serialized as a JSON string. Its structure varies by `type`.",
						},
					},
				},
				Description: "The list of node outputs.",
			},
		},
	}
}

func (d *nodeOutputsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state nodeOutputsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	experimentId := state.ExperimentId.ValueString()
	nodeId := state.NodeId.ValueString()

	action := fmt.Sprintf("/api/v1/experiments/%s/nodes/%s/outputs", experimentId, nodeId)
	response, err := d.Client().RoaGet(paiStudioProduct, paiStudioVersion, action, nil, nil, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListNodeOutputs", nodeOutputsTypeName),
			err.Error(),
		)
		return
	}

	var idRegex *regexp.Regexp
	if !state.OutputIdRegex.IsNull() && state.OutputIdRegex.ValueString() != "" {
		idRegex, err = regexp.Compile(state.OutputIdRegex.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: invalid output_id_regex", nodeOutputsTypeName),
				err.Error(),
			)
			return
		}
	}
	typeFilter := ""
	if !state.TypeFilter.IsNull() {
		typeFilter = state.TypeFilter.ValueString()
	}

	rawItems, _ := response["Outputs"].([]interface{})
	var items []nodeOutputItemModel
	var ids []string
	for _, raw := range rawItems {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		outputId := asString(item["OutputId"])
		outputType := asString(item["Type"])
		if idRegex != nil && !idRegex.MatchString(outputId) {
			continue
		}
		if typeFilter != "" && outputType != typeFilter {
			continue
		}
		items = append(items, nodeOutputItemModel{
			OutputId:     types.StringValue(outputId),
			OutputIndex:  types.StringValue(asString(item["OutputIndex"])),
			NodeName:     types.StringValue(asString(item["NodeName"])),
			AlgoName:     types.StringValue(asString(item["AlgoName"])),
			DisplayName:  types.StringValue(asString(item["DisplayName"])),
			Type:         types.StringValue(outputType),
			LocationType: types.StringValue(asString(item["LocationType"])),
			Value:        types.StringValue(valueToJSON(item["Value"])),
		})
		ids = append(ids, outputId)
	}

	state.Id = types.StringValue(fmt.Sprintf("%s:%s", experimentId, nodeId))
	state.Outputs = items
	state.Ids = types.ListValueMust(types.StringType, stringSliceToListValues(ids))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// stringSliceToListValues converts a []string to a slice of attr.Value for
// constructing a types.List.
func stringSliceToListValues(in []string) []attr.Value {
	out := make([]attr.Value, len(in))
	for i, v := range in {
		out[i] = types.StringValue(v)
	}
	return out
}
