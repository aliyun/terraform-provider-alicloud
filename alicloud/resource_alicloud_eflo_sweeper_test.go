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

// EFLO (灵骏 / Intelligent Compute) sweeper gap fill. The eflo surface spans
// three RPC products with separate versions, each used directly by the
// provider's resource files:
//
//   - "eflo"           / "2022-05-30" : er, subnet, vpd, vpd_grant_rule
//   - "eflo-controller" / "2022-12-15": cluster, node_group, vsc (+ node,
//     hyper_node, invocation, node_group_attachment — see DEFERRED below)
//   - "eflo-cnp"        / "2023-08-28": experiment_plan, experiment_plan_template
//
// None of the 14 alicloud_eflo_* resources had a registered sweeper before this
// file. This file fills the 9 that have BOTH a region-scoped list API and a
// clean (non-退订) delete API:
//
//   - alicloud_eflo_subnet               (ListSubnets      -> DeleteSubnet,      by SubnetId;        PageNumber/PageSize,       $.Content.Data)
//   - alicloud_eflo_vpd                  (ListVpds         -> DeleteVpd,         by VpdId;            PageNumber/PageSize,       $.Content.Data)
//   - alicloud_eflo_er                    (ListErs          -> DeleteEr,          by ErId;             PageNumber/PageSize,       $.Content.Data)
//   - alicloud_eflo_vpd_grant_rule        (ListVpdGrantRules-> DeleteVpdGrantRule, by GrantRuleId;     PageNumber/PageSize,       $.Content.Data)
//   - alicloud_eflo_cluster               (ListClusters     -> DeleteCluster,     by ClusterId;        NextToken/MaxResults,      $.Clusters)
//   - alicloud_eflo_vsc                   (ListVscs         -> DeleteVsc,         by VscId;            NextToken/MaxResults,      $.Vscs)
//   - alicloud_eflo_node_group            (parent-scoped: ListClusters -> ListNodeGroups per cluster -> DeleteNodeGroup, by ClusterId:NodeGroupId; NextToken/MaxResults, $.Groups)
//   - alicloud_eflo_experiment_plan       (ListExperimentPlans -> DeleteExperimentPlan, by PlanId;       Page/Size,                 $.Data)
//   - alicloud_eflo_experiment_plan_template (ListExperimentPlanTemplates -> DeleteExperimentPlanTemplate, by TemplateId; single call, $.Data; only PrivacyLevel=private is swept so shared/public templates are spared)
//
// Intentionally DEFERRED (no sweep written):
//   - alicloud_eflo_hyper_node: resource Delete is BssOpenApi RefundInstance only
//     (退订-class, requires explicit human approval). There is no clean
//     eflo-controller release; the only delete path is a BSS refund.
//   - alicloud_eflo_node: resource Delete is hybrid — Subscription/PREPAY nodes
//     are released via BssOpenApi RefundInstance (退订-class), while
//     PayAsYouGo/POSTPAY nodes use the eflo-controller DeleteNode. eflo node
//     test fixtures use payment_type "Subscription", so the only effective
//     release path is 退订. Deferred per policy to avoid a partial/ineffective
//     sweep (DeleteNode would no-op subscription nodes) and the 退订 approval
//     gate.
//   - alicloud_eflo_invocation: resource Delete is a no-op ("Cannot destroy ...
//     resources may remain"). A RunCommand invocation is ephemeral; there is no
//     delete API and no list-and-delete semantics.
//   - alicloud_eflo_node_group_attachment: resource Delete operates via the
//     cluster ShrinkCluster action (not a dedicated delete API) and there is no
//     list endpoint for attachments; it is a derived attachment managed through
//     its parent cluster/node-group.
//   - alicloud_eflo_resource: resource Delete is a no-op ("Cannot destroy ...
//     resources may remain"); it is a cluster workspace/resource registration
//     (CreateResource/ValidateResource) with no delete API.
//
// Fleet-wipe pattern mirrors the KMS sweeper: SWEEP_ALL
// (ALICLOUD_SWEEP_ALL_RESOURCES=true) bypasses the tf-prefix filter so every
// non-NODELETING resource is swept; NORMAL mode keeps the prefix filter. The
// protectedSweepVpcName NODELETING guard is always-active in both modes. All
// sweepers return nil — list/delete errors are logged and never propagated, so
// one failing item does not abort the cascade.

