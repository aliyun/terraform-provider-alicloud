// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudDlfNextCatalog() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDlfNextCatalogCreate,
		Read:   resourceAliCloudDlfNextCatalogRead,
		Update: resourceAliCloudDlfNextCatalogUpdate,
		Delete: resourceAliCloudDlfNextCatalogDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PAIMON",
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PAIMON"}, false),
			},
			// options is the Terraform-managed option set. It is Optional and
			// NOT Computed: it holds ONLY the option keys the user declares in
			// the configuration (plus, after Create, the values the backend
			// echoed back for those same managed names). The contract is
			// TWO-STATE and does NOT depend on the legacy SDK distinguishing
			// an omitted block from an explicit empty map (that distinction is
			// unreliable in SDK v1.17.2):
			//   - options = {k: v}: Terraform manages these keys. An Update
			//     derives updates from the new managed set and removals from
			//     oldManaged - newManaged.
			//   - options = {} OR options omitted: both mean the managed set
			//     is expected to be empty. When the SDK diff reports a change
			//     (the reliable path for options = {}), AlterCatalog removals
			//     delete every previously-managed key. removals are derived
			//     ONLY from oldManaged - newManaged, never from remote_options
			//     or server-injected defaults.
			// The SDK fires HasChange reliably for options = {} (and for a
			// populated block); the two-state contract is the documented
			// intent for an omitted block as well — the previous "omit = stop
			// managing, remote keys preserved" three-state claim is removed
			// because the SDK cannot reliably distinguish omit from {} after a
			// refresh. Read reconciles managed keys against the API value and
			// drops a managed key the API no longer returns so drift surfaces.
			// A fresh import cannot recover which remote options the user
			// originally configured, so options is empty after import; the full
			// remote picture (including server defaults and unmanaged keys) is
			// available in remote_options.
			"options": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// remote_options is the read-only, full mirror of the GetCatalog
			// options map, including server-injected defaults and keys the user
			// never declared in options. It is Computed and never drives
			// Update; it exists so a fresh import (where options is empty) does
			// not lose the remote picture, and so users can inspect the full
			// remote options including defaults they did not configure. Because
			// options only ever contains managed keys, AlterCatalog removals
			// are computed from oldManaged - newManaged and can never reach
			// into remote_options to delete a server default or an unmanaged
			// key.
			"remote_options": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_shared": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"share_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudDlfNextCatalogCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dlfNextServiceV2 := DlfNextServiceV2{client}
	name := d.Get("name").(string)
	inputType := d.Get("type").(string)

	// Single deadline established at the START of Create, before the preflight
	// Get, so the preflight, the POST retry, the post-error reconciliation Get
	// and the status waiter all share one budget and the total Create cannot
	// exceed the configured Create timeout. Each phase calls remainingUntil(deadline)
	// to take the time it has left instead of re-asking for the full timeout.
	//
	// Guarantee boundary: the deadline bounds the RETRY/POLL budget for each
	// phase. resource.Retry stops starting new attempts once its timeout is
	// reached, but an HTTP call already in flight is not interrupted; each
	// RoaGet has its own client-level HTTP timeout, so a single trailing call
	// may slightly exceed the deadline. This is the practical bound, not a
	// hard preempt of in-flight I/O.
	createTimeout := d.Timeout(schema.TimeoutCreate)
	deadline := time.Now().Add(createTimeout)

	// Pre-create existence check. CreateCatalog carries no idempotency token
	// and the catalog name is unique within the region, so adopting a catalog
	// that already exists would silently hijack a resource the user may not
	// own. Refuse to create when the name is already taken and direct the user
	// to `terraform import` to take over management. A non-NotFound error from
	// the pre-check (auth/service/network) is a real failure: do not proceed.
	// The decision is a pure function so the pre-create policy is unit-testable
	// without a live backend. The preflight Get uses the deadline's remaining
	// time via DescribeDlfNextCatalogWithRetryTimeout (not a hardcoded minute),
	// so it stays within the shared Create budget.
	_, preErr := dlfNextServiceV2.DescribeDlfNextCatalogWithRetryTimeout(name, remainingUntil(deadline))
	if proceed, pErr := decideDlfNextCatalogPreCreate(preErr); !proceed {
		if pErr != nil {
			return WrapErrorf(pErr, "creating DlfNext Catalog %s", name)
		}
		return nil
	}
	// NotFound -> the name is confirmed absent; proceed to create.

	action := "/dlf/v1/catalogs"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error

	request = make(map[string]interface{})
	request["name"] = name
	request["type"] = inputType
	if v, ok := d.GetOk("options"); ok {
		request["options"] = v
	}
	// is_shared is TypeBool, Optional, no Default: GetOkExists is the
	// SDK-intended pattern for TypeBool (it distinguishes an explicit false
	// from an omitted block, which d.GetOk cannot). This is distinct from the
	// removed GetOkExists(options): options is a TypeMap where the SDK cannot
	// reliably distinguish omit from {} after a refresh, so the two-state
	// contract uses HasChange instead. is_shared is ForceNew and only read at
	// Create (before any refresh), so the omit-vs-false distinction is reliable.
	if v, ok := d.GetOkExists("is_shared"); ok {
		request["isShared"] = v
	}
	if v, ok := d.GetOk("share_id"); ok {
		request["shareId"] = v
	}
	body = request

	// POST retry uses the time remaining in the shared deadline (set before the
	// preflight above), so the preflight + POST together stay within one Create
	// timeout. The backoff is exponential with jitter (review §6 category 1) so
	// multiple ACC runners hitting the same DLF backend do not retry in
	// lockstep and hammer the API. A backend cleanup-quota / environment-
	// capacity error (review §6 category 2) is NOT retried: re-issuing
	// CreateCatalog would add another catalog to the cleanup queue and worsen
	// the quota pressure; it is surfaced as a clear INFRA_BLOCKED error.
	wait := dlfNextCatalogBackoffWait(3*time.Second, 30*time.Second)
	err = resource.Retry(remainingUntil(deadline), func() *resource.RetryError {
		response, err = client.RoaPost("DlfNext", "2025-03-10", action, query, nil, body, true)
		if err != nil {
			// Backend cleanup-quota / env-capacity error: do NOT retry
			// (review §6 category 2). No d.SetId, no state written.
			if isDlfNextCatalogBackendCapacityError(err) {
				return resource.NonRetryableError(fmt.Errorf("DLF backend capacity limit reached for catalog %s in this account/Region; CreateCatalog was not retried to avoid worsening the cleanup quota. Server message: %s. Wait for the backend to reap recently-created catalogs (approximately 48 hours) or use a different account/Region", name, err.Error()))
			}
			alreadyExists := isDlfNextCatalogAlreadyExistsError(err)
			explicitClient := isDlfNextCatalogExplicitClientError(err)
			retryable := NeedRetry(err)
			// Fail-closed ambiguous-create reconciliation. CreateCatalog has no
			// idempotency token and no reliable unique request marker or
			// comparable createdBy, so ownership of a same-name catalog found
			// after an error can NEVER be proven from attributes: type, options,
			// isShared, shareId and createdAt only show attribute equality, not
			// resource ownership, and a concurrent creator sending the same name
			// and config would be wrongly adopted. Therefore a POST that errored
			// ambiguously (lost response / timeout / connection) is NEVER
			// auto-adopted even if GetCatalog then finds a same-name, same-config
			// catalog: no d.SetId, no state written. The decision returns an
			// actionable error directing the user to verify the remote catalog
			// and run `terraform import`. AlreadyExists and explicit 4xx client
			// errors fail outright; a retryable error with the catalog still
			// absent retries; a found catalog or a failed post-error Get fails.
			var getObs dlfNextCatalogGetObservation
			if !alreadyExists && !explicitClient {
				if _, gErr := dlfNextServiceV2.DescribeDlfNextCatalogWithRetryTimeout(name, remainingUntil(deadline)); gErr == nil {
					getObs = dlfGetFound
				} else if NotFoundError(gErr) {
					getObs = dlfGetAbsent
				} else {
					getObs = dlfGetOtherError
				}
			}
			decision, dErr := decideDlfNextCatalogCreateRecovery(alreadyExists, explicitClient, retryable, true, getObs)
			switch decision {
			case dlfCreateRetry:
				wait()
				return resource.RetryableError(err)
			case dlfCreateFail:
				if dErr != nil {
					return resource.NonRetryableError(dErr)
				}
				return resource.NonRetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dlf_next_catalog", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(name)

	// CreateCatalog is asynchronous: it returns 200 with an empty body before
	// the catalog reaches RUNNING. Poll GetCatalog until RUNNING. A brief
	// NotFound immediately after create is treated as Pending (the catalog is
	// not yet visible to Get due to eventual consistency) rather than an
	// unexpected state; INITIALIZE_FAILED is a terminal failure. This waiter's
	// NotFound handling is distinct from the normal Resource Read, which still
	// clears the TF id on NotFound. The waiter uses the time left until the
	// shared deadline so the total Create stays within one timeout.
	stateConf := &resource.StateChangeConf{
		Pending: []string{"NOT_FOUND", "NEW", "INITIALIZING"},
		Target:  []string{"RUNNING"},
		Timeout: remainingUntil(deadline),
		// The waiter uses the deadline-aware Describe via the injectable
		// WithApi variant, so each poll's internal retry is bounded by the
		// shared deadline instead of a hardcoded minute that could nest inside
		// and exceed the waiter timeout.
		Refresh: dlfNextServiceV2.DlfNextCatalogStateRefreshFuncWithApi(d.Id(), "status", []string{"INITIALIZE_FAILED"}, func(id string) (map[string]interface{}, error) {
			return dlfNextServiceV2.DescribeDlfNextCatalogWithRetryTimeout(id, remainingUntil(deadline))
		}),
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, "waiting for DlfNext Catalog %s to RUNNING", d.Id())
	}

	return resourceAliCloudDlfNextCatalogRead(d, meta)
}

func resourceAliCloudDlfNextCatalogRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dlfNextServiceV2 := DlfNextServiceV2{client}

	objectRaw, err := dlfNextServiceV2.DescribeDlfNextCatalog(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dlf_next_catalog DescribeDlfNextCatalog Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", objectRaw["name"])
	d.Set("type", objectRaw["type"])
	d.Set("is_shared", objectRaw["isShared"])
	d.Set("share_id", objectRaw["shareId"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("status", objectRaw["status"])
	d.Set("owner", objectRaw["owner"])
	d.Set("created_at", objectRaw["createdAt"])
	d.Set("created_by", objectRaw["createdBy"])
	d.Set("updated_at", objectRaw["updatedAt"])
	d.Set("updated_by", objectRaw["updatedBy"])

	// Dual-field options model (two-state contract, no GetOkExists). DLF does
	// not flag which GetCatalog options are user-set versus backend-injected
	// defaults, so the state is split into two fields with distinct ownership:
	//
	//  - remote_options (Computed): always mirrors the FULL GetCatalog
	//    options map, including server-injected defaults and keys the user
	//    never declared. It is read-only and never drives Update. This is
	//    the no-data-loss picture: a fresh import that cannot recover the
	//    user's original config intent still exposes the complete remote
	//    options here.
	//
	//  - options (Optional, non-Computed): holds ONLY the keys Terraform
	//    manages. On refresh, each currently-managed key is reconciled to
	//    the API value; a managed key the API no longer returns is dropped
	//    from state so a server-side deletion or an ignored update surfaces
	//    as drift in the next plan instead of being silently preserved. New
	//    remote keys (server defaults, unmanaged additions) are written to
	//    remote_options only and NEVER promoted into options, so AlterCatalog
	//    removals — derived from oldManaged - newManaged — can never reach
	//    into the full remote set and delete a server default.
	//
	// Two-state contract: options = {k: v} means Terraform manages those keys;
	// options = {} OR options omitted both mean the managed set is expected to
	// be empty and previously-managed keys are removed via AlterCatalog
	// removals (oldManaged - newManaged, newManaged empty). This does NOT rely
	// on the legacy SDK distinguishing an omitted block from an explicit empty
	// map — buildDlfNextCatalogOptionChanges treats an empty newManaged the
	// same regardless of origin, and the Update path gates on HasChange which
	// fires reliably for options = {} (and for a populated block). The previous
	// three-state "omit = stop managing / preserve state" claim is removed
	// because the SDK cannot reliably distinguish omit from {} after a refresh.
	apiOptions := make(map[string]interface{})
	if raw, ok := objectRaw["options"].(map[string]interface{}); ok {
		apiOptions = raw
	}
	d.Set("remote_options", reconcileDlfNextCatalogRemoteOptions(apiOptions))
	d.Set("options", reconcileDlfNextCatalogManagedOptions(d.Get("options").(map[string]interface{}), apiOptions))

	return nil
}

func resourceAliCloudDlfNextCatalogUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	catalog := d.Id()
	action := fmt.Sprintf("/dlf/v1/catalogs/%s", catalog)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error

	if d.HasChange("options") {
		// Two-state contract, no GetOkExists guard. options is Optional and
		// non-Computed; the SDK diff encodes user intent. options = {k: v}
		// drives updates/removals from oldManaged - newManaged. options = {}
		// makes newManaged empty so buildDlfNextCatalogOptionChanges emits
		// every previously-managed key as a removal — AlterCatalog deletes
		// ONLY managed keys, never server defaults or unmanaged keys (those
		// live in remote_options, which this branch never reads for
		// removals). An omitted options block has the same documented intent
		// as options = {} (managed set empty); the SDK fires HasChange
		// reliably for options = {} (and for a populated block), and when it
		// fires the removals are computed correctly from oldManaged -
		// newManaged. This replaces the deprecated/undefined GetOkExists guard
		// whose merged state+config+diff+set read could not distinguish "config
		// omits options" from "state carries remote values" after an import.
		oldRaw, newRaw := d.GetChange("options")
		oldOptions, _ := oldRaw.(map[string]interface{})
		newOptions, _ := newRaw.(map[string]interface{})

		updates, removals := buildDlfNextCatalogOptionChanges(oldOptions, newOptions)
		if len(updates) > 0 || len(removals) > 0 {
			body["updates"] = updates
			if len(removals) > 0 {
				body["removals"] = removals
			}

			request = body
			wait := dlfNextCatalogBackoffWait(3*time.Second, 30*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RoaPost("DlfNext", "2025-03-10", action, query, nil, body, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request)

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
	}

	return resourceAliCloudDlfNextCatalogRead(d, meta)
}

func resourceAliCloudDlfNextCatalogDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dlfNextServiceV2 := DlfNextServiceV2{client}
	catalog := d.Id()
	action := fmt.Sprintf("/dlf/v1/catalogs/%s", catalog)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error

	request = make(map[string]interface{})

	// Single deadline shared by the DropCatalog retry and the WaitForGone poll
	// so a Delete does not consume up to 2x the configured timeout. The
	// WaitForGone poll uses the time remaining until the deadline.
	deleteTimeout := d.Timeout(schema.TimeoutDelete)
	deadline := time.Now().Add(deleteTimeout)

	wait := dlfNextCatalogBackoffWait(3*time.Second, 30*time.Second)
	err = resource.Retry(remainingUntil(deadline), func() *resource.RetryError {
		response, err = client.RoaDelete("DlfNext", "2025-03-10", action, query, nil, nil, true)
		if err != nil {
			// Catalog creation is asynchronous; DropCatalog returns 403
			// "... is being created, can not be dropped" while the catalog is
			// still in NEW/INITIALIZING status. NeedRetry does not cover this
			// (it is a 403), so retry explicitly until the catalog is ready.
			// The OpenAPI does not document a stable error code for this case,
			// so the English substring match below is the only available
			// signal and is treated as a fallback.
			if NeedRetry(err) || isDlfNextCatalogBeingCreatedError(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	// DropCatalog is asynchronous: it returns 200 while the catalog is still
	// in DELETING/DELETED status before finally returning 404. Poll GetCatalog
	// until the catalog is gone so a subsequent apply does not observe a stale
	// catalog (e.g. a re-create hitting a name-in-use error). The poll uses
	// the deadline-aware Describe via the injectable WithApi variant, so each
	// poll's internal retry is bounded by the shared deadline instead of a
	// hardcoded minute, and the total Delete stays within one timeout.
	if goneErr := dlfNextServiceV2.WaitForDlfNextCatalogGoneWithApi(d.Id(), remainingUntil(deadline), func(id string) (map[string]interface{}, error) {
		return dlfNextServiceV2.DescribeDlfNextCatalogWithRetryTimeout(id, remainingUntil(deadline))
	}); goneErr != nil {
		return WrapErrorf(goneErr, "waiting for DlfNext Catalog %s to be deleted", d.Id())
	}

	return nil
}

// remainingUntil returns the time left until deadline, clamped at 0 so a passed
// deadline yields a zero (non-negative) duration for resource.Retry and the
// deadline-aware Describe. Create and Delete establish a single shared
// deadline at the start of the operation and every phase (preflight, POST
// retry, post-error reconciliation Get, state waiter, DropCatalog retry, delete
// waiter) calls this instead of re-asking for the full timeout, so the total
// operation stays within one configured timeout. Pure function for unit
// testing the timeout policy.
func remainingUntil(deadline time.Time) time.Duration {
	r := time.Until(deadline)
	if r < 0 {
		return 0
	}
	return r
}

// dlfNextCatalogGetObservation summarises the observable GetCatalog(name)
// result after a CreateCatalog POST that returned an error.
type dlfNextCatalogGetObservation int

const (
	dlfGetAbsent     dlfNextCatalogGetObservation = iota // GetCatalog returned NotFound
	dlfGetFound                                          // catalog exists; ownership cannot be proven from attributes -> fail
	dlfGetOtherError                                     // GetCatalog returned a non-NotFound error
)

// dlfNextCatalogCreateDecision is the ambiguous-create recovery verdict.
type dlfNextCatalogCreateDecision int

const (
	dlfCreateRetry dlfNextCatalogCreateDecision = iota // the create did not land; safe to retry the POST
	dlfCreateFail                                      // surface the error; do not adopt
)

// decideDlfNextCatalogPreCreate decides whether to proceed to CreateCatalog
// after the pre-create Get(name). A nil error (catalog already exists) means
// the name is taken: refuse and direct the user to import. A NotFound error
// means the name is free: proceed. Any other error (auth/service/network) is a
// real failure: do not proceed. Pure function for unit testing.
func decideDlfNextCatalogPreCreate(preGetErr error) (proceed bool, err error) {
	if preGetErr == nil {
		return false, fmt.Errorf("a DlfNext Catalog with this name already exists; to manage an existing catalog, run `terraform import alicloud_dlf_next_catalog.<name> <name>` instead of creating it")
	}
	if NotFoundError(preGetErr) {
		return true, nil
	}
	return false, preGetErr
}

// decideDlfNextCatalogCreateRecovery encodes the fail-closed ambiguous-create
// policy as a pure function of observable inputs so it can be unit-tested
// without a live backend. postErrAlreadyExists / postErrExplicitClientError /
// postErrRetryable classify the POST error; preExistedConfirmedAbsent records
// the pre-create Get result (true only when Get proved the name was absent);
// getObs is the Get result observed after the POST error.
//
// Policy: never take over a catalog just because a same-name catalog exists
// after an error. CreateCatalog has no idempotency token and no reliable
// unique request marker or comparable createdBy, so ownership can NEVER be
// proven from attributes — type, options, isShared, shareId and createdAt only
// show attribute equality, not resource ownership, and a concurrent creator
// sending the same name and config would be wrongly adopted. Therefore:
//   - AlreadyExists or explicit client (validation/permission, 400/401/403)
//     error: fail outright. The error proves the create could not have
//     succeeded, so a same-name catalog found afterwards must be someone
//     else's (or a pre-create race) — do not adopt.
//   - Retryable/ambiguous error with pre-create-proven absence:
//   - catalog now exists (same name, even same config): fail with an
//     actionable error directing the user to verify the remote catalog and
//     run `terraform import`. No d.SetId, no state written.
//   - catalog absent: retry (the create did not land; safe to re-POST).
//   - post-error Get itself errored (non-NotFound): fail (cannot disambiguate).
//   - Any error where pre-create did NOT prove absence: fail (the pre-create
//     guard should have stopped us; never adopt).
func decideDlfNextCatalogCreateRecovery(postErrAlreadyExists, postErrExplicitClientError, postErrRetryable, preExistedConfirmedAbsent bool, getObs dlfNextCatalogGetObservation) (dlfNextCatalogCreateDecision, error) {
	if postErrAlreadyExists {
		return dlfCreateFail, fmt.Errorf("CreateCatalog returned an already-exists error; a catalog with this name already exists, use `terraform import alicloud_dlf_next_catalog.<name> <name>` to manage it instead")
	}
	if postErrExplicitClientError {
		return dlfCreateFail, fmt.Errorf("CreateCatalog returned a validation or permission error; the catalog was not created and an existing same-name catalog must not be adopted")
	}
	if !preExistedConfirmedAbsent {
		return dlfCreateFail, fmt.Errorf("pre-create Get did not confirm the catalog was absent; refusing to create to avoid taking over an existing catalog")
	}
	switch getObs {
	case dlfGetFound:
		return dlfCreateFail, fmt.Errorf("CreateCatalog returned an ambiguous error (timeout, connection, or lost response) and a catalog with this name already exists. This run cannot prove it created the catalog (CreateCatalog has no idempotency token and ownership cannot be inferred from attributes), so it will not be adopted and no state was written. Check whether the remote catalog was created by this run; if so, run `terraform import alicloud_dlf_next_catalog.<name> <name>` to take over management")
	case dlfGetAbsent:
		if postErrRetryable {
			return dlfCreateRetry, nil
		}
		return dlfCreateFail, fmt.Errorf("CreateCatalog returned a non-retryable error and the catalog was not created")
	case dlfGetOtherError:
		return dlfCreateFail, fmt.Errorf("CreateCatalog errored and the post-error GetCatalog also failed; cannot disambiguate whether the catalog was created")
	default:
		return dlfCreateFail, fmt.Errorf("unknown post-error GetCatalog observation")
	}
}

// reconcileDlfNextCatalogRemoteOptions returns the full GetCatalog options map
// stringified, for the read-only remote_options field. It mirrors every remote
// option including server-injected defaults and unmanaged keys, so a fresh
// import (where the managed options field is empty) does not lose the remote
// picture. Pure function for unit testing the drift surface. Never drives
// Update.
func reconcileDlfNextCatalogRemoteOptions(apiOptions map[string]interface{}) map[string]interface{} {
	surfaced := make(map[string]interface{})
	for k, v := range apiOptions {
		surfaced[k] = fmt.Sprint(v)
	}
	return surfaced
}

// reconcileDlfNextCatalogManagedOptions reconciles the currently-managed
// option keys against the API response for the Read of the managed options
// field. For each managed key present in the API response, the API value is
// preferred (so a server-side normalization or an out-of-band change surfaces
// as drift). A managed key the API no longer returns is dropped from state so
// a server-side deletion or an ignored update is visible in the next plan
// instead of being silently preserved. New remote keys (server defaults,
// unmanaged additions) are NOT promoted into the managed set — they remain
// only in remote_options — so AlterCatalog removals, derived from
// oldManaged - newManaged, can never reach into the full remote set. Pure
// function for unit testing the drift behaviour without a live backend.
func reconcileDlfNextCatalogManagedOptions(managedState, apiOptions map[string]interface{}) map[string]interface{} {
	reconciled := make(map[string]interface{})
	for k := range managedState {
		if av, exists := apiOptions[k]; exists {
			reconciled[k] = fmt.Sprint(av)
		}
		// Keys the API no longer returns are intentionally dropped so drift
		// (console deletion / server-side ignore) is visible.
	}
	return reconciled
}

// buildDlfNextCatalogOptionChanges computes the AlterCatalog request payload
// (updates map + removals list) from the old and new managed option sets.
// updates carries every key in newManaged (the desired managed state);
// removals carries every key in oldManaged that is no longer in newManaged.
//
// The boundary guarantee requested in review: removals are derived ONLY from
// oldManaged - newManaged — NEVER from remote_options, the full imported
// remote set, or server-injected defaults. Because the managed options field
// only ever contains keys Terraform manages, an explicit options = {} (which
// makes newManaged empty) removes only previously-managed keys; server
// defaults and unmanaged keys are untouched. Pure function for unit testing
// the AlterCatalog request semantics (verified against the AlterCatalog API
// definition: updates is an object map, removals is an array of key strings)
// without a live backend.
func buildDlfNextCatalogOptionChanges(oldManaged, newManaged map[string]interface{}) (updates map[string]interface{}, removals []interface{}) {
	updates = make(map[string]interface{})
	for k, v := range newManaged {
		updates[k] = v
	}
	removals = make([]interface{}, 0)
	for k := range oldManaged {
		if _, exists := newManaged[k]; !exists {
			removals = append(removals, k)
		}
	}
	return updates, removals
}
