package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// Threat Detection (云安全中心/SAS, product "Sas" / "2018-12-03") gap sweepers.
//
// All 44 alicloud_threat_detection_* resources had ZERO pre-existing sweepers
// (verified: grep AddTestSweepers across alicloud/resource_alicloud_threat_detection_*_test.go
// returns nothing). The product is RPC-style: every call is
// client.RpcPost("Sas", "2018-12-03", action, query, body, true) with params in
// the BODY slot (query=nil). The service structs (ThreatDetectionService /
// ThreatDetectionServiceV2) only expose Describe* (Get-by-id) + a handful of
// List* helpers for the data sources, so the sweepers call the list/search APIs
// directly.
//
// Fifteen gap sweepers are added here for resources that have BOTH an
// enumerable List API AND a destructive Delete API:
//
//   anti_brute_force_rule    List DescribeAntiBruteForceRules  -> $.Rules[*] (Id/Name)       Delete DeleteAntiBruteForceRule (Ids.1)
//   attack_path_whitelist    List ListAttackPathWhitelist      -> $.AttackPathWhitelistList[*] (AttackPathWhitelistId/WhitelistName) Delete DeleteAttackPathWhitelist (AttackPathWhitelistId)
//   backup_policy            List DescribeBackupPolicies       -> $.Policies[*] (Id/Name)     Delete DeleteBackupPolicy (Id)
//   baseline_strategy        List DescribeStrategy             -> $.Strategies[*] (Id/Name)   Delete DeleteStrategy (Id)
//   client_file_protect      List ListFileProtectRule         -> $.FileProtectList[*] (Id/RuleName) Delete DeleteFileProtectRule (Id.1)
//   client_user_define_rule  List ListClientUserDefineRules   -> $.UserDefineRuleList[*] (Id/Name) Delete DeleteClientUserDefineRule (IdList.1)
//   cycle_task               List DescribeCycleTaskList       -> $.CycleScheduleResponseList[*] (ConfigId/TaskName) Delete DeleteCycleTask (ConfigId)
//   honey_pot                List ListHoneypot                -> $.List[*] (HoneypotId/HoneypotName) Delete DeleteHoneypot (HoneypotId)
//   honeypot_node            List ListHoneypotNode            -> $.HoneypotNodeList[*] (NodeId/NodeName) Delete DeleteHoneypotNode (NodeId)
//   honeypot_preset          List ListHoneypotPreset          -> $.List[*] (HoneypotPresetId/PresetName) Delete DeleteHoneypotPreset (HoneypotPresetId)
//   honeypot_probe           List ListHoneypotProbe           -> $.List[*] (ProbeId/DisplayName) Delete DeleteHoneypotProbe (ProbeId)
//   malicious_file_whitelist_config List ListMaliciousFileWhitelistConfigs -> $.List[*] (Id/Remark) Delete DeleteMaliciousFileWhitelistConfig (ConfigId in QUERY)
//   oss_scan_config          List ListOssScanConfig            -> $.Data[*] (Id/Name) Delete DeleteOssScanConfig (Id)
//   vul_whitelist            List DescribeVulWhitelist        -> $.VulWhitelists[*] (Id/Whitelist) Delete DeleteVulWhitelist (Id)
//   web_lock_config          List DescribeWebLockBindList      -> $.BindList[*] (Uuid/Uuid) Delete ModifyWebLockUnbind (Uuid)
//
// Honeypot DAG (children-first via Dependencies): honeypot_probe and
// honeypot_preset are children of honeypot_node (probes bind to honeypots,
// presets live on nodes); honey_pot is a child of honeypot_node (honeypots run
// on nodes). So probe/preset/honey_pot sweepers run before the node sweeper,
// and the probe sweeper runs before the honey_pot sweeper.
//
// DEFERRED resources (no sweeper; reason inline):
//   asset_bind                       Delete is a no-op (returns nil immediately).
//   asset_selection_config           Delete is a no-op (returns nil immediately).
//   attack_path_sensitive_asset_config  Has Delete but NO List API — only GetAttackPathSensitiveAssetConfig (single-item get by id); cannot enumerate.
//   check_config                     Delete is a no-op (ChangeCheckConfig mutates config; no destructive delete).
//   file_upload_limit                Delete is a no-op (returns nil immediately).
//   honeypot_images                  Data source only (ListAvailableHoneypot) — no singular resource, read-only.
//   image_event_operation            Has Delete but NO List API — only GetImageEventOperation (single-item get by id); cannot enumerate.
//   instance                         Subscription/PAYG instance; Delete is a no-op (returns nil). Cleanup requires BSS RefundInstance (red-line: no refund/退订).
//   log_meta                         Delete is a no-op (ModifyLogMetaStatus toggles status; no destructive delete).
//   log_shipper                      Data source wraps ModifyOpenLogShipper (a mutation, not a list+delete pair).
//   monitor_account                  Delete is a no-op (returns nil immediately).
//   rd_default_sync_list             Account-level singleton (ListRdDefaultSyncList returns ONE Data object, not a list); no name field to gate on (only numeric FolderIds); "delete" is a destructive clear-all (CreateRdDefaultSyncList with empty FolderIds) that would remove ALL folders, not just test-created ones.
//   sas_trail                        Has Delete but NO List API — only GetServiceTrail (single-item get by id); cannot enumerate.
//   service_linked_role              Has Delete but NO List API — only DescribeServiceLinkedRoleStatus (single status check); cannot enumerate.
// The plural data sources (_rules, _assets, _backup_policies, _baseline_strategies,
// _check_item_configs, _check_structures, _honey_pots, _honeypot_images,
// _honeypot_nodes, _honeypot_presets, _honeypot_probes, _instances, _vul_whitelists,
// _web_lock_configs) are read-only data sources and are inherently DEFERRED (no
// singular resource to delete).

var threatSweepPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf-example",
	"tf_example",
	"terraform-example",
	"tftestacc",
	"testAcc",
}

// threatGetField returns the first non-empty value among the given keys on item,
// trying both as-is and lowercased (SAS response casing is inconsistent).
func threatGetField(item map[string]interface{}, names ...string) string {
	for _, name := range names {
		if v, ok := item[name]; ok {
			if s := fmt.Sprint(v); s != "" && s != "<nil>" {
				return s
			}
		}
	}
	return ""
}

// threatSweepNameOk is the fleet-wipe gate: never touch NODELETING names; in
// normal (non-SWEEP_ALL) mode only sweep tf-prefixed names. SWEEP_ALL bypasses
// the prefix filter so the whole test account (cn-beijing) is fleet-wiped.
// Resources without a name field (e.g. web_lock_config whose only identifier
// is a server Uuid) pass the Uuid as the name; in normal mode such resources
// are skipped (uuids are not tf-prefixed) and only swept in SWEEP_ALL mode.
func threatSweepNameOk(name string) bool {
	if name == "" {
		return false
	}
	if protectedSweepVpcName(name) {
		log.Printf("[INFO] Skipping threat_detection resource (NODELETING): %s", name)
		return false
	}
	if !sweepAll() {
		skip := true
		for _, prefix := range threatSweepPrefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping threat_detection resource (no prefix): %s", name)
			return false
		}
	}
	return true
}

// threatRpcPost issues an RpcPost with retry.Retry + NeedRetry +
// incrementalWait. Params go in the BODY slot (query=nil), matching the
// threat_detection service pattern (RpcPost signature: product, version, action,
// query, body, autoRetry).
func threatRpcPost(client *connectivity.AliyunClient, action string, params map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, nil, params, true)
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

// threatDeletePost issues an RpcPost delete with retry. Params go in the BODY
// slot (query=nil). NotFoundError and common not-found error codes are treated
// as success (resource already gone).
func threatDeletePost(client *connectivity.AliyunClient, action string, params map[string]interface{}) error {
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, err = client.RpcPost("Sas", "2018-12-03", action, nil, params, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if NotFoundError(err) || IsExpectedErrors(err, []string{"DataNotExists", "NotFound", "ResourceNotFound", "ItemNotFound"}) {
			return nil
		}
		return err
	}
	return nil
}

// threatDeletePostQuery issues an RpcPost delete with params in the QUERY slot
// (body=nil). Used by malicious_file_whitelist_config (ConfigId in query).
func threatDeletePostQuery(client *connectivity.AliyunClient, action string, queryParams map[string]interface{}) error {
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, err = client.RpcPost("Sas", "2018-12-03", action, queryParams, map[string]interface{}{}, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if NotFoundError(err) || IsExpectedErrors(err, []string{"DataNotExists", "NotFound", "ResourceNotFound"}) {
			return nil
		}
		return err
	}
	return nil
}

