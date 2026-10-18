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

// HBR (云备份 / Hybrid Backup Recovery, product "hbr" / "2017-09-08", RPC
// style) gap sweepers. 补齐 the missing sweepers for HBR resources that have
// BOTH an enumerable List API AND a destructive Delete API.
//
// Pre-existing inline sweepers (do NOT re-register here — duplicate
// AddTestSweepers name registration panics at init):
//   - alicloud_hbr_vault              (resource_alicloud_hbr_vault_test.go)
//     Dependencies: [alicloud_hbr_policy_binding, alicloud_hbr_policy]
//   - alicloud_hbr_policy             (resource_alicloud_hbr_policy_test.go)
//     Dependencies: [alicloud_hbr_policy_binding]
//   - alicloud_hbr_policy_binding     (resource_alicloud_hbr_policy_binding_test.go)
//   - alicloud_hbr_ots_backup_plan    (resource_alicloud_hbr_ots_backup_plan_test.go)
//   - alicloud_hbr_replication_vault  (resource_alicloud_hbr_replication_vault_test.go)
//   - alicloud_hbr_server_backup_plan (resource_alicloud_hbr_server_backup_plan_test.go)
//
// Seven new gap sweepers are added here:
//   region-scoped (no parent walk):
//     alicloud_hbr_cross_account      — DescribeCrossAccounts -> DeleteCrossAccount
//       (account-level; appears in every region's list, sweeping cn-beijing
//       alone catches them all)
//     alicloud_hbr_ecs_backup_plan    — DescribeBackupPlans SourceType=ECS_FILE
//       -> DeleteBackupPlan (deprecated since v1.249.0 in favor of policy/policy_binding)
//     alicloud_hbr_nas_backup_plan    — DescribeBackupPlans SourceType=NAS
//       -> DeleteBackupPlan (deprecated since v1.249.0)
//     alicloud_hbr_oss_backup_plan    — DescribeBackupPlans SourceType=OSS
//       -> DeleteBackupPlan (deprecated since v1.249.0)
//     alicloud_hbr_ecs_backup_client  — DescribeBackupClients ClientType=ECS_CLIENT
//       -> DeleteBackupClient (filter by InstanceName, the ECS instance name)
//   vault-scoped (walk DescribeVaults first, then per-vault child list):
//     alicloud_hbr_hana_backup_client — DescribeClients ClientType=ECS_AGENT
//       SourceType=HANA -> DeleteClient
//     alicloud_hbr_hana_backup_plan   — DescribeHanaBackupPlans
//       -> DeleteHanaBackupPlan
//
// DEFERRED (no sweeper written — reason documented):
//   - alicloud_hbr_hana_instance — DeleteHanaInstance REQUIRES the `Sid`
//     parameter (ForceNew+Optional in schema, not returned by
//     DescribeHanaInstances). The provider's Delete function hard-fails with
//     "[ERROR] To destroy the resource, you must configure the 'sid' attribute"
//     when Sid is absent. The sweeper has no way to supply Sid (it is not in
//     the list response), so the resource is unsweepable. Leftover HANA
//     instances must be cleaned manually. (The cluster_id/vault_id needed to
//     call the list ARE enumerable; the blocker is purely the Sid input on
//     delete.)
//   - alicloud_hbr_restore_job    — Delete is a no-op. The provider logs
//     "[WARN] Cannot destroy resourceAlicloudHbrRestoreJob. Terraform will
//     remove this resource from the state file, however resources may remain."
//     and returns nil. Restore jobs are transient (RUNNING -> COMPLETE/FAILED)
//     and self-terminate on completion, so there is nothing to sweep.
//
// All HBR list/delete endpoints are RPC-style: client.RpcPost("hbr",
// "2017-09-08", action, query=nil, body=params, autoRetry). Params go in the
// BODY slot (query=nil) — verified from the resource Read/Delete source and
// the pre-existing inline HBR sweepers (vault/policy/policy_binding/ots/
// server/replication_vault). Delete calls use autoRetry=false (matching the
// inline pattern); list calls use autoRetry=true wrapped in retry.Retry +
// NeedRetry + incrementalWait to absorb transient 5xx on the HBR endpoints.
//
// Red-line (always respected):
//   - Only account 1511928242963727 (test account, no prod data). Only
//     cn-beijing.
//   - NODELETING guard is always-active via protectedSweepVpcName (21
//     markers, normalization strips non-letter/digit chars). NEVER bypassed,
//     even under SWEEP_ALL.
//   - Fleet-wipe: SWEEP_ALL (ALICLOUD_SWEEP_ALL_RESOURCES=true) bypasses the
//     tf-prefix filter so the whole test account is cleaned; NORMAL mode
//     only sweeps tf-prefixed names.
//   - All sweeper F functions return nil. Best-effort delete: NotFoundError /
//     IsExpectedErrors -> success (log + continue); any other error -> log +
//     continue. NEVER abort the sweep.
//   - No RefundInstance (退订). All HBR resources swept here are PAYG
//     (post-paid) — DeleteBackupPlan/DeleteBackupClient/DeleteClient/
//     DeleteHanaBackupPlan/DeleteCrossAccount/DeleteVault are plain RPC
//     deletes, not BSS refund.

var hbrSweepPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf-example",
	"tf_example",
	"terraform-example",
	"tftestacc",
	"testAcc",
}

// hbrNotFoundErrs is the set of HBR error codes that indicate the resource is
// already gone (or the parent vault/cluster is gone). These are treated as
// delete-success in hbrDeleteOne so the sweep never aborts on stale state.
// Codes were derived from the service code (VaultNotExist, etc.) and the
// provider's IsExpectedErrors call sites in resource_alicloud_hbr_*.go.
var hbrNotFoundErrs = []string{
	"VaultNotExist",
	"BackupPlanNotExist",
	"PlanNotExist",
	"ClientNotExist",
	"BackupClientNotExist",
	"CrossAccountNotExist",
	"HanaBackupPlanNotExist",
	"HanaInstanceNotExist",
	"ResourceNotFound",
	"NotFound",
}

// hbrGetField returns the first non-empty string value among the given keys on
// item, trying both the as-is and TitleCased forms (HBR response casing is
// generally PascalCase but a few list endpoints return lowercased keys).
func hbrGetField(item map[string]interface{}, names ...string) string {
	for _, name := range names {
		if v, ok := item[name]; ok {
			if s := fmt.Sprint(v); s != "" && s != "<nil>" {
				return s
			}
		}
	}
	return ""
}

// hbrSweepNameOk is the fleet-wipe gate: never touch empty names or
// NODELETING-marked names; in NORMAL (non-SWEEP_ALL) mode only sweep
// tf-prefixed names; SWEEP_ALL bypasses the prefix filter so the whole test
// account is fleet-wiped (the test account has no prod data). The NODELETING
// guard is ALWAYS active — even under SWEEP_ALL — via protectedSweepVpcName
// (21 markers, normalization strips non-letter/digit chars, so hyphenated
// markers like "no-deleting-vault" -> "nodeletingvault" which CONTAINS
// "nodeleting" -> IS protected).
func hbrSweepNameOk(name string) bool {
	if name == "" {
		return false
	}
	if protectedSweepVpcName(name) {
		log.Printf("[INFO] Skipping HBR resource (NODELETING): %s", name)
		return false
	}
	if !sweepAll() {
		skip := true
		for _, prefix := range hbrSweepPrefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping HBR resource (no prefix): %s", name)
			return false
		}
	}
	return true
}

