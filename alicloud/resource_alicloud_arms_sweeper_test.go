package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// ARMS (Application Real-Time Monitoring Service, product "ARMS" / "2019-08-08")
// child-resource sweepers — fills the gaps left by the 4 pre-existing inline
// sweepers (alert_contact, alert_contact_group, alert_robot, dispatch_rule in
// their own _test.go files). ARMS is RPC-style: RpcPost with
// map[string]interface{} params, and — like the existing ARMS sweepers — the
// params go in the BODY slot (RpcPost query=nil, body=params). The service
// structs (ArmsService / ArmsServiceV2) only expose Describe* (Get-by-id), so
// sweepers call the list/search APIs directly.
//
// Twelve gap sweepers are added here (one — remote_write — has NO list API and
// is documented unsweepable, like grey_tag_route):
//   region-scoped: environment, grafana_workspace, prometheus, synthetic_task
//   environment-scoped (walk environments first): addon_release, env_custom_job,
//     env_feature, env_pod_monitor, env_service_monitor
//   prometheus/cluster-scoped (walk prometheus instances first):
//     integration_exporter, prometheus_alert_rule, prometheus_monitoring
//
// prometheus is the HIGH-VALUE sweeper: ListPrometheusInstances returns PAYG
// test-leftover instances (code-test-tf-variable-N, tfaccrealtimecomputeNNNN —
// note these are NOT tf-prefixed, so the pre-existing prefix-only sweepers would
// skip them; the sweepAll() fleet-wipe gate here bypasses the prefix and cleans
// them) via UninstallPromCluster (a plain PAYG delete, payMethod=post_pay_gb,
// NOT a BSS refund/退订 — verified: returns {Data:"success",Code:200}). The
// environment sweeper filters RegionId==region (red-line: ListEnvironments is a
// GLOBAL list — test-account leftovers happen to be in cn-hangzhou, must NOT be
// cleaned from a cn-beijing sweep). synthetic_task is unsupported in cn-beijing
// (ParameterRegionIllegal) — graceful no-op. integration_exporter only supports
// aliyun-cs/ecs instances (remote-write-prometheus skipped — avoids 7×N useless
// ListPrometheusIntegration calls).

var armsSweepPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf-example",
	"tf_example",
	"terraform-example",
	"tftestacc",
}

// armsGetField returns the first non-empty value among the given keys on item,
// trying both as-is and lowercased (ARMS response casing is inconsistent).
func armsGetField(item map[string]interface{}, names ...string) string {
	for _, name := range names {
		if v, ok := item[name]; ok {
			if s := fmt.Sprint(v); s != "" && s != "<nil>" {
				return s
			}
		}
	}
	return ""
}

// armsSweepNameOk is the fleet-wipe gate: never touch NODELETING names; in
// normal (non-SWEEP_ALL) mode only sweep tf-prefixed names. SWEEP_ALL bypasses
// the prefix filter so the whole test account (1511928242963727, cn-beijing) is
// fleet-wiped. This matters for prometheus: test leftovers are named
// code-test-tf-variable-N / tfaccrealtimecomputeNNNN (NOT tf-prefixed), so they
// can only be cleaned in SWEEP_ALL mode.
func armsSweepNameOk(name string) bool {
	if name == "" {
		return false
	}
	if protectedSweepVpcName(name) {
		log.Printf("[INFO] Skipping arms resource (NODELETING): %s", name)
		return false
	}
	if !sweepAll() {
		skip := true
		for _, prefix := range armsSweepPrefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping arms resource (no prefix): %s", name)
			return false
		}
	}
	return true
}

// armsRpcPost issues an RpcPost with retry.Retry + NeedRetry +
// incrementalWait. Params go in the BODY slot (query=nil), matching the
// pre-existing ARMS sweepers (RpcPost signature: product, version, action,
// query, body, autoRetry).
func armsRpcPost(client *connectivity.AliyunClient, action string, params map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, params, true)
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