// threatListPaged enumerates a paginated list API (CurrentPage/PageSize style).
// listPath is a jsonpath expression (e.g. "$.Rules", "$.List", "$.Data") that
// resolves to the array of items. extra params are merged into each page
// request. Pagination stops when a page returns fewer than pageSize items.
func threatListPaged(client *connectivity.AliyunClient, action, listPath string, pageSize int, extra map[string]interface{}) ([]map[string]interface{}, error) {
	items := make([]map[string]interface{}, 0)
	currentPage := 1
	for {
		params := map[string]interface{}{
			"CurrentPage": currentPage,
			"PageSize":    pageSize,
		}
		for k, v := range extra {
			params[k] = v
		}
		response, err := threatRpcPost(client, action, params)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", action, err)
		}
		raw, err := jsonpath.Get(listPath, response)
		if err != nil {
			return items, nil
		}
		list, _ := raw.([]interface{})
		if len(list) == 0 {
			break
		}
		for _, v := range list {
			if m, ok := v.(map[string]interface{}); ok {
				items = append(items, m)
			}
		}
		if len(list) < pageSize {
			break
		}
		currentPage++
	}
	return items, nil
}

// threatListOnce enumerates a non-paginated list API (single call). Used by
// DescribeAntiBruteForceRules and DescribeStrategy which return all items in
// one response without CurrentPage/PageSize pagination.
func threatListOnce(client *connectivity.AliyunClient, action, listPath string, params map[string]interface{}) ([]map[string]interface{}, error) {
	items := make([]map[string]interface{}, 0)
	response, err := threatRpcPost(client, action, params)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", action, err)
	}
	raw, err := jsonpath.Get(listPath, response)
	if err != nil {
		return items, nil
	}
	list, _ := raw.([]interface{})
	for _, v := range list {
		if m, ok := v.(map[string]interface{}); ok {
			items = append(items, m)
		}
	}
	return items, nil
}

func init() {
	resource.AddTestSweepers("alicloud_threat_detection_anti_brute_force_rule", &resource.Sweeper{
		Name: "alicloud_threat_detection_anti_brute_force_rule", F: testSweepThreatDetectionAntiBruteForceRules,
	})
	resource.AddTestSweepers("alicloud_threat_detection_attack_path_whitelist", &resource.Sweeper{
		Name: "alicloud_threat_detection_attack_path_whitelist", F: testSweepThreatDetectionAttackPathWhitelists,
	})
	resource.AddTestSweepers("alicloud_threat_detection_backup_policy", &resource.Sweeper{
		Name: "alicloud_threat_detection_backup_policy", F: testSweepThreatDetectionBackupPolicies,
	})
	resource.AddTestSweepers("alicloud_threat_detection_baseline_strategy", &resource.Sweeper{
		Name: "alicloud_threat_detection_baseline_strategy", F: testSweepThreatDetectionBaselineStrategies,
	})
	resource.AddTestSweepers("alicloud_threat_detection_client_file_protect", &resource.Sweeper{
		Name: "alicloud_threat_detection_client_file_protect", F: testSweepThreatDetectionClientFileProtects,
	})
	resource.AddTestSweepers("alicloud_threat_detection_client_user_define_rule", &resource.Sweeper{
		Name: "alicloud_threat_detection_client_user_define_rule", F: testSweepThreatDetectionClientUserDefineRules,
	})
	resource.AddTestSweepers("alicloud_threat_detection_cycle_task", &resource.Sweeper{
		Name: "alicloud_threat_detection_cycle_task", F: testSweepThreatDetectionCycleTasks,
	})
	// honeypot_node is the parent of honey_pot, honeypot_preset, and
	// honeypot_probe — children must be swept first.
	resource.AddTestSweepers("alicloud_threat_detection_honeypot_node", &resource.Sweeper{
		Name: "alicloud_threat_detection_honeypot_node", F: testSweepThreatDetectionHoneypotNodes,
		Dependencies: []string{
			"alicloud_threat_detection_honeypot_probe",
			"alicloud_threat_detection_honey_pot",
			"alicloud_threat_detection_honeypot_preset",
		},
	})
	// honey_pot depends on honeypot_probe (probes bind to honeypots).
	resource.AddTestSweepers("alicloud_threat_detection_honey_pot", &resource.Sweeper{
		Name: "alicloud_threat_detection_honey_pot", F: testSweepThreatDetectionHoneyPots,
		Dependencies: []string{
			"alicloud_threat_detection_honeypot_probe",
		},
	})
	resource.AddTestSweepers("alicloud_threat_detection_honeypot_preset", &resource.Sweeper{
		Name: "alicloud_threat_detection_honeypot_preset", F: testSweepThreatDetectionHoneypotPresets,
	})
	resource.AddTestSweepers("alicloud_threat_detection_honeypot_probe", &resource.Sweeper{
		Name: "alicloud_threat_detection_honeypot_probe", F: testSweepThreatDetectionHoneypotProbes,
	})
	resource.AddTestSweepers("alicloud_threat_detection_malicious_file_whitelist_config", &resource.Sweeper{
		Name: "alicloud_threat_detection_malicious_file_whitelist_config", F: testSweepThreatDetectionMaliciousFileWhitelistConfigs,
	})
	resource.AddTestSweepers("alicloud_threat_detection_oss_scan_config", &resource.Sweeper{
		Name: "alicloud_threat_detection_oss_scan_config", F: testSweepThreatDetectionOssScanConfigs,
	})
	resource.AddTestSweepers("alicloud_threat_detection_vul_whitelist", &resource.Sweeper{
		Name: "alicloud_threat_detection_vul_whitelist", F: testSweepThreatDetectionVulWhitelists,
	})
	resource.AddTestSweepers("alicloud_threat_detection_web_lock_config", &resource.Sweeper{
		Name: "alicloud_threat_detection_web_lock_config", F: testSweepThreatDetectionWebLockConfigs,
	})
}

