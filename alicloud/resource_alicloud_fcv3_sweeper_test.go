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

// FCv3 (Function Compute v3, product "FC" / "2023-03-30") child-resource
// sweepers. FC v3 is a ROA (RESTful path-based) API, so these use RoaGet for
// list and RoaDelete for delete (not RpcPost/RpcGet like gpdb/amqp). Query and
// header params are map[string]*string (built via StringPointer).
//
// Ten child resources are swept:
//   - region-scoped (list directly): function, custom_domain, layer_version
//   - function-scoped (walk functions first): alias, async_invoke_config,
//     concurrency_config, function_version, provision_config, trigger,
//     vpc_binding
//
// The function sweeper is safe to run (DeleteFunction is a plain PAYG delete,
// NOT a refund/退订 unlike gpdb/amqp instances). Deleting a function cascades
// its children, so the 7 function-scoped sweepers are defense-in-depth.
//
// custom_domain is a domain-bearing resource (like esa/DCDN): it uses an
// ALWAYS-active tf-prefix gate (no SWEEP_ALL bypass) so non-prefixed domains —
// which may be real存/ICP-backed — are preserved even in fleet-wipe mode
// (red-line #5). FC auto-generated fc.devsapp.net subdomains are not real存
// but are still preserved by the prefix heuristic (low-value leftover, safe).

var fcv3SweepPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf-example",
	"tf_example",
	"terraform-example",
	"tftestacc",
}

// fcv3GetField returns the first non-empty value among the given keys on item,
// trying both as-is (PascalCase) and lowercased to tolerate FC's mixed casing
// (functions list uses nextToken lowercase; triggers list uses NextToken).
func fcv3GetField(item map[string]interface{}, names ...string) string {
	for _, name := range names {
		if v, ok := item[name]; ok {
			if s := fmt.Sprint(v); s != "" && s != "<nil>" {
				return s
			}
		}
	}
	return ""
}

// fcv3SweepNameOk is the fleet-wipe gate for non-domain FCv3 resources:
// never touch NODELETING names, and in normal (non-SWEEP_ALL) mode only sweep
// names matching a tf-test prefix. SWEEP_ALL bypasses the prefix filter so the
// whole test account (1511928242963727, cn-beijing) is fleet-wiped.
func fcv3SweepNameOk(name string) bool {
	if name == "" {
		return false
	}
	if protectedSweepVpcName(name) {
		log.Printf("[INFO] Skipping fcv3 resource (NODELETING): %s", name)
		return false
	}
	if !sweepAll() {
		skip := true
		for _, prefix := range fcv3SweepPrefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping fcv3 resource (no prefix): %s", name)
			return false
		}
	}
	return true
}

// fcv3SweepDomainNameOk is the ALWAYS-active tf-prefix gate for custom_domain
// (domain-bearing resource, red-line #5). Even in SWEEP_ALL, only tf-prefixed
// domains are swept; non-prefixed domains are preserved (possible real存/ICP),
// matching the esa/DCDN pattern. The NODELETING guard stays always-active.
func fcv3SweepDomainNameOk(name string) bool {
	if name == "" {
		return false
	}
	if protectedSweepVpcName(name) {
		log.Printf("[INFO] Skipping fcv3 domain (NODELETING): %s", name)
		return false
	}
	skip := true
	for _, prefix := range fcv3SweepPrefixes {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			skip = false
			break
		}
	}
	if skip {
		log.Printf("[INFO] Skipping fcv3 domain (no prefix, preserved): %s", name)
		return false
	}
	return true
}

// fcv3RoaGetOne issues a single RoaGet with retry.Retry + NeedRetry +
// incrementalWait. query is map[string]*string (nil for path-param-only calls).
func fcv3RoaGetOne(client *connectivity.AliyunClient, action string, query map[string]*string) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("FC", "2023-03-30", action, query, nil, nil)
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

