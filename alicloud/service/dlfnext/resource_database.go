package dlfnext

import (
	"context"
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const databaseTypeName = "alicloud_dlf_database"

var (
	_ resource.Resource                = &databaseResource{}
	_ resource.ResourceWithImportState = &databaseResource{}
)

// NewDatabaseResource constructs the DlfNext Database Framework resource.
func NewDatabaseResource() resource.Resource {
	return &databaseResource{}
}

type databaseResource struct {
	fwadapt.ResourceBase
}

type databaseResourceModel struct {
	Id         types.String `tfsdk:"id"`
	CatalogId  types.String `tfsdk:"catalog_id"`
	Name       types.String `tfsdk:"name"`
	Options    types.Map    `tfsdk:"options"`
	DatabaseId types.String `tfsdk:"database_id"`
	Location   types.String `tfsdk:"location"`
	Owner      types.String `tfsdk:"owner"`
	TableCount types.Int64  `tfsdk:"table_count"`
	Status     types.String `tfsdk:"status"`
	CreatedAt  types.Int64  `tfsdk:"created_at"`
	CreatedBy  types.String `tfsdk:"created_by"`
	UpdatedAt  types.Int64  `tfsdk:"updated_at"`
	UpdatedBy  types.String `tfsdk:"updated_by"`
}

func (r *databaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = databaseTypeName
}

func (r *databaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a DlfNext Database resource. A database lives under a catalog and contains tables.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource ID in Terraform, in the form `<catalog_id>:<name>`.",
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
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the database.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// options is Optional and intentionally NOT Computed. With Computed,
			// a null config preserves the prior-state value, so "omit to clear"
			// produces no diff and AlterDatabase removals never run. Purely
			// Optional makes config-null plan to null, diffing against a
			// populated prior state to drive the removal. Read still populates
			// state from GetDatabase; DlfNext stores only user-set keys, so a
			// no-options database reads back options: {} and stays drift-free.
			"options": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "The extended options of the database, a map of key-value strings. Omitting the block clears user-set keys; updatable via AlterDatabase.",
			},
			"database_id": schema.StringAttribute{
				Computed:    true,
				Description: "The server-generated database UUID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"location": schema.StringAttribute{
				Computed:    true,
				Description: "The storage location of the database.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"owner": schema.StringAttribute{
				Computed:    true,
				Description: "The owner principal of the database.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"table_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of tables in the database.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The database status.",
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
				Description: "The creator principal of the database.",
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
				Description: "The last modifier principal of the database.",
			},
		},
	}
}

func (r *databaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan databaseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	catalogId := plan.CatalogId.ValueString()
	body := map[string]interface{}{
		"name": plan.Name.ValueString(),
	}
	if opts := optionsToAPI(terraformMapToStringMap(plan.Options)); opts != nil {
		body["options"] = opts
	}

	path := databasesPath(catalogId)
	if _, err := r.Meta().RoaPost(dlfNextProductCode, dlfNextApiVersion, path, nil, nil, body, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateDatabase", databaseTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateDatabase", databaseTypeName).Error(),
		)
		return
	}

	object, err := getDatabase(r.Meta(), catalogId, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling GetDatabase", databaseTypeName),
			errs.WrapErrorf(err, "Creating %s: calling GetDatabase", databaseTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling GetDatabase", databaseTypeName),
			fmt.Sprintf("The database %q is not visible via GetDatabase right after creation", plan.Name.ValueString()),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, databaseState(catalogId, object))...)
}

func (r *databaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state databaseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	catalogId, name, err := parseDatabaseId(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: parsing id", databaseTypeName),
			err.Error(),
		)
		return
	}

	object, err := getDatabase(r.Meta(), catalogId, name)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetDatabase", databaseTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetDatabase", databaseTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, databaseState(catalogId, object))...)
}

func (r *databaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state databaseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// All attributes except `options` are ForceNew; only options can drift.
	updates, removals := diffOptions(
		terraformMapToStringMap(state.Options),
		terraformMapToStringMap(plan.Options),
	)
	if len(updates) != 0 || len(removals) != 0 {
		body := map[string]interface{}{}
		if len(updates) != 0 {
			body["updates"] = updates
		}
		if len(removals) != 0 {
			body["removals"] = removals
		}
		catalogId, name, err := parseDatabaseId(state.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: parsing id", databaseTypeName),
				err.Error(),
			)
			return
		}
		path := databasePath(catalogId, name)
		if _, err := r.Meta().RoaPost(dlfNextProductCode, dlfNextApiVersion, path, nil, nil, body, true); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling AlterDatabase", databaseTypeName),
				errs.WrapErrorf(err, "Updating %s: calling AlterDatabase", databaseTypeName).Error(),
			)
			return
		}
	}

	catalogId, name, err := parseDatabaseId(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: parsing id", databaseTypeName),
			err.Error(),
		)
		return
	}
	object, err := getDatabase(r.Meta(), catalogId, name)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling GetDatabase", databaseTypeName),
			errs.WrapErrorf(err, "Updating %s: calling GetDatabase", databaseTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, databaseState(catalogId, object))...)
}

func (r *databaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state databaseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	catalogId, name, err := parseDatabaseId(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: parsing id", databaseTypeName),
			err.Error(),
		)
		return
	}
	path := databasePath(catalogId, name)
	_, err = r.Meta().RoaDelete(dlfNextProductCode, dlfNextApiVersion, path, nil, nil, nil, true)
	if err != nil && !isNotFound(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DropDatabase", databaseTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DropDatabase", databaseTypeName).Error(),
		)
		return
	}
}

func (r *databaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// parseDatabaseId splits a "<catalog_id>:<name>" resource ID.
func parseDatabaseId(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid database id %q, expected <catalog_id>:<name>", id)
	}
	return parts[0], parts[1], nil
}

// databaseState maps a GetDatabase response onto the resource model. The
// Terraform `id` composes the parent catalog id and the database name so Read
// and Import can recover both path parameters.
func databaseState(catalogId string, o map[string]interface{}) databaseResourceModel {
	name := stringFromResponse(o, "name")
	return databaseResourceModel{
		Id:         types.StringValue(catalogId + ":" + name),
		CatalogId:  types.StringValue(catalogId),
		Name:       stringAttr(name),
		Options:    stringMapToTerraform(optionsFromAPI(o["options"])),
		DatabaseId: stringAttr(stringFromResponse(o, "id")),
		Location:   stringAttr(stringFromResponse(o, "location")),
		Owner:      stringAttr(stringFromResponse(o, "owner")),
		TableCount: types.Int64Value(int64FromResponse(o, "tableCount")),
		Status:     stringAttr(stringFromResponse(o, "status")),
		CreatedAt:  types.Int64Value(int64FromResponse(o, "createdAt")),
		CreatedBy:  stringAttr(stringFromResponse(o, "createdBy")),
		UpdatedAt:  types.Int64Value(int64FromResponse(o, "updatedAt")),
		UpdatedBy:  stringAttr(stringFromResponse(o, "updatedBy")),
	}
}
