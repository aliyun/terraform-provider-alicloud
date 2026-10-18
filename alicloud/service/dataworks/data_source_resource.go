package dataworks

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const resourceDataSourceTypeName = "alicloud_data_works_resource"

var (
	_ datasource.DataSource              = &resourceDataSource{}
	_ datasource.DataSourceWithConfigure = &resourceDataSource{}
)

// NewResourceDataSource is the Framework factory for the listing data source.
// Like the SDKv2 original it shares the `alicloud_data_works_resource` type
// name with the resource — resource and data source namespaces are separate.
func NewResourceDataSource() datasource.DataSource {
	return &resourceDataSource{}
}

type resourceDataSource struct {
	fwadapt.DataSourceBase
}

type resourceDataSourceModel struct {
	Id         types.String `tfsdk:"id"`
	ProjectId  types.String `tfsdk:"project_id"`
	Owner      types.String `tfsdk:"owner"`
	Type       types.String `tfsdk:"type"`
	Ids        types.List   `tfsdk:"ids"`
	OutputFile types.String `tfsdk:"output_file"`
	Resources  types.List   `tfsdk:"resources"`
}

// resourceEntryModel is one element of the `resources` list. The nested
// data_source/script objects of the SDKv2 schema are flattened to top-level
// attributes, matching the resource schema.
type resourceEntryModel struct {
	Id               types.String `tfsdk:"id"`
	ResourceId       types.String `tfsdk:"resource_id"`
	ResourceName     types.String `tfsdk:"resource_name"`
	ProjectId        types.String `tfsdk:"project_id"`
	Owner            types.String `tfsdk:"owner"`
	Type             types.String `tfsdk:"type"`
	CreateTime       types.Int64  `tfsdk:"create_time"`
	ModifyTime       types.Int64  `tfsdk:"modify_time"`
	SourceType       types.String `tfsdk:"source_type"`
	SourcePath       types.String `tfsdk:"source_path"`
	TargetType       types.String `tfsdk:"target_type"`
	TargetPath       types.String `tfsdk:"target_path"`
	DataSourceName   types.String `tfsdk:"data_source_name"`
	DataSourceType   types.String `tfsdk:"data_source_type"`
	ScriptPath       types.String `tfsdk:"script_path"`
	ScriptRuntimeCmd types.String `tfsdk:"script_runtime_command"`
	ScriptId         types.String `tfsdk:"script_id"`
}

func (d *resourceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = resourceDataSourceTypeName
}

func (d *resourceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a list of Data Works Resource items owned by the user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The data source ID (hash of the matched resource IDs).",
			},
			"project_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the DataWorks workspace to query resources for.",
			},
			"owner": schema.StringAttribute{
				Optional:    true,
				Description: "The owner of the resource, used to filter results.",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "The type of the resource. Valid values: `jar`, `python`, `file`, `archive`.",
			},
			"ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "A list of resource IDs to filter results.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "File path where results will be saved after running `terraform plan`.",
			},
			"resources": schema.ListNestedAttribute{
				Computed:    true,
				Description: "A list of Data Works Resource items.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                     schema.StringAttribute{Computed: true, Description: "The resource ID, formats as `<project_id>:<resource_id>`."},
						"resource_id":            schema.StringAttribute{Computed: true, Description: "The unique identifier of the resource."},
						"resource_name":          schema.StringAttribute{Computed: true, Description: "The name of the resource."},
						"project_id":             schema.StringAttribute{Computed: true, Description: "The ID of the workspace."},
						"owner":                  schema.StringAttribute{Computed: true, Description: "The responsible person of the resource file."},
						"type":                   schema.StringAttribute{Computed: true, Description: "The resource type."},
						"create_time":            schema.Int64Attribute{Computed: true, Description: "The creation time of the resource."},
						"modify_time":            schema.Int64Attribute{Computed: true, Description: "The last modification time of the resource."},
						"source_type":            schema.StringAttribute{Computed: true, Description: "The file resource source storage type."},
						"source_path":            schema.StringAttribute{Computed: true, Description: "The file source path."},
						"target_type":            schema.StringAttribute{Computed: true, Description: "The file destination storage type."},
						"target_path":            schema.StringAttribute{Computed: true, Description: "The file destination storage path."},
						"data_source_name":       schema.StringAttribute{Computed: true, Description: "The data source name."},
						"data_source_type":       schema.StringAttribute{Computed: true, Description: "The data source type."},
						"script_path":            schema.StringAttribute{Computed: true, Description: "The script path."},
						"script_runtime_command": schema.StringAttribute{Computed: true, Description: "The script runtime command."},
						"script_id":              schema.StringAttribute{Computed: true, Description: "The script ID."},
					},
				},
			},
		},
	}
}

