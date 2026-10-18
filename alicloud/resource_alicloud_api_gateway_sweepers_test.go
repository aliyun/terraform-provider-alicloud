package alicloud

// Consolidated sweepers for the ApiGateway (CloudAPI) product. These mirror
// the Delete path of each resource's provider implementation and reuse the
// shared name-based protection (protectedSweepVpcName) so SWEEP_ALL is safe.
//
// Child resources are swept before parents:
//   - alicloud_api_gateway_api / alicloud_api_gateway_model run before
//     alicloud_api_gateway_group (group depends on them).
//   - alicloud_api_gateway_group runs before alicloud_api_gateway_instance.
//
// Attachment resources (acl_entry_attachment, app_attachment, plugin_attachment,
// group_plugin_attachment, instance_acl_attachment) and backend_model are
// relationship/child objects that cascade from their parent's Delete, so they
// are not given standalone sweepers here.

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

const (
	apiGatewayProduct  = "CloudAPI"
	apiGatewayVersion  = "2016-07-14"
	apiGatewaySweepLog = "[SWEEPER] api_gateway"
)

// apiGatewaySweepPrefixes are the test-name prefixes matched when SWEEP_ALL is
// not enabled. In SWEEP_ALL mode the prefix filter is bypassed and only the
// name-based NODELETING protection (protectedSweepVpcName) guards resources.
var apiGatewaySweepPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf_test_",
	"tf-test-",
	"testAcc",
	"tf-example",
	"tf-test",
}

// apiGatewayNameMatchesPrefix reports whether name carries a test prefix.
func apiGatewayNameMatchesPrefix(name string) bool {
	ln := strings.ToLower(name)
	for _, p := range apiGatewaySweepPrefixes {
		if strings.HasPrefix(ln, strings.ToLower(p)) {
			return true
		}
	}
	return false
}

// apiGatewayIsBuiltInStageName reports whether name is one of the three
// system-default stage environments (RELEASE/PRE/TEST). These are created by
// the platform and must never be deleted; they are not produced by tests.
func apiGatewayIsBuiltInStageName(name string) bool {
	switch strings.ToUpper(strings.TrimSpace(name)) {
	case "RELEASE", "PRE", "TEST":
		return true
	}
	return false
}

// apiGatewayIsSharedInstance reports whether the instance is the
// platform-managed shared instance (api-shared-* id, or name carries 共享/shared).
// These cannot be deleted by the account owner and are skipped pre-emptively.
func apiGatewayIsSharedInstance(id, name string) bool {
	if strings.HasPrefix(id, "api-shared") {
		return true
	}
	ln := strings.ToLower(name)
	if strings.Contains(ln, "shared") || strings.Contains(name, "共享") {
		return true
	}
	return false
}

// apiGatewaySweepNameSkip is the shared gate: NODELETING/do-not-delete markers
// always skip; when SWEEP_ALL is off, only test-prefixed names are swept.
func apiGatewaySweepNameSkip(name string) (bool, string) {
	if protectedSweepVpcName(name) {
		return true, "do-not-delete marker in name"
	}
	if !sweepAll() && !apiGatewayNameMatchesPrefix(name) {
		return true, "name does not match test prefixes"
	}
	return false, ""
}

// apiGatewayAppIDString normalises the AppId field returned by DescribeAppAttributes
// (a JSON number -> float64) into the canonical integer string used by DeleteApp.
func apiGatewayAppIDString(v interface{}) string {
	switch x := v.(type) {
	case string:
		return x
	case float64:
		return strconv.FormatInt(int64(x), 10)
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	default:
		return fmt.Sprint(v)
	}
}

