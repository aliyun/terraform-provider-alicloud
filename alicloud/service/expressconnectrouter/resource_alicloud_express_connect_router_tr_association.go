// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package expressconnectrouter

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTRAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTRAssociationCreate,
		Read:   resourceTRAssociationRead,
		Update: resourceTRAssociationUpdate,
		Delete: resourceTRAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
			return validateTRAssociationPrefixModeRawConfig(diff.GetRawConfig())
		},
		Schema: map[string]*schema.Schema{
			"allowed_prefixes": {
				Type:             schema.TypeList,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: suppressTRAssociationPrefixReorder,
			},
			"allowed_prefixes_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringMatch(regexp.MustCompile("^(MatchMode|IncrementalMode)$"), "The prefix routing mode. Valid values: MatchMode, IncrementalMode."),
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
			"cen_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_attachment": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringMatch(regexp.MustCompile(`^[\s\S]{0,128}$`), "The description of the association resource. The length is 0-128 characters."),
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
			"transit_router_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"transit_router_owner_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func errTRAssociationModeWithoutPrefixes() error {
	return fmt.Errorf("`allowed_prefixes_mode` requires a non-empty `allowed_prefixes` list; configure at least one allowed prefix or remove `allowed_prefixes_mode`")
}

func validateTRAssociationPrefixModeRawConfig(raw cty.Value) error {
	if raw == cty.NilVal || raw.IsNull() || !raw.IsKnown() {
		return nil
	}
	prefixes := raw.GetAttr("allowed_prefixes")
	mode := raw.GetAttr("allowed_prefixes_mode")
	if mode.IsNull() || !mode.IsKnown() || strings.TrimSpace(mode.AsString()) == "" || !prefixes.IsKnown() {
		return nil
	}
	if prefixes.IsNull() || prefixes.LengthInt() == 0 {
		return errTRAssociationModeWithoutPrefixes()
	}
	return nil
}

func validateTRAssociationPrefixModeApply(d *schema.ResourceData) error {
	raw := d.GetRawConfig()
	if err := validateTRAssociationPrefixModeRawConfig(raw); err != nil {
		return err
	}
	if raw == cty.NilVal || raw.IsNull() || !raw.IsKnown() {
		return nil
	}
	prefixes := raw.GetAttr("allowed_prefixes")
	mode := raw.GetAttr("allowed_prefixes_mode")
	if mode.IsNull() || (mode.IsKnown() && strings.TrimSpace(mode.AsString()) == "") {
		return nil
	}
	if mode.IsKnown() && prefixes.IsKnown() {
		return nil
	}
	resolvedMode, ok := d.GetOk("allowed_prefixes_mode")
	if !ok || strings.TrimSpace(resolvedMode.(string)) == "" {
		return nil
	}
	resolvedPrefixes, ok := d.GetOk("allowed_prefixes")
	if !ok || len(resolvedPrefixes.([]interface{})) == 0 {
		return errTRAssociationModeWithoutPrefixes()
	}
	return nil
}

// Prefix membership and multiplicity are significant; list positions are not.
func suppressTRAssociationPrefixReorder(_, _, _ string, d *schema.ResourceData) bool {
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

func resourceTRAssociationCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	if err := validateTRAssociationPrefixModeApply(d); err != nil {
		return err
	}

	action := "CreateExpressConnectRouterAssociation"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	// One deadline bounds the create call, its retries and the state wait.
	createDeadline := trAssociationNow().Add(d.Timeout(schema.TimeoutCreate))
	ctx, cancel := context.WithDeadline(context.Background(), createDeadline)
	defer cancel()
	request = make(map[string]interface{})
	request["EcrId"] = d.Get("ecr_id")
	request["TransitRouterId"] = d.Get("transit_router_id")

	request["ClientToken"] = helper.BuildClientToken(action)

	request["AssociationRegionId"] = d.Get("association_region_id")
	if v, ok := d.GetOk("allowed_prefixes"); ok {
		allowedPrefixesMaps := v.([]interface{})
		request["AllowedPrefixes"] = allowedPrefixesMaps
	}

	if v, ok := d.GetOk("transit_router_owner_id"); ok {
		request["TransitRouterOwnerId"] = v
	}
	if v, ok := d.GetOk("cen_id"); ok {
		request["CenId"] = v
	}
	// create_attachment is Optional+Computed, so GetOk cannot tell an omitted
	// argument (leave the service default) from an explicit false (must be sent).
	if v, err := helper.GetWriteOnlyValue(d, cty.GetAttrPath("create_attachment"), cty.Bool); err != nil {
		return err
	} else if !v.IsNull() && v.IsKnown() {
		request["CreateAttachment"] = v.True()
	}
	if v, ok := d.GetOk("allowed_prefixes_mode"); ok {
		request["AllowedPrefixesMode"] = v
	}
	if rawConfig := d.GetRawConfig(); rawConfig != cty.NilVal && !rawConfig.IsNull() {
		if v := rawConfig.GetAttr("description"); v.IsKnown() && !v.IsNull() {
			request["Description"] = v.AsString()
		}
	}
	stageTimeout, err := operationTimeoutRemaining(createDeadline, "create")
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_express_connect_router_tr_association", action, errs.AlibabaCloudSdkGoERROR)
	}
	response, err = helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
		resp, callErr := trAssociationRpcPost(ctx, client, "ExpressConnectRouter", "2023-09-01", action, query, request)
		if callErr == nil {
			logs.AddDebug(action, resp, request)
		}
		return resp, callErr
	}, stageTimeout, func(err error) bool {
		// InternalErrorTemp and OperationFailed.CallingVpcYaochi are documented
		// as retryable; keep the other retryable codes and terminal/parameter
		// errors unchanged.
		return trAssociationBusinessCode(err) != "" && errs.IsExpectedErrors(err, []string{"Conflict.Lock", "IncorrectStatus.Ecr", "OperationDenied.ForceDeleting", "InternalErrorTemp", "OperationFailed.CallingVpcYaochi"}) || helper.NeedRetry(err)
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_express_connect_router_tr_association", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["EcrId"], response["AssociationId"], request["TransitRouterId"]))

	stageTimeout, err = operationTimeoutRemaining(createDeadline, "create")
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	stateConf := helper.BuildStateConf([]string{}, []string{"ACTIVE", "INACTIVE"}, stageTimeout, trAssociationStateDelay, trAssociationStateRefresh(ctx, client, d.Id()))
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceTRAssociationReadWithContext(ctx, d, meta)
}

