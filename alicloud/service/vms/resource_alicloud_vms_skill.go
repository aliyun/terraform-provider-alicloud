package vms

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const skillTypeName = "alicloud_vms_skill"

var (
	_ resource.Resource                = &skillResource{}
	_ resource.ResourceWithImportState = &skillResource{}
)

// NewSkillResource is the Framework factory for alicloud_vms_skill.
func NewSkillResource() resource.Resource {
	return &skillResource{}
}

type skillResource struct {
	fwadapt.ResourceBase
}

type skillResourceModel struct {
	Id           types.String `tfsdk:"id"`
	EnterpriseId types.String `tfsdk:"enterprise_id"`
	Name         types.String `tfsdk:"name"`
	Comment      types.String `tfsdk:"comment"`
	CreateTime   types.String `tfsdk:"create_time"`
}

func (r *skillResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = skillTypeName
}

func (r *skillResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The resource ID, in the format `<enterprise_id>:<skill_id>`. Use this format to import the resource.",
			},
			"enterprise_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The ID of the Voice Messaging Service enterprise (call center instance) that owns the skill. Changing this recreates the resource.",
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The name of the skill. Changing this recreates the resource.",
			},
			"comment": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The description of the skill. Omitting it leaves the API default; changing the value recreates the resource.",
			},
			"create_time": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The creation time of the skill, accurate to the second.",
			},
		},
	}
}

func (r *skillResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan skillResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client := r.Meta()
	if client == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: missing provider meta", skillTypeName),
			"The provider client is not configured for this resource.",
		)
		return
	}

	enterpriseId, err := strconv.ParseInt(plan.EnterpriseId.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: invalid enterprise_id", skillTypeName),
			fmt.Sprintf("enterprise_id %q is not a valid integer: %v", plan.EnterpriseId.ValueString(), err),
		)
		return
	}
	name := plan.Name.ValueString()
	comment := ""
	if !plan.Comment.IsNull() && !plan.Comment.IsUnknown() {
		comment = plan.Comment.ValueString()
	}

	ctx, cancel := context.WithTimeout(ctx, skillCreateTimeout)
	defer cancel()

	svc := NewService(client)
	data, err := svc.createSkill(ctx, enterpriseId, name, comment, skillCreateTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CloudCreateSkill", skillTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CloudCreateSkill", skillTypeName).Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, skillState(data))...)
}

func (r *skillResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state skillResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	enterpriseId, skillId, ok := parseSkillId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: parsing resource id", skillTypeName),
			fmt.Sprintf("expected the id in the format <enterprise_id>:<skill_id>, got %q", state.Id.ValueString()),
		)
		return
	}
	svc := NewService(r.Meta())
	object, err := svc.describeSkill(ctx, enterpriseId, skillId, skillReadTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling CloudListSkill", skillTypeName),
			errs.WrapErrorf(err, "Reading %s: calling CloudListSkill", skillTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, skillState(object))...)
}

// Update is present to satisfy the resource.Resource interface. Every writable
// argument (enterprise_id, name and comment) forces a replacement, so the
// framework routes any configuration change through Delete+Create and never
// reaches Update with an in-place diff. Because the Voice Messaging Service
// exposes no update API for a skill, Update mirrors the Read refresh path: it
// re-lists the skill by id and rebuilds the state, removing the resource when
// the skill no longer exists so a re-run apply converges.
func (r *skillResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan skillResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	enterpriseId, skillId, ok := parseSkillId(plan.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: parsing resource id", skillTypeName),
			fmt.Sprintf("expected the id in the format <enterprise_id>:<skill_id>, got %q", plan.Id.ValueString()),
		)
		return
	}
	svc := NewService(r.Meta())
	object, err := svc.describeSkill(ctx, enterpriseId, skillId, skillUpdateTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling CloudListSkill", skillTypeName),
			errs.WrapErrorf(err, "Updating %s: calling CloudListSkill", skillTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, skillState(object))...)
}

func (r *skillResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state skillResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	enterpriseId, skillId, ok := parseSkillId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: parsing resource id", skillTypeName),
			fmt.Sprintf("expected the id in the format <enterprise_id>:<skill_id>, got %q", state.Id.ValueString()),
		)
		return
	}
	svc := NewService(r.Meta())
	ctx, cancel := context.WithTimeout(ctx, skillDeleteTimeout)
	defer cancel()
	if err := svc.deleteSkill(ctx, enterpriseId, skillId, skillDeleteTimeout); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling CloudDeleteSkill", skillTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling CloudDeleteSkill", skillTypeName).Error(),
		)
		return
	}
}

func (r *skillResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// skillState maps a Dyvmsapi skill object (from CloudCreateSkill's Data or a
// CloudListSkill list item, both carrying Id, Name, Comment, EnterpriseId and
// CreateTime) onto the resource model. The composite id is rebuilt from the
// response so the parsed import id and the live object agree.
func skillState(object map[string]interface{}) skillResourceModel {
	m := skillResourceModel{
		Id:           types.StringNull(),
		EnterpriseId: types.StringNull(),
		Name:         types.StringNull(),
		Comment:      types.StringNull(),
		CreateTime:   types.StringNull(),
	}
	if eid, eidErr := toInt64(object["EnterpriseId"]); eidErr == nil {
		m.EnterpriseId = types.StringValue(strconv.FormatInt(eid, 10))
		if sid, sidErr := toInt64(object["Id"]); sidErr == nil {
			m.Id = types.StringValue(fmt.Sprintf("%d:%d", eid, sid))
		}
	}
	if v, ok := object["Name"].(string); ok {
		m.Name = types.StringValue(v)
	}
	if v, ok := object["Comment"].(string); ok {
		m.Comment = types.StringValue(v)
	}
	if v, ok := object["CreateTime"].(string); ok {
		m.CreateTime = types.StringValue(v)
	}
	return m
}