// armsDeletePost issues an RpcPost delete with retry (params in body slot).
func armsDeletePost(client *connectivity.AliyunClient, action string, params map[string]interface{}) error {
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, err = client.RpcPost("ARMS", "2019-08-08", action, nil, params, true)
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

// armsNavigate walks a dotted path (e.g. "Data.Environments", "PageBean.Contacts")
// on the response and returns the slice at the end, or nil if any key is missing
// or the leaf is not a slice.
func armsNavigate(response map[string]interface{}, path string) []interface{} {
	var cur interface{} = response
	for _, key := range strings.Split(path, ".") {
		m, ok := cur.(map[string]interface{})
		if !ok {
			return nil
		}
		cur = m[key]
	}
	if s, ok := cur.([]interface{}); ok {
		return s
	}
	return nil
}

// armsListPrometheusInstances enumerates prometheus instances in the region.
// ListPrometheusInstances returns Data as a STRINGIFIED JSON array (unusual —
// the HTTP response has "Data":"[{...}]" not "Data":[{...}]), so it must be
// JSON-parsed. Only instances whose regionId matches the sweep region are
// returned (red-line: never clean resources outside the sweep region).
func armsListPrometheusInstances(client *connectivity.AliyunClient) ([]map[string]interface{}, error) {
	response, err := armsRpcPost(client, "ListPrometheusInstances", map[string]interface{}{
		"ShowGlobalView": "false",
		"RegionId":       client.RegionId,
	})
	if err != nil {
		return nil, fmt.Errorf("ListPrometheusInstances: %s", err)
	}
	var raw []interface{}
	switch d := response["Data"].(type) {
	case string:
		if d != "" {
			if e := json.Unmarshal([]byte(d), &raw); e != nil {
				log.Printf("[ERROR] arms ListPrometheusInstances: parse stringified Data: %s", e)
				return nil, nil
			}
		}
	case []interface{}:
		raw = d
	default:
		log.Printf("[ERROR] arms ListPrometheusInstances: unexpected Data type %T", d)
		return nil, nil
	}
	out := make([]map[string]interface{}, 0, len(raw))
	for _, v := range raw {
		m, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		if rid := armsGetField(m, "regionId", "RegionId"); rid != "" && rid != client.RegionId {
			continue // red-line: only the sweep region
		}
		out = append(out, m)
	}
	return out, nil
}

// armsListEnvironments enumerates ARMS environments whose RegionId matches the
// sweep region. ListEnvironments is a GLOBAL list (returns environments across
// regions), so the region filter is mandatory (red-line).
func armsListEnvironments(client *connectivity.AliyunClient) ([]map[string]interface{}, error) {
	response, err := armsRpcPost(client, "ListEnvironments", map[string]interface{}{
		"RegionId": client.RegionId,
	})
	if err != nil {
		return nil, fmt.Errorf("ListEnvironments: %s", err)
	}
	envs := armsNavigate(response, "Data.Environments")
	out := make([]map[string]interface{}, 0, len(envs))
	for _, v := range envs {
		m, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		if rid := armsGetField(m, "RegionId", "regionId"); rid != "" && rid != client.RegionId {
			continue // red-line: only the sweep region
		}
		out = append(out, m)
	}
	return out, nil
}

func init() {
	// environment depends on its env-scoped children so they run first
	resource.AddTestSweepers("alicloud_arms_environment", &resource.Sweeper{
		Name: "alicloud_arms_environment", F: testSweepArmsEnvironments,
		Dependencies: []string{
			"alicloud_arms_addon_release",
			"alicloud_arms_env_custom_job",
			"alicloud_arms_env_feature",
			"alicloud_arms_env_pod_monitor",
			"alicloud_arms_env_service_monitor",
		},
	})
	resource.AddTestSweepers("alicloud_arms_grafana_workspace", &resource.Sweeper{
		Name: "alicloud_arms_grafana_workspace", F: testSweepArmsGrafanaWorkspaces,
	})
	// prometheus depends on its cluster-scoped children so they run first
	resource.AddTestSweepers("alicloud_arms_prometheus", &resource.Sweeper{
		Name: "alicloud_arms_prometheus", F: testSweepArmsPrometheus,
		Dependencies: []string{
			"alicloud_arms_integration_exporter",
			"alicloud_arms_prometheus_alert_rule",
			"alicloud_arms_prometheus_monitoring",
		},
	})
	resource.AddTestSweepers("alicloud_arms_synthetic_task", &resource.Sweeper{
		Name: "alicloud_arms_synthetic_task", F: testSweepArmsSyntheticTasks,
	})
	resource.AddTestSweepers("alicloud_arms_addon_release", &resource.Sweeper{
		Name: "alicloud_arms_addon_release", F: testSweepArmsAddonReleases,
	})
	resource.AddTestSweepers("alicloud_arms_env_custom_job", &resource.Sweeper{
		Name: "alicloud_arms_env_custom_job", F: testSweepArmsEnvCustomJobs,
	})
	resource.AddTestSweepers("alicloud_arms_env_feature", &resource.Sweeper{
		Name: "alicloud_arms_env_feature", F: testSweepArmsEnvFeatures,
	})
	resource.AddTestSweepers("alicloud_arms_env_pod_monitor", &resource.Sweeper{
		Name: "alicloud_arms_env_pod_monitor", F: testSweepArmsEnvPodMonitors,
	})
	resource.AddTestSweepers("alicloud_arms_env_service_monitor", &resource.Sweeper{
		Name: "alicloud_arms_env_service_monitor", F: testSweepArmsEnvServiceMonitors,
	})
	resource.AddTestSweepers("alicloud_arms_integration_exporter", &resource.Sweeper{
		Name: "alicloud_arms_integration_exporter", F: testSweepArmsIntegrationExporters,
	})
	resource.AddTestSweepers("alicloud_arms_prometheus_alert_rule", &resource.Sweeper{
		Name: "alicloud_arms_prometheus_alert_rule", F: testSweepArmsPrometheusAlertRules,
	})
	resource.AddTestSweepers("alicloud_arms_prometheus_monitoring", &resource.Sweeper{
		Name: "alicloud_arms_prometheus_monitoring", F: testSweepArmsPrometheusMonitorings,
	})
}

// ---- region-scoped sweepers ----

// testSweepArmsEnvironments: ListEnvironments (GLOBAL — filter RegionId, red-line)
// -> Data.Environments; delete DeleteEnvironment (EnvironmentId + DeletePromInstance=true).
func testSweepArmsEnvironments(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	envs, err := armsListEnvironments(c)
	if err != nil {
		log.Printf("[ERROR] arms list environments: %s", err)
		return nil
	}
	for _, env := range envs {
		id := armsGetField(env, "EnvironmentId", "environmentId")
		name := armsGetField(env, "EnvironmentName", "environmentName")
		if !armsSweepNameOk(name) && !armsSweepNameOk(id) {
			continue
		}
		if e := armsDeletePost(c, "DeleteEnvironment", map[string]interface{}{
			"EnvironmentId": id, "DeletePromInstance": "true",
		}); e != nil {
			log.Printf("[ERROR] Failed to delete arms environment (%s): %s", id, e)
			continue
		}
		log.Printf("[INFO] Delete arms environment success: %s", id)
	}
	return nil
}

// testSweepArmsGrafanaWorkspaces: ListGrafanaWorkspace -> Data.GrafanaWorkspace;
// delete DeleteGrafanaWorkspace (GrafanaWorkspaceId + RegionId).
func testSweepArmsGrafanaWorkspaces(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	response, err := armsRpcPost(c, "ListGrafanaWorkspace", map[string]interface{}{
		"RegionId": c.RegionId,
	})
	if err != nil {
		log.Printf("[ERROR] arms ListGrafanaWorkspace: %s", err)
		return nil
	}
	for _, v := range armsNavigate(response, "Data.GrafanaWorkspace") {
		item, _ := v.(map[string]interface{})
		id := armsGetField(item, "grafanaWorkspaceId", "GrafanaWorkspaceId")
		name := armsGetField(item, "grafanaWorkspaceName", "GrafanaWorkspaceName")
		if !armsSweepNameOk(name) && !armsSweepNameOk(id) {
			continue
		}
		if e := armsDeletePost(c, "DeleteGrafanaWorkspace", map[string]interface{}{
			"GrafanaWorkspaceId": id, "RegionId": c.RegionId,
		}); e != nil {
			log.Printf("[ERROR] Failed to delete arms grafana workspace (%s): %s", id, e)
			continue
		}
		log.Printf("[INFO] Delete arms grafana workspace success: %s", id)
	}
	return nil
}

// testSweepArmsPrometheus: ListPrometheusInstances (Data is a STRINGIFIED JSON
// array — parsed by armsListPrometheusInstances) -> UninstallPromCluster
// (ClusterId). HIGH VALUE: cleans PAYG test-leftover instances
// (code-test-tf-variable-N, tfaccrealtimecomputeNNNN). UninstallPromCluster is a
// plain PAYG delete (payMethod=post_pay_gb), NOT a BSS refund/退订 — verified.
func testSweepArmsPrometheus(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	instances, err := armsListPrometheusInstances(c)
	if err != nil {
		log.Printf("[ERROR] arms list prometheus instances: %s", err)
		return nil
	}
	for _, inst := range instances {
		id := armsGetField(inst, "clusterId", "ClusterId", "id")
		name := armsGetField(inst, "clusterName", "ClusterName")
		if !armsSweepNameOk(name) && !armsSweepNameOk(id) {
			continue
		}
		if e := armsDeletePost(c, "UninstallPromCluster", map[string]interface{}{
			"ClusterId": id, "RegionId": c.RegionId,
		}); e != nil {
			log.Printf("[ERROR] Failed to delete arms prometheus (%s): %s", id, e)
			continue
		}
		log.Printf("[INFO] Delete arms prometheus success: %s (%s)", id, name)
	}
	return nil
}

// testSweepArmsSyntheticTasks: ListTimingSyntheticTasks — synthetic testing is
// only supported in cn-hangzhou (ParameterRegionIllegal elsewhere), so this
// sweeper is a graceful no-op in cn-beijing (logs + returns nil).
func testSweepArmsSyntheticTasks(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	response, err := armsRpcPost(c, "ListTimingSyntheticTasks", map[string]interface{}{
		"RegionId": c.RegionId,
	})
	if err != nil {
		// ParameterRegionIllegal / region-unsupported is expected outside cn-hangzhou
		log.Printf("[INFO] arms synthetic_task unsupported in %s (no-op): %s", c.RegionId, err)
		return nil
	}
	for _, v := range armsNavigate(response, "Data.Tasks") {
		item, _ := v.(map[string]interface{})
		id := armsGetField(item, "TaskId", "taskId")
		name := armsGetField(item, "TaskName", "taskName")
		if !armsSweepNameOk(name) && !armsSweepNameOk(id) {
			continue
		}
		if e := armsDeletePost(c, "DeleteTimingSyntheticTask", map[string]interface{}{"TaskId": id}); e != nil {
			log.Printf("[ERROR] Failed to delete arms synthetic task (%s): %s", id, e)
			continue
		}
		log.Printf("[INFO] Delete arms synthetic task success: %s", id)
	}
	return nil
}

// ---- environment-scoped sweepers (walk environments first) ----

// armsEnvChildSweep walks ARMS environments in the region, lists the child
// resource under each, and deletes via deleteAction. The extract callback
// returns the child's name (for the sweep gate) and the delete params (without
// RegionId, which is added here). ARMS env-child list responses are
// inconsistent, so several known response paths are tried.
func armsEnvChildSweep(region, listAction, deleteAction string, extract func(item map[string]interface{}, envID string) (name string, delParams map[string]interface{})) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	envs, err := armsListEnvironments(c)
	if err != nil {
		log.Printf("[ERROR] arms list environments for %s: %s", listAction, err)
		return nil
	}
	for _, env := range envs {
		envID := armsGetField(env, "EnvironmentId", "environmentId")
		if envID == "" {
			continue
		}
		response, err := armsRpcPost(c, listAction, map[string]interface{}{
			"EnvironmentId": envID, "RegionId": c.RegionId,
		})
		if err != nil {
			log.Printf("[ERROR] arms %s (%s): %s", listAction, envID, err)
			continue
		}
		// ARMS env-child list responses are inconsistent; try the known paths.
		var items []interface{}
		for _, p := range []string{"Data", "Data.Monitoring", "Data.EnvCustomJobs", "Data.Features", "Data.PodMonitors", "Data.ServiceMonitors", "Data.AddonReleases"} {
			if s := armsNavigate(response, p); len(s) > 0 {
				items = s
				break
			}
		}
		for _, v := range items {
			item, _ := v.(map[string]interface{})
			name, delParams := extract(item, envID)
			if name == "" || !armsSweepNameOk(name) {
				continue
			}
			delParams["RegionId"] = c.RegionId
			if e := armsDeletePost(c, deleteAction, delParams); e != nil {
				log.Printf("[ERROR] Failed to delete arms %s child (%s): %s", listAction, name, e)
				continue
			}
			log.Printf("[INFO] Delete arms %s child success: %s", listAction, name)
		}
	}
	return nil
}