// apiGatewayList calls a CloudAPI list action with pagination and returns the
// raw item slice at the given jsonpath. Per-page failures are logged and
// return an error so callers can abort the list.
func apiGatewayList(client *connectivity.AliyunClient, action, path string, base map[string]interface{}) ([]map[string]interface{}, error) {
	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	for k, v := range base {
		request[k] = v
	}
	out := make([]map[string]interface{}, 0)
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		var response map[string]interface{}
		var err error
		if err := retry.Retry(3*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost(apiGatewayProduct, apiGatewayVersion, action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		}); err != nil {
			addDebug(action, response, request)
			return nil, fmt.Errorf("%s: %w", action, err)
		}
		addDebug(action, response, request)
		raw, perr := jsonpath.Get(path, response)
		if perr != nil {
			return nil, fmt.Errorf("%s response: %w", action, perr)
		}
		rows, _ := raw.([]interface{})
		for _, r := range rows {
			if m, ok := r.(map[string]interface{}); ok {
				out = append(out, m)
			}
		}
		if len(rows) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return out, nil
}

// apiGatewayDelete issues a CloudAPI delete action. Errors are returned to the
// caller, which logs and continues (sweeper returns nil even on per-item fail).
func apiGatewayDelete(client *connectivity.AliyunClient, action string, query, body map[string]interface{}) error {
	wait := incrementalWait(3*time.Second, 5*time.Second)
	var response map[string]interface{}
	var err error
	if err := retry.Retry(3*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost(apiGatewayProduct, apiGatewayVersion, action, query, body, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		addDebug(action, response, body)
		return err
	}
	addDebug(action, response, body)
	return nil
}

func init() {
	resource.AddTestSweepers("alicloud_api_gateway_instance", &resource.Sweeper{
		Name: "alicloud_api_gateway_instance",
		Dependencies: []string{
			"alicloud_api_gateway_group",
		},
		F: testSweepApiGatewayInstances,
	})
	resource.AddTestSweepers("alicloud_api_gateway_group", &resource.Sweeper{
		Name: "alicloud_api_gateway_group",
		Dependencies: []string{
			"alicloud_api_gateway_api",
			"alicloud_api_gateway_model",
		},
		F: testSweepApiGatewayGroups,
	})
	resource.AddTestSweepers("alicloud_api_gateway_api", &resource.Sweeper{
		Name: "alicloud_api_gateway_api",
		F:   testSweepApiGatewayApis,
	})
	resource.AddTestSweepers("alicloud_api_gateway_app", &resource.Sweeper{
		Name: "alicloud_api_gateway_app",
		F:   testSweepApiGatewayApps,
	})
	resource.AddTestSweepers("alicloud_api_gateway_backend", &resource.Sweeper{
		Name: "alicloud_api_gateway_backend",
		F:   testSweepApiGatewayBackends,
	})
	resource.AddTestSweepers("alicloud_api_gateway_plugin", &resource.Sweeper{
		Name: "alicloud_api_gateway_plugin",
		F:   testSweepApiGatewayPlugins,
	})
	resource.AddTestSweepers("alicloud_api_gateway_vpc_access", &resource.Sweeper{
		Name: "alicloud_api_gateway_vpc_access",
		F:   testSweepApiGatewayVpcAccesses,
	})
	resource.AddTestSweepers("alicloud_api_gateway_access_control_list", &resource.Sweeper{
		Name: "alicloud_api_gateway_access_control_list",
		F:   testSweepApiGatewayAccessControlLists,
	})
	resource.AddTestSweepers("alicloud_api_gateway_log_config", &resource.Sweeper{
		Name: "alicloud_api_gateway_log_config",
		F:   testSweepApiGatewayLogConfigs,
	})
	resource.AddTestSweepers("alicloud_api_gateway_stage_model", &resource.Sweeper{
		Name: "alicloud_api_gateway_stage_model",
		F:   testSweepApiGatewayStageModels,
	})
	resource.AddTestSweepers("alicloud_api_gateway_model", &resource.Sweeper{
		Name: "alicloud_api_gateway_model",
		F:   testSweepApiGatewayModels,
	})
}

// testSweepApiGatewayInstances mirrors resourceAliCloudApiGatewayInstanceDelete:
// DescribeInstances -> DeleteInstance(InstanceId). The instance name is the
// human-readable field used for NODELETING/prefix gating.
func testSweepApiGatewayInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := apiGatewayList(client, "DescribeInstances", "$.Instances.InstanceAttribute[*]", map[string]interface{}{})
	if err != nil {
		log.Printf("[ERROR] %s DescribeInstances failed: %v", apiGatewaySweepLog, err)
		return nil
	}
	for _, item := range items {
		id := fmt.Sprint(item["InstanceId"])
		name := fmt.Sprint(item["InstanceName"])
		// Shared instances (api-shared-*, name contains 共享/shared) are
		// platform-managed and cannot be deleted by the account owner; skip
		// them rather than hammering DeleteInstance against a managed resource.
		if apiGatewayIsSharedInstance(id, name) {
			log.Printf("[INFO] Skipping ApiGateway Instance %s (%s): shared/managed instance", name, id)
			continue
		}
		if skip, reason := apiGatewaySweepNameSkip(name); skip {
			log.Printf("[INFO] Skipping ApiGateway Instance %s (%s): %s", name, id, reason)
			continue
		}
		log.Printf("[INFO] Deleting ApiGateway Instance: %s (%s)", name, id)
		if err := apiGatewayDelete(client, "DeleteInstance", map[string]interface{}{"InstanceId": id}, map[string]interface{}{}); err != nil {
			log.Printf("[ERROR] Failed to delete ApiGateway Instance %s (%s): %v", name, id, err)
			continue
		}
	}
	return nil
}

