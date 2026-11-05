package dlfnext

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceCatalog() *schema.Resource {
	return &schema.Resource{
		Create: resourceCatalogCreate,
		Read:   resourceCatalogRead,
		Update: resourceCatalogUpdate,
		Delete: resourceCatalogDelete,
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
				ValidateFunc: validators.StringInSlice([]string{"PAIMON", "ICEBERG"}, false),
			},
			// options is the Terraform-managed option set. It is Optional and
			// NOT Computed: it holds only the option keys the user declares in
			// the configuration. Both options = {} and an omitted options block
			// mean the managed set is empty; an Update derives updates from the
			// new managed set and removals from oldManaged - newManaged.
			// remote_options holds the full remote picture.
			"options": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// remote_options is the read-only, full mirror of the GetCatalog
			// options map, including server-injected defaults and keys the user
			// never declared. It is Computed and never drives Update; it exists
			// so a fresh import (where options is empty) does not lose the
			// remote picture. AlterCatalog removals are computed from
			// oldManaged - newManaged and can never reach into remote_options
			// to delete a server default or an unmanaged key.
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCatalogCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := NewService(client)
	name := d.Get("name").(string)
	inputType := d.Get("type").(string)

	// Single deadline established at the START of Create, before the preflight
	// Get, so the preflight, the POST retry, the post-error reconciliation Get
	// and the status waiter all share one budget and the total Create cannot
	// exceed the configured Create timeout.
	deadline := time.Now().Add(d.Timeout(schema.TimeoutCreate))

	// Pre-create existence check. CreateCatalog carries no idempotency token
	// and the catalog name is unique within the region, so adopting a catalog
	// that already exists would silently hijack a resource the user may not
	// own. Refuse to create when the name is already taken and direct the user
	// to `terraform import`. A non-NotFound error from the pre-check
	// (auth/service/network) is a real failure: do not proceed.
	_, preErr := service.DescribeCatalogWithRetryTimeout(name, remainingUntil(deadline))
	if proceed, pErr := decideCatalogPreCreate(preErr); !proceed {
		if pErr != nil {
			return errs.WrapErrorf(pErr, "creating DlfNext Catalog %s", name)
		}
		return nil
	}

	action := "/dlf/v1/catalogs"
	var response map[string]interface{}
	query := make(map[string]*string)
	request := make(map[string]interface{})
	var err error

	request["name"] = name
	request["type"] = inputType
	if v, ok := d.GetOk("options"); ok {
		request["options"] = v
	}
	// is_shared is TypeBool, Optional, no Default: GetOkExists is the
	// SDK-intended pattern for TypeBool (it distinguishes an explicit false
	// from an omitted block, which d.GetOk cannot). is_shared is ForceNew and
	// only read at Create (before any refresh), so the omit-vs-false
	// distinction is reliable.
	if v, ok := d.GetOkExists("is_shared"); ok {
		request["isShared"] = v
	}
	if v, ok := d.GetOk("share_id"); ok {
		request["shareId"] = v
	}

	// The POST retry shares the Create deadline and uses exponential backoff
	// with jitter. A backend cleanup-quota / environment-capacity error is NOT
	// retried: re-issuing CreateCatalog would add another catalog to the
	// cleanup queue and worsen the quota pressure; it is surfaced clearly.
	wait := catalogBackoffWait(3*time.Second, 30*time.Second)
	err = retry.Retry(remainingUntil(deadline), func() *retry.RetryError {
		response, err = client.RoaPost("DlfNext", "2025-03-10", action, query, nil, request, true)
		if err != nil {
			if isCatalogBackendCapacityError(err) {
				return retry.NonRetryableError(fmt.Errorf("DLF backend capacity limit reached for catalog %s in this account/Region; CreateCatalog was not retried to avoid worsening the cleanup quota. Server message: %s. Wait for the backend to reap recently-created catalogs or use a different account/Region", name, err.Error()))
			}
			alreadyExists := isCatalogAlreadyExistsError(err)
			explicitClient := isCatalogExplicitClientError(err)
			retryable := helper.NeedRetry(err)
			// Fail-closed ambiguous-create reconciliation. CreateCatalog has no
			// idempotency token and no reliable unique request marker, so
			// ownership of a same-name catalog found after an error can never
			// be proven from attributes.
			var getObs catalogGetObservation
			if !alreadyExists && !explicitClient {
				if _, gErr := service.DescribeCatalogWithRetryTimeout(name, remainingUntil(deadline)); gErr == nil {
					getObs = catalogGetFound
				} else if errs.NotFoundError(gErr) {
					getObs = catalogGetAbsent
				} else {
					getObs = catalogGetOtherError
				}
			}
			decision, dErr := decideCatalogCreateRecovery(alreadyExists, explicitClient, retryable, true, getObs)
			switch decision {
			case catalogCreateRetry:
				wait()
				return retry.RetryableError(err)
			case catalogCreateFail:
				if dErr != nil {
					return retry.NonRetryableError(dErr)
				}
				return retry.NonRetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_dlf_next_catalog", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(name)

	// CreateCatalog is asynchronous: it returns 200 before the catalog reaches
	// RUNNING. Poll GetCatalog until RUNNING. A brief NotFound immediately
	// after create is treated as Pending (eventual consistency), and
	// INITIALIZE_FAILED is a terminal failure. This waiter's NotFound handling
	// is distinct from the normal Read, which clears the TF id on NotFound.
	stateConf := &retry.StateChangeConf{
		Pending: []string{"NOT_FOUND", "NEW", "INITIALIZING"},
		Target:  []string{"RUNNING"},
		Timeout: remainingUntil(deadline),
		Refresh: service.CatalogStateRefreshFuncWithApi(d.Id(), "status", []string{"INITIALIZE_FAILED"}, func(id string) (map[string]interface{}, error) {
			return service.DescribeCatalogWithRetryTimeout(id, remainingUntil(deadline))
		}),
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, "waiting for DlfNext Catalog %s to RUNNING", d.Id())
	}

	return resourceCatalogRead(d, meta)
}

func resourceCatalogRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := NewService(client)

	objectRaw, err := service.DescribeCatalog(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dlf_next_catalog DescribeCatalog Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("name", objectRaw["name"])
	d.Set("type", objectRaw["type"])
	d.Set("is_shared", objectRaw["isShared"])
	d.Set("share_id", objectRaw["shareId"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("status", objectRaw["status"])
	d.Set("owner", objectRaw["owner"])
	d.Set("created_at", catalogTimestampString(objectRaw["createdAt"]))
	d.Set("created_by", objectRaw["createdBy"])
	d.Set("updated_at", catalogTimestampString(objectRaw["updatedAt"]))
	d.Set("updated_by", objectRaw["updatedBy"])

	// DLF does not flag which GetCatalog options are user-set versus
	// backend-injected defaults, so the state is split into two fields with
	// distinct ownership: remote_options mirrors the full remote set, while
	// options holds only the keys Terraform manages. A managed key the API no
	// longer returns is dropped from options so server-side drift surfaces in
	// the next plan; new remote keys never get promoted into options, so
	// AlterCatalog removals cannot delete a server default.
	apiOptions := make(map[string]interface{})
	if raw, ok := objectRaw["options"].(map[string]interface{}); ok {
		apiOptions = raw
	}
	d.Set("remote_options", reconcileRemoteOptions(apiOptions))
	d.Set("options", reconcileManagedOptions(d.Get("options").(map[string]interface{}), apiOptions))

	return nil
}

func resourceCatalogUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	catalog := d.Id()
	action := fmt.Sprintf("/dlf/v1/catalogs/%s", catalog)
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var response map[string]interface{}
	var err error

	if d.HasChange("options") {
		// options is Optional and non-Computed; the SDK diff encodes user
		// intent. options = {} makes newManaged empty so every
		// previously-managed key is emitted as a removal; AlterCatalog deletes
		// only managed keys, never server defaults or unmanaged keys.
		oldRaw, newRaw := d.GetChange("options")
		oldOptions, _ := oldRaw.(map[string]interface{})
		newOptions, _ := newRaw.(map[string]interface{})

		updates, removals := buildOptionChanges(oldOptions, newOptions)
		if len(updates) > 0 || len(removals) > 0 {
			body["updates"] = updates
			if len(removals) > 0 {
				body["removals"] = removals
			}

			wait := catalogBackoffWait(3*time.Second, 30*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RoaPost("DlfNext", "2025-03-10", action, query, nil, body, true)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, body)

			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
	}

	return resourceCatalogRead(d, meta)
}

func resourceCatalogDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := NewService(client)
	catalog := d.Id()
	action := fmt.Sprintf("/dlf/v1/catalogs/%s", catalog)
	query := make(map[string]*string)
	var response map[string]interface{}
	var err error

	// Single deadline shared by the DropCatalog retry and the WaitForGone poll
	// so a Delete does not consume up to 2x the configured timeout.
	deadline := time.Now().Add(d.Timeout(schema.TimeoutDelete))

	wait := catalogBackoffWait(3*time.Second, 30*time.Second)
	err = retry.Retry(remainingUntil(deadline), func() *retry.RetryError {
		response, err = client.RoaDelete("DlfNext", "2025-03-10", action, query, nil, nil, true)
		if err != nil {
			// DropCatalog returns 403 "... is being created, can not be
			// dropped" while the catalog is still in NEW/INITIALIZING.
			// NeedRetry does not cover this (it is a 403), so retry it until
			// the catalog is ready.
			if helper.NeedRetry(err) || isCatalogBeingCreatedError(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, nil)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	// DropCatalog is asynchronous: it returns 200 while the catalog is still
	// in DELETING/DELETED before finally returning 404. Poll GetCatalog until
	// the catalog is gone so a subsequent apply does not observe a stale
	// catalog (e.g. a re-create hitting a name-in-use error).
	if goneErr := service.WaitForCatalogGoneWithApi(d.Id(), remainingUntil(deadline), func(id string) (map[string]interface{}, error) {
		return service.DescribeCatalogWithRetryTimeout(id, remainingUntil(deadline))
	}); goneErr != nil {
		return errs.WrapErrorf(goneErr, "waiting for DlfNext Catalog %s to be deleted", d.Id())
	}

	return nil
}

// remainingUntil returns the time left until deadline, clamped at 0 so a passed
// deadline yields a zero (non-negative) duration for retry.Retry and the
// deadline-aware Describe.
func remainingUntil(deadline time.Time) time.Duration {
	r := time.Until(deadline)
	if r < 0 {
		return 0
	}
	return r
}

// catalogTimestampString converts a GetCatalog millisecond timestamp into the
// string form the schema exposes. The API returns int64 milliseconds, which
// overflow a 32-bit int; TypeString avoids that and keeps the full value.
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

// catalogGetObservation summarises the observable GetCatalog(name) result
// after a CreateCatalog POST that returned an error.
type catalogGetObservation int

const (
	catalogGetAbsent     catalogGetObservation = iota // GetCatalog returned NotFound
	catalogGetFound                                   // catalog exists; ownership cannot be proven from attributes -> fail
	catalogGetOtherError                              // GetCatalog returned a non-NotFound error
)

// catalogCreateDecision is the ambiguous-create recovery verdict.
type catalogCreateDecision int

const (
	catalogCreateRetry catalogCreateDecision = iota // the create did not land; safe to retry the POST
	catalogCreateFail                               // surface the error; do not adopt
)

// decideCatalogPreCreate decides whether to proceed to CreateCatalog after the
// pre-create Get(name). A nil error (catalog already exists) means the name is
// taken: refuse and direct the user to import. A NotFound error means the name
// is free: proceed. Any other error is a real failure: do not proceed.
func decideCatalogPreCreate(preGetErr error) (proceed bool, err error) {
	if preGetErr == nil {
		return false, fmt.Errorf("a DlfNext Catalog with this name already exists; to manage an existing catalog, run `terraform import alicloud_dlf_next_catalog.<name> <name>` instead of creating it")
	}
	if errs.NotFoundError(preGetErr) {
		return true, nil
	}
	return false, preGetErr
}

// decideCatalogCreateRecovery encodes the fail-closed ambiguous-create policy
// as a pure function of observable inputs so it can be unit-tested without a
// live backend. Policy: never take over a catalog just because a same-name
// catalog exists after an error. CreateCatalog has no idempotency token, so
// ownership can never be proven from attributes.
//   - AlreadyExists or explicit client (validation/permission, 400/401/403)
//     error: fail outright.
//   - Retryable/ambiguous error with pre-create-proven absence:
//   - catalog now exists: fail with an actionable import error; no state.
//   - catalog absent: retry.
//   - post-error Get errored: fail.
//   - Any error where pre-create did not prove absence: fail.
func decideCatalogCreateRecovery(postErrAlreadyExists, postErrExplicitClientError, postErrRetryable, preExistedConfirmedAbsent bool, getObs catalogGetObservation) (catalogCreateDecision, error) {
	if postErrAlreadyExists {
		return catalogCreateFail, fmt.Errorf("CreateCatalog returned an already-exists error; a catalog with this name already exists, use `terraform import alicloud_dlf_next_catalog.<name> <name>` to manage it instead")
	}
	if postErrExplicitClientError {
		return catalogCreateFail, fmt.Errorf("CreateCatalog returned a validation or permission error; the catalog was not created and an existing same-name catalog must not be adopted")
	}
	if !preExistedConfirmedAbsent {
		return catalogCreateFail, fmt.Errorf("pre-create Get did not confirm the catalog was absent; refusing to create to avoid taking over an existing catalog")
	}
	switch getObs {
	case catalogGetFound:
		return catalogCreateFail, fmt.Errorf("CreateCatalog returned an ambiguous error (timeout, connection, or lost response) and a catalog with this name already exists. This run cannot prove it created the catalog, so it will not be adopted and no state was written. Check whether the remote catalog was created by this run; if so, run `terraform import alicloud_dlf_next_catalog.<name> <name>` to take over management")
	case catalogGetAbsent:
		if postErrRetryable {
			return catalogCreateRetry, nil
		}
		return catalogCreateFail, fmt.Errorf("CreateCatalog returned a non-retryable error and the catalog was not created")
	case catalogGetOtherError:
		return catalogCreateFail, fmt.Errorf("CreateCatalog errored and the post-error GetCatalog also failed; cannot disambiguate whether the catalog was created")
	default:
		return catalogCreateFail, fmt.Errorf("unknown post-error GetCatalog observation")
	}
}

// reconcileRemoteOptions returns the full GetCatalog options map stringified,
// for the read-only remote_options field. It mirrors every remote option
// including server-injected defaults and unmanaged keys. Never drives Update.
func reconcileRemoteOptions(apiOptions map[string]interface{}) map[string]interface{} {
	surfaced := make(map[string]interface{})
	for k, v := range apiOptions {
		surfaced[k] = fmt.Sprint(v)
	}
	return surfaced
}

// reconcileManagedOptions reconciles the currently-managed option keys against
// the API response. For each managed key present in the API response the API
// value is preferred; a managed key the API no longer returns is dropped so
// drift surfaces. New remote keys are not promoted into the managed set.
func reconcileManagedOptions(managedState, apiOptions map[string]interface{}) map[string]interface{} {
	reconciled := make(map[string]interface{})
	for k := range managedState {
		if av, exists := apiOptions[k]; exists {
			reconciled[k] = fmt.Sprint(av)
		}
	}
	return reconciled
}

// buildOptionChanges computes the AlterCatalog request payload (updates map +
// removals list) from the old and new managed option sets. removals are
// derived ONLY from oldManaged - newManaged, never from remote_options or
// server-injected defaults.
func buildOptionChanges(oldManaged, newManaged map[string]interface{}) (updates map[string]interface{}, removals []interface{}) {
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