// testSweepArmsAddonReleases: ListAddonReleases (EnvironmentId) ->
// DeleteAddonRelease (EnvironmentId + ReleaseName).
func testSweepArmsAddonReleases(region string) error {
	return armsEnvChildSweep(region, "ListAddonReleases", "DeleteAddonRelease", func(item map[string]interface{}, envID string) (string, map[string]interface{}) {
		name := armsGetField(item, "ReleaseName", "releaseName")
		return name, map[string]interface{}{"EnvironmentId": envID, "ReleaseName": name}
	})
}

// testSweepArmsEnvCustomJobs: ListEnvCustomJobs (EnvironmentId) ->
// DeleteEnvCustomJob (EnvironmentId + CustomJobName).
func testSweepArmsEnvCustomJobs(region string) error {
	return armsEnvChildSweep(region, "ListEnvCustomJobs", "DeleteEnvCustomJob", func(item map[string]interface{}, envID string) (string, map[string]interface{}) {
		name := armsGetField(item, "CustomJobName", "customJobName")
		return name, map[string]interface{}{"EnvironmentId": envID, "CustomJobName": name}
	})
}

// testSweepArmsEnvFeatures: ListEnvironmentFeatures (EnvironmentId) ->
// DeleteEnvironmentFeature (EnvironmentId + FeatureName).
func testSweepArmsEnvFeatures(region string) error {
	return armsEnvChildSweep(region, "ListEnvironmentFeatures", "DeleteEnvironmentFeature", func(item map[string]interface{}, envID string) (string, map[string]interface{}) {
		name := armsGetField(item, "FeatureName", "featureName")
		return name, map[string]interface{}{"EnvironmentId": envID, "FeatureName": name}
	})
}

