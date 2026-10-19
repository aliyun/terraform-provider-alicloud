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

// GPDB (AnalyticDB for PostgreSQL, product "gpdb" / "2016-05-03") child-resource
// sweepers. Nine child resources are swept:
//   - instance-scoped (walk Running instances first): account, db_instance_ip_array,
//     db_instance_plan, remote_adb_data_source, streaming_data_source, streaming_job,
//     streaming_data_service, external_data_service
//   - region-scoped (no instance walk): supabase_project
//
// The instance sweepers (alicloud_gpdb_instance + alicloud_gpdb_elastic_instance)
// already exist and are intentionally NOT re-run here: DeleteDBInstance is a refund
// (退订) for Prepaid instances and needs explicit approval (red-line).
//
// Seven gap resources are NOT sweepable (describe-one Read with no list-all API on
// the cloud side, cascade-cleaned by instance deletion): database, db_resource_group,
// hadoop_data_source, jdbc_data_source, connection, backup_policy, db_extension.
// One resource (api_key) is sweepable but workspace-scoped (ListApiKeys requires
// WorkspaceId; ListWorkspaces currently returns 0) and is deferred.

var gpdbSweeperPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf-example",
	"tf_example",
	"terraform-example",
	"tftestacc",
}

// gpdbGetField returns the first non-empty value among the given case-sensitive
// keys on item. GPDB responses use consistent PascalCase, but a few actions
// (ListApiKeys etc.) return camelCase; this helper tolerates both.
func gpdbGetField(item map[string]interface{}, names ...string) string {
	for _, name := range names {
		if v, ok := item[name]; ok {
			if s := fmt.Sprint(v); s != "" && s != "<nil>" {
				return s
			}
		}
	}
	return ""
}

// gpdbSweepNameOk applies the standard fleet-wipe gate: never touch NODELETING
// names, and in normal (non-SWEEP_ALL) mode only sweep resources whose name
// matches one of the tf-test prefixes. SWEEP_ALL bypasses the prefix filter so
// the whole test account (1511928242963727, cn-beijing) is fleet-wiped.
func gpdbSweepNameOk(name string) bool {
	if name == "" {
		return false
	}
	if protectedSweepVpcName(name) {
		log.Printf("[INFO] Skipping gpdb resource (NODELETING): %s", name)
		return false
	}
	if !sweepAll() {
		skip := true
		for _, prefix := range gpdbSweeperPrefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping gpdb resource (no prefix): %s", name)
			return false
		}
	}
	return true
}

// gpdbListItems issues a single (non-paginated) RpcPost list call and returns
// the slice at itemsPath. Most gpdb child list actions return all items in one
// response (or a small bounded set on a test account), so a single page is
// sufficient. A missing path is treated as an empty list (logged) so the
// sweeper continues with zero items rather than aborting.
func gpdbListItems(client *connectivity.AliyunClient, action string, request map[string]interface{}, itemsPath string) ([]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, nil, request, true)
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
		return nil, fmt.Errorf("%s: %s", action, err)
	}
	resp, err := jsonpath.Get(itemsPath, response)
	if err != nil {
		log.Printf("[ERROR] gpdb %s jsonpath %s: %s", action, itemsPath, err)
		return nil, nil
	}
	result, _ := resp.([]interface{})
	return result, nil
}

