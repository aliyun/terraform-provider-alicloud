package pai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	paiStudioProduct   = "PAIStudio"
	paiStudioVersion   = "2021-02-02"
	nodeOutputTypeName = "alicloud_pai_node_output"
)

var (
	_ datasource.DataSource              = (*nodeOutputDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*nodeOutputDataSource)(nil)
)

func NewNodeOutputDataSource() datasource.DataSource {
	return &nodeOutputDataSource{}
}

type nodeOutputDataSource struct {
	fwadapt.DataSourceBase
}

// nodeOutputModel is the state model for a single node output. The Value field
// is the PAI Studio node output payload, whose shape varies by Type (Model,
// MaxComputeTable, …); it is surfaced verbatim as a JSON string so callers can
// decode the structure their experiment produced.
type nodeOutputModel struct {
	Id           types.String `tfsdk:"id"`
	ExperimentId types.String `tfsdk:"experiment_id"`
	NodeId       types.String `tfsdk:"node_id"`
	OutputId     types.String `tfsdk:"output_id"`
	OutputIndex  types.String `tfsdk:"output_index"`
	NodeName     types.String `tfsdk:"node_name"`
	AlgoName     types.String `tfsdk:"algo_name"`
	DisplayName  types.String `tfsdk:"display_name"`
	Type         types.String `tfsdk:"type"`
	LocationType types.String `tfsdk:"location_type"`
	Value        types.String `tfsdk:"value"`
}

func (d *nodeOutputDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = nodeOutputTypeName
}

func (d *nodeOutputDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a single output of a node in a PAI Studio experiment.\n\n" +
			"This data source calls the PaiStudio GetNodeOutput API to read one node output identified by experiment, node and output id.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the node output, formatted as `<experiment_id>:<node_id>:<output_id>`.",
			},
			"experiment_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the PAI Studio experiment.",
			},
			"node_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the node in the experiment.",
			},
			"output_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the node output.",
			},
			"output_index": schema.StringAttribute{
				Optional:    true,
				Description: "The index of the node output. When the node has multiple outputs with the same OutputId, specify the index to pick one.",
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
	}
}

func (d *nodeOutputDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state nodeOutputModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	experimentId := state.ExperimentId.ValueString()
	nodeId := state.NodeId.ValueString()
	outputId := state.OutputId.ValueString()

	action := fmt.Sprintf("/api/v1/experiments/%s/nodes/%s/outputs/%s", experimentId, nodeId, outputId)
	query := make(map[string]*string)
	if !state.OutputIndex.IsNull() && state.OutputIndex.ValueString() != "" {
		index := state.OutputIndex.ValueString()
		query["OutputIndex"] = &index
	}

	response, err := d.Client().RoaGet(paiStudioProduct, paiStudioVersion, action, query, nil, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetNodeOutput", nodeOutputTypeName),
			err.Error(),
		)
		return
	}

	state.Id = types.StringValue(fmt.Sprintf("%s:%s:%s", experimentId, nodeId, outputId))
	state.NodeName = types.StringValue(asString(response["NodeName"]))
	state.AlgoName = types.StringValue(asString(response["AlgoName"]))
	state.DisplayName = types.StringValue(asString(response["DisplayName"]))
	state.Type = types.StringValue(asString(response["Type"]))
	state.LocationType = types.StringValue(asString(response["LocationType"]))
	state.Value = types.StringValue(valueToJSON(response["Value"]))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// asString coerces an interface{} returned by the ROA client into a string,
// tolerating the numeric/bool values the gateway may emit.
func asString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch s := v.(type) {
	case string:
		return s
	case fmt.Stringer:
		return s.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

// valueToJSON marshals the dynamic Value payload to a stable JSON string; an
// empty payload yields "".
func valueToJSON(v interface{}) string {
	if v == nil {
		return ""
	}
	raw, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(raw)
}