func init() {
	resource.AddTestSweepers("alicloud_eflo_subnet", &resource.Sweeper{
		Name: "alicloud_eflo_subnet",
		F:    testSweepEfloSubnets,
	})
	resource.AddTestSweepers("alicloud_eflo_vpd", &resource.Sweeper{
		Name: "alicloud_eflo_vpd",
		F:    testSweepEfloVpds,
		// Subnets and grant rules reference a VPD; sweep them first so
		// DeleteVpd is not blocked by child subnets/grants.
		Dependencies: []string{
			"alicloud_eflo_subnet",
			"alicloud_eflo_vpd_grant_rule",
		},
	})
	resource.AddTestSweepers("alicloud_eflo_vpd_grant_rule", &resource.Sweeper{
		Name: "alicloud_eflo_vpd_grant_rule",
		F:    testSweepEfloVpdGrantRules,
	})
	resource.AddTestSweepers("alicloud_eflo_er", &resource.Sweeper{
		Name: "alicloud_eflo_er",
		F:    testSweepEfloErs,
		// Grant rules authorize a VPD/VCC instance INTO an ER; sweep grants
		// before the ER so DeleteEr is not blocked.
		Dependencies: []string{
			"alicloud_eflo_vpd_grant_rule",
		},
	})
	resource.AddTestSweepers("alicloud_eflo_cluster", &resource.Sweeper{
		Name: "alicloud_eflo_cluster",
		F:    testSweepEfloClusters,
		// Node groups live inside a cluster; sweep them first so
		// DeleteCluster is not blocked by child node groups.
		Dependencies: []string{
			"alicloud_eflo_node_group",
		},
	})
	resource.AddTestSweepers("alicloud_eflo_node_group", &resource.Sweeper{
		Name: "alicloud_eflo_node_group",
		F:    testSweepEfloNodeGroups,
	})
	resource.AddTestSweepers("alicloud_eflo_vsc", &resource.Sweeper{
		Name: "alicloud_eflo_vsc",
		F:    testSweepEfloVscs,
	})
	resource.AddTestSweepers("alicloud_eflo_experiment_plan", &resource.Sweeper{
		Name: "alicloud_eflo_experiment_plan",
		F:    testSweepEfloExperimentPlans,
	})
	resource.AddTestSweepers("alicloud_eflo_experiment_plan_template", &resource.Sweeper{
		Name: "alicloud_eflo_experiment_plan_template",
		F:    testSweepEfloExperimentPlanTemplates,
		// Plans reference a template; sweep plans before the template so
		// DeleteExperimentPlanTemplate is not blocked by referencing plans.
		Dependencies: []string{
			"alicloud_eflo_experiment_plan",
		},
	})
}

var efloSweeperPrefixes = []string{
	"tf-testacc",
	"tf_testacc",
	"tf-testAcc",
	"tf_testAcc",
	"tf-example",
	"tf_example",
	"terraform-example",
	"tftestacc",
}

// efloGetField returns the first non-empty string value of item[name] across
// the given name variants. eflo response field casing is mostly PascalCase; the
// variant list keeps the helper tolerant of camelCase fallbacks. Integer IDs
// (PlanId/TemplateId) are stringified via fmt.Sprint, which renders whole
// float64/int64 values without a trailing decimal — matching how the resource
// Delete passes d.Id() (a string) for the same fields.
func efloGetField(item map[string]interface{}, names ...string) string {
	for _, name := range names {
		if v, ok := item[name]; ok {
			if s := fmt.Sprint(v); s != "" && s != "<nil>" {
				return s
			}
		}
	}
	return ""
}