// testSweepArmsEnvPodMonitors: ListEnvPodMonitors (EnvironmentId) ->
// DeleteEnvPodMonitor (EnvironmentId + Namespace + PodMonitorName).
func testSweepArmsEnvPodMonitors(region string) error {
	return armsEnvChildSweep(region, "ListEnvPodMonitors", "DeleteEnvPodMonitor", func(item map[string]interface{}, envID string) (string, map[string]interface{}) {
		name := armsGetField(item, "PodMonitorName", "podMonitorName")
		ns := armsGetField(item, "Namespace", "namespace")
		return name, map[string]interface{}{"EnvironmentId": envID, "Namespace": ns, "PodMonitorName": name}
	})
}

// testSweepArmsEnvServiceMonitors: ListEnvServiceMonitors (EnvironmentId) ->
// DeleteEnvServiceMonitor (EnvironmentId + Namespace + ServiceMonitorName).
func testSweepArmsEnvServiceMonitors(region string) error {
	return armsEnvChildSweep(region, "ListEnvServiceMonitors", "DeleteEnvServiceMonitor", func(item map[string]interface{}, envID string) (string, map[string]interface{}) {
		name := armsGetField(item, "ServiceMonitorName", "serviceMonitorName")
		ns := armsGetField(item, "Namespace", "namespace")
		return name, map[string]interface{}{"EnvironmentId": envID, "Namespace": ns, "ServiceMonitorName": name}
	})
}