// listGpdbRunningInstanceIDs enumerates every Running gpdb instance in the
// region via DescribeDBInstances (PageNumber/PageSize pagination). Child
// resources live on Running instances; Released instances have already
// cascade-deleted their children.
func listGpdbRunningInstanceIDs(client *connectivity.AliyunClient) ([]string, error) {
	action := "DescribeDBInstances"
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	instanceIDs := make([]string, 0)
	for {
		var response map[string]interface{}
		var err error
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("gpdb", "2016-05-03", action, nil, request, true)
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
			return nil, fmt.Errorf("DescribeDBInstances: %s", err)
		}
		resp, err := jsonpath.Get("$.Items.DBInstance", response)
		if err != nil {
			return nil, fmt.Errorf("DescribeDBInstances jsonpath $.Items.DBInstance: %s", err)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				if gpdbGetField(item, "DBInstanceStatus") == "Running" {
					if id := gpdbGetField(item, "DBInstanceId"); id != "" {
						instanceIDs = append(instanceIDs, id)
					}
				}
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return instanceIDs, nil
}

func init() {
	resource.AddTestSweepers("alicloud_gpdb_account", &resource.Sweeper{
		Name: "alicloud_gpdb_account",
		F:    testSweepGpdbAccounts,
	})
	resource.AddTestSweepers("alicloud_gpdb_db_instance_ip_array", &resource.Sweeper{
		Name: "alicloud_gpdb_db_instance_ip_array",
		F:    testSweepGpdbDBInstanceIPArrays,
	})
	resource.AddTestSweepers("alicloud_gpdb_db_instance_plan", &resource.Sweeper{
		Name: "alicloud_gpdb_db_instance_plan",
		F:    testSweepGpdbDBInstancePlans,
	})
	resource.AddTestSweepers("alicloud_gpdb_remote_adb_data_source", &resource.Sweeper{
		Name: "alicloud_gpdb_remote_adb_data_source",
		F:    testSweepGpdbRemoteADBDataSources,
	})
	resource.AddTestSweepers("alicloud_gpdb_streaming_data_source", &resource.Sweeper{
		Name: "alicloud_gpdb_streaming_data_source",
		F:    testSweepGpdbStreamingDataSources,
	})
	resource.AddTestSweepers("alicloud_gpdb_streaming_job", &resource.Sweeper{
		Name: "alicloud_gpdb_streaming_job",
		F:    testSweepGpdbStreamingJobs,
	})
	resource.AddTestSweepers("alicloud_gpdb_streaming_data_service", &resource.Sweeper{
		Name: "alicloud_gpdb_streaming_data_service",
		F:    testSweepGpdbStreamingDataServices,
	})
	resource.AddTestSweepers("alicloud_gpdb_external_data_service", &resource.Sweeper{
		Name: "alicloud_gpdb_external_data_service",
		F:    testSweepGpdbExternalDataServices,
	})
	resource.AddTestSweepers("alicloud_gpdb_supabase_project", &resource.Sweeper{
		Name: "alicloud_gpdb_supabase_project",
		F:    testSweepGpdbSupabaseProjects,
	})
}

// testSweepGpdbAccounts sweeps gpdb accounts (DescribeAccounts -> DeleteAccount)
// on every Running instance. Accounts are cascade-deleted with their instance,
// so this is defense-in-depth against leaks on shared instances.
func testSweepGpdbAccounts(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	instanceIDs, err := listGpdbRunningInstanceIDs(aliyunClient)
	if err != nil {
		log.Printf("[ERROR] list gpdb instances for account sweep: %s", err)
		return nil
	}
	for _, instID := range instanceIDs {
		req := map[string]interface{}{
			"RegionId":     aliyunClient.RegionId,
			"DBInstanceId": instID,
		}
		items, err := gpdbListItems(aliyunClient, "DescribeAccounts", req, "$.Accounts.DBInstanceAccount")
		if err != nil {
			log.Printf("[ERROR] DescribeAccounts (%s): %s", instID, err)
			continue
		}
		for _, v := range items {
			item, _ := v.(map[string]interface{})
			name := gpdbGetField(item, "AccountName")
			if !gpdbSweepNameOk(name) {
				continue
			}
			delReq := map[string]interface{}{
				"DBInstanceId": instID,
				"AccountName":  name,
				"RegionId":     aliyunClient.RegionId,
			}
			if _, e := aliyunClient.RpcPost("gpdb", "2016-05-03", "DeleteAccount", nil, delReq, true); e != nil {
				log.Printf("[ERROR] Failed to delete gpdb account (%s/%s): %s", instID, name, e)
				continue
			}
			log.Printf("[INFO] Delete gpdb account success: %s/%s", instID, name)
		}
	}
	return nil
}

// testSweepGpdbDBInstanceIPArrays sweeps gpdb IP whitelist groups
// (DescribeDBInstanceIPArrayList -> DeleteDBInstanceIPArray). The built-in
// "default" group is preserved (system-managed); only user-created groups
// (e.g. tf-test group1/group2 leaks) are deleted. Note the delete param is
// IPArrayName, not DBInstanceIPArrayName.
func testSweepGpdbDBInstanceIPArrays(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	instanceIDs, err := listGpdbRunningInstanceIDs(aliyunClient)
	if err != nil {
		log.Printf("[ERROR] list gpdb instances for ip array sweep: %s", err)
		return nil
	}
	for _, instID := range instanceIDs {
		req := map[string]interface{}{
			"RegionId":     aliyunClient.RegionId,
			"DBInstanceId": instID,
		}
		items, err := gpdbListItems(aliyunClient, "DescribeDBInstanceIPArrayList", req, "$.Items.DBInstanceIPArray")
		if err != nil {
			log.Printf("[ERROR] DescribeDBInstanceIPArrayList (%s): %s", instID, err)
			continue
		}
		for _, v := range items {
			item, _ := v.(map[string]interface{})
			name := gpdbGetField(item, "DBInstanceIPArrayName")
			if name == "default" {
				log.Printf("[INFO] Skipping gpdb ip array (system built-in): %s/default", instID)
				continue
			}
			if !gpdbSweepNameOk(name) {
				continue
			}
			delReq := map[string]interface{}{
				"DBInstanceId": instID,
				"IPArrayName":  name,
				"RegionId":     aliyunClient.RegionId,
			}
			if _, e := aliyunClient.RpcPost("gpdb", "2016-05-03", "DeleteDBInstanceIPArray", nil, delReq, true); e != nil {
				log.Printf("[ERROR] Failed to delete gpdb ip array (%s/%s): %s", instID, name, e)
				continue
			}
			log.Printf("[INFO] Delete gpdb ip array success: %s/%s", instID, name)
		}
	}
	return nil
}

// testSweepGpdbDBInstancePlans sweeps gpdb instance plans
// (DescribeDBInstancePlans -> DeleteDBInstancePlan) on every Running instance.
func testSweepGpdbDBInstancePlans(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	instanceIDs, err := listGpdbRunningInstanceIDs(aliyunClient)
	if err != nil {
		log.Printf("[ERROR] list gpdb instances for plan sweep: %s", err)
		return nil
	}
	for _, instID := range instanceIDs {
		req := map[string]interface{}{
			"RegionId":     aliyunClient.RegionId,
			"DBInstanceId": instID,
		}
		items, err := gpdbListItems(aliyunClient, "DescribeDBInstancePlans", req, "$.Items.PlanList")
		if err != nil {
			log.Printf("[ERROR] DescribeDBInstancePlans (%s): %s", instID, err)
			continue
		}
		for _, v := range items {
			item, _ := v.(map[string]interface{})
			name := gpdbGetField(item, "PlanName")
			if !gpdbSweepNameOk(name) {
				continue
			}
			planID := gpdbGetField(item, "PlanId")
			if planID == "" {
				continue
			}
			delReq := map[string]interface{}{
				"DBInstanceId": instID,
				"PlanId":       planID,
				"RegionId":     aliyunClient.RegionId,
			}
			if _, e := aliyunClient.RpcPost("gpdb", "2016-05-03", "DeleteDBInstancePlan", nil, delReq, true); e != nil {
				log.Printf("[ERROR] Failed to delete gpdb plan (%s/%s): %s", instID, name, e)
				continue
			}
			log.Printf("[INFO] Delete gpdb plan success: %s/%s", instID, name)
		}
	}
	return nil
}

// testSweepGpdbRemoteADBDataSources sweeps remote ADB data sources
// (ListRemoteADBDataSources -> DeleteRemoteADBDataSource) on every Running
// instance. Note the delete param is LocalDBInstanceId, not DBInstanceId.
func testSweepGpdbRemoteADBDataSources(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	instanceIDs, err := listGpdbRunningInstanceIDs(aliyunClient)
	if err != nil {
		log.Printf("[ERROR] list gpdb instances for remote adb sweep: %s", err)
		return nil
	}
	for _, instID := range instanceIDs {
		req := map[string]interface{}{
			"RegionId":     aliyunClient.RegionId,
			"DBInstanceId": instID,
		}
		items, err := gpdbListItems(aliyunClient, "ListRemoteADBDataSources", req, "$.DataSourceItems.RemoteDataSources")
		if err != nil {
			log.Printf("[ERROR] ListRemoteADBDataSources (%s): %s", instID, err)
			continue
		}
		for _, v := range items {
			item, _ := v.(map[string]interface{})
			name := gpdbGetField(item, "DataSourceName")
			if !gpdbSweepNameOk(name) {
				continue
			}
			dsID := gpdbGetField(item, "DataSourceId")
			if dsID == "" {
				continue
			}
			delReq := map[string]interface{}{
				"LocalDBInstanceId": instID,
				"DataSourceId":      dsID,
				"RegionId":          aliyunClient.RegionId,
			}
			if _, e := aliyunClient.RpcPost("gpdb", "2016-05-03", "DeleteRemoteADBDataSource", nil, delReq, true); e != nil {
				log.Printf("[ERROR] Failed to delete gpdb remote adb data source (%s/%s): %s", instID, name, e)
				continue
			}
			log.Printf("[INFO] Delete gpdb remote adb data source success: %s/%s", instID, name)
		}
	}
	return nil
}

// testSweepGpdbStreamingDataSources sweeps gpdb streaming data sources
// (ListStreamingDataSources -> DeleteStreamingDataSource) on every Running
// instance. Instance-scoped (DBInstanceId required for both list and delete).
func testSweepGpdbStreamingDataSources(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	instanceIDs, err := listGpdbRunningInstanceIDs(aliyunClient)
	if err != nil {
		log.Printf("[ERROR] list gpdb instances for streaming data source sweep: %s", err)
		return nil
	}
	for _, instID := range instanceIDs {
		req := map[string]interface{}{
			"RegionId":     aliyunClient.RegionId,
			"DBInstanceId": instID,
		}
		items, err := gpdbListItems(aliyunClient, "ListStreamingDataSources", req, "$.DataSourceItems")
		if err != nil {
			log.Printf("[ERROR] ListStreamingDataSources (%s): %s", instID, err)
			continue
		}
		for _, v := range items {
			item, _ := v.(map[string]interface{})
			name := gpdbGetField(item, "DataSourceName")
			if !gpdbSweepNameOk(name) {
				continue
			}
			dsID := gpdbGetField(item, "DataSourceId")
			if dsID == "" {
				continue
			}
			delReq := map[string]interface{}{
				"DBInstanceId": instID,
				"DataSourceId": dsID,
				"RegionId":     aliyunClient.RegionId,
			}
			if _, e := aliyunClient.RpcPost("gpdb", "2016-05-03", "DeleteStreamingDataSource", nil, delReq, true); e != nil {
				log.Printf("[ERROR] Failed to delete gpdb streaming data source (%s/%s): %s", instID, name, e)
				continue
			}
			log.Printf("[INFO] Delete gpdb streaming data source success: %s/%s", instID, name)
		}
	}
	return nil
}

// testSweepGpdbStreamingJobs sweeps gpdb streaming jobs
// (ListStreamingJobs -> DeleteStreamingJob) on every Running instance.
// Instance-scoped (DBInstanceId required for both list and delete).
func testSweepGpdbStreamingJobs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	instanceIDs, err := listGpdbRunningInstanceIDs(aliyunClient)
	if err != nil {
		log.Printf("[ERROR] list gpdb instances for streaming job sweep: %s", err)
		return nil
	}
	for _, instID := range instanceIDs {
		req := map[string]interface{}{
			"RegionId":     aliyunClient.RegionId,
			"DBInstanceId": instID,
		}
		items, err := gpdbListItems(aliyunClient, "ListStreamingJobs", req, "$.JobItems")
		if err != nil {
			log.Printf("[ERROR] ListStreamingJobs (%s): %s", instID, err)
			continue
		}
		for _, v := range items {
			item, _ := v.(map[string]interface{})
			name := gpdbGetField(item, "JobName")
			if !gpdbSweepNameOk(name) {
				continue
			}
			jobID := gpdbGetField(item, "JobId")
			if jobID == "" {
				continue
			}
			delReq := map[string]interface{}{
				"DBInstanceId": instID,
				"JobId":        jobID,
				"RegionId":     aliyunClient.RegionId,
			}
			if _, e := aliyunClient.RpcPost("gpdb", "2016-05-03", "DeleteStreamingJob", nil, delReq, true); e != nil {
				log.Printf("[ERROR] Failed to delete gpdb streaming job (%s/%s): %s", instID, name, e)
				continue
			}
			log.Printf("[INFO] Delete gpdb streaming job success: %s/%s", instID, name)
		}
	}
	return nil
}

