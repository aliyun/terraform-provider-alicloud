package alicloud

// PAI (Platform for AI) sweepers.
//
// 0 sweepers existed in master for any alicloud_pai_* resource (audit
// 2026-10-18). This file 补齐s sweepers for all 12 registered PAI resources:
//
//   alicloud_pai_workspace_workspace        (AIWorkSpace 2021-02-04)
//   alicloud_pai_workspace_dataset          (AIWorkSpace 2021-02-04)
//   alicloud_pai_workspace_datasetversion   (AIWorkSpace 2021-02-04, child of dataset)
//   alicloud_pai_workspace_experiment       (AIWorkSpace 2021-02-04)
//   alicloud_pai_workspace_run              (AIWorkSpace 2021-02-04)
//   alicloud_pai_workspace_code_source      (AIWorkSpace 2021-02-04)
//   alicloud_pai_workspace_member           (AIWorkSpace 2021-02-04, child of workspace)
//   alicloud_pai_workspace_model            (AIWorkSpace 2021-02-04)
//   alicloud_pai_workspace_model_version    (AIWorkSpace 2021-02-04, child of model)
//   alicloud_pai_workspace_user_config      (AIWorkSpace 2021-02-04)
//   alicloud_pai_flow_pipeline             (PAIFlow 2021-02-02)
//   alicloud_pai_service                    (EAS 2021-07-01)
//
// All PAI list endpoints are ROA RESTful paths. The AIWorkSpace / PAIFlow / EAS
// list APIs OMIT the list key entirely when the collection is empty (totalCount=0),
// which would break a naive jsonpath.Get (see rocketmq-sub PR #10925). The
// paiSweepList helper below treats a missing or malformed list key as an empty
// list and returns nil, so an empty account never fails the sweeper.
//
// Name-based NODELETING protection (protectedSweepVpcName) runs at the TOP of
// the per-instance loop in BOTH sweepAll() and normal mode, so do-not-delete
// markers are always honored. The standard `if !sweepAll() { prefix check }`
// gates deletion of non-test resources when SWEEP_ALL is not set. The prefix
// list covers every test naming convention observed in the PAI test files
// (tf_testacc..., tfacc..., tf_testaccworkspace..., tfaccpaiflow...).

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

// paiSweepPrefixes is the prefix list applied when sweepAll() is false. All PAI
// test resources are created with one of these prefixes.
var paiSweepPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf_test_",
	"tf-test-",
	"tfacc",
	"testAcc",
}

// paiNameMatchesPrefix returns true if name starts with any configured prefix
// (case-insensitive). Used by the `if !sweepAll()` gate.
func paiNameMatchesPrefix(name string) bool {
	lower := strings.ToLower(name)
	for _, prefix := range paiSweepPrefixes {
		if strings.HasPrefix(lower, strings.ToLower(prefix)) {
			return true
		}
	}
	return false
}

// paiSweepList calls a ROA GET list endpoint and returns the items at the given
// response key. If the key is absent or the response shape is malformed (the
// empty-list defect pattern: some PAI/EAS list APIs omit the list key entirely
// when totalCount=0), an empty slice is returned and no error is surfaced, so
// an empty account never fails the sweeper.
func paiSweepList(client *connectivity.AliyunClient, product, version, action string, query map[string]*string, listKey string) []map[string]interface{} {
	items := make([]map[string]interface{}, 0)
	response, err := client.RoaGet(product, version, action, query, nil, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve %s list (%s): %s", product, action, err)
		return items
	}
	raw, err := jsonpath.Get(listKey, response)
	if err != nil {
		// Missing list key: treat as empty (empty-list defect guard).
		return items
	}
	list, ok := raw.([]interface{})
	if !ok {
		return items
	}
	for _, v := range list {
		if item, ok := v.(map[string]interface{}); ok {
			items = append(items, item)
		}
	}
	return items
}