// testSweepApiGatewayGroups mirrors resourceAliyunApigatewayGroupDelete:
// DescribeApiGroups -> DeleteApiGroup(GroupId).
func testSweepApiGatewayGroups(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := apiGatewayList(client, "DescribeApiGroups", "$.ApiGroupAttributes.ApiGroupAttribute[*]", map[string]interface{}{})
	if err != nil {
		log.Printf("[ERROR] %s DescribeApiGroups failed: %v", apiGatewaySweepLog, err)
		return nil
	}
	for _, item := range items {
		id := fmt.Sprint(item["GroupId"])
		name := fmt.Sprint(item["GroupName"])
		if skip, reason := apiGatewaySweepNameSkip(name); skip {
			log.Printf("[INFO] Skipping ApiGateway Group %s (%s): %s", name, id, reason)
			continue
		}
		log.Printf("[INFO] Deleting ApiGateway Group: %s (%s)", name, id)
		if err := apiGatewayDelete(client, "DeleteApiGroup", map[string]interface{}{"GroupId": id}, map[string]interface{}{}); err != nil {
			log.Printf("[ERROR] Failed to delete ApiGateway Group %s (%s): %v", name, id, err)
			continue
		}
	}
	return nil
}

// testSweepApiGatewayApis mirrors resourceAliyunApigatewayApiDelete:
// DescribeApis -> DeleteApi(ApiId, GroupId). The composite resource id is
// group_id:api_id, but the list response carries both fields directly.
func testSweepApiGatewayApis(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := apiGatewayList(client, "DescribeApis", "$.ApiSummarys.ApiSummary[*]", map[string]interface{}{})
	if err != nil {
		log.Printf("[ERROR] %s DescribeApis failed: %v", apiGatewaySweepLog, err)
		return nil
	}
	for _, item := range items {
		apiId := fmt.Sprint(item["ApiId"])
		groupId := fmt.Sprint(item["GroupId"])
		name := fmt.Sprint(item["ApiName"])
		id := groupId + ":" + apiId
		if skip, reason := apiGatewaySweepNameSkip(name); skip {
			log.Printf("[INFO] Skipping ApiGateway Api %s (%s): %s", name, id, reason)
			continue
		}
		log.Printf("[INFO] Deleting ApiGateway Api: %s (%s)", name, id)
		if err := apiGatewayDelete(client, "DeleteApi", map[string]interface{}{"ApiId": apiId, "GroupId": groupId}, map[string]interface{}{}); err != nil {
			log.Printf("[ERROR] Failed to delete ApiGateway Api %s (%s): %v", name, id, err)
			continue
		}
	}
	return nil
}

