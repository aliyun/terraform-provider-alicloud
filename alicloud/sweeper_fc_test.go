package alicloud

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// This file 补齐 (fills in) missing sweepers for Function Compute (FC) v1
// resources that previously had none, plus all FC v3 (fcv3) resources. The
// existing v1 service/function/trigger/async-invoke-config sweepers already
// work and remain untouched. The v1 custom_domain sweeper (which previously
// deleted ALL domains with no protection) was fixed in place to honor the
// real存 preserve red-line.
//
// Red-line guards:
//   - Name-based NODELETING/do-not-delete protection via protectedSweepVpcName.
//   - SWEEP_ALL bypasses the tf- prefix filter for ordinary resources, but the
//     prefix filter ALWAYS applies for custom domains (real存 ICP preserve).
//   - Per-instance delete failures are logged and the loop continues (non-fatal);
//     the sweeper returns nil unless sharedClientForRegion itself fails.

// fcSweepPrefixes is the standard tf-testAcc name prefix list shared across FC
// sweepers. Matching is case-insensitive.
var fcSweepPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
}

// fcSweepShouldDelete decides whether to delete an FC resource by name.
// fcSweepShouldDelete decides whether to delete an FC resource by name.
// NODELETING/do-not-delete protection always applies. The tf-testAcc prefix
// filter ALWAYS applies too — FC resources are frequently real services (not
// just test leftovers) and may lack NODELETING markers, so a SWEEP_ALL bypass
// of the prefix check is unsafe for FC. This matches the canonical existing v1
// FC sweepers (alicloud_fc_service/function/trigger/...), which always apply
// their prefix list and never consult sweepAll().
func fcSweepShouldDelete(name string) (bool, string) {
	if protectedSweepVpcName(name) {
		return false, "NODELETING/do-not-delete marker in name"
	}
	for _, p := range fcSweepPrefixes {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(p)) {
			return true, ""
		}
	}
	return false, "name does not match tf-testAcc prefix"
}

// fcSweepDomainShouldDelete is the custom-domain variant. The tf- prefix filter
// ALWAYS applies, even under SWEEP_ALL, because FC custom domains may back
// real存 IP-backed sites that must never be deleted.
func fcSweepDomainShouldDelete(name string) (bool, string) {
	if protectedSweepVpcName(name) {
		return false, "NODELETING/do-not-delete marker in name"
	}
	for _, p := range fcSweepPrefixes {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(p)) {
			return true, ""
		}
	}
	return false, "domain does not match tf-testAcc prefix (real存 preserve)"
}

// =============================================================================
// FC v1 — alicloud_fc_alias
// =============================================================================

func init() {
	resource.AddTestSweepers("alicloud_fc_alias", &resource.Sweeper{
		Name: "alicloud_fc_alias",
		F:    testSweepFCAlias,
		Dependencies: []string{
			"alicloud_fc_function",
		},
	})
}

func testSweepFCAlias(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.FcNoSupportedRegions) {
		log.Printf("[INFO] Skipping Function Compute unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		return fcClient.ListServices(fc.NewListServicesInput())
	})
	if err != nil {
		log.Printf("Error retrieving FC services: %s", err)
		return nil
	}
	services, _ := raw.(*fc.ListServicesOutput)
	if services == nil {
		return nil
	}
	for _, svc := range services.Services {
		if svc.ServiceName == nil {
			continue
		}
		serviceName := *svc.ServiceName
		// Skip the whole service unless it matches the test prefix; this
		// avoids walking non-test services and trimming their aliases.
		if svcDel, _ := fcSweepShouldDelete(serviceName); !svcDel {
			continue
		}
		nextToken := ""
		for {
			args := fc.NewListAliasesInput(serviceName)
			if nextToken != "" {
				args.NextToken = &nextToken
			}
			raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
				return fcClient.ListAliases(args)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to list aliases of service (%s): %s", serviceName, err)
				break
			}
			resp, _ := raw.(*fc.ListAliasesOutput)
			if resp == nil || resp.Aliases == nil || len(resp.Aliases) < 1 {
				break
			}
			for _, alias := range resp.Aliases {
				if alias.AliasName == nil {
					continue
				}
				aliasName := *alias.AliasName
				if del, reason := fcSweepShouldDelete(aliasName); !del {
					log.Printf("[INFO] Skipping FC alias (%s): %s", aliasName, reason)
					continue
				}
				log.Printf("[INFO] Deleting FC alias: %s/%s", serviceName, aliasName)
				if _, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
					return fcClient.DeleteAlias(fc.NewDeleteAliasInput(serviceName, aliasName))
				}); err != nil {
					log.Printf("[ERROR] Failed to delete FC alias (%s/%s): %s", serviceName, aliasName, err)
				}
			}
			if resp.NextToken == nil || *resp.NextToken == "" {
				break
			}
			nextToken = *resp.NextToken
		}
	}
	return nil
}