// testSweepGpdbStreamingDataServices sweeps gpdb streaming data services
// (ListStreamingDataServices -> DeleteStreamingDataService) on every Running
// instance. Instance-scoped (DBInstanceId required for both list and delete).
func testSweepGpdbStreamingDataServices(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	instanceIDs, err := listGpdbRunningInstanceIDs(aliyunClient)
	if err != nil {
		log.Printf("[ERROR] list gpdb instances for streaming data service sweep: %s", err)
		return nil
	}
	for _, instID := range instanceIDs {
		req := map[string]interface{}{
			"RegionId":     aliyunClient.RegionId,
			"DBInstanceId": instID,
		}
		items, err := gpdbListItems(aliyunClient, "ListStreamingDataServices", req, "$.ServiceItems")
		if err != nil {
			log.Printf("[ERROR] ListStreamingDataServices (%s): %s", instID, err)
			continue
		}
		for _, v := range items {
			item, _ := v.(map[string]interface{})
			name := gpdbGetField(item, "ServiceName")
			if !gpdbSweepNameOk(name) {
				continue
			}
			svcID := gpdbGetField(item, "ServiceId")
			if svcID == "" {
				continue
			}
			delReq := map[string]interface{}{
				"DBInstanceId": instID,
				"ServiceId":    svcID,
				"RegionId":     aliyunClient.RegionId,
			}
			if _, e := aliyunClient.RpcPost("gpdb", "2016-05-03", "DeleteStreamingDataService", nil, delReq, true); e != nil {
				log.Printf("[ERROR] Failed to delete gpdb streaming data service (%s/%s): %s", instID, name, e)
				continue
			}
			log.Printf("[INFO] Delete gpdb streaming data service success: %s/%s", instID, name)
		}
	}
	return nil
}

