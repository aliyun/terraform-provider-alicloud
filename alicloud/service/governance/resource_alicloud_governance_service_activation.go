// Package governance. This file is generated automatically. Please do not modify it manually, thank you!
// Nested attributes require Terraform Plugin Protocol 6.
package governance

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const serviceActivationTypeName = "alicloud_governance_service_activation"

var (
	_ resource.Resource                = &serviceActivationResource{}
	_ resource.ResourceWithImportState = &serviceActivationResource{}
	_ resource.ResourceWithModifyPlan  = &serviceActivationResource{}
)

func NewServiceActivationResource() resource.Resource {
	return &serviceActivationResource{}
}

type serviceActivationResource struct {
	fwadapt.ResourceBase
}

type serviceActivationResourceModel struct {
	Id         types.String `tfsdk:"id"`
	FwTimeouts types.Object `tfsdk:"timeouts"`
	Region     types.String `tfsdk:"region"`
	Status     types.String `tfsdk:"status"`
}

func (r *serviceActivationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = serviceActivationTypeName
}

func (r *serviceActivationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{"timeouts": fwTimeoutBlock()},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The region in which Cloud Governance Center is activated for the account. It is the Terraform resource ID; the service is activated once per account in a single activation region and the open, decommission and status APIs must all target that region.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The activation status of Cloud Governance Center for the current account, returned by GetGovernanceServiceStatus. Enabled means activated, Disabled means not activated or closed, EnableFailed means activation failed, DisableFailed means deactivation failed, Enabling and Disabling mean the operation is in progress.",
			},
		},
	}
}

// ModifyPlan rejects a region that differs from the provider region. All three operations send
// RegionId for the activation region through the provider configuration, while region is also the
// resource ID, so a divergent value would describe one region through another region's client.
func (r *serviceActivationResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() || r.Client() == nil {
		return
	}
	var plan serviceActivationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() || plan.Region.IsNull() || plan.Region.IsUnknown() {
		return
	}
	if region := plan.Region.ValueString(); region != r.Client().RegionId {
		resp.Diagnostics.AddAttributeError(
			path.Root("region"),
			"Region must match the provider region",
			fmt.Sprintf("Cloud Governance Center must be addressed in its activation region, and this resource sends the open, decommission and status requests through the provider configuration, so region must be %q (the provider region), but %q was configured. The activation cannot move between regions: point the provider, or the provider alias applied to this resource, at the activation region %q instead of removing the resource and importing it under another region, because an import with the activation region as its ID also only succeeds while the provider is configured with that same region.", r.Client().RegionId, region, r.Client().RegionId),
		)
	}
}