// =============================================================================
// FC v1 — alicloud_fc_layer_version (FC-Open API: /2021-04-06/layers/...)
// =============================================================================

func init() {
	resource.AddTestSweepers("alicloud_fc_layer_version", &resource.Sweeper{
		Name: "alicloud_fc_layer_version",
		F:    testSweepFcLayerVersion,
	})
}

func testSweepFcLayerVersion(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.FcNoSupportedRegions) {
		log.Printf("[INFO] Skipping Function Compute unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	// List layers via the FC-Open API surface used by the resource itself.
	for {
		action := "/2021-04-06/layers"
		query := make(map[string]*string)
		query["limit"] = StringPointer(fmt.Sprint(PageSizeLarge))
		response, err := client.RoaGet("FC-Open", "2021-04-06", action, query, nil, nil)
		if err != nil {
			log.Printf("[ERROR] Failed to list FC layers: %s", err)
			return nil
		}
		rawLayers, _ := jsonpathGet("$.layers", response)
		layers, _ := rawLayers.([]interface{})
		for _, l := range layers {
			item, _ := l.(map[string]interface{})
			if item == nil {
				continue
			}
			layerName := fmt.Sprint(item["layerName"])
			if layerName == "" {
				continue
			}
			if del, reason := fcSweepShouldDelete(layerName); !del {
				log.Printf("[INFO] Skipping FC layer (%s): %s", layerName, reason)
				continue
			}
			// List versions for the layer, then delete each.
			versionsURL := fmt.Sprintf("/2021-04-06/layers/%s/versions", layerName)
			vq := make(map[string]*string)
			vq["limit"] = StringPointer(fmt.Sprint(PageSizeLarge))
			versionsResp, vErr := client.RoaGet("FC-Open", "2021-04-06", versionsURL, vq, nil, nil)
			if vErr != nil {
				log.Printf("[ERROR] Failed to list versions of FC layer (%s): %s", layerName, vErr)
				continue
			}
			rawVersions, _ := jsonpathGet("$.versions", versionsResp)
			versions, _ := rawVersions.([]interface{})
			for _, ver := range versions {
				vItem, _ := ver.(map[string]interface{})
				if vItem == nil {
					continue
				}
				version := fmt.Sprint(vItem["version"])
				if version == "" {
					continue
				}
				log.Printf("[INFO] Deleting FC layer version: %s/%s", layerName, version)
				deleteAction := fmt.Sprintf("/2021-04-06/layers/%s/versions/%s", layerName, version)
				if _, err := client.RoaDelete("FC-Open", "2021-04-06", deleteAction, nil, nil, nil, false); err != nil {
					log.Printf("[ERROR] Failed to delete FC layer version (%s/%s): %s", layerName, version, err)
				}
			}
		}
		nextToken, _ := response["nextToken"].(string)
		if nextToken == "" {
			break
		}
	}
	return nil
}

// jsonpathGet wraps jsonpath.Get so an error (e.g. path not found) returns nil
// without requiring callers to handle it inline. FC list responses reliably
// use $.<key> arrays, so a missing path just yields an empty list.
func jsonpathGet(path string, object interface{}) (interface{}, error) {
	return jsonpath.Get(path, object)
}

// keysOf returns the top-level keys of a RoaGet response map for debug logging
// (e.g. diagnosing why an FC list endpoint yields an unexpected shape).
func keysOf(m map[string]interface{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// =============================================================================
// FC v3 — alicloud_fcv3_function (top-level: lists and deletes functions)
// =============================================================================

func init() {
	resource.AddTestSweepers("alicloud_fcv3_function", &resource.Sweeper{
		Name: "alicloud_fcv3_function",
		F:    testSweepFcv3Function,
		Dependencies: []string{
			"alicloud_fcv3_trigger",
			"alicloud_fcv3_alias",
			"alicloud_fcv3_async_invoke_config",
			"alicloud_fcv3_concurrency_config",
			"alicloud_fcv3_function_version",
			"alicloud_fcv3_provision_config",
			"alicloud_fcv3_vpc_binding",
		},
	})
}

func testSweepFcv3Function(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.FcNoSupportedRegions) {
		log.Printf("[INFO] Skipping Function Compute unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	functions, err := fcv3ListFunctions(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list fcv3 functions: %s", err)
		return nil
	}
	for _, fn := range functions {
		name := fmt.Sprint(fn["functionName"])
		if name == "" {
			continue
		}
		if del, reason := fcSweepShouldDelete(name); !del {
			log.Printf("[INFO] Skipping fcv3 function (%s): %s", name, reason)
			continue
		}
		log.Printf("[INFO] Deleting fcv3 function: %s", name)
		action := fmt.Sprintf("/2023-03-30/functions/%s", name)
		if _, err := client.RoaDelete("FC", "2023-03-30", action, nil, nil, nil, true); err != nil {
			log.Printf("[ERROR] Failed to delete fcv3 function (%s): %s", name, err)
		}
	}
	return nil
}

// fcv3ListFunctions paginates through /2023-03-30/functions using nextToken.
func fcv3ListFunctions(client *connectivity.AliyunClient) ([]map[string]interface{}, error) {
	var out []map[string]interface{}
	query := make(map[string]*string)
	query["limit"] = StringPointer(fmt.Sprint(PageSizeLarge))
	for {
		response, err := client.RoaGet("FC", "2023-03-30", "/2023-03-30/functions", query, nil, nil)
		if err != nil {
			return nil, err
		}
		rawFns, _ := jsonpathGet("$.functions", response)
		fns, _ := rawFns.([]interface{})
		for _, f := range fns {
			item, _ := f.(map[string]interface{})
			if item != nil {
				out = append(out, item)
			}
		}
		nextToken, _ := response["nextToken"].(string)
		if nextToken == "" {
			break
		}
		query["nextToken"] = StringPointer(nextToken)
	}
	return out, nil
}

// =============================================================================
// FC v3 — alicloud_fcv3_alias (per-function sub-resource)
// =============================================================================

func init() {
	resource.AddTestSweepers("alicloud_fcv3_alias", &resource.Sweeper{
		Name: "alicloud_fcv3_alias",
		F:    testSweepFcv3Alias,
	})
}

func testSweepFcv3Alias(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.FcNoSupportedRegions) {
		log.Printf("[INFO] Skipping Function Compute unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	functions, err := fcv3ListFunctions(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list fcv3 functions: %s", err)
		return nil
	}
	for _, fn := range functions {
		fnName := fmt.Sprint(fn["functionName"])
		if fnName == "" {
			continue
		}
		if del, _ := fcSweepShouldDelete(fnName); !del {
			continue
		}
		nextToken := ""
		for {
			action := fmt.Sprintf("/2023-03-30/functions/%s/aliases", fnName)
			query := make(map[string]*string)
			query["limit"] = StringPointer(fmt.Sprint(PageSizeLarge))
			if nextToken != "" {
				query["nextToken"] = StringPointer(nextToken)
			}
			response, err := client.RoaGet("FC", "2023-03-30", action, query, nil, nil)
			if err != nil {
				log.Printf("[ERROR] Failed to list fcv3 aliases for function (%s): %s", fnName, err)
				break
			}
			rawAliases, _ := jsonpathGet("$.aliases", response)
			aliases, _ := rawAliases.([]interface{})
			for _, a := range aliases {
				item, _ := a.(map[string]interface{})
				if item == nil {
					continue
				}
				aliasName := fmt.Sprint(item["aliasName"])
				if aliasName == "" {
					continue
				}
				if del, reason := fcSweepShouldDelete(aliasName); !del {
					log.Printf("[INFO] Skipping fcv3 alias (%s): %s", aliasName, reason)
					continue
				}
				log.Printf("[INFO] Deleting fcv3 alias: %s/%s", fnName, aliasName)
				deleteAction := fmt.Sprintf("/2023-03-30/functions/%s/aliases/%s", fnName, aliasName)
				if _, err := client.RoaDelete("FC", "2023-03-30", deleteAction, nil, nil, nil, true); err != nil {
					log.Printf("[ERROR] Failed to delete fcv3 alias (%s/%s): %s", fnName, aliasName, err)
				}
			}
			nextToken, _ = response["nextToken"].(string)
			if nextToken == "" {
				break
			}
		}
	}
	return nil
}

// =============================================================================
// FC v3 — alicloud_fcv3_async_invoke_config (per-function, singleton config)
// =============================================================================

func init() {
	resource.AddTestSweepers("alicloud_fcv3_async_invoke_config", &resource.Sweeper{
		Name: "alicloud_fcv3_async_invoke_config",
		F:    testSweepFcv3AsyncInvokeConfig,
	})
}

func testSweepFcv3AsyncInvokeConfig(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.FcNoSupportedRegions) {
		log.Printf("[INFO] Skipping Function Compute unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	functions, err := fcv3ListFunctions(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list fcv3 functions: %s", err)
		return nil
	}
	for _, fn := range functions {
		fnName := fmt.Sprint(fn["functionName"])
		if fnName == "" {
			continue
		}
		if del, _ := fcSweepShouldDelete(fnName); !del {
			continue
		}
		action := fmt.Sprintf("/2023-03-30/functions/%s/async-invoke-config", fnName)
		// Async invoke config is a singleton per function; attempt delete and
		// tolerate AsyncConfigNotExists / NotFoundError as no-ops.
		if _, err := client.RoaDelete("FC", "2023-03-30", action, nil, nil, nil, true); err != nil {
			if !isFcv3NotFoundError(err) {
				log.Printf("[ERROR] Failed to delete fcv3 async-invoke-config for function (%s): %s", fnName, err)
			}
		} else {
			log.Printf("[INFO] Deleted fcv3 async-invoke-config for function: %s", fnName)
		}
	}
	return nil
}