// fcv3DeletePath issues a RoaDelete on the given path with retry.
func fcv3DeletePath(client *connectivity.AliyunClient, action string) error {
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, err = client.RoaDelete("FC", "2023-03-30", action, nil, nil, nil, true)
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

// fcv3ListItems issues a single (non-paginated) RoaGet and returns the slice at
// itemsPath. A missing path is treated as an empty list (logged) so the sweeper
// continues with zero items rather than aborting.
func fcv3ListItems(client *connectivity.AliyunClient, action string, query map[string]*string, itemsPath string) ([]interface{}, error) {
	response, err := fcv3RoaGetOne(client, action, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", action, err)
	}
	resp, err := jsonpath.Get(itemsPath, response)
	if err != nil {
		log.Printf("[ERROR] fcv3 %s jsonpath %s: %s", action, itemsPath, err)
		return nil, nil
	}
	result, _ := resp.([]interface{})
	return result, nil
}

// fcv3ListStringItems is like fcv3ListItems but returns []string (for endpoints
// returning a flat string array, e.g. vpc-bindings returns $.vpcIds).
func fcv3ListStringItems(client *connectivity.AliyunClient, action string, query map[string]*string, itemsPath string) ([]string, error) {
	response, err := fcv3RoaGetOne(client, action, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", action, err)
	}
	resp, err := jsonpath.Get(itemsPath, response)
	if err != nil {
		log.Printf("[ERROR] fcv3 %s jsonpath %s: %s", action, itemsPath, err)
		return nil, nil
	}
	result, _ := resp.([]interface{})
	out := make([]string, 0, len(result))
	for _, v := range result {
		if s := fmt.Sprint(v); s != "" && s != "<nil>" {
			out = append(out, s)
		}
	}
	return out, nil
}

// fcv3ListFunctions enumerates every function in the region via
// /2023-03-30/functions (limit + nextToken pagination, lowercase nextToken).
func fcv3ListFunctions(client *connectivity.AliyunClient) ([]string, error) {
	action := "/2023-03-30/functions"
	functions := make([]string, 0)
	query := map[string]*string{"limit": StringPointer(fmt.Sprint(PageSizeXLarge))}
	for {
		response, err := fcv3RoaGetOne(client, action, query)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", action, err)
		}
		resp, err := jsonpath.Get("$.functions[*]", response)
		if err != nil {
			return nil, fmt.Errorf("%s jsonpath $.functions[*]: %s", action, err)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				if name := fcv3GetField(item, "functionName"); name != "" {
					functions = append(functions, name)
				}
			}
		}
		if nextToken, ok := response["nextToken"].(string); ok && nextToken != "" {
			query["nextToken"] = StringPointer(nextToken)
		} else {
			break
		}
	}
	return functions, nil
}

// fcv3ListLayers enumerates every layer via /2023-03-30/layers. Layer versions
// are then listed per-layer (the layer_version resource is layer-scoped).
func fcv3ListLayers(client *connectivity.AliyunClient) ([]string, error) {
	action := "/2023-03-30/layers"
	query := map[string]*string{"limit": StringPointer(fmt.Sprint(PageSizeXLarge))}
	items, err := fcv3ListItems(client, action, query, "$.layers[*]")
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(items))
	for _, v := range items {
		if item, ok := v.(map[string]interface{}); ok {
			if name := fcv3GetField(item, "layerName"); name != "" {
				out = append(out, name)
			}
		}
	}
	return out, nil
}

