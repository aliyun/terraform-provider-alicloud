package cms

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &transformersDataSource{}
	_ datasource.DataSourceWithConfigure = &transformersDataSource{}
)

// NewTransformersDataSource is the Framework factory registered for
// alicloud_cms_transformers.
func NewTransformersDataSource() datasource.DataSource {
	return &transformersDataSource{}
}

// transformersDataSource lists CMS Transformer resources filtered by name, id,
// enabled state and workspace.
type transformersDataSource struct {
	fwadapt.DataSourceBase
}

type transformersDataSourceModel struct {
	Id              types.String `tfsdk:"id"`
	Ids             types.List   `tfsdk:"ids"`
	TransformerName types.String `tfsdk:"transformer_name"`
	TransformerId   types.String `tfsdk:"transformer_id"`
	Enable          types.Bool   `tfsdk:"enable"`
	Workspace       types.String `tfsdk:"workspace"`
	Transformers    types.List   `tfsdk:"transformers"`
	OutputFile      types.String `tfsdk:"output_file"`
}

// transformerItemModel mirrors the transformers NestedBlockObject. Only
// the scalar fields the List API returns are flattened here, matching the
// SDKv2 data source output.
type transformerItemModel struct {
	Id              types.String `tfsdk:"id"`
	TransformerId   types.String `tfsdk:"transformer_id"`
	TransformerName types.String `tfsdk:"transformer_name"`
	Workspace       types.String `tfsdk:"workspace"`
	Description     types.String `tfsdk:"description"`
	QuitAfterMatch  types.Bool   `tfsdk:"quit_after_match"`
	SortId          types.Int64  `tfsdk:"sort_id"`
	Enable          types.Bool   `tfsdk:"enable"`
	CreateTime      types.String `tfsdk:"create_time"`
	UpdateTime      types.String `tfsdk:"update_time"`
	UserId          types.String `tfsdk:"user_id"`
	RegionId        types.String `tfsdk:"region_id"`
}

func transformerItemObjectType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":               types.StringType,
		"transformer_id":   types.StringType,
		"transformer_name": types.StringType,
		"workspace":        types.StringType,
		"description":      types.StringType,
		"quit_after_match": types.BoolType,
		"sort_id":          types.Int64Type,
		"enable":           types.BoolType,
		"create_time":      types.StringType,
		"update_time":      types.StringType,
		"user_id":          types.StringType,
		"region_id":        types.StringType,
	}}
}

func (d *transformersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = cmsDSName
}