func resourceTRAssociationRead(d *schema.ResourceData, meta interface{}) error {
	return resourceTRAssociationReadWithContext(context.Background(), d, meta)
}

// resourceTRAssociationReadWithContext fills state through the cancellable
// describe path. Create and Update pass their operation context so the final
// refresh shares the same remaining budget and cancellation chain as the writes
// and state waits; the standalone Read and Import paths keep their original
// behaviour with a background context.
func resourceTRAssociationReadWithContext(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	objectRaw, err := describeTRAssociation(ctx, client, d.Id())
	if err != nil {
		if !d.IsNewResource() && trAssociationNotFound(err) {
			log.Printf("[DEBUG] Resource alicloud_express_connect_router_tr_association DescribeExpressConnectRouterExpressConnectRouterTrAssociation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("association_region_id", objectRaw["RegionId"])
	d.Set("cen_id", objectRaw["CenId"])
	d.Set("create_time", objectRaw["GmtCreate"])
	d.Set("status", objectRaw["Status"])
	d.Set("transit_router_owner_id", objectRaw["TransitRouterOwnerId"])
	d.Set("association_id", objectRaw["AssociationId"])
	d.Set("ecr_id", objectRaw["EcrId"])
	d.Set("transit_router_id", objectRaw["TransitRouterId"])

	allowedPrefixes1Raw := make([]interface{}, 0)
	if objectRaw["AllowedPrefixes"] != nil {
		allowedPrefixes1Raw = objectRaw["AllowedPrefixes"].([]interface{})
	}

	d.Set("allowed_prefixes", allowedPrefixes1Raw)
	d.Set("allowed_prefixes_mode", objectRaw["AllowedPrefixesMode"])
	d.Set("description", objectRaw["Description"])
	// create_attachment is a create-only request parameter that DescribeExpressConnectRouterAssociation
	// does not return; keep the configured value in state instead of overwriting it with nil.

	parts := strings.Split(d.Id(), ":")
	d.Set("ecr_id", parts[0])
	d.Set("association_id", parts[1])
	d.Set("transit_router_id", parts[2])

	return nil
}

// trAssociationNow is replaceable in tests so operation budget behaviour can be
// verified without real waits.
var trAssociationNow = time.Now

// trAssociationStateDelay is the initial wait before the first status check.
// Tests shorten it through SetTRAssociationStateDelayForTest.
var trAssociationStateDelay = 5 * time.Second

// operationTimeoutRemaining returns the budget still available before deadline.
// Every stage of one Create, Update or Delete call must share the same deadline:
// restarting d.Timeout(...) per stage would let the whole chain run several
// times longer than the timeout the user configured.
func operationTimeoutRemaining(deadline time.Time, operation string) (time.Duration, error) {
	remaining := deadline.Sub(trAssociationNow())
	if remaining <= 0 {
		return 0, fmt.Errorf("the configured %s timeout was exhausted by earlier steps of the same operation", operation)
	}
	return remaining, nil
}

func trAssociationRpcPost(ctx context.Context, client *connectivity.AliyunClient, product, version, action string, query, body map[string]interface{}) (map[string]interface{}, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return client.RpcPostWithContext(ctx, product, version, action, query, body)
}

func trAssociationSdkError(err error) *tea.SDKError {
	for err != nil {
		if sdkErr, ok := err.(*tea.SDKError); ok {
			return sdkErr
		}
		switch typed := err.(type) {
		case *errs.ComplexError:
			err = typed.Cause
		case interface{ Unwrap() error }:
			err = typed.Unwrap()
		default:
			return nil
		}
	}
	return nil
}

func trAssociationBusinessCode(err error) string {
	sdkErr := trAssociationSdkError(err)
	if sdkErr == nil {
		return ""
	}
	return strings.TrimSpace(tea.StringValue(sdkErr.Code))
}

func trAssociationNotFound(err error) bool {
	if errs.IsProviderNotFound(err) {
		return true
	}
	return trAssociationBusinessCode(err) != "" && errs.NotFoundError(err)
}

func trAssociationMatches(association map[string]interface{}, parts []string) (bool, error) {
	identityFields := []struct {
		name     string
		expected string
	}{
		{"AssociationId", parts[1]},
		{"EcrId", parts[0]},
		{"TransitRouterId", parts[2]},
		{"AssociationNodeType", "TR"},
	}
	for _, field := range identityFields {
		value, ok := association[field.name].(string)
		if !ok || value == "" {
			return false, fmt.Errorf("expected AssociationList item %s to be a nonempty string, got %T", field.name, association[field.name])
		}
		if value != field.expected {
			return false, nil
		}
	}
	return true, nil
}

func trAssociationEmptyBusinessObject(response map[string]interface{}) bool {
	for key := range response {
		if key != "_headers" {
			return false
		}
	}
	return true
}

// describeTRAssociation reads one association through the cancellable RPC path
// instead of the shared Service describe, whose retry.Retry loop and context-free
// client.RpcPost cannot be bounded by the caller's deadline. NotFound
// classification matches the shared path: a missing parent ECR or an empty
// association list both mean the target no longer exists.
func describeTRAssociation(ctx context.Context, client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		return nil, errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "DescribeExpressConnectRouterAssociation"
	request := map[string]interface{}{
		"AssociationId":       parts[1],
		"EcrId":               parts[0],
		"TransitRouterId":     parts[2],
		"AssociationNodeType": "TR",
		"ClientToken":         helper.BuildClientToken(action),
	}
	var response map[string]interface{}
	for {
		var err error
		response, err = helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
			resp, callErr := trAssociationRpcPost(ctx, client, "ExpressConnectRouter", "2023-09-01", action, map[string]interface{}{}, request)
			if callErr == nil {
				logs.AddDebug(action, resp, request)
			}
			return resp, callErr
		}, time.Minute, helper.NeedRetry)
		if err != nil {
			if trAssociationBusinessCode(err) != "" && errs.IsExpectedErrors(err, []string{"ResourceNotFound.EcrId", "ResourceNotFound.AssociationId"}) {
				return nil, errs.WrapErrorf(errs.NotFoundErr("ExpressConnectRouterTrAssociation", id), errs.NotFoundMsg, response)
			}
			return nil, errs.WrapError(err)
		}
		if trAssociationEmptyBusinessObject(response) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("ExpressConnectRouterTrAssociation", id), errs.NotFoundMsg, response)
		}
		items, ok := response["AssociationList"].([]interface{})
		if !ok {
			return nil, errs.WrapErrorf(fmt.Errorf("expected AssociationList to be an array, got %T", response["AssociationList"]), errs.FailedGetAttributeMsg, action, "$.AssociationList", response)
		}
		for _, item := range items {
			association, ok := item.(map[string]interface{})
			if !ok {
				return nil, errs.WrapErrorf(fmt.Errorf("expected AssociationList item to be an object, got %T", item), errs.FailedGetAttributeMsg, action, "$.AssociationList", response)
			}
			matches, err := trAssociationMatches(association, parts)
			if err != nil {
				return nil, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.AssociationList[*]", response)
			}
			if !matches {
				continue
			}
			return association, nil
		}
		rawNextToken, present := response["NextToken"]
		nextToken, ok := rawNextToken.(string)
		if present && !ok {
			return nil, errs.WrapErrorf(fmt.Errorf("expected NextToken to be a string, got %T", rawNextToken), errs.FailedGetAttributeMsg, action, "$.NextToken", response)
		}
		if nextToken == "" {
			return nil, errs.WrapErrorf(errs.NotFoundErr("ExpressConnectRouterTrAssociation", id), errs.NotFoundMsg, response)
		}
		request["NextToken"] = nextToken
	}
}