// fcv3GetConfig issues a RoaGet for a function-scoped single config (async
// invocation, concurrency, provision). Returns exists=false (no error) when the
// config is not set (the API returns a product-specific NotFound code).
func fcv3GetConfig(client *connectivity.AliyunClient, action string, notFoundCodes []string) (bool, error) {
	_, err := fcv3RoaGetOne(client, action, nil)
	if err != nil {
		if IsExpectedErrors(err, notFoundCodes) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func init() {
	resource.AddTestSweepers("alicloud_fcv3_function", &resource.Sweeper{
		Name: "alicloud_fcv3_function",
		F:    testSweepFcv3Functions,
	})
	resource.AddTestSweepers("alicloud_fcv3_custom_domain", &resource.Sweeper{
		Name: "alicloud_fcv3_custom_domain",
		F:    testSweepFcv3CustomDomains,
	})
	resource.AddTestSweepers("alicloud_fcv3_layer_version", &resource.Sweeper{
		Name: "alicloud_fcv3_layer_version",
		F:    testSweepFcv3LayerVersions,
	})
	resource.AddTestSweepers("alicloud_fcv3_alias", &resource.Sweeper{
		Name: "alicloud_fcv3_alias",
		F:    testSweepFcv3Aliases,
	})
	resource.AddTestSweepers("alicloud_fcv3_async_invoke_config", &resource.Sweeper{
		Name: "alicloud_fcv3_async_invoke_config",
		F:    testSweepFcv3AsyncInvokeConfigs,
	})
	resource.AddTestSweepers("alicloud_fcv3_concurrency_config", &resource.Sweeper{
		Name: "alicloud_fcv3_concurrency_config",
		F:    testSweepFcv3ConcurrencyConfigs,
	})
	resource.AddTestSweepers("alicloud_fcv3_function_version", &resource.Sweeper{
		Name: "alicloud_fcv3_function_version",
		F:    testSweepFcv3FunctionVersions,
	})
	resource.AddTestSweepers("alicloud_fcv3_provision_config", &resource.Sweeper{
		Name: "alicloud_fcv3_provision_config",
		F:    testSweepFcv3ProvisionConfigs,
	})
	resource.AddTestSweepers("alicloud_fcv3_trigger", &resource.Sweeper{
		Name: "alicloud_fcv3_trigger",
		F:    testSweepFcv3Triggers,
	})
	resource.AddTestSweepers("alicloud_fcv3_vpc_binding", &resource.Sweeper{
		Name: "alicloud_fcv3_vpc_binding",
		F:    testSweepFcv3VpcBindings,
	})
}

// testSweepFcv3Functions sweeps FCv3 functions (/functions -> DELETE /functions/%s).
// DeleteFunction is a plain PAYG delete (not a refund), so this is safe to run.
func testSweepFcv3Functions(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	functions, err := fcv3ListFunctions(aliyunClient)
	if err != nil {
		log.Printf("[ERROR] list fcv3 functions for function sweep: %s", err)
		return nil
	}
	for _, fn := range functions {
		if !fcv3SweepNameOk(fn) {
			continue
		}
		action := fmt.Sprintf("/2023-03-30/functions/%s", fn)
		if e := fcv3DeletePath(aliyunClient, action); e != nil {
			log.Printf("[ERROR] Failed to delete fcv3 function (%s): %s", fn, e)
			continue
		}
		log.Printf("[INFO] Delete fcv3 function success: %s", fn)
	}
	return nil
}

// testSweepFcv3CustomDomains sweeps FCv3 custom domains. Domain-bearing resource
// -> ALWAYS-active tf-prefix (esa/DCDN pattern); non-prefixed domains preserved
// (possible real存/ICP, red-line #5).
func testSweepFcv3CustomDomains(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	items, err := fcv3ListItems(aliyunClient, "/2023-03-30/custom-domains", map[string]*string{"limit": StringPointer(fmt.Sprint(PageSizeXLarge))}, "$.customDomains[*]")
	if err != nil {
		log.Printf("[ERROR] list fcv3 custom-domains: %s", err)
		return nil
	}
	for _, v := range items {
		item, _ := v.(map[string]interface{})
		name := fcv3GetField(item, "domainName")
		if !fcv3SweepDomainNameOk(name) {
			continue
		}
		action := fmt.Sprintf("/2023-03-30/custom-domains/%s", name)
		if e := fcv3DeletePath(aliyunClient, action); e != nil {
			log.Printf("[ERROR] Failed to delete fcv3 custom domain (%s): %s", name, e)
			continue
		}
		log.Printf("[INFO] Delete fcv3 custom domain success: %s", name)
	}
	return nil
}

// testSweepFcv3LayerVersions sweeps FCv3 layer versions. Layer-scoped: list
// layers first (/layers), then per layer list versions (/layers/%s/versions,
// response key $.layers[*], field "version"), then delete each version.
func testSweepFcv3LayerVersions(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	layers, err := fcv3ListLayers(aliyunClient)
	if err != nil {
		log.Printf("[ERROR] list fcv3 layers for layer version sweep: %s", err)
		return nil
	}
	for _, layer := range layers {
		if !fcv3SweepNameOk(layer) {
			continue
		}
		action := fmt.Sprintf("/2023-03-30/layers/%s/versions", layer)
		items, err := fcv3ListItems(aliyunClient, action, nil, "$.layers[*]")
		if err != nil {
			log.Printf("[ERROR] list fcv3 layer versions (%s): %s", layer, err)
			continue
		}
		for _, v := range items {
			item, _ := v.(map[string]interface{})
			ver := fcv3GetField(item, "version")
			if ver == "" {
				continue
			}
			delAction := fmt.Sprintf("/2023-03-30/layers/%s/versions/%s", layer, ver)
			if e := fcv3DeletePath(aliyunClient, delAction); e != nil {
				log.Printf("[ERROR] Failed to delete fcv3 layer version (%s/%s): %s", layer, ver, e)
				continue
			}
			log.Printf("[INFO] Delete fcv3 layer version success: %s/%s", layer, ver)
		}
	}
	return nil
}

// testSweepFcv3Aliases sweeps FCv3 aliases (function-scoped):
// /functions/%s/aliases -> DELETE /functions/%s/aliases/%s.
func testSweepFcv3Aliases(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	functions, err := fcv3ListFunctions(aliyunClient)
	if err != nil {
		log.Printf("[ERROR] list fcv3 functions for alias sweep: %s", err)
		return nil
	}
	for _, fn := range functions {
		if !fcv3SweepNameOk(fn) {
			continue
		}
		action := fmt.Sprintf("/2023-03-30/functions/%s/aliases", fn)
		items, err := fcv3ListItems(aliyunClient, action, nil, "$.aliases[*]")
		if err != nil {
			log.Printf("[ERROR] list fcv3 aliases (%s): %s", fn, err)
			continue
		}
		for _, v := range items {
			item, _ := v.(map[string]interface{})
			name := fcv3GetField(item, "aliasName")
			if !fcv3SweepNameOk(name) {
				continue
			}
			delAction := fmt.Sprintf("/2023-03-30/functions/%s/aliases/%s", fn, name)
			if e := fcv3DeletePath(aliyunClient, delAction); e != nil {
				log.Printf("[ERROR] Failed to delete fcv3 alias (%s/%s): %s", fn, name, e)
				continue
			}
			log.Printf("[INFO] Delete fcv3 alias success: %s/%s", fn, name)
		}
	}
	return nil
}

// testSweepFcv3AsyncInvokeConfigs sweeps FCv3 async invocation configs
// (function-scoped single config): GET /functions/%s/async-invoke-config; if set,
// DELETE the same path. NotFound code: AsyncConfigNotExists.
func testSweepFcv3AsyncInvokeConfigs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	functions, err := fcv3ListFunctions(aliyunClient)
	if err != nil {
		log.Printf("[ERROR] list fcv3 functions for async invoke config sweep: %s", err)
		return nil
	}
	for _, fn := range functions {
		if !fcv3SweepNameOk(fn) {
			continue
		}
		action := fmt.Sprintf("/2023-03-30/functions/%s/async-invoke-config", fn)
		exists, e := fcv3GetConfig(aliyunClient, action, []string{"AsyncConfigNotExists"})
		if e != nil {
			log.Printf("[ERROR] get fcv3 async invoke config (%s): %s", fn, e)
			continue
		}
		if !exists {
			continue
		}
		if de := fcv3DeletePath(aliyunClient, action); de != nil {
			log.Printf("[ERROR] Failed to delete fcv3 async invoke config (%s): %s", fn, de)
			continue
		}
		log.Printf("[INFO] Delete fcv3 async invoke config success: %s", fn)
	}
	return nil
}

// testSweepFcv3ConcurrencyConfigs sweeps FCv3 concurrency configs
// (function-scoped single config): GET /functions/%s/concurrency; if set, DELETE.
// NotFound code: OnDemandConfigNotFound.
func testSweepFcv3ConcurrencyConfigs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	functions, err := fcv3ListFunctions(aliyunClient)
	if err != nil {
		log.Printf("[ERROR] list fcv3 functions for concurrency config sweep: %s", err)
		return nil
	}
	for _, fn := range functions {
		if !fcv3SweepNameOk(fn) {
			continue
		}
		action := fmt.Sprintf("/2023-03-30/functions/%s/concurrency", fn)
		exists, e := fcv3GetConfig(aliyunClient, action, []string{"OnDemandConfigNotFound"})
		if e != nil {
			log.Printf("[ERROR] get fcv3 concurrency config (%s): %s", fn, e)
			continue
		}
		if !exists {
			continue
		}
		if de := fcv3DeletePath(aliyunClient, action); de != nil {
			log.Printf("[ERROR] Failed to delete fcv3 concurrency config (%s): %s", fn, de)
			continue
		}
		log.Printf("[INFO] Delete fcv3 concurrency config success: %s", fn)
	}
	return nil
}

