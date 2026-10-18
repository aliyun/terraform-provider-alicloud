package dataworks

import (
	"context"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
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

const resourceTypeName = "alicloud_data_works_resource"

var (
	_ resource.Resource                = &resourceResource{}
	_ resource.ResourceWithImportState = &resourceResource{}
)

// NewResourceResource is the Framework factory the registry wraps. The real
// CRUD lives in package-level functions so the unit tests can drive them with a
// credential-free httptest mock client (the embedded fwadapt client field is
// unexported and only arrives through provider Configure).
func NewResourceResource() resource.Resource {
	return &resourceResource{}
}

type resourceResource struct {
	fwadapt.ResourceBase
}

// resourceModel is the flattened Framework schema model. The nested DataSource
// and Script objects of the SDKv2 schema are flattened to top-level computed
// attributes: Framework disallows setting Attributes on a list element after the
// list itself has been set to a non-null value, which the SDKv2 single-element
// block pattern relied on. Flattening removes that hazard and matches how the
// data source already exposes the same fields.
type resourceModel struct {
	Id               types.String `tfsdk:"id"`
	ProjectId        types.String `tfsdk:"project_id"`
	ResourceId       types.String `tfsdk:"resource_id"`
	ResourceName     types.String `tfsdk:"resource_name"`
	Spec             types.String `tfsdk:"spec"`
	ResourceFile     types.String `tfsdk:"resource_file"`
	Path             types.String `tfsdk:"path"`
	CreateTime       types.Int64  `tfsdk:"create_time"`
	ModifyTime       types.Int64  `tfsdk:"modify_time"`
	Owner            types.String `tfsdk:"owner"`
	Type             types.String `tfsdk:"type"`
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

func (r *resourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = resourceTypeName
}

func (r *resourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a DataWorks Resource resource. The resource file is used when computing engine UDFs or scheduling tasks run.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. Formats as `<project_id>:<resource_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the DataWorks workspace to which the resource file belongs.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"resource_id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"spec": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "The FlowSpec definition of the resource file in JSON format.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_file": schema.StringAttribute{
				Optional:    true,
				Description: "The resource file, corresponds to the file stream or OSS download URL. Write-only: not returned on Read.",
			},
			"path": schema.StringAttribute{
				Optional:    true,
				Description: "The destination path of the resource. Write-only: not returned on Read.",
			},
			"create_time": schema.Int64Attribute{
				Computed:    true,
				Description: "The creation time of the resource.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"modify_time": schema.Int64Attribute{
				Computed:    true,
				Description: "The last modification time of the resource.",
			},
			"owner": schema.StringAttribute{
				Computed:    true,
				Description: "The responsible person of the resource file.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The resource type. Valid values: `jar`, `python`, `file`, `archive`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_type": schema.StringAttribute{
				Computed:    true,
				Description: "The file resource source storage type. Valid values: `local`, `oss`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_path": schema.StringAttribute{
				Computed:    true,
				Description: "The file source path.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"target_type": schema.StringAttribute{
				Computed:    true,
				Description: "The file destination storage type. Valid values: `gateway`, `oss`, `hdfs`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"target_path": schema.StringAttribute{
				Computed:    true,
				Description: "The file destination storage path.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"data_source_name": schema.StringAttribute{
				Computed:    true,
				Description: "The data source name (flattened from `data_source.name`).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"data_source_type": schema.StringAttribute{
				Computed:    true,
				Description: "The data source type (flattened from `data_source.type`).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"script_path": schema.StringAttribute{
				Computed:    true,
				Description: "The script path (flattened from `script.path`).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"script_runtime_command": schema.StringAttribute{
				Computed:    true,
				Description: "The script runtime command (flattened from `script.runtime_command`).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"script_id": schema.StringAttribute{
				Computed:    true,
				Description: "The script ID (flattened from `script.script_id`).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *resourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	out, err := resourceCreate(r.Client(), plan)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s", resourceTypeName),
			errs.WrapErrorf(err, "Creating %s", resourceTypeName).Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, out)...)
}

func (r *resourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	out, gone, err := resourceRead(r.Client(), state)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s", resourceTypeName),
			errs.WrapErrorf(err, "Reading %s", resourceTypeName).Error(),
		)
		return
	}
	if gone {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, out)...)
}

func (r *resourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state resourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	out, err := resourceUpdate(r.Client(), plan, state)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s", resourceTypeName),
			errs.WrapErrorf(err, "Updating %s", resourceTypeName).Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, out)...)
}

func (r *resourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := resourceDelete(r.Client(), state); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s", resourceTypeName),
			errs.WrapErrorf(err, "Deleting %s", resourceTypeName).Error(),
		)
		return
	}
}

// ImportState takes the composite id "<project_id>:<resource_id>". F2 fix: the
// id carries the project id, so Read can recover it without a config-only
// lookup that the previous SDKv2 round lost.
func (r *resourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// resourceCreate runs CreateResource (+ optional RenameResource for F3) and
// returns the refreshed state. Write-only fields (resource_file, path) are
// carried from the plan into state because the API never echoes them.
func resourceCreate(client *connectivity.AliyunClient, plan resourceModel) (resourceModel, error) {
	projectID := plan.ProjectId.ValueString()
	resourceID, err := createDataWorksResource(client, projectID, plan.Spec.ValueString(), plan.ResourceFile.ValueString())
	if err != nil {
		return resourceModel{}, err
	}
	// F3: resource_name is not a CreateResource parameter; it is applied via a
	// separate RenameResource call right after creation. The previous execute
	// skipped this, so the planned name silently never reached the server.
	if name := plan.ResourceName.ValueString(); name != "" {
		if err := renameDataWorksResource(client, projectID, resourceID, name); err != nil {
			return resourceModel{}, err
		}
	}
	object, err := getResource(client, projectID, resourceID)
	if err != nil {
		return resourceModel{}, err
	}
	if object == nil {
		return resourceModel{}, fmt.Errorf("the DataWorks Resource %s:%s is not visible via GetResource right after creation", projectID, resourceID)
	}
	return resourceState(object, plan.ResourceFile, plan.Path), nil
}

// resourceRead parses the composite id (F2), refreshes via GetResource, and
// preserves the write-only resource_file/path from the incoming state — the API
// never returns them, so echoing null would drift plan-vs-state on every refresh.
func resourceRead(client *connectivity.AliyunClient, state resourceModel) (resourceModel, bool, error) {
	projectID, resourceID, err := splitResourceID(state.Id.ValueString())
	if err != nil {
		return resourceModel{}, false, err
	}
	object, err := getResource(client, projectID, resourceID)
	if err != nil {
		return resourceModel{}, false, err
	}
	if object == nil {
		return resourceModel{}, true, nil
	}
	return resourceState(object, state.ResourceFile, state.Path), false, nil
}

// resourceUpdate applies spec/resource_file (UpdateResource), resource_name
// (RenameResource) and path (MoveResource). F4 fix: write-only fields
// (resource_file, path) are only sent when the new plan value is non-empty AND
// differs from state, so removing them from config converges instead of
// re-firing Update/Move every plan.
func resourceUpdate(client *connectivity.AliyunClient, plan, state resourceModel) (resourceModel, error) {
	projectID, resourceID, err := splitResourceID(state.Id.ValueString())
	if err != nil {
		return resourceModel{}, err
	}

	sendSpec := !plan.Spec.Equal(state.Spec)
	sendFile := !plan.ResourceFile.Equal(state.ResourceFile) && plan.ResourceFile.ValueString() != ""
	if sendSpec || sendFile {
		if err := updateDataWorksResource(client, projectID, resourceID,
			plan.Spec.ValueString(), sendSpec,
			plan.ResourceFile.ValueString(), sendFile); err != nil {
			return resourceModel{}, err
		}
	}

	if !plan.ResourceName.Equal(state.ResourceName) && plan.ResourceName.ValueString() != "" {
		if err := renameDataWorksResource(client, projectID, resourceID, plan.ResourceName.ValueString()); err != nil {
			return resourceModel{}, err
		}
	}

	if !plan.Path.Equal(state.Path) && plan.Path.ValueString() != "" {
		if err := moveDataWorksResource(client, projectID, resourceID, plan.Path.ValueString()); err != nil {
			return resourceModel{}, err
		}
	}

	object, err := getResource(client, projectID, resourceID)
	if err != nil {
		return resourceModel{}, err
	}
	if object == nil {
		return resourceModel{}, fmt.Errorf("the DataWorks Resource %s:%s vanished during update", projectID, resourceID)
	}
	// Write-only fields come from the plan (the values just applied); computed
	// fields come from the API refresh.
	return resourceState(object, plan.ResourceFile, plan.Path), nil
}

// resourceDelete removes the resource; not-found (code 5801488441564725318 or a
// 404) is treated as success so destroy is idempotent.
func resourceDelete(client *connectivity.AliyunClient, state resourceModel) error {
	projectID, resourceID, err := splitResourceID(state.Id.ValueString())
	if err != nil {
		return err
	}
	return deleteDataWorksResource(client, projectID, resourceID)
}

// resourceState maps a GetResource Resource object onto the model. Empty strings
// and missing nested objects read back as null. The write-only resource_file and
// path are carried through from the supplied preserved values (plan or state).
func resourceState(object map[string]interface{}, preservedFile, preservedPath types.String) resourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	dataSource := nestedMap(object, "DataSource")
	script := nestedMap(object, "Script")
	return resourceModel{
		Id:               types.StringValue(fmt.Sprintf("%s:%s", asString(object["ProjectId"]), asString(object["Id"]))),
		ProjectId:        str("ProjectId"),
		ResourceId:       str("Id"),
		ResourceName:     str("Name"),
		Spec:             str("Spec"),
		ResourceFile:     preservedFile,
		Path:             preservedPath,
		CreateTime:       types.Int64Value(asInt64(object["CreateTime"])),
		ModifyTime:       types.Int64Value(asInt64(object["ModifyTime"])),
		Owner:            str("Owner"),
		Type:             str("Type"),
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

// nestedMap returns object[key] as a map[string]interface{}, or nil.
func nestedMap(object map[string]interface{}, key string) map[string]interface{} {
	if v, ok := object[key].(map[string]interface{}); ok {
		return v
	}
	return nil
}

// stringOrZero reads key from a nested map and returns its string form as a
// non-null types.String (the field is absent → empty string, matching the
// SDKv2 d.Set which stored the raw interface{} value, never null).
func stringOrZero(m map[string]interface{}, key string) types.String {
	if m == nil {
		return types.StringValue("")
	}
	return types.StringValue(asString(m[key]))
}

// runtimeCommand pulls Script.Runtime.Command, matching the SDKv2 mapping.
func runtimeCommand(script map[string]interface{}) (map[string]interface{}, string) {
	if runtime, ok := script["Runtime"].(map[string]interface{}); ok {
		return runtime, "Command"
	}
	return nil, ""
}
