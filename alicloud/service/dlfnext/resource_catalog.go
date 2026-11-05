package dlfnext

import (
	"context"
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const catalogTypeName = "alicloud_dlf_catalog"

// catalogReadyTimeout bounds how long Create/Delete waits for an asynchronously
// provisioned catalog to reach RUNNING. DlfNext moves a catalog through
// NEW -> INITIALIZING -> RUNNING; DropCatalog and database/table creation
// fail until RUNNING.
const catalogReadyTimeout = 10 * time.Minute

var (
	_ resource.Resource                = &catalogResource{}
	_ resource.ResourceWithImportState = &catalogResource{}
)

// NewCatalogResource constructs the DlfNext Catalog Framework resource.
func NewCatalogResource() resource.Resource {
	return &catalogResource{}
}

type catalogResource struct {
	fwadapt.ResourceBase
}

type catalogResourceModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Type      types.String `tfsdk:"type"`
	IsShared  types.Bool   `tfsdk:"is_shared"`
	Options   types.Map    `tfsdk:"options"`
	ShareId   types.String `tfsdk:"share_id"`
	CatalogId types.String `tfsdk:"catalog_id"`
	Status    types.String `tfsdk:"status"`
	Owner     types.String `tfsdk:"owner"`
	CreatedAt types.Int64  `tfsdk:"created_at"`
	CreatedBy types.String `tfsdk:"created_by"`
	UpdatedAt types.Int64  `tfsdk:"updated_at"`
	UpdatedBy types.String `tfsdk:"updated_by"`
}

func (r *catalogResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = catalogTypeName
}

func (r *catalogResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a DlfNext Catalog resource. A catalog is the top-level container of a Data Lake Formation deployment; databases and tables live under it.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource ID in Terraform, equal to the catalog `name`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the catalog. It is the path identifier used by Get/Alter/Drop and is unique within the account.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the catalog. Valid values: `PAIMON`, `ICEBERG`.",
				Validators: []validator.String{
					stringvalidator.OneOf("PAIMON", "ICEBERG"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"is_shared": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the catalog is shared. Defaults to `false`.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"options": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "The extended options of the catalog, a map of key-value strings. Read-only: DlfNext populates storage defaults such as `dlf.trashed-file-retained-days` and `storage.data.redundancy.type`; user-set keys would conflict with these server-managed defaults, so options is surfaced only as an exported attribute.",
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"share_id": schema.StringAttribute{
				Optional:    true,
				Description: "The share ID. Required when `is_shared` is `true`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"catalog_id": schema.StringAttribute{
				Computed:    true,
				Description: "The server-generated catalog ID (e.g. `clg-paimon-xxxx`). Reference it from alicloud_dlf_database.catalog_id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The catalog status: NEW, INITIALIZING, INITIALIZE_FAILED, RUNNING, TERMINATED, DELETING, DELETE_FAILED, DELETED, STORAGE_RESTRICTED.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"owner": schema.StringAttribute{
				Computed:    true,
				Description: "The owner principal of the catalog.",
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
				Description: "The creator principal of the catalog.",
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
				Description: "The last modifier principal of the catalog.",
			},
		},
	}
}

func (r *catalogResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan catalogResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{
		"name":     plan.Name.ValueString(),
		"type":     plan.Type.ValueString(),
		"isShared": plan.IsShared.ValueBool(),
	}
	// options is read-only (Computed); DlfNext applies storage defaults itself.
	if v := plan.ShareId.ValueString(); v != "" {
		body["shareId"] = v
	}

	if _, err := r.Meta().RoaPost(dlfNextProductCode, dlfNextApiVersion, catalogsPath, nil, nil, body, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateCatalog", catalogTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateCatalog", catalogTypeName).Error(),
		)
		return
	}

	// DlfNext catalogs are provisioned asynchronously: CreateCatalog returns
	// while the catalog moves NEW -> INITIALIZING -> RUNNING. DropCatalog and
	// database/table creation fail until RUNNING, so wait for it here before
	// recording state.
	if err := waitCatalogStatus(ctx, r.Meta(), plan.Name.ValueString(), "RUNNING", catalogReadyTimeout); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for catalog to become ready", catalogTypeName),
			err.Error(),
		)
		return
	}

	object, err := getCatalog(r.Meta(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling GetCatalog", catalogTypeName),
			errs.WrapErrorf(err, "Creating %s: calling GetCatalog", catalogTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling GetCatalog", catalogTypeName),
			fmt.Sprintf("The catalog %q is not visible via GetCatalog right after creation", plan.Name.ValueString()),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, catalogState(object))...)
}

