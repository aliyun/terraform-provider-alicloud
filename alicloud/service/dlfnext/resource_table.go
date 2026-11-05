package dlfnext

import (
	"context"
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const tableTypeName = "alicloud_dlf_table"

var (
	_ resource.Resource                = &tableResource{}
	_ resource.ResourceWithImportState = &tableResource{}
)

// NewTableResource constructs the DlfNext Table Framework resource.
func NewTableResource() resource.Resource {
	return &tableResource{}
}

type tableResource struct {
	fwadapt.ResourceBase
}

// tableFieldModel mirrors one element of the `fields` nested list. The column
// `type` is a SQL-like expression (STRING, ARRAY<INT>, MAP<STRING,BIGINT>,
// ROW<id INT, name STRING>) that parseTypeString turns into the recursive
// FullDataType object the API expects.
type tableFieldModel struct {
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`
}

type tableResourceModel struct {
	Id            types.String `tfsdk:"id"`
	CatalogId     types.String `tfsdk:"catalog_id"`
	DatabaseName  types.String `tfsdk:"database_name"`
	Name          types.String `tfsdk:"name"`
	Comment       types.String `tfsdk:"comment"`
	Fields        types.List   `tfsdk:"fields"`
	Options       types.Map    `tfsdk:"options"`
	PartitionKeys types.List   `tfsdk:"partition_keys"`
	PrimaryKeys   types.List   `tfsdk:"primary_keys"`
	TableId       types.String `tfsdk:"table_id"`
	Owner         types.String `tfsdk:"owner"`
	Status        types.String `tfsdk:"status"`
	CreatedAt     types.Int64  `tfsdk:"created_at"`
	CreatedBy     types.String `tfsdk:"created_by"`
	UpdatedAt     types.Int64  `tfsdk:"updated_at"`
	UpdatedBy     types.String `tfsdk:"updated_by"`
}

func (r *tableResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = tableTypeName
}

func (r *tableResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a DlfNext Table resource. A table lives under a database and defines a schema of typed columns.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource ID in Terraform, in the form `<catalog_id>:<database_name>:<name>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"catalog_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the parent catalog (the `catalog_id` of an alicloud_dlf_catalog resource).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"database_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the parent database (the `name` of an alicloud_dlf_database resource).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the table.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Description: "A free-form comment for the table.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"options": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "Extended table options, a map of key-value strings.",
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"partition_keys": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "The partition key column names.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"primary_keys": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "The primary key column names.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"table_id": schema.StringAttribute{
				Computed:    true,
				Description: "The server-generated table ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"owner": schema.StringAttribute{
				Computed:    true,
				Description: "The owner principal of the table.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The table status.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.Int64Attribute{
				Computed:    true,
				Description: "The creation time, epoch milliseconds.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "The creator principal of the table.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.Int64Attribute{
				Computed:    true,
				Description: "The last modification time, epoch milliseconds.",
			},
			"updated_by": schema.StringAttribute{
				Computed:    true,
				Description: "The last modifier principal of the table.",
			},
		},
		Blocks: map[string]schema.Block{
			// `fields` is a ListNestedBlock (rather than a ListNestedAttribute)
			// because the alicloud provider serves Framework resources through the
			// v5 protocol bridge (tf6to5server), which cannot convert Framework
			// nested attributes into v5. A Block carries nested attributes in a v5
			// Block (which supports attributes), so it converts cleanly; the dts
			// timeouts SingleNestedBlock already relies on this path. description is
			// Optional here because block attributes are individually marked, unlike
			// an object type whose attributes are all implicitly required.
			"fields": schema.ListNestedBlock{
				Description: "The table columns. Each block has a `name`, a SQL-like `type`, and an optional `description`. Changing any field forces replacement.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "The column name.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"type": schema.StringAttribute{
							Required:    true,
							Description: "The column type, e.g. `STRING`, `BIGINT`, `ARRAY<INT>`, `MAP<STRING,BIGINT>`, `ROW<id INT, name STRING>`.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Description: "An optional description of the column.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
				},
			},
		},
	}
}

