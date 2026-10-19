package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// CMS (Cloud Monitor Service) gap sweepers — fills the gaps left by the 12
// pre-existing inline sweepers (alarm, alarm_contact, alarm_contact_group,
// dynamic_tag_group, event_rule, group_metric_rule, hybrid_monitor_sls_task,
// metric_rule_template, monitor_group, namespace, site_monitor, sls_group in
// their own *_test.go files). CMS is a HYBRID product: the RPC resources here
// use the `2019-01-01` API (RpcPost with map[string]interface{} params,
// autoRetry=false — matching the existing cms sweepers), and the ROA resources
// use the `2024-03-30` API (RoaGet/RoaPost/RoaDelete with map[string]*string
// query, StringPointer). These are the FIRST ROA cms sweepers in the repo.
//
// Eleven gap sweepers are added here (monitor_group_instances is cascade-only —
// its parent monitor_group sweeper already disassociates items on delete, so no
// standalone sweeper is needed):
//   RPC region-scoped: metric_rule_black_list, hybrid_monitor_fc_task
//   ROA region-scoped: workspace, prometheus_instance, prometheus_view,
//     integration_policy, alert_rule_v2
//   ROA parent-scoped (walk parent list, then per-parent children):
//     addon_release (parent=integration_policy), agg_task_group
//     (parent=prometheus_instance), dataset + event_notify_policy
//     (parent=workspace)
//
// Five resources that initially looked like gaps (namespace, sls_group,
// metric_rule_template, dynamic_tag_group, hybrid_monitor_sls_task) in fact
// already had inline sweepers; they are registered with a split-line
// AddTestSweepers( \n "name",) layout that a single-line grep misses — see
// arms lesson #1. They are NOT duplicated here.
//
// System-default preserve guards (red-line #4 — never force-delete managed):
//   - workspace: skip names starting "default-cms-" (system default workspace).
//   - prometheus_instance: skip name/workspace starting "default-cms-".
//   - addon_release: skip managed==true (system-managed addons).
//   - namespace: skip "acs_" prefix (built-in namespaces — defensive; the list
//     API only returns hybrid/custom namespaces anyway).
// The pre-existing dynamic_tag_group sweeper is hardened in its own file
// (always-active tf-prefix, no SWEEP_ALL bypass) so the 17 system-managed tag
// groups (ack.aliyun.com / ehpc / acs:computenest) are never fleet-wiped.
// metric_rule_template is already always-active upstream (EMR_DEFAULT_ALARM_TEMPLATE
// + si-* 代运维 templates preserved).
//
// Pagination casing is INCONSISTENT across the ROA endpoints: alert_rule_v2 +
// dataset use lowercase nextToken/maxResults; event_notify_policy uses CAPITAL
// NextToken/MaxResults. cmsRoaList parses both casings defensively.

var cmsSweepPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf-example",
	"tf_example",
	"terraform-example",
	"tftestacc",
}

// cmsGetField returns the first non-empty string among the given keys on item,
// trying each name as-is (CMS response casing is inconsistent across RPC/ROA,
// so callers pass all known casing variants).
func cmsGetField(item map[string]interface{}, names ...string) string {
	for _, name := range names {
		if v, ok := item[name]; ok {
			if s := fmt.Sprint(v); s != "" && s != "<nil>" {
				return s
			}
		}
	}
	return ""
}

// cmsSweepNameOk is the fleet-wipe gate: never touch NODELETING names; in
// normal (non-SWEEP_ALL) mode only sweep tf-prefixed names. SWEEP_ALL bypasses
// the prefix filter so the whole test account (1511928242963727, cn-beijing) is
// fleet-wiped of non-NODELETING leftovers. System-managed instances are guarded
// per-sweeper with always-active inline checks (workspace/prometheus_instance/
// addon_release) so SWEEP_ALL still preserves them.
func cmsSweepNameOk(name string) bool {
	if name == "" {
		return false
	}
	if protectedSweepVpcName(name) {
		log.Printf("[INFO] Skipping cms resource (NODELETING): %s", name)
		return false
	}
	if !sweepAll() {
		skip := true
		for _, prefix := range cmsSweepPrefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping cms resource (no prefix): %s", name)
			return false
		}
	}
	return true
}

