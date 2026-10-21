package alicloud

// resource_alicloud_slb_sweeper_test.go fills the sweeper gaps for the SLB
// (Server Load Balancer, classic) product. SLB is an RPC-style API: every
// call below goes through client.RpcPost("Slb", "2014-05-15", action, nil,
// request, autoRetry), with the parameters carried in the request body slot
// (query=nil), matching the conventions of the existing inline SLB sweepers
// (alicloud_slb_acl / alicloud_slb_load_balancer / alicloud_slb_ca_certificate
// / alicloud_slb_server_certificate / alicloud_slb_tls_cipher_policy).
//
// Gap sweepers added (List -> Delete, id field):
//
//   alicloud_slb_listener                DescribeLoadBalancerListeners -> DeleteLoadBalancerListener  (ListenerPort + LoadBalancerId + ListenerProtocol; name gate = Description)
//   alicloud_slb_rule                    DescribeRules                 -> DeleteRules                   (RuleId; name gate = RuleName)
//   alicloud_slb_domain_extension        DescribeDomainExtensions      -> DeleteDomainExtension         (DomainExtensionId; gate = parent SLB name, see below)
//   alicloud_slb_server_group            DescribeVServerGroups         -> DeleteVServerGroup            (VServerGroupId; name gate = VServerGroupName)
//   alicloud_slb_master_slave_server_group DescribeMasterSlaveServerGroups -> DeleteMasterSlaveServerGroup  (MasterSlaveServerGroupId; name gate = MasterSlaveServerGroupName)
//
// DEFERRED resources (sub-resources that mutate a parent rather than being
// independently deletable; the parent sweeper already cascades their cleanup):
//
//   alicloud_slb_acl_entry_attachment        RemoveAccessControlListEntry mutates an ACL (composite id
//                                            aclId:entry). An entry has no name to tf-prefix gate on and
//                                            deleting one entry at a time is destructive; the existing
//                                            alicloud_slb_acl sweeper deletes the whole ACL which removes
//                                            its entries. Left as leftover when the ACL cannot be deleted.
//   alicloud_slb_attachment                  RemoveBackendServers mutates an SLB (id = LoadBalancerId). There
//                                            is no standalone list API and no name field; removing backend
//                                            ECS instances from a non-test SLB is unsafe. The
//                                            alicloud_slb_load_balancer sweeper deletes the whole SLB.
//   alicloud_slb_backend_server              Same as alicloud_slb_attachment: RemoveBackendServers with
//                                            id = LoadBalancerId; no name field; covered by the
//                                            load_balancer sweeper.
//   alicloud_slb_server_group_server_attachment RemoveVServerGroupBackendServers mutates a VServerGroup
//                                            (composite id vserverGroupId:serverId:port). No region-wide
//                                            list API and no name field; the alicloud_slb_server_group
//                                            sweeper deletes the whole group which clears its servers.
//
// Children-first DAG (resource.Sweeper.Dependencies runs BEFORE this sweeper):
//   alicloud_slb_rule                  -> []                            (child of listener)
//   alicloud_slb_domain_extension      -> []                            (child of listener)
//   alicloud_slb_listener              -> [alicloud_slb_rule, alicloud_slb_domain_extension]
//   alicloud_slb_server_group          -> [alicloud_slb_listener]       (listener may bind a group)
//   alicloud_slb_master_slave_server_group -> [alicloud_slb_listener]
//
// All deletes are best-effort: NotFound / not-found error codes are logged as
// "[INFO] already gone" and the sweep continues; any other error is logged as
// "[ERROR]" and the sweep continues. The function always returns nil.
//
// SLB instances may carry DeletionProtection (DeleteProtection=on). These
// sweepers never toggle DeleteProtection on a parent SLB to force-delete a
// child; a delete that fails because the SLB is protected is logged and the
// resource is left as a leftover (recorded in the sweep log).

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

// slbSweepNameOk is the name gate shared by every SLB gap sweeper. It honours
// the always-active NODELETING protection (protectedSweepVpcName), the
// sweepAll() flag (sweep everything not protected) and the tf-testAcc /
// tf_testAcc prefix gate used in the default (non-sweepAll) mode.
func slbSweepNameOk(name string) bool {
	if protectedSweepVpcName(name) {
		return false
	}
	if sweepAll() {
		return true
	}
	for _, prefix := range []string{"tf-testAcc", "tf_testAcc"} {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			return true
		}
	}
	return false
}

