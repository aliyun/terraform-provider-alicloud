package slb

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const listenerAclAttachmentTypeName = "alicloud_slb_listener_acl_attachment"

var (
	_ resource.Resource                = &listenerAclAttachmentResource{}
	_ resource.ResourceWithImportState = &listenerAclAttachmentResource{}
	_ validator.String                 = listenerIdAttributeValidator{}
)

func NewListenerAclAttachmentResource() resource.Resource {
	return &listenerAclAttachmentResource{}
}

type listenerAclAttachmentResource struct {
	fwadapt.ResourceBase
}

type listenerAclAttachmentResourceModel struct {
	Id         types.String `tfsdk:"id"`
	ListenerId types.String `tfsdk:"listener_id"`
	AclId      types.String `tfsdk:"acl_id"`
	AclType    types.String `tfsdk:"acl_type"`
	AclStatus  types.String `tfsdk:"acl_status"`
}

// listenerIdAttributeValidator validates the listener identity of both the
// configuration and the imported ID through the same parser, so an unsupported
// protocol or an out of range port fails locally with a located diagnostic.
type listenerIdAttributeValidator struct{}

func (listenerIdAttributeValidator) Description(context.Context) string {
	return "must be in the format of <load_balancer_id>:<protocol>:<listener_port>, protocol is one of http, https, tcp and udp, port is between 1 and 65535"
}

func (v listenerIdAttributeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (listenerIdAttributeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if _, err := parseListenerAclListenerId(req.ConfigValue.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid listener_id", err.Error())
	}
}

func (r *listenerAclAttachmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = listenerAclAttachmentTypeName
}

func (r *listenerAclAttachmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an Alicloud SLB listener access control (ACL) attachment resource. " +
			"It binds an existing access control policy group to an existing SLB listener and only manages the ACL settings of that listener.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"listener_id": schema.StringAttribute{
				Required: true,
				Description: "The ID of the SLB listener, in the format of `<load_balancer_id>:<protocol>:<listener_port>`. " +
					"It is the same as the `id` of `alicloud_slb_listener` and the `id` element of `alicloud_slb_listeners`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					listenerIdAttributeValidator{},
				},
			},
			"acl_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the access control policy group bound to the listener. Exactly one policy group is managed by this resource.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[^, \t\r\n]+$`),
						"must be a single access control policy group ID; comma separated IDs are not supported",
					),
				},
			},
			"acl_type": schema.StringAttribute{
				Required:    true,
				Description: "The access control type. Valid values: `white` and `black`.",
				Validators: []validator.String{
					stringvalidator.OneOf("white", "black"),
				},
			},
			"acl_status": schema.StringAttribute{
				Computed:    true,
				Description: "The access control status of the listener.",
			},
		},
	}
}

func (r *listenerAclAttachmentResource) api() listenerAclAttachmentApi {
	return newListenerAclAttachmentApi(r.Meta())
}

func (r *listenerAclAttachmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan listenerAclAttachmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	model, warnings, err := createListenerAclAttachment(ctx, r.api(), plan)
	for _, warning := range warnings {
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("Created %s with a warning", listenerAclAttachmentTypeName),
			warning,
		)
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s", listenerAclAttachmentTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *listenerAclAttachmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state listenerAclAttachmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	model, found, err := readListenerAclAttachment(ctx, r.api(), state)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s", listenerAclAttachmentTypeName),
			err.Error(),
		)
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *listenerAclAttachmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state listenerAclAttachmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	model, warnings, err := updateListenerAclAttachment(ctx, r.api(), plan, state)
	for _, warning := range warnings {
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("Updated %s with a warning", listenerAclAttachmentTypeName),
			warning,
		)
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s", listenerAclAttachmentTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *listenerAclAttachmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state listenerAclAttachmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	listenerId, _, err := parseListenerAclAttachmentId(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: parsing resource ID", listenerAclAttachmentTypeName),
			err.Error(),
		)
		return
	}

	// One deadline covers the pre-read, the disabling write and the status
	// polling, so a Delete cannot spend the budget three times over.
	deleteCtx, cancel := context.WithTimeout(ctx, listenerAclAttachmentOperationTimeout)
	defer cancel()
	if err := deleteListenerAclAttachment(deleteCtx, r.api(), listenerId); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s", listenerAclAttachmentTypeName),
			err.Error(),
		)
	}
}

func (r *listenerAclAttachmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	listenerId, aclId, err := parseListenerAclAttachmentId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Importing %s: parsing import ID", listenerAclAttachmentTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), listenerId)...)
	if aclId != "" {
		// Keep the requested fourth segment as the managed policy group so the
		// following read never adopts a different group.
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("acl_id"), aclId)...)
	}
}