// testSweepApiGatewayApps mirrors resourceAliyunApigatewayAppDelete:
// DescribeAppAttributes -> DeleteApp(AppId). AppId is a JSON number, so it is
// normalised to the canonical integer string.
func testSweepApiGatewayApps(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := apiGatewayList(client, "DescribeAppAttributes", "$.Apps.AppAttribute[*]", map[string]interface{}{})
	if err != nil {
		log.Printf("[ERROR] %s DescribeAppAttributes failed: %v", apiGatewaySweepLog, err)
		return nil
	}
	for _, item := range items {
		id := apiGatewayAppIDString(item["AppId"])
		name := fmt.Sprint(item["AppName"])
		if skip, reason := apiGatewaySweepNameSkip(name); skip {
			log.Printf("[INFO] Skipping ApiGateway App %s (%s): %s", name, id, reason)
			continue
		}
		log.Printf("[INFO] Deleting ApiGateway App: %s (%s)", name, id)
		if err := apiGatewayDelete(client, "DeleteApp", map[string]interface{}{"AppId": id}, map[string]interface{}{}); err != nil {
			log.Printf("[ERROR] Failed to delete ApiGateway App %s (%s): %v", name, id, err)
			continue
		}
	}
	return nil
}

// testSweepApiGatewayBackends mirrors resourceAlicloudApiGatewayBackendDelete:
// DescribeBackendList -> DeleteBackend(BackendId).
func testSweepApiGatewayBackends(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := apiGatewayList(client, "DescribeBackendList", "$.BackendInfoList[*]", map[string]interface{}{})
	if err != nil {
		log.Printf("[ERROR] %s DescribeBackendList failed: %v", apiGatewaySweepLog, err)
		return nil
	}
	for _, item := range items {
		id := fmt.Sprint(item["BackendId"])
		name := fmt.Sprint(item["BackendName"])
		if skip, reason := apiGatewaySweepNameSkip(name); skip {
			log.Printf("[INFO] Skipping ApiGateway Backend %s (%s): %s", name, id, reason)
			continue
		}
		log.Printf("[INFO] Deleting ApiGateway Backend: %s (%s)", name, id)
		if err := apiGatewayDelete(client, "DeleteBackend", nil, map[string]interface{}{"BackendId": id}); err != nil {
			log.Printf("[ERROR] Failed to delete ApiGateway Backend %s (%s): %v", name, id, err)
			continue
		}
	}
	return nil
}

// testSweepApiGatewayPlugins mirrors resourceAliCloudApiGatewayPluginDelete:
// DescribePlugins -> DeletePlugin(PluginId).
func testSweepApiGatewayPlugins(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := apiGatewayList(client, "DescribePlugins", "$.Plugins.PluginAttribute[*]", map[string]interface{}{})
	if err != nil {
		log.Printf("[ERROR] %s DescribePlugins failed: %v", apiGatewaySweepLog, err)
		return nil
	}
	for _, item := range items {
		id := fmt.Sprint(item["PluginId"])
		name := fmt.Sprint(item["PluginName"])
		if skip, reason := apiGatewaySweepNameSkip(name); skip {
			log.Printf("[INFO] Skipping ApiGateway Plugin %s (%s): %s", name, id, reason)
			continue
		}
		log.Printf("[INFO] Deleting ApiGateway Plugin: %s (%s)", name, id)
		if err := apiGatewayDelete(client, "DeletePlugin", map[string]interface{}{"PluginId": id}, map[string]interface{}{}); err != nil {
			log.Printf("[ERROR] Failed to delete ApiGateway Plugin %s (%s): %v", name, id, err)
			continue
		}
	}
	return nil
}