// slbGetField reads a value from a List-API item case-insensitively and
// returns its string form. Returns "" when the field is absent. SLB list
// responses use stable PascalCase keys, but the case-insensitive lookup keeps
// the helper robust against per-API casing drift.
func slbGetField(item map[string]interface{}, field string) string {
	if v, ok := item[field]; ok {
		return fmt.Sprint(v)
	}
	lfield := strings.ToLower(field)
	for k, v := range item {
		if strings.ToLower(k) == lfield {
			return fmt.Sprint(v)
		}
	}
	return ""
}

// slbRpcPost is a thin retry wrapper around the SLB RPC POST API
// (product "Slb", version "2014-05-15"). SLB carries all parameters in the
// request body slot (query=nil), matching the existing inline SLB sweepers.
func slbRpcPost(client *connectivity.AliyunClient, action string, request map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(3*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, autoRetry)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	return response, err
}

// slbSweepDelete issues a best-effort delete through slbRpcPost. NotFound and
// the resource-specific not-found error codes are treated as success (already
// gone); every other error is returned for the caller to log and continue.
func slbSweepDelete(client *connectivity.AliyunClient, action string, request map[string]interface{}, notFound []string) error {
	_, err := slbRpcPost(client, action, request, false)
	if err != nil {
		if NotFoundError(err) || IsExpectedErrors(err, notFound) {
			return nil
		}
		return err
	}
	return nil
}

// slbListLoadBalancers paginates DescribeLoadBalancers (PageNumber/PageSize)
// and returns every SLB item (LoadBalancerId, LoadBalancerName, ...) in the
// region. Used both to gate children by their parent SLB name and to skip
// children whose parent SLB is NODELETING-protected.
func slbListLoadBalancers(client *connectivity.AliyunClient) ([]map[string]interface{}, error) {
	out := make([]map[string]interface{}, 0)
	action := "DescribeLoadBalancers"
	page := 1
	for {
		request := map[string]interface{}{
			"RegionId":   client.RegionId,
			"PageSize":   PageSizeLarge,
			"PageNumber": page,
		}
		response, err := slbRpcPost(client, action, request, true)
		if err != nil {
			return nil, err
		}
		resp, err := jsonpath.Get("$.LoadBalancers.LoadBalancer", response)
		if err != nil {
			return nil, WrapErrorf(err, FailedGetAttributeMsg, action, "$.LoadBalancers.LoadBalancer", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				out = append(out, item)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		page++
	}
	return out, nil
}

// slbListListeners paginates DescribeLoadBalancerListeners (NextToken) for the
// whole region. Each returned listener carries LoadBalancerId, ListenerPort,
// ListenerProtocol and Description. The list is used by the listener sweeper
// directly and by the rule / domain_extension sweepers to enumerate the
// (LoadBalancerId, ListenerPort, ListenerProtocol) tuples they need to list
// children under.
func slbListListeners(client *connectivity.AliyunClient) ([]map[string]interface{}, error) {
	out := make([]map[string]interface{}, 0)
	action := "DescribeLoadBalancerListeners"
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"MaxResults": PageSizeXLarge,
	}
	for {
		response, err := slbRpcPost(client, action, request, true)
		if err != nil {
			return nil, err
		}
		// DescribeLoadBalancerListeners returns the item array at the top
		// level ($.Listeners), not nested under a second key.
		resp, err := jsonpath.Get("$.Listeners", response)
		if err != nil {
			// DescribeLoadBalancerListeners omits the Listeners key entirely
			// when TotalCount is 0 (unlike DescribeLoadBalancers, which
			// returns an empty array). Treat the missing key as an empty
			// list so the sweeper continues with zero items instead of
			// aborting the whole sweep on a legitimately empty region.
			log.Printf("[INFO] slb %s: no Listeners field in response (TotalCount=%v): %s", action, response["TotalCount"], err)
			return out, nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				out = append(out, item)
			}
		}
		nextToken, _ := response["NextToken"].(string)
		if nextToken == "" {
			break
		}
		request["NextToken"] = nextToken
	}
	return out, nil
}

// slbListenerPort parses the ListenerPort field (a JSON number surfaced as
// float64) into the int the SLB RPC API expects. Returns 0 on a parse failure,
// which the caller treats as "skip this listener".
func slbListenerPort(item map[string]interface{}) int {
	port, _ := strconv.Atoi(slbGetField(item, "ListenerPort"))
	return port
}

// slbProtectedLoadBalancers returns the set of LoadBalancerId values whose name
// carries a NODELETING marker. Children discovered under such an SLB are
// skipped wholesale so a do-not-delete SLB is never mutated by a child sweep.
func slbProtectedLoadBalancers(slbs []map[string]interface{}) map[string]bool {
	protected := make(map[string]bool)
	for _, lb := range slbs {
		if protectedSweepVpcName(slbGetField(lb, "LoadBalancerName")) {
			protected[slbGetField(lb, "LoadBalancerId")] = true
		}
	}
	return protected
}