func (r *serviceActivationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serviceActivationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operationTimeout, timeoutDiagnostics := fwTimeoutDuration(plan.FwTimeouts, "create", 20*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The RPC endpoint is resolved from the provider region and RegionId names the activation,
	// which is the region itself: the activation has no identity beyond the region.
	request := map[string]interface{}{
		"ClientToken": helper.BuildClientToken("OpenGovernanceService"),
		"RegionId":    r.Client().RegionId,
	}
	id := r.Client().RegionId
	if id == "" {
		resp.Diagnostics.AddError("Creating alicloud_governance_service_activation: missing identity", "The provider region is empty, so the activation cannot be addressed")
		return
	}

	// Persist the allocated identity before any asynchronous or configuration stage.
	trackedPlan := plan
	trackedPlan.Id = types.StringValue(id)
	if trackedPlan.Region.IsUnknown() {
		trackedPlan.Region = types.StringNull()
	}
	if trackedPlan.Status.IsUnknown() {
		trackedPlan.Status = types.StringNull()
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, trackedPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// OpenGovernanceService can be accepted while a previous deactivation is still finishing, after
	// which the service settles back to Disabled instead of Enabled. Retry the open/wait cycle until
	// the operation timeout so such a stale disable does not fail the create. The service is
	// account-wide and cannot be opened twice, so the status query decides whether an open that
	// returned an error still reached Enabled.
	// One deadline covers the whole create: every write, status query and wait spends from it, and no
	// write starts once it has expired.
	opCtx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()
	deadline, _ := opCtx.Deadline()
	var object map[string]interface{}
	var waitErr error
	var openErr error
	var lastErr error
	lastStatus := ""
	lastFailure := ""
	issueOpen := true
	// OpenGovernanceService publishes InvalidStatus.AlreadyEnabled when the account is already
	// activated. Querying the status first keeps create idempotent without depending on that code.
	if current, err := serviceActivationGet(opCtx, r.Client(), id); err == nil && current != nil {
		if status, ok := current["Status"].(string); ok && status == "Enabled" {
			issueOpen = false
		}
	}
	for {
		if time.Until(deadline) <= 0 {
			break
		}
		if issueOpen {
			// A write error is not final by itself (error tables checked 2026-10-24): the status
			// query shows whether the open took effect, an in-flight activation change is waited
			// out before the next attempt, and an in-flight open either lands on Enabled or
			// reports the account as already activated when it is retried. Any other state means
			// the error is terminal here, so it is reported instead of burning the timeout.
			if _, err := retryRpcPost(opCtx, r.Client(), "governance", "2021-01-20", "OpenGovernanceService", request, time.Until(deadline)); err != nil {
				openErr = err
				current, getErr := serviceActivationGet(opCtx, r.Client(), id)
				if getErr != nil || current == nil {
					resp.Diagnostics.AddError(
						fmt.Sprintf("Creating %s: calling OpenGovernanceService", serviceActivationTypeName),
						errs.WrapErrorf(err, "Creating %s: calling OpenGovernanceService", serviceActivationTypeName).Error(),
					)
					return
				}
				status, _ := current["Status"].(string)
				switch status {
				case "Enabled":
					object = current
					lastStatus = "Enabled"
				case "Enabling", "Disabling":
					// The operation is still running; the wait below settles it and retries the open.
				default:
					resp.Diagnostics.AddError(
						fmt.Sprintf("Creating %s: calling OpenGovernanceService", serviceActivationTypeName),
						errs.WrapErrorf(err, "Creating %s: calling OpenGovernanceService", serviceActivationTypeName).Error(),
					)
					return
				}
				if lastStatus == "Enabled" {
					break
				}
			}
		}
		// Recompute the wait budget after the write: the write may have consumed most of the
		// deadline, and the wait must not receive a fresh full budget of its own.
		attemptTimeout := time.Until(deadline)
		if attemptTimeout <= 0 {
			break
		}
		if attemptTimeout > 2*time.Minute {
			attemptTimeout = 2 * time.Minute
		}
		if object, waitErr = serviceActivationWaitForState(opCtx, r.Client(), id, "Status", []string{"Enabled"}, nil, attemptTimeout); waitErr == nil && object != nil {
			lastStatus = "Enabled"
			break
		}
		lastStatus, lastFailure = "", ""
		lastErr = waitErr
		if current, err := serviceActivationGet(opCtx, r.Client(), id); err == nil && current != nil {
			if status, ok := current["Status"].(string); ok {
				lastStatus = status
			}
			lastFailure = serviceActivationFailureDetail(current)
		} else if err != nil {
			lastErr = err
		}
		if lastStatus == "EnableFailed" {
			break
		}
		// Do not race an in-flight deactivation; retry the open once it settles.
		issueOpen = lastStatus != "Disabling"
	}
	if lastStatus != "Enabled" {
		stateText := lastStatus
		if stateText == "" {
			stateText = "unreadable"
		}
		detail := lastFailure
		if lastErr != nil {
			detail += fmt.Sprintf("; last error: %v", lastErr)
		}
		if openErr != nil {
			detail += fmt.Sprintf("; OpenGovernanceService error: %v", openErr)
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for Enabled", serviceActivationTypeName),
			fmt.Sprintf("The service %s did not reach Enabled within %s (last state: %q)%s", id, operationTimeout, stateText, detail),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: after create", serviceActivationTypeName),
			fmt.Sprintf("The resource %s is not visible via GetGovernanceServiceStatus right after creation", id),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, serviceActivationState(object, &plan))...)
}

func (r *serviceActivationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state serviceActivationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := serviceActivationGet(ctx, r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetGovernanceServiceStatus", serviceActivationTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetGovernanceServiceStatus", serviceActivationTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	// GetGovernanceServiceStatus keeps answering after the service is deactivated, so the model's
	// resourceNotExistCondition is applied here instead of in serviceActivationGet, which the
	// create/delete state waits still need to observe. EnableFailed is part of it: the documented
	// Open contract returns that status only after the automatic rollback completed, so no
	// activation exists to describe and a create that failed this way must not linger in state.
	if status, ok := object["Status"].(string); ok && serviceActivationAbsent(status) {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, serviceActivationState(object, &state))...)
}

// Update only ever follows a timeouts change, because every other argument replaces the
// resource. status is Computed and therefore unknown in that plan, so the live activation is read
// back instead of persisting the unknown value; the new timeouts are carried over from the plan.
func (r *serviceActivationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state serviceActivationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := serviceActivationGet(ctx, r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling GetGovernanceServiceStatus", serviceActivationTypeName),
			errs.WrapErrorf(err, "Updating %s: calling GetGovernanceServiceStatus", serviceActivationTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	if status, ok := object["Status"].(string); ok && serviceActivationAbsent(status) {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, serviceActivationState(object, &plan))...)
}

func (r *serviceActivationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state serviceActivationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operationTimeout, timeoutDiagnostics := fwTimeoutDuration(state.FwTimeouts, "delete", 30*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	// One deadline covers the whole delete: the first write, every status query and every wait spend
	// from it, and no re-issued write starts once it has expired.
	opCtx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()
	deadline, _ := opCtx.Deadline()

	// DecommissionGovernance addresses the activation by the region that is also the resource id.
	request0 := map[string]interface{}{
		"ClientToken": helper.BuildClientToken("DecommissionGovernance"),
		"RegionId":    state.Id.ValueString(),
	}
	// issueDecommission reports accepted=true when the call was accepted (or deferred until an
	// in-flight change settles) and finished=true when no activation is left to close.
	issueDecommission := func() (accepted bool, finished bool) {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			resp.Diagnostics.AddError(
				"Deleting alicloud_governance_service_activation: operation timeout",
				fmt.Sprintf("The delete budget of %s was exhausted before DecommissionGovernance could be issued", operationTimeout),
			)
			return false, false
		}
		if _, err := retryRpcPost(opCtx, r.Client(), "governance", "2021-01-20", "DecommissionGovernance", request0, remaining); err != nil {
			// A decommission write error is resolved through the status instead of the code (error
			// tables checked 2026-10-24): the target may be reached already, because the service
			// reports an activation it cannot decommission as not enabled, or an activation change
			// may still be running and safe to retry once it settles. Both are observable in the
			// status, so the write error alone does not decide the outcome.
			if current, getErr := serviceActivationGet(opCtx, r.Client(), state.Id.ValueString()); getErr == nil && current != nil {
				status, _ := current["Status"].(string)
				switch {
				case serviceActivationAbsent(status):
					// Disabled is closed, and so is EnableFailed: OpenGovernanceService documents that
					// status only after the automatic rollback completed, which is why
					// DecommissionGovernance answers InvalidStatus.NotEnabled for it as well.
					return false, true
				case status == "Enabling" || status == "Disabling":
					return true, false
				}
			}
			resp.Diagnostics.AddError("Deleting alicloud_governance_service_activation: DecommissionGovernance", err.Error())
			return false, false
		}
		return true, false
	}
	if accepted, finished := issueDecommission(); finished {
		return
	} else if !accepted {
		return
	}
	// DecommissionGovernance is slow and can report DisableFailed, after which the service may even
	// fall back to Enabled; poll for a closed activation and re-issue the deactivation instead of
	// treating the first attempt as final. EnableFailed counts as closed: a rolled-back activation
	// leaves nothing to decommission, while Enabling still waits for that rollback to finish.
	lastStatus := ""
	lastFailure := ""
	var lastErr error
	for {
		attemptTimeout := time.Until(deadline)
		if attemptTimeout <= 0 {
			break
		}
		if attemptTimeout > 5*time.Minute {
			attemptTimeout = 5 * time.Minute
		}
		if object, err := serviceActivationWaitForState(opCtx, r.Client(), state.Id.ValueString(), "Status", []string{"Disabled", "EnableFailed"}, nil, attemptTimeout); err == nil && object != nil {
			lastStatus, _ = object["Status"].(string)
			break
		} else if err != nil {
			lastErr = err
		}
		lastStatus, lastFailure = "", ""
		if current, err := serviceActivationGet(opCtx, r.Client(), state.Id.ValueString()); err == nil && current != nil {
			if status, ok := current["Status"].(string); ok {
				lastStatus = status
			}
			lastFailure = serviceActivationFailureDetail(current)
		} else if err != nil {
			lastErr = err
		}
		if serviceActivationAbsent(lastStatus) {
			break
		}
		if lastStatus == "DisableFailed" || lastStatus == "Enabled" {
			accepted, finished := issueDecommission()
			if finished {
				lastStatus = "Disabled"
				break
			}
			if !accepted {
				return
			}
		}
	}
	if !serviceActivationAbsent(lastStatus) {
		stateText := lastStatus
		if stateText == "" {
			stateText = "unreadable"
		}
		detail := lastFailure
		if lastErr != nil {
			detail += fmt.Sprintf("; last error: %v", lastErr)
		}
		resp.Diagnostics.AddError(
			"Waiting after DecommissionGovernance",
			fmt.Sprintf("The service %s did not reach Disabled within %s (last state: %q)%s", state.Id.ValueString(), operationTimeout, stateText, detail),
		)
		return
	}
}

func (r *serviceActivationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// serviceActivationState maps the Get response object onto the resource model. Response fields are
// read by their response key; a string that is present but empty stays "" instead of becoming null,
// and the timeouts block is carried over from the prior plan or state.
func serviceActivationState(object map[string]interface{}, prior *serviceActivationResourceModel) serviceActivationResourceModel {
	// Get has already merged and normalized every response mapping.
	str := func(key string) types.String {
		cur := interface{}(object)
		for _, seg := range strings.Split(key, ".") {
			m, ok := cur.(map[string]interface{})
			if !ok {
				return types.StringNull()
			}
			cur, ok = m[seg]
			if !ok {
				return types.StringNull()
			}
		}
		// A present-but-empty string is a value, not an absent field: return "" instead of null.
		// Only the !ok branches above produce null.
		if v, ok := cur.(string); ok {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	m := serviceActivationResourceModel{
		FwTimeouts: types.ObjectNull(fwTimeoutTypes()),
		Id:         str("Region"),
		Region:     str("Region"),
		Status:     str("Status"),
	}
	if prior != nil {
		m.FwTimeouts = prior.FwTimeouts
	}
	return m
}