// testSweepFcv3FunctionVersions sweeps FCv3 function versions (function-scoped):
// /functions/%s/versions -> DELETE /functions/%s/versions/%s.
func testSweepFcv3FunctionVersions(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	functions, err := fcv3ListFunctions(aliyunClient)
	if err != nil {
		log.Printf("[ERROR] list fcv3 functions for function version sweep: %s", err)
		return nil
	}
	for _, fn := range functions {
		if !fcv3SweepNameOk(fn) {
			continue
		}
		action := fmt.Sprintf("/2023-03-30/functions/%s/versions", fn)
		items, err := fcv3ListItems(aliyunClient, action, nil, "$.versions[*]")
		if err != nil {
			log.Printf("[ERROR] list fcv3 function versions (%s): %s", fn, err)
			continue
		}
		for _, v := range items {
			item, _ := v.(map[string]interface{})
			ver := fcv3GetField(item, "version")
			if ver == "" {
				continue
			}
			delAction := fmt.Sprintf("/2023-03-30/functions/%s/versions/%s", fn, ver)
			if e := fcv3DeletePath(aliyunClient, delAction); e != nil {
				log.Printf("[ERROR] Failed to delete fcv3 function version (%s/%s): %s", fn, ver, e)
				continue
			}
			log.Printf("[INFO] Delete fcv3 function version success: %s/%s", fn, ver)
		}
	}
	return nil
}