func init() {
	resource.AddTestSweepers("alicloud_slb_rule", &resource.Sweeper{
		Name:         "alicloud_slb_rule",
		F:            testSweepSlbRule,
		Dependencies: []string{},
	})
	resource.AddTestSweepers("alicloud_slb_domain_extension", &resource.Sweeper{
		Name:         "alicloud_slb_domain_extension",
		F:            testSweepSlbDomainExtension,
		Dependencies: []string{},
	})
	resource.AddTestSweepers("alicloud_slb_listener", &resource.Sweeper{
		Name:         "alicloud_slb_listener",
		F:            testSweepSlbListener,
		Dependencies: []string{"alicloud_slb_rule", "alicloud_slb_domain_extension"},
	})
	resource.AddTestSweepers("alicloud_slb_server_group", &resource.Sweeper{
		Name:         "alicloud_slb_server_group",
		F:            testSweepSlbServerGroup,
		Dependencies: []string{"alicloud_slb_listener"},
	})
	resource.AddTestSweepers("alicloud_slb_master_slave_server_group", &resource.Sweeper{
		Name:         "alicloud_slb_master_slave_server_group",
		F:            testSweepSlbMasterSlaveServerGroup,
		Dependencies: []string{"alicloud_slb_listener"},
	})
}

// testSweepSlbListener deletes tf-prefixed (or, in sweepAll mode, all
// non-NODELETING) listeners across the region. Listeners are enumerated
// region-wide via DescribeLoadBalancerListeners and gated by their Description
// field, which the acceptance tests populate with a tf-testAcc name.
func testSweepSlbListener(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	slbs, err := slbListLoadBalancers(client)
	if err != nil {
		log.Printf("[ERROR] slb listener: failed to list load balancers: %s", err)
		return nil
	}
	protected := slbProtectedLoadBalancers(slbs)

	listeners, err := slbListListeners(client)
	if err != nil {
		log.Printf("[ERROR] slb listener: failed to list listeners: %s", err)
		return nil
	}
	for _, l := range listeners {
		lbId := slbGetField(l, "LoadBalancerId")
		if protected[lbId] {
			log.Printf("[INFO] Skipping SLB listener on protected SLB %s: %s", lbId, slbGetField(l, "Description"))
			continue
		}
		name := slbGetField(l, "Description")
		if !slbSweepNameOk(name) {
			log.Printf("[INFO] Skipping SLB listener: %s", name)
			continue
		}
		port := slbListenerPort(l)
		if port == 0 {
			log.Printf("[WARN] Skipping SLB listener with unparseable port: %v", l)
			continue
		}
		protocol := slbGetField(l, "ListenerProtocol")
		log.Printf("[INFO] Deleting SLB listener: %s (lb=%s port=%d proto=%s)", name, lbId, port, protocol)
		req := map[string]interface{}{
			"RegionId":         client.RegionId,
			"LoadBalancerId":   lbId,
			"ListenerPort":     port,
			"ListenerProtocol": protocol,
		}
		if err := slbSweepDelete(client, "DeleteLoadBalancerListener", req, []string{"ResourceNotFound.Listener"}); err != nil {
			log.Printf("[ERROR] Failed to delete SLB listener (%s): %s", name, err)
		}
	}
	return nil
}