// testSweepApiGatewayVpcAccesses mirrors resourceAliCloudApiGatewayVpcAccessDelete:
// DescribeVpcAccesses -> RemoveVpcAccess(VpcId, InstanceId, Port, Name). The
// composite resource id is name:vpc_id:instance_id:port, but the delete call
// uses the individual fields carried by the list item.
func testSweepApiGatewayVpcAccesses(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := apiGatewayList(client, "DescribeVpcAccesses", "$.VpcAccessAttributes.VpcAccessAttribute[*]", map[string]interface{}{})
	if err != nil {
		log.Printf("[ERROR] %s DescribeVpcAccesses failed: %v", apiGatewaySweepLog, err)
		return nil
	}
	for _, item := range items {
		name := fmt.Sprint(item["Name"])
		vpcId := fmt.Sprint(item["VpcId"])
		instanceId := fmt.Sprint(item["InstanceId"])
		port := fmt.Sprint(item["Port"])
		id := name + ":" + vpcId + ":" + instanceId + ":" + port
		if skip, reason := apiGatewaySweepNameSkip(name); skip {
			log.Printf("[INFO] Skipping ApiGateway VpcAccess %s (%s): %s", name, id, reason)
			continue
		}
		log.Printf("[INFO] Deleting ApiGateway VpcAccess: %s (%s)", name, id)
		body := map[string]interface{}{
			"VpcId":      vpcId,
			"InstanceId": instanceId,
			"Port":       item["Port"],
			"Name":       name,
		}
		if err := apiGatewayDelete(client, "RemoveVpcAccess", nil, body); err != nil {
			log.Printf("[ERROR] Failed to delete ApiGateway VpcAccess %s (%s): %v", name, id, err)
			continue
		}
	}
	return nil
}

// testSweepApiGatewayAccessControlLists mirrors resourceAliCloudApiGatewayAccessControlListDelete:
// DescribeAccessControlLists -> DeleteAccessControlList(AclId).
func testSweepApiGatewayAccessControlLists(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := apiGatewayList(client, "DescribeAccessControlLists", "$.Acls.Acl[*]", map[string]interface{}{})
	if err != nil {
		log.Printf("[ERROR] %s DescribeAccessControlLists failed: %v", apiGatewaySweepLog, err)
		return nil
	}
	for _, item := range items {
		id := fmt.Sprint(item["AclId"])
		name := fmt.Sprint(item["AclName"])
		if skip, reason := apiGatewaySweepNameSkip(name); skip {
			log.Printf("[INFO] Skipping ApiGateway AccessControlList %s (%s): %s", name, id, reason)
			continue
		}
		log.Printf("[INFO] Deleting ApiGateway AccessControlList: %s (%s)", name, id)
		if err := apiGatewayDelete(client, "DeleteAccessControlList", map[string]interface{}{"AclId": id}, map[string]interface{}{}); err != nil {
			log.Printf("[ERROR] Failed to delete ApiGateway AccessControlList %s (%s): %v", name, id, err)
			continue
		}
	}
	return nil
}

// testSweepApiGatewayLogConfigs mirrors resourceAlicloudApiGatewayLogConfigDelete:
// DescribeLogConfig -> DeleteLogConfig(LogType). The LogType value is both the
// resource id and the only identifier exposed by the list API.
func testSweepApiGatewayLogConfigs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	// DescribeLogConfig lists all log configs when LogType is not supplied.
	items, err := apiGatewayList(client, "DescribeLogConfig", "$.LogInfos.LogInfo", map[string]interface{}{})
	if err != nil {
		log.Printf("[ERROR] %s DescribeLogConfig failed: %v", apiGatewaySweepLog, err)
		return nil
	}
	for _, item := range items {
		logType := fmt.Sprint(item["LogType"])
		if logType == "" {
			continue
		}
		// LogType is an opaque hash, not a human name; only the NODELETING
		// marker gate applies, plus the prefix gate in non-SWEEP_ALL mode.
		if skip, reason := apiGatewaySweepNameSkip(logType); skip {
			log.Printf("[INFO] Skipping ApiGateway LogConfig (%s): %s", logType, reason)
			continue
		}
		log.Printf("[INFO] Deleting ApiGateway LogConfig: (%s)", logType)
		if err := apiGatewayDelete(client, "DeleteLogConfig", nil, map[string]interface{}{"LogType": logType}); err != nil {
			log.Printf("[ERROR] Failed to delete ApiGateway LogConfig (%s): %v", logType, err)
			continue
		}
	}
	return nil
}

