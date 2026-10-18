package dataworks

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// nodeIoEntryObjectType mirrors the four string attributes of nodeIoEntryModel so a
// types.List of entries can be built via ListValueFrom without re-deriving the schema.
var nodeIoEntryObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":         types.StringType,
		"table_name": types.StringType,
		"data":       types.StringType,
		"node_id":    types.StringType,
	},
}

const nodeIoDataSourceTypeName = "alicloud_data_works_node_io"

// ListNodeIO lives under the dataworks-public product, version 2020-05-18.
const (
	nodeIoProduct = "dataworks-public"
	nodeIoVersion = "2020-05-18"
	nodeIoAction  = "ListNodeIO"
)

var (
	_ datasource.DataSource              = &nodeIoDataSource{}
	_ datasource.DataSourceWithConfigure = &nodeIoDataSource{}
)

func NewNodeIoDataSource() datasource.DataSource {
	return &nodeIoDataSource{}
}

type nodeIoDataSource struct {
	fwadapt.DataSourceBase
}

// nodeIoDataSourceModel is the shape of the data source's state. The Optional+Computed
// ids attribute doubles as a filter (input) and the list of returned table names (output).
type nodeIoDataSourceModel struct {
	Id         types.String `tfsdk:"id"`
	NodeId     types.String `tfsdk:"node_id"`
	ProjectEnv types.String `tfsdk:"project_env"`
	IoType     types.String `tfsdk:"io_type"`
	Ids        types.List   `tfsdk:"ids"`
	OutputFile types.String `tfsdk:"output_file"`
	NodeIos    types.List   `tfsdk:"node_ios"`
}

// nodeIoEntryModel is one element of the node_ios list.
type nodeIoEntryModel struct {
	Id        types.String `tfsdk:"id"`
	TableName types.String `tfsdk:"table_name"`
	Data      types.String `tfsdk:"data"`
	NodeId    types.String `tfsdk:"node_id"`
}

func (d *nodeIoDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = nodeIoDataSourceTypeName
}

func (d *nodeIoDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides the Data Works Node Input/Output list of the current Alibaba Cloud user according to the specified node, project environment and IO type.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source. It is a hash of the returned node IO entry table names.",
			},
			"node_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the DataWorks node. ListNodeIO queries the upstream or downstream nodes of this node.",
			},
			"project_env": schema.StringAttribute{
				Required:    true,
				Description: "The environment of the DataWorks project. Valid values: `PROD`, `DEV`.",
			},
			"io_type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the node input/output. It queries the upstream (`input`) or downstream (`output`) nodes only one level.",
			},
			"ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "A list of node IO entry IDs (the related table names). If set, only the entries whose table name matches one of the provided IDs are returned.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"node_ios": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the node IO entry. It is the same value as `table_name`.",
						},
						"table_name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the upstream or downstream table associated with the node.",
						},
						"data": schema.StringAttribute{
							Computed:    true,
							Description: "The data content of the node input/output.",
						},
						"node_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the DataWorks node.",
						},
					},
				},
				Description: "A list of Data Works Node Input/Output entries.",
			},
		},
	}
}

func (d *nodeIoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config nodeIoDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := d.Client()

	request := map[string]interface{}{
		"NodeId":     config.NodeId.ValueString(),
		"ProjectEnv": config.ProjectEnv.ValueString(),
		"IoType":     config.IoType.ValueString(),
	}

	// ids is an Optional filter: a set of table names the caller cares about.
	idsFilter := make(map[string]struct{})
	for _, el := range config.Ids.Elements() {
		v := el.(types.String).ValueString()
		if v == "" {
			continue
		}
		idsFilter[v] = struct{}{}
	}

	response, err := nodeIoWithRetry(client.RpcPost, request)
	logs.AddDebug(nodeIoAction, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", nodeIoDataSourceTypeName, nodeIoAction),
			errs.WrapErrorf(err, "Reading %s: calling %s", nodeIoDataSourceTypeName, nodeIoAction).Error(),
		)
		return
	}

	rawData, ok := response["Data"]
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", nodeIoDataSourceTypeName, nodeIoAction),
			fmt.Sprintf("The response carried no Data: %v", response),
		)
		return
	}
	items, _ := rawData.([]interface{})

	entries := make([]nodeIoEntryModel, 0)
	ids := make([]string, 0)
	for _, v := range items {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		tableName := fmt.Sprint(item["TableName"])
		if len(idsFilter) > 0 {
			if _, keep := idsFilter[tableName]; !keep {
				continue
			}
		}
		nodeId := fmt.Sprint(item["NodeId"])
		entries = append(entries, nodeIoEntryModel{
			Id:        types.StringValue(tableName),
			TableName: types.StringValue(tableName),
			Data:      types.StringValue(fmt.Sprint(item["Data"])),
			NodeId:    types.StringValue(nodeId),
		})
		ids = append(ids, tableName)
	}

	config.Id = types.StringValue(nodeIoIdHash(ids))
	idsList, diags := types.ListValueFrom(ctx, types.StringType, toInterface(ids))
	resp.Diagnostics.Append(diags...)
	config.Ids = idsList

	nodeIosList, diags := types.ListValueFrom(ctx, nodeIoEntryObjectType, entries)
	resp.Diagnostics.Append(diags...)
	config.NodeIos = nodeIosList
	if resp.Diagnostics.HasError() {
		return
	}

	if path := config.OutputFile.ValueString(); path != "" {
		if err := nodeIoWriteOutputFile(path, entries); err != nil {
			resp.Diagnostics.AddWarning(
				fmt.Sprintf("Reading %s: writing output_file", nodeIoDataSourceTypeName),
				err.Error(),
			)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// nodeIoWithRetry calls ListNodeIO with bounded retry on transient errors, mirroring
// the SDK v2 resource.Retry behaviour the original data source used.
func nodeIoWithRetry(post func(product, version, action string, headers map[string]interface{}, request map[string]interface{}, useAk bool) (map[string]interface{}, error),
	request map[string]interface{}) (map[string]interface{}, error) {
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	deadline := time.Now().Add(5 * time.Minute)
	var response map[string]interface{}
	var err error
	for {
		response, err = post(nodeIoProduct, nodeIoVersion, nodeIoAction, nil, request, true)
		if err == nil {
			return response, nil
		}
		if !helper.NeedRetry(err) {
			return response, err
		}
		if time.Now().After(deadline) {
			return response, err
		}
		wait()
	}
}

func nodeIoIdHash(ids []string) string {
	var buf strings.Builder
	for _, id := range ids {
		buf.WriteString(id)
		buf.WriteByte('-')
	}
	return fmt.Sprintf("%d", helper.Hashcode(buf.String()))
}

func nodeIoWriteOutputFile(path string, entries []nodeIoEntryModel) error {
	out := make([]map[string]string, 0, len(entries))
	for _, e := range entries {
		out = append(out, map[string]string{
			"id":         e.Id.ValueString(),
			"table_name": e.TableName.ValueString(),
			"data":       e.Data.ValueString(),
			"node_id":    e.NodeId.ValueString(),
		})
	}
	if strings.HasPrefix(path, "~") {
		home, _ := os.UserHomeDir()
		if home != "" {
			path = strings.Replace(path, "~", home, 1)
		}
	}
	bs, err := json.MarshalIndent(out, "", "\t")
	if err != nil {
		return fmt.Errorf("marshal node_io entries: %w", err)
	}
	return os.WriteFile(path, bs, 0o644)
}

func toInterface(in []string) []interface{} {
	out := make([]interface{}, 0, len(in))
	for _, s := range in {
		out = append(out, s)
	}
	return out
}
