package threatdetection

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/resourcemanager"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudThreatDetectionServiceLinkedRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAliCloudThreatDetectionServiceLinkedRoleCreate,
		ReadContext:   resourceAliCloudThreatDetectionServiceLinkedRoleRead,
		DeleteContext: resourceAliCloudThreatDetectionServiceLinkedRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"role_status": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"service_linked_role": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "AliyunServiceRoleForSas",
				ValidateFunc: validators.StringInSlice([]string{"AliyunServiceRoleForSas", "AliyunServiceRoleForSasCspm", "AliyunServiceRoleForSasAgentless", "AliyunServiceRoleForOssMfd", "AliyunServiceRoleForAntiRansomwareMssp", "AliyunServiceRoleForSasSecllm", "AliyunServiceRoleForSasAI", "AliyunLogETLRole"}, false),
			},
		},
	}
}

func resourceAliCloudThreatDetectionServiceLinkedRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateServiceLinkedRole"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	roleName := d.Get("service_linked_role").(string)
	request["ServiceLinkedRole"] = roleName
	request["ClientToken"] = helper.BuildClientToken(action)

	threatDetectionServiceV2 := ServiceV2{client}
	if _, checkErr := threatDetectionServiceV2.DescribeThreatDetectionServiceLinkedRole(roleName); checkErr == nil {
		return diag.Errorf("the service linked role %s already exists, import it with: "+
			"terraform import alicloud_threat_detection_service_linked_role.example %s", roleName, roleName)
	} else if !errs.NotFoundError(checkErr) {
		return diag.FromErr(errs.WrapErrorf(checkErr, errs.DefaultErrorMsg, "alicloud_threat_detection_service_linked_role", "DescribeServiceLinkedRoleStatus", errs.AlibabaCloudSdkGoERROR))
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
		if err != nil {
			// The request carries a ClientToken built once for every attempt, so a retry of a request the server is still
			// processing is answered with 409 IdempotentRequestInProgress rather than creating the role twice. That answer
			// means the create is in flight, not that it failed, so it is retried like a throttling error instead of being
			// reported: the first attempt can time out on the client after the server accepted it.
			if errs.IsExpectedErrors(err, []string{"IdempotentRequestInProgress"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return diag.FromErr(errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_threat_detection_service_linked_role", action, errs.AlibabaCloudSdkGoERROR))
	}

	stateConf := helper.BuildStateConf([]string{}, []string{"true"}, d.Timeout(schema.TimeoutCreate), 5*time.Second,
		threatDetectionServiceV2.DescribeAsyncThreatDetectionServiceLinkedRoleStateRefreshFunc(roleName, response, "$.RoleStatus.Status", []string{}))
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(errs.WrapErrorf(err, errs.IdMsg, roleName))
	}
	if err := waitThreatDetectionServiceLinkedRoleStatus(ctx, threatDetectionServiceV2, roleName, true); err != nil {
		return diag.FromErr(errs.WrapErrorf(err, errs.IdMsg, roleName))
	}
	d.SetId(roleName)
	return resourceAliCloudThreatDetectionServiceLinkedRoleRead(ctx, d, meta)
}

func resourceAliCloudThreatDetectionServiceLinkedRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ServiceV2{client}

	objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionServiceLinkedRole(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_service_linked_role DescribeThreatDetectionServiceLinkedRole Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return diag.FromErr(errs.WrapError(err))
	}

	d.Set("role_status", objectRaw["Status"])
	d.Set("service_linked_role", d.Id())
	return nil
}

func resourceAliCloudThreatDetectionServiceLinkedRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteServiceLinkedRole"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RoleName"] = d.Id()

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if !errs.NotFoundError(err) {
			return diag.FromErr(errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR))
		}
		// The role is already gone, but the status this resource reads still trails it, so it waits for the same
		// settle as a successful deletion; otherwise a refresh straight after destroy reports it as existing again.
		if err = waitThreatDetectionServiceLinkedRoleStatus(ctx, ServiceV2{client}, d.Id(), false); err != nil {
			return diag.FromErr(errs.WrapErrorf(err, "the service linked role %s was already deleted but its status did not settle", d.Id()))
		}
		return nil
	}

	// DeleteServiceLinkedRole is asynchronous and answers with a deletion task id, the role itself is removed by a
	// background task. Only the task result carries the real outcome, for example a role that a running SAS instance
	// still uses can never be deleted. Watching the SAS role status instead hides that reason, so a deletion that is
	// rejected from the beginning looks identical to a slow one and the user only gets a timeout.
	deletionTaskId, ok := response["DeletionTaskId"].(string)
	if !ok || deletionTaskId == "" {
		return diag.Errorf("the %s for %s returned no deletion task id: %v", action, d.Id(), response)
	}
	resourceManagerService := resourcemanager.NewResourceManagerService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"SUCCEEDED"}, d.Timeout(schema.TimeoutDelete), 0*time.Second, resourceManagerService.ResourceManagerServiceLinkedRoleStateRefreshFunc(deletionTaskId, []string{"FAILED", "INTERNAL_ERROR"}))
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		deletionTask, taskErr := resourceManagerService.GetServiceLinkedRoleDeletionStatus(deletionTaskId)
		if taskErr != nil {
			return diag.FromErr(errs.WrapErrorf(err, errs.FailedToReachTargetStatusWithError, d.Id(), taskErr))
		}
		if usages := serviceLinkedRoleDeletionUsages(deletionTask); len(usages) > 0 {
			return diag.FromErr(errs.WrapErrorf(err, "the service linked role %s cannot be deleted because it is still used by %s", d.Id(), strings.Join(usages, ", ")))
		}
		return diag.FromErr(errs.WrapErrorf(err, errs.FailedToReachTargetStatusWithResponse, d.Id(), deletionTask))
	}

	if err := waitThreatDetectionServiceLinkedRoleStatus(ctx, ServiceV2{client}, d.Id(), false); err != nil {
		return diag.FromErr(errs.WrapErrorf(err, "the service linked role %s was deleted but its status did not settle", d.Id()))
	}
	return nil
}