// paiSweepSkipName returns true if the resource must be skipped: either because
// the name carries a do-not-delete marker (protectedSweepVpcName, runs in BOTH
// sweepAll and normal mode), or because sweepAll() is false and the name does
// not match any test prefix.
func paiSweepSkipName(name string) (bool, string) {
	if protectedSweepVpcName(name) {
		return true, "do-not-delete name marker"
	}
	if !sweepAll() {
		if !paiNameMatchesPrefix(name) {
			return true, "name does not match test prefixes (non-SWEEP_ALL)"
		}
	}
	return false, ""
}

// paiSweepDeleteAIWorkspace deletes a single AIWorkSpace ROA resource by DELETE
// with optional query params. Errors are logged and swallowed (sweeper returns
// nil even if a per-instance delete failed).
func paiSweepDeleteAIWorkspace(client *connectivity.AliyunClient, action string, query map[string]*string, kind, id string) {
	log.Printf("[INFO] Deleting PAI %s: %s", kind, id)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, err := client.RoaDelete("AIWorkSpace", "2021-02-04", action, query, nil, nil, true)
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
		log.Printf("[ERROR] Failed to delete PAI %s (%s): %s", kind, id, err)
	}
}

func init() {
	resource.AddTestSweepers("alicloud_pai_workspace_workspace", &resource.Sweeper{
		Name: "alicloud_pai_workspace_workspace",
		F:    testSweepPaiWorkspaceWorkspaces,
	})
	resource.AddTestSweepers("alicloud_pai_workspace_dataset", &resource.Sweeper{
		Name: "alicloud_pai_workspace_dataset",
		F:    testSweepPaiWorkspaceDatasets,
	})
	resource.AddTestSweepers("alicloud_pai_workspace_datasetversion", &resource.Sweeper{
		Name: "alicloud_pai_workspace_datasetversion",
		F:    testSweepPaiWorkspaceDatasetversions,
	})
	resource.AddTestSweepers("alicloud_pai_workspace_experiment", &resource.Sweeper{
		Name: "alicloud_pai_workspace_experiment",
		F:    testSweepPaiWorkspaceExperiments,
	})
	resource.AddTestSweepers("alicloud_pai_workspace_run", &resource.Sweeper{
		Name: "alicloud_pai_workspace_run",
		F:    testSweepPaiWorkspaceRuns,
	})
	resource.AddTestSweepers("alicloud_pai_workspace_code_source", &resource.Sweeper{
		Name: "alicloud_pai_workspace_code_source",
		F:    testSweepPaiWorkspaceCodeSources,
	})
	resource.AddTestSweepers("alicloud_pai_workspace_member", &resource.Sweeper{
		Name: "alicloud_pai_workspace_member",
		F:    testSweepPaiWorkspaceMembers,
	})
	resource.AddTestSweepers("alicloud_pai_workspace_model", &resource.Sweeper{
		Name: "alicloud_pai_workspace_model",
		F:    testSweepPaiWorkspaceModels,
	})
	resource.AddTestSweepers("alicloud_pai_workspace_model_version", &resource.Sweeper{
		Name: "alicloud_pai_workspace_model_version",
		F:    testSweepPaiWorkspaceModelVersions,
	})
	resource.AddTestSweepers("alicloud_pai_workspace_user_config", &resource.Sweeper{
		Name: "alicloud_pai_workspace_user_config",
		F:    testSweepPaiWorkspaceUserConfigs,
	})
	resource.AddTestSweepers("alicloud_pai_flow_pipeline", &resource.Sweeper{
		Name: "alicloud_pai_flow_pipeline",
		F:    testSweepPaiFlowPipelines,
	})
	resource.AddTestSweepers("alicloud_pai_service", &resource.Sweeper{
		Name: "alicloud_pai_service",
		F:    testSweepPaiServices,
	})
}