// =============================================================================
// FC v3 — alicloud_fcv3_concurrency_config (per-function, singleton config)
// =============================================================================

func init() {
	resource.AddTestSweepers("alicloud_fcv3_concurrency_config", &resource.Sweeper{
		Name: "alicloud_fcv3_concurrency_config",
		F:    testSweepFcv3ConcurrencyConfig,
	})
}

func testSweepFcv3ConcurrencyConfig(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.FcNoSupportedRegions) {
		log.Printf("[INFO] Skipping Function Compute unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	functions, err := fcv3ListFunctions(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list fcv3 functions: %s", err)
		return nil
	}
	for _, fn := range functions {
		fnName := fmt.Sprint(fn["functionName"])
		if fnName == "" {
			continue
		}
		if del, _ := fcSweepShouldDelete(fnName); !del {
			continue
		}
		action := fmt.Sprintf("/2023-03-30/functions/%s/concurrency", fnName)
		if _, err := client.RoaDelete("FC", "2023-03-30", action, nil, nil, nil, true); err != nil {
			if !isFcv3NotFoundError(err) {
				log.Printf("[ERROR] Failed to delete fcv3 concurrency-config for function (%s): %s", fnName, err)
			}
		} else {
			log.Printf("[INFO] Deleted fcv3 concurrency-config for function: %s", fnName)
		}
	}
	return nil
}