func (r *catalogResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state catalogResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The catalog `id` equals the catalog `name`, which is also the path
	// identifier; resolving from `id` keeps Read consistent right after an
	// import where `name` is not yet populated.
	object, err := getCatalog(r.Meta(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetCatalog", catalogTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetCatalog", catalogTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, catalogState(object))...)
}

// Update re-reads the catalog state. `options` is read-only (Computed) and
// every user-settable attribute (name, type, is_shared, share_id) is
// RequiresReplace, so a real change drives a destroy+recreate and this method
// is only reached on a pure refresh.
func (r *catalogResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state catalogResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := getCatalog(r.Meta(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling GetCatalog", catalogTypeName),
			errs.WrapErrorf(err, "Updating %s: calling GetCatalog", catalogTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, catalogState(object))...)
}

func (r *catalogResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state catalogResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := state.Name.ValueString()
	path := catalogPath(name)
	_, err := r.Meta().RoaDelete(dlfNextProductCode, dlfNextApiVersion, path, nil, nil, nil, true)
	if err != nil && !isNotFound(err) {
		// A catalog still provisioning returns 403 "is being created, can not
		// be dropped". CreateCatalog already waits for RUNNING, but defend
		// against a catalog created out-of-band or re-provisioned: wait for it
		// to become ready, then retry the drop once.
		if isBeingCreated(err) {
			if waitErr := waitCatalogStatus(ctx, r.Meta(), name, "RUNNING", catalogReadyTimeout); waitErr != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Deleting %s: waiting for catalog to become ready", catalogTypeName),
					waitErr.Error(),
				)
				return
			}
			_, err = r.Meta().RoaDelete(dlfNextProductCode, dlfNextApiVersion, path, nil, nil, nil, true)
		}
		if err != nil && !isNotFound(err) {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Deleting %s: calling DropCatalog", catalogTypeName),
				errs.WrapErrorf(err, "Deleting %s: calling DropCatalog", catalogTypeName).Error(),
			)
			return
		}
	}
}

func (r *catalogResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// catalogState maps a GetCatalog response onto the resource model. The
// Terraform `id` follows the catalog `name` (the path identifier), while
// `catalog_id` exposes the server-generated id for downstream references.
func catalogState(o map[string]interface{}) catalogResourceModel {
	return catalogResourceModel{
		Id:        stringAttr(stringFromResponse(o, "name")),
		Name:      stringAttr(stringFromResponse(o, "name")),
		Type:      stringAttr(stringFromResponse(o, "type")),
		IsShared:  types.BoolValue(boolFromResponse(o, "isShared")),
		Options:   stringMapToTerraform(optionsFromAPI(o["options"])),
		ShareId:   stringAttr(stringFromResponse(o, "shareId")),
		CatalogId: stringAttr(stringFromResponse(o, "id")),
		Status:    stringAttr(stringFromResponse(o, "status")),
		Owner:     stringAttr(stringFromResponse(o, "owner")),
		CreatedAt: types.Int64Value(int64FromResponse(o, "createdAt")),
		CreatedBy: stringAttr(stringFromResponse(o, "createdBy")),
		UpdatedAt: types.Int64Value(int64FromResponse(o, "updatedAt")),
		UpdatedBy: stringAttr(stringFromResponse(o, "updatedBy")),
	}
}