// ---- sweepers ----

// testSweepThreatDetectionAntiBruteForceRules: DescribeAntiBruteForceRules
// (non-paginated) -> $.Rules; delete DeleteAntiBruteForceRule (Ids.1).
func testSweepThreatDetectionAntiBruteForceRules(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	rules, err := threatListOnce(c, "DescribeAntiBruteForceRules", "$.Rules", map[string]interface{}{})
	if err != nil {
		log.Printf("[ERROR] threat_detection list anti_brute_force_rules: %s", err)
		return nil
	}
	for _, rule := range rules {
		id := threatGetField(rule, "Id", "id")
		name := threatGetField(rule, "Name", "name")
		if !threatSweepNameOk(name) && !threatSweepNameOk(id) {
			continue
		}
		if e := threatDeletePost(c, "DeleteAntiBruteForceRule", map[string]interface{}{
			"Ids.1": id,
		}); e != nil {
			log.Printf("[ERROR] Failed to delete threat_detection anti_brute_force_rule (%s): %s", id, e)
			continue
		}
	}
	return nil
}

// testSweepThreatDetectionAttackPathWhitelists: ListAttackPathWhitelist ->
// $.AttackPathWhitelistList; delete DeleteAttackPathWhitelist
// (AttackPathWhitelistId, handles DataNotExists).
func testSweepThreatDetectionAttackPathWhitelists(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	items, err := threatListPaged(c, "ListAttackPathWhitelist", "$.AttackPathWhitelistList[*]", PageSizeLarge, nil)
	if err != nil {
		log.Printf("[ERROR] threat_detection list attack_path_whitelist: %s", err)
		return nil
	}
	for _, item := range items {
		id := threatGetField(item, "AttackPathWhitelistId", "attackPathWhitelistId", "Id")
		name := threatGetField(item, "WhitelistName", "whitelistName", "PathName", "Remark")
		if !threatSweepNameOk(name) && !threatSweepNameOk(id) {
			continue
		}
		if e := threatDeletePost(c, "DeleteAttackPathWhitelist", map[string]interface{}{
			"AttackPathWhitelistId": id,
		}); e != nil {
			log.Printf("[ERROR] Failed to delete threat_detection attack_path_whitelist (%s): %s", id, e)
			continue
		}
	}
	return nil
}

// testSweepThreatDetectionBackupPolicies: DescribeBackupPolicies ->
// $.Policies; delete DeleteBackupPolicy (Id).
func testSweepThreatDetectionBackupPolicies(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	policies, err := threatListPaged(c, "DescribeBackupPolicies", "$.Policies[*]", PageSizeLarge, nil)
	if err != nil {
		log.Printf("[ERROR] threat_detection list backup_policy: %s", err)
		return nil
	}
	for _, policy := range policies {
		id := threatGetField(policy, "Id", "id")
		name := threatGetField(policy, "Name", "name")
		if !threatSweepNameOk(name) && !threatSweepNameOk(id) {
			continue
		}
		if e := threatDeletePost(c, "DeleteBackupPolicy", map[string]interface{}{
			"Id": id,
		}); e != nil {
			log.Printf("[ERROR] Failed to delete threat_detection backup_policy (%s): %s", id, e)
			continue
		}
	}
	return nil
}