// ---- prometheus/cluster-scoped sweepers (walk prometheus instances first) ----

// testSweepArmsPrometheusAlertRules: ListPrometheusAlertRules (ClusterId) ->
// PrometheusAlertRules; delete DeletePrometheusAlertRule (AlertId ONLY — not
// ClusterId+AlertId).
func testSweepArmsPrometheusAlertRules(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	instances, err := armsListPrometheusInstances(c)
	if err != nil {
		log.Printf("[ERROR] arms list prometheus for alert rule sweep: %s", err)
		return nil
	}
	for _, inst := range instances {
		clusterID := armsGetField(inst, "clusterId", "ClusterId")
		response, err := armsRpcPost(c, "ListPrometheusAlertRules", map[string]interface{}{
			"ClusterId": clusterID, "RegionId": c.RegionId,
		})
		if err != nil {
			log.Printf("[ERROR] arms ListPrometheusAlertRules (%s): %s", clusterID, err)
			continue
		}
		for _, v := range armsNavigate(response, "PrometheusAlertRules") {
			item, _ := v.(map[string]interface{})
			id := armsGetField(item, "AlertId", "alertId")
			name := armsGetField(item, "AlertName", "alertName")
			if !armsSweepNameOk(name) && !armsSweepNameOk(id) {
				continue
			}
			if e := armsDeletePost(c, "DeletePrometheusAlertRule", map[string]interface{}{"AlertId": id}); e != nil {
				log.Printf("[ERROR] Failed to delete arms prometheus alert rule (%s): %s", id, e)
				continue
			}
			log.Printf("[INFO] Delete arms prometheus alert rule success: %s", id)
		}
	}
	return nil
}