// =============================================================================
// FC v3 — alicloud_fcv3_custom_domain (top-level; real存 preserve)
// =============================================================================

func init() {
	resource.AddTestSweepers("alicloud_fcv3_custom_domain", &resource.Sweeper{
		Name: "alicloud_fcv3_custom_domain",
		F:    testSweepFcv3CustomDomain,
	})
}

func testSweepFcv3CustomDomain(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.FcNoSupportedRegions) {
		log.Printf("[INFO] Skipping Function Compute unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	nextToken := ""
	for {
		query := make(map[string]*string)
		query["limit"] = StringPointer(fmt.Sprint(PageSizeLarge))
		if nextToken != "" {
			query["nextToken"] = StringPointer(nextToken)
		}
		response, err := client.RoaGet("FC", "2023-03-30", "/2023-03-30/custom-domains", query, nil, nil)
		if err != nil {
			log.Printf("[ERROR] Failed to list fcv3 custom-domains: %s", err)
			return nil
		}
		rawDomains, _ := jsonpathGet("$.customDomains", response)
		domains, _ := rawDomains.([]interface{})
		for _, d := range domains {
			item, _ := d.(map[string]interface{})
			if item == nil {
				continue
			}
			domainName := fmt.Sprint(item["domainName"])
			if domainName == "" {
				continue
			}
			// Custom domains may back real存 ICP-bound sites; the tf- prefix
			// filter ALWAYS applies, even under SWEEP_ALL.
			if del, reason := fcSweepDomainShouldDelete(domainName); !del {
				log.Printf("[INFO] Skipping fcv3 custom-domain (%s): %s", domainName, reason)
				continue
			}
			log.Printf("[INFO] Deleting fcv3 custom-domain: %s", domainName)
			deleteAction := fmt.Sprintf("/2023-03-30/custom-domains/%s", domainName)
			if _, err := client.RoaDelete("FC", "2023-03-30", deleteAction, nil, nil, nil, true); err != nil {
				log.Printf("[ERROR] Failed to delete fcv3 custom-domain (%s): %s", domainName, err)
			}
		}
		nextToken, _ = response["nextToken"].(string)
		if nextToken == "" {
			break
		}
	}
	return nil
}

// =============================================================================
// FC v3 — alicloud_fcv3_function_version (per-function sub-resource)
// =============================================================================

func init() {
	resource.AddTestSweepers("alicloud_fcv3_function_version", &resource.Sweeper{
		Name: "alicloud_fcv3_function_version",
		F:    testSweepFcv3FunctionVersion,
	})
}

func testSweepFcv3FunctionVersion(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.FcNoSupportedRegions) {
		log.Printf("[INFO] Skipping Function Compute unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	functions, err := fcv3ListFunctions(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list fcv3 functions: %s", err)
		return nil
	}
	for _, fn := range functions {
		fnName := fmt.Sprint(fn["functionName"])
		if fnName == "" {
			continue
		}
		if del, _ := fcSweepShouldDelete(fnName); !del {
			continue
		}
		nextToken := ""
		for {
			action := fmt.Sprintf("/2023-03-30/functions/%s/versions", fnName)
			query := make(map[string]*string)
			query["limit"] = StringPointer(fmt.Sprint(PageSizeLarge))
			if nextToken != "" {
				query["nextToken"] = StringPointer(nextToken)
			}
			response, err := client.RoaGet("FC", "2023-03-30", action, query, nil, nil)
			if err != nil {
				log.Printf("[ERROR] Failed to list fcv3 versions for function (%s): %s", fnName, err)
				break
			}
			rawVers, _ := jsonpathGet("$.versions", response)
			versions, _ := rawVers.([]interface{})
			for _, vItem := range versions {
				item, _ := vItem.(map[string]interface{})
				if item == nil {
					continue
				}
				versionID := fmt.Sprint(item["version"])
				if versionID == "" {
					continue
				}
				log.Printf("[INFO] Deleting fcv3 function version: %s/%s", fnName, versionID)
				deleteAction := fmt.Sprintf("/2023-03-30/functions/%s/versions/%s", fnName, versionID)
				if _, err := client.RoaDelete("FC", "2023-03-30", deleteAction, nil, nil, nil, true); err != nil {
					log.Printf("[ERROR] Failed to delete fcv3 function version (%s/%s): %s", fnName, versionID, err)
				}
			}
			nextToken, _ = response["nextToken"].(string)
			if nextToken == "" {
				break
			}
		}
	}
	return nil
}