// testSweepSlbRule deletes tf-prefixed (or sweepAll) forwarding rules. Rules
// are listed per (LoadBalancerId, ListenerPort) via DescribeRules, so the
// sweeper first enumerates every listener in the region and then lists the
// rules under each. A rule is gated by its RuleName.
func testSweepSlbRule(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	slbs, err := slbListLoadBalancers(client)
	if err != nil {
		log.Printf("[ERROR] slb rule: failed to list load balancers: %s", err)
		return nil
	}
	protected := slbProtectedLoadBalancers(slbs)

	listeners, err := slbListListeners(client)
	if err != nil {
		log.Printf("[ERROR] slb rule: failed to list listeners: %s", err)
		return nil
	}
	for _, l := range listeners {
		lbId := slbGetField(l, "LoadBalancerId")
		if protected[lbId] {
			continue
		}
		port := slbListenerPort(l)
		if port == 0 {
			continue
		}
		protocol := slbGetField(l, "ListenerProtocol")
		listReq := map[string]interface{}{
			"RegionId":         client.RegionId,
			"LoadBalancerId":   lbId,
			"ListenerPort":     port,
			"ListenerProtocol": protocol,
		}
		response, err := slbRpcPost(client, "DescribeRules", listReq, true)
		if err != nil {
			log.Printf("[ERROR] slb rule: DescribeRules failed (lb=%s port=%d): %s", lbId, port, err)
			continue
		}
		resp, err := jsonpath.Get("$.Rules.Rule", response)
		if err != nil {
			continue
		}
		rules, _ := resp.([]interface{})
		for _, r := range rules {
			rule, ok := r.(map[string]interface{})
			if !ok {
				continue
			}
			name := slbGetField(rule, "RuleName")
			ruleId := slbGetField(rule, "RuleId")
			if !slbSweepNameOk(name) {
				log.Printf("[INFO] Skipping SLB rule: %s (%s)", name, ruleId)
				continue
			}
			log.Printf("[INFO] Deleting SLB rule: %s (%s)", name, ruleId)
			// RuleIds uses the single-quoted list form the resource delete
			// function uses (accepted by the SLB API).
			delReq := map[string]interface{}{
				"RegionId": client.RegionId,
				"RuleIds":  fmt.Sprintf("['%s']", ruleId),
			}
			if err := slbSweepDelete(client, "DeleteRules", delReq, []string{"InvalidRuleId.NotFound"}); err != nil {
				log.Printf("[ERROR] Failed to delete SLB rule (%s (%s)): %s", name, ruleId, err)
			}
		}
	}
	return nil
}

// testSweepSlbDomainExtension deletes extension domains. DescribeDomainExtensions
// is scoped to (LoadBalancerId, ListenerPort), and the Domain field itself is
// not tf-prefixed by the acceptance tests (they use "www.test.com"), so the
// sweeper gates by the parent SLB name instead: every extension domain under
// a tf-prefixed (or sweepAll, non-NODELETING) SLB is removed.
func testSweepSlbDomainExtension(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	// The listener list only carries LoadBalancerId, so build an id->name map
	// once from the SLB list and use it both for the tf-prefix / sweepAll gate
	// and for NODELETING protection (slbSweepNameOk checks both).
	slbs, err := slbListLoadBalancers(client)
	if err != nil {
		log.Printf("[ERROR] slb domain extension: failed to list load balancers: %s", err)
		return nil
	}
	lbNames := make(map[string]string, len(slbs))
	for _, lb := range slbs {
		lbNames[slbGetField(lb, "LoadBalancerId")] = slbGetField(lb, "LoadBalancerName")
	}

	listeners, err := slbListListeners(client)
	if err != nil {
		log.Printf("[ERROR] slb domain extension: failed to list listeners: %s", err)
		return nil
	}
	for _, l := range listeners {
		// Domain extensions only exist under HTTPS listeners.
		if !strings.EqualFold(slbGetField(l, "ListenerProtocol"), "https") {
			continue
		}
		lbId := slbGetField(l, "LoadBalancerId")
		lbName := lbNames[lbId]
		if !slbSweepNameOk(lbName) {
			log.Printf("[INFO] Skipping SLB domain extensions on non-test SLB: %s (%s)", lbName, lbId)
			continue
		}
		port := slbListenerPort(l)
		if port == 0 {
			continue
		}
		listReq := map[string]interface{}{
			"RegionId":       client.RegionId,
			"LoadBalancerId": lbId,
			"ListenerPort":   port,
		}
		response, err := slbRpcPost(client, "DescribeDomainExtensions", listReq, true)
		if err != nil {
			log.Printf("[ERROR] slb domain extension: DescribeDomainExtensions failed (lb=%s port=%d): %s", lbId, port, err)
			continue
		}
		resp, err := jsonpath.Get("$.DomainExtensions.DomainExtension", response)
		if err != nil {
			continue
		}
		exts, _ := resp.([]interface{})
		for _, e := range exts {
			ext, ok := e.(map[string]interface{})
			if !ok {
				continue
			}
			extId := slbGetField(ext, "DomainExtensionId")
			domain := slbGetField(ext, "Domain")
			log.Printf("[INFO] Deleting SLB domain extension: %s (%s)", domain, extId)
			delReq := map[string]interface{}{
				"RegionId":          client.RegionId,
				"DomainExtensionId": extId,
			}
			if err := slbSweepDelete(client, "DeleteDomainExtension", delReq, []string{"InvalidParameter.DomainExtensionId"}); err != nil {
				log.Printf("[ERROR] Failed to delete SLB domain extension (%s (%s)): %s", domain, extId, err)
			}
		}
	}
	return nil
}