// testSweepThreatDetectionBaselineStrategies: DescribeStrategy (non-paginated)
// -> $.Strategies; delete DeleteStrategy (Id).
func testSweepThreatDetectionBaselineStrategies(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	strategies, err := threatListOnce(c, "DescribeStrategy", "$.Strategies", map[string]interface{}{})
	if err != nil {
		log.Printf("[ERROR] threat_detection list baseline_strategy: %s", err)
		return nil
	}
	for _, strategy := range strategies {
		id := threatGetField(strategy, "Id", "id")
		name := threatGetField(strategy, "Name", "name")
		if !threatSweepNameOk(name) && !threatSweepNameOk(id) {
			continue
		}
		if e := threatDeletePost(c, "DeleteStrategy", map[string]interface{}{
			"Id": id,
		}); e != nil {
			log.Printf("[ERROR] Failed to delete threat_detection baseline_strategy (%s): %s", id, e)
			continue
		}
	}
	return nil
}

// testSweepThreatDetectionClientFileProtects: ListFileProtectRule ->
// $.FileProtectList; delete DeleteFileProtectRule (Id.1).
func testSweepThreatDetectionClientFileProtects(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	rules, err := threatListPaged(c, "ListFileProtectRule", "$.FileProtectList[*]", PageSizeLarge, nil)
	if err != nil {
		log.Printf("[ERROR] threat_detection list client_file_protect: %s", err)
		return nil
	}
	for _, rule := range rules {
		id := threatGetField(rule, "Id", "id")
		name := threatGetField(rule, "RuleName", "ruleName", "Name")
		if !threatSweepNameOk(name) && !threatSweepNameOk(id) {
			continue
		}
		if e := threatDeletePost(c, "DeleteFileProtectRule", map[string]interface{}{
			"Id.1": id,
		}); e != nil {
			log.Printf("[ERROR] Failed to delete threat_detection client_file_protect (%s): %s", id, e)
			continue
		}
	}
	return nil
}

// testSweepThreatDetectionClientUserDefineRules: ListClientUserDefineRules ->
// $.UserDefineRuleList; delete DeleteClientUserDefineRule (IdList.1).
func testSweepThreatDetectionClientUserDefineRules(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	rules, err := threatListPaged(c, "ListClientUserDefineRules", "$.UserDefineRuleList[*]", PageSizeLarge, nil)
	if err != nil {
		log.Printf("[ERROR] threat_detection list client_user_define_rule: %s", err)
		return nil
	}
	for _, rule := range rules {
		id := threatGetField(rule, "Id", "id")
		name := threatGetField(rule, "Name", "name")
		if !threatSweepNameOk(name) && !threatSweepNameOk(id) {
			continue
		}
		if e := threatDeletePost(c, "DeleteClientUserDefineRule", map[string]interface{}{
			"IdList.1": id,
		}); e != nil {
			log.Printf("[ERROR] Failed to delete threat_detection client_user_define_rule (%s): %s", id, e)
			continue
		}
	}
	return nil
}

// testSweepThreatDetectionCycleTasks: DescribeCycleTaskList ->
// $.CycleScheduleResponseList; delete DeleteCycleTask (ConfigId). No ConfigId
// is passed to the list call so ALL cycle tasks are enumerated.
func testSweepThreatDetectionCycleTasks(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	tasks, err := threatListPaged(c, "DescribeCycleTaskList", "$.CycleScheduleResponseList[*]", PageSizeLarge, nil)
	if err != nil {
		log.Printf("[ERROR] threat_detection list cycle_task: %s", err)
		return nil
	}
	for _, task := range tasks {
		id := threatGetField(task, "ConfigId", "configId")
		name := threatGetField(task, "TaskName", "taskName")
		if !threatSweepNameOk(name) && !threatSweepNameOk(id) {
			continue
		}
		if e := threatDeletePost(c, "DeleteCycleTask", map[string]interface{}{
			"ConfigId": id,
		}); e != nil {
			log.Printf("[ERROR] Failed to delete threat_detection cycle_task (%s): %s", id, e)
			continue
		}
	}
	return nil
}