// =============================================================================
// FC v3 — alicloud_fcv3_layer_version (top-level; list layers then versions)
// =============================================================================

func init() {
	resource.AddTestSweepers("alicloud_fcv3_layer_version", &resource.Sweeper{
		Name: "alicloud_fcv3_layer_version",
		F:    testSweepFcv3LayerVersion,
	})
}

func testSweepFcv3LayerVersion(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.FcNoSupportedRegions) {
		log.Printf("[INFO] Skipping Function Compute unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	nextToken := ""
	for {
		query := make(map[string]*string)
		query["limit"] = StringPointer(fmt.Sprint(PageSizeLarge))
		if nextToken != "" {
			query["nextToken"] = StringPointer(nextToken)
		}
		response, err := client.RoaGet("FC", "2023-03-30", "/2023-03-30/layers", query, nil, nil)
		if err != nil {
			log.Printf("[ERROR] Failed to list fcv3 layers: %s", err)
			return nil
		}
		// Use direct map access (the RoaGet response is map[string]interface{});
		// jsonpath is avoided here because $.layers was observed to yield nil on the
		// /layers endpoint even when layers exist, whereas direct access is reliable.
		layers := []interface{}{}
		if raw, ok := response["layers"]; ok {
			if arr, ok := raw.([]interface{}); ok {
				layers = arr
			}
		}
		log.Printf("[DEBUG] fcv3 layer list returned %d layer(s) (response keys: %v)", len(layers), keysOf(response))
		for _, l := range layers {
			item, _ := l.(map[string]interface{})
			if item == nil {
				continue
			}
			layerName := fmt.Sprint(item["layerName"])
			if layerName == "" {
				continue
			}
			if del, reason := fcSweepShouldDelete(layerName); !del {
				log.Printf("[INFO] Skipping fcv3 layer (%s): %s", layerName, reason)
				continue
			}
			// List versions for the layer.
			versionsAction := fmt.Sprintf("/2023-03-30/layers/%s/versions", layerName)
			vq := make(map[string]*string)
			vq["limit"] = StringPointer(fmt.Sprint(PageSizeLarge))
			versionsResp, vErr := client.RoaGet("FC", "2023-03-30", versionsAction, vq, nil, nil)
			if vErr != nil {
				log.Printf("[ERROR] Failed to list fcv3 layer versions for layer (%s): %s", layerName, vErr)
				continue
			}
			rawVers, _ := jsonpathGet("$.versions", versionsResp)
			versions, _ := rawVers.([]interface{})
			if versions == nil {
				if raw, ok := versionsResp["versions"]; ok {
					if arr, ok := raw.([]interface{}); ok {
						versions = arr
					}
				}
			}
			for _, vItem := range versions {
				verMap, _ := vItem.(map[string]interface{})
				if verMap == nil {
					continue
				}
				version := fmt.Sprint(verMap["version"])
				if version == "" {
					continue
				}
				log.Printf("[INFO] Deleting fcv3 layer version: %s/%s", layerName, version)
				deleteAction := fmt.Sprintf("/2023-03-30/layers/%s/versions/%s", layerName, version)
				if _, err := client.RoaDelete("FC", "2023-03-30", deleteAction, nil, nil, nil, true); err != nil {
					log.Printf("[ERROR] Failed to delete fcv3 layer version (%s/%s): %s", layerName, version, err)
				}
			}
		}
		nextToken, _ = response["nextToken"].(string)
		if nextToken == "" {
			break
		}
	}
	return nil
}

// =============================================================================
// FC v3 — alicloud_fcv3_provision_config (per-function, singleton config)
// =============================================================================

func init() {
	resource.AddTestSweepers("alicloud_fcv3_provision_config", &resource.Sweeper{
		Name: "alicloud_fcv3_provision_config",
		F:    testSweepFcv3ProvisionConfig,
	})
}

func testSweepFcv3ProvisionConfig(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.FcNoSupportedRegions) {
		log.Printf("[INFO] Skipping Function Compute unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	functions, err := fcv3ListFunctions(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list fcv3 functions: %s", err)
		return nil
	}
	for _, fn := range functions {
		fnName := fmt.Sprint(fn["functionName"])
		if fnName == "" {
			continue
		}
		if del, _ := fcSweepShouldDelete(fnName); !del {
			continue
		}
		action := fmt.Sprintf("/2023-03-30/functions/%s/provision-config", fnName)
		if _, err := client.RoaDelete("FC", "2023-03-30", action, nil, nil, nil, true); err != nil {
			if !isFcv3NotFoundError(err) {
				log.Printf("[ERROR] Failed to delete fcv3 provision-config for function (%s): %s", fnName, err)
			}
		} else {
			log.Printf("[INFO] Deleted fcv3 provision-config for function: %s", fnName)
		}
	}
	return nil
}