// testSweepApiGatewayStageModels mirrors resourceAlicloudApiGatewayStageModelDelete:
// DescribeStageModels -> DeleteStageModel(StageModelId). The StageName is used
// as the human-readable gate; if absent the opaque StageModelId gates.
func testSweepApiGatewayStageModels(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := apiGatewayList(client, "DescribeStageModels", "$.StageModelInfoList[*]", map[string]interface{}{})
	if err != nil {
		log.Printf("[ERROR] %s DescribeStageModels failed: %v", apiGatewaySweepLog, err)
		return nil
	}
	for _, item := range items {
		id := fmt.Sprint(item["StageModelId"])
		name := fmt.Sprint(item["StageName"])
		if name == "" {
			name = id
		}
		// The three system-default stage environments (RELEASE/PRE/TEST) are
		// platform-managed and must never be deleted, even in SWEEP_ALL mode.
		if apiGatewayIsBuiltInStageName(name) {
			log.Printf("[INFO] Skipping ApiGateway StageModel %s (%s): built-in default stage", name, id)
			continue
		}
		if skip, reason := apiGatewaySweepNameSkip(name); skip {
			log.Printf("[INFO] Skipping ApiGateway StageModel %s (%s): %s", name, id, reason)
			continue
		}
		log.Printf("[INFO] Deleting ApiGateway StageModel: %s (%s)", name, id)
		if err := apiGatewayDelete(client, "DeleteStageModel", nil, map[string]interface{}{"StageModelId": id}); err != nil {
			log.Printf("[ERROR] Failed to delete ApiGateway StageModel %s (%s): %v", name, id, err)
			continue
		}
	}
	return nil
}

// testSweepApiGatewayModels mirrors resourceAlicloudApiGatewayModelDelete:
// DescribeModels(Group) -> DeleteModel(GroupId, ModelName). Models are scoped
// to a group, so the sweeper first lists groups, then lists models per group.
// This runs before the group sweeper (group depends on it), so groups still
// exist when models are enumerated.
func testSweepApiGatewayModels(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	groups, err := apiGatewayList(client, "DescribeApiGroups", "$.ApiGroupAttributes.ApiGroupAttribute[*]", map[string]interface{}{})
	if err != nil {
		log.Printf("[ERROR] %s DescribeApiGroups (for models) failed: %v", apiGatewaySweepLog, err)
		return nil
	}
	for _, g := range groups {
		groupId := fmt.Sprint(g["GroupId"])
		groupName := fmt.Sprint(g["GroupName"])
		// If the group itself is protected/non-test, skip its models too so we
		// never delete a model that belongs to a preserved group.
		if skip, reason := apiGatewaySweepNameSkip(groupName); skip {
			log.Printf("[INFO] Skipping ApiGateway Models under group %s (%s): %s", groupName, groupId, reason)
			continue
		}
		models, err := apiGatewayList(client, "DescribeModels", "$.ModelDetails.ModelDetail[*]", map[string]interface{}{"GroupId": groupId})
		if err != nil {
			log.Printf("[ERROR] %s DescribeModels for group %s failed: %v", apiGatewaySweepLog, groupId, err)
			continue
		}
		for _, m := range models {
			modelName := fmt.Sprint(m["ModelName"])
			id := groupId + ":" + modelName
			if skip, reason := apiGatewaySweepNameSkip(modelName); skip {
				log.Printf("[INFO] Skipping ApiGateway Model %s (%s): %s", modelName, id, reason)
				continue
			}
			log.Printf("[INFO] Deleting ApiGateway Model: %s (%s)", modelName, id)
			if err := apiGatewayDelete(client, "DeleteModel", nil, map[string]interface{}{"GroupId": groupId, "ModelName": modelName}); err != nil {
				log.Printf("[ERROR] Failed to delete ApiGateway Model %s (%s): %v", modelName, id, err)
				continue
			}
		}
	}
	return nil
}
