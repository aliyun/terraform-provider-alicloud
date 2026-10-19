package alicloud

// resource_alicloud_esa_sweeper_test.go — ESA (Edge Security Acceleration) sweeper gap补齐.
//
// Product: ESA, Version: 2024-09-10, Style: RPC (RpcGet/RpcPost).
// alicloud_esa_site is already swept by the inline sweeper in
// resource_alicloud_esa_site_test.go (commit b5c04d195) — do NOT re-register it
// (duplicate registration → init panic).
//
// Gap sweepers registered in this file (44):
//
// Account-level (5):
//   alicloud_esa_list, alicloud_esa_page, alicloud_esa_edge_container_app,
//   alicloud_esa_routine, alicloud_esa_custom_scene_policy
//
// Site-scoped (32 — enumerated per test-prefix site, Dep: alicloud_esa_site):
//   alicloud_esa_cache_rule, alicloud_esa_certificate,
//   alicloud_esa_client_ca_certificate, alicloud_esa_client_certificate,
//   alicloud_esa_compression_rule, alicloud_esa_custom_response_code_rule,
//   alicloud_esa_http_incoming_request_header_modification_rule,
//   alicloud_esa_http_incoming_response_header_modification_rule,
//   alicloud_esa_http_request_header_modification_rule,
//   alicloud_esa_http_response_header_modification_rule,
//   alicloud_esa_https_application_configuration,
//   alicloud_esa_https_basic_configuration, alicloud_esa_image_transform,
//   alicloud_esa_load_balancer, alicloud_esa_network_optimization,
//   alicloud_esa_origin_ca_certificate, alicloud_esa_origin_client_certificate,
//   alicloud_esa_origin_pool, alicloud_esa_origin_protection,
//   alicloud_esa_origin_rule, alicloud_esa_record, alicloud_esa_redirect_rule,
//   alicloud_esa_rewrite_url_rule, alicloud_esa_scheduled_preload_job,
//   alicloud_esa_site_delivery_task, alicloud_esa_site_origin_client_certificate,
//   alicloud_esa_transport_layer_application, alicloud_esa_url_observation,
//   alicloud_esa_version, alicloud_esa_video_processing,
//   alicloud_esa_waiting_room, alicloud_esa_waf_ruleset
//
// Nested (7 — walk parent, Dep on parent sweeper):
//   alicloud_esa_waiting_room_event (Dep: alicloud_esa_waiting_room)
//   alicloud_esa_waiting_room_rule (Dep: alicloud_esa_waiting_room)
//   alicloud_esa_scheduled_preload_execution (Dep: alicloud_esa_scheduled_preload_job)
//   alicloud_esa_waf_rule (Dep: alicloud_esa_waf_ruleset)
//   alicloud_esa_edge_container_app_record (Dep: alicloud_esa_edge_container_app)
//   alicloud_esa_routine_route (Dep: alicloud_esa_routine)
//   alicloud_esa_routine_related_record (Dep: alicloud_esa_routine)
//
// DEFERRED (6 — not enumerable and/or not deletable):
//   alicloud_esa_cache_reserve_instance — no Delete API (billing instance, account-level)
//   alicloud_esa_kv_account — no Delete API (account-level KV store, lifecycle tied to plan)
//   alicloud_esa_kv_namespace — no ListKvNamespaces API to enumerate namespaces
//   alicloud_esa_kv — requires Namespace to enumerate KV pairs; namespace list unavailable
//   alicloud_esa_rate_plan_instance — no Delete API (billing/plan instance)
//   alicloud_esa_routine_code_deployment — no Delete API (deployment action, not a resource)

import (
	"fmt"
	"log"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// ---------------------------------------------------------------------------
// Shared helpers
// ---------------------------------------------------------------------------

// esaSweepPrefixes — test-resource name prefixes (case-insensitive).
var esaSweepPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf-test",
	"tf_test",
	"testAcc",
	"tf-example",
	"tf_example",
	"terraform-example",
	"tftestacc",
}

// esaGetField returns the first non-nil value among the named fields
// (case-insensitive) from a list-response item. ESA response field casing
// is inconsistent across APIs, so this tolerates variant spellings.
func esaGetField(item map[string]interface{}, names ...string) interface{} {
	for k, v := range item {
		lk := strings.ToLower(k)
		for _, name := range names {
			if lk == strings.ToLower(name) {
				return v
			}
		}
	}
	return nil
}

// esaSweepNameOk decides whether a NON-domain resource should be swept.
// NODELETING guard (protectedSweepVpcName) is always-active; sweepAll bypass
// allows SWEEP_ALL to skip the prefix check; otherwise case-insensitive prefix
// match against esaSweepPrefixes.
func esaSweepNameOk(name string) bool {
	if protectedSweepVpcName(name) {
		return false
	}
	if sweepAll() {
		return true
	}
	ln := strings.ToLower(name)
	for _, p := range esaSweepPrefixes {
		if strings.HasPrefix(ln, strings.ToLower(p)) {
			return true
		}
	}
	return false
}