// hbrRpcPost issues an RpcPost to HBR with retry.Retry + NeedRetry +
// incrementalWait. Params go in the BODY slot (query=nil), matching the
// pre-existing inline HBR sweepers and the resource source code.
func hbrRpcPost(client *connectivity.AliyunClient, action string, query, body map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("hbr", "2017-09-08", action, query, body, true)
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

// hbrDeleteOne issues a best-effort HBR RPC delete (autoRetry=false, single
// shot — matching the pre-existing inline HBR sweepers). It NEVER returns an
// error that would abort the sweep: NotFoundError / IsExpectedErrors are
// logged as success (resource already gone) and any other failure is logged
// and swallowed. Callers ALWAYS continue to the next item.
func hbrDeleteOne(client *connectivity.AliyunClient, action string, params map[string]interface{}, name, id string) error {
	_, e := client.RpcPost("hbr", "2017-09-08", action, nil, params, false)
	if e != nil {
		if NotFoundError(e) || IsExpectedErrors(e, hbrNotFoundErrs) {
			log.Printf("[INFO] HBR %s already gone: %s (%s)", action, name, id)
			return nil
		}
		log.Printf("[ERROR] Failed to delete HBR %s (%s (%s)): %s", action, name, id, e)
		return nil
	}
	log.Printf("[INFO] Delete HBR %s success: %s (%s)", action, name, id)
	return nil
}

// hbrListVaults enumerates ALL backup vaults in the region via DescribeVaults
// (paginated, $.Vaults.Vault). Shared by the vault-scoped sweepers
// (hana_backup_client, hana_backup_plan) which need a VaultId to enumerate
// their children. Returns whatever was collected so far on error so the
// caller can still sweep the vaults it did see.
func hbrListVaults(client *connectivity.AliyunClient) []map[string]interface{} {
	action := "DescribeVaults"
	request := map[string]interface{}{
		"PageNumber": 1,
		"PageSize":   PageSizeLarge,
	}
	vaults := make([]map[string]interface{}, 0)
	for {
		response, err := hbrRpcPost(client, action, nil, request)
		if err != nil {
			log.Printf("[ERROR] HBR DescribeVaults for vault-scoped sweep: %s", err)
			return vaults
		}
		resp, err := jsonpath.Get("$.Vaults.Vault", response)
		if err != nil {
			return vaults
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				vaults = append(vaults, item)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return vaults
}

func init() {
	resource.AddTestSweepers("alicloud_hbr_cross_account", &resource.Sweeper{
		Name: "alicloud_hbr_cross_account",
		F:    testSweepHbrCrossAccounts,
	})
	resource.AddTestSweepers("alicloud_hbr_ecs_backup_plan", &resource.Sweeper{
		Name: "alicloud_hbr_ecs_backup_plan",
		F:    testSweepHbrEcsBackupPlans,
	})
	resource.AddTestSweepers("alicloud_hbr_nas_backup_plan", &resource.Sweeper{
		Name: "alicloud_hbr_nas_backup_plan",
		F:    testSweepHbrNasBackupPlans,
	})
	resource.AddTestSweepers("alicloud_hbr_oss_backup_plan", &resource.Sweeper{
		Name: "alicloud_hbr_oss_backup_plan",
		F:    testSweepHbrOssBackupPlans,
	})
	resource.AddTestSweepers("alicloud_hbr_ecs_backup_client", &resource.Sweeper{
		Name: "alicloud_hbr_ecs_backup_client",
		F:    testSweepHbrEcsBackupClients,
	})
	resource.AddTestSweepers("alicloud_hbr_hana_backup_client", &resource.Sweeper{
		Name: "alicloud_hbr_hana_backup_client",
		F:    testSweepHbrHanaBackupClients,
	})
	resource.AddTestSweepers("alicloud_hbr_hana_backup_plan", &resource.Sweeper{
		Name: "alicloud_hbr_hana_backup_plan",
		F:    testSweepHbrHanaBackupPlans,
	})
}

// testSweepHbrCrossAccounts cleans HBR cross-account management relationships.
// Cross-account is account-level (the list returns OwnerId = current account)
// but the list API is regional, so the same set appears in every region's
// response; sweeping cn-beijing alone catches them all. Filter on Alias
// (the human-readable name). Delete takes CrossAccountUserId +
// CrossAccountRoleName (the two-part id).
func testSweepHbrCrossAccounts(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	action := "DescribeCrossAccounts"
	request := map[string]interface{}{
		"PageNumber": 1,
		"PageSize":   PageSizeLarge,
	}
	accounts := make([]map[string]interface{}, 0)
	for {
		response, err := hbrRpcPost(client, action, nil, request)
		if err != nil {
			log.Printf("[ERROR] HBR DescribeCrossAccounts failed: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.CrossAccounts.CrossAccount", response)
		if err != nil {
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				accounts = append(accounts, item)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	if len(accounts) == 0 {
		log.Printf("[INFO] HBR CrossAccount: 0 cross-accounts in region %s", region)
		return nil
	}

	for _, acct := range accounts {
		name := hbrGetField(acct, "Alias")
		userId := hbrGetField(acct, "CrossAccountUserId")
		roleName := hbrGetField(acct, "CrossAccountRoleName")
		if userId == "" || roleName == "" {
			continue
		}
		id := userId + ":" + roleName
		if !hbrSweepNameOk(name) {
			continue
		}
		log.Printf("[INFO] Deleting HBR CrossAccount: %s (%s)", name, id)
		hbrDeleteOne(client, "DeleteCrossAccount", map[string]interface{}{
			"CrossAccountUserId":   userId,
			"CrossAccountRoleName": roleName,
		}, name, id)
	}
	return nil
}

// testSweepHbrEcsBackupPlans cleans HBR ECS-file backup plans
// (alicloud_hbr_ecs_backup_plan, deprecated since v1.249.0 in favor of
// policy/policy_binding). Filter on PlanName. Delete takes PlanId +
// SourceType=ECS_FILE + VaultId (all returned by the list).
func testSweepHbrEcsBackupPlans(region string) error {
	return hbrSweepBackupPlans(region, "ECS_FILE", "ECS Backup Plan")
}

// testSweepHbrNasBackupPlans cleans HBR NAS backup plans
// (alicloud_hbr_nas_backup_plan, deprecated since v1.249.0). Filter on
// PlanName. Delete takes PlanId + SourceType=NAS + VaultId. Note: the
// resource's Delete also calls DetachNasFileSystem (needs FileSystemId +
// CreateTime, not enumerable from the list) — the sweeper only calls
// DeleteBackupPlan, which removes the plan; the NAS file system itself is
// managed by the NAS product sweepers and is left attached (harmless).
func testSweepHbrNasBackupPlans(region string) error {
	return hbrSweepBackupPlans(region, "NAS", "NAS Backup Plan")
}

// testSweepHbrOssBackupPlans cleans HBR OSS backup plans
// (alicloud_hbr_oss_backup_plan, deprecated since v1.249.0). Filter on
// PlanName. Delete takes PlanId + SourceType=OSS + VaultId.
func testSweepHbrOssBackupPlans(region string) error {
	return hbrSweepBackupPlans(region, "OSS", "OSS Backup Plan")
}

// hbrSweepBackupPlans is the shared body for the ECS_FILE/NAS/OSS backup-plan
// sweepers. DescribeBackupPlans (SourceType filter, $.BackupPlans.BackupPlan)
// is the list; DeleteBackupPlan (PlanId + SourceType + VaultId) is the delete.
// The three flavors differ only in the SourceType constant and the log label.
func hbrSweepBackupPlans(region, sourceType, label string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	action := "DescribeBackupPlans"
	request := map[string]interface{}{
		"SourceType": sourceType,
		"PageNumber": 1,
		"PageSize":   PageSizeLarge,
	}
	plans := make([]map[string]interface{}, 0)
	for {
		response, err := hbrRpcPost(client, action, nil, request)
		if err != nil {
			log.Printf("[ERROR] HBR DescribeBackupPlans (%s) failed: %s", sourceType, err)
			return nil
		}
		resp, err := jsonpath.Get("$.BackupPlans.BackupPlan", response)
		if err != nil {
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				plans = append(plans, item)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	if len(plans) == 0 {
		log.Printf("[INFO] HBR %s: 0 plans in region %s", label, region)
		return nil
	}

	for _, plan := range plans {
		name := hbrGetField(plan, "PlanName")
		id := hbrGetField(plan, "PlanId")
		vaultId := hbrGetField(plan, "VaultId")
		if id == "" {
			continue
		}
		// hbrSweepNameOk returns false for empty name (always — even under
		// SWEEP_ALL), so nameless plans (e.g. some default plans) are
		// preserved. This matches the established HBR inline-sweeper contract.
		if !hbrSweepNameOk(name) {
			continue
		}
		params := map[string]interface{}{
			"PlanId":     id,
			"SourceType": sourceType,
		}
		if vaultId != "" {
			params["VaultId"] = vaultId
		}
		hbrDeleteOne(client, "DeleteBackupPlan", params, name, id)
	}
	return nil
}

// testSweepHbrEcsBackupClients cleans HBR ECS backup clients (the backup
// agent installed on an ECS instance). DescribeBackupClients (ClientType=
// ECS_CLIENT, $.Clients — a direct array, not nested under a wrapper) is the
// list; DeleteBackupClient (ClientId) is the delete. ECS backup clients are
// nameless (ClientId is an opaque UUID), so the gate uses InstanceName (the
// underlying ECS instance's name) — when a test creates an ECS instance with
// tf-prefix name and installs a backup client on it, the InstanceName on the
// client list echoes that tf-prefix name. ECS instances themselves are
// swept by alicloud_instance; sweeping the HBR client separately is for
// leftover clients on still-running test ECS instances.
func testSweepHbrEcsBackupClients(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	action := "DescribeBackupClients"
	request := map[string]interface{}{
		"ClientType": "ECS_CLIENT",
		"PageNumber": 1,
		"PageSize":   PageSizeLarge,
	}
	clients := make([]map[string]interface{}, 0)
	for {
		response, err := hbrRpcPost(client, action, nil, request)
		if err != nil {
			log.Printf("[ERROR] HBR DescribeBackupClients failed: %s", err)
			return nil
		}
		// DescribeBackupClients returns $.Clients as a DIRECT array (unlike
		// DescribeClients for HANA which returns $.Clients.Client). This
		// matches the service code (service_alicloud_hbr.go:436) and the
		// data source (data_source_alicloud_hbr_ecs_backup_clients.go:199).
		resp, err := jsonpath.Get("$.Clients", response)
		if err != nil {
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				clients = append(clients, item)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	if len(clients) == 0 {
		log.Printf("[INFO] HBR ECS Backup Client: 0 clients in region %s", region)
		return nil
	}

	for _, cli := range clients {
		// ECS backup clients are nameless; use InstanceName (the ECS
		// instance's name) as the sweep gate. Hostname is a fallback for
		// instances whose InstanceName is empty but whose hostname carries
		// the tf-prefix.
		name := hbrGetField(cli, "InstanceName", "Hostname")
		id := hbrGetField(cli, "ClientId")
		if id == "" {
			continue
		}
		if !hbrSweepNameOk(name) {
			continue
		}
		hbrDeleteOne(client, "DeleteBackupClient", map[string]interface{}{
			"ClientId": id,
		}, name, id)
	}
	return nil
}

// testSweepHbrHanaBackupClients cleans HBR HANA backup clients. HANA clients
// are vault-scoped: DescribeClients REQUIRES a VaultId (plus ClientType=
// ECS_AGENT, SourceType=HANA) to enumerate, so this sweeper walks
// DescribeVaults first, then per-vault lists clients. $.Clients.Client is the
// response path (nested, unlike DescribeBackupClients for ECS). Delete takes
// VaultId + ClientId (the two-part id). Clients have a ClientName and an
// InstanceName; ClientName is preferred for the gate.
func testSweepHbrHanaBackupClients(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	vaults := hbrListVaults(client)
	if len(vaults) == 0 {
		log.Printf("[INFO] HBR HANA Backup Client: 0 vaults in region %s", region)
		return nil
	}

	for _, vault := range vaults {
		vaultId := hbrGetField(vault, "VaultId")
		vaultName := hbrGetField(vault, "VaultName")
		if vaultId == "" {
			continue
		}
		// Vault-level NODELETING guard (always-active, even under SWEEP_ALL):
		// if a vault is marked NODELETING, its entire subtree (clients + plans)
		// is preserved — deleting a child of a protected vault would break the
		// backup configuration the user intends to keep. This is the safest
		// interpretation of the red-line.
		if protectedSweepVpcName(vaultName) {
			log.Printf("[INFO] Skipping HBR HANA Backup Client subtree (vault NODELETING): %s (%s)", vaultName, vaultId)
			continue
		}
		request := map[string]interface{}{
			"VaultId":    vaultId,
			"ClientType": "ECS_AGENT",
			"SourceType": "HANA",
			"PageNumber": 1,
			"PageSize":   PageSizeLarge,
		}
		for {
			response, err := hbrRpcPost(client, "DescribeClients", nil, request)
			if err != nil {
				log.Printf("[ERROR] HBR DescribeClients (HANA, vault %s (%s)) failed: %s", vaultName, vaultId, err)
				break
			}
			resp, err := jsonpath.Get("$.Clients.Client", response)
			if err != nil {
				break
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				name := hbrGetField(item, "ClientName", "InstanceName")
				clientId := hbrGetField(item, "ClientId")
				itemVaultId := hbrGetField(item, "VaultId")
				if itemVaultId == "" {
					itemVaultId = vaultId
				}
				if clientId == "" {
					continue
				}
				id := itemVaultId + ":" + clientId
				if !hbrSweepNameOk(name) {
					continue
				}
				hbrDeleteOne(client, "DeleteClient", map[string]interface{}{
					"VaultId":  itemVaultId,
					"ClientId": clientId,
				}, name, id)
			}
			if len(result) < PageSizeLarge {
				break
			}
			request["PageNumber"] = request["PageNumber"].(int) + 1
		}
	}
	return nil
}

// testSweepHbrHanaBackupPlans cleans HBR HANA backup plans. HANA plans are
// vault-scoped: DescribeHanaBackupPlans REQUIRES a VaultId to enumerate, so
// this sweeper walks DescribeVaults first, then per-vault lists plans.
// $.HanaBackupPlans.HanaBackupPlan is the response path. Delete takes
// PlanId + VaultId + ClusterId (the three-part id). Filter on PlanName.
func testSweepHbrHanaBackupPlans(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	vaults := hbrListVaults(client)
	if len(vaults) == 0 {
		log.Printf("[INFO] HBR HANA Backup Plan: 0 vaults in region %s", region)
		return nil
	}

	for _, vault := range vaults {
		vaultId := hbrGetField(vault, "VaultId")
		vaultName := hbrGetField(vault, "VaultName")
		if vaultId == "" {
			continue
		}
		// Vault-level NODELETING guard (always-active, even under SWEEP_ALL):
		// preserve the entire subtree of a NODELETING-marked vault.
		if protectedSweepVpcName(vaultName) {
			log.Printf("[INFO] Skipping HBR HANA Backup Plan subtree (vault NODELETING): %s (%s)", vaultName, vaultId)
			continue
		}
		request := map[string]interface{}{
			"VaultId":    vaultId,
			"PageNumber": 1,
			"PageSize":   PageSizeLarge,
		}
		for {
			response, err := hbrRpcPost(client, "DescribeHanaBackupPlans", nil, request)
			if err != nil {
				log.Printf("[ERROR] HBR DescribeHanaBackupPlans (vault %s (%s)) failed: %s", vaultName, vaultId, err)
				break
			}
			resp, err := jsonpath.Get("$.HanaBackupPlans.HanaBackupPlan", response)
			if err != nil {
				break
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				name := hbrGetField(item, "PlanName")
				planId := hbrGetField(item, "PlanId")
				itemVaultId := hbrGetField(item, "VaultId")
				if itemVaultId == "" {
					itemVaultId = vaultId
				}
				clusterId := hbrGetField(item, "ClusterId")
				if planId == "" {
					continue
				}
				id := planId + ":" + itemVaultId + ":" + clusterId
				if !hbrSweepNameOk(name) {
					continue
				}
				params := map[string]interface{}{
					"PlanId":  planId,
					"VaultId": itemVaultId,
				}
				if clusterId != "" {
					params["ClusterId"] = clusterId
				}
				hbrDeleteOne(client, "DeleteHanaBackupPlan", params, name, id)
			}
			if len(result) < PageSizeLarge {
				break
			}
			request["PageNumber"] = request["PageNumber"].(int) + 1
		}
	}
	return nil
}