func testSweepPaiWorkspaceWorkspaces(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	items := paiSweepList(client, "AIWorkSpace", "2021-02-04", "/api/v1/workspaces", nil, "Workspaces")
	for _, item := range items {
		name := fmt.Sprint(item["WorkspaceName"])
		id := fmt.Sprint(item["WorkspaceId"])
		if skip, reason := paiSweepSkipName(name); skip {
			log.Printf("[INFO] Skipping PAI Workspace: %s (%s): %s", name, id, reason)
			continue
		}
		paiSweepDeleteAIWorkspace(client, fmt.Sprintf("/api/v1/workspaces/%s", id), nil, "Workspace", id)
	}
	return nil
}

func testSweepPaiWorkspaceDatasets(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	items := paiSweepList(client, "AIWorkSpace", "2021-02-04", "/api/v1/datasets", nil, "Datasets")
	for _, item := range items {
		// The dataset Read sets dataset_name from the "Name" field.
		name := fmt.Sprint(item["Name"])
		id := fmt.Sprint(item["DatasetId"])
		if skip, reason := paiSweepSkipName(name); skip {
			log.Printf("[INFO] Skipping PAI Dataset: %s (%s): %s", name, id, reason)
			continue
		}
		paiSweepDeleteAIWorkspace(client, fmt.Sprintf("/api/v1/datasets/%s", id), nil, "Dataset", id)
	}
	return nil
}

func testSweepPaiWorkspaceDatasetversions(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	// Dataset versions are children of datasets: list all datasets, then for
	// each, list and sweep its versions. The composite ID is DatasetId:VersionName.
	datasets := paiSweepList(client, "AIWorkSpace", "2021-02-04", "/api/v1/datasets", nil, "Datasets")
	for _, ds := range datasets {
		datasetId := fmt.Sprint(ds["DatasetId"])
		if datasetId == "" {
			continue
		}
		versions := paiSweepList(client, "AIWorkSpace", "2021-02-04",
			fmt.Sprintf("/api/v1/datasets/%s/versions", datasetId), nil, "Versions")
		for _, v := range versions {
			versionName := fmt.Sprint(v["VersionName"])
			if skip, reason := paiSweepSkipName(versionName); skip {
				log.Printf("[INFO] Skipping PAI DatasetVersion: %s/%s: %s", datasetId, versionName, reason)
				continue
			}
			paiSweepDeleteAIWorkspace(client,
				fmt.Sprintf("/api/v1/datasets/%s/versions/%s", datasetId, versionName),
				nil, "DatasetVersion", datasetId+":"+versionName)
		}
	}
	return nil
}

func testSweepPaiWorkspaceExperiments(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	items := paiSweepList(client, "AIWorkSpace", "2021-02-04", "/api/v1/experiments", nil, "Experiments")
	for _, item := range items {
		name := fmt.Sprint(item["Name"])
		id := fmt.Sprint(item["ExperimentId"])
		if skip, reason := paiSweepSkipName(name); skip {
			log.Printf("[INFO] Skipping PAI Experiment: %s (%s): %s", name, id, reason)
			continue
		}
		paiSweepDeleteAIWorkspace(client, fmt.Sprintf("/api/v1/experiments/%s", id), nil, "Experiment", id)
	}
	return nil
}

func testSweepPaiWorkspaceRuns(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	items := paiSweepList(client, "AIWorkSpace", "2021-02-04", "/api/v1/runs", nil, "Runs")
	for _, item := range items {
		name := fmt.Sprint(item["Name"])
		id := fmt.Sprint(item["RunId"])
		if skip, reason := paiSweepSkipName(name); skip {
			log.Printf("[INFO] Skipping PAI Run: %s (%s): %s", name, id, reason)
			continue
		}
		paiSweepDeleteAIWorkspace(client, fmt.Sprintf("/api/v1/runs/%s", id), nil, "Run", id)
	}
	return nil
}

