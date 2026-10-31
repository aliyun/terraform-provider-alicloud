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

// ALB (Application Load Balancer, product "Alb" / "2020-06-16") gap sweepers.
//
// API style note: the campaign brief describes ALB as ROA, but the actual
// provider resource code (resource_alicloud_alb_*.go) and the 7 pre-existing
// inline ALB sweepers all call client.RpcPost("Alb", "2020-06-16", action, nil,
// request, autoRetry) — ALB is RPC-style in this provider. These helpers match
// that proven convention (NOT RoaGet/RoaDelete), so they compile against the
// existing AliyunClient and behave like the inline ALB sweepers.
//
// Pre-existing inline ALB sweepers (do NOT re-register — AddTestSweepers
// log.Fatalf on duplicate names): alicloud_alb_acl,
// alicloud_alb_health_check_template, alicloud_alb_listener,
// alicloud_alb_load_balancer, alicloud_alb_rule, alicloud_alb_security_policy,
// alicloud_alb_server_group.
//
// Genuine gap sweeper added here:
//   - alicloud_alb_ascript: ListAScripts -> DeleteAScripts (id=AScriptId,
//     name=AScriptName). ListAScripts is a region-scoped list-all (ListenerIds /
//     AScriptIds / AScriptNames are all OPTIONAL per the 2020-06-16 API def), so
//     orphaned ascripts left behind when a listener is deleted can be found and
//     removed directly. Depends on alicloud_alb_listener so the parent sweeper
//     runs first: if the listener delete cascades its ascripts this sweeper is a
//     no-op; if not, it cleans the orphans that ListAScripts still returns.
//
// DEFERRED (no separate sweeper — cascade-cleaned by a parent sweeper, and no
// standalone list-all API exists; the only list endpoints are parent-scoped, or
// the attachment is a property read via GetListenerAttribute /
// GetLoadBalancerAttribute):
//   - alicloud_alb_acl_entry_attachment: entries are owned by the ACL;
//     ListAclEntries requires AclId (no list-all); the alicloud_alb_acl sweeper
//     deletes the ACL and cascades its entries.
//   - alicloud_alb_listener_acl_attachment: AclConfig is a listener property
//     (DescribeAlbListenerAclAttachment reads it via GetListenerAttribute);
//     there is no ListListenerAclAttachments API; the alicloud_alb_listener
//     sweeper deletes the listener and dissociates its ACLs.
//   - alicloud_alb_listener_additional_certificate_attachment: additional
//     certificates are a listener property (ListListenerCertificates requires
//     ListenerId; no list-all); the alicloud_alb_listener sweeper cascades.
//   - alicloud_alb_load_balancer_access_log_config_attachment: AccessLogConfig
//     is a LoadBalancer property (read via GetLoadBalancerAttribute); no
//     list-all; the alicloud_alb_load_balancer sweeper cascades (deleting the
//     LB disables its access log).
//   - alicloud_alb_load_balancer_common_bandwidth_package_attachment: the BWP
//     attachment is a LoadBalancer property; no list-all; the
//     alicloud_alb_load_balancer sweeper cascades (deleting the LB detaches the
//     BWP).
//   - alicloud_alb_load_balancer_security_group_attachment: SecurityGroupIds is
//     a LoadBalancer property; no list-all; the alicloud_alb_load_balancer
//     sweeper cascades.
//   - alicloud_alb_load_balancer_zone_shifted_attachment: the zone shift is a
//     LoadBalancer property; no list-all; the alicloud_alb_load_balancer
//     sweeper cascades (deleting the LB cancels the shift).

var albSweepPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf-example",
	"tf_example",
	"terraform-example",
	"tftestacc",
}

// albGetField returns the first non-empty value among the given keys on item,
// tolerating ALB's PascalCase response casing.
func albGetField(item map[string]interface{}, names ...string) string {
	for _, name := range names {
		if v, ok := item[name]; ok {
			if s := fmt.Sprint(v); s != "" && s != "<nil>" {
				return s
			}
		}
	}
	return ""
}