// serviceLinkedRoleStatusSettleTimeout bounds the settle wait on its own. It is deliberately not d.Timeout(...): the
// primary waits in create and delete already spend that budget, and reusing it would double the timeout the resource
// declares, while the status converges in about fifteen seconds in practice.
const serviceLinkedRoleStatusSettleTimeout = 1 * time.Minute

// waitThreatDetectionServiceLinkedRoleStatus waits until DescribeServiceLinkedRoleStatus, the call every read of this
// resource goes through, reports the expected status without interruption for a settle window. The SAS role status
// trails the role itself: the Resource Manager deletion task already succeeded while the status still reports the role
// as enabled, and on both create and delete it alternates between enabled and disabled for several seconds before it
// settles. Sampling it once, or twice in a row, therefore still lets a refresh right after create report the new role
// as absent and a refresh right after destroy report the deleted role as existing again.
func waitThreatDetectionServiceLinkedRoleStatus(ctx context.Context, service ServiceV2, roleName string, enabled bool) error {
	const settleWindow = 10 * time.Second
	var matchedSince time.Time

	wait := helper.IncrementalWait(1*time.Second, 3*time.Second)
	return retry.RetryContext(ctx, serviceLinkedRoleStatusSettleTimeout, func() *retry.RetryError {
		_, describeErr := service.DescribeThreatDetectionServiceLinkedRole(roleName)
		switch {
		case describeErr != nil && !errs.NotFoundError(describeErr):
			return retry.NonRetryableError(errs.WrapErrorf(describeErr, errs.DefaultErrorMsg, roleName, "DescribeServiceLinkedRoleStatus", errs.AlibabaCloudSdkGoERROR))
		case describeErr == nil && enabled, describeErr != nil && !enabled:
			if matchedSince.IsZero() {
				matchedSince = time.Now()
			}
			if time.Since(matchedSince) >= settleWindow {
				return nil
			}
		default:
			matchedSince = time.Time{}
		}

		wait()
		return retry.RetryableError(fmt.Errorf("the service linked role %s has not settled to the expected status yet", roleName))
	})
}

// serviceLinkedRoleDeletionUsages extracts the resources that still reference the role from the deletion task of
// DeleteServiceLinkedRole, for example a running SAS instance blocks the deletion of its service linked role.
func serviceLinkedRoleDeletionUsages(deletionTask map[string]interface{}) []string {
	reason, _ := deletionTask["Reason"].(map[string]interface{})
	roleUsages, _ := reason["RoleUsages"].(map[string]interface{})

	rawUsages := roleUsages["RoleUsage"]
	usages, ok := rawUsages.([]interface{})
	if !ok && rawUsages != nil {
		usages = []interface{}{rawUsages}
	}

	var resourceIds []string
	for _, usage := range usages {
		usageMap, ok := usage.(map[string]any)
		if !ok {
			continue
		}

		resources, _ := usageMap["Resources"].(map[string]interface{})
		rawResources := resources["Resource"]
		resourceList, ok := rawResources.([]interface{})
		if !ok && rawResources != nil {
			resourceList = []interface{}{rawResources}
		}

		for _, resourceId := range resourceList {
			resourceIds = append(resourceIds, fmt.Sprint(resourceId))
		}
	}
	return resourceIds
}