// createListenerAclAttachment binds the policy group and returns the state that
// Terraform should keep. A write whose outcome is unknown is reconciled against
// the listener instead of being reported as a plain failure, so an applied
// binding can never stay unmanaged.
func createListenerAclAttachment(ctx context.Context, api listenerAclAttachmentApi, plan listenerAclAttachmentResourceModel) (listenerAclAttachmentResourceModel, []string, error) {
	listenerId := plan.ListenerId.ValueString()
	aclId := plan.AclId.ValueString()
	aclType := plan.AclType.ValueString()

	if err := api.SetAcl(ctx, listenerId, aclId, aclType, true); err != nil {
		if !isUnknownListenerAclWriteOutcome(err) {
			return listenerAclAttachmentResourceModel{}, nil, err
		}
		object, describeErr := api.Describe(ctx, listenerId)
		if describeErr == nil && listenerAclContains(object, aclId) {
			return listenerAclAttachmentState(listenerId, aclId, plan.AclType, object), []string{
				fmt.Sprintf("The request that binds policy group %s to listener %s returned an unknown result, but the binding is present on the listener; the identity is kept in state.", aclId, listenerId),
			}, nil
		}
		return listenerAclAttachmentResourceModel{}, nil, err
	}

	// The write is confirmed, so the identity is persisted before the first
	// read. A transient read failure must not drop a binding Terraform owns;
	// the next refresh reconciles the server values.
	planState := listenerAclAttachmentPlanState(plan)
	object, err := api.Describe(ctx, listenerId)
	if err != nil {
		return planState, []string{
			fmt.Sprintf("The policy group %s was bound to listener %s, but reading the listener failed: %s. The binding is kept in state and will be reconciled on the next refresh.", aclId, listenerId, err),
		}, nil
	}
	return listenerAclAttachmentState(listenerId, aclId, plan.AclType, object), nil, nil
}

// updateListenerAclAttachment rewrites the ACL settings only when the identity
// or the type changed, and follows the same unknown outcome and read failure
// rules as create.
func updateListenerAclAttachment(ctx context.Context, api listenerAclAttachmentApi, plan, state listenerAclAttachmentResourceModel) (listenerAclAttachmentResourceModel, []string, error) {
	listenerId := plan.ListenerId.ValueString()
	aclId := plan.AclId.ValueString()
	aclType := plan.AclType.ValueString()

	if !plan.AclId.Equal(state.AclId) || !plan.AclType.Equal(state.AclType) {
		if err := api.SetAcl(ctx, listenerId, aclId, aclType, true); err != nil {
			if !isUnknownListenerAclWriteOutcome(err) {
				return listenerAclAttachmentResourceModel{}, nil, err
			}
			object, describeErr := api.Describe(ctx, listenerId)
			if describeErr != nil || !listenerAclContains(object, aclId) {
				return listenerAclAttachmentResourceModel{}, nil, err
			}
			return listenerAclAttachmentState(listenerId, aclId, plan.AclType, object), []string{
				fmt.Sprintf("The request that updates the access control settings of listener %s returned an unknown result, but the planned binding is present on the listener; the identity is kept in state.", listenerId),
			}, nil
		}
	}

	planState := listenerAclAttachmentPlanState(plan)
	object, err := api.Describe(ctx, listenerId)
	if err != nil {
		return planState, []string{
			fmt.Sprintf("The access control settings of listener %s were updated, but reading the listener failed: %s. The planned state is kept and will be reconciled on the next refresh.", listenerId, err),
		}, nil
	}
	return listenerAclAttachmentState(listenerId, aclId, plan.AclType, object), nil, nil
}