// efloSweepNameOk applies the always-active NODELETING guard and the
// sweepAll()-aware tf-prefix filter. Returns true when the resource should be
// swept. SWEEP_ALL bypasses the prefix filter (fleet-wipe: clean all
// non-NODELETING); NORMAL mode keeps the prefix filter.
func efloSweepNameOk(name string) bool {
	if name == "" {
		return false
	}
	if protectedSweepVpcName(name) {
		log.Printf("[INFO] Skipping eflo resource (NODELETING): %s", name)
		return false
	}
	if !sweepAll() {
		skip := true
		for _, prefix := range efloSweeperPrefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping eflo resource (no prefix): %s", name)
			return false
		}
	}
	return true
}

// efloRpcPost is the shared eflo list/read wrapper: RpcPost with autoRetry=false
// (nil query; all params in the request body, mirroring how the provider's own
// eflo resource/service calls populate request and leave query empty) wrapped in
// retry.Retry + NeedRetry + incrementalWait so throttling and 5xx are retried at
// the sweeper level. product/version are parameters because the eflo surface
// spans three products (eflo / eflo-controller / eflo-cnp).
func efloRpcPost(client *connectivity.AliyunClient, product, version, action string, request map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	var err error
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost(product, version, action, nil, request, false)
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
		return nil, err
	}
	return response, nil
}

// efloDelete issues a best-effort delete (single-shot, autoRetry=false) and
// swallows every error — NotFound, expected not-found codes, and transient
// failures are all logged so one item never aborts the cascade. The sweeper
// runs periodically; transient failures clear on the next pass. Mirrors the KMS
// sweeper delete style.
func efloDelete(client *connectivity.AliyunClient, product, version, action string, request map[string]interface{}, kind, id string) {
	log.Printf("[INFO] Deleting eflo %s: %s", kind, id)
	if _, e := client.RpcPost(product, version, action, nil, request, false); e != nil {
		log.Printf("[ERROR] Failed to delete eflo %s (%s): %s", kind, id, e)
		return
	}
	log.Printf("[INFO] Delete eflo %s success: %s", kind, id)
}

// efloNextToken extracts the NextToken paging cursor from an eflo-controller
// response (top-level field). Returns "" when absent/empty, which for these
// APIs means the page stream is exhausted.
func efloNextToken(response map[string]interface{}) string {
	return efloGetField(response, "NextToken")
}

// ---------------------------------------------------------------------------
// eflo / 2022-05-30 sweepers (PageNumber/PageSize, $.Content.Data)
// ---------------------------------------------------------------------------