// testSweepArmsPrometheusMonitorings: ListPrometheusMonitoring (ClusterId) ->
// Data.Monitoring; delete DeletePrometheusMonitoring (ClusterId + MonitoringName + Type).
func testSweepArmsPrometheusMonitorings(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	instances, err := armsListPrometheusInstances(c)
	if err != nil {
		log.Printf("[ERROR] arms list prometheus for monitoring sweep: %s", err)
		return nil
	}
	for _, inst := range instances {
		clusterID := armsGetField(inst, "clusterId", "ClusterId")
		response, err := armsRpcPost(c, "ListPrometheusMonitoring", map[string]interface{}{
			"ClusterId": clusterID, "RegionId": c.RegionId,
		})
		if err != nil {
			// remote-write-prometheus instances may not support monitoring configs
			continue
		}
		for _, v := range armsNavigate(response, "Data.Monitoring") {
			item, _ := v.(map[string]interface{})
			name := armsGetField(item, "MonitoringName", "monitoringName")
			typ := armsGetField(item, "Type", "type")
			if name == "" {
				continue
			}
			if !armsSweepNameOk(name) {
				continue
			}
			if e := armsDeletePost(c, "DeletePrometheusMonitoring", map[string]interface{}{
				"ClusterId": clusterID, "MonitoringName": name, "Type": typ,
			}); e != nil {
				log.Printf("[ERROR] Failed to delete arms prometheus monitoring (%s/%s): %s", clusterID, name, e)
				continue
			}
			log.Printf("[INFO] Delete arms prometheus monitoring success: %s/%s", clusterID, name)
		}
	}
	return nil
}

// testSweepArmsIntegrationExporters: ListPrometheusIntegration requires
// IntegrationType (kafka/mysql/redis/snmp/emr/nubela/tidb) and only supports
// aliyun-cs/ecs instances — remote-write-prometheus instances are skipped.
// delete DeletePrometheusIntegration (ClusterId + IntegrationType + InstanceId).
func testSweepArmsIntegrationExporters(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	instances, err := armsListPrometheusInstances(c)
	if err != nil {
		log.Printf("[ERROR] arms list prometheus for integration sweep: %s", err)
		return nil
	}
	integrationTypes := []string{"kafka", "mysql", "redis", "snmp", "emr", "nubela", "tidb"}
	for _, inst := range instances {
		clusterType := armsGetField(inst, "clusterType", "ClusterType")
		// ListPrometheusIntegration only supports aliyun-cs/ecs instances
		if clusterType != "aliyun-cs" && clusterType != "ecs" {
			continue
		}
		clusterID := armsGetField(inst, "clusterId", "ClusterId")
		for _, itype := range integrationTypes {
			response, err := armsRpcPost(c, "ListPrometheusIntegration", map[string]interface{}{
				"ClusterId": clusterID, "IntegrationType": itype, "RegionId": c.RegionId,
			})
			if err != nil {
				continue
			}
			for _, v := range armsNavigate(response, "Data.Instance") {
				item, _ := v.(map[string]interface{})
				insID := armsGetField(item, "InstanceId", "instanceId")
				if insID == "" {
					continue
				}
				if e := armsDeletePost(c, "DeletePrometheusIntegration", map[string]interface{}{
					"ClusterId": clusterID, "IntegrationType": itype, "InstanceId": insID,
				}); e != nil {
					log.Printf("[ERROR] Failed to delete arms integration exporter (%s/%s/%s): %s", clusterID, itype, insID, e)
					continue
				}
				log.Printf("[INFO] Delete arms integration exporter success: %s/%s/%s", clusterID, itype, insID)
			}
		}
	}
	return nil
}