// cmsRpcPost issues an RpcPost (params in the BODY slot, query=nil, autoRetry=false
// — matching the existing cms RPC sweepers) with retry.Retry + NeedRetry +
// incrementalWait.
func cmsRpcPost(client *connectivity.AliyunClient, action string, params map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, params, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	return response, err
}

// cmsRpcDelete issues an RpcPost delete with retry and swallows expected
// not-found codes (already gone = success).
func cmsRpcDelete(client *connectivity.AliyunClient, action string, params map[string]interface{}, notFoundCodes []string) error {
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, err = client.RpcPost("Cms", "2019-01-01", action, nil, params, false)
		if err != nil {
			if IsExpectedErrors(err, notFoundCodes) || NotFoundError(err) {
				return nil
			}
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	return err
}

// cmsRpcList paginates an RPC list action (PageNumber/PageSize) and returns the
// slice at jpath. extras (e.g. RegionId, TaskType) are merged into each page
// request.
func cmsRpcList(client *connectivity.AliyunClient, action string, pageSize int, jpath string, extras map[string]interface{}) ([]map[string]interface{}, error) {
	var all []map[string]interface{}
	request := make(map[string]interface{})
	for k, v := range extras {
		request[k] = v
	}
	request["PageSize"] = pageSize
	request["PageNumber"] = 1
	for {
		resp, err := cmsRpcPost(client, action, request)
		if err != nil {
			return all, err
		}
		v, err := jsonpath.Get(jpath, resp)
		if err != nil {
			return all, err
		}
		items, _ := v.([]interface{})
		for _, it := range items {
			if m, ok := it.(map[string]interface{}); ok {
				all = append(all, m)
			}
		}
		if len(items) < pageSize {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return all, nil
}

// cmsRoaGet issues a RoaGet (Cms 2024-03-30) with retry. RoaGet has no autoRetry
// arg, so throttling is handled here.
func cmsRoaGet(client *connectivity.AliyunClient, path string, query map[string]*string) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("Cms", "2024-03-30", path, query, nil, nil)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	return response, err
}

// cmsRoaDelete issues a RoaDelete (autoRetry=true) with retry and swallows
// expected not-found codes.
func cmsRoaDelete(client *connectivity.AliyunClient, path string, query map[string]*string, notFoundCodes []string) error {
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, err = client.RoaDelete("Cms", "2024-03-30", path, query, nil, nil, true)
		if err != nil {
			if IsExpectedErrors(err, notFoundCodes) || NotFoundError(err) {
				return nil
			}
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	return err
}

// cmsRoaPost issues a RoaPost (autoRetry=true) with retry and swallows expected
// not-found codes. Used for alert_rule_v2 list (/queryAlertRules) + delete
// (/manageAlertRules with body {uuid, action:DELETE}).
func cmsRoaPost(client *connectivity.AliyunClient, path string, query map[string]*string, body map[string]interface{}, notFoundCodes []string) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RoaPost("Cms", "2024-03-30", path, query, nil, body, true)
		if err != nil {
			if IsExpectedErrors(err, notFoundCodes) || NotFoundError(err) {
				return nil
			}
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	return response, err
}

// cmsRoaList paginates a RoaGet list (nextToken/NextToken — casing is
// inconsistent across endpoints, so both are parsed defensively) and returns
// the slice at jpath. The caller pre-sets maxResults/MaxResults + any required
// query params (workspace, regionId, etc.) on query; this helper only manages
// the pagination token.
func cmsRoaList(client *connectivity.AliyunClient, path string, query map[string]*string, jpath string) ([]map[string]interface{}, error) {
	var all []map[string]interface{}
	for {
		resp, err := cmsRoaGet(client, path, query)
		if err != nil {
			return all, err
		}
		v, err := jsonpath.Get(jpath, resp)
		if err != nil {
			return all, err
		}
		items, _ := v.([]interface{})
		for _, it := range items {
			if m, ok := it.(map[string]interface{}); ok {
				all = append(all, m)
			}
		}
		// pagination: try lowercase nextToken, then capital NextToken
		next := ""
		if t, ok := resp["nextToken"].(string); ok && t != "" {
			next = t
			query["nextToken"] = StringPointer(next)
		} else if t, ok := resp["NextToken"].(string); ok && t != "" {
			next = t
			query["NextToken"] = StringPointer(next)
		}
		if next == "" {
			break
		}
	}
	return all, nil
}

func init() {
	resource.AddTestSweepers("alicloud_cms_metric_rule_black_list", &resource.Sweeper{
		Name: "alicloud_cms_metric_rule_black_list", F: testSweepCmsMetricRuleBlackList,
	})
	resource.AddTestSweepers("alicloud_cms_hybrid_monitor_fc_task", &resource.Sweeper{
		Name: "alicloud_cms_hybrid_monitor_fc_task", F: testSweepCmsHybridMonitorFcTask,
	})
	resource.AddTestSweepers("alicloud_cms_workspace", &resource.Sweeper{
		Name: "alicloud_cms_workspace", F: testSweepCmsWorkspace,
		Dependencies: []string{"alicloud_cms_dataset", "alicloud_cms_event_notify_policy"},
	})
	resource.AddTestSweepers("alicloud_cms_prometheus_instance", &resource.Sweeper{
		Name: "alicloud_cms_prometheus_instance", F: testSweepCmsPrometheusInstance,
		Dependencies: []string{"alicloud_cms_agg_task_group"},
	})
	resource.AddTestSweepers("alicloud_cms_prometheus_view", &resource.Sweeper{
		Name: "alicloud_cms_prometheus_view", F: testSweepCmsPrometheusView,
	})
	resource.AddTestSweepers("alicloud_cms_integration_policy", &resource.Sweeper{
		Name: "alicloud_cms_integration_policy", F: testSweepCmsIntegrationPolicy,
		Dependencies: []string{"alicloud_cms_addon_release"},
	})
	resource.AddTestSweepers("alicloud_cms_alert_rule_v2", &resource.Sweeper{
		Name: "alicloud_cms_alert_rule_v2", F: testSweepCmsAlertRuleV2,
	})
	resource.AddTestSweepers("alicloud_cms_addon_release", &resource.Sweeper{
		Name: "alicloud_cms_addon_release", F: testSweepCmsAddonRelease,
	})
	resource.AddTestSweepers("alicloud_cms_agg_task_group", &resource.Sweeper{
		Name: "alicloud_cms_agg_task_group", F: testSweepCmsAggTaskGroup,
	})
	resource.AddTestSweepers("alicloud_cms_dataset", &resource.Sweeper{
		Name: "alicloud_cms_dataset", F: testSweepCmsDataset,
	})
	resource.AddTestSweepers("alicloud_cms_event_notify_policy", &resource.Sweeper{
		Name: "alicloud_cms_event_notify_policy", F: testSweepCmsEventNotifyPolicy,
	})
}

// ---- RPC region-scoped sweepers (Cms 2019-01-01) ----

func testSweepCmsMetricRuleBlackList(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	items, err := cmsRpcList(client, "DescribeMetricRuleBlackList", PageSizeLarge, "$.DescribeMetricRuleBlackList", nil)
	if err != nil {
		log.Printf("[ERROR] Failed to list cms metric_rule_black_list: %s", err)
		return nil
	}
	for _, item := range items {
		id := cmsGetField(item, "Id")
		name := cmsGetField(item, "Name")
		if !cmsSweepNameOk(name) {
			continue
		}
		log.Printf("[INFO] Deleting cms metric_rule_black_list: %s", name)
		if e := cmsRpcDelete(client, "DeleteMetricRuleBlackList", map[string]interface{}{
			"Id": id,
		}, nil); e != nil {
			log.Printf("[ERROR] Failed to delete cms metric_rule_black_list (%s): %s", name, e)
			continue
		}
		log.Printf("[INFO] Delete cms metric_rule_black_list success: %s", name)
	}
	return nil
}

func testSweepCmsHybridMonitorFcTask(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	extras := map[string]interface{}{"TaskType": "aliyun_fc"}
	items, err := cmsRpcList(client, "DescribeHybridMonitorTaskList", PageSizeLarge, "$.TaskList", extras)
	if err != nil {
		log.Printf("[ERROR] Failed to list cms hybrid_monitor_fc_task: %s", err)
		return nil
	}
	for _, item := range items {
		taskId := cmsGetField(item, "TaskId")
		targetUserId := cmsGetField(item, "TargetUserId")
		namespace := cmsGetField(item, "Namespace")
		name := cmsGetField(item, "TaskName")
		if name == "" {
			name = namespace
		}
		if !cmsSweepNameOk(name) {
			continue
		}
		if targetUserId == "" || namespace == "" {
			log.Printf("[ERROR] Skipping cms hybrid_monitor_fc_task (missing TargetUserId/Namespace): taskId=%s", taskId)
			continue
		}
		log.Printf("[INFO] Deleting cms hybrid_monitor_fc_task: %s", name)
		if e := cmsRpcDelete(client, "DeleteHybridMonitorTask", map[string]interface{}{
			"TargetUserId": targetUserId,
			"Namespace":    namespace,
		}, nil); e != nil {
			log.Printf("[ERROR] Failed to delete cms hybrid_monitor_fc_task (%s): %s", name, e)
			continue
		}
		log.Printf("[INFO] Delete cms hybrid_monitor_fc_task success: %s", name)
	}
	return nil
}

// ---- ROA region-scoped sweepers (Cms 2024-03-30) ----

func testSweepCmsWorkspace(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	query := map[string]*string{"maxResults": StringPointer(strconv.Itoa(PageSizeLarge))}
	items, err := cmsRoaList(client, "/workspace", query, "$.workspaces")
	if err != nil {
		log.Printf("[ERROR] Failed to list cms workspace: %s", err)
		return nil
	}
	for _, item := range items {
		name := cmsGetField(item, "workspaceName")
		if strings.HasPrefix(name, "default-cms-") {
			log.Printf("[INFO] Skipping cms workspace (system default): %s", name)
			continue
		}
		if !cmsSweepNameOk(name) {
			continue
		}
		log.Printf("[INFO] Deleting cms workspace: %s", name)
		path := fmt.Sprintf("/workspace/%s", name)
		if e := cmsRoaDelete(client, path, nil, []string{"WorkspaceNotExist"}); e != nil {
			log.Printf("[ERROR] Failed to delete cms workspace (%s): %s", name, e)
			continue
		}
		log.Printf("[INFO] Delete cms workspace success: %s", name)
	}
	return nil
}

func testSweepCmsPrometheusInstance(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	query := map[string]*string{"maxResults": StringPointer(strconv.Itoa(PageSizeLarge))}
	items, err := cmsRoaList(client, "/prometheus-instances", query, "$.prometheusInstances")
	if err != nil {
		log.Printf("[ERROR] Failed to list cms prometheus_instance: %s", err)
		return nil
	}
	for _, item := range items {
		id := cmsGetField(item, "prometheusInstanceId")
		name := cmsGetField(item, "prometheusInstanceName")
		workspace := cmsGetField(item, "workspace")
		if strings.HasPrefix(name, "default-cms-") || strings.HasPrefix(workspace, "default-cms-") {
			log.Printf("[INFO] Skipping cms prometheus_instance (system default): %s", name)
			continue
		}
		if !cmsSweepNameOk(name) {
			continue
		}
		log.Printf("[INFO] Deleting cms prometheus_instance: %s", name)
		path := fmt.Sprintf("/prometheus-instances/%s", id)
		if e := cmsRoaDelete(client, path, nil, []string{"404"}); e != nil {
			log.Printf("[ERROR] Failed to delete cms prometheus_instance (%s): %s", name, e)
			continue
		}
		log.Printf("[INFO] Delete cms prometheus_instance success: %s", name)
	}
	return nil
}

func testSweepCmsPrometheusView(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	query := map[string]*string{"maxResults": StringPointer(strconv.Itoa(PageSizeLarge))}
	items, err := cmsRoaList(client, "/prometheus-views", query, "$.prometheusViews")
	if err != nil {
		log.Printf("[ERROR] Failed to list cms prometheus_view: %s", err)
		return nil
	}
	for _, item := range items {
		id := cmsGetField(item, "prometheusViewId")
		name := cmsGetField(item, "prometheusViewName")
		if !cmsSweepNameOk(name) {
			continue
		}
		log.Printf("[INFO] Deleting cms prometheus_view: %s", name)
		path := fmt.Sprintf("/prometheus-views/%s", id)
		if e := cmsRoaDelete(client, path, nil, []string{"404"}); e != nil {
			log.Printf("[ERROR] Failed to delete cms prometheus_view (%s): %s", name, e)
			continue
		}
		log.Printf("[INFO] Delete cms prometheus_view success: %s", name)
	}
	return nil
}

func testSweepCmsIntegrationPolicy(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	query := map[string]*string{"maxResults": StringPointer(strconv.Itoa(PageSizeLarge))}
	items, err := cmsRoaList(client, "/integration-policies", query, "$.policies")
	if err != nil {
		log.Printf("[ERROR] Failed to list cms integration_policy: %s", err)
		return nil
	}
	for _, item := range items {
		id := cmsGetField(item, "policyId")
		name := cmsGetField(item, "policyName")
		if !cmsSweepNameOk(name) {
			continue
		}
		log.Printf("[INFO] Deleting cms integration_policy: %s", name)
		path := fmt.Sprintf("/integration-policies/%s", id)
		// force=true: addon_release child sweeper (dependency) already removed
		// children, so force is a safety net against dangling addons.
		delQuery := map[string]*string{"force": StringPointer("true")}
		if e := cmsRoaDelete(client, path, delQuery, []string{"404", "15007"}); e != nil {
			log.Printf("[ERROR] Failed to delete cms integration_policy (%s): %s", name, e)
			continue
		}
		log.Printf("[INFO] Delete cms integration_policy success: %s", name)
	}
	return nil
}

func testSweepCmsAlertRuleV2(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	// alert_rule_v2 list is a RoaPost (/queryAlertRules); delete is a RoaPost
	// (/manageAlertRules body {uuid, action:DELETE}). Both need RegionId in the
	// query. Mirror the provider's queryAlertRules call
	// (service_alicloud_cms_v2.go, Cms 2024-03-30): the query carries ONLY
	// RegionId; the body holds any filter (provider uses filter.uuid.eq for a
	// single-id lookup). For a fleet list-all we send a minimal pagination body
	// (no filter) — an empty body {} is rejected with 400 ParameterInvalid
	// "Request body is required" even though every body field is optional.
	// Putting maxResults in the query also caused a 400 (unrecognized query
	// param for this ROA endpoint). Pagination is body-based
	// (pagination.pageNumber/pageSize); the response carries data.totalCount,
	// not a nextToken.
	var all []map[string]interface{}
	// queryAlertRules (Cms 2024-03-30 ROA POST /queryAlertRules) rejects an
	// empty body with 400 ParameterInvalid "Request body is required", even
	// though every body field (filter/pagination/workspace) is optional. Send
	// a minimal pagination body so the endpoint accepts the request; the body
	// also drives page-based pagination (pageNumber/pageSize) per the schema,
	// since the response carries data.totalCount rather than a nextToken.
	query := map[string]*string{
		"RegionId": StringPointer(client.RegionId),
	}
	body := map[string]interface{}{
		"pagination": map[string]interface{}{
			"pageNumber": 1,
			"pageSize":   PageSizeLarge,
		},
	}
	for {
		resp, err := cmsRoaPost(client, "/queryAlertRules", query, body, nil)
		if err != nil {
			log.Printf("[ERROR] Failed to list cms alert_rule_v2: %s", err)
			return nil
		}
		v, e := jsonpath.Get("$.data.alertRules[*]", resp)
		if e != nil {
			log.Printf("[ERROR] Failed to parse cms alert_rule_v2 list: %s", e)
			return nil
		}
		items, _ := v.([]interface{})
		for _, it := range items {
			if m, ok := it.(map[string]interface{}); ok {
				all = append(all, m)
			}
		}
		// body-based pagination: stop when a page returns fewer than pageSize
		// items (last page); otherwise advance pageNumber and fetch again.
		if len(items) < PageSizeLarge {
			break
		}
		pag := body["pagination"].(map[string]interface{})
		pag["pageNumber"] = pag["pageNumber"].(int) + 1
	}
	for _, item := range all {
		uuid := cmsGetField(item, "uuid")
		name := cmsGetField(item, "displayName")
		if !cmsSweepNameOk(name) {
			continue
		}
		log.Printf("[INFO] Deleting cms alert_rule_v2: %s", name)
		delQuery := map[string]*string{"RegionId": StringPointer(client.RegionId)}
		// /manageAlertRules DELETE requires the workspace id: the API rejects
		// an empty workspace with 400 ParameterInvalid "workspace id must not
		// be empty". Each queryAlertRules list item carries its workspace, so
		// source it from there (matches the provider's Delete, which sets
		// request["workspace"] in resource_alicloud_cms_alert_rule_v2.go).
		workspace := cmsGetField(item, "workspace")
		delBody := map[string]interface{}{
			"uuid":      uuid,
			"workspace": workspace,
			"action":    "DELETE",
		}
		if _, e := cmsRoaPost(client, "/manageAlertRules", delQuery, delBody, []string{"ResourceNotFound"}); e != nil {
			log.Printf("[ERROR] Failed to delete cms alert_rule_v2 (%s): %s", name, e)
			continue
		}
		log.Printf("[INFO] Delete cms alert_rule_v2 success: %s", name)
	}
	return nil
}

// ---- ROA parent-scoped sweepers (walk parent list, then per-parent children) ----

func testSweepCmsAddonRelease(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	// walk integration policies (parents)
	policies, err := cmsRoaList(client, "/integration-policies",
		map[string]*string{"maxResults": StringPointer(strconv.Itoa(PageSizeLarge))}, "$.policies")
	if err != nil {
		log.Printf("[ERROR] Failed to list cms integration_policy (addon_release parent): %s", err)
		return nil
	}
	for _, pol := range policies {
		policyId := cmsGetField(pol, "policyId")
		if policyId == "" {
			continue
		}
		// per-policy addon releases (children)
		childQuery := map[string]*string{"maxResults": StringPointer(strconv.Itoa(PageSizeLarge))}
		releases, err := cmsRoaList(client, fmt.Sprintf("/integration-policies/%s/addon-releases", policyId), childQuery, "$.releases")
		if err != nil {
			log.Printf("[ERROR] Failed to list cms addon_release under policy %s: %s", policyId, err)
			continue
		}
		for _, rel := range releases {
			releaseName := cmsGetField(rel, "releaseName")
			// managed==true → system-managed addon, never delete (red-line #4)
			if managed, ok := rel["managed"]; ok && fmt.Sprint(managed) == "true" {
				log.Printf("[INFO] Skipping cms addon_release (managed): %s", releaseName)
				continue
			}
			if !cmsSweepNameOk(releaseName) {
				continue
			}
			log.Printf("[INFO] Deleting cms addon_release: %s", releaseName)
			// delete path is the COLLECTION path with releaseName in query
			delPath := fmt.Sprintf("/integration-policies/%s/addon-releases", policyId)
			delQuery := map[string]*string{"releaseName": StringPointer(releaseName)}
			if e := cmsRoaDelete(client, delPath, delQuery, []string{"404", "12002"}); e != nil {
				log.Printf("[ERROR] Failed to delete cms addon_release (%s): %s", releaseName, e)
				continue
			}
			log.Printf("[INFO] Delete cms addon_release success: %s", releaseName)
		}
	}
	return nil
}

func testSweepCmsAggTaskGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	// walk prometheus instances (parents), skip system default-cms-*
	instances, err := cmsRoaList(client, "/prometheus-instances",
		map[string]*string{"maxResults": StringPointer(strconv.Itoa(PageSizeLarge))}, "$.prometheusInstances")
	if err != nil {
		log.Printf("[ERROR] Failed to list cms prometheus_instance (agg_task_group parent): %s", err)
		return nil
	}
	for _, inst := range instances {
		instanceId := cmsGetField(inst, "prometheusInstanceId")
		instName := cmsGetField(inst, "prometheusInstanceName")
		if instanceId == "" || strings.HasPrefix(instName, "default-cms-") {
			continue
		}
		childQuery := map[string]*string{"maxResults": StringPointer(strconv.Itoa(PageSizeLarge))}
		groups, err := cmsRoaList(client, fmt.Sprintf("/prometheus-instances/%s/agg-task-groups", instanceId), childQuery, "$.aggTaskGroups")
		if err != nil {
			log.Printf("[ERROR] Failed to list cms agg_task_group under instance %s: %s", instanceId, err)
			continue
		}
		for _, grp := range groups {
			groupId := cmsGetField(grp, "aggTaskGroupId")
			name := cmsGetField(grp, "aggTaskGroupName")
			if !cmsSweepNameOk(name) {
				continue
			}
			log.Printf("[INFO] Deleting cms agg_task_group: %s", name)
			delPath := fmt.Sprintf("/prometheus-instances/%s/agg-task-groups/%s", instanceId, groupId)
			if e := cmsRoaDelete(client, delPath, nil, []string{"404"}); e != nil {
				log.Printf("[ERROR] Failed to delete cms agg_task_group (%s): %s", name, e)
				continue
			}
			log.Printf("[INFO] Delete cms agg_task_group success: %s", name)
		}
	}
	return nil
}

func testSweepCmsDataset(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	// walk workspaces (parents), skip system default-cms-
	workspaces, err := cmsRoaList(client, "/workspace",
		map[string]*string{"maxResults": StringPointer(strconv.Itoa(PageSizeLarge))}, "$.workspaces")
	if err != nil {
		log.Printf("[ERROR] Failed to list cms workspace (dataset parent): %s", err)
		return nil
	}
	for _, ws := range workspaces {
		workspace := cmsGetField(ws, "workspaceName")
		if workspace == "" || strings.HasPrefix(workspace, "default-cms-") {
			continue
		}
		childQuery := map[string]*string{"maxResults": StringPointer(strconv.Itoa(PageSizeLarge))}
		datasets, err := cmsRoaList(client, fmt.Sprintf("/workspace/%s/dataset", workspace), childQuery, "$.datasets[*]")
		if err != nil {
			// ListWorkspaces can return stale ghost entries that 404
			// (WorkspaceNotExist) when queried directly for datasets — the
			// workspace was deleted out-of-band but the list still carries it.
			// Treat that as "no datasets here" and move on at INFO, not ERROR.
			if IsExpectedErrors(err, []string{"WorkspaceNotExist"}) {
				log.Printf("[INFO] Skip cms dataset under workspace %s: workspace gone (stale list entry)", workspace)
				continue
			}
			log.Printf("[ERROR] Failed to list cms dataset under workspace %s: %s", workspace, err)
			continue
		}
		for _, ds := range datasets {
			datasetName := cmsGetField(ds, "datasetName")
			if !cmsSweepNameOk(datasetName) {
				continue
			}
			log.Printf("[INFO] Deleting cms dataset: %s", datasetName)
			delPath := fmt.Sprintf("/workspace/%s/dataset/%s", workspace, datasetName)
			if e := cmsRoaDelete(client, delPath, nil, []string{"WorkspaceNotExist", "DatasetNotExist"}); e != nil {
				log.Printf("[ERROR] Failed to delete cms dataset (%s): %s", datasetName, e)
				continue
			}
			log.Printf("[INFO] Delete cms dataset success: %s", datasetName)
		}
	}
	return nil
}

func testSweepCmsEventNotifyPolicy(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	// walk workspaces (parents), skip system default-cms-
	workspaces, err := cmsRoaList(client, "/workspace",
		map[string]*string{"maxResults": StringPointer(strconv.Itoa(PageSizeLarge))}, "$.workspaces")
	if err != nil {
		log.Printf("[ERROR] Failed to list cms workspace (event_notify_policy parent): %s", err)
		return nil
	}
	for _, ws := range workspaces {
		workspace := cmsGetField(ws, "workspaceName")
		if workspace == "" || strings.HasPrefix(workspace, "default-cms-") {
			continue
		}
		// NOTE: event_notify_policy uses CAPITAL MaxResults/NextToken (the
		// only endpoint in this set that does).
		childQuery := map[string]*string{
			"workspace":  StringPointer(workspace),
			"regionId":   StringPointer(client.RegionId),
			"MaxResults": StringPointer(strconv.Itoa(PageSizeLarge)),
		}
		policies, err := cmsRoaList(client, "/api/eventbase/notify-policies", childQuery, "$.notifyPolicyList[*]")
		if err != nil {
			log.Printf("[ERROR] Failed to list cms event_notify_policy under workspace %s: %s", workspace, err)
			continue
		}
		for _, pol := range policies {
			uuid := cmsGetField(pol, "uuid")
			name := cmsGetField(pol, "name")
			if !cmsSweepNameOk(name) {
				continue
			}
			log.Printf("[INFO] Deleting cms event_notify_policy: %s", name)
			// delete path is the COLLECTION path with uuid+workspace+regionId in query
			delQuery := map[string]*string{
				"uuid":      StringPointer(uuid),
				"workspace": StringPointer(workspace),
				"regionId":  StringPointer(client.RegionId),
			}
			if e := cmsRoaDelete(client, "/api/eventbase/notify-policy", delQuery, []string{"NotFound"}); e != nil {
				log.Printf("[ERROR] Failed to delete cms event_notify_policy (%s): %s", name, e)
				continue
			}
			log.Printf("[INFO] Delete cms event_notify_policy success: %s", name)
		}
	}
	return nil
}