func testSweepPaiWorkspaceCodeSources(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	items := paiSweepList(client, "AIWorkSpace", "2021-02-04", "/api/v1/codesources", nil, "CodeSources")
	for _, item := range items {
		// code_source Read uses "DisplayName" as the user-facing name and
		// "CodeSourceName" as the identifier; prefer CodeSourceName for prefix
		// matching (it is the field test configs populate).
		name := fmt.Sprint(item["CodeSourceName"])
		if name == "" {
			name = fmt.Sprint(item["DisplayName"])
		}
		id := fmt.Sprint(item["CodeSourceId"])
		if skip, reason := paiSweepSkipName(name); skip {
			log.Printf("[INFO] Skipping PAI CodeSource: %s (%s): %s", name, id, reason)
			continue
		}
		paiSweepDeleteAIWorkspace(client, fmt.Sprintf("/api/v1/codesources/%s", id), nil, "CodeSource", id)
	}
	return nil
}

func testSweepPaiWorkspaceMembers(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	// Members are children of workspaces and have no human-readable name of
	// their own (user_id and member_id are opaque). Scope deletion by parent
	// workspace: only sweep members of workspaces whose name passes the standard
	// prefix/NODELETING gate. This prevents SWEEP_ALL from stripping real users
	// out of NODELETING or non-test workspaces.
	workspaces := paiSweepList(client, "AIWorkSpace", "2021-02-04", "/api/v1/workspaces", nil, "Workspaces")
	for _, ws := range workspaces {
		workspaceName := fmt.Sprint(ws["WorkspaceName"])
		workspaceId := fmt.Sprint(ws["WorkspaceId"])
		if skip, reason := paiSweepSkipName(workspaceName); skip {
			log.Printf("[INFO] Skipping PAI Member (parent workspace %s (%s) skipped): %s", workspaceName, workspaceId, reason)
			continue
		}
		// The list-members endpoint is /api/v1/workspaces/{id}/members (plural).
		members := paiSweepList(client, "AIWorkSpace", "2021-02-04",
			fmt.Sprintf("/api/v1/workspaces/%s/members", workspaceId), nil, "Members")
		for _, m := range members {
			memberId := fmt.Sprint(m["MemberId"])
			if memberId == "" {
				continue
			}
			// Best-effort NODELETING guard on the member's user_id (opaque, but
			// honors markers if a test ever embeds one).
			if skip, reason := paiSweepSkipName(fmt.Sprint(m["UserId"])); skip {
				log.Printf("[INFO] Skipping PAI Member: %s (workspace %s): %s", memberId, workspaceId, reason)
				continue
			}
			query := map[string]*string{
				"MemberIds": StringPointer(memberId),
			}
			paiSweepDeleteAIWorkspace(client,
				fmt.Sprintf("/api/v1/workspaces/%s/members", workspaceId),
				query, "Member", workspaceId+":"+memberId)
		}
	}
	return nil
}

func testSweepPaiWorkspaceModels(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	items := paiSweepList(client, "AIWorkSpace", "2021-02-04", "/api/v1/models", nil, "Models")
	for _, item := range items {
		name := fmt.Sprint(item["ModelName"])
		id := fmt.Sprint(item["ModelId"])
		if skip, reason := paiSweepSkipName(name); skip {
			log.Printf("[INFO] Skipping PAI Model: %s (%s): %s", name, id, reason)
			continue
		}
		paiSweepDeleteAIWorkspace(client, fmt.Sprintf("/api/v1/models/%s", id), nil, "Model", id)
	}
	return nil
}

func testSweepPaiWorkspaceModelVersions(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	// Model versions are children of models: list all models, then for each,
	// list and sweep its versions. The composite ID is ModelId:VersionName.
	models := paiSweepList(client, "AIWorkSpace", "2021-02-04", "/api/v1/models", nil, "Models")
	for _, model := range models {
		modelId := fmt.Sprint(model["ModelId"])
		if modelId == "" {
			continue
		}
		versions := paiSweepList(client, "AIWorkSpace", "2021-02-04",
			fmt.Sprintf("/api/v1/models/%s/versions", modelId), nil, "Versions")
		for _, v := range versions {
			versionName := fmt.Sprint(v["VersionName"])
			if skip, reason := paiSweepSkipName(versionName); skip {
				log.Printf("[INFO] Skipping PAI ModelVersion: %s/%s: %s", modelId, versionName, reason)
				continue
			}
			paiSweepDeleteAIWorkspace(client,
				fmt.Sprintf("/api/v1/models/%s/versions/%s", modelId, versionName),
				nil, "ModelVersion", modelId+":"+versionName)
		}
	}
	return nil
}