func testSweepEfloSubnets(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	// Mirrors the alicloud_eflo_subnets datasource: RegionId + PageNumber/PageSize.
	// ListSubnets paginates by default (EnablePage omitted, matching the datasource).
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageNumber": 1,
		"PageSize":   PageSizeLarge,
	}
	action := "ListSubnets"
	for {
		response, err := efloRpcPost(client, "eflo", "2022-05-30", action, request)
		if err != nil {
			log.Printf("[ERROR] Failed to list eflo subnets: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.Content.Data", response)
		if err != nil {
			log.Printf("[ERROR] ListSubnets jsonpath $.Content.Data: %s", err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			name := efloGetField(item, "SubnetName")
			if !efloSweepNameOk(name) {
				continue
			}
			vpdId := efloGetField(item, "VpdId")
			subnetId := efloGetField(item, "SubnetId")
			zoneId := efloGetField(item, "ZoneId")
			if vpdId == "" || subnetId == "" {
				continue
			}
			delReq := map[string]interface{}{
				"VpdId":    vpdId,
				"SubnetId": subnetId,
				"ZoneId":   zoneId,
			}
			efloDelete(client, "eflo", "2022-05-30", "DeleteSubnet", delReq, "subnet", vpdId+":"+subnetId)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func testSweepEfloVpds(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	// Mirrors the alicloud_eflo_vpds datasource: EnablePage=true + PageNumber/PageSize.
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"EnablePage": true,
		"PageNumber": 1,
		"PageSize":   PageSizeLarge,
	}
	action := "ListVpds"
	for {
		response, err := efloRpcPost(client, "eflo", "2022-05-30", action, request)
		if err != nil {
			log.Printf("[ERROR] Failed to list eflo vpds: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.Content.Data", response)
		if err != nil {
			log.Printf("[ERROR] ListVpds jsonpath $.Content.Data: %s", err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			name := efloGetField(item, "VpdName")
			if !efloSweepNameOk(name) {
				continue
			}
			vpdId := efloGetField(item, "VpdId")
			if vpdId == "" {
				continue
			}
			delReq := map[string]interface{}{
				"VpdId":    vpdId,
				"RegionId": client.RegionId,
			}
			efloDelete(client, "eflo", "2022-05-30", "DeleteVpd", delReq, "vpd", vpdId)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func testSweepEfloErs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	// ListErs: RegionId required; EnablePage documented (default true) and set
	// explicitly like the vpds datasource. PageNumber/PageSize pagination.
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"EnablePage": true,
		"PageNumber": 1,
		"PageSize":   PageSizeLarge,
	}
	action := "ListErs"
	for {
		response, err := efloRpcPost(client, "eflo", "2022-05-30", action, request)
		if err != nil {
			log.Printf("[ERROR] Failed to list eflo ers: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.Content.Data", response)
		if err != nil {
			log.Printf("[ERROR] ListErs jsonpath $.Content.Data: %s", err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			name := efloGetField(item, "ErName")
			if !efloSweepNameOk(name) {
				continue
			}
			erId := efloGetField(item, "ErId")
			if erId == "" {
				continue
			}
			delReq := map[string]interface{}{
				"ErId":     erId,
				"RegionId": client.RegionId,
			}
			efloDelete(client, "eflo", "2022-05-30", "DeleteEr", delReq, "er", erId)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func testSweepEfloVpdGrantRules(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	// ListVpdGrantRules: RegionId required; EnablePage documented (default true).
	// Grant rules carry no name of their own; the authorized network instance
	// name (InstanceName) carries the tf-prefix because tf-test VPDs/VCCs are
	// created with a prefix name and then authorized via the grant rule.
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"EnablePage": true,
		"PageNumber": 1,
		"PageSize":   PageSizeLarge,
	}
	action := "ListVpdGrantRules"
	for {
		response, err := efloRpcPost(client, "eflo", "2022-05-30", action, request)
		if err != nil {
			log.Printf("[ERROR] Failed to list eflo vpd_grant_rules: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.Content.Data", response)
		if err != nil {
			log.Printf("[ERROR] ListVpdGrantRules jsonpath $.Content.Data: %s", err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			name := efloGetField(item, "InstanceName")
			if !efloSweepNameOk(name) {
				continue
			}
			grantRuleId := efloGetField(item, "GrantRuleId")
			if grantRuleId == "" {
				continue
			}
			// DeleteVpdGrantRule requires the full grant context (the
			// resource Delete passes InstanceId/ErId/GrantTenantId from state).
			delReq := map[string]interface{}{
				"GrantRuleId":   grantRuleId,
				"RegionId":      client.RegionId,
				"InstanceId":    efloGetField(item, "InstanceId"),
				"ErId":          efloGetField(item, "ErId"),
				"GrantTenantId": efloGetField(item, "GrantTenantId"),
			}
			efloDelete(client, "eflo", "2022-05-30", "DeleteVpdGrantRule", delReq, "vpd_grant_rule", grantRuleId)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

// ---------------------------------------------------------------------------
// eflo-controller / 2022-12-15 sweepers (NextToken/MaxResults)
// ---------------------------------------------------------------------------

func testSweepEfloClusters(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"MaxResults": PageSizeXLarge,
	}
	action := "ListClusters"
	for {
		response, err := efloRpcPost(client, "eflo-controller", "2022-12-15", action, request)
		if err != nil {
			log.Printf("[ERROR] Failed to list eflo clusters: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.Clusters", response)
		if err != nil {
			log.Printf("[ERROR] ListClusters jsonpath $.Clusters: %s", err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			name := efloGetField(item, "ClusterName")
			if !efloSweepNameOk(name) {
				continue
			}
			clusterId := efloGetField(item, "ClusterId")
			if clusterId == "" {
				continue
			}
			delReq := map[string]interface{}{
				"ClusterId": clusterId,
				"RegionId":  client.RegionId,
			}
			efloDelete(client, "eflo-controller", "2022-12-15", "DeleteCluster", delReq, "cluster", clusterId)
		}
		nt := efloNextToken(response)
		if nt == "" {
			break
		}
		request["NextToken"] = nt
	}
	return nil
}

func testSweepEfloVscs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"MaxResults": PageSizeXLarge,
	}
	action := "ListVscs"
	for {
		response, err := efloRpcPost(client, "eflo-controller", "2022-12-15", action, request)
		if err != nil {
			log.Printf("[ERROR] Failed to list eflo vscs: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.Vscs", response)
		if err != nil {
			log.Printf("[ERROR] ListVscs jsonpath $.Vscs: %s", err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			name := efloGetField(item, "VscName")
			if !efloSweepNameOk(name) {
				continue
			}
			vscId := efloGetField(item, "VscId")
			if vscId == "" {
				continue
			}
			// DeleteVsc accepts an optional ClientToken for idempotency;
			// omitted here — the sweeper is best-effort and retries next pass.
			delReq := map[string]interface{}{
				"VscId":    vscId,
				"RegionId": client.RegionId,
			}
			efloDelete(client, "eflo-controller", "2022-12-15", "DeleteVsc", delReq, "vsc", vscId)
		}
		nt := efloNextToken(response)
		if nt == "" {
			break
		}
		request["NextToken"] = nt
	}
	return nil
}

// testSweepEfloNodeGroups walks clusters (ListClusters) and for each cluster
// lists its node groups (ListNodeGroups, which is ClusterId-scoped per the
// provider's own DescribeEfloNodeGroup service call) before deleting matches.
// This parent-scoped walk mirrors the KMS client_key sweeper (walk AAPs then
// ListClientKeys per AAP). GroupId from the list maps to the NodeGroupId
// parameter of DeleteNodeGroup; the resource ID is ClusterId:NodeGroupId.
func testSweepEfloNodeGroups(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	clusterReq := map[string]interface{}{
		"RegionId":   client.RegionId,
		"MaxResults": PageSizeXLarge,
	}
	for {
		clusterResp, err := efloRpcPost(client, "eflo-controller", "2022-12-15", "ListClusters", clusterReq)
		if err != nil {
			log.Printf("[ERROR] Failed to list eflo clusters for node_group sweep: %s", err)
			return nil
		}
		clustersRaw, err := jsonpath.Get("$.Clusters", clusterResp)
		if err != nil {
			log.Printf("[ERROR] ListClusters jsonpath $.Clusters (node_group sweep): %s", err)
			return nil
		}
		clusters, _ := clustersRaw.([]interface{})
		for _, c := range clusters {
			cluster, _ := c.(map[string]interface{})
			clusterId := efloGetField(cluster, "ClusterId")
			if clusterId == "" {
				continue
			}
			ngReq := map[string]interface{}{
				"ClusterId":  clusterId,
				"RegionId":   client.RegionId,
				"MaxResults": PageSizeXLarge,
			}
			for {
				ngResp, err := efloRpcPost(client, "eflo-controller", "2022-12-15", "ListNodeGroups", ngReq)
				if err != nil {
					log.Printf("[ERROR] Failed to list eflo node_groups (cluster %s): %s", clusterId, err)
					break
				}
				groupsRaw, err := jsonpath.Get("$.Groups", ngResp)
				if err != nil {
					log.Printf("[ERROR] ListNodeGroups jsonpath $.Groups (cluster %s): %s", clusterId, err)
					break
				}
				groups, _ := groupsRaw.([]interface{})
				for _, g := range groups {
					group, _ := g.(map[string]interface{})
					name := efloGetField(group, "GroupName")
					if !efloSweepNameOk(name) {
						continue
					}
					groupId := efloGetField(group, "GroupId")
					if groupId == "" {
						continue
					}
					delReq := map[string]interface{}{
						"ClusterId":   clusterId,
						"NodeGroupId": groupId,
						"RegionId":    client.RegionId,
					}
					efloDelete(client, "eflo-controller", "2022-12-15", "DeleteNodeGroup", delReq, "node_group", clusterId+":"+groupId)
				}
				nt := efloNextToken(ngResp)
				if nt == "" {
					break
				}
				ngReq["NextToken"] = nt
			}
		}
		nt := efloNextToken(clusterResp)
		if nt == "" {
			break
		}
		clusterReq["NextToken"] = nt
	}
	return nil
}

// ---------------------------------------------------------------------------
// eflo-cnp / 2023-08-28 sweepers (Page/Size or single call, $.Data)
// ---------------------------------------------------------------------------

func testSweepEfloExperimentPlans(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	// ListExperimentPlans paginates with Page/Size (Size max observed 100).
	// The plan resource's plan_name attribute is backed by the associated
	// TemplateName (see resourceAliCloudEfloExperimentPlanRead:
	// d.Set("plan_name", objectRaw["TemplateName"])), so TemplateName is the
	// prefix-filter key. NOTE: a plan that references a non-tf-prefixed template
	// (e.g. an official "MatMul" template) is skipped under NORMAL mode; use
	// SWEEP_ALL for a complete fleet-wipe of plans on the test account.
	request := map[string]interface{}{
		"RegionId": client.RegionId,
		"Page":     1,
		"Size":     100,
	}
	action := "ListExperimentPlans"
	for {
		response, err := efloRpcPost(client, "eflo-cnp", "2023-08-28", action, request)
		if err != nil {
			log.Printf("[ERROR] Failed to list eflo experiment_plans: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.Data", response)
		if err != nil {
			log.Printf("[ERROR] ListExperimentPlans jsonpath $.Data: %s", err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			name := efloGetField(item, "TemplateName")
			if !efloSweepNameOk(name) {
				continue
			}
			planId := efloGetField(item, "PlanId")
			if planId == "" {
				continue
			}
			// PlanId/TemplateId are integers in the API; the resource Delete
			// passes d.Id() (a string), so stringify here the same way.
			delReq := map[string]interface{}{
				"PlanId":   planId,
				"RegionId": client.RegionId,
			}
			efloDelete(client, "eflo-cnp", "2023-08-28", "DeleteExperimentPlan", delReq, "experiment_plan", planId)
		}
		if len(result) < 100 {
			break
		}
		request["Page"] = request["Page"].(int) + 1
	}
	return nil
}

func testSweepEfloExperimentPlanTemplates(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	// ListExperimentPlanTemplates returns all templates in a single call (no
	// Page/Size). PrivacyLevel is pinned to "private" so shared/public
	// (potentially official) templates are never swept — even under SWEEP_ALL
	// only the account's own private templates are touched. tf-test fixtures
	// use PrivacyLevel "private".
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"PrivacyLevel": "private",
	}
	action := "ListExperimentPlanTemplates"
	response, err := efloRpcPost(client, "eflo-cnp", "2023-08-28", action, request)
	if err != nil {
		log.Printf("[ERROR] Failed to list eflo experiment_plan_templates: %s", err)
		return nil
	}
	resp, err := jsonpath.Get("$.Data", response)
	if err != nil {
		log.Printf("[ERROR] ListExperimentPlanTemplates jsonpath $.Data: %s", err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item, _ := v.(map[string]interface{})
		name := efloGetField(item, "TemplateName")
		if !efloSweepNameOk(name) {
			continue
		}
		templateId := efloGetField(item, "TemplateId")
		if templateId == "" {
			continue
		}
		delReq := map[string]interface{}{
			"TemplateId": templateId,
			"RegionId":   client.RegionId,
		}
		efloDelete(client, "eflo-cnp", "2023-08-28", "DeleteExperimentPlanTemplate", delReq, "experiment_plan_template", templateId)
	}
	return nil
}
