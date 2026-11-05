package dlfnext

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	catalogTypeName = "alicloud_dlf_next_catalog"
	// terraform-plugin-framework has no user-configurable resource timeouts.
	catalogTimeout = 5 * time.Minute
)

var (
	_ resource.Resource                = &catalogResource{}
	_ resource.ResourceWithImportState = &catalogResource{}
)

func NewCatalogResource() resource.Resource {
	return &catalogResource{}
}

type catalogResource struct {
	fwadapt.ResourceBase
}

type catalogResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Type          types.String `tfsdk:"type"`
	Options       types.Map    `tfsdk:"options"`
	RemoteOptions types.Map    `tfsdk:"remote_options"`
	IsShared      types.Bool   `tfsdk:"is_shared"`
	ShareId       types.String `tfsdk:"share_id"`
	Status        types.String `tfsdk:"status"`
	Owner         types.String `tfsdk:"owner"`
	CreatedAt     types.String `tfsdk:"created_at"`
	CreatedBy     types.String `tfsdk:"created_by"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
	UpdatedBy     types.String `tfsdk:"updated_by"`
}

func (r *catalogResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = catalogTypeName
}

func (r *catalogResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a DLF Next Catalog resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The name of the catalog. The name must be unique within the region.",
			},
			"type": schema.StringAttribute{
				Optional: true,
				Default:  stringdefault.StaticString("PAIMON"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("PAIMON", "ICEBERG"),
				},
				Description: "The catalog type. Valid values: `PAIMON`, `ICEBERG`.",
			},
			"options": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "The option keys managed by Terraform. Both `options = {}` and omitting `options` mean the managed set is empty.",
			},
			"remote_options": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "The full remote options returned by GetCatalog, including server defaults and unmanaged keys.",
			},
			"is_shared": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Description: "Whether the catalog is shared. Omitted means the field is not sent; `false` is sent explicitly.",
			},
			"share_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The share ID issued by DLF catalog sharing.",
			},
			"status":     schema.StringAttribute{Computed: true},
			"owner":      schema.StringAttribute{Computed: true},
			"created_at": schema.StringAttribute{Computed: true, Description: "Creation time in milliseconds since the epoch, as a string."},
			"created_by": schema.StringAttribute{Computed: true},
			"updated_at": schema.StringAttribute{Computed: true, Description: "Update time in milliseconds since the epoch, as a string."},
			"updated_by": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *catalogResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan catalogResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewService(r.Meta())
	name := plan.Name.ValueString()
	deadline := time.Now().Add(catalogTimeout)

	// Pre-create existence check: CreateCatalog has no idempotency token and
	// the name is unique, so an existing catalog must not be adopted.
	if _, preErr := service.DescribeCatalogWithRetryTimeout(ctx, name, remainingUntil(deadline)); preErr == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s", catalogTypeName),
			"a DlfNext Catalog with this name already exists; run `terraform import "+catalogTypeName+".<name> <name>` to manage it instead of creating it",
		)
		return
	} else if !errs.NotFoundError(preErr) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: checking the catalog name", catalogTypeName),
			errs.WrapErrorf(preErr, "checking DlfNext Catalog %s", name).Error(),
		)
		return
	}

	request, buildErr := buildCatalogCreateBody(ctx, plan)
	if buildErr != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Creating %s", catalogTypeName), buildErr.Error())
		return
	}

	action := "/dlf/v1/catalogs"
	query := make(map[string]*string)
	var response map[string]interface{}
	var err error

	// autoRetry is disabled: CreateCatalog is not idempotent and every failed
	// attempt must be reconciled through the fail-closed GetCatalog policy
	// below instead of being retried invisibly inside the ROA client.
	backoff := newCatalogBackoff(3*time.Second, 30*time.Second)
	err = retry.RetryContext(ctx, remainingUntil(deadline), func() *retry.RetryError {
		response, err = r.Meta().RoaPost("DlfNext", "2025-03-10", action, query, nil, request, false)
		if err == nil {
			return nil
		}
		if isCatalogBackendCapacityError(err) {
			return retry.NonRetryableError(fmt.Errorf("DLF backend capacity limit reached for catalog %s in this account/Region; CreateCatalog was not retried to avoid worsening the cleanup quota. Server message: %s. Wait for the backend to reap recently-created catalogs or use a different account/Region", name, err.Error()))
		}
		alreadyExists := isCatalogAlreadyExistsError(err)
		explicitClient := isCatalogExplicitClientError(err)
		retryable := needsRetry(err)
		var getObs catalogGetObservation
		if !alreadyExists && !explicitClient {
			if _, gErr := service.DescribeCatalogWithRetryTimeout(ctx, name, remainingUntil(deadline)); gErr == nil {
				getObs = catalogGetFound
			} else if errs.NotFoundError(gErr) {
				getObs = catalogGetAbsent
			} else {
				getObs = catalogGetOtherError
			}
		}
		decision, dErr := decideCatalogCreateRecovery(alreadyExists, explicitClient, retryable, getObs)
		if decision == catalogCreateRetry {
			if waitErr := backoff.Wait(ctx, remainingUntil(deadline)); waitErr != nil {
				return retry.NonRetryableError(waitErr)
			}
			return retry.RetryableError(err)
		}
		if dErr != nil {
			return retry.NonRetryableError(dErr)
		}
		return retry.NonRetryableError(err)
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s", catalogTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, catalogTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	// CreateCatalog is asynchronous. A brief NotFound is Pending (eventual
	// consistency); INITIALIZE_FAILED is terminal.
	stateConf := &retry.StateChangeConf{
		Pending: []string{"NOT_FOUND", "NEW", "INITIALIZING"},
		Target:  []string{"RUNNING"},
		Timeout: remainingUntil(deadline),
		Refresh: service.CatalogStateRefreshFuncWithApi(name, "status", []string{"INITIALIZE_FAILED"}, func(id string) (map[string]interface{}, error) {
			return service.DescribeCatalogWithRetryTimeout(ctx, id, remainingUntil(deadline))
		}),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for RUNNING", catalogTypeName),
			errs.WrapErrorf(err, "waiting for DlfNext Catalog %s to RUNNING", name).Error(),
		)
		return
	}

	model, readErr := r.readCatalog(ctx, name, plan)
	if readErr != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Reading %s after create", catalogTypeName), readErr.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *catalogResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state catalogResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	model, err := r.readCatalog(ctx, state.Id.ValueString(), state)
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Reading %s", catalogTypeName), err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *catalogResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state catalogResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Options.Equal(state.Options) {
		oldOptions, diags := stringMapFromFramework(ctx, state.Options)
		resp.Diagnostics.Append(diags...)
		newOptions, diags := stringMapFromFramework(ctx, plan.Options)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		updates, removals := buildOptionChanges(oldOptions, newOptions)
		if len(updates) > 0 || len(removals) > 0 {
			id := state.Id.ValueString()
			action := catalogAction(id)
			query := make(map[string]*string)
			body := map[string]interface{}{"updates": updates}
			if len(removals) > 0 {
				body["removals"] = removals
			}
			backoff := newCatalogBackoff(3*time.Second, 30*time.Second)
			var response map[string]interface{}
			err := retry.RetryContext(ctx, catalogTimeout, func() *retry.RetryError {
				var callErr error
				response, callErr = r.Meta().RoaPost("DlfNext", "2025-03-10", action, query, nil, body, false)
				if callErr != nil {
					if needsRetry(callErr) {
						if waitErr := backoff.Wait(ctx, catalogTimeout); waitErr != nil {
							return retry.NonRetryableError(waitErr)
						}
						return retry.RetryableError(callErr)
					}
					return retry.NonRetryableError(callErr)
				}
				return nil
			})
			logs.AddDebug(action, response, body)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Updating %s", catalogTypeName),
					errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR).Error(),
				)
				return
			}
		}
	}

	model, err := r.readCatalog(ctx, state.Id.ValueString(), plan)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Reading %s after update", catalogTypeName), err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *catalogResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state catalogResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewService(r.Meta())
	id := state.Id.ValueString()
	action := catalogAction(id)
	query := make(map[string]*string)
	deadline := time.Now().Add(catalogTimeout)

	backoff := newCatalogBackoff(3*time.Second, 30*time.Second)
	err := retry.RetryContext(ctx, remainingUntil(deadline), func() *retry.RetryError {
		var callErr error
		_, callErr = r.Meta().RoaDelete("DlfNext", "2025-03-10", action, query, nil, nil, false)
		if callErr == nil {
			return nil
		}
		// DropCatalog returns 403 while the catalog is still being created.
		if needsRetry(callErr) || isCatalogBeingCreatedError(callErr) {
			if waitErr := backoff.Wait(ctx, remainingUntil(deadline)); waitErr != nil {
				return retry.NonRetryableError(waitErr)
			}
			return retry.RetryableError(callErr)
		}
		return retry.NonRetryableError(callErr)
	})
	if err != nil {
		if errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s", catalogTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	if goneErr := service.WaitForCatalogGone(ctx, id, remainingUntil(deadline)); goneErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: waiting for the catalog to disappear", catalogTypeName),
			errs.WrapErrorf(goneErr, "waiting for DlfNext Catalog %s to be deleted", id).Error(),
		)
	}
}

func (r *catalogResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// readCatalog maps GetCatalog onto the model. managedState carries the
// currently-managed option keys so Read reconciles only those keys.
func (r *catalogResource) readCatalog(ctx context.Context, id string, managedState catalogResourceModel) (*catalogResourceModel, error) {
	object, err := NewService(r.Meta()).DescribeCatalog(ctx, id)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	managedNull := managedState.Options.IsNull()
	managed, diags := stringMapFromFramework(ctx, managedState.Options)
	if diags.HasError() {
		return nil, fmt.Errorf("reading managed options from state: %v", diags)
	}
	model, err := catalogModelFromObject(ctx, object, managed, managedNull)
	if err != nil {
		return nil, err
	}
	model.Id = types.StringValue(stringValue(object["name"]))
	return model, nil
}

// catalogModelFromObject maps the API object; remote_options mirrors all
// remote keys and options keeps only the managed keys still present remotely.
func catalogModelFromObject(ctx context.Context, object map[string]interface{}, managed map[string]string, managedNull bool) (*catalogResourceModel, error) {
	apiOptions := map[string]interface{}{}
	if raw, ok := object["options"].(map[string]interface{}); ok {
		apiOptions = raw
	}
	remote := reconcileRemoteOptions(apiOptions)
	remoteMap, diags := frameworkStringMap(ctx, remote)
	if diags.HasError() {
		return nil, fmt.Errorf("building remote_options: %v", diags)
	}
	// Preserve null when the user omitted options: framework treats null and
	// an empty map as different, and writing {} for an omitted block would
	// create a perpetual null-vs-{} diff.
	optionsMap := types.MapNull(types.StringType)
	if !managedNull {
		reconciled := reconcileManagedOptions(managed, apiOptions)
		optionsMap, diags = frameworkStringMap(ctx, reconciled)
		if diags.HasError() {
			return nil, fmt.Errorf("building options: %v", diags)
		}
	}
	model := &catalogResourceModel{
		Name:          stringOrNull(object["name"]),
		Type:          stringOrNull(object["type"]),
		Options:       optionsMap,
		RemoteOptions: remoteMap,
		IsShared:      boolOrNull(object["isShared"]),
		ShareId:       stringOrNull(object["shareId"]),
		Status:        stringOrNull(object["status"]),
		Owner:         stringOrNull(object["owner"]),
		CreatedAt:     types.StringValue(catalogTimestampString(object["createdAt"])),
		CreatedBy:     stringOrNull(object["createdBy"]),
		UpdatedAt:     types.StringValue(catalogTimestampString(object["updatedAt"])),
		UpdatedBy:     stringOrNull(object["updatedBy"]),
	}
	return model, nil
}

// buildCatalogCreateBody builds the CreateCatalog body. It is pure so the
// request shape (including the explicit-false isShared case) is unit-testable.
func buildCatalogCreateBody(ctx context.Context, plan catalogResourceModel) (map[string]interface{}, error) {
	body := map[string]interface{}{
		"name": plan.Name.ValueString(),
		"type": plan.Type.ValueString(),
	}
	options, diags := stringMapFromFramework(ctx, plan.Options)
	if diags.HasError() {
		return nil, fmt.Errorf("reading options: %v", diags)
	}
	if len(options) > 0 {
		body["options"] = options
	}
	// is_shared has no documented default: omitting it must not send the field,
	// while an explicit false must be sent as false.
	if !plan.IsShared.IsNull() && !plan.IsUnknown() {
		body["isShared"] = plan.IsShared.ValueBool()
	}
	if !plan.ShareId.IsNull() && !plan.ShareId.IsUnknown() && plan.ShareId.ValueString() != "" {
		body["shareId"] = plan.ShareId.ValueString()
	}
	return body, nil
}

// needsRetry is helper.NeedRetry behind a local name so the retry policy reads
// in one place.
func needsRetry(err error) bool {
	return helper.NeedRetry(err)
}

type catalogGetObservation int

const (
	catalogGetAbsent catalogGetObservation = iota
	catalogGetFound
	catalogGetOtherError
)

type catalogCreateDecision int

const (
	catalogCreateRetry catalogCreateDecision = iota
	catalogCreateFail
)

// decideCatalogCreateRecovery encodes the fail-closed policy: a same-name
// catalog found after an error is never adopted, because CreateCatalog has no
// idempotency token and ownership cannot be proven from attributes.
func decideCatalogCreateRecovery(postErrAlreadyExists, postErrExplicitClientError, postErrRetryable bool, getObs catalogGetObservation) (catalogCreateDecision, error) {
	if postErrAlreadyExists {
		return catalogCreateFail, fmt.Errorf("CreateCatalog returned an already-exists error; use `terraform import alicloud_dlf_next_catalog.<name> <name>` to manage it instead")
	}
	if postErrExplicitClientError {
		return catalogCreateFail, fmt.Errorf("CreateCatalog returned a validation or permission error; the catalog was not created and an existing same-name catalog must not be adopted")
	}
	if getObs == catalogGetFound {
		return catalogCreateFail, fmt.Errorf("CreateCatalog returned an ambiguous error and a catalog with this name already exists; this run cannot prove it created the catalog, so no state was written. Verify the remote catalog and run `terraform import alicloud_dlf_next_catalog.<name> <name>` if it was created by this run")
	}
	if getObs == catalogGetAbsent {
		if postErrRetryable {
			return catalogCreateRetry, nil
		}
		return catalogCreateFail, fmt.Errorf("CreateCatalog returned a non-retryable error and the catalog was not created")
	}
	return catalogCreateFail, fmt.Errorf("CreateCatalog errored and the post-error GetCatalog also failed; cannot disambiguate whether the catalog was created")
}

// catalogAction builds the catalog path with the name escaped as one path
// segment so reserved characters cannot change the request path.
func catalogAction(id string) string {
	return "/dlf/v1/catalogs/" + url.PathEscape(id)
}

// remainingUntil clamps a passed deadline to zero.
func remainingUntil(deadline time.Time) time.Duration {
	r := time.Until(deadline)
	if r < 0 {
		return 0
	}
	return r
}

// catalogTimestampString renders the int64 millisecond timestamp as a string.
func catalogTimestampString(v interface{}) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case float64:
		return strconv.FormatInt(int64(t), 10)
	case int64:
		return strconv.FormatInt(t, 10)
	case int:
		return strconv.Itoa(t)
	default:
		return fmt.Sprint(t)
	}
}

func stringValue(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

func stringOrNull(v interface{}) types.String {
	if s, ok := v.(string); ok && s != "" {
		return types.StringValue(s)
	}
	return types.StringNull()
}

func boolOrNull(v interface{}) types.Bool {
	if b, ok := v.(bool); ok {
		return types.BoolValue(b)
	}
	return types.BoolNull()
}

func stringMapFromFramework(ctx context.Context, v types.Map) (map[string]string, diag.Diagnostics) {
	out := map[string]string{}
	if v.IsNull() || v.IsUnknown() {
		return out, nil
	}
	diags := v.ElementsAs(ctx, &out, false)
	return out, diags
}

func frameworkStringMap(ctx context.Context, m map[string]string) (types.Map, diag.Diagnostics) {
	if m == nil {
		m = map[string]string{}
	}
	return types.MapValueFrom(ctx, types.StringType, m)
}

// reconcileRemoteOptions mirrors every remote option, stringified.
func reconcileRemoteOptions(apiOptions map[string]interface{}) map[string]string {
	surfaced := make(map[string]string, len(apiOptions))
	for k, v := range apiOptions {
		surfaced[k] = fmt.Sprint(v)
	}
	return surfaced
}

// reconcileManagedOptions keeps only managed keys still returned by the API;
// a key the API no longer returns is dropped so drift surfaces.
func reconcileManagedOptions(managed map[string]string, apiOptions map[string]interface{}) map[string]string {
	reconciled := make(map[string]string, len(managed))
	for k := range managed {
		if av, exists := apiOptions[k]; exists {
			reconciled[k] = fmt.Sprint(av)
		}
	}
	return reconciled
}

// buildOptionChanges derives updates from the new managed set and removals
// only from oldManaged - newManaged.
func buildOptionChanges(oldManaged, newManaged map[string]string) (map[string]interface{}, []interface{}) {
	updates := make(map[string]interface{}, len(newManaged))
	for k, v := range newManaged {
		updates[k] = v
	}
	removals := make([]interface{}, 0)
	for k := range oldManaged {
		if _, exists := newManaged[k]; !exists {
			removals = append(removals, k)
		}
	}
	return updates, removals
}