// testSweepThreatDetectionHoneyPots: ListHoneypot -> $.List; delete
// DeleteHoneypot (HoneypotId). Runs after honeypot_probe (probes bind to
// honeypots) and before honeypot_node (honeypots run on nodes).
func testSweepThreatDetectionHoneyPots(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	honeypots, err := threatListPaged(c, "ListHoneypot", "$.List[*]", PageSizeMedium, nil)
	if err != nil {
		log.Printf("[ERROR] threat_detection list honey_pot: %s", err)
		return nil
	}
	for _, honeypot := range honeypots {
		id := threatGetField(honeypot, "HoneypotId", "honeypotId")
		name := threatGetField(honeypot, "HoneypotName", "honeypotName")
		if !threatSweepNameOk(name) && !threatSweepNameOk(id) {
			continue
		}
		if e := threatDeletePost(c, "DeleteHoneypot", map[string]interface{}{
			"HoneypotId": id,
		}); e != nil {
			log.Printf("[ERROR] Failed to delete threat_detection honey_pot (%s): %s", id, e)
			continue
		}
	}
	return nil
}

// testSweepThreatDetectionHoneypotNodes: ListHoneypotNode ->
// $.HoneypotNodeList; delete DeleteHoneypotNode (NodeId). Runs last in the
// honeypot DAG (after probe/preset/honey_pot children).
func testSweepThreatDetectionHoneypotNodes(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	nodes, err := threatListPaged(c, "ListHoneypotNode", "$.HoneypotNodeList[*]", PageSizeLarge, nil)
	if err != nil {
		log.Printf("[ERROR] threat_detection list honeypot_node: %s", err)
		return nil
	}
	for _, node := range nodes {
		id := threatGetField(node, "NodeId", "nodeId")
		name := threatGetField(node, "NodeName", "nodeName")
		if !threatSweepNameOk(name) && !threatSweepNameOk(id) {
			continue
		}
		if e := threatDeletePost(c, "DeleteHoneypotNode", map[string]interface{}{
			"NodeId": id,
		}); e != nil {
			log.Printf("[ERROR] Failed to delete threat_detection honeypot_node (%s): %s", id, e)
			continue
		}
	}
	return nil
}

// testSweepThreatDetectionHoneypotPresets: ListHoneypotPreset -> $.List;
// delete DeleteHoneypotPreset (HoneypotPresetId).
func testSweepThreatDetectionHoneypotPresets(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	presets, err := threatListPaged(c, "ListHoneypotPreset", "$.List[*]", PageSizeLarge, nil)
	if err != nil {
		log.Printf("[ERROR] threat_detection list honeypot_preset: %s", err)
		return nil
	}
	for _, preset := range presets {
		id := threatGetField(preset, "HoneypotPresetId", "honeypotPresetId")
		name := threatGetField(preset, "PresetName", "presetName")
		if !threatSweepNameOk(name) && !threatSweepNameOk(id) {
			continue
		}
		if e := threatDeletePost(c, "DeleteHoneypotPreset", map[string]interface{}{
			"HoneypotPresetId": id,
		}); e != nil {
			log.Printf("[ERROR] Failed to delete threat_detection honeypot_preset (%s): %s", id, e)
			continue
		}
	}
	return nil
}

// testSweepThreatDetectionHoneypotProbes: ListHoneypotProbe -> $.List; delete
// DeleteHoneypotProbe (ProbeId). Runs first in the honeypot DAG (probes bind
// to honeypots, so they must be deleted before honeypots).
func testSweepThreatDetectionHoneypotProbes(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	probes, err := threatListPaged(c, "ListHoneypotProbe", "$.List[*]", PageSizeMedium, nil)
	if err != nil {
		log.Printf("[ERROR] threat_detection list honeypot_probe: %s", err)
		return nil
	}
	for _, probe := range probes {
		id := threatGetField(probe, "ProbeId", "probeId")
		name := threatGetField(probe, "DisplayName", "displayName")
		if !threatSweepNameOk(name) && !threatSweepNameOk(id) {
			continue
		}
		if e := threatDeletePost(c, "DeleteHoneypotProbe", map[string]interface{}{
			"ProbeId": id,
		}); e != nil {
			log.Printf("[ERROR] Failed to delete threat_detection honeypot_probe (%s): %s", id, e)
			continue
		}
	}
	return nil
}