// esaSiteNameOk — ALWAYS-active tf-prefix guard for domain-bearing sites
// (DCDN pattern). NO sweepAll bypass: real存/ICP-backed domains are preserved
// even in SWEEP_ALL. Matches the inline alicloud_esa_site sweeper (b5c04d195).
func esaSiteNameOk(name string) bool {
	if protectedSweepVpcName(name) {
		return false
	}
	ln := strings.ToLower(name)
	for _, p := range esaSweepPrefixes {
		if strings.HasPrefix(ln, strings.ToLower(p)) {
			return true
		}
	}
	return false
}

// esaExtractItems scans a list-response map for the first top-level value that
// is a []interface{} and returns its elements as typed maps. ESA List APIs use
// wildly inconsistent response keys ($.Configs, $.Result, $.Records, $.Rules,
// $.Sites, $.Rulesets, $.WaitingRooms, $.Executions, $.DataModule, …) so a
// generic scan is more robust than hard-coding each key.
func esaExtractItems(response map[string]interface{}) []map[string]interface{} {
	out := make([]map[string]interface{}, 0)
	if response == nil {
		return out
	}
	for _, v := range response {
		arr, ok := v.([]interface{})
		if !ok {
			continue
		}
		for _, e := range arr {
			if m, ok := e.(map[string]interface{}); ok {
				out = append(out, m)
			}
		}
		break // first top-level array wins
	}
	return out
}

// esaIsNotFound returns true for errors that mean the resource is already gone
// (treated as success during best-effort sweep).
func esaIsNotFound(err error) bool {
	if err == nil {
		return true
	}
	if NotFoundError(err) {
		return true
	}
	return IsExpectedErrors(err, []string{
		"NotExists", "NotExist", "NotFound", "ResourceNotFound",
		"SiteNotExist", "ConfigNotExist", "RecordNotExist",
		"RuleNotExist", "WaitingRoomNotExist", "RoutineNotExist",
		"AppNotExist", "ListNotExist", "PageNotExist", "PolicyNotExist",
		"CertificateNotExist", "RulesetNotExist", "VersionNotExist",
		"InvalidParameter.SiteId", "InvalidSiteId", "MissingSiteId",
	})
}

// esaListTestSites enumerates ALL sites and returns only those whose SiteName
// passes the always-active tf-prefix guard (esaSiteNameOk). Used by every
// site-scoped child sweeper to walk test-prefix sites only.
func esaListTestSites(client *connectivity.AliyunClient) []map[string]interface{} {
	sites := make([]map[string]interface{}, 0)
	query := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	for {
		response, err := client.RpcGet("ESA", "2024-09-10", "ListSites", query, nil)
		if err != nil {
			log.Printf("[ERROR] Failed to list ESA Sites: %s", err)
			return sites
		}
		page := esaExtractItems(response)
		for _, s := range page {
			name := fmt.Sprint(esaGetField(s, "SiteName"))
			if name == "" {
				continue
			}
			if esaSiteNameOk(name) {
				sites = append(sites, s)
			}
		}
		if len(page) < PageSizeLarge {
			break
		}
		query["PageNumber"] = query["PageNumber"].(int) + 1
	}
	return sites
}