// testSweepGpdbExternalDataServices sweeps gpdb external data services
// (ListExternalDataServices -> DeleteExternalDataService) on every Running
// instance. Instance-scoped (DBInstanceId required for both list and delete).
func testSweepGpdbExternalDataServices(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	instanceIDs, err := listGpdbRunningInstanceIDs(aliyunClient)
	if err != nil {
		log.Printf("[ERROR] list gpdb instances for external data service sweep: %s", err)
		return nil
	}
	for _, instID := range instanceIDs {
		req := map[string]interface{}{
			"RegionId":     aliyunClient.RegionId,
			"DBInstanceId": instID,
		}
		items, err := gpdbListItems(aliyunClient, "ListExternalDataServices", req, "$.ServiceItems")
		if err != nil {
			log.Printf("[ERROR] ListExternalDataServices (%s): %s", instID, err)
			continue
		}
		for _, v := range items {
			item, _ := v.(map[string]interface{})
			name := gpdbGetField(item, "ServiceName")
			if !gpdbSweepNameOk(name) {
				continue
			}
			svcID := gpdbGetField(item, "ServiceId")
			if svcID == "" {
				continue
			}
			delReq := map[string]interface{}{
				"DBInstanceId": instID,
				"ServiceId":    svcID,
				"RegionId":     aliyunClient.RegionId,
			}
			if _, e := aliyunClient.RpcPost("gpdb", "2016-05-03", "DeleteExternalDataService", nil, delReq, true); e != nil {
				log.Printf("[ERROR] Failed to delete gpdb external data service (%s/%s): %s", instID, name, e)
				continue
			}
			log.Printf("[INFO] Delete gpdb external data service success: %s/%s", instID, name)
		}
	}
	return nil
}

