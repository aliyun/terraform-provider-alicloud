// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package expressconnectrouter

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Retryable error codes are declared inline at each write action because the
// compatibility checker parses the IsExpectedErrors argument from the source
// text to detect removed retry codes. The codes are the transient entries of
// the OpenAPI error tables of the write actions (ExpressConnectRouter
// 2023-09-01); terminal errors (quota, parameter format/size, already-exists,
// unsupported feature and not-found) are deliberately absent so they fail fast.
var (
	// vpcAssociationRetryWait is a factory so every lifecycle operation owns its
	// own backoff state.
	vpcAssociationRetryWait = func() func() {
		return helper.IncrementalWait(3*time.Second, 5*time.Second)
	}
	vpcAssociationStateWaitDelay = 5 * time.Second
	// vpcAssociationNow is the clock behind the update budget shared by the
	// sequential writes and their state waits.
	vpcAssociationNow = time.Now
)

func ResourceVpcAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcAssociationCreate,
		Read:   resourceVpcAssociationRead,
		Update: resourceVpcAssociationUpdate,
		Delete: resourceVpcAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
			return validateVpcAssociationPrefixModeRawConfig(diff.GetRawConfig())
		},
		Schema: map[string]*schema.Schema{
			"allowed_prefixes": {
				Type:             schema.TypeList,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: suppressVpcAssociationPrefixReorder,
			},
			"allowed_prefixes_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"association_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"association_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ecr_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_owner_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVpcAssociationCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	if err := validateVpcAssociationPrefixModeRawConfig(d.GetRawConfig()); err != nil {
		return err
	}

	action := "CreateExpressConnectRouterAssociation"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["EcrId"] = d.Get("ecr_id")
	request["VpcId"] = d.Get("vpc_id")

	request["ClientToken"] = helper.BuildClientToken(action)

	request["AssociationRegionId"] = d.Get("association_region_id")
	if v, ok := d.GetOk("vpc_owner_id"); ok {
		request["VpcOwnerId"] = v
	}
	_, prefixesConfigured, err := vpcAssociationRawConfigAttribute(d, "allowed_prefixes")
	if err != nil {
		return err
	}
	mode, modeConfigured, err := vpcAssociationRawConfigAttribute(d, "allowed_prefixes_mode")
	if err != nil {
		return err
	}
	if prefixesConfigured {
		request["AllowedPrefixes"] = d.Get("allowed_prefixes").([]interface{})
		if modeConfigured && strings.TrimSpace(mode.AsString()) != "" {
			request["AllowedPrefixesMode"] = mode.AsString()
		}
	} else if modeConfigured && strings.TrimSpace(mode.AsString()) != "" {
		return errVpcAssociationModeWithoutPrefixes()
	}
	description, descriptionConfigured, err := vpcAssociationRawConfigAttribute(d, "description")
	if err != nil {
		return err
	}
	if descriptionConfigured {
		request["Description"] = description.AsString()
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := vpcAssociationRetryWait()
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Conflict.Lock", "IncorrectStatus.Ecr", "IncorrectStatus.Vpc", "IncorrectStatus.VpcOperating", "IncorrectStatus.Vswitch", "IncorrectStatus.RouteEntry", "OperationDenied.ForceDeleting", "OperationFailed.CallingVpcYaochi", "InternalErrorTemp"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_express_connect_router_vpc_association", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["EcrId"], response["AssociationId"], request["VpcId"]))

	expressConnectRouterServiceV2 := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutCreate), vpcAssociationStateWaitDelay, expressConnectRouterServiceV2.ExpressConnectRouterExpressConnectRouterVpcAssociationStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceVpcAssociationRead(d, meta)
}

func resourceVpcAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	expressConnectRouterServiceV2 := NewService(client)

	objectRaw, err := expressConnectRouterServiceV2.DescribeExpressConnectRouterExpressConnectRouterVpcAssociation(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_express_connect_router_vpc_association DescribeExpressConnectRouterExpressConnectRouterVpcAssociation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("association_region_id", objectRaw["RegionId"])
	d.Set("create_time", objectRaw["GmtCreate"])
	d.Set("description", objectRaw["Description"])
	d.Set("status", objectRaw["Status"])
	d.Set("vpc_owner_id", objectRaw["VpcOwnerId"])
	d.Set("association_id", objectRaw["AssociationId"])
	d.Set("ecr_id", objectRaw["EcrId"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("allowed_prefixes_mode", objectRaw["AllowedPrefixesMode"])

	allowedPrefixes1Raw := make([]interface{}, 0)
	if objectRaw["AllowedPrefixes"] != nil {
		allowedPrefixes1Raw = objectRaw["AllowedPrefixes"].([]interface{})
	}

	d.Set("allowed_prefixes", allowedPrefixes1Raw)

	parts := strings.Split(d.Id(), ":")
	d.Set("ecr_id", parts[0])
	d.Set("association_id", parts[1])
	d.Set("vpc_id", parts[2])

	return nil
}

func resourceVpcAssociationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if err := validateVpcAssociationPrefixModeRawConfig(d.GetRawConfig()); err != nil {
		return err
	}
	parts := strings.Split(d.Id(), ":")
	expressConnectRouterServiceV2 := NewService(client)
	// A single deadline is shared by the sequential modifications and their
	// state waits so a second write cannot restart the full update budget.
	budget := newVpcAssociationUpdateBudget(d.Timeout(schema.TimeoutUpdate))
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	var err error
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)

	// Update description via ModifyExpressConnectRouterAssociation.
	if d.HasChange("description") {
		action := "ModifyExpressConnectRouterAssociation"
		request = make(map[string]interface{})
		query = make(map[string]interface{})
		request["EcrId"] = parts[0]
		request["AssociationId"] = parts[1]
		request["Description"] = d.Get("description")
		request["ClientToken"] = helper.BuildClientToken(action)
		timeout := budget.remaining()
		if timeout <= 0 {
			return vpcAssociationUpdateTimeoutError(d)
		}
		wait := vpcAssociationRetryWait()
		err = retry.Retry(timeout, func() *retry.RetryError {
			response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"Conflict.Lock", "IncorrectStatus.Ecr", "OperationDenied.ForceDeleting"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		timeout = budget.remaining()
		if timeout <= 0 {
			return vpcAssociationUpdateTimeoutError(d)
		}
		stateConf := helper.BuildStateConf([]string{}, []string{"ACTIVE"}, timeout, vpcAssociationStateWaitDelay, expressConnectRouterServiceV2.ExpressConnectRouterExpressConnectRouterVpcAssociationStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	// Update allowed_prefixes / allowed_prefixes_mode via ModifyExpressConnectRouterAssociationAllowedPrefix.
	prefixes, prefixesConfigured, err := vpcAssociationRawConfigAttribute(d, "allowed_prefixes")
	if err != nil {
		return err
	}
	if d.HasChange("allowed_prefixes") || d.HasChange("allowed_prefixes_mode") {
		action := "ModifyExpressConnectRouterAssociationAllowedPrefix"
		request = make(map[string]interface{})
		query = make(map[string]interface{})
		request["EcrId"] = parts[0]
		request["AssociationId"] = parts[1]
		request["ClientToken"] = helper.BuildClientToken(action)
		if prefixesConfigured && prefixes.LengthInt() > 0 {
			request["AllowedPrefixes"] = d.Get("allowed_prefixes").([]interface{})
			// The stored mode accompanies prefix-only updates instead of relying on an
			// undeclared server-side default.
			if mode, modeOk := d.GetOk("allowed_prefixes_mode"); modeOk {
				request["AllowedPrefixesMode"] = mode
			}
		}
		timeout := budget.remaining()
		if timeout <= 0 {
			return vpcAssociationUpdateTimeoutError(d)
		}
		wait := vpcAssociationRetryWait()
		err = retry.Retry(timeout, func() *retry.RetryError {
			response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"Conflict.Lock", "IncorrectStatus.Ecr", "IncorrectStatus.EcrAssociation", "OperationDenied.ForceDeleting"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		timeout = budget.remaining()
		if timeout <= 0 {
			return vpcAssociationUpdateTimeoutError(d)
		}
		stateConf := helper.BuildStateConf([]string{}, []string{"ACTIVE"}, timeout, vpcAssociationStateWaitDelay, expressConnectRouterServiceV2.ExpressConnectRouterExpressConnectRouterVpcAssociationStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return resourceVpcAssociationRead(d, meta)
}

func resourceVpcAssociationDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteExpressConnectRouterAssociation"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["EcrId"] = parts[0]
	request["AssociationId"] = parts[1]

	request["ClientToken"] = helper.BuildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := vpcAssociationRetryWait()
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Conflict.Lock"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ResourceNotFound.AssociationId"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	expressConnectRouterServiceV2 := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), vpcAssociationStateWaitDelay, expressConnectRouterServiceV2.ExpressConnectRouterExpressConnectRouterVpcAssociationStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}

func errVpcAssociationModeWithoutPrefixes() error {
	return fmt.Errorf("`allowed_prefixes_mode` requires a non-empty `allowed_prefixes` list; configure at least one allowed prefix or remove `allowed_prefixes_mode`")
}

func errVpcAssociationEmptyPrefixes() error {
	return fmt.Errorf("`allowed_prefixes` must contain at least one prefix when configured; remove the argument instead of configuring an explicit empty list")
}

func vpcAssociationRawConfigAttribute(d *schema.ResourceData, attribute string) (cty.Value, bool, error) {
	value, diagnostics := d.GetRawConfigAt(cty.GetAttrPath(attribute))
	if diagnostics.HasError() {
		return cty.NilVal, false, errs.WrapError(fmt.Errorf("read raw configuration for `%s`: %s", attribute, diagnostics[0].Summary))
	}
	if !value.IsKnown() || value.IsNull() {
		return value, false, nil
	}
	return value, true, nil
}

// Empty slices have no RPC form representation; unknown values are rechecked at apply time.
func validateVpcAssociationPrefixModeRawConfig(raw cty.Value) error {
	if raw.IsNull() || !raw.IsKnown() {
		return nil
	}
	prefixes := raw.GetAttr("allowed_prefixes")
	if prefixes.IsKnown() && !prefixes.IsNull() && prefixes.LengthInt() == 0 {
		return errVpcAssociationEmptyPrefixes()
	}
	mode := raw.GetAttr("allowed_prefixes_mode")
	if mode.IsNull() || !mode.IsKnown() || strings.TrimSpace(mode.AsString()) == "" {
		return nil
	}
	if !prefixes.IsKnown() {
		return nil
	}
	if prefixes.IsNull() {
		return errVpcAssociationModeWithoutPrefixes()
	}
	return nil
}

// suppressVpcAssociationPrefixReorder manages prefix members without ordering them; membership and multiplicity changes remain diffs.
func suppressVpcAssociationPrefixReorder(_, _, _ string, d *schema.ResourceData) bool {
	if d == nil {
		return false
	}
	rawPrefixes, configured, err := vpcAssociationRawConfigAttribute(d, "allowed_prefixes")
	if err != nil || !configured || rawPrefixes.LengthInt() == 0 {
		return false
	}
	oldRaw, newRaw := d.GetChange("allowed_prefixes")
	oldPrefixes, ok := vpcAssociationPrefixStrings(oldRaw)
	if !ok {
		return false
	}
	newPrefixes, ok := vpcAssociationPrefixStrings(newRaw)
	if !ok {
		return false
	}
	return equalVpcAssociationPrefixSet(oldPrefixes, newPrefixes)
}

// vpcAssociationPrefixStrings accepts a fully known string list; a list that
// still carries unknown or non-string members is left for the apply-time path.
func vpcAssociationPrefixStrings(raw interface{}) ([]string, bool) {
	entries, ok := raw.([]interface{})
	if !ok || len(entries) == 0 {
		return nil, false
	}
	values := make([]string, 0, len(entries))
	for _, entry := range entries {
		value, ok := entry.(string)
		if !ok {
			return nil, false
		}
		values = append(values, value)
	}
	return values, true
}

// equalVpcAssociationPrefixSet compares the members as a multiset so a repeated
// prefix still counts as a real configuration difference.
func equalVpcAssociationPrefixSet(oldPrefixes, newPrefixes []string) bool {
	if len(oldPrefixes) != len(newPrefixes) {
		return false
	}
	oldSorted := append([]string(nil), oldPrefixes...)
	newSorted := append([]string(nil), newPrefixes...)
	sort.Strings(oldSorted)
	sort.Strings(newSorted)
	for i := range oldSorted {
		if oldSorted[i] != newSorted[i] {
			return false
		}
	}
	return true
}

type vpcAssociationUpdateBudget struct {
	deadline time.Time
}

func newVpcAssociationUpdateBudget(timeout time.Duration) *vpcAssociationUpdateBudget {
	return &vpcAssociationUpdateBudget{deadline: vpcAssociationNow().Add(timeout)}
}

func (b *vpcAssociationUpdateBudget) remaining() time.Duration {
	return remainingVpcAssociationTimeout(b.deadline, vpcAssociationNow())
}

// remainingVpcAssociationTimeout is split from the clock read so the shared
// budget regression can assert that later steps inherit the remainder instead
// of restarting the full update timeout.
func remainingVpcAssociationTimeout(deadline, now time.Time) time.Duration {
	remaining := deadline.Sub(now)
	if remaining < 0 {
		return 0
	}
	return remaining
}

func vpcAssociationUpdateTimeoutError(d *schema.ResourceData) error {
	return errs.WrapError(errs.GetTimeErrorFromString(fmt.Sprintf(errs.DefaultTimeoutMsg, d.Id(), "update", "the shared update deadline was exhausted by the preceding modifications")))
}