func (d *resourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config resourceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var idFilter []string
	resp.Diagnostics.Append(config.Ids.ElementsAs(ctx, &idFilter, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entries, err := resourceDataSourceRead(d.Client(), config.ProjectId.ValueString(), config.Owner.ValueString(), config.Type.ValueString(), idFilter)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListResources", resourceDataSourceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling ListResources", resourceDataSourceTypeName).Error(),
		)
		return
	}

	list, diags := types.ListValueFrom(ctx, resourceEntryObjectType(), entries)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids := make([]string, 0, len(entries))
	for _, e := range entries {
		ids = append(ids, e.Id.ValueString())
	}
	idsList, diags := types.ListValueFrom(ctx, types.StringType, ids)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := resourceDataSourceModel{
		Id:         types.StringValue(dataSourceID(ids)),
		ProjectId:  config.ProjectId,
		Owner:      config.Owner,
		Type:       config.Type,
		Ids:        idsList,
		OutputFile: config.OutputFile,
		Resources:  list,
	}

	if out := config.OutputFile.ValueString(); out != "" {
		writeDataWorksOutputFile(out, entries)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// resourceDataSourceRead lists resources and maps each to an entry model. The
// ids filter narrows to the requested resource ids; owner/type pass through to
// the API.
func resourceDataSourceRead(client *connectivity.AliyunClient, projectID, owner, typ string, ids []string) ([]resourceEntryModel, error) {
	idSet := make(map[string]bool, len(ids))
	for _, id := range ids {
		if id != "" {
			idSet[id] = true
		}
	}
	objects, err := listResources(client, projectID, owner, typ, idSet)
	if err != nil {
		return nil, err
	}
	entries := make([]resourceEntryModel, 0, len(objects))
	for _, object := range objects {
		entries = append(entries, resourceEntryState(object, projectID))
	}
	return entries, nil
}

// resourceEntryState maps one ListResources item to an entry model. projectID is
// the requested workspace id (the API item carries it too, but the caller knows
// it without an extra assertion).
func resourceEntryState(object map[string]interface{}, projectID string) resourceEntryModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	resourceID := asString(object["Id"])
	dataSource := nestedMap(object, "DataSource")
	script := nestedMap(object, "Script")
	return resourceEntryModel{
		Id:               types.StringValue(fmt.Sprintf("%s:%s", projectID, resourceID)),
		ResourceId:       types.StringValue(resourceID),
		ResourceName:     str("Name"),
		ProjectId:        stringOrZero(object, "ProjectId"),
		Owner:            str("Owner"),
		Type:             str("Type"),
		CreateTime:       types.Int64Value(asInt64(object["CreateTime"])),
		ModifyTime:       types.Int64Value(asInt64(object["ModifyTime"])),
		SourceType:       str("SourceType"),
		SourcePath:       str("SourcePath"),
		TargetType:       str("TargetType"),
		TargetPath:       str("TargetPath"),
		DataSourceName:   stringOrZero(dataSource, "Name"),
		DataSourceType:   stringOrZero(dataSource, "Type"),
		ScriptPath:       stringOrZero(script, "Path"),
		ScriptRuntimeCmd: stringOrZero(runtimeCommand(script)),
		ScriptId:         stringOrZero(script, "Id"),
	}
}

// resourceEntryObjectType returns the object type for the `resources` list
// element, used by types.ListValueFrom to reflect []resourceEntryModel into a
// List of objects.
func resourceEntryObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":                     types.StringType,
			"resource_id":            types.StringType,
			"resource_name":          types.StringType,
			"project_id":             types.StringType,
			"owner":                  types.StringType,
			"type":                   types.StringType,
			"create_time":            types.Int64Type,
			"modify_time":            types.Int64Type,
			"source_type":            types.StringType,
			"source_path":            types.StringType,
			"target_type":            types.StringType,
			"target_path":            types.StringType,
			"data_source_name":       types.StringType,
			"data_source_type":       types.StringType,
			"script_path":            types.StringType,
			"script_runtime_command": types.StringType,
			"script_id":              types.StringType,
		},
	}
}

// dataSourceID returns a stable hex id for the data source from the matched
// resource ids, matching the SDKv2 dataResourceIdHash intent (a hash of ids).
func dataSourceID(ids []string) string {
	if len(ids) == 0 {
		return "-"
	}
	sum := sha1.Sum([]byte(strings.Join(filterNonEmpty(ids), ",")))
	return hex.EncodeToString(sum[:])
}

func filterNonEmpty(parts []string) []string {
	out := parts[:0:0]
	for _, p := range parts {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// writeDataWorksOutputFile dumps the entries as JSON; failures are non-fatal.
func writeDataWorksOutputFile(path string, entries []resourceEntryModel) {
	if path == "" {
		return
	}
	b, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(path, b, 0o644)
}