func testSweepPaiWorkspaceUserConfigs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	// The user_config ID is CategoryName:ConfigKey; the DELETE endpoint is
	// /api/v1/userconfigs/{CategoryName}?ConfigKey={ConfigKey}. Use CategoryName
	// as the name for prefix/NODELETING gating.
	items := paiSweepList(client, "AIWorkSpace", "2021-02-04", "/api/v1/userconfigs", nil, "UserConfigs")
	for _, item := range items {
		categoryName := fmt.Sprint(item["CategoryName"])
		configKey := fmt.Sprint(item["ConfigKey"])
		if skip, reason := paiSweepSkipName(categoryName); skip {
			log.Printf("[INFO] Skipping PAI UserConfig: %s/%s: %s", categoryName, configKey, reason)
			continue
		}
		query := map[string]*string{
			"ConfigKey": StringPointer(configKey),
		}
		paiSweepDeleteAIWorkspace(client,
			fmt.Sprintf("/api/v1/userconfigs/%s", categoryName),
			query, "UserConfig", categoryName+":"+configKey)
	}
	return nil
}

func testSweepPaiFlowPipelines(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	// PAIFlow uses its own product code and version (PAIFlow 2021-02-02).
	items := paiSweepList(client, "PAIFlow", "2021-02-02", "/api/v1/pipelines", nil, "Pipelines")
	for _, item := range items {
		name := fmt.Sprint(item["PipelineName"])
		id := fmt.Sprint(item["PipelineId"])
		if skip, reason := paiSweepSkipName(name); skip {
			log.Printf("[INFO] Skipping PAI Flow Pipeline: %s (%s): %s", name, id, reason)
			continue
		}
		log.Printf("[INFO] Deleting PAI Flow Pipeline: %s (%s)", name, id)
		action := fmt.Sprintf("/api/v1/pipelines/%s", id)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err := retry.Retry(5*time.Minute, func() *retry.RetryError {
			_, err := client.RoaDelete("PAIFlow", "2021-02-02", action, nil, nil, nil, true)
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
			log.Printf("[ERROR] Failed to delete PAI Flow Pipeline (%s): %s", id, err)
		}
	}
	return nil
}

func testSweepPaiServices(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	// EAS services are addressed by ClusterId (= region) + ServiceName. The
	// list endpoint is GET /api/v2/services?ClusterId={region}. The response
	// shape varies across EAS API versions; paiSweepList treats a missing list
	// key as empty (empty-list defect guard). Try "Services" first, then fall
	// back to a top-level array shape.
	clusterId := client.RegionId
	query := map[string]*string{
		"ClusterId": StringPointer(clusterId),
	}
	items := paiSweepList(client, "EAS", "2021-07-01", "/api/v2/services", query, "Services")
	for _, item := range items {
		name := fmt.Sprint(item["ServiceName"])
		if name == "" {
			continue
		}
		if skip, reason := paiSweepSkipName(name); skip {
			log.Printf("[INFO] Skipping PAI Service: %s: %s", name, reason)
			continue
		}
		log.Printf("[INFO] Deleting PAI Service: %s", name)
		action := fmt.Sprintf("/api/v2/services/%s/%s", clusterId, name)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err := retry.Retry(5*time.Minute, func() *retry.RetryError {
			_, err := client.RoaDelete("EAS", "2021-07-01", action, query, nil, nil, true)
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
			log.Printf("[ERROR] Failed to delete PAI Service (%s): %s", name, err)
		}
	}
	return nil
}