func (r *tableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	catalogId := plan.CatalogId.ValueString()
	database := plan.DatabaseName.ValueString()
	name := plan.Name.ValueString()

	fields, err := terraformFieldsToAPI(ctx, plan.Fields)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: parsing column types", tableTypeName),
			err.Error(),
		)
		return
	}
	schemaBody := map[string]interface{}{
		"fields": fields,
	}
	if v := plan.Comment.ValueString(); v != "" {
		schemaBody["comment"] = v
	}
	if opts := optionsToAPI(terraformMapToStringMap(plan.Options)); opts != nil {
		schemaBody["options"] = opts
	}
	if pk := terraformListToStringSlice(plan.PartitionKeys); len(pk) > 0 {
		schemaBody["partitionKeys"] = pk
	}
	if pk := terraformListToStringSlice(plan.PrimaryKeys); len(pk) > 0 {
		schemaBody["primaryKeys"] = pk
	}

	body := map[string]interface{}{
		"identifier": map[string]interface{}{
			"database": database,
			"object":   name,
		},
		"schema": schemaBody,
	}

	path := tablesPath(catalogId, database)
	if _, err := r.Meta().RoaPost(dlfNextProductCode, dlfNextApiVersion, path, nil, nil, body, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateTable", tableTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateTable", tableTypeName).Error(),
		)
		return
	}

	object, err := getTable(r.Meta(), catalogId, database, name)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling GetTable", tableTypeName),
			errs.WrapErrorf(err, "Creating %s: calling GetTable", tableTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling GetTable", tableTypeName),
			fmt.Sprintf("The table %q is not visible via GetTable right after creation", name),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, tableState(catalogId, database, name, object))...)
}

func (r *tableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	catalogId, database, name, err := parseTableId(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: parsing id", tableTypeName),
			err.Error(),
		)
		return
	}

	object, err := getTable(r.Meta(), catalogId, database, name)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetTable", tableTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetTable", tableTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, tableState(catalogId, database, name, object))...)
}

// Update re-reads the table state. Every user-settable attribute is
// RequiresReplace, so a real change drives a destroy+recreate and this method
// is only reached when a pure refresh is requested.
func (r *tableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state tableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	catalogId, database, name, err := parseTableId(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: parsing id", tableTypeName),
			err.Error(),
		)
		return
	}
	object, err := getTable(r.Meta(), catalogId, database, name)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling GetTable", tableTypeName),
			errs.WrapErrorf(err, "Updating %s: calling GetTable", tableTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, tableState(catalogId, database, name, object))...)
}

func (r *tableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	catalogId, database, name, err := parseTableId(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: parsing id", tableTypeName),
			err.Error(),
		)
		return
	}
	path := tablePath(catalogId, database, name)
	_, err = r.Meta().RoaDelete(dlfNextProductCode, dlfNextApiVersion, path, nil, nil, nil, true)
	if err != nil && !isNotFound(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DropTable", tableTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DropTable", tableTypeName).Error(),
		)
		return
	}
}