// testSweepFcv3ProvisionConfigs sweeps FCv3 provision configs (function-scoped
// single config): GET /functions/%s/provision-config; if set, DELETE.
func testSweepFcv3ProvisionConfigs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	functions, err := fcv3ListFunctions(aliyunClient)
	if err != nil {
		log.Printf("[ERROR] list fcv3 functions for provision config sweep: %s", err)
		return nil
	}
	for _, fn := range functions {
		if !fcv3SweepNameOk(fn) {
			continue
		}
		action := fmt.Sprintf("/2023-03-30/functions/%s/provision-config", fn)
		exists, e := fcv3GetConfig(aliyunClient, action, []string{"ProvisionConfigNotFound", "NotFound", "ResourceNotFound"})
		if e != nil {
			log.Printf("[ERROR] get fcv3 provision config (%s): %s", fn, e)
			continue
		}
		if !exists {
			continue
		}
		if de := fcv3DeletePath(aliyunClient, action); de != nil {
			log.Printf("[ERROR] Failed to delete fcv3 provision config (%s): %s", fn, de)
			continue
		}
		log.Printf("[INFO] Delete fcv3 provision config success: %s", fn)
	}
	return nil
}

// testSweepFcv3Triggers sweeps FCv3 triggers (function-scoped):
// /functions/%s/triggers -> DELETE /functions/%s/triggers/%s.
func testSweepFcv3Triggers(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	functions, err := fcv3ListFunctions(aliyunClient)
	if err != nil {
		log.Printf("[ERROR] list fcv3 functions for trigger sweep: %s", err)
		return nil
	}
	for _, fn := range functions {
		if !fcv3SweepNameOk(fn) {
			continue
		}
		action := fmt.Sprintf("/2023-03-30/functions/%s/triggers", fn)
		items, err := fcv3ListItems(aliyunClient, action, nil, "$.triggers[*]")
		if err != nil {
			log.Printf("[ERROR] list fcv3 triggers (%s): %s", fn, err)
			continue
		}
		for _, v := range items {
			item, _ := v.(map[string]interface{})
			name := fcv3GetField(item, "triggerName")
			if !fcv3SweepNameOk(name) {
				continue
			}
			delAction := fmt.Sprintf("/2023-03-30/functions/%s/triggers/%s", fn, name)
			if e := fcv3DeletePath(aliyunClient, delAction); e != nil {
				log.Printf("[ERROR] Failed to delete fcv3 trigger (%s/%s): %s", fn, name, e)
				continue
			}
			log.Printf("[INFO] Delete fcv3 trigger success: %s/%s", fn, name)
		}
	}
	return nil
}

// testSweepFcv3VpcBindings sweeps FCv3 vpc bindings (function-scoped):
// /functions/%s/vpc-bindings returns $.vpcIds (string array);
// DELETE /functions/%s/vpc-bindings/%s.
func testSweepFcv3VpcBindings(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	functions, err := fcv3ListFunctions(aliyunClient)
	if err != nil {
		log.Printf("[ERROR] list fcv3 functions for vpc binding sweep: %s", err)
		return nil
	}
	for _, fn := range functions {
		if !fcv3SweepNameOk(fn) {
			continue
		}
		action := fmt.Sprintf("/2023-03-30/functions/%s/vpc-bindings", fn)
		vpcIDs, err := fcv3ListStringItems(aliyunClient, action, nil, "$.vpcIds")
		if err != nil {
			log.Printf("[ERROR] list fcv3 vpc bindings (%s): %s", fn, err)
			continue
		}
		for _, vpcID := range vpcIDs {
			delAction := fmt.Sprintf("/2023-03-30/functions/%s/vpc-bindings/%s", fn, vpcID)
			if e := fcv3DeletePath(aliyunClient, delAction); e != nil {
				log.Printf("[ERROR] Failed to delete fcv3 vpc binding (%s/%s): %s", fn, vpcID, e)
				continue
			}
			log.Printf("[INFO] Delete fcv3 vpc binding success: %s/%s", fn, vpcID)
		}
	}
	return nil
}