// readListenerAclAttachment resolves the managed policy group and reports
// whether the resource still exists.
func readListenerAclAttachment(ctx context.Context, api listenerAclAttachmentApi, state listenerAclAttachmentResourceModel) (listenerAclAttachmentResourceModel, bool, error) {
	listenerId, idAclId, err := parseListenerAclAttachmentId(state.Id.ValueString())
	if err != nil {
		return listenerAclAttachmentResourceModel{}, false, err
	}
	targetAclId := state.AclId.ValueString()
	if targetAclId == "" {
		targetAclId = idAclId
	}
	if targetAclId != "" {
		if err := validateSingleListenerAclId(targetAclId); err != nil {
			return listenerAclAttachmentResourceModel{}, false, fmt.Errorf(
				"invalid acl_id in state for listener %s: %v; the state holds a combined ACL value that this resource cannot manage, remove the resource from state and import it again with <load_balancer_id>:<protocol>:<listener_port>[:<acl_id>]",
				listenerId, err,
			)
		}
	}

	object, err := api.Describe(ctx, listenerId)
	if err != nil {
		if errs.NotFoundError(err) {
			return listenerAclAttachmentResourceModel{}, false, nil
		}
		return listenerAclAttachmentResourceModel{}, false, err
	}

	// SLB has no dedicated de-association API: the binding is considered gone
	// once access control is disabled or the configured policy group is no
	// longer part of the listener ACLs.
	if status, _ := object["AclStatus"].(string); !strings.EqualFold(status, "on") {
		return listenerAclAttachmentResourceModel{}, false, nil
	}

	if targetAclId == "" {
		// Imported without the ACL part. A single bound policy group is
		// unambiguous; several groups need the four segment import form so the
		// target identity is never guessed.
		aclIds := listenerAclIds(object)
		switch len(aclIds) {
		case 0:
			return listenerAclAttachmentResourceModel{}, false, nil
		case 1:
			targetAclId = aclIds[0]
		default:
			return listenerAclAttachmentResourceModel{}, false, fmt.Errorf(
				"listener %s has multiple access control policy groups (%s); import with the four segment ID <load_balancer_id>:<protocol>:<listener_port>:<acl_id> to select the managed group",
				listenerId, strings.Join(aclIds, ", "),
			)
		}
	} else if !listenerAclContains(object, targetAclId) {
		return listenerAclAttachmentResourceModel{}, false, nil
	}

	return listenerAclAttachmentState(listenerId, targetAclId, state.AclType, object), true, nil
}

// deleteListenerAclAttachment only disables access control and never deletes or
// rebuilds the listener. The caller supplies the single operation deadline.
func deleteListenerAclAttachment(ctx context.Context, api listenerAclAttachmentApi, listenerId string) error {
	object, err := api.Describe(ctx, listenerId)
	if err != nil {
		if errs.NotFoundError(err) {
			// The listener is gone, so the ACL binding no longer exists.
			return nil
		}
		return err
	}
	if status, _ := object["AclStatus"].(string); !strings.EqualFold(status, "on") {
		return nil
	}
	if err := api.SetAcl(ctx, listenerId, "", "", false); err != nil {
		return err
	}

	// Wait until the access control status is switched off. The server may keep
	// the previous AclId after disabling, so only AclStatus is asserted.
	return retry.RetryContext(ctx, listenerAclAttachmentOperationTimeout, func() *retry.RetryError {
		if err := ctx.Err(); err != nil {
			return retry.NonRetryableError(err)
		}
		current, err := api.Describe(ctx, listenerId)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil
			}
			return retry.NonRetryableError(err)
		}
		if status, _ := current["AclStatus"].(string); !strings.EqualFold(status, "on") {
			return nil
		}
		if err := waitListenerAclWithContext(ctx, listenerAclAttachmentPollInterval); err != nil {
			return retry.NonRetryableError(err)
		}
		return retry.RetryableError(fmt.Errorf("%s: access control of listener %s is still enabled", listenerAclAttachmentTypeName, listenerId))
	})
}

// listenerAclAttachmentPlanState stores the identity that is already known
// once the write succeeds. It is used before the first read so that a transient
// read failure can never leave an unmanaged binding behind; the next refresh
// replaces the assumed status with the server values.
func listenerAclAttachmentPlanState(plan listenerAclAttachmentResourceModel) listenerAclAttachmentResourceModel {
	return listenerAclAttachmentResourceModel{
		Id:         types.StringValue(plan.ListenerId.ValueString()),
		ListenerId: plan.ListenerId,
		AclId:      plan.AclId,
		AclType:    plan.AclType,
		AclStatus:  types.StringValue("on"),
	}
}

func listenerAclAttachmentState(listenerId, aclId string, aclType types.String, object map[string]interface{}) listenerAclAttachmentResourceModel {
	model := listenerAclAttachmentResourceModel{
		Id:         types.StringValue(listenerId),
		ListenerId: types.StringValue(listenerId),
		AclId:      types.StringValue(aclId),
		AclType:    aclType,
	}
	if value, ok := object["AclType"].(string); ok && value != "" {
		model.AclType = types.StringValue(value)
	}
	if value, ok := object["AclStatus"].(string); ok && value != "" {
		model.AclStatus = types.StringValue(value)
	} else {
		model.AclStatus = types.StringNull()
	}
	return model
}