func (d *transformersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides the CMS Transformers available within the specified workspace.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The data source ID.",
			},
			"ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "A list of transformer IDs, formatted `<transformer_id>:<workspace>`.",
			},
			"transformer_name": schema.StringAttribute{
				Optional:    true,
				Description: "Filter transformers by name.",
			},
			"transformer_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter transformers by id.",
			},
			"enable": schema.BoolAttribute{
				Optional:    true,
				Description: "Filter transformers by enabled state.",
			},
			"workspace": schema.StringAttribute{
				Required:    true,
				Description: "The workspace ID for resource isolation.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "Save the result to a file.",
			},
		},
		Blocks: map[string]schema.Block{
			"transformers": schema.ListNestedBlock{
				Description: "A list of CMS Transformers.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":               schema.StringAttribute{Computed: true},
						"transformer_id":   schema.StringAttribute{Computed: true},
						"transformer_name": schema.StringAttribute{Computed: true},
						"workspace":        schema.StringAttribute{Computed: true},
						"description":      schema.StringAttribute{Computed: true},
						"quit_after_match": schema.BoolAttribute{Computed: true},
						"sort_id":          schema.Int64Attribute{Computed: true},
						"enable":           schema.BoolAttribute{Computed: true},
						"create_time":      schema.StringAttribute{Computed: true},
						"update_time":      schema.StringAttribute{Computed: true},
						"user_id":          schema.StringAttribute{Computed: true},
						"region_id":        schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *transformersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config transformersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := d.Client()

	// ids filter: when the caller supplied ids, only those transformers are
	// kept in the output.
	idSet := map[string]bool{}
	if !config.Ids.IsNull() && !config.Ids.IsUnknown() {
		var ids []string
		resp.Diagnostics.Append(config.Ids.ElementsAs(ctx, &ids, false)...)
		for _, id := range ids {
			idSet[id] = true
		}
	}

	query := map[string]*string{
		"regionId":  stringPointer(client.RegionId),
		"workspace": stringPointer(config.Workspace.ValueString()),
	}
	if v := config.TransformerName.ValueString(); v != "" {
		query["transformerName"] = stringPointer(v)
	}
	if v := config.TransformerId.ValueString(); v != "" {
		query["transformerId"] = stringPointer(v)
	}
	if !config.Enable.IsNull() && !config.Enable.IsUnknown() {
		query["enable"] = stringPointer(strconv.FormatBool(config.Enable.ValueBool()))
	}
	query["MaxResults"] = stringPointer(strconv.Itoa(pageSizeLarge))

	var items []map[string]interface{}
	for {
		response, err := client.RoaGet(cmsProduct, cmsVersion, "/transformers", query, nil, nil)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: calling ListTransformers", cmsDSName),
				errs.WrapErrorf(err, "Reading %s: calling ListTransformers", cmsDSName).Error(),
			)
			return
		}

		list := toInterfaceSlice(response["transformerList"])
		for _, raw := range list {
			item := itemMap(raw)
			if item == nil {
				continue
			}
			if len(idSet) > 0 {
				if !idSet[fmt.Sprintf("%s:%s", strVal(item, "transformerId"), strVal(item, "workspace"))] {
					continue
				}
			}
			items = append(items, item)
		}

		nextToken := strVal(response, "NextToken")
		if nextToken == "" {
			break
		}
		query["NextToken"] = stringPointer(nextToken)
	}

	state := config
	state.Transformers = transformerItemsListFromAPI(ctx, items, &resp.Diagnostics)
	state.Ids = transformerIdsListFromAPI(items)
	state.Id = types.StringValue(config.Workspace.ValueString())

	if v := config.OutputFile.ValueString(); v != "" {
		if err := writeTransformersToFile(v, items); err != nil {
			resp.Diagnostics.AddWarning(
				fmt.Sprintf("Reading %s: writing output_file", cmsDSName),
				err.Error(),
			)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func transformerItemsListFromAPI(ctx context.Context, items []map[string]interface{}, diags *diag.Diagnostics) types.List {
	if len(items) == 0 {
		return types.ListNull(transformerItemObjectType())
	}
	values := make([]attr.Value, 0, len(items))
	for _, item := range items {
		im := transformerItemModel{
			Id:              strType(fmt.Sprintf("%s:%s", strVal(item, "transformerId"), strVal(item, "workspace"))),
			TransformerId:   strType(strVal(item, "transformerId")),
			TransformerName: strType(strVal(item, "transformerName")),
			Workspace:       strType(strVal(item, "workspace")),
			Description:     strType(strVal(item, "description")),
			QuitAfterMatch:  boolType(item["quitAfterMatch"]),
			SortId:          intType(item["sortId"]),
			Enable:          boolType(item["enable"]),
			CreateTime:      strType(strVal(item, "createTime")),
			UpdateTime:      strType(strVal(item, "updateTime")),
			UserId:          strType(strVal(item, "userId")),
			RegionId:        strType(strVal(item, "regionId")),
		}
		obj, d := types.ObjectValueFrom(ctx, transformerItemObjectType().AttrTypes, &im)
		diags.Append(d...)
		values = append(values, obj)
	}
	list, d := types.ListValue(transformerItemObjectType(), values)
	diags.Append(d...)
	return list
}

func transformerIdsListFromAPI(items []map[string]interface{}) types.List {
	if len(items) == 0 {
		return types.ListNull(types.StringType)
	}
	values := make([]attr.Value, 0, len(items))
	for _, item := range items {
		values = append(values, types.StringValue(fmt.Sprintf("%s:%s", strVal(item, "transformerId"), strVal(item, "workspace"))))
	}
	list, _ := types.ListValue(types.StringType, values)
	return list
}

// writeTransformersToFile persists the raw API items as indented JSON. It
// mirrors the SDKv2 writeToFile but lives in-package to avoid importing the
// provider root.
func writeTransformersToFile(path string, items []map[string]interface{}) error {
	out := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		out = append(out, map[string]interface{}{
			"id":               fmt.Sprintf("%s:%s", strVal(item, "transformerId"), strVal(item, "workspace")),
			"transformer_id":   strVal(item, "transformerId"),
			"transformer_name": strVal(item, "transformerName"),
			"workspace":        strVal(item, "workspace"),
			"description":      strVal(item, "description"),
			"quit_after_match": item["quitAfterMatch"],
			"sort_id":          item["sortId"],
			"enable":           item["enable"],
			"create_time":      strVal(item, "createTime"),
			"update_time":      strVal(item, "updateTime"),
			"user_id":          strVal(item, "userId"),
			"region_id":        strVal(item, "regionId"),
		})
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