func (r *tableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// parseTableId splits a "<catalog_id>:<database>:<name>" resource ID.
func parseTableId(id string) (string, string, string, error) {
	parts := strings.SplitN(id, ":", 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", fmt.Errorf("invalid table id %q, expected <catalog_id>:<database>:<name>", id)
	}
	return parts[0], parts[1], parts[2], nil
}

// fieldsAttrTypes is the object type of one `fields` element.
var fieldsAttrTypes = map[string]attr.Type{
	"name":        types.StringType,
	"type":        types.StringType,
	"description": types.StringType,
}

func fieldsObjectType() types.ObjectType {
	return types.ObjectType{AttrTypes: fieldsAttrTypes}
}

// terraformFieldsToAPI converts the plan's `fields` list into the DataField[]
// the CreateTable body expects, parsing each column type expression into a
// FullDataType object.
func terraformFieldsToAPI(ctx context.Context, l types.List) ([]map[string]interface{}, error) {
	if l.IsNull() || l.IsUnknown() {
		return []map[string]interface{}{}, nil
	}
	var models []tableFieldModel
	if diags := l.ElementsAs(ctx, &models, false); diags.HasError() {
		return nil, fmt.Errorf("decoding fields list: %s", diags.Errors())
	}
	out := make([]map[string]interface{}, 0, len(models))
	for i, f := range models {
		t, err := parseTypeString(f.Type.ValueString())
		if err != nil {
			return nil, fmt.Errorf("column %q: %w", f.Name.ValueString(), err)
		}
		// DataField.id is a required 1-based ordinal ("表字段号"); the API
		// rejects a null id with "provided data field id is null for field id".
		obj := map[string]interface{}{
			"id":   i + 1,
			"name": f.Name.ValueString(),
			"type": t,
		}
		if d := f.Description.ValueString(); d != "" {
			obj["description"] = d
		}
		out = append(out, obj)
	}
	return out, nil
}

// tableFieldsToTerraform converts the DataField[] of a GetTable response into a
// types.List of field objects, rendering each FullDataType back into the
// SQL-like type expression.
func tableFieldsToTerraform(v interface{}) types.List {
	arr, _ := v.([]interface{})
	if len(arr) == 0 {
		return types.ListNull(fieldsObjectType())
	}
	elems := make([]attr.Value, 0, len(arr))
	for _, f := range arr {
		fm := asMap(f)
		obj, _ := types.ObjectValue(fieldsAttrTypes, map[string]attr.Value{
			"name":        stringAttr(stringFromResponse(fm, "name")),
			"type":        stringAttr(formatTypeString(asMap(fm["type"]))),
			"description": stringAttr(stringFromResponse(fm, "description")),
		})
		elems = append(elems, obj)
	}
	list, _ := types.ListValue(fieldsObjectType(), elems)
	return list
}

// tableState maps a GetTable response onto the resource model. The Terraform
// `id` composes the catalog id, database name and table name so Read and
// Import recover all three path parameters.
func tableState(catalogId, database, name string, o map[string]interface{}) tableResourceModel {
	ident := asMap(o["identifier"])
	if ident != nil {
		if v := stringFromResponse(ident, "object"); v != "" {
			name = v
		}
		if v := stringFromResponse(ident, "database"); v != "" {
			database = v
		}
	}
	var comment string
	var fieldsList types.List = types.ListNull(fieldsObjectType())
	var options map[string]string
	var partitionKeys, primaryKeys []string
	if schemaObj := asMap(o["schema"]); schemaObj != nil {
		comment = stringFromResponse(schemaObj, "comment")
		fieldsList = tableFieldsToTerraform(schemaObj["fields"])
		options = optionsFromAPI(schemaObj["options"])
		partitionKeys = stringSlice(schemaObj["partitionKeys"])
		primaryKeys = stringSlice(schemaObj["primaryKeys"])
	}
	return tableResourceModel{
		Id:            types.StringValue(catalogId + ":" + database + ":" + name),
		CatalogId:     types.StringValue(catalogId),
		DatabaseName:  stringAttr(database),
		Name:          stringAttr(name),
		Comment:       stringAttr(comment),
		Fields:        fieldsList,
		Options:       stringMapToTerraform(options),
		PartitionKeys: stringSliceToTerraform(partitionKeys),
		PrimaryKeys:   stringSliceToTerraform(primaryKeys),
		TableId:       stringAttr(stringFromResponse(o, "id")),
		Owner:         stringAttr(stringFromResponse(o, "owner")),
		Status:        stringAttr(stringFromResponse(o, "status")),
		CreatedAt:     types.Int64Value(int64FromResponse(o, "createdAt")),
		CreatedBy:     stringAttr(stringFromResponse(o, "createdBy")),
		UpdatedAt:     types.Int64Value(int64FromResponse(o, "updatedAt")),
		UpdatedBy:     stringAttr(stringFromResponse(o, "updatedBy")),
	}
}