// albSweepNameOk is the fleet-wipe gate: never touch NODELETING names
// (always-active via protectedSweepVpcName); in normal (non-SWEEP_ALL) mode only
// sweep tf-prefixed names. SWEEP_ALL bypasses the prefix filter so the whole
// test account (1511928242963727, cn-beijing) is fleet-wiped.
func albSweepNameOk(name string) bool {
	if name == "" {
		return false
	}
	if protectedSweepVpcName(name) {
		log.Printf("[INFO] Skipping ALB resource (NODELETING): %s", name)
		return false
	}
	if !sweepAll() {
		skip := true
		for _, prefix := range albSweepPrefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping ALB resource (no prefix): %s", name)
			return false
		}
	}
	return true
}

// albRpcPost issues an RpcPost list/read with retry.Retry + NeedRetry +
// incrementalWait, matching the pre-existing inline ALB sweepers (RpcPost
// signature: product, version, action, query, body, autoRetry; autoRetry=true
// for reads). Params go in the BODY slot (query=nil).
func albRpcPost(client *connectivity.AliyunClient, action string, params map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, params, true)
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

// albDeletePost issues an RpcPost delete with retry.Retry + NeedRetry for
// transient errors only. autoRetry=true matches the alicloud_alb_ascript
// resource Delete (resource_alicloud_alb_ascript.go), which wraps the same call
// in retry.Retry and then checks NotFoundError on the returned error — so
// NotFoundError / IsExpectedErrors keep working on the retry-wrapped error
// returned here. DeleteAScripts is idempotent (already-gone returns
// ResourceNotFound.AScript), so auto-retry is safe. "Already gone" errors are
// classified by the caller; this helper never aborts the sweep.
func albDeletePost(client *connectivity.AliyunClient, action string, params map[string]interface{}) error {
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, err = client.RpcPost("Alb", "2020-06-16", action, nil, params, true)
		if err != nil {
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

func init() {
	// ascript depends on alicloud_alb_listener so the parent sweeper runs first:
	// if the listener delete cascades its ascripts this is a no-op; if not, this
	// sweeper cleans the orphans that ListAScripts (region-scoped, list-all)
	// still returns.
	resource.AddTestSweepers("alicloud_alb_ascript", &resource.Sweeper{
		Name:         "alicloud_alb_ascript",
		F:            testSweepAlbAScript,
		Dependencies: []string{"alicloud_alb_listener"},
	})
}

// testSweepAlbAScript: ListAScripts (region-scoped list-all; NextToken
// pagination, MaxResults=PageSizeXLarge) -> $.AScripts[].AScriptId; delete
// DeleteAScripts (AScriptIds.1 + ClientToken, matching the resource Delete).
// Best-effort: NotFoundError / ResourceNotFound.AScript are logged as "already
// gone"; all other errors are logged [ERROR]; the sweep always continues and
// returns nil (never aborts the whole run).
func testSweepAlbAScript(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"MaxResults": PageSizeXLarge,
	}
	for {
		response, e := albRpcPost(client, "ListAScripts", request)
		if e != nil {
			log.Printf("[ERROR] ListAScripts get an error: %#v", e)
			return nil
		}
		addDebug("ListAScripts", response, request)
		resp, jpErr := jsonpath.Get("$.AScripts", response)
		if formatInt(response["TotalCount"]) != 0 && jpErr != nil {
			log.Printf("[ERROR] Getting ALB AScripts attribute by path $.AScripts failed: %v. Body: %v.", jpErr, response)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			if item == nil {
				continue
			}
			id := albGetField(item, "AScriptId")
			name := albGetField(item, "AScriptName")
			// ascript name may be empty for some scripts; fall back to id for the
			// gate so a tf-named script is still swept even when AScriptName is
			// blank, but the NODELETING guard is never overridden.
			if !albSweepNameOk(name) && !albSweepNameOk(id) {
				continue
			}
			delReq := map[string]interface{}{
				"AScriptIds.1": id,
				"ClientToken":  buildClientToken("DeleteAScripts"),
			}
			if de := albDeletePost(client, "DeleteAScripts", delReq); de != nil {
				if NotFoundError(de) || IsExpectedErrors(de, []string{"ResourceNotFound.AScript"}) {
					log.Printf("[INFO] ALB AScript already gone: %s", id)
					continue
				}
				log.Printf("[ERROR] Failed to delete ALB AScript (%s): %s", id, de)
				continue
			}
			log.Printf("[INFO] Delete ALB AScript success: %s", id)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}