// testSweepGpdbSupabaseProjects sweeps gpdb Supabase projects
// (ListSupabaseProjects -> DeleteSupabaseProject). Region-scoped.
func testSweepGpdbSupabaseProjects(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	req := map[string]interface{}{"RegionId": aliyunClient.RegionId}
	items, err := gpdbListItems(aliyunClient, "ListSupabaseProjects", req, "$.Items")
	if err != nil {
		log.Printf("[ERROR] ListSupabaseProjects: %s", err)
		return nil
	}
	for _, v := range items {
		item, _ := v.(map[string]interface{})
		name := gpdbGetField(item, "ProjectName")
		if !gpdbSweepNameOk(name) {
			continue
		}
		projectID := gpdbGetField(item, "ProjectId")
		if projectID == "" {
			continue
		}
		delReq := map[string]interface{}{
			"ProjectId": projectID,
			"RegionId":  aliyunClient.RegionId,
		}
		if _, e := aliyunClient.RpcPost("gpdb", "2016-05-03", "DeleteSupabaseProject", nil, delReq, true); e != nil {
			log.Printf("[ERROR] Failed to delete gpdb supabase project (%s): %s", name, e)
			continue
		}
		log.Printf("[INFO] Delete gpdb supabase project success: %s", name)
	}
	return nil
}