// =============================================================================
// FC v3 — alicloud_fcv3_trigger (per-function sub-resource)
// =============================================================================

func init() {
	resource.AddTestSweepers("alicloud_fcv3_trigger", &resource.Sweeper{
		Name: "alicloud_fcv3_trigger",
		F:    testSweepFcv3Trigger,
	})
}

func testSweepFcv3Trigger(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.FcNoSupportedRegions) {
		log.Printf("[INFO] Skipping Function Compute unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	functions, err := fcv3ListFunctions(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list fcv3 functions: %s", err)
		return nil
	}
	for _, fn := range functions {
		fnName := fmt.Sprint(fn["functionName"])
		if fnName == "" {
			continue
		}
		if del, _ := fcSweepShouldDelete(fnName); !del {
			continue
		}
		nextToken := ""
		for {
			action := fmt.Sprintf("/2023-03-30/functions/%s/triggers", fnName)
			query := make(map[string]*string)
			query["limit"] = StringPointer(fmt.Sprint(PageSizeLarge))
			if nextToken != "" {
				query["nextToken"] = StringPointer(nextToken)
			}
			response, err := client.RoaGet("FC", "2023-03-30", action, query, nil, nil)
			if err != nil {
				log.Printf("[ERROR] Failed to list fcv3 triggers for function (%s): %s", fnName, err)
				break
			}
			rawTrigs, _ := jsonpathGet("$.triggers", response)
			triggers, _ := rawTrigs.([]interface{})
			for _, tItem := range triggers {
				item, _ := tItem.(map[string]interface{})
				if item == nil {
					continue
				}
				triggerName := fmt.Sprint(item["triggerName"])
				if triggerName == "" {
					continue
				}
				if del, reason := fcSweepShouldDelete(triggerName); !del {
					log.Printf("[INFO] Skipping fcv3 trigger (%s): %s", triggerName, reason)
					continue
				}
				log.Printf("[INFO] Deleting fcv3 trigger: %s/%s", fnName, triggerName)
				deleteAction := fmt.Sprintf("/2023-03-30/functions/%s/triggers/%s", fnName, triggerName)
				if _, err := client.RoaDelete("FC", "2023-03-30", deleteAction, nil, nil, nil, true); err != nil {
					log.Printf("[ERROR] Failed to delete fcv3 trigger (%s/%s): %s", fnName, triggerName, err)
				}
			}
			nextToken, _ = response["nextToken"].(string)
			if nextToken == "" {
				break
			}
		}
	}
	return nil
}