// testSweepThreatDetectionMaliciousFileWhitelistConfigs:
// ListMaliciousFileWhitelistConfigs -> $.List; delete
// DeleteMaliciousFileWhitelistConfig (ConfigId in QUERY slot — the resource's
// Delete puts ConfigId in query, not body).
func testSweepThreatDetectionMaliciousFileWhitelistConfigs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	configs, err := threatListPaged(c, "ListMaliciousFileWhitelistConfigs", "$.List[*]", PageSizeLarge, nil)
	if err != nil {
		log.Printf("[ERROR] threat_detection list malicious_file_whitelist_config: %s", err)
		return nil
	}
	for _, config := range configs {
		id := threatGetField(config, "Id", "id")
		name := threatGetField(config, "Remark", "remark", "Name", "name")
		if !threatSweepNameOk(name) && !threatSweepNameOk(id) {
			continue
		}
		if e := threatDeletePostQuery(c, "DeleteMaliciousFileWhitelistConfig", map[string]interface{}{
			"ConfigId": id,
		}); e != nil {
			log.Printf("[ERROR] Failed to delete threat_detection malicious_file_whitelist_config (%s): %s", id, e)
			continue
		}
	}
	return nil
}

// testSweepThreatDetectionOssScanConfigs: ListOssScanConfig -> $.Data; delete
// DeleteOssScanConfig (Id).
func testSweepThreatDetectionOssScanConfigs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	configs, err := threatListPaged(c, "ListOssScanConfig", "$.Data[*]", PageSizeLarge, nil)
	if err != nil {
		log.Printf("[ERROR] threat_detection list oss_scan_config: %s", err)
		return nil
	}
	for _, config := range configs {
		id := threatGetField(config, "Id", "id")
		name := threatGetField(config, "Name", "name")
		if !threatSweepNameOk(name) && !threatSweepNameOk(id) {
			continue
		}
		if e := threatDeletePost(c, "DeleteOssScanConfig", map[string]interface{}{
			"Id": id,
		}); e != nil {
			log.Printf("[ERROR] Failed to delete threat_detection oss_scan_config (%s): %s", id, e)
			continue
		}
	}
	return nil
}

// testSweepThreatDetectionVulWhitelists: DescribeVulWhitelist ->
// $.VulWhitelists; delete DeleteVulWhitelist (Id).
func testSweepThreatDetectionVulWhitelists(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	whitelists, err := threatListPaged(c, "DescribeVulWhitelist", "$.VulWhitelists[*]", PageSizeLarge, nil)
	if err != nil {
		log.Printf("[ERROR] threat_detection list vul_whitelist: %s", err)
		return nil
	}
	for _, wl := range whitelists {
		id := threatGetField(wl, "Id", "id")
		name := threatGetField(wl, "Name", "name", "Whitelist", "Reason", "Remark")
		if !threatSweepNameOk(name) && !threatSweepNameOk(id) {
			continue
		}
		if e := threatDeletePost(c, "DeleteVulWhitelist", map[string]interface{}{
			"Id": id,
		}); e != nil {
			log.Printf("[ERROR] Failed to delete threat_detection vul_whitelist (%s): %s", id, e)
			continue
		}
	}
	return nil
}

// testSweepThreatDetectionWebLockConfigs: DescribeWebLockBindList ->
// $.BindList; delete ModifyWebLockUnbind (Uuid). web_lock_config has no name
// field — the server Uuid is the only identifier, so it is used as the gate
// name. In normal (non-SWEEP_ALL) mode uuids are not tf-prefixed and are
// skipped; only SWEEP_ALL fleet-wipe unbinds all web lock configs.
func testSweepThreatDetectionWebLockConfigs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	configs, err := threatListPaged(c, "DescribeWebLockBindList", "$.BindList[*]", PageSizeLarge, nil)
	if err != nil {
		log.Printf("[ERROR] threat_detection list web_lock_config: %s", err)
		return nil
	}
	for _, config := range configs {
		uuid := threatGetField(config, "Uuid", "uuid")
		if !threatSweepNameOk(uuid) {
			continue
		}
		if e := threatDeletePost(c, "ModifyWebLockUnbind", map[string]interface{}{
			"Uuid": uuid,
		}); e != nil {
			log.Printf("[ERROR] Failed to delete threat_detection web_lock_config (%s): %s", uuid, e)
			continue
		}
	}
	return nil
}