// trAssociationStateRefresh builds the state refresh for one association using
// the same operation context. It does not use the shared Service describe path,
// whose internal retry loop cannot be cancelled by the operation deadline.
func trAssociationStateRefresh(ctx context.Context, client *connectivity.AliyunClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		association, err := describeTRAssociation(ctx, client, id)
		if err != nil {
			if trAssociationNotFound(err) {
				return nil, "", nil
			}
			return nil, "", err
		}
		return association, fmt.Sprint(association["Status"]), nil
	}
}

func resourceTRAssociationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if err := validateTRAssociationPrefixModeApply(d); err != nil {
		return err
	}
	var request map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "ModifyExpressConnectRouterAssociationAllowedPrefix"
	// One deadline bounds the prefix update, its wait and the later
	// description update of the same Terraform update operation.
	updateDeadline := trAssociationNow().Add(d.Timeout(schema.TimeoutUpdate))
	ctx, cancel := context.WithDeadline(context.Background(), updateDeadline)
	defer cancel()
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = parts[0]
	request["AssociationId"] = parts[1]
	request["ClientToken"] = helper.BuildClientToken(action)
	if d.HasChange("allowed_prefixes") || d.HasChange("allowed_prefixes_mode") {
		update = true
		if v, ok := d.GetOk("allowed_prefixes"); ok {
			allowedPrefixesMaps := v.([]interface{})
			request["AllowedPrefixes"] = allowedPrefixesMaps
			// Clearing prefixes requires omitting the computed mode as well.
			if mode, modeOk := d.GetOk("allowed_prefixes_mode"); modeOk {
				request["AllowedPrefixesMode"] = mode
			}
		}
	}

	if update {
		stageTimeout, err := operationTimeoutRemaining(updateDeadline, "update")
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		_, err = helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
			resp, callErr := trAssociationRpcPost(ctx, client, "ExpressConnectRouter", "2023-09-01", action, query, request)
			if callErr == nil {
				logs.AddDebug(action, resp, request)
			}
			return resp, callErr
		}, stageTimeout, func(err error) bool {
			return trAssociationBusinessCode(err) != "" && errs.IsExpectedErrors(err, []string{"Conflict.Lock", "IncorrectStatus.Ecr", "IncorrectStatus.EcrAssociation", "OperationDenied.ForceDeleting"}) || helper.NeedRetry(err)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		stageTimeout, err = operationTimeoutRemaining(updateDeadline, "update")
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		stateConf := helper.BuildStateConf([]string{}, []string{"ACTIVE", "INACTIVE"}, stageTimeout, trAssociationStateDelay, trAssociationStateRefresh(ctx, client, d.Id()))
		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	// ModifyExpressConnectRouterAssociation: description
	update = false
	action = "ModifyExpressConnectRouterAssociation"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = parts[0]
	request["AssociationId"] = parts[1]
	request["ClientToken"] = helper.BuildClientToken(action)
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if update {
		stageTimeout, err := operationTimeoutRemaining(updateDeadline, "update")
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		_, err = helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
			resp, callErr := trAssociationRpcPost(ctx, client, "ExpressConnectRouter", "2023-09-01", action, query, request)
			if callErr == nil {
				logs.AddDebug(action, resp, request)
			}
			return resp, callErr
		}, stageTimeout, func(err error) bool {
			// ModifyExpressConnectRouterAssociation documents only these transient
			// 400 codes (IncorrectStatus.Ecr, OperationDenied.ForceDeleting,
			// Conflict.Lock); IncorrectStatus.EcrAssociation belongs to the
			// AllowedPrefix action and is not in this action's error table.
			return trAssociationBusinessCode(err) != "" && errs.IsExpectedErrors(err, []string{"Conflict.Lock", "IncorrectStatus.Ecr", "OperationDenied.ForceDeleting"}) || helper.NeedRetry(err)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		stageTimeout, err = operationTimeoutRemaining(updateDeadline, "update")
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		stateConf := helper.BuildStateConf([]string{}, []string{"ACTIVE", "INACTIVE"}, stageTimeout, trAssociationStateDelay, trAssociationStateRefresh(ctx, client, d.Id()))
		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return resourceTRAssociationReadWithContext(ctx, d, meta)
}

func resourceTRAssociationDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 3 {
		return errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", d.Id(), 3, len(parts)))
	}
	action := "DeleteExpressConnectRouterAssociation"
	var request map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["EcrId"] = parts[0]
	request["AssociationId"] = parts[1]

	// One deadline bounds every delete attempt and the final wait of the same
	// Terraform destroy operation.
	deleteDeadline := trAssociationNow().Add(d.Timeout(schema.TimeoutDelete))
	ctx, cancel := context.WithDeadline(context.Background(), deleteDeadline)
	defer cancel()

	// Only an attachment explicitly created by this association may be deleted
	// with it. Imported or independently managed attachments are left untouched.
	if d.Get("create_attachment").(bool) {
		request["DeleteAttachment"] = true
	}

	request["ClientToken"] = helper.BuildClientToken(action)

	stageTimeout, err := operationTimeoutRemaining(deleteDeadline, "delete")
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	_, err = helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
		resp, callErr := trAssociationRpcPost(ctx, client, "ExpressConnectRouter", "2023-09-01", action, query, request)
		if callErr == nil {
			logs.AddDebug(action, resp, request)
		}
		return resp, callErr
	}, stageTimeout, func(err error) bool {
		return trAssociationBusinessCode(err) != "" && errs.IsExpectedErrors(err, []string{"Conflict.Lock", "IncorrectStatus.Ecr", "IncorrectStatus.EcrAssociation", "OperationDenied.ForceDeleting", "OperationFailed.CallingVpcYaochi"}) || helper.NeedRetry(err)
	})

	if err != nil {
		// A missing parent ECR implies the association no longer exists.
		if trAssociationBusinessCode(err) != "" && errs.IsExpectedErrors(err, []string{"ResourceNotFound.AssociationId", "ResourceNotFound.EcrId"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	stageTimeout, err = operationTimeoutRemaining(deleteDeadline, "delete")
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	stateConf := helper.BuildStateConf([]string{}, []string{}, stageTimeout, trAssociationStateDelay, trAssociationStateRefresh(ctx, client, d.Id()))
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}