// =============================================================================
// FC v3 — alicloud_fcv3_vpc_binding (per-function sub-resource)
// =============================================================================

func init() {
	resource.AddTestSweepers("alicloud_fcv3_vpc_binding", &resource.Sweeper{
		Name: "alicloud_fcv3_vpc_binding",
		F:    testSweepFcv3VpcBinding,
	})
}

func testSweepFcv3VpcBinding(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.FcNoSupportedRegions) {
		log.Printf("[INFO] Skipping Function Compute unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	functions, err := fcv3ListFunctions(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list fcv3 functions: %s", err)
		return nil
	}
	for _, fn := range functions {
		fnName := fmt.Sprint(fn["functionName"])
		if fnName == "" {
			continue
		}
		if del, _ := fcSweepShouldDelete(fnName); !del {
			continue
		}
		action := fmt.Sprintf("/2023-03-30/functions/%s/vpc-bindings", fnName)
		response, err := client.RoaGet("FC", "2023-03-30", action, nil, nil, nil)
		if err != nil {
			log.Printf("[ERROR] Failed to list fcv3 vpc-bindings for function (%s): %s", fnName, err)
			continue
		}
		// The vpc-bindings GET returns $.vpcIds as a list of vpcId strings.
		rawIds, _ := jsonpathGet("$.vpcIds", response)
		vpcIds, _ := rawIds.([]interface{})
		for _, v := range vpcIds {
			vpcId := fmt.Sprint(v)
			if vpcId == "" {
				continue
			}
			log.Printf("[INFO] Deleting fcv3 vpc-binding: %s/%s", fnName, vpcId)
			deleteAction := fmt.Sprintf("/2023-03-30/functions/%s/vpc-bindings/%s", fnName, vpcId)
			if _, err := client.RoaDelete("FC", "2023-03-30", deleteAction, nil, nil, nil, true); err != nil {
				log.Printf("[ERROR] Failed to delete fcv3 vpc-binding (%s/%s): %s", fnName, vpcId, err)
			}
		}
	}
	return nil
}

// isFcv3NotFoundError reports whether an FC v3 API error represents a
// not-found condition (so singleton-config deletions stay quiet when the
// config never existed). FC v3 surfaces these as 404-class errors carrying
// "...NotFound" error codes.
func isFcv3NotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return NotFoundError(err) || IsExpectedErrors(err, []string{
		"AsyncConfigNotExists",
		"FunctionNotFound",
		"AliasNotFound",
		"TriggerNotFound",
		"DomainNameNotFound",
		"LayerNotFound",
		"LayerVersionNotFound",
	})
}
