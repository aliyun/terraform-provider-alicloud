package slb

import (
	"context"
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const listenerAclAttachmentTypeName = "alicloud_slb_listener_acl_attachment"

var (
	_ resource.Resource                = &listenerAclAttachmentResource{}
	_ resource.ResourceWithImportState = &listenerAclAttachmentResource{}
	_ validator.String                 = listenerIdAttributeValidator{}
	_ validator.String                 = aclIdAttributeValidator{}
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
	AclPending types.Bool   `tfsdk:"acl_pending"`
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

// aclIdAttributeValidator rejects a combined ACL value locally: the resource
// manages exactly one policy group.
type aclIdAttributeValidator struct{}

func (aclIdAttributeValidator) Description(context.Context) string {
	return "must be a single access control policy group ID; comma separated IDs are not supported"
}

func (v aclIdAttributeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (aclIdAttributeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if err := validateSingleListenerAclId(req.ConfigValue.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid acl_id", err.Error())
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
					aclIdAttributeValidator{},
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
				Description: "The access control status of the listener, as reported by the last successful read.",
			},
			"acl_pending": schema.BoolAttribute{
				Computed: true,
				Description: "Whether the last access control write has not been confirmed by the listener yet. " +
					"It is true after a write whose outcome could not be verified: the planned identity is kept in state, " +
					"the refresh never drops it while it is pending, and a destroy only releases the binding that is " +
					"actually attributed to this resource.",
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
		// An acknowledged write whose readback contradicts the plan returns the
		// planned identity with the error: Terraform keeps that object and the next
		// refresh or destroy can reconcile it, so the binding is never orphaned.
		if listenerAclAttachmentIdentityKnown(model) {
			resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
		}
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
		// Same rule as create: the planned identity is persisted with the error so
		// the acknowledged change can be reconciled later.
		if listenerAclAttachmentIdentityKnown(model) {
			resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
		}
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

	listenerId, idAclId, err := parseListenerAclAttachmentId(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: parsing resource ID", listenerAclAttachmentTypeName),
			err.Error(),
		)
		return
	}
	aclId := state.AclId.ValueString()
	if aclId == "" {
		aclId = idAclId
	}

	if err := deleteListenerAclAttachment(ctx, r.api(), listenerId, aclId, state.AclPending.ValueBool()); err != nil {
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
	planState := listenerAclAttachmentPendingState(plan)

	writeCtx, cancelWrite := context.WithTimeout(ctx, listenerAclAttachmentOperationTimeout)
	defer cancelWrite()
	if err := api.SetAcl(writeCtx, listenerId, aclId, aclType, true); err != nil {
		if !isUnconfirmedListenerAclWrite(writeCtx, err) {
			// The service refused the request, so no binding was created and no
			// state is written.
			return listenerAclAttachmentResourceModel{}, nil, err
		}
		// A parent deadline or cancellation stops the reconciliation request, but
		// the write may already be on the wire. The planned identity is still kept,
		// as pending, instead of dropping a binding that could exist.
		if ctx.Err() == nil {
			// The request may have been applied. Confirm the enabled status, the
			// membership and the planned type before adopting the binding. The
			// readback gets its own bounded window so an eventually consistent
			// response does not drop an applied binding.
			confirmCtx, cancel := context.WithTimeout(ctx, listenerAclAttachmentConfirmTimeout)
			defer cancel()
			object, confirmed, _, _ := waitListenerAclAttachmentConverged(confirmCtx, api, listenerId, aclId, aclType)
			if confirmed {
				return listenerAclAttachmentState(listenerId, aclId, plan.AclType, object), []string{
					fmt.Sprintf("The request that binds policy group %s to listener %s returned an unknown result, but the listener reports the planned binding; the identity is kept in state.", aclId, listenerId),
				}, nil
			}
		}
		// The request may still have reached the service. Keep the planned identity
		// as pending so a later refresh never drops it and a destroy can release an
		// applied binding instead of leaving an enabled ACL unmanaged; the
		// diagnostic keeps the failure visible.
		return planState, nil, fmt.Errorf(
			"the request that binds policy group %s to listener %s returned an unknown result and the listener did not confirm the planned access control settings within %s; the planned identity is kept in state as pending (acl_pending = true) and will be reconciled on the next refresh: %v",
			aclId, listenerId, listenerAclAttachmentConfirmTimeout, err,
		)
	}

	// The write was acknowledged, so the identity is persisted before the first
	// read. The read waits for the listener to report the planned settings, so a
	// stale response cannot store the previous type against the new plan.
	readCtx, cancelRead := context.WithTimeout(ctx, listenerAclAttachmentOperationTimeout)
	defer cancelRead()
	object, converged, complete, err := waitListenerAclAttachmentConverged(readCtx, api, listenerId, aclId, aclType)
	if err != nil {
		return planState, []string{
			fmt.Sprintf("The policy group %s was bound to listener %s, but reading the listener failed: %s. The planned identity is kept in state as pending and will be reconciled on the next refresh.", aclId, listenerId, err),
		}, nil
	}
	if !converged {
		if complete {
			// The listener answered completely and still reports a different
			// binding: this is a contradiction, not a pending readback, so it must
			// not be reported as a successful create. The planned identity travels
			// back with the error so Terraform keeps a recoverable object: without
			// it an applied binding would stay unmanaged, and without the planned
			// fields a later refresh would reconcile against the old group.
			return planState, nil, fmt.Errorf(
				"listener %s was asked to bind policy group %s with type %s, but it reports a different access control binding (status %v, type %v, groups %v)",
				listenerId, aclId, aclType, object["AclStatus"], object["AclType"], listenerAclIdsForMessage(object),
			)
		}
		return planState, []string{
			fmt.Sprintf("The policy group %s was bound to listener %s, but the listener did not report the planned access control settings within the read budget. The planned identity is kept as pending and will be reconciled on the next refresh.", aclId, listenerId),
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
		planState := listenerAclAttachmentPendingState(plan)
		writeCtx, cancelWrite := context.WithTimeout(ctx, listenerAclAttachmentOperationTimeout)
		defer cancelWrite()
		if err := api.SetAcl(writeCtx, listenerId, aclId, aclType, true); err != nil {
			if !isUnconfirmedListenerAclWrite(writeCtx, err) {
				return listenerAclAttachmentResourceModel{}, nil, err
			}
			if ctx.Err() == nil {
				confirmCtx, cancel := context.WithTimeout(ctx, listenerAclAttachmentConfirmTimeout)
				defer cancel()
				object, confirmed, _, _ := waitListenerAclAttachmentConverged(confirmCtx, api, listenerId, aclId, aclType)
				if confirmed {
					return listenerAclAttachmentState(listenerId, aclId, plan.AclType, object), []string{
						fmt.Sprintf("The request that updates the access control settings of listener %s returned an unknown result, but the listener reports the planned binding; the identity is kept in state.", listenerId),
					}, nil
				}
			}
			// Same rule as create: keep the planned identity of a write that may
			// have reached the service, as pending, so the change stays recoverable
			// even when the parent deadline stopped the reconciliation request.
			return planState, nil, fmt.Errorf(
				"the request that updates the access control settings of listener %s returned an unknown result and the listener did not confirm the planned settings within %s; the planned identity is kept in state as pending (acl_pending = true) and will be reconciled on the next refresh: %v",
				listenerId, listenerAclAttachmentConfirmTimeout, err,
			)
		}

		// The write was acknowledged: wait for the listener to report the
		// planned settings instead of storing a stale response.
		readCtx, cancelRead := context.WithTimeout(ctx, listenerAclAttachmentOperationTimeout)
		defer cancelRead()
		object, converged, complete, err := waitListenerAclAttachmentConverged(readCtx, api, listenerId, aclId, aclType)
		if err != nil {
			return planState, []string{
				fmt.Sprintf("The access control settings of listener %s were updated, but reading the listener failed: %s. The planned state is kept as pending and will be reconciled on the next refresh.", listenerId, err),
			}, nil
		}
		if !converged {
			if complete {
				// Same rule as create: the acknowledged write stays recoverable, so
				// the planned identity is returned with the error instead of leaving
				// the previous group as the only state.
				return planState, nil, fmt.Errorf(
					"listener %s was asked to use policy group %s with type %s, but it reports a different access control binding (status %v, type %v, groups %v)",
					listenerId, aclId, aclType, object["AclStatus"], object["AclType"], listenerAclIdsForMessage(object),
				)
			}
			return planState, []string{
				fmt.Sprintf("The access control settings of listener %s were updated, but the listener did not report the planned settings within the read budget. The planned state is kept as pending and will be reconciled on the next refresh.", listenerId),
			}, nil
		}
		return listenerAclAttachmentState(listenerId, aclId, plan.AclType, object), nil, nil
	}

	// No ACL field changed, so a single read is enough to refresh the server
	// values.
	object, err := api.Describe(ctx, listenerId)
	if err != nil {
		// No write is outstanding, so the previous state is kept unchanged. Marking
		// it pending here would wrongly stop a later refresh from removing an object
		// that access control no longer binds.
		return state, []string{
			fmt.Sprintf("Reading the access control settings of listener %s failed: %s. The previous state is kept and will be refreshed again later.", listenerId, err),
		}, nil
	}
	if state.AclPending.ValueBool() && !listenerAclAttachmentConfirms(object, aclId, state.AclType.ValueString()) {
		// The outstanding write is still not confirmed by this response: keeping the
		// pending marker stops a single read from settling an unproven change.
		return state, nil, nil
	}
	return listenerAclAttachmentState(listenerId, aclId, plan.AclType, object), nil, nil
}

// readListenerAclAttachment resolves the managed policy group and reports
// whether the resource still exists.
func readListenerAclAttachment(ctx context.Context, api listenerAclAttachmentApi, state listenerAclAttachmentResourceModel) (listenerAclAttachmentResourceModel, bool, error) {
	ctx, cancel := context.WithTimeout(ctx, listenerAclAttachmentOperationTimeout)
	defer cancel()
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
		if listenerAclListenerGone(err) {
			return listenerAclAttachmentResourceModel{}, false, nil
		}
		return listenerAclAttachmentResourceModel{}, false, err
	}

	// SLB has no dedicated de-association API: the binding is considered gone
	// only when access control is explicitly off, the configured policy group is
	// no longer part of the listener ACLs, or the listener is gone. An incomplete
	// response proves none of that, so the existing state is kept instead of
	// dropping a resource that may still be bound.
	//
	// An import carries no listener_id yet: there is no existing state to keep,
	// and a partial imported resource would silently claim a binding that was
	// never proven. Such a read reports the pending confirmation instead.
	importing := state.ListenerId.IsNull()
	// A pending write may still be applied by the service. While it is pending the
	// identity is never dropped by an observation that merely does not show the
	// planned binding: first read lag, a stale group and a still disabled listener
	// all leave the outcome unproven, and dropping the object would orphan an
	// applied binding. The marker clears as soon as a complete response proves the
	// planned binding.
	pending := state.AclPending.ValueBool()
	status, statusKnown := listenerAclStatus(object)
	if !statusKnown {
		if importing {
			return listenerAclAttachmentResourceModel{}, false, fmt.Errorf(
				"cannot confirm the imported binding of listener %s: the describe response did not report a usable AclStatus; retry the import once the listener reports its access control status",
				listenerId,
			)
		}
		return state, true, nil
	}
	if strings.EqualFold(status, "off") {
		if pending {
			return state, true, nil
		}
		return listenerAclAttachmentResourceModel{}, false, nil
	}
	if !strings.EqualFold(status, "on") {
		if importing {
			return listenerAclAttachmentResourceModel{}, false, fmt.Errorf(
				"cannot confirm the imported binding of listener %s: the describe response reported AclStatus %q",
				listenerId, status,
			)
		}
		return state, true, nil
	}

	aclIds, membersReported, membersComplete := listenerAclIds(object)
	if !membersReported || !membersComplete {
		// An enabled listener always reports its members, so a response that omits
		// them, reports them unreadably or reports an empty set has not reported
		// the membership: it cannot prove that the managed group was removed.
		// Keeping the state avoids dropping a live or pending binding on a gateway
		// that serialises the list differently.
		if importing {
			return listenerAclAttachmentResourceModel{}, false, fmt.Errorf(
				"cannot confirm the imported binding of listener %s: the describe response did not report a usable access control member list",
				listenerId,
			)
		}
		return state, true, nil
	}

	// While access control is on the Describe contract reports AclType. An import
	// has no configured type to fall back on, so a response without one cannot
	// prove the imported binding and must not seed a null acl_type.
	if importing {
		if value, ok := object["AclType"].(string); !ok || strings.TrimSpace(value) == "" {
			return listenerAclAttachmentResourceModel{}, false, fmt.Errorf(
				"cannot confirm the imported binding of listener %s: the describe response did not report a usable AclType while access control is on",
				listenerId,
			)
		}
	}

	if targetAclId == "" {
		// Imported without the ACL part. A single bound policy group is
		// unambiguous; several groups need the four segment import form so the
		// target identity is never guessed. An empty member list never reaches this
		// point: it is reported as unreadable above.
		if len(aclIds) > 1 {
			return listenerAclAttachmentResourceModel{}, false, fmt.Errorf(
				"listener %s has multiple access control policy groups (%s); import with the four segment ID <load_balancer_id>:<protocol>:<listener_port>:<acl_id> to select the managed group",
				listenerId, strings.Join(aclIds, ", "),
			)
		}
		if len(aclIds) == 0 {
			// Defensive: an empty list is classified as unreadable above, and this
			// path must never panic or drop a live object if that classification
			// changes.
			return state, true, nil
		}
		targetAclId = aclIds[0]
	} else if !listenerAclContains(object, targetAclId) {
		if pending {
			// The old group or another group is still bound. That does not prove the
			// pending write failed, so the managed identity stays pending instead of
			// being removed from state.
			return state, true, nil
		}
		return listenerAclAttachmentResourceModel{}, false, nil
	}

	return listenerAclAttachmentState(listenerId, targetAclId, state.AclType, object), true, nil
}

// deleteListenerAclAttachment only disables access control and never deletes or
// rebuilds the listener. The caller supplies the single operation deadline.
//
// When the last write is still pending, the caller cannot prove that the binding
// on the listener belongs to this resource. The release is then attributed: the
// managed policy group has to be reported as a member before access control is
// disabled, so destroying a resource whose write never landed cannot switch off
// a group that another resource manages.
func deleteListenerAclAttachment(ctx context.Context, api listenerAclAttachmentApi, listenerId, aclId string, pending bool) error {
	// A single deadline covers the pre-read, the disabling write and the
	// status polling, so a Delete cannot spend the budget three times over.
	ctx, cancel := context.WithTimeout(ctx, listenerAclAttachmentOperationTimeout)
	defer cancel()
	object, err := api.Describe(ctx, listenerId)
	if err != nil {
		if listenerAclListenerGone(err) {
			// The listener is gone, so the ACL binding no longer exists.
			return nil
		}
		return err
	}
	status, statusKnown := listenerAclStatus(object)
	if !statusKnown {
		return fmt.Errorf("cannot confirm that listener %s no longer uses access control: the describe response did not report a usable AclStatus, so the binding is not proven to be removed", listenerId)
	}
	if strings.EqualFold(status, "off") {
		return nil
	}
	if !strings.EqualFold(status, "on") {
		return fmt.Errorf("cannot confirm that listener %s no longer uses access control: the describe response reported AclStatus %q", listenerId, status)
	}
	if pending {
		_, membersReported, membersComplete := listenerAclIds(object)
		if !membersReported || !membersComplete {
			return fmt.Errorf("cannot attribute the pending access control binding of listener %s: the describe response did not report a usable access control member list, so disabling it could release a group this resource does not manage", listenerId)
		}
		if aclId == "" || !listenerAclContains(object, aclId) {
			// The planned group is not bound, so the unconfirmed write never landed:
			// this resource has nothing to release, and disabling access control
			// would switch off another group.
			return nil
		}
	}
	if err := api.SetAcl(ctx, listenerId, "", "", false); err != nil {
		return err
	}

	// Wait until the access control status is switched off. The server may keep
	// the previous AclId after disabling, so only AclStatus is asserted. The
	// polling runs on the caller goroutine and stops with the operation budget.
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		current, err := api.Describe(ctx, listenerId)
		if err != nil {
			if listenerAclListenerGone(err) {
				return nil
			}
			return err
		}
		status, statusKnown := listenerAclStatus(current)
		if !statusKnown {
			return fmt.Errorf("cannot confirm that listener %s no longer uses access control: the describe response did not report a usable AclStatus, so the binding is not proven to be removed", listenerId)
		}
		if strings.EqualFold(status, "off") {
			return nil
		}
		if !strings.EqualFold(status, "on") {
			return fmt.Errorf("cannot confirm that listener %s no longer uses access control: the describe response reported AclStatus %q", listenerId, status)
		}
		if err := waitListenerAclWithContext(ctx, listenerAclAttachmentPollInterval); err != nil {
			return err
		}
	}
}

// listenerAclAttachmentIdentityKnown reports whether the model carries the
// identity of a binding that the service already accepted, so it can be stored
// even when the readback contradicts the plan. An empty model means the write
// was not acknowledged and must not create state.
func listenerAclAttachmentIdentityKnown(model listenerAclAttachmentResourceModel) bool {
	return !model.ListenerId.IsNull() || !model.Id.IsNull()
}

// listenerAclAttachmentPendingState stores the identity of a write whose
// confirmation is still outstanding. The planned fields are kept so a later
// refresh can reconcile or destroy the binding, the status stays unknown because
// the listener never confirmed it, and acl_pending marks the object so the
// refresh does not drop it when it still observes the old group or a disabled
// listener.
func listenerAclAttachmentPendingState(plan listenerAclAttachmentResourceModel) listenerAclAttachmentResourceModel {
	return listenerAclAttachmentResourceModel{
		Id:         types.StringValue(plan.ListenerId.ValueString()),
		ListenerId: plan.ListenerId,
		AclId:      plan.AclId,
		AclType:    plan.AclType,
		AclStatus:  types.StringNull(),
		AclPending: types.BoolValue(true),
	}
}

func listenerAclAttachmentState(listenerId, aclId string, aclType types.String, object map[string]interface{}) listenerAclAttachmentResourceModel {
	model := listenerAclAttachmentResourceModel{
		Id:         types.StringValue(listenerId),
		ListenerId: types.StringValue(listenerId),
		AclId:      types.StringValue(aclId),
		AclType:    aclType,
		AclPending: types.BoolValue(false),
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