// testSweepSlbServerGroup deletes tf-prefixed (or sweepAll) VServer groups.
// DescribeVServerGroups is scoped to a LoadBalancerId, so the sweeper enumerates
// every SLB and lists the groups under each. A group is gated by its
// VServerGroupName.
func testSweepSlbServerGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	slbs, err := slbListLoadBalancers(client)
	if err != nil {
		log.Printf("[ERROR] slb server group: failed to list load balancers: %s", err)
		return nil
	}
	for _, lb := range slbs {
		lbId := slbGetField(lb, "LoadBalancerId")
		lbName := slbGetField(lb, "LoadBalancerName")
		if protectedSweepVpcName(lbName) {
			log.Printf("[INFO] Skipping SLB server groups on protected SLB: %s (%s)", lbName, lbId)
			continue
		}
		listReq := map[string]interface{}{
			"RegionId":       client.RegionId,
			"LoadBalancerId": lbId,
		}
		response, err := slbRpcPost(client, "DescribeVServerGroups", listReq, true)
		if err != nil {
			log.Printf("[ERROR] slb server group: DescribeVServerGroups failed (lb=%s): %s", lbId, err)
			continue
		}
		resp, err := jsonpath.Get("$.VServerGroups.VServerGroup", response)
		if err != nil {
			continue
		}
		groups, _ := resp.([]interface{})
		for _, g := range groups {
			group, ok := g.(map[string]interface{})
			if !ok {
				continue
			}
			name := slbGetField(group, "VServerGroupName")
			id := slbGetField(group, "VServerGroupId")
			if !slbSweepNameOk(name) {
				log.Printf("[INFO] Skipping SLB server group: %s (%s)", name, id)
				continue
			}
			log.Printf("[INFO] Deleting SLB server group: %s (%s)", name, id)
			delReq := map[string]interface{}{
				"RegionId":       client.RegionId,
				"VServerGroupId": id,
			}
			if err := slbSweepDelete(client, "DeleteVServerGroup", delReq, []string{"The specified VServerGroupId does not exist", "InvalidParameter"}); err != nil {
				log.Printf("[ERROR] Failed to delete SLB server group (%s (%s)): %s", name, id, err)
			}
		}
	}
	return nil
}

// testSweepSlbMasterSlaveServerGroup deletes tf-prefixed (or sweepAll)
// master-slave server groups. DescribeMasterSlaveServerGroups is scoped to a
// LoadBalancerId, so the sweeper enumerates every SLB and lists the groups
// under each. A group is gated by its MasterSlaveServerGroupName.
func testSweepSlbMasterSlaveServerGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	slbs, err := slbListLoadBalancers(client)
	if err != nil {
		log.Printf("[ERROR] slb master-slave server group: failed to list load balancers: %s", err)
		return nil
	}
	for _, lb := range slbs {
		lbId := slbGetField(lb, "LoadBalancerId")
		lbName := slbGetField(lb, "LoadBalancerName")
		if protectedSweepVpcName(lbName) {
			log.Printf("[INFO] Skipping SLB master-slave groups on protected SLB: %s (%s)", lbName, lbId)
			continue
		}
		listReq := map[string]interface{}{
			"RegionId":       client.RegionId,
			"LoadBalancerId": lbId,
		}
		response, err := slbRpcPost(client, "DescribeMasterSlaveServerGroups", listReq, true)
		if err != nil {
			log.Printf("[ERROR] slb master-slave group: DescribeMasterSlaveServerGroups failed (lb=%s): %s", lbId, err)
			continue
		}
		resp, err := jsonpath.Get("$.MasterSlaveServerGroups.MasterSlaveServerGroup", response)
		if err != nil {
			continue
		}
		groups, _ := resp.([]interface{})
		for _, g := range groups {
			group, ok := g.(map[string]interface{})
			if !ok {
				continue
			}
			name := slbGetField(group, "MasterSlaveServerGroupName")
			id := slbGetField(group, "MasterSlaveServerGroupId")
			if !slbSweepNameOk(name) {
				log.Printf("[INFO] Skipping SLB master-slave server group: %s (%s)", name, id)
				continue
			}
			log.Printf("[INFO] Deleting SLB master-slave server group: %s (%s)", name, id)
			delReq := map[string]interface{}{
				"RegionId":                 client.RegionId,
				"MasterSlaveServerGroupId": id,
			}
			if err := slbSweepDelete(client, "DeleteMasterSlaveServerGroup", delReq, []string{"The specified MasterSlaveGroupId does not exist", "InvalidParameter"}); err != nil {
				log.Printf("[ERROR] Failed to delete SLB master-slave server group (%s (%s)): %s", name, id, err)
			}
		}
	}
	return nil
}