// esaListSiteItems paginates a site-scoped List API (RpcGet, SiteId param) and
// returns all items via esaExtractItems.
func esaListSiteItems(client *connectivity.AliyunClient, action, siteId string) []map[string]interface{} {
	out := make([]map[string]interface{}, 0)
	query := map[string]interface{}{
		"SiteId":     siteId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	for {
		response, err := client.RpcGet("ESA", "2024-09-10", action, query, nil)
		if err != nil {
			return out
		}
		page := esaExtractItems(response)
		out = append(out, page...)
		if len(page) < PageSizeLarge {
			break
		}
		query["PageNumber"] = query["PageNumber"].(int) + 1
	}
	return out
}

// esaListAccountItems paginates an account-level List API (RpcGet) and returns
// all items via esaExtractItems.
func esaListAccountItems(client *connectivity.AliyunClient, action string) []map[string]interface{} {
	out := make([]map[string]interface{}, 0)
	query := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	for {
		response, err := client.RpcGet("ESA", "2024-09-10", action, query, nil)
		if err != nil {
			return out
		}
		page := esaExtractItems(response)
		out = append(out, page...)
		if len(page) < PageSizeLarge {
			break
		}
		query["PageNumber"] = query["PageNumber"].(int) + 1
	}
	return out
}

// esaDeleteGet fires an RpcGet delete and returns nil on success or not-found.
func esaDeleteGet(client *connectivity.AliyunClient, action string, params map[string]interface{}) error {
	_, err := client.RpcGet("ESA", "2024-09-10", action, params, nil)
	if err != nil && !esaIsNotFound(err) {
		return err
	}
	return nil
}

// esaDeletePost fires an RpcPost delete and returns nil on success or not-found.
func esaDeletePost(client *connectivity.AliyunClient, action string, params map[string]interface{}) error {
	_, err := client.RpcPost("ESA", "2024-09-10", action, nil, params, false)
	if err != nil && !esaIsNotFound(err) {
		return err
	}
	return nil
}

// ---------------------------------------------------------------------------
// init — register all 44 gap sweepers
// ---------------------------------------------------------------------------

func init() {
	// --- Account-level (5) ---
	resource.AddTestSweepers("alicloud_esa_list", &resource.Sweeper{
		Name:         "alicloud_esa_list",
		F:            testSweepEsaLists,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_page", &resource.Sweeper{
		Name:         "alicloud_esa_page",
		F:            testSweepEsaPages,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_edge_container_app", &resource.Sweeper{
		Name:         "alicloud_esa_edge_container_app",
		F:            testSweepEsaEdgeContainerApps,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_routine", &resource.Sweeper{
		Name:         "alicloud_esa_routine",
		F:            testSweepEsaRoutines,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_custom_scene_policy", &resource.Sweeper{
		Name:         "alicloud_esa_custom_scene_policy",
		F:            testSweepEsaCustomScenePolicies,
		Dependencies: []string{"alicloud_esa_site"},
	})

	// --- Site-scoped (32) ---
	resource.AddTestSweepers("alicloud_esa_cache_rule", &resource.Sweeper{
		Name: "alicloud_esa_cache_rule", F: testSweepEsaCacheRules,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_certificate", &resource.Sweeper{
		Name: "alicloud_esa_certificate", F: testSweepEsaCertificates,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_client_ca_certificate", &resource.Sweeper{
		Name: "alicloud_esa_client_ca_certificate", F: testSweepEsaClientCaCertificates,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_client_certificate", &resource.Sweeper{
		Name: "alicloud_esa_client_certificate", F: testSweepEsaClientCertificates,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_compression_rule", &resource.Sweeper{
		Name: "alicloud_esa_compression_rule", F: testSweepEsaCompressionRules,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_custom_response_code_rule", &resource.Sweeper{
		Name: "alicloud_esa_custom_response_code_rule", F: testSweepEsaCustomResponseCodeRules,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_http_incoming_request_header_modification_rule", &resource.Sweeper{
		Name: "alicloud_esa_http_incoming_request_header_modification_rule", F: testSweepEsaHttpIncomingReqHeaderModRules,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_http_incoming_response_header_modification_rule", &resource.Sweeper{
		Name: "alicloud_esa_http_incoming_response_header_modification_rule", F: testSweepEsaHttpIncomingRespHeaderModRules,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_http_request_header_modification_rule", &resource.Sweeper{
		Name: "alicloud_esa_http_request_header_modification_rule", F: testSweepEsaHttpReqHeaderModRules,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_http_response_header_modification_rule", &resource.Sweeper{
		Name: "alicloud_esa_http_response_header_modification_rule", F: testSweepEsaHttpRespHeaderModRules,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_https_application_configuration", &resource.Sweeper{
		Name: "alicloud_esa_https_application_configuration", F: testSweepEsaHttpsApplicationConfigurations,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_https_basic_configuration", &resource.Sweeper{
		Name: "alicloud_esa_https_basic_configuration", F: testSweepEsaHttpsBasicConfigurations,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_image_transform", &resource.Sweeper{
		Name: "alicloud_esa_image_transform", F: testSweepEsaImageTransforms,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_load_balancer", &resource.Sweeper{
		Name: "alicloud_esa_load_balancer", F: testSweepEsaLoadBalancers,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_network_optimization", &resource.Sweeper{
		Name: "alicloud_esa_network_optimization", F: testSweepEsaNetworkOptimizations,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_origin_ca_certificate", &resource.Sweeper{
		Name: "alicloud_esa_origin_ca_certificate", F: testSweepEsaOriginCaCertificates,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_origin_client_certificate", &resource.Sweeper{
		Name: "alicloud_esa_origin_client_certificate", F: testSweepEsaOriginClientCertificates,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_origin_pool", &resource.Sweeper{
		Name: "alicloud_esa_origin_pool", F: testSweepEsaOriginPools,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_origin_protection", &resource.Sweeper{
		Name: "alicloud_esa_origin_protection", F: testSweepEsaOriginProtections,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_origin_rule", &resource.Sweeper{
		Name: "alicloud_esa_origin_rule", F: testSweepEsaOriginRules,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_record", &resource.Sweeper{
		Name: "alicloud_esa_record", F: testSweepEsaRecords,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_redirect_rule", &resource.Sweeper{
		Name: "alicloud_esa_redirect_rule", F: testSweepEsaRedirectRules,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_rewrite_url_rule", &resource.Sweeper{
		Name: "alicloud_esa_rewrite_url_rule", F: testSweepEsaRewriteUrlRules,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_scheduled_preload_job", &resource.Sweeper{
		Name: "alicloud_esa_scheduled_preload_job", F: testSweepEsaScheduledPreloadJobs,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_site_delivery_task", &resource.Sweeper{
		Name: "alicloud_esa_site_delivery_task", F: testSweepEsaSiteDeliveryTasks,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_site_origin_client_certificate", &resource.Sweeper{
		Name: "alicloud_esa_site_origin_client_certificate", F: testSweepEsaSiteOriginClientCertificates,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_transport_layer_application", &resource.Sweeper{
		Name: "alicloud_esa_transport_layer_application", F: testSweepEsaTransportLayerApplications,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_url_observation", &resource.Sweeper{
		Name: "alicloud_esa_url_observation", F: testSweepEsaUrlObservations,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_version", &resource.Sweeper{
		Name: "alicloud_esa_version", F: testSweepEsaVersions,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_video_processing", &resource.Sweeper{
		Name: "alicloud_esa_video_processing", F: testSweepEsaVideoProcessings,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_waiting_room", &resource.Sweeper{
		Name: "alicloud_esa_waiting_room", F: testSweepEsaWaitingRooms,
		Dependencies: []string{"alicloud_esa_site"},
	})
	resource.AddTestSweepers("alicloud_esa_waf_ruleset", &resource.Sweeper{
		Name: "alicloud_esa_waf_ruleset", F: testSweepEsaWafRulesets,
		Dependencies: []string{"alicloud_esa_site"},
	})

	// --- Nested (7) ---
	resource.AddTestSweepers("alicloud_esa_waiting_room_event", &resource.Sweeper{
		Name: "alicloud_esa_waiting_room_event", F: testSweepEsaWaitingRoomEvents,
		Dependencies: []string{"alicloud_esa_waiting_room"},
	})
	resource.AddTestSweepers("alicloud_esa_waiting_room_rule", &resource.Sweeper{
		Name: "alicloud_esa_waiting_room_rule", F: testSweepEsaWaitingRoomRules,
		Dependencies: []string{"alicloud_esa_waiting_room"},
	})
	resource.AddTestSweepers("alicloud_esa_scheduled_preload_execution", &resource.Sweeper{
		Name: "alicloud_esa_scheduled_preload_execution", F: testSweepEsaScheduledPreloadExecutions,
		Dependencies: []string{"alicloud_esa_scheduled_preload_job"},
	})
	resource.AddTestSweepers("alicloud_esa_waf_rule", &resource.Sweeper{
		Name: "alicloud_esa_waf_rule", F: testSweepEsaWafRules,
		Dependencies: []string{"alicloud_esa_waf_ruleset"},
	})
	resource.AddTestSweepers("alicloud_esa_edge_container_app_record", &resource.Sweeper{
		Name: "alicloud_esa_edge_container_app_record", F: testSweepEsaEdgeContainerAppRecords,
		Dependencies: []string{"alicloud_esa_edge_container_app"},
	})
	resource.AddTestSweepers("alicloud_esa_routine_route", &resource.Sweeper{
		Name: "alicloud_esa_routine_route", F: testSweepEsaRoutineRoutes,
		Dependencies: []string{"alicloud_esa_routine"},
	})
	resource.AddTestSweepers("alicloud_esa_routine_related_record", &resource.Sweeper{
		Name: "alicloud_esa_routine_related_record", F: testSweepEsaRoutineRelatedRecords,
		Dependencies: []string{"alicloud_esa_routine"},
	})
}

// ---------------------------------------------------------------------------
// Generic site-scoped sweeper
// ---------------------------------------------------------------------------

// testSweepEsaSiteScoped sweeps a site-scoped ESA resource that:
//   - lists items per test-prefix site via RpcGet(SiteId)
//   - deletes each item by ID (+ optionally SiteId)
//
// idField   — field name in list items holding the resource ID
// deleteAct — ESA Delete action
// useGet    — true: Delete via RpcGet (params in query); false: RpcPost (params in request)
// idParam   — parameter name for the ID in the Delete call
// siteParam — parameter name for SiteId in the Delete call ("" if not needed)
func testSweepEsaSiteScoped(region, listAct, idField, deleteAct string, useGet bool, idParam, siteParam string) error {
	raw, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := raw.(*connectivity.AliyunClient)

	for _, site := range esaListTestSites(client) {
		siteId := fmt.Sprint(esaGetField(site, "SiteId"))
		if siteId == "" {
			continue
		}
		for _, item := range esaListSiteItems(client, listAct, siteId) {
			idVal := esaGetField(item, idField)
			if idVal == nil {
				continue
			}
			params := map[string]interface{}{idParam: idVal}
			if siteParam != "" {
				params[siteParam] = siteId
			}
			if useGet {
				if e := esaDeleteGet(client, deleteAct, params); e != nil {
					log.Printf("[ERROR] ESA %s delete (%s): %s", deleteAct, fmt.Sprint(idVal), e)
				}
			} else {
				if e := esaDeletePost(client, deleteAct, params); e != nil {
					log.Printf("[ERROR] ESA %s delete (%s): %s", deleteAct, fmt.Sprint(idVal), e)
				}
			}
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Account-level sweepers (5)
// ---------------------------------------------------------------------------

// testSweepEsaAccountLevel sweeps an account-level ESA resource that:
//   - lists items via RpcGet (no SiteId)
//   - filters by esaSweepNameOk on the name field
//   - deletes by ID + RegionId
func testSweepEsaAccountLevel(region, listAct, idField, deleteAct string, useGet bool, idParam string, nameFields ...string) error {
	raw, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := raw.(*connectivity.AliyunClient)

	for _, item := range esaListAccountItems(client, listAct) {
		name := fmt.Sprint(esaGetField(item, nameFields...))
		if name == "" || !esaSweepNameOk(name) {
			continue
		}
		idVal := esaGetField(item, idField)
		if idVal == nil {
			continue
		}
		params := map[string]interface{}{
			idParam:    idVal,
			"RegionId": client.RegionId,
		}
		if useGet {
			if e := esaDeleteGet(client, deleteAct, params); e != nil {
				log.Printf("[ERROR] ESA %s delete (%s): %s", deleteAct, fmt.Sprint(idVal), e)
			}
		} else {
			if e := esaDeletePost(client, deleteAct, params); e != nil {
				log.Printf("[ERROR] ESA %s delete (%s): %s", deleteAct, fmt.Sprint(idVal), e)
			}
		}
	}
	return nil
}

func testSweepEsaLists(region string) error {
	return testSweepEsaAccountLevel(region, "ListLists", "Id", "DeleteList", false, "Id", "Name")
}

func testSweepEsaPages(region string) error {
	return testSweepEsaAccountLevel(region, "ListPages", "Id", "DeletePage", false, "Id", "Name")
}

func testSweepEsaEdgeContainerApps(region string) error {
	return testSweepEsaAccountLevel(region, "ListEdgeContainerApps", "AppId", "DeleteEdgeContainerApp", false, "AppId", "AppName", "Name")
}

func testSweepEsaRoutines(region string) error {
	return testSweepEsaAccountLevel(region, "ListUserRoutines", "Name", "DeleteRoutine", false, "Name", "Name")
}

func testSweepEsaCustomScenePolicies(region string) error {
	return testSweepEsaAccountLevel(region, "DescribeCustomScenePolicies", "PolicyId", "DeleteCustomScenePolicy", false, "PolicyId", "PolicyName", "Name")
}

// ---------------------------------------------------------------------------
// Site-scoped sweepers (32)
// ---------------------------------------------------------------------------

func testSweepEsaCacheRules(region string) error {
	return testSweepEsaSiteScoped(region, "ListCacheRules", "ConfigId", "DeleteCacheRule", false, "ConfigId", "SiteId")
}

func testSweepEsaCertificates(region string) error {
	return testSweepEsaSiteScoped(region, "ListCertificates", "Id", "DeleteCertificate", true, "Id", "SiteId")
}

func testSweepEsaClientCaCertificates(region string) error {
	return testSweepEsaSiteScoped(region, "ListClientCaCertificates", "Id", "DeleteClientCaCertificate", true, "Id", "SiteId")
}

func testSweepEsaClientCertificates(region string) error {
	return testSweepEsaSiteScoped(region, "ListClientCertificates", "Id", "DeleteClientCertificate", true, "Id", "SiteId")
}

func testSweepEsaCompressionRules(region string) error {
	return testSweepEsaSiteScoped(region, "ListCompressionRules", "ConfigId", "DeleteCompressionRule", false, "ConfigId", "SiteId")
}

func testSweepEsaCustomResponseCodeRules(region string) error {
	return testSweepEsaSiteScoped(region, "ListCustomResponseCodeRules", "ConfigId", "DeleteCustomResponseCodeRule", false, "ConfigId", "SiteId")
}

func testSweepEsaHttpIncomingReqHeaderModRules(region string) error {
	return testSweepEsaSiteScoped(region, "ListHttpIncomingRequestHeaderModificationRules", "ConfigId", "DeleteHttpIncomingRequestHeaderModificationRule", false, "ConfigId", "SiteId")
}

func testSweepEsaHttpIncomingRespHeaderModRules(region string) error {
	return testSweepEsaSiteScoped(region, "ListHttpIncomingResponseHeaderModificationRules", "ConfigId", "DeleteHttpIncomingResponseHeaderModificationRule", false, "ConfigId", "SiteId")
}

func testSweepEsaHttpReqHeaderModRules(region string) error {
	return testSweepEsaSiteScoped(region, "ListHttpRequestHeaderModificationRules", "ConfigId", "DeleteHttpRequestHeaderModificationRule", false, "ConfigId", "SiteId")
}

func testSweepEsaHttpRespHeaderModRules(region string) error {
	return testSweepEsaSiteScoped(region, "ListHttpResponseHeaderModificationRules", "ConfigId", "DeleteHttpResponseHeaderModificationRule", false, "ConfigId", "SiteId")
}

func testSweepEsaHttpsApplicationConfigurations(region string) error {
	return testSweepEsaSiteScoped(region, "ListHttpsApplicationConfigurations", "ConfigId", "DeleteHttpsApplicationConfiguration", false, "ConfigId", "SiteId")
}

func testSweepEsaHttpsBasicConfigurations(region string) error {
	return testSweepEsaSiteScoped(region, "ListHttpsBasicConfigurations", "ConfigId", "DeleteHttpsBasicConfiguration", false, "ConfigId", "SiteId")
}

func testSweepEsaImageTransforms(region string) error {
	return testSweepEsaSiteScoped(region, "ListImageTransforms", "ConfigId", "DeleteImageTransform", false, "ConfigId", "SiteId")
}

func testSweepEsaLoadBalancers(region string) error {
	return testSweepEsaSiteScoped(region, "ListLoadBalancers", "Id", "DeleteLoadBalancer", false, "Id", "SiteId")
}

func testSweepEsaNetworkOptimizations(region string) error {
	return testSweepEsaSiteScoped(region, "ListNetworkOptimizations", "ConfigId", "DeleteNetworkOptimization", false, "ConfigId", "SiteId")
}

func testSweepEsaOriginCaCertificates(region string) error {
	return testSweepEsaSiteScoped(region, "ListOriginCaCertificates", "Id", "DeleteOriginCaCertificate", true, "Id", "SiteId")
}

func testSweepEsaOriginClientCertificates(region string) error {
	return testSweepEsaSiteScoped(region, "ListOriginClientCertificates", "Id", "DeleteOriginClientCertificate", true, "Id", "SiteId")
}

func testSweepEsaOriginPools(region string) error {
	return testSweepEsaSiteScoped(region, "ListOriginPools", "Id", "DeleteOriginPool", false, "Id", "SiteId")
}

func testSweepEsaOriginRules(region string) error {
	return testSweepEsaSiteScoped(region, "ListOriginRules", "ConfigId", "DeleteOriginRule", false, "ConfigId", "SiteId")
}

func testSweepEsaRecords(region string) error {
	// DeleteRecord uses only RecordId (no SiteId in delete params).
	return testSweepEsaSiteScoped(region, "ListRecords", "RecordId", "DeleteRecord", false, "RecordId", "")
}

func testSweepEsaRedirectRules(region string) error {
	return testSweepEsaSiteScoped(region, "ListRedirectRules", "ConfigId", "DeleteRedirectRule", false, "ConfigId", "SiteId")
}

func testSweepEsaRewriteUrlRules(region string) error {
	return testSweepEsaSiteScoped(region, "ListRewriteUrlRules", "ConfigId", "DeleteRewriteUrlRule", false, "ConfigId", "SiteId")
}

func testSweepEsaScheduledPreloadJobs(region string) error {
	// DeleteScheduledPreloadJob uses only Id (no SiteId in delete params).
	return testSweepEsaSiteScoped(region, "ListScheduledPreloadJobs", "Id", "DeleteScheduledPreloadJob", false, "Id", "")
}

func testSweepEsaSiteDeliveryTasks(region string) error {
	return testSweepEsaSiteScoped(region, "ListSiteDeliveryTasks", "TaskName", "DeleteSiteDeliveryTask", false, "TaskName", "SiteId")
}

func testSweepEsaSiteOriginClientCertificates(region string) error {
	return testSweepEsaSiteScoped(region, "ListSiteOriginClientCertificates", "Id", "DeleteSiteOriginClientCertificate", true, "Id", "SiteId")
}

func testSweepEsaTransportLayerApplications(region string) error {
	return testSweepEsaSiteScoped(region, "ListTransportLayerApplications", "ApplicationId", "DeleteTransportLayerApplication", false, "ApplicationId", "SiteId")
}

func testSweepEsaUrlObservations(region string) error {
	return testSweepEsaSiteScoped(region, "ListUrlObservations", "ConfigId", "DeleteUrlObservation", false, "ConfigId", "SiteId")
}

func testSweepEsaVersions(region string) error {
	return testSweepEsaSiteScoped(region, "ListVersions", "SiteVersion", "DeleteVersion", false, "SiteVersion", "SiteId")
}

func testSweepEsaVideoProcessings(region string) error {
	return testSweepEsaSiteScoped(region, "ListVideoProcessings", "ConfigId", "DeleteVideoProcessing", false, "ConfigId", "SiteId")
}

func testSweepEsaWaitingRooms(region string) error {
	return testSweepEsaSiteScoped(region, "ListWaitingRooms", "WaitingRoomId", "DeleteWaitingRoom", false, "WaitingRoomId", "SiteId")
}

func testSweepEsaWafRulesets(region string) error {
	// DeleteWafRuleset ID format is reversed (RulesetId:SiteId), but for sweeping
	// we extract both fields from the list item, so order doesn't matter.
	return testSweepEsaSiteScoped(region, "ListWafRulesets", "Id", "DeleteWafRuleset", false, "Id", "SiteId")
}

// ---------------------------------------------------------------------------
// Special site-scoped: origin_protection (singleton — no List API)
// ---------------------------------------------------------------------------

func testSweepEsaOriginProtections(region string) error {
	raw, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := raw.(*connectivity.AliyunClient)

	for _, site := range esaListTestSites(client) {
		siteId := fmt.Sprint(esaGetField(site, "SiteId"))
		if siteId == "" {
			continue
		}
		// origin_protection is a singleton per site (no list API). Best-effort
		// delete: if not configured, the API returns a not-found-class error
		// which esaDeletePost treats as success.
		params := map[string]interface{}{"SiteId": siteId}
		if e := esaDeletePost(client, "DeleteOriginProtection", params); e != nil {
			log.Printf("[ERROR] ESA DeleteOriginProtection (site %s): %s", siteId, e)
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Nested sweepers (7)
// ---------------------------------------------------------------------------

// testSweepEsaWaitingRoomEvents — walk test sites → waiting rooms → events.
func testSweepEsaWaitingRoomEvents(region string) error {
	raw, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := raw.(*connectivity.AliyunClient)

	for _, site := range esaListTestSites(client) {
		siteId := fmt.Sprint(esaGetField(site, "SiteId"))
		if siteId == "" {
			continue
		}
		for _, wr := range esaListSiteItems(client, "ListWaitingRooms", siteId) {
			wrId := fmt.Sprint(esaGetField(wr, "WaitingRoomId"))
			if wrId == "" {
				continue
			}
			query := map[string]interface{}{
				"SiteId":        siteId,
				"WaitingRoomId": wrId,
				"PageSize":      PageSizeLarge,
				"PageNumber":    1,
			}
			for {
				resp, e := client.RpcGet("ESA", "2024-09-10", "ListWaitingRoomEvents", query, nil)
				if e != nil {
					break
				}
				page := esaExtractItems(resp)
				for _, item := range page {
					evId := esaGetField(item, "WaitingRoomEventId")
					if evId == nil {
						continue
					}
					params := map[string]interface{}{"SiteId": siteId, "WaitingRoomEventId": evId}
					if e := esaDeletePost(client, "DeleteWaitingRoomEvent", params); e != nil {
						log.Printf("[ERROR] ESA DeleteWaitingRoomEvent (%s): %s", fmt.Sprint(evId), e)
					}
				}
				if len(page) < PageSizeLarge {
					break
				}
				query["PageNumber"] = query["PageNumber"].(int) + 1
			}
		}
	}
	return nil
}

// testSweepEsaWaitingRoomRules — walk test sites → waiting rooms → rules.
func testSweepEsaWaitingRoomRules(region string) error {
	raw, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := raw.(*connectivity.AliyunClient)

	for _, site := range esaListTestSites(client) {
		siteId := fmt.Sprint(esaGetField(site, "SiteId"))
		if siteId == "" {
			continue
		}
		for _, wr := range esaListSiteItems(client, "ListWaitingRooms", siteId) {
			wrId := fmt.Sprint(esaGetField(wr, "WaitingRoomId"))
			if wrId == "" {
				continue
			}
			query := map[string]interface{}{
				"SiteId":        siteId,
				"WaitingRoomId": wrId,
				"PageSize":      PageSizeLarge,
				"PageNumber":    1,
			}
			for {
				resp, e := client.RpcGet("ESA", "2024-09-10", "ListWaitingRoomRules", query, nil)
				if e != nil {
					break
				}
				page := esaExtractItems(resp)
				for _, item := range page {
					rid := esaGetField(item, "WaitingRoomRuleId")
					if rid == nil {
						continue
					}
					params := map[string]interface{}{"SiteId": siteId, "WaitingRoomRuleId": rid}
					if e := esaDeletePost(client, "DeleteWaitingRoomRule", params); e != nil {
						log.Printf("[ERROR] ESA DeleteWaitingRoomRule (%s): %s", fmt.Sprint(rid), e)
					}
				}
				if len(page) < PageSizeLarge {
					break
				}
				query["PageNumber"] = query["PageNumber"].(int) + 1
			}
		}
	}
	return nil
}

// testSweepEsaScheduledPreloadExecutions — walk test sites → jobs → executions.
func testSweepEsaScheduledPreloadExecutions(region string) error {
	raw, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := raw.(*connectivity.AliyunClient)

	for _, site := range esaListTestSites(client) {
		siteId := fmt.Sprint(esaGetField(site, "SiteId"))
		if siteId == "" {
			continue
		}
		for _, job := range esaListSiteItems(client, "ListScheduledPreloadJobs", siteId) {
			jobId := fmt.Sprint(esaGetField(job, "Id"))
			if jobId == "" {
				continue
			}
			query := map[string]interface{}{
				"Id":         jobId,
				"RegionId":   client.RegionId,
				"PageSize":   PageSizeLarge,
				"PageNumber": 1,
			}
			for {
				resp, e := client.RpcGet("ESA", "2024-09-10", "ListScheduledPreloadExecutions", query, nil)
				if e != nil {
					break
				}
				page := esaExtractItems(resp)
				for _, item := range page {
					exeId := esaGetField(item, "Id")
					if exeId == nil {
						continue
					}
					params := map[string]interface{}{"Id": exeId, "RegionId": client.RegionId}
					if e := esaDeletePost(client, "DeleteScheduledPreloadExecution", params); e != nil {
						log.Printf("[ERROR] ESA DeleteScheduledPreloadExecution (%s): %s", fmt.Sprint(exeId), e)
					}
				}
				if len(page) < PageSizeLarge {
					break
				}
				query["PageNumber"] = query["PageNumber"].(int) + 1
			}
		}
	}
	return nil
}

// testSweepEsaWafRules — walk test sites × 8 phases → rules.
func testSweepEsaWafRules(region string) error {
	raw, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := raw.(*connectivity.AliyunClient)

	phases := []string{
		"http_whitelist", "http_custom", "http_managed",
		"http_anti_scan", "http_ratelimit", "ip_access_rule",
		"http_bot", "http_security_level_rule",
	}

	for _, site := range esaListTestSites(client) {
		siteId := fmt.Sprint(esaGetField(site, "SiteId"))
		if siteId == "" {
			continue
		}
		for _, phase := range phases {
			query := map[string]interface{}{
				"SiteId":     siteId,
				"Phase":      phase,
				"PageSize":   PageSizeLarge,
				"PageNumber": 1,
			}
			for {
				resp, e := client.RpcGet("ESA", "2024-09-10", "ListWafRules", query, nil)
				if e != nil {
					break
				}
				page := esaExtractItems(resp)
				for _, item := range page {
					rid := esaGetField(item, "Id")
					if rid == nil {
						continue
					}
					params := map[string]interface{}{"Id": rid, "SiteId": siteId}
					if e := esaDeletePost(client, "DeleteWafRule", params); e != nil {
						log.Printf("[ERROR] ESA DeleteWafRule (%s): %s", fmt.Sprint(rid), e)
					}
				}
				if len(page) < PageSizeLarge {
					break
				}
				query["PageNumber"] = query["PageNumber"].(int) + 1
			}
		}
	}
	return nil
}

// testSweepEsaEdgeContainerAppRecords — walk test apps → records.
func testSweepEsaEdgeContainerAppRecords(region string) error {
	raw, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := raw.(*connectivity.AliyunClient)

	for _, app := range esaListAccountItems(client, "ListEdgeContainerApps") {
		appName := fmt.Sprint(esaGetField(app, "AppName", "Name"))
		if appName == "" || !esaSweepNameOk(appName) {
			continue
		}
		appId := fmt.Sprint(esaGetField(app, "AppId"))
		if appId == "" {
			continue
		}
		query := map[string]interface{}{
			"AppId":      appId,
			"PageSize":   PageSizeLarge,
			"PageNumber": 1,
		}
		for {
			resp, e := client.RpcGet("ESA", "2024-09-10", "ListEdgeContainerAppRecords", query, nil)
			if e != nil {
				break
			}
			page := esaExtractItems(resp)
			for _, item := range page {
				recName := fmt.Sprint(esaGetField(item, "RecordName"))
				siteId := fmt.Sprint(esaGetField(item, "SiteId"))
				if recName == "" {
					continue
				}
				params := map[string]interface{}{
					"AppId":      appId,
					"SiteId":     siteId,
					"RecordName": recName,
				}
				if e := esaDeletePost(client, "DeleteEdgeContainerAppRecord", params); e != nil {
					log.Printf("[ERROR] ESA DeleteEdgeContainerAppRecord (%s): %s", recName, e)
				}
			}
			if len(page) < PageSizeLarge {
				break
			}
			query["PageNumber"] = query["PageNumber"].(int) + 1
		}
	}
	return nil
}

// testSweepEsaRoutineRoutes — walk test routines → routes.
func testSweepEsaRoutineRoutes(region string) error {
	raw, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := raw.(*connectivity.AliyunClient)

	for _, rt := range esaListAccountItems(client, "ListUserRoutines") {
		name := fmt.Sprint(esaGetField(rt, "Name"))
		if name == "" || !esaSweepNameOk(name) {
			continue
		}
		query := map[string]interface{}{
			"RoutineName": name,
			"PageSize":    PageSizeLarge,
			"PageNumber":  1,
		}
		for {
			resp, e := client.RpcGet("ESA", "2024-09-10", "ListRoutineRoutes", query, nil)
			if e != nil {
				break
			}
			page := esaExtractItems(resp)
			for _, item := range page {
				cfgId := esaGetField(item, "ConfigId")
				siteId := fmt.Sprint(esaGetField(item, "SiteId"))
				if cfgId == nil {
					continue
				}
				params := map[string]interface{}{"ConfigId": cfgId, "SiteId": siteId}
				if e := esaDeletePost(client, "DeleteRoutineRoute", params); e != nil {
					log.Printf("[ERROR] ESA DeleteRoutineRoute (%s): %s", fmt.Sprint(cfgId), e)
				}
			}
			if len(page) < PageSizeLarge {
				break
			}
			query["PageNumber"] = query["PageNumber"].(int) + 1
		}
	}
	return nil
}

// testSweepEsaRoutineRelatedRecords — walk test routines → related records.
func testSweepEsaRoutineRelatedRecords(region string) error {
	raw, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := raw.(*connectivity.AliyunClient)

	for _, rt := range esaListAccountItems(client, "ListUserRoutines") {
		name := fmt.Sprint(esaGetField(rt, "Name"))
		if name == "" || !esaSweepNameOk(name) {
			continue
		}
		req := map[string]interface{}{
			"Name":       name,
			"PageSize":   PageSizeLarge,
			"PageNumber": 1,
		}
		for {
			// ListRoutineRelatedRecords uses RpcPost in the service file.
			resp, e := client.RpcPost("ESA", "2024-09-10", "ListRoutineRelatedRecords", nil, req, false)
			if e != nil {
				break
			}
			page := esaExtractItems(resp)
			for _, item := range page {
				recId := esaGetField(item, "RecordId")
				siteId := fmt.Sprint(esaGetField(item, "SiteId"))
				recName := fmt.Sprint(esaGetField(item, "RecordName"))
				if recId == nil {
					continue
				}
				params := map[string]interface{}{
					"RecordId":   recId,
					"Name":       name,
					"SiteId":     siteId,
					"RecordName": recName,
				}
				if e := esaDeletePost(client, "DeleteRoutineRelatedRecord", params); e != nil {
					log.Printf("[ERROR] ESA DeleteRoutineRelatedRecord (%s): %s", fmt.Sprint(recId), e)
				}
			}
			if len(page) < PageSizeLarge {
				break
			}
			req["PageNumber"] = req["PageNumber"].(int) + 1
		}
	}
	return nil
}
